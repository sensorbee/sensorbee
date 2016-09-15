package server

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"gopkg.in/pfnet/jasco.v1"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"gopkg.in/sensorbee/sensorbee.v0/server/config"
	"gopkg.in/sensorbee/sensorbee.v0/server/udsstorage"
	"io"
	"io/ioutil"
)

// Context is a context object for gocraft/web.
type Context struct {
	*jasco.Context

	udsStorage udf.UDSStorage
	topologies TopologyRegistry
	config     *config.Config
	// logger is used by core.Context, not for the server's Context. This logger
	// can be shared with jasco.Context.
	logger *logrus.Logger
}

// SetTopologyRegistry sets the registry of topologies to this context. This
// method must be called in the middleware of Context.
func (c *Context) SetTopologyRegistry(r TopologyRegistry) {
	c.topologies = r
}

// ContextGlobalVariables has fields which are shared through all contexts
// allocated for each request.
type ContextGlobalVariables struct {
	// Logger is used to write log messages.
	Logger *logrus.Logger

	// LogDestination is a writer to which logs are written.
	LogDestination io.WriteCloser

	// Topologies is a registry which manages topologies to support multi
	// tenancy.
	Topologies TopologyRegistry

	// Config has configuration parameters.
	Config *config.Config
}

// SetUpContextGlobalVariables create a new ContextGlobalVariables from a config.
// DO NOT make any change on the config after calling this function. The caller
// can change other members of ContextGlobalVariables.
//
// The caller must Close LogDestination.
func SetUpContextGlobalVariables(conf *config.Config) (*ContextGlobalVariables, error) {
	logger := logrus.New()
	logLevel, err := logrus.ParseLevel(conf.Logging.MinLogLevel)
	if err != nil {
		return nil, err
	}
	logger.Level = logLevel
	w, err := conf.Logging.CreateWriter()
	if err != nil {
		return nil, err
	}
	closeWriter := true
	defer func() {
		if closeWriter {
			w.Close()
		}
	}()
	logger.Out = w

	closeWriter = false
	return &ContextGlobalVariables{
		Logger:         logger,
		LogDestination: w,
		Topologies:     NewDefaultTopologyRegistry(),
		Config:         conf,
	}, nil
}

// SetUpContextAndRouter creates a router of the API server and its context.
// jascoRoot is a root router returned from jasco.New.
//
// This function returns a new web.Router. Don't use the router returned from
// this function as a handler of HTTP server, but use jascoRoot instead.
func SetUpContextAndRouter(prefix string, jascoRoot *web.Router, gvariables *ContextGlobalVariables) (*web.Router, error) {
	gvars := *gvariables
	udsStorage, err := setUpUDSStorage(&gvars.Config.Storage.UDS)
	if err != nil {
		return nil, err
	}

	// Topologies should be created after setting up everything necessary for it.
	if err := setUpTopologies(gvars.Logger, gvars.Topologies, gvars.Config, udsStorage); err != nil {
		return nil, err
	}

	router := jascoRoot.Subrouter(Context{}, "/")
	router.Middleware(func(c *Context, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		c.logger = gvars.Logger
		c.udsStorage = udsStorage
		c.topologies = gvars.Topologies
		c.config = gvars.Config
		next(rw, req)
	})
	return router, nil
}

func setUpUDSStorage(conf *config.UDSStorage) (udf.UDSStorage, error) {
	// Parameters are already validated in conf
	switch conf.Type {
	case "in_memory":
		return udf.NewInMemoryUDSStorage(), nil
	case "fs":
		dir, _ := data.AsString(conf.Params["dir"])
		var tempDir string
		if v, ok := conf.Params["temp_dir"]; ok {
			tempDir, _ = data.AsString(v)
		}
		return udsstorage.NewFS(dir, tempDir)
	default:
		return nil, fmt.Errorf("unsupported uds storage type: %v", conf.Type)
	}
}

func setUpTopologies(logger *logrus.Logger, r TopologyRegistry, conf *config.Config, us udf.UDSStorage) error {
	stopAll := true
	defer func() {
		if stopAll {
			ts, err := r.List()
			if err != nil {
				logger.WithField("err", err).Error("Cannot list topologies for clean up")
				return
			}

			for name, t := range ts {
				if err := t.Topology().Stop(); err != nil {
					logger.WithFields(logrus.Fields{
						"err":      err,
						"topology": name,
					}).Error("Cannot stop the topology")
				}
			}
		}
	}()

	for name := range conf.Topologies {
		logger.WithField("topology", name).Info("Setting up the topology")
		tb, err := setUpTopology(logger, name, conf, us)
		if err != nil {
			return err
		}
		if err := r.Register(name, tb); err != nil {
			logger.WithFields(logrus.Fields{
				"err":      err,
				"topology": name,
			}).Error("Cannot register the topology")
			return err
		}
	}

	stopAll = false
	return nil
}

func setUpTopology(logger *logrus.Logger, name string, conf *config.Config, us udf.UDSStorage) (*bql.TopologyBuilder, error) {
	cc := &core.ContextConfig{
		Logger: logger,
	}
	cc.Flags.DroppedTupleLog.Set(conf.Logging.LogDroppedTuples)
	cc.Flags.DroppedTupleSummarization.Set(conf.Logging.SummarizeDroppedTuples)

	tp, err := core.NewDefaultTopology(core.NewContext(cc), name)
	if err != nil {
		return nil, err
	}
	tb, err := bql.NewTopologyBuilder(tp)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"topology": name,
		}).Error("Cannot create a topology builder")
		return nil, err
	}
	tb.UDSStorage = us

	bqlFilePath := conf.Topologies[name].BQLFile
	if bqlFilePath == "" {
		return tb, nil
	}

	shouldStop := true
	defer func() {
		if shouldStop {
			if err := tp.Stop(); err != nil {
				logger.WithFields(logrus.Fields{
					"err":      err,
					"topology": name,
				}).Error("Cannot stop the topology")
			}
		}
	}()

	queries, err := ioutil.ReadFile(bqlFilePath)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"topology": name,
			"path":     bqlFilePath,
		}).Error("Cannot read a BQL file")
		return nil, err
	}

	// TODO: improve error handling
	bp := parser.New()
	stmts, err := bp.ParseStmts(string(queries))
	if err != nil {
		return nil, err
	}

	for _, stmt := range stmts {
		if _, err := tb.AddStmt(stmt); err != nil {
			logger.WithFields(logrus.Fields{
				"err":      err,
				"topology": name,
				"stmt":     stmt,
			}).Error("Cannot add a statement to the topology")
			return nil, err
		}
	}

	shouldStop = false
	return tb, nil
}

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"pfi/sensorbee/sensorbee/server/config"
	"pfi/sensorbee/sensorbee/server/udsstorage"
	"runtime"
	"sync/atomic"
	"time"
)

// Context is a context object for gocraft/web.
type Context struct {
	body      []byte
	bodyError error

	requestID uint64

	response   web.ResponseWriter
	request    *web.Request
	HTTPStatus int

	logger     *logrus.Logger
	udsStorage udf.UDSStorage
	topologies TopologyRegistry
	config     *config.Config
}

// SetLogger sets the logger to the context. Must be set before any action
// is invoked.
func (c *Context) SetLogger(l *logrus.Logger) {
	c.logger = l
}

var requestIDCounter uint64

func (c *Context) setUpContext(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.requestID = atomic.AddUint64(&requestIDCounter, 1)
	c.response = rw
	c.request = req

	start := time.Now()
	defer func() {
		elapsed := time.Now().Sub(start)
		// Use custom logging because file and line aren't necessary here.
		c.logger.WithFields(logrus.Fields{
			"reqid":   c.requestID,
			"reqtime": fmt.Sprintf("%d.%03d", int(elapsed/time.Second), int(elapsed%time.Second/time.Millisecond)),
			"method":  req.Method,
			"uri":     req.URL.RequestURI(),
			"status":  c.HTTPStatus,
		}).Info("Access")
	}()
	next(rw, req)
}

// SetTopologyRegistry sets the registry of topologies to this context. This
// method must be called in the middleware of Context.
func (c *Context) SetTopologyRegistry(r TopologyRegistry) {
	c.topologies = r
}

// NotFoundHandler handles 404.
func (c *Context) NotFoundHandler(rw web.ResponseWriter, req *web.Request) {
	c.logger.WithFields(logrus.Fields{
		"method": req.Method,
		"uri":    req.URL.RequestURI(),
		"status": http.StatusNotFound,
	}).Error("The request URL not found")

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(`
{
  "error": {
    "code": "E0001",
    "message": "The request URL was not found."
  }
}`))
}

// Log returns the logger having meta information.
func (c *Context) Log() *logrus.Entry {
	return c.log(1)
}

// ErrLog returns the logger with error information.
func (c *Context) ErrLog(err error) *logrus.Entry {
	return c.log(1).WithField("err", err)
}

func (c *Context) log(depth int) *logrus.Entry {
	// TODO: This is a temporary solution until logrus support filename and line number
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		return c.logger.WithField("reqid", c.requestID)
	}
	file = filepath.Base(file) // only the filename at the moment
	return c.logger.WithFields(logrus.Fields{
		"file":  file,
		"line":  line,
		"reqid": c.requestID,
	})
}

func (c *Context) extractOptionStringFromPath(key string, target *string) error {
	s, ok := c.request.PathParams[key]
	if !ok {
		return nil
	}

	*target = s
	return nil
}

func (c *Context) renderJSON(status int, v interface{}) {
	c.HTTPStatus = status

	data, err := json.Marshal(v)
	if err != nil {
		// TODO logging
		c.response.Header().Set("Content-Type", "application/json")
		c.response.WriteHeader(http.StatusInternalServerError)
		c.response.Write([]byte(`
{
  "errors": [
    {
      "code": "E0006",
      "message": "internal server error"
    }
  ]
}
`))
		return
	}

	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status)
	_, err = c.response.Write(data)
	if err != nil {
		c.ErrLog(err).Error("Cannot write a response")
	}
}

// RenderJSON renders a successful result as a JSON.
func (c *Context) RenderJSON(v interface{}) {
	c.renderJSON(http.StatusOK, v)
}

// RenderErrorJSON renders a failing result as a JSON.
func (c *Context) RenderErrorJSON(e *Error) {
	e.SetRequestID(c.requestID)
	c.renderJSON(e.Status, map[string]interface{}{
		"error": e,
	})
}

// Body returns a slice containing whole request body.
// It caches the result so controllers can call this
// method as many time as it wants.
//
// When the request body is empty (i.e. Read(req.Body)
// returns io.EOF), this method returns and caches
// an empty body slice and a nil error.
func (c *Context) Body() ([]byte, error) {
	if c.body != nil || c.bodyError != nil {
		return c.body, c.bodyError
	}

	body, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		if err == io.EOF {
			// when body is empty, this method caches
			// an empty slice and a nil error.
			body = []byte{}
			err = nil
		}
	}

	// Close and replace with new ReadCloser for parsing
	// mime/multipart request body by Request.FormFile method.
	c.request.Body.Close()
	c.request.Body = ioutil.NopCloser(bytes.NewReader(body))
	c.body = body
	c.bodyError = err
	return body, err
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

// SetUpContextAndRouter creates a root router of the API server and its context.
func SetUpContextAndRouter(prefix string, gvariables *ContextGlobalVariables) (*web.Router, error) {
	gvars := *gvariables
	udsStorage, err := setUpUDSStorage(&gvars.Config.Storage.UDS)
	if err != nil {
		return nil, err
	}

	// Topologies should be created after setting up everything necessary for it.
	if err := setUpTopologies(gvars.Logger, gvars.Topologies, gvars.Config, udsStorage); err != nil {
		return nil, err
	}

	root := web.NewWithPrefix(Context{}, prefix)
	root.NotFound((*Context).NotFoundHandler)
	root.Middleware(func(c *Context, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		c.SetLogger(gvars.Logger)
		c.udsStorage = udsStorage
		c.topologies = gvars.Topologies
		c.config = gvars.Config
		next(rw, req)
	})
	root.Middleware((*Context).setUpContext)
	return root, nil
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

	tp := core.NewDefaultTopology(core.NewContext(cc), name)
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
			// TODO: add stmt.String()
		}).Error("Cannot add a statement to the topology")
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
				// TODO: add stmt.String()
			}).Error("Cannot add a statement to the topology")
			return nil, err
		}
	}

	shouldStop = false
	return tb, nil
}

/*
Package runfile implements sensorbee runfile command. This command is provided
separately from run command to reduce the footprint of runfile command itself.
For example, if sensorbee doesn't need to provide sensorbee run command and only
offers sensorbee runfile, the footprint of sensorbee command might be reduced a
lot due to low functionality of this command.
*/
package runfile

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"pfi/sensorbee/sensorbee/server/config"
	"pfi/sensorbee/sensorbee/server/udsstorage"
)

// SetUp sets up a command for running single BQL file.
func SetUp() cli.Command {
	cmd := cli.Command{
		Name:        "runfile",
		Usage:       "run a BQL file",
		Description: "runfile command runs a BQL file",
		Action:      Run,
	}

	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "",
			Usage:  "file path of a config file in YAML format (only logging and storage sections are used)",
			EnvVar: "SENSORBEE_CONFIG",
		},
	}
	return cmd
}

func Run(c *cli.Context) {
	// TODO: Merge this implementation with cmd/run

	if len(c.Args()) != 1 {
		cli.ShowSubcommandHelp(c)
		os.Exit(1)
	}

	conf, err := func() (*config.Config, error) {
		if c.IsSet("config") {
			p := c.String("config")
			in, err := ioutil.ReadFile(p)
			if err != nil {
				return nil, fmt.Errorf("Cannot read the config file %v: %v\n", p, err)
			}

			var yml map[string]interface{}
			if err := yaml.Unmarshal(in, &yml); err != nil {
				return nil, fmt.Errorf("Cannot parse the config file %v: %v\n", p, err)
			}
			m, err := data.NewMap(yml)
			if err != nil {
				return nil, fmt.Errorf("The config file %v has invalid values: %v\n", p, err)
			}
			return config.New(m)

		} else {
			// Currently there's no required parameters. However, when a required
			// parameter is added to Config, remove this else block and add a
			// default value like "/etc/sensorbee/config.yaml" to the config option.
			return config.New(data.Map{})
		}
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w, err := conf.Logging.CreateWriter()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger := logrus.New()
	logger.Out = w

	udsStorage, err := setUpUDSStorage(&conf.Storage.UDS)
	if err != nil {
		logger.WithField("err", err).Error("Cannot set up a storage for UDSs")
		os.Exit(1)
	}

	logger.Info("Setting up a topology")
	bqlFile := c.Args()[0]
	tb, err := setUpTopology(logger, conf, bqlFile)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"bql_file": bqlFile,
		}).Error("Cannot set up the topology")
		os.Exit(1)
	}
	defer func() {
		logger.Info("Waiting for all nodes to finish processing tuples")
		if err := tb.Topology().Stop(); err != nil {
			logger.WithField("err", err).Error("Cannot stop the topology")
			os.Exit(1)
		}
	}()
	tb.UDSStorage = udsStorage

	logger.Info("Starting the topology")
	for name, s := range tb.Topology().Sources() {
		if err := s.Resume(); err != nil {
			logger.WithFields(logrus.Fields{
				"err":    err,
				"source": name,
			}).Error("Cannot resume the source")
			os.Exit(1)
		}
	}

	for _, s := range tb.Topology().Sources() {
		// TODO: error check if necessary
		s.State().Wait(core.TSStopped)
	}
	logger.Info("All sources has been stopped.")

}

func setUpUDSStorage(conf *config.UDSStorage) (udf.UDSStorage, error) {
	// TODO: merge this implementation with server/context.go
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

func setUpTopology(logger *logrus.Logger, conf *config.Config, bqlFile string) (*bql.TopologyBuilder, error) {
	queries, err := func() (string, error) {
		f, err := os.Open(bqlFile)
		if err != nil {
			return "", err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}()
	if err != nil {
		return nil, err
	}

	bp := parser.New()
	// TODO: provide better parse error reporting using ParseStmt instead of ParseStmts
	stmts, err := bp.ParseStmts(string(queries))
	if err != nil {
		return nil, err
	}

	cc := &core.ContextConfig{
		Logger: logger,
	}
	cc.Flags.DroppedTupleLog.Set(conf.Logging.LogDroppedTuples)
	cc.Flags.DroppedTupleSummarization.Set(conf.Logging.SummarizeDroppedTuples)

	name := "runfile" // FIXME: bql file name is better?
	tp := core.NewDefaultTopology(core.NewContext(cc), name)
	tb, err := bql.NewTopologyBuilder(tp)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new topology builder: %v", err)
	}

	for _, stmt := range stmts {
		// TODO: if stmt is CREATE SOURCE, create it with PAUSED
		if n, err := tb.AddStmt(stmt); err != nil {
			logger.WithFields(logrus.Fields{
				"err":      err,
				"topology": name,
				"stmt":     stmt,
			}).Error("Cannot add a statement to the topology")
			return nil, err // FIXME: logger output "err" two twice
		} else if n != nil && n.Type() == core.NTSource {
			sn, _ := n.(core.SourceNode)
			if _, ok := sn.Source().(core.RewindableSource); ok {
				return nil, fmt.Errorf(`rewindable source "%v" isn't supported`, n.Name())
			}
		}
	}
	return tb, nil
}

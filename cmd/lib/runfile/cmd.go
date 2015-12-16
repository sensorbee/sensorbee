package runfile

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/server/config"
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

	udsStorage, err := setUpUDSStorage(conf.Storage.UDS)
	if err != nil {
		logger.WithField("err", err).Error("Cannot set up a storage for UDSs")
		os.Exit(1)
	}

	logrus.Info("Setting up a topology")
	bqlFile := c.Args()[0]
	tb, err := setUpTopology(bqlFile)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"bql_file": bqlFile,
		}).Error("Cannot set up the topology")
	}
	defer func() {
		logrus.Info("Stopping the topology")
		if err := tb.Topology().Stop(); err != nil {
			logrus.WithField("err", err).Error("Cannot stop the topology")
		}
	}()

	logrus.Info("Starting the topology")
	for name, s := range tb.Topology().Sources() {
		if err := s.Resume(); err != nil {
			logrus.WithFields(logrus.Fields{
				"err", err,
				"source": name,
			}).Error("Cannot resume the source")
			return
		}
	}

	for _, s := range tb.Topology().Sources() {
		// TODO: error check if necessary
		s.State().Wait(core.TSStopped)
	}
	logrus.Info("All sources has been stopped.")

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
			return err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}()
	if err != nil {
		return err
	}

	bp := parser.New()
	// TODO: provide better parse error reporting using ParseStmt instead of ParseStmts
	stmts, err := bp.ParseStmts(string(queries))
	if err != nil {
		return nil, err
	}

	// TODO: create context using conf and setting the logger
	// TODO: create core.Topology
	// TODO: create bql.TopologyBuilder

	for _, stmt := range stmts {
		// TODO: if stmt is CREATE SOURCE, create it with PAUSED
		if _, err := tb.AddStmt(stmt); err != nil {
			logger.WithFields(logrus.Fields{
				"err":      err,
				"topology": name,
				"stmt":     stmt,
			}).Error("Cannot add a statement to the topology")
			return nil, err
		}
	}
	// TODO: return bql.TopologyBuilder, nil
}

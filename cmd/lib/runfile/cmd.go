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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"gopkg.in/sensorbee/sensorbee.v0/server/config"
	"gopkg.in/sensorbee/sensorbee.v0/server/udsstorage"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
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
		cli.StringFlag{
			Name:  "save-uds, s",
			Value: "",
			Usage: "save UDSs after all tuples are processed",
		},
		cli.StringFlag{
			Name:  "topology, t",
			Value: "",
			Usage: "name of the topology",
		},
	}
	return cmd
}

// Run runs "runfile" command.
func Run(c *cli.Context) error {
	// TODO: Merge this implementation with cmd/run
	if len(c.Args()) != 1 {
		cli.ShowSubcommandHelp(c)
		os.Exit(1)
	}

	emptyError := fmt.Errorf("") // to provide exit code but not error message for cli
	err := func() (retErr error) {
		conf, err := func() (*config.Config, error) {
			if c.IsSet("config") {
				p := c.String("config")
				in, err := ioutil.ReadFile(p)
				if err != nil {
					return nil, fmt.Errorf("Cannot read the config file %v: %v", p, err)
				}

				var yml map[string]interface{}
				if err := yaml.Unmarshal(in, &yml); err != nil {
					return nil, fmt.Errorf("Cannot parse the config file %v: %v", p, err)
				}
				m, err := data.NewMap(yml)
				if err != nil {
					return nil, fmt.Errorf("The config file %v has invalid values: %v", p, err)
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
			return err
		}

		w, err := conf.Logging.CreateWriter()
		if err != nil {
			return err
		}
		logLevel, err := logrus.ParseLevel(conf.Logging.MinLogLevel)
		if err != nil {
			return err
		}
		logger := logrus.New()
		logger.Out = w
		logger.Level = logLevel

		udsStorage, err := setUpUDSStorage(&conf.Storage.UDS)
		if err != nil {
			logger.WithField("err", err).Error("Cannot set up a storage for UDSs")
			return emptyError
		}

		logger.Info("Setting up a topology")
		bqlFile := c.Args()[0]
		topologyName := filepath.Base(bqlFile)
		topologyName = topologyName[:len(topologyName)-len(filepath.Ext(topologyName))]
		if n := c.String("topology"); n != "" {
			topologyName = n
		}

		tb, err := setUpTopology(topologyName, logger, conf, udsStorage)
		if err != nil {
			logger.WithField("err", err).Error("Cannot set up the topology")
			return emptyError
		}

		if err := setUpBQLStmt(tb, bqlFile); err != nil {
			logger.WithFields(logrus.Fields{
				"err":      err,
				"bql_file": bqlFile,
			}).Error("Cannot set up BQL statement")
			return emptyError
		}
		if c.IsSet("save-uds") {
			if err := hasStates(tb, c.String("save-uds")); err != nil {
				logger.WithField("err", err).Error("Cannot set up 'save-uds' option")
				return emptyError
			}
		}

		defer func() {
			logger.Info("Waiting for all nodes to finish processing tuples")
			if err := tb.Topology().Stop(); err != nil {
				logger.WithField("err", err).Error("Cannot stop the topology")
				retErr = emptyError
				return
			}
			logger.Info("Topology stopped")

			if c.IsSet("save-uds") {
				saveUDSList := c.String("save-uds")
				if err := saveStates(tb, saveUDSList); err != nil {
					logger.WithFields(logrus.Fields{
						"err":      err,
						"topology": tb.Topology().Name(),
					}).Error("Cannot save UDSs")
					retErr = emptyError
				}
			}

			registry := tb.Topology().Context().SharedStates
			if states, err := registry.List(); err != nil {
				logger.WithField("err", err).Error("Cannot list shared states")
				retErr = emptyError
				return
			} else if 0 < len(states) {
				logger.Info("Terminating states")
				for name := range states {
					if _, err := registry.Remove(name); err != nil {
						logger.WithFields(logrus.Fields{
							"err":   err,
							"state": name,
						}).Error("Cannot terminate state")
						retErr = emptyError
						// Continue to terminate the next state.
					}
				}
			}
		}()

		logger.WithField("config", conf.ToMap()).Info("Starting the topology")
		for name, s := range tb.Topology().Sources() {
			if err := s.Resume(); err != nil {
				logger.WithFields(logrus.Fields{
					"err":    err,
					"source": name,
				}).Error("Cannot resume the source")
				return emptyError
			}
		}

		for _, s := range tb.Topology().Sources() {
			// TODO: error check if necessary
			s.State().Wait(core.TSStopped)
		}
		logger.Info("All sources has been stopped.")
		return nil
	}()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
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

func setUpTopology(name string, logger *logrus.Logger, conf *config.Config, us udf.UDSStorage) (
	*bql.TopologyBuilder, error) {
	cc := &core.ContextConfig{
		Logger: logger,
	}
	cc.Flags.DroppedTupleLog.Set(conf.Logging.LogDroppedTuples)
	cc.Flags.DestinationlessTupleLog.Set(conf.Logging.LogDestinationlessTuples)
	cc.Flags.DroppedTupleSummarization.Set(conf.Logging.SummarizeDroppedTuples)

	tp, err := core.NewDefaultTopology(core.NewContext(cc), name)
	if err != nil {
		return nil, err
	}
	tb, err := bql.NewTopologyBuilder(tp)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new topology builder: %v", err)
	}
	tb.UDSStorage = us

	return tb, nil
}

func setUpBQLStmt(tb *bql.TopologyBuilder, bqlFile string) error {
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
		return err
	}

	bp := parser.New()
	// TODO: provide better parse error reporting using ParseStmt instead of ParseStmts
	stmts, err := bp.ParseStmts(string(queries))
	if err != nil {
		return err
	}

	for _, stmt := range stmts {
		// TODO: if stmt is CREATE SOURCE, create it with PAUSED
		if n, err := tb.AddStmt(stmt); err != nil {
			tb.Topology().Context().ErrLog(err).WithField("stmt", stmt).Error(
				"Cannot add a statement to the topology")
			return err // FIXME: logger output "err" two twice
		} else if n != nil && n.Type() == core.NTSource {
			sn, _ := n.(core.SourceNode)
			if _, ok := sn.Source().(core.RewindableSource); ok {
				return fmt.Errorf(`rewindable source "%v" isn't supported`, n.Name())
			}
		}
	}
	return nil
}

func hasStates(tb *bql.TopologyBuilder, saveUDSList string) error {
	if saveUDSList == "" {
		return nil
	}
	states, err := tb.Topology().Context().SharedStates.List()
	if err != nil {
		return err
	}
	missing := []string{}
	for _, name := range strings.Split(saveUDSList, ",") {
		if _, ok := states[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("the topology doesn't have UDSs: %v",
			strings.Join(missing, ","))
	}
	return nil
}

func saveStates(tb *bql.TopologyBuilder, saveUDSList string) error {
	states, err := tb.Topology().Context().SharedStates.List()
	if err != nil {
		return err
	}
	// if save UDS list is empty then all UDS will be saved.
	saveUDSs := []string{}
	if saveUDSList == "" {
		for k := range states {
			saveUDSs = append(saveUDSs, k)
		}
	} else {
		saveUDSs = strings.Split(saveUDSList, ",")
		for i, n := range saveUDSs {
			saveUDSs[i] = strings.TrimSpace(n)
		}
	}

	saveErrorFlag := false
	for _, name := range saveUDSs {
		state, ok := states[name]
		if !ok {
			err := fmt.Errorf("the UDS is not found")
			tb.Topology().Context().ErrLog(err).WithField("uds", name).Error(
				"Cannot save the UDS")
			continue
		}
		target, ok := state.(core.SavableSharedState)
		if !ok {
			tb.Topology().Context().Log().WithField("uds", name).Info(
				"The UDS doesn't support Save")
			continue
		}
		// TODO get tag
		w, err := tb.UDSStorage.Save(tb.Topology().Name(), name, "default")
		if err != nil {
			tb.Topology().Context().ErrLog(err).WithField("uds", name).Error(
				"Cannot save the UDS")
			saveErrorFlag = true
			continue
		}
		err = func() error {
			defer func() {
				if err := w.Commit(); err != nil {
					tb.Topology().Context().ErrLog(err).WithField("uds", name).Error(
						"Cannot save the UDS")
					saveErrorFlag = true
				}
			}()
			// TODO get parameters
			return target.Save(tb.Topology().Context(), w, data.Map{})
		}()
		if err != nil {
			tb.Topology().Context().ErrLog(err).WithField("uds", name).Error(
				"Cannot save the UDS")
			saveErrorFlag = true
		}
	}
	if saveErrorFlag {
		return fmt.Errorf("fail to save or commit some UDSs")
	}
	return nil
}

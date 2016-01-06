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
	"strings"
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
	}
	return cmd
}

// Run runs "runfile" command.
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
	tb, err := setUpTopology(logger, conf, udsStorage)
	if err != nil {
		logger.WithField("err", err).Error("Cannot set up the topology")
		os.Exit(1)
	}

	if err := setUpBQLStmt(tb, bqlFile); err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"bql_file": bqlFile,
		}).Error("Cannot set up BQL statement")
		os.Exit(1)
	}

	defer func() {
		logger.Info("Waiting for all nodes to finish processing tuples")
		if err := tb.Topology().Stop(); err != nil {
			logger.WithField("err", err).Error("Cannot stop the topology")
			os.Exit(1)
		}

		if c.IsSet("save-uds") {
			saveUDSList := c.String("save-uds")
			if err := saveStates(tb, saveUDSList); err != nil {
				logger.WithFields(logrus.Fields{
					"err":      err,
					"topology": tb.Topology().Name(),
				}).Error("Cannot save UDSs")
				os.Exit(1)
			}
		}
		// TODO: Terminate all shared states
	}()

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

func setUpTopology(logger *logrus.Logger, conf *config.Config, us udf.UDSStorage) (
	*bql.TopologyBuilder, error) {
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

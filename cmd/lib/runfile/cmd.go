package runfile

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
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
	}

	// TODO: set up logger
	// TODO: set up storage
	// TODO: parse statements
	// TODO: bql.TopologyBuilder
	// TODO: AddStmt: Pause all sources, and resume them after build the topology
	// TODO: Wait until all sources stops
	// TODO: call core.Topology.Stop
}

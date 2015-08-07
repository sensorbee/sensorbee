package exp

import (
	"github.com/codegangsta/cli"
)

// SetUp sets up SensorBee's experimenting command.
func SetUp() cli.Command {
	cmd := cli.Command{
		Name:  "exp",
		Usage: "experiment BQL statements",
		Subcommands: []cli.Command{
			setUpRun(),
			setUpClean(),
			// TODO: file <node>: get cache file name of the node
			// TODO: hash <node>: get hash of the latest cache of the node
		},
	}
	return cmd
}

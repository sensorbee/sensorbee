package main

import (
	"github.com/codegangsta/cli"
	"os"
	_ "pfi/sensorbee/sensorbee/bql/udf/builtin"
	"pfi/sensorbee/sensorbee/cmd/lib/exp"
	"pfi/sensorbee/sensorbee/cmd/lib/run"
	"pfi/sensorbee/sensorbee/cmd/lib/shell"
	"pfi/sensorbee/sensorbee/cmd/lib/topology"
	"time"
)

func init() {
	// TODO
	time.Local = time.UTC
}

func main() {
	app := cli.NewApp()
	app.Name = "sensorbee"
	app.Usage = "SensorBee"
	app.Version = "0.3.2" // TODO get dynamic, will be get from external file
	app.Commands = []cli.Command{
		run.SetUp(),
		shell.SetUp(),
		topology.SetUp(),
		exp.SetUp(),
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

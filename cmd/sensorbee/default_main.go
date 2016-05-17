package main

import (
	"github.com/codegangsta/cli"
	_ "gopkg.in/sensorbee/sensorbee.v0/bql/udf/builtin"
	"gopkg.in/sensorbee/sensorbee.v0/cmd/lib/exp"
	"gopkg.in/sensorbee/sensorbee.v0/cmd/lib/run"
	"gopkg.in/sensorbee/sensorbee.v0/cmd/lib/runfile"
	"gopkg.in/sensorbee/sensorbee.v0/cmd/lib/shell"
	"gopkg.in/sensorbee/sensorbee.v0/cmd/lib/topology"
	"gopkg.in/sensorbee/sensorbee.v0/version"
	"os"
	"time"
)

func init() {
	// TODO
	time.Local = time.UTC
}

func main() {
	app := cli.NewApp()
	app.Name = "sensorbee"
	app.Usage = "SensorBee built with build_sensorbee 0.4.1"
	app.Version = version.Version
	app.Commands = []cli.Command{
		run.SetUp(),
		shell.SetUp(),
		topology.SetUp(),
		exp.SetUp(),
		runfile.SetUp(),
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

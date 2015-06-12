package main

import (
	"github.com/codegangsta/cli"
	"os"
	"time"
)

type CommandGenerator func() cli.Command

func init() {
	// TODO
	time.Local = time.UTC
}

func main() {
	app := setUpApp([]CommandGenerator{
		SetUpRunCommand,
	})

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func setUpApp(cmds []CommandGenerator) *cli.App {
	app := cli.NewApp()
	app.Name = "sensorbeeserver"
	app.Usage = "SenserBee server process"
	app.Version = "0.0.1" // TODO get dynmic
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "/etc/sersorbee/sensorbee.config",
			Usage:  "path to the config file",
			EnvVar: "HAWK_CONFIG",
		},
	}
	app.Before = appBeforeHook

	for _, c := range cmds {
		app.Commands = append(app.Commands, c())
	}
	return app
}

func appBeforeHook(c *cli.Context) error {
	if err := loadConfig(c); err != nil {
		return err
	}
	return nil
}

func loadConfig(c *cli.Context) error {
	// TODO load configuration file (YAML)
	return nil
}

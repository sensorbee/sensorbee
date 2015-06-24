package main

import (
	"github.com/codegangsta/cli"
	"os"
	// _ "pfi/scouter-snippets/plugin"
	"pfi/sensorbee/sensorbee/client"
	"pfi/sensorbee/sensorbee/server"
	"time"
)

type commandGenerator func() cli.Command

func init() {
	// TODO
	time.Local = time.UTC
}

func main() {
	app := setUpApp([]commandGenerator{
		server.SetUpRunCommand,
		client.SetUpCMDLineToolCommand,
	})

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func setUpApp(cmds []commandGenerator) *cli.App {
	app := cli.NewApp()
	app.Name = "sensorbee"
	app.Usage = "SenserBee"
	app.Version = "0.0.1" // TODO get dynamic, will be get from external file
	app.Flags = []cli.Flag{
		cli.StringFlag{ // TODO get configuration from external file
			Name:   "config, c",
			Value:  "/etc/sersorbee/sensorbee.config",
			Usage:  "path to the config file",
			EnvVar: "SENSORBEE_CONFIG",
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

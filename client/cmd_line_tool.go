package client

import (
	"github.com/codegangsta/cli"
	"pfi/sensorbee/sensorbee/client/cmd"
)

func SetUpCMDLineToolCommand() cli.Command {
	cmd := cli.Command{
		Name:        "cmd",
		Usage:       "command line tool",
		Description: "cmd command launch command line tool",
		Action:      Launch,
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "uri",
			Value:  "http://localhost:8090/api/v1",
			Usage:  "target URI to launch",
			EnvVar: "URI",
		},
		cli.StringFlag{
			Name:   "version,v",
			Value:  "v1",
			Usage:  "SenserBee API version",
			EnvVar: "VERSION",
		},
	}
	return cmd
}

// Launch SensorBee's command line client tool.
func Launch(c *cli.Context) {
	uri := c.String("uri")

	app := cmd.SetUpCommands(cmd.NewTopologiesCommands())
	app.Run(uri)
}

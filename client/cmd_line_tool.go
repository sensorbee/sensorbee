package client

import (
	"github.com/codegangsta/cli"
	"pfi/sensorbee/sensorbee/client/cmd"
	"strings"
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
			Value:  "http://localhost:8090/api",
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
	host := c.String("uri")
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	uri := host + c.String("version")

	cmds := []cmd.Command{}
	for _, c := range cmd.NewTopologiesCommands() {
		cmds = append(cmds, c)
	}
	for _, c := range cmd.NewFileLoadCommands() {
		cmds = append(cmds, c)
	}
	app := cmd.SetUpCommands(cmds)
	app.Run(uri)
}

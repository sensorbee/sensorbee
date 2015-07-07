package shell

import (
	"github.com/codegangsta/cli"
	"strings"
)

func SetUp() cli.Command {
	cmd := cli.Command{
		Name:   "shell",
		Usage:  "BQL shell",
		Action: Launch,
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

	cmds := []Command{}
	for _, c := range NewTopologiesCommands() {
		cmds = append(cmds, c)
	}
	for _, c := range NewFileLoadCommands() {
		cmds = append(cmds, c)
	}
	app := SetUpCommands(cmds)
	app.Run(uri)
}

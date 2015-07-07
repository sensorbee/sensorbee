package shell

import (
	"github.com/codegangsta/cli"
	"strings"
)

func SetUp() cli.Command {
	cmd := cli.Command{
		Name:        "shell",
		Usage:       "BQL shell",
		Description: "shell command launches an interactive shell for BQL",
		Action:      Launch,
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "uri",
			Value: "http://localhost:8090/",
			Usage: "target URI to launch",
		},
		cli.StringFlag{
			Name:  "version,v",
			Value: "v1",
			Usage: "SenserBee API version",
		},
	}
	return cmd
}

// Launch SensorBee's command line client tool.
func Launch(c *cli.Context) {
	host := c.String("uri") // TODO: validate URI
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	uri := host + "api/" + c.String("version")

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

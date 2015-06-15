package client

import (
	"github.com/codegangsta/cli"
	"pfi/sensorbee/sensorbee/client/cmd"
)

type REQType int

const (
	REQ_URL = "http://localhost:8090/api" // TODO from flag.Parse()
)

func SetUpCMDLineToolCommand() cli.Command {
	cmd := cli.Command{
		Name:        "cmd",
		Usage:       "command line tool",
		Description: "cmd command launch command line tool",
		Action:      Launch,
	}
	return cmd
}

// Launch SensorBee's command line client tool.
func Launch(c *cli.Context) {
	app := cmd.SetUpCommands(cmd.NewBQLCommands())
	app.Run(REQ_URL)
}

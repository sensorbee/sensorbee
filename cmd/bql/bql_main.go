package main

import (
	"gopkg.in/sensorbee/sensorbee.v0/cmd/lib/shell"
	"gopkg.in/sensorbee/sensorbee.v0/version"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "bql"
	app.Usage = "bql command launches an interactive shell for BQL"
	app.Version = version.Version
	app.Flags = shell.CmdFlags
	app.Action = shell.Launch

	app.Run(os.Args)
}

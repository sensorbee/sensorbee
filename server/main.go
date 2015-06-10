package main

import (
	"pfi/sensorbee/sensorbee/server/cmd"
)

func main() {
	app := setUpApp()
	app.Run()
}

func setUpApp() cmd.App {
	return cmd.SetUpCommands(cmd.NewBQLCommands())
}

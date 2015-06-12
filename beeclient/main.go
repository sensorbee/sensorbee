package main

import (
	"pfi/sensorbee/sensorbee/beeclient/cmd"
)

type REQType int

const (
	REQ_URL = "http://localhost:8090/api" // TODO from flag.Parse()
)

func main() {
	app := cmd.SetUpCommands(cmd.NewBQLCommands())
	app.Run(REQ_URL)
}

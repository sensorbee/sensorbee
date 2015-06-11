package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func SetUpBQLCommand() cli.Command {
	return cli.Command{
		Name:        "bql",
		Usage:       "execute BQL statements",
		Description: "execute BQL statements to run topologies",
		Action:      Execute,
	}
}

func Execute(c *cli.Context) {
	fmt.Println("hello!")
}

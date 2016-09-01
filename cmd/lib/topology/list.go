package topology

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"gopkg.in/sensorbee/sensorbee.v0/server/response"
	"gopkg.in/urfave/cli.v1"
)

func setUpList() cli.Command {
	return cli.Command{
		Name:        "list",
		Aliases:     []string{"l"},
		Usage:       "get a list of topologies",
		Description: "list commands show a list of names of topologies in the server",
		Action:      actionWrapper(runList),
		Flags:       commonFlags,
		// TODO: add flags like "ls -l"
		// TODO: maybe pagination?
	}
}

func runList(c *cli.Context) error {
	if err := validateFlags(c); err != nil {
		return err
	}

	if len(c.Args()) > 0 {
		return fmt.Errorf("too many command line arguments")
	}

	res, err := do(c, client.Get, "topologies", nil, "Cannot get a list of topologies")
	if err != nil {
		return err
	}
	ts := struct {
		Topologies []*response.Topology `json:"topologies"`
	}{}
	if err := res.ReadJSON(&ts); err != nil { // ReadJSON closes the body
		return fmt.Errorf("Cannot read a response: %v", err)
	}

	for _, t := range ts.Topologies {
		fmt.Fprintln(c.App.Writer, t.Name)
	}
	return nil
}

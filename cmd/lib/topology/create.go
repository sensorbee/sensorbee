package topology

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/sensorbee/sensorbee.v0/client"
)

func setUpCreate() cli.Command {
	return cli.Command{
		Name:        "create",
		Aliases:     []string{"c"},
		Usage:       "create a new topology",
		Description: "sensorbee topology create <topology_name> create a new topology having <topology_name>.",
		Action:      actionWrapper(runCreate),
		Flags:       commonFlags,
	}
}

func runCreate(c *cli.Context) error {
	if err := validateFlags(c); err != nil {
		return err
	}

	args := c.Args()
	switch l := len(args); l {
	case 1:
		// ok
	case 0:
		return fmt.Errorf("topology_name is missing")
	default:
		return fmt.Errorf("too many command line arguments")
	}

	res, err := do(c, client.Post, "topologies", map[string]interface{}{
		"name": args[0],
	}, "Cannot create a topology")
	if err != nil {
		return err
	}
	// TODO: show something about the created topology
	return res.Close()
}

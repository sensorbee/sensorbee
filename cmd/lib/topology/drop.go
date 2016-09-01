package topology

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/urfave/cli.v1"
	"path"
)

func setUpDrop() cli.Command {
	return cli.Command{
		Name:        "drop",
		Aliases:     []string{"d"},
		Usage:       "drop an existing topology",
		Description: "sensorbee topology drop <topology_name> drops an existing topology having <topology_name>",
		Action:      actionWrapper(runDrop),
		Flags:       commonFlags,
	}
}

func runDrop(c *cli.Context) error {
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

	name := args[0]
	if err := core.ValidateSymbol(name); err != nil {
		// This is checked here to avoid sending wrong DELETE request to different URL.
		return fmt.Errorf("The name of the topology is invalid: %v", err)
	}
	res, err := do(c, client.Delete, path.Join("topologies", name), nil, "Cannot drop a topology")
	if err != nil {
		return err
	}
	return res.Close()
}

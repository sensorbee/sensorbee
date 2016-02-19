package topology

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"os"
)

func setUpCreate() cli.Command {
	return cli.Command{
		Name:        "create",
		Aliases:     []string{"c"},
		Usage:       "create a new topology",
		Description: "sensorbee topology create <topology_name> create a new topology having <topology_name>.",
		Action:      runCreate,
		Flags:       commonFlags,
	}
}

func runCreate(c *cli.Context) {
	defer panicHandler()
	validateFlags(c)

	args := c.Args()
	switch l := len(args); l {
	case 1:
		// ok
	case 0:
		fmt.Fprintln(os.Stderr, "topology_name is missing")
		panic(1) // TODO: define exit code properly
	default:
		fmt.Fprintln(os.Stderr, "too many command line arguments")
		panic(1)
	}

	do(c, client.Post, "topologies", map[string]interface{}{
		"name": args[0],
	}, "Cannot create a topology").Close()
	// TODO: show something about the created topology
}

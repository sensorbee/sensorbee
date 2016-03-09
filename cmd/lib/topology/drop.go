package topology

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"os"
	"path"
)

func setUpDrop() cli.Command {
	return cli.Command{
		Name:        "drop",
		Aliases:     []string{"d"},
		Usage:       "drop an existing topology",
		Description: "sensorbee topology drop <topology_name> drops an existing topology having <topology_name>",
		Action:      runDrop,
		Flags:       commonFlags,
	}
}

func runDrop(c *cli.Context) {
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

	name := args[0]
	if err := core.ValidateSymbol(name); err != nil {
		// This is checked here to avoid sending wrong DELETE request to different URL.
		fmt.Fprintf(os.Stderr, "The name of the topology is invalid: %v\n", err)
		panic(1)
	}
	do(c, client.Delete, path.Join("topologies", name), nil, "Cannot drop a topology").Close()
}

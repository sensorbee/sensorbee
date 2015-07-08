package topology

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"pfi/sensorbee/sensorbee/client"
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

	req := newRequester(c)
	res, err := req.Do(client.PostRequest, "topologies", map[string]interface{}{
		"name": args[0],
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create a topology: %v\n", err)
		panic(1)
	}
	if res.IsError() {
		errRes, err := res.Error()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create a topology and failed to parse error information: %v\n", err)
			panic(1)
		}
		fmt.Fprintf(os.Stderr, "Cannot create a topology: %v, %v: %v\n", errRes.Code, errRes.RequestID, errRes.Message)
		panic(1)
	}

	// TODO: show something about the created topology
}

package topology

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path"
	"pfi/sensorbee/sensorbee/client"
	"pfi/sensorbee/sensorbee/core"
)

func setUpDelete() cli.Command {
	return cli.Command{
		Name:        "delete",
		Aliases:     []string{"d"},
		Usage:       "delete an existing topology",
		Description: "sensorbee topology delete <topology_name> delete an existing topology having <topology_name>",
		Action:      runDelete,
		Flags:       commonFlags,
	}
}

func runDelete(c *cli.Context) {
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
	if err := core.ValidateNodeName(name); err != nil {
		// This is checked here to avoid sending wrong DELETE request to different URL.
		fmt.Fprintf(os.Stderr, "The name of the topology is invalid: %v\n", err)
		panic(1)
	}

	req := newRequester(c)
	res, err := req.Do(client.DeleteRequest, path.Join("topologies", name), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot delete a topology: %v\n", err)
		panic(1)
	}
	if res.IsError() {
		errRes, err := res.Error()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot delete a topology and failed to parse error information: %v\n", err)
			panic(1)
		}
		fmt.Fprintf(os.Stderr, "Cannot delete a topology: %v, %v: %v\n", errRes.Code, errRes.RequestID, errRes.Message)
		panic(1)
	}
}

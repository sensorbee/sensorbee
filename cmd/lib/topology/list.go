package topology

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"pfi/sensorbee/sensorbee/client"
	"pfi/sensorbee/sensorbee/server/response"
)

func setUpList() cli.Command {
	return cli.Command{
		Name:        "list",
		Aliases:     []string{"l"},
		Usage:       "get a list of topologies",
		Description: "list commands show a list of names of topologies in the server",
		Action:      runList,
		Flags:       commonFlags,
		// TODO: add flags like "ls -l"
		// TODO: maybe pagination?
	}
}

func runList(c *cli.Context) {
	defer panicHandler()
	validateFlags(c)

	if len(c.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "too many command line arguments")
		panic(1)
	}

	res := do(c, client.Get, "topologies", nil, "Cannot get a list of topologies")
	ts := struct {
		Topologies []*response.Topology `json:"topologies"`
	}{}
	if err := res.ReadJSON(&ts); err != nil { // ReadJSON closes the body
		fmt.Fprintf(os.Stderr, "Cannot read a response: %v\n", err)
		panic(1)
	}

	for _, t := range ts.Topologies {
		fmt.Fprintln(c.App.Writer, t.Name)
	}
}

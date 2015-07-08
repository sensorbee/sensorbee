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

	req := newRequester(c)
	res, err := req.Do(client.GetRequest, "topologies", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get a list of topologies: %v\n", err)
		panic(1)
	}
	if res.IsError() {
		errRes, err := res.Error()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get a list of topologies and failed to parse error information: %v\n", err)
			panic(1)
		}
		fmt.Fprintf(os.Stderr, "Cannot get a list of topologies: %v, %v: %v\n", errRes.Code, errRes.RequestID, errRes.Message)
		panic(1)
	}

	ts := struct {
		Topologies []*response.Topology `json:"topologies"`
	}{}
	if err := res.ReadJSON(&ts); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read a response: %v\n", err)
		panic(1)
	}

	for _, t := range ts.Topologies {
		fmt.Fprintln(c.App.Writer, t.Name)
	}
}

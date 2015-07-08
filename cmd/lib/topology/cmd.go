package topology

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"pfi/sensorbee/sensorbee/client"
)

var (
	testMode     bool
	testExitCode int
)

func panicHandler() {
	// Because all tests uses this feature, they cannot be run in parallel.
	// Commands use this dirty logic because there's no other ways to report
	// errors due to cli library's limitation.
	if e := recover(); e != nil {
		ec, ok := e.(int)
		if !ok {
			panic(e)
		}

		if testMode {
			testExitCode = ec
		} else {
			os.Exit(e.(int))
		}

	} else if testMode {
		testExitCode = 0
	}
}

// SetUp sets up a subcommand for topology manipulation.
func SetUp() cli.Command {
	cmd := cli.Command{
		Name:        "topology",
		Aliases:     []string{"t"},
		Usage:       "manipulate topologies on SensorBee",
		Description: "topology command allows users to create, delete, list topologies",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			setUpCreate(),
			setUpList(),
			setUpDelete(),
		},
	}
	return cmd
}

var (
	commonFlags = []cli.Flag{
		cli.StringFlag{ // TODO: share this flag with others
			Name:   "uri",
			Value:  "http://localhost:8090/",
			Usage:  "the address of the target SensorBee server",
			EnvVar: "SENSORBEE_URI",
		},
		cli.StringFlag{ // TODO: share this flag with others
			Name:  "api-version",
			Value: "v1",
			Usage: "target API version",
		},
	}
)

func validateFlags(c *cli.Context) {
	if err := client.ValidateURL(c.String("uri")); err != nil {
		fmt.Fprintln(os.Stderr, "--uri flag has an invalid value: %v", err)
		panic(1)
	}
	if err := client.ValidateAPIVersion(c.String("api-version")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic(1)
	}
	// TODO: check other flags
}

func newRequester(c *cli.Context) *client.Requester {
	r, err := client.NewRequester(c.String("uri"), c.String("api-version"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create a API requester: %v\n", err)
		panic(1)
	}
	return r
}

func do(c *cli.Context, method client.Method, path string, body interface{}, baseErrMsg string) *client.Response {
	req := newRequester(c)
	res, err := req.Do(method, path, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", baseErrMsg, err)
		panic(1)
	}
	if res.IsError() {
		errRes, err := res.Error()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v and failed to parse error information: %v\n", baseErrMsg, err)
			panic(1)
		}
		fmt.Fprintf(os.Stderr, "%v: %v, %v: %v\n", baseErrMsg, errRes.Code, errRes.RequestID, errRes.Message)
		panic(1)
	}
	return res
}

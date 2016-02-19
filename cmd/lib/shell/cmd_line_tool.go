package shell

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"os"
)

// SetUp SensorBee shell tool. The tool sets up HTTP client and access to
// SensorBee server.
func SetUp() cli.Command {
	cmd := cli.Command{
		Name:        "shell",
		Usage:       "BQL shell",
		Description: "shell command launches an interactive shell for BQL",
		Action:      Launch,
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "uri",
			Value:  "http://localhost:8090/",
			Usage:  "the address of the target SensorBee server",
			EnvVar: "SENSORBEE_URI",
		},
		cli.StringFlag{
			Name:  "api-version",
			Value: "v1",
			Usage: "target API version",
		},
		cli.StringFlag{
			Name:  "topology,t",
			Usage: "the SensorBee topology to use (instead of USE command)",
		},
	}
	return cmd
}

// Launch SensorBee's command line client tool.
func Launch(c *cli.Context) {
	defer panicHandler()
	validateFlags(c)
	if c.IsSet("topology") {
		currentTopology.name = c.String("topology")
	}
	cmds := []Command{}
	for _, c := range NewTopologiesCommands() {
		cmds = append(cmds, c)
	}
	for _, c := range NewFileLoadCommands() {
		cmds = append(cmds, c)
	}
	app := SetUpCommands(cmds)
	app.Run(newRequester(c))
}

// TODO: merge following function implementations with lib/topology's
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

func validateFlags(c *cli.Context) {
	if err := client.ValidateURL(c.String("uri")); err != nil {
		fmt.Fprintf(os.Stderr, "--uri flag has an invalid value: %v\n", err)
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

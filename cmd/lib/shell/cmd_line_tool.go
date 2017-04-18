package shell

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"gopkg.in/sensorbee/sensorbee.v0/server/config"
	"gopkg.in/urfave/cli.v1"
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
	cmd.Flags = CmdFlags
	return cmd
}

// CmdFlags is list of shell command options
var CmdFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "uri",
		Value:  fmt.Sprintf("http://localhost:%d/", config.DefaultPort),
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
	cli.StringFlag{
		Name:  "output-format",
		Value: "bql",
		Usage: "the format of the output of BQL statements (bql or json)",
	},
}

// Global flags to control behavior of the shell.
// TODO: support dynamic modification of outputFormat.
var (
	outputFormat       = "bql"
	validOutputFormats = map[string]struct{}{
		"json": struct{}{},
		"bql":  struct{}{},
	}
)

// Launch SensorBee's command line client tool.
func Launch(c *cli.Context) error {
	err := func() error {
		if err := validateFlags(c); err != nil {
			return err
		}
		outputFormat = c.String("output-format")
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
		req, err := newRequester(c)
		if err != nil {
			return err
		}
		app.Run(req)
		return nil
	}()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func validateFlags(c *cli.Context) error {
	if err := client.ValidateURL(c.String("uri")); err != nil {
		return fmt.Errorf("--uri flag has an invalid value: %v", err)
	}
	if err := client.ValidateAPIVersion(c.String("api-version")); err != nil {
		return err
	}

	of := c.String("output-format")
	if _, ok := validOutputFormats[of]; !ok {
		return fmt.Errorf("invalid output format: %v", of)
	}
	return nil
}

func newRequester(c *cli.Context) (*client.Requester, error) {
	r, err := client.NewRequester(c.String("uri"), c.String("api-version"))
	if err != nil {
		return nil, fmt.Errorf("Cannot create a API requester: %v", err)
	}
	return r, nil
}

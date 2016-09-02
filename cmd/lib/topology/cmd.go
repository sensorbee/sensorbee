package topology

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/client"
	"gopkg.in/sensorbee/sensorbee.v0/server/config"
	"gopkg.in/urfave/cli.v1"
)

var (
	testMode     bool
	testExitCode int
)

// SetUp sets up a subcommand for topology manipulation.
func SetUp() cli.Command {
	cmd := cli.Command{
		Name:        "topology",
		Aliases:     []string{"t"},
		Usage:       "manipulate topologies on SensorBee",
		Description: "topology command allows users to manipulate topologies",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: []cli.Command{
			setUpCreate(),
			setUpList(),
			setUpDrop(),
		},
	}
	return cmd
}

func actionWrapper(f cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		testExitCode = 0
		if err := f(c); err != nil {
			// TODO: define exit code properly
			if testMode {
				testExitCode = 1
				return err
			}
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}
}

var (
	commonFlags = []cli.Flag{
		cli.StringFlag{ // TODO: share this flag with others
			Name:   "uri",
			Value:  fmt.Sprintf("http://localhost:%d/", config.DefaultPort),
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

func validateFlags(c *cli.Context) error {
	if err := client.ValidateURL(c.String("uri")); err != nil {
		return fmt.Errorf("--uri flag has an invalid value: %v", err)
	}
	if err := client.ValidateAPIVersion(c.String("api-version")); err != nil {
		return err
	}
	// TODO: check other flags
	return nil
}

func newRequester(c *cli.Context) (*client.Requester, error) {
	r, err := client.NewRequester(c.String("uri"), c.String("api-version"))
	if err != nil {
		return nil, fmt.Errorf("Cannot create a API requester: %v", err)
	}
	return r, nil
}

func do(c *cli.Context, method client.Method, path string, body interface{}, baseErrMsg string) (*client.Response, error) {
	req, err := newRequester(c)
	if err != nil {
		return nil, err
	}
	res, err := req.Do(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", baseErrMsg, err)
	}
	if res.IsError() {
		errRes, err := res.Error()
		if err != nil {
			return nil, fmt.Errorf("%v and failed to parse error information: %v", baseErrMsg, err)
		}
		return nil, fmt.Errorf("%v: %v, %v: %v", baseErrMsg, errRes.Code, errRes.RequestID, errRes.Message)
	}
	return res, nil
}

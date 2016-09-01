// run package implements sensorbee's subcommand command which runs an API
// server.
package run

import (
	"fmt"
	"gopkg.in/pfnet/jasco.v1"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"gopkg.in/sensorbee/sensorbee.v0/server"
	"gopkg.in/sensorbee/sensorbee.v0/server/config"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
)

// SetUp sets up SensorBee's HTTP server. The URL or port ID is set with server
// configuration file, or command line arguments.
func SetUp() cli.Command {
	cmd := cli.Command{
		Name:        "run",
		Usage:       "run the server",
		Description: "run command starts a new server process",
		Action:      Run,
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "",
			Usage:  "file path of a config file in YAML format",
			EnvVar: "SENSORBEE_CONFIG",
		},
	}
	return cmd
}

// Run run the HTTP server.
func Run(c *cli.Context) error {
	err := func() error {
		var conf *config.Config
		if c.IsSet("config") {
			p := c.String("config")
			in, err := ioutil.ReadFile(p)
			if err != nil {
				return fmt.Errorf("Cannot read the config file %v: %v", p, err)
			}

			var yml map[string]interface{}
			if err := yaml.Unmarshal(in, &yml); err != nil {
				return fmt.Errorf("Cannot parse the config file %v: %v", p, err)
			}
			m, err := data.NewMap(yml)
			if err != nil {
				return fmt.Errorf("The config file %v has invalid values: %v", p, err)
			}
			c, err := config.New(m)
			if err != nil {
				return fmt.Errorf("Cannot apply the cnofig file %v: %v", p, err)
			}
			conf = c

		} else {
			// Currently there's no required parameters. However, when a required
			// parameter is added to Config, remove this else block and add a
			// default value like "/etc/sensorbee/config.yaml" to the config option.
			c, err := config.New(data.Map{})
			if err != nil {
				return fmt.Errorf("Cannot apply the default config: %v", err)
			}
			conf = c
		}

		cgvars, err := server.SetUpContextGlobalVariables(conf)
		if err != nil {
			return fmt.Errorf("Cannot set up the server context: %v", err)
		}

		cgvars.Logger.WithField("config", conf.ToMap()).Info("Setting up the server context")

		jascoRoot := jasco.New("/", cgvars.Logger)
		router, err := server.SetUpContextAndRouter("/", jascoRoot, cgvars)
		if err != nil {
			return fmt.Errorf("Cannot set up the server context: %v", err)
		}
		server.SetUpAPIRouter("/", router, nil)

		bind := c.String("listen-on")
		if _, err := net.ResolveTCPAddr("tcp", bind); err != nil {
			return fmt.Errorf("--listen-on(-l) parameter has an invalid address: %v", err)
		}

		// TODO: support listening on multiple addresses
		// TODO: support stopping the server
		// TODO: support graceful shutdown
		s := &http.Server{
			Addr:    conf.Network.ListenOn,
			Handler: jascoRoot,
		}
		cgvars.Logger.Infof("Starting the server on %v", conf.Network.ListenOn)
		if err := s.ListenAndServe(); err != nil {
			return fmt.Errorf("Cannot start the server: %v", err)
		}
		cgvars.Logger.Infof("The server stopped")
		return nil
	}()
	if err != nil {
		// NOTE: using something like this allows tests to check exit codes.
		//   if testMode {
		//     testExitStatus = ec
		//   } else {
		//     return cli.NewExitError(err.Error(), 1)
		//   }
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// run package implements sensorbee's subcommand command which runs an API
// server.
package run

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"pfi/sensorbee/sensorbee/data"
	"pfi/sensorbee/sensorbee/server"
	"pfi/sensorbee/sensorbee/server/config"
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
func Run(c *cli.Context) {
	defer func() {
		// This logic is provided to write test codes for this command line tool like below:
		// if v := recover(); v != nil {
		//   if testMode {
		//     testExitStatus = v.(int)
		//   } else {
		//     os.Exit(v.(int))
		//   }
		// }
		if v := recover(); v != nil {
			os.Exit(v.(int))
		}
	}()

	var conf *config.Config
	if c.IsSet("config") {
		p := c.String("config")
		in, err := ioutil.ReadFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read the config file %v: %v\n", p, err)
			panic(1)
		}

		var yml map[string]interface{}
		if err := yaml.Unmarshal(in, &yml); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse the config file %v: %v\n", p, err)
			panic(1)
		}
		m, err := data.NewMap(yml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "The config file %v has invalid values: %v\n", p, err)
			panic(1)
		}
		c, err := config.New(m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot apply the cnofig file %v: %v\n", p, err)
			panic(1)
		}
		conf = c

	} else {
		// Currently there's no required parameters. However, when a required
		// parameter is added to Config, remove this else block and add a
		// default value like "/etc/sensorbee/config.yaml" to the config option.
		c, err := config.New(data.Map{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot apply the default config: %v\n", err)
			panic(1)
		}
		conf = c
	}

	cgvars, err := server.SetUpContextGlobalVariables(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot set up the server context: %v\n", err)
		panic(1)
	}

	cgvars.Logger.Info("Setting up the server context")

	router, err := server.SetUpContextAndRouter("/", cgvars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot set up the server context: %v\n", err)
		panic(1)
	}
	server.SetUpAPIRouter("/", router, nil)

	bind := c.String("listen-on")
	if _, err := net.ResolveTCPAddr("tcp", bind); err != nil {
		fmt.Fprintln(os.Stderr, "--listen-on(-l) parameter has an invalid address:", err)
		panic(1)
	}

	// TODO: support listening on multiple addresses
	// TODO: support stopping the server
	// TODO: support graceful shutdown
	s := &http.Server{
		Addr:    conf.Network.ListenOn,
		Handler: router,
	}
	cgvars.Logger.Infof("Starting the server on %v", conf.Network.ListenOn)
	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot start the server:", err)
		panic(1)
	}
	cgvars.Logger.Infof("The server stopped")
}

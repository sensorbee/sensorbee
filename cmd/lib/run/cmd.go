// run package implements sensorbee's subcommand command which runs an API
// server.
package run

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"net"
	"net/http"
	"os"
	"pfi/sensorbee/sensorbee/server"
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
			Name:  "listen-on, l",
			Value: "0.0.0.0:8090",
			Usage: "server port number",
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

	logger := logrus.New()
	// TODO: setup logger based on the config
	topologies := server.NewDefaultTopologyRegistry()

	router := server.SetUpRouter("/", server.ContextGlobalVariables{
		Logger:     logger,
		Topologies: topologies,
	})
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
		Addr:    bind,
		Handler: router,
	}
	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot start the server:", err)
		panic(1)
	}
}

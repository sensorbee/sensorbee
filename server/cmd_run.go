package server

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gocraft/web"
	"net/http"
	"os"
	"pfi/sensorbee/sensorbee/server/api"
	"strconv"
	"strings"
	"sync"
)

// SetUpRunCommand sets up SensorBee's HTTP server.
// The URL or port ID is set with server configuration file,
// or command line arguments.
func SetUpRunCommand() cli.Command {
	cmd := cli.Command{
		Name:        "run",
		Usage:       "run the server",
		Description: "run command starts a new server process",
		Action:      RunRun,
	}
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "port",
			Value:  "8090",
			Usage:  "server port number",
			EnvVar: "PORT",
		},
	}
	return cmd
}

// RunRun HTTP server.
func RunRun(c *cli.Context) {

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

	root := api.SetUpRouterWithCustomMiddleware("/", nil,
		func(c *api.Context, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
			next(rw, req)
		},
		func(prefix string, r *web.Router) {
			api.SetUpAPIRouter(prefix, r)
		})

	handler := func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			root.ServeHTTP(rw, r)
		}
	}

	port, err := strconv.Atoi(c.String("port"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get port number:", err)
	}

	mutex := &sync.Mutex{}
	cond := sync.NewCond(mutex)
	var serverErr error

	mutex.Lock()
	defer mutex.Unlock()

	ports := []int{port} // TODO do need to have several port??
	for _, p := range ports {
		p := p // create a copy of the loop variable for the closure below

		go func() {
			// TODO: We need to listen first, and then serve on it.
			s := &http.Server{
				Addr:    fmt.Sprint(":", p), // TODO Support bind
				Handler: http.HandlerFunc(handler),
			}

			err := s.ListenAndServe()
			if err != nil {
				mutex.Lock()
				defer mutex.Unlock()
				serverErr = err
				cond.Signal()
			}
		}()
	}

	cond.Wait()
	if serverErr != nil {
		fmt.Fprintln(os.Stderr, "Cannot start the server:", serverErr)
	}
}

package server

import (
	"github.com/gocraft/web"
	"os"
	"os/user"
	"runtime"
)

type serverStatus struct {
	*APIContext
}

func SetUpServerStatusRouter(prefix string, router *web.Router) {
	root := router.Subrouter(serverStatus{}, "")
	root.Get("/runtime_status", (*serverStatus).RuntimeStatus)
}

func (ss *serverStatus) RuntimeStatus(rw web.ResponseWriter, req *web.Request) {
	res := map[string]interface{}{
		"num_goroutine": runtime.NumGoroutine(),
		"num_cgo_call":  runtime.NumCgoCall(),
		"gomaxprocs":    runtime.GOMAXPROCS(0),
		"goroot":        runtime.GOROOT(),
		"num_cpu":       runtime.NumCPU(),
		"goversion":     runtime.Version(),
		"pid":           os.Getpid(),
	}
	if dir, err := os.Getwd(); err == nil {
		res["current_directory"] = dir
	}
	if host, err := os.Hostname(); err == nil {
		res["hostname"] = host
	}
	if user, err := user.Current(); err == nil {
		res["user"] = user.Username
	}
	ss.RenderJSON(res)
}

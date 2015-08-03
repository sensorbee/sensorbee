package server

import (
	"github.com/gocraft/web"
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
	ss.RenderJSON(map[string]interface{}{
		"num_goroutine": runtime.NumGoroutine(),
		"num_cgo_call":  runtime.NumCgoCall(),
		"gomaxprocs":    runtime.GOMAXPROCS(0),
		"goroot":        runtime.GOROOT(),
		"num_cpu":       runtime.NumCPU(),
		"goversion":     runtime.Version(),
	})
}

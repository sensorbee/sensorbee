package server

import (
	"github.com/gocraft/web"
	"runtime"
)

type serverStatus struct {
	*APIContext
}

func SetUpServerStatusRouter(prefix string, router *web.Router) {
	root := router.Subrouter(serverStatus{}, "/runtime_status")
	root.Get("/", (*serverStatus).Index)
}

func (ss *serverStatus) Index(rw web.ResponseWriter, req *web.Request) {
	ss.RenderJSON(map[string]interface{}{
		"num_goroutine": runtime.NumGoroutine(),
		"num_cgo_call":  runtime.NumCgoCall(),
	})
}

package api

import (
	"github.com/gocraft/web"
)

type StreamsContext struct {
	*TopologiesContext
	nodeName string
}

func SetUpStreamsRouter(prefix string, router *web.Router) {
	root := router.Subrouter(StreamsContext{}, `/:tenantName/streams/:nodeName`)
	root.Middleware((*StreamsContext).extractNodeName)

	root.Get("/", (*StreamsContext).Show)
	root.Put("/", (*StreamsContext).Update)
	root.Delete("/", (*StreamsContext).Delete)
}

func (sc *StreamsContext) extractNodeName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := sc.extractOptionStringFromPath("nodeName", &sc.nodeName); err != nil {
		return
	}
	next(rw, req)
}

func (sc *StreamsContext) Show(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

func (sc *StreamsContext) Update(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

func (sc *StreamsContext) Delete(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

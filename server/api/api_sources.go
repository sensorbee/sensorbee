package api

import (
	"github.com/gocraft/web"
)

type SourcesContext struct {
	*TopologiesContext
	nodeName string
}

func SetUpSourcesRouter(prefix string, router *web.Router) {
	root := router.Subrouter(SourcesContext{}, `/:tenantName/sources/:nodeName`)
	root.Middleware((*SourcesContext).extractNodeName)

	root.Get("/", (*SourcesContext).Show)
	root.Put("/", (*SourcesContext).Update)
	root.Delete("/", (*SourcesContext).Delete)
}

func (sc *SourcesContext) extractNodeName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := sc.extractOptionStringFromPath("nodeName", &sc.nodeName); err != nil {
		return
	}
	next(rw, req)
}

func (sc *SourcesContext) Show(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

func (sc *SourcesContext) Update(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

func (sc *SourcesContext) Delete(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

package api

import (
	"github.com/gocraft/web"
)

type SinksContext struct {
	*TopologiesContext
	nodeName string
}

func SetUpSinksRouter(prefix string, router *web.Router) {
	root := router.Subrouter(SinksContext{}, `/:tenantName/sinks/:nodeName`)
	root.Middleware((*SinksContext).extractNodeName)

	root.Get("/", (*SinksContext).Show)
	root.Put("/", (*SinksContext).Update)
	root.Delete("/", (*SinksContext).Delete)
}

func (sc *SinksContext) extractNodeName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := sc.extractOptionStringFromPath("nodeName", &sc.nodeName); err != nil {
		return
	}
	next(rw, req)
}

func (sc *SinksContext) Show(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

func (sc *SinksContext) Update(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

func (sc *SinksContext) Delete(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(&map[string]interface{}{
		"topology name": sc.tenantName,
		"node name":     sc.nodeName,
		"status":        "under the construction",
	})
}

package server

import (
	"github.com/gocraft/web"
	"gopkg.in/pfnet/jasco.v1"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/server/response"
	"net/http"
)

type sources struct {
	*topologies
	src core.SourceNode
}

func setUpSourcesRouter(prefix string, router *web.Router) {
	root := router.Subrouter(sources{}, "/:topologyName/sources")
	root.Middleware((*sources).fetchSource)
	root.Get("/", (*sources).Index)
	root.Get("/:sourceName", (*sources).Show)
}

func (sc *sources) fetchSource(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	tb := sc.fetchTopology()
	if tb == nil {
		return
	}

	if srcName := sc.PathParams().String("sourceName", ""); srcName != "" {
		src, err := tb.Topology().Source(srcName)
		if err != nil {
			sc.ErrLog(err).Error("Cannot find the source")
			sc.RenderError(jasco.NewError(requestResourceNotFoundErrorCode,
				"The source was not found", http.StatusNotFound, err))
			return
		}
		sc.src = src
		sc.AddLogField("node_type", core.NTSource.String())
		sc.AddLogField("node_name", src.Name())
	}
	next(rw, req)
}

func (sc *sources) Index(rw web.ResponseWriter, req *web.Request) {
	// TODO: support pagination

	srcs := sc.topology.Topology().Sources()
	res := make([]*response.Source, 0, len(srcs))
	for _, s := range srcs {
		res = append(res, response.NewSource(s, false))
	}
	sc.Render(map[string]interface{}{
		"topology": sc.topologyName,
		"count":    len(res),
		"sources":  res,
	})
}

func (sc *sources) Show(rw web.ResponseWriter, req *web.Request) {
	sc.Render(map[string]interface{}{
		"topology": sc.topologyName,
		"source":   response.NewSource(sc.src, true),
	})
}

// TODO: Support Update(e.g. pause/resume) and Destroy if necessary. They can be
// done by queries.

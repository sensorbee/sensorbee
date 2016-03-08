package server

import (
	"github.com/gocraft/web"
	"gopkg.in/pfnet/jasco.v1"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/server/response"
	"net/http"
)

type sinks struct {
	*topologies
	sink core.SinkNode
}

func setUpSinksRouter(prefix string, router *web.Router) {
	root := router.Subrouter(sinks{}, "/:topologyName/sinks")
	root.Middleware((*sinks).fetchSink)
	root.Get("/", (*sinks).Index)
	root.Get("/:sinkName", (*sinks).Show)
}

func (sc *sinks) fetchSink(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	tb := sc.fetchTopology()
	if tb == nil {
		return
	}

	if sinkName := sc.PathParams().String("sinkName", ""); sinkName != "" {
		sink, err := tb.Topology().Sink(sinkName)
		if err != nil {
			sc.ErrLog(err).Error("Cannot find the sink")
			sc.RenderError(jasco.NewError(requestResourceNotFoundErrorCode,
				"The sink was not found", http.StatusNotFound, err))
			return
		}
		sc.sink = sink
		sc.AddLogField("node_type", core.NTSink.String())
		sc.AddLogField("node_name", sink.Name())
	}
	next(rw, req)
}

func (sc *sinks) Index(rw web.ResponseWriter, req *web.Request) {
	// TODO: support pagination

	sinks := sc.topology.Topology().Sinks()
	res := make([]*response.Sink, 0, len(sinks))
	for _, s := range sinks {
		res = append(res, response.NewSink(s, false))
	}
	sc.Render(map[string]interface{}{
		"topology": sc.topologyName,
		"count":    len(res),
		"sinks":    res,
	})
}

func (sc *sinks) Show(rw web.ResponseWriter, req *web.Request) {
	sc.Render(map[string]interface{}{
		"topology": sc.topologyName,
		"sink":     response.NewSink(sc.sink, true),
	})
}

// TODO: Support Update(e.g. pause/resume) and Destroy if necessary. They can be
// done by queries.

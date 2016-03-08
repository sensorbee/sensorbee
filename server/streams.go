package server

import (
	"github.com/gocraft/web"
	"gopkg.in/pfnet/jasco.v1"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/server/response"
	"net/http"
)

type streams struct {
	*topologies
	stream core.BoxNode
}

func setUpStreamsRouter(prefix string, router *web.Router) {
	root := router.Subrouter(streams{}, "/:topologyName/streams")
	root.Middleware((*streams).fetchStream)
	root.Get("/", (*streams).Index)
	root.Get("/:streamName", (*streams).Show)
}

func (sc *streams) fetchStream(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	tb := sc.fetchTopology()
	if tb == nil {
		return
	}

	if strmName := sc.PathParams().String("streamName", ""); strmName != "" {
		strm, err := tb.Topology().Box(strmName)
		if err != nil {
			sc.ErrLog(err).Error("Cannot find the stream")
			sc.RenderError(jasco.NewError(requestResourceNotFoundErrorCode,
				"The stream was not found", http.StatusNotFound, err))
			return
		}
		sc.stream = strm
		sc.AddLogField("node_type", core.NTBox.String())
		sc.AddLogField("node_name", strm.Name())
	}
	next(rw, req)
}

func (sc *streams) Index(rw web.ResponseWriter, req *web.Request) {
	// TODO: support pagination

	strms := sc.topology.Topology().Boxes()
	res := make([]*response.Stream, 0, len(strms))
	for _, s := range strms {
		res = append(res, response.NewStream(s, false))
	}
	sc.Render(map[string]interface{}{
		"topology": sc.topologyName,
		"count":    len(res),
		"streams":  res,
	})
}

func (sc *streams) Show(rw web.ResponseWriter, req *web.Request) {
	sc.Render(map[string]interface{}{
		"topology": sc.topologyName,
		"stream":   response.NewStream(sc.stream, true),
	})
}

// TODO: Support Update(e.g. pause/resume) and Destroy if necessary. They can be
// done by queries.

package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
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

func (sc *streams) ErrLog(err error) *logrus.Entry {
	e := sc.topologies.ErrLog(err).WithField("topology", sc.topologyName)
	if sc.stream == nil {
		return e
	}
	return e.WithFields(logrus.Fields{
		"node_type": core.NTBox.String(),
		"node_name": sc.stream.Name(),
	})
}

func (sc *streams) fetchStream(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	tb := sc.fetchTopology()
	if tb == nil {
		return
	}

	var strmName string
	if err := sc.extractOptionStringFromPath("streamName", &strmName); err != nil {
		return
	}
	if strmName != "" {
		strm, err := tb.Topology().Box(strmName)
		if err != nil {
			sc.ErrLog(err).Error("Cannot find the stream")
			sc.RenderErrorJSON(NewError(requestURLNotFoundErrorCode, "The stream was not found",
				http.StatusNotFound, err))
			return
		}
		sc.stream = strm
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
	sc.RenderJSON(map[string]interface{}{
		"topology": sc.topologyName,
		"count":    len(res),
		"streams":  res,
	})
}

func (sc *streams) Show(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(map[string]interface{}{
		"topology": sc.topologyName,
		"stream":   response.NewStream(sc.stream, true),
	})
}

// TODO: Support Update(e.g. pause/resume) and Destroy if necessary. They can be
// done by queries.

package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
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
	root.Get("/:streamName", (*sinks).Show)
}

func (sc *sinks) ErrLog(err error) *logrus.Entry {
	e := sc.topologies.ErrLog(err).WithField("topology", sc.topologyName)
	if sc.sink == nil {
		return e
	}
	return e.WithFields(logrus.Fields{
		"node_type": core.NTSink.String(),
		"node_name": sc.sink.Name(),
	})
}

func (sc *sinks) fetchSink(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	tb := sc.fetchTopology()
	if tb == nil {
		return
	}

	var sinkName string
	if err := sc.extractOptionStringFromPath("streamName", &sinkName); err != nil {
		return
	}
	if sinkName != "" {
		sink, err := tb.Topology().Sink(sinkName)
		if err != nil {
			sc.ErrLog(err).Error("Cannot find the sink")
			sc.RenderErrorJSON(NewError(requestURLNotFoundErrorCode, "The sink was not found",
				http.StatusNotFound, err))
			return
		}
		sc.sink = sink
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
	sc.RenderJSON(map[string]interface{}{
		"topology": sc.topologyName,
		"count":    len(res),
		"sinks":    res,
	})
}

func (sc *sinks) Show(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(map[string]interface{}{
		"topology": sc.topologyName,
		"sink":     response.NewSink(sc.sink, true),
	})
}

// TODO: Support Update(e.g. pause/resume) and Destroy if necessary. They can be
// done by queries.

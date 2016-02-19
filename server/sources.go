package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
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

func (sc *sources) ErrLog(err error) *logrus.Entry {
	e := sc.topologies.ErrLog(err).WithField("topology", sc.topologyName)
	if sc.src == nil {
		return e
	}
	return e.WithFields(logrus.Fields{
		"node_type": core.NTSource.String(),
		"node_name": sc.src.Name(),
	})
}

func (sc *sources) fetchSource(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	tb := sc.fetchTopology()
	if tb == nil {
		return
	}

	var srcName string
	if err := sc.extractOptionStringFromPath("sourceName", &srcName); err != nil {
		return
	}
	if srcName != "" {
		src, err := tb.Topology().Source(srcName)
		if err != nil {
			sc.ErrLog(err).Error("Cannot find the source")
			sc.RenderErrorJSON(NewError(requestURLNotFoundErrorCode, "The source was not found",
				http.StatusNotFound, err))
			return
		}
		sc.src = src
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
	sc.RenderJSON(map[string]interface{}{
		"topology": sc.topologyName,
		"count":    len(res),
		"sources":  res,
	})
}

func (sc *sources) Show(rw web.ResponseWriter, req *web.Request) {
	sc.RenderJSON(map[string]interface{}{
		"topology": sc.topologyName,
		"source":   response.NewSource(sc.src, true),
	})
}

// TODO: Support Update(e.g. pause/resume) and Destroy if necessary. They can be
// done by queries.

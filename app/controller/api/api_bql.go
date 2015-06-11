package api

import (
	"github.com/gocraft/web"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
)

// TODO multi tenancy
var (
	topologyBuilder = bql.NewTopologyBuilder()
)

type BQLContext struct {
	*APIContext
	query string
}

func SetUpBQLRouter(prefix string, router *web.Router) {
	root := router.Subrouter(BQLContext{}, "/bql")
	root.Middleware((*BQLContext).extractQuery)
	root.Get("/", (*BQLContext).Index)
	root.Get(`/:query`, (*BQLContext).Show) // TODO validation with regex
	root.Post(`/:query`, (*BQLContext).Update)
	root.Get("/run", (*BQLContext).Run)
}

func (b *BQLContext) extractQueryFromPath(key string, query *string) error {
	s, ok := b.request.PathParams[key]
	if !ok {
		return nil
	}

	*query = s
	return nil
}

func (b *BQLContext) extractQuery(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := b.extractQueryFromPath("query", &b.query); err != nil {
		return
	}
	next(rw, req)
}

func (b *BQLContext) Index(rw web.ResponseWriter, req *web.Request) {
	b.RenderJSON(&map[string]interface{}{
		"status": "ok",
	})
}

func (b *BQLContext) Show(rw web.ResponseWriter, req *web.Request) {
	b.RenderJSON(&map[string]interface{}{
		"query": b.query,
	})
}

func (b *BQLContext) Update(rw web.ResponseWriter, req *web.Request) {
	res := map[string]interface{}{}
	if b.query == "" {
		res["status"] = "nothing to do"
		b.RenderJSON(&res)
		return
	}

	err := topologyBuilder.BQL(b.query)
	if err != nil {
		res["status"] = "fail"
		res["error"] = err.Error()
	} else {
		res["status"] = "register done"
	}
	b.RenderJSON(&res)
}

func (b *BQLContext) Run(rw web.ResponseWriter, req *web.Request) {
	res := map[string]interface{}{}

	tp, err := topologyBuilder.Build()
	if err != nil {
		res["status"] = "fail"
		res["error"] = err.Error()
		b.RenderJSON(&res)
		return
	}

	logManagment := core.NewConsolePrintLogger()
	conf := core.Configuration{
		TupleTraceEnabled: 1,
	}
	ctx := core.Context{
		Logger: logManagment,
		Config: conf,
	}

	err = tp.Run(&ctx)
	if err != nil {
		res["status"] = "fail"
		res["error"] = err.Error()
	} else {
		res["status"] = "running"
	}
	b.RenderJSON(&res)
}

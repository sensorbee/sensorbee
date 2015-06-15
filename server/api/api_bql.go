package api

import (
	"github.com/gocraft/web"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
)

var (
	topologyBuilderMap = map[string]*bql.TopologyBuilder{}
	topologyMap        = map[string]core.StaticTopology{}
	topologyContextMap = map[string]core.Context{}
)

type BQLContext struct {
	*APIContext
	tenant string
	query  string
}

func SetUpBQLRouter(prefix string, router *web.Router) {
	root := router.Subrouter(BQLContext{}, "/bql")
	root.Middleware((*BQLContext).extractQuery)
	root.Middleware((*BQLContext).extractTenant)
	// TODO validation (root can validate with regex like "\w+")
	root.Get("/", (*BQLContext).Index)
	root.Get(`/:tenant`, (*BQLContext).ShowTenant)
	root.Post(`/:tenant`, (*BQLContext).Init)
	root.Get(`/:tenant/:query`, (*BQLContext).ShowQuery)
	root.Post(`/:tenant/:query`, (*BQLContext).Eval)
	root.Post(`/:tenant/run`, (*BQLContext).Run)
	root.Post(`/:tenant/stop`, (*BQLContext).Stop)
}

func (b *BQLContext) extractStringFromPath(key string, target *string) error {
	s, ok := b.request.PathParams[key]
	if !ok {
		return nil
	}

	*target = s
	return nil
}

func (b *BQLContext) extractTenant(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := b.extractStringFromPath("tenant", &b.tenant); err != nil {
		return
	}
	next(rw, req)
}

func (b *BQLContext) extractQuery(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := b.extractStringFromPath("query", &b.query); err != nil {
		return
	}
	next(rw, req)
}

func (b *BQLContext) Index(rw web.ResponseWriter, req *web.Request) {
	tenants := []string{}
	for k, _ := range topologyBuilderMap {
		tenants = append(tenants, k)
	}
	b.RenderJSON(&map[string]interface{}{
		"status":  "ok",
		"tenants": tenants,
	})
}

func (b *BQLContext) ShowTenant(rw web.ResponseWriter, req *web.Request) {
	_, ok := topologyBuilderMap[b.tenant]
	var status string
	if !ok {
		status = "not initialized"
	} else {
		status = "initialized"
	}
	b.RenderJSON(&map[string]interface{}{
		"status": status,
		"tenant": b.tenant,
	})
}

func (b *BQLContext) Init(rw web.ResponseWriter, req *web.Request) {
	_, ok := topologyBuilderMap[b.tenant]
	var status string
	if !ok {
		tb := bql.NewTopologyBuilder()
		topologyBuilderMap[b.tenant] = tb
		status = "initialized"
	} else {
		status = "already initialized"
	}
	b.RenderJSON(&map[string]interface{}{
		"status": status,
		"tenant": b.tenant,
	})
}

// ShowQuery is just for debug
func (b *BQLContext) ShowQuery(rw web.ResponseWriter, req *web.Request) {
	b.RenderJSON(&map[string]interface{}{
		"status": "ok",
		"tenant": b.tenant,
		"query":  b.query,
	})
}

func (b *BQLContext) Eval(rw web.ResponseWriter, req *web.Request) {
	res := map[string]interface{}{}
	res["tenant"] = b.tenant
	if b.tenant == "" || b.query == "" {
		res["status"] = "none"
		b.RenderJSON(&res)
		return
	}

	tb, ok := topologyBuilderMap[b.tenant]
	if !ok {
		tb = bql.NewTopologyBuilder()
		topologyBuilderMap[b.tenant] = tb
	}
	err := tb.BQL(b.query)
	if err != nil {
		res["status"] = "fail"
		res["error"] = err.Error()
		res["query"] = b.query
	} else {
		res["status"] = "register done"
	}
	b.RenderJSON(&res)
}

func (b *BQLContext) Run(rw web.ResponseWriter, req *web.Request) {
	res := map[string]interface{}{}
	res["tenant"] = b.tenant
	if b.tenant == "" {
		res["status"] = "none"
		b.RenderJSON(&res)
		return
	}

	tb, ok := topologyBuilderMap[b.tenant]
	if !ok {
		tb = bql.NewTopologyBuilder() // occur error when Build()
	}
	tp, err := tb.Build()
	if err != nil {
		res["status"] = "fail"
		res["error"] = err.Error()
		b.RenderJSON(&res)
		return
	}
	topologyMap[b.tenant] = tp

	logManagment := core.NewConsolePrintLogger()
	conf := core.Configuration{
		TupleTraceEnabled: 1,
	}
	ctx := core.Context{
		Logger: logManagment,
		Config: conf,
	}
	topologyContextMap[b.tenant] = ctx

	err = tp.Run(&ctx) // TODO use goroutine??
	if err != nil {
		res["status"] = "fail"
		res["error"] = err.Error()
	} else {
		res["status"] = "success (running)"
	}
	b.RenderJSON(&res)
}

func (b *BQLContext) Stop(rw web.ResponseWriter, req *web.Request) {
	res := map[string]interface{}{}
	res["tenant"] = b.tenant
	if b.tenant == "" {
		res["status"] = "none"
		b.RenderJSON(&res)
		return
	}

	tp, ok := topologyMap[b.tenant]
	if !ok {
		res["status"] = "topology is not initialized"
		b.RenderJSON(&res)
		return
	}
	ctx := topologyContextMap[b.tenant]
	err := tp.Stop(&ctx)
	if err != nil {
		res["status"] = "fail"
		res["error"] = err.Error()
	} else {
		res["status"] = "success"
	}
	b.RenderJSON(&res)
}

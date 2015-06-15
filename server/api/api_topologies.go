package api

import (
	"fmt"
	"github.com/gocraft/web"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
)

var (
	topologyBuilderMap = map[string]*bql.TopologyBuilder{}
	topologyMap        = map[string]core.StaticTopology{}
	topologyContextMap = map[string]core.Context{}
)

type TopologiesContext struct {
	*APIContext
	name string
}

func SetUpBQLRouter(prefix string, router *web.Router) {
	root := router.Subrouter(TopologiesContext{}, "/topologies")
	root.Middleware((*TopologiesContext).extractName)
	// TODO validation (root can validate with regex like "\w+")
	root.Get("/", (*TopologiesContext).Index)
	root.Get(`/:name`, (*TopologiesContext).Show)
	root.Put(`/:name`, (*TopologiesContext).Update)
	root.Post(`/:name/queries`, (*TopologiesContext).Queries)
}

func (tc *TopologiesContext) extractStringFromPath(key string, target *string) error {
	s, ok := tc.request.PathParams[key]
	if !ok {
		return nil
	}

	*target = s
	return nil
}

func (tc *TopologiesContext) extractName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := tc.extractStringFromPath("name", &tc.name); err != nil {
		return
	}
	next(rw, req)
}

// Index returns registered tenant name list
func (tc *TopologiesContext) Index(rw web.ResponseWriter, req *web.Request) {
	tenants := []string{}
	for k, _ := range topologyBuilderMap {
		tenants = append(tenants, k)
	}
	tc.RenderJSON(&map[string]interface{}{
		"tenants": tenants,
	})
}

// Show returns the information of topology
func (tc *TopologiesContext) Show(rw web.ResponseWriter, req *web.Request) {
	_, ok := topologyBuilderMap[tc.name]
	var status string
	if !ok {
		status = "not initialized"
	} else {
		status = "initialized"
	}
	tc.RenderJSON(&map[string]interface{}{
		"name":          tc.name,
		"status":        status,
		"topology info": "", // TODO add topology information
	})
}

// Update execute BQLs
func (tc *TopologiesContext) Update(rw web.ResponseWriter, req *web.Request) {
	tp, ok := topologyMap[tc.name]
	if !ok {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": "not initialized or runnning",
		})
		return
	}

	s, ok := tc.request.PathParams["state"]
	if !ok {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": "cannot get updating operation",
		})
		return
	}

	ctx := topologyContextMap[tc.name]
	var err error
	switch s {
	case "stop":
		err = tp.Stop(&ctx)
	case "pause":
	case "resume":
	default:
		err = fmt.Errorf("cannot update the state: %v", s)
	}
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": err.Error(),
		})
	} else {
		tc.RenderJSON(&map[string]interface{}{
			"name":  tc.name,
			"state": s,
		})
	}
}

func (tc *TopologiesContext) Queries(rw web.ResponseWriter, req *web.Request) {
	queries, ok := tc.request.PathParams["queries"]
	if !ok {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": "cannot support to execute empty query",
		})
		return
	}

	// TODO need to apply dynamic topology
	var tb *bql.TopologyBuilder
	var ctx core.Context
	tp, ok := topologyMap[tc.name]
	if !ok {
		tb = bql.NewTopologyBuilder()
		topologyBuilderMap[tc.name] = tb

		// TODO get context configuration from BQL
		logManagment := core.NewConsolePrintLogger()
		conf := core.Configuration{
			TupleTraceEnabled: 1,
		}
		ctx := core.Context{
			Logger: logManagment,
			Config: conf,
		}
		topologyContextMap[tc.name] = ctx
	} else {
		// if already exist topology, means that topology builder also exist
		tb = topologyBuilderMap[tc.name]
		ctx = topologyContextMap[tc.name]
	}

	err := tb.BQL(queries)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": err.Error(),
		})
		return
	}

	tp, err = tb.Build()
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": err.Error(),
		})
		return
	}

	err = tp.Run(&ctx)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": err.Error(),
		})
	} else {
		tc.RenderJSON(&map[string]interface{}{
			"name":    tc.name,
			"status":  "running",
			"queries": queries,
		})
	}
}

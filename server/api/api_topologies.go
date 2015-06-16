package api

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"strconv"
)

var (
	topologyBuilderMap = map[string]*bql.TopologyBuilder{}
	topologyMap        = map[string]core.StaticTopology{}
	topologyContextMap = map[string]core.Context{}
)

type TopologiesContext struct {
	*APIContext
	name   string
	nodeId int64
}

func SetUpBQLRouter(prefix string, router *web.Router) {
	root := router.Subrouter(TopologiesContext{}, "/topologies")
	root.Middleware((*TopologiesContext).extractName)
	root.Middleware((*TopologiesContext).extractNodeId)
	// TODO validation (root can validate with regex like "\w+")
	root.Get("/", (*TopologiesContext).Index)
	root.Get(`/:name`, (*TopologiesContext).Show)
	root.Put(`/:name`, (*TopologiesContext).Update)
	root.Post(`/:name/queries`, (*TopologiesContext).Queries)

	// node controller TODO separate to another file
	root.Get(`/:name/sources/:id`, (*TopologiesContext).ShowSources)
	root.Get(`/:name/streams/:id`, (*TopologiesContext).ShowStreams)
	root.Get(`/:name/sinks/:id`, (*TopologiesContext).ShowSinks)
	root.Put(`/:name/sources/:id`, (*TopologiesContext).UpdateSources)
	root.Put(`/:name/streams/:id`, (*TopologiesContext).UpdateStreams)
	root.Put(`/:name/sinks/:id`, (*TopologiesContext).UpdateSinks)
	root.Delete(`/:name/sources/:id`, (*TopologiesContext).DeleteSources)
	root.Delete(`/:name/streams/:id`, (*TopologiesContext).DeleteStreams)
	root.Delete(`/:name/sinks/:id`, (*TopologiesContext).DeleteSinks)
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

func (tc *TopologiesContext) extractIntFromPath(key string, target *int64) error {
	i, ok := tc.request.PathParams[key]
	if !ok {
		return nil
	}

	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		return nil
	}
	*target = id
	return nil
}

func (tc *TopologiesContext) extractNodeId(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := tc.extractIntFromPath("id", &tc.nodeId); err != nil {
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

	// TODO should use ParseJSONFromRequestBoty (util.go)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": err.Error(),
		})
		return
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":       tc.name,
			"query byte": string(b),
			"status":     err.Error(),
		})
		return
	}

	state, ok := m["state"].(string)
	if !ok {
		state = ""
	}
	ctx := topologyContextMap[tc.name]
	switch state {
	case "stop":
		err = tp.Stop(&ctx)
	case "pause":
	case "resume":
	default:
		err = fmt.Errorf("cannot update the state: %v", state)
	}
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": err.Error(),
		})
	} else {
		tc.RenderJSON(&map[string]interface{}{
			"name":  tc.name,
			"state": state,
		})
	}
}

func (tc *TopologiesContext) Queries(rw web.ResponseWriter, req *web.Request) {
	// TODO should use ParseJSONFromRequestBoty (util.go)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": err.Error(),
		})
		return
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":       tc.name,
			"query byte": string(b),
			"status":     err.Error(),
		})
		return
	}
	queries, ok := m["queries"].(string)
	if !ok || queries == "" {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.name,
			"status": "not support to execute empty query",
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

	err = tb.BQL(queries)
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

func (tc *TopologiesContext) ShowSources(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) ShowStreams(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) ShowSinks(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) UpdateSources(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) UpdateStreams(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) UpdateSinks(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) DeleteSources(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) DeleteStreams(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

func (tc *TopologiesContext) DeleteSinks(rw web.ResponseWriter, req *web.Request) {
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.name,
		"node id": tc.nodeId,
		"status":  "under the construction",
	})
}

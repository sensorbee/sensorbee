package api

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
)

var (
	topologyMap        = map[string]core.DynamicTopology{}
	topologyBuilderMap = map[string]*bql.TopologyBuilder{}
	bqlParser          = parser.NewBQLParser()
)

type TopologiesContext struct {
	*APIContext
	tenantName string
}

func SetUpTopologiesRouter(prefix string, router *web.Router) {
	root := router.Subrouter(TopologiesContext{}, "/topologies")
	root.Middleware((*TopologiesContext).extractName)
	// TODO validation (root can validate with regex like "\w+")
	root.Get("/", (*TopologiesContext).Index)
	root.Get(`/:tenantName`, (*TopologiesContext).Show)
	root.Put(`/:tenantName`, (*TopologiesContext).Update)
	root.Post(`/:tenantName/queries`, (*TopologiesContext).Queries)

	SetUpSourcesRouter(prefix, root)
	SetUpStreamsRouter(prefix, root)
	SetUpSinksRouter(prefix, root)
}

func (tc *TopologiesContext) extractName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := tc.extractOptionStringFromPath("tenantName", &tc.tenantName); err != nil {
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
		"topologies": tenants,
	})
}

// Show returns the information of topology
func (tc *TopologiesContext) Show(rw web.ResponseWriter, req *web.Request) {
	_, ok := topologyBuilderMap[tc.tenantName]
	var status string
	if !ok {
		status = "not initialized"
	} else {
		status = "initialized"
	}
	tc.RenderJSON(&map[string]interface{}{
		"name":          tc.tenantName,
		"status":        status,
		"topology info": "", // TODO add topology information
	})
}

// Update dynamic nodes by BQLs
func (tc *TopologiesContext) Update(rw web.ResponseWriter, req *web.Request) {
	tp, ok := topologyMap[tc.tenantName]
	if !ok {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.tenantName,
			"status": "not initialized or running topology",
		})
		return
	}

	// TODO should use ParseJSONFromRequestBoty (util.go)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.tenantName,
			"status": err.Error(),
		})
		return
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":       tc.tenantName,
			"query byte": string(b),
			"status":     err.Error(),
		})
		return
	}

	state, ok := m["state"].(string)
	if !ok {
		state = ""
	}
	switch state {
	case "stop":
		err = tp.Stop()
	case "pause":
	case "resume":
	default:
		err = fmt.Errorf("cannot update the state: %v", state)
	}
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.tenantName,
			"status": err.Error(),
		})
	} else {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.tenantName,
			"status": "done",
			"state":  state,
		})
	}
}

func (tc *TopologiesContext) Queries(rw web.ResponseWriter, req *web.Request) {
	// TODO should use ParseJSONFromRequestBoty (util.go)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.tenantName,
			"status": err.Error(),
		})
		return
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":       tc.tenantName,
			"query byte": string(b),
			"status":     err.Error(),
		})
		return
	}
	queries, ok := m["queries"].(string)
	if !ok || queries == "" {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.tenantName,
			"status": "not support to execute empty query",
		})
		return
	}

	tb, ok := topologyBuilderMap[tc.tenantName]
	if !ok {
		// TODO get context configuration from BQL
		logManagment := core.NewConsolePrintLogger()
		conf := core.Configuration{
			TupleTraceEnabled: 1,
		}
		ctx := core.Context{
			Logger: logManagment,
			Config: conf,
		}
		tp := core.NewDefaultDynamicTopology(&ctx, tc.tenantName)
		topologyMap[tc.tenantName] = tp

		tb, _ = bql.NewTopologyBuilder(tp) // TODO: fix this by supporting Create action
		topologyBuilderMap[tc.tenantName] = tb

	}

	stmts, err := bqlParser.ParseStmts(queries)
	if err != nil {
		tc.RenderJSON(&map[string]interface{}{
			"name":   tc.tenantName,
			"status": err.Error(),
		})
		return
	}
	for _, stmt := range stmts {
		_, err = tb.AddStmt(stmt) // TODO node identifier
		if err != nil {
			tc.RenderJSON(&map[string]interface{}{
				"name":   tc.tenantName,
				"status": err.Error(),
			})
			return // TODO return error detail
		}
	}
	tc.RenderJSON(&map[string]interface{}{
		"name":    tc.tenantName,
		"status":  "running",
		"queries": queries,
	})
}

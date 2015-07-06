package api

import (
	"encoding/json"
	"github.com/gocraft/web"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type TopologiesContext struct {
	*APIContext
	tenantName string
}

func SetUpTopologiesRouter(prefix string, router *web.Router) {
	root := router.Subrouter(TopologiesContext{}, "/topologies")
	root.Middleware((*TopologiesContext).extractName)
	// TODO validation (root can validate with regex like "\w+")
	root.Post("/", (*TopologiesContext).Create)
	root.Get("/", (*TopologiesContext).Index)
	root.Get(`/:tenantName`, (*TopologiesContext).Show)
	root.Delete(`/:tenantName`, (*TopologiesContext).Destroy)
	root.Post(`/:tenantName/queries`, (*TopologiesContext).Queries)
}

func (tc *TopologiesContext) extractName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := tc.extractOptionStringFromPath("tenantName", &tc.tenantName); err != nil {
		return
	}
	next(rw, req)
}

// Create creates a new topology.
func (tc *TopologiesContext) Create(rw web.ResponseWriter, req *web.Request) {
	js, ctxErr := ParseJSONFromRequestBody(tc.Context)
	if ctxErr != nil {
		// TODO: log and render error json
		return
	}

	// TODO: use mapstructure when parameters get too many
	form, err := data.NewMap(js)
	if err != nil {
		// TODO: log and render error json
		return
	}

	n, ok := form["name"]
	if !ok {
		// TODO: log and render error json
		return
	}
	name, err := data.AsString(n)
	if err != nil {
		// TODO: log and render error json
		return
	}

	// TODO: support other parameters

	conf := core.Configuration{
		TupleTraceEnabled: 0,
	}
	ctx := core.Context{
		Logger:       core.NewConsolePrintLogger(), // TODO: set appropriate logger
		Config:       conf,
		SharedStates: core.NewDefaultSharedStateRegistry(),
	}
	tp := core.NewDefaultTopology(&ctx, tc.tenantName)
	tb, err := bql.NewTopologyBuilder(tp)
	if err != nil {
		// TODO: log and render error json
		return
	}

	if err := tc.topologies.Register(name, tb); err != nil {
		// TODO: log and render error json
		return
	}
	tc.RenderJSON(map[string]interface{}{
		"topology": map[string]interface{}{
			"name": name,
		},
	})
}

// Index returns registered tenant name list
func (tc *TopologiesContext) Index(rw web.ResponseWriter, req *web.Request) {
	tenants := []string{}
	ts, err := tc.topologies.List()
	if err != nil {
		// TODO: log and render error json
		return
	}

	for k, _ := range ts {
		tenants = append(tenants, k)
	}

	// TODO: return some statistics using Show's result
	tc.RenderJSON(map[string]interface{}{
		"topologies": tenants,
	})
}

// Show returns the information of topology
func (tc *TopologiesContext) Show(rw web.ResponseWriter, req *web.Request) {
	_, err := tc.topologies.Lookup(tc.tenantName)
	if err != nil {
		// TODO: log and render error json (404 if the error is "NotFound")
		return
	}

	// TODO: return some statistics
	tc.RenderJSON(map[string]interface{}{
		"topology": map[string]interface{}{
			"name": tc.tenantName,
		},
	})
}

// TODO: provide Update action (change state of the topology, etc.)

func (tc *TopologiesContext) Destroy(rw web.ResponseWriter, req *web.Request) {
	tb, err := tc.topologies.Unregister(tc.tenantName)
	if err != nil {
		// TODO: log and render error json
		return
	}
	if tb != nil {
		if err := tb.Topology().Stop(); err != nil {
			// TODO: log and add warning to json
		}
	}

	// TODO: return warning if the topology didn't stop correctly.
	tc.RenderJSON(map[string]interface{}{})
}

func (tc *TopologiesContext) Queries(rw web.ResponseWriter, req *web.Request) {
	tb, err := tc.topologies.Lookup(tc.tenantName)
	if err != nil {
		// TODO: log and render error json
		return
	}

	// TODO should use ParseJSONFromRequestBoty (util.go)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tc.RenderJSON(map[string]interface{}{
			"name":   tc.tenantName,
			"status": err.Error(),
		})
		return
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		tc.RenderJSON(map[string]interface{}{
			"name":       tc.tenantName,
			"query byte": string(b),
			"status":     err.Error(),
		})
		return
	}
	queries, ok := m["queries"].(string)
	if !ok || queries == "" {
		tc.RenderJSON(map[string]interface{}{
			"name":   tc.tenantName,
			"status": "not support to execute empty query",
		})
		return
	}

	bp := parser.NewBQLParser()
	stmts, err := bp.ParseStmts(queries)
	if err != nil {
		tc.RenderJSON(map[string]interface{}{
			"name":   tc.tenantName,
			"status": err.Error(),
		})
		return
	}
	for _, stmt := range stmts {
		_, err = tb.AddStmt(stmt) // TODO node identifier
		if err != nil {
			tc.RenderJSON(map[string]interface{}{
				"name":   tc.tenantName,
				"status": err.Error(),
			})
			return // TODO return error detail
		}
	}
	tc.RenderJSON(map[string]interface{}{
		"name":    tc.tenantName,
		"status":  "running",
		"queries": queries,
	})
}

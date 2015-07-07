package server

import (
	"encoding/json"
	"github.com/gocraft/web"
	"io/ioutil"
	"net/http"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type topologies struct {
	*APIContext
	topologyName string
}

func SetUpTopologiesRouter(prefix string, router *web.Router) {
	root := router.Subrouter(topologies{}, "/topologies")
	root.Middleware((*topologies).extractName)
	// TODO validation (root can validate with regex like "\w+")
	root.Post("/", (*topologies).Create)
	root.Get("/", (*topologies).Index)
	root.Get(`/:topologyName`, (*topologies).Show)
	root.Delete(`/:topologyName`, (*topologies).Destroy)
	root.Post(`/:topologyName/queries`, (*topologies).Queries)
}

func (tc *topologies) extractName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if err := tc.extractOptionStringFromPath("topologyName", &tc.topologyName); err != nil {
		return
	}
	next(rw, req)
}

// Create creates a new topology.
func (tc *topologies) Create(rw web.ResponseWriter, req *web.Request) {
	js, apiErr := ParseJSONFromRequestBody(tc.Context)
	if apiErr != nil {
		tc.ErrLog(apiErr.Err).Error("Cannot parse the request json")
		tc.RenderErrorJSON(apiErr)
		return
	}

	// TODO: use mapstructure when parameters get too many
	form, err := data.NewMap(js)
	if err != nil {
		tc.ErrLog(err).WithField("body", js).Error("The request json may contain invalid value")
		tc.RenderErrorJSON(NewError(formValidationErrorCode, "The request json may contain invalid values.",
			http.StatusBadRequest, err))
		return
	}

	// TODO: report validation errors at once (don't report each error separately) after adding other parameters

	n, ok := form["name"]
	if !ok {
		tc.Log().Error("The required 'name' field is missing")
		e := NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, nil)
		e.Meta["name"] = []string{"field is missing"}
		tc.RenderErrorJSON(e)
		return
	}
	name, err := data.AsString(n)
	if err != nil {
		tc.ErrLog(err).Error("'name' field isn't a string")
		e := NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, nil)
		e.Meta["name"] = []string{"value must be a string"}
		tc.RenderErrorJSON(e)
		return
	}
	if err := core.ValidateNodeName(name); err != nil {
		tc.ErrLog(err).Error("'name' field has invalid format")
		e := NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, nil)
		e.Meta["name"] = []string{"inavlid format"}
		tc.RenderErrorJSON(e)
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
	tp := core.NewDefaultTopology(&ctx, name)
	tb, err := bql.NewTopologyBuilder(tp)
	if err != nil {
		tc.ErrLog(err).Error("Cannot create a new topology builder")
		tc.RenderErrorJSON(NewInternalServerError(err))
		return
	}

	if err := tc.topologies.Register(name, tb); err != nil {
		if err := tp.Stop(); err != nil {
			tc.ErrLog(err).Error("Cannot stop the created topology")
		}

		l := tc.Log().WithField("topology", name)
		if os.IsExist(err) {
			l.Error("the name is already registered")
			e := NewError(formValidationErrorCode, "The request body is invalid.",
				http.StatusBadRequest, nil)
			e.Meta["name"] = []string{"already taken"}
			tc.RenderErrorJSON(e)
			return
		}
		l.WithField("err", err).Error("Cannot register the topology")
		tc.RenderJSON(NewInternalServerError(err))
		return
	}

	// TODO: return 201
	tc.RenderJSON(map[string]interface{}{
		"topology": newTopologiesShowResult(tb),
	})
}

// Index returned a list of registered topologies.
func (tc *topologies) Index(rw web.ResponseWriter, req *web.Request) {
	ts, err := tc.topologies.List()
	if err != nil {
		tc.ErrLog(err).Error("Cannot list registered topologies")
		tc.RenderErrorJSON(NewInternalServerError(err))
		return
	}

	res := []*topologiesShowResult{}
	for _, tb := range ts {
		res = append(res, newTopologiesShowResult(tb))
	}
	tc.RenderJSON(map[string]interface{}{
		"topologies": res,
	})
}

// Show returns the information of topology
func (tc *topologies) Show(rw web.ResponseWriter, req *web.Request) {
	tb, err := tc.topologies.Lookup(tc.topologyName)
	if err != nil {
		l := tc.Log().WithField("topology", tc.topologyName)
		if os.IsNotExist(err) {
			l.Error("The topology is not registered")
			tc.RenderErrorJSON(NewError(requestURLNotFoundErrorCode, "The topology doesn't exist",
				http.StatusNotFound, err))
			return
		}
		l.WithField("err", err).Error("Cannot lookup the topology")
		tc.RenderErrorJSON(NewInternalServerError(err))
		return
	}
	tc.RenderJSON(map[string]interface{}{
		"topology": newTopologiesShowResult(tb),
	})
}

type topologiesShowResult struct {
	Name string `json:"name"`
	// TODO: add other information
}

func newTopologiesShowResult(tb *bql.TopologyBuilder) *topologiesShowResult {
	return &topologiesShowResult{
		Name: tb.Topology().Name(),
	}
}

// TODO: provide Update action (change state of the topology, etc.)

func (tc *topologies) Destroy(rw web.ResponseWriter, req *web.Request) {
	tb, err := tc.topologies.Unregister(tc.topologyName)
	if err != nil {
		tc.ErrLog(err).WithField("topology", tc.topologyName).Error("Cannot unregister the topology")
		tc.RenderErrorJSON(NewInternalServerError(err))
		return
	}
	stopped := true
	if tb != nil {
		if err := tb.Topology().Stop(); err != nil {
			stopped = false
			tc.ErrLog(err).WithField("topology", tc.topologyName).Error("Cannot stop the topology")
		}
	}

	if stopped {
		// TODO: return 204 when the topology didn't exist.
		tc.RenderJSON(map[string]interface{}{})
	} else {
		tc.RenderJSON(map[string]interface{}{
			"warning": map[string]interface{}{
				"message": "the topology wasn't stopped correctly",
			},
		})
	}
}

func (tc *topologies) Queries(rw web.ResponseWriter, req *web.Request) {
	tb, err := tc.topologies.Lookup(tc.topologyName)
	if err != nil {
		// TODO: log and render error json
		return
	}

	// TODO should use ParseJSONFromRequestBoty (util.go)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tc.RenderJSON(map[string]interface{}{
			"name":   tc.topologyName,
			"status": err.Error(),
		})
		return
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		tc.RenderJSON(map[string]interface{}{
			"name":       tc.topologyName,
			"query byte": string(b),
			"status":     err.Error(),
		})
		return
	}
	queries, ok := m["queries"].(string)
	if !ok || queries == "" {
		tc.RenderJSON(map[string]interface{}{
			"name":   tc.topologyName,
			"status": "not support to execute empty query",
		})
		return
	}

	bp := parser.NewBQLParser()
	stmts, err := bp.ParseStmts(queries)
	if err != nil {
		tc.RenderJSON(map[string]interface{}{
			"name":   tc.topologyName,
			"status": err.Error(),
		})
		return
	}
	for _, stmt := range stmts {
		_, err = tb.AddStmt(stmt) // TODO node identifier
		if err != nil {
			tc.RenderJSON(map[string]interface{}{
				"name":   tc.topologyName,
				"status": err.Error(),
			})
			return // TODO return error detail
		}
	}
	tc.RenderJSON(map[string]interface{}{
		"name":    tc.topologyName,
		"status":  "running",
		"queries": queries,
	})
}

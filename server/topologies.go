package server

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"pfi/sensorbee/sensorbee/server/response"
	"strings"
)

type topologies struct {
	*APIContext
	topologyName string
	topology     *bql.TopologyBuilder
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

// fetchTopology returns the topology having tc.topologyName. When this method
// returns nil, the caller can just return from the action.
func (tc *topologies) fetchTopology() *bql.TopologyBuilder {
	tb, err := tc.topologies.Lookup(tc.topologyName)
	if err != nil {
		l := tc.Log().WithField("topology", tc.topologyName)
		if os.IsNotExist(err) {
			l.Error("The topology is not registered")
			tc.RenderErrorJSON(NewError(requestURLNotFoundErrorCode, "The topology doesn't exist",
				http.StatusNotFound, err))
			return nil
		}
		l.WithField("err", err).Error("Cannot lookup the topology")
		tc.RenderErrorJSON(NewInternalServerError(err))
		return nil
	}
	return tb
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

	ctx := core.NewContext(&core.ContextConfig{
		Logger: tc.logger,
		Flags: core.ContextFlags{
			DroppedTupleLog: 1,
		},
	})
	tp := core.NewDefaultTopology(ctx, name)
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
		"topology": response.NewTopology(tb.Topology()),
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

	res := []*response.Topology{}
	for _, tb := range ts {
		res = append(res, response.NewTopology(tb.Topology()))
	}
	tc.RenderJSON(map[string]interface{}{
		"topologies": res,
	})
}

// Show returns the information of topology
func (tc *topologies) Show(rw web.ResponseWriter, req *web.Request) {
	tb := tc.fetchTopology()
	if tb == nil {
		return
	}
	tc.RenderJSON(map[string]interface{}{
		"topology": response.NewTopology(tb.Topology()),
	})
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
	tb := tc.fetchTopology()
	if tb == nil {
		return
	}

	js, apiErr := ParseJSONFromRequestBody(tc.Context)
	if apiErr != nil {
		tc.ErrLog(apiErr.Err).WithField("topology", tc.topologyName).Error("Cannot parse the request json")
		tc.RenderErrorJSON(apiErr)
		return
	}

	// TODO: use mapstructure when parameters get too many
	form, err := data.NewMap(js)
	if err != nil {
		tc.ErrLog(err).WithField("topology", tc.topologyName).WithField("body", js).
			Error("The request json may contain invalid value")
		tc.RenderErrorJSON(NewError(formValidationErrorCode, "The request json may contain invalid values.",
			http.StatusBadRequest, err))
		return
	}

	var queries string
	if v, ok := form["queries"]; !ok {
		tc.Log().WithField("topology", tc.topologyName).Error("The request json doesn't have 'queries' field")
		tc.RenderErrorJSON(NewError(formValidationErrorCode, "'queries' field is missing",
			http.StatusBadRequest, nil))
		return
	} else if f, err := data.AsString(v); err != nil {
		tc.ErrLog(err).WithField("topology", tc.topologyName).Error("'queries' must be a string")
		tc.RenderErrorJSON(NewError(formValidationErrorCode, "'queries' field must be a string",
			http.StatusBadRequest, err))
		return
	} else {
		queries = f
	}

	bp := parser.NewBQLParser()
	type stmtWithStr struct {
		stmt    interface{}
		stmtStr string // TODO: this should be stmt.String()
	}
	stmts := []*stmtWithStr{}
	selectStmtIndex := -1
	for queries != "" {
		stmt, rest, err := bp.ParseStmt(queries)
		if err != nil {
			tc.Log().WithField("topology", tc.topologyName).WithField("parse_errors", err.Error()).
				WithField("statement", queries).Error("Cannot parse a statement")
			e := NewError(bqlStmtParseErrorCode, "Cannot parse a BQL statement", http.StatusBadRequest, err)
			e.Meta["parse_errors"] = strings.Split(err.Error(), "\n") // FIXME: too ad hoc
			e.Meta["statement"] = queries
			tc.RenderErrorJSON(e)
			return
		}
		if _, ok := stmt.(parser.SelectStmt); ok {
			selectStmtIndex = len(stmts)
		}

		stmts = append(stmts, &stmtWithStr{stmt, queries[:len(queries)-len(rest)]})
		queries = rest
	}

	if selectStmtIndex >= 0 {
		if len(stmts) != 1 {
			tc.Log().WithField("topology", tc.topologyName).Error("A SELECT statement cannot be issued with other statements")
			e := NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
			e.Meta["error"] = "a SELECT statement cannot be issued with other statements"
			e.Meta["statement"] = stmts[selectStmtIndex].stmtStr
			tc.RenderErrorJSON(e)
			return
		}
		tc.handleSelectStmt(rw, stmts[selectStmtIndex].stmt.(parser.SelectStmt), stmts[selectStmtIndex].stmtStr)
		return
	}

	// TODO: handle this atomically
	for _, stmt := range stmts {
		// TODO: change the return value of AddStmt to support the new response format.
		_, err = tb.AddStmt(stmt.stmt)
		if err != nil {
			tc.ErrLog(err).WithField("topology", tc.topologyName).Error("Cannot process a statement")
			e := NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
			e.Meta["error"] = err.Error()
			e.Meta["statement"] = stmt.stmtStr
			tc.RenderErrorJSON(e)
			return
		}
	}

	// TODO: support the new format
	tc.RenderJSON(map[string]interface{}{
		"topology_name": tc.topologyName,
		"status":        "running",
		"queries":       queries,
	})
}

func (tc *topologies) handleSelectStmt(rw web.ResponseWriter, stmt parser.SelectStmt, stmtStr string) {
	tb := tc.fetchTopology()
	if tb == nil { // just in case
		return
	}

	sn, ch, err := tb.AddSelectStmt(&stmt)
	if err != nil {
		tc.ErrLog(err).WithField("topology", tc.topologyName).Error("Cannot process a statement")
		e := NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
		e.Meta["error"] = err.Error()
		e.Meta["statement"] = stmtStr
		tc.RenderErrorJSON(e)
		return
	}
	defer func() {
		if err := sn.Stop(); err != nil {
			tc.ErrLog(err).WithFields(logrus.Fields{
				"topology":  tc.topologyName,
				"node_type": core.NTSink,
				"node_name": sn.Name(),
			}).Error("Cannot stop the temporary sink")
		}
	}()

	conn, bufrw, err := rw.Hijack()
	if err != nil {
		tc.ErrLog(err).WithField("topology", tc.topologyName).Error("Cannot hijack a connection")
		tc.RenderErrorJSON(NewInternalServerError(err))
		return
	}

	var writeErr error
	mw := multipart.NewWriter(bufrw)
	defer func() {
		if writeErr != nil {
			tc.ErrLog(writeErr).WithField("topology", tc.topologyName).Info("Cannot write contents to the hijacked connection")
		}

		if err := mw.Close(); err != nil {
			if writeErr == nil { // log it only when the write err hasn't happend
				tc.ErrLog(err).WithField("topology", tc.topologyName).Info("Cannot finish the multipart response")
			}
		}
		bufrw.Flush()
		conn.Close()

		tc.Log().WithFields(logrus.Fields{
			"topology":  tc.topologyName,
			"statement": stmtStr,
		}).Info("Streaming SELECT responses finished")
	}()

	res := []string{
		"HTTP/1.1 200 OK",
		fmt.Sprintf(`Content-Type: multipart/mixed; boundary="%v"`, mw.Boundary()),
		"\r\n",
	}
	if _, err := bufrw.WriteString(strings.Join(res, "\r\n")); err != nil {
		tc.ErrLog(err).WithField("topology", tc.topologyName).Error("Cannot write a header to the hijacked connection")
		return
	}
	bufrw.Flush()

	tc.Log().WithFields(logrus.Fields{
		"topology":  tc.topologyName,
		"statement": stmtStr,
	}).Info("Start streaming SELECT responses")

	// All error reporting logs after this is info level because they might be
	// caused by the client closing the connection.
	header := textproto.MIMEHeader{}
	header.Add("Content-Type", "application/json")
	for t := range ch {
		// TODO: provide very efficient timeout detection
		w, err := mw.CreatePart(header)
		if err != nil {
			writeErr = err
			return
		}
		if _, err := io.WriteString(w, t.Data.String()); err != nil {
			writeErr = err
			return
		}
		if err := bufrw.Flush(); err != nil {
			writeErr = err
			return
		}
	}
}

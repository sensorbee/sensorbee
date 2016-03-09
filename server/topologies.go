package server

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"golang.org/x/net/websocket"
	"gopkg.in/pfnet/jasco.v1"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"gopkg.in/sensorbee/sensorbee.v0/server/response"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"
)

type topologies struct {
	*APIContext
	topologyName string
	topology     *bql.TopologyBuilder
}

func setUpTopologiesRouter(prefix string, router *web.Router) {
	root := router.Subrouter(topologies{}, "/topologies")
	root.Middleware((*topologies).extractName)
	// TODO validation (root can validate with regex like "\w+")
	root.Post("/", (*topologies).Create)
	root.Get("/", (*topologies).Index)
	root.Get(`/:topologyName`, (*topologies).Show)
	root.Delete(`/:topologyName`, (*topologies).Destroy)
	root.Post(`/:topologyName/queries`, (*topologies).Queries)
	root.Get(`/:topologyName/wsqueries`, (*topologies).WebSocketQueries)

	setUpSourcesRouter(prefix, root)
	setUpStreamsRouter(prefix, root)
	setUpSinksRouter(prefix, root)
}

func (tc *topologies) extractName(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	tc.topologyName = tc.PathParams().String("topologyName", "")
	if tc.topologyName != "" {
		tc.AddLogField("topology", tc.topologyName)
	}
	next(rw, req)
}

// fetchTopology returns the topology having tc.topologyName. When this method
// returns nil, the caller can just return from the action.
func (tc *topologies) fetchTopology() *bql.TopologyBuilder {
	tb, err := tc.topologies.Lookup(tc.topologyName)
	if err != nil {
		if core.IsNotExist(err) {
			tc.Log().Error("The topology is not registered")
			tc.RenderError(jasco.NewError(requestResourceNotFoundErrorCode, "The topology doesn't exist",
				http.StatusNotFound, err))
			return nil
		}
		tc.ErrLog(err).Error("Cannot lookup the topology")
		tc.RenderError(jasco.NewInternalServerError(err))
		return nil
	}
	tc.topology = tb
	return tb
}

// Create creates a new topology.
func (tc *topologies) Create(rw web.ResponseWriter, req *web.Request) {
	var js map[string]interface{}
	if apiErr := tc.ParseBody(&js); apiErr != nil {
		tc.ErrLog(apiErr.Err).Error("Cannot parse the request json")
		tc.RenderError(apiErr)
		return
	}

	// TODO: use mapstructure when parameters get too many
	form, err := data.NewMap(js)
	if err != nil {
		tc.ErrLog(err).WithField("body", js).Error("The request json may contain invalid value")
		tc.RenderError(jasco.NewError(formValidationErrorCode, "The request json may contain invalid values.",
			http.StatusBadRequest, err))
		return
	}

	// TODO: report validation errors at once (don't report each error separately) after adding other parameters

	n, ok := form["name"]
	if !ok {
		tc.Log().Error("The required 'name' field is missing")
		e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, nil)
		e.Meta["name"] = []string{"field is missing"}
		tc.RenderError(e)
		return
	}
	name, err := data.AsString(n)
	if err != nil {
		tc.ErrLog(err).Error("'name' field isn't a string")
		e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, nil)
		e.Meta["name"] = []string{"value must be a string"}
		tc.RenderError(e)
		return
	}
	if err := core.ValidateSymbol(name); err != nil {
		tc.ErrLog(err).Error("'name' field has invalid format")
		e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, nil)
		e.Meta["name"] = []string{"inavlid format"}
		tc.RenderError(e)
		return
	}

	// TODO: support other parameters

	cc := &core.ContextConfig{
		Logger: tc.logger,
	}
	// TODO: Be careful of race conditions on these fields.
	cc.Flags.DroppedTupleLog.Set(tc.config.Logging.LogDroppedTuples)
	cc.Flags.DroppedTupleSummarization.Set(tc.config.Logging.SummarizeDroppedTuples)

	tp, err := core.NewDefaultTopology(core.NewContext(cc), name)
	if err != nil {
		tc.ErrLog(err).Error("Cannot create a new topology")
		tc.RenderErrorJSON(NewInternalServerError(err))
		return
	}
	tb, err := bql.NewTopologyBuilder(tp)
	if err != nil {
		tc.ErrLog(err).Error("Cannot create a new topology builder")
		tc.RenderError(jasco.NewInternalServerError(err))
		return
	}
	tb.UDSStorage = tc.udsStorage

	if err := tc.topologies.Register(name, tb); err != nil {
		if err := tp.Stop(); err != nil {
			tc.ErrLog(err).Error("Cannot stop the created topology")
		}

		if os.IsExist(err) {
			tc.Log().Error("the name is already registered")
			e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
				http.StatusBadRequest, nil)
			e.Meta["name"] = []string{"already taken"}
			tc.RenderError(e)
			return
		}
		tc.ErrLog(err).WithField("err", err).Error("Cannot register the topology")
		tc.Render(jasco.NewInternalServerError(err))
		return
	}

	// TODO: return 201
	tc.Render(map[string]interface{}{
		"topology": response.NewTopology(tb.Topology()),
	})
}

// Index returned a list of registered topologies.
func (tc *topologies) Index(rw web.ResponseWriter, req *web.Request) {
	ts, err := tc.topologies.List()
	if err != nil {
		tc.ErrLog(err).Error("Cannot list registered topologies")
		tc.RenderError(jasco.NewInternalServerError(err))
		return
	}

	res := []*response.Topology{}
	for _, tb := range ts {
		res = append(res, response.NewTopology(tb.Topology()))
	}
	tc.Render(map[string]interface{}{
		"topologies": res,
	})
}

// Show returns the information of topology
func (tc *topologies) Show(rw web.ResponseWriter, req *web.Request) {
	tb := tc.fetchTopology()
	if tb == nil {
		return
	}
	tc.Render(map[string]interface{}{
		"topology": response.NewTopology(tb.Topology()),
	})
}

// TODO: provide Update action (change state of the topology, etc.)

func (tc *topologies) Destroy(rw web.ResponseWriter, req *web.Request) {
	tb, err := tc.topologies.Unregister(tc.topologyName)
	isNotExist := core.IsNotExist(err)
	if err != nil && !isNotExist {
		tc.ErrLog(err).Error("Cannot unregister the topology")
		tc.RenderError(jasco.NewInternalServerError(err))
		return
	}
	stopped := true
	if tb != nil {
		if err := tb.Topology().Stop(); err != nil {
			stopped = false
			tc.ErrLog(err).Error("Cannot stop the topology")
		}
	}

	if stopped {
		// TODO: return 204 when the topology didn't exist.
		tc.Render(map[string]interface{}{})
	} else {
		tc.Render(map[string]interface{}{
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

	var js map[string]interface{}
	if apiErr := tc.ParseBody(&js); apiErr != nil {
		tc.ErrLog(apiErr.Err).Error("Cannot parse the request json")
		tc.RenderError(apiErr)
		return
	}

	form, err := data.NewMap(js)
	if err != nil {
		tc.ErrLog(err).WithField("body", js).
			Error("The request json may contain invalid value")
		tc.RenderError(jasco.NewError(formValidationErrorCode, "The request json may contain invalid values.",
			http.StatusBadRequest, err))
		return
	}

	var stmts []interface{}
	if ss, err := tc.parseQueries(form); err != nil {
		tc.RenderError(err)
		return
	} else if len(ss) == 0 {
		// TODO: support the new format
		tc.Render(map[string]interface{}{
			"topology_name": tc.topologyName,
			"status":        "running",
			"queries":       []interface{}{},
		})
		return
	} else {
		stmts = ss
	}

	if len(stmts) == 1 {
		stmtStr := fmt.Sprint(stmts[0])
		if stmt, ok := stmts[0].(parser.SelectStmt); ok {
			tc.handleSelectStmt(rw, stmt, stmtStr)
			return
		} else if stmt, ok := stmts[0].(parser.SelectUnionStmt); ok {
			tc.handleSelectUnionStmt(rw, stmt, stmtStr)
			return
		} else if stmt, ok := stmts[0].(parser.EvalStmt); ok {
			tc.handleEvalStmt(rw, stmt, stmtStr)
			return
		}
	}

	// TODO: handle this atomically
	for _, stmt := range stmts {
		// TODO: change the return value of AddStmt to support the new response format.
		_, err := tb.AddStmt(stmt)
		if err != nil {
			tc.ErrLog(err).Error("Cannot process a statement")
			e := jasco.NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
			e.Meta["error"] = err.Error()
			e.Meta["statement"] = fmt.Sprint(stmt)
			tc.RenderError(e)
			return
		}
	}

	// TODO: support the new format
	tc.Render(map[string]interface{}{
		"topology_name": tc.topologyName,
		"status":        "running",
		"queries":       stmts,
	})
}

func (tc *topologies) parseQueries(form data.Map) ([]interface{}, *jasco.Error) {
	// TODO: use mapstructure when parameters get too many
	var queries string
	if v, ok := form["queries"]; !ok {
		errMsg := "The request json doesn't have 'queries' field"
		tc.Log().Error(errMsg)
		e := jasco.NewError(formValidationErrorCode, "'queries' field is missing",
			http.StatusBadRequest, nil)
		return nil, e
	} else if f, err := data.AsString(v); err != nil {
		errMsg := "'queries' must be a string"
		tc.ErrLog(err).Error(errMsg)
		e := jasco.NewError(formValidationErrorCode, "'queries' field must be a string",
			http.StatusBadRequest, err)
		return nil, e
	} else {
		queries = f
	}

	bp := parser.New()
	stmts := []interface{}{}
	dataReturningStmtIndex := -1
	for queries != "" {
		stmt, rest, err := bp.ParseStmt(queries)
		if err != nil {
			tc.Log().WithField("parse_errors", err.Error()).
				WithField("statement", queries).Error("Cannot parse a statement")
			e := jasco.NewError(bqlStmtParseErrorCode, "Cannot parse a BQL statement", http.StatusBadRequest, err)
			e.Meta["parse_errors"] = strings.Split(err.Error(), "\n") // FIXME: too ad hoc
			e.Meta["statement"] = queries
			return nil, e
		}
		if _, ok := stmt.(parser.SelectStmt); ok {
			dataReturningStmtIndex = len(stmts)
		} else if _, ok := stmt.(parser.SelectUnionStmt); ok {
			dataReturningStmtIndex = len(stmts)
		} else if _, ok := stmt.(parser.EvalStmt); ok {
			dataReturningStmtIndex = len(stmts)
		}

		stmts = append(stmts, stmt)
		queries = rest
	}

	if dataReturningStmtIndex >= 0 {
		if len(stmts) != 1 {
			errMsg := "A SELECT or EVAL statement cannot be issued with other statements"
			tc.Log().Error(errMsg)
			e := jasco.NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, nil)
			e.Meta["error"] = "a SELECT or EVAL statement cannot be issued with other statements"
			e.Meta["statement"] = fmt.Sprint(stmts[dataReturningStmtIndex])
			return nil, e
		}
	}
	return stmts, nil
}

func (tc *topologies) handleSelectStmt(rw web.ResponseWriter, stmt parser.SelectStmt, stmtStr string) {
	tmpStmt := parser.SelectUnionStmt{[]parser.SelectStmt{stmt}}
	tc.handleSelectUnionStmt(rw, tmpStmt, stmtStr)
}

func (tc *topologies) handleSelectUnionStmt(rw web.ResponseWriter, stmt parser.SelectUnionStmt, stmtStr string) {
	tb := tc.fetchTopology()
	if tb == nil { // just in case
		return
	}

	sn, ch, err := tb.AddSelectUnionStmt(&stmt)
	if err != nil {
		tc.ErrLog(err).Error("Cannot process a statement")
		e := jasco.NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
		e.Meta["error"] = err.Error()
		e.Meta["statement"] = stmtStr
		tc.RenderError(e)
		return
	}
	defer func() {
		go func() {
			// vacuum all tuples to avoid blocking the sink.
			for _ = range ch {
			}
		}()
		if err := sn.Stop(); err != nil {
			tc.ErrLog(err).WithFields(logrus.Fields{
				"node_type": core.NTSink,
				"node_name": sn.Name(),
			}).Error("Cannot stop the temporary sink")
		}
	}()

	conn, bufrw, err := rw.Hijack()
	if err != nil {
		tc.ErrLog(err).Error("Cannot hijack a connection")
		tc.RenderError(jasco.NewInternalServerError(err))
		return
	}

	var (
		writeErr error
		readErr  error
	)
	mw := multipart.NewWriter(bufrw)
	defer func() {
		if writeErr != nil {
			tc.ErrLog(writeErr).Info("Cannot write contents to the hijacked connection")
		}

		if err := mw.Close(); err != nil {
			if writeErr == nil && readErr == nil { // log it only when the write err hasn't happend
				tc.ErrLog(err).Info("Cannot finish the multipart response")
			}
		}
		bufrw.Flush()
		conn.Close()

		tc.Log().WithField("statement", stmtStr).Info("Finish streaming SELECT responses")
	}()

	res := []string{
		"HTTP/1.1 200 OK",
		fmt.Sprintf(`Content-Type: multipart/mixed; boundary="%v"`, mw.Boundary()),
		"\r\n",
	}
	if _, err := bufrw.WriteString(strings.Join(res, "\r\n")); err != nil {
		tc.ErrLog(err).Error("Cannot write a header to the hijacked connection")
		return
	}
	bufrw.Flush()

	tc.Log().WithField("statement", stmtStr).Info("Start streaming SELECT responses")

	// All error reporting logs after this is info level because they might be
	// caused by the client closing the connection.
	header := textproto.MIMEHeader{}
	header.Add("Content-Type", "application/json")

	readPoll := time.After(1 * time.Minute)
	sent := false
	dummyReadBuf := make([]byte, 1024)
	for {
		var t *core.Tuple
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			t = v
			sent = true
		case <-readPoll:
			if sent {
				sent = false
				readPoll = time.After(1 * time.Minute)
				continue
			}

			// Assuming there's no more data to be read. Because no tuple was
			// written for past 1 minute, blocking read for 1ms here isn't a
			// big deal.
			// TODO: is there any better way to detect disconnection?
			// TODO: If general errors are checked before checking the deadline,
			//       this code doesn't have to add 1ms.
			if err := conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond)); err != nil {
				tc.ErrLog(err).Error("Cannot check the status of connection due to the failure of conn.SetReadDeadline. Stopping streaming.")
				// This isn't handled as a read error because some operating
				// systems don't support SetReadDeadline.
				return
			}
			if _, err := bufrw.Read(dummyReadBuf); err != nil {
				type timeout interface {
					Timeout() bool
				}
				if e, ok := err.(timeout); !ok || !e.Timeout() {
					// Something happened on this connection.
					readErr = err
					tc.ErrLog(err).Error("The connection may be closed from the client side")
					return
				}
			}
			readPoll = time.After(1 * time.Minute)
			continue
		}

		js := t.Data.String()
		// TODO: don't forget to convert \n to \r\n when returning
		// pretty-printed JSON objects.
		header.Set("Content-Length", fmt.Sprint(len(js)))

		w, err := mw.CreatePart(header)
		if err != nil {
			writeErr = err
			return
		}
		if _, err := io.WriteString(w, js); err != nil {
			writeErr = err
			return
		}
		if err := bufrw.Flush(); err != nil {
			writeErr = err
			return
		}
	}
}

func (tc *topologies) handleEvalStmt(rw web.ResponseWriter, stmt parser.EvalStmt, stmtStr string) {
	tb := tc.fetchTopology()
	if tb == nil { // just in case
		return
	}

	result, err := tb.RunEvalStmt(&stmt)
	if err != nil {
		tc.ErrLog(err).Error("Cannot process a statement")
		e := jasco.NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
		e.Meta["error"] = err.Error()
		e.Meta["statement"] = stmtStr
		tc.RenderError(e)
		return
	}

	// return value with JSON wrapper so it can be parsed on the client side
	tc.Render(map[string]interface{}{
		"result": result,
	})
}

// WebSocketQueries handles requests using WebSocket. A single WebSocket
// connection can concurrently issue multiple requests including requests
// containing a SELECT statement.
//
// All WebSocket request need to have following fields:
//
//	* rid
//	* payload
//
// "rid" field is used at the client side to identify to which request a response
// corresponds. All responses have "rid" field having the same value as the one
// in its corresponding request. rid can be any integer which is greather than
// 0 as long as the client can distinguish responses. rid 0 is used by the
// server when returning an error which happened before the actual rid can be
// obtained.
//
// "payload" field contains a request data same as the one sent to the regular
// HTTP request. Therefore, WebSocket requests have the same limitations such as
// "A SELECT statement cannot be issued with other statements including another
// SELECT statement". However, as it's mentioned earlier, a single WebSocket
// connection can concurrently send multiple requests which have a single
// SELECT statement.
//
// Example:
//
//	{
//		"rid": 1,
//		"payload": {
//			"queries": "SELECT RSTREAM * FROM my_stream [RANGE 1 TUPLES];"
//		}
//	}
//
// All WebSocket responses have following fields:
//
//	* rid
//	* type
//	* payload
//
// "rid" field contains the ID of the request to which the response corresponds.
//
// "type" field contains the type of the response:
//
//	* "result"
//	* "error"
//	* "sos"
//	* "ping"
//	* "eos"
//
// When the type is "result", "payload" field contains the result obtained by
// executing the query. The form of response depends on the type of a statement
// and some statements returns multiple responses. When the type is "error",
// "payload" has an error information which is same as the error response
// that Queries action returns. "sos", start of stream, type is used by SELECT
// statements to notify the client that a SELECT statement finishes setting up
// all necessary nodes in the topology. Its payload is always null. "ping"
// type is used by SELECT statements to validate connection. Its "payload" is
// always null. SELECT statements send "ping" responses on a regular basis.
// "eos", end of stream, responses are sent when SELECT statements has sent all
// tuples. "payload" of "eos" is always null. "eos" isn't sent when an error
// occurred.
func (tc *topologies) WebSocketQueries(rw web.ResponseWriter, req *web.Request) {
	// TODO: add a document describing which BQL statement returns which result.
	if !strings.EqualFold(req.Header.Get("Upgrade"), "WebSocket") {
		err := fmt.Errorf("the request isn't a WebSocket request")
		tc.Log().Error(err)
		tc.RenderError(jasco.NewError(nonWebSocketRequestErrorCode, "This action only accepts WebSocket connections",
			http.StatusBadRequest, err))
		return
	}

	tb := tc.fetchTopology()
	if tb == nil {
		return
	}

	tc.Log().Info("Begin WebSocket connection")
	defer tc.Log().Info("End WebSocket connection")

	websocket.Handler(func(conn *websocket.Conn) {
		for tc.processWebSocketMessage(conn, tb) {
		}
	}).ServeHTTP(rw, req.Request)
}

// processWebSocketMessage processes a request from the client. It returns true
// if the caller can call this method again, in other words, the connection is
// still alive.
func (tc *topologies) processWebSocketMessage(conn *websocket.Conn, tb *bql.TopologyBuilder) bool {
	w := &webSocketTopologyQueryHandler{
		tc:   tc,
		conn: conn,
	}

	var js map[string]interface{}
	if err := w.receive(&js); err != nil {
		e := jasco.NewError(bqlStmtParseErrorCode,
			"Cannot read or parse a JSON body received from the WebSocket connection",
			http.StatusBadRequest, err)
		tc.ErrLog(err).Error(e.Message)
		// When the error message cannot be sent back to the client, the connection
		// might be lost. So, tell the caller about it
		return w.sendErr(e)
	}

	form, err := data.NewMap(js)
	if err != nil {
		tc.ErrLog(err).WithField("body", js).Error("The request json may contain invalid value")
		return w.sendErr(jasco.NewError(formValidationErrorCode, "The request json may contain invalid values.",
			http.StatusBadRequest, err))
	}

	// TODO: use mapstructure or json schema for validation
	// TODO: return as many errors at once as possible
	var payload data.Map
	if v, ok := form["rid"]; !ok {
		tc.Log().Error("The required 'rid' field is missing")
		e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, err)
		e.Meta["rid"] = []string{"field is missing"}
		return w.sendErr(e)

	} else if r, err := data.ToInt(v); err != nil {
		tc.ErrLog(err).Error("Cannot convert 'rid' to an integer")
		e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, err)
		e.Meta["rid"] = []string{"value must be an integer"}
		return w.sendErr(e)

	} else {
		w.rid = r
	}

	// TODO: access logging

	// rid should be logged from this point. So, following logging should be
	// done by w.Log/w.ErrLog.

	if v, ok := form["payload"]; !ok {
		w.Log().Error("The required 'payload' field is missing")
		e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, err)
		e.Meta["payload"] = []string{"field is missing"}
		return w.sendErr(e)

	} else if p, err := data.AsMap(v); err != nil {
		w.ErrLog(err).Error("Cannot convert 'payload' to an integer")
		e := jasco.NewError(formValidationErrorCode, "The request body is invalid.",
			http.StatusBadRequest, err)
		e.Meta["payload"] = []string{"value must be an object"}
		return w.sendErr(e)

	} else {
		payload = p
	}

	// TODO: merge the following implementation with Queries.
	var stmts []interface{}
	if ss, err := tc.parseQueries(payload); err != nil { // TODO: logs from this method should have wsreqid, too
		return w.sendErr(err)
	} else if len(ss) == 0 {
		if err := w.send("result", map[string]interface{}{}); err != nil {
			w.ErrLog(err).Error("Cannot send a response to the WebSocket client")
			return false
		}
		return true
	} else {
		stmts = ss
	}

	// Although these requests may fail asynchronously, the connect is probably
	// still alive and next processWebSocketMessage can detect disconnection.
	// So, the following code block always returns true.
	go func() {
		if len(stmts) == 1 {
			stmtStr := fmt.Sprint(stmts[0])
			if stmt, ok := stmts[0].(parser.SelectStmt); ok {
				w.handleSelectStmtWebSocket(conn, stmt, stmtStr)
				return
			} else if stmt, ok := stmts[0].(parser.SelectUnionStmt); ok {
				w.handleSelectUnionStmtWebSocket(conn, stmt, stmtStr)
				return
			} else if stmt, ok := stmts[0].(parser.EvalStmt); ok {
				w.handleEvalStmtWebSocket(conn, stmt, stmtStr)
				return
			}
		}

		// TODO: handle this atomically
		for _, stmt := range stmts {
			// TODO: change the return value of AddStmt to support the new response format.
			_, err = tb.AddStmt(stmt)
			if err != nil {
				w.ErrLog(err).Error("Cannot process a statement")
				e := jasco.NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
				e.Meta["error"] = err.Error()
				e.Meta["statement"] = fmt.Sprint(stmt)
				w.sendErr(e)
				return
			}
		}

		// TODO: define a proper response format
		if err := w.send("result", map[string]interface{}{}); err != nil {
			w.ErrLog(err).Error("Cannot send a response to the WebSocket client")
		}
	}()
	return true
}

type webSocketTopologyQueryHandler struct {
	tc   *topologies
	conn *websocket.Conn
	rid  int64
}

func (w *webSocketTopologyQueryHandler) Log() *logrus.Entry {
	return w.tc.Log().WithField("wsreqid", w.rid)
}

func (w *webSocketTopologyQueryHandler) ErrLog(err error) *logrus.Entry {
	return w.tc.ErrLog(err).WithField("wsreqid", w.rid)
}

func (w *webSocketTopologyQueryHandler) receive(v interface{}) error {
	return websocket.JSON.Receive(w.conn, v)
}

func (w *webSocketTopologyQueryHandler) send(msgType string, v interface{}) error {
	return websocket.JSON.Send(w.conn, map[string]interface{}{
		"rid":     w.rid,
		"type":    msgType,
		"payload": v,
	})
}

// sendErr sends an error message to the client. It returns true when the
// response could be sent.
func (w *webSocketTopologyQueryHandler) sendErr(e *jasco.Error) bool {
	if err := w.send("error", e); err != nil {
		// TODO: this error message should have the caller's line number.
		w.ErrLog(err).Error("Cannot send an error response to the WebSocket connection")
		return false
	}
	return true
}

func (w *webSocketTopologyQueryHandler) handleSelectStmtWebSocket(conn *websocket.Conn, stmt parser.SelectStmt, stmtStr string) {
	tmpStmt := parser.SelectUnionStmt{[]parser.SelectStmt{stmt}}
	w.handleSelectUnionStmtWebSocket(conn, tmpStmt, stmtStr)
}

func (w *webSocketTopologyQueryHandler) handleSelectUnionStmtWebSocket(conn *websocket.Conn, stmt parser.SelectUnionStmt, stmtStr string) {
	// TODO: merge this function with handleSelectUnionStmt if possible
	tb := w.tc.fetchTopology()
	if tb == nil { // just in case
		return
	}

	sn, ch, err := tb.AddSelectUnionStmt(&stmt)
	if err != nil {
		w.ErrLog(err).Error("Cannot process a statement")
		e := jasco.NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
		e.Meta["error"] = err.Error()
		e.Meta["statement"] = stmtStr
		w.sendErr(e)
		return
	}
	defer func() {
		w.Log().WithField("statement", stmtStr).Info("Finish streaming SELECT responses")

		go func() {
			// vacuum all tuples to avoid blocking the sink.
			for _ = range ch {
			}
		}()
		if err := sn.Stop(); err != nil {
			w.ErrLog(err).WithFields(logrus.Fields{
				"node_type": core.NTSink,
				"node_name": sn.Name(),
			}).Error("Cannot stop the temporary sink")
		}
	}()

	w.Log().WithField("statement", stmtStr).Info("Start streaming SELECT responses")

	if err := w.send("sos", nil); err != nil {
		w.ErrLog(err).Error("Cannot send an sos to the WebSocket client")
		return
	}

	ping := time.After(1 * time.Minute)
	sent := false
	for {
		var t *core.Tuple
		select {
		case v, ok := <-ch:
			if !ok {
				if err := w.send("eos", nil); err != nil {
					w.ErrLog(err).Error("Cannot send an EOS message to the WebSocket client")
				}
				return
			}
			t = v
			sent = true
		case <-ping:
			if sent {
				sent = false
				ping = time.After(1 * time.Minute)
				continue
			}

			if err := w.send("ping", nil); err != nil {
				w.ErrLog(err).Error("The connection may be closed from the client side")
				return
			}
			ping = time.After(1 * time.Minute)
			continue
		}

		if err := w.send("result", t.Data); err != nil {
			w.ErrLog(err).Error("Cannot send an error response to the WebSocket client")
			return
		}
	}
}

func (w *webSocketTopologyQueryHandler) handleEvalStmtWebSocket(conn *websocket.Conn, stmt parser.EvalStmt, stmtStr string) {
	tb := w.tc.fetchTopology()
	if tb == nil { // just in case
		return
	}

	result, err := tb.RunEvalStmt(&stmt)
	if err != nil {
		w.ErrLog(err).Error("Cannot process a statement")
		e := jasco.NewError(bqlStmtProcessingErrorCode, "Cannot process a statement", http.StatusBadRequest, err)
		e.Meta["error"] = err.Error()
		e.Meta["statement"] = stmtStr
		w.sendErr(e)
		return
	}

	if err := w.send("result", map[string]interface{}{
		"result": result,
	}); err != nil {
		w.ErrLog(err).Error("Cannot send data to the WebSocket client")
		return
	}
}

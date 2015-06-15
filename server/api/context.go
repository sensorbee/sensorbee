package api

import (
	"encoding/json"
	"github.com/gocraft/web"
	"net/http"
	"sync/atomic"
)

type BaseContext struct{}

type Context struct {
	*BaseContext

	body      []byte
	bodyError error

	requestID uint64

	response   web.ResponseWriter
	request    *web.Request
	HTTPStatus int
}

var requestIDCounter uint64 = 0

func (c *Context) setUpContext(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.requestID = atomic.AddUint64(&requestIDCounter, 1)
	c.response = rw
	c.request = req

	next(rw, req)
}

func SetUpRouterWithCustomMiddleware(prefix string, parent *web.Router, middleware func(c *Context, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc), callback func(string, *web.Router)) *web.Router {
	if parent == nil {
		parent = web.New(BaseContext{})
	}

	root := parent.Subrouter(Context{}, "/")
	if middleware != nil {
		root.Middleware(middleware)
	}
	root.Middleware((*Context).setUpContext)
	callback(prefix, root)
	return parent
}

func (c *Context) renderJSON(status int, v interface{}) {
	c.HTTPStatus = status

	data, err := json.Marshal(v)
	if err != nil {
		// TODO logging
		c.response.Header().Set("Content-Type", "application/json")
		c.response.WriteHeader(http.StatusInternalServerError)
		c.response.Write([]byte(`
{
  "errors": [
    {
      "code": "E0006",
      "message": "internal server error"
    }
  ]
}
`))
		return
	}

	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(status)
	_, err = c.response.Write(data)
	if err != nil {
		// logging
	}
}

func (c *Context) RenderJSON(v interface{}) {
	c.renderJSON(http.StatusOK, v)
}

func (c *Context) RenderErrorJSON(e *Error) {
	e.SetRequestID(c.requestID)

	c.renderJSON(e.Status, &struct {
		Errors []interface{} `json:"errors"`
	}{
		Errors: []interface{}{e},
	})
}

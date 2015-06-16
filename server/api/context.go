package api

import (
	"bytes"
	"encoding/json"
	"github.com/gocraft/web"
	"io"
	"io/ioutil"
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

// Body returns a slice containing whole request body.
// It caches the result so controllers can call this
// method as many time as it wants.
//
// When the request body is empty (i.e. Read(req.Body)
// returns io.EOF), this method returns and caches
// an empty body slice and a nil error.
func (c *Context) Body() ([]byte, error) {
	if c.body != nil || c.bodyError != nil {
		return c.body, c.bodyError
	}

	body, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		if err == io.EOF {
			// when body is empty, this method caches
			// an empty slice and a nil error.
			body = []byte{}
			err = nil
		}
	}

	// Close and replace with new ReadCloser for parsing
	// mime/multipart request body by Request.FormFile method.
	c.request.Body.Close()
	c.request.Body = ioutil.NopCloser(bytes.NewReader(body))
	c.body = body
	c.bodyError = err
	return body, err
}

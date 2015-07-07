package server

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

// Context is a context object for gocraft/web.
type Context struct {
	body      []byte
	bodyError error

	requestID uint64

	response   web.ResponseWriter
	request    *web.Request
	HTTPStatus int

	logger     *logrus.Logger
	topologies TopologyRegistry
}

// SetLogger sets the logger to the context. Must be set before any action
// is invoked.
func (c *Context) SetLogger(l *logrus.Logger) {
	c.logger = l
}

var requestIDCounter uint64 = 0

func (c *Context) setUpContext(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.requestID = atomic.AddUint64(&requestIDCounter, 1)
	c.response = rw
	c.request = req

	// TODO: request logging
	next(rw, req)
}

// SetTopologyRegistry sets the registry of topologies to this context. This
// method must be called in the middleware of Context.
func (c *Context) SetTopologyRegistry(r TopologyRegistry) {
	c.topologies = r
}

func (c *Context) NotFoundHandler(rw web.ResponseWriter, req *web.Request) {
	c.logger.WithFields(logrus.Fields{
		"method": req.Method,
		"uri":    req.URL.RequestURI(),
		"status": http.StatusNotFound,
	}).Error("The request URL not found")

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(`
{
  "errors": [
    {
      "code": "E0001",
      "message": "The request URL was not found."
    }
  ]
}`))
}

// Log returns the logger having meta information.
func (c *Context) Log() *logrus.Entry {
	return c.logger.WithFields(logrus.Fields{
		"reqid": c.requestID,
	})
}

// ErrLog returns the logger with error information.
func (c *Context) ErrLog(err error) *logrus.Entry {
	return c.Log().WithField("err", err)
}

func (c *Context) extractOptionStringFromPath(key string, target *string) error {
	s, ok := c.request.PathParams[key]
	if !ok {
		return nil
	}

	*target = s
	return nil
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
	c.renderJSON(e.Status, map[string]interface{}{
		"error": e,
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

// ContextGlobalVariables has fields which are shared through all contexts
// allocated for each request.
type ContextGlobalVariables struct {
	// Logger is used to write log messages.
	Logger *logrus.Logger

	// Topologies is a registry which manages topologies to support multi
	// tenancy.
	Topologies TopologyRegistry
}

// SetUpRouter creates a root router of the API server.
func SetUpRouter(prefix string, gvars ContextGlobalVariables) *web.Router {
	root := web.NewWithPrefix(Context{}, prefix)
	root.NotFound((*Context).NotFoundHandler)
	root.Middleware(func(c *Context, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		c.logger = gvars.Logger
		c.topologies = gvars.Topologies
		next(rw, req)
	})
	root.Middleware((*Context).setUpContext)
	return root
}

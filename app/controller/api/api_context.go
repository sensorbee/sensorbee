package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"sync/atomic"
)

type APIContext struct {
	*Context
}

func SetUpAPIRouter(prefix string, router *web.Router) {
	SetUpAPIRouterWithCustomRoute(prefix, router, nil)
}

func SetUpAPIRouterWithCustomRoute(prefix string, router *web.Router, route func(prefix string, r *web.Router)) {
	root := router.Subrouter(APIContext{}, "/api") // TODO need to set version like hawk??

	SetUpBQLRouter(prefix, root)

	if route != nil {
		route(prefix, root)
	}
}

/********** Context TODO will be separated **********/

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

/********** Error TODO separate to error.go **********/

const (
	MaxErrorCode = 9999
)

var (
	InternalErrorCode = GenerateErrorCode("context", 6)
)

var packageErrorPrefixes = map[string]rune{}

// Error has error information occurred in Hawk.
// It has error messages for clients(users).
type Error struct {
	originalCode *ErrorCode `json:"-"`

	// Code has a error code of this error. It always have
	// a single alphabet prefix and 4-digit error number.
	Code string `json:"code"`

	// Message contains a user-friendly error message for users
	// of WebUI or clients of API. It MUST NOT have internal
	// error information nor debug information which are
	// completely useless for them.
	Message string `json:"message"`

	// RequestID has an ID of the request which caused this error.
	// This value is convenient for users or clients of Hawk
	// to report their problems to the development team.
	// The team can look up incidents easily by greping logs
	// with the ID.
	//
	// The type of RequestID is string because it can be
	// very large and JavaScript might not be able to handle
	// it as an integer.
	RequestID string `json:"request_id"`

	// Status has an appropriate HTTP status code for this error.
	Status int `json:"-"`

	// Err contains internal error information returned from
	// a method or a function which caused the error.
	Err error `json:"-"`

	// Meta contains arbitrary information of the error.
	// What the information means depends on each error.
	// For instance, a validation error might contain
	// error information of each field.
	Meta map[string]interface{} `json:"meta"`
}

// ErrorCode has an error code of an error defined in each package.
type ErrorCode struct {
	// Package is the name of the module which caused this error.
	Package string `json:"package"`

	// Code has an error code of this error.
	Code int `json:"code"`
}

// String returns a string representation of the error code.
// The string contains one letter package code and 4-digit
// error code (e.g. "E0105").
func (e *ErrorCode) String() string {
	p, ok := packageErrorPrefixes[e.Package]
	if !ok {
		p = 'X'
	}
	return fmt.Sprintf("%c%04d", p, e.Code)
}

func (e *Error) SetRequestID(id uint64) {
	e.RequestID = fmt.Sprint(id)
}

func NewError(code *ErrorCode, msg string, status int, err error) *Error {
	if err == nil {
		err = errors.New(msg)
	}

	return &Error{
		originalCode: code,
		Code:         code.String(),
		Message:      msg,
		Status:       status,
		Err:          err,
		Meta:         make(map[string]interface{}),
	}
}

// GenerateErrorCode creates a new error code.
func GenerateErrorCode(pkg string, code int) *ErrorCode {
	if code <= 0 || MaxErrorCode < code {
		// Since this function is only executed on initialization,
		// it can safely panic.
		panic(fmt.Errorf("An error code must be in [1, %v]: %v", MaxErrorCode, code))
	}
	return &ErrorCode{
		Package: pkg,
		Code:    code,
	}
}

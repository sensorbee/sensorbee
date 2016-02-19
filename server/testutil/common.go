package testutil

import (
	"bytes"
	"github.com/mattn/go-scan"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"gopkg.in/sensorbee/sensorbee.v0/server"
	"gopkg.in/sensorbee/sensorbee.v0/server/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var (
	// TestAPIWithRealHTTPServer is a flag to switch the use of httptest.server.
	// If this flag is true, all tests will be executed with a real HTTP server.
	// If it is false, a server will not be created and tests are done by directly
	// calling ServeHTTP of a handler.
	TestAPIWithRealHTTPServer = false
)

// Server is a temporary HTTP server for test.
type Server struct {
	// server has the information of the server
	server struct {
		realServer *httptest.Server
		router     http.Handler
		url        string
	}
}

// Close closes the server.
func (s *Server) Close() {
	if s.server.realServer != nil {
		s.server.realServer.Close()
	}
}

// URL returns the URL of the server.
func (s *Server) URL() string {
	return s.server.url
}

// HTTPClient returns the HTTP client to send requests to the server.
func (s *Server) HTTPClient() *http.Client {
	if s.server.realServer != nil {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: &testTransport{s},
	}
}

type testTransport struct {
	s *Server
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	}

	rc := httptest.NewRecorder()
	t.s.server.router.ServeHTTP(rc, req)
	return &http.Response{
		Status:        http.StatusText(rc.Code),
		StatusCode:    rc.Code,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        rc.HeaderMap,
		Body:          ioutil.NopCloser(rc.Body),
		ContentLength: int64(rc.Body.Len()),
	}, nil
}

// NewServer returns a temporary running server.
func NewServer() *Server {
	s := &Server{}

	c, err := config.New(data.Map{})
	if err != nil {
		panic(err)
	}
	gvars, err := server.SetUpContextGlobalVariables(c)
	if err != nil {
		panic(err)
	}
	root, err := server.SetUpContextAndRouter("/", gvars)
	if err != nil {
		panic(err)
	}
	server.SetUpAPIRouter("/", root, nil)

	if TestAPIWithRealHTTPServer {
		s.server.realServer = httptest.NewServer(root)
		s.server.url = s.server.realServer.URL
	} else {
		s.server.router = root
		s.server.url = "http://172.0.0.1:0602"
	}
	return s
}

// JScan traverses a json object tree.
func JScan(js interface{}, path string) interface{} {
	var v interface{}
	if err := scan.ScanTree(js, path, &v); err != nil {
		return nil
	}
	return v
}

package testutil

import (
	"bytes"
	"encoding/json"
	"github.com/mattn/go-scan"
	"io"
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

// Create creates a new client of the server
func (s *Server) Client() *Client {
	return &Client{
		s: s,
	}
}

// Client is a client of the temporary server.
type Client struct {
	uri string
	s   *Server
}

// Do sends a request to the server.
func (c *Client) Do(method string, path string, body interface{}) (*http.Response, map[string]interface{}, error) {
	req, err := c.CreateRequest(method, path, body)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	js := map[string]interface{}{}
	if err = json.Unmarshal(data, &js); err != nil {
		return nil, nil, err
	}
	return res, js, err
}

// CreateRequest creates a new request for the server.
func (c *Client) CreateRequest(method string, path string, bodyJSON interface{}) (*http.Request, error) {
	var body io.Reader
	if bodyJSON == nil {
		body = nil
	} else {
		bd, err := json.Marshal(bodyJSON)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bd)
	}

	req, err := http.NewRequest(method, c.s.server.url+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// DoRequest sends a custom http.Request to the server.
func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	if TestAPIWithRealHTTPServer {
		return http.DefaultClient.Do(req)
	}

	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	}

	rc := httptest.NewRecorder()
	c.s.server.router.ServeHTTP(rc, req)
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
func NewServer(root http.Handler) *Server {
	s := &Server{}
	if TestAPIWithRealHTTPServer {
		s.server.realServer = httptest.NewServer(root)
		s.server.url = s.server.realServer.URL
	} else {
		s.server.router = root
		s.server.url = "http://172.0.0.1:0602/api/v1"
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

package api

import (
	"bytes"
	"encoding/json"
	"github.com/gocraft/web"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"pfi/sensorbee/sensorbee/client/cmd"
)

var (
	// testAPIWithRealHTTPServer is a flag to switch the use of httptest.Server.
	// If this flag is true, all tests will be executed with a real HTTP server.
	// If it is false, a server will not be created and tests are done by directly
	// calling ServeHTTP of a handler.
	testAPIWithRealHTTPServer = false
)

type temporaryServer struct {
	// Server has the information of the server
	Server struct {
		realServer *httptest.Server
		router     http.Handler
		URL        string
	}
}

func (t *temporaryServer) close() {
	if t.Server.realServer != nil {
		t.Server.realServer.Close()
	}
}

func (t *temporaryServer) dummyClient() *temporaryClient {
	return &temporaryClient{
		s: t,
	}
}

type temporaryClient struct {
	uri string
	s   *temporaryServer
}

func (t *temporaryClient) do(reqType cmd.RequestType, path string, body interface{}) (*http.Response, map[string]interface{}, error) {
	req, err := cmd.CreateRequest(reqType, t.s.Server.URL+path, body)
	if err != nil {
		return nil, nil, err
	}

	res, err := t.doWrapper(req)
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

func (t *temporaryClient) doWrapper(req *http.Request) (*http.Response, error) {
	if testAPIWithRealHTTPServer {
		return http.DefaultClient.Do(req)
	}

	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	}

	rc := httptest.NewRecorder()
	t.s.Server.router.ServeHTTP(rc, req)
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

func createTestServer() *temporaryServer {
	return createTestServerWithCustomRoute(nil)
}

func createTestServerWithCustomRoute(route func(prefix string, r *web.Router)) *temporaryServer {
	s := &temporaryServer{}

	root := SetUpRouter("/", nil,
		func(c *Context, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
			next(rw, req)
		},
		func(prefix string, r *web.Router) {
			SetUpAPIRouter(prefix, r, route)
		})

	if testAPIWithRealHTTPServer {
		s.Server.realServer = httptest.NewServer(root)
		s.Server.URL = s.Server.realServer.URL
	} else {
		s.Server.router = root
		s.Server.URL = "http://172.0.0.1:0602/api/v1"
	}
	return s
}

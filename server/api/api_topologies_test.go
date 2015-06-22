package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/client/cmd"
	"testing"
)

func TestTopologiesWithEmptyTenants(t *testing.T) {
	s := createTestServer()
	defer s.close()
	cli := s.dummyClient()

	Convey("Given a root topologies path", t, func() {
		path := "/topologies"
		Convey("When get request to test server", func() {
			Convey("Then return topologies tenant name list", func() {
				res, js, err := cli.do(cmd.GetRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				s := js["topologies"]
				tenants, ok := s.([]interface{})
				So(ok, ShouldBeTrue)
				So(len(tenants), ShouldEqual, 0)
			})
		})
		Convey("When post request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.PostRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, 404)
				So(len(js), ShouldEqual, 1)
			})
		})
	})
	Convey("Given a dummy tenant path", t, func() {
		path := "/topologies/dummy"
		Convey("When get request with not initialized tenant", func() {
			Convey("Then return valid response", func() {
				res, js, err := cli.do(cmd.GetRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")

				st := js["status"]
				status, ok := st.(string)
				So(ok, ShouldBeTrue)
				So(status, ShouldEqual, "not initialized")
			})
		})
		Convey("When put request with not initialized tenant", func() {
			Convey("Then return valid response", func() {
				body := map[string]interface{}{
					"state": "stop",
				}
				res, js, err := cli.do(cmd.PutRequest, path, body)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")

				st := js["status"]
				status, ok := st.(string)
				So(ok, ShouldBeTrue)
				So(status, ShouldEqual, "not initialized or running topology")
			})
		})
	})
}

func TestTopologiesExecuteQueriesAndShowTenant(t *testing.T) {
	s := createTestServer()
	defer s.close()
	cli := s.dummyClient()

	Convey("Given a path for post queries", t, func() {
		path := "/topologies/dummy/queries"
		Convey("When request body has a query", func() {
			queries := "CREATE SOURCE dummy TYPE dummy;"
			body := map[string]interface{}{
				"queries": queries,
			}
			Convey("Then load queries", func() {
				res, js, err := cli.do(cmd.PostRequest, path, body)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")

				st := js["status"]
				status, ok := st.(string)
				So(ok, ShouldBeTrue)
				So(status, ShouldStartWith, "statement of type")
				Convey("Given a topologies root path", func() {
					path := "/topologies/"
					Convey("When get request, then return tenant name list", func() {
						res, js, err := cli.do(cmd.GetRequest, path, nil)
						So(err, ShouldBeNil)
						So(res.StatusCode, ShouldEqual, http.StatusOK)
						s := js["topologies"]
						tenants, ok := s.([]interface{})
						So(ok, ShouldBeTrue)
						So(len(tenants), ShouldEqual, 1)
						name, ok := tenants[0].(string)
						So(ok, ShouldBeTrue)
						So(name, ShouldEqual, "dummy")
					})
				})
				Convey("Given a tenant path", func() {
					path := "/topologies/dummy"
					Convey("When get request", func() {
						res, js, err := cli.do(cmd.GetRequest, path, nil)
						So(err, ShouldBeNil)
						So(res.StatusCode, ShouldEqual, http.StatusOK)
						st := js["status"]
						status, ok := st.(string)
						So(ok, ShouldBeTrue)
						So(status, ShouldEqual, "initialized")
					})
					Convey("When put request to update stop", func() {
						body := map[string]interface{}{
							"state": "stop",
						}
						res, js, err := cli.do(cmd.PutRequest, path, body)
						So(err, ShouldBeNil)
						So(res.StatusCode, ShouldEqual, http.StatusOK)
						st := js["status"]
						status, ok := st.(string)
						So(ok, ShouldBeTrue)
						So(status, ShouldEqual, "done")
						stt := js["state"]
						state, ok := stt.(string)
						So(ok, ShouldBeTrue)
						So(state, ShouldEqual, "stop")
					})
				})
			})
		})
		Convey("When request body has an empty query", func() {
			body := map[string]interface{}{
				"queries": "",
			}
			Convey("Then return error status", func() {
				res, js, err := cli.do(cmd.PostRequest, path, body)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")

				st := js["status"]
				status, ok := st.(string)
				So(ok, ShouldBeTrue)
				So(status, ShouldEqual, "not support to execute empty query")
			})
		})
		Convey("When request body has invalid interface", func() {
			body := "invalid body"
			Convey("Then return error status", func() {
				res, js, err := cli.do(cmd.PostRequest, path, body)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")

				q := js["query byte"]
				debugQuery, ok := q.(string)
				So(ok, ShouldBeTrue)
				So(debugQuery, ShouldEqual, `"invalid body"`)

				st := js["status"]
				status, ok := st.(string)
				So(ok, ShouldBeTrue)
				So(status, ShouldStartWith, "json: cannot unmarshal string")
			})
		})
		Convey("When request body is nil", func() {
			Convey("Then return error status", func() {
				res, js, err := cli.do(cmd.PostRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")

				q := js["query byte"]
				debugQuery, ok := q.(string)
				So(ok, ShouldBeTrue)
				So(debugQuery, ShouldEqual, "")

				st := js["status"]
				status, ok := st.(string)
				So(ok, ShouldBeTrue)
				So(status, ShouldStartWith, "unexpected")
			})
		})
	})
}

// This is testing that access path is valid, will be testing
// topology nodes access
func TestSourcesRequestWhenNotInitialized(t *testing.T) {
	s := createTestServer()
	defer s.close()
	cli := s.dummyClient()

	Convey("Given source path", t, func() {
		path := "/topologies/dummy/sources/1"
		Convey("When get request to test server", func() {
			Convey("Then return topologies tenant name list", func() {
				res, js, err := cli.do(cmd.GetRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
		Convey("When post request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.PutRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
		Convey("When detete request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.DeleteRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
	})
}

// This is testing that access path is valid, will be testing
// topology nodes access
func TestStreamRequestWhenNotInitialized(t *testing.T) {
	s := createTestServer()
	defer s.close()
	cli := s.dummyClient()

	Convey("Given stream path", t, func() {
		path := "/topologies/dummy/streams/1"
		Convey("When get request to test server", func() {
			Convey("Then return topologies tenant name list", func() {
				res, js, err := cli.do(cmd.GetRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
		Convey("When post request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.PutRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
		Convey("When detete request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.DeleteRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
	})
}

// This is testing that access path is valid, will be testing
// topology nodes access
func TestSinksRequestWhenNotInitialized(t *testing.T) {
	s := createTestServer()
	defer s.close()
	cli := s.dummyClient()

	Convey("Given sink path", t, func() {
		path := "/topologies/dummy/sinks/1"
		Convey("When get request to test server", func() {
			Convey("Then return topologies tenant name list", func() {
				res, js, err := cli.do(cmd.GetRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
		Convey("When post request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.PutRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
		Convey("When detete request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.DeleteRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node id"]
				nodeId, ok := id.(float64)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, 1)
			})
		})
	})
}

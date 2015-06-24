package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/client/cmd"
	"testing"
)

// This is testing that access path is valid, will be testing
// topology nodes access
func TestSinksRequestWhenNotInitialized(t *testing.T) {
	s := createTestServer()
	defer s.close()
	cli := s.dummyClient()

	Convey("Given sink path", t, func() {
		path := "/topologies/dummy/sinks/sink1"
		Convey("When get request to test server", func() {
			Convey("Then return topologies tenant name list", func() {
				res, js, err := cli.do(cmd.GetRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["topology name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node name"]
				nodeId, ok := id.(string)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, "sink1")
			})
		})
		Convey("When post request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.PutRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["topology name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node name"]
				nodeId, ok := id.(string)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, "sink1")
			})
		})
		Convey("When detete request to test server", func() {
			Convey("Then return 404 not found", func() {
				res, js, err := cli.do(cmd.DeleteRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				n := js["topology name"]
				name, ok := n.(string)
				So(ok, ShouldBeTrue)
				So(name, ShouldEqual, "dummy")
				id := js["node name"]
				nodeId, ok := id.(string)
				So(ok, ShouldBeTrue)
				So(nodeId, ShouldEqual, "sink1")
			})
		})
	})
}

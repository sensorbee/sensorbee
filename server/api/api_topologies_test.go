package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/client/cmd"
	"testing"
)

func TestTopologiesShowTenants(t *testing.T) {
	s := createTestServer()
	defer s.close()
	cli := s.dummyClient()

	Convey("Given a dummy server", t, func() {
		Convey("When request with root path", func() {
			path := "/topologies"
			Convey("Then return topologies tenant name list", func() {

				res, js, err := cli.do(cmd.GetRequest, path, nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				s := js["topologies"]
				tenants, ok := s.([]string)
				So(ok, ShouldBeFalse)
				So(len(tenants), ShouldEqual, 0)
			})
		})
	})
}

//TODO other test

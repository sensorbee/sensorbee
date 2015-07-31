package client

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/server/testutil"
	"testing"
)

func TestServerStatus(t *testing.T) {
	s := testutil.NewServer()
	defer s.Close()
	r := newTestRequester(s)

	Convey("Given an API server", t, func() {
		Convey("When getting runtime_status", func() {
			res, js, err := do(r, Get, "/runtime_status", nil)
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then the response should have the status", func() {
				So(js["num_goroutine"], ShouldBeGreaterThanOrEqualTo, 0)
				So(js["num_cgo_call"], ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

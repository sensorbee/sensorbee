package client

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"os"
	"os/user"
	"pfi/sensorbee/sensorbee/server/testutil"
	"runtime"
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
				So(jsonNumberToInt64(js["num_goroutine"]), ShouldBeGreaterThanOrEqualTo, 0)
				So(jsonNumberToInt64(js["num_cgo_call"]), ShouldBeGreaterThanOrEqualTo, 0)
				So(jsonNumberToInt64(js["gomaxprocs"]), ShouldEqual, runtime.GOMAXPROCS(0))
				So(js["goroot"], ShouldEqual, runtime.GOROOT())
				So(jsonNumberToInt64(js["num_cpu"]), ShouldEqual, runtime.NumCPU())
				So(js["goversion"], ShouldEqual, runtime.Version())
				So(jsonNumberToInt64(js["pid"]), ShouldEqual, os.Getpid())

				dir, err := os.Getwd()
				So(err, ShouldBeNil)
				So(js["working_directory"], ShouldEqual, dir)

				host, err := os.Hostname()
				So(err, ShouldBeNil)
				So(js["hostname"], ShouldEqual, host)

				user, err := user.Current()
				So(err, ShouldBeNil)
				So(js["user"], ShouldEqual, user.Username)
			})
		})
	})
}

func jsonNumberToInt64(n interface{}) int64 {
	ret, err := n.(json.Number).Int64()
	if err != nil {
		panic(err)
	}
	return ret
}

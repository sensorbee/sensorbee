package client

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/server/testutil"
	"testing"
)

var jscan = testutil.JScan

func TestEmptyTopologies(t *testing.T) {
	s := testutil.NewServer()
	defer s.Close()
	r := newTestRequester(s)

	Convey("Given an API server", t, func() {
		Convey("When creating a topology", func() {
			res, js, err := do(r, Post, "/topologies", map[string]interface{}{
				"name": "test_topology",
			})
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)
			Reset(func() {
				do(r, Delete, "/topologies/test_topology", nil)
			})

			Convey("Then the response should have the name", func() {
				So(jscan(js, "/topology/name"), ShouldEqual, "test_topology")
			})

			Convey("Then getting the topology should succeed", func() {
				res, js, err := do(r, Get, "/topologies/test_topology", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)
				So(jscan(js, "/topology/name"), ShouldEqual, "test_topology")
			})

			Convey("And creating another topology having the same name", func() {
				res, js, err := do(r, Post, "/topologies", map[string]interface{}{
					"name": "test_topology",
				})
				So(err, ShouldBeNil)

				Convey("Then it should fail", func() {
					So(res.Raw.StatusCode, ShouldEqual, http.StatusBadRequest)
					So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
				})
			})

			Convey("And getting a list of topologies", func() {
				res, js, err := do(r, Get, "/topologies", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

				Convey("Then it should have the new topology", func() {
					So(jscan(js, "/topologies[0]/name"), ShouldEqual, "test_topology")
				})
			})
		})

		Convey("When getting a list of topologies", func() {
			res, js, err := do(r, Get, "/topologies", nil)
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then it should have an empty response", func() {
				So(jscan(js, "/topologies"), ShouldBeEmpty)
			})
		})

		Convey("When getting a nonexistent topology", func() {
			res, js, err := do(r, Get, "/topologies/test_topology", nil)
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusNotFound)

				Convey("And the respons should have error information", func() {
					So(jscan(js, "/error"), ShouldNotBeNil)
				})
			})
		})

		Convey("When deleting a nonexistent topology", func() {
			res, _, err := do(r, Delete, "/topologies/test_topology", nil)
			So(err, ShouldBeNil)

			Convey("Then it shouldn't fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK) // TODO: This should be replaced with 204 later
			})
		})
	})
}

func TestTopologiesCreateInvalidValues(t *testing.T) {
	s := testutil.NewServer()
	defer s.Close()
	r := newTestRequester(s)

	Convey("Given an API server", t, func() {
		Convey("When posting a request missing required fields", func() {
			res, js, err := do(r, Post, "/topologies", map[string]interface{}{})
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then it should have meta information", func() {
				So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
			})
		})

		Convey("When posting a broken JSON request", func() {
			res, js, err := do(r, Post, "/topologies", json.RawMessage("{broken}"))
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(jscan(js, "/error"), ShouldNotBeNil)
			})
		})

		Convey("When posting an integer name", func() {
			res, js, err := do(r, Post, "/topologies", map[string]interface{}{
				"name": 1,
			})
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
			})
		})

		Convey("When posting an empty name", func() {
			res, js, err := do(r, Post, "/topologies", map[string]interface{}{
				"name": "",
			})
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
			})
		})
	})
}

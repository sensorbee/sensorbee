package server

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestEmptyTopologies(t *testing.T) {
	s := createTestServer()
	defer s.Close()
	cli := s.Client()

	Convey("Given an API server", t, func() {
		Convey("When creating a topology", func() {
			res, js, err := cli.Do("POST", "/topologies", map[string]interface{}{
				"name": "test_topology",
			})
			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)
			Reset(func() {
				cli.Do("DELETE", "/topologies/test_topology", nil)
			})

			Convey("Then the response should have the name", func() {
				So(jscan(js, "/topology/name"), ShouldEqual, "test_topology")
			})

			Convey("Then getting the topology should succeed", func() {
				res, js, err := cli.Do("GET", "/topologies/test_topology", nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)
				So(jscan(js, "/topology/name"), ShouldEqual, "test_topology")
			})

			Convey("And creating another topology having the same name", func() {
				res, js, err := cli.Do("POST", "/topologies", map[string]interface{}{
					"name": "test_topology",
				})
				So(err, ShouldBeNil)

				Convey("Then it should fail", func() {
					So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
					So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
				})
			})

			Convey("And getting a list of topologies", func() {
				res, js, err := cli.Do("GET", "/topologies", nil)
				So(err, ShouldBeNil)
				So(res.StatusCode, ShouldEqual, http.StatusOK)

				Convey("Then it should have the new topology", func() {
					So(jscan(js, "/topologies[0]/name"), ShouldEqual, "test_topology")
				})
			})
		})

		Convey("When getting a list of topologies", func() {
			res, js, err := cli.Do("GET", "/topologies", nil)
			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then it should have an empty response", func() {
				So(jscan(js, "/topologies"), ShouldBeEmpty)
			})
		})

		Convey("When getting a nonexistent topology", func() {
			res, js, err := cli.Do("GET", "/topologies/test_topology", nil)
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.StatusCode, ShouldEqual, http.StatusNotFound)

				Convey("And the respons should have error information", func() {
					So(jscan(js, "/error"), ShouldNotBeNil)
				})
			})
		})

		Convey("When deleting a nonexistent topology", func() {
			res, _, err := cli.Do("DELETE", "/topologies/test_topology", nil)
			So(err, ShouldBeNil)

			Convey("Then it shouldn't fail", func() {
				So(res.StatusCode, ShouldEqual, http.StatusOK) // TODO: This should be replaced with 204 later
			})
		})
	})
}

func TestTopologiesCreateInvalidValues(t *testing.T) {
	s := createTestServer()
	defer s.Close()
	cli := s.Client()

	Convey("Given an API server", t, func() {
		Convey("When posting a request missing required fields", func() {
			res, js, err := cli.Do("POST", "/topologies", map[string]interface{}{})
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then it should have meta information", func() {
				So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
			})
		})

		Convey("When posting a broken JSON request", func() {
			res, js, err := cli.Do("POST", "/topologies", json.RawMessage("{broken}"))
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(jscan(js, "/error"), ShouldNotBeNil)
			})
		})

		Convey("When posting an integer name", func() {
			res, js, err := cli.Do("POST", "/topologies", map[string]interface{}{
				"name": 1,
			})
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
			})
		})

		Convey("When posting an empty name", func() {
			res, js, err := cli.Do("POST", "/topologies", map[string]interface{}{
				"name": "",
			})
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(jscan(js, "/error/meta/name[0]"), ShouldNotBeBlank)
			})
		})
	})
}

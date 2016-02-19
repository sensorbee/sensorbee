package client

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"gopkg.in/sensorbee/sensorbee.v0/server/response"
	"gopkg.in/sensorbee/sensorbee.v0/server/testutil"
	"net/http"
	"testing"
)

func TestSources(t *testing.T) {
	s := testutil.NewServer()
	defer s.Close()
	r := newTestRequester(s)

	Convey("Given an API server with a topology", t, func() {
		res, _, err := do(r, Post, "/topologies", map[string]interface{}{
			"name": "test_topology",
		})
		So(err, ShouldBeNil)
		So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)
		Reset(func() {
			do(r, Delete, "/topologies/test_topology", nil)
		})

		type indexRes struct {
			Topology string             `json:"topology"`
			Count    int                `json:"count"`
			Sources  []*response.Source `json:"sources"`
		}

		type showRes struct {
			Topology string           `json:"topology"`
			Source   *response.Source `json:"source"`
		}

		Convey("When listing sources", func() {
			res, _, err := do(r, Get, "/topologies/test_topology/sources", nil)
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then the list should be empty", func() {
				s := indexRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				So(s.Count, ShouldEqual, 0)
				So(s.Sources, ShouldBeEmpty)
			})
		})

		Convey("When getting a nonexistent source", func() {
			res, _, err := do(r, Get, "/topologies/test_topology/sources/test_source", nil)
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When adding a source", func() {
			res, _, err := do(r, Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `CREATE PAUSED SOURCE test_source TYPE dummy;`,
			})
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then the list should contain the source", func() {
				res, _, err := do(r, Get, "/topologies/test_topology/sources", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)
				s := indexRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				So(s.Count, ShouldEqual, 1)

				src := s.Sources[0]

				Convey("And it should have correct information", func() {
					So(src.NodeType, ShouldEqual, "source")
					So(src.Name, ShouldEqual, "test_source")
					So(src.State, ShouldEqual, "paused")
				})

				Convey("And the response shouldn't contain state", func() {
					So(len(src.Status), ShouldEqual, 0)
				})

				// TODO: check Meta when it's added in the server side
			})

			Convey("Then getting the source should succeed", func() {
				res, _, err := do(r, Get, "/topologies/test_topology/sources/test_source", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

				s := showRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				src := s.Source

				Convey("And it should have correct information", func() {
					So(src.NodeType, ShouldEqual, "source")
					So(src.Name, ShouldEqual, "test_source")
					So(src.State, ShouldEqual, "paused")
				})

				Convey("And the response should contain state", func() {
					path := "output_stats.num_sent_total"
					_, err := src.Status.Get(data.MustCompilePath(path))
					So(err, ShouldBeNil)
				})

				// TODO: check Meta when it's added in the server side
			})

			Convey("Then getting the source with nonexistent topology name should fail", func() {
				res, _, err := do(r, Get, "/topologies/test_topology2/sources/test_source", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

package client

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/server/response"
	"pfi/sensorbee/sensorbee/server/testutil"
	"testing"
)

func TestSinks(t *testing.T) {
	s := testutil.NewServer()
	defer s.Close()
	r := newTestRequester(s)

	Convey("Given an API server with a topology", t, func() {
		res, _, err := do(r, Post, "/topologies", map[string]interface{}{
			"name": "test_topology",
		})
		Reset(func() {
			do(r, Delete, "/topologies/test_topology", nil)
		})
		So(err, ShouldBeNil)
		So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

		type indexRes struct {
			Topology string           `json:"topology"`
			Count    int              `json:"count"`
			Sinks    []*response.Sink `json:"sinks"`
		}

		type showRes struct {
			Topology string         `json:"topology"`
			Sink     *response.Sink `json:"sink"`
		}

		Convey("When listing sinks", func() {
			res, _, err := do(r, Get, "/topologies/test_topology/sinks", nil)
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then the list should be empty", func() {
				s := indexRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				So(s.Count, ShouldEqual, 0)
				So(s.Sinks, ShouldBeEmpty)
			})
		})

		Convey("When getting a nonexistent sink", func() {
			res, _, err := do(r, Get, "/topologies/test_topology/sinks/test_sink", nil)
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When adding a sink", func() {
			res, _, err := do(r, Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `CREATE SINK test_sink TYPE stdout;`,
			})
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then the list should contain the sink", func() {
				res, _, err := do(r, Get, "/topologies/test_topology/sinks", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)
				s := indexRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				So(s.Count, ShouldEqual, 1)

				sink := s.Sinks[0]

				Convey("And it should have correct information", func() {
					So(sink.NodeType, ShouldEqual, "sink")
					So(sink.Name, ShouldEqual, "test_sink")
					So(sink.State, ShouldEqual, "running")
				})

				Convey("And the response shouldn't contain state", func() {
					So(len(sink.Status), ShouldEqual, 0)
				})

				// TODO: check Meta when it's added in the server side
			})

			Convey("Then getting the sink should succeed", func() {
				res, _, err := do(r, Get, "/topologies/test_topology/sinks/test_sink", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

				s := showRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				sink := s.Sink

				Convey("And it should have correct information", func() {
					So(sink.NodeType, ShouldEqual, "sink")
					So(sink.Name, ShouldEqual, "test_sink")
					So(sink.State, ShouldEqual, "running")
				})

				Convey("And the response should contain state", func() {
					_, err := sink.Status.Get("input_stats.num_received_total")
					So(err, ShouldBeNil)
				})

				// TODO: check Meta when it's added in the server side
			})

			Convey("Then getting the sink with nonexistent topology name should fail", func() {
				res, _, err := do(r, Get, "/topologies/test_topology2/sinks/test_sink", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

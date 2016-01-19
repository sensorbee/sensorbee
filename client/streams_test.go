package client

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/data"
	"pfi/sensorbee/sensorbee/server/response"
	"pfi/sensorbee/sensorbee/server/testutil"
	"testing"
)

func TestStreams(t *testing.T) {
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
			Streams  []*response.Stream `json:"streams"`
		}

		type showRes struct {
			Topology string           `json:"topology"`
			Stream   *response.Stream `json:"stream"`
		}

		Convey("When listing streams", func() {
			res, _, err := do(r, Get, "/topologies/test_topology/streams", nil)
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then the list should be empty", func() {
				s := indexRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				So(s.Count, ShouldEqual, 0)
				So(s.Streams, ShouldBeEmpty)
			})
		})

		Convey("When getting a nonexistent stream", func() {
			res, _, err := do(r, Get, "/topologies/test_topology/streams/test_stream", nil)
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When adding a stream", func() {
			res, _, err := do(r, Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `CREATE PAUSED SOURCE test_source TYPE dummy;
							CREATE STREAM test_stream AS SELECT ISTREAM * FROM test_source [RANGE 1 TUPLES];`,
			})
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then the list should contain the stream", func() {
				res, _, err := do(r, Get, "/topologies/test_topology/streams", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)
				s := indexRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				So(s.Count, ShouldEqual, 1)

				strm := s.Streams[0]

				Convey("And it should have correct information", func() {
					So(strm.NodeType, ShouldEqual, "box")
					So(strm.Name, ShouldEqual, "test_stream")
					So(strm.State, ShouldEqual, "running")
				})

				Convey("And the response shouldn't contain state", func() {
					So(len(strm.Status), ShouldEqual, 0)
				})

				// TODO: check Meta when it's added in the server side
			})

			Convey("Then getting the stream should succeed", func() {
				res, _, err := do(r, Get, "/topologies/test_topology/streams/test_stream", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

				s := showRes{}
				So(res.ReadJSON(&s), ShouldBeNil)
				So(s.Topology, ShouldEqual, "test_topology")
				strm := s.Stream

				Convey("And it should have correct information", func() {
					So(strm.NodeType, ShouldEqual, "box")
					So(strm.Name, ShouldEqual, "test_stream")
					So(strm.State, ShouldEqual, "running")
				})

				Convey("And the response should contain state", func() {
					path := "output_stats.num_sent_total"
					_, err := strm.Status.Get(data.MustCompilePath(path))
					So(err, ShouldBeNil)
				})

				// TODO: check Meta when it's added in the server side
			})

			Convey("Then getting the stream with nonexistent topology name should fail", func() {
				res, _, err := do(r, Get, "/topologies/test_topology2/streams/test_stream", nil)
				So(err, ShouldBeNil)
				So(res.Raw.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

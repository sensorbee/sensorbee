package client

// TODO: replace tests with a richer client

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"pfi/sensorbee/sensorbee/data"
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
			res, _, err := do(r, Post, "/topologies", map[string]interface{}{})
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then it should have meta information", func() {
				e, err := res.Error()
				So(err, ShouldBeNil)
				So(e.Meta["name"].(data.Array)[0], ShouldNotEqual, "")
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

func TestTopologiesQueries(t *testing.T) {
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

		// TODO: add more tests
		Convey("When creating a sink", func() {
			res, _, err := do(r, Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `CREATE SINK stdout TYPE stdout;`,
			})
			So(err, ShouldBeNil)

			Convey("Then it should succeed", func() {
				So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)
			})

			// TODO: check the response json
		})
	})
}

func TestTopologiesQueriesSelectStmt(t *testing.T) {
	// TODO: Becaues results from a SELECT stmt needs to be returned through
	// hijacking, a real HTTP server is required. Support Hijack method in test
	// ResponseWriter not to use a real HTTP server.
	testutil.TestAPIWithRealHTTPServer = true

	s := testutil.NewServer()
	defer func() {
		testutil.TestAPIWithRealHTTPServer = false
		s.Close()
	}()
	r := newTestRequester(s)

	Convey("Given an API server with a topology having a paused source", t, func() {
		res, _, err := do(r, Post, "/topologies", map[string]interface{}{
			"name": "test_topology",
		})
		Reset(func() {
			do(r, Delete, "/topologies/test_topology", nil)
		})
		So(err, ShouldBeNil)
		So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

		res, _, err = do(r, Post, "/topologies/test_topology/queries", map[string]interface{}{
			"queries": `CREATE PAUSED SOURCE source TYPE dummy;`,
		})
		So(err, ShouldBeNil)
		So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

		Convey("When issueing a SELECT stmt", func() {
			streamRes, err := r.Do(Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `SELECT ISTREAM * FROM source [RANGE 1 TUPLES];`,
			})
			So(err, ShouldBeNil)
			Reset(func() {
				streamRes.Close()
			})
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			res, err := r.Do(Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `RESUME SOURCE source;`,
			})
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then it should receive all tuples and stop", func() {
				ch, err := streamRes.ReadStreamJSON()
				So(err, ShouldBeNil)

				for i := 0; i < 4; i++ {
					js, ok := <-ch
					So(ok, ShouldBeTrue)
					So(jscan(js, "/int"), ShouldEqual, i)
				}
				_, ok := <-ch
				So(ok, ShouldBeFalse)
			})
		})

		// TODO: add invalid cases
	})
}

func TestTopologiesQueriesSelectUnionStmt(t *testing.T) {
	// TODO: Becaues results from a SELECT stmt needs to be returned through
	// hijacking, a real HTTP server is required. Support Hijack method in test
	// ResponseWriter not to use a real HTTP server.
	testutil.TestAPIWithRealHTTPServer = true

	s := testutil.NewServer()
	defer func() {
		testutil.TestAPIWithRealHTTPServer = false
		s.Close()
	}()
	r := newTestRequester(s)

	Convey("Given an API server with a topology having a paused source", t, func() {
		res, _, err := do(r, Post, "/topologies", map[string]interface{}{
			"name": "test_topology",
		})
		Reset(func() {
			do(r, Delete, "/topologies/test_topology", nil)
		})
		So(err, ShouldBeNil)
		So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

		res, _, err = do(r, Post, "/topologies/test_topology/queries", map[string]interface{}{
			"queries": `CREATE PAUSED SOURCE source TYPE dummy;`,
		})
		So(err, ShouldBeNil)
		So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

		Convey("When issueing a SELECT stmt", func() {
			streamRes, err := r.Do(Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `SELECT ISTREAM * FROM source [RANGE 1 TUPLES] WHERE int % 2 = 0
					UNION ALL SELECT ISTREAM * FROM source [RANGE 1 TUPLES] WHERE int % 2 = 1;`,
			})
			So(err, ShouldBeNil)
			Reset(func() {
				streamRes.Close()
			})
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			res, err := r.Do(Post, "/topologies/test_topology/queries", map[string]interface{}{
				"queries": `RESUME SOURCE source;`,
			})
			So(err, ShouldBeNil)
			So(res.Raw.StatusCode, ShouldEqual, http.StatusOK)

			Convey("Then it should receive all tuples and stop", func() {
				ch, err := streamRes.ReadStreamJSON()
				So(err, ShouldBeNil)

				// items will not come in order, so we need to
				found := map[int64]bool{}
				for i := 0; i < 4; i++ {
					js, ok := <-ch
					So(ok, ShouldBeTrue)
					j := int64(jscan(js, "/int").(float64))
					found[j] = true
				}
				So(found, ShouldResemble, map[int64]bool{
					0: true, 1: true, 2: true, 3: true,
				})
				_, ok := <-ch
				So(ok, ShouldBeFalse)
			})
		})

		// TODO: add invalid cases
	})
}

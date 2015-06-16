package bql_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql"
	"testing"
)

func TestCreateSourceStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder", t, func() {
		dt := newTestDynamicTopology()
		Reset(func() {
			dt.Stop()
		})
		tb := bql.NewTopologyBuilder(dt)

		Convey("When running CREATE SOURCE with a dummy source", func() {
			err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE dummy`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE SOURCE using the same name", func() {
				err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE dummy`)

				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "already")
				})
			})

			Convey("And when running another CREATE SOURCE using a different name", func() {
				err := addBQLToTopology(tb, `CREATE SOURCE src TYPE dummy`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running CREATE SOURCE with valid parameters", func() {
			err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE dummy WITH num=4`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When running CREATE SOURCE with invalid parameters", func() {
			err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE dummy WITH num=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "num: cannot convert string 'bar' into integer")
			})
		})

		Convey("When running CREATE SOURCE with unknown parameters", func() {
			err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE dummy WITH foo=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source parameter: foo")
			})
		})

		Convey("When running CREATE SOURCE with an unknown source type", func() {
			err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source type: foo")
			})
		})
	})
}

func TestCreateStreamFromSourceStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder with a source", t, func() {
		dt := newTestDynamicTopology()
		Reset(func() {
			unblockSource(dt, "hoge_src")
			dt.Stop()
		})
		tb := bql.NewTopologyBuilder(dt)
		err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE blocking_dummy`)
		So(err, ShouldBeNil)

		Convey("When running CREATE STREAM FROM SOURCE on an existing source", func() {
			err := addBQLToTopology(tb, `CREATE STREAM s FROM SOURCE hoge`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE STREAM FROM SOURCE with the same name", func() {
				err := addBQLToTopology(tb, `CREATE STREAM s FROM SOURCE hoge`)

				Convey("Then an error should be returned the second time", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "already")
				})
			})

			Convey("And when running another CREATE STREAM FROM SOURCE with a different name", func() {
				err := addBQLToTopology(tb, `CREATE STREAM t FROM SOURCE hoge`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running CREATE STREAM FROM SOURCE on a non-existing source", func() {
			err := addBQLToTopology(tb, `CREATE STREAM s FROM SOURCE foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})

		Convey("When running CREATE STREAM FROM SOURCE with stream and source name the same", func() {
			err := addBQLToTopology(tb, `CREATE STREAM hoge FROM SOURCE hoge`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestCreateStreamFromSourceExtStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder", t, func() {
		dt := newTestDynamicTopology()
		Reset(func() {
			dt.Stop()
		})
		tb := bql.NewTopologyBuilder(dt)

		Convey("When running CREATE STREAM FROM SOURCE (combo) with a dummy source", func() {
			err := addBQLToTopology(tb, `CREATE STREAM hoge FROM dummy SOURCE`)
			So(err, ShouldBeNil)

			Convey("And when running another CREATE STREAM FROM SOURCE (combo) using the same name", func() {
				err := addBQLToTopology(tb, `CREATE STREAM hoge FROM dummy SOURCE`)

				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "already")
				})
			})
		})
		Convey("When running CREATE STREAM FROM SOURCE (combo) with invalid parameters", func() {
			err := addBQLToTopology(tb, `CREATE STREAM hoge FROM dummy SOURCE WITH foo=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source parameter: foo")
			})
		})
		Convey("When running CREATE STREAM FROM SOURCE (combo) with an unknown source type", func() {
			err := addBQLToTopology(tb, `CREATE STREAM hoge FROM foo SOURCE`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source type: foo")
			})
		})
	})
}

func TestCreateStreamAsSelectStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder with source and stream", t, func() {
		dt := newTestDynamicTopology()
		Reset(func() {
			unblockSource(dt, "hoge_src")
			dt.Stop()
		})
		tb := bql.NewTopologyBuilder(dt)
		err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE blocking_dummy`)
		So(err, ShouldBeNil)
		err = addBQLToTopology(tb, `CREATE STREAM s FROM SOURCE hoge`)
		So(err, ShouldBeNil)

		Convey("When running CREATE STREAM AS SELECT on an existing stream", func() {
			err := addBQLToTopology(tb, `CREATE STREAM t AS SELECT ISTREAM int FROM
                s [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE STREAM AS SELECT with the same name", func() {
				err := addBQLToTopology(tb, `CREATE STREAM t AS SELECT ISTREAM int FROM
                s [RANGE 1 TUPLES] WHERE int=1`)

				Convey("Then an error should be returned the second time", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "already")
				})
			})

			Convey("And when running another CREATE STREAM AS SELECT with a different name", func() {
				err := addBQLToTopology(tb, `CREATE STREAM u AS SELECT ISTREAM int FROM
                s [RANGE 1 TUPLES] WHERE int=1`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running CREATE STREAM AS SELECT on a source name", func() {
			err := addBQLToTopology(tb, `CREATE STREAM t AS SELECT ISTREAM int FROM
                hoge [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})

		Convey("When running CREATE STREAM AS SELECT on a non-existing stream", func() {
			err := addBQLToTopology(tb, `CREATE STREAM t AS SELECT ISTREAM int FROM
                bar [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})

		Convey("When running CREATE STREAM AS SELECT with input and output name the same", func() {
			err := addBQLToTopology(tb, `CREATE STREAM s AS SELECT ISTREAM int FROM
                s [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "already")
			})
		})
	})
}

func TestCreateSinkStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder", t, func() {
		dt := newTestDynamicTopology()
		Reset(func() {
			dt.Stop()
		})
		tb := bql.NewTopologyBuilder(dt)

		Convey("When running CREATE SINK with a collector source", func() {
			err := addBQLToTopology(tb, `CREATE SINK hoge TYPE collector`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE SINK using the same name", func() {
				err := addBQLToTopology(tb, `CREATE SINK hoge TYPE collector`)

				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, "already")
				})
			})
		})
		Convey("When running CREATE SINK with invalid parameters", func() {
			err := addBQLToTopology(tb, `CREATE SINK hoge TYPE collector WITH foo=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown sink parameter: foo")
			})
		})
		Convey("When running CREATE SINK with an unknown sink type", func() {
			err := addBQLToTopology(tb, `CREATE SINK hoge TYPE foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown sink type: foo")
			})
		})
	})
}

func TestInsertIntoSelectStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder with source, stream, and sink", t, func() {
		dt := newTestDynamicTopology()
		Reset(func() {
			unblockSource(dt, "hoge_src")
			dt.Stop()
		})
		tb := bql.NewTopologyBuilder(dt)
		err := addBQLToTopology(tb, `CREATE SOURCE hoge TYPE blocking_dummy`)
		So(err, ShouldBeNil)
		err = addBQLToTopology(tb, `CREATE STREAM s FROM SOURCE hoge`)
		So(err, ShouldBeNil)
		err = addBQLToTopology(tb, `CREATE SINK foo TYPE collector`)
		So(err, ShouldBeNil)

		Convey("When running INSERT INTO SELECT on an existing stream and sink", func() {
			err := addBQLToTopology(tb, `INSERT INTO foo SELECT int FROM s`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another INSERT INTO SELECT", func() {
				err := addBQLToTopology(tb, `INSERT INTO foo SELECT int FROM s`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running INSERT INTO SELECT on a source name", func() {
			err := addBQLToTopology(tb, `INSERT INTO foo SELECT int FROM hoge`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})

		Convey("When running INSERT INTO SELECT on a non-existing stream", func() {
			err := addBQLToTopology(tb, `INSERT INTO foo SELECT int FROM foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})

		Convey("When running INSERT INTO SELECT on a non-existing sink", func() {
			err := addBQLToTopology(tb, `INSERT INTO bar SELECT int FROM s`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})
	})
}

func TestMultipleStatements(t *testing.T) {
	Convey("Given an empty BQL TopologyBuilder", t, func() {
		dt := newTestDynamicTopology()
		Reset(func() {
			dt.Stop()
		})
		tb := bql.NewTopologyBuilder(dt)

		Convey("When issuing multiple commands in a good order", func() {
			stmts := `
			CREATE STREAM source FROM blocking_dummy SOURCE WITH num=4;
			CREATE STREAM box AS SELECT
			  ISTREAM int, str((int+1) % 3) AS x FROM source
			  [RANGE 2 SECONDS] WHERE int % 2 = 0;
			CREATE SINK snk TYPE collector;
			INSERT INTO snk SELECT * FROM box;
			`
			err := addBQLToTopology(tb, stmts)
			So(unblockSource(dt, "source"), ShouldBeNil)

			Convey("Then setup should succeed", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When issuing multiple commands in a bad order", func() {
			stmts := `
			CREATE STREAM source FROM blocking_dummy SOURCE WITH num=4;
			CREATE STREAM box AS SELECT
			  ISTREAM int, str((int+1) % 3) AS x FROM source
			  [RANGE 2 SECONDS] WHERE int % 2 = 0;
			INSERT INTO snk SELECT * FROM box;
			CREATE SINK snk TYPE collector;
			`
			err := addBQLToTopology(tb, stmts)
			So(unblockSource(dt, "source"), ShouldBeNil)

			Convey("Then setup should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

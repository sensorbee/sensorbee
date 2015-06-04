package bql_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql"
	"testing"
)

func TestCreateSourceStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder", t, func() {
		tb := bql.NewTopologyBuilder()

		Convey("When running CREATE SOURCE with a dummy source", func() {
			err := tb.BQL(`CREATE SOURCE hoge TYPE dummy`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE SOURCE using the same name", func() {
				err := tb.BQL(`CREATE SOURCE hoge TYPE dummy`)

				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "there is already a source called 'hoge-src'")
				})
			})

			Convey("And when running another CREATE SOURCE using a different name", func() {
				err := tb.BQL(`CREATE SOURCE src TYPE dummy`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running CREATE SOURCE with valid parameters", func() {
			err := tb.BQL(`CREATE SOURCE hoge TYPE dummy WITH num=4`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When running CREATE SOURCE with invalid parameters", func() {
			err := tb.BQL(`CREATE SOURCE hoge TYPE dummy WITH num=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "num: cannot convert string 'bar' into integer")
			})
		})

		Convey("When running CREATE SOURCE with unknown parameters", func() {
			err := tb.BQL(`CREATE SOURCE hoge TYPE dummy WITH foo=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source parameter: foo")
			})
		})

		Convey("When running CREATE SOURCE with an unknown source type", func() {
			err := tb.BQL(`CREATE SOURCE hoge TYPE foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source type: foo")
			})
		})
	})
}

func TestCreateStreamFromSourceStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder with a source", t, func() {
		tb := bql.NewTopologyBuilder()
		err := tb.BQL(`CREATE SOURCE hoge TYPE dummy`)
		So(err, ShouldBeNil)

		Convey("When running CREATE STREAM FROM SOURCE on an existing source", func() {
			err := tb.BQL(`CREATE STREAM s FROM SOURCE hoge`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE STREAM FROM SOURCE with the same name", func() {
				err := tb.BQL(`CREATE STREAM s FROM SOURCE hoge`)

				Convey("Then an error should be returned the second time", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "there is already a box called 's'")
				})
			})

			Convey("And when running another CREATE STREAM FROM SOURCE with a different name", func() {
				err := tb.BQL(`CREATE STREAM t FROM SOURCE hoge`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running CREATE STREAM FROM SOURCE on a non-existing source", func() {
			err := tb.BQL(`CREATE STREAM s FROM SOURCE foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "there is no box or source named 'foo-src'")
			})
		})

		Convey("When running CREATE STREAM FROM SOURCE with stream and source name the same", func() {
			err := tb.BQL(`CREATE STREAM hoge FROM SOURCE hoge`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestCreateStreamFromSourceExtStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder", t, func() {
		tb := bql.NewTopologyBuilder()

		Convey("When running CREATE STREAM FROM SOURCE (combo) with a dummy source", func() {
			err := tb.BQL(`CREATE STREAM hoge FROM dummy SOURCE`)
			So(err, ShouldBeNil)

			Convey("And when running another CREATE STREAM FROM SOURCE (combo) using the same name", func() {
				err := tb.BQL(`CREATE STREAM hoge FROM dummy SOURCE`)

				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "there is already a source called 'hoge'")
				})
			})
		})
		Convey("When running CREATE STREAM FROM SOURCE (combo) with invalid parameters", func() {
			err := tb.BQL(`CREATE STREAM hoge FROM dummy SOURCE WITH foo=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source parameter: foo")
			})
		})
		Convey("When running CREATE STREAM FROM SOURCE (combo) with an unknown source type", func() {
			err := tb.BQL(`CREATE STREAM hoge FROM foo SOURCE`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown source type: foo")
			})
		})
	})
}

func TestCreateStreamAsSelectStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder with source and stream", t, func() {
		tb := bql.NewTopologyBuilder()
		err := tb.BQL(`CREATE SOURCE hoge TYPE dummy`)
		So(err, ShouldBeNil)
		err = tb.BQL(`CREATE STREAM s FROM SOURCE hoge`)
		So(err, ShouldBeNil)

		Convey("When running CREATE STREAM AS SELECT on an existing stream", func() {
			err := tb.BQL(`CREATE STREAM t AS SELECT ISTREAM(int) FROM
                s [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE STREAM AS SELECT with the same name", func() {
				err := tb.BQL(`CREATE STREAM t AS SELECT ISTREAM(int) FROM
                s [RANGE 1 TUPLES] WHERE int=1`)

				Convey("Then an error should be returned the second time", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "there is already a box called 't'")
				})
			})

			Convey("And when running another CREATE STREAM AS SELECT with a different name", func() {
				err := tb.BQL(`CREATE STREAM u AS SELECT ISTREAM(int) FROM
                s [RANGE 1 TUPLES] WHERE int=1`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running CREATE STREAM AS SELECT on a source name", func() {
			err := tb.BQL(`CREATE STREAM t AS SELECT ISTREAM(int) FROM
                hoge [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "there is no box or source named 'hoge'")
			})
		})

		Convey("When running CREATE STREAM AS SELECT on a non-existing stream", func() {
			err := tb.BQL(`CREATE STREAM t AS SELECT ISTREAM(int) FROM
                bar [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "there is no box or source named 'bar'")
			})
		})

		Convey("When running CREATE STREAM AS SELECT with input and output name the same", func() {
			err := tb.BQL(`CREATE STREAM s AS SELECT ISTREAM(int) FROM
                s [RANGE 2 SECONDS] WHERE int=2`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "there is already a box called 's'")
			})
		})
	})
}

func TestCreateSinkStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder", t, func() {
		tb := bql.NewTopologyBuilder()

		Convey("When running CREATE SINK with a collector source", func() {
			err := tb.BQL(`CREATE SINK hoge TYPE collector`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another CREATE SINK using the same name", func() {
				err := tb.BQL(`CREATE SINK hoge TYPE collector`)

				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "there is already a sink called 'hoge'")
				})
			})
		})
		Convey("When running CREATE SINK with invalid parameters", func() {
			err := tb.BQL(`CREATE SINK hoge TYPE collector WITH foo=bar`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown sink parameter: foo")
			})
		})
		Convey("When running CREATE SINK with an unknown sink type", func() {
			err := tb.BQL(`CREATE SINK hoge TYPE foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown sink type: foo")
			})
		})
	})
}

func TestInsertIntoSelectStmt(t *testing.T) {
	Convey("Given a BQL TopologyBuilder with source, stream, and sink", t, func() {
		tb := bql.NewTopologyBuilder()
		err := tb.BQL(`CREATE SOURCE hoge TYPE dummy`)
		So(err, ShouldBeNil)
		err = tb.BQL(`CREATE STREAM s FROM SOURCE hoge`)
		So(err, ShouldBeNil)
		err = tb.BQL(`CREATE SINK foo TYPE collector`)
		So(err, ShouldBeNil)

		Convey("When running INSERT INTO SELECT on an existing stream and sink", func() {
			err := tb.BQL(`INSERT INTO foo SELECT int FROM s`)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And when running another INSERT INTO SELECT", func() {
				err := tb.BQL(`INSERT INTO foo SELECT int FROM s`)

				Convey("Then there should be no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When running INSERT INTO SELECT on a source name", func() {
			err := tb.BQL(`INSERT INTO foo SELECT int FROM hoge`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "there is no box or source named 'hoge'")
			})
		})

		Convey("When running INSERT INTO SELECT on a non-existing stream", func() {
			err := tb.BQL(`INSERT INTO foo SELECT int FROM foo`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "there is no box or source named 'foo'")
			})
		})

		Convey("When running INSERT INTO SELECT on a non-existing sink", func() {
			err := tb.BQL(`INSERT INTO bar SELECT int FROM s`)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "there is no sink named 'bar'")
			})
		})
	})
}

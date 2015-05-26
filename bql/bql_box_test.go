package bql

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
	"time"
)

func TestBqlBox(t *testing.T) {
	tup1 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(1),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	tup2 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(2),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 1, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 1, 0, time.UTC),
		BatchID:       7,
	}
	tup3 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(3),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 2, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 2, 0, time.UTC),
		BatchID:       7,
	}
	tup4 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(4),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 3, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 3, 0, time.UTC),
		BatchID:       7,
	}
	config := core.Configuration{TupleTraceEnabled: 0}
	ctx := newTestContext(config)

	setupTopology := func(stmt string) (core.StaticTopology, *TupleCollectorSink, error) {
		// source
		tb := NewTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1, &tup2, &tup3, &tup4},
		}
		tb.AddSource("source", so)
		// issue BQL statement (inserts box)
		err := tb.BQL(stmt)
		if err != nil {
			return nil, nil, err
		}
		// sink
		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("box")
		t, err := tb.Build()
		return t, si, err
	}

	Convey("Given an ISTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"ISTREAM(int) FROM source [RANGE 2 SECONDS] WHERE int = 2" // actually int % 2 == 0
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			t.Run(ctx)
			Convey("Then the sink receives 2 tuples", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 2)

				Convey("And the first tuple has tup2's data and timestamp", func() {
					si.Tuples[0].InputName = "source"
					So(*si.Tuples[0], ShouldResemble, tup2)
				})

				Convey("And the second tuple has tup4's data and timestamp", func() {
					si.Tuples[1].InputName = "source"
					So(*si.Tuples[1], ShouldResemble, tup4)
				})
			})
		})
	})

	Convey("Given a DSTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"DSTREAM(int) FROM source [RANGE 2 SECONDS] WHERE int = 2" // actually int % 2 == 0
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			t.Run(ctx)
			Convey("Then the sink receives 0 tuples", func() {
				So(si.Tuples, ShouldBeNil)
			})
		})
	})

	Convey("Given a DSTREAM/1 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"DSTREAM(int) FROM source [RANGE 1 SECONDS] WHERE int = 2" // actually int % 2 == 0
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			t.Run(ctx)
			Convey("Then the sink receives 1 tuple", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 1)

				Convey("And it has tup2's data and tup4's timestamp", func() {
					So(si.Tuples[0].Data, ShouldResemble, tup2.Data)
					So(si.Tuples[0].Timestamp, ShouldResemble, tup4.Timestamp)
				})
			})
		})
	})

	Convey("Given an RSTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"RSTREAM(int) FROM source [RANGE 2 SECONDS] WHERE int = 2" // actually int % 2 == 0
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			t.Run(ctx)
			Convey("Then the sink receives 4 tuples", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 4)

				Convey("And the first tuple has tup2's data and timestamp", func() {
					So(si.Tuples[0].Data, ShouldResemble, tup2.Data)
					So(si.Tuples[0].Timestamp, ShouldResemble, tup2.Timestamp)
				})

				Convey("And the second tuple has tup2's data and tup3's timestamp", func() {
					So(si.Tuples[1].Data, ShouldResemble, tup2.Data)
					So(si.Tuples[1].Timestamp, ShouldResemble, tup3.Timestamp)
				})

				Convey("And the third tuple has tup2's data and tup4's timestamp", func() {
					So(si.Tuples[2].Data, ShouldResemble, tup2.Data)
					So(si.Tuples[2].Timestamp, ShouldResemble, tup4.Timestamp)
				})

				Convey("And the fourth tuple has tup4's data and timestamp", func() {
					So(si.Tuples[3].Data, ShouldResemble, tup4.Data)
					So(si.Tuples[3].Timestamp, ShouldResemble, tup4.Timestamp)
				})
			})
		})
	})
}

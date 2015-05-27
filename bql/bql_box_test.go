package bql

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"testing"
)

func setupTopology(stmt string) (core.StaticTopology, *TupleCollectorSink, error) {
	// create a stream from a dummy source
	tb := NewTopologyBuilder()
	err := tb.BQL("CREATE STREAM source FROM dummy SOURCE WITH num=4")
	if err != nil {
		return nil, nil, err
	}
	// issue BQL statement (inserts box)
	err = tb.BQL(stmt)
	if err != nil {
		return nil, nil, err
	}
	// sink
	si := &TupleCollectorSink{}
	tb.stb.AddSink("sink", si).Input("box")
	t, err := tb.Build()
	return t, si, err
}

func TestIstreamSecondsBqlBox(t *testing.T) {
	tuples := getTuples(4)
	tup2 := *tuples[1]
	tup4 := *tuples[3]

	config := core.Configuration{TupleTraceEnabled: 0}
	ctx := newTestContext(config)

	Convey("Given an ISTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"ISTREAM(int) FROM source [RANGE 2 SECONDS] WHERE int % 2 = 0"
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			err := t.Run(ctx)
			So(err, ShouldBeNil)

			Convey("Then the sink receives 2 tuples", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 2)

				Convey("And the first tuple has tup2's data and timestamp", func() {
					si.Tuples[0].InputName = "input"
					So(*si.Tuples[0], ShouldResemble, tup2)
				})

				Convey("And the second tuple has tup4's data and timestamp", func() {
					si.Tuples[1].InputName = "input"
					So(*si.Tuples[1], ShouldResemble, tup4)
				})
			})
		})
	})
}

func TestDstreamSecondsBqlBox(t *testing.T) {
	tuples := getTuples(4)
	tup2 := *tuples[1]
	tup4 := *tuples[3]

	config := core.Configuration{TupleTraceEnabled: 0}
	ctx := newTestContext(config)

	Convey("Given a DSTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"DSTREAM(int) FROM source [RANGE 2 SECONDS] WHERE int % 2 = 0"
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			err := t.Run(ctx)
			So(err, ShouldBeNil)

			Convey("Then the sink receives 0 tuples", func() {
				So(si.Tuples, ShouldBeNil)
			})
		})
	})

	Convey("Given a DSTREAM/1 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"DSTREAM(int) FROM source [RANGE 1 SECONDS] WHERE int % 2 = 0"
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			err := t.Run(ctx)
			So(err, ShouldBeNil)

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
}

func TestRstreamSecondsBqlBox(t *testing.T) {
	tuples := getTuples(4)
	tup2 := *tuples[1]
	tup3 := *tuples[2]
	tup4 := *tuples[3]

	config := core.Configuration{TupleTraceEnabled: 0}
	ctx := newTestContext(config)

	Convey("Given an RSTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"RSTREAM(int) FROM source [RANGE 2 SECONDS] WHERE int % 2 = 0"
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			err := t.Run(ctx)
			So(err, ShouldBeNil)

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

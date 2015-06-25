package bql

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"testing"
)

func setupTopology(stmt string) (core.DynamicTopology, error) {
	// create a stream from a dummy source
	dt := newTestDynamicTopology()
	tb := NewTopologyBuilder(dt)
	err := addBQLToTopology(tb, "CREATE SOURCE source TYPE blocking_dummy WITH num=4")
	if err != nil {
		return nil, err
	}
	// issue BQL statement (inserts box)
	if err := addBQLToTopology(tb, stmt); err != nil {
		return nil, err
	}
	// sink
	err = addBQLToTopology(tb, "CREATE SINK snk TYPE collector")
	if err != nil {
		return nil, err
	}
	// connect box and sink
	err = addBQLToTopology(tb, "INSERT INTO snk SELECT * FROM box")
	if err != nil {
		return nil, err
	}
	return dt, err
}

func TestBasicBqlBoxConnectivity(t *testing.T) {
	tuples := mkTuples(4)
	tup2 := *tuples[1]
	tup4 := *tuples[3]

	Convey("Given an ISTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"ISTREAM int, str((int+1) % 3) AS x FROM source [RANGE 2 SECONDS] WHERE int % 2 = 0"
		dt, err := setupTopology(s)
		So(err, ShouldBeNil)
		Reset(func() {
			dt.Stop()
		})

		So(unblockSource(dt, "source"), ShouldBeNil)

		sin, err := dt.Sink("snk")
		So(err, ShouldBeNil)
		si := sin.Sink().(*tupleCollectorSink)

		Convey("When 4 tuples are emitted by the source", func() {
			tup2.Data["x"] = tuple.String(fmt.Sprintf("%d", ((2 + 1) % 3)))
			tup4.Data["x"] = tuple.String(fmt.Sprintf("%d", ((4 + 1) % 3)))

			Convey("Then the sink receives 2 tuples", func() {
				si.Wait(2)
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

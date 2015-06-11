package bql_test

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"pfi/sensorbee/sensorbee/io/dummy_sink"
	"pfi/sensorbee/sensorbee/io/dummy_source"
	"testing"
)

func setupTopology(stmt string) (core.StaticTopology, *dummy_sink.TupleCollectorSink, error) {
	// create a stream from a dummy source
	tb := bql.NewTopologyBuilder()
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
	err = tb.BQL("CREATE SINK snk TYPE collector")
	if err != nil {
		return nil, nil, err
	}
	// connect box and sink
	err = tb.BQL("INSERT INTO snk SELECT * FROM box")
	if err != nil {
		return nil, nil, err
	}
	si_, ok := tb.GetSink("snk")
	if !ok {
		return nil, nil, fmt.Errorf("unable to get sink from TopologyBuilder")
	}
	si, ok := si_.(*dummy_sink.TupleCollectorSink)
	if !ok {
		return nil, nil, fmt.Errorf("returned sink has unexpected type: %v", si_)
	}
	t, err := tb.Build()
	return t, si, err
}

func TestBasicBqlBoxConnectivity(t *testing.T) {
	tuples := dummy_source.MkTuples(4)
	tup2 := *tuples[1]
	tup4 := *tuples[3]

	config := core.Configuration{TupleTraceEnabled: 0}
	ctx := newTestContext(config)

	Convey("Given an ISTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"ISTREAM int, str((int+1) % 3) AS x FROM source [RANGE 2 SECONDS] WHERE int % 2 = 0"
		t, si, err := setupTopology(s)
		So(err, ShouldBeNil)

		Convey("When 4 tuples are emitted by the source", func() {
			err := t.Run(ctx)
			So(err, ShouldBeNil)
			tup2.Data["x"] = tuple.String(fmt.Sprintf("%d", ((2 + 1) % 3)))
			tup4.Data["x"] = tuple.String(fmt.Sprintf("%d", ((4 + 1) % 3)))

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

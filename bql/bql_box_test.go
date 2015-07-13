package bql

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func setupTopology(stmt string) (*TopologyBuilder, error) {
	// create a stream from a dummy source
	dt := newTestTopology()
	tb, err := NewTopologyBuilder(dt)
	if err != nil {
		return nil, err
	}
	err = addBQLToTopology(tb, "CREATE PAUSED SOURCE source TYPE dummy WITH num=4")
	if err != nil {
		return nil, err
	}
	// issue BQL statement (inserts box)
	if err := addBQLToTopology(tb, stmt); err != nil {
		return nil, err
	}
	// sink
	err = addBQLToTopology(tb, `
		CREATE SINK snk TYPE collector;
		INSERT INTO snk FROM box;
		RESUME SOURCE source;`)
	if err != nil {
		return nil, err
	}
	return tb, err
}

func TestBasicBQLBoxConnectivity(t *testing.T) {
	tuples := mkTuples(4)
	tup2 := *tuples[1]
	tup4 := *tuples[3]

	Convey("Given an ISTREAM/2 SECONDS BQL statement", t, func() {
		s := "CREATE STREAM box AS SELECT " +
			"ISTREAM int, str((int+1) % 3) AS x FROM source [RANGE 1 TUPLES] WHERE int % 2 = 0"
		tb, err := setupTopology(s)
		So(err, ShouldBeNil)
		dt := tb.Topology()
		Reset(func() {
			dt.Stop()
		})

		sin, err := dt.Sink("snk")
		So(err, ShouldBeNil)
		si := sin.Sink().(*tupleCollectorSink)

		Convey("When 4 tuples are emitted by the source", func() {
			tup2.Data["x"] = data.String(fmt.Sprintf("%d", ((2 + 1) % 3)))
			tup4.Data["x"] = data.String(fmt.Sprintf("%d", ((4 + 1) % 3)))

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

		Convey("When rewinding the source", func() {
			si.Wait(2)
			So(addBQLToTopology(tb, `REWIND SOURCE source;`), ShouldBeNil)

			Convey("Then the sinkreceives tuples again", func() {
				si.Wait(4)
				So(len(si.Tuples), ShouldEqual, 4)
			})
		})
	})
}

func TestBQLBoxJoinCapability(t *testing.T) {
	tuples := mkTuples(4)

	Convey("Given an RSTREAM statement with a lot of joins", t, func() {
		s := `CREATE STREAM box AS SELECT RSTREAM
		source:int AS a, s2:int AS b, duplicate:int AS c, d2:int AS d
		FROM source [RANGE 1 TUPLES],
		     source [RANGE 1 TUPLES] AS s2,
		     duplicate('source', 3) [RANGE 1 TUPLES],
		     duplicate('source', 2) [RANGE 1 TUPLES] AS d2
		`
		tb, err := setupTopology(s)
		So(err, ShouldBeNil)
		dt := tb.Topology()
		Reset(func() {
			dt.Stop()
		})

		sin, err := dt.Sink("snk")
		So(err, ShouldBeNil)
		si := sin.Sink().(*tupleCollectorSink)

		Convey("When 4 tuples are emitted by the source", func() {
			Convey("Then the sink receives a number of tuples", func() {
				si.Wait(2)
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldBeGreaterThanOrEqualTo, 2)

				// the number and order or result tuples varies,
				// so there is not a lot of stuff we can check...
				Convey("And all tuples should have keys a,b,c,d", func() {
					t := si.Tuples[0]
					// the first tuple should definitely have the same timestamp
					// as the first tuple in the input set
					So(t.Timestamp, ShouldResemble, tuples[0].Timestamp)

					for _, tup := range si.Tuples {
						_, hasA := tup.Data["a"]
						So(hasA, ShouldBeTrue)
						_, hasB := tup.Data["d"]
						So(hasB, ShouldBeTrue)
						_, hasC := tup.Data["c"]
						So(hasC, ShouldBeTrue)
						_, hasD := tup.Data["d"]
						So(hasD, ShouldBeTrue)
					}
				})
			})
		})
	})
}

func TestBQLBoxUDSF(t *testing.T) {
	Convey("Given a topology using UDSF", t, func() {
		tb, err := setupTopology(`CREATE STREAM box AS SELECT ISTREAM duplicate:int FROM duplicate('source', 3) [RANGE 1 TUPLES]`)
		So(err, ShouldBeNil)
		dt := tb.Topology()
		Reset(func() {
			dt.Stop()
		})

		sin, err := dt.Sink("snk")
		So(err, ShouldBeNil)
		si := sin.Sink().(*tupleCollectorSink)

		Convey("When 4 tuples are emitted by the source", func() {
			Convey("Then the sink should receive 12 tuples", func() {
				si.Wait(12)
				So(len(si.Tuples), ShouldEqual, 12)
			})
		})
	})
}

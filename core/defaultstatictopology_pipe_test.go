package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
	"time"
)

func freshTuples() []*tuple.Tuple {
	tup1 := &tuple.Tuple{
		Data: tuple.Map{
			"seq": tuple.Int(1),
		},
		InputName: "input",
	}
	tup2 := tup1.Copy()
	tup2.Data["seq"] = tuple.Int(2)
	tup3 := tup1.Copy()
	tup3.Data["seq"] = tuple.Int(3)
	tup4 := tup1.Copy()
	tup4.Data["seq"] = tuple.Int(4)
	tup5 := tup1.Copy()
	tup5.Data["seq"] = tuple.Int(5)
	tup6 := tup1.Copy()
	tup6.Data["seq"] = tuple.Int(6)
	tup7 := tup1.Copy()
	tup7.Data["seq"] = tuple.Int(7)
	tup8 := tup1.Copy()
	tup8.Data["seq"] = tuple.Int(8)
	return []*tuple.Tuple{tup1, tup2, tup3, tup4,
		tup5, tup6, tup7, tup8}
}

const (
	maxPar        = 1
	timeTolerance = 10 * time.Millisecond
	shortSleep    = 50 * time.Millisecond
	longSleep     = 150 * time.Millisecond
)

func TestCapacityPipeLinearTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a simple source/slow box/very slow box/sink topology", t, func() {
		/*
		 *   so -*--> b1 -*--> b2 -*--> si
		 */

		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: freshTuples(),
		}
		tb.AddSource("source", so)

		b1 := BoxFunc(slowForwardBox)
		tb.AddBox("box1", b1).Input("source")

		b2 := BoxFunc(verySlowForwardBox)
		tb.AddBox("box2", b2).Input("box1")

		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("box2")

		t, _ := tb.Build()

		for par := 1; par <= maxPar; par++ {
			par := par // safer to overlay the loop variable when used in closures
			Convey(fmt.Sprintf("When tuples are emitted with parallelism %d", par), func() {
				t.Run(ctx)

				// check that tuples arrived
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 8)

				// check that length of trace matches expectation
				So(len(si.Tuples[0].Trace), ShouldEqual, 6) // OUT-IN-OUT-IN-OUT-IN

				Convey("Then waiting time at intermediate pipes matches expectations", func() {
					// SOURCE'S PIPE
					{
						// the first `par` items should be emitted without delay
						// (because `par` items can be processed in parallel)
						for i := 0; i < par; i++ {
							waitTime := si.Tuples[i].Trace[1].Timestamp.Sub(si.Tuples[i].Trace[0].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
						// the next item should have to wait as long as the first
						// box processes (shortSleep)
						waitTime := si.Tuples[par].Trace[1].Timestamp.Sub(si.Tuples[par].Trace[0].Timestamp)
						So(waitTime, ShouldAlmostEqual, shortSleep, timeTolerance)
						// after that, every par'th item should have to wait
						// as long as the longest processing box needs (longSleep),
						// all others be processed immediately
						for i := par + 1; i < len(si.Tuples); i++ {
							waitTime := si.Tuples[i].Trace[1].Timestamp.Sub(si.Tuples[i].Trace[0].Timestamp)
							if (i-par)%par == 0 {
								So(waitTime, ShouldAlmostEqual, longSleep, timeTolerance)
							} else {
								So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
							}
						}
					}

					// FIRST BOX'S PIPE
					{
						// the first `par` items should be emitted without delay
						// (because `par` items can be processed in parallel)
						for i := 0; i < par; i++ {
							waitTime := si.Tuples[i].Trace[3].Timestamp.Sub(si.Tuples[i].Trace[2].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
						// after that, every par'th item should have to wait
						// as long as the next processing box needs, minus the
						// duration of this box itself,
						// all others be processed immediately
						for i := par; i < len(si.Tuples); i++ {
							waitTime := si.Tuples[i].Trace[3].Timestamp.Sub(si.Tuples[i].Trace[2].Timestamp)
							So(waitTime, ShouldAlmostEqual, longSleep-shortSleep, timeTolerance)
						}
					}

					// SECOND BOX'S PIPE
					{
						// the sink is much faster than the box, so waiting time should
						// 0 for all items
						for i := 0; i < len(si.Tuples); i++ {
							waitTime := si.Tuples[i].Trace[5].Timestamp.Sub(si.Tuples[i].Trace[4].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
					}
				})
			})
		}
	})

	Convey("Given a simple source/very slow box/slow box/sink topology", t, func() {
		/*
		 *   so -*--> b1 -*--> b2 -*--> si
		 */

		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: freshTuples(),
		}
		tb.AddSource("source", so)

		b1 := BoxFunc(verySlowForwardBox)
		tb.AddBox("box1", b1).Input("source")

		b2 := BoxFunc(slowForwardBox)
		tb.AddBox("box2", b2).Input("box1")

		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("box2")

		t, _ := tb.Build()

		for par := 1; par <= maxPar; par++ {
			par := par // safer to overlay the loop variable when used in closures
			Convey(fmt.Sprintf("When tuples are emitted with parallelism %d", par), func() {
				t.Run(ctx)

				// check that tuples arrived
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 8)

				// check that length of trace matches expectation
				So(len(si.Tuples[0].Trace), ShouldEqual, 6) // OUT-IN-OUT-IN-OUT-IN

				Convey("Then waiting time at intermediate pipes matches expectations", func() {
					// SOURCE'S PIPE
					{
						// the first `par` items should be emitted without delay
						// (because `par` items can be processed in parallel)
						for i := 0; i < par; i++ {
							waitTime := si.Tuples[i].Trace[1].Timestamp.Sub(si.Tuples[i].Trace[0].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
						// the next item should have to wait as long as the first
						// box processes (longSleep)
						waitTime := si.Tuples[par].Trace[1].Timestamp.Sub(si.Tuples[par].Trace[0].Timestamp)
						So(waitTime, ShouldAlmostEqual, longSleep, timeTolerance)
						// after that, every par'th item should have to wait
						// as long as the longest processing box needs (longSleep),
						// all others be processed immediately
						for i := par + 1; i < len(si.Tuples); i++ {
							waitTime := si.Tuples[i].Trace[1].Timestamp.Sub(si.Tuples[i].Trace[0].Timestamp)
							if (i-par)%par == 0 {
								So(waitTime, ShouldAlmostEqual, longSleep, timeTolerance)
							} else {
								So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
							}
						}
					}

					// FIRST BOX'S PIPE
					{
						// the second box is much faster, so waiting time should be
						// 0 for all items
						for i := 0; i < len(si.Tuples); i++ {
							waitTime := si.Tuples[i].Trace[3].Timestamp.Sub(si.Tuples[i].Trace[2].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
					}

					// SECOND BOX'S PIPE
					{
						// the sink is much faster than the box, so waiting time should
						// 0 for all items
						for i := 0; i < len(si.Tuples); i++ {
							waitTime := si.Tuples[i].Trace[5].Timestamp.Sub(si.Tuples[i].Trace[4].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
					}
				})
			})
		}
	})
}

func TestCapacityPipeForkTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)
	Convey("Given a simple source/box/sink topology with 2 boxes and 2 sinks", t, func() {
		/*
		 *        /--> b1 -*--> si1
		 *   so -*
		 *        \--> b2 -*--> si2
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: freshTuples(),
		}
		tb.AddSource("source", so)

		b1 := BoxFunc(slowForwardBox)
		tb.AddBox("box1", b1).Input("source")
		b2 := BoxFunc(verySlowForwardBox)
		tb.AddBox("box2", b2).Input("source")

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("box1")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("box2")

		t, _ := tb.Build()

		for par := 1; par <= maxPar; par++ {
			par := par // safer to overlay the loop variable when used in closures
			Convey(fmt.Sprintf("When tuples are emitted with parallelism %d", par), func() {
				t.Run(ctx)

				// check that tuples arrived
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 8)
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 8)

				// check that length of trace matches expectation
				So(len(si1.Tuples[0].Trace), ShouldEqual, 4) // OUT-IN-OUT-IN
				So(len(si2.Tuples[0].Trace), ShouldEqual, 4) // OUT-IN-OUT-IN

				Convey("Then waiting time at intermediate pipes should matches expectations", func() {
					// SOURCE'S PIPE
					{
						// the first `par` items should be emitted without delay
						// (because `par` items can be processed in parallel)
						for i := 0; i < par; i++ {
							b1waitTime := si1.Tuples[i].Trace[1].Timestamp.Sub(si2.Tuples[i].Trace[0].Timestamp)
							So(b1waitTime, ShouldAlmostEqual, 0, timeTolerance)
							b2waitTime := si1.Tuples[i].Trace[1].Timestamp.Sub(si2.Tuples[i].Trace[0].Timestamp)
							So(b2waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}

						// box1
						// the next item should have to wait as long as the first
						// box processes (shortSleep)
						waitTime := si1.Tuples[par].Trace[1].Timestamp.Sub(si1.Tuples[par].Trace[0].Timestamp)
						So(waitTime, ShouldAlmostEqual, shortSleep, timeTolerance)
						// after the first `par + 1` items, every item should
						// be processed immediately. (the source is waiting for
						// the slower box to finish before reading the next item,
						// at which point the faster box already has free capacity
						// again)
						for i := par + 1; i < len(so.Tuples); i++ {
							waitTime := si1.Tuples[i].Trace[1].Timestamp.Sub(si1.Tuples[i].Trace[0].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}

						// box2
						// after the first `par` items, every par'th item should
						// have to wait as long as the longest processing box
						// needs (longSleep), all others be processed immediately
						for i := par; i < len(so.Tuples); i++ {
							waitTime := si2.Tuples[i].Trace[1].Timestamp.Sub(si2.Tuples[i].Trace[0].Timestamp)
							if (i-par)%par == 0 {
								So(waitTime, ShouldAlmostEqual, longSleep, timeTolerance)
							} else {
								So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
							}
						}
					}

					// FIRST BOX'S PIPE
					{
						// the sink is much faster than the box, so waiting time should
						// 0 for all items
						for i := 0; i < len(si1.Tuples); i++ {
							waitTime := si1.Tuples[i].Trace[3].Timestamp.Sub(si1.Tuples[i].Trace[2].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
					}

					// SECOND BOX'S PIPE
					{
						// the sink is much faster than the box, so waiting time should
						// 0 for all items
						for i := 0; i < len(si2.Tuples); i++ {
							waitTime := si2.Tuples[i].Trace[3].Timestamp.Sub(si2.Tuples[i].Trace[2].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
					}
				})
			})
		}
	})
}

func TestCapacityPipeJoinTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)
	Convey("Given a simple source/box/sink topology with 2 sources", t, func() {
		/*
		 *   so1 -*-\
		 *           --> b -*--> si
		 *   so2 -*-/
		 */
		tb := NewDefaultStaticTopologyBuilder()

		so1 := &TupleEmitterSource{
			Tuples: freshTuples()[0:4],
		}
		tb.AddSource("source1", so1)

		so2 := &TupleEmitterSource{
			Tuples: freshTuples()[4:8],
		}
		tb.AddSource("source2", so2)

		b1 := BoxFunc(slowForwardBox)
		tb.AddBox("box", b1).
			Input("source1").
			Input("source2")
		si := &TupleCollectorSink{}
		tb.AddSink("si", si).Input("box")

		t, _ := tb.Build()

		for par := 1; par <= maxPar; par++ {
			par := par // safer to overlay the loop variable when used in closures
			Convey(fmt.Sprintf("When tuples are emitted with parallelism %d", par), func() {
				t.Run(ctx)

				// check that tuples arrived
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 8)

				// check that length of trace matches expectation
				So(len(si.Tuples[0].Trace), ShouldEqual, 4) // OUT-IN-OUT-IN

				Convey("Then waiting time at intermediate pipes should matches expectations", func() {
					// SOURCE 1'S PIPE
					{
						// the first `par` items should be emitted without delay
						// (because `par` items can be processed in parallel)
						for i := 0; i < par; i++ {
							waitTime := so1.Tuples[i].Trace[1].Timestamp.Sub(so1.Tuples[i].Trace[0].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
						// after that, every par'th item should have to wait
						// as long as the longest processing box needs (shortSleep),
						// all others be processed immediately
						for i := par; i < len(so1.Tuples); i++ {
							waitTime := so1.Tuples[i].Trace[1].Timestamp.Sub(so1.Tuples[i].Trace[0].Timestamp)
							if (i-par)%par == 0 {
								So(waitTime, ShouldAlmostEqual, shortSleep, timeTolerance)
							} else {
								So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
							}
						}
					}

					// SOURCE 2'S PIPE
					{
						// the first `par` items should be emitted without delay
						// (because `par` items can be processed in parallel)
						for i := 0; i < par; i++ {
							waitTime := so2.Tuples[i].Trace[1].Timestamp.Sub(so2.Tuples[i].Trace[0].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
						// after that, every par'th item should have to wait
						// as long as the longest processing box needs (shortSleep),
						// all others be processed immediately
						for i := par; i < len(so2.Tuples); i++ {
							waitTime := so2.Tuples[i].Trace[1].Timestamp.Sub(so2.Tuples[i].Trace[0].Timestamp)
							if (i-par)%par == 0 {
								So(waitTime, ShouldAlmostEqual, shortSleep, timeTolerance)
							} else {
								So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
							}
						}
					}

					// BOX'S PIPE
					{
						// the sink is much faster than the box, so waiting time should
						// 0 for all items
						for i := 0; i < len(si.Tuples); i++ {
							waitTime := si.Tuples[i].Trace[3].Timestamp.Sub(si.Tuples[i].Trace[2].Timestamp)
							So(waitTime, ShouldAlmostEqual, 0, timeTolerance)
						}
					}
				})
			})
		}
	})
}

func slowForwardBox(ctx *Context, t *tuple.Tuple, w Writer) error {
	time.Sleep(shortSleep)
	w.Write(ctx, t)
	return nil
}

func verySlowForwardBox(ctx *Context, t *tuple.Tuple, w Writer) error {
	time.Sleep(longSleep)
	w.Write(ctx, t)
	return nil
}

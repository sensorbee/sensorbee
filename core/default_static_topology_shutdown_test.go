package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// On one topology, there're some patterns to be tested.
//
// 1. Call Stop when no tuple was generated from Source.
// 2. Call Stop when some but not all tuples are generated from Source
//   a. and no tuples are received by Sink.
//   b. some tuples are received by Sink and others are in Boxes.
//   c. and all those tuples are received by Sink.
// 3. Call Stop when all tuples are generated from Source
//   a. and no tuples are received by Sink.
//   b. and some tuples are received by Sink and others are in Boxes.
//   c. and all tuples are received by Sink.
//
// Followings are topologies to be tested:
//
// 1. Linear
// 2. Multiple sources (including JOIN)
// 3. Multiple sinks (including fan-out)

func TestShutdownLinearDefaultStaticTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a simple linear topology", t, func() {
		/*
		 *   so -*--> b1 -*--> b2 -*--> si
		 */

		tb := NewDefaultStaticTopologyBuilder()
		so := NewTupleIncrementalEmitterSource(freshTuples())
		So(tb.AddSource("source", so).Err(), ShouldBeNil)

		b1 := &BlockingForwardBox{cnt: 8}
		So(tb.AddBox("box1", b1).Input("source").Err(), ShouldBeNil)

		b2 := BoxFunc(forwardBox)
		So(tb.AddBox("box2", b2).Input("box1").Err(), ShouldBeNil)

		si := NewTupleCollectorSink()
		So(tb.AddSink("sink", si).Input("box2").Err(), ShouldBeNil)

		ti, err := tb.Build()
		So(err, ShouldBeNil)

		t := ti.(*defaultStaticTopology)

		Convey("When generating no tuples and call stop", func() { // 1.
			go func() {
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)

			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then the sink shouldn't receive anything", func() {
				So(si.Tuples, ShouldBeEmpty)
			})
		})

		Convey("When generating some tuples and call stop before the sink receives a tuple", func() { // 2.a.
			b1.cnt = 0 // Block tuples. The box isn't running yet, so cnt can safely be changed.
			go func() {
				so.EmitTuplesNB(4)
				go func() {
					t.Stop(ctx)
				}()

				// resume b1 after the topology starts stopping all.
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(8)
			}()
			So(t.Run(ctx), ShouldBeNil)

			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then the sink should receive all of generated tuples", func() {
				So(len(si.Tuples), ShouldEqual, 4)
			})
		})

		Convey("When generating some tuples and call stop after the sink received some of them", func() { // 2.b.
			go func() {
				so.EmitTuplesNB(3)
				t.Wait(ctx, TSRunning)
				b1.EmitTuples(1)
				si.Wait(1)
				go func() {
					t.Stop(ctx)
				}()

				t.Wait(ctx, TSStopping)
				b1.EmitTuples(2)
			}()
			So(t.Run(ctx), ShouldBeNil)

			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then the sink should receive all of those tuples", func() {
				So(len(si.Tuples), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after the sink received all", func() { // 2.c.
			go func() {
				so.EmitTuples(4)
				si.Wait(4)
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)

			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then the sink should only receive those tuples", func() {
				So(len(si.Tuples), ShouldEqual, 4)
			})
		})

		Convey("When generating all tuples and call stop before the sink receives a tuple", func() { // 3.a.
			b1.cnt = 0
			go func() {
				so.EmitTuples(100) // Blocking call. Assuming the pipe's capacity is greater than or equal to 8.
				go func() {
					t.Stop(ctx)
				}()
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(8)
			}()
			So(t.Run(ctx), ShouldBeNil)

			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received some of them", func() { // 3.b.
			b1.cnt = 2
			go func() {
				so.EmitTuples(100)
				si.Wait(2)
				go func() {
					t.Stop(ctx)
				}()
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(6)
			}()
			So(t.Run(ctx), ShouldBeNil)

			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then the sink should receive all of those tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received all", func() { // 3.c.
			go func() {
				so.EmitTuples(100)
				si.Wait(8)
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)

			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})
	})
}

func TestShutdownForkTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)
	SkipConvey("Given a simple source/box/sink topology with 2 boxes and 2 sinks", t, func() {
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

		t, err := tb.Build()
		So(err, ShouldBeNil)

		for par := 1; par <= maxPar; par++ {
			par := par // safer to overlay the loop variable when used in closures
			Convey(fmt.Sprintf("When tuples are emitted with parallelism %d", par), func() {
				go func() {
					// we should be able to emit two tuples before stopping,
					// but time is not enough to process both of them. the call
					// to Stop() should wait for processing to complete, though.
					time.Sleep(time.Duration(0.5 * float64(shortSleep)))
					t.Stop(ctx)
				}()
				t.Run(ctx)

				// check that tuples arrived
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 2)
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 2)

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

						// box2
						// after the first `par` items, every par'th item should
						// have to wait as long as the longest processing box
						// needs (longSleep), all others be processed immediately
						for i := par; i < len(si2.Tuples); i++ {
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

func TestShutdownJoinTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)
	SkipConvey("Given a simple source/box/sink topology with 2 sources", t, func() {
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

		t, err := tb.Build()
		So(err, ShouldBeNil)

		for par := 1; par <= maxPar; par++ {
			par := par // safer to overlay the loop variable when used in closures
			Convey(fmt.Sprintf("When tuples are emitted with parallelism %d", par), func() {
				go func() {
					// we should be able to emit two tuples each before stopping,
					// but time is not enough to process both of them. the call
					// to Stop() should wait for processing to complete, though.
					time.Sleep(time.Duration(0.4 * float64(shortSleep)))
					t.Stop(ctx)
				}()
				t.Run(ctx)

				// check that tuples arrived
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 4)

				// check that length of trace matches expectation
				So(len(si.Tuples[0].Trace), ShouldEqual, 4) // OUT-IN-OUT-IN

				Convey("Then waiting time at intermediate pipes matches expectations", func() {
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
						for i := par; i < 2; i++ {
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
						for i := par; i < 2; i++ {
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

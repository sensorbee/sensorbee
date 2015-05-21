package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
	"testing"
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

type terminateChecker struct {
	ProxyBox

	m            sync.Mutex
	c            *sync.Cond
	terminateCnt int
}

func NewTerminateChecker(b Box) *terminateChecker {
	t := &terminateChecker{
		ProxyBox: ProxyBox{b: b},
	}
	t.c = sync.NewCond(&t.m)
	return t
}

func (t *terminateChecker) Terminate(ctx *Context) error {
	t.m.Lock()
	t.terminateCnt++
	t.c.Broadcast()
	t.m.Unlock()
	return t.ProxyBox.Terminate(ctx)
}

func (t *terminateChecker) WaitForTermination() {
	t.m.Lock()
	defer t.m.Unlock()
	for t.terminateCnt == 0 {
		t.c.Wait()
	}
}

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
// 2. Multiple sinks (including fan-out)
// 3. Multiple sources (including JOIN)

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
		tc1 := NewTerminateChecker(b1)
		So(tb.AddBox("box1", tc1).Input("source").Err(), ShouldBeNil)

		b2 := BoxFunc(forwardBox)
		tc2 := NewTerminateChecker(b2)
		So(tb.AddBox("box2", tc2).Input("box1").Err(), ShouldBeNil)

		si := NewTupleCollectorSink()
		So(tb.AddSink("sink", si).Input("box2").Err(), ShouldBeNil)

		ti, err := tb.Build()
		So(err, ShouldBeNil)

		t := ti.(*defaultStaticTopology)

		checkPostCond := func() {
			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then Box.Terminate should be called exactly once", func() {
				So(tc1.terminateCnt, ShouldEqual, 1)
				So(tc2.terminateCnt, ShouldEqual, 1)
			})
		}

		Convey("When generating no tuples and call stop", func() { // 1.
			go func() {
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

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
			checkPostCond()

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
			checkPostCond()

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
			checkPostCond()

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
			checkPostCond()

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
			checkPostCond()

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
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})
	})
}

func TestShutdownForkDefaultStaticTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a simple fork topology", t, func() {
		/*
		 *        /--> b1 -*--> si1
		 *   so -*
		 *        \--> b2 -*--> si2
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so := NewTupleIncrementalEmitterSource(freshTuples())
		So(tb.AddSource("source", so).Err(), ShouldBeNil)

		b1 := &BlockingForwardBox{cnt: 8}
		tc1 := NewTerminateChecker(b1)
		So(tb.AddBox("box1", tc1).Input("source").Err(), ShouldBeNil)
		b2 := &BlockingForwardBox{cnt: 8}
		tc2 := NewTerminateChecker(b2)
		So(tb.AddBox("box2", tc2).Input("source").Err(), ShouldBeNil)

		si1 := NewTupleCollectorSink()
		So(tb.AddSink("si1", si1).Input("box1").Err(), ShouldBeNil)
		si2 := NewTupleCollectorSink()
		So(tb.AddSink("si2", si2).Input("box2").Err(), ShouldBeNil)

		ti, err := tb.Build()
		So(err, ShouldBeNil)

		t := ti.(*defaultStaticTopology)

		checkPostCond := func() {
			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then Box.Terminate should be called exactly once", func() {
				So(tc1.terminateCnt, ShouldEqual, 1)
				So(tc2.terminateCnt, ShouldEqual, 1)
			})
		}

		Convey("When generating no tuples and call stop", func() { // 1.
			go func() {
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink shouldn't receive anything", func() {
				So(si1.Tuples, ShouldBeEmpty)
				So(si2.Tuples, ShouldBeEmpty)
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
			checkPostCond()

			Convey("Then the sink should receive all of generated tuples", func() {
				So(len(si1.Tuples), ShouldEqual, 4)
				So(len(si2.Tuples), ShouldEqual, 4)
			})
		})

		Convey("When generating some tuples and call stop after the sink received some of them", func() { // 2.b.
			go func() {
				so.EmitTuplesNB(3)
				t.Wait(ctx, TSRunning)
				b1.EmitTuples(1)
				si1.Wait(1)
				go func() {
					t.Stop(ctx)
				}()

				t.Wait(ctx, TSStopping)
				b1.EmitTuples(2)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(len(si1.Tuples), ShouldEqual, 3)
				So(len(si2.Tuples), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after both sinks received all", func() { // 2.c.
			go func() {
				so.EmitTuples(4)
				si1.Wait(4)
				si2.Wait(4)
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should only receive those tuples", func() {
				So(len(si1.Tuples), ShouldEqual, 4)
				So(len(si2.Tuples), ShouldEqual, 4)
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
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(len(si1.Tuples), ShouldEqual, 8)
				So(len(si2.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received some of them", func() { // 3.b.
			b1.cnt = 2
			go func() {
				so.EmitTuples(100)
				si1.Wait(2)
				// don't care about si2
				go func() {
					t.Stop(ctx)
				}()
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(6)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(len(si1.Tuples), ShouldEqual, 8)
				So(len(si2.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the two sinks received some of them", func() { // 3.b'.
			b1.cnt = 2
			b2.cnt = 3
			go func() {
				so.EmitTuples(100)
				si1.Wait(2)
				si2.Wait(3)
				go func() {
					t.Stop(ctx)
				}()
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(6)
				b2.EmitTuples(5)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(len(si1.Tuples), ShouldEqual, 8)
				So(len(si2.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received all", func() { // 3.c.
			go func() {
				so.EmitTuples(100)
				si1.Wait(8)
				si2.Wait(8)
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(len(si1.Tuples), ShouldEqual, 8)
				So(len(si2.Tuples), ShouldEqual, 8)
			})
		})
	})
}

func TestShutdownJoinDefaultStaticTopology(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a simple join topology", t, func() {
		/*
		 *   so1 -*-\
		 *           --> b -*--> si
		 *   so2 -*-/
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so1 := NewTupleIncrementalEmitterSource(freshTuples()[0:4])
		So(tb.AddSource("source1", so1).Err(), ShouldBeNil)
		so2 := NewTupleIncrementalEmitterSource(freshTuples()[4:8])
		So(tb.AddSource("source2", so2).Err(), ShouldBeNil)

		b1 := &BlockingForwardBox{cnt: 8}
		tc1 := NewTerminateChecker(b1)
		So(tb.AddBox("box1", tc1).Input("source1").Input("source2").Err(), ShouldBeNil)

		si := NewTupleCollectorSink()
		So(tb.AddSink("sink", si).Input("box1").Err(), ShouldBeNil)

		ti, err := tb.Build()
		So(err, ShouldBeNil)

		t := ti.(*defaultStaticTopology)

		checkPostCond := func() {
			Convey("Then the topology should be stopped", func() {
				So(t.State(ctx), ShouldEqual, TSStopped)
			})

			Convey("Then Box.Terminate should be called exactly once", func() {
				So(tc1.terminateCnt, ShouldEqual, 1)
			})
		}

		Convey("When generating no tuples and call stop", func() { // 1.
			go func() {
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink shouldn't receive anything", func() {
				So(si.Tuples, ShouldBeEmpty)
			})
		})

		Convey("When generating some tuples and call stop before the sink receives a tuple", func() { // 2.a.
			b1.cnt = 0 // Block tuples. The box isn't running yet, so cnt can safely be changed.
			go func() {
				so1.EmitTuplesNB(3)
				so2.EmitTuplesNB(1)
				go func() {
					t.Stop(ctx)
				}()

				// resume b1 after the topology starts stopping all.
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(8)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all of generated tuples", func() {
				So(len(si.Tuples), ShouldEqual, 4)
			})
		})

		Convey("When generating some tuples from one source and call stop before the sink receives a tuple", func() { // 2.a'.
			b1.cnt = 0 // Block tuples. The box isn't running yet, so cnt can safely be changed.
			go func() {
				so1.EmitTuplesNB(3)
				go func() {
					t.Stop(ctx)
				}()

				// resume b1 after the topology starts stopping all.
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(8)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all of generated tuples", func() {
				So(len(si.Tuples), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after the sink received some of them", func() { // 2.b.
			go func() {
				so1.EmitTuplesNB(1)
				so2.EmitTuplesNB(2)
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
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(len(si.Tuples), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after the sink received all", func() { // 2.c.
			go func() {
				so1.EmitTuples(2)
				so2.EmitTuples(1)
				si.Wait(3)
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should only receive those tuples", func() {
				So(len(si.Tuples), ShouldEqual, 3)
			})
		})

		Convey("When generating all tuples and call stop before the sink receives a tuple", func() { // 3.a.
			b1.cnt = 0
			go func() {
				so1.EmitTuples(100) // Blocking call. Assuming the pipe's capacity is greater than or equal to 8.
				so2.EmitTuples(100)
				go func() {
					t.Stop(ctx)
				}()
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(8)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received some of them", func() { // 3.b.
			b1.cnt = 2
			go func() {
				so1.EmitTuples(100)
				so2.EmitTuples(100)
				si.Wait(2)
				go func() {
					t.Stop(ctx)
				}()
				t.Wait(ctx, TSStopping)
				b1.EmitTuples(6)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received all", func() { // 3.c.
			go func() {
				so1.EmitTuples(100)
				so2.EmitTuples(100)
				si.Wait(8)
				t.Stop(ctx)
			}()
			So(t.Run(ctx), ShouldBeNil)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})
	})
}

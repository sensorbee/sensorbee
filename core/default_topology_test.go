package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func freshTuples() []*Tuple {
	tup1 := &Tuple{
		Data: data.Map{
			"seq": data.Int(1),
		},
		InputName: "input",
	}
	tup2 := tup1.Copy()
	tup2.Data["seq"] = data.Int(2)
	tup3 := tup1.Copy()
	tup3.Data["seq"] = data.Int(3)
	tup4 := tup1.Copy()
	tup4.Data["seq"] = data.Int(4)
	tup5 := tup1.Copy()
	tup5.Data["seq"] = data.Int(5)
	tup6 := tup1.Copy()
	tup6.Data["seq"] = data.Int(6)
	tup7 := tup1.Copy()
	tup7.Data["seq"] = data.Int(7)
	tup8 := tup1.Copy()
	tup8.Data["seq"] = data.Int(8)
	return []*Tuple{tup1, tup2, tup3, tup4,
		tup5, tup6, tup7, tup8}
}

type terminateChecker struct {
	ProxyBox

	m            sync.Mutex
	c            *sync.Cond
	terminateCnt int
}

func newTerminateChecker(b Box) *terminateChecker {
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

type sinkCloseChecker struct {
	s        Sink
	closeCnt int32
}

func (s *sinkCloseChecker) Write(ctx *Context, t *Tuple) error {
	return s.s.Write(ctx, t)
}

func (s *sinkCloseChecker) Close(ctx *Context) error {
	atomic.AddInt32(&s.closeCnt, 1)
	return s.s.Close(ctx)
}

func TestDefaultTopologySetup(t *testing.T) {
	Convey("Given a default topology", t, func() {
		dt := NewDefaultTopology(NewContext(nil), "dt1")
		t := dt.(*defaultTopology)
		Reset(func() {
			t.Stop()
		})

		dupNameTests := func(name string) {
			Convey("Then adding a source having the same name should fail", func() {
				_, err := t.AddSource(name, &DoesNothingSource{}, nil)
				So(err, ShouldNotBeNil)
			})

			Convey("Then adding a box having the same name should fail", func() {
				_, err := t.AddBox(name, &DoesNothingBox{}, nil)
				So(err, ShouldNotBeNil)
			})

			Convey("Then adding a sink having the same name should fail", func() {
				si := NewTupleCollectorSink()
				sic := &sinkCloseChecker{s: si}
				_, err := t.AddSink(name, sic, nil)
				So(err, ShouldNotBeNil)

				Convey("And the sink shoud be closed", func() {
					So(sic.closeCnt, ShouldEqual, 1)
				})
			})
		}

		Convey("When stopping it without adding anything", func() {
			t.Stop()

			Convey("Then it should stop", func() {
				So(t.state.state, ShouldEqual, TSStopped)
			})

			Convey("Then adding a source to the stopped topology should fail", func() {
				_, err := t.AddSource("test_source", &DoesNothingSource{}, nil)
				So(err, ShouldNotBeNil)
			})

			Convey("Then adding a box to the stopped topology should fail", func() {
				_, err := t.AddBox("test_box", &DoesNothingBox{}, nil)
				So(err, ShouldNotBeNil)
			})

			Convey("Then adding a sink to the stopped topology should fail", func() {
				_, err := t.AddSink("test_sink", &DoesNothingSink{}, nil)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When adding a source", func() {
			s := NewTupleIncrementalEmitterSource(freshTuples())
			sn, err := t.AddSource("source1", s, nil)
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)

				Convey("And stopping it again shouldn't fail", func() {
					So(sn.Stop(), ShouldBeNil)
				})

				Convey("And it shouldn't be removed", func() {
					_, err := t.Node("source1")
					So(err, ShouldBeNil)
				})
			})

			Convey("Then the topology should have it", func() {
				n, err := t.Source("SOURCE1")
				So(err, ShouldBeNil)
				sn1 := n.(*defaultSourceNode)
				sn2 := sn.(*defaultSourceNode)
				So(sn1, ShouldEqual, sn2)
			})

			Convey("Then it can be obtained as a node", func() {
				n, err := t.Node("SOURCE1")
				So(err, ShouldBeNil)
				sn1 := n.(*defaultSourceNode)
				sn2 := sn.(*defaultSourceNode)
				So(sn1, ShouldEqual, sn2)
			})

			Convey("Then calling Resume on the running source shouldn't fail", func() {
				So(sn.Resume(), ShouldBeNil)
			})

			Convey("And stopping it", func() {
				So(sn.Stop(), ShouldBeNil)

				Convey("Then Pause should fail", func() {
					So(sn.Pause(), ShouldNotBeNil)
				})

				Convey("Then Resume should fail", func() {
					So(sn.Resume(), ShouldNotBeNil)
				})
			})
			dupNameTests("source1")
		})

		Convey("When adding a paused source", func() {
			so := NewTupleEmitterSource(freshTuples())
			son, err := t.AddSource("source1", so, &SourceConfig{
				PausedOnStartup: true,
			})
			So(err, ShouldBeNil)

			si := NewTupleCollectorSink()
			sin, err := t.AddSink("sink", si, nil)
			So(err, ShouldBeNil)
			So(sin.Input("source1", nil), ShouldBeNil)

			Convey("Then the source's state should be paused", func() {
				So(son.State().Get(), ShouldEqual, TSPaused)

				Convey("And a redundant Pause call shouldn't fail", func() {
					So(son.Pause(), ShouldBeNil)
				})
			})

			Convey("Then the sink should've received nothing", func() {
				So(t.Stop(), ShouldBeNil)
				So(si.len(), ShouldEqual, 0)
			})

			Convey("Then the sink should receive all tuples by resuming the source", func() {
				So(son.Resume(), ShouldBeNil)
				si.Wait(8)
				So(t.Stop(), ShouldBeNil)
				So(si.len(), ShouldEqual, 8)
			})
		})

		Convey("When adding an auto remove source", func() {
			s := NewTupleIncrementalEmitterSource(freshTuples())
			sn, err := t.AddSource("source1", s, &SourceConfig{
				RemoveOnStop: true,
			})
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("source1")
						if err != nil {
							So(err, ShouldNotBeNil)
							break
						}
					}
				})
			})
		})

		Convey("When adding a source and call RemoveOnStop before stop", func() {
			s := NewTupleIncrementalEmitterSource(freshTuples())
			sn, err := t.AddSource("source1", s, &SourceConfig{
				RemoveOnStop: false,
			})
			So(err, ShouldBeNil)
			sn.RemoveOnStop()

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("source1")
						if err != nil {
							So(err, ShouldNotBeNil)
							break
						}
					}
				})
			})
		})

		Convey("When adding a source and call RemoveOnStop after stop", func() {
			s := NewTupleIncrementalEmitterSource(freshTuples())
			sn, err := t.AddSource("source1", s, &SourceConfig{
				RemoveOnStop: false,
			})
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)
				sn.RemoveOnStop()

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("source1")
						if err != nil {
							So(err, ShouldNotBeNil)
							break
						}
					}
				})
			})
		})

		Convey("When adding a box", func() {
			b := newTerminateChecker(&DoesNothingBox{})
			bn, err := t.AddBox("box1", b, nil)
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(bn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(bn.Stop(), ShouldBeNil)
				So(bn.State().Get(), ShouldEqual, TSStopped)

				Convey("And it should be terminated", func() {
					So(b.terminateCnt, ShouldEqual, 1)
				})

				Convey("And stopping it again shouldn't fail", func() {
					So(bn.Stop(), ShouldBeNil)
					So(b.terminateCnt, ShouldEqual, 1)
				})

				Convey("And it shouldn't be removed", func() {
					_, err := t.Node("box1")
					So(err, ShouldBeNil)
				})
			})

			Convey("Then Terminate should be called after stopping the topology", func() {
				So(t.Stop(), ShouldBeNil)
				So(b.terminateCnt, ShouldEqual, 1)
			})

			Convey("Then the topology should have it", func() {
				n, err := t.Box("BOX1")
				So(err, ShouldBeNil)
				So(n, ShouldPointTo, bn)
			})

			Convey("Then it can be obtained as a node", func() {
				n, err := t.Node("BOX1")
				So(err, ShouldBeNil)
				So(n, ShouldPointTo, bn)
			})

			dupNameTests("box1")
		})

		Convey("When adding an auto remove box", func() {
			b := newTerminateChecker(&DoesNothingBox{})
			bn, err := t.AddBox("box1", b, &BoxConfig{
				RemoveOnStop: true,
			})
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(bn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(bn.Stop(), ShouldBeNil)
				So(bn.State().Get(), ShouldEqual, TSStopped)

				Convey("And it should be terminated", func() {
					So(b.terminateCnt, ShouldEqual, 1)
				})

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("box1")
						if err != nil {
							So(err, ShouldNotBeNil)
							break
						}
					}
				})
			})
		})

		Convey("When adding a box and call RemoveOnStop before stop", func() {
			b := newTerminateChecker(&DoesNothingBox{})
			bn, err := t.AddBox("box1", b, &BoxConfig{
				RemoveOnStop: false,
			})
			So(err, ShouldBeNil)
			bn.RemoveOnStop()

			Convey("Then it should automatically run", func() {
				So(bn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(bn.Stop(), ShouldBeNil)
				So(bn.State().Get(), ShouldEqual, TSStopped)

				Convey("And it should be terminated", func() {
					So(b.terminateCnt, ShouldEqual, 1)
				})

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("box1")
						if err != nil {
							So(err, ShouldNotBeNil)
							break
						}
					}
				})
			})
		})

		Convey("When adding a box and call RemoveOnStop after stop", func() {
			b := newTerminateChecker(&DoesNothingBox{})
			bn, err := t.AddBox("box1", b, &BoxConfig{
				RemoveOnStop: false,
			})
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(bn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(bn.Stop(), ShouldBeNil)
				So(bn.State().Get(), ShouldEqual, TSStopped)
				bn.RemoveOnStop()

				Convey("And it should be terminated", func() {
					So(b.terminateCnt, ShouldEqual, 1)
				})

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("box1")
						if err != nil {
							So(err, ShouldNotBeNil)
							break
						}
					}
				})
			})
		})

		Convey("When adding a sink", func() {
			s := &DoesNothingSink{}
			sn, err := t.AddSink("sink1", s, nil)
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then the topology should have it", func() {
				n, err := t.Sink("SINK1")
				So(err, ShouldBeNil)
				So(n, ShouldPointTo, sn)
			})

			Convey("Then it can be obtained as a node", func() {
				n, err := t.Node("SINK1")
				So(err, ShouldBeNil)
				So(n, ShouldPointTo, sn)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)

				Convey("And stopping it again shouldn't fail", func() {
					So(sn.Stop(), ShouldBeNil)
				})

				Convey("And it shouldn't be removed", func() {
					_, err := t.Node("sink1")
					So(err, ShouldBeNil)
				})
			})

			dupNameTests("sink1")
		})

		Convey("When adding an auto remove sink", func() {
			s := &DoesNothingSink{}
			sn, err := t.AddSink("sink1", s, &SinkConfig{
				RemoveOnStop: true,
			})
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("sink1")
						if err != nil {
							So(IsNotExist(err), ShouldBeTrue)
							break
						}
					}
				})
			})
		})

		Convey("When adding a sink and call RemoveOnStop before stop", func() {
			s := &DoesNothingSink{}
			sn, err := t.AddSink("sink1", s, &SinkConfig{
				RemoveOnStop: false,
			})
			So(err, ShouldBeNil)
			sn.RemoveOnStop()

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("sink1")
						if err != nil {
							So(IsNotExist(err), ShouldBeTrue)
							break
						}
					}
				})
			})
		})

		Convey("When adding a sink and call RemoveOnStop after stop", func() {
			s := &DoesNothingSink{}
			sn, err := t.AddSink("sink1", s, &SinkConfig{
				RemoveOnStop: false,
			})
			So(err, ShouldBeNil)

			Convey("Then it should automatically run", func() {
				So(sn.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then it should be able to stop", func() {
				So(sn.Stop(), ShouldBeNil)
				So(sn.State().Get(), ShouldEqual, TSStopped)
				sn.RemoveOnStop()

				Convey("And it should be removed", func() {
					for {
						_, err := t.Node("sink1")
						if err != nil {
							So(IsNotExist(err), ShouldBeTrue)
							break
						}
					}
				})
			})
		})

		Convey("When getting nonexistent node", func() {
			_, err := t.Node("source1")

			Convey("Then it shouldn't be found", func() {
				So(IsNotExist(err), ShouldBeTrue)
			})
		})
	})
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

func TestLinearDefaultTopology(t *testing.T) {
	Convey("Given a simple linear topology", t, func() {
		/*
		 *   so -*--> b1 -*--> b2 -*--> si
		 */
		dt := NewDefaultTopology(NewContext(nil), "dt1")
		t := dt.(*defaultTopology)
		Reset(func() {
			t.Stop()
		})

		so := NewTupleIncrementalEmitterSource(freshTuples())
		son, err := t.AddSource("source", so, nil)
		So(err, ShouldBeNil)

		b1 := &BlockingForwardBox{cnt: 8}
		tc1 := newTerminateChecker(b1)
		bn1, err := t.AddBox("box1", tc1, nil)
		So(err, ShouldBeNil)
		So(bn1.Input("SOURCE", nil), ShouldBeNil)

		b2 := BoxFunc(forwardBox)
		tc2 := newTerminateChecker(b2)
		bn2, err := t.AddBox("box2", tc2, nil)
		So(err, ShouldBeNil)
		So(bn2.Input("BOX1", nil), ShouldBeNil)

		si := NewTupleCollectorSink()
		sic := &sinkCloseChecker{s: si}
		sin, err := t.AddSink("sink", sic, nil)
		So(err, ShouldBeNil)
		So(sin.Input("BOX2", nil), ShouldBeNil)

		checkPostCond := func() {
			Convey("Then the topology should be stopped", func() {
				So(t.state.Get(), ShouldEqual, TSStopped)
			})

			Convey("Then Box.Terminate should be called exactly once", func() {
				So(tc1.terminateCnt, ShouldEqual, 1)
				So(tc2.terminateCnt, ShouldEqual, 1)
			})

			Convey("Then Sink.Close should be called exactly once", func() {
				So(sic.closeCnt, ShouldEqual, 1)
			})
		}

		Convey("When getting registered nodes", func() {
			Reset(func() {
				t.Stop()
			})

			Convey("Then the topology should return all nodes", func() {
				ns := t.Nodes()
				So(len(ns), ShouldEqual, 4)
				dson1 := ns["source"].(*defaultSourceNode)
				dson2 := son.(*defaultSourceNode)
				So(dson1, ShouldEqual, dson2)
				So(ns["box1"], ShouldPointTo, bn1)
				So(ns["box2"], ShouldPointTo, bn2)
				So(ns["sink"], ShouldPointTo, sin)
			})

			Convey("Then source should be able to be obtained", func() {
				s, err := t.Source("source")
				So(err, ShouldBeNil)
				dson1 := s.(*defaultSourceNode)
				dson2 := son.(*defaultSourceNode)
				So(dson1, ShouldEqual, dson2)
			})

			Convey("Then source should be able to be obtained through Sources", func() {
				ss := t.Sources()
				So(len(ss), ShouldEqual, 1)
				dson1 := ss["source"].(*defaultSourceNode)
				dson2 := son.(*defaultSourceNode)
				So(dson1, ShouldEqual, dson2)
			})

			Convey("Then box1 should be able to be obtained", func() {
				b, err := t.Box("box1")
				So(err, ShouldBeNil)
				So(b, ShouldPointTo, bn1)
			})

			Convey("Then box2 should be able to be obtained", func() {
				b, err := t.Box("box2")
				So(err, ShouldBeNil)
				So(b, ShouldPointTo, bn2)
			})

			Convey("Then all boxes should be able to be obtained at once", func() {
				bs := t.Boxes()
				So(len(bs), ShouldEqual, 2)
				So(bs["box1"], ShouldPointTo, bn1)
				So(bs["box2"], ShouldPointTo, bn2)
			})

			Convey("Then sink should be able to be obtained", func() {
				s, err := t.Sink("sink")
				So(err, ShouldBeNil)
				So(s, ShouldPointTo, sin)
			})

			Convey("Then sink should be able to be obtained through Sinks", func() {
				ss := t.Sinks()
				So(len(ss), ShouldEqual, 1)
				So(ss["sink"], ShouldPointTo, sin)
			})

			Convey("Then source cannot be obtained via wronge methods", func() {
				_, err := t.Box("source")
				So(err, ShouldNotBeNil)
				_, err = t.Sink("source")
				So(err, ShouldNotBeNil)
			})

			Convey("Then box1 cannot be obtained via wronge methods", func() {
				_, err := t.Source("box1")
				So(err, ShouldNotBeNil)
				_, err = t.Sink("box1")
				So(err, ShouldNotBeNil)
			})

			Convey("Then sink cannot be obtained via wronge methods", func() {
				_, err := t.Source("sink")
				So(err, ShouldNotBeNil)
				_, err = t.Box("sink")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When generating no tuples and call stop", func() {
			So(t.Stop(), ShouldBeNil)
			checkPostCond()
		})

		Convey("When generating some tuples and call stop before the sink receives a tuple", func() { // 2.a.
			b1.cnt = 0
			so.EmitTuplesNB(4)
			go func() {
				t.Stop()
			}()
			t.state.Wait(TSStopping)
			b1.EmitTuples(8)
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of generated tuples", func() {
				So(si.len(), ShouldEqual, 4)
			})
		})

		Convey("When generating some tuples and call stop after the sink received some of them", func() { // 2.b.
			b1.cnt = 0
			go func() {
				so.EmitTuplesNB(3)
				b1.EmitTuples(1)
				si.Wait(1)
				go func() {
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(2)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(si.len(), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after the sink received all", func() { // 2.c.
			go func() {
				so.EmitTuples(4)
				si.Wait(4)
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should only receive those tuples", func() {
				So(si.len(), ShouldEqual, 4)
			})
		})

		Convey("When generating all tuples and call stop before the sink receives a tuple", func() { // 3.a.
			b1.cnt = 0
			go func() {
				so.EmitTuples(100) // Blocking call. Assuming the pipe's capacity is greater than or equal to 8.
				go func() {
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(8)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(si.len(), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received some of them", func() { // 3.b.
			b1.cnt = 2
			go func() {
				so.EmitTuples(100)
				si.Wait(2)
				go func() {
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(6)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(si.len(), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received all", func() { // 3.c.
			go func() {
				so.EmitTuples(100)
				si.Wait(8)
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(si.len(), ShouldEqual, 8)
			})
		})

		Convey("When removing a nonexistent node", func() {
			Convey("Then it shouldn't fail", func() {
				So(IsNotExist(t.Remove("no_such_node")), ShouldBeTrue)
			})
		})

		Convey("When removing the source after generating some tuples", func() {
			Reset(func() {
				t.Stop()
			})
			so.EmitTuples(2)
			So(t.Remove("SOURCE"), ShouldBeNil)

			Convey("Then the source should be stopped", func() {
				So(son.State().Get(), ShouldEqual, TSStopped)
			})

			Convey("Then the source shouldn't be found", func() {
				_, err := t.Source("source")
				So(err, ShouldNotBeNil)
			})

			Convey("Then the sink should receive the tuples", func() {
				si.Wait(2)
				So(si.len(), ShouldEqual, 2)
			})
		})

		Convey("When removing a box after processing some tuples", func() {
			Reset(func() {
				t.Stop()
			})
			so.EmitTuples(2)
			si.Wait(2)
			So(t.Remove("BOX1"), ShouldBeNil)
			so.EmitTuples(2)

			Convey("Then the box should be stopped", func() {
				So(bn1.State().Get(), ShouldEqual, TSStopped)
			})

			Convey("Then the box shouldn't be found", func() {
				_, err := t.Box("box1")
				So(err, ShouldNotBeNil)
			})

			Convey("Then sink should only receive the processed tuples", func() {
				So(si.len(), ShouldEqual, 2)
			})

			Convey("And connecting another box in the topology", func() {
				b3 := BoxFunc(forwardBox)
				bn3, err := t.AddBox("box3", b3, nil)
				So(err, ShouldBeNil)
				So(bn3.Input("source", nil), ShouldBeNil)
				So(bn2.Input("box3", nil), ShouldBeNil)
				so.EmitTuples(4)

				Convey("Then the sink should receive the correct number of tuples", func() {
					si.Wait(6) // 2 tuples which send just after box1 was removed were lost.
					So(si.len(), ShouldEqual, 6)
				})
			})

			Convey("And connecting the sink directly to the source", func() {
				So(sin.Input("source", nil), ShouldBeNil)
				so.EmitTuples(4)

				Convey("Then the sink should receive the correct number of tuples", func() {
					si.Wait(6) // 2 tuples which send just after box1 was removed were lost.
					So(si.len(), ShouldEqual, 6)
				})
			})
		})

		Convey("When removing a sink after receiving some tuples", func() {
			Reset(func() {
				t.Stop()
			})
			so.EmitTuples(2)
			si.Wait(2)
			So(t.Remove("SINK"), ShouldBeNil)
			so.EmitTuples(2)

			Convey("Then the sink should be stopped", func() {
				So(sin.State().Get(), ShouldEqual, TSStopped)
			})

			Convey("Then the sink shouldn't be found", func() {
				_, err := t.Sink("sink")
				So(err, ShouldNotBeNil)
			})

			Convey("Then sink shouldn't receive tuples generated after it got removed", func() {
				So(si.len(), ShouldEqual, 2)
			})
		})

		Convey("When generating some tuples and pause the source", func() { // 2.a.
			Reset(func() {
				t.Stop()
			})
			so.EmitTuples(4)
			So(son.Pause(), ShouldBeNil)
			so.EmitTuplesNB(4)

			Convey("Then the sink should only receive generated tuples", func() {
				si.Wait(4)
				So(si.len(), ShouldEqual, 4)
			})

			Convey("And resuming after that", func() {
				So(son.Resume(), ShouldBeNil)

				Convey("Then the sink should receive all tuples", func() {
					si.Wait(8)
					So(si.len(), ShouldEqual, 8)
				})
			})
		})

		Convey("When boxes stops on outbound disconnection", func() {
			bn1.StopOnDisconnect(Outbound)
			bn2.StopOnDisconnect(Outbound)

			Convey("Then bn1 should be stopped after stopping the sink", func() {
				So(sin.Stop(), ShouldBeNil)
				So(bn1.State().Wait(TSStopped), ShouldEqual, TSStopped)
			})
		})

		Convey("When only b1 stops on outbound disconnection and the sink is stopped", func() {
			bn1.StopOnDisconnect(Outbound)
			So(sin.Stop(), ShouldBeNil)

			Convey("Then bn2.StopOnDisconnect should eventually stop bn1", func() {
				bn2.StopOnDisconnect(Outbound)
				So(bn1.State().Wait(TSStopped), ShouldEqual, TSStopped)
			})
		})

		Convey("When only b1 stops on outbound disconnection and b2 is stopped", func() {
			bn1.StopOnDisconnect(Outbound)
			son.StopOnDisconnect()
			So(bn2.Stop(), ShouldBeNil)

			Convey("Then bn1.StopOnDisconnect should eventually stop son", func() {
				So(son.State().Wait(TSStopped), ShouldEqual, TSStopped)
			})
		})
	})
}

func TestForkDefaultTopology(t *testing.T) {
	Convey("Given a simple fork topology", t, func() {
		/*
		 *        /--> b1 -*--> si1
		 *   so -*
		 *        \--> b2 -*--> si2
		 */
		dt := NewDefaultTopology(NewContext(nil), "dt1")
		t := dt.(*defaultTopology)
		Reset(func() {
			t.Stop()
		})

		so := NewTupleIncrementalEmitterSource(freshTuples())
		_, err := t.AddSource("source", so, nil)
		So(err, ShouldBeNil)

		b1 := &BlockingForwardBox{cnt: 8}
		tc1 := newTerminateChecker(b1)
		bn1, err := t.AddBox("box1", tc1, nil)
		So(err, ShouldBeNil)
		So(bn1.Input("source", nil), ShouldBeNil)

		b2 := &BlockingForwardBox{cnt: 8}
		tc2 := newTerminateChecker(b2)
		bn2, err := t.AddBox("box2", tc2, nil)
		So(err, ShouldBeNil)
		So(bn2.Input("source", nil), ShouldBeNil)

		si1 := NewTupleCollectorSink()
		sic1 := &sinkCloseChecker{s: si1}
		sin1, err := t.AddSink("si1", sic1, nil)
		So(err, ShouldBeNil)
		So(sin1.Input("box1", nil), ShouldBeNil)

		si2 := NewTupleCollectorSink()
		sic2 := &sinkCloseChecker{s: si2}
		sin2, err := t.AddSink("si2", sic2, nil)
		So(err, ShouldBeNil)
		So(sin2.Input("box2", nil), ShouldBeNil)

		checkPostCond := func() {
			Convey("Then the topology should be stopped", func() {
				So(t.state.Get(), ShouldEqual, TSStopped)
			})

			Convey("Then Box.Terminate should be called exactly once", func() {
				So(tc1.terminateCnt, ShouldEqual, 1)
				So(tc2.terminateCnt, ShouldEqual, 1)
			})

			Convey("Then Sink.Close should be called exactly once", func() {
				So(sic1.closeCnt, ShouldEqual, 1)
				So(sic2.closeCnt, ShouldEqual, 1)
			})
		}

		Convey("When generating no tuples and call stop", func() { // 1.
			go func() {
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink shouldn't receive anything", func() {
				So(si1.len(), ShouldEqual, 0)
				So(si2.len(), ShouldEqual, 0)
			})
		})

		Convey("When generating some tuples and call stop before the sink receives a tuple", func() { // 2.a.
			b1.cnt = 0 // Block tuples. The box isn't running yet, so cnt can safely be changed.
			go func() {
				so.EmitTuplesNB(4)
				go func() {
					t.Stop()
				}()

				// resume b1 after the topology starts stopping all.
				t.state.Wait(TSStopping)
				b1.EmitTuples(8)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of generated tuples", func() {
				So(si1.len(), ShouldEqual, 4)
				So(si2.len(), ShouldEqual, 4)
			})
		})

		Convey("When generating some tuples and call stop after the sink received some of them", func() { // 2.b.
			go func() {
				so.EmitTuplesNB(3)
				t.state.Wait(TSRunning)
				b1.EmitTuples(1)
				si1.Wait(1)
				go func() {
					t.Stop()
				}()

				t.state.Wait(TSStopping)
				b1.EmitTuples(2)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(si1.len(), ShouldEqual, 3)
				So(si2.len(), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after both sinks received all", func() { // 2.c.
			go func() {
				so.EmitTuples(4)
				si1.Wait(4)
				si2.Wait(4)
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should only receive those tuples", func() {
				So(si1.len(), ShouldEqual, 4)
				So(si2.len(), ShouldEqual, 4)
			})
		})

		Convey("When generating all tuples and call stop before the sink receives a tuple", func() { // 3.a.
			b1.cnt = 0
			go func() {
				so.EmitTuples(100) // Blocking call. Assuming the pipe's capacity is greater than or equal to 8.
				go func() {
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(8)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(si1.len(), ShouldEqual, 8)
				So(si2.len(), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received some of them", func() { // 3.b.
			b1.cnt = 2
			go func() {
				so.EmitTuples(100)
				si1.Wait(2)
				// don't care about si2
				go func() {
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(6)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(si1.len(), ShouldEqual, 8)
				So(si2.len(), ShouldEqual, 8)
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
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(6)
				b2.EmitTuples(5)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(si1.len(), ShouldEqual, 8)
				So(si2.len(), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received all", func() { // 3.c.
			go func() {
				so.EmitTuples(100)
				si1.Wait(8)
				si2.Wait(8)
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(si1.len(), ShouldEqual, 8)
				So(si2.len(), ShouldEqual, 8)
			})
		})
	})
}

func TestJoinDefaultTopology(t *testing.T) {
	Convey("Given a simple join topology", t, func() {
		/*
		 *   so1 -*-\
		 *           --> b -*--> si
		 *   so2 -*-/
		 */
		tb := NewDefaultTopology(NewContext(nil), "dt1")
		t := tb.(*defaultTopology)
		Reset(func() {
			t.Stop()
		})

		so1 := NewTupleIncrementalEmitterSource(freshTuples()[0:4])
		_, err := t.AddSource("source1", so1, nil)
		So(err, ShouldBeNil)

		so2 := NewTupleIncrementalEmitterSource(freshTuples()[4:8])
		_, err = t.AddSource("source2", so2, nil)
		So(err, ShouldBeNil)

		b1 := &BlockingForwardBox{cnt: 8}
		tc1 := newTerminateChecker(b1)
		bn1, err := t.AddBox("box1", tc1, nil)
		So(err, ShouldBeNil)
		So(bn1.Input("source1", nil), ShouldBeNil)
		So(bn1.Input("source2", nil), ShouldBeNil)

		si := NewTupleCollectorSink()
		sic := &sinkCloseChecker{s: si}
		sin, err := t.AddSink("sink", sic, nil)
		So(err, ShouldBeNil)
		So(sin.Input("box1", nil), ShouldBeNil)

		checkPostCond := func() {
			Convey("Then the topology should be stopped", func() {
				So(t.state.Get(), ShouldEqual, TSStopped)
			})

			Convey("Then Box.Terminate should be called exactly once", func() {
				So(tc1.terminateCnt, ShouldEqual, 1)
			})

			Convey("Then Sink.Close should be called exactly once", func() {
				So(sic.closeCnt, ShouldEqual, 1)
			})
		}

		Convey("When generating no tuples and call stop", func() { // 1.
			go func() {
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink shouldn't receive anything", func() {
				So(si.len(), ShouldEqual, 0)
			})
		})

		Convey("When generating some tuples and call stop before the sink receives a tuple", func() { // 2.a.
			b1.cnt = 0 // Block tuples. The box isn't running yet, so cnt can safely be changed.
			go func() {
				so1.EmitTuplesNB(3)
				so2.EmitTuplesNB(1)
				go func() {
					t.Stop()
				}()

				// resume b1 after the topology starts stopping all.
				t.state.Wait(TSStopping)
				b1.EmitTuples(8)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of generated tuples", func() {
				So(si.len(), ShouldEqual, 4)
			})
		})

		Convey("When generating some tuples from one source and call stop before the sink receives a tuple", func() { // 2.a'.
			b1.cnt = 0 // Block tuples. The box isn't running yet, so cnt can safely be changed.
			go func() {
				so1.EmitTuplesNB(3)
				go func() {
					t.Stop()
				}()

				// resume b1 after the topology starts stopping all.
				t.state.Wait(TSStopping)
				b1.EmitTuples(8)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of generated tuples", func() {
				So(si.len(), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after the sink received some of them", func() { // 2.b.
			go func() {
				so1.EmitTuplesNB(1)
				so2.EmitTuplesNB(2)
				t.state.Wait(TSRunning)
				b1.EmitTuples(1)
				si.Wait(1)
				go func() {
					t.Stop()
				}()

				t.state.Wait(TSStopping)
				b1.EmitTuples(2)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(si.len(), ShouldEqual, 3)
			})
		})

		Convey("When generating some tuples and call stop after the sink received all", func() { // 2.c.
			go func() {
				so1.EmitTuples(2)
				so2.EmitTuples(1)
				si.Wait(3)
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should only receive those tuples", func() {
				So(si.len(), ShouldEqual, 3)
			})
		})

		Convey("When generating all tuples and call stop before the sink receives a tuple", func() { // 3.a.
			b1.cnt = 0
			go func() {
				so1.EmitTuples(100) // Blocking call. Assuming the pipe's capacity is greater than or equal to 8.
				so2.EmitTuples(100)
				go func() {
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(8)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(si.len(), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received some of them", func() { // 3.b.
			b1.cnt = 2
			go func() {
				so1.EmitTuples(100)
				so2.EmitTuples(100)
				si.Wait(2)
				go func() {
					t.Stop()
				}()
				t.state.Wait(TSStopping)
				b1.EmitTuples(6)
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all of those tuples", func() {
				So(si.len(), ShouldEqual, 8)
			})
		})

		Convey("When generating all tuples and call stop after the sink received all", func() { // 3.c.
			go func() {
				so1.EmitTuples(100)
				so2.EmitTuples(100)
				si.Wait(8)
				t.Stop()
			}()
			t.state.Wait(TSStopped)
			checkPostCond()

			Convey("Then the sink should receive all tuples", func() {
				So(si.len(), ShouldEqual, 8)
			})
		})
	})
}

func TestDefaultTopologyQueueDropMode(t *testing.T) {
	Convey("Given a simple linear topology", t, func() {
		dt := NewDefaultTopology(NewContext(nil), "dt1")
		t := dt.(*defaultTopology)
		Reset(func() {
			t.Stop()
		})

		ts := freshTuples()
		so := NewTupleIncrementalEmitterSource(ts)
		_, err := t.AddSource("source", so, nil)
		So(err, ShouldBeNil)

		Convey("When adding a box with DropLatest", func() {
			b1 := &BlockingForwardBox{cnt: 0}
			tc1 := newTerminateChecker(b1)
			bn1, err := t.AddBox("box1", tc1, nil)
			So(err, ShouldBeNil)
			So(bn1.Input("SOURCE", &BoxInputConfig{
				Capacity: 1,
				DropMode: DropLatest,
			}), ShouldBeNil)

			si := NewTupleCollectorSink()
			sic := &sinkCloseChecker{s: si}
			sin, err := t.AddSink("sink", sic, nil)
			So(err, ShouldBeNil)
			So(sin.Input("BOX1", nil), ShouldBeNil)

			Convey("Then it receives the oldest tuple", func() {
				// so.EmitTuples blocks until all tuples are written, so
				// the sink receives one or two tuples. Note that the box
				// blocks *after* receiving a tuple.
				so.EmitTuples(8)
				b1.EmitTuples(8)
				si.Wait(1)
				So(si.len(), ShouldBeBetweenOrEqual, 1, 2)
				so.m.Lock() // lock for ts[0].InputName
				Reset(so.m.Unlock)
				So(si.get(0), ShouldResemble, ts[0])
			})
		})

		Convey("When adding a box with DropOldest", func() {
			b1 := &BlockingForwardBox{cnt: 0}
			tc1 := newTerminateChecker(b1)
			bn1, err := t.AddBox("box1", tc1, nil)
			So(err, ShouldBeNil)
			So(bn1.Input("SOURCE", &BoxInputConfig{
				Capacity: 1,
				DropMode: DropOldest,
			}), ShouldBeNil)

			si := NewTupleCollectorSink()
			sic := &sinkCloseChecker{s: si}
			sin, err := t.AddSink("sink", sic, nil)
			So(err, ShouldBeNil)
			sin.State().Wait(TSRunning)
			So(sin.Input("BOX1", nil), ShouldBeNil)

			Convey("Then it receives the latest tuple", func() {
				so.EmitTuples(8)
				b1.EmitTuples(8)
				waitForInputTuplesExhausted(si, ts[len(ts)-1])
				So(si.len(), ShouldBeBetweenOrEqual, 1, 2)
				so.m.Lock() // lock for ts[len(ts)-1].InputName
				Reset(so.m.Unlock)
				So(si.getLast(), ShouldResemble, ts[len(ts)-1])
			})
		})

		Convey("When adding a sink with DropLatest", func() {
			si := NewTupleCollectorSink()
			sic := &sinkCloseChecker{s: si}
			sin, err := t.AddSink("sink", sic, nil)
			So(err, ShouldBeNil)
			So(sin.Input("SOURCE", &SinkInputConfig{
				Capacity: 1,
				DropMode: DropLatest,
			}), ShouldBeNil)

			Convey("Then it receives the oldest tuple", func() {
				so.EmitTuples(8)
				si.Wait(1)

				So(si.len(), ShouldBeGreaterThanOrEqualTo, 1)
				so.m.Lock() // lock for ts[0].InputName
				Reset(so.m.Unlock)
				So(si.get(0), ShouldResemble, ts[0])
			})
		})

		Convey("When adding a sink with DropOldest", func() {
			si := NewTupleCollectorSink()
			sic := &sinkCloseChecker{s: si}
			sin, err := t.AddSink("sink", sic, nil)
			So(err, ShouldBeNil)
			So(sin.Input("SOURCE", &SinkInputConfig{
				Capacity: 1,
				DropMode: DropOldest,
			}), ShouldBeNil)

			Convey("Then it only the latest tuple", func() {
				so.EmitTuples(8)
				waitForInputTuplesExhausted(si, ts[len(ts)-1])
				So(si.len(), ShouldBeGreaterThanOrEqualTo, 1)
				so.m.Lock() // lock for ts[len(ts)-1].InputName
				Reset(so.m.Unlock)
				So(si.getLast(), ShouldResemble, ts[len(ts)-1])
			})
		})
	})
}

func waitForInputTuplesExhausted(si *TupleCollectorSink, lastTuple *Tuple) {
	si.Wait(1)
	for si.getLast() != lastTuple {
		time.Sleep(time.Nanosecond)
	}
}

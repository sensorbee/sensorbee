package core

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
	"time"
)

func TestRewindableSource(t *testing.T) {
	Convey("Given a default topology", t, func() {
		t := NewDefaultTopology(NewContext(nil), "dt1")
		Reset(func() {
			t.Stop()
		})

		fts := freshTuples()
		so := NewTupleEmitterSource(fts)
		son, err := t.AddSource("source", NewRewindableSource(so), &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)
		son.State().Wait(TSPaused)

		b := &BlockingForwardBox{cnt: 1000}
		bn, err := t.AddBox("box", b, nil)
		So(err, ShouldBeNil)
		bn.State().Wait(TSRunning)
		So(bn.Input("source", &BoxInputConfig{
			Capacity: 1, // (almost) blocking channel
		}), ShouldBeNil)

		si := NewTupleCollectorSink()
		sin, err := t.AddSink("sink", si, nil)
		So(sin.Input("box", nil), ShouldBeNil)
		sin.State().Wait(TSRunning)

		Convey("When emitting all tuples", func() {
			So(son.Resume(), ShouldBeNil)
			si.Wait(8)

			Convey("Then the source shouldn't stop", func() {
				So(son.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then status should show that it's waiting for rewind", func() {
				waitForWaitingForRewind(son)
				st := son.Status()
				v, _ := st.Get(data.MustCompilePath("source.waiting_for_rewind"))
				So(v, ShouldEqual, data.True)
			})
		})

		Convey("When rewinding before sending any tuple", func() {
			So(son.Rewind(), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)
			si.Wait(8)

			Convey("Then the sink should receive all tuples", func() {
				So(si.len(), ShouldEqual, 8)
			})
		})

		Convey("When rewinding after sending some tuples", func() {
			// emit 2 tuples and then blocks. So, 2 tuples went to the sink,
			// 1 tuple is blocked in the box, and 2 tuples is blocked in the
			// channel between the source and box because its capacity is 1
			// (1 tuple is in the queue and the other one is blocked at sending
			// operation). So, 5 tuples in total could be emitted from the source.
			b.setCnt(2)
			So(son.Resume(), ShouldBeNil)
			si.Wait(2)
			So(son.Pause(), ShouldBeNil)
			b.EmitTuples(1000)
			So(son.Rewind(), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then all tuple should be able to be sent again", func() {
				waitForLastTuple(si, fts[len(fts)-1])

				// Due to cuncurrency, the number of tuples arriving to the sink
				// is not constant.
				offset := si.len() - len(fts)
				for i := range fts {
					So(si.get(offset+i), ShouldResemble, fts[i])
				}
			})
		})

		Convey("When sending all tuples", func() {
			So(son.Resume(), ShouldBeNil)
			si.Wait(8)

			Convey("The source should be able to be rewound", func() {
				So(son.Rewind(), ShouldBeNil)
				si.Wait(16)
				So(si.len(), ShouldEqual, 16)
			})

			Convey("Then source should stop without rewinding", func() {
				So(son.Stop(), ShouldBeNil)
			})
		})

		Convey("When rewinding the paused source", func() {
			So(son.Rewind(), ShouldBeNil)

			Convey("Then it shouldn't be resumed", func() {
				So(son.State().Get(), ShouldEqual, TSPaused)
			})

			Convey("Then it should show that it isn't waiting for rewind", func() {
				// It's still generating tuples, just the process is being paused.
				st := son.Status()
				v, _ := st.Get(data.MustCompilePath("source.waiting_for_rewind"))
				So(v, ShouldEqual, data.False)
			})
		})

		Convey("When calling Rewind on non-rewindable source", func() {
			so2 := NewTupleEmitterSource(freshTuples())
			son2, err := t.AddSource("source2", so2, nil)
			So(err, ShouldBeNil)

			Convey("Then it should fail", func() {
				So(son2.Rewind(), ShouldNotBeNil)
			})
		})

		Convey("When calling Rewind on the stopped source", func() {
			So(son.Stop(), ShouldBeNil)

			Convey("Then it should fail", func() {
				So(son.Rewind(), ShouldNotBeNil)
			})
		})
	})
}

type dummyNonstoppableSource struct {
	stopped bool
}

func (d *dummyNonstoppableSource) GenerateStream(ctx *Context, w Writer) error {
	defer func() {
		d.stopped = true
	}()
	for {
		if err := w.Write(ctx, NewTuple(data.Map{"a": data.True})); err != nil {
			return err
		}
	}
}

func (d *dummyNonstoppableSource) Stop(ctx *Context) error {
	return nil
}

func TestImplementSourceStop(t *testing.T) {
	Convey("Given a stoppable source via ImplementSourceStop", t, func() {
		ctx := NewContext(nil)
		s := ImplementSourceStop(NewTupleEmitterSource(freshTuples()))
		Reset(func() {
			s.Stop(ctx)
		})

		ch := make(chan error, 1)
		go func() {
			ch <- s.GenerateStream(ctx, WriterFunc(func(ctx *Context, t *Tuple) error {
				return nil
			}))
		}()

		Convey("When waiting for the source to be stopped", func() {
			err := <-ch

			Convey("Then it should stop without explicitly calling Stop method", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When converting it to RewindableSource", func() {
			_, ok := s.(RewindableSource)

			Convey("Then it should fail", func() {
				So(ok, ShouldBeFalse)
			})
		})
	})

	Convey("Given a non stoppable source via ImplementSourceStop", t, func() {
		ctx := NewContext(nil)
		ns := &dummyNonstoppableSource{}
		s := ImplementSourceStop(ns)
		Reset(func() {
			s.Stop(ctx)
		})

		ch := make(chan error, 1)
		go func() {
			ch <- s.GenerateStream(ctx, WriterFunc(func(ctx *Context, t *Tuple) error {
				return nil
			}))
		}()

		Convey("When stopping the source", func() {
			So(s.Stop(ctx), ShouldBeNil)

			Convey("Then the original source should stop", func() {
				So(ns.stopped, ShouldBeTrue)
			})

			Convey("Then GenerateStream should return", func() {
				So(<-ch, ShouldBeNil)
			})
		})

		Convey("When converting it to RewindableSource", func() {
			_, ok := s.(RewindableSource)

			Convey("Then it should fail", func() {
				So(ok, ShouldBeFalse)
			})
		})
	})
}

type dummyBlockingSource struct {
}

func (d *dummyBlockingSource) GenerateStream(ctx *Context, w Writer) error {
	// This sleep is only for avoiding a goroutine leak and all tests assume
	// that they finish before this method returns.
	time.Sleep(10 * time.Second)
	return nil
}

func (d *dummyBlockingSource) Stop(ctx *Context) error {
	return errors.New("cannot stop blocking source")
}

func TestRewindableSourceForceStop(t *testing.T) {
	Convey("Given a source whose GenerateStream will never return and whose Stop fails", t, func() {
		ctx := NewContext(nil)
		s := ImplementSourceStop(&dummyBlockingSource{})
		Reset(func() {
			s.Stop(ctx)
		})

		ch := make(chan error, 1)
		go func() {
			ch <- s.GenerateStream(ctx, WriterFunc(func(ctx *Context, t *Tuple) error {
				return nil
			}))
		}()

		Convey("When stopping the source", func() {
			err := s.Stop(ctx)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "cannot stop blocking source")
			})

			Convey("Then GenerateStream should stop", func() {
				err := <-ch
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func waitForWaitingForRewind(son SourceNode) {
	so := son.Source()
	rso := so.(*rewindableSource)
	for {
		rso.rwm.RLock()
		if rso.waitingForRewind {
			rso.rwm.RUnlock()
			return
		}
		rso.rwm.RUnlock()
	}
}

func waitForLastTuple(si *TupleCollectorSink, t *Tuple) {
	for si.get(si.len()-1) != t {
		time.Sleep(time.Nanosecond)
	}
}

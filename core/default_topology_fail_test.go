package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

type stubInitTerminateBox struct {
	*terminateChecker

	init struct {
		cnt         int
		shouldFail  bool
		shouldPanic bool
		shouldBlock bool
		failed      bool
		wg          sync.WaitGroup
	}

	terminate struct {
		shouldFail  bool
		shouldPanic bool
	}

	panicValue interface{}
}

var _ StatefulBox = &stubInitTerminateBox{}

func newStubInitTerminateBox(b Box) *stubInitTerminateBox {
	return &stubInitTerminateBox{
		terminateChecker: newTerminateChecker(b),
	}
}

func (s *stubInitTerminateBox) Init(ctx *Context) error {
	i := &s.init
	i.cnt++
	if i.shouldBlock {
		i.wg.Add(1)
		i.wg.Wait()
	}
	if i.shouldPanic {
		i.failed = true
		if s.panicValue != nil {
			panic(s.panicValue)
		}
		panic(fmt.Errorf("failure panic"))
	}
	if i.shouldFail {
		i.failed = true
		return fmt.Errorf("failure")
	}
	return nil
}

func (s *stubInitTerminateBox) ResumeInit() {
	s.init.wg.Done()
}

func (s *stubInitTerminateBox) Terminate(ctx *Context) error {
	if s.terminate.shouldPanic {
		panic(fmt.Errorf("failure panic"))
	}
	if s.terminate.shouldFail {
		return fmt.Errorf("failure")
	}
	if err := s.terminateChecker.Terminate(ctx); err != nil {
		return err
	}
	return nil
}

type panicBox struct {
	ProxyBox

	m            sync.Mutex
	writeFailAt  int
	fatalError   bool
	writePanicAt int
	writeCnt     int
	panicValue   interface{}
}

func (b *panicBox) Process(ctx *Context, t *Tuple, w Writer) error {
	b.m.Lock()
	defer b.m.Unlock()
	b.writeCnt++
	if b.writeCnt == b.writePanicAt {
		if b.panicValue != nil {
			panic(b.panicValue)
		}
		panic(fmt.Errorf("test failure via panic"))
	}
	if b.writeCnt == b.writeFailAt {
		err := fmt.Errorf("test failure")
		if b.fatalError {
			return FatalError(err)
		}
		return err
	}
	return b.ProxyBox.Process(ctx, t, w)
}

func TestDefaultTopologySetupFailure(t *testing.T) {
	Convey("Given a topology", t, func() {
		t := NewDefaultTopology(NewContext(nil), "dt1")
		Reset(func() {
			t.Stop()
		})

		Convey("When adding a node", func() {
			Convey("Then a source having an invalid name shouldn't be added", func() {
				_, err := t.AddSource("0invalid", &DoesNothingSource{}, nil)
				So(err, ShouldNotBeNil)
			})

			Convey("Then a box having an invalid name shouldn't be added", func() {
				_, err := t.AddBox("in-valid", &DoesNothingBox{}, nil)
				So(err, ShouldNotBeNil)
			})

			Convey("Then a sink having an invalid name shouldn't be added", func() {
				_, err := t.AddSink("in*valid", &DoesNothingSink{}, nil)
				So(err, ShouldNotBeNil)
			})
		})

		// TODO: Add Source failure tests with a duplicated name or failing GenerateStream
		//       1. Stop fails
		//       2. Stop panics

		Convey("When adding a box panicing in Init", func() {
			Convey("Then it shouldn't be added", func() {
				b := newStubInitTerminateBox(&DoesNothingBox{})
				b.init.shouldPanic = true
				tc := newTerminateChecker(b)
				_, err := t.AddBox("box", tc, nil)
				So(err, ShouldNotBeNil)
				So(tc.terminateCnt, ShouldEqual, 0)
			})

			Convey("Then it shouldn't be added even if panic doesn't have an error", func() {
				b := newStubInitTerminateBox(&DoesNothingBox{})
				b.init.shouldPanic = true
				b.panicValue = 1
				tc := newTerminateChecker(b)
				_, err := t.AddBox("box", tc, nil)
				So(err, ShouldNotBeNil)
				So(tc.terminateCnt, ShouldEqual, 0)
			})
		})

		Convey("When adding a box failing in Init", func() {
			Convey("Then it shouldn't be added", func() {
				b := newStubInitTerminateBox(&DoesNothingBox{})
				b.init.shouldFail = true
				tc := newTerminateChecker(b)
				_, err := t.AddBox("box", tc, nil)
				So(err, ShouldNotBeNil)
				So(tc.terminateCnt, ShouldEqual, 0)
			})
		})
	})
}

func TestDefaultTopologyFailure(t *testing.T) {
	Convey("Given a simple linear topology", t, func() {
		/*
		 *   so -*--> b1 -*--> si
		 */
		t := NewDefaultTopology(NewContext(nil), "dt1")
		Reset(func() {
			t.Stop()
		})

		so := NewTupleIncrementalEmitterSource(freshTuples())
		son, err := t.AddSource("source", so, nil)
		So(err, ShouldBeNil)

		sb := newStubInitTerminateBox(BoxFunc(forwardBox))
		b1 := &panicBox{
			ProxyBox: ProxyBox{
				b: sb,
			},
		}
		bn1, err := t.AddBox("box1", b1, nil)
		So(err, ShouldBeNil)
		So(bn1.Input("source", nil), ShouldBeNil)

		si := NewTupleCollectorSink()
		sic := &sinkCloseChecker{s: si}
		sin, err := t.AddSink("sink", sic, nil)
		So(err, ShouldBeNil)
		So(sin.Input("box1", nil), ShouldBeNil)

		Convey("When adding a new source with the duplicated name", func() {
			_, err := t.AddSource("SOURCE", so, nil)

			Convey("Then it should fail", func() {
				So(err.Error(), ShouldContainSubstring, "already used")
			})
		})

		Convey("When adding a new sink with the duplicated name", func() {
			_, err := t.AddSink("SINK", sic, nil)

			Convey("Then it should fail", func() {
				So(err.Error(), ShouldContainSubstring, "already used")
			})
		})

		Convey("When adding a new box with the duplicated name", func() {
			_, err := t.AddBox("BOX1", b1, nil)

			Convey("Then it should fail", func() {
				So(err.Error(), ShouldContainSubstring, "already used")
			})
		})

		Convey("When adding a new input with the duplicated name", func() {
			err := sin.Input("BOX1", nil)

			Convey("Then it should fail", func() {
				So(err.Error(), ShouldContainSubstring, "already")
			})
		})

		Convey("When adding a input after the source stops", func() {
			So(son.Stop(), ShouldBeNil)

			Convey("Then the box shouldn't be able to get input from it", func() {
				bn2, err := t.AddBox("box2", &DoesNothingBox{}, nil)
				So(err, ShouldBeNil)
				So(bn2.Input("source", nil), ShouldNotBeNil)
			})

			Convey("Then the sink sholdn't be able to get input from it", func() {
				sin2, err := t.AddSink("sink2", &DoesNothingSink{}, nil)
				So(err, ShouldBeNil)
				So(sin2.Input("source", nil), ShouldNotBeNil)
			})
		})

		Convey("When adding a new input to stopped node", func() {
			Convey("Then the stopped box shouldn't accept the new input", func() {
				bn2, err := t.AddBox("box2", &DoesNothingBox{}, nil)
				So(err, ShouldBeNil)
				So(bn2.Stop(), ShouldBeNil)
				So(bn2.Input("source", nil), ShouldNotBeNil)
			})

			Convey("Then the stopped sink shouldn't accept the new input", func() {
				sin2, err := t.AddSink("sink2", &DoesNothingSink{}, nil)
				So(err, ShouldBeNil)
				So(sin2.Stop(), ShouldBeNil)
				So(sin2.Input("source", nil), ShouldNotBeNil)
			})
		})

		Convey("When adding a nonexistent input to a node", func() {
			Convey("Then the box shouldn't accept it", func() {
				So(bn1.Input("no_such_source", nil), ShouldNotBeNil)
			})

			Convey("Then the sink shouldn't accept it", func() {
				So(sin.Input("no_such_source", nil), ShouldNotBeNil)
			})
		})

		Convey("When a box panics", func() {
			b1.writePanicAt = 1
			so.EmitTuples(5)

			Convey("Then the box should be stopped", func() {
				So(bn1.State().Wait(TSStopped), ShouldEqual, TSStopped)

				Convey("And the box should be terminated", func() {
					So(sb.terminateCnt, ShouldEqual, 1)
				})
			})

			Convey("Then the topology can be recovered by manual connection", func() {
				So(sin.Input("source", nil), ShouldBeNil)
				so.EmitTuples(3)
				si.Wait(3)
				So(len(si.Tuples), ShouldEqual, 3)
			})
		})

		Convey("When a box fails with a regular error", func() {
			b1.writeFailAt = 4
			so.EmitTuples(8)

			Convey("Then the box shouldn't be stopped", func() {
				So(bn1.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then a tuple should be lost", func() {
				si.Wait(7)
				So(len(si.Tuples), ShouldEqual, 7)
			})
		})

		Convey("When a box fails with a fatal error", func() {
			b1.writeFailAt = 4
			b1.fatalError = true
			so.EmitTuples(8)

			Convey("Then the box should be stopped", func() {
				So(bn1.State().Wait(TSStopped), ShouldEqual, TSStopped)
			})

			Convey("Then the sink should receive tuples sent before the fatal error", func() {
				si.Wait(3)
				So(len(si.Tuples), ShouldEqual, 3)
			})
		})

		// TODO: Source.Stop panic

		Convey("When a box panics on terminate", func() {
			sb.terminate.shouldPanic = true

			Convey("Then it should be stopped", func() {
				So(bn1.Stop(), ShouldBeNil)

				Convey("And the box has an error in its status", func() {
					st := bn1.Status()
					_, ok := st["error"]
					So(ok, ShouldBeTrue)
				})
			})

			Convey("Then it should be stopped when a fatal write error occurs", func() {
				b1.writeFailAt = 1
				b1.fatalError = true
				so.EmitTuples(8)
				So(bn1.State().Wait(TSStopped), ShouldEqual, TSStopped)

				Convey("And the box has an error in its status", func() {
					st := bn1.Status()
					_, ok := st["error"]
					So(ok, ShouldBeTrue)
				})
			})
		})

		Convey("When a box fails on terminate", func() {
			sb.terminate.shouldFail = true

			Convey("Then it should be stopped", func() {
				So(bn1.Stop(), ShouldBeNil)

				Convey("And the box has an error in its status", func() {
					st := bn1.Status()
					_, ok := st["error"]
					So(ok, ShouldBeTrue)
				})
			})

			Convey("Then it should be stopped when a fatal write error occurs", func() {
				b1.writeFailAt = 1
				b1.fatalError = true
				so.EmitTuples(8)
				So(bn1.State().Wait(TSStopped), ShouldEqual, TSStopped)

				Convey("And the box has an error in its status", func() {
					st := bn1.Status()
					_, ok := st["error"]
					So(ok, ShouldBeTrue)
				})
			})
		})

		// TODO: Box terminate panic and handling
		// TODO: Sink.Close panic
	})
}

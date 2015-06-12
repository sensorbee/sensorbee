package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync/atomic"
	"testing"
)

type stubSharedState struct {
	initCnt     int
	initFailAt  int
	initPanicAt int

	terminateCnt     int
	terminateFailAt  int
	terminatePanicAt int
}

func (s *stubSharedState) TypeName() string {
	return "mock_shared_state"
}

func (s *stubSharedState) Init(ctx *Context) error {
	s.initCnt++
	if s.initCnt == s.initPanicAt {
		panic(fmt.Errorf("mock shared state panic"))
	}
	if s.initCnt == s.initFailAt {
		return fmt.Errorf("mock shared state failure")
	}
	return nil
}

func (s *stubSharedState) Write(ctx *Context, t *tuple.Tuple) error {
	return nil
}

func (s *stubSharedState) Terminate(ctx *Context) error {
	s.terminateCnt++
	if s.terminateCnt == s.terminatePanicAt {
		panic(fmt.Errorf("mock shared state panic"))
	}
	if s.terminateCnt == s.terminateFailAt {
		return fmt.Errorf("mock shared state failure")
	}
	return nil
}

func TestDefaultSharedStateRegistry(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a default shared state registry", t, func() {
		r := NewDefaultSharedStateRegistry()

		Convey("When adding a state", func() {
			s := &stubSharedState{}
			So(r.Add(ctx, "test_state", s), ShouldBeNil)

			Convey("Then a state having the same name cannot be added", func() {
				So(r.Add(ctx, "test_state", &stubSharedState{}), ShouldNotBeNil)
			})

			Convey("Then a state which fails on termination and has the same name cannot be added", func() {
				s2 := &stubSharedState{}
				s2.terminateFailAt = 1
				err := r.Add(ctx, "test_state", s2)

				Convey("And the error should be about the name duplication, not termination failure", func() {
					So(err.Error(), ShouldContainSubstring, "already has")
				})
			})

			Convey("Then Get should return it", func() {
				s2, err := r.Get(ctx, "test_state")
				So(err, ShouldBeNil)
				So(s2, ShouldPointTo, s)
			})

			Convey("Then it should be listed", func() {
				m, err := r.List(ctx)
				So(err, ShouldBeNil)
				So(len(m), ShouldEqual, 1)
				So(m["test_state"], ShouldPointTo, s)
			})

			Convey("Then it can be removed", func() {
				s2, err := r.Remove(ctx, "test_state")
				So(err, ShouldBeNil)

				Convey("And the returned state should be correct", func() {
					So(s2, ShouldPointTo, s)
				})

				Convey("And it shouldn't be able to be removed twice", func() {
					s3, err := r.Remove(ctx, "test_state")
					So(err, ShouldBeNil)
					So(s3, ShouldBeNil)
				})

				Convey("And it should be terminated", func() {
					So(s.terminateCnt, ShouldEqual, 1)
				})
			})
		})

		Convey("When a state panics on initialization while adding it to registry", func() {
			s := &stubSharedState{}
			s.initPanicAt = 1

			Convey("Then Add should panic", func() {
				So(func() {
					r.Add(ctx, "test_state", s)
				}, ShouldPanic)

				Convey("And the state shouldn't be added", func() {
					_, err := r.Get(ctx, "test_state")
					So(err, ShouldNotBeNil)
				})

				Convey("And terminate shouldn't be called", func() {
					So(s.terminateCnt, ShouldEqual, 0)
				})
			})
		})

		Convey("When a state fails to initialize while adding it to registry", func() {
			s := &stubSharedState{}
			s.initFailAt = 1
			err := r.Add(ctx, "test_state", s)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the state shouldn't be added", func() {
				_, err := r.Get(ctx, "test_state")
				So(err, ShouldNotBeNil)
			})

			Convey("Then terminate shouldn't be called", func() {
				So(s.terminateCnt, ShouldEqual, 0)
			})
		})

		Convey("When getting a nonexistent state", func() {
			s, err := r.Get(ctx, "test_state")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the returned state should be nil", func() {
				So(s, ShouldBeNil)
			})
		})

		Convey("When listing an empty registry", func() {
			m, err := r.List(ctx)
			So(err, ShouldBeNil)

			Convey("Then the list should be empty", func() {
				So(m, ShouldBeEmpty)
			})
		})

		Convey("When removing a state whose termination panics", func() {
			s := &stubSharedState{}
			s.terminatePanicAt = 1
			So(r.Add(ctx, "test_state", s), ShouldBeNil)

			Convey("Then it should panic", func() {
				So(func() {
					r.Remove(ctx, "test_state")
				}, ShouldPanic)

				Convey("And it should've been removed", func() {
					_, err := r.Get(ctx, "test_state")
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When removing a state whose termination fails", func() {
			s := &stubSharedState{}
			s.terminateFailAt = 1
			So(r.Add(ctx, "test_state", s), ShouldBeNil)
			s2, err := r.Remove(ctx, "test_state")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the state should be returned even on failure", func() {
				So(s2, ShouldPointTo, s)
			})

			Convey("Then it should've been removed", func() {
				_, err := r.Get(ctx, "test_state")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When removing a nonexistent state", func() {
			s, err := r.Remove(ctx, "test_state")

			Convey("Then it shouldn't return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the returned state should be nil", func() {
				So(s, ShouldBeNil)
			})
		})
	})
}

type countingSharedState struct {
	cnt int32
}

func (c *countingSharedState) TypeName() string {
	return "counter"
}

func (c *countingSharedState) Init(ctx *Context) error {
	return nil
}

func (c *countingSharedState) Write(ctx *Context, t *tuple.Tuple) error {
	i, _ := tuple.ToInt(t.Data["n"])
	atomic.AddInt32(&c.cnt, int32(i))
	return nil
}

func (c *countingSharedState) Terminate(ctx *Context) error {
	return nil
}

func TestSharedStateInStaticTopology(t *testing.T) {
	Convey("Given a static topology and a state", t, func() {
		// Run this test with parallelism 1
		config := Configuration{TupleTraceEnabled: 1}
		ctx := newTestContext(config)
		counter := &countingSharedState{}
		So(ctx.SharedStates.Add(ctx, "test_counter", counter), ShouldBeNil)

		tb := NewDefaultStaticTopologyBuilder()

		ts := []*tuple.Tuple{}
		for i := 0; i < 4; i++ {
			ts = append(ts, &tuple.Tuple{
				Data: tuple.Map{
					"n": tuple.Int(i + 1),
				},
			})
		}
		so := NewTupleEmitterSource(ts)
		So(tb.AddSource("source", so).Err(), ShouldBeNil)

		b1 := BoxFunc(func(ctx *Context, t *tuple.Tuple, w Writer) error {
			s, err := ctx.GetSharedState("test_counter")
			if err != nil {
				return err
			}
			s.Write(ctx, t)
			return w.Write(ctx, t)
		})
		So(tb.AddBox("box1", b1).Input("source").Err(), ShouldBeNil)

		b2 := BoxFunc(func(ctx *Context, t *tuple.Tuple, w Writer) error {
			s, err := ctx.GetSharedState("test_counter")
			if err != nil {
				return err
			}
			if s.TypeName() != "counter" {
				return fmt.Errorf("unsupported state type: %v", s.TypeName())
			}
			c, ok := s.(*countingSharedState)
			if !ok {
				return fmt.Errorf("cannot convert a state to a counter")
			}

			t.Data["cur_cnt"] = tuple.Int(atomic.LoadInt32(&c.cnt))
			return w.Write(ctx, t)
		})
		So(tb.AddBox("box2", b2).Input("box1").Err(), ShouldBeNil)

		si := NewTupleCollectorSink()
		sic := &sinkCloseChecker{s: si}
		So(tb.AddSink("sink", sic).Input("box2").Err(), ShouldBeNil)

		ti, err := tb.Build()
		So(err, ShouldBeNil)

		t := ti.(*defaultStaticTopology)
		go func() {
			t.Run(ctx)
		}()
		Reset(func() {
			t.Stop(ctx)
		})
		t.state.Wait(TSRunning)
		si.Wait(4)

		Convey("When running a topology with boxes refering to the state", func() {
			Convey("Then each tuple has correct counters", func() {
				for i, t := range si.Tuples {
					So(t.Data["cur_cnt"], ShouldBeGreaterThanOrEqualTo, (i+1)*(i+2)/2)
				}
			})

			Convey("Then the cnt should be 10", func() {
				So(counter.cnt, ShouldEqual, 10)
			})
		})
	})
}

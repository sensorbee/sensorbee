package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"sync/atomic"
	"testing"
)

type stubSharedState struct {
	terminateCnt     int
	terminateFailAt  int
	terminatePanicAt int
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
	Convey("Given a default shared state registry", t, func() {
		ctx := NewContext(nil)
		r := ctx.SharedStates

		Convey("When adding a state", func() {
			s := &stubSharedState{}
			So(r.Add("test_state", "test_state_type", s), ShouldBeNil)

			Convey("Then a state having the same name cannot be added", func() {
				So(r.Add("test_state", "test_state_type", &stubSharedState{}), ShouldNotBeNil)
			})

			Convey("Then a state which fails on termination and has the same name cannot be added", func() {
				s2 := &stubSharedState{}
				s2.terminateFailAt = 1
				err := r.Add("test_state", "test_state_type", s2)

				Convey("And the error should be about the name duplication, not termination failure", func() {
					So(err.Error(), ShouldContainSubstring, "already has")
				})
			})

			Convey("Then Get should return it", func() {
				s2, err := r.Get("test_state")
				So(err, ShouldBeNil)
				So(s2, ShouldPointTo, s)
			})

			Convey("Then Type should return its type name", func() {
				tn, err := r.Type("test_state")
				So(err, ShouldBeNil)
				So(tn, ShouldEqual, "test_state_type")
			})

			Convey("Then it can be replaced", func() {
				prev, err := r.Replace("test_state", "test_state_type", &stubSharedState{}, true)
				So(err, ShouldBeNil)
				Reset(func() {
					prev.Terminate(ctx)
				})
				So(prev, ShouldPointTo, s)

				Convey("And it shouldn't have been terminated yet", func() {
					So(s.terminateCnt, ShouldEqual, 0)
				})
			})

			Convey("Then it cannot be replaced with a wrong type name", func() {
				_, err := r.Replace("test_state", "wrong_type_name", &stubSharedState{}, true)
				So(err, ShouldNotBeNil)
			})

			Convey("Then it should be listed", func() {
				m, err := r.List()
				So(err, ShouldBeNil)
				So(len(m), ShouldEqual, 1)
				So(m["test_state"], ShouldPointTo, s)
			})

			Convey("Then it can be removed", func() {
				s2, err := r.Remove("test_state")
				So(err, ShouldBeNil)

				Convey("And the returned state should be correct", func() {
					So(s2, ShouldPointTo, s)
				})

				Convey("And it shouldn't be able to be removed twice", func() {
					s3, err := r.Remove("test_state")
					So(err, ShouldBeNil)
					So(s3, ShouldBeNil)
				})

				Convey("And it should be terminated", func() {
					So(s.terminateCnt, ShouldEqual, 1)
				})
			})
		})

		Convey("When getting a nonexistent state", func() {
			s, err := r.Get("test_state")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the returned state should be nil", func() {
				So(s, ShouldBeNil)
			})
		})

		Convey("When getting a type name of a nonexistent state", func() {
			_, err := r.Type("test_state")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When replacing a nonexistent state without creating if not exists", func() {
			s := &stubSharedState{}
			_, err := r.Replace("test_state", "test_state_type", s, false)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the state should be terminated", func() {
				So(s.terminateCnt, ShouldEqual, 1)
			})
		})

		Convey("When replacing a nonexistent state with creating if not exists", func() {
			s := &stubSharedState{}
			_, err := r.Replace("test_state", "test_state_type", s, true)
			So(err, ShouldBeNil)

			Convey("Then Get should return it", func() {
				s2, err := r.Get("test_state")
				So(err, ShouldBeNil)
				So(s2, ShouldPointTo, s)
			})
		})

		Convey("When replacing fails and the state.Terminate fails", func() {
			s := &stubSharedState{}
			s.terminateFailAt = 1
			_, err := r.Replace("test_state", "test_state_type", s, false)

			Convey("Then it should fail and the error message shouldn't be about the termination", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})

		Convey("When replacing fails and the state.Terminate panics", func() {
			s := &stubSharedState{}
			s.terminatePanicAt = 1
			_, err := r.Replace("test_state", "test_state_type", s, false)

			Convey("Then it should fail and the error message shouldn't be about the termination", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "was not found")
			})
		})

		Convey("When listing an empty registry", func() {
			m, err := r.List()
			So(err, ShouldBeNil)

			Convey("Then the list should be empty", func() {
				So(m, ShouldBeEmpty)
			})
		})

		Convey("When removing a state whose termination panics", func() {
			s := &stubSharedState{}
			s.terminatePanicAt = 1
			So(r.Add("test_state", "test_state_type", s), ShouldBeNil)

			Convey("Then it should panic", func() {
				So(func() {
					r.Remove("test_state")
				}, ShouldPanic)

				Convey("And it should've been removed", func() {
					_, err := r.Get("test_state")
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When removing a state whose termination fails", func() {
			s := &stubSharedState{}
			s.terminateFailAt = 1
			So(r.Add("test_state", "test_state_type", s), ShouldBeNil)
			s2, err := r.Remove("test_state")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the state should be returned even on failure", func() {
				So(s2, ShouldPointTo, s)
			})

			Convey("Then it should've been removed", func() {
				_, err := r.Get("test_state")
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When removing a nonexistent state", func() {
			s, err := r.Remove("test_state")

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
	terminated bool
	cnt        int32
}

func (c *countingSharedState) Write(ctx *Context, t *Tuple) error {
	i, _ := data.ToInt(t.Data["n"])
	atomic.AddInt32(&c.cnt, int32(i))
	return nil
}

func (c *countingSharedState) Terminate(ctx *Context) error {
	c.terminated = true
	return nil
}

func TestSharedStateInTopology(t *testing.T) {
	Convey("Given a topology and a state", t, func() {
		// Run this test with parallelism 1
		ctx := NewContext(nil)
		counter := &countingSharedState{}
		So(ctx.SharedStates.Add("test_counter", "test_counter_type", counter), ShouldBeNil)

		t := NewDefaultTopology(ctx, "test")
		Reset(func() {
			t.Stop()
		})

		ts := []*Tuple{}
		for i := 0; i < 4; i++ {
			ts = append(ts, &Tuple{
				Data: data.Map{
					"n": data.Int(i + 1),
				},
			})
		}
		so := NewTupleEmitterSource(ts)
		son, err := t.AddSource("source", so, &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		b1 := BoxFunc(func(ctx *Context, t *Tuple, w Writer) error {
			s, err := ctx.SharedStates.Get("test_counter")
			if err != nil {
				return err
			}
			s.(Writer).Write(ctx, t)
			return w.Write(ctx, t)
		})
		bn1, err := t.AddBox("box1", b1, nil)
		So(err, ShouldBeNil)
		So(bn1.Input("source", nil), ShouldBeNil)

		b2 := BoxFunc(func(ctx *Context, t *Tuple, w Writer) error {
			s, err := ctx.SharedStates.Get("test_counter")
			if err != nil {
				return err
			}
			c, ok := s.(*countingSharedState)
			if !ok {
				return fmt.Errorf("cannot convert a state to a counter")
			}

			t.Data["cur_cnt"] = data.Int(atomic.LoadInt32(&c.cnt))
			return w.Write(ctx, t)
		})
		bn2, err := t.AddBox("box2", b2, nil)
		So(err, ShouldBeNil)
		So(bn2.Input("box1", nil), ShouldBeNil)

		si := NewTupleCollectorSink()
		sic := &sinkCloseChecker{s: si}
		sin, err := t.AddSink("sink", sic, nil)
		So(err, ShouldBeNil)
		So(sin.Input("box2", nil), ShouldBeNil)
		So(son.Resume(), ShouldBeNil)
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

func TestSharedStateSink(t *testing.T) {
	ctx := NewContext(nil)
	r := ctx.SharedStates
	counter := &countingSharedState{}

	{
		s := &stubSharedState{}
		if err := r.Add("test_state", "test_state_type", s); err != nil {
			t.Fatal("Cannot add stub shared state:", err)
		}

		if err := r.Add("test_counter", "test_counter_type", counter); err != nil {
			t.Fatal("Cannot add counting shared state:", err)
		}
	}

	Convey("Given a topology and a state registry", t, func() {
		Convey("When creating a shared state sink with non-Writer shared state", func() {
			_, err := NewSharedStateSink(ctx, "test_state")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When creating a shared state sink with non-existing state", func() {
			_, err := NewSharedStateSink(ctx, "no_such_sink")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When creating shared state sink with Writer shared state", func() {
			sink, err := NewSharedStateSink(ctx, "test_counter")
			Reset(func() {
				counter.cnt = 0
			})

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("And adding counter via sink", func() {
				t := &Tuple{
					Data: data.Map{"n": data.Int(100)},
				}
				err := sink.Write(ctx, t)

				Convey("Then it should succeed", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then the count should be be incremented", func() {
					So(counter.cnt, ShouldEqual, 100)
				})
			})

			Convey("Then it should be closed", func() {
				So(sink.Close(ctx), ShouldBeNil)

				Convey("And the state shouldn't be terminated", func() {
					So(counter.terminated, ShouldBeFalse)
				})
			})
		})
	})
}

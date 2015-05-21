package core

// Because most successful cases are tested in other tests,
// test cases in this file focus on failures.

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

type stubInitTerminateBoxSharedConfig struct {
	initFailAt int
	initCnt    int

	terminateFailAt  int
	terminatePanicAt int
	terminateCnt     int
}

type stubInitTerminateBox struct {
	*terminateChecker
	init struct {
		called bool
		block  bool
		failed bool
		wg     sync.WaitGroup
	}

	shared *stubInitTerminateBoxSharedConfig
}

func NewStubInitTerminateBox(b Box, s *stubInitTerminateBoxSharedConfig) *stubInitTerminateBox {
	return &stubInitTerminateBox{
		terminateChecker: NewTerminateChecker(b),
		shared:           s,
	}
}

func (s *stubInitTerminateBox) Init(ctx *Context) error {
	i := &s.init
	i.called = true
	if i.block {
		i.wg.Add(1)
		i.wg.Wait()
	}
	s.shared.initCnt++
	if s.shared.initCnt == s.shared.initFailAt {
		i.failed = true
		return fmt.Errorf("failure")
	}
	return nil
}

func (s *stubInitTerminateBox) ResumeInit() {
	s.init.wg.Done()
}

func (s *stubInitTerminateBox) Terminate(ctx *Context) error {
	s.shared.terminateCnt++
	if err := s.terminateChecker.Terminate(ctx); err != nil {
		return err
	}
	if s.shared.terminateCnt == s.shared.terminatePanicAt {
		panic(fmt.Errorf("failure"))
	}
	if s.shared.terminateCnt == s.shared.terminateFailAt {
		return fmt.Errorf("failure")
	}
	return nil
}

func TestDefaultStaticTopologyRun(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a default static topology", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := NewTupleEmitterSource(freshTuples())
		b1 := NewStubInitTerminateBox(BoxFunc(forwardBox), &stubInitTerminateBoxSharedConfig{})
		b2 := &BlockingForwardBox{cnt: 8}
		// Sink isn't necessary

		So(tb.AddSource("source", s).Err(), ShouldBeNil)
		So(tb.AddBox("box1", b1).Input("source").Err(), ShouldBeNil)
		So(tb.AddBox("box2", b2).Input("box1").Err(), ShouldBeNil)

		ti, err := tb.Build()
		So(err, ShouldBeNil)

		t := ti.(*defaultStaticTopology)

		Convey("When the state of the topology is TSStarting", func() {
			b1.init.block = true
			go func() {
				t.Run(ctx)
			}()
			t.Wait(ctx, TSStarting)
			Reset(func() {
				t.Stop(ctx)
			})

			// When t.State is TSStarting, t.Run blocks until t.State becomes TSRunning.
			// So, this test check the ordering of the process. Although the test could
			// get the expected order coincidentally, we can assume the logic is correct
			// if we wouldn't get the test failure for many times.
			ch := make(chan error)
			go func() {
				ch <- t.Run(ctx)
			}()
			go func() {
				ch <- fmt.Errorf("this is test specific error code")
				b1.ResumeInit()
			}()

			err1 := <-ch
			err2 := <-ch // This should be the t.Run's error

			Convey("Then the topology should fail to run again", func() {
				So(err2, ShouldNotBeNil)
			})

			Convey("Then the topology should be blocked until Box.Init resumes", func() {
				So(err1.Error(), ShouldEqual, "this is test specific error code")
			})
		})

		Convey("When the state of the topology is TSRunning", func() {
			b2.cnt = 0
			go func() {
				t.Run(ctx)
			}()
			t.Wait(ctx, TSRunning)
			Reset(func() {
				b2.EmitTuples(8)
				t.Stop(ctx)
			})

			Convey("Then the topology should fail to run again", func() {
				So(t.Run(ctx), ShouldNotBeNil)
			})
		})

		Convey("When the state of the topology is TSStopping", func() {
			b2.cnt = 0
			go func() {
				t.Run(ctx)
			}()
			go func() {
				t.Wait(ctx, TSRunning)
				t.Stop(ctx)
			}()
			t.Wait(ctx, TSStopping)
			Reset(func() {
				b2.EmitTuples(8)
				// t will stop by the goroutine above
			})

			Convey("Then the topology should fail to run again", func() {
				So(t.Run(ctx), ShouldNotBeNil)
			})
		})

		Convey("When the state of the topology is TSStopped", func() {
			go func() {
				t.Run(ctx)
			}()
			t.Wait(ctx, TSRunning)
			t.Stop(ctx)

			SkipConvey("Then the topology should fail to run again", func() {
				So(t.Run(ctx), ShouldNotBeNil)
			})
		})
	})
}

func TestDefaultStaticTopologyRunAndInit(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a default static topology", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DoesNothingSource{}
		// assuming Box.Init won't be called concurrently
		sc := &stubInitTerminateBoxSharedConfig{}
		b1 := NewStubInitTerminateBox(&DoesNothingBox{}, sc)
		b2 := NewStubInitTerminateBox(&DoesNothingBox{}, sc)
		b3 := NewStubInitTerminateBox(&DoesNothingBox{}, sc)
		bs := []*stubInitTerminateBox{b1, b2, b3}

		// Sink isn't necessary

		So(tb.AddSource("source", s).Err(), ShouldBeNil)
		So(tb.AddBox("box1", b1).Input("source").Err(), ShouldBeNil)
		So(tb.AddBox("box2", b2).Input("box1").Err(), ShouldBeNil)
		So(tb.AddBox("box3", b3).Input("box2").Err(), ShouldBeNil)
		t, err := tb.Build()
		So(err, ShouldBeNil)

		boxInitFailTestHelper := func(nth int) {
			fixtures := []struct {
				when    string
				initMsg string
			}{
				{"first", "no Box was"},
				{"second", "one Box was"},
				{"third", "two Boxes were"},
			}

			f := fixtures[nth]
			Convey(fmt.Sprintf("When the %v Box's Init fails", f.when), func() {
				sc.initFailAt = nth + 1
				err := t.Run(ctx)

				Convey("Then the topology shouldn't run", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("Then failed Box's Terminate shouldn't be called", func() {
					for _, b := range bs {
						if b.init.failed {
							So(b.terminateCnt, ShouldEqual, 0)
						}
					}
				})

				Convey("Then other boxes' Terminate should be called if its Init was called", func() {
					cnt := 0
					for _, b := range bs {
						if b.init.called && !b.init.failed {
							cnt++
							So(b.terminateCnt, ShouldEqual, 1)
						} else {
							So(b.terminateCnt, ShouldEqual, 0)
						}
					}

					Convey(fmt.Sprintf("And %v actually initialized", f.initMsg), func() {
						So(cnt, ShouldEqual, nth)
					})
				})
			})
		}
		boxInitFailTestHelper(0)
		boxInitFailTestHelper(1)
		boxInitFailTestHelper(2)

		Convey("When the third Init fails and the first Terminate fails", func() {
			sc.initFailAt = 3
			sc.terminateFailAt = 1
			err := t.Run(ctx)

			Convey("Then the topology shouldn't run", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then there should be two terminated boxes", func() {
				cnt := 0
				for _, b := range bs {
					if b.terminateCnt == 1 {
						cnt++
					}
				}
				So(cnt, ShouldEqual, 2)
			})
		})

		Convey("When the third Init fails and the first Terminate panics", func() {
			sc.initFailAt = 3
			sc.terminatePanicAt = 1
			err := t.Run(ctx)

			Convey("Then the topology shouldn't run", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then there should be two terminated boxes", func() {
				cnt := 0
				for _, b := range bs {
					if b.terminateCnt == 1 {
						cnt++
					}
				}
				So(cnt, ShouldEqual, 2)
			})
		})
	})
}

// TODO: Add a test case for Stop
// TODO: Add a test case for Source failures
// TODO: Add a test case for Box failures

package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
)

func BenchmarkDynamicPipe(b *testing.B) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)
	r, s := newDynamicPipe("test", 1024)
	go func() {
		for _ = range r.in {
		}
	}()

	t := &tuple.Tuple{}
	t.Data = tuple.Map{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Write(ctx, t)
		}
	})
}

func TestDynamicPipe(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a dynamic pipe", t, func() {
		// Use small capacity to check the sender never blocks.
		r, s := newDynamicPipe("test", 1)
		t := &tuple.Tuple{
			InputName: "hoge",
			Data: tuple.Map{
				"v": tuple.Int(1),
			},
		}

		Convey("When sending a tuple via the sender", func() {

			So(s.Write(ctx, t), ShouldBeNil)

			Convey("Then the tuple should be received by the receiver", func() {
				rt := <-r.in

				Convey("And its value should be correct", func() {
					So(rt.Data["v"], ShouldEqual, tuple.Int(1))
				})

				Convey("And its input name should be overwritten", func() {
					So(rt.InputName, ShouldEqual, "test")
				})
			})
		})

		Convey("When closing the pipe via the sender", func() {
			s.Close(ctx)

			Convey("Then it cannot no longer write a tuple", func() {
				So(s.Write(ctx, t), ShouldPointTo, errPipeClosed)
			})

			Convey("Then it can be closed again although it's not a part of the specification", func() {
				// This is only for improving the coverage.
				So(func() {
					s.Close(ctx)
				}, ShouldNotPanic)
			})
		})

		Convey("When closing the pipe via the receiver", func() {
			r.close()

			Convey("Then the sender should eventually be unable to write anymore tuple", func() {
				var err error
				for {
					// It can take many times until it will stop.
					err = s.Write(ctx, t)
					if err != nil {
						break
					}
				}
				So(err, ShouldPointTo, errPipeClosed)
			})
		})
	})
}

func TestDynamicDataSources(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given an empty data source", t, func() {
		srcs := newDynamicDataSources("test_component")

		Convey("When stopping it before starting to pour tuples", func() {
			srcs.stop(ctx)

			Convey("Then it pouring should fail", func() {
				si := NewTupleCollectorSink()
				So(srcs.pour(ctx, si, 1), ShouldNotBeNil)
			})
		})
	})

	Convey("Given an empty data source", t, func() {
		srcs := newDynamicDataSources("test_component")
		si := NewTupleCollectorSink()

		t := &tuple.Tuple{
			InputName: "some_component",
			Data: tuple.Map{
				"v": tuple.Int(1),
			},
		}

		stopped := make(chan error, 1)
		go func() {
			stopped <- srcs.pour(ctx, si, 4)
		}()
		Reset(func() {
			srcs.stop(ctx)
		})
		srcs.state.Wait(TSRunning)

		Convey("When starting it without any input", func() {
			Convey("Then it should start pouring", func() {
				// Just check if this test isn't dead-locked.
				So(true, ShouldBeTrue)
			})

			Convey("Then the sink should receive anything", func() {
				srcs.stop(ctx)
				<-stopped
				So(si.Tuples, ShouldBeEmpty)
			})
		})

		Convey("When adding an input after starting pouring and write a tuple", func() {
			r, s := newDynamicPipe("test1", 1)
			So(srcs.add("test_node_1", r), ShouldBeNil)
			So(s.Write(ctx, t), ShouldBeNil)
			s.close()
			srcs.stop(ctx)
			<-stopped

			Convey("Then the sink receive the tuple", func() {
				So(len(si.Tuples), ShouldEqual, 1)
			})
		})
	})

	Convey("Given a data source having destinations", t, func() {
		srcs := newDynamicDataSources("test_component")
		dsts := make([]*dynamicPipeSender, 2)
		for i := range dsts {
			r, s := newDynamicPipe(fmt.Sprint("test", i+1), 1)
			srcs.add(fmt.Sprint("test_node_", i+1), r)
			dsts[i] = s
		}
		Reset(func() {
			for _, d := range dsts {
				d.close() // safe to call multiple times
			}
		})
		si := NewTupleCollectorSink()

		t := &tuple.Tuple{
			InputName: "some_component",
			Data: tuple.Map{
				"v": tuple.Int(1),
			},
		}

		stopped := make(chan error, 1)
		go func() {
			stopped <- srcs.pour(ctx, si, 4)
		}()
		Reset(func() {
			srcs.stop(ctx)
		})
		srcs.state.Wait(TSRunning)

		Convey("When starting it again", func() {
			err := srcs.pour(ctx, si, 4)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When sending tuples from one source", func() {
			for i := 0; i < 5; i++ {
				So(dsts[0].Write(ctx, t), ShouldBeNil)
			}
			srcs.stop(ctx)
			So(<-stopped, ShouldBeNil)

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 5)
			})
		})

		Convey("When sending tuples from two sources", func() {
			for i := 0; i < 5; i++ {
				So(dsts[0].Write(ctx, t), ShouldBeNil)
				So(dsts[1].Write(ctx, t), ShouldBeNil)
			}
			srcs.stop(ctx)
			So(<-stopped, ShouldBeNil)

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 10)
			})
		})

		Convey("When stopping sources", func() {
			srcs.stop(ctx)

			Convey("Then it should eventually stop", func() {
				for _, d := range dsts {
					So(d.Write(ctx, t), ShouldPointTo, errPipeClosed)
				}
				So(<-stopped, ShouldBeNil)
			})
		})

		Convey("When stopping all inputs", func() {
			for _, d := range dsts {
				d.close()
			}
			Reset(func() {
				srcs.stop(ctx)
			})

			Convey("Then it shouldn't stop pouring", func() {
				// This actually doesn't guarantee anything, but test should
				// sometimes fail if the code above stopped the source.
				received := false
				select {
				case <-stopped:
					received = true
				default:
				}
				So(received, ShouldBeFalse)
			})
		})

		Convey("When stopping inputs after sending some tuples", func() {
			for i := 0; i < 5; i++ {
				So(dsts[0].Write(ctx, t), ShouldBeNil)
			}

			for i := 0; i < 3; i++ {
				So(dsts[1].Write(ctx, t), ShouldBeNil)
			}
			dsts[0].close()
			dsts[1].close()
			srcs.stop(ctx)
			So(<-stopped, ShouldBeNil)

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When adding a new input and sending a tuple", func() {
			r, s := newDynamicPipe("test3", 1)
			srcs.add("test_node_3", r)
			So(s.Write(ctx, t), ShouldBeNil)
			srcs.stop(ctx)
			So(<-stopped, ShouldBeNil)

			Convey("Then the sink should receive the tuple", func() {
				So(len(si.Tuples), ShouldEqual, 1)
			})
		})

		Convey("When adding a new input with the duplicated name", func() {
			r, _ := newDynamicPipe("test3", 1)
			err := srcs.add("test_node_1", r)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When remove an input after sending a tuple", func() {
			So(dsts[0].Write(ctx, t), ShouldBeNil)
			srcs.enableGracefulStop()
			srcs.remove("test_node_1")

			Convey("Then the input should be closed", func() {
				So(dsts[0].Write(ctx, t), ShouldPointTo, errPipeClosed)
			})

			Convey("Then the sink should receive it", func() {
				srcs.stop(ctx)
				So(<-stopped, ShouldBeNil)
				So(len(si.Tuples), ShouldEqual, 1)
			})

			Convey("Then the other input should still work", func() {
				So(dsts[1].Write(ctx, t), ShouldBeNil)
				srcs.stop(ctx)
				So(<-stopped, ShouldBeNil)
				So(len(si.Tuples), ShouldEqual, 2)
			})
		})
	})
}

// TODO: add fail tests of dynamicDataSources

func TestDynamicDataDestinations(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given an empty data destination", t, func() {
		dsts := newDynamicDataDestinations("test_component")
		t := &tuple.Tuple{
			InputName: "test_component",
			Data: tuple.Map{
				"v": tuple.Int(1),
			},
		}

		Convey("When sending a tuple", func() {
			var err error
			So(func() {
				err = dsts.Write(ctx, t)
			}, ShouldNotPanic)

			Convey("Then it shouldn't fail", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given data destinations", t, func() {
		dsts := newDynamicDataDestinations("test_component")
		recvs := make([]*dynamicPipeReceiver, 2)
		for i := range recvs {
			r, s := newDynamicPipe(fmt.Sprint("test", i+1), 1)
			recvs[i] = r
			dsts.add(fmt.Sprint("test_node_", i+1), s)
		}
		t := &tuple.Tuple{
			InputName: "test_component",
			Data: tuple.Map{
				"v": tuple.Int(1),
			},
		}

		Convey("When sending a tuple", func() {
			So(dsts.Write(ctx, t), ShouldBeNil)

			Convey("Then all destinations should receive it", func() {
				t1, ok := <-recvs[0].in
				So(ok, ShouldBeTrue)
				t2, ok := <-recvs[1].in
				So(ok, ShouldBeTrue)

				Convey("And tuples should have the correct input name", func() {
					So(t1.InputName, ShouldEqual, "test1")
					So(t2.InputName, ShouldEqual, "test2")
				})
			})
		})

		Convey("When sending closing the destinations after sending a tuple", func() {
			So(dsts.Write(ctx, t), ShouldBeNil)
			So(dsts.Close(ctx), ShouldBeNil)

			Convey("Then all receiver should receive a closing signal after the tuple", func() {
				for _, r := range recvs {
					_, ok := <-r.in
					So(ok, ShouldBeTrue)
					_, ok = <-r.in
					So(ok, ShouldBeFalse)
				}
			})
		})

		Convey("When one destination is closed by the receiver side", func() {
			recvs[0].close()
			Reset(func() {
				dsts.Close(ctx)
			})

			Convey("Then the destination receiver should eventually be removed", func() {
				go func() {
					for _ = range recvs[1].in {
					}
				}()
				for {
					if _, ok := dsts.dsts["test_node_1"]; !ok {
						break
					}
					dsts.Write(ctx, t)
				}

				_, ok := <-recvs[0].in
				So(ok, ShouldBeFalse)
			})
		})

		Convey("When adding a new destination after sending a tuple", func() {
			for _, r := range recvs {
				r := r
				go func() {
					for _ = range r.in {
					}
				}()
			}
			So(dsts.Write(ctx, t), ShouldBeNil)

			r, s := newDynamicPipe("test3", 1)
			So(dsts.add("test_node_3", s), ShouldBeNil)
			Reset(func() {
				dsts.Close(ctx)
			})

			Convey("Then the new receiver shouldn't receive the first tuple", func() {
				recved := false
				select {
				case <-r.in:
					recved = true
				default:
				}
				So(recved, ShouldBeFalse)
			})

			Convey("Then the new receiver should receive a new tuple", func() {
				So(dsts.Write(ctx, t), ShouldBeNil)
				_, ok := <-r.in
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When adding a destination with the duplicated name", func() {
			_, s := newDynamicPipe("hoge", 1)
			err := dsts.add("test_node_1", s)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When removing a destination", func() {
			dsts.remove("test_node_1")

			Convey("Then the destination should be closed", func() {
				_, ok := <-recvs[0].in
				So(ok, ShouldBeFalse)
			})

			Convey("Then Write should work", func() {
				So(dsts.Write(ctx, t), ShouldBeNil)
				_, ok := <-recvs[1].in
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When removing a destination after sending a tuple", func() {
			go func() {
				for _ = range recvs[1].in {
				}
			}()
			Reset(func() {
				dsts.Close(ctx)
			})
			So(dsts.Write(ctx, t), ShouldBeNil)
			dsts.remove("test_node_1")

			Convey("Then the destination should be able to receive the tuple", func() {
				_, ok := <-recvs[0].in
				So(ok, ShouldBeTrue)
				_, ok = <-recvs[0].in
				So(ok, ShouldBeFalse)
			})
		})

		Convey("When removing a nonexistent destination", func() {
			Convey("Then it shouldn't panic", func() {
				So(func() {
					dsts.remove("test_node_100")
				}, ShouldNotPanic)
			})
		})
	})
}

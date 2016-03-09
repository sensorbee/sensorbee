package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"sync/atomic"
	"testing"
	"time"
)

func TestNodeStatus(t *testing.T) {
	Convey("Given a topology having nodes", t, func() {
		ctx := NewContext(nil)
		t, err := NewDefaultTopology(ctx, "test")
		So(err, ShouldBeNil)
		Reset(func() {
			t.Stop()
		})

		so := NewTupleIncrementalEmitterSource(freshTuples())
		son, err := t.AddSource("source", so, nil)
		So(err, ShouldBeNil)
		so.EmitTuples(2) // send before a box is connected
		waitForNumDropped(son, 2)

		bn, err := t.AddBox("box", BoxFunc(forwardBox), nil)
		So(err, ShouldBeNil)
		So(bn.Input("source", nil), ShouldBeNil)
		bn.StopOnDisconnect(Inbound)
		bn.State().Wait(TSRunning)
		so.EmitTuples(1) // send before a sink is connected
		waitForNumDropped(bn, 1)

		si := NewTupleCollectorSink()
		sin, err := t.AddSink("sink", si, nil)
		So(err, ShouldBeNil)
		So(sin.Input("box", &SinkInputConfig{Capacity: 16}), ShouldBeNil)
		sin.StopOnDisconnect()
		so.EmitTuples(3)
		si.Wait(3)

		son.StopOnDisconnect()

		Convey("When getting status of the source while it's still running", func() {
			st := son.Status()

			Convey("Then it should have the running state", func() {
				So(st["state"], ShouldEqual, "running")
			})

			Convey("Then it should have no error", func() {
				// cannot use ShouldBeBlank because data.String isn't a standard string
				So(st["error"], ShouldBeNil)
			})

			Convey("Then it should have output_stats", func() {
				So(st["output_stats"], ShouldNotBeNil)
				os := st["output_stats"].(data.Map)

				Convey("And it should have the number of tuples sent from the source", func() {
					So(os["num_sent_total"], ShouldEqual, 4)
				})

				Convey("And it should have the number of dropped tuples", func() {
					So(os["num_dropped"], ShouldEqual, 2)
				})

				Convey("And it should have the statuses of connected nodes", func() {
					So(os["outputs"], ShouldNotBeNil)
					ns := os["outputs"].(data.Map)

					So(len(ns), ShouldEqual, 1)
					So(ns["box"], ShouldNotBeNil)

					b := ns["box"].(data.Map)
					So(b["num_sent"], ShouldEqual, 4)
					So(b["queue_size"], ShouldBeGreaterThan, 0)
					So(b["num_queued"], ShouldEqual, 0)
				})
			})

			Convey("Then it should have the status of the source implementation", func() {
				So(st["source"], ShouldNotBeNil)
				v, _ := st.Get(data.MustCompilePath("source.test"))
				So(v, ShouldEqual, "test")
			})

			Convey("Then it should have its behavior descriptions", func() {
				So(st["behaviors"], ShouldNotBeNil)
				bs := st["behaviors"].(data.Map)

				Convey("And stop_on_disconnect should be true", func() {
					So(bs["stop_on_disconnect"], ShouldEqual, data.True)
				})

				Convey("And remove_on_stop should be false", func() {
					So(bs["remove_on_stop"], ShouldEqual, data.False)
				})

				Convey("And remove_on_stop should be true after enabling it", func() {
					son.RemoveOnStop()
					st := son.Status()
					v, _ := st.Get(data.MustCompilePath("behaviors.remove_on_stop"))
					So(v, ShouldEqual, data.True)
				})
			})
		})

		Convey("When getting status of the source after the source is stopped", func() {
			so.EmitTuples(2)
			son.State().Wait(TSStopped)

			st := son.Status()

			Convey("Then it should have the stopped state", func() {
				So(st["state"], ShouldEqual, "stopped")
			})

			Convey("Then it should have no error", func() {
				// cannot use ShouldBeBlank because data.String isn't a standard string
				So(st["error"], ShouldBeNil)
			})

			Convey("Then it should have output_stats", func() {
				So(st["output_stats"], ShouldNotBeNil)
				os := st["output_stats"].(data.Map)

				Convey("And it should have the number of tuples sent from the source", func() {
					So(os["num_sent_total"], ShouldEqual, 6)
				})

				Convey("And it should have the number of dropped tuples", func() {
					So(os["num_dropped"], ShouldEqual, 2)
				})

				Convey("And it shouldn't have any connections", func() {
					So(os["outputs"], ShouldNotBeNil)
					ns := os["outputs"].(data.Map)
					So(ns, ShouldBeEmpty)
				})
			})

			Convey("Then it should have the status of the source implementation", func() {
				So(st["source"], ShouldNotBeNil)
				v, _ := st.Get(data.MustCompilePath("source.test"))
				So(v, ShouldEqual, "test")
			})
		})

		Convey("When getting status of the box while it's still running", func() {
			st := bn.Status()

			Convey("Then it should have the running state", func() {
				So(st["state"], ShouldEqual, "running")
			})

			Convey("Then it should have no error", func() {
				So(st["error"], ShouldBeNil)
			})

			Convey("Then it should have input_stats", func() {
				So(st["input_stats"], ShouldNotBeNil)
				is := st["input_stats"].(data.Map)

				Convey("And it should have the number of tuples received", func() {
					So(is["num_received_total"], ShouldEqual, 4)
				})

				Convey("And it should have the number of errors", func() {
					So(is["num_errors"], ShouldEqual, 0)
				})

				Convey("And it should have the statuses of connected nodes", func() {
					So(is["inputs"], ShouldNotBeNil)
					ns := is["inputs"].(data.Map)

					So(len(ns), ShouldEqual, 1)
					So(ns["source"], ShouldNotBeNil)

					s := ns["source"].(data.Map)
					So(s["num_received"], ShouldEqual, 4)
					So(s["queue_size"], ShouldBeGreaterThan, 0)
					So(s["num_queued"], ShouldEqual, 0)
				})
			})

			Convey("Then it should have output_stats", func() {
				So(st["output_stats"], ShouldNotBeNil)
				os := st["output_stats"].(data.Map)

				Convey("And it should have the number of tuples sent from the source", func() {
					So(os["num_sent_total"], ShouldEqual, 3)
				})

				Convey("And it should have the number of dropped tuples", func() {
					So(os["num_dropped"], ShouldEqual, 1)
				})

				Convey("And it should have the statuses of connected nodes", func() {
					So(os["outputs"], ShouldNotBeNil)
					ns := os["outputs"].(data.Map)

					So(len(ns), ShouldEqual, 1)
					So(ns["sink"], ShouldNotBeNil)

					b := ns["sink"].(data.Map)
					So(b["num_sent"], ShouldEqual, 3)
					So(b["queue_size"], ShouldEqual, 16)
					So(b["num_queued"], ShouldEqual, 0)
				})
			})

			Convey("Then it should have its behavior descriptions", func() {
				So(st["behaviors"], ShouldNotBeNil)
				bs := st["behaviors"].(data.Map)

				Convey("And stop_on_inbound_disconnect should be true", func() {
					So(bs["stop_on_inbound_disconnect"], ShouldEqual, data.True)
				})

				Convey("And stop_on_outbound_disconnect should be false", func() {
					So(bs["stop_on_outbound_disconnect"], ShouldEqual, data.False)
				})

				Convey("And graceful_stop shouldn't be enabled", func() {
					So(bs["graceful_stop"], ShouldEqual, data.False)
				})

				Convey("And graceful_stop should be true after enabling it", func() {
					bn.EnableGracefulStop()
					st := bn.Status()
					v, _ := st.Get(data.MustCompilePath("behaviors.graceful_stop"))
					So(v, ShouldEqual, data.True)
				})

				Convey("And remove_on_stop should be false", func() {
					So(bs["remove_on_stop"], ShouldEqual, data.False)
				})

				Convey("And remove_on_stop should be true after enabling it", func() {
					bn.RemoveOnStop()
					st := bn.Status()
					v, _ := st.Get(data.MustCompilePath("behaviors.remove_on_stop"))
					So(v, ShouldEqual, data.True)
				})
			})

			// TODO: check st["box"]
		})

		Convey("When getting status of the box after the box is stopped", func() {
			so.EmitTuples(2)
			si.Wait(5)
			bn.State().Wait(TSStopped)

			st := bn.Status()

			Convey("Then it should have the stopped state", func() {
				So(st["state"], ShouldEqual, "stopped")
			})

			Convey("Then it should have no error", func() {
				So(st["error"], ShouldBeNil)
			})

			Convey("Then it should have input_stats", func() {
				So(st["input_stats"], ShouldNotBeNil)
				is := st["input_stats"].(data.Map)

				Convey("And it should have the number of tuples received", func() {
					So(is["num_received_total"], ShouldEqual, 6)
				})

				Convey("And it should have the number of errors", func() {
					So(is["num_errors"], ShouldEqual, 0)
				})

				Convey("And it should have no connected nodes", func() {
					So(is["inputs"], ShouldNotBeNil)
					ns := is["inputs"].(data.Map)
					So(ns, ShouldBeEmpty)
				})
			})

			Convey("Then it should have output_stats", func() {
				So(st["output_stats"], ShouldNotBeNil)
				os := st["output_stats"].(data.Map)

				Convey("And it should have the number of tuples sent from the source", func() {
					So(os["num_sent_total"], ShouldEqual, 5)
				})

				Convey("And it should have the number of dropped tuples", func() {
					So(os["num_dropped"], ShouldEqual, 1)
				})

				Convey("And it should have no connected nodes", func() {
					So(os["outputs"], ShouldNotBeNil)
					ns := os["outputs"].(data.Map)
					So(ns, ShouldBeEmpty)
				})
			})

			// TODO: st["box"]
		})

		Convey("When getting status of the sink while it's still running", func() {
			st := sin.Status()

			Convey("Then it should have the running state", func() {
				So(st["state"], ShouldEqual, "running")
			})

			Convey("Then it should have no error", func() {
				So(st["error"], ShouldBeNil)
			})

			Convey("Then it should have input_stats", func() {
				So(st["input_stats"], ShouldNotBeNil)
				is := st["input_stats"].(data.Map)

				Convey("And it should have the number of tuples received", func() {
					So(is["num_received_total"], ShouldEqual, 3)
				})

				Convey("And it should have the number of errors", func() {
					So(is["num_errors"], ShouldEqual, 0)
				})

				Convey("And it should have the statuses of connected nodes", func() {
					So(is["inputs"], ShouldNotBeNil)
					ns := is["inputs"].(data.Map)

					So(len(ns), ShouldEqual, 1)
					So(ns["box"], ShouldNotBeNil)

					s := ns["box"].(data.Map)
					So(s["num_received"], ShouldEqual, 3)
					So(s["queue_size"], ShouldEqual, 16)
					So(s["num_queued"], ShouldEqual, 0)
				})
			})

			Convey("Then it should have its behavior descriptions", func() {
				So(st["behaviors"], ShouldNotBeNil)
				bs := st["behaviors"].(data.Map)

				Convey("And stop_on_disconnect should be true", func() {
					So(bs["stop_on_disconnect"], ShouldEqual, data.True)
				})

				Convey("And graceful_stop shouldn't be enabled", func() {
					So(bs["graceful_stop"], ShouldEqual, data.False)
				})

				Convey("And graceful_stop should be true after enabling it", func() {
					sin.EnableGracefulStop()
					st := sin.Status()
					v, _ := st.Get(data.MustCompilePath("behaviors.graceful_stop"))
					So(v, ShouldEqual, data.True)
				})

				Convey("And remove_on_stop should be false", func() {
					So(bs["remove_on_stop"], ShouldEqual, data.False)
				})

				Convey("And remove_on_stop should be true after enabling it", func() {
					sin.RemoveOnStop()
					st := sin.Status()
					v, _ := st.Get(data.MustCompilePath("behaviors.remove_on_stop"))
					So(v, ShouldEqual, data.True)
				})
			})

			// TODO: st["sink"]
		})

		Convey("When getting status of the sink after the sink is stopped", func() {
			so.EmitTuples(2)
			si.Wait(5)
			sin.State().Wait(TSStopped)

			st := sin.Status()

			Convey("Then it should have the stopped state", func() {
				So(st["state"], ShouldEqual, "stopped")
			})

			Convey("Then it should have no error", func() {
				So(st["error"], ShouldBeNil)
			})

			Convey("Then it should have input_stats", func() {
				So(st["input_stats"], ShouldNotBeNil)
				is := st["input_stats"].(data.Map)

				Convey("And it should have the number of tuples received", func() {
					So(is["num_received_total"], ShouldEqual, 5)
				})

				Convey("And it should have the number of errors", func() {
					So(is["num_errors"], ShouldEqual, 0)
				})

				Convey("And it should have no connected nodes", func() {
					So(is["inputs"], ShouldNotBeNil)
					ns := is["inputs"].(data.Map)
					So(ns, ShouldBeEmpty)
				})
			})

			// TODO: st["sink"]
		})
	})
}

// TODO: test run failures
// TODO: test Write failures of Boxes and Sinks

func waitForNumDropped(node Node, n int64) {
	var dsts *dataDestinations
	switch t := node.(type) {
	case *defaultSourceNode:
		dsts = t.dsts
	case *defaultBoxNode:
		dsts = t.dsts
	}
	for n > atomic.LoadInt64(&dsts.numDropped) {
		time.Sleep(time.Nanosecond)
	}
}

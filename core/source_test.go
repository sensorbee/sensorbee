package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestRewindableSource(t *testing.T) {
	Convey("Given a default topology", t, func() {
		t := NewDefaultTopology(NewContext(nil), "dt1")
		Reset(func() {
			t.Stop()
		})

		so := NewTupleEmitterSource(freshTuples())
		son, err := t.AddSource("source", NewRewindableSource(so), &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		b := &BlockingForwardBox{cnt: 1000}
		bn, err := t.AddBox("box", b, nil)
		So(err, ShouldBeNil)
		So(bn.Input("source", &BoxInputConfig{
			Capacity: 1, // (almost) blocking channel
		}), ShouldBeNil)

		si := NewTupleCollectorSink()
		sin, err := t.AddSink("sink", si, nil)
		So(sin.Input("box", nil), ShouldBeNil)

		Convey("When emitting all tuples", func() {
			So(son.Resume(), ShouldBeNil)
			si.Wait(8)

			Convey("Then the source shouldn't stop", func() {
				So(son.State().Get(), ShouldEqual, TSRunning)
			})

			Convey("Then status should show that it's waiting for rewind", func() {
				st := son.Status()
				v, _ := st.Get("source.waiting_for_rewind")
				So(v, ShouldEqual, data.True)
			})
		})

		Convey("When rewinding before sending any tuple", func() {
			So(son.Rewind(), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)
			si.Wait(8)

			Convey("Then the sink should receive all tuples", func() {
				So(len(si.Tuples), ShouldEqual, 8)
			})
		})

		Convey("When rewinding after sending some tuples", func() {
			// emit 2 tuples and then blocks. So, 2 tuples went to the sink,
			// 1 tuple is blocked in the box, and 2 tuples is blocked in the
			// channel between the source and box because its capacity is 1
			// (1 tuple is in the queue and the other one is blocked at sending
			// operation). So, 5 tuples in total were emitted from the source.
			b.cnt = 2
			So(son.Resume(), ShouldBeNil)
			si.Wait(2)
			So(son.Pause(), ShouldBeNil)
			b.EmitTuples(1000)
			So(son.Rewind(), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then all tuple should be able to be sent again", func() {
				si.Wait(12)

				// Due to concurrency, the number of tuples arriving to the sink
				// can be either 12 or 13. It could be 11 but very rare.
				So(len(si.Tuples), ShouldBeGreaterThanOrEqualTo, 12)
				So(len(si.Tuples), ShouldBeLessThanOrEqualTo, 13)
			})
		})

		Convey("When sending all tuples", func() {
			So(son.Resume(), ShouldBeNil)
			si.Wait(8)

			Convey("The source should be able to be rewound", func() {
				So(son.Rewind(), ShouldBeNil)
				si.Wait(16)
				So(len(si.Tuples), ShouldEqual, 16)
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
				v, _ := st.Get("source.waiting_for_rewind")
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

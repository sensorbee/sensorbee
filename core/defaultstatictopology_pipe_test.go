package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
	"time"
)

func freshTuples() []*tuple.Tuple {
	tup1 := tuple.Tuple{
		Data: tuple.Map{
			"seq": tuple.Int(1),
		},
		InputName: "input",
	}
	tup2 := tuple.Tuple{
		Data: tuple.Map{
			"seq": tuple.Int(2),
		},
		InputName: "input",
	}
	tup3 := tuple.Tuple{
		Data: tuple.Map{
			"seq": tuple.Int(3),
		},
		InputName: "input",
	}
	tup4 := tuple.Tuple{
		Data: tuple.Map{
			"seq": tuple.Int(4),
		},
		InputName: "input",
	}
	return []*tuple.Tuple{&tup1, &tup2, &tup3, &tup4}
}

func TestCapacityPipe(t *testing.T) {

	Convey("Given a simple source/slow box/very slow box/sink topology", t, func() {
		/*
		 *   so -*--> b1 -*--> b2 -*--> si
		 */

		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			// we need to get new tuples here or the same objects
			// will be used in every run and yield wrong trace
			// information
			Tuples: freshTuples(),
		}
		tb.AddSource("source", so)

		b1 := BoxFunc(slowForwardBox)
		tb.AddBox("box1", b1).Input("source")

		b2 := BoxFunc(verySlowForwardBox)
		tb.AddBox("box2", b2).Input("box1")

		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("box2")

		t := tb.Build()

		Convey("When four tuples are emitted by the source", func() {
			t.Run(&Context{})
			// wait until tuples arrived
			loopTimeout := 30
			i := 0
			for i = 0; i < loopTimeout; i++ {
				if len(si.Tuples) >= len(so.Tuples) {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			So(i, ShouldBeLessThan, loopTimeout)

			/* Processing should happen as follows:
			 *
			 * t ------------------------------------------------>
			 * so: | t1 | t2    | t3      | t4      |
			 * b1:      | t1 .. | t2 ..   | t3 ..   | t4 ..   |
			 * b2:              | t1 .... | t2 .... | t3 .... | t4 ....
			 * si:                        | t1;     | t2;     | t3;
			 */

			/*
			 That is,
			 - at the first pipe, the waiting time is
			   - about 0 for the first tuple (is handed over immediately)
			   - about 100ms for the second item (needs to wait for
			     b1-processing of first tuple to finish)
			   - about 300ms for every other item (needs to wait for
			     b2-processing of (n-2)-th item to finish)
			 - at the second pipe, the waiting time is
			   - about 0 for the first tuple (is handed over immediately)
			   - about 200ms for every other item (b1-processing of n-th
			   	 tuple takes 100ms, b2-processing of (n-1)-th tuple takes
			   	 300ms, so we need to wait 200ms)
			 - at the third pipe, the waiting time is about 0 every time
			*/

			So(si.Tuples, ShouldNotBeNil)
			So(len(si.Tuples), ShouldEqual, 4)

			So(len(si.Tuples[0].Trace), ShouldEqual, 9) // OUT-OTHER-IN-OUT-OTHER-IN-OUT-OTHER-IN

			Convey("Then waiting time at the intermediate pipes should match", func() {
				// first pipe: 0/100/300
				{
					tup1Wait := si.Tuples[0].Trace[1].Timestamp.Sub(si.Tuples[0].Trace[0].Timestamp)
					tup2Wait := si.Tuples[1].Trace[1].Timestamp.Sub(si.Tuples[1].Trace[0].Timestamp)
					tup3Wait := si.Tuples[2].Trace[1].Timestamp.Sub(si.Tuples[2].Trace[0].Timestamp)
					tup4Wait := si.Tuples[3].Trace[1].Timestamp.Sub(si.Tuples[3].Trace[0].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 100*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
				}

				// second pipe: 0/200
				{
					tup1Wait := si.Tuples[0].Trace[4].Timestamp.Sub(si.Tuples[0].Trace[3].Timestamp)
					tup2Wait := si.Tuples[1].Trace[4].Timestamp.Sub(si.Tuples[1].Trace[3].Timestamp)
					tup3Wait := si.Tuples[2].Trace[4].Timestamp.Sub(si.Tuples[2].Trace[3].Timestamp)
					tup4Wait := si.Tuples[3].Trace[4].Timestamp.Sub(si.Tuples[3].Trace[3].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 200*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 200*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 200*time.Millisecond, time.Millisecond)
				}

				// third pipe: 0
				{
					tup1Wait := si.Tuples[0].Trace[7].Timestamp.Sub(si.Tuples[0].Trace[6].Timestamp)
					tup2Wait := si.Tuples[1].Trace[7].Timestamp.Sub(si.Tuples[1].Trace[6].Timestamp)
					tup3Wait := si.Tuples[2].Trace[7].Timestamp.Sub(si.Tuples[2].Trace[6].Timestamp)
					tup4Wait := si.Tuples[3].Trace[7].Timestamp.Sub(si.Tuples[3].Trace[6].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
				}
			})
		})
	})

	Convey("Given a simple source/very slow box/slow box/sink topology", t, func() {
		/*
		 *   so -*--> b1 -*--> b2 -*--> si
		 */

		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			// we need to get new tuples here or the same objects
			// will be used in every run and yield wrong trace
			// information
			Tuples: freshTuples(),
		}
		tb.AddSource("source", so)

		b1 := BoxFunc(verySlowForwardBox)
		tb.AddBox("box1", b1).Input("source")

		b2 := BoxFunc(slowForwardBox)
		tb.AddBox("box2", b2).Input("box1")

		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("box2")

		t := tb.Build()

		Convey("When four tuples are emitted by the source", func() {
			t.Run(&Context{})
			// wait until tuples arrived
			loopTimeout := 30
			i := 0
			for i = 0; i < loopTimeout; i++ {
				if len(si.Tuples) >= len(so.Tuples) {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			So(i, ShouldBeLessThan, loopTimeout)

			/* Processing should happen as follows:
			 *
			 * t ------------------------------------------------>
			 * so: | t1 | t2      | t3      | t4      |
			 * b1:      | t1 .... | t2 .... | t3 .... | t4 ....
			 * b2:                | t1 .. ; | t2 .. ; | t3 :: ;
			 * si:                        | t1;     | t2;     | t3;
			 */

			/*
			 That is,
			 - at the first pipe, the waiting time is
			   - about 0 for the first tuple (is handed over immediately)
			   - about 300ms for every other item (needs to wait for
			     b1-processing of (n-1)-th item to finish)
			 - at the second pipe, the waiting time is about 0 every time
			   (because b2 is faster than b1)
			 - at the third pipe, the waiting time is about 0 every time
			*/

			So(si.Tuples, ShouldNotBeNil)
			So(len(si.Tuples), ShouldEqual, 4)

			So(len(si.Tuples[0].Trace), ShouldEqual, 9) // OUT-OTHER-IN-OUT-OTHER-IN-OUT-OTHER-IN

			Convey("Then waiting time at the intermediate pipes should match", func() {
				// first pipe: 0/300
				{
					tup1Wait := si.Tuples[0].Trace[1].Timestamp.Sub(si.Tuples[0].Trace[0].Timestamp)
					tup2Wait := si.Tuples[1].Trace[1].Timestamp.Sub(si.Tuples[1].Trace[0].Timestamp)
					tup3Wait := si.Tuples[2].Trace[1].Timestamp.Sub(si.Tuples[2].Trace[0].Timestamp)
					tup4Wait := si.Tuples[3].Trace[1].Timestamp.Sub(si.Tuples[3].Trace[0].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
				}

				// second pipe: 0
				{
					tup1Wait := si.Tuples[0].Trace[4].Timestamp.Sub(si.Tuples[0].Trace[3].Timestamp)
					tup2Wait := si.Tuples[1].Trace[4].Timestamp.Sub(si.Tuples[1].Trace[3].Timestamp)
					tup3Wait := si.Tuples[2].Trace[4].Timestamp.Sub(si.Tuples[2].Trace[3].Timestamp)
					tup4Wait := si.Tuples[3].Trace[4].Timestamp.Sub(si.Tuples[3].Trace[3].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
				}

				// third pipe: 0
				{
					tup1Wait := si.Tuples[0].Trace[7].Timestamp.Sub(si.Tuples[0].Trace[6].Timestamp)
					tup2Wait := si.Tuples[1].Trace[7].Timestamp.Sub(si.Tuples[1].Trace[6].Timestamp)
					tup3Wait := si.Tuples[2].Trace[7].Timestamp.Sub(si.Tuples[2].Trace[6].Timestamp)
					tup4Wait := si.Tuples[3].Trace[7].Timestamp.Sub(si.Tuples[3].Trace[6].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
				}
			})
		})
	})
}

func slowForwardBox(t *tuple.Tuple, w Writer) error {
	time.Sleep(100 * time.Millisecond)
	w.Write(t)
	return nil
}

func verySlowForwardBox(t *tuple.Tuple, w Writer) error {
	time.Sleep(300 * time.Millisecond)
	w.Write(t)
	return nil
}

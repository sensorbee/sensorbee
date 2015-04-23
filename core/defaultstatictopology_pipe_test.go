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

	Convey("Given a simple source/box/sink topology with 2 boxes and 2 sinks", t, func() {
		/*
		 *        /--> b1 -*--> si1
		 *   so -*
		 *        \--> b2 -*--> si2
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: freshTuples(),
		}
		tb.AddSource("source", so)

		b1 := BoxFunc(slowForwardBox)
		tb.AddBox("box1", b1).Input("source")
		b2 := BoxFunc(verySlowForwardBox)
		tb.AddBox("box2", b2).Input("source")

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("box1")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("box2")

		t := tb.Build()

		Convey("When four tuples are emitted by the source", func() {
			t.Run(&Context{})
			// wait until tuples arrived
			loopTimeout := 30
			i := 0
			for i = 0; i < loopTimeout; i++ {
				if len(si1.Tuples) >= len(so.Tuples) &&
					len(si2.Tuples) >= len(so.Tuples) {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			So(i, ShouldBeLessThan, loopTimeout)

			/* Processing should happen as follows:
			 *
			 * t ------------------------------------------------>
			 * so:  | t1 | t2      | t3      | t4      |
			 * b1:       | t1 .. ; | t2 .. ; | t3 .. ; | t4
			 * si1:              | t1;     | t2;     | t3;
			 * b2:       | t1 .... | t2 .... | t3 .... | t4 ....
			 * si2:                | t1;     | t2;     | t3;
			 */

			/*
				 That is,
				 - at the so-pipe, the waiting time is
				   - about 0 for the first tuple (is handed over immediately)
				   - about 300ms for every other item (needs to wait for parallel
					b1-processing (100ms) and b2-processing (300ms) of
					 (n-1)-th item to finish)
				 - at the b1-pipe, the waiting time is about 0 every time
				 - at the b2-pipe, the waiting time is about 0 every time
			*/

			So(si1.Tuples, ShouldNotBeNil)
			So(len(si1.Tuples), ShouldEqual, 4)
			So(si2.Tuples, ShouldNotBeNil)
			So(len(si2.Tuples), ShouldEqual, 4)

			So(len(si1.Tuples[0].Trace), ShouldEqual, 6) // OUT-OTHER-IN-OUT-OTHER-IN
			So(len(si2.Tuples[0].Trace), ShouldEqual, 6) // OUT-OTHER-IN-OUT-OTHER-IN

			Convey("Then waiting time at the intermediate pipes should match", func() {
				// so-pipe: 0/300
				{
					// so-b1-si1 path
					tup1Wait := si1.Tuples[0].Trace[1].Timestamp.Sub(si1.Tuples[0].Trace[0].Timestamp)
					tup2Wait := si1.Tuples[1].Trace[1].Timestamp.Sub(si1.Tuples[1].Trace[0].Timestamp)
					tup3Wait := si1.Tuples[2].Trace[1].Timestamp.Sub(si1.Tuples[2].Trace[0].Timestamp)
					tup4Wait := si1.Tuples[3].Trace[1].Timestamp.Sub(si1.Tuples[3].Trace[0].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
				}
				{
					// so-b2-si2 path
					tup1Wait := si2.Tuples[0].Trace[1].Timestamp.Sub(si2.Tuples[0].Trace[0].Timestamp)
					tup2Wait := si2.Tuples[1].Trace[1].Timestamp.Sub(si2.Tuples[1].Trace[0].Timestamp)
					tup3Wait := si2.Tuples[2].Trace[1].Timestamp.Sub(si2.Tuples[2].Trace[0].Timestamp)
					tup4Wait := si2.Tuples[3].Trace[1].Timestamp.Sub(si2.Tuples[3].Trace[0].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 300*time.Millisecond, time.Millisecond)
				}

				// b1 pipe: 0
				{
					tup1Wait := si1.Tuples[0].Trace[4].Timestamp.Sub(si1.Tuples[0].Trace[3].Timestamp)
					tup2Wait := si1.Tuples[1].Trace[4].Timestamp.Sub(si1.Tuples[1].Trace[3].Timestamp)
					tup3Wait := si1.Tuples[2].Trace[4].Timestamp.Sub(si1.Tuples[2].Trace[3].Timestamp)
					tup4Wait := si1.Tuples[3].Trace[4].Timestamp.Sub(si1.Tuples[3].Trace[3].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
				}

				// b2 pipe: 0
				{
					tup1Wait := si2.Tuples[0].Trace[4].Timestamp.Sub(si2.Tuples[0].Trace[3].Timestamp)
					tup2Wait := si2.Tuples[1].Trace[4].Timestamp.Sub(si2.Tuples[1].Trace[3].Timestamp)
					tup3Wait := si2.Tuples[2].Trace[4].Timestamp.Sub(si2.Tuples[2].Trace[3].Timestamp)
					tup4Wait := si2.Tuples[3].Trace[4].Timestamp.Sub(si2.Tuples[3].Trace[3].Timestamp)
					So(tup1Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup3Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup4Wait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
				}
			})
		})
	})

	Convey("Given a simple source/box/sink topology with 2 sources", t, func() {
		/*
		 *   so1 -*-\
		 *           --> b -*--> si
		 *   so2 -*-/
		 */
		tb := NewDefaultStaticTopologyBuilder()

		so1 := &TupleEmitterSource{
			Tuples: freshTuples()[0:3],
		}
		tb.AddSource("source1", so1)

		so2 := &TupleEmitterSource{
			Tuples: freshTuples()[1:4],
		}
		tb.AddSource("source2", so2)

		b1 := BoxFunc(slowForwardBox)
		tb.AddBox("box", b1).
			Input("source1").
			Input("source2")
		si := &TupleCollectorSink{}
		tb.AddSink("si", si).Input("box")

		t := tb.Build()

		Convey("When three tuples are emitted by each source", func() {
			t.Run(&Context{})
			// wait until tuples arrived
			loopTimeout := 30
			i := 0
			for i = 0; i < loopTimeout; i++ {
				if len(si.Tuples) >= len(so1.Tuples)+len(so2.Tuples) {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			So(i, ShouldBeLessThan, loopTimeout)

			/* Processing should happen as follows:
			 *
			 * t ------------------------------------------------>
			 * so1:  | t1 | t2    | t3     |
			 * so2:  | t2 | t3    | t4     |
			 * b:         | t1 .. | t2 ..  | t3 ..  |
			 *            \ t2 .. | t3 ..  | t4 ..  |
			 * si:                | t1+t2; | t2+t3; | t3+t4;
			 *
			 * (b is processing multiple tuples at a time)
			 */

			/*
			 That is,
			 - at the so1-pipe, the waiting time is
			   - about 0 for the first tuple (is handed over immediately)
			   - about 100ms for every other tuple (needs to wait for
			     b-processing of (n-1)-th item to finish)
			 - at the so2-pipe, the waiting time is
			   - about 0 for the first tuple (is handed over immediately)
			   - about 100ms for every other tuple (needs to wait for
			     b-processing of (n-1)-th item to finish)
			 - at the b-pipe, the waiting time is about 0 every time
			*/

			So(si.Tuples, ShouldNotBeNil)
			So(len(si.Tuples), ShouldEqual, 6)

			So(len(si.Tuples[0].Trace), ShouldEqual, 6) // OUT-OTHER-IN-OUT-OTHER-IN

			Convey("Then waiting time at the intermediate pipes should match", func() {
				// so1-pipe, so2-pipe: 0/100
				{
					// we don't know whether so1's tuple or so2's tuple will arrive
					// first at the sink, but it does not matter for this check
					tup1aWait := si.Tuples[0].Trace[1].Timestamp.Sub(si.Tuples[0].Trace[0].Timestamp)
					tup1bWait := si.Tuples[1].Trace[1].Timestamp.Sub(si.Tuples[1].Trace[0].Timestamp)
					tup2aWait := si.Tuples[2].Trace[1].Timestamp.Sub(si.Tuples[2].Trace[0].Timestamp)
					tup2bWait := si.Tuples[3].Trace[1].Timestamp.Sub(si.Tuples[3].Trace[0].Timestamp)
					tup3aWait := si.Tuples[4].Trace[1].Timestamp.Sub(si.Tuples[4].Trace[0].Timestamp)
					tup3bWait := si.Tuples[5].Trace[1].Timestamp.Sub(si.Tuples[5].Trace[0].Timestamp)
					So(tup1aWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup1bWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2aWait, ShouldAlmostEqual, 100*time.Millisecond, time.Millisecond)
					So(tup2bWait, ShouldAlmostEqual, 100*time.Millisecond, time.Millisecond)
					So(tup3aWait, ShouldAlmostEqual, 100*time.Millisecond, time.Millisecond)
					So(tup3bWait, ShouldAlmostEqual, 100*time.Millisecond, time.Millisecond)
				}

				// b-pipe: 0/100
				{
					tup1aWait := si.Tuples[0].Trace[4].Timestamp.Sub(si.Tuples[0].Trace[3].Timestamp)
					tup1bWait := si.Tuples[1].Trace[4].Timestamp.Sub(si.Tuples[1].Trace[3].Timestamp)
					tup2aWait := si.Tuples[2].Trace[4].Timestamp.Sub(si.Tuples[2].Trace[3].Timestamp)
					tup2bWait := si.Tuples[3].Trace[4].Timestamp.Sub(si.Tuples[3].Trace[3].Timestamp)
					tup3aWait := si.Tuples[4].Trace[4].Timestamp.Sub(si.Tuples[4].Trace[3].Timestamp)
					tup3bWait := si.Tuples[5].Trace[4].Timestamp.Sub(si.Tuples[5].Trace[3].Timestamp)
					So(tup1aWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup1bWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2aWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup2bWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup3aWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
					So(tup3bWait, ShouldAlmostEqual, 0*time.Millisecond, time.Millisecond)
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

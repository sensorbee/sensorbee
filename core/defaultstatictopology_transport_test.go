package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
	"sync"
	"testing"
	"time"
)

// TestDefaultTopologyTupleTransport tests that tuples are correctly
// copied/transported in various topologies.
func TestDefaultTopologyTupleTransport(t *testing.T) {
	tup1 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(1),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	tup2 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(2),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 1, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 1, 0, time.UTC),
		BatchID:       7,
	}

	Convey("Given a simple source/sink topology", t, func() {
		/*
		 *   so -*--> si
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1, &tup2},
		}
		tb.AddSource("source", so)

		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("source")

		t := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(&Context{})
			Convey("Then the sink receives the same object", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 2)
				// pointers point to the same objects
				So(so.Tuples[0], ShouldPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
		})
	})

	Convey("Given a simple source/sink topology with 2 sinks", t, func() {
		/*
		 *        /--> si1
		 *   so -*
		 *        \--> si2
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1, &tup2},
		}
		tb.AddSource("source", so)

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("source")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("source")

		t := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(&Context{})
			Convey("Then the sink 1 receives a copy", func() {
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 2)
				// contents are the same
				si := si1
				So(so.Tuples[0].Data, ShouldResemble, si.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si.Tuples[1].BatchID)
				// source has two received sinks, so tuples are copied
				So(so.Tuples[0], ShouldNotPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
			Convey("And the sink 2 receives a copy", func() {
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 2)
				// contents are the same
				si := si2
				So(so.Tuples[0].Data, ShouldResemble, si.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si.Tuples[1].BatchID)
				// pointers point to different objects
				So(so.Tuples[0], ShouldNotPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
			Convey("And the traces of tuples differ", func() {
				So(len(si1.Tuples), ShouldEqual, 2)
				So(len(si2.Tuples), ShouldEqual, 2)
				So(si1.Tuples[0].Trace, ShouldNotResemble, si2.Tuples[0].Trace)
				So(si1.Tuples[1].Trace, ShouldNotResemble, si2.Tuples[1].Trace)
			})
		})
	})

	Convey("Given a simple source/box/sink topology", t, func() {
		/*
		 *   so -*--> b -*--> si
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1, &tup2},
		}
		tb.AddSource("source", so)

		b := BoxFunc(forwardBox)
		tb.AddBox("box", b).Input("source")

		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("box")

		t := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(&Context{})
			Convey("Then the sink receives the same object", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 2)
				// pointers point to the same objects
				So(so.Tuples[0], ShouldPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
		})
	})

	Convey("Given a simple source/box/sink topology with 2 sources", t, func() {
		/*
		 *   so1 -*-(left)--\
		 *                   --> b -*--> si
		 *   so2 -*-(right)-/
		 */

		// input data
		tup1 := tuple.Tuple{
			Data: tuple.Map{
				"uid":  tuple.Int(1),
				"key1": tuple.String("tuple1"),
			},
		}
		tup2 := tuple.Tuple{
			Data: tuple.Map{
				"uid":  tuple.Int(2),
				"key2": tuple.String("tuple2"),
			},
		}
		tup3 := tuple.Tuple{
			Data: tuple.Map{
				"uid":  tuple.Int(2),
				"key3": tuple.String("tuple3"),
			},
		}
		tup4 := tuple.Tuple{
			Data: tuple.Map{
				"uid":  tuple.Int(1),
				"key4": tuple.String("tuple4"),
			},
		}
		// expected output when joined on uid (tup1 + tup4, tup2 + tup3)
		j1 := tuple.Tuple{
			Data: tuple.Map{
				"uid":  tuple.Int(1),
				"key1": tuple.String("tuple1"),
				"key4": tuple.String("tuple4"),
			},
		}
		j2 := tuple.Tuple{
			Data: tuple.Map{
				"uid":  tuple.Int(2),
				"key2": tuple.String("tuple2"),
				"key3": tuple.String("tuple3"),
			},
		}

		Convey("When two pairs of tuples are emitted by the sources", func() {
			tb := NewDefaultStaticTopologyBuilder()
			so1 := &TupleEmitterSource{
				Tuples: []*tuple.Tuple{&tup1, &tup2},
			}
			tb.AddSource("source1", so1)
			so2 := &TupleEmitterSource{
				Tuples: []*tuple.Tuple{&tup3, &tup4},
			}
			tb.AddSource("source2", so2)

			b := &SimpleJoinBox{}
			tb.AddBox("box", b).
				NamedInput("source1", "left").
				NamedInput("source2", "right")

			si := &TupleCollectorSink{}
			tb.AddSink("sink", si).Input("box")

			t := tb.Build()

			t.Run(&Context{})
			Convey("Then the sink receives two objects", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 2)

				// the items can arrive in any order, so we can't
				// use a simple ShouldResemble check
				Convey("And one of them is joined on uid 1", func() {
					ok := false
					for _, tup := range si.Tuples {
						if reflect.DeepEqual(tup.Data, j1.Data) {
							ok = true
							break
						}
					}
					So(ok, ShouldBeTrue)
				})
				Convey("And one of them is joined on uid 2", func() {
					ok := false
					for _, tup := range si.Tuples {
						if reflect.DeepEqual(tup.Data, j2.Data) {
							ok = true
							break
						}
					}
					So(ok, ShouldBeTrue)
				})

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})

			Convey("And the buffers in the box are empty", func() {
				So(b.LeftTuples, ShouldBeEmpty)
				So(b.RightTuples, ShouldBeEmpty)
			})
		})

		Convey("When a non-pair of tuples is emitted by the sources", func() {
			tb := NewDefaultStaticTopologyBuilder()
			so1 := &TupleEmitterSource{
				Tuples: []*tuple.Tuple{&tup1},
			}
			tb.AddSource("source1", so1)
			so2 := &TupleEmitterSource{
				Tuples: []*tuple.Tuple{&tup2},
			}
			tb.AddSource("source2", so2)

			b := &SimpleJoinBox{}
			tb.AddBox("box", b).
				NamedInput("source1", "left").
				NamedInput("source2", "right")

			si := &TupleCollectorSink{}
			tb.AddSink("sink", si).Input("box")

			t := tb.Build()

			t.Run(&Context{})
			Convey("Then the sink receives no objects", func() {
				So(si.Tuples, ShouldBeNil)

				Convey("And the tuples should still be in the boxes", func() {
					So(len(b.LeftTuples), ShouldEqual, 1)
					So(len(b.RightTuples), ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given a simple source/box/sink topology with 2 sinks", t, func() {
		/*
		 *                /--> si1
		 *   so -*--> b -*
		 *                \--> si2
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1, &tup2},
		}
		tb.AddSource("source", so)

		b := BoxFunc(forwardBox)
		tb.AddBox("box", b).Input("source")

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("box")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("box")

		t := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(&Context{})
			Convey("Then the sink 1 receives a copy", func() {
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 2)
				// contents are the same
				si := si1
				So(so.Tuples[0].Data, ShouldResemble, si.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si.Tuples[1].BatchID)
				// box has two received sinks, so tuples are copied
				So(so.Tuples[0], ShouldNotPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
			Convey("And the sink 2 receives a copy", func() {
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 2)
				// contents are the same
				si := si2
				So(so.Tuples[0].Data, ShouldResemble, si.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si.Tuples[1].BatchID)
				// pointers point to different objects
				So(so.Tuples[0], ShouldNotPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
			Convey("And the traces of tuples differ", func() {
				So(len(si1.Tuples), ShouldEqual, 2)
				So(len(si2.Tuples), ShouldEqual, 2)
				So(si1.Tuples[0].Trace, ShouldNotResemble, si2.Tuples[0].Trace)
				So(si1.Tuples[1].Trace, ShouldNotResemble, si2.Tuples[1].Trace)
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
			Tuples: []*tuple.Tuple{&tup1, &tup2},
		}
		tb.AddSource("source", so)

		b1 := CollectorBox{InputSchema: map[string]*Schema{"hoge": nil}}
		tb.AddBox("box1", &b1).NamedInput("source", "hoge")
		b2 := CollectorBox{}
		tb.AddBox("box2", &b2).Input("source")

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("box1")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("box2")

		t := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(&Context{})
			Convey("Then the box 1 sees this tuple", func() {
				So(b1.Tuples, ShouldNotBeNil)
				So(len(b1.Tuples), ShouldEqual, 2)
				So(b1.Tuples[0].Data, ShouldResemble, tup1.Data)
				So(b1.Tuples[1].Data, ShouldResemble, tup2.Data)

				Convey("And the InputName is set to \"hoge\"", func() {
					So(b1.Tuples[0].InputName, ShouldEqual, "hoge")
				})
			})
			Convey("And the sink 1 receives a copy", func() {
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 2)
				// contents are the same
				si := si1
				So(so.Tuples[0].Data, ShouldResemble, si.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si.Tuples[1].BatchID)
				// source has two received boxes, so tuples are copied
				So(so.Tuples[0], ShouldNotPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
			Convey("Then the box 2 sees this tuple", func() {
				So(b2.Tuples, ShouldNotBeNil)
				So(len(b2.Tuples), ShouldEqual, 2)
				So(b2.Tuples[0].Data, ShouldResemble, tup1.Data)
				So(b2.Tuples[1].Data, ShouldResemble, tup2.Data)

				Convey("And the InputName is set to \"*\"", func() {
					So(b2.Tuples[0].InputName, ShouldEqual, "*")
				})
			})
			Convey("And the sink 2 receives a copy", func() {
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 2)
				// contents are the same
				si := si2
				So(so.Tuples[0].Data, ShouldResemble, si.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si.Tuples[1].BatchID)
				// pointers point to different objects
				So(so.Tuples[0], ShouldNotPointTo, si.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si.Tuples[1])

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
			Convey("And the traces of tuples differ", func() {
				So(len(si1.Tuples), ShouldEqual, 2)
				So(len(si2.Tuples), ShouldEqual, 2)
				So(si1.Tuples[0].Trace, ShouldNotResemble, si2.Tuples[0].Trace)
				So(si1.Tuples[1].Trace, ShouldNotResemble, si2.Tuples[1].Trace)
			})
		})
	})
}

/**************************************************/

// CollectorBox is a simple forwarder box that also stores a copy
// of all forwarded data for later inspection.
type CollectorBox struct {
	mutex       *sync.Mutex
	Tuples      []*tuple.Tuple
	InputSchema map[string]*Schema
}

func (b *CollectorBox) Init(ctx *Context) error {
	b.mutex = &sync.Mutex{}
	return nil
}
func (b *CollectorBox) Process(t *tuple.Tuple, s Writer) error {
	// with multiple sources, there may be multiple concurrent calls
	// of this method (even in tests) so we need to guard the append
	// with a mutex
	b.mutex.Lock()
	b.Tuples = append(b.Tuples, t.Copy())
	b.mutex.Unlock()
	s.Write(t)
	return nil
}
func (b *CollectorBox) InputConstraints() (*BoxInputConstraints, error) {
	if b.InputSchema != nil {
		ic := &BoxInputConstraints{b.InputSchema}
		return ic, nil
	}
	return nil, nil
}
func (b *CollectorBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

/**************************************************/

// SimpleJoinBox is a box that joins two streams, called "left" and "right"
// on an Int field called "uid". When there is an item in a stream with
// a uid value that has been seen before in the other stream, a joined
// tuple is emitted and both are wiped from internal state.
type SimpleJoinBox struct {
	ctx         *Context
	mutex       *sync.Mutex
	LeftTuples  map[int64]*tuple.Tuple
	RightTuples map[int64]*tuple.Tuple
	inputSchema map[string]*Schema
}

func (b *SimpleJoinBox) Init(ctx *Context) error {
	b.mutex = &sync.Mutex{}
	b.ctx = ctx
	b.LeftTuples = make(map[int64]*tuple.Tuple, 0)
	b.RightTuples = make(map[int64]*tuple.Tuple, 0)
	return nil
}

func (b *SimpleJoinBox) Process(t *tuple.Tuple, s Writer) error {
	// get user id and convert it to int64
	userId, err := t.Data.Get("uid")
	if err != nil {
		b.ctx.Logger.DroppedTuple(t, "no uid field")
		return nil
	}
	userIdInt, err := userId.Int()
	if err != nil {
		b.ctx.Logger.DroppedTuple(t, "uid value was not an integer: %v (%v)",
			userId, err)
		return nil
	}
	uid := int64(userIdInt)
	// prevent concurrent access
	b.mutex.Lock()
	// check if we have a matching item in the other stream
	if t.InputName == "left" {
		match, exists := b.RightTuples[uid]
		if exists {
			// we found a match within the tuples from the right stream
			for key, val := range match.Data {
				t.Data[key] = val
			}
			delete(b.RightTuples, uid)
			b.mutex.Unlock()
			fmt.Printf("emit %v\n", t)
			s.Write(t)
		} else {
			// no match, store this for later
			b.LeftTuples[uid] = t
			b.mutex.Unlock()
		}
	} else if t.InputName == "right" {
		match, exists := b.LeftTuples[uid]
		if exists {
			// we found a match within the tuples from the left stream
			for key, val := range match.Data {
				t.Data[key] = val
			}
			delete(b.LeftTuples, uid)
			b.mutex.Unlock()
			s.Write(t)
		} else {
			// no match, store this for later
			b.RightTuples[uid] = t
			b.mutex.Unlock()
		}
	} else {
		b.ctx.Logger.DroppedTuple(t, "invalid input name: %s", t.InputName)
		return fmt.Errorf("tuple %v had an invalid input name: %s "+
			"(not \"left\" or \"right\")", t, t.InputName)
	}
	return nil
}

// require schemafree input from "left" and "right" named streams
func (b *SimpleJoinBox) InputConstraints() (*BoxInputConstraints, error) {
	if b.inputSchema == nil {
		b.inputSchema = map[string]*Schema{"left": nil, "right": nil}
	}
	return &BoxInputConstraints{b.inputSchema}, nil
}

func (b *SimpleJoinBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

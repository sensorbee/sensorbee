package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
	"testing"
	"time"
)

func TestDefaultTopologyBuilderInterface(t *testing.T) {
	Convey("When creating a default topology builder", t, func() {
		var tb StaticTopologyBuilder = NewDefaultStaticTopologyBuilder()
		So(tb, ShouldNotBeNil)
	})

	Convey("Given a default topology builder", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DefaultSource{}
		b := &DefaultBox{}
		si := &DefaultSink{}
		var err DeclarerError

		Convey("when using a source name twice", func() {
			err = tb.AddSource("mySource", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("mySource", s)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a source name with a box name", func() {
			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a source name with a sink name", func() {
			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name twice", func() {
			err = tb.AddBox("myBox", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("myBox", b)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name with a source name", func() {
			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name with a sink name", func() {
			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name twice", func() {
			err = tb.AddSink("mySink", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("mySink", si)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name with a source name", func() {
			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name with a box name", func() {
			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})
	})

	Convey("Given a default topology builder with a source", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DefaultSource{}
		tb.AddSource("aSource", s)
		b := &DefaultBox{}
		tb.AddBox("aBox", b)
		si := &DefaultSink{}
		var err DeclarerError

		Convey("when a new box references a non-existing item", func() {
			err = tb.AddBox("otherBox", b).
				Input("something", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing source", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing box", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references multiple items", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing source twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource", nil).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing box twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil).
				Input("aBox", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references a non-existing item", func() {
			err = tb.AddSink("aSink", si).
				Input("something")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references an existing source", func() {
			err = tb.AddSink("aSink", si).
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references an existing box", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references multiple items", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references an existing source twice", func() {
			err = tb.AddSink("aSink", si).
				Input("aSource").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references an existing box twice", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox").
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})
	})
}

func TestBasicDefaultTopologyTransport(t *testing.T) {
	Convey("Given basic topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run(&Context{})
			So(si.results[0], ShouldEqual, "VALUE")
		})
	})

	Convey("Given 2 sources topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		s2 := &DummyDefaultSource{"hoge"}
		tb.AddSource("source2", s2)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).
			Input("source1", nil).
			Input("source2", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			start := time.Now()
			t.Run(&Context{})
			So(len(si.results), ShouldEqual, 2)
			So(si.results, ShouldContain, "VALUE")
			So(si.results, ShouldContain, "HOGE")
			So(start, ShouldHappenWithin, 600*time.Millisecond, time.Now())
		})
	})

	Convey("Given 2 tuples in 1source topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s := &DummyDefaultSource2{"value", "hoge"}
		tb.AddSource("source", s)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).
			Input("source", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run(&Context{})
			So(si.results, ShouldResemble, []string{"VALUE", "HOGE"})
		})
	})

	Convey("Given 2 boxes topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1", nil)
		b2 := &DummyAddSuffixBox{}
		tb.AddBox("bBox", b2).Input("source1", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox").Input("bBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run(&Context{})
			So(si.results[0], ShouldEqual, "VALUE")
			So(si.results2[0], ShouldEqual, "value_1")
		})
	})

	Convey("Given 2 sinks topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		si2 := &DummyDefaultSink{}
		tb.AddSink("si2", si2).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run(&Context{})
			So(si.results[0], ShouldEqual, "VALUE")
			So(si2.results[0], ShouldEqual, "VALUE")
		})
	})

	Convey("Given basic topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")

		Convey("When Run topology with Context & ConsoleLogger", func() {
			ctx := &Context{
				&ConsoleLogManager{},
			}
			t := tb.Build()
			t.Run(ctx)
			So(si.results[0], ShouldEqual, "VALUE")
		})
	})

}

type DummyDefaultSource struct{ initial string }

func (s *DummyDefaultSource) GenerateStream(w Writer) error {
	time.Sleep(0.5 * 1e9) // to confirm .Run() goroutine
	t := &tuple.Tuple{}
	t.Data = tuple.Map{
		"source": tuple.String(s.initial),
	}
	w.Write(t)
	return nil
}
func (s *DummyDefaultSource) Schema() *Schema {
	var sc Schema = Schema("test")
	return &sc
}

type DummyDefaultSource2 struct {
	initial  string
	initial2 string
}

func (s *DummyDefaultSource2) GenerateStream(w Writer) error {
	t := &tuple.Tuple{}
	t.Data = tuple.Map{
		"source": tuple.String(s.initial),
	}
	w.Write(t)
	t2 := &tuple.Tuple{}
	t2.Data = tuple.Map{
		"source": tuple.String(s.initial2),
	}
	w.Write(t2)
	return nil
}
func (s *DummyDefaultSource2) Schema() *Schema {
	var sc Schema = Schema("test")
	return &sc
}

type DummyToUpperBox struct {
	ctx *Context
}

func (b *DummyToUpperBox) Init(ctx *Context) error {
	b.ctx = ctx
	return nil
}
func (b *DummyToUpperBox) Process(t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.String()
	t.Data["to-upper"] = tuple.String(strings.ToUpper(string(s)))
	w.Write(t)
	if b.ctx.Logger != nil {
		b.ctx.Logger.Log(DEBUG, "convey test %v", "ToUpperBox Processing")
	}
	return nil
}

func (b *DummyToUpperBox) RequiredInputSchema() ([]*Schema, error) {
	return []*Schema{nil}, nil
}

func (b *DummyToUpperBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

type DummyAddSuffixBox struct{}

func (b *DummyAddSuffixBox) Init(ctx *Context) error {
	return nil
}

func (b *DummyAddSuffixBox) Process(t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.String()
	t.Data["add-suffix"] = tuple.String(s + "_1")
	w.Write(t)
	return nil
}

func (b *DummyAddSuffixBox) RequiredInputSchema() ([]*Schema, error) {
	return []*Schema{nil}, nil
}

func (b *DummyAddSuffixBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

type DummyDefaultSink struct {
	results  []string
	results2 []string
}

func (s *DummyDefaultSink) Write(t *tuple.Tuple) error {
	x, err := t.Data.Get("to-upper")
	if err != nil {
		return nil
	}
	str, _ := x.String()
	s.results = append(s.results, string(str))

	x, err = t.Data.Get("add-suffix")
	if err != nil {
		return nil
	}
	str, _ = x.String()
	s.results2 = append(s.results2, string(str))
	return nil
}

func TestDefaultTopologyTupleCopying(t *testing.T) {
	tup1 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(1),
		},
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	tup2 := tuple.Tuple{
		Data: tuple.Map{
			"int": tuple.Int(2),
		},
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
			Convey("Then the sink 1 receives the same object", func() {
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 2)
				// pointers point to the same objects
				So(so.Tuples[0], ShouldPointTo, si1.Tuples[0])
				So(so.Tuples[1], ShouldPointTo, si1.Tuples[1])
			})
			Convey("And the sink 2 receives a copy", func() {
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 2)
				// contents are the same
				So(so.Tuples[0].Data, ShouldResemble, si2.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si2.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si2.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si2.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si2.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si2.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si2.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si2.Tuples[1].BatchID)
				// tracer is not equal (last input sink is different between sink1 and sink2)
				So(si1.Tuples[0].Tracers, ShouldNotResemble, si2.Tuples[0].Tracers)
				So(si1.Tuples[1].Tracers, ShouldNotResemble, si2.Tuples[1].Tracers)
				// pointers point to different objects
				So(so.Tuples[0], ShouldNotPointTo, si2.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si2.Tuples[1])
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
		tb.AddBox("box", &b).Input("source", nil)

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
		tb.AddBox("box", &b).Input("source", nil)

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("box")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("box")

		t := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(&Context{})
			Convey("Then the sink 1 receives the same object", func() {
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 2)
				// pointers point to the same objects
				So(so.Tuples[0], ShouldPointTo, si1.Tuples[0])
				So(so.Tuples[1], ShouldPointTo, si1.Tuples[1])
			})
			Convey("And the sink 2 receives a copy", func() {
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 2)
				// contents are the same
				So(so.Tuples[0].Data, ShouldResemble, si2.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si2.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si2.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si2.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si2.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si2.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si2.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si2.Tuples[1].BatchID)
				// tracer is not equal (last input sink is different between sink1 and sink2)
				So(si1.Tuples[0].Tracers, ShouldNotResemble, si2.Tuples[0].Tracers)
				So(si1.Tuples[1].Tracers, ShouldNotResemble, si2.Tuples[1].Tracers)
				// pointers point to different objects
				So(so.Tuples[0], ShouldNotPointTo, si2.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si2.Tuples[1])
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

		b1 := BoxFunc(forwardBox)
		tb.AddBox("box1", &b1).Input("source", nil)
		b2 := BoxFunc(forwardBox)
		tb.AddBox("box2", &b2).Input("source", nil)

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("box1")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("box2")

		t := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(&Context{})
			Convey("Then the sink 1 receives the same object", func() {
				So(si1.Tuples, ShouldNotBeNil)
				So(len(si1.Tuples), ShouldEqual, 2)
				// pointers point to the same objects
				So(so.Tuples[0], ShouldPointTo, si1.Tuples[0])
				So(so.Tuples[1], ShouldPointTo, si1.Tuples[1])
			})
			Convey("And the sink 2 receives a copy", func() {
				So(si2.Tuples, ShouldNotBeNil)
				So(len(si2.Tuples), ShouldEqual, 2)
				// contents are the same
				So(so.Tuples[0].Data, ShouldResemble, si2.Tuples[0].Data)
				So(so.Tuples[0].Timestamp, ShouldResemble, si2.Tuples[0].Timestamp)
				So(so.Tuples[0].ProcTimestamp, ShouldResemble, si2.Tuples[0].ProcTimestamp)
				So(so.Tuples[0].BatchID, ShouldEqual, si2.Tuples[0].BatchID)
				So(so.Tuples[1].Data, ShouldResemble, si2.Tuples[1].Data)
				So(so.Tuples[1].Timestamp, ShouldResemble, si2.Tuples[1].Timestamp)
				So(so.Tuples[1].ProcTimestamp, ShouldResemble, si2.Tuples[1].ProcTimestamp)
				So(so.Tuples[1].BatchID, ShouldEqual, si2.Tuples[1].BatchID)
				// tracer is not equal
				So(si1.Tuples[0].Tracers, ShouldNotResemble, si2.Tuples[0].Tracers)
				So(si1.Tuples[1].Tracers, ShouldNotResemble, si2.Tuples[1].Tracers)
				// pointers point to different objects
				So(so.Tuples[0], ShouldNotPointTo, si2.Tuples[0])
				So(so.Tuples[1], ShouldNotPointTo, si2.Tuples[1])
			})
		})
	})
}

type TupleEmitterSource struct {
	Tuples []*tuple.Tuple
}

func (s *TupleEmitterSource) GenerateStream(w Writer) error {
	for _, t := range s.Tuples {
		w.Write(t)
	}
	return nil
}
func (s *TupleEmitterSource) Schema() *Schema {
	return nil
}

type TupleCollectorSink struct {
	Tuples []*tuple.Tuple
}

func (s *TupleCollectorSink) Write(t *tuple.Tuple) error {
	s.Tuples = append(s.Tuples, t)
	return nil
}

func forwardBox(t *tuple.Tuple, w Writer) error {
	w.Write(t)
	return nil
}

func TestDefaultTopologyTupleTracing(t *testing.T) {
	Convey("Given complex topology, has distribution and aggregation", t, func() {

		tup1 := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(1),
			},
			Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
			BatchID:       7,
			Tracers:       make([]tuple.Tracer, 0),
		}
		tup2 := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(2),
			},
			Timestamp:     time.Date(2015, time.April, 10, 10, 23, 1, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 1, 0, time.UTC),
			BatchID:       7,
			Tracers:       make([]tuple.Tracer, 0),
		}
		/*
		 *   so1 \        /--> b2 \        /-*--> si1
		 *        *- b1 -*         *- b4 -*
		 *   so2 /        \--> b3 /        \-*--> si2
		 */
		tb := NewDefaultStaticTopologyBuilder()
		so1 := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1},
		}
		tb.AddSource("so1", so1)
		so2 := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup2},
		}
		tb.AddSource("so2", so2)

		b1 := BoxFunc(forwardBox)
		tb.AddBox("box1", &b1).
			Input("so1", nil).
			Input("so2", nil)
		b2 := BoxFunc(forwardBox)
		tb.AddBox("box2", &b2).Input("box1", nil)
		b3 := BoxFunc(forwardBox)
		tb.AddBox("box3", &b3).Input("box1", nil)
		b4 := BoxFunc(forwardBox)
		tb.AddBox("box4", &b4).
			Input("box2", nil).
			Input("box3", nil)

		si1 := &TupleCollectorSink{}
		tb.AddSink("si1", si1).Input("box4")
		si2 := &TupleCollectorSink{}
		tb.AddSink("si2", si2).Input("box4")

		to := tb.Build()
		Convey("When a tuple is emitted by the source", func() {
			to.Run(&Context{})
			Convey("Then tracer has 2 kind of route from source1", func() {
				// make expected routes
				route1 := []string{
					"OUTPUT so1", "INPUT box1", "OUTPUT box1", "INPUT box2",
					"OUTPUT box2", "INPUT box4", "OUTPUT box4", "INPUT si1",
				}
				route2 := []string{
					"OUTPUT so1", "INPUT box1", "OUTPUT box1", "INPUT box3",
					"OUTPUT box3", "INPUT box4", "OUTPUT box4", "INPUT si1",
				}
				route3 := []string{
					"OUTPUT so2", "INPUT box1", "OUTPUT box1", "INPUT box2",
					"OUTPUT box2", "INPUT box4", "OUTPUT box4", "INPUT si1",
				}
				route4 := []string{
					"OUTPUT so2", "INPUT box1", "OUTPUT box1", "INPUT box3",
					"OUTPUT box3", "INPUT box4", "OUTPUT box4", "INPUT si1",
				}
				eRoutes := []string{
					strings.Join(route1, "->"),
					strings.Join(route2, "->"),
					strings.Join(route3, "->"),
					strings.Join(route4, "->"),
				}
				aRoutes := make([]string, 0)
				for _, tu := range si1.Tuples {
					aRoute := make([]string, 0)
					for _, tr := range tu.Tracers {
						aRoute = append(aRoute, tr.Inout.String()+" "+tr.Msg)
					}
					aRoutes = append(aRoutes, strings.Join(aRoute, "->"))
				}
				So(len(aRoutes), ShouldEqual, 4)
				So(aRoutes, ShouldContain, eRoutes[0])
				So(aRoutes, ShouldContain, eRoutes[1])
				So(aRoutes, ShouldContain, eRoutes[2])
				So(aRoutes, ShouldContain, eRoutes[3])
			})
			Convey("Then tracer has 2 kind of route from source2", func() {
				// make expected routes
				route1 := []string{
					"OUTPUT so1", "INPUT box1", "OUTPUT box1", "INPUT box2",
					"OUTPUT box2", "INPUT box4", "OUTPUT box4", "INPUT si2",
				}
				route2 := []string{
					"OUTPUT so1", "INPUT box1", "OUTPUT box1", "INPUT box3",
					"OUTPUT box3", "INPUT box4", "OUTPUT box4", "INPUT si2",
				}
				route3 := []string{
					"OUTPUT so2", "INPUT box1", "OUTPUT box1", "INPUT box2",
					"OUTPUT box2", "INPUT box4", "OUTPUT box4", "INPUT si2",
				}
				route4 := []string{
					"OUTPUT so2", "INPUT box1", "OUTPUT box1", "INPUT box3",
					"OUTPUT box3", "INPUT box4", "OUTPUT box4", "INPUT si2",
				}
				eRoutes := []string{
					strings.Join(route1, "->"),
					strings.Join(route2, "->"),
					strings.Join(route3, "->"),
					strings.Join(route4, "->"),
				}
				aRoutes := make([]string, 0)
				for _, tu := range si2.Tuples {
					aRoute := make([]string, 0)
					for _, tr := range tu.Tracers {
						aRoute = append(aRoute, tr.Inout.String()+" "+tr.Msg)
					}
					aRoutes = append(aRoutes, strings.Join(aRoute, "->"))
				}
				So(len(si2.Tuples), ShouldEqual, 4)
				//fmt.Println(aRoutes[0])
				So(len(aRoutes), ShouldEqual, 4)
				So(aRoutes, ShouldContain, eRoutes[0])
				So(aRoutes, ShouldContain, eRoutes[1])
				So(aRoutes, ShouldContain, eRoutes[2])
				So(aRoutes, ShouldContain, eRoutes[3])
			})
		})
	})
}

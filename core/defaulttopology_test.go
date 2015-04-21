package core

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
	"strings"
	"sync"
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
		s := &DoesNothingSource{}
		b := &DoesNothingBox{}
		si := &DoesNothingSink{}
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
		s := &DoesNothingSource{}
		tb.AddSource("aSource", s)
		b := &DoesNothingBox{}
		tb.AddBox("aBox", b)
		si := &DoesNothingSink{}
		var err DeclarerError

		Convey("when a new box references a non-existing item", func() {
			err = tb.AddBox("otherBox", b).
				Input("something")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing source", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing box", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references multiple items", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing source twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing box twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox").
				Input("aBox")
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

func TestDefaultTopologyBuilderSchemaChecks(t *testing.T) {

	Convey("Given a default topology builder", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DoesNothingSource{}
		tb.AddSource("source", s)

		Convey("When using a box with nil input constraint", func() {
			// A box with InputConstraint() == nil should allow any and all input
			b := &DoesNothingBoxWithSchema{}

			Convey("Then adding an unnamed input should succeed", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		Convey("When using a box with {'*' => nil} input constraint", func() {
			// A box with '*' => nil should allow any and all input
			b := &DoesNothingBoxWithSchema{
				InputSchema: map[string]*Schema{"*": nil}}

			Convey("Then adding an unnamed input should succeed", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		SkipConvey("When using a box with {'*' => non-nil} input constraint", func() {
			// A box with '*' => aSchema should only allow input
			// if it matches aSchema
		})

		// Note that the following test reveals a questionable behavior: If
		// a Box declares just one input stream, there is apparently no need
		// to tell apart different streams; so what exactly is the value of
		// requiring a certain name?
		Convey("When using a box with {'hoge' => nil} input constraint", func() {
			// A box with 'hoge' => nil should only allow input
			// if it comes from an input stream called 'hoge'
			b := &DoesNothingBoxWithSchema{
				InputSchema: map[string]*Schema{"hoge": nil}}

			Convey("Then adding an unnamed input should fail", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for '*' should fail", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'foo' should fail", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "foo")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		SkipConvey("When using a box with {'hoge' => non-nil} input constraint", func() {
			// A box with 'hoge' => aSchema should only allow input
			// if it comes from an input stream called 'hoge' and
			// matches aSchema
		})

		Convey("When using a box with {'hoge' => nil, '*' => nil} input constraint", func() {
			// A box with 'hoge' => nil, *' => nil should allow any and
			// all input
			b := &DoesNothingBoxWithSchema{
				InputSchema: map[string]*Schema{"hoge": nil, "*": nil}}

			Convey("Then adding an unnamed input should succeed", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'foo' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "foo")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		SkipConvey("When using a box with {'hoge' => non-nil, '*' => non-nil} input constraint", func() {
			// A box with 'hoge' => aSchema, '*' => otherSchema should
			// allow input from arbitrarily named input streams, as long as
			// they match the corresponding schema
		})

		SkipConvey("When using a box with {'hoge' => non-nil, 'foo' => non-nil} input constraint", func() {
			// A box with 'hoge' => aSchema, 'foo' => otherSchema should
			// only allow input from input streams called 'hoge' or 'foo'
			// and matching the corresponding schema
		})
	})

}

func TestBasicDefaultTopologyTransport(t *testing.T) {
	tup1 := tuple.Tuple{
		Data: tuple.Map{
			"source": tuple.String("value"),
		}}
	tup2 := tuple.Tuple{
		Data: tuple.Map{
			"source": tuple.String("hoge"),
		}}

	Convey("Given basic topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &TupleEmitterSource{[]*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1")
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
		s1 := &TupleEmitterSource{[]*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		s2 := &TupleEmitterSource{[]*tuple.Tuple{&tup2}}
		tb.AddSource("source2", s2)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).
			Input("source1").
			Input("source2")
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
		s := &TupleEmitterSource{[]*tuple.Tuple{&tup1, &tup2}}
		tb.AddSource("source", s)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).
			Input("source")
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
		s1 := &TupleEmitterSource{[]*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1")
		b2 := &DummyAddSuffixBox{}
		tb.AddBox("bBox", b2).Input("source1")
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
		s1 := &TupleEmitterSource{[]*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1")
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
		s1 := &TupleEmitterSource{[]*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		b1 := &DummyToUpperBox{}
		tb.AddBox("aBox", b1).Input("source1")
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")

		Convey("When Run topology with Context & ConsoleLogger", func() {
			ctx := &Context{
				NewConsolePrintLogger(),
			}
			t := tb.Build()
			t.Run(ctx)
			So(si.results[0], ShouldEqual, "VALUE")
		})
	})
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

func (b *DummyToUpperBox) InputConstraints() (*BoxInputConstraints, error) {
	return nil, nil
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

func (b *DummyAddSuffixBox) InputConstraints() (*BoxInputConstraints, error) {
	return nil, nil
}

func (b *DummyAddSuffixBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

type DummyDefaultSink struct {
	results  []string
	results2 []string
}

func (s *DummyDefaultSink) Write(t *tuple.Tuple) (err error) {
	x, err := t.Data.Get("to-upper")
	if err == nil {
		str, _ := x.String()
		s.results = append(s.results, string(str))
	}

	x, err = t.Data.Get("add-suffix")
	if err == nil {
		str, _ := x.String()
		s.results2 = append(s.results2, string(str))
	}
	return err
}

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
			Trace:         make([]tuple.TraceEvent, 0),
		}
		tup2 := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(2),
			},
			Timestamp:     time.Date(2015, time.April, 10, 10, 23, 1, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 1, 0, time.UTC),
			BatchID:       7,
			Trace:         make([]tuple.TraceEvent, 0),
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
		tb.AddBox("box1", b1).
			Input("so1").
			Input("so2")
		b2 := BoxFunc(forwardBox)
		tb.AddBox("box2", b2).Input("box1")
		b3 := BoxFunc(forwardBox)
		tb.AddBox("box3", b3).Input("box1")
		b4 := BoxFunc(forwardBox)
		tb.AddBox("box4", b4).
			Input("box2").
			Input("box3")

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
					for _, ev := range tu.Trace {
						aRoute = append(aRoute, ev.Type.String()+" "+ev.Msg)
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
					for _, ev := range tu.Trace {
						aRoute = append(aRoute, ev.Type.String()+" "+ev.Msg)
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

/**************************************************/

type DoesNothingSource struct{}

func (s *DoesNothingSource) GenerateStream(w Writer) error {
	return nil
}
func (s *DoesNothingSource) Schema() *Schema {
	var sc Schema = Schema("test")
	return &sc
}

/**************************************************/

type DoesNothingBox struct {
}

func (b *DoesNothingBox) Init(ctx *Context) error {
	return nil
}
func (b *DoesNothingBox) Process(t *tuple.Tuple, s Writer) error {
	return nil
}
func (b *DoesNothingBox) InputConstraints() (*BoxInputConstraints, error) {
	return nil, nil
}
func (b *DoesNothingBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

/**************************************************/

type DoesNothingSink struct{}

func (s *DoesNothingSink) Write(t *tuple.Tuple) error {
	return nil
}

/**************************************************/

type DoesNothingBoxWithSchema struct {
	InputSchema map[string]*Schema
	DoesNothingBox
}

func (b *DoesNothingBoxWithSchema) InputConstraints() (*BoxInputConstraints, error) {
	if b.InputSchema != nil {
		ic := &BoxInputConstraints{b.InputSchema}
		return ic, nil
	}
	return nil, nil
}

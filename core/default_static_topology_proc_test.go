package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
	"testing"
	"time"
)

// TestDefaultTopologyTupleProcessing tests that tuples are correctly
// processed by boxes in various topologies. This is very similar to
// TestDefaultTopologyTupleTransport but with a focus on box processing
// results.
func TestDefaultTopologyTupleProcessing(t *testing.T) {
	tup1 := tuple.Tuple{
		Data: tuple.Map{
			"source": tuple.String("value"),
		}}
	tup2 := tuple.Tuple{
		Data: tuple.Map{
			"source": tuple.String("hoge"),
		}}

	ToUpperBox := BoxFunc(toUpper)
	AddSuffixBox := BoxFunc(addSuffix)
	ctx := newTestContext(Configuration{})
	Convey("Given a simple source/box/sink topology", t, func() {
		/*
		 *   so -*--> b -*--> si
		 */
		tb := NewDefaultStaticTopologyBuilder()
		s1 := &TupleEmitterSource{Tuples: []*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		b1 := ToUpperBox
		tb.AddBox("aBox", b1).Input("source1")
		si := &TupleContentsCollectorSink{}
		tb.AddSink("si", si).Input("aBox")
		t, _ := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(ctx)
			Convey("Then it is processed by the box", func() {
				So(si.uppercaseResults[0], ShouldEqual, "VALUE")
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
		s1 := &TupleEmitterSource{Tuples: []*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		s2 := &TupleEmitterSource{Tuples: []*tuple.Tuple{&tup2}}
		tb.AddSource("source2", s2)
		b1 := ToUpperBox
		tb.AddBox("aBox", b1).
			Input("source1").
			Input("source2")
		si := &TupleContentsCollectorSink{}
		tb.AddSink("si", si).Input("aBox")
		t, _ := tb.Build()

		Convey("When a tuple is emitted by each source", func() {
			start := time.Now()
			t.Run(ctx)
			Convey("Then they should both be processed by the box in a reasonable time", func() {
				So(len(si.uppercaseResults), ShouldEqual, 2)
				So(si.uppercaseResults, ShouldContain, "VALUE")
				So(si.uppercaseResults, ShouldContain, "HOGE")
				So(start, ShouldHappenWithin, 600*time.Millisecond, time.Now())
			})
		})
	})

	Convey("Given a simple source/box/sink topology", t, func() {
		/*
		 *   so -*--> b -*--> si
		 */
		tb := NewDefaultStaticTopologyBuilder()
		s := &TupleEmitterSource{Tuples: []*tuple.Tuple{&tup1, &tup2}}
		tb.AddSource("source", s)
		b1 := ToUpperBox
		tb.AddBox("aBox", b1).
			Input("source")
		si := &TupleContentsCollectorSink{}
		tb.AddSink("si", si).Input("aBox")
		t, _ := tb.Build()

		Convey("When two tuples are emitted by the source", func() {
			t.Run(ctx)
			Convey("Then they are processed both and in order", func() {
				So(si.uppercaseResults, ShouldResemble, []string{"VALUE", "HOGE"})
			})
		})
	})

	Convey("Given a simple source/box/sink topology with 2 boxes", t, func() {
		/*
		 *        /--> b1 -*-\
		 *   so -*            --> si
		 *        \--> b2 -*-/
		 */
		tb := NewDefaultStaticTopologyBuilder()
		s1 := &TupleEmitterSource{Tuples: []*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		b1 := ToUpperBox
		tb.AddBox("aBox", b1).Input("source1")
		b2 := AddSuffixBox
		tb.AddBox("bBox", b2).Input("source1")
		si := &TupleContentsCollectorSink{}
		tb.AddSink("si", si).Input("aBox").Input("bBox")
		t, _ := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(ctx)
			Convey("Then it is processed by both boxes", func() {
				So(si.uppercaseResults[0], ShouldEqual, "VALUE")
				So(si.suffixResults[0], ShouldEqual, "value_1")
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
		s1 := &TupleEmitterSource{Tuples: []*tuple.Tuple{&tup1}}
		tb.AddSource("source1", s1)
		b1 := ToUpperBox
		tb.AddBox("aBox", b1).Input("source1")
		si := &TupleContentsCollectorSink{}
		tb.AddSink("si", si).Input("aBox")
		si2 := &TupleContentsCollectorSink{}
		tb.AddSink("si2", si2).Input("aBox")
		t, _ := tb.Build()

		Convey("When a tuple is emitted by the source", func() {
			t.Run(ctx)
			Convey("Then the processed value arrives in both sinks", func() {
				So(si.uppercaseResults[0], ShouldEqual, "VALUE")
				So(si2.uppercaseResults[0], ShouldEqual, "VALUE")
			})
		})
	})
}

/**************************************************/

func toUpper(ctx *Context, t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.AsString()
	t.Data["to-upper"] = tuple.String(strings.ToUpper(s))
	w.Write(ctx, t)
	return nil
}

func addSuffix(ctx *Context, t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.AsString()
	t.Data["add-suffix"] = tuple.String(s + "_1")
	w.Write(ctx, t)
	return nil
}

/**************************************************/

// TupleContentsCollectorSink is a sink that will add all strings found
// in the "to-upper" field to the uppercaseResults slice, all strings
// in the "add-suffix" field to the suffixResults slice.
type TupleContentsCollectorSink struct {
	uppercaseResults []string
	suffixResults    []string
}

func (s *TupleContentsCollectorSink) Write(ctx *Context, t *tuple.Tuple) (err error) {
	x, err := t.Data.Get("to-upper")
	if err == nil {
		str, _ := x.AsString()
		s.uppercaseResults = append(s.uppercaseResults, str)
	}

	x, err = t.Data.Get("add-suffix")
	if err == nil {
		str, _ := x.AsString()
		s.suffixResults = append(s.suffixResults, str)
	}
	return err
}

func (s *TupleContentsCollectorSink) Close(ctx *Context) error {
	return nil
}

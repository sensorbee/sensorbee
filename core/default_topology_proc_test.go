package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"strings"
	"testing"
	"time"
)

// TestDefaultTopologyTupleProcessing tests that tuples are correctly
// processed by boxes in various topologies.
func TestDefaultTopologyTupleProcessing(t *testing.T) {
	tup1 := &Tuple{
		Data: data.Map{
			"source": data.String("value"),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	tup2 := &Tuple{
		Data: data.Map{
			"source": data.String("hoge"),
		},
		InputName:     "input",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 1, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 1, 0, time.UTC),
		BatchID:       7,
	}

	ToUpperBox := BoxFunc(toUpper)
	AddSuffixBox := BoxFunc(addSuffix)
	Convey("Given a simple source/box/sink topology", t, func() {
		/*
		 *   so -*--> b -*--> si
		 */
		ctx := NewContext(nil)
		t, err := NewDefaultTopology(ctx, "test")
		So(err, ShouldBeNil)
		Reset(func() {
			t.Stop()
		})

		s1 := newCustomEmitterSource()
		_, err = t.AddSource("source1", s1, nil)
		So(err, ShouldBeNil)

		bn, err := t.AddBox("aBox", ToUpperBox, nil)
		So(err, ShouldBeNil)
		So(bn.Input("source1", nil), ShouldBeNil)
		bn.StopOnDisconnect(Inbound)

		si := &TupleContentsCollectorSink{}
		sin, err := t.AddSink("si", si, nil)
		So(err, ShouldBeNil)
		So(sin.Input("aBox", nil), ShouldBeNil)
		sin.StopOnDisconnect()

		Convey("When a tuple is emitted by the source", func() {
			s1.emit(tup1)
			s1.Stop(ctx)
			sin.State().Wait(TSStopped)

			Convey("Then it is processed by the box", func() {
				So(si.uppercaseResults[0], ShouldEqual, "VALUE")
			})
		})

		Convey("When two tuples are emitted by the source", func() {
			s1.emit(tup1)
			s1.emit(tup2)
			s1.Stop(ctx)
			sin.State().Wait(TSStopped)

			Convey("Then they are processed both and in order", func() {
				So(si.uppercaseResults, ShouldResemble, []string{"VALUE", "HOGE"})
			})

			Convey("Then the sink receives the same object", func() {
				So(si.len(), ShouldEqual, 2)
				// pointers point to the same objects
				So(si.get(0), ShouldPointTo, tup1)
				So(si.get(1), ShouldPointTo, tup2)

				Convey("And the InputName is set to \"output\"", func() {
					So(si.get(0).InputName, ShouldEqual, "output")
				})
			})
		})
	})

	Convey("Given a simple source/box/sink topology with 2 sources", t, func() {
		/*
		 *   so1 -*-\
		 *           --> b -*--> si
		 *   so2 -*-/
		 */
		ctx := NewContext(nil)
		t, err := NewDefaultTopology(ctx, "test")
		So(err, ShouldBeNil)
		Reset(func() {
			t.Stop()
		})

		s1 := &TupleEmitterSource{Tuples: []*Tuple{tup1}}
		son1, err := t.AddSource("source1", s1, &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		s2 := &TupleEmitterSource{Tuples: []*Tuple{tup2}}
		son2, err := t.AddSource("source2", s2, &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		bn, err := t.AddBox("aBox", ToUpperBox, nil)
		So(err, ShouldBeNil)
		So(bn.Input("source1", nil), ShouldBeNil)
		So(bn.Input("source2", nil), ShouldBeNil)
		bn.StopOnDisconnect(Inbound)

		si := &TupleContentsCollectorSink{}
		sin, err := t.AddSink("si", si, nil)
		So(err, ShouldBeNil)
		So(sin.Input("aBox", nil), ShouldBeNil)
		sin.StopOnDisconnect()

		So(son1.Resume(), ShouldBeNil)
		So(son2.Resume(), ShouldBeNil)
		sin.State().Wait(TSStopped)

		Convey("When a tuple is emitted by each source", func() {
			start := time.Now()
			Convey("Then they should both be processed by the box in a reasonable time", func() {
				So(len(si.uppercaseResults), ShouldEqual, 2)
				So(si.uppercaseResults, ShouldContain, "VALUE")
				So(si.uppercaseResults, ShouldContain, "HOGE")
				So(start, ShouldHappenWithin, 600*time.Millisecond, time.Now())
			})
		})
	})

	Convey("Given a simple source/box/sink topology with 2 boxes", t, func() {
		/*
		 *        /--> b1 -*-\
		 *   so -*            --> si
		 *        \--> b2 -*-/
		 */

		ctx := NewContext(nil)
		t, err := NewDefaultTopology(ctx, "test")
		So(err, ShouldBeNil)
		Reset(func() {
			t.Stop()
		})

		s1 := &TupleEmitterSource{Tuples: []*Tuple{tup1}}
		son, err := t.AddSource("source1", s1, &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		bn1, err := t.AddBox("aBox", ToUpperBox, nil)
		So(err, ShouldBeNil)
		So(bn1.Input("source1", nil), ShouldBeNil)
		bn1.StopOnDisconnect(Inbound)

		bn2, err := t.AddBox("bBox", AddSuffixBox, nil)
		So(err, ShouldBeNil)
		So(bn2.Input("source1", nil), ShouldBeNil)
		bn2.StopOnDisconnect(Inbound)

		si := &TupleContentsCollectorSink{}
		sin, err := t.AddSink("si", si, nil)
		So(err, ShouldBeNil)
		So(sin.Input("aBox", nil), ShouldBeNil)
		So(sin.Input("bBox", nil), ShouldBeNil)
		sin.StopOnDisconnect()
		So(son.Resume(), ShouldBeNil)
		sin.State().Wait(TSStopped)

		Convey("When a tuple is emitted by the source", func() {
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

		ctx := NewContext(&ContextConfig{
			Flags: ContextFlags{
				TupleTrace: 1,
			},
		})
		t, err := NewDefaultTopology(ctx, "test")
		So(err, ShouldBeNil)
		Reset(func() {
			t.Stop()
		})

		s1 := &TupleEmitterSource{Tuples: []*Tuple{tup1}}
		son, err := t.AddSource("source1", s1, &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		bn, err := t.AddBox("aBox", ToUpperBox, nil)
		So(err, ShouldBeNil)
		So(bn.Input("source1", nil), ShouldBeNil)
		bn.StopOnDisconnect(Inbound)

		si1 := &TupleContentsCollectorSink{}
		sin1, err := t.AddSink("si1", si1, nil)
		So(err, ShouldBeNil)
		So(sin1.Input("aBox", nil), ShouldBeNil)
		sin1.StopOnDisconnect()

		si2 := &TupleContentsCollectorSink{}
		sin2, err := t.AddSink("si2", si2, nil)
		So(err, ShouldBeNil)
		So(sin2.Input("aBox", nil), ShouldBeNil)
		sin2.StopOnDisconnect()

		So(son.Resume(), ShouldBeNil)
		sin1.State().Wait(TSStopped)
		sin2.State().Wait(TSStopped)

		Convey("When a tuple is emitted by the source", func() {
			Convey("Then the processed value arrives in both sinks", func() {
				So(si1.uppercaseResults[0], ShouldEqual, "VALUE")
				So(si2.uppercaseResults[0], ShouldEqual, "VALUE")
			})

			Convey("Then the sink 1 receives a copy", func() {
				So(si1.len(), ShouldEqual, 1)

				// contents are the same
				t := si1.get(0)
				So(tup1.Data, ShouldResemble, t.Data)
				So(tup1.Timestamp, ShouldResemble, t.Timestamp)
				So(tup1.ProcTimestamp, ShouldResemble, t.ProcTimestamp)
				So(tup1.BatchID, ShouldEqual, t.BatchID)

				// source has two received boxes, so tuples are copied
				// TODO: This is a little to specific to the current implementation.
				So(tup1, ShouldNotPointTo, t)

				Convey("And the InputName is set to \"output\"", func() {
					So(t.InputName, ShouldEqual, "output")
				})
			})

			Convey("And the sink 2 receives a copy", func() {
				So(si2.len(), ShouldEqual, 1)

				// contents are the same
				t := si2.get(0)
				So(tup1.Data, ShouldResemble, t.Data)
				So(tup1.Timestamp, ShouldResemble, t.Timestamp)
				So(tup1.ProcTimestamp, ShouldResemble, t.ProcTimestamp)
				So(tup1.BatchID, ShouldEqual, t.BatchID)

				// pointers point to different objects
				So(tup1, ShouldNotPointTo, t)

				Convey("And the InputName is set to \"output\"", func() {
					So(t.InputName, ShouldEqual, "output")
				})
			})

			Convey("And the traces of the tuple differ", func() {
				So(si1.len(), ShouldEqual, 1)
				So(si2.len(), ShouldEqual, 1)
				So(si1.get(0).Trace, ShouldNotResemble, si2.get(0).Trace)
			})
		})
	})
}

func toUpper(ctx *Context, t *Tuple, w Writer) error {
	x, _ := t.Data.Get(data.MustCompilePath("source"))
	s, _ := data.AsString(x)
	t.Data["to_upper"] = data.String(strings.ToUpper(s))
	w.Write(ctx, t)
	return nil
}

func addSuffix(ctx *Context, t *Tuple, w Writer) error {
	x, _ := t.Data.Get(data.MustCompilePath("source"))
	s, _ := data.AsString(x)
	t.Data["add_suffix"] = data.String(s + "_1")
	w.Write(ctx, t)
	return nil
}

type customEmitterSource struct {
	ch chan *Tuple
}

func newCustomEmitterSource() *customEmitterSource {
	return &customEmitterSource{
		ch: make(chan *Tuple),
	}
}

func (s *customEmitterSource) emit(t *Tuple) {
	s.ch <- t
}

func (s *customEmitterSource) GenerateStream(ctx *Context, w Writer) error {
	for t := range s.ch {
		w.Write(ctx, t)
	}
	return nil
}

func (s *customEmitterSource) Stop(ctx *Context) error {
	close(s.ch)
	return nil
}

// TupleContentsCollectorSink is a sink that will add all strings found
// in the "to_upper" field to the uppercaseResults slice, all strings
// in the "add_suffix" field to the suffixResults slice.
type TupleContentsCollectorSink struct {
	TupleCollectorSink
	uppercaseResults []string
	suffixResults    []string
}

func (s *TupleContentsCollectorSink) Write(ctx *Context, t *Tuple) (err error) {
	s.TupleCollectorSink.Write(ctx, t)

	x, err := t.Data.Get(data.MustCompilePath("to_upper"))
	if err == nil {
		str, _ := data.AsString(x)
		s.uppercaseResults = append(s.uppercaseResults, str)
	}

	x, err = t.Data.Get(data.MustCompilePath("add_suffix"))
	if err == nil {
		str, _ := data.AsString(x)
		s.suffixResults = append(s.suffixResults, str)
	}
	return err
}

func (s *TupleContentsCollectorSink) Close(ctx *Context) error {
	return nil
}

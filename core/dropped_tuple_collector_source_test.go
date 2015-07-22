package core

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type writeFailSink struct {
}

func (f *writeFailSink) Write(ctx *Context, t *Tuple) error {
	return errors.New("sink write error")
}

func (f *writeFailSink) Close(ctx *Context) error {
	return nil
}

func TestDroppedTupleCollectorSource(t *testing.T) {
	Convey("Given a topology and a dropped tuple collector source", t, func() {
		t := NewDefaultTopology(NewContext(nil), "dt1")
		Reset(func() {
			t.Stop()
		})

		son, err := t.AddSource("source", NewTupleEmitterSource(freshTuples()), &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)
		_, err = t.AddSource("dropped_tuples", NewDroppedTupleCollectorSource(), nil)
		So(err, ShouldBeNil)
		_, err = t.AddSource("dropped_tuples2", NewDroppedTupleCollectorSource(), nil)
		So(err, ShouldBeNil)
		si := NewTupleCollectorSink()
		sin, err := t.AddSink("sink", si, nil)
		So(err, ShouldBeNil)
		So(sin.Input("dropped_tuples2", nil), ShouldBeNil)

		Convey("When a sink is connected to dropped tuple collector", func() {
			So(son.Resume(), ShouldBeNil)

			Convey("Then it should receive all dropped tuples", func() {
				si.Wait(8)
				So(len(si.Tuples), ShouldEqual, 8)

				Convey("And tuples should have the documented format", func() {
					So(si.Tuples[0].Data["node_type"], ShouldEqual, NTSource.String())
					So(si.Tuples[0].Data["node_name"], ShouldEqual, "source")
					So(si.Tuples[0].Data["event_type"], ShouldEqual, ETOutput.String())
					So(si.Tuples[0].Data["error"], ShouldNotBeNil)
					So(si.Tuples[0].Data["data"], ShouldResemble, freshTuples()[0].Data)
				})
			})
		})

		locationChecker := func(ts []*Tuple, n Node) {
			for _, t := range ts {
				So(t.Data["node_type"], ShouldEqual, n.Type().String())
				So(t.Data["node_name"], ShouldEqual, n.Name())
			}
		}

		Convey("When tuples are dropped from a Box", func() {
			bn, err := t.AddBox("box", BoxFunc(func(ctx *Context, t *Tuple, w Writer) error {
				return errors.New("box write error")
			}), nil)
			So(err, ShouldBeNil)
			So(bn.Input("source", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then they should be reported", func() {
				si.Wait(8)
				So(len(si.Tuples), ShouldEqual, 8)
				locationChecker(si.Tuples, bn)
			})
		})

		Convey("When tuples are dropped from a Sink", func() {
			sin2, err := t.AddSink("fail_sink", &writeFailSink{}, nil)
			So(err, ShouldBeNil)
			So(sin2.Input("source", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then they should be reported", func() {
				si.Wait(8)
				So(len(si.Tuples), ShouldEqual, 8)
				locationChecker(si.Tuples, sin2)
			})
		})

		Convey("When tuples are dropped from a Sink indirectly connected to the original source", func() {
			bn, err := t.AddBox("box", BoxFunc(forwardBox), nil)
			So(err, ShouldBeNil)
			So(bn.Input("source", nil), ShouldBeNil)
			sin2, err := t.AddSink("fail_sink", &writeFailSink{}, nil)
			So(err, ShouldBeNil)
			So(sin2.Input("box", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then they should be reported", func() {
				si.Wait(8)
				So(len(si.Tuples), ShouldEqual, 8)
				locationChecker(si.Tuples, sin2)
			})
		})

		Convey("When a Box connected to the collector drops all tuples", func() {
			bn, err := t.AddBox("box", BoxFunc(func(ctx *Context, t *Tuple, w Writer) error {
				return errors.New("box write error")
			}), nil)
			So(err, ShouldBeNil)
			So(bn.Input("dropped_tuples", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then tuples dropped again shouldn't be reported twice", func() {
				si.Wait(8)
				So(len(si.Tuples), ShouldEqual, 8)
				locationChecker(si.Tuples, son)
			})
		})

		Convey("When a Sink connected to the collector drops all tuples", func() {
			sin2, err := t.AddSink("fail_sink", &writeFailSink{}, nil)
			So(err, ShouldBeNil)
			So(sin2.Input("dropped_tuples", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then tuples dropped again shouldn't be reported twice", func() {
				si.Wait(8)
				So(len(si.Tuples), ShouldEqual, 8)
				locationChecker(si.Tuples, son)
			})
		})

		Convey("When a Box connected to both regular source and the collector drops all tuples", func() {
			bn, err := t.AddBox("box", BoxFunc(func(ctx *Context, t *Tuple, w Writer) error {
				return errors.New("box write error")
			}), nil)
			So(err, ShouldBeNil)
			So(bn.Input("source", nil), ShouldBeNil)
			So(bn.Input("dropped_tuples", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then no tuple should be reported", func() {
				si.Wait(0)
				So(len(si.Tuples), ShouldEqual, 0)
			})
		})

		Convey("When a Sink indirectly connected to both regular source and the collector drops all tuples", func() {
			bn, err := t.AddBox("box", BoxFunc(forwardBox), nil)
			So(err, ShouldBeNil)
			So(bn.Input("source", nil), ShouldBeNil)
			So(bn.Input("dropped_tuples", nil), ShouldBeNil)
			sin2, err := t.AddSink("fail_sink", &writeFailSink{}, nil)
			So(err, ShouldBeNil)
			So(sin2.Input("box", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then no tuple should be reported", func() {
				si.Wait(0)
				So(len(si.Tuples), ShouldEqual, 0)
			})
		})
	})
}

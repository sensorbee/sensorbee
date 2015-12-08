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
				So(si.len(), ShouldEqual, 8)

				Convey("And tuples should have the documented format", func() {
					data := si.get(0).Data
					So(data["node_type"], ShouldEqual, NTSource.String())
					So(data["node_name"], ShouldEqual, "source")
					So(data["event_type"], ShouldEqual, ETOutput.String())
					So(data["error"], ShouldNotBeNil)
					So(data["data"], ShouldResemble, freshTuples()[0].Data)
				})
			})
		})

		locationChecker := func(t *Tuple, n Node) {
			So(t.Data["node_type"], ShouldEqual, n.Type().String())
			So(t.Data["node_name"], ShouldEqual, n.Name())
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
				So(si.len(), ShouldEqual, 8)
				si.forEachTuple(func(t *Tuple) {
					locationChecker(t, bn)
				})
			})
		})

		Convey("When tuples are dropped from a Sink", func() {
			sin2, err := t.AddSink("fail_sink", &writeFailSink{}, nil)
			So(err, ShouldBeNil)
			So(sin2.Input("source", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then they should be reported", func() {
				si.Wait(8)
				So(si.len(), ShouldEqual, 8)
				si.forEachTuple(func(t *Tuple) {
					locationChecker(t, sin2)
				})
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
				So(si.len(), ShouldEqual, 8)
				si.forEachTuple(func(t *Tuple) {
					locationChecker(t, sin2)
				})
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
				So(si.len(), ShouldEqual, 8)
				si.forEachTuple(func(t *Tuple) {
					locationChecker(t, son)
				})
			})
		})

		Convey("When a Sink connected to the collector drops all tuples", func() {
			sin2, err := t.AddSink("fail_sink", &writeFailSink{}, nil)
			So(err, ShouldBeNil)
			So(sin2.Input("dropped_tuples", nil), ShouldBeNil)
			So(son.Resume(), ShouldBeNil)

			Convey("Then tuples dropped again shouldn't be reported twice", func() {
				si.Wait(8)
				So(si.len(), ShouldEqual, 8)
				si.forEachTuple(func(t *Tuple) {
					locationChecker(t, son)
				})
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

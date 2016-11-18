package core

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
		ctx := NewContext(&ContextConfig{
			Flags: ContextFlags{
				DestinationlessTupleLog: 1,
			},
		})

		t, err := NewDefaultTopology(ctx, "dt1")
		So(err, ShouldBeNil)
		Reset(func() {
			t.Stop()
		})

		son, err := t.AddSource("source", NewTupleEmitterSource(freshTuples()), &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)
		dtso := NewDroppedTupleCollectorSource().(*droppedTupleCollectorSource)
		_, err = t.AddSource("dropped_tuples", dtso, nil)
		So(err, ShouldBeNil)
		dtso.state.Wait(TSRunning)
		dtso2 := NewDroppedTupleCollectorSource().(*droppedTupleCollectorSource)
		_, err = t.AddSource("dropped_tuples2", dtso2, nil)
		So(err, ShouldBeNil)
		dtso2.state.Wait(TSRunning)
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

		Convey("When destinationless tuple logging is disabled", func() {
			ctx.Flags.DestinationlessTupleLog.Set(false)
			So(son.Resume(), ShouldBeNil)

			Convey("Then it should not receive any dropped tuple", func() {
				si.Wait(0)
				So(si.len(), ShouldEqual, 0)
			})
		})

		locationChecker := func(t *Tuple, n Node) {
			So(t.Data["node_type"], ShouldEqual, n.Type().String())
			So(t.Data["node_name"], ShouldEqual, n.Name())
		}

		Convey("When tuples are dropped from a Box", func() {
			// Because tuples in this case are dropped by an error,
			// DestinationlessTupleLog shouldn't change the behavior.
			ctx.Flags.DestinationlessTupleLog.Set(false)

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
	})
}

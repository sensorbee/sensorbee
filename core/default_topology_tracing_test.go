package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/tuple"
	"strings"
	"testing"
	"time"
)

// TestDefaultTopologyTupleTracingConfiguration tests that tracing
// information is added according to Configuration.
func TestDefaultTopologyTupleTracingConfiguration(t *testing.T) {
	Convey("Given a simple topology with tracing disabled", t, func() {
		config := Configuration{TupleTraceEnabled: 0}
		ctx := newTestContext(config)
		tup := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(1),
			},
			Timestamp:     time.Date(2015, time.May, 1, 11, 18, 0, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.May, 1, 11, 18, 0, 0, time.UTC),
			BatchID:       7,
			Trace:         []tuple.TraceEvent{},
		}

		t := NewDefaultTopology(ctx, "test")
		Reset(func() {
			t.Stop()
		})
		so1 := NewTupleIncrementalEmitterSource([]*tuple.Tuple{tup.Copy(), tup.Copy(), tup.Copy()})
		_, err := t.AddSource("so1", so1, nil)
		So(err, ShouldBeNil)

		b := BoxFunc(forwardBox)
		bn, err := t.AddBox("box", b, nil)
		So(err, ShouldBeNil)
		So(bn.Input("so1", nil), ShouldBeNil)

		si := NewTupleCollectorSink()
		sin, err := t.AddSink("si", si, nil)
		So(err, ShouldBeNil)
		So(sin.Input("box", nil), ShouldBeNil)

		Convey("When switch tracing configuration in running topology", func() {
			so1.EmitTuples(1)
			si.Wait(1)
			ctx.SetTupleTraceEnabled(true)
			so1.EmitTuples(1)
			si.Wait(2)
			ctx.SetTupleTraceEnabled(false)
			so1.EmitTuples(1)
			si.Wait(3)
			Convey("Then trace should be according to configuration", func() {
				So(len(si.Tuples), ShouldEqual, 3)
				So(len(si.Tuples[0].Trace), ShouldEqual, 0)
				So(len(si.Tuples[1].Trace), ShouldEqual, 4)
				So(len(si.Tuples[2].Trace), ShouldEqual, 0)
			})
		})
	})
}

// TestDefaultTopologyTupleTracing tests that tracing information is
// correctly added to tuples in a complex topology.
func TestDefaultTopologyTupleTracing(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)
	Convey("Given a complex topology with distribution and aggregation", t, func() {

		tup1 := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(1),
			},
			Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
			BatchID:       7,
			Trace:         []tuple.TraceEvent{},
		}
		tup2 := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(2),
			},
			Timestamp:     time.Date(2015, time.April, 10, 10, 23, 1, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 1, 0, time.UTC),
			BatchID:       7,
			Trace:         []tuple.TraceEvent{},
		}
		/*
		 *   so1 \        /--> b2 \        /-*--> si1
		 *        *- b1 -*         *- b4 -*
		 *   so2 /        \--> b3 /        \-*--> si2
		 */
		t := NewDefaultTopology(ctx, "test")
		Reset(func() {
			t.Stop()
		})
		so1 := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1},
		}
		son1, err := t.AddSource("so1", so1, &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		so2 := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup2},
		}
		son2, err := t.AddSource("so2", so2, &SourceConfig{
			PausedOnStartup: true,
		})
		So(err, ShouldBeNil)

		b1 := BoxFunc(forwardBox)
		bn1, err := t.AddBox("box1", b1, nil)
		So(err, ShouldBeNil)
		So(bn1.Input("so1", nil), ShouldBeNil)
		So(bn1.Input("so2", nil), ShouldBeNil)

		b2 := BoxFunc(forwardBox)
		bn2, err := t.AddBox("box2", b2, nil)
		So(err, ShouldBeNil)
		So(bn2.Input("box1", nil), ShouldBeNil)

		b3 := BoxFunc(forwardBox)
		bn3, err := t.AddBox("box3", b3, nil)
		So(err, ShouldBeNil)
		So(bn3.Input("box1", nil), ShouldBeNil)

		b4 := BoxFunc(forwardBox)
		bn4, err := t.AddBox("box4", b4, nil)
		So(err, ShouldBeNil)
		So(bn4.Input("box2", nil), ShouldBeNil)
		So(bn4.Input("box3", nil), ShouldBeNil)

		si1 := &TupleCollectorSink{}
		sin1, err := t.AddSink("si1", si1, nil)
		So(err, ShouldBeNil)
		So(sin1.Input("box4", nil), ShouldBeNil)

		si2 := &TupleCollectorSink{}
		sin2, err := t.AddSink("si2", si2, nil)
		So(err, ShouldBeNil)
		So(sin2.Input("box4", nil), ShouldBeNil)

		bn1.StopOnDisconnect()
		bn2.StopOnDisconnect()
		bn3.StopOnDisconnect()
		bn4.StopOnDisconnect()
		sin1.StopOnDisconnect()
		sin2.StopOnDisconnect()
		So(son1.Resume(), ShouldBeNil)
		So(son2.Resume(), ShouldBeNil)
		sin1.State().Wait(TSStopped)
		sin2.State().Wait(TSStopped)

		Convey("When a tuple is emitted by the source", func() {
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
				var aRoutes []string
				for _, tu := range si1.Tuples {
					var aRoute []string
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
				var aRoutes []string
				for _, tu := range si2.Tuples {
					var aRoute []string
					for _, ev := range tu.Trace {
						aRoute = append(aRoute, ev.Type.String()+" "+ev.Msg)
					}
					aRoutes = append(aRoutes, strings.Join(aRoute, "->"))
				}
				So(len(si2.Tuples), ShouldEqual, 4)
				So(len(aRoutes), ShouldEqual, 4)
				So(aRoutes, ShouldContain, eRoutes[0])
				So(aRoutes, ShouldContain, eRoutes[1])
				So(aRoutes, ShouldContain, eRoutes[2])
				So(aRoutes, ShouldContain, eRoutes[3])
			})
		})
	})
}

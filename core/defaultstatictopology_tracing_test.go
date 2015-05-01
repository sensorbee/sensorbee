package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
	"testing"
	"time"
)

// TestDefaultTopologyTupleTracingConfiguration test that tracing
// information is added according to Configuration.
func TestDefaultTopologyTupleTracingConfiguration(t *testing.T) {
	Convey("Given a simple topology with tracing disabled", t, func() {
		config := Configuration{TupleTraceEnabled: 0}
		ctx := Context{Config: config}
		tup := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(1),
			},
			Timestamp:     time.Date(2015, time.May, 1, 11, 18, 0, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.May, 1, 11, 18, 0, 0, time.UTC),
			BatchID:       7,
			Trace:         make([]tuple.TraceEvent, 0),
		}
		tb := NewDefaultStaticTopologyBuilder()
		so1 := TupleWaitEmitterSource{
			tuple:    tup.Copy(),
			waitTime: 0,
		}
		tb.AddSource("so1", &so1)
		so2 := TupleWaitEmitterSource{
			tuple:    tup.Copy(),
			waitTime: 50,
		}
		tb.AddSource("so2", &so2)
		so3 := TupleWaitEmitterSource{
			tuple:    tup.Copy(),
			waitTime: 100,
		}
		tb.AddSource("so3", &so3)
		b := BoxFunc(forwardBox)
		tb.AddBox("box", b).Input("so1").Input("so2").Input("so3")
		si := &TupleCollectorSink{}
		tb.AddSink("si", si).Input("box")

		tp, _ := tb.Build()
		Convey("When switch tracing configuration in running topology", func() {
			go tp.Run(&ctx)
			time.Sleep(25 * time.Millisecond)
			ctx.SetTupleTraceEnabled(true)
			time.Sleep(50 * time.Millisecond)
			ctx.SetTupleTraceEnabled(false)
			time.Sleep(50 * time.Millisecond)
			Convey("Then trace should be according to configuration", func() {
				So(len(si.Tuples), ShouldEqual, 3)
				So(len(si.Tuples[0].Trace), ShouldEqual, 0)
				So(len(si.Tuples[1].Trace), ShouldEqual, 4)
				So(len(si.Tuples[2].Trace), ShouldEqual, 0)
			})
		})
	})
}

type TupleWaitEmitterSource struct {
	tuple    *tuple.Tuple
	waitTime uint
}

func (s *TupleWaitEmitterSource) GenerateStream(ctx *Context, w Writer) error {
	sleepTime := time.Duration(s.waitTime) * time.Millisecond
	time.Sleep(sleepTime)
	w.Write(ctx, s.tuple)
	return nil
}

func (s *TupleWaitEmitterSource) Schema() *Schema {
	return nil
}

// TestDefaultTopologyTupleTracing tests that tracing information is
// correctly added to tuples in a complex topology.
func TestDefaultTopologyTupleTracing(t *testing.T) {
	config := Configuration{TupleTraceEnabled: 1}
	ctx := Context{Config: config}
	Convey("Given a complex topology with distribution and aggregation", t, func() {

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

		to, _ := tb.Build()
		Convey("When a tuple is emitted by the source", func() {
			to.Run(&ctx)
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
				So(len(aRoutes), ShouldEqual, 4)
				So(aRoutes, ShouldContain, eRoutes[0])
				So(aRoutes, ShouldContain, eRoutes[1])
				So(aRoutes, ShouldContain, eRoutes[2])
				So(aRoutes, ShouldContain, eRoutes[3])
			})
		})
	})
}

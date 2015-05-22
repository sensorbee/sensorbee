package bql

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
	"time"
)

func TestBqlBox(t *testing.T) {
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
	config := core.Configuration{TupleTraceEnabled: 1}
	ctx := newTestContext(config)

	Convey("Given a simple source/box/sink topology", t, func() {
		/*
		 *   so -*--> b -*--> si
		 */
		tb := core.NewDefaultStaticTopologyBuilder()
		so := &TupleEmitterSource{
			Tuples: []*tuple.Tuple{&tup1, &tup2},
		}
		tb.AddSource("source", so)

		p := parser.NewBQLParser()
		stmt, err := p.ParseStmt("CREATE STREAM box AS SELECT ISTREAM(int) FROM source [RANGE 2 SECONDS] WHERE int = 2")
		So(err, ShouldBeNil)
		createStreamStmt, ok := stmt.(parser.CreateStreamStmt)
		So(ok, ShouldBeTrue)
		b := NewBqlBox(&createStreamStmt)
		tb.AddBox("box", b).Input("source")

		si := &TupleCollectorSink{}
		tb.AddSink("sink", si).Input("box")

		t, err := tb.Build()
		So(err, ShouldBeNil)

		Convey("When a tuple is emitted by the source", func() {
			t.Run(ctx)
			Convey("Then the sink receives the same object", func() {
				So(si.Tuples, ShouldNotBeNil)
				So(len(si.Tuples), ShouldEqual, 2)

				Convey("And the InputName is set to \"output\"", func() {
					So(si.Tuples[0].InputName, ShouldEqual, "output")
				})
			})
		})
	})
}

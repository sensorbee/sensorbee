package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleEmitter(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains only an ISTREAM item", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.AssembleEmitter(6, 6)

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(len(comp.EmitIntervals), ShouldEqual, 0)
					})
				})
			})
		})
		Convey("When the stack contains only an ISTREAM item and an empty interval spec", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.AssembleEmitter(6, 8)

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(len(comp.EmitIntervals), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the stack contains an ISTREAM item and correct interval specs", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 8,
				StreamEmitIntervalAST{IntervalAST{NumericLiteral{2}, Tuples}, Stream{ActualStream, "x", nil}})
			ps.PushComponent(8, 9,
				StreamEmitIntervalAST{IntervalAST{NumericLiteral{3}, Tuples}, Stream{ActualStream, "y", nil}})
			ps.AssembleEmitter(6, 9)

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 9)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(len(comp.EmitIntervals), ShouldEqual, 2)
						So(comp.EmitIntervals[0].Value, ShouldEqual, 2)
						So(comp.EmitIntervals[0].Unit, ShouldEqual, Tuples)
						So(comp.EmitIntervals[0].Name, ShouldEqual, "x")
						So(comp.EmitIntervals[1].Value, ShouldEqual, 3)
						So(comp.EmitIntervals[1].Unit, ShouldEqual, Tuples)
						So(comp.EmitIntervals[1].Name, ShouldEqual, "y")
					})
				})
			})
		})

		Convey("When the stack contains a wrong emitter item and correct interval specs", func() {
			ps.PushComponent(0, 4, Raw{"wrong"})
			ps.PushComponent(6, 8,
				StreamEmitIntervalAST{IntervalAST{NumericLiteral{2}, Tuples}, Stream{ActualStream, "x", nil}})
			ps.PushComponent(8, 9,
				StreamEmitIntervalAST{IntervalAST{NumericLiteral{3}, Tuples}, Stream{ActualStream, "y", nil}})
			f := func() {
				ps.AssembleEmitter(6, 9)
			}

			Convey("Then AssembleEmitter should panic", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When the stack contains an ISTREAM item and wrong interval specs", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 8, Raw{"wrong"})
			ps.PushComponent(8, 9,
				StreamEmitIntervalAST{IntervalAST{NumericLiteral{3}, Tuples}, Stream{ActualStream, "y", nil}})
			f := func() {
				ps.AssembleEmitter(6, 9)
			}

			Convey("Then AssembleEmitter should panic", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When using ISTREAM without a specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.EmitterType, ShouldEqual, Istream)
				So(comp.EmitIntervals, ShouldEqual, nil)
			})
		})

		Convey("When using ISTREAM with a SECONDS specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 3 SECONDS] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.EmitterType, ShouldEqual, Istream)
				So(comp.EmitIntervals, ShouldResemble, []StreamEmitIntervalAST{
					{IntervalAST{NumericLiteral{3}, Seconds}, Stream{ActualStream, "*", nil}},
				})
			})
		})

		Convey("When using ISTREAM with a TUPLES specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 3 TUPLES] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.EmitterType, ShouldEqual, Istream)
				So(comp.EmitIntervals, ShouldResemble, []StreamEmitIntervalAST{
					{IntervalAST{NumericLiteral{3}, Tuples}, Stream{ActualStream, "*", nil}},
				})
			})
		})

		Convey("When using ISTREAM with a TUPLES FROM specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 3 TUPLES IN x] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.EmitterType, ShouldEqual, Istream)
				So(comp.EmitIntervals, ShouldResemble, []StreamEmitIntervalAST{
					{IntervalAST{NumericLiteral{3}, Tuples}, Stream{ActualStream, "x", nil}},
				})
			})
		})

		Convey("When using ISTREAM with multiple TUPLES FROM specifiers", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 3 TUPLES IN x, 2 TUPLES IN y] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.EmitterType, ShouldEqual, Istream)
				So(comp.EmitIntervals, ShouldResemble, []StreamEmitIntervalAST{
					{IntervalAST{NumericLiteral{3}, Tuples}, Stream{ActualStream, "x", nil}},
					{IntervalAST{NumericLiteral{2}, Tuples}, Stream{ActualStream, "y", nil}},
				})
			})
		})
	})
}

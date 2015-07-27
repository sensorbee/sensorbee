package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleWindowedFrom(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only AliasedStreamWindows in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 8, AliasedStreamWindowAST{
				StreamWindowAST{Stream{ActualStream, "a", nil}, IntervalAST{NumericLiteral{3}, Tuples}}, "",
			})
			ps.PushComponent(8, 10, AliasedStreamWindowAST{
				StreamWindowAST{Stream{ActualStream, "b", nil}, IntervalAST{NumericLiteral{2}, Seconds}}, "",
			})
			ps.AssembleWindowedFrom(6, 10)

			Convey("Then AssembleWindowedFrom transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedFromAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, WindowedFromAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(WindowedFromAST)
						So(len(comp.Relations), ShouldEqual, 2)
						So(comp.Relations[0].Name, ShouldEqual, "a")
						So(comp.Relations[0].Value, ShouldEqual, 3)
						So(comp.Relations[0].Unit, ShouldEqual, Tuples)
						So(comp.Relations[0].Alias, ShouldEqual, "")
						So(comp.Relations[1].Name, ShouldEqual, "b")
						So(comp.Relations[1].Value, ShouldEqual, 2)
						So(comp.Relations[1].Unit, ShouldEqual, Seconds)
						So(comp.Relations[1].Alias, ShouldEqual, "")
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleWindowedFrom(6, 8)

			Convey("Then AssembleWindowedFrom pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedFromAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, WindowedFromAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(WindowedFromAST)
						So(len(comp.Relations), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleWindowedFrom(6, 6)

			Convey("Then AssembleWindowedFrom pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedFromAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, WindowedFromAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(WindowedFromAST)
						So(len(comp.Relations), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the stack contains non-AliasedStreamWindows in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			f := func() {
				ps.AssembleWindowedFrom(0, 8)
			}

			Convey("Then AssembleWindowedFrom panics", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting without a FROM", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM a, b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				s := top.(CreateStreamAsSelectStmt).Select
				So(s.WindowedFromAST, ShouldResemble, WindowedFromAST{})
			})
		})

		Convey("When selecting with a FROM", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM a, b FROM c [RANGE 3 TUPLES] AS x, d [RANGE 2 SECONDS]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt).Select
				So(len(comp.Relations), ShouldEqual, 2)
				So(len(comp.Relations), ShouldEqual, 2)
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(comp.Relations[0].Value, ShouldEqual, 3)
				So(comp.Relations[0].Unit, ShouldEqual, Tuples)
				So(comp.Relations[0].Alias, ShouldEqual, "x")
				So(comp.Relations[1].Name, ShouldEqual, "d")
				So(comp.Relations[1].Value, ShouldEqual, 2)
				So(comp.Relations[1].Unit, ShouldEqual, Seconds)
				So(comp.Relations[1].Alias, ShouldEqual, "")
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleWindowedRelation(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 8, Relation{"a"})
			ps.PushComponent(8, 10, RangeAST{NumericLiteral{2}, Seconds})
			ps.AssembleWindowedRelation()

			Convey("Then AssembleWindowedRelation transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedRelationAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, WindowedRelationAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(WindowedRelationAST)
						So(comp.Name, ShouldEqual, "a")
						So(comp.Value, ShouldEqual, 2)
						So(comp.Unit, ShouldEqual, Seconds)
					})
				})
			})
		})

		Convey("When the stack contains one correct item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 8, Relation{"a"})
			ps.AssembleWindowedRelation()

			Convey("Then AssembleWindowedRelation transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedRelationAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, WindowedRelationAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(WindowedRelationAST)
						So(comp.Name, ShouldEqual, "a")
						So(comp.Value, ShouldEqual, 0)
						So(comp.Unit, ShouldEqual, Unspecified)
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleWindowedRelation panics", func() {
				So(ps.AssembleWindowedRelation, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleWindowedRelation panics", func() {
				So(ps.AssembleWindowedRelation, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting with a FROM", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM a, b FROM c [RANGE 3 TUPLES]"
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
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(comp.Relations[0].Value, ShouldEqual, 3)
				So(comp.Relations[0].Unit, ShouldEqual, Tuples)
				So(comp.Relations[0].Alias, ShouldEqual, "")
			})
		})
	})
}

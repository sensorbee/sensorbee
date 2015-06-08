package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleExpressions(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only RowValues in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, RowValue{"a"})
			ps.PushComponent(7, 8, RowValue{"b"})
			ps.AssembleExpressions(6, 8)

			Convey("Then AssembleExpressions transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ExpressionsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, ExpressionsAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(ExpressionsAST)
						So(len(comp.Expressions), ShouldEqual, 2)
						So(comp.Expressions[0], ShouldResemble, RowValue{"a"})
						So(comp.Expressions[1], ShouldResemble, RowValue{"b"})
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleExpressions(6, 8)

			Convey("Then AssembleExpressions pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ExpressionsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, ExpressionsAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(ExpressionsAST)
						So(len(comp.Expressions), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleExpressions(6, 6)

			Convey("Then AssembleExpressions pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ExpressionsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, ExpressionsAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(ExpressionsAST)
						So(len(comp.Expressions), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting a function of multiple columns", func() {
			p.Buffer = "SELECT f(a, b)"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.Projections), ShouldEqual, 1)
				So(s.Projections[0], ShouldHaveSameTypeAs, FuncAppAST{})
				fa := s.Projections[0].(FuncAppAST)
				So(len(fa.Expressions), ShouldEqual, 2)
				So(fa.Expressions[0], ShouldResemble, RowValue{"a"})
				So(fa.Expressions[1], ShouldResemble, RowValue{"b"})
			})
		})
	})
}

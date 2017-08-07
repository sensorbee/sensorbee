package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleExpressionCase(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only WhenThenPairs in the given range (and no ELSE)", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, RowValue{"", "x"})
			ps.PushComponent(6, 7, StringLiteral{"foo"})
			ps.PushComponent(7, 8, RowValue{"", "a"})
			ps.AssembleWhenThenPair()
			ps.PushComponent(8, 9, StringLiteral{"bar"})
			ps.PushComponent(9, 10, RowValue{"", "b"})
			ps.AssembleWhenThenPair()
			ps.AssembleExpressionCase(6, 10)

			Convey("Then AssembleExpressionCase transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ExpressionCaseAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, ExpressionCaseAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(ExpressionCaseAST)
						So(comp.Expr, ShouldResemble, RowValue{"", "x"})
						So(len(comp.Checks), ShouldEqual, 2)
						So(comp.Checks[0], ShouldResemble,
							WhenThenPairAST{StringLiteral{"foo"}, RowValue{"", "a"}})
						So(comp.Checks[1], ShouldResemble,
							WhenThenPairAST{StringLiteral{"bar"}, RowValue{"", "b"}})
						So(comp.Else, ShouldBeNil)
					})
				})
			})
		})

		Convey("When the stack contains only WhenThenPairs in the given range and an ELSE", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, RowValue{"", "x"})
			ps.PushComponent(6, 7, StringLiteral{"foo"})
			ps.PushComponent(7, 8, RowValue{"", "a"})
			ps.AssembleWhenThenPair()
			ps.PushComponent(8, 9, StringLiteral{"bar"})
			ps.PushComponent(9, 10, RowValue{"", "b"})
			ps.AssembleWhenThenPair()
			ps.PushComponent(10, 12, RowValue{"", "y"})
			ps.AssembleExpressionCase(6, 10)

			Convey("Then AssembleExpressionCase transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ExpressionCaseAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 12)
					So(top.comp, ShouldHaveSameTypeAs, ExpressionCaseAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(ExpressionCaseAST)
						So(comp.Expr, ShouldResemble, RowValue{"", "x"})
						So(len(comp.Checks), ShouldEqual, 2)
						So(comp.Checks[0], ShouldResemble,
							WhenThenPairAST{StringLiteral{"foo"}, RowValue{"", "a"}})
						So(comp.Checks[1], ShouldResemble,
							WhenThenPairAST{StringLiteral{"bar"}, RowValue{"", "b"}})
						So(comp.Else, ShouldResemble, RowValue{"", "y"})
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, RowValue{"", "x"})

			Convey("Then AssembleExpressionCase panics", func() {
				f := func() {
					ps.AssembleExpressionCase(6, 6)
				}
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting a function of multiple columns", func() {
			p.Buffer = "SELECT ISTREAM " +
				`CASE x WHEN "foo" THEN a WHEN "bar" THEN b ELSE y END, ` +
				`CASE z WHEN "hoge" THEN x END`
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.Projections), ShouldEqual, 2)

				So(s.Projections[0], ShouldHaveSameTypeAs, ExpressionCaseAST{})
				comp := s.Projections[0].(ExpressionCaseAST)
				So(comp.Expr, ShouldResemble, RowValue{"", "x"})
				So(len(comp.Checks), ShouldEqual, 2)
				So(comp.Checks[0], ShouldResemble,
					WhenThenPairAST{StringLiteral{"foo"}, RowValue{"", "a"}})
				So(comp.Checks[1], ShouldResemble,
					WhenThenPairAST{StringLiteral{"bar"}, RowValue{"", "b"}})
				So(comp.Else, ShouldResemble, RowValue{"", "y"})

				So(s.Projections[1], ShouldHaveSameTypeAs, ExpressionCaseAST{})
				comp = s.Projections[1].(ExpressionCaseAST)
				So(comp.Expr, ShouldResemble, RowValue{"", "z"})
				So(len(comp.Checks), ShouldEqual, 1)
				So(comp.Checks[0], ShouldResemble,
					WhenThenPairAST{StringLiteral{"hoge"}, RowValue{"", "x"}})
				So(comp.Else, ShouldBeNil)

				Convey("And String() should return the original statement", func() {
					So(s.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

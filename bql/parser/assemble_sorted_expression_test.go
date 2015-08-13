package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleSortedExpression(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, NumericLiteral{2})
			ps.PushComponent(7, 8, Yes)
			ps.AssembleSortedExpression()

			Convey("Then AssembleSortedExpression replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SortedExpressionAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, SortedExpressionAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(SortedExpressionAST)
						So(comp.Expr, ShouldResemble, NumericLiteral{2})
						So(comp.Ascending, ShouldResemble, Yes)
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleSortedExpression panics", func() {
				So(ps.AssembleSortedExpression, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleSortedExpression panics", func() {
				So(ps.AssembleSortedExpression, ShouldPanic)
			})
		})
	})
}

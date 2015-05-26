package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleFuncApp(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, FuncName("add"))
			ps.PushComponent(7, 8, ExpressionsAST{[]interface{}{
				NumericLiteral{2},
				ColumnName{"a"}}})
			ps.AssembleFuncApp()

			Convey("Then AssembleFuncApp replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a FuncAppAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, FuncAppAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(FuncAppAST)
						So(comp.Function, ShouldEqual, "add")
						So(len(comp.Expressions), ShouldEqual, 2)
						So(comp.Expressions[0], ShouldResemble, NumericLiteral{2})
						So(comp.Expressions[1], ShouldResemble, ColumnName{"a"})
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleFuncApp panics", func() {
				So(ps.AssembleFuncApp, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleFuncApp panics", func() {
				So(ps.AssembleFuncApp, ShouldPanic)
			})
		})
	})
}

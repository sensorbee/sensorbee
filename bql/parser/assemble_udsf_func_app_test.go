package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleUDSFFuncApp(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains a correct item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, FuncName("add"))
			ps.PushComponent(7, 8, ExpressionsAST{[]Expression{
				NumericLiteral{2},
				RowValue{"", "a"}}})
			ps.AssembleFuncApp()
			ps.AssembleUDSFFuncApp()

			Convey("Then AssembleFuncApp replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Stream", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Stream{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(Stream)
						So(comp.Type, ShouldEqual, USDFStream)
						So(comp.Name, ShouldEqual, "add")
						So(len(comp.Params), ShouldEqual, 2)
						So(comp.Params[0], ShouldResemble, NumericLiteral{2})
						So(comp.Params[1], ShouldResemble, RowValue{"", "a"})
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleUDSFFuncApp panics", func() {
				So(ps.AssembleUDSFFuncApp, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleUDSFFuncApp panics", func() {
				So(ps.AssembleUDSFFuncApp, ShouldPanic)
			})
		})
	})
}

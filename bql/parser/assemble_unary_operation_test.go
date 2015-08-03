package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleUnaryPrefixOperation(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When there is one item in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.AssembleUnaryPrefixOperation(2, 3)

			Convey("Then AssembleUnaryPrefixOperation does nothing to the stack", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 2)
				So(top.end, ShouldEqual, 3)
				So(top.comp, ShouldResemble, RowValue{"", "a"})
			})
		})

		Convey("When there are two correct items in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(3, 4, UnaryMinus)
			ps.PushComponent(4, 5, RowValue{"", "b"})
			ps.AssembleUnaryPrefixOperation(3, 5)

			Convey("Then AssembleUnaryPrefixOperation adds the given operator", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 3)
				So(top.end, ShouldEqual, 5)
				So(top.comp, ShouldHaveSameTypeAs, UnaryOpAST{})
				comp := top.comp.(UnaryOpAST)
				So(comp.Op, ShouldEqual, UnaryMinus)
				So(comp.Expr, ShouldResemble, RowValue{"", "b"})
			})
		})

		Convey("When there are no items in the given range", func() {
			ps.PushComponent(2, 3, RowValue{"", "a"})
			f := func() {
				ps.AssembleUnaryPrefixOperation(4, 5)
			}

			Convey("Then AssembleUnaryPrefixOperation panics", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When there are wrong items in the given range", func() {
			ps.PushComponent(4, 5, RowValue{"", "b"})
			ps.PushComponent(6, 7, RowValue{"", "c"})
			f := func() {
				ps.AssembleUnaryPrefixOperation(4, 7)
			}

			Convey("Then AssembleUnaryPrefixOperation panics", func() {
				So(f, ShouldPanic)
			})
		})
	})
}

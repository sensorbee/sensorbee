package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleTypeCast(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When there is one item in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.AssembleTypeCast(2, 3)

			Convey("Then AssembleTypeCast does nothing to the stack", func() {
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
			ps.PushComponent(3, 4, RowValue{"", "b"})
			ps.PushComponent(4, 5, Int)
			ps.AssembleTypeCast(3, 5)

			Convey("Then AssembleTypeCast adds the given operator", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 3)
				So(top.end, ShouldEqual, 5)
				So(top.comp, ShouldHaveSameTypeAs, TypeCastAST{})
				comp := top.comp.(TypeCastAST)
				So(comp.Expr, ShouldResemble, RowValue{"", "b"})
				So(comp.Target, ShouldEqual, Int)
			})
		})

		Convey("When there are no items in the given range", func() {
			ps.PushComponent(2, 3, RowValue{"", "a"})
			f := func() {
				ps.AssembleTypeCast(4, 5)
			}

			Convey("Then AssembleTypeCast panics", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When there are wrong items in the given range", func() {
			ps.PushComponent(4, 5, RowValue{"", "b"})
			ps.PushComponent(6, 7, RowValue{"", "c"})
			f := func() {
				ps.AssembleTypeCast(4, 7)
			}

			Convey("Then AssembleTypeCast panics", func() {
				So(f, ShouldPanic)
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleBinaryOperation(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When there is one item in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(2, 3, ColumnName{"a"})
			ps.AssembleBinaryOperation(2, 3, "+")

			Convey("Then AssembleBinaryOperation does nothing to the stack", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 2)
				So(top.end, ShouldEqual, 3)
				So(top.comp, ShouldResemble, ColumnName{"a"})
			})
		})

		Convey("When there are two items in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(2, 3, ColumnName{"a"})
			ps.PushComponent(4, 5, ColumnName{"b"})
			ps.AssembleBinaryOperation(2, 5, "+")

			Convey("Then AssembleBinaryOperation adds the given operator", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 2)
				So(top.end, ShouldEqual, 5)
				So(top.comp, ShouldHaveSameTypeAs, BinaryOp{})
				comp := top.comp.(BinaryOp)
				So(comp.op, ShouldEqual, "+")
				So(comp.left, ShouldResemble, ColumnName{"a"})
				So(comp.right, ShouldResemble, ColumnName{"b"})
			})
		})

		Convey("When there are no items in the given range", func() {
			ps.PushComponent(2, 3, ColumnName{"a"})
			f := func() {
				ps.AssembleBinaryOperation(4, 5, "+")
			}

			Convey("Then AssembleBinaryOperation panics", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When there are more than two items in the given range", func() {
			ps.PushComponent(2, 3, ColumnName{"a"})
			ps.PushComponent(4, 5, ColumnName{"b"})
			ps.PushComponent(6, 7, ColumnName{"c"})
			f := func() {
				ps.AssembleBinaryOperation(2, 7, "+")
			}

			Convey("Then AssembleBinaryOperation adds the given operator", func() {
				So(f, ShouldPanic)
			})
		})
	})
}

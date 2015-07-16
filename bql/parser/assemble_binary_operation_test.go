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
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.AssembleBinaryOperation(2, 3)

			Convey("Then AssembleBinaryOperation does nothing to the stack", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 2)
				So(top.end, ShouldEqual, 3)
				So(top.comp, ShouldResemble, RowValue{"", "a"})
			})
		})

		Convey("When there are three correct items in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.PushComponent(3, 4, Plus)
			ps.PushComponent(4, 5, RowValue{"", "b"})
			ps.AssembleBinaryOperation(2, 5)

			Convey("Then AssembleBinaryOperation adds the given operator", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 2)
				So(top.end, ShouldEqual, 5)
				So(top.comp, ShouldHaveSameTypeAs, BinaryOpAST{})
				comp := top.comp.(BinaryOpAST)
				So(comp.Op, ShouldEqual, Plus)
				So(comp.Left, ShouldResemble, RowValue{"", "a"})
				So(comp.Right, ShouldResemble, RowValue{"", "b"})
			})
		})

		Convey("When there are five correct items in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.PushComponent(3, 4, Plus)
			ps.PushComponent(4, 5, RowValue{"", "b"})
			ps.PushComponent(5, 6, Minus)
			ps.PushComponent(6, 7, RowValue{"", "c"})
			ps.AssembleBinaryOperation(2, 7)

			Convey("Then AssembleBinaryOperation adds the given operator", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 2)
				So(top.end, ShouldEqual, 7)
				So(top.comp, ShouldResemble,
					BinaryOpAST{Minus,
						BinaryOpAST{Plus, RowValue{"", "a"}, RowValue{"", "b"}},
						RowValue{"", "c"}})
			})
		})

		Convey("When there are seven correct items in the given range", func() {
			ps.PushComponent(0, 2, Raw{"PRE"})
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.PushComponent(3, 4, Plus)
			ps.PushComponent(4, 5, RowValue{"", "b"})
			ps.PushComponent(5, 6, Minus)
			ps.PushComponent(6, 7, RowValue{"", "c"})
			ps.PushComponent(7, 8, Plus)
			ps.PushComponent(8, 9, RowValue{"", "d"})
			ps.AssembleBinaryOperation(2, 7)

			Convey("Then AssembleBinaryOperation adds the given operator", func() {
				So(ps.Len(), ShouldEqual, 2)
				top := ps.Peek()
				So(top, ShouldNotBeNil)
				So(top.begin, ShouldEqual, 2)
				So(top.end, ShouldEqual, 7)
				So(top.comp, ShouldResemble,
					BinaryOpAST{Plus,
						BinaryOpAST{Minus,
							BinaryOpAST{Plus, RowValue{"", "a"}, RowValue{"", "b"}},
							RowValue{"", "c"}},
						RowValue{"", "d"}})
			})
		})

		Convey("When there are no items in the given range", func() {
			ps.PushComponent(2, 3, RowValue{"", "a"})
			f := func() {
				ps.AssembleBinaryOperation(4, 5)
			}

			Convey("Then AssembleBinaryOperation panics", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When there are not enough items in the given range", func() {
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.PushComponent(6, 7, RowValue{"", "c"})
			f := func() {
				ps.AssembleBinaryOperation(2, 7)
			}

			Convey("Then AssembleBinaryOperation panics", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When there are wrong items in the given range", func() {
			ps.PushComponent(2, 3, RowValue{"", "a"})
			ps.PushComponent(4, 5, RowValue{"", "b"})
			ps.PushComponent(6, 7, RowValue{"", "c"})
			f := func() {
				ps.AssembleBinaryOperation(2, 7)
			}

			Convey("Then AssembleBinaryOperation panics", func() {
				So(f, ShouldPanic)
			})
		})
	})
}

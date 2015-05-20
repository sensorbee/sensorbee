package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleRange(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, Raw{"2"})
			ps.PushComponent(7, 8, RangeUnit{"SECONDS"})
			ps.AssembleRange()

			Convey("Then AssembleRange replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Range", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Range{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(Range)
						So(comp.expr, ShouldEqual, "2")
						So(comp.unit, ShouldEqual, "SECONDS")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleRange panics", func() {
				So(ps.AssembleRange, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleRange panics", func() {
				So(ps.AssembleRange, ShouldPanic)
			})
		})
	})
}

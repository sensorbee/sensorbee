package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleStreamEmitInterval(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, IntervalAST{NumericLiteral{2}, Seconds})
			ps.PushComponent(7, 8, Stream{ActualStream, "x", nil})
			ps.AssembleStreamEmitInterval()

			Convey("Then AssembleStreamEmitInterval replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a StreamEmitIntervalAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, StreamEmitIntervalAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(StreamEmitIntervalAST)
						So(comp.Value, ShouldEqual, 2)
						So(comp.Unit, ShouldEqual, Seconds)
						So(comp.Name, ShouldEqual, "x")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleStreamEmitInterval panics", func() {
				So(ps.AssembleStreamEmitInterval, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleStreamEmitInterval panics", func() {
				So(ps.AssembleStreamEmitInterval, ShouldPanic)
			})
		})
	})
}

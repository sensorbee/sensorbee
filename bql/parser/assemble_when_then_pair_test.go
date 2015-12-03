package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleWhenThenPair(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleWhenThenPair()

			Convey("Then AssembleWhenThenPair replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WhenThenPairAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, WhenThenPairAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(WhenThenPairAST)
						So(comp.When, ShouldResemble, RowValue{"", "a"})
						So(comp.Then, ShouldResemble, RowValue{"", "b"})
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleWhenThenPair panics", func() {
				So(ps.AssembleWhenThenPair, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleWhenThenPair panics", func() {
				So(ps.AssembleWhenThenPair, ShouldPanic)
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleKeyValuePair(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, StringLiteral{"key"})
			ps.PushComponent(7, 8, RowValue{"", "a"})
			ps.AssembleKeyValuePair()

			Convey("Then AssembleKeyValuePair replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a KeyValuePairAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, KeyValuePairAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(KeyValuePairAST)
						So(comp.Key, ShouldEqual, "key")
						So(comp.Value, ShouldResemble, RowValue{"", "a"})
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleKeyValuePair panics", func() {
				So(ps.AssembleKeyValuePair, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleKeyValuePair panics", func() {
				So(ps.AssembleKeyValuePair, ShouldPanic)
			})
		})
	})
}

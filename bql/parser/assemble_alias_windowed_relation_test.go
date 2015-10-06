package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleAliasedStreamWindow(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, StreamWindowAST{Stream{ActualStream, "a", nil},
				IntervalAST{FloatLiteral{2}, Seconds}, 2})
			ps.PushComponent(7, 8, Identifier("out"))
			ps.AssembleAliasedStreamWindow()

			Convey("Then AssembleAliasedStreamWindow replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a AliasedStreamWindowAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, AliasedStreamWindowAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(AliasedStreamWindowAST)
						So(comp.StreamWindowAST, ShouldResemble,
							StreamWindowAST{Stream{ActualStream, "a", nil},
								IntervalAST{FloatLiteral{2}, Seconds}, 2})
						So(comp.Alias, ShouldEqual, "out")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleAliasedStreamWindow panics", func() {
				So(ps.AssembleAliasedStreamWindow, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleAliasedStreamWindow panics", func() {
				So(ps.AssembleAliasedStreamWindow, ShouldPanic)
			})
		})
	})
}

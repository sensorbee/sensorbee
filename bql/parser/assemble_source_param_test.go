package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleSourceParam(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, SourceParamKey("key"))
			ps.PushComponent(6, 8, SourceParamVal("value"))
			ps.AssembleSourceParam()

			Convey("Then AssembleSourceParam replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SourceParamAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, SourceParamAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SourceParamAST)
						So(comp.Key, ShouldEqual, "key")
						So(comp.Value, ShouldEqual, "value")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleSourceParam panics", func() {
				So(ps.AssembleSourceParam, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleSourceParam panics", func() {
				So(ps.AssembleSourceParam, ShouldPanic)
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
)

func TestAssembleSourceSinkParam(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, SourceSinkParamKey("key"))
			ps.PushComponent(6, 8, StringLiteral{"value"})
			ps.AssembleSourceSinkParam()

			Convey("Then AssembleSourceSinkParam replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SourceSinkParamAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, SourceSinkParamAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SourceSinkParamAST)
						So(comp.Key, ShouldEqual, "key")
						So(comp.Value, ShouldEqual, data.String("value"))
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleSourceSinkParam panics", func() {
				So(ps.AssembleSourceSinkParam, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleSourceSinkParam panics", func() {
				So(ps.AssembleSourceSinkParam, ShouldPanic)
			})
		})
	})
}

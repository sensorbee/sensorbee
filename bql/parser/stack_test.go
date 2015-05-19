package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParseStack(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}
		Convey("When the stack is empty", func() {
			Convey("And we push a component", func() {
				ps.PushComponent(1, 2, Raw{"a"})
				Convey("Then there is be one item on the stack", func() {
					So(ps.Len(), ShouldEqual, 1)
				})
			})
			Convey("And we push a component with bad begin/end bounds", func() {
				f := func() {
					ps.PushComponent(2, 1, Raw{"a"})
				}
				Convey("Then the call panics", func() {
					So(f, ShouldPanic)
				})
			})
		})
		Convey("When the stack contains one element", func() {
			ps.PushComponent(1, 2, Raw{"a"})
			Convey("And we push a component", func() {
				ps.PushComponent(2, 3, Raw{"b"})
				Convey("Then there are two items on the stack", func() {
					So(ps.Len(), ShouldEqual, 2)
				})
			})
			Convey("And we push a component with overlapping begin/end bounds", func() {
				f := func() {
					ps.PushComponent(1, 3, Raw{"bc"})
				}
				Convey("Then the call panics", func() {
					So(f, ShouldPanic)
				})
			})
		})
	})
}

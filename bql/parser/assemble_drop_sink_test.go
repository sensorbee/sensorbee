package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleDropSink(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct DROP SINK items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleDropSink()

			Convey("Then AssembleDropSink transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a DropSinkStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, DropSinkStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(DropSinkStmt)
						So(comp.Source, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleDropSink panics", func() {
				So(ps.AssembleDropSink, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full DROP SINK", func() {
			p.Buffer = "DROP SINK a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, DropSinkStmt{})
				comp := top.(DropSinkStmt)

				So(comp.Source, ShouldEqual, "a_1")
			})
		})
	})
}

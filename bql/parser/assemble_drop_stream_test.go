package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleDropStream(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct DROP STREAM items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleDropStream()

			Convey("Then AssembleDropStream transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a DropStreamStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, DropStreamStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(DropStreamStmt)
						So(comp.Source, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleDropStream panics", func() {
				So(ps.AssembleDropStream, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full DROP STREAM", func() {
			p.Buffer = "DROP STREAM a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, DropStreamStmt{})
				comp := top.(DropStreamStmt)

				So(comp.Source, ShouldEqual, "a_1")
			})
		})
	})
}

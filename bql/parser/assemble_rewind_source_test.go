package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleRewindSource(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct REWIND SOURCE items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleRewindSource()

			Convey("Then AssembleRewindSource transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a RewindSourceStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, RewindSourceStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(RewindSourceStmt)
						So(comp.Source, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleRewindSource panics", func() {
				So(ps.AssembleRewindSource, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full REWIND SOURCE", func() {
			p.Buffer = "REWIND SOURCE a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, RewindSourceStmt{})
				comp := top.(RewindSourceStmt)

				So(comp.Source, ShouldEqual, "a_1")
			})
		})
	})
}

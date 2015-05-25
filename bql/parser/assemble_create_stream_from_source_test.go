package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleCreateStreamFromSource(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct CREATE STREAM items", func() {
			ps.PushComponent(2, 4, Relation{"a"})
			ps.PushComponent(4, 6, SourceName("b"))
			ps.AssembleCreateStreamFromSource()

			Convey("Then AssembleCreateStreamFromSource transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a CreateStreamFromSourceStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, CreateStreamFromSourceStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(CreateStreamFromSourceStmt)
						So(comp.Name, ShouldEqual, "a")
						So(comp.Source, ShouldEqual, "b")
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.AssembleProjections(6, 7)
			Convey("Then AssembleCreateStreamFromSource panics", func() {
				So(ps.AssembleCreateStreamFromSource, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be Relation
			ps.PushComponent(4, 6, SourceName("b"))

			Convey("Then AssembleCreateStreamFromSource panics", func() {
				So(ps.AssembleCreateStreamFromSource, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a CREATE STREAM", func() {
			p.Buffer = "CREATE STREAM a FROM SOURCE b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamFromSourceStmt{})
				comp := top.(CreateStreamFromSourceStmt)

				So(comp.Name, ShouldEqual, "a")
				So(comp.Source, ShouldEqual, "b")
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleInsertIntoFrom(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct SELECT items with a Interval specification", func() {
			ps.PushComponent(4, 5, StreamIdentifier("x"))
			ps.PushComponent(5, 6, StreamIdentifier("y"))
			ps.AssembleInsertIntoFrom()

			Convey("Then AssembleInsertIntoFrom transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a InsertIntoFromStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, InsertIntoFromStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(InsertIntoFromStmt)
						So(comp.Sink, ShouldEqual, "x")
						So(comp.Input, ShouldEqual, "y")
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(4, 5, StreamIdentifier("x"))
			Convey("Then AssembleInsertIntoFrom panics", func() {
				So(ps.AssembleInsertIntoFrom, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(4, 5, StreamIdentifier("x"))
			ps.PushComponent(5, 6, Istream) // must be StreamIdentifier
			Convey("Then AssembleInsertIntoFrom panics", func() {
				So(ps.AssembleInsertIntoFrom, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full INSERT INTO FROM", func() {
			// the statement below does not make sense, it is just used to test
			// whether we accept optional RANGE clauses correctly
			p.Buffer = "INSERT INTO x FROM y"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, InsertIntoFromStmt{})
				comp := top.(InsertIntoFromStmt)

				So(comp.Sink, ShouldEqual, "x")
				So(comp.Input, ShouldEqual, "y")
			})
		})
	})
}

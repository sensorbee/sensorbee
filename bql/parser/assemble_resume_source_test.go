package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleResumeSource(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct RESUME SOURCE items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleResumeSource()

			Convey("Then AssembleResumeSource transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a ResumeSourceStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, ResumeSourceStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(ResumeSourceStmt)
						So(comp.Source, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleResumeSource panics", func() {
				So(ps.AssembleResumeSource, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full RESUME SOURCE", func() {
			p.Buffer = "RESUME SOURCE a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, ResumeSourceStmt{})
				comp := top.(ResumeSourceStmt)

				So(comp.Source, ShouldEqual, "a_1")

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

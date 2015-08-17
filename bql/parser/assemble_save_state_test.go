package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssemblSaveState(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct SAVE STATE items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleSaveState()

			Convey("Then AssembleSaveState transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a SaveStateStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, SaveStateStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SaveStateStmt)
						So(comp.Name, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleSaveState panics", func() {
				So(ps.AssembleSaveState, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full SAVE STATE", func() {
			p.Buffer = "SAVE STATE a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SaveStateStmt{})
				comp := top.(SaveStateStmt)

				So(comp.Name, ShouldEqual, "a_1")

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
)

func TestAssembleUpdateState(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct UPDATE STATE items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", data.String("d")})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", data.String("f")})
			ps.AssembleSourceSinkSpecs(6, 10)
			ps.AssembleUpdateState()

			Convey("Then AssembleUpdateState transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a UpdateStateStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, UpdateStateStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(UpdateStateStmt)
						So(comp.Name, ShouldEqual, "a")
						So(len(comp.Params), ShouldEqual, 2)
						So(comp.Params[0].Key, ShouldEqual, "c")
						So(comp.Params[0].Value, ShouldEqual, data.String("d"))
						So(comp.Params[1].Key, ShouldEqual, "e")
						So(comp.Params[1].Value, ShouldEqual, data.String("f"))
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.AssembleProjections(6, 7)
			Convey("Then AssembleUpdateState panics", func() {
				So(ps.AssembleUpdateState, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", data.String("d")})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", data.String("f")})
			ps.AssembleSourceSinkSpecs(6, 10)

			Convey("Then AssembleUpdateState panics", func() {
				So(ps.AssembleUpdateState, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full UPDATE STATE", func() {
			p.Buffer = "UPDATE STATE a_1 SET c=27, e_='f_1', f=[], g={}"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, UpdateStateStmt{})
				comp := top.(UpdateStateStmt)

				So(comp.Name, ShouldEqual, "a_1")
				So(len(comp.Params), ShouldEqual, 4)
				So(comp.Params[0].Key, ShouldEqual, "c")
				So(comp.Params[0].Value, ShouldEqual, data.Int(27))
				So(comp.Params[1].Key, ShouldEqual, "e_")
				So(comp.Params[1].Value, ShouldEqual, data.String("f_1"))
				So(comp.Params[2].Key, ShouldEqual, "f")
				So(comp.Params[2].Value, ShouldResemble, data.Array{})
				So(comp.Params[3].Key, ShouldEqual, "g")
				So(comp.Params[3].Value, ShouldResemble, data.Map{})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

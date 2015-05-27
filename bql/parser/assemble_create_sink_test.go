package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleCreateSink(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct CREATE SINK items", func() {
			ps.PushComponent(2, 4, SourceSinkName("a"))
			ps.PushComponent(4, 6, SourceSinkType("b"))
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", "d"})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", "f"})
			ps.AssembleSourceSinkSpecs(6, 10)
			ps.AssembleCreateSink()

			Convey("Then AssembleCreateSink transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a CreateSinkStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, CreateSinkStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(CreateSinkStmt)
						So(comp.Name, ShouldEqual, "a")
						So(comp.Type, ShouldEqual, "b")
						So(len(comp.Params), ShouldEqual, 2)
						So(comp.Params[0].Key, ShouldEqual, "c")
						So(comp.Params[0].Value, ShouldEqual, "d")
						So(comp.Params[1].Key, ShouldEqual, "e")
						So(comp.Params[1].Value, ShouldEqual, "f")
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.AssembleProjections(6, 7)
			Convey("Then AssembleCreateSink panics", func() {
				So(ps.AssembleCreateSink, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be SourceSinkName
			ps.PushComponent(4, 6, SourceSinkType("b"))
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", "d"})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", "f"})
			ps.AssembleSourceSinkSpecs(6, 10)

			Convey("Then AssembleCreateSink panics", func() {
				So(ps.AssembleCreateSink, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full CREATE SINK", func() {
			p.Buffer = "CREATE SINK a_1 TYPE b WITH c=27, e_=f_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateSinkStmt{})
				comp := top.(CreateSinkStmt)

				So(comp.Name, ShouldEqual, "a_1")
				So(comp.Type, ShouldEqual, "b")
				So(len(comp.Params), ShouldEqual, 2)
				So(comp.Params[0].Key, ShouldEqual, "c")
				So(comp.Params[0].Value, ShouldEqual, "27")
				So(comp.Params[1].Key, ShouldEqual, "e_")
				So(comp.Params[1].Value, ShouldEqual, "f_1")
			})
		})
	})
}

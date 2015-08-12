package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestAssembleLoadStateOrCreate(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct LOAD STATE OR CREATE items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.PushComponent(4, 6, SourceSinkType("b"))
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", data.String("d")})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", data.String("f")})
			ps.AssembleSourceSinkSpecs(6, 10)
			ps.PushComponent(11, 13, SourceSinkParamAST{"g", data.String("h")})
			ps.PushComponent(14, 15, SourceSinkParamAST{"i", data.String("j")})
			ps.AssembleSourceSinkSpecs(11, 15)
			ps.AssembleLoadStateOrCreate()

			Convey("Then AssembleLoadStateOrCreate transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a LoadStateOrCreateStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, LoadStateOrCreateStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(LoadStateOrCreateStmt)
						So(comp.Name, ShouldEqual, "a")
						So(comp.Type, ShouldEqual, "b")
						So(len(comp.Set.Params), ShouldEqual, 2)
						So(comp.Set.Params[0].Key, ShouldEqual, "c")
						So(comp.Set.Params[0].Value, ShouldEqual, data.String("d"))
						So(comp.Set.Params[1].Key, ShouldEqual, "e")
						So(comp.Set.Params[1].Value, ShouldEqual, data.String("f"))
						So(len(comp.With.Params), ShouldEqual, 2)
						So(comp.With.Params[0].Key, ShouldEqual, "g")
						So(comp.With.Params[0].Value, ShouldEqual, data.String("h"))
						So(comp.With.Params[1].Key, ShouldEqual, "i")
						So(comp.With.Params[1].Value, ShouldEqual, data.String("j"))
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.AssembleProjections(6, 7)
			Convey("Then AssembleLoadStateOrCreate panics", func() {
				So(ps.AssembleLoadStateOrCreate, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier
			ps.PushComponent(4, 6, SourceSinkType("b"))
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", data.String("d")})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", data.String("f")})
			ps.AssembleSourceSinkSpecs(6, 10)

			Convey("Then AssembleLoadStateOrCreate panics", func() {
				So(ps.AssembleLoadStateOrCreate, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a LOAD STATE OR CREATE without set/with option", func() {
			p.Buffer = "LOAD STATE a_1 TYPE b OR CREATE IF NOT EXISTS"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, LoadStateOrCreateStmt{})
				comp := top.(LoadStateOrCreateStmt)

				So(comp.Name, ShouldEqual, "a_1")
				So(comp.Type, ShouldEqual, "b")
				So(len(comp.Set.Params), ShouldEqual, 0)
				So(len(comp.With.Params), ShouldEqual, 0)

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When doing a full LOAD STATE OR CREATE", func() {
			p.Buffer = "LOAD STATE a_1 TYPE b SET c=27, e_='f_1', f=[7,'g'] OR CREATE IF NOT EXISTS WITH g=2"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, LoadStateOrCreateStmt{})
				comp := top.(LoadStateOrCreateStmt)

				So(comp.Name, ShouldEqual, "a_1")
				So(comp.Type, ShouldEqual, "b")
				So(len(comp.Set.Params), ShouldEqual, 3)
				So(comp.Set.Params[0].Key, ShouldEqual, "c")
				So(comp.Set.Params[0].Value, ShouldEqual, data.Int(27))
				So(comp.Set.Params[1].Key, ShouldEqual, "e_")
				So(comp.Set.Params[1].Value, ShouldEqual, data.String("f_1"))
				So(comp.Set.Params[2].Key, ShouldEqual, "f")
				So(comp.Set.Params[2].Value, ShouldResemble, data.Array{data.Int(7), data.String("g")})
				So(len(comp.With.Params), ShouldEqual, 1)
				So(comp.With.Params[0].Key, ShouldEqual, "g")
				So(comp.With.Params[0].Value, ShouldEqual, data.Int(2))

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

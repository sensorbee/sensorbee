package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
)

func TestAssembleSourceSinkSpecs(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only SourceSinkParams in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, SourceSinkParamAST{"key", data.String("val")})
			ps.PushComponent(7, 8, SourceSinkParamAST{"a", data.String("b")})
			ps.AssembleSourceSinkSpecs(6, 8)

			Convey("Then AssembleSourceSinkSpecs transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SourceSinkSpecsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, SourceSinkSpecsAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SourceSinkSpecsAST)
						So(len(comp.Params), ShouldEqual, 2)
						So(comp.Params[0].Key, ShouldEqual, "key")
						So(comp.Params[0].Value, ShouldEqual, data.String("val"))
						So(comp.Params[1].Key, ShouldEqual, "a")
						So(comp.Params[1].Value, ShouldEqual, data.String("b"))
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleSourceSinkSpecs(6, 8)

			Convey("Then AssembleSourceSinkSpecs pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SourceSinkSpecsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, SourceSinkSpecsAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(SourceSinkSpecsAST)
						So(len(comp.Params), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleSourceSinkSpecs(6, 6)

			Convey("Then AssembleSourceSinkSpecs pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SourceSinkSpecsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, SourceSinkSpecsAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(SourceSinkSpecsAST)
						So(len(comp.Params), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the stack contains non-SourceSinkParams in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			f := func() {
				ps.AssembleSourceSinkSpecs(0, 8)
			}

			Convey("Then AssembleSourceSinkSpecs panics", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When creating a source without a WITH", func() {
			p.Buffer = "CREATE SOURCE a TYPE b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateSourceStmt{})
				s := top.(CreateSourceStmt)
				So(s.Params, ShouldBeNil)

				Convey("And String() should return the original statement", func() {
					So(s.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When creating a source with a WITH", func() {
			p.Buffer = `CREATE SOURCE a TYPE b WITH port=8080, proto='http'`
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateSourceStmt{})
				s := top.(CreateSourceStmt)
				So(s.Params, ShouldNotBeNil)
				So(len(s.Params), ShouldEqual, 2)
				So(s.Params[0], ShouldResemble,
					SourceSinkParamAST{"port", data.Int(8080)})
				So(s.Params[1], ShouldResemble,
					SourceSinkParamAST{"proto", data.String("http")})

				Convey("And String() should return the original statement", func() {
					So(s.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

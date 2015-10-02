package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"strings"
	"testing"
)

func TestAssembleCreateState(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct CREATE SINK items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.PushComponent(4, 6, SourceSinkType("b"))
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", data.String("d")})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", data.String("f")})
			ps.AssembleSourceSinkSpecs(6, 10)
			ps.AssembleCreateState()

			Convey("Then AssembleCreateState transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a CreateStateStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, CreateStateStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(CreateStateStmt)
						So(comp.Name, ShouldEqual, "a")
						So(comp.Type, ShouldEqual, "b")
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
			Convey("Then AssembleCreateState panics", func() {
				So(ps.AssembleCreateState, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier
			ps.PushComponent(4, 6, SourceSinkType("b"))
			ps.PushComponent(6, 8, SourceSinkParamAST{"c", data.String("d")})
			ps.PushComponent(8, 10, SourceSinkParamAST{"e", data.String("f")})
			ps.AssembleSourceSinkSpecs(6, 10)

			Convey("Then AssembleCreateState panics", func() {
				So(ps.AssembleCreateState, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full CREATE STATE", func() {
			p.Buffer = "CREATE STATE a_1 TYPE b WITH c=27, e_='f_1', f=[7,'g']"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStateStmt{})
				comp := top.(CreateStateStmt)

				So(comp.Name, ShouldEqual, "a_1")
				So(comp.Type, ShouldEqual, "b")
				So(len(comp.Params), ShouldEqual, 3)
				So(comp.Params[0].Key, ShouldEqual, "c")
				So(comp.Params[0].Value, ShouldEqual, data.Int(27))
				So(comp.Params[1].Key, ShouldEqual, "e_")
				So(comp.Params[1].Value, ShouldEqual, data.String("f_1"))
				So(comp.Params[2].Key, ShouldEqual, "f")
				So(comp.Params[2].Value, ShouldResemble, data.Array{data.Int(7), data.String("g")})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		// ordering of map's keys are not fixed, and cannot check equality of
		// reversed query with input query, so separate map parameter test.
		Convey("When doing CREATE STATE with map parameter", func() {
			mp1 := "'i':'I_1'"
			mp2 := "'j':false"
			mp3 := "'k':8"
			mapParams := "h={" + strings.Join([]string{mp1, mp2, mp3}, ",") + "}"
			createQuery := "CREATE STATE a_1 TYPE b WITH "
			boolParam := "l=true"
			p.Buffer = createQuery + mapParams + ", " + boolParam
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStateStmt{})
				comp := top.(CreateStateStmt)

				So(comp.Name, ShouldEqual, "a_1")
				So(comp.Type, ShouldEqual, "b")
				So(len(comp.Params), ShouldEqual, 2)
				So(comp.Params[0].Key, ShouldEqual, "h")
				hmap, err := data.AsMap(comp.Params[0].Value)
				So(err, ShouldBeNil)
				So(hmap["i"], ShouldEqual, data.String("I_1"))
				So(hmap["j"], ShouldEqual, data.False)
				So(hmap["k"], ShouldEqual, data.Int(8))
				So(comp.Params[1].Key, ShouldEqual, "l")
				So(comp.Params[1].Value, ShouldEqual, data.True)

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldStartWith, createQuery)
					So(comp.String(), ShouldContainSubstring, mp1)
					So(comp.String(), ShouldContainSubstring, mp2)
					So(comp.String(), ShouldContainSubstring, mp3)
					So(comp.String(), ShouldEndWith, boolParam)
					So(len(comp.String()), ShouldEqual, len(p.Buffer))
				})
			})
		})
	})
}

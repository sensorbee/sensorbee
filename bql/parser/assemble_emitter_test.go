package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleEmitter(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains only an ISTREAM item", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.AssembleEmitterOptions(6, 6)
			ps.AssembleEmitter()

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
					})
				})
			})
		})
		Convey("When the stack contains an ISTREAM item and a LIMIT clause", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 8, NumericLiteral{7})
			ps.AssembleEmitterLimit()
			ps.AssembleEmitterOptions(6, 8)
			ps.AssembleEmitter()

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(comp.EmitterOptions, ShouldResemble, []interface{}{EmitterLimit{7}})
					})
				})
			})
		})

		Convey("When the stack contains a wrong emitter item", func() {
			ps.PushComponent(0, 4, Raw{"wrong"})
			f := func() {
				ps.AssembleEmitter()
			}

			Convey("Then AssembleEmitter should panic", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When using ISTREAM without a specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with a LIMIT specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [LIMIT 7] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{EmitterLimit{7}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleEmitProjections(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 8, ProjectionsAST{[]Expression{RowValue{"", "a"},
				RowValue{"", "b"}}})
			ps.AssembleEmitProjections()

			Convey("Then AssembleEmitProjections replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitProjectionsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, EmitProjectionsAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitProjectionsAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(len(comp.Projections), ShouldEqual, 2)
						So(comp.Projections[0], ShouldResemble, RowValue{"", "a"})
						So(comp.Projections[1], ShouldResemble, RowValue{"", "b"})
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleEmitProjections panics", func() {
				So(ps.AssembleEmitProjections, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleEmitProjections panics", func() {
				So(ps.AssembleEmitProjections, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting multiple columns", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM(a, b)"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				s := top.(CreateStreamAsSelectStmt)
				So(s.EmitterType, ShouldEqual, Istream)
				So(len(s.Projections), ShouldEqual, 2)
				So(s.Projections[0], ShouldResemble, RowValue{"", "a"})
				So(s.Projections[1], ShouldResemble, RowValue{"", "b"})
			})
		})
	})
}

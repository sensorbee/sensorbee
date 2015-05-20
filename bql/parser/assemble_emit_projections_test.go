package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleEmitProjections(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Emitter{"ISTREAM"})
			ps.PushComponent(6, 8, Projections{[]interface{}{ColumnName{"a"},
				ColumnName{"b"}}})
			ps.AssembleEmitProjections()

			Convey("Then AssembleEmitProjections replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitProjections", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, EmitProjections{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitProjections)
						So(comp.emitterType, ShouldEqual, "ISTREAM")
						So(len(comp.projections), ShouldEqual, 2)
						So(comp.projections[0], ShouldResemble, ColumnName{"a"})
						So(comp.projections[1], ShouldResemble, ColumnName{"b"})
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
		p := &Bql{}

		Convey("When selecting multiple columns", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM(a, b)"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamStmt{})
				s := top.(CreateStreamStmt)
				So(s.emitterType, ShouldEqual, "ISTREAM")
				So(len(s.projections), ShouldEqual, 2)
				So(s.projections[0], ShouldResemble, ColumnName{"a"})
				So(s.projections[1], ShouldResemble, ColumnName{"b"})
			})
		})
	})
}

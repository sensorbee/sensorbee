package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleProjections(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only ColumnNames in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.PushComponent(7, 8, ColumnName{"b"})
			ps.AssembleProjections(6, 8)

			Convey("Then AssembleProjections transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Projections", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Projections{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(Projections)
						So(len(comp.projections), ShouldEqual, 2)
						So(comp.projections[0], ShouldResemble, ColumnName{"a"})
						So(comp.projections[1], ShouldResemble, ColumnName{"b"})
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleProjections(6, 8)

			Convey("Then AssembleProjections pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Projections", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Projections{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(Projections)
						So(len(comp.projections), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleProjections(6, 6)

			Convey("Then AssembleProjections pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Projections", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, Projections{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(Projections)
						So(len(comp.projections), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting multiple columns", func() {
			p.Buffer = "SELECT a, b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.projections), ShouldEqual, 2)
				So(s.projections[0], ShouldResemble, ColumnName{"a"})
				So(s.projections[1], ShouldResemble, ColumnName{"b"})
			})
		})
	})
}

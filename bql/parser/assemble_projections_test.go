package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleProjections(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only RowValues in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleProjections(6, 8)

			Convey("Then AssembleProjections transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ProjectionsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, ProjectionsAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(ProjectionsAST)
						So(len(comp.Projections), ShouldEqual, 2)
						So(comp.Projections[0], ShouldResemble, RowValue{"", "a"})
						So(comp.Projections[1], ShouldResemble, RowValue{"", "b"})
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleProjections(6, 8)

			Convey("Then AssembleProjections pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ProjectionsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, ProjectionsAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(ProjectionsAST)
						So(len(comp.Projections), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleProjections(6, 6)

			Convey("Then AssembleProjections pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a ProjectionsAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, ProjectionsAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(ProjectionsAST)
						So(len(comp.Projections), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting multiple columns", func() {
			p.Buffer = "SELECT ISTREAM a, *, b AS c, * AS d"
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
				So(len(s.Projections), ShouldEqual, 4)
				So(s.Projections[0], ShouldResemble, RowValue{"", "a"})
				So(s.Projections[1], ShouldResemble, Wildcard{})
				So(s.Projections[2], ShouldResemble, AliasAST{RowValue{"", "b"}, "c"})
				So(s.Projections[3], ShouldResemble, AliasAST{Wildcard{}, "d"})
			})
		})
	})
}

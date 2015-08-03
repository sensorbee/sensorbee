package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleSelectUnion(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only RowValues in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, SelectStmt{ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}})
			ps.PushComponent(7, 8, SelectStmt{ProjectionsAST: ProjectionsAST{
				[]Expression{RowValue{"", "b"}}}})
			ps.AssembleSelectUnion(6, 8)

			Convey("Then AssembleSelectUnion transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SelectUnionStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, SelectUnionStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SelectUnionStmt)
						So(len(comp.Selects), ShouldEqual, 2)
						So(comp.Selects[0].Projections, ShouldResemble, []Expression{RowValue{"", "a"}})
						So(comp.Selects[1].Projections, ShouldResemble, []Expression{RowValue{"", "b"}})
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleSelectUnion(6, 8)

			Convey("Then AssembleSelectUnion pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SelectUnionStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, SelectUnionStmt{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(SelectUnionStmt)
						So(len(comp.Selects), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleSelectUnion(6, 6)

			Convey("Then AssembleSelectUnion pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a SelectUnionStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, SelectUnionStmt{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(SelectUnionStmt)
						So(len(comp.Selects), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When working with two SELECT statements", func() {
			p.Buffer = "SELECT ISTREAM a UNION ALL SELECT DSTREAM b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectUnionStmt{})
				s := top.(SelectUnionStmt)
				So(len(s.Selects), ShouldEqual, 2)
				So(s.Selects[0].EmitterType, ShouldEqual, Istream)
				So(s.Selects[0].Projections, ShouldResemble, []Expression{RowValue{"", "a"}})
				So(s.Selects[1].EmitterType, ShouldEqual, Dstream)
				So(s.Selects[1].Projections, ShouldResemble, []Expression{RowValue{"", "b"}})
			})
		})

		Convey("When working with more than two SELECT statements", func() {
			p.Buffer = "SELECT ISTREAM a UNION ALL SELECT DSTREAM b UNION ALL SELECT RSTREAM c"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectUnionStmt{})
				s := top.(SelectUnionStmt)
				So(len(s.Selects), ShouldEqual, 3)
				So(s.Selects[0].EmitterType, ShouldEqual, Istream)
				So(s.Selects[0].Projections, ShouldResemble, []Expression{RowValue{"", "a"}})
				So(s.Selects[1].EmitterType, ShouldEqual, Dstream)
				So(s.Selects[1].Projections, ShouldResemble, []Expression{RowValue{"", "b"}})
				So(s.Selects[2].EmitterType, ShouldEqual, Rstream)
				So(s.Selects[2].Projections, ShouldResemble, []Expression{RowValue{"", "c"}})
			})
		})
	})
}

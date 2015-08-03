package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleGrouping(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only RowValues in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleGrouping(6, 8)

			Convey("Then AssembleGrouping transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a GroupingAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, GroupingAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(GroupingAST)
						So(len(comp.GroupList), ShouldEqual, 2)
						So(comp.GroupList[0], ShouldResemble, RowValue{"", "a"})
						So(comp.GroupList[1], ShouldResemble, RowValue{"", "b"})
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleGrouping(6, 8)

			Convey("Then AssembleGrouping pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a GroupingAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, GroupingAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(GroupingAST)
						So(len(comp.GroupList), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleGrouping(6, 6)

			Convey("Then AssembleGrouping pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a GroupingAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, GroupingAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(GroupingAST)
						So(len(comp.GroupList), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting without a GROUP BY", func() {
			p.Buffer = "SELECT ISTREAM a, b"
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
				So(len(s.GroupList), ShouldEqual, 0)
			})
		})

		Convey("When selecting with a GROUP BY", func() {
			p.Buffer = "SELECT ISTREAM a, b GROUP BY c, d"
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
				So(len(s.GroupList), ShouldEqual, 2)
				So(s.GroupList[0], ShouldResemble, RowValue{"", "c"})
				So(s.GroupList[1], ShouldResemble, RowValue{"", "d"})
			})
		})
	})
}

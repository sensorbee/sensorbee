package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleGrouping(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}

		Convey("When the stack contains only ColumnNames in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.PushComponent(7, 8, ColumnName{"b"})
			ps.AssembleGrouping(6, 8)

			Convey("Then AssembleGrouping transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Grouping", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Grouping{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(Grouping)
						So(len(comp.groupList), ShouldEqual, 2)
						So(comp.groupList[0], ShouldResemble, ColumnName{"a"})
						So(comp.groupList[1], ShouldResemble, ColumnName{"b"})
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleGrouping(6, 8)

			Convey("Then AssembleGrouping pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Grouping", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Grouping{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(Grouping)
						So(len(comp.groupList), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleGrouping(6, 6)

			Convey("Then AssembleGrouping pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Grouping", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, Grouping{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(Grouping)
						So(len(comp.groupList), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &Bql{}

		Convey("When selecting without a GROUP BY", func() {
			p.Buffer = "SELECT a, b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.groupList), ShouldEqual, 0)
			})
		})

		Convey("When selecting with a GROUP BY", func() {
			p.Buffer = "SELECT a, b GROUP BY c, d"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.groupList), ShouldEqual, 2)
				So(s.groupList[0], ShouldResemble, ColumnName{"c"})
				So(s.groupList[1], ShouldResemble, ColumnName{"d"})
			})
		})
	})
}

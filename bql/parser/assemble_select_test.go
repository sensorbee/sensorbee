package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleSelect(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}
		Convey("When the stack contains the correct SELECT items", func() {
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.PushComponent(7, 8, ColumnName{"b"})
			ps.AssembleProjections(6, 8)
			ps.PushComponent(12, 13, Relation{"c"})
			ps.PushComponent(13, 14, Relation{"d"})
			ps.AssembleFrom(12, 14)
			ps.PushComponent(14, 15, ColumnName{"e"})
			ps.AssembleFilter(14, 15)
			ps.PushComponent(15, 16, ColumnName{"f"})
			ps.PushComponent(16, 17, ColumnName{"g"})
			ps.AssembleGrouping(15, 17)
			ps.PushComponent(17, 18, ColumnName{"h"})
			ps.AssembleHaving(17, 18)
			ps.AssembleSelect()

			Convey("Then AssembleSelect transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a SelectStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 18)
					So(top.comp, ShouldHaveSameTypeAs, SelectStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SelectStmt)
						So(len(comp.projections), ShouldEqual, 2)
						So(comp.projections[0], ShouldResemble, ColumnName{"a"})
						So(comp.projections[1], ShouldResemble, ColumnName{"b"})
						So(len(comp.relations), ShouldEqual, 2)
						So(comp.relations[0].name, ShouldEqual, "c")
						So(comp.relations[1].name, ShouldEqual, "d")
						So(comp.filter, ShouldResemble, ColumnName{"e"})
						So(len(comp.groupList), ShouldEqual, 2)
						So(comp.groupList[0], ShouldResemble, ColumnName{"f"})
						So(comp.groupList[1], ShouldResemble, ColumnName{"g"})
						So(comp.having, ShouldResemble, ColumnName{"h"})
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.AssembleProjections(6, 7)
			Convey("Then AssembleSelect panics", func() {
				So(ps.AssembleSelect, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.PushComponent(7, 8, ColumnName{"b"})
			ps.AssembleProjections(6, 8)
			ps.PushComponent(12, 13, Relation{"c"})
			ps.PushComponent(13, 14, Relation{"d"})
			ps.AssembleFrom(12, 14)
			ps.PushComponent(14, 15, ColumnName{"e"})
			ps.AssembleFilter(14, 15)
			ps.PushComponent(15, 16, ColumnName{"f"})
			ps.PushComponent(16, 17, ColumnName{"g"})
			ps.AssembleGrouping(15, 17)
			ps.PushComponent(17, 18, ColumnName{"h"})
			ps.AssembleFilter(17, 18) // must be HAVING in correct stmt
			Convey("Then AssembleSelect panics", func() {
				So(ps.AssembleSelect, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &Bql{}

		Convey("When doing a full SELECT", func() {
			p.Buffer = "SELECT a, b FROM c, d WHERE e GROUP BY f, g HAVING h"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				comp := top.(SelectStmt)

				So(len(comp.projections), ShouldEqual, 2)
				So(comp.projections[0], ShouldResemble, ColumnName{"a"})
				So(comp.projections[1], ShouldResemble, ColumnName{"b"})
				So(len(comp.relations), ShouldEqual, 2)
				So(comp.relations[0].name, ShouldEqual, "c")
				So(comp.relations[1].name, ShouldEqual, "d")
				So(comp.filter, ShouldResemble, ColumnName{"e"})
				So(len(comp.groupList), ShouldEqual, 2)
				So(comp.groupList[0], ShouldResemble, ColumnName{"f"})
				So(comp.groupList[1], ShouldResemble, ColumnName{"g"})
				So(comp.having, ShouldResemble, ColumnName{"h"})
			})
		})
	})
}

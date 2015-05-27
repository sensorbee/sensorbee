package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleInsertIntoSelect(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct SELECT items", func() {
			ps.PushComponent(4, 6, SourceSinkName("x"))
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
			ps.AssembleInsertIntoSelect()

			Convey("Then AssembleInsertIntoSelect transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a InsertIntoSelectStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 18)
					So(top.comp, ShouldHaveSameTypeAs, InsertIntoSelectStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(InsertIntoSelectStmt)
						So(comp.Sink, ShouldEqual, "x")
						So(len(comp.Projections), ShouldEqual, 2)
						So(comp.Projections[0], ShouldResemble, ColumnName{"a"})
						So(comp.Projections[1], ShouldResemble, ColumnName{"b"})
						So(len(comp.Relations), ShouldEqual, 2)
						So(comp.Relations[0].Name, ShouldEqual, "c")
						So(comp.Relations[1].Name, ShouldEqual, "d")
						So(comp.Filter, ShouldResemble, ColumnName{"e"})
						So(len(comp.GroupList), ShouldEqual, 2)
						So(comp.GroupList[0], ShouldResemble, ColumnName{"f"})
						So(comp.GroupList[1], ShouldResemble, ColumnName{"g"})
						So(comp.Having, ShouldResemble, ColumnName{"h"})
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.AssembleProjections(6, 7)
			Convey("Then AssembleInsertIntoSelect panics", func() {
				So(ps.AssembleInsertIntoSelect, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(4, 6, SourceSinkName("x"))
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
			Convey("Then AssembleInsertIntoSelect panics", func() {
				So(ps.AssembleInsertIntoSelect, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full INSERT INTO SELECT", func() {
			p.Buffer = "INSERT INTO x SELECT '日本語', b FROM c, d WHERE e GROUP BY f, g HAVING h"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, InsertIntoSelectStmt{})
				comp := top.(InsertIntoSelectStmt)

				So(comp.Sink, ShouldEqual, "x")
				So(len(comp.Projections), ShouldEqual, 2)
				So(comp.Projections[0], ShouldResemble, StringLiteral{"日本語"})
				So(comp.Projections[1], ShouldResemble, ColumnName{"b"})
				So(len(comp.Relations), ShouldEqual, 2)
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(comp.Relations[1].Name, ShouldEqual, "d")
				So(comp.Filter, ShouldResemble, ColumnName{"e"})
				So(len(comp.GroupList), ShouldEqual, 2)
				So(comp.GroupList[0], ShouldResemble, ColumnName{"f"})
				So(comp.GroupList[1], ShouldResemble, ColumnName{"g"})
				So(comp.Having, ShouldResemble, ColumnName{"h"})
			})
		})
	})
}

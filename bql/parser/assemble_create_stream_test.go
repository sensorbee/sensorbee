package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleCreateStream(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}
		Convey("When the stack contains the correct CREATE STREAM items", func() {
			ps.PushComponent(2, 4, Relation{"x"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.PushComponent(7, 8, ColumnName{"b"})
			ps.AssembleProjections(6, 8)
			ps.AssembleEmitProjections()
			ps.PushComponent(12, 13, Relation{"c"})
			ps.PushComponent(13, 14, Relation{"d"})
			ps.PushComponent(14, 15, Raw{"2"})
			ps.PushComponent(15, 20, Seconds)
			ps.AssembleRange()
			ps.AssembleWindowedFrom(12, 20)
			ps.PushComponent(20, 21, ColumnName{"e"})
			ps.AssembleFilter(20, 21)
			ps.PushComponent(21, 22, ColumnName{"f"})
			ps.PushComponent(22, 23, ColumnName{"g"})
			ps.AssembleGrouping(21, 23)
			ps.PushComponent(23, 24, ColumnName{"h"})
			ps.AssembleHaving(23, 24)
			ps.AssembleCreateStream()

			Convey("Then AssembleCreateStream transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a CreateStreamStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 24)
					So(top.comp, ShouldHaveSameTypeAs, CreateStreamStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(CreateStreamStmt)
						So(comp.name, ShouldEqual, "x")
						So(comp.emitterType, ShouldEqual, Istream)
						So(len(comp.projections), ShouldEqual, 2)
						So(comp.projections[0], ShouldResemble, ColumnName{"a"})
						So(comp.projections[1], ShouldResemble, ColumnName{"b"})
						So(len(comp.relations), ShouldEqual, 2)
						So(comp.relations[0].name, ShouldEqual, "c")
						So(comp.relations[1].name, ShouldEqual, "d")
						So(comp.expr, ShouldEqual, "2")
						So(comp.unit, ShouldEqual, Seconds)
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
			Convey("Then AssembleCreateStream panics", func() {
				So(ps.AssembleCreateStream, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Relation{"x"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.PushComponent(7, 8, ColumnName{"b"})
			ps.AssembleProjections(6, 8)
			ps.AssembleEmitProjections()
			ps.PushComponent(12, 13, Relation{"c"})
			ps.PushComponent(13, 14, Relation{"d"})
			ps.PushComponent(14, 15, Raw{"2"})
			ps.PushComponent(15, 20, Seconds)
			ps.AssembleRange()
			ps.AssembleWindowedFrom(12, 20)
			ps.PushComponent(20, 21, ColumnName{"e"})
			ps.AssembleFilter(20, 21)
			ps.PushComponent(21, 22, ColumnName{"f"})
			ps.PushComponent(22, 23, ColumnName{"g"})
			ps.AssembleGrouping(21, 23)
			ps.PushComponent(23, 24, ColumnName{"h"})
			ps.AssembleFilter(23, 24) // must be HAVING in correct stmt

			Convey("Then AssembleCreateStream panics", func() {
				So(ps.AssembleCreateStream, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &Bql{}

		Convey("When doing a full SELECT", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM(a, b) FROM c, d [RANGE 2 SECONDS] WHERE e GROUP BY f, g HAVING h"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamStmt{})
				comp := top.(CreateStreamStmt)

				So(comp.name, ShouldEqual, "x")
				So(comp.emitterType, ShouldEqual, Istream)
				So(len(comp.projections), ShouldEqual, 2)
				So(comp.projections[0], ShouldResemble, ColumnName{"a"})
				So(comp.projections[1], ShouldResemble, ColumnName{"b"})
				So(len(comp.relations), ShouldEqual, 2)
				So(comp.relations[0].name, ShouldEqual, "c")
				So(comp.relations[1].name, ShouldEqual, "d")
				So(comp.expr, ShouldEqual, "2")
				So(comp.unit, ShouldEqual, Seconds)
				So(comp.filter, ShouldResemble, ColumnName{"e"})
				So(len(comp.groupList), ShouldEqual, 2)
				So(comp.groupList[0], ShouldResemble, ColumnName{"f"})
				So(comp.groupList[1], ShouldResemble, ColumnName{"g"})
				So(comp.having, ShouldResemble, ColumnName{"h"})
			})
		})
	})
}

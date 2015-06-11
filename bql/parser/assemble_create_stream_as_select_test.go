package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleCreateStreamAsSelect(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct CREATE STREAM items", func() {
			ps.PushComponent(2, 4, Relation{"x"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.PushComponent(8, 9, Identifier("y"))
			ps.AssembleAlias()
			ps.AssembleProjections(6, 9)
			ps.AssembleEmitProjections()
			ps.PushComponent(12, 13, Relation{"c"})
			ps.EnsureAliasRelation()
			ps.PushComponent(13, 14, Relation{"d"})
			ps.EnsureAliasRelation()
			ps.PushComponent(14, 15, NumericLiteral{2})
			ps.PushComponent(15, 20, Seconds)
			ps.AssembleRange()
			ps.AssembleWindowedFrom(12, 20)
			ps.PushComponent(20, 21, RowValue{"", "e"})
			ps.AssembleFilter(20, 21)
			ps.PushComponent(21, 22, RowValue{"", "f"})
			ps.PushComponent(22, 23, RowValue{"", "g"})
			ps.AssembleGrouping(21, 23)
			ps.PushComponent(23, 24, RowValue{"", "h"})
			ps.AssembleHaving(23, 24)
			ps.AssembleCreateStreamAsSelect()

			Convey("Then AssembleCreateStreamAsSelect transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a CreateStreamAsSelectStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 24)
					So(top.comp, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(CreateStreamAsSelectStmt)
						So(comp.Name, ShouldEqual, "x")
						So(comp.EmitterType, ShouldEqual, Istream)
						So(len(comp.Projections), ShouldEqual, 2)
						So(comp.Projections[0], ShouldResemble, RowValue{"", "a"})
						So(comp.Projections[1], ShouldResemble, AliasAST{RowValue{"", "b"}, "y"})
						So(len(comp.Relations), ShouldEqual, 2)
						So(comp.Relations[0].Name, ShouldEqual, "c")
						So(comp.Relations[1].Name, ShouldEqual, "d")
						So(comp.Value, ShouldEqual, 2)
						So(comp.Unit, ShouldEqual, Seconds)
						So(comp.Filter, ShouldResemble, RowValue{"", "e"})
						So(len(comp.GroupList), ShouldEqual, 2)
						So(comp.GroupList[0], ShouldResemble, RowValue{"", "f"})
						So(comp.GroupList[1], ShouldResemble, RowValue{"", "g"})
						So(comp.Having, ShouldResemble, RowValue{"", "h"})
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.AssembleProjections(6, 7)
			Convey("Then AssembleCreateStreamAsSelect panics", func() {
				So(ps.AssembleCreateStreamAsSelect, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Relation{"x"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.PushComponent(8, 9, Identifier("y"))
			ps.AssembleAlias()
			ps.AssembleProjections(6, 9)
			ps.AssembleEmitProjections()
			ps.PushComponent(12, 13, Relation{"c"})
			ps.EnsureAliasRelation()
			ps.PushComponent(13, 14, Relation{"d"})
			ps.EnsureAliasRelation()
			ps.PushComponent(14, 15, NumericLiteral{2})
			ps.PushComponent(15, 20, Seconds)
			ps.AssembleRange()
			ps.AssembleWindowedFrom(12, 20)
			ps.PushComponent(20, 21, RowValue{"", "e"})
			ps.AssembleFilter(20, 21)
			ps.PushComponent(21, 22, RowValue{"", "f"})
			ps.PushComponent(22, 23, RowValue{"", "g"})
			ps.AssembleGrouping(21, 23)
			ps.PushComponent(23, 24, RowValue{"", "h"})
			ps.AssembleFilter(23, 24) // must be HAVING in correct stmt

			Convey("Then AssembleCreateStreamAsSelect panics", func() {
				So(ps.AssembleCreateStreamAsSelect, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full SELECT", func() {
			p.Buffer = "CREATE STREAM x_2 AS SELECT ISTREAM '日本語', b AS y FROM c, d [RANGE 2 SECONDS] WHERE e GROUP BY f, g HAVING h"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x_2")
				So(comp.EmitterType, ShouldEqual, Istream)
				So(len(comp.Projections), ShouldEqual, 2)
				So(comp.Projections[0], ShouldResemble, StringLiteral{"日本語"})
				So(comp.Projections[1], ShouldResemble, AliasAST{RowValue{"", "b"}, "y"})
				So(len(comp.Relations), ShouldEqual, 2)
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(comp.Relations[1].Name, ShouldEqual, "d")
				So(comp.Value, ShouldEqual, 2)
				So(comp.Unit, ShouldEqual, Seconds)
				So(comp.Filter, ShouldResemble, RowValue{"", "e"})
				So(len(comp.GroupList), ShouldEqual, 2)
				So(comp.GroupList[0], ShouldResemble, RowValue{"", "f"})
				So(comp.GroupList[1], ShouldResemble, RowValue{"", "g"})
				So(comp.Having, ShouldResemble, RowValue{"", "h"})
			})
		})
	})
}

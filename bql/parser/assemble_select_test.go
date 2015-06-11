package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleSelect(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct SELECT items", func() {
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleProjections(6, 8)
			ps.PushComponent(16, 18, Relation{"c"})
			ps.EnsureAliasRelation()
			ps.PushComponent(18, 20, Relation{"d"})
			ps.PushComponent(20, 22, Identifier("x"))
			ps.AssembleAliasRelation()
			ps.EnsureAliasRelation()
			ps.AssembleFrom(16, 22)
			ps.PushComponent(22, 24, RowValue{"", "e"})
			ps.AssembleFilter(22, 24)
			ps.PushComponent(24, 26, RowValue{"", "f"})
			ps.PushComponent(26, 28, RowValue{"", "g"})
			ps.AssembleGrouping(24, 28)
			ps.PushComponent(28, 30, RowValue{"", "h"})
			ps.AssembleHaving(28, 30)
			ps.AssembleSelect()

			Convey("Then AssembleSelect transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a SelectStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 30)
					So(top.comp, ShouldHaveSameTypeAs, SelectStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SelectStmt)
						So(len(comp.Projections), ShouldEqual, 2)
						So(comp.Projections[0], ShouldResemble, RowValue{"", "a"})
						So(comp.Projections[1], ShouldResemble, RowValue{"", "b"})
						So(len(comp.Relations), ShouldEqual, 2)
						So(comp.Relations[0].Name, ShouldEqual, "c")
						So(comp.Relations[0].Alias, ShouldEqual, "")
						So(comp.Relations[1].Name, ShouldEqual, "d")
						So(comp.Relations[1].Alias, ShouldEqual, "x")
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
			Convey("Then AssembleSelect panics", func() {
				So(ps.AssembleSelect, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleProjections(6, 8)
			ps.PushComponent(16, 18, Relation{"c"})
			ps.EnsureAliasRelation()
			ps.PushComponent(18, 20, Relation{"d"})
			ps.PushComponent(20, 22, Identifier("x"))
			ps.AssembleAliasRelation()
			ps.EnsureAliasRelation()
			ps.AssembleFrom(16, 22)
			ps.PushComponent(22, 24, RowValue{"", "e"})
			ps.AssembleFilter(22, 24)
			ps.PushComponent(24, 26, RowValue{"", "f"})
			ps.PushComponent(26, 28, RowValue{"", "g"})
			ps.AssembleGrouping(24, 28)
			ps.PushComponent(28, 30, RowValue{"", "h"})
			ps.AssembleFilter(28, 30) // must be HAVING in correct stmt
			Convey("Then AssembleSelect panics", func() {
				So(ps.AssembleSelect, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full SELECT", func() {
			p.Buffer = "SELECT '日本語', b FROM c, d AS x WHERE e GROUP BY f, g HAVING h"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				comp := top.(SelectStmt)

				So(len(comp.Projections), ShouldEqual, 2)
				So(comp.Projections[0], ShouldResemble, StringLiteral{"日本語"})
				So(comp.Projections[1], ShouldResemble, RowValue{"", "b"})
				So(len(comp.Relations), ShouldEqual, 2)
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(comp.Relations[0].Alias, ShouldEqual, "")
				So(comp.Relations[1].Name, ShouldEqual, "d")
				So(comp.Relations[1].Alias, ShouldEqual, "x")
				So(comp.Filter, ShouldResemble, RowValue{"", "e"})
				So(len(comp.GroupList), ShouldEqual, 2)
				So(comp.GroupList[0], ShouldResemble, RowValue{"", "f"})
				So(comp.GroupList[1], ShouldResemble, RowValue{"", "g"})
				So(comp.Having, ShouldResemble, RowValue{"", "h"})
			})
		})
	})
}

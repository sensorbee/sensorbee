package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleCreateStreamAsSelect(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct CREATE STREAM items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("x"))
			ps.PushComponent(4, 6, Istream)
			ps.AssembleEmitter()
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.PushComponent(8, 9, Identifier("y"))
			ps.AssembleAlias()
			ps.AssembleProjections(6, 9)
			ps.PushComponent(12, 13, Stream{ActualStream, "c", nil})
			ps.PushComponent(13, 14, IntervalAST{NumericLiteral{3}, Tuples})
			ps.AssembleStreamWindow()
			ps.EnsureAliasedStreamWindow()
			ps.PushComponent(14, 15, Stream{ActualStream, "d", nil})
			ps.PushComponent(16, 17, NumericLiteral{2})
			ps.PushComponent(17, 18, Seconds)
			ps.AssembleInterval()
			ps.AssembleStreamWindow()
			ps.PushComponent(18, 19, Identifier("x"))
			ps.AssembleAliasedStreamWindow()
			ps.EnsureAliasedStreamWindow()
			ps.AssembleWindowedFrom(12, 20)
			ps.PushComponent(20, 21, RowValue{"", "e"})
			ps.AssembleFilter(20, 21)
			ps.PushComponent(21, 22, RowValue{"", "f"})
			ps.PushComponent(22, 23, RowValue{"", "g"})
			ps.AssembleGrouping(21, 23)
			ps.PushComponent(23, 24, RowValue{"", "h"})
			ps.AssembleHaving(23, 24)
			ps.AssembleSelect()
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
						cssComp := top.comp.(CreateStreamAsSelectStmt)
						So(cssComp.Name, ShouldEqual, "x")
						comp := cssComp.Select
						So(comp.EmitterType, ShouldEqual, Istream)
						So(len(comp.Projections), ShouldEqual, 2)
						So(comp.Projections[0], ShouldResemble, RowValue{"", "a"})
						So(comp.Projections[1], ShouldResemble, AliasAST{RowValue{"", "b"}, "y"})
						So(len(comp.Relations), ShouldEqual, 2)
						So(comp.Relations[0].Name, ShouldEqual, "c")
						So(comp.Relations[0].Value, ShouldEqual, 3)
						So(comp.Relations[0].Unit, ShouldEqual, Tuples)
						So(comp.Relations[0].Alias, ShouldEqual, "")
						So(comp.Relations[1].Name, ShouldEqual, "d")
						So(comp.Relations[1].Value, ShouldEqual, 2)
						So(comp.Relations[1].Unit, ShouldEqual, Seconds)
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
			Convey("Then AssembleCreateStreamAsSelect panics", func() {
				So(ps.AssembleCreateStreamAsSelect, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, StreamIdentifier("x"))
			ps.PushComponent(4, 6, Istream) // must be SELECT in correct stmt

			Convey("Then AssembleCreateStreamAsSelect panics", func() {
				So(ps.AssembleCreateStreamAsSelect, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full SELECT", func() {
			p.Buffer = "CREATE STREAM x_2 AS SELECT ISTREAM '日本語', b AS y FROM c [RANGE 3 TUPLES], d [RANGE 2 SECONDS] AS x WHERE e GROUP BY f, g HAVING h"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				cssComp := top.(CreateStreamAsSelectStmt)

				So(cssComp.Name, ShouldEqual, "x_2")
				comp := cssComp.Select
				So(comp.EmitterType, ShouldEqual, Istream)
				So(len(comp.Projections), ShouldEqual, 2)
				So(comp.Projections[0], ShouldResemble, StringLiteral{"日本語"})
				So(comp.Projections[1], ShouldResemble, AliasAST{RowValue{"", "b"}, "y"})
				So(len(comp.Relations), ShouldEqual, 2)
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(comp.Relations[0].Value, ShouldEqual, 3)
				So(comp.Relations[0].Unit, ShouldEqual, Tuples)
				So(comp.Relations[0].Alias, ShouldEqual, "")
				So(comp.Relations[1].Name, ShouldEqual, "d")
				So(comp.Relations[1].Value, ShouldEqual, 2)
				So(comp.Relations[1].Unit, ShouldEqual, Seconds)
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

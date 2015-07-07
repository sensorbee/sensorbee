package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleSelect(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct SELECT items", func() {
			ps.PushComponent(4, 6, Istream)
			ps.AssembleEmitter(6, 6)
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleProjections(6, 8)
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
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 30)
					So(top.comp, ShouldHaveSameTypeAs, SelectStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(SelectStmt)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(len(comp.Projections), ShouldEqual, 2)
						So(comp.Projections[0], ShouldResemble, RowValue{"", "a"})
						So(comp.Projections[1], ShouldResemble, RowValue{"", "b"})
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
			Convey("Then AssembleSelect panics", func() {
				So(ps.AssembleSelect, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(4, 6, Istream)
			ps.AssembleEmitter(6, 6)
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleProjections(6, 8)
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
			p.Buffer = "SELECT ISTREAM '日本語', b FROM c [RANGE 3 TUPLES], d('state', 7) [RANGE 2 SECONDS] AS x WHERE e GROUP BY f, g HAVING h"
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

				So(comp.EmitterType, ShouldEqual, Istream)
				So(len(comp.Projections), ShouldEqual, 2)
				So(comp.Projections[0], ShouldResemble, StringLiteral{"日本語"})
				So(comp.Projections[1], ShouldResemble, RowValue{"", "b"})
				So(len(comp.Relations), ShouldEqual, 2)
				So(comp.Relations[0].Type, ShouldEqual, ActualStream)
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(len(comp.Relations[0].Params), ShouldEqual, 0)
				So(comp.Relations[0].Value, ShouldEqual, 3)
				So(comp.Relations[0].Unit, ShouldEqual, Tuples)
				So(comp.Relations[0].Alias, ShouldEqual, "")
				So(comp.Relations[1].Type, ShouldEqual, UDSFStream)
				So(comp.Relations[1].Name, ShouldEqual, "d")
				So(len(comp.Relations[1].Params), ShouldEqual, 2)
				So(comp.Relations[1].Params[0], ShouldResemble, StringLiteral{"state"})
				So(comp.Relations[1].Params[1], ShouldResemble, NumericLiteral{7})
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

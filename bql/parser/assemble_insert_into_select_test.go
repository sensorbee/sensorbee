package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleInsertIntoSelect(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct SELECT items with a Interval specification", func() {
			ps.PushComponent(4, 5, StreamIdentifier("x"))
			ps.PushComponent(5, 6, Istream)
			ps.AssembleEmitter()
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleProjections(6, 8)
			ps.PushComponent(12, 13, StreamWindowAST{Stream{ActualStream, "c", nil}, IntervalAST{NumericLiteral{3}, Tuples}})
			ps.EnsureAliasedStreamWindow()
			ps.PushComponent(13, 14, StreamWindowAST{Stream{ActualStream, "d", nil}, IntervalAST{NumericLiteral{2}, Seconds}})
			ps.EnsureAliasedStreamWindow()
			ps.AssembleWindowedFrom(12, 14)
			ps.PushComponent(14, 15, RowValue{"", "e"})
			ps.AssembleFilter(14, 15)
			ps.PushComponent(15, 16, RowValue{"", "f"})
			ps.PushComponent(16, 17, RowValue{"", "g"})
			ps.AssembleGrouping(15, 17)
			ps.PushComponent(17, 18, RowValue{"", "h"})
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
						So(comp.Relations[1].Alias, ShouldEqual, "")
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
			Convey("Then AssembleInsertIntoSelect panics", func() {
				So(ps.AssembleInsertIntoSelect, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(4, 5, StreamIdentifier("x"))
			ps.PushComponent(5, 6, Istream)
			ps.AssembleEmitter()
			ps.PushComponent(6, 7, RowValue{"", "a"})
			ps.PushComponent(7, 8, RowValue{"", "b"})
			ps.AssembleProjections(6, 8)
			ps.PushComponent(12, 13, StreamWindowAST{Stream{ActualStream, "c", nil}, IntervalAST{NumericLiteral{3}, Tuples}})
			ps.EnsureAliasedStreamWindow()
			ps.PushComponent(13, 14, StreamWindowAST{Stream{ActualStream, "d", nil}, IntervalAST{NumericLiteral{2}, Seconds}})
			ps.EnsureAliasedStreamWindow()
			ps.AssembleWindowedFrom(12, 14)
			ps.PushComponent(14, 15, RowValue{"", "e"})
			ps.AssembleFilter(14, 15)
			ps.PushComponent(15, 16, RowValue{"", "f"})
			ps.PushComponent(16, 17, RowValue{"", "g"})
			ps.AssembleGrouping(15, 17)
			ps.PushComponent(17, 18, RowValue{"", "h"})
			ps.AssembleFilter(17, 18) // must be HAVING in correct stmt
			Convey("Then AssembleInsertIntoSelect panics", func() {
				So(ps.AssembleInsertIntoSelect, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full INSERT INTO SELECT", func() {
			// the statement below does not make sense, it is just used to test
			// whether we accept optional RANGE clauses correctly
			p.Buffer = "INSERT INTO x SELECT ISTREAM '日本語', b FROM c [RANGE 3 TUPLES], d [RANGE 1 SECONDS] WHERE e GROUP BY f, g HAVING h"
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
				So(comp.EmitterType, ShouldEqual, Istream)
				So(len(comp.Projections), ShouldEqual, 2)
				So(comp.Projections[0], ShouldResemble, StringLiteral{"日本語"})
				So(comp.Projections[1], ShouldResemble, RowValue{"", "b"})
				So(len(comp.Relations), ShouldEqual, 2)
				So(comp.Relations[0].Name, ShouldEqual, "c")
				So(comp.Relations[0].Value, ShouldEqual, 3)
				So(comp.Relations[0].Unit, ShouldEqual, Tuples)
				So(comp.Relations[0].Alias, ShouldEqual, "")
				So(comp.Relations[1].Name, ShouldEqual, "d")
				So(comp.Relations[1].Value, ShouldEqual, 1)
				So(comp.Relations[1].Unit, ShouldEqual, Seconds)
				So(comp.Relations[1].Alias, ShouldEqual, "")
				So(comp.Filter, ShouldResemble, RowValue{"", "e"})
				So(len(comp.GroupList), ShouldEqual, 2)
				So(comp.GroupList[0], ShouldResemble, RowValue{"", "f"})
				So(comp.GroupList[1], ShouldResemble, RowValue{"", "g"})
				So(comp.Having, ShouldResemble, RowValue{"", "h"})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

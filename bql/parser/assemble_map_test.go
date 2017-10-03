package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleMap(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains only KeyValuePairs in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, StringLiteral{"foo"})
			ps.PushComponent(7, 8, RowValue{"", "a"})
			ps.AssembleKeyValuePair()
			ps.PushComponent(8, 9, StringLiteral{"bar"})
			ps.PushComponent(9, 10, RowValue{"", "b"})
			ps.AssembleKeyValuePair()
			ps.AssembleMap(6, 10)

			Convey("Then AssembleMap transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a MapAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, MapAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(MapAST)
						So(len(comp.Entries), ShouldEqual, 2)
						So(comp.Entries[0], ShouldResemble, KeyValuePairAST{"foo", RowValue{"", "a"}})
						So(comp.Entries[1], ShouldResemble, KeyValuePairAST{"bar", RowValue{"", "b"}})
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleMap(6, 8)

			Convey("Then AssembleMap pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a MapAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, MapAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(MapAST)
						So(len(comp.Entries), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleMap(6, 6)

			Convey("Then AssembleMap pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a MapAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, MapAST{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(MapAST)
						So(len(comp.Entries), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting a function of multiple columns", func() {
			p.Buffer = `SELECT ISTREAM {"foo": a,"bar":b}`
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.Projections), ShouldEqual, 1)
				So(s.Projections[0], ShouldHaveSameTypeAs, MapAST{})
				m := s.Projections[0].(MapAST)
				So(len(m.Entries), ShouldEqual, 2)
				So(m.Entries[0], ShouldResemble, KeyValuePairAST{"foo", RowValue{"", "a"}})
				So(m.Entries[1], ShouldResemble, KeyValuePairAST{"bar", RowValue{"", "b"}})
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleHaving(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains one item in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, RowValue{"a"})
			ps.AssembleHaving(6, 7)

			Convey("Then AssembleHaving replaces this with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a HavingAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 7)
					So(top.comp, ShouldHaveSameTypeAs, HavingAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(HavingAST)
						So(comp.Having, ShouldResemble, RowValue{"a"})
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleHaving(6, 6)

			Convey("Then AssembleHaving pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a HavingAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, HavingAST{})

					Convey("And it contains a nil pointer", func() {
						comp := top.comp.(HavingAST)
						So(comp.Having, ShouldBeNil)
					})
				})
			})
		})

		Convey("When the stack contains one item not in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, RowValue{"a"})
			f := func() {
				ps.AssembleHaving(5, 6)
			}
			Convey("Then AssembleHaving panics", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting without a HAVING", func() {
			p.Buffer = "SELECT a, b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(s.Having, ShouldBeNil)
			})
		})

		Convey("When selecting with a HAVING", func() {
			p.Buffer = "SELECT a, b HAVING c"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(s.Having, ShouldNotBeNil)
				So(s.Having, ShouldResemble, RowValue{"c"})
			})
		})
	})
}

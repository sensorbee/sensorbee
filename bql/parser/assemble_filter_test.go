package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleFilter(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains one item in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, ColumnName{"a"})
			ps.AssembleFilter(6, 7)

			Convey("Then AssembleFilter replaces this with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Filter", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 7)
					So(top.comp, ShouldHaveSameTypeAs, Filter{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(Filter)
						So(comp.filter, ShouldResemble, ColumnName{"a"})
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleFilter(6, 6)

			Convey("Then AssembleFilter pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Filter", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, Filter{})

					Convey("And it contains a nil pointer", func() {
						comp := top.comp.(Filter)
						So(comp.filter, ShouldBeNil)
					})
				})
			})
		})

		Convey("When the stack contains one item not in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, ColumnName{"a"})
			f := func() {
				ps.AssembleFilter(5, 6)
			}
			Convey("Then AssembleFilter panics", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When selecting without a WHERE", func() {
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
				So(s.filter, ShouldBeNil)
			})
		})

		Convey("When selecting with a WHERE", func() {
			p.Buffer = "SELECT a, b WHERE c"
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
				So(s.filter, ShouldNotBeNil)
				So(s.filter, ShouldResemble, ColumnName{"c"})
			})
		})
	})
}

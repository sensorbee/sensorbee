package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleFrom(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}

		Convey("When the stack contains only Relations in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, Relation{"a"})
			ps.PushComponent(7, 8, Relation{"b"})
			ps.AssembleFrom(6, 8)

			Convey("Then AssembleFrom transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a From", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, From{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(From)
						So(len(comp.relations), ShouldEqual, 2)
						So(comp.relations[0].name, ShouldEqual, "a")
						So(comp.relations[1].name, ShouldEqual, "b")
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleFrom(6, 8)

			Convey("Then AssembleFrom pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a From", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, From{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(From)
						So(len(comp.relations), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleFrom(6, 6)

			Convey("Then AssembleFrom pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a From", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, From{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(From)
						So(len(comp.relations), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the stack contains non-Relations in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			f := func() {
				ps.AssembleFrom(0, 8)
			}

			Convey("Then AssembleFrom panics", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &Bql{}

		Convey("When selecting without a FROM", func() {
			p.Buffer = "SELECT a, b"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.relations), ShouldEqual, 0)
			})
		})

		Convey("When selecting with a FROM", func() {
			p.Buffer = "SELECT a, b FROM c, d"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, SelectStmt{})
				s := top.(SelectStmt)
				So(len(s.relations), ShouldEqual, 2)
				So(s.relations[0].name, ShouldEqual, "c")
				So(s.relations[1].name, ShouldEqual, "d")
			})
		})
	})
}

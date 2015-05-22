package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleWindowedFrom(t *testing.T) {
	Convey("Given a ParseStack", t, func() {
		ps := ParseStack{}

		Convey("When the stack contains only Relations and a Range in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, Relation{"a"})
			ps.PushComponent(7, 8, Relation{"b"})
			ps.PushComponent(8, 10, Range{Raw{"2"}, Seconds})
			ps.AssembleWindowedFrom(6, 10)

			Convey("Then AssembleWindowedFrom transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedFrom", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, WindowedFrom{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(WindowedFrom)
						So(len(comp.relations), ShouldEqual, 2)
						So(comp.relations[0].name, ShouldEqual, "a")
						So(comp.relations[1].name, ShouldEqual, "b")
						So(comp.expr, ShouldEqual, "2")
						So(comp.unit, ShouldEqual, Seconds)
					})
				})
			})
		})

		Convey("When the stack contains no elements in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleWindowedFrom(6, 8)

			Convey("Then AssembleWindowedFrom pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedFrom", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, WindowedFrom{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(WindowedFrom)
						So(len(comp.relations), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the given range is empty", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.AssembleWindowedFrom(6, 6)

			Convey("Then AssembleWindowedFrom pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a WindowedFrom", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, WindowedFrom{})

					Convey("And it contains an empty list", func() {
						comp := top.comp.(WindowedFrom)
						So(len(comp.relations), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When the stack contains non-Relations in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 8, Range{Raw{"2"}, Seconds})
			f := func() {
				ps.AssembleWindowedFrom(0, 8)
			}

			Convey("Then AssembleWindowedFrom panics", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When the stack contains a non-Range on top", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, Relation{"a"})
			ps.PushComponent(7, 8, Relation{"b"})
			f := func() {
				ps.AssembleWindowedFrom(0, 8)
			}

			Convey("Then AssembleWindowedFrom panics", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &Bql{}

		Convey("When selecting without a FROM", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM (a, b)"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.ParseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamStmt{})
				s := top.(CreateStreamStmt)
				So(s.WindowedFrom, ShouldResemble, WindowedFrom{})
			})
		})

		Convey("When selecting with a FROM", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM(a, b) FROM c, d [RANGE 2 SECONDS]"
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
				So(len(comp.relations), ShouldEqual, 2)
				So(comp.relations[0].name, ShouldEqual, "c")
				So(comp.relations[1].name, ShouldEqual, "d")
				So(comp.expr, ShouldEqual, "2")
				So(comp.unit, ShouldEqual, Seconds)
			})
		})
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleDropSource(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct DROP SOURCE items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleDropSource()

			Convey("Then AssembleDropSource transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a DropSourceStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, DropSourceStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(DropSourceStmt)
						So(comp.Source, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleDropSource panics", func() {
				So(ps.AssembleDropSource, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full DROP SOURCE", func() {
			p.Buffer = "DROP SOURCE a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, DropSourceStmt{})
				comp := top.(DropSourceStmt)

				So(comp.Source, ShouldEqual, "a_1")

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

func TestAssembleDropStream(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct DROP STREAM items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleDropStream()

			Convey("Then AssembleDropStream transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a DropStreamStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, DropStreamStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(DropStreamStmt)
						So(comp.Stream, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleDropStream panics", func() {
				So(ps.AssembleDropStream, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full DROP STREAM", func() {
			p.Buffer = "DROP STREAM a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, DropStreamStmt{})
				comp := top.(DropStreamStmt)

				So(comp.Stream, ShouldEqual, "a_1")

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

func TestAssembleDropSink(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct DROP SINK items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleDropSink()

			Convey("Then AssembleDropSink transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a DropSinkStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, DropSinkStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(DropSinkStmt)
						So(comp.Sink, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleDropSink panics", func() {
				So(ps.AssembleDropSink, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full DROP SINK", func() {
			p.Buffer = "DROP SINK a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, DropSinkStmt{})
				comp := top.(DropSinkStmt)

				So(comp.Sink, ShouldEqual, "a_1")

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

func TestAssembleDropState(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains the correct DROP STATE items", func() {
			ps.PushComponent(2, 4, StreamIdentifier("a"))
			ps.AssembleDropState()

			Convey("Then AssembleDropState transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And that item is a DropStateStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 2)
					So(top.end, ShouldEqual, 4)
					So(top.comp, ShouldHaveSameTypeAs, DropStateStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(DropStateStmt)
						So(comp.State, ShouldEqual, "a")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"a"}) // must be StreamIdentifier

			Convey("Then AssembleDropState panics", func() {
				So(ps.AssembleDropState, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing a full DROP STATE", func() {
			p.Buffer = "DROP STATE a_1"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldEqual, nil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, DropStateStmt{})
				comp := top.(DropStateStmt)

				So(comp.State, ShouldEqual, "a_1")

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

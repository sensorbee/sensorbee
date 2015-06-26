package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEnsureKeywordPresent(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains one Yes item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 8, Yes)
			ps.EnsureKeywordPresent(6, 8)

			Convey("Then EnsureKeywordPresent does not touch it", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And the top item is still the same", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Yes)

					So(top.comp.(BinaryKeyword), ShouldEqual, Yes)
				})
			})
		})

		Convey("When the stack contains one No item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 8, No)
			ps.EnsureKeywordPresent(6, 8)

			Convey("Then EnsureKeywordPresent does not touch it", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And the top item is still the same", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Yes)

					So(top.comp.(BinaryKeyword), ShouldEqual, No)
				})
			})
		})

		Convey("When the stack contains no item in the given range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.EnsureKeywordPresent(6, 8)

			Convey("Then EnsureKeywordPresent pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And the top item is an UnspecifiedKeyword", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Yes)

					So(top.comp.(BinaryKeyword), ShouldEqual, UnspecifiedKeyword)
				})
			})
		})

		Convey("When the stack is empty", func() {
			ps.EnsureKeywordPresent(6, 8)

			Convey("Then EnsureKeywordPresent pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 1)

				Convey("And the top item is an UnspecifiedKeyword", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Yes)

					So(top.comp.(BinaryKeyword), ShouldEqual, UnspecifiedKeyword)
				})
			})
		})

		Convey("When the call has an empty range", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.EnsureKeywordPresent(6, 6)

			Convey("Then EnsureKeywordPresent pushes one item onto the stack", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And the top item is an UnspecifiedKeyword", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, Yes)

					So(top.comp.(BinaryKeyword), ShouldEqual, UnspecifiedKeyword)
				})
			})
		})
	})
}

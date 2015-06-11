package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleAliasWindowedRelation(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, WindowedRelationAST{Relation{"a"},
				RangeAST{NumericLiteral{2}, Seconds}})
			ps.PushComponent(7, 8, Identifier("out"))
			ps.AssembleAliasWindowedRelation()

			Convey("Then AssembleAliasWindowedRelation replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a AliasWindowedRelationAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, AliasWindowedRelationAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(AliasWindowedRelationAST)
						So(comp.WindowedRelationAST, ShouldResemble,
							WindowedRelationAST{Relation{"a"}, RangeAST{NumericLiteral{2}, Seconds}})
						So(comp.Alias, ShouldEqual, "out")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleAliasWindowedRelation panics", func() {
				So(ps.AssembleAliasWindowedRelation, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleAliasWindowedRelation panics", func() {
				So(ps.AssembleAliasWindowedRelation, ShouldPanic)
			})
		})
	})
}

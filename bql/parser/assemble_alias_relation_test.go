package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleAliasRelation(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains two correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, Relation{"foo"})
			ps.PushComponent(7, 8, Identifier("bar"))
			ps.AssembleAliasRelation()

			Convey("Then AssembleAliasRelation replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a AliasRelationAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, AliasRelationAST{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(AliasRelationAST)
						So(comp.Name, ShouldResemble, "foo")
						So(comp.Alias, ShouldEqual, "bar")
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleAliasRelation panics", func() {
				So(ps.AssembleAliasRelation, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleAliasRelation panics", func() {
				So(ps.AssembleAliasRelation, ShouldPanic)
			})
		})
	})
}

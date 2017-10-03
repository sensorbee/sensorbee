package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleEval(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains only one expression", func() {
			ps.PushComponent(2, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, RowValue{"", "a"})
			ps.AssembleEval(6, 6)

			Convey("Then AssembleEval transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is an EvalStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, EvalStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EvalStmt)
						So(comp.Expr, ShouldResemble, RowValue{"", "a"})
						So(comp.Input, ShouldBeNil)
					})
				})
			})
		})

		Convey("When the stack contains two expressions", func() {
			ps.PushComponent(2, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, RowValue{"", "a"})
			ps.PushComponent(6, 8, MapAST{[]KeyValuePairAST{{"a", NumericLiteral{2}}}})
			ps.AssembleEval(6, 8)

			Convey("Then AssembleEval transforms them into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is an EvalStmt", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, EvalStmt{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EvalStmt)
						So(comp.Expr, ShouldResemble, RowValue{"", "a"})
						So(*comp.Input, ShouldResemble, MapAST{[]KeyValuePairAST{{"a", NumericLiteral{2}}}})
					})
				})
			})
		})

		Convey("When the stack does not contain enough items", func() {
			f := func() { ps.AssembleEval(6, 6) }
			Convey("Then AssembleEval panics", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(2, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, RowValue{"", "a"})
			ps.PushComponent(6, 8, NumericLiteral{2})

			f := func() { ps.AssembleEval(6, 8) }
			Convey("Then AssembleEval panics", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When doing an EVAL without ON", func() {
			p.Buffer = "EVAL a"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, EvalStmt{})
				comp := top.(EvalStmt)

				So(comp.Expr, ShouldResemble, RowValue{"", "a"})
				So(comp.Input, ShouldBeNil)

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When doing an EVAL with ON", func() {
			p.Buffer = `EVAL a ON {"a":2}`
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, EvalStmt{})
				comp := top.(EvalStmt)

				So(comp.Expr, ShouldResemble, RowValue{"", "a"})
				So(*comp.Input, ShouldResemble, MapAST{[]KeyValuePairAST{{"a", NumericLiteral{2}}}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

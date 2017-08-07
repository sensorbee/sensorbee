package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleUDSFFuncApp(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}

		Convey("When the stack contains three correct items", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})
			ps.PushComponent(6, 7, FuncName("add"))
			ps.PushComponent(7, 8, ExpressionsAST{[]Expression{
				NumericLiteral{2},
				RowValue{"", "a"}}})
			ps.PushComponent(8, 8, ExpressionsAST{nil})
			ps.AssembleFuncApp()
			ps.AssembleUDSFFuncApp()

			Convey("Then AssembleFuncApp replaces them with a new item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a Stream", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 6)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, Stream{})

					Convey("And it contains the previous data", func() {
						comp := top.comp.(Stream)
						So(comp.Type, ShouldEqual, UDSFStream)
						So(comp.Name, ShouldEqual, "add")
						So(len(comp.Params), ShouldEqual, 2)
						So(comp.Params[0], ShouldResemble, NumericLiteral{2})
						So(comp.Params[1], ShouldResemble, RowValue{"", "a"})
					})
				})
			})
		})

		Convey("When the stack contains a wrong item", func() {
			ps.PushComponent(0, 6, Raw{"PRE"})

			Convey("Then AssembleUDSFFuncApp panics", func() {
				So(ps.AssembleUDSFFuncApp, ShouldPanic)
			})
		})

		Convey("When the stack is empty", func() {
			Convey("Then AssembleUDSFFuncApp panics", func() {
				So(ps.AssembleUDSFFuncApp, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When parsing a SELECT statement with a UDSF", func() {
			p.Buffer = "SELECT ISTREAM x FROM add(2, a) [RANGE 1 TUPLES]"
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
				So(len(s.WindowedFromAST.Relations), ShouldEqual, 1)
				comp := s.WindowedFromAST.Relations[0]
				So(comp.Type, ShouldEqual, UDSFStream)
				So(comp.Name, ShouldEqual, "add")
				So(len(comp.Params), ShouldEqual, 2)
				So(comp.Params[0], ShouldResemble, NumericLiteral{2})
				So(comp.Params[1], ShouldResemble, RowValue{"", "a"})

				Convey("And String() should return the original statement", func() {
					So(s.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When parsing a SELECT statement with a UDSF and ORDER BY clause", func() {
			p.Buffer = "SELECT ISTREAM x FROM add(2, a ORDER BY b) [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should fail to parse", func() {
				err := p.Parse()
				So(err, ShouldNotEqual, nil)
			})
		})
	})
}

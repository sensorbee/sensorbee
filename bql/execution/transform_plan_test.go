package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"reflect"
	"testing"
)

type analyzeTest struct {
	input         *parser.SelectStmt
	expectedError string
}

func TestRelationChecker(t *testing.T) {
	r := parser.RangeAST{parser.NumericLiteral{2}, parser.Tuples}
	singleFrom := parser.WindowedFromAST{
		[]parser.AliasedStreamWindowAST{
			{parser.StreamWindowAST{parser.Stream{"t"}, r}, ""},
		},
	}
	singleFromAlias := parser.WindowedFromAST{
		[]parser.AliasedStreamWindowAST{
			{parser.StreamWindowAST{parser.Stream{"s"}, r}, "t"},
		},
	}
	two := parser.NumericLiteral{2}
	a := parser.RowValue{"", "a"}
	b := parser.RowValue{"", "b"}
	c := parser.RowValue{"", "c"}
	t_a := parser.RowValue{"t", "a"}
	t_b := parser.RowValue{"t", "b"}
	t_c := parser.RowValue{"t", "c"}
	x_a := parser.RowValue{"x", "a"}
	x_b := parser.RowValue{"x", "b"}

	testCases := []analyzeTest{
		// SELECT a   -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{a}},
		}, "need at least one relation to select from"},
		// SELECT 2   -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{two}},
		}, "need at least one relation to select from"},
		// SELECT t.a -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{t_a}},
		}, "need at least one relation to select from"},

		////////// FROM (single input relation) //////////////

		// SELECT a        FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT 2        FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t.a      FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t.a, t.b FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a, t_b}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT 2, t.a   FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two, t_a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT a, t.b   FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, t_a}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT x.a      FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{x_a}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relation 'x' when using only 't'"},

		////////// WHERE //////////////

		// SELECT a   FROM t WHERE 2   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{two},
		}, ""},
		// SELECT 2   FROM t WHERE 2   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{two},
		}, ""},
		// SELECT t.a FROM t WHERE 2   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{two},
		}, ""},
		// SELECT a   FROM t WHERE b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{b},
		}, ""},
		// SELECT 2   FROM t WHERE b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{b},
		}, ""},
		// SELECT t.a FROM t WHERE b   -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{b},
		}, "cannot refer to relations"},
		// SELECT a   FROM t WHERE t.b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{t_b},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t WHERE t.b -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{t_b},
		}, ""},
		// SELECT t.a FROM t WHERE t.b -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{t_b},
		}, ""},
		// SELECT 2   FROM t WHERE x.b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{x_b},
		}, "cannot refer to relation 'x' when using only 't'"},

		////////// GROUP BY //////////////

		// SELECT a   FROM t GROUP BY 2        -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{two}},
		}, ""},
		// SELECT 2   FROM t GROUP BY 2        -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{two}},
		}, ""},
		// SELECT t.a FROM t GROUP BY 2        -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{two}},
		}, ""},
		// SELECT a   FROM t GROUP BY b        -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b}},
		}, ""},
		// SELECT a   FROM t GROUP BY b, c     -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b, c}},
		}, ""},
		// SELECT 2   FROM t GROUP BY b        -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b}},
		}, ""},
		// SELECT t.a FROM t GROUP BY b        -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b}},
		}, "cannot refer to relations"},
		// SELECT a   FROM t GROUP BY t.b      -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b}},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t GROUP BY t.b      -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b}},
		}, ""},
		// SELECT t.a FROM t GROUP BY t.b      -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b}},
		}, ""},
		// SELECT t.a FROM t GROUP BY t.b, t.c -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b, t_c}},
		}, ""},
		// SELECT t.a FROM t GROUP BY b, t.b   -> NG (same table with multiple aliases)
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b, t_b}},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t GROUP BY x.b      -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{x_b}},
		}, "cannot refer to relation 'x' when using only 't'"},

		////////// HAVING //////////////

		// SELECT a   FROM t HAVING 2   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{two},
		}, ""},
		// SELECT 2   FROM t HAVING 2   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{two},
		}, ""},
		// SELECT t.a FROM t HAVING 2   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{two},
		}, ""},
		// SELECT a   FROM t HAVING b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{b},
		}, ""},
		// SELECT 2   FROM t HAVING b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{b},
		}, ""},
		// SELECT t.a FROM t HAVING b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{b},
		}, ""},
		// SELECT a   FROM t HAVING t.b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{t_b},
		}, "cannot refer to input relation 't' from HAVING clause"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		selectAst := testCase.input

		Convey(fmt.Sprintf("Given the AST %+v", selectAst), t, func() {
			ast := parser.CreateStreamAsSelectStmt{
				Name: parser.StreamIdentifier("x"),
				EmitProjectionsAST: parser.EmitProjectionsAST{
					parser.Istream,
					selectAst.ProjectionsAST,
				},
				WindowedFromAST: selectAst.WindowedFromAST,
				FilterAST:       selectAst.FilterAST,
				GroupingAST:     selectAst.GroupingAST,
				HavingAST:       selectAst.HavingAST,
			}

			Convey("When we analyze it", func() {
				_, err := Analyze(ast)
				expectedError := testCase.expectedError
				if expectedError == "" {
					Convey("There is no error", func() {
						So(err, ShouldBeNil)
					})
				} else {
					Convey("There is an error", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldStartWith, expectedError)
					})
				}
			})
		})
	}

	for _, testCase := range testCases {
		testCase := testCase
		selectAst := testCase.input

		Convey(fmt.Sprintf("Given the AST %+v with an aliased relation", selectAst), t, func() {
			// we use the same test cases, but with `FROM s AS t` instead of `FROM t`
			var myFrom parser.WindowedFromAST
			if reflect.DeepEqual(selectAst.WindowedFromAST, singleFrom) {
				myFrom = singleFromAlias
			} else {
				myFrom = selectAst.WindowedFromAST
			}
			ast := parser.CreateStreamAsSelectStmt{
				Name: parser.StreamIdentifier("x"),
				EmitProjectionsAST: parser.EmitProjectionsAST{
					parser.Istream,
					selectAst.ProjectionsAST,
				},
				WindowedFromAST: myFrom,
				FilterAST:       selectAst.FilterAST,
				GroupingAST:     selectAst.GroupingAST,
				HavingAST:       selectAst.HavingAST,
			}

			Convey("When we analyze it", func() {
				_, err := Analyze(ast)
				expectedError := testCase.expectedError
				if expectedError == "" {
					Convey("There is no error", func() {
						So(err, ShouldBeNil)
					})
				} else {
					Convey("There is an error", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldStartWith, expectedError)
					})
				}
			})
		})
	}
}

func TestRelationAliasing(t *testing.T) {
	r := parser.RangeAST{parser.NumericLiteral{2}, parser.Tuples}
	two := parser.NumericLiteral{2}
	proj := parser.ProjectionsAST{[]parser.Expression{two}}

	testCases := []analyzeTest{
		// SELECT 2 FROM a              -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, ""},
				}},
		}, ""},
		// SELECT 2 FROM a AS b         -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, "b"},
				}},
		}, ""},
		// SELECT 2 FROM a AS b, a      -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, "b"},
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, ""},
				}},
		}, ""},
		// SELECT 2 FROM a AS b, c AS a -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, "b"},
					{parser.StreamWindowAST{parser.Stream{"c"}, r}, "a"},
				}},
		}, ""},
		// SELECT 2 FROM a, a           -> NG
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, ""},
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, ""},
				}},
		}, "cannot use relations"},
		// SELECT 2 FROM a, b AS a      -> NG
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{"a"}, r}, ""},
					{parser.StreamWindowAST{parser.Stream{"b"}, r}, "a"},
				}},
		}, "cannot use relations"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		selectAst := testCase.input

		Convey(fmt.Sprintf("Given the AST %+v", selectAst), t, func() {
			ast := parser.CreateStreamAsSelectStmt{
				Name: parser.StreamIdentifier("x"),
				EmitProjectionsAST: parser.EmitProjectionsAST{
					parser.Istream,
					selectAst.ProjectionsAST,
				},
				WindowedFromAST: selectAst.WindowedFromAST,
			}

			Convey("When we analyze it", func() {
				_, err := Analyze(ast)
				expectedError := testCase.expectedError
				if expectedError == "" {
					Convey("There is no error", func() {
						So(err, ShouldBeNil)
					})
				} else {
					Convey("There is an error", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldStartWith, expectedError)
					})
				}
			})
		})
	}
}

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
	r := parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}
	singleFrom := parser.WindowedFromAST{
		[]parser.AliasedStreamWindowAST{
			{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "t", nil}, r}, ""},
		},
	}
	singleFromAlias := parser.WindowedFromAST{
		[]parser.AliasedStreamWindowAST{
			{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "s", nil}, r}, "t"},
		},
	}
	two := parser.NumericLiteral{2}
	a := parser.RowValue{"", "a"}
	b := parser.RowValue{"", "b"}
	c := parser.RowValue{"", "c"}
	ts := parser.RowMeta{"", parser.TimestampMeta}
	t_a := parser.RowValue{"t", "a"}
	t_b := parser.RowValue{"t", "b"}
	t_c := parser.RowValue{"t", "c"}
	t_ts := parser.RowMeta{"t", parser.TimestampMeta}
	x_a := parser.RowValue{"x", "a"}
	x_b := parser.RowValue{"x", "b"}

	testCases := []analyzeTest{
		// SELECT a   -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{a}},
		}, "need at least one relation to select from"},
		// SELECT ts() -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{ts}},
		}, "need at least one relation to select from"},
		// SELECT 2   -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{two}},
		}, "need at least one relation to select from"},
		// SELECT t:a -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{t_a}},
		}, "need at least one relation to select from"},
		// SELECT t:ts() -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{t_ts}},
		}, "need at least one relation to select from"},

		////////// FROM (single input relation) //////////////

		// SELECT a        FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT ts()     FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{ts}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT a, ts()  FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, ts}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT 2        FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:a      FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:a, t:b FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a, t_b}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:a, t:ts() FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a, t_ts}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT 2, t:a   FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two, t_a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT a, t:b   FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, t_a}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT a, t:ts() FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, t_ts}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT x:a      FROM t -> NG
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
		// SELECT t:a FROM t WHERE 2   -> OK
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
		// SELECT t:a FROM t WHERE b   -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{b},
		}, "cannot refer to relations"},
		// SELECT a   FROM t WHERE t:b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{t_b},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t WHERE t:b -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{t_b},
		}, ""},
		// SELECT t:a FROM t WHERE t:b -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{t_b},
		}, ""},
		// SELECT 2   FROM t WHERE x:b -> NG
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
		// SELECT t:a FROM t GROUP BY 2        -> OK
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
		// SELECT t:a FROM t GROUP BY b        -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b}},
		}, "cannot refer to relations"},
		// SELECT a   FROM t GROUP BY t:b      -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b}},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t GROUP BY t:b      -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b}},
		}, ""},
		// SELECT t:a FROM t GROUP BY t:b      -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b}},
		}, ""},
		// SELECT t:a FROM t GROUP BY t:b, t:c -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{t_b, t_c}},
		}, ""},
		// SELECT t:a FROM t GROUP BY b, t:b   -> NG (same table with multiple aliases)
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b, t_b}},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t GROUP BY x:b      -> NG
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
		// SELECT t:a FROM t HAVING 2   -> OK
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
		// SELECT t:a FROM t HAVING b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{t_a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{b},
		}, ""},
		// SELECT a   FROM t HAVING t:b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{t_b},
		}, "cannot refer to input relation 't' from HAVING clause"},
	}

	emitterTestCases := []analyzeTest{
		// SELECT ISTREAM                                      a FROM t -> OK
		{&parser.SelectStmt{
			EmitterAST:      parser.EmitterAST{parser.Istream, nil},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT ISTREAM [EVERY 2 SECONDS]                    a FROM t -> OK
		{&parser.SelectStmt{
			EmitterAST: parser.EmitterAST{parser.Istream, []parser.StreamEmitIntervalAST{
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Seconds}, parser.Stream{parser.ActualStream, "*", nil}},
			}},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT ISTREAM [EVERY 2 TUPLES]                     a FROM t -> OK
		{&parser.SelectStmt{
			EmitterAST: parser.EmitterAST{parser.Istream, []parser.StreamEmitIntervalAST{
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}, parser.Stream{parser.ActualStream, "*", nil}},
			}},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT ISTREAM [EVERY 2 TUPLES IN t]                a FROM t -> OK
		{&parser.SelectStmt{
			EmitterAST: parser.EmitterAST{parser.Istream, []parser.StreamEmitIntervalAST{
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}, parser.Stream{parser.ActualStream, "t", nil}},
			}},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT ISTREAM [EVERY 2 TUPLES IN x]                a FROM t -> NG
		{&parser.SelectStmt{
			EmitterAST: parser.EmitterAST{parser.Istream, []parser.StreamEmitIntervalAST{
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}, parser.Stream{parser.ActualStream, "x", nil}},
			}},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, "the stream 'x' referenced in the ISTREAM clause is unknown"},
		// SELECT ISTREAM [EVERY 2 TUPLES IN x]             FROM t AS x -> NG
		{&parser.SelectStmt{
			EmitterAST: parser.EmitterAST{parser.Istream, []parser.StreamEmitIntervalAST{
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}, parser.Stream{parser.ActualStream, "x", nil}},
			}},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, "the stream 'x' referenced in the ISTREAM clause is unknown"},
		// SELECT ISTREAM [EVERY 2 TUPLES IN t, 3 TUPLES in x] a FROM t -> NG
		{&parser.SelectStmt{
			EmitterAST: parser.EmitterAST{parser.Istream, []parser.StreamEmitIntervalAST{
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}, parser.Stream{parser.ActualStream, "x", nil}},
			}},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, "the stream 'x' referenced in the ISTREAM clause is unknown"},
		// SELECT ISTREAM [EVERY 2 TUPLES IN t, 3 TUPLES in t] a FROM t -> NG
		{&parser.SelectStmt{
			EmitterAST: parser.EmitterAST{parser.Istream, []parser.StreamEmitIntervalAST{
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}, parser.Stream{parser.ActualStream, "t", nil}},
				{parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}, parser.Stream{parser.ActualStream, "t", nil}},
			}},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, "the stream 't' referenced in the ISTREAM clause is used more than once"},
	}

	allTestCases := append(emitterTestCases, testCases...)

	for _, testCase := range allTestCases {
		testCase := testCase
		selectAst := testCase.input

		Convey(fmt.Sprintf("Given the AST %+v", selectAst), t, func() {
			emitter := parser.EmitterAST{parser.Istream, nil}
			if selectAst.EmitterType != parser.UnspecifiedEmitter {
				emitter = selectAst.EmitterAST
			}
			ast := parser.CreateStreamAsSelectStmt{
				Name:            parser.StreamIdentifier("x"),
				EmitterAST:      emitter,
				ProjectionsAST:  selectAst.ProjectionsAST,
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
				Name:            parser.StreamIdentifier("x"),
				EmitterAST:      parser.EmitterAST{parser.Istream, nil},
				ProjectionsAST:  selectAst.ProjectionsAST,
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
	r := parser.IntervalAST{parser.NumericLiteral{2}, parser.Tuples}
	two := parser.NumericLiteral{2}
	proj := parser.ProjectionsAST{[]parser.Expression{two}}

	testCases := []analyzeTest{
		// SELECT 2 FROM a              -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, ""},
				}},
		}, ""},
		// SELECT 2 FROM a AS b         -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, "b"},
				}},
		}, ""},
		// SELECT 2 FROM a AS b, a      -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, "b"},
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, ""},
				}},
		}, ""},
		// SELECT 2 FROM a AS b, c AS a -> OK
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, "b"},
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "c", nil}, r}, "a"},
				}},
		}, ""},
		// SELECT 2 FROM a, a           -> NG
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, ""},
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, ""},
				}},
		}, "cannot use relations"},
		// SELECT 2 FROM a, b AS a      -> NG
		{&parser.SelectStmt{
			ProjectionsAST: proj,
			WindowedFromAST: parser.WindowedFromAST{
				[]parser.AliasedStreamWindowAST{
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "a", nil}, r}, ""},
					{parser.StreamWindowAST{parser.Stream{parser.ActualStream, "b", nil}, r}, "a"},
				}},
		}, "cannot use relations"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		selectAst := testCase.input

		Convey(fmt.Sprintf("Given the AST %+v", selectAst), t, func() {
			ast := parser.CreateStreamAsSelectStmt{
				Name:            parser.StreamIdentifier("x"),
				EmitterAST:      parser.EmitterAST{parser.Istream, nil},
				ProjectionsAST:  selectAst.ProjectionsAST,
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

func TestAggregateChecker(t *testing.T) {
	testCases := []struct {
		bql           string
		expectedError string
		expr          FlatExpression
		aggrs         map[string]AggFuncAppAST
	}{
		// a is no aggregate call, so the `aggrs` list is empty
		// and the selected expression is transformed normally
		{"a FROM x [RANGE 1 TUPLES]", "",
			RowValue{"x", "a"},
			nil},

		// f(a) is no aggregate call, so the `aggrs` list is empty
		// and the selected expression is transformed normally
		{"f(a) FROM x [RANGE 1 TUPLES]", "",
			FuncAppAST{"f", []FlatExpression{RowValue{"x", "a"}}},
			nil},

		// there is an aggregate call `count(a)`, so it is referenced from
		// the expression list and appears in the `aggrs` list
		{"count(a) FROM x [RANGE 1 TUPLES]", "",
			AggFuncAppRef{"#count(x:a)#"},
			map[string]AggFuncAppAST{
				"#count(x:a)#": AggFuncAppAST{"count", RowValue{"x", "a"}},
			}},

		// there is an aggregate call `count(a)`, so it is referenced from
		// the expression list and appears in the `aggrs` list
		{"a + count(a) FROM x [RANGE 1 TUPLES]", "",
			BinaryOpAST{parser.Plus,
				RowValue{"x", "a"}, AggFuncAppRef{"#count(x:a)#"}},
			map[string]AggFuncAppAST{
				"#count(x:a)#": AggFuncAppAST{"count", RowValue{"x", "a"}},
			}},

		// there is an aggregate call `udaf(a+1)`, so it is referenced from
		// the expression list and appears in the `aggrs` list
		{"a + udaf(a + 1) FROM x [RANGE 1 TUPLES]", "",
			BinaryOpAST{parser.Plus,
				RowValue{"x", "a"},
				AggFuncAppRef{"#udaf(x:a+1)#"}},
			map[string]AggFuncAppAST{
				"#udaf(x:a+1)#": AggFuncAppAST{"udaf",
					BinaryOpAST{parser.Plus, RowValue{"x", "a"}, NumericLiteral{1}}},
			}},

		// there are two aggregate calls, so both are referenced from the
		// expression list and there are two entries in the `aggrs` list
		{"udaf(a + f(1)) + g(count(a)) FROM x [RANGE 1 TUPLES]", "",
			BinaryOpAST{parser.Plus,
				AggFuncAppRef{"#udaf(x:a+f(1))#"},
				FuncAppAST{"g", []FlatExpression{
					AggFuncAppRef{"#count(x:a)#"},
				}},
			},
			map[string]AggFuncAppAST{
				"#udaf(x:a+f(1))#": AggFuncAppAST{"udaf",
					BinaryOpAST{parser.Plus,
						RowValue{"x", "a"},
						FuncAppAST{"f", []FlatExpression{NumericLiteral{1}}},
					}},
				"#count(x:a)#": AggFuncAppAST{"count", RowValue{"x", "a"}},
			}},

		// there are two aggregate calls, but they use the same value,
		// so the `aggrs` list contains only one entry
		{"count(a) + g(count(a)) FROM x [RANGE 1 TUPLES]", "",
			BinaryOpAST{parser.Plus,
				AggFuncAppRef{"#count(x:a)#"},
				FuncAppAST{"g", []FlatExpression{
					AggFuncAppRef{"#count(x:a)#"},
				}},
			},
			map[string]AggFuncAppAST{
				"#count(x:a)#": AggFuncAppAST{"count", RowValue{"x", "a"}},
			}},

		{"a + udaf(a, 1) FROM x [RANGE 1 TUPLES]",
			"aggregate functions must have exactly one parameter", nil, nil},

		{"count(udaf(a)) FROM x [RANGE 1 TUPLES]",
			"aggregate functions cannot be nested", nil, nil},

		{"a FROM x [RANGE 1 TUPLES] WHERE count(a) = 1",
			"aggregates not allowed in WHERE clause", nil, nil},

		{"a FROM x [RANGE 1 TUPLES] GROUP BY count(a)",
			"aggregates not allowed in GROUP BY clause", nil, nil},
	}

	for _, testCase := range testCases {
		testCase := testCase

		Convey(fmt.Sprintf("Given the statement", testCase.bql), t, func() {
			p := parser.NewBQLParser()
			stmt := "CREATE STREAM x AS SELECT ISTREAM " + testCase.bql
			ast_, _, err := p.ParseStmt(stmt)
			So(err, ShouldBeNil)
			So(ast_, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
			ast := ast_.(parser.CreateStreamAsSelectStmt)

			Convey("When we analyze it", func() {
				logPlan, err := Analyze(ast)
				expectedError := testCase.expectedError
				if expectedError == "" {
					Convey("There is no error", func() {
						So(err, ShouldBeNil)
						So(len(logPlan.Projections), ShouldBeGreaterThanOrEqualTo, 1)
						proj := logPlan.Projections[0]
						So(proj.expr, ShouldResemble, testCase.expr)
						So(proj.aggrInputs, ShouldResemble, testCase.aggrs)
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

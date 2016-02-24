package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"math"
	"testing"
	"time"
)

type evalTest struct {
	input    data.Value
	expected data.Value
}

func TestEvaluators(t *testing.T) {
	testCases := getTestCases()
	reg := &testFuncRegistry{ctx: core.NewContext(nil)}

	for _, testCase := range testCases {
		testCase := testCase
		ast := testCase.ast
		Convey(fmt.Sprintf("Given the AST Expression %v", ast), t, func() {

			Convey("When the expressions are evaluated", func() {
				flatExpr, err := ParserExprToFlatExpr(ast, reg)
				So(err, ShouldBeNil)
				eval, err := ExpressionToEvaluator(flatExpr, reg)
				So(err, ShouldBeNil)

				Convey("Then they should be evaluated correctly", FailureContinues, func() {
					// Although each test case should have an independent Convey block,
					// it makes tests 50 times slower. So, all inputs of each testCase
					// are evaluated within one Convey block, which makes debugging
					// very difficult.
					for _, tc := range testCase.inputs {
						input, expected := tc.input, tc.expected
						actual, err := eval.Eval(input)
						if expected == nil {
							So(err, ShouldNotBeNil)
						} else {
							So(err, ShouldBeNil)
							So(actual, ShouldResemble, expected)
						}
					}
				})
			})
		})
	}
}

func TestFoldableExecution(t *testing.T) {
	testCases := []struct {
		ast      parser.Expression
		foldable bool
		result   data.Value
	}{
		// Literals should always be independent of the input data
		{parser.NullLiteral{},
			true, data.Null{}},
		{parser.NumericLiteral{23},
			true, data.Int(23)},
		{parser.FloatLiteral{3.14},
			true, data.Float(3.14)},
		{parser.BoolLiteral{true},
			true, data.Bool(true)},
		{parser.StringLiteral{"foo"},
			true, data.String("foo")},
		// Access to column data should always be false
		{parser.RowMeta{"s", parser.TimestampMeta},
			false, nil},
		{parser.RowValue{"", "a"},
			false, nil},
		// Comparison operations
		{parser.BinaryOpAST{parser.Or, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Or, parser.BoolLiteral{false}, parser.BoolLiteral{true}},
			true, data.Bool(true)},
		{parser.BinaryOpAST{parser.And, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.And, parser.BoolLiteral{false}, parser.BoolLiteral{true}},
			true, data.Bool(false)},
		{parser.UnaryOpAST{parser.Not, parser.RowValue{"", "a"}},
			false, nil},
		{parser.UnaryOpAST{parser.Not, parser.BoolLiteral{false}},
			true, data.Bool(true)},
		{parser.BinaryOpAST{parser.Equal, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Equal, parser.NumericLiteral{7}, parser.NumericLiteral{7}},
			true, data.Bool(true)},
		{parser.BinaryOpAST{parser.Less, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Less, parser.FloatLiteral{5.5}, parser.NumericLiteral{6}},
			true, data.Bool(true)},
		// Computational Operations
		{parser.BinaryOpAST{parser.Plus, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Plus, parser.NumericLiteral{7}, parser.NumericLiteral{5}},
			true, data.Int(12)},
		{parser.BinaryOpAST{parser.Minus, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Minus, parser.NumericLiteral{7}, parser.NumericLiteral{5}},
			true, data.Int(2)},
		{parser.BinaryOpAST{parser.Multiply, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Multiply, parser.NumericLiteral{7}, parser.NumericLiteral{5}},
			true, data.Int(35)},
		{parser.BinaryOpAST{parser.Divide, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Divide, parser.NumericLiteral{6}, parser.NumericLiteral{2}},
			true, data.Int(3)},
		{parser.BinaryOpAST{parser.Modulo, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Modulo, parser.NumericLiteral{7}, parser.NumericLiteral{3}},
			true, data.Int(1)},
		{parser.BinaryOpAST{parser.Concat, parser.NumericLiteral{7}, parser.RowValue{"", "b"}},
			false, nil},
		{parser.BinaryOpAST{parser.Concat, parser.NumericLiteral{7}, parser.StringLiteral{"b"}},
			true, data.String("7b")},
		// Other
		{parser.AliasAST{parser.RowValue{"", "a"}, "hoge"},
			false, nil},
		{parser.AliasAST{parser.NumericLiteral{7}, "hoge"},
			true, data.Int(7)},
		{parser.TypeCastAST{parser.RowValue{"", "a"}, parser.Float},
			false, nil},
		{parser.TypeCastAST{parser.NumericLiteral{7}, parser.Float},
			true, data.Float(7.0)},
		{parser.FuncAppAST{parser.FuncName("now"),
			parser.ExpressionsAST{[]parser.Expression{}}, nil},
			false, nil},
		{parser.FuncAppAST{parser.FuncName("plusone"),
			parser.ExpressionsAST{[]parser.Expression{parser.RowValue{"", "a"}}}, nil},
			false, nil},
		{parser.FuncAppAST{parser.FuncName("plusone"),
			parser.ExpressionsAST{[]parser.Expression{parser.NumericLiteral{7}}}, nil},
			true, data.Int(8)},
		{parser.ArrayAST{parser.ExpressionsAST{[]parser.Expression{parser.RowValue{"", "a"}}}},
			false, nil},
		{parser.ArrayAST{parser.ExpressionsAST{[]parser.Expression{parser.NumericLiteral{7}}}},
			true, data.Array{data.Int(7)}},
		{parser.MapAST{[]parser.KeyValuePairAST{{"a", parser.RowValue{"", "a"}}}},
			false, nil},
		{parser.MapAST{[]parser.KeyValuePairAST{{"a", parser.NumericLiteral{7}}}},
			true, data.Map{"a": data.Int(7)}},
		{parser.ExpressionCaseAST{parser.RowValue{"", "a"}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{2}, parser.NumericLiteral{3}}}, parser.NullLiteral{}}},
			false, nil},
		{parser.ExpressionCaseAST{parser.NumericLiteral{2}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.RowValue{"", "a"}, parser.NumericLiteral{3}}}, parser.NullLiteral{}}},
			false, nil},
		{parser.ExpressionCaseAST{parser.NumericLiteral{3}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{2}, parser.RowValue{"", "a"}}}, parser.NullLiteral{}}},
			false, nil},
		{parser.ExpressionCaseAST{parser.NumericLiteral{3}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{2}, parser.NullLiteral{}}}, parser.RowValue{"", "a"}}},
			false, nil},
		{parser.ExpressionCaseAST{parser.NumericLiteral{3}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{3}, parser.NumericLiteral{4}}, {parser.NumericLiteral{3}, parser.RowValue{"", "a"}}}, parser.NullLiteral{}}},
			false, nil},
		{parser.ExpressionCaseAST{parser.BinaryOpAST{parser.Plus, parser.NumericLiteral{1}, parser.NumericLiteral{2}}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{3}, parser.NumericLiteral{4}}}, parser.NullLiteral{}}},
			true, data.Int(4)},
		{parser.ExpressionCaseAST{parser.BinaryOpAST{parser.Plus, parser.NumericLiteral{1}, parser.NumericLiteral{3}}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{3}, parser.NumericLiteral{4}}}, parser.NullLiteral{}}},
			true, data.Null{}},
		{parser.ExpressionCaseAST{parser.BinaryOpAST{parser.Plus, parser.NumericLiteral{1}, parser.NumericLiteral{3}}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{3}, parser.NumericLiteral{4}}, {parser.NumericLiteral{4}, parser.NumericLiteral{5}}}, parser.NullLiteral{}}},
			true, data.Int(5)},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.RowValue{"", "a"}, parser.NumericLiteral{3}}}, parser.NullLiteral{}},
			false, nil},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BoolLiteral{true}, parser.RowValue{"", "a"}}}, parser.NullLiteral{}},
			false, nil},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BoolLiteral{true}, parser.NullLiteral{}}}, parser.RowValue{"", "a"}},
			false, nil},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BoolLiteral{true}, parser.NumericLiteral{4}}, {parser.NumericLiteral{3}, parser.RowValue{"", "a"}}}, parser.NullLiteral{}},
			false, nil},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BoolLiteral{true}, parser.NumericLiteral{4}}}, parser.NullLiteral{}},
			true, data.Int(4)},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BoolLiteral{false}, parser.NumericLiteral{4}}}, parser.NullLiteral{}},
			true, data.Null{}},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BoolLiteral{false}, parser.NumericLiteral{4}}, {parser.BoolLiteral{true}, parser.NumericLiteral{5}}}, parser.NullLiteral{}},
			true, data.Int(5)},
	}

	reg := &testFuncRegistry{ctx: core.NewContext(nil)}

	for _, testCase := range testCases {
		testCase := testCase
		Convey(fmt.Sprintf("Given the AST Expression %v", testCase.ast), t, func() {

			Convey("Then the expression has the right foldability", func() {
				So(testCase.ast.Foldable(), ShouldEqual, testCase.foldable)

				Convey("And the executed result is correct", func() {
					res, err := EvaluateFoldable(testCase.ast, reg)

					if testCase.foldable {
						So(err, ShouldBeNil)
						So(res, ShouldResemble, testCase.result)
					} else {
						So(err, ShouldNotBeNil)
					}
				})

			})
		})
	}
}

func TestFuncAppConversion(t *testing.T) {
	Convey("Given a function registry", t, func() {
		reg := &testFuncRegistry{ctx: core.NewContext(nil)}

		Convey("When a function is known in the registry", func() {
			ast := parser.FuncAppAST{parser.FuncName("plusone"),
				parser.ExpressionsAST{[]parser.Expression{
					parser.RowValue{"", "a"},
				}}, nil}

			Convey("Then we obtain an evaluatable funcApp", func() {
				flatExpr, err := ParserExprToFlatExpr(ast, reg)
				So(err, ShouldBeNil)
				eval, err := ExpressionToEvaluator(flatExpr, reg)
				So(err, ShouldBeNil)
				So(eval, ShouldHaveSameTypeAs, &funcApp{})
			})
		})

		Convey("When the function is not known in the registry", func() {
			ast := parser.FuncAppAST{parser.FuncName("fun"),
				parser.ExpressionsAST{[]parser.Expression{
					parser.RowValue{"", "a"},
				}}, nil}

			Convey("Then converting to an Evaluator fails", func() {
				// we cannot even get the flat expression in that case
				_, err := ParserExprToFlatExpr(ast, reg)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When the function uses the ORDER BY clause", func() {
			ast := parser.FuncAppAST{parser.FuncName("plusone"),
				parser.ExpressionsAST{[]parser.Expression{
					parser.RowValue{"", "a"},
				}},
				[]parser.SortedExpressionAST{{parser.RowValue{"", "a"}, parser.Yes}}}

			Convey("Then converting to an Evaluator fails", func() {
				// we cannot even get the flat expression in that case
				_, err := ParserExprToFlatExpr(ast, reg)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual,
					"you cannot use ORDER BY in non-aggregate function 'plusone'")
			})
		})

		Convey("When the now() function is used", func() {
			ast := parser.FuncAppAST{parser.FuncName("now"),
				parser.ExpressionsAST{[]parser.Expression{}}, nil}

			Convey("Then we obtain an evaluatable timestampCast", func() {
				flatExpr, err := ParserExprToFlatExpr(ast, reg)
				So(err, ShouldBeNil)
				eval, err := ExpressionToEvaluator(flatExpr, reg)
				So(err, ShouldBeNil)
				So(eval, ShouldHaveSameTypeAs, &timestampCast{})
			})
		})
	})
}

func TestAggFuncAppConversion(t *testing.T) {
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))

	toString := udf.UnaryFunc(func(ctx *core.Context, v data.Value) (data.Value, error) {
		return data.String(v.String()), nil
	})
	reg.Register("f", toString)

	testCases := []struct {
		bql           string
		expectedError string
		expr          FlatExpression
		aggrs         map[string]FlatExpression
		inputs        []evalTest
	}{
		// order by a value that is already in the aggregate variables
		{"array_agg(a ORDER BY a ASC) FROM x [RANGE 1 TUPLES]", "",
			aggregateInputSorter{
				funcAppAST{"array_agg", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				[]sortExpression{sortExpression{aggInputRef{"g_f12cd6bc"}, true}},
				"ccd0ef22",
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
			},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// map does not contain the correct key
				{data.Map{"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				// map does not contain an array at that position
				{data.Map{"g_f12cd6bc": data.Int(17)}, nil},
				// correct input
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)}},
					data.Array{data.Int(1), data.Int(2)}},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(3), data.Int(1), data.Int(2)}},
					data.Array{data.Int(1), data.Int(2), data.Int(3)}},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(3), data.Int(2), data.Int(1)}},
					data.Array{data.Int(1), data.Int(2), data.Int(3)}},
			},
		},

		// order by a value that is not in the aggregate variables
		{"array_agg(a ORDER BY b DESC) FROM x [RANGE 1 TUPLES]", "",
			aggregateInputSorter{
				funcAppAST{"array_agg", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				[]sortExpression{sortExpression{aggInputRef{"g_77d2dd39"}, false}},
				"d7196f56",
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
				"g_77d2dd39": rowValue{"x", "b"},
			},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// map does not contain all correct keys
				{data.Map{"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_77d2dd39": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				// map does not contain an array at that position
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Int(17)}, nil},
				{data.Map{"g_f12cd6bc": data.Int(17),
					"g_77d2dd39": data.Array{data.Int(1), data.Int(2)}}, nil},
				// correct input
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Array{data.Int(1), data.Int(2)}},
					data.Array{data.Int(2), data.Int(1)}},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(3), data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Array{data.Int(5), data.Int(6), data.Int(4)}},
					data.Array{data.Int(1), data.Int(3), data.Int(2)}},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(3), data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Array{data.String("5"), data.Int(6), data.Int(4)}},
					data.Array{data.Int(3), data.Int(1), data.Int(2)}},
			},
		},

		// order by multiple values
		{"array_agg(a ORDER BY b DESC, a) FROM x [RANGE 1 TUPLES]", "",
			aggregateInputSorter{
				funcAppAST{"array_agg", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				[]sortExpression{sortExpression{aggInputRef{"g_77d2dd39"}, false},
					sortExpression{aggInputRef{"g_f12cd6bc"}, true}},
				"24925706",
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
				"g_77d2dd39": rowValue{"x", "b"},
			},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// map does not contain all correct keys
				{data.Map{"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_77d2dd39": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				// map does not contain an array at that position
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Int(17)}, nil},
				{data.Map{"g_f12cd6bc": data.Int(17),
					"g_77d2dd39": data.Array{data.Int(1), data.Int(2)}}, nil},
				// correct input
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Array{data.Int(1), data.Int(2)}},
					data.Array{data.Int(2), data.Int(1)}},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2), data.Int(3)},
					"g_77d2dd39": data.Array{data.Int(5), data.Int(6), data.Int(5)}},
					data.Array{data.Int(2), data.Int(1), data.Int(3)}},
			},
		},

		// use two different sorting orders for the same column
		{"array_agg(a ORDER BY b DESC)::string || array_agg(a ORDER BY b ASC)::string FROM x [RANGE 1 TUPLES]", "",
			binaryOpAST{parser.Concat,
				typeCastAST{
					aggregateInputSorter{
						funcAppAST{"array_agg",
							[]FlatExpression{aggInputRef{Ref: "g_f12cd6bc"}}},
						[]sortExpression{sortExpression{aggInputRef{"g_77d2dd39"},
							false}},
						"d7196f56"},
					parser.String},
				typeCastAST{
					aggregateInputSorter{
						funcAppAST{"array_agg",
							[]FlatExpression{aggInputRef{Ref: "g_f12cd6bc"}}},
						[]sortExpression{sortExpression{aggInputRef{"g_77d2dd39"},
							true}},
						"cd35e18d"},
					parser.String},
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
				"g_77d2dd39": rowValue{"x", "b"},
			},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// map does not contain all correct keys
				{data.Map{"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_77d2dd39": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				// map does not contain an array at that position
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Int(17)}, nil},
				{data.Map{"g_f12cd6bc": data.Int(17),
					"g_77d2dd39": data.Array{data.Int(1), data.Int(2)}}, nil},
				// correct input
				{data.Map{"g_f12cd6bc": data.Array{data.Int(1), data.Int(2)},
					"g_77d2dd39": data.Array{data.Int(3), data.Int(4)}},
					data.String("data.Array{2, 1}data.Array{1, 2}")},
			},
		},

		// order by a volatile expression
		{"array_agg(f(a) ORDER BY f(a)) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			aggregateInputSorter{
				funcAppAST{"array_agg", []FlatExpression{aggInputRef{"g_2523c3a2_0"}}},
				[]sortExpression{sortExpression{aggInputRef{"g_2523c3a2_1"}, true}},
				"cf2e24d7",
			},
			map[string]FlatExpression{
				"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
			},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// map does not contain all correct keys
				{data.Map{"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_2523c3a2_0": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				{data.Map{"g_2523c3a2_1": data.Array{data.Int(1), data.Int(2)},
					"a": data.Array{data.Int(1), data.Int(2)}}, nil},
				// map does not contain an array at that position
				{data.Map{"g_2523c3a2_0": data.Array{data.Int(1), data.Int(2)},
					"g_2523c3a2_1": data.Int(17)}, nil},
				{data.Map{"g_2523c3a2_0": data.Int(17),
					"g_2523c3a2_1": data.Array{data.Int(1), data.Int(2)}}, nil},
				// correct input
				{data.Map{"g_2523c3a2_0": data.Array{data.Int(1), data.Int(2)},
					"g_2523c3a2_1": data.Array{data.Int(4), data.Int(3)}},
					data.Array{data.Int(2), data.Int(1)}},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		Convey(fmt.Sprintf("Given the statement", testCase.bql), t, func() {
			p := parser.New()
			stmt := "CREATE STREAM x AS SELECT ISTREAM " + testCase.bql
			astUnchecked, _, err := p.ParseStmt(stmt)
			So(err, ShouldBeNil)
			So(astUnchecked, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
			ast := astUnchecked.(parser.CreateStreamAsSelectStmt).Select

			Convey("When we analyze it", func() {
				logPlan, err := Analyze(ast, reg)
				expectedError := testCase.expectedError
				if expectedError == "" {
					Convey("There is no error", func() {
						So(err, ShouldBeNil)
						So(len(logPlan.Projections), ShouldBeGreaterThanOrEqualTo, 1)
						proj := logPlan.Projections[0]
						So(proj.expr, ShouldResemble, testCase.expr)
						So(proj.aggrInputs, ShouldResemble, testCase.aggrs)

						eval, err := ExpressionToEvaluator(proj.expr, reg)
						So(err, ShouldBeNil)
						Convey("And the test cases should work", func() {
							for _, tc := range testCase.inputs {
								res, err := eval.Eval(tc.input)
								if tc.expected == nil {
									So(err, ShouldNotBeNil)
								} else {
									So(err, ShouldBeNil)
									So(res, ShouldResemble, tc.expected)
								}
							}
						})
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

var (
	// PlusOne is an example function that adds one to int and float Values.
	// It panics if the input is Null and returns an error for any other
	// type.
	PlusOne = udf.VariadicFunc(func(ctx *core.Context, vs ...data.Value) (data.Value, error) {
		if len(vs) != 1 {
			err := fmt.Errorf("cannot use %d parameters for unary function", len(vs))
			return nil, err
		}
		v := vs[0]
		if v.Type() == data.TypeInt {
			i, _ := data.AsInt(v)
			return data.Int(i + 1), nil
		} else if v.Type() == data.TypeFloat {
			f, _ := data.AsFloat(v)
			return data.Float(f + 1.0), nil
		} else if v.Type() == data.TypeNull {
			panic("null!")
		}
		return nil, fmt.Errorf("cannot add 1 to %v", v)
	})
	// MapLen is an example function that computes the length of a Map.
	// It returns an error if the input value is not a Map.
	MapLen = udf.UnaryFunc(func(ctx *core.Context, v data.Value) (data.Value, error) {
		m, err := data.AsMap(v)
		if err != nil {
			return nil, err
		}
		return data.Int(len(m)), nil
	})
)

// testFuncRegistry returns the PlusOne function above for any parameter.
type testFuncRegistry struct {
	ctx *core.Context
}

func (tfr *testFuncRegistry) Context() *core.Context {
	return tfr.ctx
}

func (tfr *testFuncRegistry) Lookup(name string, arity int) (udf.UDF, error) {
	if name == "plusone" {
		return PlusOne, nil
	} else if name == "maplen" && arity == 1 {
		return MapLen, nil
	}
	return nil, fmt.Errorf("no such function: %s", name)
}

func getTestCases() []struct {
	ast    parser.Expression
	inputs []evalTest
} {
	now := time.Now()

	// whatever binary operator we use (comparison or computation, but not
	// boolean/logical), if NULL is involved then the result should also always
	// be null
	nullOps := []evalTest{
		// null vs. *
		{data.Map{"a": data.Null{},
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.Bool(true)}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.Int(3)}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.Float(3.14)}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.String("hoge")}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.Blob("hoge")}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.Timestamp(now)}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.Array{data.Int(2)}}, data.Null{}},
		{data.Map{"a": data.Null{},
			"b": data.Map{"b": data.Int(3)}}, data.Null{}},
		// * vs. null
		{data.Map{"a": data.Bool(true),
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.Int(3),
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.Float(3.14),
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.String("hoge"),
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Null{}}, data.Null{}},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Null{}}, data.Null{}},
	}

	// these are all type combinations that are so incompatible that
	// they cannot be compared with respect to less/greater and also
	// cannot be added etc.
	incomparables := append([]evalTest{
		// not a map:
		{data.Int(17), nil},
		// keys not present:
		{data.Map{"x": data.Int(17)}, nil},
		// only left present => error
		{data.Map{"a": data.Bool(true)}, nil},
		{data.Map{"a": data.Int(17)}, nil},
		{data.Map{"a": data.Float(3.14)}, nil},
		{data.Map{"a": data.String("日本語")}, nil},
		{data.Map{"a": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
		// only right present => error
		{data.Map{"b": data.Bool(true)}, nil},
		{data.Map{"b": data.Int(17)}, nil},
		{data.Map{"b": data.Float(3.14)}, nil},
		{data.Map{"b": data.String("日本語")}, nil},
		{data.Map{"b": data.Blob("hoge")}, nil},
		{data.Map{"b": data.Timestamp(now)}, nil},
		{data.Map{"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"b": data.Map{"b": data.Int(3)}}, nil},
		// bool vs *
		{data.Map{"a": data.Bool(true),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// int vs. *
		{data.Map{"a": data.Int(3),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// float vs *
		{data.Map{"a": data.Float(3.14),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// string vs *
		{data.Map{"a": data.String("hoge"),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// blob vs *
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// timestamp vs *
		{data.Map{"a": data.Timestamp(now),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// array vs *
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Map{"b": data.Int(3)}}, nil},
		// map vs *
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Map{"b": data.Int(3)}}, nil},
	}, nullOps...)

	// we should check that every AST expression maps to
	// an evaluator with the correct behavior
	testCases := []struct {
		ast    parser.Expression
		inputs []evalTest
	}{
		// Literals should always be independent of the input data
		{parser.NullLiteral{},
			[]evalTest{
				{data.Int(17), data.Null{}},
				{data.String(""), data.Null{}},
			},
		},
		{parser.NumericLiteral{23},
			[]evalTest{
				{data.Int(17), data.Int(23)},
				{data.String(""), data.Int(23)},
			},
		},
		{parser.FloatLiteral{3.14},
			[]evalTest{
				{data.Int(17), data.Float(3.14)},
				{data.String(""), data.Float(3.14)},
			},
		},
		{parser.BoolLiteral{true},
			[]evalTest{
				{data.Int(17), data.Bool(true)},
				{data.String(""), data.Bool(true)},
			},
		},
		{parser.BoolLiteral{false},
			[]evalTest{
				{data.Int(17), data.Bool(false)},
				{data.String(""), data.Bool(false)},
			},
		},
		{parser.StringLiteral{"foo"},
			[]evalTest{
				{data.Int(17), data.String("foo")},
				{data.String(""), data.String("foo")},
			},
		},
		// Extracting the timestamp should find the timestamp at the
		// correct position
		{parser.RowMeta{"s", parser.TimestampMeta},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present, but wrong type
				{data.Map{"s:meta:TS": data.Int(17)}, nil},
				// key present and correct type
				{data.Map{"s:meta:TS": data.Timestamp(now)}, data.Timestamp(now)},
			},
		},
		// Access to columns/keys should return the same values
		{parser.RowValue{"", "a"},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present
				{data.Map{"a": data.Null{}}, data.Null{}},
				{data.Map{"a": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17)}, data.Int(17)},
				{data.Map{"a": data.Float(3.14)}, data.Float(3.14)},
				{data.Map{"a": data.String("日本語")}, data.String("日本語")},
				{data.Map{"a": data.Blob("hoge")}, data.Blob("hoge")},
				{data.Map{"a": data.Timestamp(now)}, data.Timestamp(now)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Array{data.Int(2)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Map{"b": data.Int(3)}},
			},
		},
		{parser.RowValue{"", `a["ho""ge"]`},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present
				{data.Map{"a": data.Map{`ho"ge`: data.Int(3)}}, data.Int(3)},
			},
		},
		// Access to columns/keys should return the same values
		{parser.AliasAST{parser.RowValue{"", "a"}, "hoge"},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present
				{data.Map{"a": data.Null{}}, data.Null{}},
				{data.Map{"a": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17)}, data.Int(17)},
				{data.Map{"a": data.Float(3.14)}, data.Float(3.14)},
				{data.Map{"a": data.String("日本語")}, data.String("日本語")},
				{data.Map{"a": data.Blob("hoge")}, data.Blob("hoge")},
				{data.Map{"a": data.Timestamp(now)}, data.Timestamp(now)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Array{data.Int(2)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Map{"b": data.Int(3)}},
			},
		},
		/// Combined operations
		// Or
		{parser.BinaryOpAST{parser.Or, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left key present and evaluates to true => right is not necessary
				{data.Map{"a": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// only left key present and evaluates to false => error
				{data.Map{"a": data.Bool(false)}, nil},
				{data.Map{"a": data.Int(0)}, nil},
				{data.Map{"a": data.Float(0.0)}, nil},
				{data.Map{"a": data.String("")}, nil},
				{data.Map{"a": data.Blob("")}, nil},
				{data.Map{"a": data.Timestamp{}}, nil},
				{data.Map{"a": data.Array{}}, nil},
				{data.Map{"a": data.Map{}}, nil},
				// left key evalues to false and right to true => true
				{data.Map{"a": data.Int(0),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// left key evalues to false and right to false => false
				{data.Map{"a": data.Int(0),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Map{}}, data.Bool(false)},
				// null comparison
				{data.Map{"a": data.Bool(true),
					"b": data.Null{}}, data.Bool(true)},
				{data.Map{"a": data.Bool(false),
					"b": data.Null{}}, data.Null{}},
				{data.Map{"a": data.Null{},
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Null{},
					"b": data.Bool(false)}, data.Null{}},
			},
		},
		// And
		{parser.BinaryOpAST{parser.And, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left key present and evaluates to false => right is not necessary
				{data.Map{"a": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Map{}}, data.Bool(false)},
				// only left key present and evaluates to true => error
				{data.Map{"a": data.Bool(true)}, nil},
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// left key evalues to true and right to true => true
				{data.Map{"a": data.Int(1),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// left key evalues to true and right to false => false
				{data.Map{"a": data.Int(1),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{}}, data.Bool(false)},
				// null comparison
				{data.Map{"a": data.Bool(true),
					"b": data.Null{}}, data.Null{}},
				{data.Map{"a": data.Bool(false),
					"b": data.Null{}}, data.Bool(false)},
				{data.Map{"a": data.Null{},
					"b": data.Bool(true)}, data.Null{}},
				{data.Map{"a": data.Null{},
					"b": data.Bool(false)}, data.Bool(false)},
			},
		},
		// Not
		{parser.UnaryOpAST{parser.Not, parser.RowValue{"", "a"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and false => true
				{data.Map{"a": data.Bool(false)}, data.Bool(true)},
				// key present and false-like value => error
				{data.Map{"a": data.Int(0)}, nil},
				{data.Map{"a": data.Float(0.0)}, nil},
				{data.Map{"a": data.String("")}, nil},
				{data.Map{"a": data.Blob("")}, nil},
				{data.Map{"a": data.Timestamp{}}, nil},
				{data.Map{"a": data.Array{}}, nil},
				{data.Map{"a": data.Map{}}, nil},
				// key present and true => false
				{data.Map{"a": data.Bool(true)}, data.Bool(false)},
				// key present and true-like value => error
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// null comparison
				{data.Map{"a": data.Null{}}, data.Null{}},
			},
		},
		/// Comparison Operations
		// Equal
		{parser.BinaryOpAST{parser.Equal, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left present => error
				{data.Map{"a": data.Bool(true)}, nil},
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// only right present => error
				{data.Map{"b": data.Bool(true)}, nil},
				{data.Map{"b": data.Int(17)}, nil},
				{data.Map{"b": data.Float(3.14)}, nil},
				{data.Map{"b": data.String("日本語")}, nil},
				{data.Map{"b": data.Blob("hoge")}, nil},
				{data.Map{"b": data.Timestamp(now)}, nil},
				{data.Map{"b": data.Array{data.Int(2)}}, nil},
				{data.Map{"b": data.Map{"b": data.Int(3)}}, nil},
				// left and right present and equal => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17),
					"b": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.String("日本語"),
					"b": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Blob("hoge"),
					"b": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Array{data.Int(2)},
					"b": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Map{"b": data.Int(3)},
					"b": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// left and right present and not equal => false
				{data.Map{"a": data.Int(1),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{}}, data.Bool(false)},
			}, nullOps...),
		},
		// Less
		{parser.BinaryOpAST{parser.Less, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => true
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, data.Bool(true)},
				// left and right present and comparable and equal => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and comparable and left is greater => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(false)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now.Add(time.Second)),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// LessOrEqual
		{parser.BinaryOpAST{parser.LessOrEqual, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => true
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, data.Bool(true)},
				// left and right present and comparable and equal => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and comparable and left is greater => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(false)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now.Add(time.Second)),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Greater
		{parser.BinaryOpAST{parser.Greater, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => false
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, data.Bool(false)},
				// left and right present and comparable and equal => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and comparable and left is greater => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(true)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now.Add(time.Second)),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// GreaterOrEqual
		{parser.BinaryOpAST{parser.GreaterOrEqual, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => false
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, data.Bool(false)},
				// left and right present and comparable and equal => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and comparable and left is greater => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(true)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now.Add(time.Second)),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// NotEqual
		{parser.BinaryOpAST{parser.NotEqual, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left present => error
				{data.Map{"a": data.Bool(true)}, nil},
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// only right present => error
				{data.Map{"b": data.Bool(true)}, nil},
				{data.Map{"b": data.Int(17)}, nil},
				{data.Map{"b": data.Float(3.14)}, nil},
				{data.Map{"b": data.String("日本語")}, nil},
				{data.Map{"b": data.Blob("hoge")}, nil},
				{data.Map{"b": data.Timestamp(now)}, nil},
				{data.Map{"b": data.Array{data.Int(2)}}, nil},
				{data.Map{"b": data.Map{"b": data.Int(3)}}, nil},
				// left and right present and equal => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(17),
					"b": data.Int(17)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.String("日本語"),
					"b": data.String("日本語")}, data.Bool(false)},
				{data.Map{"a": data.Blob("hoge"),
					"b": data.Blob("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(false)},
				{data.Map{"a": data.Array{data.Int(2)},
					"b": data.Array{data.Int(2)}}, data.Bool(false)},
				{data.Map{"a": data.Map{"b": data.Int(3)},
					"b": data.Map{"b": data.Int(3)}}, data.Bool(false)},
				// left and right present and not equal => true
				{data.Map{"a": data.Int(1),
					"b": data.Bool(false)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(0)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(0.0)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.String("")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp{}}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{}}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{}}, data.Bool(true)},
			}, nullOps...),
		},
		// Concatenation
		{parser.BinaryOpAST{parser.Concat, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left present => error
				{data.Map{"a": data.Bool(true)}, nil},
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// only right present => error
				{data.Map{"b": data.Bool(true)}, nil},
				{data.Map{"b": data.Int(17)}, nil},
				{data.Map{"b": data.Float(3.14)}, nil},
				{data.Map{"b": data.String("日本語")}, nil},
				{data.Map{"b": data.Blob("hoge")}, nil},
				{data.Map{"b": data.Timestamp(now)}, nil},
				{data.Map{"b": data.Array{data.Int(2)}}, nil},
				{data.Map{"b": data.Map{"b": data.Int(3)}}, nil},
				// left and right present
				{data.Map{"a": data.Int(1),
					"b": data.Bool(false)}, data.String("1false")},
				{data.Map{"a": data.Int(1),
					"b": data.Int(0)}, data.String("10")},
				{data.Map{"a": data.Int(1),
					"b": data.Float(0.0)}, data.String("10")},
				{data.Map{"a": data.Int(1),
					"b": data.String("")}, data.String("1")},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("")}, data.String("1")},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp{}}, data.String("10001-01-01T00:00:00Z")},
				// TODO this is not the best possible form...
				{data.Map{"a": data.Int(1),
					"b": data.Array{}}, data.String("1data.Array{}")},
				{data.Map{"a": data.Int(1),
					"b": data.Map{}}, data.String("1data.Map{}")},
			}, nullOps...),
		},
		// IsNull
		{parser.BinaryOpAST{parser.Is, parser.RowValue{"", "a"}, parser.NullLiteral{}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// left present and not null => false
				{data.Map{"a": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(17)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("日本語")}, data.Bool(false)},
				{data.Map{"a": data.Blob("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now)}, data.Bool(false)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Bool(false)},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Bool(false)},
				// left present and null => true
				{data.Map{"a": data.Null{}}, data.Bool(true)},
			},
		},
		// IsNotNull
		{parser.BinaryOpAST{parser.IsNot, parser.RowValue{"", "a"}, parser.NullLiteral{}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// left present and not null => false
				{data.Map{"a": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// left present and null => true
				{data.Map{"a": data.Null{}}, data.Bool(false)},
			},
		},
		/// Computational Operations
		// Plus
		{parser.BinaryOpAST{parser.Plus, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be added
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(5)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) + 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 + float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(3.14 + 3.15)},
				// left and right present and cannot be added
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Minus
		{parser.BinaryOpAST{parser.Minus, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be subtracted
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(-1)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) - 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 - float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(float64(3.14) - 3.15)},
				// left and right present and cannot be subtracted
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Multiply
		{parser.BinaryOpAST{parser.Multiply, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be multiplied
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(6)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) * 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 * float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(float64(3.14) * 3.15)},
				// left and right present and cannot be multiplied
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Divide
		{parser.BinaryOpAST{parser.Divide, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be divided
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(0)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) / 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 / float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(float64(3.14) / 3.15)},
				// division by zero
				{data.Map{"a": data.Int(2),
					"b": data.Int(0)}, nil},
				{data.Map{"a": data.Int(3),
					"b": data.Float(0)}, data.Float(math.Inf(1))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(0)}, data.Float(math.Inf(1))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(0)}, data.Float(math.Inf(1))},
				// left and right present and cannot be divided
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Modulo
		{parser.BinaryOpAST{parser.Modulo, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be moduled
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(2)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(3.0)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(3.14)},
				// modulo by zero
				{data.Map{"a": data.Int(2),
					"b": data.Int(0)}, nil},
				// TODO add a way to check for IsNaN()
				/*
					{data.Map{"a": data.Int(3),
						"b": data.Float(0)}, data.Float(math.NaN())},
					{data.Map{"a": data.Float(3.14),
						"b": data.Int(0)}, data.Float(math.NaN())},
					{data.Map{"a": data.Float(3.14),
						"b": data.Float(0)}, data.Float(math.NaN())},
				*/
				// left and right present and cannot be moduled
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now.Add(time.Second))}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Unary Minus
		{parser.UnaryOpAST{parser.UnaryMinus, parser.RowValue{"", "a"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and number-like => negative
				{data.Map{"a": data.Int(17)}, data.Int(-17)},
				{data.Map{"a": data.Float(3.14)}, data.Float(-3.14)},
				{data.Map{"a": data.Int(-17)}, data.Int(17)},
				{data.Map{"a": data.Float(-3.14)}, data.Float(3.14)},
				{data.Map{"a": data.Int(0)}, data.Int(0)},
				{data.Map{"a": data.Float(0.0)}, data.Float(-0.0)},
				// key present and other data type => error
				{data.Map{"a": data.Bool(false)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// null comparison
				{data.Map{"a": data.Null{}}, data.Null{}},
			},
		},
		// Type Cast
		{parser.TypeCastAST{parser.RowValue{"", "a"}, parser.Int},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and convertable => ok
				{data.Map{"a": data.Int(17)}, data.Int(17)},
				{data.Map{"a": data.Float(3.14)}, data.Int(3)},
				{data.Map{"a": data.Int(-17)}, data.Int(-17)},
				{data.Map{"a": data.Float(-3.14)}, data.Int(-3)},
				{data.Map{"a": data.Int(0)}, data.Int(0)},
				{data.Map{"a": data.Float(0.0)}, data.Int(0)},
				{data.Map{"a": data.Bool(false)}, data.Int(0)},
				// null propagation
				{data.Map{"a": data.Null{}}, data.Null{}},
				// key present and other data type => error
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
			},
		},
		{parser.TypeCastAST{parser.RowValue{"", "a"}, parser.String},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and number-like => conversion
				{data.Map{"a": data.Int(17)}, data.String("17")},
				{data.Map{"a": data.Float(3.14)}, data.String("3.14")},
				{data.Map{"a": data.Int(-17)}, data.String("-17")},
				{data.Map{"a": data.Float(-3.14)}, data.String("-3.14")},
				{data.Map{"a": data.Int(0)}, data.String("0")},
				{data.Map{"a": data.Float(0.0)}, data.String("0")},
				{data.Map{"a": data.Bool(false)}, data.String("false")},
				{data.Map{"a": data.String("日本語")}, data.String("日本語")},
				{data.Map{"a": data.Blob("hoge")}, data.String("hoge")},
				{data.Map{"a": data.Array{data.Int(2)}}, data.String("data.Array{2}")},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.String("data.Map{\"b\":3}")},
				// null propagation
				{data.Map{"a": data.Null{}}, data.Null{}},
			},
		},
		/// Function Application
		{parser.FuncAppAST{parser.FuncName("plusone"),
			parser.ExpressionsAST{[]parser.Expression{parser.RowValue{"", "a"}}}, nil},
			// NB. This only tests the behavior of funcApp.Eval.
			// It does *not* test the function registry, mismatch
			// in parameter counts or any particular function.
			[]evalTest{
				// function returns good result
				{data.Map{"a": data.Int(16)}, data.Int(17)},
				{data.Map{"a": data.Float(16.0)}, data.Float(17.0)},
				// function errors
				{data.Map{"x": data.Int(17)}, nil},
				{data.Map{"a": data.Bool(false)}, nil},
				// function panics
				{data.Map{"a": data.Null{}}, nil},
			},
		},
		/// JSON-like Structures
		{parser.ArrayAST{parser.ExpressionsAST{[]parser.Expression{
			parser.NumericLiteral{2}, parser.RowValue{"", "a"}}}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and number-like => conversion
				{data.Map{"a": data.Int(17)}, data.Array{data.Int(2), data.Int(17)}},
				{data.Map{"a": data.Float(3.14)}, data.Array{data.Int(2), data.Float(3.14)}},
				{data.Map{"a": data.String("日本語")}, data.Array{data.Int(2), data.String("日本語")}},
				{data.Map{"a": data.Array{data.Int(3)}}, data.Array{data.Int(2), data.Array{data.Int(3)}}},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Array{data.Int(2), data.Map{"b": data.Int(3)}}},
				{data.Map{"a": data.Null{}}, data.Array{data.Int(2), data.Null{}}},
			},
		},
		{parser.MapAST{[]parser.KeyValuePairAST{{"two", parser.NumericLiteral{2}},
			{"a", parser.RowValue{"", "a"}}}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and number-like => conversion
				{data.Map{"a": data.Int(17)}, data.Map{"two": data.Int(2), "a": data.Int(17)}},
				{data.Map{"a": data.Float(3.14)}, data.Map{"two": data.Int(2), "a": data.Float(3.14)}},
				{data.Map{"a": data.String("日本語")}, data.Map{"two": data.Int(2), "a": data.String("日本語")}},
				{data.Map{"a": data.Array{data.Int(3)}}, data.Map{"two": data.Int(2), "a": data.Array{data.Int(3)}}},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Map{"two": data.Int(2), "a": data.Map{"b": data.Int(3)}}},
				{data.Map{"a": data.Null{}}, data.Map{"two": data.Int(2), "a": data.Null{}}},
			},
		},
		// CASE a WHEN ... THEN ... ELSE ... END
		{parser.ExpressionCaseAST{parser.RowValue{"", "a"}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{2}, parser.NumericLiteral{3}}}, parser.NullLiteral{}}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and matches
				{data.Map{"a": data.Int(2)}, data.Int(3)},
				// key present and doesn't match
				{data.Map{"a": data.Int(3)}, data.Null{}},
				{data.Map{"a": data.String("3")}, data.Null{}},
			},
		},
		{parser.ExpressionCaseAST{parser.RowValue{"", "a"}, parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.NumericLiteral{2}, parser.NumericLiteral{3}},
			{parser.StringLiteral{"3"}, parser.StringLiteral{"hoge"}}}, parser.NumericLiteral{7}}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and matches first
				{data.Map{"a": data.Int(2)}, data.Int(3)},
				// key present and matches second
				{data.Map{"a": data.String("3")}, data.String("hoge")},
				// key present and doesn't match
				{data.Map{"a": data.Int(3)}, data.Int(7)},
			},
		},
		// CASE WHEN ... THEN ... ELSE ... END
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BinaryOpAST{parser.Equal, parser.RowValue{"", "a"}, parser.NumericLiteral{2}}, parser.NumericLiteral{3}}}, parser.NullLiteral{}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and matches
				{data.Map{"a": data.Int(2)}, data.Int(3)},
				// key present and doesn't match
				{data.Map{"a": data.Int(3)}, data.Null{}},
				{data.Map{"a": data.String("3")}, data.Null{}},
			},
		},
		{parser.ConditionCaseAST{[]parser.WhenThenPairAST{
			{parser.BinaryOpAST{parser.Equal, parser.RowValue{"", "a"}, parser.NumericLiteral{2}}, parser.NumericLiteral{3}},
			{parser.BinaryOpAST{parser.Greater, parser.RowValue{"", "a"}, parser.NumericLiteral{2}}, parser.NumericLiteral{7}}}, parser.StringLiteral{"hoge"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present and matches first
				{data.Map{"a": data.Int(2)}, data.Int(3)},
				// key present and matches second
				{data.Map{"a": data.Int(6)}, data.Int(7)},
				// key present and doesn't match
				{data.Map{"a": data.Int(1)}, data.String("hoge")},
			},
		},
		// Using now() should find the timestamp at the
		// correct position
		{parser.FuncAppAST{parser.FuncName("now"),
			parser.ExpressionsAST{[]parser.Expression{}}, nil},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present, but wrong type
				{data.Map{":meta:NOW": data.Int(17)}, nil},
				// key present and correct type
				{data.Map{":meta:NOW": data.Timestamp(now)}, data.Timestamp(now)},
			},
		},
		/// Wildcard
		{parser.Wildcard{},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// a map, but elements are not maps
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "b": data.Int(17)}, nil},
				// a map and the subelement is a map
				{data.Map{"a": data.Map{"b": data.Int(3)}},
					data.Map{"b": data.Int(3)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, ":meta:NOW": data.Int(17)},
					data.Map{"b": data.Int(3)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "c": data.Map{"d": data.Int(4)}},
					data.Map{"b": data.Int(3), "d": data.Int(4)}},
			}},
		{parser.Wildcard{"a"},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present
				{data.Map{"x": data.Map{"b": data.Int(3)}}, nil},
				// a map, but elements are not maps
				{data.Map{"a": data.Int(17)}, nil},
				// a map and the subelement is a map
				{data.Map{"a": data.Map{"b": data.Int(3)}},
					data.Map{"b": data.Int(3)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, ":meta:NOW": data.Int(17)},
					data.Map{"b": data.Int(3)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "c": data.Map{"d": data.Int(4)}},
					data.Map{"b": data.Int(3)}},
			}},
		{parser.ArrayAST{parser.ExpressionsAST{[]parser.Expression{parser.NumericLiteral{2},
			parser.Wildcard{}}}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// a map, but elements are not maps
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "b": data.Int(17)}, nil},
				// a map and the subelement is a map
				{data.Map{"a": data.Map{"b": data.Int(3)}},
					data.Array{data.Int(2), data.Map{"b": data.Int(3)}}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, ":meta:NOW": data.Int(17)},
					data.Array{data.Int(2), data.Map{"b": data.Int(3)}}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "c": data.Map{"d": data.Int(4)}},
					data.Array{data.Int(2), data.Map{"b": data.Int(3), "d": data.Int(4)}}},
			},
		},
		{parser.MapAST{[]parser.KeyValuePairAST{{"two", parser.NumericLiteral{2}},
			{"x", parser.Wildcard{"a"}}}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present
				{data.Map{"x": data.Map{"b": data.Int(3)}}, nil},
				// a map, but elements are not maps
				{data.Map{"a": data.Int(17)}, nil},
				// a map and the subelement is a map
				{data.Map{"a": data.Map{"b": data.Int(3)}},
					data.Map{"two": data.Int(2), "x": data.Map{"b": data.Int(3)}}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, ":meta:NOW": data.Int(17)},
					data.Map{"two": data.Int(2), "x": data.Map{"b": data.Int(3)}}},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "c": data.Map{"d": data.Int(4)}},
					data.Map{"two": data.Int(2), "x": data.Map{"b": data.Int(3)}}},
			},
		},
		{parser.FuncAppAST{parser.FuncName("maplen"),
			parser.ExpressionsAST{[]parser.Expression{parser.Wildcard{}}}, nil},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// a map, but elements are not maps
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "b": data.Int(17)}, nil},
				// a map and the subelement is a map
				{data.Map{"a": data.Map{"b": data.Int(3)}},
					data.Int(1)},
				{data.Map{"a": data.Map{"b": data.Int(3)}, ":meta:NOW": data.Int(17)},
					data.Int(1)},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "c": data.Map{"d": data.Int(4)}},
					data.Int(2)},
			},
		},
		{parser.FuncAppAST{parser.FuncName("maplen"),
			parser.ExpressionsAST{[]parser.Expression{parser.Wildcard{"a"}}}, nil},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present
				{data.Map{"x": data.Map{"b": data.Int(3)}}, nil},
				// a map, but elements are not maps
				{data.Map{"a": data.Int(17)}, nil},
				// a map and the subelement is a map
				{data.Map{"a": data.Map{"b": data.Int(3)}},
					data.Int(1)},
				{data.Map{"a": data.Map{"b": data.Int(3)}, ":meta:NOW": data.Int(17)},
					data.Int(1)},
				{data.Map{"a": data.Map{"b": data.Int(3)}, "c": data.Map{"d": data.Int(4)}},
					data.Int(1)},
			},
		},
	}
	return testCases
}

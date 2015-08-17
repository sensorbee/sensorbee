package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	_ "pfi/sensorbee/sensorbee/bql/udf/builtin"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"reflect"
	"testing"
)

type dummyAggregate struct {
}

func (f *dummyAggregate) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	if len(args) == 1 || len(args) == 2 {
		arr, err := data.AsArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("udaf needs an array input, not %v", args[0])
		}
		c := len(arr)
		if len(args) == 1 {
			return data.String(fmt.Sprintf("%d", c)), nil
		}
		return data.String(fmt.Sprintf("%d+%s", c, args[1])), nil
	}
	return nil, fmt.Errorf("udaf takes one or two arguments")
}

func (f *dummyAggregate) Accept(arity int) bool {
	return arity == 1 || arity == 2 || arity == 3
}

func (f *dummyAggregate) IsAggregationParameter(k int) bool {
	return k == 0 || k == 2
}

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
	wc := parser.Wildcard{}
	ts := parser.RowMeta{"", parser.TimestampMeta}
	tA := parser.RowValue{"t", "a"}
	tB := parser.RowValue{"t", "b"}
	tC := parser.RowValue{"t", "c"}
	tWc := parser.Wildcard{"t"}
	tTs := parser.RowMeta{"t", parser.TimestampMeta}
	xA := parser.RowValue{"x", "a"}
	xB := parser.RowValue{"x", "b"}

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
		// SELECT *   -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{wc}},
		}, "need at least one relation to select from"},
		// SELECT t:a -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{tA}},
		}, "need at least one relation to select from"},
		// SELECT t:ts() -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{tTs}},
		}, "need at least one relation to select from"},
		// SELECT t:*   -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{tWc}},
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
		// SELECT f(a ORDER BY b)  FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{
				parser.FuncAppAST{"f", parser.ExpressionsAST{[]parser.Expression{a}},
					[]parser.SortedExpressionAST{{b, parser.UnspecifiedKeyword}}},
			}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT 2        FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT *        FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{wc}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT a, *  FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, wc}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:a      FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:a, t:b FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA, tB}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:a, t:* FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA, tWc}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:a, t:ts() FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA, tTs}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT 2, t:a   FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two, tA}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT t:*      FROM t -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tWc}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT a, t:b   FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, tA}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT f(a ORDER BY t:b)  FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{
				parser.FuncAppAST{"f", parser.ExpressionsAST{[]parser.Expression{a}},
					[]parser.SortedExpressionAST{{tB, parser.UnspecifiedKeyword}}},
			}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT f(t:a ORDER BY b)  FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST: parser.ProjectionsAST{[]parser.Expression{
				parser.FuncAppAST{"f", parser.ExpressionsAST{[]parser.Expression{tA}},
					[]parser.SortedExpressionAST{{b, parser.UnspecifiedKeyword}}},
			}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT a, t:*   FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, tWc}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT t:a, *   FROM t -> OK (this is special about the wildcard!!)
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA, wc}},
			WindowedFromAST: singleFrom,
		}, ""},
		// SELECT a, t:ts() FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a, tTs}},
			WindowedFromAST: singleFrom,
		}, "cannot refer to relations"},
		// SELECT x:a      FROM t -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{xA}},
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
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
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
		// SELECT *   FROM t WHERE b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{wc}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{b},
		}, ""},
		// SELECT t:a FROM t WHERE b   -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tWc}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{b},
		}, "cannot refer to relations"},
		// SELECT t:* FROM t WHERE b   -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{b},
		}, "cannot refer to relations"},
		// SELECT a   FROM t WHERE t:b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{tB},
		}, "cannot refer to relations"},
		// SELECT *   FROM t WHERE t:b -> OK (this is special about wildcard!)
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{wc}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{tB},
		}, ""},
		// SELECT 2   FROM t WHERE t:b -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{tB},
		}, ""},
		// SELECT t:a FROM t WHERE t:b -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{tB},
		}, ""},
		// SELECT 2   FROM t WHERE x:b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			FilterAST:       parser.FilterAST{xB},
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
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
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
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b}},
		}, "cannot refer to relations"},
		// SELECT a   FROM t GROUP BY t:b      -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{tB}},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t GROUP BY t:b      -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{tB}},
		}, ""},
		// SELECT t:a FROM t GROUP BY t:b      -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{tB}},
		}, ""},
		// SELECT t:a FROM t GROUP BY t:b, t:c -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{tB, tC}},
		}, ""},
		// SELECT t:a FROM t GROUP BY b, t:b   -> NG (same table with multiple aliases)
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{b, tB}},
		}, "cannot refer to relations"},
		// SELECT 2   FROM t GROUP BY x:b      -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{two}},
			WindowedFromAST: singleFrom,
			GroupingAST:     parser.GroupingAST{[]parser.Expression{xB}},
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
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
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
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{b},
		}, "cannot refer to relations"},
		// SELECT t:a FROM t HAVING t:b   -> OK
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{tA}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{tB},
		}, ""},
		// SELECT a   FROM t HAVING t:b -> NG
		{&parser.SelectStmt{
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
			HavingAST:       parser.HavingAST{tB},
		}, "cannot refer to relations"},
	}

	emitterTestCases := []analyzeTest{
		// SELECT ISTREAM                                      a FROM t -> OK
		{&parser.SelectStmt{
			EmitterAST:      parser.EmitterAST{parser.Istream, nil},
			ProjectionsAST:  parser.ProjectionsAST{[]parser.Expression{a}},
			WindowedFromAST: singleFrom,
		}, ""},
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
			ast := parser.SelectStmt{
				EmitterAST:      emitter,
				ProjectionsAST:  selectAst.ProjectionsAST,
				WindowedFromAST: selectAst.WindowedFromAST,
				FilterAST:       selectAst.FilterAST,
				GroupingAST:     selectAst.GroupingAST,
				HavingAST:       selectAst.HavingAST,
			}

			Convey("When we analyze it", func() {
				// the two functions below were just a call to
				// `Analyze` before, but now `Analyze` does more
				// than what we want to check for here
				err := makeRelationAliases(&ast)
				if err == nil {
					err = validateReferences(&ast)
				}
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
			ast := parser.SelectStmt{
				EmitterAST:      parser.EmitterAST{parser.Istream, nil},
				ProjectionsAST:  selectAst.ProjectionsAST,
				WindowedFromAST: myFrom,
				FilterAST:       selectAst.FilterAST,
				GroupingAST:     selectAst.GroupingAST,
				HavingAST:       selectAst.HavingAST,
			}

			Convey("When we analyze it", func() {
				// the two functions below were just a call to
				// `Analyze` before, but now `Analyze` does more
				// than what we want to check for here
				err := makeRelationAliases(&ast)
				if err == nil {
					err = validateReferences(&ast)
				}
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

	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	for _, testCase := range testCases {
		testCase := testCase
		selectAst := testCase.input

		Convey(fmt.Sprintf("Given the AST %+v", selectAst), t, func() {
			ast := parser.SelectStmt{
				EmitterAST:      parser.EmitterAST{parser.Istream, nil},
				ProjectionsAST:  selectAst.ProjectionsAST,
				WindowedFromAST: selectAst.WindowedFromAST,
			}

			Convey("When we analyze it", func() {
				_, err := Analyze(ast, reg)
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
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))

	dummyFun := &dummyAggregate{}
	toString := udf.UnaryFunc(func(ctx *core.Context, v data.Value) (data.Value, error) {
		return data.String(v.String()), nil
	})

	reg.Register("udaf", dummyFun)
	reg.Register("f", toString)
	reg.Register("g", toString)

	testCases := []struct {
		bql           string
		expectedError string
		expr          FlatExpression
		aggrs         map[string]FlatExpression
	}{
		// a is no aggregate call, so the `aggrs` list is empty
		// and the selected expression is transformed normally
		{"a FROM x [RANGE 1 TUPLES]", "",
			rowValue{"x", "a"},
			nil},

		// f(a) is no aggregate call, so the `aggrs` list is empty
		// and the selected expression is transformed normally
		{"f(a) FROM x [RANGE 1 TUPLES]", "",
			funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
			nil},

		// f(*) is no aggregate call, so the `aggrs` list is empty
		// and the selected expression is transformed normally
		{"f(*) FROM x [RANGE 1 TUPLES]", "",
			funcAppAST{"f", []FlatExpression{wildcardAST{}}},
			nil},

		// there is an aggregate call `count(a)`, so it is referenced from
		// the expression list and appears in the `aggrs` list
		{"count(a) FROM x [RANGE 1 TUPLES]", "",
			funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		// there is an aggregate call `count(*)`, so it is referenced from
		// the expression list and a constant appears in the `aggrs` list
		{"count(*) FROM x [RANGE 1 TUPLES]", "",
			funcAppAST{"count", []FlatExpression{aggInputRef{"g_356a192b"}}},
			map[string]FlatExpression{
				"g_356a192b": numericLiteral{1},
			}},

		// there is an aggregate call `udaf(*, 1)`, so it is referenced from
		// the expression list and a constant appears in the `aggrs` list
		{"udaf(*, 1) FROM x [RANGE 1 TUPLES]", "",
			funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_df58248c"}, numericLiteral{1}}},
			map[string]FlatExpression{
				"g_df58248c": wildcardAST{},
			}},

		// there is an aggregate call `count(a)`, so it is referenced from
		// the expression list and appears in the `aggrs` list
		{"a + count(a) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			binaryOpAST{parser.Plus,
				rowValue{"x", "a"}, funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}}},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		// there is an aggregate call `udaf(a+1)`, so it is referenced from
		// the expression list and appears in the `aggrs` list
		{"a + udaf(a + 1) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			binaryOpAST{parser.Plus,
				rowValue{"x", "a"},
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2d5e5764"}}}},
			map[string]FlatExpression{
				"g_2d5e5764": binaryOpAST{parser.Plus, rowValue{"x", "a"}, numericLiteral{1}},
			}},

		// there are two aggregate calls, so both are referenced from the
		// expression list and there are two entries in the `aggrs` list
		{"udaf(a + 1) + g(count(a)) FROM x [RANGE 1 TUPLES]", "",
			binaryOpAST{parser.Plus,
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2d5e5764"}}},
				funcAppAST{"g", []FlatExpression{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				}},
			},
			map[string]FlatExpression{
				"g_2d5e5764": binaryOpAST{parser.Plus,
					rowValue{"x", "a"},
					numericLiteral{1},
				},
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		{"[udaf(a + 1), 3, g(count(a))] FROM x [RANGE 1 TUPLES]", "",
			arrayAST{[]FlatExpression{
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2d5e5764"}}},
				numericLiteral{3},
				funcAppAST{"g", []FlatExpression{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				}},
			}},
			map[string]FlatExpression{
				"g_2d5e5764": binaryOpAST{parser.Plus,
					rowValue{"x", "a"},
					numericLiteral{1},
				},
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		{"{'udaf': udaf(a + 1), '3': 3, 'g': g(count(a))} FROM x [RANGE 1 TUPLES]", "",
			mapAST{[]keyValuePair{
				{"udaf", funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2d5e5764"}}}},
				{"3", numericLiteral{3}},
				{"g", funcAppAST{"g", []FlatExpression{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				}}},
			}},
			map[string]FlatExpression{
				"g_2d5e5764": binaryOpAST{parser.Plus,
					rowValue{"x", "a"},
					numericLiteral{1},
				},
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		// there are two aggregate calls, but they use the same value,
		// so the `aggrs` list contains only one entry
		{"count(a) + g(count(a)) FROM x [RANGE 1 TUPLES]", "",
			binaryOpAST{parser.Plus,
				funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				funcAppAST{"g", []FlatExpression{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				}},
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		// order by a value that is already in the aggregate variables
		{"count(a ORDER BY a ASC) FROM x [RANGE 1 TUPLES]", "",
			aggregateInputSorter{
				funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				[]sortExpression{sortExpression{aggInputRef{"g_f12cd6bc"}, true}},
				"ccd0ef22",
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		// order by a value that is not in the aggregate variables
		{"count(a ORDER BY b DESC) FROM x [RANGE 1 TUPLES]", "",
			aggregateInputSorter{
				funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				[]sortExpression{sortExpression{aggInputRef{"g_77d2dd39"}, false}},
				"d7196f56",
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
				"g_77d2dd39": rowValue{"x", "b"},
			}},

		// use two different sorting orders
		{"count(a ORDER BY b DESC) + count(a ORDER BY b ASC) FROM x [RANGE 1 TUPLES]", "",
			binaryOpAST{parser.Plus,
				aggregateInputSorter{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
					[]sortExpression{sortExpression{aggInputRef{"g_77d2dd39"}, false}},
					"d7196f56",
				},
				aggregateInputSorter{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
					[]sortExpression{sortExpression{aggInputRef{"g_77d2dd39"}, true}},
					"cd35e18d",
				},
			},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
				"g_77d2dd39": rowValue{"x", "b"},
			}},

		{"count(udaf(a)) FROM x [RANGE 1 TUPLES]",
			"aggregate functions cannot be nested", nil, nil},

		{"a FROM x [RANGE 1 TUPLES] WHERE count(a) = 1",
			"aggregates not allowed in WHERE clause", nil, nil},

		{"a FROM x [RANGE 1 TUPLES] GROUP BY count(a)",
			"aggregates not allowed in GROUP BY clause", nil, nil},

		{"a FROM x [RANGE 1 TUPLES] GROUP BY a + 2",
			"grouping by expressions is not supported yet", nil, nil},

		// various grouping checks
		{"a FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			rowValue{"x", "a"},
			nil},

		{"count(a) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
			map[string]FlatExpression{
				"g_f12cd6bc": rowValue{"x", "a"},
			}},

		{"count(b) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			funcAppAST{"count", []FlatExpression{aggInputRef{"g_77d2dd39"}}},
			map[string]FlatExpression{
				"g_77d2dd39": rowValue{"x", "b"},
			}},

		{"count(b), a FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			funcAppAST{"count", []FlatExpression{aggInputRef{"g_77d2dd39"}}}, // just the first one
			map[string]FlatExpression{
				"g_77d2dd39": rowValue{"x", "b"},
			}},

		{"count(b), a, c FROM x [RANGE 1 TUPLES] GROUP BY a, c", "",
			funcAppAST{"count", []FlatExpression{aggInputRef{"g_77d2dd39"}}}, // just the first one
			map[string]FlatExpression{
				"g_77d2dd39": rowValue{"x", "b"},
			}},

		{"udaf(x, a) FROM x [RANGE 1 TUPLES] GROUP BY b",
			"column \"x:a\" must appear in the GROUP BY clause or be used in an aggregate function", nil, nil},

		{"a FROM x [RANGE 1 TUPLES] GROUP BY b",
			"column \"x:a\" must appear in the GROUP BY clause or be used in an aggregate function", nil, nil},

		{"a FROM x [RANGE 1 TUPLES] AS y GROUP BY b",
			"column \"y:a\" must appear in the GROUP BY clause or be used in an aggregate function", nil, nil},

		{"a, count(b) FROM x [RANGE 1 TUPLES]",
			"column \"x:a\" must appear in the GROUP BY clause or be used in an aggregate function", nil, nil},

		{"a + count(b) FROM x [RANGE 1 TUPLES]",
			"column \"x:a\" must appear in the GROUP BY clause or be used in an aggregate function", nil, nil},
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

func TestVolatileAggregateChecker(t *testing.T) {
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))

	dummyFun := &dummyAggregate{}
	toString := udf.UnaryFunc(func(ctx *core.Context, v data.Value) (data.Value, error) {
		return data.String(v.String()), nil
	})

	reg.Register("udaf", dummyFun)
	reg.Register("f", toString)

	// For volatile expressions as aggregate function parameters,
	// we must generate unique reference strings.

	testCases := []struct {
		bql           string
		expectedError string
		exprs         []FlatExpression
		aggrs         []map[string]FlatExpression
	}{
		// for comparison: two UDAFs referencing the same immutable
		// expression use the same reference string
		{"count(a) + udaf(a) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				binaryOpAST{parser.Plus,
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
					funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_f12cd6bc"}}}},
			},
			[]map[string]FlatExpression{{
				"g_f12cd6bc": rowValue{"x", "a"},
			}}},

		// one immutable and one volatile parameter
		{"count(a) + udaf(f(a)) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				binaryOpAST{parser.Plus,
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
					funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_1"}}}},
			},
			[]map[string]FlatExpression{{
				"g_f12cd6bc":   rowValue{"x", "a"},
				"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
			}}},

		// one volatile and one immutable parameter (order reversed)
		{"udaf(f(a)) + count(a) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				binaryOpAST{parser.Plus,
					funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_0"}}},
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}}},
			},
			[]map[string]FlatExpression{{
				"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				"g_f12cd6bc":   rowValue{"x", "a"},
			}}},

		// two UDAFs referencing the same volatile expression
		// use different reference strings
		{"count(f(a)) + udaf(f(a)) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				binaryOpAST{parser.Plus,
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_2523c3a2_0"}}},
					funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_1"}}}},
			},
			[]map[string]FlatExpression{{
				"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
			}}},

		// one UDAF referencing the same volatile expression in parameter
		// and GROUP BY uses different reference strings
		{"count(f(a) ORDER BY f(a)) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				aggregateInputSorter{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_2523c3a2_0"}}},
					[]sortExpression{sortExpression{aggInputRef{"g_2523c3a2_1"}, true}},
					"cf2e24d7",
				},
			},
			[]map[string]FlatExpression{{
				"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
			}}},

		// three UDAFs referencing the same volatile expression
		// from different columns all use different reference strings
		{"count(f(a)) + udaf(f(a)), udaf(f(a)) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				binaryOpAST{parser.Plus,
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_2523c3a2_0"}}},
					funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_1"}}}},
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_2"}}},
			},
			[]map[string]FlatExpression{
				{
					"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				},
				{
					"g_2523c3a2_2": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				}}},

		// three UDAFs referencing the same volatile expression
		// from different columns and on different parameter positions
		// all use different reference strings
		{"udaf(f(a), a, f(a)), count(f(a)) + udaf(f(a)) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_0"},
					rowValue{"x", "a"}, aggInputRef{"g_2523c3a2_1"}}},
				binaryOpAST{parser.Plus,
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_2523c3a2_2"}}},
					funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_3"}}}},
			},
			[]map[string]FlatExpression{
				{
					"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				},
				{
					"g_2523c3a2_2": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_3": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				}}},

		{"udaf(f(a), a, f(a)), [count(f(a)), udaf(f(a))] FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_0"},
					rowValue{"x", "a"}, aggInputRef{"g_2523c3a2_1"}}},
				arrayAST{[]FlatExpression{
					funcAppAST{"count", []FlatExpression{aggInputRef{"g_2523c3a2_2"}}},
					funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_3"}}},
				}},
			},
			[]map[string]FlatExpression{
				{
					"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				},
				{
					"g_2523c3a2_2": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_3": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				}}},

		{"udaf(f(a), a, f(a)), {'c': count(f(a)), 'u': udaf(f(a))} FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_0"},
					rowValue{"x", "a"}, aggInputRef{"g_2523c3a2_1"}}},
				mapAST{[]keyValuePair{
					{"c", funcAppAST{"count", []FlatExpression{aggInputRef{"g_2523c3a2_2"}}}},
					{"u", funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_3"}}}},
				}},
			},
			[]map[string]FlatExpression{
				{
					"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				},
				{
					"g_2523c3a2_2": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_3": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				}}},

		// very weird mixed combination
		{"udaf(f(a), a, f(a)), count(a), udaf(f(a)) FROM x [RANGE 1 TUPLES] GROUP BY a", "",
			[]FlatExpression{
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_0"},
					rowValue{"x", "a"}, aggInputRef{"g_2523c3a2_1"}}},
				funcAppAST{"count", []FlatExpression{aggInputRef{"g_f12cd6bc"}}},
				funcAppAST{"udaf", []FlatExpression{aggInputRef{"g_2523c3a2_3"}}},
			},
			[]map[string]FlatExpression{
				{
					"g_2523c3a2_0": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
					"g_2523c3a2_1": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				},
				{
					"g_f12cd6bc": rowValue{"x", "a"},
				},
				{
					"g_2523c3a2_3": funcAppAST{"f", []FlatExpression{rowValue{"x", "a"}}},
				}}},
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
						for i, proj := range logPlan.Projections {
							So(proj.expr, ShouldResemble, testCase.exprs[i])
							So(proj.aggrInputs, ShouldResemble, testCase.aggrs[i])
						}
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

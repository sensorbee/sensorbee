package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestFlatExpressionConverter(t *testing.T) {
	testCases := map[string]struct {
		e FlatExpression
		v VolatilityType
		r []RowValue
	}{
		// Base Expressions
		"true":  {BoolLiteral{true}, Immutable, nil},
		"NULL":  {NullLiteral{}, Immutable, nil},
		"a":     {RowValue{"", "a"}, Immutable, []RowValue{{"", "a"}}},
		"ts()":  {RowMeta{"", parser.TimestampMeta}, Immutable, nil},
		"now()": {StmtMeta{parser.NowMeta}, Stable, nil},
		"2":     {NumericLiteral{2}, Immutable, nil},
		"1.2":   {FloatLiteral{1.2}, Immutable, nil},
		`'bql'`: {StringLiteral{"bql"}, Immutable, nil},
		"*":     {WildcardAST{}, Stable, nil},
		"x:*":   {WildcardAST{"x"}, Stable, nil},
		// Type Cast
		"CAST(2 AS FLOAT)": {TypeCastAST{NumericLiteral{2}, parser.Float}, Immutable, nil},
		// Function Application
		"f(a)": {FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}, Volatile, []RowValue{{"", "a"}}},
		// Aggregate Function Application
		"count(a)": {FuncAppAST{parser.FuncName("count"),
			[]FlatExpression{AggInputRef{"_a4839edb"}}}, Volatile, nil},
		// Arrays
		"[]":  {ArrayAST{[]FlatExpression{}}, Immutable, nil},
		"[2]": {ArrayAST{[]FlatExpression{NumericLiteral{2}}}, Immutable, nil},
		"[a, now()]": {ArrayAST{[]FlatExpression{RowValue{"", "a"},
			StmtMeta{parser.NowMeta}}}, Stable, []RowValue{{"", "a"}}},
		"[f(a), true]": {ArrayAST{[]FlatExpression{FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}, BoolLiteral{true}}}, Volatile, []RowValue{{"", "a"}}},
		// Maps
		"{}":          {MapAST{[]KeyValuePair{}}, Immutable, nil},
		"{'hoge': 2}": {MapAST{[]KeyValuePair{{"hoge", NumericLiteral{2}}}}, Immutable, nil},
		"{'a':a, 'now':now()}": {MapAST{[]KeyValuePair{{"a", RowValue{"", "a"}},
			{"now", StmtMeta{parser.NowMeta}}}}, Stable, []RowValue{{"", "a"}}},
		"{'f':f(a),'b':true}": {MapAST{[]KeyValuePair{{"f", FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}}, {"b", BoolLiteral{true}}}}, Volatile, []RowValue{{"", "a"}}},
		// Composed Expressions
		"a OR 2":    {BinaryOpAST{parser.Or, RowValue{"", "a"}, NumericLiteral{2}}, Immutable, []RowValue{{"", "a"}}},
		"a IS NULL": {BinaryOpAST{parser.Is, RowValue{"", "a"}, NullLiteral{}}, Immutable, []RowValue{{"", "a"}}},
		"NOT a":     {UnaryOpAST{parser.Not, RowValue{"", "a"}}, Immutable, []RowValue{{"", "a"}}},
		"NOT f(a)": {UnaryOpAST{parser.Not, FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}}, Volatile, []RowValue{{"", "a"}}},
		// Comparisons
		"a = 2": {BinaryOpAST{parser.Equal, RowValue{"", "a"}, NumericLiteral{2}}, Immutable, []RowValue{{"", "a"}}},
		"f(a) = 2": {BinaryOpAST{parser.Equal, FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}, NumericLiteral{2}}, Volatile, []RowValue{{"", "a"}}},
	}

	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	toString := udf.UnaryFunc(func(ctx *core.Context, v data.Value) (data.Value, error) {
		return data.String(v.String()), nil
	})
	reg.Register("f", toString)

	Convey("Given a BQL parser", t, func() {
		p := parser.NewBQLParser()

		for input, expected := range testCases {
			// avoid closure over loop variables
			input, expected := input, expected

			Convey(fmt.Sprintf("When parsing %s", input), func() {
				stmt := "SELECT ISTREAM " + input
				result, _, err := p.ParseStmt(stmt)

				Convey(fmt.Sprintf("Then the result should be %v", expected), func() {
					if expected.e == nil {
						So(err, ShouldNotBeNil)
					} else {
						So(err, ShouldBeNil)
						// check we got a proper SELECT statement
						So(result, ShouldHaveSameTypeAs, parser.SelectStmt{})
						selectStmt := result.(parser.SelectStmt)
						So(len(selectStmt.Projections), ShouldBeGreaterThan, 0)
						// convert it to FlatExpression
						actual, _, err := ParserExprToMaybeAggregate(selectStmt.Projections[0], 0, reg)
						So(err, ShouldBeNil)
						// compare it against our expectation
						So(actual, ShouldResemble, expected.e)
						So(actual.Volatility(), ShouldEqual, expected.v)
						So(actual.Columns(), ShouldResemble, expected.r)
					}
				})
			})
		}
	})
}

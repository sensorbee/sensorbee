package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
)

func TestFlatExpressionConverter(t *testing.T) {
	testCases := map[string]struct {
		e FlatExpression
		v VolatilityType
		r []rowValue
	}{
		// Base Expressions
		"true":  {boolLiteral{true}, Immutable, nil},
		"NULL":  {nullLiteral{}, Immutable, nil},
		"a":     {rowValue{"", "a"}, Immutable, []rowValue{{"", "a"}}},
		"ts()":  {rowMeta{"", parser.TimestampMeta}, Immutable, nil},
		"now()": {stmtMeta{parser.NowMeta}, Stable, nil},
		"2":     {numericLiteral{2}, Immutable, nil},
		"1.2":   {floatLiteral{1.2}, Immutable, nil},
		`"bql"`: {stringLiteral{"bql"}, Immutable, nil},
		"*":     {wildcardAST{}, Stable, nil},
		"x:*":   {wildcardAST{"x"}, Stable, nil},
		// Type Cast
		"CAST(2 AS FLOAT)": {typeCastAST{numericLiteral{2}, parser.Float}, Immutable, nil},
		// Function Application
		"f(a)": {funcAppAST{parser.FuncName("f"),
			[]FlatExpression{rowValue{"", "a"}}}, Volatile, []rowValue{{"", "a"}}},
		"f(x:*)": {funcAppAST{parser.FuncName("f"),
			[]FlatExpression{wildcardAST{"x"}}}, Volatile, nil},
		// Aggregate Function Application
		"count(a)": {funcAppAST{parser.FuncName("count"),
			[]FlatExpression{aggInputRef{"g_a4839edb"}}}, Volatile, nil},
		// Arrays
		"[]":  {arrayAST{[]FlatExpression{}}, Immutable, nil},
		"[2]": {arrayAST{[]FlatExpression{numericLiteral{2}}}, Immutable, nil},
		"[a, now()]": {arrayAST{[]FlatExpression{rowValue{"", "a"},
			stmtMeta{parser.NowMeta}}}, Stable, []rowValue{{"", "a"}}},
		"[f(a), true]": {arrayAST{[]FlatExpression{funcAppAST{parser.FuncName("f"),
			[]FlatExpression{rowValue{"", "a"}}}, boolLiteral{true}}}, Volatile, []rowValue{{"", "a"}}},
		// Maps
		"{}":          {mapAST{[]keyValuePair{}}, Immutable, nil},
		`{"hoge": 2}`: {mapAST{[]keyValuePair{{"hoge", numericLiteral{2}}}}, Immutable, nil},
		`{"a": *}`:    {mapAST{[]keyValuePair{{"a", wildcardAST{}}}}, Stable, nil},
		`{"a":a, "now":now()}`: {mapAST{[]keyValuePair{{"a", rowValue{"", "a"}},
			{"now", stmtMeta{parser.NowMeta}}}}, Stable, []rowValue{{"", "a"}}},
		`{"f":f(a),"b":true}`: {mapAST{[]keyValuePair{{"f", funcAppAST{parser.FuncName("f"),
			[]FlatExpression{rowValue{"", "a"}}}}, {"b", boolLiteral{true}}}}, Volatile, []rowValue{{"", "a"}}},
		// CASE expressions
		"CASE a WHEN 2 THEN 3 END":            {caseAST{rowValue{"", "a"}, []whenThenPair{{numericLiteral{2}, numericLiteral{3}}}, nullLiteral{}}, Immutable, nil},
		"CASE WHEN true THEN 3 END":           {caseAST{boolLiteral{true}, []whenThenPair{{boolLiteral{true}, numericLiteral{3}}}, nullLiteral{}}, Immutable, nil},
		"CASE WHEN false THEN 3 ELSE 6 END":   {caseAST{boolLiteral{true}, []whenThenPair{{boolLiteral{false}, numericLiteral{3}}}, numericLiteral{6}}, Immutable, nil},
		"CASE now() WHEN 2 THEN 3 END":        {caseAST{stmtMeta{parser.NowMeta}, []whenThenPair{{numericLiteral{2}, numericLiteral{3}}}, nullLiteral{}}, Stable, nil},
		"CASE a WHEN now() THEN 3 END":        {caseAST{rowValue{"", "a"}, []whenThenPair{{stmtMeta{parser.NowMeta}, numericLiteral{3}}}, nullLiteral{}}, Stable, nil},
		"CASE a WHEN 2 THEN now() END":        {caseAST{rowValue{"", "a"}, []whenThenPair{{numericLiteral{2}, stmtMeta{parser.NowMeta}}}, nullLiteral{}}, Stable, nil},
		"CASE a WHEN 2 THEN 3 ELSE now() END": {caseAST{rowValue{"", "a"}, []whenThenPair{{numericLiteral{2}, numericLiteral{3}}}, stmtMeta{parser.NowMeta}}, Stable, nil},
		// Composed Expressions
		"a OR 2":    {binaryOpAST{parser.Or, rowValue{"", "a"}, numericLiteral{2}}, Immutable, []rowValue{{"", "a"}}},
		"a IS NULL": {binaryOpAST{parser.Is, rowValue{"", "a"}, nullLiteral{}}, Immutable, []rowValue{{"", "a"}}},
		"NOT a":     {unaryOpAST{parser.Not, rowValue{"", "a"}}, Immutable, []rowValue{{"", "a"}}},
		"NOT f(a)": {unaryOpAST{parser.Not, funcAppAST{parser.FuncName("f"),
			[]FlatExpression{rowValue{"", "a"}}}}, Volatile, []rowValue{{"", "a"}}},
		// Comparisons
		"a = 2": {binaryOpAST{parser.Equal, rowValue{"", "a"}, numericLiteral{2}}, Immutable, []rowValue{{"", "a"}}},
		"f(a) = 2": {binaryOpAST{parser.Equal, funcAppAST{parser.FuncName("f"),
			[]FlatExpression{rowValue{"", "a"}}}, numericLiteral{2}}, Volatile, []rowValue{{"", "a"}}},
	}

	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	toString := udf.UnaryFunc(func(ctx *core.Context, v data.Value) (data.Value, error) {
		return data.String(v.String()), nil
	})
	reg.Register("f", toString)

	Convey("Given a BQL parser", t, func() {
		p := parser.New()

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

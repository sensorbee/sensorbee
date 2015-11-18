package parser

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExpressionParser(t *testing.T) {
	testCases := map[string]struct {
		expr []Expression
		str  string
	}{
		/// Base Expressions
		// BooleanLiteral
		"true":  {[]Expression{BoolLiteral{true}}, "TRUE"},
		"TRUE":  {[]Expression{BoolLiteral{true}}, "TRUE"},
		"false": {[]Expression{BoolLiteral{false}}, "FALSE"},
		// NullLiteral
		"NULL": {[]Expression{NullLiteral{}}, "NULL"},
		// Function Application
		"f()": {[]Expression{FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{}}, nil}}, "f()"},
		"now()": {[]Expression{FuncAppAST{FuncName("now"),
			ExpressionsAST{[]Expression{}}, nil}}, "now()"},
		"f(a)": {[]Expression{FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{RowValue{"", "a"}}}, nil}}, "f(a)"},
		"f(*)": {[]Expression{FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{Wildcard{}}}, nil}}, "f(*)"},
		"f(x:*)": {[]Expression{FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{Wildcard{"x"}}}, nil}}, "f(x:*)"},
		"f(x:* ORDER BY a)": {[]Expression{FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{Wildcard{"x"}}},
			[]SortedExpressionAST{{RowValue{"", "a"}, UnspecifiedKeyword}}}}, "f(x:* ORDER BY a)"},
		"f(a ORDER BY a DESC, b, c ASC)": {[]Expression{FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{RowValue{"", "a"}}},
			[]SortedExpressionAST{{RowValue{"", "a"}, No}, {RowValue{"", "b"}, UnspecifiedKeyword}, {RowValue{"", "c"}, Yes}}}}, "f(a ORDER BY a DESC, b, c ASC)"},
		"count(a ORDER BY count(b))": {[]Expression{FuncAppAST{FuncName("count"),
			ExpressionsAST{[]Expression{RowValue{"", "a"}}},
			[]SortedExpressionAST{{FuncAppAST{FuncName("count"),
				ExpressionsAST{[]Expression{RowValue{"", "b"}}}, nil}, UnspecifiedKeyword}}}}, "count(a ORDER BY count(b))"},
		"f(2.1, 'a')": {[]Expression{FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{FloatLiteral{2.1}, StringLiteral{"a"}}}, nil}}, "f(2.1, 'a')"},
		// Type Cast
		"CAST(2.1 AS BOOL)":    {[]Expression{TypeCastAST{FloatLiteral{2.1}, Bool}}, "CAST(2.1 AS BOOL)"},
		"CAST(2.1 AS INT)":     {[]Expression{TypeCastAST{FloatLiteral{2.1}, Int}}, "CAST(2.1 AS INT)"},
		"CAST(a*2 AS FLOAT)":   {[]Expression{TypeCastAST{BinaryOpAST{Multiply, RowValue{"", "a"}, NumericLiteral{2}}, Float}}, "CAST(a * 2 AS FLOAT)"},
		"CAST(2.1 AS STRING)":  {[]Expression{TypeCastAST{FloatLiteral{2.1}, String}}, "CAST(2.1 AS STRING)"},
		"CAST('hoge' AS BLOB)": {[]Expression{TypeCastAST{StringLiteral{"hoge"}, Blob}}, "CAST('hoge' AS BLOB)"},
		"CAST(0 AS TIMESTAMP)": {[]Expression{TypeCastAST{NumericLiteral{0}, Timestamp}}, "CAST(0 AS TIMESTAMP)"},
		"CAST(2.1 AS ARRAY)":   {[]Expression{TypeCastAST{FloatLiteral{2.1}, Array}}, "CAST(2.1 AS ARRAY)"},
		"CAST('a' AS MAP)":     {[]Expression{TypeCastAST{StringLiteral{"a"}, Map}}, "CAST('a' AS MAP)"},
		"2.1::INT":             {[]Expression{TypeCastAST{FloatLiteral{2.1}, Int}}, "CAST(2.1 AS INT)"},
		"int::STRING":          {[]Expression{TypeCastAST{RowValue{"", "int"}, String}}, "int::STRING"},
		"x:int::STRING":        {[]Expression{TypeCastAST{RowValue{"x", "int"}, String}}, "x:int::STRING"},
		"ts()::STRING":         {[]Expression{TypeCastAST{RowMeta{"", TimestampMeta}, String}}, "ts()::STRING"},
		"tab:ts()::STRING":     {[]Expression{TypeCastAST{RowMeta{"tab", TimestampMeta}, String}}, "tab:ts()::STRING"},
		// RowValue
		"a":         {[]Expression{RowValue{"", "a"}}, "a"},
		"-a":        {[]Expression{UnaryOpAST{UnaryMinus, RowValue{"", "a"}}}, "-a"},
		"tab:a":     {[]Expression{RowValue{"tab", "a"}}, "tab:a"},
		"ts()":      {[]Expression{RowMeta{"", TimestampMeta}}, "ts()"},
		"tab:ts()":  {[]Expression{RowMeta{"tab", TimestampMeta}}, "tab:ts()"},
		"a, b":      {[]Expression{RowValue{"", "a"}, RowValue{"", "b"}}, "a, b"},
		"A":         {[]Expression{RowValue{"", "A"}}, "A"},
		"my_mem_27": {[]Expression{RowValue{"", "my_mem_27"}}, "my_mem_27"},
		/// JSON Path
		"['hoge']":       {[]Expression{RowValue{"", "['hoge']"}}, "['hoge']"},
		"['hoge'][0]..y": {[]Expression{RowValue{"", "['hoge'][0]..y"}}, "['hoge'][0]..y"},
		"['array'][0]":   {[]Expression{RowValue{"", "['array'][0]"}}, "['array'][0]"},
		"['array'][0].x": {[]Expression{RowValue{"", "['array'][0].x"}}, "['array'][0].x"},
		"['array']['x']": {[]Expression{RowValue{"", "['array']['x']"}}, "['array']['x']"},
		"array.x":        {[]Expression{RowValue{"", "array.x"}}, "array.x"},
		// Colon checks
		"array['x::int']": {[]Expression{RowValue{"", "array['x::int']"}}, "array['x::int']"},
		"[':hoge']":       {[]Expression{RowValue{"", "[':hoge']"}}, "[':hoge']"},
		"x:[':hoge']":     {[]Expression{RowValue{"x", "[':hoge']"}}, "x:[':hoge']"},
		"x:['hoge']":      {[]Expression{RowValue{"x", "['hoge']"}}, "x:['hoge']"},
		// Quote checks
		"['ar''ray']['x::int']": {[]Expression{RowValue{"", "['ar''ray']['x::int']"}}, "['ar''ray']['x::int']"},
		// Wildcard
		"*":         {[]Expression{Wildcard{}}, "*"},
		"x:*":       {[]Expression{Wildcard{"x"}}, "x:*"},
		"* IS NULL": {nil, ""}, // the wildcard is not a normal Expression!
		// Array
		"[]":          {[]Expression{ArrayAST{ExpressionsAST{[]Expression{}}}}, "[]"},
		"[2]":         {[]Expression{ArrayAST{ExpressionsAST{[]Expression{NumericLiteral{2}}}}}, "[2]"},
		"[2,]":        {[]Expression{ArrayAST{ExpressionsAST{[]Expression{NumericLiteral{2}}}}}, "[2]"},
		"[a]":         {[]Expression{ArrayAST{ExpressionsAST{[]Expression{RowValue{"", "a"}}}}}, "[a]"},
		"[a,]":        {[]Expression{ArrayAST{ExpressionsAST{[]Expression{RowValue{"", "a"}}}}}, "[a]"},
		"[a,b:*]":     {[]Expression{ArrayAST{ExpressionsAST{[]Expression{RowValue{"", "a"}, Wildcard{"b"}}}}}, "[a, b:*]"},
		"['hoge',]":   {[]Expression{ArrayAST{ExpressionsAST{[]Expression{StringLiteral{"hoge"}}}}}, "['hoge']"},
		"x:['hoge',]": {nil, ""}, // an array takes no stream prefix
		"[a, 2.3]":    {[]Expression{ArrayAST{ExpressionsAST{[]Expression{RowValue{"", "a"}, FloatLiteral{2.3}}}}}, "[a, 2.3]"},
		"[a, 2.3,]":   {[]Expression{ArrayAST{ExpressionsAST{[]Expression{RowValue{"", "a"}, FloatLiteral{2.3}}}}}, "[a, 2.3]"},
		"[['hoge'], hoge.x[0], a]": {[]Expression{ArrayAST{ExpressionsAST{[]Expression{RowValue{"", "['hoge']"},
			RowValue{"", "hoge.x[0]"}, RowValue{"", "a"}}}}}, "[['hoge'], hoge.x[0], a]"},
		"[{'key':hoge:image}]": {[]Expression{ArrayAST{ExpressionsAST{[]Expression{
			MapAST{[]KeyValuePairAST{{"key", RowValue{"hoge", "image"}}}}}}}}, "[{'key':hoge:image}]"},
		// Map
		"{}":                 {[]Expression{MapAST{[]KeyValuePairAST{}}}, "{}"},
		"{'hoge':2}":         {[]Expression{MapAST{[]KeyValuePairAST{{"hoge", NumericLiteral{2}}}}}, "{'hoge':2}"},
		"{'key':hoge:image}": {[]Expression{MapAST{[]KeyValuePairAST{{"key", RowValue{"hoge", "image"}}}}}, "{'key':hoge:image}"},
		"{'foo':x:a, 'bar':{'a':[2]}}": {[]Expression{MapAST{[]KeyValuePairAST{
			{"foo", RowValue{"x", "a"}},
			{"bar", MapAST{[]KeyValuePairAST{
				{"a", ArrayAST{ExpressionsAST{[]Expression{NumericLiteral{2}}}}},
			}}},
		}}}, "{'foo':x:a, 'bar':{'a':[2]}}"},
		"{'a': a:*, 'b': b:*}": {[]Expression{MapAST{[]KeyValuePairAST{
			{"a", Wildcard{"a"}}, {"b", Wildcard{"b"}}}}}, "{'a':a:*, 'b':b:*}"},
		// CASE
		"CASE a END":                         {nil, ""}, // WHEN-THEN is mandatory
		"CASE a WHEN 2 THEN 3 END":           {[]Expression{ExpressionCaseAST{RowValue{"", "a"}, []WhenThenPairAST{{NumericLiteral{2}, NumericLiteral{3}}}, nil}}, "CASE a WHEN 2 THEN 3 END"},
		"CASE when WHEN 2 THEN 3 ELSE 6 END": {[]Expression{ExpressionCaseAST{RowValue{"", "when"}, []WhenThenPairAST{{NumericLiteral{2}, NumericLiteral{3}}}, NumericLiteral{6}}}, "CASE when WHEN 2 THEN 3 ELSE 6 END"},
		// NumericLiteral
		"2":    {[]Expression{NumericLiteral{2}}, "2"},
		"-2":   {[]Expression{UnaryOpAST{UnaryMinus, NumericLiteral{2}}}, "-2"},
		"- -2": {[]Expression{UnaryOpAST{UnaryMinus, NumericLiteral{-2}}}, "- -2"}, // like PostgreSQL
		"999999999999999999999999999": {nil, ""}, // int64 overflow
		// FloatLiteral
		"1.2":   {[]Expression{FloatLiteral{1.2}}, "1.2"},
		"-3.14": {[]Expression{UnaryOpAST{UnaryMinus, FloatLiteral{3.14}}}, "-3.14"},
		// StringLiteral
		`'bql'`:      {[]Expression{StringLiteral{"bql"}}, "'bql'"},
		`''`:         {[]Expression{StringLiteral{""}}, `''`},
		`'Peter's'`:  {nil, ""},
		`'Peter''s'`: {[]Expression{StringLiteral{"Peter's"}}, `'Peter''s'`},
		`'日本語'`:      {[]Expression{StringLiteral{"日本語"}}, `'日本語'`},
		// Alias
		"1.2 AS x": {[]Expression{AliasAST{FloatLiteral{1.2}, "x"}}, "1.2 AS x"},
		"b AS *":   {[]Expression{AliasAST{RowValue{"", "b"}, "*"}}, "b AS *"},
		/// Composed Expressions
		// OR
		"a OR 2": {[]Expression{BinaryOpAST{Or, RowValue{"", "a"}, NumericLiteral{2}}}, "a OR 2"},
		// AND
		"a AND 2": {[]Expression{BinaryOpAST{And, RowValue{"", "a"}, NumericLiteral{2}}}, "a AND 2"},
		// NOT
		"NOT a": {[]Expression{UnaryOpAST{Not, RowValue{"", "a"}}}, "NOT a"},
		// Comparisons
		"a = 2":  {[]Expression{BinaryOpAST{Equal, RowValue{"", "a"}, NumericLiteral{2}}}, "a = 2"},
		"a < 2":  {[]Expression{BinaryOpAST{Less, RowValue{"", "a"}, NumericLiteral{2}}}, "a < 2"},
		"a <= 2": {[]Expression{BinaryOpAST{LessOrEqual, RowValue{"", "a"}, NumericLiteral{2}}}, "a <= 2"},
		"a > 2":  {[]Expression{BinaryOpAST{Greater, RowValue{"", "a"}, NumericLiteral{2}}}, "a > 2"},
		"a >= 2": {[]Expression{BinaryOpAST{GreaterOrEqual, RowValue{"", "a"}, NumericLiteral{2}}}, "a >= 2"},
		"a != 2": {[]Expression{BinaryOpAST{NotEqual, RowValue{"", "a"}, NumericLiteral{2}}}, "a != 2"},
		"a <> 2": {[]Expression{BinaryOpAST{NotEqual, RowValue{"", "a"}, NumericLiteral{2}}}, "a != 2"},
		// Other operators
		"a || 2": {[]Expression{BinaryOpAST{Concat, RowValue{"", "a"}, NumericLiteral{2}}}, "a || 2"},
		// IS Expressions
		"a IS NULL":     {[]Expression{BinaryOpAST{Is, RowValue{"", "a"}, NullLiteral{}}}, "a IS NULL"},
		"a IS NOT NULL": {[]Expression{BinaryOpAST{IsNot, RowValue{"", "a"}, NullLiteral{}}}, "a IS NOT NULL"},
		// Plus/Minus Terms
		"a + 2": {[]Expression{BinaryOpAST{Plus, RowValue{"", "a"}, NumericLiteral{2}}}, "a + 2"},
		"a - 2": {[]Expression{BinaryOpAST{Minus, RowValue{"", "a"}, NumericLiteral{2}}}, "a - 2"},
		"a + -2": {[]Expression{BinaryOpAST{Plus, RowValue{"", "a"},
			UnaryOpAST{UnaryMinus, NumericLiteral{2}}}}, "a + -2"},
		// Product/Division Terms
		"a * 2": {[]Expression{BinaryOpAST{Multiply, RowValue{"", "a"}, NumericLiteral{2}}}, "a * 2"},
		"-a * -2": {[]Expression{BinaryOpAST{Multiply, UnaryOpAST{UnaryMinus, RowValue{"", "a"}},
			UnaryOpAST{UnaryMinus, NumericLiteral{2}}}}, "-a * -2"},
		"a / 2": {[]Expression{BinaryOpAST{Divide, RowValue{"", "a"}, NumericLiteral{2}}}, "a / 2"},
		"a % 2": {[]Expression{BinaryOpAST{Modulo, RowValue{"", "a"}, NumericLiteral{2}}}, "a % 2"},
		/// Operator Precedence
		"a AND b OR 2": {[]Expression{BinaryOpAST{Or,
			BinaryOpAST{And, RowValue{"", "a"}, RowValue{"", "b"}},
			NumericLiteral{2}}}, "(a AND b) OR 2"},
		"2 OR a AND b": {[]Expression{BinaryOpAST{Or,
			NumericLiteral{2},
			BinaryOpAST{And, RowValue{"", "a"}, RowValue{"", "b"}}}}, "2 OR (a AND b)"},
		"NOT a AND b": {[]Expression{BinaryOpAST{And,
			UnaryOpAST{Not, RowValue{"", "a"}},
			RowValue{"", "b"}}}, "NOT a AND b"},
		"a AND NOT b": {[]Expression{BinaryOpAST{And,
			RowValue{"", "a"},
			UnaryOpAST{Not, RowValue{"", "b"}}}}, "a AND NOT b"},
		"a * b + 2": {[]Expression{BinaryOpAST{Plus,
			BinaryOpAST{Multiply, RowValue{"", "a"}, RowValue{"", "b"}},
			NumericLiteral{2}}}, "a * b + 2"},
		"2 - a * b": {[]Expression{BinaryOpAST{Minus,
			NumericLiteral{2},
			BinaryOpAST{Multiply, RowValue{"", "a"}, RowValue{"", "b"}}}}, "2 - a * b"},
		"2 OR a IS NULL": {[]Expression{BinaryOpAST{Or,
			NumericLiteral{2},
			BinaryOpAST{Is, RowValue{"", "a"}, NullLiteral{}}}}, "2 OR a IS NULL"},
		"2 + a IS NOT NULL": {[]Expression{BinaryOpAST{IsNot,
			BinaryOpAST{Plus, NumericLiteral{2}, RowValue{"", "a"}},
			NullLiteral{}}}, "2 + a IS NOT NULL"},
		"2 || a IS NULL": {[]Expression{BinaryOpAST{Concat,
			NumericLiteral{2},
			BinaryOpAST{Is, RowValue{"", "a"}, NullLiteral{}}}}, "2 || a IS NULL"},
		"2 || a = 4": {[]Expression{BinaryOpAST{Equal,
			BinaryOpAST{Concat, NumericLiteral{2}, RowValue{"", "a"}},
			NumericLiteral{4}}}, "2 || a = 4"},
		"-2.1::INT": {[]Expression{UnaryOpAST{UnaryMinus,
			TypeCastAST{FloatLiteral{2.1}, Int}}}, "-CAST(2.1 AS INT)"},
		/// Left-Associativity
		"a || '2' || b": {[]Expression{BinaryOpAST{Concat,
			BinaryOpAST{Concat, RowValue{"", "a"}, StringLiteral{"2"}}, RowValue{"", "b"}}}, "(a || '2') || b"},
		"a - 2 - b": {[]Expression{BinaryOpAST{Minus,
			BinaryOpAST{Minus, RowValue{"", "a"}, NumericLiteral{2}}, RowValue{"", "b"}}}, "(a - 2) - b"},
		"a - 2 - b + 4": {[]Expression{BinaryOpAST{Plus,
			BinaryOpAST{Minus,
				BinaryOpAST{Minus, RowValue{"", "a"}, NumericLiteral{2}},
				RowValue{"", "b"}},
			NumericLiteral{4}}}, "((a - 2) - b) + 4"},
		"a * 2 / b": {[]Expression{BinaryOpAST{Divide,
			BinaryOpAST{Multiply, RowValue{"", "a"}, NumericLiteral{2}}, RowValue{"", "b"}}}, "(a * 2) / b"},
		"a OR b OR 2": {[]Expression{BinaryOpAST{Or,
			BinaryOpAST{Or, RowValue{"", "a"}, RowValue{"", "b"}},
			NumericLiteral{2}}}, "(a OR b) OR 2"},
		"1 OR 2 OR 3 AND 4 AND 5 OR 6": {[]Expression{BinaryOpAST{Or,
			BinaryOpAST{Or,
				BinaryOpAST{Or, NumericLiteral{1}, NumericLiteral{2}},
				BinaryOpAST{And,
					BinaryOpAST{And, NumericLiteral{3}, NumericLiteral{4}},
					NumericLiteral{5}}},
			NumericLiteral{6}}}, "((1 OR 2) OR ((3 AND 4) AND 5)) OR 6"},
		/// Overriding Operator Precedence
		"a AND (b OR 2)": {[]Expression{BinaryOpAST{And,
			RowValue{"", "a"},
			BinaryOpAST{Or, RowValue{"", "b"}, NumericLiteral{2}}}}, "a AND (b OR 2)"},
		"(2 OR a) AND b": {[]Expression{BinaryOpAST{And,
			BinaryOpAST{Or, NumericLiteral{2}, RowValue{"", "a"}},
			RowValue{"", "b"}}}, "(2 OR a) AND b"},
		"NOT (a AND b)": {[]Expression{UnaryOpAST{Not,
			BinaryOpAST{And, RowValue{"", "a"}, RowValue{"", "b"}}}}, "NOT (a AND b)"},
		"a * (b + 2)": {[]Expression{BinaryOpAST{Multiply,
			RowValue{"", "a"},
			BinaryOpAST{Plus, RowValue{"", "b"}, NumericLiteral{2}}}}, "a * (b + 2)"},
		"(2 - a) * b": {[]Expression{BinaryOpAST{Multiply,
			BinaryOpAST{Minus, NumericLiteral{2}, RowValue{"", "a"}},
			RowValue{"", "b"}}}, "(2 - a) * b"},
		"(2 - a)::STRING": {[]Expression{TypeCastAST{
			BinaryOpAST{Minus, NumericLiteral{2}, RowValue{"", "a"}},
			String}}, "CAST(2 - a AS STRING)"},
		"a * (2 = b)": {[]Expression{BinaryOpAST{Multiply,
			RowValue{"", "a"}, BinaryOpAST{Equal, NumericLiteral{2}, RowValue{"", "b"}}}}, "a * (2 = b)"},
		"a - (2 - b)": {[]Expression{BinaryOpAST{Minus,
			RowValue{"", "a"}, BinaryOpAST{Minus, NumericLiteral{2}, RowValue{"", "b"}}}}, "a - (2 - b)"},
		"a * (2 / b)": {[]Expression{BinaryOpAST{Multiply,
			RowValue{"", "a"}, BinaryOpAST{Divide, NumericLiteral{2}, RowValue{"", "b"}}}}, "a * (2 / b)"},
		/// Operator Non-Precedence
		// TODO using more than one "same-level" operator does not work
		/*
			"a AND b AND 2": {BinaryOpAST{And,
				BinaryOpAST{And, RowValue{"","a"}, RowValue{"","b"}},
				NumericLiteral{2}}},
			"2 OR a OR b": {BinaryOpAST{Or,
				BinaryOpAST{Or, NumericLiteral{2}, RowValue{"","a"}},
				RowValue{"","b"}}},
			"a / b * 2": {BinaryOpAST{Multiply,
				BinaryOpAST{Divide, RowValue{"","a"}, RowValue{"","b"}},
				NumericLiteral{2}}},
			"a - b + 2": {BinaryOpAST{Plus,
				BinaryOpAST{Minus, RowValue{"","a"}, RowValue{"","b"}},
				NumericLiteral{2}}},
		*/
		/// Complex/Nested Expressions
		"(a+ 1) % 2= 0 OR b < 7.1, c": {
			[]Expression{BinaryOpAST{Or,
				BinaryOpAST{Equal,
					BinaryOpAST{Modulo,
						BinaryOpAST{Plus, RowValue{"", "a"}, NumericLiteral{1}},
						NumericLiteral{2}},
					NumericLiteral{0}},
				BinaryOpAST{Less, RowValue{"", "b"}, FloatLiteral{7.1}}},
				RowValue{"", "c"}}, "(a + 1) % 2 = 0 OR b < 7.1, c"},
		/// Multiple Columns
		"a, 3.1, false,-2": {[]Expression{RowValue{"", "a"}, FloatLiteral{3.1}, BoolLiteral{false}, UnaryOpAST{UnaryMinus, NumericLiteral{2}}}, "a, 3.1, FALSE, -2"},
		`'日本語', 13`:        {[]Expression{StringLiteral{"日本語"}, NumericLiteral{13}}, `'日本語', 13`},
		"concat(a, 'Pi', 3.1), b": {[]Expression{FuncAppAST{FuncName("concat"), ExpressionsAST{
			[]Expression{RowValue{"", "a"}, StringLiteral{"Pi"}, FloatLiteral{3.1}}}, nil},
			RowValue{"", "b"}}, "concat(a, 'Pi', 3.1), b"},
	}

	Convey("Given a BQL parser", t, func() {
		p := New()

		for input, expected := range testCases {
			// avoid closure over loop variables
			input, expected := input, expected

			Convey(fmt.Sprintf("When parsing %s", input), func() {
				stmt := "SELECT ISTREAM " + input
				result, rest, err := p.ParseStmt(stmt)

				Convey(fmt.Sprintf("Then the result should be %v", expected.expr), func() {
					if expected.expr == nil {
						So(err, ShouldNotBeNil)
					} else {
						// check there is no error
						So(err, ShouldBeNil)
						// check we got a proper SELECT statement
						So(result, ShouldHaveSameTypeAs, SelectStmt{})
						selectStmt := result.(SelectStmt)
						// check the statement was parsed completely
						So(rest, ShouldEqual, "")
						// compare it against our expectation
						actual := selectStmt.Projections
						So(len(actual), ShouldEqual, len(expected.expr))
						So(actual, ShouldResemble, expected.expr)

						Convey(fmt.Sprintf("And string() should be \"%s\"", expected.str), func() {
							So(selectStmt.ProjectionsAST.string(), ShouldEqual, expected.str)
						})
					}
				})
			})
		}

	})

}

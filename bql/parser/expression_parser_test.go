package parser

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExpressionParser(t *testing.T) {
	testCases := map[string][]Expression{
		/// Base Expressions
		// BooleanLiteral
		"true":  {BoolLiteral{true}},
		"TRUE":  {BoolLiteral{true}},
		"false": {BoolLiteral{false}},
		// NullLiteral
		"NULL": {NullLiteral{}},
		// Function Application
		"f()": {FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{}}}},
		"now()": {FuncAppAST{FuncName("now"),
			ExpressionsAST{[]Expression{}}}},
		"f(a)": {FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{RowValue{"", "a"}}}}},
		"f(*)": {FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{Wildcard{}}}}},
		"f(x:*)": nil, // forbidden
		"f(2.1, 'a')": {FuncAppAST{FuncName("f"),
			ExpressionsAST{[]Expression{FloatLiteral{2.1}, StringLiteral{"a"}}}}},
		// Type Cast
		"CAST(2.1 AS BOOL)":    {TypeCastAST{FloatLiteral{2.1}, Bool}},
		"CAST(2.1 AS INT)":     {TypeCastAST{FloatLiteral{2.1}, Int}},
		"CAST(a*2 AS FLOAT)":   {TypeCastAST{BinaryOpAST{Multiply, RowValue{"", "a"}, NumericLiteral{2}}, Float}},
		"CAST(2.1 AS STRING)":  {TypeCastAST{FloatLiteral{2.1}, String}},
		"CAST('hoge' AS BLOB)": {TypeCastAST{StringLiteral{"hoge"}, Blob}},
		"CAST(0 AS TIMESTAMP)": {TypeCastAST{NumericLiteral{0}, Timestamp}},
		"CAST(2.1 AS ARRAY)":   {TypeCastAST{FloatLiteral{2.1}, Array}},
		"CAST('a' AS MAP)":     {TypeCastAST{StringLiteral{"a"}, Map}},
		"2.1::INT":             {TypeCastAST{FloatLiteral{2.1}, Int}},
		"int::STRING":          {TypeCastAST{RowValue{"", "int"}, String}},
		"x:int::STRING":        {TypeCastAST{RowValue{"x", "int"}, String}},
		"ts()::STRING":         {TypeCastAST{RowMeta{"", TimestampMeta}, String}},
		"tab:ts()::STRING":     {TypeCastAST{RowMeta{"tab", TimestampMeta}, String}},
		// RowValue
		"a":         {RowValue{"", "a"}},
		"-a":        {UnaryOpAST{UnaryMinus, RowValue{"", "a"}}},
		"tab:a":     {RowValue{"tab", "a"}},
		"ts()":      {RowMeta{"", TimestampMeta}},
		"tab:ts()":  {RowMeta{"tab", TimestampMeta}},
		"a, b":      {RowValue{"", "a"}, RowValue{"", "b"}},
		"A":         {RowValue{"", "A"}},
		"my_mem_27": {RowValue{"", "my_mem_27"}},
		// Wildcard
		"*":   {Wildcard{}},
		"x:*": {Wildcard{"x"}},
		// Array
		"[]":       {ArrayAST{ExpressionsAST{[]Expression{}}}},
		"[2]":      {ArrayAST{ExpressionsAST{[]Expression{NumericLiteral{2}}}}},
		"[a, 2.3]": {ArrayAST{ExpressionsAST{[]Expression{RowValue{"", "a"}, FloatLiteral{2.3}}}}},
		// Map
		"{}":         {MapAST{[]KeyValuePairAST{}}},
		"{'hoge':2}": {MapAST{[]KeyValuePairAST{{"hoge", NumericLiteral{2}}}}},
		"{'foo':x:a, 'bar':{'a':[2]}}": {MapAST{[]KeyValuePairAST{
			{"foo", RowValue{"x", "a"}},
			{"bar", MapAST{[]KeyValuePairAST{
				{"a", ArrayAST{ExpressionsAST{[]Expression{NumericLiteral{2}}}}},
			}}},
		}}},
		// NumericLiteral
		"2":    {NumericLiteral{2}},
		"-2":   {UnaryOpAST{UnaryMinus, NumericLiteral{2}}},
		"- -2": {UnaryOpAST{UnaryMinus, NumericLiteral{-2}}}, // like PostgreSQL
		"999999999999999999999999999": nil, // int64 overflow
		// FloatLiteral
		"1.2":   {FloatLiteral{1.2}},
		"-3.14": {UnaryOpAST{UnaryMinus, FloatLiteral{3.14}}},
		// StringLiteral
		`'bql'`:      {StringLiteral{"bql"}},
		`''`:         {StringLiteral{""}},
		`'Peter's'`:  nil,
		`'Peter''s'`: {StringLiteral{"Peter's"}},
		`'日本語'`:      {StringLiteral{"日本語"}},
		/// Composed Expressions
		// OR
		"a OR 2": {BinaryOpAST{Or, RowValue{"", "a"}, NumericLiteral{2}}},
		// AND
		"a AND 2": {BinaryOpAST{And, RowValue{"", "a"}, NumericLiteral{2}}},
		// NOT
		"NOT a": {UnaryOpAST{Not, RowValue{"", "a"}}},
		// Comparisons
		"a = 2":  {BinaryOpAST{Equal, RowValue{"", "a"}, NumericLiteral{2}}},
		"a < 2":  {BinaryOpAST{Less, RowValue{"", "a"}, NumericLiteral{2}}},
		"a <= 2": {BinaryOpAST{LessOrEqual, RowValue{"", "a"}, NumericLiteral{2}}},
		"a > 2":  {BinaryOpAST{Greater, RowValue{"", "a"}, NumericLiteral{2}}},
		"a >= 2": {BinaryOpAST{GreaterOrEqual, RowValue{"", "a"}, NumericLiteral{2}}},
		"a != 2": {BinaryOpAST{NotEqual, RowValue{"", "a"}, NumericLiteral{2}}},
		"a <> 2": {BinaryOpAST{NotEqual, RowValue{"", "a"}, NumericLiteral{2}}},
		// Other operators
		"a || 2": {BinaryOpAST{Concat, RowValue{"", "a"}, NumericLiteral{2}}},
		// IS Expressions
		"a IS NULL":     {BinaryOpAST{Is, RowValue{"", "a"}, NullLiteral{}}},
		"a IS NOT NULL": {BinaryOpAST{IsNot, RowValue{"", "a"}, NullLiteral{}}},
		// Plus/Minus Terms
		"a + 2": {BinaryOpAST{Plus, RowValue{"", "a"}, NumericLiteral{2}}},
		"a - 2": {BinaryOpAST{Minus, RowValue{"", "a"}, NumericLiteral{2}}},
		"a + -2": {BinaryOpAST{Plus, RowValue{"", "a"},
			UnaryOpAST{UnaryMinus, NumericLiteral{2}}}},
		// Product/Division Terms
		"a * 2": {BinaryOpAST{Multiply, RowValue{"", "a"}, NumericLiteral{2}}},
		"-a * -2": {BinaryOpAST{Multiply, UnaryOpAST{UnaryMinus, RowValue{"", "a"}},
			UnaryOpAST{UnaryMinus, NumericLiteral{2}}}},
		"a / 2": {BinaryOpAST{Divide, RowValue{"", "a"}, NumericLiteral{2}}},
		"a % 2": {BinaryOpAST{Modulo, RowValue{"", "a"}, NumericLiteral{2}}},
		/// Operator Precedence
		"a AND b OR 2": {BinaryOpAST{Or,
			BinaryOpAST{And, RowValue{"", "a"}, RowValue{"", "b"}},
			NumericLiteral{2}}},
		"2 OR a AND b": {BinaryOpAST{Or,
			NumericLiteral{2},
			BinaryOpAST{And, RowValue{"", "a"}, RowValue{"", "b"}}}},
		"NOT a AND b": {BinaryOpAST{And,
			UnaryOpAST{Not, RowValue{"", "a"}},
			RowValue{"", "b"}}},
		"a AND NOT b": {BinaryOpAST{And,
			RowValue{"", "a"},
			UnaryOpAST{Not, RowValue{"", "b"}}}},
		"a * b + 2": {BinaryOpAST{Plus,
			BinaryOpAST{Multiply, RowValue{"", "a"}, RowValue{"", "b"}},
			NumericLiteral{2}}},
		"2 - a * b": {BinaryOpAST{Minus,
			NumericLiteral{2},
			BinaryOpAST{Multiply, RowValue{"", "a"}, RowValue{"", "b"}}}},
		"2 OR a IS NULL": {BinaryOpAST{Or,
			NumericLiteral{2},
			BinaryOpAST{Is, RowValue{"", "a"}, NullLiteral{}}}},
		"2 + a IS NOT NULL": {BinaryOpAST{IsNot,
			BinaryOpAST{Plus, NumericLiteral{2}, RowValue{"", "a"}},
			NullLiteral{}}},
		"2 || a IS NULL": {BinaryOpAST{Concat,
			NumericLiteral{2},
			BinaryOpAST{Is, RowValue{"", "a"}, NullLiteral{}}}},
		"2 || a = 4": {BinaryOpAST{Equal,
			BinaryOpAST{Concat, NumericLiteral{2}, RowValue{"", "a"}},
			NumericLiteral{4}}},
		"-2.1::INT": {UnaryOpAST{UnaryMinus,
			TypeCastAST{FloatLiteral{2.1}, Int}}},
		/// Left-Associativity
		"a || '2' || b": {BinaryOpAST{Concat,
			BinaryOpAST{Concat, RowValue{"", "a"}, StringLiteral{"2"}}, RowValue{"", "b"}}},
		"a - 2 - b": {BinaryOpAST{Minus,
			BinaryOpAST{Minus, RowValue{"", "a"}, NumericLiteral{2}}, RowValue{"", "b"}}},
		"a - 2 - b + 4": {BinaryOpAST{Plus,
			BinaryOpAST{Minus,
				BinaryOpAST{Minus, RowValue{"", "a"}, NumericLiteral{2}},
				RowValue{"", "b"}},
			NumericLiteral{4}}},
		"a * 2 / b": {BinaryOpAST{Divide,
			BinaryOpAST{Multiply, RowValue{"", "a"}, NumericLiteral{2}}, RowValue{"", "b"}}},
		"a OR b OR 2": {BinaryOpAST{Or,
			BinaryOpAST{Or, RowValue{"", "a"}, RowValue{"", "b"}},
			NumericLiteral{2}}},
		"1 OR 2 OR 3 AND 4 AND 5 OR 6": {BinaryOpAST{Or,
			BinaryOpAST{Or,
				BinaryOpAST{Or, NumericLiteral{1}, NumericLiteral{2}},
				BinaryOpAST{And,
					BinaryOpAST{And, NumericLiteral{3}, NumericLiteral{4}},
					NumericLiteral{5}}},
			NumericLiteral{6}}},
		/// Overriding Operator Precedence
		"a AND (b OR 2)": {BinaryOpAST{And,
			RowValue{"", "a"},
			BinaryOpAST{Or, RowValue{"", "b"}, NumericLiteral{2}}}},
		"(2 OR a) AND b": {BinaryOpAST{And,
			BinaryOpAST{Or, NumericLiteral{2}, RowValue{"", "a"}},
			RowValue{"", "b"}}},
		"NOT (a AND b)": {UnaryOpAST{Not,
			BinaryOpAST{And, RowValue{"", "a"}, RowValue{"", "b"}}}},
		"a * (b + 2)": {BinaryOpAST{Multiply,
			RowValue{"", "a"},
			BinaryOpAST{Plus, RowValue{"", "b"}, NumericLiteral{2}}}},
		"(2 - a) * b": {BinaryOpAST{Multiply,
			BinaryOpAST{Minus, NumericLiteral{2}, RowValue{"", "a"}},
			RowValue{"", "b"}}},
		"(2 - a)::STRING": {TypeCastAST{
			BinaryOpAST{Minus, NumericLiteral{2}, RowValue{"", "a"}},
			String}},
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
		"(a + 1) % 2 = 0 OR b < 7.1, c": {
			BinaryOpAST{Or,
				BinaryOpAST{Equal,
					BinaryOpAST{Modulo,
						BinaryOpAST{Plus, RowValue{"", "a"}, NumericLiteral{1}},
						NumericLiteral{2}},
					NumericLiteral{0}},
				BinaryOpAST{Less, RowValue{"", "b"}, FloatLiteral{7.1}}},
			RowValue{"", "c"}},
		/// Multiple Columns
		"a, 3.1, false, -2": {RowValue{"", "a"}, FloatLiteral{3.1}, BoolLiteral{false}, UnaryOpAST{UnaryMinus, NumericLiteral{2}}},
		`'日本語', 13`:         {StringLiteral{"日本語"}, NumericLiteral{13}},
		"concat(a, 'Pi', 3.1), b": {FuncAppAST{FuncName("concat"), ExpressionsAST{
			[]Expression{RowValue{"", "a"}, StringLiteral{"Pi"}, FloatLiteral{3.1}}}},
			RowValue{"", "b"}},
	}

	Convey("Given a BQL parser", t, func() {
		p := NewBQLParser()

		for input, expected := range testCases {
			// avoid closure over loop variables
			input, expected := input, expected

			Convey(fmt.Sprintf("When parsing %s", input), func() {
				stmt := "SELECT ISTREAM " + input
				result, rest, err := p.ParseStmt(stmt)

				Convey(fmt.Sprintf("Then the result should be %v", expected), func() {
					if expected == nil {
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
						So(len(actual), ShouldEqual, len(expected))
						So(actual, ShouldResemble, expected)
					}
				})
			})
		}

	})

}

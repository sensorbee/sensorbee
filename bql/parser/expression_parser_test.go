package parser

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExpressionParser(t *testing.T) {
	testCases := map[string][]interface{}{
		/// Base Expressions
		// BooleanLiteral
		"true":  {BoolLiteral{true}},
		"TRUE":  {BoolLiteral{true}},
		"false": {BoolLiteral{false}},
		// Function Application
		"f(a)": {FuncAppAST{FuncName("f"),
			ExpressionsAST{[]interface{}{ColumnName{"a"}}}}},
		"f(2.1, 'a')": {FuncAppAST{FuncName("f"),
			ExpressionsAST{[]interface{}{FloatLiteral{2.1}, StringLiteral{"a"}}}}},
		// ColumnName
		"a":         {ColumnName{"a"}},
		"a, b":      {ColumnName{"a"}, ColumnName{"b"}},
		"A":         {ColumnName{"A"}},
		"my_mem_27": {ColumnName{"my_mem_27"}},
		// Wildcard
		"*": {Wildcard{}},
		// NumericLiteral
		"2":  {NumericLiteral{2}},
		"-2": {NumericLiteral{-2}},
		"999999999999999999999999999": nil, // int64 overflow
		// FloatLiteral
		"1.2":   {FloatLiteral{1.2}},
		"-3.14": {FloatLiteral{-3.14}},
		// StringLiteral
		`'bql'`:      {StringLiteral{"bql"}},
		`''`:         {StringLiteral{""}},
		`'Peter's'`:  nil,
		`'Peter''s'`: {StringLiteral{"Peter's"}},
		`'日本語'`:      {StringLiteral{"日本語"}},
		/// Composed Expressions
		// OR
		"a OR 2": {BinaryOpAST{Or, ColumnName{"a"}, NumericLiteral{2}}},
		// AND
		"a AND 2": {BinaryOpAST{And, ColumnName{"a"}, NumericLiteral{2}}},
		// Comparisons
		"a = 2":  {BinaryOpAST{Equal, ColumnName{"a"}, NumericLiteral{2}}},
		"a < 2":  {BinaryOpAST{Less, ColumnName{"a"}, NumericLiteral{2}}},
		"a <= 2": {BinaryOpAST{LessOrEqual, ColumnName{"a"}, NumericLiteral{2}}},
		"a > 2":  {BinaryOpAST{Greater, ColumnName{"a"}, NumericLiteral{2}}},
		"a >= 2": {BinaryOpAST{GreaterOrEqual, ColumnName{"a"}, NumericLiteral{2}}},
		"a != 2": {BinaryOpAST{NotEqual, ColumnName{"a"}, NumericLiteral{2}}},
		"a <> 2": {BinaryOpAST{NotEqual, ColumnName{"a"}, NumericLiteral{2}}},
		// Plus/Minus Terms
		"a + 2": {BinaryOpAST{Plus, ColumnName{"a"}, NumericLiteral{2}}},
		"a - 2": {BinaryOpAST{Minus, ColumnName{"a"}, NumericLiteral{2}}},
		// Product/Division Terms
		"a * 2": {BinaryOpAST{Multiply, ColumnName{"a"}, NumericLiteral{2}}},
		"a / 2": {BinaryOpAST{Divide, ColumnName{"a"}, NumericLiteral{2}}},
		"a % 2": {BinaryOpAST{Modulo, ColumnName{"a"}, NumericLiteral{2}}},
		/// Operator Precedence
		"a AND b OR 2": {BinaryOpAST{Or,
			BinaryOpAST{And, ColumnName{"a"}, ColumnName{"b"}},
			NumericLiteral{2}}},
		"2 OR a AND b": {BinaryOpAST{Or,
			NumericLiteral{2},
			BinaryOpAST{And, ColumnName{"a"}, ColumnName{"b"}}}},
		"a * b + 2": {BinaryOpAST{Plus,
			BinaryOpAST{Multiply, ColumnName{"a"}, ColumnName{"b"}},
			NumericLiteral{2}}},
		"2 - a * b": {BinaryOpAST{Minus,
			NumericLiteral{2},
			BinaryOpAST{Multiply, ColumnName{"a"}, ColumnName{"b"}}}},
		/// Overriding Operator Precedence
		"a AND (b OR 2)": {BinaryOpAST{And,
			ColumnName{"a"},
			BinaryOpAST{Or, ColumnName{"b"}, NumericLiteral{2}}}},
		"(2 OR a) AND b": {BinaryOpAST{And,
			BinaryOpAST{Or, NumericLiteral{2}, ColumnName{"a"}},
			ColumnName{"b"}}},
		"a * (b + 2)": {BinaryOpAST{Multiply,
			ColumnName{"a"},
			BinaryOpAST{Plus, ColumnName{"b"}, NumericLiteral{2}}}},
		"(2 - a) * b": {BinaryOpAST{Multiply,
			BinaryOpAST{Minus, NumericLiteral{2}, ColumnName{"a"}},
			ColumnName{"b"}}},
		/// Operator Non-Precedence
		// TODO using more than one "same-level" operator does not work
		/*
			"a AND b AND 2": {BinaryOpAST{And,
				BinaryOpAST{And, ColumnName{"a"}, ColumnName{"b"}},
				NumericLiteral{2}}},
			"2 OR a OR b": {BinaryOpAST{Or,
				BinaryOpAST{Or, NumericLiteral{2}, ColumnName{"a"}},
				ColumnName{"b"}}},
			"a / b * 2": {BinaryOpAST{Multiply,
				BinaryOpAST{Divide, ColumnName{"a"}, ColumnName{"b"}},
				NumericLiteral{2}}},
			"a - b + 2": {BinaryOpAST{Plus,
				BinaryOpAST{Minus, ColumnName{"a"}, ColumnName{"b"}},
				NumericLiteral{2}}},
		*/
		/// Complex/Nested Expressions
		"(a + 1) % 2 = 0 OR b < 7.1, c": {
			BinaryOpAST{Or,
				BinaryOpAST{Equal,
					BinaryOpAST{Modulo,
						BinaryOpAST{Plus, ColumnName{"a"}, NumericLiteral{1}},
						NumericLiteral{2}},
					NumericLiteral{0}},
				BinaryOpAST{Less, ColumnName{"b"}, FloatLiteral{7.1}}},
			ColumnName{"c"}},
		/// Multiple Columns
		"a, 3.1, false, -2": {ColumnName{"a"}, FloatLiteral{3.1}, BoolLiteral{false}, NumericLiteral{-2}},
		`'日本語', 13`:         {StringLiteral{"日本語"}, NumericLiteral{13}},
		"concat(a, 'Pi', 3.1), b": {FuncAppAST{FuncName("concat"), ExpressionsAST{
			[]interface{}{ColumnName{"a"}, StringLiteral{"Pi"}, FloatLiteral{3.1}}}},
			ColumnName{"b"}},
	}

	Convey("Given a BQL parser", t, func() {
		p := NewBQLParser()

		for input, expected := range testCases {
			// avoid closure over loop variables
			input, expected := input, expected

			Convey(fmt.Sprintf("When parsing %s", input), func() {
				stmt := "SELECT " + input
				result, err := p.ParseStmt(stmt)

				Convey(fmt.Sprintf("Then the result should be %v", expected), func() {
					if expected == nil {
						So(err, ShouldNotBeNil)
					} else {
						// check there is no error
						So(err, ShouldBeNil)
						// check we got a proper SELECT statement
						So(result, ShouldHaveSameTypeAs, SelectStmt{})
						selectStmt := result.(SelectStmt)
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

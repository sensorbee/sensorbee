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
			ExpressionsAST{[]interface{}{RowValue{"a"}}}}},
		"f(2.1, 'a')": {FuncAppAST{FuncName("f"),
			ExpressionsAST{[]interface{}{FloatLiteral{2.1}, StringLiteral{"a"}}}}},
		// RowValue
		"a":         {RowValue{"a"}},
		"a, b":      {RowValue{"a"}, RowValue{"b"}},
		"A":         {RowValue{"A"}},
		"my_mem_27": {RowValue{"my_mem_27"}},
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
		"a OR 2": {BinaryOpAST{Or, RowValue{"a"}, NumericLiteral{2}}},
		// AND
		"a AND 2": {BinaryOpAST{And, RowValue{"a"}, NumericLiteral{2}}},
		// Comparisons
		"a = 2":  {BinaryOpAST{Equal, RowValue{"a"}, NumericLiteral{2}}},
		"a < 2":  {BinaryOpAST{Less, RowValue{"a"}, NumericLiteral{2}}},
		"a <= 2": {BinaryOpAST{LessOrEqual, RowValue{"a"}, NumericLiteral{2}}},
		"a > 2":  {BinaryOpAST{Greater, RowValue{"a"}, NumericLiteral{2}}},
		"a >= 2": {BinaryOpAST{GreaterOrEqual, RowValue{"a"}, NumericLiteral{2}}},
		"a != 2": {BinaryOpAST{NotEqual, RowValue{"a"}, NumericLiteral{2}}},
		"a <> 2": {BinaryOpAST{NotEqual, RowValue{"a"}, NumericLiteral{2}}},
		// Plus/Minus Terms
		"a + 2": {BinaryOpAST{Plus, RowValue{"a"}, NumericLiteral{2}}},
		"a - 2": {BinaryOpAST{Minus, RowValue{"a"}, NumericLiteral{2}}},
		// Product/Division Terms
		"a * 2": {BinaryOpAST{Multiply, RowValue{"a"}, NumericLiteral{2}}},
		"a / 2": {BinaryOpAST{Divide, RowValue{"a"}, NumericLiteral{2}}},
		"a % 2": {BinaryOpAST{Modulo, RowValue{"a"}, NumericLiteral{2}}},
		/// Operator Precedence
		"a AND b OR 2": {BinaryOpAST{Or,
			BinaryOpAST{And, RowValue{"a"}, RowValue{"b"}},
			NumericLiteral{2}}},
		"2 OR a AND b": {BinaryOpAST{Or,
			NumericLiteral{2},
			BinaryOpAST{And, RowValue{"a"}, RowValue{"b"}}}},
		"a * b + 2": {BinaryOpAST{Plus,
			BinaryOpAST{Multiply, RowValue{"a"}, RowValue{"b"}},
			NumericLiteral{2}}},
		"2 - a * b": {BinaryOpAST{Minus,
			NumericLiteral{2},
			BinaryOpAST{Multiply, RowValue{"a"}, RowValue{"b"}}}},
		/// Overriding Operator Precedence
		"a AND (b OR 2)": {BinaryOpAST{And,
			RowValue{"a"},
			BinaryOpAST{Or, RowValue{"b"}, NumericLiteral{2}}}},
		"(2 OR a) AND b": {BinaryOpAST{And,
			BinaryOpAST{Or, NumericLiteral{2}, RowValue{"a"}},
			RowValue{"b"}}},
		"a * (b + 2)": {BinaryOpAST{Multiply,
			RowValue{"a"},
			BinaryOpAST{Plus, RowValue{"b"}, NumericLiteral{2}}}},
		"(2 - a) * b": {BinaryOpAST{Multiply,
			BinaryOpAST{Minus, NumericLiteral{2}, RowValue{"a"}},
			RowValue{"b"}}},
		/// Operator Non-Precedence
		// TODO using more than one "same-level" operator does not work
		/*
			"a AND b AND 2": {BinaryOpAST{And,
				BinaryOpAST{And, RowValue{"a"}, RowValue{"b"}},
				NumericLiteral{2}}},
			"2 OR a OR b": {BinaryOpAST{Or,
				BinaryOpAST{Or, NumericLiteral{2}, RowValue{"a"}},
				RowValue{"b"}}},
			"a / b * 2": {BinaryOpAST{Multiply,
				BinaryOpAST{Divide, RowValue{"a"}, RowValue{"b"}},
				NumericLiteral{2}}},
			"a - b + 2": {BinaryOpAST{Plus,
				BinaryOpAST{Minus, RowValue{"a"}, RowValue{"b"}},
				NumericLiteral{2}}},
		*/
		/// Complex/Nested Expressions
		"(a + 1) % 2 = 0 OR b < 7.1, c": {
			BinaryOpAST{Or,
				BinaryOpAST{Equal,
					BinaryOpAST{Modulo,
						BinaryOpAST{Plus, RowValue{"a"}, NumericLiteral{1}},
						NumericLiteral{2}},
					NumericLiteral{0}},
				BinaryOpAST{Less, RowValue{"b"}, FloatLiteral{7.1}}},
			RowValue{"c"}},
		/// Multiple Columns
		"a, 3.1, false, -2": {RowValue{"a"}, FloatLiteral{3.1}, BoolLiteral{false}, NumericLiteral{-2}},
		`'日本語', 13`:         {StringLiteral{"日本語"}, NumericLiteral{13}},
		"concat(a, 'Pi', 3.1), b": {FuncAppAST{FuncName("concat"), ExpressionsAST{
			[]interface{}{RowValue{"a"}, StringLiteral{"Pi"}, FloatLiteral{3.1}}}},
			RowValue{"b"}},
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

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
	}{
		// Base Expressions
		"true":  {BoolLiteral{true}, Immutable},
		"NULL":  {NullLiteral{}, Immutable},
		"a":     {RowValue{"", "a"}, Immutable},
		"ts()":  {RowMeta{"", parser.TimestampMeta}, Immutable},
		"2":     {NumericLiteral{2}, Immutable},
		"1.2":   {FloatLiteral{1.2}, Immutable},
		`'bql'`: {StringLiteral{"bql"}, Immutable},
		"*":     {WildcardAST{}, Stable},
		// Function Application
		"f(a)": {FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}, Volatile},
		// Aggregate Function Application
		"count(a)": {FuncAppAST{parser.FuncName("count"),
			[]FlatExpression{AggInputRef{"_a4839edb"}}}, Volatile},
		// Composed Expressions
		"a OR 2":    {BinaryOpAST{parser.Or, RowValue{"", "a"}, NumericLiteral{2}}, Immutable},
		"a IS NULL": {BinaryOpAST{parser.Is, RowValue{"", "a"}, NullLiteral{}}, Immutable},
		"NOT a":     {UnaryOpAST{parser.Not, RowValue{"", "a"}}, Immutable},
		"NOT f(a)": {UnaryOpAST{parser.Not, FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}}, Volatile},
		// Comparisons
		"a = 2": {BinaryOpAST{parser.Equal, RowValue{"", "a"}, NumericLiteral{2}}, Immutable},
		"f(a) = 2": {BinaryOpAST{parser.Equal, FuncAppAST{parser.FuncName("f"),
			[]FlatExpression{RowValue{"", "a"}}}, NumericLiteral{2}}, Volatile},
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
						actual, _, err := ParserExprToMaybeAggregate(selectStmt.Projections[0], reg)
						So(err, ShouldBeNil)
						// compare it against our expectation
						So(actual, ShouldResemble, expected.e)
						So(actual.Volatility(), ShouldEqual, expected.v)
					}
				})
			})
		}
	})
}

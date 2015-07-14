package execution

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"strings"
)

// aliasedExpression represents an expression in a SELECT clause
// that has an output name associated and can be either an ordinary
// ("flat") expression or one that involves an aggregate computation.
type aliasedExpression struct {
	alias      string
	expr       FlatExpression
	aggrInputs map[string]FlatExpression
}

// Explanation of the Aggregation Workflow
// ---------------------------------------
// For a SELECT or CREATE STREAM FROM SELECT statement, we deal mostly
// with the "projections", i.e., the terms after the SELECT keyword.
// For a non-aggregate statement, we receive an instance of parser.Expression
// and transform it to an instance of execution.FlatExpression (which is
// guaranteed not to have any aggregate functions in it). This FlatExpression
// is then converted to an Evaluator that can be evaluated with a row
// as an input.
//
// Example 1: The BQL expression "x:a" will be converted to a
// `parser.RowValue{"x", "a"}` and then to an `execution.RowValue{"x", "a"}` and
// then to an `execution.PathAccess{"x.a"}`. If this `PathAccess{"x.a"}` is evaluated
// on a row like `data.Map{"x": data.Map{"a": data.Int(2), "b": data.Int(6)}}`,
// it will return `data.Int(2)`.
//
// Example 2: The BQL expression "f(x:a)" will be converted to a
// `parser.FuncAppAST{"f", parser.RowValue{"x", "a"}}` and then
// (if f is not an aggregate function) to an `execution.FuncAppAST{"f",
// execution.RowValue{"x", "a"}}` and then to a `execution.funcApp{"f", ...,
// []Evaluator{execution.PathAccess{"x.a"}}, ...}`. If this `funcApp` is evaluated
// on a row like `data.Map{"x": data.Map{"a": data.Int(2),  "b": data.Int(6)}}`,
// the `PathAccess` evaluator will first return extract `data.Int(2)` and then
// the `funcApp` evaluator will compute the result of f using `data.Int(2)`
// as an input.
//
// For an expression that involves an aggregate function call, the flow
// is necessarily different:
// - The projections are not evaluated using the original rows as an input,
//   but using groups of rows that have the same values in the GROUP BY columns.
// - That means that in a pre-processing step, the input rows need to be traversed
//   and the groups described above must be created.
// - Each of those groups consists of
//   - the values that are common to all rows in a group
//   - a list of values that are to be aggregated.
// - After the groups have been created, each list of values for aggregation
//   is transformed to an array and added to the set of "common values" of
//   that group.
// - The aggregate function is executed like a normal function, but using the
//   array mentioned before as an input.
//
// Example: The BQL expression "avg(x:a)" will be converted to a
// `parser.FuncAppAST{"avg", parser.RowValue{"x", "a"}}`and then
// to an `execution.FuncAppAST{"avg", execution.AggInputRef{"randstr"}}`
// together with a mapping `{"randstr": execution.RowValue{"x", "a"}}`.
// Those FlatExpressions are then converted to Evaluators so that we obtain
// `execution.funcApp{"avg", ..., []Evaluator{execution.PathAccess{"randstr"}}, ...}`
// and the map `{"randstr": execution.PathAccess{"x.a"}}`.
// For each row, the correct group will be computed and the data computed
// by the PathAccess{"x.a"} evaluator appended to a list, so that at the
// end of this process there is a group list such as:
//  group values | common data     | aggregate input
//  [1]          | {"x": {"b": 1}} | {"randstr": [3.5, 1.7, 0.9, 4.5, ...]}
//  [2]          | {"x": {"b": 2}} | {"randstr": [1.2, 0.8, 2.3]}
// This data is then merged so that we obtain a list such as:
//  common data
//  {"x": {"b": 1}, "randstr": [3.5, 1.7, 0.9, 4.5, ...]}
//  {"x": {"b": 2}, "randstr": [1.2, 0.8, 2.3]}
// This list looks exactly like a usual set of input rows so that the
// `funcApp{"avg", ...}` evaluator can be used normally.

// ParserExprToFlatExpr converts an expression obtained by the BQL parser
// to a FlatExpression, i.e., there are only expressions contained that
// can be evaluated on one single row and return an (unnamed) value.
// In particular, this fails for Expressions containing aggregate functions.
func ParserExprToFlatExpr(e parser.Expression, reg udf.FunctionRegistry) (FlatExpression, error) {
	switch obj := e.(type) {
	case parser.RowMeta:
		return RowMeta{obj.Relation, obj.MetaType}, nil
	case parser.RowValue:
		return RowValue{obj.Relation, obj.Column}, nil
	case parser.AliasAST:
		return ParserExprToFlatExpr(obj.Expr, reg)
	case parser.NullLiteral:
		return NullLiteral{}, nil
	case parser.NumericLiteral:
		return NumericLiteral{obj.Value}, nil
	case parser.FloatLiteral:
		return FloatLiteral{obj.Value}, nil
	case parser.BoolLiteral:
		return BoolLiteral{obj.Value}, nil
	case parser.StringLiteral:
		return StringLiteral{obj.Value}, nil
	case parser.BinaryOpAST:
		// recurse
		left, err := ParserExprToFlatExpr(obj.Left, reg)
		if err != nil {
			return nil, err
		}
		right, err := ParserExprToFlatExpr(obj.Right, reg)
		if err != nil {
			return nil, err
		}
		return BinaryOpAST{obj.Op, left, right}, nil
	case parser.UnaryOpAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return UnaryOpAST{obj.Op, expr}, nil
	case parser.FuncAppAST:
		// look up the function
		function, err := reg.Lookup(string(obj.Function), len(obj.Expressions))
		if err != nil {
			return nil, err
		}
		// fail if this is an aggregate function
		if isAggregateFunc(function, len(obj.Expressions), reg) {
			err := fmt.Errorf("you cannot use aggregate function '%s' "+
				"in a flat expression", obj.Function)
			return nil, err
		}
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			expr, err := ParserExprToFlatExpr(ast, reg)
			if err != nil {
				return nil, err
			}
			exprs[i] = expr
		}
		return FuncAppAST{obj.Function, exprs}, nil
	case parser.Wildcard:
		return WildcardAST{}, nil
	}
	err := fmt.Errorf("don't know how to convert type %#v", e)
	return nil, err
}

// ParserExprToMaybeAggregate converts an expression obtained by the BQL
// parser into a data structure where the aggregate and the non-aggregate
// parts are separated.
func ParserExprToMaybeAggregate(e parser.Expression, reg udf.FunctionRegistry) (FlatExpression, map[string]FlatExpression, error) {
	switch obj := e.(type) {
	default:
		// elementary types
		expr, err := ParserExprToFlatExpr(e, reg)
		return expr, nil, err
	case parser.AliasAST:
		return ParserExprToMaybeAggregate(obj.Expr, reg)
	case parser.BinaryOpAST:
		// recurse
		left, leftAgg, err := ParserExprToMaybeAggregate(obj.Left, reg)
		if err != nil {
			return nil, nil, err
		}
		right, rightAgg, err := ParserExprToMaybeAggregate(obj.Right, reg)
		if err != nil {
			return nil, nil, err
		}
		var returnAgg map[string]FlatExpression
		if leftAgg != nil {
			returnAgg = leftAgg
			for key, val := range rightAgg {
				returnAgg[key] = val
			}
		} else if leftAgg == nil {
			returnAgg = rightAgg
		}
		return BinaryOpAST{obj.Op, left, right}, returnAgg, nil
	case parser.UnaryOpAST:
		// recurse
		expr, agg, err := ParserExprToMaybeAggregate(obj.Expr, reg)
		if err != nil {
			return nil, nil, err
		}
		return UnaryOpAST{obj.Op, expr}, agg, nil
	case parser.FuncAppAST:
		// look up the function
		function, err := reg.Lookup(string(obj.Function), len(obj.Expressions))
		if err != nil {
			return nil, nil, err
		}
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		returnAgg := map[string]FlatExpression{}
		if isAggregateFunc(function, len(obj.Expressions), reg) {
			// we have a setting like
			//  SELECT udaf(x+1, 'state', c) ... GROUP BY c
			// where some parameters are aggregates, others aren't.
			for i, ast := range obj.Expressions {
				// this expression must be flat, there must not be other aggregates
				expr, err := ParserExprToFlatExpr(ast, reg)
				if err != nil {
					// return a prettier error message
					if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
						err = fmt.Errorf("aggregate functions cannot be nested")
					}
					return nil, nil, err
				}
				isAggr := function.IsAggregationParameter(i + 1) // parameters count from 1
				if isAggr {
					// this is an aggregation parameter, we will replace
					// it by a reference to the aggregated list of values
					h := sha1.New()
					h.Write([]byte(fmt.Sprintf("%s", expr.Repr())))
					exprId := "_" + hex.EncodeToString(h.Sum(nil))[:8]
					exprs[i] = AggInputRef{exprId}
					returnAgg[exprId] = expr
				} else {
					// this is a non-aggregate parameter, use as is
					exprs[i] = expr
				}
			}
		} else {
			for i, ast := range obj.Expressions {
				expr, agg, err := ParserExprToMaybeAggregate(ast, reg)
				if err != nil {
					return nil, nil, err
				}
				for key, val := range agg {
					returnAgg[key] = val
				}
				exprs[i] = expr
			}
			if len(returnAgg) == 0 {
				returnAgg = nil
			}
		}
		return FuncAppAST{obj.Function, exprs}, returnAgg, nil
	}
	err := fmt.Errorf("don't know how to convert type %#v", e)
	return nil, nil, err
}

// FlatExpression represents an expression that can be completely
// evaluated on a single row and results in an unnamed value. In
// particular, it cannot contain/represent a call to an aggregate
// function.
type FlatExpression interface {
	// Repr returns a string representation that can be used to
	// identify this expression (e.g., "stream:col+3") and used as
	// a dictionary key for finding duplicate expressions.
	Repr() string

	// Columns returns a list of RowValues used in this expression.
	Columns() []RowValue
}

type BinaryOpAST struct {
	Op    parser.Operator
	Left  FlatExpression
	Right FlatExpression
}

func (b BinaryOpAST) Repr() string {
	return fmt.Sprintf("%s%s%s", b.Left.Repr(), b.Op, b.Right.Repr())
}

func (b BinaryOpAST) Columns() []RowValue {
	return append(b.Left.Columns(), b.Right.Columns()...)
}

type UnaryOpAST struct {
	Op   parser.Operator
	Expr FlatExpression
}

func (u UnaryOpAST) Repr() string {
	return fmt.Sprintf("%s%s", u.Op, u.Expr.Repr())
}

func (u UnaryOpAST) Columns() []RowValue {
	return u.Expr.Columns()
}

type FuncAppAST struct {
	Function    parser.FuncName
	Expressions []FlatExpression
}

func (f FuncAppAST) Repr() string {
	reprs := make([]string, len(f.Expressions))
	for i, e := range f.Expressions {
		reprs[i] = e.Repr()
	}
	return fmt.Sprintf("%s(%s)", f.Function, strings.Join(reprs, ","))
}

func (f FuncAppAST) Columns() []RowValue {
	allColumns := []RowValue{}
	for _, e := range f.Expressions {
		allColumns = append(allColumns, e.Columns()...)
	}
	return allColumns
}

type WildcardAST struct {
}

func (w WildcardAST) Repr() string {
	return "*"
}

func (w WildcardAST) Columns() []RowValue {
	return nil
}

type AggInputRef struct {
	Ref string
}

func (a AggInputRef) Repr() string {
	return a.Ref
}

func (a AggInputRef) Columns() []RowValue {
	return nil
}

type RowValue struct {
	Relation string
	Column   string
}

func (rv RowValue) Repr() string {
	return fmt.Sprintf("%s:%s", rv.Relation, rv.Column)
}

func (rv RowValue) Columns() []RowValue {
	return []RowValue{rv}
}

type RowMeta struct {
	Relation string
	MetaType parser.MetaInformation
}

func (rm RowMeta) Repr() string {
	return fmt.Sprintf("%#v", rm)
}

func (rm RowMeta) Columns() []RowValue {
	return nil
}

type NumericLiteral struct {
	Value int64
}

func (l NumericLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l NumericLiteral) Columns() []RowValue {
	return nil
}

type FloatLiteral struct {
	Value float64
}

func (l FloatLiteral) Repr() string {
	return fmt.Sprintf("%vf", l.Value)
}

func (l FloatLiteral) Columns() []RowValue {
	return nil
}

type NullLiteral struct {
}

func (l NullLiteral) Repr() string {
	return "NULL"
}

func (l NullLiteral) Columns() []RowValue {
	return nil
}

type BoolLiteral struct {
	Value bool
}

func (l BoolLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l BoolLiteral) Columns() []RowValue {
	return nil
}

type StringLiteral struct {
	Value string
}

func (l StringLiteral) Repr() string {
	return fmt.Sprintf("%s", l.Value)
}

func (l StringLiteral) Columns() []RowValue {
	return nil
}

type AggFuncAppAST struct {
	Function   udf.UDF
	Expression FlatExpression
}

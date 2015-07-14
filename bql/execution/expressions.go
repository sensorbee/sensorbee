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
	aggrInputs map[string]AggFuncAppAST
}

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
func ParserExprToMaybeAggregate(e parser.Expression, reg udf.FunctionRegistry) (FlatExpression, map[string]AggFuncAppAST, error) {
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
		var returnAgg map[string]AggFuncAppAST
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
		if isAggregateFunc(function, len(obj.Expressions), reg) {
			// we can only have one parameter
			if len(obj.Expressions) != 1 {
				err := fmt.Errorf("aggregate functions must have exactly one parameter")
				return nil, nil, err
			}
			// this expression must be flat, there must not be other aggregates
			expr, err := ParserExprToFlatExpr(obj.Expressions[0], reg)
			if err != nil {
				// return a prettier error message
				if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
					err = fmt.Errorf("aggregate functions cannot be nested")
				}
				return nil, nil, err
			}
			// get a string that identifies this sub-expression and that
			// can be used for key access in a Map (i.e., no special chars).
			// we use the first characters of the SHA1 hash of a string
			// representation like "count(x:a+1)" and prefix it with an "a"
			// to prevent all-numeric strings
			h := sha1.New()
			h.Write([]byte(fmt.Sprintf("%s(%s)", obj.Function, expr.Repr())))
			funcId := "a" + hex.EncodeToString(h.Sum(nil))[:8]
			a := AggFuncAppAST{function, expr}
			return AggFuncAppRef{funcId}, map[string]AggFuncAppAST{funcId: a}, nil
		} else {
			// compute child expressions
			exprs := make([]FlatExpression, len(obj.Expressions))
			returnAgg := map[string]AggFuncAppAST{}
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
			return FuncAppAST{obj.Function, exprs}, returnAgg, nil
		}
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
	return nil
}

type WildcardAST struct {
}

func (w WildcardAST) Repr() string {
	return "*"
}

func (w WildcardAST) Columns() []RowValue {
	return nil
}

type AggFuncAppRef struct {
	Ref string
}

func (af AggFuncAppRef) Repr() string {
	return af.Ref
}

func (af AggFuncAppRef) Columns() []RowValue {
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

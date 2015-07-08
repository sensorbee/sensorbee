package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
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
func ParserExprToFlatExpr(e parser.Expression, isAggregate func(string) bool) (FlatExpression, error) {
	switch obj := e.(type) {
	case parser.RowMeta:
		return RowMeta{obj.Relation, obj.MetaType}, nil
	case parser.RowValue:
		return RowValue{obj.Relation, obj.Column}, nil
	case parser.AliasAST:
		return ParserExprToFlatExpr(obj.Expr, isAggregate)
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
		left, err := ParserExprToFlatExpr(obj.Left, isAggregate)
		if err != nil {
			return nil, err
		}
		right, err := ParserExprToFlatExpr(obj.Right, isAggregate)
		if err != nil {
			return nil, err
		}
		return BinaryOpAST{obj.Op, left, right}, nil
	case parser.UnaryOpAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.Expr, isAggregate)
		if err != nil {
			return nil, err
		}
		return UnaryOpAST{obj.Op, expr}, nil
	case parser.FuncAppAST:
		// fail if this is an aggregate function
		if isAggregate(string(obj.Function)) {
			err := fmt.Errorf("you cannot use aggregate function '%s' "+
				"in a flat expression", obj.Function)
			return nil, err
		}
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			expr, err := ParserExprToFlatExpr(ast, isAggregate)
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
func ParserExprToMaybeAggregate(e parser.Expression, isAggregate func(string) bool) (FlatExpression, map[string]AggFuncAppAST, error) {
	switch obj := e.(type) {
	default:
		// elementary types
		expr, err := ParserExprToFlatExpr(e, isAggregate)
		return expr, nil, err
	case parser.BinaryOpAST:
		// recurse
		left, leftAgg, err := ParserExprToMaybeAggregate(obj.Left, isAggregate)
		if err != nil {
			return nil, nil, err
		}
		right, rightAgg, err := ParserExprToMaybeAggregate(obj.Right, isAggregate)
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
		expr, agg, err := ParserExprToMaybeAggregate(obj.Expr, isAggregate)
		if err != nil {
			return nil, nil, err
		}
		return UnaryOpAST{obj.Op, expr}, agg, nil
	case parser.FuncAppAST:
		if isAggregate(string(obj.Function)) {
			// we can only have one parameter
			if len(obj.Expressions) != 1 {
				err := fmt.Errorf("aggregate functions must have exactly one parameter")
				return nil, nil, err
			}
			// this expression must be flat, there must not be other aggregates
			expr, err := ParserExprToFlatExpr(obj.Expressions[0], isAggregate)
			if err != nil {
				// return a prettier error message
				if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
					err = fmt.Errorf("aggregate functions cannot be nested")
				}
				return nil, nil, err
			}
			// get a string that identifies this sub-expression
			funcId := fmt.Sprintf("#%s(%s)#", obj.Function, expr.Repr())
			a := AggFuncAppAST{obj.Function, expr}
			return AggFuncAppRef{funcId}, map[string]AggFuncAppAST{funcId: a}, nil
		} else {
			// compute child expressions
			exprs := make([]FlatExpression, len(obj.Expressions))
			returnAgg := map[string]AggFuncAppAST{}
			for i, ast := range obj.Expressions {
				expr, agg, err := ParserExprToMaybeAggregate(ast, isAggregate)
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
	Repr() string
}

type BinaryOpAST struct {
	Op    parser.Operator
	Left  FlatExpression
	Right FlatExpression
}

func (b BinaryOpAST) Repr() string {
	return fmt.Sprintf("%s%s%s", b.Left.Repr(), b.Op, b.Right.Repr())
}

type UnaryOpAST struct {
	Op   parser.Operator
	Expr FlatExpression
}

func (u UnaryOpAST) Repr() string {
	return fmt.Sprintf("%s%s", u.Op, u.Expr.Repr())
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

type WildcardAST struct {
}

func (w WildcardAST) Repr() string {
	return "*"
}

type AggFuncAppRef struct {
	Ref string
}

func (af AggFuncAppRef) Repr() string {
	return af.Ref
}

type RowValue struct {
	Relation string
	Column   string
}

func (rv RowValue) Repr() string {
	return fmt.Sprintf("%s:%s", rv.Relation, rv.Column)
}

type RowMeta struct {
	Relation string
	MetaType parser.MetaInformation
}

func (rm RowMeta) Repr() string {
	return fmt.Sprintf("%#v", rm)
}

type NumericLiteral struct {
	Value int64
}

func (l NumericLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

type FloatLiteral struct {
	Value float64
}

func (l FloatLiteral) Repr() string {
	return fmt.Sprintf("%vf", l.Value)
}

type NullLiteral struct {
}

func (l NullLiteral) Repr() string {
	return "NULL"
}

type BoolLiteral struct {
	Value bool
}

func (l BoolLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

type StringLiteral struct {
	Value string
}

func (l StringLiteral) Repr() string {
	return fmt.Sprintf("%s", l.Value)
}

type AggFuncAppAST struct {
	Function   parser.FuncName
	Expression FlatExpression
}

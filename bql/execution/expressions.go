package execution

import (
	"fmt"
	"math/rand"
	"pfi/sensorbee/sensorbee/bql/parser"
	"strings"
)

// aliasedExpression represents an expression in a SELECT clause
// that has an output name associated and can be either an ordinary
// ("flat") expression or one that involves an aggregate computation.
type aliasedExpression struct {
	alias      string
	expr       FlatExpression
	aggrInputs map[string]FuncAppAST
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
func ParserExprToMaybeAggregate(e parser.Expression, isAggregate func(string) bool) (FlatExpression, map[string]FuncAppAST, error) {
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
		var returnAgg map[string]FuncAppAST
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
			// get a random string that identifies this sub-expression
			funcId := fmt.Sprintf("aggr_res_%d", rand.Int63())
			return AggFunResult{string(obj.Function), funcId}, map[string]FuncAppAST{funcId: {obj.Function, []FlatExpression{expr}}}, nil
		} else {
			// compute child expressions
			exprs := make([]FlatExpression, len(obj.Expressions))
			returnAgg := map[string]FuncAppAST{}
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
	Hoge() bool
}

type BinaryOpAST struct {
	Op    parser.Operator
	Left  FlatExpression
	Right FlatExpression
}

func (b BinaryOpAST) Hoge() bool {
	return b.Left.Hoge() && b.Right.Hoge()
}

type UnaryOpAST struct {
	Op   parser.Operator
	Expr FlatExpression
}

func (u UnaryOpAST) Hoge() bool {
	return u.Expr.Hoge()
}

type FuncAppAST struct {
	Function    parser.FuncName
	Expressions []FlatExpression
}

func (f FuncAppAST) Hoge() bool {
	foldable := true
	for _, expr := range f.Expressions {
		if !expr.Hoge() {
			foldable = false
			break
		}
	}
	return foldable
}

type WildcardAST struct {
}

func (w WildcardAST) Hoge() bool {
	return false
}

type AggFunResult struct {
	Function string
	ParamStr string
}

func (afr AggFunResult) Hoge() bool {
	return false
}

type RowValue struct {
	Relation string
	Column   string
}

func (rv RowValue) Hoge() bool {
	return false
}

type RowMeta struct {
	Relation string
	MetaType parser.MetaInformation
}

func (rm RowMeta) Hoge() bool {
	return false
}

type NumericLiteral struct {
	Value int64
}

func (l NumericLiteral) Hoge() bool {
	return true
}

type FloatLiteral struct {
	Value float64
}

func (l FloatLiteral) Hoge() bool {
	return true
}

type NullLiteral struct {
}

func (l NullLiteral) Hoge() bool {
	return true
}

type BoolLiteral struct {
	Value bool
}

func (l BoolLiteral) Hoge() bool {
	return true
}

type StringLiteral struct {
	Value string
}

func (l StringLiteral) Hoge() bool {
	return true
}

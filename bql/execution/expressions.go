package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
)

type aliasedExpression struct {
	alias string
	expr  FlatExpression
}

// ParserExprToFlatExpr converts and expression obtained by the BQL parser
// to a FlatExpression, i.e., there are only expressions contained that
// can be evaluated on one single row and return an (unnamed) value.
// In particular, this fails for Expressions containing aggregate functions.
//
// TODO No, it does not.
func ParserExprToFlatExpr(e parser.Expression) (FlatExpression, error) {
	switch obj := e.(type) {
	case parser.RowMeta:
		return RowMeta{obj.Relation, obj.MetaType}, nil
	case parser.RowValue:
		return RowValue{obj.Relation, obj.Column}, nil
	case parser.AliasAST:
		return ParserExprToFlatExpr(obj.Expr)
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
		left, err := ParserExprToFlatExpr(obj.Left)
		if err != nil {
			return nil, err
		}
		right, err := ParserExprToFlatExpr(obj.Right)
		if err != nil {
			return nil, err
		}
		return BinaryOpAST{obj.Op, left, right}, nil
	case parser.UnaryOpAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.Expr)
		if err != nil {
			return nil, err
		}
		return UnaryOpAST{obj.Op, expr}, nil
	case parser.FuncAppAST:
		// compute child Evaluators
		exprs := make([]FlatExpression, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			expr, err := ParserExprToFlatExpr(ast)
			if err != nil {
				return nil, err
			}
			exprs[i] = expr
		}
		// TODO fail if this is an aggregate function
		return FuncAppAST{obj.Function, exprs}, nil
	case parser.Wildcard:
		return WildcardAST{}, nil
	}
	err := fmt.Errorf("don't know how to convert type %#v", e)
	return nil, err
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

package execution

import (
	"fmt"
	"math"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
)

type Evaluator interface {
	Eval(input tuple.Value) (tuple.Value, error)
}

// ExpressionToEvaluator takes one of the Expression structs that result
// from parsing a BQL Expression (see parser/ast.go) and turns it into
// an Evaluator that can be used to evaluate an expression given a particular
// input Value.
func ExpressionToEvaluator(ast interface{}) (Evaluator, error) {
	switch obj := ast.(type) {
	case parser.ColumnName:
		return &PathAccess{obj.Name}, nil
	case parser.NumericLiteral:
		return &IntConstant{obj.Value}, nil
	case parser.FloatLiteral:
		return &FloatConstant{obj.Value}, nil
	case parser.BoolLiteral:
		return &BoolConstant{obj.Value}, nil
	case parser.BinaryOpAST:
		// recurse
		left, err := ExpressionToEvaluator(obj.Left)
		if err != nil {
			return nil, err
		}
		right, err := ExpressionToEvaluator(obj.Right)
		if err != nil {
			return nil, err
		}
		// assemble both children with the correct operator
		bo := binOp{left, right}
		switch obj.Op {
		default:
			err := fmt.Errorf("don't know how to evaluate binary operation %v", obj.Op)
			return nil, err
		case parser.Or:
			return Or(bo), nil
		case parser.And:
			return And(bo), nil
		case parser.Equal:
			return Equal(bo), nil
		case parser.Less:
			return Less(bo), nil
		case parser.LessOrEqual:
			return LessOrEqual(bo), nil
		case parser.Greater:
			return Greater(bo), nil
		case parser.GreaterOrEqual:
			return GreaterOrEqual(bo), nil
		case parser.NotEqual:
			return Not(Equal(bo)), nil
		case parser.Plus:
			return Plus(bo), nil
		case parser.Minus:
			return Minus(bo), nil
		case parser.Multiply:
			return Multiply(bo), nil
		case parser.Divide:
			return Divide(bo), nil
		case parser.Modulo:
			return Modulo(bo), nil
		}
	}
	err := fmt.Errorf("don't know how to evaluate type %#v", ast)
	return nil, err
}

// IntConstant always returns the same integer value, independent
// of the input.
type IntConstant struct {
	value int64
}

func (i *IntConstant) Eval(input tuple.Value) (tuple.Value, error) {
	return tuple.Int(i.value), nil
}

// FloatConstant always returns the same float value, independent
// of the input.
type FloatConstant struct {
	value float64
}

func (f *FloatConstant) Eval(input tuple.Value) (tuple.Value, error) {
	return tuple.Float(f.value), nil
}

// BoolConstant always returns the same boolean value, independent
// of the input.
type BoolConstant struct {
	value bool
}

func (b *BoolConstant) Eval(input tuple.Value) (tuple.Value, error) {
	return tuple.Bool(b.value), nil
}

// PathAccess only works for maps and returns the Value at the given
// JSON path.
type PathAccess struct {
	path string
}

func (fa *PathAccess) Eval(input tuple.Value) (tuple.Value, error) {
	aMap, err := input.AsMap()
	if err != nil {
		return nil, err
	}
	return aMap.Get(fa.path)
}

type binOp struct {
	left  Evaluator
	right Evaluator
}

func (bo *binOp) evalLeftAndRight(input tuple.Value) (tuple.Value, tuple.Value, error) {
	leftRes, err := bo.left.Eval(input)
	if err != nil {
		return nil, nil, err
	}
	rightRes, err := bo.right.Eval(input)
	if err != nil {
		return nil, nil, err
	}
	return leftRes, rightRes, nil
}

/// Binary Logical Operations

type logBinOp struct {
	binOp
	leftSidePreCondition bool
	leftSidePreReturn    bool
}

func (lbo *logBinOp) Eval(input tuple.Value) (tuple.Value, error) {
	leftRes, err := lbo.left.Eval(input)
	if err != nil {
		return nil, err
	}
	leftBool, err := tuple.ToBool(leftRes)
	if err != nil {
		return nil, err
	}
	// support early return if the left side has a certain value
	if leftBool == lbo.leftSidePreCondition {
		return tuple.Bool(lbo.leftSidePreReturn), nil
	}
	rightRes, err := lbo.right.Eval(input)
	if err != nil {
		return nil, err
	}
	rightBool, err := tuple.ToBool(rightRes)
	if err != nil {
		return nil, err
	}
	if rightBool {
		return tuple.Bool(true), nil
	}
	return tuple.Bool(false), nil
}

func Or(bo binOp) Evaluator {
	return &logBinOp{bo,
		// if the left side is true, we return true early
		true, true}
}

func And(bo binOp) Evaluator {
	return &logBinOp{bo,
		// if the left side is false, we return false early
		false, false}
}

/// A Unary Logical Operation

type not struct {
	neg Evaluator
}

func (n *not) Eval(input tuple.Value) (tuple.Value, error) {
	neg, err := n.neg.Eval(input)
	if err != nil {
		return nil, err
	}
	negBool, err := neg.AsBool()
	if err != nil {
		return nil, err
	}
	return tuple.Bool(!negBool), nil
}

func Not(e Evaluator) Evaluator {
	return &not{e}
}

/// Binary Comparison Operations

// compBinOp provides functionality for comparing two values
type compBinOp struct {
	binOp
	cmpOp func(tuple.Value, tuple.Value) (bool, error)
}

func (cbo *compBinOp) Eval(input tuple.Value) (tuple.Value, error) {
	leftVal, rightVal, err := cbo.evalLeftAndRight(input)
	if err != nil {
		return nil, err
	}
	res, err := cbo.cmpOp(leftVal, rightVal)
	if err != nil {
		return nil, err
	}
	return tuple.Bool(res), nil
}

func Equal(bo binOp) Evaluator {
	cmpOp := func(leftVal tuple.Value, rightVal tuple.Value) (bool, error) {
		leftType := leftVal.Type()
		rightType := rightVal.Type()
		eq := false
		if leftType == rightType {
			eq = reflect.DeepEqual(leftVal, rightVal)
		} else if leftType == tuple.TypeInt && rightType == tuple.TypeFloat {
			l, _ := leftVal.AsInt()
			r, _ := rightVal.AsFloat()
			// convert left to float to get 2 == 2.0
			eq = float64(l) == r
		} else if leftType == tuple.TypeFloat && rightType == tuple.TypeInt {
			l, _ := leftVal.AsFloat()
			r, _ := rightVal.AsInt()
			// convert right to float to get 2.0 == 2
			eq = l == float64(r)
		}
		return eq, nil

	}
	return &compBinOp{bo, cmpOp}
}

func Less(bo binOp) Evaluator {
	cmpOp := func(leftVal tuple.Value, rightVal tuple.Value) (bool, error) {
		leftType := leftVal.Type()
		rightType := rightVal.Type()
		stdErr := fmt.Errorf("cannot compare %T and %T", leftVal, rightVal)
		if leftType == rightType {
			retVal := false
			switch leftType {
			default:
				return false, stdErr
			case tuple.TypeInt:
				l, _ := leftVal.AsInt()
				r, _ := rightVal.AsInt()
				retVal = l < r
			case tuple.TypeFloat:
				l, _ := leftVal.AsFloat()
				r, _ := rightVal.AsFloat()
				retVal = l < r
			case tuple.TypeString:
				l, _ := leftVal.AsString()
				r, _ := rightVal.AsString()
				retVal = l < r
			case tuple.TypeBool:
				l, _ := leftVal.AsBool()
				r, _ := rightVal.AsBool()
				retVal = (l == false) && (r == true)
			case tuple.TypeTimestamp:
				l, _ := leftVal.AsTimestamp()
				r, _ := rightVal.AsTimestamp()
				retVal = l.Before(r)
			}
			return retVal, nil
		} else if leftType == tuple.TypeInt && rightType == tuple.TypeFloat {
			// left is integer
			l, _ := leftVal.AsInt()
			// right is float; also convert left to float to avoid overflow
			r, _ := rightVal.AsFloat()
			return float64(l) < r, nil
		} else if leftType == tuple.TypeFloat && rightType == tuple.TypeInt {
			// left is float
			l, _ := leftVal.AsFloat()
			// right is int; convert right to float to avoid overflow
			r, _ := rightVal.AsInt()
			return l < float64(r), nil
		}
		return false, stdErr
	}
	return &compBinOp{bo, cmpOp}
}

func LessOrEqual(bo binOp) Evaluator {
	return Or(binOp{Less(bo), Equal(bo)})
}

func Greater(bo binOp) Evaluator {
	return Not(LessOrEqual(bo))
}

func GreaterOrEqual(bo binOp) Evaluator {
	return Not(Less(bo))
}

func NotEqual(bo binOp) Evaluator {
	return Not(Equal(bo))
}

/// Binary Numerical Operations

// numBinOp provides functionality for evaluating binary operations
// on two numeric Values (int64 or float64 or combinations of them).
type numBinOp struct {
	binOp
	verb    string
	intOp   func(int64, int64) int64
	floatOp func(float64, float64) float64
}

func (nbo *numBinOp) Eval(input tuple.Value) (v tuple.Value, err error) {
	defer func() {
		// catch panic from (say) integer division by 0
		if r := recover(); r != nil {
			v = nil
			err = fmt.Errorf("%s", r)
		}
	}()
	// evalate both sides
	leftVal, rightVal, err := nbo.evalLeftAndRight(input)
	if err != nil {
		return nil, err
	}
	leftType := leftVal.Type()
	rightType := rightVal.Type()
	stdErr := fmt.Errorf("cannot %s %T and %T", nbo.verb, leftVal, rightVal)
	// if we have same types (both int64 or both float64, apply
	// the corresponding operation)
	if leftType == rightType {
		switch leftType {
		default:
			return nil, stdErr
		case tuple.TypeInt:
			l, _ := leftVal.AsInt()
			r, _ := rightVal.AsInt()
			return tuple.Int(nbo.intOp(l, r)), nil
		case tuple.TypeFloat:
			l, _ := leftVal.AsFloat()
			r, _ := rightVal.AsFloat()
			return tuple.Float(nbo.floatOp(l, r)), nil
		}
	} else if leftType == tuple.TypeInt && rightType == tuple.TypeFloat {
		// left is integer
		l, _ := leftVal.AsInt()
		// right is float; also convert left to float, possibly losing precision
		r, _ := rightVal.AsFloat()
		return tuple.Float(nbo.floatOp(float64(l), r)), nil
	} else if leftType == tuple.TypeFloat && rightType == tuple.TypeInt {
		// left is float
		l, _ := leftVal.AsFloat()
		// right is int; convert right to float, possibly losing precision
		r, _ := rightVal.AsInt()
		return tuple.Float(nbo.floatOp(l, float64(r))), nil
	}
	return nil, stdErr
}

func Plus(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a + b
	}
	floatOp := func(a, b float64) float64 {
		return a + b
	}
	return &numBinOp{bo, "add", intOp, floatOp}
}

func Minus(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a - b
	}
	floatOp := func(a, b float64) float64 {
		return a - b
	}
	return &numBinOp{bo, "subtract", intOp, floatOp}
}

func Multiply(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a * b
	}
	floatOp := func(a, b float64) float64 {
		return a * b
	}
	return &numBinOp{bo, "multiply", intOp, floatOp}
}

func Divide(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a / b
	}
	floatOp := func(a, b float64) float64 {
		return a / b
	}
	return &numBinOp{bo, "divide", intOp, floatOp}
}

func Modulo(bo binOp) Evaluator {
	intOp := func(a, b int64) int64 {
		return a % b
	}
	floatOp := func(a, b float64) float64 {
		return math.Mod(a, b)
	}
	return &numBinOp{bo, "compute modulo for", intOp, floatOp}
}

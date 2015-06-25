package execution

import (
	"fmt"
	"math"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"reflect"
	"strings"
)

// An Evaluator represents an expression such as `colX + 2` or
// `t1:col AND t2:col` and can be evaluated, given the actual data
// contained in one row.
type Evaluator interface {
	// Eval evaluates the expression that this Evaluator represents
	// on the given input data. Note that in order to deal with
	// joins and timestamps properly, the input data must have the shape
	//   {"alias_1": {"col_1": ..., "col_2": ...},
	//    "alias_1:meta:x": (meta datum "x" for alias_1's row),
	//    "alias_2": {"col_1": ..., "col_2": ...},
	//    "alias_2:meta:x": (meta datum "x" for alias_2's row),
	//    ...}
	// and every caller (in particular all execution plans)
	// must ensure that the data has this shape.
	Eval(input tuple.Value) (tuple.Value, error)
}

// ExpressionToEvaluator takes one of the Expression structs that result
// from parsing a BQL Expression (see parser/ast.go) and turns it into
// an Evaluator that can be used to evaluate an expression given a particular
// input Value.
func ExpressionToEvaluator(ast interface{}, reg udf.FunctionRegistry) (Evaluator, error) {
	switch obj := ast.(type) {
	case parser.RowMeta:
		// construct a key for reading as used in setMetadata() for writing
		metaKey := fmt.Sprintf("%s:meta:%s", obj.Relation, obj.MetaType)
		if obj.MetaType == parser.TimestampMeta {
			return &timestampCast{&PathAccess{metaKey}}, nil
		}
	case parser.RowValue:
		return &PathAccess{obj.Relation + "." + obj.Column}, nil
	case parser.AliasAST:
		return ExpressionToEvaluator(obj.Expr, reg)
	case parser.NumericLiteral:
		return &IntConstant{obj.Value}, nil
	case parser.FloatLiteral:
		return &FloatConstant{obj.Value}, nil
	case parser.BoolLiteral:
		return &BoolConstant{obj.Value}, nil
	case parser.BinaryOpAST:
		// recurse
		left, err := ExpressionToEvaluator(obj.Left, reg)
		if err != nil {
			return nil, err
		}
		right, err := ExpressionToEvaluator(obj.Right, reg)
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
	case parser.FuncAppAST:
		// lookup function in function registry
		// (the registry will decide if the requested function
		// is callable with the given number of arguments).
		fName := string(obj.Function)
		f, err := reg.Lookup(fName, len(obj.Expressions))
		if err != nil {
			return nil, err
		}
		// compute child Evaluators
		evals := make([]Evaluator, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			eval, err := ExpressionToEvaluator(ast, reg)
			if err != nil {
				return nil, err
			}
			evals[i] = eval
		}
		return FuncApp(fName, f, reg.Context(), evals), nil
	case parser.Wildcard:
		return &Wildcard{}, nil
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
	aMap, err := tuple.AsMap(input)
	if err != nil {
		return nil, err
	}
	return aMap.Get(fa.path)
}

// TODO this should probably be a general cast
type timestampCast struct {
	underlying Evaluator
}

func (t *timestampCast) Eval(input tuple.Value) (tuple.Value, error) {
	val, err := t.underlying.Eval(input)
	if err != nil {
		return nil, err
	}
	if val.Type() != tuple.TypeTimestamp {
		return nil, fmt.Errorf("value %v was %T, not Time", val, val)
	}
	return val, nil
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
	negBool, err := tuple.AsBool(neg)
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
			l, _ := tuple.AsInt(leftVal)
			r, _ := tuple.AsFloat(rightVal)
			// convert left to float to get 2 == 2.0
			eq = float64(l) == r
		} else if leftType == tuple.TypeFloat && rightType == tuple.TypeInt {
			l, _ := tuple.AsFloat(leftVal)
			r, _ := tuple.AsInt(rightVal)
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
				l, _ := tuple.AsInt(leftVal)
				r, _ := tuple.AsInt(rightVal)
				retVal = l < r
			case tuple.TypeFloat:
				l, _ := tuple.AsFloat(leftVal)
				r, _ := tuple.AsFloat(rightVal)
				retVal = l < r
			case tuple.TypeString:
				l, _ := tuple.AsString(leftVal)
				r, _ := tuple.AsString(rightVal)
				retVal = l < r
			case tuple.TypeBool:
				l, _ := tuple.AsBool(leftVal)
				r, _ := tuple.AsBool(rightVal)
				retVal = (l == false) && (r == true)
			case tuple.TypeTimestamp:
				l, _ := tuple.AsTimestamp(leftVal)
				r, _ := tuple.AsTimestamp(rightVal)
				retVal = l.Before(r)
			}
			return retVal, nil
		} else if leftType == tuple.TypeInt && rightType == tuple.TypeFloat {
			// left is integer
			l, _ := tuple.AsInt(leftVal)
			// right is float; also convert left to float to avoid overflow
			r, _ := tuple.AsFloat(rightVal)
			return float64(l) < r, nil
		} else if leftType == tuple.TypeFloat && rightType == tuple.TypeInt {
			// left is float
			l, _ := tuple.AsFloat(leftVal)
			// right is int; convert right to float to avoid overflow
			r, _ := tuple.AsInt(rightVal)
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
			l, _ := tuple.AsInt(leftVal)
			r, _ := tuple.AsInt(rightVal)
			return tuple.Int(nbo.intOp(l, r)), nil
		case tuple.TypeFloat:
			l, _ := tuple.AsFloat(leftVal)
			r, _ := tuple.AsFloat(rightVal)
			return tuple.Float(nbo.floatOp(l, r)), nil
		}
	} else if leftType == tuple.TypeInt && rightType == tuple.TypeFloat {
		// left is integer
		l, _ := tuple.AsInt(leftVal)
		// right is float; also convert left to float, possibly losing precision
		r, _ := tuple.AsFloat(rightVal)
		return tuple.Float(nbo.floatOp(float64(l), r)), nil
	} else if leftType == tuple.TypeFloat && rightType == tuple.TypeInt {
		// left is float
		l, _ := tuple.AsFloat(leftVal)
		// right is int; convert right to float, possibly losing precision
		r, _ := tuple.AsInt(rightVal)
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

/// Function Evaluation

type funcApp struct {
	name        string
	fVal        reflect.Value
	params      []Evaluator
	paramValues []reflect.Value
}

func (f *funcApp) Eval(input tuple.Value) (v tuple.Value, err error) {
	// catch panic (e.g., in called function)
	defer func() {
		if r := recover(); r != nil {
			v = nil
			err = fmt.Errorf("evaluating '%s' paniced: %s", f.name, r)
		}
	}()
	// evaluate all the parameters and store the results
	for i, param := range f.params {
		value, err := param.Eval(input)
		if err != nil {
			return nil, err
		}
		f.paramValues[i+1] = reflect.ValueOf(value)
	}
	// evaluate the function
	results := f.fVal.Call(f.paramValues)
	// check results
	if len(results) != 2 {
		return nil, fmt.Errorf("function %s returned %d results, not 2",
			f.name, len(results))
	}
	resultVal, errVal := results[0], results[1]
	if !errVal.IsNil() {
		err := errVal.Interface().(error)
		return nil, err
	}
	result := resultVal.Interface().(tuple.Value)
	return result, nil
}

// FuncApp represents evaluation of a function on a number
// of parameters that are expressions over an input Value.
func FuncApp(name string, f udf.VarParamFun, ctx *core.Context, params []Evaluator) Evaluator {
	fVal := reflect.ValueOf(f)
	paramValues := make([]reflect.Value, len(params)+1)
	paramValues[0] = reflect.ValueOf(ctx)
	return &funcApp{name, fVal, params, paramValues}
}

// Wildcard only works on Maps, assumes that the elements which do not contain
// ":meta:" are also Maps and pulls them up one level, so
//   {"a": {"x": ...}, "a:meta:ts": ..., "b": {"y": ..., "z": ...}}
// becomes
//   {"x": ..., "y": ..., "z": ...}.
// If there are keys appearing in multiple top-level Maps, then only one
// of them will appear in the output, but it is undefined which.
type Wildcard struct{}

func (w *Wildcard) Eval(input tuple.Value) (tuple.Value, error) {
	aMap, err := tuple.AsMap(input)
	if err != nil {
		return nil, err
	}
	output := tuple.Map{}
	for alias, subElement := range aMap {
		if strings.Contains(alias, ":meta:") {
			continue
		}
		subMap, err := tuple.AsMap(subElement)
		if err != nil {
			return nil, err
		}
		for key, value := range subMap {
			output[key] = value
		}
	}
	return output, nil
}

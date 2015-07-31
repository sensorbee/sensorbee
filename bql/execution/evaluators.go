package execution

import (
	"fmt"
	"math"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
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
	Eval(input data.Value) (data.Value, error)
}

// EvaluateFoldable evaluates a foldable expression, i.e., one that
// is independent from the input row. Note that foldable is not
// necessarily equivalent to constant (e.g., the expression `random()`
// is foldable, but not constant), and also note that this function
// should not be used for frequent evaluation of the same expression
// due to performance reasons.
func EvaluateFoldable(expr parser.Expression, reg udf.FunctionRegistry) (data.Value, error) {
	if !expr.Foldable() {
		return nil, fmt.Errorf("expression is not foldable: %s", expr)
	}
	flatExpr, err := ParserExprToFlatExpr(expr, reg)
	if err != nil {
		return nil, err
	}
	evaluator, err := ExpressionToEvaluator(flatExpr, reg)
	if err != nil {
		return nil, err
	}
	return evaluator.Eval(nil)
}

// ExpressionToEvaluator takes one of the Expression structs that result
// from parsing a BQL Expression (see parser/ast.go) and turns it into
// an Evaluator that can be used to evaluate an expression given a particular
// input Value.
func ExpressionToEvaluator(ast FlatExpression, reg udf.FunctionRegistry) (Evaluator, error) {
	switch obj := ast.(type) {
	case RowMeta:
		// construct a key for reading as used in setMetadata() for writing
		metaKey := fmt.Sprintf("%s:meta:%s", obj.Relation, obj.MetaType)
		if obj.MetaType == parser.TimestampMeta {
			return &timestampCast{&PathAccess{metaKey}}, nil
		}
	case StmtMeta:
		// construct a key for reading as used in setMetadata() for writing
		metaKey := fmt.Sprintf(":meta:%s", obj.MetaType)
		if obj.MetaType == parser.NowMeta {
			return &timestampCast{&PathAccess{metaKey}}, nil
		}
	case RowValue:
		return &PathAccess{obj.Relation + "." + obj.Column}, nil
	case AggInputRef:
		return &PathAccess{obj.Ref}, nil
	case NullLiteral:
		return &NullConstant{}, nil
	case NumericLiteral:
		return &IntConstant{obj.Value}, nil
	case FloatLiteral:
		return &FloatConstant{obj.Value}, nil
	case BoolLiteral:
		return &BoolConstant{obj.Value}, nil
	case StringLiteral:
		return &StringConstant{obj.Value}, nil
	case BinaryOpAST:
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
			return &Or{bo}, nil
		case parser.And:
			return &And{bo}, nil
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
		case parser.Concat:
			return &Concat{bo}, nil
		case parser.Is:
			// at the moment there is only NULL allowed after IS,
			// but maybe we want to allow other types later on
			if obj.Right == (NullLiteral{}) {
				return IsNull(left), nil
			}
		case parser.IsNot:
			// at the moment there is only NULL allowed after IS NOT,
			// but maybe we want to allow other types later on
			if obj.Right == (NullLiteral{}) {
				return Not(IsNull(left)), nil
			}
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
	case UnaryOpAST:
		// recurse
		expr, err := ExpressionToEvaluator(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		// combine child with the correct operator
		switch obj.Op {
		default:
			err := fmt.Errorf("don't know how to evaluate unary operation %v", obj.Op)
			return nil, err
		case parser.Not:
			return Not(expr), nil
		case parser.UnaryMinus:
			// implement negation as multiplication with -1
			bo := binOp{expr, &IntConstant{-1}}
			return Multiply(bo), nil
		}
	case TypeCastAST:
		// recurse
		expr, err := ExpressionToEvaluator(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return TypeCast(expr, obj.Target)
	case FuncAppAST:
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
	case ArrayAST:
		// compute child Evaluators
		evals := make([]Evaluator, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			eval, err := ExpressionToEvaluator(ast, reg)
			if err != nil {
				return nil, err
			}
			evals[i] = eval
		}
		return ArrayBuilder(evals), nil
	case MapAST:
		// compute child Evaluators
		names := make([]string, len(obj.Entries))
		evals := make([]Evaluator, len(obj.Entries))
		for i, pair := range obj.Entries {
			eval, err := ExpressionToEvaluator(pair.Value, reg)
			if err != nil {
				return nil, err
			}
			evals[i] = eval
			names[i] = pair.Key
		}
		return MapBuilder(names, evals)
	case WildcardAST:
		return &Wildcard{obj.Relation}, nil
	}
	err := fmt.Errorf("don't know how to evaluate type %#v", ast)
	return nil, err
}

// NullConstant always returns the same null value, independent
// of the input.
type NullConstant struct {
}

func (n *NullConstant) Eval(input data.Value) (data.Value, error) {
	return data.Null{}, nil
}

// IntConstant always returns the same integer value, independent
// of the input.
type IntConstant struct {
	value int64
}

func (i *IntConstant) Eval(input data.Value) (data.Value, error) {
	return data.Int(i.value), nil
}

// FloatConstant always returns the same float value, independent
// of the input.
type FloatConstant struct {
	value float64
}

func (f *FloatConstant) Eval(input data.Value) (data.Value, error) {
	return data.Float(f.value), nil
}

// BoolConstant always returns the same boolean value, independent
// of the input.
type BoolConstant struct {
	value bool
}

func (b *BoolConstant) Eval(input data.Value) (data.Value, error) {
	return data.Bool(b.value), nil
}

// StringConstant always returns the same string value, independent
// of the input.
type StringConstant struct {
	value string
}

func (s *StringConstant) Eval(input data.Value) (data.Value, error) {
	return data.String(s.value), nil
}

// PathAccess only works for maps and returns the Value at the given
// JSON path.
type PathAccess struct {
	path string
}

func (fa *PathAccess) Eval(input data.Value) (data.Value, error) {
	aMap, err := data.AsMap(input)
	if err != nil {
		return nil, err
	}
	return aMap.Get(fa.path)
}

type typeCast struct {
	underlying Evaluator
	converter  func(data.Value) (data.Value, error)
}

func (t *typeCast) Eval(input data.Value) (data.Value, error) {
	val, err := t.underlying.Eval(input)
	if err != nil {
		return nil, err
	}
	return t.converter(val)
}

func TypeCast(e Evaluator, t parser.Type) (Evaluator, error) {
	switch t {
	case parser.Bool:
		conv := func(v data.Value) (data.Value, error) {
			x, err := data.ToBool(v)
			if err != nil {
				return nil, err
			}
			return data.Bool(x), nil
		}
		return &typeCast{e, conv}, nil
	case parser.Int:
		conv := func(v data.Value) (data.Value, error) {
			x, err := data.ToInt(v)
			if err != nil {
				return nil, err
			}
			return data.Int(x), nil
		}
		return &typeCast{e, conv}, nil
	case parser.Float:
		conv := func(v data.Value) (data.Value, error) {
			x, err := data.ToFloat(v)
			if err != nil {
				return nil, err
			}
			return data.Float(x), nil
		}
		return &typeCast{e, conv}, nil
	case parser.String:
		conv := func(v data.Value) (data.Value, error) {
			x, err := data.ToString(v)
			if err != nil {
				return nil, err
			}
			return data.String(x), nil
		}
		return &typeCast{e, conv}, nil
	case parser.Blob:
		conv := func(v data.Value) (data.Value, error) {
			x, err := data.ToBlob(v)
			if err != nil {
				return nil, err
			}
			return data.Blob(x), nil
		}
		return &typeCast{e, conv}, nil
	case parser.Timestamp:
		conv := func(v data.Value) (data.Value, error) {
			x, err := data.ToTimestamp(v)
			if err != nil {
				return nil, err
			}
			return data.Timestamp(x), nil
		}
		return &typeCast{e, conv}, nil
	}
	return nil, fmt.Errorf("no converter for type %s known", t)
}

// TODO this should probably be a general cast
type timestampCast struct {
	underlying Evaluator
}

func (t *timestampCast) Eval(input data.Value) (data.Value, error) {
	val, err := t.underlying.Eval(input)
	if err != nil {
		return nil, err
	}
	if val.Type() != data.TypeTimestamp {
		return nil, fmt.Errorf("value %v was %T, not Time", val, val)
	}
	return val, nil
}

type binOp struct {
	left  Evaluator
	right Evaluator
}

func (bo *binOp) evalLeftAndRight(input data.Value) (data.Value, data.Value, error) {
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

type Or struct {
	binOp
}

func (o *Or) Eval(input data.Value) (data.Value, error) {
	leftRes, err := o.left.Eval(input)
	if err != nil {
		return nil, err
	}
	if leftRes.Type() == data.TypeNull {
		// if we continue here, the left side was NULL
		rightRes, err := o.right.Eval(input)
		if err != nil {
			return nil, err
		}
		if rightRes.Type() == data.TypeNull {
			// NULL OR NULL => NULL
			return data.Null{}, nil
		}
		rightBool, err := data.ToBool(rightRes)
		if err != nil {
			return nil, err
		}
		if rightBool {
			// NULL OR true => true
			return data.Bool(true), nil
		}
		// NULL OR false => NULL
		return data.Null{}, nil
	} else {
		leftBool, err := data.ToBool(leftRes)
		if err != nil {
			return nil, err
		}
		// support early return if the left side is true
		if leftBool {
			// true OR true => true
			// true OR false => true
			// true OR NULL => true
			return data.Bool(true), nil
		}
		// if we continue here, the left side was false
		rightRes, err := o.right.Eval(input)
		if err != nil {
			return nil, err
		}
		if rightRes.Type() == data.TypeNull {
			// false OR NULL => NULL
			return data.Null{}, nil
		}
		rightBool, err := data.ToBool(rightRes)
		if err != nil {
			return nil, err
		}
		// false OR true => true
		// false OR false => false
		return data.Bool(rightBool), nil
	}
}

type And struct {
	binOp
}

func (a *And) Eval(input data.Value) (data.Value, error) {
	leftRes, err := a.left.Eval(input)
	if err != nil {
		return nil, err
	}
	if leftRes.Type() == data.TypeNull {
		// if we continue here, the left side was NULL
		rightRes, err := a.right.Eval(input)
		if err != nil {
			return nil, err
		}
		if rightRes.Type() == data.TypeNull {
			// NULL AND NULL => NULL
			return data.Null{}, nil
		}
		rightBool, err := data.ToBool(rightRes)
		if err != nil {
			return nil, err
		}
		if rightBool {
			// NULL AND true => NULL
			return data.Null{}, nil
		}
		// NULL AND false => false
		return data.Bool(false), nil
	} else {
		leftBool, err := data.ToBool(leftRes)
		if err != nil {
			return nil, err
		}
		// support early return if the left side is false
		if !leftBool {
			// false AND true => false
			// false AND false => false
			// false AND NULL => false
			return data.Bool(false), nil
		}
		// if we continue here, the left side was true
		rightRes, err := a.right.Eval(input)
		if err != nil {
			return nil, err
		}
		if rightRes.Type() == data.TypeNull {
			// true AND NULL => NULL
			return data.Null{}, nil
		}
		rightBool, err := data.ToBool(rightRes)
		if err != nil {
			return nil, err
		}
		// true AND true => true
		// true AND false => false
		return data.Bool(rightBool), nil
	}
}

/// A Unary Logical Operation

type not struct {
	neg Evaluator
}

func (n *not) Eval(input data.Value) (data.Value, error) {
	neg, err := n.neg.Eval(input)
	if err != nil {
		return nil, err
	}
	// NULL propagation
	if neg.Type() == data.TypeNull {
		return data.Null{}, nil
	}
	negBool, err := data.AsBool(neg)
	if err != nil {
		return nil, err
	}
	return data.Bool(!negBool), nil
}

func Not(e Evaluator) Evaluator {
	return &not{e}
}

/// Binary Comparison Operations

// compBinOp provides functionality for comparing two values
type compBinOp struct {
	binOp
	cmpOp func(data.Value, data.Value) (bool, error)
}

func (cbo *compBinOp) Eval(input data.Value) (data.Value, error) {
	leftVal, rightVal, err := cbo.evalLeftAndRight(input)
	if err != nil {
		return nil, err
	}
	// NULL propagation
	if leftVal.Type() == data.TypeNull || rightVal.Type() == data.TypeNull {
		return data.Null{}, nil
	}
	res, err := cbo.cmpOp(leftVal, rightVal)
	if err != nil {
		return nil, err
	}
	return data.Bool(res), nil
}

func Equal(bo binOp) Evaluator {
	cmpOp := func(leftVal data.Value, rightVal data.Value) (bool, error) {
		return data.HashEqual(leftVal, rightVal), nil

	}
	return &compBinOp{bo, cmpOp}
}

func Less(bo binOp) Evaluator {
	cmpOp := func(leftVal data.Value, rightVal data.Value) (bool, error) {
		leftType := leftVal.Type()
		rightType := rightVal.Type()
		stdErr := fmt.Errorf("cannot compare %T and %T", leftVal, rightVal)
		if leftType == rightType {
			retVal := false
			switch leftType {
			default:
				return false, stdErr
			case data.TypeInt:
				l, _ := data.AsInt(leftVal)
				r, _ := data.AsInt(rightVal)
				retVal = l < r
			case data.TypeFloat:
				l, _ := data.AsFloat(leftVal)
				r, _ := data.AsFloat(rightVal)
				retVal = l < r
			case data.TypeString:
				l, _ := data.AsString(leftVal)
				r, _ := data.AsString(rightVal)
				retVal = l < r
			case data.TypeBool:
				l, _ := data.AsBool(leftVal)
				r, _ := data.AsBool(rightVal)
				retVal = (l == false) && (r == true)
			case data.TypeTimestamp:
				l, _ := data.AsTimestamp(leftVal)
				r, _ := data.AsTimestamp(rightVal)
				retVal = l.Before(r)
			}
			return retVal, nil
		} else if leftType == data.TypeInt && rightType == data.TypeFloat {
			// left is integer
			l, _ := data.AsInt(leftVal)
			// right is float; also convert left to float to avoid overflow
			r, _ := data.AsFloat(rightVal)
			return float64(l) < r, nil
		} else if leftType == data.TypeFloat && rightType == data.TypeInt {
			// left is float
			l, _ := data.AsFloat(leftVal)
			// right is int; convert right to float to avoid overflow
			r, _ := data.AsInt(rightVal)
			return l < float64(r), nil
		}
		return false, stdErr
	}
	return &compBinOp{bo, cmpOp}
}

func LessOrEqual(bo binOp) Evaluator {
	return &Or{binOp{Less(bo), Equal(bo)}}
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

/// A Unary Comparison Operation

type isNull struct {
	val Evaluator
}

func (n *isNull) Eval(input data.Value) (data.Value, error) {
	val, err := n.val.Eval(input)
	if err != nil {
		return nil, err
	}
	return data.Bool(val.Type() == data.TypeNull), nil
}

func IsNull(e Evaluator) Evaluator {
	return &isNull{e}
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

func (nbo *numBinOp) Eval(input data.Value) (v data.Value, err error) {
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
	// NULL propagation
	if leftType == data.TypeNull || rightType == data.TypeNull {
		return data.Null{}, nil
	}
	stdErr := fmt.Errorf("cannot %s %T and %T", nbo.verb, leftVal, rightVal)
	// if we have same types (both int64 or both float64, apply
	// the corresponding operation)
	if leftType == rightType {
		switch leftType {
		default:
			return nil, stdErr
		case data.TypeInt:
			l, _ := data.AsInt(leftVal)
			r, _ := data.AsInt(rightVal)
			return data.Int(nbo.intOp(l, r)), nil
		case data.TypeFloat:
			l, _ := data.AsFloat(leftVal)
			r, _ := data.AsFloat(rightVal)
			return data.Float(nbo.floatOp(l, r)), nil
		}
	} else if leftType == data.TypeInt && rightType == data.TypeFloat {
		// left is integer
		l, _ := data.AsInt(leftVal)
		// right is float; also convert left to float, possibly losing precision
		r, _ := data.AsFloat(rightVal)
		return data.Float(nbo.floatOp(float64(l), r)), nil
	} else if leftType == data.TypeFloat && rightType == data.TypeInt {
		// left is float
		l, _ := data.AsFloat(leftVal)
		// right is int; convert right to float, possibly losing precision
		r, _ := data.AsInt(rightVal)
		return data.Float(nbo.floatOp(l, float64(r))), nil
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

/// Other Binary Operations

type Concat struct {
	binOp
}

func (nbo *Concat) Eval(input data.Value) (v data.Value, err error) {
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
	// NULL propagation
	if leftVal.Type() == data.TypeNull || rightVal.Type() == data.TypeNull {
		return data.Null{}, nil
	}
	leftString, err := data.ToString(leftVal)
	if err != nil {
		return nil, err
	}
	rightString, err := data.ToString(rightVal)
	if err != nil {
		return nil, err
	}
	return data.String(leftString + rightString), nil
}

/// Function Evaluation

type funcApp struct {
	name        string
	fVal        reflect.Value
	params      []Evaluator
	paramValues []reflect.Value
}

func (f *funcApp) Eval(input data.Value) (v data.Value, err error) {
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
	result := resultVal.Interface().(data.Value)
	return result, nil
}

// FuncApp represents evaluation of a function on a number
// of parameters that are expressions over an input Value.
func FuncApp(name string, f udf.UDF, ctx *core.Context, params []Evaluator) Evaluator {
	fVal := reflect.ValueOf(f.Call)
	paramValues := make([]reflect.Value, len(params)+1)
	paramValues[0] = reflect.ValueOf(ctx)
	return &funcApp{name, fVal, params, paramValues}
}

/// JSON-like data structures

type arrayBuilder struct {
	elems []Evaluator
}

func (a *arrayBuilder) Eval(input data.Value) (v data.Value, err error) {
	results := make([]data.Value, len(a.elems))
	// evaluate all the parameters and store the results
	for i, elem := range a.elems {
		value, err := elem.Eval(input)
		if err != nil {
			return nil, err
		}
		results[i] = value
	}
	return data.Array(results), nil
}

func ArrayBuilder(elems []Evaluator) Evaluator {
	return &arrayBuilder{elems}
}

type mapBuilder struct {
	names []string
	elems []Evaluator
}

func (m *mapBuilder) Eval(input data.Value) (v data.Value, err error) {
	results := make(data.Map, len(m.elems))
	// evaluate all the parameters and store the results
	for i, elem := range m.elems {
		value, err := elem.Eval(input)
		if err != nil {
			return nil, err
		}
		results[m.names[i]] = value
	}
	return results, nil
}

func MapBuilder(names []string, elems []Evaluator) (Evaluator, error) {
	if len(names) != len(elems) {
		return nil, fmt.Errorf("number of keys and values does not match")
	}
	return &mapBuilder{names, elems}, nil
}

// Wildcard only works on Maps, assumes that the elements which do not contain
// ":meta:" are also Maps and pulls them up one level, so
//   {"a": {"x": ...}, "a:meta:ts": ..., "b": {"y": ..., "z": ...}}
// becomes
//   {"x": ..., "y": ..., "z": ...}.
// If there are keys appearing in multiple top-level Maps, then only one
// of them will appear in the output, but it is undefined which.
type Wildcard struct {
	Relation string
}

func (w *Wildcard) Eval(input data.Value) (data.Value, error) {
	aMap, err := data.AsMap(input)
	if err != nil {
		return nil, err
	}
	output := data.Map{}
	if w.Relation != "" {
		// if we have t:*, pick only the items in this submap
		subElement, exists := aMap[w.Relation]
		if !exists {
			return nil, fmt.Errorf("there is no entry with key '%s'", w.Relation)
		}
		subMap, err := data.AsMap(subElement)
		if err != nil {
			return nil, err
		}
		for key, value := range subMap {
			output[key] = value
		}
	} else {
		// if we have *, take items from all submaps
		for alias, subElement := range aMap {
			if strings.Contains(alias, ":meta:") {
				continue
			}
			subMap, err := data.AsMap(subElement)
			if err != nil {
				return nil, err
			}
			for key, value := range subMap {
				output[key] = value
			}
		}
	}
	return output, nil
}

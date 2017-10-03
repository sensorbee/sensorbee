package execution

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

// An Evaluator represents an expression such as `colX + 2` or
// `t1:col AND t2:col` and can be evaluated, given the actual data
// contained in one row.
type Evaluator interface {
	// Eval evaluates the expression that this Evaluator represents
	// on the given input data. Note that in order to deal with joins and
	// meta information such as timestamps properly, the input data must have
	// the shape:
	//   {"alias_1": {"col_0": ..., "col_1": ...},
	//    "alias_1:meta:x": (meta datum "x" for alias_1's row),
	//    "alias_2": {"col_0": ..., "col_1": ...},
	//    "alias_2:meta:x": (meta datum "x" for alias_2's row),
	//    ...}
	// and every caller (in particular all execution plans)
	// must ensure that the data has this shape even if there's only one input
	// stream.
	//
	// Eval must NOT modify the input.
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

// EvaluateOnInput evaluates a (not necessarily foldable)
// expression, given a Map that represents a row of data.
func EvaluateOnInput(expr parser.Expression, input data.Value, reg udf.FunctionRegistry) (data.Value, error) {
	flatExpr, err := ParserExprToFlatExpr(expr, reg)
	if err != nil {
		return nil, err
	}
	evaluator, err := ExpressionToEvaluator(flatExpr, reg)
	if err != nil {
		return nil, err
	}
	return evaluator.Eval(input)
}

// ExpressionToEvaluator takes one of the Expression structs that result
// from parsing a BQL Expression (see parser/ast.go) and turns it into
// an Evaluator that can be used to evaluate an expression given a particular
// input Value.
func ExpressionToEvaluator(ast FlatExpression, reg udf.FunctionRegistry) (Evaluator, error) {
	switch obj := ast.(type) {
	case rowMeta:
		// construct a key for reading as used in setMetadata() for writing
		metaKey := fmt.Sprintf(`["%s:meta:%s"]`, obj.Relation, obj.MetaType)
		if obj.MetaType == parser.TimestampMeta {
			pa, err := newPathAccess(metaKey)
			if err != nil {
				return nil, err
			}
			return &timestampCast{pa}, nil
		}
	case stmtMeta:
		// construct a key for reading as used in setMetadata() for writing
		metaKey := fmt.Sprintf(`[":meta:%s"]`, obj.MetaType)
		if obj.MetaType == parser.NowMeta {
			pa, err := newPathAccess(metaKey)
			if err != nil {
				return nil, err
			}
			return &timestampCast{pa}, nil
		}
	case rowValue:
		path := obj.Column
		if obj.Relation != "" {
			if strings.HasPrefix(path, "[") {
				path = obj.Relation + path
			} else {
				path = obj.Relation + "." + path
			}
		}
		return newPathAccess(path)
	case aggInputRef:
		return newPathAccess(obj.Ref)
	case nullLiteral:
		return &nullConstant{}, nil
	case numericLiteral:
		return &intConstant{obj.Value}, nil
	case floatLiteral:
		return &floatConstant{obj.Value}, nil
	case boolLiteral:
		return &boolConstant{obj.Value}, nil
	case stringLiteral:
		return &stringConstant{obj.Value}, nil
	case binaryOpAST:
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
			return &or{bo}, nil
		case parser.And:
			return &and{bo}, nil
		case parser.Equal:
			return newEqual(bo), nil
		case parser.Less:
			return newLess(bo), nil
		case parser.LessOrEqual:
			return newLessOrEqual(bo), nil
		case parser.Greater:
			return newGreater(bo), nil
		case parser.GreaterOrEqual:
			return newGreaterOrnewEqual(bo), nil
		case parser.NotEqual:
			return newNot(newEqual(bo)), nil
		case parser.Concat:
			return &concat{bo}, nil
		case parser.Is:
			// at the moment there is only NULL allowed after IS,
			// but maybe we want to allow other types later on
			if obj.Right == (nullLiteral{}) {
				return newIsNull(left), nil
			}
		case parser.IsNot:
			// at the moment there is only NULL allowed after IS NOT,
			// but maybe we want to allow other types later on
			if obj.Right == (nullLiteral{}) {
				return newNot(newIsNull(left)), nil
			}
		case parser.Plus:
			return newPlus(bo), nil
		case parser.Minus:
			return newMinus(bo), nil
		case parser.Multiply:
			return newMultiply(bo), nil
		case parser.Divide:
			return newDivide(bo), nil
		case parser.Modulo:
			return newModulo(bo), nil
		}
	case unaryOpAST:
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
			return newNot(expr), nil
		case parser.UnaryMinus:
			// implement negation as multiplication with -1
			bo := binOp{expr, &intConstant{-1}}
			return newMultiply(bo), nil
		}
	case missing:
		// recurse
		expr, err := ExpressionToEvaluator(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return newMissingPathCheck(expr, obj.Not)
	case typeCastAST:
		// recurse
		expr, err := ExpressionToEvaluator(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return newTypeCast(expr, obj.Target)
	case funcAppSelectorAST:
		// recurse
		expr, err := ExpressionToEvaluator(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		funcEval := expr.(*funcApp) // snip type error check
		return FuncAppSelector(funcEval, obj.Selector)
	case funcAppAST:
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
	case aggregateInputSorter:
		return newSortedInputAggFuncApp(obj.funcAppAST, obj.ID, obj.Ordering, reg)
	case arrayAST:
		// compute child Evaluators
		evals := make([]Evaluator, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			eval, err := ExpressionToEvaluator(ast, reg)
			if err != nil {
				return nil, err
			}
			evals[i] = eval
		}
		return newArrayBuilder(evals), nil
	case mapAST:
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
		return newMapBuilder(names, evals)
	case caseAST:
		// compute the Evaluator for the thing we match against
		ref, err := ExpressionToEvaluator(obj.Reference, reg)
		if err != nil {
			return nil, err
		}
		// compute Evaluators for all possible WHEN-THEN clauses
		whens := make([]Evaluator, len(obj.Checks))
		thens := make([]Evaluator, len(obj.Checks))
		for i, pair := range obj.Checks {
			eval, err := ExpressionToEvaluator(pair.When, reg)
			if err != nil {
				return nil, err
			}
			whens[i] = eval
			eval, err = ExpressionToEvaluator(pair.Then, reg)
			if err != nil {
				return nil, err
			}
			thens[i] = eval
		}
		// compute the Evaluator for the default value (if nothing matches)
		def, err := ExpressionToEvaluator(obj.Default, reg)
		if err != nil {
			return nil, err
		}
		return newCaseBuilder(ref, whens, thens, def)
	case wildcardAST:
		return &wildcard{obj.Relation}, nil
	}
	err := fmt.Errorf("don't know how to evaluate type %#v", ast)
	return nil, err
}

// nullConstant always returns the same null value, independent
// of the input.
type nullConstant struct {
}

func (n *nullConstant) Eval(input data.Value) (data.Value, error) {
	return data.Null{}, nil
}

// intConstant always returns the same integer value, independent
// of the input.
type intConstant struct {
	value int64
}

func (i *intConstant) Eval(input data.Value) (data.Value, error) {
	return data.Int(i.value), nil
}

// floatConstant always returns the same float value, independent
// of the input.
type floatConstant struct {
	value float64
}

func (f *floatConstant) Eval(input data.Value) (data.Value, error) {
	return data.Float(f.value), nil
}

// boolConstant always returns the same boolean value, independent
// of the input.
type boolConstant struct {
	value bool
}

func (b *boolConstant) Eval(input data.Value) (data.Value, error) {
	return data.Bool(b.value), nil
}

// stringConstant always returns the same string value, independent
// of the input.
type stringConstant struct {
	value string
}

func (s *stringConstant) Eval(input data.Value) (data.Value, error) {
	return data.String(s.value), nil
}

// pathAccess only works for maps and returns the Value at the given
// JSON path.
type pathAccess struct {
	path data.Path
}

func (fa *pathAccess) Eval(input data.Value) (data.Value, error) {
	aMap, err := data.AsMap(input)
	if err != nil {
		return nil, err
	}
	return aMap.Get(fa.path)
}

func newPathAccess(s string) (Evaluator, error) {
	path, err := data.CompilePath(s)
	if err != nil {
		return nil, err
	}
	return &pathAccess{path}, nil
}

type missingPathCheck struct {
	eval   pathAccess
	negate bool
}

func (m *missingPathCheck) Eval(input data.Value) (data.Value, error) {
	if input.Type() != data.TypeMap {
		return nil, fmt.Errorf("expected Map for IS MISSING check, not %s", input.Type())
	}
	// we assume that if there was any error, the value was missing
	if _, err := m.eval.Eval(input); err != nil {
		return data.Bool(true != m.negate), nil
	}
	return data.Bool(false != m.negate), nil
}

func newMissingPathCheck(eval Evaluator, negate bool) (Evaluator, error) {
	pa, ok := eval.(*pathAccess)
	if !ok {
		return nil, fmt.Errorf("expected pathAccess before IS [NOT] MISSING, not %v", eval)
	}
	return &missingPathCheck{*pa, negate}, nil
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
	// null propagation
	if val.Type() == data.TypeNull {
		return data.Null{}, nil
	}
	return t.converter(val)
}

func newTypeCast(e Evaluator, t parser.Type) (Evaluator, error) {
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

type or struct {
	binOp
}

func (o *or) Eval(input data.Value) (data.Value, error) {
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
		rightBool, err := data.AsBool(rightRes)
		if err != nil {
			return nil, err
		}
		if rightBool {
			// NULL OR true => true
			return data.Bool(true), nil
		}
		// NULL OR false => NULL
		return data.Null{}, nil
	}
	// indent the block below for symmetry reasons
	{
		leftBool, err := data.AsBool(leftRes)
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
		rightBool, err := data.AsBool(rightRes)
		if err != nil {
			return nil, err
		}
		// false OR true => true
		// false OR false => false
		return data.Bool(rightBool), nil
	}
}

type and struct {
	binOp
}

func (a *and) Eval(input data.Value) (data.Value, error) {
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
		rightBool, err := data.AsBool(rightRes)
		if err != nil {
			return nil, err
		}
		if rightBool {
			// NULL AND true => NULL
			return data.Null{}, nil
		}
		// NULL AND false => false
		return data.Bool(false), nil
	}
	// indent the block below for symmetry reasons
	{
		leftBool, err := data.AsBool(leftRes)
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
		rightBool, err := data.AsBool(rightRes)
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

func newNot(e Evaluator) Evaluator {
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

func newEqual(bo binOp) Evaluator {
	cmpOp := func(leftVal data.Value, rightVal data.Value) (bool, error) {
		return data.Equal(leftVal, rightVal), nil

	}
	return &compBinOp{bo, cmpOp}
}

func newLess(bo binOp) Evaluator {
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

func newLessOrEqual(bo binOp) Evaluator {
	return &or{binOp{newLess(bo), newEqual(bo)}}
}

func newGreater(bo binOp) Evaluator {
	return newNot(newLessOrEqual(bo))
}

func newGreaterOrnewEqual(bo binOp) Evaluator {
	return newNot(newLess(bo))
}

func newNotEqual(bo binOp) Evaluator {
	return newNot(newEqual(bo))
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

func newIsNull(e Evaluator) Evaluator {
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

func newPlus(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a + b
	}
	floatOp := func(a, b float64) float64 {
		return a + b
	}
	return &numBinOp{bo, "add", intOp, floatOp}
}

func newMinus(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a - b
	}
	floatOp := func(a, b float64) float64 {
		return a - b
	}
	return &numBinOp{bo, "subtract", intOp, floatOp}
}

func newMultiply(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a * b
	}
	floatOp := func(a, b float64) float64 {
		return a * b
	}
	return &numBinOp{bo, "multiply", intOp, floatOp}
}

func newDivide(bo binOp) Evaluator {
	// we do not check for overflows
	intOp := func(a, b int64) int64 {
		return a / b
	}
	floatOp := func(a, b float64) float64 {
		return a / b
	}
	return &numBinOp{bo, "divide", intOp, floatOp}
}

func newModulo(bo binOp) Evaluator {
	intOp := func(a, b int64) int64 {
		return a % b
	}
	floatOp := func(a, b float64) float64 {
		return math.Mod(a, b)
	}
	return &numBinOp{bo, "compute modulo for", intOp, floatOp}
}

/// Other Binary Operations

type concat struct {
	binOp
}

func (nbo *concat) Eval(input data.Value) (v data.Value, err error) {
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
	leftString, err := data.AsString(leftVal)
	if err != nil {
		return nil, fmt.Errorf("left operand of || must be string: %v", leftVal)
	}
	rightString, err := data.AsString(rightVal)
	if err != nil {
		return nil, fmt.Errorf("right operand of || must be string: %v", rightVal)
	}
	return data.String(leftString + rightString), nil
}

/// Function Evaluation

type funcApp struct {
	name        string
	fVal        reflect.Value
	params      []Evaluator
	paramValues []reflect.Value
	selector    data.Path
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
	if f.selector != nil {
		switch result.Type() {
		case data.TypeMap:
			retmap, _ := data.AsMap(result)
			selected, err := retmap.Get(f.selector)
			if err != nil {
				return nil, err
			}
			result = selected
		case data.TypeArray:
			retarr, _ := data.AsArray(result)
			selected, err := retarr.Get(f.selector)
			if err != nil {
				return nil, err
			}
			result = selected
		default:
			return nil, fmt.Errorf("type '%v' is not supported with selector", result.Type())
		}
	}
	return result, nil
}

// FuncApp represents evaluation of a function on a number
// of parameters that are expressions over an input Value.
func FuncApp(name string, f udf.UDF, ctx *core.Context, params []Evaluator) Evaluator {
	fVal := reflect.ValueOf(f.Call)
	paramValues := make([]reflect.Value, len(params)+1)
	paramValues[0] = reflect.ValueOf(ctx)
	return &funcApp{
		name:        name,
		fVal:        fVal,
		params:      params,
		paramValues: paramValues,
	}
}

// FuncAppSelector represents function and selector.
// After evaluate the function, return the selected value.
func FuncAppSelector(funcEval *funcApp, selector string) (Evaluator, error) {
	if strings.HasPrefix(selector, ".") {
		selector = selector[1:]
	}
	path, err := data.CompilePath(selector)
	if err != nil {
		return nil, err
	}
	funcEval.selector = path
	return funcEval, nil
}

/// Aggregate Function with Sorted Input

type sortEvaluator struct {
	eval      Evaluator
	ascending bool
}

type sortedInputAggFuncApp struct {
	f         Evaluator
	inOutKeys map[string]string
	ordering  []sortEvaluator
}

func (s *sortedInputAggFuncApp) Eval(input data.Value) (v data.Value, err error) {
	// catch panic (e.g., in called function)
	defer func() {
		if r := recover(); r != nil {
			v = nil
			err = fmt.Errorf("evaluating %v paniced: %s", s.f, r)
		}
	}()
	inputMap, err := data.AsMap(input)
	if err != nil {
		return nil, err
	}
	// extract the arrays that contain the data that the sort is based on
	if len(s.ordering) == 0 {
		return nil, fmt.Errorf("order definition must not be empty")
	}
	sortData := make([]sortArray, len(s.ordering))
	for i, sortEval := range s.ordering {
		val, err := sortEval.eval.Eval(input)
		if err != nil {
			return nil, fmt.Errorf("could not get data for sorting: %s", err.Error())
		}
		arr, err := data.AsArray(val)
		if err != nil {
			return nil, err
		}
		sortData[i] = sortArray{arr, sortEval.ascending}
	}

	// sort an array of indexes, then write the actual data to a new array
	// using the sorted index array
	indexes := make([]int, len(sortData[0].values))
	for i := range indexes {
		indexes[i] = i
	}
	// sort the index array
	is := &indexSlice{indexes, sortData}
	sort.Sort(is)

	// now use the sorted index array to write a sorted copy of the data
	for unsortedKey, sortedKey := range s.inOutKeys {
		unsortedData, ok := inputMap[unsortedKey]
		if !ok {
			return nil, fmt.Errorf("there was no unsorted data with key '%s'", unsortedKey)
		}
		unsortedArr, err := data.AsArray(unsortedData)
		if err != nil {
			return nil, err
		}
		if len(unsortedArr) != len(indexes) {
			return nil, fmt.Errorf("aggregate data with key '%s' had bad length (%d, not %d)",
				unsortedKey, len(unsortedArr), len(indexes))
		}
		sortedArr := make(data.Array, len(unsortedArr))
		for i := range unsortedArr {
			sortedArr[i] = unsortedArr[indexes[i]]
		}
		inputMap[sortedKey] = sortedArr
	}

	return s.f.Eval(input)
}

func newSortedInputAggFuncApp(obj funcAppAST, id string, ordering []sortExpression, reg udf.FunctionRegistry) (Evaluator, error) {
	// We may have a function call as complex as
	//  f(a, b, c ORDER BY d ASC, e DESC)
	// where a and c are aggregate parameters but b is not.
	//
	// The Eval() call will get input data like
	//   data.Map{"stream": data.Map{"a": data.Int(1)},
	//            "g_ahash": data.Array{data.Int(1), data.Int(2)},
	//            "g_chash": data.Array{data.Int(3), data.Int(4)},
	//            "g_dhash": data.Array{data.Int(5), data.Int(6)},
	//            "g_ehash": data.Array{data.Int(7), data.Int(8)}}
	//
	// With the above function clause, we have that obj.Expressions[0] and
	// obj.Expressions[2] are aggInputRef structs and b is something else.
	// Also, ordering is {{aggInputRef, true}, {aggInputRef, false}}.
	//
	// How can we call f with a sorted version of the arrays named
	// g_ahash and g_chash?
	// 1) If we sort them in-place, there may be trouble if there is
	//    another function g(a ORDER BY d DESC) which needs a different
	//    ordering.
	// 2) Another approach is to add a sorted version of the array,
	//    suffixed with some unique string, so that the Map looks like
	//     data.Map{"stream": data.Map{"a": data.Int(1)},
	//              "g_ahash":     data.Array{data.Int(1), data.Int(2)},
	//              "g_ahash_abc": data.Array{data.Int(2), data.Int(1)},
	//              "g_chash":     data.Array{data.Int(3), data.Int(4)},
	//              "g_chash_abc": data.Array{data.Int(4), data.Int(3)},
	//              "g_dhash":     data.Array{data.Int(5), data.Int(6)},
	//              "g_ehash":     data.Array{data.Int(7), data.Int(8)}}
	//    and then change the evaluators computed from obj.Expressions[0]
	//    and obj.Expressions[2] to use these sorted versions.
	// The second approach may not provide optimal performance, but it
	// avoids weird side effects when using multiple aggregates.

	if len(ordering) == 0 {
		return nil, fmt.Errorf("order definition must not be empty")
	}
	sortEvals := make([]sortEvaluator, len(ordering))
	for i, sortExpr := range ordering {
		e, err := ExpressionToEvaluator(sortExpr.Value, reg)
		if err != nil {
			return nil, err
		}
		sortEvals[i] = sortEvaluator{e, sortExpr.Ascending}
	}

	// lookup function in function registry
	// (the registry will decide if the requested function
	// is callable with the given number of arguments).
	fName := string(obj.Function)
	f, err := reg.Lookup(fName, len(obj.Expressions))
	if err != nil {
		return nil, err
	}
	// compute child Evaluators
	inOutKeys := map[string]string{}
	evals := make([]Evaluator, len(obj.Expressions))
	for i, ast := range obj.Expressions {
		if inputRef, ok := ast.(aggInputRef); ok {
			newRef := inputRef.Ref + "_" + id
			ast = aggInputRef{newRef}
			inOutKeys[inputRef.Ref] = newRef
		}
		eval, err := ExpressionToEvaluator(ast, reg)
		if err != nil {
			return nil, err
		}
		evals[i] = eval
	}
	backendFun := FuncApp(fName, f, reg.Context(), evals)

	return &sortedInputAggFuncApp{backendFun, inOutKeys, sortEvals}, nil
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

func newArrayBuilder(elems []Evaluator) Evaluator {
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

func newMapBuilder(names []string, elems []Evaluator) (Evaluator, error) {
	if len(names) != len(elems) {
		return nil, fmt.Errorf("number of keys and values does not match")
	}
	return &mapBuilder{names, elems}, nil
}

// CASE statement
type caseBuilder struct {
	reference   Evaluator
	whens       []Evaluator
	thens       []Evaluator
	defaultEval Evaluator
}

func (c *caseBuilder) Eval(input data.Value) (v data.Value, err error) {
	predicate, err := c.reference.Eval(input)
	if err != nil {
		return nil, err
	}
	resultEval := c.defaultEval
	for i, when := range c.whens {
		whenValue, err := when.Eval(input)
		if err != nil {
			return nil, err
		}
		if data.Equal(predicate, whenValue) {
			resultEval = c.thens[i]
			break
		}
	}
	return resultEval.Eval(input)
}

func newCaseBuilder(ref Evaluator, whens []Evaluator, thens []Evaluator, def Evaluator) (Evaluator, error) {
	if len(whens) != len(thens) {
		return nil, fmt.Errorf("number of when and then expressions does not match")
	}
	return &caseBuilder{ref, whens, thens, def}, nil
}

// wildcard only works on Maps, assumes that the elements which do not contain
// ":meta:" are also Maps and pulls them up one level, so
//   {"a": {"x": ...}, "a:meta:ts": ..., "b": {"y": ..., "z": ...}}
// becomes
//   {"x": ..., "y": ..., "z": ...}.
// If there are keys appearing in multiple top-level Maps, then only one
// of them will appear in the output, but it is undefined which.
// If the `Relation` member is non-empty, only the Map with that key will
// be pulled up.
type wildcard struct {
	Relation string
}

func (w *wildcard) Eval(input data.Value) (data.Value, error) {
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

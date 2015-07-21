package builtin

import (
	"fmt"
	"math"
	"math/rand"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// singleParamFunc is a template for functions that
// have exactly one parameter
type singleParamFunc struct {
}

func (f *singleParamFunc) Accept(arity int) bool {
	return arity == 1
}

func (f *singleParamFunc) IsAggregationParameter(k int) bool {
	return false
}

// twoParamFunc is a template for functions that
// have exactly two parameters
type twoParamFunc struct {
}

func (f *twoParamFunc) Accept(arity int) bool {
	return arity == 2
}

func (f *twoParamFunc) IsAggregationParameter(k int) bool {
	return false
}

// typePreservingSingleParamNumericFunc is a template for
// numeric functions that have the same return type as
// input type. If intFun is nil, then the result is computed
// by converting input and output of floatFun.
type typePreservingSingleParamNumericFunc struct {
	singleParamFunc
	intFun   func(int64) int64
	floatFun func(float64) float64
}

func (f *typePreservingSingleParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 1 {
		return nil, fmt.Errorf("function takes exactly one argument")
	}
	arg := args[0]
	if arg.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg.Type() == data.TypeInt {
		i, _ := data.AsInt(arg)
		if f.intFun == nil {
			return data.Int(f.floatFun(float64(i))), nil
		}
		return data.Int(f.intFun(i)), nil
	} else if arg.Type() == data.TypeFloat {
		d, _ := data.AsFloat(arg)
		return data.Float(f.floatFun(d)), nil
	}
	return nil, fmt.Errorf("cannot interpret %s as number", arg)
}

// floatValuedSingleParamNumericFunc is a template for
// numeric functions that return a floating point value
// even if the input is integral. If intFun is nil, then
// the result is computed by converting input of floatFun.
type floatValuedSingleParamNumericFunc struct {
	singleParamFunc
	intFun   func(int64) float64
	floatFun func(float64) float64
}

func (f *floatValuedSingleParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 1 {
		return nil, fmt.Errorf("function takes exactly one argument")
	}
	arg := args[0]
	if arg.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg.Type() == data.TypeInt {
		i, _ := data.AsInt(arg)
		if f.intFun == nil {
			return data.Float(f.floatFun(float64(i))), nil
		}
		return data.Float(f.intFun(i)), nil
	} else if arg.Type() == data.TypeFloat {
		d, _ := data.AsFloat(arg)
		return data.Float(f.floatFun(d)), nil
	}
	return nil, fmt.Errorf("cannot interpret %s as number", arg)
}

// intValuedSingleParamNumericFunc is a template for
// numeric functions that return an integer value
// even if the input is a floating point number.
// If intFun is nil, then the result is computed
// by converting input of floatFun.
type intValuedSingleParamNumericFunc struct {
	singleParamFunc
	intFun   func(int64) int64
	floatFun func(float64) int64
}

func (f *intValuedSingleParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 1 {
		return nil, fmt.Errorf("function takes exactly one argument")
	}
	arg := args[0]
	if arg.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg.Type() == data.TypeInt {
		i, _ := data.AsInt(arg)
		if f.intFun == nil {
			return data.Int(f.floatFun(float64(i))), nil
		}
		return data.Int(f.intFun(i)), nil
	} else if arg.Type() == data.TypeFloat {
		d, _ := data.AsFloat(arg)
		return data.Int(f.floatFun(d)), nil
	}
	return nil, fmt.Errorf("cannot interpret %s as number", arg)
}

// typePreservingTwoParamNumericFunc is a template for
// numeric functions that have the same return type as
// input types.
type typePreservingTwoParamNumericFunc struct {
	twoParamFunc
	intFun   func(int64, int64) int64
	floatFun func(float64, float64) float64
}

func (f *typePreservingTwoParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 2 {
		return nil, fmt.Errorf("function takes exactly two arguments")
	}
	arg1 := args[0]
	arg2 := args[1]
	if arg1.Type() == data.TypeNull || arg2.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg1.Type() == data.TypeInt && arg2.Type() == data.TypeInt {
		i1, _ := data.AsInt(arg1)
		i2, _ := data.AsInt(arg2)
		return data.Int(f.intFun(i1, i2)), nil
	} else if arg1.Type() == data.TypeFloat && arg2.Type() == data.TypeFloat {
		d1, _ := data.AsFloat(arg1)
		d2, _ := data.AsFloat(arg2)
		return data.Float(f.floatFun(d1, d2)), nil
	} else if arg1.Type() != arg2.Type() {
		return nil, fmt.Errorf("types %T and %T do not match", arg1, arg2)
	}
	return nil, fmt.Errorf("cannot interpret %s and/or %s as number", arg1, arg2)
}

// intValuedTwoParamNumericFunc is a template for
// numeric functions that have an integer as output.
type intValuedTwoParamNumericFunc struct {
	twoParamFunc
	intFun   func(int64, int64) int64
	floatFun func(float64, float64) int64
}

func (f *intValuedTwoParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 2 {
		return nil, fmt.Errorf("function takes exactly two arguments")
	}
	arg1 := args[0]
	arg2 := args[1]
	if arg1.Type() == data.TypeNull || arg2.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg1.Type() == data.TypeInt && arg2.Type() == data.TypeInt {
		i1, _ := data.AsInt(arg1)
		i2, _ := data.AsInt(arg2)
		return data.Int(f.intFun(i1, i2)), nil
	} else if arg1.Type() == data.TypeFloat && arg2.Type() == data.TypeFloat {
		d1, _ := data.AsFloat(arg1)
		d2, _ := data.AsFloat(arg2)
		return data.Int(f.floatFun(d1, d2)), nil
	} else if arg1.Type() != arg2.Type() {
		return nil, fmt.Errorf("types %T and %T do not match", arg1, arg2)
	}
	return nil, fmt.Errorf("cannot interpret %s and/or %s as integer", arg1, arg2)
}

// floatValuedTwoParamNumericFunc is a template for
// numeric functions that have a float as output.
type floatValuedTwoParamNumericFunc struct {
	twoParamFunc
	floatFun func(float64, float64) float64
}

func (f *floatValuedTwoParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 2 {
		return nil, fmt.Errorf("function takes exactly two arguments")
	}
	arg1 := args[0]
	arg2 := args[1]
	if arg1.Type() == data.TypeNull || arg2.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg1.Type() == data.TypeInt && arg2.Type() == data.TypeInt {
		i1, _ := data.AsInt(arg1)
		i2, _ := data.AsInt(arg2)
		return data.Float(f.floatFun(float64(i1), float64(i2))), nil
	} else if arg1.Type() == data.TypeFloat && arg2.Type() == data.TypeFloat {
		d1, _ := data.AsFloat(arg1)
		d2, _ := data.AsFloat(arg2)
		return data.Float(f.floatFun(d1, d2)), nil
	} else if arg1.Type() != arg2.Type() {
		return nil, fmt.Errorf("types %T and %T do not match", arg1, arg2)
	}
	return nil, fmt.Errorf("cannot interpret %s and/or %s as integer", arg1, arg2)
}

// unaryBinaryDispatcher supports the overloading of BQL functions
// by forwarding a request to the correct sub-UDF.
type unaryBinaryDispatcher struct {
	unary  udf.UDF
	binary udf.UDF
}

func (f *unaryBinaryDispatcher) Accept(arity int) bool {
	return arity == 1 || arity == 2
}

func (f *unaryBinaryDispatcher) IsAggregationParameter(k int) bool {
	return false
}

func (f *unaryBinaryDispatcher) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	if len(args) == 1 {
		return f.unary.Call(ctx, args...)
	} else if len(args) == 2 {
		return f.binary.Call(ctx, args...)
	}
	return nil, fmt.Errorf("function takes either one or two arguments")
}

var intId = func(a int64) int64 { return a }

// absFunc computes the absolute value of a number.
// See also: math.Abs.
//
// It can be used in BQL as `abs`.
//
//  Input: Int or Float
//  Return Type: same as input
var absFunc udf.UDF = &typePreservingSingleParamNumericFunc{
	floatFun: math.Abs,
}

// cbrtFunc computes the cube root of a number.
// See also: math.Cbrt.
//
// It can be used in BQL as `cbrt`.
//
//  Input: Int or Float
//  Return Type: Float
var cbrtFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Cbrt,
}

// ceilFunc computes the smallest integer not less than its argument.
// See also: math.Ceil.
//
// It can be used in BQL as `ceil`.
//
//  Input: Int or Float
//  Return Type: same as input
//
// Note: This function returns a Float for Float input in order
// to avoid truncation when converting to Int.
var ceilFunc udf.UDF = &typePreservingSingleParamNumericFunc{
	intFun:   intId,
	floatFun: math.Ceil,
}

// degreesFunc converts radians to degrees.
//
// It can be used in BQL as `degrees`.
//
//  Input: Int or Float
//  Return Type: Float
var degreesFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: func(f float64) float64 {
		return f / math.Pi * 180
	},
}

// divFunc computes the integer quotient of two numbers.
//
// It can be used in BQL as `div`.
//
//  Input: 2 * Int or Float
//  Return Type: same as input
//
// Note: This function returns a Float for Float input in order
// to avoid truncation when converting to Int.
var divFunc udf.UDF = &typePreservingTwoParamNumericFunc{
	intFun: func(a, b int64) int64 { return a / b },
	floatFun: func(a, b float64) float64 {
		if b == 0 {
			return math.NaN()
		}
		return float64(int64(a / b))
	},
}

// expFunc computes the exponential of a number.
// See also: math.Exp.
//
// It can be used in BQL as `exp`.
//
//  Input: Int or Float
//  Return Type: Float
var expFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Exp,
}

// floorFunc computes the largest integer not greater than its argument.
// See also: math.Floor.
//
// It can be used in BQL as `floor`.
//
//  Input: Int or Float
//  Return Type: same as input
//
// Note: This function returns a Float for Float input in order
// to avoid truncation when converting to Int.
var floorFunc udf.UDF = &typePreservingSingleParamNumericFunc{
	intFun:   intId,
	floatFun: math.Floor,
}

// lnFunc computes the natural logarithm of a number.
// See also: math.Log
//
// It can be used in BQL as `ln`.
//
//  Input: Int or Float
//  Return Type: Float
var lnFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Log,
}

// logFunc computes the base 10 logarithm of a number.
// See also: math.Log10
//
// It can be used in BQL as `log`.
//
//  Input: Int or Float
//  Return Type: Float
var logFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Log10,
}

// logBaseFunc(b, x) computes the logarithm of a number x to base b.
// See also: math.Log10
//
// It can be used in BQL as `log`.
//
//  Input: 2 * Int or Float
//  Return Type: Float
var logBaseFunc udf.UDF = &floatValuedTwoParamNumericFunc{
	floatFun: func(b, x float64) float64 {
		return math.Log(x) / math.Log(b)
	},
}

// modFunc computes the remainder of integer division.
// See also: math.Mod
//
// It can be used in BQL as `mod`.
//
//  Input: 2 * Int or Float
//  Return Type: same as input
var modFunc udf.UDF = &typePreservingTwoParamNumericFunc{
	intFun:   func(a, b int64) int64 { return a % b },
	floatFun: math.Mod,
}

// piFunc returns the pi constant (more or less 3.14).
//
// It can be used in BQL as `pi`.
//
//  Input: None
//  Return Type: Float
var piFunc, _ = udf.ConvertGeneric(func() float64 { return math.Pi })

// powFunc computes its first parameter raised to the power of its second.
// See also: math.Pow
//
// It can be used in BQL as `power`.
//
//  Input: 2 * Int or Float
//  Return Type: Float
//
// Note: This function always returns a Float in order to deal with
// negative integer exponents properly.
var powFunc udf.UDF = &floatValuedTwoParamNumericFunc{
	floatFun: math.Pow,
}

// radiansFunc converts degrees to radians.
//
// It can be used in BQL as `radians`.
//
//  Input: Int or Float
//  Return Type: Float
var radiansFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: func(f float64) float64 {
		return f * math.Pi / 180
	},
}

// roundFunc computes the nearest integer of a number.
// See also: math.Round.
//
// It can be used in BQL as `round`.
//
//  Input: Int or Float
//  Return Type: same as input
//
// Note: This function returns a Float for Float input in order
// to avoid truncation when converting to Int.
var roundFunc udf.UDF = &typePreservingSingleParamNumericFunc{
	intFun: intId,
	floatFun: func(a float64) float64 {
		if a < 0 {
			return math.Ceil(a - 0.5)
		}
		return math.Floor(a + 0.5)
	},
}

// signFunc computes the sign (-1, 0, +1) of a number.
//
// It can be used in BQL as `sign`.
//
//  Input: Int or Float
//  Return Type: Int
var signFunc udf.UDF = &intValuedSingleParamNumericFunc{
	intFun: func(i int64) int64 {
		if i > 0 {
			return 1
		} else if i < 0 {
			return -1
		}
		return 0
	},
	floatFun: func(f float64) int64 {
		if f > 0 {
			return 1
		} else if f < 0 {
			return -1
		}
		return 0
	},
}

// sqrtFunc computes the square root of a number.
// See also: math.Sqrt
//
// It can be used in BQL as `sqrt`.
//
//  Input: Int or Float
//  Return Type: Float
var sqrtFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Sqrt,
}

// truncFunc computes the truncated integer (towards zero) of a number.
//
// It can be used in BQL as `trunc`.
//
//  Input: Int or Float
//  Return Type: same as input
//
// Note: This function returns a Float for Float input in order
// to avoid truncation when converting to Int.
var truncFunc udf.UDF = &typePreservingSingleParamNumericFunc{
	intFun: intId,
	floatFun: func(a float64) float64 {
		if a < 0 {
			return math.Ceil(a)
		}
		return math.Floor(a)
	},
}

type widthBucketFuncTmpl struct {
}

func (f *widthBucketFuncTmpl) Accept(arity int) bool {
	return arity == 4
}

func (f *widthBucketFuncTmpl) IsAggregationParameter(k int) bool {
	return false
}

func (f *widthBucketFuncTmpl) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 4 {
		return nil, fmt.Errorf("function takes exactly four argument")
	}
	// return if any parameter is null
	for _, arg := range args {
		if arg.Type() == data.TypeNull {
			return data.Null{}, nil
		}
	}
	// convert input for our algorithm
	var x, left, right float64
	var count int64
	param := 0
	if args[param].Type() == data.TypeInt {
		i, _ := data.AsInt(args[param])
		x = float64(i)
	} else if args[param].Type() == data.TypeFloat {
		x, _ = data.AsFloat(args[param])
	} else {
		return nil, fmt.Errorf("%d-th parameter must be Int or Float", param)
	}
	param = 1
	if args[param].Type() == data.TypeInt {
		i, _ := data.AsInt(args[param])
		left = float64(i)
	} else if args[param].Type() == data.TypeFloat {
		left, _ = data.AsFloat(args[param])
	} else {
		return nil, fmt.Errorf("%d-th parameter must be Int or Float", param)
	}
	param = 2
	if args[param].Type() == data.TypeInt {
		i, _ := data.AsInt(args[param])
		right = float64(i)
	} else if args[param].Type() == data.TypeFloat {
		right, _ = data.AsFloat(args[param])
	} else {
		return nil, fmt.Errorf("%d-th parameter must be Int or Float", param)
	}
	if right <= left {
		return nil, fmt.Errorf("right interval border (%v) must be larger "+
			"than left border (%v)", right, left)
	}
	param = 3
	if args[param].Type() == data.TypeInt {
		count, _ = data.AsInt(args[param])
	} else {
		return nil, fmt.Errorf("%d-th parameter must be Int", param)
	}
	if count <= 0 {
		return nil, fmt.Errorf("bucket count must be strictly positive")
	}
	// now compute the actual bucket:
	//     left        right
	//       v           v
	//  0    | 1 | 2 | 3 |    4
	if x < left {
		return data.Int(0), nil
	} else if x >= right {
		return data.Int(count + 1), nil
	}
	width := (right - left) / float64(count)
	return data.Int(int64((x-left)/width) + 1), nil
}

// widthBucketFunc(x, left, right, count) computes the bucket to which x
// would be assigned in an equidepth histogram with count buckets in
// the range [left,right[. Points on a bucket border belong to the right
// bucket. Points outside of the [left,right[ range have bucket number
// 0 and count+1, respectively.
//
// It can be used in BQL as `width_bucket`.
//
//  Input: 3 * Int or Float, Int
//  Return Type: Int
var widthBucketFunc udf.UDF = &widthBucketFuncTmpl{}

// randomFunc returns a random number in the range [0,1[.
// See also: math/rand.Float64()
//
// It can be used in BQL as `random`.
//
//  Input: None
//  Return Type: Float
var randomFunc, _ = udf.ConvertGeneric(func() float64 { return rand.Float64() })

type setseedFuncTmpl struct {
	singleParamFunc
}

func (f *setseedFuncTmpl) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 1 {
		return nil, fmt.Errorf("function takes exactly one argument")
	}
	arg := args[0]
	if arg.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg.Type() == data.TypeFloat {
		d, _ := data.AsFloat(arg)
		if d < -1.0 || d > 1.0 {
			return nil, fmt.Errorf("seed out of range [-1,1]")
		}
		s := int64(d * float64(math.MaxInt64))
		rand.Seed(s)
		return data.Null{}, nil
	}
	return nil, fmt.Errorf("cannot interpret %s as float", arg)
}

// setseed initializes the seed for subsequent randomFunc calls.
// The argument must be a float in the range [-1,1].
// See also: math/rand.Seed()
//
// It can be used in BQL as `setseed`.
//
//  Input: Float
//  Return Type: Null
var setseedFunc udf.UDF = &setseedFuncTmpl{}

// acosFunc computes the inverse cosine of a number.
// See also: math.Acos
//
// It can be used in BQL as `acos`.
//
//  Input: Int or Float
//  Return Type: Float
var acosFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Acos,
}

// asinFunc computes the inverse sine of a number.
// See also: math.Asin
//
// It can be used in BQL as `asin`.
//
//  Input: Int or Float
//  Return Type: Float
var asinFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Asin,
}

// atanFunc computes the inverse tangent of a number.
// See also: math.Atan
//
// It can be used in BQL as `atan`.
//
//  Input: Int or Float
//  Return Type: Float
var atanFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Atan,
}

// cosFunc computes the cosine of a number.
// See also: math.Cos
//
// It can be used in BQL as `cos`.
//
//  Input: Int or Float
//  Return Type: Float
var cosFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Cos,
}

// cotFunc computes the cotangent of a number.
//
// It can be used in BQL as `cot`.
//
//  Input: Int or Float
//  Return Type: Float
var cotFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: func(x float64) float64 { return 1. / math.Tan(x) },
}

// sinFunc computes the sine of a number.
// See also: math.Sin
//
// It can be used in BQL as `sin`.
//
//  Input: Int or Float
//  Return Type: Float
var sinFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Sin,
}

// tanFunc computes the tangent of a number.
// See also: math.Tan
//
// It can be used in BQL as `tan`.
//
//  Input: Int or Float
//  Return Type: Float
var tanFunc udf.UDF = &floatValuedSingleParamNumericFunc{
	floatFun: math.Tan,
}

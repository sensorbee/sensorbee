package builtin

import (
	"fmt"
	"math"
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
	} else if arg.Type() == data.TypeInt && f.intFun == nil {
		i, _ := data.AsInt(arg)
		return data.Int(f.floatFun(float64(i))), nil
	} else if arg.Type() == data.TypeInt {
		i, _ := data.AsInt(arg)
		return data.Int(f.intFun(i)), nil
	} else if arg.Type() == data.TypeFloat {
		d, _ := data.AsFloat(arg)
		return data.Float(f.floatFun(d)), nil
	}
	return nil, fmt.Errorf("cannot interpret %s as number", arg)
}

// floatValuedSingleParamNumericFunc is a template for
// numeric functions that return a floating point value
// even if the input is integral
type floatValuedSingleParamNumericFunc struct {
	singleParamFunc
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
		return data.Float(f.floatFun(float64(i))), nil
	} else if arg.Type() == data.TypeFloat {
		d, _ := data.AsFloat(arg)
		return data.Float(f.floatFun(d)), nil
	}
	return nil, fmt.Errorf("cannot interpret %s as number", arg)
}

// intValuedTwoParamNumericFunc is a template for
// numeric functions that have integers as output.
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
var ceilFunc udf.UDF = &typePreservingSingleParamNumericFunc{
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
//  Input: Int, Int
//  Return Type: Int
var divFunc udf.UDF = &intValuedTwoParamNumericFunc{
	intFun: func(a, b int64) int64 { return a / b },
	floatFun: func(a, b float64) int64 {
		if b == 0 {
			panic("division by zero")
		}
		return int64(a / b)
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
var floorFunc udf.UDF = &typePreservingSingleParamNumericFunc{
	floatFun: math.Floor,
}

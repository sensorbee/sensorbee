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

// typePreservingSingleParamNumericFunc is a template for
// numeric functions that have the same return type as
// input type. If intFun is nil, then the result is computed
// by converting input and output of floatFun.
type typePreservingSingleParamNumericFunc struct {
	singleParamFunc
	intFun   func(int64) int64
	floatFun func(float64) float64
}

func (f *typePreservingSingleParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
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

func (f *floatValuedSingleParamNumericFunc) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
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

// AbsFunc returns a UDF to compute the absolute value of a number.
// See also: math.Abs.
//
// That UDF can be used in BQL as `abs`.
//
//  Input: Int or Float
//  Return Type: same as input
func AbsFunc() udf.UDF {
	return &typePreservingSingleParamNumericFunc{
		floatFun: math.Abs,
	}
}

// CbrtFunc returns a UDF to compute the cube root of a number.
// See also: math.Cbrt.
//
// That UDF can be used in BQL as `cbrt`.
//
//  Input: Int or Float
//  Return Type: Float
func CbrtFunc() udf.UDF {
	return &floatValuedSingleParamNumericFunc{
		floatFun: math.Cbrt,
	}
}

// CeilFunc returns a UDF to compute the smallest integer not less
// than its argument. See also: math.Ceil.
//
// That UDF can be used in BQL as `ceil`.
//
//  Input: Int or Float
//  Return Type: same as input
func CeilFunc() udf.UDF {
	return &typePreservingSingleParamNumericFunc{
		floatFun: math.Ceil,
	}
}

// DegreesFunc returns a UDF to convert radians to degrees.
//
// That UDF can be used in BQL as `degrees`.
//
//  Input: Int or Float
//  Return Type: Float
func DegreesFunc() udf.UDF {
	return &floatValuedSingleParamNumericFunc{
		floatFun: func(f float64) float64 {
			return f / math.Pi * 180
		},
	}
}

// TODO DivFunc

// ExpFunc returns a UDF to compute the exponential of a number.
// See also: math.Exp.
//
// That UDF can be used in BQL as `exp`.
//
//  Input: Int or Float
//  Return Type: Float
func ExpFunc() udf.UDF {
	return &floatValuedSingleParamNumericFunc{
		floatFun: math.Exp,
	}
}

// FloorFunc returns a UDF to compute the largest integer not greater
// than its argument. See also: math.Floor.
//
// That UDF can be used in BQL as `ceil`.
//
//  Input: Int or Float
//  Return Type: same as input
func FloorFunc() udf.UDF {
	return &typePreservingSingleParamNumericFunc{
		floatFun: math.Floor,
	}
}

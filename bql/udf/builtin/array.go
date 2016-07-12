package builtin

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

type arrayLengthFuncTmpl struct {
	singleParamFunc
}

func (f *arrayLengthFuncTmpl) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
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
	} else if arg.Type() == data.TypeArray {
		a, _ := data.AsArray(arg)
		return data.Int(len(a)), nil
	}
	return nil, fmt.Errorf("%v is not an array", arg)
}

// arrayLengthFunc returns the length of the given array.
// NULL elements are counted as well.
//
// It can be used in BQL as `array_length`.
//
//  Input: Array
//  Return Type: Int
var arrayLengthFunc udf.UDF = &arrayLengthFuncTmpl{}

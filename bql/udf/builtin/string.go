package builtin

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// singleParamStringFunc is a template for functions that
// have a single string parameter as input.
type singleParamStringFunc struct {
	singleParamFunc
	strFun func(string) data.Value
}

func (f *singleParamStringFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
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
	} else if arg.Type() == data.TypeString {
		s, _ := data.AsString(arg)
		return f.strFun(s), nil
	}
	return nil, fmt.Errorf("cannot interpret %s as a string", arg)
}

// bitLengthFunc computes the number of bits in a string.
//
// It can be used in BQL as `bit_length`.
//
//  Input: String
//  Return Type: Int
var bitLengthFunc udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		return data.Int(len(s) * 8)
	},
}

package builtin

import (
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type coalesceFuncTmpl struct {
	variadicFunc
}

func (f *coalesceFuncTmpl) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	for _, item := range args {
		if item.Type() != data.TypeNull {
			return item, nil
		}
	}
	return data.Null{}, nil
}

// coalesceFunc returns the first non-null argument.
//
// It can be used in BQL as `coalesce`.
//
//  Input: n * Any
//  Return Type: same as the first non-null argument
var coalesceFunc udf.UDF = &coalesceFuncTmpl{}

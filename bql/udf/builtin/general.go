package builtin

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

// coalesceFunc returns the first non-null argument.
//
// It can be used in BQL as `coalesce`.
//
//  Input: n * Any
//  Return Type: same as the first non-null argument
var coalesceFunc udf.UDF = &variadicFunc{
	minParams: 1,
	varFun: func(args ...data.Value) (data.Value, error) {
		for _, item := range args {
			if item.Type() != data.TypeNull {
				return item, nil
			}
		}
		return data.Null{}, nil
	},
}

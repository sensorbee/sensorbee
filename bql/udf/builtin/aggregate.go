package builtin

import (
	"bytes"
	"fmt"
	"math"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// singleParamAggFunc is a template for aggregate functions that
// have exactly one parameter
type singleParamAggFunc struct {
	aggFun func([]data.Value) (data.Value, error)
}

func (f *singleParamAggFunc) Accept(arity int) bool {
	return arity == 1
}

func (f *singleParamAggFunc) IsAggregationParameter(k int) bool {
	return k == 0
}

func (f *singleParamAggFunc) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("function takes exactly one argument")
	}
	arr, err := data.AsArray(args[0])
	if err != nil {
		return nil, fmt.Errorf("function needs array input, not %T", args[0])
	}
	return f.aggFun(arr)
}

// countFunc is an aggregate function that counts the number
// of non-null values passed in.
//
// It can be used in BQL as `count`.
//
//  Input: anything (aggregated)
//  Return Type: Int
var countFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		// count() is O(n) in the spirit of PostgreSQL
		c := int64(0)
		for _, item := range arr {
			if item.Type() != data.TypeNull {
				c++
			}
		}
		return data.Int(c), nil
	},
}

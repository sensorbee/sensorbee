package udf

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type countAggregate struct {
}

func (f *countAggregate) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("count takes exactly one argument")
	}
	arr, err := data.AsArray(args[0])
	if err != nil {
		return nil, fmt.Errorf("count needs an array input, not %v", args[0])
	}
	// count() is O(n) in the spirit of PostgreSQL
	c := int64(0)
	for _, item := range arr {
		if item.Type() != data.TypeNull {
			c++
		}
	}
	return data.Int(c), nil
}

func (f *countAggregate) Accept(arity int) bool {
	return arity == 1
}

func (f *countAggregate) IsAggregationParameter(k int) bool {
	return k == 1
}

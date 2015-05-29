package udf

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
)

// VarParamFun is an alias for a variadic function on tuple.Value.
type VarParamFun func(...tuple.Value) (tuple.Value, error)

// FunctionRegistry is an interface to lookup functions for use in BQL
// statements by their name.
type FunctionRegistry interface {
	// Lookup will return a function with the given name and arity
	// or an error if there is none. Note that this interface allows
	// multiple functions with the same name but different arity,
	// and it also allows functions with an arbitrary number of
	// parameters. However, a function returned must never be used
	// with a different arity than the one given in the Lookup call.
	Lookup(name string, arity int) (VarParamFun, error)
}

type EmptyFunctionRegistry struct{}

func (fr *EmptyFunctionRegistry) Lookup(name string, arity int) (VarParamFun, error) {
	return nil, fmt.Errorf("no function named '%s'", name)
}

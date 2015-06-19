package udf

import (
	"fmt"
	"pfi/sensorbee/sensorbee/tuple"
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

// FunctionManager is a FunctionRegistry that allows to register
// additional functions.
type FunctionManager interface {
	FunctionRegistry

	// Register allows to add a function that works exactly with
	// the given number of parameters.
	Register(name string, f VarParamFun, arity int) error

	// RegisterVariadic allows to add a variadic function.
	RegisterVariadic(name string, f VarParamFun) error
}

type defaultFunctionRegistry struct {
	funcs map[string]struct {
		fun     VarParamFun
		checker func(int) bool
	}
}

func NewDefaultFunctionRegistry() *defaultFunctionRegistry {
	empty := map[string]struct {
		fun     VarParamFun
		checker func(int) bool
	}{}
	reg := &defaultFunctionRegistry{empty}
	// register some standard functions
	toString := func(v tuple.Value) (tuple.Value, error) {
		return tuple.String(v.String()), nil
	}
	reg.RegisterUnary("str", toString)
	return reg
}

func (fr *defaultFunctionRegistry) Lookup(name string, arity int) (VarParamFun, error) {
	// look for variable-parameter functions
	funWithChecker, exists := fr.funcs[name]
	if exists {
		if funWithChecker.checker(arity) {
			return funWithChecker.fun, nil
		}
		return nil, fmt.Errorf("function '%s' is not %d-ary", name, arity)
	}
	return nil, fmt.Errorf("function '%s' is unknown", name)
}

func (fr *defaultFunctionRegistry) Register(name string, f VarParamFun, arity int) error {
	return fr.registerFlexible(name, f, func(i int) bool {
		return i == arity
	})
}

func (fr *defaultFunctionRegistry) RegisterVariadic(name string, f VarParamFun) error {
	return fr.registerFlexible(name, f, func(i int) bool {
		return true
	})
}

func (fr *defaultFunctionRegistry) registerFlexible(name string, f VarParamFun, arityOk func(int) bool) error {
	_, exists := fr.funcs[name]
	if exists {
		return fmt.Errorf("there is already a function named '%s'", name)
	}
	newData := struct {
		fun     VarParamFun
		checker func(int) bool
	}{f, arityOk}
	fr.funcs[name] = newData
	return nil
}

func (fr *defaultFunctionRegistry) RegisterNullary(name string, f func() (tuple.Value, error)) error {
	genFunc := func(vs ...tuple.Value) (tuple.Value, error) {
		if len(vs) != 0 {
			return nil, fmt.Errorf("function '%s' should be used as nullary", name)
		}
		return f()
	}
	return fr.Register(name, genFunc, 0)
}

func (fr *defaultFunctionRegistry) RegisterUnary(name string, f func(tuple.Value) (tuple.Value, error)) error {
	genFunc := func(vs ...tuple.Value) (tuple.Value, error) {
		if len(vs) != 1 {
			return nil, fmt.Errorf("function '%s' should be used as unary", name)
		}
		return f(vs[0])
	}
	return fr.Register(name, genFunc, 1)
}

func (fr *defaultFunctionRegistry) RegisterBinary(name string, f func(tuple.Value, tuple.Value) (tuple.Value, error)) error {
	genFunc := func(vs ...tuple.Value) (tuple.Value, error) {
		if len(vs) != 2 {
			return nil, fmt.Errorf("function '%s' should be used as binary", name)
		}
		return f(vs[0], vs[1])
	}
	return fr.Register(name, genFunc, 2)
}

func (fr *defaultFunctionRegistry) RegisterTernary(name string, f func(tuple.Value, tuple.Value, tuple.Value) (tuple.Value, error)) error {
	genFunc := func(vs ...tuple.Value) (tuple.Value, error) {
		if len(vs) != 3 {
			return nil, fmt.Errorf("function '%s' should be used as ternary", name)
		}
		return f(vs[0], vs[1], vs[2])
	}
	return fr.Register(name, genFunc, 3)
}

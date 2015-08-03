package udf

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"strings"
	"sync"
)

// UDF is an interface having a user defined function.
type UDF interface {
	// Call calls the UDF.
	Call(*core.Context, ...data.Value) (data.Value, error)

	// Accept checks if the function accepts the given number of arguments
	// excluding core.Context.
	Accept(arity int) bool

	// IsAggregationParameter returns true if the k-th parameter expects
	// aggregated values. A UDF with Accept(n) == true is an aggregate
	// function if and only if this function returns true for one or more
	// values of k in the range 0, ..., n-1.
	IsAggregationParameter(k int) bool
}

type function struct {
	f     func(*core.Context, ...data.Value) (data.Value, error)
	arity int
}

func (f *function) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	return f.f(ctx, args...)
}

func (f *function) Accept(arity int) bool {
	if f.arity < 0 {
		return true
	}
	return arity == f.arity
}

func (f *function) IsAggregationParameter(k int) bool {
	return false
}

// VariadicFunc creates a UDF based on a function receiving the variadic number
// of data.Values.
func VariadicFunc(f func(*core.Context, ...data.Value) (data.Value, error)) UDF {
	return &function{
		f:     f,
		arity: -1,
	}
}

func Func(f func(*core.Context, ...data.Value) (data.Value, error), arity int) UDF {
	return &function{
		f:     f,
		arity: arity,
	}
}

func NullaryFunc(f func(*core.Context) (data.Value, error)) UDF {
	genFunc := func(ctx *core.Context, vs ...data.Value) (data.Value, error) {
		if len(vs) != 0 {
			return nil, fmt.Errorf("the function should be used as nullary")
		}
		return f(ctx)
	}
	return Func(genFunc, 0)
}

func UnaryFunc(f func(*core.Context, data.Value) (data.Value, error)) UDF {
	genFunc := func(ctx *core.Context, vs ...data.Value) (data.Value, error) {
		if len(vs) != 1 {
			return nil, fmt.Errorf("the function should be used as unary")
		}
		return f(ctx, vs[0])
	}
	return Func(genFunc, 1)
}

func BinaryFunc(f func(*core.Context, data.Value, data.Value) (data.Value, error)) UDF {
	genFunc := func(ctx *core.Context, vs ...data.Value) (data.Value, error) {
		if len(vs) != 2 {
			return nil, fmt.Errorf("the function should be used as binary")
		}
		return f(ctx, vs[0], vs[1])
	}
	return Func(genFunc, 2)
}

func TernaryFunc(f func(*core.Context, data.Value, data.Value, data.Value) (data.Value, error)) UDF {
	genFunc := func(ctx *core.Context, vs ...data.Value) (data.Value, error) {
		if len(vs) != 3 {
			return nil, fmt.Errorf("the function should be used as ternary")
		}
		return f(ctx, vs[0], vs[1], vs[2])
	}
	return Func(genFunc, 3)
}

// TODO: Add magic UDF generator func NewUDF(f interface{}) (UDF, error)
//       It accepts any function whose arguments are convertible to data.Value.
//       For example, NewUDF(func(*core.Context, a, b int) (int, error) {return a + b}).
//       Even NewUDF(func(a, b int) int {return a + b}) could be valid
//       (i.e. Context and error can be optional)

// FunctionRegistry is an interface to lookup functions for use in BQL
// statements by their name.
type FunctionRegistry interface {
	// Context returns a core.Context associated with the registry.
	Context() *core.Context

	// Lookup will return a function with the given name and arity
	// or an error if there is none. Note that this interface allows
	// multiple functions with the same name but different arity,
	// and it also allows functions with an arbitrary number of
	// parameters. However, a function returned must never be used
	// with a different arity than the one given in the Lookup call.
	Lookup(name string, arity int) (UDF, error)
}

// FunctionManager is a FunctionRegistry that allows to register
// additional functions.
type FunctionManager interface {
	FunctionRegistry

	// Register allows to add a function.
	Register(name string, f UDF) error
}

type defaultFunctionRegistry struct {
	ctx   *core.Context
	m     sync.RWMutex
	funcs map[string]UDF
}

func NewDefaultFunctionRegistry(ctx *core.Context) FunctionManager {
	reg := &defaultFunctionRegistry{
		ctx:   ctx,
		funcs: map[string]UDF{},
	}
	return reg
}

func (fr *defaultFunctionRegistry) Context() *core.Context {
	return fr.ctx
}

func (fr *defaultFunctionRegistry) Lookup(name string, arity int) (UDF, error) {
	fr.m.RLock()
	defer fr.m.RUnlock()
	if f, exists := fr.funcs[strings.ToLower(name)]; exists {
		if f.Accept(arity) {
			return f, nil
		}
		return nil, fmt.Errorf("function '%s' is not %d-ary", name, arity)
	}
	return nil, fmt.Errorf("function '%s' is unknown", name)
}

func (fr *defaultFunctionRegistry) Register(name string, f UDF) error {
	fr.m.Lock()
	defer fr.m.Unlock()

	lowerName := strings.ToLower(name)
	if _, exists := fr.funcs[lowerName]; exists {
		return fmt.Errorf("there is already a function named '%s'", name)
	}
	fr.funcs[lowerName] = f
	return nil
}

var (
	// globalUDFRegistry has UDFs visible to all topologies. Do NOT use this
	// instance in running topologies because it doesn't have an actual
	// Context. Call CopyGlobalUDSRegistry to retrieve a new valid instance.
	globalUDFRegistry = NewDefaultFunctionRegistry(nil).(*defaultFunctionRegistry)
)

// RegisterGlobalUDF adds a UDF which is visible to all topologies. UDFs
// registered after running topologies might not be seen by those topologies.
// Call it from init functions to avoid such conditions.
func RegisterGlobalUDF(name string, f UDF) error {
	return globalUDFRegistry.Register(name, f)
}

// MustRegisterGlobalUDF is like RegisterGlobalUDF but
// panics if an error occurred.
func MustRegisterGlobalUDF(name string, f UDF) {
	if err := globalUDFRegistry.Register(name, f); err != nil {
		panic(fmt.Errorf("udf.MustRegisterGlobalUDF: cannot register '%v': %v", name, err))
	}
}

// CopyGlobalUDFRegistry creates a new FunctionManager which has all UDFs
// registered to the global function manager.
func CopyGlobalUDFRegistry(ctx *core.Context) FunctionManager {
	globalUDFRegistry.m.RLock()
	defer globalUDFRegistry.m.RUnlock()
	reg := &defaultFunctionRegistry{
		ctx:   ctx,
		funcs: make(map[string]UDF, len(globalUDFRegistry.funcs)),
	}

	for n, f := range globalUDFRegistry.funcs {
		reg.funcs[n] = f
	}
	return reg
}

func init() {
	// register some standard functions
	toString := func(ctx *core.Context, v data.Value) (data.Value, error) {
		return data.String(v.String()), nil
	}
	globalUDFRegistry.Register("str", UnaryFunc(toString))
}

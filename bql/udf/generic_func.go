package udf

import (
	"errors"
	"fmt"
	"math"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"reflect"
	"time"
)

// GenericFunc creates a new UDF from various form of functions. Arguments
// of the function don't have to be tuple types, but some standard types are
// allowed. The UDF returned provide a weak type conversion, that is it uses
// tuple.To{Type} function to convert values. Therefore, a string may be
// passed as an integer or vice versa. If the function wants to provide
// strict type conversion, generate UDF by Func function.
//
// Acceptable types:
//	- bool
//	- standard integers
//	- standard floats
//	- string
//	- time.Time
//	- tuple.Bool, tuple.Int, tuple.Float, tuple.String, tuple.Blob,
//	  tuple.Timestamp, tuple.Array, tuple.Map, tuple.Value
//	- a slice of types above
func GenericFunc(function interface{}) (UDF, error) {
	t := reflect.TypeOf(function)
	if t.Kind() != reflect.Func {
		return nil, errors.New("the argument must be a function")
	}

	g := &genericFunc{
		function:   reflect.ValueOf(function),
		hasContext: genericFuncHasContext(t),
		variadic:   t.IsVariadic(),
		arity:      t.NumIn(),
	}

	if g.hasContext {
		g.arity -= 1
	}

	if hasError, err := checkGenericFuncReturnTypes(t); err != nil {
		return nil, err
	} else {
		g.hasError = hasError
	}

	convs := make([]argumentConverter, 0, g.arity)
	for i := t.NumIn() - g.arity; i < t.NumIn(); i++ {
		arg := t.In(i)
		if i == t.NumIn()-1 && g.variadic {
			arg = arg.Elem()
		}

		c, err := genericFuncArgumentConverter(arg)
		if err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	g.converters = convs
	return g, nil
}

func checkGenericFuncReturnTypes(t reflect.Type) (bool, error) {
	hasError := false

	switch n := t.NumOut(); n {
	case 2:
		if !t.Out(1).Implements(reflect.TypeOf(func(error) {}).In(0)) {
			return false, fmt.Errorf("the second return value must be an error: %v", t.Out(1))
		}
		hasError = true
		fallthrough

	case 1:
		out := t.Out(0)
		if out.Kind() == reflect.Interface {
			// tuple.Value is the only interface which is accepted.
			if !out.Implements(reflect.TypeOf(tuple.NewValue).Out(0)) {
				return false, fmt.Errorf("the return value isn't convertible to tuple.Value")
			}
		}
		if _, err := tuple.NewValue(reflect.Zero(t.Out(0)).Interface()); err != nil {
			return false, fmt.Errorf("the return value isn't convertible to tuple.Value")
		}

	default:
		return false, fmt.Errorf("the number of return values must be 1 or 2: %v", n)
	}
	return hasError, nil
}

func genericFuncHasContext(t reflect.Type) bool {
	if t.NumIn() == 0 {
		return false
	}
	c := t.In(0)
	return reflect.TypeOf(&core.Context{}).AssignableTo(c)
}

type argumentConverter func(tuple.Value) (interface{}, error)

func genericFuncArgumentConverter(t reflect.Type) (argumentConverter, error) {
	// TODO: this function is too long.
	switch t.Kind() {
	case reflect.Bool:
		return func(v tuple.Value) (interface{}, error) {
			return tuple.ToBool(v)
		}, nil

	case reflect.Int:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < -1^int64(^uint(0)>>1) {
				return nil, fmt.Errorf("%v is too small for int", i)
			} else if i > int64(^uint(0)>>1) {
				return nil, fmt.Errorf("%v is too big for int", i)
			}
			return int(i), nil
		}, nil

	case reflect.Int8:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < math.MinInt8 {
				return nil, fmt.Errorf("%v is too small for int8", i)
			} else if i > math.MaxInt8 {
				return nil, fmt.Errorf("%v is too big for int8", i)
			}
			return int8(i), nil
		}, nil

	case reflect.Int16:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < math.MinInt16 {
				return nil, fmt.Errorf("%v is too small for int16", i)
			} else if i > math.MaxInt16 {
				return nil, fmt.Errorf("%v is too big for int16", i)
			}
			return int16(i), nil
		}, nil

	case reflect.Int32:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < math.MinInt32 {
				return nil, fmt.Errorf("%v is too small for int32", i)
			} else if i > math.MaxInt32 {
				return nil, fmt.Errorf("%v is too big for int32", i)
			}
			return int32(i), nil
		}, nil

	case reflect.Int64:
		return func(v tuple.Value) (interface{}, error) {
			return tuple.ToInt(v)
		}, nil

	case reflect.Uint:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < 0 {
				return nil, fmt.Errorf("%v is too small for uint", i)
			} else if i > int64(^uint(0)>>1) {
				return nil, fmt.Errorf("%v is too big for uint", i)
			}
			return uint(i), nil
		}, nil

	case reflect.Uint8:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < 0 {
				return nil, fmt.Errorf("%v is too small for uint8", i)
			} else if i > math.MaxUint8 {
				return nil, fmt.Errorf("%v is too big for uint8", i)
			}
			return uint8(i), nil
		}, nil

	case reflect.Uint16:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < 0 {
				return nil, fmt.Errorf("%v is too small for uint16", i)
			} else if i > math.MaxUint16 {
				return nil, fmt.Errorf("%v is too big for uint16", i)
			}
			return uint16(i), nil
		}, nil

	case reflect.Uint32:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < 0 {
				return nil, fmt.Errorf("%v is too small for uint32", i)
			} else if i > math.MaxUint32 {
				return nil, fmt.Errorf("%v is too big for uint32", i)
			}
			return uint32(i), nil
		}, nil

	case reflect.Uint64:
		return func(v tuple.Value) (interface{}, error) {
			i, err := tuple.ToInt(v)
			if err != nil {
				return nil, err
			}

			if i < 0 {
				return nil, fmt.Errorf("%v is too small for uint64", i)
			}
			return uint64(i), nil
		}, nil

	case reflect.Float32:
		return func(v tuple.Value) (interface{}, error) {
			f, err := tuple.ToFloat(v)
			if err != nil {
				return nil, err
			}
			return float32(f), err
		}, nil

	case reflect.Float64:
		return func(v tuple.Value) (interface{}, error) {
			return tuple.ToFloat(v)
		}, nil

	case reflect.String:
		return func(v tuple.Value) (interface{}, error) {
			return tuple.ToString(v)
		}, nil

	case reflect.Slice:
		elemType := t.Elem()
		if elemType.Kind() == reflect.Uint8 {
			// process this as a blob
			return func(v tuple.Value) (interface{}, error) {
				return tuple.ToBlob(v)
			}, nil
		}

		c, err := genericFuncArgumentConverter(elemType)
		if err != nil {
			return nil, err
		}
		return func(v tuple.Value) (interface{}, error) {
			a, err := tuple.AsArray(v)
			if err != nil {
				return nil, err
			}
			res := reflect.Zero(t)
			for _, elem := range a {
				e, err := c(elem)
				if err != nil {
					return nil, err
				}
				res = reflect.Append(res, reflect.ValueOf(e))
			}
			return res.Interface(), nil
		}, nil

	default:
		switch reflect.Zero(t).Interface().(type) {
		case tuple.Map:
			return func(v tuple.Value) (interface{}, error) {
				return tuple.AsMap(v)
			}, nil

		case time.Time:
			return func(v tuple.Value) (interface{}, error) {
				return tuple.ToTimestamp(v)
			}, nil

		default:
			if t.Implements(reflect.TypeOf(tuple.NewValue).Out(0)) { // tuple.Value
				// Zero(interface) returns nil and type assertion doesn't work for it.
				return func(v tuple.Value) (interface{}, error) {
					return v, nil
				}, nil
			}
			// other tuple types are covered in Kind() switch above
			return nil, fmt.Errorf("unsupported type: %v", t)
		}
	}
}

type genericFunc struct {
	function reflect.Value

	hasContext bool
	hasError   bool
	variadic   bool

	// arity is the number of arguments. If the function is variadic, arity
	// counts the last variadic parameter. For example, if the function is
	// func(int, float, ...string), arity is 3. It doesn't count Context.
	arity int

	converters []argumentConverter
}

func (g *genericFunc) Call(ctx *core.Context, args ...tuple.Value) (tuple.Value, error) {
	out, err := g.call(ctx, args...)
	if err != nil {
		return nil, err
	}

	if g.hasError {
		if !out[1].IsNil() {
			return nil, out[1].Interface().(error)
		}
	}
	return tuple.NewValue(out[0].Interface())
}

func (g *genericFunc) call(ctx *core.Context, args ...tuple.Value) ([]reflect.Value, error) {
	if len(args) < g.arity {
		if g.variadic && len(args) == g.arity-1 {
			// having no variadic parameter is ok.
		} else {
			return nil, fmt.Errorf("insufficient number of argumetns")
		}

	} else if len(args) != g.arity && !g.variadic {
		return nil, fmt.Errorf("too many arguments")
	}

	in := make([]reflect.Value, 0, len(args)+1) // +1 for context
	if g.hasContext {
		in = append(in, reflect.ValueOf(ctx))
	}

	noVariadic := g.arity
	if g.variadic {
		noVariadic--
	}

	for i := 0; i < noVariadic; i++ {
		v, err := g.converters[i](args[i])
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(v))
	}

	if !g.variadic {
		return g.function.Call(in), nil
	}

	for i := noVariadic; i < len(args); i++ {
		v, err := g.converters[len(g.converters)-1](args[i])
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(v))
	}
	return g.function.Call(in), nil
}

func (g *genericFunc) Accept(arity int) bool {
	if arity < g.arity {
		if g.variadic && arity == g.arity-1 {
			// having no variadic parameter is ok.
		} else {
			return false
		}

	} else if arity != g.arity && !g.variadic {
		return false
	}
	return true
}

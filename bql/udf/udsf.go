package udf

import (
	"errors"
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"reflect"
)

// UDSFCreator creates a new UDSF instance.
type UDSFCreator interface {
	// CreateUDSF returns a new UDSF instance. UDSF can be configured (via decl)
	// to accept inputs from multiple existing streams. It also receives
	// arguments passed to the UDSF in a BQL statement. The caller will call
	// UDSF.Terminate when the UDSF becomes unnecessary.
	//
	// UDSF will be discarded if it doesn't have any input.
	CreateUDSF(ctx *core.Context, decl UDSFDeclarer, args ...data.Value) (UDSF, error)

	// Accept returns true if the UDSF supports the given arity.
	Accept(arity int) bool
}

// ConvertToUDSFCreator converts a function to
func ConvertToUDSFCreator(function interface{}) (UDSFCreator, error) {
	t := reflect.TypeOf(function)
	if t.Kind() != reflect.Func {
		return nil, errors.New("the argument must be a function")
	}

	g := &genericUDSFCreator{
		function:   reflect.ValueOf(function),
		hasContext: genericFuncHasContext(t),
		variadic:   t.IsVariadic(),
		arity:      t.NumIn(),
	}

	if g.hasContext {
		g.arity -= 1
	}
	if g.arity == 0 {
		return nil, errors.New("the function must receive UDSFDeclarer as the first argument")
	}
	if c := t.In(t.NumIn() - g.arity); !reflect.TypeOf(func(UDSFDeclarer) {}).In(0).AssignableTo(c) {
		if g.hasContext {
			return nil, errors.New("the function must receive UDSFDeclarer as the second argument")
		}
		return nil, errors.New("the function must receive UDSFDeclarer as the first argument")
	}
	g.arity -= 1 // UDSFDeclarer

	if err := checkUDSFCreatorFuncReturnTypes(t); err != nil {
		return nil, err
	}

	if convs, err := createGenericConverters(t, t.NumIn()-g.arity); err != nil {
		return nil, err
	} else {
		g.converters = convs
	}
	return g, nil
}

func MustConvertToUDSFCreator(function interface{}) UDSFCreator {
	f, err := ConvertToUDSFCreator(function)
	if err != nil {
		panic(err)
	}
	return f
}

func checkUDSFCreatorFuncReturnTypes(t reflect.Type) error {
	if t.NumOut() != 2 {
		return errors.New("the function must have (udf.UDSFCreator, error) as its return type")
	}
	if !t.Out(0).Implements(reflect.TypeOf(func(UDSF) {}).In(0)) {
		return fmt.Errorf("the first return value must be a UDSF: %v", t.Out(0))
	}
	if !t.Out(1).Implements(reflect.TypeOf(func(error) {}).In(0)) {
		return fmt.Errorf("the second return value must be an error: %v", t.Out(1))
	}
	return nil
}

type genericUDSFCreator struct {
	function   reflect.Value
	hasContext bool
	variadic   bool

	// arity is the number of arguments of UDSFCreator.CreateUDSF excluding
	// *core.Context and UDSFDeclarer.
	arity int

	converters []argumentConverter
}

func (g *genericUDSFCreator) CreateUDSF(ctx *core.Context, decl UDSFDeclarer, args ...data.Value) (UDSF, error) {
	if err := g.accept(len(args)); err != nil {
		return nil, err
	}

	in := make([]reflect.Value, 0, len(args)+2) // +2 for context and declarer
	if g.hasContext {
		in = append(in, reflect.ValueOf(ctx))
	}
	in = append(in, reflect.ValueOf(decl))

	variadicBegin := g.arity
	if g.variadic {
		variadicBegin--
	}

	for i := 0; i < variadicBegin; i++ {
		v, err := g.converters[i](args[i])
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(v))
	}
	for i := variadicBegin; i < len(args); i++ {
		v, err := g.converters[len(g.converters)-1](args[i])
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(v))
	}

	out := g.function.Call(in)
	if !out[1].IsNil() {
		return nil, out[1].Interface().(error)
	}
	if out[0].IsNil() {
		return nil, fmt.Errorf("the UDSF returned from the creator was nil")
	}
	return out[0].Interface().(UDSF), nil
}

func (g *genericUDSFCreator) Accept(arity int) bool {
	return g.accept(arity) == nil
}

func (g *genericUDSFCreator) accept(arity int) error {
	if arity < g.arity {
		if g.variadic && arity == g.arity-1 {
			// having no variadic parameter is ok.
		} else {
			return errors.New("insufficient number of arguments")
		}
	} else if arity != g.arity && !g.variadic {
		return errors.New("too many arguments")
	}
	return nil
}

// UDSF is a user-defined stream-generating function. While a regular UDF
// generates a tuple from an input tuple, UDSF can generate multiple tuples
// from a tuple or even drop tuples. It's also a little different from a
// "function". UDSF can have a state. Thus, a UDSF is created by UDSFCreator
// for each invocation in a BQL statement.
//
// UDSF doesn't have Init method because initialization should be done in
// UDSFCreator.CreateUDSF.
type UDSF interface {
	// Process sends a tuple to the UDSF. This method must not block and
	// return immediately after it finished processing the received tuple.
	Process(ctx *core.Context, t *core.Tuple, w core.Writer) error

	// Terminate terminates the UDSF. Resources allocated when the UDSF is
	// created by UDFSCreator should be released in this method.
	Terminate(ctx *core.Context) error
}

// UDSFDeclarer allow UDSFs to customize their behavior.
type UDSFDeclarer interface {
	// Input adds an input from an existing stream. A UDSF doesn't have to
	// have an input. For example, there can be a UDSF which generates
	// sequential numbers and doesn't depend on any stream.
	Input(name string, config *UDSFInputConfig) error
}

// UDSFInputConfig has input configuration parameters for UDSF.
type UDSFInputConfig struct {
	// InputName is a custom name attached to incoming tuples. If this name is
	// empty, "*" will be used.
	InputName string
}

type udsfBox struct {
	f UDSF
}

var (
	_ core.StatefulBox = &udsfBox{}
)

func newUDSFBox(f UDSF) *udsfBox {
	return &udsfBox{
		f: f,
	}
}

func (b *udsfBox) Init(ctx *core.Context) error {
	return nil
}

func (b *udsfBox) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	return b.f.Process(ctx, t, w)
}

func (b *udsfBox) Terminate(ctx *core.Context) error {
	return b.f.Terminate(ctx)
}

type udsfDeclarer struct {
	inputs map[string]*UDSFInputConfig
}

func newUDSFDeclarer() *udsfDeclarer {
	return &udsfDeclarer{
		inputs: map[string]*UDSFInputConfig{},
	}
}

func (d *udsfDeclarer) Input(name string, config *UDSFInputConfig) error {
	if _, ok := d.inputs[name]; ok {
		return fmt.Errorf("the input is already added: %v", name)
	}
	if config == nil {
		config = &UDSFInputConfig{}
	}
	d.inputs[name] = config
	return nil
}

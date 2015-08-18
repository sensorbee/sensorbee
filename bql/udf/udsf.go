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
	// When a UDSF needs inputs from other sources or streams, call Input method
	// of UDSFDeclarer. For example let's assume there's a UDSF below:
	//
	//	func createMyUDSF(decl udf.UDSFDeclarer, input1, input2 string) (udf.UDSF, error) {
	//		...
	//		decl.Input(input1, &udf.UDSFInputConfig{
	//			InputName: "custom_input_name_1",
	//		})
	//		decl.Input(input2, &udf.UDSFInputConfig{
	//			InputName: "custom_input_name_2",
	//		})
	//		...
	//	}
	//
	//	func init() {
	//		udf.MustRegisterGlobalUDSFCreator("my_udsf",
	//			udf.MustConvertToUDSFCreator(createMyUDSF))
	//	}
	//
	// Then a user can specify input stream by the following statement:
	//
	//	CREATE STREAM stream1 AS SELECT ...;
	//	CREATE STREAM stream2 AS SELECT ...;
	//	CREATE STREAM join_by_udsf
	//	  SELECT RSTREAM * FROM my_udsf('stream1', 'stream2') [RANGE 1 TUPLES];
	//
	// In this example, my_udsf receives two streams: stream1 and stream2
	// created in advance.
	//
	// A UDSF doesn't have to have an input. For example, there can be a UDSF
	// which generates sequential numbers and doesn't depend on any stream.
	// Such UDSFs will be run as the source mode. See UDSF for more details.
	CreateUDSF(ctx *core.Context, decl UDSFDeclarer, args ...data.Value) (UDSF, error)

	// Accept returns true if the UDSF supports the given arity.
	Accept(arity int) bool
}

// ConvertToUDSFCreator converts a function to a UDSFCreator.
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
		g.arity--
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
	g.arity-- // UDSFDeclarer

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

// MustConvertToUDSFCreator converts a function to a UDSFCreator.
// It panics if there is an error during conversion.
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
//
// There're two kinds of processing modes of UDSFs: stream mode and source mode.
// A UDSF is processed in the stream mode when it has at least one input stream.
// Otherwise, the UDSF is processed in the source mode. Process and Terminate
// methods change their behavior based on the mode. See documentation of each
// method for details.
type UDSF interface {
	// Process sends a tuple to the UDSF.
	//
	// When the UDSF is running in the stream mode, this method must not block
	// and return immediately after it finished processing the received tuple.
	// It is called everytime a tuple is received from streams. It behaves just
	// like core.Box.
	//
	// When the UDSF is running in the source mode, this method is only called
	// once. It can block until it generates all tuples or Terminate method
	// is called. A tuple passed to this method in the source mode doesn't
	// contain anything meaningful. It behaves like core.Source although the
	// interface is like core.Box. The core.Writer returns core.ErrSourceStopped
	// when the UDSF running in the source mode is stopped. Therefore, if
	// Process method returns on that error, the implementation of Terminate
	// can just be resource deallocation.
	Process(ctx *core.Context, t *core.Tuple, w core.Writer) error

	// Terminate terminates the UDSF. Resources allocated when the UDSF is
	// created by UDFSCreator should be released in this method. Also, when
	// the UDSF is running in the source mode, Process method must return as
	// soon as Terminate is called.
	Terminate(ctx *core.Context) error
}

// UDSFDeclarer allow UDSFs to customize their behavior.
type UDSFDeclarer interface {
	// Input adds an input from an existing stream.
	Input(name string, config *UDSFInputConfig) error

	// ListInputs returns all inputs declared by a UDSF. The caller can safely
	// modify the map returned from this method.
	ListInputs() map[string]*UDSFInputConfig
}

// UDSFInputConfig has input configuration parameters for UDSF.
type UDSFInputConfig struct {
	// InputName is a custom name attached to incoming tuples. If this name is
	// empty, "*" will be used.
	InputName string
}

type udsfDeclarer struct {
	inputs map[string]*UDSFInputConfig
}

func newUDSFDeclarer() *udsfDeclarer {
	return &udsfDeclarer{
		inputs: map[string]*UDSFInputConfig{},
	}
}

// NewUDSFDeclarer creates a new declarer.
func NewUDSFDeclarer() UDSFDeclarer {
	return newUDSFDeclarer()
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

func (d *udsfDeclarer) ListInputs() map[string]*UDSFInputConfig {
	m := make(map[string]*UDSFInputConfig, len(d.inputs))
	for i, c := range d.inputs {
		m[i] = c
	}
	return m
}

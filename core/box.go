package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

// A Box is an elementary building block of a SensorBee topology.
// It is the equivalent of a StreamTask in Samza or a Bolt in Storm.
type Box interface {
	// Process is called on a Box for each item in the input stream.
	// The processed result must be written to the given Writer
	// object.
	// It is ok to modify the given Tuple object in place and
	// pass it on to the Writer. However, it is *not* ok to keep
	// a reference to the written object and then access or modify
	// it at a later point in time.
	// Note that there may be multiple concurrent calls to Process,
	// so if internal state is accessed, proper locking mechanisms
	// must be used.
	// Also note that the same Context pointer that was passed to
	// Process should be handed over to Writer.Write so that the
	// same settings will be used by that Writer.
	Process(ctx *Context, t *tuple.Tuple, w Writer) error
}

// StatefulBox is a Box having its internal state and needs to be initialized
// before it's used by a topology. Because a Box can be implemented in C or
// C++, Terminate method is also provided to deallocate resources it used
// during processing.
type StatefulBox interface {
	Box

	// Init is called on each Box in a Topology when StaticTopology.Run()
	// is executed. It can be used to keep a reference to the Context
	// object or initialize other forms of state.
	Init(ctx *Context) error

	// Terminate finalizes a Box. The Box can no longer be used after
	// this method is called. This method doesn't have to be idempotent,
	// that is the behavior is undefined when this method is called
	// more than once. Terminate isn't called if Init fails on the Box.
	//
	// As long as the Box is used by components in core package, no method
	// will be called after the invocation of Terminate method. In addition,
	// the Box may assume that Terminate is not called concurrently with
	// any other methods.
	Terminate(ctx *Context) error
}

// SchemaSet is a set of schemas required or computed by a SchemafulBox.
// A schema can be nil, which means the input or output is schemaless and
// accept any type of tuples.
type SchemaSet map[string]*Schema

// Names returns names of all the schema the set has. It returns names
// of nil-schemas.
func (s SchemaSet) Names() []string {
	var ns []string
	for n := range s {
		ns = append(ns, n)
	}
	return ns
}

// Has returns true when the set has the schema with the name. It returns
// true even when the set has a nil-schema with the name.
func (s SchemaSet) Has(name string) bool {
	_, ok := s[name]
	return ok
}

// SchemafulBox is a Box having input/output schema information so that type
// conversion error can be detected before it's executed.
type SchemafulBox interface {
	Box

	// InputSchema returns the specification of tuples which the Box receives.
	//
	// Schema is a map that declares the required input schema for
	// each logical input stream, identified by a Box-internal input
	// name. Most boxes will treat all input tuples the same way
	// (e.g., filter by a certain value, aggregate a certain key etc.);
	// such boxes do not require the concept of multiple logical streams
	// and should return
	//
	//	{"*": mySchema}
	//
	// where mySchema is either nil (to signal that any input is fine)
	// or a pointer to a valid Schema object (to signal that all input
	// tuples must match that schema). However, boxes that need to
	// treat tuples from one source different than tuples from another
	// source (e.g., when joining on `user.id == events.uid` etc.),
	// can declare that they require multiple streams with (maybe)
	// different schemas:
	//
	//	{"left": schemaA, "right": schemaB}
	//
	// If schemaA or schemaB is nil, then any input is fine on that
	// logical stream, otherwise it needs to match the given schema.
	// It is also possible to mix the forms using
	//
	//	{"special": schemaA, "*": schemaB}
	//
	// in which streams with arbitrary names can be added, but the tuples
	// in streams called "special" will have to match schemaA (if given),
	// all others schemaB (if given).
	// Also see BoxDeclarer.Input and BoxDeclarer.NamedInput.
	InputSchema() SchemaSet

	// OutputSchema is called on a Box when it is necessary to
	// determine the shape of data that will be created by this Box.
	// Since this output may vary depending on the data that is
	// received, a list of the schemas of the input streams is
	// passed as a parameter. Each schema in the argument is tied
	// to the name of the input defined in the schema returned from
	// InputSchema method. Return a nil pointer to signal that there
	// are no guarantees  about the shape of the returned data given.
	OutputSchema(SchemaSet) (*Schema, error)
}

// TODO: Support input constraints such as an acceptable frequency of tuples.

// BoxFunc can be used to add all methods required to fulfill the Box
// interface to a normal function with the signature
//   func(ctx *Context, t *tuple.Tuple, s Writer) error
//
// Example:
//
//     forward := func(ctx *Context, t *tuple.Tuple, w Writer) error {
//         w.Write(ctx, t)
//         return nil
//     }
//     var box Box = BoxFunc(forward)
func BoxFunc(b func(ctx *Context, t *tuple.Tuple, w Writer) error) Box {
	return boxFunc(b)
}

type boxFunc func(ctx *Context, t *tuple.Tuple, w Writer) error

func (b boxFunc) Process(ctx *Context, t *tuple.Tuple, w Writer) error {
	return b(ctx, t, w)
}

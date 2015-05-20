package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

// BoxInputConstraints describe requirements that a Box has towards
// its input data. At the moment, this is only the schema that
// tuples need to fulfill.
type BoxInputConstraints struct {
	// Schema is a map that declares the required input schema for
	// each logical input stream, identified by a Box-internal input
	// name. Most boxes will treat all input tuples the same way
	// (e.g., filter by a certain value, aggregate a certain key etc.);
	// such boxes do not require the concept of multiple logical streams
	// and should return
	//   {"*": mySchema}
	// where mySchema is either nil (to signal that any input is fine)
	// or a pointer to a valid Schema object (to signal that all input
	// tuples must match that schema). However, boxes that need to
	// treat tuples from one source different than tuples from another
	// source (e.g., when joining on `user.id == events.uid` etc.),
	// can declare that they require multiple streams with (maybe)
	// different schemas:
	//  {"left": schemaA, "right": schemaB}
	// If schemaA or schemaB is nil, then any input is fine on that
	// logical stream, otherwise it needs to match the given schema.
	// It is also possible to mix the forms using
	//  {"special": schemaA, "*": schemaB}
	// in which streams with arbitrary names can be added, but the tuples
	// in streams called "special" will have to match schemaA (if given),
	// all others schemaB (if given).
	// Also see BoxDeclarer.Input and BoxDeclarer.NamedInput.
	Schema map[string]*Schema
}

// A Box is an elementary building block of a SensorBee topology.
// It is the equivalent of a StreamTask in Samza or a Bolt in Storm.
type Box interface {
	// Init is called on each Box in a Topology when StaticTopology.Run()
	// is executed. It can be used to keep a reference to the Context
	// object or initialize other forms of state. It is called only
	// once, even if the same object is used in multiple places of
	// the topology.
	Init(ctx *Context) error

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
	Process(ctx *Context, t *tuple.Tuple, s Writer) error

	// InputConstraints is called on a Box to learn about the
	// requirements of this Box with respect to its input data. If
	// this function returns nil, it means that any data can serve
	// as input for this Box. Also see the documentation for the
	// BoxInputConstraints struct.
	InputConstraints() (*BoxInputConstraints, error)

	// OutputSchema is called on a Box when it is necessary to
	// determine the shape of data that will be created by this Box.
	// Since this output may vary depending on the data that is
	// received, a list of the schemas of the input streams is
	// passed as a parameter. Return a nil pointer to signal that
	// there are no guarantees about the shape of the returned data
	// given.
	OutputSchema([]*Schema) (*Schema, error)

	// Terminate finalizes a Box. The Box can no longer be used after
	// this method is called. This method doesn't have to be idempotent,
	// that is the behavior is undefined when this method is called
	// more than once.
	Terminate(ctx *Context) error
}

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
func BoxFunc(f func(ctx *Context, t *tuple.Tuple, s Writer) error) Box {
	bf := boxFunc(f)
	return &bf
}

type boxFunc func(ctx *Context, t *tuple.Tuple, s Writer) error

func (b *boxFunc) Process(ctx *Context, t *tuple.Tuple, s Writer) error {
	return (*b)(ctx, t, s)
}

func (b *boxFunc) Init(ctx *Context) error {
	return nil
}

func (b *boxFunc) InputConstraints() (*BoxInputConstraints, error) {
	return nil, nil
}

func (b *boxFunc) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

func (b *boxFunc) Terminate(ctx *Context) error {
	return nil
}

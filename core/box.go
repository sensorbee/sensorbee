package core

import (
	"fmt"
	"strings"
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
	//
	// Process can tell the caller about each error in more detail by
	// returning behavioral error types. The caller of Process must
	// follow three rules:
	//
	//	1. The caller must not call Process again if the error is fatal
	//	   (i.e. IsFatalError(err) == true). Because there can be multiple
	//	   goroutines on a single Box, this rule only applies to the goroutine
	//	   which received a fatal error.
	//	2. The caller may retry the same tuple or discard it if the error is
	//	   temporary but not fatal (i.e. IsFatalError(err) == false && IsTemporaryError(err) == true).
	//	   The caller may call Process again with the same tuple, or may even
	//	   discard the tuple and skip it. The number of retry which the caller
	//	   should attempt is not defined.
	//	3. The caller must discard the tuple and must not retry if the error
	//	   isn't temporary nor fatal (i.e. IsFatalError(err) == false && IsTemporaryError(err) == false).
	//	   The caller can call Process again with a different tuple, that is
	//	   it can just skip the tuple which Process returned the error.
	//
	// Once Process returns a fatal error, it must always return fatal errors
	// after that. Process might be called even after it returned a fatal error.
	// Terminate method will be called even if Process returns a fatal error.
	//
	// When Process returns a temporary error, the call shouldn't have any
	// side effect on the state of the Box,  that is consecutive retry calls
	// should only be reflected once regardless of the number of attempt the
	// caller makes. For example, let's assume a Box B is counting the word
	// in tuples. Then, a caller invoked B.Process with a tuple t. If B.Process
	// with t returned a temporary error, the count B has shouldn't be changed
	// until the retry succeeds.
	Process(ctx *Context, t *Tuple, w Writer) error
}

// StatefulBox is a Box having an internal state that needs to be initialized
// before it's used by a topology. Because a Box can be implemented in C or
// C++, a Terminate method is also provided to deallocate resources it used
// during processing.
//
// An instance of StatefulBox shouldn't be added to a topology or a topology
// builder more than once if it doesn't handle duplicated initialization and
// termination correctly (i.e. with something like a reference counter).
type StatefulBox interface {
	Box

	// Init is called on each Box in a Topology when Topology.Add is called.
	// It can be used to keep a reference to the Context object or initialize
	// other forms of state.
	//
	// When the same instance of the box is added to a topology more than
	// once, Init will be called multiple times.
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
	//
	// When the same instance of the box is added to a topology more than
	// once, Terminate will be called multiple times. In that case, use
	// something like a reference counter to avoid deallocating resources
	// before the Box is actually removed from a topology.
	Terminate(ctx *Context) error
}

// TODO: Support input constraints such as an acceptable frequency of tuples.

// NamedInputBox is a box whose inputs have custom input names.
type NamedInputBox interface {
	Box

	// InputNames returns a slice of names which the box accepts as input names.
	// If it's empty, any name is accepted. Otherwise, the box only accepts
	// names in the slice.
	InputNames() []string
}

// BoxFunc can be used to add all methods required to fulfill the Box
// interface to a normal function with the signature
//   func(ctx *Context, t *Tuple, s Writer) error
//
// Example:
//
//     forward := func(ctx *Context, t *Tuple, w Writer) error {
//         w.Write(ctx, t)
//         return nil
//     }
//     var box Box = BoxFunc(forward)
func BoxFunc(b func(ctx *Context, t *Tuple, w Writer) error) Box {
	return boxFunc(b)
}

type boxFunc func(ctx *Context, t *Tuple, w Writer) error

func (b boxFunc) Process(ctx *Context, t *Tuple, w Writer) error {
	return b(ctx, t, w)
}

// checkBoxInputName checks if the Box accept the given input name. It returns
// an error when the Box doesn't accept the name.
func checkBoxInputName(b Box, boxName string, inputName string) error {
	// The `Input()` caller said that we should attach the name
	// `inputName` to incoming data (or not if inputName is "*").
	// This is ok if
	// - A Box doesn't have named inputs.
	// - InputNames() is nil or empty
	// - InputNames() contains the name
	// Otherwise this is an error.
	nbox, ok := b.(NamedInputBox)
	if !ok {
		return nil // This box accepts any name
	}

	names := nbox.InputNames()
	if len(names) == 0 { // names == nil or has zero length
		return nil
	}

	for _, n := range names {
		if n == inputName {
			return nil
		}
	}
	return fmt.Errorf("an input name %s isn't defined in the box '%v': %v",
		inputName, boxName, strings.Join(names, ", "))
}

// boxWriterAdapter provides a Writer interface which writes tuples to a Box.
// It also records traces input and output tuples.
type boxWriterAdapter struct {
	box  Box
	name string
	dst  *traceWriter
}

func newBoxWriterAdapter(b Box, name string, dst WriteCloser) *boxWriterAdapter {
	return &boxWriterAdapter{
		box:  b,
		name: name,
		// An output traces is written just after the box Process writes a tuple.
		dst: newTraceWriter(dst, Output, name),
	}
}

func (wa *boxWriterAdapter) Write(ctx *Context, t *Tuple) error {
	tracing(t, ctx, Input, wa.name)
	return wa.box.Process(ctx, t, wa.dst)
}

func (wa *boxWriterAdapter) Close(ctx *Context) error {
	// TODO: handle panics
	var errb error
	if sbox, ok := wa.box.(StatefulBox); ok {
		errb = sbox.Terminate(ctx)
	}
	errw := wa.dst.w.Close(ctx)
	if errb != nil {
		return errb // An error from the Box is considered more important.
	}
	return errw
}

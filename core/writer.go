package core

import (
	"pfi/sensorbee/sensorbee/tuple"
)

// Writer describes an object that tuples can be written to
// as the output for a Box. Note that this interface was chosen
// because it also allows a Box to write multiple (or none)
// output tuples. It is expected that the ctx pointer passed in
// points to the same Context that was used by the Box that
// called Write.
type Writer interface {
	Write(ctx *Context, t *tuple.Tuple) error
}

// WriteCloser add a capability of closing to Writer.
type WriteCloser interface {
	Writer

	// Close closes the writer. An appropriate Context should be given,
	// which is usually provided by Topology. Close doesn't have to be
	// idempotent.
	Close(ctx *Context) error
}

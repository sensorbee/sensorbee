package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

// A Writer describes an object that tuples can be written to
// as the output for a Box. Note that this interface was chosen
// because it also allows a Box to write multiple (or none)
// output tuples.
type Writer interface {
	Write(ctx *Context, t *tuple.Tuple) error
}

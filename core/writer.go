package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

// A Writer describes an object that tuples can be written to
// as the output for a Box. Note that this interface was chosen
// because it also allows a Box to write multiple (or none)
// output tuples. It is expected that the ctx pointer passed in
// points to the same Context that was used by the Box that
// called Write.
type Writer interface {
	Write(ctx *Context, t *tuple.Tuple) error
}

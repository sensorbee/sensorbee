package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type LogManager interface {
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type Context interface {
	Log(level LogLevel, msg string, a ...interface{})

	// If there was an error during processing in a box and
	// you cannot process a tuple further, report this here.
	DroppedTuple(t *tuple.Tuple, msg string, a ...interface{})
}

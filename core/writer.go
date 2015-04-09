package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type Writer interface {
	Write(t *tuple.Tuple) error
}

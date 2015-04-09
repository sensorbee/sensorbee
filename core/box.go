package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type Box interface {
	Process(t *tuple.Tuple, s Sink) error
	RequiredInputSchema() ([]*Schema, error)
	OutputSchema([]*Schema) (*Schema, error)
}

type BoxFunc func(t *tuple.Tuple, s Sink) error

func (b *BoxFunc) Process(t *tuple.Tuple, s Sink) error {
	return (*b)(t, s)
}

func (b *BoxFunc) RequiredInputSchema() ([]*Schema, error) {
	return nil, nil
}

func (b *BoxFunc) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type Box interface {
	Init(ctx Context) error
	Process(t *tuple.Tuple, s Writer) error
	RequiredInputSchema() ([]*Schema, error)
	OutputSchema([]*Schema) (*Schema, error)
}

type BoxFunc func(t *tuple.Tuple, s Writer) error

func (b *BoxFunc) Process(t *tuple.Tuple, s Writer) error {
	return (*b)(t, s)
}

func (b *BoxFunc) Init(ctx Context) error {
	return nil
}

func (b *BoxFunc) RequiredInputSchema() ([]*Schema, error) {
	return nil, nil
}

func (b *BoxFunc) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

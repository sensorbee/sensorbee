package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type BoxInputConstraints struct {
	Schema map[string]*Schema
}

type Box interface {
	Init(ctx *Context) error
	Process(t *tuple.Tuple, s Writer) error
	// TODO write comment
	InputConstraints() (*BoxInputConstraints, error)
	OutputSchema([]*Schema) (*Schema, error)
}

func BoxFunc(f func(t *tuple.Tuple, s Writer) error) Box {
	bf := boxFunc(f)
	return &bf
}

type boxFunc func(t *tuple.Tuple, s Writer) error

func (b *boxFunc) Process(t *tuple.Tuple, s Writer) error {
	return (*b)(t, s)
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

package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type InputConstraints struct {
	Schema map[string]*Schema
}

type Box interface {
	Init(ctx *Context) error
	Process(t *tuple.Tuple, s Writer) error
	InputConstraints() (*InputConstraints, error)
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

func (b *boxFunc) InputConstraints() (*InputConstraints, error) {
	return nil, nil
}

func (b *boxFunc) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

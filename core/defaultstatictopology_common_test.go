package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

// DoesNothingSource is a dummy source that literally does nothing.
// It just fulfills the Source interface so that we can build a
// simple topology.
type DoesNothingSource struct{}

func (s *DoesNothingSource) GenerateStream(w Writer) error {
	return nil
}
func (s *DoesNothingSource) Schema() *Schema {
	return nil
}

/**************************************************/

// DoesNothingBox is a dummy source that literally does nothing.
// It just fulfills the Box interface so that we can build a
// simple topology.
type DoesNothingBox struct {
}

func (b *DoesNothingBox) Init(ctx *Context) error {
	return nil
}
func (b *DoesNothingBox) Process(t *tuple.Tuple, s Writer) error {
	return nil
}
func (b *DoesNothingBox) InputConstraints() (*BoxInputConstraints, error) {
	return nil, nil
}
func (b *DoesNothingBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

/**************************************************/

// DoesNothingSink is a dummy source that literally does nothing.
// It just fulfills the Sink interface so that we can build a
// simple topology.
type DoesNothingSink struct{}

func (s *DoesNothingSink) Write(t *tuple.Tuple) error {
	return nil
}

/**************************************************/

// TupleEmitterSource is a source that emits all tuples in the given
// slice once when GenerateStream is called.
type TupleEmitterSource struct {
	Tuples []*tuple.Tuple
}

func (s *TupleEmitterSource) GenerateStream(w Writer) error {
	for _, t := range s.Tuples {
		w.Write(t)
	}
	return nil
}
func (s *TupleEmitterSource) Schema() *Schema {
	return nil
}

type TupleCollectorSink struct {
	Tuples []*tuple.Tuple
}

func (s *TupleCollectorSink) Write(t *tuple.Tuple) error {
	s.Tuples = append(s.Tuples, t)
	return nil
}

/**************************************************/

// forwardBox sends an input Tuple to the given Writer without
// modification. It can be wrapped with BoxFunc to match the Box
// interface.
func forwardBox(t *tuple.Tuple, w Writer) error {
	w.Write(t)
	return nil
}

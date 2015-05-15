package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
)

func newTestContext(config Configuration) *Context {
	return &Context{
		Logger: NewConsolePrintLogger(),
		Config: config,
	}
}

// DoesNothingSource is a dummy source that literally does nothing.
// It just fulfills the Source interface so that we can build a
// simple topology.
type DoesNothingSource struct{}

func (s *DoesNothingSource) GenerateStream(ctx *Context, w Writer) error {
	return nil
}
func (s *DoesNothingSource) Stop(ctx *Context) error {
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
func (b *DoesNothingBox) Process(ctx *Context, t *tuple.Tuple, s Writer) error {
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

func (s *DoesNothingSink) Write(ctx *Context, t *tuple.Tuple) error {
	return nil
}
func (s *DoesNothingSink) Close(ctx *Context) error {
	return nil
}

/**************************************************/

// TupleEmitterSource is a source that emits all tuples in the given
// slice once when GenerateStream is called.
type TupleEmitterSource struct {
	Tuples []*tuple.Tuple
	m      sync.Mutex
	c      *sync.Cond

	// 0: running, 1: stopping, 2: stopped
	state int
}

func (s *TupleEmitterSource) GenerateStream(ctx *Context, w Writer) error {
	s.c = sync.NewCond(&s.m)
	for _, t := range s.Tuples {
		s.m.Lock()
		if s.state > 0 {
			s.state = 2
			s.c.Broadcast()
			s.m.Unlock()
			break
		}
		s.m.Unlock()

		w.Write(ctx, t)
	}
	return nil
}
func (s *TupleEmitterSource) Stop(ctx *Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.state == 2 {
		return nil
	}
	s.state = 1
	for s.state < 2 {
		s.c.Wait()
	}
	return nil
}
func (s *TupleEmitterSource) Schema() *Schema {
	return nil
}

type TupleIncrementalEmitterSource struct {
	Tuples []*tuple.Tuple

	m     sync.Mutex
	c     *sync.Cond
	cnt   int
	state int // 0: running, 1: stopping, 2: stopped
}

func NewTupleIncrementalEmitterSource(ts []*tuple.Tuple) *TupleIncrementalEmitterSource {
	s := &TupleIncrementalEmitterSource{
		Tuples: ts,
	}
	s.c = sync.NewCond(&s.m)
	return s
}

func (s *TupleIncrementalEmitterSource) GenerateStream(ctx *Context, w Writer) error {
	for _, t := range s.Tuples {
		s.m.Lock()
		for s.state == 0 && s.cnt == 0 {
			s.c.Wait()
		}
		if s.state != 0 {
			s.state = 2
			s.c.Broadcast()
			s.m.Unlock()
			return nil
		}

		func() {
			defer s.m.Unlock()
			w.Write(ctx, t)
			s.cnt--
			s.c.Broadcast()
		}()
	}
	s.m.Lock()
	defer s.m.Unlock()
	s.state = 2
	s.c.Broadcast()
	return nil
}

// EmitTuples emits n tuples and waits until all tuples are emitted.
// It doesn't wait until tuples are fully processed by the entire topology.
func (s *TupleIncrementalEmitterSource) EmitTuples(n int) {
	s.m.Lock()
	defer s.m.Unlock()
	s.cnt += n
	s.c.Broadcast()

	for s.cnt > 0 {
		s.c.Wait()
	}
}
func (s *TupleIncrementalEmitterSource) Stop(ctx *Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.state == 2 {
		return nil
	}
	s.state = 1
	s.c.Broadcast()

	for s.state < 2 {
		s.c.Wait()
	}
	return nil
}
func (s *TupleIncrementalEmitterSource) Schema() *Schema {
	return nil
}

type TupleCollectorSink struct {
	Tuples []*tuple.Tuple
	m      sync.Mutex
	c      *sync.Cond
}

func NewTupleCollectorSink() *TupleCollectorSink {
	s := &TupleCollectorSink{}
	s.c = sync.NewCond(&s.m)
	return s
}

func (s *TupleCollectorSink) Write(ctx *Context, t *tuple.Tuple) error {
	if s.c == nil { // This is for old tests
		s.Tuples = append(s.Tuples, t)
		return nil
	}
	s.m.Lock()
	defer s.m.Unlock()
	s.Tuples = append(s.Tuples, t)
	s.c.Broadcast()
	return nil
}

// Wait waits until the collector receives at least n tuples.
func (s *TupleCollectorSink) Wait(n int) {
	s.m.Lock()
	defer s.m.Unlock()
	for len(s.Tuples) < n {
		s.c.Wait()
	}
}
func (s *TupleCollectorSink) Close(ctx *Context) error {
	return nil
}

/**************************************************/

// forwardBox sends an input Tuple to the given Writer without
// modification. It can be wrapped with BoxFunc to match the Box
// interface.
func forwardBox(ctx *Context, t *tuple.Tuple, w Writer) error {
	w.Write(ctx, t)
	return nil
}

package bql

import (
	. "pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
)

func newTestContext(config Configuration) *Context {
	return &Context{
		Logger: NewConsolePrintLogger(),
		Config: config,
	}
}

//////// Below are all copies from the core package ////////

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

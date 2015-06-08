package dummy_sink

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
)

func init() {
	// register with BQL
	bql.RegisterSinkType("collector", CreateCollectorSink)
}

// CreateCollectorSink creates a sink that collects all received
// tuples in an internal array.
func CreateCollectorSink(params map[string]string) (core.Sink, error) {
	// check the given sink parameters
	for key, _ := range params {
		return nil, fmt.Errorf("unknown sink parameter: %s", key)
	}
	si := TupleCollectorSink{}
	si.c = sync.NewCond(&si.m)
	return &si, nil
}

type TupleCollectorSink struct {
	Tuples []*tuple.Tuple
	m      sync.Mutex
	c      *sync.Cond
}

func (s *TupleCollectorSink) Write(ctx *core.Context, t *tuple.Tuple) error {
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
func (s *TupleCollectorSink) Close(ctx *core.Context) error {
	return nil
}

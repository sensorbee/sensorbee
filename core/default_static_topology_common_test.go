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
func (b *DoesNothingBox) Terminate(ctx *Context) error {
	return nil
}

// ProxyBox just passes all method calls to the target Box.
type ProxyBox struct {
	b Box
}

func (b *ProxyBox) Init(ctx *Context) error {
	return b.b.Init(ctx)
}
func (b *ProxyBox) Process(ctx *Context, t *tuple.Tuple, s Writer) error {
	return b.b.Process(ctx, t, s)
}
func (b *ProxyBox) InputConstraints() (*BoxInputConstraints, error) {
	return b.b.InputConstraints()
}
func (b *ProxyBox) OutputSchema(s []*Schema) (*Schema, error) {
	return b.b.OutputSchema(s)
}
func (b *ProxyBox) Terminate(ctx *Context) error {
	return b.b.Terminate(ctx)
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

func NewTupleEmitterSource(ts []*tuple.Tuple) *TupleEmitterSource {
	s := &TupleEmitterSource{
		Tuples: ts,
	}
	s.c = sync.NewCond(&s.m)
	return s
}

func (s *TupleEmitterSource) GenerateStream(ctx *Context, w Writer) error {
	if s.c == nil {
		s.c = sync.NewCond(&s.m)
	}
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
	s.m.Lock()
	defer s.m.Unlock()
	s.state = 2
	s.c.Broadcast()
	return nil
}
func (s *TupleEmitterSource) Stop(ctx *Context) error {
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

		// To make writing tests easy, GenerateStream doesn't stop until
		// cnt get 0 or all tuples in s.Tuples are emitted.
		if s.cnt == 0 && s.state != 0 {
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
	s.cnt = 0
	s.c.Broadcast()
	return nil
}

// EmitTuples emits n tuples and waits until all tuples are emitted.
// It doesn't wait until tuples are fully processed by the entire topology.
func (s *TupleIncrementalEmitterSource) EmitTuples(n int) {
	s.EmitTuplesNB(n)
	s.WaitForEmission()
}

// EmitTupelsNB emits n tuples asynchronously and doesn't wait until
// all tuples are emitted.
func (s *TupleIncrementalEmitterSource) EmitTuplesNB(n int) {
	s.m.Lock()
	defer s.m.Unlock()
	s.cnt += n
	s.c.Broadcast()
}

// WaitForEmission waits until all tuples specified by EmitTuples(NB) are emitted.
func (s *TupleIncrementalEmitterSource) WaitForEmission() {
	s.m.Lock()
	defer s.m.Unlock()
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

// BlockingForwardBox blocks tuples in it and emit them at the timing required by users.
// BlockingForwardBox.c won't be initialized until Init method is called.
// It's safe to call EmitTuples after the topology's state becomes TSRunning.
type BlockingForwardBox struct {
	m   sync.Mutex
	c   *sync.Cond
	cnt int
}

func (b *BlockingForwardBox) Init(ctx *Context) error {
	b.c = sync.NewCond(&b.m)
	return nil
}

func (b *BlockingForwardBox) Process(ctx *Context, t *tuple.Tuple, w Writer) error {
	b.m.Lock()
	defer b.m.Unlock()
	for b.cnt == 0 {
		b.c.Wait()
	}
	b.cnt--
	b.c.Broadcast()
	w.Write(ctx, t)
	return nil
}

// EmitTuples emits n tuples. If EmitTuples(3) is called when the box isn't
// blocking any tuple, Process method won't block for three times.
// EmitTuples doesn't wait until all n tuples are emitted.
func (b *BlockingForwardBox) EmitTuples(n int) {
	b.m.Lock()
	defer b.m.Unlock()
	b.cnt += n
	b.c.Broadcast()
}

func (b *BlockingForwardBox) InputConstraints() (*BoxInputConstraints, error) {
	return nil, nil
}

func (b *BlockingForwardBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

func (b *BlockingForwardBox) Terminate(ctx *Context) error {
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

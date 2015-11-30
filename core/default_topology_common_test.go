package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/data"
	"sync"
)

// DoesNothingSource is a dummy source that literally does nothing.
// It just fulfills the Source interface so that we can build a
// simple topology.
type DoesNothingSource struct {
	// dummy is used to give a different address for each instance of this struct.
	dummy int
}

func (s *DoesNothingSource) GenerateStream(ctx *Context, w Writer) error {
	return nil
}
func (s *DoesNothingSource) Stop(ctx *Context) error {
	return nil
}

/**************************************************/

// DoesNothingBox is a dummy source that literally does nothing.
// It just fulfills the Box interface so that we can build a
// simple topology.
type DoesNothingBox struct {
	// dummy is used to give a different address for each instance of this struct.
	dummy int
}

func (b *DoesNothingBox) Process(ctx *Context, t *Tuple, s Writer) error {
	return nil
}

// ProxyBox just passes all method calls to the target Box.
type ProxyBox struct {
	b Box
}

var (
	_ StatefulBox   = &ProxyBox{}
	_ NamedInputBox = &ProxyBox{}
)

func (b *ProxyBox) Init(ctx *Context) error {
	if s, ok := b.b.(StatefulBox); ok {
		return s.Init(ctx)
	}
	return nil
}

func (b *ProxyBox) Process(ctx *Context, t *Tuple, s Writer) error {
	return b.b.Process(ctx, t, s)
}

func (b *ProxyBox) Terminate(ctx *Context) error {
	if s, ok := b.b.(StatefulBox); ok {
		return s.Terminate(ctx)
	}
	return nil
}

func (b *ProxyBox) InputNames() []string {
	if s, ok := b.b.(NamedInputBox); ok {
		return s.InputNames()
	}
	return nil
}

/**************************************************/

// DoesNothingSink is a dummy source that literally does nothing.
// It just fulfills the Sink interface so that we can build a
// simple topology.
type DoesNothingSink struct {
	// dummy is used to give a different address for each instance of this struct.
	dummy int
}

func (s *DoesNothingSink) Write(ctx *Context, t *Tuple) error {
	return nil
}
func (s *DoesNothingSink) Close(ctx *Context) error {
	return nil
}

/**************************************************/

// TupleEmitterSource is a source that emits all tuples in the given
// slice once when GenerateStream is called.
type TupleEmitterSource struct {
	Tuples []*Tuple
	m      sync.Mutex
	c      *sync.Cond

	// 0: running, 1: stopping, 2: stopped
	state int
}

func NewTupleEmitterSource(ts []*Tuple) *TupleEmitterSource {
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
	s.m.Lock()
	s.state = 0
	s.m.Unlock()

	defer func() {
		s.m.Lock()
		defer s.m.Unlock()
		s.state = 2
		s.c.Broadcast()
	}()

	for _, t := range s.Tuples {
		s.m.Lock()
		if s.state > 0 {
			s.state = 2
			s.c.Broadcast()
			s.m.Unlock()
			break
		}
		s.m.Unlock()

		if err := w.Write(ctx, t); err != nil {
			if err == ErrSourceRewound || err == ErrSourceStopped {
				return err
			}
		}
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
	s.c.Broadcast()
	for s.state < 2 {
		s.c.Wait()
	}
	return nil
}

type TupleIncrementalEmitterSource struct {
	Tuples []*Tuple

	m     sync.Mutex
	state *topologyStateHolder
	cnt   int
}

var (
	_ Resumable = &TupleIncrementalEmitterSource{}
)

func NewTupleIncrementalEmitterSource(ts []*Tuple) *TupleIncrementalEmitterSource {
	s := &TupleIncrementalEmitterSource{
		Tuples: ts,
	}
	s.state = newTopologyStateHolder(&s.m)
	return s
}

func (s *TupleIncrementalEmitterSource) GenerateStream(ctx *Context, w Writer) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.state.setWithoutLock(TSRunning)

	for _, t := range s.Tuples {
		for {
			if s.state.state == TSPaused {
				s.state.waitWithoutLock(TSRunning)

			} else if s.state.state == TSRunning {
				if s.cnt != 0 {
					break
				}
				s.state.cond.Wait()

			} else {
				break
			}
		}

		// To make writing tests easy, GenerateStream doesn't stop until
		// cnt get 0 or all tuples in s.Tuples are emitted.
		if s.cnt == 0 && s.state.state >= TSStopping {
			break
		}

		w.Write(ctx, t)
		s.cnt--
		s.state.cond.Broadcast()
	}

	s.cnt = 0
	s.state.setWithoutLock(TSStopped)
	return nil
}

// EmitTuples emits n tuples and waits until all tuples are emitted.
// It doesn't wait until tuples are fully processed by the entire topology.
func (s *TupleIncrementalEmitterSource) EmitTuples(n int) {
	s.EmitTuplesNB(n)
	s.WaitForEmission()
}

// EmitTuplesNB emits n tuples asynchronously and doesn't wait until
// all tuples are emitted.
func (s *TupleIncrementalEmitterSource) EmitTuplesNB(n int) {
	s.m.Lock()
	defer s.m.Unlock()
	s.cnt += n
	s.state.cond.Broadcast()
}

// WaitForEmission waits until all tuples specified by EmitTuples(NB) are emitted.
func (s *TupleIncrementalEmitterSource) WaitForEmission() {
	s.m.Lock()
	defer s.m.Unlock()
	for s.cnt > 0 {
		s.state.cond.Wait()
	}
}

func (s *TupleIncrementalEmitterSource) Stop(ctx *Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.state.state == TSStopped {
		return nil
	}
	s.state.setWithoutLock(TSStopping)
	s.state.waitWithoutLock(TSStopped)
	return nil
}

func (s *TupleIncrementalEmitterSource) Pause(ctx *Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.state.state >= TSStopping {
		return fmt.Errorf("source is already stopped")
	}
	s.state.setWithoutLock(TSPaused)
	return nil
}

func (s *TupleIncrementalEmitterSource) Resume(ctx *Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.state.state >= TSStopping {
		return fmt.Errorf("source is already stopped")
	}
	s.state.setWithoutLock(TSRunning)
	return nil
}

func (s *TupleIncrementalEmitterSource) Status() data.Map {
	return data.Map{
		"test": data.String("test"),
	}
}

// BlockingForwardBox blocks tuples in it and emit them at the timing required by users.
// BlockingForwardBox.c won't be initialized until Init method is called.
// It's safe to call EmitTuples after the topology's state becomes TSRunning.
type BlockingForwardBox struct {
	m   sync.Mutex
	c   *sync.Cond
	cnt int
}

var _ StatefulBox = &BlockingForwardBox{}

func (b *BlockingForwardBox) Init(ctx *Context) error {
	b.c = sync.NewCond(&b.m)
	return nil
}

func (b *BlockingForwardBox) Process(ctx *Context, t *Tuple, w Writer) error {
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

func (b *BlockingForwardBox) Terminate(ctx *Context) error {
	return nil
}

type TupleCollectorSink struct {
	Tuples []*Tuple
	m      sync.Mutex
	c      *sync.Cond
}

func NewTupleCollectorSink() *TupleCollectorSink {
	s := &TupleCollectorSink{}
	s.c = sync.NewCond(&s.m)
	return s
}

func (s *TupleCollectorSink) Write(ctx *Context, t *Tuple) error {
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

type stubForwardBox struct {
	proc func() error
}

func (c *stubForwardBox) Process(ctx *Context, t *Tuple, w Writer) error {
	if c.proc != nil {
		if err := c.proc(); err != nil {
			return err
		}
	}
	return w.Write(ctx, t)
}

// forwardBox sends an input Tuple to the given Writer without
// modification. It can be wrapped with BoxFunc to match the Box
// interface.
func forwardBox(ctx *Context, t *Tuple, w Writer) error {
	w.Write(ctx, t)
	return nil
}

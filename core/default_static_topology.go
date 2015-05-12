package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
)

type defaultStaticTopology struct {
	srcs  map[string]Source
	boxes map[string]Box
	sinks map[string]Sink

	// srcDsts have a set of destinations to which each source write tuples.
	// This is necessary because Source cannot directly close writers and
	// defaultStaticTopology has to take care of it.
	srcDsts map[string]WriteCloser

	// conns have connectors of boxes and sinks. defaultStaticTopology will run
	// a separated goroutine for each connector so that connectors can concurrently
	// send/receive tuples to/from other connectors.
	conns map[string]*staticConnector

	state      TopologyState
	stateMutex *sync.Mutex
	stateCond  *sync.Cond
}

func (t *defaultStaticTopology) Run(ctx *Context) error {
	checkState := func() error { // A closure is used to perform defer
		t.stateMutex.Lock()
		defer t.stateMutex.Unlock()
		switch t.state {
		case TSInitialized:
		case TSStarting:
			// Immediately returning an error could be confusing for callers,
			// so wait until the topology becomes at least the running state.
			if _, err := t.wait(TSRunning); err != nil {
				return err
			}

			// It's natural for Run to return an error when the state isn't
			// TSInitialized even if it's TSStarting so that only a single
			// caller will succeed.
			fallthrough

		default:
			return fmt.Errorf("the static topology has alread started")
		}
		t.state = TSStarting
		return nil
	}
	if err := checkState(); err != nil {
		return err
	}

	// Initialize boxes in advance.
	for _, box := range t.boxes {
		if err := box.Init(ctx); err != nil {
			t.stateMutex.Lock()
			defer t.stateMutex.Unlock()
			t.state = TSStopped
			t.stateCond.Broadcast()
			return err
		}
	}

	// Run goroutines for each connector(boxes and sinks).
	wg := sync.WaitGroup{}
	for _, conn := range t.conns {
		conn := conn
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn.Run(ctx)
		}()
	}

	// Start all sources
	for name, src := range t.srcs {
		name := name
		src := src
		dst := t.srcDsts[name]

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if err := dst.Close(ctx); err != nil {
					// TODO: logging
				}
			}()
			if err := src.GenerateStream(ctx, newTraceWriter(dst, tuple.OUTPUT, name)); err != nil {
				// TODO: logging
			}
		}()
	}
	t.stateMutex.Lock()
	t.state = TSRunning
	t.stateCond.Broadcast()
	t.stateMutex.Unlock()
	wg.Wait()

	t.stateMutex.Lock()
	t.state = TSStopped
	t.stateCond.Broadcast()
	t.stateMutex.Unlock()
	return nil
}

// State returns the current state of the topology. See TopologyState for details.
func (t *defaultStaticTopology) State(ctx *Context) TopologyState {
	t.stateMutex.Lock()
	defer t.stateMutex.Unlock()
	return t.state
}

// Wait waits until the topology has the specified state. It returns the
// current state of the topology. The current state may differ from the
// given state, but it's guaranteed that the current state is a successor of
// the given state when error is nil. For example, when Wait(TSStarting) is
// is called, TSRunning or TSStopped can be returned.
func (t *defaultStaticTopology) Wait(ctx *Context, s TopologyState) (TopologyState, error) {
	t.stateMutex.Lock()
	defer t.stateMutex.Unlock()
	return t.wait(s)
}

// wait is the internal version of Wait. It assumes the caller
// has already acquired the lock of stateMutex.
func (t *defaultStaticTopology) wait(s TopologyState) (TopologyState, error) {
	for t.state < s {
		t.stateCond.Wait()
	}
	return t.state, nil
}

func (t *defaultStaticTopology) Stop(ctx *Context) error {
	stopped, err := func() (bool, error) {
		t.stateMutex.Lock()
		defer t.stateMutex.Unlock()
		for {
			switch t.state {
			case TSInitialized:
				t.state = TSStopped
				t.stateCond.Broadcast()
				return true, nil

			case TSStarting:
				if _, err := t.wait(TSRunning); err != nil {
					return false, err
				}
				// If somebody else has already stopped the topology,
				// the state might be different from TSRunning. So, this process
				// continues to the next iteration.

			case TSRunning:
				t.state = TSStopping
				t.stateCond.Broadcast()
				return false, nil

			case TSStopping:
				// Someone else is trying to stop the topology. This thread
				// just waits until it's stopped.
				_, err := t.wait(TSStopped)
				return true, err

			case TSStopped:
				return true, nil

			default:
				return false, fmt.Errorf("the static topology has an invalid state: %v", t.state)
			}
		}
	}()
	if err != nil {
		return err
	} else if stopped {
		return nil
	}

	for _, src := range t.srcs {
		if err := src.Stop(ctx); err != nil {
			// TODO: logging
			// TODO: appropriate error handling
		}
	}

	// Once all sources are stopped, the stream will eventually stop.
	t.stateMutex.Lock()
	defer t.stateMutex.Unlock()
	_, err = t.wait(TSStopped)
	return err
}

// staticConnector connects sources, boxes, and sinks. It allows boxes and sinks
// to receive tuples from multiple sources and boxes. It also support fanning out
// same tuples to multiple destinations.
type staticConnector struct {
	// dst is a writer which sends tuples to the real destination.
	// dst can be a box, a sink.
	dst WriteCloser

	// inputs is the input channels
	inputs map[string]<-chan *tuple.Tuple // TODO: this should be []*tuple.Tuple for efficiency
}

func newStaticConnector(dst WriteCloser) *staticConnector {
	return &staticConnector{
		dst:    dst,
		inputs: map[string]<-chan *tuple.Tuple{},
	}
}

func (sc *staticConnector) AddInput(name string, c <-chan *tuple.Tuple) {
	sc.inputs[name] = c
}

func (sc *staticConnector) Run(ctx *Context) {
	wg := sync.WaitGroup{}
	for _, in := range sc.inputs {
		in := in

		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range in {
				if err := sc.dst.Write(ctx, t); err != nil {
					// TODO: logging
					// TODO: appropriate error handling
				}
			}
		}()
	}
	wg.Wait()

	if err := sc.dst.Close(ctx); err != nil {
		// TODO: logging
	}
}

type staticSingleChan struct {
	inputName string
	out       chan<- *tuple.Tuple
}

func newStaticSingleChan(name string, ch chan<- *tuple.Tuple) *staticSingleChan {
	return &staticSingleChan{
		inputName: name,
		out:       ch,
	}
}

func (sc *staticSingleChan) Write(ctx *Context, t *tuple.Tuple) error {
	t.InputName = sc.inputName
	sc.out <- t // TODO: writer side load shedding
	return nil
}

func (sc *staticSingleChan) Close(ctx *Context) error {
	close(sc.out)
	return nil
}

type staticDestinations struct {
	dsts []WriteCloser
}

func newStaticDestinations() *staticDestinations {
	return &staticDestinations{}
}

func (mc *staticDestinations) AddDestination(w WriteCloser) {
	mc.dsts = append(mc.dsts, w)
}

func (mc *staticDestinations) Write(ctx *Context, t *tuple.Tuple) error {
	needsCopy := len(mc.dsts) > 1
	for _, o := range mc.dsts {
		s := t
		if needsCopy {
			s = t.Copy()
		}
		if err := o.Write(ctx, s); err != nil {
			// TODO: logging and appropriate error handling
		}
	}
	return nil
}

func (mc *staticDestinations) Close(ctx *Context) error {
	var err error
	for _, o := range mc.dsts {
		if e := o.Close(ctx); e != nil {
			// TODO: logging
			// TODO: appropriate error handling
			err = e
		}
	}
	return err
}

type boxWriterAdapter struct {
	box  Box
	name string
	dst  *traceWriter
}

func newBoxWriterAdapter(b Box, name string, dst WriteCloser) *boxWriterAdapter {
	return &boxWriterAdapter{
		box:  b,
		name: name,
		// An output traces is written just after the box Process writes a tuple.
		dst: newTraceWriter(dst, tuple.OUTPUT, name),
	}
}

func (wa *boxWriterAdapter) Write(ctx *Context, t *tuple.Tuple) error {
	tracing(t, ctx, tuple.INPUT, wa.name)
	return wa.box.Process(ctx, t, wa.dst)
}

func (wa *boxWriterAdapter) Close(ctx *Context) error {
	return wa.dst.w.Close(ctx)
}

type traceWriter struct {
	w     WriteCloser
	inout tuple.EventType
	msg   string
}

func newTraceWriter(w WriteCloser, inout tuple.EventType, msg string) *traceWriter {
	return &traceWriter{
		w:     w,
		inout: inout,
		msg:   msg,
	}
}

func (tw *traceWriter) Write(ctx *Context, t *tuple.Tuple) error {
	tracing(t, ctx, tw.inout, tw.msg)
	return tw.w.Write(ctx, t)
}

func (tw *traceWriter) Close(ctx *Context) error {
	return tw.w.Close(ctx)
}

package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
	"sync"
)

type defaultStaticTopology struct {
	srcs  map[string]Source
	boxes map[string]Box
	sinks map[string]Sink

	// srcDsts have a set of destinations to which each source write tuples.
	// This is necessary because Source cannot directly close writers and
	// defaultStaticTopology has to take care of it.
	//
	// srcDsts will not be thread-safe once the topology started running.
	srcDsts      map[string]WriteCloser
	srcDstsMutex sync.Mutex

	// conns have connectors of boxes and sinks. defaultStaticTopology will run
	// a separate goroutine for each input of a connector so that connectors can concurrently
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
			return fmt.Errorf("the static topology has already started")
		}
		t.state = TSStarting
		return nil
	}
	if err := checkState(); err != nil {
		return err
	}

	// Don't move this defer to the head of the method otherwise calling
	// Run method when the state is TSRunning will set TSStopped without
	// actually stopping the topology.
	defer t.setState(TSStopped)

	// Initialize boxes in advance.
	for _, box := range t.boxes {
		if err := box.Init(ctx); err != nil {
			return err
		}
	}
	return t.run(ctx)
}

func (t *defaultStaticTopology) setState(s TopologyState) {
	t.stateMutex.Lock()
	defer t.stateMutex.Unlock()
	t.state = s
	t.stateCond.Broadcast()
}

// run spawns goroutines for sources, boxes, and sinks.
// The caller of it must set TSStop after it returns.
func (t *defaultStaticTopology) run(ctx *Context) error {
	// Run goroutines for each connector(boxes and sinks).
	var wg sync.WaitGroup
	for name, conn := range t.conns {
		wg.Add(1)
		go func(name string, conn *staticConnector) {
			defer wg.Done()
			conn.Run(name, ctx)
		}(name, conn)
	}

	// Create closures for goroutines here because srcDsts will not
	// be thread-safe once those goroutines start.
	fs := make([]func(), 0, len(t.srcs))
	for name, src := range t.srcs {
		name := name
		src := src
		dst := t.srcDsts[name]
		f := func() {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ctx.Logger.Log(Error, "%v paniced: %v", name, err)
				}

				// Because dst could be closed by defaultStaticTopology.Stop when
				// Source.Stop failed in that method, closing dst must be done
				// via closeDestination method.
				if err := t.closeDestination(ctx, name); err != nil {
					ctx.Logger.Log(Error, "%v cannot close the destination: %v", name, err)
				}
			}()
			if err := src.GenerateStream(ctx, newTraceWriter(dst, tuple.Output, name)); err != nil {
				ctx.Logger.Log(Error, "%v cannot generate tuples: %v", name, err)
			}
		}
		fs = append(fs, f)
	}
	for _, f := range fs {
		wg.Add(1)
		go f()
	}

	t.setState(TSRunning)
	wg.Wait()
	return nil
}

func (t *defaultStaticTopology) closeDestination(ctx *Context, src string) error {
	dst := func() WriteCloser {
		t.srcDstsMutex.Lock()
		defer t.srcDstsMutex.Unlock()

		// Since WriteCloser.Close doesn't have to be idempotent, it should
		// only be called exactly once.
		dst, ok := t.srcDsts[src]
		if !ok {
			return nil
		}
		delete(t.srcDsts, src)
		return dst
	}()

	// srcDstsMutex is unlocked here.
	if dst == nil {
		return nil
	}
	return dst.Close(ctx)
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

	var stopFailures []string
	for name, src := range t.srcs {
		if err := src.Stop(ctx); err != nil {
			ctx.Logger.Log(Error, "Cannot stop source %v: ", name, err)
			stopFailures = append(stopFailures, name)
			if err := t.closeDestination(ctx, name); err != nil {
				ctx.Logger.Log(Error, "Cannot close the source %v's destination: %v", name, err)
			}
		}
	}

	// TODO: There might be some WriteClosers which still haven't been closed
	// and some connectors attached to them are still running. There might
	// have to be some way to force shutdown connectors.

	// Once all sources are stopped, the stream will eventually stop.
	t.stateMutex.Lock()
	defer t.stateMutex.Unlock()
	_, err = t.wait(TSStopped)

	if err == nil && len(stopFailures) > 0 {
		return fmt.Errorf("%v sources couldn't be stopped but the topology has stopped: failed sources = %v",
			len(stopFailures), strings.Join(stopFailures, ", "))
	}
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

func (sc *staticConnector) Run(name string, ctx *Context) {
	var (
		wg           sync.WaitGroup
		panicLogging sync.Once
	)
	for _, in := range sc.inputs {
		in := in

		wg.Add(1)
		// TODO: These goroutines should be assembled to one goroutine with reflect.Select.
		go func() {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					// Once a Box panics, it'll always panic after that. So,
					// the log should only be written once.
					panicLogging.Do(func() {
						ctx.Logger.Log(Error, "%v paniced: %v", name, err)
					})
				}
			}()

			for t := range in {
				if err := sc.dst.Write(ctx, t); err != nil {
					ctx.Logger.Log(Error, "%v cannot write a tuple: %v", name, err)
					// All regular errors are considered resumable.
				}
			}
		}()
	}
	wg.Wait()

	if err := sc.dst.Close(ctx); err != nil {
		ctx.Logger.Log(Error, "%v cannot close its output channel: %v", name, err)
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
	names []string
	dsts  []WriteCloser
}

func newStaticDestinations() *staticDestinations {
	return &staticDestinations{}
}

func (mc *staticDestinations) AddDestination(name string, w WriteCloser) {
	mc.names = append(mc.names, name)
	mc.dsts = append(mc.dsts, w)
}

func (mc *staticDestinations) Write(ctx *Context, t *tuple.Tuple) error {
	e := &bulkErrors{}
	needsCopy := len(mc.dsts) > 1
	for i, d := range mc.dsts {
		s := t
		if needsCopy {
			s = t.Copy()
		}
		if err := d.Write(ctx, s); err != nil {
			// TODO: this error message could be a little redundant.
			e.append(fmt.Errorf("a tuple cannot be written to %v: %v", mc.names[i], err))
		}
	}
	return e.returnError()
}

func (mc *staticDestinations) Close(ctx *Context) error {
	e := &bulkErrors{}
	for i, d := range mc.dsts {
		if err := d.Close(ctx); err != nil {
			e.append(fmt.Errorf("output channel to %v cannot be closed: %v", mc.names[i], err))
		}
	}
	return e.returnError()
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
		dst: newTraceWriter(dst, tuple.Output, name),
	}
}

func (wa *boxWriterAdapter) Write(ctx *Context, t *tuple.Tuple) error {
	tracing(t, ctx, tuple.Input, wa.name)
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

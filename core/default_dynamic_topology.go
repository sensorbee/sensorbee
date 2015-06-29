package core

import (
	"fmt"
	"sync"
)

type defaultDynamicTopology struct {
	ctx  *Context
	name string

	// nodeMutex protects sources, boxes, and sinks from being modified
	// concurrently. DO NOT acquire this lock after locking stateMutex
	// (opposite is fine).
	nodeMutex sync.RWMutex
	sources   map[string]*defaultSourceNode
	boxes     map[string]*defaultBoxNode
	sinks     map[string]*defaultSinkNode

	state      *topologyStateHolder
	stateMutex sync.Mutex

	// TODO: support lazy invocation of GenerateStream (call it when the first
	// destination is added or a Sink is indirectly connected). Maybe graph
	// management is required.
}

// NewDefaultDynamicTopology creates a dynamic topology having a simple graph
// structure.
func NewDefaultDynamicTopology(ctx *Context, name string) DynamicTopology {
	// TODO: validate name

	t := &defaultDynamicTopology{
		ctx:  ctx,
		name: name,

		sources: map[string]*defaultSourceNode{},
		boxes:   map[string]*defaultBoxNode{},
		sinks:   map[string]*defaultSinkNode{},
	}
	t.state = newTopologyStateHolder(&t.stateMutex)
	t.state.state = TSRunning // A dynamic topology is running by default.
	return t
}

func (t *defaultDynamicTopology) Context() *Context {
	return t.ctx
}

func (t *defaultDynamicTopology) AddSource(name string, s Source, config *DynamicSourceConfig) (DynamicSourceNode, error) {
	// TODO: validate the name

	if config == nil {
		config = &DynamicSourceConfig{}
	}

	// This method assumes adding a Source having a duplicated name is rare.
	// Under this assumption, acquiring wlock without checking the existence
	// of the name with rlock doesn't degrade the performance.
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()

	// t.state is set to TSStopped while t.nodeMutex is locked. Therefore,
	// a source can safely be added when the state is TSRunning or TSPaused.
	//
	// The lock above doesn't protect t.state from being set to TSStopping.
	// So, some goroutine can go through this if block and add a source to
	// the topology when the state is TSStopping. However, adding a source
	// while the state is TSStopping is safe although the source will get
	// removed just after it's added to the topology.
	if t.state.Get() >= TSStopping {
		return nil, fmt.Errorf("the topology is already stopped")
	}

	ds := &defaultSourceNode{
		defaultNode:     newdefaultNode(t, name),
		source:          s,
		dsts:            newDataDestinations(name),
		pausedOnStartup: config.PausedOnStartup,
	}
	if err := t.checkNodeNameDuplication(name); err != nil {
		if err := ds.Stop(); err != nil { // The same variable name is intentionally used.
			t.ctx.Logger.Log(Error, "Cannot stop the source '%v': %v", name, err)
		}
		return nil, err
	}
	t.sources[name] = ds

	go func() {
		// TODO: Support lazy invocation
		if err := ds.run(); err != nil {
			t.ctx.Logger.Log(Error, "Cannot generate a stream from the source '%v': %v", name, err)
		}
	}()

	if config.PausedOnStartup {
		ds.state.Wait(TSPaused)
	} else {
		ds.state.Wait(TSRunning)
	}
	return ds, nil
}

// TODO: Add method to validate a node name

// checkNodeNameDuplication checks if the given name is unique in the topology.
// This method doesn't acquire the lock and it's the caller's responsibility
// to do it before calling this method.
func (t *defaultDynamicTopology) checkNodeNameDuplication(name string) error {
	if _, ok := t.sources[name]; ok {
		return fmt.Errorf("the name is already used by a source: %v", name)
	}
	if _, ok := t.boxes[name]; ok {
		return fmt.Errorf("the name is already used by a box: %v", name)
	}
	if _, ok := t.sinks[name]; ok {
		return fmt.Errorf("the name is already used by a sink: %v", name)
	}
	return nil
}

func (t *defaultDynamicTopology) AddBox(name string, b Box, config *DynamicBoxConfig) (DynamicBoxNode, error) {
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if t.state.Get() >= TSStopping {
		return nil, fmt.Errorf("the topology is already stopped")
	}

	if err := t.checkNodeNameDuplication(name); err != nil {
		return nil, err
	}

	if sb, ok := b.(StatefulBox); ok {
		if err := sb.Init(t.ctx); err != nil {
			return nil, err
		}
	}

	db := &defaultBoxNode{
		defaultNode: newdefaultNode(t, name),
		srcs:        newDataSources(name),
		box:         b,
		dsts:        newDataDestinations(name),
	}
	t.boxes[name] = db

	go func() {
		if err := db.run(); err != nil {
			t.ctx.Logger.Log(Error, "Box '%v' in topology '%v' failed: %v", db.name, t.name, err)
		}
	}()
	db.state.Wait(TSRunning)
	return db, nil
}

func (t *defaultDynamicTopology) AddSink(name string, s Sink, config *DynamicSinkConfig) (DynamicSinkNode, error) {
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if t.state.Get() >= TSStopping {
		return nil, fmt.Errorf("the topology is already stopped")
	}

	if err := t.checkNodeNameDuplication(name); err != nil {
		return nil, err
	}

	ds := &defaultSinkNode{
		defaultNode: newdefaultNode(t, name),
		srcs:        newDataSources(name),
		sink:        s,
	}
	t.sinks[name] = ds

	go func() {
		if err := ds.run(); err != nil {
			t.ctx.Logger.Log(Error, "Sink '%v' in topology '%v' failed: %v", ds.name, t.name, err)
		}
	}()
	ds.state.Wait(TSRunning)
	return ds, nil
}

func (t *defaultDynamicTopology) Stop() error {
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if stopped, err := t.state.checkAndPrepareForStoppingWithoutLock(false); err != nil {
		return fmt.Errorf("the dynamic topology has an invalid state: %v", t.state.Get())
	} else if stopped {
		return nil
	}

	var lastErr error
	for name, src := range t.sources {
		// TODO: this could be run concurrently
		func() {
			defer func() {
				if e := recover(); e != nil {
					t.ctx.Logger.Log(Error, "Cannot stop source '%v' due to panic: %v", name, e)
					src.dsts.Close(t.ctx)

					if err, ok := e.(error); ok {
						lastErr = err
					} else {
						lastErr = fmt.Errorf("source '%v' panicked while being stopped: %v", name, e)
					}
				}
			}()

			if err := src.Stop(); err != nil {
				src.dsts.Close(t.ctx)
				lastErr = err
				t.ctx.Logger.Log(Error, "Cannot stop source '%v': %v", name, err)
			}
		}()
	}

	var wg sync.WaitGroup
	for _, b := range t.boxes {
		b.srcs.stopOnDisconnect()
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.state.Wait(TSStopped)
		}()
	}

	for _, s := range t.sinks {
		s.srcs.stopOnDisconnect()
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.state.Wait(TSStopped)
		}()
	}
	wg.Wait()

	t.sources = nil
	t.boxes = nil
	t.sinks = nil
	t.state.Set(TSStopped)
	return nil
}

func (t *defaultDynamicTopology) State() TopologyStateHolder {
	return t.state
}

func (t *defaultDynamicTopology) Remove(name string) error {
	n := func() DynamicNode {
		t.nodeMutex.Lock()
		defer t.nodeMutex.Unlock()

		n, err := t.nodeWithoutLock(name)
		if err != nil {
			return nil // not found
		}
		switch n.Type() {
		case NTSource:
			delete(t.sources, name)
		case NTBox:
			delete(t.boxes, name)
		case NTSink:
			delete(t.sinks, name)
		}
		return n
	}()
	if n == nil {
		return nil // already removed or doesn't exist
	}

	return func() (retErr error) {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					retErr = err
				} else {
					retErr = FatalError(fmt.Errorf("%v '%v' failed to stop with panic: %v", n.Type(), name, e))
				}
			}

			if retErr != nil && n.Type() == NTSource {
				s := n.(*defaultSourceNode)
				s.dsts.Close(t.ctx)
			}
		}()
		if err := n.Stop(); err != nil {
			retErr = err
		}
		return
	}()
}

// TODO: Add method to clean up (possibly indirectly) stopped nodes

func (t *defaultDynamicTopology) Node(name string) (DynamicNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	return t.nodeWithoutLock(name)
}

func (t *defaultDynamicTopology) nodeWithoutLock(name string) (DynamicNode, error) {
	if s, ok := t.sources[name]; ok {
		return s, nil
	}
	if b, ok := t.boxes[name]; ok {
		return b, nil
	}
	if s, ok := t.sinks[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("node '%v' was not found", name)
}

func (t *defaultDynamicTopology) Nodes() map[string]DynamicNode {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]DynamicNode, len(t.sources)+len(t.boxes)+len(t.sinks))
	for name, s := range t.sources {
		m[name] = s
	}
	for name, b := range t.boxes {
		m[name] = b
	}
	for name, s := range t.sinks {
		m[name] = s
	}
	return m
}

func (t *defaultDynamicTopology) Source(name string) (DynamicSourceNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	if s, ok := t.sources[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("source '%v' was not found", name)
}

func (t *defaultDynamicTopology) Sources() map[string]DynamicSourceNode {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]DynamicSourceNode, len(t.sources))
	for name, s := range t.sources {
		m[name] = s
	}
	return m
}

func (t *defaultDynamicTopology) Box(name string) (DynamicBoxNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	if b, ok := t.boxes[name]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("box '%v' was not found", name)
}

func (t *defaultDynamicTopology) Boxes() map[string]DynamicBoxNode {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]DynamicBoxNode, len(t.boxes))
	for name, b := range t.boxes {
		m[name] = b
	}
	return m
}

func (t *defaultDynamicTopology) Sink(name string) (DynamicSinkNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	if s, ok := t.sinks[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("sink '%v' was not found", name)
}

func (t *defaultDynamicTopology) Sinks() map[string]DynamicSinkNode {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]DynamicSinkNode, len(t.sinks))
	for name, s := range t.sinks {
		m[name] = s
	}
	return m
}

type dynamicDataSource interface {
	Name() string
	destinations() *dataDestinations
}

func (t *defaultDynamicTopology) dataSource(nodeName string) (dynamicDataSource, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	if s, ok := t.sources[nodeName]; ok {
		return s, nil
	}
	if b, ok := t.boxes[nodeName]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("data source node %v was not found", nodeName)
}

type defaultNode struct {
	topology   *defaultDynamicTopology
	name       string
	state      *topologyStateHolder
	stateMutex sync.Mutex
}

func newdefaultNode(t *defaultDynamicTopology, name string) *defaultNode {
	dn := &defaultNode{
		topology: t,
		name:     name,
	}
	dn.state = newTopologyStateHolder(&dn.stateMutex)
	return dn
}

func (dn *defaultNode) Name() string {
	return dn.name
}

func (dn *defaultNode) State() TopologyStateHolder {
	return dn.state
}

func (dn *defaultNode) checkAndPrepareForRunning(nodeType string) error {
	dn.stateMutex.Lock()
	defer dn.stateMutex.Unlock()
	return dn.checkAndPrepareForRunningWithoutLock(nodeType)
}

func (dn *defaultNode) checkAndPrepareForRunningWithoutLock(nodeType string) error {
	if st, err := dn.state.checkAndPrepareForRunningWithoutLock(); err != nil {
		switch st {
		case TSRunning, TSPaused:
			return fmt.Errorf("%v '%v' is already running", nodeType, dn.name)
		case TSStopped:
			return fmt.Errorf("%v '%v' is already stopped", nodeType, dn.name)
		default:
			return fmt.Errorf("%v '%v' has an invalid state: %v", nodeType, dn.name, st)
		}
	}
	return nil
}

// checkAndPrepareStopState check the current state of the node and returns if
// the node can be stopped or is already stopped.
func (dn *defaultNode) checkAndPrepareForStopping(nodeType string) (stopped bool, err error) {
	dn.stateMutex.Lock()
	defer dn.stateMutex.Unlock()
	return dn.checkAndPrepareForStoppingWithoutLock(nodeType)
}

func (dn *defaultNode) checkAndPrepareForStoppingWithoutLock(nodeType string) (stopped bool, err error) {
	return dn.state.checkAndPrepareForStoppingWithoutLock(false)
}

package core

import (
	"fmt"
	"strings"
	"sync"
)

type defaultTopology struct {
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

// NewDefaultTopology creates a topology having a simple graph
// structure.
func NewDefaultTopology(ctx *Context, name string) Topology {
	ctx.topologyName = name
	t := &defaultTopology{
		ctx:  ctx,
		name: name,

		sources: map[string]*defaultSourceNode{},
		boxes:   map[string]*defaultBoxNode{},
		sinks:   map[string]*defaultSinkNode{},
	}
	t.state = newTopologyStateHolder(&t.stateMutex)
	t.state.state = TSRunning // A topology is running by default.
	return t
}

func (t *defaultTopology) Name() string {
	return t.name
}

func (t *defaultTopology) Context() *Context {
	return t.ctx
}

func (t *defaultTopology) AddSource(name string, s Source, config *SourceConfig) (SourceNode, error) {
	if err := ValidateNodeName(name); err != nil {
		return nil, err
	}

	if config == nil {
		config = &SourceConfig{}
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
		defaultNode:     newDefaultNode(t, name, config.Meta),
		source:          s,
		dsts:            newDataDestinations(NTSource, name),
		pausedOnStartup: config.PausedOnStartup,
	}
	ds.config = &SourceConfig{}
	*ds.config = *config
	ds.dsts.callback = ds.dstCallback
	if err := t.checkNodeNameDuplication(name); err != nil {
		// Because the source isn't started yet, it doesn't return an error.
		ds.Stop()
		return nil, err
	}
	if _, ok := s.(*droppedTupleCollectorSource); ok {
		// Tuples dropped by droppedTupleCollectorSource and nodes connected
		// to it must not be reported again.
		ds.dsts.disableDroppedTupleReporting()
	}
	t.sources[strings.ToLower(name)] = ds

	go func() {
		// TODO: Support lazy invocation
		if err := ds.run(); err != nil {
			t.ctx.ErrLog(err).WithFields(nodeLogFields(NTSource, name)).
				Error("Cannot generate a stream from the source")
		}
		if ds.config.RemoveOnStop {
			if err := t.Remove(name); err != nil {
				t.ctx.ErrLog(err).WithFields(nodeLogFields(NTSource, name)).
					Error("Cannot remove the source from topology")
			}
		}
	}()

	if config.PausedOnStartup {
		ds.state.Wait(TSPaused)
	} else {
		ds.state.Wait(TSRunning)
	}
	return ds, nil
}

// checkNodeNameDuplication checks if the given name is unique in the topology.
// This method doesn't acquire the lock and it's the caller's responsibility
// to do it before calling this method.
func (t *defaultTopology) checkNodeNameDuplication(name string) error {
	lowerName := strings.ToLower(name)
	if _, ok := t.sources[lowerName]; ok {
		return fmt.Errorf("the name is already used by a source: %v", name)
	}
	if _, ok := t.boxes[lowerName]; ok {
		return fmt.Errorf("the name is already used by a box: %v", name)
	}
	if _, ok := t.sinks[lowerName]; ok {
		return fmt.Errorf("the name is already used by a sink: %v", name)
	}
	return nil
}

func (t *defaultTopology) AddBox(name string, b Box, config *BoxConfig) (BoxNode, error) {
	if err := ValidateNodeName(name); err != nil {
		return nil, err
	}

	if config == nil {
		config = &BoxConfig{}
	}

	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if t.state.Get() >= TSStopping {
		return nil, fmt.Errorf("the topology is already stopped")
	}

	if err := t.checkNodeNameDuplication(name); err != nil {
		return nil, err
	}

	if sb, ok := b.(StatefulBox); ok {
		err := func() (err error) {
			defer func() {
				if e := recover(); e != nil {
					if er, ok := e.(error); ok {
						err = er
					} else {
						err = fmt.Errorf("the box cannot be initialized due to panic: %v", e)
					}
				}
			}()
			return sb.Init(t.ctx)
		}()
		if err != nil {
			return nil, err
		}
	}

	db := &defaultBoxNode{
		defaultNode: newDefaultNode(t, name, config.Meta),
		srcs:        newDataSources(NTBox, name),
		box:         b,
		dsts:        newDataDestinations(NTBox, name),
	}
	db.config = &BoxConfig{}
	*db.config = *config
	db.dsts.callback = db.dstCallback
	t.boxes[strings.ToLower(name)] = db

	go func() {
		if err := db.run(); err != nil {
			t.ctx.ErrLog(err).WithFields(nodeLogFields(NTBox, db.name)).
				Error("The box failed")
		}
		if db.config.RemoveOnStop {
			if err := t.Remove(name); err != nil {
				t.ctx.ErrLog(err).WithFields(nodeLogFields(NTBox, db.name)).
					Error("Cannot remove the box from topology")
			}
		}
	}()
	db.state.Wait(TSRunning)
	return db, nil
}

func (t *defaultTopology) AddSink(name string, s Sink, config *SinkConfig) (SinkNode, error) {
	// Sink must be closed when AddSink fails before creating a sink node
	closeSinkFlag := false
	defer func() {
		if !closeSinkFlag {
			return
		}
		defer func() {
			if e := recover(); e != nil {
				t.ctx.Log().WithFields(nodeLogFields(NTSink, name)).
					Errorf("Cannot close the sink which hasn't been added to the topology: %v", e)
			}
		}()
		if err := s.Close(t.ctx); err != nil {
			t.ctx.ErrLog(err).WithFields(nodeLogFields(NTSink, name)).
				Error("Cannot close the sink which hasn't been added to the topology")
		}
	}()

	if err := ValidateNodeName(name); err != nil {
		closeSinkFlag = true
		return nil, err
	}

	if config == nil {
		config = &SinkConfig{}
	}

	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if t.state.Get() >= TSStopping {
		closeSinkFlag = true
		return nil, fmt.Errorf("the topology is already stopped")
	}

	if err := t.checkNodeNameDuplication(name); err != nil {
		closeSinkFlag = true
		return nil, err
	}

	ds := &defaultSinkNode{
		defaultNode: newDefaultNode(t, name, config.Meta),
		srcs:        newDataSources(NTSink, name),
		sink:        s,
	}
	ds.config = &SinkConfig{}
	*ds.config = *config
	t.sinks[strings.ToLower(name)] = ds

	go func() {
		if err := ds.run(); err != nil {
			t.ctx.ErrLog(err).WithFields(nodeLogFields(NTSink, ds.name)).
				Error("The sink failed")
		}
		if ds.config.RemoveOnStop {
			if err := t.Remove(name); err != nil {
				t.ctx.ErrLog(err).WithFields(nodeLogFields(NTSink, ds.name)).
					Error("Cannot remove the sink from topology")
			}
		}
	}()
	ds.state.Wait(TSRunning)
	return ds, nil
}

func (t *defaultTopology) Stop() error {
	t.nodeMutex.Lock()
	defer t.nodeMutex.Unlock()
	if stopped, err := t.state.checkAndPrepareForStoppingWithoutLock(false); err != nil {
		return fmt.Errorf("the topology has an invalid state: %v", t.state.Get())
	} else if stopped {
		return nil
	}

	var lastErr error
	for name, src := range t.sources {
		// TODO: this could be run concurrently
		if err := src.Stop(); err != nil { // Stop doesn't panic
			lastErr = err
			src.dsts.Close(t.ctx)
			t.ctx.ErrLog(err).WithFields(nodeLogFields(NTSource, name)).
				Error("Cannot stop the source")
		}
	}

	var wg sync.WaitGroup
	for _, b := range t.boxes {
		b := b

		b.StopOnDisconnect(Inbound | Outbound)
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.state.Wait(TSStopped)
		}()
	}

	for _, s := range t.sinks {
		s := s

		s.StopOnDisconnect()
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
	return lastErr
}

func (t *defaultTopology) State() TopologyStateHolder {
	return t.state
}

func (t *defaultTopology) Remove(name string) error {
	lowerName := strings.ToLower(name)
	n := func() Node {
		t.nodeMutex.Lock()
		defer t.nodeMutex.Unlock()

		n, err := t.nodeWithoutLock(name)
		if err != nil {
			return nil // not found
		}
		switch n.Type() {
		case NTSource:
			delete(t.sources, lowerName)
		case NTBox:
			delete(t.boxes, lowerName)
		case NTSink:
			delete(t.sinks, lowerName)
		}
		return n
	}()
	if n == nil {
		return nil // already removed or doesn't exist
	}

	if err := n.Stop(); err != nil { // stop never panics
		if n.Type() == NTSource {
			s := n.(*defaultSourceNode)
			s.dsts.Close(t.ctx)
		}
		return err
	}
	return nil
}

// TODO: Add method to clean up (possibly indirectly) stopped nodes

func (t *defaultTopology) Node(name string) (Node, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	return t.nodeWithoutLock(name)
}

func (t *defaultTopology) nodeWithoutLock(name string) (Node, error) {
	lowerName := strings.ToLower(name)
	if s, ok := t.sources[lowerName]; ok {
		return s, nil
	}
	if b, ok := t.boxes[lowerName]; ok {
		return b, nil
	}
	if s, ok := t.sinks[lowerName]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("node '%v' was not found", name)
}

func (t *defaultTopology) Nodes() map[string]Node {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]Node, len(t.sources)+len(t.boxes)+len(t.sinks))
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

func (t *defaultTopology) Source(name string) (SourceNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	if s, ok := t.sources[strings.ToLower(name)]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("source '%v' was not found", name)
}

func (t *defaultTopology) Sources() map[string]SourceNode {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]SourceNode, len(t.sources))
	for name, s := range t.sources {
		m[name] = s
	}
	return m
}

func (t *defaultTopology) Box(name string) (BoxNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	if b, ok := t.boxes[strings.ToLower(name)]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("box '%v' was not found", name)
}

func (t *defaultTopology) Boxes() map[string]BoxNode {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]BoxNode, len(t.boxes))
	for name, b := range t.boxes {
		m[name] = b
	}
	return m
}

func (t *defaultTopology) Sink(name string) (SinkNode, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()
	if s, ok := t.sinks[strings.ToLower(name)]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("sink '%v' was not found", name)
}

func (t *defaultTopology) Sinks() map[string]SinkNode {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	m := make(map[string]SinkNode, len(t.sinks))
	for name, s := range t.sinks {
		m[name] = s
	}
	return m
}

type dataSource interface {
	Name() string
	destinations() *dataDestinations
}

func (t *defaultTopology) dataSource(nodeName string) (dataSource, error) {
	t.nodeMutex.RLock()
	defer t.nodeMutex.RUnlock()

	lowerNodeName := strings.ToLower(nodeName)
	if s, ok := t.sources[lowerNodeName]; ok {
		return s, nil
	}
	if b, ok := t.boxes[lowerNodeName]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("data source node %v was not found", nodeName)
}

type defaultNode struct {
	topology   *defaultTopology
	name       string
	state      *topologyStateHolder
	stateMutex sync.Mutex

	meta interface{}
}

func newDefaultNode(t *defaultTopology, name string, meta interface{}) *defaultNode {
	if meta == nil {
		meta = map[string]interface{}{}
	}
	dn := &defaultNode{
		topology: t,
		name:     name,
		meta:     meta,
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

func (dn *defaultNode) Meta() interface{} {
	return dn.meta
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

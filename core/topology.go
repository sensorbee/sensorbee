package core

import (
	"fmt"
	"sync"
)

// TopologyState represents a status of a topology or a node.
type TopologyState int

const (
	// TSInitialized means that a topology or a node is just initialized and
	// ready to be run.
	TSInitialized TopologyState = iota

	// TSStarting means a topology or a node is now booting itself and will run
	// shortly.
	TSStarting

	// TSRunning means a topology or a node is currently running and emitting
	// tuples to sinks.
	TSRunning

	// TSPaused means a topology or a node is temporarily stopping to emit
	// tuples and can be resumed later.
	TSPaused

	// TSStopping means a topology or a node is stopping all sources and closing
	// channels between sources, boxes, and sinks.
	TSStopping

	// TSStopped means a topology or a node is stopped. A stopped topology
	// doesn't have to be able to run again.
	TSStopped
)

func (s TopologyState) String() string {
	switch s {
	case TSInitialized:
		return "initialized"
	case TSStarting:
		return "starting"
	case TSRunning:
		return "running"
	case TSPaused:
		return "paused"
	case TSStopping:
		return "stopping"
	case TSStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// TopologyStateHolder is a struct safely referring a state of a topology or a
// node. It only provides read-only methods.
type TopologyStateHolder interface {
	// Get returns the current state of a topology or a node.
	Get() TopologyState

	// Wait waits until the topology or the node has the specified state. It
	// returns the current state. The current state may differ from the given
	// state, but it's guaranteed that the current state is a successor of
	// the given state. For example, when Wait(TSStarting) is called, TSRunning
	// or TSStopped can be returned.
	Wait(s TopologyState) TopologyState
}

type topologyStateHolder struct {
	state TopologyState
	cond  *sync.Cond
}

func newTopologyStateHolder(m sync.Locker) *topologyStateHolder {
	if m == nil {
		m = &sync.Mutex{}
	}
	return &topologyStateHolder{
		cond: sync.NewCond(m),
	}
}

func (h *topologyStateHolder) Get() TopologyState {
	h.cond.L.Lock()
	defer h.cond.L.Unlock()
	return h.getWithoutLock()
}

func (h *topologyStateHolder) getWithoutLock() TopologyState {
	return h.state
}

// Set sets a new state.
func (h *topologyStateHolder) Set(s TopologyState) error {
	h.cond.L.Lock()
	defer h.cond.L.Unlock()
	return h.setWithoutLock(s)
}

func (h *topologyStateHolder) setWithoutLock(s TopologyState) error {
	if h.state > s {
		if h.state == TSPaused && s == TSRunning {
			// TSPaused can exceptionally be reset to TSRunning
		} else {
			return fmt.Errorf("state cannot be changed from %v to %v", h.state, s)
		}
	}
	h.state = s
	h.cond.Broadcast()
	return nil
}

func (h *topologyStateHolder) Wait(s TopologyState) TopologyState {
	h.cond.L.Lock()
	defer h.cond.L.Unlock()
	return h.waitWithoutLock(s)
}

func (h *topologyStateHolder) waitWithoutLock(s TopologyState) TopologyState {
	for {
		if h.state >= s {
			if h.state == TSPaused && s == TSRunning {
				// Wait until the state becomes TSRunning
			} else {
				break
			}
		}
		h.cond.Wait()
	}
	return h.state
}

// BoxInputConfig has parameters to customize input behavior of a Box on each
// input pipe.
type BoxInputConfig struct {
	// InputName is a custom name attached to incoming tuples. When it is empty,
	// "*" will be used.
	InputName string

	// Capacity is the maximum capacity (length) of input pipe. When this
	// parameter is 0, the default value is used. This parameter is only used
	// as a hint and doesn't guarantee that the pipe can actually have the
	// specified number of tuples.
	Capacity int
}

func (c *BoxInputConfig) inputName() string {
	if c.InputName == "" {
		return "*"
	}
	return c.InputName
}

func (c *BoxInputConfig) capacity() int {
	if c.Capacity == 0 {
		return 1024
	}
	return c.Capacity
}

// SinkInputConfig has parameters to customize input behavior of a Sink on
// each input pipe.
type SinkInputConfig struct {
	// Capacity is the maximum capacity (length) of input pipe. When this
	// parameter is 0, the default value is used. This parameter is only used
	// as a hint and doesn't guarantee that the pipe can actually have the
	// specified number of tuples.
	Capacity int
}

func (c *SinkInputConfig) capacity() int {
	if c.Capacity == 0 {
		return 1024
	}
	return c.Capacity
}

var (
	defaultBoxInputConfig  = &BoxInputConfig{}
	defaultSinkInputConfig = &SinkInputConfig{}
)

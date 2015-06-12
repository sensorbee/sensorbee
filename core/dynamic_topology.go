package core

// DynamicTopology is a topology which can add Sources, Boxes, and Sinks
// dynamically. Boxes and Sinks can also add inputs dynamically from running
// Sources or Boxes.
type DynamicTopology interface {
	// Context returns the Context tied to the topology. It isn't always safe
	// to modify fields of the Context. However, some fields of it can handle
	// concurrent access properly. See the document of Context for details.
	Context() *Context

	// AddSource adds a Source to the topology. It will asynchronously call
	// source's GenerateStream and returns after the source becomes ready.
	// GenerateStream could be called lazily to avoid unnecessary computation.
	// GenerateStream and Stop might also be called when this method returns
	// an error. The caller must not call GenerateStream or Stop of the source.
	AddSource(name string, s Source, config *DynamicSourceConfig) (DynamicSourceNode, error)

	// AddBox adds a Box to the topology. It returns DynamicBoxNode and the
	// caller can configure inputs or other settings of the Box node through it.
	//
	// Don't add the same instance of a Box more than once even if they have
	// different names. However, if a Box has a reference counter and
	// its initialization and termination are done exactly once at proper
	// timing, it can be added multiple times when a builder supports duplicated
	// registration of the same instance of a Box.
	AddBox(name string, b Box, config *DynamicBoxConfig) (DynamicBoxNode, error)

	// AddSink adds a Sink to the topology. It returns DynamicSinkNode and the
	// caller can configure inputs or other settings of the Sink node through it.
	AddSink(name string, s Sink, config *DynamicSinkConfig) (DynamicSinkNode, error)

	// Remove removes a node from the topology. It doesn't return an error when
	// the topology couldn't find a node having the name. The removed node is
	// stopped by the topology and Remove methods blocks until the node actually
	// stops.
	Remove(name string) error

	// Stop stops the topology. If the topology doesn't have a cycle, it stops
	// after all tuples generated from Sources at the time of the invocation
	// are written into Sinks. Stop method returns after processing all the
	// tuples.
	//
	// BUG: Currently Stop method doesn't work if the topology has a cycle.
	Stop() error

	// State returns the current state of the topology. The topology's state
	// isn't relevant to those dynamic nodes have.
	State() TopologyStateHolder

	// TODO: low priority: Pause, Resume

	// Node returns a node registered to the topology. It returns an error
	// when the topology doesn't have the node.
	Node(name string) (DynamicNode, error)

	// Nodes returns all nodes registered to the topology. The map returned
	// from this method can safely be modified.
	Nodes() map[string]DynamicNode

	// Source returns a source registered to the topology. It returns an error
	// when the topology doesn't have the source.
	Source(name string) (DynamicSourceNode, error)

	// Sources returns all sources registered to the topology. The map returned
	// from this method can safely be modified.
	Sources() map[string]DynamicSourceNode

	// Box returns a box registered to the topology. It returns an error when
	// the topology doesn't have the box.
	Box(name string) (DynamicBoxNode, error)

	// Boxes returns all boxes registered to the topology. The map returned
	// from this method can safely be modified.
	Boxes() map[string]DynamicBoxNode

	// Sink returns a sink registereed to the topology. It returns an error
	// when the topology doesn't have the sink.
	Sink(name string) (DynamicSinkNode, error)

	// Sinks returns all sinks registered to the topology. The map returned
	// from this method can safely be modified.
	Sinks() map[string]DynamicSinkNode
}

// DynamicSourceConfig has configuration parameters of a Source node.
type DynamicSourceConfig struct {
	// TODO: default state (paused?running?)
}

// DynamicBoxConfig has configuration parameters of a Box node.
type DynamicBoxConfig struct {
	// TODO: parallelism
}

// DynamicSinkConfig has configuration parameters of a Sink node.
type DynamicSinkConfig struct {
}

// DynamicNode is a node registered to a dynamic topology. It defines methods
// common to Source, Box, and Sink nodes.
type DynamicNode interface {
	// Type returns the type of the node, which can be NTSource, NTBox,
	// or NTSink. It's safe to convert DynamicNode to a specific node interface
	// corresponding to the returned NodeType. For example, if NTBox is
	// returned, the node can be converted to DynamicBoxNode with a type
	// assertion.
	Type() NodeType

	// Name returns the name of the node in the registered topology.
	Name() string

	// State returns the current state of the node.
	State() TopologyStateHolder

	// Stop stops the node. When the node is a source, Stop waits until the
	// source actually stops generating tuples. When the node is a box a sink,
	// it waits until the box or the sink is terminated.
	//
	// The node will not be removed from the topology after it stopped.
	Stop() error
}

// DynamicSourceNode is a Source registered to a dynamic topology.
type DynamicSourceNode interface {
	DynamicNode
}

// DynamicBoxNode is a Box registered to a dynamic topology.
type DynamicBoxNode interface {
	DynamicNode

	// Input adds a new input from a Source, another Box, or even the Box
	// itself. refname refers a name of node from which the Box want to receive
	// tuples. There must be a Source or a Box having the name.
	Input(refname string, config *BoxInputConfig) error

	// EnableGracefulStop activates a graceful stop mode. If it is enabled,
	// Stop method waits until the Box doesn't have an incoming tuple. The Box
	// doesn't wait until, for example, a source generates all tuples. It only
	// waits for the moment when the Box's input queue gets empty and stops
	// even if some inputs are about to send a new tuple to the Box.
	EnableGracefulStop()

	// StopOnDisconnect tells the Box that it may automatically stop when all
	// incoming connections (channels or pipes) are closed. After calling this
	// method, the Box can automatically stop even if Stop method isn't
	// explicitly called.
	StopOnDisconnect()
}

// DynamicSinkNode is a Sink registered to a dynamic topology.
type DynamicSinkNode interface {
	DynamicNode

	// Input adds a new input from a Source or a Box. refname refers a name of
	// node from which the Box want to receive tuples. There must be a Source
	// or a Box having the name.
	Input(refname string, config *SinkInputConfig) error

	// EnableGracefulStop activates a graceful stop mode. If it is enabled,
	// Stop method waits until the Sink doesn't have an incoming tuple. The Sink
	// doesn't wait until, for example, a source generates all tuples. It only
	// waits for the moment when the Sink's input queue gets empty and stops
	// even if some inputs are about to send a new tuple to the Sink.
	EnableGracefulStop()

	// StopOnDisconnect tells the Sink that it may automatically stop when all
	// incoming connections (channels or pipes) are closed. After calling this
	// method, the Sink can automatically stop even if Stop method isn't
	// explicitly called.
	StopOnDisconnect()
}

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
	// PausedOnStartup is a flag which indicates the initial state of the
	// source. If it is true, the source is paused. Otherwise, source runs just
	// after it is added to a topology.
	PausedOnStartup bool
}

// DynamicBoxConfig has configuration parameters of a Box node.
type DynamicBoxConfig struct {
	// TODO: parallelism
}

// DynamicSinkConfig has configuration parameters of a Sink node.
type DynamicSinkConfig struct {
}

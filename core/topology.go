package core

// Topology is a topology which can add Sources, Boxes, and Sinks
// dynamically. Boxes and Sinks can also add inputs dynamically from running
// Sources or Boxes.
type Topology interface {
	// Name returns the name of this topology.
	Name() string

	// Context returns the Context tied to the topology. It isn't always safe
	// to modify fields of the Context. However, some fields of it can handle
	// concurrent access properly. See the document of Context for details.
	Context() *Context

	// AddSource adds a Source to the topology. It will asynchronously call
	// source's GenerateStream and returns after the source becomes ready.
	// GenerateStream could be called lazily to avoid unnecessary computation.
	// GenerateStream and Stop might also be called when this method returns
	// an error. The caller must not call GenerateStream or Stop of the source.
	AddSource(name string, s Source, config *SourceConfig) (SourceNode, error)

	// AddBox adds a Box to the topology. It returns BoxNode and the
	// caller can configure inputs or other settings of the Box node through it.
	//
	// Don't add the same instance of a Box more than once even if they have
	// different names. However, if a Box has a reference counter and
	// its initialization and termination are done exactly once at proper
	// timing, it can be added multiple times when a builder supports duplicated
	// registration of the same instance of a Box.
	AddBox(name string, b Box, config *BoxConfig) (BoxNode, error)

	// AddSink adds a Sink to the topology. It returns SinkNode and the
	// caller can configure inputs or other settings of the Sink node through it.
	AddSink(name string, s Sink, config *SinkConfig) (SinkNode, error)

	// Remove removes a node from the topology. It returns NotExistError when
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
	// isn't relevant to those nodes have.
	State() TopologyStateHolder

	// TODO: low priority: Pause, Resume

	// Node returns a node registered to the topology. It returns NotExistError
	// when the topology doesn't have the node.
	Node(name string) (Node, error)

	// Nodes returns all nodes registered to the topology. The map returned
	// from this method can safely be modified.
	Nodes() map[string]Node

	// Source returns a source registered to the topology. It returns
	// NotExistError when the topology doesn't have the source.
	Source(name string) (SourceNode, error)

	// Sources returns all sources registered to the topology. The map returned
	// from this method can safely be modified.
	Sources() map[string]SourceNode

	// Box returns a box registered to the topology. It returns NotExistError
	// when the topology doesn't have the box.
	Box(name string) (BoxNode, error)

	// Boxes returns all boxes registered to the topology. The map returned
	// from this method can safely be modified.
	Boxes() map[string]BoxNode

	// Sink returns a sink registereed to the topology. It returns NotExistError
	// when the topology doesn't have the sink.
	Sink(name string) (SinkNode, error)

	// Sinks returns all sinks registered to the topology. The map returned
	// from this method can safely be modified.
	Sinks() map[string]SinkNode
}

// SourceConfig has configuration parameters of a Source node.
type SourceConfig struct {
	// PausedOnStartup is a flag which indicates the initial state of the
	// source. If it is true, the source is paused. Otherwise, source runs just
	// after it is added to a topology.
	PausedOnStartup bool

	// RemoveOnStop is a flag which indicates the stop state of the topology.
	// If it is true, the source is removed.
	RemoveOnStop bool

	// Meta contains meta information of the source. This field won't be used
	// by core package and application can store any form of information
	// related to the source.
	Meta interface{}
}

// BoxConfig has configuration parameters of a Box node.
type BoxConfig struct {
	// TODO: parallelism

	// RemoveOnStop is a flag which indicates the stop state of the topology.
	// If it is true, the box is removed.
	RemoveOnStop bool

	// Meta contains meta information of the box. This field won't be used
	// by core package and application can store any form of information
	// related to the box.
	Meta interface{}
}

// SinkConfig has configuration parameters of a Sink node.
type SinkConfig struct {
	// RemoveOnStop is a flag which indicates the stop state of the topology.
	// If it is true, the sink is removed.
	RemoveOnStop bool

	// Meta contains meta information of the sink. This field won't be used
	// by core package and application can store any form of information
	// related to the sink.
	Meta interface{}
}

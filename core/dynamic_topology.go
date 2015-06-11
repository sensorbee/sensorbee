package core

type DynamicTopology interface {
	AddSource(name string, s Source, config *DynamicSourceConfig) (DynamicSourceNode, error)
	AddBox(name string, b Box, config *DynamicBoxConfig) (DynamicBoxNode, error)
	AddSink(name string, s Sink, config *DynamicSinkConfig) (DynamicSinkNode, error)
	Remove(name string) error

	Stop() error
	State() TopologyStateHolder

	// TODO: low priority: Pause, Resume

	// Node returns a node registered to the topology. It returns an error
	// when the topology doesn't have the node.
	Node(name string) (DynamicNode, error)

	// Nodes returns all nodes registered to the topology.
	Nodes() map[string]DynamicNode

	// Source returns a source registered to the topology. It returns an error
	// when the topology doesn't have the source.
	Source(name string) (DynamicSourceNode, error)

	// Sources returns all sources registered to the topology.
	Sources() map[string]DynamicSourceNode

	// Box returns a box registered to the topology. It returns an error when
	// the topology doesn't have the box.
	Box(name string) (DynamicBoxNode, error)

	// Boxes returns all boxes registered to the topology.
	Boxes() map[string]DynamicBoxNode

	// Sink returns a sink registereed to the topology. It returns an error
	// when the topology doesn't have the sink.
	Sink(name string) (DynamicSinkNode, error)

	// Sinks returns all sinks registered to the topology.
	Sinks() map[string]DynamicSinkNode
}

type DynamicSourceConfig struct {
	// TODO: default state (paused?running?)
}

type DynamicBoxConfig struct {
	// TODO: parallelism
}

type DynamicSinkConfig struct {
}

type DynamicNode interface {
	Type() NodeType
	Name() string
	State() TopologyStateHolder
	Stop() error
}

type DynamicSourceNode interface {
	DynamicNode
}

type DynamicBoxNode interface {
	DynamicNode
	Input(refname string, config *BoxInputConfig) error
}

type DynamicSinkNode interface {
	DynamicNode
	Input(refname string, config *SinkInputConfig) error
}

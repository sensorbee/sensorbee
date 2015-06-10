package core

type DynamicTopology interface {
	// TODO: Run
	// TODO: Stop
	// TODO: State

	// TODO: low priority: Pause, Resume
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

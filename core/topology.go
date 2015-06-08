package core

// TopologyState represents a status of a topology.
type TopologyState int

const (
	// TSInitialized means that a topology is just initialized and ready
	// to be run.
	TSInitialized TopologyState = iota

	// TSStarting means a topology is now booting itself and will run shortly.
	TSStarting

	// TSRunning means a topology is currently running and emitting tuples
	// to sinks.
	TSRunning

	// TSStopping means a topology is stopping all sources and closing
	// channels between sources, boxes, and sinks.
	TSStopping

	// TSStopped means a topology is stopped. A stopped topology doesn't
	// have to be able to run again.
	TSStopped
)

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

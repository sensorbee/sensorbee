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

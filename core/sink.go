package core

// A Sink describes a location that data can be written to after it
// was processed by a topology, i.e., it represents an entity
// outside of the topology (e.g., a fluentd instance).
type Sink interface {
	Writer
}

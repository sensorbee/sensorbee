package core

import (
	"pfi/sensorbee/sensorbee/data"
)

// Updater represents an entity that can update its configuration
// parameters (in particular SourceNode, SinkNode and SharedState
// instances).
type Updater interface {
	// Update updates the configuration parameters of this entity.
	// It is the updater's responsibility to check the validity
	// (e.g., data type and value) of the parameters.
	Update(data.Map) error
}

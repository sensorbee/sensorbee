package response

import (
	"gopkg.in/sensorbee/sensorbee.v0/core"
)

// Topology is a part of the response which topologies.show action returns.
type Topology struct {
	// Name is the name of the topology.
	Name string `json:"name"`
}

// NewTopology creates a new response of a topology.
func NewTopology(t core.Topology) *Topology {
	return &Topology{
		Name: t.Name(),
	}
}

// TODO: add created_at/updated_at
// TODO: add other information

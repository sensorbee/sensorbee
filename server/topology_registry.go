package server

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"os"
	"strings"
	"sync"
)

// TopologyRegistry is a registry of topologies managed in th server.
type TopologyRegistry interface {
	// Register registers a new topology. If the registry already has a
	// topology having the same name, this method fails and returns
	// os.ErrExist.
	Register(name string, tb *bql.TopologyBuilder) error

	// Lookup returns a topology having the name. It returns core.NotExistError
	// if it doesn't have the topology.
	Lookup(name string) (*bql.TopologyBuilder, error)

	// List returns all topologies the registry has. The caller can safely
	// modify the map returned from this method.
	List() (map[string]*bql.TopologyBuilder, error)

	// Unregister removes a creator from the registry. It returns a removed
	// topology. If the registry doesn't have a topology, it returns
	// core.NotExistError. If it failed to remove the topology, itreturns an error.
	//
	// Unregister doesn't stop the topology when it's removed. It's the caller's
	// responsibility to correctly stop it.
	Unregister(name string) (*bql.TopologyBuilder, error)
}

type defaultTopologyRegistry struct {
	m          sync.RWMutex
	topologies map[string]*bql.TopologyBuilder
}

// NewDefaultTopologyRegistry returns a default implementation of
// TopologyRegistry.
func NewDefaultTopologyRegistry() TopologyRegistry {
	return &defaultTopologyRegistry{
		topologies: map[string]*bql.TopologyBuilder{},
	}
}

func (r *defaultTopologyRegistry) Register(name string, tb *bql.TopologyBuilder) error {
	r.m.Lock()
	defer r.m.Unlock()

	n := strings.ToLower(name)
	if _, ok := r.topologies[n]; ok {
		return os.ErrExist
	}
	r.topologies[n] = tb
	return nil
}

func (r *defaultTopologyRegistry) Lookup(name string) (*bql.TopologyBuilder, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	if tb, ok := r.topologies[strings.ToLower(name)]; ok {
		return tb, nil
	}
	return nil, core.NotExistError(fmt.Errorf("topology '%v' is not registered", name))
}

func (r *defaultTopologyRegistry) List() (map[string]*bql.TopologyBuilder, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	m := make(map[string]*bql.TopologyBuilder, len(r.topologies))
	for n, tb := range r.topologies {
		m[n] = tb
	}
	return m, nil
}

func (r *defaultTopologyRegistry) Unregister(name string) (*bql.TopologyBuilder, error) {
	r.m.Lock()
	defer r.m.Unlock()
	n := strings.ToLower(name)
	tb, ok := r.topologies[n]
	if !ok {
		return nil, core.NotExistError(fmt.Errorf("topology '%v' is not registered", name))
	}
	delete(r.topologies, n)
	return tb, nil
}

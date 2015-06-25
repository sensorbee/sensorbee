package bql

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"sync"
)

type SourceCreator interface {
	CreateSource(ctx *core.Context, params tuple.Map) (core.Source, error)
}

type sourceCreatorFunc func(*core.Context, tuple.Map) (core.Source, error)

func (f sourceCreatorFunc) CreateSource(ctx *core.Context, params tuple.Map) (core.Source, error) {
	return f(ctx, params)
}

func SourceCreatorFunc(f func(*core.Context, tuple.Map) (core.Source, error)) SourceCreator {
	return sourceCreatorFunc(f)
}

// SourceCreatorRegistry manages creators of Sources.
type SourceCreatorRegistry interface {
	// Register adds a Source creator to the registry. It returns an error if
	// the type name is already registered.
	Register(typeName string, c SourceCreator) error

	// Lookup returns a Source creator having the type name. It returns an error
	// if it doesn't have the creator.
	Lookup(typeName string) (SourceCreator, error)

	// List returns all creators the registry has. The caller can safely modify
	// the map returned from this method.
	List() (map[string]SourceCreator, error)

	// Unregister removes a creator from the registry. It doesn't return error
	// when the registry doesn't have a creator having the type name.
	//
	// The registry itself doesn't support cascading delete. It should properly
	// done by the caller.
	Unregister(typeName string) error
}

type defaultSourceCreatorRegistry struct {
	m        sync.RWMutex
	creators map[string]SourceCreator
}

// NewDefaultSourceCreatorRegistry returns a SourceCreatorRegistry having a
// default implementation.
func NewDefaultSourceCreatorRegistry() SourceCreatorRegistry {
	return &defaultSourceCreatorRegistry{
		creators: map[string]SourceCreator{},
	}
}

func (r *defaultSourceCreatorRegistry) Register(typeName string, c SourceCreator) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.creators[typeName]; ok {
		return fmt.Errorf("source type '%v' is already registered", typeName)
	}
	r.creators[typeName] = c
	return nil
}

func (r *defaultSourceCreatorRegistry) Lookup(typeName string) (SourceCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if c, ok := r.creators[typeName]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("source type '%v' is not registered", typeName)
}

func (r *defaultSourceCreatorRegistry) List() (map[string]SourceCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	m := make(map[string]SourceCreator, len(r.creators))
	for t, c := range r.creators {
		m[t] = c
	}
	return m, nil
}

func (r *defaultSourceCreatorRegistry) Unregister(typeName string) error {
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.creators, typeName)
	return nil
}

var (
	globalSourceCreatorRegistry = NewDefaultSourceCreatorRegistry()
)

// RegisterGlobalSourceCreator adds a SourceCreator which can be referred from
// alltopologies. SourceCreators registered after running topologies might not
// be seen by those topologies. Call it from init functions to avoid such
// conditions.
func RegisterGlobalSourceCreator(typeName string, c SourceCreator) error {
	return globalSourceCreatorRegistry.Register(typeName, c)
}

// CopyGlobalSourceRegistry creates a new independent copy of the global
// SourceCreatorRegistry.
func CopyGlobalSourceRegistry() (SourceCreatorRegistry, error) {
	r := NewDefaultSourceCreatorRegistry()
	m, err := globalSourceCreatorRegistry.List()
	if err != nil {
		return nil, err
	}

	for t, c := range m {
		if err := r.Register(t, c); err != nil {
			return nil, err
		}
	}
	return r, nil
}

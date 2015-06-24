package bql

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
	"sync"
)

// UDSCreator creates a User Defined State based on core.SharedState.
type UDSCreator interface {
	// CreateState creates an instance of the state type. CreateState must not
	// call core.SharedState.Init.
	CreateState(ctx *core.Context, params tuple.Map) (core.SharedState, error)
}

type udsCreatorFunc func(*core.Context, tuple.Map) (core.SharedState, error)

func (f udsCreatorFunc) CreateState(ctx *core.Context, params tuple.Map) (core.SharedState, error) {
	return f(ctx, params)
}

// UDSCreatorFunc creates a UDSCreator from a function.
func UDSCreatorFunc(f func(*core.Context, tuple.Map) (core.SharedState, error)) UDSCreator {
	return udsCreatorFunc(f)
}

// UDSCreatorRegistry manages creators of UDSs.
type UDSCreatorRegistry interface {
	// Register adds a UDS creator to the registry. It returns an error if
	// the type name is already registered.
	Register(typeName string, c UDSCreator) error

	// Lookup returns a UDS creator having the type name. It returns an error
	// if it doesn't have the creator.
	Lookup(typeName string) (UDSCreator, error)

	// List returns all creators the registry has. The caller can safely modify
	// the map returned from this method.
	List() (map[string]UDSCreator, error)

	// Unregister removes a creator from the registry. It doesn't return error
	// when the registry doesn't have a creator having the type name.
	//
	// The registry itself doesn't support cascading delete. It should properly
	// done by the caller.
	Unregister(typeName string) error
}

type defaultUDSCreatorRegistry struct {
	m        sync.RWMutex
	creators map[string]UDSCreator
}

// NewDefaultUDSCreatorRegistry returns a UDSCreatorRegistry having a default
// implementation.
func NewDefaultUDSCreatorRegistry() UDSCreatorRegistry {
	return &defaultUDSCreatorRegistry{
		creators: map[string]UDSCreator{},
	}
}

func (r *defaultUDSCreatorRegistry) Register(typeName string, c UDSCreator) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.creators[typeName]; ok {
		return fmt.Errorf("UDS type '%v' is already registered", typeName)
	}
	r.creators[typeName] = c
	return nil
}

func (r *defaultUDSCreatorRegistry) Lookup(typeName string) (UDSCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if c, ok := r.creators[typeName]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("UDS type '%v' is not found", typeName)
}

func (r *defaultUDSCreatorRegistry) List() (map[string]UDSCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	m := make(map[string]UDSCreator, len(r.creators))
	for t, c := range r.creators {
		m[t] = c
	}
	return m, nil
}

func (r *defaultUDSCreatorRegistry) Unregister(typeName string) error {
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.creators, typeName)
	return nil
}

var (
	globalUDSCreatorRegistry = NewDefaultUDSCreatorRegistry()
)

// RegisterGlobalUDSCreator adds a UDSCreator which can be referred from all
// topologies. UDSCreators registered after running topologies might not be
// seen by those topologies. Call it from init functions to avoid such
// conditions.
func RegisterGlobalUDSCreator(typeName string, c UDSCreator) error {
	return globalUDSCreatorRegistry.Register(typeName, c)
}

func CopyGlobalUDSRegistry() (UDSCreatorRegistry, error) {
	r := NewDefaultUDSCreatorRegistry()
	m, err := globalUDSCreatorRegistry.List()
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

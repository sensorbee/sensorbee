package bql

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
)

type SinkCreator interface {
	CreateSink(ctx *core.Context, params data.Map) (core.Sink, error)
}

type sinkCreatorFunc func(*core.Context, data.Map) (core.Sink, error)

func (f sinkCreatorFunc) CreateSink(ctx *core.Context, params data.Map) (core.Sink, error) {
	return f(ctx, params)
}

func SinkCreatorFunc(f func(*core.Context, data.Map) (core.Sink, error)) SinkCreator {
	return sinkCreatorFunc(f)
}

// SinkCreatorRegistry manages creators of Sinks.
type SinkCreatorRegistry interface {
	// Register adds a Sink creator to the registry. It returns an error if
	// the type name is already registered.
	Register(typeName string, c SinkCreator) error

	// Lookup returns a Sink creator having the type name. It returns an error
	// if it doesn't have the creator.
	Lookup(typeName string) (SinkCreator, error)

	// List returns all creators the registry has. The caller can safely modify
	// the map returned from this method.
	List() (map[string]SinkCreator, error)

	// Unregister removes a creator from the registry. It doesn't return error
	// when the registry doesn't have a creator having the type name.
	//
	// The registry itself doesn't support cascading delete. It should properly
	// done by the caller.
	Unregister(typeName string) error
}

type defaultSinkCreatorRegistry struct {
	m        sync.RWMutex
	creators map[string]SinkCreator
}

// NewDefaultSinkCreatorRegistry returns a SinkCreatorRegistry having a
// default implementation.
func NewDefaultSinkCreatorRegistry() SinkCreatorRegistry {
	return &defaultSinkCreatorRegistry{
		creators: map[string]SinkCreator{},
	}
}

func (r *defaultSinkCreatorRegistry) Register(typeName string, c SinkCreator) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.creators[typeName]; ok {
		return fmt.Errorf("sink type '%v' is already registered", typeName)
	}
	r.creators[typeName] = c
	return nil
}

func (r *defaultSinkCreatorRegistry) Lookup(typeName string) (SinkCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if c, ok := r.creators[typeName]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("sink type '%v' is not registered", typeName)
}

func (r *defaultSinkCreatorRegistry) List() (map[string]SinkCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	m := make(map[string]SinkCreator, len(r.creators))
	for t, c := range r.creators {
		m[t] = c
	}
	return m, nil
}

func (r *defaultSinkCreatorRegistry) Unregister(typeName string) error {
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.creators, typeName)
	return nil
}

var (
	globalSinkCreatorRegistry = NewDefaultSinkCreatorRegistry()
)

// RegisterGlobalSinkCreator adds a SinkCreator which can be referred from
// alltopologies. SinkCreators registered after running topologies might not
// be seen by those topologies. Call it from init functions to avoid such
// conditions.
func RegisterGlobalSinkCreator(typeName string, c SinkCreator) error {
	return globalSinkCreatorRegistry.Register(typeName, c)
}

// MustRegisterGlobalSinkCreator is like RegisterGlobalSinkCreator but panics
// if an error occurred.
func MustRegisterGlobalSinkCreator(typeName string, c SinkCreator) {
	if err := globalSinkCreatorRegistry.Register(typeName, c); err != nil {
		panic(fmt.Errorf("bql.MustRegisterGlobalSinkCreator: cannot register '%v': %v", typeName, err))
	}
}

// CopyGlobalSinkCreatorRegistry creates a new independent copy of the global
// SinkCreatorRegistry.
func CopyGlobalSinkCreatorRegistry() (SinkCreatorRegistry, error) {
	r := NewDefaultSinkCreatorRegistry()
	m, err := globalSinkCreatorRegistry.List()
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

func createSharedStateSink(ctx *core.Context, params data.Map) (core.Sink, error) {
	// Get only name parameter from params
	name, ok := params["name"]
	if !ok {
		return nil, fmt.Errorf("cannot find 'name' field")
	}
	nameStr, err := data.AsString(name)
	if err != nil {
		return nil, err
	}

	return core.NewSharedStateSink(ctx, nameStr)
}

func init() {
	RegisterGlobalSinkCreator("uds", SinkCreatorFunc(createSharedStateSink))
}

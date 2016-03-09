package udf

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io"
	"strings"
	"sync"
)

// UDSCreator creates a User Defined State based on core.SharedState.
type UDSCreator interface {
	// CreateState creates an instance of the state type. CreateState must not
	// call core.SharedState.Init.
	CreateState(ctx *core.Context, params data.Map) (core.SharedState, error)
}

// UDSLoader loads a User Defined State from saved data. A UDS cannot be loaded
// if a UDSCreator doesn't implement UDSLoader even if the UDS implements
// core.LoadableSharedState.
//
// When a UDS isn't created yet, UDSLoader.LoadState will be used to load the
// state and core.LoadableSharedState.Load will not be used.
//
// When a UDS is already created or loaded and it implements
// core.LoadableSharedState, its Load method is called to load a model and
// UDSLoader.LoadState will not be called. If a UDS doesn't implement
// core.LoadableSharedState but UDSLoader is provided for its type, then
// UDSLoader.LoadState creates a new instance and the previous instance is
// replaced with it, which means loading the UDS could consume twice as much
// memory as core.LoadableSharedState.Load does. When a UDS doesn't implement
// core.LoadableSharedState and its UDSCreator doesn't implement UDSLoader,
// the UDS cannot be loaded.
type UDSLoader interface {
	UDSCreator

	// LoadState loads a state from saved data. The saved data can be read from
	// io.Reader. Parameters given by SET clause are passed as params.
	LoadState(ctx *core.Context, r io.Reader, params data.Map) (core.SharedState, error)
}

type udsCreatorFunc func(*core.Context, data.Map) (core.SharedState, error)

func (f udsCreatorFunc) CreateState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	return f(ctx, params)
}

// UDSCreatorFunc creates a UDSCreator from a function.
func UDSCreatorFunc(f func(*core.Context, data.Map) (core.SharedState, error)) UDSCreator {
	return udsCreatorFunc(f)
}

// UDSCreatorRegistry manages creators of UDSs.
type UDSCreatorRegistry interface {
	// Register adds a UDS creator to the registry. It returns an error if
	// the type name is already registered.
	Register(typeName string, c UDSCreator) error

	// Lookup returns a UDS creator having the type name. It returns
	// core.NotExistError if it doesn't have the creator.
	Lookup(typeName string) (UDSCreator, error)

	// List returns all creators the registry has. The caller can safely modify
	// the map returned from this method.
	List() (map[string]UDSCreator, error)

	// Unregister removes a creator from the registry. It returns core.NotExistError
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
	if err := core.ValidateSymbol(typeName); err != nil {
		return fmt.Errorf("invalid name for UDS type: %s", err.Error())
	}

	r.m.Lock()
	defer r.m.Unlock()

	lowerName := strings.ToLower(typeName)
	if _, ok := r.creators[lowerName]; ok {
		return fmt.Errorf("UDS type '%v' is already registered", typeName)
	}
	r.creators[lowerName] = c
	return nil
}

func (r *defaultUDSCreatorRegistry) Lookup(typeName string) (UDSCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if c, ok := r.creators[strings.ToLower(typeName)]; ok {
		return c, nil
	}
	return nil, core.NotExistError(fmt.Errorf("UDS type '%v' is not found", typeName))
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
	tn := strings.ToLower(typeName)
	if _, ok := r.creators[tn]; !ok {
		return core.NotExistError(fmt.Errorf("UDS type '%v' is not found", typeName))
	}
	delete(r.creators, tn)
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

// MustRegisterGlobalUDSCreator is like RegisterGlobalUDSCreator
// but panics if an error occurred.
func MustRegisterGlobalUDSCreator(typeName string, c UDSCreator) {
	if err := globalUDSCreatorRegistry.Register(typeName, c); err != nil {
		panic(fmt.Errorf("udf.MustRegisterGlobalUDSCreator: cannot register '%v': %v", typeName, err))
	}
}

// CopyGlobalUDSCreatorRegistry creates a new independent copy of the global
// UDSCreatorRegistry.
func CopyGlobalUDSCreatorRegistry() (UDSCreatorRegistry, error) {
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

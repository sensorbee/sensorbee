package bql

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"strings"
	"sync"
)

// IOParams has parameters for IO plugins.
type IOParams struct {
	// TypeName is the name of the type registered to SensorBee.
	TypeName string

	// Name is the name of the instance specified in a CREATE statement.
	Name string
}

// SourceCreator is an interface which creates instances of a Source.
type SourceCreator interface {
	// CreateSource creates a new Source instance using given parameters.
	CreateSource(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Source, error)
}

type sourceCreatorFunc func(*core.Context, *IOParams, data.Map) (core.Source, error)

func (f sourceCreatorFunc) CreateSource(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Source, error) {
	return f(ctx, ioParams, params)
}

// SourceCreatorFunc creates a SourceCreator from a function.
func SourceCreatorFunc(f func(*core.Context, *IOParams, data.Map) (core.Source, error)) SourceCreator {
	return sourceCreatorFunc(f)
}

// SourceCreatorRegistry manages creators of Sources.
type SourceCreatorRegistry interface {
	// Register adds a Source creator to the registry. It returns an error if
	// the type name is already registered.
	Register(typeName string, c SourceCreator) error

	// Lookup returns a Source creator having the type name. It returns
	// core.NotExistError if it doesn't have the creator.
	Lookup(typeName string) (SourceCreator, error)

	// List returns all creators the registry has. The caller can safely modify
	// the map returned from this method.
	List() (map[string]SourceCreator, error)

	// Unregister removes a creator from the registry. It returns core.NotExistError
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

	if err := core.ValidateSymbol(typeName); err != nil {
		return fmt.Errorf("invalid name for source type: %s", err.Error())
	}
	lowerName := strings.ToLower(typeName)
	if _, ok := r.creators[lowerName]; ok {
		return fmt.Errorf("source type '%v' is already registered", typeName)
	}
	r.creators[lowerName] = c
	return nil
}

func (r *defaultSourceCreatorRegistry) Lookup(typeName string) (SourceCreator, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if c, ok := r.creators[strings.ToLower(typeName)]; ok {
		return c, nil
	}
	return nil, core.NotExistError(fmt.Errorf("source type '%v' is not registered", typeName))
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
	tn := strings.ToLower(typeName)
	if _, ok := r.creators[tn]; !ok {
		return core.NotExistError(fmt.Errorf("source type '%v' is not registered", typeName))
	}
	delete(r.creators, tn)
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

// MustRegisterGlobalSourceCreator is like RegisterGlobalSourceCreator but
// panics if an error occurred.
func MustRegisterGlobalSourceCreator(typeName string, c SourceCreator) {
	if err := globalSourceCreatorRegistry.Register(typeName, c); err != nil {
		panic(fmt.Errorf("udf.MustRegisterGlobalSourceCreator: cannot register '%v': %v", typeName, err))
	}
}

// CopyGlobalSourceCreatorRegistry creates a new independent copy of the global
// SourceCreatorRegistry.
func CopyGlobalSourceCreatorRegistry() (SourceCreatorRegistry, error) {
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

package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/tuple"
	"sync"
)

// SharedState is a state which nodes in a topology can access. It can be a
// machine learning model, a data structure for aggregation (like a histgram),
// a configuration information for specific Boxes, and so on.
//
// SharedState doesn't have methods to read it's internal data because internal
// data representation heavily depends on each SharedState implementation. The
// easiest way to use SharedState from a component is to obtain the actual
// data type via the type assertion. See examples to learn more about how to
// use it.
type SharedState interface {
	// TypeName returns the name of SharedState's type. There can be multiple
	// instances of the same type. For example, a user can use "classifier"
	// as a type name and create SharedStates "category_classifier" and
	// "age_classifier" having the type name.
	TypeName() string

	// Init initialize the SharedState. Each SharedState can allocate resources
	// required for its computation. SharedState doesn't have to store ctx in it
	// because all method which are called from topologies have ctx arguments.
	//
	// Init isn't called concurrently.
	Init(ctx *Context) error

	// Write writes a tuple to the state. How tuples are processed depends on
	// each SharedState. For example, a machine learning model might use a
	// tuple as a training data, and another state could compute the average
	// of a specific field. Write may return fatal or temporary errors as
	// Box.Process does. See the documentation of Box.Process for details.
	//
	// Write method might be called after Terminate method is called. When it
	// occurs, Write should return an error. Also, Write and Terminate can be
	// called concurrently
	Write(ctx *Context, t *tuple.Tuple) error

	// Terminate finalizes the state. The state can no longer be used after
	// this method is called. This method doesn't have to be idempotent.
	// Terminate won't be called when Init fails.
	//
	// Write method might be called after Terminate method is called. When it
	// occurs, Write should return an error. Also, Write and Terminate can be
	// called concurrently.
	Terminate(ctx *Context) error
}

// TODO: Add MixiableSharedState interface
// TODO: Add SerializableSharedState (or simply Serializable) to support recovery

// SharedStateRegistry manages SharedState with names assigned to each state.
type SharedStateRegistry interface {
	// Add adds a state to the registry. It fails if the registry already has
	// a state having the same name. Add also calls SharedState.Init. If it
	// fails Add returns an error and doesn't register the SharedState. The
	// caller doesn't have to call Terminate on failure.
	//
	// Don't add the same instance of SharedState more than once to registries.
	// Otherwise, Init and Terminate methods of the state will be called
	// multiple times.
	Add(ctx *Context, name string, s SharedState) error

	// Get returns a SharedState having the name in the registry. It returns
	// an error if the registry doesn't have the state.
	Get(ctx *Context, name string) (SharedState, error)

	// List returns a map containing all SharedState the registry has.
	// The map returned from this method can safely be modified.
	List(ctx *Context) (map[string]SharedState, error)

	// Remove removes a SharedState the registry has. It automatically
	// terminates the state. If SharedState.Terminate failed, Remove returns an
	// error. However, even if it returns an error, the state is removed from
	// the registry.
	//
	// Remove also returns the removed SharedState if the registry has it. When
	// SharedState.Terminate fails, Remove returns both the removed SharedState
	// and an error. If the registry doesn't have a SharedState having the name,
	// it returns a nil SharedState and a nil error.
	Remove(ctx *Context, name string) (SharedState, error)
}

type defaultSharedStateRegistry struct {
	m      sync.RWMutex
	states map[string]SharedState
}

// NewDefaultSharedStateRegistry create a default registry of SharedStates.
func NewDefaultSharedStateRegistry() SharedStateRegistry {
	return &defaultSharedStateRegistry{
		states: map[string]SharedState{},
	}
}

func (r *defaultSharedStateRegistry) Add(ctx *Context, name string, s SharedState) error {
	if err := s.Init(ctx); err != nil {
		return err
	}

	err := func() error {
		r.m.Lock()
		defer r.m.Unlock()
		if _, ok := r.states[name]; ok {
			return fmt.Errorf("the registry already has a state '%v'", name)
		}
		r.states[name] = s
		return nil
	}()
	if err != nil {
		if err := s.Terminate(ctx); err != nil {
			ctx.Logger.Log(Error, "Cannot terminate state which couldn't be added to the registry due to name duplication '%v': %v", name, err)
		}
		return err
	}
	return nil
}

func (r *defaultSharedStateRegistry) Get(ctx *Context, name string) (SharedState, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if s, ok := r.states[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' was not found", name)
}

func (r *defaultSharedStateRegistry) List(ctx *Context) (map[string]SharedState, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	m := make(map[string]SharedState, len(r.states))
	for n, s := range r.states {
		m[n] = s
	}
	return m, nil
}

func (r *defaultSharedStateRegistry) Remove(ctx *Context, name string) (SharedState, error) {
	s := func() SharedState {
		r.m.Lock()
		defer r.m.Unlock()
		if s, ok := r.states[name]; ok {
			delete(r.states, name)
			return s
		}
		return nil
	}()
	if s == nil {
		return nil, nil
	}

	if err := s.Terminate(ctx); err != nil {
		return s, err
	}
	return s, nil
}

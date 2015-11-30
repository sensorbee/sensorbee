package core

import (
	"fmt"
	"io"
	"pfi/sensorbee/sensorbee/data"
	"sync"
)

// SharedState is a state which nodes in a topology can access. It can be a
// machine learning model, a data structure for aggregation (like a histogram),
// a configuration information for specific Boxes, and so on.
//
// SharedState doesn't have methods to read it's internal data because internal
// data representation heavily depends on each SharedState implementation. The
// easiest way to use SharedState from a component is to obtain the actual
// data type via the type assertion. See examples to learn more about how to
// use it.
//
// If a SharedState also implements Writer interface, it can be updated via
// SharedStateSink. Write method in it writes a tuple to the state. How tuples
// are processed depends on each SharedState. For example, a machine learning
// model might use a tuple as a training data, and another state could compute
// the average of a specific field. Write may return fatal or temporary errors
// as Box.Process does. See the documentation of Box.Process for details.
//
// Write method might be called after Terminate method is called. When it
// occurs, Write should return an error. Also, Write and Terminate can be
// called concurrently.
type SharedState interface {
	// Terminate finalizes the state. The state can no longer be used after
	// this method is called. This method doesn't have to be idempotent.
	//
	// Write or other methods the actual instance has might be called after
	// Terminate method is called. When it occurs, they should return an error.
	// Also, Terminate and them can be called concurrently.
	Terminate(ctx *Context) error
}

// SavableSharedState is a SharedState which can be persisted through Save
// method. Providing forward/backward compatibility of the saved file format
// is the responsibility of the author of the state.
//
// Because the best way of implementing Load method depends on each SharedState,
// it doesn't always have to be provided with Save method.
type SavableSharedState interface {
	SharedState

	// Save writes data of the state to a given writer. Save receives parameters
	// which are used to customize the behavior of the method. Parameters are
	// defined by each component and there's no common definition.
	//
	// Save and other methods can be called concurrently.
	Save(ctx *Context, w io.Writer, params data.Map) error
}

// LoadableSharedState is a SharedState which can be persisted through Save
// and Load method.
type LoadableSharedState interface {
	SavableSharedState

	// Load overwrites the state with save data. Parameters don't have to be
	// same as Save's parameters. They can even be completely different.
	// There MUST NOT be a required parameter. Values of required parameters
	// should be saved with the state itself.
	//
	// Load and other methods including Save can be called concurrently.
	Load(ctx *Context, r io.Reader, params data.Map) error
}

// TODO: Add MixiableSharedState interface

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
	Add(name, typeName string, s SharedState) error

	// Get returns a SharedState having the name in the registry. It returns
	// NotExistError if the registry doesn't have the state.
	Get(name string) (SharedState, error)

	// Type returns a type of a SharedState. It returns NotExistError if the
	// registry doesn't have the state.
	Type(name string) (string, error)

	// Replace replaces the previous SharedState instance with a new instance.
	// The previous instance is returned on success if any. The previous state
	// will not be terminated by the registry and the caller must call
	// Terminate. The type name must be same as the previous state's type name.
	//
	// The given SharedState is terminated when it cannot be replaced.
	Replace(name, typeName string, s SharedState) (SharedState, error)

	// List returns a map containing all SharedState the registry has.
	// The map returned from this method can safely be modified.
	List() (map[string]SharedState, error)

	// Remove removes a SharedState the registry has. It automatically
	// terminates the state. If SharedState.Terminate failed, Remove returns an
	// error. However, even if it returns an error, the state is removed from
	// the registry.
	//
	// Remove also returns the removed SharedState if the registry has it. When
	// SharedState.Terminate fails, Remove returns both the removed SharedState
	// and an error. If the registry doesn't have a SharedState having the name,
	// it returns a nil SharedState and NotExistError.
	Remove(name string) (SharedState, error)
}

type defaultSharedStateInfo struct {
	state    SharedState
	typeName string
}

type defaultSharedStateRegistry struct {
	ctx    *Context
	m      sync.RWMutex
	states map[string]*defaultSharedStateInfo
}

// NewDefaultSharedStateRegistry create a default registry of SharedStates.
func NewDefaultSharedStateRegistry(ctx *Context) SharedStateRegistry {
	return &defaultSharedStateRegistry{
		ctx:    ctx,
		states: map[string]*defaultSharedStateInfo{},
	}
}

func (r *defaultSharedStateRegistry) Add(name, typeName string, s SharedState) error {
	err := func() error {
		r.m.Lock()
		defer r.m.Unlock()
		if _, ok := r.states[name]; ok {
			return fmt.Errorf("the registry already has a state '%v'", name)
		}
		r.states[name] = &defaultSharedStateInfo{
			state:    s,
			typeName: typeName,
		}
		return nil
	}()
	if err != nil {
		if err := r.closeSharedState(s); err != nil {
			r.ctx.ErrLog(err).WithField("state_name", name).
				Errorf("Cannot terminate a state which couldn't be added to the registry due to name duplication")
		}
		return err // This is the original error
	}
	return nil
}

func (r *defaultSharedStateRegistry) closeSharedState(s SharedState) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if er, ok := e.(error); ok {
				err = er
			} else {
				err = fmt.Errorf("SharedState.Terminate panicked: %v", e)
			}
		}
	}()
	return s.Terminate(r.ctx)
}

func (r *defaultSharedStateRegistry) Get(name string) (SharedState, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if s, ok := r.states[name]; ok {
		return s.state, nil
	}
	return nil, NotExistError(fmt.Errorf("state '%v' was not found", name))
}

func (r *defaultSharedStateRegistry) Type(name string) (string, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if s, ok := r.states[name]; ok {
		return s.typeName, nil
	}
	return "", NotExistError(fmt.Errorf("state '%v' was not found", name))
}

func (r *defaultSharedStateRegistry) Replace(name, typeName string, s SharedState) (SharedState, error) {
	r.m.Lock()
	defer r.m.Unlock()
	prev, ok := r.states[name]
	if ok {
		if prev.typeName != typeName {
			if err := r.closeSharedState(s); err != nil {
				r.ctx.ErrLog(err).WithField("state_name", name).
					WithField("state_type", typeName).WithField("prev_state_type", prev.typeName).
					Errorf("Cannot terminate a state which couldn't be replaced due to a type mismatch")
			}
			return nil, fmt.Errorf("state '%v' has a different type from the previous state's type", name)
		}
	}
	r.states[name] = &defaultSharedStateInfo{
		state:    s,
		typeName: typeName,
	}
	if prev == nil {
		return nil, nil
	}
	return prev.state, nil
}

func (r *defaultSharedStateRegistry) List() (map[string]SharedState, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	m := make(map[string]SharedState, len(r.states))
	for n, s := range r.states {
		m[n] = s.state
	}
	return m, nil
}

func (r *defaultSharedStateRegistry) Remove(name string) (SharedState, error) {
	s := func() SharedState {
		r.m.Lock()
		defer r.m.Unlock()
		if s, ok := r.states[name]; ok {
			delete(r.states, name)
			return s.state
		}
		return nil
	}()
	if s == nil {
		return nil, NotExistError(fmt.Errorf("state '%v' was not found", name))
	}

	if err := s.Terminate(r.ctx); err != nil {
		return s, err
	}
	return s, nil
}

// sharedStateSink represents a shared state. sharedStateSink refers to a shared state by name.
type sharedStateSink struct {
	name string
}

// NewSharedStateSink creates a sink that writes to SharedState.
func NewSharedStateSink(ctx *Context, name string) (Sink, error) {
	// Get SharedState by name
	state, err := ctx.SharedStates.Get(name)
	if err != nil {
		return nil, err
	}

	// It fails if the shared state cannot be written
	_, ok := state.(Writer)
	if !ok {
		return nil, fmt.Errorf("'%v' state cannot be written", name)
	}

	// TODO: check whether the state is a LoadableSharedState for optimization.
	// When the state is a LoadableSharedState, we can omit state loading in Write() method.

	s := &sharedStateSink{
		name: name,
	}
	return s, nil
}

func (s *sharedStateSink) Write(ctx *Context, t *Tuple) error {
	state, err := ctx.SharedStates.Get(s.name)
	if err != nil {
		return err
	}

	// It fails if the shared state cannot be written
	writer, ok := state.(Writer)
	if !ok {
		return fmt.Errorf("'%v' state cannot be written", s.name)
	}

	return writer.Write(ctx, t)
}

func (s *sharedStateSink) Close(ctx *Context) error {
	// SharedState must not be terminated when this sink is closed because
	// the state is still being used by other components.
	return nil
}

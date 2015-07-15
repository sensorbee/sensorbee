package core

import (
	"fmt"
)

// writer points a shared state. sharedStateSink will point to the same shared state
// even after the state is removed from the context.
type sharedStateSink struct {
	writer Writer
}

// NewSharedStateSink creates a sink that writes to SharedState.
func NewSharedStateSink(ctx *Context, name string) (Sink, error) {
	// Get SharedState by name
	state, err := ctx.SharedStates.Get(name)
	if err != nil {
		return nil, err
	}

	// It fails if the shared state cannot be written
	writer, ok := state.(Writer)
	if !ok {
		return nil, fmt.Errorf("'%v' state cannot be written", name)
	}

	s := &sharedStateSink{
		writer: writer,
	}
	return s, nil
}

func (s *sharedStateSink) Write(ctx *Context, t *Tuple) error {
	return s.writer.Write(ctx, t)
}

func (s *sharedStateSink) Close(ctx *Context) error {
	return nil
}
}

package core

import (
	"time"
)

// EventType has a type of an event related to Tuple processing.
type EventType int

const (
	// ETInput represents an event where a tuple entered some
	// processing unit (e.g., a Box)
	ETInput EventType = iota
	// ETOutput represents an event where a tuple left some
	// processing unit (e.g., a Box)
	ETOutput
	// ETOther represents any other event
	ETOther
)

// A TraceEvent represents an event in the processing lifecycle of a
// Tuple, in particular transitions from one processing unit to the
// next.
type TraceEvent struct {
	// Timestamp is the time of the event.
	Timestamp time.Time

	// Type represents the type of the event. For transitions, the viewpoint
	// of the Tuple should be assumed. For example, when a Tuple is emitted
	// by a Source, this is an OUTPUT transition; when it enters a Box for
	// processing, this is an INPUT transition. The OTHER Type can be used
	// to add other tracing information.
	Type EventType

	// Msg is any message, but for transitions it makes sense to use the
	// name of the Source/Box/Sink that was left/entered.
	Msg string
}

func (t EventType) String() string {
	switch t {
	case ETInput:
		return "input"
	case ETOutput:
		return "output"
	case ETOther:
		return "other"
	default:
		return "unknown"
	}
}

func tracing(t *Tuple, ctx *Context, inout EventType, msg string) {
	if !ctx.Flags.TupleTrace.Enabled() {
		return
	}
	ev := newDefaultEvent(inout, msg)
	t.AddEvent(ev)
}

func newDefaultEvent(inout EventType, msg string) TraceEvent {
	return TraceEvent{
		time.Now(),
		inout,
		msg,
	}
}

type traceWriter struct {
	w     WriteCloser
	inout EventType
	msg   string
}

func newTraceWriter(w WriteCloser, inout EventType, msg string) *traceWriter {
	return &traceWriter{
		w:     w,
		inout: inout,
		msg:   msg,
	}
}

func (tw *traceWriter) Write(ctx *Context, t *Tuple) error {
	tracing(t, ctx, tw.inout, tw.msg)
	return tw.w.Write(ctx, t)
}

func (tw *traceWriter) Close(ctx *Context) error {
	return tw.w.Close(ctx)
}

package core

import (
	"pfi/sensorbee/sensorbee/tuple"
	"time"
)

func tracing(t *tuple.Tuple, ctx *Context, inout tuple.EventType, msg string) {
	if !ctx.IsTupleTraceEnabled() {
		return
	}
	ev := newDefaultEvent(inout, msg)
	t.AddEvent(ev)
}

func newDefaultEvent(inout tuple.EventType, msg string) tuple.TraceEvent {
	return tuple.TraceEvent{
		time.Now(),
		inout,
		msg,
	}
}

type traceWriter struct {
	w     WriteCloser
	inout tuple.EventType
	msg   string
}

func newTraceWriter(w WriteCloser, inout tuple.EventType, msg string) *traceWriter {
	return &traceWriter{
		w:     w,
		inout: inout,
		msg:   msg,
	}
}

func (tw *traceWriter) Write(ctx *Context, t *tuple.Tuple) error {
	tracing(t, ctx, tw.inout, tw.msg)
	return tw.w.Write(ctx, t)
}

func (tw *traceWriter) Close(ctx *Context) error {
	return tw.w.Close(ctx)
}

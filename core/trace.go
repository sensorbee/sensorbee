package core

import (
	"time"
)

func tracing(t *Tuple, ctx *Context, inout EventType, msg string) {
	if !ctx.IsTupleTraceEnabled() {
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

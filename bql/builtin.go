package bql

import (
	"errors"
	"fmt"
	"io"
	"os"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
)

// TODO: create bql/builtin directory and move components in this file to there

type writerSink struct {
	m           sync.Mutex
	w           io.Writer
	shouldClose bool
}

func (s *writerSink) Write(ctx *core.Context, t *core.Tuple) error {
	js := t.Data.String() // Format this outside the lock

	// This lock is required to avoid interleaving JSONs.
	s.m.Lock()
	defer s.m.Unlock()
	if s.w == nil {
		return errors.New("the sink is already closed")
	}
	_, err := fmt.Fprintln(s.w, js)
	return err
}

func (s *writerSink) Close(ctx *core.Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	if s.w == nil {
		return nil
	}
	if s.shouldClose {
		if c, ok := s.w.(io.Closer); ok {
			return c.Close()
		}
	}
	return nil
}

func createStdouSink(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Sink, error) {
	return &writerSink{
		w: os.Stdout,
	}, nil
}

func createFileSink(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Sink, error) {
	// TODO: currently this sink isn't secure because it accepts any path.
	// TODO: support truncation
	// TODO: support buffering

	var fpath string
	if v, ok := params["path"]; !ok {
		return nil, errors.New("path parameter is missing")
	} else if f, err := data.AsString(v); err != nil {
		return nil, errors.New("path parameter must be a string")
	} else {
		fpath = f
	}

	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &writerSink{
		w:           file,
		shouldClose: true,
	}, nil
}

func init() {
	MustRegisterGlobalSinkCreator("stdout", SinkCreatorFunc(createStdouSink))
	MustRegisterGlobalSinkCreator("file", SinkCreatorFunc(createFileSink))
}

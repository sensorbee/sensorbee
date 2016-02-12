package bql

import (
	"errors"
	"fmt"
	"io"
	"os"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
	"time"
)

// TODO: create bql/builtin directory and move components in this file to there

func createSharedStateSink(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Sink, error) {
	// Get only name parameter from params
	name, ok := params["name"]
	if !ok {
		return nil, fmt.Errorf("cannot find 'name' parameter")
	}
	nameStr, err := data.AsString(name)
	if err != nil {
		return nil, err
	}

	// TODO: Support cascading delete. Because it isn't supported yet,
	// the sink will be running even after the target state is dropped.
	// Moreover, creating a state having the same name after dropping
	// the previous state might result in a confusing behavior.
	return core.NewSharedStateSink(ctx, nameStr)
}

func init() {
	MustRegisterGlobalSinkCreator("uds", SinkCreatorFunc(createSharedStateSink))
}

type writerSink struct {
	m           sync.Mutex
	w           io.Writer
	shouldClose bool
}

func (s *writerSink) Write(ctx *core.Context, t *core.Tuple) error {
	// TODO: support custom formatting. There're several things that need to
	// be considered such as concurrent formatting, zero-copy write, and so on.
	// While encoding tuples outside the lock supports concurrent formatting,
	// it makes it difficult to support zero-copy write.

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

func createStdoutSink(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Sink, error) {
	return &writerSink{
		w: os.Stdout,
	}, nil
}

func createFileSink(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Sink, error) {
	// TODO: currently this sink isn't secure because it accepts any path.
	// TODO: support buffering
	// TODO: provide "format" parameter to support output formats other than "jsonl".
	//       "jsonl" should be the default value.
	// TODO: support "compression" parameter with values like "gz".

	var fpath string
	if v, ok := params["path"]; !ok {
		return nil, errors.New("path parameter is missing")
	} else if f, err := data.AsString(v); err != nil {
		return nil, errors.New("path parameter must be a string")
	} else {
		fpath = f
	}

	flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	if v, ok := params["truncate"]; ok {
		t, err := data.AsBool(v)
		if err != nil {
			return nil, fmt.Errorf("'truncate' parameter must be bool: %v", err)
		}
		if t {
			flags |= os.O_TRUNC
		}
	}

	file, err := os.OpenFile(fpath, flags, 0644)
	if err != nil {
		return nil, err
	}
	return &writerSink{
		w:           file,
		shouldClose: true,
	}, nil
}

func init() {
	MustRegisterGlobalSinkCreator("stdout", SinkCreatorFunc(createStdoutSink))
	MustRegisterGlobalSinkCreator("file", SinkCreatorFunc(createFileSink))
}

func createDroppedTupleCollectorSource(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Source, error) {
	return core.NewDroppedTupleCollectorSource(), nil
}

func init() {
	MustRegisterGlobalSourceCreator("dropped_tuples", SourceCreatorFunc(createDroppedTupleCollectorSource))
}

type nodeStatusSource struct {
	topology core.Topology
	interval time.Duration
	stopCh   chan struct{}
}

func (s *nodeStatusSource) GenerateStream(ctx *core.Context, w core.Writer) error {
	next := time.Now().Add(s.interval)
	for {
		select {
		case <-s.stopCh:
			return nil
		case <-time.After(next.Sub(time.Now())):
		}
		now := time.Now()

		for name, n := range s.topology.Nodes() {
			t := &core.Tuple{
				Timestamp:     now,
				ProcTimestamp: now,
				Data:          n.Status(),
			}
			t.Data["node_name"] = data.String(name)
			t.Data["node_type"] = data.String(n.Type().String())
			w.Write(ctx, t)
		}

		next = next.Add(s.interval)
		if next.Before(now) {
			// delayed too much and should be rescheduled.
			next = now.Add(s.interval)
		}
	}
}

func (s *nodeStatusSource) Stop(ctx *core.Context) error {
	close(s.stopCh)
	return nil
}

// createNodeStatusSourceCreator creates a SourceCreator which creates
// nodeStatusSource. Because it requires core.Topology, it cannot be registered
// statically. It'll be registered in a function like NewTopologyBuilder.
func createNodeStatusSourceCreator(t core.Topology) SourceCreator {
	return SourceCreatorFunc(func(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Source, error) {
		interval := 1 * time.Second
		if v, ok := params["interval"]; !ok {
		} else if d, err := data.ToDuration(v); err != nil {
			return nil, err
		} else {
			interval = d
		}

		return &nodeStatusSource{
			topology: t,
			interval: interval,
			stopCh:   make(chan struct{}),
		}, nil
	})
}

type edgeStatusSource struct {
	topology core.Topology
	interval time.Duration
	stopCh   chan struct{}
}

func (s *edgeStatusSource) GenerateStream(ctx *core.Context, w core.Writer) error {
	next := time.Now().Add(s.interval)

	inputPath := data.MustCompilePath("input_stats.inputs")

	for {
		select {
		case <-s.stopCh:
			return nil
		case <-time.After(next.Sub(time.Now())):
		}
		now := time.Now()

		// collect all nodes that can receive data
		receivers := map[string]core.Node{}
		for name, b := range s.topology.Boxes() {
			receivers[name] = b
		}
		for name, s := range s.topology.Sinks() {
			receivers[name] = s
		}

		// loop over those receiver nodes and consider all of
		// their incoming edges
		for name, n := range receivers {
			nodeStatus := n.Status()
			// get the input status
			inputs, err := nodeStatus.Get(inputPath)
			if err != nil {
				ctx.ErrLog(err).WithField("node_status", nodeStatus).
					WithField("node_name", name).
					Error("No input_stats present in node status")
				continue
			}
			inputMap, err := data.AsMap(inputs)
			if err != nil {
				ctx.ErrLog(err).WithField("inputs", inputs).
					WithField("node_name", name).
					Error("input_stats.inputs is not a Map")
				continue
			}

			// loop over the input nodes to get an edge-centric view
			for inputName, inputStats := range inputMap {
				inputNode, err := s.topology.Node(inputName)
				if err != nil {
					ctx.ErrLog(err).WithField("sender", inputName).
						WithField("receiver", name).
						Error("Node listens to non-existing node")
					continue
				}
				edgeData := data.Map{
					"sender": data.Map{
						"node_name": data.String(inputName),
						"node_type": data.String(inputNode.Type().String()),
					},
					"receiver": data.Map{
						"node_name": data.String(name),
						"node_type": data.String(n.Type().String()),
					},
					// use the input statistics for that edge from the
					// receiver as edge statistics. the data is correct,
					// but the wording may be a bit weird, e.g. "num_received"
					// should maybe rather be "num_transferred"
					"stats": inputStats,
				}
				// write tuple
				t := &core.Tuple{
					Timestamp:     now,
					ProcTimestamp: now,
					Data:          edgeData,
				}
				w.Write(ctx, t)
			}

		}

		next = next.Add(s.interval)
		if next.Before(now) {
			// delayed too much and should be rescheduled.
			next = now.Add(s.interval)
		}
	}
}

func (s *edgeStatusSource) Stop(ctx *core.Context) error {
	close(s.stopCh)
	return nil
}

func createEdgeStatusSourceCreator(t core.Topology) SourceCreator {
	return SourceCreatorFunc(func(ctx *core.Context, ioParams *IOParams, params data.Map) (core.Source, error) {
		interval := 1 * time.Second
		if v, ok := params["interval"]; !ok {
		} else if d, err := data.ToDuration(v); err != nil {
			return nil, err
		} else {
			interval = d
		}

		return &edgeStatusSource{
			topology: t,
			interval: interval,
			stopCh:   make(chan struct{}),
		}, nil
	})
}

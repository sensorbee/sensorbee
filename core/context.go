package core

import (
	"github.com/Sirupsen/logrus"
	"path/filepath"
	"pfi/sensorbee/sensorbee/data"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	temporaryIDCounter int64
)

// NewTemporaryID returns the new temporary 63bit ID. This can be used for
// any purpose.
func NewTemporaryID() int64 {
	return atomic.AddInt64(&temporaryIDCounter, 1)
}

// Context holds a set of functionality that is made available to each Topology
// at runtime. A context is created by the user before creating a Topology.
// Each Context is tied to one Topology and it must not be used by multiple
// topologies.
type Context struct {
	logger       *logrus.Logger
	topologyName string
	Flags        ContextFlags
	SharedStates SharedStateRegistry

	dtMutex   sync.RWMutex
	dtSources map[int64]*droppedTupleCollectorSource
}

// ContextConfig has configuration parameters of a Context.
type ContextConfig struct {
	// Logger provides a logrus's logger used by the Context.
	Logger *logrus.Logger
	Flags  ContextFlags
}

// NewContext creates a new Context based on the config. If config is nil,
// the default config will be used.
func NewContext(config *ContextConfig) *Context {
	if config == nil {
		config = &ContextConfig{}
	}
	logger := config.Logger
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	c := &Context{
		logger:    logger,
		Flags:     config.Flags,
		dtSources: map[int64]*droppedTupleCollectorSource{},
	}
	c.SharedStates = NewDefaultSharedStateRegistry(c)
	return c
}

// Log returns the logger tied to the Context.
func (c *Context) Log() *logrus.Entry {
	return c.log(1)
}

// ErrLog returns the logger tied to the Context having an error information.
func (c *Context) ErrLog(err error) *logrus.Entry {
	return c.log(1).WithField("err", err)
}

func (c *Context) log(depth int) *logrus.Entry {
	// TODO: This is a temporary solution until logrus support filename and line number
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		return c.logger.WithField("topology", c.topologyName)
	}
	file = filepath.Base(file) // only the filename at the moment
	return c.logger.WithFields(logrus.Fields{
		"file":     file,
		"line":     line,
		"topology": c.topologyName,
	})
}

// droppedTuple records tuples dropped by errors.
func (c *Context) droppedTuple(t *Tuple, nodeType NodeType, nodeName string, et EventType, err error) {
	if c.Flags.DroppedTupleLog.Enabled() {
		var js string
		if c.Flags.DroppedTupleSummarization.Enabled() {
			js = data.Summarize(t.Data)
		} else {
			js = t.Data.String()
		}

		l := c.Log().WithFields(nodeLogFields(nodeType, nodeName)).WithFields(logrus.Fields{
			"event_type": et.String(),
			"tuple": logrus.Fields{
				"timestamp": data.Timestamp(t.Timestamp),
				"data":      js,
				// TODO: Add trace
			},
		})
		if err != nil {
			l = l.WithField("err", err)
		}
		l.Info("A tuple was dropped from the topology") // TODO: debug should be better?
	}

	c.dtMutex.RLock()
	defer c.dtMutex.RUnlock()
	if len(c.dtSources) == 0 {
		return
	}
	// TODO: reduce copies
	dt := t.Copy()
	dt.Data = data.Map{
		"node_type":  data.String(nodeType.String()),
		"node_name":  data.String(nodeName),
		"event_type": data.String(et.String()),
		"data":       dt.Data,
	}
	if err != nil {
		dt.Data["error"] = data.String(err.Error())
	}
	shouldCopy := len(c.dtSources) > 1
	for _, s := range c.dtSources {
		// TODO: reduce copies
		copied := dt
		if shouldCopy {
			copied = dt.Copy()
		}
		s.w.Write(c, copied) // There isn't much meaning to report errors here.
	}
}

// addDroppedTupleSource is a listener which receives dropped tuples. The
// return value is the ID of the listener and it'll be required for
// removeDroppedTupleListener.
func (c *Context) addDroppedTupleSource(s *droppedTupleCollectorSource) int64 {
	c.dtMutex.Lock()
	defer c.dtMutex.Unlock()
	id := NewTemporaryID()
	c.dtSources[id] = s
	return id
}

func (c *Context) removeDroppedTupleSource(id int64) {
	c.dtMutex.Lock()
	defer c.dtMutex.Unlock()
	delete(c.dtSources, id)
}

// AtomicFlag is a boolean flag which can be read/written atomically.
type AtomicFlag int32

// Set sets a boolean value to the flag.
func (a *AtomicFlag) Set(b bool) {
	var i int32
	if b {
		i = 1
	}
	atomic.StoreInt32((*int32)(a), i)
}

// Enabled returns true if the flag is enabled.
func (a *AtomicFlag) Enabled() bool {
	return atomic.LoadInt32((*int32)(a)) != 0
}

// ContextFlags is an arrangement of SensorBee processing settings.
type ContextFlags struct {
	// TupleTrace is a Tuple's tracing on/off flag. If the flag is 0
	// (means false), a topology does not trace Tuple's events.
	// The flag can be set when creating a Context, or when the topology
	// is running. In the latter case, Context.TupleTrace.Set() should
	// be used for thread safety.
	// There is a delay between setting the flag and start/stop to trace Tuples.
	TupleTrace AtomicFlag

	// DroppedTupleLog is a flag which turns on/off logging of dropped tuple
	// events.
	DroppedTupleLog AtomicFlag

	// DroppedTupleSummarization is a flag to trun on/off summarization of
	// dropped tuple logging. If this flag is enabled, tuples being logged will
	// be a little smaller than the originals. However, they might not be parsed
	// as JSONs. If the flag is disabled, output JSONs can be parsed.
	DroppedTupleSummarization AtomicFlag
}

type droppedTupleCollectorSource struct {
	w  Writer
	id int64
	wg sync.WaitGroup
}

// NewDroppedTupleCollectorSource returns a source which generates a stream
// containing tuples dropped by other nodes. Tuples generated from this source
// won't be reported again even if they're dropped later on. So, when a sink
// is connected to two boxes and one of them is connected to the source,
// tuples dropped by the sink won't be reported even if those tuples are sent
// from another box which isn't connected to the source. Therefore, it's safe
// to have a DAG having the source and isolate it from the regular processing
// flow.
//
// Tuples generated from this source has the following fields in Data:
//
//	- node_type: the type of the node which dropped the tuple
//	- node_name: the name of the node which dropped the tuple
//	- event_type: the type of the event indicating when the tuple was dropped
//	- error(optional): the error information if any
//	- data: the original content in which the dropped tuple had
func NewDroppedTupleCollectorSource() Source {
	return &droppedTupleCollectorSource{}
}

func (s *droppedTupleCollectorSource) GenerateStream(ctx *Context, w Writer) error {
	s.w = w
	s.id = ctx.addDroppedTupleSource(s)
	s.wg.Add(1)
	s.wg.Wait()
	return nil
}

func (s *droppedTupleCollectorSource) Stop(ctx *Context) error {
	ctx.removeDroppedTupleSource(s.id)
	s.wg.Done()
	return nil
}

package core

import (
	"github.com/Sirupsen/logrus"
	"path/filepath"
	"pfi/sensorbee/sensorbee/data"
	"runtime"
	"sync/atomic"
)

// Context holds a set of functionality that is made available to each Topology
// at runtime. A context is created by the user before creating a Topology.
// Each Context is tied to one Topology and it must not be used by multiple
// topologies.
type Context struct {
	logger       *logrus.Logger
	topologyName string
	Flags        ContextFlags
	SharedStates SharedStateRegistry
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
		logger: logger,
		Flags:  config.Flags,
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
		l := c.Log().WithFields(nodeLogFields(nodeType, nodeName)).WithFields(logrus.Fields{
			"event_type": et.String(),
			"tuple": logrus.Fields{
				"timestamp": data.Timestamp(t.Timestamp),
				"data":      t.Data,
				// TODO: Add trace
			},
		})
		if err != nil {
			l = l.WithField("err", err)
		}
		l.Info("A tuple was dropped from the topology") // TODO: debug should be better?
	}

	// TODO: add listener here to notify events
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
}

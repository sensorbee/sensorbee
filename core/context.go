package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync/atomic"
)

// Context holds a set of functionality that is made available to
// each Box at runtime. A context is created by the user before
// processing starts according to specific application needs
// (e.g., central log collection) and passed in to StaticTopology.Run().
type Context struct {
	Logger LogManager
	Config Configuration
}

// IsTupleTraceEnabled can get from Context Configuration, whether the topology
// is enabled to trace Tuples' events.
func (c *Context) IsTupleTraceEnabled() bool {
	if atomic.LoadInt32(&c.Config.TupleTraceEnabled) != 0 {
		return true
	}
	return false
}

// SetTupleTraceEnabled can switch th Configuration of tracing tuples events.
// If the argument bool flag is same as Context's Configuration, does nothing.
func (c *Context) SetTupleTraceEnabled(b bool) {
	var i int32 = 0
	if b {
		i = int32(1)
	}
	atomic.StoreInt32(&c.Config.TupleTraceEnabled, i)
}

type LogLevel int

// Constants that can be used as log levels
const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// String returns a string representation of a LogLevel.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "unknown"
	}
}

// A LogManager corresponds to a typical Logger class found
// in many applications. It can be accessed from within a Box
// via the Context object.
//
// Log logs the given message respecting the given level. The
// string will be processed using the fmt.Sprintf function.
//
// DroppedTuple should be called by a Box for each Tuple
// where an error occurred that was so bad that processing
// cannot continue for that particular Tuple (e.g., schema
// mismatch), but not grave enough to abort the whole process.
type LogManager interface {
	Log(level LogLevel, msg string, a ...interface{})
	DroppedTuple(t *tuple.Tuple, msg string, a ...interface{})
}

// Configuration is an arrangement of SensorBee processing.
//
// TupleTraceEnabled is a Tuple's tracing on/off flag. If the flag is 0
// (means false), a topology does not trace Tuple's events.
// The flag can be set when creating a Context, or running the topology.
// There is a delay between setting the flag and start/stop to trace Tuples.
type Configuration struct {
	TupleTraceEnabled int32
}

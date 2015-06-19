package core

import (
	"pfi/sensorbee/sensorbee/tuple"
	"sync/atomic"
)

// Context holds a set of functionality that is made available to
// each Box at runtime. A context is created by the user before
// processing starts according to specific application needs
// (e.g., central log collection) and passed in to StaticTopology.Run().
type Context struct {
	Logger       LogManager
	Config       Configuration
	SharedStates SharedStateRegistry
}

// IsTupleTraceEnabled can get from Context Configuration, whether the topology
// is enabled to trace Tuples' events.
func (c *Context) IsTupleTraceEnabled() bool {
	if atomic.LoadInt32(&c.Config.TupleTraceEnabled) != 0 {
		return true
	}
	return false
}

// SetTupleTraceEnabled can switch the setting of tracing tuples events.
// If the argument flag is the same as Context's Configuration, does nothing.
func (c *Context) SetTupleTraceEnabled(b bool) {
	var i int32
	if b {
		i = 1
	}
	atomic.StoreInt32(&c.Config.TupleTraceEnabled, i)
}

// GetSharedState returns a state registered to this Context.
func (c *Context) GetSharedState(name string) (SharedState, error) {
	return c.SharedStates.Get(c, name)
}

// LogLevel represents severity of an event being logged.
type LogLevel int

// Constants that can be used as log levels
const (
	// Debug is a LogLevel which is used when the information is written
	// for debug purposes and shouldn't be written in a production environment.
	Debug LogLevel = iota

	// Info is a LogLevel used for events which might be useful to report
	// but no action is required.
	Info

	// Warning is a LogLevel used for events which don't require immediate
	// actions but some diagnoses are required later.
	Warning

	// Error is a LogLevel used for events which require immediate actions.
	Error

	// Fatal is a LogLevel used when the process can no longer perform
	// any operation correctly and should stop.
	Fatal
)

// String returns a string representation of a LogLevel.
func (l LogLevel) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	case Fatal:
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

// Configuration is an arrangement of SensorBee processing settings.
//
// TupleTraceEnabled is a Tuple's tracing on/off flag. If the flag is 0
// (means false), a topology does not trace Tuple's events.
// The flag can be set when creating a Context, or when the topology
// is running. In the latter case, Context.SetTupleTraceEnabled() should
// be used for thread safety.
// There is a delay between setting the flag and start/stop to trace Tuples.
type Configuration struct {
	TupleTraceEnabled int32
}

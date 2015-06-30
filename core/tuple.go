package core

import (
	"pfi/sensorbee/sensorbee/data"
	"time"
)

// Tuple is a fundamental data structure in SensorBee. All data
// that is processed is stored in tuples.
type Tuple struct {
	// Data is the actual data that is processed.
	Data data.Map

	// InputName can be used to identify the sender of a tuple when a
	// Box processes data from multiple inputs. It will be set before
	// Box.Process is called. Also see BoxDeclarer.NamedInput.
	InputName string

	// Timestamp is the time when this tuple was originally generated,
	// e.g., the timestamp of a camera image or a sensor-emitted value.
	// It is an integral part of the data and should not be changed.
	// It should be saved along with the Data when persisted to a data
	// store so that timestamp-based reprocessing can be done with the
	// same results at a later point in time.
	Timestamp time.Time

	// ProcTimestamp is the time when this tuple entered the topology
	// for processing. It should be set by the Source that emitted this
	// Tuple.
	ProcTimestamp time.Time

	// BatchID is reserved for future use.
	BatchID int64

	// Trace is used during debugging to trace to way of a Tuple through
	// a topology. See the documentation for TraceEvent.
	Trace []TraceEvent
}

// AddEvent adds a TraceEvent to this Tuple's trace. This is not
// thread-safe because it is assumed a Tuple is only processed
// by one unit at a time.
func (t *Tuple) AddEvent(ev TraceEvent) {
	t.Trace = append(t.Trace, ev)
}

// Copy creates a deep copy of a Tuple, including the contained
// data. This can be used, e.g., by fan-out pipes.
func (t *Tuple) Copy() *Tuple {
	// except for Data, there are only value types in
	// Tuple, so we can use normal copy for everything
	// except Data
	out := *t
	out.Data = out.Data.Copy()

	// the copied tuple should have new event history,
	// which is isolated from the original tuple,
	// past events are copied from the original tuple
	tr := make([]TraceEvent, len(t.Trace))
	copy(tr, t.Trace)
	out.Trace = tr

	return &out
}

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
		return "INPUT"
	case ETOutput:
		return "OUTPUT"
	case ETOther:
		return "OTHER"
	default:
		return "unknown"
	}
}

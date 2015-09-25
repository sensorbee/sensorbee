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

	// Flags has bit flags which controls behavior of this tuple. When a Box
	// emits a tuple derived from a received one, it must copy this field
	// otherwise a problem like infinite reporting of a dropped tuple could
	// occur.
	Flags TupleFlags

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
// data. This can be used, e.g., by fan-out pipes. When Tuple.Data doesn't
// need to be cloned, just use newTuple := *oldTuple instead of this method.
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

// NewTuple creates and initializes a Tuple with default
// values. And it copied Data by argument. Timestamp and
// ProcTimestamp fields will be set time.Now() value.
func NewTuple(d data.Map) *Tuple {
	now := time.Now()
	return &Tuple{
		Data:          d.Copy(),
		Timestamp:     now,
		ProcTimestamp: now,
	}
}

// TupleFlags has flags which controls behavior of a tuple.
type TupleFlags uint32

const (
	// TFDropped is a flag which is set when a tuple is dropped. Once this flag
	// is set to a tuple, the tuple will not be reported when it is dropped.
	TFDropped TupleFlags = 1 << iota
)

// Set sets a set of flags at once.
func (f *TupleFlags) Set(v TupleFlags) {
	*f |= v
}

// IsSet returns true if the all given flags are set.
func (f *TupleFlags) IsSet(v TupleFlags) bool {
	return *f&v == v
}

// Clear clears a set of flags at once.
func (f *TupleFlags) Clear(v TupleFlags) {
	*f &= ^v
}

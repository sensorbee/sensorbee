package core

import (
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"sync/atomic"
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
// need to be cloned, call ShallowCopy. NEVER do newTuple := *oldTuple.
func (t *Tuple) Copy() *Tuple {
	// except for Data, there are only value types in
	// Tuple, so we can use normal copy for everything
	// except Data
	out := t.shallowCopy()
	out.Data = out.Data.Copy()
	out.Flags.Clear(TFSharedData) // not shared
	return out
}

// ShallowCopy creates a new copy of a tuple. It only deep copies trace
// information. Because Data is shared between the old tuple and the new tuple,
// TFSharedData is set by this method. However, the tuple returned from this
// method isn't shared and its TFShared flag isn't set.
//
// It's safe to clear TFShared flag after assigning a new data.Map to Data
// field:
//
//	newT := oldT.ShallowCopy()
//	newT.Data = data.Map{} // no field other than Data is shared with oldT.
//	newT.Flags.Clear(TFShared) // therefore, it's safe to clear the flag here.
func (t *Tuple) ShallowCopy() *Tuple {
	out := t.shallowCopy()
	out.Flags.Set(TFSharedData)
	t.Flags.Set(TFSharedData)
	return out
}

func (t *Tuple) shallowCopy() *Tuple {
	out := *t
	out.Flags.Clear(TFShared)

	// the copied tuple should have new event history,
	// which is isolated from the original tuple,
	// past events are copied from the original tuple
	if t.Trace != nil {
		tr := make([]TraceEvent, len(t.Trace))
		copy(tr, t.Trace)
		out.Trace = tr
	}
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

	// TFShared is a flag which is set when a Tuple is shared by multiple node
	// so that a node changing the tuple must make a copy of it rather than
	// modifying it directly. This flag only indicates that Tuple struct itself
	// is shared. Tuple.Data might be shared even if this flag isn't set.
	//
	// To update a field of a tuple with TFShared flag, use ShallowCopy() or
	// Copy(). Also, if you keep a reference to a tuple or (parts of) its Data
	// after processing it, you must set the TFShared flag for that tuple.
	TFShared

	// TFSharedData is a flag which is set when Tuple.Data is shared by other
	// tuples. Tuple.Data must not directly modified if the flag is set.
	// To update Data of a tuple with TFSharedData flag, use Copy().
	//
	// Relations of TFShared and TFSharedData are summarized below:
	//
	//	(TFShared is set, TFSharedData is set):
	//	(true, true): *Tuple is copied
	//	(true, false): never happens
	//	(false, true): a tuple returned from ShallowCopy
	//	(false, false): a tuple returned from NewTuple or Copy
	TFSharedData
)

// Set sets a set of flags at once.
func (f *TupleFlags) Set(v TupleFlags) {
	newFlag := uint32(v)
	if v == TFShared {
		// TFSharedData needs to be set as well because Tuple.Data is shared, too.
		newFlag |= uint32(TFSharedData)
	}

	for {
		old := atomic.LoadUint32((*uint32)(f))
		if atomic.CompareAndSwapUint32((*uint32)(f), old, old|uint32(v)) {
			break
		}
	}
}

// IsSet returns true if the all given flags are set.
func (f *TupleFlags) IsSet(v TupleFlags) bool {
	return TupleFlags(atomic.LoadUint32((*uint32)(f)))&v == v
}

// Clear clears a set of flags at once.
func (f *TupleFlags) Clear(v TupleFlags) {
	for {
		old := atomic.LoadUint32((*uint32)(f))
		if atomic.CompareAndSwapUint32((*uint32)(f), old, old & ^uint32(v)) {
			break
		}
	}
}

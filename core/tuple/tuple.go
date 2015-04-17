package tuple

import (
	"time"
)

type Tuple struct {
	Data Map

	Timestamp     time.Time
	ProcTimestamp time.Time
	BatchID       int64

	Trace []TraceEvent
}

func (t *Tuple) AddEvent(ev TraceEvent) {
	t.Trace = append(t.Trace, ev)
}

func (t *Tuple) Copy() *Tuple {
	// except for Data, there are only value types in
	// Tuple, so we can use normal copy for everything
	// except Data
	out := *t
	out.Data = out.Data.Copy()
	return &out
}

type InOutType int

const (
	INPUT InOutType = iota
	OUTPUT
)

type TraceEvent struct {
	Timestanp time.Time
	Inout     InOutType
	Msg       string
}

func (t InOutType) String() string {
	switch t {
	case INPUT:
		return "INPUT"
	case OUTPUT:
		return "OUTPUT"
	default:
		return "unknown"
	}
}

package tuple

import (
	"time"
)

type Tuple struct {
	Data Map

	Timestamp     time.Time
	ProcTimestamp time.Time
	BatchID       int64
}

func (t *Tuple) Copy() *Tuple {
	// except for Data, there are only value types in
	// Tuple, so we can use normal copy for everything
	// except Data
	out := *t
	out.Data = out.Data.Copy()
	return &out
}

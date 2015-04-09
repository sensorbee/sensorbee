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

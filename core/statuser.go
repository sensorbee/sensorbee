package core

import (
	"pfi/sensorbee/sensorbee/data"
)

// Statuser is an interface which provides Status method to retrieve its status.
type Statuser interface {
	// Status returns the status of the component. Formats of statuses are not
	// strictly defined and each component can provide any information it has.
	// Node.Status has a document to show what information it returns.
	Status() data.Map
}

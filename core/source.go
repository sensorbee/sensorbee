package core

// A Source describes an entity that inserts data into a topology
// from the outside world (e.g., a fluentd instance).
//
// GenerateStream will start creating tuples and writing them to
// the given Writer in a blocking way. (Start as a gorouting to start
// the process in the background.) It will return when all tuples
// have been written (in the case of a finite data source) or if
// there was a severe error.
//
// Schema will return the schema of the data that can be expected
// from this Source. It is nil if no guarantees about that schema
// are made.
type Source interface {
	GenerateStream(ctx *Context, w Writer) error
	Schema() *Schema
}

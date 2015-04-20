package core

// A Schema describes the shape of the data associated to a Tuple.
// Is is used for two different purposes. The first is to describe
// the data that will be emitted by a Source or a Box.
// The second is to describe the required shape that data must
// have so that it can be used as input to a Box.
type Schema string

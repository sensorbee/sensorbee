package bql

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
)

var (
	// sourceTypes keeps track of all known source types.
	sourceTypes = map[string]SourceCreator{}
	// sinkTypes keeps track of all known sink types.
	sinkTypes = map[string]SinkCreator{}
)

// SourceCreator is a function that creates a source from a
// given parameter map as extracted from a CREATE SOURCE ... WITH
// or CREATE STREAM FROM ... SOURCE WITH statement.
type SourceCreator func(map[string]string) (core.Source, error)

// SinkCreator is a function that creates a sink from a
// given parameter map as extracted from a CREATE SINK ... WITH
// statement.
type SinkCreator func(map[string]string) (core.Sink, error)

// RegisterSourceType registers a new type of Source for use
// in BQL and stores a function to create an instance of that
// Source type. It will return an error if there is already
// a Source type registered with that name.
//
// This should be called in the `init` function of any module
// that provides a Source type.
func RegisterSourceType(typeName string, creator SourceCreator) error {
	if _, exists := sourceTypes[typeName]; exists {
		return fmt.Errorf("source type '%s' already exists", typeName)
	}
	sourceTypes[typeName] = creator
	return nil
}

// LookupSourceType checks if there is a creator function known
// for Sources of the given type and returns it, if present. If the
// Source type is not known, the second return value is false.
func LookupSourceType(typeName parser.SourceSinkType) (SourceCreator, bool) {
	src, ok := sourceTypes[string(typeName)]
	return src, ok
}

// RegisterSinkType registers a new type of Sink for use
// in BQL and stores a function to create an instance of that
// Sink type. It will return an error if there is already
// a Sink type registered with that name.
//
// This should be called in the `init` function of any module
// that provides a Sink type.
func RegisterSinkType(typeName string, creator SinkCreator) error {
	if _, exists := sinkTypes[typeName]; exists {
		return fmt.Errorf("sink type '%s' already exists", typeName)
	}
	sinkTypes[typeName] = creator
	return nil
}

// LookupSinkType checks if there is a creator function known
// for Sinks of the given type and returns it, if present. If the
// Sink type is not known, the second return value is false.
func LookupSinkType(typeName parser.SourceSinkType) (SinkCreator, bool) {
	sink, ok := sinkTypes[string(typeName)]
	return sink, ok
}

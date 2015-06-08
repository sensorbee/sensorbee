package core

// A StaticTopology is returned by StaticTopologyBuilder.Build and is
// used to control the actual processing.
type StaticTopology interface {
	// Run starts processing in a blocking way. It should return when
	// all Source.GenerateStream methods of this topology have finished.
	Run(ctx *Context) error

	// Stop gracefully stops processing in a blocking way. It will stop all
	// sources, finish the processing of all tuples that are still in
	// the system, and then return.
	Stop(ctx *Context) error
}

// A StaticTopologyBuilder can be used to assemble a processing
// pipeline consisting of sources that import tuples from some
// location, boxes that process tuples, and sinks that export
// tuples to some location. Sources, boxes and sinks are
// identified with a name which must be unique within a Topology,
// even across types.
type StaticTopologyBuilder interface {
	// AddSource will add a Source with the given name to the
	// Topology under construction.
	//
	// Example:
	//  tb.AddSource("fluentd-in", source)
	//
	// Don't add the same instance of a Source more than once.
	// Some implementation of StaticTopologyBuilder may reject
	// duplicated registrations of an instance.
	AddSource(name string, source Source) SourceDeclarer

	// AddBox will add a Box with the given name to the Topology
	// under construction. Input to this Box can be specified using
	// the Input or NamedInput method of the returned BoxDeclarer.
	//
	// Example:
	//  tb.AddBox("filter", box).Input("fluentd-in")
	//
	// Don't add the same instance of a Box more than once even if they have
	// different names. However, if a Box has a reference counter and
	// its initialization and termination are done exactly once at proper
	// timing, it can be added multiple times when a builder supports duplicated
	// registration of the same instance of a Box.
	AddBox(name string, box Box) BoxDeclarer

	// AddSink will add a Sink with the given name to the Topology
	// under construction. Input to this Sink can be specified using
	// the Input method of the returned SinkDeclarer.
	//
	// Example:
	//  tb.AddSink("fluentd-out", sink).Input("filter")
	//
	// Don't add the same instance of a Sink more than once.
	// Some implementation of StaticTopologyBuilder may reject
	// duplicated registrations of an instance.
	AddSink(name string, sink Sink) SinkDeclarer

	// Build will assemble a processing pipeline from the given
	// sources, boxes and sinks and return a StaticTopology object
	// that can be used to start the actual process.
	Build() (StaticTopology, error)
}

// DeclarerError provides for the various Declarer interfaces
// to be chained even though they may result in an error,
// providing a monadic interface. When building a topology,
// each operation should be checked for possible errors.
//
// Example:
//   var d DeclarerError
//   d = tb.AddSource("fluentd-in", source)
//   if d.Err() != nil { ... }
//   d = tb.AddBox("filter", box).Input("fluentd-in")
//   if d.Err() != nil { ... }
type DeclarerError interface {
	Err() error
}

// A SourceDeclarer is returned by the AddSource method of
// StaticTopologyBuilder.
type SourceDeclarer interface {
	// Err (inherited from DeclarerError) allows to check if anything
	// went wrong while adding the source in focus to the topology.
	DeclarerError
}

// A BoxDeclarer is returned by the AddBox method of StaticTopologyBuilder
// to allow for a fluid interface (method chaining).
type BoxDeclarer interface {
	// Input can be used to connect the output of the Source or Box
	// with the given name to the box in focus as an unnamed stream.
	Input(name string) BoxDeclarer

	// NamedInput can be used to connect the output of the Source or
	// Box with the given name to the box in focus as a named stream.
	// See the documentation of StaticTopologyBuilder and BoxInputConstraints
	// for details.
	NamedInput(name string, inputName string) BoxDeclarer

	// Err (inherited from DeclarerError) allows to check if anything
	// went wrong while adding the box in focus to the topology or
	// connecting inputs. Any error that happened should be transported
	// outside, so that checking an error of a chained method call
	// returns the innermost error.
	DeclarerError
}

// A SinkDeclarer is returned by the AddSink method of
// StaticTopologyBuilder to allow for a fluid interface (method
// chaining).
type SinkDeclarer interface {
	// Input can be used to connect the output of the Source or
	// Box with the given name to the sink in focus.
	Input(name string) SinkDeclarer

	// Err (inherited from DeclarerError) allows to check if anything
	// went wrong while adding the sink in focus to the topology or
	// connecting inputs. Any error that happened should be transported
	// outside, so that checking an error of a chained method call
	// returns the innermost error.
	DeclarerError
}

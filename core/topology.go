package core

// A Topology is returned by StaticTopologyBuilder.Build and is
// used to control the actual processing.
//
// Run starts processing in a blocking way. It should return when all
// Source.GenerateStream methods of this topology have finished.
type Topology interface {
	Run(ctx *Context)
}

// A StaticTopologyBuilder can be used to assemble a processing
// pipeline consisting of sources that import tuples from some
// location, boxes that process tuples, and sinks that export
// tuples to some location. Sources, boxes and sinks are
// identified with a name which must be unique within a Topology,
// even across types.
//
// AddSource will add a Source with the given name to the
// Topology under construction.
//
// Example:
//  tb.AddSource("fluentd-in", source)
//
// AddBox will add a Box with the given name to the Topology
// under construction. Input to this Box can be specified using
// the Input or NamedInput method of the returned BoxDeclarer.
//
// Example:
//  tb.AddBox("filter", box).Input("fluentd-in")
//
// AddSink will add a Sink with the given name to the Topology
// under construction. Input to this Sink can be specified using
// the Input method of the returned SinkDeclarer.
//
// Example:
//  tb.AddSink("fluentd-out", sink).Input("filter")
//
// Build will assemble a processing pipeline from the given
// sources, boxes and sinks and return a Topology object that
// can be used to start the actual process. Note that it is
// currently not safe to call Build multiple times.
type StaticTopologyBuilder interface {
	AddSource(name string, source Source) SourceDeclarer
	AddBox(name string, box Box) BoxDeclarer
	AddSink(name string, sink Sink) SinkDeclarer

	Build() Topology
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
// TopologyBuilder.
//
// Err (inherited from DeclarerError) allows to check if anything
// went wrong while adding the source in focus to the topology.
type SourceDeclarer interface {
	DeclarerError
}

// A BoxDeclarer is returned by the AddBox method of TopologyBuilder
// to allow for a fluid interface (method chaining).
//
// Input can be used to connect the output of the Source or Box
// with the given name to the box in focus as an unnamed stream.
//
// NamedInput can be used to connect the output of the Source or
// Box with the given name to the box in focus as a named stream.
// See the documentation of StaticTopologyBuilder and BoxInputConstraints
// for details.
//
// Err (inherited from DeclarerError) allows to check if anything
// went wrong while adding the box in focus to the topology or
// connecting inputs. Any error that happened should be transported
// outside, so that checking an error of a chained method call
// returns the innermost error.
type BoxDeclarer interface {
	Input(name string) BoxDeclarer
	NamedInput(name string, inputName string) BoxDeclarer
	DeclarerError
}

// A SinkDeclarer is returned by the AddSink method of TopologyBuilder
// to allow for a fluid interface (method chaining).
//
// Input can be used to connect the output of the Source or
// Box with the given name to the sink in focus.
//
// Err (inherited from DeclarerError) allows to check if anything
// went wrong while adding the sink in focus to the topology or
// connecting inputs. Any error that happened should be transported
// outside, so that checking an error of a chained method call
// returns the innermost error.
type SinkDeclarer interface {
	Input(name string) SinkDeclarer
	DeclarerError
}

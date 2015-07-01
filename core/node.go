package core

import (
	"fmt"
	"regexp"
)

// NodeType represents the type of a node in a topology.
type NodeType int

const (
	// NTSource means the node is a Source.
	NTSource NodeType = iota

	// NTBox means the node is a Box.
	NTBox

	// NTSink means the node is a Sink.
	NTSink
)

func (t NodeType) String() string {
	switch t {
	case NTSource:
		return "source"
	case NTBox:
		return "box"
	case NTSink:
		return "sink"
	default:
		return "unknown"
	}
}

var (
	nodeNameRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")
)

// ValidateNodeName validates if the node name has valid format and length. The
// minimum length of a node name is 1 and the maximum is 127. The format has to
// be [a-zA-Z][a-zA-Z0-9_]*.
func ValidateNodeName(name string) error {
	if l := len(name); l < 1 {
		return fmt.Errorf("node name is empty")
	} else if l > 127 {
		return fmt.Errorf("node name can be at most 127 letters: %v", len(name))
	}

	if !nodeNameRegexp.MatchString(name) {
		return fmt.Errorf("node name doesn't follow the format [a-zA-Z][a-zA-Z0-9_]*: %v", name)
	}
	return nil
}

// Node is a node registered to a topology. It defines methods
// common to Source, Box, and Sink nodes.
type Node interface {
	// Type returns the type of the node, which can be NTSource, NTBox,
	// or NTSink. It's safe to convert Node to a specific node interface
	// corresponding to the returned NodeType. For example, if NTBox is
	// returned, the node can be converted to BoxNode with a type
	// assertion.
	Type() NodeType

	// Name returns the name of the node in the registered topology.
	Name() string

	// State returns the current state of the node.
	State() TopologyStateHolder

	// Stop stops the node. When the node is a source, Stop waits until the
	// source actually stops generating tuples. When the node is a box a sink,
	// it waits until the box or the sink is terminated.
	//
	// The node will not be removed from the topology after it stopped.
	Stop() error

	// Meta returns meta information of the node. The meta information can be
	// updated by changing the return value. However, the meta information is
	// not protected from concurrent writes and the caller has to care about it.
	Meta() map[string]interface{}
}

// SourceNode is a Source registered to a topology.
type SourceNode interface {
	Node

	// Source returns internal source passed to Topology.AddSource.
	Source() Source

	// Pause pauses a running source. A paused source can be resumed by calling
	// Resume method. Pause is idempotent.
	Pause() error

	// Resume resumes a paused source. Resume is idempotent.
	Resume() error
}

// BoxNode is a Box registered to a topology.
type BoxNode interface {
	Node

	// Box returns internal source passed to Topology.AddBox.
	Box() Box

	// Input adds a new input from a Source, another Box, or even the Box
	// itself. refname refers a name of node from which the Box want to receive
	// tuples. There must be a Source or a Box having the name.
	Input(refname string, config *BoxInputConfig) error

	// EnableGracefulStop activates a graceful stop mode. If it is enabled,
	// Stop method waits until the Box doesn't have an incoming tuple. The Box
	// doesn't wait until, for example, a source generates all tuples. It only
	// waits for the moment when the Box's input queue gets empty and stops
	// even if some inputs are about to send a new tuple to the Box.
	EnableGracefulStop()

	// StopOnDisconnect tells the Box that it may automatically stop when all
	// inbound or outbound connections (channels or pipes) are closed. After
	// calling this method, the Box can automatically stop even if Stop method
	// isn't explicitly called.
	//
	// connDir can be InboundConnection, OutboundConnnection, or bitwise-or of
	// them. When both InboundConnection and OutboundConnection is specified,
	// the Box stops if all inbound connections are closed OR all outbound
	// connections are closed. For example, when a Box has two inbound
	// connections and one outbound connections, it stops if the outbound
	// connections is closed while two inbound connections are active.
	//
	// Currently, there's no way to disable StopOnDisconnect once it's enabled.
	// Also, it simply overwrites the direction flag as follows, so Inbound and
	// Outbound can be set separately:
	//
	//	boxNode.StopOnDisconnect(core.Inbound | core.Outbound)
	//	boxNode.StopOnDisconnect(core.Outbound) // core.Inbound is still enabled.
	StopOnDisconnect(dir ConnDir)
}

type ConnDir int

const (
	Inbound ConnDir = 1 << iota
	Outbound
)

// BoxInputConfig has parameters to customize input behavior of a Box on each
// input pipe.
type BoxInputConfig struct {
	// InputName is a custom name attached to incoming tuples. When it is empty,
	// "*" will be used.
	InputName string

	// Capacity is the maximum capacity (length) of input pipe. When this
	// parameter is 0, the default value is used. This parameter is only used
	// as a hint and doesn't guarantee that the pipe can actually have the
	// specified number of tuples.
	Capacity int
}

func (c *BoxInputConfig) inputName() string {
	if c.InputName == "" {
		return "*"
	}
	return c.InputName
}

func (c *BoxInputConfig) capacity() int {
	if c.Capacity == 0 {
		return 1024
	}
	return c.Capacity
}

var defaultBoxInputConfig = &BoxInputConfig{}

// SinkNode is a Sink registered to a topology.
type SinkNode interface {
	Node

	// Sink returns internal source passed to Topology.AddSink.
	Sink() Sink

	// Input adds a new input from a Source or a Box. refname refers a name of
	// node from which the Box want to receive tuples. There must be a Source
	// or a Box having the name.
	Input(refname string, config *SinkInputConfig) error

	// EnableGracefulStop activates a graceful stop mode. If it is enabled,
	// Stop method waits until the Sink doesn't have an incoming tuple. The Sink
	// doesn't wait until, for example, a source generates all tuples. It only
	// waits for the moment when the Sink's input queue gets empty and stops
	// even if some inputs are about to send a new tuple to the Sink.
	EnableGracefulStop()

	// StopOnDisconnect tells the Sink that it may automatically stop when all
	// incoming connections (channels or pipes) are closed. After calling this
	// method, the Sink can automatically stop even if Stop method isn't
	// explicitly called.
	StopOnDisconnect()
}

// SinkInputConfig has parameters to customize input behavior of a Sink on
// each input pipe.
type SinkInputConfig struct {
	// Capacity is the maximum capacity (length) of input pipe. When this
	// parameter is 0, the default value is used. This parameter is only used
	// as a hint and doesn't guarantee that the pipe can actually have the
	// specified number of tuples.
	Capacity int
}

func (c *SinkInputConfig) capacity() int {
	if c.Capacity == 0 {
		return 1024
	}
	return c.Capacity
}

var defaultSinkInputConfig = &SinkInputConfig{}

// ResumableNode is a node in a topology which can dynamically be paused and
// resumed at runtime.
type ResumableNode interface {
	// Pause pauses a running node. A paused node can be resumed by calling
	// Resume method. Pause is idempotent and pausing a paused node shouldn't
	// fail. Pause may be called before a node runs. For example, when a node
	// is a source, Pause could be called before calling GenerateStream. In
	// that case, GenerateStream should not generate any tuple until Resume is
	// called.
	//
	// When Stop is called while the node is paused, the node must stop without
	// waiting for Resume.
	Pause(ctx *Context) error

	// Resume resumes a paused node. Resume is idempotent and resuming a running
	// node shouldn't fail. Resume may be called before a node runs.
	Resume(ctx *Context) error
}

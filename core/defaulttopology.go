package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
	"time"
)

type DefaultTopology struct {
	// tb.boxes may contain multiple instances of the same Box object.
	// Use a set-like map to avoid calling Init() twice on the same object.
	boxpointers map[*Box]bool

	sources map[string]Source
	pipes   map[string]*SequentialPipe
}

func (t *DefaultTopology) Run(ctx *Context) {
	for box, _ := range t.boxpointers {
		(*box).Init(ctx)
	}

	var wg sync.WaitGroup
	for name, source := range t.sources {
		wg.Add(1)
		go func(name string, source Source) {
			source.GenerateStream(t.pipes[name])
			wg.Done()
		}(name, source)
	}
	wg.Wait()
}

/**************************************************/

type DefaultStaticTopologyBuilder struct {
	sources     map[string]Source
	boxes       map[string]Box
	boxpointers map[*Box]bool
	sinks       map[string]Sink
	Edges       []DataflowEdge
}

type DataflowEdge struct {
	// From is the name of the source or box at the start of this edge.
	From string

	// To is the name of the box or sink at the end of this edge.
	To string

	// InputName is the name that the box at the end of the edge
	// expects incoming tuples to have. This has no meaning when there
	// is a sink at the end of this edge.
	InputName string
}

func NewDefaultStaticTopologyBuilder() StaticTopologyBuilder {
	tb := DefaultStaticTopologyBuilder{}
	tb.sources = make(map[string]Source)
	tb.boxes = make(map[string]Box)
	tb.boxpointers = make(map[*Box]bool)
	tb.sinks = make(map[string]Sink)
	tb.Edges = make([]DataflowEdge, 0)
	return &tb
}

// check if the given name can be used as a source, box, or sink
// name (i.e., it is not used yet)
func (tb *DefaultStaticTopologyBuilder) checkName(name string) error {
	_, alreadyExists := tb.sources[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a source called '%s'", name)
		return err
	}
	_, alreadyExists = tb.boxes[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a box called '%s'", name)
		return err
	}
	_, alreadyExists = tb.sinks[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a sink called '%s'", name)
		return err
	}
	return nil
}

// check if the given name is an existing box or source
func (tb *DefaultStaticTopologyBuilder) IsValidOutputReference(name string) bool {
	_, sourceExists := tb.sources[name]
	_, boxExists := tb.boxes[name]
	return (sourceExists || boxExists)
}

func (tb *DefaultStaticTopologyBuilder) AddSource(name string, source Source) SourceDeclarer {
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &DefaultSourceDeclarer{nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of source
	tb.sources[name] = source
	return &DefaultSourceDeclarer{}
}

func (tb *DefaultStaticTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &DefaultBoxDeclarer{err: nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of box
	tb.boxes[name] = box
	tb.boxpointers[&box] = true
	return &DefaultBoxDeclarer{tb, name, box, nil}
}

func (tb *DefaultStaticTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &DefaultSinkDeclarer{err: nameErr}
	}
	// keep track of sink
	tb.sinks[name] = sink
	return &DefaultSinkDeclarer{tb, name, sink, nil}
}

func (tb *DefaultStaticTopologyBuilder) Build() Topology {
	// every source and every box gets an "output pipe"
	pipes := make(map[string]*SequentialPipe, len(tb.sources)+len(tb.boxes))
	for name, _ := range tb.sources {
		pipe := SequentialPipe{}
		pipe.FromName = name
		pipe.ReceiverBoxes = make([]ReceiverBox, 0)
		pipe.ReceiverSinks = make([]ReceiverSink, 0)
		pipes[name] = &pipe
	}
	for name, _ := range tb.boxes {
		pipe := SequentialPipe{}
		pipe.FromName = name
		pipe.ReceiverBoxes = make([]ReceiverBox, 0)
		pipe.ReceiverSinks = make([]ReceiverSink, 0)
		pipes[name] = &pipe
	}
	// add the correct receivers to each pipe
	for _, edge := range tb.Edges {
		fromName := edge.From
		toName := edge.To
		pipe := pipes[fromName]
		// add the target of the edge (is either a sink or a box) to the
		// pipe's receiver list
		sink, isSink := tb.sinks[toName]
		if isSink {
			recv := ReceiverSink{toName, sink}
			pipe.ReceiverSinks = append(pipe.ReceiverSinks, recv)
		}
		box, isBox := tb.boxes[toName]
		if isBox {
			recv := ReceiverBox{toName, box, pipes[toName], edge.InputName}
			pipe.ReceiverBoxes = append(pipe.ReceiverBoxes, recv)
		}
	}
	// TODO source and sink is reference data,
	//      so cannot call .Build() more than once
	return &DefaultTopology{tb.boxpointers, tb.sources, pipes}
}

// holds a box and the writer that will receive this box's output
type ReceiverBox struct {
	Name      string
	Box       Box
	Receiver  Writer
	InputName string
}

// holds a sink and the sink's name
type ReceiverSink struct {
	Name string
	Sink Sink
}

// receives input from a box and forwards it to registered listeners
type SequentialPipe struct {
	FromName      string
	ReceiverBoxes []ReceiverBox
	ReceiverSinks []ReceiverSink
}

func (p *SequentialPipe) Write(t *tuple.Tuple) error {
	// add tracing information
	out := newDefaultEvent(tuple.OUTPUT, p.FromName)
	t.AddEvent(out)
	// forward tuple to connected boxes
	var s *tuple.Tuple

	// copy for all receivers but if this pipe has only
	// one receiver, there is no need to copy
	notNeedsCopy := len(p.ReceiverBoxes)+len(p.ReceiverSinks) <= 1
	for _, recvBox := range p.ReceiverBoxes {
		if notNeedsCopy {
			s = t
		} else {
			s = t.Copy()
		}
		// set the name that the box is expecting
		s.InputName = recvBox.InputName
		// add tracing information and hand over to box
		in := newDefaultEvent(tuple.INPUT, recvBox.Name)
		s.AddEvent(in)
		recvBox.Box.Process(s, recvBox.Receiver)
	}
	// forward tuple to connected sinks
	for _, recvSink := range p.ReceiverSinks {
		if notNeedsCopy {
			s = t
		} else {
			s = t.Copy()
		}
		// add tracing information and hand over to sink
		in := newDefaultEvent(tuple.INPUT, recvSink.Name)
		s.AddEvent(in)
		recvSink.Sink.Write(s)
	}
	return nil
}

func newDefaultEvent(inout tuple.InOutType, msg string) tuple.TraceEvent {
	return tuple.TraceEvent{
		time.Now(),
		inout,
		msg,
	}
}

/**************************************************/

type DefaultSourceDeclarer struct {
	err error
}

func (sd *DefaultSourceDeclarer) Err() error {
	return sd.err
}

/**************************************************/

type DefaultBoxDeclarer struct {
	tb   *DefaultStaticTopologyBuilder
	name string
	box  Box
	err  error
}

func (bd *DefaultBoxDeclarer) Input(refname string) BoxDeclarer {
	return bd.NamedInput(refname, "*")
}

func (bd *DefaultBoxDeclarer) NamedInput(refname string, inputName string) BoxDeclarer {
	// if there was a previous error, do nothing
	if bd.err != nil {
		return bd
	}
	// if the name can't be used, return an error
	if !bd.tb.IsValidOutputReference(refname) {
		err := fmt.Errorf("there is no box or source named '%s'", refname)
		bd.err = err
		return bd
	}
	// The `Input()` caller said that we should attach the name
	// `inputName` to incoming data (or not if inputName is "*").
	// This is ok if
	// - InputConstraints or InputConstraints.schema is nil
	// - there is a schema (or nil) declared in the InputConstraints
	//   with that name
	// - there is a "*" schema declared in the InputConstraints
	// Otherwise this is an error.
	ok := false
	inputConstraints, err := bd.box.InputConstraints()
	if err != nil {
		bd.err = err
		return bd
	}
	if inputConstraints == nil || inputConstraints.Schema == nil {
		ok = true
	} else if _, declared := inputConstraints.Schema[inputName]; declared {
		// TODO check if given schema matches the referenced source or box
		ok = true
	} else if _, declared := inputConstraints.Schema["*"]; declared {
		// TODO check if given schema matches the referenced source or box
		ok = true
	}
	if !ok {
		err := fmt.Errorf("you cannot use %s as an input name with input constraints %v",
			inputName, inputConstraints)
		bd.err = err
		return bd
	}
	// check if this edge already exists
	edge := DataflowEdge{refname, bd.name, inputName}
	edgeAlreadyExists := false
	for _, e := range bd.tb.Edges {
		edgeAlreadyExists = edge == e
		break
	}
	if edgeAlreadyExists {
		err := fmt.Errorf("box '%s' is already connected to '%s'",
			bd.name, refname)
		bd.err = err
		return bd
	}
	// if not, store it
	bd.tb.Edges = append(bd.tb.Edges, edge)
	return bd
}

func (bd *DefaultBoxDeclarer) Err() error {
	return bd.err
}

/**************************************************/

type DefaultSinkDeclarer struct {
	tb   *DefaultStaticTopologyBuilder
	name string
	sink Sink
	err  error
}

func (sd *DefaultSinkDeclarer) Input(refname string) SinkDeclarer {
	// if there was a previous error, do nothing
	if sd.err != nil {
		return sd
	}
	// if the name can't be used, return an error
	if !sd.tb.IsValidOutputReference(refname) {
		err := fmt.Errorf("there is no box or source named '%s'", refname)
		sd.err = err
		return sd
	}
	// check if this edge already exists
	edge := DataflowEdge{refname, sd.name, ""}
	edgeAlreadyExists := false
	for _, e := range sd.tb.Edges {
		edgeAlreadyExists = edge == e
		break
	}
	if edgeAlreadyExists {
		err := fmt.Errorf("box '%s' is already connected to '%s'",
			sd.name, refname)
		sd.err = err
		return sd
	}
	// if not, store it
	sd.tb.Edges = append(sd.tb.Edges, edge)
	return sd
}

func (sd *DefaultSinkDeclarer) Err() error {
	return sd.err
}

/* TODO the default source/sink/box do not belong here.
 * They are not part of the topology implementation.
 * In order to test whether topology setup works correctly,
 * source/box/sink with a dummy implementation should be
 * part of the test suite.
 */

/**************************************************/

type DefaultSource struct{}

func (s *DefaultSource) GenerateStream(w Writer) error {
	return nil
}

func (s *DefaultSource) Schema() *Schema {
	var sc Schema = Schema("test")
	return &sc
}

/**************************************************/

type DefaultBox struct {
}

func (b *DefaultBox) Init(ctx *Context) error {
	return nil
}

func (b *DefaultBox) Process(t *tuple.Tuple, s Writer) error {
	return nil
}

func (b *DefaultBox) InputConstraints() (*InputConstraints, error) {
	return nil, nil
}

func (b *DefaultBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

/**************************************************/

type DefaultSink struct{}

func (s *DefaultSink) Write(t *tuple.Tuple) error {
	return nil
}

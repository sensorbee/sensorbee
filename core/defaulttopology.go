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
	From string
	To   string
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
			recv := ReceiverBox{toName, box, pipes[toName]}
			pipe.ReceiverBoxes = append(pipe.ReceiverBoxes, recv)
		}
	}
	// TODO source and sink is reference data,
	//      so cannot call .Build() more than once
	return &DefaultTopology{tb.boxpointers, tb.sources, pipes}
}

// holds a box and the writer that will receive this box's output
type ReceiverBox struct {
	Name     string
	Box      Box
	Receiver Writer
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
	out := newDefaultTracer(tuple.OUTPUT, p.FromName)
	t.AddTracer(out)
	// forward tuple to connected boxes
	var s *tuple.Tuple
	var tes []tuple.TraceEvent
	tupleCopies := 0
	for _, recvBox := range p.ReceiverBoxes {
		// copy for all receivers but the first so that
		// multiple receivers don't operate on the same data
		if tupleCopies == 0 {
			s = t
			tes = make([]tuple.TraceEvent, len(t.Trace))
			copy(tes, t.Trace)
		} else {
			s = t.Copy()
			s.Trace = tes
		}
		in := newDefaultTracer(tuple.INPUT, recvBox.Name)
		s.AddTracer(in)
		recvBox.Box.Process(s, recvBox.Receiver)
		tupleCopies += 1
	}
	// forward tuple to connected sinks
	for _, recvSink := range p.ReceiverSinks {
		// copy for all receivers but the first so that
		// multiple receivers don't operate on the same data
		if tupleCopies == 0 {
			s = t
			tes = make([]tuple.TraceEvent, len(t.Trace))
			copy(tes, t.Trace)
		} else {
			s = t.Copy()
			s.Trace = tes
		}
		in := newDefaultTracer(tuple.INPUT, recvSink.Name)
		s.AddTracer(in)
		recvSink.Sink.Write(s)
		tupleCopies += 1
	}
	return nil
}

func newDefaultTracer(inout tuple.InOutType, msg string) tuple.TraceEvent {
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

func (bd *DefaultBoxDeclarer) Input(refname string, schema *Schema) BoxDeclarer {
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
	// TODO check if given schema matches the referenced source or box
	// check if this edge already exists
	edge := DataflowEdge{refname, bd.name}
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
	edge := DataflowEdge{refname, sd.name}
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

func (b *DefaultBox) RequiredInputSchema() ([]*Schema, error) {
	return []*Schema{nil}, nil
}

func (b *DefaultBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

/**************************************************/

type DefaultSink struct{}

func (s *DefaultSink) Write(t *tuple.Tuple) error {
	return nil
}

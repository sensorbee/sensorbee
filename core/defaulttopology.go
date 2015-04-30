package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
	"sync"
	"time"
)

type defaultStaticTopology struct {
	// tb.boxes may contain multiple instances of the same Box object.
	// Use a set-like map to avoid calling Init() twice on the same object.
	boxpointers map[*Box]bool

	sources map[string]Source
	pipes   map[string]*capacityPipe
}

func (t *defaultStaticTopology) Run(ctx *Context) {
	for box, _ := range t.boxpointers {
		(*box).Init(ctx)
	}

	// launch a listener goroutine for each receiver of a pipe
	for _, pipe := range t.pipes {
		for _, recv := range pipe.Receivers {
			recv := recv
			go recv.ProcessItems()
		}
	}

	var wg sync.WaitGroup
	for name, source := range t.sources {
		wg.Add(1)
		go func(name string, source Source) {
			defer wg.Done()
			source.GenerateStream(t.pipes[name])
		}(name, source)
	}
	wg.Wait()
	// as a workaround, sleep a bit so that background goroutines can
	// finish their work (or else all tests will break)
	// TODO replace this by proper shutdown method
	time.Sleep(50 * time.Millisecond)
}

/**************************************************/

type defaultStaticTopologyBuilder struct {
	sources     map[string]Source
	boxes       map[string]Box
	boxpointers map[*Box]bool
	sinks       map[string]Sink
	Edges       []dataflowEdge
}

type dataflowEdge struct {
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
	tb := defaultStaticTopologyBuilder{}
	tb.sources = make(map[string]Source)
	tb.boxes = make(map[string]Box)
	tb.boxpointers = make(map[*Box]bool)
	tb.sinks = make(map[string]Sink)
	tb.Edges = make([]dataflowEdge, 0)
	return &tb
}

// check if the given name can be used as a source, box, or sink
// name (i.e., it is not used yet)
func (tb *defaultStaticTopologyBuilder) checkName(name string) error {
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
func (tb *defaultStaticTopologyBuilder) IsValidOutputReference(name string) bool {
	_, sourceExists := tb.sources[name]
	_, boxExists := tb.boxes[name]
	return (sourceExists || boxExists)
}

func (tb *defaultStaticTopologyBuilder) AddSource(name string, source Source) SourceDeclarer {
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &defaultSourceDeclarer{nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of source
	tb.sources[name] = source
	return &defaultSourceDeclarer{}
}

func (tb *defaultStaticTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &defaultBoxDeclarer{err: nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of box
	tb.boxes[name] = box
	tb.boxpointers[&box] = true
	return &defaultBoxDeclarer{tb, name, box, nil}
}

func (tb *defaultStaticTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &defaultSinkDeclarer{err: nameErr}
	}
	// keep track of sink
	tb.sinks[name] = sink
	return &defaultSinkDeclarer{tb, name, sink, nil}
}

func (tb *defaultStaticTopologyBuilder) makeCapacityPipes() map[string]*capacityPipe {
	pipes := make(map[string]*capacityPipe, len(tb.sources)+len(tb.boxes))
	for name, _ := range tb.sources {
		pipe := NewCapacityPipe()
		pipe.FromName = name
		pipes[name] = pipe
	}
	for name, _ := range tb.boxes {
		pipe := NewCapacityPipe()
		pipe.FromName = name
		pipes[name] = pipe
	}
	return pipes
}

func (tb *defaultStaticTopologyBuilder) Build() StaticTopology {
	// every source and every box gets an "output pipe"
	pipes := tb.makeCapacityPipes()
	// add the correct receivers to each pipe
	for _, edge := range tb.Edges {
		fromName := edge.From
		toName := edge.To
		pipe := pipes[fromName]
		// add the target of the edge (is either a sink or a box) to the
		// pipe's receiver list
		sink, isSink := tb.sinks[toName]
		if isSink {
			recv := receiverSink{toName, sink, make(chan *tuple.Tuple)}
			pipe.Receivers = append(pipe.Receivers, &recv)
		}
		box, isBox := tb.boxes[toName]
		if isBox {
			recv := receiverBox{toName, box, pipes[toName], edge.InputName, make(chan *tuple.Tuple)}
			pipe.Receivers = append(pipe.Receivers, &recv)
		}
	}
	// TODO source and sink is reference data,
	//      so cannot call .Build() more than once
	return &defaultStaticTopology{tb.boxpointers, tb.sources, pipes}
}

/**************************************************/

// A pipeReceiver represents an entity that can receive input from a pipe,
// i.e., a Box or a Sink.
//
// AddToChannel will add the passed tuple to the internal channel
// (i.e., it will block if there is still a process running) and
// schedule it for execution.
//
// InputName returns the name that the actual receiver Box or Sink
// expects to be associated with incoming tuples.
//
// ProcessItems() is a process to be launched as a goroutine that
// attaches to the internal channel and then passes each item in that
// channel to the actual receiver Box or Sink. Calling this multiple
// times will parallelize execution of items in the channel.
type pipeReceiver interface {
	AddToChannel(t *tuple.Tuple)
	InputName() string
	ProcessItems()
}

// holds a box and the writer that will receive this box's output
type receiverBox struct {
	Name      string
	Box       Box
	Receiver  Writer
	inputName string
	buffer    chan *tuple.Tuple
}

func (rb *receiverBox) ProcessItems() {
	// If an item is put in the queue (and there is no other processing
	// taking place), that item will be processed immediately.
	// If we are still in the middle of processing, then this item will
	// be processed later.
	// Note that if we increase the capacity of the channel here, then
	// this pipe will have a buffering functionality, but it will not
	// increase the parallelism of downstream operations.
	for t := range rb.buffer {
		// add tracing information and hand over to box
		in := newDefaultEvent(tuple.INPUT, rb.Name)
		t.AddEvent(in)
		rb.Box.Process(t, rb.Receiver)
	}
}

func (rb *receiverBox) AddToChannel(t *tuple.Tuple) {
	rb.buffer <- t
}

func (rb *receiverBox) InputName() string {
	return rb.inputName
}

// holds a sink and the sink's name
type receiverSink struct {
	Name   string
	Sink   Sink
	buffer chan *tuple.Tuple
}

func (rs *receiverSink) ProcessItems() {
	// If an item is put in the queue (and there is no other processing
	// taking place), that item will be processed immediately.
	// If we are still in the middle of processing, then this item will
	// be processed later.
	// Note that if we increase the capacity of the channel here, then
	// this pipe will have a buffering functionality, but it will not
	// increase the parallelism of downstream operations.
	for t := range rs.buffer {
		// add tracing information and hand over to sink
		in := newDefaultEvent(tuple.INPUT, rs.Name)
		t.AddEvent(in)
		rs.Sink.Write(t)
	}
}

func (rs *receiverSink) AddToChannel(t *tuple.Tuple) {
	rs.buffer <- t
}

func (rs *receiverSink) InputName() string {
	return "output"
}

/**************************************************/

func NewCapacityPipe() *capacityPipe {
	p := capacityPipe{}
	return &p
}

// receives input from a box and forwards it to registered listeners
type capacityPipe struct {
	FromName  string
	Receivers []pipeReceiver
}

func (p *capacityPipe) processItem(t *tuple.Tuple) {
	// forward tuple to connected boxes
	var s *tuple.Tuple

	// copy for all receivers but if this pipe has only
	// one receiver, there is no need to copy
	notNeedsCopy := len(p.Receivers) <= 1
	for _, recv := range p.Receivers {
		if notNeedsCopy {
			s = t
		} else {
			s = t.Copy()
		}
		// set the name that the box is expecting
		s.InputName = recv.InputName()
		// If the call below blocks, it means that the receiving
		// box/sink is still busy processing the previous item(s).
		recv.AddToChannel(s)
	}
}

func (p *capacityPipe) Write(t *tuple.Tuple) error {
	// add tracing information
	out := newDefaultEvent(tuple.OUTPUT, p.FromName)
	t.AddEvent(out)
	// distribute this item to receivers
	p.processItem(t)
	return nil
}

/**************************************************/

func newDefaultEvent(inout tuple.EventType, msg string) tuple.TraceEvent {
	return tuple.TraceEvent{
		time.Now(),
		inout,
		msg,
	}
}

/**************************************************/

type defaultSourceDeclarer struct {
	err error
}

func (sd *defaultSourceDeclarer) Err() error {
	return sd.err
}

/**************************************************/

type defaultBoxDeclarer struct {
	tb   *defaultStaticTopologyBuilder
	name string
	box  Box
	err  error
}

func (bd *defaultBoxDeclarer) Input(refname string) BoxDeclarer {
	return bd.NamedInput(refname, "*")
}

func (bd *defaultBoxDeclarer) NamedInput(refname string, inputName string) BoxDeclarer {
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
	// - InputConstraints() or InputConstraints().schema is nil
	// - there is a schema (or nil) declared in the InputConstraints()
	//   with that name
	// - there is a "*" schema declared in the InputConstraints()
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
	edge := dataflowEdge{refname, bd.name, inputName}
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

func (bd *defaultBoxDeclarer) Err() error {
	return bd.err
}

/**************************************************/

type defaultSinkDeclarer struct {
	tb   *defaultStaticTopologyBuilder
	name string
	sink Sink
	err  error
}

func (sd *defaultSinkDeclarer) Input(refname string) SinkDeclarer {
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
	edge := dataflowEdge{refname, sd.name, ""}
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

func (sd *defaultSinkDeclarer) Err() error {
	return sd.err
}

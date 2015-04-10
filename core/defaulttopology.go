package core

import (
	"fmt"
	"log"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type DefaultTopology struct {
}

func (this *DefaultTopology) Run() {}

/**************************************************/

type DefaultTopologyBuilder struct {
	sources map[string]Source
	boxes   map[string]Box
	sinks   map[string]Sink
	Edges   []DataflowEdge
}

type DataflowEdge struct {
	From string
	To   string
}

func NewDefaultTopologyBuilder() TopologyBuilder {
	tb := DefaultTopologyBuilder{}
	tb.sources = make(map[string]Source)
	tb.boxes = make(map[string]Box)
	tb.sinks = make(map[string]Sink)
	tb.Edges = make([]DataflowEdge, 0)
	return &tb
}

// check if the given name can be used as a source, box, or sink
// name (i.e., it is not used yet)
func (this *DefaultTopologyBuilder) checkName(name string) error {
	_, alreadyExists := this.sources[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a source called '%s'", name)
		return err
	}
	_, alreadyExists = this.boxes[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a box called '%s'", name)
		return err
	}
	_, alreadyExists = this.sinks[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a sink called '%s'", name)
		return err
	}
	return nil
}

// check if the given name is an existing box or source
func (this *DefaultTopologyBuilder) IsValidOutputReference(name string) bool {
	_, sourceExists := this.sources[name]
	_, boxExists := this.boxes[name]
	return (sourceExists || boxExists)
}

func (this *DefaultTopologyBuilder) AddSource(name string, source Source) SourceDeclarer {
	// check name
	if nameErr := this.checkName(name); nameErr != nil {
		return &DefaultSourceDeclarer{nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of source
	this.sources[name] = source
	return &DefaultSourceDeclarer{}
}

func (this *DefaultTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	// check name
	if nameErr := this.checkName(name); nameErr != nil {
		return &DefaultBoxDeclarer{err: nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of box
	this.boxes[name] = box
	return &DefaultBoxDeclarer{this, name, box, nil}
}

func (this *DefaultTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	// check name
	if nameErr := this.checkName(name); nameErr != nil {
		return &DefaultSinkDeclarer{err: nameErr}
	}
	// keep track of sink
	this.sinks[name] = sink
	return &DefaultSinkDeclarer{this, name, sink, nil}
}

func (this *DefaultTopologyBuilder) Build() Topology {
	// every source and every box gets an "output pipe"
	pipes := make(map[string]*SequentialPipe, len(this.sources)+len(this.boxes))
	for name, _ := range this.sources {
		pipe := SequentialPipe{}
		pipe.ReceiverBoxes = make([]ReceiverBox, 0)
		pipe.ReceiverSinks = make([]Sink, 0)
		pipes[name] = &pipe
	}
	for name, _ := range this.boxes {
		pipe := SequentialPipe{}
		pipe.ReceiverBoxes = make([]ReceiverBox, 0)
		pipe.ReceiverSinks = make([]Sink, 0)
		pipes[name] = &pipe
	}
	// add the correct receivers to each pipe
	for _, edge := range this.Edges {
		fromName := edge.From
		toName := edge.To
		pipe := pipes[fromName]
		// add the target of the edge (is either a sink or a box) to the
		// pipe's receiver list
		sink, isSink := this.sinks[toName]
		if isSink {
			pipe.ReceiverSinks = append(pipe.ReceiverSinks, sink)
		}
		box, isBox := this.boxes[toName]
		if isBox {
			recv := ReceiverBox{box, pipes[toName]}
			pipe.ReceiverBoxes = append(pipe.ReceiverBoxes, recv)
		}
	}
	// debug
	for _, pipe := range pipes {
		log.Printf("%#v\n", pipe)
	}
	// more debug
	for name, source := range this.sources {
		source.GenerateStream(pipes[name])
	}
	return &DefaultTopology{}
}

// holds a box and the writer that will receive this box's output
type ReceiverBox struct {
	Box      Box
	Receiver Writer
}

// receives input from a box and forwards it to registered listeners
type SequentialPipe struct {
	ReceiverBoxes []ReceiverBox
	ReceiverSinks []Sink
}

func (this *SequentialPipe) Write(t *tuple.Tuple) error {
	log.Printf("pipe %p: received %v\n", this, t)
	// forward tuple to connected boxes
	for _, recvBox := range this.ReceiverBoxes {
		recvBox.Box.Process(t, recvBox.Receiver)
	}
	// forward tuple to connected sinks
	for _, sink := range this.ReceiverSinks {
		sink.Write(t)
	}
	return nil
}

/**************************************************/

type DefaultSourceDeclarer struct {
	err error
}

func (this *DefaultSourceDeclarer) Err() error {
	return this.err
}

/**************************************************/

type DefaultBoxDeclarer struct {
	tb   *DefaultTopologyBuilder
	name string
	box  Box
	err  error
}

func (this *DefaultBoxDeclarer) Input(refname string, schema *Schema) BoxDeclarer {
	// if there was a previous error, do nothing
	if this.err != nil {
		return this
	}
	// if the name can't be used, return an error
	if !this.tb.IsValidOutputReference(refname) {
		err := fmt.Errorf("there is no box or source named '%s'", refname)
		this.err = err
		return this
	}
	// TODO check if given schema matches the referenced source or box
	// check if this edge already exists
	edge := DataflowEdge{refname, this.name}
	edgeAlreadyExists := false
	for _, e := range this.tb.Edges {
		edgeAlreadyExists = edge == e
		break
	}
	if edgeAlreadyExists {
		err := fmt.Errorf("box '%s' is already connected to '%s'",
			this.name, refname)
		this.err = err
		return this
	}
	// if not, store it
	this.tb.Edges = append(this.tb.Edges, edge)
	return this
}

func (this *DefaultBoxDeclarer) Err() error {
	return this.err
}

/**************************************************/

type DefaultSinkDeclarer struct {
	tb   *DefaultTopologyBuilder
	name string
	sink Sink
	err  error
}

func (this *DefaultSinkDeclarer) Input(refname string) SinkDeclarer {
	// if there was a previous error, do nothing
	if this.err != nil {
		return this
	}
	// if the name can't be used, return an error
	if !this.tb.IsValidOutputReference(refname) {
		err := fmt.Errorf("there is no box or source named '%s'", refname)
		this.err = err
		return this
	}
	// check if this edge already exists
	edge := DataflowEdge{refname, this.name}
	edgeAlreadyExists := false
	for _, e := range this.tb.Edges {
		edgeAlreadyExists = edge == e
		break
	}
	if edgeAlreadyExists {
		err := fmt.Errorf("box '%s' is already connected to '%s'",
			this.name, refname)
		this.err = err
		return this
	}
	// if not, store it
	this.tb.Edges = append(this.tb.Edges, edge)
	return this
}

func (this *DefaultSinkDeclarer) Err() error {
	return this.err
}

/**************************************************/

type DefaultSource struct{}

func (this *DefaultSource) GenerateStream(w Writer) error {
	log.Printf("start processing\n")
	log.Printf("emit first item\n")
	t := &tuple.Tuple{}
	w.Write(t)
	log.Printf("emit second item\n")
	t2 := &tuple.Tuple{}
	w.Write(t2)
	return nil
}

func (this *DefaultSource) Schema() *Schema {
	var s Schema = Schema("test")
	return &s
}

/**************************************************/

type DefaultBox struct{}

func (this *DefaultBox) Process(t *tuple.Tuple, s Writer) error {
	log.Printf("box  %p: processing %v\n", this, t)
	t.BatchID = 1
	s.Write(t)
	return nil
}

func (this *DefaultBox) RequiredInputSchema() ([]*Schema, error) {
	return []*Schema{nil}, nil
}

func (this *DefaultBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

/**************************************************/

type DefaultSink struct {
	x int
}

func (this *DefaultSink) Write(t *tuple.Tuple) error {
	log.Printf("sink %p: received %v\n", this, t)
	return nil
}

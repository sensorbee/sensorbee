package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type DefaultTopology struct {
}

func (this *DefaultTopology) Run() {}

type DefaultTopologyBuilder struct {
	sources map[string]Source
	boxes   map[string]Box
	sinks   map[string]Sink
}

func NewDefaultTopologyBuilder() TopologyBuilder {
	tb := DefaultTopologyBuilder{}
	tb.sources = make(map[string]Source)
	tb.boxes = make(map[string]Box)
	tb.sinks = make(map[string]Sink)
	return &tb
}

func (this *DefaultTopologyBuilder) AddSource(name string, source Source) SourceDeclarer {
	// check name
	_, alreadyExists := this.sources[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a source called '%s'", name)
		return &DefaultSourceDeclarer{err}
	}
	_, alreadyExists = this.boxes[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a box called '%s'", name)
		return &DefaultSourceDeclarer{err}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of source
	this.sources[name] = source
	return &DefaultSourceDeclarer{}
}

func (this *DefaultTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	// check name
	_, alreadyExists := this.sources[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a source called '%s'", name)
		return &DefaultBoxDeclarer{err}
	}
	_, alreadyExists = this.boxes[name]
	if alreadyExists {
		err := fmt.Errorf("there is already a box called '%s'", name)
		return &DefaultBoxDeclarer{err}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of box
	this.boxes[name] = box
	return &DefaultBoxDeclarer{}
}

func (this *DefaultTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	return &DefaultSinkDeclarer{}
}

func (this *DefaultTopologyBuilder) Build() Topology {
	return &DefaultTopology{}
}

type DefaultSourceDeclarer struct {
	err error
}

func (this *DefaultSourceDeclarer) Err() error {
	return this.err
}

type DefaultBoxDeclarer struct {
	err error
}

func (this *DefaultBoxDeclarer) Input(name string, schema *Schema) BoxDeclarer {
	return &DefaultBoxDeclarer{}
}

func (this *DefaultBoxDeclarer) Err() error {
	return this.err
}

type DefaultSinkDeclarer struct{}

func (this *DefaultSinkDeclarer) Input(name string) SinkDeclarer {
	return &DefaultSinkDeclarer{}
}

func (this *DefaultSinkDeclarer) Err() error {
	return nil
}

type DefaultSource struct{}

func (this *DefaultSource) GenerateStream(w Writer) error {
	return nil
}

func (this *DefaultSource) Schema() *Schema {
	var s Schema = Schema("test")
	return &s
}

type DefaultBox struct{}

func (this *DefaultBox) Process(t *tuple.Tuple, s Sink) error {
	return nil
}

func (this *DefaultBox) RequiredInputSchema() ([]*Schema, error) {
	return []*Schema{nil}, nil
}

func (this *DefaultBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

type DefaultSink struct{}

func (this *DefaultSink) Write(t *tuple.Tuple) error {
	return nil
}

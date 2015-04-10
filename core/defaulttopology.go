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
		return &DefaultBoxDeclarer{nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	// keep track of box
	this.boxes[name] = box
	return &DefaultBoxDeclarer{}
}

func (this *DefaultTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	// check name
	if nameErr := this.checkName(name); nameErr != nil {
		return &DefaultSinkDeclarer{nameErr}
	}
	// keep track of sink
	this.sinks[name] = sink
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

type DefaultSinkDeclarer struct {
	err error
}

func (this *DefaultSinkDeclarer) Input(name string) SinkDeclarer {
	return &DefaultSinkDeclarer{}
}

func (this *DefaultSinkDeclarer) Err() error {
	return this.err
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

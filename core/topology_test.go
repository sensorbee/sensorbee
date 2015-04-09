package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
)

type DummyTopology struct{}

func (this *DummyTopology) Run() {}

type DummyTopologyBuilder struct{}

func (this *DummyTopologyBuilder) AddSource(name string, source Source) SourceDeclarer {
	return nil
}
func (this *DummyTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	return &DummyBoxDeclarer{}
}
func (this *DummyTopologyBuilder) Build() Topology {
	return &DummyTopology{}
}

type DummyBoxDeclarer struct{}

func (this *DummyBoxDeclarer) Input(name string, schema *Schema) BoxDeclarer {
	return nil
}

type DummySource struct{}

func (this *DummySource) GenerateStream(w Writer) error {
	return nil
}
func (this *DummySource) Schema() *Schema {
	var s Schema = Schema("test")
	return &s
}

type DummyBox struct{}

func (this *DummyBox) Process(t *tuple.Tuple, s Sink) error {
	return nil
}
func (this *DummyBox) RequiredInputSchema() ([]*Schema, error) {
	return []*Schema{nil}, nil
}
func (this *DummyBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

func TestTopology(t *testing.T) {
	Convey("Given dummy topology builder, to build topology", t, func() {

		source := &DummySource{}
		tb := DummyTopologyBuilder{}
		tb.AddSource("test_source1", source)

		box := &DummyBox{}
		inputSchema, _ := box.RequiredInputSchema()
		tb.AddBox("test_box1", box).
			Input("test_input_schema", inputSchema[0])

		t := tb.Build()

		Convey("It should satisfy Topology interface", func() {
			So(t, ShouldNotBeNil)
		})
	})
}

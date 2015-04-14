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
	return &DummySourceDeclarer{}
}
func (this *DummyTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	return &DummyBoxDeclarer{}
}
func (this *DummyTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	return &DummySinkDeclarer{}
}
func (this *DummyTopologyBuilder) Build() Topology {
	return &DummyTopology{}
}

type DummySourceDeclarer struct{}

func (this *DummySourceDeclarer) Err() error {
	return nil
}

type DummyBoxDeclarer struct{}

func (this *DummyBoxDeclarer) Input(name string, schema *Schema) BoxDeclarer {
	return &DummyBoxDeclarer{}
}

func (this *DummyBoxDeclarer) Err() error {
	return nil
}

type DummySinkDeclarer struct{}

func (this *DummySinkDeclarer) Input(name string) SinkDeclarer {
	return &DummySinkDeclarer{}
}

func (this *DummySinkDeclarer) Err() error {
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

func (this *DummyBox) Init(ctx *Context) error {
	return nil
}
func (this *DummyBox) Process(t *tuple.Tuple, s Writer) error {
	return nil
}
func (this *DummyBox) RequiredInputSchema() ([]*Schema, error) {
	return []*Schema{nil}, nil
}
func (this *DummyBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}

type DummySink struct{}

func (this *DummySink) Write(t *tuple.Tuple) error {
	return nil
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

	Convey("Given dummy topology builder, to build complex topology", t, func() {

		so1 := &DummySource{}
		so2 := &DummySource{}
		b1 := &DummyBox{}
		b2 := &DummyBox{}
		b3 := &DummyBox{}
		si1 := &DummySink{}
		si2 := &DummySink{}

		tb := DummyTopologyBuilder{}

		var err DeclarerError
		err = tb.AddSource("so1", so1)
		So(err.Err(), ShouldBeNil)
		err = tb.AddSource("so2", so2)
		So(err.Err(), ShouldBeNil)
		err = tb.AddBox("b1", b1).
			Input("so1", nil).
			Input("so2", nil)
		So(err.Err(), ShouldBeNil)
		err = tb.AddBox("b2", b2).
			Input("so1", nil)
		So(err.Err(), ShouldBeNil)
		err = tb.AddBox("b3", b3).
			Input("b1", nil)
		So(err.Err(), ShouldBeNil)
		err = tb.AddSink("si1", si1).
			Input("b2").
			Input("b1").
			Input("so1")
		So(err.Err(), ShouldBeNil)
		err = tb.AddSink("si2", si2).
			Input("b2").
			Input("b3")
		So(err.Err(), ShouldBeNil)

	})
}

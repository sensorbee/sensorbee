package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
)

type DummyTopology struct{}

func (dt *DummyTopology) Run(ctx *Context) error {
	return nil
}
func (dt *DummyTopology) Stop(ctx *Context) error {
	return nil
}

type DummyTopologyBuilder struct{}

func (dtb *DummyTopologyBuilder) AddSource(name string, source Source) SourceDeclarer {
	return &DummySourceDeclarer{}
}
func (dtb *DummyTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	return &DummyBoxDeclarer{}
}
func (dtb *DummyTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	return &DummySinkDeclarer{}
}
func (dtb *DummyTopologyBuilder) Build() StaticTopology {
	return &DummyTopology{}
}

type DummySourceDeclarer struct{}

func (dsd *DummySourceDeclarer) Err() error {
	return nil
}

type DummyBoxDeclarer struct{}

func (dbd *DummyBoxDeclarer) Input(name string) BoxDeclarer {
	return &DummyBoxDeclarer{}
}

func (dbd *DummyBoxDeclarer) NamedInput(name string, inputName string) BoxDeclarer {
	return &DummyBoxDeclarer{}
}

func (dbd *DummyBoxDeclarer) Err() error {
	return nil
}

type DummySinkDeclarer struct{}

func (dsd *DummySinkDeclarer) Input(name string) SinkDeclarer {
	return &DummySinkDeclarer{}
}

func (dsd *DummySinkDeclarer) Err() error {
	return nil
}

type DummySource struct{}

func (ds *DummySource) GenerateStream(ctx *Context, w Writer) error {
	return nil
}
func (ds *DummySource) Stop(ctx *Context) error {
	return nil
}
func (ds *DummySource) Schema() *Schema {
	var s Schema = Schema("test")
	return &s
}

type DummyBox struct{}

func (db *DummyBox) Init(ctx *Context) error {
	return nil
}
func (db *DummyBox) Process(ctx *Context, t *tuple.Tuple, s Writer) error {
	return nil
}
func (db *DummyBox) InputConstraints() (*BoxInputConstraints, error) {
	return nil, nil
}
func (db *DummyBox) OutputSchema(s []*Schema) (*Schema, error) {
	return nil, nil
}
func (db *DummyBox) Terminate(ctx *Context) error {
	return nil
}

type DummySink struct{}

func (ds *DummySink) Write(ctx *Context, t *tuple.Tuple) error {
	return nil
}
func (ds *DummySink) Close(ctx *Context) error {
	return nil
}

func TestTopology(t *testing.T) {
	Convey("Given dummy topology builder, to build topology", t, func() {

		source := &DummySource{}
		tb := DummyTopologyBuilder{}
		tb.AddSource("test_source1", source)

		box := &DummyBox{}
		tb.AddBox("test_box1", box).
			Input("test_input_schema")

		t := tb.Build()

		Convey("It should satisfy StaticTopology interface", func() {
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
			Input("so1").
			Input("so2")
		So(err.Err(), ShouldBeNil)
		err = tb.AddBox("b2", b2).
			Input("so1")
		So(err.Err(), ShouldBeNil)
		err = tb.AddBox("b3", b3).
			Input("b1")
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

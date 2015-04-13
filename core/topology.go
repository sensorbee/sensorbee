package core

type Topology interface {
	Run()
}

type StaticTopologyBuilder interface {
	AddSource(name string, source Source) SourceDeclarer
	AddBox(name string, box Box) BoxDeclarer
	AddSink(name string, sink Sink) SinkDeclarer

	Build() Topology
}

type DeclarerError interface {
	Err() error
}

type SourceDeclarer interface {
	DeclarerError
}

type BoxDeclarer interface {
	Input(name string, schema *Schema) BoxDeclarer
	DeclarerError
}

type SinkDeclarer interface {
	Input(name string) SinkDeclarer
	DeclarerError
}

package core

type Topology interface {
	Run()
}

type TopologyBuilder interface {
	AddSource(name string, source Source) SourceDeclarer
	AddBox(name string, box Box) BoxDeclarer
	AddSink(name string, sink Sink) SinkDeclarer

	Build() Topology
}

type SourceDeclarer interface {
}

type BoxDeclarer interface {
	Input(name string, schema *Schema) BoxDeclarer
}

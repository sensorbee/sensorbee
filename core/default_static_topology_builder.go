package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
	"strings"
	"sync"
)

type defaultStaticTopologyBuilder struct {
	sources   map[string]Source
	boxes     map[string]Box
	sinks     map[string]Sink
	Edges     []dataflowEdge
	builtFlag bool

	// addedComponents has pointers of components added to the builder.
	// It has pointers to sources, boxes, and sinks.
	addedComponents map[interface{}]string

	// TODO: provide a flag whether the builder accepts duplicated registrations
	// of the same instance of a Source, a Box, or a Sink (i.e. provide a way
	// to disable addComponent method)
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

// NewDefaultStaticTopologyBuilder creates a default implementation of a
// StaticTopologyBuilder. Note that this implementation is not thread-safe,
// i.e., it is not safe to have, say, multiple calls to AddBox in parallel.
// Also, this implementation will return an error when calling
// Build more than once.
func NewDefaultStaticTopologyBuilder() StaticTopologyBuilder {
	return &defaultStaticTopologyBuilder{
		sources:         map[string]Source{},
		boxes:           map[string]Box{},
		sinks:           map[string]Sink{},
		Edges:           []dataflowEdge{},
		builtFlag:       false,
		addedComponents: map[interface{}]string{},
	}
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
	if tb.builtFlag {
		err := fmt.Errorf(topologyBuilderAlreadyCalledBuildMsg)
		return &defaultSourceDeclarer{err}
	}
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &defaultSourceDeclarer{nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	if err := tb.addComponent(name, source); err != nil {
		return &defaultSourceDeclarer{err}
	}
	// keep track of source
	tb.sources[name] = source
	return &defaultSourceDeclarer{}
}

func (tb *defaultStaticTopologyBuilder) addComponent(name string, c interface{}) error {
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Func { // TODO: Use reflect.Type.Comparable if it's available (> Go 1.4)
		// When the component isn't comparable (e.g. functions), its uniqueness
		// cannot be guaranteed.
		return nil
	}

	{
		v := reflect.ValueOf(c)
		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct && v.NumField() == 0 {
			// When the version of Go is < 1.5, all instances of structs having
			// no field will have the same address and uniqueness cannot be
			// determined.
			return nil
		}
	}

	if name, ok := tb.addedComponents[c]; ok {
		return fmt.Errorf("the component is already added to the builder with the name '%v'", name)
	}
	tb.addedComponents[c] = name
	return nil
}

func (tb *defaultStaticTopologyBuilder) AddBox(name string, box Box) BoxDeclarer {
	if tb.builtFlag {
		err := fmt.Errorf(topologyBuilderAlreadyCalledBuildMsg)
		return &defaultBoxDeclarer{err: err}
	}
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &defaultBoxDeclarer{err: nameErr}
	}
	// TODO check that declared schema is a valid JSON Schema string
	if err := tb.addComponent(name, box); err != nil {
		return &defaultBoxDeclarer{err: err}
	}
	// keep track of box
	tb.boxes[name] = box
	return &defaultBoxDeclarer{tb, name, box, nil}
}

func (tb *defaultStaticTopologyBuilder) AddSink(name string, sink Sink) SinkDeclarer {
	if tb.builtFlag {
		err := fmt.Errorf(topologyBuilderAlreadyCalledBuildMsg)
		return &defaultSinkDeclarer{err: err}
	}
	// check name
	if nameErr := tb.checkName(name); nameErr != nil {
		return &defaultSinkDeclarer{err: nameErr}
	}
	if err := tb.addComponent(name, sink); err != nil {
		return &defaultSinkDeclarer{err: err}
	}
	// keep track of sink
	tb.sinks[name] = sink
	return &defaultSinkDeclarer{tb, name, sink, nil}
}

func (tb *defaultStaticTopologyBuilder) Build() (StaticTopology, error) {
	if tb.builtFlag {
		return nil, fmt.Errorf(topologyBuilderAlreadyCalledBuildMsg)
	}
	if len(tb.sources) == 0 {
		return nil, fmt.Errorf("there must be at least one source")
	}
	if has, path := tb.hasCycle(); has {
		return nil, fmt.Errorf("the topology has a cycle: %v", strings.Join(path, "->"))
	}

	stateMutex := &sync.Mutex{}
	st := &defaultStaticTopology{
		srcs:  tb.sources,
		boxes: tb.boxes,
		sinks: tb.sinks,

		srcDsts: map[string]WriteCloser{},
		nodes:   map[string]*staticNode{},

		state:      newTopologyStateHolder(stateMutex),
		stateMutex: stateMutex,
	}

	// Create st.nodes and its next writer
	dsts := map[string]*staticDestinations{}
	for name := range tb.sources {
		dsts[name] = newStaticDestinations()
	}
	for name, box := range tb.boxes {
		dst := newStaticDestinations()
		st.nodes[name] = newStaticNode(newBoxWriterAdapter(box, name, dst))
		dsts[name] = dst
	}
	for name, sink := range tb.sinks {
		st.nodes[name] = newStaticNode(newTraceWriter(sink, tuple.Input, name))
	}

	for _, e := range tb.Edges {
		r, s := newStaticPipe(e.InputName, 1024) // TODO: make capacity customizable
		st.nodes[e.To].addInput(e.From, r)
		dsts[e.From].addDestination(e.To, s)
	}

	for name := range tb.sources {
		st.srcDsts[name] = dsts[name]
	}

	tb.builtFlag = true
	return st, nil
}

// hasCycle returns true when the topology has a cycle.
// It also returns the path on a cycle.
func (tb *defaultStaticTopologyBuilder) hasCycle() (bool, []string) {
	// assumes there's at least one source.
	adj := map[string][]string{}
	for _, e := range tb.Edges {
		adj[e.From] = append(adj[e.From], e.To)
	}

	visited := map[string]int{} // 0: not yet, 1: visiting, 2: visited
	for s := range tb.sources {
		path := tb.detectCycle(s, adj, visited)
		if len(path) != 0 {
			for i := 0; i < len(path)/2; i++ {
				p := len(path) - i - 1
				path[i], path[p] = path[p], path[i]
			}
			return true, path
		}
	}

	// TODO: visited can be used to detect unused boxes or sinks
	return false, nil
}

// detectCycle returns non-empty path in the reverse order when it detected a cycle in the graph.
func (tb *defaultStaticTopologyBuilder) detectCycle(node string, adj map[string][]string, visited map[string]int) []string {
	switch visited[node] {
	case 0:
	case 1:
		return []string{node}
	default:
		return nil
	}
	visited[node] = 1
	for _, n := range adj[node] {
		if path := tb.detectCycle(n, adj, visited); path != nil {
			if len(path) > 1 && path[0] == path[len(path)-1] {
				return path
			}
			return append(path, node)
		}
	}
	visited[node] = 2
	return nil
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

	if err := checkBoxInputName(bd.box, bd.name, inputName); err != nil {
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

	// Setting InputName "output" prevents names of boxes from accidentally being leaked.
	edge := dataflowEdge{refname, sd.name, "output"}

	// check if this edge already exists
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

/**************************************************/

const (
	topologyBuilderAlreadyCalledBuildMsg = "this topology builder has already built the topology"
)

package bql

import (
	"fmt"
	"math/rand"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strconv"
)

type TopologyBuilder struct {
	stb      core.StaticTopologyBuilder
	sinkDecl map[string]core.SinkDeclarer
}

func NewTopologyBuilder() *TopologyBuilder {
	tb := TopologyBuilder{
		core.NewDefaultStaticTopologyBuilder(),
		map[string]core.SinkDeclarer{},
	}
	return &tb
}

func (tb *TopologyBuilder) BQL(s string) error {
	p := parser.NewBQLParser()
	_stmt, err := p.ParseStmt(s)
	if err != nil {
		return err
	}
	return tb.processStmt(_stmt)
}

func (tb *TopologyBuilder) Build() (core.StaticTopology, error) {
	return tb.stb.Build()
}

func (tb *TopologyBuilder) processStmt(_stmt interface{}) error {
	// check the type of statement
	switch stmt := _stmt.(type) {

	case parser.CreateSourceStmt:
		// we create a source with an internal name
		// since BQL-wise it is not allowed to directly
		// act on a "source" (but only on a "stream" derived
		// from a "source")
		srcName := string(stmt.Name + "-src")
		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)
		// pick the type of source
		if stmt.Type == "dummy" {
			return tb.createDummySource(srcName, paramsMap)
		}
		return fmt.Errorf("unknown source type: %s", stmt.Type)

	case parser.CreateStreamFromSourceStmt:
		// we create a stream by inserting a forwarding box
		// behind the source created before
		f := func(ctx *core.Context, t *tuple.Tuple, s core.Writer) error {
			s.Write(ctx, t)
			return nil
		}
		b := core.BoxFunc(f)
		refSrcName := string(stmt.Source) + "-src"
		// if the referenced source does not exist, `decl`
		// will have an error
		decl := tb.stb.AddBox(stmt.Name, b).Input(refSrcName)
		return decl.Err()

	case parser.CreateStreamFromSourceExtStmt:
		// this is a shortcut version of CREATE SOURCE
		// and CREATE STREAM FROM SOURCE.

		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)
		if stmt.Type == "dummy" {
			return tb.createDummySource(string(stmt.Name), paramsMap)
		}
		return fmt.Errorf("unknown source type: %s", stmt.Type)

	case parser.CreateStreamAsSelectStmt:
		// insert a bqlBox that executes the SELECT statement
		outName := stmt.Name
		box := NewBqlBox(&stmt)
		// add all the referenced relations as named inputs
		decl := tb.stb.AddBox(outName, box)
		for _, rel := range stmt.Relations {
			decl = decl.NamedInput(rel.Name, rel.Name)
		}
		return decl.Err()

	case parser.CreateSinkStmt:
		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)
		// we insert a sink, but cannot connect it to
		// any streams yet, therefore we have to keep track
		// of the SinkDeclarer
		if stmt.Type == "collector" {
			return tb.createCollectorSink(string(stmt.Name), paramsMap)
		}
		return fmt.Errorf("unknown sink type: %s", stmt.Type)

	case parser.InsertIntoSelectStmt:
		// get the sink to add an input to
		decl, ok := tb.sinkDecl[string(stmt.Sink)]
		if !ok {
			return fmt.Errorf("there is no sink named '%s'", stmt.Sink)
		}
		// construct an intermediate box doing the SELECT computation.
		//   INSERT INTO sink SELECT a, b FROM c WHERE d
		// becomes
		//   CREATE STREAM (random_string) AS SELECT ISTREAM(a, b)
		//   FROM c [RANGE 1 TUPLES] WHERE d
		//  + a connection (random_string -> sink)
		tmpName := fmt.Sprintf("tmp-%v", rand.Int())
		tmpRel := parser.Relation{tmpName}
		tmpStmt := parser.CreateStreamAsSelectStmt{
			tmpRel,
			parser.EmitProjectionsAST{parser.Istream, stmt.ProjectionsAST},
			parser.WindowedFromAST{stmt.FromAST,
				parser.RangeAST{parser.NumericLiteral{1}, parser.Tuples}},
			stmt.FilterAST,
			stmt.GroupingAST,
			stmt.HavingAST,
		}
		err := tb.processStmt(tmpStmt)
		if err != nil {
			return err
		}
		// now connect the sink to that box
		decl = decl.Input(tmpName)
		tb.sinkDecl[string(stmt.Sink)] = decl
		return decl.Err()
	}

	return fmt.Errorf("statement of type %T is unimplemented", _stmt)
}

func (tb *TopologyBuilder) mkParamsMap(params []parser.SourceSinkParamAST) map[string]string {
	paramsMap := make(map[string]string, len(params))
	for _, kv := range params {
		paramsMap[string(kv.Key)] = string(kv.Value)
	}
	return paramsMap
}

func (tb *TopologyBuilder) createDummySource(srcName string, paramsMap map[string]string) error {
	numTuples := 4
	// check the given source parameters
	for key, value := range paramsMap {
		if key == "num" {
			numTuples64, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				msg := "num: cannot convert string '%s' into integer"
				return fmt.Errorf(msg, value)
			}
			numTuples = int(numTuples64)
		} else {
			return fmt.Errorf("unknown source parameter: %s", key)
		}
	}
	// create source and add it to the topology
	so := NewTupleEmitterSource(numTuples)
	decl := tb.stb.AddSource(srcName, so)
	return decl.Err()
}

func (tb *TopologyBuilder) createCollectorSink(sinkName string, paramsMap map[string]string) error {
	// check the given sink parameters
	for key, _ := range paramsMap {
		return fmt.Errorf("unknown sink parameter: %s", key)
	}
	si := &TupleCollectorSink{}
	decl := tb.stb.AddSink(sinkName, si)
	tb.sinkDecl[sinkName] = decl
	return decl.Err()
}

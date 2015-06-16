package bql

import (
	"fmt"
	"math/rand"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type TopologyBuilder struct {
	stb      core.StaticTopologyBuilder
	sinks    map[string]core.Sink
	sinkDecl map[string]core.SinkDeclarer
	Reg      udf.FunctionManager
}

func NewTopologyBuilder() *TopologyBuilder {
	// create topology
	tb := TopologyBuilder{
		core.NewDefaultStaticTopologyBuilder(),
		map[string]core.Sink{},
		map[string]core.SinkDeclarer{},
		udf.NewDefaultFunctionRegistry(),
	}
	return &tb
}

func (tb *TopologyBuilder) BQL(s string) error {
	p := parser.NewBQLParser()
	// execute all parsed statements
	_stmts, err := p.ParseStmts(s)
	if err != nil {
		return err
	}
	for _, _stmt := range _stmts {
		err := tb.processStmt(_stmt)
		if err != nil {
			return err
		}
	}
	return nil
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
		// check if we know whis type of source
		if creator, ok := LookupSourceType(stmt.Type); ok {
			// if so, try to create such a source
			source, err := creator(paramsMap)
			if err != nil {
				return err
			}
			decl := tb.stb.AddSource(srcName, source)
			return decl.Err()
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
		decl := tb.stb.AddBox(string(stmt.Name), b).Input(refSrcName)
		return decl.Err()

	case parser.CreateStreamFromSourceExtStmt:
		// this is a shortcut version of CREATE SOURCE
		// and CREATE STREAM FROM SOURCE.

		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)
		// check if we know whis type of source
		if creator, ok := LookupSourceType(stmt.Type); ok {
			// if so, try to create such a source
			source, err := creator(paramsMap)
			if err != nil {
				return err
			}
			decl := tb.stb.AddSource(string(stmt.Name), source)
			return decl.Err()
		}
		return fmt.Errorf("unknown source type: %s", stmt.Type)

	case parser.CreateStreamAsSelectStmt:
		// insert a bqlBox that executes the SELECT statement
		outName := string(stmt.Name)
		box := NewBQLBox(&stmt, tb.Reg)
		// add all the referenced relations as named inputs
		decl := tb.stb.AddBox(outName, box)
		for _, rel := range stmt.Relations {
			decl = decl.NamedInput(rel.Name, rel.Name)
		}
		return decl.Err()

	case parser.CreateSinkStmt:
		// load params into map for faster access
		paramsMap := tb.mkParamsMap(stmt.Params)

		// check if we know whis type of sink
		if creator, ok := LookupSinkType(stmt.Type); ok {
			// if so, try to create such a sink
			sink, err := creator(paramsMap)
			if err != nil {
				return err
			}
			// we insert a sink, but cannot connect it to
			// any streams yet, therefore we have to keep track
			// of the SinkDeclarer
			sinkName := string(stmt.Name)
			decl := tb.stb.AddSink(sinkName, sink)
			tb.sinkDecl[sinkName] = decl
			// also remember the sink itself for programmatic access
			tb.sinks[sinkName] = sink
			return decl.Err()
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
		newRels := make([]parser.AliasedStreamWindowAST, len(stmt.Relations))
		for i, from := range stmt.Relations {
			if from.Unit != parser.UnspecifiedIntervalUnit {
				err := fmt.Errorf("you cannot use a RANGE clause with an INSERT INTO " +
					"statement at the moment")
				return err
			} else {
				newRels[i] = from
				newRels[i].IntervalAST = parser.IntervalAST{
					parser.NumericLiteral{1}, parser.Tuples}
			}
		}
		tmpStmt := parser.CreateStreamAsSelectStmt{
			parser.StreamIdentifier(tmpName),
			stmt.EmitterAST,
			stmt.ProjectionsAST,
			parser.WindowedFromAST{newRels},
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

// GetSink can be used to get programmatic access to a sink created via
// a CREATE SINK statement.
func (tb *TopologyBuilder) GetSink(name string) (core.Sink, bool) {
	sink, ok := tb.sinks[name]
	return sink, ok
}

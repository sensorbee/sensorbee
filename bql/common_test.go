package bql

import (
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/tuple"
)

func newTestContext(config core.Configuration) *core.Context {
	return &core.Context{
		Logger:       core.NewConsolePrintLogger(),
		Config:       config,
		SharedStates: core.NewDefaultSharedStateRegistry(),
	}
}

func newTestDynamicTopology() core.DynamicTopology {
	config := core.Configuration{TupleTraceEnabled: 0}
	ctx := newTestContext(config)
	return core.NewDefaultDynamicTopology(ctx, "testTopology")
}

func addBQLToTopology(tb *TopologyBuilder, bql string) error {
	p := parser.NewBQLParser()
	// execute all parsed statements
	stmts, err := p.ParseStmts(bql)
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		_, err := tb.AddStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

type dummyUDS struct {
	num int64
}

func newDummyUDS(ctx *core.Context, params tuple.Map) (core.SharedState, error) {
	s := &dummyUDS{}
	if v, ok := params["num"]; ok {
		if n, err := tuple.ToInt(v); err != nil {
			return nil, err
		} else {
			s.num = n
		}
	}
	return s, nil
}

func (s *dummyUDS) TypeName() string {
	return "dummy_uds"
}

func (s *dummyUDS) Init(ctx *core.Context) error {
	return nil
}

func (s *dummyUDS) Write(ctx *core.Context, t *tuple.Tuple) error {
	return nil
}

func (s *dummyUDS) Terminate(ctx *core.Context) error {
	return nil
}

func init() {
	if err := udf.RegisterGlobalUDSCreator("dummy_uds", udf.UDSCreatorFunc(newDummyUDS)); err != nil {
		panic(err)
	}
}

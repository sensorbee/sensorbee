package bql_test

import (
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/io/dummy_source"
	"pfi/sensorbee/sensorbee/tuple"
	"sync"
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

func addBQLToTopology(tb *bql.TopologyBuilder, bql string) error {
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

type blockingSource struct {
	wg        sync.WaitGroup
	s         core.Source
	unblocked bool
}

// TODO: Remove blockingSource after supporting pause&resume of sources in BQL.

var (
	_ core.ResumableNode = &blockingSource{}
)

func (s *blockingSource) GenerateStream(ctx *core.Context, w core.Writer) error {
	s.wg.Wait()
	return s.s.GenerateStream(ctx, w)
}

func (s *blockingSource) Stop(ctx *core.Context) error {
	defer func() {
		if e := recover(); e != nil {
			panic(e)
		}
	}()
	return s.s.Stop(ctx)
}

func (s *blockingSource) Pause(ctx *core.Context) error {
	// TODO: this is temporary implementation
	return nil
}

func (s *blockingSource) Resume(ctx *core.Context) error {
	if s.unblocked {
		return nil
	}
	s.unblocked = true
	s.wg.Done()
	return nil
}

func init() {
	bql.RegisterSourceType("blocking_dummy",
		func(ctx *core.Context, params tuple.Map) (core.Source, error) {
			// TODO: remove dummy_source package and move the source to this package.
			s, err := dummy_source.CreateDummySource(ctx, params)
			if err != nil {
				return nil, err
			}
			bs := &blockingSource{
				s: s,
			}
			bs.wg.Add(1)
			return bs, nil
		})
}

func unblockSource(dt core.DynamicTopology, srcName string) error {
	// TODO: this is a temporary function. Remove this after provoding pause&resume in BQL.
	sn, err := dt.Source(srcName)
	if err != nil {
		return err
	}
	if err := sn.Pause(); err != nil {
		return err
	}
	return sn.Resume()
}

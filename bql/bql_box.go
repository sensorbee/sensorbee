package bql

import (
	"pfi/sensorbee/sensorbee/bql/execution"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"sync"
)

type bqlBox struct {
	// stmt is the BQL statement executed by this box
	stmt *parser.SelectStmt
	// reg holds functions that can be used in this box
	reg udf.FunctionRegistry
	// plan is the execution plan for the SELECT statement in there
	execPlan execution.ExecutionPlan
	// mutex protects access to shared state
	mutex sync.Mutex
}

func NewBQLBox(stmt *parser.SelectStmt, reg udf.FunctionRegistry) *bqlBox {
	return &bqlBox{stmt: stmt, reg: reg}
}

func (b *bqlBox) Init(ctx *core.Context) error {
	// create the execution plan
	analyzedPlan, err := execution.Analyze(*b.stmt, b.reg)
	if err != nil {
		return err
	}
	optimizedPlan, err := analyzedPlan.LogicalOptimize()
	if err != nil {
		return err
	}
	b.execPlan, err = optimizedPlan.MakePhysicalPlan(b.reg)
	if err != nil {
		return err
	}
	return nil
}

func (b *bqlBox) Process(ctx *core.Context, t *core.Tuple, s core.Writer) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// feed tuple into plan
	resultData, err := b.execPlan.Process(t)
	if err != nil {
		return err
	}

	// emit result data as tuples
	for _, data := range resultData {
		tup := &core.Tuple{
			Data:          data,
			Timestamp:     t.Timestamp,
			ProcTimestamp: t.ProcTimestamp,
			BatchID:       t.BatchID,
		}
		if err := s.Write(ctx, tup); err != nil {
			return err
		}
	}
	return nil
}

func (b *bqlBox) Terminate(ctx *core.Context) error {
	// TODO cleanup
	return nil
}

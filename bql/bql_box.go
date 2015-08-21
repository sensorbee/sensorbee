package bql

import (
	"math/rand"
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
	execPlan execution.PhysicalPlan
	// mutex protects access to shared state
	mutex sync.Mutex
	// emitterLimit holds a positive value if this box should
	// stop emitting items after a certain number of items
	emitterLimit int64
	// emitterSampling holds a positive value if this box should only
	// emit a certain subset of items (defined by emitterSamplingType)
	emitterSampling int64
	// emitterSamplingType holds a value different from
	// parser.UnspecifiedSamplingType if output sampling is active
	emitterSamplingType parser.EmitterSamplingType
	// genCount holds the number of items generated so far
	// (i.e. computed by the underlying execution plan)
	genCount int64
	// emitCount holds the number of items emitted so far
	emitCount int64
	// removeMe is a function to remove this bqlBox from its
	// topology. A nil check must be done before calling.
	removeMe func()
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
	b.emitterLimit = analyzedPlan.EmitterLimit
	b.emitterSampling = analyzedPlan.EmitterSampling
	b.emitterSamplingType = analyzedPlan.EmitterSamplingType
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

	// deal with statements that have an emitter limit. in particular,
	// if we are already over the limit, exit here
	if b.emitterLimit >= 0 && b.emitCount >= b.emitterLimit {
		return nil
	}

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
		if len(t.Trace) != 0 {
			tup.Trace = make([]core.TraceEvent, len(t.Trace))
			copy(tup.Trace, t.Trace)
		}
		// decide if we should emit a tuple for this item
		shouldWriteTuple := true
		if b.emitterSamplingType == parser.CountBasedSampling {
			shouldWriteTuple = b.genCount%b.emitterSampling == 0
		} else if b.emitterSamplingType == parser.RandomizedSampling {
			shouldWriteTuple = rand.Int63n(100) < b.emitterSampling
		}
		// with 1,000,000 items per second, the counter below will overflow
		// after running for 292,471 years. probably ok.
		b.genCount += 1
		// write the tuple to the connected box
		if shouldWriteTuple {
			if err := s.Write(ctx, tup); err != nil {
				return err
			}
			b.emitCount += 1
		}
		// stop emitting if we have hit the limit
		if b.emitterLimit >= 0 && b.emitCount >= b.emitterLimit {
			break
		}
	}

	// remove this box if we are over the limit
	if b.emitterLimit >= 0 && b.emitCount >= b.emitterLimit {
		if b.removeMe != nil {
			b.removeMe()
			// don't call twice
			b.removeMe = nil
		}
	}

	return nil
}

func (b *bqlBox) Terminate(ctx *core.Context) error {
	// TODO cleanup
	return nil
}

package execution

import (
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type defaultSelectExecutionPlan struct {
	streamRelationStreamExecutionPlan
}

// CanBuildDefaultSelectExecutionPlan checks whether the given statement
// allows to use an defaultSelectExecutionPlan.
func CanBuildDefaultSelectExecutionPlan(lp *LogicalPlan, reg udf.FunctionRegistry) bool {
	return !lp.GroupingStmt
}

// defaultSelectExecutionPlan is a very simple plan that follows the
// theoretical processing model. It does not support aggregration.
//
// After each tuple arrives,
// - compute the contents of the current window using the
//   specified window size/type,
// - perform a SELECT query on that data,
// - compute the data that need to be emitted by comparison with
//   the previous run's results.
func NewDefaultSelectExecutionPlan(lp *LogicalPlan, reg udf.FunctionRegistry) (ExecutionPlan, error) {
	underlying, err := newStreamRelationStreamExecutionPlan(lp, reg)
	if err != nil {
		return nil, err
	}
	return &defaultSelectExecutionPlan{
		*underlying,
	}, nil
}

// Process takes an input tuple and returns a slice of Map values that
// correspond to the results of the query represented by this execution
// plan. Note that the order of items in the returned slice is undefined
// and cannot be relied on.
func (ep *defaultSelectExecutionPlan) Process(input *core.Tuple) ([]data.Map, error) {
	return ep.process(input, ep.performQueryOnBuffer)
}

// performQueryOnBuffer computes the projections of a SELECT query on the data
// stored in `ep.filteredInputRows`. The query results (which is a set of
// data.Value, not core.Tuple) is stored in ep.curResults. The data
// that was stored in ep.curResults before this method was called is
// moved to ep.prevResults. Note that the order of values in ep.curResults
// is undefined.
//
// In case of an error the contents of ep.curResults will still be
// the same as before the call (so that the next run performs as
// if no error had happened), but the contents of ep.curResults are
// undefined.
func (ep *defaultSelectExecutionPlan) performQueryOnBuffer() error {
	// reuse the allocated memory
	output := ep.prevResults[0:0]
	// remember the previous results
	ep.prevResults = ep.curResults

	rollback := func() {
		// NB. ep.prevResults currently points to an slice with
		//     results from the previous run. ep.curResults points
		//     to the same slice. output points to a different slice
		//     with a different underlying array.
		//     in the next run, output will be reusing the underlying
		//     storage of the current ep.prevResults to hold results.
		//     therefore when we leave this function we must make
		//     sure that ep.prevResults and ep.curResults have
		//     different underlying arrays or ISTREAM/DSTREAM will
		//     return wrong results.
		ep.prevResults = output
	}

	// function to compute the projection values and store
	// the result in the `output` slice
	evalItem := func(d data.Map) error {
		// otherwise, compute all the expressions
		result := data.Map(make(map[string]data.Value, len(ep.projections)))
		for _, proj := range ep.projections {
			value, err := proj.evaluator.Eval(d)
			if err != nil {
				return err
			}
			if err := assignOutputValue(result, proj.alias, value); err != nil {
				return err
			}
		}
		output = append(output, result)
		return nil
	}

	// compute the output for each item in ep.filteredInputRows
	for e := ep.filteredInputRows.Front(); e != nil; e = e.Next() {
		item := e.Value.(*data.Map)
		if err := evalItem(*item); err != nil {
			rollback()
			return err
		}
	}

	ep.curResults = output
	return nil
}

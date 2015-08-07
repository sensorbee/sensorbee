package execution

import (
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type filterPlan struct {
	commonExecutionPlan
	relAlias string
}

// CanBuildFilterPlan checks whether the given statement
// allows to use a filterPlan.
func CanBuildFilterPlan(lp *LogicalPlan, reg udf.FunctionRegistry) bool {
	if len(lp.Relations) != 1 {
		return false
	}
	return !lp.GroupingStmt &&
		lp.EmitterType == parser.Rstream &&
		lp.Relations[0].Unit == parser.Tuples &&
		lp.Relations[0].Value == 1
}

// filterPlan is a fast and simple plan for the case where the
// BQL statement has an Rstream emitter, a [RANGE 1 TUPLES] and (maybe)
// a WHERE clause (no GROUP BY/aggregate functions). In that case we can
// perform the check with less memory and faster than the default plan.
func NewFilterPlan(lp *LogicalPlan, reg udf.FunctionRegistry) (ExecutionPlan, error) {
	// prepare projection components
	projs, err := prepareProjections(lp.Projections, reg)
	if err != nil {
		return nil, err
	}
	// compute evaluator for the filter
	filter, err := prepareFilter(lp.Filter, reg)
	if err != nil {
		return nil, err
	}
	return &filterPlan{commonExecutionPlan{
		projections: projs,
		filter:      filter,
	}, lp.Relations[0].Alias}, nil
}

func (ep *filterPlan) Process(input *core.Tuple) ([]data.Map, error) {
	// nest the data in a one-element map using the alias as the key
	input.Data = data.Map{ep.relAlias: input.Data}
	setMetadata(input.Data, ep.relAlias, input)

	// add the information accessed by the now() function
	// to each item
	input.Data[":meta:NOW"] = data.Timestamp(time.Now().In(time.UTC))

	// evaluate filter condition and convert to bool
	if ep.filter != nil {
		filterResult, err := ep.filter.Eval(input.Data)
		if err != nil {
			return nil, err
		}
		filterResultBool, err := data.ToBool(filterResult)
		if err != nil {
			return nil, err
		}
		// if it evaluated to false, do not further process this tuple
		if !filterResultBool {
			return nil, nil
		}
	}
	// otherwise, compute all the expressions
	result := data.Map(make(map[string]data.Value, len(ep.projections)))
	for _, proj := range ep.projections {
		value, err := proj.evaluator.Eval(input.Data)
		if err != nil {
			return nil, err
		}
		if err := assignOutputValue(result, proj.alias, value); err != nil {
			return nil, err
		}
	}

	return []data.Map{result}, nil
}

package execution

import (
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type filterIstreamPlan struct {
	commonExecutionPlan
}

// CanBuildFilterIstreamPlan checks whether the given statement
// allows to use a filterIstreamPlan.
func CanBuildFilterIstreamPlan(lp *LogicalPlan, reg udf.FunctionRegistry) bool {
	// TODO check that there are no aggregate functions
	return len(lp.Relations) == 1 &&
		len(lp.GroupList) == 0 &&
		lp.Having == nil &&
		lp.EmitterType == parser.Istream
}

// filterIstreamPlan is a fast and simple plan for the case where
// the BQL statement has an Istream emitter and only a WHERE clause
// (no GROUP BY/aggregate functions). In that case we do not need
// an internal buffer, but we can just evaluate filter on projections
// on the incoming tuple and emit it right away if the filter matches.
func NewFilterIstreamPlan(lp *LogicalPlan, reg udf.FunctionRegistry) (ExecutionPlan, error) {
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
	return &filterIstreamPlan{commonExecutionPlan{
		projections: projs,
		filter:      filter,
	}}, nil
}

func (ep *filterIstreamPlan) Process(input *tuple.Tuple) ([]tuple.Map, error) {
	// evaluate filter condition and convert to bool
	filterResult, err := ep.filter.Eval(input.Data)
	if err != nil {
		return nil, err
	}
	filterResultBool, err := tuple.ToBool(filterResult)
	if err != nil {
		return nil, err
	}
	// if it evaluated to false, do not further process this tuple
	if !filterResultBool {
		return []tuple.Map{}, nil
	}
	// otherwise, compute all the expressions
	result := tuple.Map(make(map[string]tuple.Value, len(ep.projections)))
	for _, proj := range ep.projections {
		value, err := proj.evaluator.Eval(input.Data)
		if err != nil {
			return nil, err
		}
		if err := assignOutputValue(result, proj.alias, value); err != nil {
			return nil, err
		}
	}
	return []tuple.Map{result}, nil
}

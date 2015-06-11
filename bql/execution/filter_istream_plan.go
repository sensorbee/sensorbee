package execution

import (
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
)

type filterIstreamPlan struct {
	commonExecutionPlan
	relAlias    string
	windowSize  int64
	prevResults []tuple.Map
	procItems   int64
}

// CanBuildFilterIstreamPlan checks whether the given statement
// allows to use a filterIstreamPlan.
func CanBuildFilterIstreamPlan(lp *LogicalPlan, reg udf.FunctionRegistry) bool {
	// TODO check that there are no aggregate functions
	return len(lp.Relations) == 1 &&
		len(lp.GroupList) == 0 &&
		lp.Having == nil &&
		lp.EmitterType == parser.Istream &&
		lp.Unit == parser.Tuples
}

// filterIstreamPlan is a fast and simple plan for the case where the
// BQL statement has an Istream emitter, a TUPLES range clause and only
// a WHERE clause (no GROUP BY/aggregate functions). In that case we can
// perform the check with less memory and faster than the default plan.
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
	}, lp.Relations[0].Alias, lp.Value, make([]tuple.Map, 0, lp.Value), 0}, nil
}

func (ep *filterIstreamPlan) Process(input *tuple.Tuple) ([]tuple.Map, error) {
	// remove the oldest item from prevResults after it has once
	// been filled completely
	if ep.procItems >= ep.windowSize {
		if len(ep.prevResults) >= 1 {
			ep.prevResults = ep.prevResults[1:]
		}
	} else {
		ep.procItems++
	}

	// nest the data in a one-element map using the alias as the key
	input.Data = tuple.Map{ep.relAlias: input.Data}

	// evaluate filter condition and convert to bool
	if ep.filter != nil {
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
			return nil, nil
		}
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
	// check if the item is in the previous results
	alreadyEmitted := false
	for _, item := range ep.prevResults {
		if reflect.DeepEqual(result, item) {
			alreadyEmitted = true
			break
		}
	}
	// in any way, append it the the previous results
	ep.prevResults = append(ep.prevResults, result)
	if alreadyEmitted {
		return nil, nil
	}
	return []tuple.Map{result}, nil
}

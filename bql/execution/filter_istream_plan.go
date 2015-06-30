package execution

import (
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"reflect"
)

type filterIstreamPlan struct {
	commonExecutionPlan
	relAlias    string
	windowSize  int64
	prevResults []data.Map
	procItems   int64
}

// CanBuildFilterIstreamPlan checks whether the given statement
// allows to use a filterIstreamPlan.
func CanBuildFilterIstreamPlan(lp *LogicalPlan, reg udf.FunctionRegistry) bool {
	if len(lp.Relations) != 1 {
		return false
	}
	rangeUnit := lp.Relations[0].Unit
	// TODO check that there are no aggregate functions
	return len(lp.GroupList) == 0 &&
		lp.Having == nil &&
		lp.EmitterType == parser.Istream &&
		rangeUnit == parser.Tuples
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
	rangeValue := lp.Relations[0].Value
	return &filterIstreamPlan{commonExecutionPlan{
		projections: projs,
		filter:      filter,
	}, lp.Relations[0].Alias, rangeValue, make([]data.Map, 0, rangeValue), 0}, nil
}

func (ep *filterIstreamPlan) Process(input *core.Tuple) ([]data.Map, error) {
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
	input.Data = data.Map{ep.relAlias: input.Data}
	setMetadata(input.Data, ep.relAlias, input)

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
	return []data.Map{result}, nil
}

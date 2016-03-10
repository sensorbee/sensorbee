package execution

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
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

// NewFilterPlan creates a fast and simple plan for the case where the
// BQL statement has an Rstream emitter, a [RANGE 1 TUPLES] and (maybe)
// a WHERE clause (no GROUP BY/aggregate functions). In that case we can
// perform the check with less memory and faster than the default plan.
func NewFilterPlan(lp *LogicalPlan, reg udf.FunctionRegistry) (PhysicalPlan, error) {
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
	d := data.Map{ep.relAlias: input.Data}
	setMetadata(d, ep.relAlias, input)

	// because this plan doesn't cache data, TFSharedData flag doesn't have to be set.

	// add the information accessed by the now() function
	// to each item
	d[":meta:NOW"] = data.Timestamp(time.Now().In(time.UTC))

	// evaluate filter condition and convert to bool
	if ep.filter != nil {
		filterResult, err := ep.filter.Eval(d)
		if err != nil {
			return nil, err
		}
		// a NULL value is definitely not "true", so since we
		// have only a binary decision, we should drop tuples
		// where the filter condition evaluates to NULL
		filterResultBool := false
		if filterResult.Type() != data.TypeNull {
			filterResultBool, err = data.AsBool(filterResult)
			if err != nil {
				return nil, err
			}
		}
		// if it evaluated to false, do not further process this tuple
		if !filterResultBool {
			return nil, nil
		}
	}
	// otherwise, compute all the expressions
	result := data.Map(make(map[string]data.Value, len(ep.projections)))
	for _, proj := range ep.projections {
		value, err := proj.evaluator.Eval(d)
		if err != nil {
			return nil, err
		}
		if err := assignOutputValue(result, proj.alias, proj.aliasPath, value); err != nil {
			return nil, err
		}
	}

	return []data.Map{result}, nil
}

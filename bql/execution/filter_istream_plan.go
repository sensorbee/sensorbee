package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type filterIstreamPlan struct {
	// TODO turn this into a list of structs to ensure same length
	// colHeaders stores the names of the result columns.
	colHeaders []string
	// selectors stores the evaluators of the result columns.
	selectors []Evaluator
	// filter stores the evaluator of the filter condition,
	// or nil if there is no WHERE clause.
	filter Evaluator
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
	projs := make([]Evaluator, len(lp.Projections))
	colHeaders := make([]string, len(lp.Projections))
	for i, proj := range lp.Projections {
		// compute evaluators for each column
		plan, err := ExpressionToEvaluator(proj, reg)
		if err != nil {
			return nil, err
		}
		projs[i] = plan
		// compute column name
		colHeader := fmt.Sprintf("col_%v", i+1)
		switch projType := proj.(type) {
		case parser.ColumnName:
			colHeader = projType.Name
		case parser.FuncAppAST:
			colHeader = string(projType.Function)
		}
		colHeaders[i] = colHeader
	}
	// compute evaluator for the filter
	var filter Evaluator
	if lp.Filter != nil {
		f, err := ExpressionToEvaluator(lp.Filter, reg)
		if err != nil {
			return nil, err
		}
		filter = f
	}
	return &filterIstreamPlan{
		colHeaders: colHeaders,
		selectors:  projs,
		filter:     filter,
	}, nil
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
	result := tuple.Map(make(map[string]tuple.Value, len(ep.colHeaders)))
	for idx, selector := range ep.selectors {
		colName := ep.colHeaders[idx]
		value, err := selector.Eval(input.Data)
		if err != nil {
			return nil, err
		}
		result[colName] = value
	}
	return []tuple.Map{result}, nil
}

package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
)

// defaultSelectExecutionPlan can only process statements consisting of
// expressions in the SELECT and WHERE clause; it is not able to JOIN
// or GROUP BY. Therefore the processing is very simple: First compute
// the tuples for which the filter expression evaluates to true, then
// evaluate all the expressions in the SELECT clause.
type defaultSelectExecutionPlan struct {
	colHeaders []string
	selectors  []Evaluator
	filter     Evaluator
}

func NewDefaultSelectExecutionPlan(sProjections []interface{},
	sFilter interface{}, reg udf.FunctionRegistry) (ExecutionPlan, error) {
	projs := make([]Evaluator, len(sProjections))
	colHeaders := make([]string, len(sProjections))
	for i, proj := range sProjections {
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
	if sFilter != nil {
		f, err := ExpressionToEvaluator(sFilter, reg)
		if err != nil {
			return nil, err
		}
		filter = f
	}

	return &defaultSelectExecutionPlan{colHeaders, projs, filter}, nil
}

func (ep *defaultSelectExecutionPlan) Run(input []*tuple.Tuple) ([]tuple.Map, error) {
	if len(ep.colHeaders) != len(ep.selectors) {
		return nil, fmt.Errorf("number of columns (%v) doesn't match selectors (%v)",
			len(ep.colHeaders), len(ep.selectors))
	}
	output := []tuple.Map{}
	for _, t := range input {
		// evaluate filter condition and convert to bool
		filterResult, err := ep.filter.Eval(t.Data)
		if err != nil {
			return nil, err
		}
		filterResultBool, err := tuple.ToBool(filterResult)
		if err != nil {
			return nil, err
		}
		// if it evaluated to false, do not further process this tuple
		if !filterResultBool {
			continue
		}
		// otherwise, compute all the expressions
		result := tuple.Map(make(map[string]tuple.Value, len(ep.colHeaders)))
		for idx, selector := range ep.selectors {
			colName := ep.colHeaders[idx]
			value, err := selector.Eval(t.Data)
			if err != nil {
				return nil, err
			}
			result[colName] = value
		}
		output = append(output, result)
	}
	return output, nil
}

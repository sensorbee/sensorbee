package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
)

type colDesc struct {
	alias     string
	evaluator Evaluator
}

type commonExecutionPlan struct {
	// TODO turn this into a list of structs to ensure same length
	projections []colDesc
	// filter stores the evaluator of the filter condition,
	// or nil if there is no WHERE clause.
	filter Evaluator
}

func prepareProjections(projections []interface{}, reg udf.FunctionRegistry) ([]colDesc, error) {
	output := make([]colDesc, len(projections))
	for i, proj := range projections {
		// compute evaluators for each column
		plan, err := ExpressionToEvaluator(proj, reg)
		if err != nil {
			return nil, err
		}
		// compute column name
		colHeader := fmt.Sprintf("col_%v", i+1)
		switch projType := proj.(type) {
		case parser.ColumnName:
			colHeader = projType.Name
		case parser.AliasAST:
			colHeader = projType.Alias
		case parser.FuncAppAST:
			colHeader = string(projType.Function)
		}
		output[i] = colDesc{colHeader, plan}
	}
	return output, nil
}

func prepareFilter(filter interface{}, reg udf.FunctionRegistry) (Evaluator, error) {
	if filter != nil {
		return ExpressionToEvaluator(filter, reg)
	}
	return nil, nil
}

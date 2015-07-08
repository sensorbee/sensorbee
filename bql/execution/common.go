package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type colDesc struct {
	alias     string
	evaluator Evaluator
}

type commonExecutionPlan struct {
	projections []colDesc
	// filter stores the evaluator of the filter condition,
	// or nil if there is no WHERE clause.
	filter Evaluator
}

func prepareProjections(projections []parser.Expression, reg udf.FunctionRegistry) ([]colDesc, error) {
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
		case parser.RowMeta:
			if projType.MetaType == parser.TimestampMeta {
				colHeader = "ts"
			}
		case parser.RowValue:
			colHeader = projType.Column
		case parser.AliasAST:
			colHeader = projType.Alias
		case parser.FuncAppAST:
			colHeader = string(projType.Function)
		case parser.Wildcard:
			// The wildcard projection (without AS) is very special in that
			// it is the only case where the BQL user does not determine
			// the output key names (implicitly or explicitly). The
			// Evaluator interface is designed such that Evaluator
			// has 100% control over the returned value, but 0% control
			// over how it is named, therefore the wildcard evaluation
			// requires handling in multiple locations.
			// As a workaround, we will return the complete Map from
			// the wildcard Evaluator, nest it under a hard-coded key
			// called "*" and flatten them later (this is done correctly
			// by the assignOutputValue function).
			// Note that if it is desired at some point that there are
			// more evaluators with that behavior, we should change the
			// Evaluator.Eval interface.
			colHeader = "*"
		}
		output[i] = colDesc{colHeader, plan}
	}
	return output, nil
}

func prepareFilter(filter parser.Expression, reg udf.FunctionRegistry) (Evaluator, error) {
	if filter != nil {
		return ExpressionToEvaluator(filter, reg)
	}
	return nil, nil
}

// setMetadata adds the metadata contained in the given Tuple into the
// given Map with a key constructed using the given alias string. For example,
//   {"alias": {"col_1": ..., "col_2": ...}}
// is transformed into
//   {"alias": {"col_1": ..., "col_2": ...},
//    "alias:meta:TS": (timestamp of the given tuple)}
// so that the Evaluator created from a RowMeta AST struct works correctly.
func setMetadata(where data.Map, alias string, t *core.Tuple) {
	// this key format is also used in ExpressionToEvaluator()
	tsKey := fmt.Sprintf("%s:meta:%s", alias, parser.TimestampMeta)
	where[tsKey] = data.Timestamp(t.Timestamp)
}

// assignOutputValue writes the given Value `value` to the given
// Map `where` using the given key.
// If the key is "*" and the value is itself a Map, its contents
// will be "pulled up" and directly assigned to `where` (not
// nested) in order to provide wildcard functionality.
func assignOutputValue(where data.Map, key string, value data.Value) error {
	if key == "*" {
		valMap, err := data.AsMap(value)
		if err != nil {
			return err
		}
		for k, v := range valMap {
			where[k] = v
		}
	} else {
		return where.Set(key, value)
	}
	return nil
}

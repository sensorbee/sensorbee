package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type aliasedEvaluator struct {
	alias     string
	evaluator Evaluator
}

type commonExecutionPlan struct {
	projections []aliasedEvaluator
	// filter stores the evaluator of the filter condition,
	// or nil if there is no WHERE clause.
	filter Evaluator
}

func prepareProjections(projections []aliasedExpression, reg udf.FunctionRegistry) ([]aliasedEvaluator, error) {
	output := make([]aliasedEvaluator, len(projections))
	for i, proj := range projections {
		// compute evaluators for each column
		plan, err := ExpressionToEvaluator(proj.expr, reg)
		if err != nil {
			return nil, err
		}
		output[i] = aliasedEvaluator{proj.alias, plan}
	}
	return output, nil
}

func prepareFilter(filter FlatExpression, reg udf.FunctionRegistry) (Evaluator, error) {
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

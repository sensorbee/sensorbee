package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type aliasedEvaluator struct {
	alias        string
	aliasPath    data.Path
	evaluator    Evaluator
	hasAggregate bool
	aggrEvals    map[string]Evaluator
}

type commonExecutionPlan struct {
	projections []aliasedEvaluator
	groupList   []Evaluator
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
		containsAggregate := len(proj.aggrInputs) > 0
		// compute evaluators for the aggregate inputs
		var aggrEvals map[string]Evaluator
		if containsAggregate {
			aggrEvals = make(map[string]Evaluator, len(proj.aggrInputs))
			for key, aggrInput := range proj.aggrInputs {
				aggrEval, err := ExpressionToEvaluator(aggrInput, reg)
				if err != nil {
					return nil, err
				}
				aggrEvals[key] = aggrEval
			}
		}
		var path data.Path
		if proj.alias != "*" && proj.alias != ":having:" {
			path, err = data.CompilePath(proj.alias)
			if err != nil {
				return nil, err
			}
		}
		output[i] = aliasedEvaluator{proj.alias, path, plan, containsAggregate, aggrEvals}
	}
	return output, nil
}

func prepareFilter(filter FlatExpression, reg udf.FunctionRegistry) (Evaluator, error) {
	if filter != nil {
		return ExpressionToEvaluator(filter, reg)
	}
	return nil, nil
}

func prepareGroupList(groupList []FlatExpression, reg udf.FunctionRegistry) ([]Evaluator, error) {
	output := make([]Evaluator, len(groupList))
	for i, expr := range groupList {
		// compute evaluators for each expression
		plan, err := ExpressionToEvaluator(expr, reg)
		if err != nil {
			return nil, err
		}
		output[i] = plan
	}
	return output, nil
}

// setMetadata adds the metadata contained in the given Tuple into the
// given Map with a key constructed using the given alias string. For example,
//   {"alias": {"col_1": ..., "col_2": ...}}
// is transformed into
//   {"alias": {"col_1": ..., "col_2": ...},
//    "alias:meta:TS": (timestamp of the given tuple)}
// so that the Evaluator created from a parser.RowMeta AST struct works correctly.
func setMetadata(where data.Map, alias string, t *core.Tuple) {
	// this key format is also used in ExpressionToEvaluator()
	tsKey := fmt.Sprintf("%s:meta:%s", alias, parser.TimestampMeta)
	where[tsKey] = data.Timestamp(t.Timestamp)
}

// assignOutputValue writes the given Value `value` to the given
// Map `where` using the given `path`.
// If the `key` is "*" and the value is itself a Map, its contents
// will be "pulled up" and directly assigned to `where` (not
// nested) in order to provide wildcard functionality.
func assignOutputValue(where data.Map, key string, path data.Path, value data.Value) error {
	// TODO this wildcard thing should be solved differently, we
	//      should get rid of the legacy `key` parameter
	if key == "*" {
		valMap, err := data.AsMap(value)
		if err != nil {
			return err
		}
		for k, v := range valMap {
			where[k] = v
		}
	} else {
		if path == nil {
			return fmt.Errorf("trying to assign output column '%s', path was nil", key)
		}
		return where.Set(path, value)
	}
	return nil
}

package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
)

/*
The functions in this file transform an AST as returned by
`parser.bqlParser.ParseStmt()` into a physical execution plan
that can be used to query a set of Tuples. This works similar
to the way it is done in Spark SQL as outlined on
  https://databricks.com/blog/2015/04/13/deep-dive-into-spark-sqls-catalyst-optimizer.html
in three phases:
- Analyze
- LogicalOptimize
- MakePhysicalPlan

The ExecutionPlan that is returned has a method
  Run(input []*tuple.Tuple) ([]tuple.Map, error)
which can be used to run the query and obtain its results.
*/

type LogicalPlan struct {
	parser.ProjectionsAST
	parser.FilterAST
}

type ExecutionPlan interface {
	Run(input []*tuple.Tuple) ([]tuple.Map, error)
}

func Analyze(s *parser.CreateStreamAsSelectStmt) (*LogicalPlan, error) {
	/*
	   In Spark, this does the following:

	   > - Looking up relations by name from the catalog.
	   > - Mapping named attributes, such as col, to the input provided
	   >   given operator’s children.
	   > - Determining which attributes refer to the same value to give
	   >   them a unique ID (which later allows optimization of expressions
	   >   such as col = col).
	   > - Propagating and coercing types through expressions: for
	   >   example, we cannot know the return type of 1 + col until we
	   >   have resolved col and possibly casted its subexpressions to a
	   >   compatible types.
	*/
	if len(s.Relations) != 1 {
		return nil, fmt.Errorf("Using multiple relations is not implemented yet")
	}
	return &LogicalPlan{
		s.ProjectionsAST,
		s.FilterAST,
	}, nil
}

func (lp *LogicalPlan) LogicalOptimize() (*LogicalPlan, error) {
	/*
	   In Spark, this does the following:

	   > These include constant folding, predicate pushdown, projection
	   > pruning, null propagation, Boolean expression simplification,
	   > and other rules.
	*/
	return lp, nil
}

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

func (lp *LogicalPlan) MakePhysicalPlan(reg udf.FunctionRegistry) (ExecutionPlan, error) {
	/*
	   In Spark, this does the following:

	   > In the physical planning phase, Spark SQL takes a logical plan
	   > and generates one or more physical plans, using physical operators
	   > that match the Spark execution engine.
	*/
	projs := make([]Evaluator, len(lp.Projections))
	colHeaders := make([]string, len(lp.Projections))
	for i, proj := range lp.Projections {
		plan, err := ExpressionToEvaluator(proj, reg)
		if err != nil {
			return nil, err
		}
		projs[i] = plan
		colHeader := fmt.Sprintf("col_%v", i+1)
		switch projType := proj.(type) {
		case parser.ColumnName:
			colHeader = projType.Name

		}
		colHeaders[i] = colHeader
	}
	var filter Evaluator
	if lp.Filter != nil {
		f, err := ExpressionToEvaluator(lp.Filter, reg)
		if err != nil {
			return nil, err
		}
		filter = f
	}
	return &defaultSelectExecutionPlan{colHeaders, projs, filter}, nil
}

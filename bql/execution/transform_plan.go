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

func (lp *LogicalPlan) MakePhysicalPlan(reg udf.FunctionRegistry) (ExecutionPlan, error) {
	/*
	   In Spark, this does the following:

	   > In the physical planning phase, Spark SQL takes a logical plan
	   > and generates one or more physical plans, using physical operators
	   > that match the Spark execution engine.
	*/
	return NewDefaultSelectExecutionPlan(lp.Projections, lp.Filter, reg)
}

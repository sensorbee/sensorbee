package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
)

/*
The functions in this file transform an AST as returned by
`parser.bqlParser.ParseStmt()` into a physical execution plan.
This works similar to the way it is done in Spark SQL as
outlined on
  https://databricks.com/blog/2015/04/13/deep-dive-into-spark-sqls-catalyst-optimizer.html
in three phases:
- Analyze
- LogicalOptimize
- MakePhysicalPlan
*/

type LogicalPlan struct {
	parser.EmitProjectionsAST
	parser.WindowedFromAST
	parser.FilterAST
	parser.GroupingAST
	parser.HavingAST
}

// ExecutionPlan is a physical interface that is capable of
// computing the data that needs to be emitted into an output
// stream when a new tuple arrives in the input stream.
type ExecutionPlan interface {
	// Process must be called whenever a new tuple arrives in
	// the input stream. It will return a list of tuple.Map
	// items where each of these items is to be emitted as
	// a tuple. It is the caller's task to create those tuples
	// and set appropriate meta information such as timestamps.
	//
	// NB. Process is not thread-safe, i.e., it must be called in
	// a single-threaded context.
	Process(input *tuple.Tuple) ([]tuple.Map, error)
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
	return &LogicalPlan{
		s.EmitProjectionsAST,
		s.WindowedFromAST,
		s.FilterAST,
		s.GroupingAST,
		s.HavingAST,
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
	if CanBuildFilterIstreamPlan(lp, reg) {
		return NewFilterIstreamPlan(lp, reg)
	} else if CanBuildDefaultSelectExecutionPlan(lp, reg) {
		return NewDefaultSelectExecutionPlan(lp, reg)
	}
	return nil, fmt.Errorf("no plan can deal with such a statement")
}

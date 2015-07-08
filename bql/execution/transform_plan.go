package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"strings"
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
	parser.EmitterAST
	Projections []aliasedExpression
	parser.WindowedFromAST
	Filter    FlatExpression
	GroupList []FlatExpression
	parser.HavingAST
}

// ExecutionPlan is a physical interface that is capable of
// computing the data that needs to be emitted into an output
// stream when a new tuple arrives in the input stream.
type ExecutionPlan interface {
	// Process must be called whenever a new tuple arrives in
	// the input stream. It will return a list of data.Map
	// items where each of these items is to be emitted as
	// a tuple. It is the caller's task to create those tuples
	// and set appropriate meta information such as timestamps.
	//
	// NB. Process is not thread-safe, i.e., it must be called in
	// a single-threaded context.
	Process(input *core.Tuple) ([]data.Map, error)
}

func Analyze(s parser.CreateStreamAsSelectStmt) (*LogicalPlan, error) {
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

	if err := makeRelationAliases(&s); err != nil {
		return nil, err
	}

	if err := validateReferences(&s); err != nil {
		return nil, err
	}

	return flattenExpressions(&s)
}

// dummy implementation of aggregate function lookup
// TODO think of the real lookup mechanism
func isAggregateDummy(fName string) bool {
	if fName == "count" || fName == "udaf" {
		return true
	}
	return false
}

// flattenExpressions separates the aggregate and non-aggregate
// part in a statement and returns with an error if there are
// aggregates in structures that may not have some
func flattenExpressions(s *parser.CreateStreamAsSelectStmt) (*LogicalPlan, error) {
	flatProjExprs := make([]aliasedExpression, len(s.Projections))
	for i, expr := range s.Projections {
		// convert the parser Expression to a FlatExpression
		flatExpr, aggrs, err := ParserExprToMaybeAggregate(expr, isAggregateDummy)
		if err != nil {
			return nil, err
		}
		// compute column name
		colHeader := fmt.Sprintf("col_%v", i+1)
		switch projType := expr.(type) {
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
		flatProjExprs[i] = aliasedExpression{colHeader, flatExpr, aggrs}
	}

	var filterExpr FlatExpression
	if s.Filter != nil {
		filterFlatExpr, err := ParserExprToFlatExpr(s.Filter, isAggregateDummy)
		if err != nil {
			// return a prettier error message
			if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
				err = fmt.Errorf("aggregates not allowed in WHERE clause")
			}
			return nil, err
		}
		filterExpr = filterFlatExpr
	}

	flatGroupExprs := make([]FlatExpression, len(s.GroupList))
	for i, expr := range s.GroupList {
		// convert the parser Expression to a FlatExpression
		flatExpr, err := ParserExprToFlatExpr(expr, isAggregateDummy)
		if err != nil {
			// return a prettier error message
			if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
				err = fmt.Errorf("aggregates not allowed in GROUP BY clause")
			}
			return nil, err
		}
		flatGroupExprs[i] = flatExpr
	}

	return &LogicalPlan{
		s.EmitterAST,
		flatProjExprs,
		s.WindowedFromAST,
		filterExpr,
		flatGroupExprs,
		s.HavingAST,
	}, nil
}

// makeRelationAliases will assign an internal alias to every relation
// does not yet have one (given by the user). It will also detect if
// there is a conflict between aliases.
func makeRelationAliases(s *parser.CreateStreamAsSelectStmt) error {
	relNames := make(map[string]parser.AliasedStreamWindowAST, len(s.Relations))
	newRels := make([]parser.AliasedStreamWindowAST, len(s.Relations))
	for i, aliasedRel := range s.Relations {
		// if the relation does not yet have an internal alias, use
		// the relation name itself
		if aliasedRel.Alias == "" {
			aliasedRel.Alias = aliasedRel.Name
		}
		otherRel, exists := relNames[aliasedRel.Alias]
		if exists {
			return fmt.Errorf("cannot use relations '%s' and '%s' with the "+
				"same alias '%s'", aliasedRel.Name, otherRel.Name, aliasedRel.Alias)
		} else {
			relNames[aliasedRel.Alias] = aliasedRel
		}
		newRels[i] = aliasedRel
	}
	s.Relations = newRels
	return nil
}

// validateReferences checks if the references to input relations
// in SELECT, WHERE, GROUP BY and HAVING clauses of the given
// statement are matching the relations mentioned in the FROM
// clause.
func validateReferences(s *parser.CreateStreamAsSelectStmt) error {

	/* We want to check if we can access all relations properly.
	   If there is just one input relation, we ask that none of the
	   RowValue structs has a Relation string different from ""
	   (as in `SELECT col FROM stream`).
	   If there are multiple input relations, we ask that *all*
	   RowValue structs have a Relation string that matches one of
	   the input relations (as in `SELECT a.col, b.col FROM a, b`).
	*/

	// collect the referenced relations in SELECT, WHERE, GROUP BY clauses
	// and store them in the given map
	refRels := map[string]bool{}
	for _, proj := range s.Projections {
		for rel := range proj.ReferencedRelations() {
			refRels[rel] = true
		}
	}
	if s.Filter != nil {
		for rel := range s.Filter.ReferencedRelations() {
			refRels[rel] = true
		}
	}
	for _, group := range s.GroupList {
		for rel := range group.ReferencedRelations() {
			refRels[rel] = true
		}
	}

	// do the correctness check for SELECT, WHERE, GROUP BY clauses
	if len(s.Relations) == 0 {
		// Sample: SELECT a (no FROM clause)
		// this case should never happen due to parser setup
		return fmt.Errorf("need at least one relation to select from")

	} else if len(s.Relations) == 1 {
		inputRel := s.Relations[0].Alias
		if len(refRels) == 1 {
			// Sample: SELECT a FROM b // SELECT b.a FROM b
			// check if the one map item is either "" or the name of
			// the input relation
			for rel := range refRels {
				if rel != "" && rel != inputRel {
					err := fmt.Errorf("cannot refer to relation '%s' "+
						"when using only '%s'", rel, inputRel)
					return err
				}
			}
			// we need to make the references more explicit,
			// i.e., change all "" references to the name
			// of the only input relation
			newProjs := make([]parser.Expression, len(s.Projections))
			for i, proj := range s.Projections {
				newProjs[i] = proj.RenameReferencedRelation("", inputRel)
			}
			s.Projections = newProjs
			if s.Filter != nil {
				s.Filter = s.Filter.RenameReferencedRelation("", inputRel)
			}
			newGroup := make([]parser.Expression, len(s.GroupList))
			for i, group := range s.GroupList {
				newGroup[i] = group.RenameReferencedRelation("", inputRel)
			}
			s.GroupList = newGroup

		} else if len(refRels) > 1 {
			// Sample: SELECT a, b.a FROM b // SELECT b.a, x.a FROM b
			// this is an invalid statement
			failRels := make([]string, 0, len(refRels))
			for rel := range refRels {
				failRels = append(failRels, fmt.Sprintf("'%s'", rel))
			}
			failRelsStr := strings.Join(failRels, ", ")
			err := fmt.Errorf("cannot refer to relations %s "+
				"when using only '%s'", failRelsStr, inputRel)
			return err
		}
		// if we arrive here, the only referenced relation is valid or
		// we do not actually reference anything

	} else if len(s.Relations) > 1 {
		// Sample: SELECT b.a, c.d FROM b, c
		// check if all referenced relations are actually listed in FROM
		for rel := range refRels {
			found := false
			for _, inputRel := range s.Relations {
				if rel == inputRel.Alias {
					found = true
					break
				}
			}
			if !found {
				prettyRels := make([]string, 0, len(s.Relations))
				for _, rel_ := range s.Relations {
					prettyRels = append(prettyRels, fmt.Sprintf("'%s'", rel_.Alias))
				}
				prettyRelsStr := strings.Join(prettyRels, ", ")
				err := fmt.Errorf("cannot reference relation '%s' "+
					"when using input relations %v", rel, prettyRelsStr)
				return err
			}
		}
		// if we arrive here, all referenced relations exist in the
		// FROM clause -> OK
	}

	// HAVING is fundamentally different in that it does refer to
	// output columns, therefore must not contain references to input
	// relations
	havingRels := map[string]bool{}
	if s.Having != nil {
		for rel := range s.Having.ReferencedRelations() {
			havingRels[rel] = true
		}
	}
	for rel := range havingRels {
		if rel != "" {
			err := fmt.Errorf("cannot refer to input relation '%s' "+
				"from HAVING clause", rel)
			return err
		}
	}

	if s.EmitIntervals != nil {
		emitStreams := map[string]bool{}

		for _, emitSpec := range s.EmitIntervals {
			// check that no stream is used twice in the emitter clause
			if _, exists := emitStreams[emitSpec.Name]; exists {
				err := fmt.Errorf("the stream '%s' referenced in the %s "+
					"clause is used more than once", emitSpec.Name, s.EmitterType)
				return err
			}
			emitStreams[emitSpec.Name] = true

			// check that the emitter does not refer to any streams we don't know
			found := false
			for _, rel := range s.Relations {
				if emitSpec.Name == "*" {
					// just using `XSTREAM [EVERY k SECONDS/TUPLES]`
					// without a FROM clause will insert a "*" as the
					// stream name, which is always correct
					found = true
					break
				} else if emitSpec.Name == rel.Name {
					// when using `XSTREAM [EVERY k TUPLES FROM x]`
					// we must make sure that "x" is an input stream
					// we refer to (and *not* an alias, because that
					// is naming a relation, not a stream!)
					found = true
					break
				}
			}
			if !found {
				err := fmt.Errorf("the stream '%s' referenced in the %s "+
					"clause is unknown", emitSpec.Name, s.EmitterType)
				return err
			}
		}
	}

	return nil
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

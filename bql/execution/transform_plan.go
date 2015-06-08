package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
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

	return &LogicalPlan{
		s.EmitProjectionsAST,
		s.WindowedFromAST,
		s.FilterAST,
		s.GroupingAST,
		s.HavingAST,
	}, nil
}

// makeRelationAliases will assign an internal alias to every relation
// does not yet have one (given by the user). It will also detect if
// there is a conflict between aliases.
func makeRelationAliases(s *parser.CreateStreamAsSelectStmt) error {
	relNames := make(map[string]parser.AliasRelationAST, len(s.Relations))
	newRels := make([]parser.AliasRelationAST, len(s.Relations))
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

	// define a recursive function to collect all referenced relations
	// and store them in the given map
	var collectRels func(interface{}, map[string]bool) error
	collectRels = func(_expr interface{}, rels map[string]bool) error {
		if _expr == nil {
			return nil
		}
		switch expr := _expr.(type) {
		default:
			return fmt.Errorf("don't know how to collect referenced "+
				"relations from AST object %T", _expr)
		case parser.RowValue:
			rels[expr.Relation] = true
		case parser.AliasAST:
			return collectRels(expr.Expr, rels)
		case parser.NumericLiteral, parser.FloatLiteral, parser.BoolLiteral:
			// no referenced relations
		case parser.BinaryOpAST:
			if err := collectRels(expr.Left, rels); err != nil {
				return err
			}
			if err := collectRels(expr.Right, rels); err != nil {
				return err
			}
		case parser.FuncAppAST:
			for _, e := range expr.Expressions {
				if err := collectRels(e, rels); err != nil {
					return err
				}
			}
		case parser.Wildcard:
			// this is special, we can't use `rel.*` at the moment
		}
		return nil
	}

	// collect the referenced relations in SELECT, WHERE, GROUP BY clauses
	refRels := map[string]bool{}
	for _, proj := range s.Projections {
		if err := collectRels(proj, refRels); err != nil {
			return err
		}
	}
	if err := collectRels(s.Filter, refRels); err != nil {
		return err
	}
	for _, group := range s.GroupList {
		if err := collectRels(group, refRels); err != nil {
			return err
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
				for rel := range s.Relations {
					prettyRels = append(prettyRels, fmt.Sprintf("'%s'", rel))
				}
				prettyRelsStr := strings.Join(prettyRels, ", ")
				err := fmt.Errorf("cannot reference relation '%s' "+
					"when using input relations %s", rel, prettyRelsStr)
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
	if err := collectRels(s.Having, havingRels); err != nil {
		return err
	}
	for rel := range havingRels {
		if rel != "" {
			err := fmt.Errorf("cannot refer to input relation '%s' "+
				"from HAVING clause", rel)
			return err
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

package execution

import (
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"math"
	"regexp"
	"strings"
)

const (
	MaxRangeTuples   float64 = 1<<20 - 1
	MaxRangeSec      float64 = 60 * 60 * 24
	MaxRangeMillisec float64 = 60 * 60 * 24 * 1000
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

var (
	simpleColumnNameRe = regexp.MustCompile("^[a-z][a-zA-Z0-9_]*$")
)

// LogicalPlan represents a parsed and analyzed version of a SELECT
// statement. A LogicalPlan as returned by `Analyze` should not contain
// logical errors such as "... must appear in GROUP BY clause" etc.
type LogicalPlan struct {
	GroupingStmt        bool
	EmitterType         parser.Emitter
	EmitterLimit        int64
	EmitterSampling     float64
	EmitterSamplingType parser.EmitterSamplingType
	Projections         []aliasedExpression
	parser.WindowedFromAST
	Filter    FlatExpression
	GroupList []FlatExpression
	parser.HavingAST
}

// PhysicalPlan is a physical interface that is capable of
// computing the data that needs to be emitted into an output
// stream when a new tuple arrives in the input stream.
type PhysicalPlan interface {
	// Process must be called whenever a new tuple arrives in
	// the input stream. It will return a list of data.Map
	// items where each of these items is to be emitted as
	// a tuple. It is the caller's task to create those tuples
	// and set appropriate meta information such as timestamps.
	//
	// Process must NOT modify any field of the input tuple when its
	// core.TFShared flag is set. To modify the tuple, create a shallow copy of
	// it. Moreover, when Tuple.Data field is cached in the plan,
	// core.TFSharedData flag of the input tuple must be set and the plan must
	// not modify the Data.
	//
	// NB. Process is not thread-safe, i.e., it must be called in
	// a single-threaded context.
	Process(input *core.Tuple) ([]data.Map, error)
}

// Analyze checks the given SELECT statement for logical errors
// (references to unknown tables etc.) and creates a LogicalPlan
// that is internally consistent.
func Analyze(s parser.SelectStmt, reg udf.FunctionRegistry) (*LogicalPlan, error) {
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

	return flattenExpressions(&s, reg)
}

// isAggregateFunc is a helper function to check if one of
// the parameters of the given function is an aggregate
// parameter.
func isAggregateFunc(f udf.UDF, arity int) bool {
	agg := false
	for k := 0; k < arity; k++ {
		if f.IsAggregationParameter(k) {
			agg = true
		}
	}
	return agg
}

// flattenExpressions separates the aggregate and non-aggregate
// part in a statement and returns with an error if there are
// aggregates in structures that may not have some
func flattenExpressions(s *parser.SelectStmt, reg udf.FunctionRegistry) (*LogicalPlan, error) {
	// groupingMode is active when aggregate functions or the
	// GROUP BY clause are used
	groupingMode := false

	flatProjExprs := make([]aliasedExpression, len(s.Projections))
	numAggParams := 0
	for i, expr := range s.Projections {
		// convert the parser Expression to a FlatExpression
		flatExpr, aggrs, err := ParserExprToMaybeAggregate(expr, numAggParams, reg)
		numAggParams += len(aggrs)
		if err != nil {
			return nil, err
		}
		// remember if we have aggregates at all
		if len(aggrs) > 0 {
			groupingMode = true
		}
		// compute column name
		colHeader := fmt.Sprintf("col_%v", i+1)
		switch projType := expr.(type) {
		case parser.RowMeta:
			if projType.MetaType == parser.TimestampMeta {
				colHeader = "ts"
			}
		case parser.RowValue:
			// We can only use the column name as an alias if it is not
			// a complex JSON Path. For example, `SELECT a` will be treated
			// like `SELECT a AS a`, but for `SELECT a..b` we will have to
			// use the col_N form.
			if simpleColumnNameRe.MatchString(projType.Column) {
				colHeader = projType.Column
			}
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

	if s.Having != nil {
		// convert the parser Expression to a FlatExpression
		flatExpr, aggrs, err := ParserExprToMaybeAggregate(s.Having, numAggParams, reg)
		if err != nil {
			return nil, err
		}
		// use a special column name
		colHeader := ":having:"
		flatProjExprs = append(flatProjExprs,
			aliasedExpression{colHeader, flatExpr, aggrs})
		// we are definitely in grouping mode
		groupingMode = true
	}

	var filterExpr FlatExpression
	if s.Filter != nil {
		filterFlatExpr, err := ParserExprToFlatExpr(s.Filter, reg)
		if err != nil {
			// return a prettier error message
			if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
				err = fmt.Errorf("aggregates not allowed in WHERE clause")
			}
			return nil, err
		}
		filterExpr = filterFlatExpr
	}

	groupCols := make([]rowValue, len(s.GroupList))
	flatGroupExprs := make([]FlatExpression, len(s.GroupList))
	for i, expr := range s.GroupList {
		// convert the parser Expression to a FlatExpression
		flatExpr, err := ParserExprToFlatExpr(expr, reg)
		if err != nil {
			// return a prettier error message
			if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
				err = fmt.Errorf("aggregates not allowed in GROUP BY clause")
			}
			return nil, err
		}
		// at the moment we only support grouping by single columns,
		// not expressions
		col, ok := flatExpr.(rowValue)
		if !ok {
			err := fmt.Errorf("grouping by expressions is not supported yet")
			return nil, err
		}
		groupCols[i] = col
		flatGroupExprs[i] = flatExpr
	}
	groupingMode = groupingMode || len(flatGroupExprs) > 0

	// check if grouping is done correctly
	if groupingMode {
		for _, expr := range flatProjExprs {
			// the wildcard operator cannot be used with GROUP BY
			if expr.expr.ContainsWildcard() {
				err := fmt.Errorf("* cannot be used in GROUP BY statements")
				return nil, err
			}
			// all columns mentioned outside of an aggregate
			// function must be in the GROUP BY clause
			if rm, ok := expr.expr.(rowMeta); ok {
				err := fmt.Errorf("using metadata '%s' in GROUP BY statements is "+
					"not supported yet", rm.MetaType)
				return nil, err
			}
			usedCols := expr.expr.Columns()
			for _, usedCol := range usedCols {
				// look for this col in the GROUP BY clause
				mentioned := false
				for _, groupCol := range groupCols {
					if usedCol == groupCol {
						mentioned = true
						break
					}
				}
				if !mentioned {
					err := fmt.Errorf("column \"%s\" must appear in the GROUP BY "+
						"clause or be used in an aggregate function", usedCol.Repr())
					return nil, err
				}
			}
		}
	}

	// validate the emitter parameters
	emitLimit := int64(-1)
	emitSampling := float64(-1)
	emitSamplingType := parser.UnspecifiedSamplingType
	for _, opt := range s.EmitterAST.EmitterOptions {
		switch obj := opt.(type) {
		default:
			return nil, fmt.Errorf("unknown emitter options: %+v", obj)
		case parser.EmitterLimit:
			l := obj.Limit
			if l < 0 {
				return nil, fmt.Errorf("LIMIT parameter must have a "+
					"positive value, not %d", l)
			}
			emitLimit = l
		case parser.EmitterSampling:
			v := obj.Value
			switch obj.Type {
			default:
				return nil, fmt.Errorf("unknown emitter sampling type: %+v", obj.Type)
			case parser.CountBasedSampling:
				if v <= 0 {
					return nil, fmt.Errorf("EVERY parameter must have a "+
						"positive value, not %d", v)
				}
				if math.Trunc(v) != v {
					// this should be prevented by the parser, but better
					// check here again
					return nil, fmt.Errorf("EVERY parameter must have an "+
						"integral value for TUPLE, not %d", v)
				}
				emitSampling = v
			case parser.TimeBasedSampling:
				if v <= 0 {
					return nil, fmt.Errorf("EVERY parameter must have a "+
						"positive value, not %d", v)
				}
				emitSampling = v
			case parser.RandomizedSampling:
				if v < 0 || v > 100 {
					return nil, fmt.Errorf("SAMPLE parameter must have a "+
						"value between 0 and 100, not %d", v)
				}
				emitSampling = v / 100 // project to [0,1] interval
			}
			emitSamplingType = obj.Type
		}
	}

	return &LogicalPlan{
		groupingMode,
		s.EmitterAST.EmitterType,
		emitLimit,
		emitSampling,
		emitSamplingType,
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
func makeRelationAliases(s *parser.SelectStmt) error {
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
		}
		relNames[aliasedRel.Alias] = aliasedRel
		newRels[i] = aliasedRel
	}
	s.Relations = newRels
	return nil
}

// validateReferences checks if the references to input relations
// in SELECT, WHERE, GROUP BY and HAVING clauses of the given
// statement are matching the relations mentioned in the FROM
// clause.
func validateReferences(s *parser.SelectStmt) error {

	/* We want to check if we can access all relations properly.
	   If there is just one input relation, we ask that none of the
	   rowValue structs has a Relation string different from ""
	   (as in `SELECT col FROM stream`).
	   If there are multiple input relations, we ask that *all*
	   rowValue structs have a Relation string that matches one of
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
	if s.Having != nil {
		for rel := range s.Having.ReferencedRelations() {
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
			if s.Having != nil {
				s.Having = s.Having.RenameReferencedRelation("", inputRel)
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
				for _, inRel := range s.Relations {
					prettyRels = append(prettyRels, fmt.Sprintf("'%s'", inRel.Alias))
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

	for _, rel := range s.Relations {
		if rel.Value <= 0 {
			err := fmt.Errorf("number in RANGE clause must be positive, not %v", rel.Value)
			return err
		}
		if rel.Unit == parser.Tuples && math.Trunc(rel.Value) != rel.Value {
			// actually the parser should not allow fractional numbers,
			// but we check anyway
			err := fmt.Errorf("number in RANGE clause must be integral "+
				"for TUPLES, not %v", rel.Value)
			return err
		}
		switch rel.Unit {
		case parser.Tuples:
			if rel.Value > MaxRangeTuples {
				err := fmt.Errorf("RANGE value %d is too large for TUPLES (must be at most %d)",
					int64(rel.Value), int64(MaxRangeTuples))
				return err
			}
		case parser.Seconds:
			if rel.Value > MaxRangeSec {
				err := fmt.Errorf("RANGE value %v is too large for SECONDS (must be at most %d)",
					rel.Value, int64(MaxRangeSec))
				return err
			}
		case parser.Milliseconds:
			if rel.Value > MaxRangeMillisec {
				err := fmt.Errorf("RANGE value %v is too large for MILLISECONDS (must be at most %d)",
					rel.Value, int64(MaxRangeMillisec))
				return err
			}
		}
	}

	return nil
}

// LogicalOptimize does nothing at the moment. In the future, logical
// optimizations (evaluation of foldable terms etc.) can be added here.
func (lp *LogicalPlan) LogicalOptimize() (*LogicalPlan, error) {
	/*
	   In Spark, this does the following:

	   > These include constant folding, predicate pushdown, projection
	   > pruning, null propagation, Boolean expression simplification,
	   > and other rules.
	*/
	return lp, nil
}

// MakePhysicalPlan creates a physical execution plan that is able to
// deal with the statement under consideration.
func (lp *LogicalPlan) MakePhysicalPlan(reg udf.FunctionRegistry) (PhysicalPlan, error) {
	/*
	   In Spark, this does the following:

	   > In the physical planning phase, Spark SQL takes a logical plan
	   > and generates one or more physical plans, using physical operators
	   > that match the Spark execution engine.
	*/
	if CanBuildFilterPlan(lp, reg) {
		return NewFilterPlan(lp, reg)
	} else if CanBuildDefaultSelectExecutionPlan(lp, reg) {
		return NewDefaultSelectExecutionPlan(lp, reg)
	} else if CanBuildGroupbyExecutionPlan(lp, reg) {
		return NewGroupbyExecutionPlan(lp, reg)
	}
	return nil, fmt.Errorf("no plan can deal with such a statement")
}

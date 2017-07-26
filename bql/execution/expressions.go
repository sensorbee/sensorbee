package execution

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
)

// aliasedExpression represents an expression in a SELECT clause
// that has an output name associated and can be either an ordinary
// ("flat") expression or one that involves an aggregate computation.
type aliasedExpression struct {
	alias      string
	expr       FlatExpression
	aggrInputs map[string]FlatExpression
}

// Explanation of the Aggregation Workflow
// ---------------------------------------
// For a SELECT or CREATE STREAM FROM SELECT statement, we deal mostly
// with the "projections", i.e., the terms after the SELECT keyword.
// For a non-aggregate statement, we receive an instance of parser.Expression
// and transform it to an instance of execution.FlatExpression (which is
// guaranteed not to have any aggregate functions in it). This FlatExpression
// is then converted to an Evaluator that can be evaluated with a row
// as an input.
//
// Example 1: The BQL expression "x:a" will be converted to a
// `parser.RowValue{"x", "a"}` and then to an `execution.rowValue{"x", "a"}` and
// then to an `execution.PathAccess{"x.a"}`. If this `PathAccess{"x.a"}` is evaluated
// on a row like `data.Map{"x": data.Map{"a": data.Int(2), "b": data.Int(6)}}`,
// it will return `data.Int(2)`.
//
// Example 2: The BQL expression "f(x:a)" will be converted to a
// `parser.FuncAppAST{"f", parser.RowValue{"x", "a"}}` and then
// (if f is not an aggregate function) to an `execution.funcAppAST{"f",
// execution.rowValue{"x", "a"}}` and then to a `execution.funcApp{"f", ...,
// []Evaluator{execution.PathAccess{"x.a"}}, ...}`. If this `funcApp` is evaluated
// on a row like `data.Map{"x": data.Map{"a": data.Int(2),  "b": data.Int(6)}}`,
// the `PathAccess` evaluator will first return extract `data.Int(2)` and then
// the `funcApp` evaluator will compute the result of f using `data.Int(2)`
// as an input.
//
// For an expression that involves an aggregate function call, the flow
// is necessarily different:
// - The projections are not evaluated using the original rows as an input,
//   but using groups of rows that have the same values in the GROUP BY columns.
// - That means that in a pre-processing step, the input rows need to be traversed
//   and the groups described above must be created.
// - Each of those groups consists of
//   - the values that are common to all rows in a group
//   - a list of values that are to be aggregated.
// - After the groups have been created, each list of values for aggregation
//   is transformed to an array and added to the set of "common values" of
//   that group.
// - The aggregate function is executed like a normal function, but using the
//   array mentioned before as an input.
//
// Example: The BQL expression "avg(x:a)" will be converted to a
// `parser.FuncAppAST{"avg", parser.RowValue{"x", "a"}}`and then
// to an `execution.funcAppAST{"avg", execution.aggInputRef{"randstr"}}`
// together with a mapping `{"randstr": execution.rowValue{"x", "a"}}`.
// Those FlatExpressions are then converted to Evaluators so that we obtain
// `execution.funcApp{"avg", ..., []Evaluator{execution.PathAccess{"randstr"}}, ...}`
// and the map `{"randstr": execution.PathAccess{"x.a"}}`.
// For each row, the correct group will be computed and the data computed
// by the PathAccess{"x.a"} evaluator appended to a list, so that at the
// end of this process there is a group list such as:
//  group values | common data     | aggregate input
//  [1]          | {"x": {"b": 1}} | {"randstr": [3.5, 1.7, 0.9, 4.5, ...]}
//  [2]          | {"x": {"b": 2}} | {"randstr": [1.2, 0.8, 2.3]}
// This data is then merged so that we obtain a list such as:
//  common data
//  {"x": {"b": 1}, "randstr": [3.5, 1.7, 0.9, 4.5, ...]}
//  {"x": {"b": 2}, "randstr": [1.2, 0.8, 2.3]}
// This list looks exactly like a usual set of input rows so that the
// `funcApp{"avg", ...}` evaluator can be used normally.

// ParserExprToFlatExpr converts an expression obtained by the BQL parser
// to a FlatExpression, i.e., there are only expressions contained that
// can be evaluated on one single row and return an (unnamed) value.
// In particular, this fails for Expressions containing aggregate functions.
func ParserExprToFlatExpr(e parser.Expression, reg udf.FunctionRegistry) (FlatExpression, error) {
	switch obj := e.(type) {
	case parser.RowMeta:
		return rowMeta{obj.Relation, obj.MetaType}, nil
	case parser.RowValue:
		return rowValue{obj.Relation, obj.Column}, nil
	case parser.AliasAST:
		return ParserExprToFlatExpr(obj.Expr, reg)
	case parser.NullLiteral:
		return nullLiteral{}, nil
	case parser.NumericLiteral:
		return numericLiteral{obj.Value}, nil
	case parser.FloatLiteral:
		return floatLiteral{obj.Value}, nil
	case parser.BoolLiteral:
		return boolLiteral{obj.Value}, nil
	case parser.StringLiteral:
		return stringLiteral{obj.Value}, nil
	case parser.BinaryOpAST:
		// recurse left
		left, err := ParserExprToFlatExpr(obj.Left, reg)
		if err != nil {
			return nil, err
		}
		// In the case of IS [NOT] MISSING, we cannot use the standard
		// approach for "evaluate children and do something with the
		// result" since the evaluation will result in an error if the
		// specified path is actually missing. There are various methods
		// to work around this; we have chosen to use a special FlatExpr
		// for the IS [NOT] MISSING expression:
		if _, ok := obj.Right.(parser.Missing); ok {
			if rv, ok := left.(rowValue); ok {
				if obj.Op == parser.Is {
					return missing{rv, false}, nil
				} else if obj.Op == parser.IsNot {
					return missing{rv, true}, nil
				}
				// actually the parser should not allow
				// any operators except IS and IS NOT
				return nil, fmt.Errorf("MISSING requires IS [NOT], not: %s", obj.Op)
			}
			return nil, fmt.Errorf("IS [NOT] MISSING does not work with complex expression: %s", obj.Left.String())
		}
		// recurse right
		right, err := ParserExprToFlatExpr(obj.Right, reg)
		if err != nil {
			return nil, err
		}
		return binaryOpAST{obj.Op, left, right}, nil
	case parser.UnaryOpAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return unaryOpAST{obj.Op, expr}, nil
	case parser.TypeCastAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return typeCastAST{expr, obj.Target}, nil
	case parser.FuncAppSelectorAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.FuncAppAST, reg)
		if err != nil {
			return nil, err
		}
		// TODO: Relation is always empty
		return funcAppSelectorAST{
			Expr:     expr,
			Selector: fmt.Sprintf("%s", obj.Selector.Column),
		}, nil
	case parser.FuncAppAST:
		// exception for now()
		if string(obj.Function) == "now" && len(obj.Expressions) == 0 && len(obj.Ordering) == 0 {
			return stmtMeta{parser.NowMeta}, nil
		}
		// look up the function
		function, err := reg.Lookup(string(obj.Function), len(obj.Expressions))
		if err != nil {
			return nil, err
		}
		// fail if this is an aggregate function
		if isAggregateFunc(function, len(obj.Expressions)) {
			err := fmt.Errorf("you cannot use aggregate function '%s' "+
				"in a flat expression", obj.Function)
			return nil, err
		} else if len(obj.Ordering) > 0 {
			err := fmt.Errorf("you cannot use ORDER BY in non-aggregate "+
				"function '%s'", obj.Function)
			return nil, err
		}
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			expr, err := ParserExprToFlatExpr(ast, reg)
			if err != nil {
				return nil, err
			}
			exprs[i] = expr
		}
		return funcAppAST{obj.Function, exprs}, nil
	case parser.ArrayAST:
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		for i, ast := range obj.Expressions {
			expr, err := ParserExprToFlatExpr(ast, reg)
			if err != nil {
				return nil, err
			}
			exprs[i] = expr
		}
		return arrayAST{exprs}, nil
	case parser.MapAST:
		// compute child expressions
		pairs := make([]keyValuePair, len(obj.Entries))
		for i, pair := range obj.Entries {
			expr, err := ParserExprToFlatExpr(pair.Value, reg)
			if err != nil {
				return nil, err
			}
			pairs[i] = keyValuePair{pair.Key, expr}
		}
		return mapAST{pairs}, nil
	case parser.ConditionCaseAST:
		// compute child expressions
		pairs := make([]whenThenPair, len(obj.Checks))
		for i, pair := range obj.Checks {
			expr1, err := ParserExprToFlatExpr(pair.When, reg)
			if err != nil {
				return nil, err
			}
			expr2, err := ParserExprToFlatExpr(pair.Then, reg)
			if err != nil {
				return nil, err
			}
			pairs[i] = whenThenPair{expr1, expr2}
		}
		var defaultExpr FlatExpression = nullLiteral{}
		if obj.Else != nil {
			var err error
			defaultExpr, err = ParserExprToFlatExpr(obj.Else, reg)
			if err != nil {
				return nil, err
			}
		}
		return caseAST{boolLiteral{true}, pairs, defaultExpr}, nil
	case parser.ExpressionCaseAST:
		// compute child expressions
		_c, err := ParserExprToFlatExpr(obj.ConditionCaseAST, reg)
		if err != nil {
			return nil, err
		}
		c := _c.(caseAST)
		ref, err := ParserExprToFlatExpr(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		// return a new object
		return caseAST{ref, c.Checks, c.Default}, nil
	case parser.Wildcard:
		return wildcardAST{obj.Relation}, nil
	}
	err := fmt.Errorf("don't know how to convert type %#v", e)
	return nil, err
}

// ParserExprToMaybeAggregate converts an expression obtained by the BQL
// parser into a data structure where the aggregate and the non-aggregate
// parts are separated.
func ParserExprToMaybeAggregate(e parser.Expression, aggIdx int, reg udf.FunctionRegistry) (FlatExpression, map[string]FlatExpression, error) {
	switch obj := e.(type) {
	default:
		// elementary types
		expr, err := ParserExprToFlatExpr(e, reg)
		return expr, nil, err
	case parser.AliasAST:
		return ParserExprToMaybeAggregate(obj.Expr, aggIdx, reg)
	case parser.BinaryOpAST:
		// recurse left
		left, leftAgg, err := ParserExprToMaybeAggregate(obj.Left, aggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		// In the case of IS [NOT] MISSING, we cannot use the standard
		// approach for "evaluate children and do something with the
		// result" since the evaluation will result in an error if the
		// specified path is actually missing. There are various methods
		// to work around this; we have chosen to use a special FlatExpr
		// for the IS [NOT] MISSING expression:
		if _, ok := obj.Right.(parser.Missing); ok {
			if rv, ok := left.(rowValue); ok {
				if obj.Op == parser.Is {
					return missing{rv, false}, leftAgg, nil
				} else if obj.Op == parser.IsNot {
					return missing{rv, true}, leftAgg, nil
				}
				// actually the parser should not allow
				// any operators except IS and IS NOT
				return nil, nil, fmt.Errorf("MISSING requires IS [NOT], not: %s", obj.Op)
			}
			return nil, nil, fmt.Errorf("IS [NOT] MISSING does not work with complex expression: %s", obj.Left.String())
		}
		// return right
		right, rightAgg, err := ParserExprToMaybeAggregate(obj.Right, aggIdx+len(leftAgg), reg)
		if err != nil {
			return nil, nil, err
		}
		var returnAgg map[string]FlatExpression
		if leftAgg != nil {
			returnAgg = leftAgg
			for key, val := range rightAgg {
				returnAgg[key] = val
			}
		} else if leftAgg == nil {
			returnAgg = rightAgg
		}
		return binaryOpAST{obj.Op, left, right}, returnAgg, nil
	case parser.UnaryOpAST:
		// recurse
		expr, agg, err := ParserExprToMaybeAggregate(obj.Expr, aggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		return unaryOpAST{obj.Op, expr}, agg, nil
	case parser.TypeCastAST:
		// recurse
		expr, agg, err := ParserExprToMaybeAggregate(obj.Expr, aggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		return typeCastAST{expr, obj.Target}, agg, nil
	case parser.FuncAppSelectorAST:
		// recurse
		expr, agg, err := ParserExprToMaybeAggregate(obj.FuncAppAST, aggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		return funcAppSelectorAST{
			Expr:     expr,
			Selector: fmt.Sprintf("%s", obj.Selector.Column),
		}, agg, nil
	case parser.FuncAppAST:
		// exception for now()
		if string(obj.Function) == "now" && len(obj.Expressions) == 0 {
			return stmtMeta{parser.NowMeta}, nil, nil
		}
		// look up the function
		function, err := reg.Lookup(string(obj.Function), len(obj.Expressions))
		if err != nil {
			return nil, nil, err
		}
		// replace the "*" by 1 for the count function
		for i, ast := range obj.Expressions {
			if _, ok := ast.(parser.Wildcard); ok {
				if string(obj.Function) == "count" {
					// replace the wildcard by an always non-null expression
					obj.Expressions[i] = parser.NumericLiteral{1}
				}
			}
		}
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		returnAgg := map[string]FlatExpression{}
		if isAggregateFunc(function, len(obj.Expressions)) {
			// we have a setting like
			//  SELECT udaf(x+1, "state", c ORDER BY d + e, f DESC) ... GROUP BY c
			// where some parameters are aggregates, others aren't.
			for i, ast := range obj.Expressions {
				// this expression must be flat, there must not be other aggregates
				expr, err := ParserExprToFlatExpr(ast, reg)
				if err != nil {
					// return a prettier error message
					if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
						err = fmt.Errorf("aggregate functions cannot be nested")
					}
					return nil, nil, err
				}
				isAggr := function.IsAggregationParameter(i)
				if isAggr {
					// this is an aggregation parameter, we will replace
					// it by a reference to the aggregated list of values
					h := sha1.New()
					h.Write([]byte(fmt.Sprintf("%s", expr.Repr())))
					exprID := "g_" + hex.EncodeToString(h.Sum(nil))[:8]
					// For stable or immutable expressions, we compute a reference
					// string that depends only on the expression's string
					// representation so that we have to compute and store values
					// only once per row (e.g., for `sum(a)/count(a)`).
					// For volatile expressions (e.g., `sum(random())/avg(random())`)
					// we add a numeric suffix that represents the index
					// of this aggregate in the whole projection list.
					if expr.Volatility() == Volatile {
						exprID += fmt.Sprintf("_%d", aggIdx+len(returnAgg))
					}
					exprs[i] = aggInputRef{exprID}
					returnAgg[exprID] = expr
				} else {
					// this is a non-aggregate parameter, use as is
					exprs[i] = expr
				}
			}

			// deal with ORDER BY specifications
			if len(obj.Ordering) > 0 {
				ordering := make([]sortExpression, len(obj.Ordering))
				// we need a string that uniquely identifies this ordering in order
				// to allow `SELECT f(a ORDER BY b), f(a ORDER BY c)`
				orderHash := sha1.New()
				// also add all the expressions in obj.Ordering to the list
				// of aggregate values to be computed
				for i, sortExpr := range obj.Ordering {
					// this expression must be flat, there must not be other aggregates
					expr, err := ParserExprToFlatExpr(sortExpr.Expr, reg)
					if err != nil {
						// return a prettier error message
						if strings.HasPrefix(err.Error(), "you cannot use aggregate") {
							err = fmt.Errorf("aggregate functions cannot be used in ORDER BY")
						}
						return nil, nil, err
					}
					// we will replace this value by a reference to the
					// aggregated list of values
					h := sha1.New()
					h.Write([]byte(fmt.Sprintf("%s", expr.Repr())))
					orderHash.Write([]byte(fmt.Sprintf("%s", expr.Repr())))
					exprID := "g_" + hex.EncodeToString(h.Sum(nil))[:8]
					// For stable or immutable expressions, we compute a reference
					// string that depends only on the expression's string
					// representation so that we have to compute and store values
					// only once per row (e.g., for `sum(a)/count(a)`).
					// For volatile expressions (e.g., `sum(random())/avg(random())`)
					// we add a numeric suffix that represents the index
					// of this aggregate in the whole projection list.
					if expr.Volatility() == Volatile {
						exprID += fmt.Sprintf("_%d", aggIdx+len(returnAgg))
					}
					// remember the information about exprID and the order direction
					ascending := true
					if sortExpr.Ascending == parser.No {
						ascending = false
						orderHash.Write([]byte(" DESC"))
					}
					orderHash.Write([]byte(","))
					ordering[i] = sortExpression{
						aggInputRef{exprID}, ascending,
					}
					returnAgg[exprID] = expr
				}
				return aggregateInputSorter{
					funcAppAST{obj.Function, exprs},
					ordering,
					hex.EncodeToString(orderHash.Sum(nil))[:8],
				}, returnAgg, nil
			}

		} else {
			for i, ast := range obj.Expressions {
				expr, agg, err := ParserExprToMaybeAggregate(ast, aggIdx, reg)
				if err != nil {
					return nil, nil, err
				}
				for key, val := range agg {
					returnAgg[key] = val
				}
				exprs[i] = expr
			}
			if len(returnAgg) == 0 {
				returnAgg = nil
			}
		}
		return funcAppAST{obj.Function, exprs}, returnAgg, nil
	case parser.ArrayAST:
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		returnAgg := map[string]FlatExpression{}
		for i, ast := range obj.Expressions {
			// compute the correct aggIdx
			newAggIdx := aggIdx + len(returnAgg)
			expr, agg, err := ParserExprToMaybeAggregate(ast, newAggIdx, reg)
			if err != nil {
				return nil, nil, err
			}
			for key, val := range agg {
				returnAgg[key] = val
			}
			exprs[i] = expr
		}
		if len(returnAgg) == 0 {
			returnAgg = nil
		}
		return arrayAST{exprs}, returnAgg, nil
	case parser.MapAST:
		// compute child expressions
		pairs := make([]keyValuePair, len(obj.Entries))
		returnAgg := map[string]FlatExpression{}
		for i, pair := range obj.Entries {
			// compute the correct aggIdx
			newAggIdx := aggIdx + len(returnAgg)
			expr, agg, err := ParserExprToMaybeAggregate(pair.Value, newAggIdx, reg)
			if err != nil {
				return nil, nil, err
			}
			for key, val := range agg {
				returnAgg[key] = val
			}
			pairs[i] = keyValuePair{pair.Key, expr}
		}
		if len(returnAgg) == 0 {
			returnAgg = nil
		}
		return mapAST{pairs}, returnAgg, nil
	case parser.ConditionCaseAST:
		// compute child expressions
		pairs := make([]whenThenPair, len(obj.Checks))
		returnAgg := map[string]FlatExpression{}
		for i, pair := range obj.Checks {
			// compute the correct aggIdx
			newAggIdx := aggIdx + len(returnAgg)
			expr1, agg, err := ParserExprToMaybeAggregate(pair.When, newAggIdx, reg)
			if err != nil {
				return nil, nil, err
			}
			for key, val := range agg {
				returnAgg[key] = val
			}

			newAggIdx = aggIdx + len(returnAgg)
			expr2, agg, err := ParserExprToMaybeAggregate(pair.Then, newAggIdx, reg)
			if err != nil {
				return nil, nil, err
			}
			for key, val := range agg {
				returnAgg[key] = val
			}
			pairs[i] = whenThenPair{expr1, expr2}
		}

		newAggIdx := aggIdx + len(returnAgg)
		var defaultExpr FlatExpression = nullLiteral{}
		if obj.Else != nil {
			var err error
			var agg map[string]FlatExpression
			defaultExpr, agg, err = ParserExprToMaybeAggregate(obj.Else, newAggIdx, reg)
			if err != nil {
				return nil, nil, err
			}
			for key, val := range agg {
				returnAgg[key] = val
			}
		}

		if len(returnAgg) == 0 {
			returnAgg = nil
		}
		return caseAST{boolLiteral{true}, pairs, defaultExpr}, returnAgg, nil

	case parser.ExpressionCaseAST:
		// compute child expressions
		returnAgg := map[string]FlatExpression{}
		newAggIdx := aggIdx + len(returnAgg)
		_c, agg, err := ParserExprToMaybeAggregate(obj.ConditionCaseAST, newAggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		for key, val := range agg {
			returnAgg[key] = val
		}
		c := _c.(caseAST)

		newAggIdx = aggIdx + len(returnAgg)
		ref, agg, err := ParserExprToMaybeAggregate(obj.Expr, newAggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		for key, val := range agg {
			returnAgg[key] = val
		}

		if len(returnAgg) == 0 {
			returnAgg = nil
		}
		return caseAST{ref, c.Checks, c.Default}, returnAgg, nil
	}
	err := fmt.Errorf("don't know how to convert type %#v", e)
	return nil, nil, err
}

// FlatExpression represents an expression that can be completely
// evaluated on a single row and results in an unnamed value. In
// particular, it cannot contain/represent a call to an aggregate
// function.
type FlatExpression interface {
	// Repr returns a string representation that can be used to
	// identify this expression (e.g., "stream:col+3") and used as
	// a dictionary key for finding duplicate expressions.
	Repr() string

	// Columns returns a list of rowValues used in this expression.
	Columns() []rowValue

	// Volatility returns the volatility of an expression.
	Volatility() VolatilityType

	// ContainsWildcard returns whether this expression contains
	// a wildcard symbol.
	ContainsWildcard() bool
}

// VolatilityType describes the volatility of an expression as per
// the PostgreSQL classification.
type VolatilityType int

const (
	// UnknownVolatility describes an unset value. This is not
	// a valid return value for FlatExpression.Volatility().
	UnknownVolatility = iota
	// Volatile expressions can do anything, in particular return a
	// different result on every call
	Volatile
	// Stable expressions return the same result for the same input
	// values within a single statement execution
	Stable
	// Immutable expressions return the same result for the same input
	// values forever. One good hint to distinguish between Stable
	// and Immutable is that (in PostgreSQL) Immutable functions can
	// be used in functional indexes, while Stable functions can't.
	Immutable
)

func (v VolatilityType) String() string {
	s := "UNKNOWN"
	switch v {
	case Volatile:
		s = "VOLATILE"
	case Stable:
		s = "STABLE"
	case Immutable:
		s = "IMMUTABLE"
	}
	return s
}

type binaryOpAST struct {
	Op    parser.Operator
	Left  FlatExpression
	Right FlatExpression
}

func (b binaryOpAST) Repr() string {
	return fmt.Sprintf("(%s)%s(%s)", b.Left.Repr(), b.Op, b.Right.Repr())
}

func (b binaryOpAST) Columns() []rowValue {
	return append(b.Left.Columns(), b.Right.Columns()...)
}

func (b binaryOpAST) Volatility() VolatilityType {
	// take the lower level of both sub-expressions
	// (all of the operators have IMMUTABLE behavior)
	l := b.Left.Volatility()
	r := b.Right.Volatility()
	if l < r {
		return l
	}
	return r
}

func (b binaryOpAST) ContainsWildcard() bool {
	return b.Left.ContainsWildcard() || b.Right.ContainsWildcard()
}

type unaryOpAST struct {
	Op   parser.Operator
	Expr FlatExpression
}

func (u unaryOpAST) Repr() string {
	return fmt.Sprintf("%s(%s)", u.Op, u.Expr.Repr())
}

func (u unaryOpAST) Columns() []rowValue {
	return u.Expr.Columns()
}

func (u unaryOpAST) Volatility() VolatilityType {
	return u.Expr.Volatility()
}

func (u unaryOpAST) ContainsWildcard() bool {
	return u.Expr.ContainsWildcard()
}

type typeCastAST struct {
	Expr   FlatExpression
	Target parser.Type
}

func (t typeCastAST) Repr() string {
	return fmt.Sprintf("CAST(%s AS %s)", t.Expr.Repr(), t.Target)
}

func (t typeCastAST) Columns() []rowValue {
	return t.Expr.Columns()
}

func (t typeCastAST) Volatility() VolatilityType {
	return t.Expr.Volatility()
}

func (t typeCastAST) ContainsWildcard() bool {
	return t.Expr.ContainsWildcard()
}

type funcAppAST struct {
	Function    parser.FuncName
	Expressions []FlatExpression
}

func (f funcAppAST) Repr() string {
	reprs := make([]string, len(f.Expressions))
	for i, e := range f.Expressions {
		reprs[i] = e.Repr()
	}
	return fmt.Sprintf("%s(%s)", f.Function, strings.Join(reprs, ","))
}

func (f funcAppAST) Columns() []rowValue {
	var allColumns []rowValue
	for _, e := range f.Expressions {
		allColumns = append(allColumns, e.Columns()...)
	}
	return allColumns
}

func (f funcAppAST) Volatility() VolatilityType {
	// we can change this later, but for the moment we
	// cannot assume that UDFs are stable or immutable
	// in general
	return Volatile
}

func (f funcAppAST) ContainsWildcard() bool {
	for _, e := range f.Expressions {
		if e.ContainsWildcard() {
			return true
		}
	}
	return false
}

type funcAppSelectorAST struct {
	Expr     FlatExpression
	Selector string
}

func (f funcAppSelectorAST) Repr() string {
	return fmt.Sprintf("%s%s", f.Expr.Repr(), f.Selector)
}

func (f funcAppSelectorAST) Columns() []rowValue {
	return f.Expr.Columns()
}

func (f funcAppSelectorAST) Volatility() VolatilityType {
	return f.Expr.Volatility()
}

func (f funcAppSelectorAST) ContainsWildcard() bool {
	return f.Expr.ContainsWildcard()
}

type sortExpression struct {
	Value     aggInputRef
	Ascending bool
}

type aggregateInputSorter struct {
	funcAppAST
	Ordering []sortExpression
	ID       string
}

func (a aggregateInputSorter) Repr() string {
	reprs := make([]string, len(a.Expressions))
	for i, e := range a.Expressions {
		reprs[i] = e.Repr()
	}
	ordering := make([]string, len(a.Ordering))
	for i, e := range a.Ordering {
		ordering[i] = e.Value.Repr()
		if e.Ascending {
			ordering[i] += " ASC"
		} else {
			ordering[i] += " DESC"
		}
	}
	return fmt.Sprintf("%s(%s ORDER BY %s)", a.Function,
		strings.Join(reprs, ","), strings.Join(ordering, ","))
}

type arrayAST struct {
	Expressions []FlatExpression
}

func (a arrayAST) Repr() string {
	reprs := make([]string, len(a.Expressions))
	for i, e := range a.Expressions {
		reprs[i] = e.Repr()
	}
	return fmt.Sprintf("[%s]", strings.Join(reprs, ", "))
}

func (a arrayAST) Columns() []rowValue {
	var allColumns []rowValue
	for _, e := range a.Expressions {
		allColumns = append(allColumns, e.Columns()...)
	}
	return allColumns
}

func (a arrayAST) Volatility() VolatilityType {
	lv := VolatilityType(Immutable)
	for _, e := range a.Expressions {
		v := e.Volatility()
		if v < lv {
			lv = v
		}
	}
	return lv
}

func (a arrayAST) ContainsWildcard() bool {
	for _, e := range a.Expressions {
		if e.ContainsWildcard() {
			return true
		}
	}
	return false
}

type mapAST struct {
	Entries []keyValuePair
}

func (m mapAST) Repr() string {
	reprs := make([]string, len(m.Entries))
	for i, p := range m.Entries {
		reprs[i] = fmt.Sprintf("\"%s\":%s", p.Key, p.Value.Repr())
	}
	return fmt.Sprintf("{%s}", strings.Join(reprs, ", "))
}

func (m mapAST) Columns() []rowValue {
	var allColumns []rowValue
	for _, p := range m.Entries {
		allColumns = append(allColumns, p.Value.Columns()...)
	}
	return allColumns
}

func (m mapAST) Volatility() VolatilityType {
	lv := VolatilityType(Immutable)
	for _, p := range m.Entries {
		v := p.Value.Volatility()
		if v < lv {
			lv = v
		}
	}
	return lv
}

func (m mapAST) ContainsWildcard() bool {
	for _, p := range m.Entries {
		if p.Value.ContainsWildcard() {
			return true
		}
	}
	return false
}

type keyValuePair struct {
	Key   string
	Value FlatExpression
}

type caseAST struct {
	Reference FlatExpression
	Checks    []whenThenPair
	Default   FlatExpression
}

func (c caseAST) Repr() string {
	reprs := make([]string, len(c.Checks))
	for i, p := range c.Checks {
		reprs[i] = fmt.Sprintf("%s->%s", p.When.Repr(), p.Then.Repr())
	}
	return fmt.Sprintf("case(%s:%s,%s)",
		c.Reference.Repr(), strings.Join(reprs, ","), c.Default.Repr())
}

func (c caseAST) Columns() []rowValue {
	allColumns := c.Default.Columns()
	for _, p := range c.Checks {
		allColumns = append(allColumns, p.When.Columns()...)
		allColumns = append(allColumns, p.Then.Columns()...)
	}
	allColumns = append(allColumns, c.Default.Columns()...)
	return allColumns
}

func (c caseAST) Volatility() VolatilityType {
	lv := c.Reference.Volatility()
	for _, p := range c.Checks {
		v := p.When.Volatility()
		if v < lv {
			lv = v
		}
		v = p.Then.Volatility()
		if v < lv {
			lv = v
		}
	}
	v := c.Default.Volatility()
	if v < lv {
		lv = v
	}
	return lv
}

func (c caseAST) ContainsWildcard() bool {
	if c.Reference.ContainsWildcard() {
		return true
	}
	for _, p := range c.Checks {
		if p.When.ContainsWildcard() {
			return true
		}
		if p.Then.ContainsWildcard() {
			return true
		}
	}
	if c.Default.ContainsWildcard() {
		return true
	}
	return false
}

type whenThenPair struct {
	When FlatExpression
	Then FlatExpression
}

type wildcardAST struct {
	Relation string
}

func (w wildcardAST) Repr() string {
	if w.Relation == "" {
		return "*"
	}
	return fmt.Sprintf("%s:*", w.Relation)
}

func (w wildcardAST) Columns() []rowValue {
	return nil
}

func (w wildcardAST) Volatility() VolatilityType {
	// this selects all elements of a tuple, which is
	// a stable operation (maybe it is even immutable,
	// but probably we will never evaluate this)
	return Stable
}

func (w wildcardAST) ContainsWildcard() bool {
	return true
}

type aggInputRef struct {
	Ref string
}

func (a aggInputRef) Repr() string {
	return a.Ref
}

func (a aggInputRef) Columns() []rowValue {
	return nil
}

func (a aggInputRef) Volatility() VolatilityType {
	// TODO make this the same volatility level as the
	//      aggregate function used
	return Volatile
}

func (a aggInputRef) ContainsWildcard() bool {
	// this references the result of an aggregate function,
	// and whether there is a wildcard used in that aggregate
	// function is irrelevant for the reference
	return false
}

type rowValue struct {
	Relation string
	Column   string
}

func (rv rowValue) Repr() string {
	return fmt.Sprintf("%s:%s", rv.Relation, rv.Column)
}

func (rv rowValue) Columns() []rowValue {
	return []rowValue{rv}
}

func (rv rowValue) Volatility() VolatilityType {
	return Immutable
}

func (rv rowValue) ContainsWildcard() bool {
	return false
}

type stmtMeta struct {
	MetaType parser.MetaInformation
}

func (sm stmtMeta) Repr() string {
	return fmt.Sprintf("%#v", sm)
}

func (sm stmtMeta) Columns() []rowValue {
	return nil
}

func (sm stmtMeta) Volatility() VolatilityType {
	return Stable
}

func (sm stmtMeta) ContainsWildcard() bool {
	return false
}

type rowMeta struct {
	Relation string
	MetaType parser.MetaInformation
}

func (rm rowMeta) Repr() string {
	return fmt.Sprintf("%#v", rm)
}

func (rm rowMeta) Columns() []rowValue {
	return nil
}

func (rm rowMeta) Volatility() VolatilityType {
	return Immutable
}

func (rm rowMeta) ContainsWildcard() bool {
	return false
}

type numericLiteral struct {
	Value int64
}

func (l numericLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l numericLiteral) Columns() []rowValue {
	return nil
}

func (l numericLiteral) Volatility() VolatilityType {
	return Immutable
}

func (l numericLiteral) ContainsWildcard() bool {
	return false
}

type floatLiteral struct {
	Value float64
}

func (l floatLiteral) Repr() string {
	return fmt.Sprintf("%vf", l.Value)
}

func (l floatLiteral) Columns() []rowValue {
	return nil
}

func (l floatLiteral) Volatility() VolatilityType {
	return Immutable
}

func (l floatLiteral) ContainsWildcard() bool {
	return false
}

type nullLiteral struct {
}

func (l nullLiteral) Repr() string {
	return "NULL"
}

func (l nullLiteral) Columns() []rowValue {
	return nil
}

func (l nullLiteral) Volatility() VolatilityType {
	return Immutable
}

func (l nullLiteral) ContainsWildcard() bool {
	return false
}

type missing struct {
	Expr rowValue
	Not  bool
}

func (m missing) Repr() string {
	if m.Not {
		return fmt.Sprintf("(%s)-NOT-MISSING", m.Expr.Repr())
	}
	return fmt.Sprintf("(%s)-MISSING", m.Expr.Repr())
}

func (m missing) Columns() []rowValue {
	return m.Expr.Columns()
}

func (m missing) Volatility() VolatilityType {
	return m.Expr.Volatility()
}

func (m missing) ContainsWildcard() bool {
	return m.Expr.ContainsWildcard()
}

type boolLiteral struct {
	Value bool
}

func (l boolLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l boolLiteral) Columns() []rowValue {
	return nil
}

func (l boolLiteral) Volatility() VolatilityType {
	return Immutable
}

func (l boolLiteral) ContainsWildcard() bool {
	return false
}

type stringLiteral struct {
	Value string
}

func (l stringLiteral) Repr() string {
	return fmt.Sprintf("%s", l.Value)
}

func (l stringLiteral) Columns() []rowValue {
	return nil
}

func (l stringLiteral) Volatility() VolatilityType {
	return Immutable
}

func (l stringLiteral) ContainsWildcard() bool {
	return false
}

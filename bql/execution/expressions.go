package execution

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"strings"
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
// `parser.RowValue{"x", "a"}` and then to an `execution.RowValue{"x", "a"}` and
// then to an `execution.PathAccess{"x.a"}`. If this `PathAccess{"x.a"}` is evaluated
// on a row like `data.Map{"x": data.Map{"a": data.Int(2), "b": data.Int(6)}}`,
// it will return `data.Int(2)`.
//
// Example 2: The BQL expression "f(x:a)" will be converted to a
// `parser.FuncAppAST{"f", parser.RowValue{"x", "a"}}` and then
// (if f is not an aggregate function) to an `execution.FuncAppAST{"f",
// execution.RowValue{"x", "a"}}` and then to a `execution.funcApp{"f", ...,
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
// to an `execution.FuncAppAST{"avg", execution.AggInputRef{"randstr"}}`
// together with a mapping `{"randstr": execution.RowValue{"x", "a"}}`.
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
		return RowMeta{obj.Relation, obj.MetaType}, nil
	case parser.RowValue:
		return RowValue{obj.Relation, obj.Column}, nil
	case parser.AliasAST:
		return ParserExprToFlatExpr(obj.Expr, reg)
	case parser.NullLiteral:
		return NullLiteral{}, nil
	case parser.NumericLiteral:
		return NumericLiteral{obj.Value}, nil
	case parser.FloatLiteral:
		return FloatLiteral{obj.Value}, nil
	case parser.BoolLiteral:
		return BoolLiteral{obj.Value}, nil
	case parser.StringLiteral:
		return StringLiteral{obj.Value}, nil
	case parser.BinaryOpAST:
		// recurse
		left, err := ParserExprToFlatExpr(obj.Left, reg)
		if err != nil {
			return nil, err
		}
		right, err := ParserExprToFlatExpr(obj.Right, reg)
		if err != nil {
			return nil, err
		}
		return BinaryOpAST{obj.Op, left, right}, nil
	case parser.UnaryOpAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return UnaryOpAST{obj.Op, expr}, nil
	case parser.TypeCastAST:
		// recurse
		expr, err := ParserExprToFlatExpr(obj.Expr, reg)
		if err != nil {
			return nil, err
		}
		return TypeCastAST{expr, obj.Target}, nil
	case parser.FuncAppAST:
		// exception for now()
		if string(obj.Function) == "now" && len(obj.Expressions) == 0 {
			return StmtMeta{parser.NowMeta}, nil
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
		}
		// fail if this uses the "*" parameter
		for _, ast := range obj.Expressions {
			if _, ok := ast.(parser.Wildcard); ok {
				return nil, fmt.Errorf("* can only be used " +
					"as a parameter in count()")
			}
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
		return FuncAppAST{obj.Function, exprs}, nil
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
		return ArrayAST{exprs}, nil
	case parser.MapAST:
		// compute child expressions
		pairs := make([]KeyValuePair, len(obj.Entries))
		for i, pair := range obj.Entries {
			expr, err := ParserExprToFlatExpr(pair.Value, reg)
			if err != nil {
				return nil, err
			}
			pairs[i] = KeyValuePair{pair.Key, expr}
		}
		return MapAST{pairs}, nil
	case parser.Wildcard:
		return WildcardAST{obj.Relation}, nil
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
		// recurse
		left, leftAgg, err := ParserExprToMaybeAggregate(obj.Left, aggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
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
		return BinaryOpAST{obj.Op, left, right}, returnAgg, nil
	case parser.UnaryOpAST:
		// recurse
		expr, agg, err := ParserExprToMaybeAggregate(obj.Expr, aggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		return UnaryOpAST{obj.Op, expr}, agg, nil
	case parser.TypeCastAST:
		// recurse
		expr, agg, err := ParserExprToMaybeAggregate(obj.Expr, aggIdx, reg)
		if err != nil {
			return nil, nil, err
		}
		return TypeCastAST{expr, obj.Target}, agg, nil
	case parser.FuncAppAST:
		// exception for now()
		if string(obj.Function) == "now" && len(obj.Expressions) == 0 {
			return StmtMeta{parser.NowMeta}, nil, nil
		}
		// look up the function
		function, err := reg.Lookup(string(obj.Function), len(obj.Expressions))
		if err != nil {
			return nil, nil, err
		}
		// fail if this uses the "*" parameter for anything else than count()
		for i, ast := range obj.Expressions {
			if _, ok := ast.(parser.Wildcard); ok {
				if string(obj.Function) == "count" {
					// replace the wildcard by an always non-null expression
					obj.Expressions[i] = parser.NumericLiteral{1}
				} else {
					return nil, nil, fmt.Errorf("* can only be used " +
						"as a parameter in count()")
				}
			}
		}
		// compute child expressions
		exprs := make([]FlatExpression, len(obj.Expressions))
		returnAgg := map[string]FlatExpression{}
		if isAggregateFunc(function, len(obj.Expressions)) {
			// we have a setting like
			//  SELECT udaf(x+1, 'state', c) ... GROUP BY c
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
					exprId := "_" + hex.EncodeToString(h.Sum(nil))[:8]
					// For stable or immutable expressions, we compute a reference
					// string that depends only on the expression's string
					// representation so that we have to compute and store values
					// only once per row (e.g., for `sum(a)/count(a)`).
					// For volatile expressions (e.g., `sum(random())/avg(random())`)
					// we add a numeric suffix that represents the index
					// of this aggregate in the whole projection list.
					if expr.Volatility() == Volatile {
						exprId += fmt.Sprintf("_%d", aggIdx+len(returnAgg))
					}
					exprs[i] = AggInputRef{exprId}
					returnAgg[exprId] = expr
				} else {
					// this is a non-aggregate parameter, use as is
					exprs[i] = expr
				}
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
		return FuncAppAST{obj.Function, exprs}, returnAgg, nil
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
		return ArrayAST{exprs}, returnAgg, nil
	case parser.MapAST:
		// compute child expressions
		pairs := make([]KeyValuePair, len(obj.Entries))
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
			pairs[i] = KeyValuePair{pair.Key, expr}
		}
		if len(returnAgg) == 0 {
			returnAgg = nil
		}
		return MapAST{pairs}, returnAgg, nil
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

	// Columns returns a list of RowValues used in this expression.
	Columns() []RowValue

	// Volatility returns the volatility type of an expression as
	// per the PostgreSQL classification:
	// - VOLATILE: can do anything, in particular return a different
	//             result on every call
	// - STABLE: returns the same result for the same input values
	//           within a single statement execution
	// - IMMUTABLE: returns the same result for the same input values
	//              forever
	// One good hint to distinguish between Stable and Immutable is
	// that (in PostgreSQL) Immutable functions can be used in
	// functional indexes, while Stable functions can't.
	Volatility() VolatilityType
}

type VolatilityType int

const (
	UnknownVolatility = iota
	Volatile
	Stable
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

type BinaryOpAST struct {
	Op    parser.Operator
	Left  FlatExpression
	Right FlatExpression
}

func (b BinaryOpAST) Repr() string {
	return fmt.Sprintf("(%s)%s(%s)", b.Left.Repr(), b.Op, b.Right.Repr())
}

func (b BinaryOpAST) Columns() []RowValue {
	return append(b.Left.Columns(), b.Right.Columns()...)
}

func (b BinaryOpAST) Volatility() VolatilityType {
	// take the lower level of both sub-expressions
	// (all of the operators have IMMUTABLE behavior)
	l := b.Left.Volatility()
	r := b.Right.Volatility()
	if l < r {
		return l
	}
	return r
}

type UnaryOpAST struct {
	Op   parser.Operator
	Expr FlatExpression
}

func (u UnaryOpAST) Repr() string {
	return fmt.Sprintf("%s(%s)", u.Op, u.Expr.Repr())
}

func (u UnaryOpAST) Columns() []RowValue {
	return u.Expr.Columns()
}

func (u UnaryOpAST) Volatility() VolatilityType {
	return u.Expr.Volatility()
}

type TypeCastAST struct {
	Expr   FlatExpression
	Target parser.Type
}

func (t TypeCastAST) Repr() string {
	return fmt.Sprintf("CAST(%s AS %s)", t.Expr.Repr(), t.Target)
}

func (t TypeCastAST) Columns() []RowValue {
	return t.Expr.Columns()
}

func (t TypeCastAST) Volatility() VolatilityType {
	return t.Expr.Volatility()
}

type FuncAppAST struct {
	Function    parser.FuncName
	Expressions []FlatExpression
}

func (f FuncAppAST) Repr() string {
	reprs := make([]string, len(f.Expressions))
	for i, e := range f.Expressions {
		reprs[i] = e.Repr()
	}
	return fmt.Sprintf("%s(%s)", f.Function, strings.Join(reprs, ","))
}

func (f FuncAppAST) Columns() []RowValue {
	var allColumns []RowValue
	for _, e := range f.Expressions {
		allColumns = append(allColumns, e.Columns()...)
	}
	return allColumns
}

func (f FuncAppAST) Volatility() VolatilityType {
	// we can change this later, but for the moment we
	// cannot assume that UDFs are stable or immutable
	// in general
	return Volatile
}

type ArrayAST struct {
	Expressions []FlatExpression
}

func (a ArrayAST) Repr() string {
	reprs := make([]string, len(a.Expressions))
	for i, e := range a.Expressions {
		reprs[i] = e.Repr()
	}
	return fmt.Sprintf("[%s]", strings.Join(reprs, ", "))
}

func (a ArrayAST) Columns() []RowValue {
	var allColumns []RowValue
	for _, e := range a.Expressions {
		allColumns = append(allColumns, e.Columns()...)
	}
	return allColumns
}

func (a ArrayAST) Volatility() VolatilityType {
	lv := VolatilityType(Immutable)
	for _, e := range a.Expressions {
		v := e.Volatility()
		if v < lv {
			lv = v
		}
	}
	return lv
}

type MapAST struct {
	Entries []KeyValuePair
}

func (m MapAST) Repr() string {
	reprs := make([]string, len(m.Entries))
	for i, p := range m.Entries {
		reprs[i] = fmt.Sprintf("\"%s\":%s", p.Key, p.Value.Repr())
	}
	return fmt.Sprintf("{%s}", strings.Join(reprs, ", "))
}

func (m MapAST) Columns() []RowValue {
	var allColumns []RowValue
	for _, p := range m.Entries {
		allColumns = append(allColumns, p.Value.Columns()...)
	}
	return allColumns
}

func (m MapAST) Volatility() VolatilityType {
	lv := VolatilityType(Immutable)
	for _, p := range m.Entries {
		v := p.Value.Volatility()
		if v < lv {
			lv = v
		}
	}
	return lv
}

type KeyValuePair struct {
	Key   string
	Value FlatExpression
}

type WildcardAST struct {
	Relation string
}

func (w WildcardAST) Repr() string {
	if w.Relation == "" {
		return "*"
	}
	return fmt.Sprintf("%s:*", w.Relation)
}

func (w WildcardAST) Columns() []RowValue {
	return nil
}

func (w WildcardAST) Volatility() VolatilityType {
	// this selects all elements of a tuple, which is
	// a stable operation (maybe it is even immutable,
	// but probably we will never evaluate this)
	return Stable
}

type AggInputRef struct {
	Ref string
}

func (a AggInputRef) Repr() string {
	return a.Ref
}

func (a AggInputRef) Columns() []RowValue {
	return nil
}

func (a AggInputRef) Volatility() VolatilityType {
	// TODO make this the same volatility level as the
	//      aggregate function used
	return Volatile
}

type RowValue struct {
	Relation string
	Column   string
}

func (rv RowValue) Repr() string {
	return fmt.Sprintf("%s:%s", rv.Relation, rv.Column)
}

func (rv RowValue) Columns() []RowValue {
	return []RowValue{rv}
}

func (rv RowValue) Volatility() VolatilityType {
	return Immutable
}

type StmtMeta struct {
	MetaType parser.MetaInformation
}

func (sm StmtMeta) Repr() string {
	return fmt.Sprintf("%#v", sm)
}

func (sm StmtMeta) Columns() []RowValue {
	return nil
}

func (sm StmtMeta) Volatility() VolatilityType {
	return Stable
}

type RowMeta struct {
	Relation string
	MetaType parser.MetaInformation
}

func (rm RowMeta) Repr() string {
	return fmt.Sprintf("%#v", rm)
}

func (rm RowMeta) Columns() []RowValue {
	return nil
}

func (rm RowMeta) Volatility() VolatilityType {
	return Immutable
}

type NumericLiteral struct {
	Value int64
}

func (l NumericLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l NumericLiteral) Columns() []RowValue {
	return nil
}

func (l NumericLiteral) Volatility() VolatilityType {
	return Immutable
}

type FloatLiteral struct {
	Value float64
}

func (l FloatLiteral) Repr() string {
	return fmt.Sprintf("%vf", l.Value)
}

func (l FloatLiteral) Columns() []RowValue {
	return nil
}

func (l FloatLiteral) Volatility() VolatilityType {
	return Immutable
}

type NullLiteral struct {
}

func (l NullLiteral) Repr() string {
	return "NULL"
}

func (l NullLiteral) Columns() []RowValue {
	return nil
}

func (l NullLiteral) Volatility() VolatilityType {
	return Immutable
}

type BoolLiteral struct {
	Value bool
}

func (l BoolLiteral) Repr() string {
	return fmt.Sprintf("%v", l.Value)
}

func (l BoolLiteral) Columns() []RowValue {
	return nil
}

func (l BoolLiteral) Volatility() VolatilityType {
	return Immutable
}

type StringLiteral struct {
	Value string
}

func (l StringLiteral) Repr() string {
	return fmt.Sprintf("%s", l.Value)
}

func (l StringLiteral) Columns() []RowValue {
	return nil
}

func (l StringLiteral) Volatility() VolatilityType {
	return Immutable
}

type AggFuncAppAST struct {
	Function   udf.UDF
	Expression FlatExpression
}

package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
	"time"
)

type evalTest struct {
	input    data.Value
	expected data.Value
}

func newTestContext() *core.Context {
	return &core.Context{
		Logger:       core.NewConsolePrintLogger(),
		Config:       core.Configuration{},
		SharedStates: core.NewDefaultSharedStateRegistry(),
	}
}

func TestEvaluators(t *testing.T) {
	testCases := getTestCases()
	reg := &testFuncRegistry{ctx: newTestContext()}

	for _, testCase := range testCases {
		testCase := testCase
		ast := testCase.ast
		Convey(fmt.Sprintf("Given the AST Expression %v", ast), t, func() {

			Convey("Then an Evaluator can be computed", func() {
				eval, err := ExpressionToEvaluator(ast, reg)
				So(err, ShouldBeNil)

				for i, tc := range testCase.inputs {
					input, expected := tc.input, tc.expected

					Convey(fmt.Sprintf("And when applied to input %v [%v]", input, i), func() {
						actual, err := eval.Eval(input)

						Convey(fmt.Sprintf("Then the result should be %v", expected), func() {
							i++
							if expected == nil {
								So(err, ShouldNotBeNil)
							} else {
								So(err, ShouldBeNil)
								So(actual, ShouldResemble, expected)
							}
						})
					})
				}

			})
		})
	}
}

func TestFuncAppConversion(t *testing.T) {
	Convey("Given a function registry", t, func() {
		reg := &testFuncRegistry{ctx: newTestContext()}

		Convey("When a function is known in the registry", func() {
			ast := parser.FuncAppAST{parser.FuncName("plusone"),
				parser.ExpressionsAST{[]parser.Expression{
					parser.RowValue{"", "a"},
				}}}

			Convey("Then we obtain an evaluatable funcApp", func() {
				eval, err := ExpressionToEvaluator(ast, reg)
				So(err, ShouldBeNil)
				So(eval, ShouldHaveSameTypeAs, &funcApp{})
			})
		})

		Convey("When the function is not known in the registry", func() {
			ast := parser.FuncAppAST{parser.FuncName("fun"),
				parser.ExpressionsAST{[]parser.Expression{
					parser.RowValue{"", "a"},
				}}}

			Convey("Then converting to an Evaluator fails", func() {
				_, err := ExpressionToEvaluator(ast, reg)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

// PlusOne is an example function that adds one to int and float Values.
// It panics if the input is Null and returns an error for any other
// type.
var (
	PlusOne = udf.VariadicFunc(func(ctx *core.Context, vs ...data.Value) (data.Value, error) {
		if len(vs) != 1 {
			err := fmt.Errorf("cannot use %d parameters for unary function", len(vs))
			return nil, err
		}
		v := vs[0]
		if v.Type() == data.TypeInt {
			i, _ := data.AsInt(v)
			return data.Int(i + 1), nil
		} else if v.Type() == data.TypeFloat {
			f, _ := data.AsFloat(v)
			return data.Float(f + 1.0), nil
		} else if v.Type() == data.TypeNull {
			panic("null!")
		}
		return nil, fmt.Errorf("cannot add 1 to %v", v)
	})
)

// testFuncRegistry returns the PlusOne function above for any parameter.
type testFuncRegistry struct {
	ctx *core.Context
}

func (tfr *testFuncRegistry) Context() *core.Context {
	return tfr.ctx
}

func (tfr *testFuncRegistry) Lookup(name string, arity int) (udf.UDF, error) {
	if name == "plusone" {
		return PlusOne, nil
	}
	return nil, fmt.Errorf("no such function: %s", name)
}

func getTestCases() []struct {
	ast    interface{}
	inputs []evalTest
} {
	now := time.Now()

	// these are all type combinations that are so incompatible that
	// they cannot be compared with respect to less/greater and also
	// cannot be added etc.
	incomparables := []evalTest{
		// not a map:
		{data.Int(17), nil},
		// keys not present:
		{data.Map{"x": data.Int(17)}, nil},
		// only left present => error
		{data.Map{"a": data.Bool(true)}, nil},
		{data.Map{"a": data.Int(17)}, nil},
		{data.Map{"a": data.Float(3.14)}, nil},
		{data.Map{"a": data.String("日本語")}, nil},
		{data.Map{"a": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
		// only right present => error
		{data.Map{"b": data.Bool(true)}, nil},
		{data.Map{"b": data.Int(17)}, nil},
		{data.Map{"b": data.Float(3.14)}, nil},
		{data.Map{"b": data.String("日本語")}, nil},
		{data.Map{"b": data.Blob("hoge")}, nil},
		{data.Map{"b": data.Timestamp(now)}, nil},
		{data.Map{"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"b": data.Map{"b": data.Int(3)}}, nil},
		// null vs. *
		{data.Map{"a": data.Null{},
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Null{},
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Null{},
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Null{},
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Null{},
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Null{},
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Null{},
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Null{},
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Null{},
			"b": data.Map{"b": data.Int(3)}}, nil},
		// bool vs *
		{data.Map{"a": data.Bool(true),
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Bool(true),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// int vs. *
		{data.Map{"a": data.Int(3),
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Int(3),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// float vs *
		{data.Map{"a": data.Float(3.14),
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Float(3.14),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// string vs *
		{data.Map{"a": data.String("hoge"),
			"b": data.Null{}}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.String("hoge"),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// blob vs *
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Blob("hoge"),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// timestamp vs *
		{data.Map{"a": data.Timestamp(now),
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Timestamp(now),
			"b": data.Map{"b": data.Int(3)}}, nil},
		// array vs *
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Array{data.Int(2)},
			"b": data.Map{"b": data.Int(3)}}, nil},
		// map vs *
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Null{}}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Bool(true)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Int(3)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Float(3.14)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.String("hoge")}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Blob("hoge")}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Timestamp(now)}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Array{data.Int(2)}}, nil},
		{data.Map{"a": data.Map{"b": data.Int(3)},
			"b": data.Map{"b": data.Int(3)}}, nil},
	}

	// we should check that every AST expression maps to
	// an evaluator with the correct behavior
	testCases := []struct {
		ast    interface{}
		inputs []evalTest
	}{
		// Literals should always be independent of the input data
		{parser.NumericLiteral{23},
			[]evalTest{
				{data.Int(17), data.Int(23)},
				{data.String(""), data.Int(23)},
			},
		},
		{parser.FloatLiteral{3.14},
			[]evalTest{
				{data.Int(17), data.Float(3.14)},
				{data.String(""), data.Float(3.14)},
			},
		},
		{parser.BoolLiteral{true},
			[]evalTest{
				{data.Int(17), data.Bool(true)},
				{data.String(""), data.Bool(true)},
			},
		},
		{parser.BoolLiteral{false},
			[]evalTest{
				{data.Int(17), data.Bool(false)},
				{data.String(""), data.Bool(false)},
			},
		},
		{parser.StringLiteral{"foo"},
			[]evalTest{
				{data.Int(17), data.String("foo")},
				{data.String(""), data.String("foo")},
			},
		},
		// Extracting the timestamp should find the timestamp at the
		// correct position
		{parser.RowMeta{"s", parser.TimestampMeta},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present, but wrong type
				{data.Map{"s:meta:TS": data.Int(17)}, nil},
				// key present and correct type
				{data.Map{"s:meta:TS": data.Timestamp(now)}, data.Timestamp(now)},
			},
		},
		// Access to columns/keys should return the same values
		{parser.RowValue{"", "a"},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present
				{data.Map{"a": data.Null{}}, data.Null{}},
				{data.Map{"a": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17)}, data.Int(17)},
				{data.Map{"a": data.Float(3.14)}, data.Float(3.14)},
				{data.Map{"a": data.String("日本語")}, data.String("日本語")},
				{data.Map{"a": data.Blob("hoge")}, data.Blob("hoge")},
				{data.Map{"a": data.Timestamp(now)}, data.Timestamp(now)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Array{data.Int(2)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Map{"b": data.Int(3)}},
			},
		},
		// Access to columns/keys should return the same values
		{parser.AliasAST{parser.RowValue{"", "a"}, "hoge"},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// key not present:
				{data.Map{"x": data.Int(17)}, nil},
				// key present
				{data.Map{"a": data.Null{}}, data.Null{}},
				{data.Map{"a": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17)}, data.Int(17)},
				{data.Map{"a": data.Float(3.14)}, data.Float(3.14)},
				{data.Map{"a": data.String("日本語")}, data.String("日本語")},
				{data.Map{"a": data.Blob("hoge")}, data.Blob("hoge")},
				{data.Map{"a": data.Timestamp(now)}, data.Timestamp(now)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Array{data.Int(2)}},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Map{"b": data.Int(3)}},
			},
		},
		/// Combined operations
		// Or
		{parser.BinaryOpAST{parser.Or, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left key present and evaluates to true => right is not necessary
				{data.Map{"a": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// only left key present and evaluates to false => error
				{data.Map{"a": data.Null{}}, nil},
				{data.Map{"a": data.Bool(false)}, nil},
				{data.Map{"a": data.Int(0)}, nil},
				{data.Map{"a": data.Float(0.0)}, nil},
				{data.Map{"a": data.String("")}, nil},
				{data.Map{"a": data.Blob("")}, nil},
				{data.Map{"a": data.Timestamp{}}, nil},
				{data.Map{"a": data.Array{}}, nil},
				{data.Map{"a": data.Map{}}, nil},
				// left key evalues to false and right to true => true
				{data.Map{"a": data.Int(0),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Int(0),
					"b": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// left key evalues to false and right to false => false
				{data.Map{"a": data.Int(0),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Int(0),
					"b": data.Map{}}, data.Bool(false)},
			},
		},
		// And
		{parser.BinaryOpAST{parser.And, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left key present and evaluates to false => right is not necessary
				{data.Map{"a": data.Null{}}, data.Bool(false)},
				{data.Map{"a": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Map{}}, data.Bool(false)},
				// only left key present and evaluates to true => error
				{data.Map{"a": data.Bool(true)}, nil},
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// left key evalues to true and right to true => true
				{data.Map{"a": data.Int(1),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// left key evalues to true and right to false => false
				{data.Map{"a": data.Int(1),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{}}, data.Bool(false)},
			},
		},
		/// Comparison Operations
		// Equal
		{parser.BinaryOpAST{parser.Equal, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left present => error
				{data.Map{"a": data.Bool(true)}, nil},
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// only right present => error
				{data.Map{"b": data.Bool(true)}, nil},
				{data.Map{"b": data.Int(17)}, nil},
				{data.Map{"b": data.Float(3.14)}, nil},
				{data.Map{"b": data.String("日本語")}, nil},
				{data.Map{"b": data.Blob("hoge")}, nil},
				{data.Map{"b": data.Timestamp(now)}, nil},
				{data.Map{"b": data.Array{data.Int(2)}}, nil},
				{data.Map{"b": data.Map{"b": data.Int(3)}}, nil},
				// left and right present and equal => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(17),
					"b": data.Int(17)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.String("日本語"),
					"b": data.String("日本語")}, data.Bool(true)},
				{data.Map{"a": data.Blob("hoge"),
					"b": data.Blob("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(true)},
				{data.Map{"a": data.Array{data.Int(2)},
					"b": data.Array{data.Int(2)}}, data.Bool(true)},
				{data.Map{"a": data.Map{"b": data.Int(3)},
					"b": data.Map{"b": data.Int(3)}}, data.Bool(true)},
				// left and right present and not equal => false
				{data.Map{"a": data.Int(1),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(0.0)}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.String("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("")}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{}}, data.Bool(false)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{}}, data.Bool(false)},
			},
		},
		// Less
		{parser.BinaryOpAST{parser.Less, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => true
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, data.Bool(true)},
				// left and right present and comparable and equal => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and comparable and left is greater => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(false)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(time.Now()),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// LessOrEqual
		{parser.BinaryOpAST{parser.LessOrEqual, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => true
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, data.Bool(true)},
				// left and right present and comparable and equal => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and comparable and left is greater => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(false)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(time.Now()),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Greater
		{parser.BinaryOpAST{parser.Greater, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => false
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, data.Bool(false)},
				// left and right present and comparable and equal => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(false)},
				// left and right present and comparable and left is greater => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(true)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(time.Now()),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// GreaterOrEqual
		{parser.BinaryOpAST{parser.GreaterOrEqual, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and comparable and left is less => false
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Bool(false)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, data.Bool(false)},
				// left and right present and comparable and equal => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and comparable and left is greater => true
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(false)}, data.Bool(true)},
				{data.Map{"a": data.Int(3),
					"b": data.Int(2)}, data.Bool(true)},
				{data.Map{"a": data.Int(4),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(3)}, data.Bool(true)},
				{data.Map{"a": data.Float(3.15),
					"b": data.Float(3.14)}, data.Bool(true)},
				{data.Map{"a": data.String("hogee"),
					"b": data.String("hoge")}, data.Bool(true)},
				{data.Map{"a": data.Timestamp(time.Now()),
					"b": data.Timestamp(now)}, data.Bool(true)},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// NotEqual
		{parser.BinaryOpAST{parser.NotEqual, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			[]evalTest{
				// not a map:
				{data.Int(17), nil},
				// keys not present:
				{data.Map{"x": data.Int(17)}, nil},
				// only left present => error
				{data.Map{"a": data.Bool(true)}, nil},
				{data.Map{"a": data.Int(17)}, nil},
				{data.Map{"a": data.Float(3.14)}, nil},
				{data.Map{"a": data.String("日本語")}, nil},
				{data.Map{"a": data.Blob("hoge")}, nil},
				{data.Map{"a": data.Timestamp(now)}, nil},
				{data.Map{"a": data.Array{data.Int(2)}}, nil},
				{data.Map{"a": data.Map{"b": data.Int(3)}}, nil},
				// only right present => error
				{data.Map{"b": data.Bool(true)}, nil},
				{data.Map{"b": data.Int(17)}, nil},
				{data.Map{"b": data.Float(3.14)}, nil},
				{data.Map{"b": data.String("日本語")}, nil},
				{data.Map{"b": data.Blob("hoge")}, nil},
				{data.Map{"b": data.Timestamp(now)}, nil},
				{data.Map{"b": data.Array{data.Int(2)}}, nil},
				{data.Map{"b": data.Map{"b": data.Int(3)}}, nil},
				// left and right present and equal => false
				{data.Map{"a": data.Bool(true),
					"b": data.Bool(true)}, data.Bool(false)},
				{data.Map{"a": data.Int(17),
					"b": data.Int(17)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.14)}, data.Bool(false)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.0)}, data.Bool(false)},
				{data.Map{"a": data.Float(3.0),
					"b": data.Int(3)}, data.Bool(false)},
				{data.Map{"a": data.String("日本語"),
					"b": data.String("日本語")}, data.Bool(false)},
				{data.Map{"a": data.Blob("hoge"),
					"b": data.Blob("hoge")}, data.Bool(false)},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(now)}, data.Bool(false)},
				{data.Map{"a": data.Array{data.Int(2)},
					"b": data.Array{data.Int(2)}}, data.Bool(false)},
				{data.Map{"a": data.Map{"b": data.Int(3)},
					"b": data.Map{"b": data.Int(3)}}, data.Bool(false)},
				// left and right present and not equal => true
				{data.Map{"a": data.Int(1),
					"b": data.Bool(false)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Int(0)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Float(0.0)}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.String("")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Blob("")}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Timestamp{}}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Array{}}, data.Bool(true)},
				{data.Map{"a": data.Int(1),
					"b": data.Map{}}, data.Bool(true)},
			},
		},
		/// Computational Operations
		// Plus
		{parser.BinaryOpAST{parser.Plus, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be added
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(5)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) + 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 + float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(3.14 + 3.15)},
				// left and right present and cannot be added
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Minus
		{parser.BinaryOpAST{parser.Minus, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be subtracted
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(-1)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) - 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 - float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(float64(3.14) - 3.15)},
				// left and right present and cannot be subtracted
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Multiply
		{parser.BinaryOpAST{parser.Multiply, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be multiplied
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(6)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) * 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 * float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(float64(3.14) * 3.15)},
				// left and right present and cannot be multiplied
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Divide
		{parser.BinaryOpAST{parser.Divide, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be divided
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(0)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(float64(3) / 3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14 / float64(4))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(float64(3.14) / 3.15)},
				// division by zero
				{data.Map{"a": data.Int(2),
					"b": data.Int(0)}, nil},
				{data.Map{"a": data.Int(3),
					"b": data.Float(0)}, data.Float(math.Inf(1))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(0)}, data.Float(math.Inf(1))},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(0)}, data.Float(math.Inf(1))},
				// left and right present and cannot be divided
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		// Modulo
		{parser.BinaryOpAST{parser.Modulo, parser.RowValue{"", "a"}, parser.RowValue{"", "b"}},
			append([]evalTest{
				// left and right present and can be moduled
				{data.Map{"a": data.Int(2),
					"b": data.Int(3)}, data.Int(2)},
				{data.Map{"a": data.Int(3),
					"b": data.Float(3.14)}, data.Float(3.0)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Int(4)}, data.Float(3.14)},
				{data.Map{"a": data.Float(3.14),
					"b": data.Float(3.15)}, data.Float(3.14)},
				// modulo by zero
				{data.Map{"a": data.Int(2),
					"b": data.Int(0)}, nil},
				// TODO add a way to check for IsNaN()
				/*
					{data.Map{"a": data.Int(3),
						"b": data.Float(0)}, data.Float(math.NaN())},
					{data.Map{"a": data.Float(3.14),
						"b": data.Int(0)}, data.Float(math.NaN())},
					{data.Map{"a": data.Float(3.14),
						"b": data.Float(0)}, data.Float(math.NaN())},
				*/
				// left and right present and cannot be moduled
				{data.Map{"a": data.Bool(false),
					"b": data.Bool(true)}, nil},
				{data.Map{"a": data.String("hoge"),
					"b": data.String("hogee")}, nil},
				{data.Map{"a": data.Timestamp(now),
					"b": data.Timestamp(time.Now())}, nil},
				// left and right present and not comparable => error
			}, incomparables...),
		},
		/// Function Application
		{parser.FuncAppAST{parser.FuncName("plusone"),
			parser.ExpressionsAST{[]parser.Expression{parser.RowValue{"", "a"}}}},
			// NB. This only tests the behavior of funcApp.Eval.
			// It does *not* test the function registry, mismatch
			// in parameter counts or any particular function.
			[]evalTest{
				// function returns good result
				{data.Map{"a": data.Int(16)}, data.Int(17)},
				{data.Map{"a": data.Float(16.0)}, data.Float(17.0)},
				// function errors
				{data.Map{"x": data.Int(17)}, nil},
				{data.Map{"a": data.Bool(false)}, nil},
				// function panics
				{data.Map{"a": data.Null{}}, nil},
			},
		},
	}
	return testCases
}

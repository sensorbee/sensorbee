package builtin

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/data"
	"testing"
	"time"
)

func TestUnaryAggregateFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfUnaryTestCaseInput{
		{data.Null{}, nil},
		{data.Blob{}, nil},
		{data.Bool(true), nil},
		{data.Float(2.7), nil},
		{data.Int(3), nil},
		{data.Map{}, nil},
		{data.String("hoge"), nil},
		{data.Timestamp(someTime), nil},
	}

	udfUnaryTestCases := []udfUnaryTestCase{
		{"count", countFunc, []udfUnaryTestCaseInput{
			// empty array: 0
			{data.Array{}, data.Int(0)},
			// normal inputs
			{data.Array{data.Int(7)}, data.Int(1)},
			{data.Array{data.Int(7), data.Int(3)}, data.Int(2)},
			// do not count 0
			{data.Array{data.Int(7), data.Null{}, data.Int(3)}, data.Int(2)},
		}},
		{"array_agg", arrayAggFunc, []udfUnaryTestCaseInput{
			// empty array: Null
			{data.Array{}, data.Null{}},
			// normal inputs
			{data.Array{data.Int(7), data.Int(3)},
				data.Array{data.Int(7), data.Int(3)}},
			{data.Array{data.Int(7), data.Null{}, data.Int(3)},
				data.Array{data.Int(7), data.Null{}, data.Int(3)}},
		}},
		{"avg", avgFunc, []udfUnaryTestCaseInput{
			// empty array: Null
			{data.Array{}, data.Null{}},
			// normal inputs
			{data.Array{data.Int(7), data.Int(3)}, data.Float(5.0)},
			{data.Array{data.Int(7), data.Null{}, data.Float(3.0)}, data.Float(5.0)},
			// incompatible data
			{data.Array{data.Int(7), data.Timestamp(someTime)}, nil},
		}},
		{"bool_and", boolAndFunc, []udfUnaryTestCaseInput{
			// empty array: Null
			{data.Array{}, data.Null{}},
			// normal inputs
			{data.Array{data.Bool(true), data.Bool(true)}, data.Bool(true)},
			{data.Array{data.Bool(false), data.Bool(true)}, data.Bool(false)},
			{data.Array{data.Bool(true), data.Bool(false)}, data.Bool(false)},
			{data.Array{data.Bool(false), data.Bool(false)}, data.Bool(false)},
			// normal inputs with nulls
			{data.Array{data.Null{}, data.Bool(true), data.Bool(true)}, data.Bool(true)},
			{data.Array{data.Bool(false), data.Null{}, data.Bool(true)}, data.Bool(false)},
			{data.Array{data.Bool(true), data.Bool(false), data.Null{}}, data.Bool(false)},
			// incompatible data
			{data.Array{data.Bool(true), data.Int(7)}, nil},
			{data.Array{data.Int(7), data.Bool(true)}, nil},
		}},
		{"bool_or", boolOrFunc, []udfUnaryTestCaseInput{
			// empty array: Null
			{data.Array{}, data.Null{}},
			// normal inputs
			{data.Array{data.Bool(true), data.Bool(true)}, data.Bool(true)},
			{data.Array{data.Bool(false), data.Bool(true)}, data.Bool(true)},
			{data.Array{data.Bool(true), data.Bool(false)}, data.Bool(true)},
			{data.Array{data.Bool(false), data.Bool(false)}, data.Bool(false)},
			// normal inputs with nulls
			{data.Array{data.Null{}, data.Bool(true), data.Bool(true)}, data.Bool(true)},
			{data.Array{data.Bool(false), data.Null{}, data.Bool(true)}, data.Bool(true)},
			{data.Array{data.Bool(true), data.Bool(false), data.Null{}}, data.Bool(true)},
			{data.Array{data.Bool(false), data.Bool(false), data.Null{}}, data.Bool(false)},
			// incompatible data
			{data.Array{data.Bool(true), data.Int(7)}, nil},
			{data.Array{data.Int(7), data.Bool(true)}, nil},
		}},
		{"max", maxFunc, []udfUnaryTestCaseInput{
			// empty array: Null
			{data.Array{}, data.Null{}},
			/// normal inputs
			// single values
			{data.Array{data.Float(2.3)}, data.Float(2.3)},
			{data.Array{data.Int(2)}, data.Int(2)},
			// homogeneous inputs
			{data.Array{data.Int(2), data.Int(3)}, data.Int(3)},
			{data.Array{data.Int(3), data.Int(2)}, data.Int(3)},
			{data.Array{data.Float(2.3), data.Float(3.2)}, data.Float(3.2)},
			{data.Array{data.Float(3.2), data.Float(2.3)}, data.Float(3.2)},
			// mixed type
			{data.Array{data.Float(2.3), data.Int(3)}, data.Int(3)},
			{data.Array{data.Int(3), data.Float(2.3)}, data.Int(3)},
			{data.Array{data.Int(2), data.Float(3.2)}, data.Float(3.2)},
			{data.Array{data.Float(3.2), data.Int(2)}, data.Float(3.2)},
			// same value: int wins
			{data.Array{data.Float(3.0), data.Int(3)}, data.Int(3)},
			{data.Array{data.Int(3), data.Float(3.0)}, data.Int(3)},
			// extreme values
			{data.Array{data.Float(math.MaxFloat64), data.Int(math.MaxInt64)}, data.Float(math.MaxFloat64)},
			{data.Array{data.Float(-math.MaxFloat64), data.Int(math.MinInt64)}, data.Int(math.MinInt64)},
			// incompatible data
			{data.Array{data.Int(7), data.Timestamp(someTime)}, nil},
		}},
		{"min", minFunc, []udfUnaryTestCaseInput{
			// empty array: Null
			{data.Array{}, data.Null{}},
			/// normal inputs
			// single values
			{data.Array{data.Float(2.3)}, data.Float(2.3)},
			{data.Array{data.Int(2)}, data.Int(2)},
			// homogeneous inputs
			{data.Array{data.Int(2), data.Int(3)}, data.Int(2)},
			{data.Array{data.Int(3), data.Int(2)}, data.Int(2)},
			{data.Array{data.Float(2.3), data.Float(3.2)}, data.Float(2.3)},
			{data.Array{data.Float(3.2), data.Float(2.3)}, data.Float(2.3)},
			// mixed type
			{data.Array{data.Float(2.3), data.Int(3)}, data.Float(2.3)},
			{data.Array{data.Int(3), data.Float(2.3)}, data.Float(2.3)},
			{data.Array{data.Int(2), data.Float(3.2)}, data.Int(2)},
			{data.Array{data.Float(3.2), data.Int(2)}, data.Int(2)},
			// same value: int wins
			{data.Array{data.Float(3.0), data.Int(3)}, data.Int(3)},
			{data.Array{data.Int(3), data.Float(3.0)}, data.Int(3)},
			// extreme values
			{data.Array{data.Float(-math.MaxFloat64), data.Int(math.MinInt64)}, data.Float(-math.MaxFloat64)},
			{data.Array{data.Float(math.MaxFloat64), data.Int(math.MaxInt64)}, data.Int(math.MaxInt64)},
			// incompatible data
			{data.Array{data.Int(7), data.Timestamp(someTime)}, nil},
		}},
		{"sum", sumFunc, []udfUnaryTestCaseInput{
			// empty array: Null
			{data.Array{}, data.Null{}},
			/// normal inputs
			// single values
			{data.Array{data.Float(2.3)}, data.Float(2.3)},
			{data.Array{data.Int(2)}, data.Int(2)},
			// homogeneous inputs
			{data.Array{data.Int(2), data.Int(3)}, data.Int(5)},
			{data.Array{data.Int(3), data.Int(2)}, data.Int(5)},
			{data.Array{data.Float(2.3), data.Float(3.2)}, data.Float(5.5)},
			{data.Array{data.Float(3.2), data.Float(2.3)}, data.Float(5.5)},
			// mixed type
			{data.Array{data.Float(2.3), data.Int(3)}, data.Float(5.3)},
			{data.Array{data.Int(3), data.Float(2.3)}, data.Float(5.3)},
			/// overflow
			// the integer case is commutative
			{data.Array{data.Int(math.MaxInt64), data.Int(10), data.Int(-20)}, data.Int(math.MaxInt64 - 10)},
			{data.Array{data.Int(10), data.Int(math.MaxInt64), data.Int(-20)}, data.Int(math.MaxInt64 - 10)},
			{data.Array{data.Int(10), data.Int(-20), data.Int(math.MaxInt64)}, data.Int(math.MaxInt64 - 10)},
			// the float case is not
			{data.Array{data.Float(math.MaxFloat64 / 4), data.Float(-math.MaxFloat64 / 2), data.Float(math.MaxFloat64)},
				data.Float(3. / 4 * math.MaxFloat64)},
			{data.Array{data.Float(math.MaxFloat64 / 4), data.Float(math.MaxFloat64), data.Float(-math.MaxFloat64 / 2)},
				data.Float(math.Inf(1))},
			// incompatible data
			{data.Array{data.Int(7), data.Timestamp(someTime)}, nil},
		}},
	}

	for _, testCase := range udfUnaryTestCases {
		f := testCase.f
		allInputs := append(testCase.inputs, invalidInputs...)

		Convey(fmt.Sprintf("Given the %s function", testCase.name), t, func() {
			Convey("Then it should be an aggregate in the first parameter", func() {
				So(f.IsAggregationParameter(0), ShouldBeTrue)
			})

			for i, tc := range allInputs {
				tc := tc

				Convey(fmt.Sprintf("[%d] When evaluating it on %s (%T)", i, tc.input, tc.input), func() {
					val, err := f.Call(nil, tc.input)

					if tc.expected == nil {
						Convey("Then evaluation should fail", func() {
							So(err, ShouldNotBeNil)
						})
					} else {
						Convey(fmt.Sprintf("Then the result should be %s", tc.expected), func() {
							So(err, ShouldBeNil)
							if val.Type() == data.TypeFloat && tc.expected.Type() == data.TypeFloat {
								f_actual, _ := data.AsFloat(val)
								f_expected, _ := data.AsFloat(tc.expected)
								if math.IsNaN(f_expected) {
									So(math.IsNaN(f_actual), ShouldBeTrue)
								} else if math.IsInf(f_expected, 0) {
									pos := 1
									if f_expected < 0 {
										pos = -1
									}
									So(math.IsInf(f_actual, pos), ShouldBeTrue)
								} else {
									So(val, ShouldAlmostEqual, tc.expected, 0.0000001)
								}
							} else {
								So(val, ShouldResemble, tc.expected)
							}
						})
					}
				})
			}

			Convey("Then it should equal the one in the default registry", func() {
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 1)
				if dispatcher, ok := regFun.(*arityDispatcher); ok {
					regFun = dispatcher.unary
				}
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

func TestBinaryAggregateFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfBinaryTestCaseInput{
		{data.Null{}, data.Int(1), nil},
		{data.Float(1), data.Null{}, nil},
		{data.Blob{}, data.Blob{}, nil},
		{data.Bool(true), data.Bool(true), nil},
		{data.Float(2.7), data.Float(2.7), nil},
		{data.Int(3), data.Int(3), nil},
		{data.Map{}, data.Map{}, nil},
		{data.String("hoge"), data.String("hoge"), nil},
		{data.Timestamp(someTime), data.Timestamp(someTime), nil},
	}

	udfBinaryTestCases := []udfBinaryTestCase{
		{"json_object_agg", jsonObjectAggFunc, []udfBinaryTestCaseInput{
			{data.Array{}, data.Array{}, data.Null{}},
			// normal cases
			{data.Array{data.String("foo")}, data.Array{data.Int(7)},
				data.Map{"foo": data.Int(7)}},
			{data.Array{data.String("foo"), data.String("bar")},
				data.Array{data.Int(7), data.Int(3)},
				data.Map{"foo": data.Int(7), "bar": data.Int(3)}},
			{data.Array{data.String("foo"), data.Null{}, data.String("bar")},
				data.Array{data.Int(7), data.Null{}, data.Null{}},
				data.Map{"foo": data.Int(7), "bar": data.Null{}}},
			/// fail cases
			// duplicate keys
			{data.Array{data.String("foo"), data.String("foo")},
				data.Array{data.Int(7), data.Int(3)}, nil},
			// different length
			{data.Array{data.String("foo")},
				data.Array{data.Int(7), data.Int(3)}, nil},
			{data.Array{data.String("foo"), data.String("bar")},
				data.Array{data.Int(7)}, nil},
			// key is null
			{data.Array{data.String("foo"), data.Null{}},
				data.Array{data.Int(7), data.Int(3)}, nil},
			// key is non-string
			{data.Array{data.String("foo"), data.Int(17)},
				data.Array{data.Int(7), data.Int(3)}, nil},
		}},
		{"string_agg", stringAggFunc, []udfBinaryTestCaseInput{
			{data.Array{}, data.String(", "), data.Null{}},
			// normal cases
			{data.Array{data.String("foo")}, data.String(", "),
				data.String("foo")},
			{data.Array{data.String("foo"), data.String("bar")}, data.String(", "),
				data.String("foo, bar")},
			{data.Array{data.String("foo"), data.Null{}, data.String("bar")}, data.String(", "),
				data.String("foo, bar")},
			{data.Array{data.Null{}, data.String("foo"), data.String("bar")}, data.String(", "),
				data.String("foo, bar")},
			/// fail cases
			// delimiter is null
			{data.Array{data.String("foo"), data.String("bar")}, data.Null{}, nil},
			// delimiter is non-string
			{data.Array{data.String("foo"), data.String("bar")}, data.Int(7), nil},
			// array contains non-string
			{data.Array{data.String("foo"), data.Int(7)}, data.String(", "), nil},
		}},
	}

	for _, testCase := range udfBinaryTestCases {
		f := testCase.f
		allInputs := append(testCase.inputs, invalidInputs...)

		Convey(fmt.Sprintf("Given the %s function", testCase.name), t, func() {
			for _, tc := range allInputs {
				tc := tc

				Convey(fmt.Sprintf("When evaluating it on %s (%T) and %s (%T)",
					tc.input1, tc.input1, tc.input2, tc.input2), func() {
					val, err := f.Call(nil, tc.input1, tc.input2)

					if tc.expected == nil {
						Convey("Then evaluation should fail", func() {
							So(err, ShouldNotBeNil)
						})
					} else {
						Convey(fmt.Sprintf("Then the result should be %s", tc.expected), func() {
							So(err, ShouldBeNil)
							if val.Type() == data.TypeFloat && tc.expected.Type() == data.TypeFloat {
								f_actual, _ := data.AsFloat(val)
								f_expected, _ := data.AsFloat(tc.expected)
								if math.IsNaN(f_expected) {
									So(math.IsNaN(f_actual), ShouldBeTrue)
								} else {
									So(val, ShouldAlmostEqual, tc.expected, 0.0000001)
								}
							} else {
								So(val, ShouldResemble, tc.expected)
							}
						})
					}
				})
			}

			Convey("Then it should equal the one in the default registry", func() {
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 2)
				if dispatcher, ok := regFun.(*arityDispatcher); ok {
					regFun = dispatcher.binary
				}
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

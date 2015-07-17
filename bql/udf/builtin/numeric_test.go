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

type udfUnaryTestCase struct {
	name   string
	f      udf.UDF
	inputs []udfUnaryTestCaseInput
}

type udfUnaryTestCaseInput struct {
	input    data.Value
	expected data.Value
}

func TestUnaryNumericFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfUnaryTestCaseInput{
		// NULL input -> NULL output
		{data.Null{}, data.Null{}},
		// cannot process the following
		{data.Array{}, nil},
		{data.Blob{}, nil},
		{data.Bool(true), nil},
		{data.Map{}, nil},
		{data.String("hoge"), nil},
		{data.Timestamp(someTime), nil},
	}

	udfUnaryTestCases := []udfUnaryTestCase{
		{"abs", absFunc, []udfUnaryTestCaseInput{
			{data.Int(7), data.Int(7)},
			{data.Int(-7), data.Int(7)},
			{data.Float(2.3), data.Float(2.3)},
			{data.Float(-2.3), data.Float(2.3)},
		}},
		{"cbrt", cbrtFunc, []udfUnaryTestCaseInput{
			{data.Int(27), data.Float(3.0)},
			{data.Int(-27), data.Float(-3.0)},
			{data.Float(27.0), data.Float(3.0)},
			{data.Float(-27.0), data.Float(-3.0)},
		}},
		{"ceil", ceilFunc, []udfUnaryTestCaseInput{
			{data.Int(27), data.Int(27)},
			{data.Int(-27), data.Int(-27)},
			{data.Float(42.8), data.Float(43.0)},
			{data.Float(-42.8), data.Float(-42.0)},
		}},
		{"degrees", degreesFunc, []udfUnaryTestCaseInput{
			{data.Int(1), data.Float(57.29577951308232)},
			{data.Float(0.5), data.Float(28.6478897565412)},
		}},
		{"exp", expFunc, []udfUnaryTestCaseInput{
			{data.Int(1), data.Float(math.E)},
			{data.Float(math.Log(1.0)), data.Float(1.0)},
		}},
		{"floor", floorFunc, []udfUnaryTestCaseInput{
			{data.Int(27), data.Int(27)},
			{data.Int(-27), data.Int(-27)},
			{data.Float(42.8), data.Float(42.0)},
			{data.Float(-42.8), data.Float(-43.0)},
		}},
	}

	for _, testCase := range udfUnaryTestCases {
		f := testCase.f
		allInputs := append(testCase.inputs, invalidInputs...)

		Convey(fmt.Sprintf("Given the %s function", testCase.name), t, func() {
			for _, tc := range allInputs {
				tc := tc

				Convey(fmt.Sprintf("When evaluating it on %s (%T)", tc.input, tc.input), func() {
					val, err := f.Call(nil, tc.input)

					if tc.expected == nil {
						Convey("Then evaluation should fail", func() {
							So(err, ShouldNotBeNil)
						})
					} else {
						Convey(fmt.Sprintf("Then the result should be %s", tc.expected), func() {
							So(err, ShouldBeNil)
							if val.Type() == data.TypeFloat && tc.expected.Type() == data.TypeFloat {
								So(val, ShouldAlmostEqual, tc.expected, 0.0000001)
							} else {
								So(val, ShouldResemble, tc.expected)
							}
						})
					}
				})
			}

			Convey("Then it should equal in the default registry", func() {
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 1)
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

type udfBinaryTestCase struct {
	name   string
	f      udf.UDF
	inputs []udfBinaryTestCaseInput
}

type udfBinaryTestCaseInput struct {
	input1   data.Value
	input2   data.Value
	expected data.Value
}

func TestBinaryNumericFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfBinaryTestCaseInput{
		// NULL input -> NULL output
		{data.Null{}, data.Int(1), data.Null{}},
		{data.Float(1), data.Null{}, data.Null{}},
		// cannot process the following
		{data.Array{}, data.Array{}, nil},
		{data.Blob{}, data.Blob{}, nil},
		{data.Bool(true), data.Bool(true), nil},
		{data.Map{}, data.Map{}, nil},
		{data.String("hoge"), data.String("hoge"), nil},
		{data.Timestamp(someTime), data.Timestamp(someTime), nil},
		// also cannot process mixed types
		{data.Int(7), data.Float(7), nil},
		{data.Float(7), data.Int(7), nil},
	}

	udfBinaryTestCases := []udfBinaryTestCase{
		{"div", divFunc, []udfBinaryTestCaseInput{
			{data.Int(9), data.Int(4), data.Int(2)},
			{data.Int(9), data.Int(-4), data.Int(-2)},
			{data.Int(9), data.Int(0), nil},
			{data.Float(2.7), data.Float(1.3), data.Int(2)},
			{data.Float(2.7), data.Float(-1.3), data.Int(-2)},
			{data.Float(-2.7), data.Float(1.3), data.Int(-2)},
			{data.Float(-2.7), data.Float(-1.3), data.Int(2)},
			{data.Float(2.7), data.Float(0.0), nil},
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
								So(val, ShouldAlmostEqual, tc.expected, 0.0000001)
							} else {
								So(val, ShouldResemble, tc.expected)
							}
						})
					}
				})
			}

			Convey("Then it should equal in the default registry", func() {
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 2)
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

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

func TestNullaryNumericFuncs(t *testing.T) {
	Convey("pi() should return pi", t, func() {
		f := piFunc
		val, err := f.Call(nil)
		So(err, ShouldBeNil)
		So(val, ShouldAlmostEqual, math.Pi, 0.00000001)
	})

	Convey("piFunc should equal the one the default registry", t, func() {
		regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup("pi", 0)
		So(err, ShouldBeNil)
		So(regFun, ShouldHaveSameTypeAs, piFunc)
	})

	Convey("random() should return a random number", t, func() {
		f := randomFunc
		val, err := f.Call(nil)
		So(err, ShouldBeNil)
		So(val, ShouldBeGreaterThanOrEqualTo, 0.0)
		So(val, ShouldBeLessThan, 1.0)
	})

	Convey("random should equal the one the default registry", t, func() {
		regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup("random", 0)
		So(err, ShouldBeNil)
		So(regFun, ShouldHaveSameTypeAs, randomFunc)
	})
}

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
		{"ln", lnFunc, []udfUnaryTestCaseInput{
			{data.Float(-1.0), data.Float(math.NaN())},
			{data.Int(2), data.Float(0.693147180559945)},
			{data.Float(math.E), data.Float(1.0)},
		}},
		{"log", logFunc, []udfUnaryTestCaseInput{
			{data.Float(-1.0), data.Float(math.NaN())},
			{data.Int(10), data.Float(1.0)},
			{data.Float(100.0), data.Float(2.0)},
		}},
		{"radians", radiansFunc, []udfUnaryTestCaseInput{
			{data.Int(45), data.Float(0.785398163397448)},
			{data.Float(28.6478897565412), data.Float(0.5)},
		}},
		{"round", roundFunc, []udfUnaryTestCaseInput{
			{data.Int(27), data.Int(27)},
			{data.Int(-27), data.Int(-27)},
			{data.Float(42.8), data.Float(43.0)},
			{data.Float(-42.8), data.Float(-43.0)},
			{data.Float(42.5), data.Float(43.0)},
			{data.Float(-42.5), data.Float(-43.0)},
			{data.Float(42.2), data.Float(42.0)},
			{data.Float(-42.2), data.Float(-42.0)},
		}},
		{"sign", signFunc, []udfUnaryTestCaseInput{
			{data.Int(27), data.Int(1)},
			{data.Int(0), data.Int(0)},
			{data.Int(-27), data.Int(-1)},
			{data.Float(42.8), data.Int(1)},
			{data.Float(0.0), data.Int(0)},
			{data.Float(-42.8), data.Int(-1)},
		}},
		{"sqrt", sqrtFunc, []udfUnaryTestCaseInput{
			{data.Int(2), data.Float(math.Sqrt2)},
			{data.Int(-2), data.Float(math.NaN())},
			{data.Float(9.0), data.Float(3.0)},
			{data.Float(-9.0), data.Float(math.NaN())},
		}},
		{"trunc", truncFunc, []udfUnaryTestCaseInput{
			{data.Int(27), data.Int(27)},
			{data.Int(-27), data.Int(-27)},
			{data.Float(42.8), data.Float(42.0)},
			{data.Float(-42.8), data.Float(-42.0)},
		}},
		{"setseed", setseedFunc, []udfUnaryTestCaseInput{
			{data.Int(27), nil},
			{data.Float(42.8), nil},
			{data.Float(-1.0), data.Null{}},
			{data.Float(0.5), data.Null{}},
			{data.Float(1.0), data.Null{}},
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
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 1)
				if dispatcher, ok := regFun.(*unaryBinaryDispatcher); ok {
					regFun = dispatcher.unary
				}
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
			{data.Float(2.7), data.Float(1.3), data.Float(2.0)},
			{data.Float(2.7), data.Float(-1.3), data.Float(-2.0)},
			{data.Float(-2.7), data.Float(1.3), data.Float(-2.0)},
			{data.Float(-2.7), data.Float(-1.3), data.Float(2.0)},
			{data.Float(2.7), data.Float(0.0), data.Float(math.NaN())},
		}},
		{"log", logBaseFunc, []udfBinaryTestCaseInput{
			{data.Float(2.0), data.Float(-1.0), data.Float(math.NaN())},
			{data.Int(2), data.Int(64), data.Float(6.0)},
			{data.Float(1.5), data.Float(2.25), data.Float(2.0)},
		}},
		{"mod", modFunc, []udfBinaryTestCaseInput{
			{data.Int(9), data.Int(4), data.Int(1)},
			{data.Int(9), data.Int(-4), data.Int(1)},
			{data.Int(-9), data.Int(4), data.Int(-1)},
			{data.Int(-9), data.Int(-4), data.Int(-1)},
			{data.Int(9), data.Int(0), nil},
			{data.Float(2.7), data.Float(1.3), data.Float(0.1)},
			{data.Float(2.7), data.Float(-1.3), data.Float(0.1)},
			{data.Float(-2.7), data.Float(1.3), data.Float(-0.1)},
			{data.Float(-2.7), data.Float(-1.3), data.Float(-0.1)},
			{data.Float(2.7), data.Float(0.0), data.Float(math.NaN())},
		}},
		{"log", logBaseFunc, []udfBinaryTestCaseInput{
			{data.Float(2.0), data.Float(-1.0), data.Float(math.NaN())},
			{data.Int(2), data.Int(64), data.Float(6.0)},
			{data.Float(1.5), data.Float(2.25), data.Float(2.0)},
		}},
		{"power", powFunc, []udfBinaryTestCaseInput{
			{data.Int(2), data.Int(6), data.Float(64.0)},
			{data.Int(-2), data.Int(3), data.Float(-8.0)},
			{data.Int(2), data.Int(-6), data.Float(1. / 64.0)},
			{data.Float(9.0), data.Float(3.0), data.Float(729.0)},
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
				if dispatcher, ok := regFun.(*unaryBinaryDispatcher); ok {
					regFun = dispatcher.binary
				}
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

type udf4aryTestCase struct {
	name   string
	f      udf.UDF
	inputs []udf4aryTestCaseInput
}

type udf4aryTestCaseInput struct {
	input1   data.Value
	input2   data.Value
	input3   data.Value
	input4   data.Value
	expected data.Value
}

func Test4aryNumericFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udf4aryTestCaseInput{
		// NULL input -> NULL output
		{data.Int(1), data.Null{}, data.Null{}, data.Null{}, data.Null{}},
		{data.Null{}, data.Int(1), data.Null{}, data.Null{}, data.Null{}},
		{data.Null{}, data.Null{}, data.Int(1), data.Null{}, data.Null{}},
		{data.Null{}, data.Null{}, data.Null{}, data.Int(1), data.Null{}},
		// cannot process the following
		{data.Array{}, data.Array{}, data.Array{}, data.Array{}, nil},
		{data.Blob{}, data.Blob{}, data.Blob{}, data.Blob{}, nil},
		{data.Bool(true), data.Bool(true), data.Bool(true), data.Bool(true), nil},
		{data.Map{}, data.Map{}, data.Map{}, data.Map{}, nil},
		{data.String("hoge"), data.String("hoge"), data.String("hoge"), data.String("hoge"), nil},
		{data.Timestamp(someTime), data.Timestamp(someTime), data.Timestamp(someTime), data.Timestamp(someTime), nil},
	}

	udf4aryTestCases := []udf4aryTestCase{
		{"width_bucket", widthBucketFunc, []udf4aryTestCaseInput{
			//     -1.5   0.0   1.5   3.0
			//       v     v     v     v
			//  0    |  1  |  2  |  3  |    4
			// -------------------------------
			// left of left border
			{data.Float(-2.5), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(0)},
			// on left border -> belongs to first bucket
			{data.Float(-1.5), data.Float(-1.5), data.Int(3.0), data.Int(3), data.Int(1)},
			// in first bucket
			{data.Float(-0.5), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(1)},
			// on bucket border -> belongs to right bucket
			{data.Int(0), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(2)},
			// in second bucket
			{data.Int(1), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(2)},
			// on bucket border -> belongs fo right bucket
			{data.Float(1.5), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(3)},
			// in third bucket
			{data.Int(2), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(3)},
			// on right bucket border -> belongs to outside
			{data.Float(3.0), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(4)},
			// right of right border
			{data.Int(5), data.Float(-1.5), data.Float(3.0), data.Int(3), data.Int(4)},
			// invalid: border limits broken
			{data.Float(5), data.Int(3), data.Float(-2.0), data.Int(3), nil},
			{data.Float(5), data.Int(3), data.Float(3.0), data.Int(3), nil},
			// invalid: count broken
			{data.Int(5), data.Float(-1.5), data.Float(3.0), data.Float(4.0), nil},
			{data.Int(5), data.Float(-1.5), data.Float(3.0), data.Int(0), nil},
		}},
	}

	for _, testCase := range udf4aryTestCases {
		f := testCase.f
		allInputs := append(testCase.inputs, invalidInputs...)

		Convey(fmt.Sprintf("Given the %s function", testCase.name), t, func() {
			for _, tc := range allInputs {
				tc := tc

				Convey(fmt.Sprintf("When evaluating it on %#v",
					[]data.Value{tc.input1, tc.input2, tc.input3, tc.input4}), func() {
					val, err := f.Call(nil, tc.input1, tc.input2, tc.input3, tc.input4)

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
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 4)
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

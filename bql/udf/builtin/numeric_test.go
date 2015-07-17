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

type udfTestCase struct {
	name   string
	f      udf.UDF
	inputs []udfTestCaseInput
}

type udfTestCaseInput struct {
	input    data.Value
	expected data.Value
}

func TestNumericFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfTestCaseInput{
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

	udfTestCases := []udfTestCase{
		{"abs", absFunc, []udfTestCaseInput{
			{data.Int(7), data.Int(7)},
			{data.Int(-7), data.Int(7)},
			{data.Float(2.3), data.Float(2.3)},
			{data.Float(-2.3), data.Float(2.3)},
		}},
		{"cbrt", cbrtFunc, []udfTestCaseInput{
			{data.Int(27), data.Float(3.0)},
			{data.Int(-27), data.Float(-3.0)},
			{data.Float(27.0), data.Float(3.0)},
			{data.Float(-27.0), data.Float(-3.0)},
		}},
		{"ceil", ceilFunc, []udfTestCaseInput{
			{data.Int(27), data.Int(27)},
			{data.Int(-27), data.Int(-27)},
			{data.Float(42.8), data.Float(43.0)},
			{data.Float(-42.8), data.Float(-42.0)},
		}},
		{"degrees", degreesFunc, []udfTestCaseInput{
			{data.Int(1), data.Float(57.29577951308232)},
			{data.Float(0.5), data.Float(28.6478897565412)},
		}},
		{"exp", expFunc, []udfTestCaseInput{
			{data.Int(1), data.Float(math.E)},
			{data.Float(math.Log(1.0)), data.Float(1.0)},
		}},
		{"floor", floorFunc, []udfTestCaseInput{
			{data.Int(27), data.Int(27)},
			{data.Int(-27), data.Int(-27)},
			{data.Float(42.8), data.Float(42.0)},
			{data.Float(-42.8), data.Float(-43.0)},
		}},
	}

	for _, udfTestCase := range udfTestCases {
		udfTestCase := udfTestCase
		f := udfTestCase.f
		testCases := append(udfTestCase.inputs, invalidInputs...)

		Convey(fmt.Sprintf("Given the %s function", udfTestCase.name), t, func() {
			for _, tc := range testCases {
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
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(udfTestCase.name, 1)
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

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
	}

	for _, testCase := range udfUnaryTestCases {
		f := testCase.f
		allInputs := append(testCase.inputs, invalidInputs...)

		Convey(fmt.Sprintf("Given the %s function", testCase.name), t, func() {
			Convey("Then it should be an aggregate in the first parameter", func() {
				So(f.IsAggregationParameter(0), ShouldBeTrue)
			})

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

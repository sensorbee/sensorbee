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

func TestUnaryStringFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfUnaryTestCaseInput{
		// NULL input -> NULL output
		{data.Null{}, data.Null{}},
		// cannot process the following
		{data.Array{}, nil},
		{data.Blob{}, nil},
		{data.Bool(true), nil},
		{data.Float(2.3), nil},
		{data.Int(7), nil},
		{data.Map{}, nil},
		{data.Timestamp(someTime), nil},
	}

	udfUnaryTestCases := []udfUnaryTestCase{
		{"bit_length", bitLengthFunc, []udfUnaryTestCaseInput{
			{data.String(""), data.Int(0)},
			{data.String("jose"), data.Int(32)},
			{data.String("日本語"), data.Int(72)},
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

package builtin

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"math"
	"testing"
	"time"
)

func TestUnaryArrayFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfUnaryTestCaseInput{
		// NULL input -> NULL output
		{data.Null{}, data.Null{}},
		// cannot process the following
		{data.Blob{}, nil},
		{data.Bool(true), nil},
		{data.Float(1.0), nil},
		{data.Int(1), nil},
		{data.Map{}, nil},
		{data.String("hoge"), nil},
		{data.Timestamp(someTime), nil},
	}

	udfUnaryTestCases := []udfUnaryTestCase{
		{"array_length", arrayLengthFunc, []udfUnaryTestCaseInput{
			{data.Array{}, data.Int(0)},
			{data.Array{data.Int(2)}, data.Int(1)},
			{data.Array{data.Null{}}, data.Int(1)},
			{data.Array{data.Int(2), data.Float(3)}, data.Int(2)},
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
								fActual, _ := data.AsFloat(val)
								fExpected, _ := data.AsFloat(tc.expected)
								if math.IsNaN(fExpected) {
									So(math.IsNaN(fActual), ShouldBeTrue)
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

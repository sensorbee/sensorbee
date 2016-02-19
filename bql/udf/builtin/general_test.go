package builtin

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"math"
	"testing"
)

func TestVariadicGeneralFuncs(t *testing.T) {
	udfVariadicTestCases := []udfVariadicTestCase{
		{"coalesce", coalesceFunc, []udfVariadicTestCaseInput{
			{[]data.Value{data.Null{}}, data.Null{}},
			{[]data.Value{data.String("a"), data.Null{}, data.String("b")},
				data.String("a")},
			{[]data.Value{data.Null{}, data.Int(7), data.String("b")},
				data.Int(7)},
		}},
	}

	for _, testCase := range udfVariadicTestCases {
		f := testCase.f
		allInputs := testCase.inputs

		Convey(fmt.Sprintf("Given the %s function", testCase.name), t, func() {
			for _, tc := range allInputs {
				tc := tc

				Convey(fmt.Sprintf("When evaluating it on %#v", tc.input), func() {
					val, err := f.Call(nil, tc.input...)

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
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 4)
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

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
		{"char_length", charLengthFunc, []udfUnaryTestCaseInput{
			{data.String(""), data.Int(0)},
			{data.String("jose"), data.Int(4)},
			{data.String("日本語"), data.Int(3)},
		}},
		{"lower", lowerFunc, []udfUnaryTestCaseInput{
			{data.String(""), data.String("")},
			{data.String("JosÉ"), data.String("josé")},
			{data.String("日本語"), data.String("日本語")},
		}},
		{"upper", upperFunc, []udfUnaryTestCaseInput{
			{data.String(""), data.String("")},
			{data.String("José"), data.String("JOSÉ")},
			{data.String("日本語"), data.String("日本語")},
		}},
		{"octet_length", octetLengthFunc, []udfUnaryTestCaseInput{
			{data.String(""), data.Int(0)},
			{data.String("jose"), data.Int(4)},
			{data.String("日本語"), data.Int(9)},
		}},
		{"ltrim", ltrimSpaceFunc, []udfUnaryTestCaseInput{
			{data.String("  trim"), data.String("trim")},
			{data.String(" \n trim "), data.String("trim ")},
		}},
		{"rtrim", rtrimSpaceFunc, []udfUnaryTestCaseInput{
			{data.String("trim  "), data.String("trim")},
			{data.String(" trim \n "), data.String(" trim")},
		}},
		{"btrim", btrimSpaceFunc, []udfUnaryTestCaseInput{
			{data.String(" \t trim \n "), data.String("trim")},
		}},
		{"md5", md5Func, []udfUnaryTestCaseInput{
			{data.String("abc"), data.String("900150983cd24fb0d6963f7d28e17f72")},
			{data.String("日本語\n"), data.String("2123035863e00ab6633d0f429fd9aefa")},
		}},
		{"sha1", sha1Func, []udfUnaryTestCaseInput{
			{data.String("abc\n"), data.String("03cfd743661f07975fa2f1220c5194cbaff48451")},
			{data.String("日本語\n"), data.String("eb2ac8ca4ec0f3913058ccd239dd984bda31a821")},
		}},
		{"sha256", sha256Func, []udfUnaryTestCaseInput{
			{data.String("abc\n"), data.String("edeaaff3f1774ad2888673770c6d64097e391bc362d7d6fb34982ddf0efd18cb")},
			{data.String("日本語\n"), data.String("a43d56ae90ff2daebd847bf06f9c0a7b416f48f89b6e9dbefe7886540c94b550")},
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
				if dispatcher, ok := regFun.(*arityDispatcher); ok {
					regFun = dispatcher.unary
				}
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

func TestBinaryStringFuncs(t *testing.T) {
	someTime := time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)

	invalidInputs := []udfBinaryTestCaseInput{
		// NULL input -> NULL output
		{data.Null{}, data.String("1"), data.Null{}},
		{data.String("1"), data.Null{}, data.Null{}},
		// cannot process the following
		{data.Array{}, data.Array{}, nil},
		{data.Blob{}, data.Blob{}, nil},
		{data.Bool(true), data.Bool(true), nil},
		{data.Float(2.7), data.Float(2.7), nil},
		{data.Int(3), data.Int(3), nil},
		{data.Map{}, data.Map{}, nil},
		{data.Timestamp(someTime), data.Timestamp(someTime), nil},
	}

	udfBinaryTestCases := []udfBinaryTestCase{
		{"strpos", strposFunc, []udfBinaryTestCaseInput{
			{data.String("high"), data.String("ig"), data.Int(2)},
			{data.String(""), data.String("ig"), data.Int(0)},
			{data.String("high"), data.String(""), data.Int(1)},
			{data.String("high"), data.String("x"), data.Int(0)},
		}},
		{"substring", substringFunc, []udfBinaryTestCaseInput{
			// substring(string, regex)
			{data.String("high"), data.String("ig"), data.String("ig")},
			{data.String("Thomas"), data.String("...$"), data.String("mas")},
			{data.String("Thomas"), data.String("h.+s"), data.String("homas")},
			{data.String("Thomas"), data.String(""), data.String("")},
			{data.String("Thomas"), data.String("?"), nil},
			// substring(string, fromIdx)
			{data.String("Thomas"), data.Int(1), data.String("Thomas")},
			{data.String("Thomas"), data.Int(2), data.String("homas")},
			{data.String("Thomas"), data.Int(0), nil},
			{data.String("Thomas"), data.Int(6), data.String("s")},
			{data.String("Thomas"), data.Int(7), data.String("")},
			{data.String("Thomas"), data.Int(8), data.String("")},
			{data.String("日本語"), data.Int(3), data.String("語")},
		}},
		{"ltrim", ltrimFunc, []udfBinaryTestCaseInput{
			{data.String("zzzytrim"), data.String("xyz"), data.String("trim")},
			{data.String("zzzytrimz"), data.String("xyz"), data.String("trimz")},
		}},
		{"rtrim", rtrimFunc, []udfBinaryTestCaseInput{
			{data.String("trimzzzy"), data.String("xyz"), data.String("trim")},
			{data.String("zzzytrimz"), data.String("xyz"), data.String("zzzytrim")},
		}},
		{"btrim", btrimFunc, []udfBinaryTestCaseInput{
			{data.String("zzzytrimz"), data.String("xyz"), data.String("trim")},
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

type udf3aryTestCase struct {
	name   string
	f      udf.UDF
	inputs []udf3aryTestCaseInput
}

type udf3aryTestCaseInput struct {
	input1   data.Value
	input2   data.Value
	input3   data.Value
	expected data.Value
}

func Test3aryStringFuncs(t *testing.T) {
	invalidInputs := []udf3aryTestCaseInput{
		// NULL input -> NULL output
		{data.String("1"), data.Null{}, data.Null{}, data.Null{}},
		{data.Null{}, data.String("1"), data.Null{}, data.Null{}},
		{data.Null{}, data.Null{}, data.Int(1), data.Null{}},
		{data.Null{}, data.Null{}, data.Null{}, data.Null{}},
	}

	udf3aryTestCases := []udf3aryTestCase{
		{"overlay", overlayFunc, []udf3aryTestCaseInput{
			{data.String("Txxxxas"), data.String("hom"), data.Int(2), data.String("Thomxas")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(1), data.String("homxxas")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(7), data.String("Txxxxahom")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(8), data.String("Txxxxashom")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(9), data.String("Txxxxashom")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(100), data.String("Txxxxashom")},
			{data.String("日本語"), data.String("中国"), data.Int(1), data.String("中国語")},
			{data.String("Txxxxas"), data.String(""), data.Int(2), data.String("Txxxxas")},
			{data.String(""), data.String("hom"), data.Int(2), data.String("hom")},
			// invalid cases
			{data.String("Txxxxas"), data.String("hom"), data.Int(0), nil},
			{data.Int(3), data.String("hom"), data.Int(1), nil},
			{data.String("Txxxxas"), data.Int(4), data.Int(1), nil},
		}},
		{"substring", substringFunc, []udf3aryTestCaseInput{
			// substring(string, fromIdx, length)
			{data.String("Thomas"), data.Int(1), data.Int(2), data.String("Th")},
			{data.String("Thomas"), data.Int(2), data.Int(3), data.String("hom")},
			{data.String("Thomas"), data.Int(0), data.Int(2), nil},
			{data.String("Thomas"), data.Int(1), data.Int(-1), nil},
			{data.String("Thomas"), data.Int(6), data.Int(0), data.String("")},
			{data.String("Thomas"), data.Int(5), data.Int(2), data.String("as")},
			{data.String("Thomas"), data.Int(6), data.Int(1), data.String("s")},
			{data.String("Thomas"), data.Int(6), data.Int(2), data.String("s")},
			{data.String("Thomas"), data.Int(7), data.Int(30), data.String("")},
			{data.String("日本語"), data.Int(1), data.Int(2), data.String("日本")},
		}},
	}

	for _, testCase := range udf3aryTestCases {
		f := testCase.f
		allInputs := append(testCase.inputs, invalidInputs...)

		Convey(fmt.Sprintf("Given the %s function", testCase.name), t, func() {
			for _, tc := range allInputs {
				tc := tc

				Convey(fmt.Sprintf("When evaluating it on %#v",
					[]data.Value{tc.input1, tc.input2, tc.input3}), func() {
					val, err := f.Call(nil, tc.input1, tc.input2, tc.input3)

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
				regFun, err := udf.CopyGlobalUDFRegistry(nil).Lookup(testCase.name, 3)
				if dispatcher, ok := regFun.(*arityDispatcher); ok {
					regFun = dispatcher.ternary
				}
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

func Test4aryStringFuncs(t *testing.T) {
	invalidInputs := []udf4aryTestCaseInput{
		// NULL input -> NULL output
		{data.String("1"), data.Null{}, data.Null{}, data.Null{}, data.Null{}},
		{data.Null{}, data.String("1"), data.Null{}, data.Null{}, data.Null{}},
		{data.Null{}, data.Null{}, data.Int(1), data.Null{}, data.Null{}},
		{data.Null{}, data.Null{}, data.Null{}, data.Int(1), data.Null{}},
	}

	udf4aryTestCases := []udf4aryTestCase{
		{"overlay", overlayFunc, []udf4aryTestCaseInput{
			// these are the same cases as in the ternary case:
			{data.String("Txxxxas"), data.String("hom"), data.Int(2), data.Int(3), data.String("Thomxas")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(1), data.Int(3), data.String("homxxas")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(7), data.Int(3), data.String("Txxxxahom")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(8), data.Int(3), data.String("Txxxxashom")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(9), data.Int(3), data.String("Txxxxashom")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(100), data.Int(3), data.String("Txxxxashom")},
			{data.String("日本語"), data.String("中国"), data.Int(1), data.Int(2), data.String("中国語")},
			{data.String("Txxxxas"), data.String(""), data.Int(2), data.Int(0), data.String("Txxxxas")},
			{data.String(""), data.String("hom"), data.Int(2), data.Int(3), data.String("hom")},
			// these are variations:
			{data.String("Txxxxas"), data.String("hom"), data.Int(2), data.Int(4), data.String("Thomas")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(2), data.Int(6), data.String("Thom")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(1), data.Int(1), data.String("homxxxxas")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(7), data.Int(0), data.String("Txxxxahoms")},
			{data.String("Txxxxas"), data.String("hom"), data.Int(100), data.Int(1), data.String("Txxxxashom")},
			{data.String("日本語"), data.String("ドイツ"), data.Int(1), data.Int(2), data.String("ドイツ語")},
			{data.String("Txxxxas"), data.String(""), data.Int(2), data.Int(4), data.String("Tas")},
			{data.String(""), data.String("hom"), data.Int(2), data.Int(1), data.String("hom")},
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
				if dispatcher, ok := regFun.(*arityDispatcher); ok {
					regFun = dispatcher.quaternary
				}
				So(err, ShouldBeNil)
				So(regFun, ShouldHaveSameTypeAs, f)
			})
		})
	}
}

type udfVariadicTestCase struct {
	name   string
	f      udf.UDF
	inputs []udfVariadicTestCaseInput
}

type udfVariadicTestCaseInput struct {
	input    []data.Value
	expected data.Value
}

func TestVariadicStringFuncs(t *testing.T) {
	udfVariadicTestCases := []udfVariadicTestCase{
		{"concat", concatFunc, []udfVariadicTestCaseInput{
			{[]data.Value{data.Null{}}, data.String("")},
			{[]data.Value{data.String("a")}, data.String("a")},
			{[]data.Value{data.String("a"), data.String("b")},
				data.String("ab")},
			{[]data.Value{data.String("a"), data.Null{}, data.String("b")},
				data.String("ab")},
			// invalid cases
			{[]data.Value{data.String("a"), data.Int(2)}, nil},
		}},
		{"concat_ws", concatWsFunc, []udfVariadicTestCaseInput{
			{[]data.Value{data.String(","), data.Null{}}, data.String("")},
			{[]data.Value{data.String(","), data.String("a")}, data.String("a")},
			{[]data.Value{data.String(","), data.String("a"), data.String("b")},
				data.String("a,b")},
			{[]data.Value{data.String(","), data.String("a"), data.Null{}, data.String("b")},
				data.String("a,b")},
			{[]data.Value{data.String(","), data.Null{}, data.String("a"), data.String("b")},
				data.String("a,b")},
			{[]data.Value{data.Null{}, data.String("a"), data.String("b")},
				data.Null{}},
			// invalid cases
			{[]data.Value{data.String(",")}, nil},
			{[]data.Value{data.String(","), data.String("a"), data.Int(2)}, nil},
			{[]data.Value{data.Int(2), data.String("a")}, nil},
		}},
		{"format", formatFunc, []udfVariadicTestCaseInput{
			{[]data.Value{data.Null{}, data.String("b")},
				data.Null{}},
			{[]data.Value{data.String("Hello %s, %d"), data.String("b"), data.Int(17)},
				data.String("Hello b, 17")},
			{[]data.Value{data.String("%f, %s, %.2f"), data.Float(math.Inf(1)), data.Null{}, data.Float(2.3)},
				data.String("+Inf, null, 2.30")},
			{[]data.Value{data.String("%s %% %s"), data.Map{"a": data.Int(3)}, data.Array{data.Int(2)}},
				data.String("{\"a\":3} % [2]")},
			// number of parameters does not match
			{[]data.Value{data.String("a"), data.String("b")},
				data.String("a%!(EXTRA string=b)")},
			{[]data.Value{data.String("%s a %d"), data.String("b")},
				data.String("b a %!d(MISSING)")},
			// type does not match
			{[]data.Value{data.String("%f"), data.Int(6)},
				data.String("%!f(data.Int=6)")},
			// invalid cases
			{[]data.Value{data.String("a")}, nil},
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

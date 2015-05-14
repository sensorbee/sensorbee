package tuple

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
	"time"
)

type convTestInput struct {
	ValueDesc string
	Value     Value
	Expected  interface{}
}

func TestToBool(t *testing.T) {
	testCases := map[string]([]convTestInput){
		"Null": []convTestInput{
			convTestInput{"Null", Null{}, false},
		},
		"Bool": []convTestInput{
			convTestInput{"true", Bool(true), true},
			convTestInput{"false", Bool(false), false},
		},
		"Int": []convTestInput{
			convTestInput{"positive", Int(2), true},
			convTestInput{"negative", Int(-2), true},
			convTestInput{"zero", Int(0), false},
		},
		"Float": []convTestInput{
			convTestInput{"positive", Float(3.14), true},
			convTestInput{"negative", Float(-3.14), true},
			convTestInput{"zero", Float(0.0), false},
			convTestInput{"NaN", Float(math.NaN()), false},
		},
		"String": []convTestInput{
			convTestInput{"empty", String(""), false},
			convTestInput{"non-empty", String("hoge"), true},
		},
		"Blob": []convTestInput{
			convTestInput{"empty", Blob(""), false},
			convTestInput{"non-empty", Blob("hoge"), true},
		},
		"Timestamp": []convTestInput{
			convTestInput{"zero", Timestamp(time.Time{}), false},
			convTestInput{"now", Timestamp(time.Now()), true},
		},
		"Array": []convTestInput{
			convTestInput{"empty", Array{}, false},
			convTestInput{"non-empty", Array{Int(2), String("foo")}, true},
		},
		"Map": []convTestInput{
			convTestInput{"empty", Map{}, false},
			convTestInput{"non-empty", Map{"a": Int(2), "b": String("foo")}, true},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToBool(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToBool", testCases)
}

func TestToInt(t *testing.T) {
	now := time.Now()
	negTime := time.Unix(-1, -1000)

	testCases := map[string]([]convTestInput){
		"Null": []convTestInput{
			convTestInput{"Null", Null{}, nil},
		},
		"Bool": []convTestInput{
			convTestInput{"true", Bool(true), int64(1)},
			convTestInput{"false", Bool(false), int64(0)},
		},
		"Int": []convTestInput{
			convTestInput{"positive", Int(2), int64(2)},
			convTestInput{"negative", Int(-2), int64(-2)},
			convTestInput{"zero", Int(0), int64(0)},
		},
		"Float": []convTestInput{
			// normal conversion
			convTestInput{"positive", Float(3.14), int64(3)},
			convTestInput{"negative", Float(-3.14), int64(-3)},
			convTestInput{"zero", Float(0.0), int64(0)},
			// we truncate and don't round
			convTestInput{"positive (> x.5)", Float(2.71), int64(2)},
			convTestInput{"negative (< x.5)", Float(-2.71), int64(-2)},
			// we cannot convert all numbers
			convTestInput{"maximal positive", Float(math.MaxFloat64), nil},
			convTestInput{"maximal negative", Float(-math.MaxFloat64), nil},
		},
		"String": []convTestInput{
			convTestInput{"empty", String(""), nil},
			convTestInput{"non-empty", String("hoge"), nil},
			convTestInput{"numeric", String("123456"), int64(123456)},
		},
		"Blob": []convTestInput{
			convTestInput{"empty", Blob(""), nil},
			convTestInput{"non-empty", Blob("hoge"), nil},
		},
		"Timestamp": []convTestInput{
			// The zero value for a time.Time is *not* the timestamp
			// that has unix time zero!
			convTestInput{"zero", Timestamp(time.Time{}), int64(-62135596800000000)},
			convTestInput{"now", Timestamp(now), now.UnixNano() / 1000},
			convTestInput{"negative", Timestamp(negTime), negTime.UnixNano() / 1000},
		},
		"Array": []convTestInput{
			convTestInput{"empty", Array{}, nil},
			convTestInput{"non-empty", Array{Int(2), String("foo")}, nil},
		},
		"Map": []convTestInput{
			convTestInput{"empty", Map{}, nil},
			convTestInput{"non-empty", Map{"a": Int(2), "b": String("foo")}, nil},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToInt(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToInt", testCases)
}

func TestToFloat(t *testing.T) {
	now := time.Now()
	// We would like to check whether ToFloat(timestamp) equals the
	// value of `float64(timestamp.UnixNano()) / 1e9` but actually
	// this occasionally fails due to numerical rounding. Therefore
	// we use a computation method closer to the actual computation,
	// method, but still based on UnixNano().
	nowSeconds := now.UnixNano() / 1e9
	nowNanoSeconds := now.UnixNano() % 1e9
	nowFloatSeconds := float64(nowSeconds) + float64(nowNanoSeconds)/1e9

	testCases := map[string]([]convTestInput){
		"Null": []convTestInput{
			convTestInput{"Null", Null{}, nil},
		},
		"Bool": []convTestInput{
			convTestInput{"true", Bool(true), float64(1.0)},
			convTestInput{"false", Bool(false), float64(0.0)},
		},
		"Int": []convTestInput{
			// normal conversion
			convTestInput{"positive", Int(2), float64(2.0)},
			convTestInput{"negative", Int(-2), float64(-2.0)},
			convTestInput{"zero", Int(0), float64(0.0)},
			// float64 can represent all int64 numbers (with loss
			// of precision, though)
			convTestInput{"maximal positive", Int(math.MaxInt64), float64(9.223372036854776e+18)},
			convTestInput{"maximal negative", Int(math.MinInt64), float64(-9.223372036854776e+18)},
		},
		"Float": []convTestInput{
			convTestInput{"positive", Float(3.14), float64(3.14)},
			convTestInput{"negative", Float(-3.14), float64(-3.14)},
			convTestInput{"zero", Float(0.0), float64(0)},
		},
		"String": []convTestInput{
			convTestInput{"empty", String(""), nil},
			convTestInput{"non-empty", String("hoge"), nil},
			convTestInput{"numeric", String("123.456"), float64(123.456)},
		},
		"Blob": []convTestInput{
			convTestInput{"empty", Blob(""), nil},
			convTestInput{"non-empty", Blob("hoge"), nil},
		},
		"Timestamp": []convTestInput{
			// The zero value for a time.Time is *not* the timestamp
			// that has unix time zero!
			convTestInput{"zero", Timestamp(time.Time{}), float64(-62135596800)},
			convTestInput{"now", Timestamp(now), nowFloatSeconds},
		},
		"Array": []convTestInput{
			convTestInput{"empty", Array{}, nil},
			convTestInput{"non-empty", Array{Int(2), String("foo")}, nil},
		},
		"Map": []convTestInput{
			convTestInput{"empty", Map{}, nil},
			convTestInput{"non-empty", Map{"a": Int(2), "b": String("foo")}, nil},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToFloat(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToFloat", testCases)
}

func TestToString(t *testing.T) {
	now := time.Now()

	testCases := map[string]([]convTestInput){
		"Null": []convTestInput{
			convTestInput{"Null", Null{}, "null"},
		},
		"Bool": []convTestInput{
			convTestInput{"true", Bool(true), "true"},
			convTestInput{"false", Bool(false), "false"},
		},
		"Int": []convTestInput{
			convTestInput{"positive", Int(2), "2"},
			convTestInput{"negative", Int(-2), "-2"},
			convTestInput{"zero", Int(0), "0"},
		},
		"Float": []convTestInput{
			convTestInput{"positive", Float(3.14), "3.14"},
			convTestInput{"negative", Float(-3.14), "-3.14"},
			convTestInput{"zero", Float(0.0), "0"},
		},
		"String": []convTestInput{
			convTestInput{"empty", String(""), ""},
			convTestInput{"non-empty", String("hoge"), "hoge"},
		},
		"Blob": []convTestInput{
			convTestInput{"empty", Blob(""), "tuple.Blob{}"},
			convTestInput{"non-empty", Blob("hoge"), "tuple.Blob{0x68, 0x6f, 0x67, 0x65}"},
		},
		"Timestamp": []convTestInput{
			convTestInput{"zero", Timestamp(time.Time{}), "0001-01-01T00:00:00Z"},
			convTestInput{"now", Timestamp(now), now.Format(time.RFC3339Nano)},
		},
		"Array": []convTestInput{
			convTestInput{"empty", Array{}, "tuple.Array{}"},
			convTestInput{"non-empty", Array{Int(2), String("foo")}, `tuple.Array{2, "foo"}`},
		},
		"Map": []convTestInput{
			convTestInput{"empty", Map{}, `tuple.Map{}`},
			// the following test would fail once in a while because
			// golang randomizes the keys for Maps, i.e., we cannot be sure
			// that we will always get the same string
			//convTestInput{"non-empty", Map{"a": Int(2), "b": String("foo")}, `tuple.Map{"a":2, "b":"foo"}`},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToString(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToString", testCases)
}

func TestToTime(t *testing.T) {
	now := time.Now()
	jst, err := time.LoadLocation("JST")
	if err != nil {
		t.Fatal("Cannot load JST location:", err)
	}

	testCases := map[string]([]convTestInput){
		"Null": []convTestInput{
			convTestInput{"Null", Null{}, time.Time{}},
		},
		"Bool": []convTestInput{
			convTestInput{"true", Bool(true), nil},
			convTestInput{"false", Bool(false), nil},
		},
		"Int": []convTestInput{
			convTestInput{"positive", Int(2), time.Unix(0, 2000)},
			convTestInput{"negative", Int(-2), time.Unix(0, -2000)},
			convTestInput{"large and positive", Int(2e6), time.Unix(2, 0)},
			convTestInput{"large and negative", Int(-2e6), time.Unix(-2, 0)},
			convTestInput{"zero", Int(0), time.Unix(0, 0)},
		},
		"Float": []convTestInput{
			convTestInput{"positive", Float(3.14), time.Unix(3, 14e7)},
			convTestInput{"negative", Float(-3.14), time.Unix(-3, -14e7)},
			convTestInput{"negative (alternative)", Float(-3.14), time.Unix(-4, 86e7)},
			convTestInput{"zero", Float(0.0), time.Unix(0, 0)},
		},
		"String": []convTestInput{
			convTestInput{"empty", String(""), nil},
			convTestInput{"non-empty", String("hoge"), nil},
			convTestInput{"valid time string", String("1970-01-01T09:00:02+09:00"), time.Unix(2, 0).In(jst)},
			convTestInput{"valid time string with ns", String(now.Format(time.RFC3339Nano)), now},
		},
		"Blob": []convTestInput{
			convTestInput{"empty", Blob(""), nil},
			convTestInput{"non-empty", Blob("hoge"), nil},
		},
		"Timestamp": []convTestInput{
			convTestInput{"zero", Timestamp(time.Time{}), time.Time{}},
			convTestInput{"now", Timestamp(now), now},
		},
		"Array": []convTestInput{
			convTestInput{"empty", Array{}, nil},
			convTestInput{"non-empty", Array{Int(2), String("foo")}, nil},
		},
		"Map": []convTestInput{
			convTestInput{"empty", Map{}, nil},
			convTestInput{"non-empty", Map{"a": Int(2), "b": String("foo")}, nil},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToTime(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToTime", testCases)
}

func runConversionTestCases(t *testing.T,
	toFun func(v Value) (interface{}, error),
	funcName string,
	testCases map[string][]convTestInput) {
	for valType, cases := range testCases {
		cases := cases
		Convey(fmt.Sprintf("Given a %s value", valType), t, func() {
			for _, testCase := range cases {
				tc := testCase
				Convey(fmt.Sprintf("When it is %s", tc.ValueDesc), func() {
					inVal := tc.Value
					exp := tc.Expected
					if exp == nil {
						Convey(fmt.Sprintf("Then %s returns an error", funcName), func() {
							_, err := toFun(inVal)
							So(err, ShouldNotBeNil)
						})
					} else {
						Convey(fmt.Sprintf("Then %s returns %v", funcName, exp), func() {
							val, err := toFun(inVal)
							So(err, ShouldBeNil)
							switch fval := val.(type) {
							case float64:
								switch fexp := exp.(type) {
								case float64:
									// equal up to machine precision
									So(fval-fexp, ShouldAlmostEqual, 0.0, 1e-15)
								default:
									So(val, ShouldResemble, exp)
								}
							default:
								So(val, ShouldResemble, exp)
							}
						})
					}
				})
			}
		})
	}
}

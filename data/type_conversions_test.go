package data

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
	testCases := map[string][]convTestInput{
		"Null": {
			{"Null", Null{}, false},
		},
		"Bool": {
			{"true", Bool(true), true},
			{"false", Bool(false), false},
		},
		"Int": {
			{"positive", Int(2), true},
			{"negative", Int(-2), true},
			{"zero", Int(0), false},
		},
		"Float": {
			{"positive", Float(3.14), true},
			{"negative", Float(-3.14), true},
			{"zero", Float(0.0), false},
			{"NaN", Float(math.NaN()), false},
		},
		"String": {
			{"empty", String(""), false},
			{"non-empty", String("hoge"), true},
		},
		"Blob": {
			{"empty", Blob(""), false},
			{"nil", Blob(nil), false},
			{"non-empty", Blob("hoge"), true},
		},
		"Timestamp": {
			{"zero", Timestamp(time.Time{}), false},
			{"now", Timestamp(time.Now()), true},
		},
		"Array": {
			{"empty", Array{}, false},
			{"non-empty", Array{Int(2), String("foo")}, true},
		},
		"Map": {
			{"empty", Map{}, false},
			{"non-empty", Map{"a": Int(2), "b": String("foo")}, true},
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

	testCases := map[string][]convTestInput{
		"Null": {
			{"Null", Null{}, nil},
		},
		"Bool": {
			{"true", Bool(true), int64(1)},
			{"false", Bool(false), int64(0)},
		},
		"Int": {
			{"positive", Int(2), int64(2)},
			{"negative", Int(-2), int64(-2)},
			{"zero", Int(0), int64(0)},
		},
		"Float": {
			// normal conversion
			{"positive", Float(3.14), int64(3)},
			{"negative", Float(-3.14), int64(-3)},
			{"zero", Float(0.0), int64(0)},
			// we truncate and don't round
			{"positive (> x.5)", Float(2.71), int64(2)},
			{"negative (< x.5)", Float(-2.71), int64(-2)},
			// we cannot convert all numbers
			{"maximal positive", Float(math.MaxFloat64), nil},
			{"maximal negative", Float(-math.MaxFloat64), nil},
		},
		"String": {
			{"empty", String(""), nil},
			{"non-empty", String("hoge"), nil},
			{"numeric", String("123456"), int64(123456)},
		},
		"Blob": {
			{"empty", Blob(""), nil},
			{"non-empty", Blob("hoge"), nil},
		},
		"Timestamp": {
			// The zero value for a time.Time is *not* the timestamp
			// that has unix time zero!
			{"zero", Timestamp(time.Time{}), int64(-62135596800000000)},
			{"now", Timestamp(now), now.UnixNano() / 1000},
			{"negative", Timestamp(negTime), negTime.UnixNano() / 1000},
		},
		"Array": {
			{"empty", Array{}, nil},
			{"non-empty", Array{Int(2), String("foo")}, nil},
		},
		"Map": {
			{"empty", Map{}, nil},
			{"non-empty", Map{"a": Int(2), "b": String("foo")}, nil},
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

	testCases := map[string][]convTestInput{
		"Null": {
			{"Null", Null{}, nil},
		},
		"Bool": {
			{"true", Bool(true), float64(1.0)},
			{"false", Bool(false), float64(0.0)},
		},
		"Int": {
			// normal conversion
			{"positive", Int(2), float64(2.0)},
			{"negative", Int(-2), float64(-2.0)},
			{"zero", Int(0), float64(0.0)},
			// float64 can represent all int64 numbers (with loss
			// of precision, though)
			{"maximal positive", Int(math.MaxInt64), float64(9.223372036854776e+18)},
			{"maximal negative", Int(math.MinInt64), float64(-9.223372036854776e+18)},
		},
		"Float": {
			{"positive", Float(3.14), float64(3.14)},
			{"negative", Float(-3.14), float64(-3.14)},
			{"zero", Float(0.0), float64(0)},
		},
		"String": {
			{"empty", String(""), nil},
			{"non-empty", String("hoge"), nil},
			{"numeric", String("123.456"), float64(123.456)},
		},
		"Blob": {
			{"empty", Blob(""), nil},
			{"non-empty", Blob("hoge"), nil},
		},
		"Timestamp": {
			// The zero value for a time.Time is *not* the timestamp
			// that has unix time zero!
			{"zero", Timestamp(time.Time{}), float64(-62135596800)},
			{"now", Timestamp(now), nowFloatSeconds},
		},
		"Array": {
			{"empty", Array{}, nil},
			{"non-empty", Array{Int(2), String("foo")}, nil},
		},
		"Map": {
			{"empty", Map{}, nil},
			{"non-empty", Map{"a": Int(2), "b": String("foo")}, nil},
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

	testCases := map[string][]convTestInput{
		"Null": {
			{"Null", Null{}, "null"},
		},
		"Bool": {
			{"true", Bool(true), "true"},
			{"false", Bool(false), "false"},
		},
		"Int": {
			{"positive", Int(2), "2"},
			{"negative", Int(-2), "-2"},
			{"zero", Int(0), "0"},
		},
		"Float": {
			{"positive", Float(3.14), "3.14"},
			{"negative", Float(-3.14), "-3.14"},
			{"zero", Float(0.0), "0"},
		},
		"String": {
			{"empty", String(""), ""},
			{"non-empty", String("hoge"), "hoge"},
		},
		"Blob": {
			{"empty", Blob(""), ""},
			{"non-empty", Blob("hoge"), "hoge"},
		},
		"Timestamp": {
			{"zero", Timestamp(time.Time{}), "0001-01-01T00:00:00Z"},
			{"now", Timestamp(now), now.Format(time.RFC3339Nano)},
		},
		"Array": {
			{"empty", Array{}, "data.Array{}"},
			{"non-empty", Array{Int(2), String("foo")}, `data.Array{2, "foo"}`},
		},
		"Map": {
			{"empty", Map{}, `data.Map{}`},
			// the following test would fail once in a while because
			// golang randomizes the keys for Maps, i.e., we cannot be sure
			// that we will always get the same string
			// convTestInput{"non-empty", Map{"a": Int(2), "b": String("foo")}, `data.Map{"a":2, "b":"foo"}`},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToString(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToString", testCases)
}

func TestToBlob(t *testing.T) {
	testCases := map[string][]convTestInput{
		"Null": {
			{"Null", Null{}, []byte(nil)},
		},
		"Bool": {
			{"true", Bool(true), nil},
			{"false", Bool(false), nil},
		},
		"Int": {
			{"positive", Int(2), nil},
			{"negative", Int(-2), nil},
			{"zero", Int(0), nil},
		},
		"Float": {
			{"positive", Float(3.14), nil},
			{"negative", Float(-3.14), nil},
			{"zero", Float(0.0), nil},
		},
		"String": {
			{"empty", String(""), []byte{}},
			{"non-empty", String("hoge"), []byte("hoge")},
			{"numeric", String("123.456"), []byte("123.456")},
		},
		"Blob": {
			{"nil", Blob(nil), []byte(nil)},
			{"empty", Blob{}, []byte{}},
			{"non-empty", Blob("hoge"), []byte("hoge")},
		},
		"Timestamp": {
			{"zero", Timestamp(time.Time{}), nil},
			{"now", Timestamp(time.Now()), nil},
		},
		"Array": {
			{"empty", Array{}, nil},
			{"non-empty", Array{Int(2), String("foo")}, nil},
		},
		"Map": {
			{"empty", Map{}, nil},
			{"non-empty", Map{"a": Int(2), "b": String("foo")}, nil},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToBlob(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToBlob", testCases)
}

func TestToTimestamp(t *testing.T) {
	now := time.Now()

	testCases := map[string][]convTestInput{
		"Null": {
			{"Null", Null{}, time.Time{}},
		},
		"Bool": {
			{"true", Bool(true), nil},
			{"false", Bool(false), nil},
		},
		"Int": {
			{"positive", Int(2), time.Unix(0, 2000)},
			{"negative", Int(-2), time.Unix(0, -2000)},
			{"large and positive", Int(2e6), time.Unix(2, 0)},
			{"large and negative", Int(-2e6), time.Unix(-2, 0)},
			{"zero", Int(0), time.Unix(0, 0)},
		},
		"Float": {
			{"positive", Float(3.14), time.Unix(3, 14e7)},
			{"negative", Float(-3.14), time.Unix(-3, -14e7)},
			{"negative (alternative)", Float(-3.14), time.Unix(-4, 86e7)},
			{"zero", Float(0.0), time.Unix(0, 0)},
		},
		"String": {
			{"empty", String(""), nil},
			{"non-empty", String("hoge"), nil},
			{"valid time string", String("1970-01-01T09:00:02+09:00"), time.Unix(2, 0)},
			{"valid time string with ns", String(now.Format(time.RFC3339Nano)), now},
		},
		"Blob": {
			{"empty", Blob(""), nil},
			{"non-empty", Blob("hoge"), nil},
		},
		"Timestamp": {
			{"zero", Timestamp(time.Time{}), time.Time{}},
			{"now", Timestamp(now), now},
		},
		"Array": {
			{"empty", Array{}, nil},
			{"non-empty", Array{Int(2), String("foo")}, nil},
		},
		"Map": {
			{"empty", Map{}, nil},
			{"non-empty", Map{"a": Int(2), "b": String("foo")}, nil},
		},
	}

	toFun := func(v Value) (interface{}, error) {
		val, err := ToTimestamp(v)
		return val, err
	}
	runConversionTestCases(t, toFun, "ToTimestamp", testCases)
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
							case time.Time:
								So(val, ShouldHappenOnOrBetween, exp, exp)
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

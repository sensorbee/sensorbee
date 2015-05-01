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

	testCases := map[string]([]convTestInput){
		"Null": []convTestInput{
			convTestInput{"Null", Null{}, int64(0)},
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
		},
		"Blob": []convTestInput{
			convTestInput{"empty", Blob(""), nil},
			convTestInput{"non-empty", Blob("hoge"), nil},
		},
		"Timestamp": []convTestInput{
			// The zero value for a time.Time is *not* the timestamp
			// that has unix time zero!
			convTestInput{"zero", Timestamp(time.Time{}), int64(-62135596800)},
			convTestInput{"now", Timestamp(now), now.Unix()},
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
							So(val, ShouldEqual, exp)
						})
					}
				})
			}
		})
	}
}

package tuple

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type convTestInput struct {
	ValueDesc string
	Value     Value
	Expected  interface{}
}

func TestToBool(t *testing.T) {
	toFun := ToBool
	funcName := "ToBool"

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

	for valType, cases := range testCases {
		cases := cases
		Convey(fmt.Sprintf("Given a %s value", valType), t, func() {
			for _, testCase := range cases {
				tc := testCase
				Convey(fmt.Sprintf("When it is %s", tc.ValueDesc), func() {
					inVal := tc.Value
					exp := tc.Expected
					Convey(fmt.Sprintf("Then %s returns %v", funcName, exp), func() {
						val, err := toFun(inVal)
						So(err, ShouldBeNil)
						So(val, ShouldEqual, exp)
					})
				})
			}
		})
	}
}

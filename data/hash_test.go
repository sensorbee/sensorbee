package data

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"reflect"
	"testing"
)

var testCases = []struct {
	input    Value
	expected string
}{
	// Blob
	{Blob(""), "B0:"},
	{Blob("hoge"), "B4:hoge"},
	{Blob("日本語"), "B9:日本語"},
	// Bool
	{Bool(false), "b0"},
	{Bool(true), "b1"},
	// Float
	{Float(2.34), "n2.34;"},
	{Float(math.NaN()), "nNaN;"},
	{Float(math.Inf(-1)), "n-Inf;"},
	{Float(2.000000000000000000000000000000000000000000000000000000000001), "n2;"},
	{Float(2.0), "n2;"},
	{Float(0.0), "n0;"},
	{Float(-0.0), "n0;"},
	// Int
	{Int(2), "n2;"}, // same as Float(2.0)
	{Int(-6), "n-6;"},
	// Null
	{Null{}, "N"},
	// String
	{String(""), "s0:"},
	{String("hoge"), "s4:hoge"},
	{String("日本語"), "s9:日本語"},
	// Timestamp
	{Timestamp{}, "t20:0001-01-01T00:00:00Z"},
	// Array
	{Array{}, "a0:"},
	{Array{String("hoge")}, "a1:s4:hoge"},
	{Array{Int(2), Float(3.0)}, "a2:n2;n3;"},
	{Array{Float(2.0), Int(3)}, "a2:n2;n3;"},
	{Array{Int(-6), Array{String("hoge")}}, "a2:n-6;a1:s4:hoge"},
	// Map
	{Map{}, "m0:"},
	{Map{"hoge": String("hoge")}, "m1:s4:hoges4:hoge"},
	{Map{"a": Int(2), "b": Float(3.0)}, "m2:s1:an2;s1:bn3;"},
	{Map{"b": Int(3), "a": Float(2.0)}, "m2:s1:an2;s1:bn3;"},
	{Map{"i": Int(-6), "xy": Map{"h": String("hoge")}}, "m2:s1:in-6;s2:xym1:s1:hs4:hoge"},
}

func TestUpdateHash(t *testing.T) {
	for _, testCase := range testCases {
		testCase := testCase
		Convey(fmt.Sprintf("When serializing %#v", testCase.input), t, func() {
			var b bytes.Buffer
			updateHash(testCase.input, &b)

			Convey(fmt.Sprintf("Then the result should be %s", testCase.expected), func() {
				So(b.String(), ShouldEqual, testCase.expected)
			})
		})
	}
}

func TestEquality(t *testing.T) {
	for i, tc1 := range testCases {
		for j, tc2 := range testCases {
			left := tc1.input
			right := tc2.input
			Convey(fmt.Sprintf("When comparing %#v and %#v", left, right), t, func() {
				de := reflect.DeepEqual(left, right)
				he := HashEqual(left, right)

				Convey("Then the output should be the same", func() {
					if // int vs float
					((i == 8 || i == 9) && (j == 12)) ||
						((j == 8 || j == 9) && (i == 12)) ||
						// array
						((i == 21 && j == 22) || (j == 21 && i == 22)) ||
						// map
						((i == 26 && j == 27) || (j == 26 && i == 27)) ||
						// NaN
						(i == 6 && j == 6) {
						So(de, ShouldBeFalse)
						So(he, ShouldBeTrue)
					} else {
						if de != he {
							fmt.Printf("%v vs %v: %t/%t\n", left, right, de, he)
						}
						So(de, ShouldEqual, he)
					}
				})
			})
		}
	}
}

func BenchmarkDeepEqual(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc1 := range testCases {
			for _, tc2 := range testCases {
				reflect.DeepEqual(tc1.input, tc2.input)
			}
		}
	}
}

func BenchmarkHashEqual(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc1 := range testCases {
			for _, tc2 := range testCases {
				HashEqual(tc1.input, tc2.input)
			}
		}
	}
}

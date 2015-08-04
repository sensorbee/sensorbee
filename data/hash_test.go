package data

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"reflect"
	"testing"
)

var testCases = []struct {
	input Value
}{
	// Blob
	{Blob("")},
	{Blob("hoge")},
	{Blob("日本語")},
	// Bool
	{Bool(false)},
	{Bool(true)},
	// Float
	{Float(2.34)},
	{Float(math.NaN())}, // NaN == NaN is false
	{Float(math.Inf(-1))},
	{Float(2.000000000000000000000000000000000000000000000000000000000001)},
	{Float(2.0)},
	{Float(0.0)},
	{Float(-0.0)},
	// Int
	{Int(2)}, // same as Float(2.0)
	{Int(-6)},
	// Null
	{Null{}},
	// String
	{String("")},
	{String("hoge")},
	{String("日本語")},
	// Timestamp
	{Timestamp{}},
	// Array
	{Array{}},
	{Array{String("hoge")}},
	{Array{Int(2), Float(3.0)}},
	{Array{Float(2.0), Int(3)}},
	{Array{Int(-6), Array{String("hoge")}}},
	// Map
	{Map{}},
	{Map{"hoge": String("hoge")}},
	{Map{"a": Int(2), "b": Float(3.0)}},
	{Map{"b": Int(3), "a": Float(2.0)}},
	{Map{"i": Int(-6), "xy": Map{"h": String("hoge")}}},
}

func TestEquality(t *testing.T) {
	for i, tc1 := range testCases {
		for j, tc2 := range testCases {
			left := tc1.input
			right := tc2.input
			Convey(fmt.Sprintf("When comparing %#v and %#v", left, right), t, func() {
				de := reflect.DeepEqual(left, right)
				he := Equal(left, right)

				Convey("Then the output should be the same", func() {
					if // int vs float
					((i == 8 || i == 9) && (j == 12)) ||
						((j == 8 || j == 9) && (i == 12)) ||
						// array
						((i == 21 && j == 22) || (j == 21 && i == 22)) ||
						// map
						((i == 26 && j == 27) || (j == 26 && i == 27)) {
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

func BenchmarkEqual(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc1 := range testCases {
			for _, tc2 := range testCases {
				Equal(tc1.input, tc2.input)
			}
		}
	}
}

func BenchmarkHash(b *testing.B) {
	var h HashValue
	for n := 0; n < b.N; n++ {
		for _, tc := range testCases {
			h += Hash(tc.input)
		}
	}
}

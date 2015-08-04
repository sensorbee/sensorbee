package data

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"reflect"
	"testing"
	"time"
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
	{Map{"b": Int(3), "a": Float(2.5)}},
	{Map{"i": Int(-6), "xy": Map{"h": String("hoge")}}},
}

func TestHash(t *testing.T) {
	now := time.Now()
	for now.Nanosecond()%1000 == 999 {
		now = time.Now()
	}

	Convey("Given a map containing all types", t, func() {
		m := Map{
			"null":      Null{},
			"true":      True,
			"false":     False,
			"int":       Int(10),
			"float":     Float(2.5),
			"string":    String("hoge"),
			"blob":      Blob("hoge"),
			"timestamp": Timestamp(now),
			"array": Array{Null{}, True, False, Int(10), Float(2.5), String("hoge"),
				Blob("hoge"), Timestamp(now), Array{Null{}, True, False, Int(10), Float(2.5), String("hoge")},
				Map{
					"null":   Null{},
					"true":   True,
					"false":  False,
					"int":    Int(10),
					"float":  Float(2.5),
					"string": String("hoge"),
				},
			},
			"map": Map{
				"null":      Null{},
				"true":      True,
				"false":     False,
				"int":       Int(10),
				"float":     Float(2.5),
				"string":    String("hoge"),
				"blob":      Blob("hoge"),
				"timestamp": Timestamp(now),
				"array": Array{Null{}, True, False, Int(10), Float(2.5), String("hoge"),
					Blob("hoge"), Timestamp(now), Array{Null{}, True, False, Int(10), Float(2.5), String("hoge")},
					Map{
						"null":   Null{},
						"true":   True,
						"false":  False,
						"int":    Int(10),
						"float":  Float(2.5),
						"string": String("hoge"),
					},
				},
			},
		}

		Convey("When a map doesn't contain something invalid", func() {
			Convey("Then Hash should always return the same value", func() {
				So(Hash(m), ShouldEqual, Hash(m))
			})
		})

		Convey("When comparing a float having an integer value to int and vice vasa", func() {
			c := m.Copy()
			c["int"] = Float(10.0)
			m["float"] = Float(3.0)
			c["float"] = Int(3)
			So(c.Set("array[3]", Float(10.0)), ShouldBeNil)
			So(m.Set("map.int", Float(10.0)), ShouldBeNil)

			Convey("Then Hash should return the same value", func() {
				So(Hash(c), ShouldEqual, Hash(m))
			})
		})

		Convey("When comparing timestamps whose differences are less than 1us", func() {
			t := Timestamp(now.Add(1))
			c := m.Copy()
			Convey("Then Hash should return the same value", func() {
				c["timestamp"] = t
				So(Hash(c), ShouldEqual, Hash(m))
			})

			Convey("Then Hash should behave samely when they're in an array", func() {
				So(m.Set("array[7]", t), ShouldBeNil)
				So(Hash(c), ShouldEqual, Hash(m))
			})

			Convey("Then Hash should behave samely when they're in a map", func() {
				So(m.Set("map.timestamp", t), ShouldBeNil)
				So(Hash(m), ShouldEqual, Hash(m))
			})
		})

		Convey("When a map contains NaN", func() {
			Convey("Then Hash should alwasy return different values", func() {
				m["nan"] = Float(math.NaN())
				So(Hash(m), ShouldNotEqual, Hash(m))
			})

			Convey("Then NaN in an array should behave samely", func() {
				So(m.Set("array[0]", Float(math.NaN())), ShouldBeNil)
				So(Hash(m), ShouldNotEqual, Hash(m))
			})

			Convey("Then NaN in a map should behave samely", func() {
				So(m.Set("map.float", Float(math.NaN())), ShouldBeNil)
				So(Hash(m), ShouldNotEqual, Hash(m))
			})
		})
	})
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

func BenchmarkHashLargeMap(b *testing.B) {
	nested := Map{}
	prefixes := []string{"hoge", "moge", "fuga", "pfn", "sensorbee"}
	for i := 0; i < 2500; i++ {
		p := fmt.Sprintf("%v%v", prefixes[i%len(prefixes)], i+1)

		nested[p+"int"] = Int(((-2 * (i % 2)) + 1) * i)
		nested[p+"float"] = Float(math.Exp(math.Pow(float64(i), 0.3)))
		nested[p+"string"] = String("hogehogehogehogehogehogehgoehoge")
		nested[p+"timestamp"] = Timestamp(time.Now().AddDate(0, 0, i))
	}
	m := Map{}
	for i := 0; i < 10; i++ {
		m[fmt.Sprint(i+1)] = nested.Copy()
	}
	b.ResetTimer()

	var h HashValue
	for n := 0; n < b.N; n++ {
		h += Hash(m)
	}
}

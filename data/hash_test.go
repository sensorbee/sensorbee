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
			So(c.Set(MustCompilePath("array[3]"), Float(10.0)), ShouldBeNil)
			So(m.Set(MustCompilePath("map.int"), Float(10.0)), ShouldBeNil)

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
				So(c.Set(MustCompilePath("array[7]"), t), ShouldBeNil)
				So(Hash(c), ShouldEqual, Hash(m))
			})

			Convey("Then Hash should behave samely when they're in a map", func() {
				So(c.Set(MustCompilePath("map.timestamp"), t), ShouldBeNil)
				So(Hash(c), ShouldEqual, Hash(m))
			})
		})

		Convey("When a map contains NaN", func() {
			Convey("Then Hash should alwasy return different values", func() {
				m["nan"] = Float(math.NaN())
				So(Hash(m), ShouldNotEqual, Hash(m))
			})

			Convey("Then NaN in an array should behave samely", func() {
				So(m.Set(MustCompilePath("array[0]"), Float(math.NaN())), ShouldBeNil)
				So(Hash(m), ShouldNotEqual, Hash(m))
			})

			Convey("Then NaN in a map should behave samely", func() {
				So(m.Set(MustCompilePath("map.float"), Float(math.NaN())), ShouldBeNil)
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

				Convey("Then the results of Equal and Less should be consistent", func() {
					if he {
						So(Less(left, right) || Less(right, left), ShouldBeFalse)
					} else {
						if i == 6 { // NaN
							if right.Type() == TypeInt || right.Type() == TypeFloat {
								So(Less(left, right), ShouldBeFalse)
								So(Less(right, left), ShouldBeFalse)
							}
						} else if j == 6 { // NaN
							if left.Type() == TypeInt || left.Type() == TypeFloat {
								So(Less(left, right), ShouldBeFalse)
								So(Less(right, left), ShouldBeFalse)
							}
						} else {
							// if the values are not equal, exactly one of them should be less
							So(Less(left, right), ShouldNotEqual, Less(right, left))
						}
					}
				})
			})
		}
	}
}

func TestLess(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Second)

	ltTestCases := []struct {
		l    Value
		r    Value
		less bool
	}{
		{Null{}, Null{}, false}, // equal
		{Null{}, Bool(true), true},
		{Null{}, Int(1), true},
		{Null{}, Float(3.14), true},
		{Null{}, String("hoge"), true},
		{Null{}, Blob("hello"), true},
		{Null{}, Timestamp(now), true},
		{Null{}, Array{Int(1)}, true},
		{Null{}, Map{"a": Int(1)}, true},
		{Bool(true), Null{}, false},
		{Bool(false), Bool(false), false}, // equal
		{Bool(true), Bool(false), false},
		{Bool(false), Bool(true), true},
		{Bool(true), Bool(true), false},
		{Bool(true), Int(1), true},
		{Bool(true), Float(3.14), true},
		{Bool(true), String("hoge"), true},
		{Bool(true), Blob("hello"), true},
		{Bool(true), Timestamp(now), true},
		{Bool(true), Array{Int(1)}, true},
		{Bool(true), Map{"a": Int(1)}, true},
		{Int(1), Null{}, false},
		{Int(1), Bool(true), false},
		{Int(1), Int(1), false}, // equal
		{Int(1), Int(2), true},
		{Int(2), Int(1), false},
		{Int(1), Float(1.0), false},
		{Int(1), Float(-3.14), false},
		{Int(1), Float(3.14), true},
		{Int(1), String("hoge"), true},
		{Int(1), Blob("hello"), true},
		{Int(1), Timestamp(now), true},
		{Int(1), Array{Int(1)}, true},
		{Int(1), Map{"a": Int(1)}, true},
		{Float(3.14), Null{}, false},
		{Float(3.14), Bool(true), false},
		{Float(3.14), Int(1), false},
		{Float(3.14), Int(4), true},
		{Float(3.14), Float(3.14), false}, // equal
		{Float(3.14), Float(2.14), false},
		{Float(3.14), Float(4.14), true},
		{Float(3.14), String("hoge"), true},
		{Float(3.14), Blob("hello"), true},
		{Float(3.14), Timestamp(now), true},
		{Float(3.14), Array{Int(1)}, true},
		{Float(3.14), Map{"a": Int(1)}, true},
		{String("hoge"), Null{}, false},
		{String("hoge"), Bool(true), false},
		{String("hoge"), Int(1), false},
		{String("hoge"), Float(3.14), false},
		{String("hoge"), String("hoge"), false}, // equal
		{String("hoge"), String("a"), false},
		{String("hoge"), String("k"), true},
		{String("hoge"), Blob("hello"), true},
		{String("hoge"), Timestamp(now), true},
		{String("hoge"), Array{Int(1)}, true},
		{String("hoge"), Map{"a": Int(1)}, true},
		{Blob("hello"), Null{}, false},
		{Blob("hello"), Bool(true), false},
		{Blob("hello"), Int(1), false},
		{Blob("hello"), Float(3.14), false},
		{Blob("hello"), String("hoge"), false},
		{Blob("hello"), Blob("hello"), false}, // equal
		{Blob("hello"), Blob("a"), false},     // more bytes
		{Blob("hello"), Blob("abcdef"), true}, // less bytes
		{Blob("hello"), Blob("hallo"), false}, // hash is smaller
		{Blob("hello"), Blob("abcde"), true},  // hash is larger
		{Blob("hello"), Timestamp(now), true},
		{Blob("hello"), Array{Int(1)}, true},
		{Blob("hello"), Map{"a": Int(1)}, true},
		{Timestamp(now), Null{}, false},
		{Timestamp(now), Bool(true), false},
		{Timestamp(now), Int(1), false},
		{Timestamp(now), Float(3.14), false},
		{Timestamp(now), String("hoge"), false},
		{Timestamp(now), Blob("hello"), false},
		{Timestamp(now), Timestamp(now), false},
		{Timestamp(now), Timestamp(later), true},
		{Timestamp(later), Timestamp(now), false},
		{Timestamp(now), Array{Int(1)}, true},
		{Timestamp(now), Map{"a": Int(1)}, true},
		{Array{Int(1)}, Null{}, false},
		{Array{Int(1)}, Bool(true), false},
		{Array{Int(1)}, Int(1), false},
		{Array{Int(1)}, Float(3.14), false},
		{Array{Int(1)}, String("hoge"), false},
		{Array{Int(1)}, Blob("hello"), false},
		{Array{Int(1)}, Timestamp(now), false},
		{Array{Int(1)}, Array{Int(1)}, false},         // equal
		{Array{Int(1)}, Array{Int(1), Int(2)}, true},  // less entries
		{Array{Int(1), Int(2)}, Array{Int(1)}, false}, // more entries
		{Array{Int(1)}, Array{Int(6)}, true},          // hash is smaller
		{Array{Int(1)}, Array{Int(2)}, false},         // hash is larger
		{Array{Int(1)}, Map{"a": Int(1)}, true},
		{Map{"a": Int(1)}, Null{}, false},
		{Map{"a": Int(1)}, Bool(true), false},
		{Map{"a": Int(1)}, Int(1), false},
		{Map{"a": Int(1)}, Float(3.14), false},
		{Map{"a": Int(1)}, String("hoge"), false},
		{Map{"a": Int(1)}, Blob("hello"), false},
		{Map{"a": Int(1)}, Timestamp(now), false},
		{Map{"a": Int(1)}, Array{Int(1)}, false},
		{Map{"a": Int(1)}, Map{"a": Int(1)}, false},              // equal
		{Map{"a": Int(1)}, Map{"a": Int(1), "b": Int(2)}, true},  // less entries
		{Map{"a": Int(1), "b": Int(2)}, Map{"a": Int(1)}, false}, // more entries
		{Map{"a": Int(1)}, Map{"a": Int(8)}, true},               // hash is smaller
		{Map{"a": Int(1)}, Map{"a": Int(3)}, false},              // hash is larger
	}

	Convey("When comparing a < b", t, func() {
		Convey("Then the result should be correct in all cases", func() {
			for _, tc := range ltTestCases {
				So(Less(tc.l, tc.r), ShouldEqual, tc.less)
				if tc.less {
					So(Less(tc.r, tc.l), ShouldNotEqual, tc.less)
				} else if Equal(tc.l, tc.r) {
					So(Less(tc.r, tc.l), ShouldEqual, tc.less)
				}

			}
		})
	})
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

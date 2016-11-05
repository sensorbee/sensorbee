package data

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecoder(t *testing.T) {
	Convey("Given a decoder with the default config", t, func() {
		d := NewDecoder(nil)

		type nested struct {
			X int     `bql:"nested_int"`
			Y float64 `bql:"nested_float"`
			Z string  `bql:"nested_str"`
		}
		s := &struct {
			B          bool
			WB         bool `bql:",weaklytyped"`
			I          int  `bql:",required"`
			F          float64
			S          string `bql:"str_key"`
			V          Value
			Map        Map
			FloatMap   map[string]float64
			Array      Array
			IntArray   []int
			Blob       []byte
			Struct     nested `bql:"nested"`
			Time       time.Time
			Timestamp  Timestamp // just in case
			Duration   time.Duration
			IPtr       *int
			MissingPtr *float64
			NestedPtr  *nested
		}{}

		Convey("When decoding a map", func() {
			now := time.Now()
			tsStr, _ := ToString(Timestamp(now))
			So(d.Decode(Map{
				"b":       True,
				"wb":      Array{Int(1)},
				"i":       Int(10),
				"f":       Float(3.14),
				"str_key": String("str"),
				"v": Map{
					"a": String("b"),
					"c": Int(1),
				},
				"map": Map{
					"key": String("value"),
				},
				"float_map": Map{
					"a": Float(1.2),
					"b": Float(3.4),
					"c": Float(5.6),
				},
				"array":     Array{Int(1), Float(2.3), String("4")},
				"int_array": Array{Int(1), Int(2), Int(3)},
				"blob":      Blob{4, 5, 6},
				"nested": Map{
					"nested_int":   Int(1),
					"nested_float": Float(2.3),
					"nested_str":   String("4"),
				},
				"time":      String(tsStr),
				"timestamp": Int(1),
				"duration":  String("5s"),
				"i_ptr":     Int(99),
				"nested_ptr": Map{
					"nested_int":   Int(10),
					"nested_float": Float(20.3),
					"nested_str":   String("40"),
				},
			}, s), ShouldBeNil)

			Convey("Then it should decode a boolean", func() {
				So(s.B, ShouldBeTrue)
			})

			Convey("Then it should decode a weakly typed boolean", func() {
				So(s.WB, ShouldBeTrue)
			})

			Convey("Then it should decode an integer", func() {
				So(s.I, ShouldEqual, 10)
			})

			Convey("Then it should decode a float", func() {
				So(s.F, ShouldEqual, 3.14)
			})

			Convey("Then it should decode a string", func() {
				So(s.S, ShouldEqual, "str")
			})

			Convey("Then it should decode a Value", func() {
				So(s.V, ShouldResemble, Map{
					"a": String("b"),
					"c": Int(1),
				})
			})

			Convey("Then it should decode a generic map", func() {
				So(s.Map, ShouldResemble, Map{
					"key": String("value"),
				})
			})

			Convey("Then it should decode a typed map", func() {
				So(s.FloatMap, ShouldResemble, map[string]float64{
					"a": 1.2,
					"b": 3.4,
					"c": 5.6,
				})
			})

			Convey("Then it should decode a generic array", func() {
				So(s.Array, ShouldResemble, Array{Int(1), Float(2.3), String("4")})
			})

			Convey("Then it should decode a typed array", func() {
				So(s.IntArray, ShouldResemble, []int{1, 2, 3})
			})

			Convey("Then it should decode a blob", func() {
				So(s.Blob, ShouldResemble, []byte{4, 5, 6})
			})

			Convey("Then it should decode a nested struct", func() {
				So(s.Struct, ShouldResemble, nested{1, 2.3, "4"})
			})

			Convey("Then it should decode a timestamp to time.Time", func() {
				So(s.Time, ShouldHappenOnOrBetween, now, now)
			})

			Convey("Then it should decode a timestamp to Timestamp", func() {
				So(time.Time(s.Timestamp), ShouldHappenOnOrBetween, time.Unix(1, 0), time.Unix(1, 0))
			})

			Convey("Then it should decode a duration", func() {
				So(s.Duration, ShouldEqual, 5*time.Second)
			})

			Convey("Then it should decode an integer to an *int", func() {
				So(*s.IPtr, ShouldEqual, 99)
			})

			Convey("Then it should not decode a float when the field is not given", func() {
				So(s.MissingPtr, ShouldBeNil)
			})

			Convey("Then it should decode a nested struct to *struct", func() {
				So(s.NestedPtr, ShouldResemble, &nested{10, 20.3, "40"})
			})
		})

		Convey("When decoding an integer/a float to float/int", func() {
			So(d.Decode(Map{
				"i": Float(10.0),
				"f": Int(3),
			}, s), ShouldBeNil)

			Convey("Then it should be converted implicitly", func() {
				So(s.I, ShouldEqual, 10)
				So(s.F, ShouldEqual, 3)
			})
		})

		Convey("When a required field is missing", func() {
			err := d.Decode(Map{
				"f": Float(3.14),
			}, s)

			Convey("Then it should fail", func() {
				fmt.Println(err)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When decoding non-addresable value", func() {
			err := d.Decode(Map{}, struct {
				I int
			}{})

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When decoding a non-struct value", func() {
			i := 10
			err := d.Decode(Map{}, &i)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When decoding a struct from a non-map value", func() {
			err := d.decodeStruct("", Int(1), reflect.Indirect(reflect.ValueOf(&struct {
				I int
			}{})))

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When a struct has an undefined option", func() {
			err := d.Decode(Map{"i": Int(1)}, &struct {
				I int `bql:",nosuchoption"`
			}{})

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When decoding time.Time directly", func() {
			err := d.Decode(Map{}, &time.Time{})

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When decoding Timestamp directly", func() {
			err := d.Decode(Map{}, &Timestamp{})

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestDecoderErrorReporting(t *testing.T) {
	Convey("Given a decode failure error report", t, func() {
		s := struct {
			I int `bql:",required"`
			A struct {
				B []struct {
					C int
					D int
				}
				M map[string]struct {
					E int
				}
			}
		}{}
		err := Decode(Map{
			"a": Map{
				"b": Array{
					Map{"c": Int(1)},
					Int(1),
					Map{"c": String("1"), "d": Array{}},
					Map{"c": Float(1.2)},
				},
				"m": Map{
					"a": Map{
						"d": Int(1),
					},
					"b": Map{
						"e": String("1"),
					},
					"c": Map{
						"e": Float(1.2),
					},
				},
			},
		}, &s)
		So(err, ShouldNotBeNil)
		report := err.Error()

		Convey("Then it should contain error about i", func() {
			So(report, ShouldContainSubstring, "i: ")
		})

		Convey("Then it should contain error about a.b[1]", func() {
			So(report, ShouldContainSubstring, "a.b[1]: ")
		})

		Convey("Then it should contain error about a.b[2].c", func() {
			So(report, ShouldContainSubstring, "a.b[2].c: ")
		})

		Convey("Then it should contain error about a.b[2].d", func() {
			So(report, ShouldContainSubstring, "a.b[2].d: ")
		})

		Convey("Then it should contain error about a.b[3].c", func() {
			So(report, ShouldContainSubstring, "a.b[3].c: ")
		})

		Convey(`Then it should contain error about a.m["a"]`, func() { // unused key
			So(report, ShouldContainSubstring, `a.m["a"]: `)
		})

		Convey(`Then it should contain error about a.m["b"].e`, func() {
			So(report, ShouldContainSubstring, `a.m["b"].e: `)
		})

		Convey(`Then it should contain error about a.m["c"].e`, func() {
			So(report, ShouldContainSubstring, `a.m["c"].e: `)
		})
	})
}

func TestDecoderMetadata(t *testing.T) {
	Convey("Given a decoder with Metadata + ErrorUnused", t, func() {
		s := struct {
			A int
			B int
			C int
		}{}

		md := &DecoderMetadata{}
		d := NewDecoder(&DecoderConfig{
			ErrorUnused: true,
			Metadata:    md,
		})

		Convey("When decoding all valid values", func() {
			err := d.Decode(Map{
				"a": Int(1),
				"b": Int(2),
				"c": Int(4),
			}, &s)
			So(err, ShouldBeNil)

			Convey("Then it should succeed", func() {
				So(s.A, ShouldEqual, 1)
				So(s.B, ShouldEqual, 2)
				So(s.C, ShouldEqual, 4)
			})

			Convey("Then the metadata should have all keys", func() {
				So(len(md.Keys), ShouldEqual, 3)
				So(md.Keys, ShouldContain, "a")
				So(md.Keys, ShouldContain, "b")
				So(md.Keys, ShouldContain, "c")
			})

			Convey("Then the metadata shouldn't contain any unused key", func() {
				So(md.Unused, ShouldBeEmpty)
			})
		})

		Convey("When decoding partially", func() {
			s.B = 3 // default value
			err := d.Decode(Map{
				"a": Int(1),
				"c": Int(4),
			}, &s)
			So(err, ShouldBeNil)

			Convey("Then it should succeed", func() {
				So(s.A, ShouldEqual, 1)
				So(s.B, ShouldEqual, 3)
				So(s.C, ShouldEqual, 4)
			})

			Convey("Then the metadata should have all processed keys", func() {
				So(len(md.Keys), ShouldEqual, 2)
				So(md.Keys, ShouldContain, "a")
				So(md.Keys, ShouldNotContain, "b")
				So(md.Keys, ShouldContain, "c")
			})

			Convey("Then the metadata shouldn't contain any unused key", func() {
				So(md.Unused, ShouldBeEmpty)
			})
		})

		Convey("When a map contains a undefined key", func() {
			err := d.Decode(Map{
				"a":  Int(1),
				"bb": Int(2),
				"c":  Int(4),
				"d":  Int(8),
			}, &s)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the metadata should only have processed keys", func() {
				So(len(md.Keys), ShouldEqual, 2)
				So(md.Keys, ShouldContain, "a")
				So(md.Keys, ShouldNotContain, "b")
				So(md.Keys, ShouldContain, "c")
			})

			Convey("Then the metadata should contain unused keys", func() {
				So(len(md.Unused), ShouldEqual, 2)
				So(md.Unused, ShouldContain, "bb")
				So(md.Unused, ShouldContain, "d")
			})
		})
	})
}

func TestDecodeWeaklytyped(t *testing.T) {
	Convey("Given decoding a struct filled with weaklytyped options", t, func() {
		s := &struct {
			B    bool               `bql:",weaklytyped"`
			I    int                `bql:",weaklytyped"`
			F    float64            `bql:",weaklytyped"`
			S    string             `bql:",weaklytyped"`
			M    map[string]float64 `bql:",weaklytyped"`
			A    []int              `bql:",weaklytyped"`
			Blob []byte             `bql:",weaklytyped"`
			IPtr *int               `bql:",weaklytyped"`
		}{}
		So(Decode(Map{
			"b": Map{"a": Int(1)},
			"i": Timestamp(time.Unix(2, 0)),
			"f": String("3.14"),
			"s": Int(4),
			"m": Map{
				"a": True,
				"b": String("2.3"),
			},
			"a":     Array{Float(3.4), String("5")},
			"blob":  String("MTIzLjQ1Ng=="),
			"i_ptr": String("6"),
		}, s), ShouldBeNil)

		Convey("Decoding a bool should succeed", func() {
			So(s.B, ShouldBeTrue)
		})

		Convey("Decoding an int should succeed", func() {
			So(s.I, ShouldEqual, 2)
		})

		Convey("Decoding a float should succeed", func() {
			So(s.F, ShouldEqual, 3.14)
		})

		Convey("Decoding a string should succeed", func() {
			So(s.S, ShouldEqual, "4")
		})

		Convey("Decoding a map should succeed", func() {
			So(s.M, ShouldResemble, map[string]float64{
				"a": 1,
				"b": 2.3,
			})
		})

		Convey("Decoding an array should succeed", func() {
			So(s.A, ShouldResemble, []int{3, 5})
		})

		Convey("Decoding a blob should succeed", func() {
			So(s.Blob, ShouldResemble, []byte("123.456"))
		})

		Convey("Decoding a pointer should succeed", func() {
			So(*s.IPtr, ShouldEqual, 6)
		})
	})
}

func TestDecodeUnsupportedType(t *testing.T) {
	// All possible combinations are tested in type conversion tests, so these
	// tests only have to check if ToType is called.

	Convey("Given a struct having unsupported type", t, func() {
		Convey("An interface other than Value should fail", func() {
			So(Decode(Map{"v": Int(1)}, &struct{ V interface{} }{}), ShouldNotBeNil)
		})

		Convey("Non supported type should fail", func() {
			So(Decode(Map{"v": Int(1)}, &struct{ V chan int }{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-boolean value to bool should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V bool }{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-integer value to int should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V int }{}), ShouldNotBeNil)
		})

		Convey("Decoding a value that is not convertible to int should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Array{}}, &struct {
				V int `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-float value to float should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V float64 }{}), ShouldNotBeNil)
		})

		Convey("Decoding a value that is not convertible to float should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Array{}}, &struct {
				V float64 `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-string value to string should fail", func() {
			So(Decode(Map{"v": Int(1)}, &struct{ V string }{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-map value to map should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V map[string]int }{}), ShouldNotBeNil)
		})

		Convey("Decoding a value that is not convertible to map should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Array{}}, &struct {
				V map[string]int `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding a map value having incompatible value type should fail", func() {
			So(Decode(Map{"v": Map{"a": String("1")}}, &struct{ V map[string]int }{}), ShouldNotBeNil)
		})

		Convey("Decoding a map value having weakly-incompatible value type should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Map{"a": String("a")}}, &struct {
				V map[string]int `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding to a map with non-string key should fail", func() {
			So(Decode(Map{"v": Map{"a": Int(1)}}, &struct{ V map[int]int }{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-array value to array should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V []int }{}), ShouldNotBeNil)
		})

		Convey("Decoding a value that is not convertible to array should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Map{}}, &struct {
				V []int `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding an array value having incompatible value type should fail", func() {
			So(Decode(Map{"v": Array{String("1")}}, &struct{ V []int }{}), ShouldNotBeNil)
		})

		Convey("Decoding an array value having weakly-incompatible value type should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Array{Map{}}}, &struct {
				V []int `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-blob value to blob should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V []byte }{}), ShouldNotBeNil)
		})

		Convey("Decoding a value that is not convertible to blob should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Map{}}, &struct {
				V []byte `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-timestamp value to time.Time should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V time.Time }{}), ShouldNotBeNil)
		})

		Convey("Decoding a value that is not convertible to time.Time should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Array{}}, &struct {
				V time.Time `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding a non-timestamp value to data.Timestamp should fail", func() {
			So(Decode(Map{"v": String("1")}, &struct{ V Timestamp }{}), ShouldNotBeNil)
		})

		Convey("Decoding a value that is not convertible to data.Timestamp should fail even with weaklytyped", func() {
			So(Decode(Map{"v": Array{}}, &struct {
				V Timestamp `bql:",weaklytyped"`
			}{}), ShouldNotBeNil)
		})

		Convey("Decoding a timestamp value to unsupported time struct should fail", func() {
			type MyTime time.Time
			So(Decode(Map{"v": Timestamp(time.Now())}, &struct{ V MyTime }{}), ShouldNotBeNil)
		})

		Convey("Decoding duration with a wrong format should fail", func() {
			So(Decode(Map{"v": String("10gs")}, &struct{ V time.Duration }{}), ShouldNotBeNil)
		})

		Convey("Decoding duration from an unsupported should fail", func() {
			So(Decode(Map{"v": Map{}}, &struct{ V time.Duration }{}), ShouldNotBeNil)
		})

		Convey("Decoding a pointer with incompatible value should fail", func() {
			So(Decode(Map{"v": Map{}}, &struct{ V *int }{}), ShouldNotBeNil)
		})
	})
}

func TestToSnakeCase(t *testing.T) {
	Convey("toSnakeCase should transform camelcase to snake case", t, func() {
		cases := [][]string{
			{"Test", "test"},
			{"ParseBQL", "parse_bql"},
			{"B2B", "b_2_b"},
		}

		for _, c := range cases {
			So(toSnakeCase(c[0]), ShouldEqual, c[1])
		}
	})
}

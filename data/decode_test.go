package data

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecoder(t *testing.T) {
	Convey("Given a decoder with the default config", t, func() {
		d := NewDecoder(nil)

		s := &struct {
			B bool
			I int     `bql:",required"`
			F float64 `bql:",weaklytyped"`
			S string  `bql:"str_key"`
			// TODO: support generic array when decoder supports Value
			IntArray []int
		}{}

		Convey("When decoding a map", func() {
			So(d.Decode(Map{
				"b":         True,
				"i":         Int(10),
				"f":         Float(3.14),
				"str_key":   String("str"),
				"int_array": Array{Int(1), Int(2), Int(3)},
			}, s), ShouldBeNil)

			Convey("Then it should decode a boolean", func() {
				So(s.B, ShouldBeTrue)
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

			Convey("Then it should decode a typed array", func() {
				So(s.IntArray, ShouldResemble, []int{1, 2, 3})
			})
		})

		Convey("When decoding a weakly typed value", func() {
			So(d.Decode(Map{
				"i": Int(10),
				"f": String("3.14"),
			}, s), ShouldBeNil)

			Convey("Then it should be decoded as a float", func() {
				So(s.F, ShouldEqual, 3.14)
			})
		})

		Convey("When a required field is missing", func() {
			err := d.Decode(Map{
				"f": String("3.14"),
			}, s)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

// TODO: weaklytype tests for each type
// All possible combinations are tested in type conversion tests, so these
// tests only have to check if ToType is called.

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

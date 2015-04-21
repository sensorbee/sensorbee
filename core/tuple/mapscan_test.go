package tuple

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScanMap(t *testing.T) {
	var testData = Map{
		"string": String("homhom"),
		"array":  Array([]Value{String("saysay"), String("mammam")}),
		"map": Map{
			"nested string": String("nested foo"),
			"nested_string": String("nested hoo"),
			"nested.string": String("nested loo"),
		},
	}
	Convey("Given a Map with values in it", t, func() {
		Convey("When accessing a empty string key", func() {
			var v Value
			err := scanMap(testData, "", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "empty key is not supported")
			})
		})
		Convey("When accessing an invalid string key", func() {
			var v Value
			err := scanMap(testData, "ab[a", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid path phrase")
			})
		})
		Convey("When accessing a non-existing key", func() {
			var v Value
			err := scanMap(testData, "str", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "not found the key in map: str")
			})
		})
		Convey("When accessing an invalid index in array key", func() {
			var v Value
			err := scanMap(testData, "array[2147483648]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "overflow index number: array[2147483648]")
			})
		})
		Convey("When accessing an invalid array key", func() {
			var v Value
			err := scanMap(testData, "string[0]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unsupported cast array from string")
			})
		})
		Convey("When accessing an out-of-range index", func() {
			var v Value
			err := scanMap(testData, "array[2]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "out of range access: array[2]")
			})
		})

		Convey("When accessing nonsense bracket holders", func() {
			var x1, x2 Value
			err1 := scanMap(testData, `["string']`, &x1)
			err2 := scanMap(testData, `['string"]`, &x2)
			Convey("Then lookup should fail", func() {
				So(err1, ShouldNotBeNil)
				So(err2, ShouldNotBeNil)
			})
		})
		Convey("When accessing bracket holders", func() {
			var x1, x2 Value
			err1 := scanMap(testData, `["string"]`, &x1)
			err2 := scanMap(testData, "['string']", &x2)
			Convey("Then lookup should exist", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				Convey("and is should be match the original value", func() {
					s1, _ := x1.String()
					s2, _ := x2.String()
					So(s1, ShouldEqual, "homhom")
					So(s2, ShouldEqual, "homhom")
				})
			})
		})
		Convey("When accessing nested bracket holders", func() {
			var x1, x2, x3, x4, x5 Value
			err1 := scanMap(testData, `["map"]["nested string"]`, &x1)
			err2 := scanMap(testData, `map['nested.string']`, &x2)
			err3 := scanMap(testData, `["map"]nested_string`, &x3)
			err4 := scanMap(testData, `map.['nested.string']`, &x4)
			err5 := scanMap(testData, `["map"].nested_string`, &x5)
			Convey("Then lookup should exist", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				So(err3, ShouldBeNil)
				So(err4, ShouldBeNil)
				So(err5, ShouldBeNil)
				Convey("and is should be match the original value", func() {
					s1, _ := x1.String()
					s2, _ := x2.String()
					s3, _ := x3.String()
					s4, _ := x4.String()
					s5, _ := x5.String()
					So(s1, ShouldEqual, "nested foo")
					So(s2, ShouldEqual, "nested loo")
					So(s3, ShouldEqual, "nested hoo")
					So(s4, ShouldEqual, "nested loo")
					So(s5, ShouldEqual, "nested hoo")
				})
			})
		})
		Convey("When accessing bracket array holder", func() {
			var v Value
			err := scanMap(testData, "['array'][0]", &v)
			Convey("Then lookup should exist", func() {
				So(err, ShouldBeNil)
				Convey("and is should be match the original value", func() {
					s, _ := v.String()
					So(s, ShouldEqual, "saysay")
				})
			})
		})
	})
}

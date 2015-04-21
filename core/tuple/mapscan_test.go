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
		"string['":       String("funnypath"),
		"string['aaa":    String("funnypath2"),
		"string['aaa\"]": String("funnypath3"),
		"string[\"aaa]":  String("funnypath4"),
	}
	Convey("Given a Map with values in it", t, func() {
		Convey("When accessing a empty string key", func() {
			_, err := testData.Get("")
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "empty key is not supported")
			})
		})
		Convey("When accessing an invalid string key", func() {
			_, err := testData.Get("ab[a")
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid path phrase")
			})
		})
		Convey("When accessing a non-existing key", func() {
			_, err := testData.Get("str")
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "not found the key in map: str")
			})
		})
		Convey("When accessing an invalid index in array key", func() {
			_, err := testData.Get("array[2147483648]")
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "overflow index number: array[2147483648]")
			})
		})
		Convey("When accessing an invalid array key", func() {
			_, err := testData.Get("string[0]")
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unsupported cast array from string")
			})
		})
		Convey("When accessing an out-of-range index", func() {
			_, err := testData.Get("array[2]")
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "out of range access: array[2]")
			})
		})

		Convey("When accessing nonsense bracket holders", func() {
			_, err1 := testData.Get(`["string']`)
			_, err2 := testData.Get(`['string"]`)
			Convey("Then lookup should fail", func() {
				So(err1, ShouldNotBeNil)
				So(err2, ShouldNotBeNil)
			})
		})
		Convey("When accessing bracket holders", func() {
			x1, err1 := testData.Get(`["string"]`)
			x2, err2 := testData.Get("['string']")
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
			x1, err1 := testData.Get(`["map"]["nested string"]`)
			x2, err2 := testData.Get(`map['nested.string']`)
			x3, err3 := testData.Get(`["map"]nested_string`)
			x4, err4 := testData.Get(`map.['nested.string']`)
			x5, err5 := testData.Get(`["map"].nested_string`)
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
	})
}

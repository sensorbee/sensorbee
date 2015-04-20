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
			"string": String("homhom"),
		},
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
				So(err.Error(), ShouldEqual, "not found the key in map")
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
	})
}

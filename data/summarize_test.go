package data

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSummarize(t *testing.T) {
	Convey("Given a map", t, func() {
		m, err := NewMap(map[string]interface{}{
			"int":    1,
			"bool":   false,
			"string": "hoge",
			"blob":   []byte("this is a blob"),
			"map": map[string]interface{}{
				"nested_blob": []byte("another blob"),
			},
			"array": []interface{}{
				[]byte("yet another blob"),
			},
		})
		So(err, ShouldBeNil)

		Convey("When summarizing it", func() {
			s := summarize(m)

			Convey("Then a blob should be replaced with (blob)", func() {
				v, err := s.(Map).Get("blob")
				So(err, ShouldBeNil)
				So(v, ShouldEqual, "(blob)")
			})

			Convey("Then a blob nested in a map should be replaced with (blob)", func() {
				v, err := s.(Map).Get("map.nested_blob")
				So(err, ShouldBeNil)
				So(v, ShouldEqual, "(blob)")
			})

			Convey("Then a blob nested in an array should be replaced with (blob)", func() {
				v, err := s.(Map).Get("array[0]")
				So(err, ShouldBeNil)
				So(v, ShouldEqual, "(blob)")
			})
		})

		Convey("When summarizing it to a JSON-like string", func() {
			s := Summarize(m)

			Convey("Then a blob should be replaced with (blob)", func() {
				So(s, ShouldNotContainSubstring, "this is a blob")
			})

			Convey("Then a blob nested in a map should be replaced with (blob)", func() {
				So(s, ShouldNotContainSubstring, "another blob")
			})

			Convey("Then a blob nested in an array should be replaced with (blob)", func() {
				So(s, ShouldNotContainSubstring, "yet")
			})
		})
	})
}

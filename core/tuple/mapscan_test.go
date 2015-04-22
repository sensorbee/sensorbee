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

func TestGetArrayIndex(t *testing.T) {
	Convey("Given strings including array index", t, func() {
		s := "[0]aaa[1]"
		Convey("Then return valid index", func() {
			index1 := getArrayIndex([]rune(s), 1)
			index2 := getArrayIndex([]rune(s), 7)
			So(index1, ShouldEqual, "[0]")
			So(index2, ShouldEqual, "[1]")
		})
		Convey("Then return empty index", func() {
			index1 := getArrayIndex([]rune(s), 0)
			index2 := getArrayIndex([]rune(s), 2)
			So(index1, ShouldEqual, "")
			So(index2, ShouldEqual, "")

		})
	})
}

func TestSplitBracket(t *testing.T) {
	Convey("Given strings start with bracket", t, func() {
		Convey("When valid expressions", func() {
			s := []string{`[""]abc`, `["a"]"]`, `['a']['`, `["ab\"]`, `["a]b"]`, `["a"b"]`}
			Convey("Then return inner strings", func() {
				expects := []string{"", "a", "a", "ab\\", "a]b", "a\"b"}
				actuals := make([]string, len(s))
				actuals[0] = splitBracket([]rune(s[0]), 2, '"')
				actuals[1] = splitBracket([]rune(s[1]), 2, '"')
				actuals[2] = splitBracket([]rune(s[2]), 2, '\'')
				actuals[3] = splitBracket([]rune(s[3]), 2, '"')
				actuals[4] = splitBracket([]rune(s[4]), 2, '"')
				actuals[5] = splitBracket([]rune(s[5]), 2, '"')
				for i, a := range actuals {
					So(a, ShouldEqual, expects[i])
				}
			})
		})
		Convey("When invalid invalid expressions", func() {
			s := []string{`["a`, `['a`, `['a"]`, `["ab']`, `["a]`, `["a"`}
			Convey("Then return inner strings", func() {
				actuals := make([]string, len(s))
				actuals[0] = splitBracket([]rune(s[0]), 2, '"')
				actuals[1] = splitBracket([]rune(s[1]), 2, '\'')
				actuals[2] = splitBracket([]rune(s[2]), 2, '\'')
				actuals[3] = splitBracket([]rune(s[3]), 2, '"')
				actuals[4] = splitBracket([]rune(s[4]), 2, '"')
				actuals[5] = splitBracket([]rune(s[5]), 2, '"')
				for _, a := range actuals {
					So(a, ShouldEqual, "")
				}
			})
		})
	})
}

func TestSplit(t *testing.T) {
	Convey("Given path expressions", t, func() {
		Convey("When valid expressions", func() {
			phrase1 := `path1.pa\.th2.[pa]th3["path4"]['path5'][0]path6[0]`
			phrase2 := `path1..path2.["path3\"].['pat.h4'].path5`
			phrase3 := `path1["path1']`
			phrase4 := `path2['path2"]`
			Convey("Then split", func() {
				expected1 := []string{"path1", "pa.th2", "[pa]th3", "path4", "path5[0]", "path6[0]"}
				expected2 := []string{"path1", "path2", "path3\\", "pat.h4", "path5"}
				expected3 := []string{"path1[\"path1']"}
				expected4 := []string{"path2['path2\"]"}
				actual1 := split(phrase1)
				actual2 := split(phrase2)
				actual3 := split(phrase3)
				actual4 := split(phrase4)
				So(actual1, ShouldResemble, expected1)
				So(actual2, ShouldResemble, expected2)
				So(actual3, ShouldResemble, expected3)
				So(actual4, ShouldResemble, expected4)
			})
		})
	})
}

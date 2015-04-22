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
			"nested.string":    String("nested foo"),
			"nested..string":   String("nested loo"),
			"nestedstring]":    String("nested qoo"),
			"nested\"string":   String("nested woo"),
			"nested'string":    String("nested roo"),
			"nested\\\"string": String("nested zoo"),
			"nestedstring\\":   String("nested xoo"),
			"nested\nstring":   String("nested coo"),
			"内部マップ":            String("内部ﾏｯﾌﾟ"),
			"nestedstring":     String("nested hoo"),
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
			Convey("With valid nested bracket", func() {
				var x1, x2, x3, x4, x5, x6, x7, x8, x9 Value
				err1 := scanMap(testData, `map["nested.string"]`, &x1)
				err2 := scanMap(testData, `map['nested..string']`, &x2)
				err3 := scanMap(testData, `map['nestedstring]']`, &x3)
				err4 := scanMap(testData, `map['nested"string']`, &x4)
				err5 := scanMap(testData, `map['nested'string']`, &x5)
				err6 := scanMap(testData, `map['nested\"string']`, &x6)
				err7 := scanMap(testData, `map['nestedstring\']`, &x7)
				err8 := scanMap(testData, "map['nested\nstring']", &x8)
				err9 := scanMap(testData, `map['内部マップ']`, &x9)
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				So(err3, ShouldBeNil)
				So(err4, ShouldBeNil)
				So(err5, ShouldBeNil)
				So(err6, ShouldBeNil)
				So(err7, ShouldBeNil)
				So(err8, ShouldBeNil)
				So(err9, ShouldBeNil)
				Convey("Then lookup should be match the original value", func() {
					s1, _ := x1.String()
					s2, _ := x2.String()
					s3, _ := x3.String()
					s4, _ := x4.String()
					s5, _ := x5.String()
					s6, _ := x6.String()
					s7, _ := x7.String()
					s8, _ := x8.String()
					s9, _ := x9.String()
					So(s1, ShouldEqual, "nested foo")
					So(s2, ShouldEqual, "nested loo")
					So(s3, ShouldEqual, "nested qoo")
					So(s4, ShouldEqual, "nested woo")
					So(s5, ShouldEqual, "nested roo")
					So(s6, ShouldEqual, "nested zoo")
					So(s7, ShouldEqual, "nested xoo")
					So(s8, ShouldEqual, "nested coo")
					So(s9, ShouldEqual, "内部ﾏｯﾌﾟ")
				})
			})
			Convey("With invalid nested bracket", func() {
				var x1, x2, x3, x4, x5, x6 Value
				err1 := scanMap(testData, `map[]`, &x1)
				err2 := scanMap(testData, `map..nestedstring`, &x2)
				err3 := scanMap(testData, `map[nestedstring]`, &x3)
				err4 := scanMap(testData, `map["nestedstring']`, &x4)
				err5 := scanMap(testData, `map['nestedstring"]`, &x5)
				err6 := scanMap(testData, `map[nested.string]`, &x6)
				So(err1, ShouldNotBeNil)
				SkipSo(err2, ShouldNotBeNil)
				So(err3, ShouldNotBeNil)
				So(err4, ShouldNotBeNil)
				So(err5, ShouldNotBeNil)
				So(err6, ShouldNotBeNil)
			})
		})
		Convey("When accessing bracket array holder", func() {
			Convey("With valid array index should exist", func() {
				var v1, v2, v3, v4 Value
				err1 := scanMap(testData, "['array'][0]", &v1)
				err2 := scanMap(testData, "array[0]", &v2)
				err3 := scanMap(testData, "array.[0]", &v3)
				err4 := scanMap(testData, "array[0].", &v4)
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				So(err3, ShouldBeNil)
				So(err4, ShouldBeNil)
				Convey("Then lookup should be match the original value", func() {
					s1, _ := v1.String()
					s2, _ := v2.String()
					s3, _ := v3.String()
					s4, _ := v4.String()
					So(s1, ShouldEqual, "saysay")
					So(s2, ShouldEqual, "saysay")
					So(s3, ShouldEqual, "saysay")
					So(s4, ShouldEqual, "saysay")
				})
			})
			Convey("With invalid array index", func() {
				var v Value
				err1 := scanMap(testData, "array[０]", &v)
				err2 := scanMap(testData, "array[0][", &v)
				err3 := scanMap(testData, "array[0]]", &v)
				err4 := scanMap(testData, "array[-1]", &v)
				err5 := scanMap(testData, "array[0:1]", &v)
				So(err1.Error(), ShouldEqual, "invalid path phrase")
				So(err2.Error(), ShouldEqual, "invalid path phrase")
				So(err3.Error(), ShouldEqual, "invalid path phrase")
				So(err4.Error(), ShouldEqual, "invalid path phrase")
				So(err5.Error(), ShouldEqual, "invalid path phrase")
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
		Convey("When invalid expressions", func() {
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

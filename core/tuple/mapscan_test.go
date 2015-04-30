package tuple

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMapscanDocstrings(t *testing.T) {
	storeData := Map{
		"name": String("store name"),
		"book": Array([]Value{
			Map{
				"title": String("book name"),
			},
		}),
	}
	m := Map{"store": storeData}

	Convey("Map.Get examples should be correct", t, func() {
		examples := map[string]interface{}{
			"store":                         storeData,
			"store.name":                    String("store name"),
			"store.book[0].title":           String("book name"),
			`["store"]`:                     storeData,
			`["store"]["name"]`:             String("store name"),
			`["store"]["book"][0]["title"]`: String("book name"),
		}
		for input, expected := range examples {
			actual, err := m.Get(input)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		}
	})

	Convey("split examples should be correct", t, func() {
		examples := map[string]interface{}{
			"store.book[0].title":           []string{"store", "book[0]", "title"},
			`["store"]["book"][0]["title"]`: []string{"store", "book[0]", "title"},
		}
		for input, expected := range examples {
			actual := split(input)
			So(actual, ShouldResemble, expected)
		}
	})

	Convey("splitBracket examples should be correct", t, func() {
		examples := map[string]interface{}{
			`a["hoge"].b`:    "hoge",
			`a["hoge"][123]`: "hoge[123]",
		}
		for input, expected := range examples {
			actual := splitBracket([]rune(input), 3, '"')
			So(actual, ShouldResemble, expected)
		}
	})

	Convey("getArrayIndex examples should be correct", t, func() {
		examples := map[string]interface{}{
			`hoge[123]`: "[123]",
		}
		for input, expected := range examples {
			actual := getArrayIndex([]rune(input), 5)
			So(actual, ShouldResemble, expected)
		}
	})
}

func TestScanMap(t *testing.T) {
	nestedData := Map{
		"nested.string":    String("keywithdot"),
		"nested..string":   String("keywithtwodots"),
		"nestedstring]":    String("keywithbracket"),
		"nested\"string":   String("keywithdoublequote"),
		"nested'string":    String("keywithsinglequote"),
		"nested\\\"string": String("keywithescapeddoublequote"),
		"nestedstring\\":   String("keywithbackslash"),
		"nested\nstring":   String("keywithnewline"),
		"内部マップ":            String("内部ﾏｯﾌﾟ"),
		"nestedstring":     String("normalkey"),
	}
	var testData = Map{
		"string":   String("homhom"),
		"array":    Array([]Value{String("saysay"), String("mammam")}),
		"map":      nestedData,
		"arraymap": Array([]Value{Map{"mappedstring": String("boo")}}),
	}
	Convey("Given a Map with values in it", t, func() {
		Convey("When accessing an empty string key", func() {
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
				So(err.Error(), ShouldEqual, "invalid path component: ab[a")
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
				So(err.Error(), ShouldEqual, "cannot access a tuple.String using index 0")
			})
		})
		Convey("When accessing only an array key", func() {
			var v Value
			err := scanMap(testData, "[0]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "cannot access a tuple.Map using index 0")
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
			Convey("Then lookup should succeed", func() {
				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				Convey("and it should be match the original value", func() {
					s1, _ := x1.String()
					s2, _ := x2.String()
					So(s1, ShouldEqual, "homhom")
					So(s2, ShouldEqual, "homhom")
				})
			})
		})
		Convey("When accessing nested bracket holders", func() {
			Convey("With valid nested bracket", func() {
				samples := map[string]string{
					`map["nested.string"]`:  "keywithdot",
					`map['nested..string']`: "keywithtwodots",
					`map['nestedstring]']`:  "keywithbracket",
					`map['nested"string']`:  "keywithdoublequote",
					`map['nested'string']`:  "keywithsinglequote",
					`map['nested\"string']`: "keywithescapeddoublequote",
					`map['nestedstring\']`:  "keywithbackslash",
					"map['nested\nstring']": "keywithnewline",
					`map['内部マップ']`:          "内部ﾏｯﾌﾟ",
				}
				Convey("Then lookup should succeed and match the original value", func() {
					for input, expected := range samples {
						var v Value
						err := scanMap(testData, input, &v)
						So(err, ShouldBeNil)
						actual, _ := v.String()
						So(actual, ShouldEqual, expected)
					}
				})
			})
			Convey("With invalid nested bracket", func() {
				samples := []string{
					`map[]`,
					`map..nestedstring`, // currently invalid path
					`map[nestedstring]`,
					`map["nestedstring']`,
					`map['nestedstring"]`,
					`map[nested.string]`,
					`string["string"]`,
				}
				Convey("Then lookup should fail", func() {
					for idx, input := range samples {
						var v Value
						err := scanMap(testData, input, &v)
						if idx == 1 {
							SkipSo(err, ShouldNotBeNil)
						} else {
							So(err, ShouldNotBeNil)
						}

					}
				})
			})
		})
		Convey("When accessing bracket array holder", func() {
			Convey("With valid array index", func() {
				samples := map[string]string{
					`['array'][0]`:                    "saysay",
					`array[0]`:                        "saysay",
					`array.[0]`:                       "saysay",
					`array[0].`:                       "saysay",
					`['arraymap'][0]['mappedstring']`: "boo",
				}
				Convey("Then lookup should succeed and match the original value", func() {
					for input, expected := range samples {
						var v Value
						err := scanMap(testData, input, &v)
						So(err, ShouldBeNil)
						actual, _ := v.String()
						So(actual, ShouldEqual, expected)
					}
				})
			})
			Convey("With invalid array index", func() {
				samples := []string{
					`array[０]`, // zenkaku
					`array[0][`,
					`array[0]]`,
					`array[-1]`,
					`array[0:1]`,
				}
				Convey("Then lookup should fail", func() {
					for _, input := range samples {
						var v Value
						err := scanMap(testData, input, &v)
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "invalid path component: "+input)

					}
				})
			})
		})
	})
}

func TestGetArrayIndex(t *testing.T) {
	Convey("Given strings including array index", t, func() {
		s := "[0]aaa[1]"
		Convey("Then getArrayIndex returns that index", func() {
			index1 := getArrayIndex([]rune(s), 1)
			index2 := getArrayIndex([]rune(s), 7)
			So(index1, ShouldEqual, "[0]")
			So(index2, ShouldEqual, "[1]")
		})
		Convey("Then getArrayIndex returns an empty string", func() {
			index1 := getArrayIndex([]rune(s), 0)
			index2 := getArrayIndex([]rune(s), 2)
			So(index1, ShouldEqual, "")
			So(index2, ShouldEqual, "")

		})
	})
}

func TestSplitBracket(t *testing.T) {
	Convey("Given a string that starts with bracket", t, func() {
		Convey("When the expression is valid", func() {
			samples := map[string]string{
				`[""]abc`:       "",
				`["a"]"]`:       "a",
				`['a']['`:       "a",
				`["ab\"]`:       "ab\\",
				`["a]b"]`:       "a]b",
				`["a"b"]`:       "a\"b",
				"a[\"hoge\"].b": "hoge",
				`a["hog"e"].b`:  "hog\"e",
				`a["b"][123]`:   "b[123]",
				`a["b"][]`:      "b",
			}
			Convey("Then splitBracket() returns the contained string", func() {
				for input, expected := range samples {
					from := 2
					if input[0] == 'a' {
						from = 3
					}
					quote := '"'
					if input[1] == '\'' {
						quote = '\''
					}
					actual := splitBracket([]rune(input), from, quote)
					So(actual, ShouldEqual, expected)
				}
			})
		})
		Convey("When the expression is invalid", func() {
			samples := []string{
				`["a`,
				`['a`,
				`['a"]`,
				`["ab']`,
				`["a]`,
				`["a"`,
			}
			Convey("Then splitBracket() returns an empty string", func() {
				for _, input := range samples {
					quote := '"'
					if input[1] == '\'' {
						quote = '\''
					}
					actual := splitBracket([]rune(input), 2, quote)
					So(actual, ShouldEqual, "")
				}
			})
		})
	})
}

func TestSplit(t *testing.T) {
	Convey("Given a path expression", t, func() {
		samples := map[string][]string{
			`path1.pa\.th2.[pa]th3["path4"]['path5'][0]path6[0]`: []string{"path1", "pa.th2", "[pa]th3", "path4", "path5[0]", "path6[0]"},
			`path1..path2.["path3\"].['pat.h4'].path5`:           []string{"path1", "path2", "path3\\", "pat.h4", "path5"},
			`path1["path1']`:                                     []string{"path1[\"path1']"},
			`path2['path2"]`:                                     []string{"path2['path2\"]"},
		}
		Convey("Then split() returns a proper list of components", func() {
			for input, expected := range samples {
				actual := split(input)
				So(actual, ShouldResemble, expected)
			}
		})
	})
}

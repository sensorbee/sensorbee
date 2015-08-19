package data

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// scanMap is a legacy method that only exists to avoid
// large changes to the tests.
func scanMap(m Map, p string, v *Value) (err error) {
	path, err := CompilePath(p)
	if err != nil {
		return err
	}
	val, err := path.evaluate(m)
	if err != nil {
		return err
	}
	*v = val
	return nil
}

// setInMap is a legacy method that only exists to avoid
// large changes to the tests.
func setInMap(m Map, p string, v Value) (err error) {
	path, err := CompilePath(p)
	if err != nil {
		return err
	}
	return path.set(m, v)
}

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
			path, err := CompilePath(input)
			So(err, ShouldBeNil)
			actual, err := m.Get(path)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		}
	})
}

func TestArraySlicing(t *testing.T) {
	elem0 := Map{"hoge": Array{
		Map{"a": Int(1), "b": Int(2)},
		Map{"a": Int(3), "b": Int(4)},
	}, "bar": Int(5)}
	elem1 := Map{"hoge": Array{
		Map{"a": Int(5), "b": Int(6)},
		Map{"a": Int(7), "b": Int(8)},
	}, "bar": Int(2)}
	elem2 := Map{"hoge": Array{
		Map{"a": Int(9), "b": Int(10)},
	}, "bar": Int(8)}

	data := Map{
		"foo": Array{
			elem0, elem1, elem2,
		},
		"nantoka": Map{"x": String("y")},
	}

	Convey("Slicing should work as expected", t, func() {
		examples := map[string]interface{}{
			// normal array slicing
			"foo[0:1]":  Array{elem0},
			"foo[0:2]":  Array{elem0, elem1},
			"foo[1:3]":  Array{elem1, elem2},
			"foo[2:2]":  Array{},
			"foo[2:17]": Array{elem2},
			"foo[3:17]": Array{},
			// slicing with sub-elements
			"foo[0:1].bar":    Array{Int(5)},
			"foo[0:2].bar":    Array{Int(5), Int(2)},
			"foo[1:3].bar":    Array{Int(2), Int(8)},
			"foo[1:3]['bar']": Array{Int(2), Int(8)},
			"foo[2:2].bar":    Array{},
			"foo[2:17].bar":   Array{Int(8)},
			"foo[3:17].bar":   Array{},
			// slicing with further descend
			"foo[0:1].hoge[0].b":    Array{Int(2)},
			"foo[0:2].hoge[0].b":    Array{Int(2), Int(6)},
			"foo[0:2].hoge[1].b":    Array{Int(4), Int(8)},
			"foo[1:3].hoge[0].b":    Array{Int(6), Int(10)},
			"foo[1:3]['hoge'][0].b": Array{Int(6), Int(10)},
			"foo[1:3].hoge[1].b":    nil, // foo[2] has no hoge[1]
			"foo[2:2].hoge[0].b":    Array{},
			"foo[2:17].hoge[0].b":   Array{Int(10)},
			"foo[3:17].hoge[0].b":   Array{},
			// recursion
			"foo..bar":              Array{Int(5), Int(2), Int(8)},
			"foo..hoge":             Array{elem0["hoge"], elem1["hoge"], elem2["hoge"]},
			"foo..hoge[0].b":        Array{Int(2), Int(6), Int(10)},
			"foo..['hoge'][0]['b']": Array{Int(2), Int(6), Int(10)},
			"foo..b":                Array{Int(2), Int(4), Int(6), Int(8), Int(10)},
			"nantoka..x":            Array{String("y")},
			"foo..x":                Array{},
			"nantoka.x..a":          nil, // recursive access on non-container
		}
		for input, expected := range examples {
			path, err := CompilePath(input)
			So(err, ShouldBeNil)
			actual, err := data.Get(path)
			if expected == nil {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			}
		}
	})

	Convey("Illegal slice paths should return an error", t, func() {
		examples := []string{
			"foo[0:1][2:3]",
			"foo[0:1].hoge[2:3]",
			"foo[0:1].hoge[2:3].bar",
			"foo[0:1]..bar",
			"foo[3:2]",
			"foo[3:2].hoge[0].b",
		}

		for _, input := range examples {
			_, err := CompilePath(input)
			So(err, ShouldNotBeNil)
		}
	})
}

func TestSetInMap(t *testing.T) {
	testCases := []struct {
		key    string
		val    Value
		errmsg string
	}{
		/// successes
		// set top-level element
		{"store", Int(13), ""},
		// add top-level element
		{"string", Int(27), ""},
		// set second-level element
		{"store.name", Int(13), ""},
		// add second-level element
		{"store.id", Int(27), ""},
		// set item in list
		{"store.book[0]", Int(27), ""},
		// change item in list
		{"store.book[0].title", Int(27), ""},
		{"store.book[0].hoge", Int(13), ""},
		// append item to list
		{"store.book[1]", Int(27), ""},
		// add item to list with null-padding
		{"store.book[5]", Int(27), ""},
		// create parent map
		{"store.owner.name", String("bar foo"), ""},
		// create parent list
		{"store.owners[1].nickname", String("ore"), ""},
		// nested lists
		{"store.owners[1][2]", String("ore"), ""},
		/// fails
		// fail: add element below non-map
		{"store.name.hoge", Int(13), "cannot access a data.String using key \"hoge\""},
		// fail: set index in map
		{"store.book.hoge", Int(13), "cannot access a data.Array using key \"hoge\""},
		// fail: set key in array
		{"store[5]", Int(27), "cannot access a data.Map using index 5"},
	}

	Convey("Given a Map with values in it", t, func() {
		testData := Map{
			"store": Map{
				"name": String("store name"),
				"book": Array([]Value{
					Map{
						"title": String("book name"),
					},
				}),
			}}

		for _, testCase := range testCases {
			tc := testCase
			Convey(fmt.Sprintf("When setting the value at '%s' to %v", tc.key, tc.val), func() {
				path, err := CompilePath(tc.key)
				So(err, ShouldBeNil)
				err = testData.Set(path, tc.val)
				if tc.errmsg == "" {
					Convey("There should be no error", func() {
						So(err, ShouldBeNil)

						Convey("And Get() should get back the result", func() {
							getVal, err := testData.Get(path)
							So(err, ShouldBeNil)
							So(getVal, ShouldResemble, tc.val)
						})
					})
				} else {
					Convey("There should be an error", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, tc.errmsg)
					})
				}
			})
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
		"'nestedstring":    String("keywithsinglequoteA"),
		"nestedstring'":    String("keywithsinglequoteB"),
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
			})
		})
		Convey("When accessing an invalid string key", func() {
			var v Value
			err := scanMap(testData, "ab[a", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When accessing a non-existing key", func() {
			var v Value
			err := scanMap(testData, "str", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "key 'str' was not found in map")
			})
		})
		Convey("When accessing an invalid index in array key", func() {
			var v Value
			err := scanMap(testData, "array[2147483648]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "overflow index number: 2147483648")
			})
		})
		Convey("When accessing an invalid array key", func() {
			var v Value
			err := scanMap(testData, "string[0]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "cannot access a data.String using index 0")
			})
		})
		Convey("When accessing only an array key", func() {
			var v Value
			err := scanMap(testData, "[0]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When accessing an out-of-range index", func() {
			var v Value
			err := scanMap(testData, "array[2]", &v)
			Convey("Then lookup should fail", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "out of range access: 2")
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
					s1, err3 := x1.asString()
					s2, err4 := x2.asString()
					So(err3, ShouldBeNil)
					So(err4, ShouldBeNil)
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
						actual, err := v.asString()
						So(err, ShouldBeNil)
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
					`['arraymap'][0]['mappedstring']`: "boo",
				}
				Convey("Then lookup should succeed and match the original value", func() {
					for input, expected := range samples {
						var v Value
						err := scanMap(testData, input, &v)
						So(err, ShouldBeNil)
						actual, err := v.asString()
						So(err, ShouldBeNil)
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
				}
				Convey("Then lookup should fail", func() {
					for _, input := range samples {
						var v Value
						err := scanMap(testData, input, &v)
						So(err, ShouldNotBeNil)

					}
				})
			})
		})
	})
}

func BenchmarkMapAccess(b *testing.B) {
	m := Map{
		"store": Map{
			"name": String("store name"),
			"book": Array([]Value{
				Map{
					"title": String("book name"),
				},
			}),
		},
	}
	examples := map[string]interface{}{
		"store.name":                    String("store name"),
		"store.book[0].title":           String("book name"),
		`["store"]["name"]`:             String("store name"),
		`["store"]["book"][0]["title"]`: String("book name"),
	}
	for n := 0; n < b.N; n++ {
		for input, expected := range examples {
			path, err := CompilePath(input)
			So(err, ShouldBeNil)
			actual, err := m.Get(path)
			if err != nil {
				panic(err)
			}
			if actual != expected {
				panic("result mismatch")
			}
		}
	}
}

package data

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ugorji/go/codec"
	"math"
	"testing"
	"time"
)

// TestUnmarshalMsgpack tests that byte array encoded by msgpack
// is deserialized to Map object correctly
func TestUnmarshalMsgpack(t *testing.T) {
	Convey("Given a msgpack byte data", t, func() {
		now := time.Now()
		unixTime := now.Unix()
		var testMap = map[string]interface{}{
			"bool":    true,
			"int32":   int32(1),
			"int64":   int64(2),
			"float32": float32(0.1),
			"float64": float64(0.2),
			"string":  "homhom",
			"time":    int64(unixTime),
			"array": []interface{}{true, 10, "inarray",
				map[string]interface{}{
					"mapinarray": "arraymap",
				}},
			"map": map[string]interface{}{
				"map_a": "a",
				"map_b": 2,
			},
			"null": nil,
			// TODO add []byte
		}
		var testData []byte
		codec.NewEncoderBytes(&testData, msgpackHandle).Encode(testMap)
		Convey("When convert to Map object", func() {
			m, err := UnmarshalMsgpack(testData)
			So(err, ShouldBeNil)
			Convey("Then decode data should be match with Map data", func() {
				var expected = Map{
					"bool":    Bool(true),
					"int32":   Int(1),
					"int64":   Int(2),
					"float32": Float(float32(0.1)),
					"float64": Float(0.2),
					"string":  String("homhom"),
					"time":    Int(unixTime),
					"array": Array([]Value{Bool(true), Int(10), String("inarray"),
						Map{
							"mapinarray": String("arraymap"),
						}}),
					"map": Map{
						"map_a": String("a"),
						"map_b": Int(2),
					},
					"null": Null{},
				}
				So(m, ShouldResemble, expected)
			})
		})
	})
}

// TestNewMapDocMaps tests that supported type by SensorBee can be converted
// correctly. A test data is same with doc example.
func TestNewMapDocMaps(t *testing.T) {
	Convey("Given a map[string]interface{} including variable type value", t, func() {
		var m = map[string]interface{}{
			"bool":   true,
			"int":    int64(1),
			"float":  float64(0.1),
			"string": "homhom",
			"time":   time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC),
			"array": []interface{}{true, 10, "inarray",
				map[string]interface{}{
					"mapinarray": "arraymap",
				}},
			"map": map[string]interface{}{
				"map_a": "a",
				"map_b": 2,
			},
			"imap": map[interface{}]interface{}{
				"map_a": "b",
				"map_b": 3,
			},
			"byte":  []byte("test byte"),
			"null":  nil,
			"value": Array{Int(1), Float(2), String("hoge")},
		}
		var expected = Map{
			"bool":   Bool(true),
			"int":    Int(1),
			"float":  Float(0.1),
			"string": String("homhom"),
			"time":   Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)),
			"array": Array([]Value{Bool(true), Int(10), String("inarray"),
				Map{
					"mapinarray": String("arraymap"),
				}}),
			"map": Map{
				"map_a": String("a"),
				"map_b": Int(2),
			},
			"imap": Map{
				"map_a": String("b"),
				"map_b": Int(3),
			},
			"byte":  Blob([]byte("test byte")),
			"null":  Null{},
			"value": Array{Int(1), Float(2), String("hoge")},
		}
		Convey("When convert to Map object", func() {
			actual, err := NewMap(m)
			Convey("Then test data should be converted correctly", func() {
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})
	})
}

// TestNewMapIncludeUnsupportedValue tests that unsupported value
// (e.g. complex64) can not be converted
func TestNewMapIncludeUnsupportedValue(t *testing.T) {
	Convey("Given a map[string]interface{} including unsupported value", t, func() {
		var m = map[string]interface{}{
			"errortype": complex64(1 + 1i),
		}
		Convey("When convert to Map object", func() {
			_, err := NewMap(m)
			Convey("Then error should be occurred", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
	Convey("Given a map[string]interface{} including overflow value", t, func() {
		mxint64 := uint64(math.MaxInt64 + 1)
		var m = map[string]interface{}{
			"errorvalue": mxint64,
		}
		Convey("When convert to Map object", func() {
			_, err := NewMap(m)
			Convey("Then error should be occurred", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, fmt.Sprintf(
					"an int value must be less than 2^63: %v", mxint64))
			})
		})
	})
}

// TestMarshalMsgpack tests that Map object is serialized to
// byte array correctly.
func TestMarshalMsgpack(t *testing.T) {
	Convey("Given a Map object data", t, func() {
		now := time.Now()
		unitTime := now.Unix()
		var testMap = Map{
			"bool":   Bool(true),
			"int":    Int(1),
			"float":  Float(0.1),
			"string": String("homhom"),
			"time":   Timestamp(now),
			"array": Array([]Value{Bool(true), Int(10), String("inarray"),
				Map{
					"mapinarray": String("arraymap"),
				}}),
			"map": Map{
				"map_a": String("a"),
				"map_b": Int(2),
			},
			"null": Null{},
			// TODO add Blob
		}
		Convey("When convert to []byte", func() {
			b, err := MarshalMsgpack(testMap)
			So(err, ShouldBeNil)
			Convey("Then encode data should be match with expected bytes", func() {
				var expected = map[string]interface{}{
					"bool":   true,
					"int":    int64(1),
					"float":  float64(0.1),
					"string": "homhom",
					"time":   unitTime,
					"array": []interface{}{true, 10, "inarray",
						map[string]interface{}{
							"mapinarray": "arraymap",
						}},
					"map": map[string]interface{}{
						"map_a": "a",
						"map_b": 2,
					},
					"null": nil,
				}
				var expectedBytes []byte
				codec.NewEncoderBytes(&expectedBytes, msgpackHandle).Encode(expected)

				// it should compare b and expectedBytes, but byte array order is not
				// always correspond in converting map to bytes.
				var actualMap, expectedMap map[string]interface{}
				codec.NewDecoderBytes(expectedBytes, msgpackHandle).Decode(&expectedMap)
				codec.NewDecoderBytes(b, msgpackHandle).Decode(&actualMap)
				So(actualMap, ShouldResemble, expectedMap)
			})
		})
	})
}

func TestValue(t *testing.T) {
	var testData = Map{
		"bool":   Bool(true),
		"int":    Int(1),
		"float":  Float(0.1),
		"string": String("homhom"),
		"byte":   Blob([]byte("madmad")),
		"time":   Timestamp(time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC)),
		"array":  Array([]Value{String("saysay"), String("mammam")}),
		"map": Map{
			"string": String("homhom2"),
		},
		"null": Null{},
	}

	Convey("Given a Map with values in it", t, func() {
		Convey("When accessing a non-existing key", func() {
			_, getErr := testData.Get(MustCompilePath("key"))
			Convey("Then lookup should fail", func() {
				So(getErr, ShouldNotBeNil)
			})
		})

		Convey("When deep-copying the value", func() {
			copy := testData.Copy()
			Convey("Then all values should be the same", func() {
				simpleTypes := []string{"bool", "int", "float", "string",
					"array[0]", "map.string"}
				for _, typeName := range simpleTypes {
					path, err := CompilePath(typeName)
					So(err, ShouldBeNil)
					a, getErrA := testData.Get(path)
					So(getErrA, ShouldBeNil)
					b, getErrB := copy.Get(path)
					So(getErrB, ShouldBeNil)
					// objects should have the same value
					So(a, ShouldEqual, b)
					// pointers should not be the same
					So(&a, ShouldNotPointTo, &b)
				}

				complexTypes := []string{"byte", "time", "null"}
				for _, typeName := range complexTypes {
					path, err := CompilePath(typeName)
					So(err, ShouldBeNil)
					a, getErrA := testData.Get(path)
					So(getErrA, ShouldBeNil)
					b, getErrB := copy.Get(path)
					So(getErrB, ShouldBeNil)
					// objects should have the same value
					So(a, ShouldResemble, b)
					// pointers should not be the same
					So(&a, ShouldNotPointTo, &b)
				}
			})
		})
	})

	Convey("Given a Map with a Bool value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get(MustCompilePath("bool"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.asBool()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, true)
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asString()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with an Int value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get(MustCompilePath("int"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should be accessible as an int", func() {
					val, typeErr := x.asInt()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 1)
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asBool()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asString()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a Float value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get(MustCompilePath("float"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should be accessible as a float", func() {
					val, typeErr := x.asFloat()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 0.1)
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asBool()
					So(err, ShouldNotBeNil)
					_, err = x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asString()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a String value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get(MustCompilePath("string"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.asString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "homhom")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asBool()
					So(err, ShouldNotBeNil)
					_, err = x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a Blob value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get(MustCompilePath("byte"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.asBlob()
					So(typeErr, ShouldBeNil)
					So(val, ShouldResemble, []byte("madmad"))
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asBool()
					So(err, ShouldNotBeNil)
					_, err = x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asString()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a Timestamp value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get(MustCompilePath("time"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.asTimestamp()
					So(typeErr, ShouldBeNil)
					const layout = "2006-01-02 15:04:00"
					expected, _ := time.Parse(layout, "2015-04-10 10:23:00")
					So(time.Time(val).Format(layout), ShouldEqual, expected.Format(layout))
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asBool()
					So(err, ShouldNotBeNil)
					_, err = x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asString()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with an Array value in it", t, func() {
		Convey("When accessing the value by key", func() {
			a, getErr := testData.Get(MustCompilePath("array"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original array", func() {
					e := Array([]Value{String("saysay"), String("mammam")})
					So(a, ShouldResemble, e)
				})
			})
		})
	})

	Convey("Given a Map with an Array value in it", t, func() {
		Convey("When accessing the first array element by key and index", func() {
			x, getErr := testData.Get(MustCompilePath("array[0]"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.asString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "saysay")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asBool()
					So(err, ShouldNotBeNil)
					_, err = x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When accessing the second array element by key and index", func() {
			x, getErr := testData.Get(MustCompilePath("array[1]"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.asString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "mammam")
				})
			})
		})

		Convey("When accessing an array element out of bounds", func() {
			_, getErr := testData.Get(MustCompilePath("array[2]"))
			Convey("Then the lookup should fail", func() {
				So(getErr, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a Map with a Map value in it", t, func() {
		Convey("When accessing the value by key", func() {
			m, getErr := testData.Get(MustCompilePath("map"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original map", func() {
					e := Map{
						"string": String("homhom2"),
					}
					So(m, ShouldResemble, e)
				})
			})
		})
	})

	Convey("Given a Map with a Map value in it", t, func() {
		Convey("When accessing a map element by nested key", func() {
			x, getErr := testData.Get(MustCompilePath("map.string"))
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.asString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "homhom2")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.asBool()
					So(err, ShouldNotBeNil)
					_, err = x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When accessing non-existing map element", func() {
			_, getErr := testData.Get(MustCompilePath("map.key"))
			Convey("Then the lookup should fail", func() {
				So(getErr, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a Map with a Null value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get(MustCompilePath("null"))
			Convey("Then the value should Null type object", func() {
				So(getErr, ShouldBeNil)
				So(x.Type(), ShouldEqual, TypeNull)
				Convey("and all type conversion should fail", func() {
					_, err := x.asInt()
					So(err, ShouldNotBeNil)
					_, err = x.asFloat()
					So(err, ShouldNotBeNil)
					_, err = x.asString()
					So(err, ShouldNotBeNil)
					_, err = x.asBlob()
					So(err, ShouldNotBeNil)
					_, err = x.asTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.asArray()
					So(err, ShouldNotBeNil)
					_, err = x.asMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestValueString(t *testing.T) {
	now := time.Now()

	testCases := map[string]([]convTestInput){
		"Null": {
			{"Null", Null{}, "null"},
		},
		"Bool": {
			{"true", Bool(true), "true"},
			{"false", Bool(false), "false"},
		},
		"Int": {
			{"positive", Int(2), "2"},
			{"negative", Int(-2), "-2"},
			{"zero", Int(0), "0"},
		},
		"Float": {
			{"positive", Float(3.14), "3.14"},
			{"negative", Float(-3.14), "-3.14"},
			{"zero", Float(0.0), "0"},
			// "NaN and Infinity regardless of sign are represented
			// as the String null." (ECMA-262)
			{"NaN", Float(math.NaN()), "null"},
		},
		"String": {
			{"empty", String(""), `""`},
			{"non-empty", String("hoge"), `"hoge"`},
			{"containing non-ASCII chars", String("日本語"), `"日本語"`},
		},
		"Blob": {
			{"empty", Blob(""), `""`}, // base64 of []
			{"nil", Blob(nil), "null"},
			{"non-empty", Blob("hoge"), `"aG9nZQ=="`}, // base64 of ['h','o','g','e']
		},
		"Timestamp": {
			{"zero", Timestamp(time.Time{}), `"0001-01-01T00:00:00Z"`},
			{"now", Timestamp(now), fmt.Sprintf(`"%s"`, now.Format(time.RFC3339Nano))},
		},
		"Array": []convTestInput{
			{"empty", Array{}, `[]`},
			{"non-empty", Array{Int(2), String("foo")}, `[2,"foo"]`},
		},
		"Map": []convTestInput{
			{"empty", Map{}, `{}`},
			// NB. this may fail randomly as map order is not guaranteed in Go
			{"non-empty", Map{"a": Int(2), "b": String("foo")}, `{"a":2,"b":"foo"}`},
		},
	}

	for valType, cases := range testCases {
		cases := cases
		Convey(fmt.Sprintf("Given a %s value", valType), t, func() {
			for _, testCase := range cases {
				tc := testCase
				Convey(fmt.Sprintf("When it is %s", tc.ValueDesc), func() {
					inVal := tc.Value
					exp := tc.Expected
					Convey(fmt.Sprintf("Then String returns %v", exp), func() {
						So(inVal.String(), ShouldResemble, exp)
					})
					Convey("And String returns the same as Marshal", func() {
						j, err := json.Marshal(inVal)
						So(err, ShouldBeNil)
						So([]byte(inVal.String()), ShouldResemble, j)
					})
				})
			}
		})
	}

	Convey("Given a nested value structure", t, func() {
		m := Map{
			"bool":    Bool(true),
			"int32":   Int(1),
			"int64":   Int(2),
			"float32": Float(float32(0.1)),
			"float64": Float(0.2),
			"string":  String("homhom"),
			"time":    Timestamp(now),
			"array": Array([]Value{Bool(true), Int(10), String("inarray"),
				Map{
					"mapinarray": String("arraymap"),
				}}),
			"map": Map{
				"map_a": String("a"),
				"map_b": Int(2),
			},
			"null": Null{},
		}
		Convey("String outputs a JSON representation", func() {
			s := m.String()
			// the order of the items in the Map may change, so we have to
			// use the following checks instead of just a single string comparison
			So(s, ShouldStartWith, `{"`)
			So(s, ShouldEndWith, `}`)
			So(s, ShouldContainSubstring, `"bool":true`)
			So(s, ShouldContainSubstring, `"int32":1`)
			So(s, ShouldContainSubstring, `"int64":2`)
			// the float32 below will actually be output as 0.10000000149011612
			So(s, ShouldContainSubstring, `"float32":0.1`)
			So(s, ShouldContainSubstring, `"float64":0.2`)
			So(s, ShouldContainSubstring, `"string":"homhom"`)
			So(s, ShouldContainSubstring, `"time":"`+now.Format(time.RFC3339Nano)+`"`)
			So(s, ShouldContainSubstring, `"array":[true,10,"inarray",{"mapinarray":"arraymap"}]`)
			So(s, ShouldContainSubstring, `"map":{"`)
			So(s, ShouldContainSubstring, `"map_a":"a"`)
			So(s, ShouldContainSubstring, `"map_b":2`)
			So(s, ShouldContainSubstring, `"null":null`)
		})
		Convey("Unmarshal does not result in an error", func() {
			v := map[string]interface{}{}
			err := json.Unmarshal([]byte(m.String()), &v)
			So(err, ShouldBeNil)
		})
	})
}

func TestNewValueFromSlice(t *testing.T) {
	Convey("Given NewValue function", t, func() {
		Convey("When passing a slice of integers", func() {
			v, err := NewValue([]int{1, 2, 3})
			So(err, ShouldBeNil)

			Convey("Then it should return an Array", func() {
				a, err := v.asArray()
				So(err, ShouldBeNil)
				So(len(a), ShouldEqual, 3)
				So(a[0], ShouldEqual, 1)
				So(a[1], ShouldEqual, 2)
				So(a[2], ShouldEqual, 3)
			})
		})

		Convey("When passing a slice of strings", func() {
			v, err := NewValue([]string{"a", "b", "c"})
			So(err, ShouldBeNil)

			Convey("Then it should return an Array", func() {
				a, err := v.asArray()
				So(err, ShouldBeNil)
				So(len(a), ShouldEqual, 3)
				So(a[0], ShouldEqual, "a")
				So(a[1], ShouldEqual, "b")
				So(a[2], ShouldEqual, "c")
			})
		})

		Convey("When passing a slice of slices", func() {
			v, err := NewValue([][]int{[]int{1}, []int{2, 3}})
			So(err, ShouldBeNil)

			Convey("Then it should return an Array", func() {
				a, err := v.asArray()
				So(err, ShouldBeNil)
				So(len(a), ShouldEqual, 2)
				So(a[0], ShouldResemble, Array{Int(1)})
				So(a[1], ShouldResemble, Array{Int(2), Int(3)})
			})
		})

		Convey("When passing a slice of maps", func() {
			v, err := NewValue([]map[string]interface{}{
				{"a": 1},
				{"b": "2", "c": 3.5},
			})
			So(err, ShouldBeNil)

			Convey("Then it should return an Array", func() {
				a, err := v.asArray()
				So(err, ShouldBeNil)
				So(len(a), ShouldEqual, 2)
				So(a[0], ShouldResemble, Map{"a": Int(1)})
				So(a[1], ShouldResemble, Map{"b": String("2"), "c": Float(3.5)})
			})
		})

		Convey("When passing a nil slice", func() {
			v, err := NewValue([]int(nil))
			So(err, ShouldBeNil)

			Convey("Then it should return an empty Array", func() {
				a, err := v.asArray()
				So(err, ShouldBeNil)
				So(len(a), ShouldEqual, 0)
			})
		})

		Convey("When passing a slice of inconvertible type", func() {
			_, err := NewValue([]struct{}{struct{}{}})

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

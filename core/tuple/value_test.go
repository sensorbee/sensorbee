package tuple

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ugorji/go/codec"
	"testing"
	"time"
)

// TestUnmarshalMsgpack tests that byte array encoded by msgpack
// is deserialized to Map object correctly
func TestUnmarshalMsgpack(t *testing.T) {
	Convey("Given a msgpack byte data", t, func() {
		var testMap = map[string]interface{}{
			"bool":    true,
			"int32":   int32(1),
			"int64":   int64(2),
			"float32": float32(0.1),
			"float64": float64(0.2),
			"string":  "homhom",
			"time":    time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC),
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
		codec.NewEncoderBytes(&testData, mh).Encode(testMap)
		fmt.Println(testData)
		Convey("When convert to Map object", func() {
			m, _ := UnmarshalMsgpack(testData)
			Convey("Then decode data should be match with Map data", func() {
				var expected = Map{
					"bool":    Bool(true),
					"int32":   Int(1),
					"int64":   Int(2),
					"float32": Float(float32(0.1)),
					"float64": Float(0.2),
					"string":  String("homhom"),
					"time":    Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)),
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

// TestMarshalMsgpack tests that Map object is serialized to
// byte array correctly.
func TestMarshalMsgpack(t *testing.T) {
	Convey("Given a Map object data", t, func() {
		var testMap = Map{
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
			"null": Null{},
			// TODO add Blob
		}
		Convey("When convert to []byte", func() {
			b, _ := MarshalMsgpack(testMap)
			Convey("Then encode data should be match with expected bytes", func() {
				var expected = map[string]interface{}{
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
					"null": nil,
				}
				var expectedBytes []byte
				codec.NewEncoderBytes(&expectedBytes, mh).Encode(expected)

				// it should compare b and expectedBytes, but byte array order is not
				// always correspond in converting map to bytes.
				var actualMap, expectedMap map[string]interface{}
				codec.NewDecoderBytes(expectedBytes, mh).Decode(&expectedMap)
				codec.NewDecoderBytes(b, mh).Decode(&actualMap)
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
			_, getErr := testData.Get("key")
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
					a, getErrA := testData.Get(typeName)
					So(getErrA, ShouldBeNil)
					b, getErrB := copy.Get(typeName)
					So(getErrB, ShouldBeNil)
					// objects should have the same value
					So(a, ShouldEqual, b)
					// pointers should not be the same
					So(&a, ShouldNotPointTo, &b)
				}

				complexTypes := []string{"byte", "time", "null"}
				for _, typeName := range complexTypes {
					a, getErrA := testData.Get(typeName)
					So(getErrA, ShouldBeNil)
					b, getErrB := copy.Get(typeName)
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
			x, getErr := testData.Get("bool")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.AsBool()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, true)
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsInt()
					So(err, ShouldNotBeNil)
					_, err = x.AsFloat()
					So(err, ShouldNotBeNil)
					_, err = x.AsString()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with an Int value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get("int")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should be accessible as an int", func() {
					val, typeErr := x.AsInt()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 1)
				})
				Convey("and it should be accessible as a float", func() {
					val, typeErr := x.AsFloat()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 1.0)
				})
				Convey("and it should be accessible as a string", func() {
					val, typeErr := x.AsString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "1")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsBool()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a Float value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get("float")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should be accessible as a float", func() {
					val, typeErr := x.AsFloat()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 0.1)
				})
				Convey("and it should be accessible as an int", func() {
					val, typeErr := x.AsInt()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 0)
				})
				Convey("and it should be accessible as a string", func() {
					val, typeErr := x.AsString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "0.1")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsBool()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a String value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get("string")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.AsString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "homhom")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsBool()
					So(err, ShouldNotBeNil)
					_, err = x.AsInt()
					So(err, ShouldNotBeNil)
					_, err = x.AsFloat()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a Blob value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get("byte")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.AsBlob()
					So(typeErr, ShouldBeNil)
					So(val, ShouldResemble, []byte("madmad"))
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsBool()
					So(err, ShouldNotBeNil)
					_, err = x.AsInt()
					So(err, ShouldNotBeNil)
					_, err = x.AsFloat()
					So(err, ShouldNotBeNil)
					_, err = x.AsString()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with a Timestamp value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get("time")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.AsTimestamp()
					So(typeErr, ShouldBeNil)
					const layout = "2006-01-02 15:04:00"
					expected, _ := time.Parse(layout, "2015-04-10 10:23:00")
					So(time.Time(val).Format(layout), ShouldEqual, expected.Format(layout))
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsBool()
					So(err, ShouldNotBeNil)
					_, err = x.AsInt()
					So(err, ShouldNotBeNil)
					_, err = x.AsFloat()
					So(err, ShouldNotBeNil)
					_, err = x.AsString()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Given a Map with an Array value in it", t, func() {
		Convey("When accessing the value by key", func() {
			a, getErr := testData.Get("array")
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
			x, getErr := testData.Get("array[0]")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.AsString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "saysay")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsBool()
					So(err, ShouldNotBeNil)
					_, err = x.AsInt()
					So(err, ShouldNotBeNil)
					_, err = x.AsFloat()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When accessing the second array element by key and index", func() {
			x, getErr := testData.Get("array[1]")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.AsString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "mammam")
				})
			})
		})

		Convey("When accessing an array element out of bounds", func() {
			_, getErr := testData.Get("array[2]")
			Convey("Then the lookup should fail", func() {
				So(getErr, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a Map with a Map value in it", t, func() {
		Convey("When accessing the value by key", func() {
			m, getErr := testData.Get("map")
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
			x, getErr := testData.Get("map.string")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.AsString()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "homhom2")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.AsBool()
					So(err, ShouldNotBeNil)
					_, err = x.AsInt()
					So(err, ShouldNotBeNil)
					_, err = x.AsFloat()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When accessing non-existing map element", func() {
			_, getErr := testData.Get("map/key")
			Convey("Then the lookup should fail", func() {
				So(getErr, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a Map with a Null value in it", t, func() {
		Convey("When accessing the value by key", func() {
			x, getErr := testData.Get("null")
			Convey("Then the value should Null type object", func() {
				So(getErr, ShouldBeNil)
				So(x.Type(), ShouldEqual, TypeNull)
				Convey("and all type conversion should fail", func() {
					_, err := x.AsInt()
					So(err, ShouldNotBeNil)
					_, err = x.AsFloat()
					So(err, ShouldNotBeNil)
					_, err = x.AsString()
					So(err, ShouldNotBeNil)
					_, err = x.AsBlob()
					So(err, ShouldNotBeNil)
					_, err = x.AsTimestamp()
					So(err, ShouldNotBeNil)
					_, err = x.AsArray()
					So(err, ShouldNotBeNil)
					_, err = x.AsMap()
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

package tuple

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

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
			"string": String("homhom"),
		},
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
					"array[0]", "map/string"}
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

				complexTypes := []string{"byte", "time"}
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
					val, typeErr := x.Bool()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, true)
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Int()
					So(err, ShouldNotBeNil)
					_, err = x.Float()
					So(err, ShouldNotBeNil)
					_, err = x.String()
					So(err, ShouldNotBeNil)
					_, err = x.Blob()
					So(err, ShouldNotBeNil)
					_, err = x.Timestamp()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
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
					val, typeErr := x.Int()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 1)
				})
				Convey("and it should be accessible as a float", func() {
					val, typeErr := x.Float()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 1.0)
				})
				Convey("and it should be accessible as a string", func() {
					val, typeErr := x.String()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "1")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Bool()
					So(err, ShouldNotBeNil)
					_, err = x.Blob()
					So(err, ShouldNotBeNil)
					_, err = x.Timestamp()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
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
					val, typeErr := x.Float()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 0.1)
				})
				Convey("and it should be accessible as an int", func() {
					val, typeErr := x.Int()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, 0)
				})
				Convey("and it should be accessible as a string", func() {
					val, typeErr := x.String()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "0.1")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Bool()
					So(err, ShouldNotBeNil)
					_, err = x.Blob()
					So(err, ShouldNotBeNil)
					_, err = x.Timestamp()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
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
					val, typeErr := x.String()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "homhom")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Bool()
					So(err, ShouldNotBeNil)
					_, err = x.Int()
					So(err, ShouldNotBeNil)
					_, err = x.Float()
					So(err, ShouldNotBeNil)
					_, err = x.Blob()
					So(err, ShouldNotBeNil)
					_, err = x.Timestamp()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
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
					val, typeErr := x.Blob()
					So(typeErr, ShouldBeNil)
					So(val, ShouldResemble, Blob([]byte("madmad")))
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Bool()
					So(err, ShouldNotBeNil)
					_, err = x.Int()
					So(err, ShouldNotBeNil)
					_, err = x.Float()
					So(err, ShouldNotBeNil)
					_, err = x.String()
					So(err, ShouldNotBeNil)
					_, err = x.Timestamp()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
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
					val, typeErr := x.Timestamp()
					So(typeErr, ShouldBeNil)
					const layout = "2006-01-02 15:04:00"
					expected, _ := time.Parse(layout, "2015-04-10 10:23:00")
					So(time.Time(val).Format(layout), ShouldEqual, expected.Format(layout))
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Bool()
					So(err, ShouldNotBeNil)
					_, err = x.Int()
					So(err, ShouldNotBeNil)
					_, err = x.Float()
					So(err, ShouldNotBeNil)
					_, err = x.String()
					So(err, ShouldNotBeNil)
					_, err = x.Blob()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
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
				Convey("and is should match the original array", func() {
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
					val, typeErr := x.String()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "saysay")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Bool()
					So(err, ShouldNotBeNil)
					_, err = x.Int()
					So(err, ShouldNotBeNil)
					_, err = x.Float()
					So(err, ShouldNotBeNil)
					_, err = x.Blob()
					So(err, ShouldNotBeNil)
					_, err = x.Timestamp()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When accessing the second array element by key and index", func() {
			x, getErr := testData.Get("array[1]")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.String()
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
				Convey("and is should match the original map", func() {
					e := Map{
						"string": String("homhom"),
					}
					So(m, ShouldResemble, e)
				})
			})
		})
	})

	Convey("Given a Map with a Map value in it", t, func() {
		Convey("When accessing a map element by nested key", func() {
			x, getErr := testData.Get("map/string")
			Convey("Then the value should exist", func() {
				So(getErr, ShouldBeNil)
				Convey("and it should match the original value", func() {
					val, typeErr := x.String()
					So(typeErr, ShouldBeNil)
					So(val, ShouldEqual, "homhom")
				})
				Convey("and other type conversions should fail", func() {
					_, err := x.Bool()
					So(err, ShouldNotBeNil)
					_, err = x.Int()
					So(err, ShouldNotBeNil)
					_, err = x.Float()
					So(err, ShouldNotBeNil)
					_, err = x.Blob()
					So(err, ShouldNotBeNil)
					_, err = x.Timestamp()
					So(err, ShouldNotBeNil)
					_, err = x.Array()
					So(err, ShouldNotBeNil)
					_, err = x.Map()
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
}

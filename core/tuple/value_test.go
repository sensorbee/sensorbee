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
	}

	Convey("Not found from test data", t, func() {
		_, err := testData.Get("key")
		So(err, ShouldNotBeNil)
	})

	Convey("Get bool value from test data", t, func() {
		x, err := testData.Get("bool")
		So(err, ShouldBeNil)

		Convey("Cast bool value", func() {
			b, err := x.Bool()
			So(err, ShouldBeNil)
			_, err = x.Int()
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

			Convey("The value should be equal to test data", func() {
				So(b, ShouldEqual, true)
			})
		})
	})

	Convey("Get int value from test data", t, func() {
		x, err := testData.Get("int")
		So(err, ShouldBeNil)

		Convey("Cast int value", func() {
			_, err = x.Bool()
			So(err, ShouldNotBeNil)
			i, err := x.Int()
			So(err, ShouldBeNil)
			f, err := x.Float()
			So(err, ShouldBeNil)
			s, err := x.String()
			So(err, ShouldBeNil)
			_, err = x.Blob()
			So(err, ShouldNotBeNil)
			_, err = x.Timestamp()
			So(err, ShouldNotBeNil)
			_, err = x.Array()
			So(err, ShouldNotBeNil)
			_, err = x.Map()
			So(err, ShouldNotBeNil)

			Convey("The value should be equal to test data", func() {
				So(i, ShouldEqual, 1)
				So(f, ShouldEqual, 1.0)
				So(s, ShouldEqual, "1")
			})
		})
	})

	Convey("Get float value from test data", t, func() {
		x, err := testData.Get("float")
		So(err, ShouldBeNil)

		Convey("Cast float value", func() {
			_, err = x.Bool()
			So(err, ShouldNotBeNil)
			i, err := x.Int()
			So(err, ShouldBeNil)
			f, err := x.Float()
			So(err, ShouldBeNil)
			s, err := x.String()
			So(err, ShouldBeNil)
			_, err = x.Blob()
			So(err, ShouldNotBeNil)
			_, err = x.Timestamp()
			So(err, ShouldNotBeNil)
			_, err = x.Array()
			So(err, ShouldNotBeNil)
			_, err = x.Map()
			So(err, ShouldNotBeNil)

			Convey("The value should be equal to test data", func() {
				So(i, ShouldEqual, 0)
				So(f, ShouldEqual, 0.1)
				So(s, ShouldEqual, "0.1")
			})
		})
	})

	Convey("Get string value from test data", t, func() {
		x, err := testData.Get("string")
		So(err, ShouldBeNil)

		Convey("Cast string value", func() {
			_, err = x.Bool()
			So(err, ShouldNotBeNil)
			_, err = x.Int()
			So(err, ShouldNotBeNil)
			_, err = x.Float()
			So(err, ShouldNotBeNil)
			s, err := x.String()
			So(err, ShouldBeNil)
			_, err = x.Blob()
			So(err, ShouldNotBeNil)
			_, err = x.Timestamp()
			So(err, ShouldNotBeNil)
			_, err = x.Array()
			So(err, ShouldNotBeNil)
			_, err = x.Map()
			So(err, ShouldNotBeNil)

			Convey("The value should be equal to test data", func() {
				So(s, ShouldEqual, "homhom")
			})
		})
	})

	Convey("Get byte value from test data", t, func() {
		x, err := testData.Get("byte")
		So(err, ShouldBeNil)

		Convey("Cast byte value", func() {
			_, err = x.Bool()
			So(err, ShouldNotBeNil)
			_, err = x.Int()
			So(err, ShouldNotBeNil)
			_, err = x.Float()
			So(err, ShouldNotBeNil)
			_, err = x.String()
			So(err, ShouldNotBeNil)
			b, err := x.Blob()
			So(err, ShouldBeNil)
			_, err = x.Timestamp()
			So(err, ShouldNotBeNil)
			_, err = x.Array()
			So(err, ShouldNotBeNil)
			_, err = x.Map()
			So(err, ShouldNotBeNil)

			Convey("The value should be equal to test data", func() {
				So([]byte(b), ShouldResemble, []byte("madmad"))
			})
		})
	})

	Convey("Get timestamp value from test data", t, func() {
		x, err := testData.Get("time")
		So(err, ShouldBeNil)

		Convey("Cast timestamp value", func() {
			_, err = x.Bool()
			So(err, ShouldNotBeNil)
			_, err = x.Int()
			So(err, ShouldNotBeNil)
			_, err = x.Float()
			So(err, ShouldNotBeNil)
			_, err = x.String()
			So(err, ShouldNotBeNil)
			_, err = x.Blob()
			So(err, ShouldNotBeNil)
			t, err := x.Timestamp()
			So(err, ShouldBeNil)
			_, err = x.Array()
			So(err, ShouldNotBeNil)
			_, err = x.Map()
			So(err, ShouldNotBeNil)

			Convey("The value should be equal to test data", func() {
				const layout = "2006-01-02 15:04:00"
				expected, _ := time.Parse(layout, "2015-04-10 10:23:00")
				So(time.Time(t).Format(layout), ShouldEqual, expected.Format(layout))
			})
		})
	})

	Convey("Get array value from test data", t, func() {
		_, err := testData.Get("array")
		So(err, ShouldNotBeNil) // TODO expected not occur error
		x, err := testData.Get("array[0]")

		So(err, ShouldBeNil)
		Convey("Cast string value", func() {
			s, err := x.String()
			So(err, ShouldBeNil)

			Convey("The value should be equal to test data", func() {
				So(s, ShouldEqual, "saysay")
			})
		})
	})

	SkipConvey("Get array from test data", t, func() {
		_, err := testData.Get("array")
		So(err, ShouldBeNil)
	})

}

func TestNestedValue(t *testing.T) {
	var testData = Map{
		"map": Map{
			"string": String("homhom"),
		},
	}

	Convey("Not found from nested map data", t, func() {
		_, err := testData.Get("map/key")
		So(err, ShouldNotBeNil)
	})

	Convey("Get nested value from test data", t, func() {
		_, err := testData.Get("map") // TODO expected not occur error
		So(err, ShouldNotBeNil)
		x, err := testData.Get("map/string")
		So(err, ShouldBeNil)

		Convey("Cast string value", func() {
			sx, err := x.String()
			So(err, ShouldBeNil)

			Convey("The value shoulld be equal to test data", func() {
				So(sx, ShouldEqual, "homhom")
			})
		})
	})

	SkipConvey("Get map from test data", t, func() {
		_, err := testData.Get("map")
		So(err, ShouldBeNil)
	})
}

package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"testing"
	"time"
)

func TestTuple(t *testing.T) {
	var testData = data.Map{
		"bool":   data.Bool(true),
		"int":    data.Int(1),
		"float":  data.Float(0.1),
		"string": data.String("homhom"),
		"byte":   data.Blob([]byte("madmad")),
		"time":   data.Timestamp(time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC)),
		"array":  data.Array([]data.Value{data.String("saysay"), data.String("mammam")}),
		"map": data.Map{
			"string": data.String("homhom"),
		},
	}
	tup := Tuple{
		Data:          testData,
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}

	dataShouldBeTheSame := func(t *Tuple) {
		simpleTypes := []string{"bool", "int", "float", "string",
			"array[0]", "map.string"}
		for _, typeName := range simpleTypes {
			a, getErrA := t.Data.Get(typeName)
			So(getErrA, ShouldBeNil)
			b, getErrB := t.Data.Get(typeName)
			So(getErrB, ShouldBeNil)
			// objects should have the same value
			So(a, ShouldEqual, b)
			// pointers should not be the same
			So(&a, ShouldNotPointTo, &b)
		}

		complexTypes := []string{"byte", "time"}
		for _, typeName := range complexTypes {
			a, getErrA := t.Data.Get(typeName)
			So(getErrA, ShouldBeNil)
			b, getErrB := t.Data.Get(typeName)
			So(getErrB, ShouldBeNil)
			// objects should have the same value
			So(a, ShouldResemble, b)
			// pointers should not be the same
			So(&a, ShouldNotPointTo, &b)
		}
	}

	Convey("Given a Tuple with values in it", t, func() {
		Convey("When deep-copying the Tuple", func() {
			copy := tup.Copy()

			Convey("Then tuple metadata should be the same", func() {
				So(tup.Timestamp, ShouldResemble, copy.Timestamp)
				So(&tup.Timestamp, ShouldNotPointTo, &copy.Timestamp)

				So(tup.ProcTimestamp, ShouldResemble, copy.ProcTimestamp)
				So(&tup.ProcTimestamp, ShouldNotPointTo, &copy.ProcTimestamp)

				So(tup.BatchID, ShouldResemble, copy.BatchID)
				So(&tup.BatchID, ShouldNotPointTo, &copy.BatchID)
			})

			Convey("Then all values should be the same", func() {
				dataShouldBeTheSame(copy)
			})
		})

		Convey("When creating a Tuple by NewTuple", func() {
			t := NewTuple(testData)

			Convey("Then tuple metadata should be initialized", func() {
				So(t.Timestamp, ShouldNotEqual, time.Time{})
				So(t.ProcTimestamp, ShouldNotEqual, time.Time{})

				So(t.InputName, ShouldBeEmpty)
				So(t.BatchID, ShouldEqual, 0)
				So(t.Trace, ShouldBeEmpty)
			})

			Convey("Then all values should be the same", func() {
				dataShouldBeTheSame(t)
			})
		})
	})
}

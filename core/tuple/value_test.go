package tuple

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestValue(t *testing.T) {
	var testData = Map{
		"int":    Int(1),
		"float":  Float(0.1),
		"string": String("homhom"),
		// TODO add Blob, Timestamp, Array, Map
	}

	// Only pass t into top-level Convey calls
	Convey("Get int value from test data", t, func() {
		x, err := testData.Get("int")
		if err != nil {
			t.Fatal(err)
		}

		Convey("Cast int value", func() {
			i, err2 := x.Int()
			if err2 != nil {
				t.Fatal(err2)
			}

			Convey("The value should be equal to test data", func() {
				So(i, ShouldEqual, 1)
			})
		})
	})
}

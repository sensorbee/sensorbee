package tuple

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestMapMarshalJSON(t *testing.T) {
	type testStruct struct {
		M Map `json:"map"`
	}

	Convey("Given a struct embedding a Map", t, func() {
		now := time.Now()
		s := &testStruct{
			M: Map{
				"c": Int(2),
				"d": Timestamp(now),
			},
		}

		Convey("When marshaling it to a JSON", func() {
			js, err := json.Marshal(s)
			So(err, ShouldBeNil)

			Convey("Then it should contain an int field", func() {
				So(string(js), ShouldContainSubstring, `"c":2`)
			})

			Convey("Then it should contain a timestamp field", func() {
				So(string(js), ShouldContainSubstring, `"d":"`+now.Format(time.RFC3339Nano)+`"`)
			})
		})

		Convey("When unmarshaling it from a JSON", func() {
			js := `
			{
				"map": {
					"a": 1,
					"b": "hoge"
				}
			}
			`
			So(json.Unmarshal([]byte(js), s), ShouldBeNil)

			Convey("Then the struct should a correct int value", func() {
				v, err := ToInt(s.M["a"])
				So(err, ShouldBeNil)
				So(v, ShouldEqual, 1)
			})

			Convey("Then the struct should a correct string value", func() {
				v, err := ToString(s.M["b"])
				So(err, ShouldBeNil)
				So(v, ShouldEqual, "hoge")
			})
		})
	})
}

package tuple

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestArrayMarshalJSON(t *testing.T) {
	type testStruct struct {
		A Array `json:"array"`
	}

	Convey("Given a struct embedding an Array", t, func() {
		now := time.Now()
		s := &testStruct{
			A: Array{
				Int(2),
				Timestamp(now),
			},
		}

		Convey("When marshaling it to a JSON", func() {
			js, err := json.Marshal(s)
			So(err, ShouldBeNil)

			Convey("Then it should contain fields", func() {
				So(string(js), ShouldContainSubstring, `[2,"`+now.Format(time.RFC3339Nano)+`"]`)
			})
		})

		Convey("When unmarshaling it from a JSON", func() {
			js := `
			{
				"array": [
					1,
					"hoge"
				]
			}
			`
			So(json.Unmarshal([]byte(js), s), ShouldBeNil)

			Convey("Then the struct should a correct int value", func() {
				v, err := ToInt(s.A[0])
				So(err, ShouldBeNil)
				So(v, ShouldEqual, 1)
			})

			Convey("Then the struct should a correct string value", func() {
				v, err := ToString(s.A[1])
				So(err, ShouldBeNil)
				So(v, ShouldEqual, "hoge")
			})
		})
	})
}

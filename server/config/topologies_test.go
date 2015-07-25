package config

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTopologies(t *testing.T) {
	Convey("Given a JSON config for logging section", t, func() {
		Convey("When the config is valid", func() {
			ts, err := NewTopologies(toMap(`{"test1":{},"test2":{"bql_file":"/path/to/hoge.bql"}}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters", func() {
				So(ts["test1"].Name, ShouldEqual, "test1")
				So(ts["test1"].BQLFile, ShouldEqual, "")
				So(ts["test2"].Name, ShouldEqual, "test2")
				So(ts["test2"].BQLFile, ShouldEqual, "/path/to/hoge.bql")
			})
		})

		Convey("When the config only has required parameters", func() {
			// no required parameter at the moment
			ts, err := NewTopologies(toMap(`{}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters and default values", func() {
				So(ts, ShouldBeEmpty)
			})
		})

		Convey("When the config has an undefined field", func() {
			_, err := NewTopologies(toMap(`{"test":{"bql_path":"/path/to/hoge.bql"}}`))

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When validating bql_file", func() {
			for _, b := range []string{"a", "test.bql", "/path/to/hoge.bql"} {
				Convey(fmt.Sprint("Then it should accept ", b), func() {
					ts, err := NewTopologies(toMap(fmt.Sprintf(`{"test":{"bql_file":"%v"}}`, b)))
					So(err, ShouldBeNil)
					So(ts["test"].BQLFile, ShouldEqual, b)
				})
			}

			for _, b := range [][]interface{}{{"empty", `""`}, {"invalid type", 1}} {
				Convey(fmt.Sprintf("Then it should reject %v value", b[0]), func() {
					_, err := NewTopologies(toMap(fmt.Sprintf(`{"test":{"bql_file":%v}}`, b[1])))
					So(err, ShouldNotBeNil)
					fmt.Println(err)
				})
			}
		})
	})
}

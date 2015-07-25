package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConfig(t *testing.T) {
	Convey("Given a JSON config", t, func() {
		base := toMap(`{
	"network": {
		"listen_on": ":12345"
	},
	"topologies": {
		"test1": {
		},
		"test2": {
			"bql_file": "/path/to/hoge.bql"
		}
	},
	"logging": {
		"target": "stdout"
	}
}`)
		Convey("When the config is valid", func() {
			c, err := New(base)
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters", func() {
				So(c.Network.ListenOn, ShouldEqual, ":12345")
				So(c.Topologies["test1"].Name, ShouldEqual, "test1")
				So(c.Topologies["test2"].BQLFile, ShouldEqual, "/path/to/hoge.bql")
				So(c.Logging.Target, ShouldEqual, "stdout")
			})
		})

		// Because detailed cases are covered in other test, this test case
		// only check additional properties.

		Convey("When the config has an undefined field", func() {
			base["loggin"] = base["logging"]
			delete(base, "logging")
			_, err := New(base)

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

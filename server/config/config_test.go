package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
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

func TestConfigToMap(t *testing.T) {
	Convey("Given a server config", t, func() {
		c := Config{
			Network: &Network{
				ListenOn: "12345",
			},
			Topologies: Topologies{
				"t1": &Topology{
					BQLFile: "t1.bql",
				},
				"t2": &Topology{
					BQLFile: "t2.bql",
				},
			},
			Storage: &Storage{
				UDS: UDSStorage{
					Type: "fs",
					Params: data.Map{
						"dir": data.String("uds"),
					},
				},
			},
			Logging: &Logging{
				Target:                 "stderr",
				MinLogLevel:            "info",
				LogDroppedTuples:       true,
				SummarizeDroppedTuples: true,
			},
		}
		Convey("When convert to data.Map", func() {
			ac := c.ToMap()
			Convey("Then map should be equal as the config", func() {
				ex := data.Map{
					"network": data.Map{
						"listen_on": data.String("12345"),
					},
					"topologies": data.Array{
						data.Map{
							"name":     data.String("t1"),
							"bql_file": data.String("t1.bql"),
						}, data.Map{
							"name":     data.String("t2"),
							"bql_file": data.String("t2.bql"),
						},
					},
					"storage": data.Map{
						"uds": data.Map{
							"type": data.String("fs"),
							"params": data.Map{
								"dir": data.String("uds"),
							},
						},
					},
					"logging": data.Map{
						"target":                   data.String("stderr"),
						"min_log_level":            data.String("info"),
						"log_dropped_tuples":       data.True,
						"summarize_dropped_tuples": data.True,
					},
				}
				So(ac, ShouldResemble, ex)
			})
		})
	})
}

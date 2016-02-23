package config

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNetwork(t *testing.T) {
	Convey("Given a JSON config for network section", t, func() {
		Convey("When the config is valid", func() {
			n, err := NewNetwork(toMap(`{"listen_on":":12345"}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters", func() {
				So(n.ListenOn, ShouldEqual, ":12345")
			})
		})

		Convey("When the config only has required parameters", func() {
			// no required parameter at the moment
			n, err := NewNetwork(toMap(`{}`))

			Convey("Then it should have given parameters and default values", func() {
				So(err, ShouldBeNil)
				So(n.ListenOn, ShouldEqual, fmt.Sprintf(":%d", DefaultPort))
			})
		})

		Convey("When the config has an undefined field", func() {
			_, err := NewNetwork(toMap(`{"listen_on":":12345","listenon":":12345"}`))

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When validating listen_on", func() {
			for _, addr := range []string{fmt.Sprintf("127.0.0.1:%d", DefaultPort), fmt.Sprintf("localhost:%d", DefaultPort), fmt.Sprintf(":%d", DefaultPort)} {
				Convey(fmt.Sprint("Then it should accept ", addr), func() {
					n, err := NewNetwork(toMap(fmt.Sprintf(`{"listen_on":"%v"}`, addr)))
					So(err, ShouldBeNil)
					So(n.ListenOn, ShouldEqual, addr)
				})
			}

			for _, lv := range [][]interface{}{{"empty addr", `""`},
				{"no port", `":"`},
				{"no :", fmt.Sprintf(`"localhost%d"`, DefaultPort)},
				{"invalid type", 1}} {
				Convey(fmt.Sprintf("Then it should reject %v", lv[0]), func() {
					_, err := NewNetwork(toMap(fmt.Sprintf(`{"listen_on":%v}`, lv[1])))
					So(err, ShouldNotBeNil)
				})
			}
		})
	})
}

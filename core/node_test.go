package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestValidateNodeName(t *testing.T) {
	Convey("Given a node name validation function", t, func() {
		cases := []struct {
			title   string
			name    string
			success bool
		}{
			{"an empty string", "", false},
			{"one alphabet name", "a", true},
			{"one digit name", "0", false},
			{"one underscore name", "_", false},
			{"a name starting from a digit", "0test", false},
			{"a name starting from an underscore", "_test", false},
			{"a name only containing alphabet", "test", true},
			{"a name containing alphabet-digit", "t1e2s3t4", true},
			{"a name containing alphabet-digit-underscore", "t1e2_s3t4", true},
			{"a name containing an invalid letter", "da-me", false},
			{"a name has maximum length", strings.Repeat("a", 127), true},
			{"a name longer than the maximum length", strings.Repeat("a", 128), false},
		}

		for _, c := range cases {
			c := c

			Convey("When validating "+c.title, func() {
				if c.success {
					Convey("Then it should succeed", func() {
						So(ValidateNodeName(c.name), ShouldBeNil)
					})
				} else {
					Convey("Then it should fail", func() {
						So(ValidateNodeName(c.name), ShouldNotBeNil)
					})
				}
			})
		}
	})
}

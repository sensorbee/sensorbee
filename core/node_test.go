package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestValidateSymbol(t *testing.T) {
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
			{"a reserved word", "UNTIL", false},
			{"a reserved word with different capitalization", "vaLIdaTE", false},
		}

		for _, c := range cases {
			c := c

			Convey("When validating "+c.title, func() {
				if c.success {
					Convey("Then it should succeed", func() {
						So(ValidateSymbol(c.name), ShouldBeNil)
					})
				} else {
					Convey("Then it should fail", func() {
						So(ValidateSymbol(c.name), ShouldNotBeNil)
					})
				}
			})
		}
	})
}

func TestValidateCapacity(t *testing.T) {
	Convey("Given validateCapacity function", t, func() {
		Convey("When passing a valid value to it", func() {
			Convey("Then it should accept 0", func() {
				So(validateCapacity(0), ShouldBeNil)
			})

			Convey("Then it should accept a maximum value", func() {
				So(validateCapacity(MaxCapacity), ShouldBeNil)
			})
		})

		Convey("When passing a too large value", func() {
			Convey("Then it should fail", func() {
				So(validateCapacity(MaxCapacity+1), ShouldNotBeNil)
			})
		})

		Convey("When passing a negative value", func() {
			Convey("Then it should fail", func() {
				So(validateCapacity(-1), ShouldNotBeNil)
			})
		})
	})
}

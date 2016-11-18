package config

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLogging(t *testing.T) {
	Convey("Given a JSON config for logging section", t, func() {
		Convey("When the config is valid", func() {
			l, err := NewLogging(toMap(`{"target":"stdout","min_log_level":"error","log_dropped_tuples":true,"log_destinationless_tuples":true,"summarize_dropped_tuples":true}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters", func() {
				So(l.Target, ShouldEqual, "stdout")
				So(l.MinLogLevel, ShouldEqual, "error")
				So(l.LogDroppedTuples, ShouldBeTrue)
				So(l.SummarizeDroppedTuples, ShouldBeTrue)
			})
		})

		Convey("When the config only has required parameters", func() {
			// no required parameter at the moment
			l, err := NewLogging(toMap(`{}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters and default values", func() {
				So(l.Target, ShouldEqual, "stderr")
				So(l.MinLogLevel, ShouldEqual, "info")
				So(l.LogDroppedTuples, ShouldBeFalse)
				So(l.SummarizeDroppedTuples, ShouldBeFalse)
			})
		})

		Convey("When the config has an undefined field", func() {
			_, err := NewLogging(toMap(`{"target":"stdout","max_log_level":"info"}`))

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When validating target parameter", func() {
			Convey("Then it should accept stdout", func() {
				l, err := NewLogging(toMap(`{"target":"stdout"}`))
				So(err, ShouldBeNil)
				w, err := l.CreateWriter()
				So(err, ShouldBeNil)
				So(w.(*nopCloser).w, ShouldPointTo, os.Stdout)
			})

			Convey("Then it should accept stderr", func() {
				l, err := NewLogging(toMap(`{"target":"stderr"}`))
				So(err, ShouldBeNil)
				w, err := l.CreateWriter()
				So(err, ShouldBeNil)
				So(w.(*nopCloser).w, ShouldPointTo, os.Stderr)
			})

			// Testing Target having file isn't currently done for SSDs.
		})

		Convey("When validating min_log_level", func() {
			for _, lv := range []string{"debug", "info", "warn", "warning", "error", "fatal"} {
				Convey(fmt.Sprint("Then it should accept ", lv), func() {
					l, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","min_log_level":"%v"}`, lv)))
					So(err, ShouldBeNil)
					So(l.MinLogLevel, ShouldEqual, lv)
				})
			}

			for _, lv := range [][]interface{}{{"empty", ""}, {"invalid", "ainfob"}, {"invalid type", 1}} {
				Convey(fmt.Sprintf("Then it should reject %v value", lv[0]), func() {
					_, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","min_log_level":"%v"}`, lv[1])))
					So(err, ShouldNotBeNil)
				})
			}
		})

		Convey("When validating log_dropped_tuples", func() {
			for _, v := range []bool{true, false} {
				Convey(fmt.Sprint("Then it should accept ", v), func() {
					l, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","log_dropped_tuples":%v}`, v)))
					So(err, ShouldBeNil)
					So(l.LogDroppedTuples, ShouldEqual, v)
				})
			}

			for _, v := range [][]interface{}{{"an integer", 1}, {"a string", `"true"`}} {
				Convey(fmt.Sprintf("Then it should reject %v value", v[0]), func() {
					_, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","log_dropped_tuples":%v}`, v[1])))
					So(err, ShouldNotBeNil)
				})
			}
		})

		Convey("When validating log_destinationless_tuples", func() {
			for _, v := range []bool{true, false} {
				Convey(fmt.Sprint("Then it should accept ", v), func() {
					l, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","log_destinationless_tuples":%v}`, v)))
					So(err, ShouldBeNil)
					So(l.LogDestinationlessTuples, ShouldEqual, v)
				})
			}

			for _, v := range [][]interface{}{{"an integer", 1}, {"a string", `"true"`}} {
				Convey(fmt.Sprintf("Then it should reject %v value", v[0]), func() {
					_, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","log_destinationless_tuples":%v}`, v[1])))
					So(err, ShouldNotBeNil)
				})
			}
		})

		Convey("When validating summarize_dropped_tuples", func() {
			for _, v := range []bool{true, false} {
				Convey(fmt.Sprint("Then it should accept ", v), func() {
					l, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","summarize_dropped_tuples":%v}`, v)))
					So(err, ShouldBeNil)
					So(l.SummarizeDroppedTuples, ShouldEqual, v)
				})
			}

			for _, v := range [][]interface{}{{"an integer", 1}, {"a string", `"true"`}} {
				Convey(fmt.Sprintf("Then it should reject %v value", v[0]), func() {
					_, err := NewLogging(toMap(fmt.Sprintf(`{"target":"stderr","summarize_dropped_tuples":%v}`, v[1])))
					So(err, ShouldNotBeNil)
				})
			}
		})
	})
}

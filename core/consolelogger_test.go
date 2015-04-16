package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
)

func TestConsoleLogger(t *testing.T) {
	Convey("Given target log manager", t, func() {
		logger := ConsoleLogManager{}

		Convey("When given several log type messages", func() {
			dl := "debug log %v"
			el := "error log %v"

			Convey("Logger should write　message on console", func() {
				logger.Log(DEBUG, dl, "aaa")
				logger.Log(ERROR, el, "bbb")
			})
		})

		Convey("When given dropped tuple", func() {
			t := &tuple.Tuple{BatchID: 1}

			Convey("Logger should write tuple dump on console", func() {
				logger.DroppedTuple(t, "for debug")
			})
		})
	})
}

package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConsoleLogger(t *testing.T) {
	Convey("Given target log manager", t, func() {
		logger := ConsoleLogManager{}

		Convey("Given several log type messages", func() {
			dl := "debug log %v"
			el := "error log %v"

			Convey("Logger should write console", func() {
				logger.Log(DEBUG, dl, "aaa")
				logger.Log(ERROR, el, "bbb")
			})
		})
	})
}

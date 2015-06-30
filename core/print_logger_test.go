package core

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPrintLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewPrintLogger(buf)

	Convey("Given a print log manager bound to a Writer", t, func() {

		Convey("When given several log type messages", func() {
			dl := "debug log %v"
			el := "error log %v"

			Convey("Logger should write messages to Writer", func() {
				logger.Log(Debug, dl, "aaa")
				line1, err := buf.ReadString('\n')
				So(err, ShouldBeNil)
				So(line1, ShouldEndWith, "[DEBUG  ] debug log aaa\n")

				logger.Log(Error, el, "bbb")
				line2, err := buf.ReadString('\n')
				So(err, ShouldBeNil)
				So(line2, ShouldEndWith, "[ERROR  ] error log bbb\n")
			})
		})

		Convey("When given dropped tuple", func() {
			t := &Tuple{BatchID: 1}

			Convey("Logger should write tuple dump to Writer", func() {
				logger.DroppedTuple(t, "for debug")
				line, err := buf.ReadString('\n')
				So(err, ShouldBeNil)
				So(line, ShouldEndWith, "[DROPPED] Tuple Batch ID is 1, for debug\n")
			})
		})
	})
}

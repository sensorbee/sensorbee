package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultTopology(t *testing.T) {
	Convey("When creating a default topology builder", t, func() {
		var tb TopologyBuilder = NewDefaultTopologyBuilder()
		So(tb, ShouldNotBeNil)
	})

	Convey("Given a default topology builder", t, func() {
		tb := NewDefaultTopologyBuilder()
		s := &DefaultSource{}
		b := &DefaultBox{}
		var err DeclarerError

		Convey("when using a source name twice", func() {
			err = tb.AddSource("mySource", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("mySource", s)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a source name with a box name", func() {
			err = tb.AddBox("myBoxSource", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("myBoxSource", s)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name twice", func() {
			err = tb.AddBox("myBox", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("myBox", b)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

	})
}

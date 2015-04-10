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
		si := &DefaultSink{}
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
			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a source name with a sink name", func() {
			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("someName", s)
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

		Convey("when using a box name with a source name", func() {
			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name with a sink name", func() {
			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name twice", func() {
			err = tb.AddSink("mySink", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("mySink", si)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name with a source name", func() {
			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name with a box name", func() {
			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})
	})

	Convey("Given a default topology builder with a source", t, func() {
		tb := NewDefaultTopologyBuilder()
		s := &DefaultSource{}
		tb.AddSource("aSource", s)
		b := &DefaultBox{}
		tb.AddBox("aBox", b)
		var err DeclarerError

		Convey("when a new box references a non-existing item", func() {
			err = tb.AddBox("otherBox", b).
				Input("something", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing source", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing box", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references multiple items", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing source twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource", nil).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing box twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil).
				Input("aBox", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})
	})
}

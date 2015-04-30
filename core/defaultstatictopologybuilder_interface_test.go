package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestDefaultTopologyBuilderInterface tests that checks for names
// and dependencies are done correctly when building a topology.
func TestDefaultTopologyBuilderInterface(t *testing.T) {
	Convey("When creating a default topology builder", t, func() {
		var tb StaticTopologyBuilder = NewDefaultStaticTopologyBuilder()
		So(tb, ShouldNotBeNil)
	})

	Convey("Given a default topology builder", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DoesNothingSource{}
		b := &DoesNothingBox{}
		si := &DoesNothingSink{}
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
		tb := NewDefaultStaticTopologyBuilder()
		s := &DoesNothingSource{}
		tb.AddSource("aSource", s)
		b := &DoesNothingBox{}
		tb.AddBox("aBox", b)
		si := &DoesNothingSink{}
		var err DeclarerError

		Convey("when a new box references a non-existing item", func() {
			err = tb.AddBox("otherBox", b).
				Input("something")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing source", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing box", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references multiple items", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing source twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing box twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox").
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references a non-existing item", func() {
			err = tb.AddSink("aSink", si).
				Input("something")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references an existing source", func() {
			err = tb.AddSink("aSink", si).
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references an existing box", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references multiple items", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references an existing source twice", func() {
			err = tb.AddSink("aSink", si).
				Input("aSource").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references an existing box twice", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox").
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})
	})
}

// TestMultipleBuild tests that checks default static topology builder
// can build only once.
func TestMultipleBuild(t *testing.T) {
	Convey("Given basic topology builder called build() once", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		tp, _ := tb.Build()
		So(tp, ShouldNotBeNil)
		Convey("When add source", func() {
			sd := tb.AddSource("src", &DoesNothingSource{})
			Convey("Then it should occur non-buildable error", func() {
				So(sd.Err(), ShouldNotBeNil)
			})
		})
		Convey("When add box", func() {
			bd := tb.AddBox("box", &DoesNothingBox{})
			Convey("Then it should occur non-buildable error", func() {
				So(bd.Err(), ShouldNotBeNil)
			})
		})
		Convey("When add sink", func() {
			sd := tb.AddSink("si", &DoesNothingSink{})
			Convey("Then it should occur non-buildable error", func() {
				So(sd.Err(), ShouldNotBeNil)
			})
		})
		Convey("When build topology once again", func() {
			_, err := tb.Build()
			Convey("Then it should occur non-buildable error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

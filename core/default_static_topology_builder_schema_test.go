package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestDefaultStaticTopologyBuilderInputNamesChecks tests that checks for matching
// named inputs are done correctly when building a topology.
func TestDefaultStaticTopologyBuilderInputNamesChecks(t *testing.T) {
	Convey("Given a default topology builder", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DoesNothingSource{}
		tb.AddSource("source", s)

		Convey("When using a box with no input name constraint", func() {
			b := &DoesNothingBoxWithInputNames{}

			Convey("Then adding an unnamed input should succeed", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		Convey("When using a box with empty input names", func() {
			b := &DoesNothingBoxWithInputNames{InNames: []string{}}

			Convey("Then adding an unnamed input should succeed", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		// Note that the following test reveals a questionable behavior: If
		// a Box declares just one input stream, there is apparently no need
		// to tell apart different streams; so what exactly is the value of
		// requiring a certain name?
		Convey("When using a box with a named input 'hoge'", func() {
			b := &DoesNothingBoxWithInputNames{InNames: []string{"hoge"}}

			Convey("Then adding an unnamed input should fail", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for '*' should fail", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'foo' should fail", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "foo")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		Convey("When using a box with a named inputs 'hoge' and 'fuga'", func() {
			b := &DoesNothingBoxWithInputNames{InNames: []string{"hoge", "fuga"}}

			Convey("Then adding an unnamed input should fail", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for '*' should fail", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'foo' should fail", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "foo")
				So(bdecl.Err(), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'fuga' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "fuga")
				So(bdecl.Err(), ShouldBeNil)
			})
		})
	})
}

/**************************************************/

// DoesNothingBoxWithInputNames is basically like DoNothingBox, but it is
// possible to add a custom input name that will be returned by InputNames.
// This way it is possible to simulate a Box with input requirements
type DoesNothingBoxWithInputNames struct {
	InNames []string
	DoesNothingBox
}

var _ NamedInputBox = &DoesNothingBoxWithInputNames{}

func (b *DoesNothingBoxWithInputNames) InputNames() []string {
	return b.InNames
}

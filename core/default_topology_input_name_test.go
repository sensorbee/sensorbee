package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestDefaultTopologyInputNamesChecks tests that checks for matching
// named inputs are done correctly when building a topology.
func TestDefaultTopologyInputNamesChecks(t *testing.T) {
	Convey("Given a default topology builder", t, func() {
		t := NewDefaultTopology(NewContext(nil), "dt1")
		Reset(func() {
			t.Stop()
		})

		s := NewTupleEmitterSource(freshTuples())
		t.AddSource("source", s, &SourceConfig{
			PausedOnStartup: true,
		})

		Convey("When using a box with no input name constraint", func() {
			b := &DoesNothingBoxWithInputNames{}
			bn, err := t.AddBox("box", b, nil)
			So(err, ShouldBeNil)

			Convey("Then adding an unnamed input should succeed", func() {
				So(bn.Input("source", nil), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "*"}), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "hoge"}), ShouldBeNil)
			})
		})

		Convey("When using a box with empty input names", func() {
			b := &DoesNothingBoxWithInputNames{InNames: []string{}}
			bn, err := t.AddBox("box", b, nil)
			So(err, ShouldBeNil)

			Convey("Then adding an unnamed input should succeed", func() {
				So(bn.Input("source", nil), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "*"}), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "hoge"}), ShouldBeNil)
			})
		})

		// Note that the following test reveals a questionable behavior: If
		// a Box declares just one input stream, there is apparently no need
		// to tell apart different streams; so what exactly is the value of
		// requiring a certain name?
		Convey("When using a box with a named input 'hoge'", func() {
			b := &DoesNothingBoxWithInputNames{InNames: []string{"hoge"}}
			bn, err := t.AddBox("box", b, nil)
			So(err, ShouldBeNil)

			Convey("Then adding an unnamed input should fail", func() {
				So(bn.Input("source", nil), ShouldNotBeNil)
			})

			Convey("Then adding a named input for '*' should fail", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "*"}), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'foo' should fail", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "foo"}), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "hoge"}), ShouldBeNil)
			})
		})

		Convey("When using a box with a named inputs 'hoge' and 'fuga'", func() {
			b := &DoesNothingBoxWithInputNames{InNames: []string{"hoge", "fuga"}}
			bn, err := t.AddBox("box", b, nil)
			So(err, ShouldBeNil)

			Convey("Then adding an unnamed input should fail", func() {
				So(bn.Input("source", nil), ShouldNotBeNil)
			})

			Convey("Then adding a named input for '*' should fail", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "*"}), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'foo' should fail", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "foo"}), ShouldNotBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "hoge"}), ShouldBeNil)
			})

			Convey("Then adding a named input for 'fuga' should succeed", func() {
				So(bn.Input("source", &BoxInputConfig{InputName: "fuga"}), ShouldBeNil)
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

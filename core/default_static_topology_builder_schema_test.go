package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestDefaultStaticTopologyBuilderSchemaChecks tests that checks for matching
// input and output schema are done correctly when building a topology.
func TestDefaultStaticTopologyBuilderSchemaChecks(t *testing.T) {
	Convey("Given a default topology builder", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DoesNothingSource{}
		tb.AddSource("source", s)

		Convey("When using a box with nil input constraint", func() {
			// A box with InputConstraint() == nil should allow any and all input
			b := &DoesNothingBoxWithSchema{}

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

		Convey("When using a box with {'*' => nil} input constraint", func() {
			// A box with '*' => nil should allow any and all input
			b := &DoesNothingBoxWithSchema{
				InSchema: SchemaSet{"*": nil}}

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

		SkipConvey("When using a box with {'*' => non-nil} input constraint", func() {
			// A box with '*' => aSchema should only allow input
			// if it matches aSchema
		})

		// Note that the following test reveals a questionable behavior: If
		// a Box declares just one input stream, there is apparently no need
		// to tell apart different streams; so what exactly is the value of
		// requiring a certain name?
		Convey("When using a box with {'hoge' => nil} input constraint", func() {
			// A box with 'hoge' => nil should only allow input
			// if it comes from an input stream called 'hoge'
			b := &DoesNothingBoxWithSchema{
				InSchema: SchemaSet{"hoge": nil}}

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

		SkipConvey("When using a box with {'hoge' => non-nil} input constraint", func() {
			// A box with 'hoge' => aSchema should only allow input
			// if it comes from an input stream called 'hoge' and
			// matches aSchema
		})

		Convey("When using a box with {'hoge' => nil, '*' => nil} input constraint", func() {
			// A box with 'hoge' => nil, *' => nil should allow any and
			// all input
			b := &DoesNothingBoxWithSchema{
				InSchema: SchemaSet{"hoge": nil, "*": nil}}

			Convey("Then adding an unnamed input should succeed", func() {
				bdecl := tb.AddBox("box", b).Input("source")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for '*' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "*")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'foo' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "foo")
				So(bdecl.Err(), ShouldBeNil)
			})

			Convey("Then adding a named input for 'hoge' should succeed", func() {
				bdecl := tb.AddBox("box", b).NamedInput("source", "hoge")
				So(bdecl.Err(), ShouldBeNil)
			})
		})

		SkipConvey("When using a box with {'hoge' => non-nil, '*' => non-nil} input constraint", func() {
			// A box with 'hoge' => aSchema, '*' => otherSchema should
			// allow input from arbitrarily named input streams, as long as
			// they match the corresponding schema
		})

		SkipConvey("When using a box with {'hoge' => non-nil, 'foo' => non-nil} input constraint", func() {
			// A box with 'hoge' => aSchema, 'foo' => otherSchema should
			// only allow input from input streams called 'hoge' or 'foo'
			// and matching the corresponding schema
		})
	})
}

/**************************************************/

// DoesNothingBoxWithSchema is basically like DoNothingBox, but it is
// possible to add a schema that will be returned by InputConstraints.
// This way it is possible to simulate a Box with input requirements
type DoesNothingBoxWithSchema struct {
	InSchema SchemaSet
	DoesNothingBox
}

var _ SchemafulBox = &DoesNothingBoxWithSchema{}

func (b *DoesNothingBoxWithSchema) InputSchema() SchemaSet {
	return b.InSchema
}

func (b *DoesNothingBoxWithSchema) OutputSchema(s SchemaSet) (*Schema, error) {
	return nil, nil
}

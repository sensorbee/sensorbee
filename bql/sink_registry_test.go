package bql

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"testing"
)

func TestEmptyDefaultSinkCreatorRegistry(t *testing.T) {
	Convey("Given an empty default Sink registry", t, func() {
		r := NewDefaultSinkCreatorRegistry()

		Convey("When adding a creator function", func() {
			err := r.Register("test_sink", SinkCreatorFunc(createCollectorSink))

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When looking up a nonexistent creator", func() {
			_, err := r.Lookup("test_sink")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When retrieving a list of creators", func() {
			m, err := r.List()

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And the list should be empty", func() {
					So(m, ShouldBeEmpty)
				})
			})
		})

		Convey("When unregistering a nonexistent creator", func() {
			err := r.Unregister("test_sink")

			Convey("Then it shouldn't fail", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestDefaultSinkCreatorRegistry(t *testing.T) {
	ctx := core.NewContext(nil)

	Convey("Given an default Sink registry having two types", t, func() {
		r := NewDefaultSinkCreatorRegistry()
		So(r.Register("test_sink", SinkCreatorFunc(createCollectorSink)), ShouldBeNil)
		So(r.Register("test_sink2", SinkCreatorFunc(createCollectorSink)), ShouldBeNil)

		Convey("When adding a new type having the registered type name", func() {
			err := r.Register("test_sink", SinkCreatorFunc(createCollectorSink))

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When looking up a creator", func() {
			c, err := r.Lookup("test_sink2")

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And it should have the expected type", func() {
					s, err := c.CreateSink(ctx, nil)
					So(err, ShouldBeNil)
					So(s, ShouldHaveSameTypeAs, &tupleCollectorSink{})
				})
			})
		})

		Convey("When retrieving a list of creators", func() {
			m, err := r.List()

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And the list should have all creators", func() {
					So(len(m), ShouldEqual, 2)
					So(m["test_sink"], ShouldNotBeNil)
					So(m["test_sink2"], ShouldNotBeNil)
				})
			})
		})

		Convey("When unregistering a creator", func() {
			err := r.Unregister("test_sink")

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And the unregistered creator shouldn't be found", func() {
					_, err := r.Lookup("test_sink")
					So(err, ShouldNotBeNil)
				})

				Convey("And the other creator should be found", func() {
					_, err := r.Lookup("test_sink2")
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

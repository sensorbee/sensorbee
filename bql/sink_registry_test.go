package bql

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
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

func TestGlobalSinkCreatorRegistry(t *testing.T) {
	ctx := core.NewContext(nil)

	Convey("Given an default Global Sink registry", t, func() {
		r, err := CopyGlobalSinkCreatorRegistry()
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)
		Reset(func() {
			r, _ = CopyGlobalSinkCreatorRegistry()
		})

		Convey("When adding a new type having the registered type name", func() {
			err := r.Register("uds", SinkCreatorFunc(createCollectorSink))

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When looking up a creator", func() {
			c, err := r.Lookup("uds")

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And it should return error when parameter is empty", func() {
					_, err := c.CreateSink(ctx, nil)
					So(err.Error(), ShouldContainSubstring, "parameter")
				})

				Convey("And it should return error because the specified parameter is invalid type", func() {
					_, err := c.CreateSink(ctx, data.Map{"name": data.Int(100)})
					So(err.Error(), ShouldContainSubstring, "unsupported")
				})

				Convey("And it should return error because the specified shared state is missing", func() {
					_, err := c.CreateSink(ctx, data.Map{"name": data.String("shared_state_not_found")})
					So(err.Error(), ShouldContainSubstring, "was not found")
				})
			})
		})
	})
}

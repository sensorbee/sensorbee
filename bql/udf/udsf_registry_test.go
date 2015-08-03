package udf

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestEmptyDefaultUDSFCreatorRegistry(t *testing.T) {
	Convey("Given an empty default UDSF registry", t, func() {
		r := NewDefaultUDSFCreatorRegistry()

		Convey("When adding a creator function", func() {
			err := r.Register("duplicate", MustConvertToUDSFCreator(createDuplicateUDSF))

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When looking up a nonexistent creator", func() {
			_, err := r.Lookup("duplicate", 2)

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
			err := r.Unregister("duplicate")

			Convey("Then it shouldn't fail", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestDefaultUDSFCreatorRegistry(t *testing.T) {
	ctx := core.NewContext(nil)

	Convey("Given an default UDSF registry having two types", t, func() {
		r := NewDefaultUDSFCreatorRegistry()
		So(r.Register("DUPLICATE", MustConvertToUDSFCreator(createDuplicateUDSF)), ShouldBeNil)
		So(r.Register("DUPLICATE2", MustConvertToUDSFCreator(createDuplicateUDSF)), ShouldBeNil)

		Convey("When adding a new type having the registered type name", func() {
			err := r.Register("duplicate", MustConvertToUDSFCreator(createDuplicateUDSF))

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When looking up a creator", func() {
			c, err := r.Lookup("Duplicate2", 2)

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And it should have the expected type", func() {
					s, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.String("test_stream"), data.Int(2))
					So(err, ShouldBeNil)
					So(s, ShouldHaveSameTypeAs, &duplicateUDSF{})
				})
			})
		})

		Convey("When looking up a creator with a wrong arity", func() {
			_, err := r.Lookup("duplicate2", 1)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When retrieving a list of creators", func() {
			m, err := r.List()

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And the list should have all creators", func() {
					So(len(m), ShouldEqual, 2)
					So(m["duplicate"], ShouldNotBeNil)
					So(m["duplicate2"], ShouldNotBeNil)
				})
			})
		})

		Convey("When unregistering a creator", func() {
			err := r.Unregister("Duplicate")

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And the unregistered creator shouldn't be found", func() {
					_, err := r.Lookup("duplicate", 2)
					So(err, ShouldNotBeNil)
				})

				Convey("And the other creator should be found", func() {
					_, err := r.Lookup("duplicate2", 2)
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

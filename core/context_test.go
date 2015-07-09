package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAtomicFlag(t *testing.T) {
	Convey("Given an enabled atomic flag", t, func() {
		f := AtomicFlag(1)

		Convey("When switching to off", func() {
			f.Set(false)
			Convey("Then it should be disabled", func() {
				b := f.Enabled()
				So(b, ShouldBeFalse)
			})
		})
		Convey("When switching to on", func() {
			f.Set(true)
			Convey("Then it should be enabled", func() {
				b := f.Enabled()
				So(b, ShouldBeTrue)
			})
		})
	})

	Convey("Given a disabled atomic flag", t, func() {
		f := AtomicFlag(0)

		Convey("When switching to on", func() {
			f.Set(true)
			Convey("Then it should be enabled", func() {
				b := f.Enabled()
				So(b, ShouldBeTrue)
			})
		})
		Convey("When switching to off", func() {
			f.Set(false)
			Convey("Then it should be disabled", func() {
				b := f.Enabled()
				So(b, ShouldBeFalse)
			})
		})
	})
}

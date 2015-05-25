package core

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type fakeFatalError struct {
	fatal bool
}

func (f *fakeFatalError) Error() string {
	return "fake error message"
}

func (f *fakeFatalError) Fatal() bool {
	return f.fatal
}

func TestFatalError(t *testing.T) {
	Convey("Given an error", t, func() {
		Convey("When the error is nil", func() {
			Convey("Then it shouldn't be a fatal error", func() {
				So(IsFatalError(nil), ShouldBeFalse)
			})

			Convey("Then FatalError should panic", func() {
				So(func() {
					FatalError(nil)
				}, ShouldPanic)
			})
		})

		Convey("When the error implements Fatal method", func() {
			Convey("Then it can be a fatal error", func() {
				So(IsFatalError(&fakeFatalError{true}), ShouldBeTrue)
			})

			Convey("Then it can also be a non-fatal error by configuration", func() {
				So(IsFatalError(&fakeFatalError{false}), ShouldBeFalse)
			})

			Convey("Then a non-fatal error can be wrapped as another fatal error", func() {
				err := FatalError(&fakeFatalError{false})
				So(err, ShouldNotBeNil)
				So(IsFatalError(err), ShouldBeTrue)

				Convey("And the original error message shouldn't be changed", func() {
					So(err.Error(), ShouldEqual, "fake error message")
				})
			})
		})

		Convey("When the error doesn't implements Fatal method", func() {
			err := errors.New("test failure")

			Convey("Then it shouldn't be a fatal error", func() {
				So(IsFatalError(err), ShouldBeFalse)
			})

			Convey("Then it can be wrapped as a fatal error", func() {
				e := FatalError(err)
				So(e, ShouldNotBeNil)

				Convey("And the error should be fatal", func() {
					So(IsFatalError(e), ShouldBeTrue)
				})

				Convey("And the original error message shouldn't be changed", func() {
					So(e.Error(), ShouldEqual, "test failure")
				})
			})
		})
	})
}

type fakeTemporaryError struct {
	temporary bool
}

func (f *fakeTemporaryError) Error() string {
	return "fake error message"
}

func (f *fakeTemporaryError) Temporary() bool {
	return f.temporary
}

func TestTemporaryError(t *testing.T) {
	Convey("Given an error", t, func() {
		Convey("When the error is nil", func() {
			Convey("Then it shouldn't be a temporary error", func() {
				So(IsTemporaryError(nil), ShouldBeFalse)
			})

			Convey("Then TemporaryError should panic", func() {
				So(func() {
					TemporaryError(nil)
				}, ShouldPanic)
			})
		})

		Convey("When the error implements Temporary method", func() {
			Convey("Then it can be a temporary error", func() {
				So(IsTemporaryError(&fakeTemporaryError{true}), ShouldBeTrue)
			})

			Convey("Then it can also be a non-temporary error by configuration", func() {
				So(IsTemporaryError(&fakeTemporaryError{false}), ShouldBeFalse)

				Convey("And it can be wrapped as another temporary error", func() {
					err := TemporaryError(&fakeTemporaryError{false})
					So(err, ShouldNotBeNil)
					So(IsTemporaryError(err), ShouldBeTrue)
				})
			})

			Convey("Then a non-temporary error can be wrapped as another temporary error", func() {
				err := TemporaryError(&fakeTemporaryError{false})
				So(err, ShouldNotBeNil)
				So(IsTemporaryError(err), ShouldBeTrue)

				Convey("And the original error message shouldn't be changed", func() {
					So(err.Error(), ShouldEqual, "fake error message")
				})
			})
		})

		Convey("When the error doesn't implements Temporary method", func() {
			err := errors.New("test failure")

			Convey("Then it shouldn't be a temporary error", func() {
				So(IsTemporaryError(err), ShouldBeFalse)
			})

			Convey("Then it can be wrapped as a temporary error", func() {
				e := TemporaryError(err)
				So(e, ShouldNotBeNil)

				Convey("And the error should be temporary", func() {
					So(IsTemporaryError(e), ShouldBeTrue)
				})

				Convey("And the original error message shouldn't be changed", func() {
					So(e.Error(), ShouldEqual, "test failure")
				})
			})
		})
	})
}

package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSetTupleTraceEnabled(t *testing.T) {
	Convey("Given a context set tracing enabled", t, func() {
		config := Configuration{TupleTraceEnabled: 1}
		ctx := Context{Config: config}
		Convey("When switch to trace off", func() {
			ctx.SetTupleTraceEnabled(false)
			Convey("Then context's tracing configuration should return false", func() {
				f := ctx.IsTupleTraceEnabled()
				So(f, ShouldBeFalse)
			})
		})
		Convey("When switch to trace on", func() {
			ctx.SetTupleTraceEnabled(true)
			Convey("Then context's tracing configuration should return true", func() {
				f := ctx.IsTupleTraceEnabled()
				So(f, ShouldBeTrue)
			})
		})
	})
	Convey("Given a context set tracing not enabled", t, func() {
		config := Configuration{TupleTraceEnabled: 0}
		ctx := Context{Config: config}
		Convey("When switch to trace on", func() {
			ctx.SetTupleTraceEnabled(true)
			Convey("Then context's tracing configuration should return true", func() {
				f := ctx.IsTupleTraceEnabled()
				So(f, ShouldBeTrue)
			})
		})
		Convey("When switch to trace off", func() {
			ctx.SetTupleTraceEnabled(false)
			Convey("Then context's tracing configuration should return false", func() {
				f := ctx.IsTupleTraceEnabled()
				So(f, ShouldBeFalse)
			})
		})
	})
}

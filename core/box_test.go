package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
)

func TestBox(t *testing.T) {
	Convey("Given return nil BoxFunc implementation", t, func() {

		a := func(ctx *Context, t *tuple.Tuple, s Writer) error {
			return nil
		}
		bf := boxFunc(a)
		var b Box = &bf

		Convey("It should satisfy Box interface", func() {
			So(b, ShouldNotBeNil)
		})
	})
}

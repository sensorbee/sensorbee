package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestSharedStateSink(t *testing.T) {
	ctx := NewContext(nil)
	r := ctx.SharedStates
	counter := &countingSharedState{}

	{
		s := &stubSharedState{}
		if err := r.Add("test_state", s); err != nil {
			t.Fatal("Cannot add stub shared state:", err)
		}

		if err := r.Add("test_counter", counter); err != nil {
			t.Fatal("Cannot add counting shared state:", err)
		}
	}

	Convey("Given a topology and a state registry", t, func() {
		Convey("When creating shared state sink with non-Writer shared state", func() {
			_, err := NewSharedStateSink(ctx, "test_state")

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When creating shared state sink with Writer shared state", func() {
			sink, err := NewSharedStateSink(ctx, "test_counter")
			Reset(func() {
				counter.cnt = 0
			})

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("And adding counter via sink", func() {
				t := &Tuple{
					Data: data.Map{"n": data.Int(100)},
				}
				err := sink.Write(ctx, t)

				Convey("Then it should succeed", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then the count should be be incremented", func() {
					So(counter.cnt, ShouldEqual, 100)
				})
			})
		})
	})
}

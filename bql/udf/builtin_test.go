package udf

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestCountAggregate(t *testing.T) {
	Convey("Given an instance of countAggregate", t, func() {
		f := countAggregate{}

		Convey("Then it should be unary", func() {
			So(f.Accept(0), ShouldBeFalse)
			So(f.Accept(1), ShouldBeTrue)
			So(f.Accept(2), ShouldBeFalse)
		})

		Convey("Then it should be an aggregate in the first argument", func() {
			So(f.IsAggregationParameter(0), ShouldBeFalse)
			So(f.IsAggregationParameter(1), ShouldBeTrue)
			So(f.IsAggregationParameter(2), ShouldBeFalse)
		})

		Convey("When feeding it with a non-array", func() {
			input := data.Int(27)
			_, err := f.Call(nil, input)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When feeding it with multiple parameters", func() {
			input1 := data.Array{data.Int(1), data.Int(2)}
			input2 := data.Int(27)
			_, err := f.Call(nil, input1, input2)

			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When feeding it with an array", func() {
			input := data.Array{data.Int(1), data.Null{}, data.String("2")}
			result, err := f.Call(nil, input)

			Convey("Then it counts the number of non-null items", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, data.Int(2))
			})
		})
	})
}

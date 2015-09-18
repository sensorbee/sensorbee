package execution

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func TestMultiplicityHandling(t *testing.T) {
	Convey("Given an execution plan", t, func() {
		ep := streamRelationStreamExecutionPlan{}
		hashes := map[data.HashValue][]resultRowCount{}

		Convey("and an item/hash pairs", func() {
			a := resultRow{
				data.Map{"a": data.Int(5)},
				data.HashValue(17),
			}
			b := resultRow{
				data.Map{"a": data.Int(6)},
				data.HashValue(17),
			}
			c := resultRow{
				data.Map{"a": data.Int(7)},
				data.HashValue(18),
			}

			Convey("Then adding and counting should work correctly", func() {
				So(ep.currentMultiplicity(&a, hashes), ShouldEqual, 0)

				// add an item with data A and hash X -> 1
				So(ep.incrAndGetMultiplicity(&a, hashes), ShouldEqual, 1)
				So(ep.currentMultiplicity(&a, hashes), ShouldEqual, 1)

				// add an item with data B and hash X -> 1
				So(ep.incrAndGetMultiplicity(&b, hashes), ShouldEqual, 1)
				So(ep.currentMultiplicity(&b, hashes), ShouldEqual, 1)

				// add an item with data C and hash Y -> 1
				So(ep.incrAndGetMultiplicity(&c, hashes), ShouldEqual, 1)
				So(ep.currentMultiplicity(&c, hashes), ShouldEqual, 1)

				// add an item with data A and hash X -> 2
				So(ep.incrAndGetMultiplicity(&a, hashes), ShouldEqual, 2)
				So(ep.currentMultiplicity(&a, hashes), ShouldEqual, 2)
				So(ep.currentMultiplicity(&b, hashes), ShouldEqual, 1)
				So(ep.currentMultiplicity(&c, hashes), ShouldEqual, 1)
			})
		})
	})
}

package server

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"testing"
)

func TestDefaultTopologyRegistry(t *testing.T) {
	ctx := core.NewContext(nil)

	Convey("Given a default topology registry which has two topology builder", t, func() {
		r := NewDefaultTopologyRegistry()
		tpName1 := "test_TOPOLOGY"
		tp1, err := core.NewDefaultTopology(ctx, tpName1)
		So(err, ShouldBeNil)
		tb1, err := bql.NewTopologyBuilder(tp1)
		So(err, ShouldBeNil)
		So(r.Register(tpName1, tb1), ShouldBeNil)
		tpName2 := "test_TOPOLOGY2"
		tp2, err := core.NewDefaultTopology(ctx, tpName2)
		So(err, ShouldBeNil)
		tb2, err := bql.NewTopologyBuilder(tp2)
		So(err, ShouldBeNil)
		So(r.Register(tpName2, tb2), ShouldBeNil)

		Convey("When adding a new topology builder with registered name", func() {
			name := "test_topology"
			tp, err := core.NewDefaultTopology(ctx, name)
			So(err, ShouldBeNil)
			tb, err := bql.NewTopologyBuilder(tp)
			So(err, ShouldBeNil)
			err = r.Register(name, tb)
			Convey("Then it should fail", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When looking up a topology builder", func() {
			tb, err := r.Lookup("test_topology2")

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
				So(tb, ShouldNotBeNil)
			})
		})

		Convey("When retrieving a list of topology builders", func() {
			m, err := r.List()

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And the list should have all topology builder", func() {
					So(len(m), ShouldEqual, 2)
					So(m["test_topology"], ShouldNotBeNil)
					So(m["test_topology2"], ShouldNotBeNil)
				})
			})
		})

		Convey("When unregistering a topology builder", func() {
			tb, err := r.Unregister("TEST_topology")

			Convey("Then it should succeed", func() {
				So(err, ShouldBeNil)
				So(tb, ShouldNotBeNil)

				Convey("And the unregistered topology builder should not be found", func() {
					_, err := r.Lookup("test_topology")
					So(core.IsNotExist(err), ShouldBeTrue)
				})

				Convey("And the other topology builder should be found", func() {
					_, err := r.Lookup("test_topology2")
					So(err, ShouldBeNil)
				})

				Convey("And when unregistering the unregistered topology builder", func() {
					_, err := r.Unregister("test_TOPOLOGY")

					Convey("Then it should fail", func() {
						So(core.IsNotExist(err), ShouldBeTrue)
					})
				})
			})
		})
	})
}

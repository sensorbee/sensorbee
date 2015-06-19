package udf

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/tuple"
	"testing"
)

func TestDefaultFunctionRegistry(t *testing.T) {
	Convey("Given a default function registry", t, func() {
		fr := NewDefaultFunctionRegistry()

		Convey("When asking for an unknown function", func() {
			_, err := fr.Lookup("hoge", 17)
			Convey("An error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When adding a unary function via Register", func() {
			fun := func(vs ...tuple.Value) (tuple.Value, error) {
				return tuple.Bool(true), nil
			}
			fr.Register("test", fun, 1)

			Convey("Then it can be looked up as unary", func() {
				_, err := fr.Lookup("test", 1)
				So(err, ShouldBeNil)
			})

			Convey("And it won't be found as binary", func() {
				_, err := fr.Lookup("test", 2)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When adding a variadic function via RegisterVariadic", func() {
			fun := func(vs ...tuple.Value) (tuple.Value, error) {
				return tuple.Bool(true), nil
			}
			fr.RegisterVariadic("hello", fun)

			Convey("Then it can be looked up as unary", func() {
				_, err := fr.Lookup("hello", 1)
				So(err, ShouldBeNil)
			})
			Convey("And it can be looked up as binary", func() {
				_, err := fr.Lookup("hello", 2)
				So(err, ShouldBeNil)
			})
		})

		Convey("When adding a unary function via RegisterNullary", func() {
			fun := func() (tuple.Value, error) {
				return tuple.Bool(true), nil
			}
			fr.RegisterNullary("test0", fun)

			Convey("Then it can be looked up as nullary", func() {
				_, err := fr.Lookup("test0", 0)
				So(err, ShouldBeNil)
			})

			Convey("And it won't be found as binary", func() {
				_, err := fr.Lookup("test0", 2)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When adding a unary function via RegisterUnary", func() {
			fun := func(tuple.Value) (tuple.Value, error) {
				return tuple.Bool(true), nil
			}
			fr.RegisterUnary("test1", fun)

			Convey("Then it can be looked up as unary", func() {
				_, err := fr.Lookup("test1", 1)
				So(err, ShouldBeNil)
			})

			Convey("And it won't be found as binary", func() {
				_, err := fr.Lookup("test1", 2)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When adding a binary function via RegisterBinary", func() {
			fun := func(tuple.Value, tuple.Value) (tuple.Value, error) {
				return tuple.Bool(true), nil
			}
			fr.RegisterBinary("test2", fun)

			Convey("Then it can be looked up as binary", func() {
				_, err := fr.Lookup("test2", 2)
				So(err, ShouldBeNil)
			})

			Convey("And it won't be found as unary", func() {
				_, err := fr.Lookup("test2", 1)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When adding a binary function via RegisterTernary", func() {
			fun := func(tuple.Value, tuple.Value, tuple.Value) (tuple.Value, error) {
				return tuple.Bool(true), nil
			}
			fr.RegisterTernary("test3", fun)

			Convey("Then it can be looked up as ternary", func() {
				_, err := fr.Lookup("test3", 3)
				So(err, ShouldBeNil)
			})

			Convey("And it won't be found as unary", func() {
				_, err := fr.Lookup("test3", 1)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

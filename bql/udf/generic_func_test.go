package udf

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"strings"
	"testing"
)

func TestGenericFunc(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a generic UDF generator", t, func() {
		normalCases := []struct {
			title string
			f     interface{}
		}{
			{
				title: "When passing a function receiving a context and two args and returning an error",
				f: func(ctx *core.Context, i int, f float32) (float32, error) {
					return float32(i) + f, nil
				},
			},
			{
				title: "When passing a function receiving two args and returning an error",
				f: func(i int, f float32) (float32, error) {
					return float32(i) + f, nil
				},
			},
			{
				title: "When passing a function receiving a context and two args and not returning an error",
				f: func(ctx *core.Context, i int, f float32) float32 {
					return float32(i) + f
				},
			},
			{
				title: "When passing a function receiving two args and not returning an error",
				f: func(i int, f float32) float32 {
					return float32(i) + f
				},
			},
		}

		for _, c := range normalCases {
			c := c
			Convey(c.title, func() {
				f, err := ConvertGeneric(c.f)
				So(err, ShouldBeNil)

				Convey("Then the udf should return a correct value", func() {
					v, err := f.Call(ctx, data.Int(1), data.Float(1.5))
					So(err, ShouldBeNil)
					res, err := data.ToFloat(v)
					So(err, ShouldBeNil)
					So(res, ShouldEqual, 2.5)
				})

				Convey("Then the udf's arity should be 2", func() {
					So(f.Accept(2), ShouldBeTrue)
					So(f.Accept(1), ShouldBeFalse)
					So(f.Accept(3), ShouldBeFalse)
				})
			})
		}

		variadicCases := []struct {
			title string
			f     interface{}
		}{
			{
				title: "When passing a variadic function with a context and returning an error",
				f: func(ctx *core.Context, ss ...string) (string, error) {
					return strings.Join(ss, ""), nil
				},
			},
			{
				title: "When passing a variadic function without a context and returning an error",
				f: func(ss ...string) (string, error) {
					return strings.Join(ss, ""), nil
				},
			},
			{
				title: "When passing a variadic function with a context and not returing an error",
				f: func(ctx *core.Context, ss ...string) string {
					return strings.Join(ss, "")
				},
			},
			{
				title: "When passing a variadic function without a context and not returing an error",
				f: func(ss ...string) string {
					return strings.Join(ss, "")
				},
			},
		}

		for _, c := range variadicCases {
			c := c
			Convey(c.title, func() {
				f, err := ConvertGeneric(c.f)
				So(err, ShouldBeNil)

				Convey("And passing no arguments", func() {
					res, err := f.Call(ctx)

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldBeBlank)
					})
				})

				Convey("And passing one arguments", func() {
					res, err := f.Call(ctx, data.String("a"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldEqual, "a")
					})
				})

				Convey("And passing many arguments", func() {
					res, err := f.Call(ctx, data.String("a"), data.String("b"), data.String("c"), data.String("d"), data.String("e"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldEqual, "abcde")
					})
				})

				Convey("And passing a convertible value", func() {
					res, err := f.Call(ctx, data.String("a"), data.Int(1), data.String("c"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldEqual, "a1c")
					})
				})

				Convey("Then it should accept any arity", func() {
					So(f.Accept(0), ShouldBeTrue)
					So(f.Accept(1), ShouldBeTrue)
					So(f.Accept(123456789), ShouldBeTrue)
				})
			})
		}

		variadicExtraCases := []struct {
			title string
			f     interface{}
		}{
			{
				title: "When passing a variadic function having additional args with a context and returning values include error",
				f: func(ctx *core.Context, rep int, ss ...string) (string, error) {
					return strings.Repeat(strings.Join(ss, ""), rep), nil
				},
			},
			{
				title: "When passing a variadic function having additional args without a context and returning values include error",
				f: func(rep int, ss ...string) (string, error) {
					return strings.Repeat(strings.Join(ss, ""), rep), nil
				},
			},
			{
				title: "When passing a variadic function having additional args with a context and not returing an error",
				f: func(ctx *core.Context, rep int, ss ...string) string {
					return strings.Repeat(strings.Join(ss, ""), rep)
				},
			},
			{
				title: "When passing a variadic function having additional args without a context and not returing an error",
				f: func(rep int, ss ...string) string {
					return strings.Repeat(strings.Join(ss, ""), rep)
				},
			},
		}

		for _, c := range variadicExtraCases {
			c := c
			Convey(c.title, func() {
				f, err := ConvertGeneric(c.f)
				So(err, ShouldBeNil)

				Convey("And passing no arguments", func() {
					_, err := f.Call(ctx)

					Convey("Then it should fail", func() {
						So(err, ShouldNotBeNil)
					})
				})

				Convey("And only passing required arguments", func() {
					res, err := f.Call(ctx, data.Int(1))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldEqual, "")
					})
				})

				Convey("And passing one argument", func() {
					res, err := f.Call(ctx, data.Int(5), data.String("a"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldEqual, "aaaaa")
					})
				})

				Convey("And passing many arguments", func() {
					res, err := f.Call(ctx, data.Int(2), data.String("a"), data.String("b"), data.String("c"), data.String("d"), data.String("e"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldEqual, "abcdeabcde")
					})
				})

				Convey("And passing a convertible value", func() {
					res, err := f.Call(ctx, data.Int(3), data.String("a"), data.Int(1), data.String("c"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						s, err := data.AsString(res)
						So(err, ShouldBeNil)
						So(s, ShouldEqual, "a1ca1ca1c")
					})
				})

				Convey("Then it should accept any arity greater than 0", func() {
					So(f.Accept(0), ShouldBeFalse)
					So(f.Accept(1), ShouldBeTrue)
					So(f.Accept(123456789), ShouldBeTrue)
				})
			})
		}

		Convey("When creating a function returning an error", func() {
			f, err := ConvertGeneric(func() (int, error) {
				return 0, fmt.Errorf("test failure")
			})
			So(err, ShouldBeNil)

			Convey("Then calling it should fail", func() {
				_, err := f.Call(ctx)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "test failure")
			})
		})

		Convey("When creating a valid UDF with MustConvertGeneric", func() {
			Convey("Then it shouldn't panic", func() {
				So(func() {
					MustConvertGeneric(func() int { return 0 })
				}, ShouldNotPanic)
			})
		})

		Convey("When creating a invalid UDF with MustConvertGeneric", func() {
			Convey("Then it should panic", func() {
				So(func() {
					MustConvertGeneric(func() {})
				}, ShouldPanic)
			})
		})
	})
}

func TestGenericFuncInvalidCases(t *testing.T) {
	ctx := &core.Context{}
	toArgs := func(vs ...interface{}) data.Array {
		a, err := data.NewArray(vs)
		if err != nil {
			t.Fatal(err) // don't want to increase So's count unnecessarily
		}
		return a
	}

	Convey("Given a generic UDF generator", t, func() {
		genCases := []struct {
			title string
			f     interface{}
		}{
			{"with no return value", func() {}},
			{"with an unsupported type", func(error) int { return 0 }},
			{"with non-function type", 10},
			{"with non-error second return type", func() (int, int) { return 0, 0 }},
			{"with an unsupported type and non-error second return type", func(error) (int, int) { return 0, 0 }},
			{"with an unsupported return type", func() *core.Context { return nil }},
			{"with an unsupported interface return type", func() error { return nil }},
		}

		for _, c := range genCases {
			c := c
			Convey("When passing a function "+c.title, func() {
				_, err := ConvertGeneric(c.f)

				Convey("Then it should fail", func() {
					So(err, ShouldNotBeNil)
				})
			})
		}

		callCases := []struct {
			title string
			f     interface{}
			args  data.Array
		}{
			{
				title: "When calling a function with too few arguments",
				f:     func(int, int) int { return 0 },
				args:  toArgs(1),
			},
			{
				title: "When calling a function too many arguments",
				f:     func(int, int) int { return 0 },
				args:  toArgs(1, 2, 3, 4),
			},
			{
				title: "When calling a function with inconvertible arguments",
				f:     func(data.Map) int { return 0 },
				args:  toArgs("hoge"),
			},
			{
				title: "When calling a variadic function with inconvertible regular arguments",
				f:     func(data.Array, ...data.Map) int { return 0 },
				args:  toArgs("owata", data.Map{}, data.Map{}),
			},
			{
				title: "When calling a variadic function with inconvertible variadic arguments",
				f:     func(data.Array, ...data.Map) int { return 0 },
				args:  toArgs(data.Array{}, data.Map{}, "damepo", data.Map{}),
			},
		}

		for _, c := range callCases {
			c := c
			Convey(c.title, func() {
				f, err := ConvertGeneric(c.f)
				So(err, ShouldBeNil)

				Convey("Then it should fail", func() {
					_, err := f.Call(ctx, c.args...)
					So(err, ShouldNotBeNil)
				})
			})
		}
	})
}

func TestIntGenericInt8Func(t *testing.T) {
	ctx := &core.Context{} // not used in this test

	Convey("Given a function receiving int8", t, func() {
		f, err := ConvertGeneric(func(i int8) int8 {
			return i * 2
		})
		So(err, ShouldBeNil)

		Convey("When passing a valid value", func() {
			v, err := f.Call(ctx, data.Int(10))
			So(err, ShouldBeNil)

			Convey("Then it should be doubled", func() {
				i, err := data.ToInt(v)
				So(err, ShouldBeNil)
				So(i, ShouldEqual, 20)
			})
		})

		// TODO: add tests passing invalid values (out of range, inconvertible types)
	})
}

// TODO: add tests for all types.

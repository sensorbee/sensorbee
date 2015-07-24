package udf

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"strings"
	"testing"
	"time"
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

func TestGenericFuncInconvertibleType(t *testing.T) {
	ctx := &core.Context{} // not used in this test

	udfs := []UDF{
		MustConvertGeneric(func(i int8) int8 {
			return i * 2
		}),
		MustConvertGeneric(func(i int16) int16 {
			return i * 2
		}),
		MustConvertGeneric(func(i int32) int32 {
			return i * 2
		}),
		MustConvertGeneric(func(i int64) int64 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint8) uint8 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint16) uint16 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint32) uint32 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint64) uint64 {
			return i * 2
		}),
		MustConvertGeneric(func(i float32) float32 {
			return i * 2
		}),
		MustConvertGeneric(func(i float64) float64 {
			return i * 2
		}),
		MustConvertGeneric(func(b []byte) []byte {
			return b
		}),
		MustConvertGeneric(func(t time.Time) time.Time {
			return t
		}),
	}
	funcTypes := []string{"int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "blob", "timestamp"}
	if len(udfs) != 12 {
		t.Fatal("len of udfs isn't 12 but", len(udfs))
	}

	type InputType struct {
		typeName string
		value    data.Value
	}

	inconvertible := [][]InputType{
		{ // int8
			{"int", data.Int(int(math.MaxInt8) + 1)},
			{"negative int", data.Int(int(math.MinInt8) - 1)},
			{"float", data.Float(float64(math.MaxInt8) + 1.0)},
			{"negative float", data.Float(float64(math.MinInt8) - 1.0)},
			{"time", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC))},
		},
		{ // int16
			{"int", data.Int(int32(math.MaxInt16) + 1)},
			{"negative int", data.Int(int32(math.MinInt16) - 1)},
			{"float", data.Float(float64(math.MaxInt16) + 1.0)},
			{"negative float", data.Float(float64(math.MinInt16) - 1.0)},
			{"time", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC))},
		},
		{ // int32
			{"int", data.Int(int64(math.MaxInt32) + 1)},
			{"negative int", data.Int(int64(math.MinInt32) - 1)},
			{"float", data.Float(float64(math.MaxInt32) + 1.0)},
			{"negative float", data.Float(float64(math.MinInt32) - 1.0)},
			{"time", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC))},
		},
		{ // int64
			{"float", data.Float(float64(math.MaxUint64))},
			{"negative float", data.Float(float64(math.MinInt64) * 2)},
		},
		{ // uint8
			{"int", data.Int(int(math.MaxUint8) + 1)},
			{"negative int", data.Int(-1)},
			{"float", data.Float(float64(math.MaxUint8) + 1.0)},
			{"negative float", data.Float(-1.0)},
			{"time", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC))},
		},
		{ // uint16
			{"int", data.Int(int32(math.MaxUint16) + 1)},
			{"negative int", data.Int(-1)},
			{"float", data.Float(float64(math.MaxUint16) + 1.0)},
			{"negative float", data.Float(-1.0)},
			{"time", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC))},
		},
		{ // uint32
			{"int", data.Int(int64(math.MaxUint32) + 1)},
			{"negative int", data.Int(-1)},
			{"float", data.Float(float64(math.MaxUint32) + 1.0)},
			{"negative float", data.Float(-1.0)},
			{"time", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC))},
		},
		{ // uint64
			{"negative int", data.Int(-1)},
			{"float", data.Float(float64(math.MaxUint64))},
			{"negative float", data.Float(-1.0)},
		},
		{ // float32 covered by common values
		},
		{ // float64 covered by common values
		},
		{ // blob
			{"int", data.Int(1)},
			{"float", data.Float(1.0)},
			{"time", data.Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC))},
			{"array", data.Array([]data.Value{data.Int(10)})},
			{"map", data.Map{"key": data.Int(10)}},
		},
		{ // timestamp
			{"string", data.String("str")},
			{"blob", data.Blob([]byte("blob"))},
			{"array", data.Array([]data.Value{data.Int(10)})},
			{"map", data.Map{"key": data.Int(10)}},
		},
	}

	// common incovertible values for integer and float
	numInconvertibleValues := []InputType{
		{"null", data.Null{}},
		{"string", data.String("str")},
		{"blob", data.Blob([]byte("blob"))},
		{"array", data.Array([]data.Value{data.Int(10)})},
		{"map", data.Map{"key": data.Int(10)}},
	}
	// append common values for int8, int16, int32, int64/ uint8, uint16, uint32, uint64/ float32, float64 (4 + 4 + 2 patterns)
	for i := 0; i < 10; i++ {
		inconvertible[i] = append(inconvertible[i], numInconvertibleValues...)
	}

	Convey("Given UDFs and inconvertible values", t, func() {
		for i, f := range udfs {
			f := f
			for _, inc := range inconvertible[i] {
				t := inc.typeName
				v := inc.value
				Convey(fmt.Sprintf("When passing inconvertible value of %v for %v", t, funcTypes[i]), func() {
					_, err := f.Call(ctx, v)

					Convey("Then it should fail", func() {
						So(err, ShouldNotBeNil)
					})
				})
			}
		}
	})
}

func TestGenericIntAndFloatFunc(t *testing.T) {
	ctx := &core.Context{} // not used in this test

	udfs := []UDF{
		MustConvertGeneric(func(i int8) int8 {
			return i * 2
		}),
		MustConvertGeneric(func(i int16) int16 {
			return i * 2
		}),
		MustConvertGeneric(func(i int32) int32 {
			return i * 2
		}),
		MustConvertGeneric(func(i int64) int64 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint8) uint8 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint16) uint16 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint32) uint32 {
			return i * 2
		}),
		MustConvertGeneric(func(i uint64) uint64 {
			return i * 2
		}),
		MustConvertGeneric(func(i float32) float32 {
			return i * 2
		}),
		MustConvertGeneric(func(i float64) float64 {
			return i * 2
		}),
	}

	funcTypes := []string{"int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64"}
	Convey("Given a function receiving integer", t, func() {
		for i, f := range udfs {
			f := f
			i := i
			Convey(fmt.Sprintf("When passing a valid value for %v", funcTypes[i]), func() {
				v, err := f.Call(ctx, data.String("1"))
				So(err, ShouldBeNil)

				Convey("Then it should be doubled", func() {
					i, err := data.ToInt(v)
					So(err, ShouldBeNil)
					So(i, ShouldEqual, 2)
				})
			})
		}
	})
}

func TestGenericBoolFunc(t *testing.T) {
	ctx := &core.Context{} // not used in this test

	Convey("Given a function receiving bool", t, func() {
		f, err := ConvertGeneric(func(b bool) bool {
			return !b
		})
		So(err, ShouldBeNil)

		Convey("When passing a valid value", func() {
			v, err := f.Call(ctx, data.Int(1))
			So(err, ShouldBeNil)

			Convey("Then it should be false", func() {
				b, err := data.ToBool(v)
				So(err, ShouldBeNil)
				So(b, ShouldBeFalse)
			})
		})
	})
}

// TODO: add tests for all types.

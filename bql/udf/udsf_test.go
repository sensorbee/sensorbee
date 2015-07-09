package udf

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"strings"
	"testing"
)

type duplicateUDSF struct {
	dup int
}

func (d *duplicateUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	for i := 0; i < d.dup; i++ {
		w.Write(ctx, t)
	}
	return nil
}

func (d *duplicateUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createDuplicateUDSF(decl UDSFDeclarer, stream string, dup int) (UDSF, error) {
	if err := decl.Input(stream, &UDSFInputConfig{
		InputName: "test",
	}); err != nil {
		return nil, err
	}

	return &duplicateUDSF{
		dup: dup,
	}, nil
}

func init() {
	RegisterGlobalUDSFCreator("duplicate", MustConvertToUDSFCreator(createDuplicateUDSF))
}

func TestDuplicateUDSF(t *testing.T) {
	ctx := core.NewContext(nil)

	Convey("Given a udsf in a UDSF registry", t, func() {
		r, err := CopyGlobalUDSFCreatorRegistry()
		So(err, ShouldBeNil)

		c, err := r.Lookup("duplicate", 2)
		So(err, ShouldBeNil)

		Convey("When creating a UDSF", func() {
			decl := newUDSFDeclarer()
			f, err := c.CreateUDSF(ctx, decl, data.String("test_stream"), data.Int(2))
			So(err, ShouldBeNil)

			Convey("Then the declarer should have input information", func() {
				So(len(decl.inputs), ShouldEqual, 1)
				So(decl.inputs["test_stream"], ShouldNotBeNil)
				So(decl.inputs["test_stream"].InputName, ShouldEqual, "test")
			})

			Convey("Then it should duplicates tuples", func() {
				cnt := 0
				w := core.WriterFunc(func(ctx *core.Context, t *core.Tuple) error {
					cnt++
					return nil
				})
				So(f.Process(ctx, &core.Tuple{}, w), ShouldBeNil)
				So(cnt, ShouldEqual, 2)
			})
		})
	})
}

type testUDSF struct {
	f float32
	s string
}

func (t *testUDSF) Process(ctx *core.Context, _ *core.Tuple, w core.Writer) error {
	return nil
}

func (t *testUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func TestConvertToUDSFCretor(t *testing.T) {
	ctx := &core.Context{}
	Convey("Given a generic UDSFCreator generator", t, func() {
		normalCases := []struct {
			title string
			f     interface{}
		}{
			{
				title: "When passing a function receiving a context and two args",
				f: func(ctx *core.Context, decl UDSFDeclarer, i int, f float32) (UDSF, error) {
					return &testUDSF{
						f: float32(i) + f,
					}, nil
				},
			},
			{
				title: "When passing a function receiving two args",
				f: func(decl UDSFDeclarer, i int, f float32) (UDSF, error) {
					return &testUDSF{
						f: float32(i) + f,
					}, nil
				},
			},
		}

		for _, c := range normalCases {
			c := c
			Convey(c.title, func() {
				c, err := ConvertToUDSFCreator(c.f)
				So(err, ShouldBeNil)

				Convey("Then the creator should return a correct udsf", func() {
					v, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.Int(1), data.Float(1.5))
					So(err, ShouldBeNil)
					So(v.(*testUDSF).f, ShouldEqual, 2.5)
				})

				Convey("Then the creator's arity should be 2", func() {
					So(c.Accept(2), ShouldBeTrue)
					So(c.Accept(1), ShouldBeFalse)
					So(c.Accept(3), ShouldBeFalse)
				})
			})
		}

		variadicCases := []struct {
			title string
			f     interface{}
		}{
			{
				title: "When passing a variadic function with a context",
				f: func(ctx *core.Context, decl UDSFDeclarer, ss ...string) (UDSF, error) {
					return &testUDSF{
						s: strings.Join(ss, ""),
					}, nil
				},
			},
			{
				title: "When passing a variadic function without a context",
				f: func(decl UDSFDeclarer, ss ...string) (UDSF, error) {
					return &testUDSF{
						s: strings.Join(ss, ""),
					}, nil
				},
			},
		}

		for _, c := range variadicCases {
			c := c
			Convey(c.title, func() {
				c, err := ConvertToUDSFCreator(c.f)
				So(err, ShouldBeNil)

				Convey("And passing no arguments", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer())

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldBeBlank)
					})
				})

				Convey("And passing one arguments", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.String("a"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldEqual, "a")
					})
				})

				Convey("And passing many arguments", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.String("a"), data.String("b"), data.String("c"), data.String("d"), data.String("e"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldEqual, "abcde")
					})
				})

				Convey("And passing a convertible value", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.String("a"), data.Int(1), data.String("c"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldEqual, "a1c")
					})
				})

				Convey("Then it should accept any arity", func() {
					So(c.Accept(0), ShouldBeTrue)
					So(c.Accept(1), ShouldBeTrue)
					So(c.Accept(123456789), ShouldBeTrue)
				})
			})
		}

		variadicExtraCases := []struct {
			title string
			f     interface{}
		}{
			{
				title: "When passing a variadic function having additional args with a context",
				f: func(ctx *core.Context, decl UDSFDeclarer, rep int, ss ...string) (UDSF, error) {
					return &testUDSF{
						s: strings.Repeat(strings.Join(ss, ""), rep),
					}, nil
				},
			},
			{
				title: "When passing a variadic function having additional args without a context",
				f: func(decl UDSFDeclarer, rep int, ss ...string) (UDSF, error) {
					return &testUDSF{
						s: strings.Repeat(strings.Join(ss, ""), rep),
					}, nil
				},
			},
		}

		for _, c := range variadicExtraCases {
			c := c
			Convey(c.title, func() {
				c, err := ConvertToUDSFCreator(c.f)
				So(err, ShouldBeNil)

				Convey("And passing no arguments", func() {
					_, err := c.CreateUDSF(ctx, newUDSFDeclarer())

					Convey("Then it should fail", func() {
						So(err, ShouldNotBeNil)
					})
				})

				Convey("And only passing required arguments", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.Int(1))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldEqual, "")
					})
				})

				Convey("And passing one argument", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.Int(5), data.String("a"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldEqual, "aaaaa")
					})
				})

				Convey("And passing many arguments", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.Int(2), data.String("a"), data.String("b"), data.String("c"), data.String("d"), data.String("e"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldEqual, "abcdeabcde")
					})
				})

				Convey("And passing a convertible value", func() {
					res, err := c.CreateUDSF(ctx, newUDSFDeclarer(), data.Int(3), data.String("a"), data.Int(1), data.String("c"))

					Convey("Then it should succeed", func() {
						So(err, ShouldBeNil)
						So(res.(*testUDSF).s, ShouldEqual, "a1ca1ca1c")
					})
				})

				Convey("Then it should accept any arity greater than 0", func() {
					So(c.Accept(0), ShouldBeFalse)
					So(c.Accept(1), ShouldBeTrue)
					So(c.Accept(123456789), ShouldBeTrue)
				})
			})
		}

		Convey("When creating a function returning an error", func() {
			c, err := ConvertToUDSFCreator(func(decl UDSFDeclarer) (UDSF, error) {
				return nil, fmt.Errorf("test failure")
			})
			So(err, ShouldBeNil)

			Convey("Then calling it should fail", func() {
				_, err := c.CreateUDSF(ctx, newUDSFDeclarer())
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "test failure")
			})
		})
	})
}

func TestConvertToUDSFCreatorInvalidCases(t *testing.T) {
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
			{"with no declarer", func() (UDSF, error) { return nil, nil }},
			{"with no return value", func(UDSFDeclarer) {}},
			{"with no return error value", func(UDSFDeclarer) UDSF { return nil }},
			{"with an unsupported type", func(UDSFDeclarer, error) (UDSF, error) { return nil, nil }},
			{"with non-function type", 10},
			{"with non-error second return type", func(UDSFDeclarer) (UDSF, int) { return nil, 0 }},
			{"with an unsupported type and non-error second return type", func(UDSFDeclarer, error) (UDSF, int) { return nil, 0 }},
			{"with an unsupported return type", func(UDSFDeclarer) (*core.Context, error) { return nil, nil }},
			{"with an unsupported interface return type", func(UDSFDeclarer) (error, error) { return nil, nil }},
			{"with Context and without UDSFDeclarer", func(*core.Context, int) (UDSF, error) { return nil, nil }},
			{"without Context UDSFDeclarer", func(int) (UDSF, error) { return nil, nil }},
		}

		for _, c := range genCases {
			c := c
			Convey("When passing a function "+c.title, func() {
				_, err := ConvertToUDSFCreator(c.f)

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
				f:     func(UDSFDeclarer, int, int) (UDSF, error) { return &testUDSF{}, nil },
				args:  toArgs(1),
			},
			{
				title: "When calling a function too many arguments",
				f:     func(UDSFDeclarer, int, int) (UDSF, error) { return &testUDSF{}, nil },
				args:  toArgs(1, 2, 3, 4),
			},
			{
				title: "When calling a function with inconvertible arguments",
				f:     func(UDSFDeclarer, data.Map) (UDSF, error) { return &testUDSF{}, nil },
				args:  toArgs("hoge"),
			},
			{
				title: "When calling a variadic function with inconvertible regular arguments",
				f:     func(UDSFDeclarer, data.Array, ...data.Map) (UDSF, error) { return &testUDSF{}, nil },
				args:  toArgs("owata", data.Map{}, data.Map{}),
			},
			{
				title: "When calling a variadic function with inconvertible variadic arguments",
				f:     func(UDSFDeclarer, data.Array, ...data.Map) (UDSF, error) { return &testUDSF{}, nil },
				args:  toArgs(data.Array{}, data.Map{}, "damepo", data.Map{}),
			},
			{
				title: "When calling a function returning nil",
				f:     func(UDSFDeclarer) (UDSF, error) { return nil, nil },
				args:  toArgs(),
			},
		}

		for _, c := range callCases {
			c := c
			Convey(c.title, func() {
				creator, err := ConvertToUDSFCreator(c.f)
				So(err, ShouldBeNil)

				Convey("Then it should fail", func() {
					_, err := creator.CreateUDSF(ctx, newUDSFDeclarer(), c.args...)
					So(err, ShouldNotBeNil)
				})
			})
		}
	})
}

package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAssembleEmitter(t *testing.T) {
	Convey("Given a parseStack", t, func() {
		ps := parseStack{}
		Convey("When the stack contains only an ISTREAM item", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.AssembleEmitterOptions(6, 6)
			ps.AssembleEmitter()

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 6)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
					})
				})
			})
		})

		Convey("When the stack contains an ISTREAM item and a LIMIT clause", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 8, NumericLiteral{7})
			ps.AssembleEmitterLimit()
			ps.AssembleEmitterOptions(6, 8)
			ps.AssembleEmitter()

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(comp.EmitterOptions, ShouldResemble, []interface{}{EmitterLimit{7}})
					})
				})
			})
		})

		Convey("When the stack contains an ISTREAM item and an EVERY clause", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 8, NumericLiteral{7})
			ps.AssembleEmitterSampling(CountBasedSampling, 1)
			ps.AssembleEmitterOptions(6, 8)
			ps.AssembleEmitter()

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 8)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(comp.EmitterOptions, ShouldResemble, []interface{}{EmitterSampling{7, CountBasedSampling}})
					})
				})
			})
		})

		Convey("When the stack contains an ISTREAM item and an EVERY and a LIMIT clause", func() {
			ps.PushComponent(0, 4, Raw{"PRE"})
			ps.PushComponent(4, 6, Istream)
			ps.PushComponent(6, 8, NumericLiteral{2})
			ps.AssembleEmitterSampling(CountBasedSampling, 1)
			ps.PushComponent(8, 10, NumericLiteral{7})
			ps.AssembleEmitterLimit()
			ps.AssembleEmitterOptions(6, 10)
			ps.AssembleEmitter()

			Convey("Then AssembleEmitter transforms it into one item", func() {
				So(ps.Len(), ShouldEqual, 2)

				Convey("And that item is a EmitterAST", func() {
					top := ps.Peek()
					So(top, ShouldNotBeNil)
					So(top.begin, ShouldEqual, 4)
					So(top.end, ShouldEqual, 10)
					So(top.comp, ShouldHaveSameTypeAs, EmitterAST{})

					Convey("And it contains the previously pushed data", func() {
						comp := top.comp.(EmitterAST)
						So(comp.EmitterType, ShouldEqual, Istream)
						So(comp.EmitterOptions, ShouldResemble, []interface{}{
							EmitterSampling{2, CountBasedSampling}, EmitterLimit{7}})
					})
				})
			})
		})

		Convey("When the stack contains a wrong emitter item", func() {
			ps.PushComponent(0, 4, Raw{"wrong"})
			f := func() {
				ps.AssembleEmitter()
			}

			Convey("Then AssembleEmitter should panic", func() {
				So(f, ShouldPanic)
			})
		})
	})

	Convey("Given a parser", t, func() {
		p := &bqlPeg{}

		Convey("When using ISTREAM without a specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with a LIMIT specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [LIMIT 7] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{EmitterLimit{7}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with an EVERY k-TH TUPLE specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 2-ND TUPLE] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{2, CountBasedSampling}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with an EVERY k MILLISECONDS specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 200 MILLISECONDS] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{0.2, TimeBasedSampling}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with an EVERY k MILLISECONDS (float) specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 2.5 MILLISECONDS] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{0.0025, TimeBasedSampling}})

				Convey("And String() should almost return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with an EVERY k SECONDS specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 2 SECONDS] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{2, TimeBasedSampling}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with an EVERY k SECONDS specifier (float)", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 2.5 SECONDS] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{2.5, TimeBasedSampling}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with a SAMPLE specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [SAMPLE 20%] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{20, RandomizedSampling}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with a SAMPLE specifier (float)", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [SAMPLE 0.01%] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{0.01, RandomizedSampling}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})

		Convey("When using ISTREAM with EVERY and LIMIT specifier", func() {
			p.Buffer = "CREATE STREAM x AS SELECT ISTREAM [EVERY 4-TH TUPLE LIMIT 7] 2 FROM a [RANGE 1 TUPLES]"
			p.Init()

			Convey("Then the statement should be parsed correctly", func() {
				err := p.Parse()
				So(err, ShouldBeNil)
				p.Execute()

				ps := p.parseStack
				So(ps.Len(), ShouldEqual, 1)
				top := ps.Peek().comp
				So(top, ShouldHaveSameTypeAs, CreateStreamAsSelectStmt{})
				comp := top.(CreateStreamAsSelectStmt)

				So(comp.Name, ShouldEqual, "x")
				So(comp.Select.EmitterType, ShouldEqual, Istream)
				So(comp.Select.EmitterOptions, ShouldResemble, []interface{}{
					EmitterSampling{4, CountBasedSampling}, EmitterLimit{7}})

				Convey("And String() should return the original statement", func() {
					So(comp.String(), ShouldEqual, p.Buffer)
				})
			})
		})
	})
}

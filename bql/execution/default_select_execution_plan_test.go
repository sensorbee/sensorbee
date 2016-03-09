package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"sort"
	"testing"
	"time"
)

func getTuples(num int) []*core.Tuple {
	tuples := make([]*core.Tuple, 0, num)
	for i := 0; i < num; i++ {
		tup := core.Tuple{
			Data: data.Map{
				"int": data.Int(i + 1),
			},
			InputName:     "src",
			Timestamp:     time.Date(2015, time.April, 10, 10, 23, i, 0, time.UTC),
			ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, i, 0, time.UTC),
			BatchID:       7,
		}
		tuples = append(tuples, &tup)
	}
	return tuples
}

func createDefaultSelectPlan(s string, t *testing.T) (PhysicalPlan, error) {
	p := parser.New()
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	_stmt, _, err := p.ParseStmt(s)
	if err != nil {
		return nil, err
	}
	So(_stmt, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
	stmt := _stmt.(parser.CreateStreamAsSelectStmt).Select
	logicalPlan, err := Analyze(stmt, reg)
	if err != nil {
		return nil, err
	}
	canBuild := CanBuildDefaultSelectExecutionPlan(logicalPlan, reg)
	if !canBuild {
		err := fmt.Errorf("defaultSelectExecutionPlan cannot be used for statement: %s", s)
		return nil, err
	}
	return NewDefaultSelectExecutionPlan(logicalPlan, reg)
}

func TestDefaultSelectExecutionPlan(t *testing.T) {
	// Select constant
	Convey("Given a SELECT clause with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then that constant should appear in %v", idx), func() {
					if idx <= 2 {
						// items should be emitted as the number of
						// rows increases
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"col_1": data.Float(2.0), "col_2": data.Bool(true),
								"col_3": data.Null{}, "col_4": data.Bool(true)})
					} else {
						// nothing should be emitted because no new
						// data appears
						So(len(out), ShouldEqual, 0)
					}
				})
			}

		})
	})

	// Select a column with changing values
	Convey("Given a SELECT clause with only a column", t, func() {
		tuples := getTuples(4)
		tuples[2].Data["int"] = data.Null{}
		s := `CREATE STREAM box AS SELECT ISTREAM int::string AS int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					if idx == 2 {
						So(out[0], ShouldResemble,
							data.Map{"int": data.Null{}})
					} else {
						So(out[0], ShouldResemble,
							data.Map{"int": data.String(fmt.Sprintf("%d", idx+1))})
					}
				})
			}

		})
	})

	Convey("Given a SELECT clause with a column and timestamp", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int, now(), clock_timestamp() AS t FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0]["int"], ShouldEqual, data.Int(1))
						So(out[0]["now"], ShouldHaveSameTypeAs, data.Timestamp{})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						// the result of now() and clock_timestamp() does
						// not change for the previous tuple after it was
						// emitted, so it is not emitted via ISTREAM
						So(out[0]["int"], ShouldEqual, data.Int(2))
						So(out[0]["now"], ShouldHaveSameTypeAs, data.Timestamp{})
					}
				})
			}

		})
	})

	Convey("Given a SELECT clause with only a column using the table name", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM src:int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"int": data.Int(idx + 1)})
				})
			}

		})
	})

	// Select the tuple's timestamp
	Convey("Given a SELECT clause with only the timestamp", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM ts() FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"ts": data.Timestamp(time.Date(2015, time.April, 10,
							10, 23, idx, 0, time.UTC))})
				})
			}

		})
	})

	// Select a non-existing column
	Convey("Given a SELECT clause with a non-existing column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM hoge FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for _, inTup := range tuples {
				_, err := plan.Process(inTup)
				So(err, ShouldNotBeNil) // hoge not found
			}

		})
	})

	Convey("Given a SELECT clause with a non-existing column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM hoge + 1 FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for _, inTup := range tuples {
				_, err := plan.Process(inTup)
				So(err, ShouldNotBeNil) // hoge not found
			}

		})
	})

	// Select constant and a column with changing values
	Convey("Given a SELECT clause with a constant and a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2, int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"col_1": data.Int(2), "int": data.Int(idx + 1)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with an operation on a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 0 - int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"col_1": data.Int(-(idx + 1))})
				})
			}

		})
	})

	// Select constant and a column with changing values from aliased relation
	Convey("Given a SELECT clause with a constant, a column, and a table alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2, int FROM src [RANGE 2 SECONDS] AS x`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"col_1": data.Int(2), "int": data.Int(idx + 1)})
				})
			}

		})
	})

	// Select NULL-related operations
	Convey("Given a SELECT clause with NULL operations", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM null IS NULL, null + 2 = 2 FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the null operations should be correct %v", idx), func() {
					if idx <= 2 {
						// items should be emitted as the number of
						// rows increases
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"col_1": data.Bool(true), "col_2": data.Null{}})
					} else {
						// nothing should be emitted because no new
						// data appears
						So(len(out), ShouldEqual, 0)
					}
				})
			}

		})
	})

	Convey("Given a SELECT clause with NULL filter", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int FROM src [RANGE 2 SECONDS] WHERE null`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then there should be no rows in the output %v", idx), func() {
					So(len(out), ShouldEqual, 0)
				})
			}

		})
	})

	Convey("Given a SELECT clause with a non-boolean filter", t, func() {
		tuples := getTuples(1)
		s := `CREATE STREAM box AS SELECT ISTREAM int FROM src [RANGE 2 SECONDS] WHERE 6`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for _, inTup := range tuples {
				_, err := plan.Process(inTup)

				Convey("Then there should be an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "unsupported cast bool from int")
				})
			}

		})
	})

	// Select constant and a column with changing values from aliased relation
	// using that alias
	Convey("Given a SELECT clause with a constant, a table alias, and a column using it", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2, x:int FROM src [RANGE 2 SECONDS] AS x`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"col_1": data.Int(2), "int": data.Int(idx + 1)})
				})
			}

		})
	})

	// Use alias
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int-1 AS a, int AS b FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"a": data.Int(idx), "b": data.Int(idx + 1)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with a nested column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int-1 AS a.c, int+1 AS a["d"], int AS b[1] FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"a": data.Map{"c": data.Int(idx), "d": data.Int(idx + 2)},
							"b": data.Array{data.Null{}, data.Int(idx + 1)}})
				})
			}

		})
	})

	Convey("Given a SELECT clause with various expressions", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM CASE int WHEN 1 THEN int+1 WHEN 3 THEN "b" ELSE "c" END AS x FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"x": data.Int(2)})
					} else if idx == 2 {
						So(out[0], ShouldResemble, data.Map{"x": data.String("b")})
					} else {
						So(out[0], ShouldResemble, data.Map{"x": data.String("c")})
					}
				})
			}

		})
	})

	Convey("Given a SELECT clause with a complex JSON Path", t, func() {
		tuples := getTuples(4)
		for i, t := range tuples {
			k := i + 1
			t.Data["d"] = data.Map{
				"foo": data.Array{
					data.Map{"hoge": data.Array{
						data.Map{"a": data.Int(1 + k), "b": data.Int(2 * k)},
						data.Map{"a": data.Int(3 + k), "b": data.Int(4 * k)},
					}, "bar": data.Int(5)},
					data.Map{"hoge": data.Array{
						data.Map{"a": data.Int(5 + k), "b": data.Int(6 * k)},
						data.Map{"a": data.Int(7 + k), "b": data.Int(8 * k)},
					}, "bar": data.Int(2)},
					data.Map{"hoge": data.Array{
						data.Map{"a": data.Int(9 + k), "b": data.Int(10 * k)},
					}, "bar": data.Int(8)},
				},
				"nantoka": data.Map{"x": data.String("y")},
			}
		}
		s := `CREATE STREAM box AS SELECT ISTREAM d..b AS bs, d..b, d.foo[1:].hoge[0].a AS as, d.foo[1:].hoge[0].a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					if idx == 0 {
						So(out[0], ShouldResemble,
							data.Map{"bs": data.Array{
								data.Int(2), data.Int(4), data.Int(6), data.Int(8), data.Int(10),
							},
								"col_2": data.Array{
									data.Int(2), data.Int(4), data.Int(6), data.Int(8), data.Int(10),
								},
								"as": data.Array{
									data.Int(6), data.Int(10),
								},
								"col_4": data.Array{
									data.Int(6), data.Int(10),
								},
							})
					} else if idx == 1 {
						So(out[0], ShouldResemble,
							data.Map{"bs": data.Array{
								data.Int(4), data.Int(8), data.Int(12), data.Int(16), data.Int(20),
							},
								"col_2": data.Array{
									data.Int(4), data.Int(8), data.Int(12), data.Int(16), data.Int(20),
								},
								"as": data.Array{
									data.Int(7), data.Int(11),
								},
								"col_4": data.Array{
									data.Int(7), data.Int(11),
								},
							})
					}
				})
			}

		})
	})

	Convey("Given a SELECT clause with an AS * flattening", t, func() {
		tuples := getTuples(4)
		for i := range tuples {
			if i == 1 {
				tuples[i].Data["nest"] = data.Int(i)
			} else {
				tuples[i].Data["nest"] = data.Map{"c": data.Int(i), "d": data.Float(i + 1)}
			}
		}
		s := `CREATE STREAM box AS SELECT ISTREAM nest AS * FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				if idx == 0 {
					So(err, ShouldBeNil)
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"c": data.Int(0), "d": data.Float(1.0)})
					})

				} else if idx == 1 || idx == 2 { // window size is 2!
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "tried to use value 1 as columns, but is not a map")

				} else {
					So(err, ShouldBeNil)
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"c": data.Int(2), "d": data.Float(3.0)})
						So(out[1], ShouldResemble,
							data.Map{"c": data.Int(3), "d": data.Float(4.0)})
					})
				}
			}

		})
	})

	// Use wildcard
	Convey("Given a SELECT clause with a wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM * FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"int": data.Int(idx + 1)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with a wildcard using stream prefix", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM src:* FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"int": data.Int(idx + 1)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with a wildcard in a map", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM {"a": src:*} AS x FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"x": data.Map{"a": data.Map{"int": data.Int(idx + 1)}}})
				})
			}

		})
	})

	Convey("Given a SELECT clause with a wildcard and a named column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM *, (int-1)*2 AS int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the column is prioritized %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"int": data.Int(2 * idx)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with a stream wildcard and a named column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM src:*, (src:int-1)*2 AS int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the column is prioritized %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"int": data.Int(2 * idx)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with a named column and a wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM (int-1)*2 AS int, * FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the column is prioritzed %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"int": data.Int(2 * idx)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with an aliased wildcard and an anonymous column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM * AS x, (int-1)*2 FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"col_2": data.Int(2 * idx), "x": data.Map{"int": data.Int(idx + 1)}})
				})
			}

		})
	})

	Convey("Given a SELECT clause with an aliased wildcard and an anonymous column (both prefixed)", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM src:* AS x, (src:int-1)*2 FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"col_2": data.Int(2 * idx), "x": data.Map{"int": data.Int(idx + 1)}})
				})
			}

		})
	})

	Convey("Given a SELECT clause with an aliased wildcard and a prefixed anonymous column (both prefixed)", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM * AS x, (src:int-1)*2 FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"col_2": data.Int(2 * idx), "x": data.Map{"int": data.Int(idx + 1)}})
				})
			}

		})
	})

	// Use a filter
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int AS b FROM src [RANGE 2 SECONDS] 
            WHERE int % 2 = 0`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if (idx+1)%2 == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"b": data.Int(idx + 1)})
					} else {
						So(len(out), ShouldEqual, 0)
					}
				})
			}

		})
	})
}

func TestDefaultSelectExecutionPlanEmitters(t *testing.T) {
	// Recovery from errors in tuples
	Convey("Given a SELECT clause with a column that does not exist in one tuple (RSTREAM)", t, func() {
		tuples := getTuples(6)
		// remove the selected key from one tuple
		delete(tuples[1].Data, "int")

		s := `CREATE STREAM box AS SELECT RSTREAM int FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)

				if idx == 0 {
					// In the idx==0 run, the window contains only item 0.
					// That item is fine, no problem.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				} else if idx == 1 || idx == 2 {
					// In the idx==1 run, the window contains item 0 and item 1,
					// the latter is broken, therefore the query fails.
					// In the idx==2 run, the window contains item 1 and item 2,
					// the latter is broken, therefore the query fails.
					Convey(fmt.Sprintf("Then there should be an error for a queries in %v", idx), func() {
						So(err, ShouldNotBeNil)
					})
				} else {
					// In later runs, we have recovered from the error in item 1
					// and emit one item per run as normal.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx)})
						So(out[1], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				}
			}

		})
	})

	Convey("Given a SELECT clause with a column that does not exist in one tuple (ISTREAM)", t, func() {
		tuples := getTuples(6)
		// remove the selected key from one tuple
		delete(tuples[1].Data, "int")

		s := `CREATE STREAM box AS SELECT ISTREAM int FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)

				if idx == 0 {
					// In the idx==0 run, the window contains only item 0.
					// That item is fine, no problem.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				} else if idx == 1 || idx == 2 {
					// In the idx==1 run, the window contains item 0 and item 1,
					// the latter is broken, therefore the query fails.
					// In the idx==2 run, the window contains item 1 and item 2,
					// the former is broken, therefore the query fails.
					Convey(fmt.Sprintf("Then there should be an error for a queries in %v", idx), func() {
						So(err, ShouldNotBeNil)
					})
				} else if idx == 3 {
					// In the idx==3 run, the window contains item 2 and item 3.
					// Both items are fine and have not been emitted before, so
					// both are emitted now.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx)})
						So(out[1], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				} else {
					// In later runs, we have recovered from the error in item 1
					// and emit one item per run as normal.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				}
			}

		})
	})

	Convey("Given a WHERE clause with a column that does not exist in one tuple (ISTREAM)", t, func() {
		tuples := getTuples(6)
		// remove the selected key from one tuple
		delete(tuples[1].Data, "int")

		s := `CREATE STREAM box AS SELECT ISTREAM int FROM src [RANGE 2 TUPLES] WHERE int > 0`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)

				if idx == 0 {
					// In the idx==0 run, the window contains only item 0.
					// That item is fine, no problem.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				} else if idx == 1 {
					// In the idx==1 run, the window contains item 0 and item 1,
					// the latter is broken, therefore the query fails.
					Convey(fmt.Sprintf("Then there should be an error for a queries in %v", idx), func() {
						So(err, ShouldNotBeNil)
					})
				} else if idx == 2 {
					// In the idx==2 run, the window contains item 1 and item 2,
					// only the latter is evaluated so there is no problem.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				} else if idx == 3 {
					// In the idx==3 run, the window contains item 2 and item 3.
					// Both items are fine and have not been emitted before, so
					// both are emitted now.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				} else {
					// In later runs, we have recovered from the error in item 1
					// and emit one item per run as normal.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx + 1)})
					})
				}
			}

		})
	})

	Convey("Given a SELECT clause with a column that does not exist in one tuple (DSTREAM)", t, func() {
		tuples := getTuples(6)
		// remove the selected key from one tuple
		delete(tuples[1].Data, "int")

		s := `CREATE STREAM box AS SELECT DSTREAM int FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)

				if idx == 0 {
					// In the idx==0 run, the window contains only item 0.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 0)
					})
				} else if idx == 1 || idx == 2 {
					// In the idx==1 run, the window contains item 0 and item 1,
					// the latter is broken, therefore the query fails.
					// In the idx==2 run, the window contains item 1 and item 2,
					// the latter is broken, therefore the query fails.
					Convey(fmt.Sprintf("Then there should be an error for a queries in %v", idx), func() {
						So(err, ShouldNotBeNil)
					})
				} else if idx == 3 {
					// In the idx==3 run, the window contains item 2 and item 3.
					// Both items are fine and so item 0 is emitted.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(1)})
					})
				} else {
					// In later runs, we have recovered from the error in item 1
					// and emit one item per run as normal.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"int": data.Int(idx - 1)})
					})
				}
			}

		})
	})

	// RSTREAM/2 SECONDS window
	Convey("Given an RSTREAM emitter selecting a constant and a 2 SECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2 AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[1][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 3)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][2], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[3]), ShouldEqual, 3)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[3][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[3][2], ShouldResemble, data.Map{"a": data.Int(2)})
			})

		})
	})

	Convey("Given an RSTREAM emitter selecting a column and a 2 SECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(output[1][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 3)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(output[2][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][2], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 3)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[3][1], ShouldResemble, data.Map{"a": data.Int(3)})
				So(output[3][2], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	Convey("Given an RSTREAM emitter selecting a column and a 2000 MILLISECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 2000 MILLISECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(output[1][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 3)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(output[2][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][2], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 3)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[3][1], ShouldResemble, data.Map{"a": data.Int(3)})
				So(output[3][2], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	Convey("Given an RSTREAM emitter selecting a column and a 20 MILLISECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 20 MILLISECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 1)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	// RSTREAM/2 TUPLES window
	Convey("Given an RSTREAM emitter selecting a constant and a 2 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2 AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[1][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 2)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[3]), ShouldEqual, 2)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[3][1], ShouldResemble, data.Map{"a": data.Int(2)})
			})

		})
	})

	Convey("Given an RSTREAM emitter selecting a column and a 2 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole window state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(output[1][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 2)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][1], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 2)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(3)})
				So(output[3][1], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	// ISTREAM/2 SECONDS window
	Convey("Given an ISTREAM emitter selecting a constant and a 2 SECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2 AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(2)})
				// items should be emitted as (long as) the number of
				// rows increases
				So(len(output[1]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[3]), ShouldEqual, 0)
			})

		})
	})

	Convey("Given an ISTREAM emitter selecting a column and a 2 SECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 1)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	// ISTREAM/2 TUPLES window
	Convey("Given an ISTREAM emitter selecting a constant and a 2 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2 AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(2)})
				// items should be emitted as (long as) the number of
				// rows increases
				So(len(output[1]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 0)
				So(len(output[3]), ShouldEqual, 0)
			})

		})
	})

	Convey("Given an ISTREAM emitter selecting a column and a 2 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 1)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	// DSTREAM/2 SECONDS window
	Convey("Given a DSTREAM emitter selecting a constant and a 2 SECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM 2 AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then items dropped from state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 0)
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 0)
				So(len(output[3]), ShouldEqual, 0)
			})

		})
	})

	Convey("Given a DSTREAM emitter selecting a column and a 2 SECONDS window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM int AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then items dropped from state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 0)
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 0)
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(1)})
			})

		})
	})

	// DSTREAM/2 TUPLES window
	Convey("Given a DSTREAM emitter selecting a constant and a 2 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM 2 AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then items dropped from state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 0)
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 0)
				So(len(output[3]), ShouldEqual, 0)
			})

		})
	})

	Convey("Given a DSTREAM emitter selecting a column and a 2 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then items dropped from state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 0)
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 1)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(2)})
			})

		})
	})

	Convey("Given a DSTREAM emitter selecting a column (varying multiplicities) and a 2 TUPLES window", t, func() {
		tuples := getTuples(6)
		tuples[1].Data["int"] = data.Int(1)
		s := `CREATE STREAM box AS SELECT DSTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then items dropped from state should be emitted", func() {
				So(len(output), ShouldEqual, 6)
				So(len(output[0]), ShouldEqual, 0)
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 1)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[4][0], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[5][0], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	// fractional RANGE sizes
	Convey("Given an RSTREAM emitter selecting a column and a 1.9 SECONDS (=2 TUPLES) window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 1.9 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole window state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(output[1][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 2)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][1], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 2)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(3)})
				So(output[3][1], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})

	Convey("Given an RSTREAM emitter selecting a column and a 1900.9 MILLISECONDS (=2 TUPLES) window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 1900.9 MILLISECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole window state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, data.Map{"a": data.Int(1)})
				So(output[1][1], ShouldResemble, data.Map{"a": data.Int(2)})
				So(len(output[2]), ShouldEqual, 2)
				So(output[2][0], ShouldResemble, data.Map{"a": data.Int(2)})
				So(output[2][1], ShouldResemble, data.Map{"a": data.Int(3)})
				So(len(output[3]), ShouldEqual, 2)
				So(output[3][0], ShouldResemble, data.Map{"a": data.Int(3)})
				So(output[3][1], ShouldResemble, data.Map{"a": data.Int(4)})
			})

		})
	})
}

// sortedMapString computes a reliable string representation,
// i.e., always with the same order
// (does not work with nested maps)
func sortedMapString(inMap data.Map) string {
	// get the keys in correct order
	keys := make(sort.StringSlice, 0, len(inMap))
	for key := range inMap {
		keys = append(keys, key)
	}
	keys.Sort()
	// now build a reproducible string
	out := "{"
	for idx, key := range keys {
		if idx != 0 {
			out += ", "
		}
		out += fmt.Sprintf(`"%s": %s`, key, inMap[key])
	}
	out += "}"
	return out
}

// tupleList implements sort.Interface for []data.Map based on
// its string representation as per sortedMapString().
type tupleList []data.Map

func (tl tupleList) Len() int {
	return len(tl)
}
func (tl tupleList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}
func (tl tupleList) Less(i, j int) bool {
	return sortedMapString(tl[i]) < sortedMapString(tl[j])
}

func TestDefaultSelectExecutionPlanJoin(t *testing.T) {
	Convey("Given a JOIN selecting from left and right", t, func() {
		tuples := getTuples(8)
		// rearrange the tuples
		for i, t := range tuples {
			if i%2 == 0 {
				t.InputName = "src1"
				t.Data["l"] = data.String(fmt.Sprintf("l%d", i))
			} else {
				t.InputName = "src2"
				t.Data["r"] = data.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT ISTREAM src1:l, src2:r FROM src1 [RANGE 2 TUPLES], src2 [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l0"),
							"r": data.String("r1"),
						})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l2"),
							"r": data.String("r1"),
						})
					} else if idx%2 == 1 {
						// a tuple from src2 (=right) was just added
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String(fmt.Sprintf("l%d", idx-3)),
							"r": data.String(fmt.Sprintf("r%d", idx)),
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String(fmt.Sprintf("l%d", idx-1)),
							"r": data.String(fmt.Sprintf("r%d", idx)),
						})
					} else {
						// a tuple from src1 (=left) was just added
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String(fmt.Sprintf("l%d", idx)),
							"r": data.String(fmt.Sprintf("r%d", idx-3)),
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String(fmt.Sprintf("l%d", idx)),
							"r": data.String(fmt.Sprintf("r%d", idx-1)),
						})
					}
				})
			}
		})
	})

	Convey("Given a JOIN selecting with a big wildcard", t, func() {
		tuples := getTuples(4)
		// rearrange the tuples
		for i, t := range tuples {
			if i%2 == 0 {
				t.InputName = "src1"
				t.Data["l"] = data.String(fmt.Sprintf("l%d", i))
				delete(t.Data, "int")
				t.Data["hoge"] = data.Int(i)
			} else {
				t.InputName = "src2"
				t.Data["r"] = data.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT RSTREAM * FROM src1 [RANGE 2 TUPLES], src2 [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l":    data.String("l0"),
							"r":    data.String("r1"),
							"hoge": data.Int(0),
							"int":  data.Int(2),
						})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, data.Map{
							"l":    data.String("l0"),
							"r":    data.String("r1"),
							"hoge": data.Int(0),
							"int":  data.Int(2),
						})
						So(out[1], ShouldResemble, data.Map{
							"l":    data.String("l2"),
							"r":    data.String("r1"),
							"hoge": data.Int(2),
							"int":  data.Int(2),
						})
					}
				})
			}
		})
	})

	Convey("Given a JOIN selecting with a left wildcard", t, func() {
		tuples := getTuples(4)
		// rearrange the tuples
		for i, t := range tuples {
			if i%2 == 0 {
				t.InputName = "src1"
				t.Data["l"] = data.String(fmt.Sprintf("l%d", i))
				delete(t.Data, "int")
				t.Data["hoge"] = data.Int(i)
			} else {
				t.InputName = "src2"
				t.Data["r"] = data.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT RSTREAM src1:*, src2:r FROM src1 [RANGE 2 TUPLES], src2 [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l":    data.String("l0"),
							"r":    data.String("r1"),
							"hoge": data.Int(0),
						})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, data.Map{
							"l":    data.String("l0"),
							"r":    data.String("r1"),
							"hoge": data.Int(0),
						})
						So(out[1], ShouldResemble, data.Map{
							"l":    data.String("l2"),
							"r":    data.String("r1"),
							"hoge": data.Int(2),
						})
					}
				})
			}
		})
	})

	Convey("Given a JOIN selecting with a nested right wildcard", t, func() {
		tuples := getTuples(4)
		// rearrange the tuples
		for i, t := range tuples {
			if i%2 == 0 {
				t.InputName = "src1"
				t.Data["l"] = data.String(fmt.Sprintf("l%d", i))
				delete(t.Data, "int")
				t.Data["hoge"] = data.Int(i)
			} else {
				t.InputName = "src2"
				t.Data["r"] = data.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT RSTREAM src1:l, src2:* AS b FROM src1 [RANGE 2 TUPLES], src2 [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l0"),
							"b": data.Map{
								"r":   data.String("r1"),
								"int": data.Int(2),
							},
						})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l0"),
							"b": data.Map{
								"r":   data.String("r1"),
								"int": data.Int(2),
							},
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String("l2"),
							"b": data.Map{
								"r":   data.String("r1"),
								"int": data.Int(2),
							},
						})
					}
				})
			}
		})
	})

	Convey("Given a JOIN selecting from left and right with different ranges", t, func() {
		tuples := getTuples(8)
		// rearrange the tuples
		for i, t := range tuples {
			if i%2 == 0 {
				t.InputName = "src1"
				t.Data["l"] = data.String(fmt.Sprintf("l%d", i))
			} else {
				t.InputName = "src2"
				t.Data["r"] = data.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT RSTREAM src1:l, src2:r FROM src1 [RANGE 1 TUPLES], src2 [RANGE 5 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 { // l0
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 { // r1
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l0"),
							"r": data.String("r1"),
						})
					} else if idx == 2 { // l2
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l2"),
							"r": data.String("r1"),
						})
					} else if idx == 3 { // r3
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l2"),
							"r": data.String("r1"),
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String("l2"),
							"r": data.String("r3"),
						})
					} else if idx == 4 { // l4
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l4"),
							"r": data.String("r1"),
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String("l4"),
							"r": data.String("r3"),
						})
					} else if idx == 5 { // r5
						So(len(out), ShouldEqual, 3)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l4"),
							"r": data.String("r1"),
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String("l4"),
							"r": data.String("r3"),
						})
						So(out[2], ShouldResemble, data.Map{
							"l": data.String("l4"),
							"r": data.String("r5"),
						})
					} else if idx == 6 { // l6
						So(len(out), ShouldEqual, 3)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l6"),
							"r": data.String("r1"),
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String("l6"),
							"r": data.String("r3"),
						})
						So(out[2], ShouldResemble, data.Map{
							"l": data.String("l6"),
							"r": data.String("r5"),
						})
					} else if idx == 7 { // r7
						So(len(out), ShouldEqual, 3)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l6"),
							"r": data.String("r3"),
						})
						So(out[1], ShouldResemble, data.Map{
							"l": data.String("l6"),
							"r": data.String("r5"),
						})
						So(out[2], ShouldResemble, data.Map{
							"l": data.String("l6"),
							"r": data.String("r7"),
						})
					}
				})
			}
		})
	})

	Convey("Given a JOIN selecting from left and right with a join condition", t, func() {
		tuples := getTuples(8)
		// rearrange the tuples
		for i, t := range tuples {
			if i%2 == 0 {
				t.InputName = "src1"
				t.Data["l"] = data.String(fmt.Sprintf("l%d", i))
			} else {
				t.InputName = "src2"
				t.Data["r"] = data.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT ISTREAM src1:l, src2:r FROM src1 [RANGE 2 TUPLES], src2 [RANGE 2 TUPLES] ` +
			`WHERE src1:int + 1 = src2:int`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String("l0"), // int: 1
							"r": data.String("r1"), // int: 2
						})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 0)
					} else if idx%2 == 1 {
						// a tuple from src2 (=right) was just added
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String(fmt.Sprintf("l%d", idx-1)), // int: x
							"r": data.String(fmt.Sprintf("r%d", idx)),   // int: x+1
						})
					} else {
						// a tuple from src1 (=left) was just added
						So(len(out), ShouldEqual, 0)
					}
				})
			}
		})
	})

	Convey("Given a self-join with a join condition", t, func() {
		tuples := getTuples(8)
		// rearrange the tuples
		for i, t := range tuples {
			t.InputName = "src"
			t.Data["x"] = data.String(fmt.Sprintf("x%d", i))
		}
		s := `CREATE STREAM box AS SELECT ISTREAM src1:x AS l, src2:x AS r ` +
			`FROM src [RANGE 2 TUPLES] AS src1, src [RANGE 2 TUPLES] AS src2 ` +
			`WHERE src1:int + 1 = src2:int`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 {
						// join condition fails
						So(len(out), ShouldEqual, 0)
					} else {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String(fmt.Sprintf("x%d", idx-1)), // int: x
							"r": data.String(fmt.Sprintf("x%d", idx)),   // int: x+1
						})
					}
				})
			}
		})
	})

	Convey("Given a self-self-join with a join condition", t, func() {
		tuples := getTuples(8)
		// rearrange the tuples
		for i, t := range tuples {
			t.InputName = "src"
			t.Data["x"] = data.String(fmt.Sprintf("x%d", i))
		}
		s := `CREATE STREAM box AS SELECT ISTREAM src1:x AS l, src2:x AS r, src3:x AS x ` +
			`FROM src [RANGE 2 TUPLES] AS src1, src [RANGE 2 TUPLES] AS src2, src [RANGE 3 TUPLES] AS src3 ` +
			`WHERE src1:int + 1 = src2:int AND src2:int = src3:int + 1`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				// sort the output by the "l" and then the "r" key before
				// checking if it resembles the expected value
				sort.Sort(tupleList(out))

				Convey(fmt.Sprintf("Then joined values should appear in %v", idx), func() {
					if idx == 0 {
						// join condition fails
						So(len(out), ShouldEqual, 0)
					} else {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{
							"l": data.String(fmt.Sprintf("x%d", idx-1)), // int: x
							"r": data.String(fmt.Sprintf("x%d", idx)),   // int: x+1
							"x": data.String(fmt.Sprintf("x%d", idx-1)), // int: x
						})
					}
				})
			}
		})
	})
}

func createDefaultSelectPlan2(s string) (PhysicalPlan, error) {
	p := parser.New()
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	_stmt, _, err := p.ParseStmt(s)
	if err != nil {
		return nil, err
	}
	stmt := _stmt.(parser.CreateStreamAsSelectStmt).Select
	logicalPlan, err := Analyze(stmt, reg)
	if err != nil {
		return nil, err
	}
	canBuild := CanBuildDefaultSelectExecutionPlan(logicalPlan, reg)
	if !canBuild {
		err := fmt.Errorf("defaultSelectExecutionPlan cannot be used for statement: %s", s)
		return nil, err
	}
	return NewDefaultSelectExecutionPlan(logicalPlan, reg)
}

func BenchmarkNormalExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT ISTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 5 TUPLES]`
	plan, err := createDefaultSelectPlan2(s)
	if err != nil {
		panic(err.Error())
	}
	tmplTup := core.Tuple{
		Data:          data.Map{"int": data.Int(-1)},
		InputName:     "src",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	for n := 0; n < b.N; n++ {
		inTup := tmplTup.Copy()
		inTup.Data["int"] = data.Int(n)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkNormalExecutionBigWindow(b *testing.B) {
	s := `CREATE STREAM box AS SELECT ISTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 50 TUPLES]`
	plan, err := createDefaultSelectPlan2(s)
	if err != nil {
		panic(err.Error())
	}
	tmplTup := core.Tuple{
		Data:          data.Map{"int": data.Int(-1)},
		InputName:     "src",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	for n := 0; n < b.N; n++ {
		inTup := tmplTup.Copy()
		inTup.Data["int"] = data.Int(n)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkTimeBasedExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT ISTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 5 SECONDS]`
	plan, err := createDefaultSelectPlan2(s)
	if err != nil {
		panic(err.Error())
	}
	tmplTup := core.Tuple{
		Data:          data.Map{"int": data.Int(-1)},
		InputName:     "src",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	for n := 0; n < b.N; n++ {
		inTup := tmplTup.Copy()
		inTup.Data["int"] = data.Int(n)
		inTup.Timestamp = inTup.Timestamp.Add(time.Duration(n) * time.Second)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkTimeBasedExecutionBigWindow(b *testing.B) {
	s := `CREATE STREAM box AS SELECT ISTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 50 SECONDS]`
	plan, err := createDefaultSelectPlan2(s)
	if err != nil {
		panic(err.Error())
	}
	tmplTup := core.Tuple{
		Data:          data.Map{"int": data.Int(-1)},
		InputName:     "src",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	for n := 0; n < b.N; n++ {
		inTup := tmplTup.Copy()
		inTup.Data["int"] = data.Int(n)
		inTup.Timestamp = inTup.Timestamp.Add(time.Duration(n) * time.Second)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkWithWhere(b *testing.B) {
	s := `CREATE STREAM box AS SELECT ISTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 5 TUPLES] WHERE int % 2 = 0`
	plan, err := createDefaultSelectPlan2(s)
	if err != nil {
		panic(err.Error())
	}
	tmplTup := core.Tuple{
		Data:          data.Map{"int": data.Int(-1)},
		InputName:     "src",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	for n := 0; n < b.N; n++ {
		inTup := tmplTup.Copy()
		inTup.Data["int"] = data.Int(n)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkNormalJoin(b *testing.B) {
	s := `CREATE STREAM box AS SELECT ISTREAM left:int, right:int
	FROM src [RANGE 5 TUPLES] AS left, src [RANGE 5 TUPLES] AS right
	WHERE left:int - right:int < 2`
	plan, err := createDefaultSelectPlan2(s)
	if err != nil {
		panic(err.Error())
	}
	tmplTup := core.Tuple{
		Data:          data.Map{"int": data.Int(-1)},
		InputName:     "src",
		Timestamp:     time.Date(2015, time.April, 10, 10, 23, 0, 0, time.UTC),
		ProcTimestamp: time.Date(2015, time.April, 10, 10, 24, 0, 0, time.UTC),
		BatchID:       7,
	}
	for n := 0; n < b.N; n++ {
		inTup := tmplTup.Copy()
		inTup.Data["int"] = data.Int(n)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

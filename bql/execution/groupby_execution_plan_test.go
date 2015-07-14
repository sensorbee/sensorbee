package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func createGroupbyPlan(s string, t *testing.T) (ExecutionPlan, error) {
	p := parser.NewBQLParser()
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	reg.Register("udaf", &dummyAggregate{})
	_stmt, _, err := p.ParseStmt(s)
	if err != nil {
		return nil, err
	}
	So(_stmt, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
	stmt := _stmt.(parser.CreateStreamAsSelectStmt)
	logicalPlan, err := Analyze(stmt, reg)
	if err != nil {
		return nil, err
	}
	canBuild := CanBuildGroupbyExecutionPlan(logicalPlan, reg)
	if !canBuild {
		err := fmt.Errorf("groupByExecutionPlan cannot be used for statement: %s", s)
		return nil, err
	}
	return NewGroupbyExecutionPlan(logicalPlan, reg)
}

func getOtherTuples() []*core.Tuple {
	tuples := getTuples(4)
	tuples[0].Data["foo"] = data.Int(1)
	tuples[1].Data["foo"] = data.Int(1)
	tuples[2].Data["foo"] = data.Int(2)
	tuples[3].Data["foo"] = data.Int(2)
	return tuples
}

func TestGroupbyExecutionPlan(t *testing.T) {
	Convey("Given a SELECT clause with GROUP BY but no aggregation", t, func() {
		tuples := getOtherTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM foo FROM src [RANGE 3 TUPLES] GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2)})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with aggregation but no GROUP BY", t, func() {
		tuples := getOtherTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM count(foo) FROM src [RANGE 3 TUPLES]`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"count": data.Int(math.Min(float64(idx+1), 3))})
				})
			}
		})
	})

	Convey("Given a SELECT clause with aggregation but no GROUP BY on empty input", t, func() {
		tuples := getOtherTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM count(foo) FROM src [RANGE 3 TUPLES] WHERE foo=7`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"count": data.Int(0)})
				})
			}
		})
	})

	Convey("Given a SELECT clause with constant aggregation but no GROUP BY", t, func() {
		tuples := getOtherTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM count(17) FROM src [RANGE 3 TUPLES]`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"count": data.Int(math.Min(float64(idx+1), 3))})
				})
			}
		})
	})

	Convey("Given a SELECT clause with count(*) but no GROUP BY", t, func() {
		tuples := getOtherTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM count(*) FROM src [RANGE 3 TUPLES]`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"count": data.Int(math.Min(float64(idx+1), 3))})
				})
			}
		})
	})

	Convey("Given a SELECT clause with a simple aggregation and GROUP BY", t, func() {
		tuples := getOtherTuples()
		tuples[3].Data["int"] = data.Null{} // NULL should not be counted
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 3 TUPLES] GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "count": data.Int(1)})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
						So(out[1], ShouldResemble,
							// the below is just 1 because NULL isn't counted
							data.Map{"foo": data.Int(2), "count": data.Int(1)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a multiple-parameters aggregation (agg+const) and GROUP BY", t, func() {
		tuples := getOtherTuples()
		s := `CREATE STREAM box AS SELECT RSTREAM foo, udaf(int, 'hoge') FROM src [RANGE 3 TUPLES] GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`1+"hoge"`)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`2+"hoge"`)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`2+"hoge"`)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "udaf": data.String(`1+"hoge"`)})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`1+"hoge"`)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "udaf": data.String(`2+"hoge"`)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a multiple-parameters aggregation (agg+groupby) and GROUP BY", t, func() {
		tuples := getOtherTuples()
		s := `CREATE STREAM box AS SELECT RSTREAM foo, udaf(int, foo) FROM src [RANGE 3 TUPLES] GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`1+1`)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`2+1`)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`2+1`)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "udaf": data.String(`1+2`)})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "udaf": data.String(`1+1`)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "udaf": data.String(`2+2`)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a simple aggregation and GROUP BY and WHERE on the aggregated column", t, func() {
		tuples := getOtherTuples()
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 3 TUPLES] WHERE int%2=1 GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "count": data.Int(1)})
					} else {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(2), "count": data.Int(1)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a simple aggregation and GROUP BY and WHERE on the grouped column", t, func() {
		tuples := getOtherTuples()
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 3 TUPLES] WHERE foo=1 GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2)})
					} else {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a computation inside the aggregation", t, func() {
		tuples := getOtherTuples()
		// TODO we should not use "count" or we don't see the effect
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int+foo) FROM src [RANGE 3 TUPLES] GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "count": data.Int(1)})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "count": data.Int(2)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a computation after the aggregation", t, func() {
		tuples := getOtherTuples()
		// TODO we should not use "count" or we don't see the effect
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int)+foo AS count FROM src [RANGE 3 TUPLES] GROUP BY foo`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1 + 1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2 + 1)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(2 + 1)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "count": data.Int(1 + 2)})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "count": data.Int(1 + 1)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "count": data.Int(2 + 2)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with two columns in the GROUP BY clause", t, func() {
		tuples := getOtherTuples()
		tuples[0].Data["bar"] = data.Int(1)
		tuples[1].Data["bar"] = data.Int(1)
		tuples[2].Data["bar"] = data.Int(1)
		tuples[3].Data["bar"] = data.Int(2)
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) + 2 AS x FROM src [RANGE 3 TUPLES] GROUP BY foo, bar`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "x": data.Int(3)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "x": data.Int(4)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "x": data.Int(4)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "x": data.Int(3)})
					} else {
						So(len(out), ShouldEqual, 3)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "x": data.Int(3)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "x": data.Int(3)})
						So(out[2], ShouldResemble,
							data.Map{"foo": data.Int(2), "x": data.Int(3)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a JOIN and two columns in the GROUP BY clause", t, func() {
		tuples := getOtherTuples()
		s := `CREATE STREAM box AS SELECT RSTREAM count(b:foo) FROM
			src [RANGE 4 TUPLES] AS a, src [RANGE 4 TUPLES] AS b WHERE a:foo=1 GROUP BY b:int`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 4 {
						So(len(out), ShouldEqual, 4)
						So(out[0], ShouldResemble,
							data.Map{"count": data.Int(2)})
						So(out[1], ShouldResemble,
							data.Map{"count": data.Int(2)})
						So(out[2], ShouldResemble,
							data.Map{"count": data.Int(2)})
						So(out[3], ShouldResemble,
							data.Map{"count": data.Int(2)})
					}
				})
			}
		})
	})

	/// things that do not work

	// limitations of this plan (must use default plan)

	Convey("Given an SELECT statement without aggregation or GROUP BY", t, func() {
		// not actually an aggregation statement
		s := `CREATE STREAM box AS SELECT RSTREAM foo FROM src [RANGE 3 TUPLES]`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "groupByExecutionPlan cannot be used for statement")
	})

	// normal SQL GROUP BY limitations

	Convey("Given an SELECT statement with an invalid GROUP BY (1)", t, func() {
		// not referencing the group-by column
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 3 TUPLES]`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `column "src:foo" must appear in the GROUP BY clause or be used in an aggregate function`)
	})

	Convey("Given an SELECT statement with an invalid GROUP BY (2)", t, func() {
		// not referencing the group-by column
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int)+int AS count FROM src [RANGE 3 TUPLES] GROUP BY foo`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `column "src:int" must appear in the GROUP BY clause or be used in an aggregate function`)
	})

	// BQL limitations (less features than SQL)

	Convey("Given an SELECT statement with an invalid GROUP BY (3)", t, func() {
		// using an expression in the GROUP BY clause
		s := `CREATE STREAM box AS SELECT RSTREAM foo+2 FROM src [RANGE 3 TUPLES] GROUP BY foo+2`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `grouping by expressions is not supported yet`)
	})

	Convey("Given an SELECT statement with an invalid GROUP BY (4)", t, func() {
		// using an expression in the GROUP BY clause
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 3 TUPLES] GROUP BY ts()`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `grouping by expressions is not supported yet`)
	})

	Convey("Given an SELECT statement with HAVING", t, func() {
		// using HAVING
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 3 TUPLES] GROUP BY foo HAVING count=2`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "groupByExecutionPlan cannot be used for statement")
	})

	Convey("Given an SELECT statement with an unknown UDAF", t, func() {
		// using an unknown UDAF
		s := `CREATE STREAM box AS SELECT RSTREAM foo, unknownUDAF(int) FROM src [RANGE 3 TUPLES] GROUP BY foo`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "function 'unknownUDAF' is unknown")
	})

	Convey("Given an SELECT statement with a timestamp", t, func() {
		// using metadata in the SELECT clause (outside of aggregate)
		s := `CREATE STREAM box AS SELECT RSTREAM ts() FROM src [RANGE 3 TUPLES] GROUP BY foo`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `using metadata 'TS' in GROUP BY statements is not supported yet`)
	})
}

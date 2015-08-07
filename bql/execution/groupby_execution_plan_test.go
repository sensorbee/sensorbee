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
	"time"
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
	stmt := _stmt.(parser.CreateStreamAsSelectStmt).Select
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
	Convey("Given a SELECT clause with GROUP BY and HAVING but no aggregation", t, func() {
		tuples := getOtherTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM foo FROM src [RANGE 3 TUPLES] GROUP BY foo HAVING foo=1`
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
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1)})
					} else {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1)})
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

	Convey("Given a SELECT clause with aggregation and HAVING but no GROUP BY", t, func() {
		tuples := getOtherTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM count(foo) FROM src [RANGE 3 TUPLES] HAVING count(foo)%2=0`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					if idx == 1 {
						// count(foo) is 2
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, data.Map{"count": data.Int(2)})
					} else {
						// count(foo) is 1 or 3
						So(len(out), ShouldEqual, 0)
					}
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

	Convey("Given a SELECT clause with an structure of aggregations and GROUP BY", t, func() {
		tuples := getOtherTuples()
		tuples[3].Data["int"] = data.Null{} // NULL should not be counted
		s := `CREATE STREAM box AS SELECT RSTREAM foo, [count(int), max(int)] AS a,
			{'c': count(int), 'm': min(int)} AS b FROM src [RANGE 3 TUPLES] GROUP BY foo`
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
							data.Map{"foo": data.Int(1), "a": data.Array{data.Int(1), data.Int(1)},
								"b": data.Map{"c": data.Int(1), "m": data.Int(1)}})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "a": data.Array{data.Int(2), data.Int(2)},
								"b": data.Map{"c": data.Int(2), "m": data.Int(1)}})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "a": data.Array{data.Int(2), data.Int(2)},
								"b": data.Map{"c": data.Int(2), "m": data.Int(1)}})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "a": data.Array{data.Int(1), data.Int(3)},
								"b": data.Map{"c": data.Int(1), "m": data.Int(3)}})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "a": data.Array{data.Int(1), data.Int(2)},
								"b": data.Map{"c": data.Int(1), "m": data.Int(2)}})
						So(out[1], ShouldResemble,
							// the below is just 1 because NULL isn't counted
							data.Map{"foo": data.Int(2), "a": data.Array{data.Int(1), data.Int(3)},
								"b": data.Map{"c": data.Int(1), "m": data.Int(3)}})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with a simple aggregation and GROUP BY and HAVING", t, func() {
		tuples := getOtherTuples()
		tuples[2].Data["int"] = data.Null{}
		s := `CREATE STREAM box AS SELECT RSTREAM foo, max(int) FROM src [RANGE 3 TUPLES] GROUP BY foo HAVING count(*)=1`
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
							data.Map{"foo": data.Int(1), "max": data.Int(1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 2 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(2), "max": data.Null{}})
					} else {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "max": data.Int(2)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with two identical aggregations and GROUP BY", t, func() {
		tuples := getOtherTuples()
		tuples[3].Data["int"] = data.Null{} // NULL should not be counted
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) AS x, count(int) AS y FROM src [RANGE 3 TUPLES] GROUP BY foo`
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
							data.Map{"foo": data.Int(1), "x": data.Int(1), "y": data.Int(1)})
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "x": data.Int(2), "y": data.Int(2)})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "x": data.Int(2), "y": data.Int(2)})
						So(out[1], ShouldResemble,
							data.Map{"foo": data.Int(2), "x": data.Int(1), "y": data.Int(1)})
					} else {
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble,
							data.Map{"foo": data.Int(1), "x": data.Int(1), "y": data.Int(1)})
						So(out[1], ShouldResemble,
							// the below is just 1 because NULL isn't counted
							data.Map{"foo": data.Int(2), "x": data.Int(1), "y": data.Int(1)})
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

	Convey("Given an SELECT statement with an invalid GROUP BY (3)", t, func() {
		// not referencing the group-by column
		s := `CREATE STREAM box AS SELECT RSTREAM foo FROM src [RANGE 3 TUPLES] HAVING foo`
		_, err := createGroupbyPlan(s, t)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `column "src:foo" must appear in the GROUP BY clause or be used in an aggregate function`)
	})

	Convey("Given an SELECT statement with an invalid GROUP BY (4)", t, func() {
		// not referencing the group-by column
		s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 3 TUPLES] GROUP BY foo HAVING int=2`
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

func TestAggregateFunctions(t *testing.T) {
	getExtTuples := func() []*core.Tuple {
		tuples := getOtherTuples()
		tuples[0].Data["bar"] = data.String("a")
		tuples[1].Data["bar"] = data.String("b")
		tuples[2].Data["bar"] = data.String("c")
		tuples[3].Data["bar"] = data.String("d")
		return tuples
	}

	Convey("Given a SELECT clause with array_agg", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM array_agg(int) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Array{
							data.Int(2), data.Int(3), data.Int(4)}})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with avg", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM avg(int) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Float(3.0)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with bool_and", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM bool_and(int = 2) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Bool(false)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with bool_or", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM bool_or(int = 2) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Bool(true)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with json_object_agg", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM json_object_agg(bar, int) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Map{
							"b": data.Int(2), "c": data.Int(3), "d": data.Int(4)}})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with max", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM max(int) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Int(4)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with min", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM min(int) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Int(2)})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with string_agg", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM string_agg(bar, ', ') AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.String("b, c, d")})
					}
				})
			}
		})
	})

	Convey("Given a SELECT clause with sum", t, func() {
		tuples := getExtTuples()

		s := `CREATE STREAM box AS SELECT RSTREAM sum(int) AS result
			FROM src [RANGE 3 TUPLES] WHERE int > 1`
		plan, err := createGroupbyPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)

					if idx == 0 {
						So(out[0], ShouldResemble, data.Map{"result": data.Null{}})
					} else if idx == 3 {
						So(out[0], ShouldResemble, data.Map{"result": data.Int(9)})
					}
				})
			}
		})
	})
}

func createGroupbyPlan2(s string) (ExecutionPlan, error) {
	p := parser.NewBQLParser()
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	reg.Register("udaf", &dummyAggregate{})
	_stmt, _, err := p.ParseStmt(s)
	if err != nil {
		return nil, err
	}
	stmt := _stmt.(parser.CreateStreamAsSelectStmt).Select
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

// ca. 65000 ns/op
func BenchmarkGroupingExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 5 TUPLES] GROUP BY foo`
	plan, err := createGroupbyPlan2(s)
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
		inTup.Data["foo"] = data.Int(n % 2)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

// ca. 71000 ns/op
func BenchmarkGroupingTimeBasedExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 5 SECONDS] GROUP BY foo`
	plan, err := createGroupbyPlan2(s)
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
		inTup.Data["foo"] = data.Int(n % 2)
		inTup.Timestamp = inTup.Timestamp.Add(time.Duration(n) * time.Second)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

// ca. 240000 ns/op
func BenchmarkLargeGroupExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 50 TUPLES] GROUP BY foo`
	plan, err := createGroupbyPlan2(s)
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
		inTup.Data["foo"] = data.Int(n % 3)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

// ca. 250000 ns/op
func BenchmarkLargeGroupTimeBasedExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT RSTREAM foo, count(int) FROM src [RANGE 50 SECONDS] GROUP BY foo`
	plan, err := createGroupbyPlan2(s)
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
		inTup.Data["foo"] = data.Int(n % 3)
		inTup.Timestamp = inTup.Timestamp.Add(time.Duration(n) * time.Second)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

// ca. 110000 ns/op
func BenchmarkComplicatedGroupExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT RSTREAM foo, udaf(int, foo) FROM src [RANGE 10 TUPLES] GROUP BY foo`
	plan, err := createGroupbyPlan2(s)
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
		inTup.Data["foo"] = data.Int(n % 3)
		_, err := plan.Process(inTup)
		if err != nil {
			panic(err.Error())
		}
	}
}

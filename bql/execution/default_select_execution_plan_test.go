package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
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

func createDefaultSelectPlan(s string, t *testing.T) (ExecutionPlan, error) {
	p := parser.NewBQLParser()
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
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
		s := `CREATE STREAM box AS SELECT ISTREAM 3+4-6+1, 3.0*4/2+1=7, null FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then that constant should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							data.Map{"col_1": data.Int(2), "col_2": data.Bool(true), "col_3": data.Null{}})
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
		s := `CREATE STREAM box AS SELECT ISTREAM int FROM src [RANGE 2 SECONDS]`
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

	Convey("Given a SELECT clause with a column and timestamp", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int, now() FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					var prevTime data.Value = nil
					if idx == 0 {
						So(len(out), ShouldEqual, 1)
						So(out[0]["int"], ShouldEqual, data.Int(1))
						So(out[0]["now"], ShouldHaveSameTypeAs, data.Timestamp{})
						prevTime = out[0]["now"]
					} else if idx == 1 {
						So(len(out), ShouldEqual, 2)
						So(out[0]["int"], ShouldEqual, data.Int(1))
						So(out[0]["now"], ShouldHaveSameTypeAs, data.Timestamp{})
						So(out[1]["int"], ShouldEqual, data.Int(2))
						So(out[1]["now"], ShouldHaveSameTypeAs, data.Timestamp{})
						So(out[1]["now"], ShouldResemble, out[0]["now"])
						So(out[1]["now"], ShouldNotResemble, prevTime)
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
					if idx == 0 {
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

	Convey("Given a SELECT clause with a wildcard and an overriding column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM *, (int-1)*2 AS int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						data.Map{"int": data.Int(2 * idx)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with a column and an overriding wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM (int-1)*2 AS int, * FROM src [RANGE 2 SECONDS]`
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
					// the latter is broken, therefore the query fails.
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

	// RSTREAM/2 TUPLES window
	Convey("Given an RSTREAM emitter selecting a constant and a 2 SECONDS window", t, func() {
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

	Convey("Given an RSTREAM emitter selecting a column and a 2 SECONDS window", t, func() {
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
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 0)
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
				So(len(output[1]), ShouldEqual, 0)
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
}

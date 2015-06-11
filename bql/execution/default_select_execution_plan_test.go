package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
	"time"
)

func getTuples(num int) []*tuple.Tuple {
	tuples := make([]*tuple.Tuple, 0, num)
	for i := 0; i < num; i++ {
		tup := tuple.Tuple{
			Data: tuple.Map{
				"int": tuple.Int(i + 1),
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
	reg := udf.NewDefaultFunctionRegistry()
	_stmt, err := p.ParseStmt(s)
	So(err, ShouldBeNil)
	So(_stmt, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
	stmt := _stmt.(parser.CreateStreamAsSelectStmt)
	logicalPlan, err := Analyze(stmt)
	So(err, ShouldBeNil)
	canBuild := CanBuildDefaultSelectExecutionPlan(logicalPlan, reg)
	So(canBuild, ShouldBeTrue)
	return NewDefaultSelectExecutionPlan(logicalPlan, reg)
}

func TestDefaultSelectExecutionPlan(t *testing.T) {
	// Select constant
	Convey("Given a SELECT clause with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2 FROM src [RANGE 2 SECONDS]`
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
							tuple.Map{"col_1": tuple.Int(2)})
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
						tuple.Map{"int": tuple.Int(idx + 1)})
				})
			}

		})
	})

	Convey("Given a SELECT clause with only a column using the table name", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM src.int FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						tuple.Map{"int": tuple.Int(idx + 1)})
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
						tuple.Map{"col_1": tuple.Int(2), "int": tuple.Int(idx + 1)})
				})
			}

		})
	})

	// Select constant and a column with changing values from aliased relation
	Convey("Given a SELECT clause with a constant, a column, and a table alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2, int FROM src AS x [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						tuple.Map{"col_1": tuple.Int(2), "int": tuple.Int(idx + 1)})
				})
			}

		})
	})

	// Select constant and a column with changing values from aliased relation
	// using that alias
	Convey("Given a SELECT clause with a constant, a table alias, and a column using it", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2, x.int FROM src AS x [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
					So(len(out), ShouldEqual, 1)
					So(out[0], ShouldResemble,
						tuple.Map{"col_1": tuple.Int(2), "int": tuple.Int(idx + 1)})
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
						tuple.Map{"a": tuple.Int(idx), "b": tuple.Int(idx + 1)})
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
						tuple.Map{"int": tuple.Int(idx + 1)})
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
						tuple.Map{"int": tuple.Int(2 * idx)})
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
						tuple.Map{"int": tuple.Int(idx + 1)})
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
						tuple.Map{"col_2": tuple.Int(2 * idx), "x": tuple.Map{"int": tuple.Int(idx + 1)}})
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
							tuple.Map{"b": tuple.Int(idx + 1)})
					} else {
						So(len(out), ShouldEqual, 0)
					}
				})
			}

		})
	})

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
							tuple.Map{"int": tuple.Int(idx + 1)})
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
							tuple.Map{"int": tuple.Int(idx)})
						So(out[1], ShouldResemble,
							tuple.Map{"int": tuple.Int(idx + 1)})
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
							tuple.Map{"int": tuple.Int(idx + 1)})
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
							tuple.Map{"int": tuple.Int(idx)})
						So(out[1], ShouldResemble,
							tuple.Map{"int": tuple.Int(idx + 1)})
					})
				} else {
					// In later runs, we have recovered from the error in item 1
					// and emit one item per run as normal.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							tuple.Map{"int": tuple.Int(idx + 1)})
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
							tuple.Map{"int": tuple.Int(1)})
					})
				} else {
					// In later runs, we have recovered from the error in item 1
					// and emit one item per run as normal.
					Convey(fmt.Sprintf("Then those values should appear in %v", idx), func() {
						So(err, ShouldBeNil)
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble,
							tuple.Map{"int": tuple.Int(idx - 1)})
					})
				}
			}

		})
	})

	// RSTREAM/2 SECONDS
	Convey("Given an RSTREAM/2 SECONDS statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2 AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[1][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[2]), ShouldEqual, 3)
				So(output[2][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[2][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[2][2], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[3]), ShouldEqual, 3)
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[3][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[3][2], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
			})

		})
	})

	Convey("Given an RSTREAM/2 SECONDS statement with a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(output[1][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[2]), ShouldEqual, 3)
				So(output[2][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(output[2][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[2][2], ShouldResemble, tuple.Map{"a": tuple.Int(3)})
				So(len(output[3]), ShouldEqual, 3)
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[3][1], ShouldResemble, tuple.Map{"a": tuple.Int(3)})
				So(output[3][2], ShouldResemble, tuple.Map{"a": tuple.Int(4)})
			})

		})
	})

	// RSTREAM/2 TUPLES
	Convey("Given an RSTREAM/2 SECONDS statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2 AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[1][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[2]), ShouldEqual, 2)
				So(output[2][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[2][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[3]), ShouldEqual, 2)
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[3][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
			})

		})
	})

	Convey("Given an RSTREAM/2 SECONDS statement with a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then the whole window state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(len(output[1]), ShouldEqual, 2)
				So(output[1][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(output[1][1], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[2]), ShouldEqual, 2)
				So(output[2][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(output[2][1], ShouldResemble, tuple.Map{"a": tuple.Int(3)})
				So(len(output[3]), ShouldEqual, 2)
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(3)})
				So(output[3][1], ShouldResemble, tuple.Map{"a": tuple.Int(4)})
			})

		})
	})

	// ISTREAM/2 SECONDS
	Convey("Given an ISTREAM/2 SECONDS statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2 AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 0)
				So(len(output[3]), ShouldEqual, 0)
			})

		})
	})

	Convey("Given an ISTREAM/2 SECONDS statement with a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(len(output[1]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[2]), ShouldEqual, 1)
				So(output[2][0], ShouldResemble, tuple.Map{"a": tuple.Int(3)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(4)})
			})

		})
	})

	// ISTREAM/2 TUPLES
	Convey("Given an ISTREAM/2 TUPLES statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2 AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[1]), ShouldEqual, 0)
				So(len(output[2]), ShouldEqual, 0)
				So(len(output[3]), ShouldEqual, 0)
			})

		})
	})

	Convey("Given an ISTREAM/2 TUPLES statement with a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			Convey("Then new items in state should be emitted", func() {
				So(len(output), ShouldEqual, 4)
				So(len(output[0]), ShouldEqual, 1)
				So(output[0][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(len(output[1]), ShouldEqual, 1)
				So(output[1][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
				So(len(output[2]), ShouldEqual, 1)
				So(output[2][0], ShouldResemble, tuple.Map{"a": tuple.Int(3)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(4)})
			})

		})
	})

	// DSTREAM/2 SECONDS
	Convey("Given a DSTREAM/2 SECONDS statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM 2 AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
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

	Convey("Given a DSTREAM/2 SECONDS statement with a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM int AS a FROM src [RANGE 2 SECONDS]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
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
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
			})

		})
	})

	// DSTREAM/2 TUPLES
	Convey("Given a DSTREAM/2 TUPLES statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM 2 AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
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

	Convey("Given a DSTREAM/2 TUPLES statement with a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT DSTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
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
				So(output[2][0], ShouldResemble, tuple.Map{"a": tuple.Int(1)})
				So(len(output[3]), ShouldEqual, 1)
				So(output[3][0], ShouldResemble, tuple.Map{"a": tuple.Int(2)})
			})

		})
	})
}

func TestDefaultSelectExecutionPlanJoin(t *testing.T) {
	Convey("Given a JOIN selecting from left and right", t, func() {
		tuples := getTuples(8)
		// rearrange the tuples
		for i, t := range tuples {
			if i%2 == 0 {
				t.InputName = "src1"
				t.Data["l"] = tuple.String(fmt.Sprintf("l%d", i))
			} else {
				t.InputName = "src2"
				t.Data["r"] = tuple.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT ISTREAM src1.l, src2.r FROM src1, src2 [RANGE 2 TUPLES]`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then that constant should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, tuple.Map{
							"l": tuple.String("l0"),
							"r": tuple.String("r1"),
						})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, tuple.Map{
							"l": tuple.String("l2"),
							"r": tuple.String("r1"),
						})
					} else if idx%2 == 1 {
						// a tuple from src2 (=right) was just added
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, tuple.Map{
							"l": tuple.String(fmt.Sprintf("l%d", idx-3)),
							"r": tuple.String(fmt.Sprintf("r%d", idx)),
						})
						So(out[1], ShouldResemble, tuple.Map{
							"l": tuple.String(fmt.Sprintf("l%d", idx-1)),
							"r": tuple.String(fmt.Sprintf("r%d", idx)),
						})
					} else {
						// a tuple from src1 (=left) was just added
						So(len(out), ShouldEqual, 2)
						So(out[0], ShouldResemble, tuple.Map{
							"l": tuple.String(fmt.Sprintf("l%d", idx)),
							"r": tuple.String(fmt.Sprintf("r%d", idx-3)),
						})
						So(out[1], ShouldResemble, tuple.Map{
							"l": tuple.String(fmt.Sprintf("l%d", idx)),
							"r": tuple.String(fmt.Sprintf("r%d", idx-1)),
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
				t.Data["l"] = tuple.String(fmt.Sprintf("l%d", i))
			} else {
				t.InputName = "src2"
				t.Data["r"] = tuple.String(fmt.Sprintf("r%d", i))
			}
		}
		s := `CREATE STREAM box AS SELECT ISTREAM src1.l, src2.r FROM src1, src2 [RANGE 2 TUPLES] ` +
			`WHERE src1.int + 1 = src2.int`
		plan, err := createDefaultSelectPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then that constant should appear in %v", idx), func() {
					if idx == 0 {
						So(len(out), ShouldEqual, 0)
					} else if idx == 1 {
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, tuple.Map{
							"l": tuple.String("l0"), // int: 1
							"r": tuple.String("r1"), // int: 2
						})
					} else if idx == 2 {
						So(len(out), ShouldEqual, 0)
					} else if idx%2 == 1 {
						// a tuple from src2 (=right) was just added
						So(len(out), ShouldEqual, 1)
						So(out[0], ShouldResemble, tuple.Map{
							"l": tuple.String(fmt.Sprintf("l%d", idx-1)), // int: x
							"r": tuple.String(fmt.Sprintf("r%d", idx)),   // int: x+1
						})
					} else {
						// a tuple from src1 (=left) was just added
						So(len(out), ShouldEqual, 0)
					}
				})
			}
		})
	})
}

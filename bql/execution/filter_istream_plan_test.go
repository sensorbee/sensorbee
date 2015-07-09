package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
)

func createFilterIstreamPlan(s string, t *testing.T) (ExecutionPlan, ExecutionPlan, error) {
	p := parser.NewBQLParser()
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	_stmt, _, err := p.ParseStmt(s)
	So(err, ShouldBeNil)
	So(_stmt, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
	stmt := _stmt.(parser.CreateStreamAsSelectStmt)
	logicalPlan, err := Analyze(stmt)
	So(err, ShouldBeNil)
	canBuild := CanBuildFilterIstreamPlan(logicalPlan, reg)
	So(canBuild, ShouldBeTrue)
	refCanBuild := CanBuildDefaultSelectExecutionPlan(logicalPlan, reg)
	So(refCanBuild, ShouldBeTrue)
	plan, err := NewFilterIstreamPlan(logicalPlan, reg)
	if err != nil {
		return nil, nil, err
	}
	refPlan, err := NewDefaultSelectExecutionPlan(logicalPlan, reg)
	if err != nil {
		return nil, nil, err
	}
	return plan, refPlan, err
}

func TestFilterIstreamPlan(t *testing.T) {
	// Select constant
	Convey("Given a SELECT clause with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2 FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Select a column with changing values
	Convey("Given a SELECT clause with only a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Select the tuple's timestamp
	Convey("Given a SELECT clause with only the timestamp", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM ts() FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Select a non-existing column
	Convey("Given a SELECT clause with a non-existing column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM hoge FROM src [RANGE 2 TUPLES]`
		plan, _, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for _, inTup := range tuples {
				_, err := plan.Process(inTup.Copy())
				So(err, ShouldNotBeNil) // hoge not found
			}

		})
	})

	Convey("Given a SELECT clause with a non-existing column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM hoge + 1 FROM src [RANGE 2 TUPLES]`
		plan, _, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for _, inTup := range tuples {
				_, err := plan.Process(inTup.Copy())
				So(err, ShouldNotBeNil) // hoge not found
			}

		})
	})

	// Select constant and a column with changing values
	Convey("Given a SELECT clause with a constant and a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2, int FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Use alias
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int-1 AS a, int AS b FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Use wildcard
	Convey("Given a SELECT clause with a wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM * FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	Convey("Given a SELECT clause with a wildcard and an overriding column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM *, (int-1)*2 AS int FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	Convey("Given a SELECT clause with a column and an overriding wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM (int-1)*2 AS int, * FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	Convey("Given a SELECT clause with an aliased wildcard and an anonymous column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM * AS x, (int-1)*2 FROM src [RANGE 2 TUPLES]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Use a filter
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int AS b FROM src [RANGE 2 TUPLES]
            WHERE int % 2 = 0`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				refOut, err := refPlan.Process(inTup.Copy())
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Recovery from errors in tuples
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

	// ISTREAM/2 SECONDS: (not supported)

	// ISTREAM/2 TUPLES
	Convey("Given an ISTREAM/2 TUPLES statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM 2 AS a FROM src [RANGE 2 TUPLES]`
		plan, _, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			// this is not implemented correctly in FilterIstreamPlan
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

	Convey("Given an ISTREAM/2 TUPLES statement with a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM int AS a FROM src [RANGE 2 TUPLES]`
		plan, _, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]data.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup.Copy())
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
}

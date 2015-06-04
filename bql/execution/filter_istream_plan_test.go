package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
	"testing"
)

func createFilterIstreamPlan(s string, t *testing.T) (ExecutionPlan, ExecutionPlan, error) {
	p := parser.NewBQLParser()
	reg := udf.NewDefaultFunctionRegistry()
	_stmt, err := p.ParseStmt(s)
	So(err, ShouldBeNil)
	So(_stmt, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
	stmt := _stmt.(parser.CreateStreamAsSelectStmt)
	logicalPlan, err := Analyze(&stmt)
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
		s := `CREATE STREAM box AS SELECT ISTREAM(2) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		// TODO this is not implemented correctly in FilterIstreamPlan
		SkipConvey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Select a column with changing values
	Convey("Given a SELECT clause with only a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(int) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Select a non-existing column
	Convey("Given a SELECT clause with a non-existing column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(hoge) FROM src [RANGE 2 SECONDS]`
		plan, _, err := createFilterIstreamPlan(s, t)
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
		s := `CREATE STREAM box AS SELECT ISTREAM(hoge + 1) FROM src [RANGE 2 SECONDS]`
		plan, _, err := createFilterIstreamPlan(s, t)
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
		s := `CREATE STREAM box AS SELECT ISTREAM(2, int) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Use alias
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(int-1 AS a, int AS b) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Use wildcard
	Convey("Given a SELECT clause with a wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(*) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	Convey("Given a SELECT clause with a wildcard and an overriding column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(*, (int-1)*2 AS int) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	Convey("Given a SELECT clause with a column and an overriding wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM((int-1)*2 AS int, *) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	Convey("Given a SELECT clause with an aliased wildcard and an anonymous column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(* AS x, (int-1)*2) FROM src [RANGE 2 SECONDS]`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// Use a filter
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(int AS b) FROM src [RANGE 2 SECONDS] 
            WHERE int % 2 = 0`
		plan, refPlan, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			for idx, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					refOut, err := refPlan.Process(inTup)
					So(err, ShouldBeNil)
					So(out, ShouldResemble, refOut)
				})
			}

		})
	})

	// ISTREAM/2 SECONDS
	Convey("Given an ISTREAM/2 SECONDS statement with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT ISTREAM(2 AS a) FROM src [RANGE 2 SECONDS]`
		plan, _, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			// TODO this is not implemented correctly in FilterIstreamPlan
			SkipConvey("Then new items in state should be emitted", func() {
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
		s := `CREATE STREAM box AS SELECT ISTREAM(int AS a) FROM src [RANGE 2 SECONDS]`
		plan, _, err := createFilterIstreamPlan(s, t)
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
		s := `CREATE STREAM box AS SELECT ISTREAM(2 AS a) FROM src [RANGE 2 TUPLES]`
		plan, _, err := createFilterIstreamPlan(s, t)
		So(err, ShouldBeNil)

		Convey("When feeding it with tuples", func() {
			output := [][]tuple.Map{}
			for _, inTup := range tuples {
				out, err := plan.Process(inTup)
				So(err, ShouldBeNil)
				output = append(output, out)
			}

			// TODO this is not implemented correctly in FilterIstreamPlan
			SkipConvey("Then new items in state should be emitted", func() {
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
		s := `CREATE STREAM box AS SELECT ISTREAM(int AS a) FROM src [RANGE 2 TUPLES]`
		plan, _, err := createFilterIstreamPlan(s, t)
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
}

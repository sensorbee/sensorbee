package execution

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"testing"
	"time"
)

func createFilterPlan(s string, t *testing.T) (PhysicalPlan, PhysicalPlan, error) {
	p := parser.New()
	reg := udf.CopyGlobalUDFRegistry(core.NewContext(nil))
	_stmt, _, err := p.ParseStmt(s)
	So(err, ShouldBeNil)
	So(_stmt, ShouldHaveSameTypeAs, parser.CreateStreamAsSelectStmt{})
	stmt := _stmt.(parser.CreateStreamAsSelectStmt).Select
	logicalPlan, err := Analyze(stmt, reg)
	So(err, ShouldBeNil)
	canBuild := CanBuildFilterPlan(logicalPlan, reg)
	So(canBuild, ShouldBeTrue)
	refCanBuild := CanBuildDefaultSelectExecutionPlan(logicalPlan, reg)
	So(refCanBuild, ShouldBeTrue)
	plan, err := NewFilterPlan(logicalPlan, reg)
	if err != nil {
		return nil, nil, err
	}
	refPlan, err := NewDefaultSelectExecutionPlan(logicalPlan, reg)
	if err != nil {
		return nil, nil, err
	}
	return plan, refPlan, err
}

func compareWithRef(t *testing.T, plan, refPlan PhysicalPlan, tuples []*core.Tuple) {
	Convey("When feeding it with tuples", func() {
		for idx, inTup := range tuples {
			out, err := plan.Process(inTup.Copy())
			refOut, refErr := refPlan.Process(inTup.Copy())
			if err != nil {
				// if our plan fails, so should the reference plan
				So(err, ShouldResemble, refErr)
			} else {
				Convey(fmt.Sprintf("Then the result should match the reference in %v", idx), func() {
					So(out, ShouldResemble, refOut)
				})
			}
		}

	})
}

func TestFilterPlan(t *testing.T) {
	// Select constant
	Convey("Given a SELECT clause with a constant", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
            null, [2.0,3] = [2,3.0] FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		_ = refPlan
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Select a column with changing values
	Convey("Given a SELECT clause with only a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int::string AS int FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with a column and timestamp", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int, distance_us(now(), now()), distance_us(clock_timestamp(), clock_timestamp()) < 100 AS t FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with only a column using the table name", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM src:int FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Select the tuple's timestamp
	Convey("Given a SELECT clause with only the timestamp", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM ts() FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Select a non-existing column
	Convey("Given a SELECT clause with a non-existing column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM hoge FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with a non-existing column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM hoge + 1 FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Select constant and a column with changing values
	Convey("Given a SELECT clause with a constant and a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2, int FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with an operation on a column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 0 - int FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Select constant and a column with changing values from aliased relation
	Convey("Given a SELECT clause with a constant, a column, and a table alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2, int FROM src [RANGE 1 TUPLES] AS x`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Select NULL-related operations
	Convey("Given a SELECT clause with NULL operations", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM null IS NULL, null + 2 = 2 FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with NULL filter", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int FROM src [RANGE 1 TUPLES] WHERE null`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Select constant and a column with changing values from aliased relation
	// using that alias
	Convey("Given a SELECT clause with a constant, a table alias, and a column using it", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2, x:int FROM src [RANGE 1 TUPLES] AS x`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Use alias
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int-1 AS a, int AS b FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with a nested column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int-1 AS a.c, int+1 AS a['d'], int AS b[1] FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Use wildcard
	Convey("Given a SELECT clause with a wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM * FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with a wildcard using stream prefix", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM src:* FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with a wildcard and an overriding column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM *, (int-1)*2 AS int FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with a stream wildcard and an overriding column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM src:*, (src:int-1)*2 AS int FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with a column and an overriding wildcard", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM (int-1)*2 AS int, * FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with an aliased wildcard and an anonymous column", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM * AS x, (int-1)*2 FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with an aliased wildcard and an anonymous column (both prefixed)", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM src:* AS x, (src:int-1)*2 FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a SELECT clause with an aliased wildcard and a prefixed anonymous column (both prefixed)", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM * AS x, (src:int-1)*2 FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// Use a filter
	Convey("Given a SELECT clause with a column alias", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS b FROM src [RANGE 1 TUPLES]
		            WHERE int % 2 = 0`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})
}

func TestFilterPlanEmitters(t *testing.T) {
	// Recovery from errors in tuples
	Convey("Given a SELECT clause with a column that does not exist in one tuple (RSTREAM)", t, func() {
		tuples := getTuples(6)
		// remove the selected key from one tuple
		delete(tuples[1].Data, "int")

		s := `CREATE STREAM box AS SELECT RSTREAM int FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given a WHERE clause with a column that does not exist in one tuple (RSTREAM)", t, func() {
		tuples := getTuples(6)
		// remove the selected key from one tuple
		delete(tuples[1].Data, "int")

		s := `CREATE STREAM box AS SELECT RSTREAM int FROM src [RANGE 1 TUPLES] WHERE int > 0`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	// RSTREAM/1 TUPLES window
	Convey("Given an RSTREAM emitter selecting a constant and a 1 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM 2 AS a FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

	Convey("Given an RSTREAM emitter selecting a column and a 1 TUPLES window", t, func() {
		tuples := getTuples(4)
		s := `CREATE STREAM box AS SELECT RSTREAM int AS a FROM src [RANGE 1 TUPLES]`
		plan, refPlan, err := createFilterPlan(s, t)
		So(err, ShouldBeNil)

		compareWithRef(t, plan, refPlan, tuples)
	})

}

func createFilterPlan2(s string) (PhysicalPlan, error) {
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
	canBuild := CanBuildFilterPlan(logicalPlan, reg)
	if !canBuild {
		err := fmt.Errorf("filterPlan cannot be used for statement: %s", s)
		return nil, err
	}
	return NewFilterPlan(logicalPlan, reg)
}

func BenchmarkFilterExecution(b *testing.B) {
	s := `CREATE STREAM box AS SELECT RSTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 1 TUPLES]`
	plan, err := createFilterPlan2(s)
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

func BenchmarkFilterWithWhere(b *testing.B) {
	s := `CREATE STREAM box AS SELECT RSTREAM cast(3+4-6+1 as float), 3.0::int*4/2+1=7.0,
			null, [2.0,3] = [2,3.0] FROM src [RANGE 1 TUPLES] WHERE int % 2 = 0`
	plan, err := createFilterPlan2(s)
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

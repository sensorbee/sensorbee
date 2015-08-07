package parser

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParserErrorMessage(t *testing.T) {
	testCases := map[string]string{
		// rubbish
		`HELLO`: `failed to parse string as BQL statement
statement has an unlocatable syntax error`,
		// wrong keyword -> unrecognizable statement
		`CREATE STRAEM x AS SELECT ISTREAM x`: `failed to parse string as BQL statement
statement has a syntax error near line 1, symbol 8:
  CREATE STRAEM x AS SELECT ISTREAM x
         ^`,
		`REWIND SSOURCE ab`: `statement has a syntax error near line 1, symbol 8:
  REWIND SSOURCE ab
         ^`,
		// some remainder after the statement is complete
		`REWIND SOURCE ab cd`: `statement has a syntax error near line 1, symbol 17:
  REWIND SOURCE ab cd
                  ^
consider to look up the documentation for RewindSourceStmt`,
		// unicode characters
		`SELECT ISTREAM '日本語', b FROM c [RANGE 3 UPLES]`: `failed to parse string as BQL statement
statement has a syntax error near line 1, symbol 24:
  ...ECT ISTREAM '日本語', b FROM c [RANGE 3 UPLES]
                            ^
consider to look up the documentation for SelectStmt`,
		// wrong expression
		`CREATE STREAM x AS SELECT ISTREAM f(x:*)`: `failed to parse string as BQL statement
statement has a syntax error near line 1, symbol 36:
  ... AS SELECT ISTREAM f(x:*)
                         ^
consider to look up the documentation for CreateStreamAsSelectStmt`,
		`INSERT INTO x SELECT ISTREAM 'Peter's'`: `failed to parse string as BQL statement
statement has a syntax error near line 1, symbol 37:
  ...LECT ISTREAM 'Peter's'
                         ^
consider to look up the documentation for InsertIntoSelectStmt`,
	}

	Convey("Given a BQL parser", t, func() {
		p := New()

		for stmt, expected := range testCases {
			// avoid closure over loop variables
			stmt, expected := stmt, expected

			Convey(fmt.Sprintf("When parsing %s", stmt), func() {
				_, _, err := p.ParseStmt(stmt)

				Convey("Then parsing should fail", func() {
					So(err, ShouldNotBeNil)

					Convey("And the error message should match", func() {
						So(err.Error(), ShouldEndWith, expected)
					})
				})
			})
		}

	})

}

package parser

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMultipleStmtParser(t *testing.T) {
	testCases := map[string][]interface{}{
		// empty statement
		"   ": []interface{}{},
		" ; ": nil,
		// single statement
		"SELECT RSTREAM a": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement, terminated with semicolon
		"SELECT RSTREAM a;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement, terminated with semicolon and some spaces
		"SELECT RSTREAM a ; ": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement starting with a space
		"  SELECT RSTREAM a ;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// two statements with various space separations
		"SELECT RSTREAM a;SELECT RSTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT RSTREAM a ;SELECT RSTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT RSTREAM a; SELECT RSTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT RSTREAM a ; SELECT RSTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		// three statements
		"SELECT RSTREAM a ; SELECT RSTREAM b; SELECT RSTREAM c": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"}}}},
		},
		// using semi-colons within the statements
		"SELECT RSTREAM a ; SELECT RSTREAM b; SELECT RSTREAM c, ';'": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"},
					StringLiteral{";"}}}},
		},
		"SELECT RSTREAM a;SELECT RSTREAM c, ';';SELECT RSTREAM b;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"},
					StringLiteral{";"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
	}

	Convey("Given a BQL parser", t, func() {
		p := NewBQLParser()

		for input, expected := range testCases {
			// avoid closure over loop variables
			input, expected := input, expected

			Convey(fmt.Sprintf("When parsing %s", input), func() {
				results, err := p.ParseStmts(input)

				Convey(fmt.Sprintf("Then the result should be %v", expected), func() {
					if expected == nil {
						So(err, ShouldNotBeNil)
					} else {
						// check there is no error
						So(err, ShouldBeNil)
						// check we go what we expected
						So(results, ShouldResemble, expected)
					}
				})
			})
		}

	})

}

func TestSingleStmtParser(t *testing.T) {
	testCases := map[string]string{
		// single statement
		"SELECT RSTREAM a": "",
		// single statement, terminated with semicolon
		"SELECT RSTREAM a;": "",
		// single statement, terminated with semicolon and some spaces
		"SELECT RSTREAM a ; ": "",
		// single statement starting with a space
		"  SELECT RSTREAM a ;": "",
		// two statements with various space separations
		"SELECT RSTREAM a;SELECT RSTREAM b":   "SELECT RSTREAM b",
		"SELECT RSTREAM a ;SELECT RSTREAM b":  "SELECT RSTREAM b",
		"SELECT RSTREAM a; SELECT RSTREAM b":  "SELECT RSTREAM b",
		"SELECT RSTREAM a ; SELECT RSTREAM b": "SELECT RSTREAM b",
		// three statements
		"SELECT RSTREAM a ; SELECT RSTREAM b; SELECT RSTREAM c": "SELECT RSTREAM b; SELECT RSTREAM c",
		// using semi-colons within the statements
		"SELECT RSTREAM a ; SELECT RSTREAM b; SELECT RSTREAM c, ';'": "SELECT RSTREAM b; SELECT RSTREAM c, ';'",
		"SELECT RSTREAM a;SELECT RSTREAM c, ';';SELECT RSTREAM b;":   "SELECT RSTREAM c, ';';SELECT RSTREAM b;",
	}

	Convey("Given a BQL parser", t, func() {
		p := NewBQLParser()

		for input, expected := range testCases {
			// avoid closure over loop variables
			input, expected := input, expected

			Convey(fmt.Sprintf("When parsing %s", input), func() {
				_, rest, err := p.ParseStmt(input)

				Convey(fmt.Sprintf("Then the result should be %v", expected), func() {
					// check there is no error
					So(err, ShouldBeNil)
					// check we go what we expected
					So(rest, ShouldEqual, expected)
				})
			})
		}

	})

}

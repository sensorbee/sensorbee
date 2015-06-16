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
		"SELECT a": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement, terminated with semicolon
		"SELECT a;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement, terminated with semicolon and some spaces
		"SELECT a ; ": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement starting with a space
		"  SELECT a ;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// two statements with various space separations
		"SELECT a;SELECT b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT a ;SELECT b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT a; SELECT b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT a ; SELECT b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		// three statements
		"SELECT a ; SELECT b; SELECT c": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"}}}},
		},
		// using semi-colons within the statements
		"SELECT a ; SELECT b; SELECT c, ';'": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
			SelectStmt{EmitterAST: EmitterAST{Rstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"},
					StringLiteral{";"}}}},
		},
		"SELECT a;SELECT c, ';';SELECT b;": []interface{}{
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
		"SELECT a": "",
		// single statement, terminated with semicolon
		"SELECT a;": "",
		// single statement, terminated with semicolon and some spaces
		"SELECT a ; ": "",
		// single statement starting with a space
		"  SELECT a ;": "",
		// two statements with various space separations
		"SELECT a;SELECT b":   "SELECT b",
		"SELECT a ;SELECT b":  "SELECT b",
		"SELECT a; SELECT b":  "SELECT b",
		"SELECT a ; SELECT b": "SELECT b",
		// three statements
		"SELECT a ; SELECT b; SELECT c": "SELECT b; SELECT c",
		// using semi-colons within the statements
		"SELECT a ; SELECT b; SELECT c, ';'": "SELECT b; SELECT c, ';'",
		"SELECT a;SELECT c, ';';SELECT b;":   "SELECT c, ';';SELECT b;",
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

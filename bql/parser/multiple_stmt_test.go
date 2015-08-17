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
		"SELECT ISTREAM a": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement, terminated with semicolon
		"SELECT ISTREAM a;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement, terminated with semicolon and some spaces
		"SELECT ISTREAM a ; ": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement starting with a space
		"  SELECT ISTREAM a ;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// two statements with various space separations
		"SELECT ISTREAM a;SELECT ISTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT ISTREAM a ;SELECT ISTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT ISTREAM a; SELECT ISTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		"SELECT ISTREAM a ; SELECT ISTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		// three statements
		"SELECT ISTREAM a ; SELECT ISTREAM b; SELECT ISTREAM c": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"}}}},
		},
		// using semi-colons within the statements
		"SELECT ISTREAM a ; SELECT ISTREAM b; SELECT ISTREAM c, ';'": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"},
					StringLiteral{";"}}}},
		},
		"SELECT ISTREAM a;SELECT ISTREAM c, ';';SELECT ISTREAM b;": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "c"},
					StringLiteral{";"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
	}

	Convey("Given a BQL parser", t, func() {
		p := New()

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

func TestComment(t *testing.T) {
	testCases := map[string][]interface{}{
		// single statement
		"SELECT ISTREAM a": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement on two lines
		"SELECT ISTREAM \na": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement with a blank line following
		"SELECT ISTREAM a\n": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement with a comment following on a new line
		"SELECT DSTREAM a\n--comment": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Dstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single semi-colon terminated statement with a comment following on a new line
		"SELECT DSTREAM a;\n--comment": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Dstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement with a comment following on a new line and more newlines
		"SELECT DSTREAM a\n--comment\n\n": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Dstream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement with an empty comment on two lines
		"SELECT ISTREAM --\na": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement with a space-only comment on two lines
		"SELECT ISTREAM --   \na": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement with a normal comment on two lines
		"SELECT ISTREAM -- comment here\na": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement with multiple comments
		"SELECT ISTREAM -- comment\n  -- continues\n--even longer\na": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// single statement introduced by comment
		" -- comment\nSELECT ISTREAM a": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
		},
		// multiple statements separated by comment
		"SELECT ISTREAM a;\n--comment\nSELECT ISTREAM b": []interface{}{
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}}}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "b"}}}},
		},
		// non-select statements as well
		("-- do some setup\nCREATE STATE hoge TYPE test;\nSELECT ISTREAM\n  --cols\n" +
			"  a,b;\nDROP STATE hoge;\n--done"): []interface{}{
			CreateStateStmt{StreamIdentifier("hoge"), SourceSinkType("test"), SourceSinkSpecsAST{nil}},
			SelectStmt{EmitterAST: EmitterAST{Istream, nil},
				ProjectionsAST: ProjectionsAST{[]Expression{RowValue{"", "a"}, RowValue{"", "b"}}}},
			DropStateStmt{StreamIdentifier("hoge")},
		},
	}

	Convey("Given a BQL parser", t, func() {
		p := New()

		for input, expected := range testCases {
			// avoid closure over loop variables
			input, expected := input, expected

			Convey(fmt.Sprintf("When parsing <%s>", input), func() {
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
		"SELECT ISTREAM a": "",
		// single statement, terminated with semicolon
		"SELECT ISTREAM a;": "",
		// single statement, terminated with semicolon and some spaces
		"SELECT ISTREAM a ; ": "",
		// single statement starting with a space
		"  SELECT ISTREAM a ;": "",
		// two statements with various space separations
		"SELECT ISTREAM a;SELECT ISTREAM b":   "SELECT ISTREAM b",
		"SELECT ISTREAM a ;SELECT ISTREAM b":  "SELECT ISTREAM b",
		"SELECT ISTREAM a; SELECT ISTREAM b":  "SELECT ISTREAM b",
		"SELECT ISTREAM a ; SELECT ISTREAM b": "SELECT ISTREAM b",
		// three statements
		"SELECT ISTREAM a ; SELECT ISTREAM b; SELECT ISTREAM c": "SELECT ISTREAM b; SELECT ISTREAM c",
		// using semi-colons within the statements
		"SELECT ISTREAM a ; SELECT ISTREAM b; SELECT ISTREAM c, ';'": "SELECT ISTREAM b; SELECT ISTREAM c, ';'",
		"SELECT ISTREAM a;SELECT ISTREAM c, ';';SELECT ISTREAM b;":   "SELECT ISTREAM c, ';';SELECT ISTREAM b;",
	}

	Convey("Given a BQL parser", t, func() {
		p := New()

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

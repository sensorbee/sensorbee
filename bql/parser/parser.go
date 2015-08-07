package parser

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
	"unicode"
)

type bqlParser struct {
	b bqlPeg
}

func New() *bqlParser {
	return &bqlParser{}
}

func (p *bqlParser) ParseStmt(s string) (result interface{}, rest string, err error) {
	// catch any parser errors
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error in BQL parser: %v", r)
		}
	}()
	// parse the statement
	b := p.b
	b.Buffer = s
	b.Init()
	if err := b.Parse(); err != nil {
		return nil, "", err
	}
	b.Execute()
	if b.parseStack.Peek() == nil {
		// the statement was parsed ok, but not put on the stack?
		// this should never occur.
		return nil, "", fmt.Errorf("no valid BQL statement could be parsed")
	}
	stackElem := b.parseStack.Pop()
	// we look at the part of the string right of the parsed
	// statement. note that we expect that trailing whitespace
	// or comments are already included in the range [0:stackElem.end]
	// as done by IncludeTrailingWhitespace() so that we do not
	// return a comment-only string as rest.
	isSpaceOrSemicolon := func(r rune) bool {
		return unicode.IsSpace(r) || r == rune(';')
	}
	rest = strings.TrimLeftFunc(string([]rune(s)[stackElem.end:]), isSpaceOrSemicolon)
	// pop it from the parse stack
	return stackElem.comp, rest, nil
}

func (p *bqlParser) ParseStmts(s string) ([]interface{}, error) {
	// parse all statements
	results := make([]interface{}, 0)
	rest := strings.TrimSpace(s)
	for rest != "" {
		result, rest_, err := p.ParseStmt(rest)
		if err != nil {
			return nil, err
		}
		// append the parsed statement to the result list
		results = append(results, result)
		rest = rest_
	}
	return results, nil
}

type bqlPeg struct {
	bqlPegBackend
}

func (b *bqlPeg) Parse(rule ...int) error {
	// override the Parse method from the bqlPegBackend in order
	// to place our own error before returning
	if err := b.bqlPegBackend.Parse(rule...); err != nil {
		if pErr, ok := err.(*parseError); ok {
			return &bqlParseError{pErr}
		}
		return err
	}
	return nil
}

type bqlParseError struct {
	*parseError
}

func (e *bqlParseError) Error() string {
	error := "failed to parse string as BQL statement\n"
	stmt := []rune(e.p.Buffer)
	tokens := e.p.tokenTree.Error()
	// collect the current stack of tokens and translate their
	// string indexes into line/symbol pairs
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	// now find the offensive line
	foundError := false
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		if end == 0 {
			// these are '' matches we cannot exploit for a useful error message
			continue
		} else if foundError {
			// if we found an error, the next tokens may give some additional
			// information about what kind of statement we have here. the first
			// rule that starts at 0 is (often?) the description we want.
			ruleName := rul3s[token.pegRule]
			if begin == 0 && end > 0 {
				error += fmt.Sprintf("\nconsider to look up the documentation for %s",
					ruleName)
				break
			}
		} else if end > 0 {
			error += fmt.Sprintf("statement has a syntax error near line %v, symbol %v:\n",
				translations[end].line, translations[end].symbol)
			// we want some output like:
			//
			//   ... FROM x [RANGE 7 UPLES] WHERE ...
			//                       ^
			//
			snipStartIdx := end - 20
			snipStart := "..."
			if snipStartIdx < 0 {
				snipStartIdx = 0
				snipStart = ""
			}
			snipEndIdx := end + 30
			snipEnd := "..."
			if snipEndIdx > len(stmt) {
				snipEndIdx = len(stmt)
				snipEnd = ""
			}
			// first line: an excerpt from the statement
			error += "  " + snipStart
			snipBeforeErr := strings.Replace(string(stmt[snipStartIdx:end]), "\n", " ", -1)
			snipAfterInclErr := strings.Replace(string(stmt[end:snipEndIdx]), "\n", " ", -1)
			error += snipBeforeErr + snipAfterInclErr
			error += snipEnd + "\n"
			// second line: a ^ marker at the correct position
			error += strings.Repeat(" ", len(snipStart)+2)
			error += strings.Repeat(" ", runewidth.StringWidth(snipBeforeErr))
			error += "^"
			foundError = true
		}
	}
	if !foundError {
		error += "statement has an unlocatable syntax error"
	}

	return error
}

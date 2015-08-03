package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type bqlParser struct {
	b bqlPeg
}

func NewBQLParser() *bqlParser {
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

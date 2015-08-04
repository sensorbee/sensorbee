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
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

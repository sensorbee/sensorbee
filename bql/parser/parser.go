package parser

import (
	"fmt"
)

type bqlParser struct {
	b bqlPeg
}

func NewBQLParser() *bqlParser {
	return &bqlParser{}
}

func (p *bqlParser) ParseStmt(s string) (res interface{}, err error) {
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
		return nil, err
	}
	b.Execute()
	if b.parseStack.Peek() == nil {
		// the statement was parsed ok, but not put on the stack?
		// this should never occur.
		return nil, fmt.Errorf("no valid BQL statement could be parsed")
	}
	// pop it from the parse stack
	return b.parseStack.Pop().comp, nil
}

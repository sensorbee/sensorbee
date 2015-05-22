package parser

import (
	"fmt"
)

func (b *Bql) ParseStmt(s string) (interface{}, error) {
	// parse the statement
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

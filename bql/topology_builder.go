package bql

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
)

type TopologyBuilder struct {
	core.StaticTopologyBuilder
}

func NewTopologyBuilder() *TopologyBuilder {
	tb := TopologyBuilder{core.NewDefaultStaticTopologyBuilder()}
	return &tb
}

func (tb *TopologyBuilder) BQL(s string) error {
	p := parser.NewBQLParser()
	_stmt, err := p.ParseStmt(s)
	if err != nil {
		return err
	}
	// check the type of statement
	switch stmt := _stmt.(type) {
	case parser.CreateStreamStmt:
		// the string "x" in "CREATE STREAM x AS" will be the box name
		outName := stmt.Name
		box := NewBqlBox(&stmt)
		// add all the referenced relations as named inputs
		decl := tb.AddBox(outName, box)
		for _, rel := range stmt.Relations {
			decl = decl.NamedInput(rel.Name, rel.Name)
		}
		return decl.Err()

	}
	return fmt.Errorf("statement of type %T is unimplemented", _stmt)
}

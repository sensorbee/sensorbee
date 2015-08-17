package exp

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/execution"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sort"
)

// Statement is a wrapper of BQL statement structs.
type Statement struct {
	Stmt interface{}
}

// NewStatement creates a Statement from a BQL statement. Currently, following
// statements are not supported:
//
//	* SELECT (CREATE STREAM is supported)
//	* UPDATE SOURCE/SINK/STATE
//	* DROP SOURCE/STREAM/SINK/STATE
//	* SAVE/LOAD STATE (they are automatically saved)
func NewStatement(s interface{}) (*Statement, error) {
	switch s.(type) {
	case parser.SelectStmt, parser.SelectUnionStmt,
		parser.UpdateSourceStmt, parser.UpdateSinkStmt, parser.UpdateStateStmt,
		parser.DropSourceStmt, parser.DropStreamStmt, parser.DropSinkStmt, parser.DropStateStmt,
		parser.SaveStateStmt, parser.LoadStateStmt, parser.LoadStateOrCreateStmt:
		return nil, fmt.Errorf("statement isn't supported by sensorbee exp command: %v", s)
	}

	return &Statement{
		Stmt: s,
	}, nil
}

// Input returns input node names of SELECT statements. It doesn't support Sinks.
func (s *Statement) Input() ([]string, error) {
	var (
		names []string
		err   error
	)
	switch stmt := s.Stmt.(type) {
	case parser.CreateStreamAsSelectStmt:
		names, err = inputFromSelect(&stmt.Select)
	case parser.CreateStreamAsSelectUnionStmt:
		for _, s := range stmt.Selects {
			ns, err := inputFromSelect(&s)
			if err != nil {
				return nil, err
			}
			names = append(names, ns...)
		}
	case parser.InsertIntoFromStmt:
		return []string{string(stmt.Input)}, nil
	case parser.InsertIntoSelectStmt:
		names, err = inputFromSelect(&stmt.SelectStmt)
	}
	if err != nil {
		return nil, err
	}
	return uniqueInputNames(names), nil
}

func inputFromSelect(stmt *parser.SelectStmt) ([]string, error) {
	var names []string
	for _, rel := range stmt.Relations {
		switch rel.Type {
		case parser.ActualStream:
			names = append(names, rel.Name)

		case parser.UDSFStream:
			ctx := core.NewContext(nil)
			udfReg := udf.CopyGlobalUDFRegistry(ctx)
			params := make([]data.Value, len(rel.Params))
			for i, expr := range rel.Params {
				p, err := execution.EvaluateFoldable(expr, udfReg)
				if err != nil {
					return nil, err
				}
				params[i] = p
			}

			reg, err := udf.CopyGlobalUDSFCreatorRegistry()
			if err != nil {
				return nil, err
			}

			udsfc, err := reg.Lookup(rel.Name, len(params))
			if err != nil {
				return nil, err
			}

			decl := udf.NewUDSFDeclarer()
			udsf, err := udsfc.CreateUDSF(ctx, decl, params...)
			if err != nil {
				return nil, err
			}
			defer udsf.Terminate(ctx)

			for n := range decl.ListInputs() {
				names = append(names, n)
			}
		}
	}
	return names, nil
}

func uniqueInputNames(names []string) []string {
	if len(names) == 0 {
		return names
	}

	sort.Sort(sort.StringSlice(names))
	prev := ""
	k := 0
	for _, n := range names {
		if n == prev {
			continue
		}

		names[k] = n
		k++
		prev = n
	}
	return names[:k]
}

// NodeName returns the name of node that the statement creates or manipulates.
// If the statement doesn't create a node or the statement creates a sink, it
// returns an empty string. Instead, INSERT INTO statements returns a name
// of a sink.
func (s *Statement) NodeName() string {
	switch stmt := s.Stmt.(type) {
	case parser.CreateSourceStmt:
		return string(stmt.Name)
	case parser.CreateStreamAsSelectStmt:
		return string(stmt.Name)
	case parser.CreateStreamAsSelectUnionStmt:
		return string(stmt.Name)
	case parser.InsertIntoFromStmt:
		return string(stmt.Sink)
	case parser.InsertIntoSelectStmt:
		return string(stmt.Sink)
	}
	// There's no need to cache the result of sinks.
	// CREATE STATE returns "" because it doesn't create a node.
	return ""
}

// String returns a string representation of the statement.
func (s *Statement) String() string {
	return fmt.Sprint(s.Stmt)
}

// IsStream returns true if the statement is a stream.
func (s *Statement) IsStream() bool {
	switch s.Stmt.(type) {
	case parser.CreateStreamAsSelectStmt, parser.CreateStreamAsSelectUnionStmt:
		return true
	}
	return false
}

// IsDataSourceNodeQuery returns true if the statement creates or manipulate a
// data source node, which is a source or a stream.
func (s *Statement) IsDataSourceNodeQuery() bool {
	switch s.Stmt.(type) {
	case parser.CreateSourceStmt, parser.CreateStreamAsSelectStmt,
		parser.CreateStreamAsSelectUnionStmt:
		return true
	}
	return false
}

// IsInsertStatement returns true when the statement is INSERT INTO.
func (s *Statement) IsInsertStatement() bool {
	switch s.Stmt.(type) {
	case parser.InsertIntoFromStmt, parser.InsertIntoSelectStmt:
		return true
	}
	return false
}

// IsCreateStateStatement returns true when the statement is CREATE STATE.
// Because LOAD STATE isn't supported by the command, it only check if the
// statement is CREATE STATE.
func (s *Statement) IsCreateStateStatement() bool {
	switch s.Stmt.(type) {
	case parser.CreateStateStmt:
		return true
	}
	return false
}

// Statements is a collection of Statement.
type Statements struct {
	Stmts []*Statement
}

// Parse parses BQL statements and create Statements.
func Parse(bql string) (*Statements, error) {
	p := parser.New()
	ss := &Statements{}
	for bql != "" {
		stmt, rest, err := p.ParseStmt(bql)
		if err != nil {
			// TODO: more detailed error reporting may be required
			return nil, err
		}

		s, err := NewStatement(stmt)
		if err != nil {
			return nil, err
		}
		ss.Stmts = append(ss.Stmts, s)
		bql = rest
	}
	return ss, nil
}

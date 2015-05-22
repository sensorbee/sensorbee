package parser

// This file holds a set of structs that make up the Abstract
// Syntax Tree of a BQL statement. Usually, for every rule in
// the PEG file, the left side should correspond to a struct
// in this file with the same name.

// Combined Structures (all with *AST)

type SelectStmt struct {
	ProjectionsAST
	FromAST
	FilterAST
	GroupingAST
	HavingAST
}

type CreateStreamStmt struct {
	Relation
	EmitProjectionsAST
	WindowedFromAST
	FilterAST
	GroupingAST
	HavingAST
}

type EmitProjectionsAST struct {
	EmitterType Emitter
	ProjectionsAST
}

type ProjectionsAST struct {
	Projections []interface{}
}

type WindowedFromAST struct {
	FromAST
	RangeAST
}

type RangeAST struct {
	Raw
	Unit RangeUnit
}

type FromAST struct {
	Relations []Relation
}

type FilterAST struct {
	Filter interface{}
}

type GroupingAST struct {
	GroupList []interface{}
}

type HavingAST struct {
	Having interface{}
}

type BinaryOpAST struct {
	Op    string
	Left  interface{}
	Right interface{}
}

// Elementary Structures (all without *AST for now)

// Note that we need the constructors for the elementary structures
// because we cannot use curly brackets for Expr{...} style
// initialization in the .peg file.

type Relation struct {
	Name string
}

func NewRelation(s string) Relation {
	return Relation{s}
}

type ColumnName struct {
	Name string
}

func NewColumnName(s string) ColumnName {
	return ColumnName{s}
}

type Raw struct {
	Expr string
}

func NewRaw(s string) Raw {
	return Raw{s}
}

type Emitter int

const (
	Istream Emitter = iota
	Dstream
	Rstream
)

type RangeUnit int

const (
	Tuples RangeUnit = iota
	Seconds
)

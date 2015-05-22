package parser

// This file holds a set of structs that make up the Abstract
// Syntax Tree of a BQL statement. Usually, for every rule in
// the PEG file, the left side should correspond to a struct
// in this file with the same name.

// Combined Structures

type SelectStmt struct {
	Projections
	From
	Filter
	Grouping
	Having
}

type CreateStreamStmt struct {
	Relation
	EmitProjections
	WindowedFrom
	Filter
	Grouping
	Having
}

type EmitProjections struct {
	emitterType Emitter
	Projections
}

type Projections struct {
	projections []interface{}
}

type WindowedFrom struct {
	From
	Range
}

type Range struct {
	Raw
	unit RangeUnit
}

type From struct {
	relations []Relation
}

type Filter struct {
	filter interface{}
}

type Grouping struct {
	groupList []interface{}
}

type Having struct {
	having interface{}
}

type BinaryOp struct {
	op    string
	left  interface{}
	right interface{}
}

// Elementary Structures

// Note that we need the constructors for the elementary structures
// because we cannot use curly brackets for Expr{...} style
// initialization in the .peg file.

type Relation struct {
	name string
}

func NewRelation(s string) Relation {
	return Relation{s}
}

type ColumnName struct {
	name string
}

func NewColumnName(s string) ColumnName {
	return ColumnName{s}
}

type Raw struct {
	expr string
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

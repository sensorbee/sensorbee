package parser

import (
	"strconv"
	"strings"
)

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

type CreateStreamAsSelectStmt struct {
	Relation
	EmitProjectionsAST
	WindowedFromAST
	FilterAST
	GroupingAST
	HavingAST
}

type CreateSourceStmt struct {
	Name SourceSinkName
	Type SourceSinkType
	SourceSinkSpecsAST
}

type CreateSinkStmt struct {
	Name SourceSinkName
	Type SourceSinkType
	SourceSinkSpecsAST
}

type CreateStreamFromSourceStmt struct {
	Relation
	Source SourceSinkName
}

type CreateStreamFromSourceExtStmt struct {
	Relation
	Type SourceSinkType
	SourceSinkSpecsAST
}

type InsertIntoSelectStmt struct {
	Sink SourceSinkName
	SelectStmt
}

type EmitProjectionsAST struct {
	EmitterType Emitter
	ProjectionsAST
}

type ProjectionsAST struct {
	Projections []interface{}
}

type AliasAST struct {
	Expr  interface{}
	Alias string
}

type WindowedFromAST struct {
	FromAST
	RangeAST
}

type RangeAST struct {
	NumericLiteral
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

type SourceSinkSpecsAST struct {
	Params []SourceSinkParamAST
}

type SourceSinkParamAST struct {
	Key   SourceSinkParamKey
	Value SourceSinkParamVal
}

type BinaryOpAST struct {
	Op    Operator
	Left  interface{}
	Right interface{}
}

type FuncAppAST struct {
	Function FuncName
	ExpressionsAST
}

type ExpressionsAST struct {
	Expressions []interface{}
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

type Wildcard struct {
}

func NewWildcard() Wildcard {
	return Wildcard{}
}

type RowValue struct {
	Relation string
	Column   string
}

func NewRowValue(s string) RowValue {
	return RowValue{"", s}
}

type Raw struct {
	Expr string
}

func NewRaw(s string) Raw {
	return Raw{s}
}

type NumericLiteral struct {
	Value int64
}

func NewNumericLiteral(s string) NumericLiteral {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return NumericLiteral{val}
}

type FloatLiteral struct {
	Value float64
}

func NewFloatLiteral(s string) FloatLiteral {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return FloatLiteral{val}
}

type BoolLiteral struct {
	Value bool
}

func NewBoolLiteral(b bool) BoolLiteral {
	return BoolLiteral{b}
}

type StringLiteral struct {
	Value string
}

func NewStringLiteral(s string) StringLiteral {
	runes := []rune(s)
	stripped := string(runes[1 : len(runes)-1])
	unescaped := strings.Replace(stripped, "''", "'", -1)
	return StringLiteral{unescaped}
}

type FuncName string

type SourceSinkName string

type SourceSinkType string

type SourceSinkParamKey string

type SourceSinkParamVal string

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

type Operator int

const (
	Or Operator = iota
	And
	Equal
	Less
	LessOrEqual
	Greater
	GreaterOrEqual
	NotEqual
	Plus
	Minus
	Multiply
	Divide
	Modulo
)

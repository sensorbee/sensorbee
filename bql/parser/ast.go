package parser

import (
	"strconv"
	"strings"
)

type Expression interface {
	ReferencedRelations() map[string]bool
	RenameReferencedRelation(string, string) Expression
}

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
	Projections []Expression
}

type AliasAST struct {
	Expr  Expression
	Alias string
}

func (a AliasAST) ReferencedRelations() map[string]bool {
	return a.Expr.ReferencedRelations()
}

func (a AliasAST) RenameReferencedRelation(from, to string) Expression {
	return AliasAST{a.Expr.RenameReferencedRelation(from, to), a.Alias}
}

type WindowedFromAST struct {
	Relations []AliasWindowedRelationAST
}

type AliasWindowedRelationAST struct {
	WindowedRelationAST
	Alias string
}

type WindowedRelationAST struct {
	Relation
	RangeAST
}

type RangeAST struct {
	NumericLiteral
	Unit RangeUnit
}

type FromAST struct {
	Relations []AliasRelationAST
}

func (f FromAST) ToWindowedFrom(size int64, unit RangeUnit) WindowedFromAST {
	output := make([]AliasWindowedRelationAST, len(f.Relations))
	for i, rel := range f.Relations {
		output[i] = AliasWindowedRelationAST{
			WindowedRelationAST{Relation{rel.Name},
				RangeAST{NumericLiteral{size}, unit}},
			rel.Alias,
		}
	}
	return WindowedFromAST{output}
}

type FilterAST struct {
	Filter Expression
}

type GroupingAST struct {
	GroupList []Expression
}

type HavingAST struct {
	Having Expression
}

type AliasRelationAST struct {
	Name  string
	Alias string
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
	Left  Expression
	Right Expression
}

func (b BinaryOpAST) ReferencedRelations() map[string]bool {
	rels := b.Left.ReferencedRelations()
	for rel := range b.Right.ReferencedRelations() {
		rels[rel] = true
	}
	return rels
}

func (b BinaryOpAST) RenameReferencedRelation(from, to string) Expression {
	return BinaryOpAST{b.Op,
		b.Left.RenameReferencedRelation(from, to),
		b.Right.RenameReferencedRelation(from, to)}
}

type FuncAppAST struct {
	Function FuncName
	ExpressionsAST
}

func (f FuncAppAST) ReferencedRelations() map[string]bool {
	rels := map[string]bool{}
	for _, expr := range f.Expressions {
		for rel := range expr.ReferencedRelations() {
			rels[rel] = true
		}
	}
	return rels
}

func (f FuncAppAST) RenameReferencedRelation(from, to string) Expression {
	newExprs := make([]Expression, len(f.Expressions))
	for i, expr := range f.Expressions {
		newExprs[i] = expr.RenameReferencedRelation(from, to)
	}
	return FuncAppAST{f.Function, ExpressionsAST{newExprs}}
}

type ExpressionsAST struct {
	Expressions []Expression
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

func (w Wildcard) ReferencedRelations() map[string]bool {
	return map[string]bool{"": true}
}

func (w Wildcard) RenameReferencedRelation(from, to string) Expression {
	return Wildcard{}
}

func NewWildcard() Wildcard {
	return Wildcard{}
}

type RowValue struct {
	Relation string
	Column   string
}

func (rv RowValue) ReferencedRelations() map[string]bool {
	return map[string]bool{rv.Relation: true}
}

func (rv RowValue) RenameReferencedRelation(from, to string) Expression {
	if rv.Relation == from {
		return RowValue{to, rv.Column}
	}
	return rv
}

func NewRowValue(s string) RowValue {
	components := strings.SplitN(s, ".", 2)
	if len(components) == 1 {
		// just "col"
		return RowValue{"", components[0]}
	}
	// "table.col"
	return RowValue{components[0], components[1]}

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

func (l NumericLiteral) ReferencedRelations() map[string]bool {
	return nil
}

func (l NumericLiteral) RenameReferencedRelation(from, to string) Expression {
	return l
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

func (l FloatLiteral) ReferencedRelations() map[string]bool {
	return nil
}

func (l FloatLiteral) RenameReferencedRelation(from, to string) Expression {
	return l
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

func (l BoolLiteral) ReferencedRelations() map[string]bool {
	return nil
}

func (l BoolLiteral) RenameReferencedRelation(from, to string) Expression {
	return l
}

func NewBoolLiteral(b bool) BoolLiteral {
	return BoolLiteral{b}
}

type StringLiteral struct {
	Value string
}

func (l StringLiteral) ReferencedRelations() map[string]bool {
	return nil
}

func (l StringLiteral) RenameReferencedRelation(from, to string) Expression {
	return l
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

type Identifier string

package parser

import (
	"pfi/sensorbee/sensorbee/data"
	"strconv"
	"strings"
)

type Expression interface {
	ReferencedRelations() map[string]bool
	RenameReferencedRelation(string, string) Expression
	Foldable() bool
}

// This file holds a set of structs that make up the Abstract
// Syntax Tree of a BQL statement. Usually, for every rule in
// the PEG file, the left side should correspond to a struct
// in this file with the same name.

// Combined Structures (all with *AST)

type SelectStmt struct {
	EmitterAST
	ProjectionsAST
	WindowedFromAST
	FilterAST
	GroupingAST
	HavingAST
}

type CreateStreamAsSelectStmt struct {
	Name StreamIdentifier
	EmitterAST
	ProjectionsAST
	WindowedFromAST
	FilterAST
	GroupingAST
	HavingAST
}

type CreateSourceStmt struct {
	Paused BinaryKeyword
	Name   StreamIdentifier
	Type   SourceSinkType
	SourceSinkSpecsAST
}

type CreateSinkStmt struct {
	Name StreamIdentifier
	Type SourceSinkType
	SourceSinkSpecsAST
}

type CreateStateStmt struct {
	Name StreamIdentifier
	Type SourceSinkType
	SourceSinkSpecsAST
}

type InsertIntoSelectStmt struct {
	Sink StreamIdentifier
	SelectStmt
}

type InsertIntoFromStmt struct {
	Sink  StreamIdentifier
	Input StreamIdentifier
}

type PauseSourceStmt struct {
	Source StreamIdentifier
}

type ResumeSourceStmt struct {
	Source StreamIdentifier
}

type RewindSourceStmt struct {
	Source StreamIdentifier
}

type DropSourceStmt struct {
	Source StreamIdentifier
}

type DropStreamStmt struct {
	Source StreamIdentifier
}

type DropSinkStmt struct {
	Source StreamIdentifier
}

type EmitterAST struct {
	EmitterType   Emitter
	EmitIntervals []StreamEmitIntervalAST
}

type StreamEmitIntervalAST struct {
	IntervalAST
	Stream
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

func (a AliasAST) Foldable() bool {
	return a.Expr.Foldable()
}

type WindowedFromAST struct {
	Relations []AliasedStreamWindowAST
}

type AliasedStreamWindowAST struct {
	StreamWindowAST
	Alias string
}

type StreamWindowAST struct {
	Stream
	IntervalAST
}

type IntervalAST struct {
	NumericLiteral
	Unit IntervalUnit
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

type SourceSinkSpecsAST struct {
	Params []SourceSinkParamAST
}

type SourceSinkParamAST struct {
	Key   SourceSinkParamKey
	Value data.Value
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

func (b BinaryOpAST) Foldable() bool {
	return b.Left.Foldable() && b.Right.Foldable()
}

type UnaryOpAST struct {
	Op   Operator
	Expr Expression
}

func (u UnaryOpAST) ReferencedRelations() map[string]bool {
	return u.Expr.ReferencedRelations()
}

func (u UnaryOpAST) RenameReferencedRelation(from, to string) Expression {
	return UnaryOpAST{u.Op,
		u.Expr.RenameReferencedRelation(from, to)}
}

func (u UnaryOpAST) Foldable() bool {
	return u.Expr.Foldable()
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

func (f FuncAppAST) Foldable() bool {
	foldable := true
	for _, expr := range f.Expressions {
		if !expr.Foldable() {
			foldable = false
			break
		}
	}
	return foldable
}

type ExpressionsAST struct {
	Expressions []Expression
}

// Elementary Structures (all without *AST for now)

// Note that we need the constructors for the elementary structures
// because we cannot use curly brackets for Expr{...} style
// initialization in the .peg file.

// It seems not possible in Go to have a variable that says "this is
// either struct A or struct B or struct C", so we build one struct
// that serves both for "real" streams (as in `FROM x`) and stream-
// generating functions (as in `FROM series(1, 5)`).
type Stream struct {
	Type   StreamType
	Name   string
	Params []Expression
}

func NewStream(s string) Stream {
	return Stream{ActualStream, s, nil}
}

type Wildcard struct {
}

func (w Wildcard) ReferencedRelations() map[string]bool {
	return map[string]bool{"": true}
}

func (w Wildcard) RenameReferencedRelation(from, to string) Expression {
	return Wildcard{}
}

func (w Wildcard) Foldable() bool {
	return false
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

func (rv RowValue) Foldable() bool {
	return false
}

func NewRowValue(s string) RowValue {
	// TODO when we support full JSONPath this must become more
	//      sophisticated in order to deal, for example, with:
	//        `SELECT elem["foo:bar"] FROM mystream`
	components := strings.SplitN(s, ":", 2)
	if len(components) == 1 {
		// just "col"
		return RowValue{"", components[0]}
	}
	// "table.col"
	return RowValue{components[0], components[1]}
}

type RowMeta struct {
	Relation string
	MetaType MetaInformation
}

func (rm RowMeta) ReferencedRelations() map[string]bool {
	return map[string]bool{rm.Relation: true}
}

func (rm RowMeta) RenameReferencedRelation(from, to string) Expression {
	if rm.Relation == from {
		return RowMeta{to, rm.MetaType}
	}
	return rm
}

func (rm RowMeta) Foldable() bool {
	return false
}

func NewRowMeta(s string, t MetaInformation) RowMeta {
	components := strings.SplitN(s, ":", 2)
	if len(components) == 1 {
		// just the meta information
		return RowMeta{"", t}
	}
	// relation name and meta information
	return RowMeta{components[0], t}
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

func (l NumericLiteral) Foldable() bool {
	return true
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

func (l FloatLiteral) Foldable() bool {
	return true
}

func NewFloatLiteral(s string) FloatLiteral {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return FloatLiteral{val}
}

type NullLiteral struct {
}

func (l NullLiteral) ReferencedRelations() map[string]bool {
	return nil
}

func (l NullLiteral) RenameReferencedRelation(from, to string) Expression {
	return l
}

func (l NullLiteral) Foldable() bool {
	return true
}

func NewNullLiteral() NullLiteral {
	return NullLiteral{}
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

func (l BoolLiteral) Foldable() bool {
	return true
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

func (l StringLiteral) Foldable() bool {
	return true
}

func NewStringLiteral(s string) StringLiteral {
	runes := []rune(s)
	stripped := string(runes[1 : len(runes)-1])
	unescaped := strings.Replace(stripped, "''", "'", -1)
	return StringLiteral{unescaped}
}

type FuncName string

type StreamIdentifier string

type SourceSinkType string

type SourceSinkParamKey string

type Emitter int

const (
	UnspecifiedEmitter Emitter = iota
	Istream
	Dstream
	Rstream
)

func (e Emitter) String() string {
	s := "UNSPECIFIED"
	switch e {
	case Istream:
		s = "ISTREAM"
	case Dstream:
		s = "DSTREAM"
	case Rstream:
		s = "RSTREAM"
	}
	return s
}

type StreamType int

const (
	UnknownStreamType StreamType = iota
	ActualStream
	UDSFStream
)

func (st StreamType) String() string {
	s := "UNKNOWN"
	switch st {
	case ActualStream:
		s = "ActualStream"
	case UDSFStream:
		s = "UDSFStream"
	}
	return s
}

type IntervalUnit int

const (
	UnspecifiedIntervalUnit IntervalUnit = iota
	Tuples
	Seconds
)

func (i IntervalUnit) String() string {
	s := "UNSPECIFIED"
	switch i {
	case Tuples:
		s = "TUPLES"
	case Seconds:
		s = "SECONDS"
	}
	return s
}

type MetaInformation int

const (
	UnknownMeta MetaInformation = iota
	TimestampMeta
)

func (m MetaInformation) String() string {
	s := "UnknownMeta"
	switch m {
	case TimestampMeta:
		s = "TS"
	}
	return s
}

type BinaryKeyword int

const (
	UnspecifiedKeyword BinaryKeyword = iota
	Yes
	No
)

func (k BinaryKeyword) String() string {
	s := "Unspecified"
	switch k {
	case Yes:
		s = "Yes"
	case No:
		s = "No"
	}
	return s
}

type Operator int

const (
	UnknownOperator Operator = iota
	Or
	And
	Not
	Equal
	Less
	LessOrEqual
	Greater
	GreaterOrEqual
	NotEqual
	Is
	IsNot
	Plus
	Minus
	Multiply
	Divide
	Modulo
	UnaryMinus
)

func (o Operator) String() string {
	s := "UnknownOperator"
	switch o {
	case Or:
		s = "OR"
	case And:
		s = "AND"
	case Not:
		s = "NOT"
	case Equal:
		s = "="
	case Less:
		s = "<"
	case LessOrEqual:
		s = "<="
	case Greater:
		s = ">"
	case GreaterOrEqual:
		s = ">="
	case NotEqual:
		s = "!="
	case Is:
		s = "IS"
	case IsNot:
		s = "IS NOT"
	case Plus:
		s = "+"
	case Minus:
		s = "-"
	case Multiply:
		s = "*"
	case Divide:
		s = "/"
	case Modulo:
		s = "%"
	case UnaryMinus:
		s = "-"
	}
	return s
}

type Identifier string

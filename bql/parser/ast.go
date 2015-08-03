package parser

import (
	"fmt"
	"pfi/sensorbee/sensorbee/data"
	"strconv"
	"strings"
)

type Expression interface {
	ReferencedRelations() map[string]bool
	RenameReferencedRelation(string, string) Expression
	Foldable() bool
	string() string
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

func (s SelectStmt) String() string {
	str := []string{"SELECT", s.EmitterAST.string()}
	str = append(str, s.ProjectionsAST.string())
	str = append(str, s.WindowedFromAST.string())
	str = append(str, s.FilterAST.string())
	str = append(str, s.GroupingAST.string())
	str = append(str, s.HavingAST.string())

	st := []string{}
	for _, s := range str {
		if s != "" {
			st = append(st, s)
		}
	}
	return strings.Join(st, " ")
}

type SelectUnionStmt struct {
	Selects []SelectStmt
}

type CreateStreamAsSelectStmt struct {
	Name   StreamIdentifier
	Select SelectStmt
}

func (s CreateStreamAsSelectStmt) String() string {
	str := []string{"CREATE", "STREAM", string(s.Name), "AS", s.Select.String()}
	return strings.Join(str, " ")
}

type CreateStreamAsSelectUnionStmt struct {
	Name StreamIdentifier
	SelectUnionStmt
}

type CreateSourceStmt struct {
	Paused BinaryKeyword
	Name   StreamIdentifier
	Type   SourceSinkType
	SourceSinkSpecsAST
}

func (s CreateSourceStmt) String() string {
	str := []string{"CREATE", "SOURCE", string(s.Name), "TYPE", string(s.Type)}
	paused := s.Paused.string("PAUSED", "UNPAUSED")
	if paused != "" {
		str = append(str[:1], append([]string{paused}, str[1:]...)...)
	}
	specs := s.SourceSinkSpecsAST.string("WITH")
	if specs != "" {
		str = append(str, specs)
	}
	return strings.Join(str, " ")
}

type CreateSinkStmt struct {
	Name StreamIdentifier
	Type SourceSinkType
	SourceSinkSpecsAST
}

func (s CreateSinkStmt) String() string {
	str := []string{"CREATE", "SINK", string(s.Name), "TYPE", string(s.Type)}
	specs := s.SourceSinkSpecsAST.string("WITH")
	if specs != "" {
		str = append(str, specs)
	}
	return strings.Join(str, " ")
}

type CreateStateStmt struct {
	Name StreamIdentifier
	Type SourceSinkType
	SourceSinkSpecsAST
}

func (s CreateStateStmt) String() string {
	str := []string{"CREATE", "STATE", string(s.Name), "TYPE", string(s.Type)}
	specs := s.SourceSinkSpecsAST.string("WITH")
	if specs != "" {
		str = append(str, specs)
	}
	return strings.Join(str, " ")
}

type UpdateStateStmt struct {
	Name StreamIdentifier
	SourceSinkSpecsAST
}

func (s UpdateStateStmt) String() string {
	str := []string{"UPDATE", "STATE", string(s.Name)}
	specs := s.SourceSinkSpecsAST.string("SET")
	if specs != "" {
		str = append(str, specs)
	}
	return strings.Join(str, " ")
}

type UpdateSourceStmt struct {
	Name StreamIdentifier
	SourceSinkSpecsAST
}

func (s UpdateSourceStmt) String() string {
	str := []string{"UPDATE", "SOURCE", string(s.Name)}
	specs := s.SourceSinkSpecsAST.string("SET")
	if specs != "" {
		str = append(str, specs)
	}
	return strings.Join(str, " ")
}

type UpdateSinkStmt struct {
	Name StreamIdentifier
	SourceSinkSpecsAST
}

func (s UpdateSinkStmt) String() string {
	str := []string{"UPDATE", "SINK", string(s.Name)}
	specs := s.SourceSinkSpecsAST.string("SET")
	if specs != "" {
		str = append(str, specs)
	}
	return strings.Join(str, " ")
}

type InsertIntoSelectStmt struct {
	Sink StreamIdentifier
	SelectStmt
}

func (s InsertIntoSelectStmt) String() string {
	str := []string{"INSERT", "INTO", string(s.Sink), s.SelectStmt.String()}
	return strings.Join(str, " ")
}

type InsertIntoFromStmt struct {
	Sink  StreamIdentifier
	Input StreamIdentifier
}

func (s InsertIntoFromStmt) String() string {
	str := []string{"INSERT", "INTO", string(s.Sink), "FROM", string(s.Input)}
	return strings.Join(str, " ")
}

type PauseSourceStmt struct {
	Source StreamIdentifier
}

func (s PauseSourceStmt) String() string {
	str := []string{"PAUSE", "SOURCE", string(s.Source)}
	return strings.Join(str, " ")
}

type ResumeSourceStmt struct {
	Source StreamIdentifier
}

func (s ResumeSourceStmt) String() string {
	str := []string{"RESUME", "SOURCE", string(s.Source)}
	return strings.Join(str, " ")
}

type RewindSourceStmt struct {
	Source StreamIdentifier
}

func (s RewindSourceStmt) String() string {
	str := []string{"REWIND", "SOURCE", string(s.Source)}
	return strings.Join(str, " ")
}

type DropSourceStmt struct {
	Source StreamIdentifier
}

func (s DropSourceStmt) String() string {
	str := []string{"DROP", "SOURCE", string(s.Source)}
	return strings.Join(str, " ")
}

type DropStreamStmt struct {
	Stream StreamIdentifier
}

func (s DropStreamStmt) String() string {
	str := []string{"DROP", "STREAM", string(s.Stream)}
	return strings.Join(str, " ")
}

type DropSinkStmt struct {
	Sink StreamIdentifier
}

func (s DropSinkStmt) String() string {
	str := []string{"DROP", "SINK", string(s.Sink)}
	return strings.Join(str, " ")
}

type DropStateStmt struct {
	State StreamIdentifier
}

func (s DropStateStmt) String() string {
	str := []string{"DROP", "STATE", string(s.State)}
	return strings.Join(str, " ")
}

type EmitterAST struct {
	EmitterType Emitter
	// here is space for some emit options later on
}

func (a EmitterAST) string() string {
	return a.EmitterType.String()
}

type ProjectionsAST struct {
	Projections []Expression
}

func (a ProjectionsAST) string() string {
	prj := []string{}
	for _, e := range a.Projections {
		prj = append(prj, e.string())
	}
	return strings.Join(prj, ", ")
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

func (a AliasAST) string() string {
	return a.Expr.string() + " AS " + a.Alias
}

type WindowedFromAST struct {
	Relations []AliasedStreamWindowAST
}

func (a WindowedFromAST) string() string {
	if len(a.Relations) == 0 {
		return ""
	}

	str := []string{}
	for _, r := range a.Relations {
		str = append(str, r.string())
	}
	return "FROM " + strings.Join(str, ", ")
}

type AliasedStreamWindowAST struct {
	StreamWindowAST
	Alias string
}

func (a AliasedStreamWindowAST) string() string {
	str := a.StreamWindowAST.string()
	if a.Alias != "" {
		str = str + " AS " + a.Alias
	}
	return str
}

type StreamWindowAST struct {
	Stream
	IntervalAST
}

func (a StreamWindowAST) string() string {
	interval := a.IntervalAST.string()

	switch a.Stream.Type {
	case ActualStream:
		return a.Stream.Name + " " + interval

	case UDSFStream:
		ps := []string{}
		for _, p := range a.Stream.Params {
			ps = append(ps, p.string())
		}
		return a.Stream.Name + "(" + strings.Join(ps, ", ") + ") " + interval
	}

	return "UnknownStreamType"
}

type IntervalAST struct {
	NumericLiteral
	Unit IntervalUnit
}

func (a IntervalAST) string() string {
	return "[RANGE " + a.NumericLiteral.string() + " " + a.Unit.String() + "]"
}

type FilterAST struct {
	Filter Expression
}

func (a FilterAST) string() string {
	if a.Filter == nil {
		return ""
	}
	return "WHERE " + a.Filter.string()
}

type GroupingAST struct {
	GroupList []Expression
}

func (a GroupingAST) string() string {
	if len(a.GroupList) == 0 {
		return ""
	}

	str := []string{}
	for _, e := range a.GroupList {
		str = append(str, e.string())
	}
	return "GROUP BY " + strings.Join(str, ", ")
}

type HavingAST struct {
	Having Expression
}

func (a HavingAST) string() string {
	if a.Having == nil {
		return ""
	}
	return "HAVING " + a.Having.string()
}

type SourceSinkSpecsAST struct {
	Params []SourceSinkParamAST
}

func (a SourceSinkSpecsAST) string(keyword string) string {
	if len(a.Params) == 0 {
		return ""
	}
	ps := make([]string, len(a.Params))
	for i, p := range a.Params {
		ps[i] = p.string()
	}
	return keyword + " " + strings.Join(ps, ", ")
}

type SourceSinkParamAST struct {
	Key   SourceSinkParamKey
	Value data.Value
}

func (a SourceSinkParamAST) string() string {
	s, _ := data.ToString(a.Value)
	if a.Value.Type() == data.TypeString {
		s = "'" + strings.Replace(s, "'", "''", -1) + "'"
	}
	return string(a.Key) + "=" + s
}

type BinaryOpAST struct {
	Op    Operator
	Left  Expression
	Right Expression
}

func (b BinaryOpAST) ReferencedRelations() map[string]bool {
	rels := b.Left.ReferencedRelations()
	if rels == nil {
		return b.Right.ReferencedRelations()
	}
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

func (b BinaryOpAST) string() string {
	str := []string{b.Left.string(), b.Op.String(), b.Right.string()}

	// TODO: This implementation may add unnecessary parentheses.
	// For example, in
	//  input:  "a * 2 / b"
	//  output: "(a * 2) / b"
	// we could omit output parentehsis.

	// Enclose expression in parentheses for operator precedence
	encloseLeft, encloseRight := false, false

	if left, ok := b.Left.(BinaryOpAST); ok {
		if left.Op.hasHigherPrecedenceThan(b.Op) {
			// we need no parentheses
		} else {
			// we probably need parentheses
			encloseLeft = true
		}
	}

	if right, ok := b.Right.(BinaryOpAST); ok {
		if right.Op.hasHigherPrecedenceThan(b.Op) {
			// we need no parentheses
		} else {
			// we probably need parentheses
			encloseRight = true
		}
	}

	if encloseLeft {
		str[0] = "(" + str[0] + ")"
	}
	if encloseRight {
		str[2] = "(" + str[2] + ")"
	}

	return strings.Join(str, " ")
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

func (u UnaryOpAST) string() string {
	op := u.Op.String()
	expr := u.Expr.string()

	// Unary minus operator such as "- - 2"
	if u.Op != UnaryMinus || strings.HasPrefix(expr, "-") {
		op = op + " "
	}

	// Enclose expression in parentheses for "NOT (a AND B)" like case
	if _, ok := u.Expr.(BinaryOpAST); ok {
		expr = "(" + expr + ")"
	}

	return op + expr
}

type TypeCastAST struct {
	Expr   Expression
	Target Type
}

func (u TypeCastAST) ReferencedRelations() map[string]bool {
	return u.Expr.ReferencedRelations()
}

func (u TypeCastAST) RenameReferencedRelation(from, to string) Expression {
	return TypeCastAST{u.Expr.RenameReferencedRelation(from, to),
		u.Target}
}

func (u TypeCastAST) Foldable() bool {
	return u.Expr.Foldable()
}

func (u TypeCastAST) string() string {
	if rv, ok := u.Expr.(RowValue); ok {
		return rv.string() + "::" + u.Target.String()
	}

	if rm, ok := u.Expr.(RowMeta); ok {
		return rm.string() + "::" + u.Target.String()
	}

	return "CAST(" + u.Expr.string() + " AS " + u.Target.String() + ")"
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
	// now() is not evaluable outside of some execution context
	if string(f.Function) == "now" && len(f.Expressions) == 0 {
		return false
	}
	for _, expr := range f.Expressions {
		if !expr.Foldable() {
			foldable = false
			break
		}
	}
	return foldable
}

type ArrayAST struct {
	ExpressionsAST
}

func (a ArrayAST) ReferencedRelations() map[string]bool {
	rels := map[string]bool{}
	for _, expr := range a.Expressions {
		for rel := range expr.ReferencedRelations() {
			rels[rel] = true
		}
	}
	return rels
}

func (a ArrayAST) RenameReferencedRelation(from, to string) Expression {
	newExprs := make([]Expression, len(a.Expressions))
	for i, expr := range a.Expressions {
		newExprs[i] = expr.RenameReferencedRelation(from, to)
	}
	return ArrayAST{ExpressionsAST{newExprs}}
}

func (a ArrayAST) Foldable() bool {
	foldable := true
	for _, expr := range a.Expressions {
		if !expr.Foldable() {
			foldable = false
			break
		}
	}
	return foldable
}

func (a ArrayAST) string() string {
	return "[" + a.ExpressionsAST.string() + "]"
}

func (f FuncAppAST) string() string {
	return string(f.Function) + "(" + f.ExpressionsAST.string() + ")"
}

type ExpressionsAST struct {
	Expressions []Expression
}

func (a ExpressionsAST) string() string {
	str := []string{}
	for _, e := range a.Expressions {
		str = append(str, e.string())
	}
	return strings.Join(str, ", ")
}

type MapAST struct {
	Entries []KeyValuePairAST
}

func (m MapAST) ReferencedRelations() map[string]bool {
	rels := map[string]bool{}
	for _, pair := range m.Entries {
		for rel := range pair.Value.ReferencedRelations() {
			rels[rel] = true
		}
	}
	return rels
}

func (m MapAST) RenameReferencedRelation(from, to string) Expression {
	newEntries := make([]KeyValuePairAST, len(m.Entries))
	for i, pair := range m.Entries {
		newEntries[i] = KeyValuePairAST{
			pair.Key,
			pair.Value.RenameReferencedRelation(from, to),
		}
	}
	return MapAST{newEntries}
}

func (m MapAST) Foldable() bool {
	foldable := true
	for _, pair := range m.Entries {
		if !pair.Value.Foldable() {
			foldable = false
			break
		}
	}
	return foldable
}

func (m MapAST) string() string {
	entries := []string{}
	for _, pair := range m.Entries {
		entries = append(entries, pair.string())
	}
	return "{" + strings.Join(entries, ", ") + "}"
}

type KeyValuePairAST struct {
	Key   string
	Value Expression
}

func (k KeyValuePairAST) string() string {
	return `'` + k.Key + `':` + k.Value.string()
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
	Relation string
}

func (w Wildcard) ReferencedRelations() map[string]bool {
	if w.Relation == "" {
		// the wildcard does not reference any relation
		// (this is different to referencing the "" relation)
		return nil
	}
	return map[string]bool{w.Relation: true}
}

func (w Wildcard) RenameReferencedRelation(from, to string) Expression {
	if w.Relation == from {
		return Wildcard{to}
	}
	return Wildcard{w.Relation}
}

func (w Wildcard) Foldable() bool {
	return false
}

func NewWildcard(relation string) Wildcard {
	return Wildcard{strings.TrimRight(relation, ":*")}
}

func (w Wildcard) string() string {
	if w.Relation != "" {
		return w.Relation + ":*"
	}
	return "*"
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

func (rv RowValue) string() string {
	if rv.Relation != "" {
		return rv.Relation + ":" + rv.Column
	}
	return rv.Column
}

func NewRowValue(s string) RowValue {
	bracketPos := strings.Index(s, "[")
	components := strings.SplitN(s, ":", 2)
	if bracketPos >= 0 && bracketPos < len(components[0]) {
		// if there is a bracket, then it is definitely on the right
		// side of the colon. therefore, if the part before the first
		// found colon is longer than where the first bracket is,
		// then the colon is part of the JSON path, not the stream
		// separator.
		return RowValue{"", s}
	} else if len(components) == 1 {
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

func (rm RowMeta) string() string {
	if rm.Relation != "" {
		return rm.Relation + ":" + rm.MetaType.string()
	}
	return rm.MetaType.string()
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

func (l NumericLiteral) string() string {
	return fmt.Sprintf("%v", l.Value)
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

func (l FloatLiteral) string() string {
	return fmt.Sprintf("%v", l.Value)
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

func (l NullLiteral) string() string {
	return "NULL"
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

func (l BoolLiteral) string() string {
	if l.Value {
		return "TRUE"
	}
	return "FALSE"
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

func (l StringLiteral) string() string {
	return "'" + strings.Replace(l.Value, "'", "''", -1) + "'"
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
	Milliseconds
)

func (i IntervalUnit) String() string {
	s := "UNSPECIFIED"
	switch i {
	case Tuples:
		s = "TUPLES"
	case Seconds:
		s = "SECONDS"
	case Milliseconds:
		s = "MILLISECONDS"
	}
	return s
}

type MetaInformation int

const (
	UnknownMeta MetaInformation = iota
	TimestampMeta
	NowMeta
)

func (m MetaInformation) String() string {
	s := "UnknownMeta"
	switch m {
	case TimestampMeta:
		s = "TS"
	case NowMeta:
		s = "NOW"
	}
	return s
}

func (m MetaInformation) string() string {
	s := "UnknownMeta"
	switch m {
	case TimestampMeta:
		s = "ts()"
	case NowMeta:
		s = "now()"
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

func (k BinaryKeyword) string(yes, no string) string {
	switch k {
	case Yes:
		return yes
	case No:
		return no
	}
	return ""
}

type Type int

const (
	UnknownType Type = iota
	Bool
	Int
	Float
	String
	Blob
	Timestamp
	Array
	Map
)

func (t Type) String() string {
	s := "UnknownType"
	switch t {
	case Bool:
		s = "BOOL"
	case Int:
		s = "INT"
	case Float:
		s = "FLOAT"
	case String:
		s = "STRING"
	case Blob:
		s = "BLOB"
	case Timestamp:
		s = "TIMESTAMP"
	case Array:
		s = "ARRAY"
	case Map:
		s = "MAP"
	}
	return s
}

type Operator int

const (
	// Operators are defined in precedence order (increasing). These
	// values can be compared using the hasHigherPrecedenceThan method.
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
	Concat
	Is
	IsNot
	Plus
	Minus
	Multiply
	Divide
	Modulo
	UnaryMinus
)

// hasSamePrecedenceAs checks if the arguement operator has the same precedence.
func (op Operator) hasSamePrecedenceAs(rhs Operator) bool {
	if Or <= op && op <= Not && Or <= rhs && rhs <= Not {
		return true
	}
	if Less <= op && op <= GreaterOrEqual && Less <= rhs && rhs <= GreaterOrEqual {
		return true
	}
	if Is <= op && op <= IsNot && Is <= rhs && rhs <= IsNot {
		return true
	}
	if Plus <= op && op <= Minus && Plus <= rhs && rhs <= Minus {
		return true
	}
	if Multiply <= op && op <= Modulo && Multiply <= rhs && rhs <= Modulo {
		return true
	}

	return false
}

func (op Operator) hasHigherPrecedenceThan(rhs Operator) bool {
	if op.hasSamePrecedenceAs(rhs) {
		return false
	}

	return op > rhs
}

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
	case Concat:
		s = "||"
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

package parser

import (
	"fmt"

	"gopkg.in/sensorbee/sensorbee.v0/data"
)

// parseStack is a standard stack implementation, but also holds
// methods for transforming the top k elements into a new element.
type parseStack struct {
	top  *stackElement
	size int
}

// stackElement is a stack-internal data structure that is used
// as a wrapper for the actual data.
type stackElement struct {
	value *ParsedComponent
	next  *stackElement
}

// ParsedComponent is an element of the parse stack that represents
// a section of the input string that was successfully parsed.
type ParsedComponent struct {
	// begin is the index of the first character that belongs to
	// the parsed statement
	begin int
	// end is the index of the last character that belongs to the
	// parsed statement + 1
	end int
	// comp stores the struct that the string was parsed into
	comp interface{}
}

// Len return the stack's size.
func (ps *parseStack) Len() int {
	return ps.size
}

// Push pushes a new element onto the stack.
func (ps *parseStack) Push(value *ParsedComponent) {
	ps.top = &stackElement{value, ps.top}
	ps.size++
}

// Pop removes the top element from the stack and returns its value.
// If the stack is empty, returns nil.
func (ps *parseStack) Pop() (value *ParsedComponent) {
	if ps.size > 0 {
		value, ps.top = ps.top.value, ps.top.next
		ps.size--
		return
	}
	return nil
}

// Peek returns the top element from the stack but doesn't remove it.
// If the stack is empty, returns nil.
func (ps *parseStack) Peek() (value *ParsedComponent) {
	if ps.size > 0 {
		return ps.top.value
	}
	return nil
}

// AssembleSelect takes the topmost elements from the stack, assuming
// they are components of a SELECT statement, and replaces them by
// a single SelectStmt element.
//
//  EmitterAST
//  HavingAST
//  GroupingAST
//  FilterAST
//  WindowedFromAST
//  ProjectionsAST
//   =>
//  SelectStmt{EmitterAST, ProjectionsAST, WindowedFromAST, FilterAST, GroupingAST, HavingAST}
func (ps *parseStack) AssembleSelect() {
	// pop the components from the stack in reverse order
	_having, _grouping, _filter, _from, _projections, _emitter := ps.pop6()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	having := _having.comp.(HavingAST)
	grouping := _grouping.comp.(GroupingAST)
	filter := _filter.comp.(FilterAST)
	from := _from.comp.(WindowedFromAST)
	projections := _projections.comp.(ProjectionsAST)
	emitter := _emitter.comp.(EmitterAST)

	// assemble the SelectStmt and push it back
	s := SelectStmt{emitter, projections, from, filter, grouping, having}
	se := ParsedComponent{_emitter.begin, _having.end, s}
	ps.Push(&se)
}

// AssembleSelectUnion takes the elements from the stack that
// correspond to the input[begin:end] string and wraps a
// SelectUnionStmt struct around them.
//
//  SelectStmt
//  SelectStmt
//  SelectStmt
//   =>
//  SelectUnionStmt{[SelectStmt, SelectStmt, SelectStmt]}
func (ps *parseStack) AssembleSelectUnion(begin int, end int) {
	elems := ps.collectElements(begin, end)
	selects := make([]SelectStmt, len(elems))
	for i := range elems {
		selects[i] = elems[i].(SelectStmt)
	}
	// push the grouped list back
	ps.PushComponent(begin, end, SelectUnionStmt{selects})
}

// AssembleCreateStreamAsSelect takes the topmost elements from the stack,
// assuming they are components of a CREATE STREAM statement, and
// replaces them by a single CreateStreamAsSelectStmt element.
//
//  SelectStmt
//  StreamIdentifier
//   =>
//  CreateStreamAsSelectStmt{StreamIdentifier, SelectStmt}
func (ps *parseStack) AssembleCreateStreamAsSelect() {
	// now pop the components from the stack in reverse order
	_select, _name := ps.pop2()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	s := _select.comp.(SelectStmt)
	name := _name.comp.(StreamIdentifier)

	// assemble the SelectStmt and push it back
	css := CreateStreamAsSelectStmt{name, s}
	se := ParsedComponent{_name.begin, _select.end, css}
	ps.Push(&se)
}

// AssembleCreateStreamAsSelectUnion takes the topmost elements from the
// stack, assuming they are components of a CREATE STREAM statement, and
// replaces them by a single CreateStreamAsSelectUnionStmt element.
//
//  SelectUnionStmt
//  StreamIdentifier
//   =>
//  CreateStreamAsSelectUnionStmt{StreamIdentifier, SelectUnionStmt}
func (ps *parseStack) AssembleCreateStreamAsSelectUnion() {
	// now pop the components from the stack in reverse order
	_selectUnion, _name := ps.pop2()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	selectUnion := _selectUnion.comp.(SelectUnionStmt)
	name := _name.comp.(StreamIdentifier)

	// assemble the SelectUnionStmt and push it back
	css := CreateStreamAsSelectUnionStmt{name, selectUnion}
	se := ParsedComponent{_name.begin, _selectUnion.end, css}
	ps.Push(&se)
}

// AssembleCreateSource takes the topmost elements from the stack,
// assuming they are components of a CREATE SOURCE statement, and
// replaces them by a single CreateSourceStmt element.
//
//  BinaryKeyword
//  SourceSinkSpecsAST
//  SourceSinkType
//  StreamIdentifier
//   =>
//  CreateSourceStmt{BinaryKeyword, StreamIdentifier, SourceSinkType,
//    SourceSinkSpecsAST}
func (ps *parseStack) AssembleCreateSource() {
	// pop the components from the stack in reverse order
	_specs, _sourceType, _name, _paused := ps.pop4()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	specs := _specs.comp.(SourceSinkSpecsAST)
	sourceType := _sourceType.comp.(SourceSinkType)
	name := _name.comp.(StreamIdentifier)
	paused := _paused.comp.(BinaryKeyword)

	// assemble the CreateSourceStmt and push it back
	s := CreateSourceStmt{paused, name, sourceType, specs}
	se := ParsedComponent{_paused.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleCreateSink takes the topmost elements from the stack,
// assuming they are components of a CREATE SINK statement, and
// replaces them by a single CreateSinkStmt element.
//
//  SourceSinkSpecsAST
//  SourceSinkType
//  StreamIdentifier
//   =>
//  CreateSinkStmt{StreamIdentifier, SourceSinkType, SourceSinkSpecsAST}
func (ps *parseStack) AssembleCreateSink() {
	_specs, _sinkType, _name := ps.pop3()

	specs := _specs.comp.(SourceSinkSpecsAST)
	sinkType := _sinkType.comp.(SourceSinkType)
	name := _name.comp.(StreamIdentifier)

	s := CreateSinkStmt{name, sinkType, specs}
	se := ParsedComponent{_name.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleCreateState takes the topmost elements from the stack,
// assuming they are components of a CREATE STATE statement, and
// replaces them by a single CreateStateStmt element.
//
//  SourceSinkSpecsAST
//  SourceSinkType
//  StreamIdentifier
//   =>
//  CreateStateStmt{StreamIdentifier, SourceSinkType, SourceSinkSpecsAST}
func (ps *parseStack) AssembleCreateState() {
	_specs, _sinkType, _name := ps.pop3()

	specs := _specs.comp.(SourceSinkSpecsAST)
	sinkType := _sinkType.comp.(SourceSinkType)
	name := _name.comp.(StreamIdentifier)

	s := CreateStateStmt{name, sinkType, specs}
	se := ParsedComponent{_name.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleUpdateState takes the topmost elements from the stack,
// assuming they are components of a UPDATE STATE statement, and
// replaces them by a single UpdateStateStmt element.
//
//  SourceSinkSpecsAST
//  StreamIdentifier
//   =>
//  UpdateStateStmt{StreamIdentifier, SourceSinkSpecsAST}
func (ps *parseStack) AssembleUpdateState() {
	_specs, _name := ps.pop2()

	specs := _specs.comp.(SourceSinkSpecsAST)
	name := _name.comp.(StreamIdentifier)

	s := UpdateStateStmt{name, specs}
	se := ParsedComponent{_name.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleUpdateSource takes the topmost elements from the stack,
// assuming they are components of a UPDATE SOURCE statement, and
// replaces them by a single UpdateSourceStmt element.
//
//  SourceSinkSpecsAST
//  StreamIdentifier
//   =>
//  UpdateSourceStmt{StreamIdentifier, SourceSinkSpecsAST}
func (ps *parseStack) AssembleUpdateSource() {
	_specs, _name := ps.pop2()

	specs := _specs.comp.(SourceSinkSpecsAST)
	name := _name.comp.(StreamIdentifier)

	s := UpdateSourceStmt{name, specs}
	se := ParsedComponent{_name.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleUpdateSink takes the topmost elements from the stack,
// assuming they are components of a UPDATE SINK statement, and
// replaces them by a single UpdateSinkStmt element.
//
//  SourceSinkSpecsAST
//  StreamIdentifier
//   =>
//  UpdateSinkStmt{StreamIdentifier, SourceSinkSpecsAST}
func (ps *parseStack) AssembleUpdateSink() {
	_specs, _name := ps.pop2()

	specs := _specs.comp.(SourceSinkSpecsAST)
	name := _name.comp.(StreamIdentifier)

	s := UpdateSinkStmt{name, specs}
	se := ParsedComponent{_name.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleInsertIntoFrom takes the topmost elements from the stack,
// assuming they are components of a INSERT ... FROM ... statement, and
// replaces them by a single InsertIntoFromStmt element.
//
//  StreamIdentifier
//  StreamIdentifier
//   =>
//  InsertIntoFromStmt{StreamIdentifier, StreamIdentifier}
func (ps *parseStack) AssembleInsertIntoFrom() {
	_input, _sink := ps.pop2()

	input := _input.comp.(StreamIdentifier)
	sink := _sink.comp.(StreamIdentifier)

	s := InsertIntoFromStmt{sink, input}
	se := ParsedComponent{_sink.begin, _input.end, s}
	ps.Push(&se)
}

// AssemblePauseSource takes the topmost elements from the stack,
// assuming they are components of a PAUSE SOURCE statement, and
// replaces them by a single PauseSourceStmt element.
//
//  StreamIdentifier
//   =>
//  PauseSourceStmt{StreamIdentifier}
func (ps *parseStack) AssemblePauseSource() {
	// pop the components from the stack in reverse order
	_name := ps.Pop()

	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _name.end, PauseSourceStmt{name}}
	ps.Push(&se)
}

// AssembleResumeSource takes the topmost elements from the stack,
// assuming they are components of a RESUME SOURCE statement, and
// replaces them by a single ResumeSourceStmt element.
//
//  StreamIdentifier
//   =>
//  ResumeSourceStmt{StreamIdentifier}
func (ps *parseStack) AssembleResumeSource() {
	// pop the components from the stack in reverse order
	_name := ps.Pop()

	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _name.end, ResumeSourceStmt{name}}
	ps.Push(&se)
}

// AssembleRewindSource takes the topmost elements from the stack,
// assuming they are components of a REWIND SOURCE statement, and
// replaces them by a single RewindSourceStmt element.
//
//  StreamIdentifier
//   =>
//  RewindSourceStmt{StreamIdentifier}
func (ps *parseStack) AssembleRewindSource() {
	// pop the components from the stack in reverse order
	_name := ps.Pop()

	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _name.end, RewindSourceStmt{name}}
	ps.Push(&se)
}

// AssembleDropSource takes the topmost elements from the stack,
// assuming they are components of a DROP SOURCE statement, and
// replaces them by a single DropSourceStmt element.
//
//  StreamIdentifier
//   =>
//  DropSourceStmt{StreamIdentifier}
func (ps *parseStack) AssembleDropSource() {
	// pop the components from the stack in reverse order
	_name := ps.Pop()

	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _name.end, DropSourceStmt{name}}
	ps.Push(&se)
}

// AssembleDropStream takes the topmost elements from the stack,
// assuming they are components of a DROP STREAM statement, and
// replaces them by a single DropStreamStmt element.
//
//  StreamIdentifier
//   =>
//  DropStreamStmt{StreamIdentifier}
func (ps *parseStack) AssembleDropStream() {
	// pop the components from the stack in reverse order
	_name := ps.Pop()

	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _name.end, DropStreamStmt{name}}
	ps.Push(&se)
}

// AssembleDropSink takes the topmost elements from the stack,
// assuming they are components of a DROP SINK statement, and
// replaces them by a single DropSinkStmt element.
//
//  StreamIdentifier
//   =>
//  DropSinkStmt{StreamIdentifier}
func (ps *parseStack) AssembleDropSink() {
	// pop the components from the stack in reverse order
	_name := ps.Pop()

	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _name.end, DropSinkStmt{name}}
	ps.Push(&se)
}

// AssembleDropState takes the topmost elements from the stack,
// assuming they are components of a DROP STATE statement, and
// replaces them by a single DropStateStmt element.
//
//  StreamIdentifier
//   =>
//  DropStateStmt{StreamIdentifier}
func (ps *parseStack) AssembleDropState() {
	// pop the components from the stack in reverse order
	_name := ps.Pop()

	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _name.end, DropStateStmt{name}}
	ps.Push(&se)
}

// AssembleLoadState takes the topmost elements from the stack,
// assuming they are components of a LOAD STATE statement, and
// replaces them by a single LoadStateStmt element.
//
//  SourceSinkSpecsAST
//  Identifier
//  SourceSinkType
//  StreamIdentifier
//   =>
//  LoadStateStmt{StreamIdentifier, SourceSinkType,
//    string, SourceSinkSpecsAST}
func (ps *parseStack) AssembleLoadState() {
	// pop the components from the stack in reverse order
	_specs, _tag, _sinkType, _name := ps.pop4()

	specs := _specs.comp.(SourceSinkSpecsAST)
	tag := _tag.comp.(Identifier)
	sinkType := _sinkType.comp.(SourceSinkType)
	name := _name.comp.(StreamIdentifier)

	s := LoadStateStmt{name, sinkType, string(tag), specs}
	se := ParsedComponent{_name.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleLoadStateOrCreate takes the topmost elements from the stack,
// assuming they are components of a LOAD STATE OR CREATE statement, and
// replaces them by a single LoadStateOrCreateStmt element.
//
//  SourceSinkSpecsAST
//  LoadStateStmt
//   =>
//  LoadStateOrCreateStmt{StreamIdentifier, SourceSinkType,
//    string, SourceSinkSpecsAST, SourceSinkSpecsAST}
func (ps *parseStack) AssembleLoadStateOrCreate() {
	// pop the components from the stack in reverse order
	_createSpecs, _loadStateStmt := ps.pop2()
	loadStateStmt := _loadStateStmt.comp.(LoadStateStmt)

	createSpecs := _createSpecs.comp.(SourceSinkSpecsAST)
	specs := loadStateStmt.SourceSinkSpecsAST
	tag := loadStateStmt.Tag
	sinkType := loadStateStmt.Type
	name := loadStateStmt.Name

	s := LoadStateOrCreateStmt{name, sinkType, tag, specs, createSpecs}
	se := ParsedComponent{_loadStateStmt.begin, _createSpecs.end, s}
	ps.Push(&se)
}

// AssembleSaveState takes the topmost elements from the stack,
// assuming they are components of a SAVE STATE statement, and
// replaces them by a single SaveStateStmt element.
//
//  Identifier
//  StreamIdentifier
//   =>
//  SaveStateStmt{StreamIdentifier, string}
func (ps *parseStack) AssembleSaveState() {
	// pop the components from the stack in reverse order
	_tag, _name := ps.pop2()

	tag := _tag.comp.(Identifier)
	name := _name.comp.(StreamIdentifier)

	se := ParsedComponent{_name.begin, _tag.end, SaveStateStmt{name, string(tag)}}
	ps.Push(&se)
}

// AssembleEval takes the topmost one or two elements from the
// stack, assuming they are components of an EVAL statement, and
// replaces them by a single EvalStmt element.
//
//  Expression
//  Expression
//   =>
//  EvalStmt{Expression, Expression}
// or
//  Expression
//   =>
//  EvalStmt{Expression, nil}
func (ps *parseStack) AssembleEval(begin, end int) {
	var exprBegin int
	var expr Expression
	var inputRow *MapAST

	// pop a different number of items depending on whether we have the
	// optional ON clause or not
	if begin == end {
		// the `... ON input` clause is empty
		_expr := ps.Pop()
		exprBegin = _expr.begin
		expr = _expr.comp.(Expression)
	} else {
		_input, _expr := ps.pop2()
		exprBegin = _expr.begin
		expr = _expr.comp.(Expression)
		input := _input.comp.(MapAST)
		inputRow = &input
	}

	se := ParsedComponent{exprBegin, end, EvalStmt{expr, inputRow}}
	ps.Push(&se)
}

/* Projections/Columns */

// AssembleEmitter takes the topmost elements from the stack, assuming
// they are components of a emitter clause, and replaces them by
// a single EmitterAST element.
//
//  Emitter
//  ...
//   =>
//  EmitterAST{Emitter}
func (ps *parseStack) AssembleEmitter() {
	// pop the components from the stack in reverse order
	_options, _emitter := ps.pop2()

	emitter := _emitter.comp.(Emitter)
	options := _options.comp.([]interface{})

	ps.PushComponent(_emitter.begin, _options.end, EmitterAST{emitter, options})
}

// AssembleEmitterOptions takes the elements from the stack that
// correspond to the input[begin:end] string and pushes a slice
// with all of them back to the stack.
//
//  Any
//  Any
//  Any
//   =>
//  []{Any, Any, Any}
func (ps *parseStack) AssembleEmitterOptions(begin int, end int) {
	elems := ps.collectElements(begin, end)
	if len(elems) == 0 {
		elems = nil
	}
	// push the grouped list back
	ps.PushComponent(begin, end, elems)
}

// AssembleEmitterLimit takes the topmost elements from the stack,
// assuming they are components of a emitter LIMIT option, and replaces
// them by a single EmitterLimit element.
//
//  NumericLiteral
//  ...
//   =>
//  EmitterLimit{NumericLiteral}
func (ps *parseStack) AssembleEmitterLimit() {
	_limit := ps.Pop()

	limit := _limit.comp.(NumericLiteral)

	ps.PushComponent(_limit.begin, _limit.end, EmitterLimit{limit.Value})
}

// AssembleEmitterSampling takes the topmost elements from the stack,
// assuming they are components of a emitter EVERY/SAMPLE option, and replaces
// them by a single EmitterSampling element.
//
//  NumericLiteral
//  ...
//   =>
//  EmitterSampling{NumericLiteral, EmitterSamplingType}
func (ps *parseStack) AssembleEmitterSampling(samplingType EmitterSamplingType, factor float64) {
	_value := ps.Pop()

	var value float64
	if num, ok := _value.comp.(NumericLiteral); ok {
		value = float64(num.Value)
	} else {
		num := _value.comp.(FloatLiteral)
		value = num.Value
	}

	ps.PushComponent(_value.begin, _value.end, EmitterSampling{value * factor, samplingType})
}

// AssembleProjections takes the elements from the stack that
// correspond to the input[begin:end] string and wraps a
// ProjectionsAST struct around them.
//
//  Any
//  Any
//  Any
//   =>
//  ProjectionsAST{[Any, Any, Any]}
func (ps *parseStack) AssembleProjections(begin int, end int) {
	elems := ps.collectElements(begin, end)
	exprs := make([]Expression, len(elems))
	for i := range elems {
		exprs[i] = elems[i].(Expression)
	}
	// push the grouped list back
	ps.PushComponent(begin, end, ProjectionsAST{exprs})
}

// AssembleAlias takes the topmost elements from the stack, assuming
// they are components of an AS clause, and replaces them by
// a single AliasAST element.
//
//  Identifier
//  Any
//   =>
//  AliasAST{Any, Identifier}
func (ps *parseStack) AssembleAlias() {
	// pop the components from the stack in reverse order
	_name, _expr := ps.pop2()

	name := _name.comp.(Identifier)
	expr := _expr.comp.(Expression)

	ps.PushComponent(_expr.begin, _name.end, AliasAST{expr, string(name)})
}

/* FROM clause */

// AssembleWindowedFrom takes the elements from the stack that
// correspond to the input[begin:end] string, makes sure they are all
// AliasedStreamWindowAST elements and wraps a WindowedFromAST struct
// around them. If there are no such elements, adds an
// empty WindowedFromAST struct to the stack.
//
//  AliasedStreamWindowAST
//  AliasedStreamWindowAST
//   =>
//  WindowedFromAST{[AliasedStreamWindowAST, AliasedStreamWindowAST]}
func (ps *parseStack) AssembleWindowedFrom(begin int, end int) {
	if begin == end {
		// push an empty FROM clause
		ps.PushComponent(begin, end, WindowedFromAST{})
	} else {
		elems := ps.collectElements(begin, end)
		rels := make([]AliasedStreamWindowAST, len(elems), len(elems))
		for i, elem := range elems {
			// (if this conversion fails, this is a fundamental parser bug)
			e := elem.(AliasedStreamWindowAST)
			rels[i] = e
		}
		// push the grouped list back
		ps.PushComponent(begin, end, WindowedFromAST{rels})
	}
}

// AssembleInterval takes the topmost elements from the stack, assuming
// they are components of a RANGE clause, and replaces them by
// a single IntervalAST element.
//
//  IntervalUnit
//  NumericLiteral
//   =>
//  IntervalAST{FloatLiteral, IntervalUnit}
// or
//  IntervalUnit
//  FloatLiteral
//   =>
//  IntervalAST{FloatLiteral, IntervalUnit}
func (ps *parseStack) AssembleInterval() {
	// pop the components from the stack in reverse order
	_unit, _num := ps.pop2()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	unit := _unit.comp.(IntervalUnit)
	var val float64
	if num, ok := _num.comp.(NumericLiteral); ok {
		val = float64(num.Value)
	} else {
		num := _num.comp.(FloatLiteral)
		val = num.Value
	}

	// assemble the IntervalAST and push it back
	ps.PushComponent(_num.begin, _unit.end, IntervalAST{FloatLiteral{val}, unit})
}

/* WHERE clause */

// AssembleFilter takes the expression on top of the stack
// (if there is a WHERE clause) and wraps a FilterAST struct
// around it. If there is no WHERE clause, an empty FilterAST
// struct is used.
//
//  Any
//   =>
//  FilterAST{Any}
func (ps *parseStack) AssembleFilter(begin int, end int) {
	if begin == end {
		// push an empty from clause
		ps.PushComponent(begin, end, FilterAST{})
	} else {
		// if the stack is empty at this point, this is
		// a serious parser bug
		f := ps.Pop()
		if begin > f.begin || end < f.end {
			panic("the item on top of the stack is not within given range")
		}
		ps.PushComponent(begin, end, FilterAST{f.comp.(Expression)})
	}
}

/* GROUP BY clause */

// AssembleGrouping takes the elements from the stack that
// correspond to the input[begin:end] string and wraps a
// GroupingAST struct around them. If there are no such elements,
// adds an empty GroupingAST struct to the stack.
//
//  Any
//  Any
//  Any
//   =>
//  GroupingAST{[Any, Any, Any]}
func (ps *parseStack) AssembleGrouping(begin int, end int) {
	elems := ps.collectElements(begin, end)
	var exprs []Expression
	if len(elems) > 0 {
		exprs = make([]Expression, len(elems))
	}
	for i := range elems {
		exprs[i] = elems[i].(Expression)
	}
	// push the grouped list back
	ps.PushComponent(begin, end, GroupingAST{exprs})
}

/* HAVING clause */

// AssembleHaving takes the expression on top of the stack
// (if there is a HAVING clause) and wraps a HavingAST struct
// around it. If there is no HAVING clause, an empty HavingAST
// struct is used.
//
//  Any
//   =>
//  HavingAST{Any}
func (ps *parseStack) AssembleHaving(begin int, end int) {
	if begin == end {
		// push an empty from clause
		ps.PushComponent(begin, end, HavingAST{})
	} else {
		// if the stack is empty at this point, this is
		// a serious parser bug
		h := ps.Pop()
		if begin > h.begin || end < h.end {
			panic("the item on top of the stack is not within given range")
		}
		ps.PushComponent(begin, end, HavingAST{h.comp.(Expression)})
	}
}

// AssembleAliasedStreamWindow takes the topmost elements from the stack, assuming
// they are components of an AS clause, and replaces them by
// a single AliasedStreamWindowAST element.
//
//  Identifier
//  StreamWindowAST
//   =>
//  AliasedStreamWindowAST{StreamWindowAST, Identifier}
func (ps *parseStack) AssembleAliasedStreamWindow() {
	// pop the components from the stack in reverse order
	_name, _rel := ps.pop2()

	name := _name.comp.(Identifier)
	rel := _rel.comp.(StreamWindowAST)

	ps.PushComponent(_rel.begin, _name.end, AliasedStreamWindowAST{rel, string(name)})
}

// EnsureAliasedStreamWindow takes the top element from the stack. If it is a
// StreamWindowAST element, it wraps it into an AliasedStreamWindowAST struct; if it
// is already an AliasedStreamWindowAST it just pushes it back. This helps to
// ensure we only deal with AliasedStreamWindowAST objects in the collection step.
func (ps *parseStack) EnsureAliasedStreamWindow() {
	_elem := ps.Pop()
	elem := _elem.comp

	var aliasRel AliasedStreamWindowAST
	e, ok := elem.(AliasedStreamWindowAST)
	if ok {
		aliasRel = e
	} else {
		e := elem.(StreamWindowAST)
		aliasRel = AliasedStreamWindowAST{e, ""}
	}
	ps.PushComponent(_elem.begin, _elem.end, aliasRel)
}

// AssembleStreamWindow takes the topmost elements from the stack, assuming
// they are components of an AS clause, and replaces them by
// a single StreamWindowAST element.
//
//  IntervalAST
//  Stream
//   =>
//  StreamWindowAST{Stream, IntervalAST}
func (ps *parseStack) AssembleStreamWindow() {
	// pop the components from the stack in reverse order
	_shedding, _capacity, _range, _rel := ps.pop4()

	rel := _rel.comp.(Stream)
	rangeAst := _range.comp.(IntervalAST)
	capacity := _capacity.comp.(NumericLiteral)
	shedding := _shedding.comp.(SheddingOption)

	ps.PushComponent(_rel.begin, _shedding.end, StreamWindowAST{rel, rangeAst,
		capacity.Value, shedding})
}

// AssembleUDSFFuncApp takes the topmost elements from the stack,
// assuming they are components of a UDSF clause, and
// replaces them by a single Stream element.
//
//  FuncAppAST{Function, ExpressionsAST}
//   =>
//  Stream{UDSFStream, Function, ExpressionAST.Expressions}
func (ps *parseStack) AssembleUDSFFuncApp() {
	_fun := ps.Pop()

	fun := _fun.comp.(FuncAppAST)

	se := ParsedComponent{_fun.begin, _fun.end,
		Stream{UDSFStream, string(fun.Function), fun.Expressions}}
	ps.Push(&se)
}

// EnsureCapacitySpec makes sure that the top element of the stack
// is a NumericLiteral element.
func (ps *parseStack) EnsureCapacitySpec(begin int, end int) {
	top := ps.Peek()
	if top == nil || top.end <= begin {
		// there is no item in the given range
		ps.PushComponent(begin, end, NumericLiteral{UnspecifiedCapacity})
	} else {
		// there is an item in the given range
		_, ok := top.comp.(NumericLiteral)
		if !ok {
			panic(fmt.Sprintf("begin (%d) != end (%d), but there "+
				"was a %T on the stack", begin, end, top.comp))
		}
	}
}

// EnsureSheddingSpec makes sure that the top element of the stack
// is a SheddingOption element.
func (ps *parseStack) EnsureSheddingSpec(begin int, end int) {
	top := ps.Peek()
	if top == nil || top.end <= begin {
		// there is no item in the given range
		ps.PushComponent(begin, end, UnspecifiedSheddingOption)
	} else {
		// there is an item in the given range
		_, ok := top.comp.(SheddingOption)
		if !ok {
			panic(fmt.Sprintf("begin (%d) != end (%d), but there "+
				"was a %T on the stack", begin, end, top.comp))
		}
	}
}

// AssembleSourceSinkSpecs takes the elements from the stack that
// correspond to the input[begin:end] string, makes sure
// they are all SourceSinkParamAST elements and wraps a SourceSinkSpecsAST
// struct around them. If there are no such elements, adds an
// empty SourceSpecAST struct to the stack.
//
//  SourceSinkParamAST
//  SourceSinkParamAST
//  SourceSinkParamAST
//   =>
//  SourceSinkSpecsAST{[SourceSpecAST, SourceSpecAST, SourceSpecAST]}
func (ps *parseStack) AssembleSourceSinkSpecs(begin int, end int) {
	if begin == end {
		// push an empty from clause
		ps.PushComponent(begin, end, SourceSinkSpecsAST{})
	} else {
		elems := ps.collectElements(begin, end)
		params := make([]SourceSinkParamAST, len(elems), len(elems))
		for i, elem := range elems {
			// (if this conversion fails, this is a fundamental parser bug)
			e := elem.(SourceSinkParamAST)
			params[i] = e
		}
		// push the grouped list back
		ps.PushComponent(begin, end, SourceSinkSpecsAST{params})
	}
}

// AssembleSourceSinkParam takes the topmost elements from the
// stack, assuming they are part of a WITH clause in a CREATE SOURCE
// statement and replaces them by a single SourceSinkParamAST element.
//
//  Any
//  SourceSinkParamKey
//   =>
//  SourceSinkParamAST{SourceSinkParamKey, data.Value}
func (ps *parseStack) AssembleSourceSinkParam() {
	_value, _key := ps.pop2()

	var toValue func(obj interface{}) data.Value
	toValue = func(obj interface{}) data.Value {
		var value data.Value
		switch lit := obj.(type) {
		default:
			panic(fmt.Sprintf("cannot deal with a %T here", lit))
		case StringLiteral:
			value = data.String(lit.Value)
		case BoolLiteral:
			value = data.Bool(lit.Value)
		case NumericLiteral:
			value = data.Int(lit.Value)
		case FloatLiteral:
			value = data.Float(lit.Value)
		case ArrayAST:
			arr := make(data.Array, len(lit.Expressions))
			for i, item := range lit.Expressions {
				arr[i] = toValue(item)
			}
			value = arr
		case MapAST:
			m := data.Map{}
			for _, item := range lit.Entries {
				m[item.Key] = toValue(item.Value)
			}
			value = m
		}
		return value
	}

	value := toValue(_value.comp)
	key := _key.comp.(SourceSinkParamKey)

	ss := SourceSinkParamAST{key, value}
	ps.PushComponent(_key.begin, _value.end, ss)
}

// EnsureIdentifier makes sure that the top element of the stack
// is an Identifier element. If there was no Identifier on top,
// an empty Identifier is pushed there. Note that since regular
// identifiers must always be empty, an empty identifier is
// semantically equivalent to "not given".
func (ps *parseStack) EnsureIdentifier(begin int, end int) {
	top := ps.Peek()
	if top == nil || top.end <= begin {
		// there is no item in the given range
		ps.PushComponent(begin, end, Identifier(""))
	} else {
		// there is an item in the given range, don't touch it
		_, ok := top.comp.(Identifier)
		if !ok {
			panic(fmt.Sprintf("begin (%d) != end (%d), but there "+
				"was a %T on the stack", begin, end, top.comp))
		}
	}
}

// EnsureKeywordPresent makes sure that the top element of the stack
// is a BinaryKeyword element.
func (ps *parseStack) EnsureKeywordPresent(begin int, end int) {
	top := ps.Peek()
	if top == nil || top.end <= begin {
		// there is no item in the given range
		ps.PushComponent(begin, end, UnspecifiedKeyword)
	} else {
		// there is an item in the given range
		_, ok := top.comp.(BinaryKeyword)
		if !ok {
			panic(fmt.Sprintf("begin (%d) != end (%d), but there "+
				"was a %T on the stack", begin, end, top.comp))
		}
	}
}

/* Expressions */

// AssembleBinaryOperation takes the (odd number of) elements from
// the stack that correspond to the input[begin:end] string and adds
// the given binary operator between each two of them. At the moment,
// all operators are treated as left-associative (which is why the
// parser only allows this for +,-,*,/,OR,AND)
//
//  Any
//   =>
//  Any
// or
//  Any
//  Operator
//  Any
//   =>
//  BinaryOpAST{Operator, Any, Any}
// or
//  Any
//  Operator
//  Any
//  Operator
//  Any
//   =>
//  BinaryOpAST{Operator, BinaryOpAST{Operator, Any, Any}, Any}
func (ps *parseStack) AssembleBinaryOperation(begin int, end int) {
	elems := ps.collectElements(begin, end)
	if len(elems) == 1 {
		// there is no "binary" operation, push back the single element
		ps.PushComponent(begin, end, elems[0])
	} else if len(elems) == 3 {
		op := elems[1].(Operator)
		// connect left and right with the given operator
		ps.PushComponent(begin, end, BinaryOpAST{op, elems[0].(Expression), elems[2].(Expression)})
	} else if len(elems) > 3 {
		op := elems[1].(Operator)
		// left-associativity: process three leftmost items,
		// then nest them with the next item
		leftmost := BinaryOpAST{op, elems[0].(Expression), elems[2].(Expression)}
		for i := 4; i < len(elems); i += 2 {
			leftmost = BinaryOpAST{elems[i-1].(Operator), leftmost, elems[i].(Expression)}
		}
		ps.PushComponent(begin, end, leftmost)
	} else {
		panic(fmt.Sprintf("cannot turn %+v into a binary operation", elems))
	}
}

// AssembleUnaryPrefixOperation takes the two elements from the stack that
// correspond to the input[begin:end] string and adds the given
// unary operator. If there is just one element, push it back unmodified.
//
//  Any
//   =>
//  Any
// or
//  Operator
//  Any
//   =>
//  UnaryOpAST{Operator, Any}
func (ps *parseStack) AssembleUnaryPrefixOperation(begin int, end int) {
	elems := ps.collectElements(begin, end)
	if len(elems) == 1 {
		// there is no operation, push back the single element
		ps.PushComponent(begin, end, elems[0])
	} else if len(elems) == 2 {
		op := elems[0].(Operator)
		// connect the expression with the given operator
		ps.PushComponent(begin, end, UnaryOpAST{op, elems[1].(Expression)})
	} else {
		panic(fmt.Sprintf("cannot turn %+v into a unary operation", elems))
	}
}

// AssembleTypeCast takes the two elements from the stack that
// correspond to the input[begin:end] string and replaces them by
// a single TypeCastAST element. If there is just one element, push
// it back unmodified.
//
//  Any
//   =>
//  Any
// or
//  Any
//  Type
//   =>
//  AssembleTypeCastAST{Any, Type}
func (ps *parseStack) AssembleTypeCast(begin int, end int) {

	elems := ps.collectElements(begin, end)
	if len(elems) == 1 {
		// there is no operation, push back the single element
		ps.PushComponent(begin, end, elems[0])
	} else if len(elems) == 2 {
		target := elems[1].(Type)
		// connect the expression with the given operator
		ps.PushComponent(begin, end, TypeCastAST{elems[0].(Expression), target})
	} else {
		panic(fmt.Sprintf("cannot turn %+v into a type cast", elems))
	}
}

// AssembleFuncApp takes the topmost elements from the stack, assuming
// they are components of a function application clause, and replaces
// them by a single FuncAppAST element.
//
//  ExpressionsAST
//  ExpressionsAST
//  FuncName
//   =>
//  FuncAppAST{FuncName, ExpressionsAST}
func (ps *parseStack) AssembleFuncApp() {
	_ordering, _exprs, _funcName := ps.pop3()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	ordering := _ordering.comp.(ExpressionsAST)
	exprs := _exprs.comp.(ExpressionsAST)
	funcName := _funcName.comp.(FuncName)

	orderExprs := make([]SortedExpressionAST, len(ordering.Expressions))
	for i, e := range ordering.Expressions {
		orderExprs[i] = e.(SortedExpressionAST)
	}
	if len(orderExprs) == 0 {
		orderExprs = nil
	}

	// assemble the FuncAppAST and push it back
	ps.PushComponent(_funcName.begin, _exprs.end, FuncAppAST{funcName, exprs, orderExprs})
}

func (ps *parseStack) AssembleFuncAppSelector() {
	_selector, _func := ps.pop2()

	selector := _selector.comp.(RowValue)
	funcApp := _func.comp.(FuncAppAST)

	ps.PushComponent(_func.begin, _selector.end, FuncAppSelectorAST{
		FuncAppAST: funcApp,
		Selector:   selector,
	})
}

// AssembleSortedExpression takes the topmost elements from the stack,
// assuming they are components of an ORDER BY clause, and replaces
// them by a single SortedExpressionAST element.
//
//  BinaryKeyword
//  Expression
//   =>
//  SortedExpressionAST{Expression, BinaryKeyword}
func (ps *parseStack) AssembleSortedExpression() {
	_sortOrder, _expr := ps.pop2()

	expr := _expr.comp.(Expression)
	sortOrder := _sortOrder.comp.(BinaryKeyword)

	ps.PushComponent(_expr.begin, _sortOrder.end, SortedExpressionAST{expr, sortOrder})
}

// AssembleArray takes the topmost elements from the stack, assuming
// they are components of an array expression, and replaces them
// by a single ArrayAST element.
//
//  ExpressionsAST
//   =>
//  ArrayAST{ExpressionsAST}
func (ps *parseStack) AssembleArray() {
	// pop the components from the stack in reverse order
	_exprs := ps.Pop()

	exprs := _exprs.comp.(ExpressionsAST)

	se := ParsedComponent{_exprs.begin, _exprs.end, ArrayAST{exprs}}
	ps.Push(&se)
}

// AssembleMap takes the elements from the stack that
// correspond to the input[begin:end] string and wraps a
// MapAST struct around them.
//
//  KeyValuePairAST
//  KeyValuePairAST
//  KeyValuePairAST
//   =>
//  MapAST{[KeyValuePairAST, KeyValuePairAST, KeyValuePairAST]}
func (ps *parseStack) AssembleMap(begin int, end int) {
	elems := ps.collectElements(begin, end)
	pairs := make([]KeyValuePairAST, len(elems))
	for i := range elems {
		pairs[i] = elems[i].(KeyValuePairAST)
	}
	// push the grouped list back
	ps.PushComponent(begin, end, MapAST{pairs})
}

// AssembleKeyValuePair takes the topmost elements from the stack,
// assuming they are components of a key-value pair expression, and
// replaces them by a single KeyValuePairAST element.
//
//  Expression
//  StringLiteral
//   =>
//  KeyValuePairAST{string, Expression}
func (ps *parseStack) AssembleKeyValuePair() {
	_expr, _key := ps.pop2()

	expr := _expr.comp.(Expression)
	key := _key.comp.(StringLiteral).Value

	ps.PushComponent(_key.begin, _expr.end, KeyValuePairAST{key, expr})
}

// AssembleConditionCase takes the elements from the stack that
// correspond to the input[begin:end] string plus an additional
// Expression element and wraps an ConditionCaseAST struct around
// them.
//
//  Expression
//  WhenThenPairAST
//  WhenThenPairAST
//   =>
//  ConditionCaseAST{[KeyValuePairAST, KeyValuePairAST], Expression}
//
// or
//
//  WhenThenPairAST
//  WhenThenPairAST
//   =>
//  ConditionCaseAST{[KeyValuePairAST, KeyValuePairAST], nil}
func (ps *parseStack) AssembleConditionCase(begin int, end int) {
	top := ps.Peek()
	// check if the top element is an expression (then it is
	// the ELSE part
	var elseExpr Expression
	if _else, ok := top.comp.(Expression); ok {
		ps.Pop()
		elseExpr = _else
	}
	// collect the WHEN ... THEN pairs
	elems := ps.collectElements(begin, end)
	if len(elems) < 1 {
		panic("no WHEN-THEN pairs on the stack!")
	}
	pairs := make([]WhenThenPairAST, len(elems))
	for i := range elems {
		pairs[i] = elems[i].(WhenThenPairAST)
	}
	// push the grouped list back
	ps.PushComponent(begin, top.end, ConditionCaseAST{pairs, elseExpr})
}

// AssembleExpressionCase takes the elements from the stack that
// correspond to the input[begin:end] string plus the surrounding
// Expression elements and wraps an ExpressionCaseAST struct around
// them.
//
//  Expression
//  WhenThenPairAST
//  WhenThenPairAST
//  Expression
//   =>
//  ExpressionCaseAST{Expression, ConditionCaseAST{[KeyValuePairAST, KeyValuePairAST], Expression}}
//
// or
//
//  WhenThenPairAST
//  WhenThenPairAST
//  Expression
//   =>
//  ExpressionCaseAST{Expression, ConditionCaseAST{[KeyValuePairAST, KeyValuePairAST], nil}}
func (ps *parseStack) AssembleExpressionCase(begin int, end int) {
	// transform the WHEN-THEN and ELSE parts
	ps.AssembleConditionCase(begin, end)
	_cc := ps.Pop()
	cc := _cc.comp.(ConditionCaseAST)
	// get the CASE part
	_case := ps.Pop()
	caseExpr := _case.comp.(Expression)
	// push the grouped list back
	ps.PushComponent(_case.begin, _cc.end, ExpressionCaseAST{caseExpr, cc})
}

// AssembleWhenThenPair takes the topmost elements from the stack,
// assuming they are components of a when-then pair expression, and
// replaces them by a single WhenThenPairAST element.
//
//  Expression
//  Expression
//   =>
//  WhenThenPairAST{Expression, Expression}
func (ps *parseStack) AssembleWhenThenPair() {
	_then, _when := ps.pop2()

	then := _then.comp.(Expression)
	when := _when.comp.(Expression)

	ps.PushComponent(_when.begin, _then.end, WhenThenPairAST{when, then})
}

// AssembleExpressions takes the elements from the stack that
// correspond to the input[begin:end] string and wraps a
// ProjectionsAST struct around them.
//
//  Expression
//  Expression
//  Expression
//   =>
//  ExpressionsAST{[Expression, Expression, Expression]}
func (ps *parseStack) AssembleExpressions(begin int, end int) {
	elems := ps.collectElements(begin, end)
	exprs := make([]Expression, len(elems))
	for i := range elems {
		exprs[i] = elems[i].(Expression)
	}
	// push the grouped list back
	ps.PushComponent(begin, end, ExpressionsAST{exprs})
}

// PushComponent pushes the given component to the top of the stack
// wrapped in a ParsedComponent struct. It's the caller's responsibility
// to make sure that the parameter is one of the AST classes, or there
// will almost surely be a panic at a later point in the parsing process.
func (ps *parseStack) PushComponent(begin int, end int, comp interface{}) {
	if begin > end {
		panic("begin must be less or equal to end")
	}
	if top := ps.Peek(); top != nil && top.end > begin {
		panic("begin must be larger or equal to the previous item's end")
	}
	se := ParsedComponent{begin, end, comp}
	ps.Push(&se)
}

// IncludeTrailingWhitespace updates the `end` value of the top of
// the stack to match the given `end` parameter. This is required
// so that we "eat" trailing comments and do not try to parse them
// into a follow-up statement.
func (ps *parseStack) IncludeTrailingWhitespace(begin int, end int) {
	topElem := ps.Pop()
	topElem.end = end
	ps.Push(topElem)
}

/* helper functions to reduce code duplication */

// collectElements pops all elements with begin/end contained in
// the parameter range from the stack, reverses their order and
// returns them.
func (ps *parseStack) collectElements(begin int, end int) []interface{} {
	elems := []interface{}{}
	// look at elements on the stack as long as there are some and
	// they are contained in our interval
	for ps.Peek() != nil {
		if ps.Peek().end <= begin {
			break
		}
		top := ps.Pop().comp
		elems = append(elems, top)
	}
	// reverse the list to restore original order
	size := len(elems)
	for i := 0; i < size/2; i++ {
		elems[i], elems[size-i-1] = elems[size-i-1], elems[i]
	}
	return elems
}

func (ps *parseStack) pop2() (*ParsedComponent, *ParsedComponent) {
	if ps.Len() < 2 {
		panic("not enough elements on stack to pop 2 of them")
	}
	return ps.Pop(), ps.Pop()
}

func (ps *parseStack) pop3() (*ParsedComponent, *ParsedComponent,
	*ParsedComponent) {
	if ps.Len() < 3 {
		panic("not enough elements on stack to pop 3 of them")
	}
	return ps.Pop(), ps.Pop(), ps.Pop()
}

func (ps *parseStack) pop4() (*ParsedComponent, *ParsedComponent,
	*ParsedComponent, *ParsedComponent) {
	if ps.Len() < 4 {
		panic("not enough elements on stack to pop 4 of them")
	}
	return ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop()
}

func (ps *parseStack) pop5() (*ParsedComponent, *ParsedComponent,
	*ParsedComponent, *ParsedComponent, *ParsedComponent) {
	if ps.Len() < 5 {
		panic("not enough elements on stack to pop 5 of them")
	}
	return ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop()
}

func (ps *parseStack) pop6() (*ParsedComponent, *ParsedComponent,
	*ParsedComponent, *ParsedComponent, *ParsedComponent, *ParsedComponent) {
	if ps.Len() < 6 {
		panic("not enough elements on stack to pop 6 of them")
	}
	return ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop()
}

func (ps *parseStack) pop7() (*ParsedComponent, *ParsedComponent,
	*ParsedComponent, *ParsedComponent, *ParsedComponent, *ParsedComponent,
	*ParsedComponent) {
	if ps.Len() < 7 {
		panic("not enough elements on stack to pop 7 of them")
	}
	return ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop()
}

func (ps *parseStack) pop8() (*ParsedComponent, *ParsedComponent,
	*ParsedComponent, *ParsedComponent, *ParsedComponent, *ParsedComponent,
	*ParsedComponent, *ParsedComponent) {
	if ps.Len() < 8 {
		panic("not enough elements on stack to pop 8 of them")
	}
	return ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop(), ps.Pop()
}

package parser

import (
	"fmt"
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
//  HavingAST
//  GroupingAST
//  FilterAST
//  FromAST
//  ProjectionsAST
//   =>
//  SelectStmt{ProjectionsAST, FromAST, FilterAST, GroupingAST, HavingAST}
func (ps *parseStack) AssembleSelect() {
	// pop the components from the stack in reverse order
	_having, _grouping, _filter, _from, _projections := ps.pop5()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	having := _having.comp.(HavingAST)
	grouping := _grouping.comp.(GroupingAST)
	filter := _filter.comp.(FilterAST)
	from := _from.comp.(FromAST)
	projections := _projections.comp.(ProjectionsAST)

	// assemble the SelectStmt and push it back
	s := SelectStmt{projections, from, filter, grouping, having}
	se := ParsedComponent{_projections.begin, _having.end, s}
	ps.Push(&se)
}

// AssembleCreateStreamAsSelect takes the topmost elements from the stack,
// assuming they are components of a CREATE STREAM statement, and
// replaces them by a single CreateStreamAsSelectStmt element.
//
//  HavingAST
//  GroupingAST
//  FilterAST
//  WindowedFromAST
//  EmitProjectionsAST
//  Relation
//   =>
//  CreateStreamAsSelectStmt{Relation, EmitProjectionsAST, WindowedFromAST, FilterAST,
//    GroupingAST, HavingAST}
func (ps *parseStack) AssembleCreateStreamAsSelect() {
	// pop the components from the stack in reverse order
	_having, _grouping, _filter, _from, _projections, _rel := ps.pop6()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	having := _having.comp.(HavingAST)
	grouping := _grouping.comp.(GroupingAST)
	filter := _filter.comp.(FilterAST)
	from := _from.comp.(WindowedFromAST)
	projections := _projections.comp.(EmitProjectionsAST)
	rel := _rel.comp.(Relation)

	// assemble the SelectStmt and push it back
	s := CreateStreamAsSelectStmt{rel, projections, from, filter, grouping, having}
	se := ParsedComponent{_rel.begin, _having.end, s}
	ps.Push(&se)
}

// AssembleCreateSource takes the topmost elements from the stack,
// assuming they are components of a CREATE SOURCE statement, and
// replaces them by a single CreateSourceStmt element.
//
//  SourceSpecsAST
//  SourceType
//  SourceName
//   =>
//  CreateSourceStmt{Relation, SourceType, SourceSpecsAST}
func (ps *parseStack) AssembleCreateSource() {
	// pop the components from the stack in reverse order
	_specs, _sourceType, _name := ps.pop3()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	specs := _specs.comp.(SourceSpecsAST)
	sourceType := _sourceType.comp.(SourceType)
	name := _name.comp.(SourceName)

	// assemble the CreateSourceStmt and push it back
	s := CreateSourceStmt{name, sourceType, specs}
	se := ParsedComponent{_name.begin, _specs.end, s}
	ps.Push(&se)
}

// AssembleCreateStreamFromSource takes the topmost elements from the stack,
// assuming they are components of a CREATE STREAM statement, and
// replaces them by a single CreateStreamFromSourceStmt element.
//
//  SourceName
//  Relation
//   =>
//  CreateStreamFromSourceStmt{Relation, SourceName}
func (ps *parseStack) AssembleCreateStreamFromSource() {
	_src, _rel := ps.pop2()

	src := _src.comp.(SourceName)
	rel := _rel.comp.(Relation)

	s := CreateStreamFromSourceStmt{rel, src}
	se := ParsedComponent{_rel.begin, _src.end, s}
	ps.Push(&se)
}

/* Projections/Columns */

// AssembleEmitProjections takes the topmost elements from the
// stack, assuming they are part of a SELECT clause in a CREATE STREAM
// statement and replaces them by a single EmitProjectionsAST element.
//
//  ProjectionsAST
//  Emitter
//   =>
//  EmitProjectionsAST{Emitter, ProjectionsAST}
func (ps *parseStack) AssembleEmitProjections() {
	// pop the components from the stack in reverse order
	_projections, _emitter := ps.pop2()

	// extract and convert the contained structure
	projections := _projections.comp.(ProjectionsAST)
	emitter := _emitter.comp.(Emitter)

	// assemble the EmitProjectionsAST and push it back
	ep := EmitProjectionsAST{emitter, projections}
	ps.PushComponent(_emitter.begin, _projections.end, ep)
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
	// push the grouped list back
	ps.PushComponent(begin, end, ProjectionsAST{elems})
}

/* FROM clause */

// AssembleWindowedFrom assumes that the string input[begin:end]
// holds a number of Relation elements followed by a RangeAST, pops all
// of them and turns them into a WindowedFromAST.
//
//  RangeAST
//  Relation
//  Relation
//  Relation
//   =>
//  WindowedFromAST{FromAST, RangeAST}
func (ps *parseStack) AssembleWindowedFrom(begin int, end int) {
	if begin == end {
		// push an empty FROM clause
		ps.PushComponent(begin, end, WindowedFromAST{})
	} else {
		elems := ps.collectElements(begin, end)
		numElems := len(elems)
		if numElems == 0 {
			ps.PushComponent(begin, end, WindowedFromAST{})
		} else {
			// convert the last item to a RangeAST (if this does
			// not work, it is a fundamental parser bug)
			rangeClause := elems[numElems-1].(RangeAST)
			// convert the rest to Relation structs
			rels := make([]Relation, numElems-1, numElems-1)
			for i, elem := range elems[:numElems-1] {
				// (if this conversion fails, this is a fundamental parser bug)
				e := elem.(Relation)
				rels[i] = e
			}
			// push the grouped list wrapped in a FromAST and the RangeAST
			// struct back to the stack
			ps.PushComponent(begin, end, WindowedFromAST{FromAST{rels}, rangeClause})
		}
	}
}

// AssembleRange takes the topmost elements from the stack, assuming
// they are components of a RANGE clause, and replaces them by
// a single RangeAST element.
//
//  RangeUnit
//  Raw
//   =>
//  RangeAST{Raw, RangeUnit}
func (ps *parseStack) AssembleRange() {
	// pop the components from the stack in reverse order
	_unit, _num := ps.pop2()

	// extract and convert the contained structure
	// (if this fails, this is a fundamental parser bug => panic ok)
	unit := _unit.comp.(RangeUnit)
	num := _num.comp.(Raw)

	// assemble the RangeAST and push it back
	ps.PushComponent(_num.begin, _unit.end, RangeAST{num, unit})
}

// AssembleFrom takes the elements from the stack that
// correspond to the input[begin:end] string, makes sure
// they are all Relation elements and wraps a FromAST struct
// around them. If there are no such elements, adds an
// empty FromAST struct to the stack.
//
//  Relation
//  Relation
//  Relation
//   =>
//  FromAST{[Relation, Relation, Relation]}
func (ps *parseStack) AssembleFrom(begin int, end int) {
	if begin == end {
		// push an empty from clause
		ps.PushComponent(begin, end, FromAST{})
	} else {
		elems := ps.collectElements(begin, end)
		rels := make([]Relation, len(elems), len(elems))
		for i, elem := range elems {
			// (if this conversion fails, this is a fundamental parser bug)
			e := elem.(Relation)
			rels[i] = e
		}
		// push the grouped list back
		ps.PushComponent(begin, end, FromAST{rels})
	}
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
		ps.PushComponent(begin, end, FilterAST{f.comp})
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
	// push the grouped list back
	ps.PushComponent(begin, end, GroupingAST{elems})
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
		ps.PushComponent(begin, end, HavingAST{h.comp})
	}
}

// AssembleSourceSpecs takes the elements from the stack that
// correspond to the input[begin:end] string, makes sure
// they are all SourceParamAST elements and wraps a SourceSpecsAST
// struct around them. If there are no such elements, adds an
// empty SourceSpecAST struct to the stack.
//
//  SourceParamAST
//  SourceParamAST
//  SourceParamAST
//   =>
//  SourceSpecsAST{[SourceSpecAST, SourceSpecAST, SourceSpecAST]}
func (ps *parseStack) AssembleSourceSpecs(begin int, end int) {
	if begin == end {
		// push an empty from clause
		ps.PushComponent(begin, end, SourceSpecsAST{})
	} else {
		elems := ps.collectElements(begin, end)
		params := make([]SourceParamAST, len(elems), len(elems))
		for i, elem := range elems {
			// (if this conversion fails, this is a fundamental parser bug)
			e := elem.(SourceParamAST)
			params[i] = e
		}
		// push the grouped list back
		ps.PushComponent(begin, end, SourceSpecsAST{params})
	}
}

// AssembleSourceParam takes the topmost elements from the
// stack, assuming they are part of a WITH clause in a CREATE SOURCE
// statement and replaces them by a single SourceParamAST element.
//
//  SourceParamVal
//  SourceParamKey
//   =>
//  SourceParamAST{SourceParamKey, SourceParamVal}
func (ps *parseStack) AssembleSourceParam() {
	_value, _key := ps.pop2()

	value := _value.comp.(SourceParamVal)
	key := _key.comp.(SourceParamKey)

	ss := SourceParamAST{key, value}
	ps.PushComponent(_key.begin, _value.end, ss)
}

/* Expressions */

// AssembleBinaryOperation takes the two elements from the stack that
// correspond to the input[begin:end] string and adds the given
// binary operator in between. If there is just one element, push
// it back unmodified.
//
//  Any
//  Any
//   =>
//  BinaryOpAST{op, Any, Any}
func (ps *parseStack) AssembleBinaryOperation(begin int, end int, op string) {
	elems := ps.collectElements(begin, end)
	if len(elems) == 1 {
		// there is no "binary" operation, push back the single element
		ps.PushComponent(begin, end, elems[0])
	} else if len(elems) == 2 {
		// connect left and right with the given operator
		ps.PushComponent(begin, end, BinaryOpAST{op, elems[0], elems[1]})
	} else {
		panic(fmt.Sprintf("cannot turn %+v into a binary operation", elems))
	}
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

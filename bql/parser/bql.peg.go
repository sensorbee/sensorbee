package parser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleStatements
	ruleStatement
	ruleSelectStmt
	ruleCreateStreamAsSelectStmt
	ruleCreateSourceStmt
	ruleCreateSinkStmt
	ruleCreateStateStmt
	ruleInsertIntoSelectStmt
	ruleEmitter
	ruleEmitterIntervals
	ruleTimeEmitterInterval
	ruleTupleEmitterInterval
	ruleTupleEmitterFromInterval
	ruleProjections
	ruleProjection
	ruleAliasExpression
	ruleWindowedFrom
	ruleDefWindowedFrom
	ruleInterval
	ruleTimeInterval
	ruleTuplesInterval
	ruleRelations
	ruleDefRelations
	ruleFilter
	ruleGrouping
	ruleGroupList
	ruleHaving
	ruleRelationLike
	ruleDefRelationLike
	ruleAliasedStreamWindow
	ruleDefAliasedStreamWindow
	ruleStreamWindow
	ruleDefStreamWindow
	ruleSourceSinkSpecs
	ruleSourceSinkParam
	ruleSourceSinkParamVal
	ruleExpression
	ruleorExpr
	ruleandExpr
	rulecomparisonExpr
	ruletermExpr
	ruleproductExpr
	rulebaseExpr
	ruleFuncApp
	ruleFuncParams
	ruleLiteral
	ruleComparisonOp
	rulePlusMinusOp
	ruleMultDivOp
	ruleStream
	ruleRowValue
	ruleNumericLiteral
	ruleFloatLiteral
	ruleFunction
	ruleBooleanLiteral
	ruleTRUE
	ruleFALSE
	ruleWildcard
	ruleStringLiteral
	ruleISTREAM
	ruleDSTREAM
	ruleRSTREAM
	ruleTUPLES
	ruleSECONDS
	ruleStreamIdentifier
	ruleSourceSinkType
	ruleSourceSinkParamKey
	ruleOr
	ruleAnd
	ruleEqual
	ruleLess
	ruleLessOrEqual
	ruleGreater
	ruleGreaterOrEqual
	ruleNotEqual
	rulePlus
	ruleMinus
	ruleMultiply
	ruleDivide
	ruleModulo
	ruleIdentifier
	ruleident
	rulesp
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	rulePegText
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35
	ruleAction36
	ruleAction37
	ruleAction38
	ruleAction39
	ruleAction40
	ruleAction41
	ruleAction42
	ruleAction43
	ruleAction44
	ruleAction45
	ruleAction46
	ruleAction47
	ruleAction48
	ruleAction49
	ruleAction50
	ruleAction51
	ruleAction52
	ruleAction53
	ruleAction54
	ruleAction55
	ruleAction56
	ruleAction57
	ruleAction58
	ruleAction59
	ruleAction60
	ruleAction61
	ruleAction62
	ruleAction63
	ruleAction64

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Statements",
	"Statement",
	"SelectStmt",
	"CreateStreamAsSelectStmt",
	"CreateSourceStmt",
	"CreateSinkStmt",
	"CreateStateStmt",
	"InsertIntoSelectStmt",
	"Emitter",
	"EmitterIntervals",
	"TimeEmitterInterval",
	"TupleEmitterInterval",
	"TupleEmitterFromInterval",
	"Projections",
	"Projection",
	"AliasExpression",
	"WindowedFrom",
	"DefWindowedFrom",
	"Interval",
	"TimeInterval",
	"TuplesInterval",
	"Relations",
	"DefRelations",
	"Filter",
	"Grouping",
	"GroupList",
	"Having",
	"RelationLike",
	"DefRelationLike",
	"AliasedStreamWindow",
	"DefAliasedStreamWindow",
	"StreamWindow",
	"DefStreamWindow",
	"SourceSinkSpecs",
	"SourceSinkParam",
	"SourceSinkParamVal",
	"Expression",
	"orExpr",
	"andExpr",
	"comparisonExpr",
	"termExpr",
	"productExpr",
	"baseExpr",
	"FuncApp",
	"FuncParams",
	"Literal",
	"ComparisonOp",
	"PlusMinusOp",
	"MultDivOp",
	"Stream",
	"RowValue",
	"NumericLiteral",
	"FloatLiteral",
	"Function",
	"BooleanLiteral",
	"TRUE",
	"FALSE",
	"Wildcard",
	"StringLiteral",
	"ISTREAM",
	"DSTREAM",
	"RSTREAM",
	"TUPLES",
	"SECONDS",
	"StreamIdentifier",
	"SourceSinkType",
	"SourceSinkParamKey",
	"Or",
	"And",
	"Equal",
	"Less",
	"LessOrEqual",
	"Greater",
	"GreaterOrEqual",
	"NotEqual",
	"Plus",
	"Minus",
	"Multiply",
	"Divide",
	"Modulo",
	"Identifier",
	"ident",
	"sp",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"PegText",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
	"Action39",
	"Action40",
	"Action41",
	"Action42",
	"Action43",
	"Action44",
	"Action45",
	"Action46",
	"Action47",
	"Action48",
	"Action49",
	"Action50",
	"Action51",
	"Action52",
	"Action53",
	"Action54",
	"Action55",
	"Action56",
	"Action57",
	"Action58",
	"Action59",
	"Action60",
	"Action61",
	"Action62",
	"Action63",
	"Action64",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next uint32, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/*func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2 * len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}*/

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type bqlPeg struct {
	parseStack

	Buffer string
	buffer []rune
	rules  [150]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer[0:] {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *bqlPeg
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *bqlPeg) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *bqlPeg) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *bqlPeg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.AssembleSelect()

		case ruleAction1:

			p.AssembleCreateStreamAsSelect()

		case ruleAction2:

			p.AssembleCreateSource()

		case ruleAction3:

			p.AssembleCreateSink()

		case ruleAction4:

			p.AssembleCreateState()

		case ruleAction5:

			p.AssembleInsertIntoSelect()

		case ruleAction6:

			p.AssembleEmitter(begin, end)

		case ruleAction7:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction8:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction9:

			p.AssembleStreamEmitInterval()

		case ruleAction10:

			p.AssembleProjections(begin, end)

		case ruleAction11:

			p.AssembleAlias()

		case ruleAction12:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction13:

			p.AssembleWindowedFrom(begin, end)

		case ruleAction14:

			p.AssembleInterval()

		case ruleAction15:

			p.AssembleInterval()

		case ruleAction16:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction17:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction18:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction19:

			p.EnsureAliasedStreamWindow()

		case ruleAction20:

			p.EnsureAliasedStreamWindow()

		case ruleAction21:

			p.AssembleAliasedStreamWindow()

		case ruleAction22:

			p.AssembleAliasedStreamWindow()

		case ruleAction23:

			p.AssembleStreamWindow()

		case ruleAction24:

			p.AssembleStreamWindow()

		case ruleAction25:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction26:

			p.AssembleSourceSinkParam()

		case ruleAction27:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction28:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction29:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction30:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction31:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction32:

			p.AssembleFuncApp()

		case ruleAction33:

			p.AssembleExpressions(begin, end)

		case ruleAction34:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction35:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction36:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction37:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction38:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction39:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction40:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction41:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction43:

			p.PushComponent(begin, end, Istream)

		case ruleAction44:

			p.PushComponent(begin, end, Dstream)

		case ruleAction45:

			p.PushComponent(begin, end, Rstream)

		case ruleAction46:

			p.PushComponent(begin, end, Tuples)

		case ruleAction47:

			p.PushComponent(begin, end, Seconds)

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction51:

			p.PushComponent(begin, end, Or)

		case ruleAction52:

			p.PushComponent(begin, end, And)

		case ruleAction53:

			p.PushComponent(begin, end, Equal)

		case ruleAction54:

			p.PushComponent(begin, end, Less)

		case ruleAction55:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction56:

			p.PushComponent(begin, end, Greater)

		case ruleAction57:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction58:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction59:

			p.PushComponent(begin, end, Plus)

		case ruleAction60:

			p.PushComponent(begin, end, Minus)

		case ruleAction61:

			p.PushComponent(begin, end, Multiply)

		case ruleAction62:

			p.PushComponent(begin, end, Divide)

		case ruleAction63:

			p.PushComponent(begin, end, Modulo)

		case ruleAction64:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		}
	}
	_, _, _, _ = buffer, text, begin, end
}

func (p *bqlPeg) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens32{tree: make([]token32, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Statements <- <(sp ((Statement sp ';' .*) / Statement) !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[rulesp]() {
					goto l0
				}
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[ruleStatement]() {
						goto l3
					}
					if !_rules[rulesp]() {
						goto l3
					}
					if buffer[position] != rune(';') {
						goto l3
					}
					position++
				l4:
					{
						position5, tokenIndex5, depth5 := position, tokenIndex, depth
						if !matchDot() {
							goto l5
						}
						goto l4
					l5:
						position, tokenIndex, depth = position5, tokenIndex5, depth5
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleStatement]() {
						goto l0
					}
				}
			l2:
				{
					position6, tokenIndex6, depth6 := position, tokenIndex, depth
					if !matchDot() {
						goto l6
					}
					goto l0
				l6:
					position, tokenIndex, depth = position6, tokenIndex6, depth6
				}
				depth--
				add(ruleStatements, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / CreateStateStmt)> */
		func() bool {
			position7, tokenIndex7, depth7 := position, tokenIndex, depth
			{
				position8 := position
				depth++
				{
					position9, tokenIndex9, depth9 := position, tokenIndex, depth
					if !_rules[ruleSelectStmt]() {
						goto l10
					}
					goto l9
				l10:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l11
					}
					goto l9
				l11:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleCreateSourceStmt]() {
						goto l12
					}
					goto l9
				l12:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleCreateSinkStmt]() {
						goto l13
					}
					goto l9
				l13:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l14
					}
					goto l9
				l14:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleCreateStateStmt]() {
						goto l7
					}
				}
			l9:
				depth--
				add(ruleStatement, position8)
			}
			return true
		l7:
			position, tokenIndex, depth = position7, tokenIndex7, depth7
			return false
		},
		/* 2 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter? sp Projections sp DefWindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
		func() bool {
			position15, tokenIndex15, depth15 := position, tokenIndex, depth
			{
				position16 := position
				depth++
				{
					position17, tokenIndex17, depth17 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l18
					}
					position++
					goto l17
				l18:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if buffer[position] != rune('S') {
						goto l15
					}
					position++
				}
			l17:
				{
					position19, tokenIndex19, depth19 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l20
					}
					position++
					goto l19
				l20:
					position, tokenIndex, depth = position19, tokenIndex19, depth19
					if buffer[position] != rune('E') {
						goto l15
					}
					position++
				}
			l19:
				{
					position21, tokenIndex21, depth21 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l22
					}
					position++
					goto l21
				l22:
					position, tokenIndex, depth = position21, tokenIndex21, depth21
					if buffer[position] != rune('L') {
						goto l15
					}
					position++
				}
			l21:
				{
					position23, tokenIndex23, depth23 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l24
					}
					position++
					goto l23
				l24:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if buffer[position] != rune('E') {
						goto l15
					}
					position++
				}
			l23:
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l26
					}
					position++
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if buffer[position] != rune('C') {
						goto l15
					}
					position++
				}
			l25:
				{
					position27, tokenIndex27, depth27 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l28
					}
					position++
					goto l27
				l28:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					if buffer[position] != rune('T') {
						goto l15
					}
					position++
				}
			l27:
				if !_rules[rulesp]() {
					goto l15
				}
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if !_rules[ruleEmitter]() {
						goto l29
					}
					goto l30
				l29:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
				}
			l30:
				if !_rules[rulesp]() {
					goto l15
				}
				if !_rules[ruleProjections]() {
					goto l15
				}
				if !_rules[rulesp]() {
					goto l15
				}
				if !_rules[ruleDefWindowedFrom]() {
					goto l15
				}
				if !_rules[rulesp]() {
					goto l15
				}
				if !_rules[ruleFilter]() {
					goto l15
				}
				if !_rules[rulesp]() {
					goto l15
				}
				if !_rules[ruleGrouping]() {
					goto l15
				}
				if !_rules[rulesp]() {
					goto l15
				}
				if !_rules[ruleHaving]() {
					goto l15
				}
				if !_rules[rulesp]() {
					goto l15
				}
				if !_rules[ruleAction0]() {
					goto l15
				}
				depth--
				add(ruleSelectStmt, position16)
			}
			return true
		l15:
			position, tokenIndex, depth = position15, tokenIndex15, depth15
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position31, tokenIndex31, depth31 := position, tokenIndex, depth
			{
				position32 := position
				depth++
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l34
					}
					position++
					goto l33
				l34:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
					if buffer[position] != rune('C') {
						goto l31
					}
					position++
				}
			l33:
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l36
					}
					position++
					goto l35
				l36:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if buffer[position] != rune('R') {
						goto l31
					}
					position++
				}
			l35:
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l38
					}
					position++
					goto l37
				l38:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if buffer[position] != rune('E') {
						goto l31
					}
					position++
				}
			l37:
				{
					position39, tokenIndex39, depth39 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l40
					}
					position++
					goto l39
				l40:
					position, tokenIndex, depth = position39, tokenIndex39, depth39
					if buffer[position] != rune('A') {
						goto l31
					}
					position++
				}
			l39:
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l42
					}
					position++
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if buffer[position] != rune('T') {
						goto l31
					}
					position++
				}
			l41:
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l44
					}
					position++
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if buffer[position] != rune('E') {
						goto l31
					}
					position++
				}
			l43:
				if !_rules[rulesp]() {
					goto l31
				}
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l46
					}
					position++
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if buffer[position] != rune('S') {
						goto l31
					}
					position++
				}
			l45:
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l48
					}
					position++
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					if buffer[position] != rune('T') {
						goto l31
					}
					position++
				}
			l47:
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l50
					}
					position++
					goto l49
				l50:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
					if buffer[position] != rune('R') {
						goto l31
					}
					position++
				}
			l49:
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l52
					}
					position++
					goto l51
				l52:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
					if buffer[position] != rune('E') {
						goto l31
					}
					position++
				}
			l51:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					if buffer[position] != rune('A') {
						goto l31
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
					if buffer[position] != rune('M') {
						goto l31
					}
					position++
				}
			l55:
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l31
				}
				if !_rules[rulesp]() {
					goto l31
				}
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
					if buffer[position] != rune('A') {
						goto l31
					}
					position++
				}
			l57:
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l60
					}
					position++
					goto l59
				l60:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
					if buffer[position] != rune('S') {
						goto l31
					}
					position++
				}
			l59:
				if !_rules[rulesp]() {
					goto l31
				}
				{
					position61, tokenIndex61, depth61 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l62
					}
					position++
					goto l61
				l62:
					position, tokenIndex, depth = position61, tokenIndex61, depth61
					if buffer[position] != rune('S') {
						goto l31
					}
					position++
				}
			l61:
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l64
					}
					position++
					goto l63
				l64:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
					if buffer[position] != rune('E') {
						goto l31
					}
					position++
				}
			l63:
				{
					position65, tokenIndex65, depth65 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l66
					}
					position++
					goto l65
				l66:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
					if buffer[position] != rune('L') {
						goto l31
					}
					position++
				}
			l65:
				{
					position67, tokenIndex67, depth67 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l68
					}
					position++
					goto l67
				l68:
					position, tokenIndex, depth = position67, tokenIndex67, depth67
					if buffer[position] != rune('E') {
						goto l31
					}
					position++
				}
			l67:
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l70
					}
					position++
					goto l69
				l70:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
					if buffer[position] != rune('C') {
						goto l31
					}
					position++
				}
			l69:
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l72
					}
					position++
					goto l71
				l72:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					if buffer[position] != rune('T') {
						goto l31
					}
					position++
				}
			l71:
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleEmitter]() {
					goto l31
				}
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleProjections]() {
					goto l31
				}
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleWindowedFrom]() {
					goto l31
				}
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleFilter]() {
					goto l31
				}
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleGrouping]() {
					goto l31
				}
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleHaving]() {
					goto l31
				}
				if !_rules[rulesp]() {
					goto l31
				}
				if !_rules[ruleAction1]() {
					goto l31
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position32)
			}
			return true
		l31:
			position, tokenIndex, depth = position31, tokenIndex31, depth31
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position73, tokenIndex73, depth73 := position, tokenIndex, depth
			{
				position74 := position
				depth++
				{
					position75, tokenIndex75, depth75 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l76
					}
					position++
					goto l75
				l76:
					position, tokenIndex, depth = position75, tokenIndex75, depth75
					if buffer[position] != rune('C') {
						goto l73
					}
					position++
				}
			l75:
				{
					position77, tokenIndex77, depth77 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l78
					}
					position++
					goto l77
				l78:
					position, tokenIndex, depth = position77, tokenIndex77, depth77
					if buffer[position] != rune('R') {
						goto l73
					}
					position++
				}
			l77:
				{
					position79, tokenIndex79, depth79 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l80
					}
					position++
					goto l79
				l80:
					position, tokenIndex, depth = position79, tokenIndex79, depth79
					if buffer[position] != rune('E') {
						goto l73
					}
					position++
				}
			l79:
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l82
					}
					position++
					goto l81
				l82:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if buffer[position] != rune('A') {
						goto l73
					}
					position++
				}
			l81:
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l84
					}
					position++
					goto l83
				l84:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
					if buffer[position] != rune('T') {
						goto l73
					}
					position++
				}
			l83:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l86
					}
					position++
					goto l85
				l86:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
					if buffer[position] != rune('E') {
						goto l73
					}
					position++
				}
			l85:
				if !_rules[rulesp]() {
					goto l73
				}
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l88
					}
					position++
					goto l87
				l88:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
					if buffer[position] != rune('S') {
						goto l73
					}
					position++
				}
			l87:
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l90
					}
					position++
					goto l89
				l90:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
					if buffer[position] != rune('O') {
						goto l73
					}
					position++
				}
			l89:
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l92
					}
					position++
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					if buffer[position] != rune('U') {
						goto l73
					}
					position++
				}
			l91:
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('R') {
						goto l73
					}
					position++
				}
			l93:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if buffer[position] != rune('C') {
						goto l73
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('E') {
						goto l73
					}
					position++
				}
			l97:
				if !_rules[rulesp]() {
					goto l73
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l73
				}
				if !_rules[rulesp]() {
					goto l73
				}
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('T') {
						goto l73
					}
					position++
				}
			l99:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('Y') {
						goto l73
					}
					position++
				}
			l101:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if buffer[position] != rune('P') {
						goto l73
					}
					position++
				}
			l103:
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('E') {
						goto l73
					}
					position++
				}
			l105:
				if !_rules[rulesp]() {
					goto l73
				}
				if !_rules[ruleSourceSinkType]() {
					goto l73
				}
				if !_rules[rulesp]() {
					goto l73
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l73
				}
				if !_rules[ruleAction2]() {
					goto l73
				}
				depth--
				add(ruleCreateSourceStmt, position74)
			}
			return true
		l73:
			position, tokenIndex, depth = position73, tokenIndex73, depth73
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position107, tokenIndex107, depth107 := position, tokenIndex, depth
			{
				position108 := position
				depth++
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('C') {
						goto l107
					}
					position++
				}
			l109:
				{
					position111, tokenIndex111, depth111 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l112
					}
					position++
					goto l111
				l112:
					position, tokenIndex, depth = position111, tokenIndex111, depth111
					if buffer[position] != rune('R') {
						goto l107
					}
					position++
				}
			l111:
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l114
					}
					position++
					goto l113
				l114:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
					if buffer[position] != rune('E') {
						goto l107
					}
					position++
				}
			l113:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('A') {
						goto l107
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('T') {
						goto l107
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('E') {
						goto l107
					}
					position++
				}
			l119:
				if !_rules[rulesp]() {
					goto l107
				}
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if buffer[position] != rune('S') {
						goto l107
					}
					position++
				}
			l121:
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('I') {
						goto l107
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('N') {
						goto l107
					}
					position++
				}
			l125:
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('K') {
						goto l107
					}
					position++
				}
			l127:
				if !_rules[rulesp]() {
					goto l107
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l107
				}
				if !_rules[rulesp]() {
					goto l107
				}
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('T') {
						goto l107
					}
					position++
				}
			l129:
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('Y') {
						goto l107
					}
					position++
				}
			l131:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('P') {
						goto l107
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('E') {
						goto l107
					}
					position++
				}
			l135:
				if !_rules[rulesp]() {
					goto l107
				}
				if !_rules[ruleSourceSinkType]() {
					goto l107
				}
				if !_rules[rulesp]() {
					goto l107
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l107
				}
				if !_rules[ruleAction3]() {
					goto l107
				}
				depth--
				add(ruleCreateSinkStmt, position108)
			}
			return true
		l107:
			position, tokenIndex, depth = position107, tokenIndex107, depth107
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position137, tokenIndex137, depth137 := position, tokenIndex, depth
			{
				position138 := position
				depth++
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('C') {
						goto l137
					}
					position++
				}
			l139:
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l142
					}
					position++
					goto l141
				l142:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					if buffer[position] != rune('R') {
						goto l137
					}
					position++
				}
			l141:
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					if buffer[position] != rune('E') {
						goto l137
					}
					position++
				}
			l143:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('A') {
						goto l137
					}
					position++
				}
			l145:
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('T') {
						goto l137
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('E') {
						goto l137
					}
					position++
				}
			l149:
				if !_rules[rulesp]() {
					goto l137
				}
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('S') {
						goto l137
					}
					position++
				}
			l151:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('T') {
						goto l137
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('A') {
						goto l137
					}
					position++
				}
			l155:
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('T') {
						goto l137
					}
					position++
				}
			l157:
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l160
					}
					position++
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if buffer[position] != rune('E') {
						goto l137
					}
					position++
				}
			l159:
				if !_rules[rulesp]() {
					goto l137
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l137
				}
				if !_rules[rulesp]() {
					goto l137
				}
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l162
					}
					position++
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if buffer[position] != rune('T') {
						goto l137
					}
					position++
				}
			l161:
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l164
					}
					position++
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if buffer[position] != rune('Y') {
						goto l137
					}
					position++
				}
			l163:
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('P') {
						goto l137
					}
					position++
				}
			l165:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('E') {
						goto l137
					}
					position++
				}
			l167:
				if !_rules[rulesp]() {
					goto l137
				}
				if !_rules[ruleSourceSinkType]() {
					goto l137
				}
				if !_rules[rulesp]() {
					goto l137
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l137
				}
				if !_rules[ruleAction4]() {
					goto l137
				}
				depth--
				add(ruleCreateStateStmt, position138)
			}
			return true
		l137:
			position, tokenIndex, depth = position137, tokenIndex137, depth137
			return false
		},
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action5)> */
		func() bool {
			position169, tokenIndex169, depth169 := position, tokenIndex, depth
			{
				position170 := position
				depth++
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('I') {
						goto l169
					}
					position++
				}
			l171:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('N') {
						goto l169
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('S') {
						goto l169
					}
					position++
				}
			l175:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('E') {
						goto l169
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('R') {
						goto l169
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('T') {
						goto l169
					}
					position++
				}
			l181:
				if !_rules[rulesp]() {
					goto l169
				}
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('I') {
						goto l169
					}
					position++
				}
			l183:
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('N') {
						goto l169
					}
					position++
				}
			l185:
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if buffer[position] != rune('T') {
						goto l169
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('O') {
						goto l169
					}
					position++
				}
			l189:
				if !_rules[rulesp]() {
					goto l169
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l169
				}
				if !_rules[rulesp]() {
					goto l169
				}
				if !_rules[ruleSelectStmt]() {
					goto l169
				}
				if !_rules[ruleAction5]() {
					goto l169
				}
				depth--
				add(ruleInsertIntoSelectStmt, position170)
			}
			return true
		l169:
			position, tokenIndex, depth = position169, tokenIndex169, depth169
			return false
		},
		/* 8 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action6)> */
		func() bool {
			position191, tokenIndex191, depth191 := position, tokenIndex, depth
			{
				position192 := position
				depth++
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l194
					}
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if !_rules[ruleDSTREAM]() {
						goto l195
					}
					goto l193
				l195:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if !_rules[ruleRSTREAM]() {
						goto l191
					}
				}
			l193:
				{
					position196 := position
					depth++
					{
						position197, tokenIndex197, depth197 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l197
						}
						if buffer[position] != rune('[') {
							goto l197
						}
						position++
						if !_rules[rulesp]() {
							goto l197
						}
						{
							position199, tokenIndex199, depth199 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l200
							}
							position++
							goto l199
						l200:
							position, tokenIndex, depth = position199, tokenIndex199, depth199
							if buffer[position] != rune('E') {
								goto l197
							}
							position++
						}
					l199:
						{
							position201, tokenIndex201, depth201 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l202
							}
							position++
							goto l201
						l202:
							position, tokenIndex, depth = position201, tokenIndex201, depth201
							if buffer[position] != rune('V') {
								goto l197
							}
							position++
						}
					l201:
						{
							position203, tokenIndex203, depth203 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l204
							}
							position++
							goto l203
						l204:
							position, tokenIndex, depth = position203, tokenIndex203, depth203
							if buffer[position] != rune('E') {
								goto l197
							}
							position++
						}
					l203:
						{
							position205, tokenIndex205, depth205 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l206
							}
							position++
							goto l205
						l206:
							position, tokenIndex, depth = position205, tokenIndex205, depth205
							if buffer[position] != rune('R') {
								goto l197
							}
							position++
						}
					l205:
						{
							position207, tokenIndex207, depth207 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l208
							}
							position++
							goto l207
						l208:
							position, tokenIndex, depth = position207, tokenIndex207, depth207
							if buffer[position] != rune('Y') {
								goto l197
							}
							position++
						}
					l207:
						if !_rules[rulesp]() {
							goto l197
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l197
						}
						if !_rules[rulesp]() {
							goto l197
						}
						if buffer[position] != rune(']') {
							goto l197
						}
						position++
						goto l198
					l197:
						position, tokenIndex, depth = position197, tokenIndex197, depth197
					}
				l198:
					depth--
					add(rulePegText, position196)
				}
				if !_rules[ruleAction6]() {
					goto l191
				}
				depth--
				add(ruleEmitter, position192)
			}
			return true
		l191:
			position, tokenIndex, depth = position191, tokenIndex191, depth191
			return false
		},
		/* 9 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position209, tokenIndex209, depth209 := position, tokenIndex, depth
			{
				position210 := position
				depth++
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l212
					}
				l213:
					{
						position214, tokenIndex214, depth214 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l214
						}
						if buffer[position] != rune(',') {
							goto l214
						}
						position++
						if !_rules[rulesp]() {
							goto l214
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l214
						}
						goto l213
					l214:
						position, tokenIndex, depth = position214, tokenIndex214, depth214
					}
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if !_rules[ruleTimeEmitterInterval]() {
						goto l215
					}
					goto l211
				l215:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if !_rules[ruleTupleEmitterInterval]() {
						goto l209
					}
				}
			l211:
				depth--
				add(ruleEmitterIntervals, position210)
			}
			return true
		l209:
			position, tokenIndex, depth = position209, tokenIndex209, depth209
			return false
		},
		/* 10 TimeEmitterInterval <- <(<TimeInterval> Action7)> */
		func() bool {
			position216, tokenIndex216, depth216 := position, tokenIndex, depth
			{
				position217 := position
				depth++
				{
					position218 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l216
					}
					depth--
					add(rulePegText, position218)
				}
				if !_rules[ruleAction7]() {
					goto l216
				}
				depth--
				add(ruleTimeEmitterInterval, position217)
			}
			return true
		l216:
			position, tokenIndex, depth = position216, tokenIndex216, depth216
			return false
		},
		/* 11 TupleEmitterInterval <- <(<TuplesInterval> Action8)> */
		func() bool {
			position219, tokenIndex219, depth219 := position, tokenIndex, depth
			{
				position220 := position
				depth++
				{
					position221 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l219
					}
					depth--
					add(rulePegText, position221)
				}
				if !_rules[ruleAction8]() {
					goto l219
				}
				depth--
				add(ruleTupleEmitterInterval, position220)
			}
			return true
		l219:
			position, tokenIndex, depth = position219, tokenIndex219, depth219
			return false
		},
		/* 12 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action9)> */
		func() bool {
			position222, tokenIndex222, depth222 := position, tokenIndex, depth
			{
				position223 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l222
				}
				if !_rules[rulesp]() {
					goto l222
				}
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('I') {
						goto l222
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('N') {
						goto l222
					}
					position++
				}
			l226:
				if !_rules[rulesp]() {
					goto l222
				}
				if !_rules[ruleStream]() {
					goto l222
				}
				if !_rules[ruleAction9]() {
					goto l222
				}
				depth--
				add(ruleTupleEmitterFromInterval, position223)
			}
			return true
		l222:
			position, tokenIndex, depth = position222, tokenIndex222, depth222
			return false
		},
		/* 13 Projections <- <(<(Projection sp (',' sp Projection)*)> Action10)> */
		func() bool {
			position228, tokenIndex228, depth228 := position, tokenIndex, depth
			{
				position229 := position
				depth++
				{
					position230 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l228
					}
					if !_rules[rulesp]() {
						goto l228
					}
				l231:
					{
						position232, tokenIndex232, depth232 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l232
						}
						position++
						if !_rules[rulesp]() {
							goto l232
						}
						if !_rules[ruleProjection]() {
							goto l232
						}
						goto l231
					l232:
						position, tokenIndex, depth = position232, tokenIndex232, depth232
					}
					depth--
					add(rulePegText, position230)
				}
				if !_rules[ruleAction10]() {
					goto l228
				}
				depth--
				add(ruleProjections, position229)
			}
			return true
		l228:
			position, tokenIndex, depth = position228, tokenIndex228, depth228
			return false
		},
		/* 14 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position233, tokenIndex233, depth233 := position, tokenIndex, depth
			{
				position234 := position
				depth++
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l236
					}
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if !_rules[ruleExpression]() {
						goto l237
					}
					goto l235
				l237:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if !_rules[ruleWildcard]() {
						goto l233
					}
				}
			l235:
				depth--
				add(ruleProjection, position234)
			}
			return true
		l233:
			position, tokenIndex, depth = position233, tokenIndex233, depth233
			return false
		},
		/* 15 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action11)> */
		func() bool {
			position238, tokenIndex238, depth238 := position, tokenIndex, depth
			{
				position239 := position
				depth++
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l241
					}
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if !_rules[ruleWildcard]() {
						goto l238
					}
				}
			l240:
				if !_rules[rulesp]() {
					goto l238
				}
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('A') {
						goto l238
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('S') {
						goto l238
					}
					position++
				}
			l244:
				if !_rules[rulesp]() {
					goto l238
				}
				if !_rules[ruleIdentifier]() {
					goto l238
				}
				if !_rules[ruleAction11]() {
					goto l238
				}
				depth--
				add(ruleAliasExpression, position239)
			}
			return true
		l238:
			position, tokenIndex, depth = position238, tokenIndex238, depth238
			return false
		},
		/* 16 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action12)> */
		func() bool {
			position246, tokenIndex246, depth246 := position, tokenIndex, depth
			{
				position247 := position
				depth++
				{
					position248 := position
					depth++
					{
						position249, tokenIndex249, depth249 := position, tokenIndex, depth
						{
							position251, tokenIndex251, depth251 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l252
							}
							position++
							goto l251
						l252:
							position, tokenIndex, depth = position251, tokenIndex251, depth251
							if buffer[position] != rune('F') {
								goto l249
							}
							position++
						}
					l251:
						{
							position253, tokenIndex253, depth253 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l254
							}
							position++
							goto l253
						l254:
							position, tokenIndex, depth = position253, tokenIndex253, depth253
							if buffer[position] != rune('R') {
								goto l249
							}
							position++
						}
					l253:
						{
							position255, tokenIndex255, depth255 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l256
							}
							position++
							goto l255
						l256:
							position, tokenIndex, depth = position255, tokenIndex255, depth255
							if buffer[position] != rune('O') {
								goto l249
							}
							position++
						}
					l255:
						{
							position257, tokenIndex257, depth257 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l258
							}
							position++
							goto l257
						l258:
							position, tokenIndex, depth = position257, tokenIndex257, depth257
							if buffer[position] != rune('M') {
								goto l249
							}
							position++
						}
					l257:
						if !_rules[rulesp]() {
							goto l249
						}
						if !_rules[ruleRelations]() {
							goto l249
						}
						if !_rules[rulesp]() {
							goto l249
						}
						goto l250
					l249:
						position, tokenIndex, depth = position249, tokenIndex249, depth249
					}
				l250:
					depth--
					add(rulePegText, position248)
				}
				if !_rules[ruleAction12]() {
					goto l246
				}
				depth--
				add(ruleWindowedFrom, position247)
			}
			return true
		l246:
			position, tokenIndex, depth = position246, tokenIndex246, depth246
			return false
		},
		/* 17 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action13)> */
		func() bool {
			position259, tokenIndex259, depth259 := position, tokenIndex, depth
			{
				position260 := position
				depth++
				{
					position261 := position
					depth++
					{
						position262, tokenIndex262, depth262 := position, tokenIndex, depth
						{
							position264, tokenIndex264, depth264 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l265
							}
							position++
							goto l264
						l265:
							position, tokenIndex, depth = position264, tokenIndex264, depth264
							if buffer[position] != rune('F') {
								goto l262
							}
							position++
						}
					l264:
						{
							position266, tokenIndex266, depth266 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l267
							}
							position++
							goto l266
						l267:
							position, tokenIndex, depth = position266, tokenIndex266, depth266
							if buffer[position] != rune('R') {
								goto l262
							}
							position++
						}
					l266:
						{
							position268, tokenIndex268, depth268 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l269
							}
							position++
							goto l268
						l269:
							position, tokenIndex, depth = position268, tokenIndex268, depth268
							if buffer[position] != rune('O') {
								goto l262
							}
							position++
						}
					l268:
						{
							position270, tokenIndex270, depth270 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l271
							}
							position++
							goto l270
						l271:
							position, tokenIndex, depth = position270, tokenIndex270, depth270
							if buffer[position] != rune('M') {
								goto l262
							}
							position++
						}
					l270:
						if !_rules[rulesp]() {
							goto l262
						}
						if !_rules[ruleDefRelations]() {
							goto l262
						}
						if !_rules[rulesp]() {
							goto l262
						}
						goto l263
					l262:
						position, tokenIndex, depth = position262, tokenIndex262, depth262
					}
				l263:
					depth--
					add(rulePegText, position261)
				}
				if !_rules[ruleAction13]() {
					goto l259
				}
				depth--
				add(ruleDefWindowedFrom, position260)
			}
			return true
		l259:
			position, tokenIndex, depth = position259, tokenIndex259, depth259
			return false
		},
		/* 18 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l275
					}
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if !_rules[ruleTuplesInterval]() {
						goto l272
					}
				}
			l274:
				depth--
				add(ruleInterval, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 19 TimeInterval <- <(NumericLiteral sp SECONDS Action14)> */
		func() bool {
			position276, tokenIndex276, depth276 := position, tokenIndex, depth
			{
				position277 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l276
				}
				if !_rules[rulesp]() {
					goto l276
				}
				if !_rules[ruleSECONDS]() {
					goto l276
				}
				if !_rules[ruleAction14]() {
					goto l276
				}
				depth--
				add(ruleTimeInterval, position277)
			}
			return true
		l276:
			position, tokenIndex, depth = position276, tokenIndex276, depth276
			return false
		},
		/* 20 TuplesInterval <- <(NumericLiteral sp TUPLES Action15)> */
		func() bool {
			position278, tokenIndex278, depth278 := position, tokenIndex, depth
			{
				position279 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l278
				}
				if !_rules[rulesp]() {
					goto l278
				}
				if !_rules[ruleTUPLES]() {
					goto l278
				}
				if !_rules[ruleAction15]() {
					goto l278
				}
				depth--
				add(ruleTuplesInterval, position279)
			}
			return true
		l278:
			position, tokenIndex, depth = position278, tokenIndex278, depth278
			return false
		},
		/* 21 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l280
				}
				if !_rules[rulesp]() {
					goto l280
				}
			l282:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l283
					}
					position++
					if !_rules[rulesp]() {
						goto l283
					}
					if !_rules[ruleRelationLike]() {
						goto l283
					}
					goto l282
				l283:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
				}
				depth--
				add(ruleRelations, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 22 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position284, tokenIndex284, depth284 := position, tokenIndex, depth
			{
				position285 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l284
				}
				if !_rules[rulesp]() {
					goto l284
				}
			l286:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l287
					}
					position++
					if !_rules[rulesp]() {
						goto l287
					}
					if !_rules[ruleDefRelationLike]() {
						goto l287
					}
					goto l286
				l287:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
				}
				depth--
				add(ruleDefRelations, position285)
			}
			return true
		l284:
			position, tokenIndex, depth = position284, tokenIndex284, depth284
			return false
		},
		/* 23 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action16)> */
		func() bool {
			position288, tokenIndex288, depth288 := position, tokenIndex, depth
			{
				position289 := position
				depth++
				{
					position290 := position
					depth++
					{
						position291, tokenIndex291, depth291 := position, tokenIndex, depth
						{
							position293, tokenIndex293, depth293 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l294
							}
							position++
							goto l293
						l294:
							position, tokenIndex, depth = position293, tokenIndex293, depth293
							if buffer[position] != rune('W') {
								goto l291
							}
							position++
						}
					l293:
						{
							position295, tokenIndex295, depth295 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l296
							}
							position++
							goto l295
						l296:
							position, tokenIndex, depth = position295, tokenIndex295, depth295
							if buffer[position] != rune('H') {
								goto l291
							}
							position++
						}
					l295:
						{
							position297, tokenIndex297, depth297 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l298
							}
							position++
							goto l297
						l298:
							position, tokenIndex, depth = position297, tokenIndex297, depth297
							if buffer[position] != rune('E') {
								goto l291
							}
							position++
						}
					l297:
						{
							position299, tokenIndex299, depth299 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l300
							}
							position++
							goto l299
						l300:
							position, tokenIndex, depth = position299, tokenIndex299, depth299
							if buffer[position] != rune('R') {
								goto l291
							}
							position++
						}
					l299:
						{
							position301, tokenIndex301, depth301 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l302
							}
							position++
							goto l301
						l302:
							position, tokenIndex, depth = position301, tokenIndex301, depth301
							if buffer[position] != rune('E') {
								goto l291
							}
							position++
						}
					l301:
						if !_rules[rulesp]() {
							goto l291
						}
						if !_rules[ruleExpression]() {
							goto l291
						}
						goto l292
					l291:
						position, tokenIndex, depth = position291, tokenIndex291, depth291
					}
				l292:
					depth--
					add(rulePegText, position290)
				}
				if !_rules[ruleAction16]() {
					goto l288
				}
				depth--
				add(ruleFilter, position289)
			}
			return true
		l288:
			position, tokenIndex, depth = position288, tokenIndex288, depth288
			return false
		},
		/* 24 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action17)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				{
					position305 := position
					depth++
					{
						position306, tokenIndex306, depth306 := position, tokenIndex, depth
						{
							position308, tokenIndex308, depth308 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l309
							}
							position++
							goto l308
						l309:
							position, tokenIndex, depth = position308, tokenIndex308, depth308
							if buffer[position] != rune('G') {
								goto l306
							}
							position++
						}
					l308:
						{
							position310, tokenIndex310, depth310 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l311
							}
							position++
							goto l310
						l311:
							position, tokenIndex, depth = position310, tokenIndex310, depth310
							if buffer[position] != rune('R') {
								goto l306
							}
							position++
						}
					l310:
						{
							position312, tokenIndex312, depth312 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l313
							}
							position++
							goto l312
						l313:
							position, tokenIndex, depth = position312, tokenIndex312, depth312
							if buffer[position] != rune('O') {
								goto l306
							}
							position++
						}
					l312:
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('U') {
								goto l306
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('P') {
								goto l306
							}
							position++
						}
					l316:
						if !_rules[rulesp]() {
							goto l306
						}
						{
							position318, tokenIndex318, depth318 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l319
							}
							position++
							goto l318
						l319:
							position, tokenIndex, depth = position318, tokenIndex318, depth318
							if buffer[position] != rune('B') {
								goto l306
							}
							position++
						}
					l318:
						{
							position320, tokenIndex320, depth320 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l321
							}
							position++
							goto l320
						l321:
							position, tokenIndex, depth = position320, tokenIndex320, depth320
							if buffer[position] != rune('Y') {
								goto l306
							}
							position++
						}
					l320:
						if !_rules[rulesp]() {
							goto l306
						}
						if !_rules[ruleGroupList]() {
							goto l306
						}
						goto l307
					l306:
						position, tokenIndex, depth = position306, tokenIndex306, depth306
					}
				l307:
					depth--
					add(rulePegText, position305)
				}
				if !_rules[ruleAction17]() {
					goto l303
				}
				depth--
				add(ruleGrouping, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 25 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l322
				}
				if !_rules[rulesp]() {
					goto l322
				}
			l324:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l325
					}
					position++
					if !_rules[rulesp]() {
						goto l325
					}
					if !_rules[ruleExpression]() {
						goto l325
					}
					goto l324
				l325:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
				}
				depth--
				add(ruleGroupList, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 26 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action18)> */
		func() bool {
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				{
					position328 := position
					depth++
					{
						position329, tokenIndex329, depth329 := position, tokenIndex, depth
						{
							position331, tokenIndex331, depth331 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l332
							}
							position++
							goto l331
						l332:
							position, tokenIndex, depth = position331, tokenIndex331, depth331
							if buffer[position] != rune('H') {
								goto l329
							}
							position++
						}
					l331:
						{
							position333, tokenIndex333, depth333 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l334
							}
							position++
							goto l333
						l334:
							position, tokenIndex, depth = position333, tokenIndex333, depth333
							if buffer[position] != rune('A') {
								goto l329
							}
							position++
						}
					l333:
						{
							position335, tokenIndex335, depth335 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l336
							}
							position++
							goto l335
						l336:
							position, tokenIndex, depth = position335, tokenIndex335, depth335
							if buffer[position] != rune('V') {
								goto l329
							}
							position++
						}
					l335:
						{
							position337, tokenIndex337, depth337 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l338
							}
							position++
							goto l337
						l338:
							position, tokenIndex, depth = position337, tokenIndex337, depth337
							if buffer[position] != rune('I') {
								goto l329
							}
							position++
						}
					l337:
						{
							position339, tokenIndex339, depth339 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l340
							}
							position++
							goto l339
						l340:
							position, tokenIndex, depth = position339, tokenIndex339, depth339
							if buffer[position] != rune('N') {
								goto l329
							}
							position++
						}
					l339:
						{
							position341, tokenIndex341, depth341 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l342
							}
							position++
							goto l341
						l342:
							position, tokenIndex, depth = position341, tokenIndex341, depth341
							if buffer[position] != rune('G') {
								goto l329
							}
							position++
						}
					l341:
						if !_rules[rulesp]() {
							goto l329
						}
						if !_rules[ruleExpression]() {
							goto l329
						}
						goto l330
					l329:
						position, tokenIndex, depth = position329, tokenIndex329, depth329
					}
				l330:
					depth--
					add(rulePegText, position328)
				}
				if !_rules[ruleAction18]() {
					goto l326
				}
				depth--
				add(ruleHaving, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 27 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action19))> */
		func() bool {
			position343, tokenIndex343, depth343 := position, tokenIndex, depth
			{
				position344 := position
				depth++
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l346
					}
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if !_rules[ruleStreamWindow]() {
						goto l343
					}
					if !_rules[ruleAction19]() {
						goto l343
					}
				}
			l345:
				depth--
				add(ruleRelationLike, position344)
			}
			return true
		l343:
			position, tokenIndex, depth = position343, tokenIndex343, depth343
			return false
		},
		/* 28 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action20))> */
		func() bool {
			position347, tokenIndex347, depth347 := position, tokenIndex, depth
			{
				position348 := position
				depth++
				{
					position349, tokenIndex349, depth349 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l350
					}
					goto l349
				l350:
					position, tokenIndex, depth = position349, tokenIndex349, depth349
					if !_rules[ruleDefStreamWindow]() {
						goto l347
					}
					if !_rules[ruleAction20]() {
						goto l347
					}
				}
			l349:
				depth--
				add(ruleDefRelationLike, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 29 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action21)> */
		func() bool {
			position351, tokenIndex351, depth351 := position, tokenIndex, depth
			{
				position352 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l351
				}
				if !_rules[rulesp]() {
					goto l351
				}
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l354
					}
					position++
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if buffer[position] != rune('A') {
						goto l351
					}
					position++
				}
			l353:
				{
					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l356
					}
					position++
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune('S') {
						goto l351
					}
					position++
				}
			l355:
				if !_rules[rulesp]() {
					goto l351
				}
				if !_rules[ruleIdentifier]() {
					goto l351
				}
				if !_rules[ruleAction21]() {
					goto l351
				}
				depth--
				add(ruleAliasedStreamWindow, position352)
			}
			return true
		l351:
			position, tokenIndex, depth = position351, tokenIndex351, depth351
			return false
		},
		/* 30 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action22)> */
		func() bool {
			position357, tokenIndex357, depth357 := position, tokenIndex, depth
			{
				position358 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l357
				}
				if !_rules[rulesp]() {
					goto l357
				}
				{
					position359, tokenIndex359, depth359 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l360
					}
					position++
					goto l359
				l360:
					position, tokenIndex, depth = position359, tokenIndex359, depth359
					if buffer[position] != rune('A') {
						goto l357
					}
					position++
				}
			l359:
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l362
					}
					position++
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if buffer[position] != rune('S') {
						goto l357
					}
					position++
				}
			l361:
				if !_rules[rulesp]() {
					goto l357
				}
				if !_rules[ruleIdentifier]() {
					goto l357
				}
				if !_rules[ruleAction22]() {
					goto l357
				}
				depth--
				add(ruleDefAliasedStreamWindow, position358)
			}
			return true
		l357:
			position, tokenIndex, depth = position357, tokenIndex357, depth357
			return false
		},
		/* 31 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action23)> */
		func() bool {
			position363, tokenIndex363, depth363 := position, tokenIndex, depth
			{
				position364 := position
				depth++
				if !_rules[ruleStream]() {
					goto l363
				}
				if !_rules[rulesp]() {
					goto l363
				}
				if buffer[position] != rune('[') {
					goto l363
				}
				position++
				if !_rules[rulesp]() {
					goto l363
				}
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('R') {
						goto l363
					}
					position++
				}
			l365:
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					if buffer[position] != rune('A') {
						goto l363
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('N') {
						goto l363
					}
					position++
				}
			l369:
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l372
					}
					position++
					goto l371
				l372:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
					if buffer[position] != rune('G') {
						goto l363
					}
					position++
				}
			l371:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('E') {
						goto l363
					}
					position++
				}
			l373:
				if !_rules[rulesp]() {
					goto l363
				}
				if !_rules[ruleInterval]() {
					goto l363
				}
				if !_rules[rulesp]() {
					goto l363
				}
				if buffer[position] != rune(']') {
					goto l363
				}
				position++
				if !_rules[ruleAction23]() {
					goto l363
				}
				depth--
				add(ruleStreamWindow, position364)
			}
			return true
		l363:
			position, tokenIndex, depth = position363, tokenIndex363, depth363
			return false
		},
		/* 32 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action24)> */
		func() bool {
			position375, tokenIndex375, depth375 := position, tokenIndex, depth
			{
				position376 := position
				depth++
				if !_rules[ruleStream]() {
					goto l375
				}
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l377
					}
					if buffer[position] != rune('[') {
						goto l377
					}
					position++
					if !_rules[rulesp]() {
						goto l377
					}
					{
						position379, tokenIndex379, depth379 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l380
						}
						position++
						goto l379
					l380:
						position, tokenIndex, depth = position379, tokenIndex379, depth379
						if buffer[position] != rune('R') {
							goto l377
						}
						position++
					}
				l379:
					{
						position381, tokenIndex381, depth381 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l382
						}
						position++
						goto l381
					l382:
						position, tokenIndex, depth = position381, tokenIndex381, depth381
						if buffer[position] != rune('A') {
							goto l377
						}
						position++
					}
				l381:
					{
						position383, tokenIndex383, depth383 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l384
						}
						position++
						goto l383
					l384:
						position, tokenIndex, depth = position383, tokenIndex383, depth383
						if buffer[position] != rune('N') {
							goto l377
						}
						position++
					}
				l383:
					{
						position385, tokenIndex385, depth385 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l386
						}
						position++
						goto l385
					l386:
						position, tokenIndex, depth = position385, tokenIndex385, depth385
						if buffer[position] != rune('G') {
							goto l377
						}
						position++
					}
				l385:
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l388
						}
						position++
						goto l387
					l388:
						position, tokenIndex, depth = position387, tokenIndex387, depth387
						if buffer[position] != rune('E') {
							goto l377
						}
						position++
					}
				l387:
					if !_rules[rulesp]() {
						goto l377
					}
					if !_rules[ruleInterval]() {
						goto l377
					}
					if !_rules[rulesp]() {
						goto l377
					}
					if buffer[position] != rune(']') {
						goto l377
					}
					position++
					goto l378
				l377:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
				}
			l378:
				if !_rules[ruleAction24]() {
					goto l375
				}
				depth--
				add(ruleDefStreamWindow, position376)
			}
			return true
		l375:
			position, tokenIndex, depth = position375, tokenIndex375, depth375
			return false
		},
		/* 33 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action25)> */
		func() bool {
			position389, tokenIndex389, depth389 := position, tokenIndex, depth
			{
				position390 := position
				depth++
				{
					position391 := position
					depth++
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						{
							position394, tokenIndex394, depth394 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l395
							}
							position++
							goto l394
						l395:
							position, tokenIndex, depth = position394, tokenIndex394, depth394
							if buffer[position] != rune('W') {
								goto l392
							}
							position++
						}
					l394:
						{
							position396, tokenIndex396, depth396 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l397
							}
							position++
							goto l396
						l397:
							position, tokenIndex, depth = position396, tokenIndex396, depth396
							if buffer[position] != rune('I') {
								goto l392
							}
							position++
						}
					l396:
						{
							position398, tokenIndex398, depth398 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l399
							}
							position++
							goto l398
						l399:
							position, tokenIndex, depth = position398, tokenIndex398, depth398
							if buffer[position] != rune('T') {
								goto l392
							}
							position++
						}
					l398:
						{
							position400, tokenIndex400, depth400 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l401
							}
							position++
							goto l400
						l401:
							position, tokenIndex, depth = position400, tokenIndex400, depth400
							if buffer[position] != rune('H') {
								goto l392
							}
							position++
						}
					l400:
						if !_rules[rulesp]() {
							goto l392
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l392
						}
						if !_rules[rulesp]() {
							goto l392
						}
					l402:
						{
							position403, tokenIndex403, depth403 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l403
							}
							position++
							if !_rules[rulesp]() {
								goto l403
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l403
							}
							goto l402
						l403:
							position, tokenIndex, depth = position403, tokenIndex403, depth403
						}
						goto l393
					l392:
						position, tokenIndex, depth = position392, tokenIndex392, depth392
					}
				l393:
					depth--
					add(rulePegText, position391)
				}
				if !_rules[ruleAction25]() {
					goto l389
				}
				depth--
				add(ruleSourceSinkSpecs, position390)
			}
			return true
		l389:
			position, tokenIndex, depth = position389, tokenIndex389, depth389
			return false
		},
		/* 34 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action26)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l404
				}
				if buffer[position] != rune('=') {
					goto l404
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l404
				}
				if !_rules[ruleAction26]() {
					goto l404
				}
				depth--
				add(ruleSourceSinkParam, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 35 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position406, tokenIndex406, depth406 := position, tokenIndex, depth
			{
				position407 := position
				depth++
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l409
					}
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if !_rules[ruleLiteral]() {
						goto l406
					}
				}
			l408:
				depth--
				add(ruleSourceSinkParamVal, position407)
			}
			return true
		l406:
			position, tokenIndex, depth = position406, tokenIndex406, depth406
			return false
		},
		/* 36 Expression <- <orExpr> */
		func() bool {
			position410, tokenIndex410, depth410 := position, tokenIndex, depth
			{
				position411 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l410
				}
				depth--
				add(ruleExpression, position411)
			}
			return true
		l410:
			position, tokenIndex, depth = position410, tokenIndex410, depth410
			return false
		},
		/* 37 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action27)> */
		func() bool {
			position412, tokenIndex412, depth412 := position, tokenIndex, depth
			{
				position413 := position
				depth++
				{
					position414 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l412
					}
					if !_rules[rulesp]() {
						goto l412
					}
					{
						position415, tokenIndex415, depth415 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l415
						}
						if !_rules[rulesp]() {
							goto l415
						}
						if !_rules[ruleandExpr]() {
							goto l415
						}
						goto l416
					l415:
						position, tokenIndex, depth = position415, tokenIndex415, depth415
					}
				l416:
					depth--
					add(rulePegText, position414)
				}
				if !_rules[ruleAction27]() {
					goto l412
				}
				depth--
				add(ruleorExpr, position413)
			}
			return true
		l412:
			position, tokenIndex, depth = position412, tokenIndex412, depth412
			return false
		},
		/* 38 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action28)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				{
					position419 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l417
					}
					if !_rules[rulesp]() {
						goto l417
					}
					{
						position420, tokenIndex420, depth420 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l420
						}
						if !_rules[rulesp]() {
							goto l420
						}
						if !_rules[rulecomparisonExpr]() {
							goto l420
						}
						goto l421
					l420:
						position, tokenIndex, depth = position420, tokenIndex420, depth420
					}
				l421:
					depth--
					add(rulePegText, position419)
				}
				if !_rules[ruleAction28]() {
					goto l417
				}
				depth--
				add(ruleandExpr, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 39 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action29)> */
		func() bool {
			position422, tokenIndex422, depth422 := position, tokenIndex, depth
			{
				position423 := position
				depth++
				{
					position424 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l422
					}
					if !_rules[rulesp]() {
						goto l422
					}
					{
						position425, tokenIndex425, depth425 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l425
						}
						if !_rules[rulesp]() {
							goto l425
						}
						if !_rules[ruletermExpr]() {
							goto l425
						}
						goto l426
					l425:
						position, tokenIndex, depth = position425, tokenIndex425, depth425
					}
				l426:
					depth--
					add(rulePegText, position424)
				}
				if !_rules[ruleAction29]() {
					goto l422
				}
				depth--
				add(rulecomparisonExpr, position423)
			}
			return true
		l422:
			position, tokenIndex, depth = position422, tokenIndex422, depth422
			return false
		},
		/* 40 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action30)> */
		func() bool {
			position427, tokenIndex427, depth427 := position, tokenIndex, depth
			{
				position428 := position
				depth++
				{
					position429 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l427
					}
					if !_rules[rulesp]() {
						goto l427
					}
					{
						position430, tokenIndex430, depth430 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l430
						}
						if !_rules[rulesp]() {
							goto l430
						}
						if !_rules[ruleproductExpr]() {
							goto l430
						}
						goto l431
					l430:
						position, tokenIndex, depth = position430, tokenIndex430, depth430
					}
				l431:
					depth--
					add(rulePegText, position429)
				}
				if !_rules[ruleAction30]() {
					goto l427
				}
				depth--
				add(ruletermExpr, position428)
			}
			return true
		l427:
			position, tokenIndex, depth = position427, tokenIndex427, depth427
			return false
		},
		/* 41 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action31)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				{
					position434 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l432
					}
					if !_rules[rulesp]() {
						goto l432
					}
					{
						position435, tokenIndex435, depth435 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l435
						}
						if !_rules[rulesp]() {
							goto l435
						}
						if !_rules[rulebaseExpr]() {
							goto l435
						}
						goto l436
					l435:
						position, tokenIndex, depth = position435, tokenIndex435, depth435
					}
				l436:
					depth--
					add(rulePegText, position434)
				}
				if !_rules[ruleAction31]() {
					goto l432
				}
				depth--
				add(ruleproductExpr, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 42 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position437, tokenIndex437, depth437 := position, tokenIndex, depth
			{
				position438 := position
				depth++
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l440
					}
					position++
					if !_rules[rulesp]() {
						goto l440
					}
					if !_rules[ruleExpression]() {
						goto l440
					}
					if !_rules[rulesp]() {
						goto l440
					}
					if buffer[position] != rune(')') {
						goto l440
					}
					position++
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleBooleanLiteral]() {
						goto l441
					}
					goto l439
				l441:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleFuncApp]() {
						goto l442
					}
					goto l439
				l442:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleRowValue]() {
						goto l443
					}
					goto l439
				l443:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleLiteral]() {
						goto l437
					}
				}
			l439:
				depth--
				add(rulebaseExpr, position438)
			}
			return true
		l437:
			position, tokenIndex, depth = position437, tokenIndex437, depth437
			return false
		},
		/* 43 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action32)> */
		func() bool {
			position444, tokenIndex444, depth444 := position, tokenIndex, depth
			{
				position445 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l444
				}
				if !_rules[rulesp]() {
					goto l444
				}
				if buffer[position] != rune('(') {
					goto l444
				}
				position++
				if !_rules[rulesp]() {
					goto l444
				}
				if !_rules[ruleFuncParams]() {
					goto l444
				}
				if !_rules[rulesp]() {
					goto l444
				}
				if buffer[position] != rune(')') {
					goto l444
				}
				position++
				if !_rules[ruleAction32]() {
					goto l444
				}
				depth--
				add(ruleFuncApp, position445)
			}
			return true
		l444:
			position, tokenIndex, depth = position444, tokenIndex444, depth444
			return false
		},
		/* 44 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action33)> */
		func() bool {
			position446, tokenIndex446, depth446 := position, tokenIndex, depth
			{
				position447 := position
				depth++
				{
					position448 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l446
					}
					if !_rules[rulesp]() {
						goto l446
					}
				l449:
					{
						position450, tokenIndex450, depth450 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l450
						}
						position++
						if !_rules[rulesp]() {
							goto l450
						}
						if !_rules[ruleExpression]() {
							goto l450
						}
						goto l449
					l450:
						position, tokenIndex, depth = position450, tokenIndex450, depth450
					}
					depth--
					add(rulePegText, position448)
				}
				if !_rules[ruleAction33]() {
					goto l446
				}
				depth--
				add(ruleFuncParams, position447)
			}
			return true
		l446:
			position, tokenIndex, depth = position446, tokenIndex446, depth446
			return false
		},
		/* 45 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position451, tokenIndex451, depth451 := position, tokenIndex, depth
			{
				position452 := position
				depth++
				{
					position453, tokenIndex453, depth453 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l454
					}
					goto l453
				l454:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
					if !_rules[ruleNumericLiteral]() {
						goto l455
					}
					goto l453
				l455:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
					if !_rules[ruleStringLiteral]() {
						goto l451
					}
				}
			l453:
				depth--
				add(ruleLiteral, position452)
			}
			return true
		l451:
			position, tokenIndex, depth = position451, tokenIndex451, depth451
			return false
		},
		/* 46 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l459
					}
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !_rules[ruleNotEqual]() {
						goto l460
					}
					goto l458
				l460:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !_rules[ruleLessOrEqual]() {
						goto l461
					}
					goto l458
				l461:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !_rules[ruleLess]() {
						goto l462
					}
					goto l458
				l462:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !_rules[ruleGreaterOrEqual]() {
						goto l463
					}
					goto l458
				l463:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !_rules[ruleGreater]() {
						goto l464
					}
					goto l458
				l464:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !_rules[ruleNotEqual]() {
						goto l456
					}
				}
			l458:
				depth--
				add(ruleComparisonOp, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 47 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position465, tokenIndex465, depth465 := position, tokenIndex, depth
			{
				position466 := position
				depth++
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l468
					}
					goto l467
				l468:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
					if !_rules[ruleMinus]() {
						goto l465
					}
				}
			l467:
				depth--
				add(rulePlusMinusOp, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 48 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l472
					}
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if !_rules[ruleDivide]() {
						goto l473
					}
					goto l471
				l473:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if !_rules[ruleModulo]() {
						goto l469
					}
				}
			l471:
				depth--
				add(ruleMultDivOp, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 49 Stream <- <(<ident> Action34)> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				{
					position476 := position
					depth++
					if !_rules[ruleident]() {
						goto l474
					}
					depth--
					add(rulePegText, position476)
				}
				if !_rules[ruleAction34]() {
					goto l474
				}
				depth--
				add(ruleStream, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 50 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action35)> */
		func() bool {
			position477, tokenIndex477, depth477 := position, tokenIndex, depth
			{
				position478 := position
				depth++
				{
					position479 := position
					depth++
					{
						position480, tokenIndex480, depth480 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l480
						}
						if buffer[position] != rune(':') {
							goto l480
						}
						position++
						goto l481
					l480:
						position, tokenIndex, depth = position480, tokenIndex480, depth480
					}
				l481:
					{
						position482, tokenIndex482, depth482 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l483
						}
						position++
						goto l482
					l483:
						position, tokenIndex, depth = position482, tokenIndex482, depth482
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l477
						}
						position++
					}
				l482:
				l484:
					{
						position485, tokenIndex485, depth485 := position, tokenIndex, depth
						{
							position486, tokenIndex486, depth486 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l487
							}
							position++
							goto l486
						l487:
							position, tokenIndex, depth = position486, tokenIndex486, depth486
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l488
							}
							position++
							goto l486
						l488:
							position, tokenIndex, depth = position486, tokenIndex486, depth486
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l489
							}
							position++
							goto l486
						l489:
							position, tokenIndex, depth = position486, tokenIndex486, depth486
							if buffer[position] != rune('_') {
								goto l490
							}
							position++
							goto l486
						l490:
							position, tokenIndex, depth = position486, tokenIndex486, depth486
							if buffer[position] != rune('.') {
								goto l485
							}
							position++
						}
					l486:
						goto l484
					l485:
						position, tokenIndex, depth = position485, tokenIndex485, depth485
					}
					depth--
					add(rulePegText, position479)
				}
				if !_rules[ruleAction35]() {
					goto l477
				}
				depth--
				add(ruleRowValue, position478)
			}
			return true
		l477:
			position, tokenIndex, depth = position477, tokenIndex477, depth477
			return false
		},
		/* 51 NumericLiteral <- <(<('-'? [0-9]+)> Action36)> */
		func() bool {
			position491, tokenIndex491, depth491 := position, tokenIndex, depth
			{
				position492 := position
				depth++
				{
					position493 := position
					depth++
					{
						position494, tokenIndex494, depth494 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l494
						}
						position++
						goto l495
					l494:
						position, tokenIndex, depth = position494, tokenIndex494, depth494
					}
				l495:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l491
					}
					position++
				l496:
					{
						position497, tokenIndex497, depth497 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l497
						}
						position++
						goto l496
					l497:
						position, tokenIndex, depth = position497, tokenIndex497, depth497
					}
					depth--
					add(rulePegText, position493)
				}
				if !_rules[ruleAction36]() {
					goto l491
				}
				depth--
				add(ruleNumericLiteral, position492)
			}
			return true
		l491:
			position, tokenIndex, depth = position491, tokenIndex491, depth491
			return false
		},
		/* 52 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action37)> */
		func() bool {
			position498, tokenIndex498, depth498 := position, tokenIndex, depth
			{
				position499 := position
				depth++
				{
					position500 := position
					depth++
					{
						position501, tokenIndex501, depth501 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l501
						}
						position++
						goto l502
					l501:
						position, tokenIndex, depth = position501, tokenIndex501, depth501
					}
				l502:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l498
					}
					position++
				l503:
					{
						position504, tokenIndex504, depth504 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l504
						}
						position++
						goto l503
					l504:
						position, tokenIndex, depth = position504, tokenIndex504, depth504
					}
					if buffer[position] != rune('.') {
						goto l498
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l498
					}
					position++
				l505:
					{
						position506, tokenIndex506, depth506 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l506
						}
						position++
						goto l505
					l506:
						position, tokenIndex, depth = position506, tokenIndex506, depth506
					}
					depth--
					add(rulePegText, position500)
				}
				if !_rules[ruleAction37]() {
					goto l498
				}
				depth--
				add(ruleFloatLiteral, position499)
			}
			return true
		l498:
			position, tokenIndex, depth = position498, tokenIndex498, depth498
			return false
		},
		/* 53 Function <- <(<ident> Action38)> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				{
					position509 := position
					depth++
					if !_rules[ruleident]() {
						goto l507
					}
					depth--
					add(rulePegText, position509)
				}
				if !_rules[ruleAction38]() {
					goto l507
				}
				depth--
				add(ruleFunction, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 54 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position510, tokenIndex510, depth510 := position, tokenIndex, depth
			{
				position511 := position
				depth++
				{
					position512, tokenIndex512, depth512 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l513
					}
					goto l512
				l513:
					position, tokenIndex, depth = position512, tokenIndex512, depth512
					if !_rules[ruleFALSE]() {
						goto l510
					}
				}
			l512:
				depth--
				add(ruleBooleanLiteral, position511)
			}
			return true
		l510:
			position, tokenIndex, depth = position510, tokenIndex510, depth510
			return false
		},
		/* 55 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action39)> */
		func() bool {
			position514, tokenIndex514, depth514 := position, tokenIndex, depth
			{
				position515 := position
				depth++
				{
					position516 := position
					depth++
					{
						position517, tokenIndex517, depth517 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l518
						}
						position++
						goto l517
					l518:
						position, tokenIndex, depth = position517, tokenIndex517, depth517
						if buffer[position] != rune('T') {
							goto l514
						}
						position++
					}
				l517:
					{
						position519, tokenIndex519, depth519 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l520
						}
						position++
						goto l519
					l520:
						position, tokenIndex, depth = position519, tokenIndex519, depth519
						if buffer[position] != rune('R') {
							goto l514
						}
						position++
					}
				l519:
					{
						position521, tokenIndex521, depth521 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l522
						}
						position++
						goto l521
					l522:
						position, tokenIndex, depth = position521, tokenIndex521, depth521
						if buffer[position] != rune('U') {
							goto l514
						}
						position++
					}
				l521:
					{
						position523, tokenIndex523, depth523 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l524
						}
						position++
						goto l523
					l524:
						position, tokenIndex, depth = position523, tokenIndex523, depth523
						if buffer[position] != rune('E') {
							goto l514
						}
						position++
					}
				l523:
					depth--
					add(rulePegText, position516)
				}
				if !_rules[ruleAction39]() {
					goto l514
				}
				depth--
				add(ruleTRUE, position515)
			}
			return true
		l514:
			position, tokenIndex, depth = position514, tokenIndex514, depth514
			return false
		},
		/* 56 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action40)> */
		func() bool {
			position525, tokenIndex525, depth525 := position, tokenIndex, depth
			{
				position526 := position
				depth++
				{
					position527 := position
					depth++
					{
						position528, tokenIndex528, depth528 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l529
						}
						position++
						goto l528
					l529:
						position, tokenIndex, depth = position528, tokenIndex528, depth528
						if buffer[position] != rune('F') {
							goto l525
						}
						position++
					}
				l528:
					{
						position530, tokenIndex530, depth530 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l531
						}
						position++
						goto l530
					l531:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
						if buffer[position] != rune('A') {
							goto l525
						}
						position++
					}
				l530:
					{
						position532, tokenIndex532, depth532 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l533
						}
						position++
						goto l532
					l533:
						position, tokenIndex, depth = position532, tokenIndex532, depth532
						if buffer[position] != rune('L') {
							goto l525
						}
						position++
					}
				l532:
					{
						position534, tokenIndex534, depth534 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l535
						}
						position++
						goto l534
					l535:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if buffer[position] != rune('S') {
							goto l525
						}
						position++
					}
				l534:
					{
						position536, tokenIndex536, depth536 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l537
						}
						position++
						goto l536
					l537:
						position, tokenIndex, depth = position536, tokenIndex536, depth536
						if buffer[position] != rune('E') {
							goto l525
						}
						position++
					}
				l536:
					depth--
					add(rulePegText, position527)
				}
				if !_rules[ruleAction40]() {
					goto l525
				}
				depth--
				add(ruleFALSE, position526)
			}
			return true
		l525:
			position, tokenIndex, depth = position525, tokenIndex525, depth525
			return false
		},
		/* 57 Wildcard <- <(<'*'> Action41)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					position540 := position
					depth++
					if buffer[position] != rune('*') {
						goto l538
					}
					position++
					depth--
					add(rulePegText, position540)
				}
				if !_rules[ruleAction41]() {
					goto l538
				}
				depth--
				add(ruleWildcard, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 58 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action42)> */
		func() bool {
			position541, tokenIndex541, depth541 := position, tokenIndex, depth
			{
				position542 := position
				depth++
				{
					position543 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l541
					}
					position++
				l544:
					{
						position545, tokenIndex545, depth545 := position, tokenIndex, depth
						{
							position546, tokenIndex546, depth546 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l547
							}
							position++
							if buffer[position] != rune('\'') {
								goto l547
							}
							position++
							goto l546
						l547:
							position, tokenIndex, depth = position546, tokenIndex546, depth546
							{
								position548, tokenIndex548, depth548 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l548
								}
								position++
								goto l545
							l548:
								position, tokenIndex, depth = position548, tokenIndex548, depth548
							}
							if !matchDot() {
								goto l545
							}
						}
					l546:
						goto l544
					l545:
						position, tokenIndex, depth = position545, tokenIndex545, depth545
					}
					if buffer[position] != rune('\'') {
						goto l541
					}
					position++
					depth--
					add(rulePegText, position543)
				}
				if !_rules[ruleAction42]() {
					goto l541
				}
				depth--
				add(ruleStringLiteral, position542)
			}
			return true
		l541:
			position, tokenIndex, depth = position541, tokenIndex541, depth541
			return false
		},
		/* 59 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action43)> */
		func() bool {
			position549, tokenIndex549, depth549 := position, tokenIndex, depth
			{
				position550 := position
				depth++
				{
					position551 := position
					depth++
					{
						position552, tokenIndex552, depth552 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l553
						}
						position++
						goto l552
					l553:
						position, tokenIndex, depth = position552, tokenIndex552, depth552
						if buffer[position] != rune('I') {
							goto l549
						}
						position++
					}
				l552:
					{
						position554, tokenIndex554, depth554 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l555
						}
						position++
						goto l554
					l555:
						position, tokenIndex, depth = position554, tokenIndex554, depth554
						if buffer[position] != rune('S') {
							goto l549
						}
						position++
					}
				l554:
					{
						position556, tokenIndex556, depth556 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l557
						}
						position++
						goto l556
					l557:
						position, tokenIndex, depth = position556, tokenIndex556, depth556
						if buffer[position] != rune('T') {
							goto l549
						}
						position++
					}
				l556:
					{
						position558, tokenIndex558, depth558 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l559
						}
						position++
						goto l558
					l559:
						position, tokenIndex, depth = position558, tokenIndex558, depth558
						if buffer[position] != rune('R') {
							goto l549
						}
						position++
					}
				l558:
					{
						position560, tokenIndex560, depth560 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l561
						}
						position++
						goto l560
					l561:
						position, tokenIndex, depth = position560, tokenIndex560, depth560
						if buffer[position] != rune('E') {
							goto l549
						}
						position++
					}
				l560:
					{
						position562, tokenIndex562, depth562 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l563
						}
						position++
						goto l562
					l563:
						position, tokenIndex, depth = position562, tokenIndex562, depth562
						if buffer[position] != rune('A') {
							goto l549
						}
						position++
					}
				l562:
					{
						position564, tokenIndex564, depth564 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l565
						}
						position++
						goto l564
					l565:
						position, tokenIndex, depth = position564, tokenIndex564, depth564
						if buffer[position] != rune('M') {
							goto l549
						}
						position++
					}
				l564:
					depth--
					add(rulePegText, position551)
				}
				if !_rules[ruleAction43]() {
					goto l549
				}
				depth--
				add(ruleISTREAM, position550)
			}
			return true
		l549:
			position, tokenIndex, depth = position549, tokenIndex549, depth549
			return false
		},
		/* 60 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action44)> */
		func() bool {
			position566, tokenIndex566, depth566 := position, tokenIndex, depth
			{
				position567 := position
				depth++
				{
					position568 := position
					depth++
					{
						position569, tokenIndex569, depth569 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l570
						}
						position++
						goto l569
					l570:
						position, tokenIndex, depth = position569, tokenIndex569, depth569
						if buffer[position] != rune('D') {
							goto l566
						}
						position++
					}
				l569:
					{
						position571, tokenIndex571, depth571 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l572
						}
						position++
						goto l571
					l572:
						position, tokenIndex, depth = position571, tokenIndex571, depth571
						if buffer[position] != rune('S') {
							goto l566
						}
						position++
					}
				l571:
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l574
						}
						position++
						goto l573
					l574:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
						if buffer[position] != rune('T') {
							goto l566
						}
						position++
					}
				l573:
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l576
						}
						position++
						goto l575
					l576:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
						if buffer[position] != rune('R') {
							goto l566
						}
						position++
					}
				l575:
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('E') {
							goto l566
						}
						position++
					}
				l577:
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l580
						}
						position++
						goto l579
					l580:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
						if buffer[position] != rune('A') {
							goto l566
						}
						position++
					}
				l579:
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l582
						}
						position++
						goto l581
					l582:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
						if buffer[position] != rune('M') {
							goto l566
						}
						position++
					}
				l581:
					depth--
					add(rulePegText, position568)
				}
				if !_rules[ruleAction44]() {
					goto l566
				}
				depth--
				add(ruleDSTREAM, position567)
			}
			return true
		l566:
			position, tokenIndex, depth = position566, tokenIndex566, depth566
			return false
		},
		/* 61 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action45)> */
		func() bool {
			position583, tokenIndex583, depth583 := position, tokenIndex, depth
			{
				position584 := position
				depth++
				{
					position585 := position
					depth++
					{
						position586, tokenIndex586, depth586 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l587
						}
						position++
						goto l586
					l587:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
						if buffer[position] != rune('R') {
							goto l583
						}
						position++
					}
				l586:
					{
						position588, tokenIndex588, depth588 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l589
						}
						position++
						goto l588
					l589:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
						if buffer[position] != rune('S') {
							goto l583
						}
						position++
					}
				l588:
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
						if buffer[position] != rune('T') {
							goto l583
						}
						position++
					}
				l590:
					{
						position592, tokenIndex592, depth592 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex, depth = position592, tokenIndex592, depth592
						if buffer[position] != rune('R') {
							goto l583
						}
						position++
					}
				l592:
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
						if buffer[position] != rune('E') {
							goto l583
						}
						position++
					}
				l594:
					{
						position596, tokenIndex596, depth596 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l597
						}
						position++
						goto l596
					l597:
						position, tokenIndex, depth = position596, tokenIndex596, depth596
						if buffer[position] != rune('A') {
							goto l583
						}
						position++
					}
				l596:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l599
						}
						position++
						goto l598
					l599:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
						if buffer[position] != rune('M') {
							goto l583
						}
						position++
					}
				l598:
					depth--
					add(rulePegText, position585)
				}
				if !_rules[ruleAction45]() {
					goto l583
				}
				depth--
				add(ruleRSTREAM, position584)
			}
			return true
		l583:
			position, tokenIndex, depth = position583, tokenIndex583, depth583
			return false
		},
		/* 62 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action46)> */
		func() bool {
			position600, tokenIndex600, depth600 := position, tokenIndex, depth
			{
				position601 := position
				depth++
				{
					position602 := position
					depth++
					{
						position603, tokenIndex603, depth603 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l604
						}
						position++
						goto l603
					l604:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
						if buffer[position] != rune('T') {
							goto l600
						}
						position++
					}
				l603:
					{
						position605, tokenIndex605, depth605 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l606
						}
						position++
						goto l605
					l606:
						position, tokenIndex, depth = position605, tokenIndex605, depth605
						if buffer[position] != rune('U') {
							goto l600
						}
						position++
					}
				l605:
					{
						position607, tokenIndex607, depth607 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l608
						}
						position++
						goto l607
					l608:
						position, tokenIndex, depth = position607, tokenIndex607, depth607
						if buffer[position] != rune('P') {
							goto l600
						}
						position++
					}
				l607:
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('L') {
							goto l600
						}
						position++
					}
				l609:
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('E') {
							goto l600
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('S') {
							goto l600
						}
						position++
					}
				l613:
					depth--
					add(rulePegText, position602)
				}
				if !_rules[ruleAction46]() {
					goto l600
				}
				depth--
				add(ruleTUPLES, position601)
			}
			return true
		l600:
			position, tokenIndex, depth = position600, tokenIndex600, depth600
			return false
		},
		/* 63 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action47)> */
		func() bool {
			position615, tokenIndex615, depth615 := position, tokenIndex, depth
			{
				position616 := position
				depth++
				{
					position617 := position
					depth++
					{
						position618, tokenIndex618, depth618 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l619
						}
						position++
						goto l618
					l619:
						position, tokenIndex, depth = position618, tokenIndex618, depth618
						if buffer[position] != rune('S') {
							goto l615
						}
						position++
					}
				l618:
					{
						position620, tokenIndex620, depth620 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l621
						}
						position++
						goto l620
					l621:
						position, tokenIndex, depth = position620, tokenIndex620, depth620
						if buffer[position] != rune('E') {
							goto l615
						}
						position++
					}
				l620:
					{
						position622, tokenIndex622, depth622 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l623
						}
						position++
						goto l622
					l623:
						position, tokenIndex, depth = position622, tokenIndex622, depth622
						if buffer[position] != rune('C') {
							goto l615
						}
						position++
					}
				l622:
					{
						position624, tokenIndex624, depth624 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('O') {
							goto l615
						}
						position++
					}
				l624:
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('N') {
							goto l615
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('D') {
							goto l615
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('S') {
							goto l615
						}
						position++
					}
				l630:
					depth--
					add(rulePegText, position617)
				}
				if !_rules[ruleAction47]() {
					goto l615
				}
				depth--
				add(ruleSECONDS, position616)
			}
			return true
		l615:
			position, tokenIndex, depth = position615, tokenIndex615, depth615
			return false
		},
		/* 64 StreamIdentifier <- <(<ident> Action48)> */
		func() bool {
			position632, tokenIndex632, depth632 := position, tokenIndex, depth
			{
				position633 := position
				depth++
				{
					position634 := position
					depth++
					if !_rules[ruleident]() {
						goto l632
					}
					depth--
					add(rulePegText, position634)
				}
				if !_rules[ruleAction48]() {
					goto l632
				}
				depth--
				add(ruleStreamIdentifier, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 65 SourceSinkType <- <(<ident> Action49)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				{
					position637 := position
					depth++
					if !_rules[ruleident]() {
						goto l635
					}
					depth--
					add(rulePegText, position637)
				}
				if !_rules[ruleAction49]() {
					goto l635
				}
				depth--
				add(ruleSourceSinkType, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 66 SourceSinkParamKey <- <(<ident> Action50)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					if !_rules[ruleident]() {
						goto l638
					}
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction50]() {
					goto l638
				}
				depth--
				add(ruleSourceSinkParamKey, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 67 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action51)> */
		func() bool {
			position641, tokenIndex641, depth641 := position, tokenIndex, depth
			{
				position642 := position
				depth++
				{
					position643 := position
					depth++
					{
						position644, tokenIndex644, depth644 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l645
						}
						position++
						goto l644
					l645:
						position, tokenIndex, depth = position644, tokenIndex644, depth644
						if buffer[position] != rune('O') {
							goto l641
						}
						position++
					}
				l644:
					{
						position646, tokenIndex646, depth646 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l647
						}
						position++
						goto l646
					l647:
						position, tokenIndex, depth = position646, tokenIndex646, depth646
						if buffer[position] != rune('R') {
							goto l641
						}
						position++
					}
				l646:
					depth--
					add(rulePegText, position643)
				}
				if !_rules[ruleAction51]() {
					goto l641
				}
				depth--
				add(ruleOr, position642)
			}
			return true
		l641:
			position, tokenIndex, depth = position641, tokenIndex641, depth641
			return false
		},
		/* 68 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action52)> */
		func() bool {
			position648, tokenIndex648, depth648 := position, tokenIndex, depth
			{
				position649 := position
				depth++
				{
					position650 := position
					depth++
					{
						position651, tokenIndex651, depth651 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l652
						}
						position++
						goto l651
					l652:
						position, tokenIndex, depth = position651, tokenIndex651, depth651
						if buffer[position] != rune('A') {
							goto l648
						}
						position++
					}
				l651:
					{
						position653, tokenIndex653, depth653 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l654
						}
						position++
						goto l653
					l654:
						position, tokenIndex, depth = position653, tokenIndex653, depth653
						if buffer[position] != rune('N') {
							goto l648
						}
						position++
					}
				l653:
					{
						position655, tokenIndex655, depth655 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l656
						}
						position++
						goto l655
					l656:
						position, tokenIndex, depth = position655, tokenIndex655, depth655
						if buffer[position] != rune('D') {
							goto l648
						}
						position++
					}
				l655:
					depth--
					add(rulePegText, position650)
				}
				if !_rules[ruleAction52]() {
					goto l648
				}
				depth--
				add(ruleAnd, position649)
			}
			return true
		l648:
			position, tokenIndex, depth = position648, tokenIndex648, depth648
			return false
		},
		/* 69 Equal <- <(<'='> Action53)> */
		func() bool {
			position657, tokenIndex657, depth657 := position, tokenIndex, depth
			{
				position658 := position
				depth++
				{
					position659 := position
					depth++
					if buffer[position] != rune('=') {
						goto l657
					}
					position++
					depth--
					add(rulePegText, position659)
				}
				if !_rules[ruleAction53]() {
					goto l657
				}
				depth--
				add(ruleEqual, position658)
			}
			return true
		l657:
			position, tokenIndex, depth = position657, tokenIndex657, depth657
			return false
		},
		/* 70 Less <- <(<'<'> Action54)> */
		func() bool {
			position660, tokenIndex660, depth660 := position, tokenIndex, depth
			{
				position661 := position
				depth++
				{
					position662 := position
					depth++
					if buffer[position] != rune('<') {
						goto l660
					}
					position++
					depth--
					add(rulePegText, position662)
				}
				if !_rules[ruleAction54]() {
					goto l660
				}
				depth--
				add(ruleLess, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 71 LessOrEqual <- <(<('<' '=')> Action55)> */
		func() bool {
			position663, tokenIndex663, depth663 := position, tokenIndex, depth
			{
				position664 := position
				depth++
				{
					position665 := position
					depth++
					if buffer[position] != rune('<') {
						goto l663
					}
					position++
					if buffer[position] != rune('=') {
						goto l663
					}
					position++
					depth--
					add(rulePegText, position665)
				}
				if !_rules[ruleAction55]() {
					goto l663
				}
				depth--
				add(ruleLessOrEqual, position664)
			}
			return true
		l663:
			position, tokenIndex, depth = position663, tokenIndex663, depth663
			return false
		},
		/* 72 Greater <- <(<'>'> Action56)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				{
					position668 := position
					depth++
					if buffer[position] != rune('>') {
						goto l666
					}
					position++
					depth--
					add(rulePegText, position668)
				}
				if !_rules[ruleAction56]() {
					goto l666
				}
				depth--
				add(ruleGreater, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 73 GreaterOrEqual <- <(<('>' '=')> Action57)> */
		func() bool {
			position669, tokenIndex669, depth669 := position, tokenIndex, depth
			{
				position670 := position
				depth++
				{
					position671 := position
					depth++
					if buffer[position] != rune('>') {
						goto l669
					}
					position++
					if buffer[position] != rune('=') {
						goto l669
					}
					position++
					depth--
					add(rulePegText, position671)
				}
				if !_rules[ruleAction57]() {
					goto l669
				}
				depth--
				add(ruleGreaterOrEqual, position670)
			}
			return true
		l669:
			position, tokenIndex, depth = position669, tokenIndex669, depth669
			return false
		},
		/* 74 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action58)> */
		func() bool {
			position672, tokenIndex672, depth672 := position, tokenIndex, depth
			{
				position673 := position
				depth++
				{
					position674 := position
					depth++
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l676
						}
						position++
						if buffer[position] != rune('=') {
							goto l676
						}
						position++
						goto l675
					l676:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('<') {
							goto l672
						}
						position++
						if buffer[position] != rune('>') {
							goto l672
						}
						position++
					}
				l675:
					depth--
					add(rulePegText, position674)
				}
				if !_rules[ruleAction58]() {
					goto l672
				}
				depth--
				add(ruleNotEqual, position673)
			}
			return true
		l672:
			position, tokenIndex, depth = position672, tokenIndex672, depth672
			return false
		},
		/* 75 Plus <- <(<'+'> Action59)> */
		func() bool {
			position677, tokenIndex677, depth677 := position, tokenIndex, depth
			{
				position678 := position
				depth++
				{
					position679 := position
					depth++
					if buffer[position] != rune('+') {
						goto l677
					}
					position++
					depth--
					add(rulePegText, position679)
				}
				if !_rules[ruleAction59]() {
					goto l677
				}
				depth--
				add(rulePlus, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 76 Minus <- <(<'-'> Action60)> */
		func() bool {
			position680, tokenIndex680, depth680 := position, tokenIndex, depth
			{
				position681 := position
				depth++
				{
					position682 := position
					depth++
					if buffer[position] != rune('-') {
						goto l680
					}
					position++
					depth--
					add(rulePegText, position682)
				}
				if !_rules[ruleAction60]() {
					goto l680
				}
				depth--
				add(ruleMinus, position681)
			}
			return true
		l680:
			position, tokenIndex, depth = position680, tokenIndex680, depth680
			return false
		},
		/* 77 Multiply <- <(<'*'> Action61)> */
		func() bool {
			position683, tokenIndex683, depth683 := position, tokenIndex, depth
			{
				position684 := position
				depth++
				{
					position685 := position
					depth++
					if buffer[position] != rune('*') {
						goto l683
					}
					position++
					depth--
					add(rulePegText, position685)
				}
				if !_rules[ruleAction61]() {
					goto l683
				}
				depth--
				add(ruleMultiply, position684)
			}
			return true
		l683:
			position, tokenIndex, depth = position683, tokenIndex683, depth683
			return false
		},
		/* 78 Divide <- <(<'/'> Action62)> */
		func() bool {
			position686, tokenIndex686, depth686 := position, tokenIndex, depth
			{
				position687 := position
				depth++
				{
					position688 := position
					depth++
					if buffer[position] != rune('/') {
						goto l686
					}
					position++
					depth--
					add(rulePegText, position688)
				}
				if !_rules[ruleAction62]() {
					goto l686
				}
				depth--
				add(ruleDivide, position687)
			}
			return true
		l686:
			position, tokenIndex, depth = position686, tokenIndex686, depth686
			return false
		},
		/* 79 Modulo <- <(<'%'> Action63)> */
		func() bool {
			position689, tokenIndex689, depth689 := position, tokenIndex, depth
			{
				position690 := position
				depth++
				{
					position691 := position
					depth++
					if buffer[position] != rune('%') {
						goto l689
					}
					position++
					depth--
					add(rulePegText, position691)
				}
				if !_rules[ruleAction63]() {
					goto l689
				}
				depth--
				add(ruleModulo, position690)
			}
			return true
		l689:
			position, tokenIndex, depth = position689, tokenIndex689, depth689
			return false
		},
		/* 80 Identifier <- <(<ident> Action64)> */
		func() bool {
			position692, tokenIndex692, depth692 := position, tokenIndex, depth
			{
				position693 := position
				depth++
				{
					position694 := position
					depth++
					if !_rules[ruleident]() {
						goto l692
					}
					depth--
					add(rulePegText, position694)
				}
				if !_rules[ruleAction64]() {
					goto l692
				}
				depth--
				add(ruleIdentifier, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 81 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position695, tokenIndex695, depth695 := position, tokenIndex, depth
			{
				position696 := position
				depth++
				{
					position697, tokenIndex697, depth697 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l698
					}
					position++
					goto l697
				l698:
					position, tokenIndex, depth = position697, tokenIndex697, depth697
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l695
					}
					position++
				}
			l697:
			l699:
				{
					position700, tokenIndex700, depth700 := position, tokenIndex, depth
					{
						position701, tokenIndex701, depth701 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l702
						}
						position++
						goto l701
					l702:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l703
						}
						position++
						goto l701
					l703:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l704
						}
						position++
						goto l701
					l704:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
						if buffer[position] != rune('_') {
							goto l700
						}
						position++
					}
				l701:
					goto l699
				l700:
					position, tokenIndex, depth = position700, tokenIndex700, depth700
				}
				depth--
				add(ruleident, position696)
			}
			return true
		l695:
			position, tokenIndex, depth = position695, tokenIndex695, depth695
			return false
		},
		/* 82 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position706 := position
				depth++
			l707:
				{
					position708, tokenIndex708, depth708 := position, tokenIndex, depth
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l710
						}
						position++
						goto l709
					l710:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if buffer[position] != rune('\t') {
							goto l711
						}
						position++
						goto l709
					l711:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if buffer[position] != rune('\n') {
							goto l708
						}
						position++
					}
				l709:
					goto l707
				l708:
					position, tokenIndex, depth = position708, tokenIndex708, depth708
				}
				depth--
				add(rulesp, position706)
			}
			return true
		},
		/* 84 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 85 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 86 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 87 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 88 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 89 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		nil,
		/* 91 Action6 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 92 Action7 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 93 Action8 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 94 Action9 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 95 Action10 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 96 Action11 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 97 Action12 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 98 Action13 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 99 Action14 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 100 Action15 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 101 Action16 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 102 Action17 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 103 Action18 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 104 Action19 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 105 Action20 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 106 Action21 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 107 Action22 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 108 Action23 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 109 Action24 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 110 Action25 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 111 Action26 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 112 Action27 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 113 Action28 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 114 Action29 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 115 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 116 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 117 Action32 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 118 Action33 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 119 Action34 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 120 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 121 Action36 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 122 Action37 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 123 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 124 Action39 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 125 Action40 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 126 Action41 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 127 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 128 Action43 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 129 Action44 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 130 Action45 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 131 Action46 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 132 Action47 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 133 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 134 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 135 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 136 Action51 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 137 Action52 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 138 Action53 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 139 Action54 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 140 Action55 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 141 Action56 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 142 Action57 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 143 Action58 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 144 Action59 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 145 Action60 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 146 Action61 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 147 Action62 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 148 Action63 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 149 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
	}
	p.rules = _rules
}

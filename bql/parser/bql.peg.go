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
	ruleCreateStreamFromSourceStmt
	ruleCreateStreamFromSourceExtStmt
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
	ruleAction6
	ruleAction7
	rulePegText
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
	ruleAction65
	ruleAction66

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
	"CreateStreamFromSourceStmt",
	"CreateStreamFromSourceExtStmt",
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
	"Action6",
	"Action7",
	"PegText",
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
	"Action65",
	"Action66",

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
	rules  [154]func() bool
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

			p.AssembleCreateStreamFromSource()

		case ruleAction6:

			p.AssembleCreateStreamFromSourceExt()

		case ruleAction7:

			p.AssembleInsertIntoSelect()

		case ruleAction8:

			p.AssembleEmitter(begin, end)

		case ruleAction9:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction10:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction11:

			p.AssembleStreamEmitInterval()

		case ruleAction12:

			p.AssembleProjections(begin, end)

		case ruleAction13:

			p.AssembleAlias()

		case ruleAction14:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction15:

			p.AssembleWindowedFrom(begin, end)

		case ruleAction16:

			p.AssembleInterval()

		case ruleAction17:

			p.AssembleInterval()

		case ruleAction18:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction19:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction20:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction21:

			p.EnsureAliasedStreamWindow()

		case ruleAction22:

			p.EnsureAliasedStreamWindow()

		case ruleAction23:

			p.AssembleAliasedStreamWindow()

		case ruleAction24:

			p.AssembleAliasedStreamWindow()

		case ruleAction25:

			p.AssembleStreamWindow()

		case ruleAction26:

			p.AssembleStreamWindow()

		case ruleAction27:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction28:

			p.AssembleSourceSinkParam()

		case ruleAction29:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction30:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction31:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction32:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction33:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction34:

			p.AssembleFuncApp()

		case ruleAction35:

			p.AssembleExpressions(begin, end)

		case ruleAction36:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction37:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction38:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction41:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction42:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction43:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction44:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction45:

			p.PushComponent(begin, end, Istream)

		case ruleAction46:

			p.PushComponent(begin, end, Dstream)

		case ruleAction47:

			p.PushComponent(begin, end, Rstream)

		case ruleAction48:

			p.PushComponent(begin, end, Tuples)

		case ruleAction49:

			p.PushComponent(begin, end, Seconds)

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction53:

			p.PushComponent(begin, end, Or)

		case ruleAction54:

			p.PushComponent(begin, end, And)

		case ruleAction55:

			p.PushComponent(begin, end, Equal)

		case ruleAction56:

			p.PushComponent(begin, end, Less)

		case ruleAction57:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction58:

			p.PushComponent(begin, end, Greater)

		case ruleAction59:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction60:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction61:

			p.PushComponent(begin, end, Plus)

		case ruleAction62:

			p.PushComponent(begin, end, Minus)

		case ruleAction63:

			p.PushComponent(begin, end, Multiply)

		case ruleAction64:

			p.PushComponent(begin, end, Divide)

		case ruleAction65:

			p.PushComponent(begin, end, Modulo)

		case ruleAction66:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateStreamFromSourceStmt / CreateStreamFromSourceExtStmt / CreateSinkStmt / InsertIntoSelectStmt / CreateStateStmt)> */
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
					if !_rules[ruleCreateStreamFromSourceStmt]() {
						goto l13
					}
					goto l9
				l13:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleCreateStreamFromSourceExtStmt]() {
						goto l14
					}
					goto l9
				l14:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleCreateSinkStmt]() {
						goto l15
					}
					goto l9
				l15:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l16
					}
					goto l9
				l16:
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
			position17, tokenIndex17, depth17 := position, tokenIndex, depth
			{
				position18 := position
				depth++
				{
					position19, tokenIndex19, depth19 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l20
					}
					position++
					goto l19
				l20:
					position, tokenIndex, depth = position19, tokenIndex19, depth19
					if buffer[position] != rune('S') {
						goto l17
					}
					position++
				}
			l19:
				{
					position21, tokenIndex21, depth21 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l22
					}
					position++
					goto l21
				l22:
					position, tokenIndex, depth = position21, tokenIndex21, depth21
					if buffer[position] != rune('E') {
						goto l17
					}
					position++
				}
			l21:
				{
					position23, tokenIndex23, depth23 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l24
					}
					position++
					goto l23
				l24:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if buffer[position] != rune('L') {
						goto l17
					}
					position++
				}
			l23:
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l26
					}
					position++
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if buffer[position] != rune('E') {
						goto l17
					}
					position++
				}
			l25:
				{
					position27, tokenIndex27, depth27 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l28
					}
					position++
					goto l27
				l28:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					if buffer[position] != rune('C') {
						goto l17
					}
					position++
				}
			l27:
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l30
					}
					position++
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if buffer[position] != rune('T') {
						goto l17
					}
					position++
				}
			l29:
				if !_rules[rulesp]() {
					goto l17
				}
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					if !_rules[ruleEmitter]() {
						goto l31
					}
					goto l32
				l31:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
				}
			l32:
				if !_rules[rulesp]() {
					goto l17
				}
				if !_rules[ruleProjections]() {
					goto l17
				}
				if !_rules[rulesp]() {
					goto l17
				}
				if !_rules[ruleDefWindowedFrom]() {
					goto l17
				}
				if !_rules[rulesp]() {
					goto l17
				}
				if !_rules[ruleFilter]() {
					goto l17
				}
				if !_rules[rulesp]() {
					goto l17
				}
				if !_rules[ruleGrouping]() {
					goto l17
				}
				if !_rules[rulesp]() {
					goto l17
				}
				if !_rules[ruleHaving]() {
					goto l17
				}
				if !_rules[rulesp]() {
					goto l17
				}
				if !_rules[ruleAction0]() {
					goto l17
				}
				depth--
				add(ruleSelectStmt, position18)
			}
			return true
		l17:
			position, tokenIndex, depth = position17, tokenIndex17, depth17
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position33, tokenIndex33, depth33 := position, tokenIndex, depth
			{
				position34 := position
				depth++
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l36
					}
					position++
					goto l35
				l36:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if buffer[position] != rune('C') {
						goto l33
					}
					position++
				}
			l35:
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l38
					}
					position++
					goto l37
				l38:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if buffer[position] != rune('R') {
						goto l33
					}
					position++
				}
			l37:
				{
					position39, tokenIndex39, depth39 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l40
					}
					position++
					goto l39
				l40:
					position, tokenIndex, depth = position39, tokenIndex39, depth39
					if buffer[position] != rune('E') {
						goto l33
					}
					position++
				}
			l39:
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l42
					}
					position++
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if buffer[position] != rune('A') {
						goto l33
					}
					position++
				}
			l41:
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l44
					}
					position++
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if buffer[position] != rune('T') {
						goto l33
					}
					position++
				}
			l43:
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l46
					}
					position++
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if buffer[position] != rune('E') {
						goto l33
					}
					position++
				}
			l45:
				if !_rules[rulesp]() {
					goto l33
				}
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l48
					}
					position++
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					if buffer[position] != rune('S') {
						goto l33
					}
					position++
				}
			l47:
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l50
					}
					position++
					goto l49
				l50:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
					if buffer[position] != rune('T') {
						goto l33
					}
					position++
				}
			l49:
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l52
					}
					position++
					goto l51
				l52:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
					if buffer[position] != rune('R') {
						goto l33
					}
					position++
				}
			l51:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					if buffer[position] != rune('E') {
						goto l33
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
					if buffer[position] != rune('A') {
						goto l33
					}
					position++
				}
			l55:
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
					if buffer[position] != rune('M') {
						goto l33
					}
					position++
				}
			l57:
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l33
				}
				if !_rules[rulesp]() {
					goto l33
				}
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l60
					}
					position++
					goto l59
				l60:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
					if buffer[position] != rune('A') {
						goto l33
					}
					position++
				}
			l59:
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
						goto l33
					}
					position++
				}
			l61:
				if !_rules[rulesp]() {
					goto l33
				}
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l64
					}
					position++
					goto l63
				l64:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
					if buffer[position] != rune('S') {
						goto l33
					}
					position++
				}
			l63:
				{
					position65, tokenIndex65, depth65 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l66
					}
					position++
					goto l65
				l66:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
					if buffer[position] != rune('E') {
						goto l33
					}
					position++
				}
			l65:
				{
					position67, tokenIndex67, depth67 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l68
					}
					position++
					goto l67
				l68:
					position, tokenIndex, depth = position67, tokenIndex67, depth67
					if buffer[position] != rune('L') {
						goto l33
					}
					position++
				}
			l67:
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l70
					}
					position++
					goto l69
				l70:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
					if buffer[position] != rune('E') {
						goto l33
					}
					position++
				}
			l69:
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l72
					}
					position++
					goto l71
				l72:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					if buffer[position] != rune('C') {
						goto l33
					}
					position++
				}
			l71:
				{
					position73, tokenIndex73, depth73 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l74
					}
					position++
					goto l73
				l74:
					position, tokenIndex, depth = position73, tokenIndex73, depth73
					if buffer[position] != rune('T') {
						goto l33
					}
					position++
				}
			l73:
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleEmitter]() {
					goto l33
				}
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleProjections]() {
					goto l33
				}
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleWindowedFrom]() {
					goto l33
				}
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleFilter]() {
					goto l33
				}
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleGrouping]() {
					goto l33
				}
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleHaving]() {
					goto l33
				}
				if !_rules[rulesp]() {
					goto l33
				}
				if !_rules[ruleAction1]() {
					goto l33
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position34)
			}
			return true
		l33:
			position, tokenIndex, depth = position33, tokenIndex33, depth33
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position75, tokenIndex75, depth75 := position, tokenIndex, depth
			{
				position76 := position
				depth++
				{
					position77, tokenIndex77, depth77 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l78
					}
					position++
					goto l77
				l78:
					position, tokenIndex, depth = position77, tokenIndex77, depth77
					if buffer[position] != rune('C') {
						goto l75
					}
					position++
				}
			l77:
				{
					position79, tokenIndex79, depth79 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l80
					}
					position++
					goto l79
				l80:
					position, tokenIndex, depth = position79, tokenIndex79, depth79
					if buffer[position] != rune('R') {
						goto l75
					}
					position++
				}
			l79:
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l82
					}
					position++
					goto l81
				l82:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if buffer[position] != rune('E') {
						goto l75
					}
					position++
				}
			l81:
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l84
					}
					position++
					goto l83
				l84:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
					if buffer[position] != rune('A') {
						goto l75
					}
					position++
				}
			l83:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l86
					}
					position++
					goto l85
				l86:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
					if buffer[position] != rune('T') {
						goto l75
					}
					position++
				}
			l85:
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l88
					}
					position++
					goto l87
				l88:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
					if buffer[position] != rune('E') {
						goto l75
					}
					position++
				}
			l87:
				if !_rules[rulesp]() {
					goto l75
				}
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l90
					}
					position++
					goto l89
				l90:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
					if buffer[position] != rune('S') {
						goto l75
					}
					position++
				}
			l89:
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l92
					}
					position++
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					if buffer[position] != rune('O') {
						goto l75
					}
					position++
				}
			l91:
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('U') {
						goto l75
					}
					position++
				}
			l93:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if buffer[position] != rune('R') {
						goto l75
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('C') {
						goto l75
					}
					position++
				}
			l97:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('E') {
						goto l75
					}
					position++
				}
			l99:
				if !_rules[rulesp]() {
					goto l75
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l75
				}
				if !_rules[rulesp]() {
					goto l75
				}
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('T') {
						goto l75
					}
					position++
				}
			l101:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if buffer[position] != rune('Y') {
						goto l75
					}
					position++
				}
			l103:
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('P') {
						goto l75
					}
					position++
				}
			l105:
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l108
					}
					position++
					goto l107
				l108:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
					if buffer[position] != rune('E') {
						goto l75
					}
					position++
				}
			l107:
				if !_rules[rulesp]() {
					goto l75
				}
				if !_rules[ruleSourceSinkType]() {
					goto l75
				}
				if !_rules[rulesp]() {
					goto l75
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l75
				}
				if !_rules[ruleAction2]() {
					goto l75
				}
				depth--
				add(ruleCreateSourceStmt, position76)
			}
			return true
		l75:
			position, tokenIndex, depth = position75, tokenIndex75, depth75
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position109, tokenIndex109, depth109 := position, tokenIndex, depth
			{
				position110 := position
				depth++
				{
					position111, tokenIndex111, depth111 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l112
					}
					position++
					goto l111
				l112:
					position, tokenIndex, depth = position111, tokenIndex111, depth111
					if buffer[position] != rune('C') {
						goto l109
					}
					position++
				}
			l111:
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l114
					}
					position++
					goto l113
				l114:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
					if buffer[position] != rune('R') {
						goto l109
					}
					position++
				}
			l113:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('E') {
						goto l109
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('A') {
						goto l109
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('T') {
						goto l109
					}
					position++
				}
			l119:
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if buffer[position] != rune('E') {
						goto l109
					}
					position++
				}
			l121:
				if !_rules[rulesp]() {
					goto l109
				}
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('S') {
						goto l109
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('I') {
						goto l109
					}
					position++
				}
			l125:
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('N') {
						goto l109
					}
					position++
				}
			l127:
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('K') {
						goto l109
					}
					position++
				}
			l129:
				if !_rules[rulesp]() {
					goto l109
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l109
				}
				if !_rules[rulesp]() {
					goto l109
				}
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('T') {
						goto l109
					}
					position++
				}
			l131:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('Y') {
						goto l109
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('P') {
						goto l109
					}
					position++
				}
			l135:
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l138
					}
					position++
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if buffer[position] != rune('E') {
						goto l109
					}
					position++
				}
			l137:
				if !_rules[rulesp]() {
					goto l109
				}
				if !_rules[ruleSourceSinkType]() {
					goto l109
				}
				if !_rules[rulesp]() {
					goto l109
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l109
				}
				if !_rules[ruleAction3]() {
					goto l109
				}
				depth--
				add(ruleCreateSinkStmt, position110)
			}
			return true
		l109:
			position, tokenIndex, depth = position109, tokenIndex109, depth109
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position139, tokenIndex139, depth139 := position, tokenIndex, depth
			{
				position140 := position
				depth++
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l142
					}
					position++
					goto l141
				l142:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					if buffer[position] != rune('C') {
						goto l139
					}
					position++
				}
			l141:
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					if buffer[position] != rune('R') {
						goto l139
					}
					position++
				}
			l143:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('E') {
						goto l139
					}
					position++
				}
			l145:
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('A') {
						goto l139
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('T') {
						goto l139
					}
					position++
				}
			l149:
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('E') {
						goto l139
					}
					position++
				}
			l151:
				if !_rules[rulesp]() {
					goto l139
				}
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('S') {
						goto l139
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('T') {
						goto l139
					}
					position++
				}
			l155:
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('A') {
						goto l139
					}
					position++
				}
			l157:
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l160
					}
					position++
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if buffer[position] != rune('T') {
						goto l139
					}
					position++
				}
			l159:
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l162
					}
					position++
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if buffer[position] != rune('E') {
						goto l139
					}
					position++
				}
			l161:
				if !_rules[rulesp]() {
					goto l139
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l139
				}
				if !_rules[rulesp]() {
					goto l139
				}
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l164
					}
					position++
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if buffer[position] != rune('T') {
						goto l139
					}
					position++
				}
			l163:
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('Y') {
						goto l139
					}
					position++
				}
			l165:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('P') {
						goto l139
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('E') {
						goto l139
					}
					position++
				}
			l169:
				if !_rules[rulesp]() {
					goto l139
				}
				if !_rules[ruleSourceSinkType]() {
					goto l139
				}
				if !_rules[rulesp]() {
					goto l139
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l139
				}
				if !_rules[ruleAction4]() {
					goto l139
				}
				depth--
				add(ruleCreateStateStmt, position140)
			}
			return true
		l139:
			position, tokenIndex, depth = position139, tokenIndex139, depth139
			return false
		},
		/* 7 CreateStreamFromSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action5)> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('C') {
						goto l171
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('R') {
						goto l171
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
						goto l171
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('A') {
						goto l171
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
						goto l171
					}
					position++
				}
			l181:
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('E') {
						goto l171
					}
					position++
				}
			l183:
				if !_rules[rulesp]() {
					goto l171
				}
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('S') {
						goto l171
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
						goto l171
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('R') {
						goto l171
					}
					position++
				}
			l189:
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('E') {
						goto l171
					}
					position++
				}
			l191:
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('A') {
						goto l171
					}
					position++
				}
			l193:
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('M') {
						goto l171
					}
					position++
				}
			l195:
				if !_rules[rulesp]() {
					goto l171
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l171
				}
				if !_rules[rulesp]() {
					goto l171
				}
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('F') {
						goto l171
					}
					position++
				}
			l197:
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('R') {
						goto l171
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('O') {
						goto l171
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('M') {
						goto l171
					}
					position++
				}
			l203:
				if !_rules[rulesp]() {
					goto l171
				}
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('S') {
						goto l171
					}
					position++
				}
			l205:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('O') {
						goto l171
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('U') {
						goto l171
					}
					position++
				}
			l209:
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('R') {
						goto l171
					}
					position++
				}
			l211:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('C') {
						goto l171
					}
					position++
				}
			l213:
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('E') {
						goto l171
					}
					position++
				}
			l215:
				if !_rules[rulesp]() {
					goto l171
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l171
				}
				if !_rules[ruleAction5]() {
					goto l171
				}
				depth--
				add(ruleCreateStreamFromSourceStmt, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 8 CreateStreamFromSourceExtStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp SourceSinkType sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp SourceSinkSpecs Action6)> */
		func() bool {
			position217, tokenIndex217, depth217 := position, tokenIndex, depth
			{
				position218 := position
				depth++
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('C') {
						goto l217
					}
					position++
				}
			l219:
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if buffer[position] != rune('R') {
						goto l217
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('E') {
						goto l217
					}
					position++
				}
			l223:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('A') {
						goto l217
					}
					position++
				}
			l225:
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l228
					}
					position++
					goto l227
				l228:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
					if buffer[position] != rune('T') {
						goto l217
					}
					position++
				}
			l227:
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('E') {
						goto l217
					}
					position++
				}
			l229:
				if !_rules[rulesp]() {
					goto l217
				}
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if buffer[position] != rune('S') {
						goto l217
					}
					position++
				}
			l231:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('T') {
						goto l217
					}
					position++
				}
			l233:
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l236
					}
					position++
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if buffer[position] != rune('R') {
						goto l217
					}
					position++
				}
			l235:
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if buffer[position] != rune('E') {
						goto l217
					}
					position++
				}
			l237:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					if buffer[position] != rune('A') {
						goto l217
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('M') {
						goto l217
					}
					position++
				}
			l241:
				if !_rules[rulesp]() {
					goto l217
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l217
				}
				if !_rules[rulesp]() {
					goto l217
				}
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('F') {
						goto l217
					}
					position++
				}
			l243:
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l246
					}
					position++
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if buffer[position] != rune('R') {
						goto l217
					}
					position++
				}
			l245:
				{
					position247, tokenIndex247, depth247 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l248
					}
					position++
					goto l247
				l248:
					position, tokenIndex, depth = position247, tokenIndex247, depth247
					if buffer[position] != rune('O') {
						goto l217
					}
					position++
				}
			l247:
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('M') {
						goto l217
					}
					position++
				}
			l249:
				if !_rules[rulesp]() {
					goto l217
				}
				if !_rules[ruleSourceSinkType]() {
					goto l217
				}
				if !_rules[rulesp]() {
					goto l217
				}
				{
					position251, tokenIndex251, depth251 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l252
					}
					position++
					goto l251
				l252:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('S') {
						goto l217
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('O') {
						goto l217
					}
					position++
				}
			l253:
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l256
					}
					position++
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('U') {
						goto l217
					}
					position++
				}
			l255:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('R') {
						goto l217
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('C') {
						goto l217
					}
					position++
				}
			l259:
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('E') {
						goto l217
					}
					position++
				}
			l261:
				if !_rules[rulesp]() {
					goto l217
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l217
				}
				if !_rules[ruleAction6]() {
					goto l217
				}
				depth--
				add(ruleCreateStreamFromSourceExtStmt, position218)
			}
			return true
		l217:
			position, tokenIndex, depth = position217, tokenIndex217, depth217
			return false
		},
		/* 9 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action7)> */
		func() bool {
			position263, tokenIndex263, depth263 := position, tokenIndex, depth
			{
				position264 := position
				depth++
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('I') {
						goto l263
					}
					position++
				}
			l265:
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l268
					}
					position++
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if buffer[position] != rune('N') {
						goto l263
					}
					position++
				}
			l267:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('S') {
						goto l263
					}
					position++
				}
			l269:
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('E') {
						goto l263
					}
					position++
				}
			l271:
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('R') {
						goto l263
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('T') {
						goto l263
					}
					position++
				}
			l275:
				if !_rules[rulesp]() {
					goto l263
				}
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('I') {
						goto l263
					}
					position++
				}
			l277:
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('N') {
						goto l263
					}
					position++
				}
			l279:
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('T') {
						goto l263
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('O') {
						goto l263
					}
					position++
				}
			l283:
				if !_rules[rulesp]() {
					goto l263
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l263
				}
				if !_rules[rulesp]() {
					goto l263
				}
				if !_rules[ruleSelectStmt]() {
					goto l263
				}
				if !_rules[ruleAction7]() {
					goto l263
				}
				depth--
				add(ruleInsertIntoSelectStmt, position264)
			}
			return true
		l263:
			position, tokenIndex, depth = position263, tokenIndex263, depth263
			return false
		},
		/* 10 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action8)> */
		func() bool {
			position285, tokenIndex285, depth285 := position, tokenIndex, depth
			{
				position286 := position
				depth++
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l288
					}
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if !_rules[ruleDSTREAM]() {
						goto l289
					}
					goto l287
				l289:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if !_rules[ruleRSTREAM]() {
						goto l285
					}
				}
			l287:
				{
					position290 := position
					depth++
					{
						position291, tokenIndex291, depth291 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l291
						}
						if buffer[position] != rune('[') {
							goto l291
						}
						position++
						if !_rules[rulesp]() {
							goto l291
						}
						{
							position293, tokenIndex293, depth293 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l294
							}
							position++
							goto l293
						l294:
							position, tokenIndex, depth = position293, tokenIndex293, depth293
							if buffer[position] != rune('E') {
								goto l291
							}
							position++
						}
					l293:
						{
							position295, tokenIndex295, depth295 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l296
							}
							position++
							goto l295
						l296:
							position, tokenIndex, depth = position295, tokenIndex295, depth295
							if buffer[position] != rune('V') {
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
							if buffer[position] != rune('y') {
								goto l302
							}
							position++
							goto l301
						l302:
							position, tokenIndex, depth = position301, tokenIndex301, depth301
							if buffer[position] != rune('Y') {
								goto l291
							}
							position++
						}
					l301:
						if !_rules[rulesp]() {
							goto l291
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l291
						}
						if !_rules[rulesp]() {
							goto l291
						}
						if buffer[position] != rune(']') {
							goto l291
						}
						position++
						goto l292
					l291:
						position, tokenIndex, depth = position291, tokenIndex291, depth291
					}
				l292:
					depth--
					add(rulePegText, position290)
				}
				if !_rules[ruleAction8]() {
					goto l285
				}
				depth--
				add(ruleEmitter, position286)
			}
			return true
		l285:
			position, tokenIndex, depth = position285, tokenIndex285, depth285
			return false
		},
		/* 11 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l306
					}
				l307:
					{
						position308, tokenIndex308, depth308 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l308
						}
						if buffer[position] != rune(',') {
							goto l308
						}
						position++
						if !_rules[rulesp]() {
							goto l308
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l308
						}
						goto l307
					l308:
						position, tokenIndex, depth = position308, tokenIndex308, depth308
					}
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if !_rules[ruleTimeEmitterInterval]() {
						goto l309
					}
					goto l305
				l309:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if !_rules[ruleTupleEmitterInterval]() {
						goto l303
					}
				}
			l305:
				depth--
				add(ruleEmitterIntervals, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 12 TimeEmitterInterval <- <(<TimeInterval> Action9)> */
		func() bool {
			position310, tokenIndex310, depth310 := position, tokenIndex, depth
			{
				position311 := position
				depth++
				{
					position312 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l310
					}
					depth--
					add(rulePegText, position312)
				}
				if !_rules[ruleAction9]() {
					goto l310
				}
				depth--
				add(ruleTimeEmitterInterval, position311)
			}
			return true
		l310:
			position, tokenIndex, depth = position310, tokenIndex310, depth310
			return false
		},
		/* 13 TupleEmitterInterval <- <(<TuplesInterval> Action10)> */
		func() bool {
			position313, tokenIndex313, depth313 := position, tokenIndex, depth
			{
				position314 := position
				depth++
				{
					position315 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l313
					}
					depth--
					add(rulePegText, position315)
				}
				if !_rules[ruleAction10]() {
					goto l313
				}
				depth--
				add(ruleTupleEmitterInterval, position314)
			}
			return true
		l313:
			position, tokenIndex, depth = position313, tokenIndex313, depth313
			return false
		},
		/* 14 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action11)> */
		func() bool {
			position316, tokenIndex316, depth316 := position, tokenIndex, depth
			{
				position317 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l316
				}
				if !_rules[rulesp]() {
					goto l316
				}
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
					if buffer[position] != rune('I') {
						goto l316
					}
					position++
				}
			l318:
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('N') {
						goto l316
					}
					position++
				}
			l320:
				if !_rules[rulesp]() {
					goto l316
				}
				if !_rules[ruleStream]() {
					goto l316
				}
				if !_rules[ruleAction11]() {
					goto l316
				}
				depth--
				add(ruleTupleEmitterFromInterval, position317)
			}
			return true
		l316:
			position, tokenIndex, depth = position316, tokenIndex316, depth316
			return false
		},
		/* 15 Projections <- <(<(Projection sp (',' sp Projection)*)> Action12)> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				{
					position324 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l322
					}
					if !_rules[rulesp]() {
						goto l322
					}
				l325:
					{
						position326, tokenIndex326, depth326 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l326
						}
						position++
						if !_rules[rulesp]() {
							goto l326
						}
						if !_rules[ruleProjection]() {
							goto l326
						}
						goto l325
					l326:
						position, tokenIndex, depth = position326, tokenIndex326, depth326
					}
					depth--
					add(rulePegText, position324)
				}
				if !_rules[ruleAction12]() {
					goto l322
				}
				depth--
				add(ruleProjections, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 16 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position327, tokenIndex327, depth327 := position, tokenIndex, depth
			{
				position328 := position
				depth++
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l330
					}
					goto l329
				l330:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
					if !_rules[ruleExpression]() {
						goto l331
					}
					goto l329
				l331:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
					if !_rules[ruleWildcard]() {
						goto l327
					}
				}
			l329:
				depth--
				add(ruleProjection, position328)
			}
			return true
		l327:
			position, tokenIndex, depth = position327, tokenIndex327, depth327
			return false
		},
		/* 17 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action13)> */
		func() bool {
			position332, tokenIndex332, depth332 := position, tokenIndex, depth
			{
				position333 := position
				depth++
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l335
					}
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if !_rules[ruleWildcard]() {
						goto l332
					}
				}
			l334:
				if !_rules[rulesp]() {
					goto l332
				}
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if buffer[position] != rune('A') {
						goto l332
					}
					position++
				}
			l336:
				{
					position338, tokenIndex338, depth338 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l339
					}
					position++
					goto l338
				l339:
					position, tokenIndex, depth = position338, tokenIndex338, depth338
					if buffer[position] != rune('S') {
						goto l332
					}
					position++
				}
			l338:
				if !_rules[rulesp]() {
					goto l332
				}
				if !_rules[ruleIdentifier]() {
					goto l332
				}
				if !_rules[ruleAction13]() {
					goto l332
				}
				depth--
				add(ruleAliasExpression, position333)
			}
			return true
		l332:
			position, tokenIndex, depth = position332, tokenIndex332, depth332
			return false
		},
		/* 18 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action14)> */
		func() bool {
			position340, tokenIndex340, depth340 := position, tokenIndex, depth
			{
				position341 := position
				depth++
				{
					position342 := position
					depth++
					{
						position343, tokenIndex343, depth343 := position, tokenIndex, depth
						{
							position345, tokenIndex345, depth345 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l346
							}
							position++
							goto l345
						l346:
							position, tokenIndex, depth = position345, tokenIndex345, depth345
							if buffer[position] != rune('F') {
								goto l343
							}
							position++
						}
					l345:
						{
							position347, tokenIndex347, depth347 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l348
							}
							position++
							goto l347
						l348:
							position, tokenIndex, depth = position347, tokenIndex347, depth347
							if buffer[position] != rune('R') {
								goto l343
							}
							position++
						}
					l347:
						{
							position349, tokenIndex349, depth349 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l350
							}
							position++
							goto l349
						l350:
							position, tokenIndex, depth = position349, tokenIndex349, depth349
							if buffer[position] != rune('O') {
								goto l343
							}
							position++
						}
					l349:
						{
							position351, tokenIndex351, depth351 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l352
							}
							position++
							goto l351
						l352:
							position, tokenIndex, depth = position351, tokenIndex351, depth351
							if buffer[position] != rune('M') {
								goto l343
							}
							position++
						}
					l351:
						if !_rules[rulesp]() {
							goto l343
						}
						if !_rules[ruleRelations]() {
							goto l343
						}
						if !_rules[rulesp]() {
							goto l343
						}
						goto l344
					l343:
						position, tokenIndex, depth = position343, tokenIndex343, depth343
					}
				l344:
					depth--
					add(rulePegText, position342)
				}
				if !_rules[ruleAction14]() {
					goto l340
				}
				depth--
				add(ruleWindowedFrom, position341)
			}
			return true
		l340:
			position, tokenIndex, depth = position340, tokenIndex340, depth340
			return false
		},
		/* 19 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action15)> */
		func() bool {
			position353, tokenIndex353, depth353 := position, tokenIndex, depth
			{
				position354 := position
				depth++
				{
					position355 := position
					depth++
					{
						position356, tokenIndex356, depth356 := position, tokenIndex, depth
						{
							position358, tokenIndex358, depth358 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex, depth = position358, tokenIndex358, depth358
							if buffer[position] != rune('F') {
								goto l356
							}
							position++
						}
					l358:
						{
							position360, tokenIndex360, depth360 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l361
							}
							position++
							goto l360
						l361:
							position, tokenIndex, depth = position360, tokenIndex360, depth360
							if buffer[position] != rune('R') {
								goto l356
							}
							position++
						}
					l360:
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('O') {
								goto l356
							}
							position++
						}
					l362:
						{
							position364, tokenIndex364, depth364 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l365
							}
							position++
							goto l364
						l365:
							position, tokenIndex, depth = position364, tokenIndex364, depth364
							if buffer[position] != rune('M') {
								goto l356
							}
							position++
						}
					l364:
						if !_rules[rulesp]() {
							goto l356
						}
						if !_rules[ruleDefRelations]() {
							goto l356
						}
						if !_rules[rulesp]() {
							goto l356
						}
						goto l357
					l356:
						position, tokenIndex, depth = position356, tokenIndex356, depth356
					}
				l357:
					depth--
					add(rulePegText, position355)
				}
				if !_rules[ruleAction15]() {
					goto l353
				}
				depth--
				add(ruleDefWindowedFrom, position354)
			}
			return true
		l353:
			position, tokenIndex, depth = position353, tokenIndex353, depth353
			return false
		},
		/* 20 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position366, tokenIndex366, depth366 := position, tokenIndex, depth
			{
				position367 := position
				depth++
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l369
					}
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if !_rules[ruleTuplesInterval]() {
						goto l366
					}
				}
			l368:
				depth--
				add(ruleInterval, position367)
			}
			return true
		l366:
			position, tokenIndex, depth = position366, tokenIndex366, depth366
			return false
		},
		/* 21 TimeInterval <- <(NumericLiteral sp SECONDS Action16)> */
		func() bool {
			position370, tokenIndex370, depth370 := position, tokenIndex, depth
			{
				position371 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l370
				}
				if !_rules[rulesp]() {
					goto l370
				}
				if !_rules[ruleSECONDS]() {
					goto l370
				}
				if !_rules[ruleAction16]() {
					goto l370
				}
				depth--
				add(ruleTimeInterval, position371)
			}
			return true
		l370:
			position, tokenIndex, depth = position370, tokenIndex370, depth370
			return false
		},
		/* 22 TuplesInterval <- <(NumericLiteral sp TUPLES Action17)> */
		func() bool {
			position372, tokenIndex372, depth372 := position, tokenIndex, depth
			{
				position373 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l372
				}
				if !_rules[rulesp]() {
					goto l372
				}
				if !_rules[ruleTUPLES]() {
					goto l372
				}
				if !_rules[ruleAction17]() {
					goto l372
				}
				depth--
				add(ruleTuplesInterval, position373)
			}
			return true
		l372:
			position, tokenIndex, depth = position372, tokenIndex372, depth372
			return false
		},
		/* 23 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l374
				}
				if !_rules[rulesp]() {
					goto l374
				}
			l376:
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l377
					}
					position++
					if !_rules[rulesp]() {
						goto l377
					}
					if !_rules[ruleRelationLike]() {
						goto l377
					}
					goto l376
				l377:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
				}
				depth--
				add(ruleRelations, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
			return false
		},
		/* 24 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l378
				}
				if !_rules[rulesp]() {
					goto l378
				}
			l380:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l381
					}
					position++
					if !_rules[rulesp]() {
						goto l381
					}
					if !_rules[ruleDefRelationLike]() {
						goto l381
					}
					goto l380
				l381:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
				}
				depth--
				add(ruleDefRelations, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 25 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action18)> */
		func() bool {
			position382, tokenIndex382, depth382 := position, tokenIndex, depth
			{
				position383 := position
				depth++
				{
					position384 := position
					depth++
					{
						position385, tokenIndex385, depth385 := position, tokenIndex, depth
						{
							position387, tokenIndex387, depth387 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l388
							}
							position++
							goto l387
						l388:
							position, tokenIndex, depth = position387, tokenIndex387, depth387
							if buffer[position] != rune('W') {
								goto l385
							}
							position++
						}
					l387:
						{
							position389, tokenIndex389, depth389 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l390
							}
							position++
							goto l389
						l390:
							position, tokenIndex, depth = position389, tokenIndex389, depth389
							if buffer[position] != rune('H') {
								goto l385
							}
							position++
						}
					l389:
						{
							position391, tokenIndex391, depth391 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l392
							}
							position++
							goto l391
						l392:
							position, tokenIndex, depth = position391, tokenIndex391, depth391
							if buffer[position] != rune('E') {
								goto l385
							}
							position++
						}
					l391:
						{
							position393, tokenIndex393, depth393 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l394
							}
							position++
							goto l393
						l394:
							position, tokenIndex, depth = position393, tokenIndex393, depth393
							if buffer[position] != rune('R') {
								goto l385
							}
							position++
						}
					l393:
						{
							position395, tokenIndex395, depth395 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l396
							}
							position++
							goto l395
						l396:
							position, tokenIndex, depth = position395, tokenIndex395, depth395
							if buffer[position] != rune('E') {
								goto l385
							}
							position++
						}
					l395:
						if !_rules[rulesp]() {
							goto l385
						}
						if !_rules[ruleExpression]() {
							goto l385
						}
						goto l386
					l385:
						position, tokenIndex, depth = position385, tokenIndex385, depth385
					}
				l386:
					depth--
					add(rulePegText, position384)
				}
				if !_rules[ruleAction18]() {
					goto l382
				}
				depth--
				add(ruleFilter, position383)
			}
			return true
		l382:
			position, tokenIndex, depth = position382, tokenIndex382, depth382
			return false
		},
		/* 26 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action19)> */
		func() bool {
			position397, tokenIndex397, depth397 := position, tokenIndex, depth
			{
				position398 := position
				depth++
				{
					position399 := position
					depth++
					{
						position400, tokenIndex400, depth400 := position, tokenIndex, depth
						{
							position402, tokenIndex402, depth402 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l403
							}
							position++
							goto l402
						l403:
							position, tokenIndex, depth = position402, tokenIndex402, depth402
							if buffer[position] != rune('G') {
								goto l400
							}
							position++
						}
					l402:
						{
							position404, tokenIndex404, depth404 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l405
							}
							position++
							goto l404
						l405:
							position, tokenIndex, depth = position404, tokenIndex404, depth404
							if buffer[position] != rune('R') {
								goto l400
							}
							position++
						}
					l404:
						{
							position406, tokenIndex406, depth406 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l407
							}
							position++
							goto l406
						l407:
							position, tokenIndex, depth = position406, tokenIndex406, depth406
							if buffer[position] != rune('O') {
								goto l400
							}
							position++
						}
					l406:
						{
							position408, tokenIndex408, depth408 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l409
							}
							position++
							goto l408
						l409:
							position, tokenIndex, depth = position408, tokenIndex408, depth408
							if buffer[position] != rune('U') {
								goto l400
							}
							position++
						}
					l408:
						{
							position410, tokenIndex410, depth410 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l411
							}
							position++
							goto l410
						l411:
							position, tokenIndex, depth = position410, tokenIndex410, depth410
							if buffer[position] != rune('P') {
								goto l400
							}
							position++
						}
					l410:
						if !_rules[rulesp]() {
							goto l400
						}
						{
							position412, tokenIndex412, depth412 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l413
							}
							position++
							goto l412
						l413:
							position, tokenIndex, depth = position412, tokenIndex412, depth412
							if buffer[position] != rune('B') {
								goto l400
							}
							position++
						}
					l412:
						{
							position414, tokenIndex414, depth414 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l415
							}
							position++
							goto l414
						l415:
							position, tokenIndex, depth = position414, tokenIndex414, depth414
							if buffer[position] != rune('Y') {
								goto l400
							}
							position++
						}
					l414:
						if !_rules[rulesp]() {
							goto l400
						}
						if !_rules[ruleGroupList]() {
							goto l400
						}
						goto l401
					l400:
						position, tokenIndex, depth = position400, tokenIndex400, depth400
					}
				l401:
					depth--
					add(rulePegText, position399)
				}
				if !_rules[ruleAction19]() {
					goto l397
				}
				depth--
				add(ruleGrouping, position398)
			}
			return true
		l397:
			position, tokenIndex, depth = position397, tokenIndex397, depth397
			return false
		},
		/* 27 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position416, tokenIndex416, depth416 := position, tokenIndex, depth
			{
				position417 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l416
				}
				if !_rules[rulesp]() {
					goto l416
				}
			l418:
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l419
					}
					position++
					if !_rules[rulesp]() {
						goto l419
					}
					if !_rules[ruleExpression]() {
						goto l419
					}
					goto l418
				l419:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
				}
				depth--
				add(ruleGroupList, position417)
			}
			return true
		l416:
			position, tokenIndex, depth = position416, tokenIndex416, depth416
			return false
		},
		/* 28 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action20)> */
		func() bool {
			position420, tokenIndex420, depth420 := position, tokenIndex, depth
			{
				position421 := position
				depth++
				{
					position422 := position
					depth++
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						{
							position425, tokenIndex425, depth425 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l426
							}
							position++
							goto l425
						l426:
							position, tokenIndex, depth = position425, tokenIndex425, depth425
							if buffer[position] != rune('H') {
								goto l423
							}
							position++
						}
					l425:
						{
							position427, tokenIndex427, depth427 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l428
							}
							position++
							goto l427
						l428:
							position, tokenIndex, depth = position427, tokenIndex427, depth427
							if buffer[position] != rune('A') {
								goto l423
							}
							position++
						}
					l427:
						{
							position429, tokenIndex429, depth429 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l430
							}
							position++
							goto l429
						l430:
							position, tokenIndex, depth = position429, tokenIndex429, depth429
							if buffer[position] != rune('V') {
								goto l423
							}
							position++
						}
					l429:
						{
							position431, tokenIndex431, depth431 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l432
							}
							position++
							goto l431
						l432:
							position, tokenIndex, depth = position431, tokenIndex431, depth431
							if buffer[position] != rune('I') {
								goto l423
							}
							position++
						}
					l431:
						{
							position433, tokenIndex433, depth433 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l434
							}
							position++
							goto l433
						l434:
							position, tokenIndex, depth = position433, tokenIndex433, depth433
							if buffer[position] != rune('N') {
								goto l423
							}
							position++
						}
					l433:
						{
							position435, tokenIndex435, depth435 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l436
							}
							position++
							goto l435
						l436:
							position, tokenIndex, depth = position435, tokenIndex435, depth435
							if buffer[position] != rune('G') {
								goto l423
							}
							position++
						}
					l435:
						if !_rules[rulesp]() {
							goto l423
						}
						if !_rules[ruleExpression]() {
							goto l423
						}
						goto l424
					l423:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
					}
				l424:
					depth--
					add(rulePegText, position422)
				}
				if !_rules[ruleAction20]() {
					goto l420
				}
				depth--
				add(ruleHaving, position421)
			}
			return true
		l420:
			position, tokenIndex, depth = position420, tokenIndex420, depth420
			return false
		},
		/* 29 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action21))> */
		func() bool {
			position437, tokenIndex437, depth437 := position, tokenIndex, depth
			{
				position438 := position
				depth++
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l440
					}
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleStreamWindow]() {
						goto l437
					}
					if !_rules[ruleAction21]() {
						goto l437
					}
				}
			l439:
				depth--
				add(ruleRelationLike, position438)
			}
			return true
		l437:
			position, tokenIndex, depth = position437, tokenIndex437, depth437
			return false
		},
		/* 30 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action22))> */
		func() bool {
			position441, tokenIndex441, depth441 := position, tokenIndex, depth
			{
				position442 := position
				depth++
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l444
					}
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if !_rules[ruleDefStreamWindow]() {
						goto l441
					}
					if !_rules[ruleAction22]() {
						goto l441
					}
				}
			l443:
				depth--
				add(ruleDefRelationLike, position442)
			}
			return true
		l441:
			position, tokenIndex, depth = position441, tokenIndex441, depth441
			return false
		},
		/* 31 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action23)> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l445
				}
				if !_rules[rulesp]() {
					goto l445
				}
				{
					position447, tokenIndex447, depth447 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l448
					}
					position++
					goto l447
				l448:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if buffer[position] != rune('A') {
						goto l445
					}
					position++
				}
			l447:
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l450
					}
					position++
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					if buffer[position] != rune('S') {
						goto l445
					}
					position++
				}
			l449:
				if !_rules[rulesp]() {
					goto l445
				}
				if !_rules[ruleIdentifier]() {
					goto l445
				}
				if !_rules[ruleAction23]() {
					goto l445
				}
				depth--
				add(ruleAliasedStreamWindow, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 32 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action24)> */
		func() bool {
			position451, tokenIndex451, depth451 := position, tokenIndex, depth
			{
				position452 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l451
				}
				if !_rules[rulesp]() {
					goto l451
				}
				{
					position453, tokenIndex453, depth453 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l454
					}
					position++
					goto l453
				l454:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
					if buffer[position] != rune('A') {
						goto l451
					}
					position++
				}
			l453:
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l456
					}
					position++
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if buffer[position] != rune('S') {
						goto l451
					}
					position++
				}
			l455:
				if !_rules[rulesp]() {
					goto l451
				}
				if !_rules[ruleIdentifier]() {
					goto l451
				}
				if !_rules[ruleAction24]() {
					goto l451
				}
				depth--
				add(ruleDefAliasedStreamWindow, position452)
			}
			return true
		l451:
			position, tokenIndex, depth = position451, tokenIndex451, depth451
			return false
		},
		/* 33 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action25)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				if !_rules[ruleStream]() {
					goto l457
				}
				if !_rules[rulesp]() {
					goto l457
				}
				if buffer[position] != rune('[') {
					goto l457
				}
				position++
				if !_rules[rulesp]() {
					goto l457
				}
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l460
					}
					position++
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if buffer[position] != rune('R') {
						goto l457
					}
					position++
				}
			l459:
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l462
					}
					position++
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('A') {
						goto l457
					}
					position++
				}
			l461:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('N') {
						goto l457
					}
					position++
				}
			l463:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l466
					}
					position++
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					if buffer[position] != rune('G') {
						goto l457
					}
					position++
				}
			l465:
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l468
					}
					position++
					goto l467
				l468:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
					if buffer[position] != rune('E') {
						goto l457
					}
					position++
				}
			l467:
				if !_rules[rulesp]() {
					goto l457
				}
				if !_rules[ruleInterval]() {
					goto l457
				}
				if !_rules[rulesp]() {
					goto l457
				}
				if buffer[position] != rune(']') {
					goto l457
				}
				position++
				if !_rules[ruleAction25]() {
					goto l457
				}
				depth--
				add(ruleStreamWindow, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 34 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action26)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				if !_rules[ruleStream]() {
					goto l469
				}
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l471
					}
					if buffer[position] != rune('[') {
						goto l471
					}
					position++
					if !_rules[rulesp]() {
						goto l471
					}
					{
						position473, tokenIndex473, depth473 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l474
						}
						position++
						goto l473
					l474:
						position, tokenIndex, depth = position473, tokenIndex473, depth473
						if buffer[position] != rune('R') {
							goto l471
						}
						position++
					}
				l473:
					{
						position475, tokenIndex475, depth475 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l476
						}
						position++
						goto l475
					l476:
						position, tokenIndex, depth = position475, tokenIndex475, depth475
						if buffer[position] != rune('A') {
							goto l471
						}
						position++
					}
				l475:
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l478
						}
						position++
						goto l477
					l478:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
						if buffer[position] != rune('N') {
							goto l471
						}
						position++
					}
				l477:
					{
						position479, tokenIndex479, depth479 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l480
						}
						position++
						goto l479
					l480:
						position, tokenIndex, depth = position479, tokenIndex479, depth479
						if buffer[position] != rune('G') {
							goto l471
						}
						position++
					}
				l479:
					{
						position481, tokenIndex481, depth481 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l482
						}
						position++
						goto l481
					l482:
						position, tokenIndex, depth = position481, tokenIndex481, depth481
						if buffer[position] != rune('E') {
							goto l471
						}
						position++
					}
				l481:
					if !_rules[rulesp]() {
						goto l471
					}
					if !_rules[ruleInterval]() {
						goto l471
					}
					if !_rules[rulesp]() {
						goto l471
					}
					if buffer[position] != rune(']') {
						goto l471
					}
					position++
					goto l472
				l471:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
				}
			l472:
				if !_rules[ruleAction26]() {
					goto l469
				}
				depth--
				add(ruleDefStreamWindow, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 35 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action27)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				{
					position485 := position
					depth++
					{
						position486, tokenIndex486, depth486 := position, tokenIndex, depth
						{
							position488, tokenIndex488, depth488 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l489
							}
							position++
							goto l488
						l489:
							position, tokenIndex, depth = position488, tokenIndex488, depth488
							if buffer[position] != rune('W') {
								goto l486
							}
							position++
						}
					l488:
						{
							position490, tokenIndex490, depth490 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l491
							}
							position++
							goto l490
						l491:
							position, tokenIndex, depth = position490, tokenIndex490, depth490
							if buffer[position] != rune('I') {
								goto l486
							}
							position++
						}
					l490:
						{
							position492, tokenIndex492, depth492 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l493
							}
							position++
							goto l492
						l493:
							position, tokenIndex, depth = position492, tokenIndex492, depth492
							if buffer[position] != rune('T') {
								goto l486
							}
							position++
						}
					l492:
						{
							position494, tokenIndex494, depth494 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l495
							}
							position++
							goto l494
						l495:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							if buffer[position] != rune('H') {
								goto l486
							}
							position++
						}
					l494:
						if !_rules[rulesp]() {
							goto l486
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l486
						}
						if !_rules[rulesp]() {
							goto l486
						}
					l496:
						{
							position497, tokenIndex497, depth497 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l497
							}
							position++
							if !_rules[rulesp]() {
								goto l497
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l497
							}
							goto l496
						l497:
							position, tokenIndex, depth = position497, tokenIndex497, depth497
						}
						goto l487
					l486:
						position, tokenIndex, depth = position486, tokenIndex486, depth486
					}
				l487:
					depth--
					add(rulePegText, position485)
				}
				if !_rules[ruleAction27]() {
					goto l483
				}
				depth--
				add(ruleSourceSinkSpecs, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 36 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action28)> */
		func() bool {
			position498, tokenIndex498, depth498 := position, tokenIndex, depth
			{
				position499 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l498
				}
				if buffer[position] != rune('=') {
					goto l498
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l498
				}
				if !_rules[ruleAction28]() {
					goto l498
				}
				depth--
				add(ruleSourceSinkParam, position499)
			}
			return true
		l498:
			position, tokenIndex, depth = position498, tokenIndex498, depth498
			return false
		},
		/* 37 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				{
					position502, tokenIndex502, depth502 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l503
					}
					goto l502
				l503:
					position, tokenIndex, depth = position502, tokenIndex502, depth502
					if !_rules[ruleLiteral]() {
						goto l500
					}
				}
			l502:
				depth--
				add(ruleSourceSinkParamVal, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 38 Expression <- <orExpr> */
		func() bool {
			position504, tokenIndex504, depth504 := position, tokenIndex, depth
			{
				position505 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l504
				}
				depth--
				add(ruleExpression, position505)
			}
			return true
		l504:
			position, tokenIndex, depth = position504, tokenIndex504, depth504
			return false
		},
		/* 39 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action29)> */
		func() bool {
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				{
					position508 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l506
					}
					if !_rules[rulesp]() {
						goto l506
					}
					{
						position509, tokenIndex509, depth509 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l509
						}
						if !_rules[rulesp]() {
							goto l509
						}
						if !_rules[ruleandExpr]() {
							goto l509
						}
						goto l510
					l509:
						position, tokenIndex, depth = position509, tokenIndex509, depth509
					}
				l510:
					depth--
					add(rulePegText, position508)
				}
				if !_rules[ruleAction29]() {
					goto l506
				}
				depth--
				add(ruleorExpr, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 40 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action30)> */
		func() bool {
			position511, tokenIndex511, depth511 := position, tokenIndex, depth
			{
				position512 := position
				depth++
				{
					position513 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l511
					}
					if !_rules[rulesp]() {
						goto l511
					}
					{
						position514, tokenIndex514, depth514 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l514
						}
						if !_rules[rulesp]() {
							goto l514
						}
						if !_rules[rulecomparisonExpr]() {
							goto l514
						}
						goto l515
					l514:
						position, tokenIndex, depth = position514, tokenIndex514, depth514
					}
				l515:
					depth--
					add(rulePegText, position513)
				}
				if !_rules[ruleAction30]() {
					goto l511
				}
				depth--
				add(ruleandExpr, position512)
			}
			return true
		l511:
			position, tokenIndex, depth = position511, tokenIndex511, depth511
			return false
		},
		/* 41 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action31)> */
		func() bool {
			position516, tokenIndex516, depth516 := position, tokenIndex, depth
			{
				position517 := position
				depth++
				{
					position518 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l516
					}
					if !_rules[rulesp]() {
						goto l516
					}
					{
						position519, tokenIndex519, depth519 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l519
						}
						if !_rules[rulesp]() {
							goto l519
						}
						if !_rules[ruletermExpr]() {
							goto l519
						}
						goto l520
					l519:
						position, tokenIndex, depth = position519, tokenIndex519, depth519
					}
				l520:
					depth--
					add(rulePegText, position518)
				}
				if !_rules[ruleAction31]() {
					goto l516
				}
				depth--
				add(rulecomparisonExpr, position517)
			}
			return true
		l516:
			position, tokenIndex, depth = position516, tokenIndex516, depth516
			return false
		},
		/* 42 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action32)> */
		func() bool {
			position521, tokenIndex521, depth521 := position, tokenIndex, depth
			{
				position522 := position
				depth++
				{
					position523 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l521
					}
					if !_rules[rulesp]() {
						goto l521
					}
					{
						position524, tokenIndex524, depth524 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l524
						}
						if !_rules[rulesp]() {
							goto l524
						}
						if !_rules[ruleproductExpr]() {
							goto l524
						}
						goto l525
					l524:
						position, tokenIndex, depth = position524, tokenIndex524, depth524
					}
				l525:
					depth--
					add(rulePegText, position523)
				}
				if !_rules[ruleAction32]() {
					goto l521
				}
				depth--
				add(ruletermExpr, position522)
			}
			return true
		l521:
			position, tokenIndex, depth = position521, tokenIndex521, depth521
			return false
		},
		/* 43 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action33)> */
		func() bool {
			position526, tokenIndex526, depth526 := position, tokenIndex, depth
			{
				position527 := position
				depth++
				{
					position528 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l526
					}
					if !_rules[rulesp]() {
						goto l526
					}
					{
						position529, tokenIndex529, depth529 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l529
						}
						if !_rules[rulesp]() {
							goto l529
						}
						if !_rules[rulebaseExpr]() {
							goto l529
						}
						goto l530
					l529:
						position, tokenIndex, depth = position529, tokenIndex529, depth529
					}
				l530:
					depth--
					add(rulePegText, position528)
				}
				if !_rules[ruleAction33]() {
					goto l526
				}
				depth--
				add(ruleproductExpr, position527)
			}
			return true
		l526:
			position, tokenIndex, depth = position526, tokenIndex526, depth526
			return false
		},
		/* 44 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position531, tokenIndex531, depth531 := position, tokenIndex, depth
			{
				position532 := position
				depth++
				{
					position533, tokenIndex533, depth533 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l534
					}
					position++
					if !_rules[rulesp]() {
						goto l534
					}
					if !_rules[ruleExpression]() {
						goto l534
					}
					if !_rules[rulesp]() {
						goto l534
					}
					if buffer[position] != rune(')') {
						goto l534
					}
					position++
					goto l533
				l534:
					position, tokenIndex, depth = position533, tokenIndex533, depth533
					if !_rules[ruleBooleanLiteral]() {
						goto l535
					}
					goto l533
				l535:
					position, tokenIndex, depth = position533, tokenIndex533, depth533
					if !_rules[ruleFuncApp]() {
						goto l536
					}
					goto l533
				l536:
					position, tokenIndex, depth = position533, tokenIndex533, depth533
					if !_rules[ruleRowValue]() {
						goto l537
					}
					goto l533
				l537:
					position, tokenIndex, depth = position533, tokenIndex533, depth533
					if !_rules[ruleLiteral]() {
						goto l531
					}
				}
			l533:
				depth--
				add(rulebaseExpr, position532)
			}
			return true
		l531:
			position, tokenIndex, depth = position531, tokenIndex531, depth531
			return false
		},
		/* 45 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action34)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l538
				}
				if !_rules[rulesp]() {
					goto l538
				}
				if buffer[position] != rune('(') {
					goto l538
				}
				position++
				if !_rules[rulesp]() {
					goto l538
				}
				if !_rules[ruleFuncParams]() {
					goto l538
				}
				if !_rules[rulesp]() {
					goto l538
				}
				if buffer[position] != rune(')') {
					goto l538
				}
				position++
				if !_rules[ruleAction34]() {
					goto l538
				}
				depth--
				add(ruleFuncApp, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 46 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action35)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l540
					}
					if !_rules[rulesp]() {
						goto l540
					}
				l543:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l544
						}
						position++
						if !_rules[rulesp]() {
							goto l544
						}
						if !_rules[ruleExpression]() {
							goto l544
						}
						goto l543
					l544:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
					}
					depth--
					add(rulePegText, position542)
				}
				if !_rules[ruleAction35]() {
					goto l540
				}
				depth--
				add(ruleFuncParams, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 47 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				{
					position547, tokenIndex547, depth547 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l548
					}
					goto l547
				l548:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
					if !_rules[ruleNumericLiteral]() {
						goto l549
					}
					goto l547
				l549:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
					if !_rules[ruleStringLiteral]() {
						goto l545
					}
				}
			l547:
				depth--
				add(ruleLiteral, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 48 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position550, tokenIndex550, depth550 := position, tokenIndex, depth
			{
				position551 := position
				depth++
				{
					position552, tokenIndex552, depth552 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l553
					}
					goto l552
				l553:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if !_rules[ruleNotEqual]() {
						goto l554
					}
					goto l552
				l554:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if !_rules[ruleLessOrEqual]() {
						goto l555
					}
					goto l552
				l555:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if !_rules[ruleLess]() {
						goto l556
					}
					goto l552
				l556:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if !_rules[ruleGreaterOrEqual]() {
						goto l557
					}
					goto l552
				l557:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if !_rules[ruleGreater]() {
						goto l558
					}
					goto l552
				l558:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if !_rules[ruleNotEqual]() {
						goto l550
					}
				}
			l552:
				depth--
				add(ruleComparisonOp, position551)
			}
			return true
		l550:
			position, tokenIndex, depth = position550, tokenIndex550, depth550
			return false
		},
		/* 49 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position559, tokenIndex559, depth559 := position, tokenIndex, depth
			{
				position560 := position
				depth++
				{
					position561, tokenIndex561, depth561 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l562
					}
					goto l561
				l562:
					position, tokenIndex, depth = position561, tokenIndex561, depth561
					if !_rules[ruleMinus]() {
						goto l559
					}
				}
			l561:
				depth--
				add(rulePlusMinusOp, position560)
			}
			return true
		l559:
			position, tokenIndex, depth = position559, tokenIndex559, depth559
			return false
		},
		/* 50 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position563, tokenIndex563, depth563 := position, tokenIndex, depth
			{
				position564 := position
				depth++
				{
					position565, tokenIndex565, depth565 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l566
					}
					goto l565
				l566:
					position, tokenIndex, depth = position565, tokenIndex565, depth565
					if !_rules[ruleDivide]() {
						goto l567
					}
					goto l565
				l567:
					position, tokenIndex, depth = position565, tokenIndex565, depth565
					if !_rules[ruleModulo]() {
						goto l563
					}
				}
			l565:
				depth--
				add(ruleMultDivOp, position564)
			}
			return true
		l563:
			position, tokenIndex, depth = position563, tokenIndex563, depth563
			return false
		},
		/* 51 Stream <- <(<ident> Action36)> */
		func() bool {
			position568, tokenIndex568, depth568 := position, tokenIndex, depth
			{
				position569 := position
				depth++
				{
					position570 := position
					depth++
					if !_rules[ruleident]() {
						goto l568
					}
					depth--
					add(rulePegText, position570)
				}
				if !_rules[ruleAction36]() {
					goto l568
				}
				depth--
				add(ruleStream, position569)
			}
			return true
		l568:
			position, tokenIndex, depth = position568, tokenIndex568, depth568
			return false
		},
		/* 52 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action37)> */
		func() bool {
			position571, tokenIndex571, depth571 := position, tokenIndex, depth
			{
				position572 := position
				depth++
				{
					position573 := position
					depth++
					{
						position574, tokenIndex574, depth574 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l574
						}
						if buffer[position] != rune(':') {
							goto l574
						}
						position++
						goto l575
					l574:
						position, tokenIndex, depth = position574, tokenIndex574, depth574
					}
				l575:
					{
						position576, tokenIndex576, depth576 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l577
						}
						position++
						goto l576
					l577:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l571
						}
						position++
					}
				l576:
				l578:
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						{
							position580, tokenIndex580, depth580 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l581
							}
							position++
							goto l580
						l581:
							position, tokenIndex, depth = position580, tokenIndex580, depth580
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l582
							}
							position++
							goto l580
						l582:
							position, tokenIndex, depth = position580, tokenIndex580, depth580
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l583
							}
							position++
							goto l580
						l583:
							position, tokenIndex, depth = position580, tokenIndex580, depth580
							if buffer[position] != rune('_') {
								goto l584
							}
							position++
							goto l580
						l584:
							position, tokenIndex, depth = position580, tokenIndex580, depth580
							if buffer[position] != rune('.') {
								goto l579
							}
							position++
						}
					l580:
						goto l578
					l579:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
					}
					depth--
					add(rulePegText, position573)
				}
				if !_rules[ruleAction37]() {
					goto l571
				}
				depth--
				add(ruleRowValue, position572)
			}
			return true
		l571:
			position, tokenIndex, depth = position571, tokenIndex571, depth571
			return false
		},
		/* 53 NumericLiteral <- <(<('-'? [0-9]+)> Action38)> */
		func() bool {
			position585, tokenIndex585, depth585 := position, tokenIndex, depth
			{
				position586 := position
				depth++
				{
					position587 := position
					depth++
					{
						position588, tokenIndex588, depth588 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l588
						}
						position++
						goto l589
					l588:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
					}
				l589:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l585
					}
					position++
				l590:
					{
						position591, tokenIndex591, depth591 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex, depth = position591, tokenIndex591, depth591
					}
					depth--
					add(rulePegText, position587)
				}
				if !_rules[ruleAction38]() {
					goto l585
				}
				depth--
				add(ruleNumericLiteral, position586)
			}
			return true
		l585:
			position, tokenIndex, depth = position585, tokenIndex585, depth585
			return false
		},
		/* 54 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action39)> */
		func() bool {
			position592, tokenIndex592, depth592 := position, tokenIndex, depth
			{
				position593 := position
				depth++
				{
					position594 := position
					depth++
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l595
						}
						position++
						goto l596
					l595:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
					}
				l596:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l592
					}
					position++
				l597:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l598
						}
						position++
						goto l597
					l598:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
					}
					if buffer[position] != rune('.') {
						goto l592
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l592
					}
					position++
				l599:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l600
						}
						position++
						goto l599
					l600:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
					}
					depth--
					add(rulePegText, position594)
				}
				if !_rules[ruleAction39]() {
					goto l592
				}
				depth--
				add(ruleFloatLiteral, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 55 Function <- <(<ident> Action40)> */
		func() bool {
			position601, tokenIndex601, depth601 := position, tokenIndex, depth
			{
				position602 := position
				depth++
				{
					position603 := position
					depth++
					if !_rules[ruleident]() {
						goto l601
					}
					depth--
					add(rulePegText, position603)
				}
				if !_rules[ruleAction40]() {
					goto l601
				}
				depth--
				add(ruleFunction, position602)
			}
			return true
		l601:
			position, tokenIndex, depth = position601, tokenIndex601, depth601
			return false
		},
		/* 56 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position604, tokenIndex604, depth604 := position, tokenIndex, depth
			{
				position605 := position
				depth++
				{
					position606, tokenIndex606, depth606 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l607
					}
					goto l606
				l607:
					position, tokenIndex, depth = position606, tokenIndex606, depth606
					if !_rules[ruleFALSE]() {
						goto l604
					}
				}
			l606:
				depth--
				add(ruleBooleanLiteral, position605)
			}
			return true
		l604:
			position, tokenIndex, depth = position604, tokenIndex604, depth604
			return false
		},
		/* 57 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action41)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('T') {
							goto l608
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('R') {
							goto l608
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('U') {
							goto l608
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('E') {
							goto l608
						}
						position++
					}
				l617:
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction41]() {
					goto l608
				}
				depth--
				add(ruleTRUE, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 58 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action42)> */
		func() bool {
			position619, tokenIndex619, depth619 := position, tokenIndex, depth
			{
				position620 := position
				depth++
				{
					position621 := position
					depth++
					{
						position622, tokenIndex622, depth622 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l623
						}
						position++
						goto l622
					l623:
						position, tokenIndex, depth = position622, tokenIndex622, depth622
						if buffer[position] != rune('F') {
							goto l619
						}
						position++
					}
				l622:
					{
						position624, tokenIndex624, depth624 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('A') {
							goto l619
						}
						position++
					}
				l624:
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('L') {
							goto l619
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('S') {
							goto l619
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('E') {
							goto l619
						}
						position++
					}
				l630:
					depth--
					add(rulePegText, position621)
				}
				if !_rules[ruleAction42]() {
					goto l619
				}
				depth--
				add(ruleFALSE, position620)
			}
			return true
		l619:
			position, tokenIndex, depth = position619, tokenIndex619, depth619
			return false
		},
		/* 59 Wildcard <- <(<'*'> Action43)> */
		func() bool {
			position632, tokenIndex632, depth632 := position, tokenIndex, depth
			{
				position633 := position
				depth++
				{
					position634 := position
					depth++
					if buffer[position] != rune('*') {
						goto l632
					}
					position++
					depth--
					add(rulePegText, position634)
				}
				if !_rules[ruleAction43]() {
					goto l632
				}
				depth--
				add(ruleWildcard, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 60 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action44)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				{
					position637 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l635
					}
					position++
				l638:
					{
						position639, tokenIndex639, depth639 := position, tokenIndex, depth
						{
							position640, tokenIndex640, depth640 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l641
							}
							position++
							if buffer[position] != rune('\'') {
								goto l641
							}
							position++
							goto l640
						l641:
							position, tokenIndex, depth = position640, tokenIndex640, depth640
							{
								position642, tokenIndex642, depth642 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l642
								}
								position++
								goto l639
							l642:
								position, tokenIndex, depth = position642, tokenIndex642, depth642
							}
							if !matchDot() {
								goto l639
							}
						}
					l640:
						goto l638
					l639:
						position, tokenIndex, depth = position639, tokenIndex639, depth639
					}
					if buffer[position] != rune('\'') {
						goto l635
					}
					position++
					depth--
					add(rulePegText, position637)
				}
				if !_rules[ruleAction44]() {
					goto l635
				}
				depth--
				add(ruleStringLiteral, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 61 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action45)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645 := position
					depth++
					{
						position646, tokenIndex646, depth646 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l647
						}
						position++
						goto l646
					l647:
						position, tokenIndex, depth = position646, tokenIndex646, depth646
						if buffer[position] != rune('I') {
							goto l643
						}
						position++
					}
				l646:
					{
						position648, tokenIndex648, depth648 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l649
						}
						position++
						goto l648
					l649:
						position, tokenIndex, depth = position648, tokenIndex648, depth648
						if buffer[position] != rune('S') {
							goto l643
						}
						position++
					}
				l648:
					{
						position650, tokenIndex650, depth650 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l651
						}
						position++
						goto l650
					l651:
						position, tokenIndex, depth = position650, tokenIndex650, depth650
						if buffer[position] != rune('T') {
							goto l643
						}
						position++
					}
				l650:
					{
						position652, tokenIndex652, depth652 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l653
						}
						position++
						goto l652
					l653:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
						if buffer[position] != rune('R') {
							goto l643
						}
						position++
					}
				l652:
					{
						position654, tokenIndex654, depth654 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l655
						}
						position++
						goto l654
					l655:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
						if buffer[position] != rune('E') {
							goto l643
						}
						position++
					}
				l654:
					{
						position656, tokenIndex656, depth656 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l657
						}
						position++
						goto l656
					l657:
						position, tokenIndex, depth = position656, tokenIndex656, depth656
						if buffer[position] != rune('A') {
							goto l643
						}
						position++
					}
				l656:
					{
						position658, tokenIndex658, depth658 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l659
						}
						position++
						goto l658
					l659:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
						if buffer[position] != rune('M') {
							goto l643
						}
						position++
					}
				l658:
					depth--
					add(rulePegText, position645)
				}
				if !_rules[ruleAction45]() {
					goto l643
				}
				depth--
				add(ruleISTREAM, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 62 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action46)> */
		func() bool {
			position660, tokenIndex660, depth660 := position, tokenIndex, depth
			{
				position661 := position
				depth++
				{
					position662 := position
					depth++
					{
						position663, tokenIndex663, depth663 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l664
						}
						position++
						goto l663
					l664:
						position, tokenIndex, depth = position663, tokenIndex663, depth663
						if buffer[position] != rune('D') {
							goto l660
						}
						position++
					}
				l663:
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l666
						}
						position++
						goto l665
					l666:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if buffer[position] != rune('S') {
							goto l660
						}
						position++
					}
				l665:
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if buffer[position] != rune('T') {
							goto l660
						}
						position++
					}
				l667:
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l670
						}
						position++
						goto l669
					l670:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if buffer[position] != rune('R') {
							goto l660
						}
						position++
					}
				l669:
					{
						position671, tokenIndex671, depth671 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l672
						}
						position++
						goto l671
					l672:
						position, tokenIndex, depth = position671, tokenIndex671, depth671
						if buffer[position] != rune('E') {
							goto l660
						}
						position++
					}
				l671:
					{
						position673, tokenIndex673, depth673 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l674
						}
						position++
						goto l673
					l674:
						position, tokenIndex, depth = position673, tokenIndex673, depth673
						if buffer[position] != rune('A') {
							goto l660
						}
						position++
					}
				l673:
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l676
						}
						position++
						goto l675
					l676:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('M') {
							goto l660
						}
						position++
					}
				l675:
					depth--
					add(rulePegText, position662)
				}
				if !_rules[ruleAction46]() {
					goto l660
				}
				depth--
				add(ruleDSTREAM, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 63 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action47)> */
		func() bool {
			position677, tokenIndex677, depth677 := position, tokenIndex, depth
			{
				position678 := position
				depth++
				{
					position679 := position
					depth++
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l681
						}
						position++
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if buffer[position] != rune('R') {
							goto l677
						}
						position++
					}
				l680:
					{
						position682, tokenIndex682, depth682 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l683
						}
						position++
						goto l682
					l683:
						position, tokenIndex, depth = position682, tokenIndex682, depth682
						if buffer[position] != rune('S') {
							goto l677
						}
						position++
					}
				l682:
					{
						position684, tokenIndex684, depth684 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l685
						}
						position++
						goto l684
					l685:
						position, tokenIndex, depth = position684, tokenIndex684, depth684
						if buffer[position] != rune('T') {
							goto l677
						}
						position++
					}
				l684:
					{
						position686, tokenIndex686, depth686 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
						if buffer[position] != rune('R') {
							goto l677
						}
						position++
					}
				l686:
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if buffer[position] != rune('E') {
							goto l677
						}
						position++
					}
				l688:
					{
						position690, tokenIndex690, depth690 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l691
						}
						position++
						goto l690
					l691:
						position, tokenIndex, depth = position690, tokenIndex690, depth690
						if buffer[position] != rune('A') {
							goto l677
						}
						position++
					}
				l690:
					{
						position692, tokenIndex692, depth692 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l693
						}
						position++
						goto l692
					l693:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						if buffer[position] != rune('M') {
							goto l677
						}
						position++
					}
				l692:
					depth--
					add(rulePegText, position679)
				}
				if !_rules[ruleAction47]() {
					goto l677
				}
				depth--
				add(ruleRSTREAM, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 64 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action48)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				{
					position696 := position
					depth++
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('T') {
							goto l694
						}
						position++
					}
				l697:
					{
						position699, tokenIndex699, depth699 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l700
						}
						position++
						goto l699
					l700:
						position, tokenIndex, depth = position699, tokenIndex699, depth699
						if buffer[position] != rune('U') {
							goto l694
						}
						position++
					}
				l699:
					{
						position701, tokenIndex701, depth701 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l702
						}
						position++
						goto l701
					l702:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
						if buffer[position] != rune('P') {
							goto l694
						}
						position++
					}
				l701:
					{
						position703, tokenIndex703, depth703 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l704
						}
						position++
						goto l703
					l704:
						position, tokenIndex, depth = position703, tokenIndex703, depth703
						if buffer[position] != rune('L') {
							goto l694
						}
						position++
					}
				l703:
					{
						position705, tokenIndex705, depth705 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l706
						}
						position++
						goto l705
					l706:
						position, tokenIndex, depth = position705, tokenIndex705, depth705
						if buffer[position] != rune('E') {
							goto l694
						}
						position++
					}
				l705:
					{
						position707, tokenIndex707, depth707 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l708
						}
						position++
						goto l707
					l708:
						position, tokenIndex, depth = position707, tokenIndex707, depth707
						if buffer[position] != rune('S') {
							goto l694
						}
						position++
					}
				l707:
					depth--
					add(rulePegText, position696)
				}
				if !_rules[ruleAction48]() {
					goto l694
				}
				depth--
				add(ruleTUPLES, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 65 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action49)> */
		func() bool {
			position709, tokenIndex709, depth709 := position, tokenIndex, depth
			{
				position710 := position
				depth++
				{
					position711 := position
					depth++
					{
						position712, tokenIndex712, depth712 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l713
						}
						position++
						goto l712
					l713:
						position, tokenIndex, depth = position712, tokenIndex712, depth712
						if buffer[position] != rune('S') {
							goto l709
						}
						position++
					}
				l712:
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l715
						}
						position++
						goto l714
					l715:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
						if buffer[position] != rune('E') {
							goto l709
						}
						position++
					}
				l714:
					{
						position716, tokenIndex716, depth716 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l717
						}
						position++
						goto l716
					l717:
						position, tokenIndex, depth = position716, tokenIndex716, depth716
						if buffer[position] != rune('C') {
							goto l709
						}
						position++
					}
				l716:
					{
						position718, tokenIndex718, depth718 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l719
						}
						position++
						goto l718
					l719:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
						if buffer[position] != rune('O') {
							goto l709
						}
						position++
					}
				l718:
					{
						position720, tokenIndex720, depth720 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l721
						}
						position++
						goto l720
					l721:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
						if buffer[position] != rune('N') {
							goto l709
						}
						position++
					}
				l720:
					{
						position722, tokenIndex722, depth722 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l723
						}
						position++
						goto l722
					l723:
						position, tokenIndex, depth = position722, tokenIndex722, depth722
						if buffer[position] != rune('D') {
							goto l709
						}
						position++
					}
				l722:
					{
						position724, tokenIndex724, depth724 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l725
						}
						position++
						goto l724
					l725:
						position, tokenIndex, depth = position724, tokenIndex724, depth724
						if buffer[position] != rune('S') {
							goto l709
						}
						position++
					}
				l724:
					depth--
					add(rulePegText, position711)
				}
				if !_rules[ruleAction49]() {
					goto l709
				}
				depth--
				add(ruleSECONDS, position710)
			}
			return true
		l709:
			position, tokenIndex, depth = position709, tokenIndex709, depth709
			return false
		},
		/* 66 StreamIdentifier <- <(<ident> Action50)> */
		func() bool {
			position726, tokenIndex726, depth726 := position, tokenIndex, depth
			{
				position727 := position
				depth++
				{
					position728 := position
					depth++
					if !_rules[ruleident]() {
						goto l726
					}
					depth--
					add(rulePegText, position728)
				}
				if !_rules[ruleAction50]() {
					goto l726
				}
				depth--
				add(ruleStreamIdentifier, position727)
			}
			return true
		l726:
			position, tokenIndex, depth = position726, tokenIndex726, depth726
			return false
		},
		/* 67 SourceSinkType <- <(<ident> Action51)> */
		func() bool {
			position729, tokenIndex729, depth729 := position, tokenIndex, depth
			{
				position730 := position
				depth++
				{
					position731 := position
					depth++
					if !_rules[ruleident]() {
						goto l729
					}
					depth--
					add(rulePegText, position731)
				}
				if !_rules[ruleAction51]() {
					goto l729
				}
				depth--
				add(ruleSourceSinkType, position730)
			}
			return true
		l729:
			position, tokenIndex, depth = position729, tokenIndex729, depth729
			return false
		},
		/* 68 SourceSinkParamKey <- <(<ident> Action52)> */
		func() bool {
			position732, tokenIndex732, depth732 := position, tokenIndex, depth
			{
				position733 := position
				depth++
				{
					position734 := position
					depth++
					if !_rules[ruleident]() {
						goto l732
					}
					depth--
					add(rulePegText, position734)
				}
				if !_rules[ruleAction52]() {
					goto l732
				}
				depth--
				add(ruleSourceSinkParamKey, position733)
			}
			return true
		l732:
			position, tokenIndex, depth = position732, tokenIndex732, depth732
			return false
		},
		/* 69 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action53)> */
		func() bool {
			position735, tokenIndex735, depth735 := position, tokenIndex, depth
			{
				position736 := position
				depth++
				{
					position737 := position
					depth++
					{
						position738, tokenIndex738, depth738 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l739
						}
						position++
						goto l738
					l739:
						position, tokenIndex, depth = position738, tokenIndex738, depth738
						if buffer[position] != rune('O') {
							goto l735
						}
						position++
					}
				l738:
					{
						position740, tokenIndex740, depth740 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l741
						}
						position++
						goto l740
					l741:
						position, tokenIndex, depth = position740, tokenIndex740, depth740
						if buffer[position] != rune('R') {
							goto l735
						}
						position++
					}
				l740:
					depth--
					add(rulePegText, position737)
				}
				if !_rules[ruleAction53]() {
					goto l735
				}
				depth--
				add(ruleOr, position736)
			}
			return true
		l735:
			position, tokenIndex, depth = position735, tokenIndex735, depth735
			return false
		},
		/* 70 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action54)> */
		func() bool {
			position742, tokenIndex742, depth742 := position, tokenIndex, depth
			{
				position743 := position
				depth++
				{
					position744 := position
					depth++
					{
						position745, tokenIndex745, depth745 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l746
						}
						position++
						goto l745
					l746:
						position, tokenIndex, depth = position745, tokenIndex745, depth745
						if buffer[position] != rune('A') {
							goto l742
						}
						position++
					}
				l745:
					{
						position747, tokenIndex747, depth747 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l748
						}
						position++
						goto l747
					l748:
						position, tokenIndex, depth = position747, tokenIndex747, depth747
						if buffer[position] != rune('N') {
							goto l742
						}
						position++
					}
				l747:
					{
						position749, tokenIndex749, depth749 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l750
						}
						position++
						goto l749
					l750:
						position, tokenIndex, depth = position749, tokenIndex749, depth749
						if buffer[position] != rune('D') {
							goto l742
						}
						position++
					}
				l749:
					depth--
					add(rulePegText, position744)
				}
				if !_rules[ruleAction54]() {
					goto l742
				}
				depth--
				add(ruleAnd, position743)
			}
			return true
		l742:
			position, tokenIndex, depth = position742, tokenIndex742, depth742
			return false
		},
		/* 71 Equal <- <(<'='> Action55)> */
		func() bool {
			position751, tokenIndex751, depth751 := position, tokenIndex, depth
			{
				position752 := position
				depth++
				{
					position753 := position
					depth++
					if buffer[position] != rune('=') {
						goto l751
					}
					position++
					depth--
					add(rulePegText, position753)
				}
				if !_rules[ruleAction55]() {
					goto l751
				}
				depth--
				add(ruleEqual, position752)
			}
			return true
		l751:
			position, tokenIndex, depth = position751, tokenIndex751, depth751
			return false
		},
		/* 72 Less <- <(<'<'> Action56)> */
		func() bool {
			position754, tokenIndex754, depth754 := position, tokenIndex, depth
			{
				position755 := position
				depth++
				{
					position756 := position
					depth++
					if buffer[position] != rune('<') {
						goto l754
					}
					position++
					depth--
					add(rulePegText, position756)
				}
				if !_rules[ruleAction56]() {
					goto l754
				}
				depth--
				add(ruleLess, position755)
			}
			return true
		l754:
			position, tokenIndex, depth = position754, tokenIndex754, depth754
			return false
		},
		/* 73 LessOrEqual <- <(<('<' '=')> Action57)> */
		func() bool {
			position757, tokenIndex757, depth757 := position, tokenIndex, depth
			{
				position758 := position
				depth++
				{
					position759 := position
					depth++
					if buffer[position] != rune('<') {
						goto l757
					}
					position++
					if buffer[position] != rune('=') {
						goto l757
					}
					position++
					depth--
					add(rulePegText, position759)
				}
				if !_rules[ruleAction57]() {
					goto l757
				}
				depth--
				add(ruleLessOrEqual, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 74 Greater <- <(<'>'> Action58)> */
		func() bool {
			position760, tokenIndex760, depth760 := position, tokenIndex, depth
			{
				position761 := position
				depth++
				{
					position762 := position
					depth++
					if buffer[position] != rune('>') {
						goto l760
					}
					position++
					depth--
					add(rulePegText, position762)
				}
				if !_rules[ruleAction58]() {
					goto l760
				}
				depth--
				add(ruleGreater, position761)
			}
			return true
		l760:
			position, tokenIndex, depth = position760, tokenIndex760, depth760
			return false
		},
		/* 75 GreaterOrEqual <- <(<('>' '=')> Action59)> */
		func() bool {
			position763, tokenIndex763, depth763 := position, tokenIndex, depth
			{
				position764 := position
				depth++
				{
					position765 := position
					depth++
					if buffer[position] != rune('>') {
						goto l763
					}
					position++
					if buffer[position] != rune('=') {
						goto l763
					}
					position++
					depth--
					add(rulePegText, position765)
				}
				if !_rules[ruleAction59]() {
					goto l763
				}
				depth--
				add(ruleGreaterOrEqual, position764)
			}
			return true
		l763:
			position, tokenIndex, depth = position763, tokenIndex763, depth763
			return false
		},
		/* 76 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action60)> */
		func() bool {
			position766, tokenIndex766, depth766 := position, tokenIndex, depth
			{
				position767 := position
				depth++
				{
					position768 := position
					depth++
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l770
						}
						position++
						if buffer[position] != rune('=') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('<') {
							goto l766
						}
						position++
						if buffer[position] != rune('>') {
							goto l766
						}
						position++
					}
				l769:
					depth--
					add(rulePegText, position768)
				}
				if !_rules[ruleAction60]() {
					goto l766
				}
				depth--
				add(ruleNotEqual, position767)
			}
			return true
		l766:
			position, tokenIndex, depth = position766, tokenIndex766, depth766
			return false
		},
		/* 77 Plus <- <(<'+'> Action61)> */
		func() bool {
			position771, tokenIndex771, depth771 := position, tokenIndex, depth
			{
				position772 := position
				depth++
				{
					position773 := position
					depth++
					if buffer[position] != rune('+') {
						goto l771
					}
					position++
					depth--
					add(rulePegText, position773)
				}
				if !_rules[ruleAction61]() {
					goto l771
				}
				depth--
				add(rulePlus, position772)
			}
			return true
		l771:
			position, tokenIndex, depth = position771, tokenIndex771, depth771
			return false
		},
		/* 78 Minus <- <(<'-'> Action62)> */
		func() bool {
			position774, tokenIndex774, depth774 := position, tokenIndex, depth
			{
				position775 := position
				depth++
				{
					position776 := position
					depth++
					if buffer[position] != rune('-') {
						goto l774
					}
					position++
					depth--
					add(rulePegText, position776)
				}
				if !_rules[ruleAction62]() {
					goto l774
				}
				depth--
				add(ruleMinus, position775)
			}
			return true
		l774:
			position, tokenIndex, depth = position774, tokenIndex774, depth774
			return false
		},
		/* 79 Multiply <- <(<'*'> Action63)> */
		func() bool {
			position777, tokenIndex777, depth777 := position, tokenIndex, depth
			{
				position778 := position
				depth++
				{
					position779 := position
					depth++
					if buffer[position] != rune('*') {
						goto l777
					}
					position++
					depth--
					add(rulePegText, position779)
				}
				if !_rules[ruleAction63]() {
					goto l777
				}
				depth--
				add(ruleMultiply, position778)
			}
			return true
		l777:
			position, tokenIndex, depth = position777, tokenIndex777, depth777
			return false
		},
		/* 80 Divide <- <(<'/'> Action64)> */
		func() bool {
			position780, tokenIndex780, depth780 := position, tokenIndex, depth
			{
				position781 := position
				depth++
				{
					position782 := position
					depth++
					if buffer[position] != rune('/') {
						goto l780
					}
					position++
					depth--
					add(rulePegText, position782)
				}
				if !_rules[ruleAction64]() {
					goto l780
				}
				depth--
				add(ruleDivide, position781)
			}
			return true
		l780:
			position, tokenIndex, depth = position780, tokenIndex780, depth780
			return false
		},
		/* 81 Modulo <- <(<'%'> Action65)> */
		func() bool {
			position783, tokenIndex783, depth783 := position, tokenIndex, depth
			{
				position784 := position
				depth++
				{
					position785 := position
					depth++
					if buffer[position] != rune('%') {
						goto l783
					}
					position++
					depth--
					add(rulePegText, position785)
				}
				if !_rules[ruleAction65]() {
					goto l783
				}
				depth--
				add(ruleModulo, position784)
			}
			return true
		l783:
			position, tokenIndex, depth = position783, tokenIndex783, depth783
			return false
		},
		/* 82 Identifier <- <(<ident> Action66)> */
		func() bool {
			position786, tokenIndex786, depth786 := position, tokenIndex, depth
			{
				position787 := position
				depth++
				{
					position788 := position
					depth++
					if !_rules[ruleident]() {
						goto l786
					}
					depth--
					add(rulePegText, position788)
				}
				if !_rules[ruleAction66]() {
					goto l786
				}
				depth--
				add(ruleIdentifier, position787)
			}
			return true
		l786:
			position, tokenIndex, depth = position786, tokenIndex786, depth786
			return false
		},
		/* 83 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position789, tokenIndex789, depth789 := position, tokenIndex, depth
			{
				position790 := position
				depth++
				{
					position791, tokenIndex791, depth791 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l792
					}
					position++
					goto l791
				l792:
					position, tokenIndex, depth = position791, tokenIndex791, depth791
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l789
					}
					position++
				}
			l791:
			l793:
				{
					position794, tokenIndex794, depth794 := position, tokenIndex, depth
					{
						position795, tokenIndex795, depth795 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l796
						}
						position++
						goto l795
					l796:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l797
						}
						position++
						goto l795
					l797:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l798
						}
						position++
						goto l795
					l798:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						if buffer[position] != rune('_') {
							goto l794
						}
						position++
					}
				l795:
					goto l793
				l794:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
				}
				depth--
				add(ruleident, position790)
			}
			return true
		l789:
			position, tokenIndex, depth = position789, tokenIndex789, depth789
			return false
		},
		/* 84 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position800 := position
				depth++
			l801:
				{
					position802, tokenIndex802, depth802 := position, tokenIndex, depth
					{
						position803, tokenIndex803, depth803 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l804
						}
						position++
						goto l803
					l804:
						position, tokenIndex, depth = position803, tokenIndex803, depth803
						if buffer[position] != rune('\t') {
							goto l805
						}
						position++
						goto l803
					l805:
						position, tokenIndex, depth = position803, tokenIndex803, depth803
						if buffer[position] != rune('\n') {
							goto l802
						}
						position++
					}
				l803:
					goto l801
				l802:
					position, tokenIndex, depth = position802, tokenIndex802, depth802
				}
				depth--
				add(rulesp, position800)
			}
			return true
		},
		/* 86 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 87 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 88 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 89 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 90 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 91 Action5 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 92 Action6 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 93 Action7 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 95 Action8 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 96 Action9 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 97 Action10 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 98 Action11 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 99 Action12 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 100 Action13 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 101 Action14 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 102 Action15 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 103 Action16 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 104 Action17 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 105 Action18 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 106 Action19 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 107 Action20 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 108 Action21 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 109 Action22 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 110 Action23 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 111 Action24 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 112 Action25 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 113 Action26 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 114 Action27 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 115 Action28 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 116 Action29 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 117 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 118 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 119 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 120 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 121 Action34 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 122 Action35 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 123 Action36 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 124 Action37 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 125 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 126 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 127 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 128 Action41 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 129 Action42 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 130 Action43 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 131 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 132 Action45 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 133 Action46 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 134 Action47 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 135 Action48 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 136 Action49 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 137 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 138 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 139 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 140 Action53 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 141 Action54 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 142 Action55 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 143 Action56 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 144 Action57 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 145 Action58 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 146 Action59 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 147 Action60 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 148 Action61 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 149 Action62 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 150 Action63 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 151 Action64 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 152 Action65 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 153 Action66 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
	}
	p.rules = _rules
}

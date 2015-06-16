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
	ruleSourceSinkParamVal
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
	rulePegText
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
	"SourceSinkParamVal",
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
	"PegText",
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
	rules  [153]func() bool
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

			p.AssembleCreateStreamFromSource()

		case ruleAction5:

			p.AssembleCreateStreamFromSourceExt()

		case ruleAction6:

			p.AssembleInsertIntoSelect()

		case ruleAction7:

			p.AssembleEmitter(begin, end)

		case ruleAction8:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction9:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction10:

			p.AssembleStreamEmitInterval()

		case ruleAction11:

			p.AssembleProjections(begin, end)

		case ruleAction12:

			p.AssembleAlias()

		case ruleAction13:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction14:

			p.AssembleWindowedFrom(begin, end)

		case ruleAction15:

			p.AssembleInterval()

		case ruleAction16:

			p.AssembleInterval()

		case ruleAction17:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction18:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction19:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction20:

			p.EnsureAliasedStreamWindow()

		case ruleAction21:

			p.EnsureAliasedStreamWindow()

		case ruleAction22:

			p.AssembleAliasedStreamWindow()

		case ruleAction23:

			p.AssembleAliasedStreamWindow()

		case ruleAction24:

			p.AssembleStreamWindow()

		case ruleAction25:

			p.AssembleStreamWindow()

		case ruleAction26:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction27:

			p.AssembleSourceSinkParam()

		case ruleAction28:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction29:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction30:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction31:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction32:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction33:

			p.AssembleFuncApp()

		case ruleAction34:

			p.AssembleExpressions(begin, end)

		case ruleAction35:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction36:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction37:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction38:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction40:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction41:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction42:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction44:

			p.PushComponent(begin, end, Istream)

		case ruleAction45:

			p.PushComponent(begin, end, Dstream)

		case ruleAction46:

			p.PushComponent(begin, end, Rstream)

		case ruleAction47:

			p.PushComponent(begin, end, Tuples)

		case ruleAction48:

			p.PushComponent(begin, end, Seconds)

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamVal(substr))

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateStreamFromSourceStmt / CreateStreamFromSourceExtStmt / CreateSinkStmt / InsertIntoSelectStmt)> */
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
		/* 2 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp DefWindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
		func() bool {
			position16, tokenIndex16, depth16 := position, tokenIndex, depth
			{
				position17 := position
				depth++
				{
					position18, tokenIndex18, depth18 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l19
					}
					position++
					goto l18
				l19:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					if buffer[position] != rune('S') {
						goto l16
					}
					position++
				}
			l18:
				{
					position20, tokenIndex20, depth20 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l21
					}
					position++
					goto l20
				l21:
					position, tokenIndex, depth = position20, tokenIndex20, depth20
					if buffer[position] != rune('E') {
						goto l16
					}
					position++
				}
			l20:
				{
					position22, tokenIndex22, depth22 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l23
					}
					position++
					goto l22
				l23:
					position, tokenIndex, depth = position22, tokenIndex22, depth22
					if buffer[position] != rune('L') {
						goto l16
					}
					position++
				}
			l22:
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l25
					}
					position++
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if buffer[position] != rune('E') {
						goto l16
					}
					position++
				}
			l24:
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l27
					}
					position++
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					if buffer[position] != rune('C') {
						goto l16
					}
					position++
				}
			l26:
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l29
					}
					position++
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if buffer[position] != rune('T') {
						goto l16
					}
					position++
				}
			l28:
				if !_rules[rulesp]() {
					goto l16
				}
				if !_rules[ruleEmitter]() {
					goto l16
				}
				if !_rules[rulesp]() {
					goto l16
				}
				if !_rules[ruleProjections]() {
					goto l16
				}
				if !_rules[rulesp]() {
					goto l16
				}
				if !_rules[ruleDefWindowedFrom]() {
					goto l16
				}
				if !_rules[rulesp]() {
					goto l16
				}
				if !_rules[ruleFilter]() {
					goto l16
				}
				if !_rules[rulesp]() {
					goto l16
				}
				if !_rules[ruleGrouping]() {
					goto l16
				}
				if !_rules[rulesp]() {
					goto l16
				}
				if !_rules[ruleHaving]() {
					goto l16
				}
				if !_rules[rulesp]() {
					goto l16
				}
				if !_rules[ruleAction0]() {
					goto l16
				}
				depth--
				add(ruleSelectStmt, position17)
			}
			return true
		l16:
			position, tokenIndex, depth = position16, tokenIndex16, depth16
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position30, tokenIndex30, depth30 := position, tokenIndex, depth
			{
				position31 := position
				depth++
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l33
					}
					position++
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if buffer[position] != rune('C') {
						goto l30
					}
					position++
				}
			l32:
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l35
					}
					position++
					goto l34
				l35:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if buffer[position] != rune('R') {
						goto l30
					}
					position++
				}
			l34:
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l37
					}
					position++
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if buffer[position] != rune('E') {
						goto l30
					}
					position++
				}
			l36:
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l39
					}
					position++
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if buffer[position] != rune('A') {
						goto l30
					}
					position++
				}
			l38:
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('T') {
						goto l30
					}
					position++
				}
			l40:
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('E') {
						goto l30
					}
					position++
				}
			l42:
				if !_rules[rulesp]() {
					goto l30
				}
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l45
					}
					position++
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if buffer[position] != rune('S') {
						goto l30
					}
					position++
				}
			l44:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('T') {
						goto l30
					}
					position++
				}
			l46:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('R') {
						goto l30
					}
					position++
				}
			l48:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('E') {
						goto l30
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('A') {
						goto l30
					}
					position++
				}
			l52:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('M') {
						goto l30
					}
					position++
				}
			l54:
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l30
				}
				if !_rules[rulesp]() {
					goto l30
				}
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('A') {
						goto l30
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('S') {
						goto l30
					}
					position++
				}
			l58:
				if !_rules[rulesp]() {
					goto l30
				}
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l61
					}
					position++
					goto l60
				l61:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
					if buffer[position] != rune('S') {
						goto l30
					}
					position++
				}
			l60:
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
					if buffer[position] != rune('E') {
						goto l30
					}
					position++
				}
			l62:
				{
					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l65
					}
					position++
					goto l64
				l65:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
					if buffer[position] != rune('L') {
						goto l30
					}
					position++
				}
			l64:
				{
					position66, tokenIndex66, depth66 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex, depth = position66, tokenIndex66, depth66
					if buffer[position] != rune('E') {
						goto l30
					}
					position++
				}
			l66:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
					if buffer[position] != rune('C') {
						goto l30
					}
					position++
				}
			l68:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l71
					}
					position++
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if buffer[position] != rune('T') {
						goto l30
					}
					position++
				}
			l70:
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleEmitter]() {
					goto l30
				}
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleProjections]() {
					goto l30
				}
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleWindowedFrom]() {
					goto l30
				}
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleFilter]() {
					goto l30
				}
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleGrouping]() {
					goto l30
				}
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleHaving]() {
					goto l30
				}
				if !_rules[rulesp]() {
					goto l30
				}
				if !_rules[ruleAction1]() {
					goto l30
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position31)
			}
			return true
		l30:
			position, tokenIndex, depth = position30, tokenIndex30, depth30
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position72, tokenIndex72, depth72 := position, tokenIndex, depth
			{
				position73 := position
				depth++
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if buffer[position] != rune('C') {
						goto l72
					}
					position++
				}
			l74:
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('R') {
						goto l72
					}
					position++
				}
			l76:
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l79
					}
					position++
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('E') {
						goto l72
					}
					position++
				}
			l78:
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('A') {
						goto l72
					}
					position++
				}
			l80:
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l83
					}
					position++
					goto l82
				l83:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
					if buffer[position] != rune('T') {
						goto l72
					}
					position++
				}
			l82:
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
					if buffer[position] != rune('E') {
						goto l72
					}
					position++
				}
			l84:
				if !_rules[rulesp]() {
					goto l72
				}
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if buffer[position] != rune('S') {
						goto l72
					}
					position++
				}
			l86:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l89
					}
					position++
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if buffer[position] != rune('O') {
						goto l72
					}
					position++
				}
			l88:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('U') {
						goto l72
					}
					position++
				}
			l90:
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('R') {
						goto l72
					}
					position++
				}
			l92:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('C') {
						goto l72
					}
					position++
				}
			l94:
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('E') {
						goto l72
					}
					position++
				}
			l96:
				if !_rules[rulesp]() {
					goto l72
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l72
				}
				if !_rules[rulesp]() {
					goto l72
				}
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('T') {
						goto l72
					}
					position++
				}
			l98:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('Y') {
						goto l72
					}
					position++
				}
			l100:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('P') {
						goto l72
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('E') {
						goto l72
					}
					position++
				}
			l104:
				if !_rules[rulesp]() {
					goto l72
				}
				if !_rules[ruleSourceSinkType]() {
					goto l72
				}
				if !_rules[rulesp]() {
					goto l72
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l72
				}
				if !_rules[ruleAction2]() {
					goto l72
				}
				depth--
				add(ruleCreateSourceStmt, position73)
			}
			return true
		l72:
			position, tokenIndex, depth = position72, tokenIndex72, depth72
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position106, tokenIndex106, depth106 := position, tokenIndex, depth
			{
				position107 := position
				depth++
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('C') {
						goto l106
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('R') {
						goto l106
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('E') {
						goto l106
					}
					position++
				}
			l112:
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('A') {
						goto l106
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('T') {
						goto l106
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('E') {
						goto l106
					}
					position++
				}
			l118:
				if !_rules[rulesp]() {
					goto l106
				}
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
					if buffer[position] != rune('S') {
						goto l106
					}
					position++
				}
			l120:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('I') {
						goto l106
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('N') {
						goto l106
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('K') {
						goto l106
					}
					position++
				}
			l126:
				if !_rules[rulesp]() {
					goto l106
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l106
				}
				if !_rules[rulesp]() {
					goto l106
				}
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('T') {
						goto l106
					}
					position++
				}
			l128:
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('Y') {
						goto l106
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('P') {
						goto l106
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('E') {
						goto l106
					}
					position++
				}
			l134:
				if !_rules[rulesp]() {
					goto l106
				}
				if !_rules[ruleSourceSinkType]() {
					goto l106
				}
				if !_rules[rulesp]() {
					goto l106
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l106
				}
				if !_rules[ruleAction3]() {
					goto l106
				}
				depth--
				add(ruleCreateSinkStmt, position107)
			}
			return true
		l106:
			position, tokenIndex, depth = position106, tokenIndex106, depth106
			return false
		},
		/* 6 CreateStreamFromSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action4)> */
		func() bool {
			position136, tokenIndex136, depth136 := position, tokenIndex, depth
			{
				position137 := position
				depth++
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('C') {
						goto l136
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('R') {
						goto l136
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('E') {
						goto l136
					}
					position++
				}
			l142:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('A') {
						goto l136
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('T') {
						goto l136
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('E') {
						goto l136
					}
					position++
				}
			l148:
				if !_rules[rulesp]() {
					goto l136
				}
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('S') {
						goto l136
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('T') {
						goto l136
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('R') {
						goto l136
					}
					position++
				}
			l154:
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('E') {
						goto l136
					}
					position++
				}
			l156:
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if buffer[position] != rune('A') {
						goto l136
					}
					position++
				}
			l158:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('M') {
						goto l136
					}
					position++
				}
			l160:
				if !_rules[rulesp]() {
					goto l136
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l136
				}
				if !_rules[rulesp]() {
					goto l136
				}
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('F') {
						goto l136
					}
					position++
				}
			l162:
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('R') {
						goto l136
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('O') {
						goto l136
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('M') {
						goto l136
					}
					position++
				}
			l168:
				if !_rules[rulesp]() {
					goto l136
				}
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('S') {
						goto l136
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('O') {
						goto l136
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('U') {
						goto l136
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('R') {
						goto l136
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('C') {
						goto l136
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('E') {
						goto l136
					}
					position++
				}
			l180:
				if !_rules[rulesp]() {
					goto l136
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l136
				}
				if !_rules[ruleAction4]() {
					goto l136
				}
				depth--
				add(ruleCreateStreamFromSourceStmt, position137)
			}
			return true
		l136:
			position, tokenIndex, depth = position136, tokenIndex136, depth136
			return false
		},
		/* 7 CreateStreamFromSourceExtStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp SourceSinkType sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp SourceSinkSpecs Action5)> */
		func() bool {
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('C') {
						goto l182
					}
					position++
				}
			l184:
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('R') {
						goto l182
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('E') {
						goto l182
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('A') {
						goto l182
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('T') {
						goto l182
					}
					position++
				}
			l192:
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('E') {
						goto l182
					}
					position++
				}
			l194:
				if !_rules[rulesp]() {
					goto l182
				}
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('S') {
						goto l182
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('T') {
						goto l182
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('R') {
						goto l182
					}
					position++
				}
			l200:
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('E') {
						goto l182
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('A') {
						goto l182
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('M') {
						goto l182
					}
					position++
				}
			l206:
				if !_rules[rulesp]() {
					goto l182
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l182
				}
				if !_rules[rulesp]() {
					goto l182
				}
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('F') {
						goto l182
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('R') {
						goto l182
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('O') {
						goto l182
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('M') {
						goto l182
					}
					position++
				}
			l214:
				if !_rules[rulesp]() {
					goto l182
				}
				if !_rules[ruleSourceSinkType]() {
					goto l182
				}
				if !_rules[rulesp]() {
					goto l182
				}
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('S') {
						goto l182
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('O') {
						goto l182
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('U') {
						goto l182
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('R') {
						goto l182
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('C') {
						goto l182
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('E') {
						goto l182
					}
					position++
				}
			l226:
				if !_rules[rulesp]() {
					goto l182
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l182
				}
				if !_rules[ruleAction5]() {
					goto l182
				}
				depth--
				add(ruleCreateStreamFromSourceExtStmt, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
			return false
		},
		/* 8 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action6)> */
		func() bool {
			position228, tokenIndex228, depth228 := position, tokenIndex, depth
			{
				position229 := position
				depth++
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
					if buffer[position] != rune('I') {
						goto l228
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('N') {
						goto l228
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('S') {
						goto l228
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('E') {
						goto l228
					}
					position++
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('R') {
						goto l228
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('T') {
						goto l228
					}
					position++
				}
			l240:
				if !_rules[rulesp]() {
					goto l228
				}
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('I') {
						goto l228
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('N') {
						goto l228
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('T') {
						goto l228
					}
					position++
				}
			l246:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('O') {
						goto l228
					}
					position++
				}
			l248:
				if !_rules[rulesp]() {
					goto l228
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l228
				}
				if !_rules[rulesp]() {
					goto l228
				}
				if !_rules[ruleSelectStmt]() {
					goto l228
				}
				if !_rules[ruleAction6]() {
					goto l228
				}
				depth--
				add(ruleInsertIntoSelectStmt, position229)
			}
			return true
		l228:
			position, tokenIndex, depth = position228, tokenIndex228, depth228
			return false
		},
		/* 9 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action7)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l253
					}
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if !_rules[ruleDSTREAM]() {
						goto l254
					}
					goto l252
				l254:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if !_rules[ruleRSTREAM]() {
						goto l250
					}
				}
			l252:
				{
					position255 := position
					depth++
					{
						position256, tokenIndex256, depth256 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l256
						}
						if buffer[position] != rune('[') {
							goto l256
						}
						position++
						if !_rules[rulesp]() {
							goto l256
						}
						{
							position258, tokenIndex258, depth258 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l259
							}
							position++
							goto l258
						l259:
							position, tokenIndex, depth = position258, tokenIndex258, depth258
							if buffer[position] != rune('E') {
								goto l256
							}
							position++
						}
					l258:
						{
							position260, tokenIndex260, depth260 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l261
							}
							position++
							goto l260
						l261:
							position, tokenIndex, depth = position260, tokenIndex260, depth260
							if buffer[position] != rune('V') {
								goto l256
							}
							position++
						}
					l260:
						{
							position262, tokenIndex262, depth262 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l263
							}
							position++
							goto l262
						l263:
							position, tokenIndex, depth = position262, tokenIndex262, depth262
							if buffer[position] != rune('E') {
								goto l256
							}
							position++
						}
					l262:
						{
							position264, tokenIndex264, depth264 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l265
							}
							position++
							goto l264
						l265:
							position, tokenIndex, depth = position264, tokenIndex264, depth264
							if buffer[position] != rune('R') {
								goto l256
							}
							position++
						}
					l264:
						{
							position266, tokenIndex266, depth266 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l267
							}
							position++
							goto l266
						l267:
							position, tokenIndex, depth = position266, tokenIndex266, depth266
							if buffer[position] != rune('Y') {
								goto l256
							}
							position++
						}
					l266:
						if !_rules[rulesp]() {
							goto l256
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l256
						}
						if !_rules[rulesp]() {
							goto l256
						}
						if buffer[position] != rune(']') {
							goto l256
						}
						position++
						goto l257
					l256:
						position, tokenIndex, depth = position256, tokenIndex256, depth256
					}
				l257:
					depth--
					add(rulePegText, position255)
				}
				if !_rules[ruleAction7]() {
					goto l250
				}
				depth--
				add(ruleEmitter, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 10 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position268, tokenIndex268, depth268 := position, tokenIndex, depth
			{
				position269 := position
				depth++
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l271
					}
				l272:
					{
						position273, tokenIndex273, depth273 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l273
						}
						if buffer[position] != rune(',') {
							goto l273
						}
						position++
						if !_rules[rulesp]() {
							goto l273
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l273
						}
						goto l272
					l273:
						position, tokenIndex, depth = position273, tokenIndex273, depth273
					}
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if !_rules[ruleTimeEmitterInterval]() {
						goto l274
					}
					goto l270
				l274:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if !_rules[ruleTupleEmitterInterval]() {
						goto l268
					}
				}
			l270:
				depth--
				add(ruleEmitterIntervals, position269)
			}
			return true
		l268:
			position, tokenIndex, depth = position268, tokenIndex268, depth268
			return false
		},
		/* 11 TimeEmitterInterval <- <(<TimeInterval> Action8)> */
		func() bool {
			position275, tokenIndex275, depth275 := position, tokenIndex, depth
			{
				position276 := position
				depth++
				{
					position277 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l275
					}
					depth--
					add(rulePegText, position277)
				}
				if !_rules[ruleAction8]() {
					goto l275
				}
				depth--
				add(ruleTimeEmitterInterval, position276)
			}
			return true
		l275:
			position, tokenIndex, depth = position275, tokenIndex275, depth275
			return false
		},
		/* 12 TupleEmitterInterval <- <(<TuplesInterval> Action9)> */
		func() bool {
			position278, tokenIndex278, depth278 := position, tokenIndex, depth
			{
				position279 := position
				depth++
				{
					position280 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l278
					}
					depth--
					add(rulePegText, position280)
				}
				if !_rules[ruleAction9]() {
					goto l278
				}
				depth--
				add(ruleTupleEmitterInterval, position279)
			}
			return true
		l278:
			position, tokenIndex, depth = position278, tokenIndex278, depth278
			return false
		},
		/* 13 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action10)> */
		func() bool {
			position281, tokenIndex281, depth281 := position, tokenIndex, depth
			{
				position282 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l281
				}
				if !_rules[rulesp]() {
					goto l281
				}
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('I') {
						goto l281
					}
					position++
				}
			l283:
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('N') {
						goto l281
					}
					position++
				}
			l285:
				if !_rules[rulesp]() {
					goto l281
				}
				if !_rules[ruleStream]() {
					goto l281
				}
				if !_rules[ruleAction10]() {
					goto l281
				}
				depth--
				add(ruleTupleEmitterFromInterval, position282)
			}
			return true
		l281:
			position, tokenIndex, depth = position281, tokenIndex281, depth281
			return false
		},
		/* 14 Projections <- <(<(Projection sp (',' sp Projection)*)> Action11)> */
		func() bool {
			position287, tokenIndex287, depth287 := position, tokenIndex, depth
			{
				position288 := position
				depth++
				{
					position289 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l287
					}
					if !_rules[rulesp]() {
						goto l287
					}
				l290:
					{
						position291, tokenIndex291, depth291 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l291
						}
						position++
						if !_rules[rulesp]() {
							goto l291
						}
						if !_rules[ruleProjection]() {
							goto l291
						}
						goto l290
					l291:
						position, tokenIndex, depth = position291, tokenIndex291, depth291
					}
					depth--
					add(rulePegText, position289)
				}
				if !_rules[ruleAction11]() {
					goto l287
				}
				depth--
				add(ruleProjections, position288)
			}
			return true
		l287:
			position, tokenIndex, depth = position287, tokenIndex287, depth287
			return false
		},
		/* 15 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position292, tokenIndex292, depth292 := position, tokenIndex, depth
			{
				position293 := position
				depth++
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l295
					}
					goto l294
				l295:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if !_rules[ruleExpression]() {
						goto l296
					}
					goto l294
				l296:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if !_rules[ruleWildcard]() {
						goto l292
					}
				}
			l294:
				depth--
				add(ruleProjection, position293)
			}
			return true
		l292:
			position, tokenIndex, depth = position292, tokenIndex292, depth292
			return false
		},
		/* 16 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action12)> */
		func() bool {
			position297, tokenIndex297, depth297 := position, tokenIndex, depth
			{
				position298 := position
				depth++
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l300
					}
					goto l299
				l300:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
					if !_rules[ruleWildcard]() {
						goto l297
					}
				}
			l299:
				if !_rules[rulesp]() {
					goto l297
				}
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l302
					}
					position++
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if buffer[position] != rune('A') {
						goto l297
					}
					position++
				}
			l301:
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('S') {
						goto l297
					}
					position++
				}
			l303:
				if !_rules[rulesp]() {
					goto l297
				}
				if !_rules[ruleIdentifier]() {
					goto l297
				}
				if !_rules[ruleAction12]() {
					goto l297
				}
				depth--
				add(ruleAliasExpression, position298)
			}
			return true
		l297:
			position, tokenIndex, depth = position297, tokenIndex297, depth297
			return false
		},
		/* 17 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action13)> */
		func() bool {
			position305, tokenIndex305, depth305 := position, tokenIndex, depth
			{
				position306 := position
				depth++
				{
					position307 := position
					depth++
					{
						position308, tokenIndex308, depth308 := position, tokenIndex, depth
						{
							position310, tokenIndex310, depth310 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l311
							}
							position++
							goto l310
						l311:
							position, tokenIndex, depth = position310, tokenIndex310, depth310
							if buffer[position] != rune('F') {
								goto l308
							}
							position++
						}
					l310:
						{
							position312, tokenIndex312, depth312 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l313
							}
							position++
							goto l312
						l313:
							position, tokenIndex, depth = position312, tokenIndex312, depth312
							if buffer[position] != rune('R') {
								goto l308
							}
							position++
						}
					l312:
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('O') {
								goto l308
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('M') {
								goto l308
							}
							position++
						}
					l316:
						if !_rules[rulesp]() {
							goto l308
						}
						if !_rules[ruleRelations]() {
							goto l308
						}
						if !_rules[rulesp]() {
							goto l308
						}
						goto l309
					l308:
						position, tokenIndex, depth = position308, tokenIndex308, depth308
					}
				l309:
					depth--
					add(rulePegText, position307)
				}
				if !_rules[ruleAction13]() {
					goto l305
				}
				depth--
				add(ruleWindowedFrom, position306)
			}
			return true
		l305:
			position, tokenIndex, depth = position305, tokenIndex305, depth305
			return false
		},
		/* 18 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action14)> */
		func() bool {
			position318, tokenIndex318, depth318 := position, tokenIndex, depth
			{
				position319 := position
				depth++
				{
					position320 := position
					depth++
					{
						position321, tokenIndex321, depth321 := position, tokenIndex, depth
						{
							position323, tokenIndex323, depth323 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l324
							}
							position++
							goto l323
						l324:
							position, tokenIndex, depth = position323, tokenIndex323, depth323
							if buffer[position] != rune('F') {
								goto l321
							}
							position++
						}
					l323:
						{
							position325, tokenIndex325, depth325 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l326
							}
							position++
							goto l325
						l326:
							position, tokenIndex, depth = position325, tokenIndex325, depth325
							if buffer[position] != rune('R') {
								goto l321
							}
							position++
						}
					l325:
						{
							position327, tokenIndex327, depth327 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l328
							}
							position++
							goto l327
						l328:
							position, tokenIndex, depth = position327, tokenIndex327, depth327
							if buffer[position] != rune('O') {
								goto l321
							}
							position++
						}
					l327:
						{
							position329, tokenIndex329, depth329 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l330
							}
							position++
							goto l329
						l330:
							position, tokenIndex, depth = position329, tokenIndex329, depth329
							if buffer[position] != rune('M') {
								goto l321
							}
							position++
						}
					l329:
						if !_rules[rulesp]() {
							goto l321
						}
						if !_rules[ruleDefRelations]() {
							goto l321
						}
						if !_rules[rulesp]() {
							goto l321
						}
						goto l322
					l321:
						position, tokenIndex, depth = position321, tokenIndex321, depth321
					}
				l322:
					depth--
					add(rulePegText, position320)
				}
				if !_rules[ruleAction14]() {
					goto l318
				}
				depth--
				add(ruleDefWindowedFrom, position319)
			}
			return true
		l318:
			position, tokenIndex, depth = position318, tokenIndex318, depth318
			return false
		},
		/* 19 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position331, tokenIndex331, depth331 := position, tokenIndex, depth
			{
				position332 := position
				depth++
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l334
					}
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if !_rules[ruleTuplesInterval]() {
						goto l331
					}
				}
			l333:
				depth--
				add(ruleInterval, position332)
			}
			return true
		l331:
			position, tokenIndex, depth = position331, tokenIndex331, depth331
			return false
		},
		/* 20 TimeInterval <- <(NumericLiteral sp SECONDS Action15)> */
		func() bool {
			position335, tokenIndex335, depth335 := position, tokenIndex, depth
			{
				position336 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l335
				}
				if !_rules[rulesp]() {
					goto l335
				}
				if !_rules[ruleSECONDS]() {
					goto l335
				}
				if !_rules[ruleAction15]() {
					goto l335
				}
				depth--
				add(ruleTimeInterval, position336)
			}
			return true
		l335:
			position, tokenIndex, depth = position335, tokenIndex335, depth335
			return false
		},
		/* 21 TuplesInterval <- <(NumericLiteral sp TUPLES Action16)> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l337
				}
				if !_rules[rulesp]() {
					goto l337
				}
				if !_rules[ruleTUPLES]() {
					goto l337
				}
				if !_rules[ruleAction16]() {
					goto l337
				}
				depth--
				add(ruleTuplesInterval, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 22 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position339, tokenIndex339, depth339 := position, tokenIndex, depth
			{
				position340 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l339
				}
				if !_rules[rulesp]() {
					goto l339
				}
			l341:
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l342
					}
					position++
					if !_rules[rulesp]() {
						goto l342
					}
					if !_rules[ruleRelationLike]() {
						goto l342
					}
					goto l341
				l342:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
				}
				depth--
				add(ruleRelations, position340)
			}
			return true
		l339:
			position, tokenIndex, depth = position339, tokenIndex339, depth339
			return false
		},
		/* 23 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position343, tokenIndex343, depth343 := position, tokenIndex, depth
			{
				position344 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l343
				}
				if !_rules[rulesp]() {
					goto l343
				}
			l345:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l346
					}
					position++
					if !_rules[rulesp]() {
						goto l346
					}
					if !_rules[ruleDefRelationLike]() {
						goto l346
					}
					goto l345
				l346:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
				}
				depth--
				add(ruleDefRelations, position344)
			}
			return true
		l343:
			position, tokenIndex, depth = position343, tokenIndex343, depth343
			return false
		},
		/* 24 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action17)> */
		func() bool {
			position347, tokenIndex347, depth347 := position, tokenIndex, depth
			{
				position348 := position
				depth++
				{
					position349 := position
					depth++
					{
						position350, tokenIndex350, depth350 := position, tokenIndex, depth
						{
							position352, tokenIndex352, depth352 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l353
							}
							position++
							goto l352
						l353:
							position, tokenIndex, depth = position352, tokenIndex352, depth352
							if buffer[position] != rune('W') {
								goto l350
							}
							position++
						}
					l352:
						{
							position354, tokenIndex354, depth354 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l355
							}
							position++
							goto l354
						l355:
							position, tokenIndex, depth = position354, tokenIndex354, depth354
							if buffer[position] != rune('H') {
								goto l350
							}
							position++
						}
					l354:
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l357
							}
							position++
							goto l356
						l357:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
							if buffer[position] != rune('E') {
								goto l350
							}
							position++
						}
					l356:
						{
							position358, tokenIndex358, depth358 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex, depth = position358, tokenIndex358, depth358
							if buffer[position] != rune('R') {
								goto l350
							}
							position++
						}
					l358:
						{
							position360, tokenIndex360, depth360 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l361
							}
							position++
							goto l360
						l361:
							position, tokenIndex, depth = position360, tokenIndex360, depth360
							if buffer[position] != rune('E') {
								goto l350
							}
							position++
						}
					l360:
						if !_rules[rulesp]() {
							goto l350
						}
						if !_rules[ruleExpression]() {
							goto l350
						}
						goto l351
					l350:
						position, tokenIndex, depth = position350, tokenIndex350, depth350
					}
				l351:
					depth--
					add(rulePegText, position349)
				}
				if !_rules[ruleAction17]() {
					goto l347
				}
				depth--
				add(ruleFilter, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 25 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action18)> */
		func() bool {
			position362, tokenIndex362, depth362 := position, tokenIndex, depth
			{
				position363 := position
				depth++
				{
					position364 := position
					depth++
					{
						position365, tokenIndex365, depth365 := position, tokenIndex, depth
						{
							position367, tokenIndex367, depth367 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l368
							}
							position++
							goto l367
						l368:
							position, tokenIndex, depth = position367, tokenIndex367, depth367
							if buffer[position] != rune('G') {
								goto l365
							}
							position++
						}
					l367:
						{
							position369, tokenIndex369, depth369 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l370
							}
							position++
							goto l369
						l370:
							position, tokenIndex, depth = position369, tokenIndex369, depth369
							if buffer[position] != rune('R') {
								goto l365
							}
							position++
						}
					l369:
						{
							position371, tokenIndex371, depth371 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l372
							}
							position++
							goto l371
						l372:
							position, tokenIndex, depth = position371, tokenIndex371, depth371
							if buffer[position] != rune('O') {
								goto l365
							}
							position++
						}
					l371:
						{
							position373, tokenIndex373, depth373 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l374
							}
							position++
							goto l373
						l374:
							position, tokenIndex, depth = position373, tokenIndex373, depth373
							if buffer[position] != rune('U') {
								goto l365
							}
							position++
						}
					l373:
						{
							position375, tokenIndex375, depth375 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l376
							}
							position++
							goto l375
						l376:
							position, tokenIndex, depth = position375, tokenIndex375, depth375
							if buffer[position] != rune('P') {
								goto l365
							}
							position++
						}
					l375:
						if !_rules[rulesp]() {
							goto l365
						}
						{
							position377, tokenIndex377, depth377 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l378
							}
							position++
							goto l377
						l378:
							position, tokenIndex, depth = position377, tokenIndex377, depth377
							if buffer[position] != rune('B') {
								goto l365
							}
							position++
						}
					l377:
						{
							position379, tokenIndex379, depth379 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l380
							}
							position++
							goto l379
						l380:
							position, tokenIndex, depth = position379, tokenIndex379, depth379
							if buffer[position] != rune('Y') {
								goto l365
							}
							position++
						}
					l379:
						if !_rules[rulesp]() {
							goto l365
						}
						if !_rules[ruleGroupList]() {
							goto l365
						}
						goto l366
					l365:
						position, tokenIndex, depth = position365, tokenIndex365, depth365
					}
				l366:
					depth--
					add(rulePegText, position364)
				}
				if !_rules[ruleAction18]() {
					goto l362
				}
				depth--
				add(ruleGrouping, position363)
			}
			return true
		l362:
			position, tokenIndex, depth = position362, tokenIndex362, depth362
			return false
		},
		/* 26 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position381, tokenIndex381, depth381 := position, tokenIndex, depth
			{
				position382 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l381
				}
				if !_rules[rulesp]() {
					goto l381
				}
			l383:
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l384
					}
					position++
					if !_rules[rulesp]() {
						goto l384
					}
					if !_rules[ruleExpression]() {
						goto l384
					}
					goto l383
				l384:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
				}
				depth--
				add(ruleGroupList, position382)
			}
			return true
		l381:
			position, tokenIndex, depth = position381, tokenIndex381, depth381
			return false
		},
		/* 27 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action19)> */
		func() bool {
			position385, tokenIndex385, depth385 := position, tokenIndex, depth
			{
				position386 := position
				depth++
				{
					position387 := position
					depth++
					{
						position388, tokenIndex388, depth388 := position, tokenIndex, depth
						{
							position390, tokenIndex390, depth390 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l391
							}
							position++
							goto l390
						l391:
							position, tokenIndex, depth = position390, tokenIndex390, depth390
							if buffer[position] != rune('H') {
								goto l388
							}
							position++
						}
					l390:
						{
							position392, tokenIndex392, depth392 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l393
							}
							position++
							goto l392
						l393:
							position, tokenIndex, depth = position392, tokenIndex392, depth392
							if buffer[position] != rune('A') {
								goto l388
							}
							position++
						}
					l392:
						{
							position394, tokenIndex394, depth394 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l395
							}
							position++
							goto l394
						l395:
							position, tokenIndex, depth = position394, tokenIndex394, depth394
							if buffer[position] != rune('V') {
								goto l388
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
								goto l388
							}
							position++
						}
					l396:
						{
							position398, tokenIndex398, depth398 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l399
							}
							position++
							goto l398
						l399:
							position, tokenIndex, depth = position398, tokenIndex398, depth398
							if buffer[position] != rune('N') {
								goto l388
							}
							position++
						}
					l398:
						{
							position400, tokenIndex400, depth400 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l401
							}
							position++
							goto l400
						l401:
							position, tokenIndex, depth = position400, tokenIndex400, depth400
							if buffer[position] != rune('G') {
								goto l388
							}
							position++
						}
					l400:
						if !_rules[rulesp]() {
							goto l388
						}
						if !_rules[ruleExpression]() {
							goto l388
						}
						goto l389
					l388:
						position, tokenIndex, depth = position388, tokenIndex388, depth388
					}
				l389:
					depth--
					add(rulePegText, position387)
				}
				if !_rules[ruleAction19]() {
					goto l385
				}
				depth--
				add(ruleHaving, position386)
			}
			return true
		l385:
			position, tokenIndex, depth = position385, tokenIndex385, depth385
			return false
		},
		/* 28 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action20))> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l405
					}
					goto l404
				l405:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					if !_rules[ruleStreamWindow]() {
						goto l402
					}
					if !_rules[ruleAction20]() {
						goto l402
					}
				}
			l404:
				depth--
				add(ruleRelationLike, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 29 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action21))> */
		func() bool {
			position406, tokenIndex406, depth406 := position, tokenIndex, depth
			{
				position407 := position
				depth++
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l409
					}
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if !_rules[ruleDefStreamWindow]() {
						goto l406
					}
					if !_rules[ruleAction21]() {
						goto l406
					}
				}
			l408:
				depth--
				add(ruleDefRelationLike, position407)
			}
			return true
		l406:
			position, tokenIndex, depth = position406, tokenIndex406, depth406
			return false
		},
		/* 30 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action22)> */
		func() bool {
			position410, tokenIndex410, depth410 := position, tokenIndex, depth
			{
				position411 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l410
				}
				if !_rules[rulesp]() {
					goto l410
				}
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l413
					}
					position++
					goto l412
				l413:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
					if buffer[position] != rune('A') {
						goto l410
					}
					position++
				}
			l412:
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('S') {
						goto l410
					}
					position++
				}
			l414:
				if !_rules[rulesp]() {
					goto l410
				}
				if !_rules[ruleIdentifier]() {
					goto l410
				}
				if !_rules[ruleAction22]() {
					goto l410
				}
				depth--
				add(ruleAliasedStreamWindow, position411)
			}
			return true
		l410:
			position, tokenIndex, depth = position410, tokenIndex410, depth410
			return false
		},
		/* 31 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action23)> */
		func() bool {
			position416, tokenIndex416, depth416 := position, tokenIndex, depth
			{
				position417 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l416
				}
				if !_rules[rulesp]() {
					goto l416
				}
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
					if buffer[position] != rune('A') {
						goto l416
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('S') {
						goto l416
					}
					position++
				}
			l420:
				if !_rules[rulesp]() {
					goto l416
				}
				if !_rules[ruleIdentifier]() {
					goto l416
				}
				if !_rules[ruleAction23]() {
					goto l416
				}
				depth--
				add(ruleDefAliasedStreamWindow, position417)
			}
			return true
		l416:
			position, tokenIndex, depth = position416, tokenIndex416, depth416
			return false
		},
		/* 32 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action24)> */
		func() bool {
			position422, tokenIndex422, depth422 := position, tokenIndex, depth
			{
				position423 := position
				depth++
				if !_rules[ruleStream]() {
					goto l422
				}
				if !_rules[rulesp]() {
					goto l422
				}
				if buffer[position] != rune('[') {
					goto l422
				}
				position++
				if !_rules[rulesp]() {
					goto l422
				}
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l425
					}
					position++
					goto l424
				l425:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
					if buffer[position] != rune('R') {
						goto l422
					}
					position++
				}
			l424:
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l427
					}
					position++
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if buffer[position] != rune('A') {
						goto l422
					}
					position++
				}
			l426:
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l429
					}
					position++
					goto l428
				l429:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
					if buffer[position] != rune('N') {
						goto l422
					}
					position++
				}
			l428:
				{
					position430, tokenIndex430, depth430 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l431
					}
					position++
					goto l430
				l431:
					position, tokenIndex, depth = position430, tokenIndex430, depth430
					if buffer[position] != rune('G') {
						goto l422
					}
					position++
				}
			l430:
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('E') {
						goto l422
					}
					position++
				}
			l432:
				if !_rules[rulesp]() {
					goto l422
				}
				if !_rules[ruleInterval]() {
					goto l422
				}
				if !_rules[rulesp]() {
					goto l422
				}
				if buffer[position] != rune(']') {
					goto l422
				}
				position++
				if !_rules[ruleAction24]() {
					goto l422
				}
				depth--
				add(ruleStreamWindow, position423)
			}
			return true
		l422:
			position, tokenIndex, depth = position422, tokenIndex422, depth422
			return false
		},
		/* 33 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action25)> */
		func() bool {
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
			{
				position435 := position
				depth++
				if !_rules[ruleStream]() {
					goto l434
				}
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l436
					}
					if buffer[position] != rune('[') {
						goto l436
					}
					position++
					if !_rules[rulesp]() {
						goto l436
					}
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l439
						}
						position++
						goto l438
					l439:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
						if buffer[position] != rune('R') {
							goto l436
						}
						position++
					}
				l438:
					{
						position440, tokenIndex440, depth440 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l441
						}
						position++
						goto l440
					l441:
						position, tokenIndex, depth = position440, tokenIndex440, depth440
						if buffer[position] != rune('A') {
							goto l436
						}
						position++
					}
				l440:
					{
						position442, tokenIndex442, depth442 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l443
						}
						position++
						goto l442
					l443:
						position, tokenIndex, depth = position442, tokenIndex442, depth442
						if buffer[position] != rune('N') {
							goto l436
						}
						position++
					}
				l442:
					{
						position444, tokenIndex444, depth444 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l445
						}
						position++
						goto l444
					l445:
						position, tokenIndex, depth = position444, tokenIndex444, depth444
						if buffer[position] != rune('G') {
							goto l436
						}
						position++
					}
				l444:
					{
						position446, tokenIndex446, depth446 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l447
						}
						position++
						goto l446
					l447:
						position, tokenIndex, depth = position446, tokenIndex446, depth446
						if buffer[position] != rune('E') {
							goto l436
						}
						position++
					}
				l446:
					if !_rules[rulesp]() {
						goto l436
					}
					if !_rules[ruleInterval]() {
						goto l436
					}
					if !_rules[rulesp]() {
						goto l436
					}
					if buffer[position] != rune(']') {
						goto l436
					}
					position++
					goto l437
				l436:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
				}
			l437:
				if !_rules[ruleAction25]() {
					goto l434
				}
				depth--
				add(ruleDefStreamWindow, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
			return false
		},
		/* 34 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action26)> */
		func() bool {
			position448, tokenIndex448, depth448 := position, tokenIndex, depth
			{
				position449 := position
				depth++
				{
					position450 := position
					depth++
					{
						position451, tokenIndex451, depth451 := position, tokenIndex, depth
						{
							position453, tokenIndex453, depth453 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l454
							}
							position++
							goto l453
						l454:
							position, tokenIndex, depth = position453, tokenIndex453, depth453
							if buffer[position] != rune('W') {
								goto l451
							}
							position++
						}
					l453:
						{
							position455, tokenIndex455, depth455 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l456
							}
							position++
							goto l455
						l456:
							position, tokenIndex, depth = position455, tokenIndex455, depth455
							if buffer[position] != rune('I') {
								goto l451
							}
							position++
						}
					l455:
						{
							position457, tokenIndex457, depth457 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l458
							}
							position++
							goto l457
						l458:
							position, tokenIndex, depth = position457, tokenIndex457, depth457
							if buffer[position] != rune('T') {
								goto l451
							}
							position++
						}
					l457:
						{
							position459, tokenIndex459, depth459 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l460
							}
							position++
							goto l459
						l460:
							position, tokenIndex, depth = position459, tokenIndex459, depth459
							if buffer[position] != rune('H') {
								goto l451
							}
							position++
						}
					l459:
						if !_rules[rulesp]() {
							goto l451
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l451
						}
						if !_rules[rulesp]() {
							goto l451
						}
					l461:
						{
							position462, tokenIndex462, depth462 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l462
							}
							position++
							if !_rules[rulesp]() {
								goto l462
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l462
							}
							goto l461
						l462:
							position, tokenIndex, depth = position462, tokenIndex462, depth462
						}
						goto l452
					l451:
						position, tokenIndex, depth = position451, tokenIndex451, depth451
					}
				l452:
					depth--
					add(rulePegText, position450)
				}
				if !_rules[ruleAction26]() {
					goto l448
				}
				depth--
				add(ruleSourceSinkSpecs, position449)
			}
			return true
		l448:
			position, tokenIndex, depth = position448, tokenIndex448, depth448
			return false
		},
		/* 35 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action27)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l463
				}
				if buffer[position] != rune('=') {
					goto l463
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l463
				}
				if !_rules[ruleAction27]() {
					goto l463
				}
				depth--
				add(ruleSourceSinkParam, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
			return false
		},
		/* 36 Expression <- <orExpr> */
		func() bool {
			position465, tokenIndex465, depth465 := position, tokenIndex, depth
			{
				position466 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l465
				}
				depth--
				add(ruleExpression, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 37 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action28)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				{
					position469 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l467
					}
					if !_rules[rulesp]() {
						goto l467
					}
					{
						position470, tokenIndex470, depth470 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l470
						}
						if !_rules[rulesp]() {
							goto l470
						}
						if !_rules[ruleandExpr]() {
							goto l470
						}
						goto l471
					l470:
						position, tokenIndex, depth = position470, tokenIndex470, depth470
					}
				l471:
					depth--
					add(rulePegText, position469)
				}
				if !_rules[ruleAction28]() {
					goto l467
				}
				depth--
				add(ruleorExpr, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 38 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action29)> */
		func() bool {
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				{
					position474 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l472
					}
					if !_rules[rulesp]() {
						goto l472
					}
					{
						position475, tokenIndex475, depth475 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l475
						}
						if !_rules[rulesp]() {
							goto l475
						}
						if !_rules[rulecomparisonExpr]() {
							goto l475
						}
						goto l476
					l475:
						position, tokenIndex, depth = position475, tokenIndex475, depth475
					}
				l476:
					depth--
					add(rulePegText, position474)
				}
				if !_rules[ruleAction29]() {
					goto l472
				}
				depth--
				add(ruleandExpr, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 39 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action30)> */
		func() bool {
			position477, tokenIndex477, depth477 := position, tokenIndex, depth
			{
				position478 := position
				depth++
				{
					position479 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l477
					}
					if !_rules[rulesp]() {
						goto l477
					}
					{
						position480, tokenIndex480, depth480 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l480
						}
						if !_rules[rulesp]() {
							goto l480
						}
						if !_rules[ruletermExpr]() {
							goto l480
						}
						goto l481
					l480:
						position, tokenIndex, depth = position480, tokenIndex480, depth480
					}
				l481:
					depth--
					add(rulePegText, position479)
				}
				if !_rules[ruleAction30]() {
					goto l477
				}
				depth--
				add(rulecomparisonExpr, position478)
			}
			return true
		l477:
			position, tokenIndex, depth = position477, tokenIndex477, depth477
			return false
		},
		/* 40 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action31)> */
		func() bool {
			position482, tokenIndex482, depth482 := position, tokenIndex, depth
			{
				position483 := position
				depth++
				{
					position484 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l482
					}
					if !_rules[rulesp]() {
						goto l482
					}
					{
						position485, tokenIndex485, depth485 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l485
						}
						if !_rules[rulesp]() {
							goto l485
						}
						if !_rules[ruleproductExpr]() {
							goto l485
						}
						goto l486
					l485:
						position, tokenIndex, depth = position485, tokenIndex485, depth485
					}
				l486:
					depth--
					add(rulePegText, position484)
				}
				if !_rules[ruleAction31]() {
					goto l482
				}
				depth--
				add(ruletermExpr, position483)
			}
			return true
		l482:
			position, tokenIndex, depth = position482, tokenIndex482, depth482
			return false
		},
		/* 41 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action32)> */
		func() bool {
			position487, tokenIndex487, depth487 := position, tokenIndex, depth
			{
				position488 := position
				depth++
				{
					position489 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l487
					}
					if !_rules[rulesp]() {
						goto l487
					}
					{
						position490, tokenIndex490, depth490 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l490
						}
						if !_rules[rulesp]() {
							goto l490
						}
						if !_rules[rulebaseExpr]() {
							goto l490
						}
						goto l491
					l490:
						position, tokenIndex, depth = position490, tokenIndex490, depth490
					}
				l491:
					depth--
					add(rulePegText, position489)
				}
				if !_rules[ruleAction32]() {
					goto l487
				}
				depth--
				add(ruleproductExpr, position488)
			}
			return true
		l487:
			position, tokenIndex, depth = position487, tokenIndex487, depth487
			return false
		},
		/* 42 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position492, tokenIndex492, depth492 := position, tokenIndex, depth
			{
				position493 := position
				depth++
				{
					position494, tokenIndex494, depth494 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l495
					}
					position++
					if !_rules[rulesp]() {
						goto l495
					}
					if !_rules[ruleExpression]() {
						goto l495
					}
					if !_rules[rulesp]() {
						goto l495
					}
					if buffer[position] != rune(')') {
						goto l495
					}
					position++
					goto l494
				l495:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if !_rules[ruleBooleanLiteral]() {
						goto l496
					}
					goto l494
				l496:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if !_rules[ruleFuncApp]() {
						goto l497
					}
					goto l494
				l497:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if !_rules[ruleRowValue]() {
						goto l498
					}
					goto l494
				l498:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if !_rules[ruleLiteral]() {
						goto l492
					}
				}
			l494:
				depth--
				add(rulebaseExpr, position493)
			}
			return true
		l492:
			position, tokenIndex, depth = position492, tokenIndex492, depth492
			return false
		},
		/* 43 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action33)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l499
				}
				if !_rules[rulesp]() {
					goto l499
				}
				if buffer[position] != rune('(') {
					goto l499
				}
				position++
				if !_rules[rulesp]() {
					goto l499
				}
				if !_rules[ruleFuncParams]() {
					goto l499
				}
				if !_rules[rulesp]() {
					goto l499
				}
				if buffer[position] != rune(')') {
					goto l499
				}
				position++
				if !_rules[ruleAction33]() {
					goto l499
				}
				depth--
				add(ruleFuncApp, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 44 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action34)> */
		func() bool {
			position501, tokenIndex501, depth501 := position, tokenIndex, depth
			{
				position502 := position
				depth++
				{
					position503 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l501
					}
					if !_rules[rulesp]() {
						goto l501
					}
				l504:
					{
						position505, tokenIndex505, depth505 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l505
						}
						position++
						if !_rules[rulesp]() {
							goto l505
						}
						if !_rules[ruleExpression]() {
							goto l505
						}
						goto l504
					l505:
						position, tokenIndex, depth = position505, tokenIndex505, depth505
					}
					depth--
					add(rulePegText, position503)
				}
				if !_rules[ruleAction34]() {
					goto l501
				}
				depth--
				add(ruleFuncParams, position502)
			}
			return true
		l501:
			position, tokenIndex, depth = position501, tokenIndex501, depth501
			return false
		},
		/* 45 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				{
					position508, tokenIndex508, depth508 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l509
					}
					goto l508
				l509:
					position, tokenIndex, depth = position508, tokenIndex508, depth508
					if !_rules[ruleNumericLiteral]() {
						goto l510
					}
					goto l508
				l510:
					position, tokenIndex, depth = position508, tokenIndex508, depth508
					if !_rules[ruleStringLiteral]() {
						goto l506
					}
				}
			l508:
				depth--
				add(ruleLiteral, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 46 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position511, tokenIndex511, depth511 := position, tokenIndex, depth
			{
				position512 := position
				depth++
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l514
					}
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleNotEqual]() {
						goto l515
					}
					goto l513
				l515:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleLessOrEqual]() {
						goto l516
					}
					goto l513
				l516:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleLess]() {
						goto l517
					}
					goto l513
				l517:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleGreaterOrEqual]() {
						goto l518
					}
					goto l513
				l518:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleGreater]() {
						goto l519
					}
					goto l513
				l519:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleNotEqual]() {
						goto l511
					}
				}
			l513:
				depth--
				add(ruleComparisonOp, position512)
			}
			return true
		l511:
			position, tokenIndex, depth = position511, tokenIndex511, depth511
			return false
		},
		/* 47 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position520, tokenIndex520, depth520 := position, tokenIndex, depth
			{
				position521 := position
				depth++
				{
					position522, tokenIndex522, depth522 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l523
					}
					goto l522
				l523:
					position, tokenIndex, depth = position522, tokenIndex522, depth522
					if !_rules[ruleMinus]() {
						goto l520
					}
				}
			l522:
				depth--
				add(rulePlusMinusOp, position521)
			}
			return true
		l520:
			position, tokenIndex, depth = position520, tokenIndex520, depth520
			return false
		},
		/* 48 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				{
					position526, tokenIndex526, depth526 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l527
					}
					goto l526
				l527:
					position, tokenIndex, depth = position526, tokenIndex526, depth526
					if !_rules[ruleDivide]() {
						goto l528
					}
					goto l526
				l528:
					position, tokenIndex, depth = position526, tokenIndex526, depth526
					if !_rules[ruleModulo]() {
						goto l524
					}
				}
			l526:
				depth--
				add(ruleMultDivOp, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 49 Stream <- <(<ident> Action35)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531 := position
					depth++
					if !_rules[ruleident]() {
						goto l529
					}
					depth--
					add(rulePegText, position531)
				}
				if !_rules[ruleAction35]() {
					goto l529
				}
				depth--
				add(ruleStream, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 50 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action36)> */
		func() bool {
			position532, tokenIndex532, depth532 := position, tokenIndex, depth
			{
				position533 := position
				depth++
				{
					position534 := position
					depth++
					{
						position535, tokenIndex535, depth535 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l535
						}
						if buffer[position] != rune(':') {
							goto l535
						}
						position++
						goto l536
					l535:
						position, tokenIndex, depth = position535, tokenIndex535, depth535
					}
				l536:
					{
						position537, tokenIndex537, depth537 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l538
						}
						position++
						goto l537
					l538:
						position, tokenIndex, depth = position537, tokenIndex537, depth537
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l532
						}
						position++
					}
				l537:
				l539:
					{
						position540, tokenIndex540, depth540 := position, tokenIndex, depth
						{
							position541, tokenIndex541, depth541 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l542
							}
							position++
							goto l541
						l542:
							position, tokenIndex, depth = position541, tokenIndex541, depth541
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l543
							}
							position++
							goto l541
						l543:
							position, tokenIndex, depth = position541, tokenIndex541, depth541
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l544
							}
							position++
							goto l541
						l544:
							position, tokenIndex, depth = position541, tokenIndex541, depth541
							if buffer[position] != rune('_') {
								goto l545
							}
							position++
							goto l541
						l545:
							position, tokenIndex, depth = position541, tokenIndex541, depth541
							if buffer[position] != rune('.') {
								goto l540
							}
							position++
						}
					l541:
						goto l539
					l540:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
					}
					depth--
					add(rulePegText, position534)
				}
				if !_rules[ruleAction36]() {
					goto l532
				}
				depth--
				add(ruleRowValue, position533)
			}
			return true
		l532:
			position, tokenIndex, depth = position532, tokenIndex532, depth532
			return false
		},
		/* 51 NumericLiteral <- <(<('-'? [0-9]+)> Action37)> */
		func() bool {
			position546, tokenIndex546, depth546 := position, tokenIndex, depth
			{
				position547 := position
				depth++
				{
					position548 := position
					depth++
					{
						position549, tokenIndex549, depth549 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l549
						}
						position++
						goto l550
					l549:
						position, tokenIndex, depth = position549, tokenIndex549, depth549
					}
				l550:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l546
					}
					position++
				l551:
					{
						position552, tokenIndex552, depth552 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l552
						}
						position++
						goto l551
					l552:
						position, tokenIndex, depth = position552, tokenIndex552, depth552
					}
					depth--
					add(rulePegText, position548)
				}
				if !_rules[ruleAction37]() {
					goto l546
				}
				depth--
				add(ruleNumericLiteral, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 52 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action38)> */
		func() bool {
			position553, tokenIndex553, depth553 := position, tokenIndex, depth
			{
				position554 := position
				depth++
				{
					position555 := position
					depth++
					{
						position556, tokenIndex556, depth556 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l556
						}
						position++
						goto l557
					l556:
						position, tokenIndex, depth = position556, tokenIndex556, depth556
					}
				l557:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l553
					}
					position++
				l558:
					{
						position559, tokenIndex559, depth559 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l559
						}
						position++
						goto l558
					l559:
						position, tokenIndex, depth = position559, tokenIndex559, depth559
					}
					if buffer[position] != rune('.') {
						goto l553
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l553
					}
					position++
				l560:
					{
						position561, tokenIndex561, depth561 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l561
						}
						position++
						goto l560
					l561:
						position, tokenIndex, depth = position561, tokenIndex561, depth561
					}
					depth--
					add(rulePegText, position555)
				}
				if !_rules[ruleAction38]() {
					goto l553
				}
				depth--
				add(ruleFloatLiteral, position554)
			}
			return true
		l553:
			position, tokenIndex, depth = position553, tokenIndex553, depth553
			return false
		},
		/* 53 Function <- <(<ident> Action39)> */
		func() bool {
			position562, tokenIndex562, depth562 := position, tokenIndex, depth
			{
				position563 := position
				depth++
				{
					position564 := position
					depth++
					if !_rules[ruleident]() {
						goto l562
					}
					depth--
					add(rulePegText, position564)
				}
				if !_rules[ruleAction39]() {
					goto l562
				}
				depth--
				add(ruleFunction, position563)
			}
			return true
		l562:
			position, tokenIndex, depth = position562, tokenIndex562, depth562
			return false
		},
		/* 54 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567, tokenIndex567, depth567 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l568
					}
					goto l567
				l568:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleFALSE]() {
						goto l565
					}
				}
			l567:
				depth--
				add(ruleBooleanLiteral, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 55 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action40)> */
		func() bool {
			position569, tokenIndex569, depth569 := position, tokenIndex, depth
			{
				position570 := position
				depth++
				{
					position571 := position
					depth++
					{
						position572, tokenIndex572, depth572 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l573
						}
						position++
						goto l572
					l573:
						position, tokenIndex, depth = position572, tokenIndex572, depth572
						if buffer[position] != rune('T') {
							goto l569
						}
						position++
					}
				l572:
					{
						position574, tokenIndex574, depth574 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l575
						}
						position++
						goto l574
					l575:
						position, tokenIndex, depth = position574, tokenIndex574, depth574
						if buffer[position] != rune('R') {
							goto l569
						}
						position++
					}
				l574:
					{
						position576, tokenIndex576, depth576 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l577
						}
						position++
						goto l576
					l577:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						if buffer[position] != rune('U') {
							goto l569
						}
						position++
					}
				l576:
					{
						position578, tokenIndex578, depth578 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l579
						}
						position++
						goto l578
					l579:
						position, tokenIndex, depth = position578, tokenIndex578, depth578
						if buffer[position] != rune('E') {
							goto l569
						}
						position++
					}
				l578:
					depth--
					add(rulePegText, position571)
				}
				if !_rules[ruleAction40]() {
					goto l569
				}
				depth--
				add(ruleTRUE, position570)
			}
			return true
		l569:
			position, tokenIndex, depth = position569, tokenIndex569, depth569
			return false
		},
		/* 56 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action41)> */
		func() bool {
			position580, tokenIndex580, depth580 := position, tokenIndex, depth
			{
				position581 := position
				depth++
				{
					position582 := position
					depth++
					{
						position583, tokenIndex583, depth583 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l584
						}
						position++
						goto l583
					l584:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
						if buffer[position] != rune('F') {
							goto l580
						}
						position++
					}
				l583:
					{
						position585, tokenIndex585, depth585 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l586
						}
						position++
						goto l585
					l586:
						position, tokenIndex, depth = position585, tokenIndex585, depth585
						if buffer[position] != rune('A') {
							goto l580
						}
						position++
					}
				l585:
					{
						position587, tokenIndex587, depth587 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l588
						}
						position++
						goto l587
					l588:
						position, tokenIndex, depth = position587, tokenIndex587, depth587
						if buffer[position] != rune('L') {
							goto l580
						}
						position++
					}
				l587:
					{
						position589, tokenIndex589, depth589 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l590
						}
						position++
						goto l589
					l590:
						position, tokenIndex, depth = position589, tokenIndex589, depth589
						if buffer[position] != rune('S') {
							goto l580
						}
						position++
					}
				l589:
					{
						position591, tokenIndex591, depth591 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l592
						}
						position++
						goto l591
					l592:
						position, tokenIndex, depth = position591, tokenIndex591, depth591
						if buffer[position] != rune('E') {
							goto l580
						}
						position++
					}
				l591:
					depth--
					add(rulePegText, position582)
				}
				if !_rules[ruleAction41]() {
					goto l580
				}
				depth--
				add(ruleFALSE, position581)
			}
			return true
		l580:
			position, tokenIndex, depth = position580, tokenIndex580, depth580
			return false
		},
		/* 57 Wildcard <- <(<'*'> Action42)> */
		func() bool {
			position593, tokenIndex593, depth593 := position, tokenIndex, depth
			{
				position594 := position
				depth++
				{
					position595 := position
					depth++
					if buffer[position] != rune('*') {
						goto l593
					}
					position++
					depth--
					add(rulePegText, position595)
				}
				if !_rules[ruleAction42]() {
					goto l593
				}
				depth--
				add(ruleWildcard, position594)
			}
			return true
		l593:
			position, tokenIndex, depth = position593, tokenIndex593, depth593
			return false
		},
		/* 58 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action43)> */
		func() bool {
			position596, tokenIndex596, depth596 := position, tokenIndex, depth
			{
				position597 := position
				depth++
				{
					position598 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l596
					}
					position++
				l599:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						{
							position601, tokenIndex601, depth601 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l602
							}
							position++
							if buffer[position] != rune('\'') {
								goto l602
							}
							position++
							goto l601
						l602:
							position, tokenIndex, depth = position601, tokenIndex601, depth601
							{
								position603, tokenIndex603, depth603 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l603
								}
								position++
								goto l600
							l603:
								position, tokenIndex, depth = position603, tokenIndex603, depth603
							}
							if !matchDot() {
								goto l600
							}
						}
					l601:
						goto l599
					l600:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
					}
					if buffer[position] != rune('\'') {
						goto l596
					}
					position++
					depth--
					add(rulePegText, position598)
				}
				if !_rules[ruleAction43]() {
					goto l596
				}
				depth--
				add(ruleStringLiteral, position597)
			}
			return true
		l596:
			position, tokenIndex, depth = position596, tokenIndex596, depth596
			return false
		},
		/* 59 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action44)> */
		func() bool {
			position604, tokenIndex604, depth604 := position, tokenIndex, depth
			{
				position605 := position
				depth++
				{
					position606 := position
					depth++
					{
						position607, tokenIndex607, depth607 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l608
						}
						position++
						goto l607
					l608:
						position, tokenIndex, depth = position607, tokenIndex607, depth607
						if buffer[position] != rune('I') {
							goto l604
						}
						position++
					}
				l607:
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('S') {
							goto l604
						}
						position++
					}
				l609:
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
							goto l604
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
							goto l604
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('E') {
							goto l604
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('A') {
							goto l604
						}
						position++
					}
				l617:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('M') {
							goto l604
						}
						position++
					}
				l619:
					depth--
					add(rulePegText, position606)
				}
				if !_rules[ruleAction44]() {
					goto l604
				}
				depth--
				add(ruleISTREAM, position605)
			}
			return true
		l604:
			position, tokenIndex, depth = position604, tokenIndex604, depth604
			return false
		},
		/* 60 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action45)> */
		func() bool {
			position621, tokenIndex621, depth621 := position, tokenIndex, depth
			{
				position622 := position
				depth++
				{
					position623 := position
					depth++
					{
						position624, tokenIndex624, depth624 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('D') {
							goto l621
						}
						position++
					}
				l624:
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('S') {
							goto l621
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('T') {
							goto l621
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('R') {
							goto l621
						}
						position++
					}
				l630:
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l633
						}
						position++
						goto l632
					l633:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
						if buffer[position] != rune('E') {
							goto l621
						}
						position++
					}
				l632:
					{
						position634, tokenIndex634, depth634 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l635
						}
						position++
						goto l634
					l635:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
						if buffer[position] != rune('A') {
							goto l621
						}
						position++
					}
				l634:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l637
						}
						position++
						goto l636
					l637:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
						if buffer[position] != rune('M') {
							goto l621
						}
						position++
					}
				l636:
					depth--
					add(rulePegText, position623)
				}
				if !_rules[ruleAction45]() {
					goto l621
				}
				depth--
				add(ruleDSTREAM, position622)
			}
			return true
		l621:
			position, tokenIndex, depth = position621, tokenIndex621, depth621
			return false
		},
		/* 61 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action46)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					{
						position641, tokenIndex641, depth641 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l642
						}
						position++
						goto l641
					l642:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
						if buffer[position] != rune('R') {
							goto l638
						}
						position++
					}
				l641:
					{
						position643, tokenIndex643, depth643 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l644
						}
						position++
						goto l643
					l644:
						position, tokenIndex, depth = position643, tokenIndex643, depth643
						if buffer[position] != rune('S') {
							goto l638
						}
						position++
					}
				l643:
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l646
						}
						position++
						goto l645
					l646:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
						if buffer[position] != rune('T') {
							goto l638
						}
						position++
					}
				l645:
					{
						position647, tokenIndex647, depth647 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l648
						}
						position++
						goto l647
					l648:
						position, tokenIndex, depth = position647, tokenIndex647, depth647
						if buffer[position] != rune('R') {
							goto l638
						}
						position++
					}
				l647:
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l650
						}
						position++
						goto l649
					l650:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						if buffer[position] != rune('E') {
							goto l638
						}
						position++
					}
				l649:
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
							goto l638
						}
						position++
					}
				l651:
					{
						position653, tokenIndex653, depth653 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l654
						}
						position++
						goto l653
					l654:
						position, tokenIndex, depth = position653, tokenIndex653, depth653
						if buffer[position] != rune('M') {
							goto l638
						}
						position++
					}
				l653:
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction46]() {
					goto l638
				}
				depth--
				add(ruleRSTREAM, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 62 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action47)> */
		func() bool {
			position655, tokenIndex655, depth655 := position, tokenIndex, depth
			{
				position656 := position
				depth++
				{
					position657 := position
					depth++
					{
						position658, tokenIndex658, depth658 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l659
						}
						position++
						goto l658
					l659:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
						if buffer[position] != rune('T') {
							goto l655
						}
						position++
					}
				l658:
					{
						position660, tokenIndex660, depth660 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l661
						}
						position++
						goto l660
					l661:
						position, tokenIndex, depth = position660, tokenIndex660, depth660
						if buffer[position] != rune('U') {
							goto l655
						}
						position++
					}
				l660:
					{
						position662, tokenIndex662, depth662 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l663
						}
						position++
						goto l662
					l663:
						position, tokenIndex, depth = position662, tokenIndex662, depth662
						if buffer[position] != rune('P') {
							goto l655
						}
						position++
					}
				l662:
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l665
						}
						position++
						goto l664
					l665:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if buffer[position] != rune('L') {
							goto l655
						}
						position++
					}
				l664:
					{
						position666, tokenIndex666, depth666 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l667
						}
						position++
						goto l666
					l667:
						position, tokenIndex, depth = position666, tokenIndex666, depth666
						if buffer[position] != rune('E') {
							goto l655
						}
						position++
					}
				l666:
					{
						position668, tokenIndex668, depth668 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l669
						}
						position++
						goto l668
					l669:
						position, tokenIndex, depth = position668, tokenIndex668, depth668
						if buffer[position] != rune('S') {
							goto l655
						}
						position++
					}
				l668:
					depth--
					add(rulePegText, position657)
				}
				if !_rules[ruleAction47]() {
					goto l655
				}
				depth--
				add(ruleTUPLES, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 63 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action48)> */
		func() bool {
			position670, tokenIndex670, depth670 := position, tokenIndex, depth
			{
				position671 := position
				depth++
				{
					position672 := position
					depth++
					{
						position673, tokenIndex673, depth673 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l674
						}
						position++
						goto l673
					l674:
						position, tokenIndex, depth = position673, tokenIndex673, depth673
						if buffer[position] != rune('S') {
							goto l670
						}
						position++
					}
				l673:
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l676
						}
						position++
						goto l675
					l676:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('E') {
							goto l670
						}
						position++
					}
				l675:
					{
						position677, tokenIndex677, depth677 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l678
						}
						position++
						goto l677
					l678:
						position, tokenIndex, depth = position677, tokenIndex677, depth677
						if buffer[position] != rune('C') {
							goto l670
						}
						position++
					}
				l677:
					{
						position679, tokenIndex679, depth679 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l680
						}
						position++
						goto l679
					l680:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
						if buffer[position] != rune('O') {
							goto l670
						}
						position++
					}
				l679:
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if buffer[position] != rune('N') {
							goto l670
						}
						position++
					}
				l681:
					{
						position683, tokenIndex683, depth683 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l684
						}
						position++
						goto l683
					l684:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
						if buffer[position] != rune('D') {
							goto l670
						}
						position++
					}
				l683:
					{
						position685, tokenIndex685, depth685 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l686
						}
						position++
						goto l685
					l686:
						position, tokenIndex, depth = position685, tokenIndex685, depth685
						if buffer[position] != rune('S') {
							goto l670
						}
						position++
					}
				l685:
					depth--
					add(rulePegText, position672)
				}
				if !_rules[ruleAction48]() {
					goto l670
				}
				depth--
				add(ruleSECONDS, position671)
			}
			return true
		l670:
			position, tokenIndex, depth = position670, tokenIndex670, depth670
			return false
		},
		/* 64 StreamIdentifier <- <(<ident> Action49)> */
		func() bool {
			position687, tokenIndex687, depth687 := position, tokenIndex, depth
			{
				position688 := position
				depth++
				{
					position689 := position
					depth++
					if !_rules[ruleident]() {
						goto l687
					}
					depth--
					add(rulePegText, position689)
				}
				if !_rules[ruleAction49]() {
					goto l687
				}
				depth--
				add(ruleStreamIdentifier, position688)
			}
			return true
		l687:
			position, tokenIndex, depth = position687, tokenIndex687, depth687
			return false
		},
		/* 65 SourceSinkType <- <(<ident> Action50)> */
		func() bool {
			position690, tokenIndex690, depth690 := position, tokenIndex, depth
			{
				position691 := position
				depth++
				{
					position692 := position
					depth++
					if !_rules[ruleident]() {
						goto l690
					}
					depth--
					add(rulePegText, position692)
				}
				if !_rules[ruleAction50]() {
					goto l690
				}
				depth--
				add(ruleSourceSinkType, position691)
			}
			return true
		l690:
			position, tokenIndex, depth = position690, tokenIndex690, depth690
			return false
		},
		/* 66 SourceSinkParamKey <- <(<ident> Action51)> */
		func() bool {
			position693, tokenIndex693, depth693 := position, tokenIndex, depth
			{
				position694 := position
				depth++
				{
					position695 := position
					depth++
					if !_rules[ruleident]() {
						goto l693
					}
					depth--
					add(rulePegText, position695)
				}
				if !_rules[ruleAction51]() {
					goto l693
				}
				depth--
				add(ruleSourceSinkParamKey, position694)
			}
			return true
		l693:
			position, tokenIndex, depth = position693, tokenIndex693, depth693
			return false
		},
		/* 67 SourceSinkParamVal <- <(<([a-z] / [A-Z] / [0-9] / '_')+> Action52)> */
		func() bool {
			position696, tokenIndex696, depth696 := position, tokenIndex, depth
			{
				position697 := position
				depth++
				{
					position698 := position
					depth++
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
							goto l696
						}
						position++
					}
				l701:
				l699:
					{
						position700, tokenIndex700, depth700 := position, tokenIndex, depth
						{
							position705, tokenIndex705, depth705 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l706
							}
							position++
							goto l705
						l706:
							position, tokenIndex, depth = position705, tokenIndex705, depth705
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l707
							}
							position++
							goto l705
						l707:
							position, tokenIndex, depth = position705, tokenIndex705, depth705
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l708
							}
							position++
							goto l705
						l708:
							position, tokenIndex, depth = position705, tokenIndex705, depth705
							if buffer[position] != rune('_') {
								goto l700
							}
							position++
						}
					l705:
						goto l699
					l700:
						position, tokenIndex, depth = position700, tokenIndex700, depth700
					}
					depth--
					add(rulePegText, position698)
				}
				if !_rules[ruleAction52]() {
					goto l696
				}
				depth--
				add(ruleSourceSinkParamVal, position697)
			}
			return true
		l696:
			position, tokenIndex, depth = position696, tokenIndex696, depth696
			return false
		},
		/* 68 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action53)> */
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
						if buffer[position] != rune('o') {
							goto l713
						}
						position++
						goto l712
					l713:
						position, tokenIndex, depth = position712, tokenIndex712, depth712
						if buffer[position] != rune('O') {
							goto l709
						}
						position++
					}
				l712:
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l715
						}
						position++
						goto l714
					l715:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
						if buffer[position] != rune('R') {
							goto l709
						}
						position++
					}
				l714:
					depth--
					add(rulePegText, position711)
				}
				if !_rules[ruleAction53]() {
					goto l709
				}
				depth--
				add(ruleOr, position710)
			}
			return true
		l709:
			position, tokenIndex, depth = position709, tokenIndex709, depth709
			return false
		},
		/* 69 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action54)> */
		func() bool {
			position716, tokenIndex716, depth716 := position, tokenIndex, depth
			{
				position717 := position
				depth++
				{
					position718 := position
					depth++
					{
						position719, tokenIndex719, depth719 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l720
						}
						position++
						goto l719
					l720:
						position, tokenIndex, depth = position719, tokenIndex719, depth719
						if buffer[position] != rune('A') {
							goto l716
						}
						position++
					}
				l719:
					{
						position721, tokenIndex721, depth721 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l722
						}
						position++
						goto l721
					l722:
						position, tokenIndex, depth = position721, tokenIndex721, depth721
						if buffer[position] != rune('N') {
							goto l716
						}
						position++
					}
				l721:
					{
						position723, tokenIndex723, depth723 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l724
						}
						position++
						goto l723
					l724:
						position, tokenIndex, depth = position723, tokenIndex723, depth723
						if buffer[position] != rune('D') {
							goto l716
						}
						position++
					}
				l723:
					depth--
					add(rulePegText, position718)
				}
				if !_rules[ruleAction54]() {
					goto l716
				}
				depth--
				add(ruleAnd, position717)
			}
			return true
		l716:
			position, tokenIndex, depth = position716, tokenIndex716, depth716
			return false
		},
		/* 70 Equal <- <(<'='> Action55)> */
		func() bool {
			position725, tokenIndex725, depth725 := position, tokenIndex, depth
			{
				position726 := position
				depth++
				{
					position727 := position
					depth++
					if buffer[position] != rune('=') {
						goto l725
					}
					position++
					depth--
					add(rulePegText, position727)
				}
				if !_rules[ruleAction55]() {
					goto l725
				}
				depth--
				add(ruleEqual, position726)
			}
			return true
		l725:
			position, tokenIndex, depth = position725, tokenIndex725, depth725
			return false
		},
		/* 71 Less <- <(<'<'> Action56)> */
		func() bool {
			position728, tokenIndex728, depth728 := position, tokenIndex, depth
			{
				position729 := position
				depth++
				{
					position730 := position
					depth++
					if buffer[position] != rune('<') {
						goto l728
					}
					position++
					depth--
					add(rulePegText, position730)
				}
				if !_rules[ruleAction56]() {
					goto l728
				}
				depth--
				add(ruleLess, position729)
			}
			return true
		l728:
			position, tokenIndex, depth = position728, tokenIndex728, depth728
			return false
		},
		/* 72 LessOrEqual <- <(<('<' '=')> Action57)> */
		func() bool {
			position731, tokenIndex731, depth731 := position, tokenIndex, depth
			{
				position732 := position
				depth++
				{
					position733 := position
					depth++
					if buffer[position] != rune('<') {
						goto l731
					}
					position++
					if buffer[position] != rune('=') {
						goto l731
					}
					position++
					depth--
					add(rulePegText, position733)
				}
				if !_rules[ruleAction57]() {
					goto l731
				}
				depth--
				add(ruleLessOrEqual, position732)
			}
			return true
		l731:
			position, tokenIndex, depth = position731, tokenIndex731, depth731
			return false
		},
		/* 73 Greater <- <(<'>'> Action58)> */
		func() bool {
			position734, tokenIndex734, depth734 := position, tokenIndex, depth
			{
				position735 := position
				depth++
				{
					position736 := position
					depth++
					if buffer[position] != rune('>') {
						goto l734
					}
					position++
					depth--
					add(rulePegText, position736)
				}
				if !_rules[ruleAction58]() {
					goto l734
				}
				depth--
				add(ruleGreater, position735)
			}
			return true
		l734:
			position, tokenIndex, depth = position734, tokenIndex734, depth734
			return false
		},
		/* 74 GreaterOrEqual <- <(<('>' '=')> Action59)> */
		func() bool {
			position737, tokenIndex737, depth737 := position, tokenIndex, depth
			{
				position738 := position
				depth++
				{
					position739 := position
					depth++
					if buffer[position] != rune('>') {
						goto l737
					}
					position++
					if buffer[position] != rune('=') {
						goto l737
					}
					position++
					depth--
					add(rulePegText, position739)
				}
				if !_rules[ruleAction59]() {
					goto l737
				}
				depth--
				add(ruleGreaterOrEqual, position738)
			}
			return true
		l737:
			position, tokenIndex, depth = position737, tokenIndex737, depth737
			return false
		},
		/* 75 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action60)> */
		func() bool {
			position740, tokenIndex740, depth740 := position, tokenIndex, depth
			{
				position741 := position
				depth++
				{
					position742 := position
					depth++
					{
						position743, tokenIndex743, depth743 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l744
						}
						position++
						if buffer[position] != rune('=') {
							goto l744
						}
						position++
						goto l743
					l744:
						position, tokenIndex, depth = position743, tokenIndex743, depth743
						if buffer[position] != rune('<') {
							goto l740
						}
						position++
						if buffer[position] != rune('>') {
							goto l740
						}
						position++
					}
				l743:
					depth--
					add(rulePegText, position742)
				}
				if !_rules[ruleAction60]() {
					goto l740
				}
				depth--
				add(ruleNotEqual, position741)
			}
			return true
		l740:
			position, tokenIndex, depth = position740, tokenIndex740, depth740
			return false
		},
		/* 76 Plus <- <(<'+'> Action61)> */
		func() bool {
			position745, tokenIndex745, depth745 := position, tokenIndex, depth
			{
				position746 := position
				depth++
				{
					position747 := position
					depth++
					if buffer[position] != rune('+') {
						goto l745
					}
					position++
					depth--
					add(rulePegText, position747)
				}
				if !_rules[ruleAction61]() {
					goto l745
				}
				depth--
				add(rulePlus, position746)
			}
			return true
		l745:
			position, tokenIndex, depth = position745, tokenIndex745, depth745
			return false
		},
		/* 77 Minus <- <(<'-'> Action62)> */
		func() bool {
			position748, tokenIndex748, depth748 := position, tokenIndex, depth
			{
				position749 := position
				depth++
				{
					position750 := position
					depth++
					if buffer[position] != rune('-') {
						goto l748
					}
					position++
					depth--
					add(rulePegText, position750)
				}
				if !_rules[ruleAction62]() {
					goto l748
				}
				depth--
				add(ruleMinus, position749)
			}
			return true
		l748:
			position, tokenIndex, depth = position748, tokenIndex748, depth748
			return false
		},
		/* 78 Multiply <- <(<'*'> Action63)> */
		func() bool {
			position751, tokenIndex751, depth751 := position, tokenIndex, depth
			{
				position752 := position
				depth++
				{
					position753 := position
					depth++
					if buffer[position] != rune('*') {
						goto l751
					}
					position++
					depth--
					add(rulePegText, position753)
				}
				if !_rules[ruleAction63]() {
					goto l751
				}
				depth--
				add(ruleMultiply, position752)
			}
			return true
		l751:
			position, tokenIndex, depth = position751, tokenIndex751, depth751
			return false
		},
		/* 79 Divide <- <(<'/'> Action64)> */
		func() bool {
			position754, tokenIndex754, depth754 := position, tokenIndex, depth
			{
				position755 := position
				depth++
				{
					position756 := position
					depth++
					if buffer[position] != rune('/') {
						goto l754
					}
					position++
					depth--
					add(rulePegText, position756)
				}
				if !_rules[ruleAction64]() {
					goto l754
				}
				depth--
				add(ruleDivide, position755)
			}
			return true
		l754:
			position, tokenIndex, depth = position754, tokenIndex754, depth754
			return false
		},
		/* 80 Modulo <- <(<'%'> Action65)> */
		func() bool {
			position757, tokenIndex757, depth757 := position, tokenIndex, depth
			{
				position758 := position
				depth++
				{
					position759 := position
					depth++
					if buffer[position] != rune('%') {
						goto l757
					}
					position++
					depth--
					add(rulePegText, position759)
				}
				if !_rules[ruleAction65]() {
					goto l757
				}
				depth--
				add(ruleModulo, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 81 Identifier <- <(<ident> Action66)> */
		func() bool {
			position760, tokenIndex760, depth760 := position, tokenIndex, depth
			{
				position761 := position
				depth++
				{
					position762 := position
					depth++
					if !_rules[ruleident]() {
						goto l760
					}
					depth--
					add(rulePegText, position762)
				}
				if !_rules[ruleAction66]() {
					goto l760
				}
				depth--
				add(ruleIdentifier, position761)
			}
			return true
		l760:
			position, tokenIndex, depth = position760, tokenIndex760, depth760
			return false
		},
		/* 82 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position763, tokenIndex763, depth763 := position, tokenIndex, depth
			{
				position764 := position
				depth++
				{
					position765, tokenIndex765, depth765 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l766
					}
					position++
					goto l765
				l766:
					position, tokenIndex, depth = position765, tokenIndex765, depth765
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l763
					}
					position++
				}
			l765:
			l767:
				{
					position768, tokenIndex768, depth768 := position, tokenIndex, depth
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l771
						}
						position++
						goto l769
					l771:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l772
						}
						position++
						goto l769
					l772:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('_') {
							goto l768
						}
						position++
					}
				l769:
					goto l767
				l768:
					position, tokenIndex, depth = position768, tokenIndex768, depth768
				}
				depth--
				add(ruleident, position764)
			}
			return true
		l763:
			position, tokenIndex, depth = position763, tokenIndex763, depth763
			return false
		},
		/* 83 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position774 := position
				depth++
			l775:
				{
					position776, tokenIndex776, depth776 := position, tokenIndex, depth
					{
						position777, tokenIndex777, depth777 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l778
						}
						position++
						goto l777
					l778:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
						if buffer[position] != rune('\t') {
							goto l779
						}
						position++
						goto l777
					l779:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
						if buffer[position] != rune('\n') {
							goto l776
						}
						position++
					}
				l777:
					goto l775
				l776:
					position, tokenIndex, depth = position776, tokenIndex776, depth776
				}
				depth--
				add(rulesp, position774)
			}
			return true
		},
		/* 85 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 86 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 87 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 88 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 89 Action4 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 90 Action5 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 91 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		nil,
		/* 93 Action7 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 94 Action8 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 95 Action9 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 96 Action10 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 97 Action11 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 98 Action12 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 99 Action13 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 100 Action14 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 101 Action15 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 102 Action16 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 103 Action17 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 104 Action18 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 105 Action19 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 106 Action20 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 107 Action21 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 108 Action22 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 109 Action23 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 110 Action24 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 111 Action25 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 112 Action26 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 113 Action27 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 114 Action28 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 115 Action29 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 116 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 117 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 118 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 119 Action33 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 120 Action34 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 121 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 122 Action36 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 123 Action37 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 124 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 125 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 126 Action40 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 127 Action41 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 128 Action42 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 129 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 130 Action44 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 131 Action45 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 132 Action46 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 133 Action47 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 134 Action48 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 135 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 136 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 137 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 138 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamVal(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 139 Action53 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 140 Action54 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 141 Action55 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 142 Action56 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 143 Action57 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 144 Action58 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 145 Action59 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 146 Action60 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 147 Action61 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 148 Action62 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 149 Action63 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 150 Action64 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 151 Action65 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 152 Action66 <- <{
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

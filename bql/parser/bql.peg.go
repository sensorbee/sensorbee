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
	rulePauseSourceStmt
	ruleResumeSourceStmt
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
	ruleRowMeta
	ruleRowTimestamp
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
	ruleAction67

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
	"PauseSourceStmt",
	"ResumeSourceStmt",
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
	"RowMeta",
	"RowTimestamp",
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
	"Action67",

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
	rules  [157]func() bool
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

			p.AssemblePauseSource()

		case ruleAction7:

			p.AssembleResumeSource()

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
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction38:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction42:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction43:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction44:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction46:

			p.PushComponent(begin, end, Istream)

		case ruleAction47:

			p.PushComponent(begin, end, Dstream)

		case ruleAction48:

			p.PushComponent(begin, end, Rstream)

		case ruleAction49:

			p.PushComponent(begin, end, Tuples)

		case ruleAction50:

			p.PushComponent(begin, end, Seconds)

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction53:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction54:

			p.PushComponent(begin, end, Or)

		case ruleAction55:

			p.PushComponent(begin, end, And)

		case ruleAction56:

			p.PushComponent(begin, end, Equal)

		case ruleAction57:

			p.PushComponent(begin, end, Less)

		case ruleAction58:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction59:

			p.PushComponent(begin, end, Greater)

		case ruleAction60:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction61:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction62:

			p.PushComponent(begin, end, Plus)

		case ruleAction63:

			p.PushComponent(begin, end, Minus)

		case ruleAction64:

			p.PushComponent(begin, end, Multiply)

		case ruleAction65:

			p.PushComponent(begin, end, Divide)

		case ruleAction66:

			p.PushComponent(begin, end, Modulo)

		case ruleAction67:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt)> */
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
						goto l15
					}
					goto l9
				l15:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[rulePauseSourceStmt]() {
						goto l16
					}
					goto l9
				l16:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleResumeSourceStmt]() {
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
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action5)> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('I') {
						goto l171
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('N') {
						goto l171
					}
					position++
				}
			l175:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('S') {
						goto l171
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('E') {
						goto l171
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('R') {
						goto l171
					}
					position++
				}
			l181:
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('T') {
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
					if buffer[position] != rune('i') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('I') {
						goto l171
					}
					position++
				}
			l185:
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if buffer[position] != rune('N') {
						goto l171
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('T') {
						goto l171
					}
					position++
				}
			l189:
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('O') {
						goto l171
					}
					position++
				}
			l191:
				if !_rules[rulesp]() {
					goto l171
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l171
				}
				if !_rules[rulesp]() {
					goto l171
				}
				if !_rules[ruleSelectStmt]() {
					goto l171
				}
				if !_rules[ruleAction5]() {
					goto l171
				}
				depth--
				add(ruleInsertIntoSelectStmt, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 8 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action6)> */
		func() bool {
			position193, tokenIndex193, depth193 := position, tokenIndex, depth
			{
				position194 := position
				depth++
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('P') {
						goto l193
					}
					position++
				}
			l195:
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('A') {
						goto l193
					}
					position++
				}
			l197:
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('U') {
						goto l193
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('S') {
						goto l193
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
						goto l193
					}
					position++
				}
			l203:
				if !_rules[rulesp]() {
					goto l193
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
						goto l193
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
						goto l193
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
						goto l193
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
						goto l193
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
						goto l193
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
						goto l193
					}
					position++
				}
			l215:
				if !_rules[rulesp]() {
					goto l193
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l193
				}
				if !_rules[ruleAction6]() {
					goto l193
				}
				depth--
				add(rulePauseSourceStmt, position194)
			}
			return true
		l193:
			position, tokenIndex, depth = position193, tokenIndex193, depth193
			return false
		},
		/* 9 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action7)> */
		func() bool {
			position217, tokenIndex217, depth217 := position, tokenIndex, depth
			{
				position218 := position
				depth++
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('R') {
						goto l217
					}
					position++
				}
			l219:
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if buffer[position] != rune('E') {
						goto l217
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('S') {
						goto l217
					}
					position++
				}
			l223:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('U') {
						goto l217
					}
					position++
				}
			l225:
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l228
					}
					position++
					goto l227
				l228:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
					if buffer[position] != rune('M') {
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
					if buffer[position] != rune('o') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('O') {
						goto l217
					}
					position++
				}
			l233:
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l236
					}
					position++
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if buffer[position] != rune('U') {
						goto l217
					}
					position++
				}
			l235:
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if buffer[position] != rune('R') {
						goto l217
					}
					position++
				}
			l237:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					if buffer[position] != rune('C') {
						goto l217
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('E') {
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
				if !_rules[ruleAction7]() {
					goto l217
				}
				depth--
				add(ruleResumeSourceStmt, position218)
			}
			return true
		l217:
			position, tokenIndex, depth = position217, tokenIndex217, depth217
			return false
		},
		/* 10 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action8)> */
		func() bool {
			position243, tokenIndex243, depth243 := position, tokenIndex, depth
			{
				position244 := position
				depth++
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l246
					}
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if !_rules[ruleDSTREAM]() {
						goto l247
					}
					goto l245
				l247:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if !_rules[ruleRSTREAM]() {
						goto l243
					}
				}
			l245:
				{
					position248 := position
					depth++
					{
						position249, tokenIndex249, depth249 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l249
						}
						if buffer[position] != rune('[') {
							goto l249
						}
						position++
						if !_rules[rulesp]() {
							goto l249
						}
						{
							position251, tokenIndex251, depth251 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l252
							}
							position++
							goto l251
						l252:
							position, tokenIndex, depth = position251, tokenIndex251, depth251
							if buffer[position] != rune('E') {
								goto l249
							}
							position++
						}
					l251:
						{
							position253, tokenIndex253, depth253 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l254
							}
							position++
							goto l253
						l254:
							position, tokenIndex, depth = position253, tokenIndex253, depth253
							if buffer[position] != rune('V') {
								goto l249
							}
							position++
						}
					l253:
						{
							position255, tokenIndex255, depth255 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l256
							}
							position++
							goto l255
						l256:
							position, tokenIndex, depth = position255, tokenIndex255, depth255
							if buffer[position] != rune('E') {
								goto l249
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
								goto l249
							}
							position++
						}
					l257:
						{
							position259, tokenIndex259, depth259 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l260
							}
							position++
							goto l259
						l260:
							position, tokenIndex, depth = position259, tokenIndex259, depth259
							if buffer[position] != rune('Y') {
								goto l249
							}
							position++
						}
					l259:
						if !_rules[rulesp]() {
							goto l249
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l249
						}
						if !_rules[rulesp]() {
							goto l249
						}
						if buffer[position] != rune(']') {
							goto l249
						}
						position++
						goto l250
					l249:
						position, tokenIndex, depth = position249, tokenIndex249, depth249
					}
				l250:
					depth--
					add(rulePegText, position248)
				}
				if !_rules[ruleAction8]() {
					goto l243
				}
				depth--
				add(ruleEmitter, position244)
			}
			return true
		l243:
			position, tokenIndex, depth = position243, tokenIndex243, depth243
			return false
		},
		/* 11 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position261, tokenIndex261, depth261 := position, tokenIndex, depth
			{
				position262 := position
				depth++
				{
					position263, tokenIndex263, depth263 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l264
					}
				l265:
					{
						position266, tokenIndex266, depth266 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l266
						}
						if buffer[position] != rune(',') {
							goto l266
						}
						position++
						if !_rules[rulesp]() {
							goto l266
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l266
						}
						goto l265
					l266:
						position, tokenIndex, depth = position266, tokenIndex266, depth266
					}
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if !_rules[ruleTimeEmitterInterval]() {
						goto l267
					}
					goto l263
				l267:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if !_rules[ruleTupleEmitterInterval]() {
						goto l261
					}
				}
			l263:
				depth--
				add(ruleEmitterIntervals, position262)
			}
			return true
		l261:
			position, tokenIndex, depth = position261, tokenIndex261, depth261
			return false
		},
		/* 12 TimeEmitterInterval <- <(<TimeInterval> Action9)> */
		func() bool {
			position268, tokenIndex268, depth268 := position, tokenIndex, depth
			{
				position269 := position
				depth++
				{
					position270 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l268
					}
					depth--
					add(rulePegText, position270)
				}
				if !_rules[ruleAction9]() {
					goto l268
				}
				depth--
				add(ruleTimeEmitterInterval, position269)
			}
			return true
		l268:
			position, tokenIndex, depth = position268, tokenIndex268, depth268
			return false
		},
		/* 13 TupleEmitterInterval <- <(<TuplesInterval> Action10)> */
		func() bool {
			position271, tokenIndex271, depth271 := position, tokenIndex, depth
			{
				position272 := position
				depth++
				{
					position273 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l271
					}
					depth--
					add(rulePegText, position273)
				}
				if !_rules[ruleAction10]() {
					goto l271
				}
				depth--
				add(ruleTupleEmitterInterval, position272)
			}
			return true
		l271:
			position, tokenIndex, depth = position271, tokenIndex271, depth271
			return false
		},
		/* 14 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action11)> */
		func() bool {
			position274, tokenIndex274, depth274 := position, tokenIndex, depth
			{
				position275 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l274
				}
				if !_rules[rulesp]() {
					goto l274
				}
				{
					position276, tokenIndex276, depth276 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l277
					}
					position++
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					if buffer[position] != rune('I') {
						goto l274
					}
					position++
				}
			l276:
				{
					position278, tokenIndex278, depth278 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l279
					}
					position++
					goto l278
				l279:
					position, tokenIndex, depth = position278, tokenIndex278, depth278
					if buffer[position] != rune('N') {
						goto l274
					}
					position++
				}
			l278:
				if !_rules[rulesp]() {
					goto l274
				}
				if !_rules[ruleStream]() {
					goto l274
				}
				if !_rules[ruleAction11]() {
					goto l274
				}
				depth--
				add(ruleTupleEmitterFromInterval, position275)
			}
			return true
		l274:
			position, tokenIndex, depth = position274, tokenIndex274, depth274
			return false
		},
		/* 15 Projections <- <(<(Projection sp (',' sp Projection)*)> Action12)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				{
					position282 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l280
					}
					if !_rules[rulesp]() {
						goto l280
					}
				l283:
					{
						position284, tokenIndex284, depth284 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l284
						}
						position++
						if !_rules[rulesp]() {
							goto l284
						}
						if !_rules[ruleProjection]() {
							goto l284
						}
						goto l283
					l284:
						position, tokenIndex, depth = position284, tokenIndex284, depth284
					}
					depth--
					add(rulePegText, position282)
				}
				if !_rules[ruleAction12]() {
					goto l280
				}
				depth--
				add(ruleProjections, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 16 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position285, tokenIndex285, depth285 := position, tokenIndex, depth
			{
				position286 := position
				depth++
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l288
					}
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if !_rules[ruleExpression]() {
						goto l289
					}
					goto l287
				l289:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if !_rules[ruleWildcard]() {
						goto l285
					}
				}
			l287:
				depth--
				add(ruleProjection, position286)
			}
			return true
		l285:
			position, tokenIndex, depth = position285, tokenIndex285, depth285
			return false
		},
		/* 17 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action13)> */
		func() bool {
			position290, tokenIndex290, depth290 := position, tokenIndex, depth
			{
				position291 := position
				depth++
				{
					position292, tokenIndex292, depth292 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l293
					}
					goto l292
				l293:
					position, tokenIndex, depth = position292, tokenIndex292, depth292
					if !_rules[ruleWildcard]() {
						goto l290
					}
				}
			l292:
				if !_rules[rulesp]() {
					goto l290
				}
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l295
					}
					position++
					goto l294
				l295:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if buffer[position] != rune('A') {
						goto l290
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if buffer[position] != rune('S') {
						goto l290
					}
					position++
				}
			l296:
				if !_rules[rulesp]() {
					goto l290
				}
				if !_rules[ruleIdentifier]() {
					goto l290
				}
				if !_rules[ruleAction13]() {
					goto l290
				}
				depth--
				add(ruleAliasExpression, position291)
			}
			return true
		l290:
			position, tokenIndex, depth = position290, tokenIndex290, depth290
			return false
		},
		/* 18 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action14)> */
		func() bool {
			position298, tokenIndex298, depth298 := position, tokenIndex, depth
			{
				position299 := position
				depth++
				{
					position300 := position
					depth++
					{
						position301, tokenIndex301, depth301 := position, tokenIndex, depth
						{
							position303, tokenIndex303, depth303 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l304
							}
							position++
							goto l303
						l304:
							position, tokenIndex, depth = position303, tokenIndex303, depth303
							if buffer[position] != rune('F') {
								goto l301
							}
							position++
						}
					l303:
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l306
							}
							position++
							goto l305
						l306:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
							if buffer[position] != rune('R') {
								goto l301
							}
							position++
						}
					l305:
						{
							position307, tokenIndex307, depth307 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l308
							}
							position++
							goto l307
						l308:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if buffer[position] != rune('O') {
								goto l301
							}
							position++
						}
					l307:
						{
							position309, tokenIndex309, depth309 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l310
							}
							position++
							goto l309
						l310:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if buffer[position] != rune('M') {
								goto l301
							}
							position++
						}
					l309:
						if !_rules[rulesp]() {
							goto l301
						}
						if !_rules[ruleRelations]() {
							goto l301
						}
						if !_rules[rulesp]() {
							goto l301
						}
						goto l302
					l301:
						position, tokenIndex, depth = position301, tokenIndex301, depth301
					}
				l302:
					depth--
					add(rulePegText, position300)
				}
				if !_rules[ruleAction14]() {
					goto l298
				}
				depth--
				add(ruleWindowedFrom, position299)
			}
			return true
		l298:
			position, tokenIndex, depth = position298, tokenIndex298, depth298
			return false
		},
		/* 19 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action15)> */
		func() bool {
			position311, tokenIndex311, depth311 := position, tokenIndex, depth
			{
				position312 := position
				depth++
				{
					position313 := position
					depth++
					{
						position314, tokenIndex314, depth314 := position, tokenIndex, depth
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('F') {
								goto l314
							}
							position++
						}
					l316:
						{
							position318, tokenIndex318, depth318 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l319
							}
							position++
							goto l318
						l319:
							position, tokenIndex, depth = position318, tokenIndex318, depth318
							if buffer[position] != rune('R') {
								goto l314
							}
							position++
						}
					l318:
						{
							position320, tokenIndex320, depth320 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l321
							}
							position++
							goto l320
						l321:
							position, tokenIndex, depth = position320, tokenIndex320, depth320
							if buffer[position] != rune('O') {
								goto l314
							}
							position++
						}
					l320:
						{
							position322, tokenIndex322, depth322 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l323
							}
							position++
							goto l322
						l323:
							position, tokenIndex, depth = position322, tokenIndex322, depth322
							if buffer[position] != rune('M') {
								goto l314
							}
							position++
						}
					l322:
						if !_rules[rulesp]() {
							goto l314
						}
						if !_rules[ruleDefRelations]() {
							goto l314
						}
						if !_rules[rulesp]() {
							goto l314
						}
						goto l315
					l314:
						position, tokenIndex, depth = position314, tokenIndex314, depth314
					}
				l315:
					depth--
					add(rulePegText, position313)
				}
				if !_rules[ruleAction15]() {
					goto l311
				}
				depth--
				add(ruleDefWindowedFrom, position312)
			}
			return true
		l311:
			position, tokenIndex, depth = position311, tokenIndex311, depth311
			return false
		},
		/* 20 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position324, tokenIndex324, depth324 := position, tokenIndex, depth
			{
				position325 := position
				depth++
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l327
					}
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					if !_rules[ruleTuplesInterval]() {
						goto l324
					}
				}
			l326:
				depth--
				add(ruleInterval, position325)
			}
			return true
		l324:
			position, tokenIndex, depth = position324, tokenIndex324, depth324
			return false
		},
		/* 21 TimeInterval <- <(NumericLiteral sp SECONDS Action16)> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l328
				}
				if !_rules[rulesp]() {
					goto l328
				}
				if !_rules[ruleSECONDS]() {
					goto l328
				}
				if !_rules[ruleAction16]() {
					goto l328
				}
				depth--
				add(ruleTimeInterval, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 22 TuplesInterval <- <(NumericLiteral sp TUPLES Action17)> */
		func() bool {
			position330, tokenIndex330, depth330 := position, tokenIndex, depth
			{
				position331 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l330
				}
				if !_rules[rulesp]() {
					goto l330
				}
				if !_rules[ruleTUPLES]() {
					goto l330
				}
				if !_rules[ruleAction17]() {
					goto l330
				}
				depth--
				add(ruleTuplesInterval, position331)
			}
			return true
		l330:
			position, tokenIndex, depth = position330, tokenIndex330, depth330
			return false
		},
		/* 23 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position332, tokenIndex332, depth332 := position, tokenIndex, depth
			{
				position333 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l332
				}
				if !_rules[rulesp]() {
					goto l332
				}
			l334:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l335
					}
					position++
					if !_rules[rulesp]() {
						goto l335
					}
					if !_rules[ruleRelationLike]() {
						goto l335
					}
					goto l334
				l335:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
				}
				depth--
				add(ruleRelations, position333)
			}
			return true
		l332:
			position, tokenIndex, depth = position332, tokenIndex332, depth332
			return false
		},
		/* 24 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position336, tokenIndex336, depth336 := position, tokenIndex, depth
			{
				position337 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l336
				}
				if !_rules[rulesp]() {
					goto l336
				}
			l338:
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l339
					}
					position++
					if !_rules[rulesp]() {
						goto l339
					}
					if !_rules[ruleDefRelationLike]() {
						goto l339
					}
					goto l338
				l339:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
				}
				depth--
				add(ruleDefRelations, position337)
			}
			return true
		l336:
			position, tokenIndex, depth = position336, tokenIndex336, depth336
			return false
		},
		/* 25 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action18)> */
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
							if buffer[position] != rune('w') {
								goto l346
							}
							position++
							goto l345
						l346:
							position, tokenIndex, depth = position345, tokenIndex345, depth345
							if buffer[position] != rune('W') {
								goto l343
							}
							position++
						}
					l345:
						{
							position347, tokenIndex347, depth347 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l348
							}
							position++
							goto l347
						l348:
							position, tokenIndex, depth = position347, tokenIndex347, depth347
							if buffer[position] != rune('H') {
								goto l343
							}
							position++
						}
					l347:
						{
							position349, tokenIndex349, depth349 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l350
							}
							position++
							goto l349
						l350:
							position, tokenIndex, depth = position349, tokenIndex349, depth349
							if buffer[position] != rune('E') {
								goto l343
							}
							position++
						}
					l349:
						{
							position351, tokenIndex351, depth351 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l352
							}
							position++
							goto l351
						l352:
							position, tokenIndex, depth = position351, tokenIndex351, depth351
							if buffer[position] != rune('R') {
								goto l343
							}
							position++
						}
					l351:
						{
							position353, tokenIndex353, depth353 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l354
							}
							position++
							goto l353
						l354:
							position, tokenIndex, depth = position353, tokenIndex353, depth353
							if buffer[position] != rune('E') {
								goto l343
							}
							position++
						}
					l353:
						if !_rules[rulesp]() {
							goto l343
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction18]() {
					goto l340
				}
				depth--
				add(ruleFilter, position341)
			}
			return true
		l340:
			position, tokenIndex, depth = position340, tokenIndex340, depth340
			return false
		},
		/* 26 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action19)> */
		func() bool {
			position355, tokenIndex355, depth355 := position, tokenIndex, depth
			{
				position356 := position
				depth++
				{
					position357 := position
					depth++
					{
						position358, tokenIndex358, depth358 := position, tokenIndex, depth
						{
							position360, tokenIndex360, depth360 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l361
							}
							position++
							goto l360
						l361:
							position, tokenIndex, depth = position360, tokenIndex360, depth360
							if buffer[position] != rune('G') {
								goto l358
							}
							position++
						}
					l360:
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('R') {
								goto l358
							}
							position++
						}
					l362:
						{
							position364, tokenIndex364, depth364 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l365
							}
							position++
							goto l364
						l365:
							position, tokenIndex, depth = position364, tokenIndex364, depth364
							if buffer[position] != rune('O') {
								goto l358
							}
							position++
						}
					l364:
						{
							position366, tokenIndex366, depth366 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l367
							}
							position++
							goto l366
						l367:
							position, tokenIndex, depth = position366, tokenIndex366, depth366
							if buffer[position] != rune('U') {
								goto l358
							}
							position++
						}
					l366:
						{
							position368, tokenIndex368, depth368 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l369
							}
							position++
							goto l368
						l369:
							position, tokenIndex, depth = position368, tokenIndex368, depth368
							if buffer[position] != rune('P') {
								goto l358
							}
							position++
						}
					l368:
						if !_rules[rulesp]() {
							goto l358
						}
						{
							position370, tokenIndex370, depth370 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l371
							}
							position++
							goto l370
						l371:
							position, tokenIndex, depth = position370, tokenIndex370, depth370
							if buffer[position] != rune('B') {
								goto l358
							}
							position++
						}
					l370:
						{
							position372, tokenIndex372, depth372 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l373
							}
							position++
							goto l372
						l373:
							position, tokenIndex, depth = position372, tokenIndex372, depth372
							if buffer[position] != rune('Y') {
								goto l358
							}
							position++
						}
					l372:
						if !_rules[rulesp]() {
							goto l358
						}
						if !_rules[ruleGroupList]() {
							goto l358
						}
						goto l359
					l358:
						position, tokenIndex, depth = position358, tokenIndex358, depth358
					}
				l359:
					depth--
					add(rulePegText, position357)
				}
				if !_rules[ruleAction19]() {
					goto l355
				}
				depth--
				add(ruleGrouping, position356)
			}
			return true
		l355:
			position, tokenIndex, depth = position355, tokenIndex355, depth355
			return false
		},
		/* 27 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				if !_rules[ruleExpression]() {
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
					if !_rules[ruleExpression]() {
						goto l377
					}
					goto l376
				l377:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
				}
				depth--
				add(ruleGroupList, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
			return false
		},
		/* 28 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action20)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				{
					position380 := position
					depth++
					{
						position381, tokenIndex381, depth381 := position, tokenIndex, depth
						{
							position383, tokenIndex383, depth383 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l384
							}
							position++
							goto l383
						l384:
							position, tokenIndex, depth = position383, tokenIndex383, depth383
							if buffer[position] != rune('H') {
								goto l381
							}
							position++
						}
					l383:
						{
							position385, tokenIndex385, depth385 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l386
							}
							position++
							goto l385
						l386:
							position, tokenIndex, depth = position385, tokenIndex385, depth385
							if buffer[position] != rune('A') {
								goto l381
							}
							position++
						}
					l385:
						{
							position387, tokenIndex387, depth387 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l388
							}
							position++
							goto l387
						l388:
							position, tokenIndex, depth = position387, tokenIndex387, depth387
							if buffer[position] != rune('V') {
								goto l381
							}
							position++
						}
					l387:
						{
							position389, tokenIndex389, depth389 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l390
							}
							position++
							goto l389
						l390:
							position, tokenIndex, depth = position389, tokenIndex389, depth389
							if buffer[position] != rune('I') {
								goto l381
							}
							position++
						}
					l389:
						{
							position391, tokenIndex391, depth391 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l392
							}
							position++
							goto l391
						l392:
							position, tokenIndex, depth = position391, tokenIndex391, depth391
							if buffer[position] != rune('N') {
								goto l381
							}
							position++
						}
					l391:
						{
							position393, tokenIndex393, depth393 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l394
							}
							position++
							goto l393
						l394:
							position, tokenIndex, depth = position393, tokenIndex393, depth393
							if buffer[position] != rune('G') {
								goto l381
							}
							position++
						}
					l393:
						if !_rules[rulesp]() {
							goto l381
						}
						if !_rules[ruleExpression]() {
							goto l381
						}
						goto l382
					l381:
						position, tokenIndex, depth = position381, tokenIndex381, depth381
					}
				l382:
					depth--
					add(rulePegText, position380)
				}
				if !_rules[ruleAction20]() {
					goto l378
				}
				depth--
				add(ruleHaving, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 29 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action21))> */
		func() bool {
			position395, tokenIndex395, depth395 := position, tokenIndex, depth
			{
				position396 := position
				depth++
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l398
					}
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if !_rules[ruleStreamWindow]() {
						goto l395
					}
					if !_rules[ruleAction21]() {
						goto l395
					}
				}
			l397:
				depth--
				add(ruleRelationLike, position396)
			}
			return true
		l395:
			position, tokenIndex, depth = position395, tokenIndex395, depth395
			return false
		},
		/* 30 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action22))> */
		func() bool {
			position399, tokenIndex399, depth399 := position, tokenIndex, depth
			{
				position400 := position
				depth++
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l402
					}
					goto l401
				l402:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if !_rules[ruleDefStreamWindow]() {
						goto l399
					}
					if !_rules[ruleAction22]() {
						goto l399
					}
				}
			l401:
				depth--
				add(ruleDefRelationLike, position400)
			}
			return true
		l399:
			position, tokenIndex, depth = position399, tokenIndex399, depth399
			return false
		},
		/* 31 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action23)> */
		func() bool {
			position403, tokenIndex403, depth403 := position, tokenIndex, depth
			{
				position404 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l403
				}
				if !_rules[rulesp]() {
					goto l403
				}
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l406
					}
					position++
					goto l405
				l406:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if buffer[position] != rune('A') {
						goto l403
					}
					position++
				}
			l405:
				{
					position407, tokenIndex407, depth407 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l408
					}
					position++
					goto l407
				l408:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
					if buffer[position] != rune('S') {
						goto l403
					}
					position++
				}
			l407:
				if !_rules[rulesp]() {
					goto l403
				}
				if !_rules[ruleIdentifier]() {
					goto l403
				}
				if !_rules[ruleAction23]() {
					goto l403
				}
				depth--
				add(ruleAliasedStreamWindow, position404)
			}
			return true
		l403:
			position, tokenIndex, depth = position403, tokenIndex403, depth403
			return false
		},
		/* 32 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action24)> */
		func() bool {
			position409, tokenIndex409, depth409 := position, tokenIndex, depth
			{
				position410 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l409
				}
				if !_rules[rulesp]() {
					goto l409
				}
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l412
					}
					position++
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if buffer[position] != rune('A') {
						goto l409
					}
					position++
				}
			l411:
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l414
					}
					position++
					goto l413
				l414:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
					if buffer[position] != rune('S') {
						goto l409
					}
					position++
				}
			l413:
				if !_rules[rulesp]() {
					goto l409
				}
				if !_rules[ruleIdentifier]() {
					goto l409
				}
				if !_rules[ruleAction24]() {
					goto l409
				}
				depth--
				add(ruleDefAliasedStreamWindow, position410)
			}
			return true
		l409:
			position, tokenIndex, depth = position409, tokenIndex409, depth409
			return false
		},
		/* 33 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action25)> */
		func() bool {
			position415, tokenIndex415, depth415 := position, tokenIndex, depth
			{
				position416 := position
				depth++
				if !_rules[ruleStream]() {
					goto l415
				}
				if !_rules[rulesp]() {
					goto l415
				}
				if buffer[position] != rune('[') {
					goto l415
				}
				position++
				if !_rules[rulesp]() {
					goto l415
				}
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l418
					}
					position++
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if buffer[position] != rune('R') {
						goto l415
					}
					position++
				}
			l417:
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l420
					}
					position++
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if buffer[position] != rune('A') {
						goto l415
					}
					position++
				}
			l419:
				{
					position421, tokenIndex421, depth421 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l422
					}
					position++
					goto l421
				l422:
					position, tokenIndex, depth = position421, tokenIndex421, depth421
					if buffer[position] != rune('N') {
						goto l415
					}
					position++
				}
			l421:
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l424
					}
					position++
					goto l423
				l424:
					position, tokenIndex, depth = position423, tokenIndex423, depth423
					if buffer[position] != rune('G') {
						goto l415
					}
					position++
				}
			l423:
				{
					position425, tokenIndex425, depth425 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l426
					}
					position++
					goto l425
				l426:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if buffer[position] != rune('E') {
						goto l415
					}
					position++
				}
			l425:
				if !_rules[rulesp]() {
					goto l415
				}
				if !_rules[ruleInterval]() {
					goto l415
				}
				if !_rules[rulesp]() {
					goto l415
				}
				if buffer[position] != rune(']') {
					goto l415
				}
				position++
				if !_rules[ruleAction25]() {
					goto l415
				}
				depth--
				add(ruleStreamWindow, position416)
			}
			return true
		l415:
			position, tokenIndex, depth = position415, tokenIndex415, depth415
			return false
		},
		/* 34 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action26)> */
		func() bool {
			position427, tokenIndex427, depth427 := position, tokenIndex, depth
			{
				position428 := position
				depth++
				if !_rules[ruleStream]() {
					goto l427
				}
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l429
					}
					if buffer[position] != rune('[') {
						goto l429
					}
					position++
					if !_rules[rulesp]() {
						goto l429
					}
					{
						position431, tokenIndex431, depth431 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l432
						}
						position++
						goto l431
					l432:
						position, tokenIndex, depth = position431, tokenIndex431, depth431
						if buffer[position] != rune('R') {
							goto l429
						}
						position++
					}
				l431:
					{
						position433, tokenIndex433, depth433 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l434
						}
						position++
						goto l433
					l434:
						position, tokenIndex, depth = position433, tokenIndex433, depth433
						if buffer[position] != rune('A') {
							goto l429
						}
						position++
					}
				l433:
					{
						position435, tokenIndex435, depth435 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l436
						}
						position++
						goto l435
					l436:
						position, tokenIndex, depth = position435, tokenIndex435, depth435
						if buffer[position] != rune('N') {
							goto l429
						}
						position++
					}
				l435:
					{
						position437, tokenIndex437, depth437 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l438
						}
						position++
						goto l437
					l438:
						position, tokenIndex, depth = position437, tokenIndex437, depth437
						if buffer[position] != rune('G') {
							goto l429
						}
						position++
					}
				l437:
					{
						position439, tokenIndex439, depth439 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l440
						}
						position++
						goto l439
					l440:
						position, tokenIndex, depth = position439, tokenIndex439, depth439
						if buffer[position] != rune('E') {
							goto l429
						}
						position++
					}
				l439:
					if !_rules[rulesp]() {
						goto l429
					}
					if !_rules[ruleInterval]() {
						goto l429
					}
					if !_rules[rulesp]() {
						goto l429
					}
					if buffer[position] != rune(']') {
						goto l429
					}
					position++
					goto l430
				l429:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
				}
			l430:
				if !_rules[ruleAction26]() {
					goto l427
				}
				depth--
				add(ruleDefStreamWindow, position428)
			}
			return true
		l427:
			position, tokenIndex, depth = position427, tokenIndex427, depth427
			return false
		},
		/* 35 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action27)> */
		func() bool {
			position441, tokenIndex441, depth441 := position, tokenIndex, depth
			{
				position442 := position
				depth++
				{
					position443 := position
					depth++
					{
						position444, tokenIndex444, depth444 := position, tokenIndex, depth
						{
							position446, tokenIndex446, depth446 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l447
							}
							position++
							goto l446
						l447:
							position, tokenIndex, depth = position446, tokenIndex446, depth446
							if buffer[position] != rune('W') {
								goto l444
							}
							position++
						}
					l446:
						{
							position448, tokenIndex448, depth448 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l449
							}
							position++
							goto l448
						l449:
							position, tokenIndex, depth = position448, tokenIndex448, depth448
							if buffer[position] != rune('I') {
								goto l444
							}
							position++
						}
					l448:
						{
							position450, tokenIndex450, depth450 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l451
							}
							position++
							goto l450
						l451:
							position, tokenIndex, depth = position450, tokenIndex450, depth450
							if buffer[position] != rune('T') {
								goto l444
							}
							position++
						}
					l450:
						{
							position452, tokenIndex452, depth452 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l453
							}
							position++
							goto l452
						l453:
							position, tokenIndex, depth = position452, tokenIndex452, depth452
							if buffer[position] != rune('H') {
								goto l444
							}
							position++
						}
					l452:
						if !_rules[rulesp]() {
							goto l444
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l444
						}
						if !_rules[rulesp]() {
							goto l444
						}
					l454:
						{
							position455, tokenIndex455, depth455 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l455
							}
							position++
							if !_rules[rulesp]() {
								goto l455
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l455
							}
							goto l454
						l455:
							position, tokenIndex, depth = position455, tokenIndex455, depth455
						}
						goto l445
					l444:
						position, tokenIndex, depth = position444, tokenIndex444, depth444
					}
				l445:
					depth--
					add(rulePegText, position443)
				}
				if !_rules[ruleAction27]() {
					goto l441
				}
				depth--
				add(ruleSourceSinkSpecs, position442)
			}
			return true
		l441:
			position, tokenIndex, depth = position441, tokenIndex441, depth441
			return false
		},
		/* 36 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action28)> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l456
				}
				if buffer[position] != rune('=') {
					goto l456
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l456
				}
				if !_rules[ruleAction28]() {
					goto l456
				}
				depth--
				add(ruleSourceSinkParam, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 37 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position458, tokenIndex458, depth458 := position, tokenIndex, depth
			{
				position459 := position
				depth++
				{
					position460, tokenIndex460, depth460 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l461
					}
					goto l460
				l461:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
					if !_rules[ruleLiteral]() {
						goto l458
					}
				}
			l460:
				depth--
				add(ruleSourceSinkParamVal, position459)
			}
			return true
		l458:
			position, tokenIndex, depth = position458, tokenIndex458, depth458
			return false
		},
		/* 38 Expression <- <orExpr> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l462
				}
				depth--
				add(ruleExpression, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 39 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action29)> */
		func() bool {
			position464, tokenIndex464, depth464 := position, tokenIndex, depth
			{
				position465 := position
				depth++
				{
					position466 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l464
					}
					if !_rules[rulesp]() {
						goto l464
					}
					{
						position467, tokenIndex467, depth467 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l467
						}
						if !_rules[rulesp]() {
							goto l467
						}
						if !_rules[ruleandExpr]() {
							goto l467
						}
						goto l468
					l467:
						position, tokenIndex, depth = position467, tokenIndex467, depth467
					}
				l468:
					depth--
					add(rulePegText, position466)
				}
				if !_rules[ruleAction29]() {
					goto l464
				}
				depth--
				add(ruleorExpr, position465)
			}
			return true
		l464:
			position, tokenIndex, depth = position464, tokenIndex464, depth464
			return false
		},
		/* 40 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action30)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				{
					position471 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l469
					}
					if !_rules[rulesp]() {
						goto l469
					}
					{
						position472, tokenIndex472, depth472 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l472
						}
						if !_rules[rulesp]() {
							goto l472
						}
						if !_rules[rulecomparisonExpr]() {
							goto l472
						}
						goto l473
					l472:
						position, tokenIndex, depth = position472, tokenIndex472, depth472
					}
				l473:
					depth--
					add(rulePegText, position471)
				}
				if !_rules[ruleAction30]() {
					goto l469
				}
				depth--
				add(ruleandExpr, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 41 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action31)> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				{
					position476 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l474
					}
					if !_rules[rulesp]() {
						goto l474
					}
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l477
						}
						if !_rules[rulesp]() {
							goto l477
						}
						if !_rules[ruletermExpr]() {
							goto l477
						}
						goto l478
					l477:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
					}
				l478:
					depth--
					add(rulePegText, position476)
				}
				if !_rules[ruleAction31]() {
					goto l474
				}
				depth--
				add(rulecomparisonExpr, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 42 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action32)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				{
					position481 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l479
					}
					if !_rules[rulesp]() {
						goto l479
					}
					{
						position482, tokenIndex482, depth482 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l482
						}
						if !_rules[rulesp]() {
							goto l482
						}
						if !_rules[ruleproductExpr]() {
							goto l482
						}
						goto l483
					l482:
						position, tokenIndex, depth = position482, tokenIndex482, depth482
					}
				l483:
					depth--
					add(rulePegText, position481)
				}
				if !_rules[ruleAction32]() {
					goto l479
				}
				depth--
				add(ruletermExpr, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 43 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action33)> */
		func() bool {
			position484, tokenIndex484, depth484 := position, tokenIndex, depth
			{
				position485 := position
				depth++
				{
					position486 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l484
					}
					if !_rules[rulesp]() {
						goto l484
					}
					{
						position487, tokenIndex487, depth487 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l487
						}
						if !_rules[rulesp]() {
							goto l487
						}
						if !_rules[rulebaseExpr]() {
							goto l487
						}
						goto l488
					l487:
						position, tokenIndex, depth = position487, tokenIndex487, depth487
					}
				l488:
					depth--
					add(rulePegText, position486)
				}
				if !_rules[ruleAction33]() {
					goto l484
				}
				depth--
				add(ruleproductExpr, position485)
			}
			return true
		l484:
			position, tokenIndex, depth = position484, tokenIndex484, depth484
			return false
		},
		/* 44 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position489, tokenIndex489, depth489 := position, tokenIndex, depth
			{
				position490 := position
				depth++
				{
					position491, tokenIndex491, depth491 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l492
					}
					position++
					if !_rules[rulesp]() {
						goto l492
					}
					if !_rules[ruleExpression]() {
						goto l492
					}
					if !_rules[rulesp]() {
						goto l492
					}
					if buffer[position] != rune(')') {
						goto l492
					}
					position++
					goto l491
				l492:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if !_rules[ruleBooleanLiteral]() {
						goto l493
					}
					goto l491
				l493:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if !_rules[ruleFuncApp]() {
						goto l494
					}
					goto l491
				l494:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if !_rules[ruleRowMeta]() {
						goto l495
					}
					goto l491
				l495:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if !_rules[ruleRowValue]() {
						goto l496
					}
					goto l491
				l496:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if !_rules[ruleLiteral]() {
						goto l489
					}
				}
			l491:
				depth--
				add(rulebaseExpr, position490)
			}
			return true
		l489:
			position, tokenIndex, depth = position489, tokenIndex489, depth489
			return false
		},
		/* 45 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action34)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l497
				}
				if !_rules[rulesp]() {
					goto l497
				}
				if buffer[position] != rune('(') {
					goto l497
				}
				position++
				if !_rules[rulesp]() {
					goto l497
				}
				if !_rules[ruleFuncParams]() {
					goto l497
				}
				if !_rules[rulesp]() {
					goto l497
				}
				if buffer[position] != rune(')') {
					goto l497
				}
				position++
				if !_rules[ruleAction34]() {
					goto l497
				}
				depth--
				add(ruleFuncApp, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 46 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action35)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				{
					position501 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l499
					}
					if !_rules[rulesp]() {
						goto l499
					}
				l502:
					{
						position503, tokenIndex503, depth503 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l503
						}
						position++
						if !_rules[rulesp]() {
							goto l503
						}
						if !_rules[ruleExpression]() {
							goto l503
						}
						goto l502
					l503:
						position, tokenIndex, depth = position503, tokenIndex503, depth503
					}
					depth--
					add(rulePegText, position501)
				}
				if !_rules[ruleAction35]() {
					goto l499
				}
				depth--
				add(ruleFuncParams, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 47 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position504, tokenIndex504, depth504 := position, tokenIndex, depth
			{
				position505 := position
				depth++
				{
					position506, tokenIndex506, depth506 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l507
					}
					goto l506
				l507:
					position, tokenIndex, depth = position506, tokenIndex506, depth506
					if !_rules[ruleNumericLiteral]() {
						goto l508
					}
					goto l506
				l508:
					position, tokenIndex, depth = position506, tokenIndex506, depth506
					if !_rules[ruleStringLiteral]() {
						goto l504
					}
				}
			l506:
				depth--
				add(ruleLiteral, position505)
			}
			return true
		l504:
			position, tokenIndex, depth = position504, tokenIndex504, depth504
			return false
		},
		/* 48 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l512
					}
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if !_rules[ruleNotEqual]() {
						goto l513
					}
					goto l511
				l513:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if !_rules[ruleLessOrEqual]() {
						goto l514
					}
					goto l511
				l514:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if !_rules[ruleLess]() {
						goto l515
					}
					goto l511
				l515:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if !_rules[ruleGreaterOrEqual]() {
						goto l516
					}
					goto l511
				l516:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if !_rules[ruleGreater]() {
						goto l517
					}
					goto l511
				l517:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if !_rules[ruleNotEqual]() {
						goto l509
					}
				}
			l511:
				depth--
				add(ruleComparisonOp, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 49 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position518, tokenIndex518, depth518 := position, tokenIndex, depth
			{
				position519 := position
				depth++
				{
					position520, tokenIndex520, depth520 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l521
					}
					goto l520
				l521:
					position, tokenIndex, depth = position520, tokenIndex520, depth520
					if !_rules[ruleMinus]() {
						goto l518
					}
				}
			l520:
				depth--
				add(rulePlusMinusOp, position519)
			}
			return true
		l518:
			position, tokenIndex, depth = position518, tokenIndex518, depth518
			return false
		},
		/* 50 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position522, tokenIndex522, depth522 := position, tokenIndex, depth
			{
				position523 := position
				depth++
				{
					position524, tokenIndex524, depth524 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l525
					}
					goto l524
				l525:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleDivide]() {
						goto l526
					}
					goto l524
				l526:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleModulo]() {
						goto l522
					}
				}
			l524:
				depth--
				add(ruleMultDivOp, position523)
			}
			return true
		l522:
			position, tokenIndex, depth = position522, tokenIndex522, depth522
			return false
		},
		/* 51 Stream <- <(<ident> Action36)> */
		func() bool {
			position527, tokenIndex527, depth527 := position, tokenIndex, depth
			{
				position528 := position
				depth++
				{
					position529 := position
					depth++
					if !_rules[ruleident]() {
						goto l527
					}
					depth--
					add(rulePegText, position529)
				}
				if !_rules[ruleAction36]() {
					goto l527
				}
				depth--
				add(ruleStream, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 52 RowMeta <- <RowTimestamp> */
		func() bool {
			position530, tokenIndex530, depth530 := position, tokenIndex, depth
			{
				position531 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l530
				}
				depth--
				add(ruleRowMeta, position531)
			}
			return true
		l530:
			position, tokenIndex, depth = position530, tokenIndex530, depth530
			return false
		},
		/* 53 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action37)> */
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
					if buffer[position] != rune('t') {
						goto l532
					}
					position++
					if buffer[position] != rune('s') {
						goto l532
					}
					position++
					if buffer[position] != rune('(') {
						goto l532
					}
					position++
					if buffer[position] != rune(')') {
						goto l532
					}
					position++
					depth--
					add(rulePegText, position534)
				}
				if !_rules[ruleAction37]() {
					goto l532
				}
				depth--
				add(ruleRowTimestamp, position533)
			}
			return true
		l532:
			position, tokenIndex, depth = position532, tokenIndex532, depth532
			return false
		},
		/* 54 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action38)> */
		func() bool {
			position537, tokenIndex537, depth537 := position, tokenIndex, depth
			{
				position538 := position
				depth++
				{
					position539 := position
					depth++
					{
						position540, tokenIndex540, depth540 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l540
						}
						if buffer[position] != rune(':') {
							goto l540
						}
						position++
						goto l541
					l540:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
					}
				l541:
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l543
						}
						position++
						goto l542
					l543:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l537
						}
						position++
					}
				l542:
				l544:
					{
						position545, tokenIndex545, depth545 := position, tokenIndex, depth
						{
							position546, tokenIndex546, depth546 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l547
							}
							position++
							goto l546
						l547:
							position, tokenIndex, depth = position546, tokenIndex546, depth546
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l548
							}
							position++
							goto l546
						l548:
							position, tokenIndex, depth = position546, tokenIndex546, depth546
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l549
							}
							position++
							goto l546
						l549:
							position, tokenIndex, depth = position546, tokenIndex546, depth546
							if buffer[position] != rune('_') {
								goto l550
							}
							position++
							goto l546
						l550:
							position, tokenIndex, depth = position546, tokenIndex546, depth546
							if buffer[position] != rune('.') {
								goto l545
							}
							position++
						}
					l546:
						goto l544
					l545:
						position, tokenIndex, depth = position545, tokenIndex545, depth545
					}
					depth--
					add(rulePegText, position539)
				}
				if !_rules[ruleAction38]() {
					goto l537
				}
				depth--
				add(ruleRowValue, position538)
			}
			return true
		l537:
			position, tokenIndex, depth = position537, tokenIndex537, depth537
			return false
		},
		/* 55 NumericLiteral <- <(<('-'? [0-9]+)> Action39)> */
		func() bool {
			position551, tokenIndex551, depth551 := position, tokenIndex, depth
			{
				position552 := position
				depth++
				{
					position553 := position
					depth++
					{
						position554, tokenIndex554, depth554 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l554
						}
						position++
						goto l555
					l554:
						position, tokenIndex, depth = position554, tokenIndex554, depth554
					}
				l555:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l551
					}
					position++
				l556:
					{
						position557, tokenIndex557, depth557 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l557
						}
						position++
						goto l556
					l557:
						position, tokenIndex, depth = position557, tokenIndex557, depth557
					}
					depth--
					add(rulePegText, position553)
				}
				if !_rules[ruleAction39]() {
					goto l551
				}
				depth--
				add(ruleNumericLiteral, position552)
			}
			return true
		l551:
			position, tokenIndex, depth = position551, tokenIndex551, depth551
			return false
		},
		/* 56 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action40)> */
		func() bool {
			position558, tokenIndex558, depth558 := position, tokenIndex, depth
			{
				position559 := position
				depth++
				{
					position560 := position
					depth++
					{
						position561, tokenIndex561, depth561 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l561
						}
						position++
						goto l562
					l561:
						position, tokenIndex, depth = position561, tokenIndex561, depth561
					}
				l562:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l558
					}
					position++
				l563:
					{
						position564, tokenIndex564, depth564 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l564
						}
						position++
						goto l563
					l564:
						position, tokenIndex, depth = position564, tokenIndex564, depth564
					}
					if buffer[position] != rune('.') {
						goto l558
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l558
					}
					position++
				l565:
					{
						position566, tokenIndex566, depth566 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l566
						}
						position++
						goto l565
					l566:
						position, tokenIndex, depth = position566, tokenIndex566, depth566
					}
					depth--
					add(rulePegText, position560)
				}
				if !_rules[ruleAction40]() {
					goto l558
				}
				depth--
				add(ruleFloatLiteral, position559)
			}
			return true
		l558:
			position, tokenIndex, depth = position558, tokenIndex558, depth558
			return false
		},
		/* 57 Function <- <(<ident> Action41)> */
		func() bool {
			position567, tokenIndex567, depth567 := position, tokenIndex, depth
			{
				position568 := position
				depth++
				{
					position569 := position
					depth++
					if !_rules[ruleident]() {
						goto l567
					}
					depth--
					add(rulePegText, position569)
				}
				if !_rules[ruleAction41]() {
					goto l567
				}
				depth--
				add(ruleFunction, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 58 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position570, tokenIndex570, depth570 := position, tokenIndex, depth
			{
				position571 := position
				depth++
				{
					position572, tokenIndex572, depth572 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l573
					}
					goto l572
				l573:
					position, tokenIndex, depth = position572, tokenIndex572, depth572
					if !_rules[ruleFALSE]() {
						goto l570
					}
				}
			l572:
				depth--
				add(ruleBooleanLiteral, position571)
			}
			return true
		l570:
			position, tokenIndex, depth = position570, tokenIndex570, depth570
			return false
		},
		/* 59 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action42)> */
		func() bool {
			position574, tokenIndex574, depth574 := position, tokenIndex, depth
			{
				position575 := position
				depth++
				{
					position576 := position
					depth++
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('T') {
							goto l574
						}
						position++
					}
				l577:
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l580
						}
						position++
						goto l579
					l580:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
						if buffer[position] != rune('R') {
							goto l574
						}
						position++
					}
				l579:
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l582
						}
						position++
						goto l581
					l582:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
						if buffer[position] != rune('U') {
							goto l574
						}
						position++
					}
				l581:
					{
						position583, tokenIndex583, depth583 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l584
						}
						position++
						goto l583
					l584:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
						if buffer[position] != rune('E') {
							goto l574
						}
						position++
					}
				l583:
					depth--
					add(rulePegText, position576)
				}
				if !_rules[ruleAction42]() {
					goto l574
				}
				depth--
				add(ruleTRUE, position575)
			}
			return true
		l574:
			position, tokenIndex, depth = position574, tokenIndex574, depth574
			return false
		},
		/* 60 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action43)> */
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
						if buffer[position] != rune('f') {
							goto l589
						}
						position++
						goto l588
					l589:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
						if buffer[position] != rune('F') {
							goto l585
						}
						position++
					}
				l588:
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
						if buffer[position] != rune('A') {
							goto l585
						}
						position++
					}
				l590:
					{
						position592, tokenIndex592, depth592 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex, depth = position592, tokenIndex592, depth592
						if buffer[position] != rune('L') {
							goto l585
						}
						position++
					}
				l592:
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
						if buffer[position] != rune('S') {
							goto l585
						}
						position++
					}
				l594:
					{
						position596, tokenIndex596, depth596 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l597
						}
						position++
						goto l596
					l597:
						position, tokenIndex, depth = position596, tokenIndex596, depth596
						if buffer[position] != rune('E') {
							goto l585
						}
						position++
					}
				l596:
					depth--
					add(rulePegText, position587)
				}
				if !_rules[ruleAction43]() {
					goto l585
				}
				depth--
				add(ruleFALSE, position586)
			}
			return true
		l585:
			position, tokenIndex, depth = position585, tokenIndex585, depth585
			return false
		},
		/* 61 Wildcard <- <(<'*'> Action44)> */
		func() bool {
			position598, tokenIndex598, depth598 := position, tokenIndex, depth
			{
				position599 := position
				depth++
				{
					position600 := position
					depth++
					if buffer[position] != rune('*') {
						goto l598
					}
					position++
					depth--
					add(rulePegText, position600)
				}
				if !_rules[ruleAction44]() {
					goto l598
				}
				depth--
				add(ruleWildcard, position599)
			}
			return true
		l598:
			position, tokenIndex, depth = position598, tokenIndex598, depth598
			return false
		},
		/* 62 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action45)> */
		func() bool {
			position601, tokenIndex601, depth601 := position, tokenIndex, depth
			{
				position602 := position
				depth++
				{
					position603 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l601
					}
					position++
				l604:
					{
						position605, tokenIndex605, depth605 := position, tokenIndex, depth
						{
							position606, tokenIndex606, depth606 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l607
							}
							position++
							if buffer[position] != rune('\'') {
								goto l607
							}
							position++
							goto l606
						l607:
							position, tokenIndex, depth = position606, tokenIndex606, depth606
							{
								position608, tokenIndex608, depth608 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l608
								}
								position++
								goto l605
							l608:
								position, tokenIndex, depth = position608, tokenIndex608, depth608
							}
							if !matchDot() {
								goto l605
							}
						}
					l606:
						goto l604
					l605:
						position, tokenIndex, depth = position605, tokenIndex605, depth605
					}
					if buffer[position] != rune('\'') {
						goto l601
					}
					position++
					depth--
					add(rulePegText, position603)
				}
				if !_rules[ruleAction45]() {
					goto l601
				}
				depth--
				add(ruleStringLiteral, position602)
			}
			return true
		l601:
			position, tokenIndex, depth = position601, tokenIndex601, depth601
			return false
		},
		/* 63 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action46)> */
		func() bool {
			position609, tokenIndex609, depth609 := position, tokenIndex, depth
			{
				position610 := position
				depth++
				{
					position611 := position
					depth++
					{
						position612, tokenIndex612, depth612 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l613
						}
						position++
						goto l612
					l613:
						position, tokenIndex, depth = position612, tokenIndex612, depth612
						if buffer[position] != rune('I') {
							goto l609
						}
						position++
					}
				l612:
					{
						position614, tokenIndex614, depth614 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l615
						}
						position++
						goto l614
					l615:
						position, tokenIndex, depth = position614, tokenIndex614, depth614
						if buffer[position] != rune('S') {
							goto l609
						}
						position++
					}
				l614:
					{
						position616, tokenIndex616, depth616 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l617
						}
						position++
						goto l616
					l617:
						position, tokenIndex, depth = position616, tokenIndex616, depth616
						if buffer[position] != rune('T') {
							goto l609
						}
						position++
					}
				l616:
					{
						position618, tokenIndex618, depth618 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l619
						}
						position++
						goto l618
					l619:
						position, tokenIndex, depth = position618, tokenIndex618, depth618
						if buffer[position] != rune('R') {
							goto l609
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
							goto l609
						}
						position++
					}
				l620:
					{
						position622, tokenIndex622, depth622 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l623
						}
						position++
						goto l622
					l623:
						position, tokenIndex, depth = position622, tokenIndex622, depth622
						if buffer[position] != rune('A') {
							goto l609
						}
						position++
					}
				l622:
					{
						position624, tokenIndex624, depth624 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('M') {
							goto l609
						}
						position++
					}
				l624:
					depth--
					add(rulePegText, position611)
				}
				if !_rules[ruleAction46]() {
					goto l609
				}
				depth--
				add(ruleISTREAM, position610)
			}
			return true
		l609:
			position, tokenIndex, depth = position609, tokenIndex609, depth609
			return false
		},
		/* 64 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action47)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
				{
					position628 := position
					depth++
					{
						position629, tokenIndex629, depth629 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l630
						}
						position++
						goto l629
					l630:
						position, tokenIndex, depth = position629, tokenIndex629, depth629
						if buffer[position] != rune('D') {
							goto l626
						}
						position++
					}
				l629:
					{
						position631, tokenIndex631, depth631 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l632
						}
						position++
						goto l631
					l632:
						position, tokenIndex, depth = position631, tokenIndex631, depth631
						if buffer[position] != rune('S') {
							goto l626
						}
						position++
					}
				l631:
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l634
						}
						position++
						goto l633
					l634:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
						if buffer[position] != rune('T') {
							goto l626
						}
						position++
					}
				l633:
					{
						position635, tokenIndex635, depth635 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l636
						}
						position++
						goto l635
					l636:
						position, tokenIndex, depth = position635, tokenIndex635, depth635
						if buffer[position] != rune('R') {
							goto l626
						}
						position++
					}
				l635:
					{
						position637, tokenIndex637, depth637 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l638
						}
						position++
						goto l637
					l638:
						position, tokenIndex, depth = position637, tokenIndex637, depth637
						if buffer[position] != rune('E') {
							goto l626
						}
						position++
					}
				l637:
					{
						position639, tokenIndex639, depth639 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l640
						}
						position++
						goto l639
					l640:
						position, tokenIndex, depth = position639, tokenIndex639, depth639
						if buffer[position] != rune('A') {
							goto l626
						}
						position++
					}
				l639:
					{
						position641, tokenIndex641, depth641 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l642
						}
						position++
						goto l641
					l642:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
						if buffer[position] != rune('M') {
							goto l626
						}
						position++
					}
				l641:
					depth--
					add(rulePegText, position628)
				}
				if !_rules[ruleAction47]() {
					goto l626
				}
				depth--
				add(ruleDSTREAM, position627)
			}
			return true
		l626:
			position, tokenIndex, depth = position626, tokenIndex626, depth626
			return false
		},
		/* 65 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action48)> */
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
						if buffer[position] != rune('r') {
							goto l647
						}
						position++
						goto l646
					l647:
						position, tokenIndex, depth = position646, tokenIndex646, depth646
						if buffer[position] != rune('R') {
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
				if !_rules[ruleAction48]() {
					goto l643
				}
				depth--
				add(ruleRSTREAM, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 66 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action49)> */
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
						if buffer[position] != rune('t') {
							goto l664
						}
						position++
						goto l663
					l664:
						position, tokenIndex, depth = position663, tokenIndex663, depth663
						if buffer[position] != rune('T') {
							goto l660
						}
						position++
					}
				l663:
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l666
						}
						position++
						goto l665
					l666:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if buffer[position] != rune('U') {
							goto l660
						}
						position++
					}
				l665:
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if buffer[position] != rune('P') {
							goto l660
						}
						position++
					}
				l667:
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l670
						}
						position++
						goto l669
					l670:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if buffer[position] != rune('L') {
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
						if buffer[position] != rune('s') {
							goto l674
						}
						position++
						goto l673
					l674:
						position, tokenIndex, depth = position673, tokenIndex673, depth673
						if buffer[position] != rune('S') {
							goto l660
						}
						position++
					}
				l673:
					depth--
					add(rulePegText, position662)
				}
				if !_rules[ruleAction49]() {
					goto l660
				}
				depth--
				add(ruleTUPLES, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 67 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action50)> */
		func() bool {
			position675, tokenIndex675, depth675 := position, tokenIndex, depth
			{
				position676 := position
				depth++
				{
					position677 := position
					depth++
					{
						position678, tokenIndex678, depth678 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l679
						}
						position++
						goto l678
					l679:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
						if buffer[position] != rune('S') {
							goto l675
						}
						position++
					}
				l678:
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l681
						}
						position++
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if buffer[position] != rune('E') {
							goto l675
						}
						position++
					}
				l680:
					{
						position682, tokenIndex682, depth682 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l683
						}
						position++
						goto l682
					l683:
						position, tokenIndex, depth = position682, tokenIndex682, depth682
						if buffer[position] != rune('C') {
							goto l675
						}
						position++
					}
				l682:
					{
						position684, tokenIndex684, depth684 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l685
						}
						position++
						goto l684
					l685:
						position, tokenIndex, depth = position684, tokenIndex684, depth684
						if buffer[position] != rune('O') {
							goto l675
						}
						position++
					}
				l684:
					{
						position686, tokenIndex686, depth686 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
						if buffer[position] != rune('N') {
							goto l675
						}
						position++
					}
				l686:
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if buffer[position] != rune('D') {
							goto l675
						}
						position++
					}
				l688:
					{
						position690, tokenIndex690, depth690 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l691
						}
						position++
						goto l690
					l691:
						position, tokenIndex, depth = position690, tokenIndex690, depth690
						if buffer[position] != rune('S') {
							goto l675
						}
						position++
					}
				l690:
					depth--
					add(rulePegText, position677)
				}
				if !_rules[ruleAction50]() {
					goto l675
				}
				depth--
				add(ruleSECONDS, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 68 StreamIdentifier <- <(<ident> Action51)> */
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
				if !_rules[ruleAction51]() {
					goto l692
				}
				depth--
				add(ruleStreamIdentifier, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 69 SourceSinkType <- <(<ident> Action52)> */
		func() bool {
			position695, tokenIndex695, depth695 := position, tokenIndex, depth
			{
				position696 := position
				depth++
				{
					position697 := position
					depth++
					if !_rules[ruleident]() {
						goto l695
					}
					depth--
					add(rulePegText, position697)
				}
				if !_rules[ruleAction52]() {
					goto l695
				}
				depth--
				add(ruleSourceSinkType, position696)
			}
			return true
		l695:
			position, tokenIndex, depth = position695, tokenIndex695, depth695
			return false
		},
		/* 70 SourceSinkParamKey <- <(<ident> Action53)> */
		func() bool {
			position698, tokenIndex698, depth698 := position, tokenIndex, depth
			{
				position699 := position
				depth++
				{
					position700 := position
					depth++
					if !_rules[ruleident]() {
						goto l698
					}
					depth--
					add(rulePegText, position700)
				}
				if !_rules[ruleAction53]() {
					goto l698
				}
				depth--
				add(ruleSourceSinkParamKey, position699)
			}
			return true
		l698:
			position, tokenIndex, depth = position698, tokenIndex698, depth698
			return false
		},
		/* 71 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action54)> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				{
					position703 := position
					depth++
					{
						position704, tokenIndex704, depth704 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l705
						}
						position++
						goto l704
					l705:
						position, tokenIndex, depth = position704, tokenIndex704, depth704
						if buffer[position] != rune('O') {
							goto l701
						}
						position++
					}
				l704:
					{
						position706, tokenIndex706, depth706 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l707
						}
						position++
						goto l706
					l707:
						position, tokenIndex, depth = position706, tokenIndex706, depth706
						if buffer[position] != rune('R') {
							goto l701
						}
						position++
					}
				l706:
					depth--
					add(rulePegText, position703)
				}
				if !_rules[ruleAction54]() {
					goto l701
				}
				depth--
				add(ruleOr, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 72 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action55)> */
		func() bool {
			position708, tokenIndex708, depth708 := position, tokenIndex, depth
			{
				position709 := position
				depth++
				{
					position710 := position
					depth++
					{
						position711, tokenIndex711, depth711 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l712
						}
						position++
						goto l711
					l712:
						position, tokenIndex, depth = position711, tokenIndex711, depth711
						if buffer[position] != rune('A') {
							goto l708
						}
						position++
					}
				l711:
					{
						position713, tokenIndex713, depth713 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l714
						}
						position++
						goto l713
					l714:
						position, tokenIndex, depth = position713, tokenIndex713, depth713
						if buffer[position] != rune('N') {
							goto l708
						}
						position++
					}
				l713:
					{
						position715, tokenIndex715, depth715 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l716
						}
						position++
						goto l715
					l716:
						position, tokenIndex, depth = position715, tokenIndex715, depth715
						if buffer[position] != rune('D') {
							goto l708
						}
						position++
					}
				l715:
					depth--
					add(rulePegText, position710)
				}
				if !_rules[ruleAction55]() {
					goto l708
				}
				depth--
				add(ruleAnd, position709)
			}
			return true
		l708:
			position, tokenIndex, depth = position708, tokenIndex708, depth708
			return false
		},
		/* 73 Equal <- <(<'='> Action56)> */
		func() bool {
			position717, tokenIndex717, depth717 := position, tokenIndex, depth
			{
				position718 := position
				depth++
				{
					position719 := position
					depth++
					if buffer[position] != rune('=') {
						goto l717
					}
					position++
					depth--
					add(rulePegText, position719)
				}
				if !_rules[ruleAction56]() {
					goto l717
				}
				depth--
				add(ruleEqual, position718)
			}
			return true
		l717:
			position, tokenIndex, depth = position717, tokenIndex717, depth717
			return false
		},
		/* 74 Less <- <(<'<'> Action57)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722 := position
					depth++
					if buffer[position] != rune('<') {
						goto l720
					}
					position++
					depth--
					add(rulePegText, position722)
				}
				if !_rules[ruleAction57]() {
					goto l720
				}
				depth--
				add(ruleLess, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 75 LessOrEqual <- <(<('<' '=')> Action58)> */
		func() bool {
			position723, tokenIndex723, depth723 := position, tokenIndex, depth
			{
				position724 := position
				depth++
				{
					position725 := position
					depth++
					if buffer[position] != rune('<') {
						goto l723
					}
					position++
					if buffer[position] != rune('=') {
						goto l723
					}
					position++
					depth--
					add(rulePegText, position725)
				}
				if !_rules[ruleAction58]() {
					goto l723
				}
				depth--
				add(ruleLessOrEqual, position724)
			}
			return true
		l723:
			position, tokenIndex, depth = position723, tokenIndex723, depth723
			return false
		},
		/* 76 Greater <- <(<'>'> Action59)> */
		func() bool {
			position726, tokenIndex726, depth726 := position, tokenIndex, depth
			{
				position727 := position
				depth++
				{
					position728 := position
					depth++
					if buffer[position] != rune('>') {
						goto l726
					}
					position++
					depth--
					add(rulePegText, position728)
				}
				if !_rules[ruleAction59]() {
					goto l726
				}
				depth--
				add(ruleGreater, position727)
			}
			return true
		l726:
			position, tokenIndex, depth = position726, tokenIndex726, depth726
			return false
		},
		/* 77 GreaterOrEqual <- <(<('>' '=')> Action60)> */
		func() bool {
			position729, tokenIndex729, depth729 := position, tokenIndex, depth
			{
				position730 := position
				depth++
				{
					position731 := position
					depth++
					if buffer[position] != rune('>') {
						goto l729
					}
					position++
					if buffer[position] != rune('=') {
						goto l729
					}
					position++
					depth--
					add(rulePegText, position731)
				}
				if !_rules[ruleAction60]() {
					goto l729
				}
				depth--
				add(ruleGreaterOrEqual, position730)
			}
			return true
		l729:
			position, tokenIndex, depth = position729, tokenIndex729, depth729
			return false
		},
		/* 78 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action61)> */
		func() bool {
			position732, tokenIndex732, depth732 := position, tokenIndex, depth
			{
				position733 := position
				depth++
				{
					position734 := position
					depth++
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l736
						}
						position++
						if buffer[position] != rune('=') {
							goto l736
						}
						position++
						goto l735
					l736:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if buffer[position] != rune('<') {
							goto l732
						}
						position++
						if buffer[position] != rune('>') {
							goto l732
						}
						position++
					}
				l735:
					depth--
					add(rulePegText, position734)
				}
				if !_rules[ruleAction61]() {
					goto l732
				}
				depth--
				add(ruleNotEqual, position733)
			}
			return true
		l732:
			position, tokenIndex, depth = position732, tokenIndex732, depth732
			return false
		},
		/* 79 Plus <- <(<'+'> Action62)> */
		func() bool {
			position737, tokenIndex737, depth737 := position, tokenIndex, depth
			{
				position738 := position
				depth++
				{
					position739 := position
					depth++
					if buffer[position] != rune('+') {
						goto l737
					}
					position++
					depth--
					add(rulePegText, position739)
				}
				if !_rules[ruleAction62]() {
					goto l737
				}
				depth--
				add(rulePlus, position738)
			}
			return true
		l737:
			position, tokenIndex, depth = position737, tokenIndex737, depth737
			return false
		},
		/* 80 Minus <- <(<'-'> Action63)> */
		func() bool {
			position740, tokenIndex740, depth740 := position, tokenIndex, depth
			{
				position741 := position
				depth++
				{
					position742 := position
					depth++
					if buffer[position] != rune('-') {
						goto l740
					}
					position++
					depth--
					add(rulePegText, position742)
				}
				if !_rules[ruleAction63]() {
					goto l740
				}
				depth--
				add(ruleMinus, position741)
			}
			return true
		l740:
			position, tokenIndex, depth = position740, tokenIndex740, depth740
			return false
		},
		/* 81 Multiply <- <(<'*'> Action64)> */
		func() bool {
			position743, tokenIndex743, depth743 := position, tokenIndex, depth
			{
				position744 := position
				depth++
				{
					position745 := position
					depth++
					if buffer[position] != rune('*') {
						goto l743
					}
					position++
					depth--
					add(rulePegText, position745)
				}
				if !_rules[ruleAction64]() {
					goto l743
				}
				depth--
				add(ruleMultiply, position744)
			}
			return true
		l743:
			position, tokenIndex, depth = position743, tokenIndex743, depth743
			return false
		},
		/* 82 Divide <- <(<'/'> Action65)> */
		func() bool {
			position746, tokenIndex746, depth746 := position, tokenIndex, depth
			{
				position747 := position
				depth++
				{
					position748 := position
					depth++
					if buffer[position] != rune('/') {
						goto l746
					}
					position++
					depth--
					add(rulePegText, position748)
				}
				if !_rules[ruleAction65]() {
					goto l746
				}
				depth--
				add(ruleDivide, position747)
			}
			return true
		l746:
			position, tokenIndex, depth = position746, tokenIndex746, depth746
			return false
		},
		/* 83 Modulo <- <(<'%'> Action66)> */
		func() bool {
			position749, tokenIndex749, depth749 := position, tokenIndex, depth
			{
				position750 := position
				depth++
				{
					position751 := position
					depth++
					if buffer[position] != rune('%') {
						goto l749
					}
					position++
					depth--
					add(rulePegText, position751)
				}
				if !_rules[ruleAction66]() {
					goto l749
				}
				depth--
				add(ruleModulo, position750)
			}
			return true
		l749:
			position, tokenIndex, depth = position749, tokenIndex749, depth749
			return false
		},
		/* 84 Identifier <- <(<ident> Action67)> */
		func() bool {
			position752, tokenIndex752, depth752 := position, tokenIndex, depth
			{
				position753 := position
				depth++
				{
					position754 := position
					depth++
					if !_rules[ruleident]() {
						goto l752
					}
					depth--
					add(rulePegText, position754)
				}
				if !_rules[ruleAction67]() {
					goto l752
				}
				depth--
				add(ruleIdentifier, position753)
			}
			return true
		l752:
			position, tokenIndex, depth = position752, tokenIndex752, depth752
			return false
		},
		/* 85 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position755, tokenIndex755, depth755 := position, tokenIndex, depth
			{
				position756 := position
				depth++
				{
					position757, tokenIndex757, depth757 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l758
					}
					position++
					goto l757
				l758:
					position, tokenIndex, depth = position757, tokenIndex757, depth757
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l755
					}
					position++
				}
			l757:
			l759:
				{
					position760, tokenIndex760, depth760 := position, tokenIndex, depth
					{
						position761, tokenIndex761, depth761 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l762
						}
						position++
						goto l761
					l762:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l763
						}
						position++
						goto l761
					l763:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l764
						}
						position++
						goto l761
					l764:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if buffer[position] != rune('_') {
							goto l760
						}
						position++
					}
				l761:
					goto l759
				l760:
					position, tokenIndex, depth = position760, tokenIndex760, depth760
				}
				depth--
				add(ruleident, position756)
			}
			return true
		l755:
			position, tokenIndex, depth = position755, tokenIndex755, depth755
			return false
		},
		/* 86 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position766 := position
				depth++
			l767:
				{
					position768, tokenIndex768, depth768 := position, tokenIndex, depth
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('\t') {
							goto l771
						}
						position++
						goto l769
					l771:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('\n') {
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
				add(rulesp, position766)
			}
			return true
		},
		/* 88 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 89 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 90 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 91 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 92 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 93 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 94 Action6 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 95 Action7 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 97 Action8 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 98 Action9 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 99 Action10 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 100 Action11 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 101 Action12 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 102 Action13 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 103 Action14 <- <{
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
		/* 104 Action15 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 105 Action16 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 106 Action17 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 107 Action18 <- <{
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
		/* 108 Action19 <- <{
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
		/* 109 Action20 <- <{
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
		/* 110 Action21 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 111 Action22 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 112 Action23 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 113 Action24 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 114 Action25 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 115 Action26 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 116 Action27 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 117 Action28 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 118 Action29 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 119 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 120 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 121 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 122 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 123 Action34 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 124 Action35 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 125 Action36 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 126 Action37 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 127 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 128 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 129 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 130 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 131 Action42 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 132 Action43 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 133 Action44 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 134 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 135 Action46 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 136 Action47 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 137 Action48 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 138 Action49 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 139 Action50 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 140 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 141 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 142 Action53 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 143 Action54 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 144 Action55 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 145 Action56 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 146 Action57 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 147 Action58 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 148 Action59 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 149 Action60 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 150 Action61 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 151 Action62 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 152 Action63 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 153 Action64 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 154 Action65 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 155 Action66 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 156 Action67 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
	}
	p.rules = _rules
}

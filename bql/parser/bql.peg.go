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
	ruleProjections
	ruleProjection
	ruleAliasExpression
	ruleWindowedFrom
	ruleDefWindowedFrom
	ruleInterval
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
	ruleIntervalUnit
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
	"Projections",
	"Projection",
	"AliasExpression",
	"WindowedFrom",
	"DefWindowedFrom",
	"Interval",
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
	"IntervalUnit",
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
	rules  [144]func() bool
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

			p.AssembleEmitter()

		case ruleAction8:

			p.AssembleProjections(begin, end)

		case ruleAction9:

			p.AssembleAlias()

		case ruleAction10:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction11:

			p.AssembleWindowedFrom(begin, end)

		case ruleAction12:

			p.AssembleInterval()

		case ruleAction13:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction14:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction15:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction16:

			p.EnsureAliasedStreamWindow()

		case ruleAction17:

			p.EnsureAliasedStreamWindow()

		case ruleAction18:

			p.AssembleAliasedStreamWindow()

		case ruleAction19:

			p.AssembleAliasedStreamWindow()

		case ruleAction20:

			p.AssembleStreamWindow()

		case ruleAction21:

			p.AssembleStreamWindow()

		case ruleAction22:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction23:

			p.AssembleSourceSinkParam()

		case ruleAction24:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction25:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction26:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction27:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction28:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction29:

			p.AssembleFuncApp()

		case ruleAction30:

			p.AssembleExpressions(begin, end)

		case ruleAction31:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction32:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction33:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction34:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction35:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction36:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction37:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction38:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction40:

			p.PushComponent(begin, end, Istream)

		case ruleAction41:

			p.PushComponent(begin, end, Dstream)

		case ruleAction42:

			p.PushComponent(begin, end, Rstream)

		case ruleAction43:

			p.PushComponent(begin, end, Tuples)

		case ruleAction44:

			p.PushComponent(begin, end, Seconds)

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamVal(substr))

		case ruleAction49:

			p.PushComponent(begin, end, Or)

		case ruleAction50:

			p.PushComponent(begin, end, And)

		case ruleAction51:

			p.PushComponent(begin, end, Equal)

		case ruleAction52:

			p.PushComponent(begin, end, Less)

		case ruleAction53:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction54:

			p.PushComponent(begin, end, Greater)

		case ruleAction55:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction56:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction57:

			p.PushComponent(begin, end, Plus)

		case ruleAction58:

			p.PushComponent(begin, end, Minus)

		case ruleAction59:

			p.PushComponent(begin, end, Multiply)

		case ruleAction60:

			p.PushComponent(begin, end, Divide)

		case ruleAction61:

			p.PushComponent(begin, end, Modulo)

		case ruleAction62:

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
		/* 2 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Projections sp DefWindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
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
		/* 9 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) Action7)> */
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
		/* 10 Projections <- <(<(Projection sp (',' sp Projection)*)> Action8)> */
		func() bool {
			position255, tokenIndex255, depth255 := position, tokenIndex, depth
			{
				position256 := position
				depth++
				{
					position257 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l255
					}
					if !_rules[rulesp]() {
						goto l255
					}
				l258:
					{
						position259, tokenIndex259, depth259 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l259
						}
						position++
						if !_rules[rulesp]() {
							goto l259
						}
						if !_rules[ruleProjection]() {
							goto l259
						}
						goto l258
					l259:
						position, tokenIndex, depth = position259, tokenIndex259, depth259
					}
					depth--
					add(rulePegText, position257)
				}
				if !_rules[ruleAction8]() {
					goto l255
				}
				depth--
				add(ruleProjections, position256)
			}
			return true
		l255:
			position, tokenIndex, depth = position255, tokenIndex255, depth255
			return false
		},
		/* 11 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position260, tokenIndex260, depth260 := position, tokenIndex, depth
			{
				position261 := position
				depth++
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l263
					}
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if !_rules[ruleExpression]() {
						goto l264
					}
					goto l262
				l264:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if !_rules[ruleWildcard]() {
						goto l260
					}
				}
			l262:
				depth--
				add(ruleProjection, position261)
			}
			return true
		l260:
			position, tokenIndex, depth = position260, tokenIndex260, depth260
			return false
		},
		/* 12 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action9)> */
		func() bool {
			position265, tokenIndex265, depth265 := position, tokenIndex, depth
			{
				position266 := position
				depth++
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l268
					}
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if !_rules[ruleWildcard]() {
						goto l265
					}
				}
			l267:
				if !_rules[rulesp]() {
					goto l265
				}
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('A') {
						goto l265
					}
					position++
				}
			l269:
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('S') {
						goto l265
					}
					position++
				}
			l271:
				if !_rules[rulesp]() {
					goto l265
				}
				if !_rules[ruleIdentifier]() {
					goto l265
				}
				if !_rules[ruleAction9]() {
					goto l265
				}
				depth--
				add(ruleAliasExpression, position266)
			}
			return true
		l265:
			position, tokenIndex, depth = position265, tokenIndex265, depth265
			return false
		},
		/* 13 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action10)> */
		func() bool {
			position273, tokenIndex273, depth273 := position, tokenIndex, depth
			{
				position274 := position
				depth++
				{
					position275 := position
					depth++
					{
						position276, tokenIndex276, depth276 := position, tokenIndex, depth
						{
							position278, tokenIndex278, depth278 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l279
							}
							position++
							goto l278
						l279:
							position, tokenIndex, depth = position278, tokenIndex278, depth278
							if buffer[position] != rune('F') {
								goto l276
							}
							position++
						}
					l278:
						{
							position280, tokenIndex280, depth280 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l281
							}
							position++
							goto l280
						l281:
							position, tokenIndex, depth = position280, tokenIndex280, depth280
							if buffer[position] != rune('R') {
								goto l276
							}
							position++
						}
					l280:
						{
							position282, tokenIndex282, depth282 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l283
							}
							position++
							goto l282
						l283:
							position, tokenIndex, depth = position282, tokenIndex282, depth282
							if buffer[position] != rune('O') {
								goto l276
							}
							position++
						}
					l282:
						{
							position284, tokenIndex284, depth284 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l285
							}
							position++
							goto l284
						l285:
							position, tokenIndex, depth = position284, tokenIndex284, depth284
							if buffer[position] != rune('M') {
								goto l276
							}
							position++
						}
					l284:
						if !_rules[rulesp]() {
							goto l276
						}
						if !_rules[ruleRelations]() {
							goto l276
						}
						if !_rules[rulesp]() {
							goto l276
						}
						goto l277
					l276:
						position, tokenIndex, depth = position276, tokenIndex276, depth276
					}
				l277:
					depth--
					add(rulePegText, position275)
				}
				if !_rules[ruleAction10]() {
					goto l273
				}
				depth--
				add(ruleWindowedFrom, position274)
			}
			return true
		l273:
			position, tokenIndex, depth = position273, tokenIndex273, depth273
			return false
		},
		/* 14 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action11)> */
		func() bool {
			position286, tokenIndex286, depth286 := position, tokenIndex, depth
			{
				position287 := position
				depth++
				{
					position288 := position
					depth++
					{
						position289, tokenIndex289, depth289 := position, tokenIndex, depth
						{
							position291, tokenIndex291, depth291 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l292
							}
							position++
							goto l291
						l292:
							position, tokenIndex, depth = position291, tokenIndex291, depth291
							if buffer[position] != rune('F') {
								goto l289
							}
							position++
						}
					l291:
						{
							position293, tokenIndex293, depth293 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l294
							}
							position++
							goto l293
						l294:
							position, tokenIndex, depth = position293, tokenIndex293, depth293
							if buffer[position] != rune('R') {
								goto l289
							}
							position++
						}
					l293:
						{
							position295, tokenIndex295, depth295 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l296
							}
							position++
							goto l295
						l296:
							position, tokenIndex, depth = position295, tokenIndex295, depth295
							if buffer[position] != rune('O') {
								goto l289
							}
							position++
						}
					l295:
						{
							position297, tokenIndex297, depth297 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l298
							}
							position++
							goto l297
						l298:
							position, tokenIndex, depth = position297, tokenIndex297, depth297
							if buffer[position] != rune('M') {
								goto l289
							}
							position++
						}
					l297:
						if !_rules[rulesp]() {
							goto l289
						}
						if !_rules[ruleDefRelations]() {
							goto l289
						}
						if !_rules[rulesp]() {
							goto l289
						}
						goto l290
					l289:
						position, tokenIndex, depth = position289, tokenIndex289, depth289
					}
				l290:
					depth--
					add(rulePegText, position288)
				}
				if !_rules[ruleAction11]() {
					goto l286
				}
				depth--
				add(ruleDefWindowedFrom, position287)
			}
			return true
		l286:
			position, tokenIndex, depth = position286, tokenIndex286, depth286
			return false
		},
		/* 15 Interval <- <(NumericLiteral sp IntervalUnit Action12)> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{
				position300 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l299
				}
				if !_rules[rulesp]() {
					goto l299
				}
				if !_rules[ruleIntervalUnit]() {
					goto l299
				}
				if !_rules[ruleAction12]() {
					goto l299
				}
				depth--
				add(ruleInterval, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 16 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position301, tokenIndex301, depth301 := position, tokenIndex, depth
			{
				position302 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l301
				}
				if !_rules[rulesp]() {
					goto l301
				}
			l303:
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l304
					}
					position++
					if !_rules[rulesp]() {
						goto l304
					}
					if !_rules[ruleRelationLike]() {
						goto l304
					}
					goto l303
				l304:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
				}
				depth--
				add(ruleRelations, position302)
			}
			return true
		l301:
			position, tokenIndex, depth = position301, tokenIndex301, depth301
			return false
		},
		/* 17 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position305, tokenIndex305, depth305 := position, tokenIndex, depth
			{
				position306 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l305
				}
				if !_rules[rulesp]() {
					goto l305
				}
			l307:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l308
					}
					position++
					if !_rules[rulesp]() {
						goto l308
					}
					if !_rules[ruleDefRelationLike]() {
						goto l308
					}
					goto l307
				l308:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
				}
				depth--
				add(ruleDefRelations, position306)
			}
			return true
		l305:
			position, tokenIndex, depth = position305, tokenIndex305, depth305
			return false
		},
		/* 18 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action13)> */
		func() bool {
			position309, tokenIndex309, depth309 := position, tokenIndex, depth
			{
				position310 := position
				depth++
				{
					position311 := position
					depth++
					{
						position312, tokenIndex312, depth312 := position, tokenIndex, depth
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('W') {
								goto l312
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('H') {
								goto l312
							}
							position++
						}
					l316:
						{
							position318, tokenIndex318, depth318 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l319
							}
							position++
							goto l318
						l319:
							position, tokenIndex, depth = position318, tokenIndex318, depth318
							if buffer[position] != rune('E') {
								goto l312
							}
							position++
						}
					l318:
						{
							position320, tokenIndex320, depth320 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l321
							}
							position++
							goto l320
						l321:
							position, tokenIndex, depth = position320, tokenIndex320, depth320
							if buffer[position] != rune('R') {
								goto l312
							}
							position++
						}
					l320:
						{
							position322, tokenIndex322, depth322 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l323
							}
							position++
							goto l322
						l323:
							position, tokenIndex, depth = position322, tokenIndex322, depth322
							if buffer[position] != rune('E') {
								goto l312
							}
							position++
						}
					l322:
						if !_rules[rulesp]() {
							goto l312
						}
						if !_rules[ruleExpression]() {
							goto l312
						}
						goto l313
					l312:
						position, tokenIndex, depth = position312, tokenIndex312, depth312
					}
				l313:
					depth--
					add(rulePegText, position311)
				}
				if !_rules[ruleAction13]() {
					goto l309
				}
				depth--
				add(ruleFilter, position310)
			}
			return true
		l309:
			position, tokenIndex, depth = position309, tokenIndex309, depth309
			return false
		},
		/* 19 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action14)> */
		func() bool {
			position324, tokenIndex324, depth324 := position, tokenIndex, depth
			{
				position325 := position
				depth++
				{
					position326 := position
					depth++
					{
						position327, tokenIndex327, depth327 := position, tokenIndex, depth
						{
							position329, tokenIndex329, depth329 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l330
							}
							position++
							goto l329
						l330:
							position, tokenIndex, depth = position329, tokenIndex329, depth329
							if buffer[position] != rune('G') {
								goto l327
							}
							position++
						}
					l329:
						{
							position331, tokenIndex331, depth331 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l332
							}
							position++
							goto l331
						l332:
							position, tokenIndex, depth = position331, tokenIndex331, depth331
							if buffer[position] != rune('R') {
								goto l327
							}
							position++
						}
					l331:
						{
							position333, tokenIndex333, depth333 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l334
							}
							position++
							goto l333
						l334:
							position, tokenIndex, depth = position333, tokenIndex333, depth333
							if buffer[position] != rune('O') {
								goto l327
							}
							position++
						}
					l333:
						{
							position335, tokenIndex335, depth335 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l336
							}
							position++
							goto l335
						l336:
							position, tokenIndex, depth = position335, tokenIndex335, depth335
							if buffer[position] != rune('U') {
								goto l327
							}
							position++
						}
					l335:
						{
							position337, tokenIndex337, depth337 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l338
							}
							position++
							goto l337
						l338:
							position, tokenIndex, depth = position337, tokenIndex337, depth337
							if buffer[position] != rune('P') {
								goto l327
							}
							position++
						}
					l337:
						if !_rules[rulesp]() {
							goto l327
						}
						{
							position339, tokenIndex339, depth339 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l340
							}
							position++
							goto l339
						l340:
							position, tokenIndex, depth = position339, tokenIndex339, depth339
							if buffer[position] != rune('B') {
								goto l327
							}
							position++
						}
					l339:
						{
							position341, tokenIndex341, depth341 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l342
							}
							position++
							goto l341
						l342:
							position, tokenIndex, depth = position341, tokenIndex341, depth341
							if buffer[position] != rune('Y') {
								goto l327
							}
							position++
						}
					l341:
						if !_rules[rulesp]() {
							goto l327
						}
						if !_rules[ruleGroupList]() {
							goto l327
						}
						goto l328
					l327:
						position, tokenIndex, depth = position327, tokenIndex327, depth327
					}
				l328:
					depth--
					add(rulePegText, position326)
				}
				if !_rules[ruleAction14]() {
					goto l324
				}
				depth--
				add(ruleGrouping, position325)
			}
			return true
		l324:
			position, tokenIndex, depth = position324, tokenIndex324, depth324
			return false
		},
		/* 20 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position343, tokenIndex343, depth343 := position, tokenIndex, depth
			{
				position344 := position
				depth++
				if !_rules[ruleExpression]() {
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
					if !_rules[ruleExpression]() {
						goto l346
					}
					goto l345
				l346:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
				}
				depth--
				add(ruleGroupList, position344)
			}
			return true
		l343:
			position, tokenIndex, depth = position343, tokenIndex343, depth343
			return false
		},
		/* 21 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action15)> */
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
							if buffer[position] != rune('h') {
								goto l353
							}
							position++
							goto l352
						l353:
							position, tokenIndex, depth = position352, tokenIndex352, depth352
							if buffer[position] != rune('H') {
								goto l350
							}
							position++
						}
					l352:
						{
							position354, tokenIndex354, depth354 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l355
							}
							position++
							goto l354
						l355:
							position, tokenIndex, depth = position354, tokenIndex354, depth354
							if buffer[position] != rune('A') {
								goto l350
							}
							position++
						}
					l354:
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l357
							}
							position++
							goto l356
						l357:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
							if buffer[position] != rune('V') {
								goto l350
							}
							position++
						}
					l356:
						{
							position358, tokenIndex358, depth358 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex, depth = position358, tokenIndex358, depth358
							if buffer[position] != rune('I') {
								goto l350
							}
							position++
						}
					l358:
						{
							position360, tokenIndex360, depth360 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l361
							}
							position++
							goto l360
						l361:
							position, tokenIndex, depth = position360, tokenIndex360, depth360
							if buffer[position] != rune('N') {
								goto l350
							}
							position++
						}
					l360:
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('G') {
								goto l350
							}
							position++
						}
					l362:
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
				if !_rules[ruleAction15]() {
					goto l347
				}
				depth--
				add(ruleHaving, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 22 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action16))> */
		func() bool {
			position364, tokenIndex364, depth364 := position, tokenIndex, depth
			{
				position365 := position
				depth++
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l367
					}
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if !_rules[ruleStreamWindow]() {
						goto l364
					}
					if !_rules[ruleAction16]() {
						goto l364
					}
				}
			l366:
				depth--
				add(ruleRelationLike, position365)
			}
			return true
		l364:
			position, tokenIndex, depth = position364, tokenIndex364, depth364
			return false
		},
		/* 23 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action17))> */
		func() bool {
			position368, tokenIndex368, depth368 := position, tokenIndex, depth
			{
				position369 := position
				depth++
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l371
					}
					goto l370
				l371:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
					if !_rules[ruleDefStreamWindow]() {
						goto l368
					}
					if !_rules[ruleAction17]() {
						goto l368
					}
				}
			l370:
				depth--
				add(ruleDefRelationLike, position369)
			}
			return true
		l368:
			position, tokenIndex, depth = position368, tokenIndex368, depth368
			return false
		},
		/* 24 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action18)> */
		func() bool {
			position372, tokenIndex372, depth372 := position, tokenIndex, depth
			{
				position373 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l372
				}
				if !_rules[rulesp]() {
					goto l372
				}
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('A') {
						goto l372
					}
					position++
				}
			l374:
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l377
					}
					position++
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					if buffer[position] != rune('S') {
						goto l372
					}
					position++
				}
			l376:
				if !_rules[rulesp]() {
					goto l372
				}
				if !_rules[ruleIdentifier]() {
					goto l372
				}
				if !_rules[ruleAction18]() {
					goto l372
				}
				depth--
				add(ruleAliasedStreamWindow, position373)
			}
			return true
		l372:
			position, tokenIndex, depth = position372, tokenIndex372, depth372
			return false
		},
		/* 25 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action19)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l378
				}
				if !_rules[rulesp]() {
					goto l378
				}
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('A') {
						goto l378
					}
					position++
				}
			l380:
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l383
					}
					position++
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if buffer[position] != rune('S') {
						goto l378
					}
					position++
				}
			l382:
				if !_rules[rulesp]() {
					goto l378
				}
				if !_rules[ruleIdentifier]() {
					goto l378
				}
				if !_rules[ruleAction19]() {
					goto l378
				}
				depth--
				add(ruleDefAliasedStreamWindow, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 26 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action20)> */
		func() bool {
			position384, tokenIndex384, depth384 := position, tokenIndex, depth
			{
				position385 := position
				depth++
				if !_rules[ruleStream]() {
					goto l384
				}
				if !_rules[rulesp]() {
					goto l384
				}
				if buffer[position] != rune('[') {
					goto l384
				}
				position++
				if !_rules[rulesp]() {
					goto l384
				}
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('R') {
						goto l384
					}
					position++
				}
			l386:
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l389
					}
					position++
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					if buffer[position] != rune('A') {
						goto l384
					}
					position++
				}
			l388:
				{
					position390, tokenIndex390, depth390 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l391
					}
					position++
					goto l390
				l391:
					position, tokenIndex, depth = position390, tokenIndex390, depth390
					if buffer[position] != rune('N') {
						goto l384
					}
					position++
				}
			l390:
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if buffer[position] != rune('G') {
						goto l384
					}
					position++
				}
			l392:
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l395
					}
					position++
					goto l394
				l395:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
					if buffer[position] != rune('E') {
						goto l384
					}
					position++
				}
			l394:
				if !_rules[rulesp]() {
					goto l384
				}
				if !_rules[ruleInterval]() {
					goto l384
				}
				if !_rules[rulesp]() {
					goto l384
				}
				if buffer[position] != rune(']') {
					goto l384
				}
				position++
				if !_rules[ruleAction20]() {
					goto l384
				}
				depth--
				add(ruleStreamWindow, position385)
			}
			return true
		l384:
			position, tokenIndex, depth = position384, tokenIndex384, depth384
			return false
		},
		/* 27 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action21)> */
		func() bool {
			position396, tokenIndex396, depth396 := position, tokenIndex, depth
			{
				position397 := position
				depth++
				if !_rules[ruleStream]() {
					goto l396
				}
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l398
					}
					if buffer[position] != rune('[') {
						goto l398
					}
					position++
					if !_rules[rulesp]() {
						goto l398
					}
					{
						position400, tokenIndex400, depth400 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l401
						}
						position++
						goto l400
					l401:
						position, tokenIndex, depth = position400, tokenIndex400, depth400
						if buffer[position] != rune('R') {
							goto l398
						}
						position++
					}
				l400:
					{
						position402, tokenIndex402, depth402 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l403
						}
						position++
						goto l402
					l403:
						position, tokenIndex, depth = position402, tokenIndex402, depth402
						if buffer[position] != rune('A') {
							goto l398
						}
						position++
					}
				l402:
					{
						position404, tokenIndex404, depth404 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l405
						}
						position++
						goto l404
					l405:
						position, tokenIndex, depth = position404, tokenIndex404, depth404
						if buffer[position] != rune('N') {
							goto l398
						}
						position++
					}
				l404:
					{
						position406, tokenIndex406, depth406 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l407
						}
						position++
						goto l406
					l407:
						position, tokenIndex, depth = position406, tokenIndex406, depth406
						if buffer[position] != rune('G') {
							goto l398
						}
						position++
					}
				l406:
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l409
						}
						position++
						goto l408
					l409:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
						if buffer[position] != rune('E') {
							goto l398
						}
						position++
					}
				l408:
					if !_rules[rulesp]() {
						goto l398
					}
					if !_rules[ruleInterval]() {
						goto l398
					}
					if !_rules[rulesp]() {
						goto l398
					}
					if buffer[position] != rune(']') {
						goto l398
					}
					position++
					goto l399
				l398:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
				}
			l399:
				if !_rules[ruleAction21]() {
					goto l396
				}
				depth--
				add(ruleDefStreamWindow, position397)
			}
			return true
		l396:
			position, tokenIndex, depth = position396, tokenIndex396, depth396
			return false
		},
		/* 28 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action22)> */
		func() bool {
			position410, tokenIndex410, depth410 := position, tokenIndex, depth
			{
				position411 := position
				depth++
				{
					position412 := position
					depth++
					{
						position413, tokenIndex413, depth413 := position, tokenIndex, depth
						{
							position415, tokenIndex415, depth415 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l416
							}
							position++
							goto l415
						l416:
							position, tokenIndex, depth = position415, tokenIndex415, depth415
							if buffer[position] != rune('W') {
								goto l413
							}
							position++
						}
					l415:
						{
							position417, tokenIndex417, depth417 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l418
							}
							position++
							goto l417
						l418:
							position, tokenIndex, depth = position417, tokenIndex417, depth417
							if buffer[position] != rune('I') {
								goto l413
							}
							position++
						}
					l417:
						{
							position419, tokenIndex419, depth419 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l420
							}
							position++
							goto l419
						l420:
							position, tokenIndex, depth = position419, tokenIndex419, depth419
							if buffer[position] != rune('T') {
								goto l413
							}
							position++
						}
					l419:
						{
							position421, tokenIndex421, depth421 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l422
							}
							position++
							goto l421
						l422:
							position, tokenIndex, depth = position421, tokenIndex421, depth421
							if buffer[position] != rune('H') {
								goto l413
							}
							position++
						}
					l421:
						if !_rules[rulesp]() {
							goto l413
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l413
						}
						if !_rules[rulesp]() {
							goto l413
						}
					l423:
						{
							position424, tokenIndex424, depth424 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l424
							}
							position++
							if !_rules[rulesp]() {
								goto l424
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l424
							}
							goto l423
						l424:
							position, tokenIndex, depth = position424, tokenIndex424, depth424
						}
						goto l414
					l413:
						position, tokenIndex, depth = position413, tokenIndex413, depth413
					}
				l414:
					depth--
					add(rulePegText, position412)
				}
				if !_rules[ruleAction22]() {
					goto l410
				}
				depth--
				add(ruleSourceSinkSpecs, position411)
			}
			return true
		l410:
			position, tokenIndex, depth = position410, tokenIndex410, depth410
			return false
		},
		/* 29 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action23)> */
		func() bool {
			position425, tokenIndex425, depth425 := position, tokenIndex, depth
			{
				position426 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l425
				}
				if buffer[position] != rune('=') {
					goto l425
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l425
				}
				if !_rules[ruleAction23]() {
					goto l425
				}
				depth--
				add(ruleSourceSinkParam, position426)
			}
			return true
		l425:
			position, tokenIndex, depth = position425, tokenIndex425, depth425
			return false
		},
		/* 30 Expression <- <orExpr> */
		func() bool {
			position427, tokenIndex427, depth427 := position, tokenIndex, depth
			{
				position428 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l427
				}
				depth--
				add(ruleExpression, position428)
			}
			return true
		l427:
			position, tokenIndex, depth = position427, tokenIndex427, depth427
			return false
		},
		/* 31 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action24)> */
		func() bool {
			position429, tokenIndex429, depth429 := position, tokenIndex, depth
			{
				position430 := position
				depth++
				{
					position431 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l429
					}
					if !_rules[rulesp]() {
						goto l429
					}
					{
						position432, tokenIndex432, depth432 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l432
						}
						if !_rules[rulesp]() {
							goto l432
						}
						if !_rules[ruleandExpr]() {
							goto l432
						}
						goto l433
					l432:
						position, tokenIndex, depth = position432, tokenIndex432, depth432
					}
				l433:
					depth--
					add(rulePegText, position431)
				}
				if !_rules[ruleAction24]() {
					goto l429
				}
				depth--
				add(ruleorExpr, position430)
			}
			return true
		l429:
			position, tokenIndex, depth = position429, tokenIndex429, depth429
			return false
		},
		/* 32 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action25)> */
		func() bool {
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
			{
				position435 := position
				depth++
				{
					position436 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l434
					}
					if !_rules[rulesp]() {
						goto l434
					}
					{
						position437, tokenIndex437, depth437 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l437
						}
						if !_rules[rulesp]() {
							goto l437
						}
						if !_rules[rulecomparisonExpr]() {
							goto l437
						}
						goto l438
					l437:
						position, tokenIndex, depth = position437, tokenIndex437, depth437
					}
				l438:
					depth--
					add(rulePegText, position436)
				}
				if !_rules[ruleAction25]() {
					goto l434
				}
				depth--
				add(ruleandExpr, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
			return false
		},
		/* 33 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action26)> */
		func() bool {
			position439, tokenIndex439, depth439 := position, tokenIndex, depth
			{
				position440 := position
				depth++
				{
					position441 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l439
					}
					if !_rules[rulesp]() {
						goto l439
					}
					{
						position442, tokenIndex442, depth442 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l442
						}
						if !_rules[rulesp]() {
							goto l442
						}
						if !_rules[ruletermExpr]() {
							goto l442
						}
						goto l443
					l442:
						position, tokenIndex, depth = position442, tokenIndex442, depth442
					}
				l443:
					depth--
					add(rulePegText, position441)
				}
				if !_rules[ruleAction26]() {
					goto l439
				}
				depth--
				add(rulecomparisonExpr, position440)
			}
			return true
		l439:
			position, tokenIndex, depth = position439, tokenIndex439, depth439
			return false
		},
		/* 34 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action27)> */
		func() bool {
			position444, tokenIndex444, depth444 := position, tokenIndex, depth
			{
				position445 := position
				depth++
				{
					position446 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l444
					}
					if !_rules[rulesp]() {
						goto l444
					}
					{
						position447, tokenIndex447, depth447 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l447
						}
						if !_rules[rulesp]() {
							goto l447
						}
						if !_rules[ruleproductExpr]() {
							goto l447
						}
						goto l448
					l447:
						position, tokenIndex, depth = position447, tokenIndex447, depth447
					}
				l448:
					depth--
					add(rulePegText, position446)
				}
				if !_rules[ruleAction27]() {
					goto l444
				}
				depth--
				add(ruletermExpr, position445)
			}
			return true
		l444:
			position, tokenIndex, depth = position444, tokenIndex444, depth444
			return false
		},
		/* 35 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action28)> */
		func() bool {
			position449, tokenIndex449, depth449 := position, tokenIndex, depth
			{
				position450 := position
				depth++
				{
					position451 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l449
					}
					if !_rules[rulesp]() {
						goto l449
					}
					{
						position452, tokenIndex452, depth452 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l452
						}
						if !_rules[rulesp]() {
							goto l452
						}
						if !_rules[rulebaseExpr]() {
							goto l452
						}
						goto l453
					l452:
						position, tokenIndex, depth = position452, tokenIndex452, depth452
					}
				l453:
					depth--
					add(rulePegText, position451)
				}
				if !_rules[ruleAction28]() {
					goto l449
				}
				depth--
				add(ruleproductExpr, position450)
			}
			return true
		l449:
			position, tokenIndex, depth = position449, tokenIndex449, depth449
			return false
		},
		/* 36 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position454, tokenIndex454, depth454 := position, tokenIndex, depth
			{
				position455 := position
				depth++
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l457
					}
					position++
					if !_rules[rulesp]() {
						goto l457
					}
					if !_rules[ruleExpression]() {
						goto l457
					}
					if !_rules[rulesp]() {
						goto l457
					}
					if buffer[position] != rune(')') {
						goto l457
					}
					position++
					goto l456
				l457:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if !_rules[ruleBooleanLiteral]() {
						goto l458
					}
					goto l456
				l458:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if !_rules[ruleFuncApp]() {
						goto l459
					}
					goto l456
				l459:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if !_rules[ruleRowValue]() {
						goto l460
					}
					goto l456
				l460:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if !_rules[ruleLiteral]() {
						goto l454
					}
				}
			l456:
				depth--
				add(rulebaseExpr, position455)
			}
			return true
		l454:
			position, tokenIndex, depth = position454, tokenIndex454, depth454
			return false
		},
		/* 37 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action29)> */
		func() bool {
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l461
				}
				if !_rules[rulesp]() {
					goto l461
				}
				if buffer[position] != rune('(') {
					goto l461
				}
				position++
				if !_rules[rulesp]() {
					goto l461
				}
				if !_rules[ruleFuncParams]() {
					goto l461
				}
				if !_rules[rulesp]() {
					goto l461
				}
				if buffer[position] != rune(')') {
					goto l461
				}
				position++
				if !_rules[ruleAction29]() {
					goto l461
				}
				depth--
				add(ruleFuncApp, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 38 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action30)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				{
					position465 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l463
					}
					if !_rules[rulesp]() {
						goto l463
					}
				l466:
					{
						position467, tokenIndex467, depth467 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l467
						}
						position++
						if !_rules[rulesp]() {
							goto l467
						}
						if !_rules[ruleExpression]() {
							goto l467
						}
						goto l466
					l467:
						position, tokenIndex, depth = position467, tokenIndex467, depth467
					}
					depth--
					add(rulePegText, position465)
				}
				if !_rules[ruleAction30]() {
					goto l463
				}
				depth--
				add(ruleFuncParams, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
			return false
		},
		/* 39 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position468, tokenIndex468, depth468 := position, tokenIndex, depth
			{
				position469 := position
				depth++
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l471
					}
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if !_rules[ruleNumericLiteral]() {
						goto l472
					}
					goto l470
				l472:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if !_rules[ruleStringLiteral]() {
						goto l468
					}
				}
			l470:
				depth--
				add(ruleLiteral, position469)
			}
			return true
		l468:
			position, tokenIndex, depth = position468, tokenIndex468, depth468
			return false
		},
		/* 40 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position473, tokenIndex473, depth473 := position, tokenIndex, depth
			{
				position474 := position
				depth++
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l476
					}
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !_rules[ruleNotEqual]() {
						goto l477
					}
					goto l475
				l477:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !_rules[ruleLessOrEqual]() {
						goto l478
					}
					goto l475
				l478:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !_rules[ruleLess]() {
						goto l479
					}
					goto l475
				l479:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !_rules[ruleGreaterOrEqual]() {
						goto l480
					}
					goto l475
				l480:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !_rules[ruleGreater]() {
						goto l481
					}
					goto l475
				l481:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !_rules[ruleNotEqual]() {
						goto l473
					}
				}
			l475:
				depth--
				add(ruleComparisonOp, position474)
			}
			return true
		l473:
			position, tokenIndex, depth = position473, tokenIndex473, depth473
			return false
		},
		/* 41 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position482, tokenIndex482, depth482 := position, tokenIndex, depth
			{
				position483 := position
				depth++
				{
					position484, tokenIndex484, depth484 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l485
					}
					goto l484
				l485:
					position, tokenIndex, depth = position484, tokenIndex484, depth484
					if !_rules[ruleMinus]() {
						goto l482
					}
				}
			l484:
				depth--
				add(rulePlusMinusOp, position483)
			}
			return true
		l482:
			position, tokenIndex, depth = position482, tokenIndex482, depth482
			return false
		},
		/* 42 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position486, tokenIndex486, depth486 := position, tokenIndex, depth
			{
				position487 := position
				depth++
				{
					position488, tokenIndex488, depth488 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l489
					}
					goto l488
				l489:
					position, tokenIndex, depth = position488, tokenIndex488, depth488
					if !_rules[ruleDivide]() {
						goto l490
					}
					goto l488
				l490:
					position, tokenIndex, depth = position488, tokenIndex488, depth488
					if !_rules[ruleModulo]() {
						goto l486
					}
				}
			l488:
				depth--
				add(ruleMultDivOp, position487)
			}
			return true
		l486:
			position, tokenIndex, depth = position486, tokenIndex486, depth486
			return false
		},
		/* 43 Stream <- <(<ident> Action31)> */
		func() bool {
			position491, tokenIndex491, depth491 := position, tokenIndex, depth
			{
				position492 := position
				depth++
				{
					position493 := position
					depth++
					if !_rules[ruleident]() {
						goto l491
					}
					depth--
					add(rulePegText, position493)
				}
				if !_rules[ruleAction31]() {
					goto l491
				}
				depth--
				add(ruleStream, position492)
			}
			return true
		l491:
			position, tokenIndex, depth = position491, tokenIndex491, depth491
			return false
		},
		/* 44 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action32)> */
		func() bool {
			position494, tokenIndex494, depth494 := position, tokenIndex, depth
			{
				position495 := position
				depth++
				{
					position496 := position
					depth++
					{
						position497, tokenIndex497, depth497 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l497
						}
						if buffer[position] != rune(':') {
							goto l497
						}
						position++
						goto l498
					l497:
						position, tokenIndex, depth = position497, tokenIndex497, depth497
					}
				l498:
					{
						position499, tokenIndex499, depth499 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l500
						}
						position++
						goto l499
					l500:
						position, tokenIndex, depth = position499, tokenIndex499, depth499
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l494
						}
						position++
					}
				l499:
				l501:
					{
						position502, tokenIndex502, depth502 := position, tokenIndex, depth
						{
							position503, tokenIndex503, depth503 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l504
							}
							position++
							goto l503
						l504:
							position, tokenIndex, depth = position503, tokenIndex503, depth503
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l505
							}
							position++
							goto l503
						l505:
							position, tokenIndex, depth = position503, tokenIndex503, depth503
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l506
							}
							position++
							goto l503
						l506:
							position, tokenIndex, depth = position503, tokenIndex503, depth503
							if buffer[position] != rune('_') {
								goto l507
							}
							position++
							goto l503
						l507:
							position, tokenIndex, depth = position503, tokenIndex503, depth503
							if buffer[position] != rune('.') {
								goto l502
							}
							position++
						}
					l503:
						goto l501
					l502:
						position, tokenIndex, depth = position502, tokenIndex502, depth502
					}
					depth--
					add(rulePegText, position496)
				}
				if !_rules[ruleAction32]() {
					goto l494
				}
				depth--
				add(ruleRowValue, position495)
			}
			return true
		l494:
			position, tokenIndex, depth = position494, tokenIndex494, depth494
			return false
		},
		/* 45 NumericLiteral <- <(<('-'? [0-9]+)> Action33)> */
		func() bool {
			position508, tokenIndex508, depth508 := position, tokenIndex, depth
			{
				position509 := position
				depth++
				{
					position510 := position
					depth++
					{
						position511, tokenIndex511, depth511 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l511
						}
						position++
						goto l512
					l511:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
					}
				l512:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l508
					}
					position++
				l513:
					{
						position514, tokenIndex514, depth514 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l514
						}
						position++
						goto l513
					l514:
						position, tokenIndex, depth = position514, tokenIndex514, depth514
					}
					depth--
					add(rulePegText, position510)
				}
				if !_rules[ruleAction33]() {
					goto l508
				}
				depth--
				add(ruleNumericLiteral, position509)
			}
			return true
		l508:
			position, tokenIndex, depth = position508, tokenIndex508, depth508
			return false
		},
		/* 46 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action34)> */
		func() bool {
			position515, tokenIndex515, depth515 := position, tokenIndex, depth
			{
				position516 := position
				depth++
				{
					position517 := position
					depth++
					{
						position518, tokenIndex518, depth518 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l518
						}
						position++
						goto l519
					l518:
						position, tokenIndex, depth = position518, tokenIndex518, depth518
					}
				l519:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l515
					}
					position++
				l520:
					{
						position521, tokenIndex521, depth521 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l521
						}
						position++
						goto l520
					l521:
						position, tokenIndex, depth = position521, tokenIndex521, depth521
					}
					if buffer[position] != rune('.') {
						goto l515
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l515
					}
					position++
				l522:
					{
						position523, tokenIndex523, depth523 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l523
						}
						position++
						goto l522
					l523:
						position, tokenIndex, depth = position523, tokenIndex523, depth523
					}
					depth--
					add(rulePegText, position517)
				}
				if !_rules[ruleAction34]() {
					goto l515
				}
				depth--
				add(ruleFloatLiteral, position516)
			}
			return true
		l515:
			position, tokenIndex, depth = position515, tokenIndex515, depth515
			return false
		},
		/* 47 Function <- <(<ident> Action35)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				{
					position526 := position
					depth++
					if !_rules[ruleident]() {
						goto l524
					}
					depth--
					add(rulePegText, position526)
				}
				if !_rules[ruleAction35]() {
					goto l524
				}
				depth--
				add(ruleFunction, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 48 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position527, tokenIndex527, depth527 := position, tokenIndex, depth
			{
				position528 := position
				depth++
				{
					position529, tokenIndex529, depth529 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l530
					}
					goto l529
				l530:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if !_rules[ruleFALSE]() {
						goto l527
					}
				}
			l529:
				depth--
				add(ruleBooleanLiteral, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 49 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action36)> */
		func() bool {
			position531, tokenIndex531, depth531 := position, tokenIndex, depth
			{
				position532 := position
				depth++
				{
					position533 := position
					depth++
					{
						position534, tokenIndex534, depth534 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l535
						}
						position++
						goto l534
					l535:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if buffer[position] != rune('T') {
							goto l531
						}
						position++
					}
				l534:
					{
						position536, tokenIndex536, depth536 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l537
						}
						position++
						goto l536
					l537:
						position, tokenIndex, depth = position536, tokenIndex536, depth536
						if buffer[position] != rune('R') {
							goto l531
						}
						position++
					}
				l536:
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l539
						}
						position++
						goto l538
					l539:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
						if buffer[position] != rune('U') {
							goto l531
						}
						position++
					}
				l538:
					{
						position540, tokenIndex540, depth540 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l541
						}
						position++
						goto l540
					l541:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
						if buffer[position] != rune('E') {
							goto l531
						}
						position++
					}
				l540:
					depth--
					add(rulePegText, position533)
				}
				if !_rules[ruleAction36]() {
					goto l531
				}
				depth--
				add(ruleTRUE, position532)
			}
			return true
		l531:
			position, tokenIndex, depth = position531, tokenIndex531, depth531
			return false
		},
		/* 50 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action37)> */
		func() bool {
			position542, tokenIndex542, depth542 := position, tokenIndex, depth
			{
				position543 := position
				depth++
				{
					position544 := position
					depth++
					{
						position545, tokenIndex545, depth545 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l546
						}
						position++
						goto l545
					l546:
						position, tokenIndex, depth = position545, tokenIndex545, depth545
						if buffer[position] != rune('F') {
							goto l542
						}
						position++
					}
				l545:
					{
						position547, tokenIndex547, depth547 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l548
						}
						position++
						goto l547
					l548:
						position, tokenIndex, depth = position547, tokenIndex547, depth547
						if buffer[position] != rune('A') {
							goto l542
						}
						position++
					}
				l547:
					{
						position549, tokenIndex549, depth549 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l550
						}
						position++
						goto l549
					l550:
						position, tokenIndex, depth = position549, tokenIndex549, depth549
						if buffer[position] != rune('L') {
							goto l542
						}
						position++
					}
				l549:
					{
						position551, tokenIndex551, depth551 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l552
						}
						position++
						goto l551
					l552:
						position, tokenIndex, depth = position551, tokenIndex551, depth551
						if buffer[position] != rune('S') {
							goto l542
						}
						position++
					}
				l551:
					{
						position553, tokenIndex553, depth553 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l554
						}
						position++
						goto l553
					l554:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						if buffer[position] != rune('E') {
							goto l542
						}
						position++
					}
				l553:
					depth--
					add(rulePegText, position544)
				}
				if !_rules[ruleAction37]() {
					goto l542
				}
				depth--
				add(ruleFALSE, position543)
			}
			return true
		l542:
			position, tokenIndex, depth = position542, tokenIndex542, depth542
			return false
		},
		/* 51 Wildcard <- <(<'*'> Action38)> */
		func() bool {
			position555, tokenIndex555, depth555 := position, tokenIndex, depth
			{
				position556 := position
				depth++
				{
					position557 := position
					depth++
					if buffer[position] != rune('*') {
						goto l555
					}
					position++
					depth--
					add(rulePegText, position557)
				}
				if !_rules[ruleAction38]() {
					goto l555
				}
				depth--
				add(ruleWildcard, position556)
			}
			return true
		l555:
			position, tokenIndex, depth = position555, tokenIndex555, depth555
			return false
		},
		/* 52 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action39)> */
		func() bool {
			position558, tokenIndex558, depth558 := position, tokenIndex, depth
			{
				position559 := position
				depth++
				{
					position560 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l558
					}
					position++
				l561:
					{
						position562, tokenIndex562, depth562 := position, tokenIndex, depth
						{
							position563, tokenIndex563, depth563 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l564
							}
							position++
							if buffer[position] != rune('\'') {
								goto l564
							}
							position++
							goto l563
						l564:
							position, tokenIndex, depth = position563, tokenIndex563, depth563
							{
								position565, tokenIndex565, depth565 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l565
								}
								position++
								goto l562
							l565:
								position, tokenIndex, depth = position565, tokenIndex565, depth565
							}
							if !matchDot() {
								goto l562
							}
						}
					l563:
						goto l561
					l562:
						position, tokenIndex, depth = position562, tokenIndex562, depth562
					}
					if buffer[position] != rune('\'') {
						goto l558
					}
					position++
					depth--
					add(rulePegText, position560)
				}
				if !_rules[ruleAction39]() {
					goto l558
				}
				depth--
				add(ruleStringLiteral, position559)
			}
			return true
		l558:
			position, tokenIndex, depth = position558, tokenIndex558, depth558
			return false
		},
		/* 53 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action40)> */
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
						if buffer[position] != rune('i') {
							goto l570
						}
						position++
						goto l569
					l570:
						position, tokenIndex, depth = position569, tokenIndex569, depth569
						if buffer[position] != rune('I') {
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
				if !_rules[ruleAction40]() {
					goto l566
				}
				depth--
				add(ruleISTREAM, position567)
			}
			return true
		l566:
			position, tokenIndex, depth = position566, tokenIndex566, depth566
			return false
		},
		/* 54 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action41)> */
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
						if buffer[position] != rune('d') {
							goto l587
						}
						position++
						goto l586
					l587:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
						if buffer[position] != rune('D') {
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
				if !_rules[ruleAction41]() {
					goto l583
				}
				depth--
				add(ruleDSTREAM, position584)
			}
			return true
		l583:
			position, tokenIndex, depth = position583, tokenIndex583, depth583
			return false
		},
		/* 55 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action42)> */
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
						if buffer[position] != rune('r') {
							goto l604
						}
						position++
						goto l603
					l604:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
						if buffer[position] != rune('R') {
							goto l600
						}
						position++
					}
				l603:
					{
						position605, tokenIndex605, depth605 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l606
						}
						position++
						goto l605
					l606:
						position, tokenIndex, depth = position605, tokenIndex605, depth605
						if buffer[position] != rune('S') {
							goto l600
						}
						position++
					}
				l605:
					{
						position607, tokenIndex607, depth607 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l608
						}
						position++
						goto l607
					l608:
						position, tokenIndex, depth = position607, tokenIndex607, depth607
						if buffer[position] != rune('T') {
							goto l600
						}
						position++
					}
				l607:
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('R') {
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
						if buffer[position] != rune('a') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('A') {
							goto l600
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('M') {
							goto l600
						}
						position++
					}
				l615:
					depth--
					add(rulePegText, position602)
				}
				if !_rules[ruleAction42]() {
					goto l600
				}
				depth--
				add(ruleRSTREAM, position601)
			}
			return true
		l600:
			position, tokenIndex, depth = position600, tokenIndex600, depth600
			return false
		},
		/* 56 IntervalUnit <- <(TUPLES / SECONDS)> */
		func() bool {
			position617, tokenIndex617, depth617 := position, tokenIndex, depth
			{
				position618 := position
				depth++
				{
					position619, tokenIndex619, depth619 := position, tokenIndex, depth
					if !_rules[ruleTUPLES]() {
						goto l620
					}
					goto l619
				l620:
					position, tokenIndex, depth = position619, tokenIndex619, depth619
					if !_rules[ruleSECONDS]() {
						goto l617
					}
				}
			l619:
				depth--
				add(ruleIntervalUnit, position618)
			}
			return true
		l617:
			position, tokenIndex, depth = position617, tokenIndex617, depth617
			return false
		},
		/* 57 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action43)> */
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
						if buffer[position] != rune('t') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('T') {
							goto l621
						}
						position++
					}
				l624:
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('U') {
							goto l621
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('P') {
							goto l621
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('L') {
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
						if buffer[position] != rune('s') {
							goto l635
						}
						position++
						goto l634
					l635:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
						if buffer[position] != rune('S') {
							goto l621
						}
						position++
					}
				l634:
					depth--
					add(rulePegText, position623)
				}
				if !_rules[ruleAction43]() {
					goto l621
				}
				depth--
				add(ruleTUPLES, position622)
			}
			return true
		l621:
			position, tokenIndex, depth = position621, tokenIndex621, depth621
			return false
		},
		/* 58 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action44)> */
		func() bool {
			position636, tokenIndex636, depth636 := position, tokenIndex, depth
			{
				position637 := position
				depth++
				{
					position638 := position
					depth++
					{
						position639, tokenIndex639, depth639 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l640
						}
						position++
						goto l639
					l640:
						position, tokenIndex, depth = position639, tokenIndex639, depth639
						if buffer[position] != rune('S') {
							goto l636
						}
						position++
					}
				l639:
					{
						position641, tokenIndex641, depth641 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l642
						}
						position++
						goto l641
					l642:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
						if buffer[position] != rune('E') {
							goto l636
						}
						position++
					}
				l641:
					{
						position643, tokenIndex643, depth643 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l644
						}
						position++
						goto l643
					l644:
						position, tokenIndex, depth = position643, tokenIndex643, depth643
						if buffer[position] != rune('C') {
							goto l636
						}
						position++
					}
				l643:
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l646
						}
						position++
						goto l645
					l646:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
						if buffer[position] != rune('O') {
							goto l636
						}
						position++
					}
				l645:
					{
						position647, tokenIndex647, depth647 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l648
						}
						position++
						goto l647
					l648:
						position, tokenIndex, depth = position647, tokenIndex647, depth647
						if buffer[position] != rune('N') {
							goto l636
						}
						position++
					}
				l647:
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l650
						}
						position++
						goto l649
					l650:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						if buffer[position] != rune('D') {
							goto l636
						}
						position++
					}
				l649:
					{
						position651, tokenIndex651, depth651 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l652
						}
						position++
						goto l651
					l652:
						position, tokenIndex, depth = position651, tokenIndex651, depth651
						if buffer[position] != rune('S') {
							goto l636
						}
						position++
					}
				l651:
					depth--
					add(rulePegText, position638)
				}
				if !_rules[ruleAction44]() {
					goto l636
				}
				depth--
				add(ruleSECONDS, position637)
			}
			return true
		l636:
			position, tokenIndex, depth = position636, tokenIndex636, depth636
			return false
		},
		/* 59 StreamIdentifier <- <(<ident> Action45)> */
		func() bool {
			position653, tokenIndex653, depth653 := position, tokenIndex, depth
			{
				position654 := position
				depth++
				{
					position655 := position
					depth++
					if !_rules[ruleident]() {
						goto l653
					}
					depth--
					add(rulePegText, position655)
				}
				if !_rules[ruleAction45]() {
					goto l653
				}
				depth--
				add(ruleStreamIdentifier, position654)
			}
			return true
		l653:
			position, tokenIndex, depth = position653, tokenIndex653, depth653
			return false
		},
		/* 60 SourceSinkType <- <(<ident> Action46)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				{
					position658 := position
					depth++
					if !_rules[ruleident]() {
						goto l656
					}
					depth--
					add(rulePegText, position658)
				}
				if !_rules[ruleAction46]() {
					goto l656
				}
				depth--
				add(ruleSourceSinkType, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 61 SourceSinkParamKey <- <(<ident> Action47)> */
		func() bool {
			position659, tokenIndex659, depth659 := position, tokenIndex, depth
			{
				position660 := position
				depth++
				{
					position661 := position
					depth++
					if !_rules[ruleident]() {
						goto l659
					}
					depth--
					add(rulePegText, position661)
				}
				if !_rules[ruleAction47]() {
					goto l659
				}
				depth--
				add(ruleSourceSinkParamKey, position660)
			}
			return true
		l659:
			position, tokenIndex, depth = position659, tokenIndex659, depth659
			return false
		},
		/* 62 SourceSinkParamVal <- <(<([a-z] / [A-Z] / [0-9] / '_')+> Action48)> */
		func() bool {
			position662, tokenIndex662, depth662 := position, tokenIndex, depth
			{
				position663 := position
				depth++
				{
					position664 := position
					depth++
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l669
						}
						position++
						goto l667
					l669:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l670
						}
						position++
						goto l667
					l670:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if buffer[position] != rune('_') {
							goto l662
						}
						position++
					}
				l667:
				l665:
					{
						position666, tokenIndex666, depth666 := position, tokenIndex, depth
						{
							position671, tokenIndex671, depth671 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l672
							}
							position++
							goto l671
						l672:
							position, tokenIndex, depth = position671, tokenIndex671, depth671
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l673
							}
							position++
							goto l671
						l673:
							position, tokenIndex, depth = position671, tokenIndex671, depth671
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l674
							}
							position++
							goto l671
						l674:
							position, tokenIndex, depth = position671, tokenIndex671, depth671
							if buffer[position] != rune('_') {
								goto l666
							}
							position++
						}
					l671:
						goto l665
					l666:
						position, tokenIndex, depth = position666, tokenIndex666, depth666
					}
					depth--
					add(rulePegText, position664)
				}
				if !_rules[ruleAction48]() {
					goto l662
				}
				depth--
				add(ruleSourceSinkParamVal, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 63 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action49)> */
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
						if buffer[position] != rune('o') {
							goto l679
						}
						position++
						goto l678
					l679:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
						if buffer[position] != rune('O') {
							goto l675
						}
						position++
					}
				l678:
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
							goto l675
						}
						position++
					}
				l680:
					depth--
					add(rulePegText, position677)
				}
				if !_rules[ruleAction49]() {
					goto l675
				}
				depth--
				add(ruleOr, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 64 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action50)> */
		func() bool {
			position682, tokenIndex682, depth682 := position, tokenIndex, depth
			{
				position683 := position
				depth++
				{
					position684 := position
					depth++
					{
						position685, tokenIndex685, depth685 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l686
						}
						position++
						goto l685
					l686:
						position, tokenIndex, depth = position685, tokenIndex685, depth685
						if buffer[position] != rune('A') {
							goto l682
						}
						position++
					}
				l685:
					{
						position687, tokenIndex687, depth687 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l688
						}
						position++
						goto l687
					l688:
						position, tokenIndex, depth = position687, tokenIndex687, depth687
						if buffer[position] != rune('N') {
							goto l682
						}
						position++
					}
				l687:
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l690
						}
						position++
						goto l689
					l690:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
						if buffer[position] != rune('D') {
							goto l682
						}
						position++
					}
				l689:
					depth--
					add(rulePegText, position684)
				}
				if !_rules[ruleAction50]() {
					goto l682
				}
				depth--
				add(ruleAnd, position683)
			}
			return true
		l682:
			position, tokenIndex, depth = position682, tokenIndex682, depth682
			return false
		},
		/* 65 Equal <- <(<'='> Action51)> */
		func() bool {
			position691, tokenIndex691, depth691 := position, tokenIndex, depth
			{
				position692 := position
				depth++
				{
					position693 := position
					depth++
					if buffer[position] != rune('=') {
						goto l691
					}
					position++
					depth--
					add(rulePegText, position693)
				}
				if !_rules[ruleAction51]() {
					goto l691
				}
				depth--
				add(ruleEqual, position692)
			}
			return true
		l691:
			position, tokenIndex, depth = position691, tokenIndex691, depth691
			return false
		},
		/* 66 Less <- <(<'<'> Action52)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				{
					position696 := position
					depth++
					if buffer[position] != rune('<') {
						goto l694
					}
					position++
					depth--
					add(rulePegText, position696)
				}
				if !_rules[ruleAction52]() {
					goto l694
				}
				depth--
				add(ruleLess, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 67 LessOrEqual <- <(<('<' '=')> Action53)> */
		func() bool {
			position697, tokenIndex697, depth697 := position, tokenIndex, depth
			{
				position698 := position
				depth++
				{
					position699 := position
					depth++
					if buffer[position] != rune('<') {
						goto l697
					}
					position++
					if buffer[position] != rune('=') {
						goto l697
					}
					position++
					depth--
					add(rulePegText, position699)
				}
				if !_rules[ruleAction53]() {
					goto l697
				}
				depth--
				add(ruleLessOrEqual, position698)
			}
			return true
		l697:
			position, tokenIndex, depth = position697, tokenIndex697, depth697
			return false
		},
		/* 68 Greater <- <(<'>'> Action54)> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				{
					position702 := position
					depth++
					if buffer[position] != rune('>') {
						goto l700
					}
					position++
					depth--
					add(rulePegText, position702)
				}
				if !_rules[ruleAction54]() {
					goto l700
				}
				depth--
				add(ruleGreater, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 69 GreaterOrEqual <- <(<('>' '=')> Action55)> */
		func() bool {
			position703, tokenIndex703, depth703 := position, tokenIndex, depth
			{
				position704 := position
				depth++
				{
					position705 := position
					depth++
					if buffer[position] != rune('>') {
						goto l703
					}
					position++
					if buffer[position] != rune('=') {
						goto l703
					}
					position++
					depth--
					add(rulePegText, position705)
				}
				if !_rules[ruleAction55]() {
					goto l703
				}
				depth--
				add(ruleGreaterOrEqual, position704)
			}
			return true
		l703:
			position, tokenIndex, depth = position703, tokenIndex703, depth703
			return false
		},
		/* 70 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action56)> */
		func() bool {
			position706, tokenIndex706, depth706 := position, tokenIndex, depth
			{
				position707 := position
				depth++
				{
					position708 := position
					depth++
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l710
						}
						position++
						if buffer[position] != rune('=') {
							goto l710
						}
						position++
						goto l709
					l710:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if buffer[position] != rune('<') {
							goto l706
						}
						position++
						if buffer[position] != rune('>') {
							goto l706
						}
						position++
					}
				l709:
					depth--
					add(rulePegText, position708)
				}
				if !_rules[ruleAction56]() {
					goto l706
				}
				depth--
				add(ruleNotEqual, position707)
			}
			return true
		l706:
			position, tokenIndex, depth = position706, tokenIndex706, depth706
			return false
		},
		/* 71 Plus <- <(<'+'> Action57)> */
		func() bool {
			position711, tokenIndex711, depth711 := position, tokenIndex, depth
			{
				position712 := position
				depth++
				{
					position713 := position
					depth++
					if buffer[position] != rune('+') {
						goto l711
					}
					position++
					depth--
					add(rulePegText, position713)
				}
				if !_rules[ruleAction57]() {
					goto l711
				}
				depth--
				add(rulePlus, position712)
			}
			return true
		l711:
			position, tokenIndex, depth = position711, tokenIndex711, depth711
			return false
		},
		/* 72 Minus <- <(<'-'> Action58)> */
		func() bool {
			position714, tokenIndex714, depth714 := position, tokenIndex, depth
			{
				position715 := position
				depth++
				{
					position716 := position
					depth++
					if buffer[position] != rune('-') {
						goto l714
					}
					position++
					depth--
					add(rulePegText, position716)
				}
				if !_rules[ruleAction58]() {
					goto l714
				}
				depth--
				add(ruleMinus, position715)
			}
			return true
		l714:
			position, tokenIndex, depth = position714, tokenIndex714, depth714
			return false
		},
		/* 73 Multiply <- <(<'*'> Action59)> */
		func() bool {
			position717, tokenIndex717, depth717 := position, tokenIndex, depth
			{
				position718 := position
				depth++
				{
					position719 := position
					depth++
					if buffer[position] != rune('*') {
						goto l717
					}
					position++
					depth--
					add(rulePegText, position719)
				}
				if !_rules[ruleAction59]() {
					goto l717
				}
				depth--
				add(ruleMultiply, position718)
			}
			return true
		l717:
			position, tokenIndex, depth = position717, tokenIndex717, depth717
			return false
		},
		/* 74 Divide <- <(<'/'> Action60)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722 := position
					depth++
					if buffer[position] != rune('/') {
						goto l720
					}
					position++
					depth--
					add(rulePegText, position722)
				}
				if !_rules[ruleAction60]() {
					goto l720
				}
				depth--
				add(ruleDivide, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 75 Modulo <- <(<'%'> Action61)> */
		func() bool {
			position723, tokenIndex723, depth723 := position, tokenIndex, depth
			{
				position724 := position
				depth++
				{
					position725 := position
					depth++
					if buffer[position] != rune('%') {
						goto l723
					}
					position++
					depth--
					add(rulePegText, position725)
				}
				if !_rules[ruleAction61]() {
					goto l723
				}
				depth--
				add(ruleModulo, position724)
			}
			return true
		l723:
			position, tokenIndex, depth = position723, tokenIndex723, depth723
			return false
		},
		/* 76 Identifier <- <(<ident> Action62)> */
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
				if !_rules[ruleAction62]() {
					goto l726
				}
				depth--
				add(ruleIdentifier, position727)
			}
			return true
		l726:
			position, tokenIndex, depth = position726, tokenIndex726, depth726
			return false
		},
		/* 77 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position729, tokenIndex729, depth729 := position, tokenIndex, depth
			{
				position730 := position
				depth++
				{
					position731, tokenIndex731, depth731 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l732
					}
					position++
					goto l731
				l732:
					position, tokenIndex, depth = position731, tokenIndex731, depth731
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l729
					}
					position++
				}
			l731:
			l733:
				{
					position734, tokenIndex734, depth734 := position, tokenIndex, depth
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l736
						}
						position++
						goto l735
					l736:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l737
						}
						position++
						goto l735
					l737:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l738
						}
						position++
						goto l735
					l738:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if buffer[position] != rune('_') {
							goto l734
						}
						position++
					}
				l735:
					goto l733
				l734:
					position, tokenIndex, depth = position734, tokenIndex734, depth734
				}
				depth--
				add(ruleident, position730)
			}
			return true
		l729:
			position, tokenIndex, depth = position729, tokenIndex729, depth729
			return false
		},
		/* 78 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position740 := position
				depth++
			l741:
				{
					position742, tokenIndex742, depth742 := position, tokenIndex, depth
					{
						position743, tokenIndex743, depth743 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l744
						}
						position++
						goto l743
					l744:
						position, tokenIndex, depth = position743, tokenIndex743, depth743
						if buffer[position] != rune('\t') {
							goto l745
						}
						position++
						goto l743
					l745:
						position, tokenIndex, depth = position743, tokenIndex743, depth743
						if buffer[position] != rune('\n') {
							goto l742
						}
						position++
					}
				l743:
					goto l741
				l742:
					position, tokenIndex, depth = position742, tokenIndex742, depth742
				}
				depth--
				add(rulesp, position740)
			}
			return true
		},
		/* 80 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 81 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 82 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 83 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 84 Action4 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 85 Action5 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 86 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 87 Action7 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 89 Action8 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 90 Action9 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 91 Action10 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 92 Action11 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 93 Action12 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 94 Action13 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 95 Action14 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 96 Action15 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 97 Action16 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 98 Action17 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 99 Action18 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 100 Action19 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 101 Action20 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 102 Action21 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 103 Action22 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 104 Action23 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 105 Action24 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 106 Action25 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 107 Action26 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 108 Action27 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 109 Action28 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 110 Action29 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 111 Action30 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 112 Action31 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 113 Action32 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 114 Action33 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 115 Action34 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 116 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 117 Action36 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 118 Action37 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 119 Action38 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 120 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 121 Action40 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 122 Action41 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 123 Action42 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 124 Action43 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 125 Action44 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 126 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 127 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 128 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 129 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamVal(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 130 Action49 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 131 Action50 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 132 Action51 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 133 Action52 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 134 Action53 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 135 Action54 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 136 Action55 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 137 Action56 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 138 Action57 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 139 Action58 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 140 Action59 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 141 Action60 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 142 Action61 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 143 Action62 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
	}
	p.rules = _rules
}

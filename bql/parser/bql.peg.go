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
	rules  [146]func() bool
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

			p.AssembleInterval()

		case ruleAction14:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction15:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction16:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction17:

			p.EnsureAliasedStreamWindow()

		case ruleAction18:

			p.EnsureAliasedStreamWindow()

		case ruleAction19:

			p.AssembleAliasedStreamWindow()

		case ruleAction20:

			p.AssembleAliasedStreamWindow()

		case ruleAction21:

			p.AssembleStreamWindow()

		case ruleAction22:

			p.AssembleStreamWindow()

		case ruleAction23:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction24:

			p.AssembleSourceSinkParam()

		case ruleAction25:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction26:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction27:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction28:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction29:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction30:

			p.AssembleFuncApp()

		case ruleAction31:

			p.AssembleExpressions(begin, end)

		case ruleAction32:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction33:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction34:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction35:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction36:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction37:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction38:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction39:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction41:

			p.PushComponent(begin, end, Istream)

		case ruleAction42:

			p.PushComponent(begin, end, Dstream)

		case ruleAction43:

			p.PushComponent(begin, end, Rstream)

		case ruleAction44:

			p.PushComponent(begin, end, Tuples)

		case ruleAction45:

			p.PushComponent(begin, end, Seconds)

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamVal(substr))

		case ruleAction50:

			p.PushComponent(begin, end, Or)

		case ruleAction51:

			p.PushComponent(begin, end, And)

		case ruleAction52:

			p.PushComponent(begin, end, Equal)

		case ruleAction53:

			p.PushComponent(begin, end, Less)

		case ruleAction54:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction55:

			p.PushComponent(begin, end, Greater)

		case ruleAction56:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction57:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction58:

			p.PushComponent(begin, end, Plus)

		case ruleAction59:

			p.PushComponent(begin, end, Minus)

		case ruleAction60:

			p.PushComponent(begin, end, Multiply)

		case ruleAction61:

			p.PushComponent(begin, end, Divide)

		case ruleAction62:

			p.PushComponent(begin, end, Modulo)

		case ruleAction63:

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
		/* 15 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{
				position300 := position
				depth++
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l302
					}
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if !_rules[ruleTuplesInterval]() {
						goto l299
					}
				}
			l301:
				depth--
				add(ruleInterval, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 16 TimeInterval <- <(NumericLiteral sp SECONDS Action12)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l303
				}
				if !_rules[rulesp]() {
					goto l303
				}
				if !_rules[ruleSECONDS]() {
					goto l303
				}
				if !_rules[ruleAction12]() {
					goto l303
				}
				depth--
				add(ruleTimeInterval, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 17 TuplesInterval <- <(NumericLiteral sp TUPLES Action13)> */
		func() bool {
			position305, tokenIndex305, depth305 := position, tokenIndex, depth
			{
				position306 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l305
				}
				if !_rules[rulesp]() {
					goto l305
				}
				if !_rules[ruleTUPLES]() {
					goto l305
				}
				if !_rules[ruleAction13]() {
					goto l305
				}
				depth--
				add(ruleTuplesInterval, position306)
			}
			return true
		l305:
			position, tokenIndex, depth = position305, tokenIndex305, depth305
			return false
		},
		/* 18 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position307, tokenIndex307, depth307 := position, tokenIndex, depth
			{
				position308 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l307
				}
				if !_rules[rulesp]() {
					goto l307
				}
			l309:
				{
					position310, tokenIndex310, depth310 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l310
					}
					position++
					if !_rules[rulesp]() {
						goto l310
					}
					if !_rules[ruleRelationLike]() {
						goto l310
					}
					goto l309
				l310:
					position, tokenIndex, depth = position310, tokenIndex310, depth310
				}
				depth--
				add(ruleRelations, position308)
			}
			return true
		l307:
			position, tokenIndex, depth = position307, tokenIndex307, depth307
			return false
		},
		/* 19 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position311, tokenIndex311, depth311 := position, tokenIndex, depth
			{
				position312 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l311
				}
				if !_rules[rulesp]() {
					goto l311
				}
			l313:
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l314
					}
					position++
					if !_rules[rulesp]() {
						goto l314
					}
					if !_rules[ruleDefRelationLike]() {
						goto l314
					}
					goto l313
				l314:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
				}
				depth--
				add(ruleDefRelations, position312)
			}
			return true
		l311:
			position, tokenIndex, depth = position311, tokenIndex311, depth311
			return false
		},
		/* 20 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action14)> */
		func() bool {
			position315, tokenIndex315, depth315 := position, tokenIndex, depth
			{
				position316 := position
				depth++
				{
					position317 := position
					depth++
					{
						position318, tokenIndex318, depth318 := position, tokenIndex, depth
						{
							position320, tokenIndex320, depth320 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l321
							}
							position++
							goto l320
						l321:
							position, tokenIndex, depth = position320, tokenIndex320, depth320
							if buffer[position] != rune('W') {
								goto l318
							}
							position++
						}
					l320:
						{
							position322, tokenIndex322, depth322 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l323
							}
							position++
							goto l322
						l323:
							position, tokenIndex, depth = position322, tokenIndex322, depth322
							if buffer[position] != rune('H') {
								goto l318
							}
							position++
						}
					l322:
						{
							position324, tokenIndex324, depth324 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l325
							}
							position++
							goto l324
						l325:
							position, tokenIndex, depth = position324, tokenIndex324, depth324
							if buffer[position] != rune('E') {
								goto l318
							}
							position++
						}
					l324:
						{
							position326, tokenIndex326, depth326 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l327
							}
							position++
							goto l326
						l327:
							position, tokenIndex, depth = position326, tokenIndex326, depth326
							if buffer[position] != rune('R') {
								goto l318
							}
							position++
						}
					l326:
						{
							position328, tokenIndex328, depth328 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l329
							}
							position++
							goto l328
						l329:
							position, tokenIndex, depth = position328, tokenIndex328, depth328
							if buffer[position] != rune('E') {
								goto l318
							}
							position++
						}
					l328:
						if !_rules[rulesp]() {
							goto l318
						}
						if !_rules[ruleExpression]() {
							goto l318
						}
						goto l319
					l318:
						position, tokenIndex, depth = position318, tokenIndex318, depth318
					}
				l319:
					depth--
					add(rulePegText, position317)
				}
				if !_rules[ruleAction14]() {
					goto l315
				}
				depth--
				add(ruleFilter, position316)
			}
			return true
		l315:
			position, tokenIndex, depth = position315, tokenIndex315, depth315
			return false
		},
		/* 21 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action15)> */
		func() bool {
			position330, tokenIndex330, depth330 := position, tokenIndex, depth
			{
				position331 := position
				depth++
				{
					position332 := position
					depth++
					{
						position333, tokenIndex333, depth333 := position, tokenIndex, depth
						{
							position335, tokenIndex335, depth335 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l336
							}
							position++
							goto l335
						l336:
							position, tokenIndex, depth = position335, tokenIndex335, depth335
							if buffer[position] != rune('G') {
								goto l333
							}
							position++
						}
					l335:
						{
							position337, tokenIndex337, depth337 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l338
							}
							position++
							goto l337
						l338:
							position, tokenIndex, depth = position337, tokenIndex337, depth337
							if buffer[position] != rune('R') {
								goto l333
							}
							position++
						}
					l337:
						{
							position339, tokenIndex339, depth339 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l340
							}
							position++
							goto l339
						l340:
							position, tokenIndex, depth = position339, tokenIndex339, depth339
							if buffer[position] != rune('O') {
								goto l333
							}
							position++
						}
					l339:
						{
							position341, tokenIndex341, depth341 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l342
							}
							position++
							goto l341
						l342:
							position, tokenIndex, depth = position341, tokenIndex341, depth341
							if buffer[position] != rune('U') {
								goto l333
							}
							position++
						}
					l341:
						{
							position343, tokenIndex343, depth343 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l344
							}
							position++
							goto l343
						l344:
							position, tokenIndex, depth = position343, tokenIndex343, depth343
							if buffer[position] != rune('P') {
								goto l333
							}
							position++
						}
					l343:
						if !_rules[rulesp]() {
							goto l333
						}
						{
							position345, tokenIndex345, depth345 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l346
							}
							position++
							goto l345
						l346:
							position, tokenIndex, depth = position345, tokenIndex345, depth345
							if buffer[position] != rune('B') {
								goto l333
							}
							position++
						}
					l345:
						{
							position347, tokenIndex347, depth347 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l348
							}
							position++
							goto l347
						l348:
							position, tokenIndex, depth = position347, tokenIndex347, depth347
							if buffer[position] != rune('Y') {
								goto l333
							}
							position++
						}
					l347:
						if !_rules[rulesp]() {
							goto l333
						}
						if !_rules[ruleGroupList]() {
							goto l333
						}
						goto l334
					l333:
						position, tokenIndex, depth = position333, tokenIndex333, depth333
					}
				l334:
					depth--
					add(rulePegText, position332)
				}
				if !_rules[ruleAction15]() {
					goto l330
				}
				depth--
				add(ruleGrouping, position331)
			}
			return true
		l330:
			position, tokenIndex, depth = position330, tokenIndex330, depth330
			return false
		},
		/* 22 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l349
				}
				if !_rules[rulesp]() {
					goto l349
				}
			l351:
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l352
					}
					position++
					if !_rules[rulesp]() {
						goto l352
					}
					if !_rules[ruleExpression]() {
						goto l352
					}
					goto l351
				l352:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
				}
				depth--
				add(ruleGroupList, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 23 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action16)> */
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
							if buffer[position] != rune('h') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex, depth = position358, tokenIndex358, depth358
							if buffer[position] != rune('H') {
								goto l356
							}
							position++
						}
					l358:
						{
							position360, tokenIndex360, depth360 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l361
							}
							position++
							goto l360
						l361:
							position, tokenIndex, depth = position360, tokenIndex360, depth360
							if buffer[position] != rune('A') {
								goto l356
							}
							position++
						}
					l360:
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('V') {
								goto l356
							}
							position++
						}
					l362:
						{
							position364, tokenIndex364, depth364 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l365
							}
							position++
							goto l364
						l365:
							position, tokenIndex, depth = position364, tokenIndex364, depth364
							if buffer[position] != rune('I') {
								goto l356
							}
							position++
						}
					l364:
						{
							position366, tokenIndex366, depth366 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l367
							}
							position++
							goto l366
						l367:
							position, tokenIndex, depth = position366, tokenIndex366, depth366
							if buffer[position] != rune('N') {
								goto l356
							}
							position++
						}
					l366:
						{
							position368, tokenIndex368, depth368 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l369
							}
							position++
							goto l368
						l369:
							position, tokenIndex, depth = position368, tokenIndex368, depth368
							if buffer[position] != rune('G') {
								goto l356
							}
							position++
						}
					l368:
						if !_rules[rulesp]() {
							goto l356
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction16]() {
					goto l353
				}
				depth--
				add(ruleHaving, position354)
			}
			return true
		l353:
			position, tokenIndex, depth = position353, tokenIndex353, depth353
			return false
		},
		/* 24 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action17))> */
		func() bool {
			position370, tokenIndex370, depth370 := position, tokenIndex, depth
			{
				position371 := position
				depth++
				{
					position372, tokenIndex372, depth372 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l373
					}
					goto l372
				l373:
					position, tokenIndex, depth = position372, tokenIndex372, depth372
					if !_rules[ruleStreamWindow]() {
						goto l370
					}
					if !_rules[ruleAction17]() {
						goto l370
					}
				}
			l372:
				depth--
				add(ruleRelationLike, position371)
			}
			return true
		l370:
			position, tokenIndex, depth = position370, tokenIndex370, depth370
			return false
		},
		/* 25 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action18))> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l377
					}
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					if !_rules[ruleDefStreamWindow]() {
						goto l374
					}
					if !_rules[ruleAction18]() {
						goto l374
					}
				}
			l376:
				depth--
				add(ruleDefRelationLike, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
			return false
		},
		/* 26 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action19)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				if !_rules[ruleStreamWindow]() {
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
				add(ruleAliasedStreamWindow, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 27 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action20)> */
		func() bool {
			position384, tokenIndex384, depth384 := position, tokenIndex, depth
			{
				position385 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l384
				}
				if !_rules[rulesp]() {
					goto l384
				}
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('A') {
						goto l384
					}
					position++
				}
			l386:
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l389
					}
					position++
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					if buffer[position] != rune('S') {
						goto l384
					}
					position++
				}
			l388:
				if !_rules[rulesp]() {
					goto l384
				}
				if !_rules[ruleIdentifier]() {
					goto l384
				}
				if !_rules[ruleAction20]() {
					goto l384
				}
				depth--
				add(ruleDefAliasedStreamWindow, position385)
			}
			return true
		l384:
			position, tokenIndex, depth = position384, tokenIndex384, depth384
			return false
		},
		/* 28 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action21)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				if !_rules[ruleStream]() {
					goto l390
				}
				if !_rules[rulesp]() {
					goto l390
				}
				if buffer[position] != rune('[') {
					goto l390
				}
				position++
				if !_rules[rulesp]() {
					goto l390
				}
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if buffer[position] != rune('R') {
						goto l390
					}
					position++
				}
			l392:
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l395
					}
					position++
					goto l394
				l395:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
					if buffer[position] != rune('A') {
						goto l390
					}
					position++
				}
			l394:
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('N') {
						goto l390
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if buffer[position] != rune('G') {
						goto l390
					}
					position++
				}
			l398:
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l401
					}
					position++
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					if buffer[position] != rune('E') {
						goto l390
					}
					position++
				}
			l400:
				if !_rules[rulesp]() {
					goto l390
				}
				if !_rules[ruleInterval]() {
					goto l390
				}
				if !_rules[rulesp]() {
					goto l390
				}
				if buffer[position] != rune(']') {
					goto l390
				}
				position++
				if !_rules[ruleAction21]() {
					goto l390
				}
				depth--
				add(ruleStreamWindow, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 29 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action22)> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				if !_rules[ruleStream]() {
					goto l402
				}
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l404
					}
					if buffer[position] != rune('[') {
						goto l404
					}
					position++
					if !_rules[rulesp]() {
						goto l404
					}
					{
						position406, tokenIndex406, depth406 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l407
						}
						position++
						goto l406
					l407:
						position, tokenIndex, depth = position406, tokenIndex406, depth406
						if buffer[position] != rune('R') {
							goto l404
						}
						position++
					}
				l406:
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l409
						}
						position++
						goto l408
					l409:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
						if buffer[position] != rune('A') {
							goto l404
						}
						position++
					}
				l408:
					{
						position410, tokenIndex410, depth410 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l411
						}
						position++
						goto l410
					l411:
						position, tokenIndex, depth = position410, tokenIndex410, depth410
						if buffer[position] != rune('N') {
							goto l404
						}
						position++
					}
				l410:
					{
						position412, tokenIndex412, depth412 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l413
						}
						position++
						goto l412
					l413:
						position, tokenIndex, depth = position412, tokenIndex412, depth412
						if buffer[position] != rune('G') {
							goto l404
						}
						position++
					}
				l412:
					{
						position414, tokenIndex414, depth414 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l415
						}
						position++
						goto l414
					l415:
						position, tokenIndex, depth = position414, tokenIndex414, depth414
						if buffer[position] != rune('E') {
							goto l404
						}
						position++
					}
				l414:
					if !_rules[rulesp]() {
						goto l404
					}
					if !_rules[ruleInterval]() {
						goto l404
					}
					if !_rules[rulesp]() {
						goto l404
					}
					if buffer[position] != rune(']') {
						goto l404
					}
					position++
					goto l405
				l404:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
				}
			l405:
				if !_rules[ruleAction22]() {
					goto l402
				}
				depth--
				add(ruleDefStreamWindow, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 30 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action23)> */
		func() bool {
			position416, tokenIndex416, depth416 := position, tokenIndex, depth
			{
				position417 := position
				depth++
				{
					position418 := position
					depth++
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						{
							position421, tokenIndex421, depth421 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l422
							}
							position++
							goto l421
						l422:
							position, tokenIndex, depth = position421, tokenIndex421, depth421
							if buffer[position] != rune('W') {
								goto l419
							}
							position++
						}
					l421:
						{
							position423, tokenIndex423, depth423 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l424
							}
							position++
							goto l423
						l424:
							position, tokenIndex, depth = position423, tokenIndex423, depth423
							if buffer[position] != rune('I') {
								goto l419
							}
							position++
						}
					l423:
						{
							position425, tokenIndex425, depth425 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l426
							}
							position++
							goto l425
						l426:
							position, tokenIndex, depth = position425, tokenIndex425, depth425
							if buffer[position] != rune('T') {
								goto l419
							}
							position++
						}
					l425:
						{
							position427, tokenIndex427, depth427 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l428
							}
							position++
							goto l427
						l428:
							position, tokenIndex, depth = position427, tokenIndex427, depth427
							if buffer[position] != rune('H') {
								goto l419
							}
							position++
						}
					l427:
						if !_rules[rulesp]() {
							goto l419
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l419
						}
						if !_rules[rulesp]() {
							goto l419
						}
					l429:
						{
							position430, tokenIndex430, depth430 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l430
							}
							position++
							if !_rules[rulesp]() {
								goto l430
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l430
							}
							goto l429
						l430:
							position, tokenIndex, depth = position430, tokenIndex430, depth430
						}
						goto l420
					l419:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
					}
				l420:
					depth--
					add(rulePegText, position418)
				}
				if !_rules[ruleAction23]() {
					goto l416
				}
				depth--
				add(ruleSourceSinkSpecs, position417)
			}
			return true
		l416:
			position, tokenIndex, depth = position416, tokenIndex416, depth416
			return false
		},
		/* 31 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action24)> */
		func() bool {
			position431, tokenIndex431, depth431 := position, tokenIndex, depth
			{
				position432 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l431
				}
				if buffer[position] != rune('=') {
					goto l431
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l431
				}
				if !_rules[ruleAction24]() {
					goto l431
				}
				depth--
				add(ruleSourceSinkParam, position432)
			}
			return true
		l431:
			position, tokenIndex, depth = position431, tokenIndex431, depth431
			return false
		},
		/* 32 Expression <- <orExpr> */
		func() bool {
			position433, tokenIndex433, depth433 := position, tokenIndex, depth
			{
				position434 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l433
				}
				depth--
				add(ruleExpression, position434)
			}
			return true
		l433:
			position, tokenIndex, depth = position433, tokenIndex433, depth433
			return false
		},
		/* 33 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action25)> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				{
					position437 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l435
					}
					if !_rules[rulesp]() {
						goto l435
					}
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l438
						}
						if !_rules[rulesp]() {
							goto l438
						}
						if !_rules[ruleandExpr]() {
							goto l438
						}
						goto l439
					l438:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
					}
				l439:
					depth--
					add(rulePegText, position437)
				}
				if !_rules[ruleAction25]() {
					goto l435
				}
				depth--
				add(ruleorExpr, position436)
			}
			return true
		l435:
			position, tokenIndex, depth = position435, tokenIndex435, depth435
			return false
		},
		/* 34 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action26)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				{
					position442 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l440
					}
					if !_rules[rulesp]() {
						goto l440
					}
					{
						position443, tokenIndex443, depth443 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l443
						}
						if !_rules[rulesp]() {
							goto l443
						}
						if !_rules[rulecomparisonExpr]() {
							goto l443
						}
						goto l444
					l443:
						position, tokenIndex, depth = position443, tokenIndex443, depth443
					}
				l444:
					depth--
					add(rulePegText, position442)
				}
				if !_rules[ruleAction26]() {
					goto l440
				}
				depth--
				add(ruleandExpr, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 35 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action27)> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				{
					position447 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l445
					}
					if !_rules[rulesp]() {
						goto l445
					}
					{
						position448, tokenIndex448, depth448 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l448
						}
						if !_rules[rulesp]() {
							goto l448
						}
						if !_rules[ruletermExpr]() {
							goto l448
						}
						goto l449
					l448:
						position, tokenIndex, depth = position448, tokenIndex448, depth448
					}
				l449:
					depth--
					add(rulePegText, position447)
				}
				if !_rules[ruleAction27]() {
					goto l445
				}
				depth--
				add(rulecomparisonExpr, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 36 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action28)> */
		func() bool {
			position450, tokenIndex450, depth450 := position, tokenIndex, depth
			{
				position451 := position
				depth++
				{
					position452 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l450
					}
					if !_rules[rulesp]() {
						goto l450
					}
					{
						position453, tokenIndex453, depth453 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l453
						}
						if !_rules[rulesp]() {
							goto l453
						}
						if !_rules[ruleproductExpr]() {
							goto l453
						}
						goto l454
					l453:
						position, tokenIndex, depth = position453, tokenIndex453, depth453
					}
				l454:
					depth--
					add(rulePegText, position452)
				}
				if !_rules[ruleAction28]() {
					goto l450
				}
				depth--
				add(ruletermExpr, position451)
			}
			return true
		l450:
			position, tokenIndex, depth = position450, tokenIndex450, depth450
			return false
		},
		/* 37 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action29)> */
		func() bool {
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				{
					position457 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l455
					}
					if !_rules[rulesp]() {
						goto l455
					}
					{
						position458, tokenIndex458, depth458 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l458
						}
						if !_rules[rulesp]() {
							goto l458
						}
						if !_rules[rulebaseExpr]() {
							goto l458
						}
						goto l459
					l458:
						position, tokenIndex, depth = position458, tokenIndex458, depth458
					}
				l459:
					depth--
					add(rulePegText, position457)
				}
				if !_rules[ruleAction29]() {
					goto l455
				}
				depth--
				add(ruleproductExpr, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
			return false
		},
		/* 38 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position460, tokenIndex460, depth460 := position, tokenIndex, depth
			{
				position461 := position
				depth++
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l463
					}
					position++
					if !_rules[rulesp]() {
						goto l463
					}
					if !_rules[ruleExpression]() {
						goto l463
					}
					if !_rules[rulesp]() {
						goto l463
					}
					if buffer[position] != rune(')') {
						goto l463
					}
					position++
					goto l462
				l463:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if !_rules[ruleBooleanLiteral]() {
						goto l464
					}
					goto l462
				l464:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if !_rules[ruleFuncApp]() {
						goto l465
					}
					goto l462
				l465:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if !_rules[ruleRowValue]() {
						goto l466
					}
					goto l462
				l466:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if !_rules[ruleLiteral]() {
						goto l460
					}
				}
			l462:
				depth--
				add(rulebaseExpr, position461)
			}
			return true
		l460:
			position, tokenIndex, depth = position460, tokenIndex460, depth460
			return false
		},
		/* 39 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action30)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l467
				}
				if !_rules[rulesp]() {
					goto l467
				}
				if buffer[position] != rune('(') {
					goto l467
				}
				position++
				if !_rules[rulesp]() {
					goto l467
				}
				if !_rules[ruleFuncParams]() {
					goto l467
				}
				if !_rules[rulesp]() {
					goto l467
				}
				if buffer[position] != rune(')') {
					goto l467
				}
				position++
				if !_rules[ruleAction30]() {
					goto l467
				}
				depth--
				add(ruleFuncApp, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 40 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action31)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				{
					position471 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l469
					}
					if !_rules[rulesp]() {
						goto l469
					}
				l472:
					{
						position473, tokenIndex473, depth473 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l473
						}
						position++
						if !_rules[rulesp]() {
							goto l473
						}
						if !_rules[ruleExpression]() {
							goto l473
						}
						goto l472
					l473:
						position, tokenIndex, depth = position473, tokenIndex473, depth473
					}
					depth--
					add(rulePegText, position471)
				}
				if !_rules[ruleAction31]() {
					goto l469
				}
				depth--
				add(ruleFuncParams, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 41 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				{
					position476, tokenIndex476, depth476 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l477
					}
					goto l476
				l477:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if !_rules[ruleNumericLiteral]() {
						goto l478
					}
					goto l476
				l478:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if !_rules[ruleStringLiteral]() {
						goto l474
					}
				}
			l476:
				depth--
				add(ruleLiteral, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 42 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l482
					}
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleNotEqual]() {
						goto l483
					}
					goto l481
				l483:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleLessOrEqual]() {
						goto l484
					}
					goto l481
				l484:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleLess]() {
						goto l485
					}
					goto l481
				l485:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleGreaterOrEqual]() {
						goto l486
					}
					goto l481
				l486:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleGreater]() {
						goto l487
					}
					goto l481
				l487:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleNotEqual]() {
						goto l479
					}
				}
			l481:
				depth--
				add(ruleComparisonOp, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 43 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position488, tokenIndex488, depth488 := position, tokenIndex, depth
			{
				position489 := position
				depth++
				{
					position490, tokenIndex490, depth490 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l491
					}
					goto l490
				l491:
					position, tokenIndex, depth = position490, tokenIndex490, depth490
					if !_rules[ruleMinus]() {
						goto l488
					}
				}
			l490:
				depth--
				add(rulePlusMinusOp, position489)
			}
			return true
		l488:
			position, tokenIndex, depth = position488, tokenIndex488, depth488
			return false
		},
		/* 44 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position492, tokenIndex492, depth492 := position, tokenIndex, depth
			{
				position493 := position
				depth++
				{
					position494, tokenIndex494, depth494 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l495
					}
					goto l494
				l495:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if !_rules[ruleDivide]() {
						goto l496
					}
					goto l494
				l496:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if !_rules[ruleModulo]() {
						goto l492
					}
				}
			l494:
				depth--
				add(ruleMultDivOp, position493)
			}
			return true
		l492:
			position, tokenIndex, depth = position492, tokenIndex492, depth492
			return false
		},
		/* 45 Stream <- <(<ident> Action32)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				{
					position499 := position
					depth++
					if !_rules[ruleident]() {
						goto l497
					}
					depth--
					add(rulePegText, position499)
				}
				if !_rules[ruleAction32]() {
					goto l497
				}
				depth--
				add(ruleStream, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 46 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action33)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				{
					position502 := position
					depth++
					{
						position503, tokenIndex503, depth503 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l503
						}
						if buffer[position] != rune(':') {
							goto l503
						}
						position++
						goto l504
					l503:
						position, tokenIndex, depth = position503, tokenIndex503, depth503
					}
				l504:
					{
						position505, tokenIndex505, depth505 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l506
						}
						position++
						goto l505
					l506:
						position, tokenIndex, depth = position505, tokenIndex505, depth505
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l500
						}
						position++
					}
				l505:
				l507:
					{
						position508, tokenIndex508, depth508 := position, tokenIndex, depth
						{
							position509, tokenIndex509, depth509 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l510
							}
							position++
							goto l509
						l510:
							position, tokenIndex, depth = position509, tokenIndex509, depth509
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l511
							}
							position++
							goto l509
						l511:
							position, tokenIndex, depth = position509, tokenIndex509, depth509
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l512
							}
							position++
							goto l509
						l512:
							position, tokenIndex, depth = position509, tokenIndex509, depth509
							if buffer[position] != rune('_') {
								goto l513
							}
							position++
							goto l509
						l513:
							position, tokenIndex, depth = position509, tokenIndex509, depth509
							if buffer[position] != rune('.') {
								goto l508
							}
							position++
						}
					l509:
						goto l507
					l508:
						position, tokenIndex, depth = position508, tokenIndex508, depth508
					}
					depth--
					add(rulePegText, position502)
				}
				if !_rules[ruleAction33]() {
					goto l500
				}
				depth--
				add(ruleRowValue, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 47 NumericLiteral <- <(<('-'? [0-9]+)> Action34)> */
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
						if buffer[position] != rune('-') {
							goto l517
						}
						position++
						goto l518
					l517:
						position, tokenIndex, depth = position517, tokenIndex517, depth517
					}
				l518:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l514
					}
					position++
				l519:
					{
						position520, tokenIndex520, depth520 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l520
						}
						position++
						goto l519
					l520:
						position, tokenIndex, depth = position520, tokenIndex520, depth520
					}
					depth--
					add(rulePegText, position516)
				}
				if !_rules[ruleAction34]() {
					goto l514
				}
				depth--
				add(ruleNumericLiteral, position515)
			}
			return true
		l514:
			position, tokenIndex, depth = position514, tokenIndex514, depth514
			return false
		},
		/* 48 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action35)> */
		func() bool {
			position521, tokenIndex521, depth521 := position, tokenIndex, depth
			{
				position522 := position
				depth++
				{
					position523 := position
					depth++
					{
						position524, tokenIndex524, depth524 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l524
						}
						position++
						goto l525
					l524:
						position, tokenIndex, depth = position524, tokenIndex524, depth524
					}
				l525:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l521
					}
					position++
				l526:
					{
						position527, tokenIndex527, depth527 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l527
						}
						position++
						goto l526
					l527:
						position, tokenIndex, depth = position527, tokenIndex527, depth527
					}
					if buffer[position] != rune('.') {
						goto l521
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l521
					}
					position++
				l528:
					{
						position529, tokenIndex529, depth529 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l529
						}
						position++
						goto l528
					l529:
						position, tokenIndex, depth = position529, tokenIndex529, depth529
					}
					depth--
					add(rulePegText, position523)
				}
				if !_rules[ruleAction35]() {
					goto l521
				}
				depth--
				add(ruleFloatLiteral, position522)
			}
			return true
		l521:
			position, tokenIndex, depth = position521, tokenIndex521, depth521
			return false
		},
		/* 49 Function <- <(<ident> Action36)> */
		func() bool {
			position530, tokenIndex530, depth530 := position, tokenIndex, depth
			{
				position531 := position
				depth++
				{
					position532 := position
					depth++
					if !_rules[ruleident]() {
						goto l530
					}
					depth--
					add(rulePegText, position532)
				}
				if !_rules[ruleAction36]() {
					goto l530
				}
				depth--
				add(ruleFunction, position531)
			}
			return true
		l530:
			position, tokenIndex, depth = position530, tokenIndex530, depth530
			return false
		},
		/* 50 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position533, tokenIndex533, depth533 := position, tokenIndex, depth
			{
				position534 := position
				depth++
				{
					position535, tokenIndex535, depth535 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l536
					}
					goto l535
				l536:
					position, tokenIndex, depth = position535, tokenIndex535, depth535
					if !_rules[ruleFALSE]() {
						goto l533
					}
				}
			l535:
				depth--
				add(ruleBooleanLiteral, position534)
			}
			return true
		l533:
			position, tokenIndex, depth = position533, tokenIndex533, depth533
			return false
		},
		/* 51 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action37)> */
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
						if buffer[position] != rune('t') {
							goto l541
						}
						position++
						goto l540
					l541:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
						if buffer[position] != rune('T') {
							goto l537
						}
						position++
					}
				l540:
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l543
						}
						position++
						goto l542
					l543:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
						if buffer[position] != rune('R') {
							goto l537
						}
						position++
					}
				l542:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l545
						}
						position++
						goto l544
					l545:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
						if buffer[position] != rune('U') {
							goto l537
						}
						position++
					}
				l544:
					{
						position546, tokenIndex546, depth546 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l547
						}
						position++
						goto l546
					l547:
						position, tokenIndex, depth = position546, tokenIndex546, depth546
						if buffer[position] != rune('E') {
							goto l537
						}
						position++
					}
				l546:
					depth--
					add(rulePegText, position539)
				}
				if !_rules[ruleAction37]() {
					goto l537
				}
				depth--
				add(ruleTRUE, position538)
			}
			return true
		l537:
			position, tokenIndex, depth = position537, tokenIndex537, depth537
			return false
		},
		/* 52 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action38)> */
		func() bool {
			position548, tokenIndex548, depth548 := position, tokenIndex, depth
			{
				position549 := position
				depth++
				{
					position550 := position
					depth++
					{
						position551, tokenIndex551, depth551 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l552
						}
						position++
						goto l551
					l552:
						position, tokenIndex, depth = position551, tokenIndex551, depth551
						if buffer[position] != rune('F') {
							goto l548
						}
						position++
					}
				l551:
					{
						position553, tokenIndex553, depth553 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l554
						}
						position++
						goto l553
					l554:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						if buffer[position] != rune('A') {
							goto l548
						}
						position++
					}
				l553:
					{
						position555, tokenIndex555, depth555 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l556
						}
						position++
						goto l555
					l556:
						position, tokenIndex, depth = position555, tokenIndex555, depth555
						if buffer[position] != rune('L') {
							goto l548
						}
						position++
					}
				l555:
					{
						position557, tokenIndex557, depth557 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l558
						}
						position++
						goto l557
					l558:
						position, tokenIndex, depth = position557, tokenIndex557, depth557
						if buffer[position] != rune('S') {
							goto l548
						}
						position++
					}
				l557:
					{
						position559, tokenIndex559, depth559 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l560
						}
						position++
						goto l559
					l560:
						position, tokenIndex, depth = position559, tokenIndex559, depth559
						if buffer[position] != rune('E') {
							goto l548
						}
						position++
					}
				l559:
					depth--
					add(rulePegText, position550)
				}
				if !_rules[ruleAction38]() {
					goto l548
				}
				depth--
				add(ruleFALSE, position549)
			}
			return true
		l548:
			position, tokenIndex, depth = position548, tokenIndex548, depth548
			return false
		},
		/* 53 Wildcard <- <(<'*'> Action39)> */
		func() bool {
			position561, tokenIndex561, depth561 := position, tokenIndex, depth
			{
				position562 := position
				depth++
				{
					position563 := position
					depth++
					if buffer[position] != rune('*') {
						goto l561
					}
					position++
					depth--
					add(rulePegText, position563)
				}
				if !_rules[ruleAction39]() {
					goto l561
				}
				depth--
				add(ruleWildcard, position562)
			}
			return true
		l561:
			position, tokenIndex, depth = position561, tokenIndex561, depth561
			return false
		},
		/* 54 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action40)> */
		func() bool {
			position564, tokenIndex564, depth564 := position, tokenIndex, depth
			{
				position565 := position
				depth++
				{
					position566 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l564
					}
					position++
				l567:
					{
						position568, tokenIndex568, depth568 := position, tokenIndex, depth
						{
							position569, tokenIndex569, depth569 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l570
							}
							position++
							if buffer[position] != rune('\'') {
								goto l570
							}
							position++
							goto l569
						l570:
							position, tokenIndex, depth = position569, tokenIndex569, depth569
							{
								position571, tokenIndex571, depth571 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l571
								}
								position++
								goto l568
							l571:
								position, tokenIndex, depth = position571, tokenIndex571, depth571
							}
							if !matchDot() {
								goto l568
							}
						}
					l569:
						goto l567
					l568:
						position, tokenIndex, depth = position568, tokenIndex568, depth568
					}
					if buffer[position] != rune('\'') {
						goto l564
					}
					position++
					depth--
					add(rulePegText, position566)
				}
				if !_rules[ruleAction40]() {
					goto l564
				}
				depth--
				add(ruleStringLiteral, position565)
			}
			return true
		l564:
			position, tokenIndex, depth = position564, tokenIndex564, depth564
			return false
		},
		/* 55 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action41)> */
		func() bool {
			position572, tokenIndex572, depth572 := position, tokenIndex, depth
			{
				position573 := position
				depth++
				{
					position574 := position
					depth++
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l576
						}
						position++
						goto l575
					l576:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
						if buffer[position] != rune('I') {
							goto l572
						}
						position++
					}
				l575:
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('S') {
							goto l572
						}
						position++
					}
				l577:
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l580
						}
						position++
						goto l579
					l580:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
						if buffer[position] != rune('T') {
							goto l572
						}
						position++
					}
				l579:
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l582
						}
						position++
						goto l581
					l582:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
						if buffer[position] != rune('R') {
							goto l572
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
							goto l572
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
							goto l572
						}
						position++
					}
				l585:
					{
						position587, tokenIndex587, depth587 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l588
						}
						position++
						goto l587
					l588:
						position, tokenIndex, depth = position587, tokenIndex587, depth587
						if buffer[position] != rune('M') {
							goto l572
						}
						position++
					}
				l587:
					depth--
					add(rulePegText, position574)
				}
				if !_rules[ruleAction41]() {
					goto l572
				}
				depth--
				add(ruleISTREAM, position573)
			}
			return true
		l572:
			position, tokenIndex, depth = position572, tokenIndex572, depth572
			return false
		},
		/* 56 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action42)> */
		func() bool {
			position589, tokenIndex589, depth589 := position, tokenIndex, depth
			{
				position590 := position
				depth++
				{
					position591 := position
					depth++
					{
						position592, tokenIndex592, depth592 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex, depth = position592, tokenIndex592, depth592
						if buffer[position] != rune('D') {
							goto l589
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
							goto l589
						}
						position++
					}
				l594:
					{
						position596, tokenIndex596, depth596 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l597
						}
						position++
						goto l596
					l597:
						position, tokenIndex, depth = position596, tokenIndex596, depth596
						if buffer[position] != rune('T') {
							goto l589
						}
						position++
					}
				l596:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l599
						}
						position++
						goto l598
					l599:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
						if buffer[position] != rune('R') {
							goto l589
						}
						position++
					}
				l598:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l601
						}
						position++
						goto l600
					l601:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
						if buffer[position] != rune('E') {
							goto l589
						}
						position++
					}
				l600:
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l603
						}
						position++
						goto l602
					l603:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if buffer[position] != rune('A') {
							goto l589
						}
						position++
					}
				l602:
					{
						position604, tokenIndex604, depth604 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l605
						}
						position++
						goto l604
					l605:
						position, tokenIndex, depth = position604, tokenIndex604, depth604
						if buffer[position] != rune('M') {
							goto l589
						}
						position++
					}
				l604:
					depth--
					add(rulePegText, position591)
				}
				if !_rules[ruleAction42]() {
					goto l589
				}
				depth--
				add(ruleDSTREAM, position590)
			}
			return true
		l589:
			position, tokenIndex, depth = position589, tokenIndex589, depth589
			return false
		},
		/* 57 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action43)> */
		func() bool {
			position606, tokenIndex606, depth606 := position, tokenIndex, depth
			{
				position607 := position
				depth++
				{
					position608 := position
					depth++
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
							goto l606
						}
						position++
					}
				l609:
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('S') {
							goto l606
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('T') {
							goto l606
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('R') {
							goto l606
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
							goto l606
						}
						position++
					}
				l617:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('A') {
							goto l606
						}
						position++
					}
				l619:
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l622
						}
						position++
						goto l621
					l622:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
						if buffer[position] != rune('M') {
							goto l606
						}
						position++
					}
				l621:
					depth--
					add(rulePegText, position608)
				}
				if !_rules[ruleAction43]() {
					goto l606
				}
				depth--
				add(ruleRSTREAM, position607)
			}
			return true
		l606:
			position, tokenIndex, depth = position606, tokenIndex606, depth606
			return false
		},
		/* 58 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action44)> */
		func() bool {
			position623, tokenIndex623, depth623 := position, tokenIndex, depth
			{
				position624 := position
				depth++
				{
					position625 := position
					depth++
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('T') {
							goto l623
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('U') {
							goto l623
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('P') {
							goto l623
						}
						position++
					}
				l630:
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l633
						}
						position++
						goto l632
					l633:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
						if buffer[position] != rune('L') {
							goto l623
						}
						position++
					}
				l632:
					{
						position634, tokenIndex634, depth634 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l635
						}
						position++
						goto l634
					l635:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
						if buffer[position] != rune('E') {
							goto l623
						}
						position++
					}
				l634:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l637
						}
						position++
						goto l636
					l637:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
						if buffer[position] != rune('S') {
							goto l623
						}
						position++
					}
				l636:
					depth--
					add(rulePegText, position625)
				}
				if !_rules[ruleAction44]() {
					goto l623
				}
				depth--
				add(ruleTUPLES, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 59 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action45)> */
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
						if buffer[position] != rune('s') {
							goto l642
						}
						position++
						goto l641
					l642:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
						if buffer[position] != rune('S') {
							goto l638
						}
						position++
					}
				l641:
					{
						position643, tokenIndex643, depth643 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l644
						}
						position++
						goto l643
					l644:
						position, tokenIndex, depth = position643, tokenIndex643, depth643
						if buffer[position] != rune('E') {
							goto l638
						}
						position++
					}
				l643:
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l646
						}
						position++
						goto l645
					l646:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
						if buffer[position] != rune('C') {
							goto l638
						}
						position++
					}
				l645:
					{
						position647, tokenIndex647, depth647 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l648
						}
						position++
						goto l647
					l648:
						position, tokenIndex, depth = position647, tokenIndex647, depth647
						if buffer[position] != rune('O') {
							goto l638
						}
						position++
					}
				l647:
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l650
						}
						position++
						goto l649
					l650:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						if buffer[position] != rune('N') {
							goto l638
						}
						position++
					}
				l649:
					{
						position651, tokenIndex651, depth651 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l652
						}
						position++
						goto l651
					l652:
						position, tokenIndex, depth = position651, tokenIndex651, depth651
						if buffer[position] != rune('D') {
							goto l638
						}
						position++
					}
				l651:
					{
						position653, tokenIndex653, depth653 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l654
						}
						position++
						goto l653
					l654:
						position, tokenIndex, depth = position653, tokenIndex653, depth653
						if buffer[position] != rune('S') {
							goto l638
						}
						position++
					}
				l653:
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction45]() {
					goto l638
				}
				depth--
				add(ruleSECONDS, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 60 StreamIdentifier <- <(<ident> Action46)> */
		func() bool {
			position655, tokenIndex655, depth655 := position, tokenIndex, depth
			{
				position656 := position
				depth++
				{
					position657 := position
					depth++
					if !_rules[ruleident]() {
						goto l655
					}
					depth--
					add(rulePegText, position657)
				}
				if !_rules[ruleAction46]() {
					goto l655
				}
				depth--
				add(ruleStreamIdentifier, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 61 SourceSinkType <- <(<ident> Action47)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				{
					position660 := position
					depth++
					if !_rules[ruleident]() {
						goto l658
					}
					depth--
					add(rulePegText, position660)
				}
				if !_rules[ruleAction47]() {
					goto l658
				}
				depth--
				add(ruleSourceSinkType, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 62 SourceSinkParamKey <- <(<ident> Action48)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				{
					position663 := position
					depth++
					if !_rules[ruleident]() {
						goto l661
					}
					depth--
					add(rulePegText, position663)
				}
				if !_rules[ruleAction48]() {
					goto l661
				}
				depth--
				add(ruleSourceSinkParamKey, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 63 SourceSinkParamVal <- <(<([a-z] / [A-Z] / [0-9] / '_')+> Action49)> */
		func() bool {
			position664, tokenIndex664, depth664 := position, tokenIndex, depth
			{
				position665 := position
				depth++
				{
					position666 := position
					depth++
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l670
						}
						position++
						goto l669
					l670:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l671
						}
						position++
						goto l669
					l671:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l672
						}
						position++
						goto l669
					l672:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if buffer[position] != rune('_') {
							goto l664
						}
						position++
					}
				l669:
				l667:
					{
						position668, tokenIndex668, depth668 := position, tokenIndex, depth
						{
							position673, tokenIndex673, depth673 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l674
							}
							position++
							goto l673
						l674:
							position, tokenIndex, depth = position673, tokenIndex673, depth673
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l675
							}
							position++
							goto l673
						l675:
							position, tokenIndex, depth = position673, tokenIndex673, depth673
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l676
							}
							position++
							goto l673
						l676:
							position, tokenIndex, depth = position673, tokenIndex673, depth673
							if buffer[position] != rune('_') {
								goto l668
							}
							position++
						}
					l673:
						goto l667
					l668:
						position, tokenIndex, depth = position668, tokenIndex668, depth668
					}
					depth--
					add(rulePegText, position666)
				}
				if !_rules[ruleAction49]() {
					goto l664
				}
				depth--
				add(ruleSourceSinkParamVal, position665)
			}
			return true
		l664:
			position, tokenIndex, depth = position664, tokenIndex664, depth664
			return false
		},
		/* 64 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action50)> */
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
						if buffer[position] != rune('o') {
							goto l681
						}
						position++
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if buffer[position] != rune('O') {
							goto l677
						}
						position++
					}
				l680:
					{
						position682, tokenIndex682, depth682 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l683
						}
						position++
						goto l682
					l683:
						position, tokenIndex, depth = position682, tokenIndex682, depth682
						if buffer[position] != rune('R') {
							goto l677
						}
						position++
					}
				l682:
					depth--
					add(rulePegText, position679)
				}
				if !_rules[ruleAction50]() {
					goto l677
				}
				depth--
				add(ruleOr, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 65 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action51)> */
		func() bool {
			position684, tokenIndex684, depth684 := position, tokenIndex, depth
			{
				position685 := position
				depth++
				{
					position686 := position
					depth++
					{
						position687, tokenIndex687, depth687 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l688
						}
						position++
						goto l687
					l688:
						position, tokenIndex, depth = position687, tokenIndex687, depth687
						if buffer[position] != rune('A') {
							goto l684
						}
						position++
					}
				l687:
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l690
						}
						position++
						goto l689
					l690:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
						if buffer[position] != rune('N') {
							goto l684
						}
						position++
					}
				l689:
					{
						position691, tokenIndex691, depth691 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l692
						}
						position++
						goto l691
					l692:
						position, tokenIndex, depth = position691, tokenIndex691, depth691
						if buffer[position] != rune('D') {
							goto l684
						}
						position++
					}
				l691:
					depth--
					add(rulePegText, position686)
				}
				if !_rules[ruleAction51]() {
					goto l684
				}
				depth--
				add(ruleAnd, position685)
			}
			return true
		l684:
			position, tokenIndex, depth = position684, tokenIndex684, depth684
			return false
		},
		/* 66 Equal <- <(<'='> Action52)> */
		func() bool {
			position693, tokenIndex693, depth693 := position, tokenIndex, depth
			{
				position694 := position
				depth++
				{
					position695 := position
					depth++
					if buffer[position] != rune('=') {
						goto l693
					}
					position++
					depth--
					add(rulePegText, position695)
				}
				if !_rules[ruleAction52]() {
					goto l693
				}
				depth--
				add(ruleEqual, position694)
			}
			return true
		l693:
			position, tokenIndex, depth = position693, tokenIndex693, depth693
			return false
		},
		/* 67 Less <- <(<'<'> Action53)> */
		func() bool {
			position696, tokenIndex696, depth696 := position, tokenIndex, depth
			{
				position697 := position
				depth++
				{
					position698 := position
					depth++
					if buffer[position] != rune('<') {
						goto l696
					}
					position++
					depth--
					add(rulePegText, position698)
				}
				if !_rules[ruleAction53]() {
					goto l696
				}
				depth--
				add(ruleLess, position697)
			}
			return true
		l696:
			position, tokenIndex, depth = position696, tokenIndex696, depth696
			return false
		},
		/* 68 LessOrEqual <- <(<('<' '=')> Action54)> */
		func() bool {
			position699, tokenIndex699, depth699 := position, tokenIndex, depth
			{
				position700 := position
				depth++
				{
					position701 := position
					depth++
					if buffer[position] != rune('<') {
						goto l699
					}
					position++
					if buffer[position] != rune('=') {
						goto l699
					}
					position++
					depth--
					add(rulePegText, position701)
				}
				if !_rules[ruleAction54]() {
					goto l699
				}
				depth--
				add(ruleLessOrEqual, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 69 Greater <- <(<'>'> Action55)> */
		func() bool {
			position702, tokenIndex702, depth702 := position, tokenIndex, depth
			{
				position703 := position
				depth++
				{
					position704 := position
					depth++
					if buffer[position] != rune('>') {
						goto l702
					}
					position++
					depth--
					add(rulePegText, position704)
				}
				if !_rules[ruleAction55]() {
					goto l702
				}
				depth--
				add(ruleGreater, position703)
			}
			return true
		l702:
			position, tokenIndex, depth = position702, tokenIndex702, depth702
			return false
		},
		/* 70 GreaterOrEqual <- <(<('>' '=')> Action56)> */
		func() bool {
			position705, tokenIndex705, depth705 := position, tokenIndex, depth
			{
				position706 := position
				depth++
				{
					position707 := position
					depth++
					if buffer[position] != rune('>') {
						goto l705
					}
					position++
					if buffer[position] != rune('=') {
						goto l705
					}
					position++
					depth--
					add(rulePegText, position707)
				}
				if !_rules[ruleAction56]() {
					goto l705
				}
				depth--
				add(ruleGreaterOrEqual, position706)
			}
			return true
		l705:
			position, tokenIndex, depth = position705, tokenIndex705, depth705
			return false
		},
		/* 71 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action57)> */
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
						if buffer[position] != rune('!') {
							goto l712
						}
						position++
						if buffer[position] != rune('=') {
							goto l712
						}
						position++
						goto l711
					l712:
						position, tokenIndex, depth = position711, tokenIndex711, depth711
						if buffer[position] != rune('<') {
							goto l708
						}
						position++
						if buffer[position] != rune('>') {
							goto l708
						}
						position++
					}
				l711:
					depth--
					add(rulePegText, position710)
				}
				if !_rules[ruleAction57]() {
					goto l708
				}
				depth--
				add(ruleNotEqual, position709)
			}
			return true
		l708:
			position, tokenIndex, depth = position708, tokenIndex708, depth708
			return false
		},
		/* 72 Plus <- <(<'+'> Action58)> */
		func() bool {
			position713, tokenIndex713, depth713 := position, tokenIndex, depth
			{
				position714 := position
				depth++
				{
					position715 := position
					depth++
					if buffer[position] != rune('+') {
						goto l713
					}
					position++
					depth--
					add(rulePegText, position715)
				}
				if !_rules[ruleAction58]() {
					goto l713
				}
				depth--
				add(rulePlus, position714)
			}
			return true
		l713:
			position, tokenIndex, depth = position713, tokenIndex713, depth713
			return false
		},
		/* 73 Minus <- <(<'-'> Action59)> */
		func() bool {
			position716, tokenIndex716, depth716 := position, tokenIndex, depth
			{
				position717 := position
				depth++
				{
					position718 := position
					depth++
					if buffer[position] != rune('-') {
						goto l716
					}
					position++
					depth--
					add(rulePegText, position718)
				}
				if !_rules[ruleAction59]() {
					goto l716
				}
				depth--
				add(ruleMinus, position717)
			}
			return true
		l716:
			position, tokenIndex, depth = position716, tokenIndex716, depth716
			return false
		},
		/* 74 Multiply <- <(<'*'> Action60)> */
		func() bool {
			position719, tokenIndex719, depth719 := position, tokenIndex, depth
			{
				position720 := position
				depth++
				{
					position721 := position
					depth++
					if buffer[position] != rune('*') {
						goto l719
					}
					position++
					depth--
					add(rulePegText, position721)
				}
				if !_rules[ruleAction60]() {
					goto l719
				}
				depth--
				add(ruleMultiply, position720)
			}
			return true
		l719:
			position, tokenIndex, depth = position719, tokenIndex719, depth719
			return false
		},
		/* 75 Divide <- <(<'/'> Action61)> */
		func() bool {
			position722, tokenIndex722, depth722 := position, tokenIndex, depth
			{
				position723 := position
				depth++
				{
					position724 := position
					depth++
					if buffer[position] != rune('/') {
						goto l722
					}
					position++
					depth--
					add(rulePegText, position724)
				}
				if !_rules[ruleAction61]() {
					goto l722
				}
				depth--
				add(ruleDivide, position723)
			}
			return true
		l722:
			position, tokenIndex, depth = position722, tokenIndex722, depth722
			return false
		},
		/* 76 Modulo <- <(<'%'> Action62)> */
		func() bool {
			position725, tokenIndex725, depth725 := position, tokenIndex, depth
			{
				position726 := position
				depth++
				{
					position727 := position
					depth++
					if buffer[position] != rune('%') {
						goto l725
					}
					position++
					depth--
					add(rulePegText, position727)
				}
				if !_rules[ruleAction62]() {
					goto l725
				}
				depth--
				add(ruleModulo, position726)
			}
			return true
		l725:
			position, tokenIndex, depth = position725, tokenIndex725, depth725
			return false
		},
		/* 77 Identifier <- <(<ident> Action63)> */
		func() bool {
			position728, tokenIndex728, depth728 := position, tokenIndex, depth
			{
				position729 := position
				depth++
				{
					position730 := position
					depth++
					if !_rules[ruleident]() {
						goto l728
					}
					depth--
					add(rulePegText, position730)
				}
				if !_rules[ruleAction63]() {
					goto l728
				}
				depth--
				add(ruleIdentifier, position729)
			}
			return true
		l728:
			position, tokenIndex, depth = position728, tokenIndex728, depth728
			return false
		},
		/* 78 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position731, tokenIndex731, depth731 := position, tokenIndex, depth
			{
				position732 := position
				depth++
				{
					position733, tokenIndex733, depth733 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l734
					}
					position++
					goto l733
				l734:
					position, tokenIndex, depth = position733, tokenIndex733, depth733
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l731
					}
					position++
				}
			l733:
			l735:
				{
					position736, tokenIndex736, depth736 := position, tokenIndex, depth
					{
						position737, tokenIndex737, depth737 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l738
						}
						position++
						goto l737
					l738:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l739
						}
						position++
						goto l737
					l739:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l740
						}
						position++
						goto l737
					l740:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
						if buffer[position] != rune('_') {
							goto l736
						}
						position++
					}
				l737:
					goto l735
				l736:
					position, tokenIndex, depth = position736, tokenIndex736, depth736
				}
				depth--
				add(ruleident, position732)
			}
			return true
		l731:
			position, tokenIndex, depth = position731, tokenIndex731, depth731
			return false
		},
		/* 79 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position742 := position
				depth++
			l743:
				{
					position744, tokenIndex744, depth744 := position, tokenIndex, depth
					{
						position745, tokenIndex745, depth745 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l746
						}
						position++
						goto l745
					l746:
						position, tokenIndex, depth = position745, tokenIndex745, depth745
						if buffer[position] != rune('\t') {
							goto l747
						}
						position++
						goto l745
					l747:
						position, tokenIndex, depth = position745, tokenIndex745, depth745
						if buffer[position] != rune('\n') {
							goto l744
						}
						position++
					}
				l745:
					goto l743
				l744:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
				}
				depth--
				add(rulesp, position742)
			}
			return true
		},
		/* 81 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 82 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 83 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 84 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 85 Action4 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 86 Action5 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 87 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 88 Action7 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 90 Action8 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 91 Action9 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 92 Action10 <- <{
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
		/* 93 Action11 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 94 Action12 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 95 Action13 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 96 Action14 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 97 Action15 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 98 Action16 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 99 Action17 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 100 Action18 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 101 Action19 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 102 Action20 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 103 Action21 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 104 Action22 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 105 Action23 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 106 Action24 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 107 Action25 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 108 Action26 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 109 Action27 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 110 Action28 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 111 Action29 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 112 Action30 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 113 Action31 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 114 Action32 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 115 Action33 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 116 Action34 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 117 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 118 Action36 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 119 Action37 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 120 Action38 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 121 Action39 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 122 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 123 Action41 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 124 Action42 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 125 Action43 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 126 Action44 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 127 Action45 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 128 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 129 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 130 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 131 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamVal(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 132 Action50 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 133 Action51 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 134 Action52 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 135 Action53 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 136 Action54 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 137 Action55 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 138 Action56 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 139 Action57 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 140 Action58 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 141 Action59 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 142 Action60 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 143 Action61 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 144 Action62 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 145 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
	}
	p.rules = _rules
}

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
	ruleAction65

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
	"Action65",

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
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

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

			p.PushComponent(begin, end, Or)

		case ruleAction53:

			p.PushComponent(begin, end, And)

		case ruleAction54:

			p.PushComponent(begin, end, Equal)

		case ruleAction55:

			p.PushComponent(begin, end, Less)

		case ruleAction56:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction57:

			p.PushComponent(begin, end, Greater)

		case ruleAction58:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction59:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction60:

			p.PushComponent(begin, end, Plus)

		case ruleAction61:

			p.PushComponent(begin, end, Minus)

		case ruleAction62:

			p.PushComponent(begin, end, Multiply)

		case ruleAction63:

			p.PushComponent(begin, end, Divide)

		case ruleAction64:

			p.PushComponent(begin, end, Modulo)

		case ruleAction65:

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
		/* 42 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowMeta / RowValue / Literal)> */
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
					if !_rules[ruleRowMeta]() {
						goto l443
					}
					goto l439
				l443:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleRowValue]() {
						goto l444
					}
					goto l439
				l444:
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
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l445
				}
				if !_rules[rulesp]() {
					goto l445
				}
				if buffer[position] != rune('(') {
					goto l445
				}
				position++
				if !_rules[rulesp]() {
					goto l445
				}
				if !_rules[ruleFuncParams]() {
					goto l445
				}
				if !_rules[rulesp]() {
					goto l445
				}
				if buffer[position] != rune(')') {
					goto l445
				}
				position++
				if !_rules[ruleAction32]() {
					goto l445
				}
				depth--
				add(ruleFuncApp, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 44 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action33)> */
		func() bool {
			position447, tokenIndex447, depth447 := position, tokenIndex, depth
			{
				position448 := position
				depth++
				{
					position449 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l447
					}
					if !_rules[rulesp]() {
						goto l447
					}
				l450:
					{
						position451, tokenIndex451, depth451 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l451
						}
						position++
						if !_rules[rulesp]() {
							goto l451
						}
						if !_rules[ruleExpression]() {
							goto l451
						}
						goto l450
					l451:
						position, tokenIndex, depth = position451, tokenIndex451, depth451
					}
					depth--
					add(rulePegText, position449)
				}
				if !_rules[ruleAction33]() {
					goto l447
				}
				depth--
				add(ruleFuncParams, position448)
			}
			return true
		l447:
			position, tokenIndex, depth = position447, tokenIndex447, depth447
			return false
		},
		/* 45 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position452, tokenIndex452, depth452 := position, tokenIndex, depth
			{
				position453 := position
				depth++
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l455
					}
					goto l454
				l455:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if !_rules[ruleNumericLiteral]() {
						goto l456
					}
					goto l454
				l456:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if !_rules[ruleStringLiteral]() {
						goto l452
					}
				}
			l454:
				depth--
				add(ruleLiteral, position453)
			}
			return true
		l452:
			position, tokenIndex, depth = position452, tokenIndex452, depth452
			return false
		},
		/* 46 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l460
					}
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleNotEqual]() {
						goto l461
					}
					goto l459
				l461:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleLessOrEqual]() {
						goto l462
					}
					goto l459
				l462:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleLess]() {
						goto l463
					}
					goto l459
				l463:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleGreaterOrEqual]() {
						goto l464
					}
					goto l459
				l464:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleGreater]() {
						goto l465
					}
					goto l459
				l465:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleNotEqual]() {
						goto l457
					}
				}
			l459:
				depth--
				add(ruleComparisonOp, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 47 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position466, tokenIndex466, depth466 := position, tokenIndex, depth
			{
				position467 := position
				depth++
				{
					position468, tokenIndex468, depth468 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l469
					}
					goto l468
				l469:
					position, tokenIndex, depth = position468, tokenIndex468, depth468
					if !_rules[ruleMinus]() {
						goto l466
					}
				}
			l468:
				depth--
				add(rulePlusMinusOp, position467)
			}
			return true
		l466:
			position, tokenIndex, depth = position466, tokenIndex466, depth466
			return false
		},
		/* 48 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position470, tokenIndex470, depth470 := position, tokenIndex, depth
			{
				position471 := position
				depth++
				{
					position472, tokenIndex472, depth472 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l473
					}
					goto l472
				l473:
					position, tokenIndex, depth = position472, tokenIndex472, depth472
					if !_rules[ruleDivide]() {
						goto l474
					}
					goto l472
				l474:
					position, tokenIndex, depth = position472, tokenIndex472, depth472
					if !_rules[ruleModulo]() {
						goto l470
					}
				}
			l472:
				depth--
				add(ruleMultDivOp, position471)
			}
			return true
		l470:
			position, tokenIndex, depth = position470, tokenIndex470, depth470
			return false
		},
		/* 49 Stream <- <(<ident> Action34)> */
		func() bool {
			position475, tokenIndex475, depth475 := position, tokenIndex, depth
			{
				position476 := position
				depth++
				{
					position477 := position
					depth++
					if !_rules[ruleident]() {
						goto l475
					}
					depth--
					add(rulePegText, position477)
				}
				if !_rules[ruleAction34]() {
					goto l475
				}
				depth--
				add(ruleStream, position476)
			}
			return true
		l475:
			position, tokenIndex, depth = position475, tokenIndex475, depth475
			return false
		},
		/* 50 RowMeta <- <RowTimestamp> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l478
				}
				depth--
				add(ruleRowMeta, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 51 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action35)> */
		func() bool {
			position480, tokenIndex480, depth480 := position, tokenIndex, depth
			{
				position481 := position
				depth++
				{
					position482 := position
					depth++
					{
						position483, tokenIndex483, depth483 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l483
						}
						if buffer[position] != rune(':') {
							goto l483
						}
						position++
						goto l484
					l483:
						position, tokenIndex, depth = position483, tokenIndex483, depth483
					}
				l484:
					if buffer[position] != rune('t') {
						goto l480
					}
					position++
					if buffer[position] != rune('s') {
						goto l480
					}
					position++
					if buffer[position] != rune('(') {
						goto l480
					}
					position++
					if buffer[position] != rune(')') {
						goto l480
					}
					position++
					depth--
					add(rulePegText, position482)
				}
				if !_rules[ruleAction35]() {
					goto l480
				}
				depth--
				add(ruleRowTimestamp, position481)
			}
			return true
		l480:
			position, tokenIndex, depth = position480, tokenIndex480, depth480
			return false
		},
		/* 52 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action36)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				{
					position487 := position
					depth++
					{
						position488, tokenIndex488, depth488 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l488
						}
						if buffer[position] != rune(':') {
							goto l488
						}
						position++
						goto l489
					l488:
						position, tokenIndex, depth = position488, tokenIndex488, depth488
					}
				l489:
					{
						position490, tokenIndex490, depth490 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l491
						}
						position++
						goto l490
					l491:
						position, tokenIndex, depth = position490, tokenIndex490, depth490
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l485
						}
						position++
					}
				l490:
				l492:
					{
						position493, tokenIndex493, depth493 := position, tokenIndex, depth
						{
							position494, tokenIndex494, depth494 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l495
							}
							position++
							goto l494
						l495:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l496
							}
							position++
							goto l494
						l496:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l497
							}
							position++
							goto l494
						l497:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							if buffer[position] != rune('_') {
								goto l498
							}
							position++
							goto l494
						l498:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							if buffer[position] != rune('.') {
								goto l493
							}
							position++
						}
					l494:
						goto l492
					l493:
						position, tokenIndex, depth = position493, tokenIndex493, depth493
					}
					depth--
					add(rulePegText, position487)
				}
				if !_rules[ruleAction36]() {
					goto l485
				}
				depth--
				add(ruleRowValue, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 53 NumericLiteral <- <(<('-'? [0-9]+)> Action37)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				{
					position501 := position
					depth++
					{
						position502, tokenIndex502, depth502 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l502
						}
						position++
						goto l503
					l502:
						position, tokenIndex, depth = position502, tokenIndex502, depth502
					}
				l503:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l499
					}
					position++
				l504:
					{
						position505, tokenIndex505, depth505 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l505
						}
						position++
						goto l504
					l505:
						position, tokenIndex, depth = position505, tokenIndex505, depth505
					}
					depth--
					add(rulePegText, position501)
				}
				if !_rules[ruleAction37]() {
					goto l499
				}
				depth--
				add(ruleNumericLiteral, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 54 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action38)> */
		func() bool {
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				{
					position508 := position
					depth++
					{
						position509, tokenIndex509, depth509 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l509
						}
						position++
						goto l510
					l509:
						position, tokenIndex, depth = position509, tokenIndex509, depth509
					}
				l510:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l506
					}
					position++
				l511:
					{
						position512, tokenIndex512, depth512 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l512
						}
						position++
						goto l511
					l512:
						position, tokenIndex, depth = position512, tokenIndex512, depth512
					}
					if buffer[position] != rune('.') {
						goto l506
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l506
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
					add(rulePegText, position508)
				}
				if !_rules[ruleAction38]() {
					goto l506
				}
				depth--
				add(ruleFloatLiteral, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 55 Function <- <(<ident> Action39)> */
		func() bool {
			position515, tokenIndex515, depth515 := position, tokenIndex, depth
			{
				position516 := position
				depth++
				{
					position517 := position
					depth++
					if !_rules[ruleident]() {
						goto l515
					}
					depth--
					add(rulePegText, position517)
				}
				if !_rules[ruleAction39]() {
					goto l515
				}
				depth--
				add(ruleFunction, position516)
			}
			return true
		l515:
			position, tokenIndex, depth = position515, tokenIndex515, depth515
			return false
		},
		/* 56 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position518, tokenIndex518, depth518 := position, tokenIndex, depth
			{
				position519 := position
				depth++
				{
					position520, tokenIndex520, depth520 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l521
					}
					goto l520
				l521:
					position, tokenIndex, depth = position520, tokenIndex520, depth520
					if !_rules[ruleFALSE]() {
						goto l518
					}
				}
			l520:
				depth--
				add(ruleBooleanLiteral, position519)
			}
			return true
		l518:
			position, tokenIndex, depth = position518, tokenIndex518, depth518
			return false
		},
		/* 57 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action40)> */
		func() bool {
			position522, tokenIndex522, depth522 := position, tokenIndex, depth
			{
				position523 := position
				depth++
				{
					position524 := position
					depth++
					{
						position525, tokenIndex525, depth525 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l526
						}
						position++
						goto l525
					l526:
						position, tokenIndex, depth = position525, tokenIndex525, depth525
						if buffer[position] != rune('T') {
							goto l522
						}
						position++
					}
				l525:
					{
						position527, tokenIndex527, depth527 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l528
						}
						position++
						goto l527
					l528:
						position, tokenIndex, depth = position527, tokenIndex527, depth527
						if buffer[position] != rune('R') {
							goto l522
						}
						position++
					}
				l527:
					{
						position529, tokenIndex529, depth529 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l530
						}
						position++
						goto l529
					l530:
						position, tokenIndex, depth = position529, tokenIndex529, depth529
						if buffer[position] != rune('U') {
							goto l522
						}
						position++
					}
				l529:
					{
						position531, tokenIndex531, depth531 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l532
						}
						position++
						goto l531
					l532:
						position, tokenIndex, depth = position531, tokenIndex531, depth531
						if buffer[position] != rune('E') {
							goto l522
						}
						position++
					}
				l531:
					depth--
					add(rulePegText, position524)
				}
				if !_rules[ruleAction40]() {
					goto l522
				}
				depth--
				add(ruleTRUE, position523)
			}
			return true
		l522:
			position, tokenIndex, depth = position522, tokenIndex522, depth522
			return false
		},
		/* 58 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action41)> */
		func() bool {
			position533, tokenIndex533, depth533 := position, tokenIndex, depth
			{
				position534 := position
				depth++
				{
					position535 := position
					depth++
					{
						position536, tokenIndex536, depth536 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l537
						}
						position++
						goto l536
					l537:
						position, tokenIndex, depth = position536, tokenIndex536, depth536
						if buffer[position] != rune('F') {
							goto l533
						}
						position++
					}
				l536:
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l539
						}
						position++
						goto l538
					l539:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
						if buffer[position] != rune('A') {
							goto l533
						}
						position++
					}
				l538:
					{
						position540, tokenIndex540, depth540 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l541
						}
						position++
						goto l540
					l541:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
						if buffer[position] != rune('L') {
							goto l533
						}
						position++
					}
				l540:
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l543
						}
						position++
						goto l542
					l543:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
						if buffer[position] != rune('S') {
							goto l533
						}
						position++
					}
				l542:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l545
						}
						position++
						goto l544
					l545:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
						if buffer[position] != rune('E') {
							goto l533
						}
						position++
					}
				l544:
					depth--
					add(rulePegText, position535)
				}
				if !_rules[ruleAction41]() {
					goto l533
				}
				depth--
				add(ruleFALSE, position534)
			}
			return true
		l533:
			position, tokenIndex, depth = position533, tokenIndex533, depth533
			return false
		},
		/* 59 Wildcard <- <(<'*'> Action42)> */
		func() bool {
			position546, tokenIndex546, depth546 := position, tokenIndex, depth
			{
				position547 := position
				depth++
				{
					position548 := position
					depth++
					if buffer[position] != rune('*') {
						goto l546
					}
					position++
					depth--
					add(rulePegText, position548)
				}
				if !_rules[ruleAction42]() {
					goto l546
				}
				depth--
				add(ruleWildcard, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 60 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action43)> */
		func() bool {
			position549, tokenIndex549, depth549 := position, tokenIndex, depth
			{
				position550 := position
				depth++
				{
					position551 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l549
					}
					position++
				l552:
					{
						position553, tokenIndex553, depth553 := position, tokenIndex, depth
						{
							position554, tokenIndex554, depth554 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l555
							}
							position++
							if buffer[position] != rune('\'') {
								goto l555
							}
							position++
							goto l554
						l555:
							position, tokenIndex, depth = position554, tokenIndex554, depth554
							{
								position556, tokenIndex556, depth556 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l556
								}
								position++
								goto l553
							l556:
								position, tokenIndex, depth = position556, tokenIndex556, depth556
							}
							if !matchDot() {
								goto l553
							}
						}
					l554:
						goto l552
					l553:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
					}
					if buffer[position] != rune('\'') {
						goto l549
					}
					position++
					depth--
					add(rulePegText, position551)
				}
				if !_rules[ruleAction43]() {
					goto l549
				}
				depth--
				add(ruleStringLiteral, position550)
			}
			return true
		l549:
			position, tokenIndex, depth = position549, tokenIndex549, depth549
			return false
		},
		/* 61 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action44)> */
		func() bool {
			position557, tokenIndex557, depth557 := position, tokenIndex, depth
			{
				position558 := position
				depth++
				{
					position559 := position
					depth++
					{
						position560, tokenIndex560, depth560 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l561
						}
						position++
						goto l560
					l561:
						position, tokenIndex, depth = position560, tokenIndex560, depth560
						if buffer[position] != rune('I') {
							goto l557
						}
						position++
					}
				l560:
					{
						position562, tokenIndex562, depth562 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l563
						}
						position++
						goto l562
					l563:
						position, tokenIndex, depth = position562, tokenIndex562, depth562
						if buffer[position] != rune('S') {
							goto l557
						}
						position++
					}
				l562:
					{
						position564, tokenIndex564, depth564 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l565
						}
						position++
						goto l564
					l565:
						position, tokenIndex, depth = position564, tokenIndex564, depth564
						if buffer[position] != rune('T') {
							goto l557
						}
						position++
					}
				l564:
					{
						position566, tokenIndex566, depth566 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l567
						}
						position++
						goto l566
					l567:
						position, tokenIndex, depth = position566, tokenIndex566, depth566
						if buffer[position] != rune('R') {
							goto l557
						}
						position++
					}
				l566:
					{
						position568, tokenIndex568, depth568 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l569
						}
						position++
						goto l568
					l569:
						position, tokenIndex, depth = position568, tokenIndex568, depth568
						if buffer[position] != rune('E') {
							goto l557
						}
						position++
					}
				l568:
					{
						position570, tokenIndex570, depth570 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l571
						}
						position++
						goto l570
					l571:
						position, tokenIndex, depth = position570, tokenIndex570, depth570
						if buffer[position] != rune('A') {
							goto l557
						}
						position++
					}
				l570:
					{
						position572, tokenIndex572, depth572 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l573
						}
						position++
						goto l572
					l573:
						position, tokenIndex, depth = position572, tokenIndex572, depth572
						if buffer[position] != rune('M') {
							goto l557
						}
						position++
					}
				l572:
					depth--
					add(rulePegText, position559)
				}
				if !_rules[ruleAction44]() {
					goto l557
				}
				depth--
				add(ruleISTREAM, position558)
			}
			return true
		l557:
			position, tokenIndex, depth = position557, tokenIndex557, depth557
			return false
		},
		/* 62 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action45)> */
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
						if buffer[position] != rune('d') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('D') {
							goto l574
						}
						position++
					}
				l577:
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l580
						}
						position++
						goto l579
					l580:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
						if buffer[position] != rune('S') {
							goto l574
						}
						position++
					}
				l579:
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l582
						}
						position++
						goto l581
					l582:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
						if buffer[position] != rune('T') {
							goto l574
						}
						position++
					}
				l581:
					{
						position583, tokenIndex583, depth583 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l584
						}
						position++
						goto l583
					l584:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
						if buffer[position] != rune('R') {
							goto l574
						}
						position++
					}
				l583:
					{
						position585, tokenIndex585, depth585 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l586
						}
						position++
						goto l585
					l586:
						position, tokenIndex, depth = position585, tokenIndex585, depth585
						if buffer[position] != rune('E') {
							goto l574
						}
						position++
					}
				l585:
					{
						position587, tokenIndex587, depth587 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l588
						}
						position++
						goto l587
					l588:
						position, tokenIndex, depth = position587, tokenIndex587, depth587
						if buffer[position] != rune('A') {
							goto l574
						}
						position++
					}
				l587:
					{
						position589, tokenIndex589, depth589 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l590
						}
						position++
						goto l589
					l590:
						position, tokenIndex, depth = position589, tokenIndex589, depth589
						if buffer[position] != rune('M') {
							goto l574
						}
						position++
					}
				l589:
					depth--
					add(rulePegText, position576)
				}
				if !_rules[ruleAction45]() {
					goto l574
				}
				depth--
				add(ruleDSTREAM, position575)
			}
			return true
		l574:
			position, tokenIndex, depth = position574, tokenIndex574, depth574
			return false
		},
		/* 63 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action46)> */
		func() bool {
			position591, tokenIndex591, depth591 := position, tokenIndex, depth
			{
				position592 := position
				depth++
				{
					position593 := position
					depth++
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
						if buffer[position] != rune('R') {
							goto l591
						}
						position++
					}
				l594:
					{
						position596, tokenIndex596, depth596 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l597
						}
						position++
						goto l596
					l597:
						position, tokenIndex, depth = position596, tokenIndex596, depth596
						if buffer[position] != rune('S') {
							goto l591
						}
						position++
					}
				l596:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l599
						}
						position++
						goto l598
					l599:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
						if buffer[position] != rune('T') {
							goto l591
						}
						position++
					}
				l598:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l601
						}
						position++
						goto l600
					l601:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
						if buffer[position] != rune('R') {
							goto l591
						}
						position++
					}
				l600:
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l603
						}
						position++
						goto l602
					l603:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if buffer[position] != rune('E') {
							goto l591
						}
						position++
					}
				l602:
					{
						position604, tokenIndex604, depth604 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l605
						}
						position++
						goto l604
					l605:
						position, tokenIndex, depth = position604, tokenIndex604, depth604
						if buffer[position] != rune('A') {
							goto l591
						}
						position++
					}
				l604:
					{
						position606, tokenIndex606, depth606 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l607
						}
						position++
						goto l606
					l607:
						position, tokenIndex, depth = position606, tokenIndex606, depth606
						if buffer[position] != rune('M') {
							goto l591
						}
						position++
					}
				l606:
					depth--
					add(rulePegText, position593)
				}
				if !_rules[ruleAction46]() {
					goto l591
				}
				depth--
				add(ruleRSTREAM, position592)
			}
			return true
		l591:
			position, tokenIndex, depth = position591, tokenIndex591, depth591
			return false
		},
		/* 64 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action47)> */
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
						if buffer[position] != rune('u') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('U') {
							goto l608
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('P') {
							goto l608
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('L') {
							goto l608
						}
						position++
					}
				l617:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('E') {
							goto l608
						}
						position++
					}
				l619:
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l622
						}
						position++
						goto l621
					l622:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
						if buffer[position] != rune('S') {
							goto l608
						}
						position++
					}
				l621:
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction47]() {
					goto l608
				}
				depth--
				add(ruleTUPLES, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 65 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action48)> */
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
						if buffer[position] != rune('s') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('S') {
							goto l623
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('E') {
							goto l623
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('C') {
							goto l623
						}
						position++
					}
				l630:
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l633
						}
						position++
						goto l632
					l633:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
						if buffer[position] != rune('O') {
							goto l623
						}
						position++
					}
				l632:
					{
						position634, tokenIndex634, depth634 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l635
						}
						position++
						goto l634
					l635:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
						if buffer[position] != rune('N') {
							goto l623
						}
						position++
					}
				l634:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l637
						}
						position++
						goto l636
					l637:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
						if buffer[position] != rune('D') {
							goto l623
						}
						position++
					}
				l636:
					{
						position638, tokenIndex638, depth638 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l639
						}
						position++
						goto l638
					l639:
						position, tokenIndex, depth = position638, tokenIndex638, depth638
						if buffer[position] != rune('S') {
							goto l623
						}
						position++
					}
				l638:
					depth--
					add(rulePegText, position625)
				}
				if !_rules[ruleAction48]() {
					goto l623
				}
				depth--
				add(ruleSECONDS, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 66 StreamIdentifier <- <(<ident> Action49)> */
		func() bool {
			position640, tokenIndex640, depth640 := position, tokenIndex, depth
			{
				position641 := position
				depth++
				{
					position642 := position
					depth++
					if !_rules[ruleident]() {
						goto l640
					}
					depth--
					add(rulePegText, position642)
				}
				if !_rules[ruleAction49]() {
					goto l640
				}
				depth--
				add(ruleStreamIdentifier, position641)
			}
			return true
		l640:
			position, tokenIndex, depth = position640, tokenIndex640, depth640
			return false
		},
		/* 67 SourceSinkType <- <(<ident> Action50)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645 := position
					depth++
					if !_rules[ruleident]() {
						goto l643
					}
					depth--
					add(rulePegText, position645)
				}
				if !_rules[ruleAction50]() {
					goto l643
				}
				depth--
				add(ruleSourceSinkType, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 68 SourceSinkParamKey <- <(<ident> Action51)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				{
					position648 := position
					depth++
					if !_rules[ruleident]() {
						goto l646
					}
					depth--
					add(rulePegText, position648)
				}
				if !_rules[ruleAction51]() {
					goto l646
				}
				depth--
				add(ruleSourceSinkParamKey, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 69 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action52)> */
		func() bool {
			position649, tokenIndex649, depth649 := position, tokenIndex, depth
			{
				position650 := position
				depth++
				{
					position651 := position
					depth++
					{
						position652, tokenIndex652, depth652 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l653
						}
						position++
						goto l652
					l653:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
						if buffer[position] != rune('O') {
							goto l649
						}
						position++
					}
				l652:
					{
						position654, tokenIndex654, depth654 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l655
						}
						position++
						goto l654
					l655:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
						if buffer[position] != rune('R') {
							goto l649
						}
						position++
					}
				l654:
					depth--
					add(rulePegText, position651)
				}
				if !_rules[ruleAction52]() {
					goto l649
				}
				depth--
				add(ruleOr, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 70 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action53)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				{
					position658 := position
					depth++
					{
						position659, tokenIndex659, depth659 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l660
						}
						position++
						goto l659
					l660:
						position, tokenIndex, depth = position659, tokenIndex659, depth659
						if buffer[position] != rune('A') {
							goto l656
						}
						position++
					}
				l659:
					{
						position661, tokenIndex661, depth661 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l662
						}
						position++
						goto l661
					l662:
						position, tokenIndex, depth = position661, tokenIndex661, depth661
						if buffer[position] != rune('N') {
							goto l656
						}
						position++
					}
				l661:
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
							goto l656
						}
						position++
					}
				l663:
					depth--
					add(rulePegText, position658)
				}
				if !_rules[ruleAction53]() {
					goto l656
				}
				depth--
				add(ruleAnd, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 71 Equal <- <(<'='> Action54)> */
		func() bool {
			position665, tokenIndex665, depth665 := position, tokenIndex, depth
			{
				position666 := position
				depth++
				{
					position667 := position
					depth++
					if buffer[position] != rune('=') {
						goto l665
					}
					position++
					depth--
					add(rulePegText, position667)
				}
				if !_rules[ruleAction54]() {
					goto l665
				}
				depth--
				add(ruleEqual, position666)
			}
			return true
		l665:
			position, tokenIndex, depth = position665, tokenIndex665, depth665
			return false
		},
		/* 72 Less <- <(<'<'> Action55)> */
		func() bool {
			position668, tokenIndex668, depth668 := position, tokenIndex, depth
			{
				position669 := position
				depth++
				{
					position670 := position
					depth++
					if buffer[position] != rune('<') {
						goto l668
					}
					position++
					depth--
					add(rulePegText, position670)
				}
				if !_rules[ruleAction55]() {
					goto l668
				}
				depth--
				add(ruleLess, position669)
			}
			return true
		l668:
			position, tokenIndex, depth = position668, tokenIndex668, depth668
			return false
		},
		/* 73 LessOrEqual <- <(<('<' '=')> Action56)> */
		func() bool {
			position671, tokenIndex671, depth671 := position, tokenIndex, depth
			{
				position672 := position
				depth++
				{
					position673 := position
					depth++
					if buffer[position] != rune('<') {
						goto l671
					}
					position++
					if buffer[position] != rune('=') {
						goto l671
					}
					position++
					depth--
					add(rulePegText, position673)
				}
				if !_rules[ruleAction56]() {
					goto l671
				}
				depth--
				add(ruleLessOrEqual, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 74 Greater <- <(<'>'> Action57)> */
		func() bool {
			position674, tokenIndex674, depth674 := position, tokenIndex, depth
			{
				position675 := position
				depth++
				{
					position676 := position
					depth++
					if buffer[position] != rune('>') {
						goto l674
					}
					position++
					depth--
					add(rulePegText, position676)
				}
				if !_rules[ruleAction57]() {
					goto l674
				}
				depth--
				add(ruleGreater, position675)
			}
			return true
		l674:
			position, tokenIndex, depth = position674, tokenIndex674, depth674
			return false
		},
		/* 75 GreaterOrEqual <- <(<('>' '=')> Action58)> */
		func() bool {
			position677, tokenIndex677, depth677 := position, tokenIndex, depth
			{
				position678 := position
				depth++
				{
					position679 := position
					depth++
					if buffer[position] != rune('>') {
						goto l677
					}
					position++
					if buffer[position] != rune('=') {
						goto l677
					}
					position++
					depth--
					add(rulePegText, position679)
				}
				if !_rules[ruleAction58]() {
					goto l677
				}
				depth--
				add(ruleGreaterOrEqual, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 76 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action59)> */
		func() bool {
			position680, tokenIndex680, depth680 := position, tokenIndex, depth
			{
				position681 := position
				depth++
				{
					position682 := position
					depth++
					{
						position683, tokenIndex683, depth683 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l684
						}
						position++
						if buffer[position] != rune('=') {
							goto l684
						}
						position++
						goto l683
					l684:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
						if buffer[position] != rune('<') {
							goto l680
						}
						position++
						if buffer[position] != rune('>') {
							goto l680
						}
						position++
					}
				l683:
					depth--
					add(rulePegText, position682)
				}
				if !_rules[ruleAction59]() {
					goto l680
				}
				depth--
				add(ruleNotEqual, position681)
			}
			return true
		l680:
			position, tokenIndex, depth = position680, tokenIndex680, depth680
			return false
		},
		/* 77 Plus <- <(<'+'> Action60)> */
		func() bool {
			position685, tokenIndex685, depth685 := position, tokenIndex, depth
			{
				position686 := position
				depth++
				{
					position687 := position
					depth++
					if buffer[position] != rune('+') {
						goto l685
					}
					position++
					depth--
					add(rulePegText, position687)
				}
				if !_rules[ruleAction60]() {
					goto l685
				}
				depth--
				add(rulePlus, position686)
			}
			return true
		l685:
			position, tokenIndex, depth = position685, tokenIndex685, depth685
			return false
		},
		/* 78 Minus <- <(<'-'> Action61)> */
		func() bool {
			position688, tokenIndex688, depth688 := position, tokenIndex, depth
			{
				position689 := position
				depth++
				{
					position690 := position
					depth++
					if buffer[position] != rune('-') {
						goto l688
					}
					position++
					depth--
					add(rulePegText, position690)
				}
				if !_rules[ruleAction61]() {
					goto l688
				}
				depth--
				add(ruleMinus, position689)
			}
			return true
		l688:
			position, tokenIndex, depth = position688, tokenIndex688, depth688
			return false
		},
		/* 79 Multiply <- <(<'*'> Action62)> */
		func() bool {
			position691, tokenIndex691, depth691 := position, tokenIndex, depth
			{
				position692 := position
				depth++
				{
					position693 := position
					depth++
					if buffer[position] != rune('*') {
						goto l691
					}
					position++
					depth--
					add(rulePegText, position693)
				}
				if !_rules[ruleAction62]() {
					goto l691
				}
				depth--
				add(ruleMultiply, position692)
			}
			return true
		l691:
			position, tokenIndex, depth = position691, tokenIndex691, depth691
			return false
		},
		/* 80 Divide <- <(<'/'> Action63)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				{
					position696 := position
					depth++
					if buffer[position] != rune('/') {
						goto l694
					}
					position++
					depth--
					add(rulePegText, position696)
				}
				if !_rules[ruleAction63]() {
					goto l694
				}
				depth--
				add(ruleDivide, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 81 Modulo <- <(<'%'> Action64)> */
		func() bool {
			position697, tokenIndex697, depth697 := position, tokenIndex, depth
			{
				position698 := position
				depth++
				{
					position699 := position
					depth++
					if buffer[position] != rune('%') {
						goto l697
					}
					position++
					depth--
					add(rulePegText, position699)
				}
				if !_rules[ruleAction64]() {
					goto l697
				}
				depth--
				add(ruleModulo, position698)
			}
			return true
		l697:
			position, tokenIndex, depth = position697, tokenIndex697, depth697
			return false
		},
		/* 82 Identifier <- <(<ident> Action65)> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				{
					position702 := position
					depth++
					if !_rules[ruleident]() {
						goto l700
					}
					depth--
					add(rulePegText, position702)
				}
				if !_rules[ruleAction65]() {
					goto l700
				}
				depth--
				add(ruleIdentifier, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 83 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position703, tokenIndex703, depth703 := position, tokenIndex, depth
			{
				position704 := position
				depth++
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
						goto l703
					}
					position++
				}
			l705:
			l707:
				{
					position708, tokenIndex708, depth708 := position, tokenIndex, depth
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l710
						}
						position++
						goto l709
					l710:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l711
						}
						position++
						goto l709
					l711:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l712
						}
						position++
						goto l709
					l712:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if buffer[position] != rune('_') {
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
				add(ruleident, position704)
			}
			return true
		l703:
			position, tokenIndex, depth = position703, tokenIndex703, depth703
			return false
		},
		/* 84 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position714 := position
				depth++
			l715:
				{
					position716, tokenIndex716, depth716 := position, tokenIndex, depth
					{
						position717, tokenIndex717, depth717 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l718
						}
						position++
						goto l717
					l718:
						position, tokenIndex, depth = position717, tokenIndex717, depth717
						if buffer[position] != rune('\t') {
							goto l719
						}
						position++
						goto l717
					l719:
						position, tokenIndex, depth = position717, tokenIndex717, depth717
						if buffer[position] != rune('\n') {
							goto l716
						}
						position++
					}
				l717:
					goto l715
				l716:
					position, tokenIndex, depth = position716, tokenIndex716, depth716
				}
				depth--
				add(rulesp, position714)
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
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		nil,
		/* 93 Action6 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 94 Action7 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 95 Action8 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 96 Action9 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 97 Action10 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 98 Action11 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 99 Action12 <- <{
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
		/* 100 Action13 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 101 Action14 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 102 Action15 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 103 Action16 <- <{
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
		/* 104 Action17 <- <{
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
		/* 105 Action18 <- <{
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
		/* 106 Action19 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 107 Action20 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 108 Action21 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 109 Action22 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 110 Action23 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 111 Action24 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 112 Action25 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 113 Action26 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 114 Action27 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 115 Action28 <- <{
		    p.AssembleBinaryOperation(begin, end)
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
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 120 Action33 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 121 Action34 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 122 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 123 Action36 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 124 Action37 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 125 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 126 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 127 Action40 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 128 Action41 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 129 Action42 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 130 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 131 Action44 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 132 Action45 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 133 Action46 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 134 Action47 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 135 Action48 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 136 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 137 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 138 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 139 Action52 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 140 Action53 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 141 Action54 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 142 Action55 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 143 Action56 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 144 Action57 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 145 Action58 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 146 Action59 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 147 Action60 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 148 Action61 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 149 Action62 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 150 Action63 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 151 Action64 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 152 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
	}
	p.rules = _rules
}

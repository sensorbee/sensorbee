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
	rules  [152]func() bool
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
		/* 2 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter? sp Projections sp DefWindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
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
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if !_rules[ruleEmitter]() {
						goto l30
					}
					goto l31
				l30:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
				}
			l31:
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
			position32, tokenIndex32, depth32 := position, tokenIndex, depth
			{
				position33 := position
				depth++
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l35
					}
					position++
					goto l34
				l35:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if buffer[position] != rune('C') {
						goto l32
					}
					position++
				}
			l34:
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l37
					}
					position++
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if buffer[position] != rune('R') {
						goto l32
					}
					position++
				}
			l36:
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l39
					}
					position++
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if buffer[position] != rune('E') {
						goto l32
					}
					position++
				}
			l38:
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('A') {
						goto l32
					}
					position++
				}
			l40:
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('T') {
						goto l32
					}
					position++
				}
			l42:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l45
					}
					position++
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if buffer[position] != rune('E') {
						goto l32
					}
					position++
				}
			l44:
				if !_rules[rulesp]() {
					goto l32
				}
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('S') {
						goto l32
					}
					position++
				}
			l46:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('T') {
						goto l32
					}
					position++
				}
			l48:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('R') {
						goto l32
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('E') {
						goto l32
					}
					position++
				}
			l52:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('A') {
						goto l32
					}
					position++
				}
			l54:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('M') {
						goto l32
					}
					position++
				}
			l56:
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l32
				}
				if !_rules[rulesp]() {
					goto l32
				}
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('A') {
						goto l32
					}
					position++
				}
			l58:
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
						goto l32
					}
					position++
				}
			l60:
				if !_rules[rulesp]() {
					goto l32
				}
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
					if buffer[position] != rune('S') {
						goto l32
					}
					position++
				}
			l62:
				{
					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l65
					}
					position++
					goto l64
				l65:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
					if buffer[position] != rune('E') {
						goto l32
					}
					position++
				}
			l64:
				{
					position66, tokenIndex66, depth66 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex, depth = position66, tokenIndex66, depth66
					if buffer[position] != rune('L') {
						goto l32
					}
					position++
				}
			l66:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
					if buffer[position] != rune('E') {
						goto l32
					}
					position++
				}
			l68:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l71
					}
					position++
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if buffer[position] != rune('C') {
						goto l32
					}
					position++
				}
			l70:
				{
					position72, tokenIndex72, depth72 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l73
					}
					position++
					goto l72
				l73:
					position, tokenIndex, depth = position72, tokenIndex72, depth72
					if buffer[position] != rune('T') {
						goto l32
					}
					position++
				}
			l72:
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleEmitter]() {
					goto l32
				}
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleProjections]() {
					goto l32
				}
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleWindowedFrom]() {
					goto l32
				}
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleFilter]() {
					goto l32
				}
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleGrouping]() {
					goto l32
				}
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleHaving]() {
					goto l32
				}
				if !_rules[rulesp]() {
					goto l32
				}
				if !_rules[ruleAction1]() {
					goto l32
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position33)
			}
			return true
		l32:
			position, tokenIndex, depth = position32, tokenIndex32, depth32
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position74, tokenIndex74, depth74 := position, tokenIndex, depth
			{
				position75 := position
				depth++
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('C') {
						goto l74
					}
					position++
				}
			l76:
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l79
					}
					position++
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('R') {
						goto l74
					}
					position++
				}
			l78:
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('E') {
						goto l74
					}
					position++
				}
			l80:
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l83
					}
					position++
					goto l82
				l83:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
					if buffer[position] != rune('A') {
						goto l74
					}
					position++
				}
			l82:
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
					if buffer[position] != rune('T') {
						goto l74
					}
					position++
				}
			l84:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if buffer[position] != rune('E') {
						goto l74
					}
					position++
				}
			l86:
				if !_rules[rulesp]() {
					goto l74
				}
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l89
					}
					position++
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if buffer[position] != rune('S') {
						goto l74
					}
					position++
				}
			l88:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('O') {
						goto l74
					}
					position++
				}
			l90:
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('U') {
						goto l74
					}
					position++
				}
			l92:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('R') {
						goto l74
					}
					position++
				}
			l94:
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('C') {
						goto l74
					}
					position++
				}
			l96:
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('E') {
						goto l74
					}
					position++
				}
			l98:
				if !_rules[rulesp]() {
					goto l74
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l74
				}
				if !_rules[rulesp]() {
					goto l74
				}
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('T') {
						goto l74
					}
					position++
				}
			l100:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('Y') {
						goto l74
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('P') {
						goto l74
					}
					position++
				}
			l104:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('E') {
						goto l74
					}
					position++
				}
			l106:
				if !_rules[rulesp]() {
					goto l74
				}
				if !_rules[ruleSourceSinkType]() {
					goto l74
				}
				if !_rules[rulesp]() {
					goto l74
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l74
				}
				if !_rules[ruleAction2]() {
					goto l74
				}
				depth--
				add(ruleCreateSourceStmt, position75)
			}
			return true
		l74:
			position, tokenIndex, depth = position74, tokenIndex74, depth74
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position108, tokenIndex108, depth108 := position, tokenIndex, depth
			{
				position109 := position
				depth++
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('C') {
						goto l108
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('R') {
						goto l108
					}
					position++
				}
			l112:
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('E') {
						goto l108
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('A') {
						goto l108
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('T') {
						goto l108
					}
					position++
				}
			l118:
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
					if buffer[position] != rune('E') {
						goto l108
					}
					position++
				}
			l120:
				if !_rules[rulesp]() {
					goto l108
				}
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('S') {
						goto l108
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('I') {
						goto l108
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('N') {
						goto l108
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('K') {
						goto l108
					}
					position++
				}
			l128:
				if !_rules[rulesp]() {
					goto l108
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l108
				}
				if !_rules[rulesp]() {
					goto l108
				}
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('T') {
						goto l108
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('Y') {
						goto l108
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('P') {
						goto l108
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('E') {
						goto l108
					}
					position++
				}
			l136:
				if !_rules[rulesp]() {
					goto l108
				}
				if !_rules[ruleSourceSinkType]() {
					goto l108
				}
				if !_rules[rulesp]() {
					goto l108
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l108
				}
				if !_rules[ruleAction3]() {
					goto l108
				}
				depth--
				add(ruleCreateSinkStmt, position109)
			}
			return true
		l108:
			position, tokenIndex, depth = position108, tokenIndex108, depth108
			return false
		},
		/* 6 CreateStreamFromSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action4)> */
		func() bool {
			position138, tokenIndex138, depth138 := position, tokenIndex, depth
			{
				position139 := position
				depth++
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('C') {
						goto l138
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('R') {
						goto l138
					}
					position++
				}
			l142:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('E') {
						goto l138
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('A') {
						goto l138
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('T') {
						goto l138
					}
					position++
				}
			l148:
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('E') {
						goto l138
					}
					position++
				}
			l150:
				if !_rules[rulesp]() {
					goto l138
				}
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('S') {
						goto l138
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('T') {
						goto l138
					}
					position++
				}
			l154:
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('R') {
						goto l138
					}
					position++
				}
			l156:
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if buffer[position] != rune('E') {
						goto l138
					}
					position++
				}
			l158:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('A') {
						goto l138
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('M') {
						goto l138
					}
					position++
				}
			l162:
				if !_rules[rulesp]() {
					goto l138
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l138
				}
				if !_rules[rulesp]() {
					goto l138
				}
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('F') {
						goto l138
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('R') {
						goto l138
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('O') {
						goto l138
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('M') {
						goto l138
					}
					position++
				}
			l170:
				if !_rules[rulesp]() {
					goto l138
				}
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('S') {
						goto l138
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('O') {
						goto l138
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('U') {
						goto l138
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('R') {
						goto l138
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('C') {
						goto l138
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('E') {
						goto l138
					}
					position++
				}
			l182:
				if !_rules[rulesp]() {
					goto l138
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l138
				}
				if !_rules[ruleAction4]() {
					goto l138
				}
				depth--
				add(ruleCreateStreamFromSourceStmt, position139)
			}
			return true
		l138:
			position, tokenIndex, depth = position138, tokenIndex138, depth138
			return false
		},
		/* 7 CreateStreamFromSourceExtStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp SourceSinkType sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp SourceSinkSpecs Action5)> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('C') {
						goto l184
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('R') {
						goto l184
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('E') {
						goto l184
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('A') {
						goto l184
					}
					position++
				}
			l192:
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('T') {
						goto l184
					}
					position++
				}
			l194:
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('E') {
						goto l184
					}
					position++
				}
			l196:
				if !_rules[rulesp]() {
					goto l184
				}
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('S') {
						goto l184
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('T') {
						goto l184
					}
					position++
				}
			l200:
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('R') {
						goto l184
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('E') {
						goto l184
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('A') {
						goto l184
					}
					position++
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('M') {
						goto l184
					}
					position++
				}
			l208:
				if !_rules[rulesp]() {
					goto l184
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l184
				}
				if !_rules[rulesp]() {
					goto l184
				}
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('F') {
						goto l184
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('R') {
						goto l184
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('O') {
						goto l184
					}
					position++
				}
			l214:
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('M') {
						goto l184
					}
					position++
				}
			l216:
				if !_rules[rulesp]() {
					goto l184
				}
				if !_rules[ruleSourceSinkType]() {
					goto l184
				}
				if !_rules[rulesp]() {
					goto l184
				}
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('S') {
						goto l184
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('O') {
						goto l184
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('U') {
						goto l184
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('R') {
						goto l184
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('C') {
						goto l184
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('E') {
						goto l184
					}
					position++
				}
			l228:
				if !_rules[rulesp]() {
					goto l184
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l184
				}
				if !_rules[ruleAction5]() {
					goto l184
				}
				depth--
				add(ruleCreateStreamFromSourceExtStmt, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
		/* 8 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action6)> */
		func() bool {
			position230, tokenIndex230, depth230 := position, tokenIndex, depth
			{
				position231 := position
				depth++
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('I') {
						goto l230
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('N') {
						goto l230
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('S') {
						goto l230
					}
					position++
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('E') {
						goto l230
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('R') {
						goto l230
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('T') {
						goto l230
					}
					position++
				}
			l242:
				if !_rules[rulesp]() {
					goto l230
				}
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('I') {
						goto l230
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('N') {
						goto l230
					}
					position++
				}
			l246:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('T') {
						goto l230
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('O') {
						goto l230
					}
					position++
				}
			l250:
				if !_rules[rulesp]() {
					goto l230
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l230
				}
				if !_rules[rulesp]() {
					goto l230
				}
				if !_rules[ruleSelectStmt]() {
					goto l230
				}
				if !_rules[ruleAction6]() {
					goto l230
				}
				depth--
				add(ruleInsertIntoSelectStmt, position231)
			}
			return true
		l230:
			position, tokenIndex, depth = position230, tokenIndex230, depth230
			return false
		},
		/* 9 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action7)> */
		func() bool {
			position252, tokenIndex252, depth252 := position, tokenIndex, depth
			{
				position253 := position
				depth++
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l255
					}
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if !_rules[ruleDSTREAM]() {
						goto l256
					}
					goto l254
				l256:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if !_rules[ruleRSTREAM]() {
						goto l252
					}
				}
			l254:
				{
					position257 := position
					depth++
					{
						position258, tokenIndex258, depth258 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l258
						}
						if buffer[position] != rune('[') {
							goto l258
						}
						position++
						if !_rules[rulesp]() {
							goto l258
						}
						{
							position260, tokenIndex260, depth260 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l261
							}
							position++
							goto l260
						l261:
							position, tokenIndex, depth = position260, tokenIndex260, depth260
							if buffer[position] != rune('E') {
								goto l258
							}
							position++
						}
					l260:
						{
							position262, tokenIndex262, depth262 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l263
							}
							position++
							goto l262
						l263:
							position, tokenIndex, depth = position262, tokenIndex262, depth262
							if buffer[position] != rune('V') {
								goto l258
							}
							position++
						}
					l262:
						{
							position264, tokenIndex264, depth264 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l265
							}
							position++
							goto l264
						l265:
							position, tokenIndex, depth = position264, tokenIndex264, depth264
							if buffer[position] != rune('E') {
								goto l258
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
								goto l258
							}
							position++
						}
					l266:
						{
							position268, tokenIndex268, depth268 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l269
							}
							position++
							goto l268
						l269:
							position, tokenIndex, depth = position268, tokenIndex268, depth268
							if buffer[position] != rune('Y') {
								goto l258
							}
							position++
						}
					l268:
						if !_rules[rulesp]() {
							goto l258
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l258
						}
						if !_rules[rulesp]() {
							goto l258
						}
						if buffer[position] != rune(']') {
							goto l258
						}
						position++
						goto l259
					l258:
						position, tokenIndex, depth = position258, tokenIndex258, depth258
					}
				l259:
					depth--
					add(rulePegText, position257)
				}
				if !_rules[ruleAction7]() {
					goto l252
				}
				depth--
				add(ruleEmitter, position253)
			}
			return true
		l252:
			position, tokenIndex, depth = position252, tokenIndex252, depth252
			return false
		},
		/* 10 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position270, tokenIndex270, depth270 := position, tokenIndex, depth
			{
				position271 := position
				depth++
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l273
					}
				l274:
					{
						position275, tokenIndex275, depth275 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l275
						}
						if buffer[position] != rune(',') {
							goto l275
						}
						position++
						if !_rules[rulesp]() {
							goto l275
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l275
						}
						goto l274
					l275:
						position, tokenIndex, depth = position275, tokenIndex275, depth275
					}
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if !_rules[ruleTimeEmitterInterval]() {
						goto l276
					}
					goto l272
				l276:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if !_rules[ruleTupleEmitterInterval]() {
						goto l270
					}
				}
			l272:
				depth--
				add(ruleEmitterIntervals, position271)
			}
			return true
		l270:
			position, tokenIndex, depth = position270, tokenIndex270, depth270
			return false
		},
		/* 11 TimeEmitterInterval <- <(<TimeInterval> Action8)> */
		func() bool {
			position277, tokenIndex277, depth277 := position, tokenIndex, depth
			{
				position278 := position
				depth++
				{
					position279 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l277
					}
					depth--
					add(rulePegText, position279)
				}
				if !_rules[ruleAction8]() {
					goto l277
				}
				depth--
				add(ruleTimeEmitterInterval, position278)
			}
			return true
		l277:
			position, tokenIndex, depth = position277, tokenIndex277, depth277
			return false
		},
		/* 12 TupleEmitterInterval <- <(<TuplesInterval> Action9)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				{
					position282 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l280
					}
					depth--
					add(rulePegText, position282)
				}
				if !_rules[ruleAction9]() {
					goto l280
				}
				depth--
				add(ruleTupleEmitterInterval, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 13 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action10)> */
		func() bool {
			position283, tokenIndex283, depth283 := position, tokenIndex, depth
			{
				position284 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l283
				}
				if !_rules[rulesp]() {
					goto l283
				}
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('I') {
						goto l283
					}
					position++
				}
			l285:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('N') {
						goto l283
					}
					position++
				}
			l287:
				if !_rules[rulesp]() {
					goto l283
				}
				if !_rules[ruleStream]() {
					goto l283
				}
				if !_rules[ruleAction10]() {
					goto l283
				}
				depth--
				add(ruleTupleEmitterFromInterval, position284)
			}
			return true
		l283:
			position, tokenIndex, depth = position283, tokenIndex283, depth283
			return false
		},
		/* 14 Projections <- <(<(Projection sp (',' sp Projection)*)> Action11)> */
		func() bool {
			position289, tokenIndex289, depth289 := position, tokenIndex, depth
			{
				position290 := position
				depth++
				{
					position291 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l289
					}
					if !_rules[rulesp]() {
						goto l289
					}
				l292:
					{
						position293, tokenIndex293, depth293 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l293
						}
						position++
						if !_rules[rulesp]() {
							goto l293
						}
						if !_rules[ruleProjection]() {
							goto l293
						}
						goto l292
					l293:
						position, tokenIndex, depth = position293, tokenIndex293, depth293
					}
					depth--
					add(rulePegText, position291)
				}
				if !_rules[ruleAction11]() {
					goto l289
				}
				depth--
				add(ruleProjections, position290)
			}
			return true
		l289:
			position, tokenIndex, depth = position289, tokenIndex289, depth289
			return false
		},
		/* 15 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position294, tokenIndex294, depth294 := position, tokenIndex, depth
			{
				position295 := position
				depth++
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l297
					}
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if !_rules[ruleExpression]() {
						goto l298
					}
					goto l296
				l298:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if !_rules[ruleWildcard]() {
						goto l294
					}
				}
			l296:
				depth--
				add(ruleProjection, position295)
			}
			return true
		l294:
			position, tokenIndex, depth = position294, tokenIndex294, depth294
			return false
		},
		/* 16 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action12)> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{
				position300 := position
				depth++
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l302
					}
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if !_rules[ruleWildcard]() {
						goto l299
					}
				}
			l301:
				if !_rules[rulesp]() {
					goto l299
				}
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('A') {
						goto l299
					}
					position++
				}
			l303:
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l306
					}
					position++
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if buffer[position] != rune('S') {
						goto l299
					}
					position++
				}
			l305:
				if !_rules[rulesp]() {
					goto l299
				}
				if !_rules[ruleIdentifier]() {
					goto l299
				}
				if !_rules[ruleAction12]() {
					goto l299
				}
				depth--
				add(ruleAliasExpression, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 17 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action13)> */
		func() bool {
			position307, tokenIndex307, depth307 := position, tokenIndex, depth
			{
				position308 := position
				depth++
				{
					position309 := position
					depth++
					{
						position310, tokenIndex310, depth310 := position, tokenIndex, depth
						{
							position312, tokenIndex312, depth312 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l313
							}
							position++
							goto l312
						l313:
							position, tokenIndex, depth = position312, tokenIndex312, depth312
							if buffer[position] != rune('F') {
								goto l310
							}
							position++
						}
					l312:
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('R') {
								goto l310
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('O') {
								goto l310
							}
							position++
						}
					l316:
						{
							position318, tokenIndex318, depth318 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l319
							}
							position++
							goto l318
						l319:
							position, tokenIndex, depth = position318, tokenIndex318, depth318
							if buffer[position] != rune('M') {
								goto l310
							}
							position++
						}
					l318:
						if !_rules[rulesp]() {
							goto l310
						}
						if !_rules[ruleRelations]() {
							goto l310
						}
						if !_rules[rulesp]() {
							goto l310
						}
						goto l311
					l310:
						position, tokenIndex, depth = position310, tokenIndex310, depth310
					}
				l311:
					depth--
					add(rulePegText, position309)
				}
				if !_rules[ruleAction13]() {
					goto l307
				}
				depth--
				add(ruleWindowedFrom, position308)
			}
			return true
		l307:
			position, tokenIndex, depth = position307, tokenIndex307, depth307
			return false
		},
		/* 18 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action14)> */
		func() bool {
			position320, tokenIndex320, depth320 := position, tokenIndex, depth
			{
				position321 := position
				depth++
				{
					position322 := position
					depth++
					{
						position323, tokenIndex323, depth323 := position, tokenIndex, depth
						{
							position325, tokenIndex325, depth325 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l326
							}
							position++
							goto l325
						l326:
							position, tokenIndex, depth = position325, tokenIndex325, depth325
							if buffer[position] != rune('F') {
								goto l323
							}
							position++
						}
					l325:
						{
							position327, tokenIndex327, depth327 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l328
							}
							position++
							goto l327
						l328:
							position, tokenIndex, depth = position327, tokenIndex327, depth327
							if buffer[position] != rune('R') {
								goto l323
							}
							position++
						}
					l327:
						{
							position329, tokenIndex329, depth329 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l330
							}
							position++
							goto l329
						l330:
							position, tokenIndex, depth = position329, tokenIndex329, depth329
							if buffer[position] != rune('O') {
								goto l323
							}
							position++
						}
					l329:
						{
							position331, tokenIndex331, depth331 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l332
							}
							position++
							goto l331
						l332:
							position, tokenIndex, depth = position331, tokenIndex331, depth331
							if buffer[position] != rune('M') {
								goto l323
							}
							position++
						}
					l331:
						if !_rules[rulesp]() {
							goto l323
						}
						if !_rules[ruleDefRelations]() {
							goto l323
						}
						if !_rules[rulesp]() {
							goto l323
						}
						goto l324
					l323:
						position, tokenIndex, depth = position323, tokenIndex323, depth323
					}
				l324:
					depth--
					add(rulePegText, position322)
				}
				if !_rules[ruleAction14]() {
					goto l320
				}
				depth--
				add(ruleDefWindowedFrom, position321)
			}
			return true
		l320:
			position, tokenIndex, depth = position320, tokenIndex320, depth320
			return false
		},
		/* 19 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position333, tokenIndex333, depth333 := position, tokenIndex, depth
			{
				position334 := position
				depth++
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l336
					}
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if !_rules[ruleTuplesInterval]() {
						goto l333
					}
				}
			l335:
				depth--
				add(ruleInterval, position334)
			}
			return true
		l333:
			position, tokenIndex, depth = position333, tokenIndex333, depth333
			return false
		},
		/* 20 TimeInterval <- <(NumericLiteral sp SECONDS Action15)> */
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
				if !_rules[ruleSECONDS]() {
					goto l337
				}
				if !_rules[ruleAction15]() {
					goto l337
				}
				depth--
				add(ruleTimeInterval, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 21 TuplesInterval <- <(NumericLiteral sp TUPLES Action16)> */
		func() bool {
			position339, tokenIndex339, depth339 := position, tokenIndex, depth
			{
				position340 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l339
				}
				if !_rules[rulesp]() {
					goto l339
				}
				if !_rules[ruleTUPLES]() {
					goto l339
				}
				if !_rules[ruleAction16]() {
					goto l339
				}
				depth--
				add(ruleTuplesInterval, position340)
			}
			return true
		l339:
			position, tokenIndex, depth = position339, tokenIndex339, depth339
			return false
		},
		/* 22 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
			{
				position342 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l341
				}
				if !_rules[rulesp]() {
					goto l341
				}
			l343:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l344
					}
					position++
					if !_rules[rulesp]() {
						goto l344
					}
					if !_rules[ruleRelationLike]() {
						goto l344
					}
					goto l343
				l344:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
				}
				depth--
				add(ruleRelations, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
			return false
		},
		/* 23 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position345, tokenIndex345, depth345 := position, tokenIndex, depth
			{
				position346 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l345
				}
				if !_rules[rulesp]() {
					goto l345
				}
			l347:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l348
					}
					position++
					if !_rules[rulesp]() {
						goto l348
					}
					if !_rules[ruleDefRelationLike]() {
						goto l348
					}
					goto l347
				l348:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
				}
				depth--
				add(ruleDefRelations, position346)
			}
			return true
		l345:
			position, tokenIndex, depth = position345, tokenIndex345, depth345
			return false
		},
		/* 24 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action17)> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				{
					position351 := position
					depth++
					{
						position352, tokenIndex352, depth352 := position, tokenIndex, depth
						{
							position354, tokenIndex354, depth354 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l355
							}
							position++
							goto l354
						l355:
							position, tokenIndex, depth = position354, tokenIndex354, depth354
							if buffer[position] != rune('W') {
								goto l352
							}
							position++
						}
					l354:
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l357
							}
							position++
							goto l356
						l357:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
							if buffer[position] != rune('H') {
								goto l352
							}
							position++
						}
					l356:
						{
							position358, tokenIndex358, depth358 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex, depth = position358, tokenIndex358, depth358
							if buffer[position] != rune('E') {
								goto l352
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
								goto l352
							}
							position++
						}
					l360:
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('E') {
								goto l352
							}
							position++
						}
					l362:
						if !_rules[rulesp]() {
							goto l352
						}
						if !_rules[ruleExpression]() {
							goto l352
						}
						goto l353
					l352:
						position, tokenIndex, depth = position352, tokenIndex352, depth352
					}
				l353:
					depth--
					add(rulePegText, position351)
				}
				if !_rules[ruleAction17]() {
					goto l349
				}
				depth--
				add(ruleFilter, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 25 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action18)> */
		func() bool {
			position364, tokenIndex364, depth364 := position, tokenIndex, depth
			{
				position365 := position
				depth++
				{
					position366 := position
					depth++
					{
						position367, tokenIndex367, depth367 := position, tokenIndex, depth
						{
							position369, tokenIndex369, depth369 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l370
							}
							position++
							goto l369
						l370:
							position, tokenIndex, depth = position369, tokenIndex369, depth369
							if buffer[position] != rune('G') {
								goto l367
							}
							position++
						}
					l369:
						{
							position371, tokenIndex371, depth371 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l372
							}
							position++
							goto l371
						l372:
							position, tokenIndex, depth = position371, tokenIndex371, depth371
							if buffer[position] != rune('R') {
								goto l367
							}
							position++
						}
					l371:
						{
							position373, tokenIndex373, depth373 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l374
							}
							position++
							goto l373
						l374:
							position, tokenIndex, depth = position373, tokenIndex373, depth373
							if buffer[position] != rune('O') {
								goto l367
							}
							position++
						}
					l373:
						{
							position375, tokenIndex375, depth375 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l376
							}
							position++
							goto l375
						l376:
							position, tokenIndex, depth = position375, tokenIndex375, depth375
							if buffer[position] != rune('U') {
								goto l367
							}
							position++
						}
					l375:
						{
							position377, tokenIndex377, depth377 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l378
							}
							position++
							goto l377
						l378:
							position, tokenIndex, depth = position377, tokenIndex377, depth377
							if buffer[position] != rune('P') {
								goto l367
							}
							position++
						}
					l377:
						if !_rules[rulesp]() {
							goto l367
						}
						{
							position379, tokenIndex379, depth379 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l380
							}
							position++
							goto l379
						l380:
							position, tokenIndex, depth = position379, tokenIndex379, depth379
							if buffer[position] != rune('B') {
								goto l367
							}
							position++
						}
					l379:
						{
							position381, tokenIndex381, depth381 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l382
							}
							position++
							goto l381
						l382:
							position, tokenIndex, depth = position381, tokenIndex381, depth381
							if buffer[position] != rune('Y') {
								goto l367
							}
							position++
						}
					l381:
						if !_rules[rulesp]() {
							goto l367
						}
						if !_rules[ruleGroupList]() {
							goto l367
						}
						goto l368
					l367:
						position, tokenIndex, depth = position367, tokenIndex367, depth367
					}
				l368:
					depth--
					add(rulePegText, position366)
				}
				if !_rules[ruleAction18]() {
					goto l364
				}
				depth--
				add(ruleGrouping, position365)
			}
			return true
		l364:
			position, tokenIndex, depth = position364, tokenIndex364, depth364
			return false
		},
		/* 26 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position383, tokenIndex383, depth383 := position, tokenIndex, depth
			{
				position384 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l383
				}
				if !_rules[rulesp]() {
					goto l383
				}
			l385:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l386
					}
					position++
					if !_rules[rulesp]() {
						goto l386
					}
					if !_rules[ruleExpression]() {
						goto l386
					}
					goto l385
				l386:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
				}
				depth--
				add(ruleGroupList, position384)
			}
			return true
		l383:
			position, tokenIndex, depth = position383, tokenIndex383, depth383
			return false
		},
		/* 27 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action19)> */
		func() bool {
			position387, tokenIndex387, depth387 := position, tokenIndex, depth
			{
				position388 := position
				depth++
				{
					position389 := position
					depth++
					{
						position390, tokenIndex390, depth390 := position, tokenIndex, depth
						{
							position392, tokenIndex392, depth392 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l393
							}
							position++
							goto l392
						l393:
							position, tokenIndex, depth = position392, tokenIndex392, depth392
							if buffer[position] != rune('H') {
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
							if buffer[position] != rune('v') {
								goto l397
							}
							position++
							goto l396
						l397:
							position, tokenIndex, depth = position396, tokenIndex396, depth396
							if buffer[position] != rune('V') {
								goto l390
							}
							position++
						}
					l396:
						{
							position398, tokenIndex398, depth398 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l399
							}
							position++
							goto l398
						l399:
							position, tokenIndex, depth = position398, tokenIndex398, depth398
							if buffer[position] != rune('I') {
								goto l390
							}
							position++
						}
					l398:
						{
							position400, tokenIndex400, depth400 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l401
							}
							position++
							goto l400
						l401:
							position, tokenIndex, depth = position400, tokenIndex400, depth400
							if buffer[position] != rune('N') {
								goto l390
							}
							position++
						}
					l400:
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
								goto l390
							}
							position++
						}
					l402:
						if !_rules[rulesp]() {
							goto l390
						}
						if !_rules[ruleExpression]() {
							goto l390
						}
						goto l391
					l390:
						position, tokenIndex, depth = position390, tokenIndex390, depth390
					}
				l391:
					depth--
					add(rulePegText, position389)
				}
				if !_rules[ruleAction19]() {
					goto l387
				}
				depth--
				add(ruleHaving, position388)
			}
			return true
		l387:
			position, tokenIndex, depth = position387, tokenIndex387, depth387
			return false
		},
		/* 28 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action20))> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l407
					}
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if !_rules[ruleStreamWindow]() {
						goto l404
					}
					if !_rules[ruleAction20]() {
						goto l404
					}
				}
			l406:
				depth--
				add(ruleRelationLike, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 29 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action21))> */
		func() bool {
			position408, tokenIndex408, depth408 := position, tokenIndex, depth
			{
				position409 := position
				depth++
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l411
					}
					goto l410
				l411:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
					if !_rules[ruleDefStreamWindow]() {
						goto l408
					}
					if !_rules[ruleAction21]() {
						goto l408
					}
				}
			l410:
				depth--
				add(ruleDefRelationLike, position409)
			}
			return true
		l408:
			position, tokenIndex, depth = position408, tokenIndex408, depth408
			return false
		},
		/* 30 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action22)> */
		func() bool {
			position412, tokenIndex412, depth412 := position, tokenIndex, depth
			{
				position413 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l412
				}
				if !_rules[rulesp]() {
					goto l412
				}
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('A') {
						goto l412
					}
					position++
				}
			l414:
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l417
					}
					position++
					goto l416
				l417:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if buffer[position] != rune('S') {
						goto l412
					}
					position++
				}
			l416:
				if !_rules[rulesp]() {
					goto l412
				}
				if !_rules[ruleIdentifier]() {
					goto l412
				}
				if !_rules[ruleAction22]() {
					goto l412
				}
				depth--
				add(ruleAliasedStreamWindow, position413)
			}
			return true
		l412:
			position, tokenIndex, depth = position412, tokenIndex412, depth412
			return false
		},
		/* 31 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action23)> */
		func() bool {
			position418, tokenIndex418, depth418 := position, tokenIndex, depth
			{
				position419 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l418
				}
				if !_rules[rulesp]() {
					goto l418
				}
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('A') {
						goto l418
					}
					position++
				}
			l420:
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l423
					}
					position++
					goto l422
				l423:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
					if buffer[position] != rune('S') {
						goto l418
					}
					position++
				}
			l422:
				if !_rules[rulesp]() {
					goto l418
				}
				if !_rules[ruleIdentifier]() {
					goto l418
				}
				if !_rules[ruleAction23]() {
					goto l418
				}
				depth--
				add(ruleDefAliasedStreamWindow, position419)
			}
			return true
		l418:
			position, tokenIndex, depth = position418, tokenIndex418, depth418
			return false
		},
		/* 32 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action24)> */
		func() bool {
			position424, tokenIndex424, depth424 := position, tokenIndex, depth
			{
				position425 := position
				depth++
				if !_rules[ruleStream]() {
					goto l424
				}
				if !_rules[rulesp]() {
					goto l424
				}
				if buffer[position] != rune('[') {
					goto l424
				}
				position++
				if !_rules[rulesp]() {
					goto l424
				}
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l427
					}
					position++
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if buffer[position] != rune('R') {
						goto l424
					}
					position++
				}
			l426:
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l429
					}
					position++
					goto l428
				l429:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
					if buffer[position] != rune('A') {
						goto l424
					}
					position++
				}
			l428:
				{
					position430, tokenIndex430, depth430 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l431
					}
					position++
					goto l430
				l431:
					position, tokenIndex, depth = position430, tokenIndex430, depth430
					if buffer[position] != rune('N') {
						goto l424
					}
					position++
				}
			l430:
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('G') {
						goto l424
					}
					position++
				}
			l432:
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l435
					}
					position++
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if buffer[position] != rune('E') {
						goto l424
					}
					position++
				}
			l434:
				if !_rules[rulesp]() {
					goto l424
				}
				if !_rules[ruleInterval]() {
					goto l424
				}
				if !_rules[rulesp]() {
					goto l424
				}
				if buffer[position] != rune(']') {
					goto l424
				}
				position++
				if !_rules[ruleAction24]() {
					goto l424
				}
				depth--
				add(ruleStreamWindow, position425)
			}
			return true
		l424:
			position, tokenIndex, depth = position424, tokenIndex424, depth424
			return false
		},
		/* 33 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action25)> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				if !_rules[ruleStream]() {
					goto l436
				}
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l438
					}
					if buffer[position] != rune('[') {
						goto l438
					}
					position++
					if !_rules[rulesp]() {
						goto l438
					}
					{
						position440, tokenIndex440, depth440 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l441
						}
						position++
						goto l440
					l441:
						position, tokenIndex, depth = position440, tokenIndex440, depth440
						if buffer[position] != rune('R') {
							goto l438
						}
						position++
					}
				l440:
					{
						position442, tokenIndex442, depth442 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l443
						}
						position++
						goto l442
					l443:
						position, tokenIndex, depth = position442, tokenIndex442, depth442
						if buffer[position] != rune('A') {
							goto l438
						}
						position++
					}
				l442:
					{
						position444, tokenIndex444, depth444 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l445
						}
						position++
						goto l444
					l445:
						position, tokenIndex, depth = position444, tokenIndex444, depth444
						if buffer[position] != rune('N') {
							goto l438
						}
						position++
					}
				l444:
					{
						position446, tokenIndex446, depth446 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l447
						}
						position++
						goto l446
					l447:
						position, tokenIndex, depth = position446, tokenIndex446, depth446
						if buffer[position] != rune('G') {
							goto l438
						}
						position++
					}
				l446:
					{
						position448, tokenIndex448, depth448 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l449
						}
						position++
						goto l448
					l449:
						position, tokenIndex, depth = position448, tokenIndex448, depth448
						if buffer[position] != rune('E') {
							goto l438
						}
						position++
					}
				l448:
					if !_rules[rulesp]() {
						goto l438
					}
					if !_rules[ruleInterval]() {
						goto l438
					}
					if !_rules[rulesp]() {
						goto l438
					}
					if buffer[position] != rune(']') {
						goto l438
					}
					position++
					goto l439
				l438:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
				}
			l439:
				if !_rules[ruleAction25]() {
					goto l436
				}
				depth--
				add(ruleDefStreamWindow, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
			return false
		},
		/* 34 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action26)> */
		func() bool {
			position450, tokenIndex450, depth450 := position, tokenIndex, depth
			{
				position451 := position
				depth++
				{
					position452 := position
					depth++
					{
						position453, tokenIndex453, depth453 := position, tokenIndex, depth
						{
							position455, tokenIndex455, depth455 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l456
							}
							position++
							goto l455
						l456:
							position, tokenIndex, depth = position455, tokenIndex455, depth455
							if buffer[position] != rune('W') {
								goto l453
							}
							position++
						}
					l455:
						{
							position457, tokenIndex457, depth457 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l458
							}
							position++
							goto l457
						l458:
							position, tokenIndex, depth = position457, tokenIndex457, depth457
							if buffer[position] != rune('I') {
								goto l453
							}
							position++
						}
					l457:
						{
							position459, tokenIndex459, depth459 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l460
							}
							position++
							goto l459
						l460:
							position, tokenIndex, depth = position459, tokenIndex459, depth459
							if buffer[position] != rune('T') {
								goto l453
							}
							position++
						}
					l459:
						{
							position461, tokenIndex461, depth461 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l462
							}
							position++
							goto l461
						l462:
							position, tokenIndex, depth = position461, tokenIndex461, depth461
							if buffer[position] != rune('H') {
								goto l453
							}
							position++
						}
					l461:
						if !_rules[rulesp]() {
							goto l453
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l453
						}
						if !_rules[rulesp]() {
							goto l453
						}
					l463:
						{
							position464, tokenIndex464, depth464 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l464
							}
							position++
							if !_rules[rulesp]() {
								goto l464
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l464
							}
							goto l463
						l464:
							position, tokenIndex, depth = position464, tokenIndex464, depth464
						}
						goto l454
					l453:
						position, tokenIndex, depth = position453, tokenIndex453, depth453
					}
				l454:
					depth--
					add(rulePegText, position452)
				}
				if !_rules[ruleAction26]() {
					goto l450
				}
				depth--
				add(ruleSourceSinkSpecs, position451)
			}
			return true
		l450:
			position, tokenIndex, depth = position450, tokenIndex450, depth450
			return false
		},
		/* 35 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action27)> */
		func() bool {
			position465, tokenIndex465, depth465 := position, tokenIndex, depth
			{
				position466 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l465
				}
				if buffer[position] != rune('=') {
					goto l465
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l465
				}
				if !_rules[ruleAction27]() {
					goto l465
				}
				depth--
				add(ruleSourceSinkParam, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 36 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l470
					}
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if !_rules[ruleLiteral]() {
						goto l467
					}
				}
			l469:
				depth--
				add(ruleSourceSinkParamVal, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 37 Expression <- <orExpr> */
		func() bool {
			position471, tokenIndex471, depth471 := position, tokenIndex, depth
			{
				position472 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l471
				}
				depth--
				add(ruleExpression, position472)
			}
			return true
		l471:
			position, tokenIndex, depth = position471, tokenIndex471, depth471
			return false
		},
		/* 38 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action28)> */
		func() bool {
			position473, tokenIndex473, depth473 := position, tokenIndex, depth
			{
				position474 := position
				depth++
				{
					position475 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l473
					}
					if !_rules[rulesp]() {
						goto l473
					}
					{
						position476, tokenIndex476, depth476 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l476
						}
						if !_rules[rulesp]() {
							goto l476
						}
						if !_rules[ruleandExpr]() {
							goto l476
						}
						goto l477
					l476:
						position, tokenIndex, depth = position476, tokenIndex476, depth476
					}
				l477:
					depth--
					add(rulePegText, position475)
				}
				if !_rules[ruleAction28]() {
					goto l473
				}
				depth--
				add(ruleorExpr, position474)
			}
			return true
		l473:
			position, tokenIndex, depth = position473, tokenIndex473, depth473
			return false
		},
		/* 39 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action29)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				{
					position480 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l478
					}
					if !_rules[rulesp]() {
						goto l478
					}
					{
						position481, tokenIndex481, depth481 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l481
						}
						if !_rules[rulesp]() {
							goto l481
						}
						if !_rules[rulecomparisonExpr]() {
							goto l481
						}
						goto l482
					l481:
						position, tokenIndex, depth = position481, tokenIndex481, depth481
					}
				l482:
					depth--
					add(rulePegText, position480)
				}
				if !_rules[ruleAction29]() {
					goto l478
				}
				depth--
				add(ruleandExpr, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 40 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action30)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				{
					position485 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l483
					}
					if !_rules[rulesp]() {
						goto l483
					}
					{
						position486, tokenIndex486, depth486 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l486
						}
						if !_rules[rulesp]() {
							goto l486
						}
						if !_rules[ruletermExpr]() {
							goto l486
						}
						goto l487
					l486:
						position, tokenIndex, depth = position486, tokenIndex486, depth486
					}
				l487:
					depth--
					add(rulePegText, position485)
				}
				if !_rules[ruleAction30]() {
					goto l483
				}
				depth--
				add(rulecomparisonExpr, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 41 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action31)> */
		func() bool {
			position488, tokenIndex488, depth488 := position, tokenIndex, depth
			{
				position489 := position
				depth++
				{
					position490 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l488
					}
					if !_rules[rulesp]() {
						goto l488
					}
					{
						position491, tokenIndex491, depth491 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l491
						}
						if !_rules[rulesp]() {
							goto l491
						}
						if !_rules[ruleproductExpr]() {
							goto l491
						}
						goto l492
					l491:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
					}
				l492:
					depth--
					add(rulePegText, position490)
				}
				if !_rules[ruleAction31]() {
					goto l488
				}
				depth--
				add(ruletermExpr, position489)
			}
			return true
		l488:
			position, tokenIndex, depth = position488, tokenIndex488, depth488
			return false
		},
		/* 42 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action32)> */
		func() bool {
			position493, tokenIndex493, depth493 := position, tokenIndex, depth
			{
				position494 := position
				depth++
				{
					position495 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l493
					}
					if !_rules[rulesp]() {
						goto l493
					}
					{
						position496, tokenIndex496, depth496 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l496
						}
						if !_rules[rulesp]() {
							goto l496
						}
						if !_rules[rulebaseExpr]() {
							goto l496
						}
						goto l497
					l496:
						position, tokenIndex, depth = position496, tokenIndex496, depth496
					}
				l497:
					depth--
					add(rulePegText, position495)
				}
				if !_rules[ruleAction32]() {
					goto l493
				}
				depth--
				add(ruleproductExpr, position494)
			}
			return true
		l493:
			position, tokenIndex, depth = position493, tokenIndex493, depth493
			return false
		},
		/* 43 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position498, tokenIndex498, depth498 := position, tokenIndex, depth
			{
				position499 := position
				depth++
				{
					position500, tokenIndex500, depth500 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l501
					}
					position++
					if !_rules[rulesp]() {
						goto l501
					}
					if !_rules[ruleExpression]() {
						goto l501
					}
					if !_rules[rulesp]() {
						goto l501
					}
					if buffer[position] != rune(')') {
						goto l501
					}
					position++
					goto l500
				l501:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if !_rules[ruleBooleanLiteral]() {
						goto l502
					}
					goto l500
				l502:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if !_rules[ruleFuncApp]() {
						goto l503
					}
					goto l500
				l503:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if !_rules[ruleRowValue]() {
						goto l504
					}
					goto l500
				l504:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if !_rules[ruleLiteral]() {
						goto l498
					}
				}
			l500:
				depth--
				add(rulebaseExpr, position499)
			}
			return true
		l498:
			position, tokenIndex, depth = position498, tokenIndex498, depth498
			return false
		},
		/* 44 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action33)> */
		func() bool {
			position505, tokenIndex505, depth505 := position, tokenIndex, depth
			{
				position506 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l505
				}
				if !_rules[rulesp]() {
					goto l505
				}
				if buffer[position] != rune('(') {
					goto l505
				}
				position++
				if !_rules[rulesp]() {
					goto l505
				}
				if !_rules[ruleFuncParams]() {
					goto l505
				}
				if !_rules[rulesp]() {
					goto l505
				}
				if buffer[position] != rune(')') {
					goto l505
				}
				position++
				if !_rules[ruleAction33]() {
					goto l505
				}
				depth--
				add(ruleFuncApp, position506)
			}
			return true
		l505:
			position, tokenIndex, depth = position505, tokenIndex505, depth505
			return false
		},
		/* 45 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action34)> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				{
					position509 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l507
					}
					if !_rules[rulesp]() {
						goto l507
					}
				l510:
					{
						position511, tokenIndex511, depth511 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l511
						}
						position++
						if !_rules[rulesp]() {
							goto l511
						}
						if !_rules[ruleExpression]() {
							goto l511
						}
						goto l510
					l511:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
					}
					depth--
					add(rulePegText, position509)
				}
				if !_rules[ruleAction34]() {
					goto l507
				}
				depth--
				add(ruleFuncParams, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 46 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position512, tokenIndex512, depth512 := position, tokenIndex, depth
			{
				position513 := position
				depth++
				{
					position514, tokenIndex514, depth514 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l515
					}
					goto l514
				l515:
					position, tokenIndex, depth = position514, tokenIndex514, depth514
					if !_rules[ruleNumericLiteral]() {
						goto l516
					}
					goto l514
				l516:
					position, tokenIndex, depth = position514, tokenIndex514, depth514
					if !_rules[ruleStringLiteral]() {
						goto l512
					}
				}
			l514:
				depth--
				add(ruleLiteral, position513)
			}
			return true
		l512:
			position, tokenIndex, depth = position512, tokenIndex512, depth512
			return false
		},
		/* 47 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position517, tokenIndex517, depth517 := position, tokenIndex, depth
			{
				position518 := position
				depth++
				{
					position519, tokenIndex519, depth519 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l520
					}
					goto l519
				l520:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if !_rules[ruleNotEqual]() {
						goto l521
					}
					goto l519
				l521:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if !_rules[ruleLessOrEqual]() {
						goto l522
					}
					goto l519
				l522:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if !_rules[ruleLess]() {
						goto l523
					}
					goto l519
				l523:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if !_rules[ruleGreaterOrEqual]() {
						goto l524
					}
					goto l519
				l524:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if !_rules[ruleGreater]() {
						goto l525
					}
					goto l519
				l525:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if !_rules[ruleNotEqual]() {
						goto l517
					}
				}
			l519:
				depth--
				add(ruleComparisonOp, position518)
			}
			return true
		l517:
			position, tokenIndex, depth = position517, tokenIndex517, depth517
			return false
		},
		/* 48 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position526, tokenIndex526, depth526 := position, tokenIndex, depth
			{
				position527 := position
				depth++
				{
					position528, tokenIndex528, depth528 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l529
					}
					goto l528
				l529:
					position, tokenIndex, depth = position528, tokenIndex528, depth528
					if !_rules[ruleMinus]() {
						goto l526
					}
				}
			l528:
				depth--
				add(rulePlusMinusOp, position527)
			}
			return true
		l526:
			position, tokenIndex, depth = position526, tokenIndex526, depth526
			return false
		},
		/* 49 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position530, tokenIndex530, depth530 := position, tokenIndex, depth
			{
				position531 := position
				depth++
				{
					position532, tokenIndex532, depth532 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l533
					}
					goto l532
				l533:
					position, tokenIndex, depth = position532, tokenIndex532, depth532
					if !_rules[ruleDivide]() {
						goto l534
					}
					goto l532
				l534:
					position, tokenIndex, depth = position532, tokenIndex532, depth532
					if !_rules[ruleModulo]() {
						goto l530
					}
				}
			l532:
				depth--
				add(ruleMultDivOp, position531)
			}
			return true
		l530:
			position, tokenIndex, depth = position530, tokenIndex530, depth530
			return false
		},
		/* 50 Stream <- <(<ident> Action35)> */
		func() bool {
			position535, tokenIndex535, depth535 := position, tokenIndex, depth
			{
				position536 := position
				depth++
				{
					position537 := position
					depth++
					if !_rules[ruleident]() {
						goto l535
					}
					depth--
					add(rulePegText, position537)
				}
				if !_rules[ruleAction35]() {
					goto l535
				}
				depth--
				add(ruleStream, position536)
			}
			return true
		l535:
			position, tokenIndex, depth = position535, tokenIndex535, depth535
			return false
		},
		/* 51 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action36)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					position540 := position
					depth++
					{
						position541, tokenIndex541, depth541 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l541
						}
						if buffer[position] != rune(':') {
							goto l541
						}
						position++
						goto l542
					l541:
						position, tokenIndex, depth = position541, tokenIndex541, depth541
					}
				l542:
					{
						position543, tokenIndex543, depth543 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l544
						}
						position++
						goto l543
					l544:
						position, tokenIndex, depth = position543, tokenIndex543, depth543
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l538
						}
						position++
					}
				l543:
				l545:
					{
						position546, tokenIndex546, depth546 := position, tokenIndex, depth
						{
							position547, tokenIndex547, depth547 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l548
							}
							position++
							goto l547
						l548:
							position, tokenIndex, depth = position547, tokenIndex547, depth547
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l549
							}
							position++
							goto l547
						l549:
							position, tokenIndex, depth = position547, tokenIndex547, depth547
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l550
							}
							position++
							goto l547
						l550:
							position, tokenIndex, depth = position547, tokenIndex547, depth547
							if buffer[position] != rune('_') {
								goto l551
							}
							position++
							goto l547
						l551:
							position, tokenIndex, depth = position547, tokenIndex547, depth547
							if buffer[position] != rune('.') {
								goto l546
							}
							position++
						}
					l547:
						goto l545
					l546:
						position, tokenIndex, depth = position546, tokenIndex546, depth546
					}
					depth--
					add(rulePegText, position540)
				}
				if !_rules[ruleAction36]() {
					goto l538
				}
				depth--
				add(ruleRowValue, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 52 NumericLiteral <- <(<('-'? [0-9]+)> Action37)> */
		func() bool {
			position552, tokenIndex552, depth552 := position, tokenIndex, depth
			{
				position553 := position
				depth++
				{
					position554 := position
					depth++
					{
						position555, tokenIndex555, depth555 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l555
						}
						position++
						goto l556
					l555:
						position, tokenIndex, depth = position555, tokenIndex555, depth555
					}
				l556:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l552
					}
					position++
				l557:
					{
						position558, tokenIndex558, depth558 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l558
						}
						position++
						goto l557
					l558:
						position, tokenIndex, depth = position558, tokenIndex558, depth558
					}
					depth--
					add(rulePegText, position554)
				}
				if !_rules[ruleAction37]() {
					goto l552
				}
				depth--
				add(ruleNumericLiteral, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 53 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action38)> */
		func() bool {
			position559, tokenIndex559, depth559 := position, tokenIndex, depth
			{
				position560 := position
				depth++
				{
					position561 := position
					depth++
					{
						position562, tokenIndex562, depth562 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l562
						}
						position++
						goto l563
					l562:
						position, tokenIndex, depth = position562, tokenIndex562, depth562
					}
				l563:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l559
					}
					position++
				l564:
					{
						position565, tokenIndex565, depth565 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l565
						}
						position++
						goto l564
					l565:
						position, tokenIndex, depth = position565, tokenIndex565, depth565
					}
					if buffer[position] != rune('.') {
						goto l559
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l559
					}
					position++
				l566:
					{
						position567, tokenIndex567, depth567 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l567
						}
						position++
						goto l566
					l567:
						position, tokenIndex, depth = position567, tokenIndex567, depth567
					}
					depth--
					add(rulePegText, position561)
				}
				if !_rules[ruleAction38]() {
					goto l559
				}
				depth--
				add(ruleFloatLiteral, position560)
			}
			return true
		l559:
			position, tokenIndex, depth = position559, tokenIndex559, depth559
			return false
		},
		/* 54 Function <- <(<ident> Action39)> */
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
				if !_rules[ruleAction39]() {
					goto l568
				}
				depth--
				add(ruleFunction, position569)
			}
			return true
		l568:
			position, tokenIndex, depth = position568, tokenIndex568, depth568
			return false
		},
		/* 55 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position571, tokenIndex571, depth571 := position, tokenIndex, depth
			{
				position572 := position
				depth++
				{
					position573, tokenIndex573, depth573 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l574
					}
					goto l573
				l574:
					position, tokenIndex, depth = position573, tokenIndex573, depth573
					if !_rules[ruleFALSE]() {
						goto l571
					}
				}
			l573:
				depth--
				add(ruleBooleanLiteral, position572)
			}
			return true
		l571:
			position, tokenIndex, depth = position571, tokenIndex571, depth571
			return false
		},
		/* 56 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action40)> */
		func() bool {
			position575, tokenIndex575, depth575 := position, tokenIndex, depth
			{
				position576 := position
				depth++
				{
					position577 := position
					depth++
					{
						position578, tokenIndex578, depth578 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l579
						}
						position++
						goto l578
					l579:
						position, tokenIndex, depth = position578, tokenIndex578, depth578
						if buffer[position] != rune('T') {
							goto l575
						}
						position++
					}
				l578:
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l581
						}
						position++
						goto l580
					l581:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
						if buffer[position] != rune('R') {
							goto l575
						}
						position++
					}
				l580:
					{
						position582, tokenIndex582, depth582 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l583
						}
						position++
						goto l582
					l583:
						position, tokenIndex, depth = position582, tokenIndex582, depth582
						if buffer[position] != rune('U') {
							goto l575
						}
						position++
					}
				l582:
					{
						position584, tokenIndex584, depth584 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l585
						}
						position++
						goto l584
					l585:
						position, tokenIndex, depth = position584, tokenIndex584, depth584
						if buffer[position] != rune('E') {
							goto l575
						}
						position++
					}
				l584:
					depth--
					add(rulePegText, position577)
				}
				if !_rules[ruleAction40]() {
					goto l575
				}
				depth--
				add(ruleTRUE, position576)
			}
			return true
		l575:
			position, tokenIndex, depth = position575, tokenIndex575, depth575
			return false
		},
		/* 57 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action41)> */
		func() bool {
			position586, tokenIndex586, depth586 := position, tokenIndex, depth
			{
				position587 := position
				depth++
				{
					position588 := position
					depth++
					{
						position589, tokenIndex589, depth589 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l590
						}
						position++
						goto l589
					l590:
						position, tokenIndex, depth = position589, tokenIndex589, depth589
						if buffer[position] != rune('F') {
							goto l586
						}
						position++
					}
				l589:
					{
						position591, tokenIndex591, depth591 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l592
						}
						position++
						goto l591
					l592:
						position, tokenIndex, depth = position591, tokenIndex591, depth591
						if buffer[position] != rune('A') {
							goto l586
						}
						position++
					}
				l591:
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l594
						}
						position++
						goto l593
					l594:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
						if buffer[position] != rune('L') {
							goto l586
						}
						position++
					}
				l593:
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l596
						}
						position++
						goto l595
					l596:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
						if buffer[position] != rune('S') {
							goto l586
						}
						position++
					}
				l595:
					{
						position597, tokenIndex597, depth597 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l598
						}
						position++
						goto l597
					l598:
						position, tokenIndex, depth = position597, tokenIndex597, depth597
						if buffer[position] != rune('E') {
							goto l586
						}
						position++
					}
				l597:
					depth--
					add(rulePegText, position588)
				}
				if !_rules[ruleAction41]() {
					goto l586
				}
				depth--
				add(ruleFALSE, position587)
			}
			return true
		l586:
			position, tokenIndex, depth = position586, tokenIndex586, depth586
			return false
		},
		/* 58 Wildcard <- <(<'*'> Action42)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				{
					position601 := position
					depth++
					if buffer[position] != rune('*') {
						goto l599
					}
					position++
					depth--
					add(rulePegText, position601)
				}
				if !_rules[ruleAction42]() {
					goto l599
				}
				depth--
				add(ruleWildcard, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 59 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action43)> */
		func() bool {
			position602, tokenIndex602, depth602 := position, tokenIndex, depth
			{
				position603 := position
				depth++
				{
					position604 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l602
					}
					position++
				l605:
					{
						position606, tokenIndex606, depth606 := position, tokenIndex, depth
						{
							position607, tokenIndex607, depth607 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l608
							}
							position++
							if buffer[position] != rune('\'') {
								goto l608
							}
							position++
							goto l607
						l608:
							position, tokenIndex, depth = position607, tokenIndex607, depth607
							{
								position609, tokenIndex609, depth609 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l609
								}
								position++
								goto l606
							l609:
								position, tokenIndex, depth = position609, tokenIndex609, depth609
							}
							if !matchDot() {
								goto l606
							}
						}
					l607:
						goto l605
					l606:
						position, tokenIndex, depth = position606, tokenIndex606, depth606
					}
					if buffer[position] != rune('\'') {
						goto l602
					}
					position++
					depth--
					add(rulePegText, position604)
				}
				if !_rules[ruleAction43]() {
					goto l602
				}
				depth--
				add(ruleStringLiteral, position603)
			}
			return true
		l602:
			position, tokenIndex, depth = position602, tokenIndex602, depth602
			return false
		},
		/* 60 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action44)> */
		func() bool {
			position610, tokenIndex610, depth610 := position, tokenIndex, depth
			{
				position611 := position
				depth++
				{
					position612 := position
					depth++
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('I') {
							goto l610
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('S') {
							goto l610
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('T') {
							goto l610
						}
						position++
					}
				l617:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('R') {
							goto l610
						}
						position++
					}
				l619:
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l622
						}
						position++
						goto l621
					l622:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
						if buffer[position] != rune('E') {
							goto l610
						}
						position++
					}
				l621:
					{
						position623, tokenIndex623, depth623 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l624
						}
						position++
						goto l623
					l624:
						position, tokenIndex, depth = position623, tokenIndex623, depth623
						if buffer[position] != rune('A') {
							goto l610
						}
						position++
					}
				l623:
					{
						position625, tokenIndex625, depth625 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l626
						}
						position++
						goto l625
					l626:
						position, tokenIndex, depth = position625, tokenIndex625, depth625
						if buffer[position] != rune('M') {
							goto l610
						}
						position++
					}
				l625:
					depth--
					add(rulePegText, position612)
				}
				if !_rules[ruleAction44]() {
					goto l610
				}
				depth--
				add(ruleISTREAM, position611)
			}
			return true
		l610:
			position, tokenIndex, depth = position610, tokenIndex610, depth610
			return false
		},
		/* 61 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action45)> */
		func() bool {
			position627, tokenIndex627, depth627 := position, tokenIndex, depth
			{
				position628 := position
				depth++
				{
					position629 := position
					depth++
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('D') {
							goto l627
						}
						position++
					}
				l630:
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l633
						}
						position++
						goto l632
					l633:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
						if buffer[position] != rune('S') {
							goto l627
						}
						position++
					}
				l632:
					{
						position634, tokenIndex634, depth634 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l635
						}
						position++
						goto l634
					l635:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
						if buffer[position] != rune('T') {
							goto l627
						}
						position++
					}
				l634:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l637
						}
						position++
						goto l636
					l637:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
						if buffer[position] != rune('R') {
							goto l627
						}
						position++
					}
				l636:
					{
						position638, tokenIndex638, depth638 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l639
						}
						position++
						goto l638
					l639:
						position, tokenIndex, depth = position638, tokenIndex638, depth638
						if buffer[position] != rune('E') {
							goto l627
						}
						position++
					}
				l638:
					{
						position640, tokenIndex640, depth640 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l641
						}
						position++
						goto l640
					l641:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
						if buffer[position] != rune('A') {
							goto l627
						}
						position++
					}
				l640:
					{
						position642, tokenIndex642, depth642 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l643
						}
						position++
						goto l642
					l643:
						position, tokenIndex, depth = position642, tokenIndex642, depth642
						if buffer[position] != rune('M') {
							goto l627
						}
						position++
					}
				l642:
					depth--
					add(rulePegText, position629)
				}
				if !_rules[ruleAction45]() {
					goto l627
				}
				depth--
				add(ruleDSTREAM, position628)
			}
			return true
		l627:
			position, tokenIndex, depth = position627, tokenIndex627, depth627
			return false
		},
		/* 62 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action46)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				{
					position646 := position
					depth++
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
							goto l644
						}
						position++
					}
				l647:
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l650
						}
						position++
						goto l649
					l650:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						if buffer[position] != rune('S') {
							goto l644
						}
						position++
					}
				l649:
					{
						position651, tokenIndex651, depth651 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l652
						}
						position++
						goto l651
					l652:
						position, tokenIndex, depth = position651, tokenIndex651, depth651
						if buffer[position] != rune('T') {
							goto l644
						}
						position++
					}
				l651:
					{
						position653, tokenIndex653, depth653 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l654
						}
						position++
						goto l653
					l654:
						position, tokenIndex, depth = position653, tokenIndex653, depth653
						if buffer[position] != rune('R') {
							goto l644
						}
						position++
					}
				l653:
					{
						position655, tokenIndex655, depth655 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l656
						}
						position++
						goto l655
					l656:
						position, tokenIndex, depth = position655, tokenIndex655, depth655
						if buffer[position] != rune('E') {
							goto l644
						}
						position++
					}
				l655:
					{
						position657, tokenIndex657, depth657 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l658
						}
						position++
						goto l657
					l658:
						position, tokenIndex, depth = position657, tokenIndex657, depth657
						if buffer[position] != rune('A') {
							goto l644
						}
						position++
					}
				l657:
					{
						position659, tokenIndex659, depth659 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l660
						}
						position++
						goto l659
					l660:
						position, tokenIndex, depth = position659, tokenIndex659, depth659
						if buffer[position] != rune('M') {
							goto l644
						}
						position++
					}
				l659:
					depth--
					add(rulePegText, position646)
				}
				if !_rules[ruleAction46]() {
					goto l644
				}
				depth--
				add(ruleRSTREAM, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 63 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action47)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				{
					position663 := position
					depth++
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l665
						}
						position++
						goto l664
					l665:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if buffer[position] != rune('T') {
							goto l661
						}
						position++
					}
				l664:
					{
						position666, tokenIndex666, depth666 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l667
						}
						position++
						goto l666
					l667:
						position, tokenIndex, depth = position666, tokenIndex666, depth666
						if buffer[position] != rune('U') {
							goto l661
						}
						position++
					}
				l666:
					{
						position668, tokenIndex668, depth668 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l669
						}
						position++
						goto l668
					l669:
						position, tokenIndex, depth = position668, tokenIndex668, depth668
						if buffer[position] != rune('P') {
							goto l661
						}
						position++
					}
				l668:
					{
						position670, tokenIndex670, depth670 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l671
						}
						position++
						goto l670
					l671:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
						if buffer[position] != rune('L') {
							goto l661
						}
						position++
					}
				l670:
					{
						position672, tokenIndex672, depth672 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l673
						}
						position++
						goto l672
					l673:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						if buffer[position] != rune('E') {
							goto l661
						}
						position++
					}
				l672:
					{
						position674, tokenIndex674, depth674 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l675
						}
						position++
						goto l674
					l675:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
						if buffer[position] != rune('S') {
							goto l661
						}
						position++
					}
				l674:
					depth--
					add(rulePegText, position663)
				}
				if !_rules[ruleAction47]() {
					goto l661
				}
				depth--
				add(ruleTUPLES, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 64 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action48)> */
		func() bool {
			position676, tokenIndex676, depth676 := position, tokenIndex, depth
			{
				position677 := position
				depth++
				{
					position678 := position
					depth++
					{
						position679, tokenIndex679, depth679 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l680
						}
						position++
						goto l679
					l680:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
						if buffer[position] != rune('S') {
							goto l676
						}
						position++
					}
				l679:
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if buffer[position] != rune('E') {
							goto l676
						}
						position++
					}
				l681:
					{
						position683, tokenIndex683, depth683 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l684
						}
						position++
						goto l683
					l684:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
						if buffer[position] != rune('C') {
							goto l676
						}
						position++
					}
				l683:
					{
						position685, tokenIndex685, depth685 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l686
						}
						position++
						goto l685
					l686:
						position, tokenIndex, depth = position685, tokenIndex685, depth685
						if buffer[position] != rune('O') {
							goto l676
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
							goto l676
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
							goto l676
						}
						position++
					}
				l689:
					{
						position691, tokenIndex691, depth691 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l692
						}
						position++
						goto l691
					l692:
						position, tokenIndex, depth = position691, tokenIndex691, depth691
						if buffer[position] != rune('S') {
							goto l676
						}
						position++
					}
				l691:
					depth--
					add(rulePegText, position678)
				}
				if !_rules[ruleAction48]() {
					goto l676
				}
				depth--
				add(ruleSECONDS, position677)
			}
			return true
		l676:
			position, tokenIndex, depth = position676, tokenIndex676, depth676
			return false
		},
		/* 65 StreamIdentifier <- <(<ident> Action49)> */
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
				if !_rules[ruleAction49]() {
					goto l693
				}
				depth--
				add(ruleStreamIdentifier, position694)
			}
			return true
		l693:
			position, tokenIndex, depth = position693, tokenIndex693, depth693
			return false
		},
		/* 66 SourceSinkType <- <(<ident> Action50)> */
		func() bool {
			position696, tokenIndex696, depth696 := position, tokenIndex, depth
			{
				position697 := position
				depth++
				{
					position698 := position
					depth++
					if !_rules[ruleident]() {
						goto l696
					}
					depth--
					add(rulePegText, position698)
				}
				if !_rules[ruleAction50]() {
					goto l696
				}
				depth--
				add(ruleSourceSinkType, position697)
			}
			return true
		l696:
			position, tokenIndex, depth = position696, tokenIndex696, depth696
			return false
		},
		/* 67 SourceSinkParamKey <- <(<ident> Action51)> */
		func() bool {
			position699, tokenIndex699, depth699 := position, tokenIndex, depth
			{
				position700 := position
				depth++
				{
					position701 := position
					depth++
					if !_rules[ruleident]() {
						goto l699
					}
					depth--
					add(rulePegText, position701)
				}
				if !_rules[ruleAction51]() {
					goto l699
				}
				depth--
				add(ruleSourceSinkParamKey, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 68 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action52)> */
		func() bool {
			position702, tokenIndex702, depth702 := position, tokenIndex, depth
			{
				position703 := position
				depth++
				{
					position704 := position
					depth++
					{
						position705, tokenIndex705, depth705 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l706
						}
						position++
						goto l705
					l706:
						position, tokenIndex, depth = position705, tokenIndex705, depth705
						if buffer[position] != rune('O') {
							goto l702
						}
						position++
					}
				l705:
					{
						position707, tokenIndex707, depth707 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l708
						}
						position++
						goto l707
					l708:
						position, tokenIndex, depth = position707, tokenIndex707, depth707
						if buffer[position] != rune('R') {
							goto l702
						}
						position++
					}
				l707:
					depth--
					add(rulePegText, position704)
				}
				if !_rules[ruleAction52]() {
					goto l702
				}
				depth--
				add(ruleOr, position703)
			}
			return true
		l702:
			position, tokenIndex, depth = position702, tokenIndex702, depth702
			return false
		},
		/* 69 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action53)> */
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
						if buffer[position] != rune('a') {
							goto l713
						}
						position++
						goto l712
					l713:
						position, tokenIndex, depth = position712, tokenIndex712, depth712
						if buffer[position] != rune('A') {
							goto l709
						}
						position++
					}
				l712:
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l715
						}
						position++
						goto l714
					l715:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
						if buffer[position] != rune('N') {
							goto l709
						}
						position++
					}
				l714:
					{
						position716, tokenIndex716, depth716 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l717
						}
						position++
						goto l716
					l717:
						position, tokenIndex, depth = position716, tokenIndex716, depth716
						if buffer[position] != rune('D') {
							goto l709
						}
						position++
					}
				l716:
					depth--
					add(rulePegText, position711)
				}
				if !_rules[ruleAction53]() {
					goto l709
				}
				depth--
				add(ruleAnd, position710)
			}
			return true
		l709:
			position, tokenIndex, depth = position709, tokenIndex709, depth709
			return false
		},
		/* 70 Equal <- <(<'='> Action54)> */
		func() bool {
			position718, tokenIndex718, depth718 := position, tokenIndex, depth
			{
				position719 := position
				depth++
				{
					position720 := position
					depth++
					if buffer[position] != rune('=') {
						goto l718
					}
					position++
					depth--
					add(rulePegText, position720)
				}
				if !_rules[ruleAction54]() {
					goto l718
				}
				depth--
				add(ruleEqual, position719)
			}
			return true
		l718:
			position, tokenIndex, depth = position718, tokenIndex718, depth718
			return false
		},
		/* 71 Less <- <(<'<'> Action55)> */
		func() bool {
			position721, tokenIndex721, depth721 := position, tokenIndex, depth
			{
				position722 := position
				depth++
				{
					position723 := position
					depth++
					if buffer[position] != rune('<') {
						goto l721
					}
					position++
					depth--
					add(rulePegText, position723)
				}
				if !_rules[ruleAction55]() {
					goto l721
				}
				depth--
				add(ruleLess, position722)
			}
			return true
		l721:
			position, tokenIndex, depth = position721, tokenIndex721, depth721
			return false
		},
		/* 72 LessOrEqual <- <(<('<' '=')> Action56)> */
		func() bool {
			position724, tokenIndex724, depth724 := position, tokenIndex, depth
			{
				position725 := position
				depth++
				{
					position726 := position
					depth++
					if buffer[position] != rune('<') {
						goto l724
					}
					position++
					if buffer[position] != rune('=') {
						goto l724
					}
					position++
					depth--
					add(rulePegText, position726)
				}
				if !_rules[ruleAction56]() {
					goto l724
				}
				depth--
				add(ruleLessOrEqual, position725)
			}
			return true
		l724:
			position, tokenIndex, depth = position724, tokenIndex724, depth724
			return false
		},
		/* 73 Greater <- <(<'>'> Action57)> */
		func() bool {
			position727, tokenIndex727, depth727 := position, tokenIndex, depth
			{
				position728 := position
				depth++
				{
					position729 := position
					depth++
					if buffer[position] != rune('>') {
						goto l727
					}
					position++
					depth--
					add(rulePegText, position729)
				}
				if !_rules[ruleAction57]() {
					goto l727
				}
				depth--
				add(ruleGreater, position728)
			}
			return true
		l727:
			position, tokenIndex, depth = position727, tokenIndex727, depth727
			return false
		},
		/* 74 GreaterOrEqual <- <(<('>' '=')> Action58)> */
		func() bool {
			position730, tokenIndex730, depth730 := position, tokenIndex, depth
			{
				position731 := position
				depth++
				{
					position732 := position
					depth++
					if buffer[position] != rune('>') {
						goto l730
					}
					position++
					if buffer[position] != rune('=') {
						goto l730
					}
					position++
					depth--
					add(rulePegText, position732)
				}
				if !_rules[ruleAction58]() {
					goto l730
				}
				depth--
				add(ruleGreaterOrEqual, position731)
			}
			return true
		l730:
			position, tokenIndex, depth = position730, tokenIndex730, depth730
			return false
		},
		/* 75 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action59)> */
		func() bool {
			position733, tokenIndex733, depth733 := position, tokenIndex, depth
			{
				position734 := position
				depth++
				{
					position735 := position
					depth++
					{
						position736, tokenIndex736, depth736 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l737
						}
						position++
						if buffer[position] != rune('=') {
							goto l737
						}
						position++
						goto l736
					l737:
						position, tokenIndex, depth = position736, tokenIndex736, depth736
						if buffer[position] != rune('<') {
							goto l733
						}
						position++
						if buffer[position] != rune('>') {
							goto l733
						}
						position++
					}
				l736:
					depth--
					add(rulePegText, position735)
				}
				if !_rules[ruleAction59]() {
					goto l733
				}
				depth--
				add(ruleNotEqual, position734)
			}
			return true
		l733:
			position, tokenIndex, depth = position733, tokenIndex733, depth733
			return false
		},
		/* 76 Plus <- <(<'+'> Action60)> */
		func() bool {
			position738, tokenIndex738, depth738 := position, tokenIndex, depth
			{
				position739 := position
				depth++
				{
					position740 := position
					depth++
					if buffer[position] != rune('+') {
						goto l738
					}
					position++
					depth--
					add(rulePegText, position740)
				}
				if !_rules[ruleAction60]() {
					goto l738
				}
				depth--
				add(rulePlus, position739)
			}
			return true
		l738:
			position, tokenIndex, depth = position738, tokenIndex738, depth738
			return false
		},
		/* 77 Minus <- <(<'-'> Action61)> */
		func() bool {
			position741, tokenIndex741, depth741 := position, tokenIndex, depth
			{
				position742 := position
				depth++
				{
					position743 := position
					depth++
					if buffer[position] != rune('-') {
						goto l741
					}
					position++
					depth--
					add(rulePegText, position743)
				}
				if !_rules[ruleAction61]() {
					goto l741
				}
				depth--
				add(ruleMinus, position742)
			}
			return true
		l741:
			position, tokenIndex, depth = position741, tokenIndex741, depth741
			return false
		},
		/* 78 Multiply <- <(<'*'> Action62)> */
		func() bool {
			position744, tokenIndex744, depth744 := position, tokenIndex, depth
			{
				position745 := position
				depth++
				{
					position746 := position
					depth++
					if buffer[position] != rune('*') {
						goto l744
					}
					position++
					depth--
					add(rulePegText, position746)
				}
				if !_rules[ruleAction62]() {
					goto l744
				}
				depth--
				add(ruleMultiply, position745)
			}
			return true
		l744:
			position, tokenIndex, depth = position744, tokenIndex744, depth744
			return false
		},
		/* 79 Divide <- <(<'/'> Action63)> */
		func() bool {
			position747, tokenIndex747, depth747 := position, tokenIndex, depth
			{
				position748 := position
				depth++
				{
					position749 := position
					depth++
					if buffer[position] != rune('/') {
						goto l747
					}
					position++
					depth--
					add(rulePegText, position749)
				}
				if !_rules[ruleAction63]() {
					goto l747
				}
				depth--
				add(ruleDivide, position748)
			}
			return true
		l747:
			position, tokenIndex, depth = position747, tokenIndex747, depth747
			return false
		},
		/* 80 Modulo <- <(<'%'> Action64)> */
		func() bool {
			position750, tokenIndex750, depth750 := position, tokenIndex, depth
			{
				position751 := position
				depth++
				{
					position752 := position
					depth++
					if buffer[position] != rune('%') {
						goto l750
					}
					position++
					depth--
					add(rulePegText, position752)
				}
				if !_rules[ruleAction64]() {
					goto l750
				}
				depth--
				add(ruleModulo, position751)
			}
			return true
		l750:
			position, tokenIndex, depth = position750, tokenIndex750, depth750
			return false
		},
		/* 81 Identifier <- <(<ident> Action65)> */
		func() bool {
			position753, tokenIndex753, depth753 := position, tokenIndex, depth
			{
				position754 := position
				depth++
				{
					position755 := position
					depth++
					if !_rules[ruleident]() {
						goto l753
					}
					depth--
					add(rulePegText, position755)
				}
				if !_rules[ruleAction65]() {
					goto l753
				}
				depth--
				add(ruleIdentifier, position754)
			}
			return true
		l753:
			position, tokenIndex, depth = position753, tokenIndex753, depth753
			return false
		},
		/* 82 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position756, tokenIndex756, depth756 := position, tokenIndex, depth
			{
				position757 := position
				depth++
				{
					position758, tokenIndex758, depth758 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l759
					}
					position++
					goto l758
				l759:
					position, tokenIndex, depth = position758, tokenIndex758, depth758
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l756
					}
					position++
				}
			l758:
			l760:
				{
					position761, tokenIndex761, depth761 := position, tokenIndex, depth
					{
						position762, tokenIndex762, depth762 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l763
						}
						position++
						goto l762
					l763:
						position, tokenIndex, depth = position762, tokenIndex762, depth762
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l764
						}
						position++
						goto l762
					l764:
						position, tokenIndex, depth = position762, tokenIndex762, depth762
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l765
						}
						position++
						goto l762
					l765:
						position, tokenIndex, depth = position762, tokenIndex762, depth762
						if buffer[position] != rune('_') {
							goto l761
						}
						position++
					}
				l762:
					goto l760
				l761:
					position, tokenIndex, depth = position761, tokenIndex761, depth761
				}
				depth--
				add(ruleident, position757)
			}
			return true
		l756:
			position, tokenIndex, depth = position756, tokenIndex756, depth756
			return false
		},
		/* 83 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position767 := position
				depth++
			l768:
				{
					position769, tokenIndex769, depth769 := position, tokenIndex, depth
					{
						position770, tokenIndex770, depth770 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l771
						}
						position++
						goto l770
					l771:
						position, tokenIndex, depth = position770, tokenIndex770, depth770
						if buffer[position] != rune('\t') {
							goto l772
						}
						position++
						goto l770
					l772:
						position, tokenIndex, depth = position770, tokenIndex770, depth770
						if buffer[position] != rune('\n') {
							goto l769
						}
						position++
					}
				l770:
					goto l768
				l769:
					position, tokenIndex, depth = position769, tokenIndex769, depth769
				}
				depth--
				add(rulesp, position767)
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
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 139 Action53 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 140 Action54 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 141 Action55 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 142 Action56 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 143 Action57 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 144 Action58 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 145 Action59 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 146 Action60 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 147 Action61 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 148 Action62 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 149 Action63 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 150 Action64 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 151 Action65 <- <{
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

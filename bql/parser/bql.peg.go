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
	rulePausedOpt
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
	rulePaused
	ruleUnpaused
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
	ruleAction68
	ruleAction69
	ruleAction70

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
	"PausedOpt",
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
	"Paused",
	"Unpaused",
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
	"Action68",
	"Action69",
	"Action70",

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
	rules  [163]func() bool
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

			p.EnsureKeywordPresent(begin, end)

		case ruleAction30:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction31:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction32:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction33:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction34:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction35:

			p.AssembleFuncApp()

		case ruleAction36:

			p.AssembleExpressions(begin, end)

		case ruleAction37:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction38:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction43:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction44:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction45:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction47:

			p.PushComponent(begin, end, Istream)

		case ruleAction48:

			p.PushComponent(begin, end, Dstream)

		case ruleAction49:

			p.PushComponent(begin, end, Rstream)

		case ruleAction50:

			p.PushComponent(begin, end, Tuples)

		case ruleAction51:

			p.PushComponent(begin, end, Seconds)

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction53:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction54:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction55:

			p.PushComponent(begin, end, Yes)

		case ruleAction56:

			p.PushComponent(begin, end, No)

		case ruleAction57:

			p.PushComponent(begin, end, Or)

		case ruleAction58:

			p.PushComponent(begin, end, And)

		case ruleAction59:

			p.PushComponent(begin, end, Equal)

		case ruleAction60:

			p.PushComponent(begin, end, Less)

		case ruleAction61:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction62:

			p.PushComponent(begin, end, Greater)

		case ruleAction63:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction64:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction65:

			p.PushComponent(begin, end, Plus)

		case ruleAction66:

			p.PushComponent(begin, end, Minus)

		case ruleAction67:

			p.PushComponent(begin, end, Multiply)

		case ruleAction68:

			p.PushComponent(begin, end, Divide)

		case ruleAction69:

			p.PushComponent(begin, end, Modulo)

		case ruleAction70:

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
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
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
				if !_rules[rulePausedOpt]() {
					goto l75
				}
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
		/* 38 PausedOpt <- <(<(Paused / Unpaused)?> Action29)> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				{
					position464 := position
					depth++
					{
						position465, tokenIndex465, depth465 := position, tokenIndex, depth
						{
							position467, tokenIndex467, depth467 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l468
							}
							goto l467
						l468:
							position, tokenIndex, depth = position467, tokenIndex467, depth467
							if !_rules[ruleUnpaused]() {
								goto l465
							}
						}
					l467:
						goto l466
					l465:
						position, tokenIndex, depth = position465, tokenIndex465, depth465
					}
				l466:
					depth--
					add(rulePegText, position464)
				}
				if !_rules[ruleAction29]() {
					goto l462
				}
				depth--
				add(rulePausedOpt, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 39 Expression <- <orExpr> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l469
				}
				depth--
				add(ruleExpression, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 40 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action30)> */
		func() bool {
			position471, tokenIndex471, depth471 := position, tokenIndex, depth
			{
				position472 := position
				depth++
				{
					position473 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l471
					}
					if !_rules[rulesp]() {
						goto l471
					}
					{
						position474, tokenIndex474, depth474 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l474
						}
						if !_rules[rulesp]() {
							goto l474
						}
						if !_rules[ruleandExpr]() {
							goto l474
						}
						goto l475
					l474:
						position, tokenIndex, depth = position474, tokenIndex474, depth474
					}
				l475:
					depth--
					add(rulePegText, position473)
				}
				if !_rules[ruleAction30]() {
					goto l471
				}
				depth--
				add(ruleorExpr, position472)
			}
			return true
		l471:
			position, tokenIndex, depth = position471, tokenIndex471, depth471
			return false
		},
		/* 41 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action31)> */
		func() bool {
			position476, tokenIndex476, depth476 := position, tokenIndex, depth
			{
				position477 := position
				depth++
				{
					position478 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l476
					}
					if !_rules[rulesp]() {
						goto l476
					}
					{
						position479, tokenIndex479, depth479 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l479
						}
						if !_rules[rulesp]() {
							goto l479
						}
						if !_rules[rulecomparisonExpr]() {
							goto l479
						}
						goto l480
					l479:
						position, tokenIndex, depth = position479, tokenIndex479, depth479
					}
				l480:
					depth--
					add(rulePegText, position478)
				}
				if !_rules[ruleAction31]() {
					goto l476
				}
				depth--
				add(ruleandExpr, position477)
			}
			return true
		l476:
			position, tokenIndex, depth = position476, tokenIndex476, depth476
			return false
		},
		/* 42 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action32)> */
		func() bool {
			position481, tokenIndex481, depth481 := position, tokenIndex, depth
			{
				position482 := position
				depth++
				{
					position483 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l481
					}
					if !_rules[rulesp]() {
						goto l481
					}
					{
						position484, tokenIndex484, depth484 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l484
						}
						if !_rules[rulesp]() {
							goto l484
						}
						if !_rules[ruletermExpr]() {
							goto l484
						}
						goto l485
					l484:
						position, tokenIndex, depth = position484, tokenIndex484, depth484
					}
				l485:
					depth--
					add(rulePegText, position483)
				}
				if !_rules[ruleAction32]() {
					goto l481
				}
				depth--
				add(rulecomparisonExpr, position482)
			}
			return true
		l481:
			position, tokenIndex, depth = position481, tokenIndex481, depth481
			return false
		},
		/* 43 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action33)> */
		func() bool {
			position486, tokenIndex486, depth486 := position, tokenIndex, depth
			{
				position487 := position
				depth++
				{
					position488 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l486
					}
					if !_rules[rulesp]() {
						goto l486
					}
					{
						position489, tokenIndex489, depth489 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l489
						}
						if !_rules[rulesp]() {
							goto l489
						}
						if !_rules[ruleproductExpr]() {
							goto l489
						}
						goto l490
					l489:
						position, tokenIndex, depth = position489, tokenIndex489, depth489
					}
				l490:
					depth--
					add(rulePegText, position488)
				}
				if !_rules[ruleAction33]() {
					goto l486
				}
				depth--
				add(ruletermExpr, position487)
			}
			return true
		l486:
			position, tokenIndex, depth = position486, tokenIndex486, depth486
			return false
		},
		/* 44 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action34)> */
		func() bool {
			position491, tokenIndex491, depth491 := position, tokenIndex, depth
			{
				position492 := position
				depth++
				{
					position493 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l491
					}
					if !_rules[rulesp]() {
						goto l491
					}
					{
						position494, tokenIndex494, depth494 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l494
						}
						if !_rules[rulesp]() {
							goto l494
						}
						if !_rules[rulebaseExpr]() {
							goto l494
						}
						goto l495
					l494:
						position, tokenIndex, depth = position494, tokenIndex494, depth494
					}
				l495:
					depth--
					add(rulePegText, position493)
				}
				if !_rules[ruleAction34]() {
					goto l491
				}
				depth--
				add(ruleproductExpr, position492)
			}
			return true
		l491:
			position, tokenIndex, depth = position491, tokenIndex491, depth491
			return false
		},
		/* 45 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position496, tokenIndex496, depth496 := position, tokenIndex, depth
			{
				position497 := position
				depth++
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l499
					}
					position++
					if !_rules[rulesp]() {
						goto l499
					}
					if !_rules[ruleExpression]() {
						goto l499
					}
					if !_rules[rulesp]() {
						goto l499
					}
					if buffer[position] != rune(')') {
						goto l499
					}
					position++
					goto l498
				l499:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleBooleanLiteral]() {
						goto l500
					}
					goto l498
				l500:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleFuncApp]() {
						goto l501
					}
					goto l498
				l501:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleRowMeta]() {
						goto l502
					}
					goto l498
				l502:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleRowValue]() {
						goto l503
					}
					goto l498
				l503:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleLiteral]() {
						goto l496
					}
				}
			l498:
				depth--
				add(rulebaseExpr, position497)
			}
			return true
		l496:
			position, tokenIndex, depth = position496, tokenIndex496, depth496
			return false
		},
		/* 46 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action35)> */
		func() bool {
			position504, tokenIndex504, depth504 := position, tokenIndex, depth
			{
				position505 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l504
				}
				if !_rules[rulesp]() {
					goto l504
				}
				if buffer[position] != rune('(') {
					goto l504
				}
				position++
				if !_rules[rulesp]() {
					goto l504
				}
				if !_rules[ruleFuncParams]() {
					goto l504
				}
				if !_rules[rulesp]() {
					goto l504
				}
				if buffer[position] != rune(')') {
					goto l504
				}
				position++
				if !_rules[ruleAction35]() {
					goto l504
				}
				depth--
				add(ruleFuncApp, position505)
			}
			return true
		l504:
			position, tokenIndex, depth = position504, tokenIndex504, depth504
			return false
		},
		/* 47 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action36)> */
		func() bool {
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				{
					position508 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l506
					}
					if !_rules[rulesp]() {
						goto l506
					}
				l509:
					{
						position510, tokenIndex510, depth510 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l510
						}
						position++
						if !_rules[rulesp]() {
							goto l510
						}
						if !_rules[ruleExpression]() {
							goto l510
						}
						goto l509
					l510:
						position, tokenIndex, depth = position510, tokenIndex510, depth510
					}
					depth--
					add(rulePegText, position508)
				}
				if !_rules[ruleAction36]() {
					goto l506
				}
				depth--
				add(ruleFuncParams, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 48 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position511, tokenIndex511, depth511 := position, tokenIndex, depth
			{
				position512 := position
				depth++
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l514
					}
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleNumericLiteral]() {
						goto l515
					}
					goto l513
				l515:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if !_rules[ruleStringLiteral]() {
						goto l511
					}
				}
			l513:
				depth--
				add(ruleLiteral, position512)
			}
			return true
		l511:
			position, tokenIndex, depth = position511, tokenIndex511, depth511
			return false
		},
		/* 49 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position516, tokenIndex516, depth516 := position, tokenIndex, depth
			{
				position517 := position
				depth++
				{
					position518, tokenIndex518, depth518 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l519
					}
					goto l518
				l519:
					position, tokenIndex, depth = position518, tokenIndex518, depth518
					if !_rules[ruleNotEqual]() {
						goto l520
					}
					goto l518
				l520:
					position, tokenIndex, depth = position518, tokenIndex518, depth518
					if !_rules[ruleLessOrEqual]() {
						goto l521
					}
					goto l518
				l521:
					position, tokenIndex, depth = position518, tokenIndex518, depth518
					if !_rules[ruleLess]() {
						goto l522
					}
					goto l518
				l522:
					position, tokenIndex, depth = position518, tokenIndex518, depth518
					if !_rules[ruleGreaterOrEqual]() {
						goto l523
					}
					goto l518
				l523:
					position, tokenIndex, depth = position518, tokenIndex518, depth518
					if !_rules[ruleGreater]() {
						goto l524
					}
					goto l518
				l524:
					position, tokenIndex, depth = position518, tokenIndex518, depth518
					if !_rules[ruleNotEqual]() {
						goto l516
					}
				}
			l518:
				depth--
				add(ruleComparisonOp, position517)
			}
			return true
		l516:
			position, tokenIndex, depth = position516, tokenIndex516, depth516
			return false
		},
		/* 50 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position525, tokenIndex525, depth525 := position, tokenIndex, depth
			{
				position526 := position
				depth++
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l528
					}
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if !_rules[ruleMinus]() {
						goto l525
					}
				}
			l527:
				depth--
				add(rulePlusMinusOp, position526)
			}
			return true
		l525:
			position, tokenIndex, depth = position525, tokenIndex525, depth525
			return false
		},
		/* 51 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l532
					}
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleDivide]() {
						goto l533
					}
					goto l531
				l533:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleModulo]() {
						goto l529
					}
				}
			l531:
				depth--
				add(ruleMultDivOp, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 52 Stream <- <(<ident> Action37)> */
		func() bool {
			position534, tokenIndex534, depth534 := position, tokenIndex, depth
			{
				position535 := position
				depth++
				{
					position536 := position
					depth++
					if !_rules[ruleident]() {
						goto l534
					}
					depth--
					add(rulePegText, position536)
				}
				if !_rules[ruleAction37]() {
					goto l534
				}
				depth--
				add(ruleStream, position535)
			}
			return true
		l534:
			position, tokenIndex, depth = position534, tokenIndex534, depth534
			return false
		},
		/* 53 RowMeta <- <RowTimestamp> */
		func() bool {
			position537, tokenIndex537, depth537 := position, tokenIndex, depth
			{
				position538 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l537
				}
				depth--
				add(ruleRowMeta, position538)
			}
			return true
		l537:
			position, tokenIndex, depth = position537, tokenIndex537, depth537
			return false
		},
		/* 54 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action38)> */
		func() bool {
			position539, tokenIndex539, depth539 := position, tokenIndex, depth
			{
				position540 := position
				depth++
				{
					position541 := position
					depth++
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l542
						}
						if buffer[position] != rune(':') {
							goto l542
						}
						position++
						goto l543
					l542:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
					}
				l543:
					if buffer[position] != rune('t') {
						goto l539
					}
					position++
					if buffer[position] != rune('s') {
						goto l539
					}
					position++
					if buffer[position] != rune('(') {
						goto l539
					}
					position++
					if buffer[position] != rune(')') {
						goto l539
					}
					position++
					depth--
					add(rulePegText, position541)
				}
				if !_rules[ruleAction38]() {
					goto l539
				}
				depth--
				add(ruleRowTimestamp, position540)
			}
			return true
		l539:
			position, tokenIndex, depth = position539, tokenIndex539, depth539
			return false
		},
		/* 55 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action39)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				{
					position546 := position
					depth++
					{
						position547, tokenIndex547, depth547 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l547
						}
						if buffer[position] != rune(':') {
							goto l547
						}
						position++
						goto l548
					l547:
						position, tokenIndex, depth = position547, tokenIndex547, depth547
					}
				l548:
					{
						position549, tokenIndex549, depth549 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l550
						}
						position++
						goto l549
					l550:
						position, tokenIndex, depth = position549, tokenIndex549, depth549
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l544
						}
						position++
					}
				l549:
				l551:
					{
						position552, tokenIndex552, depth552 := position, tokenIndex, depth
						{
							position553, tokenIndex553, depth553 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l554
							}
							position++
							goto l553
						l554:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l555
							}
							position++
							goto l553
						l555:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l556
							}
							position++
							goto l553
						l556:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if buffer[position] != rune('_') {
								goto l557
							}
							position++
							goto l553
						l557:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if buffer[position] != rune('.') {
								goto l552
							}
							position++
						}
					l553:
						goto l551
					l552:
						position, tokenIndex, depth = position552, tokenIndex552, depth552
					}
					depth--
					add(rulePegText, position546)
				}
				if !_rules[ruleAction39]() {
					goto l544
				}
				depth--
				add(ruleRowValue, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 56 NumericLiteral <- <(<('-'? [0-9]+)> Action40)> */
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
					depth--
					add(rulePegText, position560)
				}
				if !_rules[ruleAction40]() {
					goto l558
				}
				depth--
				add(ruleNumericLiteral, position559)
			}
			return true
		l558:
			position, tokenIndex, depth = position558, tokenIndex558, depth558
			return false
		},
		/* 57 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action41)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567 := position
					depth++
					{
						position568, tokenIndex568, depth568 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l568
						}
						position++
						goto l569
					l568:
						position, tokenIndex, depth = position568, tokenIndex568, depth568
					}
				l569:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l565
					}
					position++
				l570:
					{
						position571, tokenIndex571, depth571 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l571
						}
						position++
						goto l570
					l571:
						position, tokenIndex, depth = position571, tokenIndex571, depth571
					}
					if buffer[position] != rune('.') {
						goto l565
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l565
					}
					position++
				l572:
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l573
						}
						position++
						goto l572
					l573:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
					}
					depth--
					add(rulePegText, position567)
				}
				if !_rules[ruleAction41]() {
					goto l565
				}
				depth--
				add(ruleFloatLiteral, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 58 Function <- <(<ident> Action42)> */
		func() bool {
			position574, tokenIndex574, depth574 := position, tokenIndex, depth
			{
				position575 := position
				depth++
				{
					position576 := position
					depth++
					if !_rules[ruleident]() {
						goto l574
					}
					depth--
					add(rulePegText, position576)
				}
				if !_rules[ruleAction42]() {
					goto l574
				}
				depth--
				add(ruleFunction, position575)
			}
			return true
		l574:
			position, tokenIndex, depth = position574, tokenIndex574, depth574
			return false
		},
		/* 59 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position577, tokenIndex577, depth577 := position, tokenIndex, depth
			{
				position578 := position
				depth++
				{
					position579, tokenIndex579, depth579 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l580
					}
					goto l579
				l580:
					position, tokenIndex, depth = position579, tokenIndex579, depth579
					if !_rules[ruleFALSE]() {
						goto l577
					}
				}
			l579:
				depth--
				add(ruleBooleanLiteral, position578)
			}
			return true
		l577:
			position, tokenIndex, depth = position577, tokenIndex577, depth577
			return false
		},
		/* 60 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action43)> */
		func() bool {
			position581, tokenIndex581, depth581 := position, tokenIndex, depth
			{
				position582 := position
				depth++
				{
					position583 := position
					depth++
					{
						position584, tokenIndex584, depth584 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l585
						}
						position++
						goto l584
					l585:
						position, tokenIndex, depth = position584, tokenIndex584, depth584
						if buffer[position] != rune('T') {
							goto l581
						}
						position++
					}
				l584:
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
							goto l581
						}
						position++
					}
				l586:
					{
						position588, tokenIndex588, depth588 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l589
						}
						position++
						goto l588
					l589:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
						if buffer[position] != rune('U') {
							goto l581
						}
						position++
					}
				l588:
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
						if buffer[position] != rune('E') {
							goto l581
						}
						position++
					}
				l590:
					depth--
					add(rulePegText, position583)
				}
				if !_rules[ruleAction43]() {
					goto l581
				}
				depth--
				add(ruleTRUE, position582)
			}
			return true
		l581:
			position, tokenIndex, depth = position581, tokenIndex581, depth581
			return false
		},
		/* 61 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action44)> */
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
						if buffer[position] != rune('f') {
							goto l596
						}
						position++
						goto l595
					l596:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
						if buffer[position] != rune('F') {
							goto l592
						}
						position++
					}
				l595:
					{
						position597, tokenIndex597, depth597 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l598
						}
						position++
						goto l597
					l598:
						position, tokenIndex, depth = position597, tokenIndex597, depth597
						if buffer[position] != rune('A') {
							goto l592
						}
						position++
					}
				l597:
					{
						position599, tokenIndex599, depth599 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l600
						}
						position++
						goto l599
					l600:
						position, tokenIndex, depth = position599, tokenIndex599, depth599
						if buffer[position] != rune('L') {
							goto l592
						}
						position++
					}
				l599:
					{
						position601, tokenIndex601, depth601 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l602
						}
						position++
						goto l601
					l602:
						position, tokenIndex, depth = position601, tokenIndex601, depth601
						if buffer[position] != rune('S') {
							goto l592
						}
						position++
					}
				l601:
					{
						position603, tokenIndex603, depth603 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l604
						}
						position++
						goto l603
					l604:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
						if buffer[position] != rune('E') {
							goto l592
						}
						position++
					}
				l603:
					depth--
					add(rulePegText, position594)
				}
				if !_rules[ruleAction44]() {
					goto l592
				}
				depth--
				add(ruleFALSE, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 62 Wildcard <- <(<'*'> Action45)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				{
					position607 := position
					depth++
					if buffer[position] != rune('*') {
						goto l605
					}
					position++
					depth--
					add(rulePegText, position607)
				}
				if !_rules[ruleAction45]() {
					goto l605
				}
				depth--
				add(ruleWildcard, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 63 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action46)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l608
					}
					position++
				l611:
					{
						position612, tokenIndex612, depth612 := position, tokenIndex, depth
						{
							position613, tokenIndex613, depth613 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l614
							}
							position++
							if buffer[position] != rune('\'') {
								goto l614
							}
							position++
							goto l613
						l614:
							position, tokenIndex, depth = position613, tokenIndex613, depth613
							{
								position615, tokenIndex615, depth615 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l615
								}
								position++
								goto l612
							l615:
								position, tokenIndex, depth = position615, tokenIndex615, depth615
							}
							if !matchDot() {
								goto l612
							}
						}
					l613:
						goto l611
					l612:
						position, tokenIndex, depth = position612, tokenIndex612, depth612
					}
					if buffer[position] != rune('\'') {
						goto l608
					}
					position++
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction46]() {
					goto l608
				}
				depth--
				add(ruleStringLiteral, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 64 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action47)> */
		func() bool {
			position616, tokenIndex616, depth616 := position, tokenIndex, depth
			{
				position617 := position
				depth++
				{
					position618 := position
					depth++
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('I') {
							goto l616
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
							goto l616
						}
						position++
					}
				l621:
					{
						position623, tokenIndex623, depth623 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l624
						}
						position++
						goto l623
					l624:
						position, tokenIndex, depth = position623, tokenIndex623, depth623
						if buffer[position] != rune('T') {
							goto l616
						}
						position++
					}
				l623:
					{
						position625, tokenIndex625, depth625 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l626
						}
						position++
						goto l625
					l626:
						position, tokenIndex, depth = position625, tokenIndex625, depth625
						if buffer[position] != rune('R') {
							goto l616
						}
						position++
					}
				l625:
					{
						position627, tokenIndex627, depth627 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l628
						}
						position++
						goto l627
					l628:
						position, tokenIndex, depth = position627, tokenIndex627, depth627
						if buffer[position] != rune('E') {
							goto l616
						}
						position++
					}
				l627:
					{
						position629, tokenIndex629, depth629 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l630
						}
						position++
						goto l629
					l630:
						position, tokenIndex, depth = position629, tokenIndex629, depth629
						if buffer[position] != rune('A') {
							goto l616
						}
						position++
					}
				l629:
					{
						position631, tokenIndex631, depth631 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l632
						}
						position++
						goto l631
					l632:
						position, tokenIndex, depth = position631, tokenIndex631, depth631
						if buffer[position] != rune('M') {
							goto l616
						}
						position++
					}
				l631:
					depth--
					add(rulePegText, position618)
				}
				if !_rules[ruleAction47]() {
					goto l616
				}
				depth--
				add(ruleISTREAM, position617)
			}
			return true
		l616:
			position, tokenIndex, depth = position616, tokenIndex616, depth616
			return false
		},
		/* 65 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action48)> */
		func() bool {
			position633, tokenIndex633, depth633 := position, tokenIndex, depth
			{
				position634 := position
				depth++
				{
					position635 := position
					depth++
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
							goto l633
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
							goto l633
						}
						position++
					}
				l638:
					{
						position640, tokenIndex640, depth640 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l641
						}
						position++
						goto l640
					l641:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
						if buffer[position] != rune('T') {
							goto l633
						}
						position++
					}
				l640:
					{
						position642, tokenIndex642, depth642 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l643
						}
						position++
						goto l642
					l643:
						position, tokenIndex, depth = position642, tokenIndex642, depth642
						if buffer[position] != rune('R') {
							goto l633
						}
						position++
					}
				l642:
					{
						position644, tokenIndex644, depth644 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l645
						}
						position++
						goto l644
					l645:
						position, tokenIndex, depth = position644, tokenIndex644, depth644
						if buffer[position] != rune('E') {
							goto l633
						}
						position++
					}
				l644:
					{
						position646, tokenIndex646, depth646 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l647
						}
						position++
						goto l646
					l647:
						position, tokenIndex, depth = position646, tokenIndex646, depth646
						if buffer[position] != rune('A') {
							goto l633
						}
						position++
					}
				l646:
					{
						position648, tokenIndex648, depth648 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l649
						}
						position++
						goto l648
					l649:
						position, tokenIndex, depth = position648, tokenIndex648, depth648
						if buffer[position] != rune('M') {
							goto l633
						}
						position++
					}
				l648:
					depth--
					add(rulePegText, position635)
				}
				if !_rules[ruleAction48]() {
					goto l633
				}
				depth--
				add(ruleDSTREAM, position634)
			}
			return true
		l633:
			position, tokenIndex, depth = position633, tokenIndex633, depth633
			return false
		},
		/* 66 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action49)> */
		func() bool {
			position650, tokenIndex650, depth650 := position, tokenIndex, depth
			{
				position651 := position
				depth++
				{
					position652 := position
					depth++
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
							goto l650
						}
						position++
					}
				l653:
					{
						position655, tokenIndex655, depth655 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l656
						}
						position++
						goto l655
					l656:
						position, tokenIndex, depth = position655, tokenIndex655, depth655
						if buffer[position] != rune('S') {
							goto l650
						}
						position++
					}
				l655:
					{
						position657, tokenIndex657, depth657 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l658
						}
						position++
						goto l657
					l658:
						position, tokenIndex, depth = position657, tokenIndex657, depth657
						if buffer[position] != rune('T') {
							goto l650
						}
						position++
					}
				l657:
					{
						position659, tokenIndex659, depth659 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l660
						}
						position++
						goto l659
					l660:
						position, tokenIndex, depth = position659, tokenIndex659, depth659
						if buffer[position] != rune('R') {
							goto l650
						}
						position++
					}
				l659:
					{
						position661, tokenIndex661, depth661 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l662
						}
						position++
						goto l661
					l662:
						position, tokenIndex, depth = position661, tokenIndex661, depth661
						if buffer[position] != rune('E') {
							goto l650
						}
						position++
					}
				l661:
					{
						position663, tokenIndex663, depth663 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l664
						}
						position++
						goto l663
					l664:
						position, tokenIndex, depth = position663, tokenIndex663, depth663
						if buffer[position] != rune('A') {
							goto l650
						}
						position++
					}
				l663:
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l666
						}
						position++
						goto l665
					l666:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if buffer[position] != rune('M') {
							goto l650
						}
						position++
					}
				l665:
					depth--
					add(rulePegText, position652)
				}
				if !_rules[ruleAction49]() {
					goto l650
				}
				depth--
				add(ruleRSTREAM, position651)
			}
			return true
		l650:
			position, tokenIndex, depth = position650, tokenIndex650, depth650
			return false
		},
		/* 67 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action50)> */
		func() bool {
			position667, tokenIndex667, depth667 := position, tokenIndex, depth
			{
				position668 := position
				depth++
				{
					position669 := position
					depth++
					{
						position670, tokenIndex670, depth670 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l671
						}
						position++
						goto l670
					l671:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
						if buffer[position] != rune('T') {
							goto l667
						}
						position++
					}
				l670:
					{
						position672, tokenIndex672, depth672 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l673
						}
						position++
						goto l672
					l673:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						if buffer[position] != rune('U') {
							goto l667
						}
						position++
					}
				l672:
					{
						position674, tokenIndex674, depth674 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l675
						}
						position++
						goto l674
					l675:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
						if buffer[position] != rune('P') {
							goto l667
						}
						position++
					}
				l674:
					{
						position676, tokenIndex676, depth676 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l677
						}
						position++
						goto l676
					l677:
						position, tokenIndex, depth = position676, tokenIndex676, depth676
						if buffer[position] != rune('L') {
							goto l667
						}
						position++
					}
				l676:
					{
						position678, tokenIndex678, depth678 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l679
						}
						position++
						goto l678
					l679:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
						if buffer[position] != rune('E') {
							goto l667
						}
						position++
					}
				l678:
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l681
						}
						position++
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if buffer[position] != rune('S') {
							goto l667
						}
						position++
					}
				l680:
					depth--
					add(rulePegText, position669)
				}
				if !_rules[ruleAction50]() {
					goto l667
				}
				depth--
				add(ruleTUPLES, position668)
			}
			return true
		l667:
			position, tokenIndex, depth = position667, tokenIndex667, depth667
			return false
		},
		/* 68 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action51)> */
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
						if buffer[position] != rune('s') {
							goto l686
						}
						position++
						goto l685
					l686:
						position, tokenIndex, depth = position685, tokenIndex685, depth685
						if buffer[position] != rune('S') {
							goto l682
						}
						position++
					}
				l685:
					{
						position687, tokenIndex687, depth687 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l688
						}
						position++
						goto l687
					l688:
						position, tokenIndex, depth = position687, tokenIndex687, depth687
						if buffer[position] != rune('E') {
							goto l682
						}
						position++
					}
				l687:
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l690
						}
						position++
						goto l689
					l690:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
						if buffer[position] != rune('C') {
							goto l682
						}
						position++
					}
				l689:
					{
						position691, tokenIndex691, depth691 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l692
						}
						position++
						goto l691
					l692:
						position, tokenIndex, depth = position691, tokenIndex691, depth691
						if buffer[position] != rune('O') {
							goto l682
						}
						position++
					}
				l691:
					{
						position693, tokenIndex693, depth693 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l694
						}
						position++
						goto l693
					l694:
						position, tokenIndex, depth = position693, tokenIndex693, depth693
						if buffer[position] != rune('N') {
							goto l682
						}
						position++
					}
				l693:
					{
						position695, tokenIndex695, depth695 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
						if buffer[position] != rune('D') {
							goto l682
						}
						position++
					}
				l695:
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('S') {
							goto l682
						}
						position++
					}
				l697:
					depth--
					add(rulePegText, position684)
				}
				if !_rules[ruleAction51]() {
					goto l682
				}
				depth--
				add(ruleSECONDS, position683)
			}
			return true
		l682:
			position, tokenIndex, depth = position682, tokenIndex682, depth682
			return false
		},
		/* 69 StreamIdentifier <- <(<ident> Action52)> */
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
				if !_rules[ruleAction52]() {
					goto l699
				}
				depth--
				add(ruleStreamIdentifier, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 70 SourceSinkType <- <(<ident> Action53)> */
		func() bool {
			position702, tokenIndex702, depth702 := position, tokenIndex, depth
			{
				position703 := position
				depth++
				{
					position704 := position
					depth++
					if !_rules[ruleident]() {
						goto l702
					}
					depth--
					add(rulePegText, position704)
				}
				if !_rules[ruleAction53]() {
					goto l702
				}
				depth--
				add(ruleSourceSinkType, position703)
			}
			return true
		l702:
			position, tokenIndex, depth = position702, tokenIndex702, depth702
			return false
		},
		/* 71 SourceSinkParamKey <- <(<ident> Action54)> */
		func() bool {
			position705, tokenIndex705, depth705 := position, tokenIndex, depth
			{
				position706 := position
				depth++
				{
					position707 := position
					depth++
					if !_rules[ruleident]() {
						goto l705
					}
					depth--
					add(rulePegText, position707)
				}
				if !_rules[ruleAction54]() {
					goto l705
				}
				depth--
				add(ruleSourceSinkParamKey, position706)
			}
			return true
		l705:
			position, tokenIndex, depth = position705, tokenIndex705, depth705
			return false
		},
		/* 72 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action55)> */
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
						if buffer[position] != rune('p') {
							goto l712
						}
						position++
						goto l711
					l712:
						position, tokenIndex, depth = position711, tokenIndex711, depth711
						if buffer[position] != rune('P') {
							goto l708
						}
						position++
					}
				l711:
					{
						position713, tokenIndex713, depth713 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l714
						}
						position++
						goto l713
					l714:
						position, tokenIndex, depth = position713, tokenIndex713, depth713
						if buffer[position] != rune('A') {
							goto l708
						}
						position++
					}
				l713:
					{
						position715, tokenIndex715, depth715 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l716
						}
						position++
						goto l715
					l716:
						position, tokenIndex, depth = position715, tokenIndex715, depth715
						if buffer[position] != rune('U') {
							goto l708
						}
						position++
					}
				l715:
					{
						position717, tokenIndex717, depth717 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l718
						}
						position++
						goto l717
					l718:
						position, tokenIndex, depth = position717, tokenIndex717, depth717
						if buffer[position] != rune('S') {
							goto l708
						}
						position++
					}
				l717:
					{
						position719, tokenIndex719, depth719 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l720
						}
						position++
						goto l719
					l720:
						position, tokenIndex, depth = position719, tokenIndex719, depth719
						if buffer[position] != rune('E') {
							goto l708
						}
						position++
					}
				l719:
					{
						position721, tokenIndex721, depth721 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l722
						}
						position++
						goto l721
					l722:
						position, tokenIndex, depth = position721, tokenIndex721, depth721
						if buffer[position] != rune('D') {
							goto l708
						}
						position++
					}
				l721:
					depth--
					add(rulePegText, position710)
				}
				if !_rules[ruleAction55]() {
					goto l708
				}
				depth--
				add(rulePaused, position709)
			}
			return true
		l708:
			position, tokenIndex, depth = position708, tokenIndex708, depth708
			return false
		},
		/* 73 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action56)> */
		func() bool {
			position723, tokenIndex723, depth723 := position, tokenIndex, depth
			{
				position724 := position
				depth++
				{
					position725 := position
					depth++
					{
						position726, tokenIndex726, depth726 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l727
						}
						position++
						goto l726
					l727:
						position, tokenIndex, depth = position726, tokenIndex726, depth726
						if buffer[position] != rune('U') {
							goto l723
						}
						position++
					}
				l726:
					{
						position728, tokenIndex728, depth728 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l729
						}
						position++
						goto l728
					l729:
						position, tokenIndex, depth = position728, tokenIndex728, depth728
						if buffer[position] != rune('N') {
							goto l723
						}
						position++
					}
				l728:
					{
						position730, tokenIndex730, depth730 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l731
						}
						position++
						goto l730
					l731:
						position, tokenIndex, depth = position730, tokenIndex730, depth730
						if buffer[position] != rune('P') {
							goto l723
						}
						position++
					}
				l730:
					{
						position732, tokenIndex732, depth732 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l733
						}
						position++
						goto l732
					l733:
						position, tokenIndex, depth = position732, tokenIndex732, depth732
						if buffer[position] != rune('A') {
							goto l723
						}
						position++
					}
				l732:
					{
						position734, tokenIndex734, depth734 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l735
						}
						position++
						goto l734
					l735:
						position, tokenIndex, depth = position734, tokenIndex734, depth734
						if buffer[position] != rune('U') {
							goto l723
						}
						position++
					}
				l734:
					{
						position736, tokenIndex736, depth736 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l737
						}
						position++
						goto l736
					l737:
						position, tokenIndex, depth = position736, tokenIndex736, depth736
						if buffer[position] != rune('S') {
							goto l723
						}
						position++
					}
				l736:
					{
						position738, tokenIndex738, depth738 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l739
						}
						position++
						goto l738
					l739:
						position, tokenIndex, depth = position738, tokenIndex738, depth738
						if buffer[position] != rune('E') {
							goto l723
						}
						position++
					}
				l738:
					{
						position740, tokenIndex740, depth740 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l741
						}
						position++
						goto l740
					l741:
						position, tokenIndex, depth = position740, tokenIndex740, depth740
						if buffer[position] != rune('D') {
							goto l723
						}
						position++
					}
				l740:
					depth--
					add(rulePegText, position725)
				}
				if !_rules[ruleAction56]() {
					goto l723
				}
				depth--
				add(ruleUnpaused, position724)
			}
			return true
		l723:
			position, tokenIndex, depth = position723, tokenIndex723, depth723
			return false
		},
		/* 74 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action57)> */
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
						if buffer[position] != rune('o') {
							goto l746
						}
						position++
						goto l745
					l746:
						position, tokenIndex, depth = position745, tokenIndex745, depth745
						if buffer[position] != rune('O') {
							goto l742
						}
						position++
					}
				l745:
					{
						position747, tokenIndex747, depth747 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l748
						}
						position++
						goto l747
					l748:
						position, tokenIndex, depth = position747, tokenIndex747, depth747
						if buffer[position] != rune('R') {
							goto l742
						}
						position++
					}
				l747:
					depth--
					add(rulePegText, position744)
				}
				if !_rules[ruleAction57]() {
					goto l742
				}
				depth--
				add(ruleOr, position743)
			}
			return true
		l742:
			position, tokenIndex, depth = position742, tokenIndex742, depth742
			return false
		},
		/* 75 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action58)> */
		func() bool {
			position749, tokenIndex749, depth749 := position, tokenIndex, depth
			{
				position750 := position
				depth++
				{
					position751 := position
					depth++
					{
						position752, tokenIndex752, depth752 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l753
						}
						position++
						goto l752
					l753:
						position, tokenIndex, depth = position752, tokenIndex752, depth752
						if buffer[position] != rune('A') {
							goto l749
						}
						position++
					}
				l752:
					{
						position754, tokenIndex754, depth754 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l755
						}
						position++
						goto l754
					l755:
						position, tokenIndex, depth = position754, tokenIndex754, depth754
						if buffer[position] != rune('N') {
							goto l749
						}
						position++
					}
				l754:
					{
						position756, tokenIndex756, depth756 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l757
						}
						position++
						goto l756
					l757:
						position, tokenIndex, depth = position756, tokenIndex756, depth756
						if buffer[position] != rune('D') {
							goto l749
						}
						position++
					}
				l756:
					depth--
					add(rulePegText, position751)
				}
				if !_rules[ruleAction58]() {
					goto l749
				}
				depth--
				add(ruleAnd, position750)
			}
			return true
		l749:
			position, tokenIndex, depth = position749, tokenIndex749, depth749
			return false
		},
		/* 76 Equal <- <(<'='> Action59)> */
		func() bool {
			position758, tokenIndex758, depth758 := position, tokenIndex, depth
			{
				position759 := position
				depth++
				{
					position760 := position
					depth++
					if buffer[position] != rune('=') {
						goto l758
					}
					position++
					depth--
					add(rulePegText, position760)
				}
				if !_rules[ruleAction59]() {
					goto l758
				}
				depth--
				add(ruleEqual, position759)
			}
			return true
		l758:
			position, tokenIndex, depth = position758, tokenIndex758, depth758
			return false
		},
		/* 77 Less <- <(<'<'> Action60)> */
		func() bool {
			position761, tokenIndex761, depth761 := position, tokenIndex, depth
			{
				position762 := position
				depth++
				{
					position763 := position
					depth++
					if buffer[position] != rune('<') {
						goto l761
					}
					position++
					depth--
					add(rulePegText, position763)
				}
				if !_rules[ruleAction60]() {
					goto l761
				}
				depth--
				add(ruleLess, position762)
			}
			return true
		l761:
			position, tokenIndex, depth = position761, tokenIndex761, depth761
			return false
		},
		/* 78 LessOrEqual <- <(<('<' '=')> Action61)> */
		func() bool {
			position764, tokenIndex764, depth764 := position, tokenIndex, depth
			{
				position765 := position
				depth++
				{
					position766 := position
					depth++
					if buffer[position] != rune('<') {
						goto l764
					}
					position++
					if buffer[position] != rune('=') {
						goto l764
					}
					position++
					depth--
					add(rulePegText, position766)
				}
				if !_rules[ruleAction61]() {
					goto l764
				}
				depth--
				add(ruleLessOrEqual, position765)
			}
			return true
		l764:
			position, tokenIndex, depth = position764, tokenIndex764, depth764
			return false
		},
		/* 79 Greater <- <(<'>'> Action62)> */
		func() bool {
			position767, tokenIndex767, depth767 := position, tokenIndex, depth
			{
				position768 := position
				depth++
				{
					position769 := position
					depth++
					if buffer[position] != rune('>') {
						goto l767
					}
					position++
					depth--
					add(rulePegText, position769)
				}
				if !_rules[ruleAction62]() {
					goto l767
				}
				depth--
				add(ruleGreater, position768)
			}
			return true
		l767:
			position, tokenIndex, depth = position767, tokenIndex767, depth767
			return false
		},
		/* 80 GreaterOrEqual <- <(<('>' '=')> Action63)> */
		func() bool {
			position770, tokenIndex770, depth770 := position, tokenIndex, depth
			{
				position771 := position
				depth++
				{
					position772 := position
					depth++
					if buffer[position] != rune('>') {
						goto l770
					}
					position++
					if buffer[position] != rune('=') {
						goto l770
					}
					position++
					depth--
					add(rulePegText, position772)
				}
				if !_rules[ruleAction63]() {
					goto l770
				}
				depth--
				add(ruleGreaterOrEqual, position771)
			}
			return true
		l770:
			position, tokenIndex, depth = position770, tokenIndex770, depth770
			return false
		},
		/* 81 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action64)> */
		func() bool {
			position773, tokenIndex773, depth773 := position, tokenIndex, depth
			{
				position774 := position
				depth++
				{
					position775 := position
					depth++
					{
						position776, tokenIndex776, depth776 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l777
						}
						position++
						if buffer[position] != rune('=') {
							goto l777
						}
						position++
						goto l776
					l777:
						position, tokenIndex, depth = position776, tokenIndex776, depth776
						if buffer[position] != rune('<') {
							goto l773
						}
						position++
						if buffer[position] != rune('>') {
							goto l773
						}
						position++
					}
				l776:
					depth--
					add(rulePegText, position775)
				}
				if !_rules[ruleAction64]() {
					goto l773
				}
				depth--
				add(ruleNotEqual, position774)
			}
			return true
		l773:
			position, tokenIndex, depth = position773, tokenIndex773, depth773
			return false
		},
		/* 82 Plus <- <(<'+'> Action65)> */
		func() bool {
			position778, tokenIndex778, depth778 := position, tokenIndex, depth
			{
				position779 := position
				depth++
				{
					position780 := position
					depth++
					if buffer[position] != rune('+') {
						goto l778
					}
					position++
					depth--
					add(rulePegText, position780)
				}
				if !_rules[ruleAction65]() {
					goto l778
				}
				depth--
				add(rulePlus, position779)
			}
			return true
		l778:
			position, tokenIndex, depth = position778, tokenIndex778, depth778
			return false
		},
		/* 83 Minus <- <(<'-'> Action66)> */
		func() bool {
			position781, tokenIndex781, depth781 := position, tokenIndex, depth
			{
				position782 := position
				depth++
				{
					position783 := position
					depth++
					if buffer[position] != rune('-') {
						goto l781
					}
					position++
					depth--
					add(rulePegText, position783)
				}
				if !_rules[ruleAction66]() {
					goto l781
				}
				depth--
				add(ruleMinus, position782)
			}
			return true
		l781:
			position, tokenIndex, depth = position781, tokenIndex781, depth781
			return false
		},
		/* 84 Multiply <- <(<'*'> Action67)> */
		func() bool {
			position784, tokenIndex784, depth784 := position, tokenIndex, depth
			{
				position785 := position
				depth++
				{
					position786 := position
					depth++
					if buffer[position] != rune('*') {
						goto l784
					}
					position++
					depth--
					add(rulePegText, position786)
				}
				if !_rules[ruleAction67]() {
					goto l784
				}
				depth--
				add(ruleMultiply, position785)
			}
			return true
		l784:
			position, tokenIndex, depth = position784, tokenIndex784, depth784
			return false
		},
		/* 85 Divide <- <(<'/'> Action68)> */
		func() bool {
			position787, tokenIndex787, depth787 := position, tokenIndex, depth
			{
				position788 := position
				depth++
				{
					position789 := position
					depth++
					if buffer[position] != rune('/') {
						goto l787
					}
					position++
					depth--
					add(rulePegText, position789)
				}
				if !_rules[ruleAction68]() {
					goto l787
				}
				depth--
				add(ruleDivide, position788)
			}
			return true
		l787:
			position, tokenIndex, depth = position787, tokenIndex787, depth787
			return false
		},
		/* 86 Modulo <- <(<'%'> Action69)> */
		func() bool {
			position790, tokenIndex790, depth790 := position, tokenIndex, depth
			{
				position791 := position
				depth++
				{
					position792 := position
					depth++
					if buffer[position] != rune('%') {
						goto l790
					}
					position++
					depth--
					add(rulePegText, position792)
				}
				if !_rules[ruleAction69]() {
					goto l790
				}
				depth--
				add(ruleModulo, position791)
			}
			return true
		l790:
			position, tokenIndex, depth = position790, tokenIndex790, depth790
			return false
		},
		/* 87 Identifier <- <(<ident> Action70)> */
		func() bool {
			position793, tokenIndex793, depth793 := position, tokenIndex, depth
			{
				position794 := position
				depth++
				{
					position795 := position
					depth++
					if !_rules[ruleident]() {
						goto l793
					}
					depth--
					add(rulePegText, position795)
				}
				if !_rules[ruleAction70]() {
					goto l793
				}
				depth--
				add(ruleIdentifier, position794)
			}
			return true
		l793:
			position, tokenIndex, depth = position793, tokenIndex793, depth793
			return false
		},
		/* 88 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position796, tokenIndex796, depth796 := position, tokenIndex, depth
			{
				position797 := position
				depth++
				{
					position798, tokenIndex798, depth798 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l799
					}
					position++
					goto l798
				l799:
					position, tokenIndex, depth = position798, tokenIndex798, depth798
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l796
					}
					position++
				}
			l798:
			l800:
				{
					position801, tokenIndex801, depth801 := position, tokenIndex, depth
					{
						position802, tokenIndex802, depth802 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l803
						}
						position++
						goto l802
					l803:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l804
						}
						position++
						goto l802
					l804:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l805
						}
						position++
						goto l802
					l805:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if buffer[position] != rune('_') {
							goto l801
						}
						position++
					}
				l802:
					goto l800
				l801:
					position, tokenIndex, depth = position801, tokenIndex801, depth801
				}
				depth--
				add(ruleident, position797)
			}
			return true
		l796:
			position, tokenIndex, depth = position796, tokenIndex796, depth796
			return false
		},
		/* 89 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position807 := position
				depth++
			l808:
				{
					position809, tokenIndex809, depth809 := position, tokenIndex, depth
					{
						position810, tokenIndex810, depth810 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l811
						}
						position++
						goto l810
					l811:
						position, tokenIndex, depth = position810, tokenIndex810, depth810
						if buffer[position] != rune('\t') {
							goto l812
						}
						position++
						goto l810
					l812:
						position, tokenIndex, depth = position810, tokenIndex810, depth810
						if buffer[position] != rune('\n') {
							goto l809
						}
						position++
					}
				l810:
					goto l808
				l809:
					position, tokenIndex, depth = position809, tokenIndex809, depth809
				}
				depth--
				add(rulesp, position807)
			}
			return true
		},
		/* 91 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 92 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 93 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 94 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 95 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 96 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 97 Action6 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 98 Action7 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 100 Action8 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 101 Action9 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 102 Action10 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 103 Action11 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 104 Action12 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 105 Action13 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 106 Action14 <- <{
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
		/* 107 Action15 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 108 Action16 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 109 Action17 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 110 Action18 <- <{
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
		/* 111 Action19 <- <{
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
		/* 112 Action20 <- <{
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
		/* 113 Action21 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 114 Action22 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 115 Action23 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 116 Action24 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 117 Action25 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 118 Action26 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 119 Action27 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 120 Action28 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 121 Action29 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 122 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 123 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 124 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 125 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 126 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 127 Action35 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 128 Action36 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 129 Action37 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 130 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 131 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 132 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 133 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 134 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 135 Action43 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 136 Action44 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 137 Action45 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 138 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 139 Action47 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 140 Action48 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 141 Action49 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 142 Action50 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 143 Action51 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 144 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 145 Action53 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 146 Action54 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 147 Action55 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 148 Action56 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 149 Action57 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 150 Action58 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 151 Action59 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 152 Action60 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 153 Action61 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 154 Action62 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 155 Action63 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 156 Action64 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 157 Action65 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 158 Action66 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 159 Action67 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 160 Action68 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 161 Action69 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 162 Action70 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
	}
	p.rules = _rules
}

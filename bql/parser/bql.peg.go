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
	ruleRewindSourceStmt
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
	ruleStreamLike
	ruleUDSFFuncApp
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
	ruleAction8
	rulePegText
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
	ruleAction71
	ruleAction72

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
	"RewindSourceStmt",
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
	"StreamLike",
	"UDSFFuncApp",
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
	"Action8",
	"PegText",
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
	"Action71",
	"Action72",

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
	rules  [168]func() bool
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

			p.AssembleRewindSource()

		case ruleAction9:

			p.AssembleEmitter(begin, end)

		case ruleAction10:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction11:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction12:

			p.AssembleStreamEmitInterval()

		case ruleAction13:

			p.AssembleProjections(begin, end)

		case ruleAction14:

			p.AssembleAlias()

		case ruleAction15:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction16:

			p.AssembleWindowedFrom(begin, end)

		case ruleAction17:

			p.AssembleInterval()

		case ruleAction18:

			p.AssembleInterval()

		case ruleAction19:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction20:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction21:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction22:

			p.EnsureAliasedStreamWindow()

		case ruleAction23:

			p.EnsureAliasedStreamWindow()

		case ruleAction24:

			p.AssembleAliasedStreamWindow()

		case ruleAction25:

			p.AssembleAliasedStreamWindow()

		case ruleAction26:

			p.AssembleStreamWindow()

		case ruleAction27:

			p.AssembleStreamWindow()

		case ruleAction28:

			p.AssembleUDSFFuncApp()

		case ruleAction29:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction30:

			p.AssembleSourceSinkParam()

		case ruleAction31:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction32:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction33:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction34:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleFuncApp()

		case ruleAction38:

			p.AssembleExpressions(begin, end)

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction44:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction45:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction46:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction47:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction49:

			p.PushComponent(begin, end, Istream)

		case ruleAction50:

			p.PushComponent(begin, end, Dstream)

		case ruleAction51:

			p.PushComponent(begin, end, Rstream)

		case ruleAction52:

			p.PushComponent(begin, end, Tuples)

		case ruleAction53:

			p.PushComponent(begin, end, Seconds)

		case ruleAction54:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction55:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction56:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction57:

			p.PushComponent(begin, end, Yes)

		case ruleAction58:

			p.PushComponent(begin, end, No)

		case ruleAction59:

			p.PushComponent(begin, end, Or)

		case ruleAction60:

			p.PushComponent(begin, end, And)

		case ruleAction61:

			p.PushComponent(begin, end, Equal)

		case ruleAction62:

			p.PushComponent(begin, end, Less)

		case ruleAction63:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction64:

			p.PushComponent(begin, end, Greater)

		case ruleAction65:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction66:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction67:

			p.PushComponent(begin, end, Plus)

		case ruleAction68:

			p.PushComponent(begin, end, Minus)

		case ruleAction69:

			p.PushComponent(begin, end, Multiply)

		case ruleAction70:

			p.PushComponent(begin, end, Divide)

		case ruleAction71:

			p.PushComponent(begin, end, Modulo)

		case ruleAction72:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt)> */
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
						goto l17
					}
					goto l9
				l17:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleRewindSourceStmt]() {
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
			position18, tokenIndex18, depth18 := position, tokenIndex, depth
			{
				position19 := position
				depth++
				{
					position20, tokenIndex20, depth20 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l21
					}
					position++
					goto l20
				l21:
					position, tokenIndex, depth = position20, tokenIndex20, depth20
					if buffer[position] != rune('S') {
						goto l18
					}
					position++
				}
			l20:
				{
					position22, tokenIndex22, depth22 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l23
					}
					position++
					goto l22
				l23:
					position, tokenIndex, depth = position22, tokenIndex22, depth22
					if buffer[position] != rune('E') {
						goto l18
					}
					position++
				}
			l22:
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l25
					}
					position++
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if buffer[position] != rune('L') {
						goto l18
					}
					position++
				}
			l24:
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l27
					}
					position++
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					if buffer[position] != rune('E') {
						goto l18
					}
					position++
				}
			l26:
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l29
					}
					position++
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if buffer[position] != rune('C') {
						goto l18
					}
					position++
				}
			l28:
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l31
					}
					position++
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if buffer[position] != rune('T') {
						goto l18
					}
					position++
				}
			l30:
				if !_rules[rulesp]() {
					goto l18
				}
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if !_rules[ruleEmitter]() {
						goto l32
					}
					goto l33
				l32:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
				}
			l33:
				if !_rules[rulesp]() {
					goto l18
				}
				if !_rules[ruleProjections]() {
					goto l18
				}
				if !_rules[rulesp]() {
					goto l18
				}
				if !_rules[ruleDefWindowedFrom]() {
					goto l18
				}
				if !_rules[rulesp]() {
					goto l18
				}
				if !_rules[ruleFilter]() {
					goto l18
				}
				if !_rules[rulesp]() {
					goto l18
				}
				if !_rules[ruleGrouping]() {
					goto l18
				}
				if !_rules[rulesp]() {
					goto l18
				}
				if !_rules[ruleHaving]() {
					goto l18
				}
				if !_rules[rulesp]() {
					goto l18
				}
				if !_rules[ruleAction0]() {
					goto l18
				}
				depth--
				add(ruleSelectStmt, position19)
			}
			return true
		l18:
			position, tokenIndex, depth = position18, tokenIndex18, depth18
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position34, tokenIndex34, depth34 := position, tokenIndex, depth
			{
				position35 := position
				depth++
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l37
					}
					position++
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if buffer[position] != rune('C') {
						goto l34
					}
					position++
				}
			l36:
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l39
					}
					position++
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if buffer[position] != rune('R') {
						goto l34
					}
					position++
				}
			l38:
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('E') {
						goto l34
					}
					position++
				}
			l40:
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('A') {
						goto l34
					}
					position++
				}
			l42:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l45
					}
					position++
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if buffer[position] != rune('T') {
						goto l34
					}
					position++
				}
			l44:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('E') {
						goto l34
					}
					position++
				}
			l46:
				if !_rules[rulesp]() {
					goto l34
				}
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('S') {
						goto l34
					}
					position++
				}
			l48:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('T') {
						goto l34
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('R') {
						goto l34
					}
					position++
				}
			l52:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('E') {
						goto l34
					}
					position++
				}
			l54:
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
						goto l34
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('M') {
						goto l34
					}
					position++
				}
			l58:
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l34
				}
				if !_rules[rulesp]() {
					goto l34
				}
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l61
					}
					position++
					goto l60
				l61:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
					if buffer[position] != rune('A') {
						goto l34
					}
					position++
				}
			l60:
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
						goto l34
					}
					position++
				}
			l62:
				if !_rules[rulesp]() {
					goto l34
				}
				{
					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l65
					}
					position++
					goto l64
				l65:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
					if buffer[position] != rune('S') {
						goto l34
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
						goto l34
					}
					position++
				}
			l66:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
					if buffer[position] != rune('L') {
						goto l34
					}
					position++
				}
			l68:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l71
					}
					position++
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if buffer[position] != rune('E') {
						goto l34
					}
					position++
				}
			l70:
				{
					position72, tokenIndex72, depth72 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l73
					}
					position++
					goto l72
				l73:
					position, tokenIndex, depth = position72, tokenIndex72, depth72
					if buffer[position] != rune('C') {
						goto l34
					}
					position++
				}
			l72:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if buffer[position] != rune('T') {
						goto l34
					}
					position++
				}
			l74:
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleEmitter]() {
					goto l34
				}
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleProjections]() {
					goto l34
				}
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleWindowedFrom]() {
					goto l34
				}
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleFilter]() {
					goto l34
				}
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleGrouping]() {
					goto l34
				}
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleHaving]() {
					goto l34
				}
				if !_rules[rulesp]() {
					goto l34
				}
				if !_rules[ruleAction1]() {
					goto l34
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position35)
			}
			return true
		l34:
			position, tokenIndex, depth = position34, tokenIndex34, depth34
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position76, tokenIndex76, depth76 := position, tokenIndex, depth
			{
				position77 := position
				depth++
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l79
					}
					position++
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('C') {
						goto l76
					}
					position++
				}
			l78:
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('R') {
						goto l76
					}
					position++
				}
			l80:
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l83
					}
					position++
					goto l82
				l83:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
					if buffer[position] != rune('E') {
						goto l76
					}
					position++
				}
			l82:
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
					if buffer[position] != rune('A') {
						goto l76
					}
					position++
				}
			l84:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if buffer[position] != rune('T') {
						goto l76
					}
					position++
				}
			l86:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l89
					}
					position++
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if buffer[position] != rune('E') {
						goto l76
					}
					position++
				}
			l88:
				if !_rules[rulesp]() {
					goto l76
				}
				if !_rules[rulePausedOpt]() {
					goto l76
				}
				if !_rules[rulesp]() {
					goto l76
				}
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('S') {
						goto l76
					}
					position++
				}
			l90:
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('O') {
						goto l76
					}
					position++
				}
			l92:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('U') {
						goto l76
					}
					position++
				}
			l94:
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('R') {
						goto l76
					}
					position++
				}
			l96:
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('C') {
						goto l76
					}
					position++
				}
			l98:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('E') {
						goto l76
					}
					position++
				}
			l100:
				if !_rules[rulesp]() {
					goto l76
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l76
				}
				if !_rules[rulesp]() {
					goto l76
				}
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('T') {
						goto l76
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('Y') {
						goto l76
					}
					position++
				}
			l104:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('P') {
						goto l76
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('E') {
						goto l76
					}
					position++
				}
			l108:
				if !_rules[rulesp]() {
					goto l76
				}
				if !_rules[ruleSourceSinkType]() {
					goto l76
				}
				if !_rules[rulesp]() {
					goto l76
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l76
				}
				if !_rules[ruleAction2]() {
					goto l76
				}
				depth--
				add(ruleCreateSourceStmt, position77)
			}
			return true
		l76:
			position, tokenIndex, depth = position76, tokenIndex76, depth76
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position110, tokenIndex110, depth110 := position, tokenIndex, depth
			{
				position111 := position
				depth++
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('C') {
						goto l110
					}
					position++
				}
			l112:
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('R') {
						goto l110
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('E') {
						goto l110
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('A') {
						goto l110
					}
					position++
				}
			l118:
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
					if buffer[position] != rune('T') {
						goto l110
					}
					position++
				}
			l120:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('E') {
						goto l110
					}
					position++
				}
			l122:
				if !_rules[rulesp]() {
					goto l110
				}
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('S') {
						goto l110
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('I') {
						goto l110
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('N') {
						goto l110
					}
					position++
				}
			l128:
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('K') {
						goto l110
					}
					position++
				}
			l130:
				if !_rules[rulesp]() {
					goto l110
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l110
				}
				if !_rules[rulesp]() {
					goto l110
				}
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('T') {
						goto l110
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('Y') {
						goto l110
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('P') {
						goto l110
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('E') {
						goto l110
					}
					position++
				}
			l138:
				if !_rules[rulesp]() {
					goto l110
				}
				if !_rules[ruleSourceSinkType]() {
					goto l110
				}
				if !_rules[rulesp]() {
					goto l110
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l110
				}
				if !_rules[ruleAction3]() {
					goto l110
				}
				depth--
				add(ruleCreateSinkStmt, position111)
			}
			return true
		l110:
			position, tokenIndex, depth = position110, tokenIndex110, depth110
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position140, tokenIndex140, depth140 := position, tokenIndex, depth
			{
				position141 := position
				depth++
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('C') {
						goto l140
					}
					position++
				}
			l142:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('R') {
						goto l140
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('E') {
						goto l140
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('A') {
						goto l140
					}
					position++
				}
			l148:
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('T') {
						goto l140
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('E') {
						goto l140
					}
					position++
				}
			l152:
				if !_rules[rulesp]() {
					goto l140
				}
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('S') {
						goto l140
					}
					position++
				}
			l154:
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('T') {
						goto l140
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
						goto l140
					}
					position++
				}
			l158:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('T') {
						goto l140
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('E') {
						goto l140
					}
					position++
				}
			l162:
				if !_rules[rulesp]() {
					goto l140
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l140
				}
				if !_rules[rulesp]() {
					goto l140
				}
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('T') {
						goto l140
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('Y') {
						goto l140
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('P') {
						goto l140
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('E') {
						goto l140
					}
					position++
				}
			l170:
				if !_rules[rulesp]() {
					goto l140
				}
				if !_rules[ruleSourceSinkType]() {
					goto l140
				}
				if !_rules[rulesp]() {
					goto l140
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l140
				}
				if !_rules[ruleAction4]() {
					goto l140
				}
				depth--
				add(ruleCreateStateStmt, position141)
			}
			return true
		l140:
			position, tokenIndex, depth = position140, tokenIndex140, depth140
			return false
		},
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action5)> */
		func() bool {
			position172, tokenIndex172, depth172 := position, tokenIndex, depth
			{
				position173 := position
				depth++
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('I') {
						goto l172
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('N') {
						goto l172
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('S') {
						goto l172
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
						goto l172
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('R') {
						goto l172
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('T') {
						goto l172
					}
					position++
				}
			l184:
				if !_rules[rulesp]() {
					goto l172
				}
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('I') {
						goto l172
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('N') {
						goto l172
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('T') {
						goto l172
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('O') {
						goto l172
					}
					position++
				}
			l192:
				if !_rules[rulesp]() {
					goto l172
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l172
				}
				if !_rules[rulesp]() {
					goto l172
				}
				if !_rules[ruleSelectStmt]() {
					goto l172
				}
				if !_rules[ruleAction5]() {
					goto l172
				}
				depth--
				add(ruleInsertIntoSelectStmt, position173)
			}
			return true
		l172:
			position, tokenIndex, depth = position172, tokenIndex172, depth172
			return false
		},
		/* 8 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action6)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('P') {
						goto l194
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('A') {
						goto l194
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('U') {
						goto l194
					}
					position++
				}
			l200:
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('S') {
						goto l194
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
						goto l194
					}
					position++
				}
			l204:
				if !_rules[rulesp]() {
					goto l194
				}
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('S') {
						goto l194
					}
					position++
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('O') {
						goto l194
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('U') {
						goto l194
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
						goto l194
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('C') {
						goto l194
					}
					position++
				}
			l214:
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('E') {
						goto l194
					}
					position++
				}
			l216:
				if !_rules[rulesp]() {
					goto l194
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l194
				}
				if !_rules[ruleAction6]() {
					goto l194
				}
				depth--
				add(rulePauseSourceStmt, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 9 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action7)> */
		func() bool {
			position218, tokenIndex218, depth218 := position, tokenIndex, depth
			{
				position219 := position
				depth++
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('R') {
						goto l218
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('E') {
						goto l218
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('S') {
						goto l218
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('U') {
						goto l218
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('M') {
						goto l218
					}
					position++
				}
			l228:
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
					if buffer[position] != rune('E') {
						goto l218
					}
					position++
				}
			l230:
				if !_rules[rulesp]() {
					goto l218
				}
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('S') {
						goto l218
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('O') {
						goto l218
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('U') {
						goto l218
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
						goto l218
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('C') {
						goto l218
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('E') {
						goto l218
					}
					position++
				}
			l242:
				if !_rules[rulesp]() {
					goto l218
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l218
				}
				if !_rules[ruleAction7]() {
					goto l218
				}
				depth--
				add(ruleResumeSourceStmt, position219)
			}
			return true
		l218:
			position, tokenIndex, depth = position218, tokenIndex218, depth218
			return false
		},
		/* 10 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action8)> */
		func() bool {
			position244, tokenIndex244, depth244 := position, tokenIndex, depth
			{
				position245 := position
				depth++
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('R') {
						goto l244
					}
					position++
				}
			l246:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('E') {
						goto l244
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('W') {
						goto l244
					}
					position++
				}
			l250:
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if buffer[position] != rune('I') {
						goto l244
					}
					position++
				}
			l252:
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if buffer[position] != rune('N') {
						goto l244
					}
					position++
				}
			l254:
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l257
					}
					position++
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if buffer[position] != rune('D') {
						goto l244
					}
					position++
				}
			l256:
				if !_rules[rulesp]() {
					goto l244
				}
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if buffer[position] != rune('S') {
						goto l244
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('O') {
						goto l244
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('U') {
						goto l244
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
						goto l244
					}
					position++
				}
			l264:
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('C') {
						goto l244
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('E') {
						goto l244
					}
					position++
				}
			l268:
				if !_rules[rulesp]() {
					goto l244
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l244
				}
				if !_rules[ruleAction8]() {
					goto l244
				}
				depth--
				add(ruleRewindSourceStmt, position245)
			}
			return true
		l244:
			position, tokenIndex, depth = position244, tokenIndex244, depth244
			return false
		},
		/* 11 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action9)> */
		func() bool {
			position270, tokenIndex270, depth270 := position, tokenIndex, depth
			{
				position271 := position
				depth++
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l273
					}
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if !_rules[ruleDSTREAM]() {
						goto l274
					}
					goto l272
				l274:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if !_rules[ruleRSTREAM]() {
						goto l270
					}
				}
			l272:
				{
					position275 := position
					depth++
					{
						position276, tokenIndex276, depth276 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l276
						}
						if buffer[position] != rune('[') {
							goto l276
						}
						position++
						if !_rules[rulesp]() {
							goto l276
						}
						{
							position278, tokenIndex278, depth278 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l279
							}
							position++
							goto l278
						l279:
							position, tokenIndex, depth = position278, tokenIndex278, depth278
							if buffer[position] != rune('E') {
								goto l276
							}
							position++
						}
					l278:
						{
							position280, tokenIndex280, depth280 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l281
							}
							position++
							goto l280
						l281:
							position, tokenIndex, depth = position280, tokenIndex280, depth280
							if buffer[position] != rune('V') {
								goto l276
							}
							position++
						}
					l280:
						{
							position282, tokenIndex282, depth282 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l283
							}
							position++
							goto l282
						l283:
							position, tokenIndex, depth = position282, tokenIndex282, depth282
							if buffer[position] != rune('E') {
								goto l276
							}
							position++
						}
					l282:
						{
							position284, tokenIndex284, depth284 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l285
							}
							position++
							goto l284
						l285:
							position, tokenIndex, depth = position284, tokenIndex284, depth284
							if buffer[position] != rune('R') {
								goto l276
							}
							position++
						}
					l284:
						{
							position286, tokenIndex286, depth286 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l287
							}
							position++
							goto l286
						l287:
							position, tokenIndex, depth = position286, tokenIndex286, depth286
							if buffer[position] != rune('Y') {
								goto l276
							}
							position++
						}
					l286:
						if !_rules[rulesp]() {
							goto l276
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l276
						}
						if !_rules[rulesp]() {
							goto l276
						}
						if buffer[position] != rune(']') {
							goto l276
						}
						position++
						goto l277
					l276:
						position, tokenIndex, depth = position276, tokenIndex276, depth276
					}
				l277:
					depth--
					add(rulePegText, position275)
				}
				if !_rules[ruleAction9]() {
					goto l270
				}
				depth--
				add(ruleEmitter, position271)
			}
			return true
		l270:
			position, tokenIndex, depth = position270, tokenIndex270, depth270
			return false
		},
		/* 12 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position288, tokenIndex288, depth288 := position, tokenIndex, depth
			{
				position289 := position
				depth++
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l291
					}
				l292:
					{
						position293, tokenIndex293, depth293 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l293
						}
						if buffer[position] != rune(',') {
							goto l293
						}
						position++
						if !_rules[rulesp]() {
							goto l293
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l293
						}
						goto l292
					l293:
						position, tokenIndex, depth = position293, tokenIndex293, depth293
					}
					goto l290
				l291:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
					if !_rules[ruleTimeEmitterInterval]() {
						goto l294
					}
					goto l290
				l294:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
					if !_rules[ruleTupleEmitterInterval]() {
						goto l288
					}
				}
			l290:
				depth--
				add(ruleEmitterIntervals, position289)
			}
			return true
		l288:
			position, tokenIndex, depth = position288, tokenIndex288, depth288
			return false
		},
		/* 13 TimeEmitterInterval <- <(<TimeInterval> Action10)> */
		func() bool {
			position295, tokenIndex295, depth295 := position, tokenIndex, depth
			{
				position296 := position
				depth++
				{
					position297 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l295
					}
					depth--
					add(rulePegText, position297)
				}
				if !_rules[ruleAction10]() {
					goto l295
				}
				depth--
				add(ruleTimeEmitterInterval, position296)
			}
			return true
		l295:
			position, tokenIndex, depth = position295, tokenIndex295, depth295
			return false
		},
		/* 14 TupleEmitterInterval <- <(<TuplesInterval> Action11)> */
		func() bool {
			position298, tokenIndex298, depth298 := position, tokenIndex, depth
			{
				position299 := position
				depth++
				{
					position300 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l298
					}
					depth--
					add(rulePegText, position300)
				}
				if !_rules[ruleAction11]() {
					goto l298
				}
				depth--
				add(ruleTupleEmitterInterval, position299)
			}
			return true
		l298:
			position, tokenIndex, depth = position298, tokenIndex298, depth298
			return false
		},
		/* 15 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action12)> */
		func() bool {
			position301, tokenIndex301, depth301 := position, tokenIndex, depth
			{
				position302 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l301
				}
				if !_rules[rulesp]() {
					goto l301
				}
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('I') {
						goto l301
					}
					position++
				}
			l303:
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l306
					}
					position++
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if buffer[position] != rune('N') {
						goto l301
					}
					position++
				}
			l305:
				if !_rules[rulesp]() {
					goto l301
				}
				if !_rules[ruleStream]() {
					goto l301
				}
				if !_rules[ruleAction12]() {
					goto l301
				}
				depth--
				add(ruleTupleEmitterFromInterval, position302)
			}
			return true
		l301:
			position, tokenIndex, depth = position301, tokenIndex301, depth301
			return false
		},
		/* 16 Projections <- <(<(Projection sp (',' sp Projection)*)> Action13)> */
		func() bool {
			position307, tokenIndex307, depth307 := position, tokenIndex, depth
			{
				position308 := position
				depth++
				{
					position309 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l307
					}
					if !_rules[rulesp]() {
						goto l307
					}
				l310:
					{
						position311, tokenIndex311, depth311 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l311
						}
						position++
						if !_rules[rulesp]() {
							goto l311
						}
						if !_rules[ruleProjection]() {
							goto l311
						}
						goto l310
					l311:
						position, tokenIndex, depth = position311, tokenIndex311, depth311
					}
					depth--
					add(rulePegText, position309)
				}
				if !_rules[ruleAction13]() {
					goto l307
				}
				depth--
				add(ruleProjections, position308)
			}
			return true
		l307:
			position, tokenIndex, depth = position307, tokenIndex307, depth307
			return false
		},
		/* 17 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position312, tokenIndex312, depth312 := position, tokenIndex, depth
			{
				position313 := position
				depth++
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l315
					}
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if !_rules[ruleExpression]() {
						goto l316
					}
					goto l314
				l316:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if !_rules[ruleWildcard]() {
						goto l312
					}
				}
			l314:
				depth--
				add(ruleProjection, position313)
			}
			return true
		l312:
			position, tokenIndex, depth = position312, tokenIndex312, depth312
			return false
		},
		/* 18 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action14)> */
		func() bool {
			position317, tokenIndex317, depth317 := position, tokenIndex, depth
			{
				position318 := position
				depth++
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l320
					}
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if !_rules[ruleWildcard]() {
						goto l317
					}
				}
			l319:
				if !_rules[rulesp]() {
					goto l317
				}
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l322
					}
					position++
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if buffer[position] != rune('A') {
						goto l317
					}
					position++
				}
			l321:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l324
					}
					position++
					goto l323
				l324:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					if buffer[position] != rune('S') {
						goto l317
					}
					position++
				}
			l323:
				if !_rules[rulesp]() {
					goto l317
				}
				if !_rules[ruleIdentifier]() {
					goto l317
				}
				if !_rules[ruleAction14]() {
					goto l317
				}
				depth--
				add(ruleAliasExpression, position318)
			}
			return true
		l317:
			position, tokenIndex, depth = position317, tokenIndex317, depth317
			return false
		},
		/* 19 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action15)> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				{
					position327 := position
					depth++
					{
						position328, tokenIndex328, depth328 := position, tokenIndex, depth
						{
							position330, tokenIndex330, depth330 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l331
							}
							position++
							goto l330
						l331:
							position, tokenIndex, depth = position330, tokenIndex330, depth330
							if buffer[position] != rune('F') {
								goto l328
							}
							position++
						}
					l330:
						{
							position332, tokenIndex332, depth332 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l333
							}
							position++
							goto l332
						l333:
							position, tokenIndex, depth = position332, tokenIndex332, depth332
							if buffer[position] != rune('R') {
								goto l328
							}
							position++
						}
					l332:
						{
							position334, tokenIndex334, depth334 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l335
							}
							position++
							goto l334
						l335:
							position, tokenIndex, depth = position334, tokenIndex334, depth334
							if buffer[position] != rune('O') {
								goto l328
							}
							position++
						}
					l334:
						{
							position336, tokenIndex336, depth336 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l337
							}
							position++
							goto l336
						l337:
							position, tokenIndex, depth = position336, tokenIndex336, depth336
							if buffer[position] != rune('M') {
								goto l328
							}
							position++
						}
					l336:
						if !_rules[rulesp]() {
							goto l328
						}
						if !_rules[ruleRelations]() {
							goto l328
						}
						if !_rules[rulesp]() {
							goto l328
						}
						goto l329
					l328:
						position, tokenIndex, depth = position328, tokenIndex328, depth328
					}
				l329:
					depth--
					add(rulePegText, position327)
				}
				if !_rules[ruleAction15]() {
					goto l325
				}
				depth--
				add(ruleWindowedFrom, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
			return false
		},
		/* 20 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action16)> */
		func() bool {
			position338, tokenIndex338, depth338 := position, tokenIndex, depth
			{
				position339 := position
				depth++
				{
					position340 := position
					depth++
					{
						position341, tokenIndex341, depth341 := position, tokenIndex, depth
						{
							position343, tokenIndex343, depth343 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l344
							}
							position++
							goto l343
						l344:
							position, tokenIndex, depth = position343, tokenIndex343, depth343
							if buffer[position] != rune('F') {
								goto l341
							}
							position++
						}
					l343:
						{
							position345, tokenIndex345, depth345 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l346
							}
							position++
							goto l345
						l346:
							position, tokenIndex, depth = position345, tokenIndex345, depth345
							if buffer[position] != rune('R') {
								goto l341
							}
							position++
						}
					l345:
						{
							position347, tokenIndex347, depth347 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l348
							}
							position++
							goto l347
						l348:
							position, tokenIndex, depth = position347, tokenIndex347, depth347
							if buffer[position] != rune('O') {
								goto l341
							}
							position++
						}
					l347:
						{
							position349, tokenIndex349, depth349 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l350
							}
							position++
							goto l349
						l350:
							position, tokenIndex, depth = position349, tokenIndex349, depth349
							if buffer[position] != rune('M') {
								goto l341
							}
							position++
						}
					l349:
						if !_rules[rulesp]() {
							goto l341
						}
						if !_rules[ruleDefRelations]() {
							goto l341
						}
						if !_rules[rulesp]() {
							goto l341
						}
						goto l342
					l341:
						position, tokenIndex, depth = position341, tokenIndex341, depth341
					}
				l342:
					depth--
					add(rulePegText, position340)
				}
				if !_rules[ruleAction16]() {
					goto l338
				}
				depth--
				add(ruleDefWindowedFrom, position339)
			}
			return true
		l338:
			position, tokenIndex, depth = position338, tokenIndex338, depth338
			return false
		},
		/* 21 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position351, tokenIndex351, depth351 := position, tokenIndex, depth
			{
				position352 := position
				depth++
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l354
					}
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if !_rules[ruleTuplesInterval]() {
						goto l351
					}
				}
			l353:
				depth--
				add(ruleInterval, position352)
			}
			return true
		l351:
			position, tokenIndex, depth = position351, tokenIndex351, depth351
			return false
		},
		/* 22 TimeInterval <- <(NumericLiteral sp SECONDS Action17)> */
		func() bool {
			position355, tokenIndex355, depth355 := position, tokenIndex, depth
			{
				position356 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l355
				}
				if !_rules[rulesp]() {
					goto l355
				}
				if !_rules[ruleSECONDS]() {
					goto l355
				}
				if !_rules[ruleAction17]() {
					goto l355
				}
				depth--
				add(ruleTimeInterval, position356)
			}
			return true
		l355:
			position, tokenIndex, depth = position355, tokenIndex355, depth355
			return false
		},
		/* 23 TuplesInterval <- <(NumericLiteral sp TUPLES Action18)> */
		func() bool {
			position357, tokenIndex357, depth357 := position, tokenIndex, depth
			{
				position358 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l357
				}
				if !_rules[rulesp]() {
					goto l357
				}
				if !_rules[ruleTUPLES]() {
					goto l357
				}
				if !_rules[ruleAction18]() {
					goto l357
				}
				depth--
				add(ruleTuplesInterval, position358)
			}
			return true
		l357:
			position, tokenIndex, depth = position357, tokenIndex357, depth357
			return false
		},
		/* 24 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position359, tokenIndex359, depth359 := position, tokenIndex, depth
			{
				position360 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l359
				}
				if !_rules[rulesp]() {
					goto l359
				}
			l361:
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l362
					}
					position++
					if !_rules[rulesp]() {
						goto l362
					}
					if !_rules[ruleRelationLike]() {
						goto l362
					}
					goto l361
				l362:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
				}
				depth--
				add(ruleRelations, position360)
			}
			return true
		l359:
			position, tokenIndex, depth = position359, tokenIndex359, depth359
			return false
		},
		/* 25 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position363, tokenIndex363, depth363 := position, tokenIndex, depth
			{
				position364 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l363
				}
				if !_rules[rulesp]() {
					goto l363
				}
			l365:
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l366
					}
					position++
					if !_rules[rulesp]() {
						goto l366
					}
					if !_rules[ruleDefRelationLike]() {
						goto l366
					}
					goto l365
				l366:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
				}
				depth--
				add(ruleDefRelations, position364)
			}
			return true
		l363:
			position, tokenIndex, depth = position363, tokenIndex363, depth363
			return false
		},
		/* 26 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action19)> */
		func() bool {
			position367, tokenIndex367, depth367 := position, tokenIndex, depth
			{
				position368 := position
				depth++
				{
					position369 := position
					depth++
					{
						position370, tokenIndex370, depth370 := position, tokenIndex, depth
						{
							position372, tokenIndex372, depth372 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l373
							}
							position++
							goto l372
						l373:
							position, tokenIndex, depth = position372, tokenIndex372, depth372
							if buffer[position] != rune('W') {
								goto l370
							}
							position++
						}
					l372:
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l375
							}
							position++
							goto l374
						l375:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
							if buffer[position] != rune('H') {
								goto l370
							}
							position++
						}
					l374:
						{
							position376, tokenIndex376, depth376 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l377
							}
							position++
							goto l376
						l377:
							position, tokenIndex, depth = position376, tokenIndex376, depth376
							if buffer[position] != rune('E') {
								goto l370
							}
							position++
						}
					l376:
						{
							position378, tokenIndex378, depth378 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l379
							}
							position++
							goto l378
						l379:
							position, tokenIndex, depth = position378, tokenIndex378, depth378
							if buffer[position] != rune('R') {
								goto l370
							}
							position++
						}
					l378:
						{
							position380, tokenIndex380, depth380 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l381
							}
							position++
							goto l380
						l381:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
							if buffer[position] != rune('E') {
								goto l370
							}
							position++
						}
					l380:
						if !_rules[rulesp]() {
							goto l370
						}
						if !_rules[ruleExpression]() {
							goto l370
						}
						goto l371
					l370:
						position, tokenIndex, depth = position370, tokenIndex370, depth370
					}
				l371:
					depth--
					add(rulePegText, position369)
				}
				if !_rules[ruleAction19]() {
					goto l367
				}
				depth--
				add(ruleFilter, position368)
			}
			return true
		l367:
			position, tokenIndex, depth = position367, tokenIndex367, depth367
			return false
		},
		/* 27 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action20)> */
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
							if buffer[position] != rune('g') {
								goto l388
							}
							position++
							goto l387
						l388:
							position, tokenIndex, depth = position387, tokenIndex387, depth387
							if buffer[position] != rune('G') {
								goto l385
							}
							position++
						}
					l387:
						{
							position389, tokenIndex389, depth389 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l390
							}
							position++
							goto l389
						l390:
							position, tokenIndex, depth = position389, tokenIndex389, depth389
							if buffer[position] != rune('R') {
								goto l385
							}
							position++
						}
					l389:
						{
							position391, tokenIndex391, depth391 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l392
							}
							position++
							goto l391
						l392:
							position, tokenIndex, depth = position391, tokenIndex391, depth391
							if buffer[position] != rune('O') {
								goto l385
							}
							position++
						}
					l391:
						{
							position393, tokenIndex393, depth393 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l394
							}
							position++
							goto l393
						l394:
							position, tokenIndex, depth = position393, tokenIndex393, depth393
							if buffer[position] != rune('U') {
								goto l385
							}
							position++
						}
					l393:
						{
							position395, tokenIndex395, depth395 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l396
							}
							position++
							goto l395
						l396:
							position, tokenIndex, depth = position395, tokenIndex395, depth395
							if buffer[position] != rune('P') {
								goto l385
							}
							position++
						}
					l395:
						if !_rules[rulesp]() {
							goto l385
						}
						{
							position397, tokenIndex397, depth397 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l398
							}
							position++
							goto l397
						l398:
							position, tokenIndex, depth = position397, tokenIndex397, depth397
							if buffer[position] != rune('B') {
								goto l385
							}
							position++
						}
					l397:
						{
							position399, tokenIndex399, depth399 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l400
							}
							position++
							goto l399
						l400:
							position, tokenIndex, depth = position399, tokenIndex399, depth399
							if buffer[position] != rune('Y') {
								goto l385
							}
							position++
						}
					l399:
						if !_rules[rulesp]() {
							goto l385
						}
						if !_rules[ruleGroupList]() {
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
				if !_rules[ruleAction20]() {
					goto l382
				}
				depth--
				add(ruleGrouping, position383)
			}
			return true
		l382:
			position, tokenIndex, depth = position382, tokenIndex382, depth382
			return false
		},
		/* 28 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position401, tokenIndex401, depth401 := position, tokenIndex, depth
			{
				position402 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l401
				}
				if !_rules[rulesp]() {
					goto l401
				}
			l403:
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l404
					}
					position++
					if !_rules[rulesp]() {
						goto l404
					}
					if !_rules[ruleExpression]() {
						goto l404
					}
					goto l403
				l404:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
				}
				depth--
				add(ruleGroupList, position402)
			}
			return true
		l401:
			position, tokenIndex, depth = position401, tokenIndex401, depth401
			return false
		},
		/* 29 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action21)> */
		func() bool {
			position405, tokenIndex405, depth405 := position, tokenIndex, depth
			{
				position406 := position
				depth++
				{
					position407 := position
					depth++
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						{
							position410, tokenIndex410, depth410 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l411
							}
							position++
							goto l410
						l411:
							position, tokenIndex, depth = position410, tokenIndex410, depth410
							if buffer[position] != rune('H') {
								goto l408
							}
							position++
						}
					l410:
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
								goto l408
							}
							position++
						}
					l412:
						{
							position414, tokenIndex414, depth414 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l415
							}
							position++
							goto l414
						l415:
							position, tokenIndex, depth = position414, tokenIndex414, depth414
							if buffer[position] != rune('V') {
								goto l408
							}
							position++
						}
					l414:
						{
							position416, tokenIndex416, depth416 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l417
							}
							position++
							goto l416
						l417:
							position, tokenIndex, depth = position416, tokenIndex416, depth416
							if buffer[position] != rune('I') {
								goto l408
							}
							position++
						}
					l416:
						{
							position418, tokenIndex418, depth418 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l419
							}
							position++
							goto l418
						l419:
							position, tokenIndex, depth = position418, tokenIndex418, depth418
							if buffer[position] != rune('N') {
								goto l408
							}
							position++
						}
					l418:
						{
							position420, tokenIndex420, depth420 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l421
							}
							position++
							goto l420
						l421:
							position, tokenIndex, depth = position420, tokenIndex420, depth420
							if buffer[position] != rune('G') {
								goto l408
							}
							position++
						}
					l420:
						if !_rules[rulesp]() {
							goto l408
						}
						if !_rules[ruleExpression]() {
							goto l408
						}
						goto l409
					l408:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
					}
				l409:
					depth--
					add(rulePegText, position407)
				}
				if !_rules[ruleAction21]() {
					goto l405
				}
				depth--
				add(ruleHaving, position406)
			}
			return true
		l405:
			position, tokenIndex, depth = position405, tokenIndex405, depth405
			return false
		},
		/* 30 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action22))> */
		func() bool {
			position422, tokenIndex422, depth422 := position, tokenIndex, depth
			{
				position423 := position
				depth++
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l425
					}
					goto l424
				l425:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
					if !_rules[ruleStreamWindow]() {
						goto l422
					}
					if !_rules[ruleAction22]() {
						goto l422
					}
				}
			l424:
				depth--
				add(ruleRelationLike, position423)
			}
			return true
		l422:
			position, tokenIndex, depth = position422, tokenIndex422, depth422
			return false
		},
		/* 31 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action23))> */
		func() bool {
			position426, tokenIndex426, depth426 := position, tokenIndex, depth
			{
				position427 := position
				depth++
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l429
					}
					goto l428
				l429:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
					if !_rules[ruleDefStreamWindow]() {
						goto l426
					}
					if !_rules[ruleAction23]() {
						goto l426
					}
				}
			l428:
				depth--
				add(ruleDefRelationLike, position427)
			}
			return true
		l426:
			position, tokenIndex, depth = position426, tokenIndex426, depth426
			return false
		},
		/* 32 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action24)> */
		func() bool {
			position430, tokenIndex430, depth430 := position, tokenIndex, depth
			{
				position431 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l430
				}
				if !_rules[rulesp]() {
					goto l430
				}
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('A') {
						goto l430
					}
					position++
				}
			l432:
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l435
					}
					position++
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if buffer[position] != rune('S') {
						goto l430
					}
					position++
				}
			l434:
				if !_rules[rulesp]() {
					goto l430
				}
				if !_rules[ruleIdentifier]() {
					goto l430
				}
				if !_rules[ruleAction24]() {
					goto l430
				}
				depth--
				add(ruleAliasedStreamWindow, position431)
			}
			return true
		l430:
			position, tokenIndex, depth = position430, tokenIndex430, depth430
			return false
		},
		/* 33 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action25)> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l436
				}
				if !_rules[rulesp]() {
					goto l436
				}
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if buffer[position] != rune('A') {
						goto l436
					}
					position++
				}
			l438:
				{
					position440, tokenIndex440, depth440 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l441
					}
					position++
					goto l440
				l441:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if buffer[position] != rune('S') {
						goto l436
					}
					position++
				}
			l440:
				if !_rules[rulesp]() {
					goto l436
				}
				if !_rules[ruleIdentifier]() {
					goto l436
				}
				if !_rules[ruleAction25]() {
					goto l436
				}
				depth--
				add(ruleDefAliasedStreamWindow, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
			return false
		},
		/* 34 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action26)> */
		func() bool {
			position442, tokenIndex442, depth442 := position, tokenIndex, depth
			{
				position443 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l442
				}
				if !_rules[rulesp]() {
					goto l442
				}
				if buffer[position] != rune('[') {
					goto l442
				}
				position++
				if !_rules[rulesp]() {
					goto l442
				}
				{
					position444, tokenIndex444, depth444 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l445
					}
					position++
					goto l444
				l445:
					position, tokenIndex, depth = position444, tokenIndex444, depth444
					if buffer[position] != rune('R') {
						goto l442
					}
					position++
				}
			l444:
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
					if buffer[position] != rune('A') {
						goto l442
					}
					position++
				}
			l446:
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
					if buffer[position] != rune('N') {
						goto l442
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('G') {
						goto l442
					}
					position++
				}
			l450:
				{
					position452, tokenIndex452, depth452 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l453
					}
					position++
					goto l452
				l453:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if buffer[position] != rune('E') {
						goto l442
					}
					position++
				}
			l452:
				if !_rules[rulesp]() {
					goto l442
				}
				if !_rules[ruleInterval]() {
					goto l442
				}
				if !_rules[rulesp]() {
					goto l442
				}
				if buffer[position] != rune(']') {
					goto l442
				}
				position++
				if !_rules[ruleAction26]() {
					goto l442
				}
				depth--
				add(ruleStreamWindow, position443)
			}
			return true
		l442:
			position, tokenIndex, depth = position442, tokenIndex442, depth442
			return false
		},
		/* 35 DefStreamWindow <- <(StreamLike (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action27)> */
		func() bool {
			position454, tokenIndex454, depth454 := position, tokenIndex, depth
			{
				position455 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l454
				}
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l456
					}
					if buffer[position] != rune('[') {
						goto l456
					}
					position++
					if !_rules[rulesp]() {
						goto l456
					}
					{
						position458, tokenIndex458, depth458 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l459
						}
						position++
						goto l458
					l459:
						position, tokenIndex, depth = position458, tokenIndex458, depth458
						if buffer[position] != rune('R') {
							goto l456
						}
						position++
					}
				l458:
					{
						position460, tokenIndex460, depth460 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l461
						}
						position++
						goto l460
					l461:
						position, tokenIndex, depth = position460, tokenIndex460, depth460
						if buffer[position] != rune('A') {
							goto l456
						}
						position++
					}
				l460:
					{
						position462, tokenIndex462, depth462 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l463
						}
						position++
						goto l462
					l463:
						position, tokenIndex, depth = position462, tokenIndex462, depth462
						if buffer[position] != rune('N') {
							goto l456
						}
						position++
					}
				l462:
					{
						position464, tokenIndex464, depth464 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l465
						}
						position++
						goto l464
					l465:
						position, tokenIndex, depth = position464, tokenIndex464, depth464
						if buffer[position] != rune('G') {
							goto l456
						}
						position++
					}
				l464:
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l467
						}
						position++
						goto l466
					l467:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
						if buffer[position] != rune('E') {
							goto l456
						}
						position++
					}
				l466:
					if !_rules[rulesp]() {
						goto l456
					}
					if !_rules[ruleInterval]() {
						goto l456
					}
					if !_rules[rulesp]() {
						goto l456
					}
					if buffer[position] != rune(']') {
						goto l456
					}
					position++
					goto l457
				l456:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
				}
			l457:
				if !_rules[ruleAction27]() {
					goto l454
				}
				depth--
				add(ruleDefStreamWindow, position455)
			}
			return true
		l454:
			position, tokenIndex, depth = position454, tokenIndex454, depth454
			return false
		},
		/* 36 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position468, tokenIndex468, depth468 := position, tokenIndex, depth
			{
				position469 := position
				depth++
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l471
					}
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if !_rules[ruleStream]() {
						goto l468
					}
				}
			l470:
				depth--
				add(ruleStreamLike, position469)
			}
			return true
		l468:
			position, tokenIndex, depth = position468, tokenIndex468, depth468
			return false
		},
		/* 37 UDSFFuncApp <- <(FuncApp Action28)> */
		func() bool {
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l472
				}
				if !_rules[ruleAction28]() {
					goto l472
				}
				depth--
				add(ruleUDSFFuncApp, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 38 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action29)> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				{
					position476 := position
					depth++
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						{
							position479, tokenIndex479, depth479 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l480
							}
							position++
							goto l479
						l480:
							position, tokenIndex, depth = position479, tokenIndex479, depth479
							if buffer[position] != rune('W') {
								goto l477
							}
							position++
						}
					l479:
						{
							position481, tokenIndex481, depth481 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l482
							}
							position++
							goto l481
						l482:
							position, tokenIndex, depth = position481, tokenIndex481, depth481
							if buffer[position] != rune('I') {
								goto l477
							}
							position++
						}
					l481:
						{
							position483, tokenIndex483, depth483 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l484
							}
							position++
							goto l483
						l484:
							position, tokenIndex, depth = position483, tokenIndex483, depth483
							if buffer[position] != rune('T') {
								goto l477
							}
							position++
						}
					l483:
						{
							position485, tokenIndex485, depth485 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l486
							}
							position++
							goto l485
						l486:
							position, tokenIndex, depth = position485, tokenIndex485, depth485
							if buffer[position] != rune('H') {
								goto l477
							}
							position++
						}
					l485:
						if !_rules[rulesp]() {
							goto l477
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l477
						}
						if !_rules[rulesp]() {
							goto l477
						}
					l487:
						{
							position488, tokenIndex488, depth488 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l488
							}
							position++
							if !_rules[rulesp]() {
								goto l488
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l488
							}
							goto l487
						l488:
							position, tokenIndex, depth = position488, tokenIndex488, depth488
						}
						goto l478
					l477:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
					}
				l478:
					depth--
					add(rulePegText, position476)
				}
				if !_rules[ruleAction29]() {
					goto l474
				}
				depth--
				add(ruleSourceSinkSpecs, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 39 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action30)> */
		func() bool {
			position489, tokenIndex489, depth489 := position, tokenIndex, depth
			{
				position490 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l489
				}
				if buffer[position] != rune('=') {
					goto l489
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l489
				}
				if !_rules[ruleAction30]() {
					goto l489
				}
				depth--
				add(ruleSourceSinkParam, position490)
			}
			return true
		l489:
			position, tokenIndex, depth = position489, tokenIndex489, depth489
			return false
		},
		/* 40 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position491, tokenIndex491, depth491 := position, tokenIndex, depth
			{
				position492 := position
				depth++
				{
					position493, tokenIndex493, depth493 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l494
					}
					goto l493
				l494:
					position, tokenIndex, depth = position493, tokenIndex493, depth493
					if !_rules[ruleLiteral]() {
						goto l491
					}
				}
			l493:
				depth--
				add(ruleSourceSinkParamVal, position492)
			}
			return true
		l491:
			position, tokenIndex, depth = position491, tokenIndex491, depth491
			return false
		},
		/* 41 PausedOpt <- <(<(Paused / Unpaused)?> Action31)> */
		func() bool {
			position495, tokenIndex495, depth495 := position, tokenIndex, depth
			{
				position496 := position
				depth++
				{
					position497 := position
					depth++
					{
						position498, tokenIndex498, depth498 := position, tokenIndex, depth
						{
							position500, tokenIndex500, depth500 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l501
							}
							goto l500
						l501:
							position, tokenIndex, depth = position500, tokenIndex500, depth500
							if !_rules[ruleUnpaused]() {
								goto l498
							}
						}
					l500:
						goto l499
					l498:
						position, tokenIndex, depth = position498, tokenIndex498, depth498
					}
				l499:
					depth--
					add(rulePegText, position497)
				}
				if !_rules[ruleAction31]() {
					goto l495
				}
				depth--
				add(rulePausedOpt, position496)
			}
			return true
		l495:
			position, tokenIndex, depth = position495, tokenIndex495, depth495
			return false
		},
		/* 42 Expression <- <orExpr> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l502
				}
				depth--
				add(ruleExpression, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 43 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action32)> */
		func() bool {
			position504, tokenIndex504, depth504 := position, tokenIndex, depth
			{
				position505 := position
				depth++
				{
					position506 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l504
					}
					if !_rules[rulesp]() {
						goto l504
					}
					{
						position507, tokenIndex507, depth507 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l507
						}
						if !_rules[rulesp]() {
							goto l507
						}
						if !_rules[ruleandExpr]() {
							goto l507
						}
						goto l508
					l507:
						position, tokenIndex, depth = position507, tokenIndex507, depth507
					}
				l508:
					depth--
					add(rulePegText, position506)
				}
				if !_rules[ruleAction32]() {
					goto l504
				}
				depth--
				add(ruleorExpr, position505)
			}
			return true
		l504:
			position, tokenIndex, depth = position504, tokenIndex504, depth504
			return false
		},
		/* 44 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action33)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				{
					position511 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l509
					}
					if !_rules[rulesp]() {
						goto l509
					}
					{
						position512, tokenIndex512, depth512 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l512
						}
						if !_rules[rulesp]() {
							goto l512
						}
						if !_rules[rulecomparisonExpr]() {
							goto l512
						}
						goto l513
					l512:
						position, tokenIndex, depth = position512, tokenIndex512, depth512
					}
				l513:
					depth--
					add(rulePegText, position511)
				}
				if !_rules[ruleAction33]() {
					goto l509
				}
				depth--
				add(ruleandExpr, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 45 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action34)> */
		func() bool {
			position514, tokenIndex514, depth514 := position, tokenIndex, depth
			{
				position515 := position
				depth++
				{
					position516 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l514
					}
					if !_rules[rulesp]() {
						goto l514
					}
					{
						position517, tokenIndex517, depth517 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l517
						}
						if !_rules[rulesp]() {
							goto l517
						}
						if !_rules[ruletermExpr]() {
							goto l517
						}
						goto l518
					l517:
						position, tokenIndex, depth = position517, tokenIndex517, depth517
					}
				l518:
					depth--
					add(rulePegText, position516)
				}
				if !_rules[ruleAction34]() {
					goto l514
				}
				depth--
				add(rulecomparisonExpr, position515)
			}
			return true
		l514:
			position, tokenIndex, depth = position514, tokenIndex514, depth514
			return false
		},
		/* 46 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action35)> */
		func() bool {
			position519, tokenIndex519, depth519 := position, tokenIndex, depth
			{
				position520 := position
				depth++
				{
					position521 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l519
					}
					if !_rules[rulesp]() {
						goto l519
					}
					{
						position522, tokenIndex522, depth522 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l522
						}
						if !_rules[rulesp]() {
							goto l522
						}
						if !_rules[ruleproductExpr]() {
							goto l522
						}
						goto l523
					l522:
						position, tokenIndex, depth = position522, tokenIndex522, depth522
					}
				l523:
					depth--
					add(rulePegText, position521)
				}
				if !_rules[ruleAction35]() {
					goto l519
				}
				depth--
				add(ruletermExpr, position520)
			}
			return true
		l519:
			position, tokenIndex, depth = position519, tokenIndex519, depth519
			return false
		},
		/* 47 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action36)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				{
					position526 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l524
					}
					if !_rules[rulesp]() {
						goto l524
					}
					{
						position527, tokenIndex527, depth527 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l527
						}
						if !_rules[rulesp]() {
							goto l527
						}
						if !_rules[rulebaseExpr]() {
							goto l527
						}
						goto l528
					l527:
						position, tokenIndex, depth = position527, tokenIndex527, depth527
					}
				l528:
					depth--
					add(rulePegText, position526)
				}
				if !_rules[ruleAction36]() {
					goto l524
				}
				depth--
				add(ruleproductExpr, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 48 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l532
					}
					position++
					if !_rules[rulesp]() {
						goto l532
					}
					if !_rules[ruleExpression]() {
						goto l532
					}
					if !_rules[rulesp]() {
						goto l532
					}
					if buffer[position] != rune(')') {
						goto l532
					}
					position++
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleBooleanLiteral]() {
						goto l533
					}
					goto l531
				l533:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleFuncApp]() {
						goto l534
					}
					goto l531
				l534:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleRowMeta]() {
						goto l535
					}
					goto l531
				l535:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleRowValue]() {
						goto l536
					}
					goto l531
				l536:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleLiteral]() {
						goto l529
					}
				}
			l531:
				depth--
				add(rulebaseExpr, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 49 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action37)> */
		func() bool {
			position537, tokenIndex537, depth537 := position, tokenIndex, depth
			{
				position538 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l537
				}
				if !_rules[rulesp]() {
					goto l537
				}
				if buffer[position] != rune('(') {
					goto l537
				}
				position++
				if !_rules[rulesp]() {
					goto l537
				}
				if !_rules[ruleFuncParams]() {
					goto l537
				}
				if !_rules[rulesp]() {
					goto l537
				}
				if buffer[position] != rune(')') {
					goto l537
				}
				position++
				if !_rules[ruleAction37]() {
					goto l537
				}
				depth--
				add(ruleFuncApp, position538)
			}
			return true
		l537:
			position, tokenIndex, depth = position537, tokenIndex537, depth537
			return false
		},
		/* 50 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action38)> */
		func() bool {
			position539, tokenIndex539, depth539 := position, tokenIndex, depth
			{
				position540 := position
				depth++
				{
					position541 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l539
					}
					if !_rules[rulesp]() {
						goto l539
					}
				l542:
					{
						position543, tokenIndex543, depth543 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l543
						}
						position++
						if !_rules[rulesp]() {
							goto l543
						}
						if !_rules[ruleExpression]() {
							goto l543
						}
						goto l542
					l543:
						position, tokenIndex, depth = position543, tokenIndex543, depth543
					}
					depth--
					add(rulePegText, position541)
				}
				if !_rules[ruleAction38]() {
					goto l539
				}
				depth--
				add(ruleFuncParams, position540)
			}
			return true
		l539:
			position, tokenIndex, depth = position539, tokenIndex539, depth539
			return false
		},
		/* 51 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				{
					position546, tokenIndex546, depth546 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l547
					}
					goto l546
				l547:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
					if !_rules[ruleNumericLiteral]() {
						goto l548
					}
					goto l546
				l548:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
					if !_rules[ruleStringLiteral]() {
						goto l544
					}
				}
			l546:
				depth--
				add(ruleLiteral, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 52 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position549, tokenIndex549, depth549 := position, tokenIndex, depth
			{
				position550 := position
				depth++
				{
					position551, tokenIndex551, depth551 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l552
					}
					goto l551
				l552:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if !_rules[ruleNotEqual]() {
						goto l553
					}
					goto l551
				l553:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if !_rules[ruleLessOrEqual]() {
						goto l554
					}
					goto l551
				l554:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if !_rules[ruleLess]() {
						goto l555
					}
					goto l551
				l555:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if !_rules[ruleGreaterOrEqual]() {
						goto l556
					}
					goto l551
				l556:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if !_rules[ruleGreater]() {
						goto l557
					}
					goto l551
				l557:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if !_rules[ruleNotEqual]() {
						goto l549
					}
				}
			l551:
				depth--
				add(ruleComparisonOp, position550)
			}
			return true
		l549:
			position, tokenIndex, depth = position549, tokenIndex549, depth549
			return false
		},
		/* 53 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position558, tokenIndex558, depth558 := position, tokenIndex, depth
			{
				position559 := position
				depth++
				{
					position560, tokenIndex560, depth560 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l561
					}
					goto l560
				l561:
					position, tokenIndex, depth = position560, tokenIndex560, depth560
					if !_rules[ruleMinus]() {
						goto l558
					}
				}
			l560:
				depth--
				add(rulePlusMinusOp, position559)
			}
			return true
		l558:
			position, tokenIndex, depth = position558, tokenIndex558, depth558
			return false
		},
		/* 54 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position562, tokenIndex562, depth562 := position, tokenIndex, depth
			{
				position563 := position
				depth++
				{
					position564, tokenIndex564, depth564 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l565
					}
					goto l564
				l565:
					position, tokenIndex, depth = position564, tokenIndex564, depth564
					if !_rules[ruleDivide]() {
						goto l566
					}
					goto l564
				l566:
					position, tokenIndex, depth = position564, tokenIndex564, depth564
					if !_rules[ruleModulo]() {
						goto l562
					}
				}
			l564:
				depth--
				add(ruleMultDivOp, position563)
			}
			return true
		l562:
			position, tokenIndex, depth = position562, tokenIndex562, depth562
			return false
		},
		/* 55 Stream <- <(<ident> Action39)> */
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
				if !_rules[ruleAction39]() {
					goto l567
				}
				depth--
				add(ruleStream, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 56 RowMeta <- <RowTimestamp> */
		func() bool {
			position570, tokenIndex570, depth570 := position, tokenIndex, depth
			{
				position571 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l570
				}
				depth--
				add(ruleRowMeta, position571)
			}
			return true
		l570:
			position, tokenIndex, depth = position570, tokenIndex570, depth570
			return false
		},
		/* 57 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action40)> */
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
						if !_rules[ruleident]() {
							goto l575
						}
						if buffer[position] != rune(':') {
							goto l575
						}
						position++
						goto l576
					l575:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
					}
				l576:
					if buffer[position] != rune('t') {
						goto l572
					}
					position++
					if buffer[position] != rune('s') {
						goto l572
					}
					position++
					if buffer[position] != rune('(') {
						goto l572
					}
					position++
					if buffer[position] != rune(')') {
						goto l572
					}
					position++
					depth--
					add(rulePegText, position574)
				}
				if !_rules[ruleAction40]() {
					goto l572
				}
				depth--
				add(ruleRowTimestamp, position573)
			}
			return true
		l572:
			position, tokenIndex, depth = position572, tokenIndex572, depth572
			return false
		},
		/* 58 RowValue <- <(<((ident ':')? ([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.')*)> Action41)> */
		func() bool {
			position577, tokenIndex577, depth577 := position, tokenIndex, depth
			{
				position578 := position
				depth++
				{
					position579 := position
					depth++
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l580
						}
						if buffer[position] != rune(':') {
							goto l580
						}
						position++
						goto l581
					l580:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
					}
				l581:
					{
						position582, tokenIndex582, depth582 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l583
						}
						position++
						goto l582
					l583:
						position, tokenIndex, depth = position582, tokenIndex582, depth582
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l577
						}
						position++
					}
				l582:
				l584:
					{
						position585, tokenIndex585, depth585 := position, tokenIndex, depth
						{
							position586, tokenIndex586, depth586 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l587
							}
							position++
							goto l586
						l587:
							position, tokenIndex, depth = position586, tokenIndex586, depth586
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l588
							}
							position++
							goto l586
						l588:
							position, tokenIndex, depth = position586, tokenIndex586, depth586
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l589
							}
							position++
							goto l586
						l589:
							position, tokenIndex, depth = position586, tokenIndex586, depth586
							if buffer[position] != rune('_') {
								goto l590
							}
							position++
							goto l586
						l590:
							position, tokenIndex, depth = position586, tokenIndex586, depth586
							if buffer[position] != rune('.') {
								goto l585
							}
							position++
						}
					l586:
						goto l584
					l585:
						position, tokenIndex, depth = position585, tokenIndex585, depth585
					}
					depth--
					add(rulePegText, position579)
				}
				if !_rules[ruleAction41]() {
					goto l577
				}
				depth--
				add(ruleRowValue, position578)
			}
			return true
		l577:
			position, tokenIndex, depth = position577, tokenIndex577, depth577
			return false
		},
		/* 59 NumericLiteral <- <(<('-'? [0-9]+)> Action42)> */
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
						if buffer[position] != rune('-') {
							goto l594
						}
						position++
						goto l595
					l594:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
					}
				l595:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l591
					}
					position++
				l596:
					{
						position597, tokenIndex597, depth597 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l597
						}
						position++
						goto l596
					l597:
						position, tokenIndex, depth = position597, tokenIndex597, depth597
					}
					depth--
					add(rulePegText, position593)
				}
				if !_rules[ruleAction42]() {
					goto l591
				}
				depth--
				add(ruleNumericLiteral, position592)
			}
			return true
		l591:
			position, tokenIndex, depth = position591, tokenIndex591, depth591
			return false
		},
		/* 60 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action43)> */
		func() bool {
			position598, tokenIndex598, depth598 := position, tokenIndex, depth
			{
				position599 := position
				depth++
				{
					position600 := position
					depth++
					{
						position601, tokenIndex601, depth601 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l601
						}
						position++
						goto l602
					l601:
						position, tokenIndex, depth = position601, tokenIndex601, depth601
					}
				l602:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l598
					}
					position++
				l603:
					{
						position604, tokenIndex604, depth604 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l604
						}
						position++
						goto l603
					l604:
						position, tokenIndex, depth = position604, tokenIndex604, depth604
					}
					if buffer[position] != rune('.') {
						goto l598
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l598
					}
					position++
				l605:
					{
						position606, tokenIndex606, depth606 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l606
						}
						position++
						goto l605
					l606:
						position, tokenIndex, depth = position606, tokenIndex606, depth606
					}
					depth--
					add(rulePegText, position600)
				}
				if !_rules[ruleAction43]() {
					goto l598
				}
				depth--
				add(ruleFloatLiteral, position599)
			}
			return true
		l598:
			position, tokenIndex, depth = position598, tokenIndex598, depth598
			return false
		},
		/* 61 Function <- <(<ident> Action44)> */
		func() bool {
			position607, tokenIndex607, depth607 := position, tokenIndex, depth
			{
				position608 := position
				depth++
				{
					position609 := position
					depth++
					if !_rules[ruleident]() {
						goto l607
					}
					depth--
					add(rulePegText, position609)
				}
				if !_rules[ruleAction44]() {
					goto l607
				}
				depth--
				add(ruleFunction, position608)
			}
			return true
		l607:
			position, tokenIndex, depth = position607, tokenIndex607, depth607
			return false
		},
		/* 62 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position610, tokenIndex610, depth610 := position, tokenIndex, depth
			{
				position611 := position
				depth++
				{
					position612, tokenIndex612, depth612 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l613
					}
					goto l612
				l613:
					position, tokenIndex, depth = position612, tokenIndex612, depth612
					if !_rules[ruleFALSE]() {
						goto l610
					}
				}
			l612:
				depth--
				add(ruleBooleanLiteral, position611)
			}
			return true
		l610:
			position, tokenIndex, depth = position610, tokenIndex610, depth610
			return false
		},
		/* 63 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action45)> */
		func() bool {
			position614, tokenIndex614, depth614 := position, tokenIndex, depth
			{
				position615 := position
				depth++
				{
					position616 := position
					depth++
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
							goto l614
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
							goto l614
						}
						position++
					}
				l619:
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l622
						}
						position++
						goto l621
					l622:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
						if buffer[position] != rune('U') {
							goto l614
						}
						position++
					}
				l621:
					{
						position623, tokenIndex623, depth623 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l624
						}
						position++
						goto l623
					l624:
						position, tokenIndex, depth = position623, tokenIndex623, depth623
						if buffer[position] != rune('E') {
							goto l614
						}
						position++
					}
				l623:
					depth--
					add(rulePegText, position616)
				}
				if !_rules[ruleAction45]() {
					goto l614
				}
				depth--
				add(ruleTRUE, position615)
			}
			return true
		l614:
			position, tokenIndex, depth = position614, tokenIndex614, depth614
			return false
		},
		/* 64 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action46)> */
		func() bool {
			position625, tokenIndex625, depth625 := position, tokenIndex, depth
			{
				position626 := position
				depth++
				{
					position627 := position
					depth++
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('F') {
							goto l625
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('A') {
							goto l625
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
							goto l625
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
							goto l625
						}
						position++
					}
				l634:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l637
						}
						position++
						goto l636
					l637:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
						if buffer[position] != rune('E') {
							goto l625
						}
						position++
					}
				l636:
					depth--
					add(rulePegText, position627)
				}
				if !_rules[ruleAction46]() {
					goto l625
				}
				depth--
				add(ruleFALSE, position626)
			}
			return true
		l625:
			position, tokenIndex, depth = position625, tokenIndex625, depth625
			return false
		},
		/* 65 Wildcard <- <(<'*'> Action47)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					if buffer[position] != rune('*') {
						goto l638
					}
					position++
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction47]() {
					goto l638
				}
				depth--
				add(ruleWildcard, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 66 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action48)> */
		func() bool {
			position641, tokenIndex641, depth641 := position, tokenIndex, depth
			{
				position642 := position
				depth++
				{
					position643 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l641
					}
					position++
				l644:
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						{
							position646, tokenIndex646, depth646 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l647
							}
							position++
							if buffer[position] != rune('\'') {
								goto l647
							}
							position++
							goto l646
						l647:
							position, tokenIndex, depth = position646, tokenIndex646, depth646
							{
								position648, tokenIndex648, depth648 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l648
								}
								position++
								goto l645
							l648:
								position, tokenIndex, depth = position648, tokenIndex648, depth648
							}
							if !matchDot() {
								goto l645
							}
						}
					l646:
						goto l644
					l645:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
					}
					if buffer[position] != rune('\'') {
						goto l641
					}
					position++
					depth--
					add(rulePegText, position643)
				}
				if !_rules[ruleAction48]() {
					goto l641
				}
				depth--
				add(ruleStringLiteral, position642)
			}
			return true
		l641:
			position, tokenIndex, depth = position641, tokenIndex641, depth641
			return false
		},
		/* 67 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action49)> */
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
						if buffer[position] != rune('i') {
							goto l653
						}
						position++
						goto l652
					l653:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
						if buffer[position] != rune('I') {
							goto l649
						}
						position++
					}
				l652:
					{
						position654, tokenIndex654, depth654 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l655
						}
						position++
						goto l654
					l655:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
						if buffer[position] != rune('S') {
							goto l649
						}
						position++
					}
				l654:
					{
						position656, tokenIndex656, depth656 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l657
						}
						position++
						goto l656
					l657:
						position, tokenIndex, depth = position656, tokenIndex656, depth656
						if buffer[position] != rune('T') {
							goto l649
						}
						position++
					}
				l656:
					{
						position658, tokenIndex658, depth658 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l659
						}
						position++
						goto l658
					l659:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
						if buffer[position] != rune('R') {
							goto l649
						}
						position++
					}
				l658:
					{
						position660, tokenIndex660, depth660 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l661
						}
						position++
						goto l660
					l661:
						position, tokenIndex, depth = position660, tokenIndex660, depth660
						if buffer[position] != rune('E') {
							goto l649
						}
						position++
					}
				l660:
					{
						position662, tokenIndex662, depth662 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l663
						}
						position++
						goto l662
					l663:
						position, tokenIndex, depth = position662, tokenIndex662, depth662
						if buffer[position] != rune('A') {
							goto l649
						}
						position++
					}
				l662:
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l665
						}
						position++
						goto l664
					l665:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if buffer[position] != rune('M') {
							goto l649
						}
						position++
					}
				l664:
					depth--
					add(rulePegText, position651)
				}
				if !_rules[ruleAction49]() {
					goto l649
				}
				depth--
				add(ruleISTREAM, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 68 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action50)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				{
					position668 := position
					depth++
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l670
						}
						position++
						goto l669
					l670:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if buffer[position] != rune('D') {
							goto l666
						}
						position++
					}
				l669:
					{
						position671, tokenIndex671, depth671 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l672
						}
						position++
						goto l671
					l672:
						position, tokenIndex, depth = position671, tokenIndex671, depth671
						if buffer[position] != rune('S') {
							goto l666
						}
						position++
					}
				l671:
					{
						position673, tokenIndex673, depth673 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l674
						}
						position++
						goto l673
					l674:
						position, tokenIndex, depth = position673, tokenIndex673, depth673
						if buffer[position] != rune('T') {
							goto l666
						}
						position++
					}
				l673:
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l676
						}
						position++
						goto l675
					l676:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('R') {
							goto l666
						}
						position++
					}
				l675:
					{
						position677, tokenIndex677, depth677 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l678
						}
						position++
						goto l677
					l678:
						position, tokenIndex, depth = position677, tokenIndex677, depth677
						if buffer[position] != rune('E') {
							goto l666
						}
						position++
					}
				l677:
					{
						position679, tokenIndex679, depth679 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l680
						}
						position++
						goto l679
					l680:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
						if buffer[position] != rune('A') {
							goto l666
						}
						position++
					}
				l679:
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if buffer[position] != rune('M') {
							goto l666
						}
						position++
					}
				l681:
					depth--
					add(rulePegText, position668)
				}
				if !_rules[ruleAction50]() {
					goto l666
				}
				depth--
				add(ruleDSTREAM, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 69 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action51)> */
		func() bool {
			position683, tokenIndex683, depth683 := position, tokenIndex, depth
			{
				position684 := position
				depth++
				{
					position685 := position
					depth++
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
							goto l683
						}
						position++
					}
				l686:
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if buffer[position] != rune('S') {
							goto l683
						}
						position++
					}
				l688:
					{
						position690, tokenIndex690, depth690 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l691
						}
						position++
						goto l690
					l691:
						position, tokenIndex, depth = position690, tokenIndex690, depth690
						if buffer[position] != rune('T') {
							goto l683
						}
						position++
					}
				l690:
					{
						position692, tokenIndex692, depth692 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l693
						}
						position++
						goto l692
					l693:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						if buffer[position] != rune('R') {
							goto l683
						}
						position++
					}
				l692:
					{
						position694, tokenIndex694, depth694 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l695
						}
						position++
						goto l694
					l695:
						position, tokenIndex, depth = position694, tokenIndex694, depth694
						if buffer[position] != rune('E') {
							goto l683
						}
						position++
					}
				l694:
					{
						position696, tokenIndex696, depth696 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l697
						}
						position++
						goto l696
					l697:
						position, tokenIndex, depth = position696, tokenIndex696, depth696
						if buffer[position] != rune('A') {
							goto l683
						}
						position++
					}
				l696:
					{
						position698, tokenIndex698, depth698 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l699
						}
						position++
						goto l698
					l699:
						position, tokenIndex, depth = position698, tokenIndex698, depth698
						if buffer[position] != rune('M') {
							goto l683
						}
						position++
					}
				l698:
					depth--
					add(rulePegText, position685)
				}
				if !_rules[ruleAction51]() {
					goto l683
				}
				depth--
				add(ruleRSTREAM, position684)
			}
			return true
		l683:
			position, tokenIndex, depth = position683, tokenIndex683, depth683
			return false
		},
		/* 70 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action52)> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				{
					position702 := position
					depth++
					{
						position703, tokenIndex703, depth703 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l704
						}
						position++
						goto l703
					l704:
						position, tokenIndex, depth = position703, tokenIndex703, depth703
						if buffer[position] != rune('T') {
							goto l700
						}
						position++
					}
				l703:
					{
						position705, tokenIndex705, depth705 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l706
						}
						position++
						goto l705
					l706:
						position, tokenIndex, depth = position705, tokenIndex705, depth705
						if buffer[position] != rune('U') {
							goto l700
						}
						position++
					}
				l705:
					{
						position707, tokenIndex707, depth707 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l708
						}
						position++
						goto l707
					l708:
						position, tokenIndex, depth = position707, tokenIndex707, depth707
						if buffer[position] != rune('P') {
							goto l700
						}
						position++
					}
				l707:
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l710
						}
						position++
						goto l709
					l710:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if buffer[position] != rune('L') {
							goto l700
						}
						position++
					}
				l709:
					{
						position711, tokenIndex711, depth711 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l712
						}
						position++
						goto l711
					l712:
						position, tokenIndex, depth = position711, tokenIndex711, depth711
						if buffer[position] != rune('E') {
							goto l700
						}
						position++
					}
				l711:
					{
						position713, tokenIndex713, depth713 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l714
						}
						position++
						goto l713
					l714:
						position, tokenIndex, depth = position713, tokenIndex713, depth713
						if buffer[position] != rune('S') {
							goto l700
						}
						position++
					}
				l713:
					depth--
					add(rulePegText, position702)
				}
				if !_rules[ruleAction52]() {
					goto l700
				}
				depth--
				add(ruleTUPLES, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 71 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action53)> */
		func() bool {
			position715, tokenIndex715, depth715 := position, tokenIndex, depth
			{
				position716 := position
				depth++
				{
					position717 := position
					depth++
					{
						position718, tokenIndex718, depth718 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l719
						}
						position++
						goto l718
					l719:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
						if buffer[position] != rune('S') {
							goto l715
						}
						position++
					}
				l718:
					{
						position720, tokenIndex720, depth720 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l721
						}
						position++
						goto l720
					l721:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
						if buffer[position] != rune('E') {
							goto l715
						}
						position++
					}
				l720:
					{
						position722, tokenIndex722, depth722 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l723
						}
						position++
						goto l722
					l723:
						position, tokenIndex, depth = position722, tokenIndex722, depth722
						if buffer[position] != rune('C') {
							goto l715
						}
						position++
					}
				l722:
					{
						position724, tokenIndex724, depth724 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l725
						}
						position++
						goto l724
					l725:
						position, tokenIndex, depth = position724, tokenIndex724, depth724
						if buffer[position] != rune('O') {
							goto l715
						}
						position++
					}
				l724:
					{
						position726, tokenIndex726, depth726 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l727
						}
						position++
						goto l726
					l727:
						position, tokenIndex, depth = position726, tokenIndex726, depth726
						if buffer[position] != rune('N') {
							goto l715
						}
						position++
					}
				l726:
					{
						position728, tokenIndex728, depth728 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l729
						}
						position++
						goto l728
					l729:
						position, tokenIndex, depth = position728, tokenIndex728, depth728
						if buffer[position] != rune('D') {
							goto l715
						}
						position++
					}
				l728:
					{
						position730, tokenIndex730, depth730 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l731
						}
						position++
						goto l730
					l731:
						position, tokenIndex, depth = position730, tokenIndex730, depth730
						if buffer[position] != rune('S') {
							goto l715
						}
						position++
					}
				l730:
					depth--
					add(rulePegText, position717)
				}
				if !_rules[ruleAction53]() {
					goto l715
				}
				depth--
				add(ruleSECONDS, position716)
			}
			return true
		l715:
			position, tokenIndex, depth = position715, tokenIndex715, depth715
			return false
		},
		/* 72 StreamIdentifier <- <(<ident> Action54)> */
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
				if !_rules[ruleAction54]() {
					goto l732
				}
				depth--
				add(ruleStreamIdentifier, position733)
			}
			return true
		l732:
			position, tokenIndex, depth = position732, tokenIndex732, depth732
			return false
		},
		/* 73 SourceSinkType <- <(<ident> Action55)> */
		func() bool {
			position735, tokenIndex735, depth735 := position, tokenIndex, depth
			{
				position736 := position
				depth++
				{
					position737 := position
					depth++
					if !_rules[ruleident]() {
						goto l735
					}
					depth--
					add(rulePegText, position737)
				}
				if !_rules[ruleAction55]() {
					goto l735
				}
				depth--
				add(ruleSourceSinkType, position736)
			}
			return true
		l735:
			position, tokenIndex, depth = position735, tokenIndex735, depth735
			return false
		},
		/* 74 SourceSinkParamKey <- <(<ident> Action56)> */
		func() bool {
			position738, tokenIndex738, depth738 := position, tokenIndex, depth
			{
				position739 := position
				depth++
				{
					position740 := position
					depth++
					if !_rules[ruleident]() {
						goto l738
					}
					depth--
					add(rulePegText, position740)
				}
				if !_rules[ruleAction56]() {
					goto l738
				}
				depth--
				add(ruleSourceSinkParamKey, position739)
			}
			return true
		l738:
			position, tokenIndex, depth = position738, tokenIndex738, depth738
			return false
		},
		/* 75 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action57)> */
		func() bool {
			position741, tokenIndex741, depth741 := position, tokenIndex, depth
			{
				position742 := position
				depth++
				{
					position743 := position
					depth++
					{
						position744, tokenIndex744, depth744 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l745
						}
						position++
						goto l744
					l745:
						position, tokenIndex, depth = position744, tokenIndex744, depth744
						if buffer[position] != rune('P') {
							goto l741
						}
						position++
					}
				l744:
					{
						position746, tokenIndex746, depth746 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l747
						}
						position++
						goto l746
					l747:
						position, tokenIndex, depth = position746, tokenIndex746, depth746
						if buffer[position] != rune('A') {
							goto l741
						}
						position++
					}
				l746:
					{
						position748, tokenIndex748, depth748 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l749
						}
						position++
						goto l748
					l749:
						position, tokenIndex, depth = position748, tokenIndex748, depth748
						if buffer[position] != rune('U') {
							goto l741
						}
						position++
					}
				l748:
					{
						position750, tokenIndex750, depth750 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l751
						}
						position++
						goto l750
					l751:
						position, tokenIndex, depth = position750, tokenIndex750, depth750
						if buffer[position] != rune('S') {
							goto l741
						}
						position++
					}
				l750:
					{
						position752, tokenIndex752, depth752 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l753
						}
						position++
						goto l752
					l753:
						position, tokenIndex, depth = position752, tokenIndex752, depth752
						if buffer[position] != rune('E') {
							goto l741
						}
						position++
					}
				l752:
					{
						position754, tokenIndex754, depth754 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l755
						}
						position++
						goto l754
					l755:
						position, tokenIndex, depth = position754, tokenIndex754, depth754
						if buffer[position] != rune('D') {
							goto l741
						}
						position++
					}
				l754:
					depth--
					add(rulePegText, position743)
				}
				if !_rules[ruleAction57]() {
					goto l741
				}
				depth--
				add(rulePaused, position742)
			}
			return true
		l741:
			position, tokenIndex, depth = position741, tokenIndex741, depth741
			return false
		},
		/* 76 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action58)> */
		func() bool {
			position756, tokenIndex756, depth756 := position, tokenIndex, depth
			{
				position757 := position
				depth++
				{
					position758 := position
					depth++
					{
						position759, tokenIndex759, depth759 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l760
						}
						position++
						goto l759
					l760:
						position, tokenIndex, depth = position759, tokenIndex759, depth759
						if buffer[position] != rune('U') {
							goto l756
						}
						position++
					}
				l759:
					{
						position761, tokenIndex761, depth761 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l762
						}
						position++
						goto l761
					l762:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if buffer[position] != rune('N') {
							goto l756
						}
						position++
					}
				l761:
					{
						position763, tokenIndex763, depth763 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l764
						}
						position++
						goto l763
					l764:
						position, tokenIndex, depth = position763, tokenIndex763, depth763
						if buffer[position] != rune('P') {
							goto l756
						}
						position++
					}
				l763:
					{
						position765, tokenIndex765, depth765 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l766
						}
						position++
						goto l765
					l766:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
						if buffer[position] != rune('A') {
							goto l756
						}
						position++
					}
				l765:
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('U') {
							goto l756
						}
						position++
					}
				l767:
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('S') {
							goto l756
						}
						position++
					}
				l769:
					{
						position771, tokenIndex771, depth771 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l772
						}
						position++
						goto l771
					l772:
						position, tokenIndex, depth = position771, tokenIndex771, depth771
						if buffer[position] != rune('E') {
							goto l756
						}
						position++
					}
				l771:
					{
						position773, tokenIndex773, depth773 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l774
						}
						position++
						goto l773
					l774:
						position, tokenIndex, depth = position773, tokenIndex773, depth773
						if buffer[position] != rune('D') {
							goto l756
						}
						position++
					}
				l773:
					depth--
					add(rulePegText, position758)
				}
				if !_rules[ruleAction58]() {
					goto l756
				}
				depth--
				add(ruleUnpaused, position757)
			}
			return true
		l756:
			position, tokenIndex, depth = position756, tokenIndex756, depth756
			return false
		},
		/* 77 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action59)> */
		func() bool {
			position775, tokenIndex775, depth775 := position, tokenIndex, depth
			{
				position776 := position
				depth++
				{
					position777 := position
					depth++
					{
						position778, tokenIndex778, depth778 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l779
						}
						position++
						goto l778
					l779:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						if buffer[position] != rune('O') {
							goto l775
						}
						position++
					}
				l778:
					{
						position780, tokenIndex780, depth780 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l781
						}
						position++
						goto l780
					l781:
						position, tokenIndex, depth = position780, tokenIndex780, depth780
						if buffer[position] != rune('R') {
							goto l775
						}
						position++
					}
				l780:
					depth--
					add(rulePegText, position777)
				}
				if !_rules[ruleAction59]() {
					goto l775
				}
				depth--
				add(ruleOr, position776)
			}
			return true
		l775:
			position, tokenIndex, depth = position775, tokenIndex775, depth775
			return false
		},
		/* 78 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action60)> */
		func() bool {
			position782, tokenIndex782, depth782 := position, tokenIndex, depth
			{
				position783 := position
				depth++
				{
					position784 := position
					depth++
					{
						position785, tokenIndex785, depth785 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l786
						}
						position++
						goto l785
					l786:
						position, tokenIndex, depth = position785, tokenIndex785, depth785
						if buffer[position] != rune('A') {
							goto l782
						}
						position++
					}
				l785:
					{
						position787, tokenIndex787, depth787 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l788
						}
						position++
						goto l787
					l788:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
						if buffer[position] != rune('N') {
							goto l782
						}
						position++
					}
				l787:
					{
						position789, tokenIndex789, depth789 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l790
						}
						position++
						goto l789
					l790:
						position, tokenIndex, depth = position789, tokenIndex789, depth789
						if buffer[position] != rune('D') {
							goto l782
						}
						position++
					}
				l789:
					depth--
					add(rulePegText, position784)
				}
				if !_rules[ruleAction60]() {
					goto l782
				}
				depth--
				add(ruleAnd, position783)
			}
			return true
		l782:
			position, tokenIndex, depth = position782, tokenIndex782, depth782
			return false
		},
		/* 79 Equal <- <(<'='> Action61)> */
		func() bool {
			position791, tokenIndex791, depth791 := position, tokenIndex, depth
			{
				position792 := position
				depth++
				{
					position793 := position
					depth++
					if buffer[position] != rune('=') {
						goto l791
					}
					position++
					depth--
					add(rulePegText, position793)
				}
				if !_rules[ruleAction61]() {
					goto l791
				}
				depth--
				add(ruleEqual, position792)
			}
			return true
		l791:
			position, tokenIndex, depth = position791, tokenIndex791, depth791
			return false
		},
		/* 80 Less <- <(<'<'> Action62)> */
		func() bool {
			position794, tokenIndex794, depth794 := position, tokenIndex, depth
			{
				position795 := position
				depth++
				{
					position796 := position
					depth++
					if buffer[position] != rune('<') {
						goto l794
					}
					position++
					depth--
					add(rulePegText, position796)
				}
				if !_rules[ruleAction62]() {
					goto l794
				}
				depth--
				add(ruleLess, position795)
			}
			return true
		l794:
			position, tokenIndex, depth = position794, tokenIndex794, depth794
			return false
		},
		/* 81 LessOrEqual <- <(<('<' '=')> Action63)> */
		func() bool {
			position797, tokenIndex797, depth797 := position, tokenIndex, depth
			{
				position798 := position
				depth++
				{
					position799 := position
					depth++
					if buffer[position] != rune('<') {
						goto l797
					}
					position++
					if buffer[position] != rune('=') {
						goto l797
					}
					position++
					depth--
					add(rulePegText, position799)
				}
				if !_rules[ruleAction63]() {
					goto l797
				}
				depth--
				add(ruleLessOrEqual, position798)
			}
			return true
		l797:
			position, tokenIndex, depth = position797, tokenIndex797, depth797
			return false
		},
		/* 82 Greater <- <(<'>'> Action64)> */
		func() bool {
			position800, tokenIndex800, depth800 := position, tokenIndex, depth
			{
				position801 := position
				depth++
				{
					position802 := position
					depth++
					if buffer[position] != rune('>') {
						goto l800
					}
					position++
					depth--
					add(rulePegText, position802)
				}
				if !_rules[ruleAction64]() {
					goto l800
				}
				depth--
				add(ruleGreater, position801)
			}
			return true
		l800:
			position, tokenIndex, depth = position800, tokenIndex800, depth800
			return false
		},
		/* 83 GreaterOrEqual <- <(<('>' '=')> Action65)> */
		func() bool {
			position803, tokenIndex803, depth803 := position, tokenIndex, depth
			{
				position804 := position
				depth++
				{
					position805 := position
					depth++
					if buffer[position] != rune('>') {
						goto l803
					}
					position++
					if buffer[position] != rune('=') {
						goto l803
					}
					position++
					depth--
					add(rulePegText, position805)
				}
				if !_rules[ruleAction65]() {
					goto l803
				}
				depth--
				add(ruleGreaterOrEqual, position804)
			}
			return true
		l803:
			position, tokenIndex, depth = position803, tokenIndex803, depth803
			return false
		},
		/* 84 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action66)> */
		func() bool {
			position806, tokenIndex806, depth806 := position, tokenIndex, depth
			{
				position807 := position
				depth++
				{
					position808 := position
					depth++
					{
						position809, tokenIndex809, depth809 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l810
						}
						position++
						if buffer[position] != rune('=') {
							goto l810
						}
						position++
						goto l809
					l810:
						position, tokenIndex, depth = position809, tokenIndex809, depth809
						if buffer[position] != rune('<') {
							goto l806
						}
						position++
						if buffer[position] != rune('>') {
							goto l806
						}
						position++
					}
				l809:
					depth--
					add(rulePegText, position808)
				}
				if !_rules[ruleAction66]() {
					goto l806
				}
				depth--
				add(ruleNotEqual, position807)
			}
			return true
		l806:
			position, tokenIndex, depth = position806, tokenIndex806, depth806
			return false
		},
		/* 85 Plus <- <(<'+'> Action67)> */
		func() bool {
			position811, tokenIndex811, depth811 := position, tokenIndex, depth
			{
				position812 := position
				depth++
				{
					position813 := position
					depth++
					if buffer[position] != rune('+') {
						goto l811
					}
					position++
					depth--
					add(rulePegText, position813)
				}
				if !_rules[ruleAction67]() {
					goto l811
				}
				depth--
				add(rulePlus, position812)
			}
			return true
		l811:
			position, tokenIndex, depth = position811, tokenIndex811, depth811
			return false
		},
		/* 86 Minus <- <(<'-'> Action68)> */
		func() bool {
			position814, tokenIndex814, depth814 := position, tokenIndex, depth
			{
				position815 := position
				depth++
				{
					position816 := position
					depth++
					if buffer[position] != rune('-') {
						goto l814
					}
					position++
					depth--
					add(rulePegText, position816)
				}
				if !_rules[ruleAction68]() {
					goto l814
				}
				depth--
				add(ruleMinus, position815)
			}
			return true
		l814:
			position, tokenIndex, depth = position814, tokenIndex814, depth814
			return false
		},
		/* 87 Multiply <- <(<'*'> Action69)> */
		func() bool {
			position817, tokenIndex817, depth817 := position, tokenIndex, depth
			{
				position818 := position
				depth++
				{
					position819 := position
					depth++
					if buffer[position] != rune('*') {
						goto l817
					}
					position++
					depth--
					add(rulePegText, position819)
				}
				if !_rules[ruleAction69]() {
					goto l817
				}
				depth--
				add(ruleMultiply, position818)
			}
			return true
		l817:
			position, tokenIndex, depth = position817, tokenIndex817, depth817
			return false
		},
		/* 88 Divide <- <(<'/'> Action70)> */
		func() bool {
			position820, tokenIndex820, depth820 := position, tokenIndex, depth
			{
				position821 := position
				depth++
				{
					position822 := position
					depth++
					if buffer[position] != rune('/') {
						goto l820
					}
					position++
					depth--
					add(rulePegText, position822)
				}
				if !_rules[ruleAction70]() {
					goto l820
				}
				depth--
				add(ruleDivide, position821)
			}
			return true
		l820:
			position, tokenIndex, depth = position820, tokenIndex820, depth820
			return false
		},
		/* 89 Modulo <- <(<'%'> Action71)> */
		func() bool {
			position823, tokenIndex823, depth823 := position, tokenIndex, depth
			{
				position824 := position
				depth++
				{
					position825 := position
					depth++
					if buffer[position] != rune('%') {
						goto l823
					}
					position++
					depth--
					add(rulePegText, position825)
				}
				if !_rules[ruleAction71]() {
					goto l823
				}
				depth--
				add(ruleModulo, position824)
			}
			return true
		l823:
			position, tokenIndex, depth = position823, tokenIndex823, depth823
			return false
		},
		/* 90 Identifier <- <(<ident> Action72)> */
		func() bool {
			position826, tokenIndex826, depth826 := position, tokenIndex, depth
			{
				position827 := position
				depth++
				{
					position828 := position
					depth++
					if !_rules[ruleident]() {
						goto l826
					}
					depth--
					add(rulePegText, position828)
				}
				if !_rules[ruleAction72]() {
					goto l826
				}
				depth--
				add(ruleIdentifier, position827)
			}
			return true
		l826:
			position, tokenIndex, depth = position826, tokenIndex826, depth826
			return false
		},
		/* 91 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position829, tokenIndex829, depth829 := position, tokenIndex, depth
			{
				position830 := position
				depth++
				{
					position831, tokenIndex831, depth831 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l832
					}
					position++
					goto l831
				l832:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l829
					}
					position++
				}
			l831:
			l833:
				{
					position834, tokenIndex834, depth834 := position, tokenIndex, depth
					{
						position835, tokenIndex835, depth835 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l836
						}
						position++
						goto l835
					l836:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l837
						}
						position++
						goto l835
					l837:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l838
						}
						position++
						goto l835
					l838:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
						if buffer[position] != rune('_') {
							goto l834
						}
						position++
					}
				l835:
					goto l833
				l834:
					position, tokenIndex, depth = position834, tokenIndex834, depth834
				}
				depth--
				add(ruleident, position830)
			}
			return true
		l829:
			position, tokenIndex, depth = position829, tokenIndex829, depth829
			return false
		},
		/* 92 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position840 := position
				depth++
			l841:
				{
					position842, tokenIndex842, depth842 := position, tokenIndex, depth
					{
						position843, tokenIndex843, depth843 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l844
						}
						position++
						goto l843
					l844:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						if buffer[position] != rune('\t') {
							goto l845
						}
						position++
						goto l843
					l845:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						if buffer[position] != rune('\n') {
							goto l842
						}
						position++
					}
				l843:
					goto l841
				l842:
					position, tokenIndex, depth = position842, tokenIndex842, depth842
				}
				depth--
				add(rulesp, position840)
			}
			return true
		},
		/* 94 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 95 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 96 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 97 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 98 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 99 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 100 Action6 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 101 Action7 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 102 Action8 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		nil,
		/* 104 Action9 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 105 Action10 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 106 Action11 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 107 Action12 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 108 Action13 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 109 Action14 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 110 Action15 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 111 Action16 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 112 Action17 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 113 Action18 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 114 Action19 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 115 Action20 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 116 Action21 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 117 Action22 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 118 Action23 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 119 Action24 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 120 Action25 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 121 Action26 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 122 Action27 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 123 Action28 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 124 Action29 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 125 Action30 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 126 Action31 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 127 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 128 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 129 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 130 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 131 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 132 Action37 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 133 Action38 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 134 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 135 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 136 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 137 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 138 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 139 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 140 Action45 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 141 Action46 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 142 Action47 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 143 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 144 Action49 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 145 Action50 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 146 Action51 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 147 Action52 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 148 Action53 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 149 Action54 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 150 Action55 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 151 Action56 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 152 Action57 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 153 Action58 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 154 Action59 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 155 Action60 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 156 Action61 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 157 Action62 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 158 Action63 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 159 Action64 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 160 Action65 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 161 Action66 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 162 Action67 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 163 Action68 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 164 Action69 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 165 Action70 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 166 Action71 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 167 Action72 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
	}
	p.rules = _rules
}

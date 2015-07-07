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
	ruleInsertIntoFromStmt
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
	ruleisExpr
	ruletermExpr
	ruleproductExpr
	rulebaseExpr
	ruleFuncApp
	ruleFuncParams
	ruleLiteral
	ruleComparisonOp
	ruleIsOp
	rulePlusMinusOp
	ruleMultDivOp
	ruleStream
	ruleRowMeta
	ruleRowTimestamp
	ruleRowValue
	ruleNumericLiteral
	ruleFloatLiteral
	ruleFunction
	ruleNullLiteral
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
	ruleIs
	ruleIsNot
	rulePlus
	ruleMinus
	ruleMultiply
	ruleDivide
	ruleModulo
	ruleIdentifier
	ruleTargetIdentifier
	ruleident
	rulejsonPath
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
	ruleAction9
	rulePegText
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
	ruleAction73
	ruleAction74
	ruleAction75
	ruleAction76
	ruleAction77
	ruleAction78

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
	"InsertIntoFromStmt",
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
	"isExpr",
	"termExpr",
	"productExpr",
	"baseExpr",
	"FuncApp",
	"FuncParams",
	"Literal",
	"ComparisonOp",
	"IsOp",
	"PlusMinusOp",
	"MultDivOp",
	"Stream",
	"RowMeta",
	"RowTimestamp",
	"RowValue",
	"NumericLiteral",
	"FloatLiteral",
	"Function",
	"NullLiteral",
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
	"Is",
	"IsNot",
	"Plus",
	"Minus",
	"Multiply",
	"Divide",
	"Modulo",
	"Identifier",
	"TargetIdentifier",
	"ident",
	"jsonPath",
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
	"Action9",
	"PegText",
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
	"Action73",
	"Action74",
	"Action75",
	"Action76",
	"Action77",
	"Action78",

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
	rules  [182]func() bool
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

			p.AssembleInsertIntoFrom()

		case ruleAction7:

			p.AssemblePauseSource()

		case ruleAction8:

			p.AssembleResumeSource()

		case ruleAction9:

			p.AssembleRewindSource()

		case ruleAction10:

			p.AssembleEmitter(begin, end)

		case ruleAction11:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction12:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction13:

			p.AssembleStreamEmitInterval()

		case ruleAction14:

			p.AssembleProjections(begin, end)

		case ruleAction15:

			p.AssembleAlias()

		case ruleAction16:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction17:

			p.AssembleWindowedFrom(begin, end)

		case ruleAction18:

			p.AssembleInterval()

		case ruleAction19:

			p.AssembleInterval()

		case ruleAction20:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction21:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction22:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction23:

			p.EnsureAliasedStreamWindow()

		case ruleAction24:

			p.EnsureAliasedStreamWindow()

		case ruleAction25:

			p.AssembleAliasedStreamWindow()

		case ruleAction26:

			p.AssembleAliasedStreamWindow()

		case ruleAction27:

			p.AssembleStreamWindow()

		case ruleAction28:

			p.AssembleStreamWindow()

		case ruleAction29:

			p.AssembleUDSFFuncApp()

		case ruleAction30:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction31:

			p.AssembleSourceSinkParam()

		case ruleAction32:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction33:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction34:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction38:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction39:

			p.AssembleFuncApp()

		case ruleAction40:

			p.AssembleExpressions(begin, end)

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction44:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction47:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction48:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction49:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction50:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction52:

			p.PushComponent(begin, end, Istream)

		case ruleAction53:

			p.PushComponent(begin, end, Dstream)

		case ruleAction54:

			p.PushComponent(begin, end, Rstream)

		case ruleAction55:

			p.PushComponent(begin, end, Tuples)

		case ruleAction56:

			p.PushComponent(begin, end, Seconds)

		case ruleAction57:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction58:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction59:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction60:

			p.PushComponent(begin, end, Yes)

		case ruleAction61:

			p.PushComponent(begin, end, No)

		case ruleAction62:

			p.PushComponent(begin, end, Or)

		case ruleAction63:

			p.PushComponent(begin, end, And)

		case ruleAction64:

			p.PushComponent(begin, end, Equal)

		case ruleAction65:

			p.PushComponent(begin, end, Less)

		case ruleAction66:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction67:

			p.PushComponent(begin, end, Greater)

		case ruleAction68:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction69:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction70:

			p.PushComponent(begin, end, Is)

		case ruleAction71:

			p.PushComponent(begin, end, IsNot)

		case ruleAction72:

			p.PushComponent(begin, end, Plus)

		case ruleAction73:

			p.PushComponent(begin, end, Minus)

		case ruleAction74:

			p.PushComponent(begin, end, Multiply)

		case ruleAction75:

			p.PushComponent(begin, end, Divide)

		case ruleAction76:

			p.PushComponent(begin, end, Modulo)

		case ruleAction77:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction78:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt)> */
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
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l15
					}
					goto l9
				l15:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleCreateStateStmt]() {
						goto l16
					}
					goto l9
				l16:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[rulePauseSourceStmt]() {
						goto l17
					}
					goto l9
				l17:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleResumeSourceStmt]() {
						goto l18
					}
					goto l9
				l18:
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
			position19, tokenIndex19, depth19 := position, tokenIndex, depth
			{
				position20 := position
				depth++
				{
					position21, tokenIndex21, depth21 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l22
					}
					position++
					goto l21
				l22:
					position, tokenIndex, depth = position21, tokenIndex21, depth21
					if buffer[position] != rune('S') {
						goto l19
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
						goto l19
					}
					position++
				}
			l23:
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l26
					}
					position++
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if buffer[position] != rune('L') {
						goto l19
					}
					position++
				}
			l25:
				{
					position27, tokenIndex27, depth27 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l28
					}
					position++
					goto l27
				l28:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					if buffer[position] != rune('E') {
						goto l19
					}
					position++
				}
			l27:
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l30
					}
					position++
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if buffer[position] != rune('C') {
						goto l19
					}
					position++
				}
			l29:
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l32
					}
					position++
					goto l31
				l32:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
					if buffer[position] != rune('T') {
						goto l19
					}
					position++
				}
			l31:
				if !_rules[rulesp]() {
					goto l19
				}
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					if !_rules[ruleEmitter]() {
						goto l33
					}
					goto l34
				l33:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
				}
			l34:
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleProjections]() {
					goto l19
				}
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleDefWindowedFrom]() {
					goto l19
				}
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleFilter]() {
					goto l19
				}
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleGrouping]() {
					goto l19
				}
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleHaving]() {
					goto l19
				}
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleAction0]() {
					goto l19
				}
				depth--
				add(ruleSelectStmt, position20)
			}
			return true
		l19:
			position, tokenIndex, depth = position19, tokenIndex19, depth19
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l38
					}
					position++
					goto l37
				l38:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if buffer[position] != rune('C') {
						goto l35
					}
					position++
				}
			l37:
				{
					position39, tokenIndex39, depth39 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l40
					}
					position++
					goto l39
				l40:
					position, tokenIndex, depth = position39, tokenIndex39, depth39
					if buffer[position] != rune('R') {
						goto l35
					}
					position++
				}
			l39:
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l42
					}
					position++
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if buffer[position] != rune('E') {
						goto l35
					}
					position++
				}
			l41:
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l44
					}
					position++
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if buffer[position] != rune('A') {
						goto l35
					}
					position++
				}
			l43:
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l46
					}
					position++
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if buffer[position] != rune('T') {
						goto l35
					}
					position++
				}
			l45:
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l48
					}
					position++
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					if buffer[position] != rune('E') {
						goto l35
					}
					position++
				}
			l47:
				if !_rules[rulesp]() {
					goto l35
				}
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l50
					}
					position++
					goto l49
				l50:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
					if buffer[position] != rune('S') {
						goto l35
					}
					position++
				}
			l49:
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l52
					}
					position++
					goto l51
				l52:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
					if buffer[position] != rune('T') {
						goto l35
					}
					position++
				}
			l51:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					if buffer[position] != rune('R') {
						goto l35
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
					if buffer[position] != rune('E') {
						goto l35
					}
					position++
				}
			l55:
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
						goto l35
					}
					position++
				}
			l57:
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l60
					}
					position++
					goto l59
				l60:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
					if buffer[position] != rune('M') {
						goto l35
					}
					position++
				}
			l59:
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l35
				}
				if !_rules[rulesp]() {
					goto l35
				}
				{
					position61, tokenIndex61, depth61 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l62
					}
					position++
					goto l61
				l62:
					position, tokenIndex, depth = position61, tokenIndex61, depth61
					if buffer[position] != rune('A') {
						goto l35
					}
					position++
				}
			l61:
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
						goto l35
					}
					position++
				}
			l63:
				if !_rules[rulesp]() {
					goto l35
				}
				{
					position65, tokenIndex65, depth65 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l66
					}
					position++
					goto l65
				l66:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
					if buffer[position] != rune('S') {
						goto l35
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
						goto l35
					}
					position++
				}
			l67:
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l70
					}
					position++
					goto l69
				l70:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
					if buffer[position] != rune('L') {
						goto l35
					}
					position++
				}
			l69:
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l72
					}
					position++
					goto l71
				l72:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					if buffer[position] != rune('E') {
						goto l35
					}
					position++
				}
			l71:
				{
					position73, tokenIndex73, depth73 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l74
					}
					position++
					goto l73
				l74:
					position, tokenIndex, depth = position73, tokenIndex73, depth73
					if buffer[position] != rune('C') {
						goto l35
					}
					position++
				}
			l73:
				{
					position75, tokenIndex75, depth75 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l76
					}
					position++
					goto l75
				l76:
					position, tokenIndex, depth = position75, tokenIndex75, depth75
					if buffer[position] != rune('T') {
						goto l35
					}
					position++
				}
			l75:
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleEmitter]() {
					goto l35
				}
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleProjections]() {
					goto l35
				}
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleWindowedFrom]() {
					goto l35
				}
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleFilter]() {
					goto l35
				}
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleGrouping]() {
					goto l35
				}
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleHaving]() {
					goto l35
				}
				if !_rules[rulesp]() {
					goto l35
				}
				if !_rules[ruleAction1]() {
					goto l35
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position77, tokenIndex77, depth77 := position, tokenIndex, depth
			{
				position78 := position
				depth++
				{
					position79, tokenIndex79, depth79 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l80
					}
					position++
					goto l79
				l80:
					position, tokenIndex, depth = position79, tokenIndex79, depth79
					if buffer[position] != rune('C') {
						goto l77
					}
					position++
				}
			l79:
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l82
					}
					position++
					goto l81
				l82:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if buffer[position] != rune('R') {
						goto l77
					}
					position++
				}
			l81:
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l84
					}
					position++
					goto l83
				l84:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
					if buffer[position] != rune('E') {
						goto l77
					}
					position++
				}
			l83:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l86
					}
					position++
					goto l85
				l86:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
					if buffer[position] != rune('A') {
						goto l77
					}
					position++
				}
			l85:
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l88
					}
					position++
					goto l87
				l88:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
					if buffer[position] != rune('T') {
						goto l77
					}
					position++
				}
			l87:
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l90
					}
					position++
					goto l89
				l90:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
					if buffer[position] != rune('E') {
						goto l77
					}
					position++
				}
			l89:
				if !_rules[rulesp]() {
					goto l77
				}
				if !_rules[rulePausedOpt]() {
					goto l77
				}
				if !_rules[rulesp]() {
					goto l77
				}
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l92
					}
					position++
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					if buffer[position] != rune('S') {
						goto l77
					}
					position++
				}
			l91:
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('O') {
						goto l77
					}
					position++
				}
			l93:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if buffer[position] != rune('U') {
						goto l77
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('R') {
						goto l77
					}
					position++
				}
			l97:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('C') {
						goto l77
					}
					position++
				}
			l99:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('E') {
						goto l77
					}
					position++
				}
			l101:
				if !_rules[rulesp]() {
					goto l77
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l77
				}
				if !_rules[rulesp]() {
					goto l77
				}
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if buffer[position] != rune('T') {
						goto l77
					}
					position++
				}
			l103:
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('Y') {
						goto l77
					}
					position++
				}
			l105:
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l108
					}
					position++
					goto l107
				l108:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
					if buffer[position] != rune('P') {
						goto l77
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('E') {
						goto l77
					}
					position++
				}
			l109:
				if !_rules[rulesp]() {
					goto l77
				}
				if !_rules[ruleSourceSinkType]() {
					goto l77
				}
				if !_rules[rulesp]() {
					goto l77
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l77
				}
				if !_rules[ruleAction2]() {
					goto l77
				}
				depth--
				add(ruleCreateSourceStmt, position78)
			}
			return true
		l77:
			position, tokenIndex, depth = position77, tokenIndex77, depth77
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position111, tokenIndex111, depth111 := position, tokenIndex, depth
			{
				position112 := position
				depth++
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l114
					}
					position++
					goto l113
				l114:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
					if buffer[position] != rune('C') {
						goto l111
					}
					position++
				}
			l113:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('R') {
						goto l111
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('E') {
						goto l111
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('A') {
						goto l111
					}
					position++
				}
			l119:
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if buffer[position] != rune('T') {
						goto l111
					}
					position++
				}
			l121:
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('E') {
						goto l111
					}
					position++
				}
			l123:
				if !_rules[rulesp]() {
					goto l111
				}
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('S') {
						goto l111
					}
					position++
				}
			l125:
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('I') {
						goto l111
					}
					position++
				}
			l127:
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('N') {
						goto l111
					}
					position++
				}
			l129:
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('K') {
						goto l111
					}
					position++
				}
			l131:
				if !_rules[rulesp]() {
					goto l111
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l111
				}
				if !_rules[rulesp]() {
					goto l111
				}
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('T') {
						goto l111
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('Y') {
						goto l111
					}
					position++
				}
			l135:
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l138
					}
					position++
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if buffer[position] != rune('P') {
						goto l111
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('E') {
						goto l111
					}
					position++
				}
			l139:
				if !_rules[rulesp]() {
					goto l111
				}
				if !_rules[ruleSourceSinkType]() {
					goto l111
				}
				if !_rules[rulesp]() {
					goto l111
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l111
				}
				if !_rules[ruleAction3]() {
					goto l111
				}
				depth--
				add(ruleCreateSinkStmt, position112)
			}
			return true
		l111:
			position, tokenIndex, depth = position111, tokenIndex111, depth111
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position141, tokenIndex141, depth141 := position, tokenIndex, depth
			{
				position142 := position
				depth++
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					if buffer[position] != rune('C') {
						goto l141
					}
					position++
				}
			l143:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('R') {
						goto l141
					}
					position++
				}
			l145:
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('E') {
						goto l141
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('A') {
						goto l141
					}
					position++
				}
			l149:
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('T') {
						goto l141
					}
					position++
				}
			l151:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('E') {
						goto l141
					}
					position++
				}
			l153:
				if !_rules[rulesp]() {
					goto l141
				}
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('S') {
						goto l141
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
						goto l141
					}
					position++
				}
			l157:
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l160
					}
					position++
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if buffer[position] != rune('A') {
						goto l141
					}
					position++
				}
			l159:
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
						goto l141
					}
					position++
				}
			l161:
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l164
					}
					position++
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if buffer[position] != rune('E') {
						goto l141
					}
					position++
				}
			l163:
				if !_rules[rulesp]() {
					goto l141
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l141
				}
				if !_rules[rulesp]() {
					goto l141
				}
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('T') {
						goto l141
					}
					position++
				}
			l165:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('Y') {
						goto l141
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('P') {
						goto l141
					}
					position++
				}
			l169:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('E') {
						goto l141
					}
					position++
				}
			l171:
				if !_rules[rulesp]() {
					goto l141
				}
				if !_rules[ruleSourceSinkType]() {
					goto l141
				}
				if !_rules[rulesp]() {
					goto l141
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l141
				}
				if !_rules[ruleAction4]() {
					goto l141
				}
				depth--
				add(ruleCreateStateStmt, position142)
			}
			return true
		l141:
			position, tokenIndex, depth = position141, tokenIndex141, depth141
			return false
		},
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action5)> */
		func() bool {
			position173, tokenIndex173, depth173 := position, tokenIndex, depth
			{
				position174 := position
				depth++
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('I') {
						goto l173
					}
					position++
				}
			l175:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('N') {
						goto l173
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('S') {
						goto l173
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('E') {
						goto l173
					}
					position++
				}
			l181:
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('R') {
						goto l173
					}
					position++
				}
			l183:
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('T') {
						goto l173
					}
					position++
				}
			l185:
				if !_rules[rulesp]() {
					goto l173
				}
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if buffer[position] != rune('I') {
						goto l173
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('N') {
						goto l173
					}
					position++
				}
			l189:
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('T') {
						goto l173
					}
					position++
				}
			l191:
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('O') {
						goto l173
					}
					position++
				}
			l193:
				if !_rules[rulesp]() {
					goto l173
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l173
				}
				if !_rules[rulesp]() {
					goto l173
				}
				if !_rules[ruleSelectStmt]() {
					goto l173
				}
				if !_rules[ruleAction5]() {
					goto l173
				}
				depth--
				add(ruleInsertIntoSelectStmt, position174)
			}
			return true
		l173:
			position, tokenIndex, depth = position173, tokenIndex173, depth173
			return false
		},
		/* 8 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action6)> */
		func() bool {
			position195, tokenIndex195, depth195 := position, tokenIndex, depth
			{
				position196 := position
				depth++
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('I') {
						goto l195
					}
					position++
				}
			l197:
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('N') {
						goto l195
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
						goto l195
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
						goto l195
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
						goto l195
					}
					position++
				}
			l205:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('T') {
						goto l195
					}
					position++
				}
			l207:
				if !_rules[rulesp]() {
					goto l195
				}
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('I') {
						goto l195
					}
					position++
				}
			l209:
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('N') {
						goto l195
					}
					position++
				}
			l211:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('T') {
						goto l195
					}
					position++
				}
			l213:
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('O') {
						goto l195
					}
					position++
				}
			l215:
				if !_rules[rulesp]() {
					goto l195
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l195
				}
				if !_rules[rulesp]() {
					goto l195
				}
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('F') {
						goto l195
					}
					position++
				}
			l217:
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
						goto l195
					}
					position++
				}
			l219:
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if buffer[position] != rune('O') {
						goto l195
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('M') {
						goto l195
					}
					position++
				}
			l223:
				if !_rules[rulesp]() {
					goto l195
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l195
				}
				if !_rules[ruleAction6]() {
					goto l195
				}
				depth--
				add(ruleInsertIntoFromStmt, position196)
			}
			return true
		l195:
			position, tokenIndex, depth = position195, tokenIndex195, depth195
			return false
		},
		/* 9 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action7)> */
		func() bool {
			position225, tokenIndex225, depth225 := position, tokenIndex, depth
			{
				position226 := position
				depth++
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l228
					}
					position++
					goto l227
				l228:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
					if buffer[position] != rune('P') {
						goto l225
					}
					position++
				}
			l227:
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('A') {
						goto l225
					}
					position++
				}
			l229:
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if buffer[position] != rune('U') {
						goto l225
					}
					position++
				}
			l231:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('S') {
						goto l225
					}
					position++
				}
			l233:
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l236
					}
					position++
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if buffer[position] != rune('E') {
						goto l225
					}
					position++
				}
			l235:
				if !_rules[rulesp]() {
					goto l225
				}
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if buffer[position] != rune('S') {
						goto l225
					}
					position++
				}
			l237:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					if buffer[position] != rune('O') {
						goto l225
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('U') {
						goto l225
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('R') {
						goto l225
					}
					position++
				}
			l243:
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l246
					}
					position++
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if buffer[position] != rune('C') {
						goto l225
					}
					position++
				}
			l245:
				{
					position247, tokenIndex247, depth247 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l248
					}
					position++
					goto l247
				l248:
					position, tokenIndex, depth = position247, tokenIndex247, depth247
					if buffer[position] != rune('E') {
						goto l225
					}
					position++
				}
			l247:
				if !_rules[rulesp]() {
					goto l225
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l225
				}
				if !_rules[ruleAction7]() {
					goto l225
				}
				depth--
				add(rulePauseSourceStmt, position226)
			}
			return true
		l225:
			position, tokenIndex, depth = position225, tokenIndex225, depth225
			return false
		},
		/* 10 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action8)> */
		func() bool {
			position249, tokenIndex249, depth249 := position, tokenIndex, depth
			{
				position250 := position
				depth++
				{
					position251, tokenIndex251, depth251 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l252
					}
					position++
					goto l251
				l252:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('R') {
						goto l249
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('E') {
						goto l249
					}
					position++
				}
			l253:
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l256
					}
					position++
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('S') {
						goto l249
					}
					position++
				}
			l255:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('U') {
						goto l249
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('M') {
						goto l249
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
						goto l249
					}
					position++
				}
			l261:
				if !_rules[rulesp]() {
					goto l249
				}
				{
					position263, tokenIndex263, depth263 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if buffer[position] != rune('S') {
						goto l249
					}
					position++
				}
			l263:
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('O') {
						goto l249
					}
					position++
				}
			l265:
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l268
					}
					position++
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if buffer[position] != rune('U') {
						goto l249
					}
					position++
				}
			l267:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('R') {
						goto l249
					}
					position++
				}
			l269:
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('C') {
						goto l249
					}
					position++
				}
			l271:
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('E') {
						goto l249
					}
					position++
				}
			l273:
				if !_rules[rulesp]() {
					goto l249
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l249
				}
				if !_rules[ruleAction8]() {
					goto l249
				}
				depth--
				add(ruleResumeSourceStmt, position250)
			}
			return true
		l249:
			position, tokenIndex, depth = position249, tokenIndex249, depth249
			return false
		},
		/* 11 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action9)> */
		func() bool {
			position275, tokenIndex275, depth275 := position, tokenIndex, depth
			{
				position276 := position
				depth++
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('R') {
						goto l275
					}
					position++
				}
			l277:
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('E') {
						goto l275
					}
					position++
				}
			l279:
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('W') {
						goto l275
					}
					position++
				}
			l281:
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
						goto l275
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
						goto l275
					}
					position++
				}
			l285:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('D') {
						goto l275
					}
					position++
				}
			l287:
				if !_rules[rulesp]() {
					goto l275
				}
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('S') {
						goto l275
					}
					position++
				}
			l289:
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('O') {
						goto l275
					}
					position++
				}
			l291:
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					if buffer[position] != rune('U') {
						goto l275
					}
					position++
				}
			l293:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('R') {
						goto l275
					}
					position++
				}
			l295:
				{
					position297, tokenIndex297, depth297 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l298
					}
					position++
					goto l297
				l298:
					position, tokenIndex, depth = position297, tokenIndex297, depth297
					if buffer[position] != rune('C') {
						goto l275
					}
					position++
				}
			l297:
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l300
					}
					position++
					goto l299
				l300:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
					if buffer[position] != rune('E') {
						goto l275
					}
					position++
				}
			l299:
				if !_rules[rulesp]() {
					goto l275
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l275
				}
				if !_rules[ruleAction9]() {
					goto l275
				}
				depth--
				add(ruleRewindSourceStmt, position276)
			}
			return true
		l275:
			position, tokenIndex, depth = position275, tokenIndex275, depth275
			return false
		},
		/* 12 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action10)> */
		func() bool {
			position301, tokenIndex301, depth301 := position, tokenIndex, depth
			{
				position302 := position
				depth++
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l304
					}
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if !_rules[ruleDSTREAM]() {
						goto l305
					}
					goto l303
				l305:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if !_rules[ruleRSTREAM]() {
						goto l301
					}
				}
			l303:
				{
					position306 := position
					depth++
					{
						position307, tokenIndex307, depth307 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l307
						}
						if buffer[position] != rune('[') {
							goto l307
						}
						position++
						if !_rules[rulesp]() {
							goto l307
						}
						{
							position309, tokenIndex309, depth309 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l310
							}
							position++
							goto l309
						l310:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if buffer[position] != rune('E') {
								goto l307
							}
							position++
						}
					l309:
						{
							position311, tokenIndex311, depth311 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l312
							}
							position++
							goto l311
						l312:
							position, tokenIndex, depth = position311, tokenIndex311, depth311
							if buffer[position] != rune('V') {
								goto l307
							}
							position++
						}
					l311:
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l314
							}
							position++
							goto l313
						l314:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
							if buffer[position] != rune('E') {
								goto l307
							}
							position++
						}
					l313:
						{
							position315, tokenIndex315, depth315 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l316
							}
							position++
							goto l315
						l316:
							position, tokenIndex, depth = position315, tokenIndex315, depth315
							if buffer[position] != rune('R') {
								goto l307
							}
							position++
						}
					l315:
						{
							position317, tokenIndex317, depth317 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l318
							}
							position++
							goto l317
						l318:
							position, tokenIndex, depth = position317, tokenIndex317, depth317
							if buffer[position] != rune('Y') {
								goto l307
							}
							position++
						}
					l317:
						if !_rules[rulesp]() {
							goto l307
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l307
						}
						if !_rules[rulesp]() {
							goto l307
						}
						if buffer[position] != rune(']') {
							goto l307
						}
						position++
						goto l308
					l307:
						position, tokenIndex, depth = position307, tokenIndex307, depth307
					}
				l308:
					depth--
					add(rulePegText, position306)
				}
				if !_rules[ruleAction10]() {
					goto l301
				}
				depth--
				add(ruleEmitter, position302)
			}
			return true
		l301:
			position, tokenIndex, depth = position301, tokenIndex301, depth301
			return false
		},
		/* 13 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l322
					}
				l323:
					{
						position324, tokenIndex324, depth324 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l324
						}
						if buffer[position] != rune(',') {
							goto l324
						}
						position++
						if !_rules[rulesp]() {
							goto l324
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l324
						}
						goto l323
					l324:
						position, tokenIndex, depth = position324, tokenIndex324, depth324
					}
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if !_rules[ruleTimeEmitterInterval]() {
						goto l325
					}
					goto l321
				l325:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if !_rules[ruleTupleEmitterInterval]() {
						goto l319
					}
				}
			l321:
				depth--
				add(ruleEmitterIntervals, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 14 TimeEmitterInterval <- <(<TimeInterval> Action11)> */
		func() bool {
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				{
					position328 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l326
					}
					depth--
					add(rulePegText, position328)
				}
				if !_rules[ruleAction11]() {
					goto l326
				}
				depth--
				add(ruleTimeEmitterInterval, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 15 TupleEmitterInterval <- <(<TuplesInterval> Action12)> */
		func() bool {
			position329, tokenIndex329, depth329 := position, tokenIndex, depth
			{
				position330 := position
				depth++
				{
					position331 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l329
					}
					depth--
					add(rulePegText, position331)
				}
				if !_rules[ruleAction12]() {
					goto l329
				}
				depth--
				add(ruleTupleEmitterInterval, position330)
			}
			return true
		l329:
			position, tokenIndex, depth = position329, tokenIndex329, depth329
			return false
		},
		/* 16 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action13)> */
		func() bool {
			position332, tokenIndex332, depth332 := position, tokenIndex, depth
			{
				position333 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l332
				}
				if !_rules[rulesp]() {
					goto l332
				}
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if buffer[position] != rune('I') {
						goto l332
					}
					position++
				}
			l334:
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if buffer[position] != rune('N') {
						goto l332
					}
					position++
				}
			l336:
				if !_rules[rulesp]() {
					goto l332
				}
				if !_rules[ruleStream]() {
					goto l332
				}
				if !_rules[ruleAction13]() {
					goto l332
				}
				depth--
				add(ruleTupleEmitterFromInterval, position333)
			}
			return true
		l332:
			position, tokenIndex, depth = position332, tokenIndex332, depth332
			return false
		},
		/* 17 Projections <- <(<(Projection sp (',' sp Projection)*)> Action14)> */
		func() bool {
			position338, tokenIndex338, depth338 := position, tokenIndex, depth
			{
				position339 := position
				depth++
				{
					position340 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l338
					}
					if !_rules[rulesp]() {
						goto l338
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
						if !_rules[ruleProjection]() {
							goto l342
						}
						goto l341
					l342:
						position, tokenIndex, depth = position342, tokenIndex342, depth342
					}
					depth--
					add(rulePegText, position340)
				}
				if !_rules[ruleAction14]() {
					goto l338
				}
				depth--
				add(ruleProjections, position339)
			}
			return true
		l338:
			position, tokenIndex, depth = position338, tokenIndex338, depth338
			return false
		},
		/* 18 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position343, tokenIndex343, depth343 := position, tokenIndex, depth
			{
				position344 := position
				depth++
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l346
					}
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if !_rules[ruleExpression]() {
						goto l347
					}
					goto l345
				l347:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if !_rules[ruleWildcard]() {
						goto l343
					}
				}
			l345:
				depth--
				add(ruleProjection, position344)
			}
			return true
		l343:
			position, tokenIndex, depth = position343, tokenIndex343, depth343
			return false
		},
		/* 19 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action15)> */
		func() bool {
			position348, tokenIndex348, depth348 := position, tokenIndex, depth
			{
				position349 := position
				depth++
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l351
					}
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if !_rules[ruleWildcard]() {
						goto l348
					}
				}
			l350:
				if !_rules[rulesp]() {
					goto l348
				}
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('A') {
						goto l348
					}
					position++
				}
			l352:
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
					if buffer[position] != rune('S') {
						goto l348
					}
					position++
				}
			l354:
				if !_rules[rulesp]() {
					goto l348
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l348
				}
				if !_rules[ruleAction15]() {
					goto l348
				}
				depth--
				add(ruleAliasExpression, position349)
			}
			return true
		l348:
			position, tokenIndex, depth = position348, tokenIndex348, depth348
			return false
		},
		/* 20 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action16)> */
		func() bool {
			position356, tokenIndex356, depth356 := position, tokenIndex, depth
			{
				position357 := position
				depth++
				{
					position358 := position
					depth++
					{
						position359, tokenIndex359, depth359 := position, tokenIndex, depth
						{
							position361, tokenIndex361, depth361 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l362
							}
							position++
							goto l361
						l362:
							position, tokenIndex, depth = position361, tokenIndex361, depth361
							if buffer[position] != rune('F') {
								goto l359
							}
							position++
						}
					l361:
						{
							position363, tokenIndex363, depth363 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l364
							}
							position++
							goto l363
						l364:
							position, tokenIndex, depth = position363, tokenIndex363, depth363
							if buffer[position] != rune('R') {
								goto l359
							}
							position++
						}
					l363:
						{
							position365, tokenIndex365, depth365 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l366
							}
							position++
							goto l365
						l366:
							position, tokenIndex, depth = position365, tokenIndex365, depth365
							if buffer[position] != rune('O') {
								goto l359
							}
							position++
						}
					l365:
						{
							position367, tokenIndex367, depth367 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l368
							}
							position++
							goto l367
						l368:
							position, tokenIndex, depth = position367, tokenIndex367, depth367
							if buffer[position] != rune('M') {
								goto l359
							}
							position++
						}
					l367:
						if !_rules[rulesp]() {
							goto l359
						}
						if !_rules[ruleRelations]() {
							goto l359
						}
						if !_rules[rulesp]() {
							goto l359
						}
						goto l360
					l359:
						position, tokenIndex, depth = position359, tokenIndex359, depth359
					}
				l360:
					depth--
					add(rulePegText, position358)
				}
				if !_rules[ruleAction16]() {
					goto l356
				}
				depth--
				add(ruleWindowedFrom, position357)
			}
			return true
		l356:
			position, tokenIndex, depth = position356, tokenIndex356, depth356
			return false
		},
		/* 21 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action17)> */
		func() bool {
			position369, tokenIndex369, depth369 := position, tokenIndex, depth
			{
				position370 := position
				depth++
				{
					position371 := position
					depth++
					{
						position372, tokenIndex372, depth372 := position, tokenIndex, depth
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l375
							}
							position++
							goto l374
						l375:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
							if buffer[position] != rune('F') {
								goto l372
							}
							position++
						}
					l374:
						{
							position376, tokenIndex376, depth376 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l377
							}
							position++
							goto l376
						l377:
							position, tokenIndex, depth = position376, tokenIndex376, depth376
							if buffer[position] != rune('R') {
								goto l372
							}
							position++
						}
					l376:
						{
							position378, tokenIndex378, depth378 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l379
							}
							position++
							goto l378
						l379:
							position, tokenIndex, depth = position378, tokenIndex378, depth378
							if buffer[position] != rune('O') {
								goto l372
							}
							position++
						}
					l378:
						{
							position380, tokenIndex380, depth380 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l381
							}
							position++
							goto l380
						l381:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
							if buffer[position] != rune('M') {
								goto l372
							}
							position++
						}
					l380:
						if !_rules[rulesp]() {
							goto l372
						}
						if !_rules[ruleDefRelations]() {
							goto l372
						}
						if !_rules[rulesp]() {
							goto l372
						}
						goto l373
					l372:
						position, tokenIndex, depth = position372, tokenIndex372, depth372
					}
				l373:
					depth--
					add(rulePegText, position371)
				}
				if !_rules[ruleAction17]() {
					goto l369
				}
				depth--
				add(ruleDefWindowedFrom, position370)
			}
			return true
		l369:
			position, tokenIndex, depth = position369, tokenIndex369, depth369
			return false
		},
		/* 22 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position382, tokenIndex382, depth382 := position, tokenIndex, depth
			{
				position383 := position
				depth++
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l385
					}
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if !_rules[ruleTuplesInterval]() {
						goto l382
					}
				}
			l384:
				depth--
				add(ruleInterval, position383)
			}
			return true
		l382:
			position, tokenIndex, depth = position382, tokenIndex382, depth382
			return false
		},
		/* 23 TimeInterval <- <(NumericLiteral sp SECONDS Action18)> */
		func() bool {
			position386, tokenIndex386, depth386 := position, tokenIndex, depth
			{
				position387 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l386
				}
				if !_rules[rulesp]() {
					goto l386
				}
				if !_rules[ruleSECONDS]() {
					goto l386
				}
				if !_rules[ruleAction18]() {
					goto l386
				}
				depth--
				add(ruleTimeInterval, position387)
			}
			return true
		l386:
			position, tokenIndex, depth = position386, tokenIndex386, depth386
			return false
		},
		/* 24 TuplesInterval <- <(NumericLiteral sp TUPLES Action19)> */
		func() bool {
			position388, tokenIndex388, depth388 := position, tokenIndex, depth
			{
				position389 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l388
				}
				if !_rules[rulesp]() {
					goto l388
				}
				if !_rules[ruleTUPLES]() {
					goto l388
				}
				if !_rules[ruleAction19]() {
					goto l388
				}
				depth--
				add(ruleTuplesInterval, position389)
			}
			return true
		l388:
			position, tokenIndex, depth = position388, tokenIndex388, depth388
			return false
		},
		/* 25 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l390
				}
				if !_rules[rulesp]() {
					goto l390
				}
			l392:
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l393
					}
					position++
					if !_rules[rulesp]() {
						goto l393
					}
					if !_rules[ruleRelationLike]() {
						goto l393
					}
					goto l392
				l393:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
				}
				depth--
				add(ruleRelations, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 26 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l394
				}
				if !_rules[rulesp]() {
					goto l394
				}
			l396:
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l397
					}
					position++
					if !_rules[rulesp]() {
						goto l397
					}
					if !_rules[ruleDefRelationLike]() {
						goto l397
					}
					goto l396
				l397:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
				}
				depth--
				add(ruleDefRelations, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 27 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action20)> */
		func() bool {
			position398, tokenIndex398, depth398 := position, tokenIndex, depth
			{
				position399 := position
				depth++
				{
					position400 := position
					depth++
					{
						position401, tokenIndex401, depth401 := position, tokenIndex, depth
						{
							position403, tokenIndex403, depth403 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l404
							}
							position++
							goto l403
						l404:
							position, tokenIndex, depth = position403, tokenIndex403, depth403
							if buffer[position] != rune('W') {
								goto l401
							}
							position++
						}
					l403:
						{
							position405, tokenIndex405, depth405 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l406
							}
							position++
							goto l405
						l406:
							position, tokenIndex, depth = position405, tokenIndex405, depth405
							if buffer[position] != rune('H') {
								goto l401
							}
							position++
						}
					l405:
						{
							position407, tokenIndex407, depth407 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l408
							}
							position++
							goto l407
						l408:
							position, tokenIndex, depth = position407, tokenIndex407, depth407
							if buffer[position] != rune('E') {
								goto l401
							}
							position++
						}
					l407:
						{
							position409, tokenIndex409, depth409 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l410
							}
							position++
							goto l409
						l410:
							position, tokenIndex, depth = position409, tokenIndex409, depth409
							if buffer[position] != rune('R') {
								goto l401
							}
							position++
						}
					l409:
						{
							position411, tokenIndex411, depth411 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l412
							}
							position++
							goto l411
						l412:
							position, tokenIndex, depth = position411, tokenIndex411, depth411
							if buffer[position] != rune('E') {
								goto l401
							}
							position++
						}
					l411:
						if !_rules[rulesp]() {
							goto l401
						}
						if !_rules[ruleExpression]() {
							goto l401
						}
						goto l402
					l401:
						position, tokenIndex, depth = position401, tokenIndex401, depth401
					}
				l402:
					depth--
					add(rulePegText, position400)
				}
				if !_rules[ruleAction20]() {
					goto l398
				}
				depth--
				add(ruleFilter, position399)
			}
			return true
		l398:
			position, tokenIndex, depth = position398, tokenIndex398, depth398
			return false
		},
		/* 28 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action21)> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				{
					position415 := position
					depth++
					{
						position416, tokenIndex416, depth416 := position, tokenIndex, depth
						{
							position418, tokenIndex418, depth418 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l419
							}
							position++
							goto l418
						l419:
							position, tokenIndex, depth = position418, tokenIndex418, depth418
							if buffer[position] != rune('G') {
								goto l416
							}
							position++
						}
					l418:
						{
							position420, tokenIndex420, depth420 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l421
							}
							position++
							goto l420
						l421:
							position, tokenIndex, depth = position420, tokenIndex420, depth420
							if buffer[position] != rune('R') {
								goto l416
							}
							position++
						}
					l420:
						{
							position422, tokenIndex422, depth422 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l423
							}
							position++
							goto l422
						l423:
							position, tokenIndex, depth = position422, tokenIndex422, depth422
							if buffer[position] != rune('O') {
								goto l416
							}
							position++
						}
					l422:
						{
							position424, tokenIndex424, depth424 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l425
							}
							position++
							goto l424
						l425:
							position, tokenIndex, depth = position424, tokenIndex424, depth424
							if buffer[position] != rune('U') {
								goto l416
							}
							position++
						}
					l424:
						{
							position426, tokenIndex426, depth426 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l427
							}
							position++
							goto l426
						l427:
							position, tokenIndex, depth = position426, tokenIndex426, depth426
							if buffer[position] != rune('P') {
								goto l416
							}
							position++
						}
					l426:
						if !_rules[rulesp]() {
							goto l416
						}
						{
							position428, tokenIndex428, depth428 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l429
							}
							position++
							goto l428
						l429:
							position, tokenIndex, depth = position428, tokenIndex428, depth428
							if buffer[position] != rune('B') {
								goto l416
							}
							position++
						}
					l428:
						{
							position430, tokenIndex430, depth430 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l431
							}
							position++
							goto l430
						l431:
							position, tokenIndex, depth = position430, tokenIndex430, depth430
							if buffer[position] != rune('Y') {
								goto l416
							}
							position++
						}
					l430:
						if !_rules[rulesp]() {
							goto l416
						}
						if !_rules[ruleGroupList]() {
							goto l416
						}
						goto l417
					l416:
						position, tokenIndex, depth = position416, tokenIndex416, depth416
					}
				l417:
					depth--
					add(rulePegText, position415)
				}
				if !_rules[ruleAction21]() {
					goto l413
				}
				depth--
				add(ruleGrouping, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 29 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l432
				}
				if !_rules[rulesp]() {
					goto l432
				}
			l434:
				{
					position435, tokenIndex435, depth435 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l435
					}
					position++
					if !_rules[rulesp]() {
						goto l435
					}
					if !_rules[ruleExpression]() {
						goto l435
					}
					goto l434
				l435:
					position, tokenIndex, depth = position435, tokenIndex435, depth435
				}
				depth--
				add(ruleGroupList, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 30 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action22)> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				{
					position438 := position
					depth++
					{
						position439, tokenIndex439, depth439 := position, tokenIndex, depth
						{
							position441, tokenIndex441, depth441 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l442
							}
							position++
							goto l441
						l442:
							position, tokenIndex, depth = position441, tokenIndex441, depth441
							if buffer[position] != rune('H') {
								goto l439
							}
							position++
						}
					l441:
						{
							position443, tokenIndex443, depth443 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l444
							}
							position++
							goto l443
						l444:
							position, tokenIndex, depth = position443, tokenIndex443, depth443
							if buffer[position] != rune('A') {
								goto l439
							}
							position++
						}
					l443:
						{
							position445, tokenIndex445, depth445 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l446
							}
							position++
							goto l445
						l446:
							position, tokenIndex, depth = position445, tokenIndex445, depth445
							if buffer[position] != rune('V') {
								goto l439
							}
							position++
						}
					l445:
						{
							position447, tokenIndex447, depth447 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l448
							}
							position++
							goto l447
						l448:
							position, tokenIndex, depth = position447, tokenIndex447, depth447
							if buffer[position] != rune('I') {
								goto l439
							}
							position++
						}
					l447:
						{
							position449, tokenIndex449, depth449 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l450
							}
							position++
							goto l449
						l450:
							position, tokenIndex, depth = position449, tokenIndex449, depth449
							if buffer[position] != rune('N') {
								goto l439
							}
							position++
						}
					l449:
						{
							position451, tokenIndex451, depth451 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l452
							}
							position++
							goto l451
						l452:
							position, tokenIndex, depth = position451, tokenIndex451, depth451
							if buffer[position] != rune('G') {
								goto l439
							}
							position++
						}
					l451:
						if !_rules[rulesp]() {
							goto l439
						}
						if !_rules[ruleExpression]() {
							goto l439
						}
						goto l440
					l439:
						position, tokenIndex, depth = position439, tokenIndex439, depth439
					}
				l440:
					depth--
					add(rulePegText, position438)
				}
				if !_rules[ruleAction22]() {
					goto l436
				}
				depth--
				add(ruleHaving, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
			return false
		},
		/* 31 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action23))> */
		func() bool {
			position453, tokenIndex453, depth453 := position, tokenIndex, depth
			{
				position454 := position
				depth++
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l456
					}
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if !_rules[ruleStreamWindow]() {
						goto l453
					}
					if !_rules[ruleAction23]() {
						goto l453
					}
				}
			l455:
				depth--
				add(ruleRelationLike, position454)
			}
			return true
		l453:
			position, tokenIndex, depth = position453, tokenIndex453, depth453
			return false
		},
		/* 32 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action24))> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l460
					}
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleDefStreamWindow]() {
						goto l457
					}
					if !_rules[ruleAction24]() {
						goto l457
					}
				}
			l459:
				depth--
				add(ruleDefRelationLike, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 33 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action25)> */
		func() bool {
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l461
				}
				if !_rules[rulesp]() {
					goto l461
				}
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('A') {
						goto l461
					}
					position++
				}
			l463:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l466
					}
					position++
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					if buffer[position] != rune('S') {
						goto l461
					}
					position++
				}
			l465:
				if !_rules[rulesp]() {
					goto l461
				}
				if !_rules[ruleIdentifier]() {
					goto l461
				}
				if !_rules[ruleAction25]() {
					goto l461
				}
				depth--
				add(ruleAliasedStreamWindow, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 34 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action26)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l467
				}
				if !_rules[rulesp]() {
					goto l467
				}
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l470
					}
					position++
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if buffer[position] != rune('A') {
						goto l467
					}
					position++
				}
			l469:
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l472
					}
					position++
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if buffer[position] != rune('S') {
						goto l467
					}
					position++
				}
			l471:
				if !_rules[rulesp]() {
					goto l467
				}
				if !_rules[ruleIdentifier]() {
					goto l467
				}
				if !_rules[ruleAction26]() {
					goto l467
				}
				depth--
				add(ruleDefAliasedStreamWindow, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 35 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action27)> */
		func() bool {
			position473, tokenIndex473, depth473 := position, tokenIndex, depth
			{
				position474 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l473
				}
				if !_rules[rulesp]() {
					goto l473
				}
				if buffer[position] != rune('[') {
					goto l473
				}
				position++
				if !_rules[rulesp]() {
					goto l473
				}
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l476
					}
					position++
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if buffer[position] != rune('R') {
						goto l473
					}
					position++
				}
			l475:
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l478
					}
					position++
					goto l477
				l478:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
					if buffer[position] != rune('A') {
						goto l473
					}
					position++
				}
			l477:
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l480
					}
					position++
					goto l479
				l480:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if buffer[position] != rune('N') {
						goto l473
					}
					position++
				}
			l479:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if buffer[position] != rune('G') {
						goto l473
					}
					position++
				}
			l481:
				{
					position483, tokenIndex483, depth483 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l484
					}
					position++
					goto l483
				l484:
					position, tokenIndex, depth = position483, tokenIndex483, depth483
					if buffer[position] != rune('E') {
						goto l473
					}
					position++
				}
			l483:
				if !_rules[rulesp]() {
					goto l473
				}
				if !_rules[ruleInterval]() {
					goto l473
				}
				if !_rules[rulesp]() {
					goto l473
				}
				if buffer[position] != rune(']') {
					goto l473
				}
				position++
				if !_rules[ruleAction27]() {
					goto l473
				}
				depth--
				add(ruleStreamWindow, position474)
			}
			return true
		l473:
			position, tokenIndex, depth = position473, tokenIndex473, depth473
			return false
		},
		/* 36 DefStreamWindow <- <(StreamLike (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action28)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l485
				}
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l487
					}
					if buffer[position] != rune('[') {
						goto l487
					}
					position++
					if !_rules[rulesp]() {
						goto l487
					}
					{
						position489, tokenIndex489, depth489 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l490
						}
						position++
						goto l489
					l490:
						position, tokenIndex, depth = position489, tokenIndex489, depth489
						if buffer[position] != rune('R') {
							goto l487
						}
						position++
					}
				l489:
					{
						position491, tokenIndex491, depth491 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l492
						}
						position++
						goto l491
					l492:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
						if buffer[position] != rune('A') {
							goto l487
						}
						position++
					}
				l491:
					{
						position493, tokenIndex493, depth493 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l494
						}
						position++
						goto l493
					l494:
						position, tokenIndex, depth = position493, tokenIndex493, depth493
						if buffer[position] != rune('N') {
							goto l487
						}
						position++
					}
				l493:
					{
						position495, tokenIndex495, depth495 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l496
						}
						position++
						goto l495
					l496:
						position, tokenIndex, depth = position495, tokenIndex495, depth495
						if buffer[position] != rune('G') {
							goto l487
						}
						position++
					}
				l495:
					{
						position497, tokenIndex497, depth497 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l498
						}
						position++
						goto l497
					l498:
						position, tokenIndex, depth = position497, tokenIndex497, depth497
						if buffer[position] != rune('E') {
							goto l487
						}
						position++
					}
				l497:
					if !_rules[rulesp]() {
						goto l487
					}
					if !_rules[ruleInterval]() {
						goto l487
					}
					if !_rules[rulesp]() {
						goto l487
					}
					if buffer[position] != rune(']') {
						goto l487
					}
					position++
					goto l488
				l487:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
				}
			l488:
				if !_rules[ruleAction28]() {
					goto l485
				}
				depth--
				add(ruleDefStreamWindow, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 37 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				{
					position501, tokenIndex501, depth501 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l502
					}
					goto l501
				l502:
					position, tokenIndex, depth = position501, tokenIndex501, depth501
					if !_rules[ruleStream]() {
						goto l499
					}
				}
			l501:
				depth--
				add(ruleStreamLike, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 38 UDSFFuncApp <- <(FuncApp Action29)> */
		func() bool {
			position503, tokenIndex503, depth503 := position, tokenIndex, depth
			{
				position504 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l503
				}
				if !_rules[ruleAction29]() {
					goto l503
				}
				depth--
				add(ruleUDSFFuncApp, position504)
			}
			return true
		l503:
			position, tokenIndex, depth = position503, tokenIndex503, depth503
			return false
		},
		/* 39 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action30)> */
		func() bool {
			position505, tokenIndex505, depth505 := position, tokenIndex, depth
			{
				position506 := position
				depth++
				{
					position507 := position
					depth++
					{
						position508, tokenIndex508, depth508 := position, tokenIndex, depth
						{
							position510, tokenIndex510, depth510 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l511
							}
							position++
							goto l510
						l511:
							position, tokenIndex, depth = position510, tokenIndex510, depth510
							if buffer[position] != rune('W') {
								goto l508
							}
							position++
						}
					l510:
						{
							position512, tokenIndex512, depth512 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l513
							}
							position++
							goto l512
						l513:
							position, tokenIndex, depth = position512, tokenIndex512, depth512
							if buffer[position] != rune('I') {
								goto l508
							}
							position++
						}
					l512:
						{
							position514, tokenIndex514, depth514 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l515
							}
							position++
							goto l514
						l515:
							position, tokenIndex, depth = position514, tokenIndex514, depth514
							if buffer[position] != rune('T') {
								goto l508
							}
							position++
						}
					l514:
						{
							position516, tokenIndex516, depth516 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l517
							}
							position++
							goto l516
						l517:
							position, tokenIndex, depth = position516, tokenIndex516, depth516
							if buffer[position] != rune('H') {
								goto l508
							}
							position++
						}
					l516:
						if !_rules[rulesp]() {
							goto l508
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l508
						}
						if !_rules[rulesp]() {
							goto l508
						}
					l518:
						{
							position519, tokenIndex519, depth519 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l519
							}
							position++
							if !_rules[rulesp]() {
								goto l519
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l519
							}
							goto l518
						l519:
							position, tokenIndex, depth = position519, tokenIndex519, depth519
						}
						goto l509
					l508:
						position, tokenIndex, depth = position508, tokenIndex508, depth508
					}
				l509:
					depth--
					add(rulePegText, position507)
				}
				if !_rules[ruleAction30]() {
					goto l505
				}
				depth--
				add(ruleSourceSinkSpecs, position506)
			}
			return true
		l505:
			position, tokenIndex, depth = position505, tokenIndex505, depth505
			return false
		},
		/* 40 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action31)> */
		func() bool {
			position520, tokenIndex520, depth520 := position, tokenIndex, depth
			{
				position521 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l520
				}
				if buffer[position] != rune('=') {
					goto l520
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l520
				}
				if !_rules[ruleAction31]() {
					goto l520
				}
				depth--
				add(ruleSourceSinkParam, position521)
			}
			return true
		l520:
			position, tokenIndex, depth = position520, tokenIndex520, depth520
			return false
		},
		/* 41 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position522, tokenIndex522, depth522 := position, tokenIndex, depth
			{
				position523 := position
				depth++
				{
					position524, tokenIndex524, depth524 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l525
					}
					goto l524
				l525:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleLiteral]() {
						goto l522
					}
				}
			l524:
				depth--
				add(ruleSourceSinkParamVal, position523)
			}
			return true
		l522:
			position, tokenIndex, depth = position522, tokenIndex522, depth522
			return false
		},
		/* 42 PausedOpt <- <(<(Paused / Unpaused)?> Action32)> */
		func() bool {
			position526, tokenIndex526, depth526 := position, tokenIndex, depth
			{
				position527 := position
				depth++
				{
					position528 := position
					depth++
					{
						position529, tokenIndex529, depth529 := position, tokenIndex, depth
						{
							position531, tokenIndex531, depth531 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l532
							}
							goto l531
						l532:
							position, tokenIndex, depth = position531, tokenIndex531, depth531
							if !_rules[ruleUnpaused]() {
								goto l529
							}
						}
					l531:
						goto l530
					l529:
						position, tokenIndex, depth = position529, tokenIndex529, depth529
					}
				l530:
					depth--
					add(rulePegText, position528)
				}
				if !_rules[ruleAction32]() {
					goto l526
				}
				depth--
				add(rulePausedOpt, position527)
			}
			return true
		l526:
			position, tokenIndex, depth = position526, tokenIndex526, depth526
			return false
		},
		/* 43 Expression <- <orExpr> */
		func() bool {
			position533, tokenIndex533, depth533 := position, tokenIndex, depth
			{
				position534 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l533
				}
				depth--
				add(ruleExpression, position534)
			}
			return true
		l533:
			position, tokenIndex, depth = position533, tokenIndex533, depth533
			return false
		},
		/* 44 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action33)> */
		func() bool {
			position535, tokenIndex535, depth535 := position, tokenIndex, depth
			{
				position536 := position
				depth++
				{
					position537 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l535
					}
					if !_rules[rulesp]() {
						goto l535
					}
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l538
						}
						if !_rules[rulesp]() {
							goto l538
						}
						if !_rules[ruleandExpr]() {
							goto l538
						}
						goto l539
					l538:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
					}
				l539:
					depth--
					add(rulePegText, position537)
				}
				if !_rules[ruleAction33]() {
					goto l535
				}
				depth--
				add(ruleorExpr, position536)
			}
			return true
		l535:
			position, tokenIndex, depth = position535, tokenIndex535, depth535
			return false
		},
		/* 45 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action34)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l540
					}
					if !_rules[rulesp]() {
						goto l540
					}
					{
						position543, tokenIndex543, depth543 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l543
						}
						if !_rules[rulesp]() {
							goto l543
						}
						if !_rules[rulecomparisonExpr]() {
							goto l543
						}
						goto l544
					l543:
						position, tokenIndex, depth = position543, tokenIndex543, depth543
					}
				l544:
					depth--
					add(rulePegText, position542)
				}
				if !_rules[ruleAction34]() {
					goto l540
				}
				depth--
				add(ruleandExpr, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 46 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action35)> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				{
					position547 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l545
					}
					if !_rules[rulesp]() {
						goto l545
					}
					{
						position548, tokenIndex548, depth548 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l548
						}
						if !_rules[rulesp]() {
							goto l548
						}
						if !_rules[ruleisExpr]() {
							goto l548
						}
						goto l549
					l548:
						position, tokenIndex, depth = position548, tokenIndex548, depth548
					}
				l549:
					depth--
					add(rulePegText, position547)
				}
				if !_rules[ruleAction35]() {
					goto l545
				}
				depth--
				add(rulecomparisonExpr, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 47 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action36)> */
		func() bool {
			position550, tokenIndex550, depth550 := position, tokenIndex, depth
			{
				position551 := position
				depth++
				{
					position552 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l550
					}
					if !_rules[rulesp]() {
						goto l550
					}
					{
						position553, tokenIndex553, depth553 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l553
						}
						if !_rules[rulesp]() {
							goto l553
						}
						if !_rules[ruleNullLiteral]() {
							goto l553
						}
						goto l554
					l553:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
					}
				l554:
					depth--
					add(rulePegText, position552)
				}
				if !_rules[ruleAction36]() {
					goto l550
				}
				depth--
				add(ruleisExpr, position551)
			}
			return true
		l550:
			position, tokenIndex, depth = position550, tokenIndex550, depth550
			return false
		},
		/* 48 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action37)> */
		func() bool {
			position555, tokenIndex555, depth555 := position, tokenIndex, depth
			{
				position556 := position
				depth++
				{
					position557 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l555
					}
					if !_rules[rulesp]() {
						goto l555
					}
					{
						position558, tokenIndex558, depth558 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l558
						}
						if !_rules[rulesp]() {
							goto l558
						}
						if !_rules[ruleproductExpr]() {
							goto l558
						}
						goto l559
					l558:
						position, tokenIndex, depth = position558, tokenIndex558, depth558
					}
				l559:
					depth--
					add(rulePegText, position557)
				}
				if !_rules[ruleAction37]() {
					goto l555
				}
				depth--
				add(ruletermExpr, position556)
			}
			return true
		l555:
			position, tokenIndex, depth = position555, tokenIndex555, depth555
			return false
		},
		/* 49 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action38)> */
		func() bool {
			position560, tokenIndex560, depth560 := position, tokenIndex, depth
			{
				position561 := position
				depth++
				{
					position562 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l560
					}
					if !_rules[rulesp]() {
						goto l560
					}
					{
						position563, tokenIndex563, depth563 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l563
						}
						if !_rules[rulesp]() {
							goto l563
						}
						if !_rules[rulebaseExpr]() {
							goto l563
						}
						goto l564
					l563:
						position, tokenIndex, depth = position563, tokenIndex563, depth563
					}
				l564:
					depth--
					add(rulePegText, position562)
				}
				if !_rules[ruleAction38]() {
					goto l560
				}
				depth--
				add(ruleproductExpr, position561)
			}
			return true
		l560:
			position, tokenIndex, depth = position560, tokenIndex560, depth560
			return false
		},
		/* 50 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567, tokenIndex567, depth567 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l568
					}
					position++
					if !_rules[rulesp]() {
						goto l568
					}
					if !_rules[ruleExpression]() {
						goto l568
					}
					if !_rules[rulesp]() {
						goto l568
					}
					if buffer[position] != rune(')') {
						goto l568
					}
					position++
					goto l567
				l568:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleBooleanLiteral]() {
						goto l569
					}
					goto l567
				l569:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleNullLiteral]() {
						goto l570
					}
					goto l567
				l570:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleFuncApp]() {
						goto l571
					}
					goto l567
				l571:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleRowMeta]() {
						goto l572
					}
					goto l567
				l572:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleRowValue]() {
						goto l573
					}
					goto l567
				l573:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleLiteral]() {
						goto l565
					}
				}
			l567:
				depth--
				add(rulebaseExpr, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 51 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action39)> */
		func() bool {
			position574, tokenIndex574, depth574 := position, tokenIndex, depth
			{
				position575 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l574
				}
				if !_rules[rulesp]() {
					goto l574
				}
				if buffer[position] != rune('(') {
					goto l574
				}
				position++
				if !_rules[rulesp]() {
					goto l574
				}
				if !_rules[ruleFuncParams]() {
					goto l574
				}
				if !_rules[rulesp]() {
					goto l574
				}
				if buffer[position] != rune(')') {
					goto l574
				}
				position++
				if !_rules[ruleAction39]() {
					goto l574
				}
				depth--
				add(ruleFuncApp, position575)
			}
			return true
		l574:
			position, tokenIndex, depth = position574, tokenIndex574, depth574
			return false
		},
		/* 52 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action40)> */
		func() bool {
			position576, tokenIndex576, depth576 := position, tokenIndex, depth
			{
				position577 := position
				depth++
				{
					position578 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l576
					}
					if !_rules[rulesp]() {
						goto l576
					}
				l579:
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l580
						}
						position++
						if !_rules[rulesp]() {
							goto l580
						}
						if !_rules[ruleExpression]() {
							goto l580
						}
						goto l579
					l580:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
					}
					depth--
					add(rulePegText, position578)
				}
				if !_rules[ruleAction40]() {
					goto l576
				}
				depth--
				add(ruleFuncParams, position577)
			}
			return true
		l576:
			position, tokenIndex, depth = position576, tokenIndex576, depth576
			return false
		},
		/* 53 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position581, tokenIndex581, depth581 := position, tokenIndex, depth
			{
				position582 := position
				depth++
				{
					position583, tokenIndex583, depth583 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l584
					}
					goto l583
				l584:
					position, tokenIndex, depth = position583, tokenIndex583, depth583
					if !_rules[ruleNumericLiteral]() {
						goto l585
					}
					goto l583
				l585:
					position, tokenIndex, depth = position583, tokenIndex583, depth583
					if !_rules[ruleStringLiteral]() {
						goto l581
					}
				}
			l583:
				depth--
				add(ruleLiteral, position582)
			}
			return true
		l581:
			position, tokenIndex, depth = position581, tokenIndex581, depth581
			return false
		},
		/* 54 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position586, tokenIndex586, depth586 := position, tokenIndex, depth
			{
				position587 := position
				depth++
				{
					position588, tokenIndex588, depth588 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l589
					}
					goto l588
				l589:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					if !_rules[ruleNotEqual]() {
						goto l590
					}
					goto l588
				l590:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					if !_rules[ruleLessOrEqual]() {
						goto l591
					}
					goto l588
				l591:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					if !_rules[ruleLess]() {
						goto l592
					}
					goto l588
				l592:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					if !_rules[ruleGreaterOrEqual]() {
						goto l593
					}
					goto l588
				l593:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					if !_rules[ruleGreater]() {
						goto l594
					}
					goto l588
				l594:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					if !_rules[ruleNotEqual]() {
						goto l586
					}
				}
			l588:
				depth--
				add(ruleComparisonOp, position587)
			}
			return true
		l586:
			position, tokenIndex, depth = position586, tokenIndex586, depth586
			return false
		},
		/* 55 IsOp <- <(IsNot / Is)> */
		func() bool {
			position595, tokenIndex595, depth595 := position, tokenIndex, depth
			{
				position596 := position
				depth++
				{
					position597, tokenIndex597, depth597 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l598
					}
					goto l597
				l598:
					position, tokenIndex, depth = position597, tokenIndex597, depth597
					if !_rules[ruleIs]() {
						goto l595
					}
				}
			l597:
				depth--
				add(ruleIsOp, position596)
			}
			return true
		l595:
			position, tokenIndex, depth = position595, tokenIndex595, depth595
			return false
		},
		/* 56 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				{
					position601, tokenIndex601, depth601 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l602
					}
					goto l601
				l602:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if !_rules[ruleMinus]() {
						goto l599
					}
				}
			l601:
				depth--
				add(rulePlusMinusOp, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 57 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position603, tokenIndex603, depth603 := position, tokenIndex, depth
			{
				position604 := position
				depth++
				{
					position605, tokenIndex605, depth605 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l606
					}
					goto l605
				l606:
					position, tokenIndex, depth = position605, tokenIndex605, depth605
					if !_rules[ruleDivide]() {
						goto l607
					}
					goto l605
				l607:
					position, tokenIndex, depth = position605, tokenIndex605, depth605
					if !_rules[ruleModulo]() {
						goto l603
					}
				}
			l605:
				depth--
				add(ruleMultDivOp, position604)
			}
			return true
		l603:
			position, tokenIndex, depth = position603, tokenIndex603, depth603
			return false
		},
		/* 58 Stream <- <(<ident> Action41)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					if !_rules[ruleident]() {
						goto l608
					}
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction41]() {
					goto l608
				}
				depth--
				add(ruleStream, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 59 RowMeta <- <RowTimestamp> */
		func() bool {
			position611, tokenIndex611, depth611 := position, tokenIndex, depth
			{
				position612 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l611
				}
				depth--
				add(ruleRowMeta, position612)
			}
			return true
		l611:
			position, tokenIndex, depth = position611, tokenIndex611, depth611
			return false
		},
		/* 60 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action42)> */
		func() bool {
			position613, tokenIndex613, depth613 := position, tokenIndex, depth
			{
				position614 := position
				depth++
				{
					position615 := position
					depth++
					{
						position616, tokenIndex616, depth616 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l616
						}
						if buffer[position] != rune(':') {
							goto l616
						}
						position++
						goto l617
					l616:
						position, tokenIndex, depth = position616, tokenIndex616, depth616
					}
				l617:
					if buffer[position] != rune('t') {
						goto l613
					}
					position++
					if buffer[position] != rune('s') {
						goto l613
					}
					position++
					if buffer[position] != rune('(') {
						goto l613
					}
					position++
					if buffer[position] != rune(')') {
						goto l613
					}
					position++
					depth--
					add(rulePegText, position615)
				}
				if !_rules[ruleAction42]() {
					goto l613
				}
				depth--
				add(ruleRowTimestamp, position614)
			}
			return true
		l613:
			position, tokenIndex, depth = position613, tokenIndex613, depth613
			return false
		},
		/* 61 RowValue <- <(<((ident ':')? jsonPath)> Action43)> */
		func() bool {
			position618, tokenIndex618, depth618 := position, tokenIndex, depth
			{
				position619 := position
				depth++
				{
					position620 := position
					depth++
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l621
						}
						if buffer[position] != rune(':') {
							goto l621
						}
						position++
						goto l622
					l621:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
					}
				l622:
					if !_rules[rulejsonPath]() {
						goto l618
					}
					depth--
					add(rulePegText, position620)
				}
				if !_rules[ruleAction43]() {
					goto l618
				}
				depth--
				add(ruleRowValue, position619)
			}
			return true
		l618:
			position, tokenIndex, depth = position618, tokenIndex618, depth618
			return false
		},
		/* 62 NumericLiteral <- <(<('-'? [0-9]+)> Action44)> */
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
						if buffer[position] != rune('-') {
							goto l626
						}
						position++
						goto l627
					l626:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
					}
				l627:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l623
					}
					position++
				l628:
					{
						position629, tokenIndex629, depth629 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position629, tokenIndex629, depth629
					}
					depth--
					add(rulePegText, position625)
				}
				if !_rules[ruleAction44]() {
					goto l623
				}
				depth--
				add(ruleNumericLiteral, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 63 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action45)> */
		func() bool {
			position630, tokenIndex630, depth630 := position, tokenIndex, depth
			{
				position631 := position
				depth++
				{
					position632 := position
					depth++
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l633
						}
						position++
						goto l634
					l633:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
					}
				l634:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l630
					}
					position++
				l635:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l636
						}
						position++
						goto l635
					l636:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
					}
					if buffer[position] != rune('.') {
						goto l630
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l630
					}
					position++
				l637:
					{
						position638, tokenIndex638, depth638 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l638
						}
						position++
						goto l637
					l638:
						position, tokenIndex, depth = position638, tokenIndex638, depth638
					}
					depth--
					add(rulePegText, position632)
				}
				if !_rules[ruleAction45]() {
					goto l630
				}
				depth--
				add(ruleFloatLiteral, position631)
			}
			return true
		l630:
			position, tokenIndex, depth = position630, tokenIndex630, depth630
			return false
		},
		/* 64 Function <- <(<ident> Action46)> */
		func() bool {
			position639, tokenIndex639, depth639 := position, tokenIndex, depth
			{
				position640 := position
				depth++
				{
					position641 := position
					depth++
					if !_rules[ruleident]() {
						goto l639
					}
					depth--
					add(rulePegText, position641)
				}
				if !_rules[ruleAction46]() {
					goto l639
				}
				depth--
				add(ruleFunction, position640)
			}
			return true
		l639:
			position, tokenIndex, depth = position639, tokenIndex639, depth639
			return false
		},
		/* 65 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action47)> */
		func() bool {
			position642, tokenIndex642, depth642 := position, tokenIndex, depth
			{
				position643 := position
				depth++
				{
					position644 := position
					depth++
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l646
						}
						position++
						goto l645
					l646:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
						if buffer[position] != rune('N') {
							goto l642
						}
						position++
					}
				l645:
					{
						position647, tokenIndex647, depth647 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l648
						}
						position++
						goto l647
					l648:
						position, tokenIndex, depth = position647, tokenIndex647, depth647
						if buffer[position] != rune('U') {
							goto l642
						}
						position++
					}
				l647:
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l650
						}
						position++
						goto l649
					l650:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						if buffer[position] != rune('L') {
							goto l642
						}
						position++
					}
				l649:
					{
						position651, tokenIndex651, depth651 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l652
						}
						position++
						goto l651
					l652:
						position, tokenIndex, depth = position651, tokenIndex651, depth651
						if buffer[position] != rune('L') {
							goto l642
						}
						position++
					}
				l651:
					depth--
					add(rulePegText, position644)
				}
				if !_rules[ruleAction47]() {
					goto l642
				}
				depth--
				add(ruleNullLiteral, position643)
			}
			return true
		l642:
			position, tokenIndex, depth = position642, tokenIndex642, depth642
			return false
		},
		/* 66 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position653, tokenIndex653, depth653 := position, tokenIndex, depth
			{
				position654 := position
				depth++
				{
					position655, tokenIndex655, depth655 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l656
					}
					goto l655
				l656:
					position, tokenIndex, depth = position655, tokenIndex655, depth655
					if !_rules[ruleFALSE]() {
						goto l653
					}
				}
			l655:
				depth--
				add(ruleBooleanLiteral, position654)
			}
			return true
		l653:
			position, tokenIndex, depth = position653, tokenIndex653, depth653
			return false
		},
		/* 67 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action48)> */
		func() bool {
			position657, tokenIndex657, depth657 := position, tokenIndex, depth
			{
				position658 := position
				depth++
				{
					position659 := position
					depth++
					{
						position660, tokenIndex660, depth660 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l661
						}
						position++
						goto l660
					l661:
						position, tokenIndex, depth = position660, tokenIndex660, depth660
						if buffer[position] != rune('T') {
							goto l657
						}
						position++
					}
				l660:
					{
						position662, tokenIndex662, depth662 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l663
						}
						position++
						goto l662
					l663:
						position, tokenIndex, depth = position662, tokenIndex662, depth662
						if buffer[position] != rune('R') {
							goto l657
						}
						position++
					}
				l662:
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l665
						}
						position++
						goto l664
					l665:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if buffer[position] != rune('U') {
							goto l657
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
							goto l657
						}
						position++
					}
				l666:
					depth--
					add(rulePegText, position659)
				}
				if !_rules[ruleAction48]() {
					goto l657
				}
				depth--
				add(ruleTRUE, position658)
			}
			return true
		l657:
			position, tokenIndex, depth = position657, tokenIndex657, depth657
			return false
		},
		/* 68 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action49)> */
		func() bool {
			position668, tokenIndex668, depth668 := position, tokenIndex, depth
			{
				position669 := position
				depth++
				{
					position670 := position
					depth++
					{
						position671, tokenIndex671, depth671 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l672
						}
						position++
						goto l671
					l672:
						position, tokenIndex, depth = position671, tokenIndex671, depth671
						if buffer[position] != rune('F') {
							goto l668
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
							goto l668
						}
						position++
					}
				l673:
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l676
						}
						position++
						goto l675
					l676:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('L') {
							goto l668
						}
						position++
					}
				l675:
					{
						position677, tokenIndex677, depth677 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l678
						}
						position++
						goto l677
					l678:
						position, tokenIndex, depth = position677, tokenIndex677, depth677
						if buffer[position] != rune('S') {
							goto l668
						}
						position++
					}
				l677:
					{
						position679, tokenIndex679, depth679 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l680
						}
						position++
						goto l679
					l680:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
						if buffer[position] != rune('E') {
							goto l668
						}
						position++
					}
				l679:
					depth--
					add(rulePegText, position670)
				}
				if !_rules[ruleAction49]() {
					goto l668
				}
				depth--
				add(ruleFALSE, position669)
			}
			return true
		l668:
			position, tokenIndex, depth = position668, tokenIndex668, depth668
			return false
		},
		/* 69 Wildcard <- <(<'*'> Action50)> */
		func() bool {
			position681, tokenIndex681, depth681 := position, tokenIndex, depth
			{
				position682 := position
				depth++
				{
					position683 := position
					depth++
					if buffer[position] != rune('*') {
						goto l681
					}
					position++
					depth--
					add(rulePegText, position683)
				}
				if !_rules[ruleAction50]() {
					goto l681
				}
				depth--
				add(ruleWildcard, position682)
			}
			return true
		l681:
			position, tokenIndex, depth = position681, tokenIndex681, depth681
			return false
		},
		/* 70 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action51)> */
		func() bool {
			position684, tokenIndex684, depth684 := position, tokenIndex, depth
			{
				position685 := position
				depth++
				{
					position686 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l684
					}
					position++
				l687:
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						{
							position689, tokenIndex689, depth689 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l690
							}
							position++
							if buffer[position] != rune('\'') {
								goto l690
							}
							position++
							goto l689
						l690:
							position, tokenIndex, depth = position689, tokenIndex689, depth689
							{
								position691, tokenIndex691, depth691 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l691
								}
								position++
								goto l688
							l691:
								position, tokenIndex, depth = position691, tokenIndex691, depth691
							}
							if !matchDot() {
								goto l688
							}
						}
					l689:
						goto l687
					l688:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
					}
					if buffer[position] != rune('\'') {
						goto l684
					}
					position++
					depth--
					add(rulePegText, position686)
				}
				if !_rules[ruleAction51]() {
					goto l684
				}
				depth--
				add(ruleStringLiteral, position685)
			}
			return true
		l684:
			position, tokenIndex, depth = position684, tokenIndex684, depth684
			return false
		},
		/* 71 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action52)> */
		func() bool {
			position692, tokenIndex692, depth692 := position, tokenIndex, depth
			{
				position693 := position
				depth++
				{
					position694 := position
					depth++
					{
						position695, tokenIndex695, depth695 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
						if buffer[position] != rune('I') {
							goto l692
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
							goto l692
						}
						position++
					}
				l697:
					{
						position699, tokenIndex699, depth699 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l700
						}
						position++
						goto l699
					l700:
						position, tokenIndex, depth = position699, tokenIndex699, depth699
						if buffer[position] != rune('T') {
							goto l692
						}
						position++
					}
				l699:
					{
						position701, tokenIndex701, depth701 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l702
						}
						position++
						goto l701
					l702:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
						if buffer[position] != rune('R') {
							goto l692
						}
						position++
					}
				l701:
					{
						position703, tokenIndex703, depth703 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l704
						}
						position++
						goto l703
					l704:
						position, tokenIndex, depth = position703, tokenIndex703, depth703
						if buffer[position] != rune('E') {
							goto l692
						}
						position++
					}
				l703:
					{
						position705, tokenIndex705, depth705 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l706
						}
						position++
						goto l705
					l706:
						position, tokenIndex, depth = position705, tokenIndex705, depth705
						if buffer[position] != rune('A') {
							goto l692
						}
						position++
					}
				l705:
					{
						position707, tokenIndex707, depth707 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l708
						}
						position++
						goto l707
					l708:
						position, tokenIndex, depth = position707, tokenIndex707, depth707
						if buffer[position] != rune('M') {
							goto l692
						}
						position++
					}
				l707:
					depth--
					add(rulePegText, position694)
				}
				if !_rules[ruleAction52]() {
					goto l692
				}
				depth--
				add(ruleISTREAM, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 72 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action53)> */
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
						if buffer[position] != rune('d') {
							goto l713
						}
						position++
						goto l712
					l713:
						position, tokenIndex, depth = position712, tokenIndex712, depth712
						if buffer[position] != rune('D') {
							goto l709
						}
						position++
					}
				l712:
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l715
						}
						position++
						goto l714
					l715:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
						if buffer[position] != rune('S') {
							goto l709
						}
						position++
					}
				l714:
					{
						position716, tokenIndex716, depth716 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l717
						}
						position++
						goto l716
					l717:
						position, tokenIndex, depth = position716, tokenIndex716, depth716
						if buffer[position] != rune('T') {
							goto l709
						}
						position++
					}
				l716:
					{
						position718, tokenIndex718, depth718 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l719
						}
						position++
						goto l718
					l719:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
						if buffer[position] != rune('R') {
							goto l709
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
							goto l709
						}
						position++
					}
				l720:
					{
						position722, tokenIndex722, depth722 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l723
						}
						position++
						goto l722
					l723:
						position, tokenIndex, depth = position722, tokenIndex722, depth722
						if buffer[position] != rune('A') {
							goto l709
						}
						position++
					}
				l722:
					{
						position724, tokenIndex724, depth724 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l725
						}
						position++
						goto l724
					l725:
						position, tokenIndex, depth = position724, tokenIndex724, depth724
						if buffer[position] != rune('M') {
							goto l709
						}
						position++
					}
				l724:
					depth--
					add(rulePegText, position711)
				}
				if !_rules[ruleAction53]() {
					goto l709
				}
				depth--
				add(ruleDSTREAM, position710)
			}
			return true
		l709:
			position, tokenIndex, depth = position709, tokenIndex709, depth709
			return false
		},
		/* 73 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action54)> */
		func() bool {
			position726, tokenIndex726, depth726 := position, tokenIndex, depth
			{
				position727 := position
				depth++
				{
					position728 := position
					depth++
					{
						position729, tokenIndex729, depth729 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l730
						}
						position++
						goto l729
					l730:
						position, tokenIndex, depth = position729, tokenIndex729, depth729
						if buffer[position] != rune('R') {
							goto l726
						}
						position++
					}
				l729:
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l732
						}
						position++
						goto l731
					l732:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
						if buffer[position] != rune('S') {
							goto l726
						}
						position++
					}
				l731:
					{
						position733, tokenIndex733, depth733 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l734
						}
						position++
						goto l733
					l734:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
						if buffer[position] != rune('T') {
							goto l726
						}
						position++
					}
				l733:
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l736
						}
						position++
						goto l735
					l736:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if buffer[position] != rune('R') {
							goto l726
						}
						position++
					}
				l735:
					{
						position737, tokenIndex737, depth737 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l738
						}
						position++
						goto l737
					l738:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
						if buffer[position] != rune('E') {
							goto l726
						}
						position++
					}
				l737:
					{
						position739, tokenIndex739, depth739 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l740
						}
						position++
						goto l739
					l740:
						position, tokenIndex, depth = position739, tokenIndex739, depth739
						if buffer[position] != rune('A') {
							goto l726
						}
						position++
					}
				l739:
					{
						position741, tokenIndex741, depth741 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l742
						}
						position++
						goto l741
					l742:
						position, tokenIndex, depth = position741, tokenIndex741, depth741
						if buffer[position] != rune('M') {
							goto l726
						}
						position++
					}
				l741:
					depth--
					add(rulePegText, position728)
				}
				if !_rules[ruleAction54]() {
					goto l726
				}
				depth--
				add(ruleRSTREAM, position727)
			}
			return true
		l726:
			position, tokenIndex, depth = position726, tokenIndex726, depth726
			return false
		},
		/* 74 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action55)> */
		func() bool {
			position743, tokenIndex743, depth743 := position, tokenIndex, depth
			{
				position744 := position
				depth++
				{
					position745 := position
					depth++
					{
						position746, tokenIndex746, depth746 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l747
						}
						position++
						goto l746
					l747:
						position, tokenIndex, depth = position746, tokenIndex746, depth746
						if buffer[position] != rune('T') {
							goto l743
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
							goto l743
						}
						position++
					}
				l748:
					{
						position750, tokenIndex750, depth750 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l751
						}
						position++
						goto l750
					l751:
						position, tokenIndex, depth = position750, tokenIndex750, depth750
						if buffer[position] != rune('P') {
							goto l743
						}
						position++
					}
				l750:
					{
						position752, tokenIndex752, depth752 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l753
						}
						position++
						goto l752
					l753:
						position, tokenIndex, depth = position752, tokenIndex752, depth752
						if buffer[position] != rune('L') {
							goto l743
						}
						position++
					}
				l752:
					{
						position754, tokenIndex754, depth754 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l755
						}
						position++
						goto l754
					l755:
						position, tokenIndex, depth = position754, tokenIndex754, depth754
						if buffer[position] != rune('E') {
							goto l743
						}
						position++
					}
				l754:
					{
						position756, tokenIndex756, depth756 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l757
						}
						position++
						goto l756
					l757:
						position, tokenIndex, depth = position756, tokenIndex756, depth756
						if buffer[position] != rune('S') {
							goto l743
						}
						position++
					}
				l756:
					depth--
					add(rulePegText, position745)
				}
				if !_rules[ruleAction55]() {
					goto l743
				}
				depth--
				add(ruleTUPLES, position744)
			}
			return true
		l743:
			position, tokenIndex, depth = position743, tokenIndex743, depth743
			return false
		},
		/* 75 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action56)> */
		func() bool {
			position758, tokenIndex758, depth758 := position, tokenIndex, depth
			{
				position759 := position
				depth++
				{
					position760 := position
					depth++
					{
						position761, tokenIndex761, depth761 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l762
						}
						position++
						goto l761
					l762:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if buffer[position] != rune('S') {
							goto l758
						}
						position++
					}
				l761:
					{
						position763, tokenIndex763, depth763 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l764
						}
						position++
						goto l763
					l764:
						position, tokenIndex, depth = position763, tokenIndex763, depth763
						if buffer[position] != rune('E') {
							goto l758
						}
						position++
					}
				l763:
					{
						position765, tokenIndex765, depth765 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l766
						}
						position++
						goto l765
					l766:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
						if buffer[position] != rune('C') {
							goto l758
						}
						position++
					}
				l765:
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('O') {
							goto l758
						}
						position++
					}
				l767:
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('N') {
							goto l758
						}
						position++
					}
				l769:
					{
						position771, tokenIndex771, depth771 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l772
						}
						position++
						goto l771
					l772:
						position, tokenIndex, depth = position771, tokenIndex771, depth771
						if buffer[position] != rune('D') {
							goto l758
						}
						position++
					}
				l771:
					{
						position773, tokenIndex773, depth773 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l774
						}
						position++
						goto l773
					l774:
						position, tokenIndex, depth = position773, tokenIndex773, depth773
						if buffer[position] != rune('S') {
							goto l758
						}
						position++
					}
				l773:
					depth--
					add(rulePegText, position760)
				}
				if !_rules[ruleAction56]() {
					goto l758
				}
				depth--
				add(ruleSECONDS, position759)
			}
			return true
		l758:
			position, tokenIndex, depth = position758, tokenIndex758, depth758
			return false
		},
		/* 76 StreamIdentifier <- <(<ident> Action57)> */
		func() bool {
			position775, tokenIndex775, depth775 := position, tokenIndex, depth
			{
				position776 := position
				depth++
				{
					position777 := position
					depth++
					if !_rules[ruleident]() {
						goto l775
					}
					depth--
					add(rulePegText, position777)
				}
				if !_rules[ruleAction57]() {
					goto l775
				}
				depth--
				add(ruleStreamIdentifier, position776)
			}
			return true
		l775:
			position, tokenIndex, depth = position775, tokenIndex775, depth775
			return false
		},
		/* 77 SourceSinkType <- <(<ident> Action58)> */
		func() bool {
			position778, tokenIndex778, depth778 := position, tokenIndex, depth
			{
				position779 := position
				depth++
				{
					position780 := position
					depth++
					if !_rules[ruleident]() {
						goto l778
					}
					depth--
					add(rulePegText, position780)
				}
				if !_rules[ruleAction58]() {
					goto l778
				}
				depth--
				add(ruleSourceSinkType, position779)
			}
			return true
		l778:
			position, tokenIndex, depth = position778, tokenIndex778, depth778
			return false
		},
		/* 78 SourceSinkParamKey <- <(<ident> Action59)> */
		func() bool {
			position781, tokenIndex781, depth781 := position, tokenIndex, depth
			{
				position782 := position
				depth++
				{
					position783 := position
					depth++
					if !_rules[ruleident]() {
						goto l781
					}
					depth--
					add(rulePegText, position783)
				}
				if !_rules[ruleAction59]() {
					goto l781
				}
				depth--
				add(ruleSourceSinkParamKey, position782)
			}
			return true
		l781:
			position, tokenIndex, depth = position781, tokenIndex781, depth781
			return false
		},
		/* 79 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action60)> */
		func() bool {
			position784, tokenIndex784, depth784 := position, tokenIndex, depth
			{
				position785 := position
				depth++
				{
					position786 := position
					depth++
					{
						position787, tokenIndex787, depth787 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l788
						}
						position++
						goto l787
					l788:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
						if buffer[position] != rune('P') {
							goto l784
						}
						position++
					}
				l787:
					{
						position789, tokenIndex789, depth789 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l790
						}
						position++
						goto l789
					l790:
						position, tokenIndex, depth = position789, tokenIndex789, depth789
						if buffer[position] != rune('A') {
							goto l784
						}
						position++
					}
				l789:
					{
						position791, tokenIndex791, depth791 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l792
						}
						position++
						goto l791
					l792:
						position, tokenIndex, depth = position791, tokenIndex791, depth791
						if buffer[position] != rune('U') {
							goto l784
						}
						position++
					}
				l791:
					{
						position793, tokenIndex793, depth793 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l794
						}
						position++
						goto l793
					l794:
						position, tokenIndex, depth = position793, tokenIndex793, depth793
						if buffer[position] != rune('S') {
							goto l784
						}
						position++
					}
				l793:
					{
						position795, tokenIndex795, depth795 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l796
						}
						position++
						goto l795
					l796:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						if buffer[position] != rune('E') {
							goto l784
						}
						position++
					}
				l795:
					{
						position797, tokenIndex797, depth797 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l798
						}
						position++
						goto l797
					l798:
						position, tokenIndex, depth = position797, tokenIndex797, depth797
						if buffer[position] != rune('D') {
							goto l784
						}
						position++
					}
				l797:
					depth--
					add(rulePegText, position786)
				}
				if !_rules[ruleAction60]() {
					goto l784
				}
				depth--
				add(rulePaused, position785)
			}
			return true
		l784:
			position, tokenIndex, depth = position784, tokenIndex784, depth784
			return false
		},
		/* 80 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action61)> */
		func() bool {
			position799, tokenIndex799, depth799 := position, tokenIndex, depth
			{
				position800 := position
				depth++
				{
					position801 := position
					depth++
					{
						position802, tokenIndex802, depth802 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l803
						}
						position++
						goto l802
					l803:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if buffer[position] != rune('U') {
							goto l799
						}
						position++
					}
				l802:
					{
						position804, tokenIndex804, depth804 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l805
						}
						position++
						goto l804
					l805:
						position, tokenIndex, depth = position804, tokenIndex804, depth804
						if buffer[position] != rune('N') {
							goto l799
						}
						position++
					}
				l804:
					{
						position806, tokenIndex806, depth806 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l807
						}
						position++
						goto l806
					l807:
						position, tokenIndex, depth = position806, tokenIndex806, depth806
						if buffer[position] != rune('P') {
							goto l799
						}
						position++
					}
				l806:
					{
						position808, tokenIndex808, depth808 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l809
						}
						position++
						goto l808
					l809:
						position, tokenIndex, depth = position808, tokenIndex808, depth808
						if buffer[position] != rune('A') {
							goto l799
						}
						position++
					}
				l808:
					{
						position810, tokenIndex810, depth810 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l811
						}
						position++
						goto l810
					l811:
						position, tokenIndex, depth = position810, tokenIndex810, depth810
						if buffer[position] != rune('U') {
							goto l799
						}
						position++
					}
				l810:
					{
						position812, tokenIndex812, depth812 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l813
						}
						position++
						goto l812
					l813:
						position, tokenIndex, depth = position812, tokenIndex812, depth812
						if buffer[position] != rune('S') {
							goto l799
						}
						position++
					}
				l812:
					{
						position814, tokenIndex814, depth814 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l815
						}
						position++
						goto l814
					l815:
						position, tokenIndex, depth = position814, tokenIndex814, depth814
						if buffer[position] != rune('E') {
							goto l799
						}
						position++
					}
				l814:
					{
						position816, tokenIndex816, depth816 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l817
						}
						position++
						goto l816
					l817:
						position, tokenIndex, depth = position816, tokenIndex816, depth816
						if buffer[position] != rune('D') {
							goto l799
						}
						position++
					}
				l816:
					depth--
					add(rulePegText, position801)
				}
				if !_rules[ruleAction61]() {
					goto l799
				}
				depth--
				add(ruleUnpaused, position800)
			}
			return true
		l799:
			position, tokenIndex, depth = position799, tokenIndex799, depth799
			return false
		},
		/* 81 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action62)> */
		func() bool {
			position818, tokenIndex818, depth818 := position, tokenIndex, depth
			{
				position819 := position
				depth++
				{
					position820 := position
					depth++
					{
						position821, tokenIndex821, depth821 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l822
						}
						position++
						goto l821
					l822:
						position, tokenIndex, depth = position821, tokenIndex821, depth821
						if buffer[position] != rune('O') {
							goto l818
						}
						position++
					}
				l821:
					{
						position823, tokenIndex823, depth823 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l824
						}
						position++
						goto l823
					l824:
						position, tokenIndex, depth = position823, tokenIndex823, depth823
						if buffer[position] != rune('R') {
							goto l818
						}
						position++
					}
				l823:
					depth--
					add(rulePegText, position820)
				}
				if !_rules[ruleAction62]() {
					goto l818
				}
				depth--
				add(ruleOr, position819)
			}
			return true
		l818:
			position, tokenIndex, depth = position818, tokenIndex818, depth818
			return false
		},
		/* 82 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action63)> */
		func() bool {
			position825, tokenIndex825, depth825 := position, tokenIndex, depth
			{
				position826 := position
				depth++
				{
					position827 := position
					depth++
					{
						position828, tokenIndex828, depth828 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l829
						}
						position++
						goto l828
					l829:
						position, tokenIndex, depth = position828, tokenIndex828, depth828
						if buffer[position] != rune('A') {
							goto l825
						}
						position++
					}
				l828:
					{
						position830, tokenIndex830, depth830 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l831
						}
						position++
						goto l830
					l831:
						position, tokenIndex, depth = position830, tokenIndex830, depth830
						if buffer[position] != rune('N') {
							goto l825
						}
						position++
					}
				l830:
					{
						position832, tokenIndex832, depth832 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l833
						}
						position++
						goto l832
					l833:
						position, tokenIndex, depth = position832, tokenIndex832, depth832
						if buffer[position] != rune('D') {
							goto l825
						}
						position++
					}
				l832:
					depth--
					add(rulePegText, position827)
				}
				if !_rules[ruleAction63]() {
					goto l825
				}
				depth--
				add(ruleAnd, position826)
			}
			return true
		l825:
			position, tokenIndex, depth = position825, tokenIndex825, depth825
			return false
		},
		/* 83 Equal <- <(<'='> Action64)> */
		func() bool {
			position834, tokenIndex834, depth834 := position, tokenIndex, depth
			{
				position835 := position
				depth++
				{
					position836 := position
					depth++
					if buffer[position] != rune('=') {
						goto l834
					}
					position++
					depth--
					add(rulePegText, position836)
				}
				if !_rules[ruleAction64]() {
					goto l834
				}
				depth--
				add(ruleEqual, position835)
			}
			return true
		l834:
			position, tokenIndex, depth = position834, tokenIndex834, depth834
			return false
		},
		/* 84 Less <- <(<'<'> Action65)> */
		func() bool {
			position837, tokenIndex837, depth837 := position, tokenIndex, depth
			{
				position838 := position
				depth++
				{
					position839 := position
					depth++
					if buffer[position] != rune('<') {
						goto l837
					}
					position++
					depth--
					add(rulePegText, position839)
				}
				if !_rules[ruleAction65]() {
					goto l837
				}
				depth--
				add(ruleLess, position838)
			}
			return true
		l837:
			position, tokenIndex, depth = position837, tokenIndex837, depth837
			return false
		},
		/* 85 LessOrEqual <- <(<('<' '=')> Action66)> */
		func() bool {
			position840, tokenIndex840, depth840 := position, tokenIndex, depth
			{
				position841 := position
				depth++
				{
					position842 := position
					depth++
					if buffer[position] != rune('<') {
						goto l840
					}
					position++
					if buffer[position] != rune('=') {
						goto l840
					}
					position++
					depth--
					add(rulePegText, position842)
				}
				if !_rules[ruleAction66]() {
					goto l840
				}
				depth--
				add(ruleLessOrEqual, position841)
			}
			return true
		l840:
			position, tokenIndex, depth = position840, tokenIndex840, depth840
			return false
		},
		/* 86 Greater <- <(<'>'> Action67)> */
		func() bool {
			position843, tokenIndex843, depth843 := position, tokenIndex, depth
			{
				position844 := position
				depth++
				{
					position845 := position
					depth++
					if buffer[position] != rune('>') {
						goto l843
					}
					position++
					depth--
					add(rulePegText, position845)
				}
				if !_rules[ruleAction67]() {
					goto l843
				}
				depth--
				add(ruleGreater, position844)
			}
			return true
		l843:
			position, tokenIndex, depth = position843, tokenIndex843, depth843
			return false
		},
		/* 87 GreaterOrEqual <- <(<('>' '=')> Action68)> */
		func() bool {
			position846, tokenIndex846, depth846 := position, tokenIndex, depth
			{
				position847 := position
				depth++
				{
					position848 := position
					depth++
					if buffer[position] != rune('>') {
						goto l846
					}
					position++
					if buffer[position] != rune('=') {
						goto l846
					}
					position++
					depth--
					add(rulePegText, position848)
				}
				if !_rules[ruleAction68]() {
					goto l846
				}
				depth--
				add(ruleGreaterOrEqual, position847)
			}
			return true
		l846:
			position, tokenIndex, depth = position846, tokenIndex846, depth846
			return false
		},
		/* 88 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action69)> */
		func() bool {
			position849, tokenIndex849, depth849 := position, tokenIndex, depth
			{
				position850 := position
				depth++
				{
					position851 := position
					depth++
					{
						position852, tokenIndex852, depth852 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l853
						}
						position++
						if buffer[position] != rune('=') {
							goto l853
						}
						position++
						goto l852
					l853:
						position, tokenIndex, depth = position852, tokenIndex852, depth852
						if buffer[position] != rune('<') {
							goto l849
						}
						position++
						if buffer[position] != rune('>') {
							goto l849
						}
						position++
					}
				l852:
					depth--
					add(rulePegText, position851)
				}
				if !_rules[ruleAction69]() {
					goto l849
				}
				depth--
				add(ruleNotEqual, position850)
			}
			return true
		l849:
			position, tokenIndex, depth = position849, tokenIndex849, depth849
			return false
		},
		/* 89 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action70)> */
		func() bool {
			position854, tokenIndex854, depth854 := position, tokenIndex, depth
			{
				position855 := position
				depth++
				{
					position856 := position
					depth++
					{
						position857, tokenIndex857, depth857 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l858
						}
						position++
						goto l857
					l858:
						position, tokenIndex, depth = position857, tokenIndex857, depth857
						if buffer[position] != rune('I') {
							goto l854
						}
						position++
					}
				l857:
					{
						position859, tokenIndex859, depth859 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l860
						}
						position++
						goto l859
					l860:
						position, tokenIndex, depth = position859, tokenIndex859, depth859
						if buffer[position] != rune('S') {
							goto l854
						}
						position++
					}
				l859:
					depth--
					add(rulePegText, position856)
				}
				if !_rules[ruleAction70]() {
					goto l854
				}
				depth--
				add(ruleIs, position855)
			}
			return true
		l854:
			position, tokenIndex, depth = position854, tokenIndex854, depth854
			return false
		},
		/* 90 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action71)> */
		func() bool {
			position861, tokenIndex861, depth861 := position, tokenIndex, depth
			{
				position862 := position
				depth++
				{
					position863 := position
					depth++
					{
						position864, tokenIndex864, depth864 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l865
						}
						position++
						goto l864
					l865:
						position, tokenIndex, depth = position864, tokenIndex864, depth864
						if buffer[position] != rune('I') {
							goto l861
						}
						position++
					}
				l864:
					{
						position866, tokenIndex866, depth866 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l867
						}
						position++
						goto l866
					l867:
						position, tokenIndex, depth = position866, tokenIndex866, depth866
						if buffer[position] != rune('S') {
							goto l861
						}
						position++
					}
				l866:
					if !_rules[rulesp]() {
						goto l861
					}
					{
						position868, tokenIndex868, depth868 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l869
						}
						position++
						goto l868
					l869:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('N') {
							goto l861
						}
						position++
					}
				l868:
					{
						position870, tokenIndex870, depth870 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l871
						}
						position++
						goto l870
					l871:
						position, tokenIndex, depth = position870, tokenIndex870, depth870
						if buffer[position] != rune('O') {
							goto l861
						}
						position++
					}
				l870:
					{
						position872, tokenIndex872, depth872 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l873
						}
						position++
						goto l872
					l873:
						position, tokenIndex, depth = position872, tokenIndex872, depth872
						if buffer[position] != rune('T') {
							goto l861
						}
						position++
					}
				l872:
					depth--
					add(rulePegText, position863)
				}
				if !_rules[ruleAction71]() {
					goto l861
				}
				depth--
				add(ruleIsNot, position862)
			}
			return true
		l861:
			position, tokenIndex, depth = position861, tokenIndex861, depth861
			return false
		},
		/* 91 Plus <- <(<'+'> Action72)> */
		func() bool {
			position874, tokenIndex874, depth874 := position, tokenIndex, depth
			{
				position875 := position
				depth++
				{
					position876 := position
					depth++
					if buffer[position] != rune('+') {
						goto l874
					}
					position++
					depth--
					add(rulePegText, position876)
				}
				if !_rules[ruleAction72]() {
					goto l874
				}
				depth--
				add(rulePlus, position875)
			}
			return true
		l874:
			position, tokenIndex, depth = position874, tokenIndex874, depth874
			return false
		},
		/* 92 Minus <- <(<'-'> Action73)> */
		func() bool {
			position877, tokenIndex877, depth877 := position, tokenIndex, depth
			{
				position878 := position
				depth++
				{
					position879 := position
					depth++
					if buffer[position] != rune('-') {
						goto l877
					}
					position++
					depth--
					add(rulePegText, position879)
				}
				if !_rules[ruleAction73]() {
					goto l877
				}
				depth--
				add(ruleMinus, position878)
			}
			return true
		l877:
			position, tokenIndex, depth = position877, tokenIndex877, depth877
			return false
		},
		/* 93 Multiply <- <(<'*'> Action74)> */
		func() bool {
			position880, tokenIndex880, depth880 := position, tokenIndex, depth
			{
				position881 := position
				depth++
				{
					position882 := position
					depth++
					if buffer[position] != rune('*') {
						goto l880
					}
					position++
					depth--
					add(rulePegText, position882)
				}
				if !_rules[ruleAction74]() {
					goto l880
				}
				depth--
				add(ruleMultiply, position881)
			}
			return true
		l880:
			position, tokenIndex, depth = position880, tokenIndex880, depth880
			return false
		},
		/* 94 Divide <- <(<'/'> Action75)> */
		func() bool {
			position883, tokenIndex883, depth883 := position, tokenIndex, depth
			{
				position884 := position
				depth++
				{
					position885 := position
					depth++
					if buffer[position] != rune('/') {
						goto l883
					}
					position++
					depth--
					add(rulePegText, position885)
				}
				if !_rules[ruleAction75]() {
					goto l883
				}
				depth--
				add(ruleDivide, position884)
			}
			return true
		l883:
			position, tokenIndex, depth = position883, tokenIndex883, depth883
			return false
		},
		/* 95 Modulo <- <(<'%'> Action76)> */
		func() bool {
			position886, tokenIndex886, depth886 := position, tokenIndex, depth
			{
				position887 := position
				depth++
				{
					position888 := position
					depth++
					if buffer[position] != rune('%') {
						goto l886
					}
					position++
					depth--
					add(rulePegText, position888)
				}
				if !_rules[ruleAction76]() {
					goto l886
				}
				depth--
				add(ruleModulo, position887)
			}
			return true
		l886:
			position, tokenIndex, depth = position886, tokenIndex886, depth886
			return false
		},
		/* 96 Identifier <- <(<ident> Action77)> */
		func() bool {
			position889, tokenIndex889, depth889 := position, tokenIndex, depth
			{
				position890 := position
				depth++
				{
					position891 := position
					depth++
					if !_rules[ruleident]() {
						goto l889
					}
					depth--
					add(rulePegText, position891)
				}
				if !_rules[ruleAction77]() {
					goto l889
				}
				depth--
				add(ruleIdentifier, position890)
			}
			return true
		l889:
			position, tokenIndex, depth = position889, tokenIndex889, depth889
			return false
		},
		/* 97 TargetIdentifier <- <(<jsonPath> Action78)> */
		func() bool {
			position892, tokenIndex892, depth892 := position, tokenIndex, depth
			{
				position893 := position
				depth++
				{
					position894 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l892
					}
					depth--
					add(rulePegText, position894)
				}
				if !_rules[ruleAction78]() {
					goto l892
				}
				depth--
				add(ruleTargetIdentifier, position893)
			}
			return true
		l892:
			position, tokenIndex, depth = position892, tokenIndex892, depth892
			return false
		},
		/* 98 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position895, tokenIndex895, depth895 := position, tokenIndex, depth
			{
				position896 := position
				depth++
				{
					position897, tokenIndex897, depth897 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l898
					}
					position++
					goto l897
				l898:
					position, tokenIndex, depth = position897, tokenIndex897, depth897
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l895
					}
					position++
				}
			l897:
			l899:
				{
					position900, tokenIndex900, depth900 := position, tokenIndex, depth
					{
						position901, tokenIndex901, depth901 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l902
						}
						position++
						goto l901
					l902:
						position, tokenIndex, depth = position901, tokenIndex901, depth901
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l903
						}
						position++
						goto l901
					l903:
						position, tokenIndex, depth = position901, tokenIndex901, depth901
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l904
						}
						position++
						goto l901
					l904:
						position, tokenIndex, depth = position901, tokenIndex901, depth901
						if buffer[position] != rune('_') {
							goto l900
						}
						position++
					}
				l901:
					goto l899
				l900:
					position, tokenIndex, depth = position900, tokenIndex900, depth900
				}
				depth--
				add(ruleident, position896)
			}
			return true
		l895:
			position, tokenIndex, depth = position895, tokenIndex895, depth895
			return false
		},
		/* 99 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position905, tokenIndex905, depth905 := position, tokenIndex, depth
			{
				position906 := position
				depth++
				{
					position907, tokenIndex907, depth907 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l908
					}
					position++
					goto l907
				l908:
					position, tokenIndex, depth = position907, tokenIndex907, depth907
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l905
					}
					position++
				}
			l907:
			l909:
				{
					position910, tokenIndex910, depth910 := position, tokenIndex, depth
					{
						position911, tokenIndex911, depth911 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l912
						}
						position++
						goto l911
					l912:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l913
						}
						position++
						goto l911
					l913:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l914
						}
						position++
						goto l911
					l914:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('_') {
							goto l915
						}
						position++
						goto l911
					l915:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('.') {
							goto l916
						}
						position++
						goto l911
					l916:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('[') {
							goto l917
						}
						position++
						goto l911
					l917:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune(']') {
							goto l918
						}
						position++
						goto l911
					l918:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('"') {
							goto l910
						}
						position++
					}
				l911:
					goto l909
				l910:
					position, tokenIndex, depth = position910, tokenIndex910, depth910
				}
				depth--
				add(rulejsonPath, position906)
			}
			return true
		l905:
			position, tokenIndex, depth = position905, tokenIndex905, depth905
			return false
		},
		/* 100 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position920 := position
				depth++
			l921:
				{
					position922, tokenIndex922, depth922 := position, tokenIndex, depth
					{
						position923, tokenIndex923, depth923 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l924
						}
						position++
						goto l923
					l924:
						position, tokenIndex, depth = position923, tokenIndex923, depth923
						if buffer[position] != rune('\t') {
							goto l925
						}
						position++
						goto l923
					l925:
						position, tokenIndex, depth = position923, tokenIndex923, depth923
						if buffer[position] != rune('\n') {
							goto l922
						}
						position++
					}
				l923:
					goto l921
				l922:
					position, tokenIndex, depth = position922, tokenIndex922, depth922
				}
				depth--
				add(rulesp, position920)
			}
			return true
		},
		/* 102 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 103 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 104 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 105 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 106 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 107 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 108 Action6 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 109 Action7 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 110 Action8 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 111 Action9 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		nil,
		/* 113 Action10 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 114 Action11 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 115 Action12 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 116 Action13 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 117 Action14 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 118 Action15 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 119 Action16 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 120 Action17 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 121 Action18 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 122 Action19 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 123 Action20 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 124 Action21 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 125 Action22 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 126 Action23 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 127 Action24 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 128 Action25 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 129 Action26 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 130 Action27 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 131 Action28 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 132 Action29 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 133 Action30 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 134 Action31 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 135 Action32 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 136 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 137 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 138 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 139 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 140 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 141 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 142 Action39 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 143 Action40 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 144 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 145 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 146 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 147 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 148 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 149 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 150 Action47 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 151 Action48 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 152 Action49 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 153 Action50 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 154 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 155 Action52 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 156 Action53 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 157 Action54 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 158 Action55 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 159 Action56 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 160 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 161 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 162 Action59 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 163 Action60 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 164 Action61 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 165 Action62 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 166 Action63 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 167 Action64 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 168 Action65 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 169 Action66 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 170 Action67 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 171 Action68 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 172 Action69 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 173 Action70 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 174 Action71 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 175 Action72 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 176 Action73 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 177 Action74 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 178 Action75 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 179 Action76 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 180 Action77 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 181 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
	}
	p.rules = _rules
}

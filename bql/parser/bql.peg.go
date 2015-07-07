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
	ruleInterval
	ruleTimeInterval
	ruleTuplesInterval
	ruleRelations
	ruleFilter
	ruleGrouping
	ruleGroupList
	ruleHaving
	ruleRelationLike
	ruleAliasedStreamWindow
	ruleStreamWindow
	ruleStreamLike
	ruleUDSFFuncApp
	ruleSourceSinkSpecs
	ruleSourceSinkParam
	ruleSourceSinkParamVal
	rulePausedOpt
	ruleExpression
	ruleorExpr
	ruleandExpr
	rulenotExpr
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
	ruleNot
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
	"Interval",
	"TimeInterval",
	"TuplesInterval",
	"Relations",
	"Filter",
	"Grouping",
	"GroupList",
	"Having",
	"RelationLike",
	"AliasedStreamWindow",
	"StreamWindow",
	"StreamLike",
	"UDSFFuncApp",
	"SourceSinkSpecs",
	"SourceSinkParam",
	"SourceSinkParamVal",
	"PausedOpt",
	"Expression",
	"orExpr",
	"andExpr",
	"notExpr",
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
	"Not",
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
	rules  [177]func() bool
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

			p.AssembleAliasedStreamWindow()

		case ruleAction24:

			p.AssembleStreamWindow()

		case ruleAction25:

			p.AssembleUDSFFuncApp()

		case ruleAction26:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction27:

			p.AssembleSourceSinkParam()

		case ruleAction28:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction29:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction30:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction31:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction32:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction33:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction34:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleFuncApp()

		case ruleAction37:

			p.AssembleExpressions(begin, end)

		case ruleAction38:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction44:

			p.PushComponent(begin, end, NewNullLiteral())

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

			p.PushComponent(begin, end, Not)

		case ruleAction62:

			p.PushComponent(begin, end, Equal)

		case ruleAction63:

			p.PushComponent(begin, end, Less)

		case ruleAction64:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction65:

			p.PushComponent(begin, end, Greater)

		case ruleAction66:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction67:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction68:

			p.PushComponent(begin, end, Is)

		case ruleAction69:

			p.PushComponent(begin, end, IsNot)

		case ruleAction70:

			p.PushComponent(begin, end, Plus)

		case ruleAction71:

			p.PushComponent(begin, end, Minus)

		case ruleAction72:

			p.PushComponent(begin, end, Multiply)

		case ruleAction73:

			p.PushComponent(begin, end, Divide)

		case ruleAction74:

			p.PushComponent(begin, end, Modulo)

		case ruleAction75:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction76:

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
		/* 2 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
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
				if !_rules[ruleEmitter]() {
					goto l19
				}
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleProjections]() {
					goto l19
				}
				if !_rules[rulesp]() {
					goto l19
				}
				if !_rules[ruleWindowedFrom]() {
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
		/* 8 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action6)> */
		func() bool {
			position193, tokenIndex193, depth193 := position, tokenIndex, depth
			{
				position194 := position
				depth++
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('I') {
						goto l193
					}
					position++
				}
			l195:
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('N') {
						goto l193
					}
					position++
				}
			l197:
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('S') {
						goto l193
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('E') {
						goto l193
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('R') {
						goto l193
					}
					position++
				}
			l203:
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('T') {
						goto l193
					}
					position++
				}
			l205:
				if !_rules[rulesp]() {
					goto l193
				}
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('I') {
						goto l193
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('N') {
						goto l193
					}
					position++
				}
			l209:
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('T') {
						goto l193
					}
					position++
				}
			l211:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('O') {
						goto l193
					}
					position++
				}
			l213:
				if !_rules[rulesp]() {
					goto l193
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l193
				}
				if !_rules[rulesp]() {
					goto l193
				}
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('F') {
						goto l193
					}
					position++
				}
			l215:
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('R') {
						goto l193
					}
					position++
				}
			l217:
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('O') {
						goto l193
					}
					position++
				}
			l219:
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if buffer[position] != rune('M') {
						goto l193
					}
					position++
				}
			l221:
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
				add(ruleInsertIntoFromStmt, position194)
			}
			return true
		l193:
			position, tokenIndex, depth = position193, tokenIndex193, depth193
			return false
		},
		/* 9 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action7)> */
		func() bool {
			position223, tokenIndex223, depth223 := position, tokenIndex, depth
			{
				position224 := position
				depth++
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('P') {
						goto l223
					}
					position++
				}
			l225:
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l228
					}
					position++
					goto l227
				l228:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
					if buffer[position] != rune('A') {
						goto l223
					}
					position++
				}
			l227:
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('U') {
						goto l223
					}
					position++
				}
			l229:
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
						goto l223
					}
					position++
				}
			l231:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('E') {
						goto l223
					}
					position++
				}
			l233:
				if !_rules[rulesp]() {
					goto l223
				}
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l236
					}
					position++
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if buffer[position] != rune('S') {
						goto l223
					}
					position++
				}
			l235:
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if buffer[position] != rune('O') {
						goto l223
					}
					position++
				}
			l237:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					if buffer[position] != rune('U') {
						goto l223
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('R') {
						goto l223
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('C') {
						goto l223
					}
					position++
				}
			l243:
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l246
					}
					position++
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if buffer[position] != rune('E') {
						goto l223
					}
					position++
				}
			l245:
				if !_rules[rulesp]() {
					goto l223
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l223
				}
				if !_rules[ruleAction7]() {
					goto l223
				}
				depth--
				add(rulePauseSourceStmt, position224)
			}
			return true
		l223:
			position, tokenIndex, depth = position223, tokenIndex223, depth223
			return false
		},
		/* 10 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action8)> */
		func() bool {
			position247, tokenIndex247, depth247 := position, tokenIndex, depth
			{
				position248 := position
				depth++
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('R') {
						goto l247
					}
					position++
				}
			l249:
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
						goto l247
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('S') {
						goto l247
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
						goto l247
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
						goto l247
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('E') {
						goto l247
					}
					position++
				}
			l259:
				if !_rules[rulesp]() {
					goto l247
				}
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('S') {
						goto l247
					}
					position++
				}
			l261:
				{
					position263, tokenIndex263, depth263 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if buffer[position] != rune('O') {
						goto l247
					}
					position++
				}
			l263:
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('U') {
						goto l247
					}
					position++
				}
			l265:
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l268
					}
					position++
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if buffer[position] != rune('R') {
						goto l247
					}
					position++
				}
			l267:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('C') {
						goto l247
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
						goto l247
					}
					position++
				}
			l271:
				if !_rules[rulesp]() {
					goto l247
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l247
				}
				if !_rules[ruleAction8]() {
					goto l247
				}
				depth--
				add(ruleResumeSourceStmt, position248)
			}
			return true
		l247:
			position, tokenIndex, depth = position247, tokenIndex247, depth247
			return false
		},
		/* 11 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action9)> */
		func() bool {
			position273, tokenIndex273, depth273 := position, tokenIndex, depth
			{
				position274 := position
				depth++
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('R') {
						goto l273
					}
					position++
				}
			l275:
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('E') {
						goto l273
					}
					position++
				}
			l277:
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('W') {
						goto l273
					}
					position++
				}
			l279:
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('I') {
						goto l273
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('N') {
						goto l273
					}
					position++
				}
			l283:
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('D') {
						goto l273
					}
					position++
				}
			l285:
				if !_rules[rulesp]() {
					goto l273
				}
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('S') {
						goto l273
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('O') {
						goto l273
					}
					position++
				}
			l289:
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('U') {
						goto l273
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
						goto l273
					}
					position++
				}
			l293:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('C') {
						goto l273
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
						goto l273
					}
					position++
				}
			l297:
				if !_rules[rulesp]() {
					goto l273
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l273
				}
				if !_rules[ruleAction9]() {
					goto l273
				}
				depth--
				add(ruleRewindSourceStmt, position274)
			}
			return true
		l273:
			position, tokenIndex, depth = position273, tokenIndex273, depth273
			return false
		},
		/* 12 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action10)> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{
				position300 := position
				depth++
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l302
					}
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if !_rules[ruleDSTREAM]() {
						goto l303
					}
					goto l301
				l303:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if !_rules[ruleRSTREAM]() {
						goto l299
					}
				}
			l301:
				{
					position304 := position
					depth++
					{
						position305, tokenIndex305, depth305 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l305
						}
						if buffer[position] != rune('[') {
							goto l305
						}
						position++
						if !_rules[rulesp]() {
							goto l305
						}
						{
							position307, tokenIndex307, depth307 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l308
							}
							position++
							goto l307
						l308:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if buffer[position] != rune('E') {
								goto l305
							}
							position++
						}
					l307:
						{
							position309, tokenIndex309, depth309 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l310
							}
							position++
							goto l309
						l310:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if buffer[position] != rune('V') {
								goto l305
							}
							position++
						}
					l309:
						{
							position311, tokenIndex311, depth311 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l312
							}
							position++
							goto l311
						l312:
							position, tokenIndex, depth = position311, tokenIndex311, depth311
							if buffer[position] != rune('E') {
								goto l305
							}
							position++
						}
					l311:
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l314
							}
							position++
							goto l313
						l314:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
							if buffer[position] != rune('R') {
								goto l305
							}
							position++
						}
					l313:
						{
							position315, tokenIndex315, depth315 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l316
							}
							position++
							goto l315
						l316:
							position, tokenIndex, depth = position315, tokenIndex315, depth315
							if buffer[position] != rune('Y') {
								goto l305
							}
							position++
						}
					l315:
						if !_rules[rulesp]() {
							goto l305
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l305
						}
						if !_rules[rulesp]() {
							goto l305
						}
						if buffer[position] != rune(']') {
							goto l305
						}
						position++
						goto l306
					l305:
						position, tokenIndex, depth = position305, tokenIndex305, depth305
					}
				l306:
					depth--
					add(rulePegText, position304)
				}
				if !_rules[ruleAction10]() {
					goto l299
				}
				depth--
				add(ruleEmitter, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 13 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position317, tokenIndex317, depth317 := position, tokenIndex, depth
			{
				position318 := position
				depth++
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l320
					}
				l321:
					{
						position322, tokenIndex322, depth322 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l322
						}
						if buffer[position] != rune(',') {
							goto l322
						}
						position++
						if !_rules[rulesp]() {
							goto l322
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l322
						}
						goto l321
					l322:
						position, tokenIndex, depth = position322, tokenIndex322, depth322
					}
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if !_rules[ruleTimeEmitterInterval]() {
						goto l323
					}
					goto l319
				l323:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if !_rules[ruleTupleEmitterInterval]() {
						goto l317
					}
				}
			l319:
				depth--
				add(ruleEmitterIntervals, position318)
			}
			return true
		l317:
			position, tokenIndex, depth = position317, tokenIndex317, depth317
			return false
		},
		/* 14 TimeEmitterInterval <- <(<TimeInterval> Action11)> */
		func() bool {
			position324, tokenIndex324, depth324 := position, tokenIndex, depth
			{
				position325 := position
				depth++
				{
					position326 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l324
					}
					depth--
					add(rulePegText, position326)
				}
				if !_rules[ruleAction11]() {
					goto l324
				}
				depth--
				add(ruleTimeEmitterInterval, position325)
			}
			return true
		l324:
			position, tokenIndex, depth = position324, tokenIndex324, depth324
			return false
		},
		/* 15 TupleEmitterInterval <- <(<TuplesInterval> Action12)> */
		func() bool {
			position327, tokenIndex327, depth327 := position, tokenIndex, depth
			{
				position328 := position
				depth++
				{
					position329 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l327
					}
					depth--
					add(rulePegText, position329)
				}
				if !_rules[ruleAction12]() {
					goto l327
				}
				depth--
				add(ruleTupleEmitterInterval, position328)
			}
			return true
		l327:
			position, tokenIndex, depth = position327, tokenIndex327, depth327
			return false
		},
		/* 16 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action13)> */
		func() bool {
			position330, tokenIndex330, depth330 := position, tokenIndex, depth
			{
				position331 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l330
				}
				if !_rules[rulesp]() {
					goto l330
				}
				{
					position332, tokenIndex332, depth332 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l333
					}
					position++
					goto l332
				l333:
					position, tokenIndex, depth = position332, tokenIndex332, depth332
					if buffer[position] != rune('I') {
						goto l330
					}
					position++
				}
			l332:
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if buffer[position] != rune('N') {
						goto l330
					}
					position++
				}
			l334:
				if !_rules[rulesp]() {
					goto l330
				}
				if !_rules[ruleStream]() {
					goto l330
				}
				if !_rules[ruleAction13]() {
					goto l330
				}
				depth--
				add(ruleTupleEmitterFromInterval, position331)
			}
			return true
		l330:
			position, tokenIndex, depth = position330, tokenIndex330, depth330
			return false
		},
		/* 17 Projections <- <(<(Projection sp (',' sp Projection)*)> Action14)> */
		func() bool {
			position336, tokenIndex336, depth336 := position, tokenIndex, depth
			{
				position337 := position
				depth++
				{
					position338 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l336
					}
					if !_rules[rulesp]() {
						goto l336
					}
				l339:
					{
						position340, tokenIndex340, depth340 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l340
						}
						position++
						if !_rules[rulesp]() {
							goto l340
						}
						if !_rules[ruleProjection]() {
							goto l340
						}
						goto l339
					l340:
						position, tokenIndex, depth = position340, tokenIndex340, depth340
					}
					depth--
					add(rulePegText, position338)
				}
				if !_rules[ruleAction14]() {
					goto l336
				}
				depth--
				add(ruleProjections, position337)
			}
			return true
		l336:
			position, tokenIndex, depth = position336, tokenIndex336, depth336
			return false
		},
		/* 18 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
			{
				position342 := position
				depth++
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l344
					}
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if !_rules[ruleExpression]() {
						goto l345
					}
					goto l343
				l345:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if !_rules[ruleWildcard]() {
						goto l341
					}
				}
			l343:
				depth--
				add(ruleProjection, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
			return false
		},
		/* 19 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action15)> */
		func() bool {
			position346, tokenIndex346, depth346 := position, tokenIndex, depth
			{
				position347 := position
				depth++
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l349
					}
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if !_rules[ruleWildcard]() {
						goto l346
					}
				}
			l348:
				if !_rules[rulesp]() {
					goto l346
				}
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if buffer[position] != rune('A') {
						goto l346
					}
					position++
				}
			l350:
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('S') {
						goto l346
					}
					position++
				}
			l352:
				if !_rules[rulesp]() {
					goto l346
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l346
				}
				if !_rules[ruleAction15]() {
					goto l346
				}
				depth--
				add(ruleAliasExpression, position347)
			}
			return true
		l346:
			position, tokenIndex, depth = position346, tokenIndex346, depth346
			return false
		},
		/* 20 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action16)> */
		func() bool {
			position354, tokenIndex354, depth354 := position, tokenIndex, depth
			{
				position355 := position
				depth++
				{
					position356 := position
					depth++
					{
						position357, tokenIndex357, depth357 := position, tokenIndex, depth
						{
							position359, tokenIndex359, depth359 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l360
							}
							position++
							goto l359
						l360:
							position, tokenIndex, depth = position359, tokenIndex359, depth359
							if buffer[position] != rune('F') {
								goto l357
							}
							position++
						}
					l359:
						{
							position361, tokenIndex361, depth361 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l362
							}
							position++
							goto l361
						l362:
							position, tokenIndex, depth = position361, tokenIndex361, depth361
							if buffer[position] != rune('R') {
								goto l357
							}
							position++
						}
					l361:
						{
							position363, tokenIndex363, depth363 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l364
							}
							position++
							goto l363
						l364:
							position, tokenIndex, depth = position363, tokenIndex363, depth363
							if buffer[position] != rune('O') {
								goto l357
							}
							position++
						}
					l363:
						{
							position365, tokenIndex365, depth365 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l366
							}
							position++
							goto l365
						l366:
							position, tokenIndex, depth = position365, tokenIndex365, depth365
							if buffer[position] != rune('M') {
								goto l357
							}
							position++
						}
					l365:
						if !_rules[rulesp]() {
							goto l357
						}
						if !_rules[ruleRelations]() {
							goto l357
						}
						if !_rules[rulesp]() {
							goto l357
						}
						goto l358
					l357:
						position, tokenIndex, depth = position357, tokenIndex357, depth357
					}
				l358:
					depth--
					add(rulePegText, position356)
				}
				if !_rules[ruleAction16]() {
					goto l354
				}
				depth--
				add(ruleWindowedFrom, position355)
			}
			return true
		l354:
			position, tokenIndex, depth = position354, tokenIndex354, depth354
			return false
		},
		/* 21 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position367, tokenIndex367, depth367 := position, tokenIndex, depth
			{
				position368 := position
				depth++
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l370
					}
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if !_rules[ruleTuplesInterval]() {
						goto l367
					}
				}
			l369:
				depth--
				add(ruleInterval, position368)
			}
			return true
		l367:
			position, tokenIndex, depth = position367, tokenIndex367, depth367
			return false
		},
		/* 22 TimeInterval <- <(NumericLiteral sp SECONDS Action17)> */
		func() bool {
			position371, tokenIndex371, depth371 := position, tokenIndex, depth
			{
				position372 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l371
				}
				if !_rules[rulesp]() {
					goto l371
				}
				if !_rules[ruleSECONDS]() {
					goto l371
				}
				if !_rules[ruleAction17]() {
					goto l371
				}
				depth--
				add(ruleTimeInterval, position372)
			}
			return true
		l371:
			position, tokenIndex, depth = position371, tokenIndex371, depth371
			return false
		},
		/* 23 TuplesInterval <- <(NumericLiteral sp TUPLES Action18)> */
		func() bool {
			position373, tokenIndex373, depth373 := position, tokenIndex, depth
			{
				position374 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l373
				}
				if !_rules[rulesp]() {
					goto l373
				}
				if !_rules[ruleTUPLES]() {
					goto l373
				}
				if !_rules[ruleAction18]() {
					goto l373
				}
				depth--
				add(ruleTuplesInterval, position374)
			}
			return true
		l373:
			position, tokenIndex, depth = position373, tokenIndex373, depth373
			return false
		},
		/* 24 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position375, tokenIndex375, depth375 := position, tokenIndex, depth
			{
				position376 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l375
				}
				if !_rules[rulesp]() {
					goto l375
				}
			l377:
				{
					position378, tokenIndex378, depth378 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l378
					}
					position++
					if !_rules[rulesp]() {
						goto l378
					}
					if !_rules[ruleRelationLike]() {
						goto l378
					}
					goto l377
				l378:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
				}
				depth--
				add(ruleRelations, position376)
			}
			return true
		l375:
			position, tokenIndex, depth = position375, tokenIndex375, depth375
			return false
		},
		/* 25 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action19)> */
		func() bool {
			position379, tokenIndex379, depth379 := position, tokenIndex, depth
			{
				position380 := position
				depth++
				{
					position381 := position
					depth++
					{
						position382, tokenIndex382, depth382 := position, tokenIndex, depth
						{
							position384, tokenIndex384, depth384 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l385
							}
							position++
							goto l384
						l385:
							position, tokenIndex, depth = position384, tokenIndex384, depth384
							if buffer[position] != rune('W') {
								goto l382
							}
							position++
						}
					l384:
						{
							position386, tokenIndex386, depth386 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l387
							}
							position++
							goto l386
						l387:
							position, tokenIndex, depth = position386, tokenIndex386, depth386
							if buffer[position] != rune('H') {
								goto l382
							}
							position++
						}
					l386:
						{
							position388, tokenIndex388, depth388 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l389
							}
							position++
							goto l388
						l389:
							position, tokenIndex, depth = position388, tokenIndex388, depth388
							if buffer[position] != rune('E') {
								goto l382
							}
							position++
						}
					l388:
						{
							position390, tokenIndex390, depth390 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l391
							}
							position++
							goto l390
						l391:
							position, tokenIndex, depth = position390, tokenIndex390, depth390
							if buffer[position] != rune('R') {
								goto l382
							}
							position++
						}
					l390:
						{
							position392, tokenIndex392, depth392 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l393
							}
							position++
							goto l392
						l393:
							position, tokenIndex, depth = position392, tokenIndex392, depth392
							if buffer[position] != rune('E') {
								goto l382
							}
							position++
						}
					l392:
						if !_rules[rulesp]() {
							goto l382
						}
						if !_rules[ruleExpression]() {
							goto l382
						}
						goto l383
					l382:
						position, tokenIndex, depth = position382, tokenIndex382, depth382
					}
				l383:
					depth--
					add(rulePegText, position381)
				}
				if !_rules[ruleAction19]() {
					goto l379
				}
				depth--
				add(ruleFilter, position380)
			}
			return true
		l379:
			position, tokenIndex, depth = position379, tokenIndex379, depth379
			return false
		},
		/* 26 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action20)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				{
					position396 := position
					depth++
					{
						position397, tokenIndex397, depth397 := position, tokenIndex, depth
						{
							position399, tokenIndex399, depth399 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l400
							}
							position++
							goto l399
						l400:
							position, tokenIndex, depth = position399, tokenIndex399, depth399
							if buffer[position] != rune('G') {
								goto l397
							}
							position++
						}
					l399:
						{
							position401, tokenIndex401, depth401 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l402
							}
							position++
							goto l401
						l402:
							position, tokenIndex, depth = position401, tokenIndex401, depth401
							if buffer[position] != rune('R') {
								goto l397
							}
							position++
						}
					l401:
						{
							position403, tokenIndex403, depth403 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l404
							}
							position++
							goto l403
						l404:
							position, tokenIndex, depth = position403, tokenIndex403, depth403
							if buffer[position] != rune('O') {
								goto l397
							}
							position++
						}
					l403:
						{
							position405, tokenIndex405, depth405 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l406
							}
							position++
							goto l405
						l406:
							position, tokenIndex, depth = position405, tokenIndex405, depth405
							if buffer[position] != rune('U') {
								goto l397
							}
							position++
						}
					l405:
						{
							position407, tokenIndex407, depth407 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l408
							}
							position++
							goto l407
						l408:
							position, tokenIndex, depth = position407, tokenIndex407, depth407
							if buffer[position] != rune('P') {
								goto l397
							}
							position++
						}
					l407:
						if !_rules[rulesp]() {
							goto l397
						}
						{
							position409, tokenIndex409, depth409 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l410
							}
							position++
							goto l409
						l410:
							position, tokenIndex, depth = position409, tokenIndex409, depth409
							if buffer[position] != rune('B') {
								goto l397
							}
							position++
						}
					l409:
						{
							position411, tokenIndex411, depth411 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l412
							}
							position++
							goto l411
						l412:
							position, tokenIndex, depth = position411, tokenIndex411, depth411
							if buffer[position] != rune('Y') {
								goto l397
							}
							position++
						}
					l411:
						if !_rules[rulesp]() {
							goto l397
						}
						if !_rules[ruleGroupList]() {
							goto l397
						}
						goto l398
					l397:
						position, tokenIndex, depth = position397, tokenIndex397, depth397
					}
				l398:
					depth--
					add(rulePegText, position396)
				}
				if !_rules[ruleAction20]() {
					goto l394
				}
				depth--
				add(ruleGrouping, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 27 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l413
				}
				if !_rules[rulesp]() {
					goto l413
				}
			l415:
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l416
					}
					position++
					if !_rules[rulesp]() {
						goto l416
					}
					if !_rules[ruleExpression]() {
						goto l416
					}
					goto l415
				l416:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
				}
				depth--
				add(ruleGroupList, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 28 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action21)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				{
					position419 := position
					depth++
					{
						position420, tokenIndex420, depth420 := position, tokenIndex, depth
						{
							position422, tokenIndex422, depth422 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l423
							}
							position++
							goto l422
						l423:
							position, tokenIndex, depth = position422, tokenIndex422, depth422
							if buffer[position] != rune('H') {
								goto l420
							}
							position++
						}
					l422:
						{
							position424, tokenIndex424, depth424 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l425
							}
							position++
							goto l424
						l425:
							position, tokenIndex, depth = position424, tokenIndex424, depth424
							if buffer[position] != rune('A') {
								goto l420
							}
							position++
						}
					l424:
						{
							position426, tokenIndex426, depth426 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l427
							}
							position++
							goto l426
						l427:
							position, tokenIndex, depth = position426, tokenIndex426, depth426
							if buffer[position] != rune('V') {
								goto l420
							}
							position++
						}
					l426:
						{
							position428, tokenIndex428, depth428 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l429
							}
							position++
							goto l428
						l429:
							position, tokenIndex, depth = position428, tokenIndex428, depth428
							if buffer[position] != rune('I') {
								goto l420
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
								goto l420
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
								goto l420
							}
							position++
						}
					l432:
						if !_rules[rulesp]() {
							goto l420
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction21]() {
					goto l417
				}
				depth--
				add(ruleHaving, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 29 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action22))> */
		func() bool {
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
			{
				position435 := position
				depth++
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l437
					}
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if !_rules[ruleStreamWindow]() {
						goto l434
					}
					if !_rules[ruleAction22]() {
						goto l434
					}
				}
			l436:
				depth--
				add(ruleRelationLike, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
			return false
		},
		/* 30 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action23)> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l438
				}
				if !_rules[rulesp]() {
					goto l438
				}
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
						goto l438
					}
					position++
				}
			l440:
				{
					position442, tokenIndex442, depth442 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l443
					}
					position++
					goto l442
				l443:
					position, tokenIndex, depth = position442, tokenIndex442, depth442
					if buffer[position] != rune('S') {
						goto l438
					}
					position++
				}
			l442:
				if !_rules[rulesp]() {
					goto l438
				}
				if !_rules[ruleIdentifier]() {
					goto l438
				}
				if !_rules[ruleAction23]() {
					goto l438
				}
				depth--
				add(ruleAliasedStreamWindow, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 31 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action24)> */
		func() bool {
			position444, tokenIndex444, depth444 := position, tokenIndex, depth
			{
				position445 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l444
				}
				if !_rules[rulesp]() {
					goto l444
				}
				if buffer[position] != rune('[') {
					goto l444
				}
				position++
				if !_rules[rulesp]() {
					goto l444
				}
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
					if buffer[position] != rune('R') {
						goto l444
					}
					position++
				}
			l446:
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
					if buffer[position] != rune('A') {
						goto l444
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('N') {
						goto l444
					}
					position++
				}
			l450:
				{
					position452, tokenIndex452, depth452 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l453
					}
					position++
					goto l452
				l453:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if buffer[position] != rune('G') {
						goto l444
					}
					position++
				}
			l452:
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l455
					}
					position++
					goto l454
				l455:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if buffer[position] != rune('E') {
						goto l444
					}
					position++
				}
			l454:
				if !_rules[rulesp]() {
					goto l444
				}
				if !_rules[ruleInterval]() {
					goto l444
				}
				if !_rules[rulesp]() {
					goto l444
				}
				if buffer[position] != rune(']') {
					goto l444
				}
				position++
				if !_rules[ruleAction24]() {
					goto l444
				}
				depth--
				add(ruleStreamWindow, position445)
			}
			return true
		l444:
			position, tokenIndex, depth = position444, tokenIndex444, depth444
			return false
		},
		/* 32 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l459
					}
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !_rules[ruleStream]() {
						goto l456
					}
				}
			l458:
				depth--
				add(ruleStreamLike, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 33 UDSFFuncApp <- <(FuncApp Action25)> */
		func() bool {
			position460, tokenIndex460, depth460 := position, tokenIndex, depth
			{
				position461 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l460
				}
				if !_rules[ruleAction25]() {
					goto l460
				}
				depth--
				add(ruleUDSFFuncApp, position461)
			}
			return true
		l460:
			position, tokenIndex, depth = position460, tokenIndex460, depth460
			return false
		},
		/* 34 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action26)> */
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
							if buffer[position] != rune('w') {
								goto l468
							}
							position++
							goto l467
						l468:
							position, tokenIndex, depth = position467, tokenIndex467, depth467
							if buffer[position] != rune('W') {
								goto l465
							}
							position++
						}
					l467:
						{
							position469, tokenIndex469, depth469 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l470
							}
							position++
							goto l469
						l470:
							position, tokenIndex, depth = position469, tokenIndex469, depth469
							if buffer[position] != rune('I') {
								goto l465
							}
							position++
						}
					l469:
						{
							position471, tokenIndex471, depth471 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l472
							}
							position++
							goto l471
						l472:
							position, tokenIndex, depth = position471, tokenIndex471, depth471
							if buffer[position] != rune('T') {
								goto l465
							}
							position++
						}
					l471:
						{
							position473, tokenIndex473, depth473 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l474
							}
							position++
							goto l473
						l474:
							position, tokenIndex, depth = position473, tokenIndex473, depth473
							if buffer[position] != rune('H') {
								goto l465
							}
							position++
						}
					l473:
						if !_rules[rulesp]() {
							goto l465
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l465
						}
						if !_rules[rulesp]() {
							goto l465
						}
					l475:
						{
							position476, tokenIndex476, depth476 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l476
							}
							position++
							if !_rules[rulesp]() {
								goto l476
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l476
							}
							goto l475
						l476:
							position, tokenIndex, depth = position476, tokenIndex476, depth476
						}
						goto l466
					l465:
						position, tokenIndex, depth = position465, tokenIndex465, depth465
					}
				l466:
					depth--
					add(rulePegText, position464)
				}
				if !_rules[ruleAction26]() {
					goto l462
				}
				depth--
				add(ruleSourceSinkSpecs, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 35 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action27)> */
		func() bool {
			position477, tokenIndex477, depth477 := position, tokenIndex, depth
			{
				position478 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l477
				}
				if buffer[position] != rune('=') {
					goto l477
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l477
				}
				if !_rules[ruleAction27]() {
					goto l477
				}
				depth--
				add(ruleSourceSinkParam, position478)
			}
			return true
		l477:
			position, tokenIndex, depth = position477, tokenIndex477, depth477
			return false
		},
		/* 36 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l482
					}
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleLiteral]() {
						goto l479
					}
				}
			l481:
				depth--
				add(ruleSourceSinkParamVal, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 37 PausedOpt <- <(<(Paused / Unpaused)?> Action28)> */
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
							if !_rules[rulePaused]() {
								goto l489
							}
							goto l488
						l489:
							position, tokenIndex, depth = position488, tokenIndex488, depth488
							if !_rules[ruleUnpaused]() {
								goto l486
							}
						}
					l488:
						goto l487
					l486:
						position, tokenIndex, depth = position486, tokenIndex486, depth486
					}
				l487:
					depth--
					add(rulePegText, position485)
				}
				if !_rules[ruleAction28]() {
					goto l483
				}
				depth--
				add(rulePausedOpt, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 38 Expression <- <orExpr> */
		func() bool {
			position490, tokenIndex490, depth490 := position, tokenIndex, depth
			{
				position491 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l490
				}
				depth--
				add(ruleExpression, position491)
			}
			return true
		l490:
			position, tokenIndex, depth = position490, tokenIndex490, depth490
			return false
		},
		/* 39 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action29)> */
		func() bool {
			position492, tokenIndex492, depth492 := position, tokenIndex, depth
			{
				position493 := position
				depth++
				{
					position494 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l492
					}
					if !_rules[rulesp]() {
						goto l492
					}
					{
						position495, tokenIndex495, depth495 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l495
						}
						if !_rules[rulesp]() {
							goto l495
						}
						if !_rules[ruleandExpr]() {
							goto l495
						}
						goto l496
					l495:
						position, tokenIndex, depth = position495, tokenIndex495, depth495
					}
				l496:
					depth--
					add(rulePegText, position494)
				}
				if !_rules[ruleAction29]() {
					goto l492
				}
				depth--
				add(ruleorExpr, position493)
			}
			return true
		l492:
			position, tokenIndex, depth = position492, tokenIndex492, depth492
			return false
		},
		/* 40 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action30)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				{
					position499 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l497
					}
					if !_rules[rulesp]() {
						goto l497
					}
					{
						position500, tokenIndex500, depth500 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l500
						}
						if !_rules[rulesp]() {
							goto l500
						}
						if !_rules[rulenotExpr]() {
							goto l500
						}
						goto l501
					l500:
						position, tokenIndex, depth = position500, tokenIndex500, depth500
					}
				l501:
					depth--
					add(rulePegText, position499)
				}
				if !_rules[ruleAction30]() {
					goto l497
				}
				depth--
				add(ruleandExpr, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 41 notExpr <- <(<((Not sp)? comparisonExpr)> Action31)> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				{
					position504 := position
					depth++
					{
						position505, tokenIndex505, depth505 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l505
						}
						if !_rules[rulesp]() {
							goto l505
						}
						goto l506
					l505:
						position, tokenIndex, depth = position505, tokenIndex505, depth505
					}
				l506:
					if !_rules[rulecomparisonExpr]() {
						goto l502
					}
					depth--
					add(rulePegText, position504)
				}
				if !_rules[ruleAction31]() {
					goto l502
				}
				depth--
				add(rulenotExpr, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 42 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action32)> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				{
					position509 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l507
					}
					if !_rules[rulesp]() {
						goto l507
					}
					{
						position510, tokenIndex510, depth510 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l510
						}
						if !_rules[rulesp]() {
							goto l510
						}
						if !_rules[ruleisExpr]() {
							goto l510
						}
						goto l511
					l510:
						position, tokenIndex, depth = position510, tokenIndex510, depth510
					}
				l511:
					depth--
					add(rulePegText, position509)
				}
				if !_rules[ruleAction32]() {
					goto l507
				}
				depth--
				add(rulecomparisonExpr, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 43 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action33)> */
		func() bool {
			position512, tokenIndex512, depth512 := position, tokenIndex, depth
			{
				position513 := position
				depth++
				{
					position514 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l512
					}
					if !_rules[rulesp]() {
						goto l512
					}
					{
						position515, tokenIndex515, depth515 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l515
						}
						if !_rules[rulesp]() {
							goto l515
						}
						if !_rules[ruleNullLiteral]() {
							goto l515
						}
						goto l516
					l515:
						position, tokenIndex, depth = position515, tokenIndex515, depth515
					}
				l516:
					depth--
					add(rulePegText, position514)
				}
				if !_rules[ruleAction33]() {
					goto l512
				}
				depth--
				add(ruleisExpr, position513)
			}
			return true
		l512:
			position, tokenIndex, depth = position512, tokenIndex512, depth512
			return false
		},
		/* 44 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action34)> */
		func() bool {
			position517, tokenIndex517, depth517 := position, tokenIndex, depth
			{
				position518 := position
				depth++
				{
					position519 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l517
					}
					if !_rules[rulesp]() {
						goto l517
					}
					{
						position520, tokenIndex520, depth520 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l520
						}
						if !_rules[rulesp]() {
							goto l520
						}
						if !_rules[ruleproductExpr]() {
							goto l520
						}
						goto l521
					l520:
						position, tokenIndex, depth = position520, tokenIndex520, depth520
					}
				l521:
					depth--
					add(rulePegText, position519)
				}
				if !_rules[ruleAction34]() {
					goto l517
				}
				depth--
				add(ruletermExpr, position518)
			}
			return true
		l517:
			position, tokenIndex, depth = position517, tokenIndex517, depth517
			return false
		},
		/* 45 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action35)> */
		func() bool {
			position522, tokenIndex522, depth522 := position, tokenIndex, depth
			{
				position523 := position
				depth++
				{
					position524 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l522
					}
					if !_rules[rulesp]() {
						goto l522
					}
					{
						position525, tokenIndex525, depth525 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l525
						}
						if !_rules[rulesp]() {
							goto l525
						}
						if !_rules[rulebaseExpr]() {
							goto l525
						}
						goto l526
					l525:
						position, tokenIndex, depth = position525, tokenIndex525, depth525
					}
				l526:
					depth--
					add(rulePegText, position524)
				}
				if !_rules[ruleAction35]() {
					goto l522
				}
				depth--
				add(ruleproductExpr, position523)
			}
			return true
		l522:
			position, tokenIndex, depth = position522, tokenIndex522, depth522
			return false
		},
		/* 46 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position527, tokenIndex527, depth527 := position, tokenIndex, depth
			{
				position528 := position
				depth++
				{
					position529, tokenIndex529, depth529 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l530
					}
					position++
					if !_rules[rulesp]() {
						goto l530
					}
					if !_rules[ruleExpression]() {
						goto l530
					}
					if !_rules[rulesp]() {
						goto l530
					}
					if buffer[position] != rune(')') {
						goto l530
					}
					position++
					goto l529
				l530:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if !_rules[ruleBooleanLiteral]() {
						goto l531
					}
					goto l529
				l531:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if !_rules[ruleNullLiteral]() {
						goto l532
					}
					goto l529
				l532:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if !_rules[ruleFuncApp]() {
						goto l533
					}
					goto l529
				l533:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if !_rules[ruleRowMeta]() {
						goto l534
					}
					goto l529
				l534:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if !_rules[ruleRowValue]() {
						goto l535
					}
					goto l529
				l535:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if !_rules[ruleLiteral]() {
						goto l527
					}
				}
			l529:
				depth--
				add(rulebaseExpr, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 47 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action36)> */
		func() bool {
			position536, tokenIndex536, depth536 := position, tokenIndex, depth
			{
				position537 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l536
				}
				if !_rules[rulesp]() {
					goto l536
				}
				if buffer[position] != rune('(') {
					goto l536
				}
				position++
				if !_rules[rulesp]() {
					goto l536
				}
				if !_rules[ruleFuncParams]() {
					goto l536
				}
				if !_rules[rulesp]() {
					goto l536
				}
				if buffer[position] != rune(')') {
					goto l536
				}
				position++
				if !_rules[ruleAction36]() {
					goto l536
				}
				depth--
				add(ruleFuncApp, position537)
			}
			return true
		l536:
			position, tokenIndex, depth = position536, tokenIndex536, depth536
			return false
		},
		/* 48 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action37)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					position540 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l538
					}
					if !_rules[rulesp]() {
						goto l538
					}
				l541:
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l542
						}
						position++
						if !_rules[rulesp]() {
							goto l542
						}
						if !_rules[ruleExpression]() {
							goto l542
						}
						goto l541
					l542:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
					}
					depth--
					add(rulePegText, position540)
				}
				if !_rules[ruleAction37]() {
					goto l538
				}
				depth--
				add(ruleFuncParams, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 49 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position543, tokenIndex543, depth543 := position, tokenIndex, depth
			{
				position544 := position
				depth++
				{
					position545, tokenIndex545, depth545 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l546
					}
					goto l545
				l546:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleNumericLiteral]() {
						goto l547
					}
					goto l545
				l547:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleStringLiteral]() {
						goto l543
					}
				}
			l545:
				depth--
				add(ruleLiteral, position544)
			}
			return true
		l543:
			position, tokenIndex, depth = position543, tokenIndex543, depth543
			return false
		},
		/* 50 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position548, tokenIndex548, depth548 := position, tokenIndex, depth
			{
				position549 := position
				depth++
				{
					position550, tokenIndex550, depth550 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l551
					}
					goto l550
				l551:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleNotEqual]() {
						goto l552
					}
					goto l550
				l552:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleLessOrEqual]() {
						goto l553
					}
					goto l550
				l553:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleLess]() {
						goto l554
					}
					goto l550
				l554:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleGreaterOrEqual]() {
						goto l555
					}
					goto l550
				l555:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleGreater]() {
						goto l556
					}
					goto l550
				l556:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleNotEqual]() {
						goto l548
					}
				}
			l550:
				depth--
				add(ruleComparisonOp, position549)
			}
			return true
		l548:
			position, tokenIndex, depth = position548, tokenIndex548, depth548
			return false
		},
		/* 51 IsOp <- <(IsNot / Is)> */
		func() bool {
			position557, tokenIndex557, depth557 := position, tokenIndex, depth
			{
				position558 := position
				depth++
				{
					position559, tokenIndex559, depth559 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l560
					}
					goto l559
				l560:
					position, tokenIndex, depth = position559, tokenIndex559, depth559
					if !_rules[ruleIs]() {
						goto l557
					}
				}
			l559:
				depth--
				add(ruleIsOp, position558)
			}
			return true
		l557:
			position, tokenIndex, depth = position557, tokenIndex557, depth557
			return false
		},
		/* 52 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position561, tokenIndex561, depth561 := position, tokenIndex, depth
			{
				position562 := position
				depth++
				{
					position563, tokenIndex563, depth563 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l564
					}
					goto l563
				l564:
					position, tokenIndex, depth = position563, tokenIndex563, depth563
					if !_rules[ruleMinus]() {
						goto l561
					}
				}
			l563:
				depth--
				add(rulePlusMinusOp, position562)
			}
			return true
		l561:
			position, tokenIndex, depth = position561, tokenIndex561, depth561
			return false
		},
		/* 53 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567, tokenIndex567, depth567 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l568
					}
					goto l567
				l568:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleDivide]() {
						goto l569
					}
					goto l567
				l569:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleModulo]() {
						goto l565
					}
				}
			l567:
				depth--
				add(ruleMultDivOp, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 54 Stream <- <(<ident> Action38)> */
		func() bool {
			position570, tokenIndex570, depth570 := position, tokenIndex, depth
			{
				position571 := position
				depth++
				{
					position572 := position
					depth++
					if !_rules[ruleident]() {
						goto l570
					}
					depth--
					add(rulePegText, position572)
				}
				if !_rules[ruleAction38]() {
					goto l570
				}
				depth--
				add(ruleStream, position571)
			}
			return true
		l570:
			position, tokenIndex, depth = position570, tokenIndex570, depth570
			return false
		},
		/* 55 RowMeta <- <RowTimestamp> */
		func() bool {
			position573, tokenIndex573, depth573 := position, tokenIndex, depth
			{
				position574 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l573
				}
				depth--
				add(ruleRowMeta, position574)
			}
			return true
		l573:
			position, tokenIndex, depth = position573, tokenIndex573, depth573
			return false
		},
		/* 56 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action39)> */
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
						if !_rules[ruleident]() {
							goto l578
						}
						if buffer[position] != rune(':') {
							goto l578
						}
						position++
						goto l579
					l578:
						position, tokenIndex, depth = position578, tokenIndex578, depth578
					}
				l579:
					if buffer[position] != rune('t') {
						goto l575
					}
					position++
					if buffer[position] != rune('s') {
						goto l575
					}
					position++
					if buffer[position] != rune('(') {
						goto l575
					}
					position++
					if buffer[position] != rune(')') {
						goto l575
					}
					position++
					depth--
					add(rulePegText, position577)
				}
				if !_rules[ruleAction39]() {
					goto l575
				}
				depth--
				add(ruleRowTimestamp, position576)
			}
			return true
		l575:
			position, tokenIndex, depth = position575, tokenIndex575, depth575
			return false
		},
		/* 57 RowValue <- <(<((ident ':')? jsonPath)> Action40)> */
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
						if !_rules[ruleident]() {
							goto l583
						}
						if buffer[position] != rune(':') {
							goto l583
						}
						position++
						goto l584
					l583:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
					}
				l584:
					if !_rules[rulejsonPath]() {
						goto l580
					}
					depth--
					add(rulePegText, position582)
				}
				if !_rules[ruleAction40]() {
					goto l580
				}
				depth--
				add(ruleRowValue, position581)
			}
			return true
		l580:
			position, tokenIndex, depth = position580, tokenIndex580, depth580
			return false
		},
		/* 58 NumericLiteral <- <(<('-'? [0-9]+)> Action41)> */
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
				if !_rules[ruleAction41]() {
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
		/* 59 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action42)> */
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
				if !_rules[ruleAction42]() {
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
		/* 60 Function <- <(<ident> Action43)> */
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
				if !_rules[ruleAction43]() {
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
		/* 61 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action44)> */
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
						if buffer[position] != rune('n') {
							goto l608
						}
						position++
						goto l607
					l608:
						position, tokenIndex, depth = position607, tokenIndex607, depth607
						if buffer[position] != rune('N') {
							goto l604
						}
						position++
					}
				l607:
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('U') {
							goto l604
						}
						position++
					}
				l609:
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('L') {
							goto l604
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('L') {
							goto l604
						}
						position++
					}
				l613:
					depth--
					add(rulePegText, position606)
				}
				if !_rules[ruleAction44]() {
					goto l604
				}
				depth--
				add(ruleNullLiteral, position605)
			}
			return true
		l604:
			position, tokenIndex, depth = position604, tokenIndex604, depth604
			return false
		},
		/* 62 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position615, tokenIndex615, depth615 := position, tokenIndex, depth
			{
				position616 := position
				depth++
				{
					position617, tokenIndex617, depth617 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l618
					}
					goto l617
				l618:
					position, tokenIndex, depth = position617, tokenIndex617, depth617
					if !_rules[ruleFALSE]() {
						goto l615
					}
				}
			l617:
				depth--
				add(ruleBooleanLiteral, position616)
			}
			return true
		l615:
			position, tokenIndex, depth = position615, tokenIndex615, depth615
			return false
		},
		/* 63 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action45)> */
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
						if buffer[position] != rune('t') {
							goto l623
						}
						position++
						goto l622
					l623:
						position, tokenIndex, depth = position622, tokenIndex622, depth622
						if buffer[position] != rune('T') {
							goto l619
						}
						position++
					}
				l622:
					{
						position624, tokenIndex624, depth624 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('R') {
							goto l619
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
							goto l619
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
							goto l619
						}
						position++
					}
				l628:
					depth--
					add(rulePegText, position621)
				}
				if !_rules[ruleAction45]() {
					goto l619
				}
				depth--
				add(ruleTRUE, position620)
			}
			return true
		l619:
			position, tokenIndex, depth = position619, tokenIndex619, depth619
			return false
		},
		/* 64 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action46)> */
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
						if buffer[position] != rune('f') {
							goto l634
						}
						position++
						goto l633
					l634:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
						if buffer[position] != rune('F') {
							goto l630
						}
						position++
					}
				l633:
					{
						position635, tokenIndex635, depth635 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l636
						}
						position++
						goto l635
					l636:
						position, tokenIndex, depth = position635, tokenIndex635, depth635
						if buffer[position] != rune('A') {
							goto l630
						}
						position++
					}
				l635:
					{
						position637, tokenIndex637, depth637 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l638
						}
						position++
						goto l637
					l638:
						position, tokenIndex, depth = position637, tokenIndex637, depth637
						if buffer[position] != rune('L') {
							goto l630
						}
						position++
					}
				l637:
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
							goto l630
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
							goto l630
						}
						position++
					}
				l641:
					depth--
					add(rulePegText, position632)
				}
				if !_rules[ruleAction46]() {
					goto l630
				}
				depth--
				add(ruleFALSE, position631)
			}
			return true
		l630:
			position, tokenIndex, depth = position630, tokenIndex630, depth630
			return false
		},
		/* 65 Wildcard <- <(<'*'> Action47)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645 := position
					depth++
					if buffer[position] != rune('*') {
						goto l643
					}
					position++
					depth--
					add(rulePegText, position645)
				}
				if !_rules[ruleAction47]() {
					goto l643
				}
				depth--
				add(ruleWildcard, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 66 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action48)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				{
					position648 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l646
					}
					position++
				l649:
					{
						position650, tokenIndex650, depth650 := position, tokenIndex, depth
						{
							position651, tokenIndex651, depth651 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l652
							}
							position++
							if buffer[position] != rune('\'') {
								goto l652
							}
							position++
							goto l651
						l652:
							position, tokenIndex, depth = position651, tokenIndex651, depth651
							{
								position653, tokenIndex653, depth653 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l653
								}
								position++
								goto l650
							l653:
								position, tokenIndex, depth = position653, tokenIndex653, depth653
							}
							if !matchDot() {
								goto l650
							}
						}
					l651:
						goto l649
					l650:
						position, tokenIndex, depth = position650, tokenIndex650, depth650
					}
					if buffer[position] != rune('\'') {
						goto l646
					}
					position++
					depth--
					add(rulePegText, position648)
				}
				if !_rules[ruleAction48]() {
					goto l646
				}
				depth--
				add(ruleStringLiteral, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 67 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action49)> */
		func() bool {
			position654, tokenIndex654, depth654 := position, tokenIndex, depth
			{
				position655 := position
				depth++
				{
					position656 := position
					depth++
					{
						position657, tokenIndex657, depth657 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l658
						}
						position++
						goto l657
					l658:
						position, tokenIndex, depth = position657, tokenIndex657, depth657
						if buffer[position] != rune('I') {
							goto l654
						}
						position++
					}
				l657:
					{
						position659, tokenIndex659, depth659 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l660
						}
						position++
						goto l659
					l660:
						position, tokenIndex, depth = position659, tokenIndex659, depth659
						if buffer[position] != rune('S') {
							goto l654
						}
						position++
					}
				l659:
					{
						position661, tokenIndex661, depth661 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l662
						}
						position++
						goto l661
					l662:
						position, tokenIndex, depth = position661, tokenIndex661, depth661
						if buffer[position] != rune('T') {
							goto l654
						}
						position++
					}
				l661:
					{
						position663, tokenIndex663, depth663 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l664
						}
						position++
						goto l663
					l664:
						position, tokenIndex, depth = position663, tokenIndex663, depth663
						if buffer[position] != rune('R') {
							goto l654
						}
						position++
					}
				l663:
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l666
						}
						position++
						goto l665
					l666:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if buffer[position] != rune('E') {
							goto l654
						}
						position++
					}
				l665:
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if buffer[position] != rune('A') {
							goto l654
						}
						position++
					}
				l667:
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l670
						}
						position++
						goto l669
					l670:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if buffer[position] != rune('M') {
							goto l654
						}
						position++
					}
				l669:
					depth--
					add(rulePegText, position656)
				}
				if !_rules[ruleAction49]() {
					goto l654
				}
				depth--
				add(ruleISTREAM, position655)
			}
			return true
		l654:
			position, tokenIndex, depth = position654, tokenIndex654, depth654
			return false
		},
		/* 68 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action50)> */
		func() bool {
			position671, tokenIndex671, depth671 := position, tokenIndex, depth
			{
				position672 := position
				depth++
				{
					position673 := position
					depth++
					{
						position674, tokenIndex674, depth674 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l675
						}
						position++
						goto l674
					l675:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
						if buffer[position] != rune('D') {
							goto l671
						}
						position++
					}
				l674:
					{
						position676, tokenIndex676, depth676 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l677
						}
						position++
						goto l676
					l677:
						position, tokenIndex, depth = position676, tokenIndex676, depth676
						if buffer[position] != rune('S') {
							goto l671
						}
						position++
					}
				l676:
					{
						position678, tokenIndex678, depth678 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l679
						}
						position++
						goto l678
					l679:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
						if buffer[position] != rune('T') {
							goto l671
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
							goto l671
						}
						position++
					}
				l680:
					{
						position682, tokenIndex682, depth682 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l683
						}
						position++
						goto l682
					l683:
						position, tokenIndex, depth = position682, tokenIndex682, depth682
						if buffer[position] != rune('E') {
							goto l671
						}
						position++
					}
				l682:
					{
						position684, tokenIndex684, depth684 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l685
						}
						position++
						goto l684
					l685:
						position, tokenIndex, depth = position684, tokenIndex684, depth684
						if buffer[position] != rune('A') {
							goto l671
						}
						position++
					}
				l684:
					{
						position686, tokenIndex686, depth686 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
						if buffer[position] != rune('M') {
							goto l671
						}
						position++
					}
				l686:
					depth--
					add(rulePegText, position673)
				}
				if !_rules[ruleAction50]() {
					goto l671
				}
				depth--
				add(ruleDSTREAM, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 69 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action51)> */
		func() bool {
			position688, tokenIndex688, depth688 := position, tokenIndex, depth
			{
				position689 := position
				depth++
				{
					position690 := position
					depth++
					{
						position691, tokenIndex691, depth691 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l692
						}
						position++
						goto l691
					l692:
						position, tokenIndex, depth = position691, tokenIndex691, depth691
						if buffer[position] != rune('R') {
							goto l688
						}
						position++
					}
				l691:
					{
						position693, tokenIndex693, depth693 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l694
						}
						position++
						goto l693
					l694:
						position, tokenIndex, depth = position693, tokenIndex693, depth693
						if buffer[position] != rune('S') {
							goto l688
						}
						position++
					}
				l693:
					{
						position695, tokenIndex695, depth695 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
						if buffer[position] != rune('T') {
							goto l688
						}
						position++
					}
				l695:
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('R') {
							goto l688
						}
						position++
					}
				l697:
					{
						position699, tokenIndex699, depth699 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l700
						}
						position++
						goto l699
					l700:
						position, tokenIndex, depth = position699, tokenIndex699, depth699
						if buffer[position] != rune('E') {
							goto l688
						}
						position++
					}
				l699:
					{
						position701, tokenIndex701, depth701 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l702
						}
						position++
						goto l701
					l702:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
						if buffer[position] != rune('A') {
							goto l688
						}
						position++
					}
				l701:
					{
						position703, tokenIndex703, depth703 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l704
						}
						position++
						goto l703
					l704:
						position, tokenIndex, depth = position703, tokenIndex703, depth703
						if buffer[position] != rune('M') {
							goto l688
						}
						position++
					}
				l703:
					depth--
					add(rulePegText, position690)
				}
				if !_rules[ruleAction51]() {
					goto l688
				}
				depth--
				add(ruleRSTREAM, position689)
			}
			return true
		l688:
			position, tokenIndex, depth = position688, tokenIndex688, depth688
			return false
		},
		/* 70 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action52)> */
		func() bool {
			position705, tokenIndex705, depth705 := position, tokenIndex, depth
			{
				position706 := position
				depth++
				{
					position707 := position
					depth++
					{
						position708, tokenIndex708, depth708 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l709
						}
						position++
						goto l708
					l709:
						position, tokenIndex, depth = position708, tokenIndex708, depth708
						if buffer[position] != rune('T') {
							goto l705
						}
						position++
					}
				l708:
					{
						position710, tokenIndex710, depth710 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l711
						}
						position++
						goto l710
					l711:
						position, tokenIndex, depth = position710, tokenIndex710, depth710
						if buffer[position] != rune('U') {
							goto l705
						}
						position++
					}
				l710:
					{
						position712, tokenIndex712, depth712 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l713
						}
						position++
						goto l712
					l713:
						position, tokenIndex, depth = position712, tokenIndex712, depth712
						if buffer[position] != rune('P') {
							goto l705
						}
						position++
					}
				l712:
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l715
						}
						position++
						goto l714
					l715:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
						if buffer[position] != rune('L') {
							goto l705
						}
						position++
					}
				l714:
					{
						position716, tokenIndex716, depth716 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l717
						}
						position++
						goto l716
					l717:
						position, tokenIndex, depth = position716, tokenIndex716, depth716
						if buffer[position] != rune('E') {
							goto l705
						}
						position++
					}
				l716:
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
							goto l705
						}
						position++
					}
				l718:
					depth--
					add(rulePegText, position707)
				}
				if !_rules[ruleAction52]() {
					goto l705
				}
				depth--
				add(ruleTUPLES, position706)
			}
			return true
		l705:
			position, tokenIndex, depth = position705, tokenIndex705, depth705
			return false
		},
		/* 71 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action53)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722 := position
					depth++
					{
						position723, tokenIndex723, depth723 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l724
						}
						position++
						goto l723
					l724:
						position, tokenIndex, depth = position723, tokenIndex723, depth723
						if buffer[position] != rune('S') {
							goto l720
						}
						position++
					}
				l723:
					{
						position725, tokenIndex725, depth725 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l726
						}
						position++
						goto l725
					l726:
						position, tokenIndex, depth = position725, tokenIndex725, depth725
						if buffer[position] != rune('E') {
							goto l720
						}
						position++
					}
				l725:
					{
						position727, tokenIndex727, depth727 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l728
						}
						position++
						goto l727
					l728:
						position, tokenIndex, depth = position727, tokenIndex727, depth727
						if buffer[position] != rune('C') {
							goto l720
						}
						position++
					}
				l727:
					{
						position729, tokenIndex729, depth729 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l730
						}
						position++
						goto l729
					l730:
						position, tokenIndex, depth = position729, tokenIndex729, depth729
						if buffer[position] != rune('O') {
							goto l720
						}
						position++
					}
				l729:
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l732
						}
						position++
						goto l731
					l732:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
						if buffer[position] != rune('N') {
							goto l720
						}
						position++
					}
				l731:
					{
						position733, tokenIndex733, depth733 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l734
						}
						position++
						goto l733
					l734:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
						if buffer[position] != rune('D') {
							goto l720
						}
						position++
					}
				l733:
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l736
						}
						position++
						goto l735
					l736:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if buffer[position] != rune('S') {
							goto l720
						}
						position++
					}
				l735:
					depth--
					add(rulePegText, position722)
				}
				if !_rules[ruleAction53]() {
					goto l720
				}
				depth--
				add(ruleSECONDS, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 72 StreamIdentifier <- <(<ident> Action54)> */
		func() bool {
			position737, tokenIndex737, depth737 := position, tokenIndex, depth
			{
				position738 := position
				depth++
				{
					position739 := position
					depth++
					if !_rules[ruleident]() {
						goto l737
					}
					depth--
					add(rulePegText, position739)
				}
				if !_rules[ruleAction54]() {
					goto l737
				}
				depth--
				add(ruleStreamIdentifier, position738)
			}
			return true
		l737:
			position, tokenIndex, depth = position737, tokenIndex737, depth737
			return false
		},
		/* 73 SourceSinkType <- <(<ident> Action55)> */
		func() bool {
			position740, tokenIndex740, depth740 := position, tokenIndex, depth
			{
				position741 := position
				depth++
				{
					position742 := position
					depth++
					if !_rules[ruleident]() {
						goto l740
					}
					depth--
					add(rulePegText, position742)
				}
				if !_rules[ruleAction55]() {
					goto l740
				}
				depth--
				add(ruleSourceSinkType, position741)
			}
			return true
		l740:
			position, tokenIndex, depth = position740, tokenIndex740, depth740
			return false
		},
		/* 74 SourceSinkParamKey <- <(<ident> Action56)> */
		func() bool {
			position743, tokenIndex743, depth743 := position, tokenIndex, depth
			{
				position744 := position
				depth++
				{
					position745 := position
					depth++
					if !_rules[ruleident]() {
						goto l743
					}
					depth--
					add(rulePegText, position745)
				}
				if !_rules[ruleAction56]() {
					goto l743
				}
				depth--
				add(ruleSourceSinkParamKey, position744)
			}
			return true
		l743:
			position, tokenIndex, depth = position743, tokenIndex743, depth743
			return false
		},
		/* 75 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action57)> */
		func() bool {
			position746, tokenIndex746, depth746 := position, tokenIndex, depth
			{
				position747 := position
				depth++
				{
					position748 := position
					depth++
					{
						position749, tokenIndex749, depth749 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l750
						}
						position++
						goto l749
					l750:
						position, tokenIndex, depth = position749, tokenIndex749, depth749
						if buffer[position] != rune('P') {
							goto l746
						}
						position++
					}
				l749:
					{
						position751, tokenIndex751, depth751 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l752
						}
						position++
						goto l751
					l752:
						position, tokenIndex, depth = position751, tokenIndex751, depth751
						if buffer[position] != rune('A') {
							goto l746
						}
						position++
					}
				l751:
					{
						position753, tokenIndex753, depth753 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l754
						}
						position++
						goto l753
					l754:
						position, tokenIndex, depth = position753, tokenIndex753, depth753
						if buffer[position] != rune('U') {
							goto l746
						}
						position++
					}
				l753:
					{
						position755, tokenIndex755, depth755 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l756
						}
						position++
						goto l755
					l756:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						if buffer[position] != rune('S') {
							goto l746
						}
						position++
					}
				l755:
					{
						position757, tokenIndex757, depth757 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l758
						}
						position++
						goto l757
					l758:
						position, tokenIndex, depth = position757, tokenIndex757, depth757
						if buffer[position] != rune('E') {
							goto l746
						}
						position++
					}
				l757:
					{
						position759, tokenIndex759, depth759 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l760
						}
						position++
						goto l759
					l760:
						position, tokenIndex, depth = position759, tokenIndex759, depth759
						if buffer[position] != rune('D') {
							goto l746
						}
						position++
					}
				l759:
					depth--
					add(rulePegText, position748)
				}
				if !_rules[ruleAction57]() {
					goto l746
				}
				depth--
				add(rulePaused, position747)
			}
			return true
		l746:
			position, tokenIndex, depth = position746, tokenIndex746, depth746
			return false
		},
		/* 76 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action58)> */
		func() bool {
			position761, tokenIndex761, depth761 := position, tokenIndex, depth
			{
				position762 := position
				depth++
				{
					position763 := position
					depth++
					{
						position764, tokenIndex764, depth764 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l765
						}
						position++
						goto l764
					l765:
						position, tokenIndex, depth = position764, tokenIndex764, depth764
						if buffer[position] != rune('U') {
							goto l761
						}
						position++
					}
				l764:
					{
						position766, tokenIndex766, depth766 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l767
						}
						position++
						goto l766
					l767:
						position, tokenIndex, depth = position766, tokenIndex766, depth766
						if buffer[position] != rune('N') {
							goto l761
						}
						position++
					}
				l766:
					{
						position768, tokenIndex768, depth768 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l769
						}
						position++
						goto l768
					l769:
						position, tokenIndex, depth = position768, tokenIndex768, depth768
						if buffer[position] != rune('P') {
							goto l761
						}
						position++
					}
				l768:
					{
						position770, tokenIndex770, depth770 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l771
						}
						position++
						goto l770
					l771:
						position, tokenIndex, depth = position770, tokenIndex770, depth770
						if buffer[position] != rune('A') {
							goto l761
						}
						position++
					}
				l770:
					{
						position772, tokenIndex772, depth772 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l773
						}
						position++
						goto l772
					l773:
						position, tokenIndex, depth = position772, tokenIndex772, depth772
						if buffer[position] != rune('U') {
							goto l761
						}
						position++
					}
				l772:
					{
						position774, tokenIndex774, depth774 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l775
						}
						position++
						goto l774
					l775:
						position, tokenIndex, depth = position774, tokenIndex774, depth774
						if buffer[position] != rune('S') {
							goto l761
						}
						position++
					}
				l774:
					{
						position776, tokenIndex776, depth776 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l777
						}
						position++
						goto l776
					l777:
						position, tokenIndex, depth = position776, tokenIndex776, depth776
						if buffer[position] != rune('E') {
							goto l761
						}
						position++
					}
				l776:
					{
						position778, tokenIndex778, depth778 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l779
						}
						position++
						goto l778
					l779:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						if buffer[position] != rune('D') {
							goto l761
						}
						position++
					}
				l778:
					depth--
					add(rulePegText, position763)
				}
				if !_rules[ruleAction58]() {
					goto l761
				}
				depth--
				add(ruleUnpaused, position762)
			}
			return true
		l761:
			position, tokenIndex, depth = position761, tokenIndex761, depth761
			return false
		},
		/* 77 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action59)> */
		func() bool {
			position780, tokenIndex780, depth780 := position, tokenIndex, depth
			{
				position781 := position
				depth++
				{
					position782 := position
					depth++
					{
						position783, tokenIndex783, depth783 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l784
						}
						position++
						goto l783
					l784:
						position, tokenIndex, depth = position783, tokenIndex783, depth783
						if buffer[position] != rune('O') {
							goto l780
						}
						position++
					}
				l783:
					{
						position785, tokenIndex785, depth785 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l786
						}
						position++
						goto l785
					l786:
						position, tokenIndex, depth = position785, tokenIndex785, depth785
						if buffer[position] != rune('R') {
							goto l780
						}
						position++
					}
				l785:
					depth--
					add(rulePegText, position782)
				}
				if !_rules[ruleAction59]() {
					goto l780
				}
				depth--
				add(ruleOr, position781)
			}
			return true
		l780:
			position, tokenIndex, depth = position780, tokenIndex780, depth780
			return false
		},
		/* 78 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action60)> */
		func() bool {
			position787, tokenIndex787, depth787 := position, tokenIndex, depth
			{
				position788 := position
				depth++
				{
					position789 := position
					depth++
					{
						position790, tokenIndex790, depth790 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l791
						}
						position++
						goto l790
					l791:
						position, tokenIndex, depth = position790, tokenIndex790, depth790
						if buffer[position] != rune('A') {
							goto l787
						}
						position++
					}
				l790:
					{
						position792, tokenIndex792, depth792 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l793
						}
						position++
						goto l792
					l793:
						position, tokenIndex, depth = position792, tokenIndex792, depth792
						if buffer[position] != rune('N') {
							goto l787
						}
						position++
					}
				l792:
					{
						position794, tokenIndex794, depth794 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l795
						}
						position++
						goto l794
					l795:
						position, tokenIndex, depth = position794, tokenIndex794, depth794
						if buffer[position] != rune('D') {
							goto l787
						}
						position++
					}
				l794:
					depth--
					add(rulePegText, position789)
				}
				if !_rules[ruleAction60]() {
					goto l787
				}
				depth--
				add(ruleAnd, position788)
			}
			return true
		l787:
			position, tokenIndex, depth = position787, tokenIndex787, depth787
			return false
		},
		/* 79 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action61)> */
		func() bool {
			position796, tokenIndex796, depth796 := position, tokenIndex, depth
			{
				position797 := position
				depth++
				{
					position798 := position
					depth++
					{
						position799, tokenIndex799, depth799 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l800
						}
						position++
						goto l799
					l800:
						position, tokenIndex, depth = position799, tokenIndex799, depth799
						if buffer[position] != rune('N') {
							goto l796
						}
						position++
					}
				l799:
					{
						position801, tokenIndex801, depth801 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l802
						}
						position++
						goto l801
					l802:
						position, tokenIndex, depth = position801, tokenIndex801, depth801
						if buffer[position] != rune('O') {
							goto l796
						}
						position++
					}
				l801:
					{
						position803, tokenIndex803, depth803 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l804
						}
						position++
						goto l803
					l804:
						position, tokenIndex, depth = position803, tokenIndex803, depth803
						if buffer[position] != rune('T') {
							goto l796
						}
						position++
					}
				l803:
					depth--
					add(rulePegText, position798)
				}
				if !_rules[ruleAction61]() {
					goto l796
				}
				depth--
				add(ruleNot, position797)
			}
			return true
		l796:
			position, tokenIndex, depth = position796, tokenIndex796, depth796
			return false
		},
		/* 80 Equal <- <(<'='> Action62)> */
		func() bool {
			position805, tokenIndex805, depth805 := position, tokenIndex, depth
			{
				position806 := position
				depth++
				{
					position807 := position
					depth++
					if buffer[position] != rune('=') {
						goto l805
					}
					position++
					depth--
					add(rulePegText, position807)
				}
				if !_rules[ruleAction62]() {
					goto l805
				}
				depth--
				add(ruleEqual, position806)
			}
			return true
		l805:
			position, tokenIndex, depth = position805, tokenIndex805, depth805
			return false
		},
		/* 81 Less <- <(<'<'> Action63)> */
		func() bool {
			position808, tokenIndex808, depth808 := position, tokenIndex, depth
			{
				position809 := position
				depth++
				{
					position810 := position
					depth++
					if buffer[position] != rune('<') {
						goto l808
					}
					position++
					depth--
					add(rulePegText, position810)
				}
				if !_rules[ruleAction63]() {
					goto l808
				}
				depth--
				add(ruleLess, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 82 LessOrEqual <- <(<('<' '=')> Action64)> */
		func() bool {
			position811, tokenIndex811, depth811 := position, tokenIndex, depth
			{
				position812 := position
				depth++
				{
					position813 := position
					depth++
					if buffer[position] != rune('<') {
						goto l811
					}
					position++
					if buffer[position] != rune('=') {
						goto l811
					}
					position++
					depth--
					add(rulePegText, position813)
				}
				if !_rules[ruleAction64]() {
					goto l811
				}
				depth--
				add(ruleLessOrEqual, position812)
			}
			return true
		l811:
			position, tokenIndex, depth = position811, tokenIndex811, depth811
			return false
		},
		/* 83 Greater <- <(<'>'> Action65)> */
		func() bool {
			position814, tokenIndex814, depth814 := position, tokenIndex, depth
			{
				position815 := position
				depth++
				{
					position816 := position
					depth++
					if buffer[position] != rune('>') {
						goto l814
					}
					position++
					depth--
					add(rulePegText, position816)
				}
				if !_rules[ruleAction65]() {
					goto l814
				}
				depth--
				add(ruleGreater, position815)
			}
			return true
		l814:
			position, tokenIndex, depth = position814, tokenIndex814, depth814
			return false
		},
		/* 84 GreaterOrEqual <- <(<('>' '=')> Action66)> */
		func() bool {
			position817, tokenIndex817, depth817 := position, tokenIndex, depth
			{
				position818 := position
				depth++
				{
					position819 := position
					depth++
					if buffer[position] != rune('>') {
						goto l817
					}
					position++
					if buffer[position] != rune('=') {
						goto l817
					}
					position++
					depth--
					add(rulePegText, position819)
				}
				if !_rules[ruleAction66]() {
					goto l817
				}
				depth--
				add(ruleGreaterOrEqual, position818)
			}
			return true
		l817:
			position, tokenIndex, depth = position817, tokenIndex817, depth817
			return false
		},
		/* 85 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action67)> */
		func() bool {
			position820, tokenIndex820, depth820 := position, tokenIndex, depth
			{
				position821 := position
				depth++
				{
					position822 := position
					depth++
					{
						position823, tokenIndex823, depth823 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l824
						}
						position++
						if buffer[position] != rune('=') {
							goto l824
						}
						position++
						goto l823
					l824:
						position, tokenIndex, depth = position823, tokenIndex823, depth823
						if buffer[position] != rune('<') {
							goto l820
						}
						position++
						if buffer[position] != rune('>') {
							goto l820
						}
						position++
					}
				l823:
					depth--
					add(rulePegText, position822)
				}
				if !_rules[ruleAction67]() {
					goto l820
				}
				depth--
				add(ruleNotEqual, position821)
			}
			return true
		l820:
			position, tokenIndex, depth = position820, tokenIndex820, depth820
			return false
		},
		/* 86 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action68)> */
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
						if buffer[position] != rune('i') {
							goto l829
						}
						position++
						goto l828
					l829:
						position, tokenIndex, depth = position828, tokenIndex828, depth828
						if buffer[position] != rune('I') {
							goto l825
						}
						position++
					}
				l828:
					{
						position830, tokenIndex830, depth830 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l831
						}
						position++
						goto l830
					l831:
						position, tokenIndex, depth = position830, tokenIndex830, depth830
						if buffer[position] != rune('S') {
							goto l825
						}
						position++
					}
				l830:
					depth--
					add(rulePegText, position827)
				}
				if !_rules[ruleAction68]() {
					goto l825
				}
				depth--
				add(ruleIs, position826)
			}
			return true
		l825:
			position, tokenIndex, depth = position825, tokenIndex825, depth825
			return false
		},
		/* 87 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action69)> */
		func() bool {
			position832, tokenIndex832, depth832 := position, tokenIndex, depth
			{
				position833 := position
				depth++
				{
					position834 := position
					depth++
					{
						position835, tokenIndex835, depth835 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l836
						}
						position++
						goto l835
					l836:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
						if buffer[position] != rune('I') {
							goto l832
						}
						position++
					}
				l835:
					{
						position837, tokenIndex837, depth837 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l838
						}
						position++
						goto l837
					l838:
						position, tokenIndex, depth = position837, tokenIndex837, depth837
						if buffer[position] != rune('S') {
							goto l832
						}
						position++
					}
				l837:
					if !_rules[rulesp]() {
						goto l832
					}
					{
						position839, tokenIndex839, depth839 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l840
						}
						position++
						goto l839
					l840:
						position, tokenIndex, depth = position839, tokenIndex839, depth839
						if buffer[position] != rune('N') {
							goto l832
						}
						position++
					}
				l839:
					{
						position841, tokenIndex841, depth841 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l842
						}
						position++
						goto l841
					l842:
						position, tokenIndex, depth = position841, tokenIndex841, depth841
						if buffer[position] != rune('O') {
							goto l832
						}
						position++
					}
				l841:
					{
						position843, tokenIndex843, depth843 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l844
						}
						position++
						goto l843
					l844:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						if buffer[position] != rune('T') {
							goto l832
						}
						position++
					}
				l843:
					depth--
					add(rulePegText, position834)
				}
				if !_rules[ruleAction69]() {
					goto l832
				}
				depth--
				add(ruleIsNot, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 88 Plus <- <(<'+'> Action70)> */
		func() bool {
			position845, tokenIndex845, depth845 := position, tokenIndex, depth
			{
				position846 := position
				depth++
				{
					position847 := position
					depth++
					if buffer[position] != rune('+') {
						goto l845
					}
					position++
					depth--
					add(rulePegText, position847)
				}
				if !_rules[ruleAction70]() {
					goto l845
				}
				depth--
				add(rulePlus, position846)
			}
			return true
		l845:
			position, tokenIndex, depth = position845, tokenIndex845, depth845
			return false
		},
		/* 89 Minus <- <(<'-'> Action71)> */
		func() bool {
			position848, tokenIndex848, depth848 := position, tokenIndex, depth
			{
				position849 := position
				depth++
				{
					position850 := position
					depth++
					if buffer[position] != rune('-') {
						goto l848
					}
					position++
					depth--
					add(rulePegText, position850)
				}
				if !_rules[ruleAction71]() {
					goto l848
				}
				depth--
				add(ruleMinus, position849)
			}
			return true
		l848:
			position, tokenIndex, depth = position848, tokenIndex848, depth848
			return false
		},
		/* 90 Multiply <- <(<'*'> Action72)> */
		func() bool {
			position851, tokenIndex851, depth851 := position, tokenIndex, depth
			{
				position852 := position
				depth++
				{
					position853 := position
					depth++
					if buffer[position] != rune('*') {
						goto l851
					}
					position++
					depth--
					add(rulePegText, position853)
				}
				if !_rules[ruleAction72]() {
					goto l851
				}
				depth--
				add(ruleMultiply, position852)
			}
			return true
		l851:
			position, tokenIndex, depth = position851, tokenIndex851, depth851
			return false
		},
		/* 91 Divide <- <(<'/'> Action73)> */
		func() bool {
			position854, tokenIndex854, depth854 := position, tokenIndex, depth
			{
				position855 := position
				depth++
				{
					position856 := position
					depth++
					if buffer[position] != rune('/') {
						goto l854
					}
					position++
					depth--
					add(rulePegText, position856)
				}
				if !_rules[ruleAction73]() {
					goto l854
				}
				depth--
				add(ruleDivide, position855)
			}
			return true
		l854:
			position, tokenIndex, depth = position854, tokenIndex854, depth854
			return false
		},
		/* 92 Modulo <- <(<'%'> Action74)> */
		func() bool {
			position857, tokenIndex857, depth857 := position, tokenIndex, depth
			{
				position858 := position
				depth++
				{
					position859 := position
					depth++
					if buffer[position] != rune('%') {
						goto l857
					}
					position++
					depth--
					add(rulePegText, position859)
				}
				if !_rules[ruleAction74]() {
					goto l857
				}
				depth--
				add(ruleModulo, position858)
			}
			return true
		l857:
			position, tokenIndex, depth = position857, tokenIndex857, depth857
			return false
		},
		/* 93 Identifier <- <(<ident> Action75)> */
		func() bool {
			position860, tokenIndex860, depth860 := position, tokenIndex, depth
			{
				position861 := position
				depth++
				{
					position862 := position
					depth++
					if !_rules[ruleident]() {
						goto l860
					}
					depth--
					add(rulePegText, position862)
				}
				if !_rules[ruleAction75]() {
					goto l860
				}
				depth--
				add(ruleIdentifier, position861)
			}
			return true
		l860:
			position, tokenIndex, depth = position860, tokenIndex860, depth860
			return false
		},
		/* 94 TargetIdentifier <- <(<jsonPath> Action76)> */
		func() bool {
			position863, tokenIndex863, depth863 := position, tokenIndex, depth
			{
				position864 := position
				depth++
				{
					position865 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l863
					}
					depth--
					add(rulePegText, position865)
				}
				if !_rules[ruleAction76]() {
					goto l863
				}
				depth--
				add(ruleTargetIdentifier, position864)
			}
			return true
		l863:
			position, tokenIndex, depth = position863, tokenIndex863, depth863
			return false
		},
		/* 95 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position866, tokenIndex866, depth866 := position, tokenIndex, depth
			{
				position867 := position
				depth++
				{
					position868, tokenIndex868, depth868 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l869
					}
					position++
					goto l868
				l869:
					position, tokenIndex, depth = position868, tokenIndex868, depth868
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l866
					}
					position++
				}
			l868:
			l870:
				{
					position871, tokenIndex871, depth871 := position, tokenIndex, depth
					{
						position872, tokenIndex872, depth872 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l873
						}
						position++
						goto l872
					l873:
						position, tokenIndex, depth = position872, tokenIndex872, depth872
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l874
						}
						position++
						goto l872
					l874:
						position, tokenIndex, depth = position872, tokenIndex872, depth872
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l875
						}
						position++
						goto l872
					l875:
						position, tokenIndex, depth = position872, tokenIndex872, depth872
						if buffer[position] != rune('_') {
							goto l871
						}
						position++
					}
				l872:
					goto l870
				l871:
					position, tokenIndex, depth = position871, tokenIndex871, depth871
				}
				depth--
				add(ruleident, position867)
			}
			return true
		l866:
			position, tokenIndex, depth = position866, tokenIndex866, depth866
			return false
		},
		/* 96 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position876, tokenIndex876, depth876 := position, tokenIndex, depth
			{
				position877 := position
				depth++
				{
					position878, tokenIndex878, depth878 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l879
					}
					position++
					goto l878
				l879:
					position, tokenIndex, depth = position878, tokenIndex878, depth878
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l876
					}
					position++
				}
			l878:
			l880:
				{
					position881, tokenIndex881, depth881 := position, tokenIndex, depth
					{
						position882, tokenIndex882, depth882 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l883
						}
						position++
						goto l882
					l883:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l884
						}
						position++
						goto l882
					l884:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l885
						}
						position++
						goto l882
					l885:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if buffer[position] != rune('_') {
							goto l886
						}
						position++
						goto l882
					l886:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if buffer[position] != rune('.') {
							goto l887
						}
						position++
						goto l882
					l887:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if buffer[position] != rune('[') {
							goto l888
						}
						position++
						goto l882
					l888:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if buffer[position] != rune(']') {
							goto l889
						}
						position++
						goto l882
					l889:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if buffer[position] != rune('"') {
							goto l881
						}
						position++
					}
				l882:
					goto l880
				l881:
					position, tokenIndex, depth = position881, tokenIndex881, depth881
				}
				depth--
				add(rulejsonPath, position877)
			}
			return true
		l876:
			position, tokenIndex, depth = position876, tokenIndex876, depth876
			return false
		},
		/* 97 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position891 := position
				depth++
			l892:
				{
					position893, tokenIndex893, depth893 := position, tokenIndex, depth
					{
						position894, tokenIndex894, depth894 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l895
						}
						position++
						goto l894
					l895:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
						if buffer[position] != rune('\t') {
							goto l896
						}
						position++
						goto l894
					l896:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
						if buffer[position] != rune('\n') {
							goto l893
						}
						position++
					}
				l894:
					goto l892
				l893:
					position, tokenIndex, depth = position893, tokenIndex893, depth893
				}
				depth--
				add(rulesp, position891)
			}
			return true
		},
		/* 99 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 100 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 101 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 102 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 103 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 104 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 105 Action6 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 106 Action7 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 107 Action8 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 108 Action9 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		nil,
		/* 110 Action10 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 111 Action11 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 112 Action12 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 113 Action13 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 114 Action14 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 115 Action15 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 116 Action16 <- <{
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
		/* 117 Action17 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 118 Action18 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 119 Action19 <- <{
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
		/* 120 Action20 <- <{
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
		/* 121 Action21 <- <{
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
		/* 122 Action22 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 123 Action23 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 124 Action24 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 125 Action25 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 126 Action26 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 127 Action27 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 128 Action28 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 129 Action29 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 130 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 131 Action31 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 132 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 133 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 134 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 135 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 136 Action36 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 137 Action37 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 138 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 139 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 140 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 141 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 142 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 143 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 144 Action44 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 145 Action45 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 146 Action46 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 147 Action47 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 148 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 149 Action49 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 150 Action50 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 151 Action51 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 152 Action52 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 153 Action53 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 154 Action54 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 155 Action55 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 156 Action56 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 157 Action57 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 158 Action58 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 159 Action59 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 160 Action60 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 161 Action61 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 162 Action62 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 163 Action63 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 164 Action64 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 165 Action65 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 166 Action66 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 167 Action67 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 168 Action68 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 169 Action69 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 170 Action70 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 171 Action71 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 172 Action72 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 173 Action73 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 174 Action74 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 175 Action75 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 176 Action76 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
	}
	p.rules = _rules
}

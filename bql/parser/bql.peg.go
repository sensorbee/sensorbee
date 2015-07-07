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
	rules  [173]func() bool
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

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction44:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction45:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction46:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction48:

			p.PushComponent(begin, end, Istream)

		case ruleAction49:

			p.PushComponent(begin, end, Dstream)

		case ruleAction50:

			p.PushComponent(begin, end, Rstream)

		case ruleAction51:

			p.PushComponent(begin, end, Tuples)

		case ruleAction52:

			p.PushComponent(begin, end, Seconds)

		case ruleAction53:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction54:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction55:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction56:

			p.PushComponent(begin, end, Yes)

		case ruleAction57:

			p.PushComponent(begin, end, No)

		case ruleAction58:

			p.PushComponent(begin, end, Or)

		case ruleAction59:

			p.PushComponent(begin, end, And)

		case ruleAction60:

			p.PushComponent(begin, end, Equal)

		case ruleAction61:

			p.PushComponent(begin, end, Less)

		case ruleAction62:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction63:

			p.PushComponent(begin, end, Greater)

		case ruleAction64:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction65:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction66:

			p.PushComponent(begin, end, Is)

		case ruleAction67:

			p.PushComponent(begin, end, IsNot)

		case ruleAction68:

			p.PushComponent(begin, end, Plus)

		case ruleAction69:

			p.PushComponent(begin, end, Minus)

		case ruleAction70:

			p.PushComponent(begin, end, Multiply)

		case ruleAction71:

			p.PushComponent(begin, end, Divide)

		case ruleAction72:

			p.PushComponent(begin, end, Modulo)

		case ruleAction73:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction74:

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
		/* 40 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action30)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				{
					position499 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
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
						if !_rules[rulecomparisonExpr]() {
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
		/* 41 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action31)> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				{
					position504 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l502
					}
					if !_rules[rulesp]() {
						goto l502
					}
					{
						position505, tokenIndex505, depth505 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l505
						}
						if !_rules[rulesp]() {
							goto l505
						}
						if !_rules[ruleisExpr]() {
							goto l505
						}
						goto l506
					l505:
						position, tokenIndex, depth = position505, tokenIndex505, depth505
					}
				l506:
					depth--
					add(rulePegText, position504)
				}
				if !_rules[ruleAction31]() {
					goto l502
				}
				depth--
				add(rulecomparisonExpr, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 42 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action32)> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				{
					position509 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l507
					}
					if !_rules[rulesp]() {
						goto l507
					}
					{
						position510, tokenIndex510, depth510 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l510
						}
						if !_rules[rulesp]() {
							goto l510
						}
						if !_rules[ruleNullLiteral]() {
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
				add(ruleisExpr, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 43 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action33)> */
		func() bool {
			position512, tokenIndex512, depth512 := position, tokenIndex, depth
			{
				position513 := position
				depth++
				{
					position514 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l512
					}
					if !_rules[rulesp]() {
						goto l512
					}
					{
						position515, tokenIndex515, depth515 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l515
						}
						if !_rules[rulesp]() {
							goto l515
						}
						if !_rules[ruleproductExpr]() {
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
				add(ruletermExpr, position513)
			}
			return true
		l512:
			position, tokenIndex, depth = position512, tokenIndex512, depth512
			return false
		},
		/* 44 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action34)> */
		func() bool {
			position517, tokenIndex517, depth517 := position, tokenIndex, depth
			{
				position518 := position
				depth++
				{
					position519 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l517
					}
					if !_rules[rulesp]() {
						goto l517
					}
					{
						position520, tokenIndex520, depth520 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l520
						}
						if !_rules[rulesp]() {
							goto l520
						}
						if !_rules[rulebaseExpr]() {
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
				add(ruleproductExpr, position518)
			}
			return true
		l517:
			position, tokenIndex, depth = position517, tokenIndex517, depth517
			return false
		},
		/* 45 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position522, tokenIndex522, depth522 := position, tokenIndex, depth
			{
				position523 := position
				depth++
				{
					position524, tokenIndex524, depth524 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l525
					}
					position++
					if !_rules[rulesp]() {
						goto l525
					}
					if !_rules[ruleExpression]() {
						goto l525
					}
					if !_rules[rulesp]() {
						goto l525
					}
					if buffer[position] != rune(')') {
						goto l525
					}
					position++
					goto l524
				l525:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleBooleanLiteral]() {
						goto l526
					}
					goto l524
				l526:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleNullLiteral]() {
						goto l527
					}
					goto l524
				l527:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleFuncApp]() {
						goto l528
					}
					goto l524
				l528:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleRowMeta]() {
						goto l529
					}
					goto l524
				l529:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleRowValue]() {
						goto l530
					}
					goto l524
				l530:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if !_rules[ruleLiteral]() {
						goto l522
					}
				}
			l524:
				depth--
				add(rulebaseExpr, position523)
			}
			return true
		l522:
			position, tokenIndex, depth = position522, tokenIndex522, depth522
			return false
		},
		/* 46 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action35)> */
		func() bool {
			position531, tokenIndex531, depth531 := position, tokenIndex, depth
			{
				position532 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l531
				}
				if !_rules[rulesp]() {
					goto l531
				}
				if buffer[position] != rune('(') {
					goto l531
				}
				position++
				if !_rules[rulesp]() {
					goto l531
				}
				if !_rules[ruleFuncParams]() {
					goto l531
				}
				if !_rules[rulesp]() {
					goto l531
				}
				if buffer[position] != rune(')') {
					goto l531
				}
				position++
				if !_rules[ruleAction35]() {
					goto l531
				}
				depth--
				add(ruleFuncApp, position532)
			}
			return true
		l531:
			position, tokenIndex, depth = position531, tokenIndex531, depth531
			return false
		},
		/* 47 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action36)> */
		func() bool {
			position533, tokenIndex533, depth533 := position, tokenIndex, depth
			{
				position534 := position
				depth++
				{
					position535 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l533
					}
					if !_rules[rulesp]() {
						goto l533
					}
				l536:
					{
						position537, tokenIndex537, depth537 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l537
						}
						position++
						if !_rules[rulesp]() {
							goto l537
						}
						if !_rules[ruleExpression]() {
							goto l537
						}
						goto l536
					l537:
						position, tokenIndex, depth = position537, tokenIndex537, depth537
					}
					depth--
					add(rulePegText, position535)
				}
				if !_rules[ruleAction36]() {
					goto l533
				}
				depth--
				add(ruleFuncParams, position534)
			}
			return true
		l533:
			position, tokenIndex, depth = position533, tokenIndex533, depth533
			return false
		},
		/* 48 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					position540, tokenIndex540, depth540 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l541
					}
					goto l540
				l541:
					position, tokenIndex, depth = position540, tokenIndex540, depth540
					if !_rules[ruleNumericLiteral]() {
						goto l542
					}
					goto l540
				l542:
					position, tokenIndex, depth = position540, tokenIndex540, depth540
					if !_rules[ruleStringLiteral]() {
						goto l538
					}
				}
			l540:
				depth--
				add(ruleLiteral, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 49 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position543, tokenIndex543, depth543 := position, tokenIndex, depth
			{
				position544 := position
				depth++
				{
					position545, tokenIndex545, depth545 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l546
					}
					goto l545
				l546:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleNotEqual]() {
						goto l547
					}
					goto l545
				l547:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleLessOrEqual]() {
						goto l548
					}
					goto l545
				l548:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleLess]() {
						goto l549
					}
					goto l545
				l549:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleGreaterOrEqual]() {
						goto l550
					}
					goto l545
				l550:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleGreater]() {
						goto l551
					}
					goto l545
				l551:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !_rules[ruleNotEqual]() {
						goto l543
					}
				}
			l545:
				depth--
				add(ruleComparisonOp, position544)
			}
			return true
		l543:
			position, tokenIndex, depth = position543, tokenIndex543, depth543
			return false
		},
		/* 50 IsOp <- <(IsNot / Is)> */
		func() bool {
			position552, tokenIndex552, depth552 := position, tokenIndex, depth
			{
				position553 := position
				depth++
				{
					position554, tokenIndex554, depth554 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l555
					}
					goto l554
				l555:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
					if !_rules[ruleIs]() {
						goto l552
					}
				}
			l554:
				depth--
				add(ruleIsOp, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 51 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position556, tokenIndex556, depth556 := position, tokenIndex, depth
			{
				position557 := position
				depth++
				{
					position558, tokenIndex558, depth558 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l559
					}
					goto l558
				l559:
					position, tokenIndex, depth = position558, tokenIndex558, depth558
					if !_rules[ruleMinus]() {
						goto l556
					}
				}
			l558:
				depth--
				add(rulePlusMinusOp, position557)
			}
			return true
		l556:
			position, tokenIndex, depth = position556, tokenIndex556, depth556
			return false
		},
		/* 52 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position560, tokenIndex560, depth560 := position, tokenIndex, depth
			{
				position561 := position
				depth++
				{
					position562, tokenIndex562, depth562 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l563
					}
					goto l562
				l563:
					position, tokenIndex, depth = position562, tokenIndex562, depth562
					if !_rules[ruleDivide]() {
						goto l564
					}
					goto l562
				l564:
					position, tokenIndex, depth = position562, tokenIndex562, depth562
					if !_rules[ruleModulo]() {
						goto l560
					}
				}
			l562:
				depth--
				add(ruleMultDivOp, position561)
			}
			return true
		l560:
			position, tokenIndex, depth = position560, tokenIndex560, depth560
			return false
		},
		/* 53 Stream <- <(<ident> Action37)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567 := position
					depth++
					if !_rules[ruleident]() {
						goto l565
					}
					depth--
					add(rulePegText, position567)
				}
				if !_rules[ruleAction37]() {
					goto l565
				}
				depth--
				add(ruleStream, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 54 RowMeta <- <RowTimestamp> */
		func() bool {
			position568, tokenIndex568, depth568 := position, tokenIndex, depth
			{
				position569 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l568
				}
				depth--
				add(ruleRowMeta, position569)
			}
			return true
		l568:
			position, tokenIndex, depth = position568, tokenIndex568, depth568
			return false
		},
		/* 55 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action38)> */
		func() bool {
			position570, tokenIndex570, depth570 := position, tokenIndex, depth
			{
				position571 := position
				depth++
				{
					position572 := position
					depth++
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l573
						}
						if buffer[position] != rune(':') {
							goto l573
						}
						position++
						goto l574
					l573:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
					}
				l574:
					if buffer[position] != rune('t') {
						goto l570
					}
					position++
					if buffer[position] != rune('s') {
						goto l570
					}
					position++
					if buffer[position] != rune('(') {
						goto l570
					}
					position++
					if buffer[position] != rune(')') {
						goto l570
					}
					position++
					depth--
					add(rulePegText, position572)
				}
				if !_rules[ruleAction38]() {
					goto l570
				}
				depth--
				add(ruleRowTimestamp, position571)
			}
			return true
		l570:
			position, tokenIndex, depth = position570, tokenIndex570, depth570
			return false
		},
		/* 56 RowValue <- <(<((ident ':')? jsonPath)> Action39)> */
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
					if !_rules[rulejsonPath]() {
						goto l575
					}
					depth--
					add(rulePegText, position577)
				}
				if !_rules[ruleAction39]() {
					goto l575
				}
				depth--
				add(ruleRowValue, position576)
			}
			return true
		l575:
			position, tokenIndex, depth = position575, tokenIndex575, depth575
			return false
		},
		/* 57 NumericLiteral <- <(<('-'? [0-9]+)> Action40)> */
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
						if buffer[position] != rune('-') {
							goto l583
						}
						position++
						goto l584
					l583:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
					}
				l584:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l580
					}
					position++
				l585:
					{
						position586, tokenIndex586, depth586 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l586
						}
						position++
						goto l585
					l586:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
					}
					depth--
					add(rulePegText, position582)
				}
				if !_rules[ruleAction40]() {
					goto l580
				}
				depth--
				add(ruleNumericLiteral, position581)
			}
			return true
		l580:
			position, tokenIndex, depth = position580, tokenIndex580, depth580
			return false
		},
		/* 58 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action41)> */
		func() bool {
			position587, tokenIndex587, depth587 := position, tokenIndex, depth
			{
				position588 := position
				depth++
				{
					position589 := position
					depth++
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l590
						}
						position++
						goto l591
					l590:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
					}
				l591:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l587
					}
					position++
				l592:
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
					}
					if buffer[position] != rune('.') {
						goto l587
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l587
					}
					position++
				l594:
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
					}
					depth--
					add(rulePegText, position589)
				}
				if !_rules[ruleAction41]() {
					goto l587
				}
				depth--
				add(ruleFloatLiteral, position588)
			}
			return true
		l587:
			position, tokenIndex, depth = position587, tokenIndex587, depth587
			return false
		},
		/* 59 Function <- <(<ident> Action42)> */
		func() bool {
			position596, tokenIndex596, depth596 := position, tokenIndex, depth
			{
				position597 := position
				depth++
				{
					position598 := position
					depth++
					if !_rules[ruleident]() {
						goto l596
					}
					depth--
					add(rulePegText, position598)
				}
				if !_rules[ruleAction42]() {
					goto l596
				}
				depth--
				add(ruleFunction, position597)
			}
			return true
		l596:
			position, tokenIndex, depth = position596, tokenIndex596, depth596
			return false
		},
		/* 60 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action43)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				{
					position601 := position
					depth++
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l603
						}
						position++
						goto l602
					l603:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if buffer[position] != rune('N') {
							goto l599
						}
						position++
					}
				l602:
					{
						position604, tokenIndex604, depth604 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l605
						}
						position++
						goto l604
					l605:
						position, tokenIndex, depth = position604, tokenIndex604, depth604
						if buffer[position] != rune('U') {
							goto l599
						}
						position++
					}
				l604:
					{
						position606, tokenIndex606, depth606 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l607
						}
						position++
						goto l606
					l607:
						position, tokenIndex, depth = position606, tokenIndex606, depth606
						if buffer[position] != rune('L') {
							goto l599
						}
						position++
					}
				l606:
					{
						position608, tokenIndex608, depth608 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l609
						}
						position++
						goto l608
					l609:
						position, tokenIndex, depth = position608, tokenIndex608, depth608
						if buffer[position] != rune('L') {
							goto l599
						}
						position++
					}
				l608:
					depth--
					add(rulePegText, position601)
				}
				if !_rules[ruleAction43]() {
					goto l599
				}
				depth--
				add(ruleNullLiteral, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 61 BooleanLiteral <- <(TRUE / FALSE)> */
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
		/* 62 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action44)> */
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
				if !_rules[ruleAction44]() {
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
		/* 63 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action45)> */
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
				if !_rules[ruleAction45]() {
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
		/* 64 Wildcard <- <(<'*'> Action46)> */
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
				if !_rules[ruleAction46]() {
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
		/* 65 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action47)> */
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
				if !_rules[ruleAction47]() {
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
		/* 66 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action48)> */
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
				if !_rules[ruleAction48]() {
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
		/* 67 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action49)> */
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
				if !_rules[ruleAction49]() {
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
		/* 68 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action50)> */
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
				if !_rules[ruleAction50]() {
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
		/* 69 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action51)> */
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
				if !_rules[ruleAction51]() {
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
		/* 70 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action52)> */
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
				if !_rules[ruleAction52]() {
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
		/* 71 StreamIdentifier <- <(<ident> Action53)> */
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
				if !_rules[ruleAction53]() {
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
		/* 72 SourceSinkType <- <(<ident> Action54)> */
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
				if !_rules[ruleAction54]() {
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
		/* 73 SourceSinkParamKey <- <(<ident> Action55)> */
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
				if !_rules[ruleAction55]() {
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
		/* 74 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action56)> */
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
				if !_rules[ruleAction56]() {
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
		/* 75 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action57)> */
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
				if !_rules[ruleAction57]() {
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
		/* 76 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action58)> */
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
				if !_rules[ruleAction58]() {
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
		/* 77 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action59)> */
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
				if !_rules[ruleAction59]() {
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
		/* 78 Equal <- <(<'='> Action60)> */
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
				if !_rules[ruleAction60]() {
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
		/* 79 Less <- <(<'<'> Action61)> */
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
				if !_rules[ruleAction61]() {
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
		/* 80 LessOrEqual <- <(<('<' '=')> Action62)> */
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
				if !_rules[ruleAction62]() {
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
		/* 81 Greater <- <(<'>'> Action63)> */
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
				if !_rules[ruleAction63]() {
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
		/* 82 GreaterOrEqual <- <(<('>' '=')> Action64)> */
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
				if !_rules[ruleAction64]() {
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
		/* 83 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action65)> */
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
				if !_rules[ruleAction65]() {
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
		/* 84 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action66)> */
		func() bool {
			position811, tokenIndex811, depth811 := position, tokenIndex, depth
			{
				position812 := position
				depth++
				{
					position813 := position
					depth++
					{
						position814, tokenIndex814, depth814 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l815
						}
						position++
						goto l814
					l815:
						position, tokenIndex, depth = position814, tokenIndex814, depth814
						if buffer[position] != rune('I') {
							goto l811
						}
						position++
					}
				l814:
					{
						position816, tokenIndex816, depth816 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l817
						}
						position++
						goto l816
					l817:
						position, tokenIndex, depth = position816, tokenIndex816, depth816
						if buffer[position] != rune('S') {
							goto l811
						}
						position++
					}
				l816:
					depth--
					add(rulePegText, position813)
				}
				if !_rules[ruleAction66]() {
					goto l811
				}
				depth--
				add(ruleIs, position812)
			}
			return true
		l811:
			position, tokenIndex, depth = position811, tokenIndex811, depth811
			return false
		},
		/* 85 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action67)> */
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
						if buffer[position] != rune('i') {
							goto l822
						}
						position++
						goto l821
					l822:
						position, tokenIndex, depth = position821, tokenIndex821, depth821
						if buffer[position] != rune('I') {
							goto l818
						}
						position++
					}
				l821:
					{
						position823, tokenIndex823, depth823 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l824
						}
						position++
						goto l823
					l824:
						position, tokenIndex, depth = position823, tokenIndex823, depth823
						if buffer[position] != rune('S') {
							goto l818
						}
						position++
					}
				l823:
					if !_rules[rulesp]() {
						goto l818
					}
					{
						position825, tokenIndex825, depth825 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l826
						}
						position++
						goto l825
					l826:
						position, tokenIndex, depth = position825, tokenIndex825, depth825
						if buffer[position] != rune('N') {
							goto l818
						}
						position++
					}
				l825:
					{
						position827, tokenIndex827, depth827 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l828
						}
						position++
						goto l827
					l828:
						position, tokenIndex, depth = position827, tokenIndex827, depth827
						if buffer[position] != rune('O') {
							goto l818
						}
						position++
					}
				l827:
					{
						position829, tokenIndex829, depth829 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l830
						}
						position++
						goto l829
					l830:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						if buffer[position] != rune('T') {
							goto l818
						}
						position++
					}
				l829:
					depth--
					add(rulePegText, position820)
				}
				if !_rules[ruleAction67]() {
					goto l818
				}
				depth--
				add(ruleIsNot, position819)
			}
			return true
		l818:
			position, tokenIndex, depth = position818, tokenIndex818, depth818
			return false
		},
		/* 86 Plus <- <(<'+'> Action68)> */
		func() bool {
			position831, tokenIndex831, depth831 := position, tokenIndex, depth
			{
				position832 := position
				depth++
				{
					position833 := position
					depth++
					if buffer[position] != rune('+') {
						goto l831
					}
					position++
					depth--
					add(rulePegText, position833)
				}
				if !_rules[ruleAction68]() {
					goto l831
				}
				depth--
				add(rulePlus, position832)
			}
			return true
		l831:
			position, tokenIndex, depth = position831, tokenIndex831, depth831
			return false
		},
		/* 87 Minus <- <(<'-'> Action69)> */
		func() bool {
			position834, tokenIndex834, depth834 := position, tokenIndex, depth
			{
				position835 := position
				depth++
				{
					position836 := position
					depth++
					if buffer[position] != rune('-') {
						goto l834
					}
					position++
					depth--
					add(rulePegText, position836)
				}
				if !_rules[ruleAction69]() {
					goto l834
				}
				depth--
				add(ruleMinus, position835)
			}
			return true
		l834:
			position, tokenIndex, depth = position834, tokenIndex834, depth834
			return false
		},
		/* 88 Multiply <- <(<'*'> Action70)> */
		func() bool {
			position837, tokenIndex837, depth837 := position, tokenIndex, depth
			{
				position838 := position
				depth++
				{
					position839 := position
					depth++
					if buffer[position] != rune('*') {
						goto l837
					}
					position++
					depth--
					add(rulePegText, position839)
				}
				if !_rules[ruleAction70]() {
					goto l837
				}
				depth--
				add(ruleMultiply, position838)
			}
			return true
		l837:
			position, tokenIndex, depth = position837, tokenIndex837, depth837
			return false
		},
		/* 89 Divide <- <(<'/'> Action71)> */
		func() bool {
			position840, tokenIndex840, depth840 := position, tokenIndex, depth
			{
				position841 := position
				depth++
				{
					position842 := position
					depth++
					if buffer[position] != rune('/') {
						goto l840
					}
					position++
					depth--
					add(rulePegText, position842)
				}
				if !_rules[ruleAction71]() {
					goto l840
				}
				depth--
				add(ruleDivide, position841)
			}
			return true
		l840:
			position, tokenIndex, depth = position840, tokenIndex840, depth840
			return false
		},
		/* 90 Modulo <- <(<'%'> Action72)> */
		func() bool {
			position843, tokenIndex843, depth843 := position, tokenIndex, depth
			{
				position844 := position
				depth++
				{
					position845 := position
					depth++
					if buffer[position] != rune('%') {
						goto l843
					}
					position++
					depth--
					add(rulePegText, position845)
				}
				if !_rules[ruleAction72]() {
					goto l843
				}
				depth--
				add(ruleModulo, position844)
			}
			return true
		l843:
			position, tokenIndex, depth = position843, tokenIndex843, depth843
			return false
		},
		/* 91 Identifier <- <(<ident> Action73)> */
		func() bool {
			position846, tokenIndex846, depth846 := position, tokenIndex, depth
			{
				position847 := position
				depth++
				{
					position848 := position
					depth++
					if !_rules[ruleident]() {
						goto l846
					}
					depth--
					add(rulePegText, position848)
				}
				if !_rules[ruleAction73]() {
					goto l846
				}
				depth--
				add(ruleIdentifier, position847)
			}
			return true
		l846:
			position, tokenIndex, depth = position846, tokenIndex846, depth846
			return false
		},
		/* 92 TargetIdentifier <- <(<jsonPath> Action74)> */
		func() bool {
			position849, tokenIndex849, depth849 := position, tokenIndex, depth
			{
				position850 := position
				depth++
				{
					position851 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l849
					}
					depth--
					add(rulePegText, position851)
				}
				if !_rules[ruleAction74]() {
					goto l849
				}
				depth--
				add(ruleTargetIdentifier, position850)
			}
			return true
		l849:
			position, tokenIndex, depth = position849, tokenIndex849, depth849
			return false
		},
		/* 93 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position852, tokenIndex852, depth852 := position, tokenIndex, depth
			{
				position853 := position
				depth++
				{
					position854, tokenIndex854, depth854 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l855
					}
					position++
					goto l854
				l855:
					position, tokenIndex, depth = position854, tokenIndex854, depth854
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l852
					}
					position++
				}
			l854:
			l856:
				{
					position857, tokenIndex857, depth857 := position, tokenIndex, depth
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l859
						}
						position++
						goto l858
					l859:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l860
						}
						position++
						goto l858
					l860:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l861
						}
						position++
						goto l858
					l861:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if buffer[position] != rune('_') {
							goto l857
						}
						position++
					}
				l858:
					goto l856
				l857:
					position, tokenIndex, depth = position857, tokenIndex857, depth857
				}
				depth--
				add(ruleident, position853)
			}
			return true
		l852:
			position, tokenIndex, depth = position852, tokenIndex852, depth852
			return false
		},
		/* 94 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position862, tokenIndex862, depth862 := position, tokenIndex, depth
			{
				position863 := position
				depth++
				{
					position864, tokenIndex864, depth864 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l865
					}
					position++
					goto l864
				l865:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l862
					}
					position++
				}
			l864:
			l866:
				{
					position867, tokenIndex867, depth867 := position, tokenIndex, depth
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
							goto l870
						}
						position++
						goto l868
					l870:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l871
						}
						position++
						goto l868
					l871:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('_') {
							goto l872
						}
						position++
						goto l868
					l872:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('.') {
							goto l873
						}
						position++
						goto l868
					l873:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('[') {
							goto l874
						}
						position++
						goto l868
					l874:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune(']') {
							goto l875
						}
						position++
						goto l868
					l875:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('"') {
							goto l867
						}
						position++
					}
				l868:
					goto l866
				l867:
					position, tokenIndex, depth = position867, tokenIndex867, depth867
				}
				depth--
				add(rulejsonPath, position863)
			}
			return true
		l862:
			position, tokenIndex, depth = position862, tokenIndex862, depth862
			return false
		},
		/* 95 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position877 := position
				depth++
			l878:
				{
					position879, tokenIndex879, depth879 := position, tokenIndex, depth
					{
						position880, tokenIndex880, depth880 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l881
						}
						position++
						goto l880
					l881:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
						if buffer[position] != rune('\t') {
							goto l882
						}
						position++
						goto l880
					l882:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
						if buffer[position] != rune('\n') {
							goto l879
						}
						position++
					}
				l880:
					goto l878
				l879:
					position, tokenIndex, depth = position879, tokenIndex879, depth879
				}
				depth--
				add(rulesp, position877)
			}
			return true
		},
		/* 97 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 98 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 99 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 100 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 101 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 102 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 103 Action6 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 104 Action7 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 105 Action8 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 106 Action9 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		nil,
		/* 108 Action10 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 109 Action11 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 110 Action12 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 111 Action13 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 112 Action14 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 113 Action15 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 114 Action16 <- <{
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
		/* 115 Action17 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 116 Action18 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 117 Action19 <- <{
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
		/* 118 Action20 <- <{
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
		/* 119 Action21 <- <{
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
		/* 120 Action22 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 121 Action23 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 122 Action24 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 123 Action25 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 124 Action26 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 125 Action27 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 126 Action28 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 127 Action29 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 128 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 129 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 130 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 131 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 132 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 133 Action35 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 134 Action36 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 135 Action37 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 136 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 137 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 138 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 139 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 140 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 141 Action43 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 142 Action44 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 143 Action45 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 144 Action46 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 145 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 146 Action48 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 147 Action49 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 148 Action50 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 149 Action51 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 150 Action52 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 151 Action53 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 152 Action54 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 153 Action55 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 154 Action56 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 155 Action57 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 156 Action58 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 157 Action59 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 158 Action60 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 159 Action61 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 160 Action62 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 161 Action63 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 162 Action64 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 163 Action65 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 164 Action66 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 165 Action67 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 166 Action68 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 167 Action69 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 168 Action70 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 169 Action71 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 170 Action72 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 171 Action73 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 172 Action74 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
	}
	p.rules = _rules
}

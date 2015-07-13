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
	ruleDropSourceStmt
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
	ruleminusExpr
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
	ruleUnaryMinus
	ruleIdentifier
	ruleTargetIdentifier
	ruleident
	rulejsonPath
	rulesp
	rulecomment
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
	ruleAction10
	rulePegText
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
	ruleAction79

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
	"DropSourceStmt",
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
	"minusExpr",
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
	"UnaryMinus",
	"Identifier",
	"TargetIdentifier",
	"ident",
	"jsonPath",
	"sp",
	"comment",
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
	"Action10",
	"PegText",
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
	"Action79",

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
	rules  [184]func() bool
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

			p.AssembleDropSource()

		case ruleAction11:

			p.AssembleEmitter(begin, end)

		case ruleAction12:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction13:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction14:

			p.AssembleStreamEmitInterval()

		case ruleAction15:

			p.AssembleProjections(begin, end)

		case ruleAction16:

			p.AssembleAlias()

		case ruleAction17:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
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

			p.AssembleAliasedStreamWindow()

		case ruleAction25:

			p.AssembleStreamWindow()

		case ruleAction26:

			p.AssembleUDSFFuncApp()

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

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction33:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction34:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction38:

			p.AssembleFuncApp()

		case ruleAction39:

			p.AssembleExpressions(begin, end)

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction44:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction46:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction47:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction48:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction49:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction51:

			p.PushComponent(begin, end, Istream)

		case ruleAction52:

			p.PushComponent(begin, end, Dstream)

		case ruleAction53:

			p.PushComponent(begin, end, Rstream)

		case ruleAction54:

			p.PushComponent(begin, end, Tuples)

		case ruleAction55:

			p.PushComponent(begin, end, Seconds)

		case ruleAction56:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction57:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction58:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction59:

			p.PushComponent(begin, end, Yes)

		case ruleAction60:

			p.PushComponent(begin, end, No)

		case ruleAction61:

			p.PushComponent(begin, end, Or)

		case ruleAction62:

			p.PushComponent(begin, end, And)

		case ruleAction63:

			p.PushComponent(begin, end, Not)

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

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction78:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction79:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt / DropSourceStmt)> */
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
						goto l19
					}
					goto l9
				l19:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleDropSourceStmt]() {
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
			position20, tokenIndex20, depth20 := position, tokenIndex, depth
			{
				position21 := position
				depth++
				{
					position22, tokenIndex22, depth22 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l23
					}
					position++
					goto l22
				l23:
					position, tokenIndex, depth = position22, tokenIndex22, depth22
					if buffer[position] != rune('S') {
						goto l20
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
						goto l20
					}
					position++
				}
			l24:
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l27
					}
					position++
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					if buffer[position] != rune('L') {
						goto l20
					}
					position++
				}
			l26:
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l29
					}
					position++
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if buffer[position] != rune('E') {
						goto l20
					}
					position++
				}
			l28:
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l31
					}
					position++
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if buffer[position] != rune('C') {
						goto l20
					}
					position++
				}
			l30:
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l33
					}
					position++
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if buffer[position] != rune('T') {
						goto l20
					}
					position++
				}
			l32:
				if !_rules[rulesp]() {
					goto l20
				}
				if !_rules[ruleEmitter]() {
					goto l20
				}
				if !_rules[rulesp]() {
					goto l20
				}
				if !_rules[ruleProjections]() {
					goto l20
				}
				if !_rules[rulesp]() {
					goto l20
				}
				if !_rules[ruleWindowedFrom]() {
					goto l20
				}
				if !_rules[rulesp]() {
					goto l20
				}
				if !_rules[ruleFilter]() {
					goto l20
				}
				if !_rules[rulesp]() {
					goto l20
				}
				if !_rules[ruleGrouping]() {
					goto l20
				}
				if !_rules[rulesp]() {
					goto l20
				}
				if !_rules[ruleHaving]() {
					goto l20
				}
				if !_rules[rulesp]() {
					goto l20
				}
				if !_rules[ruleAction0]() {
					goto l20
				}
				depth--
				add(ruleSelectStmt, position21)
			}
			return true
		l20:
			position, tokenIndex, depth = position20, tokenIndex20, depth20
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
		/* 8 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action6)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('I') {
						goto l194
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('N') {
						goto l194
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('S') {
						goto l194
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
						goto l194
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('R') {
						goto l194
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('T') {
						goto l194
					}
					position++
				}
			l206:
				if !_rules[rulesp]() {
					goto l194
				}
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('I') {
						goto l194
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('N') {
						goto l194
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('T') {
						goto l194
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
						goto l194
					}
					position++
				}
			l214:
				if !_rules[rulesp]() {
					goto l194
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l194
				}
				if !_rules[rulesp]() {
					goto l194
				}
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('F') {
						goto l194
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('R') {
						goto l194
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
						goto l194
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('M') {
						goto l194
					}
					position++
				}
			l222:
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
				add(ruleInsertIntoFromStmt, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 9 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action7)> */
		func() bool {
			position224, tokenIndex224, depth224 := position, tokenIndex, depth
			{
				position225 := position
				depth++
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('P') {
						goto l224
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('A') {
						goto l224
					}
					position++
				}
			l228:
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
					if buffer[position] != rune('U') {
						goto l224
					}
					position++
				}
			l230:
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
						goto l224
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('E') {
						goto l224
					}
					position++
				}
			l234:
				if !_rules[rulesp]() {
					goto l224
				}
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
						goto l224
					}
					position++
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('O') {
						goto l224
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('U') {
						goto l224
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('R') {
						goto l224
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('C') {
						goto l224
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('E') {
						goto l224
					}
					position++
				}
			l246:
				if !_rules[rulesp]() {
					goto l224
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l224
				}
				if !_rules[ruleAction7]() {
					goto l224
				}
				depth--
				add(rulePauseSourceStmt, position225)
			}
			return true
		l224:
			position, tokenIndex, depth = position224, tokenIndex224, depth224
			return false
		},
		/* 10 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action8)> */
		func() bool {
			position248, tokenIndex248, depth248 := position, tokenIndex, depth
			{
				position249 := position
				depth++
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('R') {
						goto l248
					}
					position++
				}
			l250:
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if buffer[position] != rune('E') {
						goto l248
					}
					position++
				}
			l252:
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if buffer[position] != rune('S') {
						goto l248
					}
					position++
				}
			l254:
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l257
					}
					position++
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if buffer[position] != rune('U') {
						goto l248
					}
					position++
				}
			l256:
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if buffer[position] != rune('M') {
						goto l248
					}
					position++
				}
			l258:
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
						goto l248
					}
					position++
				}
			l260:
				if !_rules[rulesp]() {
					goto l248
				}
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('S') {
						goto l248
					}
					position++
				}
			l262:
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('O') {
						goto l248
					}
					position++
				}
			l264:
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('U') {
						goto l248
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('R') {
						goto l248
					}
					position++
				}
			l268:
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if buffer[position] != rune('C') {
						goto l248
					}
					position++
				}
			l270:
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l273
					}
					position++
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if buffer[position] != rune('E') {
						goto l248
					}
					position++
				}
			l272:
				if !_rules[rulesp]() {
					goto l248
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l248
				}
				if !_rules[ruleAction8]() {
					goto l248
				}
				depth--
				add(ruleResumeSourceStmt, position249)
			}
			return true
		l248:
			position, tokenIndex, depth = position248, tokenIndex248, depth248
			return false
		},
		/* 11 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action9)> */
		func() bool {
			position274, tokenIndex274, depth274 := position, tokenIndex, depth
			{
				position275 := position
				depth++
				{
					position276, tokenIndex276, depth276 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l277
					}
					position++
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					if buffer[position] != rune('R') {
						goto l274
					}
					position++
				}
			l276:
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
						goto l274
					}
					position++
				}
			l278:
				{
					position280, tokenIndex280, depth280 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l281
					}
					position++
					goto l280
				l281:
					position, tokenIndex, depth = position280, tokenIndex280, depth280
					if buffer[position] != rune('W') {
						goto l274
					}
					position++
				}
			l280:
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('I') {
						goto l274
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
					if buffer[position] != rune('N') {
						goto l274
					}
					position++
				}
			l284:
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l287
					}
					position++
					goto l286
				l287:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
					if buffer[position] != rune('D') {
						goto l274
					}
					position++
				}
			l286:
				if !_rules[rulesp]() {
					goto l274
				}
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l289
					}
					position++
					goto l288
				l289:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
					if buffer[position] != rune('S') {
						goto l274
					}
					position++
				}
			l288:
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l291
					}
					position++
					goto l290
				l291:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
					if buffer[position] != rune('O') {
						goto l274
					}
					position++
				}
			l290:
				{
					position292, tokenIndex292, depth292 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l293
					}
					position++
					goto l292
				l293:
					position, tokenIndex, depth = position292, tokenIndex292, depth292
					if buffer[position] != rune('U') {
						goto l274
					}
					position++
				}
			l292:
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l295
					}
					position++
					goto l294
				l295:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if buffer[position] != rune('R') {
						goto l274
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if buffer[position] != rune('C') {
						goto l274
					}
					position++
				}
			l296:
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					if buffer[position] != rune('E') {
						goto l274
					}
					position++
				}
			l298:
				if !_rules[rulesp]() {
					goto l274
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l274
				}
				if !_rules[ruleAction9]() {
					goto l274
				}
				depth--
				add(ruleRewindSourceStmt, position275)
			}
			return true
		l274:
			position, tokenIndex, depth = position274, tokenIndex274, depth274
			return false
		},
		/* 12 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position300, tokenIndex300, depth300 := position, tokenIndex, depth
			{
				position301 := position
				depth++
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l303
					}
					position++
					goto l302
				l303:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					if buffer[position] != rune('D') {
						goto l300
					}
					position++
				}
			l302:
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
					if buffer[position] != rune('R') {
						goto l300
					}
					position++
				}
			l304:
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l307
					}
					position++
					goto l306
				l307:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					if buffer[position] != rune('O') {
						goto l300
					}
					position++
				}
			l306:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					if buffer[position] != rune('P') {
						goto l300
					}
					position++
				}
			l308:
				if !_rules[rulesp]() {
					goto l300
				}
				{
					position310, tokenIndex310, depth310 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l311
					}
					position++
					goto l310
				l311:
					position, tokenIndex, depth = position310, tokenIndex310, depth310
					if buffer[position] != rune('S') {
						goto l300
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
						goto l300
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
						goto l300
					}
					position++
				}
			l314:
				{
					position316, tokenIndex316, depth316 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l317
					}
					position++
					goto l316
				l317:
					position, tokenIndex, depth = position316, tokenIndex316, depth316
					if buffer[position] != rune('R') {
						goto l300
					}
					position++
				}
			l316:
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
					if buffer[position] != rune('C') {
						goto l300
					}
					position++
				}
			l318:
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('E') {
						goto l300
					}
					position++
				}
			l320:
				if !_rules[rulesp]() {
					goto l300
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l300
				}
				if !_rules[ruleAction10]() {
					goto l300
				}
				depth--
				add(ruleDropSourceStmt, position301)
			}
			return true
		l300:
			position, tokenIndex, depth = position300, tokenIndex300, depth300
			return false
		},
		/* 13 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action11)> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l325
					}
					goto l324
				l325:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if !_rules[ruleDSTREAM]() {
						goto l326
					}
					goto l324
				l326:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if !_rules[ruleRSTREAM]() {
						goto l322
					}
				}
			l324:
				{
					position327 := position
					depth++
					{
						position328, tokenIndex328, depth328 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l328
						}
						if buffer[position] != rune('[') {
							goto l328
						}
						position++
						if !_rules[rulesp]() {
							goto l328
						}
						{
							position330, tokenIndex330, depth330 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l331
							}
							position++
							goto l330
						l331:
							position, tokenIndex, depth = position330, tokenIndex330, depth330
							if buffer[position] != rune('E') {
								goto l328
							}
							position++
						}
					l330:
						{
							position332, tokenIndex332, depth332 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l333
							}
							position++
							goto l332
						l333:
							position, tokenIndex, depth = position332, tokenIndex332, depth332
							if buffer[position] != rune('V') {
								goto l328
							}
							position++
						}
					l332:
						{
							position334, tokenIndex334, depth334 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l335
							}
							position++
							goto l334
						l335:
							position, tokenIndex, depth = position334, tokenIndex334, depth334
							if buffer[position] != rune('E') {
								goto l328
							}
							position++
						}
					l334:
						{
							position336, tokenIndex336, depth336 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l337
							}
							position++
							goto l336
						l337:
							position, tokenIndex, depth = position336, tokenIndex336, depth336
							if buffer[position] != rune('R') {
								goto l328
							}
							position++
						}
					l336:
						{
							position338, tokenIndex338, depth338 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l339
							}
							position++
							goto l338
						l339:
							position, tokenIndex, depth = position338, tokenIndex338, depth338
							if buffer[position] != rune('Y') {
								goto l328
							}
							position++
						}
					l338:
						if !_rules[rulesp]() {
							goto l328
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l328
						}
						if !_rules[rulesp]() {
							goto l328
						}
						if buffer[position] != rune(']') {
							goto l328
						}
						position++
						goto l329
					l328:
						position, tokenIndex, depth = position328, tokenIndex328, depth328
					}
				l329:
					depth--
					add(rulePegText, position327)
				}
				if !_rules[ruleAction11]() {
					goto l322
				}
				depth--
				add(ruleEmitter, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 14 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position340, tokenIndex340, depth340 := position, tokenIndex, depth
			{
				position341 := position
				depth++
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l343
					}
				l344:
					{
						position345, tokenIndex345, depth345 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l345
						}
						if buffer[position] != rune(',') {
							goto l345
						}
						position++
						if !_rules[rulesp]() {
							goto l345
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l345
						}
						goto l344
					l345:
						position, tokenIndex, depth = position345, tokenIndex345, depth345
					}
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if !_rules[ruleTimeEmitterInterval]() {
						goto l346
					}
					goto l342
				l346:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if !_rules[ruleTupleEmitterInterval]() {
						goto l340
					}
				}
			l342:
				depth--
				add(ruleEmitterIntervals, position341)
			}
			return true
		l340:
			position, tokenIndex, depth = position340, tokenIndex340, depth340
			return false
		},
		/* 15 TimeEmitterInterval <- <(<TimeInterval> Action12)> */
		func() bool {
			position347, tokenIndex347, depth347 := position, tokenIndex, depth
			{
				position348 := position
				depth++
				{
					position349 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l347
					}
					depth--
					add(rulePegText, position349)
				}
				if !_rules[ruleAction12]() {
					goto l347
				}
				depth--
				add(ruleTimeEmitterInterval, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 16 TupleEmitterInterval <- <(<TuplesInterval> Action13)> */
		func() bool {
			position350, tokenIndex350, depth350 := position, tokenIndex, depth
			{
				position351 := position
				depth++
				{
					position352 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l350
					}
					depth--
					add(rulePegText, position352)
				}
				if !_rules[ruleAction13]() {
					goto l350
				}
				depth--
				add(ruleTupleEmitterInterval, position351)
			}
			return true
		l350:
			position, tokenIndex, depth = position350, tokenIndex350, depth350
			return false
		},
		/* 17 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action14)> */
		func() bool {
			position353, tokenIndex353, depth353 := position, tokenIndex, depth
			{
				position354 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l353
				}
				if !_rules[rulesp]() {
					goto l353
				}
				{
					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l356
					}
					position++
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune('I') {
						goto l353
					}
					position++
				}
			l355:
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l358
					}
					position++
					goto l357
				l358:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
					if buffer[position] != rune('N') {
						goto l353
					}
					position++
				}
			l357:
				if !_rules[rulesp]() {
					goto l353
				}
				if !_rules[ruleStream]() {
					goto l353
				}
				if !_rules[ruleAction14]() {
					goto l353
				}
				depth--
				add(ruleTupleEmitterFromInterval, position354)
			}
			return true
		l353:
			position, tokenIndex, depth = position353, tokenIndex353, depth353
			return false
		},
		/* 18 Projections <- <(<(Projection sp (',' sp Projection)*)> Action15)> */
		func() bool {
			position359, tokenIndex359, depth359 := position, tokenIndex, depth
			{
				position360 := position
				depth++
				{
					position361 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l359
					}
					if !_rules[rulesp]() {
						goto l359
					}
				l362:
					{
						position363, tokenIndex363, depth363 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l363
						}
						position++
						if !_rules[rulesp]() {
							goto l363
						}
						if !_rules[ruleProjection]() {
							goto l363
						}
						goto l362
					l363:
						position, tokenIndex, depth = position363, tokenIndex363, depth363
					}
					depth--
					add(rulePegText, position361)
				}
				if !_rules[ruleAction15]() {
					goto l359
				}
				depth--
				add(ruleProjections, position360)
			}
			return true
		l359:
			position, tokenIndex, depth = position359, tokenIndex359, depth359
			return false
		},
		/* 19 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position364, tokenIndex364, depth364 := position, tokenIndex, depth
			{
				position365 := position
				depth++
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l367
					}
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if !_rules[ruleExpression]() {
						goto l368
					}
					goto l366
				l368:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if !_rules[ruleWildcard]() {
						goto l364
					}
				}
			l366:
				depth--
				add(ruleProjection, position365)
			}
			return true
		l364:
			position, tokenIndex, depth = position364, tokenIndex364, depth364
			return false
		},
		/* 20 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action16)> */
		func() bool {
			position369, tokenIndex369, depth369 := position, tokenIndex, depth
			{
				position370 := position
				depth++
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l372
					}
					goto l371
				l372:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
					if !_rules[ruleWildcard]() {
						goto l369
					}
				}
			l371:
				if !_rules[rulesp]() {
					goto l369
				}
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('A') {
						goto l369
					}
					position++
				}
			l373:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l376
					}
					position++
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					if buffer[position] != rune('S') {
						goto l369
					}
					position++
				}
			l375:
				if !_rules[rulesp]() {
					goto l369
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l369
				}
				if !_rules[ruleAction16]() {
					goto l369
				}
				depth--
				add(ruleAliasExpression, position370)
			}
			return true
		l369:
			position, tokenIndex, depth = position369, tokenIndex369, depth369
			return false
		},
		/* 21 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action17)> */
		func() bool {
			position377, tokenIndex377, depth377 := position, tokenIndex, depth
			{
				position378 := position
				depth++
				{
					position379 := position
					depth++
					{
						position380, tokenIndex380, depth380 := position, tokenIndex, depth
						{
							position382, tokenIndex382, depth382 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l383
							}
							position++
							goto l382
						l383:
							position, tokenIndex, depth = position382, tokenIndex382, depth382
							if buffer[position] != rune('F') {
								goto l380
							}
							position++
						}
					l382:
						{
							position384, tokenIndex384, depth384 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l385
							}
							position++
							goto l384
						l385:
							position, tokenIndex, depth = position384, tokenIndex384, depth384
							if buffer[position] != rune('R') {
								goto l380
							}
							position++
						}
					l384:
						{
							position386, tokenIndex386, depth386 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l387
							}
							position++
							goto l386
						l387:
							position, tokenIndex, depth = position386, tokenIndex386, depth386
							if buffer[position] != rune('O') {
								goto l380
							}
							position++
						}
					l386:
						{
							position388, tokenIndex388, depth388 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l389
							}
							position++
							goto l388
						l389:
							position, tokenIndex, depth = position388, tokenIndex388, depth388
							if buffer[position] != rune('M') {
								goto l380
							}
							position++
						}
					l388:
						if !_rules[rulesp]() {
							goto l380
						}
						if !_rules[ruleRelations]() {
							goto l380
						}
						if !_rules[rulesp]() {
							goto l380
						}
						goto l381
					l380:
						position, tokenIndex, depth = position380, tokenIndex380, depth380
					}
				l381:
					depth--
					add(rulePegText, position379)
				}
				if !_rules[ruleAction17]() {
					goto l377
				}
				depth--
				add(ruleWindowedFrom, position378)
			}
			return true
		l377:
			position, tokenIndex, depth = position377, tokenIndex377, depth377
			return false
		},
		/* 22 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l393
					}
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if !_rules[ruleTuplesInterval]() {
						goto l390
					}
				}
			l392:
				depth--
				add(ruleInterval, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 23 TimeInterval <- <(NumericLiteral sp SECONDS Action18)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l394
				}
				if !_rules[rulesp]() {
					goto l394
				}
				if !_rules[ruleSECONDS]() {
					goto l394
				}
				if !_rules[ruleAction18]() {
					goto l394
				}
				depth--
				add(ruleTimeInterval, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 24 TuplesInterval <- <(NumericLiteral sp TUPLES Action19)> */
		func() bool {
			position396, tokenIndex396, depth396 := position, tokenIndex, depth
			{
				position397 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l396
				}
				if !_rules[rulesp]() {
					goto l396
				}
				if !_rules[ruleTUPLES]() {
					goto l396
				}
				if !_rules[ruleAction19]() {
					goto l396
				}
				depth--
				add(ruleTuplesInterval, position397)
			}
			return true
		l396:
			position, tokenIndex, depth = position396, tokenIndex396, depth396
			return false
		},
		/* 25 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position398, tokenIndex398, depth398 := position, tokenIndex, depth
			{
				position399 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l398
				}
				if !_rules[rulesp]() {
					goto l398
				}
			l400:
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l401
					}
					position++
					if !_rules[rulesp]() {
						goto l401
					}
					if !_rules[ruleRelationLike]() {
						goto l401
					}
					goto l400
				l401:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
				}
				depth--
				add(ruleRelations, position399)
			}
			return true
		l398:
			position, tokenIndex, depth = position398, tokenIndex398, depth398
			return false
		},
		/* 26 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action20)> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				{
					position404 := position
					depth++
					{
						position405, tokenIndex405, depth405 := position, tokenIndex, depth
						{
							position407, tokenIndex407, depth407 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l408
							}
							position++
							goto l407
						l408:
							position, tokenIndex, depth = position407, tokenIndex407, depth407
							if buffer[position] != rune('W') {
								goto l405
							}
							position++
						}
					l407:
						{
							position409, tokenIndex409, depth409 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l410
							}
							position++
							goto l409
						l410:
							position, tokenIndex, depth = position409, tokenIndex409, depth409
							if buffer[position] != rune('H') {
								goto l405
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
								goto l405
							}
							position++
						}
					l411:
						{
							position413, tokenIndex413, depth413 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l414
							}
							position++
							goto l413
						l414:
							position, tokenIndex, depth = position413, tokenIndex413, depth413
							if buffer[position] != rune('R') {
								goto l405
							}
							position++
						}
					l413:
						{
							position415, tokenIndex415, depth415 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l416
							}
							position++
							goto l415
						l416:
							position, tokenIndex, depth = position415, tokenIndex415, depth415
							if buffer[position] != rune('E') {
								goto l405
							}
							position++
						}
					l415:
						if !_rules[rulesp]() {
							goto l405
						}
						if !_rules[ruleExpression]() {
							goto l405
						}
						goto l406
					l405:
						position, tokenIndex, depth = position405, tokenIndex405, depth405
					}
				l406:
					depth--
					add(rulePegText, position404)
				}
				if !_rules[ruleAction20]() {
					goto l402
				}
				depth--
				add(ruleFilter, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 27 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action21)> */
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
							if buffer[position] != rune('g') {
								goto l423
							}
							position++
							goto l422
						l423:
							position, tokenIndex, depth = position422, tokenIndex422, depth422
							if buffer[position] != rune('G') {
								goto l420
							}
							position++
						}
					l422:
						{
							position424, tokenIndex424, depth424 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l425
							}
							position++
							goto l424
						l425:
							position, tokenIndex, depth = position424, tokenIndex424, depth424
							if buffer[position] != rune('R') {
								goto l420
							}
							position++
						}
					l424:
						{
							position426, tokenIndex426, depth426 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l427
							}
							position++
							goto l426
						l427:
							position, tokenIndex, depth = position426, tokenIndex426, depth426
							if buffer[position] != rune('O') {
								goto l420
							}
							position++
						}
					l426:
						{
							position428, tokenIndex428, depth428 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l429
							}
							position++
							goto l428
						l429:
							position, tokenIndex, depth = position428, tokenIndex428, depth428
							if buffer[position] != rune('U') {
								goto l420
							}
							position++
						}
					l428:
						{
							position430, tokenIndex430, depth430 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l431
							}
							position++
							goto l430
						l431:
							position, tokenIndex, depth = position430, tokenIndex430, depth430
							if buffer[position] != rune('P') {
								goto l420
							}
							position++
						}
					l430:
						if !_rules[rulesp]() {
							goto l420
						}
						{
							position432, tokenIndex432, depth432 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l433
							}
							position++
							goto l432
						l433:
							position, tokenIndex, depth = position432, tokenIndex432, depth432
							if buffer[position] != rune('B') {
								goto l420
							}
							position++
						}
					l432:
						{
							position434, tokenIndex434, depth434 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l435
							}
							position++
							goto l434
						l435:
							position, tokenIndex, depth = position434, tokenIndex434, depth434
							if buffer[position] != rune('Y') {
								goto l420
							}
							position++
						}
					l434:
						if !_rules[rulesp]() {
							goto l420
						}
						if !_rules[ruleGroupList]() {
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
				add(ruleGrouping, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 28 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l436
				}
				if !_rules[rulesp]() {
					goto l436
				}
			l438:
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l439
					}
					position++
					if !_rules[rulesp]() {
						goto l439
					}
					if !_rules[ruleExpression]() {
						goto l439
					}
					goto l438
				l439:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
				}
				depth--
				add(ruleGroupList, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
			return false
		},
		/* 29 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action22)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				{
					position442 := position
					depth++
					{
						position443, tokenIndex443, depth443 := position, tokenIndex, depth
						{
							position445, tokenIndex445, depth445 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l446
							}
							position++
							goto l445
						l446:
							position, tokenIndex, depth = position445, tokenIndex445, depth445
							if buffer[position] != rune('H') {
								goto l443
							}
							position++
						}
					l445:
						{
							position447, tokenIndex447, depth447 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l448
							}
							position++
							goto l447
						l448:
							position, tokenIndex, depth = position447, tokenIndex447, depth447
							if buffer[position] != rune('A') {
								goto l443
							}
							position++
						}
					l447:
						{
							position449, tokenIndex449, depth449 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l450
							}
							position++
							goto l449
						l450:
							position, tokenIndex, depth = position449, tokenIndex449, depth449
							if buffer[position] != rune('V') {
								goto l443
							}
							position++
						}
					l449:
						{
							position451, tokenIndex451, depth451 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l452
							}
							position++
							goto l451
						l452:
							position, tokenIndex, depth = position451, tokenIndex451, depth451
							if buffer[position] != rune('I') {
								goto l443
							}
							position++
						}
					l451:
						{
							position453, tokenIndex453, depth453 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l454
							}
							position++
							goto l453
						l454:
							position, tokenIndex, depth = position453, tokenIndex453, depth453
							if buffer[position] != rune('N') {
								goto l443
							}
							position++
						}
					l453:
						{
							position455, tokenIndex455, depth455 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l456
							}
							position++
							goto l455
						l456:
							position, tokenIndex, depth = position455, tokenIndex455, depth455
							if buffer[position] != rune('G') {
								goto l443
							}
							position++
						}
					l455:
						if !_rules[rulesp]() {
							goto l443
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction22]() {
					goto l440
				}
				depth--
				add(ruleHaving, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 30 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action23))> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l460
					}
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleStreamWindow]() {
						goto l457
					}
					if !_rules[ruleAction23]() {
						goto l457
					}
				}
			l459:
				depth--
				add(ruleRelationLike, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 31 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action24)> */
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
				if !_rules[ruleAction24]() {
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
		/* 32 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action25)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l467
				}
				if !_rules[rulesp]() {
					goto l467
				}
				if buffer[position] != rune('[') {
					goto l467
				}
				position++
				if !_rules[rulesp]() {
					goto l467
				}
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l470
					}
					position++
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if buffer[position] != rune('R') {
						goto l467
					}
					position++
				}
			l469:
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l472
					}
					position++
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if buffer[position] != rune('A') {
						goto l467
					}
					position++
				}
			l471:
				{
					position473, tokenIndex473, depth473 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l474
					}
					position++
					goto l473
				l474:
					position, tokenIndex, depth = position473, tokenIndex473, depth473
					if buffer[position] != rune('N') {
						goto l467
					}
					position++
				}
			l473:
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l476
					}
					position++
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if buffer[position] != rune('G') {
						goto l467
					}
					position++
				}
			l475:
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l478
					}
					position++
					goto l477
				l478:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
					if buffer[position] != rune('E') {
						goto l467
					}
					position++
				}
			l477:
				if !_rules[rulesp]() {
					goto l467
				}
				if !_rules[ruleInterval]() {
					goto l467
				}
				if !_rules[rulesp]() {
					goto l467
				}
				if buffer[position] != rune(']') {
					goto l467
				}
				position++
				if !_rules[ruleAction25]() {
					goto l467
				}
				depth--
				add(ruleStreamWindow, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 33 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l482
					}
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleStream]() {
						goto l479
					}
				}
			l481:
				depth--
				add(ruleStreamLike, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 34 UDSFFuncApp <- <(FuncApp Action26)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l483
				}
				if !_rules[ruleAction26]() {
					goto l483
				}
				depth--
				add(ruleUDSFFuncApp, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 35 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action27)> */
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
						{
							position490, tokenIndex490, depth490 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l491
							}
							position++
							goto l490
						l491:
							position, tokenIndex, depth = position490, tokenIndex490, depth490
							if buffer[position] != rune('W') {
								goto l488
							}
							position++
						}
					l490:
						{
							position492, tokenIndex492, depth492 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l493
							}
							position++
							goto l492
						l493:
							position, tokenIndex, depth = position492, tokenIndex492, depth492
							if buffer[position] != rune('I') {
								goto l488
							}
							position++
						}
					l492:
						{
							position494, tokenIndex494, depth494 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l495
							}
							position++
							goto l494
						l495:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							if buffer[position] != rune('T') {
								goto l488
							}
							position++
						}
					l494:
						{
							position496, tokenIndex496, depth496 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l497
							}
							position++
							goto l496
						l497:
							position, tokenIndex, depth = position496, tokenIndex496, depth496
							if buffer[position] != rune('H') {
								goto l488
							}
							position++
						}
					l496:
						if !_rules[rulesp]() {
							goto l488
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l488
						}
						if !_rules[rulesp]() {
							goto l488
						}
					l498:
						{
							position499, tokenIndex499, depth499 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l499
							}
							position++
							if !_rules[rulesp]() {
								goto l499
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l499
							}
							goto l498
						l499:
							position, tokenIndex, depth = position499, tokenIndex499, depth499
						}
						goto l489
					l488:
						position, tokenIndex, depth = position488, tokenIndex488, depth488
					}
				l489:
					depth--
					add(rulePegText, position487)
				}
				if !_rules[ruleAction27]() {
					goto l485
				}
				depth--
				add(ruleSourceSinkSpecs, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 36 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action28)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l500
				}
				if buffer[position] != rune('=') {
					goto l500
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l500
				}
				if !_rules[ruleAction28]() {
					goto l500
				}
				depth--
				add(ruleSourceSinkParam, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 37 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l505
					}
					goto l504
				l505:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
					if !_rules[ruleLiteral]() {
						goto l502
					}
				}
			l504:
				depth--
				add(ruleSourceSinkParamVal, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 38 PausedOpt <- <(<(Paused / Unpaused)?> Action29)> */
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
						{
							position511, tokenIndex511, depth511 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l512
							}
							goto l511
						l512:
							position, tokenIndex, depth = position511, tokenIndex511, depth511
							if !_rules[ruleUnpaused]() {
								goto l509
							}
						}
					l511:
						goto l510
					l509:
						position, tokenIndex, depth = position509, tokenIndex509, depth509
					}
				l510:
					depth--
					add(rulePegText, position508)
				}
				if !_rules[ruleAction29]() {
					goto l506
				}
				depth--
				add(rulePausedOpt, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 39 Expression <- <orExpr> */
		func() bool {
			position513, tokenIndex513, depth513 := position, tokenIndex, depth
			{
				position514 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l513
				}
				depth--
				add(ruleExpression, position514)
			}
			return true
		l513:
			position, tokenIndex, depth = position513, tokenIndex513, depth513
			return false
		},
		/* 40 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action30)> */
		func() bool {
			position515, tokenIndex515, depth515 := position, tokenIndex, depth
			{
				position516 := position
				depth++
				{
					position517 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l515
					}
					if !_rules[rulesp]() {
						goto l515
					}
					{
						position518, tokenIndex518, depth518 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l518
						}
						if !_rules[rulesp]() {
							goto l518
						}
						if !_rules[ruleandExpr]() {
							goto l518
						}
						goto l519
					l518:
						position, tokenIndex, depth = position518, tokenIndex518, depth518
					}
				l519:
					depth--
					add(rulePegText, position517)
				}
				if !_rules[ruleAction30]() {
					goto l515
				}
				depth--
				add(ruleorExpr, position516)
			}
			return true
		l515:
			position, tokenIndex, depth = position515, tokenIndex515, depth515
			return false
		},
		/* 41 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action31)> */
		func() bool {
			position520, tokenIndex520, depth520 := position, tokenIndex, depth
			{
				position521 := position
				depth++
				{
					position522 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l520
					}
					if !_rules[rulesp]() {
						goto l520
					}
					{
						position523, tokenIndex523, depth523 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l523
						}
						if !_rules[rulesp]() {
							goto l523
						}
						if !_rules[rulenotExpr]() {
							goto l523
						}
						goto l524
					l523:
						position, tokenIndex, depth = position523, tokenIndex523, depth523
					}
				l524:
					depth--
					add(rulePegText, position522)
				}
				if !_rules[ruleAction31]() {
					goto l520
				}
				depth--
				add(ruleandExpr, position521)
			}
			return true
		l520:
			position, tokenIndex, depth = position520, tokenIndex520, depth520
			return false
		},
		/* 42 notExpr <- <(<((Not sp)? comparisonExpr)> Action32)> */
		func() bool {
			position525, tokenIndex525, depth525 := position, tokenIndex, depth
			{
				position526 := position
				depth++
				{
					position527 := position
					depth++
					{
						position528, tokenIndex528, depth528 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l528
						}
						if !_rules[rulesp]() {
							goto l528
						}
						goto l529
					l528:
						position, tokenIndex, depth = position528, tokenIndex528, depth528
					}
				l529:
					if !_rules[rulecomparisonExpr]() {
						goto l525
					}
					depth--
					add(rulePegText, position527)
				}
				if !_rules[ruleAction32]() {
					goto l525
				}
				depth--
				add(rulenotExpr, position526)
			}
			return true
		l525:
			position, tokenIndex, depth = position525, tokenIndex525, depth525
			return false
		},
		/* 43 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action33)> */
		func() bool {
			position530, tokenIndex530, depth530 := position, tokenIndex, depth
			{
				position531 := position
				depth++
				{
					position532 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l530
					}
					if !_rules[rulesp]() {
						goto l530
					}
					{
						position533, tokenIndex533, depth533 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l533
						}
						if !_rules[rulesp]() {
							goto l533
						}
						if !_rules[ruleisExpr]() {
							goto l533
						}
						goto l534
					l533:
						position, tokenIndex, depth = position533, tokenIndex533, depth533
					}
				l534:
					depth--
					add(rulePegText, position532)
				}
				if !_rules[ruleAction33]() {
					goto l530
				}
				depth--
				add(rulecomparisonExpr, position531)
			}
			return true
		l530:
			position, tokenIndex, depth = position530, tokenIndex530, depth530
			return false
		},
		/* 44 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action34)> */
		func() bool {
			position535, tokenIndex535, depth535 := position, tokenIndex, depth
			{
				position536 := position
				depth++
				{
					position537 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l535
					}
					if !_rules[rulesp]() {
						goto l535
					}
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l538
						}
						if !_rules[rulesp]() {
							goto l538
						}
						if !_rules[ruleNullLiteral]() {
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
				if !_rules[ruleAction34]() {
					goto l535
				}
				depth--
				add(ruleisExpr, position536)
			}
			return true
		l535:
			position, tokenIndex, depth = position535, tokenIndex535, depth535
			return false
		},
		/* 45 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action35)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l540
					}
					if !_rules[rulesp]() {
						goto l540
					}
					{
						position543, tokenIndex543, depth543 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l543
						}
						if !_rules[rulesp]() {
							goto l543
						}
						if !_rules[ruleproductExpr]() {
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
				if !_rules[ruleAction35]() {
					goto l540
				}
				depth--
				add(ruletermExpr, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 46 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action36)> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				{
					position547 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l545
					}
					if !_rules[rulesp]() {
						goto l545
					}
					{
						position548, tokenIndex548, depth548 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l548
						}
						if !_rules[rulesp]() {
							goto l548
						}
						if !_rules[ruleminusExpr]() {
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
				if !_rules[ruleAction36]() {
					goto l545
				}
				depth--
				add(ruleproductExpr, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 47 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action37)> */
		func() bool {
			position550, tokenIndex550, depth550 := position, tokenIndex, depth
			{
				position551 := position
				depth++
				{
					position552 := position
					depth++
					{
						position553, tokenIndex553, depth553 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l553
						}
						if !_rules[rulesp]() {
							goto l553
						}
						goto l554
					l553:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
					}
				l554:
					if !_rules[rulebaseExpr]() {
						goto l550
					}
					depth--
					add(rulePegText, position552)
				}
				if !_rules[ruleAction37]() {
					goto l550
				}
				depth--
				add(ruleminusExpr, position551)
			}
			return true
		l550:
			position, tokenIndex, depth = position550, tokenIndex550, depth550
			return false
		},
		/* 48 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position555, tokenIndex555, depth555 := position, tokenIndex, depth
			{
				position556 := position
				depth++
				{
					position557, tokenIndex557, depth557 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l558
					}
					position++
					if !_rules[rulesp]() {
						goto l558
					}
					if !_rules[ruleExpression]() {
						goto l558
					}
					if !_rules[rulesp]() {
						goto l558
					}
					if buffer[position] != rune(')') {
						goto l558
					}
					position++
					goto l557
				l558:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if !_rules[ruleBooleanLiteral]() {
						goto l559
					}
					goto l557
				l559:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if !_rules[ruleNullLiteral]() {
						goto l560
					}
					goto l557
				l560:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if !_rules[ruleFuncApp]() {
						goto l561
					}
					goto l557
				l561:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if !_rules[ruleRowMeta]() {
						goto l562
					}
					goto l557
				l562:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if !_rules[ruleRowValue]() {
						goto l563
					}
					goto l557
				l563:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if !_rules[ruleLiteral]() {
						goto l555
					}
				}
			l557:
				depth--
				add(rulebaseExpr, position556)
			}
			return true
		l555:
			position, tokenIndex, depth = position555, tokenIndex555, depth555
			return false
		},
		/* 49 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action38)> */
		func() bool {
			position564, tokenIndex564, depth564 := position, tokenIndex, depth
			{
				position565 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l564
				}
				if !_rules[rulesp]() {
					goto l564
				}
				if buffer[position] != rune('(') {
					goto l564
				}
				position++
				if !_rules[rulesp]() {
					goto l564
				}
				if !_rules[ruleFuncParams]() {
					goto l564
				}
				if !_rules[rulesp]() {
					goto l564
				}
				if buffer[position] != rune(')') {
					goto l564
				}
				position++
				if !_rules[ruleAction38]() {
					goto l564
				}
				depth--
				add(ruleFuncApp, position565)
			}
			return true
		l564:
			position, tokenIndex, depth = position564, tokenIndex564, depth564
			return false
		},
		/* 50 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action39)> */
		func() bool {
			position566, tokenIndex566, depth566 := position, tokenIndex, depth
			{
				position567 := position
				depth++
				{
					position568 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l566
					}
					if !_rules[rulesp]() {
						goto l566
					}
				l569:
					{
						position570, tokenIndex570, depth570 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l570
						}
						position++
						if !_rules[rulesp]() {
							goto l570
						}
						if !_rules[ruleExpression]() {
							goto l570
						}
						goto l569
					l570:
						position, tokenIndex, depth = position570, tokenIndex570, depth570
					}
					depth--
					add(rulePegText, position568)
				}
				if !_rules[ruleAction39]() {
					goto l566
				}
				depth--
				add(ruleFuncParams, position567)
			}
			return true
		l566:
			position, tokenIndex, depth = position566, tokenIndex566, depth566
			return false
		},
		/* 51 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position571, tokenIndex571, depth571 := position, tokenIndex, depth
			{
				position572 := position
				depth++
				{
					position573, tokenIndex573, depth573 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l574
					}
					goto l573
				l574:
					position, tokenIndex, depth = position573, tokenIndex573, depth573
					if !_rules[ruleNumericLiteral]() {
						goto l575
					}
					goto l573
				l575:
					position, tokenIndex, depth = position573, tokenIndex573, depth573
					if !_rules[ruleStringLiteral]() {
						goto l571
					}
				}
			l573:
				depth--
				add(ruleLiteral, position572)
			}
			return true
		l571:
			position, tokenIndex, depth = position571, tokenIndex571, depth571
			return false
		},
		/* 52 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position576, tokenIndex576, depth576 := position, tokenIndex, depth
			{
				position577 := position
				depth++
				{
					position578, tokenIndex578, depth578 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l579
					}
					goto l578
				l579:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
					if !_rules[ruleNotEqual]() {
						goto l580
					}
					goto l578
				l580:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
					if !_rules[ruleLessOrEqual]() {
						goto l581
					}
					goto l578
				l581:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
					if !_rules[ruleLess]() {
						goto l582
					}
					goto l578
				l582:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
					if !_rules[ruleGreaterOrEqual]() {
						goto l583
					}
					goto l578
				l583:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
					if !_rules[ruleGreater]() {
						goto l584
					}
					goto l578
				l584:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
					if !_rules[ruleNotEqual]() {
						goto l576
					}
				}
			l578:
				depth--
				add(ruleComparisonOp, position577)
			}
			return true
		l576:
			position, tokenIndex, depth = position576, tokenIndex576, depth576
			return false
		},
		/* 53 IsOp <- <(IsNot / Is)> */
		func() bool {
			position585, tokenIndex585, depth585 := position, tokenIndex, depth
			{
				position586 := position
				depth++
				{
					position587, tokenIndex587, depth587 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l588
					}
					goto l587
				l588:
					position, tokenIndex, depth = position587, tokenIndex587, depth587
					if !_rules[ruleIs]() {
						goto l585
					}
				}
			l587:
				depth--
				add(ruleIsOp, position586)
			}
			return true
		l585:
			position, tokenIndex, depth = position585, tokenIndex585, depth585
			return false
		},
		/* 54 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position589, tokenIndex589, depth589 := position, tokenIndex, depth
			{
				position590 := position
				depth++
				{
					position591, tokenIndex591, depth591 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l592
					}
					goto l591
				l592:
					position, tokenIndex, depth = position591, tokenIndex591, depth591
					if !_rules[ruleMinus]() {
						goto l589
					}
				}
			l591:
				depth--
				add(rulePlusMinusOp, position590)
			}
			return true
		l589:
			position, tokenIndex, depth = position589, tokenIndex589, depth589
			return false
		},
		/* 55 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position593, tokenIndex593, depth593 := position, tokenIndex, depth
			{
				position594 := position
				depth++
				{
					position595, tokenIndex595, depth595 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l596
					}
					goto l595
				l596:
					position, tokenIndex, depth = position595, tokenIndex595, depth595
					if !_rules[ruleDivide]() {
						goto l597
					}
					goto l595
				l597:
					position, tokenIndex, depth = position595, tokenIndex595, depth595
					if !_rules[ruleModulo]() {
						goto l593
					}
				}
			l595:
				depth--
				add(ruleMultDivOp, position594)
			}
			return true
		l593:
			position, tokenIndex, depth = position593, tokenIndex593, depth593
			return false
		},
		/* 56 Stream <- <(<ident> Action40)> */
		func() bool {
			position598, tokenIndex598, depth598 := position, tokenIndex, depth
			{
				position599 := position
				depth++
				{
					position600 := position
					depth++
					if !_rules[ruleident]() {
						goto l598
					}
					depth--
					add(rulePegText, position600)
				}
				if !_rules[ruleAction40]() {
					goto l598
				}
				depth--
				add(ruleStream, position599)
			}
			return true
		l598:
			position, tokenIndex, depth = position598, tokenIndex598, depth598
			return false
		},
		/* 57 RowMeta <- <RowTimestamp> */
		func() bool {
			position601, tokenIndex601, depth601 := position, tokenIndex, depth
			{
				position602 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l601
				}
				depth--
				add(ruleRowMeta, position602)
			}
			return true
		l601:
			position, tokenIndex, depth = position601, tokenIndex601, depth601
			return false
		},
		/* 58 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action41)> */
		func() bool {
			position603, tokenIndex603, depth603 := position, tokenIndex, depth
			{
				position604 := position
				depth++
				{
					position605 := position
					depth++
					{
						position606, tokenIndex606, depth606 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l606
						}
						if buffer[position] != rune(':') {
							goto l606
						}
						position++
						goto l607
					l606:
						position, tokenIndex, depth = position606, tokenIndex606, depth606
					}
				l607:
					if buffer[position] != rune('t') {
						goto l603
					}
					position++
					if buffer[position] != rune('s') {
						goto l603
					}
					position++
					if buffer[position] != rune('(') {
						goto l603
					}
					position++
					if buffer[position] != rune(')') {
						goto l603
					}
					position++
					depth--
					add(rulePegText, position605)
				}
				if !_rules[ruleAction41]() {
					goto l603
				}
				depth--
				add(ruleRowTimestamp, position604)
			}
			return true
		l603:
			position, tokenIndex, depth = position603, tokenIndex603, depth603
			return false
		},
		/* 59 RowValue <- <(<((ident ':')? jsonPath)> Action42)> */
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
						if !_rules[ruleident]() {
							goto l611
						}
						if buffer[position] != rune(':') {
							goto l611
						}
						position++
						goto l612
					l611:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
					}
				l612:
					if !_rules[rulejsonPath]() {
						goto l608
					}
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction42]() {
					goto l608
				}
				depth--
				add(ruleRowValue, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 60 NumericLiteral <- <(<('-'? [0-9]+)> Action43)> */
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
						if buffer[position] != rune('-') {
							goto l616
						}
						position++
						goto l617
					l616:
						position, tokenIndex, depth = position616, tokenIndex616, depth616
					}
				l617:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l613
					}
					position++
				l618:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l619
						}
						position++
						goto l618
					l619:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
					}
					depth--
					add(rulePegText, position615)
				}
				if !_rules[ruleAction43]() {
					goto l613
				}
				depth--
				add(ruleNumericLiteral, position614)
			}
			return true
		l613:
			position, tokenIndex, depth = position613, tokenIndex613, depth613
			return false
		},
		/* 61 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action44)> */
		func() bool {
			position620, tokenIndex620, depth620 := position, tokenIndex, depth
			{
				position621 := position
				depth++
				{
					position622 := position
					depth++
					{
						position623, tokenIndex623, depth623 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l623
						}
						position++
						goto l624
					l623:
						position, tokenIndex, depth = position623, tokenIndex623, depth623
					}
				l624:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l620
					}
					position++
				l625:
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l626
						}
						position++
						goto l625
					l626:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
					}
					if buffer[position] != rune('.') {
						goto l620
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l620
					}
					position++
				l627:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l628
						}
						position++
						goto l627
					l628:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
					}
					depth--
					add(rulePegText, position622)
				}
				if !_rules[ruleAction44]() {
					goto l620
				}
				depth--
				add(ruleFloatLiteral, position621)
			}
			return true
		l620:
			position, tokenIndex, depth = position620, tokenIndex620, depth620
			return false
		},
		/* 62 Function <- <(<ident> Action45)> */
		func() bool {
			position629, tokenIndex629, depth629 := position, tokenIndex, depth
			{
				position630 := position
				depth++
				{
					position631 := position
					depth++
					if !_rules[ruleident]() {
						goto l629
					}
					depth--
					add(rulePegText, position631)
				}
				if !_rules[ruleAction45]() {
					goto l629
				}
				depth--
				add(ruleFunction, position630)
			}
			return true
		l629:
			position, tokenIndex, depth = position629, tokenIndex629, depth629
			return false
		},
		/* 63 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action46)> */
		func() bool {
			position632, tokenIndex632, depth632 := position, tokenIndex, depth
			{
				position633 := position
				depth++
				{
					position634 := position
					depth++
					{
						position635, tokenIndex635, depth635 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l636
						}
						position++
						goto l635
					l636:
						position, tokenIndex, depth = position635, tokenIndex635, depth635
						if buffer[position] != rune('N') {
							goto l632
						}
						position++
					}
				l635:
					{
						position637, tokenIndex637, depth637 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l638
						}
						position++
						goto l637
					l638:
						position, tokenIndex, depth = position637, tokenIndex637, depth637
						if buffer[position] != rune('U') {
							goto l632
						}
						position++
					}
				l637:
					{
						position639, tokenIndex639, depth639 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l640
						}
						position++
						goto l639
					l640:
						position, tokenIndex, depth = position639, tokenIndex639, depth639
						if buffer[position] != rune('L') {
							goto l632
						}
						position++
					}
				l639:
					{
						position641, tokenIndex641, depth641 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l642
						}
						position++
						goto l641
					l642:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
						if buffer[position] != rune('L') {
							goto l632
						}
						position++
					}
				l641:
					depth--
					add(rulePegText, position634)
				}
				if !_rules[ruleAction46]() {
					goto l632
				}
				depth--
				add(ruleNullLiteral, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 64 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645, tokenIndex645, depth645 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l646
					}
					goto l645
				l646:
					position, tokenIndex, depth = position645, tokenIndex645, depth645
					if !_rules[ruleFALSE]() {
						goto l643
					}
				}
			l645:
				depth--
				add(ruleBooleanLiteral, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 65 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action47)> */
		func() bool {
			position647, tokenIndex647, depth647 := position, tokenIndex, depth
			{
				position648 := position
				depth++
				{
					position649 := position
					depth++
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
							goto l647
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
							goto l647
						}
						position++
					}
				l652:
					{
						position654, tokenIndex654, depth654 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l655
						}
						position++
						goto l654
					l655:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
						if buffer[position] != rune('U') {
							goto l647
						}
						position++
					}
				l654:
					{
						position656, tokenIndex656, depth656 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l657
						}
						position++
						goto l656
					l657:
						position, tokenIndex, depth = position656, tokenIndex656, depth656
						if buffer[position] != rune('E') {
							goto l647
						}
						position++
					}
				l656:
					depth--
					add(rulePegText, position649)
				}
				if !_rules[ruleAction47]() {
					goto l647
				}
				depth--
				add(ruleTRUE, position648)
			}
			return true
		l647:
			position, tokenIndex, depth = position647, tokenIndex647, depth647
			return false
		},
		/* 66 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action48)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				{
					position660 := position
					depth++
					{
						position661, tokenIndex661, depth661 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l662
						}
						position++
						goto l661
					l662:
						position, tokenIndex, depth = position661, tokenIndex661, depth661
						if buffer[position] != rune('F') {
							goto l658
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
							goto l658
						}
						position++
					}
				l663:
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l666
						}
						position++
						goto l665
					l666:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if buffer[position] != rune('L') {
							goto l658
						}
						position++
					}
				l665:
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if buffer[position] != rune('S') {
							goto l658
						}
						position++
					}
				l667:
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l670
						}
						position++
						goto l669
					l670:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
						if buffer[position] != rune('E') {
							goto l658
						}
						position++
					}
				l669:
					depth--
					add(rulePegText, position660)
				}
				if !_rules[ruleAction48]() {
					goto l658
				}
				depth--
				add(ruleFALSE, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 67 Wildcard <- <(<'*'> Action49)> */
		func() bool {
			position671, tokenIndex671, depth671 := position, tokenIndex, depth
			{
				position672 := position
				depth++
				{
					position673 := position
					depth++
					if buffer[position] != rune('*') {
						goto l671
					}
					position++
					depth--
					add(rulePegText, position673)
				}
				if !_rules[ruleAction49]() {
					goto l671
				}
				depth--
				add(ruleWildcard, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 68 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action50)> */
		func() bool {
			position674, tokenIndex674, depth674 := position, tokenIndex, depth
			{
				position675 := position
				depth++
				{
					position676 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l674
					}
					position++
				l677:
					{
						position678, tokenIndex678, depth678 := position, tokenIndex, depth
						{
							position679, tokenIndex679, depth679 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l680
							}
							position++
							if buffer[position] != rune('\'') {
								goto l680
							}
							position++
							goto l679
						l680:
							position, tokenIndex, depth = position679, tokenIndex679, depth679
							{
								position681, tokenIndex681, depth681 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l681
								}
								position++
								goto l678
							l681:
								position, tokenIndex, depth = position681, tokenIndex681, depth681
							}
							if !matchDot() {
								goto l678
							}
						}
					l679:
						goto l677
					l678:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
					}
					if buffer[position] != rune('\'') {
						goto l674
					}
					position++
					depth--
					add(rulePegText, position676)
				}
				if !_rules[ruleAction50]() {
					goto l674
				}
				depth--
				add(ruleStringLiteral, position675)
			}
			return true
		l674:
			position, tokenIndex, depth = position674, tokenIndex674, depth674
			return false
		},
		/* 69 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action51)> */
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
						if buffer[position] != rune('i') {
							goto l686
						}
						position++
						goto l685
					l686:
						position, tokenIndex, depth = position685, tokenIndex685, depth685
						if buffer[position] != rune('I') {
							goto l682
						}
						position++
					}
				l685:
					{
						position687, tokenIndex687, depth687 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l688
						}
						position++
						goto l687
					l688:
						position, tokenIndex, depth = position687, tokenIndex687, depth687
						if buffer[position] != rune('S') {
							goto l682
						}
						position++
					}
				l687:
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l690
						}
						position++
						goto l689
					l690:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
						if buffer[position] != rune('T') {
							goto l682
						}
						position++
					}
				l689:
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
							goto l682
						}
						position++
					}
				l691:
					{
						position693, tokenIndex693, depth693 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l694
						}
						position++
						goto l693
					l694:
						position, tokenIndex, depth = position693, tokenIndex693, depth693
						if buffer[position] != rune('E') {
							goto l682
						}
						position++
					}
				l693:
					{
						position695, tokenIndex695, depth695 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
						if buffer[position] != rune('A') {
							goto l682
						}
						position++
					}
				l695:
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('M') {
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
				add(ruleISTREAM, position683)
			}
			return true
		l682:
			position, tokenIndex, depth = position682, tokenIndex682, depth682
			return false
		},
		/* 70 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action52)> */
		func() bool {
			position699, tokenIndex699, depth699 := position, tokenIndex, depth
			{
				position700 := position
				depth++
				{
					position701 := position
					depth++
					{
						position702, tokenIndex702, depth702 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l703
						}
						position++
						goto l702
					l703:
						position, tokenIndex, depth = position702, tokenIndex702, depth702
						if buffer[position] != rune('D') {
							goto l699
						}
						position++
					}
				l702:
					{
						position704, tokenIndex704, depth704 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l705
						}
						position++
						goto l704
					l705:
						position, tokenIndex, depth = position704, tokenIndex704, depth704
						if buffer[position] != rune('S') {
							goto l699
						}
						position++
					}
				l704:
					{
						position706, tokenIndex706, depth706 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l707
						}
						position++
						goto l706
					l707:
						position, tokenIndex, depth = position706, tokenIndex706, depth706
						if buffer[position] != rune('T') {
							goto l699
						}
						position++
					}
				l706:
					{
						position708, tokenIndex708, depth708 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l709
						}
						position++
						goto l708
					l709:
						position, tokenIndex, depth = position708, tokenIndex708, depth708
						if buffer[position] != rune('R') {
							goto l699
						}
						position++
					}
				l708:
					{
						position710, tokenIndex710, depth710 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l711
						}
						position++
						goto l710
					l711:
						position, tokenIndex, depth = position710, tokenIndex710, depth710
						if buffer[position] != rune('E') {
							goto l699
						}
						position++
					}
				l710:
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
							goto l699
						}
						position++
					}
				l712:
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l715
						}
						position++
						goto l714
					l715:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
						if buffer[position] != rune('M') {
							goto l699
						}
						position++
					}
				l714:
					depth--
					add(rulePegText, position701)
				}
				if !_rules[ruleAction52]() {
					goto l699
				}
				depth--
				add(ruleDSTREAM, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 71 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action53)> */
		func() bool {
			position716, tokenIndex716, depth716 := position, tokenIndex, depth
			{
				position717 := position
				depth++
				{
					position718 := position
					depth++
					{
						position719, tokenIndex719, depth719 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l720
						}
						position++
						goto l719
					l720:
						position, tokenIndex, depth = position719, tokenIndex719, depth719
						if buffer[position] != rune('R') {
							goto l716
						}
						position++
					}
				l719:
					{
						position721, tokenIndex721, depth721 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l722
						}
						position++
						goto l721
					l722:
						position, tokenIndex, depth = position721, tokenIndex721, depth721
						if buffer[position] != rune('S') {
							goto l716
						}
						position++
					}
				l721:
					{
						position723, tokenIndex723, depth723 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l724
						}
						position++
						goto l723
					l724:
						position, tokenIndex, depth = position723, tokenIndex723, depth723
						if buffer[position] != rune('T') {
							goto l716
						}
						position++
					}
				l723:
					{
						position725, tokenIndex725, depth725 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l726
						}
						position++
						goto l725
					l726:
						position, tokenIndex, depth = position725, tokenIndex725, depth725
						if buffer[position] != rune('R') {
							goto l716
						}
						position++
					}
				l725:
					{
						position727, tokenIndex727, depth727 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l728
						}
						position++
						goto l727
					l728:
						position, tokenIndex, depth = position727, tokenIndex727, depth727
						if buffer[position] != rune('E') {
							goto l716
						}
						position++
					}
				l727:
					{
						position729, tokenIndex729, depth729 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l730
						}
						position++
						goto l729
					l730:
						position, tokenIndex, depth = position729, tokenIndex729, depth729
						if buffer[position] != rune('A') {
							goto l716
						}
						position++
					}
				l729:
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l732
						}
						position++
						goto l731
					l732:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
						if buffer[position] != rune('M') {
							goto l716
						}
						position++
					}
				l731:
					depth--
					add(rulePegText, position718)
				}
				if !_rules[ruleAction53]() {
					goto l716
				}
				depth--
				add(ruleRSTREAM, position717)
			}
			return true
		l716:
			position, tokenIndex, depth = position716, tokenIndex716, depth716
			return false
		},
		/* 72 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action54)> */
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
						if buffer[position] != rune('t') {
							goto l737
						}
						position++
						goto l736
					l737:
						position, tokenIndex, depth = position736, tokenIndex736, depth736
						if buffer[position] != rune('T') {
							goto l733
						}
						position++
					}
				l736:
					{
						position738, tokenIndex738, depth738 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l739
						}
						position++
						goto l738
					l739:
						position, tokenIndex, depth = position738, tokenIndex738, depth738
						if buffer[position] != rune('U') {
							goto l733
						}
						position++
					}
				l738:
					{
						position740, tokenIndex740, depth740 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l741
						}
						position++
						goto l740
					l741:
						position, tokenIndex, depth = position740, tokenIndex740, depth740
						if buffer[position] != rune('P') {
							goto l733
						}
						position++
					}
				l740:
					{
						position742, tokenIndex742, depth742 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l743
						}
						position++
						goto l742
					l743:
						position, tokenIndex, depth = position742, tokenIndex742, depth742
						if buffer[position] != rune('L') {
							goto l733
						}
						position++
					}
				l742:
					{
						position744, tokenIndex744, depth744 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l745
						}
						position++
						goto l744
					l745:
						position, tokenIndex, depth = position744, tokenIndex744, depth744
						if buffer[position] != rune('E') {
							goto l733
						}
						position++
					}
				l744:
					{
						position746, tokenIndex746, depth746 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l747
						}
						position++
						goto l746
					l747:
						position, tokenIndex, depth = position746, tokenIndex746, depth746
						if buffer[position] != rune('S') {
							goto l733
						}
						position++
					}
				l746:
					depth--
					add(rulePegText, position735)
				}
				if !_rules[ruleAction54]() {
					goto l733
				}
				depth--
				add(ruleTUPLES, position734)
			}
			return true
		l733:
			position, tokenIndex, depth = position733, tokenIndex733, depth733
			return false
		},
		/* 73 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action55)> */
		func() bool {
			position748, tokenIndex748, depth748 := position, tokenIndex, depth
			{
				position749 := position
				depth++
				{
					position750 := position
					depth++
					{
						position751, tokenIndex751, depth751 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l752
						}
						position++
						goto l751
					l752:
						position, tokenIndex, depth = position751, tokenIndex751, depth751
						if buffer[position] != rune('S') {
							goto l748
						}
						position++
					}
				l751:
					{
						position753, tokenIndex753, depth753 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l754
						}
						position++
						goto l753
					l754:
						position, tokenIndex, depth = position753, tokenIndex753, depth753
						if buffer[position] != rune('E') {
							goto l748
						}
						position++
					}
				l753:
					{
						position755, tokenIndex755, depth755 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l756
						}
						position++
						goto l755
					l756:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						if buffer[position] != rune('C') {
							goto l748
						}
						position++
					}
				l755:
					{
						position757, tokenIndex757, depth757 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l758
						}
						position++
						goto l757
					l758:
						position, tokenIndex, depth = position757, tokenIndex757, depth757
						if buffer[position] != rune('O') {
							goto l748
						}
						position++
					}
				l757:
					{
						position759, tokenIndex759, depth759 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l760
						}
						position++
						goto l759
					l760:
						position, tokenIndex, depth = position759, tokenIndex759, depth759
						if buffer[position] != rune('N') {
							goto l748
						}
						position++
					}
				l759:
					{
						position761, tokenIndex761, depth761 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l762
						}
						position++
						goto l761
					l762:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if buffer[position] != rune('D') {
							goto l748
						}
						position++
					}
				l761:
					{
						position763, tokenIndex763, depth763 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l764
						}
						position++
						goto l763
					l764:
						position, tokenIndex, depth = position763, tokenIndex763, depth763
						if buffer[position] != rune('S') {
							goto l748
						}
						position++
					}
				l763:
					depth--
					add(rulePegText, position750)
				}
				if !_rules[ruleAction55]() {
					goto l748
				}
				depth--
				add(ruleSECONDS, position749)
			}
			return true
		l748:
			position, tokenIndex, depth = position748, tokenIndex748, depth748
			return false
		},
		/* 74 StreamIdentifier <- <(<ident> Action56)> */
		func() bool {
			position765, tokenIndex765, depth765 := position, tokenIndex, depth
			{
				position766 := position
				depth++
				{
					position767 := position
					depth++
					if !_rules[ruleident]() {
						goto l765
					}
					depth--
					add(rulePegText, position767)
				}
				if !_rules[ruleAction56]() {
					goto l765
				}
				depth--
				add(ruleStreamIdentifier, position766)
			}
			return true
		l765:
			position, tokenIndex, depth = position765, tokenIndex765, depth765
			return false
		},
		/* 75 SourceSinkType <- <(<ident> Action57)> */
		func() bool {
			position768, tokenIndex768, depth768 := position, tokenIndex, depth
			{
				position769 := position
				depth++
				{
					position770 := position
					depth++
					if !_rules[ruleident]() {
						goto l768
					}
					depth--
					add(rulePegText, position770)
				}
				if !_rules[ruleAction57]() {
					goto l768
				}
				depth--
				add(ruleSourceSinkType, position769)
			}
			return true
		l768:
			position, tokenIndex, depth = position768, tokenIndex768, depth768
			return false
		},
		/* 76 SourceSinkParamKey <- <(<ident> Action58)> */
		func() bool {
			position771, tokenIndex771, depth771 := position, tokenIndex, depth
			{
				position772 := position
				depth++
				{
					position773 := position
					depth++
					if !_rules[ruleident]() {
						goto l771
					}
					depth--
					add(rulePegText, position773)
				}
				if !_rules[ruleAction58]() {
					goto l771
				}
				depth--
				add(ruleSourceSinkParamKey, position772)
			}
			return true
		l771:
			position, tokenIndex, depth = position771, tokenIndex771, depth771
			return false
		},
		/* 77 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action59)> */
		func() bool {
			position774, tokenIndex774, depth774 := position, tokenIndex, depth
			{
				position775 := position
				depth++
				{
					position776 := position
					depth++
					{
						position777, tokenIndex777, depth777 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l778
						}
						position++
						goto l777
					l778:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
						if buffer[position] != rune('P') {
							goto l774
						}
						position++
					}
				l777:
					{
						position779, tokenIndex779, depth779 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l780
						}
						position++
						goto l779
					l780:
						position, tokenIndex, depth = position779, tokenIndex779, depth779
						if buffer[position] != rune('A') {
							goto l774
						}
						position++
					}
				l779:
					{
						position781, tokenIndex781, depth781 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l782
						}
						position++
						goto l781
					l782:
						position, tokenIndex, depth = position781, tokenIndex781, depth781
						if buffer[position] != rune('U') {
							goto l774
						}
						position++
					}
				l781:
					{
						position783, tokenIndex783, depth783 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l784
						}
						position++
						goto l783
					l784:
						position, tokenIndex, depth = position783, tokenIndex783, depth783
						if buffer[position] != rune('S') {
							goto l774
						}
						position++
					}
				l783:
					{
						position785, tokenIndex785, depth785 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l786
						}
						position++
						goto l785
					l786:
						position, tokenIndex, depth = position785, tokenIndex785, depth785
						if buffer[position] != rune('E') {
							goto l774
						}
						position++
					}
				l785:
					{
						position787, tokenIndex787, depth787 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l788
						}
						position++
						goto l787
					l788:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
						if buffer[position] != rune('D') {
							goto l774
						}
						position++
					}
				l787:
					depth--
					add(rulePegText, position776)
				}
				if !_rules[ruleAction59]() {
					goto l774
				}
				depth--
				add(rulePaused, position775)
			}
			return true
		l774:
			position, tokenIndex, depth = position774, tokenIndex774, depth774
			return false
		},
		/* 78 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action60)> */
		func() bool {
			position789, tokenIndex789, depth789 := position, tokenIndex, depth
			{
				position790 := position
				depth++
				{
					position791 := position
					depth++
					{
						position792, tokenIndex792, depth792 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l793
						}
						position++
						goto l792
					l793:
						position, tokenIndex, depth = position792, tokenIndex792, depth792
						if buffer[position] != rune('U') {
							goto l789
						}
						position++
					}
				l792:
					{
						position794, tokenIndex794, depth794 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l795
						}
						position++
						goto l794
					l795:
						position, tokenIndex, depth = position794, tokenIndex794, depth794
						if buffer[position] != rune('N') {
							goto l789
						}
						position++
					}
				l794:
					{
						position796, tokenIndex796, depth796 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l797
						}
						position++
						goto l796
					l797:
						position, tokenIndex, depth = position796, tokenIndex796, depth796
						if buffer[position] != rune('P') {
							goto l789
						}
						position++
					}
				l796:
					{
						position798, tokenIndex798, depth798 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l799
						}
						position++
						goto l798
					l799:
						position, tokenIndex, depth = position798, tokenIndex798, depth798
						if buffer[position] != rune('A') {
							goto l789
						}
						position++
					}
				l798:
					{
						position800, tokenIndex800, depth800 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l801
						}
						position++
						goto l800
					l801:
						position, tokenIndex, depth = position800, tokenIndex800, depth800
						if buffer[position] != rune('U') {
							goto l789
						}
						position++
					}
				l800:
					{
						position802, tokenIndex802, depth802 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l803
						}
						position++
						goto l802
					l803:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if buffer[position] != rune('S') {
							goto l789
						}
						position++
					}
				l802:
					{
						position804, tokenIndex804, depth804 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l805
						}
						position++
						goto l804
					l805:
						position, tokenIndex, depth = position804, tokenIndex804, depth804
						if buffer[position] != rune('E') {
							goto l789
						}
						position++
					}
				l804:
					{
						position806, tokenIndex806, depth806 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l807
						}
						position++
						goto l806
					l807:
						position, tokenIndex, depth = position806, tokenIndex806, depth806
						if buffer[position] != rune('D') {
							goto l789
						}
						position++
					}
				l806:
					depth--
					add(rulePegText, position791)
				}
				if !_rules[ruleAction60]() {
					goto l789
				}
				depth--
				add(ruleUnpaused, position790)
			}
			return true
		l789:
			position, tokenIndex, depth = position789, tokenIndex789, depth789
			return false
		},
		/* 79 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action61)> */
		func() bool {
			position808, tokenIndex808, depth808 := position, tokenIndex, depth
			{
				position809 := position
				depth++
				{
					position810 := position
					depth++
					{
						position811, tokenIndex811, depth811 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l812
						}
						position++
						goto l811
					l812:
						position, tokenIndex, depth = position811, tokenIndex811, depth811
						if buffer[position] != rune('O') {
							goto l808
						}
						position++
					}
				l811:
					{
						position813, tokenIndex813, depth813 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l814
						}
						position++
						goto l813
					l814:
						position, tokenIndex, depth = position813, tokenIndex813, depth813
						if buffer[position] != rune('R') {
							goto l808
						}
						position++
					}
				l813:
					depth--
					add(rulePegText, position810)
				}
				if !_rules[ruleAction61]() {
					goto l808
				}
				depth--
				add(ruleOr, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 80 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action62)> */
		func() bool {
			position815, tokenIndex815, depth815 := position, tokenIndex, depth
			{
				position816 := position
				depth++
				{
					position817 := position
					depth++
					{
						position818, tokenIndex818, depth818 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l819
						}
						position++
						goto l818
					l819:
						position, tokenIndex, depth = position818, tokenIndex818, depth818
						if buffer[position] != rune('A') {
							goto l815
						}
						position++
					}
				l818:
					{
						position820, tokenIndex820, depth820 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l821
						}
						position++
						goto l820
					l821:
						position, tokenIndex, depth = position820, tokenIndex820, depth820
						if buffer[position] != rune('N') {
							goto l815
						}
						position++
					}
				l820:
					{
						position822, tokenIndex822, depth822 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l823
						}
						position++
						goto l822
					l823:
						position, tokenIndex, depth = position822, tokenIndex822, depth822
						if buffer[position] != rune('D') {
							goto l815
						}
						position++
					}
				l822:
					depth--
					add(rulePegText, position817)
				}
				if !_rules[ruleAction62]() {
					goto l815
				}
				depth--
				add(ruleAnd, position816)
			}
			return true
		l815:
			position, tokenIndex, depth = position815, tokenIndex815, depth815
			return false
		},
		/* 81 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action63)> */
		func() bool {
			position824, tokenIndex824, depth824 := position, tokenIndex, depth
			{
				position825 := position
				depth++
				{
					position826 := position
					depth++
					{
						position827, tokenIndex827, depth827 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l828
						}
						position++
						goto l827
					l828:
						position, tokenIndex, depth = position827, tokenIndex827, depth827
						if buffer[position] != rune('N') {
							goto l824
						}
						position++
					}
				l827:
					{
						position829, tokenIndex829, depth829 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l830
						}
						position++
						goto l829
					l830:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						if buffer[position] != rune('O') {
							goto l824
						}
						position++
					}
				l829:
					{
						position831, tokenIndex831, depth831 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l832
						}
						position++
						goto l831
					l832:
						position, tokenIndex, depth = position831, tokenIndex831, depth831
						if buffer[position] != rune('T') {
							goto l824
						}
						position++
					}
				l831:
					depth--
					add(rulePegText, position826)
				}
				if !_rules[ruleAction63]() {
					goto l824
				}
				depth--
				add(ruleNot, position825)
			}
			return true
		l824:
			position, tokenIndex, depth = position824, tokenIndex824, depth824
			return false
		},
		/* 82 Equal <- <(<'='> Action64)> */
		func() bool {
			position833, tokenIndex833, depth833 := position, tokenIndex, depth
			{
				position834 := position
				depth++
				{
					position835 := position
					depth++
					if buffer[position] != rune('=') {
						goto l833
					}
					position++
					depth--
					add(rulePegText, position835)
				}
				if !_rules[ruleAction64]() {
					goto l833
				}
				depth--
				add(ruleEqual, position834)
			}
			return true
		l833:
			position, tokenIndex, depth = position833, tokenIndex833, depth833
			return false
		},
		/* 83 Less <- <(<'<'> Action65)> */
		func() bool {
			position836, tokenIndex836, depth836 := position, tokenIndex, depth
			{
				position837 := position
				depth++
				{
					position838 := position
					depth++
					if buffer[position] != rune('<') {
						goto l836
					}
					position++
					depth--
					add(rulePegText, position838)
				}
				if !_rules[ruleAction65]() {
					goto l836
				}
				depth--
				add(ruleLess, position837)
			}
			return true
		l836:
			position, tokenIndex, depth = position836, tokenIndex836, depth836
			return false
		},
		/* 84 LessOrEqual <- <(<('<' '=')> Action66)> */
		func() bool {
			position839, tokenIndex839, depth839 := position, tokenIndex, depth
			{
				position840 := position
				depth++
				{
					position841 := position
					depth++
					if buffer[position] != rune('<') {
						goto l839
					}
					position++
					if buffer[position] != rune('=') {
						goto l839
					}
					position++
					depth--
					add(rulePegText, position841)
				}
				if !_rules[ruleAction66]() {
					goto l839
				}
				depth--
				add(ruleLessOrEqual, position840)
			}
			return true
		l839:
			position, tokenIndex, depth = position839, tokenIndex839, depth839
			return false
		},
		/* 85 Greater <- <(<'>'> Action67)> */
		func() bool {
			position842, tokenIndex842, depth842 := position, tokenIndex, depth
			{
				position843 := position
				depth++
				{
					position844 := position
					depth++
					if buffer[position] != rune('>') {
						goto l842
					}
					position++
					depth--
					add(rulePegText, position844)
				}
				if !_rules[ruleAction67]() {
					goto l842
				}
				depth--
				add(ruleGreater, position843)
			}
			return true
		l842:
			position, tokenIndex, depth = position842, tokenIndex842, depth842
			return false
		},
		/* 86 GreaterOrEqual <- <(<('>' '=')> Action68)> */
		func() bool {
			position845, tokenIndex845, depth845 := position, tokenIndex, depth
			{
				position846 := position
				depth++
				{
					position847 := position
					depth++
					if buffer[position] != rune('>') {
						goto l845
					}
					position++
					if buffer[position] != rune('=') {
						goto l845
					}
					position++
					depth--
					add(rulePegText, position847)
				}
				if !_rules[ruleAction68]() {
					goto l845
				}
				depth--
				add(ruleGreaterOrEqual, position846)
			}
			return true
		l845:
			position, tokenIndex, depth = position845, tokenIndex845, depth845
			return false
		},
		/* 87 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action69)> */
		func() bool {
			position848, tokenIndex848, depth848 := position, tokenIndex, depth
			{
				position849 := position
				depth++
				{
					position850 := position
					depth++
					{
						position851, tokenIndex851, depth851 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l852
						}
						position++
						if buffer[position] != rune('=') {
							goto l852
						}
						position++
						goto l851
					l852:
						position, tokenIndex, depth = position851, tokenIndex851, depth851
						if buffer[position] != rune('<') {
							goto l848
						}
						position++
						if buffer[position] != rune('>') {
							goto l848
						}
						position++
					}
				l851:
					depth--
					add(rulePegText, position850)
				}
				if !_rules[ruleAction69]() {
					goto l848
				}
				depth--
				add(ruleNotEqual, position849)
			}
			return true
		l848:
			position, tokenIndex, depth = position848, tokenIndex848, depth848
			return false
		},
		/* 88 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action70)> */
		func() bool {
			position853, tokenIndex853, depth853 := position, tokenIndex, depth
			{
				position854 := position
				depth++
				{
					position855 := position
					depth++
					{
						position856, tokenIndex856, depth856 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l857
						}
						position++
						goto l856
					l857:
						position, tokenIndex, depth = position856, tokenIndex856, depth856
						if buffer[position] != rune('I') {
							goto l853
						}
						position++
					}
				l856:
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l859
						}
						position++
						goto l858
					l859:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if buffer[position] != rune('S') {
							goto l853
						}
						position++
					}
				l858:
					depth--
					add(rulePegText, position855)
				}
				if !_rules[ruleAction70]() {
					goto l853
				}
				depth--
				add(ruleIs, position854)
			}
			return true
		l853:
			position, tokenIndex, depth = position853, tokenIndex853, depth853
			return false
		},
		/* 89 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action71)> */
		func() bool {
			position860, tokenIndex860, depth860 := position, tokenIndex, depth
			{
				position861 := position
				depth++
				{
					position862 := position
					depth++
					{
						position863, tokenIndex863, depth863 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l864
						}
						position++
						goto l863
					l864:
						position, tokenIndex, depth = position863, tokenIndex863, depth863
						if buffer[position] != rune('I') {
							goto l860
						}
						position++
					}
				l863:
					{
						position865, tokenIndex865, depth865 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l866
						}
						position++
						goto l865
					l866:
						position, tokenIndex, depth = position865, tokenIndex865, depth865
						if buffer[position] != rune('S') {
							goto l860
						}
						position++
					}
				l865:
					if !_rules[rulesp]() {
						goto l860
					}
					{
						position867, tokenIndex867, depth867 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l868
						}
						position++
						goto l867
					l868:
						position, tokenIndex, depth = position867, tokenIndex867, depth867
						if buffer[position] != rune('N') {
							goto l860
						}
						position++
					}
				l867:
					{
						position869, tokenIndex869, depth869 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l870
						}
						position++
						goto l869
					l870:
						position, tokenIndex, depth = position869, tokenIndex869, depth869
						if buffer[position] != rune('O') {
							goto l860
						}
						position++
					}
				l869:
					{
						position871, tokenIndex871, depth871 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l872
						}
						position++
						goto l871
					l872:
						position, tokenIndex, depth = position871, tokenIndex871, depth871
						if buffer[position] != rune('T') {
							goto l860
						}
						position++
					}
				l871:
					depth--
					add(rulePegText, position862)
				}
				if !_rules[ruleAction71]() {
					goto l860
				}
				depth--
				add(ruleIsNot, position861)
			}
			return true
		l860:
			position, tokenIndex, depth = position860, tokenIndex860, depth860
			return false
		},
		/* 90 Plus <- <(<'+'> Action72)> */
		func() bool {
			position873, tokenIndex873, depth873 := position, tokenIndex, depth
			{
				position874 := position
				depth++
				{
					position875 := position
					depth++
					if buffer[position] != rune('+') {
						goto l873
					}
					position++
					depth--
					add(rulePegText, position875)
				}
				if !_rules[ruleAction72]() {
					goto l873
				}
				depth--
				add(rulePlus, position874)
			}
			return true
		l873:
			position, tokenIndex, depth = position873, tokenIndex873, depth873
			return false
		},
		/* 91 Minus <- <(<'-'> Action73)> */
		func() bool {
			position876, tokenIndex876, depth876 := position, tokenIndex, depth
			{
				position877 := position
				depth++
				{
					position878 := position
					depth++
					if buffer[position] != rune('-') {
						goto l876
					}
					position++
					depth--
					add(rulePegText, position878)
				}
				if !_rules[ruleAction73]() {
					goto l876
				}
				depth--
				add(ruleMinus, position877)
			}
			return true
		l876:
			position, tokenIndex, depth = position876, tokenIndex876, depth876
			return false
		},
		/* 92 Multiply <- <(<'*'> Action74)> */
		func() bool {
			position879, tokenIndex879, depth879 := position, tokenIndex, depth
			{
				position880 := position
				depth++
				{
					position881 := position
					depth++
					if buffer[position] != rune('*') {
						goto l879
					}
					position++
					depth--
					add(rulePegText, position881)
				}
				if !_rules[ruleAction74]() {
					goto l879
				}
				depth--
				add(ruleMultiply, position880)
			}
			return true
		l879:
			position, tokenIndex, depth = position879, tokenIndex879, depth879
			return false
		},
		/* 93 Divide <- <(<'/'> Action75)> */
		func() bool {
			position882, tokenIndex882, depth882 := position, tokenIndex, depth
			{
				position883 := position
				depth++
				{
					position884 := position
					depth++
					if buffer[position] != rune('/') {
						goto l882
					}
					position++
					depth--
					add(rulePegText, position884)
				}
				if !_rules[ruleAction75]() {
					goto l882
				}
				depth--
				add(ruleDivide, position883)
			}
			return true
		l882:
			position, tokenIndex, depth = position882, tokenIndex882, depth882
			return false
		},
		/* 94 Modulo <- <(<'%'> Action76)> */
		func() bool {
			position885, tokenIndex885, depth885 := position, tokenIndex, depth
			{
				position886 := position
				depth++
				{
					position887 := position
					depth++
					if buffer[position] != rune('%') {
						goto l885
					}
					position++
					depth--
					add(rulePegText, position887)
				}
				if !_rules[ruleAction76]() {
					goto l885
				}
				depth--
				add(ruleModulo, position886)
			}
			return true
		l885:
			position, tokenIndex, depth = position885, tokenIndex885, depth885
			return false
		},
		/* 95 UnaryMinus <- <(<'-'> Action77)> */
		func() bool {
			position888, tokenIndex888, depth888 := position, tokenIndex, depth
			{
				position889 := position
				depth++
				{
					position890 := position
					depth++
					if buffer[position] != rune('-') {
						goto l888
					}
					position++
					depth--
					add(rulePegText, position890)
				}
				if !_rules[ruleAction77]() {
					goto l888
				}
				depth--
				add(ruleUnaryMinus, position889)
			}
			return true
		l888:
			position, tokenIndex, depth = position888, tokenIndex888, depth888
			return false
		},
		/* 96 Identifier <- <(<ident> Action78)> */
		func() bool {
			position891, tokenIndex891, depth891 := position, tokenIndex, depth
			{
				position892 := position
				depth++
				{
					position893 := position
					depth++
					if !_rules[ruleident]() {
						goto l891
					}
					depth--
					add(rulePegText, position893)
				}
				if !_rules[ruleAction78]() {
					goto l891
				}
				depth--
				add(ruleIdentifier, position892)
			}
			return true
		l891:
			position, tokenIndex, depth = position891, tokenIndex891, depth891
			return false
		},
		/* 97 TargetIdentifier <- <(<jsonPath> Action79)> */
		func() bool {
			position894, tokenIndex894, depth894 := position, tokenIndex, depth
			{
				position895 := position
				depth++
				{
					position896 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l894
					}
					depth--
					add(rulePegText, position896)
				}
				if !_rules[ruleAction79]() {
					goto l894
				}
				depth--
				add(ruleTargetIdentifier, position895)
			}
			return true
		l894:
			position, tokenIndex, depth = position894, tokenIndex894, depth894
			return false
		},
		/* 98 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position897, tokenIndex897, depth897 := position, tokenIndex, depth
			{
				position898 := position
				depth++
				{
					position899, tokenIndex899, depth899 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l900
					}
					position++
					goto l899
				l900:
					position, tokenIndex, depth = position899, tokenIndex899, depth899
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l897
					}
					position++
				}
			l899:
			l901:
				{
					position902, tokenIndex902, depth902 := position, tokenIndex, depth
					{
						position903, tokenIndex903, depth903 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l904
						}
						position++
						goto l903
					l904:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l905
						}
						position++
						goto l903
					l905:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l906
						}
						position++
						goto l903
					l906:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if buffer[position] != rune('_') {
							goto l902
						}
						position++
					}
				l903:
					goto l901
				l902:
					position, tokenIndex, depth = position902, tokenIndex902, depth902
				}
				depth--
				add(ruleident, position898)
			}
			return true
		l897:
			position, tokenIndex, depth = position897, tokenIndex897, depth897
			return false
		},
		/* 99 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position907, tokenIndex907, depth907 := position, tokenIndex, depth
			{
				position908 := position
				depth++
				{
					position909, tokenIndex909, depth909 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l910
					}
					position++
					goto l909
				l910:
					position, tokenIndex, depth = position909, tokenIndex909, depth909
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l907
					}
					position++
				}
			l909:
			l911:
				{
					position912, tokenIndex912, depth912 := position, tokenIndex, depth
					{
						position913, tokenIndex913, depth913 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l914
						}
						position++
						goto l913
					l914:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l915
						}
						position++
						goto l913
					l915:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l916
						}
						position++
						goto l913
					l916:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('_') {
							goto l917
						}
						position++
						goto l913
					l917:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('.') {
							goto l918
						}
						position++
						goto l913
					l918:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('[') {
							goto l919
						}
						position++
						goto l913
					l919:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune(']') {
							goto l920
						}
						position++
						goto l913
					l920:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('"') {
							goto l912
						}
						position++
					}
				l913:
					goto l911
				l912:
					position, tokenIndex, depth = position912, tokenIndex912, depth912
				}
				depth--
				add(rulejsonPath, position908)
			}
			return true
		l907:
			position, tokenIndex, depth = position907, tokenIndex907, depth907
			return false
		},
		/* 100 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position922 := position
				depth++
			l923:
				{
					position924, tokenIndex924, depth924 := position, tokenIndex, depth
					{
						position925, tokenIndex925, depth925 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l926
						}
						position++
						goto l925
					l926:
						position, tokenIndex, depth = position925, tokenIndex925, depth925
						if buffer[position] != rune('\t') {
							goto l927
						}
						position++
						goto l925
					l927:
						position, tokenIndex, depth = position925, tokenIndex925, depth925
						if buffer[position] != rune('\n') {
							goto l928
						}
						position++
						goto l925
					l928:
						position, tokenIndex, depth = position925, tokenIndex925, depth925
						if buffer[position] != rune('\r') {
							goto l929
						}
						position++
						goto l925
					l929:
						position, tokenIndex, depth = position925, tokenIndex925, depth925
						if !_rules[rulecomment]() {
							goto l924
						}
					}
				l925:
					goto l923
				l924:
					position, tokenIndex, depth = position924, tokenIndex924, depth924
				}
				depth--
				add(rulesp, position922)
			}
			return true
		},
		/* 101 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position930, tokenIndex930, depth930 := position, tokenIndex, depth
			{
				position931 := position
				depth++
				if buffer[position] != rune('-') {
					goto l930
				}
				position++
				if buffer[position] != rune('-') {
					goto l930
				}
				position++
			l932:
				{
					position933, tokenIndex933, depth933 := position, tokenIndex, depth
					{
						position934, tokenIndex934, depth934 := position, tokenIndex, depth
						{
							position935, tokenIndex935, depth935 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l936
							}
							position++
							goto l935
						l936:
							position, tokenIndex, depth = position935, tokenIndex935, depth935
							if buffer[position] != rune('\n') {
								goto l934
							}
							position++
						}
					l935:
						goto l933
					l934:
						position, tokenIndex, depth = position934, tokenIndex934, depth934
					}
					if !matchDot() {
						goto l933
					}
					goto l932
				l933:
					position, tokenIndex, depth = position933, tokenIndex933, depth933
				}
				{
					position937, tokenIndex937, depth937 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l938
					}
					position++
					goto l937
				l938:
					position, tokenIndex, depth = position937, tokenIndex937, depth937
					if buffer[position] != rune('\n') {
						goto l930
					}
					position++
				}
			l937:
				depth--
				add(rulecomment, position931)
			}
			return true
		l930:
			position, tokenIndex, depth = position930, tokenIndex930, depth930
			return false
		},
		/* 103 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 104 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 105 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 106 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 107 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 108 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 109 Action6 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 110 Action7 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 111 Action8 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 112 Action9 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 113 Action10 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		nil,
		/* 115 Action11 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 116 Action12 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 117 Action13 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 118 Action14 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 119 Action15 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 120 Action16 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 121 Action17 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 122 Action18 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 123 Action19 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 124 Action20 <- <{
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
		/* 125 Action21 <- <{
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
		/* 126 Action22 <- <{
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
		/* 127 Action23 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 128 Action24 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 129 Action25 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 130 Action26 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 131 Action27 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 132 Action28 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 133 Action29 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 134 Action30 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 135 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 136 Action32 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 137 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 138 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 139 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 140 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 141 Action37 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 142 Action38 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 143 Action39 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 144 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 145 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 146 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 147 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 148 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 149 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 150 Action46 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 151 Action47 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 152 Action48 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 153 Action49 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 154 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 155 Action51 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 156 Action52 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 157 Action53 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 158 Action54 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 159 Action55 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 160 Action56 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 161 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 162 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 163 Action59 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 164 Action60 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 165 Action61 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 166 Action62 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 167 Action63 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 168 Action64 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 169 Action65 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 170 Action66 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 171 Action67 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 172 Action68 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 173 Action69 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 174 Action70 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 175 Action71 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 176 Action72 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 177 Action73 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 178 Action74 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 179 Action75 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 180 Action76 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 181 Action77 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 182 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 183 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
	}
	p.rules = _rules
}

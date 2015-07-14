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
	ruleSourceStmt
	ruleSinkStmt
	ruleStateStmt
	ruleStreamStmt
	ruleSelectStmt
	ruleCreateStreamAsSelectStmt
	ruleCreateSourceStmt
	ruleCreateSinkStmt
	ruleCreateStateStmt
	ruleUpdateStateStmt
	ruleUpdateSourceStmt
	ruleUpdateSinkStmt
	ruleInsertIntoSelectStmt
	ruleInsertIntoFromStmt
	rulePauseSourceStmt
	ruleResumeSourceStmt
	ruleRewindSourceStmt
	ruleDropSourceStmt
	ruleDropStreamStmt
	ruleDropSinkStmt
	ruleDropStateStmt
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
	ruleUpdateSourceSinkSpecs
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
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	rulePegText
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
	ruleAction80
	ruleAction81
	ruleAction82
	ruleAction83
	ruleAction84
	ruleAction85
	ruleAction86

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Statements",
	"Statement",
	"SourceStmt",
	"SinkStmt",
	"StateStmt",
	"StreamStmt",
	"SelectStmt",
	"CreateStreamAsSelectStmt",
	"CreateSourceStmt",
	"CreateSinkStmt",
	"CreateStateStmt",
	"UpdateStateStmt",
	"UpdateSourceStmt",
	"UpdateSinkStmt",
	"InsertIntoSelectStmt",
	"InsertIntoFromStmt",
	"PauseSourceStmt",
	"ResumeSourceStmt",
	"RewindSourceStmt",
	"DropSourceStmt",
	"DropStreamStmt",
	"DropSinkStmt",
	"DropStateStmt",
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
	"UpdateSourceSinkSpecs",
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
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"PegText",
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
	"Action80",
	"Action81",
	"Action82",
	"Action83",
	"Action84",
	"Action85",
	"Action86",

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
	rules  [202]func() bool
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

			p.AssembleUpdateState()

		case ruleAction6:

			p.AssembleUpdateSource()

		case ruleAction7:

			p.AssembleUpdateSink()

		case ruleAction8:

			p.AssembleInsertIntoSelect()

		case ruleAction9:

			p.AssembleInsertIntoFrom()

		case ruleAction10:

			p.AssemblePauseSource()

		case ruleAction11:

			p.AssembleResumeSource()

		case ruleAction12:

			p.AssembleRewindSource()

		case ruleAction13:

			p.AssembleDropSource()

		case ruleAction14:

			p.AssembleDropStream()

		case ruleAction15:

			p.AssembleDropSink()

		case ruleAction16:

			p.AssembleDropState()

		case ruleAction17:

			p.AssembleEmitter(begin, end)

		case ruleAction18:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction19:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction20:

			p.AssembleStreamEmitInterval()

		case ruleAction21:

			p.AssembleProjections(begin, end)

		case ruleAction22:

			p.AssembleAlias()

		case ruleAction23:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction24:

			p.AssembleInterval()

		case ruleAction25:

			p.AssembleInterval()

		case ruleAction26:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction27:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction28:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction29:

			p.EnsureAliasedStreamWindow()

		case ruleAction30:

			p.AssembleAliasedStreamWindow()

		case ruleAction31:

			p.AssembleStreamWindow()

		case ruleAction32:

			p.AssembleUDSFFuncApp()

		case ruleAction33:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction34:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction35:

			p.AssembleSourceSinkParam()

		case ruleAction36:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction38:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction39:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction40:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction41:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction42:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction43:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction44:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction45:

			p.AssembleFuncApp()

		case ruleAction46:

			p.AssembleExpressions(begin, end)

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction53:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction54:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction55:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction56:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction57:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction58:

			p.PushComponent(begin, end, Istream)

		case ruleAction59:

			p.PushComponent(begin, end, Dstream)

		case ruleAction60:

			p.PushComponent(begin, end, Rstream)

		case ruleAction61:

			p.PushComponent(begin, end, Tuples)

		case ruleAction62:

			p.PushComponent(begin, end, Seconds)

		case ruleAction63:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction64:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction65:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction66:

			p.PushComponent(begin, end, Yes)

		case ruleAction67:

			p.PushComponent(begin, end, No)

		case ruleAction68:

			p.PushComponent(begin, end, Or)

		case ruleAction69:

			p.PushComponent(begin, end, And)

		case ruleAction70:

			p.PushComponent(begin, end, Not)

		case ruleAction71:

			p.PushComponent(begin, end, Equal)

		case ruleAction72:

			p.PushComponent(begin, end, Less)

		case ruleAction73:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction74:

			p.PushComponent(begin, end, Greater)

		case ruleAction75:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction76:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction77:

			p.PushComponent(begin, end, Is)

		case ruleAction78:

			p.PushComponent(begin, end, IsNot)

		case ruleAction79:

			p.PushComponent(begin, end, Plus)

		case ruleAction80:

			p.PushComponent(begin, end, Minus)

		case ruleAction81:

			p.PushComponent(begin, end, Multiply)

		case ruleAction82:

			p.PushComponent(begin, end, Divide)

		case ruleAction83:

			p.PushComponent(begin, end, Modulo)

		case ruleAction84:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction85:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction86:

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
		/* 1 Statement <- <(SelectStmt / SourceStmt / SinkStmt / StateStmt / StreamStmt)> */
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
					if !_rules[ruleSourceStmt]() {
						goto l11
					}
					goto l9
				l11:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleSinkStmt]() {
						goto l12
					}
					goto l9
				l12:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleStateStmt]() {
						goto l13
					}
					goto l9
				l13:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleStreamStmt]() {
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
		/* 2 SourceStmt <- <(CreateSourceStmt / UpdateSourceStmt / DropSourceStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt)> */
		func() bool {
			position14, tokenIndex14, depth14 := position, tokenIndex, depth
			{
				position15 := position
				depth++
				{
					position16, tokenIndex16, depth16 := position, tokenIndex, depth
					if !_rules[ruleCreateSourceStmt]() {
						goto l17
					}
					goto l16
				l17:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
					if !_rules[ruleUpdateSourceStmt]() {
						goto l18
					}
					goto l16
				l18:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
					if !_rules[ruleDropSourceStmt]() {
						goto l19
					}
					goto l16
				l19:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
					if !_rules[rulePauseSourceStmt]() {
						goto l20
					}
					goto l16
				l20:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
					if !_rules[ruleResumeSourceStmt]() {
						goto l21
					}
					goto l16
				l21:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
					if !_rules[ruleRewindSourceStmt]() {
						goto l14
					}
				}
			l16:
				depth--
				add(ruleSourceStmt, position15)
			}
			return true
		l14:
			position, tokenIndex, depth = position14, tokenIndex14, depth14
			return false
		},
		/* 3 SinkStmt <- <(CreateSinkStmt / UpdateSinkStmt / DropSinkStmt)> */
		func() bool {
			position22, tokenIndex22, depth22 := position, tokenIndex, depth
			{
				position23 := position
				depth++
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if !_rules[ruleCreateSinkStmt]() {
						goto l25
					}
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleUpdateSinkStmt]() {
						goto l26
					}
					goto l24
				l26:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleDropSinkStmt]() {
						goto l22
					}
				}
			l24:
				depth--
				add(ruleSinkStmt, position23)
			}
			return true
		l22:
			position, tokenIndex, depth = position22, tokenIndex22, depth22
			return false
		},
		/* 4 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt)> */
		func() bool {
			position27, tokenIndex27, depth27 := position, tokenIndex, depth
			{
				position28 := position
				depth++
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if !_rules[ruleCreateStateStmt]() {
						goto l30
					}
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if !_rules[ruleUpdateStateStmt]() {
						goto l31
					}
					goto l29
				l31:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if !_rules[ruleDropStateStmt]() {
						goto l27
					}
				}
			l29:
				depth--
				add(ruleStateStmt, position28)
			}
			return true
		l27:
			position, tokenIndex, depth = position27, tokenIndex27, depth27
			return false
		},
		/* 5 StreamStmt <- <(CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoSelectStmt / InsertIntoFromStmt)> */
		func() bool {
			position32, tokenIndex32, depth32 := position, tokenIndex, depth
			{
				position33 := position
				depth++
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l35
					}
					goto l34
				l35:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleDropStreamStmt]() {
						goto l36
					}
					goto l34
				l36:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l37
					}
					goto l34
				l37:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l32
					}
				}
			l34:
				depth--
				add(ruleStreamStmt, position33)
			}
			return true
		l32:
			position, tokenIndex, depth = position32, tokenIndex32, depth32
			return false
		},
		/* 6 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
		func() bool {
			position38, tokenIndex38, depth38 := position, tokenIndex, depth
			{
				position39 := position
				depth++
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('S') {
						goto l38
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
						goto l38
					}
					position++
				}
			l42:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l45
					}
					position++
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if buffer[position] != rune('L') {
						goto l38
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
						goto l38
					}
					position++
				}
			l46:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('C') {
						goto l38
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
						goto l38
					}
					position++
				}
			l50:
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleEmitter]() {
					goto l38
				}
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleProjections]() {
					goto l38
				}
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleWindowedFrom]() {
					goto l38
				}
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleFilter]() {
					goto l38
				}
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleGrouping]() {
					goto l38
				}
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleHaving]() {
					goto l38
				}
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleAction0]() {
					goto l38
				}
				depth--
				add(ruleSelectStmt, position39)
			}
			return true
		l38:
			position, tokenIndex, depth = position38, tokenIndex38, depth38
			return false
		},
		/* 7 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position52, tokenIndex52, depth52 := position, tokenIndex, depth
			{
				position53 := position
				depth++
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('C') {
						goto l52
					}
					position++
				}
			l54:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('R') {
						goto l52
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('E') {
						goto l52
					}
					position++
				}
			l58:
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
						goto l52
					}
					position++
				}
			l60:
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
					if buffer[position] != rune('T') {
						goto l52
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
						goto l52
					}
					position++
				}
			l64:
				if !_rules[rulesp]() {
					goto l52
				}
				{
					position66, tokenIndex66, depth66 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex, depth = position66, tokenIndex66, depth66
					if buffer[position] != rune('S') {
						goto l52
					}
					position++
				}
			l66:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
					if buffer[position] != rune('T') {
						goto l52
					}
					position++
				}
			l68:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l71
					}
					position++
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if buffer[position] != rune('R') {
						goto l52
					}
					position++
				}
			l70:
				{
					position72, tokenIndex72, depth72 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l73
					}
					position++
					goto l72
				l73:
					position, tokenIndex, depth = position72, tokenIndex72, depth72
					if buffer[position] != rune('E') {
						goto l52
					}
					position++
				}
			l72:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if buffer[position] != rune('A') {
						goto l52
					}
					position++
				}
			l74:
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('M') {
						goto l52
					}
					position++
				}
			l76:
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l52
				}
				if !_rules[rulesp]() {
					goto l52
				}
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l79
					}
					position++
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('A') {
						goto l52
					}
					position++
				}
			l78:
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('S') {
						goto l52
					}
					position++
				}
			l80:
				if !_rules[rulesp]() {
					goto l52
				}
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l83
					}
					position++
					goto l82
				l83:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
					if buffer[position] != rune('S') {
						goto l52
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
						goto l52
					}
					position++
				}
			l84:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if buffer[position] != rune('L') {
						goto l52
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
						goto l52
					}
					position++
				}
			l88:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('C') {
						goto l52
					}
					position++
				}
			l90:
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('T') {
						goto l52
					}
					position++
				}
			l92:
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleEmitter]() {
					goto l52
				}
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleProjections]() {
					goto l52
				}
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleWindowedFrom]() {
					goto l52
				}
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleFilter]() {
					goto l52
				}
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleGrouping]() {
					goto l52
				}
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleHaving]() {
					goto l52
				}
				if !_rules[rulesp]() {
					goto l52
				}
				if !_rules[ruleAction1]() {
					goto l52
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position53)
			}
			return true
		l52:
			position, tokenIndex, depth = position52, tokenIndex52, depth52
			return false
		},
		/* 8 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position94, tokenIndex94, depth94 := position, tokenIndex, depth
			{
				position95 := position
				depth++
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
						goto l94
					}
					position++
				}
			l96:
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('R') {
						goto l94
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
						goto l94
					}
					position++
				}
			l100:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('A') {
						goto l94
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('T') {
						goto l94
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
						goto l94
					}
					position++
				}
			l106:
				if !_rules[rulesp]() {
					goto l94
				}
				if !_rules[rulePausedOpt]() {
					goto l94
				}
				if !_rules[rulesp]() {
					goto l94
				}
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('S') {
						goto l94
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('O') {
						goto l94
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('U') {
						goto l94
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
						goto l94
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('C') {
						goto l94
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
						goto l94
					}
					position++
				}
			l118:
				if !_rules[rulesp]() {
					goto l94
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l94
				}
				if !_rules[rulesp]() {
					goto l94
				}
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
						goto l94
					}
					position++
				}
			l120:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('Y') {
						goto l94
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('P') {
						goto l94
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('E') {
						goto l94
					}
					position++
				}
			l126:
				if !_rules[rulesp]() {
					goto l94
				}
				if !_rules[ruleSourceSinkType]() {
					goto l94
				}
				if !_rules[rulesp]() {
					goto l94
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l94
				}
				if !_rules[ruleAction2]() {
					goto l94
				}
				depth--
				add(ruleCreateSourceStmt, position95)
			}
			return true
		l94:
			position, tokenIndex, depth = position94, tokenIndex94, depth94
			return false
		},
		/* 9 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position128, tokenIndex128, depth128 := position, tokenIndex, depth
			{
				position129 := position
				depth++
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('C') {
						goto l128
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('R') {
						goto l128
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
						goto l128
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('A') {
						goto l128
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('T') {
						goto l128
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('E') {
						goto l128
					}
					position++
				}
			l140:
				if !_rules[rulesp]() {
					goto l128
				}
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('S') {
						goto l128
					}
					position++
				}
			l142:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('I') {
						goto l128
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('N') {
						goto l128
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('K') {
						goto l128
					}
					position++
				}
			l148:
				if !_rules[rulesp]() {
					goto l128
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l128
				}
				if !_rules[rulesp]() {
					goto l128
				}
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
						goto l128
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('Y') {
						goto l128
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('P') {
						goto l128
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
						goto l128
					}
					position++
				}
			l156:
				if !_rules[rulesp]() {
					goto l128
				}
				if !_rules[ruleSourceSinkType]() {
					goto l128
				}
				if !_rules[rulesp]() {
					goto l128
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l128
				}
				if !_rules[ruleAction3]() {
					goto l128
				}
				depth--
				add(ruleCreateSinkStmt, position129)
			}
			return true
		l128:
			position, tokenIndex, depth = position128, tokenIndex128, depth128
			return false
		},
		/* 10 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position158, tokenIndex158, depth158 := position, tokenIndex, depth
			{
				position159 := position
				depth++
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('C') {
						goto l158
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('R') {
						goto l158
					}
					position++
				}
			l162:
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('E') {
						goto l158
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('A') {
						goto l158
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('T') {
						goto l158
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
						goto l158
					}
					position++
				}
			l170:
				if !_rules[rulesp]() {
					goto l158
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
						goto l158
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('T') {
						goto l158
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('A') {
						goto l158
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('T') {
						goto l158
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
						goto l158
					}
					position++
				}
			l180:
				if !_rules[rulesp]() {
					goto l158
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l158
				}
				if !_rules[rulesp]() {
					goto l158
				}
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('T') {
						goto l158
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('Y') {
						goto l158
					}
					position++
				}
			l184:
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('P') {
						goto l158
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
						goto l158
					}
					position++
				}
			l188:
				if !_rules[rulesp]() {
					goto l158
				}
				if !_rules[ruleSourceSinkType]() {
					goto l158
				}
				if !_rules[rulesp]() {
					goto l158
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l158
				}
				if !_rules[ruleAction4]() {
					goto l158
				}
				depth--
				add(ruleCreateStateStmt, position159)
			}
			return true
		l158:
			position, tokenIndex, depth = position158, tokenIndex158, depth158
			return false
		},
		/* 11 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action5)> */
		func() bool {
			position190, tokenIndex190, depth190 := position, tokenIndex, depth
			{
				position191 := position
				depth++
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('U') {
						goto l190
					}
					position++
				}
			l192:
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('P') {
						goto l190
					}
					position++
				}
			l194:
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('D') {
						goto l190
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
						goto l190
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
						goto l190
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
						goto l190
					}
					position++
				}
			l202:
				if !_rules[rulesp]() {
					goto l190
				}
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('S') {
						goto l190
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
						goto l190
					}
					position++
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('A') {
						goto l190
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('T') {
						goto l190
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('E') {
						goto l190
					}
					position++
				}
			l212:
				if !_rules[rulesp]() {
					goto l190
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l190
				}
				if !_rules[rulesp]() {
					goto l190
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l190
				}
				if !_rules[ruleAction5]() {
					goto l190
				}
				depth--
				add(ruleUpdateStateStmt, position191)
			}
			return true
		l190:
			position, tokenIndex, depth = position190, tokenIndex190, depth190
			return false
		},
		/* 12 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action6)> */
		func() bool {
			position214, tokenIndex214, depth214 := position, tokenIndex, depth
			{
				position215 := position
				depth++
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('U') {
						goto l214
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('P') {
						goto l214
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('D') {
						goto l214
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('A') {
						goto l214
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('T') {
						goto l214
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
						goto l214
					}
					position++
				}
			l226:
				if !_rules[rulesp]() {
					goto l214
				}
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('S') {
						goto l214
					}
					position++
				}
			l228:
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
					if buffer[position] != rune('O') {
						goto l214
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('U') {
						goto l214
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('R') {
						goto l214
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('C') {
						goto l214
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
						goto l214
					}
					position++
				}
			l238:
				if !_rules[rulesp]() {
					goto l214
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l214
				}
				if !_rules[rulesp]() {
					goto l214
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l214
				}
				if !_rules[ruleAction6]() {
					goto l214
				}
				depth--
				add(ruleUpdateSourceStmt, position215)
			}
			return true
		l214:
			position, tokenIndex, depth = position214, tokenIndex214, depth214
			return false
		},
		/* 13 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action7)> */
		func() bool {
			position240, tokenIndex240, depth240 := position, tokenIndex, depth
			{
				position241 := position
				depth++
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('U') {
						goto l240
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('P') {
						goto l240
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('D') {
						goto l240
					}
					position++
				}
			l246:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('A') {
						goto l240
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('T') {
						goto l240
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
						goto l240
					}
					position++
				}
			l252:
				if !_rules[rulesp]() {
					goto l240
				}
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
						goto l240
					}
					position++
				}
			l254:
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l257
					}
					position++
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if buffer[position] != rune('I') {
						goto l240
					}
					position++
				}
			l256:
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if buffer[position] != rune('N') {
						goto l240
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('K') {
						goto l240
					}
					position++
				}
			l260:
				if !_rules[rulesp]() {
					goto l240
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l240
				}
				if !_rules[rulesp]() {
					goto l240
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l240
				}
				if !_rules[ruleAction7]() {
					goto l240
				}
				depth--
				add(ruleUpdateSinkStmt, position241)
			}
			return true
		l240:
			position, tokenIndex, depth = position240, tokenIndex240, depth240
			return false
		},
		/* 14 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action8)> */
		func() bool {
			position262, tokenIndex262, depth262 := position, tokenIndex, depth
			{
				position263 := position
				depth++
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('I') {
						goto l262
					}
					position++
				}
			l264:
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('N') {
						goto l262
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('S') {
						goto l262
					}
					position++
				}
			l268:
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if buffer[position] != rune('E') {
						goto l262
					}
					position++
				}
			l270:
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l273
					}
					position++
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if buffer[position] != rune('R') {
						goto l262
					}
					position++
				}
			l272:
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l275
					}
					position++
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('T') {
						goto l262
					}
					position++
				}
			l274:
				if !_rules[rulesp]() {
					goto l262
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
						goto l262
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
						goto l262
					}
					position++
				}
			l278:
				{
					position280, tokenIndex280, depth280 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l281
					}
					position++
					goto l280
				l281:
					position, tokenIndex, depth = position280, tokenIndex280, depth280
					if buffer[position] != rune('T') {
						goto l262
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
						goto l262
					}
					position++
				}
			l282:
				if !_rules[rulesp]() {
					goto l262
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l262
				}
				if !_rules[rulesp]() {
					goto l262
				}
				if !_rules[ruleSelectStmt]() {
					goto l262
				}
				if !_rules[ruleAction8]() {
					goto l262
				}
				depth--
				add(ruleInsertIntoSelectStmt, position263)
			}
			return true
		l262:
			position, tokenIndex, depth = position262, tokenIndex262, depth262
			return false
		},
		/* 15 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action9)> */
		func() bool {
			position284, tokenIndex284, depth284 := position, tokenIndex, depth
			{
				position285 := position
				depth++
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l287
					}
					position++
					goto l286
				l287:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
					if buffer[position] != rune('I') {
						goto l284
					}
					position++
				}
			l286:
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l289
					}
					position++
					goto l288
				l289:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
					if buffer[position] != rune('N') {
						goto l284
					}
					position++
				}
			l288:
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l291
					}
					position++
					goto l290
				l291:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
					if buffer[position] != rune('S') {
						goto l284
					}
					position++
				}
			l290:
				{
					position292, tokenIndex292, depth292 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l293
					}
					position++
					goto l292
				l293:
					position, tokenIndex, depth = position292, tokenIndex292, depth292
					if buffer[position] != rune('E') {
						goto l284
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
						goto l284
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if buffer[position] != rune('T') {
						goto l284
					}
					position++
				}
			l296:
				if !_rules[rulesp]() {
					goto l284
				}
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					if buffer[position] != rune('I') {
						goto l284
					}
					position++
				}
			l298:
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l301
					}
					position++
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if buffer[position] != rune('N') {
						goto l284
					}
					position++
				}
			l300:
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l303
					}
					position++
					goto l302
				l303:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					if buffer[position] != rune('T') {
						goto l284
					}
					position++
				}
			l302:
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
					if buffer[position] != rune('O') {
						goto l284
					}
					position++
				}
			l304:
				if !_rules[rulesp]() {
					goto l284
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l284
				}
				if !_rules[rulesp]() {
					goto l284
				}
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l307
					}
					position++
					goto l306
				l307:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					if buffer[position] != rune('F') {
						goto l284
					}
					position++
				}
			l306:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					if buffer[position] != rune('R') {
						goto l284
					}
					position++
				}
			l308:
				{
					position310, tokenIndex310, depth310 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l311
					}
					position++
					goto l310
				l311:
					position, tokenIndex, depth = position310, tokenIndex310, depth310
					if buffer[position] != rune('O') {
						goto l284
					}
					position++
				}
			l310:
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('M') {
						goto l284
					}
					position++
				}
			l312:
				if !_rules[rulesp]() {
					goto l284
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l284
				}
				if !_rules[ruleAction9]() {
					goto l284
				}
				depth--
				add(ruleInsertIntoFromStmt, position285)
			}
			return true
		l284:
			position, tokenIndex, depth = position284, tokenIndex284, depth284
			return false
		},
		/* 16 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position314, tokenIndex314, depth314 := position, tokenIndex, depth
			{
				position315 := position
				depth++
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
						goto l314
					}
					position++
				}
			l316:
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
					if buffer[position] != rune('A') {
						goto l314
					}
					position++
				}
			l318:
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('U') {
						goto l314
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if buffer[position] != rune('S') {
						goto l314
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
						goto l314
					}
					position++
				}
			l324:
				if !_rules[rulesp]() {
					goto l314
				}
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l327
					}
					position++
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					if buffer[position] != rune('S') {
						goto l314
					}
					position++
				}
			l326:
				{
					position328, tokenIndex328, depth328 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l329
					}
					position++
					goto l328
				l329:
					position, tokenIndex, depth = position328, tokenIndex328, depth328
					if buffer[position] != rune('O') {
						goto l314
					}
					position++
				}
			l328:
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l331
					}
					position++
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					if buffer[position] != rune('U') {
						goto l314
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
						goto l314
					}
					position++
				}
			l332:
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if buffer[position] != rune('C') {
						goto l314
					}
					position++
				}
			l334:
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if buffer[position] != rune('E') {
						goto l314
					}
					position++
				}
			l336:
				if !_rules[rulesp]() {
					goto l314
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l314
				}
				if !_rules[ruleAction10]() {
					goto l314
				}
				depth--
				add(rulePauseSourceStmt, position315)
			}
			return true
		l314:
			position, tokenIndex, depth = position314, tokenIndex314, depth314
			return false
		},
		/* 17 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action11)> */
		func() bool {
			position338, tokenIndex338, depth338 := position, tokenIndex, depth
			{
				position339 := position
				depth++
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l341
					}
					position++
					goto l340
				l341:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
					if buffer[position] != rune('R') {
						goto l338
					}
					position++
				}
			l340:
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if buffer[position] != rune('E') {
						goto l338
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
					if buffer[position] != rune('S') {
						goto l338
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					if buffer[position] != rune('U') {
						goto l338
					}
					position++
				}
			l346:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if buffer[position] != rune('M') {
						goto l338
					}
					position++
				}
			l348:
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if buffer[position] != rune('E') {
						goto l338
					}
					position++
				}
			l350:
				if !_rules[rulesp]() {
					goto l338
				}
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
						goto l338
					}
					position++
				}
			l352:
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
					if buffer[position] != rune('O') {
						goto l338
					}
					position++
				}
			l354:
				{
					position356, tokenIndex356, depth356 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l357
					}
					position++
					goto l356
				l357:
					position, tokenIndex, depth = position356, tokenIndex356, depth356
					if buffer[position] != rune('U') {
						goto l338
					}
					position++
				}
			l356:
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l359
					}
					position++
					goto l358
				l359:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					if buffer[position] != rune('R') {
						goto l338
					}
					position++
				}
			l358:
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if buffer[position] != rune('C') {
						goto l338
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
						goto l338
					}
					position++
				}
			l362:
				if !_rules[rulesp]() {
					goto l338
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l338
				}
				if !_rules[ruleAction11]() {
					goto l338
				}
				depth--
				add(ruleResumeSourceStmt, position339)
			}
			return true
		l338:
			position, tokenIndex, depth = position338, tokenIndex338, depth338
			return false
		},
		/* 18 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action12)> */
		func() bool {
			position364, tokenIndex364, depth364 := position, tokenIndex, depth
			{
				position365 := position
				depth++
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l367
					}
					position++
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if buffer[position] != rune('R') {
						goto l364
					}
					position++
				}
			l366:
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if buffer[position] != rune('E') {
						goto l364
					}
					position++
				}
			l368:
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l371
					}
					position++
					goto l370
				l371:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
					if buffer[position] != rune('W') {
						goto l364
					}
					position++
				}
			l370:
				{
					position372, tokenIndex372, depth372 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l373
					}
					position++
					goto l372
				l373:
					position, tokenIndex, depth = position372, tokenIndex372, depth372
					if buffer[position] != rune('I') {
						goto l364
					}
					position++
				}
			l372:
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('N') {
						goto l364
					}
					position++
				}
			l374:
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l377
					}
					position++
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					if buffer[position] != rune('D') {
						goto l364
					}
					position++
				}
			l376:
				if !_rules[rulesp]() {
					goto l364
				}
				{
					position378, tokenIndex378, depth378 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l379
					}
					position++
					goto l378
				l379:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
					if buffer[position] != rune('S') {
						goto l364
					}
					position++
				}
			l378:
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('O') {
						goto l364
					}
					position++
				}
			l380:
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l383
					}
					position++
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if buffer[position] != rune('U') {
						goto l364
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
						goto l364
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('C') {
						goto l364
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
						goto l364
					}
					position++
				}
			l388:
				if !_rules[rulesp]() {
					goto l364
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l364
				}
				if !_rules[ruleAction12]() {
					goto l364
				}
				depth--
				add(ruleRewindSourceStmt, position365)
			}
			return true
		l364:
			position, tokenIndex, depth = position364, tokenIndex364, depth364
			return false
		},
		/* 19 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if buffer[position] != rune('D') {
						goto l390
					}
					position++
				}
			l392:
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l395
					}
					position++
					goto l394
				l395:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
					if buffer[position] != rune('R') {
						goto l390
					}
					position++
				}
			l394:
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('O') {
						goto l390
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if buffer[position] != rune('P') {
						goto l390
					}
					position++
				}
			l398:
				if !_rules[rulesp]() {
					goto l390
				}
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l401
					}
					position++
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					if buffer[position] != rune('S') {
						goto l390
					}
					position++
				}
			l400:
				{
					position402, tokenIndex402, depth402 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l403
					}
					position++
					goto l402
				l403:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
					if buffer[position] != rune('O') {
						goto l390
					}
					position++
				}
			l402:
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l405
					}
					position++
					goto l404
				l405:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					if buffer[position] != rune('U') {
						goto l390
					}
					position++
				}
			l404:
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
						goto l390
					}
					position++
				}
			l406:
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l409
					}
					position++
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if buffer[position] != rune('C') {
						goto l390
					}
					position++
				}
			l408:
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l411
					}
					position++
					goto l410
				l411:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
					if buffer[position] != rune('E') {
						goto l390
					}
					position++
				}
			l410:
				if !_rules[rulesp]() {
					goto l390
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l390
				}
				if !_rules[ruleAction13]() {
					goto l390
				}
				depth--
				add(ruleDropSourceStmt, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 20 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action14)> */
		func() bool {
			position412, tokenIndex412, depth412 := position, tokenIndex, depth
			{
				position413 := position
				depth++
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('D') {
						goto l412
					}
					position++
				}
			l414:
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l417
					}
					position++
					goto l416
				l417:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if buffer[position] != rune('R') {
						goto l412
					}
					position++
				}
			l416:
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
					if buffer[position] != rune('O') {
						goto l412
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('P') {
						goto l412
					}
					position++
				}
			l420:
				if !_rules[rulesp]() {
					goto l412
				}
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
						goto l412
					}
					position++
				}
			l422:
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l425
					}
					position++
					goto l424
				l425:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
					if buffer[position] != rune('T') {
						goto l412
					}
					position++
				}
			l424:
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
						goto l412
					}
					position++
				}
			l426:
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l429
					}
					position++
					goto l428
				l429:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
					if buffer[position] != rune('E') {
						goto l412
					}
					position++
				}
			l428:
				{
					position430, tokenIndex430, depth430 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l431
					}
					position++
					goto l430
				l431:
					position, tokenIndex, depth = position430, tokenIndex430, depth430
					if buffer[position] != rune('A') {
						goto l412
					}
					position++
				}
			l430:
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('M') {
						goto l412
					}
					position++
				}
			l432:
				if !_rules[rulesp]() {
					goto l412
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l412
				}
				if !_rules[ruleAction14]() {
					goto l412
				}
				depth--
				add(ruleDropStreamStmt, position413)
			}
			return true
		l412:
			position, tokenIndex, depth = position412, tokenIndex412, depth412
			return false
		},
		/* 21 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action15)> */
		func() bool {
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
			{
				position435 := position
				depth++
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l437
					}
					position++
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if buffer[position] != rune('D') {
						goto l434
					}
					position++
				}
			l436:
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if buffer[position] != rune('R') {
						goto l434
					}
					position++
				}
			l438:
				{
					position440, tokenIndex440, depth440 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l441
					}
					position++
					goto l440
				l441:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if buffer[position] != rune('O') {
						goto l434
					}
					position++
				}
			l440:
				{
					position442, tokenIndex442, depth442 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l443
					}
					position++
					goto l442
				l443:
					position, tokenIndex, depth = position442, tokenIndex442, depth442
					if buffer[position] != rune('P') {
						goto l434
					}
					position++
				}
			l442:
				if !_rules[rulesp]() {
					goto l434
				}
				{
					position444, tokenIndex444, depth444 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l445
					}
					position++
					goto l444
				l445:
					position, tokenIndex, depth = position444, tokenIndex444, depth444
					if buffer[position] != rune('S') {
						goto l434
					}
					position++
				}
			l444:
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
					if buffer[position] != rune('I') {
						goto l434
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
						goto l434
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('K') {
						goto l434
					}
					position++
				}
			l450:
				if !_rules[rulesp]() {
					goto l434
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l434
				}
				if !_rules[ruleAction15]() {
					goto l434
				}
				depth--
				add(ruleDropSinkStmt, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
			return false
		},
		/* 22 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position452, tokenIndex452, depth452 := position, tokenIndex, depth
			{
				position453 := position
				depth++
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l455
					}
					position++
					goto l454
				l455:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if buffer[position] != rune('D') {
						goto l452
					}
					position++
				}
			l454:
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l457
					}
					position++
					goto l456
				l457:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if buffer[position] != rune('R') {
						goto l452
					}
					position++
				}
			l456:
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l459
					}
					position++
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if buffer[position] != rune('O') {
						goto l452
					}
					position++
				}
			l458:
				{
					position460, tokenIndex460, depth460 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l461
					}
					position++
					goto l460
				l461:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
					if buffer[position] != rune('P') {
						goto l452
					}
					position++
				}
			l460:
				if !_rules[rulesp]() {
					goto l452
				}
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l463
					}
					position++
					goto l462
				l463:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if buffer[position] != rune('S') {
						goto l452
					}
					position++
				}
			l462:
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l465
					}
					position++
					goto l464
				l465:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					if buffer[position] != rune('T') {
						goto l452
					}
					position++
				}
			l464:
				{
					position466, tokenIndex466, depth466 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l467
					}
					position++
					goto l466
				l467:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if buffer[position] != rune('A') {
						goto l452
					}
					position++
				}
			l466:
				{
					position468, tokenIndex468, depth468 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l469
					}
					position++
					goto l468
				l469:
					position, tokenIndex, depth = position468, tokenIndex468, depth468
					if buffer[position] != rune('T') {
						goto l452
					}
					position++
				}
			l468:
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l471
					}
					position++
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if buffer[position] != rune('E') {
						goto l452
					}
					position++
				}
			l470:
				if !_rules[rulesp]() {
					goto l452
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l452
				}
				if !_rules[ruleAction16]() {
					goto l452
				}
				depth--
				add(ruleDropStateStmt, position453)
			}
			return true
		l452:
			position, tokenIndex, depth = position452, tokenIndex452, depth452
			return false
		},
		/* 23 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action17)> */
		func() bool {
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				{
					position474, tokenIndex474, depth474 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l475
					}
					goto l474
				l475:
					position, tokenIndex, depth = position474, tokenIndex474, depth474
					if !_rules[ruleDSTREAM]() {
						goto l476
					}
					goto l474
				l476:
					position, tokenIndex, depth = position474, tokenIndex474, depth474
					if !_rules[ruleRSTREAM]() {
						goto l472
					}
				}
			l474:
				{
					position477 := position
					depth++
					{
						position478, tokenIndex478, depth478 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l478
						}
						if buffer[position] != rune('[') {
							goto l478
						}
						position++
						if !_rules[rulesp]() {
							goto l478
						}
						{
							position480, tokenIndex480, depth480 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l481
							}
							position++
							goto l480
						l481:
							position, tokenIndex, depth = position480, tokenIndex480, depth480
							if buffer[position] != rune('E') {
								goto l478
							}
							position++
						}
					l480:
						{
							position482, tokenIndex482, depth482 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l483
							}
							position++
							goto l482
						l483:
							position, tokenIndex, depth = position482, tokenIndex482, depth482
							if buffer[position] != rune('V') {
								goto l478
							}
							position++
						}
					l482:
						{
							position484, tokenIndex484, depth484 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l485
							}
							position++
							goto l484
						l485:
							position, tokenIndex, depth = position484, tokenIndex484, depth484
							if buffer[position] != rune('E') {
								goto l478
							}
							position++
						}
					l484:
						{
							position486, tokenIndex486, depth486 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l487
							}
							position++
							goto l486
						l487:
							position, tokenIndex, depth = position486, tokenIndex486, depth486
							if buffer[position] != rune('R') {
								goto l478
							}
							position++
						}
					l486:
						{
							position488, tokenIndex488, depth488 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l489
							}
							position++
							goto l488
						l489:
							position, tokenIndex, depth = position488, tokenIndex488, depth488
							if buffer[position] != rune('Y') {
								goto l478
							}
							position++
						}
					l488:
						if !_rules[rulesp]() {
							goto l478
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l478
						}
						if !_rules[rulesp]() {
							goto l478
						}
						if buffer[position] != rune(']') {
							goto l478
						}
						position++
						goto l479
					l478:
						position, tokenIndex, depth = position478, tokenIndex478, depth478
					}
				l479:
					depth--
					add(rulePegText, position477)
				}
				if !_rules[ruleAction17]() {
					goto l472
				}
				depth--
				add(ruleEmitter, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 24 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position490, tokenIndex490, depth490 := position, tokenIndex, depth
			{
				position491 := position
				depth++
				{
					position492, tokenIndex492, depth492 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l493
					}
				l494:
					{
						position495, tokenIndex495, depth495 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l495
						}
						if buffer[position] != rune(',') {
							goto l495
						}
						position++
						if !_rules[rulesp]() {
							goto l495
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l495
						}
						goto l494
					l495:
						position, tokenIndex, depth = position495, tokenIndex495, depth495
					}
					goto l492
				l493:
					position, tokenIndex, depth = position492, tokenIndex492, depth492
					if !_rules[ruleTimeEmitterInterval]() {
						goto l496
					}
					goto l492
				l496:
					position, tokenIndex, depth = position492, tokenIndex492, depth492
					if !_rules[ruleTupleEmitterInterval]() {
						goto l490
					}
				}
			l492:
				depth--
				add(ruleEmitterIntervals, position491)
			}
			return true
		l490:
			position, tokenIndex, depth = position490, tokenIndex490, depth490
			return false
		},
		/* 25 TimeEmitterInterval <- <(<TimeInterval> Action18)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				{
					position499 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l497
					}
					depth--
					add(rulePegText, position499)
				}
				if !_rules[ruleAction18]() {
					goto l497
				}
				depth--
				add(ruleTimeEmitterInterval, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 26 TupleEmitterInterval <- <(<TuplesInterval> Action19)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				{
					position502 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l500
					}
					depth--
					add(rulePegText, position502)
				}
				if !_rules[ruleAction19]() {
					goto l500
				}
				depth--
				add(ruleTupleEmitterInterval, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 27 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action20)> */
		func() bool {
			position503, tokenIndex503, depth503 := position, tokenIndex, depth
			{
				position504 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l503
				}
				if !_rules[rulesp]() {
					goto l503
				}
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l506
					}
					position++
					goto l505
				l506:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
					if buffer[position] != rune('I') {
						goto l503
					}
					position++
				}
			l505:
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('N') {
						goto l503
					}
					position++
				}
			l507:
				if !_rules[rulesp]() {
					goto l503
				}
				if !_rules[ruleStream]() {
					goto l503
				}
				if !_rules[ruleAction20]() {
					goto l503
				}
				depth--
				add(ruleTupleEmitterFromInterval, position504)
			}
			return true
		l503:
			position, tokenIndex, depth = position503, tokenIndex503, depth503
			return false
		},
		/* 28 Projections <- <(<(Projection sp (',' sp Projection)*)> Action21)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				{
					position511 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l509
					}
					if !_rules[rulesp]() {
						goto l509
					}
				l512:
					{
						position513, tokenIndex513, depth513 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l513
						}
						position++
						if !_rules[rulesp]() {
							goto l513
						}
						if !_rules[ruleProjection]() {
							goto l513
						}
						goto l512
					l513:
						position, tokenIndex, depth = position513, tokenIndex513, depth513
					}
					depth--
					add(rulePegText, position511)
				}
				if !_rules[ruleAction21]() {
					goto l509
				}
				depth--
				add(ruleProjections, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 29 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position514, tokenIndex514, depth514 := position, tokenIndex, depth
			{
				position515 := position
				depth++
				{
					position516, tokenIndex516, depth516 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l517
					}
					goto l516
				l517:
					position, tokenIndex, depth = position516, tokenIndex516, depth516
					if !_rules[ruleExpression]() {
						goto l518
					}
					goto l516
				l518:
					position, tokenIndex, depth = position516, tokenIndex516, depth516
					if !_rules[ruleWildcard]() {
						goto l514
					}
				}
			l516:
				depth--
				add(ruleProjection, position515)
			}
			return true
		l514:
			position, tokenIndex, depth = position514, tokenIndex514, depth514
			return false
		},
		/* 30 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action22)> */
		func() bool {
			position519, tokenIndex519, depth519 := position, tokenIndex, depth
			{
				position520 := position
				depth++
				{
					position521, tokenIndex521, depth521 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l522
					}
					goto l521
				l522:
					position, tokenIndex, depth = position521, tokenIndex521, depth521
					if !_rules[ruleWildcard]() {
						goto l519
					}
				}
			l521:
				if !_rules[rulesp]() {
					goto l519
				}
				{
					position523, tokenIndex523, depth523 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l524
					}
					position++
					goto l523
				l524:
					position, tokenIndex, depth = position523, tokenIndex523, depth523
					if buffer[position] != rune('A') {
						goto l519
					}
					position++
				}
			l523:
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l526
					}
					position++
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if buffer[position] != rune('S') {
						goto l519
					}
					position++
				}
			l525:
				if !_rules[rulesp]() {
					goto l519
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l519
				}
				if !_rules[ruleAction22]() {
					goto l519
				}
				depth--
				add(ruleAliasExpression, position520)
			}
			return true
		l519:
			position, tokenIndex, depth = position519, tokenIndex519, depth519
			return false
		},
		/* 31 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action23)> */
		func() bool {
			position527, tokenIndex527, depth527 := position, tokenIndex, depth
			{
				position528 := position
				depth++
				{
					position529 := position
					depth++
					{
						position530, tokenIndex530, depth530 := position, tokenIndex, depth
						{
							position532, tokenIndex532, depth532 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l533
							}
							position++
							goto l532
						l533:
							position, tokenIndex, depth = position532, tokenIndex532, depth532
							if buffer[position] != rune('F') {
								goto l530
							}
							position++
						}
					l532:
						{
							position534, tokenIndex534, depth534 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l535
							}
							position++
							goto l534
						l535:
							position, tokenIndex, depth = position534, tokenIndex534, depth534
							if buffer[position] != rune('R') {
								goto l530
							}
							position++
						}
					l534:
						{
							position536, tokenIndex536, depth536 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l537
							}
							position++
							goto l536
						l537:
							position, tokenIndex, depth = position536, tokenIndex536, depth536
							if buffer[position] != rune('O') {
								goto l530
							}
							position++
						}
					l536:
						{
							position538, tokenIndex538, depth538 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l539
							}
							position++
							goto l538
						l539:
							position, tokenIndex, depth = position538, tokenIndex538, depth538
							if buffer[position] != rune('M') {
								goto l530
							}
							position++
						}
					l538:
						if !_rules[rulesp]() {
							goto l530
						}
						if !_rules[ruleRelations]() {
							goto l530
						}
						if !_rules[rulesp]() {
							goto l530
						}
						goto l531
					l530:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
					}
				l531:
					depth--
					add(rulePegText, position529)
				}
				if !_rules[ruleAction23]() {
					goto l527
				}
				depth--
				add(ruleWindowedFrom, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 32 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542, tokenIndex542, depth542 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l543
					}
					goto l542
				l543:
					position, tokenIndex, depth = position542, tokenIndex542, depth542
					if !_rules[ruleTuplesInterval]() {
						goto l540
					}
				}
			l542:
				depth--
				add(ruleInterval, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 33 TimeInterval <- <(NumericLiteral sp SECONDS Action24)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l544
				}
				if !_rules[rulesp]() {
					goto l544
				}
				if !_rules[ruleSECONDS]() {
					goto l544
				}
				if !_rules[ruleAction24]() {
					goto l544
				}
				depth--
				add(ruleTimeInterval, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 34 TuplesInterval <- <(NumericLiteral sp TUPLES Action25)> */
		func() bool {
			position546, tokenIndex546, depth546 := position, tokenIndex, depth
			{
				position547 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l546
				}
				if !_rules[rulesp]() {
					goto l546
				}
				if !_rules[ruleTUPLES]() {
					goto l546
				}
				if !_rules[ruleAction25]() {
					goto l546
				}
				depth--
				add(ruleTuplesInterval, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 35 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position548, tokenIndex548, depth548 := position, tokenIndex, depth
			{
				position549 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l548
				}
				if !_rules[rulesp]() {
					goto l548
				}
			l550:
				{
					position551, tokenIndex551, depth551 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l551
					}
					position++
					if !_rules[rulesp]() {
						goto l551
					}
					if !_rules[ruleRelationLike]() {
						goto l551
					}
					goto l550
				l551:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
				}
				depth--
				add(ruleRelations, position549)
			}
			return true
		l548:
			position, tokenIndex, depth = position548, tokenIndex548, depth548
			return false
		},
		/* 36 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action26)> */
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
						{
							position557, tokenIndex557, depth557 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l558
							}
							position++
							goto l557
						l558:
							position, tokenIndex, depth = position557, tokenIndex557, depth557
							if buffer[position] != rune('W') {
								goto l555
							}
							position++
						}
					l557:
						{
							position559, tokenIndex559, depth559 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l560
							}
							position++
							goto l559
						l560:
							position, tokenIndex, depth = position559, tokenIndex559, depth559
							if buffer[position] != rune('H') {
								goto l555
							}
							position++
						}
					l559:
						{
							position561, tokenIndex561, depth561 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l562
							}
							position++
							goto l561
						l562:
							position, tokenIndex, depth = position561, tokenIndex561, depth561
							if buffer[position] != rune('E') {
								goto l555
							}
							position++
						}
					l561:
						{
							position563, tokenIndex563, depth563 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l564
							}
							position++
							goto l563
						l564:
							position, tokenIndex, depth = position563, tokenIndex563, depth563
							if buffer[position] != rune('R') {
								goto l555
							}
							position++
						}
					l563:
						{
							position565, tokenIndex565, depth565 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l566
							}
							position++
							goto l565
						l566:
							position, tokenIndex, depth = position565, tokenIndex565, depth565
							if buffer[position] != rune('E') {
								goto l555
							}
							position++
						}
					l565:
						if !_rules[rulesp]() {
							goto l555
						}
						if !_rules[ruleExpression]() {
							goto l555
						}
						goto l556
					l555:
						position, tokenIndex, depth = position555, tokenIndex555, depth555
					}
				l556:
					depth--
					add(rulePegText, position554)
				}
				if !_rules[ruleAction26]() {
					goto l552
				}
				depth--
				add(ruleFilter, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 37 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action27)> */
		func() bool {
			position567, tokenIndex567, depth567 := position, tokenIndex, depth
			{
				position568 := position
				depth++
				{
					position569 := position
					depth++
					{
						position570, tokenIndex570, depth570 := position, tokenIndex, depth
						{
							position572, tokenIndex572, depth572 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l573
							}
							position++
							goto l572
						l573:
							position, tokenIndex, depth = position572, tokenIndex572, depth572
							if buffer[position] != rune('G') {
								goto l570
							}
							position++
						}
					l572:
						{
							position574, tokenIndex574, depth574 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l575
							}
							position++
							goto l574
						l575:
							position, tokenIndex, depth = position574, tokenIndex574, depth574
							if buffer[position] != rune('R') {
								goto l570
							}
							position++
						}
					l574:
						{
							position576, tokenIndex576, depth576 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l577
							}
							position++
							goto l576
						l577:
							position, tokenIndex, depth = position576, tokenIndex576, depth576
							if buffer[position] != rune('O') {
								goto l570
							}
							position++
						}
					l576:
						{
							position578, tokenIndex578, depth578 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l579
							}
							position++
							goto l578
						l579:
							position, tokenIndex, depth = position578, tokenIndex578, depth578
							if buffer[position] != rune('U') {
								goto l570
							}
							position++
						}
					l578:
						{
							position580, tokenIndex580, depth580 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l581
							}
							position++
							goto l580
						l581:
							position, tokenIndex, depth = position580, tokenIndex580, depth580
							if buffer[position] != rune('P') {
								goto l570
							}
							position++
						}
					l580:
						if !_rules[rulesp]() {
							goto l570
						}
						{
							position582, tokenIndex582, depth582 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l583
							}
							position++
							goto l582
						l583:
							position, tokenIndex, depth = position582, tokenIndex582, depth582
							if buffer[position] != rune('B') {
								goto l570
							}
							position++
						}
					l582:
						{
							position584, tokenIndex584, depth584 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l585
							}
							position++
							goto l584
						l585:
							position, tokenIndex, depth = position584, tokenIndex584, depth584
							if buffer[position] != rune('Y') {
								goto l570
							}
							position++
						}
					l584:
						if !_rules[rulesp]() {
							goto l570
						}
						if !_rules[ruleGroupList]() {
							goto l570
						}
						goto l571
					l570:
						position, tokenIndex, depth = position570, tokenIndex570, depth570
					}
				l571:
					depth--
					add(rulePegText, position569)
				}
				if !_rules[ruleAction27]() {
					goto l567
				}
				depth--
				add(ruleGrouping, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 38 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position586, tokenIndex586, depth586 := position, tokenIndex, depth
			{
				position587 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l586
				}
				if !_rules[rulesp]() {
					goto l586
				}
			l588:
				{
					position589, tokenIndex589, depth589 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l589
					}
					position++
					if !_rules[rulesp]() {
						goto l589
					}
					if !_rules[ruleExpression]() {
						goto l589
					}
					goto l588
				l589:
					position, tokenIndex, depth = position589, tokenIndex589, depth589
				}
				depth--
				add(ruleGroupList, position587)
			}
			return true
		l586:
			position, tokenIndex, depth = position586, tokenIndex586, depth586
			return false
		},
		/* 39 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action28)> */
		func() bool {
			position590, tokenIndex590, depth590 := position, tokenIndex, depth
			{
				position591 := position
				depth++
				{
					position592 := position
					depth++
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						{
							position595, tokenIndex595, depth595 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l596
							}
							position++
							goto l595
						l596:
							position, tokenIndex, depth = position595, tokenIndex595, depth595
							if buffer[position] != rune('H') {
								goto l593
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
								goto l593
							}
							position++
						}
					l597:
						{
							position599, tokenIndex599, depth599 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l600
							}
							position++
							goto l599
						l600:
							position, tokenIndex, depth = position599, tokenIndex599, depth599
							if buffer[position] != rune('V') {
								goto l593
							}
							position++
						}
					l599:
						{
							position601, tokenIndex601, depth601 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l602
							}
							position++
							goto l601
						l602:
							position, tokenIndex, depth = position601, tokenIndex601, depth601
							if buffer[position] != rune('I') {
								goto l593
							}
							position++
						}
					l601:
						{
							position603, tokenIndex603, depth603 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l604
							}
							position++
							goto l603
						l604:
							position, tokenIndex, depth = position603, tokenIndex603, depth603
							if buffer[position] != rune('N') {
								goto l593
							}
							position++
						}
					l603:
						{
							position605, tokenIndex605, depth605 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l606
							}
							position++
							goto l605
						l606:
							position, tokenIndex, depth = position605, tokenIndex605, depth605
							if buffer[position] != rune('G') {
								goto l593
							}
							position++
						}
					l605:
						if !_rules[rulesp]() {
							goto l593
						}
						if !_rules[ruleExpression]() {
							goto l593
						}
						goto l594
					l593:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
					}
				l594:
					depth--
					add(rulePegText, position592)
				}
				if !_rules[ruleAction28]() {
					goto l590
				}
				depth--
				add(ruleHaving, position591)
			}
			return true
		l590:
			position, tokenIndex, depth = position590, tokenIndex590, depth590
			return false
		},
		/* 40 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action29))> */
		func() bool {
			position607, tokenIndex607, depth607 := position, tokenIndex, depth
			{
				position608 := position
				depth++
				{
					position609, tokenIndex609, depth609 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l610
					}
					goto l609
				l610:
					position, tokenIndex, depth = position609, tokenIndex609, depth609
					if !_rules[ruleStreamWindow]() {
						goto l607
					}
					if !_rules[ruleAction29]() {
						goto l607
					}
				}
			l609:
				depth--
				add(ruleRelationLike, position608)
			}
			return true
		l607:
			position, tokenIndex, depth = position607, tokenIndex607, depth607
			return false
		},
		/* 41 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action30)> */
		func() bool {
			position611, tokenIndex611, depth611 := position, tokenIndex, depth
			{
				position612 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l611
				}
				if !_rules[rulesp]() {
					goto l611
				}
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
						goto l611
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
						goto l611
					}
					position++
				}
			l615:
				if !_rules[rulesp]() {
					goto l611
				}
				if !_rules[ruleIdentifier]() {
					goto l611
				}
				if !_rules[ruleAction30]() {
					goto l611
				}
				depth--
				add(ruleAliasedStreamWindow, position612)
			}
			return true
		l611:
			position, tokenIndex, depth = position611, tokenIndex611, depth611
			return false
		},
		/* 42 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action31)> */
		func() bool {
			position617, tokenIndex617, depth617 := position, tokenIndex, depth
			{
				position618 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l617
				}
				if !_rules[rulesp]() {
					goto l617
				}
				if buffer[position] != rune('[') {
					goto l617
				}
				position++
				if !_rules[rulesp]() {
					goto l617
				}
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
						goto l617
					}
					position++
				}
			l619:
				{
					position621, tokenIndex621, depth621 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l622
					}
					position++
					goto l621
				l622:
					position, tokenIndex, depth = position621, tokenIndex621, depth621
					if buffer[position] != rune('A') {
						goto l617
					}
					position++
				}
			l621:
				{
					position623, tokenIndex623, depth623 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l624
					}
					position++
					goto l623
				l624:
					position, tokenIndex, depth = position623, tokenIndex623, depth623
					if buffer[position] != rune('N') {
						goto l617
					}
					position++
				}
			l623:
				{
					position625, tokenIndex625, depth625 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l626
					}
					position++
					goto l625
				l626:
					position, tokenIndex, depth = position625, tokenIndex625, depth625
					if buffer[position] != rune('G') {
						goto l617
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
						goto l617
					}
					position++
				}
			l627:
				if !_rules[rulesp]() {
					goto l617
				}
				if !_rules[ruleInterval]() {
					goto l617
				}
				if !_rules[rulesp]() {
					goto l617
				}
				if buffer[position] != rune(']') {
					goto l617
				}
				position++
				if !_rules[ruleAction31]() {
					goto l617
				}
				depth--
				add(ruleStreamWindow, position618)
			}
			return true
		l617:
			position, tokenIndex, depth = position617, tokenIndex617, depth617
			return false
		},
		/* 43 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position629, tokenIndex629, depth629 := position, tokenIndex, depth
			{
				position630 := position
				depth++
				{
					position631, tokenIndex631, depth631 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l632
					}
					goto l631
				l632:
					position, tokenIndex, depth = position631, tokenIndex631, depth631
					if !_rules[ruleStream]() {
						goto l629
					}
				}
			l631:
				depth--
				add(ruleStreamLike, position630)
			}
			return true
		l629:
			position, tokenIndex, depth = position629, tokenIndex629, depth629
			return false
		},
		/* 44 UDSFFuncApp <- <(FuncApp Action32)> */
		func() bool {
			position633, tokenIndex633, depth633 := position, tokenIndex, depth
			{
				position634 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l633
				}
				if !_rules[ruleAction32]() {
					goto l633
				}
				depth--
				add(ruleUDSFFuncApp, position634)
			}
			return true
		l633:
			position, tokenIndex, depth = position633, tokenIndex633, depth633
			return false
		},
		/* 45 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action33)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				{
					position637 := position
					depth++
					{
						position638, tokenIndex638, depth638 := position, tokenIndex, depth
						{
							position640, tokenIndex640, depth640 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l641
							}
							position++
							goto l640
						l641:
							position, tokenIndex, depth = position640, tokenIndex640, depth640
							if buffer[position] != rune('W') {
								goto l638
							}
							position++
						}
					l640:
						{
							position642, tokenIndex642, depth642 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l643
							}
							position++
							goto l642
						l643:
							position, tokenIndex, depth = position642, tokenIndex642, depth642
							if buffer[position] != rune('I') {
								goto l638
							}
							position++
						}
					l642:
						{
							position644, tokenIndex644, depth644 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l645
							}
							position++
							goto l644
						l645:
							position, tokenIndex, depth = position644, tokenIndex644, depth644
							if buffer[position] != rune('T') {
								goto l638
							}
							position++
						}
					l644:
						{
							position646, tokenIndex646, depth646 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l647
							}
							position++
							goto l646
						l647:
							position, tokenIndex, depth = position646, tokenIndex646, depth646
							if buffer[position] != rune('H') {
								goto l638
							}
							position++
						}
					l646:
						if !_rules[rulesp]() {
							goto l638
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l638
						}
						if !_rules[rulesp]() {
							goto l638
						}
					l648:
						{
							position649, tokenIndex649, depth649 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l649
							}
							position++
							if !_rules[rulesp]() {
								goto l649
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l649
							}
							goto l648
						l649:
							position, tokenIndex, depth = position649, tokenIndex649, depth649
						}
						goto l639
					l638:
						position, tokenIndex, depth = position638, tokenIndex638, depth638
					}
				l639:
					depth--
					add(rulePegText, position637)
				}
				if !_rules[ruleAction33]() {
					goto l635
				}
				depth--
				add(ruleSourceSinkSpecs, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 46 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action34)> */
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
						if buffer[position] != rune('s') {
							goto l654
						}
						position++
						goto l653
					l654:
						position, tokenIndex, depth = position653, tokenIndex653, depth653
						if buffer[position] != rune('S') {
							goto l650
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
					if !_rules[rulesp]() {
						goto l650
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l650
					}
					if !_rules[rulesp]() {
						goto l650
					}
				l659:
					{
						position660, tokenIndex660, depth660 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l660
						}
						position++
						if !_rules[rulesp]() {
							goto l660
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l660
						}
						goto l659
					l660:
						position, tokenIndex, depth = position660, tokenIndex660, depth660
					}
					depth--
					add(rulePegText, position652)
				}
				if !_rules[ruleAction34]() {
					goto l650
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position651)
			}
			return true
		l650:
			position, tokenIndex, depth = position650, tokenIndex650, depth650
			return false
		},
		/* 47 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action35)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l661
				}
				if buffer[position] != rune('=') {
					goto l661
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l661
				}
				if !_rules[ruleAction35]() {
					goto l661
				}
				depth--
				add(ruleSourceSinkParam, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 48 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position663, tokenIndex663, depth663 := position, tokenIndex, depth
			{
				position664 := position
				depth++
				{
					position665, tokenIndex665, depth665 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l666
					}
					goto l665
				l666:
					position, tokenIndex, depth = position665, tokenIndex665, depth665
					if !_rules[ruleLiteral]() {
						goto l663
					}
				}
			l665:
				depth--
				add(ruleSourceSinkParamVal, position664)
			}
			return true
		l663:
			position, tokenIndex, depth = position663, tokenIndex663, depth663
			return false
		},
		/* 49 PausedOpt <- <(<(Paused / Unpaused)?> Action36)> */
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
						{
							position672, tokenIndex672, depth672 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l673
							}
							goto l672
						l673:
							position, tokenIndex, depth = position672, tokenIndex672, depth672
							if !_rules[ruleUnpaused]() {
								goto l670
							}
						}
					l672:
						goto l671
					l670:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
					}
				l671:
					depth--
					add(rulePegText, position669)
				}
				if !_rules[ruleAction36]() {
					goto l667
				}
				depth--
				add(rulePausedOpt, position668)
			}
			return true
		l667:
			position, tokenIndex, depth = position667, tokenIndex667, depth667
			return false
		},
		/* 50 Expression <- <orExpr> */
		func() bool {
			position674, tokenIndex674, depth674 := position, tokenIndex, depth
			{
				position675 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l674
				}
				depth--
				add(ruleExpression, position675)
			}
			return true
		l674:
			position, tokenIndex, depth = position674, tokenIndex674, depth674
			return false
		},
		/* 51 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action37)> */
		func() bool {
			position676, tokenIndex676, depth676 := position, tokenIndex, depth
			{
				position677 := position
				depth++
				{
					position678 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l676
					}
					if !_rules[rulesp]() {
						goto l676
					}
					{
						position679, tokenIndex679, depth679 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l679
						}
						if !_rules[rulesp]() {
							goto l679
						}
						if !_rules[ruleandExpr]() {
							goto l679
						}
						goto l680
					l679:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
					}
				l680:
					depth--
					add(rulePegText, position678)
				}
				if !_rules[ruleAction37]() {
					goto l676
				}
				depth--
				add(ruleorExpr, position677)
			}
			return true
		l676:
			position, tokenIndex, depth = position676, tokenIndex676, depth676
			return false
		},
		/* 52 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action38)> */
		func() bool {
			position681, tokenIndex681, depth681 := position, tokenIndex, depth
			{
				position682 := position
				depth++
				{
					position683 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l681
					}
					if !_rules[rulesp]() {
						goto l681
					}
					{
						position684, tokenIndex684, depth684 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l684
						}
						if !_rules[rulesp]() {
							goto l684
						}
						if !_rules[rulenotExpr]() {
							goto l684
						}
						goto l685
					l684:
						position, tokenIndex, depth = position684, tokenIndex684, depth684
					}
				l685:
					depth--
					add(rulePegText, position683)
				}
				if !_rules[ruleAction38]() {
					goto l681
				}
				depth--
				add(ruleandExpr, position682)
			}
			return true
		l681:
			position, tokenIndex, depth = position681, tokenIndex681, depth681
			return false
		},
		/* 53 notExpr <- <(<((Not sp)? comparisonExpr)> Action39)> */
		func() bool {
			position686, tokenIndex686, depth686 := position, tokenIndex, depth
			{
				position687 := position
				depth++
				{
					position688 := position
					depth++
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l689
						}
						if !_rules[rulesp]() {
							goto l689
						}
						goto l690
					l689:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
					}
				l690:
					if !_rules[rulecomparisonExpr]() {
						goto l686
					}
					depth--
					add(rulePegText, position688)
				}
				if !_rules[ruleAction39]() {
					goto l686
				}
				depth--
				add(rulenotExpr, position687)
			}
			return true
		l686:
			position, tokenIndex, depth = position686, tokenIndex686, depth686
			return false
		},
		/* 54 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action40)> */
		func() bool {
			position691, tokenIndex691, depth691 := position, tokenIndex, depth
			{
				position692 := position
				depth++
				{
					position693 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l691
					}
					if !_rules[rulesp]() {
						goto l691
					}
					{
						position694, tokenIndex694, depth694 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l694
						}
						if !_rules[rulesp]() {
							goto l694
						}
						if !_rules[ruleisExpr]() {
							goto l694
						}
						goto l695
					l694:
						position, tokenIndex, depth = position694, tokenIndex694, depth694
					}
				l695:
					depth--
					add(rulePegText, position693)
				}
				if !_rules[ruleAction40]() {
					goto l691
				}
				depth--
				add(rulecomparisonExpr, position692)
			}
			return true
		l691:
			position, tokenIndex, depth = position691, tokenIndex691, depth691
			return false
		},
		/* 55 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action41)> */
		func() bool {
			position696, tokenIndex696, depth696 := position, tokenIndex, depth
			{
				position697 := position
				depth++
				{
					position698 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l696
					}
					if !_rules[rulesp]() {
						goto l696
					}
					{
						position699, tokenIndex699, depth699 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l699
						}
						if !_rules[rulesp]() {
							goto l699
						}
						if !_rules[ruleNullLiteral]() {
							goto l699
						}
						goto l700
					l699:
						position, tokenIndex, depth = position699, tokenIndex699, depth699
					}
				l700:
					depth--
					add(rulePegText, position698)
				}
				if !_rules[ruleAction41]() {
					goto l696
				}
				depth--
				add(ruleisExpr, position697)
			}
			return true
		l696:
			position, tokenIndex, depth = position696, tokenIndex696, depth696
			return false
		},
		/* 56 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action42)> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				{
					position703 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l701
					}
					if !_rules[rulesp]() {
						goto l701
					}
					{
						position704, tokenIndex704, depth704 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l704
						}
						if !_rules[rulesp]() {
							goto l704
						}
						if !_rules[ruleproductExpr]() {
							goto l704
						}
						goto l705
					l704:
						position, tokenIndex, depth = position704, tokenIndex704, depth704
					}
				l705:
					depth--
					add(rulePegText, position703)
				}
				if !_rules[ruleAction42]() {
					goto l701
				}
				depth--
				add(ruletermExpr, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 57 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action43)> */
		func() bool {
			position706, tokenIndex706, depth706 := position, tokenIndex, depth
			{
				position707 := position
				depth++
				{
					position708 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l706
					}
					if !_rules[rulesp]() {
						goto l706
					}
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l709
						}
						if !_rules[rulesp]() {
							goto l709
						}
						if !_rules[ruleminusExpr]() {
							goto l709
						}
						goto l710
					l709:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
					}
				l710:
					depth--
					add(rulePegText, position708)
				}
				if !_rules[ruleAction43]() {
					goto l706
				}
				depth--
				add(ruleproductExpr, position707)
			}
			return true
		l706:
			position, tokenIndex, depth = position706, tokenIndex706, depth706
			return false
		},
		/* 58 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action44)> */
		func() bool {
			position711, tokenIndex711, depth711 := position, tokenIndex, depth
			{
				position712 := position
				depth++
				{
					position713 := position
					depth++
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l714
						}
						if !_rules[rulesp]() {
							goto l714
						}
						goto l715
					l714:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
					}
				l715:
					if !_rules[rulebaseExpr]() {
						goto l711
					}
					depth--
					add(rulePegText, position713)
				}
				if !_rules[ruleAction44]() {
					goto l711
				}
				depth--
				add(ruleminusExpr, position712)
			}
			return true
		l711:
			position, tokenIndex, depth = position711, tokenIndex711, depth711
			return false
		},
		/* 59 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position716, tokenIndex716, depth716 := position, tokenIndex, depth
			{
				position717 := position
				depth++
				{
					position718, tokenIndex718, depth718 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l719
					}
					position++
					if !_rules[rulesp]() {
						goto l719
					}
					if !_rules[ruleExpression]() {
						goto l719
					}
					if !_rules[rulesp]() {
						goto l719
					}
					if buffer[position] != rune(')') {
						goto l719
					}
					position++
					goto l718
				l719:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if !_rules[ruleBooleanLiteral]() {
						goto l720
					}
					goto l718
				l720:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if !_rules[ruleNullLiteral]() {
						goto l721
					}
					goto l718
				l721:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if !_rules[ruleFuncApp]() {
						goto l722
					}
					goto l718
				l722:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if !_rules[ruleRowMeta]() {
						goto l723
					}
					goto l718
				l723:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if !_rules[ruleRowValue]() {
						goto l724
					}
					goto l718
				l724:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if !_rules[ruleLiteral]() {
						goto l716
					}
				}
			l718:
				depth--
				add(rulebaseExpr, position717)
			}
			return true
		l716:
			position, tokenIndex, depth = position716, tokenIndex716, depth716
			return false
		},
		/* 60 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action45)> */
		func() bool {
			position725, tokenIndex725, depth725 := position, tokenIndex, depth
			{
				position726 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l725
				}
				if !_rules[rulesp]() {
					goto l725
				}
				if buffer[position] != rune('(') {
					goto l725
				}
				position++
				if !_rules[rulesp]() {
					goto l725
				}
				if !_rules[ruleFuncParams]() {
					goto l725
				}
				if !_rules[rulesp]() {
					goto l725
				}
				if buffer[position] != rune(')') {
					goto l725
				}
				position++
				if !_rules[ruleAction45]() {
					goto l725
				}
				depth--
				add(ruleFuncApp, position726)
			}
			return true
		l725:
			position, tokenIndex, depth = position725, tokenIndex725, depth725
			return false
		},
		/* 61 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action46)> */
		func() bool {
			position727, tokenIndex727, depth727 := position, tokenIndex, depth
			{
				position728 := position
				depth++
				{
					position729 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l727
					}
					if !_rules[rulesp]() {
						goto l727
					}
				l730:
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l731
						}
						position++
						if !_rules[rulesp]() {
							goto l731
						}
						if !_rules[ruleExpression]() {
							goto l731
						}
						goto l730
					l731:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
					}
					depth--
					add(rulePegText, position729)
				}
				if !_rules[ruleAction46]() {
					goto l727
				}
				depth--
				add(ruleFuncParams, position728)
			}
			return true
		l727:
			position, tokenIndex, depth = position727, tokenIndex727, depth727
			return false
		},
		/* 62 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position732, tokenIndex732, depth732 := position, tokenIndex, depth
			{
				position733 := position
				depth++
				{
					position734, tokenIndex734, depth734 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l735
					}
					goto l734
				l735:
					position, tokenIndex, depth = position734, tokenIndex734, depth734
					if !_rules[ruleNumericLiteral]() {
						goto l736
					}
					goto l734
				l736:
					position, tokenIndex, depth = position734, tokenIndex734, depth734
					if !_rules[ruleStringLiteral]() {
						goto l732
					}
				}
			l734:
				depth--
				add(ruleLiteral, position733)
			}
			return true
		l732:
			position, tokenIndex, depth = position732, tokenIndex732, depth732
			return false
		},
		/* 63 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position737, tokenIndex737, depth737 := position, tokenIndex, depth
			{
				position738 := position
				depth++
				{
					position739, tokenIndex739, depth739 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l740
					}
					goto l739
				l740:
					position, tokenIndex, depth = position739, tokenIndex739, depth739
					if !_rules[ruleNotEqual]() {
						goto l741
					}
					goto l739
				l741:
					position, tokenIndex, depth = position739, tokenIndex739, depth739
					if !_rules[ruleLessOrEqual]() {
						goto l742
					}
					goto l739
				l742:
					position, tokenIndex, depth = position739, tokenIndex739, depth739
					if !_rules[ruleLess]() {
						goto l743
					}
					goto l739
				l743:
					position, tokenIndex, depth = position739, tokenIndex739, depth739
					if !_rules[ruleGreaterOrEqual]() {
						goto l744
					}
					goto l739
				l744:
					position, tokenIndex, depth = position739, tokenIndex739, depth739
					if !_rules[ruleGreater]() {
						goto l745
					}
					goto l739
				l745:
					position, tokenIndex, depth = position739, tokenIndex739, depth739
					if !_rules[ruleNotEqual]() {
						goto l737
					}
				}
			l739:
				depth--
				add(ruleComparisonOp, position738)
			}
			return true
		l737:
			position, tokenIndex, depth = position737, tokenIndex737, depth737
			return false
		},
		/* 64 IsOp <- <(IsNot / Is)> */
		func() bool {
			position746, tokenIndex746, depth746 := position, tokenIndex, depth
			{
				position747 := position
				depth++
				{
					position748, tokenIndex748, depth748 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l749
					}
					goto l748
				l749:
					position, tokenIndex, depth = position748, tokenIndex748, depth748
					if !_rules[ruleIs]() {
						goto l746
					}
				}
			l748:
				depth--
				add(ruleIsOp, position747)
			}
			return true
		l746:
			position, tokenIndex, depth = position746, tokenIndex746, depth746
			return false
		},
		/* 65 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position750, tokenIndex750, depth750 := position, tokenIndex, depth
			{
				position751 := position
				depth++
				{
					position752, tokenIndex752, depth752 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l753
					}
					goto l752
				l753:
					position, tokenIndex, depth = position752, tokenIndex752, depth752
					if !_rules[ruleMinus]() {
						goto l750
					}
				}
			l752:
				depth--
				add(rulePlusMinusOp, position751)
			}
			return true
		l750:
			position, tokenIndex, depth = position750, tokenIndex750, depth750
			return false
		},
		/* 66 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position754, tokenIndex754, depth754 := position, tokenIndex, depth
			{
				position755 := position
				depth++
				{
					position756, tokenIndex756, depth756 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l757
					}
					goto l756
				l757:
					position, tokenIndex, depth = position756, tokenIndex756, depth756
					if !_rules[ruleDivide]() {
						goto l758
					}
					goto l756
				l758:
					position, tokenIndex, depth = position756, tokenIndex756, depth756
					if !_rules[ruleModulo]() {
						goto l754
					}
				}
			l756:
				depth--
				add(ruleMultDivOp, position755)
			}
			return true
		l754:
			position, tokenIndex, depth = position754, tokenIndex754, depth754
			return false
		},
		/* 67 Stream <- <(<ident> Action47)> */
		func() bool {
			position759, tokenIndex759, depth759 := position, tokenIndex, depth
			{
				position760 := position
				depth++
				{
					position761 := position
					depth++
					if !_rules[ruleident]() {
						goto l759
					}
					depth--
					add(rulePegText, position761)
				}
				if !_rules[ruleAction47]() {
					goto l759
				}
				depth--
				add(ruleStream, position760)
			}
			return true
		l759:
			position, tokenIndex, depth = position759, tokenIndex759, depth759
			return false
		},
		/* 68 RowMeta <- <RowTimestamp> */
		func() bool {
			position762, tokenIndex762, depth762 := position, tokenIndex, depth
			{
				position763 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l762
				}
				depth--
				add(ruleRowMeta, position763)
			}
			return true
		l762:
			position, tokenIndex, depth = position762, tokenIndex762, depth762
			return false
		},
		/* 69 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action48)> */
		func() bool {
			position764, tokenIndex764, depth764 := position, tokenIndex, depth
			{
				position765 := position
				depth++
				{
					position766 := position
					depth++
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l767
						}
						if buffer[position] != rune(':') {
							goto l767
						}
						position++
						goto l768
					l767:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
					}
				l768:
					if buffer[position] != rune('t') {
						goto l764
					}
					position++
					if buffer[position] != rune('s') {
						goto l764
					}
					position++
					if buffer[position] != rune('(') {
						goto l764
					}
					position++
					if buffer[position] != rune(')') {
						goto l764
					}
					position++
					depth--
					add(rulePegText, position766)
				}
				if !_rules[ruleAction48]() {
					goto l764
				}
				depth--
				add(ruleRowTimestamp, position765)
			}
			return true
		l764:
			position, tokenIndex, depth = position764, tokenIndex764, depth764
			return false
		},
		/* 70 RowValue <- <(<((ident ':')? jsonPath)> Action49)> */
		func() bool {
			position769, tokenIndex769, depth769 := position, tokenIndex, depth
			{
				position770 := position
				depth++
				{
					position771 := position
					depth++
					{
						position772, tokenIndex772, depth772 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l772
						}
						if buffer[position] != rune(':') {
							goto l772
						}
						position++
						goto l773
					l772:
						position, tokenIndex, depth = position772, tokenIndex772, depth772
					}
				l773:
					if !_rules[rulejsonPath]() {
						goto l769
					}
					depth--
					add(rulePegText, position771)
				}
				if !_rules[ruleAction49]() {
					goto l769
				}
				depth--
				add(ruleRowValue, position770)
			}
			return true
		l769:
			position, tokenIndex, depth = position769, tokenIndex769, depth769
			return false
		},
		/* 71 NumericLiteral <- <(<('-'? [0-9]+)> Action50)> */
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
						if buffer[position] != rune('-') {
							goto l777
						}
						position++
						goto l778
					l777:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
					}
				l778:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l774
					}
					position++
				l779:
					{
						position780, tokenIndex780, depth780 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l780
						}
						position++
						goto l779
					l780:
						position, tokenIndex, depth = position780, tokenIndex780, depth780
					}
					depth--
					add(rulePegText, position776)
				}
				if !_rules[ruleAction50]() {
					goto l774
				}
				depth--
				add(ruleNumericLiteral, position775)
			}
			return true
		l774:
			position, tokenIndex, depth = position774, tokenIndex774, depth774
			return false
		},
		/* 72 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action51)> */
		func() bool {
			position781, tokenIndex781, depth781 := position, tokenIndex, depth
			{
				position782 := position
				depth++
				{
					position783 := position
					depth++
					{
						position784, tokenIndex784, depth784 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l784
						}
						position++
						goto l785
					l784:
						position, tokenIndex, depth = position784, tokenIndex784, depth784
					}
				l785:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l781
					}
					position++
				l786:
					{
						position787, tokenIndex787, depth787 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l787
						}
						position++
						goto l786
					l787:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
					}
					if buffer[position] != rune('.') {
						goto l781
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l781
					}
					position++
				l788:
					{
						position789, tokenIndex789, depth789 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l789
						}
						position++
						goto l788
					l789:
						position, tokenIndex, depth = position789, tokenIndex789, depth789
					}
					depth--
					add(rulePegText, position783)
				}
				if !_rules[ruleAction51]() {
					goto l781
				}
				depth--
				add(ruleFloatLiteral, position782)
			}
			return true
		l781:
			position, tokenIndex, depth = position781, tokenIndex781, depth781
			return false
		},
		/* 73 Function <- <(<ident> Action52)> */
		func() bool {
			position790, tokenIndex790, depth790 := position, tokenIndex, depth
			{
				position791 := position
				depth++
				{
					position792 := position
					depth++
					if !_rules[ruleident]() {
						goto l790
					}
					depth--
					add(rulePegText, position792)
				}
				if !_rules[ruleAction52]() {
					goto l790
				}
				depth--
				add(ruleFunction, position791)
			}
			return true
		l790:
			position, tokenIndex, depth = position790, tokenIndex790, depth790
			return false
		},
		/* 74 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action53)> */
		func() bool {
			position793, tokenIndex793, depth793 := position, tokenIndex, depth
			{
				position794 := position
				depth++
				{
					position795 := position
					depth++
					{
						position796, tokenIndex796, depth796 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l797
						}
						position++
						goto l796
					l797:
						position, tokenIndex, depth = position796, tokenIndex796, depth796
						if buffer[position] != rune('N') {
							goto l793
						}
						position++
					}
				l796:
					{
						position798, tokenIndex798, depth798 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l799
						}
						position++
						goto l798
					l799:
						position, tokenIndex, depth = position798, tokenIndex798, depth798
						if buffer[position] != rune('U') {
							goto l793
						}
						position++
					}
				l798:
					{
						position800, tokenIndex800, depth800 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l801
						}
						position++
						goto l800
					l801:
						position, tokenIndex, depth = position800, tokenIndex800, depth800
						if buffer[position] != rune('L') {
							goto l793
						}
						position++
					}
				l800:
					{
						position802, tokenIndex802, depth802 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l803
						}
						position++
						goto l802
					l803:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if buffer[position] != rune('L') {
							goto l793
						}
						position++
					}
				l802:
					depth--
					add(rulePegText, position795)
				}
				if !_rules[ruleAction53]() {
					goto l793
				}
				depth--
				add(ruleNullLiteral, position794)
			}
			return true
		l793:
			position, tokenIndex, depth = position793, tokenIndex793, depth793
			return false
		},
		/* 75 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position804, tokenIndex804, depth804 := position, tokenIndex, depth
			{
				position805 := position
				depth++
				{
					position806, tokenIndex806, depth806 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l807
					}
					goto l806
				l807:
					position, tokenIndex, depth = position806, tokenIndex806, depth806
					if !_rules[ruleFALSE]() {
						goto l804
					}
				}
			l806:
				depth--
				add(ruleBooleanLiteral, position805)
			}
			return true
		l804:
			position, tokenIndex, depth = position804, tokenIndex804, depth804
			return false
		},
		/* 76 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action54)> */
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
						if buffer[position] != rune('t') {
							goto l812
						}
						position++
						goto l811
					l812:
						position, tokenIndex, depth = position811, tokenIndex811, depth811
						if buffer[position] != rune('T') {
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
					{
						position815, tokenIndex815, depth815 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l816
						}
						position++
						goto l815
					l816:
						position, tokenIndex, depth = position815, tokenIndex815, depth815
						if buffer[position] != rune('U') {
							goto l808
						}
						position++
					}
				l815:
					{
						position817, tokenIndex817, depth817 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l818
						}
						position++
						goto l817
					l818:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						if buffer[position] != rune('E') {
							goto l808
						}
						position++
					}
				l817:
					depth--
					add(rulePegText, position810)
				}
				if !_rules[ruleAction54]() {
					goto l808
				}
				depth--
				add(ruleTRUE, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 77 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action55)> */
		func() bool {
			position819, tokenIndex819, depth819 := position, tokenIndex, depth
			{
				position820 := position
				depth++
				{
					position821 := position
					depth++
					{
						position822, tokenIndex822, depth822 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l823
						}
						position++
						goto l822
					l823:
						position, tokenIndex, depth = position822, tokenIndex822, depth822
						if buffer[position] != rune('F') {
							goto l819
						}
						position++
					}
				l822:
					{
						position824, tokenIndex824, depth824 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l825
						}
						position++
						goto l824
					l825:
						position, tokenIndex, depth = position824, tokenIndex824, depth824
						if buffer[position] != rune('A') {
							goto l819
						}
						position++
					}
				l824:
					{
						position826, tokenIndex826, depth826 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l827
						}
						position++
						goto l826
					l827:
						position, tokenIndex, depth = position826, tokenIndex826, depth826
						if buffer[position] != rune('L') {
							goto l819
						}
						position++
					}
				l826:
					{
						position828, tokenIndex828, depth828 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l829
						}
						position++
						goto l828
					l829:
						position, tokenIndex, depth = position828, tokenIndex828, depth828
						if buffer[position] != rune('S') {
							goto l819
						}
						position++
					}
				l828:
					{
						position830, tokenIndex830, depth830 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l831
						}
						position++
						goto l830
					l831:
						position, tokenIndex, depth = position830, tokenIndex830, depth830
						if buffer[position] != rune('E') {
							goto l819
						}
						position++
					}
				l830:
					depth--
					add(rulePegText, position821)
				}
				if !_rules[ruleAction55]() {
					goto l819
				}
				depth--
				add(ruleFALSE, position820)
			}
			return true
		l819:
			position, tokenIndex, depth = position819, tokenIndex819, depth819
			return false
		},
		/* 78 Wildcard <- <(<'*'> Action56)> */
		func() bool {
			position832, tokenIndex832, depth832 := position, tokenIndex, depth
			{
				position833 := position
				depth++
				{
					position834 := position
					depth++
					if buffer[position] != rune('*') {
						goto l832
					}
					position++
					depth--
					add(rulePegText, position834)
				}
				if !_rules[ruleAction56]() {
					goto l832
				}
				depth--
				add(ruleWildcard, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 79 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action57)> */
		func() bool {
			position835, tokenIndex835, depth835 := position, tokenIndex, depth
			{
				position836 := position
				depth++
				{
					position837 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l835
					}
					position++
				l838:
					{
						position839, tokenIndex839, depth839 := position, tokenIndex, depth
						{
							position840, tokenIndex840, depth840 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l841
							}
							position++
							if buffer[position] != rune('\'') {
								goto l841
							}
							position++
							goto l840
						l841:
							position, tokenIndex, depth = position840, tokenIndex840, depth840
							{
								position842, tokenIndex842, depth842 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l842
								}
								position++
								goto l839
							l842:
								position, tokenIndex, depth = position842, tokenIndex842, depth842
							}
							if !matchDot() {
								goto l839
							}
						}
					l840:
						goto l838
					l839:
						position, tokenIndex, depth = position839, tokenIndex839, depth839
					}
					if buffer[position] != rune('\'') {
						goto l835
					}
					position++
					depth--
					add(rulePegText, position837)
				}
				if !_rules[ruleAction57]() {
					goto l835
				}
				depth--
				add(ruleStringLiteral, position836)
			}
			return true
		l835:
			position, tokenIndex, depth = position835, tokenIndex835, depth835
			return false
		},
		/* 80 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action58)> */
		func() bool {
			position843, tokenIndex843, depth843 := position, tokenIndex, depth
			{
				position844 := position
				depth++
				{
					position845 := position
					depth++
					{
						position846, tokenIndex846, depth846 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l847
						}
						position++
						goto l846
					l847:
						position, tokenIndex, depth = position846, tokenIndex846, depth846
						if buffer[position] != rune('I') {
							goto l843
						}
						position++
					}
				l846:
					{
						position848, tokenIndex848, depth848 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l849
						}
						position++
						goto l848
					l849:
						position, tokenIndex, depth = position848, tokenIndex848, depth848
						if buffer[position] != rune('S') {
							goto l843
						}
						position++
					}
				l848:
					{
						position850, tokenIndex850, depth850 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l851
						}
						position++
						goto l850
					l851:
						position, tokenIndex, depth = position850, tokenIndex850, depth850
						if buffer[position] != rune('T') {
							goto l843
						}
						position++
					}
				l850:
					{
						position852, tokenIndex852, depth852 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l853
						}
						position++
						goto l852
					l853:
						position, tokenIndex, depth = position852, tokenIndex852, depth852
						if buffer[position] != rune('R') {
							goto l843
						}
						position++
					}
				l852:
					{
						position854, tokenIndex854, depth854 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l855
						}
						position++
						goto l854
					l855:
						position, tokenIndex, depth = position854, tokenIndex854, depth854
						if buffer[position] != rune('E') {
							goto l843
						}
						position++
					}
				l854:
					{
						position856, tokenIndex856, depth856 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l857
						}
						position++
						goto l856
					l857:
						position, tokenIndex, depth = position856, tokenIndex856, depth856
						if buffer[position] != rune('A') {
							goto l843
						}
						position++
					}
				l856:
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l859
						}
						position++
						goto l858
					l859:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if buffer[position] != rune('M') {
							goto l843
						}
						position++
					}
				l858:
					depth--
					add(rulePegText, position845)
				}
				if !_rules[ruleAction58]() {
					goto l843
				}
				depth--
				add(ruleISTREAM, position844)
			}
			return true
		l843:
			position, tokenIndex, depth = position843, tokenIndex843, depth843
			return false
		},
		/* 81 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action59)> */
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
						if buffer[position] != rune('d') {
							goto l864
						}
						position++
						goto l863
					l864:
						position, tokenIndex, depth = position863, tokenIndex863, depth863
						if buffer[position] != rune('D') {
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
					{
						position867, tokenIndex867, depth867 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l868
						}
						position++
						goto l867
					l868:
						position, tokenIndex, depth = position867, tokenIndex867, depth867
						if buffer[position] != rune('T') {
							goto l860
						}
						position++
					}
				l867:
					{
						position869, tokenIndex869, depth869 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l870
						}
						position++
						goto l869
					l870:
						position, tokenIndex, depth = position869, tokenIndex869, depth869
						if buffer[position] != rune('R') {
							goto l860
						}
						position++
					}
				l869:
					{
						position871, tokenIndex871, depth871 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l872
						}
						position++
						goto l871
					l872:
						position, tokenIndex, depth = position871, tokenIndex871, depth871
						if buffer[position] != rune('E') {
							goto l860
						}
						position++
					}
				l871:
					{
						position873, tokenIndex873, depth873 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l874
						}
						position++
						goto l873
					l874:
						position, tokenIndex, depth = position873, tokenIndex873, depth873
						if buffer[position] != rune('A') {
							goto l860
						}
						position++
					}
				l873:
					{
						position875, tokenIndex875, depth875 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l876
						}
						position++
						goto l875
					l876:
						position, tokenIndex, depth = position875, tokenIndex875, depth875
						if buffer[position] != rune('M') {
							goto l860
						}
						position++
					}
				l875:
					depth--
					add(rulePegText, position862)
				}
				if !_rules[ruleAction59]() {
					goto l860
				}
				depth--
				add(ruleDSTREAM, position861)
			}
			return true
		l860:
			position, tokenIndex, depth = position860, tokenIndex860, depth860
			return false
		},
		/* 82 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action60)> */
		func() bool {
			position877, tokenIndex877, depth877 := position, tokenIndex, depth
			{
				position878 := position
				depth++
				{
					position879 := position
					depth++
					{
						position880, tokenIndex880, depth880 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l881
						}
						position++
						goto l880
					l881:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
						if buffer[position] != rune('R') {
							goto l877
						}
						position++
					}
				l880:
					{
						position882, tokenIndex882, depth882 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l883
						}
						position++
						goto l882
					l883:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if buffer[position] != rune('S') {
							goto l877
						}
						position++
					}
				l882:
					{
						position884, tokenIndex884, depth884 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l885
						}
						position++
						goto l884
					l885:
						position, tokenIndex, depth = position884, tokenIndex884, depth884
						if buffer[position] != rune('T') {
							goto l877
						}
						position++
					}
				l884:
					{
						position886, tokenIndex886, depth886 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l887
						}
						position++
						goto l886
					l887:
						position, tokenIndex, depth = position886, tokenIndex886, depth886
						if buffer[position] != rune('R') {
							goto l877
						}
						position++
					}
				l886:
					{
						position888, tokenIndex888, depth888 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l889
						}
						position++
						goto l888
					l889:
						position, tokenIndex, depth = position888, tokenIndex888, depth888
						if buffer[position] != rune('E') {
							goto l877
						}
						position++
					}
				l888:
					{
						position890, tokenIndex890, depth890 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l891
						}
						position++
						goto l890
					l891:
						position, tokenIndex, depth = position890, tokenIndex890, depth890
						if buffer[position] != rune('A') {
							goto l877
						}
						position++
					}
				l890:
					{
						position892, tokenIndex892, depth892 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l893
						}
						position++
						goto l892
					l893:
						position, tokenIndex, depth = position892, tokenIndex892, depth892
						if buffer[position] != rune('M') {
							goto l877
						}
						position++
					}
				l892:
					depth--
					add(rulePegText, position879)
				}
				if !_rules[ruleAction60]() {
					goto l877
				}
				depth--
				add(ruleRSTREAM, position878)
			}
			return true
		l877:
			position, tokenIndex, depth = position877, tokenIndex877, depth877
			return false
		},
		/* 83 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action61)> */
		func() bool {
			position894, tokenIndex894, depth894 := position, tokenIndex, depth
			{
				position895 := position
				depth++
				{
					position896 := position
					depth++
					{
						position897, tokenIndex897, depth897 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l898
						}
						position++
						goto l897
					l898:
						position, tokenIndex, depth = position897, tokenIndex897, depth897
						if buffer[position] != rune('T') {
							goto l894
						}
						position++
					}
				l897:
					{
						position899, tokenIndex899, depth899 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l900
						}
						position++
						goto l899
					l900:
						position, tokenIndex, depth = position899, tokenIndex899, depth899
						if buffer[position] != rune('U') {
							goto l894
						}
						position++
					}
				l899:
					{
						position901, tokenIndex901, depth901 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l902
						}
						position++
						goto l901
					l902:
						position, tokenIndex, depth = position901, tokenIndex901, depth901
						if buffer[position] != rune('P') {
							goto l894
						}
						position++
					}
				l901:
					{
						position903, tokenIndex903, depth903 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l904
						}
						position++
						goto l903
					l904:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if buffer[position] != rune('L') {
							goto l894
						}
						position++
					}
				l903:
					{
						position905, tokenIndex905, depth905 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l906
						}
						position++
						goto l905
					l906:
						position, tokenIndex, depth = position905, tokenIndex905, depth905
						if buffer[position] != rune('E') {
							goto l894
						}
						position++
					}
				l905:
					{
						position907, tokenIndex907, depth907 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l908
						}
						position++
						goto l907
					l908:
						position, tokenIndex, depth = position907, tokenIndex907, depth907
						if buffer[position] != rune('S') {
							goto l894
						}
						position++
					}
				l907:
					depth--
					add(rulePegText, position896)
				}
				if !_rules[ruleAction61]() {
					goto l894
				}
				depth--
				add(ruleTUPLES, position895)
			}
			return true
		l894:
			position, tokenIndex, depth = position894, tokenIndex894, depth894
			return false
		},
		/* 84 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action62)> */
		func() bool {
			position909, tokenIndex909, depth909 := position, tokenIndex, depth
			{
				position910 := position
				depth++
				{
					position911 := position
					depth++
					{
						position912, tokenIndex912, depth912 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l913
						}
						position++
						goto l912
					l913:
						position, tokenIndex, depth = position912, tokenIndex912, depth912
						if buffer[position] != rune('S') {
							goto l909
						}
						position++
					}
				l912:
					{
						position914, tokenIndex914, depth914 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l915
						}
						position++
						goto l914
					l915:
						position, tokenIndex, depth = position914, tokenIndex914, depth914
						if buffer[position] != rune('E') {
							goto l909
						}
						position++
					}
				l914:
					{
						position916, tokenIndex916, depth916 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l917
						}
						position++
						goto l916
					l917:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						if buffer[position] != rune('C') {
							goto l909
						}
						position++
					}
				l916:
					{
						position918, tokenIndex918, depth918 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l919
						}
						position++
						goto l918
					l919:
						position, tokenIndex, depth = position918, tokenIndex918, depth918
						if buffer[position] != rune('O') {
							goto l909
						}
						position++
					}
				l918:
					{
						position920, tokenIndex920, depth920 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l921
						}
						position++
						goto l920
					l921:
						position, tokenIndex, depth = position920, tokenIndex920, depth920
						if buffer[position] != rune('N') {
							goto l909
						}
						position++
					}
				l920:
					{
						position922, tokenIndex922, depth922 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l923
						}
						position++
						goto l922
					l923:
						position, tokenIndex, depth = position922, tokenIndex922, depth922
						if buffer[position] != rune('D') {
							goto l909
						}
						position++
					}
				l922:
					{
						position924, tokenIndex924, depth924 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l925
						}
						position++
						goto l924
					l925:
						position, tokenIndex, depth = position924, tokenIndex924, depth924
						if buffer[position] != rune('S') {
							goto l909
						}
						position++
					}
				l924:
					depth--
					add(rulePegText, position911)
				}
				if !_rules[ruleAction62]() {
					goto l909
				}
				depth--
				add(ruleSECONDS, position910)
			}
			return true
		l909:
			position, tokenIndex, depth = position909, tokenIndex909, depth909
			return false
		},
		/* 85 StreamIdentifier <- <(<ident> Action63)> */
		func() bool {
			position926, tokenIndex926, depth926 := position, tokenIndex, depth
			{
				position927 := position
				depth++
				{
					position928 := position
					depth++
					if !_rules[ruleident]() {
						goto l926
					}
					depth--
					add(rulePegText, position928)
				}
				if !_rules[ruleAction63]() {
					goto l926
				}
				depth--
				add(ruleStreamIdentifier, position927)
			}
			return true
		l926:
			position, tokenIndex, depth = position926, tokenIndex926, depth926
			return false
		},
		/* 86 SourceSinkType <- <(<ident> Action64)> */
		func() bool {
			position929, tokenIndex929, depth929 := position, tokenIndex, depth
			{
				position930 := position
				depth++
				{
					position931 := position
					depth++
					if !_rules[ruleident]() {
						goto l929
					}
					depth--
					add(rulePegText, position931)
				}
				if !_rules[ruleAction64]() {
					goto l929
				}
				depth--
				add(ruleSourceSinkType, position930)
			}
			return true
		l929:
			position, tokenIndex, depth = position929, tokenIndex929, depth929
			return false
		},
		/* 87 SourceSinkParamKey <- <(<ident> Action65)> */
		func() bool {
			position932, tokenIndex932, depth932 := position, tokenIndex, depth
			{
				position933 := position
				depth++
				{
					position934 := position
					depth++
					if !_rules[ruleident]() {
						goto l932
					}
					depth--
					add(rulePegText, position934)
				}
				if !_rules[ruleAction65]() {
					goto l932
				}
				depth--
				add(ruleSourceSinkParamKey, position933)
			}
			return true
		l932:
			position, tokenIndex, depth = position932, tokenIndex932, depth932
			return false
		},
		/* 88 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action66)> */
		func() bool {
			position935, tokenIndex935, depth935 := position, tokenIndex, depth
			{
				position936 := position
				depth++
				{
					position937 := position
					depth++
					{
						position938, tokenIndex938, depth938 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l939
						}
						position++
						goto l938
					l939:
						position, tokenIndex, depth = position938, tokenIndex938, depth938
						if buffer[position] != rune('P') {
							goto l935
						}
						position++
					}
				l938:
					{
						position940, tokenIndex940, depth940 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l941
						}
						position++
						goto l940
					l941:
						position, tokenIndex, depth = position940, tokenIndex940, depth940
						if buffer[position] != rune('A') {
							goto l935
						}
						position++
					}
				l940:
					{
						position942, tokenIndex942, depth942 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l943
						}
						position++
						goto l942
					l943:
						position, tokenIndex, depth = position942, tokenIndex942, depth942
						if buffer[position] != rune('U') {
							goto l935
						}
						position++
					}
				l942:
					{
						position944, tokenIndex944, depth944 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l945
						}
						position++
						goto l944
					l945:
						position, tokenIndex, depth = position944, tokenIndex944, depth944
						if buffer[position] != rune('S') {
							goto l935
						}
						position++
					}
				l944:
					{
						position946, tokenIndex946, depth946 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l947
						}
						position++
						goto l946
					l947:
						position, tokenIndex, depth = position946, tokenIndex946, depth946
						if buffer[position] != rune('E') {
							goto l935
						}
						position++
					}
				l946:
					{
						position948, tokenIndex948, depth948 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l949
						}
						position++
						goto l948
					l949:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
						if buffer[position] != rune('D') {
							goto l935
						}
						position++
					}
				l948:
					depth--
					add(rulePegText, position937)
				}
				if !_rules[ruleAction66]() {
					goto l935
				}
				depth--
				add(rulePaused, position936)
			}
			return true
		l935:
			position, tokenIndex, depth = position935, tokenIndex935, depth935
			return false
		},
		/* 89 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action67)> */
		func() bool {
			position950, tokenIndex950, depth950 := position, tokenIndex, depth
			{
				position951 := position
				depth++
				{
					position952 := position
					depth++
					{
						position953, tokenIndex953, depth953 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l954
						}
						position++
						goto l953
					l954:
						position, tokenIndex, depth = position953, tokenIndex953, depth953
						if buffer[position] != rune('U') {
							goto l950
						}
						position++
					}
				l953:
					{
						position955, tokenIndex955, depth955 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l956
						}
						position++
						goto l955
					l956:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('N') {
							goto l950
						}
						position++
					}
				l955:
					{
						position957, tokenIndex957, depth957 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l958
						}
						position++
						goto l957
					l958:
						position, tokenIndex, depth = position957, tokenIndex957, depth957
						if buffer[position] != rune('P') {
							goto l950
						}
						position++
					}
				l957:
					{
						position959, tokenIndex959, depth959 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l960
						}
						position++
						goto l959
					l960:
						position, tokenIndex, depth = position959, tokenIndex959, depth959
						if buffer[position] != rune('A') {
							goto l950
						}
						position++
					}
				l959:
					{
						position961, tokenIndex961, depth961 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l962
						}
						position++
						goto l961
					l962:
						position, tokenIndex, depth = position961, tokenIndex961, depth961
						if buffer[position] != rune('U') {
							goto l950
						}
						position++
					}
				l961:
					{
						position963, tokenIndex963, depth963 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l964
						}
						position++
						goto l963
					l964:
						position, tokenIndex, depth = position963, tokenIndex963, depth963
						if buffer[position] != rune('S') {
							goto l950
						}
						position++
					}
				l963:
					{
						position965, tokenIndex965, depth965 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l966
						}
						position++
						goto l965
					l966:
						position, tokenIndex, depth = position965, tokenIndex965, depth965
						if buffer[position] != rune('E') {
							goto l950
						}
						position++
					}
				l965:
					{
						position967, tokenIndex967, depth967 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l968
						}
						position++
						goto l967
					l968:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if buffer[position] != rune('D') {
							goto l950
						}
						position++
					}
				l967:
					depth--
					add(rulePegText, position952)
				}
				if !_rules[ruleAction67]() {
					goto l950
				}
				depth--
				add(ruleUnpaused, position951)
			}
			return true
		l950:
			position, tokenIndex, depth = position950, tokenIndex950, depth950
			return false
		},
		/* 90 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action68)> */
		func() bool {
			position969, tokenIndex969, depth969 := position, tokenIndex, depth
			{
				position970 := position
				depth++
				{
					position971 := position
					depth++
					{
						position972, tokenIndex972, depth972 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l973
						}
						position++
						goto l972
					l973:
						position, tokenIndex, depth = position972, tokenIndex972, depth972
						if buffer[position] != rune('O') {
							goto l969
						}
						position++
					}
				l972:
					{
						position974, tokenIndex974, depth974 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l975
						}
						position++
						goto l974
					l975:
						position, tokenIndex, depth = position974, tokenIndex974, depth974
						if buffer[position] != rune('R') {
							goto l969
						}
						position++
					}
				l974:
					depth--
					add(rulePegText, position971)
				}
				if !_rules[ruleAction68]() {
					goto l969
				}
				depth--
				add(ruleOr, position970)
			}
			return true
		l969:
			position, tokenIndex, depth = position969, tokenIndex969, depth969
			return false
		},
		/* 91 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action69)> */
		func() bool {
			position976, tokenIndex976, depth976 := position, tokenIndex, depth
			{
				position977 := position
				depth++
				{
					position978 := position
					depth++
					{
						position979, tokenIndex979, depth979 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l980
						}
						position++
						goto l979
					l980:
						position, tokenIndex, depth = position979, tokenIndex979, depth979
						if buffer[position] != rune('A') {
							goto l976
						}
						position++
					}
				l979:
					{
						position981, tokenIndex981, depth981 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l982
						}
						position++
						goto l981
					l982:
						position, tokenIndex, depth = position981, tokenIndex981, depth981
						if buffer[position] != rune('N') {
							goto l976
						}
						position++
					}
				l981:
					{
						position983, tokenIndex983, depth983 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l984
						}
						position++
						goto l983
					l984:
						position, tokenIndex, depth = position983, tokenIndex983, depth983
						if buffer[position] != rune('D') {
							goto l976
						}
						position++
					}
				l983:
					depth--
					add(rulePegText, position978)
				}
				if !_rules[ruleAction69]() {
					goto l976
				}
				depth--
				add(ruleAnd, position977)
			}
			return true
		l976:
			position, tokenIndex, depth = position976, tokenIndex976, depth976
			return false
		},
		/* 92 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action70)> */
		func() bool {
			position985, tokenIndex985, depth985 := position, tokenIndex, depth
			{
				position986 := position
				depth++
				{
					position987 := position
					depth++
					{
						position988, tokenIndex988, depth988 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l989
						}
						position++
						goto l988
					l989:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('N') {
							goto l985
						}
						position++
					}
				l988:
					{
						position990, tokenIndex990, depth990 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l991
						}
						position++
						goto l990
					l991:
						position, tokenIndex, depth = position990, tokenIndex990, depth990
						if buffer[position] != rune('O') {
							goto l985
						}
						position++
					}
				l990:
					{
						position992, tokenIndex992, depth992 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l993
						}
						position++
						goto l992
					l993:
						position, tokenIndex, depth = position992, tokenIndex992, depth992
						if buffer[position] != rune('T') {
							goto l985
						}
						position++
					}
				l992:
					depth--
					add(rulePegText, position987)
				}
				if !_rules[ruleAction70]() {
					goto l985
				}
				depth--
				add(ruleNot, position986)
			}
			return true
		l985:
			position, tokenIndex, depth = position985, tokenIndex985, depth985
			return false
		},
		/* 93 Equal <- <(<'='> Action71)> */
		func() bool {
			position994, tokenIndex994, depth994 := position, tokenIndex, depth
			{
				position995 := position
				depth++
				{
					position996 := position
					depth++
					if buffer[position] != rune('=') {
						goto l994
					}
					position++
					depth--
					add(rulePegText, position996)
				}
				if !_rules[ruleAction71]() {
					goto l994
				}
				depth--
				add(ruleEqual, position995)
			}
			return true
		l994:
			position, tokenIndex, depth = position994, tokenIndex994, depth994
			return false
		},
		/* 94 Less <- <(<'<'> Action72)> */
		func() bool {
			position997, tokenIndex997, depth997 := position, tokenIndex, depth
			{
				position998 := position
				depth++
				{
					position999 := position
					depth++
					if buffer[position] != rune('<') {
						goto l997
					}
					position++
					depth--
					add(rulePegText, position999)
				}
				if !_rules[ruleAction72]() {
					goto l997
				}
				depth--
				add(ruleLess, position998)
			}
			return true
		l997:
			position, tokenIndex, depth = position997, tokenIndex997, depth997
			return false
		},
		/* 95 LessOrEqual <- <(<('<' '=')> Action73)> */
		func() bool {
			position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
			{
				position1001 := position
				depth++
				{
					position1002 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1000
					}
					position++
					if buffer[position] != rune('=') {
						goto l1000
					}
					position++
					depth--
					add(rulePegText, position1002)
				}
				if !_rules[ruleAction73]() {
					goto l1000
				}
				depth--
				add(ruleLessOrEqual, position1001)
			}
			return true
		l1000:
			position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
			return false
		},
		/* 96 Greater <- <(<'>'> Action74)> */
		func() bool {
			position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
			{
				position1004 := position
				depth++
				{
					position1005 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1003
					}
					position++
					depth--
					add(rulePegText, position1005)
				}
				if !_rules[ruleAction74]() {
					goto l1003
				}
				depth--
				add(ruleGreater, position1004)
			}
			return true
		l1003:
			position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
			return false
		},
		/* 97 GreaterOrEqual <- <(<('>' '=')> Action75)> */
		func() bool {
			position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
			{
				position1007 := position
				depth++
				{
					position1008 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1006
					}
					position++
					if buffer[position] != rune('=') {
						goto l1006
					}
					position++
					depth--
					add(rulePegText, position1008)
				}
				if !_rules[ruleAction75]() {
					goto l1006
				}
				depth--
				add(ruleGreaterOrEqual, position1007)
			}
			return true
		l1006:
			position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
			return false
		},
		/* 98 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action76)> */
		func() bool {
			position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
			{
				position1010 := position
				depth++
				{
					position1011 := position
					depth++
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1013
						}
						position++
						if buffer[position] != rune('=') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('<') {
							goto l1009
						}
						position++
						if buffer[position] != rune('>') {
							goto l1009
						}
						position++
					}
				l1012:
					depth--
					add(rulePegText, position1011)
				}
				if !_rules[ruleAction76]() {
					goto l1009
				}
				depth--
				add(ruleNotEqual, position1010)
			}
			return true
		l1009:
			position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
			return false
		},
		/* 99 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action77)> */
		func() bool {
			position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
			{
				position1015 := position
				depth++
				{
					position1016 := position
					depth++
					{
						position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1018
						}
						position++
						goto l1017
					l1018:
						position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
						if buffer[position] != rune('I') {
							goto l1014
						}
						position++
					}
				l1017:
					{
						position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1020
						}
						position++
						goto l1019
					l1020:
						position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
						if buffer[position] != rune('S') {
							goto l1014
						}
						position++
					}
				l1019:
					depth--
					add(rulePegText, position1016)
				}
				if !_rules[ruleAction77]() {
					goto l1014
				}
				depth--
				add(ruleIs, position1015)
			}
			return true
		l1014:
			position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
			return false
		},
		/* 100 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action78)> */
		func() bool {
			position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
			{
				position1022 := position
				depth++
				{
					position1023 := position
					depth++
					{
						position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1025
						}
						position++
						goto l1024
					l1025:
						position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
						if buffer[position] != rune('I') {
							goto l1021
						}
						position++
					}
				l1024:
					{
						position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1027
						}
						position++
						goto l1026
					l1027:
						position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
						if buffer[position] != rune('S') {
							goto l1021
						}
						position++
					}
				l1026:
					if !_rules[rulesp]() {
						goto l1021
					}
					{
						position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1029
						}
						position++
						goto l1028
					l1029:
						position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
						if buffer[position] != rune('N') {
							goto l1021
						}
						position++
					}
				l1028:
					{
						position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1031
						}
						position++
						goto l1030
					l1031:
						position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
						if buffer[position] != rune('O') {
							goto l1021
						}
						position++
					}
				l1030:
					{
						position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1033
						}
						position++
						goto l1032
					l1033:
						position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
						if buffer[position] != rune('T') {
							goto l1021
						}
						position++
					}
				l1032:
					depth--
					add(rulePegText, position1023)
				}
				if !_rules[ruleAction78]() {
					goto l1021
				}
				depth--
				add(ruleIsNot, position1022)
			}
			return true
		l1021:
			position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
			return false
		},
		/* 101 Plus <- <(<'+'> Action79)> */
		func() bool {
			position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
			{
				position1035 := position
				depth++
				{
					position1036 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1034
					}
					position++
					depth--
					add(rulePegText, position1036)
				}
				if !_rules[ruleAction79]() {
					goto l1034
				}
				depth--
				add(rulePlus, position1035)
			}
			return true
		l1034:
			position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
			return false
		},
		/* 102 Minus <- <(<'-'> Action80)> */
		func() bool {
			position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
			{
				position1038 := position
				depth++
				{
					position1039 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1037
					}
					position++
					depth--
					add(rulePegText, position1039)
				}
				if !_rules[ruleAction80]() {
					goto l1037
				}
				depth--
				add(ruleMinus, position1038)
			}
			return true
		l1037:
			position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
			return false
		},
		/* 103 Multiply <- <(<'*'> Action81)> */
		func() bool {
			position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
			{
				position1041 := position
				depth++
				{
					position1042 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1040
					}
					position++
					depth--
					add(rulePegText, position1042)
				}
				if !_rules[ruleAction81]() {
					goto l1040
				}
				depth--
				add(ruleMultiply, position1041)
			}
			return true
		l1040:
			position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
			return false
		},
		/* 104 Divide <- <(<'/'> Action82)> */
		func() bool {
			position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
			{
				position1044 := position
				depth++
				{
					position1045 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1043
					}
					position++
					depth--
					add(rulePegText, position1045)
				}
				if !_rules[ruleAction82]() {
					goto l1043
				}
				depth--
				add(ruleDivide, position1044)
			}
			return true
		l1043:
			position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
			return false
		},
		/* 105 Modulo <- <(<'%'> Action83)> */
		func() bool {
			position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
			{
				position1047 := position
				depth++
				{
					position1048 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1046
					}
					position++
					depth--
					add(rulePegText, position1048)
				}
				if !_rules[ruleAction83]() {
					goto l1046
				}
				depth--
				add(ruleModulo, position1047)
			}
			return true
		l1046:
			position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
			return false
		},
		/* 106 UnaryMinus <- <(<'-'> Action84)> */
		func() bool {
			position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
			{
				position1050 := position
				depth++
				{
					position1051 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1049
					}
					position++
					depth--
					add(rulePegText, position1051)
				}
				if !_rules[ruleAction84]() {
					goto l1049
				}
				depth--
				add(ruleUnaryMinus, position1050)
			}
			return true
		l1049:
			position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
			return false
		},
		/* 107 Identifier <- <(<ident> Action85)> */
		func() bool {
			position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
			{
				position1053 := position
				depth++
				{
					position1054 := position
					depth++
					if !_rules[ruleident]() {
						goto l1052
					}
					depth--
					add(rulePegText, position1054)
				}
				if !_rules[ruleAction85]() {
					goto l1052
				}
				depth--
				add(ruleIdentifier, position1053)
			}
			return true
		l1052:
			position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
			return false
		},
		/* 108 TargetIdentifier <- <(<jsonPath> Action86)> */
		func() bool {
			position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
			{
				position1056 := position
				depth++
				{
					position1057 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1055
					}
					depth--
					add(rulePegText, position1057)
				}
				if !_rules[ruleAction86]() {
					goto l1055
				}
				depth--
				add(ruleTargetIdentifier, position1056)
			}
			return true
		l1055:
			position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
			return false
		},
		/* 109 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
			{
				position1059 := position
				depth++
				{
					position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1061
					}
					position++
					goto l1060
				l1061:
					position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1058
					}
					position++
				}
			l1060:
			l1062:
				{
					position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
					{
						position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1065
						}
						position++
						goto l1064
					l1065:
						position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1066
						}
						position++
						goto l1064
					l1066:
						position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1067
						}
						position++
						goto l1064
					l1067:
						position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
						if buffer[position] != rune('_') {
							goto l1063
						}
						position++
					}
				l1064:
					goto l1062
				l1063:
					position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
				}
				depth--
				add(ruleident, position1059)
			}
			return true
		l1058:
			position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
			return false
		},
		/* 110 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
			{
				position1069 := position
				depth++
				{
					position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1071
					}
					position++
					goto l1070
				l1071:
					position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1068
					}
					position++
				}
			l1070:
			l1072:
				{
					position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
					{
						position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1075
						}
						position++
						goto l1074
					l1075:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1076
						}
						position++
						goto l1074
					l1076:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1077
						}
						position++
						goto l1074
					l1077:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('_') {
							goto l1078
						}
						position++
						goto l1074
					l1078:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('.') {
							goto l1079
						}
						position++
						goto l1074
					l1079:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('[') {
							goto l1080
						}
						position++
						goto l1074
					l1080:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune(']') {
							goto l1081
						}
						position++
						goto l1074
					l1081:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('"') {
							goto l1073
						}
						position++
					}
				l1074:
					goto l1072
				l1073:
					position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
				}
				depth--
				add(rulejsonPath, position1069)
			}
			return true
		l1068:
			position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
			return false
		},
		/* 111 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1083 := position
				depth++
			l1084:
				{
					position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
					{
						position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1087
						}
						position++
						goto l1086
					l1087:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if buffer[position] != rune('\t') {
							goto l1088
						}
						position++
						goto l1086
					l1088:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if buffer[position] != rune('\n') {
							goto l1089
						}
						position++
						goto l1086
					l1089:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if buffer[position] != rune('\r') {
							goto l1090
						}
						position++
						goto l1086
					l1090:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if !_rules[rulecomment]() {
							goto l1085
						}
					}
				l1086:
					goto l1084
				l1085:
					position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
				}
				depth--
				add(rulesp, position1083)
			}
			return true
		},
		/* 112 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
			{
				position1092 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1091
				}
				position++
				if buffer[position] != rune('-') {
					goto l1091
				}
				position++
			l1093:
				{
					position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
					{
						position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
						{
							position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1097
							}
							position++
							goto l1096
						l1097:
							position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
							if buffer[position] != rune('\n') {
								goto l1095
							}
							position++
						}
					l1096:
						goto l1094
					l1095:
						position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
					}
					if !matchDot() {
						goto l1094
					}
					goto l1093
				l1094:
					position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
				}
				{
					position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1099
					}
					position++
					goto l1098
				l1099:
					position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
					if buffer[position] != rune('\n') {
						goto l1091
					}
					position++
				}
			l1098:
				depth--
				add(rulecomment, position1092)
			}
			return true
		l1091:
			position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
			return false
		},
		/* 114 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 115 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 116 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 117 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 118 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 119 Action5 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 120 Action6 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 121 Action7 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 122 Action8 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 123 Action9 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 124 Action10 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 125 Action11 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 126 Action12 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 127 Action13 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 128 Action14 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 129 Action15 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 130 Action16 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		nil,
		/* 132 Action17 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 133 Action18 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 134 Action19 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 135 Action20 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 136 Action21 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 137 Action22 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 138 Action23 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 139 Action24 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 140 Action25 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 141 Action26 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 142 Action27 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 143 Action28 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 144 Action29 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 145 Action30 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 146 Action31 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 147 Action32 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 148 Action33 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 149 Action34 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 150 Action35 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 151 Action36 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 152 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 153 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 154 Action39 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 155 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 156 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 157 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 158 Action43 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 159 Action44 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 160 Action45 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 161 Action46 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 162 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 163 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 164 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 165 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 166 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 167 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 168 Action53 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 169 Action54 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 170 Action55 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 171 Action56 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 172 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 173 Action58 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 174 Action59 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 175 Action60 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 176 Action61 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 177 Action62 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 178 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 179 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 180 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 181 Action66 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 182 Action67 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 183 Action68 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 184 Action69 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 185 Action70 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 186 Action71 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 187 Action72 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 188 Action73 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 189 Action74 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 190 Action75 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 191 Action76 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 192 Action77 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 193 Action78 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 194 Action79 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 195 Action80 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 196 Action81 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 197 Action82 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 198 Action83 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 199 Action84 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 200 Action85 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 201 Action86 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
	}
	p.rules = _rules
}

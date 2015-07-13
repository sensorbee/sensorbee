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
	ruleUpdateStateStmt
	ruleUpdatSourceStmt
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
	rulePegText
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
	ruleAction80
	ruleAction81
	ruleAction82
	ruleAction83
	ruleAction84
	ruleAction85

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
	"UpdateStateStmt",
	"UpdatSourceStmt",
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
	"PegText",
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
	"Action80",
	"Action81",
	"Action82",
	"Action83",
	"Action84",
	"Action85",

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
	rules  [196]func() bool
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

			p.AssembleInsertIntoSelect()

		case ruleAction8:

			p.AssembleInsertIntoFrom()

		case ruleAction9:

			p.AssemblePauseSource()

		case ruleAction10:

			p.AssembleResumeSource()

		case ruleAction11:

			p.AssembleRewindSource()

		case ruleAction12:

			p.AssembleDropSource()

		case ruleAction13:

			p.AssembleDropStream()

		case ruleAction14:

			p.AssembleDropSink()

		case ruleAction15:

			p.AssembleDropState()

		case ruleAction16:

			p.AssembleEmitter(begin, end)

		case ruleAction17:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction18:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction19:

			p.AssembleStreamEmitInterval()

		case ruleAction20:

			p.AssembleProjections(begin, end)

		case ruleAction21:

			p.AssembleAlias()

		case ruleAction22:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction23:

			p.AssembleInterval()

		case ruleAction24:

			p.AssembleInterval()

		case ruleAction25:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction26:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction27:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction28:

			p.EnsureAliasedStreamWindow()

		case ruleAction29:

			p.AssembleAliasedStreamWindow()

		case ruleAction30:

			p.AssembleStreamWindow()

		case ruleAction31:

			p.AssembleUDSFFuncApp()

		case ruleAction32:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction33:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction34:

			p.AssembleSourceSinkParam()

		case ruleAction35:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction38:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction39:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction40:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction41:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction42:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction43:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction44:

			p.AssembleFuncApp()

		case ruleAction45:

			p.AssembleExpressions(begin, end)

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction52:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction53:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction54:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction55:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction56:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction57:

			p.PushComponent(begin, end, Istream)

		case ruleAction58:

			p.PushComponent(begin, end, Dstream)

		case ruleAction59:

			p.PushComponent(begin, end, Rstream)

		case ruleAction60:

			p.PushComponent(begin, end, Tuples)

		case ruleAction61:

			p.PushComponent(begin, end, Seconds)

		case ruleAction62:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction63:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction64:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction65:

			p.PushComponent(begin, end, Yes)

		case ruleAction66:

			p.PushComponent(begin, end, No)

		case ruleAction67:

			p.PushComponent(begin, end, Or)

		case ruleAction68:

			p.PushComponent(begin, end, And)

		case ruleAction69:

			p.PushComponent(begin, end, Not)

		case ruleAction70:

			p.PushComponent(begin, end, Equal)

		case ruleAction71:

			p.PushComponent(begin, end, Less)

		case ruleAction72:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction73:

			p.PushComponent(begin, end, Greater)

		case ruleAction74:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction75:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction76:

			p.PushComponent(begin, end, Is)

		case ruleAction77:

			p.PushComponent(begin, end, IsNot)

		case ruleAction78:

			p.PushComponent(begin, end, Plus)

		case ruleAction79:

			p.PushComponent(begin, end, Minus)

		case ruleAction80:

			p.PushComponent(begin, end, Multiply)

		case ruleAction81:

			p.PushComponent(begin, end, Divide)

		case ruleAction82:

			p.PushComponent(begin, end, Modulo)

		case ruleAction83:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction84:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction85:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt / DropSourceStmt / DropStreamStmt / DropSinkStmt / DropStateStmt / UpdateStateStmt / UpdatSourceStmt)> */
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
						goto l20
					}
					goto l9
				l20:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleDropStreamStmt]() {
						goto l21
					}
					goto l9
				l21:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleDropSinkStmt]() {
						goto l22
					}
					goto l9
				l22:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleDropStateStmt]() {
						goto l23
					}
					goto l9
				l23:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleUpdateStateStmt]() {
						goto l24
					}
					goto l9
				l24:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleUpdatSourceStmt]() {
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
			position25, tokenIndex25, depth25 := position, tokenIndex, depth
			{
				position26 := position
				depth++
				{
					position27, tokenIndex27, depth27 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l28
					}
					position++
					goto l27
				l28:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					if buffer[position] != rune('S') {
						goto l25
					}
					position++
				}
			l27:
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l30
					}
					position++
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if buffer[position] != rune('E') {
						goto l25
					}
					position++
				}
			l29:
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l32
					}
					position++
					goto l31
				l32:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
					if buffer[position] != rune('L') {
						goto l25
					}
					position++
				}
			l31:
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l34
					}
					position++
					goto l33
				l34:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
					if buffer[position] != rune('E') {
						goto l25
					}
					position++
				}
			l33:
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
						goto l25
					}
					position++
				}
			l35:
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l38
					}
					position++
					goto l37
				l38:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if buffer[position] != rune('T') {
						goto l25
					}
					position++
				}
			l37:
				if !_rules[rulesp]() {
					goto l25
				}
				if !_rules[ruleEmitter]() {
					goto l25
				}
				if !_rules[rulesp]() {
					goto l25
				}
				if !_rules[ruleProjections]() {
					goto l25
				}
				if !_rules[rulesp]() {
					goto l25
				}
				if !_rules[ruleWindowedFrom]() {
					goto l25
				}
				if !_rules[rulesp]() {
					goto l25
				}
				if !_rules[ruleFilter]() {
					goto l25
				}
				if !_rules[rulesp]() {
					goto l25
				}
				if !_rules[ruleGrouping]() {
					goto l25
				}
				if !_rules[rulesp]() {
					goto l25
				}
				if !_rules[ruleHaving]() {
					goto l25
				}
				if !_rules[rulesp]() {
					goto l25
				}
				if !_rules[ruleAction0]() {
					goto l25
				}
				depth--
				add(ruleSelectStmt, position26)
			}
			return true
		l25:
			position, tokenIndex, depth = position25, tokenIndex25, depth25
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position39, tokenIndex39, depth39 := position, tokenIndex, depth
			{
				position40 := position
				depth++
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l42
					}
					position++
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if buffer[position] != rune('C') {
						goto l39
					}
					position++
				}
			l41:
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l44
					}
					position++
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if buffer[position] != rune('R') {
						goto l39
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
						goto l39
					}
					position++
				}
			l45:
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l48
					}
					position++
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					if buffer[position] != rune('A') {
						goto l39
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
						goto l39
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
						goto l39
					}
					position++
				}
			l51:
				if !_rules[rulesp]() {
					goto l39
				}
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					if buffer[position] != rune('S') {
						goto l39
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
					if buffer[position] != rune('T') {
						goto l39
					}
					position++
				}
			l55:
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
					if buffer[position] != rune('R') {
						goto l39
					}
					position++
				}
			l57:
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l60
					}
					position++
					goto l59
				l60:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
					if buffer[position] != rune('E') {
						goto l39
					}
					position++
				}
			l59:
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
						goto l39
					}
					position++
				}
			l61:
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l64
					}
					position++
					goto l63
				l64:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
					if buffer[position] != rune('M') {
						goto l39
					}
					position++
				}
			l63:
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				{
					position65, tokenIndex65, depth65 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l66
					}
					position++
					goto l65
				l66:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
					if buffer[position] != rune('A') {
						goto l39
					}
					position++
				}
			l65:
				{
					position67, tokenIndex67, depth67 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l68
					}
					position++
					goto l67
				l68:
					position, tokenIndex, depth = position67, tokenIndex67, depth67
					if buffer[position] != rune('S') {
						goto l39
					}
					position++
				}
			l67:
				if !_rules[rulesp]() {
					goto l39
				}
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l70
					}
					position++
					goto l69
				l70:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
					if buffer[position] != rune('S') {
						goto l39
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
						goto l39
					}
					position++
				}
			l71:
				{
					position73, tokenIndex73, depth73 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l74
					}
					position++
					goto l73
				l74:
					position, tokenIndex, depth = position73, tokenIndex73, depth73
					if buffer[position] != rune('L') {
						goto l39
					}
					position++
				}
			l73:
				{
					position75, tokenIndex75, depth75 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l76
					}
					position++
					goto l75
				l76:
					position, tokenIndex, depth = position75, tokenIndex75, depth75
					if buffer[position] != rune('E') {
						goto l39
					}
					position++
				}
			l75:
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
						goto l39
					}
					position++
				}
			l77:
				{
					position79, tokenIndex79, depth79 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l80
					}
					position++
					goto l79
				l80:
					position, tokenIndex, depth = position79, tokenIndex79, depth79
					if buffer[position] != rune('T') {
						goto l39
					}
					position++
				}
			l79:
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleEmitter]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleProjections]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleWindowedFrom]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleFilter]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleGrouping]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleHaving]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleAction1]() {
					goto l39
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position40)
			}
			return true
		l39:
			position, tokenIndex, depth = position39, tokenIndex39, depth39
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position81, tokenIndex81, depth81 := position, tokenIndex, depth
			{
				position82 := position
				depth++
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l84
					}
					position++
					goto l83
				l84:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
					if buffer[position] != rune('C') {
						goto l81
					}
					position++
				}
			l83:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l86
					}
					position++
					goto l85
				l86:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
					if buffer[position] != rune('R') {
						goto l81
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
						goto l81
					}
					position++
				}
			l87:
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l90
					}
					position++
					goto l89
				l90:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
					if buffer[position] != rune('A') {
						goto l81
					}
					position++
				}
			l89:
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l92
					}
					position++
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					if buffer[position] != rune('T') {
						goto l81
					}
					position++
				}
			l91:
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('E') {
						goto l81
					}
					position++
				}
			l93:
				if !_rules[rulesp]() {
					goto l81
				}
				if !_rules[rulePausedOpt]() {
					goto l81
				}
				if !_rules[rulesp]() {
					goto l81
				}
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if buffer[position] != rune('S') {
						goto l81
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('O') {
						goto l81
					}
					position++
				}
			l97:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('U') {
						goto l81
					}
					position++
				}
			l99:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('R') {
						goto l81
					}
					position++
				}
			l101:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if buffer[position] != rune('C') {
						goto l81
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
						goto l81
					}
					position++
				}
			l105:
				if !_rules[rulesp]() {
					goto l81
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l81
				}
				if !_rules[rulesp]() {
					goto l81
				}
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l108
					}
					position++
					goto l107
				l108:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
					if buffer[position] != rune('T') {
						goto l81
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('Y') {
						goto l81
					}
					position++
				}
			l109:
				{
					position111, tokenIndex111, depth111 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l112
					}
					position++
					goto l111
				l112:
					position, tokenIndex, depth = position111, tokenIndex111, depth111
					if buffer[position] != rune('P') {
						goto l81
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
						goto l81
					}
					position++
				}
			l113:
				if !_rules[rulesp]() {
					goto l81
				}
				if !_rules[ruleSourceSinkType]() {
					goto l81
				}
				if !_rules[rulesp]() {
					goto l81
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l81
				}
				if !_rules[ruleAction2]() {
					goto l81
				}
				depth--
				add(ruleCreateSourceStmt, position82)
			}
			return true
		l81:
			position, tokenIndex, depth = position81, tokenIndex81, depth81
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position115, tokenIndex115, depth115 := position, tokenIndex, depth
			{
				position116 := position
				depth++
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('C') {
						goto l115
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('R') {
						goto l115
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
						goto l115
					}
					position++
				}
			l121:
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('A') {
						goto l115
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('T') {
						goto l115
					}
					position++
				}
			l125:
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('E') {
						goto l115
					}
					position++
				}
			l127:
				if !_rules[rulesp]() {
					goto l115
				}
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('S') {
						goto l115
					}
					position++
				}
			l129:
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('I') {
						goto l115
					}
					position++
				}
			l131:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('N') {
						goto l115
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('K') {
						goto l115
					}
					position++
				}
			l135:
				if !_rules[rulesp]() {
					goto l115
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l115
				}
				if !_rules[rulesp]() {
					goto l115
				}
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l138
					}
					position++
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if buffer[position] != rune('T') {
						goto l115
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('Y') {
						goto l115
					}
					position++
				}
			l139:
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l142
					}
					position++
					goto l141
				l142:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					if buffer[position] != rune('P') {
						goto l115
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
						goto l115
					}
					position++
				}
			l143:
				if !_rules[rulesp]() {
					goto l115
				}
				if !_rules[ruleSourceSinkType]() {
					goto l115
				}
				if !_rules[rulesp]() {
					goto l115
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l115
				}
				if !_rules[ruleAction3]() {
					goto l115
				}
				depth--
				add(ruleCreateSinkStmt, position116)
			}
			return true
		l115:
			position, tokenIndex, depth = position115, tokenIndex115, depth115
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position145, tokenIndex145, depth145 := position, tokenIndex, depth
			{
				position146 := position
				depth++
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('C') {
						goto l145
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('R') {
						goto l145
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
						goto l145
					}
					position++
				}
			l151:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('A') {
						goto l145
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
						goto l145
					}
					position++
				}
			l155:
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('E') {
						goto l145
					}
					position++
				}
			l157:
				if !_rules[rulesp]() {
					goto l145
				}
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l160
					}
					position++
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if buffer[position] != rune('S') {
						goto l145
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
						goto l145
					}
					position++
				}
			l161:
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l164
					}
					position++
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if buffer[position] != rune('A') {
						goto l145
					}
					position++
				}
			l163:
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
						goto l145
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
						goto l145
					}
					position++
				}
			l167:
				if !_rules[rulesp]() {
					goto l145
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l145
				}
				if !_rules[rulesp]() {
					goto l145
				}
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('T') {
						goto l145
					}
					position++
				}
			l169:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('Y') {
						goto l145
					}
					position++
				}
			l171:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('P') {
						goto l145
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('E') {
						goto l145
					}
					position++
				}
			l175:
				if !_rules[rulesp]() {
					goto l145
				}
				if !_rules[ruleSourceSinkType]() {
					goto l145
				}
				if !_rules[rulesp]() {
					goto l145
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l145
				}
				if !_rules[ruleAction4]() {
					goto l145
				}
				depth--
				add(ruleCreateStateStmt, position146)
			}
			return true
		l145:
			position, tokenIndex, depth = position145, tokenIndex145, depth145
			return false
		},
		/* 7 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action5)> */
		func() bool {
			position177, tokenIndex177, depth177 := position, tokenIndex, depth
			{
				position178 := position
				depth++
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('U') {
						goto l177
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('P') {
						goto l177
					}
					position++
				}
			l181:
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('D') {
						goto l177
					}
					position++
				}
			l183:
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('A') {
						goto l177
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
						goto l177
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('E') {
						goto l177
					}
					position++
				}
			l189:
				if !_rules[rulesp]() {
					goto l177
				}
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('S') {
						goto l177
					}
					position++
				}
			l191:
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('T') {
						goto l177
					}
					position++
				}
			l193:
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('A') {
						goto l177
					}
					position++
				}
			l195:
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('T') {
						goto l177
					}
					position++
				}
			l197:
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
						goto l177
					}
					position++
				}
			l199:
				if !_rules[rulesp]() {
					goto l177
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l177
				}
				if !_rules[rulesp]() {
					goto l177
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l177
				}
				if !_rules[ruleAction5]() {
					goto l177
				}
				depth--
				add(ruleUpdateStateStmt, position178)
			}
			return true
		l177:
			position, tokenIndex, depth = position177, tokenIndex177, depth177
			return false
		},
		/* 8 UpdatSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action6)> */
		func() bool {
			position201, tokenIndex201, depth201 := position, tokenIndex, depth
			{
				position202 := position
				depth++
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('U') {
						goto l201
					}
					position++
				}
			l203:
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('P') {
						goto l201
					}
					position++
				}
			l205:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('D') {
						goto l201
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('A') {
						goto l201
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
						goto l201
					}
					position++
				}
			l211:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('E') {
						goto l201
					}
					position++
				}
			l213:
				if !_rules[rulesp]() {
					goto l201
				}
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('S') {
						goto l201
					}
					position++
				}
			l215:
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('O') {
						goto l201
					}
					position++
				}
			l217:
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('U') {
						goto l201
					}
					position++
				}
			l219:
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if buffer[position] != rune('R') {
						goto l201
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('C') {
						goto l201
					}
					position++
				}
			l223:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('E') {
						goto l201
					}
					position++
				}
			l225:
				if !_rules[rulesp]() {
					goto l201
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l201
				}
				if !_rules[rulesp]() {
					goto l201
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l201
				}
				if !_rules[ruleAction6]() {
					goto l201
				}
				depth--
				add(ruleUpdatSourceStmt, position202)
			}
			return true
		l201:
			position, tokenIndex, depth = position201, tokenIndex201, depth201
			return false
		},
		/* 9 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action7)> */
		func() bool {
			position227, tokenIndex227, depth227 := position, tokenIndex, depth
			{
				position228 := position
				depth++
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('I') {
						goto l227
					}
					position++
				}
			l229:
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if buffer[position] != rune('N') {
						goto l227
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
						goto l227
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
						goto l227
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
						goto l227
					}
					position++
				}
			l237:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					if buffer[position] != rune('T') {
						goto l227
					}
					position++
				}
			l239:
				if !_rules[rulesp]() {
					goto l227
				}
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('I') {
						goto l227
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('N') {
						goto l227
					}
					position++
				}
			l243:
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l246
					}
					position++
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if buffer[position] != rune('T') {
						goto l227
					}
					position++
				}
			l245:
				{
					position247, tokenIndex247, depth247 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l248
					}
					position++
					goto l247
				l248:
					position, tokenIndex, depth = position247, tokenIndex247, depth247
					if buffer[position] != rune('O') {
						goto l227
					}
					position++
				}
			l247:
				if !_rules[rulesp]() {
					goto l227
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l227
				}
				if !_rules[rulesp]() {
					goto l227
				}
				if !_rules[ruleSelectStmt]() {
					goto l227
				}
				if !_rules[ruleAction7]() {
					goto l227
				}
				depth--
				add(ruleInsertIntoSelectStmt, position228)
			}
			return true
		l227:
			position, tokenIndex, depth = position227, tokenIndex227, depth227
			return false
		},
		/* 10 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action8)> */
		func() bool {
			position249, tokenIndex249, depth249 := position, tokenIndex, depth
			{
				position250 := position
				depth++
				{
					position251, tokenIndex251, depth251 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l252
					}
					position++
					goto l251
				l252:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('I') {
						goto l249
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('N') {
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
					if buffer[position] != rune('e') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('E') {
						goto l249
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('R') {
						goto l249
					}
					position++
				}
			l259:
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('T') {
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
					if buffer[position] != rune('i') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if buffer[position] != rune('I') {
						goto l249
					}
					position++
				}
			l263:
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('N') {
						goto l249
					}
					position++
				}
			l265:
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l268
					}
					position++
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if buffer[position] != rune('T') {
						goto l249
					}
					position++
				}
			l267:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('O') {
						goto l249
					}
					position++
				}
			l269:
				if !_rules[rulesp]() {
					goto l249
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l249
				}
				if !_rules[rulesp]() {
					goto l249
				}
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('F') {
						goto l249
					}
					position++
				}
			l271:
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('R') {
						goto l249
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('O') {
						goto l249
					}
					position++
				}
			l275:
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('M') {
						goto l249
					}
					position++
				}
			l277:
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
				add(ruleInsertIntoFromStmt, position250)
			}
			return true
		l249:
			position, tokenIndex, depth = position249, tokenIndex249, depth249
			return false
		},
		/* 11 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action9)> */
		func() bool {
			position279, tokenIndex279, depth279 := position, tokenIndex, depth
			{
				position280 := position
				depth++
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('P') {
						goto l279
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('A') {
						goto l279
					}
					position++
				}
			l283:
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('U') {
						goto l279
					}
					position++
				}
			l285:
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
						goto l279
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('E') {
						goto l279
					}
					position++
				}
			l289:
				if !_rules[rulesp]() {
					goto l279
				}
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('S') {
						goto l279
					}
					position++
				}
			l291:
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					if buffer[position] != rune('O') {
						goto l279
					}
					position++
				}
			l293:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('U') {
						goto l279
					}
					position++
				}
			l295:
				{
					position297, tokenIndex297, depth297 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l298
					}
					position++
					goto l297
				l298:
					position, tokenIndex, depth = position297, tokenIndex297, depth297
					if buffer[position] != rune('R') {
						goto l279
					}
					position++
				}
			l297:
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l300
					}
					position++
					goto l299
				l300:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
					if buffer[position] != rune('C') {
						goto l279
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
						goto l279
					}
					position++
				}
			l301:
				if !_rules[rulesp]() {
					goto l279
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l279
				}
				if !_rules[ruleAction9]() {
					goto l279
				}
				depth--
				add(rulePauseSourceStmt, position280)
			}
			return true
		l279:
			position, tokenIndex, depth = position279, tokenIndex279, depth279
			return false
		},
		/* 12 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
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
						goto l303
					}
					position++
				}
			l305:
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
						goto l303
					}
					position++
				}
			l307:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('S') {
						goto l303
					}
					position++
				}
			l309:
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l312
					}
					position++
					goto l311
				l312:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
					if buffer[position] != rune('U') {
						goto l303
					}
					position++
				}
			l311:
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('M') {
						goto l303
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('E') {
						goto l303
					}
					position++
				}
			l315:
				if !_rules[rulesp]() {
					goto l303
				}
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('S') {
						goto l303
					}
					position++
				}
			l317:
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l320
					}
					position++
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if buffer[position] != rune('O') {
						goto l303
					}
					position++
				}
			l319:
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l322
					}
					position++
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if buffer[position] != rune('U') {
						goto l303
					}
					position++
				}
			l321:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l324
					}
					position++
					goto l323
				l324:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					if buffer[position] != rune('R') {
						goto l303
					}
					position++
				}
			l323:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l326
					}
					position++
					goto l325
				l326:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
					if buffer[position] != rune('C') {
						goto l303
					}
					position++
				}
			l325:
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l328
					}
					position++
					goto l327
				l328:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if buffer[position] != rune('E') {
						goto l303
					}
					position++
				}
			l327:
				if !_rules[rulesp]() {
					goto l303
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l303
				}
				if !_rules[ruleAction10]() {
					goto l303
				}
				depth--
				add(ruleResumeSourceStmt, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 13 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action11)> */
		func() bool {
			position329, tokenIndex329, depth329 := position, tokenIndex, depth
			{
				position330 := position
				depth++
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('R') {
						goto l329
					}
					position++
				}
			l331:
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('E') {
						goto l329
					}
					position++
				}
			l333:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l336
					}
					position++
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if buffer[position] != rune('W') {
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
					if buffer[position] != rune('d') {
						goto l342
					}
					position++
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					if buffer[position] != rune('D') {
						goto l329
					}
					position++
				}
			l341:
				if !_rules[rulesp]() {
					goto l329
				}
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('S') {
						goto l329
					}
					position++
				}
			l343:
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l346
					}
					position++
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if buffer[position] != rune('O') {
						goto l329
					}
					position++
				}
			l345:
				{
					position347, tokenIndex347, depth347 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l348
					}
					position++
					goto l347
				l348:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if buffer[position] != rune('U') {
						goto l329
					}
					position++
				}
			l347:
				{
					position349, tokenIndex349, depth349 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l350
					}
					position++
					goto l349
				l350:
					position, tokenIndex, depth = position349, tokenIndex349, depth349
					if buffer[position] != rune('R') {
						goto l329
					}
					position++
				}
			l349:
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l352
					}
					position++
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					if buffer[position] != rune('C') {
						goto l329
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
						goto l329
					}
					position++
				}
			l353:
				if !_rules[rulesp]() {
					goto l329
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l329
				}
				if !_rules[ruleAction11]() {
					goto l329
				}
				depth--
				add(ruleRewindSourceStmt, position330)
			}
			return true
		l329:
			position, tokenIndex, depth = position329, tokenIndex329, depth329
			return false
		},
		/* 14 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action12)> */
		func() bool {
			position355, tokenIndex355, depth355 := position, tokenIndex, depth
			{
				position356 := position
				depth++
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l358
					}
					position++
					goto l357
				l358:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
					if buffer[position] != rune('D') {
						goto l355
					}
					position++
				}
			l357:
				{
					position359, tokenIndex359, depth359 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l360
					}
					position++
					goto l359
				l360:
					position, tokenIndex, depth = position359, tokenIndex359, depth359
					if buffer[position] != rune('R') {
						goto l355
					}
					position++
				}
			l359:
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l362
					}
					position++
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if buffer[position] != rune('O') {
						goto l355
					}
					position++
				}
			l361:
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l364
					}
					position++
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					if buffer[position] != rune('P') {
						goto l355
					}
					position++
				}
			l363:
				if !_rules[rulesp]() {
					goto l355
				}
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('S') {
						goto l355
					}
					position++
				}
			l365:
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					if buffer[position] != rune('O') {
						goto l355
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('U') {
						goto l355
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
						goto l355
					}
					position++
				}
			l371:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('C') {
						goto l355
					}
					position++
				}
			l373:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l376
					}
					position++
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					if buffer[position] != rune('E') {
						goto l355
					}
					position++
				}
			l375:
				if !_rules[rulesp]() {
					goto l355
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l355
				}
				if !_rules[ruleAction12]() {
					goto l355
				}
				depth--
				add(ruleDropSourceStmt, position356)
			}
			return true
		l355:
			position, tokenIndex, depth = position355, tokenIndex355, depth355
			return false
		},
		/* 15 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action13)> */
		func() bool {
			position377, tokenIndex377, depth377 := position, tokenIndex, depth
			{
				position378 := position
				depth++
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('D') {
						goto l377
					}
					position++
				}
			l379:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('R') {
						goto l377
					}
					position++
				}
			l381:
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('O') {
						goto l377
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('P') {
						goto l377
					}
					position++
				}
			l385:
				if !_rules[rulesp]() {
					goto l377
				}
				{
					position387, tokenIndex387, depth387 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l388
					}
					position++
					goto l387
				l388:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
					if buffer[position] != rune('S') {
						goto l377
					}
					position++
				}
			l387:
				{
					position389, tokenIndex389, depth389 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l390
					}
					position++
					goto l389
				l390:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
					if buffer[position] != rune('T') {
						goto l377
					}
					position++
				}
			l389:
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
					if buffer[position] != rune('R') {
						goto l377
					}
					position++
				}
			l391:
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l394
					}
					position++
					goto l393
				l394:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
					if buffer[position] != rune('E') {
						goto l377
					}
					position++
				}
			l393:
				{
					position395, tokenIndex395, depth395 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l396
					}
					position++
					goto l395
				l396:
					position, tokenIndex, depth = position395, tokenIndex395, depth395
					if buffer[position] != rune('A') {
						goto l377
					}
					position++
				}
			l395:
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l398
					}
					position++
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if buffer[position] != rune('M') {
						goto l377
					}
					position++
				}
			l397:
				if !_rules[rulesp]() {
					goto l377
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l377
				}
				if !_rules[ruleAction13]() {
					goto l377
				}
				depth--
				add(ruleDropStreamStmt, position378)
			}
			return true
		l377:
			position, tokenIndex, depth = position377, tokenIndex377, depth377
			return false
		},
		/* 16 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action14)> */
		func() bool {
			position399, tokenIndex399, depth399 := position, tokenIndex, depth
			{
				position400 := position
				depth++
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l402
					}
					position++
					goto l401
				l402:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if buffer[position] != rune('D') {
						goto l399
					}
					position++
				}
			l401:
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l404
					}
					position++
					goto l403
				l404:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
					if buffer[position] != rune('R') {
						goto l399
					}
					position++
				}
			l403:
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l406
					}
					position++
					goto l405
				l406:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if buffer[position] != rune('O') {
						goto l399
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
						goto l399
					}
					position++
				}
			l407:
				if !_rules[rulesp]() {
					goto l399
				}
				{
					position409, tokenIndex409, depth409 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l410
					}
					position++
					goto l409
				l410:
					position, tokenIndex, depth = position409, tokenIndex409, depth409
					if buffer[position] != rune('S') {
						goto l399
					}
					position++
				}
			l409:
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l412
					}
					position++
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if buffer[position] != rune('I') {
						goto l399
					}
					position++
				}
			l411:
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l414
					}
					position++
					goto l413
				l414:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
					if buffer[position] != rune('N') {
						goto l399
					}
					position++
				}
			l413:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l416
					}
					position++
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if buffer[position] != rune('K') {
						goto l399
					}
					position++
				}
			l415:
				if !_rules[rulesp]() {
					goto l399
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l399
				}
				if !_rules[ruleAction14]() {
					goto l399
				}
				depth--
				add(ruleDropSinkStmt, position400)
			}
			return true
		l399:
			position, tokenIndex, depth = position399, tokenIndex399, depth399
			return false
		},
		/* 17 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l420
					}
					position++
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if buffer[position] != rune('D') {
						goto l417
					}
					position++
				}
			l419:
				{
					position421, tokenIndex421, depth421 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l422
					}
					position++
					goto l421
				l422:
					position, tokenIndex, depth = position421, tokenIndex421, depth421
					if buffer[position] != rune('R') {
						goto l417
					}
					position++
				}
			l421:
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l424
					}
					position++
					goto l423
				l424:
					position, tokenIndex, depth = position423, tokenIndex423, depth423
					if buffer[position] != rune('O') {
						goto l417
					}
					position++
				}
			l423:
				{
					position425, tokenIndex425, depth425 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l426
					}
					position++
					goto l425
				l426:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if buffer[position] != rune('P') {
						goto l417
					}
					position++
				}
			l425:
				if !_rules[rulesp]() {
					goto l417
				}
				{
					position427, tokenIndex427, depth427 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l428
					}
					position++
					goto l427
				l428:
					position, tokenIndex, depth = position427, tokenIndex427, depth427
					if buffer[position] != rune('S') {
						goto l417
					}
					position++
				}
			l427:
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l430
					}
					position++
					goto l429
				l430:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if buffer[position] != rune('T') {
						goto l417
					}
					position++
				}
			l429:
				{
					position431, tokenIndex431, depth431 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l432
					}
					position++
					goto l431
				l432:
					position, tokenIndex, depth = position431, tokenIndex431, depth431
					if buffer[position] != rune('A') {
						goto l417
					}
					position++
				}
			l431:
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l434
					}
					position++
					goto l433
				l434:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if buffer[position] != rune('T') {
						goto l417
					}
					position++
				}
			l433:
				{
					position435, tokenIndex435, depth435 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l436
					}
					position++
					goto l435
				l436:
					position, tokenIndex, depth = position435, tokenIndex435, depth435
					if buffer[position] != rune('E') {
						goto l417
					}
					position++
				}
			l435:
				if !_rules[rulesp]() {
					goto l417
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l417
				}
				if !_rules[ruleAction15]() {
					goto l417
				}
				depth--
				add(ruleDropStateStmt, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 18 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action16)> */
		func() bool {
			position437, tokenIndex437, depth437 := position, tokenIndex, depth
			{
				position438 := position
				depth++
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l440
					}
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleDSTREAM]() {
						goto l441
					}
					goto l439
				l441:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if !_rules[ruleRSTREAM]() {
						goto l437
					}
				}
			l439:
				{
					position442 := position
					depth++
					{
						position443, tokenIndex443, depth443 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l443
						}
						if buffer[position] != rune('[') {
							goto l443
						}
						position++
						if !_rules[rulesp]() {
							goto l443
						}
						{
							position445, tokenIndex445, depth445 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l446
							}
							position++
							goto l445
						l446:
							position, tokenIndex, depth = position445, tokenIndex445, depth445
							if buffer[position] != rune('E') {
								goto l443
							}
							position++
						}
					l445:
						{
							position447, tokenIndex447, depth447 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l448
							}
							position++
							goto l447
						l448:
							position, tokenIndex, depth = position447, tokenIndex447, depth447
							if buffer[position] != rune('V') {
								goto l443
							}
							position++
						}
					l447:
						{
							position449, tokenIndex449, depth449 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l450
							}
							position++
							goto l449
						l450:
							position, tokenIndex, depth = position449, tokenIndex449, depth449
							if buffer[position] != rune('E') {
								goto l443
							}
							position++
						}
					l449:
						{
							position451, tokenIndex451, depth451 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l452
							}
							position++
							goto l451
						l452:
							position, tokenIndex, depth = position451, tokenIndex451, depth451
							if buffer[position] != rune('R') {
								goto l443
							}
							position++
						}
					l451:
						{
							position453, tokenIndex453, depth453 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l454
							}
							position++
							goto l453
						l454:
							position, tokenIndex, depth = position453, tokenIndex453, depth453
							if buffer[position] != rune('Y') {
								goto l443
							}
							position++
						}
					l453:
						if !_rules[rulesp]() {
							goto l443
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l443
						}
						if !_rules[rulesp]() {
							goto l443
						}
						if buffer[position] != rune(']') {
							goto l443
						}
						position++
						goto l444
					l443:
						position, tokenIndex, depth = position443, tokenIndex443, depth443
					}
				l444:
					depth--
					add(rulePegText, position442)
				}
				if !_rules[ruleAction16]() {
					goto l437
				}
				depth--
				add(ruleEmitter, position438)
			}
			return true
		l437:
			position, tokenIndex, depth = position437, tokenIndex437, depth437
			return false
		},
		/* 19 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				{
					position457, tokenIndex457, depth457 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l458
					}
				l459:
					{
						position460, tokenIndex460, depth460 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l460
						}
						if buffer[position] != rune(',') {
							goto l460
						}
						position++
						if !_rules[rulesp]() {
							goto l460
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l460
						}
						goto l459
					l460:
						position, tokenIndex, depth = position460, tokenIndex460, depth460
					}
					goto l457
				l458:
					position, tokenIndex, depth = position457, tokenIndex457, depth457
					if !_rules[ruleTimeEmitterInterval]() {
						goto l461
					}
					goto l457
				l461:
					position, tokenIndex, depth = position457, tokenIndex457, depth457
					if !_rules[ruleTupleEmitterInterval]() {
						goto l455
					}
				}
			l457:
				depth--
				add(ruleEmitterIntervals, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
			return false
		},
		/* 20 TimeEmitterInterval <- <(<TimeInterval> Action17)> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				{
					position464 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l462
					}
					depth--
					add(rulePegText, position464)
				}
				if !_rules[ruleAction17]() {
					goto l462
				}
				depth--
				add(ruleTimeEmitterInterval, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 21 TupleEmitterInterval <- <(<TuplesInterval> Action18)> */
		func() bool {
			position465, tokenIndex465, depth465 := position, tokenIndex, depth
			{
				position466 := position
				depth++
				{
					position467 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l465
					}
					depth--
					add(rulePegText, position467)
				}
				if !_rules[ruleAction18]() {
					goto l465
				}
				depth--
				add(ruleTupleEmitterInterval, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 22 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action19)> */
		func() bool {
			position468, tokenIndex468, depth468 := position, tokenIndex, depth
			{
				position469 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l468
				}
				if !_rules[rulesp]() {
					goto l468
				}
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l471
					}
					position++
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if buffer[position] != rune('I') {
						goto l468
					}
					position++
				}
			l470:
				{
					position472, tokenIndex472, depth472 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l473
					}
					position++
					goto l472
				l473:
					position, tokenIndex, depth = position472, tokenIndex472, depth472
					if buffer[position] != rune('N') {
						goto l468
					}
					position++
				}
			l472:
				if !_rules[rulesp]() {
					goto l468
				}
				if !_rules[ruleStream]() {
					goto l468
				}
				if !_rules[ruleAction19]() {
					goto l468
				}
				depth--
				add(ruleTupleEmitterFromInterval, position469)
			}
			return true
		l468:
			position, tokenIndex, depth = position468, tokenIndex468, depth468
			return false
		},
		/* 23 Projections <- <(<(Projection sp (',' sp Projection)*)> Action20)> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				{
					position476 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l474
					}
					if !_rules[rulesp]() {
						goto l474
					}
				l477:
					{
						position478, tokenIndex478, depth478 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l478
						}
						position++
						if !_rules[rulesp]() {
							goto l478
						}
						if !_rules[ruleProjection]() {
							goto l478
						}
						goto l477
					l478:
						position, tokenIndex, depth = position478, tokenIndex478, depth478
					}
					depth--
					add(rulePegText, position476)
				}
				if !_rules[ruleAction20]() {
					goto l474
				}
				depth--
				add(ruleProjections, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 24 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l482
					}
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleExpression]() {
						goto l483
					}
					goto l481
				l483:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if !_rules[ruleWildcard]() {
						goto l479
					}
				}
			l481:
				depth--
				add(ruleProjection, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 25 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action21)> */
		func() bool {
			position484, tokenIndex484, depth484 := position, tokenIndex, depth
			{
				position485 := position
				depth++
				{
					position486, tokenIndex486, depth486 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l487
					}
					goto l486
				l487:
					position, tokenIndex, depth = position486, tokenIndex486, depth486
					if !_rules[ruleWildcard]() {
						goto l484
					}
				}
			l486:
				if !_rules[rulesp]() {
					goto l484
				}
				{
					position488, tokenIndex488, depth488 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l489
					}
					position++
					goto l488
				l489:
					position, tokenIndex, depth = position488, tokenIndex488, depth488
					if buffer[position] != rune('A') {
						goto l484
					}
					position++
				}
			l488:
				{
					position490, tokenIndex490, depth490 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l491
					}
					position++
					goto l490
				l491:
					position, tokenIndex, depth = position490, tokenIndex490, depth490
					if buffer[position] != rune('S') {
						goto l484
					}
					position++
				}
			l490:
				if !_rules[rulesp]() {
					goto l484
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l484
				}
				if !_rules[ruleAction21]() {
					goto l484
				}
				depth--
				add(ruleAliasExpression, position485)
			}
			return true
		l484:
			position, tokenIndex, depth = position484, tokenIndex484, depth484
			return false
		},
		/* 26 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action22)> */
		func() bool {
			position492, tokenIndex492, depth492 := position, tokenIndex, depth
			{
				position493 := position
				depth++
				{
					position494 := position
					depth++
					{
						position495, tokenIndex495, depth495 := position, tokenIndex, depth
						{
							position497, tokenIndex497, depth497 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l498
							}
							position++
							goto l497
						l498:
							position, tokenIndex, depth = position497, tokenIndex497, depth497
							if buffer[position] != rune('F') {
								goto l495
							}
							position++
						}
					l497:
						{
							position499, tokenIndex499, depth499 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l500
							}
							position++
							goto l499
						l500:
							position, tokenIndex, depth = position499, tokenIndex499, depth499
							if buffer[position] != rune('R') {
								goto l495
							}
							position++
						}
					l499:
						{
							position501, tokenIndex501, depth501 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l502
							}
							position++
							goto l501
						l502:
							position, tokenIndex, depth = position501, tokenIndex501, depth501
							if buffer[position] != rune('O') {
								goto l495
							}
							position++
						}
					l501:
						{
							position503, tokenIndex503, depth503 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l504
							}
							position++
							goto l503
						l504:
							position, tokenIndex, depth = position503, tokenIndex503, depth503
							if buffer[position] != rune('M') {
								goto l495
							}
							position++
						}
					l503:
						if !_rules[rulesp]() {
							goto l495
						}
						if !_rules[ruleRelations]() {
							goto l495
						}
						if !_rules[rulesp]() {
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
				if !_rules[ruleAction22]() {
					goto l492
				}
				depth--
				add(ruleWindowedFrom, position493)
			}
			return true
		l492:
			position, tokenIndex, depth = position492, tokenIndex492, depth492
			return false
		},
		/* 27 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position505, tokenIndex505, depth505 := position, tokenIndex, depth
			{
				position506 := position
				depth++
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l508
					}
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if !_rules[ruleTuplesInterval]() {
						goto l505
					}
				}
			l507:
				depth--
				add(ruleInterval, position506)
			}
			return true
		l505:
			position, tokenIndex, depth = position505, tokenIndex505, depth505
			return false
		},
		/* 28 TimeInterval <- <(NumericLiteral sp SECONDS Action23)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l509
				}
				if !_rules[rulesp]() {
					goto l509
				}
				if !_rules[ruleSECONDS]() {
					goto l509
				}
				if !_rules[ruleAction23]() {
					goto l509
				}
				depth--
				add(ruleTimeInterval, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 29 TuplesInterval <- <(NumericLiteral sp TUPLES Action24)> */
		func() bool {
			position511, tokenIndex511, depth511 := position, tokenIndex, depth
			{
				position512 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l511
				}
				if !_rules[rulesp]() {
					goto l511
				}
				if !_rules[ruleTUPLES]() {
					goto l511
				}
				if !_rules[ruleAction24]() {
					goto l511
				}
				depth--
				add(ruleTuplesInterval, position512)
			}
			return true
		l511:
			position, tokenIndex, depth = position511, tokenIndex511, depth511
			return false
		},
		/* 30 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position513, tokenIndex513, depth513 := position, tokenIndex, depth
			{
				position514 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l513
				}
				if !_rules[rulesp]() {
					goto l513
				}
			l515:
				{
					position516, tokenIndex516, depth516 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l516
					}
					position++
					if !_rules[rulesp]() {
						goto l516
					}
					if !_rules[ruleRelationLike]() {
						goto l516
					}
					goto l515
				l516:
					position, tokenIndex, depth = position516, tokenIndex516, depth516
				}
				depth--
				add(ruleRelations, position514)
			}
			return true
		l513:
			position, tokenIndex, depth = position513, tokenIndex513, depth513
			return false
		},
		/* 31 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action25)> */
		func() bool {
			position517, tokenIndex517, depth517 := position, tokenIndex, depth
			{
				position518 := position
				depth++
				{
					position519 := position
					depth++
					{
						position520, tokenIndex520, depth520 := position, tokenIndex, depth
						{
							position522, tokenIndex522, depth522 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l523
							}
							position++
							goto l522
						l523:
							position, tokenIndex, depth = position522, tokenIndex522, depth522
							if buffer[position] != rune('W') {
								goto l520
							}
							position++
						}
					l522:
						{
							position524, tokenIndex524, depth524 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l525
							}
							position++
							goto l524
						l525:
							position, tokenIndex, depth = position524, tokenIndex524, depth524
							if buffer[position] != rune('H') {
								goto l520
							}
							position++
						}
					l524:
						{
							position526, tokenIndex526, depth526 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l527
							}
							position++
							goto l526
						l527:
							position, tokenIndex, depth = position526, tokenIndex526, depth526
							if buffer[position] != rune('E') {
								goto l520
							}
							position++
						}
					l526:
						{
							position528, tokenIndex528, depth528 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l529
							}
							position++
							goto l528
						l529:
							position, tokenIndex, depth = position528, tokenIndex528, depth528
							if buffer[position] != rune('R') {
								goto l520
							}
							position++
						}
					l528:
						{
							position530, tokenIndex530, depth530 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l531
							}
							position++
							goto l530
						l531:
							position, tokenIndex, depth = position530, tokenIndex530, depth530
							if buffer[position] != rune('E') {
								goto l520
							}
							position++
						}
					l530:
						if !_rules[rulesp]() {
							goto l520
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction25]() {
					goto l517
				}
				depth--
				add(ruleFilter, position518)
			}
			return true
		l517:
			position, tokenIndex, depth = position517, tokenIndex517, depth517
			return false
		},
		/* 32 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action26)> */
		func() bool {
			position532, tokenIndex532, depth532 := position, tokenIndex, depth
			{
				position533 := position
				depth++
				{
					position534 := position
					depth++
					{
						position535, tokenIndex535, depth535 := position, tokenIndex, depth
						{
							position537, tokenIndex537, depth537 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l538
							}
							position++
							goto l537
						l538:
							position, tokenIndex, depth = position537, tokenIndex537, depth537
							if buffer[position] != rune('G') {
								goto l535
							}
							position++
						}
					l537:
						{
							position539, tokenIndex539, depth539 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l540
							}
							position++
							goto l539
						l540:
							position, tokenIndex, depth = position539, tokenIndex539, depth539
							if buffer[position] != rune('R') {
								goto l535
							}
							position++
						}
					l539:
						{
							position541, tokenIndex541, depth541 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l542
							}
							position++
							goto l541
						l542:
							position, tokenIndex, depth = position541, tokenIndex541, depth541
							if buffer[position] != rune('O') {
								goto l535
							}
							position++
						}
					l541:
						{
							position543, tokenIndex543, depth543 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l544
							}
							position++
							goto l543
						l544:
							position, tokenIndex, depth = position543, tokenIndex543, depth543
							if buffer[position] != rune('U') {
								goto l535
							}
							position++
						}
					l543:
						{
							position545, tokenIndex545, depth545 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l546
							}
							position++
							goto l545
						l546:
							position, tokenIndex, depth = position545, tokenIndex545, depth545
							if buffer[position] != rune('P') {
								goto l535
							}
							position++
						}
					l545:
						if !_rules[rulesp]() {
							goto l535
						}
						{
							position547, tokenIndex547, depth547 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l548
							}
							position++
							goto l547
						l548:
							position, tokenIndex, depth = position547, tokenIndex547, depth547
							if buffer[position] != rune('B') {
								goto l535
							}
							position++
						}
					l547:
						{
							position549, tokenIndex549, depth549 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l550
							}
							position++
							goto l549
						l550:
							position, tokenIndex, depth = position549, tokenIndex549, depth549
							if buffer[position] != rune('Y') {
								goto l535
							}
							position++
						}
					l549:
						if !_rules[rulesp]() {
							goto l535
						}
						if !_rules[ruleGroupList]() {
							goto l535
						}
						goto l536
					l535:
						position, tokenIndex, depth = position535, tokenIndex535, depth535
					}
				l536:
					depth--
					add(rulePegText, position534)
				}
				if !_rules[ruleAction26]() {
					goto l532
				}
				depth--
				add(ruleGrouping, position533)
			}
			return true
		l532:
			position, tokenIndex, depth = position532, tokenIndex532, depth532
			return false
		},
		/* 33 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position551, tokenIndex551, depth551 := position, tokenIndex, depth
			{
				position552 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l551
				}
				if !_rules[rulesp]() {
					goto l551
				}
			l553:
				{
					position554, tokenIndex554, depth554 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l554
					}
					position++
					if !_rules[rulesp]() {
						goto l554
					}
					if !_rules[ruleExpression]() {
						goto l554
					}
					goto l553
				l554:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
				}
				depth--
				add(ruleGroupList, position552)
			}
			return true
		l551:
			position, tokenIndex, depth = position551, tokenIndex551, depth551
			return false
		},
		/* 34 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action27)> */
		func() bool {
			position555, tokenIndex555, depth555 := position, tokenIndex, depth
			{
				position556 := position
				depth++
				{
					position557 := position
					depth++
					{
						position558, tokenIndex558, depth558 := position, tokenIndex, depth
						{
							position560, tokenIndex560, depth560 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l561
							}
							position++
							goto l560
						l561:
							position, tokenIndex, depth = position560, tokenIndex560, depth560
							if buffer[position] != rune('H') {
								goto l558
							}
							position++
						}
					l560:
						{
							position562, tokenIndex562, depth562 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l563
							}
							position++
							goto l562
						l563:
							position, tokenIndex, depth = position562, tokenIndex562, depth562
							if buffer[position] != rune('A') {
								goto l558
							}
							position++
						}
					l562:
						{
							position564, tokenIndex564, depth564 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l565
							}
							position++
							goto l564
						l565:
							position, tokenIndex, depth = position564, tokenIndex564, depth564
							if buffer[position] != rune('V') {
								goto l558
							}
							position++
						}
					l564:
						{
							position566, tokenIndex566, depth566 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l567
							}
							position++
							goto l566
						l567:
							position, tokenIndex, depth = position566, tokenIndex566, depth566
							if buffer[position] != rune('I') {
								goto l558
							}
							position++
						}
					l566:
						{
							position568, tokenIndex568, depth568 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l569
							}
							position++
							goto l568
						l569:
							position, tokenIndex, depth = position568, tokenIndex568, depth568
							if buffer[position] != rune('N') {
								goto l558
							}
							position++
						}
					l568:
						{
							position570, tokenIndex570, depth570 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l571
							}
							position++
							goto l570
						l571:
							position, tokenIndex, depth = position570, tokenIndex570, depth570
							if buffer[position] != rune('G') {
								goto l558
							}
							position++
						}
					l570:
						if !_rules[rulesp]() {
							goto l558
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction27]() {
					goto l555
				}
				depth--
				add(ruleHaving, position556)
			}
			return true
		l555:
			position, tokenIndex, depth = position555, tokenIndex555, depth555
			return false
		},
		/* 35 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action28))> */
		func() bool {
			position572, tokenIndex572, depth572 := position, tokenIndex, depth
			{
				position573 := position
				depth++
				{
					position574, tokenIndex574, depth574 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l575
					}
					goto l574
				l575:
					position, tokenIndex, depth = position574, tokenIndex574, depth574
					if !_rules[ruleStreamWindow]() {
						goto l572
					}
					if !_rules[ruleAction28]() {
						goto l572
					}
				}
			l574:
				depth--
				add(ruleRelationLike, position573)
			}
			return true
		l572:
			position, tokenIndex, depth = position572, tokenIndex572, depth572
			return false
		},
		/* 36 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action29)> */
		func() bool {
			position576, tokenIndex576, depth576 := position, tokenIndex, depth
			{
				position577 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l576
				}
				if !_rules[rulesp]() {
					goto l576
				}
				{
					position578, tokenIndex578, depth578 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l579
					}
					position++
					goto l578
				l579:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
					if buffer[position] != rune('A') {
						goto l576
					}
					position++
				}
			l578:
				{
					position580, tokenIndex580, depth580 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l581
					}
					position++
					goto l580
				l581:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if buffer[position] != rune('S') {
						goto l576
					}
					position++
				}
			l580:
				if !_rules[rulesp]() {
					goto l576
				}
				if !_rules[ruleIdentifier]() {
					goto l576
				}
				if !_rules[ruleAction29]() {
					goto l576
				}
				depth--
				add(ruleAliasedStreamWindow, position577)
			}
			return true
		l576:
			position, tokenIndex, depth = position576, tokenIndex576, depth576
			return false
		},
		/* 37 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action30)> */
		func() bool {
			position582, tokenIndex582, depth582 := position, tokenIndex, depth
			{
				position583 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l582
				}
				if !_rules[rulesp]() {
					goto l582
				}
				if buffer[position] != rune('[') {
					goto l582
				}
				position++
				if !_rules[rulesp]() {
					goto l582
				}
				{
					position584, tokenIndex584, depth584 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l585
					}
					position++
					goto l584
				l585:
					position, tokenIndex, depth = position584, tokenIndex584, depth584
					if buffer[position] != rune('R') {
						goto l582
					}
					position++
				}
			l584:
				{
					position586, tokenIndex586, depth586 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l587
					}
					position++
					goto l586
				l587:
					position, tokenIndex, depth = position586, tokenIndex586, depth586
					if buffer[position] != rune('A') {
						goto l582
					}
					position++
				}
			l586:
				{
					position588, tokenIndex588, depth588 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l589
					}
					position++
					goto l588
				l589:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					if buffer[position] != rune('N') {
						goto l582
					}
					position++
				}
			l588:
				{
					position590, tokenIndex590, depth590 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l591
					}
					position++
					goto l590
				l591:
					position, tokenIndex, depth = position590, tokenIndex590, depth590
					if buffer[position] != rune('G') {
						goto l582
					}
					position++
				}
			l590:
				{
					position592, tokenIndex592, depth592 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l593
					}
					position++
					goto l592
				l593:
					position, tokenIndex, depth = position592, tokenIndex592, depth592
					if buffer[position] != rune('E') {
						goto l582
					}
					position++
				}
			l592:
				if !_rules[rulesp]() {
					goto l582
				}
				if !_rules[ruleInterval]() {
					goto l582
				}
				if !_rules[rulesp]() {
					goto l582
				}
				if buffer[position] != rune(']') {
					goto l582
				}
				position++
				if !_rules[ruleAction30]() {
					goto l582
				}
				depth--
				add(ruleStreamWindow, position583)
			}
			return true
		l582:
			position, tokenIndex, depth = position582, tokenIndex582, depth582
			return false
		},
		/* 38 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position594, tokenIndex594, depth594 := position, tokenIndex, depth
			{
				position595 := position
				depth++
				{
					position596, tokenIndex596, depth596 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l597
					}
					goto l596
				l597:
					position, tokenIndex, depth = position596, tokenIndex596, depth596
					if !_rules[ruleStream]() {
						goto l594
					}
				}
			l596:
				depth--
				add(ruleStreamLike, position595)
			}
			return true
		l594:
			position, tokenIndex, depth = position594, tokenIndex594, depth594
			return false
		},
		/* 39 UDSFFuncApp <- <(FuncApp Action31)> */
		func() bool {
			position598, tokenIndex598, depth598 := position, tokenIndex, depth
			{
				position599 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l598
				}
				if !_rules[ruleAction31]() {
					goto l598
				}
				depth--
				add(ruleUDSFFuncApp, position599)
			}
			return true
		l598:
			position, tokenIndex, depth = position598, tokenIndex598, depth598
			return false
		},
		/* 40 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action32)> */
		func() bool {
			position600, tokenIndex600, depth600 := position, tokenIndex, depth
			{
				position601 := position
				depth++
				{
					position602 := position
					depth++
					{
						position603, tokenIndex603, depth603 := position, tokenIndex, depth
						{
							position605, tokenIndex605, depth605 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l606
							}
							position++
							goto l605
						l606:
							position, tokenIndex, depth = position605, tokenIndex605, depth605
							if buffer[position] != rune('W') {
								goto l603
							}
							position++
						}
					l605:
						{
							position607, tokenIndex607, depth607 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l608
							}
							position++
							goto l607
						l608:
							position, tokenIndex, depth = position607, tokenIndex607, depth607
							if buffer[position] != rune('I') {
								goto l603
							}
							position++
						}
					l607:
						{
							position609, tokenIndex609, depth609 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l610
							}
							position++
							goto l609
						l610:
							position, tokenIndex, depth = position609, tokenIndex609, depth609
							if buffer[position] != rune('T') {
								goto l603
							}
							position++
						}
					l609:
						{
							position611, tokenIndex611, depth611 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l612
							}
							position++
							goto l611
						l612:
							position, tokenIndex, depth = position611, tokenIndex611, depth611
							if buffer[position] != rune('H') {
								goto l603
							}
							position++
						}
					l611:
						if !_rules[rulesp]() {
							goto l603
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l603
						}
						if !_rules[rulesp]() {
							goto l603
						}
					l613:
						{
							position614, tokenIndex614, depth614 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l614
							}
							position++
							if !_rules[rulesp]() {
								goto l614
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l614
							}
							goto l613
						l614:
							position, tokenIndex, depth = position614, tokenIndex614, depth614
						}
						goto l604
					l603:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
					}
				l604:
					depth--
					add(rulePegText, position602)
				}
				if !_rules[ruleAction32]() {
					goto l600
				}
				depth--
				add(ruleSourceSinkSpecs, position601)
			}
			return true
		l600:
			position, tokenIndex, depth = position600, tokenIndex600, depth600
			return false
		},
		/* 41 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action33)> */
		func() bool {
			position615, tokenIndex615, depth615 := position, tokenIndex, depth
			{
				position616 := position
				depth++
				{
					position617 := position
					depth++
					{
						position618, tokenIndex618, depth618 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l619
						}
						position++
						goto l618
					l619:
						position, tokenIndex, depth = position618, tokenIndex618, depth618
						if buffer[position] != rune('S') {
							goto l615
						}
						position++
					}
				l618:
					{
						position620, tokenIndex620, depth620 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l621
						}
						position++
						goto l620
					l621:
						position, tokenIndex, depth = position620, tokenIndex620, depth620
						if buffer[position] != rune('E') {
							goto l615
						}
						position++
					}
				l620:
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
							goto l615
						}
						position++
					}
				l622:
					if !_rules[rulesp]() {
						goto l615
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l615
					}
					if !_rules[rulesp]() {
						goto l615
					}
				l624:
					{
						position625, tokenIndex625, depth625 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l625
						}
						position++
						if !_rules[rulesp]() {
							goto l625
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l625
						}
						goto l624
					l625:
						position, tokenIndex, depth = position625, tokenIndex625, depth625
					}
					depth--
					add(rulePegText, position617)
				}
				if !_rules[ruleAction33]() {
					goto l615
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position616)
			}
			return true
		l615:
			position, tokenIndex, depth = position615, tokenIndex615, depth615
			return false
		},
		/* 42 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action34)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l626
				}
				if buffer[position] != rune('=') {
					goto l626
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l626
				}
				if !_rules[ruleAction34]() {
					goto l626
				}
				depth--
				add(ruleSourceSinkParam, position627)
			}
			return true
		l626:
			position, tokenIndex, depth = position626, tokenIndex626, depth626
			return false
		},
		/* 43 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position628, tokenIndex628, depth628 := position, tokenIndex, depth
			{
				position629 := position
				depth++
				{
					position630, tokenIndex630, depth630 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l631
					}
					goto l630
				l631:
					position, tokenIndex, depth = position630, tokenIndex630, depth630
					if !_rules[ruleLiteral]() {
						goto l628
					}
				}
			l630:
				depth--
				add(ruleSourceSinkParamVal, position629)
			}
			return true
		l628:
			position, tokenIndex, depth = position628, tokenIndex628, depth628
			return false
		},
		/* 44 PausedOpt <- <(<(Paused / Unpaused)?> Action35)> */
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
						{
							position637, tokenIndex637, depth637 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l638
							}
							goto l637
						l638:
							position, tokenIndex, depth = position637, tokenIndex637, depth637
							if !_rules[ruleUnpaused]() {
								goto l635
							}
						}
					l637:
						goto l636
					l635:
						position, tokenIndex, depth = position635, tokenIndex635, depth635
					}
				l636:
					depth--
					add(rulePegText, position634)
				}
				if !_rules[ruleAction35]() {
					goto l632
				}
				depth--
				add(rulePausedOpt, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 45 Expression <- <orExpr> */
		func() bool {
			position639, tokenIndex639, depth639 := position, tokenIndex, depth
			{
				position640 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l639
				}
				depth--
				add(ruleExpression, position640)
			}
			return true
		l639:
			position, tokenIndex, depth = position639, tokenIndex639, depth639
			return false
		},
		/* 46 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action36)> */
		func() bool {
			position641, tokenIndex641, depth641 := position, tokenIndex, depth
			{
				position642 := position
				depth++
				{
					position643 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l641
					}
					if !_rules[rulesp]() {
						goto l641
					}
					{
						position644, tokenIndex644, depth644 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l644
						}
						if !_rules[rulesp]() {
							goto l644
						}
						if !_rules[ruleandExpr]() {
							goto l644
						}
						goto l645
					l644:
						position, tokenIndex, depth = position644, tokenIndex644, depth644
					}
				l645:
					depth--
					add(rulePegText, position643)
				}
				if !_rules[ruleAction36]() {
					goto l641
				}
				depth--
				add(ruleorExpr, position642)
			}
			return true
		l641:
			position, tokenIndex, depth = position641, tokenIndex641, depth641
			return false
		},
		/* 47 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action37)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				{
					position648 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l646
					}
					if !_rules[rulesp]() {
						goto l646
					}
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l649
						}
						if !_rules[rulesp]() {
							goto l649
						}
						if !_rules[rulenotExpr]() {
							goto l649
						}
						goto l650
					l649:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
					}
				l650:
					depth--
					add(rulePegText, position648)
				}
				if !_rules[ruleAction37]() {
					goto l646
				}
				depth--
				add(ruleandExpr, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 48 notExpr <- <(<((Not sp)? comparisonExpr)> Action38)> */
		func() bool {
			position651, tokenIndex651, depth651 := position, tokenIndex, depth
			{
				position652 := position
				depth++
				{
					position653 := position
					depth++
					{
						position654, tokenIndex654, depth654 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l654
						}
						if !_rules[rulesp]() {
							goto l654
						}
						goto l655
					l654:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
					}
				l655:
					if !_rules[rulecomparisonExpr]() {
						goto l651
					}
					depth--
					add(rulePegText, position653)
				}
				if !_rules[ruleAction38]() {
					goto l651
				}
				depth--
				add(rulenotExpr, position652)
			}
			return true
		l651:
			position, tokenIndex, depth = position651, tokenIndex651, depth651
			return false
		},
		/* 49 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action39)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				{
					position658 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l656
					}
					if !_rules[rulesp]() {
						goto l656
					}
					{
						position659, tokenIndex659, depth659 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l659
						}
						if !_rules[rulesp]() {
							goto l659
						}
						if !_rules[ruleisExpr]() {
							goto l659
						}
						goto l660
					l659:
						position, tokenIndex, depth = position659, tokenIndex659, depth659
					}
				l660:
					depth--
					add(rulePegText, position658)
				}
				if !_rules[ruleAction39]() {
					goto l656
				}
				depth--
				add(rulecomparisonExpr, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 50 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action40)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				{
					position663 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l661
					}
					if !_rules[rulesp]() {
						goto l661
					}
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l664
						}
						if !_rules[rulesp]() {
							goto l664
						}
						if !_rules[ruleNullLiteral]() {
							goto l664
						}
						goto l665
					l664:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
					}
				l665:
					depth--
					add(rulePegText, position663)
				}
				if !_rules[ruleAction40]() {
					goto l661
				}
				depth--
				add(ruleisExpr, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 51 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action41)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				{
					position668 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l666
					}
					if !_rules[rulesp]() {
						goto l666
					}
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l669
						}
						if !_rules[rulesp]() {
							goto l669
						}
						if !_rules[ruleproductExpr]() {
							goto l669
						}
						goto l670
					l669:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
					}
				l670:
					depth--
					add(rulePegText, position668)
				}
				if !_rules[ruleAction41]() {
					goto l666
				}
				depth--
				add(ruletermExpr, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 52 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action42)> */
		func() bool {
			position671, tokenIndex671, depth671 := position, tokenIndex, depth
			{
				position672 := position
				depth++
				{
					position673 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l671
					}
					if !_rules[rulesp]() {
						goto l671
					}
					{
						position674, tokenIndex674, depth674 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l674
						}
						if !_rules[rulesp]() {
							goto l674
						}
						if !_rules[ruleminusExpr]() {
							goto l674
						}
						goto l675
					l674:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
					}
				l675:
					depth--
					add(rulePegText, position673)
				}
				if !_rules[ruleAction42]() {
					goto l671
				}
				depth--
				add(ruleproductExpr, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 53 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action43)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l679
						}
						if !_rules[rulesp]() {
							goto l679
						}
						goto l680
					l679:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
					}
				l680:
					if !_rules[rulebaseExpr]() {
						goto l676
					}
					depth--
					add(rulePegText, position678)
				}
				if !_rules[ruleAction43]() {
					goto l676
				}
				depth--
				add(ruleminusExpr, position677)
			}
			return true
		l676:
			position, tokenIndex, depth = position676, tokenIndex676, depth676
			return false
		},
		/* 54 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position681, tokenIndex681, depth681 := position, tokenIndex, depth
			{
				position682 := position
				depth++
				{
					position683, tokenIndex683, depth683 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l684
					}
					position++
					if !_rules[rulesp]() {
						goto l684
					}
					if !_rules[ruleExpression]() {
						goto l684
					}
					if !_rules[rulesp]() {
						goto l684
					}
					if buffer[position] != rune(')') {
						goto l684
					}
					position++
					goto l683
				l684:
					position, tokenIndex, depth = position683, tokenIndex683, depth683
					if !_rules[ruleBooleanLiteral]() {
						goto l685
					}
					goto l683
				l685:
					position, tokenIndex, depth = position683, tokenIndex683, depth683
					if !_rules[ruleNullLiteral]() {
						goto l686
					}
					goto l683
				l686:
					position, tokenIndex, depth = position683, tokenIndex683, depth683
					if !_rules[ruleFuncApp]() {
						goto l687
					}
					goto l683
				l687:
					position, tokenIndex, depth = position683, tokenIndex683, depth683
					if !_rules[ruleRowMeta]() {
						goto l688
					}
					goto l683
				l688:
					position, tokenIndex, depth = position683, tokenIndex683, depth683
					if !_rules[ruleRowValue]() {
						goto l689
					}
					goto l683
				l689:
					position, tokenIndex, depth = position683, tokenIndex683, depth683
					if !_rules[ruleLiteral]() {
						goto l681
					}
				}
			l683:
				depth--
				add(rulebaseExpr, position682)
			}
			return true
		l681:
			position, tokenIndex, depth = position681, tokenIndex681, depth681
			return false
		},
		/* 55 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action44)> */
		func() bool {
			position690, tokenIndex690, depth690 := position, tokenIndex, depth
			{
				position691 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l690
				}
				if !_rules[rulesp]() {
					goto l690
				}
				if buffer[position] != rune('(') {
					goto l690
				}
				position++
				if !_rules[rulesp]() {
					goto l690
				}
				if !_rules[ruleFuncParams]() {
					goto l690
				}
				if !_rules[rulesp]() {
					goto l690
				}
				if buffer[position] != rune(')') {
					goto l690
				}
				position++
				if !_rules[ruleAction44]() {
					goto l690
				}
				depth--
				add(ruleFuncApp, position691)
			}
			return true
		l690:
			position, tokenIndex, depth = position690, tokenIndex690, depth690
			return false
		},
		/* 56 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action45)> */
		func() bool {
			position692, tokenIndex692, depth692 := position, tokenIndex, depth
			{
				position693 := position
				depth++
				{
					position694 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l692
					}
					if !_rules[rulesp]() {
						goto l692
					}
				l695:
					{
						position696, tokenIndex696, depth696 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l696
						}
						position++
						if !_rules[rulesp]() {
							goto l696
						}
						if !_rules[ruleExpression]() {
							goto l696
						}
						goto l695
					l696:
						position, tokenIndex, depth = position696, tokenIndex696, depth696
					}
					depth--
					add(rulePegText, position694)
				}
				if !_rules[ruleAction45]() {
					goto l692
				}
				depth--
				add(ruleFuncParams, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 57 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position697, tokenIndex697, depth697 := position, tokenIndex, depth
			{
				position698 := position
				depth++
				{
					position699, tokenIndex699, depth699 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l700
					}
					goto l699
				l700:
					position, tokenIndex, depth = position699, tokenIndex699, depth699
					if !_rules[ruleNumericLiteral]() {
						goto l701
					}
					goto l699
				l701:
					position, tokenIndex, depth = position699, tokenIndex699, depth699
					if !_rules[ruleStringLiteral]() {
						goto l697
					}
				}
			l699:
				depth--
				add(ruleLiteral, position698)
			}
			return true
		l697:
			position, tokenIndex, depth = position697, tokenIndex697, depth697
			return false
		},
		/* 58 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position702, tokenIndex702, depth702 := position, tokenIndex, depth
			{
				position703 := position
				depth++
				{
					position704, tokenIndex704, depth704 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l705
					}
					goto l704
				l705:
					position, tokenIndex, depth = position704, tokenIndex704, depth704
					if !_rules[ruleNotEqual]() {
						goto l706
					}
					goto l704
				l706:
					position, tokenIndex, depth = position704, tokenIndex704, depth704
					if !_rules[ruleLessOrEqual]() {
						goto l707
					}
					goto l704
				l707:
					position, tokenIndex, depth = position704, tokenIndex704, depth704
					if !_rules[ruleLess]() {
						goto l708
					}
					goto l704
				l708:
					position, tokenIndex, depth = position704, tokenIndex704, depth704
					if !_rules[ruleGreaterOrEqual]() {
						goto l709
					}
					goto l704
				l709:
					position, tokenIndex, depth = position704, tokenIndex704, depth704
					if !_rules[ruleGreater]() {
						goto l710
					}
					goto l704
				l710:
					position, tokenIndex, depth = position704, tokenIndex704, depth704
					if !_rules[ruleNotEqual]() {
						goto l702
					}
				}
			l704:
				depth--
				add(ruleComparisonOp, position703)
			}
			return true
		l702:
			position, tokenIndex, depth = position702, tokenIndex702, depth702
			return false
		},
		/* 59 IsOp <- <(IsNot / Is)> */
		func() bool {
			position711, tokenIndex711, depth711 := position, tokenIndex, depth
			{
				position712 := position
				depth++
				{
					position713, tokenIndex713, depth713 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l714
					}
					goto l713
				l714:
					position, tokenIndex, depth = position713, tokenIndex713, depth713
					if !_rules[ruleIs]() {
						goto l711
					}
				}
			l713:
				depth--
				add(ruleIsOp, position712)
			}
			return true
		l711:
			position, tokenIndex, depth = position711, tokenIndex711, depth711
			return false
		},
		/* 60 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position715, tokenIndex715, depth715 := position, tokenIndex, depth
			{
				position716 := position
				depth++
				{
					position717, tokenIndex717, depth717 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l718
					}
					goto l717
				l718:
					position, tokenIndex, depth = position717, tokenIndex717, depth717
					if !_rules[ruleMinus]() {
						goto l715
					}
				}
			l717:
				depth--
				add(rulePlusMinusOp, position716)
			}
			return true
		l715:
			position, tokenIndex, depth = position715, tokenIndex715, depth715
			return false
		},
		/* 61 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position719, tokenIndex719, depth719 := position, tokenIndex, depth
			{
				position720 := position
				depth++
				{
					position721, tokenIndex721, depth721 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l722
					}
					goto l721
				l722:
					position, tokenIndex, depth = position721, tokenIndex721, depth721
					if !_rules[ruleDivide]() {
						goto l723
					}
					goto l721
				l723:
					position, tokenIndex, depth = position721, tokenIndex721, depth721
					if !_rules[ruleModulo]() {
						goto l719
					}
				}
			l721:
				depth--
				add(ruleMultDivOp, position720)
			}
			return true
		l719:
			position, tokenIndex, depth = position719, tokenIndex719, depth719
			return false
		},
		/* 62 Stream <- <(<ident> Action46)> */
		func() bool {
			position724, tokenIndex724, depth724 := position, tokenIndex, depth
			{
				position725 := position
				depth++
				{
					position726 := position
					depth++
					if !_rules[ruleident]() {
						goto l724
					}
					depth--
					add(rulePegText, position726)
				}
				if !_rules[ruleAction46]() {
					goto l724
				}
				depth--
				add(ruleStream, position725)
			}
			return true
		l724:
			position, tokenIndex, depth = position724, tokenIndex724, depth724
			return false
		},
		/* 63 RowMeta <- <RowTimestamp> */
		func() bool {
			position727, tokenIndex727, depth727 := position, tokenIndex, depth
			{
				position728 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l727
				}
				depth--
				add(ruleRowMeta, position728)
			}
			return true
		l727:
			position, tokenIndex, depth = position727, tokenIndex727, depth727
			return false
		},
		/* 64 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action47)> */
		func() bool {
			position729, tokenIndex729, depth729 := position, tokenIndex, depth
			{
				position730 := position
				depth++
				{
					position731 := position
					depth++
					{
						position732, tokenIndex732, depth732 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l732
						}
						if buffer[position] != rune(':') {
							goto l732
						}
						position++
						goto l733
					l732:
						position, tokenIndex, depth = position732, tokenIndex732, depth732
					}
				l733:
					if buffer[position] != rune('t') {
						goto l729
					}
					position++
					if buffer[position] != rune('s') {
						goto l729
					}
					position++
					if buffer[position] != rune('(') {
						goto l729
					}
					position++
					if buffer[position] != rune(')') {
						goto l729
					}
					position++
					depth--
					add(rulePegText, position731)
				}
				if !_rules[ruleAction47]() {
					goto l729
				}
				depth--
				add(ruleRowTimestamp, position730)
			}
			return true
		l729:
			position, tokenIndex, depth = position729, tokenIndex729, depth729
			return false
		},
		/* 65 RowValue <- <(<((ident ':')? jsonPath)> Action48)> */
		func() bool {
			position734, tokenIndex734, depth734 := position, tokenIndex, depth
			{
				position735 := position
				depth++
				{
					position736 := position
					depth++
					{
						position737, tokenIndex737, depth737 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l737
						}
						if buffer[position] != rune(':') {
							goto l737
						}
						position++
						goto l738
					l737:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
					}
				l738:
					if !_rules[rulejsonPath]() {
						goto l734
					}
					depth--
					add(rulePegText, position736)
				}
				if !_rules[ruleAction48]() {
					goto l734
				}
				depth--
				add(ruleRowValue, position735)
			}
			return true
		l734:
			position, tokenIndex, depth = position734, tokenIndex734, depth734
			return false
		},
		/* 66 NumericLiteral <- <(<('-'? [0-9]+)> Action49)> */
		func() bool {
			position739, tokenIndex739, depth739 := position, tokenIndex, depth
			{
				position740 := position
				depth++
				{
					position741 := position
					depth++
					{
						position742, tokenIndex742, depth742 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l742
						}
						position++
						goto l743
					l742:
						position, tokenIndex, depth = position742, tokenIndex742, depth742
					}
				l743:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l739
					}
					position++
				l744:
					{
						position745, tokenIndex745, depth745 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l745
						}
						position++
						goto l744
					l745:
						position, tokenIndex, depth = position745, tokenIndex745, depth745
					}
					depth--
					add(rulePegText, position741)
				}
				if !_rules[ruleAction49]() {
					goto l739
				}
				depth--
				add(ruleNumericLiteral, position740)
			}
			return true
		l739:
			position, tokenIndex, depth = position739, tokenIndex739, depth739
			return false
		},
		/* 67 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action50)> */
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
						if buffer[position] != rune('-') {
							goto l749
						}
						position++
						goto l750
					l749:
						position, tokenIndex, depth = position749, tokenIndex749, depth749
					}
				l750:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l746
					}
					position++
				l751:
					{
						position752, tokenIndex752, depth752 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l752
						}
						position++
						goto l751
					l752:
						position, tokenIndex, depth = position752, tokenIndex752, depth752
					}
					if buffer[position] != rune('.') {
						goto l746
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l746
					}
					position++
				l753:
					{
						position754, tokenIndex754, depth754 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l754
						}
						position++
						goto l753
					l754:
						position, tokenIndex, depth = position754, tokenIndex754, depth754
					}
					depth--
					add(rulePegText, position748)
				}
				if !_rules[ruleAction50]() {
					goto l746
				}
				depth--
				add(ruleFloatLiteral, position747)
			}
			return true
		l746:
			position, tokenIndex, depth = position746, tokenIndex746, depth746
			return false
		},
		/* 68 Function <- <(<ident> Action51)> */
		func() bool {
			position755, tokenIndex755, depth755 := position, tokenIndex, depth
			{
				position756 := position
				depth++
				{
					position757 := position
					depth++
					if !_rules[ruleident]() {
						goto l755
					}
					depth--
					add(rulePegText, position757)
				}
				if !_rules[ruleAction51]() {
					goto l755
				}
				depth--
				add(ruleFunction, position756)
			}
			return true
		l755:
			position, tokenIndex, depth = position755, tokenIndex755, depth755
			return false
		},
		/* 69 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action52)> */
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
						if buffer[position] != rune('n') {
							goto l762
						}
						position++
						goto l761
					l762:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if buffer[position] != rune('N') {
							goto l758
						}
						position++
					}
				l761:
					{
						position763, tokenIndex763, depth763 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l764
						}
						position++
						goto l763
					l764:
						position, tokenIndex, depth = position763, tokenIndex763, depth763
						if buffer[position] != rune('U') {
							goto l758
						}
						position++
					}
				l763:
					{
						position765, tokenIndex765, depth765 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l766
						}
						position++
						goto l765
					l766:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
						if buffer[position] != rune('L') {
							goto l758
						}
						position++
					}
				l765:
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('L') {
							goto l758
						}
						position++
					}
				l767:
					depth--
					add(rulePegText, position760)
				}
				if !_rules[ruleAction52]() {
					goto l758
				}
				depth--
				add(ruleNullLiteral, position759)
			}
			return true
		l758:
			position, tokenIndex, depth = position758, tokenIndex758, depth758
			return false
		},
		/* 70 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position769, tokenIndex769, depth769 := position, tokenIndex, depth
			{
				position770 := position
				depth++
				{
					position771, tokenIndex771, depth771 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l772
					}
					goto l771
				l772:
					position, tokenIndex, depth = position771, tokenIndex771, depth771
					if !_rules[ruleFALSE]() {
						goto l769
					}
				}
			l771:
				depth--
				add(ruleBooleanLiteral, position770)
			}
			return true
		l769:
			position, tokenIndex, depth = position769, tokenIndex769, depth769
			return false
		},
		/* 71 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action53)> */
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
						if buffer[position] != rune('t') {
							goto l777
						}
						position++
						goto l776
					l777:
						position, tokenIndex, depth = position776, tokenIndex776, depth776
						if buffer[position] != rune('T') {
							goto l773
						}
						position++
					}
				l776:
					{
						position778, tokenIndex778, depth778 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l779
						}
						position++
						goto l778
					l779:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						if buffer[position] != rune('R') {
							goto l773
						}
						position++
					}
				l778:
					{
						position780, tokenIndex780, depth780 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l781
						}
						position++
						goto l780
					l781:
						position, tokenIndex, depth = position780, tokenIndex780, depth780
						if buffer[position] != rune('U') {
							goto l773
						}
						position++
					}
				l780:
					{
						position782, tokenIndex782, depth782 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l783
						}
						position++
						goto l782
					l783:
						position, tokenIndex, depth = position782, tokenIndex782, depth782
						if buffer[position] != rune('E') {
							goto l773
						}
						position++
					}
				l782:
					depth--
					add(rulePegText, position775)
				}
				if !_rules[ruleAction53]() {
					goto l773
				}
				depth--
				add(ruleTRUE, position774)
			}
			return true
		l773:
			position, tokenIndex, depth = position773, tokenIndex773, depth773
			return false
		},
		/* 72 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action54)> */
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
						if buffer[position] != rune('f') {
							goto l788
						}
						position++
						goto l787
					l788:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
						if buffer[position] != rune('F') {
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
						if buffer[position] != rune('l') {
							goto l792
						}
						position++
						goto l791
					l792:
						position, tokenIndex, depth = position791, tokenIndex791, depth791
						if buffer[position] != rune('L') {
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
					depth--
					add(rulePegText, position786)
				}
				if !_rules[ruleAction54]() {
					goto l784
				}
				depth--
				add(ruleFALSE, position785)
			}
			return true
		l784:
			position, tokenIndex, depth = position784, tokenIndex784, depth784
			return false
		},
		/* 73 Wildcard <- <(<'*'> Action55)> */
		func() bool {
			position797, tokenIndex797, depth797 := position, tokenIndex, depth
			{
				position798 := position
				depth++
				{
					position799 := position
					depth++
					if buffer[position] != rune('*') {
						goto l797
					}
					position++
					depth--
					add(rulePegText, position799)
				}
				if !_rules[ruleAction55]() {
					goto l797
				}
				depth--
				add(ruleWildcard, position798)
			}
			return true
		l797:
			position, tokenIndex, depth = position797, tokenIndex797, depth797
			return false
		},
		/* 74 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action56)> */
		func() bool {
			position800, tokenIndex800, depth800 := position, tokenIndex, depth
			{
				position801 := position
				depth++
				{
					position802 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l800
					}
					position++
				l803:
					{
						position804, tokenIndex804, depth804 := position, tokenIndex, depth
						{
							position805, tokenIndex805, depth805 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l806
							}
							position++
							if buffer[position] != rune('\'') {
								goto l806
							}
							position++
							goto l805
						l806:
							position, tokenIndex, depth = position805, tokenIndex805, depth805
							{
								position807, tokenIndex807, depth807 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l807
								}
								position++
								goto l804
							l807:
								position, tokenIndex, depth = position807, tokenIndex807, depth807
							}
							if !matchDot() {
								goto l804
							}
						}
					l805:
						goto l803
					l804:
						position, tokenIndex, depth = position804, tokenIndex804, depth804
					}
					if buffer[position] != rune('\'') {
						goto l800
					}
					position++
					depth--
					add(rulePegText, position802)
				}
				if !_rules[ruleAction56]() {
					goto l800
				}
				depth--
				add(ruleStringLiteral, position801)
			}
			return true
		l800:
			position, tokenIndex, depth = position800, tokenIndex800, depth800
			return false
		},
		/* 75 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action57)> */
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
						if buffer[position] != rune('i') {
							goto l812
						}
						position++
						goto l811
					l812:
						position, tokenIndex, depth = position811, tokenIndex811, depth811
						if buffer[position] != rune('I') {
							goto l808
						}
						position++
					}
				l811:
					{
						position813, tokenIndex813, depth813 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l814
						}
						position++
						goto l813
					l814:
						position, tokenIndex, depth = position813, tokenIndex813, depth813
						if buffer[position] != rune('S') {
							goto l808
						}
						position++
					}
				l813:
					{
						position815, tokenIndex815, depth815 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l816
						}
						position++
						goto l815
					l816:
						position, tokenIndex, depth = position815, tokenIndex815, depth815
						if buffer[position] != rune('T') {
							goto l808
						}
						position++
					}
				l815:
					{
						position817, tokenIndex817, depth817 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l818
						}
						position++
						goto l817
					l818:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						if buffer[position] != rune('R') {
							goto l808
						}
						position++
					}
				l817:
					{
						position819, tokenIndex819, depth819 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l820
						}
						position++
						goto l819
					l820:
						position, tokenIndex, depth = position819, tokenIndex819, depth819
						if buffer[position] != rune('E') {
							goto l808
						}
						position++
					}
				l819:
					{
						position821, tokenIndex821, depth821 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l822
						}
						position++
						goto l821
					l822:
						position, tokenIndex, depth = position821, tokenIndex821, depth821
						if buffer[position] != rune('A') {
							goto l808
						}
						position++
					}
				l821:
					{
						position823, tokenIndex823, depth823 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l824
						}
						position++
						goto l823
					l824:
						position, tokenIndex, depth = position823, tokenIndex823, depth823
						if buffer[position] != rune('M') {
							goto l808
						}
						position++
					}
				l823:
					depth--
					add(rulePegText, position810)
				}
				if !_rules[ruleAction57]() {
					goto l808
				}
				depth--
				add(ruleISTREAM, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 76 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action58)> */
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
						if buffer[position] != rune('d') {
							goto l829
						}
						position++
						goto l828
					l829:
						position, tokenIndex, depth = position828, tokenIndex828, depth828
						if buffer[position] != rune('D') {
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
					{
						position832, tokenIndex832, depth832 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l833
						}
						position++
						goto l832
					l833:
						position, tokenIndex, depth = position832, tokenIndex832, depth832
						if buffer[position] != rune('T') {
							goto l825
						}
						position++
					}
				l832:
					{
						position834, tokenIndex834, depth834 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l835
						}
						position++
						goto l834
					l835:
						position, tokenIndex, depth = position834, tokenIndex834, depth834
						if buffer[position] != rune('R') {
							goto l825
						}
						position++
					}
				l834:
					{
						position836, tokenIndex836, depth836 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l837
						}
						position++
						goto l836
					l837:
						position, tokenIndex, depth = position836, tokenIndex836, depth836
						if buffer[position] != rune('E') {
							goto l825
						}
						position++
					}
				l836:
					{
						position838, tokenIndex838, depth838 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l839
						}
						position++
						goto l838
					l839:
						position, tokenIndex, depth = position838, tokenIndex838, depth838
						if buffer[position] != rune('A') {
							goto l825
						}
						position++
					}
				l838:
					{
						position840, tokenIndex840, depth840 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l841
						}
						position++
						goto l840
					l841:
						position, tokenIndex, depth = position840, tokenIndex840, depth840
						if buffer[position] != rune('M') {
							goto l825
						}
						position++
					}
				l840:
					depth--
					add(rulePegText, position827)
				}
				if !_rules[ruleAction58]() {
					goto l825
				}
				depth--
				add(ruleDSTREAM, position826)
			}
			return true
		l825:
			position, tokenIndex, depth = position825, tokenIndex825, depth825
			return false
		},
		/* 77 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action59)> */
		func() bool {
			position842, tokenIndex842, depth842 := position, tokenIndex, depth
			{
				position843 := position
				depth++
				{
					position844 := position
					depth++
					{
						position845, tokenIndex845, depth845 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l846
						}
						position++
						goto l845
					l846:
						position, tokenIndex, depth = position845, tokenIndex845, depth845
						if buffer[position] != rune('R') {
							goto l842
						}
						position++
					}
				l845:
					{
						position847, tokenIndex847, depth847 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l848
						}
						position++
						goto l847
					l848:
						position, tokenIndex, depth = position847, tokenIndex847, depth847
						if buffer[position] != rune('S') {
							goto l842
						}
						position++
					}
				l847:
					{
						position849, tokenIndex849, depth849 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l850
						}
						position++
						goto l849
					l850:
						position, tokenIndex, depth = position849, tokenIndex849, depth849
						if buffer[position] != rune('T') {
							goto l842
						}
						position++
					}
				l849:
					{
						position851, tokenIndex851, depth851 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l852
						}
						position++
						goto l851
					l852:
						position, tokenIndex, depth = position851, tokenIndex851, depth851
						if buffer[position] != rune('R') {
							goto l842
						}
						position++
					}
				l851:
					{
						position853, tokenIndex853, depth853 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l854
						}
						position++
						goto l853
					l854:
						position, tokenIndex, depth = position853, tokenIndex853, depth853
						if buffer[position] != rune('E') {
							goto l842
						}
						position++
					}
				l853:
					{
						position855, tokenIndex855, depth855 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l856
						}
						position++
						goto l855
					l856:
						position, tokenIndex, depth = position855, tokenIndex855, depth855
						if buffer[position] != rune('A') {
							goto l842
						}
						position++
					}
				l855:
					{
						position857, tokenIndex857, depth857 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l858
						}
						position++
						goto l857
					l858:
						position, tokenIndex, depth = position857, tokenIndex857, depth857
						if buffer[position] != rune('M') {
							goto l842
						}
						position++
					}
				l857:
					depth--
					add(rulePegText, position844)
				}
				if !_rules[ruleAction59]() {
					goto l842
				}
				depth--
				add(ruleRSTREAM, position843)
			}
			return true
		l842:
			position, tokenIndex, depth = position842, tokenIndex842, depth842
			return false
		},
		/* 78 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action60)> */
		func() bool {
			position859, tokenIndex859, depth859 := position, tokenIndex, depth
			{
				position860 := position
				depth++
				{
					position861 := position
					depth++
					{
						position862, tokenIndex862, depth862 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l863
						}
						position++
						goto l862
					l863:
						position, tokenIndex, depth = position862, tokenIndex862, depth862
						if buffer[position] != rune('T') {
							goto l859
						}
						position++
					}
				l862:
					{
						position864, tokenIndex864, depth864 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l865
						}
						position++
						goto l864
					l865:
						position, tokenIndex, depth = position864, tokenIndex864, depth864
						if buffer[position] != rune('U') {
							goto l859
						}
						position++
					}
				l864:
					{
						position866, tokenIndex866, depth866 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l867
						}
						position++
						goto l866
					l867:
						position, tokenIndex, depth = position866, tokenIndex866, depth866
						if buffer[position] != rune('P') {
							goto l859
						}
						position++
					}
				l866:
					{
						position868, tokenIndex868, depth868 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l869
						}
						position++
						goto l868
					l869:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('L') {
							goto l859
						}
						position++
					}
				l868:
					{
						position870, tokenIndex870, depth870 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l871
						}
						position++
						goto l870
					l871:
						position, tokenIndex, depth = position870, tokenIndex870, depth870
						if buffer[position] != rune('E') {
							goto l859
						}
						position++
					}
				l870:
					{
						position872, tokenIndex872, depth872 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l873
						}
						position++
						goto l872
					l873:
						position, tokenIndex, depth = position872, tokenIndex872, depth872
						if buffer[position] != rune('S') {
							goto l859
						}
						position++
					}
				l872:
					depth--
					add(rulePegText, position861)
				}
				if !_rules[ruleAction60]() {
					goto l859
				}
				depth--
				add(ruleTUPLES, position860)
			}
			return true
		l859:
			position, tokenIndex, depth = position859, tokenIndex859, depth859
			return false
		},
		/* 79 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action61)> */
		func() bool {
			position874, tokenIndex874, depth874 := position, tokenIndex, depth
			{
				position875 := position
				depth++
				{
					position876 := position
					depth++
					{
						position877, tokenIndex877, depth877 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l878
						}
						position++
						goto l877
					l878:
						position, tokenIndex, depth = position877, tokenIndex877, depth877
						if buffer[position] != rune('S') {
							goto l874
						}
						position++
					}
				l877:
					{
						position879, tokenIndex879, depth879 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l880
						}
						position++
						goto l879
					l880:
						position, tokenIndex, depth = position879, tokenIndex879, depth879
						if buffer[position] != rune('E') {
							goto l874
						}
						position++
					}
				l879:
					{
						position881, tokenIndex881, depth881 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l882
						}
						position++
						goto l881
					l882:
						position, tokenIndex, depth = position881, tokenIndex881, depth881
						if buffer[position] != rune('C') {
							goto l874
						}
						position++
					}
				l881:
					{
						position883, tokenIndex883, depth883 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l884
						}
						position++
						goto l883
					l884:
						position, tokenIndex, depth = position883, tokenIndex883, depth883
						if buffer[position] != rune('O') {
							goto l874
						}
						position++
					}
				l883:
					{
						position885, tokenIndex885, depth885 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l886
						}
						position++
						goto l885
					l886:
						position, tokenIndex, depth = position885, tokenIndex885, depth885
						if buffer[position] != rune('N') {
							goto l874
						}
						position++
					}
				l885:
					{
						position887, tokenIndex887, depth887 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l888
						}
						position++
						goto l887
					l888:
						position, tokenIndex, depth = position887, tokenIndex887, depth887
						if buffer[position] != rune('D') {
							goto l874
						}
						position++
					}
				l887:
					{
						position889, tokenIndex889, depth889 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l890
						}
						position++
						goto l889
					l890:
						position, tokenIndex, depth = position889, tokenIndex889, depth889
						if buffer[position] != rune('S') {
							goto l874
						}
						position++
					}
				l889:
					depth--
					add(rulePegText, position876)
				}
				if !_rules[ruleAction61]() {
					goto l874
				}
				depth--
				add(ruleSECONDS, position875)
			}
			return true
		l874:
			position, tokenIndex, depth = position874, tokenIndex874, depth874
			return false
		},
		/* 80 StreamIdentifier <- <(<ident> Action62)> */
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
				if !_rules[ruleAction62]() {
					goto l891
				}
				depth--
				add(ruleStreamIdentifier, position892)
			}
			return true
		l891:
			position, tokenIndex, depth = position891, tokenIndex891, depth891
			return false
		},
		/* 81 SourceSinkType <- <(<ident> Action63)> */
		func() bool {
			position894, tokenIndex894, depth894 := position, tokenIndex, depth
			{
				position895 := position
				depth++
				{
					position896 := position
					depth++
					if !_rules[ruleident]() {
						goto l894
					}
					depth--
					add(rulePegText, position896)
				}
				if !_rules[ruleAction63]() {
					goto l894
				}
				depth--
				add(ruleSourceSinkType, position895)
			}
			return true
		l894:
			position, tokenIndex, depth = position894, tokenIndex894, depth894
			return false
		},
		/* 82 SourceSinkParamKey <- <(<ident> Action64)> */
		func() bool {
			position897, tokenIndex897, depth897 := position, tokenIndex, depth
			{
				position898 := position
				depth++
				{
					position899 := position
					depth++
					if !_rules[ruleident]() {
						goto l897
					}
					depth--
					add(rulePegText, position899)
				}
				if !_rules[ruleAction64]() {
					goto l897
				}
				depth--
				add(ruleSourceSinkParamKey, position898)
			}
			return true
		l897:
			position, tokenIndex, depth = position897, tokenIndex897, depth897
			return false
		},
		/* 83 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action65)> */
		func() bool {
			position900, tokenIndex900, depth900 := position, tokenIndex, depth
			{
				position901 := position
				depth++
				{
					position902 := position
					depth++
					{
						position903, tokenIndex903, depth903 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l904
						}
						position++
						goto l903
					l904:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if buffer[position] != rune('P') {
							goto l900
						}
						position++
					}
				l903:
					{
						position905, tokenIndex905, depth905 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l906
						}
						position++
						goto l905
					l906:
						position, tokenIndex, depth = position905, tokenIndex905, depth905
						if buffer[position] != rune('A') {
							goto l900
						}
						position++
					}
				l905:
					{
						position907, tokenIndex907, depth907 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l908
						}
						position++
						goto l907
					l908:
						position, tokenIndex, depth = position907, tokenIndex907, depth907
						if buffer[position] != rune('U') {
							goto l900
						}
						position++
					}
				l907:
					{
						position909, tokenIndex909, depth909 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l910
						}
						position++
						goto l909
					l910:
						position, tokenIndex, depth = position909, tokenIndex909, depth909
						if buffer[position] != rune('S') {
							goto l900
						}
						position++
					}
				l909:
					{
						position911, tokenIndex911, depth911 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l912
						}
						position++
						goto l911
					l912:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('E') {
							goto l900
						}
						position++
					}
				l911:
					{
						position913, tokenIndex913, depth913 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l914
						}
						position++
						goto l913
					l914:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('D') {
							goto l900
						}
						position++
					}
				l913:
					depth--
					add(rulePegText, position902)
				}
				if !_rules[ruleAction65]() {
					goto l900
				}
				depth--
				add(rulePaused, position901)
			}
			return true
		l900:
			position, tokenIndex, depth = position900, tokenIndex900, depth900
			return false
		},
		/* 84 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action66)> */
		func() bool {
			position915, tokenIndex915, depth915 := position, tokenIndex, depth
			{
				position916 := position
				depth++
				{
					position917 := position
					depth++
					{
						position918, tokenIndex918, depth918 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l919
						}
						position++
						goto l918
					l919:
						position, tokenIndex, depth = position918, tokenIndex918, depth918
						if buffer[position] != rune('U') {
							goto l915
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
							goto l915
						}
						position++
					}
				l920:
					{
						position922, tokenIndex922, depth922 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l923
						}
						position++
						goto l922
					l923:
						position, tokenIndex, depth = position922, tokenIndex922, depth922
						if buffer[position] != rune('P') {
							goto l915
						}
						position++
					}
				l922:
					{
						position924, tokenIndex924, depth924 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l925
						}
						position++
						goto l924
					l925:
						position, tokenIndex, depth = position924, tokenIndex924, depth924
						if buffer[position] != rune('A') {
							goto l915
						}
						position++
					}
				l924:
					{
						position926, tokenIndex926, depth926 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('U') {
							goto l915
						}
						position++
					}
				l926:
					{
						position928, tokenIndex928, depth928 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l929
						}
						position++
						goto l928
					l929:
						position, tokenIndex, depth = position928, tokenIndex928, depth928
						if buffer[position] != rune('S') {
							goto l915
						}
						position++
					}
				l928:
					{
						position930, tokenIndex930, depth930 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l931
						}
						position++
						goto l930
					l931:
						position, tokenIndex, depth = position930, tokenIndex930, depth930
						if buffer[position] != rune('E') {
							goto l915
						}
						position++
					}
				l930:
					{
						position932, tokenIndex932, depth932 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l933
						}
						position++
						goto l932
					l933:
						position, tokenIndex, depth = position932, tokenIndex932, depth932
						if buffer[position] != rune('D') {
							goto l915
						}
						position++
					}
				l932:
					depth--
					add(rulePegText, position917)
				}
				if !_rules[ruleAction66]() {
					goto l915
				}
				depth--
				add(ruleUnpaused, position916)
			}
			return true
		l915:
			position, tokenIndex, depth = position915, tokenIndex915, depth915
			return false
		},
		/* 85 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action67)> */
		func() bool {
			position934, tokenIndex934, depth934 := position, tokenIndex, depth
			{
				position935 := position
				depth++
				{
					position936 := position
					depth++
					{
						position937, tokenIndex937, depth937 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l938
						}
						position++
						goto l937
					l938:
						position, tokenIndex, depth = position937, tokenIndex937, depth937
						if buffer[position] != rune('O') {
							goto l934
						}
						position++
					}
				l937:
					{
						position939, tokenIndex939, depth939 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l940
						}
						position++
						goto l939
					l940:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						if buffer[position] != rune('R') {
							goto l934
						}
						position++
					}
				l939:
					depth--
					add(rulePegText, position936)
				}
				if !_rules[ruleAction67]() {
					goto l934
				}
				depth--
				add(ruleOr, position935)
			}
			return true
		l934:
			position, tokenIndex, depth = position934, tokenIndex934, depth934
			return false
		},
		/* 86 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action68)> */
		func() bool {
			position941, tokenIndex941, depth941 := position, tokenIndex, depth
			{
				position942 := position
				depth++
				{
					position943 := position
					depth++
					{
						position944, tokenIndex944, depth944 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l945
						}
						position++
						goto l944
					l945:
						position, tokenIndex, depth = position944, tokenIndex944, depth944
						if buffer[position] != rune('A') {
							goto l941
						}
						position++
					}
				l944:
					{
						position946, tokenIndex946, depth946 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l947
						}
						position++
						goto l946
					l947:
						position, tokenIndex, depth = position946, tokenIndex946, depth946
						if buffer[position] != rune('N') {
							goto l941
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
							goto l941
						}
						position++
					}
				l948:
					depth--
					add(rulePegText, position943)
				}
				if !_rules[ruleAction68]() {
					goto l941
				}
				depth--
				add(ruleAnd, position942)
			}
			return true
		l941:
			position, tokenIndex, depth = position941, tokenIndex941, depth941
			return false
		},
		/* 87 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action69)> */
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
						if buffer[position] != rune('n') {
							goto l954
						}
						position++
						goto l953
					l954:
						position, tokenIndex, depth = position953, tokenIndex953, depth953
						if buffer[position] != rune('N') {
							goto l950
						}
						position++
					}
				l953:
					{
						position955, tokenIndex955, depth955 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l956
						}
						position++
						goto l955
					l956:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('O') {
							goto l950
						}
						position++
					}
				l955:
					{
						position957, tokenIndex957, depth957 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l958
						}
						position++
						goto l957
					l958:
						position, tokenIndex, depth = position957, tokenIndex957, depth957
						if buffer[position] != rune('T') {
							goto l950
						}
						position++
					}
				l957:
					depth--
					add(rulePegText, position952)
				}
				if !_rules[ruleAction69]() {
					goto l950
				}
				depth--
				add(ruleNot, position951)
			}
			return true
		l950:
			position, tokenIndex, depth = position950, tokenIndex950, depth950
			return false
		},
		/* 88 Equal <- <(<'='> Action70)> */
		func() bool {
			position959, tokenIndex959, depth959 := position, tokenIndex, depth
			{
				position960 := position
				depth++
				{
					position961 := position
					depth++
					if buffer[position] != rune('=') {
						goto l959
					}
					position++
					depth--
					add(rulePegText, position961)
				}
				if !_rules[ruleAction70]() {
					goto l959
				}
				depth--
				add(ruleEqual, position960)
			}
			return true
		l959:
			position, tokenIndex, depth = position959, tokenIndex959, depth959
			return false
		},
		/* 89 Less <- <(<'<'> Action71)> */
		func() bool {
			position962, tokenIndex962, depth962 := position, tokenIndex, depth
			{
				position963 := position
				depth++
				{
					position964 := position
					depth++
					if buffer[position] != rune('<') {
						goto l962
					}
					position++
					depth--
					add(rulePegText, position964)
				}
				if !_rules[ruleAction71]() {
					goto l962
				}
				depth--
				add(ruleLess, position963)
			}
			return true
		l962:
			position, tokenIndex, depth = position962, tokenIndex962, depth962
			return false
		},
		/* 90 LessOrEqual <- <(<('<' '=')> Action72)> */
		func() bool {
			position965, tokenIndex965, depth965 := position, tokenIndex, depth
			{
				position966 := position
				depth++
				{
					position967 := position
					depth++
					if buffer[position] != rune('<') {
						goto l965
					}
					position++
					if buffer[position] != rune('=') {
						goto l965
					}
					position++
					depth--
					add(rulePegText, position967)
				}
				if !_rules[ruleAction72]() {
					goto l965
				}
				depth--
				add(ruleLessOrEqual, position966)
			}
			return true
		l965:
			position, tokenIndex, depth = position965, tokenIndex965, depth965
			return false
		},
		/* 91 Greater <- <(<'>'> Action73)> */
		func() bool {
			position968, tokenIndex968, depth968 := position, tokenIndex, depth
			{
				position969 := position
				depth++
				{
					position970 := position
					depth++
					if buffer[position] != rune('>') {
						goto l968
					}
					position++
					depth--
					add(rulePegText, position970)
				}
				if !_rules[ruleAction73]() {
					goto l968
				}
				depth--
				add(ruleGreater, position969)
			}
			return true
		l968:
			position, tokenIndex, depth = position968, tokenIndex968, depth968
			return false
		},
		/* 92 GreaterOrEqual <- <(<('>' '=')> Action74)> */
		func() bool {
			position971, tokenIndex971, depth971 := position, tokenIndex, depth
			{
				position972 := position
				depth++
				{
					position973 := position
					depth++
					if buffer[position] != rune('>') {
						goto l971
					}
					position++
					if buffer[position] != rune('=') {
						goto l971
					}
					position++
					depth--
					add(rulePegText, position973)
				}
				if !_rules[ruleAction74]() {
					goto l971
				}
				depth--
				add(ruleGreaterOrEqual, position972)
			}
			return true
		l971:
			position, tokenIndex, depth = position971, tokenIndex971, depth971
			return false
		},
		/* 93 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action75)> */
		func() bool {
			position974, tokenIndex974, depth974 := position, tokenIndex, depth
			{
				position975 := position
				depth++
				{
					position976 := position
					depth++
					{
						position977, tokenIndex977, depth977 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l978
						}
						position++
						if buffer[position] != rune('=') {
							goto l978
						}
						position++
						goto l977
					l978:
						position, tokenIndex, depth = position977, tokenIndex977, depth977
						if buffer[position] != rune('<') {
							goto l974
						}
						position++
						if buffer[position] != rune('>') {
							goto l974
						}
						position++
					}
				l977:
					depth--
					add(rulePegText, position976)
				}
				if !_rules[ruleAction75]() {
					goto l974
				}
				depth--
				add(ruleNotEqual, position975)
			}
			return true
		l974:
			position, tokenIndex, depth = position974, tokenIndex974, depth974
			return false
		},
		/* 94 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action76)> */
		func() bool {
			position979, tokenIndex979, depth979 := position, tokenIndex, depth
			{
				position980 := position
				depth++
				{
					position981 := position
					depth++
					{
						position982, tokenIndex982, depth982 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l983
						}
						position++
						goto l982
					l983:
						position, tokenIndex, depth = position982, tokenIndex982, depth982
						if buffer[position] != rune('I') {
							goto l979
						}
						position++
					}
				l982:
					{
						position984, tokenIndex984, depth984 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l985
						}
						position++
						goto l984
					l985:
						position, tokenIndex, depth = position984, tokenIndex984, depth984
						if buffer[position] != rune('S') {
							goto l979
						}
						position++
					}
				l984:
					depth--
					add(rulePegText, position981)
				}
				if !_rules[ruleAction76]() {
					goto l979
				}
				depth--
				add(ruleIs, position980)
			}
			return true
		l979:
			position, tokenIndex, depth = position979, tokenIndex979, depth979
			return false
		},
		/* 95 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action77)> */
		func() bool {
			position986, tokenIndex986, depth986 := position, tokenIndex, depth
			{
				position987 := position
				depth++
				{
					position988 := position
					depth++
					{
						position989, tokenIndex989, depth989 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l990
						}
						position++
						goto l989
					l990:
						position, tokenIndex, depth = position989, tokenIndex989, depth989
						if buffer[position] != rune('I') {
							goto l986
						}
						position++
					}
				l989:
					{
						position991, tokenIndex991, depth991 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l992
						}
						position++
						goto l991
					l992:
						position, tokenIndex, depth = position991, tokenIndex991, depth991
						if buffer[position] != rune('S') {
							goto l986
						}
						position++
					}
				l991:
					if !_rules[rulesp]() {
						goto l986
					}
					{
						position993, tokenIndex993, depth993 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l994
						}
						position++
						goto l993
					l994:
						position, tokenIndex, depth = position993, tokenIndex993, depth993
						if buffer[position] != rune('N') {
							goto l986
						}
						position++
					}
				l993:
					{
						position995, tokenIndex995, depth995 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l996
						}
						position++
						goto l995
					l996:
						position, tokenIndex, depth = position995, tokenIndex995, depth995
						if buffer[position] != rune('O') {
							goto l986
						}
						position++
					}
				l995:
					{
						position997, tokenIndex997, depth997 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l998
						}
						position++
						goto l997
					l998:
						position, tokenIndex, depth = position997, tokenIndex997, depth997
						if buffer[position] != rune('T') {
							goto l986
						}
						position++
					}
				l997:
					depth--
					add(rulePegText, position988)
				}
				if !_rules[ruleAction77]() {
					goto l986
				}
				depth--
				add(ruleIsNot, position987)
			}
			return true
		l986:
			position, tokenIndex, depth = position986, tokenIndex986, depth986
			return false
		},
		/* 96 Plus <- <(<'+'> Action78)> */
		func() bool {
			position999, tokenIndex999, depth999 := position, tokenIndex, depth
			{
				position1000 := position
				depth++
				{
					position1001 := position
					depth++
					if buffer[position] != rune('+') {
						goto l999
					}
					position++
					depth--
					add(rulePegText, position1001)
				}
				if !_rules[ruleAction78]() {
					goto l999
				}
				depth--
				add(rulePlus, position1000)
			}
			return true
		l999:
			position, tokenIndex, depth = position999, tokenIndex999, depth999
			return false
		},
		/* 97 Minus <- <(<'-'> Action79)> */
		func() bool {
			position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
			{
				position1003 := position
				depth++
				{
					position1004 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1002
					}
					position++
					depth--
					add(rulePegText, position1004)
				}
				if !_rules[ruleAction79]() {
					goto l1002
				}
				depth--
				add(ruleMinus, position1003)
			}
			return true
		l1002:
			position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
			return false
		},
		/* 98 Multiply <- <(<'*'> Action80)> */
		func() bool {
			position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
			{
				position1006 := position
				depth++
				{
					position1007 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1005
					}
					position++
					depth--
					add(rulePegText, position1007)
				}
				if !_rules[ruleAction80]() {
					goto l1005
				}
				depth--
				add(ruleMultiply, position1006)
			}
			return true
		l1005:
			position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
			return false
		},
		/* 99 Divide <- <(<'/'> Action81)> */
		func() bool {
			position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
			{
				position1009 := position
				depth++
				{
					position1010 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1008
					}
					position++
					depth--
					add(rulePegText, position1010)
				}
				if !_rules[ruleAction81]() {
					goto l1008
				}
				depth--
				add(ruleDivide, position1009)
			}
			return true
		l1008:
			position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
			return false
		},
		/* 100 Modulo <- <(<'%'> Action82)> */
		func() bool {
			position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
			{
				position1012 := position
				depth++
				{
					position1013 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1011
					}
					position++
					depth--
					add(rulePegText, position1013)
				}
				if !_rules[ruleAction82]() {
					goto l1011
				}
				depth--
				add(ruleModulo, position1012)
			}
			return true
		l1011:
			position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
			return false
		},
		/* 101 UnaryMinus <- <(<'-'> Action83)> */
		func() bool {
			position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
			{
				position1015 := position
				depth++
				{
					position1016 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1014
					}
					position++
					depth--
					add(rulePegText, position1016)
				}
				if !_rules[ruleAction83]() {
					goto l1014
				}
				depth--
				add(ruleUnaryMinus, position1015)
			}
			return true
		l1014:
			position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
			return false
		},
		/* 102 Identifier <- <(<ident> Action84)> */
		func() bool {
			position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
			{
				position1018 := position
				depth++
				{
					position1019 := position
					depth++
					if !_rules[ruleident]() {
						goto l1017
					}
					depth--
					add(rulePegText, position1019)
				}
				if !_rules[ruleAction84]() {
					goto l1017
				}
				depth--
				add(ruleIdentifier, position1018)
			}
			return true
		l1017:
			position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
			return false
		},
		/* 103 TargetIdentifier <- <(<jsonPath> Action85)> */
		func() bool {
			position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
			{
				position1021 := position
				depth++
				{
					position1022 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1020
					}
					depth--
					add(rulePegText, position1022)
				}
				if !_rules[ruleAction85]() {
					goto l1020
				}
				depth--
				add(ruleTargetIdentifier, position1021)
			}
			return true
		l1020:
			position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
			return false
		},
		/* 104 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
			{
				position1024 := position
				depth++
				{
					position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1026
					}
					position++
					goto l1025
				l1026:
					position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1023
					}
					position++
				}
			l1025:
			l1027:
				{
					position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
					{
						position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1030
						}
						position++
						goto l1029
					l1030:
						position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1031
						}
						position++
						goto l1029
					l1031:
						position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1032
						}
						position++
						goto l1029
					l1032:
						position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
						if buffer[position] != rune('_') {
							goto l1028
						}
						position++
					}
				l1029:
					goto l1027
				l1028:
					position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
				}
				depth--
				add(ruleident, position1024)
			}
			return true
		l1023:
			position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
			return false
		},
		/* 105 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
			{
				position1034 := position
				depth++
				{
					position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1036
					}
					position++
					goto l1035
				l1036:
					position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1033
					}
					position++
				}
			l1035:
			l1037:
				{
					position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
					{
						position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1040
						}
						position++
						goto l1039
					l1040:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1041
						}
						position++
						goto l1039
					l1041:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1042
						}
						position++
						goto l1039
					l1042:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if buffer[position] != rune('_') {
							goto l1043
						}
						position++
						goto l1039
					l1043:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if buffer[position] != rune('.') {
							goto l1044
						}
						position++
						goto l1039
					l1044:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if buffer[position] != rune('[') {
							goto l1045
						}
						position++
						goto l1039
					l1045:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if buffer[position] != rune(']') {
							goto l1046
						}
						position++
						goto l1039
					l1046:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if buffer[position] != rune('"') {
							goto l1038
						}
						position++
					}
				l1039:
					goto l1037
				l1038:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
				}
				depth--
				add(rulejsonPath, position1034)
			}
			return true
		l1033:
			position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
			return false
		},
		/* 106 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1048 := position
				depth++
			l1049:
				{
					position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
					{
						position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1052
						}
						position++
						goto l1051
					l1052:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if buffer[position] != rune('\t') {
							goto l1053
						}
						position++
						goto l1051
					l1053:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if buffer[position] != rune('\n') {
							goto l1054
						}
						position++
						goto l1051
					l1054:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if buffer[position] != rune('\r') {
							goto l1055
						}
						position++
						goto l1051
					l1055:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if !_rules[rulecomment]() {
							goto l1050
						}
					}
				l1051:
					goto l1049
				l1050:
					position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
				}
				depth--
				add(rulesp, position1048)
			}
			return true
		},
		/* 107 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
			{
				position1057 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1056
				}
				position++
				if buffer[position] != rune('-') {
					goto l1056
				}
				position++
			l1058:
				{
					position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
					{
						position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
						{
							position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1062
							}
							position++
							goto l1061
						l1062:
							position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
							if buffer[position] != rune('\n') {
								goto l1060
							}
							position++
						}
					l1061:
						goto l1059
					l1060:
						position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
					}
					if !matchDot() {
						goto l1059
					}
					goto l1058
				l1059:
					position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
				}
				{
					position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1064
					}
					position++
					goto l1063
				l1064:
					position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
					if buffer[position] != rune('\n') {
						goto l1056
					}
					position++
				}
			l1063:
				depth--
				add(rulecomment, position1057)
			}
			return true
		l1056:
			position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
			return false
		},
		/* 109 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 110 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 111 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 112 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 113 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 114 Action5 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 115 Action6 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 116 Action7 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 117 Action8 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 118 Action9 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 119 Action10 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 120 Action11 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 121 Action12 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 122 Action13 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 123 Action14 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 124 Action15 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		nil,
		/* 126 Action16 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 127 Action17 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 128 Action18 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 129 Action19 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 130 Action20 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 131 Action21 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 132 Action22 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 133 Action23 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 134 Action24 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 135 Action25 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 136 Action26 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 137 Action27 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 138 Action28 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 139 Action29 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 140 Action30 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 141 Action31 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 142 Action32 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 143 Action33 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 144 Action34 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 145 Action35 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 146 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 147 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 148 Action38 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 149 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 150 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 151 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 152 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 153 Action43 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 154 Action44 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 155 Action45 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 156 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 157 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 158 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 159 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 160 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 161 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 162 Action52 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 163 Action53 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 164 Action54 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 165 Action55 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 166 Action56 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 167 Action57 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 168 Action58 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 169 Action59 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 170 Action60 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 171 Action61 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 172 Action62 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 173 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 174 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 175 Action65 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 176 Action66 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 177 Action67 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 178 Action68 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 179 Action69 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 180 Action70 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 181 Action71 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 182 Action72 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 183 Action73 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 184 Action74 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 185 Action75 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 186 Action76 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 187 Action77 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 188 Action78 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 189 Action79 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 190 Action80 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 191 Action81 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 192 Action82 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 193 Action83 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 194 Action84 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 195 Action85 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
	}
	p.rules = _rules
}

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
	rules  [198]func() bool
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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt / DropSourceStmt / DropStreamStmt / DropSinkStmt / DropStateStmt / UpdateStateStmt / UpdateSourceStmt / UpdateSinkStmt)> */
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
					if !_rules[ruleUpdateSourceStmt]() {
						goto l25
					}
					goto l9
				l25:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleUpdateSinkStmt]() {
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
			position26, tokenIndex26, depth26 := position, tokenIndex, depth
			{
				position27 := position
				depth++
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l29
					}
					position++
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if buffer[position] != rune('S') {
						goto l26
					}
					position++
				}
			l28:
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l31
					}
					position++
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if buffer[position] != rune('E') {
						goto l26
					}
					position++
				}
			l30:
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l33
					}
					position++
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if buffer[position] != rune('L') {
						goto l26
					}
					position++
				}
			l32:
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l35
					}
					position++
					goto l34
				l35:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if buffer[position] != rune('E') {
						goto l26
					}
					position++
				}
			l34:
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
						goto l26
					}
					position++
				}
			l36:
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l39
					}
					position++
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if buffer[position] != rune('T') {
						goto l26
					}
					position++
				}
			l38:
				if !_rules[rulesp]() {
					goto l26
				}
				if !_rules[ruleEmitter]() {
					goto l26
				}
				if !_rules[rulesp]() {
					goto l26
				}
				if !_rules[ruleProjections]() {
					goto l26
				}
				if !_rules[rulesp]() {
					goto l26
				}
				if !_rules[ruleWindowedFrom]() {
					goto l26
				}
				if !_rules[rulesp]() {
					goto l26
				}
				if !_rules[ruleFilter]() {
					goto l26
				}
				if !_rules[rulesp]() {
					goto l26
				}
				if !_rules[ruleGrouping]() {
					goto l26
				}
				if !_rules[rulesp]() {
					goto l26
				}
				if !_rules[ruleHaving]() {
					goto l26
				}
				if !_rules[rulesp]() {
					goto l26
				}
				if !_rules[ruleAction0]() {
					goto l26
				}
				depth--
				add(ruleSelectStmt, position27)
			}
			return true
		l26:
			position, tokenIndex, depth = position26, tokenIndex26, depth26
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position40, tokenIndex40, depth40 := position, tokenIndex, depth
			{
				position41 := position
				depth++
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('C') {
						goto l40
					}
					position++
				}
			l42:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l45
					}
					position++
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if buffer[position] != rune('R') {
						goto l40
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
						goto l40
					}
					position++
				}
			l46:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('A') {
						goto l40
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
						goto l40
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
						goto l40
					}
					position++
				}
			l52:
				if !_rules[rulesp]() {
					goto l40
				}
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('S') {
						goto l40
					}
					position++
				}
			l54:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('T') {
						goto l40
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('R') {
						goto l40
					}
					position++
				}
			l58:
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l61
					}
					position++
					goto l60
				l61:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
					if buffer[position] != rune('E') {
						goto l40
					}
					position++
				}
			l60:
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
					if buffer[position] != rune('A') {
						goto l40
					}
					position++
				}
			l62:
				{
					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l65
					}
					position++
					goto l64
				l65:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
					if buffer[position] != rune('M') {
						goto l40
					}
					position++
				}
			l64:
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				{
					position66, tokenIndex66, depth66 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex, depth = position66, tokenIndex66, depth66
					if buffer[position] != rune('A') {
						goto l40
					}
					position++
				}
			l66:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
					if buffer[position] != rune('S') {
						goto l40
					}
					position++
				}
			l68:
				if !_rules[rulesp]() {
					goto l40
				}
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l71
					}
					position++
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if buffer[position] != rune('S') {
						goto l40
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
						goto l40
					}
					position++
				}
			l72:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if buffer[position] != rune('L') {
						goto l40
					}
					position++
				}
			l74:
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('E') {
						goto l40
					}
					position++
				}
			l76:
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
						goto l40
					}
					position++
				}
			l78:
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('T') {
						goto l40
					}
					position++
				}
			l80:
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleEmitter]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleProjections]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleWindowedFrom]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleFilter]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleGrouping]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleHaving]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleAction1]() {
					goto l40
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position41)
			}
			return true
		l40:
			position, tokenIndex, depth = position40, tokenIndex40, depth40
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position82, tokenIndex82, depth82 := position, tokenIndex, depth
			{
				position83 := position
				depth++
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
					if buffer[position] != rune('C') {
						goto l82
					}
					position++
				}
			l84:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if buffer[position] != rune('R') {
						goto l82
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
						goto l82
					}
					position++
				}
			l88:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('A') {
						goto l82
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
						goto l82
					}
					position++
				}
			l92:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('E') {
						goto l82
					}
					position++
				}
			l94:
				if !_rules[rulesp]() {
					goto l82
				}
				if !_rules[rulePausedOpt]() {
					goto l82
				}
				if !_rules[rulesp]() {
					goto l82
				}
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('S') {
						goto l82
					}
					position++
				}
			l96:
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('O') {
						goto l82
					}
					position++
				}
			l98:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('U') {
						goto l82
					}
					position++
				}
			l100:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('R') {
						goto l82
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('C') {
						goto l82
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
						goto l82
					}
					position++
				}
			l106:
				if !_rules[rulesp]() {
					goto l82
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l82
				}
				if !_rules[rulesp]() {
					goto l82
				}
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('T') {
						goto l82
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('Y') {
						goto l82
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('P') {
						goto l82
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
						goto l82
					}
					position++
				}
			l114:
				if !_rules[rulesp]() {
					goto l82
				}
				if !_rules[ruleSourceSinkType]() {
					goto l82
				}
				if !_rules[rulesp]() {
					goto l82
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l82
				}
				if !_rules[ruleAction2]() {
					goto l82
				}
				depth--
				add(ruleCreateSourceStmt, position83)
			}
			return true
		l82:
			position, tokenIndex, depth = position82, tokenIndex82, depth82
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position116, tokenIndex116, depth116 := position, tokenIndex, depth
			{
				position117 := position
				depth++
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('C') {
						goto l116
					}
					position++
				}
			l118:
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
					if buffer[position] != rune('R') {
						goto l116
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
						goto l116
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('A') {
						goto l116
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('T') {
						goto l116
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('E') {
						goto l116
					}
					position++
				}
			l128:
				if !_rules[rulesp]() {
					goto l116
				}
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('S') {
						goto l116
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('I') {
						goto l116
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('N') {
						goto l116
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('K') {
						goto l116
					}
					position++
				}
			l136:
				if !_rules[rulesp]() {
					goto l116
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l116
				}
				if !_rules[rulesp]() {
					goto l116
				}
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
						goto l116
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('Y') {
						goto l116
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('P') {
						goto l116
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
						goto l116
					}
					position++
				}
			l144:
				if !_rules[rulesp]() {
					goto l116
				}
				if !_rules[ruleSourceSinkType]() {
					goto l116
				}
				if !_rules[rulesp]() {
					goto l116
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l116
				}
				if !_rules[ruleAction3]() {
					goto l116
				}
				depth--
				add(ruleCreateSinkStmt, position117)
			}
			return true
		l116:
			position, tokenIndex, depth = position116, tokenIndex116, depth116
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position146, tokenIndex146, depth146 := position, tokenIndex, depth
			{
				position147 := position
				depth++
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('C') {
						goto l146
					}
					position++
				}
			l148:
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('R') {
						goto l146
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
						goto l146
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('A') {
						goto l146
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
						goto l146
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
						goto l146
					}
					position++
				}
			l158:
				if !_rules[rulesp]() {
					goto l146
				}
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('S') {
						goto l146
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('T') {
						goto l146
					}
					position++
				}
			l162:
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('A') {
						goto l146
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('T') {
						goto l146
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('E') {
						goto l146
					}
					position++
				}
			l168:
				if !_rules[rulesp]() {
					goto l146
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l146
				}
				if !_rules[rulesp]() {
					goto l146
				}
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('T') {
						goto l146
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('Y') {
						goto l146
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('P') {
						goto l146
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('E') {
						goto l146
					}
					position++
				}
			l176:
				if !_rules[rulesp]() {
					goto l146
				}
				if !_rules[ruleSourceSinkType]() {
					goto l146
				}
				if !_rules[rulesp]() {
					goto l146
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l146
				}
				if !_rules[ruleAction4]() {
					goto l146
				}
				depth--
				add(ruleCreateStateStmt, position147)
			}
			return true
		l146:
			position, tokenIndex, depth = position146, tokenIndex146, depth146
			return false
		},
		/* 7 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action5)> */
		func() bool {
			position178, tokenIndex178, depth178 := position, tokenIndex, depth
			{
				position179 := position
				depth++
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('U') {
						goto l178
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('P') {
						goto l178
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('D') {
						goto l178
					}
					position++
				}
			l184:
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('A') {
						goto l178
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('T') {
						goto l178
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
						goto l178
					}
					position++
				}
			l190:
				if !_rules[rulesp]() {
					goto l178
				}
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('S') {
						goto l178
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
						goto l178
					}
					position++
				}
			l194:
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('A') {
						goto l178
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('T') {
						goto l178
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('E') {
						goto l178
					}
					position++
				}
			l200:
				if !_rules[rulesp]() {
					goto l178
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l178
				}
				if !_rules[rulesp]() {
					goto l178
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l178
				}
				if !_rules[ruleAction5]() {
					goto l178
				}
				depth--
				add(ruleUpdateStateStmt, position179)
			}
			return true
		l178:
			position, tokenIndex, depth = position178, tokenIndex178, depth178
			return false
		},
		/* 8 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action6)> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('U') {
						goto l202
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('P') {
						goto l202
					}
					position++
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('D') {
						goto l202
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('A') {
						goto l202
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
						goto l202
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('E') {
						goto l202
					}
					position++
				}
			l214:
				if !_rules[rulesp]() {
					goto l202
				}
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('S') {
						goto l202
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('O') {
						goto l202
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('U') {
						goto l202
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('R') {
						goto l202
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('C') {
						goto l202
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
						goto l202
					}
					position++
				}
			l226:
				if !_rules[rulesp]() {
					goto l202
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l202
				}
				if !_rules[rulesp]() {
					goto l202
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l202
				}
				if !_rules[ruleAction6]() {
					goto l202
				}
				depth--
				add(ruleUpdateSourceStmt, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 9 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action7)> */
		func() bool {
			position228, tokenIndex228, depth228 := position, tokenIndex, depth
			{
				position229 := position
				depth++
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
						goto l228
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('P') {
						goto l228
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('D') {
						goto l228
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('A') {
						goto l228
					}
					position++
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('T') {
						goto l228
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('E') {
						goto l228
					}
					position++
				}
			l240:
				if !_rules[rulesp]() {
					goto l228
				}
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('S') {
						goto l228
					}
					position++
				}
			l242:
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
						goto l228
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
						goto l228
					}
					position++
				}
			l246:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('K') {
						goto l228
					}
					position++
				}
			l248:
				if !_rules[rulesp]() {
					goto l228
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l228
				}
				if !_rules[rulesp]() {
					goto l228
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l228
				}
				if !_rules[ruleAction7]() {
					goto l228
				}
				depth--
				add(ruleUpdateSinkStmt, position229)
			}
			return true
		l228:
			position, tokenIndex, depth = position228, tokenIndex228, depth228
			return false
		},
		/* 10 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action8)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
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
						goto l250
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
						goto l250
					}
					position++
				}
			l254:
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l257
					}
					position++
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if buffer[position] != rune('S') {
						goto l250
					}
					position++
				}
			l256:
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if buffer[position] != rune('E') {
						goto l250
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('R') {
						goto l250
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('T') {
						goto l250
					}
					position++
				}
			l262:
				if !_rules[rulesp]() {
					goto l250
				}
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
						goto l250
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
						goto l250
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('T') {
						goto l250
					}
					position++
				}
			l268:
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if buffer[position] != rune('O') {
						goto l250
					}
					position++
				}
			l270:
				if !_rules[rulesp]() {
					goto l250
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l250
				}
				if !_rules[rulesp]() {
					goto l250
				}
				if !_rules[ruleSelectStmt]() {
					goto l250
				}
				if !_rules[ruleAction8]() {
					goto l250
				}
				depth--
				add(ruleInsertIntoSelectStmt, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 11 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action9)> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l275
					}
					position++
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('I') {
						goto l272
					}
					position++
				}
			l274:
				{
					position276, tokenIndex276, depth276 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l277
					}
					position++
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					if buffer[position] != rune('N') {
						goto l272
					}
					position++
				}
			l276:
				{
					position278, tokenIndex278, depth278 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l279
					}
					position++
					goto l278
				l279:
					position, tokenIndex, depth = position278, tokenIndex278, depth278
					if buffer[position] != rune('S') {
						goto l272
					}
					position++
				}
			l278:
				{
					position280, tokenIndex280, depth280 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l281
					}
					position++
					goto l280
				l281:
					position, tokenIndex, depth = position280, tokenIndex280, depth280
					if buffer[position] != rune('E') {
						goto l272
					}
					position++
				}
			l280:
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('R') {
						goto l272
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
					if buffer[position] != rune('T') {
						goto l272
					}
					position++
				}
			l284:
				if !_rules[rulesp]() {
					goto l272
				}
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
						goto l272
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
						goto l272
					}
					position++
				}
			l288:
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l291
					}
					position++
					goto l290
				l291:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
					if buffer[position] != rune('T') {
						goto l272
					}
					position++
				}
			l290:
				{
					position292, tokenIndex292, depth292 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l293
					}
					position++
					goto l292
				l293:
					position, tokenIndex, depth = position292, tokenIndex292, depth292
					if buffer[position] != rune('O') {
						goto l272
					}
					position++
				}
			l292:
				if !_rules[rulesp]() {
					goto l272
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l272
				}
				if !_rules[rulesp]() {
					goto l272
				}
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l295
					}
					position++
					goto l294
				l295:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if buffer[position] != rune('F') {
						goto l272
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if buffer[position] != rune('R') {
						goto l272
					}
					position++
				}
			l296:
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					if buffer[position] != rune('O') {
						goto l272
					}
					position++
				}
			l298:
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l301
					}
					position++
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if buffer[position] != rune('M') {
						goto l272
					}
					position++
				}
			l300:
				if !_rules[rulesp]() {
					goto l272
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l272
				}
				if !_rules[ruleAction9]() {
					goto l272
				}
				depth--
				add(ruleInsertIntoFromStmt, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 12 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position302, tokenIndex302, depth302 := position, tokenIndex, depth
			{
				position303 := position
				depth++
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
					if buffer[position] != rune('P') {
						goto l302
					}
					position++
				}
			l304:
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l307
					}
					position++
					goto l306
				l307:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					if buffer[position] != rune('A') {
						goto l302
					}
					position++
				}
			l306:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					if buffer[position] != rune('U') {
						goto l302
					}
					position++
				}
			l308:
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
						goto l302
					}
					position++
				}
			l310:
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('E') {
						goto l302
					}
					position++
				}
			l312:
				if !_rules[rulesp]() {
					goto l302
				}
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l315
					}
					position++
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if buffer[position] != rune('S') {
						goto l302
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
						goto l302
					}
					position++
				}
			l316:
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
					if buffer[position] != rune('U') {
						goto l302
					}
					position++
				}
			l318:
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('R') {
						goto l302
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if buffer[position] != rune('C') {
						goto l302
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
						goto l302
					}
					position++
				}
			l324:
				if !_rules[rulesp]() {
					goto l302
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l302
				}
				if !_rules[ruleAction10]() {
					goto l302
				}
				depth--
				add(rulePauseSourceStmt, position303)
			}
			return true
		l302:
			position, tokenIndex, depth = position302, tokenIndex302, depth302
			return false
		},
		/* 13 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action11)> */
		func() bool {
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				{
					position328, tokenIndex328, depth328 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l329
					}
					position++
					goto l328
				l329:
					position, tokenIndex, depth = position328, tokenIndex328, depth328
					if buffer[position] != rune('R') {
						goto l326
					}
					position++
				}
			l328:
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
						goto l326
					}
					position++
				}
			l330:
				{
					position332, tokenIndex332, depth332 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l333
					}
					position++
					goto l332
				l333:
					position, tokenIndex, depth = position332, tokenIndex332, depth332
					if buffer[position] != rune('S') {
						goto l326
					}
					position++
				}
			l332:
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if buffer[position] != rune('U') {
						goto l326
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
						goto l326
					}
					position++
				}
			l336:
				{
					position338, tokenIndex338, depth338 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l339
					}
					position++
					goto l338
				l339:
					position, tokenIndex, depth = position338, tokenIndex338, depth338
					if buffer[position] != rune('E') {
						goto l326
					}
					position++
				}
			l338:
				if !_rules[rulesp]() {
					goto l326
				}
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l341
					}
					position++
					goto l340
				l341:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
					if buffer[position] != rune('S') {
						goto l326
					}
					position++
				}
			l340:
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if buffer[position] != rune('O') {
						goto l326
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
					if buffer[position] != rune('U') {
						goto l326
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					if buffer[position] != rune('R') {
						goto l326
					}
					position++
				}
			l346:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if buffer[position] != rune('C') {
						goto l326
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
						goto l326
					}
					position++
				}
			l350:
				if !_rules[rulesp]() {
					goto l326
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l326
				}
				if !_rules[ruleAction11]() {
					goto l326
				}
				depth--
				add(ruleResumeSourceStmt, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 14 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action12)> */
		func() bool {
			position352, tokenIndex352, depth352 := position, tokenIndex, depth
			{
				position353 := position
				depth++
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
					if buffer[position] != rune('R') {
						goto l352
					}
					position++
				}
			l354:
				{
					position356, tokenIndex356, depth356 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l357
					}
					position++
					goto l356
				l357:
					position, tokenIndex, depth = position356, tokenIndex356, depth356
					if buffer[position] != rune('E') {
						goto l352
					}
					position++
				}
			l356:
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l359
					}
					position++
					goto l358
				l359:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					if buffer[position] != rune('W') {
						goto l352
					}
					position++
				}
			l358:
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if buffer[position] != rune('I') {
						goto l352
					}
					position++
				}
			l360:
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l363
					}
					position++
					goto l362
				l363:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
					if buffer[position] != rune('N') {
						goto l352
					}
					position++
				}
			l362:
				{
					position364, tokenIndex364, depth364 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l365
					}
					position++
					goto l364
				l365:
					position, tokenIndex, depth = position364, tokenIndex364, depth364
					if buffer[position] != rune('D') {
						goto l352
					}
					position++
				}
			l364:
				if !_rules[rulesp]() {
					goto l352
				}
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l367
					}
					position++
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if buffer[position] != rune('S') {
						goto l352
					}
					position++
				}
			l366:
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if buffer[position] != rune('O') {
						goto l352
					}
					position++
				}
			l368:
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l371
					}
					position++
					goto l370
				l371:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
					if buffer[position] != rune('U') {
						goto l352
					}
					position++
				}
			l370:
				{
					position372, tokenIndex372, depth372 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l373
					}
					position++
					goto l372
				l373:
					position, tokenIndex, depth = position372, tokenIndex372, depth372
					if buffer[position] != rune('R') {
						goto l352
					}
					position++
				}
			l372:
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('C') {
						goto l352
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
						goto l352
					}
					position++
				}
			l376:
				if !_rules[rulesp]() {
					goto l352
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l352
				}
				if !_rules[ruleAction12]() {
					goto l352
				}
				depth--
				add(ruleRewindSourceStmt, position353)
			}
			return true
		l352:
			position, tokenIndex, depth = position352, tokenIndex352, depth352
			return false
		},
		/* 15 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('D') {
						goto l378
					}
					position++
				}
			l380:
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l383
					}
					position++
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if buffer[position] != rune('R') {
						goto l378
					}
					position++
				}
			l382:
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l385
					}
					position++
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if buffer[position] != rune('O') {
						goto l378
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('P') {
						goto l378
					}
					position++
				}
			l386:
				if !_rules[rulesp]() {
					goto l378
				}
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l389
					}
					position++
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					if buffer[position] != rune('S') {
						goto l378
					}
					position++
				}
			l388:
				{
					position390, tokenIndex390, depth390 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l391
					}
					position++
					goto l390
				l391:
					position, tokenIndex, depth = position390, tokenIndex390, depth390
					if buffer[position] != rune('O') {
						goto l378
					}
					position++
				}
			l390:
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if buffer[position] != rune('U') {
						goto l378
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
						goto l378
					}
					position++
				}
			l394:
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('C') {
						goto l378
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if buffer[position] != rune('E') {
						goto l378
					}
					position++
				}
			l398:
				if !_rules[rulesp]() {
					goto l378
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l378
				}
				if !_rules[ruleAction13]() {
					goto l378
				}
				depth--
				add(ruleDropSourceStmt, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 16 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action14)> */
		func() bool {
			position400, tokenIndex400, depth400 := position, tokenIndex, depth
			{
				position401 := position
				depth++
				{
					position402, tokenIndex402, depth402 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l403
					}
					position++
					goto l402
				l403:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
					if buffer[position] != rune('D') {
						goto l400
					}
					position++
				}
			l402:
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l405
					}
					position++
					goto l404
				l405:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					if buffer[position] != rune('R') {
						goto l400
					}
					position++
				}
			l404:
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l407
					}
					position++
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if buffer[position] != rune('O') {
						goto l400
					}
					position++
				}
			l406:
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l409
					}
					position++
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if buffer[position] != rune('P') {
						goto l400
					}
					position++
				}
			l408:
				if !_rules[rulesp]() {
					goto l400
				}
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l411
					}
					position++
					goto l410
				l411:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
					if buffer[position] != rune('S') {
						goto l400
					}
					position++
				}
			l410:
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l413
					}
					position++
					goto l412
				l413:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
					if buffer[position] != rune('T') {
						goto l400
					}
					position++
				}
			l412:
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('R') {
						goto l400
					}
					position++
				}
			l414:
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l417
					}
					position++
					goto l416
				l417:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if buffer[position] != rune('E') {
						goto l400
					}
					position++
				}
			l416:
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
					if buffer[position] != rune('A') {
						goto l400
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('M') {
						goto l400
					}
					position++
				}
			l420:
				if !_rules[rulesp]() {
					goto l400
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l400
				}
				if !_rules[ruleAction14]() {
					goto l400
				}
				depth--
				add(ruleDropStreamStmt, position401)
			}
			return true
		l400:
			position, tokenIndex, depth = position400, tokenIndex400, depth400
			return false
		},
		/* 17 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action15)> */
		func() bool {
			position422, tokenIndex422, depth422 := position, tokenIndex, depth
			{
				position423 := position
				depth++
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l425
					}
					position++
					goto l424
				l425:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
					if buffer[position] != rune('D') {
						goto l422
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
						goto l422
					}
					position++
				}
			l426:
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l429
					}
					position++
					goto l428
				l429:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
					if buffer[position] != rune('O') {
						goto l422
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
						goto l422
					}
					position++
				}
			l430:
				if !_rules[rulesp]() {
					goto l422
				}
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('S') {
						goto l422
					}
					position++
				}
			l432:
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l435
					}
					position++
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if buffer[position] != rune('I') {
						goto l422
					}
					position++
				}
			l434:
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l437
					}
					position++
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if buffer[position] != rune('N') {
						goto l422
					}
					position++
				}
			l436:
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if buffer[position] != rune('K') {
						goto l422
					}
					position++
				}
			l438:
				if !_rules[rulesp]() {
					goto l422
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l422
				}
				if !_rules[ruleAction15]() {
					goto l422
				}
				depth--
				add(ruleDropSinkStmt, position423)
			}
			return true
		l422:
			position, tokenIndex, depth = position422, tokenIndex422, depth422
			return false
		},
		/* 18 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				{
					position442, tokenIndex442, depth442 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l443
					}
					position++
					goto l442
				l443:
					position, tokenIndex, depth = position442, tokenIndex442, depth442
					if buffer[position] != rune('D') {
						goto l440
					}
					position++
				}
			l442:
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
						goto l440
					}
					position++
				}
			l444:
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
					if buffer[position] != rune('O') {
						goto l440
					}
					position++
				}
			l446:
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
					if buffer[position] != rune('P') {
						goto l440
					}
					position++
				}
			l448:
				if !_rules[rulesp]() {
					goto l440
				}
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('S') {
						goto l440
					}
					position++
				}
			l450:
				{
					position452, tokenIndex452, depth452 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l453
					}
					position++
					goto l452
				l453:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if buffer[position] != rune('T') {
						goto l440
					}
					position++
				}
			l452:
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l455
					}
					position++
					goto l454
				l455:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if buffer[position] != rune('A') {
						goto l440
					}
					position++
				}
			l454:
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l457
					}
					position++
					goto l456
				l457:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if buffer[position] != rune('T') {
						goto l440
					}
					position++
				}
			l456:
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l459
					}
					position++
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if buffer[position] != rune('E') {
						goto l440
					}
					position++
				}
			l458:
				if !_rules[rulesp]() {
					goto l440
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l440
				}
				if !_rules[ruleAction16]() {
					goto l440
				}
				depth--
				add(ruleDropStateStmt, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 19 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action17)> */
		func() bool {
			position460, tokenIndex460, depth460 := position, tokenIndex, depth
			{
				position461 := position
				depth++
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l463
					}
					goto l462
				l463:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if !_rules[ruleDSTREAM]() {
						goto l464
					}
					goto l462
				l464:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if !_rules[ruleRSTREAM]() {
						goto l460
					}
				}
			l462:
				{
					position465 := position
					depth++
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l466
						}
						if buffer[position] != rune('[') {
							goto l466
						}
						position++
						if !_rules[rulesp]() {
							goto l466
						}
						{
							position468, tokenIndex468, depth468 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l469
							}
							position++
							goto l468
						l469:
							position, tokenIndex, depth = position468, tokenIndex468, depth468
							if buffer[position] != rune('E') {
								goto l466
							}
							position++
						}
					l468:
						{
							position470, tokenIndex470, depth470 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l471
							}
							position++
							goto l470
						l471:
							position, tokenIndex, depth = position470, tokenIndex470, depth470
							if buffer[position] != rune('V') {
								goto l466
							}
							position++
						}
					l470:
						{
							position472, tokenIndex472, depth472 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l473
							}
							position++
							goto l472
						l473:
							position, tokenIndex, depth = position472, tokenIndex472, depth472
							if buffer[position] != rune('E') {
								goto l466
							}
							position++
						}
					l472:
						{
							position474, tokenIndex474, depth474 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l475
							}
							position++
							goto l474
						l475:
							position, tokenIndex, depth = position474, tokenIndex474, depth474
							if buffer[position] != rune('R') {
								goto l466
							}
							position++
						}
					l474:
						{
							position476, tokenIndex476, depth476 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l477
							}
							position++
							goto l476
						l477:
							position, tokenIndex, depth = position476, tokenIndex476, depth476
							if buffer[position] != rune('Y') {
								goto l466
							}
							position++
						}
					l476:
						if !_rules[rulesp]() {
							goto l466
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l466
						}
						if !_rules[rulesp]() {
							goto l466
						}
						if buffer[position] != rune(']') {
							goto l466
						}
						position++
						goto l467
					l466:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
					}
				l467:
					depth--
					add(rulePegText, position465)
				}
				if !_rules[ruleAction17]() {
					goto l460
				}
				depth--
				add(ruleEmitter, position461)
			}
			return true
		l460:
			position, tokenIndex, depth = position460, tokenIndex460, depth460
			return false
		},
		/* 20 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				{
					position480, tokenIndex480, depth480 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l481
					}
				l482:
					{
						position483, tokenIndex483, depth483 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l483
						}
						if buffer[position] != rune(',') {
							goto l483
						}
						position++
						if !_rules[rulesp]() {
							goto l483
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l483
						}
						goto l482
					l483:
						position, tokenIndex, depth = position483, tokenIndex483, depth483
					}
					goto l480
				l481:
					position, tokenIndex, depth = position480, tokenIndex480, depth480
					if !_rules[ruleTimeEmitterInterval]() {
						goto l484
					}
					goto l480
				l484:
					position, tokenIndex, depth = position480, tokenIndex480, depth480
					if !_rules[ruleTupleEmitterInterval]() {
						goto l478
					}
				}
			l480:
				depth--
				add(ruleEmitterIntervals, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 21 TimeEmitterInterval <- <(<TimeInterval> Action18)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				{
					position487 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l485
					}
					depth--
					add(rulePegText, position487)
				}
				if !_rules[ruleAction18]() {
					goto l485
				}
				depth--
				add(ruleTimeEmitterInterval, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 22 TupleEmitterInterval <- <(<TuplesInterval> Action19)> */
		func() bool {
			position488, tokenIndex488, depth488 := position, tokenIndex, depth
			{
				position489 := position
				depth++
				{
					position490 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l488
					}
					depth--
					add(rulePegText, position490)
				}
				if !_rules[ruleAction19]() {
					goto l488
				}
				depth--
				add(ruleTupleEmitterInterval, position489)
			}
			return true
		l488:
			position, tokenIndex, depth = position488, tokenIndex488, depth488
			return false
		},
		/* 23 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action20)> */
		func() bool {
			position491, tokenIndex491, depth491 := position, tokenIndex, depth
			{
				position492 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l491
				}
				if !_rules[rulesp]() {
					goto l491
				}
				{
					position493, tokenIndex493, depth493 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l494
					}
					position++
					goto l493
				l494:
					position, tokenIndex, depth = position493, tokenIndex493, depth493
					if buffer[position] != rune('I') {
						goto l491
					}
					position++
				}
			l493:
				{
					position495, tokenIndex495, depth495 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l496
					}
					position++
					goto l495
				l496:
					position, tokenIndex, depth = position495, tokenIndex495, depth495
					if buffer[position] != rune('N') {
						goto l491
					}
					position++
				}
			l495:
				if !_rules[rulesp]() {
					goto l491
				}
				if !_rules[ruleStream]() {
					goto l491
				}
				if !_rules[ruleAction20]() {
					goto l491
				}
				depth--
				add(ruleTupleEmitterFromInterval, position492)
			}
			return true
		l491:
			position, tokenIndex, depth = position491, tokenIndex491, depth491
			return false
		},
		/* 24 Projections <- <(<(Projection sp (',' sp Projection)*)> Action21)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				{
					position499 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l497
					}
					if !_rules[rulesp]() {
						goto l497
					}
				l500:
					{
						position501, tokenIndex501, depth501 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l501
						}
						position++
						if !_rules[rulesp]() {
							goto l501
						}
						if !_rules[ruleProjection]() {
							goto l501
						}
						goto l500
					l501:
						position, tokenIndex, depth = position501, tokenIndex501, depth501
					}
					depth--
					add(rulePegText, position499)
				}
				if !_rules[ruleAction21]() {
					goto l497
				}
				depth--
				add(ruleProjections, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 25 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l505
					}
					goto l504
				l505:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
					if !_rules[ruleExpression]() {
						goto l506
					}
					goto l504
				l506:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
					if !_rules[ruleWildcard]() {
						goto l502
					}
				}
			l504:
				depth--
				add(ruleProjection, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 26 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action22)> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				{
					position509, tokenIndex509, depth509 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l510
					}
					goto l509
				l510:
					position, tokenIndex, depth = position509, tokenIndex509, depth509
					if !_rules[ruleWildcard]() {
						goto l507
					}
				}
			l509:
				if !_rules[rulesp]() {
					goto l507
				}
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if buffer[position] != rune('A') {
						goto l507
					}
					position++
				}
			l511:
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l514
					}
					position++
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if buffer[position] != rune('S') {
						goto l507
					}
					position++
				}
			l513:
				if !_rules[rulesp]() {
					goto l507
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l507
				}
				if !_rules[ruleAction22]() {
					goto l507
				}
				depth--
				add(ruleAliasExpression, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 27 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action23)> */
		func() bool {
			position515, tokenIndex515, depth515 := position, tokenIndex, depth
			{
				position516 := position
				depth++
				{
					position517 := position
					depth++
					{
						position518, tokenIndex518, depth518 := position, tokenIndex, depth
						{
							position520, tokenIndex520, depth520 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l521
							}
							position++
							goto l520
						l521:
							position, tokenIndex, depth = position520, tokenIndex520, depth520
							if buffer[position] != rune('F') {
								goto l518
							}
							position++
						}
					l520:
						{
							position522, tokenIndex522, depth522 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l523
							}
							position++
							goto l522
						l523:
							position, tokenIndex, depth = position522, tokenIndex522, depth522
							if buffer[position] != rune('R') {
								goto l518
							}
							position++
						}
					l522:
						{
							position524, tokenIndex524, depth524 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l525
							}
							position++
							goto l524
						l525:
							position, tokenIndex, depth = position524, tokenIndex524, depth524
							if buffer[position] != rune('O') {
								goto l518
							}
							position++
						}
					l524:
						{
							position526, tokenIndex526, depth526 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l527
							}
							position++
							goto l526
						l527:
							position, tokenIndex, depth = position526, tokenIndex526, depth526
							if buffer[position] != rune('M') {
								goto l518
							}
							position++
						}
					l526:
						if !_rules[rulesp]() {
							goto l518
						}
						if !_rules[ruleRelations]() {
							goto l518
						}
						if !_rules[rulesp]() {
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
				if !_rules[ruleAction23]() {
					goto l515
				}
				depth--
				add(ruleWindowedFrom, position516)
			}
			return true
		l515:
			position, tokenIndex, depth = position515, tokenIndex515, depth515
			return false
		},
		/* 28 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position528, tokenIndex528, depth528 := position, tokenIndex, depth
			{
				position529 := position
				depth++
				{
					position530, tokenIndex530, depth530 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l531
					}
					goto l530
				l531:
					position, tokenIndex, depth = position530, tokenIndex530, depth530
					if !_rules[ruleTuplesInterval]() {
						goto l528
					}
				}
			l530:
				depth--
				add(ruleInterval, position529)
			}
			return true
		l528:
			position, tokenIndex, depth = position528, tokenIndex528, depth528
			return false
		},
		/* 29 TimeInterval <- <(NumericLiteral sp SECONDS Action24)> */
		func() bool {
			position532, tokenIndex532, depth532 := position, tokenIndex, depth
			{
				position533 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l532
				}
				if !_rules[rulesp]() {
					goto l532
				}
				if !_rules[ruleSECONDS]() {
					goto l532
				}
				if !_rules[ruleAction24]() {
					goto l532
				}
				depth--
				add(ruleTimeInterval, position533)
			}
			return true
		l532:
			position, tokenIndex, depth = position532, tokenIndex532, depth532
			return false
		},
		/* 30 TuplesInterval <- <(NumericLiteral sp TUPLES Action25)> */
		func() bool {
			position534, tokenIndex534, depth534 := position, tokenIndex, depth
			{
				position535 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l534
				}
				if !_rules[rulesp]() {
					goto l534
				}
				if !_rules[ruleTUPLES]() {
					goto l534
				}
				if !_rules[ruleAction25]() {
					goto l534
				}
				depth--
				add(ruleTuplesInterval, position535)
			}
			return true
		l534:
			position, tokenIndex, depth = position534, tokenIndex534, depth534
			return false
		},
		/* 31 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position536, tokenIndex536, depth536 := position, tokenIndex, depth
			{
				position537 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l536
				}
				if !_rules[rulesp]() {
					goto l536
				}
			l538:
				{
					position539, tokenIndex539, depth539 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l539
					}
					position++
					if !_rules[rulesp]() {
						goto l539
					}
					if !_rules[ruleRelationLike]() {
						goto l539
					}
					goto l538
				l539:
					position, tokenIndex, depth = position539, tokenIndex539, depth539
				}
				depth--
				add(ruleRelations, position537)
			}
			return true
		l536:
			position, tokenIndex, depth = position536, tokenIndex536, depth536
			return false
		},
		/* 32 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action26)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					{
						position543, tokenIndex543, depth543 := position, tokenIndex, depth
						{
							position545, tokenIndex545, depth545 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l546
							}
							position++
							goto l545
						l546:
							position, tokenIndex, depth = position545, tokenIndex545, depth545
							if buffer[position] != rune('W') {
								goto l543
							}
							position++
						}
					l545:
						{
							position547, tokenIndex547, depth547 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l548
							}
							position++
							goto l547
						l548:
							position, tokenIndex, depth = position547, tokenIndex547, depth547
							if buffer[position] != rune('H') {
								goto l543
							}
							position++
						}
					l547:
						{
							position549, tokenIndex549, depth549 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l550
							}
							position++
							goto l549
						l550:
							position, tokenIndex, depth = position549, tokenIndex549, depth549
							if buffer[position] != rune('E') {
								goto l543
							}
							position++
						}
					l549:
						{
							position551, tokenIndex551, depth551 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l552
							}
							position++
							goto l551
						l552:
							position, tokenIndex, depth = position551, tokenIndex551, depth551
							if buffer[position] != rune('R') {
								goto l543
							}
							position++
						}
					l551:
						{
							position553, tokenIndex553, depth553 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l554
							}
							position++
							goto l553
						l554:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if buffer[position] != rune('E') {
								goto l543
							}
							position++
						}
					l553:
						if !_rules[rulesp]() {
							goto l543
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction26]() {
					goto l540
				}
				depth--
				add(ruleFilter, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 33 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action27)> */
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
							if buffer[position] != rune('g') {
								goto l561
							}
							position++
							goto l560
						l561:
							position, tokenIndex, depth = position560, tokenIndex560, depth560
							if buffer[position] != rune('G') {
								goto l558
							}
							position++
						}
					l560:
						{
							position562, tokenIndex562, depth562 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l563
							}
							position++
							goto l562
						l563:
							position, tokenIndex, depth = position562, tokenIndex562, depth562
							if buffer[position] != rune('R') {
								goto l558
							}
							position++
						}
					l562:
						{
							position564, tokenIndex564, depth564 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l565
							}
							position++
							goto l564
						l565:
							position, tokenIndex, depth = position564, tokenIndex564, depth564
							if buffer[position] != rune('O') {
								goto l558
							}
							position++
						}
					l564:
						{
							position566, tokenIndex566, depth566 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l567
							}
							position++
							goto l566
						l567:
							position, tokenIndex, depth = position566, tokenIndex566, depth566
							if buffer[position] != rune('U') {
								goto l558
							}
							position++
						}
					l566:
						{
							position568, tokenIndex568, depth568 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l569
							}
							position++
							goto l568
						l569:
							position, tokenIndex, depth = position568, tokenIndex568, depth568
							if buffer[position] != rune('P') {
								goto l558
							}
							position++
						}
					l568:
						if !_rules[rulesp]() {
							goto l558
						}
						{
							position570, tokenIndex570, depth570 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l571
							}
							position++
							goto l570
						l571:
							position, tokenIndex, depth = position570, tokenIndex570, depth570
							if buffer[position] != rune('B') {
								goto l558
							}
							position++
						}
					l570:
						{
							position572, tokenIndex572, depth572 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l573
							}
							position++
							goto l572
						l573:
							position, tokenIndex, depth = position572, tokenIndex572, depth572
							if buffer[position] != rune('Y') {
								goto l558
							}
							position++
						}
					l572:
						if !_rules[rulesp]() {
							goto l558
						}
						if !_rules[ruleGroupList]() {
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
				add(ruleGrouping, position556)
			}
			return true
		l555:
			position, tokenIndex, depth = position555, tokenIndex555, depth555
			return false
		},
		/* 34 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position574, tokenIndex574, depth574 := position, tokenIndex, depth
			{
				position575 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l574
				}
				if !_rules[rulesp]() {
					goto l574
				}
			l576:
				{
					position577, tokenIndex577, depth577 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l577
					}
					position++
					if !_rules[rulesp]() {
						goto l577
					}
					if !_rules[ruleExpression]() {
						goto l577
					}
					goto l576
				l577:
					position, tokenIndex, depth = position577, tokenIndex577, depth577
				}
				depth--
				add(ruleGroupList, position575)
			}
			return true
		l574:
			position, tokenIndex, depth = position574, tokenIndex574, depth574
			return false
		},
		/* 35 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action28)> */
		func() bool {
			position578, tokenIndex578, depth578 := position, tokenIndex, depth
			{
				position579 := position
				depth++
				{
					position580 := position
					depth++
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						{
							position583, tokenIndex583, depth583 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l584
							}
							position++
							goto l583
						l584:
							position, tokenIndex, depth = position583, tokenIndex583, depth583
							if buffer[position] != rune('H') {
								goto l581
							}
							position++
						}
					l583:
						{
							position585, tokenIndex585, depth585 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l586
							}
							position++
							goto l585
						l586:
							position, tokenIndex, depth = position585, tokenIndex585, depth585
							if buffer[position] != rune('A') {
								goto l581
							}
							position++
						}
					l585:
						{
							position587, tokenIndex587, depth587 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l588
							}
							position++
							goto l587
						l588:
							position, tokenIndex, depth = position587, tokenIndex587, depth587
							if buffer[position] != rune('V') {
								goto l581
							}
							position++
						}
					l587:
						{
							position589, tokenIndex589, depth589 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l590
							}
							position++
							goto l589
						l590:
							position, tokenIndex, depth = position589, tokenIndex589, depth589
							if buffer[position] != rune('I') {
								goto l581
							}
							position++
						}
					l589:
						{
							position591, tokenIndex591, depth591 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l592
							}
							position++
							goto l591
						l592:
							position, tokenIndex, depth = position591, tokenIndex591, depth591
							if buffer[position] != rune('N') {
								goto l581
							}
							position++
						}
					l591:
						{
							position593, tokenIndex593, depth593 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l594
							}
							position++
							goto l593
						l594:
							position, tokenIndex, depth = position593, tokenIndex593, depth593
							if buffer[position] != rune('G') {
								goto l581
							}
							position++
						}
					l593:
						if !_rules[rulesp]() {
							goto l581
						}
						if !_rules[ruleExpression]() {
							goto l581
						}
						goto l582
					l581:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
					}
				l582:
					depth--
					add(rulePegText, position580)
				}
				if !_rules[ruleAction28]() {
					goto l578
				}
				depth--
				add(ruleHaving, position579)
			}
			return true
		l578:
			position, tokenIndex, depth = position578, tokenIndex578, depth578
			return false
		},
		/* 36 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action29))> */
		func() bool {
			position595, tokenIndex595, depth595 := position, tokenIndex, depth
			{
				position596 := position
				depth++
				{
					position597, tokenIndex597, depth597 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l598
					}
					goto l597
				l598:
					position, tokenIndex, depth = position597, tokenIndex597, depth597
					if !_rules[ruleStreamWindow]() {
						goto l595
					}
					if !_rules[ruleAction29]() {
						goto l595
					}
				}
			l597:
				depth--
				add(ruleRelationLike, position596)
			}
			return true
		l595:
			position, tokenIndex, depth = position595, tokenIndex595, depth595
			return false
		},
		/* 37 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action30)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l599
				}
				if !_rules[rulesp]() {
					goto l599
				}
				{
					position601, tokenIndex601, depth601 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l602
					}
					position++
					goto l601
				l602:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if buffer[position] != rune('A') {
						goto l599
					}
					position++
				}
			l601:
				{
					position603, tokenIndex603, depth603 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l604
					}
					position++
					goto l603
				l604:
					position, tokenIndex, depth = position603, tokenIndex603, depth603
					if buffer[position] != rune('S') {
						goto l599
					}
					position++
				}
			l603:
				if !_rules[rulesp]() {
					goto l599
				}
				if !_rules[ruleIdentifier]() {
					goto l599
				}
				if !_rules[ruleAction30]() {
					goto l599
				}
				depth--
				add(ruleAliasedStreamWindow, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 38 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action31)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l605
				}
				if !_rules[rulesp]() {
					goto l605
				}
				if buffer[position] != rune('[') {
					goto l605
				}
				position++
				if !_rules[rulesp]() {
					goto l605
				}
				{
					position607, tokenIndex607, depth607 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l608
					}
					position++
					goto l607
				l608:
					position, tokenIndex, depth = position607, tokenIndex607, depth607
					if buffer[position] != rune('R') {
						goto l605
					}
					position++
				}
			l607:
				{
					position609, tokenIndex609, depth609 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l610
					}
					position++
					goto l609
				l610:
					position, tokenIndex, depth = position609, tokenIndex609, depth609
					if buffer[position] != rune('A') {
						goto l605
					}
					position++
				}
			l609:
				{
					position611, tokenIndex611, depth611 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l612
					}
					position++
					goto l611
				l612:
					position, tokenIndex, depth = position611, tokenIndex611, depth611
					if buffer[position] != rune('N') {
						goto l605
					}
					position++
				}
			l611:
				{
					position613, tokenIndex613, depth613 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l614
					}
					position++
					goto l613
				l614:
					position, tokenIndex, depth = position613, tokenIndex613, depth613
					if buffer[position] != rune('G') {
						goto l605
					}
					position++
				}
			l613:
				{
					position615, tokenIndex615, depth615 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l616
					}
					position++
					goto l615
				l616:
					position, tokenIndex, depth = position615, tokenIndex615, depth615
					if buffer[position] != rune('E') {
						goto l605
					}
					position++
				}
			l615:
				if !_rules[rulesp]() {
					goto l605
				}
				if !_rules[ruleInterval]() {
					goto l605
				}
				if !_rules[rulesp]() {
					goto l605
				}
				if buffer[position] != rune(']') {
					goto l605
				}
				position++
				if !_rules[ruleAction31]() {
					goto l605
				}
				depth--
				add(ruleStreamWindow, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 39 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position617, tokenIndex617, depth617 := position, tokenIndex, depth
			{
				position618 := position
				depth++
				{
					position619, tokenIndex619, depth619 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l620
					}
					goto l619
				l620:
					position, tokenIndex, depth = position619, tokenIndex619, depth619
					if !_rules[ruleStream]() {
						goto l617
					}
				}
			l619:
				depth--
				add(ruleStreamLike, position618)
			}
			return true
		l617:
			position, tokenIndex, depth = position617, tokenIndex617, depth617
			return false
		},
		/* 40 UDSFFuncApp <- <(FuncApp Action32)> */
		func() bool {
			position621, tokenIndex621, depth621 := position, tokenIndex, depth
			{
				position622 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l621
				}
				if !_rules[ruleAction32]() {
					goto l621
				}
				depth--
				add(ruleUDSFFuncApp, position622)
			}
			return true
		l621:
			position, tokenIndex, depth = position621, tokenIndex621, depth621
			return false
		},
		/* 41 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action33)> */
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
						{
							position628, tokenIndex628, depth628 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l629
							}
							position++
							goto l628
						l629:
							position, tokenIndex, depth = position628, tokenIndex628, depth628
							if buffer[position] != rune('W') {
								goto l626
							}
							position++
						}
					l628:
						{
							position630, tokenIndex630, depth630 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l631
							}
							position++
							goto l630
						l631:
							position, tokenIndex, depth = position630, tokenIndex630, depth630
							if buffer[position] != rune('I') {
								goto l626
							}
							position++
						}
					l630:
						{
							position632, tokenIndex632, depth632 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l633
							}
							position++
							goto l632
						l633:
							position, tokenIndex, depth = position632, tokenIndex632, depth632
							if buffer[position] != rune('T') {
								goto l626
							}
							position++
						}
					l632:
						{
							position634, tokenIndex634, depth634 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l635
							}
							position++
							goto l634
						l635:
							position, tokenIndex, depth = position634, tokenIndex634, depth634
							if buffer[position] != rune('H') {
								goto l626
							}
							position++
						}
					l634:
						if !_rules[rulesp]() {
							goto l626
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l626
						}
						if !_rules[rulesp]() {
							goto l626
						}
					l636:
						{
							position637, tokenIndex637, depth637 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l637
							}
							position++
							if !_rules[rulesp]() {
								goto l637
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l637
							}
							goto l636
						l637:
							position, tokenIndex, depth = position637, tokenIndex637, depth637
						}
						goto l627
					l626:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
					}
				l627:
					depth--
					add(rulePegText, position625)
				}
				if !_rules[ruleAction33]() {
					goto l623
				}
				depth--
				add(ruleSourceSinkSpecs, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 42 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action34)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					{
						position641, tokenIndex641, depth641 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l642
						}
						position++
						goto l641
					l642:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
						if buffer[position] != rune('S') {
							goto l638
						}
						position++
					}
				l641:
					{
						position643, tokenIndex643, depth643 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l644
						}
						position++
						goto l643
					l644:
						position, tokenIndex, depth = position643, tokenIndex643, depth643
						if buffer[position] != rune('E') {
							goto l638
						}
						position++
					}
				l643:
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l646
						}
						position++
						goto l645
					l646:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
						if buffer[position] != rune('T') {
							goto l638
						}
						position++
					}
				l645:
					if !_rules[rulesp]() {
						goto l638
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l638
					}
					if !_rules[rulesp]() {
						goto l638
					}
				l647:
					{
						position648, tokenIndex648, depth648 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l648
						}
						position++
						if !_rules[rulesp]() {
							goto l648
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l648
						}
						goto l647
					l648:
						position, tokenIndex, depth = position648, tokenIndex648, depth648
					}
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction34]() {
					goto l638
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 43 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action35)> */
		func() bool {
			position649, tokenIndex649, depth649 := position, tokenIndex, depth
			{
				position650 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l649
				}
				if buffer[position] != rune('=') {
					goto l649
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l649
				}
				if !_rules[ruleAction35]() {
					goto l649
				}
				depth--
				add(ruleSourceSinkParam, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 44 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position651, tokenIndex651, depth651 := position, tokenIndex, depth
			{
				position652 := position
				depth++
				{
					position653, tokenIndex653, depth653 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l654
					}
					goto l653
				l654:
					position, tokenIndex, depth = position653, tokenIndex653, depth653
					if !_rules[ruleLiteral]() {
						goto l651
					}
				}
			l653:
				depth--
				add(ruleSourceSinkParamVal, position652)
			}
			return true
		l651:
			position, tokenIndex, depth = position651, tokenIndex651, depth651
			return false
		},
		/* 45 PausedOpt <- <(<(Paused / Unpaused)?> Action36)> */
		func() bool {
			position655, tokenIndex655, depth655 := position, tokenIndex, depth
			{
				position656 := position
				depth++
				{
					position657 := position
					depth++
					{
						position658, tokenIndex658, depth658 := position, tokenIndex, depth
						{
							position660, tokenIndex660, depth660 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l661
							}
							goto l660
						l661:
							position, tokenIndex, depth = position660, tokenIndex660, depth660
							if !_rules[ruleUnpaused]() {
								goto l658
							}
						}
					l660:
						goto l659
					l658:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
					}
				l659:
					depth--
					add(rulePegText, position657)
				}
				if !_rules[ruleAction36]() {
					goto l655
				}
				depth--
				add(rulePausedOpt, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 46 Expression <- <orExpr> */
		func() bool {
			position662, tokenIndex662, depth662 := position, tokenIndex, depth
			{
				position663 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l662
				}
				depth--
				add(ruleExpression, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 47 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action37)> */
		func() bool {
			position664, tokenIndex664, depth664 := position, tokenIndex, depth
			{
				position665 := position
				depth++
				{
					position666 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l664
					}
					if !_rules[rulesp]() {
						goto l664
					}
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l667
						}
						if !_rules[rulesp]() {
							goto l667
						}
						if !_rules[ruleandExpr]() {
							goto l667
						}
						goto l668
					l667:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
					}
				l668:
					depth--
					add(rulePegText, position666)
				}
				if !_rules[ruleAction37]() {
					goto l664
				}
				depth--
				add(ruleorExpr, position665)
			}
			return true
		l664:
			position, tokenIndex, depth = position664, tokenIndex664, depth664
			return false
		},
		/* 48 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action38)> */
		func() bool {
			position669, tokenIndex669, depth669 := position, tokenIndex, depth
			{
				position670 := position
				depth++
				{
					position671 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l669
					}
					if !_rules[rulesp]() {
						goto l669
					}
					{
						position672, tokenIndex672, depth672 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l672
						}
						if !_rules[rulesp]() {
							goto l672
						}
						if !_rules[rulenotExpr]() {
							goto l672
						}
						goto l673
					l672:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
					}
				l673:
					depth--
					add(rulePegText, position671)
				}
				if !_rules[ruleAction38]() {
					goto l669
				}
				depth--
				add(ruleandExpr, position670)
			}
			return true
		l669:
			position, tokenIndex, depth = position669, tokenIndex669, depth669
			return false
		},
		/* 49 notExpr <- <(<((Not sp)? comparisonExpr)> Action39)> */
		func() bool {
			position674, tokenIndex674, depth674 := position, tokenIndex, depth
			{
				position675 := position
				depth++
				{
					position676 := position
					depth++
					{
						position677, tokenIndex677, depth677 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l677
						}
						if !_rules[rulesp]() {
							goto l677
						}
						goto l678
					l677:
						position, tokenIndex, depth = position677, tokenIndex677, depth677
					}
				l678:
					if !_rules[rulecomparisonExpr]() {
						goto l674
					}
					depth--
					add(rulePegText, position676)
				}
				if !_rules[ruleAction39]() {
					goto l674
				}
				depth--
				add(rulenotExpr, position675)
			}
			return true
		l674:
			position, tokenIndex, depth = position674, tokenIndex674, depth674
			return false
		},
		/* 50 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action40)> */
		func() bool {
			position679, tokenIndex679, depth679 := position, tokenIndex, depth
			{
				position680 := position
				depth++
				{
					position681 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l679
					}
					if !_rules[rulesp]() {
						goto l679
					}
					{
						position682, tokenIndex682, depth682 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l682
						}
						if !_rules[rulesp]() {
							goto l682
						}
						if !_rules[ruleisExpr]() {
							goto l682
						}
						goto l683
					l682:
						position, tokenIndex, depth = position682, tokenIndex682, depth682
					}
				l683:
					depth--
					add(rulePegText, position681)
				}
				if !_rules[ruleAction40]() {
					goto l679
				}
				depth--
				add(rulecomparisonExpr, position680)
			}
			return true
		l679:
			position, tokenIndex, depth = position679, tokenIndex679, depth679
			return false
		},
		/* 51 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action41)> */
		func() bool {
			position684, tokenIndex684, depth684 := position, tokenIndex, depth
			{
				position685 := position
				depth++
				{
					position686 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l684
					}
					if !_rules[rulesp]() {
						goto l684
					}
					{
						position687, tokenIndex687, depth687 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l687
						}
						if !_rules[rulesp]() {
							goto l687
						}
						if !_rules[ruleNullLiteral]() {
							goto l687
						}
						goto l688
					l687:
						position, tokenIndex, depth = position687, tokenIndex687, depth687
					}
				l688:
					depth--
					add(rulePegText, position686)
				}
				if !_rules[ruleAction41]() {
					goto l684
				}
				depth--
				add(ruleisExpr, position685)
			}
			return true
		l684:
			position, tokenIndex, depth = position684, tokenIndex684, depth684
			return false
		},
		/* 52 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action42)> */
		func() bool {
			position689, tokenIndex689, depth689 := position, tokenIndex, depth
			{
				position690 := position
				depth++
				{
					position691 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l689
					}
					if !_rules[rulesp]() {
						goto l689
					}
					{
						position692, tokenIndex692, depth692 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l692
						}
						if !_rules[rulesp]() {
							goto l692
						}
						if !_rules[ruleproductExpr]() {
							goto l692
						}
						goto l693
					l692:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
					}
				l693:
					depth--
					add(rulePegText, position691)
				}
				if !_rules[ruleAction42]() {
					goto l689
				}
				depth--
				add(ruletermExpr, position690)
			}
			return true
		l689:
			position, tokenIndex, depth = position689, tokenIndex689, depth689
			return false
		},
		/* 53 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action43)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				{
					position696 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l694
					}
					if !_rules[rulesp]() {
						goto l694
					}
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l697
						}
						if !_rules[rulesp]() {
							goto l697
						}
						if !_rules[ruleminusExpr]() {
							goto l697
						}
						goto l698
					l697:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
					}
				l698:
					depth--
					add(rulePegText, position696)
				}
				if !_rules[ruleAction43]() {
					goto l694
				}
				depth--
				add(ruleproductExpr, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 54 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action44)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l702
						}
						if !_rules[rulesp]() {
							goto l702
						}
						goto l703
					l702:
						position, tokenIndex, depth = position702, tokenIndex702, depth702
					}
				l703:
					if !_rules[rulebaseExpr]() {
						goto l699
					}
					depth--
					add(rulePegText, position701)
				}
				if !_rules[ruleAction44]() {
					goto l699
				}
				depth--
				add(ruleminusExpr, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 55 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position704, tokenIndex704, depth704 := position, tokenIndex, depth
			{
				position705 := position
				depth++
				{
					position706, tokenIndex706, depth706 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l707
					}
					position++
					if !_rules[rulesp]() {
						goto l707
					}
					if !_rules[ruleExpression]() {
						goto l707
					}
					if !_rules[rulesp]() {
						goto l707
					}
					if buffer[position] != rune(')') {
						goto l707
					}
					position++
					goto l706
				l707:
					position, tokenIndex, depth = position706, tokenIndex706, depth706
					if !_rules[ruleBooleanLiteral]() {
						goto l708
					}
					goto l706
				l708:
					position, tokenIndex, depth = position706, tokenIndex706, depth706
					if !_rules[ruleNullLiteral]() {
						goto l709
					}
					goto l706
				l709:
					position, tokenIndex, depth = position706, tokenIndex706, depth706
					if !_rules[ruleFuncApp]() {
						goto l710
					}
					goto l706
				l710:
					position, tokenIndex, depth = position706, tokenIndex706, depth706
					if !_rules[ruleRowMeta]() {
						goto l711
					}
					goto l706
				l711:
					position, tokenIndex, depth = position706, tokenIndex706, depth706
					if !_rules[ruleRowValue]() {
						goto l712
					}
					goto l706
				l712:
					position, tokenIndex, depth = position706, tokenIndex706, depth706
					if !_rules[ruleLiteral]() {
						goto l704
					}
				}
			l706:
				depth--
				add(rulebaseExpr, position705)
			}
			return true
		l704:
			position, tokenIndex, depth = position704, tokenIndex704, depth704
			return false
		},
		/* 56 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action45)> */
		func() bool {
			position713, tokenIndex713, depth713 := position, tokenIndex, depth
			{
				position714 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l713
				}
				if !_rules[rulesp]() {
					goto l713
				}
				if buffer[position] != rune('(') {
					goto l713
				}
				position++
				if !_rules[rulesp]() {
					goto l713
				}
				if !_rules[ruleFuncParams]() {
					goto l713
				}
				if !_rules[rulesp]() {
					goto l713
				}
				if buffer[position] != rune(')') {
					goto l713
				}
				position++
				if !_rules[ruleAction45]() {
					goto l713
				}
				depth--
				add(ruleFuncApp, position714)
			}
			return true
		l713:
			position, tokenIndex, depth = position713, tokenIndex713, depth713
			return false
		},
		/* 57 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action46)> */
		func() bool {
			position715, tokenIndex715, depth715 := position, tokenIndex, depth
			{
				position716 := position
				depth++
				{
					position717 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l715
					}
					if !_rules[rulesp]() {
						goto l715
					}
				l718:
					{
						position719, tokenIndex719, depth719 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l719
						}
						position++
						if !_rules[rulesp]() {
							goto l719
						}
						if !_rules[ruleExpression]() {
							goto l719
						}
						goto l718
					l719:
						position, tokenIndex, depth = position719, tokenIndex719, depth719
					}
					depth--
					add(rulePegText, position717)
				}
				if !_rules[ruleAction46]() {
					goto l715
				}
				depth--
				add(ruleFuncParams, position716)
			}
			return true
		l715:
			position, tokenIndex, depth = position715, tokenIndex715, depth715
			return false
		},
		/* 58 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722, tokenIndex722, depth722 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l723
					}
					goto l722
				l723:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleNumericLiteral]() {
						goto l724
					}
					goto l722
				l724:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleStringLiteral]() {
						goto l720
					}
				}
			l722:
				depth--
				add(ruleLiteral, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 59 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position725, tokenIndex725, depth725 := position, tokenIndex, depth
			{
				position726 := position
				depth++
				{
					position727, tokenIndex727, depth727 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l728
					}
					goto l727
				l728:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
					if !_rules[ruleNotEqual]() {
						goto l729
					}
					goto l727
				l729:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
					if !_rules[ruleLessOrEqual]() {
						goto l730
					}
					goto l727
				l730:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
					if !_rules[ruleLess]() {
						goto l731
					}
					goto l727
				l731:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
					if !_rules[ruleGreaterOrEqual]() {
						goto l732
					}
					goto l727
				l732:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
					if !_rules[ruleGreater]() {
						goto l733
					}
					goto l727
				l733:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
					if !_rules[ruleNotEqual]() {
						goto l725
					}
				}
			l727:
				depth--
				add(ruleComparisonOp, position726)
			}
			return true
		l725:
			position, tokenIndex, depth = position725, tokenIndex725, depth725
			return false
		},
		/* 60 IsOp <- <(IsNot / Is)> */
		func() bool {
			position734, tokenIndex734, depth734 := position, tokenIndex, depth
			{
				position735 := position
				depth++
				{
					position736, tokenIndex736, depth736 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l737
					}
					goto l736
				l737:
					position, tokenIndex, depth = position736, tokenIndex736, depth736
					if !_rules[ruleIs]() {
						goto l734
					}
				}
			l736:
				depth--
				add(ruleIsOp, position735)
			}
			return true
		l734:
			position, tokenIndex, depth = position734, tokenIndex734, depth734
			return false
		},
		/* 61 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position738, tokenIndex738, depth738 := position, tokenIndex, depth
			{
				position739 := position
				depth++
				{
					position740, tokenIndex740, depth740 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l741
					}
					goto l740
				l741:
					position, tokenIndex, depth = position740, tokenIndex740, depth740
					if !_rules[ruleMinus]() {
						goto l738
					}
				}
			l740:
				depth--
				add(rulePlusMinusOp, position739)
			}
			return true
		l738:
			position, tokenIndex, depth = position738, tokenIndex738, depth738
			return false
		},
		/* 62 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position742, tokenIndex742, depth742 := position, tokenIndex, depth
			{
				position743 := position
				depth++
				{
					position744, tokenIndex744, depth744 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l745
					}
					goto l744
				l745:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
					if !_rules[ruleDivide]() {
						goto l746
					}
					goto l744
				l746:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
					if !_rules[ruleModulo]() {
						goto l742
					}
				}
			l744:
				depth--
				add(ruleMultDivOp, position743)
			}
			return true
		l742:
			position, tokenIndex, depth = position742, tokenIndex742, depth742
			return false
		},
		/* 63 Stream <- <(<ident> Action47)> */
		func() bool {
			position747, tokenIndex747, depth747 := position, tokenIndex, depth
			{
				position748 := position
				depth++
				{
					position749 := position
					depth++
					if !_rules[ruleident]() {
						goto l747
					}
					depth--
					add(rulePegText, position749)
				}
				if !_rules[ruleAction47]() {
					goto l747
				}
				depth--
				add(ruleStream, position748)
			}
			return true
		l747:
			position, tokenIndex, depth = position747, tokenIndex747, depth747
			return false
		},
		/* 64 RowMeta <- <RowTimestamp> */
		func() bool {
			position750, tokenIndex750, depth750 := position, tokenIndex, depth
			{
				position751 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l750
				}
				depth--
				add(ruleRowMeta, position751)
			}
			return true
		l750:
			position, tokenIndex, depth = position750, tokenIndex750, depth750
			return false
		},
		/* 65 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action48)> */
		func() bool {
			position752, tokenIndex752, depth752 := position, tokenIndex, depth
			{
				position753 := position
				depth++
				{
					position754 := position
					depth++
					{
						position755, tokenIndex755, depth755 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l755
						}
						if buffer[position] != rune(':') {
							goto l755
						}
						position++
						goto l756
					l755:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
					}
				l756:
					if buffer[position] != rune('t') {
						goto l752
					}
					position++
					if buffer[position] != rune('s') {
						goto l752
					}
					position++
					if buffer[position] != rune('(') {
						goto l752
					}
					position++
					if buffer[position] != rune(')') {
						goto l752
					}
					position++
					depth--
					add(rulePegText, position754)
				}
				if !_rules[ruleAction48]() {
					goto l752
				}
				depth--
				add(ruleRowTimestamp, position753)
			}
			return true
		l752:
			position, tokenIndex, depth = position752, tokenIndex752, depth752
			return false
		},
		/* 66 RowValue <- <(<((ident ':')? jsonPath)> Action49)> */
		func() bool {
			position757, tokenIndex757, depth757 := position, tokenIndex, depth
			{
				position758 := position
				depth++
				{
					position759 := position
					depth++
					{
						position760, tokenIndex760, depth760 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l760
						}
						if buffer[position] != rune(':') {
							goto l760
						}
						position++
						goto l761
					l760:
						position, tokenIndex, depth = position760, tokenIndex760, depth760
					}
				l761:
					if !_rules[rulejsonPath]() {
						goto l757
					}
					depth--
					add(rulePegText, position759)
				}
				if !_rules[ruleAction49]() {
					goto l757
				}
				depth--
				add(ruleRowValue, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 67 NumericLiteral <- <(<('-'? [0-9]+)> Action50)> */
		func() bool {
			position762, tokenIndex762, depth762 := position, tokenIndex, depth
			{
				position763 := position
				depth++
				{
					position764 := position
					depth++
					{
						position765, tokenIndex765, depth765 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l765
						}
						position++
						goto l766
					l765:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
					}
				l766:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l762
					}
					position++
				l767:
					{
						position768, tokenIndex768, depth768 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position768, tokenIndex768, depth768
					}
					depth--
					add(rulePegText, position764)
				}
				if !_rules[ruleAction50]() {
					goto l762
				}
				depth--
				add(ruleNumericLiteral, position763)
			}
			return true
		l762:
			position, tokenIndex, depth = position762, tokenIndex762, depth762
			return false
		},
		/* 68 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action51)> */
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
						if buffer[position] != rune('-') {
							goto l772
						}
						position++
						goto l773
					l772:
						position, tokenIndex, depth = position772, tokenIndex772, depth772
					}
				l773:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l769
					}
					position++
				l774:
					{
						position775, tokenIndex775, depth775 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l775
						}
						position++
						goto l774
					l775:
						position, tokenIndex, depth = position775, tokenIndex775, depth775
					}
					if buffer[position] != rune('.') {
						goto l769
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l769
					}
					position++
				l776:
					{
						position777, tokenIndex777, depth777 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l777
						}
						position++
						goto l776
					l777:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
					}
					depth--
					add(rulePegText, position771)
				}
				if !_rules[ruleAction51]() {
					goto l769
				}
				depth--
				add(ruleFloatLiteral, position770)
			}
			return true
		l769:
			position, tokenIndex, depth = position769, tokenIndex769, depth769
			return false
		},
		/* 69 Function <- <(<ident> Action52)> */
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
				if !_rules[ruleAction52]() {
					goto l778
				}
				depth--
				add(ruleFunction, position779)
			}
			return true
		l778:
			position, tokenIndex, depth = position778, tokenIndex778, depth778
			return false
		},
		/* 70 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action53)> */
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
						if buffer[position] != rune('n') {
							goto l785
						}
						position++
						goto l784
					l785:
						position, tokenIndex, depth = position784, tokenIndex784, depth784
						if buffer[position] != rune('N') {
							goto l781
						}
						position++
					}
				l784:
					{
						position786, tokenIndex786, depth786 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l787
						}
						position++
						goto l786
					l787:
						position, tokenIndex, depth = position786, tokenIndex786, depth786
						if buffer[position] != rune('U') {
							goto l781
						}
						position++
					}
				l786:
					{
						position788, tokenIndex788, depth788 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l789
						}
						position++
						goto l788
					l789:
						position, tokenIndex, depth = position788, tokenIndex788, depth788
						if buffer[position] != rune('L') {
							goto l781
						}
						position++
					}
				l788:
					{
						position790, tokenIndex790, depth790 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l791
						}
						position++
						goto l790
					l791:
						position, tokenIndex, depth = position790, tokenIndex790, depth790
						if buffer[position] != rune('L') {
							goto l781
						}
						position++
					}
				l790:
					depth--
					add(rulePegText, position783)
				}
				if !_rules[ruleAction53]() {
					goto l781
				}
				depth--
				add(ruleNullLiteral, position782)
			}
			return true
		l781:
			position, tokenIndex, depth = position781, tokenIndex781, depth781
			return false
		},
		/* 71 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position792, tokenIndex792, depth792 := position, tokenIndex, depth
			{
				position793 := position
				depth++
				{
					position794, tokenIndex794, depth794 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l795
					}
					goto l794
				l795:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if !_rules[ruleFALSE]() {
						goto l792
					}
				}
			l794:
				depth--
				add(ruleBooleanLiteral, position793)
			}
			return true
		l792:
			position, tokenIndex, depth = position792, tokenIndex792, depth792
			return false
		},
		/* 72 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action54)> */
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
						if buffer[position] != rune('t') {
							goto l800
						}
						position++
						goto l799
					l800:
						position, tokenIndex, depth = position799, tokenIndex799, depth799
						if buffer[position] != rune('T') {
							goto l796
						}
						position++
					}
				l799:
					{
						position801, tokenIndex801, depth801 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l802
						}
						position++
						goto l801
					l802:
						position, tokenIndex, depth = position801, tokenIndex801, depth801
						if buffer[position] != rune('R') {
							goto l796
						}
						position++
					}
				l801:
					{
						position803, tokenIndex803, depth803 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l804
						}
						position++
						goto l803
					l804:
						position, tokenIndex, depth = position803, tokenIndex803, depth803
						if buffer[position] != rune('U') {
							goto l796
						}
						position++
					}
				l803:
					{
						position805, tokenIndex805, depth805 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l806
						}
						position++
						goto l805
					l806:
						position, tokenIndex, depth = position805, tokenIndex805, depth805
						if buffer[position] != rune('E') {
							goto l796
						}
						position++
					}
				l805:
					depth--
					add(rulePegText, position798)
				}
				if !_rules[ruleAction54]() {
					goto l796
				}
				depth--
				add(ruleTRUE, position797)
			}
			return true
		l796:
			position, tokenIndex, depth = position796, tokenIndex796, depth796
			return false
		},
		/* 73 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action55)> */
		func() bool {
			position807, tokenIndex807, depth807 := position, tokenIndex, depth
			{
				position808 := position
				depth++
				{
					position809 := position
					depth++
					{
						position810, tokenIndex810, depth810 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l811
						}
						position++
						goto l810
					l811:
						position, tokenIndex, depth = position810, tokenIndex810, depth810
						if buffer[position] != rune('F') {
							goto l807
						}
						position++
					}
				l810:
					{
						position812, tokenIndex812, depth812 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l813
						}
						position++
						goto l812
					l813:
						position, tokenIndex, depth = position812, tokenIndex812, depth812
						if buffer[position] != rune('A') {
							goto l807
						}
						position++
					}
				l812:
					{
						position814, tokenIndex814, depth814 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l815
						}
						position++
						goto l814
					l815:
						position, tokenIndex, depth = position814, tokenIndex814, depth814
						if buffer[position] != rune('L') {
							goto l807
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
							goto l807
						}
						position++
					}
				l816:
					{
						position818, tokenIndex818, depth818 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l819
						}
						position++
						goto l818
					l819:
						position, tokenIndex, depth = position818, tokenIndex818, depth818
						if buffer[position] != rune('E') {
							goto l807
						}
						position++
					}
				l818:
					depth--
					add(rulePegText, position809)
				}
				if !_rules[ruleAction55]() {
					goto l807
				}
				depth--
				add(ruleFALSE, position808)
			}
			return true
		l807:
			position, tokenIndex, depth = position807, tokenIndex807, depth807
			return false
		},
		/* 74 Wildcard <- <(<'*'> Action56)> */
		func() bool {
			position820, tokenIndex820, depth820 := position, tokenIndex, depth
			{
				position821 := position
				depth++
				{
					position822 := position
					depth++
					if buffer[position] != rune('*') {
						goto l820
					}
					position++
					depth--
					add(rulePegText, position822)
				}
				if !_rules[ruleAction56]() {
					goto l820
				}
				depth--
				add(ruleWildcard, position821)
			}
			return true
		l820:
			position, tokenIndex, depth = position820, tokenIndex820, depth820
			return false
		},
		/* 75 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action57)> */
		func() bool {
			position823, tokenIndex823, depth823 := position, tokenIndex, depth
			{
				position824 := position
				depth++
				{
					position825 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l823
					}
					position++
				l826:
					{
						position827, tokenIndex827, depth827 := position, tokenIndex, depth
						{
							position828, tokenIndex828, depth828 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l829
							}
							position++
							if buffer[position] != rune('\'') {
								goto l829
							}
							position++
							goto l828
						l829:
							position, tokenIndex, depth = position828, tokenIndex828, depth828
							{
								position830, tokenIndex830, depth830 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l830
								}
								position++
								goto l827
							l830:
								position, tokenIndex, depth = position830, tokenIndex830, depth830
							}
							if !matchDot() {
								goto l827
							}
						}
					l828:
						goto l826
					l827:
						position, tokenIndex, depth = position827, tokenIndex827, depth827
					}
					if buffer[position] != rune('\'') {
						goto l823
					}
					position++
					depth--
					add(rulePegText, position825)
				}
				if !_rules[ruleAction57]() {
					goto l823
				}
				depth--
				add(ruleStringLiteral, position824)
			}
			return true
		l823:
			position, tokenIndex, depth = position823, tokenIndex823, depth823
			return false
		},
		/* 76 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action58)> */
		func() bool {
			position831, tokenIndex831, depth831 := position, tokenIndex, depth
			{
				position832 := position
				depth++
				{
					position833 := position
					depth++
					{
						position834, tokenIndex834, depth834 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l835
						}
						position++
						goto l834
					l835:
						position, tokenIndex, depth = position834, tokenIndex834, depth834
						if buffer[position] != rune('I') {
							goto l831
						}
						position++
					}
				l834:
					{
						position836, tokenIndex836, depth836 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l837
						}
						position++
						goto l836
					l837:
						position, tokenIndex, depth = position836, tokenIndex836, depth836
						if buffer[position] != rune('S') {
							goto l831
						}
						position++
					}
				l836:
					{
						position838, tokenIndex838, depth838 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l839
						}
						position++
						goto l838
					l839:
						position, tokenIndex, depth = position838, tokenIndex838, depth838
						if buffer[position] != rune('T') {
							goto l831
						}
						position++
					}
				l838:
					{
						position840, tokenIndex840, depth840 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l841
						}
						position++
						goto l840
					l841:
						position, tokenIndex, depth = position840, tokenIndex840, depth840
						if buffer[position] != rune('R') {
							goto l831
						}
						position++
					}
				l840:
					{
						position842, tokenIndex842, depth842 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l843
						}
						position++
						goto l842
					l843:
						position, tokenIndex, depth = position842, tokenIndex842, depth842
						if buffer[position] != rune('E') {
							goto l831
						}
						position++
					}
				l842:
					{
						position844, tokenIndex844, depth844 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l845
						}
						position++
						goto l844
					l845:
						position, tokenIndex, depth = position844, tokenIndex844, depth844
						if buffer[position] != rune('A') {
							goto l831
						}
						position++
					}
				l844:
					{
						position846, tokenIndex846, depth846 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l847
						}
						position++
						goto l846
					l847:
						position, tokenIndex, depth = position846, tokenIndex846, depth846
						if buffer[position] != rune('M') {
							goto l831
						}
						position++
					}
				l846:
					depth--
					add(rulePegText, position833)
				}
				if !_rules[ruleAction58]() {
					goto l831
				}
				depth--
				add(ruleISTREAM, position832)
			}
			return true
		l831:
			position, tokenIndex, depth = position831, tokenIndex831, depth831
			return false
		},
		/* 77 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action59)> */
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
						if buffer[position] != rune('d') {
							goto l852
						}
						position++
						goto l851
					l852:
						position, tokenIndex, depth = position851, tokenIndex851, depth851
						if buffer[position] != rune('D') {
							goto l848
						}
						position++
					}
				l851:
					{
						position853, tokenIndex853, depth853 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l854
						}
						position++
						goto l853
					l854:
						position, tokenIndex, depth = position853, tokenIndex853, depth853
						if buffer[position] != rune('S') {
							goto l848
						}
						position++
					}
				l853:
					{
						position855, tokenIndex855, depth855 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l856
						}
						position++
						goto l855
					l856:
						position, tokenIndex, depth = position855, tokenIndex855, depth855
						if buffer[position] != rune('T') {
							goto l848
						}
						position++
					}
				l855:
					{
						position857, tokenIndex857, depth857 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l858
						}
						position++
						goto l857
					l858:
						position, tokenIndex, depth = position857, tokenIndex857, depth857
						if buffer[position] != rune('R') {
							goto l848
						}
						position++
					}
				l857:
					{
						position859, tokenIndex859, depth859 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l860
						}
						position++
						goto l859
					l860:
						position, tokenIndex, depth = position859, tokenIndex859, depth859
						if buffer[position] != rune('E') {
							goto l848
						}
						position++
					}
				l859:
					{
						position861, tokenIndex861, depth861 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l862
						}
						position++
						goto l861
					l862:
						position, tokenIndex, depth = position861, tokenIndex861, depth861
						if buffer[position] != rune('A') {
							goto l848
						}
						position++
					}
				l861:
					{
						position863, tokenIndex863, depth863 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l864
						}
						position++
						goto l863
					l864:
						position, tokenIndex, depth = position863, tokenIndex863, depth863
						if buffer[position] != rune('M') {
							goto l848
						}
						position++
					}
				l863:
					depth--
					add(rulePegText, position850)
				}
				if !_rules[ruleAction59]() {
					goto l848
				}
				depth--
				add(ruleDSTREAM, position849)
			}
			return true
		l848:
			position, tokenIndex, depth = position848, tokenIndex848, depth848
			return false
		},
		/* 78 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action60)> */
		func() bool {
			position865, tokenIndex865, depth865 := position, tokenIndex, depth
			{
				position866 := position
				depth++
				{
					position867 := position
					depth++
					{
						position868, tokenIndex868, depth868 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l869
						}
						position++
						goto l868
					l869:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('R') {
							goto l865
						}
						position++
					}
				l868:
					{
						position870, tokenIndex870, depth870 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l871
						}
						position++
						goto l870
					l871:
						position, tokenIndex, depth = position870, tokenIndex870, depth870
						if buffer[position] != rune('S') {
							goto l865
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
							goto l865
						}
						position++
					}
				l872:
					{
						position874, tokenIndex874, depth874 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l875
						}
						position++
						goto l874
					l875:
						position, tokenIndex, depth = position874, tokenIndex874, depth874
						if buffer[position] != rune('R') {
							goto l865
						}
						position++
					}
				l874:
					{
						position876, tokenIndex876, depth876 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l877
						}
						position++
						goto l876
					l877:
						position, tokenIndex, depth = position876, tokenIndex876, depth876
						if buffer[position] != rune('E') {
							goto l865
						}
						position++
					}
				l876:
					{
						position878, tokenIndex878, depth878 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l879
						}
						position++
						goto l878
					l879:
						position, tokenIndex, depth = position878, tokenIndex878, depth878
						if buffer[position] != rune('A') {
							goto l865
						}
						position++
					}
				l878:
					{
						position880, tokenIndex880, depth880 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l881
						}
						position++
						goto l880
					l881:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
						if buffer[position] != rune('M') {
							goto l865
						}
						position++
					}
				l880:
					depth--
					add(rulePegText, position867)
				}
				if !_rules[ruleAction60]() {
					goto l865
				}
				depth--
				add(ruleRSTREAM, position866)
			}
			return true
		l865:
			position, tokenIndex, depth = position865, tokenIndex865, depth865
			return false
		},
		/* 79 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action61)> */
		func() bool {
			position882, tokenIndex882, depth882 := position, tokenIndex, depth
			{
				position883 := position
				depth++
				{
					position884 := position
					depth++
					{
						position885, tokenIndex885, depth885 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l886
						}
						position++
						goto l885
					l886:
						position, tokenIndex, depth = position885, tokenIndex885, depth885
						if buffer[position] != rune('T') {
							goto l882
						}
						position++
					}
				l885:
					{
						position887, tokenIndex887, depth887 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l888
						}
						position++
						goto l887
					l888:
						position, tokenIndex, depth = position887, tokenIndex887, depth887
						if buffer[position] != rune('U') {
							goto l882
						}
						position++
					}
				l887:
					{
						position889, tokenIndex889, depth889 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l890
						}
						position++
						goto l889
					l890:
						position, tokenIndex, depth = position889, tokenIndex889, depth889
						if buffer[position] != rune('P') {
							goto l882
						}
						position++
					}
				l889:
					{
						position891, tokenIndex891, depth891 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l892
						}
						position++
						goto l891
					l892:
						position, tokenIndex, depth = position891, tokenIndex891, depth891
						if buffer[position] != rune('L') {
							goto l882
						}
						position++
					}
				l891:
					{
						position893, tokenIndex893, depth893 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l894
						}
						position++
						goto l893
					l894:
						position, tokenIndex, depth = position893, tokenIndex893, depth893
						if buffer[position] != rune('E') {
							goto l882
						}
						position++
					}
				l893:
					{
						position895, tokenIndex895, depth895 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l896
						}
						position++
						goto l895
					l896:
						position, tokenIndex, depth = position895, tokenIndex895, depth895
						if buffer[position] != rune('S') {
							goto l882
						}
						position++
					}
				l895:
					depth--
					add(rulePegText, position884)
				}
				if !_rules[ruleAction61]() {
					goto l882
				}
				depth--
				add(ruleTUPLES, position883)
			}
			return true
		l882:
			position, tokenIndex, depth = position882, tokenIndex882, depth882
			return false
		},
		/* 80 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action62)> */
		func() bool {
			position897, tokenIndex897, depth897 := position, tokenIndex, depth
			{
				position898 := position
				depth++
				{
					position899 := position
					depth++
					{
						position900, tokenIndex900, depth900 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l901
						}
						position++
						goto l900
					l901:
						position, tokenIndex, depth = position900, tokenIndex900, depth900
						if buffer[position] != rune('S') {
							goto l897
						}
						position++
					}
				l900:
					{
						position902, tokenIndex902, depth902 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l903
						}
						position++
						goto l902
					l903:
						position, tokenIndex, depth = position902, tokenIndex902, depth902
						if buffer[position] != rune('E') {
							goto l897
						}
						position++
					}
				l902:
					{
						position904, tokenIndex904, depth904 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l905
						}
						position++
						goto l904
					l905:
						position, tokenIndex, depth = position904, tokenIndex904, depth904
						if buffer[position] != rune('C') {
							goto l897
						}
						position++
					}
				l904:
					{
						position906, tokenIndex906, depth906 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l907
						}
						position++
						goto l906
					l907:
						position, tokenIndex, depth = position906, tokenIndex906, depth906
						if buffer[position] != rune('O') {
							goto l897
						}
						position++
					}
				l906:
					{
						position908, tokenIndex908, depth908 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l909
						}
						position++
						goto l908
					l909:
						position, tokenIndex, depth = position908, tokenIndex908, depth908
						if buffer[position] != rune('N') {
							goto l897
						}
						position++
					}
				l908:
					{
						position910, tokenIndex910, depth910 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l911
						}
						position++
						goto l910
					l911:
						position, tokenIndex, depth = position910, tokenIndex910, depth910
						if buffer[position] != rune('D') {
							goto l897
						}
						position++
					}
				l910:
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
							goto l897
						}
						position++
					}
				l912:
					depth--
					add(rulePegText, position899)
				}
				if !_rules[ruleAction62]() {
					goto l897
				}
				depth--
				add(ruleSECONDS, position898)
			}
			return true
		l897:
			position, tokenIndex, depth = position897, tokenIndex897, depth897
			return false
		},
		/* 81 StreamIdentifier <- <(<ident> Action63)> */
		func() bool {
			position914, tokenIndex914, depth914 := position, tokenIndex, depth
			{
				position915 := position
				depth++
				{
					position916 := position
					depth++
					if !_rules[ruleident]() {
						goto l914
					}
					depth--
					add(rulePegText, position916)
				}
				if !_rules[ruleAction63]() {
					goto l914
				}
				depth--
				add(ruleStreamIdentifier, position915)
			}
			return true
		l914:
			position, tokenIndex, depth = position914, tokenIndex914, depth914
			return false
		},
		/* 82 SourceSinkType <- <(<ident> Action64)> */
		func() bool {
			position917, tokenIndex917, depth917 := position, tokenIndex, depth
			{
				position918 := position
				depth++
				{
					position919 := position
					depth++
					if !_rules[ruleident]() {
						goto l917
					}
					depth--
					add(rulePegText, position919)
				}
				if !_rules[ruleAction64]() {
					goto l917
				}
				depth--
				add(ruleSourceSinkType, position918)
			}
			return true
		l917:
			position, tokenIndex, depth = position917, tokenIndex917, depth917
			return false
		},
		/* 83 SourceSinkParamKey <- <(<ident> Action65)> */
		func() bool {
			position920, tokenIndex920, depth920 := position, tokenIndex, depth
			{
				position921 := position
				depth++
				{
					position922 := position
					depth++
					if !_rules[ruleident]() {
						goto l920
					}
					depth--
					add(rulePegText, position922)
				}
				if !_rules[ruleAction65]() {
					goto l920
				}
				depth--
				add(ruleSourceSinkParamKey, position921)
			}
			return true
		l920:
			position, tokenIndex, depth = position920, tokenIndex920, depth920
			return false
		},
		/* 84 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action66)> */
		func() bool {
			position923, tokenIndex923, depth923 := position, tokenIndex, depth
			{
				position924 := position
				depth++
				{
					position925 := position
					depth++
					{
						position926, tokenIndex926, depth926 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('P') {
							goto l923
						}
						position++
					}
				l926:
					{
						position928, tokenIndex928, depth928 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l929
						}
						position++
						goto l928
					l929:
						position, tokenIndex, depth = position928, tokenIndex928, depth928
						if buffer[position] != rune('A') {
							goto l923
						}
						position++
					}
				l928:
					{
						position930, tokenIndex930, depth930 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l931
						}
						position++
						goto l930
					l931:
						position, tokenIndex, depth = position930, tokenIndex930, depth930
						if buffer[position] != rune('U') {
							goto l923
						}
						position++
					}
				l930:
					{
						position932, tokenIndex932, depth932 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l933
						}
						position++
						goto l932
					l933:
						position, tokenIndex, depth = position932, tokenIndex932, depth932
						if buffer[position] != rune('S') {
							goto l923
						}
						position++
					}
				l932:
					{
						position934, tokenIndex934, depth934 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l935
						}
						position++
						goto l934
					l935:
						position, tokenIndex, depth = position934, tokenIndex934, depth934
						if buffer[position] != rune('E') {
							goto l923
						}
						position++
					}
				l934:
					{
						position936, tokenIndex936, depth936 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l937
						}
						position++
						goto l936
					l937:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if buffer[position] != rune('D') {
							goto l923
						}
						position++
					}
				l936:
					depth--
					add(rulePegText, position925)
				}
				if !_rules[ruleAction66]() {
					goto l923
				}
				depth--
				add(rulePaused, position924)
			}
			return true
		l923:
			position, tokenIndex, depth = position923, tokenIndex923, depth923
			return false
		},
		/* 85 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action67)> */
		func() bool {
			position938, tokenIndex938, depth938 := position, tokenIndex, depth
			{
				position939 := position
				depth++
				{
					position940 := position
					depth++
					{
						position941, tokenIndex941, depth941 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l942
						}
						position++
						goto l941
					l942:
						position, tokenIndex, depth = position941, tokenIndex941, depth941
						if buffer[position] != rune('U') {
							goto l938
						}
						position++
					}
				l941:
					{
						position943, tokenIndex943, depth943 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l944
						}
						position++
						goto l943
					l944:
						position, tokenIndex, depth = position943, tokenIndex943, depth943
						if buffer[position] != rune('N') {
							goto l938
						}
						position++
					}
				l943:
					{
						position945, tokenIndex945, depth945 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l946
						}
						position++
						goto l945
					l946:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
						if buffer[position] != rune('P') {
							goto l938
						}
						position++
					}
				l945:
					{
						position947, tokenIndex947, depth947 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l948
						}
						position++
						goto l947
					l948:
						position, tokenIndex, depth = position947, tokenIndex947, depth947
						if buffer[position] != rune('A') {
							goto l938
						}
						position++
					}
				l947:
					{
						position949, tokenIndex949, depth949 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l950
						}
						position++
						goto l949
					l950:
						position, tokenIndex, depth = position949, tokenIndex949, depth949
						if buffer[position] != rune('U') {
							goto l938
						}
						position++
					}
				l949:
					{
						position951, tokenIndex951, depth951 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l952
						}
						position++
						goto l951
					l952:
						position, tokenIndex, depth = position951, tokenIndex951, depth951
						if buffer[position] != rune('S') {
							goto l938
						}
						position++
					}
				l951:
					{
						position953, tokenIndex953, depth953 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l954
						}
						position++
						goto l953
					l954:
						position, tokenIndex, depth = position953, tokenIndex953, depth953
						if buffer[position] != rune('E') {
							goto l938
						}
						position++
					}
				l953:
					{
						position955, tokenIndex955, depth955 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l956
						}
						position++
						goto l955
					l956:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('D') {
							goto l938
						}
						position++
					}
				l955:
					depth--
					add(rulePegText, position940)
				}
				if !_rules[ruleAction67]() {
					goto l938
				}
				depth--
				add(ruleUnpaused, position939)
			}
			return true
		l938:
			position, tokenIndex, depth = position938, tokenIndex938, depth938
			return false
		},
		/* 86 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action68)> */
		func() bool {
			position957, tokenIndex957, depth957 := position, tokenIndex, depth
			{
				position958 := position
				depth++
				{
					position959 := position
					depth++
					{
						position960, tokenIndex960, depth960 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l961
						}
						position++
						goto l960
					l961:
						position, tokenIndex, depth = position960, tokenIndex960, depth960
						if buffer[position] != rune('O') {
							goto l957
						}
						position++
					}
				l960:
					{
						position962, tokenIndex962, depth962 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l963
						}
						position++
						goto l962
					l963:
						position, tokenIndex, depth = position962, tokenIndex962, depth962
						if buffer[position] != rune('R') {
							goto l957
						}
						position++
					}
				l962:
					depth--
					add(rulePegText, position959)
				}
				if !_rules[ruleAction68]() {
					goto l957
				}
				depth--
				add(ruleOr, position958)
			}
			return true
		l957:
			position, tokenIndex, depth = position957, tokenIndex957, depth957
			return false
		},
		/* 87 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action69)> */
		func() bool {
			position964, tokenIndex964, depth964 := position, tokenIndex, depth
			{
				position965 := position
				depth++
				{
					position966 := position
					depth++
					{
						position967, tokenIndex967, depth967 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l968
						}
						position++
						goto l967
					l968:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if buffer[position] != rune('A') {
							goto l964
						}
						position++
					}
				l967:
					{
						position969, tokenIndex969, depth969 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l970
						}
						position++
						goto l969
					l970:
						position, tokenIndex, depth = position969, tokenIndex969, depth969
						if buffer[position] != rune('N') {
							goto l964
						}
						position++
					}
				l969:
					{
						position971, tokenIndex971, depth971 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l972
						}
						position++
						goto l971
					l972:
						position, tokenIndex, depth = position971, tokenIndex971, depth971
						if buffer[position] != rune('D') {
							goto l964
						}
						position++
					}
				l971:
					depth--
					add(rulePegText, position966)
				}
				if !_rules[ruleAction69]() {
					goto l964
				}
				depth--
				add(ruleAnd, position965)
			}
			return true
		l964:
			position, tokenIndex, depth = position964, tokenIndex964, depth964
			return false
		},
		/* 88 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action70)> */
		func() bool {
			position973, tokenIndex973, depth973 := position, tokenIndex, depth
			{
				position974 := position
				depth++
				{
					position975 := position
					depth++
					{
						position976, tokenIndex976, depth976 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l977
						}
						position++
						goto l976
					l977:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('N') {
							goto l973
						}
						position++
					}
				l976:
					{
						position978, tokenIndex978, depth978 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l979
						}
						position++
						goto l978
					l979:
						position, tokenIndex, depth = position978, tokenIndex978, depth978
						if buffer[position] != rune('O') {
							goto l973
						}
						position++
					}
				l978:
					{
						position980, tokenIndex980, depth980 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l981
						}
						position++
						goto l980
					l981:
						position, tokenIndex, depth = position980, tokenIndex980, depth980
						if buffer[position] != rune('T') {
							goto l973
						}
						position++
					}
				l980:
					depth--
					add(rulePegText, position975)
				}
				if !_rules[ruleAction70]() {
					goto l973
				}
				depth--
				add(ruleNot, position974)
			}
			return true
		l973:
			position, tokenIndex, depth = position973, tokenIndex973, depth973
			return false
		},
		/* 89 Equal <- <(<'='> Action71)> */
		func() bool {
			position982, tokenIndex982, depth982 := position, tokenIndex, depth
			{
				position983 := position
				depth++
				{
					position984 := position
					depth++
					if buffer[position] != rune('=') {
						goto l982
					}
					position++
					depth--
					add(rulePegText, position984)
				}
				if !_rules[ruleAction71]() {
					goto l982
				}
				depth--
				add(ruleEqual, position983)
			}
			return true
		l982:
			position, tokenIndex, depth = position982, tokenIndex982, depth982
			return false
		},
		/* 90 Less <- <(<'<'> Action72)> */
		func() bool {
			position985, tokenIndex985, depth985 := position, tokenIndex, depth
			{
				position986 := position
				depth++
				{
					position987 := position
					depth++
					if buffer[position] != rune('<') {
						goto l985
					}
					position++
					depth--
					add(rulePegText, position987)
				}
				if !_rules[ruleAction72]() {
					goto l985
				}
				depth--
				add(ruleLess, position986)
			}
			return true
		l985:
			position, tokenIndex, depth = position985, tokenIndex985, depth985
			return false
		},
		/* 91 LessOrEqual <- <(<('<' '=')> Action73)> */
		func() bool {
			position988, tokenIndex988, depth988 := position, tokenIndex, depth
			{
				position989 := position
				depth++
				{
					position990 := position
					depth++
					if buffer[position] != rune('<') {
						goto l988
					}
					position++
					if buffer[position] != rune('=') {
						goto l988
					}
					position++
					depth--
					add(rulePegText, position990)
				}
				if !_rules[ruleAction73]() {
					goto l988
				}
				depth--
				add(ruleLessOrEqual, position989)
			}
			return true
		l988:
			position, tokenIndex, depth = position988, tokenIndex988, depth988
			return false
		},
		/* 92 Greater <- <(<'>'> Action74)> */
		func() bool {
			position991, tokenIndex991, depth991 := position, tokenIndex, depth
			{
				position992 := position
				depth++
				{
					position993 := position
					depth++
					if buffer[position] != rune('>') {
						goto l991
					}
					position++
					depth--
					add(rulePegText, position993)
				}
				if !_rules[ruleAction74]() {
					goto l991
				}
				depth--
				add(ruleGreater, position992)
			}
			return true
		l991:
			position, tokenIndex, depth = position991, tokenIndex991, depth991
			return false
		},
		/* 93 GreaterOrEqual <- <(<('>' '=')> Action75)> */
		func() bool {
			position994, tokenIndex994, depth994 := position, tokenIndex, depth
			{
				position995 := position
				depth++
				{
					position996 := position
					depth++
					if buffer[position] != rune('>') {
						goto l994
					}
					position++
					if buffer[position] != rune('=') {
						goto l994
					}
					position++
					depth--
					add(rulePegText, position996)
				}
				if !_rules[ruleAction75]() {
					goto l994
				}
				depth--
				add(ruleGreaterOrEqual, position995)
			}
			return true
		l994:
			position, tokenIndex, depth = position994, tokenIndex994, depth994
			return false
		},
		/* 94 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action76)> */
		func() bool {
			position997, tokenIndex997, depth997 := position, tokenIndex, depth
			{
				position998 := position
				depth++
				{
					position999 := position
					depth++
					{
						position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1001
						}
						position++
						if buffer[position] != rune('=') {
							goto l1001
						}
						position++
						goto l1000
					l1001:
						position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
						if buffer[position] != rune('<') {
							goto l997
						}
						position++
						if buffer[position] != rune('>') {
							goto l997
						}
						position++
					}
				l1000:
					depth--
					add(rulePegText, position999)
				}
				if !_rules[ruleAction76]() {
					goto l997
				}
				depth--
				add(ruleNotEqual, position998)
			}
			return true
		l997:
			position, tokenIndex, depth = position997, tokenIndex997, depth997
			return false
		},
		/* 95 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action77)> */
		func() bool {
			position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
			{
				position1003 := position
				depth++
				{
					position1004 := position
					depth++
					{
						position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1006
						}
						position++
						goto l1005
					l1006:
						position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
						if buffer[position] != rune('I') {
							goto l1002
						}
						position++
					}
				l1005:
					{
						position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1008
						}
						position++
						goto l1007
					l1008:
						position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
						if buffer[position] != rune('S') {
							goto l1002
						}
						position++
					}
				l1007:
					depth--
					add(rulePegText, position1004)
				}
				if !_rules[ruleAction77]() {
					goto l1002
				}
				depth--
				add(ruleIs, position1003)
			}
			return true
		l1002:
			position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
			return false
		},
		/* 96 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action78)> */
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
						if buffer[position] != rune('i') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('I') {
							goto l1009
						}
						position++
					}
				l1012:
					{
						position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1015
						}
						position++
						goto l1014
					l1015:
						position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
						if buffer[position] != rune('S') {
							goto l1009
						}
						position++
					}
				l1014:
					if !_rules[rulesp]() {
						goto l1009
					}
					{
						position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1017
						}
						position++
						goto l1016
					l1017:
						position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
						if buffer[position] != rune('N') {
							goto l1009
						}
						position++
					}
				l1016:
					{
						position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1019
						}
						position++
						goto l1018
					l1019:
						position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
						if buffer[position] != rune('O') {
							goto l1009
						}
						position++
					}
				l1018:
					{
						position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1021
						}
						position++
						goto l1020
					l1021:
						position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
						if buffer[position] != rune('T') {
							goto l1009
						}
						position++
					}
				l1020:
					depth--
					add(rulePegText, position1011)
				}
				if !_rules[ruleAction78]() {
					goto l1009
				}
				depth--
				add(ruleIsNot, position1010)
			}
			return true
		l1009:
			position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
			return false
		},
		/* 97 Plus <- <(<'+'> Action79)> */
		func() bool {
			position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
			{
				position1023 := position
				depth++
				{
					position1024 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1022
					}
					position++
					depth--
					add(rulePegText, position1024)
				}
				if !_rules[ruleAction79]() {
					goto l1022
				}
				depth--
				add(rulePlus, position1023)
			}
			return true
		l1022:
			position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
			return false
		},
		/* 98 Minus <- <(<'-'> Action80)> */
		func() bool {
			position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
			{
				position1026 := position
				depth++
				{
					position1027 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1025
					}
					position++
					depth--
					add(rulePegText, position1027)
				}
				if !_rules[ruleAction80]() {
					goto l1025
				}
				depth--
				add(ruleMinus, position1026)
			}
			return true
		l1025:
			position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
			return false
		},
		/* 99 Multiply <- <(<'*'> Action81)> */
		func() bool {
			position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
			{
				position1029 := position
				depth++
				{
					position1030 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1028
					}
					position++
					depth--
					add(rulePegText, position1030)
				}
				if !_rules[ruleAction81]() {
					goto l1028
				}
				depth--
				add(ruleMultiply, position1029)
			}
			return true
		l1028:
			position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
			return false
		},
		/* 100 Divide <- <(<'/'> Action82)> */
		func() bool {
			position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
			{
				position1032 := position
				depth++
				{
					position1033 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1031
					}
					position++
					depth--
					add(rulePegText, position1033)
				}
				if !_rules[ruleAction82]() {
					goto l1031
				}
				depth--
				add(ruleDivide, position1032)
			}
			return true
		l1031:
			position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
			return false
		},
		/* 101 Modulo <- <(<'%'> Action83)> */
		func() bool {
			position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
			{
				position1035 := position
				depth++
				{
					position1036 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1034
					}
					position++
					depth--
					add(rulePegText, position1036)
				}
				if !_rules[ruleAction83]() {
					goto l1034
				}
				depth--
				add(ruleModulo, position1035)
			}
			return true
		l1034:
			position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
			return false
		},
		/* 102 UnaryMinus <- <(<'-'> Action84)> */
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
				if !_rules[ruleAction84]() {
					goto l1037
				}
				depth--
				add(ruleUnaryMinus, position1038)
			}
			return true
		l1037:
			position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
			return false
		},
		/* 103 Identifier <- <(<ident> Action85)> */
		func() bool {
			position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
			{
				position1041 := position
				depth++
				{
					position1042 := position
					depth++
					if !_rules[ruleident]() {
						goto l1040
					}
					depth--
					add(rulePegText, position1042)
				}
				if !_rules[ruleAction85]() {
					goto l1040
				}
				depth--
				add(ruleIdentifier, position1041)
			}
			return true
		l1040:
			position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
			return false
		},
		/* 104 TargetIdentifier <- <(<jsonPath> Action86)> */
		func() bool {
			position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
			{
				position1044 := position
				depth++
				{
					position1045 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1043
					}
					depth--
					add(rulePegText, position1045)
				}
				if !_rules[ruleAction86]() {
					goto l1043
				}
				depth--
				add(ruleTargetIdentifier, position1044)
			}
			return true
		l1043:
			position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
			return false
		},
		/* 105 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
			{
				position1047 := position
				depth++
				{
					position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1049
					}
					position++
					goto l1048
				l1049:
					position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1046
					}
					position++
				}
			l1048:
			l1050:
				{
					position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
					{
						position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1053
						}
						position++
						goto l1052
					l1053:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1054
						}
						position++
						goto l1052
					l1054:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1055
						}
						position++
						goto l1052
					l1055:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
						if buffer[position] != rune('_') {
							goto l1051
						}
						position++
					}
				l1052:
					goto l1050
				l1051:
					position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
				}
				depth--
				add(ruleident, position1047)
			}
			return true
		l1046:
			position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
			return false
		},
		/* 106 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
			{
				position1057 := position
				depth++
				{
					position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1059
					}
					position++
					goto l1058
				l1059:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1056
					}
					position++
				}
			l1058:
			l1060:
				{
					position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
					{
						position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1063
						}
						position++
						goto l1062
					l1063:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1064
						}
						position++
						goto l1062
					l1064:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1065
						}
						position++
						goto l1062
					l1065:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune('_') {
							goto l1066
						}
						position++
						goto l1062
					l1066:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune('.') {
							goto l1067
						}
						position++
						goto l1062
					l1067:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune('[') {
							goto l1068
						}
						position++
						goto l1062
					l1068:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune(']') {
							goto l1069
						}
						position++
						goto l1062
					l1069:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune('"') {
							goto l1061
						}
						position++
					}
				l1062:
					goto l1060
				l1061:
					position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
				}
				depth--
				add(rulejsonPath, position1057)
			}
			return true
		l1056:
			position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
			return false
		},
		/* 107 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1071 := position
				depth++
			l1072:
				{
					position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
					{
						position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1075
						}
						position++
						goto l1074
					l1075:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('\t') {
							goto l1076
						}
						position++
						goto l1074
					l1076:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('\n') {
							goto l1077
						}
						position++
						goto l1074
					l1077:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('\r') {
							goto l1078
						}
						position++
						goto l1074
					l1078:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if !_rules[rulecomment]() {
							goto l1073
						}
					}
				l1074:
					goto l1072
				l1073:
					position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
				}
				depth--
				add(rulesp, position1071)
			}
			return true
		},
		/* 108 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
			{
				position1080 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1079
				}
				position++
				if buffer[position] != rune('-') {
					goto l1079
				}
				position++
			l1081:
				{
					position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
					{
						position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
						{
							position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1085
							}
							position++
							goto l1084
						l1085:
							position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
							if buffer[position] != rune('\n') {
								goto l1083
							}
							position++
						}
					l1084:
						goto l1082
					l1083:
						position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
					}
					if !matchDot() {
						goto l1082
					}
					goto l1081
				l1082:
					position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
				}
				{
					position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1087
					}
					position++
					goto l1086
				l1087:
					position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
					if buffer[position] != rune('\n') {
						goto l1079
					}
					position++
				}
			l1086:
				depth--
				add(rulecomment, position1080)
			}
			return true
		l1079:
			position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
			return false
		},
		/* 110 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 111 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 112 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 113 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 114 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 115 Action5 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 116 Action6 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 117 Action7 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 118 Action8 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 119 Action9 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 120 Action10 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 121 Action11 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 122 Action12 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 123 Action13 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 124 Action14 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 125 Action15 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 126 Action16 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		nil,
		/* 128 Action17 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 129 Action18 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 130 Action19 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 131 Action20 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 132 Action21 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 133 Action22 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 134 Action23 <- <{
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
		/* 135 Action24 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 136 Action25 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 137 Action26 <- <{
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
		/* 138 Action27 <- <{
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
		/* 139 Action28 <- <{
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
		/* 140 Action29 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 141 Action30 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 142 Action31 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 143 Action32 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 144 Action33 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 145 Action34 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 146 Action35 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 147 Action36 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 148 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 149 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 150 Action39 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 151 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 152 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 153 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 154 Action43 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 155 Action44 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 156 Action45 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 157 Action46 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 158 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 159 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 160 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 161 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 162 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 163 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 164 Action53 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 165 Action54 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 166 Action55 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 167 Action56 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 168 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 169 Action58 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 170 Action59 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 171 Action60 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 172 Action61 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 173 Action62 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 174 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 175 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 176 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 177 Action66 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 178 Action67 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 179 Action68 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 180 Action69 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 181 Action70 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 182 Action71 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 183 Action72 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 184 Action73 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 185 Action74 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 186 Action75 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 187 Action76 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 188 Action77 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 189 Action78 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 190 Action79 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 191 Action80 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 192 Action81 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 193 Action82 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 194 Action83 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 195 Action84 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 196 Action85 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 197 Action86 <- <{
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

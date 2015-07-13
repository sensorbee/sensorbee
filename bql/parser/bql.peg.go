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
	rulePegText
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
	ruleAction80
	ruleAction81
	ruleAction82
	ruleAction83
	ruleAction84

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
	"PegText",
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
	"Action80",
	"Action81",
	"Action82",
	"Action83",
	"Action84",

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
	rules  [194]func() bool
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

			p.AssembleInsertIntoSelect()

		case ruleAction7:

			p.AssembleInsertIntoFrom()

		case ruleAction8:

			p.AssemblePauseSource()

		case ruleAction9:

			p.AssembleResumeSource()

		case ruleAction10:

			p.AssembleRewindSource()

		case ruleAction11:

			p.AssembleDropSource()

		case ruleAction12:

			p.AssembleDropStream()

		case ruleAction13:

			p.AssembleDropSink()

		case ruleAction14:

			p.AssembleDropState()

		case ruleAction15:

			p.AssembleEmitter(begin, end)

		case ruleAction16:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction17:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction18:

			p.AssembleStreamEmitInterval()

		case ruleAction19:

			p.AssembleProjections(begin, end)

		case ruleAction20:

			p.AssembleAlias()

		case ruleAction21:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction22:

			p.AssembleInterval()

		case ruleAction23:

			p.AssembleInterval()

		case ruleAction24:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction25:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction26:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction27:

			p.EnsureAliasedStreamWindow()

		case ruleAction28:

			p.AssembleAliasedStreamWindow()

		case ruleAction29:

			p.AssembleStreamWindow()

		case ruleAction30:

			p.AssembleUDSFFuncApp()

		case ruleAction31:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction32:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction33:

			p.AssembleSourceSinkParam()

		case ruleAction34:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction38:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction39:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction40:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction41:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction42:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction43:

			p.AssembleFuncApp()

		case ruleAction44:

			p.AssembleExpressions(begin, end)

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction51:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction52:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction53:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction54:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction55:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction56:

			p.PushComponent(begin, end, Istream)

		case ruleAction57:

			p.PushComponent(begin, end, Dstream)

		case ruleAction58:

			p.PushComponent(begin, end, Rstream)

		case ruleAction59:

			p.PushComponent(begin, end, Tuples)

		case ruleAction60:

			p.PushComponent(begin, end, Seconds)

		case ruleAction61:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction62:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction63:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction64:

			p.PushComponent(begin, end, Yes)

		case ruleAction65:

			p.PushComponent(begin, end, No)

		case ruleAction66:

			p.PushComponent(begin, end, Or)

		case ruleAction67:

			p.PushComponent(begin, end, And)

		case ruleAction68:

			p.PushComponent(begin, end, Not)

		case ruleAction69:

			p.PushComponent(begin, end, Equal)

		case ruleAction70:

			p.PushComponent(begin, end, Less)

		case ruleAction71:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction72:

			p.PushComponent(begin, end, Greater)

		case ruleAction73:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction74:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction75:

			p.PushComponent(begin, end, Is)

		case ruleAction76:

			p.PushComponent(begin, end, IsNot)

		case ruleAction77:

			p.PushComponent(begin, end, Plus)

		case ruleAction78:

			p.PushComponent(begin, end, Minus)

		case ruleAction79:

			p.PushComponent(begin, end, Multiply)

		case ruleAction80:

			p.PushComponent(begin, end, Divide)

		case ruleAction81:

			p.PushComponent(begin, end, Modulo)

		case ruleAction82:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction83:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction84:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt / DropSourceStmt / DropStreamStmt / DropSinkStmt / DropStateStmt / UpdateStateStmt)> */
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
			position24, tokenIndex24, depth24 := position, tokenIndex, depth
			{
				position25 := position
				depth++
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l27
					}
					position++
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					if buffer[position] != rune('S') {
						goto l24
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
						goto l24
					}
					position++
				}
			l28:
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l31
					}
					position++
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if buffer[position] != rune('L') {
						goto l24
					}
					position++
				}
			l30:
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l33
					}
					position++
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if buffer[position] != rune('E') {
						goto l24
					}
					position++
				}
			l32:
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
						goto l24
					}
					position++
				}
			l34:
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l37
					}
					position++
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if buffer[position] != rune('T') {
						goto l24
					}
					position++
				}
			l36:
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleEmitter]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleProjections]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleWindowedFrom]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleFilter]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleGrouping]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleHaving]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleAction0]() {
					goto l24
				}
				depth--
				add(ruleSelectStmt, position25)
			}
			return true
		l24:
			position, tokenIndex, depth = position24, tokenIndex24, depth24
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position38, tokenIndex38, depth38 := position, tokenIndex, depth
			{
				position39 := position
				depth++
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('C') {
						goto l38
					}
					position++
				}
			l40:
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('R') {
						goto l38
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
						goto l38
					}
					position++
				}
			l44:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('A') {
						goto l38
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
						goto l38
					}
					position++
				}
			l48:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('E') {
						goto l38
					}
					position++
				}
			l50:
				if !_rules[rulesp]() {
					goto l38
				}
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('S') {
						goto l38
					}
					position++
				}
			l52:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('T') {
						goto l38
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
						goto l38
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
						goto l38
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
						goto l38
					}
					position++
				}
			l60:
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
					if buffer[position] != rune('M') {
						goto l38
					}
					position++
				}
			l62:
				if !_rules[rulesp]() {
					goto l38
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l38
				}
				if !_rules[rulesp]() {
					goto l38
				}
				{
					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l65
					}
					position++
					goto l64
				l65:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
					if buffer[position] != rune('A') {
						goto l38
					}
					position++
				}
			l64:
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
						goto l38
					}
					position++
				}
			l66:
				if !_rules[rulesp]() {
					goto l38
				}
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
						goto l38
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
						goto l38
					}
					position++
				}
			l70:
				{
					position72, tokenIndex72, depth72 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l73
					}
					position++
					goto l72
				l73:
					position, tokenIndex, depth = position72, tokenIndex72, depth72
					if buffer[position] != rune('L') {
						goto l38
					}
					position++
				}
			l72:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if buffer[position] != rune('E') {
						goto l38
					}
					position++
				}
			l74:
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
						goto l38
					}
					position++
				}
			l76:
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l79
					}
					position++
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('T') {
						goto l38
					}
					position++
				}
			l78:
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
				if !_rules[ruleAction1]() {
					goto l38
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position39)
			}
			return true
		l38:
			position, tokenIndex, depth = position38, tokenIndex38, depth38
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position80, tokenIndex80, depth80 := position, tokenIndex, depth
			{
				position81 := position
				depth++
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l83
					}
					position++
					goto l82
				l83:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
					if buffer[position] != rune('C') {
						goto l80
					}
					position++
				}
			l82:
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
					if buffer[position] != rune('R') {
						goto l80
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
						goto l80
					}
					position++
				}
			l86:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l89
					}
					position++
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if buffer[position] != rune('A') {
						goto l80
					}
					position++
				}
			l88:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('T') {
						goto l80
					}
					position++
				}
			l90:
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('E') {
						goto l80
					}
					position++
				}
			l92:
				if !_rules[rulesp]() {
					goto l80
				}
				if !_rules[rulePausedOpt]() {
					goto l80
				}
				if !_rules[rulesp]() {
					goto l80
				}
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('S') {
						goto l80
					}
					position++
				}
			l94:
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('O') {
						goto l80
					}
					position++
				}
			l96:
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('U') {
						goto l80
					}
					position++
				}
			l98:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('R') {
						goto l80
					}
					position++
				}
			l100:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('C') {
						goto l80
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('E') {
						goto l80
					}
					position++
				}
			l104:
				if !_rules[rulesp]() {
					goto l80
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l80
				}
				if !_rules[rulesp]() {
					goto l80
				}
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('T') {
						goto l80
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('Y') {
						goto l80
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('P') {
						goto l80
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('E') {
						goto l80
					}
					position++
				}
			l112:
				if !_rules[rulesp]() {
					goto l80
				}
				if !_rules[ruleSourceSinkType]() {
					goto l80
				}
				if !_rules[rulesp]() {
					goto l80
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l80
				}
				if !_rules[ruleAction2]() {
					goto l80
				}
				depth--
				add(ruleCreateSourceStmt, position81)
			}
			return true
		l80:
			position, tokenIndex, depth = position80, tokenIndex80, depth80
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position114, tokenIndex114, depth114 := position, tokenIndex, depth
			{
				position115 := position
				depth++
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
						goto l114
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('R') {
						goto l114
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
						goto l114
					}
					position++
				}
			l120:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('A') {
						goto l114
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('T') {
						goto l114
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
						goto l114
					}
					position++
				}
			l126:
				if !_rules[rulesp]() {
					goto l114
				}
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('S') {
						goto l114
					}
					position++
				}
			l128:
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('I') {
						goto l114
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('N') {
						goto l114
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('K') {
						goto l114
					}
					position++
				}
			l134:
				if !_rules[rulesp]() {
					goto l114
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l114
				}
				if !_rules[rulesp]() {
					goto l114
				}
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('T') {
						goto l114
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('Y') {
						goto l114
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('P') {
						goto l114
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('E') {
						goto l114
					}
					position++
				}
			l142:
				if !_rules[rulesp]() {
					goto l114
				}
				if !_rules[ruleSourceSinkType]() {
					goto l114
				}
				if !_rules[rulesp]() {
					goto l114
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l114
				}
				if !_rules[ruleAction3]() {
					goto l114
				}
				depth--
				add(ruleCreateSinkStmt, position115)
			}
			return true
		l114:
			position, tokenIndex, depth = position114, tokenIndex114, depth114
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('C') {
						goto l144
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('R') {
						goto l144
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
						goto l144
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('A') {
						goto l144
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
						goto l144
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
						goto l144
					}
					position++
				}
			l156:
				if !_rules[rulesp]() {
					goto l144
				}
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if buffer[position] != rune('S') {
						goto l144
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
						goto l144
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('A') {
						goto l144
					}
					position++
				}
			l162:
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
						goto l144
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('E') {
						goto l144
					}
					position++
				}
			l166:
				if !_rules[rulesp]() {
					goto l144
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l144
				}
				if !_rules[rulesp]() {
					goto l144
				}
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
						goto l144
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('Y') {
						goto l144
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('P') {
						goto l144
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('E') {
						goto l144
					}
					position++
				}
			l174:
				if !_rules[rulesp]() {
					goto l144
				}
				if !_rules[ruleSourceSinkType]() {
					goto l144
				}
				if !_rules[rulesp]() {
					goto l144
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l144
				}
				if !_rules[ruleAction4]() {
					goto l144
				}
				depth--
				add(ruleCreateStateStmt, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 7 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action5)> */
		func() bool {
			position176, tokenIndex176, depth176 := position, tokenIndex, depth
			{
				position177 := position
				depth++
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('U') {
						goto l176
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('P') {
						goto l176
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('D') {
						goto l176
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('A') {
						goto l176
					}
					position++
				}
			l184:
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('T') {
						goto l176
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
						goto l176
					}
					position++
				}
			l188:
				if !_rules[rulesp]() {
					goto l176
				}
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('S') {
						goto l176
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('T') {
						goto l176
					}
					position++
				}
			l192:
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('A') {
						goto l176
					}
					position++
				}
			l194:
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('T') {
						goto l176
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('E') {
						goto l176
					}
					position++
				}
			l198:
				if !_rules[rulesp]() {
					goto l176
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l176
				}
				if !_rules[rulesp]() {
					goto l176
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l176
				}
				if !_rules[ruleAction5]() {
					goto l176
				}
				depth--
				add(ruleUpdateStateStmt, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 8 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action6)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('I') {
						goto l200
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('N') {
						goto l200
					}
					position++
				}
			l204:
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
						goto l200
					}
					position++
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('E') {
						goto l200
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('R') {
						goto l200
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
						goto l200
					}
					position++
				}
			l212:
				if !_rules[rulesp]() {
					goto l200
				}
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('I') {
						goto l200
					}
					position++
				}
			l214:
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('N') {
						goto l200
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('T') {
						goto l200
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
						goto l200
					}
					position++
				}
			l220:
				if !_rules[rulesp]() {
					goto l200
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l200
				}
				if !_rules[rulesp]() {
					goto l200
				}
				if !_rules[ruleSelectStmt]() {
					goto l200
				}
				if !_rules[ruleAction6]() {
					goto l200
				}
				depth--
				add(ruleInsertIntoSelectStmt, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 9 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action7)> */
		func() bool {
			position222, tokenIndex222, depth222 := position, tokenIndex, depth
			{
				position223 := position
				depth++
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('I') {
						goto l222
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('N') {
						goto l222
					}
					position++
				}
			l226:
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
						goto l222
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
						goto l222
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('R') {
						goto l222
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('T') {
						goto l222
					}
					position++
				}
			l234:
				if !_rules[rulesp]() {
					goto l222
				}
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('I') {
						goto l222
					}
					position++
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('N') {
						goto l222
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('T') {
						goto l222
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('O') {
						goto l222
					}
					position++
				}
			l242:
				if !_rules[rulesp]() {
					goto l222
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l222
				}
				if !_rules[rulesp]() {
					goto l222
				}
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('F') {
						goto l222
					}
					position++
				}
			l244:
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
						goto l222
					}
					position++
				}
			l246:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('O') {
						goto l222
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('M') {
						goto l222
					}
					position++
				}
			l250:
				if !_rules[rulesp]() {
					goto l222
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l222
				}
				if !_rules[ruleAction7]() {
					goto l222
				}
				depth--
				add(ruleInsertIntoFromStmt, position223)
			}
			return true
		l222:
			position, tokenIndex, depth = position222, tokenIndex222, depth222
			return false
		},
		/* 10 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action8)> */
		func() bool {
			position252, tokenIndex252, depth252 := position, tokenIndex, depth
			{
				position253 := position
				depth++
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if buffer[position] != rune('P') {
						goto l252
					}
					position++
				}
			l254:
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l257
					}
					position++
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if buffer[position] != rune('A') {
						goto l252
					}
					position++
				}
			l256:
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if buffer[position] != rune('U') {
						goto l252
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('S') {
						goto l252
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('E') {
						goto l252
					}
					position++
				}
			l262:
				if !_rules[rulesp]() {
					goto l252
				}
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('S') {
						goto l252
					}
					position++
				}
			l264:
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('O') {
						goto l252
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('U') {
						goto l252
					}
					position++
				}
			l268:
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if buffer[position] != rune('R') {
						goto l252
					}
					position++
				}
			l270:
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l273
					}
					position++
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if buffer[position] != rune('C') {
						goto l252
					}
					position++
				}
			l272:
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l275
					}
					position++
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('E') {
						goto l252
					}
					position++
				}
			l274:
				if !_rules[rulesp]() {
					goto l252
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l252
				}
				if !_rules[ruleAction8]() {
					goto l252
				}
				depth--
				add(rulePauseSourceStmt, position253)
			}
			return true
		l252:
			position, tokenIndex, depth = position252, tokenIndex252, depth252
			return false
		},
		/* 11 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action9)> */
		func() bool {
			position276, tokenIndex276, depth276 := position, tokenIndex, depth
			{
				position277 := position
				depth++
				{
					position278, tokenIndex278, depth278 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l279
					}
					position++
					goto l278
				l279:
					position, tokenIndex, depth = position278, tokenIndex278, depth278
					if buffer[position] != rune('R') {
						goto l276
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
						goto l276
					}
					position++
				}
			l280:
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('S') {
						goto l276
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
					if buffer[position] != rune('U') {
						goto l276
					}
					position++
				}
			l284:
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l287
					}
					position++
					goto l286
				l287:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
					if buffer[position] != rune('M') {
						goto l276
					}
					position++
				}
			l286:
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l289
					}
					position++
					goto l288
				l289:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
					if buffer[position] != rune('E') {
						goto l276
					}
					position++
				}
			l288:
				if !_rules[rulesp]() {
					goto l276
				}
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
						goto l276
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
						goto l276
					}
					position++
				}
			l292:
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l295
					}
					position++
					goto l294
				l295:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if buffer[position] != rune('U') {
						goto l276
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
						goto l276
					}
					position++
				}
			l296:
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					if buffer[position] != rune('C') {
						goto l276
					}
					position++
				}
			l298:
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l301
					}
					position++
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if buffer[position] != rune('E') {
						goto l276
					}
					position++
				}
			l300:
				if !_rules[rulesp]() {
					goto l276
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l276
				}
				if !_rules[ruleAction9]() {
					goto l276
				}
				depth--
				add(ruleResumeSourceStmt, position277)
			}
			return true
		l276:
			position, tokenIndex, depth = position276, tokenIndex276, depth276
			return false
		},
		/* 12 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position302, tokenIndex302, depth302 := position, tokenIndex, depth
			{
				position303 := position
				depth++
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
						goto l302
					}
					position++
				}
			l304:
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l307
					}
					position++
					goto l306
				l307:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					if buffer[position] != rune('E') {
						goto l302
					}
					position++
				}
			l306:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					if buffer[position] != rune('W') {
						goto l302
					}
					position++
				}
			l308:
				{
					position310, tokenIndex310, depth310 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l311
					}
					position++
					goto l310
				l311:
					position, tokenIndex, depth = position310, tokenIndex310, depth310
					if buffer[position] != rune('I') {
						goto l302
					}
					position++
				}
			l310:
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('N') {
						goto l302
					}
					position++
				}
			l312:
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l315
					}
					position++
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if buffer[position] != rune('D') {
						goto l302
					}
					position++
				}
			l314:
				if !_rules[rulesp]() {
					goto l302
				}
				{
					position316, tokenIndex316, depth316 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l317
					}
					position++
					goto l316
				l317:
					position, tokenIndex, depth = position316, tokenIndex316, depth316
					if buffer[position] != rune('S') {
						goto l302
					}
					position++
				}
			l316:
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
					if buffer[position] != rune('O') {
						goto l302
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
						goto l302
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if buffer[position] != rune('R') {
						goto l302
					}
					position++
				}
			l322:
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l325
					}
					position++
					goto l324
				l325:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if buffer[position] != rune('C') {
						goto l302
					}
					position++
				}
			l324:
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l327
					}
					position++
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					if buffer[position] != rune('E') {
						goto l302
					}
					position++
				}
			l326:
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
				add(ruleRewindSourceStmt, position303)
			}
			return true
		l302:
			position, tokenIndex, depth = position302, tokenIndex302, depth302
			return false
		},
		/* 13 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action11)> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l331
					}
					position++
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					if buffer[position] != rune('D') {
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
					if buffer[position] != rune('p') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if buffer[position] != rune('P') {
						goto l328
					}
					position++
				}
			l336:
				if !_rules[rulesp]() {
					goto l328
				}
				{
					position338, tokenIndex338, depth338 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l339
					}
					position++
					goto l338
				l339:
					position, tokenIndex, depth = position338, tokenIndex338, depth338
					if buffer[position] != rune('S') {
						goto l328
					}
					position++
				}
			l338:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l341
					}
					position++
					goto l340
				l341:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
					if buffer[position] != rune('O') {
						goto l328
					}
					position++
				}
			l340:
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if buffer[position] != rune('U') {
						goto l328
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
					if buffer[position] != rune('R') {
						goto l328
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					if buffer[position] != rune('C') {
						goto l328
					}
					position++
				}
			l346:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if buffer[position] != rune('E') {
						goto l328
					}
					position++
				}
			l348:
				if !_rules[rulesp]() {
					goto l328
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l328
				}
				if !_rules[ruleAction11]() {
					goto l328
				}
				depth--
				add(ruleDropSourceStmt, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 14 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action12)> */
		func() bool {
			position350, tokenIndex350, depth350 := position, tokenIndex, depth
			{
				position351 := position
				depth++
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('D') {
						goto l350
					}
					position++
				}
			l352:
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
						goto l350
					}
					position++
				}
			l354:
				{
					position356, tokenIndex356, depth356 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l357
					}
					position++
					goto l356
				l357:
					position, tokenIndex, depth = position356, tokenIndex356, depth356
					if buffer[position] != rune('O') {
						goto l350
					}
					position++
				}
			l356:
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l359
					}
					position++
					goto l358
				l359:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					if buffer[position] != rune('P') {
						goto l350
					}
					position++
				}
			l358:
				if !_rules[rulesp]() {
					goto l350
				}
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if buffer[position] != rune('S') {
						goto l350
					}
					position++
				}
			l360:
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l363
					}
					position++
					goto l362
				l363:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
					if buffer[position] != rune('T') {
						goto l350
					}
					position++
				}
			l362:
				{
					position364, tokenIndex364, depth364 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l365
					}
					position++
					goto l364
				l365:
					position, tokenIndex, depth = position364, tokenIndex364, depth364
					if buffer[position] != rune('R') {
						goto l350
					}
					position++
				}
			l364:
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l367
					}
					position++
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if buffer[position] != rune('E') {
						goto l350
					}
					position++
				}
			l366:
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if buffer[position] != rune('A') {
						goto l350
					}
					position++
				}
			l368:
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l371
					}
					position++
					goto l370
				l371:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
					if buffer[position] != rune('M') {
						goto l350
					}
					position++
				}
			l370:
				if !_rules[rulesp]() {
					goto l350
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l350
				}
				if !_rules[ruleAction12]() {
					goto l350
				}
				depth--
				add(ruleDropStreamStmt, position351)
			}
			return true
		l350:
			position, tokenIndex, depth = position350, tokenIndex350, depth350
			return false
		},
		/* 15 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action13)> */
		func() bool {
			position372, tokenIndex372, depth372 := position, tokenIndex, depth
			{
				position373 := position
				depth++
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('D') {
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
					if buffer[position] != rune('p') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('P') {
						goto l372
					}
					position++
				}
			l380:
				if !_rules[rulesp]() {
					goto l372
				}
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l383
					}
					position++
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if buffer[position] != rune('S') {
						goto l372
					}
					position++
				}
			l382:
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l385
					}
					position++
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if buffer[position] != rune('I') {
						goto l372
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('N') {
						goto l372
					}
					position++
				}
			l386:
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l389
					}
					position++
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					if buffer[position] != rune('K') {
						goto l372
					}
					position++
				}
			l388:
				if !_rules[rulesp]() {
					goto l372
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l372
				}
				if !_rules[ruleAction13]() {
					goto l372
				}
				depth--
				add(ruleDropSinkStmt, position373)
			}
			return true
		l372:
			position, tokenIndex, depth = position372, tokenIndex372, depth372
			return false
		},
		/* 16 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action14)> */
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
					if buffer[position] != rune('t') {
						goto l403
					}
					position++
					goto l402
				l403:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
					if buffer[position] != rune('T') {
						goto l390
					}
					position++
				}
			l402:
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l405
					}
					position++
					goto l404
				l405:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					if buffer[position] != rune('A') {
						goto l390
					}
					position++
				}
			l404:
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l407
					}
					position++
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if buffer[position] != rune('T') {
						goto l390
					}
					position++
				}
			l406:
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l409
					}
					position++
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if buffer[position] != rune('E') {
						goto l390
					}
					position++
				}
			l408:
				if !_rules[rulesp]() {
					goto l390
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l390
				}
				if !_rules[ruleAction14]() {
					goto l390
				}
				depth--
				add(ruleDropStateStmt, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 17 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action15)> */
		func() bool {
			position410, tokenIndex410, depth410 := position, tokenIndex, depth
			{
				position411 := position
				depth++
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l413
					}
					goto l412
				l413:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
					if !_rules[ruleDSTREAM]() {
						goto l414
					}
					goto l412
				l414:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
					if !_rules[ruleRSTREAM]() {
						goto l410
					}
				}
			l412:
				{
					position415 := position
					depth++
					{
						position416, tokenIndex416, depth416 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l416
						}
						if buffer[position] != rune('[') {
							goto l416
						}
						position++
						if !_rules[rulesp]() {
							goto l416
						}
						{
							position418, tokenIndex418, depth418 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l419
							}
							position++
							goto l418
						l419:
							position, tokenIndex, depth = position418, tokenIndex418, depth418
							if buffer[position] != rune('E') {
								goto l416
							}
							position++
						}
					l418:
						{
							position420, tokenIndex420, depth420 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l421
							}
							position++
							goto l420
						l421:
							position, tokenIndex, depth = position420, tokenIndex420, depth420
							if buffer[position] != rune('V') {
								goto l416
							}
							position++
						}
					l420:
						{
							position422, tokenIndex422, depth422 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l423
							}
							position++
							goto l422
						l423:
							position, tokenIndex, depth = position422, tokenIndex422, depth422
							if buffer[position] != rune('E') {
								goto l416
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
								goto l416
							}
							position++
						}
					l424:
						{
							position426, tokenIndex426, depth426 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l427
							}
							position++
							goto l426
						l427:
							position, tokenIndex, depth = position426, tokenIndex426, depth426
							if buffer[position] != rune('Y') {
								goto l416
							}
							position++
						}
					l426:
						if !_rules[rulesp]() {
							goto l416
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l416
						}
						if !_rules[rulesp]() {
							goto l416
						}
						if buffer[position] != rune(']') {
							goto l416
						}
						position++
						goto l417
					l416:
						position, tokenIndex, depth = position416, tokenIndex416, depth416
					}
				l417:
					depth--
					add(rulePegText, position415)
				}
				if !_rules[ruleAction15]() {
					goto l410
				}
				depth--
				add(ruleEmitter, position411)
			}
			return true
		l410:
			position, tokenIndex, depth = position410, tokenIndex410, depth410
			return false
		},
		/* 18 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position428, tokenIndex428, depth428 := position, tokenIndex, depth
			{
				position429 := position
				depth++
				{
					position430, tokenIndex430, depth430 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l431
					}
				l432:
					{
						position433, tokenIndex433, depth433 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l433
						}
						if buffer[position] != rune(',') {
							goto l433
						}
						position++
						if !_rules[rulesp]() {
							goto l433
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l433
						}
						goto l432
					l433:
						position, tokenIndex, depth = position433, tokenIndex433, depth433
					}
					goto l430
				l431:
					position, tokenIndex, depth = position430, tokenIndex430, depth430
					if !_rules[ruleTimeEmitterInterval]() {
						goto l434
					}
					goto l430
				l434:
					position, tokenIndex, depth = position430, tokenIndex430, depth430
					if !_rules[ruleTupleEmitterInterval]() {
						goto l428
					}
				}
			l430:
				depth--
				add(ruleEmitterIntervals, position429)
			}
			return true
		l428:
			position, tokenIndex, depth = position428, tokenIndex428, depth428
			return false
		},
		/* 19 TimeEmitterInterval <- <(<TimeInterval> Action16)> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				{
					position437 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l435
					}
					depth--
					add(rulePegText, position437)
				}
				if !_rules[ruleAction16]() {
					goto l435
				}
				depth--
				add(ruleTimeEmitterInterval, position436)
			}
			return true
		l435:
			position, tokenIndex, depth = position435, tokenIndex435, depth435
			return false
		},
		/* 20 TupleEmitterInterval <- <(<TuplesInterval> Action17)> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				{
					position440 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l438
					}
					depth--
					add(rulePegText, position440)
				}
				if !_rules[ruleAction17]() {
					goto l438
				}
				depth--
				add(ruleTupleEmitterInterval, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 21 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action18)> */
		func() bool {
			position441, tokenIndex441, depth441 := position, tokenIndex, depth
			{
				position442 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l441
				}
				if !_rules[rulesp]() {
					goto l441
				}
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l444
					}
					position++
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if buffer[position] != rune('I') {
						goto l441
					}
					position++
				}
			l443:
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l446
					}
					position++
					goto l445
				l446:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
					if buffer[position] != rune('N') {
						goto l441
					}
					position++
				}
			l445:
				if !_rules[rulesp]() {
					goto l441
				}
				if !_rules[ruleStream]() {
					goto l441
				}
				if !_rules[ruleAction18]() {
					goto l441
				}
				depth--
				add(ruleTupleEmitterFromInterval, position442)
			}
			return true
		l441:
			position, tokenIndex, depth = position441, tokenIndex441, depth441
			return false
		},
		/* 22 Projections <- <(<(Projection sp (',' sp Projection)*)> Action19)> */
		func() bool {
			position447, tokenIndex447, depth447 := position, tokenIndex, depth
			{
				position448 := position
				depth++
				{
					position449 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l447
					}
					if !_rules[rulesp]() {
						goto l447
					}
				l450:
					{
						position451, tokenIndex451, depth451 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l451
						}
						position++
						if !_rules[rulesp]() {
							goto l451
						}
						if !_rules[ruleProjection]() {
							goto l451
						}
						goto l450
					l451:
						position, tokenIndex, depth = position451, tokenIndex451, depth451
					}
					depth--
					add(rulePegText, position449)
				}
				if !_rules[ruleAction19]() {
					goto l447
				}
				depth--
				add(ruleProjections, position448)
			}
			return true
		l447:
			position, tokenIndex, depth = position447, tokenIndex447, depth447
			return false
		},
		/* 23 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position452, tokenIndex452, depth452 := position, tokenIndex, depth
			{
				position453 := position
				depth++
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l455
					}
					goto l454
				l455:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if !_rules[ruleExpression]() {
						goto l456
					}
					goto l454
				l456:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if !_rules[ruleWildcard]() {
						goto l452
					}
				}
			l454:
				depth--
				add(ruleProjection, position453)
			}
			return true
		l452:
			position, tokenIndex, depth = position452, tokenIndex452, depth452
			return false
		},
		/* 24 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action20)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l460
					}
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleWildcard]() {
						goto l457
					}
				}
			l459:
				if !_rules[rulesp]() {
					goto l457
				}
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l462
					}
					position++
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('A') {
						goto l457
					}
					position++
				}
			l461:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('S') {
						goto l457
					}
					position++
				}
			l463:
				if !_rules[rulesp]() {
					goto l457
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l457
				}
				if !_rules[ruleAction20]() {
					goto l457
				}
				depth--
				add(ruleAliasExpression, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 25 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action21)> */
		func() bool {
			position465, tokenIndex465, depth465 := position, tokenIndex, depth
			{
				position466 := position
				depth++
				{
					position467 := position
					depth++
					{
						position468, tokenIndex468, depth468 := position, tokenIndex, depth
						{
							position470, tokenIndex470, depth470 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l471
							}
							position++
							goto l470
						l471:
							position, tokenIndex, depth = position470, tokenIndex470, depth470
							if buffer[position] != rune('F') {
								goto l468
							}
							position++
						}
					l470:
						{
							position472, tokenIndex472, depth472 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l473
							}
							position++
							goto l472
						l473:
							position, tokenIndex, depth = position472, tokenIndex472, depth472
							if buffer[position] != rune('R') {
								goto l468
							}
							position++
						}
					l472:
						{
							position474, tokenIndex474, depth474 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l475
							}
							position++
							goto l474
						l475:
							position, tokenIndex, depth = position474, tokenIndex474, depth474
							if buffer[position] != rune('O') {
								goto l468
							}
							position++
						}
					l474:
						{
							position476, tokenIndex476, depth476 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l477
							}
							position++
							goto l476
						l477:
							position, tokenIndex, depth = position476, tokenIndex476, depth476
							if buffer[position] != rune('M') {
								goto l468
							}
							position++
						}
					l476:
						if !_rules[rulesp]() {
							goto l468
						}
						if !_rules[ruleRelations]() {
							goto l468
						}
						if !_rules[rulesp]() {
							goto l468
						}
						goto l469
					l468:
						position, tokenIndex, depth = position468, tokenIndex468, depth468
					}
				l469:
					depth--
					add(rulePegText, position467)
				}
				if !_rules[ruleAction21]() {
					goto l465
				}
				depth--
				add(ruleWindowedFrom, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 26 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				{
					position480, tokenIndex480, depth480 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l481
					}
					goto l480
				l481:
					position, tokenIndex, depth = position480, tokenIndex480, depth480
					if !_rules[ruleTuplesInterval]() {
						goto l478
					}
				}
			l480:
				depth--
				add(ruleInterval, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 27 TimeInterval <- <(NumericLiteral sp SECONDS Action22)> */
		func() bool {
			position482, tokenIndex482, depth482 := position, tokenIndex, depth
			{
				position483 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l482
				}
				if !_rules[rulesp]() {
					goto l482
				}
				if !_rules[ruleSECONDS]() {
					goto l482
				}
				if !_rules[ruleAction22]() {
					goto l482
				}
				depth--
				add(ruleTimeInterval, position483)
			}
			return true
		l482:
			position, tokenIndex, depth = position482, tokenIndex482, depth482
			return false
		},
		/* 28 TuplesInterval <- <(NumericLiteral sp TUPLES Action23)> */
		func() bool {
			position484, tokenIndex484, depth484 := position, tokenIndex, depth
			{
				position485 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l484
				}
				if !_rules[rulesp]() {
					goto l484
				}
				if !_rules[ruleTUPLES]() {
					goto l484
				}
				if !_rules[ruleAction23]() {
					goto l484
				}
				depth--
				add(ruleTuplesInterval, position485)
			}
			return true
		l484:
			position, tokenIndex, depth = position484, tokenIndex484, depth484
			return false
		},
		/* 29 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position486, tokenIndex486, depth486 := position, tokenIndex, depth
			{
				position487 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l486
				}
				if !_rules[rulesp]() {
					goto l486
				}
			l488:
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l489
					}
					position++
					if !_rules[rulesp]() {
						goto l489
					}
					if !_rules[ruleRelationLike]() {
						goto l489
					}
					goto l488
				l489:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
				}
				depth--
				add(ruleRelations, position487)
			}
			return true
		l486:
			position, tokenIndex, depth = position486, tokenIndex486, depth486
			return false
		},
		/* 30 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action24)> */
		func() bool {
			position490, tokenIndex490, depth490 := position, tokenIndex, depth
			{
				position491 := position
				depth++
				{
					position492 := position
					depth++
					{
						position493, tokenIndex493, depth493 := position, tokenIndex, depth
						{
							position495, tokenIndex495, depth495 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l496
							}
							position++
							goto l495
						l496:
							position, tokenIndex, depth = position495, tokenIndex495, depth495
							if buffer[position] != rune('W') {
								goto l493
							}
							position++
						}
					l495:
						{
							position497, tokenIndex497, depth497 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l498
							}
							position++
							goto l497
						l498:
							position, tokenIndex, depth = position497, tokenIndex497, depth497
							if buffer[position] != rune('H') {
								goto l493
							}
							position++
						}
					l497:
						{
							position499, tokenIndex499, depth499 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l500
							}
							position++
							goto l499
						l500:
							position, tokenIndex, depth = position499, tokenIndex499, depth499
							if buffer[position] != rune('E') {
								goto l493
							}
							position++
						}
					l499:
						{
							position501, tokenIndex501, depth501 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l502
							}
							position++
							goto l501
						l502:
							position, tokenIndex, depth = position501, tokenIndex501, depth501
							if buffer[position] != rune('R') {
								goto l493
							}
							position++
						}
					l501:
						{
							position503, tokenIndex503, depth503 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l504
							}
							position++
							goto l503
						l504:
							position, tokenIndex, depth = position503, tokenIndex503, depth503
							if buffer[position] != rune('E') {
								goto l493
							}
							position++
						}
					l503:
						if !_rules[rulesp]() {
							goto l493
						}
						if !_rules[ruleExpression]() {
							goto l493
						}
						goto l494
					l493:
						position, tokenIndex, depth = position493, tokenIndex493, depth493
					}
				l494:
					depth--
					add(rulePegText, position492)
				}
				if !_rules[ruleAction24]() {
					goto l490
				}
				depth--
				add(ruleFilter, position491)
			}
			return true
		l490:
			position, tokenIndex, depth = position490, tokenIndex490, depth490
			return false
		},
		/* 31 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action25)> */
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
							if buffer[position] != rune('g') {
								goto l511
							}
							position++
							goto l510
						l511:
							position, tokenIndex, depth = position510, tokenIndex510, depth510
							if buffer[position] != rune('G') {
								goto l508
							}
							position++
						}
					l510:
						{
							position512, tokenIndex512, depth512 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l513
							}
							position++
							goto l512
						l513:
							position, tokenIndex, depth = position512, tokenIndex512, depth512
							if buffer[position] != rune('R') {
								goto l508
							}
							position++
						}
					l512:
						{
							position514, tokenIndex514, depth514 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l515
							}
							position++
							goto l514
						l515:
							position, tokenIndex, depth = position514, tokenIndex514, depth514
							if buffer[position] != rune('O') {
								goto l508
							}
							position++
						}
					l514:
						{
							position516, tokenIndex516, depth516 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l517
							}
							position++
							goto l516
						l517:
							position, tokenIndex, depth = position516, tokenIndex516, depth516
							if buffer[position] != rune('U') {
								goto l508
							}
							position++
						}
					l516:
						{
							position518, tokenIndex518, depth518 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l519
							}
							position++
							goto l518
						l519:
							position, tokenIndex, depth = position518, tokenIndex518, depth518
							if buffer[position] != rune('P') {
								goto l508
							}
							position++
						}
					l518:
						if !_rules[rulesp]() {
							goto l508
						}
						{
							position520, tokenIndex520, depth520 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l521
							}
							position++
							goto l520
						l521:
							position, tokenIndex, depth = position520, tokenIndex520, depth520
							if buffer[position] != rune('B') {
								goto l508
							}
							position++
						}
					l520:
						{
							position522, tokenIndex522, depth522 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l523
							}
							position++
							goto l522
						l523:
							position, tokenIndex, depth = position522, tokenIndex522, depth522
							if buffer[position] != rune('Y') {
								goto l508
							}
							position++
						}
					l522:
						if !_rules[rulesp]() {
							goto l508
						}
						if !_rules[ruleGroupList]() {
							goto l508
						}
						goto l509
					l508:
						position, tokenIndex, depth = position508, tokenIndex508, depth508
					}
				l509:
					depth--
					add(rulePegText, position507)
				}
				if !_rules[ruleAction25]() {
					goto l505
				}
				depth--
				add(ruleGrouping, position506)
			}
			return true
		l505:
			position, tokenIndex, depth = position505, tokenIndex505, depth505
			return false
		},
		/* 32 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l524
				}
				if !_rules[rulesp]() {
					goto l524
				}
			l526:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l527
					}
					position++
					if !_rules[rulesp]() {
						goto l527
					}
					if !_rules[ruleExpression]() {
						goto l527
					}
					goto l526
				l527:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
				}
				depth--
				add(ruleGroupList, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 33 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action26)> */
		func() bool {
			position528, tokenIndex528, depth528 := position, tokenIndex, depth
			{
				position529 := position
				depth++
				{
					position530 := position
					depth++
					{
						position531, tokenIndex531, depth531 := position, tokenIndex, depth
						{
							position533, tokenIndex533, depth533 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l534
							}
							position++
							goto l533
						l534:
							position, tokenIndex, depth = position533, tokenIndex533, depth533
							if buffer[position] != rune('H') {
								goto l531
							}
							position++
						}
					l533:
						{
							position535, tokenIndex535, depth535 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l536
							}
							position++
							goto l535
						l536:
							position, tokenIndex, depth = position535, tokenIndex535, depth535
							if buffer[position] != rune('A') {
								goto l531
							}
							position++
						}
					l535:
						{
							position537, tokenIndex537, depth537 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l538
							}
							position++
							goto l537
						l538:
							position, tokenIndex, depth = position537, tokenIndex537, depth537
							if buffer[position] != rune('V') {
								goto l531
							}
							position++
						}
					l537:
						{
							position539, tokenIndex539, depth539 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l540
							}
							position++
							goto l539
						l540:
							position, tokenIndex, depth = position539, tokenIndex539, depth539
							if buffer[position] != rune('I') {
								goto l531
							}
							position++
						}
					l539:
						{
							position541, tokenIndex541, depth541 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l542
							}
							position++
							goto l541
						l542:
							position, tokenIndex, depth = position541, tokenIndex541, depth541
							if buffer[position] != rune('N') {
								goto l531
							}
							position++
						}
					l541:
						{
							position543, tokenIndex543, depth543 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l544
							}
							position++
							goto l543
						l544:
							position, tokenIndex, depth = position543, tokenIndex543, depth543
							if buffer[position] != rune('G') {
								goto l531
							}
							position++
						}
					l543:
						if !_rules[rulesp]() {
							goto l531
						}
						if !_rules[ruleExpression]() {
							goto l531
						}
						goto l532
					l531:
						position, tokenIndex, depth = position531, tokenIndex531, depth531
					}
				l532:
					depth--
					add(rulePegText, position530)
				}
				if !_rules[ruleAction26]() {
					goto l528
				}
				depth--
				add(ruleHaving, position529)
			}
			return true
		l528:
			position, tokenIndex, depth = position528, tokenIndex528, depth528
			return false
		},
		/* 34 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action27))> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				{
					position547, tokenIndex547, depth547 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l548
					}
					goto l547
				l548:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
					if !_rules[ruleStreamWindow]() {
						goto l545
					}
					if !_rules[ruleAction27]() {
						goto l545
					}
				}
			l547:
				depth--
				add(ruleRelationLike, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 35 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action28)> */
		func() bool {
			position549, tokenIndex549, depth549 := position, tokenIndex, depth
			{
				position550 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l549
				}
				if !_rules[rulesp]() {
					goto l549
				}
				{
					position551, tokenIndex551, depth551 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l552
					}
					position++
					goto l551
				l552:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if buffer[position] != rune('A') {
						goto l549
					}
					position++
				}
			l551:
				{
					position553, tokenIndex553, depth553 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l554
					}
					position++
					goto l553
				l554:
					position, tokenIndex, depth = position553, tokenIndex553, depth553
					if buffer[position] != rune('S') {
						goto l549
					}
					position++
				}
			l553:
				if !_rules[rulesp]() {
					goto l549
				}
				if !_rules[ruleIdentifier]() {
					goto l549
				}
				if !_rules[ruleAction28]() {
					goto l549
				}
				depth--
				add(ruleAliasedStreamWindow, position550)
			}
			return true
		l549:
			position, tokenIndex, depth = position549, tokenIndex549, depth549
			return false
		},
		/* 36 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action29)> */
		func() bool {
			position555, tokenIndex555, depth555 := position, tokenIndex, depth
			{
				position556 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l555
				}
				if !_rules[rulesp]() {
					goto l555
				}
				if buffer[position] != rune('[') {
					goto l555
				}
				position++
				if !_rules[rulesp]() {
					goto l555
				}
				{
					position557, tokenIndex557, depth557 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l558
					}
					position++
					goto l557
				l558:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if buffer[position] != rune('R') {
						goto l555
					}
					position++
				}
			l557:
				{
					position559, tokenIndex559, depth559 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l560
					}
					position++
					goto l559
				l560:
					position, tokenIndex, depth = position559, tokenIndex559, depth559
					if buffer[position] != rune('A') {
						goto l555
					}
					position++
				}
			l559:
				{
					position561, tokenIndex561, depth561 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l562
					}
					position++
					goto l561
				l562:
					position, tokenIndex, depth = position561, tokenIndex561, depth561
					if buffer[position] != rune('N') {
						goto l555
					}
					position++
				}
			l561:
				{
					position563, tokenIndex563, depth563 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l564
					}
					position++
					goto l563
				l564:
					position, tokenIndex, depth = position563, tokenIndex563, depth563
					if buffer[position] != rune('G') {
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
				if !_rules[ruleInterval]() {
					goto l555
				}
				if !_rules[rulesp]() {
					goto l555
				}
				if buffer[position] != rune(']') {
					goto l555
				}
				position++
				if !_rules[ruleAction29]() {
					goto l555
				}
				depth--
				add(ruleStreamWindow, position556)
			}
			return true
		l555:
			position, tokenIndex, depth = position555, tokenIndex555, depth555
			return false
		},
		/* 37 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position567, tokenIndex567, depth567 := position, tokenIndex, depth
			{
				position568 := position
				depth++
				{
					position569, tokenIndex569, depth569 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l570
					}
					goto l569
				l570:
					position, tokenIndex, depth = position569, tokenIndex569, depth569
					if !_rules[ruleStream]() {
						goto l567
					}
				}
			l569:
				depth--
				add(ruleStreamLike, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 38 UDSFFuncApp <- <(FuncApp Action30)> */
		func() bool {
			position571, tokenIndex571, depth571 := position, tokenIndex, depth
			{
				position572 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l571
				}
				if !_rules[ruleAction30]() {
					goto l571
				}
				depth--
				add(ruleUDSFFuncApp, position572)
			}
			return true
		l571:
			position, tokenIndex, depth = position571, tokenIndex571, depth571
			return false
		},
		/* 39 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action31)> */
		func() bool {
			position573, tokenIndex573, depth573 := position, tokenIndex, depth
			{
				position574 := position
				depth++
				{
					position575 := position
					depth++
					{
						position576, tokenIndex576, depth576 := position, tokenIndex, depth
						{
							position578, tokenIndex578, depth578 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l579
							}
							position++
							goto l578
						l579:
							position, tokenIndex, depth = position578, tokenIndex578, depth578
							if buffer[position] != rune('W') {
								goto l576
							}
							position++
						}
					l578:
						{
							position580, tokenIndex580, depth580 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l581
							}
							position++
							goto l580
						l581:
							position, tokenIndex, depth = position580, tokenIndex580, depth580
							if buffer[position] != rune('I') {
								goto l576
							}
							position++
						}
					l580:
						{
							position582, tokenIndex582, depth582 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l583
							}
							position++
							goto l582
						l583:
							position, tokenIndex, depth = position582, tokenIndex582, depth582
							if buffer[position] != rune('T') {
								goto l576
							}
							position++
						}
					l582:
						{
							position584, tokenIndex584, depth584 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l585
							}
							position++
							goto l584
						l585:
							position, tokenIndex, depth = position584, tokenIndex584, depth584
							if buffer[position] != rune('H') {
								goto l576
							}
							position++
						}
					l584:
						if !_rules[rulesp]() {
							goto l576
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l576
						}
						if !_rules[rulesp]() {
							goto l576
						}
					l586:
						{
							position587, tokenIndex587, depth587 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l587
							}
							position++
							if !_rules[rulesp]() {
								goto l587
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l587
							}
							goto l586
						l587:
							position, tokenIndex, depth = position587, tokenIndex587, depth587
						}
						goto l577
					l576:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
					}
				l577:
					depth--
					add(rulePegText, position575)
				}
				if !_rules[ruleAction31]() {
					goto l573
				}
				depth--
				add(ruleSourceSinkSpecs, position574)
			}
			return true
		l573:
			position, tokenIndex, depth = position573, tokenIndex573, depth573
			return false
		},
		/* 40 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action32)> */
		func() bool {
			position588, tokenIndex588, depth588 := position, tokenIndex, depth
			{
				position589 := position
				depth++
				{
					position590 := position
					depth++
					{
						position591, tokenIndex591, depth591 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l592
						}
						position++
						goto l591
					l592:
						position, tokenIndex, depth = position591, tokenIndex591, depth591
						if buffer[position] != rune('S') {
							goto l588
						}
						position++
					}
				l591:
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l594
						}
						position++
						goto l593
					l594:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
						if buffer[position] != rune('E') {
							goto l588
						}
						position++
					}
				l593:
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l596
						}
						position++
						goto l595
					l596:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
						if buffer[position] != rune('T') {
							goto l588
						}
						position++
					}
				l595:
					if !_rules[rulesp]() {
						goto l588
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l588
					}
					if !_rules[rulesp]() {
						goto l588
					}
				l597:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l598
						}
						position++
						if !_rules[rulesp]() {
							goto l598
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l598
						}
						goto l597
					l598:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
					}
					depth--
					add(rulePegText, position590)
				}
				if !_rules[ruleAction32]() {
					goto l588
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position589)
			}
			return true
		l588:
			position, tokenIndex, depth = position588, tokenIndex588, depth588
			return false
		},
		/* 41 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action33)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l599
				}
				if buffer[position] != rune('=') {
					goto l599
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l599
				}
				if !_rules[ruleAction33]() {
					goto l599
				}
				depth--
				add(ruleSourceSinkParam, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 42 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position601, tokenIndex601, depth601 := position, tokenIndex, depth
			{
				position602 := position
				depth++
				{
					position603, tokenIndex603, depth603 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l604
					}
					goto l603
				l604:
					position, tokenIndex, depth = position603, tokenIndex603, depth603
					if !_rules[ruleLiteral]() {
						goto l601
					}
				}
			l603:
				depth--
				add(ruleSourceSinkParamVal, position602)
			}
			return true
		l601:
			position, tokenIndex, depth = position601, tokenIndex601, depth601
			return false
		},
		/* 43 PausedOpt <- <(<(Paused / Unpaused)?> Action34)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				{
					position607 := position
					depth++
					{
						position608, tokenIndex608, depth608 := position, tokenIndex, depth
						{
							position610, tokenIndex610, depth610 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l611
							}
							goto l610
						l611:
							position, tokenIndex, depth = position610, tokenIndex610, depth610
							if !_rules[ruleUnpaused]() {
								goto l608
							}
						}
					l610:
						goto l609
					l608:
						position, tokenIndex, depth = position608, tokenIndex608, depth608
					}
				l609:
					depth--
					add(rulePegText, position607)
				}
				if !_rules[ruleAction34]() {
					goto l605
				}
				depth--
				add(rulePausedOpt, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 44 Expression <- <orExpr> */
		func() bool {
			position612, tokenIndex612, depth612 := position, tokenIndex, depth
			{
				position613 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l612
				}
				depth--
				add(ruleExpression, position613)
			}
			return true
		l612:
			position, tokenIndex, depth = position612, tokenIndex612, depth612
			return false
		},
		/* 45 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action35)> */
		func() bool {
			position614, tokenIndex614, depth614 := position, tokenIndex, depth
			{
				position615 := position
				depth++
				{
					position616 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l614
					}
					if !_rules[rulesp]() {
						goto l614
					}
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l617
						}
						if !_rules[rulesp]() {
							goto l617
						}
						if !_rules[ruleandExpr]() {
							goto l617
						}
						goto l618
					l617:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
					}
				l618:
					depth--
					add(rulePegText, position616)
				}
				if !_rules[ruleAction35]() {
					goto l614
				}
				depth--
				add(ruleorExpr, position615)
			}
			return true
		l614:
			position, tokenIndex, depth = position614, tokenIndex614, depth614
			return false
		},
		/* 46 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action36)> */
		func() bool {
			position619, tokenIndex619, depth619 := position, tokenIndex, depth
			{
				position620 := position
				depth++
				{
					position621 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l619
					}
					if !_rules[rulesp]() {
						goto l619
					}
					{
						position622, tokenIndex622, depth622 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l622
						}
						if !_rules[rulesp]() {
							goto l622
						}
						if !_rules[rulenotExpr]() {
							goto l622
						}
						goto l623
					l622:
						position, tokenIndex, depth = position622, tokenIndex622, depth622
					}
				l623:
					depth--
					add(rulePegText, position621)
				}
				if !_rules[ruleAction36]() {
					goto l619
				}
				depth--
				add(ruleandExpr, position620)
			}
			return true
		l619:
			position, tokenIndex, depth = position619, tokenIndex619, depth619
			return false
		},
		/* 47 notExpr <- <(<((Not sp)? comparisonExpr)> Action37)> */
		func() bool {
			position624, tokenIndex624, depth624 := position, tokenIndex, depth
			{
				position625 := position
				depth++
				{
					position626 := position
					depth++
					{
						position627, tokenIndex627, depth627 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l627
						}
						if !_rules[rulesp]() {
							goto l627
						}
						goto l628
					l627:
						position, tokenIndex, depth = position627, tokenIndex627, depth627
					}
				l628:
					if !_rules[rulecomparisonExpr]() {
						goto l624
					}
					depth--
					add(rulePegText, position626)
				}
				if !_rules[ruleAction37]() {
					goto l624
				}
				depth--
				add(rulenotExpr, position625)
			}
			return true
		l624:
			position, tokenIndex, depth = position624, tokenIndex624, depth624
			return false
		},
		/* 48 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action38)> */
		func() bool {
			position629, tokenIndex629, depth629 := position, tokenIndex, depth
			{
				position630 := position
				depth++
				{
					position631 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l629
					}
					if !_rules[rulesp]() {
						goto l629
					}
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l632
						}
						if !_rules[rulesp]() {
							goto l632
						}
						if !_rules[ruleisExpr]() {
							goto l632
						}
						goto l633
					l632:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
					}
				l633:
					depth--
					add(rulePegText, position631)
				}
				if !_rules[ruleAction38]() {
					goto l629
				}
				depth--
				add(rulecomparisonExpr, position630)
			}
			return true
		l629:
			position, tokenIndex, depth = position629, tokenIndex629, depth629
			return false
		},
		/* 49 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action39)> */
		func() bool {
			position634, tokenIndex634, depth634 := position, tokenIndex, depth
			{
				position635 := position
				depth++
				{
					position636 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l634
					}
					if !_rules[rulesp]() {
						goto l634
					}
					{
						position637, tokenIndex637, depth637 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l637
						}
						if !_rules[rulesp]() {
							goto l637
						}
						if !_rules[ruleNullLiteral]() {
							goto l637
						}
						goto l638
					l637:
						position, tokenIndex, depth = position637, tokenIndex637, depth637
					}
				l638:
					depth--
					add(rulePegText, position636)
				}
				if !_rules[ruleAction39]() {
					goto l634
				}
				depth--
				add(ruleisExpr, position635)
			}
			return true
		l634:
			position, tokenIndex, depth = position634, tokenIndex634, depth634
			return false
		},
		/* 50 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action40)> */
		func() bool {
			position639, tokenIndex639, depth639 := position, tokenIndex, depth
			{
				position640 := position
				depth++
				{
					position641 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l639
					}
					if !_rules[rulesp]() {
						goto l639
					}
					{
						position642, tokenIndex642, depth642 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l642
						}
						if !_rules[rulesp]() {
							goto l642
						}
						if !_rules[ruleproductExpr]() {
							goto l642
						}
						goto l643
					l642:
						position, tokenIndex, depth = position642, tokenIndex642, depth642
					}
				l643:
					depth--
					add(rulePegText, position641)
				}
				if !_rules[ruleAction40]() {
					goto l639
				}
				depth--
				add(ruletermExpr, position640)
			}
			return true
		l639:
			position, tokenIndex, depth = position639, tokenIndex639, depth639
			return false
		},
		/* 51 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action41)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				{
					position646 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l644
					}
					if !_rules[rulesp]() {
						goto l644
					}
					{
						position647, tokenIndex647, depth647 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l647
						}
						if !_rules[rulesp]() {
							goto l647
						}
						if !_rules[ruleminusExpr]() {
							goto l647
						}
						goto l648
					l647:
						position, tokenIndex, depth = position647, tokenIndex647, depth647
					}
				l648:
					depth--
					add(rulePegText, position646)
				}
				if !_rules[ruleAction41]() {
					goto l644
				}
				depth--
				add(ruleproductExpr, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 52 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action42)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l652
						}
						if !_rules[rulesp]() {
							goto l652
						}
						goto l653
					l652:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
					}
				l653:
					if !_rules[rulebaseExpr]() {
						goto l649
					}
					depth--
					add(rulePegText, position651)
				}
				if !_rules[ruleAction42]() {
					goto l649
				}
				depth--
				add(ruleminusExpr, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 53 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position654, tokenIndex654, depth654 := position, tokenIndex, depth
			{
				position655 := position
				depth++
				{
					position656, tokenIndex656, depth656 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l657
					}
					position++
					if !_rules[rulesp]() {
						goto l657
					}
					if !_rules[ruleExpression]() {
						goto l657
					}
					if !_rules[rulesp]() {
						goto l657
					}
					if buffer[position] != rune(')') {
						goto l657
					}
					position++
					goto l656
				l657:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleBooleanLiteral]() {
						goto l658
					}
					goto l656
				l658:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleNullLiteral]() {
						goto l659
					}
					goto l656
				l659:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleFuncApp]() {
						goto l660
					}
					goto l656
				l660:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleRowMeta]() {
						goto l661
					}
					goto l656
				l661:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleRowValue]() {
						goto l662
					}
					goto l656
				l662:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleLiteral]() {
						goto l654
					}
				}
			l656:
				depth--
				add(rulebaseExpr, position655)
			}
			return true
		l654:
			position, tokenIndex, depth = position654, tokenIndex654, depth654
			return false
		},
		/* 54 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action43)> */
		func() bool {
			position663, tokenIndex663, depth663 := position, tokenIndex, depth
			{
				position664 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l663
				}
				if !_rules[rulesp]() {
					goto l663
				}
				if buffer[position] != rune('(') {
					goto l663
				}
				position++
				if !_rules[rulesp]() {
					goto l663
				}
				if !_rules[ruleFuncParams]() {
					goto l663
				}
				if !_rules[rulesp]() {
					goto l663
				}
				if buffer[position] != rune(')') {
					goto l663
				}
				position++
				if !_rules[ruleAction43]() {
					goto l663
				}
				depth--
				add(ruleFuncApp, position664)
			}
			return true
		l663:
			position, tokenIndex, depth = position663, tokenIndex663, depth663
			return false
		},
		/* 55 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action44)> */
		func() bool {
			position665, tokenIndex665, depth665 := position, tokenIndex, depth
			{
				position666 := position
				depth++
				{
					position667 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l665
					}
					if !_rules[rulesp]() {
						goto l665
					}
				l668:
					{
						position669, tokenIndex669, depth669 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l669
						}
						position++
						if !_rules[rulesp]() {
							goto l669
						}
						if !_rules[ruleExpression]() {
							goto l669
						}
						goto l668
					l669:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
					}
					depth--
					add(rulePegText, position667)
				}
				if !_rules[ruleAction44]() {
					goto l665
				}
				depth--
				add(ruleFuncParams, position666)
			}
			return true
		l665:
			position, tokenIndex, depth = position665, tokenIndex665, depth665
			return false
		},
		/* 56 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position670, tokenIndex670, depth670 := position, tokenIndex, depth
			{
				position671 := position
				depth++
				{
					position672, tokenIndex672, depth672 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l673
					}
					goto l672
				l673:
					position, tokenIndex, depth = position672, tokenIndex672, depth672
					if !_rules[ruleNumericLiteral]() {
						goto l674
					}
					goto l672
				l674:
					position, tokenIndex, depth = position672, tokenIndex672, depth672
					if !_rules[ruleStringLiteral]() {
						goto l670
					}
				}
			l672:
				depth--
				add(ruleLiteral, position671)
			}
			return true
		l670:
			position, tokenIndex, depth = position670, tokenIndex670, depth670
			return false
		},
		/* 57 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position675, tokenIndex675, depth675 := position, tokenIndex, depth
			{
				position676 := position
				depth++
				{
					position677, tokenIndex677, depth677 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l678
					}
					goto l677
				l678:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if !_rules[ruleNotEqual]() {
						goto l679
					}
					goto l677
				l679:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if !_rules[ruleLessOrEqual]() {
						goto l680
					}
					goto l677
				l680:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if !_rules[ruleLess]() {
						goto l681
					}
					goto l677
				l681:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if !_rules[ruleGreaterOrEqual]() {
						goto l682
					}
					goto l677
				l682:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if !_rules[ruleGreater]() {
						goto l683
					}
					goto l677
				l683:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if !_rules[ruleNotEqual]() {
						goto l675
					}
				}
			l677:
				depth--
				add(ruleComparisonOp, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 58 IsOp <- <(IsNot / Is)> */
		func() bool {
			position684, tokenIndex684, depth684 := position, tokenIndex, depth
			{
				position685 := position
				depth++
				{
					position686, tokenIndex686, depth686 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l687
					}
					goto l686
				l687:
					position, tokenIndex, depth = position686, tokenIndex686, depth686
					if !_rules[ruleIs]() {
						goto l684
					}
				}
			l686:
				depth--
				add(ruleIsOp, position685)
			}
			return true
		l684:
			position, tokenIndex, depth = position684, tokenIndex684, depth684
			return false
		},
		/* 59 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position688, tokenIndex688, depth688 := position, tokenIndex, depth
			{
				position689 := position
				depth++
				{
					position690, tokenIndex690, depth690 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l691
					}
					goto l690
				l691:
					position, tokenIndex, depth = position690, tokenIndex690, depth690
					if !_rules[ruleMinus]() {
						goto l688
					}
				}
			l690:
				depth--
				add(rulePlusMinusOp, position689)
			}
			return true
		l688:
			position, tokenIndex, depth = position688, tokenIndex688, depth688
			return false
		},
		/* 60 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position692, tokenIndex692, depth692 := position, tokenIndex, depth
			{
				position693 := position
				depth++
				{
					position694, tokenIndex694, depth694 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l695
					}
					goto l694
				l695:
					position, tokenIndex, depth = position694, tokenIndex694, depth694
					if !_rules[ruleDivide]() {
						goto l696
					}
					goto l694
				l696:
					position, tokenIndex, depth = position694, tokenIndex694, depth694
					if !_rules[ruleModulo]() {
						goto l692
					}
				}
			l694:
				depth--
				add(ruleMultDivOp, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 61 Stream <- <(<ident> Action45)> */
		func() bool {
			position697, tokenIndex697, depth697 := position, tokenIndex, depth
			{
				position698 := position
				depth++
				{
					position699 := position
					depth++
					if !_rules[ruleident]() {
						goto l697
					}
					depth--
					add(rulePegText, position699)
				}
				if !_rules[ruleAction45]() {
					goto l697
				}
				depth--
				add(ruleStream, position698)
			}
			return true
		l697:
			position, tokenIndex, depth = position697, tokenIndex697, depth697
			return false
		},
		/* 62 RowMeta <- <RowTimestamp> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l700
				}
				depth--
				add(ruleRowMeta, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 63 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action46)> */
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
						if !_rules[ruleident]() {
							goto l705
						}
						if buffer[position] != rune(':') {
							goto l705
						}
						position++
						goto l706
					l705:
						position, tokenIndex, depth = position705, tokenIndex705, depth705
					}
				l706:
					if buffer[position] != rune('t') {
						goto l702
					}
					position++
					if buffer[position] != rune('s') {
						goto l702
					}
					position++
					if buffer[position] != rune('(') {
						goto l702
					}
					position++
					if buffer[position] != rune(')') {
						goto l702
					}
					position++
					depth--
					add(rulePegText, position704)
				}
				if !_rules[ruleAction46]() {
					goto l702
				}
				depth--
				add(ruleRowTimestamp, position703)
			}
			return true
		l702:
			position, tokenIndex, depth = position702, tokenIndex702, depth702
			return false
		},
		/* 64 RowValue <- <(<((ident ':')? jsonPath)> Action47)> */
		func() bool {
			position707, tokenIndex707, depth707 := position, tokenIndex, depth
			{
				position708 := position
				depth++
				{
					position709 := position
					depth++
					{
						position710, tokenIndex710, depth710 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l710
						}
						if buffer[position] != rune(':') {
							goto l710
						}
						position++
						goto l711
					l710:
						position, tokenIndex, depth = position710, tokenIndex710, depth710
					}
				l711:
					if !_rules[rulejsonPath]() {
						goto l707
					}
					depth--
					add(rulePegText, position709)
				}
				if !_rules[ruleAction47]() {
					goto l707
				}
				depth--
				add(ruleRowValue, position708)
			}
			return true
		l707:
			position, tokenIndex, depth = position707, tokenIndex707, depth707
			return false
		},
		/* 65 NumericLiteral <- <(<('-'? [0-9]+)> Action48)> */
		func() bool {
			position712, tokenIndex712, depth712 := position, tokenIndex, depth
			{
				position713 := position
				depth++
				{
					position714 := position
					depth++
					{
						position715, tokenIndex715, depth715 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l715
						}
						position++
						goto l716
					l715:
						position, tokenIndex, depth = position715, tokenIndex715, depth715
					}
				l716:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l712
					}
					position++
				l717:
					{
						position718, tokenIndex718, depth718 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l718
						}
						position++
						goto l717
					l718:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
					}
					depth--
					add(rulePegText, position714)
				}
				if !_rules[ruleAction48]() {
					goto l712
				}
				depth--
				add(ruleNumericLiteral, position713)
			}
			return true
		l712:
			position, tokenIndex, depth = position712, tokenIndex712, depth712
			return false
		},
		/* 66 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action49)> */
		func() bool {
			position719, tokenIndex719, depth719 := position, tokenIndex, depth
			{
				position720 := position
				depth++
				{
					position721 := position
					depth++
					{
						position722, tokenIndex722, depth722 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l722
						}
						position++
						goto l723
					l722:
						position, tokenIndex, depth = position722, tokenIndex722, depth722
					}
				l723:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l719
					}
					position++
				l724:
					{
						position725, tokenIndex725, depth725 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l725
						}
						position++
						goto l724
					l725:
						position, tokenIndex, depth = position725, tokenIndex725, depth725
					}
					if buffer[position] != rune('.') {
						goto l719
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l719
					}
					position++
				l726:
					{
						position727, tokenIndex727, depth727 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l727
						}
						position++
						goto l726
					l727:
						position, tokenIndex, depth = position727, tokenIndex727, depth727
					}
					depth--
					add(rulePegText, position721)
				}
				if !_rules[ruleAction49]() {
					goto l719
				}
				depth--
				add(ruleFloatLiteral, position720)
			}
			return true
		l719:
			position, tokenIndex, depth = position719, tokenIndex719, depth719
			return false
		},
		/* 67 Function <- <(<ident> Action50)> */
		func() bool {
			position728, tokenIndex728, depth728 := position, tokenIndex, depth
			{
				position729 := position
				depth++
				{
					position730 := position
					depth++
					if !_rules[ruleident]() {
						goto l728
					}
					depth--
					add(rulePegText, position730)
				}
				if !_rules[ruleAction50]() {
					goto l728
				}
				depth--
				add(ruleFunction, position729)
			}
			return true
		l728:
			position, tokenIndex, depth = position728, tokenIndex728, depth728
			return false
		},
		/* 68 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action51)> */
		func() bool {
			position731, tokenIndex731, depth731 := position, tokenIndex, depth
			{
				position732 := position
				depth++
				{
					position733 := position
					depth++
					{
						position734, tokenIndex734, depth734 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l735
						}
						position++
						goto l734
					l735:
						position, tokenIndex, depth = position734, tokenIndex734, depth734
						if buffer[position] != rune('N') {
							goto l731
						}
						position++
					}
				l734:
					{
						position736, tokenIndex736, depth736 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l737
						}
						position++
						goto l736
					l737:
						position, tokenIndex, depth = position736, tokenIndex736, depth736
						if buffer[position] != rune('U') {
							goto l731
						}
						position++
					}
				l736:
					{
						position738, tokenIndex738, depth738 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l739
						}
						position++
						goto l738
					l739:
						position, tokenIndex, depth = position738, tokenIndex738, depth738
						if buffer[position] != rune('L') {
							goto l731
						}
						position++
					}
				l738:
					{
						position740, tokenIndex740, depth740 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l741
						}
						position++
						goto l740
					l741:
						position, tokenIndex, depth = position740, tokenIndex740, depth740
						if buffer[position] != rune('L') {
							goto l731
						}
						position++
					}
				l740:
					depth--
					add(rulePegText, position733)
				}
				if !_rules[ruleAction51]() {
					goto l731
				}
				depth--
				add(ruleNullLiteral, position732)
			}
			return true
		l731:
			position, tokenIndex, depth = position731, tokenIndex731, depth731
			return false
		},
		/* 69 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position742, tokenIndex742, depth742 := position, tokenIndex, depth
			{
				position743 := position
				depth++
				{
					position744, tokenIndex744, depth744 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l745
					}
					goto l744
				l745:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
					if !_rules[ruleFALSE]() {
						goto l742
					}
				}
			l744:
				depth--
				add(ruleBooleanLiteral, position743)
			}
			return true
		l742:
			position, tokenIndex, depth = position742, tokenIndex742, depth742
			return false
		},
		/* 70 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action52)> */
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
						if buffer[position] != rune('t') {
							goto l750
						}
						position++
						goto l749
					l750:
						position, tokenIndex, depth = position749, tokenIndex749, depth749
						if buffer[position] != rune('T') {
							goto l746
						}
						position++
					}
				l749:
					{
						position751, tokenIndex751, depth751 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l752
						}
						position++
						goto l751
					l752:
						position, tokenIndex, depth = position751, tokenIndex751, depth751
						if buffer[position] != rune('R') {
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
						if buffer[position] != rune('e') {
							goto l756
						}
						position++
						goto l755
					l756:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						if buffer[position] != rune('E') {
							goto l746
						}
						position++
					}
				l755:
					depth--
					add(rulePegText, position748)
				}
				if !_rules[ruleAction52]() {
					goto l746
				}
				depth--
				add(ruleTRUE, position747)
			}
			return true
		l746:
			position, tokenIndex, depth = position746, tokenIndex746, depth746
			return false
		},
		/* 71 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action53)> */
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
						if buffer[position] != rune('f') {
							goto l761
						}
						position++
						goto l760
					l761:
						position, tokenIndex, depth = position760, tokenIndex760, depth760
						if buffer[position] != rune('F') {
							goto l757
						}
						position++
					}
				l760:
					{
						position762, tokenIndex762, depth762 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l763
						}
						position++
						goto l762
					l763:
						position, tokenIndex, depth = position762, tokenIndex762, depth762
						if buffer[position] != rune('A') {
							goto l757
						}
						position++
					}
				l762:
					{
						position764, tokenIndex764, depth764 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l765
						}
						position++
						goto l764
					l765:
						position, tokenIndex, depth = position764, tokenIndex764, depth764
						if buffer[position] != rune('L') {
							goto l757
						}
						position++
					}
				l764:
					{
						position766, tokenIndex766, depth766 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l767
						}
						position++
						goto l766
					l767:
						position, tokenIndex, depth = position766, tokenIndex766, depth766
						if buffer[position] != rune('S') {
							goto l757
						}
						position++
					}
				l766:
					{
						position768, tokenIndex768, depth768 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l769
						}
						position++
						goto l768
					l769:
						position, tokenIndex, depth = position768, tokenIndex768, depth768
						if buffer[position] != rune('E') {
							goto l757
						}
						position++
					}
				l768:
					depth--
					add(rulePegText, position759)
				}
				if !_rules[ruleAction53]() {
					goto l757
				}
				depth--
				add(ruleFALSE, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 72 Wildcard <- <(<'*'> Action54)> */
		func() bool {
			position770, tokenIndex770, depth770 := position, tokenIndex, depth
			{
				position771 := position
				depth++
				{
					position772 := position
					depth++
					if buffer[position] != rune('*') {
						goto l770
					}
					position++
					depth--
					add(rulePegText, position772)
				}
				if !_rules[ruleAction54]() {
					goto l770
				}
				depth--
				add(ruleWildcard, position771)
			}
			return true
		l770:
			position, tokenIndex, depth = position770, tokenIndex770, depth770
			return false
		},
		/* 73 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action55)> */
		func() bool {
			position773, tokenIndex773, depth773 := position, tokenIndex, depth
			{
				position774 := position
				depth++
				{
					position775 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l773
					}
					position++
				l776:
					{
						position777, tokenIndex777, depth777 := position, tokenIndex, depth
						{
							position778, tokenIndex778, depth778 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l779
							}
							position++
							if buffer[position] != rune('\'') {
								goto l779
							}
							position++
							goto l778
						l779:
							position, tokenIndex, depth = position778, tokenIndex778, depth778
							{
								position780, tokenIndex780, depth780 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l780
								}
								position++
								goto l777
							l780:
								position, tokenIndex, depth = position780, tokenIndex780, depth780
							}
							if !matchDot() {
								goto l777
							}
						}
					l778:
						goto l776
					l777:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
					}
					if buffer[position] != rune('\'') {
						goto l773
					}
					position++
					depth--
					add(rulePegText, position775)
				}
				if !_rules[ruleAction55]() {
					goto l773
				}
				depth--
				add(ruleStringLiteral, position774)
			}
			return true
		l773:
			position, tokenIndex, depth = position773, tokenIndex773, depth773
			return false
		},
		/* 74 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action56)> */
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
						if buffer[position] != rune('i') {
							goto l785
						}
						position++
						goto l784
					l785:
						position, tokenIndex, depth = position784, tokenIndex784, depth784
						if buffer[position] != rune('I') {
							goto l781
						}
						position++
					}
				l784:
					{
						position786, tokenIndex786, depth786 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l787
						}
						position++
						goto l786
					l787:
						position, tokenIndex, depth = position786, tokenIndex786, depth786
						if buffer[position] != rune('S') {
							goto l781
						}
						position++
					}
				l786:
					{
						position788, tokenIndex788, depth788 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l789
						}
						position++
						goto l788
					l789:
						position, tokenIndex, depth = position788, tokenIndex788, depth788
						if buffer[position] != rune('T') {
							goto l781
						}
						position++
					}
				l788:
					{
						position790, tokenIndex790, depth790 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l791
						}
						position++
						goto l790
					l791:
						position, tokenIndex, depth = position790, tokenIndex790, depth790
						if buffer[position] != rune('R') {
							goto l781
						}
						position++
					}
				l790:
					{
						position792, tokenIndex792, depth792 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l793
						}
						position++
						goto l792
					l793:
						position, tokenIndex, depth = position792, tokenIndex792, depth792
						if buffer[position] != rune('E') {
							goto l781
						}
						position++
					}
				l792:
					{
						position794, tokenIndex794, depth794 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l795
						}
						position++
						goto l794
					l795:
						position, tokenIndex, depth = position794, tokenIndex794, depth794
						if buffer[position] != rune('A') {
							goto l781
						}
						position++
					}
				l794:
					{
						position796, tokenIndex796, depth796 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l797
						}
						position++
						goto l796
					l797:
						position, tokenIndex, depth = position796, tokenIndex796, depth796
						if buffer[position] != rune('M') {
							goto l781
						}
						position++
					}
				l796:
					depth--
					add(rulePegText, position783)
				}
				if !_rules[ruleAction56]() {
					goto l781
				}
				depth--
				add(ruleISTREAM, position782)
			}
			return true
		l781:
			position, tokenIndex, depth = position781, tokenIndex781, depth781
			return false
		},
		/* 75 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action57)> */
		func() bool {
			position798, tokenIndex798, depth798 := position, tokenIndex, depth
			{
				position799 := position
				depth++
				{
					position800 := position
					depth++
					{
						position801, tokenIndex801, depth801 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l802
						}
						position++
						goto l801
					l802:
						position, tokenIndex, depth = position801, tokenIndex801, depth801
						if buffer[position] != rune('D') {
							goto l798
						}
						position++
					}
				l801:
					{
						position803, tokenIndex803, depth803 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l804
						}
						position++
						goto l803
					l804:
						position, tokenIndex, depth = position803, tokenIndex803, depth803
						if buffer[position] != rune('S') {
							goto l798
						}
						position++
					}
				l803:
					{
						position805, tokenIndex805, depth805 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l806
						}
						position++
						goto l805
					l806:
						position, tokenIndex, depth = position805, tokenIndex805, depth805
						if buffer[position] != rune('T') {
							goto l798
						}
						position++
					}
				l805:
					{
						position807, tokenIndex807, depth807 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l808
						}
						position++
						goto l807
					l808:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						if buffer[position] != rune('R') {
							goto l798
						}
						position++
					}
				l807:
					{
						position809, tokenIndex809, depth809 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l810
						}
						position++
						goto l809
					l810:
						position, tokenIndex, depth = position809, tokenIndex809, depth809
						if buffer[position] != rune('E') {
							goto l798
						}
						position++
					}
				l809:
					{
						position811, tokenIndex811, depth811 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l812
						}
						position++
						goto l811
					l812:
						position, tokenIndex, depth = position811, tokenIndex811, depth811
						if buffer[position] != rune('A') {
							goto l798
						}
						position++
					}
				l811:
					{
						position813, tokenIndex813, depth813 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l814
						}
						position++
						goto l813
					l814:
						position, tokenIndex, depth = position813, tokenIndex813, depth813
						if buffer[position] != rune('M') {
							goto l798
						}
						position++
					}
				l813:
					depth--
					add(rulePegText, position800)
				}
				if !_rules[ruleAction57]() {
					goto l798
				}
				depth--
				add(ruleDSTREAM, position799)
			}
			return true
		l798:
			position, tokenIndex, depth = position798, tokenIndex798, depth798
			return false
		},
		/* 76 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action58)> */
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
						if buffer[position] != rune('r') {
							goto l819
						}
						position++
						goto l818
					l819:
						position, tokenIndex, depth = position818, tokenIndex818, depth818
						if buffer[position] != rune('R') {
							goto l815
						}
						position++
					}
				l818:
					{
						position820, tokenIndex820, depth820 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l821
						}
						position++
						goto l820
					l821:
						position, tokenIndex, depth = position820, tokenIndex820, depth820
						if buffer[position] != rune('S') {
							goto l815
						}
						position++
					}
				l820:
					{
						position822, tokenIndex822, depth822 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l823
						}
						position++
						goto l822
					l823:
						position, tokenIndex, depth = position822, tokenIndex822, depth822
						if buffer[position] != rune('T') {
							goto l815
						}
						position++
					}
				l822:
					{
						position824, tokenIndex824, depth824 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l825
						}
						position++
						goto l824
					l825:
						position, tokenIndex, depth = position824, tokenIndex824, depth824
						if buffer[position] != rune('R') {
							goto l815
						}
						position++
					}
				l824:
					{
						position826, tokenIndex826, depth826 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l827
						}
						position++
						goto l826
					l827:
						position, tokenIndex, depth = position826, tokenIndex826, depth826
						if buffer[position] != rune('E') {
							goto l815
						}
						position++
					}
				l826:
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
							goto l815
						}
						position++
					}
				l828:
					{
						position830, tokenIndex830, depth830 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l831
						}
						position++
						goto l830
					l831:
						position, tokenIndex, depth = position830, tokenIndex830, depth830
						if buffer[position] != rune('M') {
							goto l815
						}
						position++
					}
				l830:
					depth--
					add(rulePegText, position817)
				}
				if !_rules[ruleAction58]() {
					goto l815
				}
				depth--
				add(ruleRSTREAM, position816)
			}
			return true
		l815:
			position, tokenIndex, depth = position815, tokenIndex815, depth815
			return false
		},
		/* 77 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action59)> */
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
						if buffer[position] != rune('t') {
							goto l836
						}
						position++
						goto l835
					l836:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
						if buffer[position] != rune('T') {
							goto l832
						}
						position++
					}
				l835:
					{
						position837, tokenIndex837, depth837 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l838
						}
						position++
						goto l837
					l838:
						position, tokenIndex, depth = position837, tokenIndex837, depth837
						if buffer[position] != rune('U') {
							goto l832
						}
						position++
					}
				l837:
					{
						position839, tokenIndex839, depth839 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l840
						}
						position++
						goto l839
					l840:
						position, tokenIndex, depth = position839, tokenIndex839, depth839
						if buffer[position] != rune('P') {
							goto l832
						}
						position++
					}
				l839:
					{
						position841, tokenIndex841, depth841 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l842
						}
						position++
						goto l841
					l842:
						position, tokenIndex, depth = position841, tokenIndex841, depth841
						if buffer[position] != rune('L') {
							goto l832
						}
						position++
					}
				l841:
					{
						position843, tokenIndex843, depth843 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l844
						}
						position++
						goto l843
					l844:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						if buffer[position] != rune('E') {
							goto l832
						}
						position++
					}
				l843:
					{
						position845, tokenIndex845, depth845 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l846
						}
						position++
						goto l845
					l846:
						position, tokenIndex, depth = position845, tokenIndex845, depth845
						if buffer[position] != rune('S') {
							goto l832
						}
						position++
					}
				l845:
					depth--
					add(rulePegText, position834)
				}
				if !_rules[ruleAction59]() {
					goto l832
				}
				depth--
				add(ruleTUPLES, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 78 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action60)> */
		func() bool {
			position847, tokenIndex847, depth847 := position, tokenIndex, depth
			{
				position848 := position
				depth++
				{
					position849 := position
					depth++
					{
						position850, tokenIndex850, depth850 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l851
						}
						position++
						goto l850
					l851:
						position, tokenIndex, depth = position850, tokenIndex850, depth850
						if buffer[position] != rune('S') {
							goto l847
						}
						position++
					}
				l850:
					{
						position852, tokenIndex852, depth852 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l853
						}
						position++
						goto l852
					l853:
						position, tokenIndex, depth = position852, tokenIndex852, depth852
						if buffer[position] != rune('E') {
							goto l847
						}
						position++
					}
				l852:
					{
						position854, tokenIndex854, depth854 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l855
						}
						position++
						goto l854
					l855:
						position, tokenIndex, depth = position854, tokenIndex854, depth854
						if buffer[position] != rune('C') {
							goto l847
						}
						position++
					}
				l854:
					{
						position856, tokenIndex856, depth856 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l857
						}
						position++
						goto l856
					l857:
						position, tokenIndex, depth = position856, tokenIndex856, depth856
						if buffer[position] != rune('O') {
							goto l847
						}
						position++
					}
				l856:
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l859
						}
						position++
						goto l858
					l859:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if buffer[position] != rune('N') {
							goto l847
						}
						position++
					}
				l858:
					{
						position860, tokenIndex860, depth860 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l861
						}
						position++
						goto l860
					l861:
						position, tokenIndex, depth = position860, tokenIndex860, depth860
						if buffer[position] != rune('D') {
							goto l847
						}
						position++
					}
				l860:
					{
						position862, tokenIndex862, depth862 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l863
						}
						position++
						goto l862
					l863:
						position, tokenIndex, depth = position862, tokenIndex862, depth862
						if buffer[position] != rune('S') {
							goto l847
						}
						position++
					}
				l862:
					depth--
					add(rulePegText, position849)
				}
				if !_rules[ruleAction60]() {
					goto l847
				}
				depth--
				add(ruleSECONDS, position848)
			}
			return true
		l847:
			position, tokenIndex, depth = position847, tokenIndex847, depth847
			return false
		},
		/* 79 StreamIdentifier <- <(<ident> Action61)> */
		func() bool {
			position864, tokenIndex864, depth864 := position, tokenIndex, depth
			{
				position865 := position
				depth++
				{
					position866 := position
					depth++
					if !_rules[ruleident]() {
						goto l864
					}
					depth--
					add(rulePegText, position866)
				}
				if !_rules[ruleAction61]() {
					goto l864
				}
				depth--
				add(ruleStreamIdentifier, position865)
			}
			return true
		l864:
			position, tokenIndex, depth = position864, tokenIndex864, depth864
			return false
		},
		/* 80 SourceSinkType <- <(<ident> Action62)> */
		func() bool {
			position867, tokenIndex867, depth867 := position, tokenIndex, depth
			{
				position868 := position
				depth++
				{
					position869 := position
					depth++
					if !_rules[ruleident]() {
						goto l867
					}
					depth--
					add(rulePegText, position869)
				}
				if !_rules[ruleAction62]() {
					goto l867
				}
				depth--
				add(ruleSourceSinkType, position868)
			}
			return true
		l867:
			position, tokenIndex, depth = position867, tokenIndex867, depth867
			return false
		},
		/* 81 SourceSinkParamKey <- <(<ident> Action63)> */
		func() bool {
			position870, tokenIndex870, depth870 := position, tokenIndex, depth
			{
				position871 := position
				depth++
				{
					position872 := position
					depth++
					if !_rules[ruleident]() {
						goto l870
					}
					depth--
					add(rulePegText, position872)
				}
				if !_rules[ruleAction63]() {
					goto l870
				}
				depth--
				add(ruleSourceSinkParamKey, position871)
			}
			return true
		l870:
			position, tokenIndex, depth = position870, tokenIndex870, depth870
			return false
		},
		/* 82 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action64)> */
		func() bool {
			position873, tokenIndex873, depth873 := position, tokenIndex, depth
			{
				position874 := position
				depth++
				{
					position875 := position
					depth++
					{
						position876, tokenIndex876, depth876 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l877
						}
						position++
						goto l876
					l877:
						position, tokenIndex, depth = position876, tokenIndex876, depth876
						if buffer[position] != rune('P') {
							goto l873
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
							goto l873
						}
						position++
					}
				l878:
					{
						position880, tokenIndex880, depth880 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l881
						}
						position++
						goto l880
					l881:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
						if buffer[position] != rune('U') {
							goto l873
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
							goto l873
						}
						position++
					}
				l882:
					{
						position884, tokenIndex884, depth884 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l885
						}
						position++
						goto l884
					l885:
						position, tokenIndex, depth = position884, tokenIndex884, depth884
						if buffer[position] != rune('E') {
							goto l873
						}
						position++
					}
				l884:
					{
						position886, tokenIndex886, depth886 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l887
						}
						position++
						goto l886
					l887:
						position, tokenIndex, depth = position886, tokenIndex886, depth886
						if buffer[position] != rune('D') {
							goto l873
						}
						position++
					}
				l886:
					depth--
					add(rulePegText, position875)
				}
				if !_rules[ruleAction64]() {
					goto l873
				}
				depth--
				add(rulePaused, position874)
			}
			return true
		l873:
			position, tokenIndex, depth = position873, tokenIndex873, depth873
			return false
		},
		/* 83 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action65)> */
		func() bool {
			position888, tokenIndex888, depth888 := position, tokenIndex, depth
			{
				position889 := position
				depth++
				{
					position890 := position
					depth++
					{
						position891, tokenIndex891, depth891 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l892
						}
						position++
						goto l891
					l892:
						position, tokenIndex, depth = position891, tokenIndex891, depth891
						if buffer[position] != rune('U') {
							goto l888
						}
						position++
					}
				l891:
					{
						position893, tokenIndex893, depth893 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l894
						}
						position++
						goto l893
					l894:
						position, tokenIndex, depth = position893, tokenIndex893, depth893
						if buffer[position] != rune('N') {
							goto l888
						}
						position++
					}
				l893:
					{
						position895, tokenIndex895, depth895 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l896
						}
						position++
						goto l895
					l896:
						position, tokenIndex, depth = position895, tokenIndex895, depth895
						if buffer[position] != rune('P') {
							goto l888
						}
						position++
					}
				l895:
					{
						position897, tokenIndex897, depth897 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l898
						}
						position++
						goto l897
					l898:
						position, tokenIndex, depth = position897, tokenIndex897, depth897
						if buffer[position] != rune('A') {
							goto l888
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
							goto l888
						}
						position++
					}
				l899:
					{
						position901, tokenIndex901, depth901 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l902
						}
						position++
						goto l901
					l902:
						position, tokenIndex, depth = position901, tokenIndex901, depth901
						if buffer[position] != rune('S') {
							goto l888
						}
						position++
					}
				l901:
					{
						position903, tokenIndex903, depth903 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l904
						}
						position++
						goto l903
					l904:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if buffer[position] != rune('E') {
							goto l888
						}
						position++
					}
				l903:
					{
						position905, tokenIndex905, depth905 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l906
						}
						position++
						goto l905
					l906:
						position, tokenIndex, depth = position905, tokenIndex905, depth905
						if buffer[position] != rune('D') {
							goto l888
						}
						position++
					}
				l905:
					depth--
					add(rulePegText, position890)
				}
				if !_rules[ruleAction65]() {
					goto l888
				}
				depth--
				add(ruleUnpaused, position889)
			}
			return true
		l888:
			position, tokenIndex, depth = position888, tokenIndex888, depth888
			return false
		},
		/* 84 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action66)> */
		func() bool {
			position907, tokenIndex907, depth907 := position, tokenIndex, depth
			{
				position908 := position
				depth++
				{
					position909 := position
					depth++
					{
						position910, tokenIndex910, depth910 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l911
						}
						position++
						goto l910
					l911:
						position, tokenIndex, depth = position910, tokenIndex910, depth910
						if buffer[position] != rune('O') {
							goto l907
						}
						position++
					}
				l910:
					{
						position912, tokenIndex912, depth912 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l913
						}
						position++
						goto l912
					l913:
						position, tokenIndex, depth = position912, tokenIndex912, depth912
						if buffer[position] != rune('R') {
							goto l907
						}
						position++
					}
				l912:
					depth--
					add(rulePegText, position909)
				}
				if !_rules[ruleAction66]() {
					goto l907
				}
				depth--
				add(ruleOr, position908)
			}
			return true
		l907:
			position, tokenIndex, depth = position907, tokenIndex907, depth907
			return false
		},
		/* 85 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action67)> */
		func() bool {
			position914, tokenIndex914, depth914 := position, tokenIndex, depth
			{
				position915 := position
				depth++
				{
					position916 := position
					depth++
					{
						position917, tokenIndex917, depth917 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l918
						}
						position++
						goto l917
					l918:
						position, tokenIndex, depth = position917, tokenIndex917, depth917
						if buffer[position] != rune('A') {
							goto l914
						}
						position++
					}
				l917:
					{
						position919, tokenIndex919, depth919 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l920
						}
						position++
						goto l919
					l920:
						position, tokenIndex, depth = position919, tokenIndex919, depth919
						if buffer[position] != rune('N') {
							goto l914
						}
						position++
					}
				l919:
					{
						position921, tokenIndex921, depth921 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l922
						}
						position++
						goto l921
					l922:
						position, tokenIndex, depth = position921, tokenIndex921, depth921
						if buffer[position] != rune('D') {
							goto l914
						}
						position++
					}
				l921:
					depth--
					add(rulePegText, position916)
				}
				if !_rules[ruleAction67]() {
					goto l914
				}
				depth--
				add(ruleAnd, position915)
			}
			return true
		l914:
			position, tokenIndex, depth = position914, tokenIndex914, depth914
			return false
		},
		/* 86 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action68)> */
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
						if buffer[position] != rune('n') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('N') {
							goto l923
						}
						position++
					}
				l926:
					{
						position928, tokenIndex928, depth928 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l929
						}
						position++
						goto l928
					l929:
						position, tokenIndex, depth = position928, tokenIndex928, depth928
						if buffer[position] != rune('O') {
							goto l923
						}
						position++
					}
				l928:
					{
						position930, tokenIndex930, depth930 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l931
						}
						position++
						goto l930
					l931:
						position, tokenIndex, depth = position930, tokenIndex930, depth930
						if buffer[position] != rune('T') {
							goto l923
						}
						position++
					}
				l930:
					depth--
					add(rulePegText, position925)
				}
				if !_rules[ruleAction68]() {
					goto l923
				}
				depth--
				add(ruleNot, position924)
			}
			return true
		l923:
			position, tokenIndex, depth = position923, tokenIndex923, depth923
			return false
		},
		/* 87 Equal <- <(<'='> Action69)> */
		func() bool {
			position932, tokenIndex932, depth932 := position, tokenIndex, depth
			{
				position933 := position
				depth++
				{
					position934 := position
					depth++
					if buffer[position] != rune('=') {
						goto l932
					}
					position++
					depth--
					add(rulePegText, position934)
				}
				if !_rules[ruleAction69]() {
					goto l932
				}
				depth--
				add(ruleEqual, position933)
			}
			return true
		l932:
			position, tokenIndex, depth = position932, tokenIndex932, depth932
			return false
		},
		/* 88 Less <- <(<'<'> Action70)> */
		func() bool {
			position935, tokenIndex935, depth935 := position, tokenIndex, depth
			{
				position936 := position
				depth++
				{
					position937 := position
					depth++
					if buffer[position] != rune('<') {
						goto l935
					}
					position++
					depth--
					add(rulePegText, position937)
				}
				if !_rules[ruleAction70]() {
					goto l935
				}
				depth--
				add(ruleLess, position936)
			}
			return true
		l935:
			position, tokenIndex, depth = position935, tokenIndex935, depth935
			return false
		},
		/* 89 LessOrEqual <- <(<('<' '=')> Action71)> */
		func() bool {
			position938, tokenIndex938, depth938 := position, tokenIndex, depth
			{
				position939 := position
				depth++
				{
					position940 := position
					depth++
					if buffer[position] != rune('<') {
						goto l938
					}
					position++
					if buffer[position] != rune('=') {
						goto l938
					}
					position++
					depth--
					add(rulePegText, position940)
				}
				if !_rules[ruleAction71]() {
					goto l938
				}
				depth--
				add(ruleLessOrEqual, position939)
			}
			return true
		l938:
			position, tokenIndex, depth = position938, tokenIndex938, depth938
			return false
		},
		/* 90 Greater <- <(<'>'> Action72)> */
		func() bool {
			position941, tokenIndex941, depth941 := position, tokenIndex, depth
			{
				position942 := position
				depth++
				{
					position943 := position
					depth++
					if buffer[position] != rune('>') {
						goto l941
					}
					position++
					depth--
					add(rulePegText, position943)
				}
				if !_rules[ruleAction72]() {
					goto l941
				}
				depth--
				add(ruleGreater, position942)
			}
			return true
		l941:
			position, tokenIndex, depth = position941, tokenIndex941, depth941
			return false
		},
		/* 91 GreaterOrEqual <- <(<('>' '=')> Action73)> */
		func() bool {
			position944, tokenIndex944, depth944 := position, tokenIndex, depth
			{
				position945 := position
				depth++
				{
					position946 := position
					depth++
					if buffer[position] != rune('>') {
						goto l944
					}
					position++
					if buffer[position] != rune('=') {
						goto l944
					}
					position++
					depth--
					add(rulePegText, position946)
				}
				if !_rules[ruleAction73]() {
					goto l944
				}
				depth--
				add(ruleGreaterOrEqual, position945)
			}
			return true
		l944:
			position, tokenIndex, depth = position944, tokenIndex944, depth944
			return false
		},
		/* 92 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action74)> */
		func() bool {
			position947, tokenIndex947, depth947 := position, tokenIndex, depth
			{
				position948 := position
				depth++
				{
					position949 := position
					depth++
					{
						position950, tokenIndex950, depth950 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l951
						}
						position++
						if buffer[position] != rune('=') {
							goto l951
						}
						position++
						goto l950
					l951:
						position, tokenIndex, depth = position950, tokenIndex950, depth950
						if buffer[position] != rune('<') {
							goto l947
						}
						position++
						if buffer[position] != rune('>') {
							goto l947
						}
						position++
					}
				l950:
					depth--
					add(rulePegText, position949)
				}
				if !_rules[ruleAction74]() {
					goto l947
				}
				depth--
				add(ruleNotEqual, position948)
			}
			return true
		l947:
			position, tokenIndex, depth = position947, tokenIndex947, depth947
			return false
		},
		/* 93 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action75)> */
		func() bool {
			position952, tokenIndex952, depth952 := position, tokenIndex, depth
			{
				position953 := position
				depth++
				{
					position954 := position
					depth++
					{
						position955, tokenIndex955, depth955 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l956
						}
						position++
						goto l955
					l956:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('I') {
							goto l952
						}
						position++
					}
				l955:
					{
						position957, tokenIndex957, depth957 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l958
						}
						position++
						goto l957
					l958:
						position, tokenIndex, depth = position957, tokenIndex957, depth957
						if buffer[position] != rune('S') {
							goto l952
						}
						position++
					}
				l957:
					depth--
					add(rulePegText, position954)
				}
				if !_rules[ruleAction75]() {
					goto l952
				}
				depth--
				add(ruleIs, position953)
			}
			return true
		l952:
			position, tokenIndex, depth = position952, tokenIndex952, depth952
			return false
		},
		/* 94 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action76)> */
		func() bool {
			position959, tokenIndex959, depth959 := position, tokenIndex, depth
			{
				position960 := position
				depth++
				{
					position961 := position
					depth++
					{
						position962, tokenIndex962, depth962 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l963
						}
						position++
						goto l962
					l963:
						position, tokenIndex, depth = position962, tokenIndex962, depth962
						if buffer[position] != rune('I') {
							goto l959
						}
						position++
					}
				l962:
					{
						position964, tokenIndex964, depth964 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l965
						}
						position++
						goto l964
					l965:
						position, tokenIndex, depth = position964, tokenIndex964, depth964
						if buffer[position] != rune('S') {
							goto l959
						}
						position++
					}
				l964:
					if !_rules[rulesp]() {
						goto l959
					}
					{
						position966, tokenIndex966, depth966 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l967
						}
						position++
						goto l966
					l967:
						position, tokenIndex, depth = position966, tokenIndex966, depth966
						if buffer[position] != rune('N') {
							goto l959
						}
						position++
					}
				l966:
					{
						position968, tokenIndex968, depth968 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l969
						}
						position++
						goto l968
					l969:
						position, tokenIndex, depth = position968, tokenIndex968, depth968
						if buffer[position] != rune('O') {
							goto l959
						}
						position++
					}
				l968:
					{
						position970, tokenIndex970, depth970 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l971
						}
						position++
						goto l970
					l971:
						position, tokenIndex, depth = position970, tokenIndex970, depth970
						if buffer[position] != rune('T') {
							goto l959
						}
						position++
					}
				l970:
					depth--
					add(rulePegText, position961)
				}
				if !_rules[ruleAction76]() {
					goto l959
				}
				depth--
				add(ruleIsNot, position960)
			}
			return true
		l959:
			position, tokenIndex, depth = position959, tokenIndex959, depth959
			return false
		},
		/* 95 Plus <- <(<'+'> Action77)> */
		func() bool {
			position972, tokenIndex972, depth972 := position, tokenIndex, depth
			{
				position973 := position
				depth++
				{
					position974 := position
					depth++
					if buffer[position] != rune('+') {
						goto l972
					}
					position++
					depth--
					add(rulePegText, position974)
				}
				if !_rules[ruleAction77]() {
					goto l972
				}
				depth--
				add(rulePlus, position973)
			}
			return true
		l972:
			position, tokenIndex, depth = position972, tokenIndex972, depth972
			return false
		},
		/* 96 Minus <- <(<'-'> Action78)> */
		func() bool {
			position975, tokenIndex975, depth975 := position, tokenIndex, depth
			{
				position976 := position
				depth++
				{
					position977 := position
					depth++
					if buffer[position] != rune('-') {
						goto l975
					}
					position++
					depth--
					add(rulePegText, position977)
				}
				if !_rules[ruleAction78]() {
					goto l975
				}
				depth--
				add(ruleMinus, position976)
			}
			return true
		l975:
			position, tokenIndex, depth = position975, tokenIndex975, depth975
			return false
		},
		/* 97 Multiply <- <(<'*'> Action79)> */
		func() bool {
			position978, tokenIndex978, depth978 := position, tokenIndex, depth
			{
				position979 := position
				depth++
				{
					position980 := position
					depth++
					if buffer[position] != rune('*') {
						goto l978
					}
					position++
					depth--
					add(rulePegText, position980)
				}
				if !_rules[ruleAction79]() {
					goto l978
				}
				depth--
				add(ruleMultiply, position979)
			}
			return true
		l978:
			position, tokenIndex, depth = position978, tokenIndex978, depth978
			return false
		},
		/* 98 Divide <- <(<'/'> Action80)> */
		func() bool {
			position981, tokenIndex981, depth981 := position, tokenIndex, depth
			{
				position982 := position
				depth++
				{
					position983 := position
					depth++
					if buffer[position] != rune('/') {
						goto l981
					}
					position++
					depth--
					add(rulePegText, position983)
				}
				if !_rules[ruleAction80]() {
					goto l981
				}
				depth--
				add(ruleDivide, position982)
			}
			return true
		l981:
			position, tokenIndex, depth = position981, tokenIndex981, depth981
			return false
		},
		/* 99 Modulo <- <(<'%'> Action81)> */
		func() bool {
			position984, tokenIndex984, depth984 := position, tokenIndex, depth
			{
				position985 := position
				depth++
				{
					position986 := position
					depth++
					if buffer[position] != rune('%') {
						goto l984
					}
					position++
					depth--
					add(rulePegText, position986)
				}
				if !_rules[ruleAction81]() {
					goto l984
				}
				depth--
				add(ruleModulo, position985)
			}
			return true
		l984:
			position, tokenIndex, depth = position984, tokenIndex984, depth984
			return false
		},
		/* 100 UnaryMinus <- <(<'-'> Action82)> */
		func() bool {
			position987, tokenIndex987, depth987 := position, tokenIndex, depth
			{
				position988 := position
				depth++
				{
					position989 := position
					depth++
					if buffer[position] != rune('-') {
						goto l987
					}
					position++
					depth--
					add(rulePegText, position989)
				}
				if !_rules[ruleAction82]() {
					goto l987
				}
				depth--
				add(ruleUnaryMinus, position988)
			}
			return true
		l987:
			position, tokenIndex, depth = position987, tokenIndex987, depth987
			return false
		},
		/* 101 Identifier <- <(<ident> Action83)> */
		func() bool {
			position990, tokenIndex990, depth990 := position, tokenIndex, depth
			{
				position991 := position
				depth++
				{
					position992 := position
					depth++
					if !_rules[ruleident]() {
						goto l990
					}
					depth--
					add(rulePegText, position992)
				}
				if !_rules[ruleAction83]() {
					goto l990
				}
				depth--
				add(ruleIdentifier, position991)
			}
			return true
		l990:
			position, tokenIndex, depth = position990, tokenIndex990, depth990
			return false
		},
		/* 102 TargetIdentifier <- <(<jsonPath> Action84)> */
		func() bool {
			position993, tokenIndex993, depth993 := position, tokenIndex, depth
			{
				position994 := position
				depth++
				{
					position995 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l993
					}
					depth--
					add(rulePegText, position995)
				}
				if !_rules[ruleAction84]() {
					goto l993
				}
				depth--
				add(ruleTargetIdentifier, position994)
			}
			return true
		l993:
			position, tokenIndex, depth = position993, tokenIndex993, depth993
			return false
		},
		/* 103 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position996, tokenIndex996, depth996 := position, tokenIndex, depth
			{
				position997 := position
				depth++
				{
					position998, tokenIndex998, depth998 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l999
					}
					position++
					goto l998
				l999:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l996
					}
					position++
				}
			l998:
			l1000:
				{
					position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
					{
						position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1003
						}
						position++
						goto l1002
					l1003:
						position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1004
						}
						position++
						goto l1002
					l1004:
						position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1005
						}
						position++
						goto l1002
					l1005:
						position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
						if buffer[position] != rune('_') {
							goto l1001
						}
						position++
					}
				l1002:
					goto l1000
				l1001:
					position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
				}
				depth--
				add(ruleident, position997)
			}
			return true
		l996:
			position, tokenIndex, depth = position996, tokenIndex996, depth996
			return false
		},
		/* 104 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
			{
				position1007 := position
				depth++
				{
					position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1009
					}
					position++
					goto l1008
				l1009:
					position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1006
					}
					position++
				}
			l1008:
			l1010:
				{
					position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1014
						}
						position++
						goto l1012
					l1014:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1015
						}
						position++
						goto l1012
					l1015:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('_') {
							goto l1016
						}
						position++
						goto l1012
					l1016:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('.') {
							goto l1017
						}
						position++
						goto l1012
					l1017:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('[') {
							goto l1018
						}
						position++
						goto l1012
					l1018:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune(']') {
							goto l1019
						}
						position++
						goto l1012
					l1019:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('"') {
							goto l1011
						}
						position++
					}
				l1012:
					goto l1010
				l1011:
					position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
				}
				depth--
				add(rulejsonPath, position1007)
			}
			return true
		l1006:
			position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
			return false
		},
		/* 105 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1021 := position
				depth++
			l1022:
				{
					position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
					{
						position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1025
						}
						position++
						goto l1024
					l1025:
						position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
						if buffer[position] != rune('\t') {
							goto l1026
						}
						position++
						goto l1024
					l1026:
						position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
						if buffer[position] != rune('\n') {
							goto l1027
						}
						position++
						goto l1024
					l1027:
						position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
						if buffer[position] != rune('\r') {
							goto l1028
						}
						position++
						goto l1024
					l1028:
						position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
						if !_rules[rulecomment]() {
							goto l1023
						}
					}
				l1024:
					goto l1022
				l1023:
					position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
				}
				depth--
				add(rulesp, position1021)
			}
			return true
		},
		/* 106 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
			{
				position1030 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1029
				}
				position++
				if buffer[position] != rune('-') {
					goto l1029
				}
				position++
			l1031:
				{
					position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
					{
						position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
						{
							position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1035
							}
							position++
							goto l1034
						l1035:
							position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
							if buffer[position] != rune('\n') {
								goto l1033
							}
							position++
						}
					l1034:
						goto l1032
					l1033:
						position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
					}
					if !matchDot() {
						goto l1032
					}
					goto l1031
				l1032:
					position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
				}
				{
					position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1037
					}
					position++
					goto l1036
				l1037:
					position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
					if buffer[position] != rune('\n') {
						goto l1029
					}
					position++
				}
			l1036:
				depth--
				add(rulecomment, position1030)
			}
			return true
		l1029:
			position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
			return false
		},
		/* 108 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 109 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 110 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 111 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 112 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 113 Action5 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 114 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 115 Action7 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 116 Action8 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 117 Action9 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 118 Action10 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 119 Action11 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 120 Action12 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 121 Action13 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 122 Action14 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		nil,
		/* 124 Action15 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 125 Action16 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 126 Action17 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 127 Action18 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 128 Action19 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 129 Action20 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 130 Action21 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 131 Action22 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 132 Action23 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 133 Action24 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 134 Action25 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 135 Action26 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 136 Action27 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 137 Action28 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 138 Action29 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 139 Action30 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 140 Action31 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 141 Action32 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 142 Action33 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 143 Action34 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 144 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 145 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 146 Action37 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 147 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 148 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 149 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 150 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 151 Action42 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 152 Action43 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 153 Action44 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 154 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 155 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 156 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 157 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 158 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 159 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 160 Action51 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 161 Action52 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 162 Action53 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 163 Action54 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 164 Action55 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 165 Action56 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 166 Action57 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 167 Action58 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 168 Action59 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 169 Action60 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 170 Action61 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 171 Action62 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 172 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 173 Action64 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 174 Action65 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 175 Action66 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 176 Action67 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 177 Action68 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 178 Action69 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 179 Action70 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 180 Action71 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 181 Action72 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 182 Action73 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 183 Action74 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 184 Action75 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 185 Action76 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 186 Action77 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 187 Action78 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 188 Action79 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 189 Action80 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 190 Action81 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 191 Action82 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 192 Action83 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 193 Action84 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
	}
	p.rules = _rules
}

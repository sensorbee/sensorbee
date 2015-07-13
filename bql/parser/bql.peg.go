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
	ruleDropStreamStmt
	ruleDropSinkStmt
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
	ruleAction11
	ruleAction12
	rulePegText
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
	ruleAction80
	ruleAction81

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
	"DropStreamStmt",
	"DropSinkStmt",
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
	"Action11",
	"Action12",
	"PegText",
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
	"Action80",
	"Action81",

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
	rules  [188]func() bool
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

			p.AssembleDropStream()

		case ruleAction12:

			p.AssembleDropSink()

		case ruleAction13:

			p.AssembleEmitter(begin, end)

		case ruleAction14:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction15:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction16:

			p.AssembleStreamEmitInterval()

		case ruleAction17:

			p.AssembleProjections(begin, end)

		case ruleAction18:

			p.AssembleAlias()

		case ruleAction19:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction20:

			p.AssembleInterval()

		case ruleAction21:

			p.AssembleInterval()

		case ruleAction22:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction23:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction24:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction25:

			p.EnsureAliasedStreamWindow()

		case ruleAction26:

			p.AssembleAliasedStreamWindow()

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

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction38:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction39:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction40:

			p.AssembleFuncApp()

		case ruleAction41:

			p.AssembleExpressions(begin, end)

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction44:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction48:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction49:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction50:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction51:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction53:

			p.PushComponent(begin, end, Istream)

		case ruleAction54:

			p.PushComponent(begin, end, Dstream)

		case ruleAction55:

			p.PushComponent(begin, end, Rstream)

		case ruleAction56:

			p.PushComponent(begin, end, Tuples)

		case ruleAction57:

			p.PushComponent(begin, end, Seconds)

		case ruleAction58:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction59:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction60:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction61:

			p.PushComponent(begin, end, Yes)

		case ruleAction62:

			p.PushComponent(begin, end, No)

		case ruleAction63:

			p.PushComponent(begin, end, Or)

		case ruleAction64:

			p.PushComponent(begin, end, And)

		case ruleAction65:

			p.PushComponent(begin, end, Not)

		case ruleAction66:

			p.PushComponent(begin, end, Equal)

		case ruleAction67:

			p.PushComponent(begin, end, Less)

		case ruleAction68:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction69:

			p.PushComponent(begin, end, Greater)

		case ruleAction70:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction71:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction72:

			p.PushComponent(begin, end, Is)

		case ruleAction73:

			p.PushComponent(begin, end, IsNot)

		case ruleAction74:

			p.PushComponent(begin, end, Plus)

		case ruleAction75:

			p.PushComponent(begin, end, Minus)

		case ruleAction76:

			p.PushComponent(begin, end, Multiply)

		case ruleAction77:

			p.PushComponent(begin, end, Divide)

		case ruleAction78:

			p.PushComponent(begin, end, Modulo)

		case ruleAction79:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction80:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction81:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt / DropSourceStmt / DropStreamStmt / DropSinkStmt)> */
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
			position22, tokenIndex22, depth22 := position, tokenIndex, depth
			{
				position23 := position
				depth++
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l25
					}
					position++
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if buffer[position] != rune('S') {
						goto l22
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
						goto l22
					}
					position++
				}
			l26:
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l29
					}
					position++
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if buffer[position] != rune('L') {
						goto l22
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
						goto l22
					}
					position++
				}
			l30:
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l33
					}
					position++
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if buffer[position] != rune('C') {
						goto l22
					}
					position++
				}
			l32:
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l35
					}
					position++
					goto l34
				l35:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if buffer[position] != rune('T') {
						goto l22
					}
					position++
				}
			l34:
				if !_rules[rulesp]() {
					goto l22
				}
				if !_rules[ruleEmitter]() {
					goto l22
				}
				if !_rules[rulesp]() {
					goto l22
				}
				if !_rules[ruleProjections]() {
					goto l22
				}
				if !_rules[rulesp]() {
					goto l22
				}
				if !_rules[ruleWindowedFrom]() {
					goto l22
				}
				if !_rules[rulesp]() {
					goto l22
				}
				if !_rules[ruleFilter]() {
					goto l22
				}
				if !_rules[rulesp]() {
					goto l22
				}
				if !_rules[ruleGrouping]() {
					goto l22
				}
				if !_rules[rulesp]() {
					goto l22
				}
				if !_rules[ruleHaving]() {
					goto l22
				}
				if !_rules[rulesp]() {
					goto l22
				}
				if !_rules[ruleAction0]() {
					goto l22
				}
				depth--
				add(ruleSelectStmt, position23)
			}
			return true
		l22:
			position, tokenIndex, depth = position22, tokenIndex22, depth22
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position36, tokenIndex36, depth36 := position, tokenIndex, depth
			{
				position37 := position
				depth++
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l39
					}
					position++
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if buffer[position] != rune('C') {
						goto l36
					}
					position++
				}
			l38:
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('R') {
						goto l36
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
						goto l36
					}
					position++
				}
			l42:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l45
					}
					position++
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if buffer[position] != rune('A') {
						goto l36
					}
					position++
				}
			l44:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('T') {
						goto l36
					}
					position++
				}
			l46:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('E') {
						goto l36
					}
					position++
				}
			l48:
				if !_rules[rulesp]() {
					goto l36
				}
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('S') {
						goto l36
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('T') {
						goto l36
					}
					position++
				}
			l52:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('R') {
						goto l36
					}
					position++
				}
			l54:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('E') {
						goto l36
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('A') {
						goto l36
					}
					position++
				}
			l58:
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l61
					}
					position++
					goto l60
				l61:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
					if buffer[position] != rune('M') {
						goto l36
					}
					position++
				}
			l60:
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l36
				}
				if !_rules[rulesp]() {
					goto l36
				}
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
						goto l36
					}
					position++
				}
			l62:
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
						goto l36
					}
					position++
				}
			l64:
				if !_rules[rulesp]() {
					goto l36
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
						goto l36
					}
					position++
				}
			l66:
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
					if buffer[position] != rune('E') {
						goto l36
					}
					position++
				}
			l68:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l71
					}
					position++
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if buffer[position] != rune('L') {
						goto l36
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
						goto l36
					}
					position++
				}
			l72:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if buffer[position] != rune('C') {
						goto l36
					}
					position++
				}
			l74:
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('T') {
						goto l36
					}
					position++
				}
			l76:
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleEmitter]() {
					goto l36
				}
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleProjections]() {
					goto l36
				}
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleWindowedFrom]() {
					goto l36
				}
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleFilter]() {
					goto l36
				}
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleGrouping]() {
					goto l36
				}
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleHaving]() {
					goto l36
				}
				if !_rules[rulesp]() {
					goto l36
				}
				if !_rules[ruleAction1]() {
					goto l36
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position37)
			}
			return true
		l36:
			position, tokenIndex, depth = position36, tokenIndex36, depth36
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position78, tokenIndex78, depth78 := position, tokenIndex, depth
			{
				position79 := position
				depth++
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('C') {
						goto l78
					}
					position++
				}
			l80:
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l83
					}
					position++
					goto l82
				l83:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
					if buffer[position] != rune('R') {
						goto l78
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
						goto l78
					}
					position++
				}
			l84:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if buffer[position] != rune('A') {
						goto l78
					}
					position++
				}
			l86:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l89
					}
					position++
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if buffer[position] != rune('T') {
						goto l78
					}
					position++
				}
			l88:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('E') {
						goto l78
					}
					position++
				}
			l90:
				if !_rules[rulesp]() {
					goto l78
				}
				if !_rules[rulePausedOpt]() {
					goto l78
				}
				if !_rules[rulesp]() {
					goto l78
				}
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('S') {
						goto l78
					}
					position++
				}
			l92:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('O') {
						goto l78
					}
					position++
				}
			l94:
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('U') {
						goto l78
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
						goto l78
					}
					position++
				}
			l98:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('C') {
						goto l78
					}
					position++
				}
			l100:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('E') {
						goto l78
					}
					position++
				}
			l102:
				if !_rules[rulesp]() {
					goto l78
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l78
				}
				if !_rules[rulesp]() {
					goto l78
				}
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
						goto l78
					}
					position++
				}
			l104:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('Y') {
						goto l78
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('P') {
						goto l78
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('E') {
						goto l78
					}
					position++
				}
			l110:
				if !_rules[rulesp]() {
					goto l78
				}
				if !_rules[ruleSourceSinkType]() {
					goto l78
				}
				if !_rules[rulesp]() {
					goto l78
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l78
				}
				if !_rules[ruleAction2]() {
					goto l78
				}
				depth--
				add(ruleCreateSourceStmt, position79)
			}
			return true
		l78:
			position, tokenIndex, depth = position78, tokenIndex78, depth78
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position112, tokenIndex112, depth112 := position, tokenIndex, depth
			{
				position113 := position
				depth++
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('C') {
						goto l112
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('R') {
						goto l112
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
						goto l112
					}
					position++
				}
			l118:
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
					if buffer[position] != rune('A') {
						goto l112
					}
					position++
				}
			l120:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('T') {
						goto l112
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('E') {
						goto l112
					}
					position++
				}
			l124:
				if !_rules[rulesp]() {
					goto l112
				}
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('S') {
						goto l112
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('I') {
						goto l112
					}
					position++
				}
			l128:
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('N') {
						goto l112
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('K') {
						goto l112
					}
					position++
				}
			l132:
				if !_rules[rulesp]() {
					goto l112
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l112
				}
				if !_rules[rulesp]() {
					goto l112
				}
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('T') {
						goto l112
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('Y') {
						goto l112
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('P') {
						goto l112
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
						goto l112
					}
					position++
				}
			l140:
				if !_rules[rulesp]() {
					goto l112
				}
				if !_rules[ruleSourceSinkType]() {
					goto l112
				}
				if !_rules[rulesp]() {
					goto l112
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l112
				}
				if !_rules[ruleAction3]() {
					goto l112
				}
				depth--
				add(ruleCreateSinkStmt, position113)
			}
			return true
		l112:
			position, tokenIndex, depth = position112, tokenIndex112, depth112
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position142, tokenIndex142, depth142 := position, tokenIndex, depth
			{
				position143 := position
				depth++
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('C') {
						goto l142
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('R') {
						goto l142
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('E') {
						goto l142
					}
					position++
				}
			l148:
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('A') {
						goto l142
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('T') {
						goto l142
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('E') {
						goto l142
					}
					position++
				}
			l154:
				if !_rules[rulesp]() {
					goto l142
				}
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('S') {
						goto l142
					}
					position++
				}
			l156:
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if buffer[position] != rune('T') {
						goto l142
					}
					position++
				}
			l158:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('A') {
						goto l142
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
						goto l142
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
						goto l142
					}
					position++
				}
			l164:
				if !_rules[rulesp]() {
					goto l142
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l142
				}
				if !_rules[rulesp]() {
					goto l142
				}
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
						goto l142
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('Y') {
						goto l142
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('P') {
						goto l142
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('E') {
						goto l142
					}
					position++
				}
			l172:
				if !_rules[rulesp]() {
					goto l142
				}
				if !_rules[ruleSourceSinkType]() {
					goto l142
				}
				if !_rules[rulesp]() {
					goto l142
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l142
				}
				if !_rules[ruleAction4]() {
					goto l142
				}
				depth--
				add(ruleCreateStateStmt, position143)
			}
			return true
		l142:
			position, tokenIndex, depth = position142, tokenIndex142, depth142
			return false
		},
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action5)> */
		func() bool {
			position174, tokenIndex174, depth174 := position, tokenIndex, depth
			{
				position175 := position
				depth++
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('I') {
						goto l174
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('N') {
						goto l174
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('S') {
						goto l174
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('E') {
						goto l174
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('R') {
						goto l174
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
						goto l174
					}
					position++
				}
			l186:
				if !_rules[rulesp]() {
					goto l174
				}
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('I') {
						goto l174
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('N') {
						goto l174
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
						goto l174
					}
					position++
				}
			l192:
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('O') {
						goto l174
					}
					position++
				}
			l194:
				if !_rules[rulesp]() {
					goto l174
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l174
				}
				if !_rules[rulesp]() {
					goto l174
				}
				if !_rules[ruleSelectStmt]() {
					goto l174
				}
				if !_rules[ruleAction5]() {
					goto l174
				}
				depth--
				add(ruleInsertIntoSelectStmt, position175)
			}
			return true
		l174:
			position, tokenIndex, depth = position174, tokenIndex174, depth174
			return false
		},
		/* 8 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action6)> */
		func() bool {
			position196, tokenIndex196, depth196 := position, tokenIndex, depth
			{
				position197 := position
				depth++
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('I') {
						goto l196
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('N') {
						goto l196
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
						goto l196
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
						goto l196
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('R') {
						goto l196
					}
					position++
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('T') {
						goto l196
					}
					position++
				}
			l208:
				if !_rules[rulesp]() {
					goto l196
				}
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('I') {
						goto l196
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('N') {
						goto l196
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('T') {
						goto l196
					}
					position++
				}
			l214:
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('O') {
						goto l196
					}
					position++
				}
			l216:
				if !_rules[rulesp]() {
					goto l196
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l196
				}
				if !_rules[rulesp]() {
					goto l196
				}
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('F') {
						goto l196
					}
					position++
				}
			l218:
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
						goto l196
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('O') {
						goto l196
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('M') {
						goto l196
					}
					position++
				}
			l224:
				if !_rules[rulesp]() {
					goto l196
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l196
				}
				if !_rules[ruleAction6]() {
					goto l196
				}
				depth--
				add(ruleInsertIntoFromStmt, position197)
			}
			return true
		l196:
			position, tokenIndex, depth = position196, tokenIndex196, depth196
			return false
		},
		/* 9 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action7)> */
		func() bool {
			position226, tokenIndex226, depth226 := position, tokenIndex, depth
			{
				position227 := position
				depth++
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('P') {
						goto l226
					}
					position++
				}
			l228:
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
					if buffer[position] != rune('A') {
						goto l226
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
						goto l226
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('S') {
						goto l226
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('E') {
						goto l226
					}
					position++
				}
			l236:
				if !_rules[rulesp]() {
					goto l226
				}
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('S') {
						goto l226
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('O') {
						goto l226
					}
					position++
				}
			l240:
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
						goto l226
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('R') {
						goto l226
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('C') {
						goto l226
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
						goto l226
					}
					position++
				}
			l248:
				if !_rules[rulesp]() {
					goto l226
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l226
				}
				if !_rules[ruleAction7]() {
					goto l226
				}
				depth--
				add(rulePauseSourceStmt, position227)
			}
			return true
		l226:
			position, tokenIndex, depth = position226, tokenIndex226, depth226
			return false
		},
		/* 10 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action8)> */
		func() bool {
			position250, tokenIndex250, depth250 := position, tokenIndex, depth
			{
				position251 := position
				depth++
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if buffer[position] != rune('R') {
						goto l250
					}
					position++
				}
			l252:
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if buffer[position] != rune('E') {
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
					if buffer[position] != rune('u') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if buffer[position] != rune('U') {
						goto l250
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('M') {
						goto l250
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
					if buffer[position] != rune('s') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('S') {
						goto l250
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
						goto l250
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
						goto l250
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
						goto l250
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
						goto l250
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
						goto l250
					}
					position++
				}
			l274:
				if !_rules[rulesp]() {
					goto l250
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l250
				}
				if !_rules[ruleAction8]() {
					goto l250
				}
				depth--
				add(ruleResumeSourceStmt, position251)
			}
			return true
		l250:
			position, tokenIndex, depth = position250, tokenIndex250, depth250
			return false
		},
		/* 11 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action9)> */
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
					if buffer[position] != rune('w') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('W') {
						goto l276
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
					if buffer[position] != rune('I') {
						goto l276
					}
					position++
				}
			l284:
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l287
					}
					position++
					goto l286
				l287:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
					if buffer[position] != rune('N') {
						goto l276
					}
					position++
				}
			l286:
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l289
					}
					position++
					goto l288
				l289:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
					if buffer[position] != rune('D') {
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
				add(ruleRewindSourceStmt, position277)
			}
			return true
		l276:
			position, tokenIndex, depth = position276, tokenIndex276, depth276
			return false
		},
		/* 12 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position302, tokenIndex302, depth302 := position, tokenIndex, depth
			{
				position303 := position
				depth++
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
					if buffer[position] != rune('D') {
						goto l302
					}
					position++
				}
			l304:
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l307
					}
					position++
					goto l306
				l307:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					if buffer[position] != rune('R') {
						goto l302
					}
					position++
				}
			l306:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					if buffer[position] != rune('O') {
						goto l302
					}
					position++
				}
			l308:
				{
					position310, tokenIndex310, depth310 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l311
					}
					position++
					goto l310
				l311:
					position, tokenIndex, depth = position310, tokenIndex310, depth310
					if buffer[position] != rune('P') {
						goto l302
					}
					position++
				}
			l310:
				if !_rules[rulesp]() {
					goto l302
				}
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('S') {
						goto l302
					}
					position++
				}
			l312:
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l315
					}
					position++
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if buffer[position] != rune('O') {
						goto l302
					}
					position++
				}
			l314:
				{
					position316, tokenIndex316, depth316 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l317
					}
					position++
					goto l316
				l317:
					position, tokenIndex, depth = position316, tokenIndex316, depth316
					if buffer[position] != rune('U') {
						goto l302
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
						goto l302
					}
					position++
				}
			l318:
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('C') {
						goto l302
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if buffer[position] != rune('E') {
						goto l302
					}
					position++
				}
			l322:
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
				add(ruleDropSourceStmt, position303)
			}
			return true
		l302:
			position, tokenIndex, depth = position302, tokenIndex302, depth302
			return false
		},
		/* 13 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action11)> */
		func() bool {
			position324, tokenIndex324, depth324 := position, tokenIndex, depth
			{
				position325 := position
				depth++
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l327
					}
					position++
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					if buffer[position] != rune('D') {
						goto l324
					}
					position++
				}
			l326:
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
						goto l324
					}
					position++
				}
			l328:
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l331
					}
					position++
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					if buffer[position] != rune('O') {
						goto l324
					}
					position++
				}
			l330:
				{
					position332, tokenIndex332, depth332 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l333
					}
					position++
					goto l332
				l333:
					position, tokenIndex, depth = position332, tokenIndex332, depth332
					if buffer[position] != rune('P') {
						goto l324
					}
					position++
				}
			l332:
				if !_rules[rulesp]() {
					goto l324
				}
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if buffer[position] != rune('S') {
						goto l324
					}
					position++
				}
			l334:
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if buffer[position] != rune('T') {
						goto l324
					}
					position++
				}
			l336:
				{
					position338, tokenIndex338, depth338 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l339
					}
					position++
					goto l338
				l339:
					position, tokenIndex, depth = position338, tokenIndex338, depth338
					if buffer[position] != rune('R') {
						goto l324
					}
					position++
				}
			l338:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l341
					}
					position++
					goto l340
				l341:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
					if buffer[position] != rune('E') {
						goto l324
					}
					position++
				}
			l340:
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if buffer[position] != rune('A') {
						goto l324
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
					if buffer[position] != rune('M') {
						goto l324
					}
					position++
				}
			l344:
				if !_rules[rulesp]() {
					goto l324
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l324
				}
				if !_rules[ruleAction11]() {
					goto l324
				}
				depth--
				add(ruleDropStreamStmt, position325)
			}
			return true
		l324:
			position, tokenIndex, depth = position324, tokenIndex324, depth324
			return false
		},
		/* 14 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action12)> */
		func() bool {
			position346, tokenIndex346, depth346 := position, tokenIndex, depth
			{
				position347 := position
				depth++
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if buffer[position] != rune('D') {
						goto l346
					}
					position++
				}
			l348:
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if buffer[position] != rune('R') {
						goto l346
					}
					position++
				}
			l350:
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('O') {
						goto l346
					}
					position++
				}
			l352:
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
					if buffer[position] != rune('P') {
						goto l346
					}
					position++
				}
			l354:
				if !_rules[rulesp]() {
					goto l346
				}
				{
					position356, tokenIndex356, depth356 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l357
					}
					position++
					goto l356
				l357:
					position, tokenIndex, depth = position356, tokenIndex356, depth356
					if buffer[position] != rune('S') {
						goto l346
					}
					position++
				}
			l356:
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l359
					}
					position++
					goto l358
				l359:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					if buffer[position] != rune('I') {
						goto l346
					}
					position++
				}
			l358:
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if buffer[position] != rune('N') {
						goto l346
					}
					position++
				}
			l360:
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l363
					}
					position++
					goto l362
				l363:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
					if buffer[position] != rune('K') {
						goto l346
					}
					position++
				}
			l362:
				if !_rules[rulesp]() {
					goto l346
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l346
				}
				if !_rules[ruleAction12]() {
					goto l346
				}
				depth--
				add(ruleDropSinkStmt, position347)
			}
			return true
		l346:
			position, tokenIndex, depth = position346, tokenIndex346, depth346
			return false
		},
		/* 15 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action13)> */
		func() bool {
			position364, tokenIndex364, depth364 := position, tokenIndex, depth
			{
				position365 := position
				depth++
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l367
					}
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if !_rules[ruleDSTREAM]() {
						goto l368
					}
					goto l366
				l368:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if !_rules[ruleRSTREAM]() {
						goto l364
					}
				}
			l366:
				{
					position369 := position
					depth++
					{
						position370, tokenIndex370, depth370 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l370
						}
						if buffer[position] != rune('[') {
							goto l370
						}
						position++
						if !_rules[rulesp]() {
							goto l370
						}
						{
							position372, tokenIndex372, depth372 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l373
							}
							position++
							goto l372
						l373:
							position, tokenIndex, depth = position372, tokenIndex372, depth372
							if buffer[position] != rune('E') {
								goto l370
							}
							position++
						}
					l372:
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l375
							}
							position++
							goto l374
						l375:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
							if buffer[position] != rune('V') {
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
							if buffer[position] != rune('y') {
								goto l381
							}
							position++
							goto l380
						l381:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
							if buffer[position] != rune('Y') {
								goto l370
							}
							position++
						}
					l380:
						if !_rules[rulesp]() {
							goto l370
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l370
						}
						if !_rules[rulesp]() {
							goto l370
						}
						if buffer[position] != rune(']') {
							goto l370
						}
						position++
						goto l371
					l370:
						position, tokenIndex, depth = position370, tokenIndex370, depth370
					}
				l371:
					depth--
					add(rulePegText, position369)
				}
				if !_rules[ruleAction13]() {
					goto l364
				}
				depth--
				add(ruleEmitter, position365)
			}
			return true
		l364:
			position, tokenIndex, depth = position364, tokenIndex364, depth364
			return false
		},
		/* 16 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position382, tokenIndex382, depth382 := position, tokenIndex, depth
			{
				position383 := position
				depth++
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l385
					}
				l386:
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l387
						}
						if buffer[position] != rune(',') {
							goto l387
						}
						position++
						if !_rules[rulesp]() {
							goto l387
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l387
						}
						goto l386
					l387:
						position, tokenIndex, depth = position387, tokenIndex387, depth387
					}
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if !_rules[ruleTimeEmitterInterval]() {
						goto l388
					}
					goto l384
				l388:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if !_rules[ruleTupleEmitterInterval]() {
						goto l382
					}
				}
			l384:
				depth--
				add(ruleEmitterIntervals, position383)
			}
			return true
		l382:
			position, tokenIndex, depth = position382, tokenIndex382, depth382
			return false
		},
		/* 17 TimeEmitterInterval <- <(<TimeInterval> Action14)> */
		func() bool {
			position389, tokenIndex389, depth389 := position, tokenIndex, depth
			{
				position390 := position
				depth++
				{
					position391 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l389
					}
					depth--
					add(rulePegText, position391)
				}
				if !_rules[ruleAction14]() {
					goto l389
				}
				depth--
				add(ruleTimeEmitterInterval, position390)
			}
			return true
		l389:
			position, tokenIndex, depth = position389, tokenIndex389, depth389
			return false
		},
		/* 18 TupleEmitterInterval <- <(<TuplesInterval> Action15)> */
		func() bool {
			position392, tokenIndex392, depth392 := position, tokenIndex, depth
			{
				position393 := position
				depth++
				{
					position394 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l392
					}
					depth--
					add(rulePegText, position394)
				}
				if !_rules[ruleAction15]() {
					goto l392
				}
				depth--
				add(ruleTupleEmitterInterval, position393)
			}
			return true
		l392:
			position, tokenIndex, depth = position392, tokenIndex392, depth392
			return false
		},
		/* 19 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action16)> */
		func() bool {
			position395, tokenIndex395, depth395 := position, tokenIndex, depth
			{
				position396 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l395
				}
				if !_rules[rulesp]() {
					goto l395
				}
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l398
					}
					position++
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if buffer[position] != rune('I') {
						goto l395
					}
					position++
				}
			l397:
				{
					position399, tokenIndex399, depth399 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l400
					}
					position++
					goto l399
				l400:
					position, tokenIndex, depth = position399, tokenIndex399, depth399
					if buffer[position] != rune('N') {
						goto l395
					}
					position++
				}
			l399:
				if !_rules[rulesp]() {
					goto l395
				}
				if !_rules[ruleStream]() {
					goto l395
				}
				if !_rules[ruleAction16]() {
					goto l395
				}
				depth--
				add(ruleTupleEmitterFromInterval, position396)
			}
			return true
		l395:
			position, tokenIndex, depth = position395, tokenIndex395, depth395
			return false
		},
		/* 20 Projections <- <(<(Projection sp (',' sp Projection)*)> Action17)> */
		func() bool {
			position401, tokenIndex401, depth401 := position, tokenIndex, depth
			{
				position402 := position
				depth++
				{
					position403 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l401
					}
					if !_rules[rulesp]() {
						goto l401
					}
				l404:
					{
						position405, tokenIndex405, depth405 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l405
						}
						position++
						if !_rules[rulesp]() {
							goto l405
						}
						if !_rules[ruleProjection]() {
							goto l405
						}
						goto l404
					l405:
						position, tokenIndex, depth = position405, tokenIndex405, depth405
					}
					depth--
					add(rulePegText, position403)
				}
				if !_rules[ruleAction17]() {
					goto l401
				}
				depth--
				add(ruleProjections, position402)
			}
			return true
		l401:
			position, tokenIndex, depth = position401, tokenIndex401, depth401
			return false
		},
		/* 21 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position406, tokenIndex406, depth406 := position, tokenIndex, depth
			{
				position407 := position
				depth++
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l409
					}
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if !_rules[ruleExpression]() {
						goto l410
					}
					goto l408
				l410:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if !_rules[ruleWildcard]() {
						goto l406
					}
				}
			l408:
				depth--
				add(ruleProjection, position407)
			}
			return true
		l406:
			position, tokenIndex, depth = position406, tokenIndex406, depth406
			return false
		},
		/* 22 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action18)> */
		func() bool {
			position411, tokenIndex411, depth411 := position, tokenIndex, depth
			{
				position412 := position
				depth++
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l414
					}
					goto l413
				l414:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
					if !_rules[ruleWildcard]() {
						goto l411
					}
				}
			l413:
				if !_rules[rulesp]() {
					goto l411
				}
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l416
					}
					position++
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if buffer[position] != rune('A') {
						goto l411
					}
					position++
				}
			l415:
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l418
					}
					position++
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if buffer[position] != rune('S') {
						goto l411
					}
					position++
				}
			l417:
				if !_rules[rulesp]() {
					goto l411
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l411
				}
				if !_rules[ruleAction18]() {
					goto l411
				}
				depth--
				add(ruleAliasExpression, position412)
			}
			return true
		l411:
			position, tokenIndex, depth = position411, tokenIndex411, depth411
			return false
		},
		/* 23 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action19)> */
		func() bool {
			position419, tokenIndex419, depth419 := position, tokenIndex, depth
			{
				position420 := position
				depth++
				{
					position421 := position
					depth++
					{
						position422, tokenIndex422, depth422 := position, tokenIndex, depth
						{
							position424, tokenIndex424, depth424 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l425
							}
							position++
							goto l424
						l425:
							position, tokenIndex, depth = position424, tokenIndex424, depth424
							if buffer[position] != rune('F') {
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
							if buffer[position] != rune('m') {
								goto l431
							}
							position++
							goto l430
						l431:
							position, tokenIndex, depth = position430, tokenIndex430, depth430
							if buffer[position] != rune('M') {
								goto l422
							}
							position++
						}
					l430:
						if !_rules[rulesp]() {
							goto l422
						}
						if !_rules[ruleRelations]() {
							goto l422
						}
						if !_rules[rulesp]() {
							goto l422
						}
						goto l423
					l422:
						position, tokenIndex, depth = position422, tokenIndex422, depth422
					}
				l423:
					depth--
					add(rulePegText, position421)
				}
				if !_rules[ruleAction19]() {
					goto l419
				}
				depth--
				add(ruleWindowedFrom, position420)
			}
			return true
		l419:
			position, tokenIndex, depth = position419, tokenIndex419, depth419
			return false
		},
		/* 24 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l435
					}
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if !_rules[ruleTuplesInterval]() {
						goto l432
					}
				}
			l434:
				depth--
				add(ruleInterval, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 25 TimeInterval <- <(NumericLiteral sp SECONDS Action20)> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l436
				}
				if !_rules[rulesp]() {
					goto l436
				}
				if !_rules[ruleSECONDS]() {
					goto l436
				}
				if !_rules[ruleAction20]() {
					goto l436
				}
				depth--
				add(ruleTimeInterval, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
			return false
		},
		/* 26 TuplesInterval <- <(NumericLiteral sp TUPLES Action21)> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l438
				}
				if !_rules[rulesp]() {
					goto l438
				}
				if !_rules[ruleTUPLES]() {
					goto l438
				}
				if !_rules[ruleAction21]() {
					goto l438
				}
				depth--
				add(ruleTuplesInterval, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 27 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l440
				}
				if !_rules[rulesp]() {
					goto l440
				}
			l442:
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l443
					}
					position++
					if !_rules[rulesp]() {
						goto l443
					}
					if !_rules[ruleRelationLike]() {
						goto l443
					}
					goto l442
				l443:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
				}
				depth--
				add(ruleRelations, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 28 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action22)> */
		func() bool {
			position444, tokenIndex444, depth444 := position, tokenIndex, depth
			{
				position445 := position
				depth++
				{
					position446 := position
					depth++
					{
						position447, tokenIndex447, depth447 := position, tokenIndex, depth
						{
							position449, tokenIndex449, depth449 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l450
							}
							position++
							goto l449
						l450:
							position, tokenIndex, depth = position449, tokenIndex449, depth449
							if buffer[position] != rune('W') {
								goto l447
							}
							position++
						}
					l449:
						{
							position451, tokenIndex451, depth451 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l452
							}
							position++
							goto l451
						l452:
							position, tokenIndex, depth = position451, tokenIndex451, depth451
							if buffer[position] != rune('H') {
								goto l447
							}
							position++
						}
					l451:
						{
							position453, tokenIndex453, depth453 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l454
							}
							position++
							goto l453
						l454:
							position, tokenIndex, depth = position453, tokenIndex453, depth453
							if buffer[position] != rune('E') {
								goto l447
							}
							position++
						}
					l453:
						{
							position455, tokenIndex455, depth455 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l456
							}
							position++
							goto l455
						l456:
							position, tokenIndex, depth = position455, tokenIndex455, depth455
							if buffer[position] != rune('R') {
								goto l447
							}
							position++
						}
					l455:
						{
							position457, tokenIndex457, depth457 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l458
							}
							position++
							goto l457
						l458:
							position, tokenIndex, depth = position457, tokenIndex457, depth457
							if buffer[position] != rune('E') {
								goto l447
							}
							position++
						}
					l457:
						if !_rules[rulesp]() {
							goto l447
						}
						if !_rules[ruleExpression]() {
							goto l447
						}
						goto l448
					l447:
						position, tokenIndex, depth = position447, tokenIndex447, depth447
					}
				l448:
					depth--
					add(rulePegText, position446)
				}
				if !_rules[ruleAction22]() {
					goto l444
				}
				depth--
				add(ruleFilter, position445)
			}
			return true
		l444:
			position, tokenIndex, depth = position444, tokenIndex444, depth444
			return false
		},
		/* 29 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action23)> */
		func() bool {
			position459, tokenIndex459, depth459 := position, tokenIndex, depth
			{
				position460 := position
				depth++
				{
					position461 := position
					depth++
					{
						position462, tokenIndex462, depth462 := position, tokenIndex, depth
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
								goto l462
							}
							position++
						}
					l464:
						{
							position466, tokenIndex466, depth466 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l467
							}
							position++
							goto l466
						l467:
							position, tokenIndex, depth = position466, tokenIndex466, depth466
							if buffer[position] != rune('R') {
								goto l462
							}
							position++
						}
					l466:
						{
							position468, tokenIndex468, depth468 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l469
							}
							position++
							goto l468
						l469:
							position, tokenIndex, depth = position468, tokenIndex468, depth468
							if buffer[position] != rune('O') {
								goto l462
							}
							position++
						}
					l468:
						{
							position470, tokenIndex470, depth470 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l471
							}
							position++
							goto l470
						l471:
							position, tokenIndex, depth = position470, tokenIndex470, depth470
							if buffer[position] != rune('U') {
								goto l462
							}
							position++
						}
					l470:
						{
							position472, tokenIndex472, depth472 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l473
							}
							position++
							goto l472
						l473:
							position, tokenIndex, depth = position472, tokenIndex472, depth472
							if buffer[position] != rune('P') {
								goto l462
							}
							position++
						}
					l472:
						if !_rules[rulesp]() {
							goto l462
						}
						{
							position474, tokenIndex474, depth474 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l475
							}
							position++
							goto l474
						l475:
							position, tokenIndex, depth = position474, tokenIndex474, depth474
							if buffer[position] != rune('B') {
								goto l462
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
								goto l462
							}
							position++
						}
					l476:
						if !_rules[rulesp]() {
							goto l462
						}
						if !_rules[ruleGroupList]() {
							goto l462
						}
						goto l463
					l462:
						position, tokenIndex, depth = position462, tokenIndex462, depth462
					}
				l463:
					depth--
					add(rulePegText, position461)
				}
				if !_rules[ruleAction23]() {
					goto l459
				}
				depth--
				add(ruleGrouping, position460)
			}
			return true
		l459:
			position, tokenIndex, depth = position459, tokenIndex459, depth459
			return false
		},
		/* 30 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l478
				}
				if !_rules[rulesp]() {
					goto l478
				}
			l480:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l481
					}
					position++
					if !_rules[rulesp]() {
						goto l481
					}
					if !_rules[ruleExpression]() {
						goto l481
					}
					goto l480
				l481:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
				}
				depth--
				add(ruleGroupList, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 31 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action24)> */
		func() bool {
			position482, tokenIndex482, depth482 := position, tokenIndex, depth
			{
				position483 := position
				depth++
				{
					position484 := position
					depth++
					{
						position485, tokenIndex485, depth485 := position, tokenIndex, depth
						{
							position487, tokenIndex487, depth487 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l488
							}
							position++
							goto l487
						l488:
							position, tokenIndex, depth = position487, tokenIndex487, depth487
							if buffer[position] != rune('H') {
								goto l485
							}
							position++
						}
					l487:
						{
							position489, tokenIndex489, depth489 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l490
							}
							position++
							goto l489
						l490:
							position, tokenIndex, depth = position489, tokenIndex489, depth489
							if buffer[position] != rune('A') {
								goto l485
							}
							position++
						}
					l489:
						{
							position491, tokenIndex491, depth491 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l492
							}
							position++
							goto l491
						l492:
							position, tokenIndex, depth = position491, tokenIndex491, depth491
							if buffer[position] != rune('V') {
								goto l485
							}
							position++
						}
					l491:
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
								goto l485
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
								goto l485
							}
							position++
						}
					l495:
						{
							position497, tokenIndex497, depth497 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l498
							}
							position++
							goto l497
						l498:
							position, tokenIndex, depth = position497, tokenIndex497, depth497
							if buffer[position] != rune('G') {
								goto l485
							}
							position++
						}
					l497:
						if !_rules[rulesp]() {
							goto l485
						}
						if !_rules[ruleExpression]() {
							goto l485
						}
						goto l486
					l485:
						position, tokenIndex, depth = position485, tokenIndex485, depth485
					}
				l486:
					depth--
					add(rulePegText, position484)
				}
				if !_rules[ruleAction24]() {
					goto l482
				}
				depth--
				add(ruleHaving, position483)
			}
			return true
		l482:
			position, tokenIndex, depth = position482, tokenIndex482, depth482
			return false
		},
		/* 32 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action25))> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				{
					position501, tokenIndex501, depth501 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l502
					}
					goto l501
				l502:
					position, tokenIndex, depth = position501, tokenIndex501, depth501
					if !_rules[ruleStreamWindow]() {
						goto l499
					}
					if !_rules[ruleAction25]() {
						goto l499
					}
				}
			l501:
				depth--
				add(ruleRelationLike, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 33 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action26)> */
		func() bool {
			position503, tokenIndex503, depth503 := position, tokenIndex, depth
			{
				position504 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l503
				}
				if !_rules[rulesp]() {
					goto l503
				}
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l506
					}
					position++
					goto l505
				l506:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
					if buffer[position] != rune('A') {
						goto l503
					}
					position++
				}
			l505:
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('S') {
						goto l503
					}
					position++
				}
			l507:
				if !_rules[rulesp]() {
					goto l503
				}
				if !_rules[ruleIdentifier]() {
					goto l503
				}
				if !_rules[ruleAction26]() {
					goto l503
				}
				depth--
				add(ruleAliasedStreamWindow, position504)
			}
			return true
		l503:
			position, tokenIndex, depth = position503, tokenIndex503, depth503
			return false
		},
		/* 34 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action27)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l509
				}
				if !_rules[rulesp]() {
					goto l509
				}
				if buffer[position] != rune('[') {
					goto l509
				}
				position++
				if !_rules[rulesp]() {
					goto l509
				}
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if buffer[position] != rune('R') {
						goto l509
					}
					position++
				}
			l511:
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l514
					}
					position++
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if buffer[position] != rune('A') {
						goto l509
					}
					position++
				}
			l513:
				{
					position515, tokenIndex515, depth515 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l516
					}
					position++
					goto l515
				l516:
					position, tokenIndex, depth = position515, tokenIndex515, depth515
					if buffer[position] != rune('N') {
						goto l509
					}
					position++
				}
			l515:
				{
					position517, tokenIndex517, depth517 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l518
					}
					position++
					goto l517
				l518:
					position, tokenIndex, depth = position517, tokenIndex517, depth517
					if buffer[position] != rune('G') {
						goto l509
					}
					position++
				}
			l517:
				{
					position519, tokenIndex519, depth519 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l520
					}
					position++
					goto l519
				l520:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if buffer[position] != rune('E') {
						goto l509
					}
					position++
				}
			l519:
				if !_rules[rulesp]() {
					goto l509
				}
				if !_rules[ruleInterval]() {
					goto l509
				}
				if !_rules[rulesp]() {
					goto l509
				}
				if buffer[position] != rune(']') {
					goto l509
				}
				position++
				if !_rules[ruleAction27]() {
					goto l509
				}
				depth--
				add(ruleStreamWindow, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 35 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position521, tokenIndex521, depth521 := position, tokenIndex, depth
			{
				position522 := position
				depth++
				{
					position523, tokenIndex523, depth523 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l524
					}
					goto l523
				l524:
					position, tokenIndex, depth = position523, tokenIndex523, depth523
					if !_rules[ruleStream]() {
						goto l521
					}
				}
			l523:
				depth--
				add(ruleStreamLike, position522)
			}
			return true
		l521:
			position, tokenIndex, depth = position521, tokenIndex521, depth521
			return false
		},
		/* 36 UDSFFuncApp <- <(FuncApp Action28)> */
		func() bool {
			position525, tokenIndex525, depth525 := position, tokenIndex, depth
			{
				position526 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l525
				}
				if !_rules[ruleAction28]() {
					goto l525
				}
				depth--
				add(ruleUDSFFuncApp, position526)
			}
			return true
		l525:
			position, tokenIndex, depth = position525, tokenIndex525, depth525
			return false
		},
		/* 37 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action29)> */
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
							if buffer[position] != rune('w') {
								goto l533
							}
							position++
							goto l532
						l533:
							position, tokenIndex, depth = position532, tokenIndex532, depth532
							if buffer[position] != rune('W') {
								goto l530
							}
							position++
						}
					l532:
						{
							position534, tokenIndex534, depth534 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l535
							}
							position++
							goto l534
						l535:
							position, tokenIndex, depth = position534, tokenIndex534, depth534
							if buffer[position] != rune('I') {
								goto l530
							}
							position++
						}
					l534:
						{
							position536, tokenIndex536, depth536 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l537
							}
							position++
							goto l536
						l537:
							position, tokenIndex, depth = position536, tokenIndex536, depth536
							if buffer[position] != rune('T') {
								goto l530
							}
							position++
						}
					l536:
						{
							position538, tokenIndex538, depth538 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l539
							}
							position++
							goto l538
						l539:
							position, tokenIndex, depth = position538, tokenIndex538, depth538
							if buffer[position] != rune('H') {
								goto l530
							}
							position++
						}
					l538:
						if !_rules[rulesp]() {
							goto l530
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l530
						}
						if !_rules[rulesp]() {
							goto l530
						}
					l540:
						{
							position541, tokenIndex541, depth541 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l541
							}
							position++
							if !_rules[rulesp]() {
								goto l541
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l541
							}
							goto l540
						l541:
							position, tokenIndex, depth = position541, tokenIndex541, depth541
						}
						goto l531
					l530:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
					}
				l531:
					depth--
					add(rulePegText, position529)
				}
				if !_rules[ruleAction29]() {
					goto l527
				}
				depth--
				add(ruleSourceSinkSpecs, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 38 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action30)> */
		func() bool {
			position542, tokenIndex542, depth542 := position, tokenIndex, depth
			{
				position543 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l542
				}
				if buffer[position] != rune('=') {
					goto l542
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l542
				}
				if !_rules[ruleAction30]() {
					goto l542
				}
				depth--
				add(ruleSourceSinkParam, position543)
			}
			return true
		l542:
			position, tokenIndex, depth = position542, tokenIndex542, depth542
			return false
		},
		/* 39 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				{
					position546, tokenIndex546, depth546 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l547
					}
					goto l546
				l547:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
					if !_rules[ruleLiteral]() {
						goto l544
					}
				}
			l546:
				depth--
				add(ruleSourceSinkParamVal, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 40 PausedOpt <- <(<(Paused / Unpaused)?> Action31)> */
		func() bool {
			position548, tokenIndex548, depth548 := position, tokenIndex, depth
			{
				position549 := position
				depth++
				{
					position550 := position
					depth++
					{
						position551, tokenIndex551, depth551 := position, tokenIndex, depth
						{
							position553, tokenIndex553, depth553 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l554
							}
							goto l553
						l554:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if !_rules[ruleUnpaused]() {
								goto l551
							}
						}
					l553:
						goto l552
					l551:
						position, tokenIndex, depth = position551, tokenIndex551, depth551
					}
				l552:
					depth--
					add(rulePegText, position550)
				}
				if !_rules[ruleAction31]() {
					goto l548
				}
				depth--
				add(rulePausedOpt, position549)
			}
			return true
		l548:
			position, tokenIndex, depth = position548, tokenIndex548, depth548
			return false
		},
		/* 41 Expression <- <orExpr> */
		func() bool {
			position555, tokenIndex555, depth555 := position, tokenIndex, depth
			{
				position556 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l555
				}
				depth--
				add(ruleExpression, position556)
			}
			return true
		l555:
			position, tokenIndex, depth = position555, tokenIndex555, depth555
			return false
		},
		/* 42 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action32)> */
		func() bool {
			position557, tokenIndex557, depth557 := position, tokenIndex, depth
			{
				position558 := position
				depth++
				{
					position559 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l557
					}
					if !_rules[rulesp]() {
						goto l557
					}
					{
						position560, tokenIndex560, depth560 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l560
						}
						if !_rules[rulesp]() {
							goto l560
						}
						if !_rules[ruleandExpr]() {
							goto l560
						}
						goto l561
					l560:
						position, tokenIndex, depth = position560, tokenIndex560, depth560
					}
				l561:
					depth--
					add(rulePegText, position559)
				}
				if !_rules[ruleAction32]() {
					goto l557
				}
				depth--
				add(ruleorExpr, position558)
			}
			return true
		l557:
			position, tokenIndex, depth = position557, tokenIndex557, depth557
			return false
		},
		/* 43 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action33)> */
		func() bool {
			position562, tokenIndex562, depth562 := position, tokenIndex, depth
			{
				position563 := position
				depth++
				{
					position564 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l562
					}
					if !_rules[rulesp]() {
						goto l562
					}
					{
						position565, tokenIndex565, depth565 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l565
						}
						if !_rules[rulesp]() {
							goto l565
						}
						if !_rules[rulenotExpr]() {
							goto l565
						}
						goto l566
					l565:
						position, tokenIndex, depth = position565, tokenIndex565, depth565
					}
				l566:
					depth--
					add(rulePegText, position564)
				}
				if !_rules[ruleAction33]() {
					goto l562
				}
				depth--
				add(ruleandExpr, position563)
			}
			return true
		l562:
			position, tokenIndex, depth = position562, tokenIndex562, depth562
			return false
		},
		/* 44 notExpr <- <(<((Not sp)? comparisonExpr)> Action34)> */
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
						if !_rules[ruleNot]() {
							goto l570
						}
						if !_rules[rulesp]() {
							goto l570
						}
						goto l571
					l570:
						position, tokenIndex, depth = position570, tokenIndex570, depth570
					}
				l571:
					if !_rules[rulecomparisonExpr]() {
						goto l567
					}
					depth--
					add(rulePegText, position569)
				}
				if !_rules[ruleAction34]() {
					goto l567
				}
				depth--
				add(rulenotExpr, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 45 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action35)> */
		func() bool {
			position572, tokenIndex572, depth572 := position, tokenIndex, depth
			{
				position573 := position
				depth++
				{
					position574 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l572
					}
					if !_rules[rulesp]() {
						goto l572
					}
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l575
						}
						if !_rules[rulesp]() {
							goto l575
						}
						if !_rules[ruleisExpr]() {
							goto l575
						}
						goto l576
					l575:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
					}
				l576:
					depth--
					add(rulePegText, position574)
				}
				if !_rules[ruleAction35]() {
					goto l572
				}
				depth--
				add(rulecomparisonExpr, position573)
			}
			return true
		l572:
			position, tokenIndex, depth = position572, tokenIndex572, depth572
			return false
		},
		/* 46 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action36)> */
		func() bool {
			position577, tokenIndex577, depth577 := position, tokenIndex, depth
			{
				position578 := position
				depth++
				{
					position579 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l577
					}
					if !_rules[rulesp]() {
						goto l577
					}
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l580
						}
						if !_rules[rulesp]() {
							goto l580
						}
						if !_rules[ruleNullLiteral]() {
							goto l580
						}
						goto l581
					l580:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
					}
				l581:
					depth--
					add(rulePegText, position579)
				}
				if !_rules[ruleAction36]() {
					goto l577
				}
				depth--
				add(ruleisExpr, position578)
			}
			return true
		l577:
			position, tokenIndex, depth = position577, tokenIndex577, depth577
			return false
		},
		/* 47 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action37)> */
		func() bool {
			position582, tokenIndex582, depth582 := position, tokenIndex, depth
			{
				position583 := position
				depth++
				{
					position584 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l582
					}
					if !_rules[rulesp]() {
						goto l582
					}
					{
						position585, tokenIndex585, depth585 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l585
						}
						if !_rules[rulesp]() {
							goto l585
						}
						if !_rules[ruleproductExpr]() {
							goto l585
						}
						goto l586
					l585:
						position, tokenIndex, depth = position585, tokenIndex585, depth585
					}
				l586:
					depth--
					add(rulePegText, position584)
				}
				if !_rules[ruleAction37]() {
					goto l582
				}
				depth--
				add(ruletermExpr, position583)
			}
			return true
		l582:
			position, tokenIndex, depth = position582, tokenIndex582, depth582
			return false
		},
		/* 48 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action38)> */
		func() bool {
			position587, tokenIndex587, depth587 := position, tokenIndex, depth
			{
				position588 := position
				depth++
				{
					position589 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l587
					}
					if !_rules[rulesp]() {
						goto l587
					}
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l590
						}
						if !_rules[rulesp]() {
							goto l590
						}
						if !_rules[ruleminusExpr]() {
							goto l590
						}
						goto l591
					l590:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
					}
				l591:
					depth--
					add(rulePegText, position589)
				}
				if !_rules[ruleAction38]() {
					goto l587
				}
				depth--
				add(ruleproductExpr, position588)
			}
			return true
		l587:
			position, tokenIndex, depth = position587, tokenIndex587, depth587
			return false
		},
		/* 49 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action39)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l595
						}
						if !_rules[rulesp]() {
							goto l595
						}
						goto l596
					l595:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
					}
				l596:
					if !_rules[rulebaseExpr]() {
						goto l592
					}
					depth--
					add(rulePegText, position594)
				}
				if !_rules[ruleAction39]() {
					goto l592
				}
				depth--
				add(ruleminusExpr, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 50 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position597, tokenIndex597, depth597 := position, tokenIndex, depth
			{
				position598 := position
				depth++
				{
					position599, tokenIndex599, depth599 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l600
					}
					position++
					if !_rules[rulesp]() {
						goto l600
					}
					if !_rules[ruleExpression]() {
						goto l600
					}
					if !_rules[rulesp]() {
						goto l600
					}
					if buffer[position] != rune(')') {
						goto l600
					}
					position++
					goto l599
				l600:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if !_rules[ruleBooleanLiteral]() {
						goto l601
					}
					goto l599
				l601:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if !_rules[ruleNullLiteral]() {
						goto l602
					}
					goto l599
				l602:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if !_rules[ruleFuncApp]() {
						goto l603
					}
					goto l599
				l603:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if !_rules[ruleRowMeta]() {
						goto l604
					}
					goto l599
				l604:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if !_rules[ruleRowValue]() {
						goto l605
					}
					goto l599
				l605:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if !_rules[ruleLiteral]() {
						goto l597
					}
				}
			l599:
				depth--
				add(rulebaseExpr, position598)
			}
			return true
		l597:
			position, tokenIndex, depth = position597, tokenIndex597, depth597
			return false
		},
		/* 51 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action40)> */
		func() bool {
			position606, tokenIndex606, depth606 := position, tokenIndex, depth
			{
				position607 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l606
				}
				if !_rules[rulesp]() {
					goto l606
				}
				if buffer[position] != rune('(') {
					goto l606
				}
				position++
				if !_rules[rulesp]() {
					goto l606
				}
				if !_rules[ruleFuncParams]() {
					goto l606
				}
				if !_rules[rulesp]() {
					goto l606
				}
				if buffer[position] != rune(')') {
					goto l606
				}
				position++
				if !_rules[ruleAction40]() {
					goto l606
				}
				depth--
				add(ruleFuncApp, position607)
			}
			return true
		l606:
			position, tokenIndex, depth = position606, tokenIndex606, depth606
			return false
		},
		/* 52 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action41)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l608
					}
					if !_rules[rulesp]() {
						goto l608
					}
				l611:
					{
						position612, tokenIndex612, depth612 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l612
						}
						position++
						if !_rules[rulesp]() {
							goto l612
						}
						if !_rules[ruleExpression]() {
							goto l612
						}
						goto l611
					l612:
						position, tokenIndex, depth = position612, tokenIndex612, depth612
					}
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction41]() {
					goto l608
				}
				depth--
				add(ruleFuncParams, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 53 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position613, tokenIndex613, depth613 := position, tokenIndex, depth
			{
				position614 := position
				depth++
				{
					position615, tokenIndex615, depth615 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l616
					}
					goto l615
				l616:
					position, tokenIndex, depth = position615, tokenIndex615, depth615
					if !_rules[ruleNumericLiteral]() {
						goto l617
					}
					goto l615
				l617:
					position, tokenIndex, depth = position615, tokenIndex615, depth615
					if !_rules[ruleStringLiteral]() {
						goto l613
					}
				}
			l615:
				depth--
				add(ruleLiteral, position614)
			}
			return true
		l613:
			position, tokenIndex, depth = position613, tokenIndex613, depth613
			return false
		},
		/* 54 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position618, tokenIndex618, depth618 := position, tokenIndex, depth
			{
				position619 := position
				depth++
				{
					position620, tokenIndex620, depth620 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l621
					}
					goto l620
				l621:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleNotEqual]() {
						goto l622
					}
					goto l620
				l622:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleLessOrEqual]() {
						goto l623
					}
					goto l620
				l623:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleLess]() {
						goto l624
					}
					goto l620
				l624:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleGreaterOrEqual]() {
						goto l625
					}
					goto l620
				l625:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleGreater]() {
						goto l626
					}
					goto l620
				l626:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleNotEqual]() {
						goto l618
					}
				}
			l620:
				depth--
				add(ruleComparisonOp, position619)
			}
			return true
		l618:
			position, tokenIndex, depth = position618, tokenIndex618, depth618
			return false
		},
		/* 55 IsOp <- <(IsNot / Is)> */
		func() bool {
			position627, tokenIndex627, depth627 := position, tokenIndex, depth
			{
				position628 := position
				depth++
				{
					position629, tokenIndex629, depth629 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l630
					}
					goto l629
				l630:
					position, tokenIndex, depth = position629, tokenIndex629, depth629
					if !_rules[ruleIs]() {
						goto l627
					}
				}
			l629:
				depth--
				add(ruleIsOp, position628)
			}
			return true
		l627:
			position, tokenIndex, depth = position627, tokenIndex627, depth627
			return false
		},
		/* 56 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position631, tokenIndex631, depth631 := position, tokenIndex, depth
			{
				position632 := position
				depth++
				{
					position633, tokenIndex633, depth633 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l634
					}
					goto l633
				l634:
					position, tokenIndex, depth = position633, tokenIndex633, depth633
					if !_rules[ruleMinus]() {
						goto l631
					}
				}
			l633:
				depth--
				add(rulePlusMinusOp, position632)
			}
			return true
		l631:
			position, tokenIndex, depth = position631, tokenIndex631, depth631
			return false
		},
		/* 57 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				{
					position637, tokenIndex637, depth637 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l638
					}
					goto l637
				l638:
					position, tokenIndex, depth = position637, tokenIndex637, depth637
					if !_rules[ruleDivide]() {
						goto l639
					}
					goto l637
				l639:
					position, tokenIndex, depth = position637, tokenIndex637, depth637
					if !_rules[ruleModulo]() {
						goto l635
					}
				}
			l637:
				depth--
				add(ruleMultDivOp, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 58 Stream <- <(<ident> Action42)> */
		func() bool {
			position640, tokenIndex640, depth640 := position, tokenIndex, depth
			{
				position641 := position
				depth++
				{
					position642 := position
					depth++
					if !_rules[ruleident]() {
						goto l640
					}
					depth--
					add(rulePegText, position642)
				}
				if !_rules[ruleAction42]() {
					goto l640
				}
				depth--
				add(ruleStream, position641)
			}
			return true
		l640:
			position, tokenIndex, depth = position640, tokenIndex640, depth640
			return false
		},
		/* 59 RowMeta <- <RowTimestamp> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l643
				}
				depth--
				add(ruleRowMeta, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 60 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action43)> */
		func() bool {
			position645, tokenIndex645, depth645 := position, tokenIndex, depth
			{
				position646 := position
				depth++
				{
					position647 := position
					depth++
					{
						position648, tokenIndex648, depth648 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l648
						}
						if buffer[position] != rune(':') {
							goto l648
						}
						position++
						goto l649
					l648:
						position, tokenIndex, depth = position648, tokenIndex648, depth648
					}
				l649:
					if buffer[position] != rune('t') {
						goto l645
					}
					position++
					if buffer[position] != rune('s') {
						goto l645
					}
					position++
					if buffer[position] != rune('(') {
						goto l645
					}
					position++
					if buffer[position] != rune(')') {
						goto l645
					}
					position++
					depth--
					add(rulePegText, position647)
				}
				if !_rules[ruleAction43]() {
					goto l645
				}
				depth--
				add(ruleRowTimestamp, position646)
			}
			return true
		l645:
			position, tokenIndex, depth = position645, tokenIndex645, depth645
			return false
		},
		/* 61 RowValue <- <(<((ident ':')? jsonPath)> Action44)> */
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
						if !_rules[ruleident]() {
							goto l653
						}
						if buffer[position] != rune(':') {
							goto l653
						}
						position++
						goto l654
					l653:
						position, tokenIndex, depth = position653, tokenIndex653, depth653
					}
				l654:
					if !_rules[rulejsonPath]() {
						goto l650
					}
					depth--
					add(rulePegText, position652)
				}
				if !_rules[ruleAction44]() {
					goto l650
				}
				depth--
				add(ruleRowValue, position651)
			}
			return true
		l650:
			position, tokenIndex, depth = position650, tokenIndex650, depth650
			return false
		},
		/* 62 NumericLiteral <- <(<('-'? [0-9]+)> Action45)> */
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
						if buffer[position] != rune('-') {
							goto l658
						}
						position++
						goto l659
					l658:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
					}
				l659:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l655
					}
					position++
				l660:
					{
						position661, tokenIndex661, depth661 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l661
						}
						position++
						goto l660
					l661:
						position, tokenIndex, depth = position661, tokenIndex661, depth661
					}
					depth--
					add(rulePegText, position657)
				}
				if !_rules[ruleAction45]() {
					goto l655
				}
				depth--
				add(ruleNumericLiteral, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 63 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action46)> */
		func() bool {
			position662, tokenIndex662, depth662 := position, tokenIndex, depth
			{
				position663 := position
				depth++
				{
					position664 := position
					depth++
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l665
						}
						position++
						goto l666
					l665:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
					}
				l666:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l662
					}
					position++
				l667:
					{
						position668, tokenIndex668, depth668 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position668, tokenIndex668, depth668
					}
					if buffer[position] != rune('.') {
						goto l662
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l662
					}
					position++
				l669:
					{
						position670, tokenIndex670, depth670 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l670
						}
						position++
						goto l669
					l670:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
					}
					depth--
					add(rulePegText, position664)
				}
				if !_rules[ruleAction46]() {
					goto l662
				}
				depth--
				add(ruleFloatLiteral, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 64 Function <- <(<ident> Action47)> */
		func() bool {
			position671, tokenIndex671, depth671 := position, tokenIndex, depth
			{
				position672 := position
				depth++
				{
					position673 := position
					depth++
					if !_rules[ruleident]() {
						goto l671
					}
					depth--
					add(rulePegText, position673)
				}
				if !_rules[ruleAction47]() {
					goto l671
				}
				depth--
				add(ruleFunction, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 65 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action48)> */
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
						if buffer[position] != rune('n') {
							goto l678
						}
						position++
						goto l677
					l678:
						position, tokenIndex, depth = position677, tokenIndex677, depth677
						if buffer[position] != rune('N') {
							goto l674
						}
						position++
					}
				l677:
					{
						position679, tokenIndex679, depth679 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l680
						}
						position++
						goto l679
					l680:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
						if buffer[position] != rune('U') {
							goto l674
						}
						position++
					}
				l679:
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if buffer[position] != rune('L') {
							goto l674
						}
						position++
					}
				l681:
					{
						position683, tokenIndex683, depth683 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l684
						}
						position++
						goto l683
					l684:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
						if buffer[position] != rune('L') {
							goto l674
						}
						position++
					}
				l683:
					depth--
					add(rulePegText, position676)
				}
				if !_rules[ruleAction48]() {
					goto l674
				}
				depth--
				add(ruleNullLiteral, position675)
			}
			return true
		l674:
			position, tokenIndex, depth = position674, tokenIndex674, depth674
			return false
		},
		/* 66 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position685, tokenIndex685, depth685 := position, tokenIndex, depth
			{
				position686 := position
				depth++
				{
					position687, tokenIndex687, depth687 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l688
					}
					goto l687
				l688:
					position, tokenIndex, depth = position687, tokenIndex687, depth687
					if !_rules[ruleFALSE]() {
						goto l685
					}
				}
			l687:
				depth--
				add(ruleBooleanLiteral, position686)
			}
			return true
		l685:
			position, tokenIndex, depth = position685, tokenIndex685, depth685
			return false
		},
		/* 67 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action49)> */
		func() bool {
			position689, tokenIndex689, depth689 := position, tokenIndex, depth
			{
				position690 := position
				depth++
				{
					position691 := position
					depth++
					{
						position692, tokenIndex692, depth692 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l693
						}
						position++
						goto l692
					l693:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						if buffer[position] != rune('T') {
							goto l689
						}
						position++
					}
				l692:
					{
						position694, tokenIndex694, depth694 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l695
						}
						position++
						goto l694
					l695:
						position, tokenIndex, depth = position694, tokenIndex694, depth694
						if buffer[position] != rune('R') {
							goto l689
						}
						position++
					}
				l694:
					{
						position696, tokenIndex696, depth696 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l697
						}
						position++
						goto l696
					l697:
						position, tokenIndex, depth = position696, tokenIndex696, depth696
						if buffer[position] != rune('U') {
							goto l689
						}
						position++
					}
				l696:
					{
						position698, tokenIndex698, depth698 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l699
						}
						position++
						goto l698
					l699:
						position, tokenIndex, depth = position698, tokenIndex698, depth698
						if buffer[position] != rune('E') {
							goto l689
						}
						position++
					}
				l698:
					depth--
					add(rulePegText, position691)
				}
				if !_rules[ruleAction49]() {
					goto l689
				}
				depth--
				add(ruleTRUE, position690)
			}
			return true
		l689:
			position, tokenIndex, depth = position689, tokenIndex689, depth689
			return false
		},
		/* 68 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action50)> */
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
						if buffer[position] != rune('f') {
							goto l704
						}
						position++
						goto l703
					l704:
						position, tokenIndex, depth = position703, tokenIndex703, depth703
						if buffer[position] != rune('F') {
							goto l700
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
							goto l700
						}
						position++
					}
				l705:
					{
						position707, tokenIndex707, depth707 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l708
						}
						position++
						goto l707
					l708:
						position, tokenIndex, depth = position707, tokenIndex707, depth707
						if buffer[position] != rune('L') {
							goto l700
						}
						position++
					}
				l707:
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l710
						}
						position++
						goto l709
					l710:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if buffer[position] != rune('S') {
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
					depth--
					add(rulePegText, position702)
				}
				if !_rules[ruleAction50]() {
					goto l700
				}
				depth--
				add(ruleFALSE, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 69 Wildcard <- <(<'*'> Action51)> */
		func() bool {
			position713, tokenIndex713, depth713 := position, tokenIndex, depth
			{
				position714 := position
				depth++
				{
					position715 := position
					depth++
					if buffer[position] != rune('*') {
						goto l713
					}
					position++
					depth--
					add(rulePegText, position715)
				}
				if !_rules[ruleAction51]() {
					goto l713
				}
				depth--
				add(ruleWildcard, position714)
			}
			return true
		l713:
			position, tokenIndex, depth = position713, tokenIndex713, depth713
			return false
		},
		/* 70 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action52)> */
		func() bool {
			position716, tokenIndex716, depth716 := position, tokenIndex, depth
			{
				position717 := position
				depth++
				{
					position718 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l716
					}
					position++
				l719:
					{
						position720, tokenIndex720, depth720 := position, tokenIndex, depth
						{
							position721, tokenIndex721, depth721 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l722
							}
							position++
							if buffer[position] != rune('\'') {
								goto l722
							}
							position++
							goto l721
						l722:
							position, tokenIndex, depth = position721, tokenIndex721, depth721
							{
								position723, tokenIndex723, depth723 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l723
								}
								position++
								goto l720
							l723:
								position, tokenIndex, depth = position723, tokenIndex723, depth723
							}
							if !matchDot() {
								goto l720
							}
						}
					l721:
						goto l719
					l720:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
					}
					if buffer[position] != rune('\'') {
						goto l716
					}
					position++
					depth--
					add(rulePegText, position718)
				}
				if !_rules[ruleAction52]() {
					goto l716
				}
				depth--
				add(ruleStringLiteral, position717)
			}
			return true
		l716:
			position, tokenIndex, depth = position716, tokenIndex716, depth716
			return false
		},
		/* 71 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action53)> */
		func() bool {
			position724, tokenIndex724, depth724 := position, tokenIndex, depth
			{
				position725 := position
				depth++
				{
					position726 := position
					depth++
					{
						position727, tokenIndex727, depth727 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l728
						}
						position++
						goto l727
					l728:
						position, tokenIndex, depth = position727, tokenIndex727, depth727
						if buffer[position] != rune('I') {
							goto l724
						}
						position++
					}
				l727:
					{
						position729, tokenIndex729, depth729 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l730
						}
						position++
						goto l729
					l730:
						position, tokenIndex, depth = position729, tokenIndex729, depth729
						if buffer[position] != rune('S') {
							goto l724
						}
						position++
					}
				l729:
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l732
						}
						position++
						goto l731
					l732:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
						if buffer[position] != rune('T') {
							goto l724
						}
						position++
					}
				l731:
					{
						position733, tokenIndex733, depth733 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l734
						}
						position++
						goto l733
					l734:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
						if buffer[position] != rune('R') {
							goto l724
						}
						position++
					}
				l733:
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l736
						}
						position++
						goto l735
					l736:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if buffer[position] != rune('E') {
							goto l724
						}
						position++
					}
				l735:
					{
						position737, tokenIndex737, depth737 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l738
						}
						position++
						goto l737
					l738:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
						if buffer[position] != rune('A') {
							goto l724
						}
						position++
					}
				l737:
					{
						position739, tokenIndex739, depth739 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l740
						}
						position++
						goto l739
					l740:
						position, tokenIndex, depth = position739, tokenIndex739, depth739
						if buffer[position] != rune('M') {
							goto l724
						}
						position++
					}
				l739:
					depth--
					add(rulePegText, position726)
				}
				if !_rules[ruleAction53]() {
					goto l724
				}
				depth--
				add(ruleISTREAM, position725)
			}
			return true
		l724:
			position, tokenIndex, depth = position724, tokenIndex724, depth724
			return false
		},
		/* 72 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action54)> */
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
						if buffer[position] != rune('d') {
							goto l745
						}
						position++
						goto l744
					l745:
						position, tokenIndex, depth = position744, tokenIndex744, depth744
						if buffer[position] != rune('D') {
							goto l741
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
							goto l741
						}
						position++
					}
				l746:
					{
						position748, tokenIndex748, depth748 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l749
						}
						position++
						goto l748
					l749:
						position, tokenIndex, depth = position748, tokenIndex748, depth748
						if buffer[position] != rune('T') {
							goto l741
						}
						position++
					}
				l748:
					{
						position750, tokenIndex750, depth750 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l751
						}
						position++
						goto l750
					l751:
						position, tokenIndex, depth = position750, tokenIndex750, depth750
						if buffer[position] != rune('R') {
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
						if buffer[position] != rune('a') {
							goto l755
						}
						position++
						goto l754
					l755:
						position, tokenIndex, depth = position754, tokenIndex754, depth754
						if buffer[position] != rune('A') {
							goto l741
						}
						position++
					}
				l754:
					{
						position756, tokenIndex756, depth756 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l757
						}
						position++
						goto l756
					l757:
						position, tokenIndex, depth = position756, tokenIndex756, depth756
						if buffer[position] != rune('M') {
							goto l741
						}
						position++
					}
				l756:
					depth--
					add(rulePegText, position743)
				}
				if !_rules[ruleAction54]() {
					goto l741
				}
				depth--
				add(ruleDSTREAM, position742)
			}
			return true
		l741:
			position, tokenIndex, depth = position741, tokenIndex741, depth741
			return false
		},
		/* 73 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action55)> */
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
						if buffer[position] != rune('r') {
							goto l762
						}
						position++
						goto l761
					l762:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if buffer[position] != rune('R') {
							goto l758
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
							goto l758
						}
						position++
					}
				l763:
					{
						position765, tokenIndex765, depth765 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l766
						}
						position++
						goto l765
					l766:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
						if buffer[position] != rune('T') {
							goto l758
						}
						position++
					}
				l765:
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('R') {
							goto l758
						}
						position++
					}
				l767:
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('E') {
							goto l758
						}
						position++
					}
				l769:
					{
						position771, tokenIndex771, depth771 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l772
						}
						position++
						goto l771
					l772:
						position, tokenIndex, depth = position771, tokenIndex771, depth771
						if buffer[position] != rune('A') {
							goto l758
						}
						position++
					}
				l771:
					{
						position773, tokenIndex773, depth773 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l774
						}
						position++
						goto l773
					l774:
						position, tokenIndex, depth = position773, tokenIndex773, depth773
						if buffer[position] != rune('M') {
							goto l758
						}
						position++
					}
				l773:
					depth--
					add(rulePegText, position760)
				}
				if !_rules[ruleAction55]() {
					goto l758
				}
				depth--
				add(ruleRSTREAM, position759)
			}
			return true
		l758:
			position, tokenIndex, depth = position758, tokenIndex758, depth758
			return false
		},
		/* 74 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action56)> */
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
						if buffer[position] != rune('t') {
							goto l779
						}
						position++
						goto l778
					l779:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						if buffer[position] != rune('T') {
							goto l775
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
							goto l775
						}
						position++
					}
				l780:
					{
						position782, tokenIndex782, depth782 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l783
						}
						position++
						goto l782
					l783:
						position, tokenIndex, depth = position782, tokenIndex782, depth782
						if buffer[position] != rune('P') {
							goto l775
						}
						position++
					}
				l782:
					{
						position784, tokenIndex784, depth784 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l785
						}
						position++
						goto l784
					l785:
						position, tokenIndex, depth = position784, tokenIndex784, depth784
						if buffer[position] != rune('L') {
							goto l775
						}
						position++
					}
				l784:
					{
						position786, tokenIndex786, depth786 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l787
						}
						position++
						goto l786
					l787:
						position, tokenIndex, depth = position786, tokenIndex786, depth786
						if buffer[position] != rune('E') {
							goto l775
						}
						position++
					}
				l786:
					{
						position788, tokenIndex788, depth788 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l789
						}
						position++
						goto l788
					l789:
						position, tokenIndex, depth = position788, tokenIndex788, depth788
						if buffer[position] != rune('S') {
							goto l775
						}
						position++
					}
				l788:
					depth--
					add(rulePegText, position777)
				}
				if !_rules[ruleAction56]() {
					goto l775
				}
				depth--
				add(ruleTUPLES, position776)
			}
			return true
		l775:
			position, tokenIndex, depth = position775, tokenIndex775, depth775
			return false
		},
		/* 75 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action57)> */
		func() bool {
			position790, tokenIndex790, depth790 := position, tokenIndex, depth
			{
				position791 := position
				depth++
				{
					position792 := position
					depth++
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
							goto l790
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
							goto l790
						}
						position++
					}
				l795:
					{
						position797, tokenIndex797, depth797 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l798
						}
						position++
						goto l797
					l798:
						position, tokenIndex, depth = position797, tokenIndex797, depth797
						if buffer[position] != rune('C') {
							goto l790
						}
						position++
					}
				l797:
					{
						position799, tokenIndex799, depth799 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l800
						}
						position++
						goto l799
					l800:
						position, tokenIndex, depth = position799, tokenIndex799, depth799
						if buffer[position] != rune('O') {
							goto l790
						}
						position++
					}
				l799:
					{
						position801, tokenIndex801, depth801 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l802
						}
						position++
						goto l801
					l802:
						position, tokenIndex, depth = position801, tokenIndex801, depth801
						if buffer[position] != rune('N') {
							goto l790
						}
						position++
					}
				l801:
					{
						position803, tokenIndex803, depth803 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l804
						}
						position++
						goto l803
					l804:
						position, tokenIndex, depth = position803, tokenIndex803, depth803
						if buffer[position] != rune('D') {
							goto l790
						}
						position++
					}
				l803:
					{
						position805, tokenIndex805, depth805 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l806
						}
						position++
						goto l805
					l806:
						position, tokenIndex, depth = position805, tokenIndex805, depth805
						if buffer[position] != rune('S') {
							goto l790
						}
						position++
					}
				l805:
					depth--
					add(rulePegText, position792)
				}
				if !_rules[ruleAction57]() {
					goto l790
				}
				depth--
				add(ruleSECONDS, position791)
			}
			return true
		l790:
			position, tokenIndex, depth = position790, tokenIndex790, depth790
			return false
		},
		/* 76 StreamIdentifier <- <(<ident> Action58)> */
		func() bool {
			position807, tokenIndex807, depth807 := position, tokenIndex, depth
			{
				position808 := position
				depth++
				{
					position809 := position
					depth++
					if !_rules[ruleident]() {
						goto l807
					}
					depth--
					add(rulePegText, position809)
				}
				if !_rules[ruleAction58]() {
					goto l807
				}
				depth--
				add(ruleStreamIdentifier, position808)
			}
			return true
		l807:
			position, tokenIndex, depth = position807, tokenIndex807, depth807
			return false
		},
		/* 77 SourceSinkType <- <(<ident> Action59)> */
		func() bool {
			position810, tokenIndex810, depth810 := position, tokenIndex, depth
			{
				position811 := position
				depth++
				{
					position812 := position
					depth++
					if !_rules[ruleident]() {
						goto l810
					}
					depth--
					add(rulePegText, position812)
				}
				if !_rules[ruleAction59]() {
					goto l810
				}
				depth--
				add(ruleSourceSinkType, position811)
			}
			return true
		l810:
			position, tokenIndex, depth = position810, tokenIndex810, depth810
			return false
		},
		/* 78 SourceSinkParamKey <- <(<ident> Action60)> */
		func() bool {
			position813, tokenIndex813, depth813 := position, tokenIndex, depth
			{
				position814 := position
				depth++
				{
					position815 := position
					depth++
					if !_rules[ruleident]() {
						goto l813
					}
					depth--
					add(rulePegText, position815)
				}
				if !_rules[ruleAction60]() {
					goto l813
				}
				depth--
				add(ruleSourceSinkParamKey, position814)
			}
			return true
		l813:
			position, tokenIndex, depth = position813, tokenIndex813, depth813
			return false
		},
		/* 79 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action61)> */
		func() bool {
			position816, tokenIndex816, depth816 := position, tokenIndex, depth
			{
				position817 := position
				depth++
				{
					position818 := position
					depth++
					{
						position819, tokenIndex819, depth819 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l820
						}
						position++
						goto l819
					l820:
						position, tokenIndex, depth = position819, tokenIndex819, depth819
						if buffer[position] != rune('P') {
							goto l816
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
							goto l816
						}
						position++
					}
				l821:
					{
						position823, tokenIndex823, depth823 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l824
						}
						position++
						goto l823
					l824:
						position, tokenIndex, depth = position823, tokenIndex823, depth823
						if buffer[position] != rune('U') {
							goto l816
						}
						position++
					}
				l823:
					{
						position825, tokenIndex825, depth825 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l826
						}
						position++
						goto l825
					l826:
						position, tokenIndex, depth = position825, tokenIndex825, depth825
						if buffer[position] != rune('S') {
							goto l816
						}
						position++
					}
				l825:
					{
						position827, tokenIndex827, depth827 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l828
						}
						position++
						goto l827
					l828:
						position, tokenIndex, depth = position827, tokenIndex827, depth827
						if buffer[position] != rune('E') {
							goto l816
						}
						position++
					}
				l827:
					{
						position829, tokenIndex829, depth829 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l830
						}
						position++
						goto l829
					l830:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						if buffer[position] != rune('D') {
							goto l816
						}
						position++
					}
				l829:
					depth--
					add(rulePegText, position818)
				}
				if !_rules[ruleAction61]() {
					goto l816
				}
				depth--
				add(rulePaused, position817)
			}
			return true
		l816:
			position, tokenIndex, depth = position816, tokenIndex816, depth816
			return false
		},
		/* 80 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action62)> */
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
						if buffer[position] != rune('u') {
							goto l835
						}
						position++
						goto l834
					l835:
						position, tokenIndex, depth = position834, tokenIndex834, depth834
						if buffer[position] != rune('U') {
							goto l831
						}
						position++
					}
				l834:
					{
						position836, tokenIndex836, depth836 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l837
						}
						position++
						goto l836
					l837:
						position, tokenIndex, depth = position836, tokenIndex836, depth836
						if buffer[position] != rune('N') {
							goto l831
						}
						position++
					}
				l836:
					{
						position838, tokenIndex838, depth838 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l839
						}
						position++
						goto l838
					l839:
						position, tokenIndex, depth = position838, tokenIndex838, depth838
						if buffer[position] != rune('P') {
							goto l831
						}
						position++
					}
				l838:
					{
						position840, tokenIndex840, depth840 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l841
						}
						position++
						goto l840
					l841:
						position, tokenIndex, depth = position840, tokenIndex840, depth840
						if buffer[position] != rune('A') {
							goto l831
						}
						position++
					}
				l840:
					{
						position842, tokenIndex842, depth842 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l843
						}
						position++
						goto l842
					l843:
						position, tokenIndex, depth = position842, tokenIndex842, depth842
						if buffer[position] != rune('U') {
							goto l831
						}
						position++
					}
				l842:
					{
						position844, tokenIndex844, depth844 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l845
						}
						position++
						goto l844
					l845:
						position, tokenIndex, depth = position844, tokenIndex844, depth844
						if buffer[position] != rune('S') {
							goto l831
						}
						position++
					}
				l844:
					{
						position846, tokenIndex846, depth846 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l847
						}
						position++
						goto l846
					l847:
						position, tokenIndex, depth = position846, tokenIndex846, depth846
						if buffer[position] != rune('E') {
							goto l831
						}
						position++
					}
				l846:
					{
						position848, tokenIndex848, depth848 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l849
						}
						position++
						goto l848
					l849:
						position, tokenIndex, depth = position848, tokenIndex848, depth848
						if buffer[position] != rune('D') {
							goto l831
						}
						position++
					}
				l848:
					depth--
					add(rulePegText, position833)
				}
				if !_rules[ruleAction62]() {
					goto l831
				}
				depth--
				add(ruleUnpaused, position832)
			}
			return true
		l831:
			position, tokenIndex, depth = position831, tokenIndex831, depth831
			return false
		},
		/* 81 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action63)> */
		func() bool {
			position850, tokenIndex850, depth850 := position, tokenIndex, depth
			{
				position851 := position
				depth++
				{
					position852 := position
					depth++
					{
						position853, tokenIndex853, depth853 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l854
						}
						position++
						goto l853
					l854:
						position, tokenIndex, depth = position853, tokenIndex853, depth853
						if buffer[position] != rune('O') {
							goto l850
						}
						position++
					}
				l853:
					{
						position855, tokenIndex855, depth855 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l856
						}
						position++
						goto l855
					l856:
						position, tokenIndex, depth = position855, tokenIndex855, depth855
						if buffer[position] != rune('R') {
							goto l850
						}
						position++
					}
				l855:
					depth--
					add(rulePegText, position852)
				}
				if !_rules[ruleAction63]() {
					goto l850
				}
				depth--
				add(ruleOr, position851)
			}
			return true
		l850:
			position, tokenIndex, depth = position850, tokenIndex850, depth850
			return false
		},
		/* 82 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action64)> */
		func() bool {
			position857, tokenIndex857, depth857 := position, tokenIndex, depth
			{
				position858 := position
				depth++
				{
					position859 := position
					depth++
					{
						position860, tokenIndex860, depth860 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l861
						}
						position++
						goto l860
					l861:
						position, tokenIndex, depth = position860, tokenIndex860, depth860
						if buffer[position] != rune('A') {
							goto l857
						}
						position++
					}
				l860:
					{
						position862, tokenIndex862, depth862 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l863
						}
						position++
						goto l862
					l863:
						position, tokenIndex, depth = position862, tokenIndex862, depth862
						if buffer[position] != rune('N') {
							goto l857
						}
						position++
					}
				l862:
					{
						position864, tokenIndex864, depth864 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l865
						}
						position++
						goto l864
					l865:
						position, tokenIndex, depth = position864, tokenIndex864, depth864
						if buffer[position] != rune('D') {
							goto l857
						}
						position++
					}
				l864:
					depth--
					add(rulePegText, position859)
				}
				if !_rules[ruleAction64]() {
					goto l857
				}
				depth--
				add(ruleAnd, position858)
			}
			return true
		l857:
			position, tokenIndex, depth = position857, tokenIndex857, depth857
			return false
		},
		/* 83 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action65)> */
		func() bool {
			position866, tokenIndex866, depth866 := position, tokenIndex, depth
			{
				position867 := position
				depth++
				{
					position868 := position
					depth++
					{
						position869, tokenIndex869, depth869 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l870
						}
						position++
						goto l869
					l870:
						position, tokenIndex, depth = position869, tokenIndex869, depth869
						if buffer[position] != rune('N') {
							goto l866
						}
						position++
					}
				l869:
					{
						position871, tokenIndex871, depth871 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l872
						}
						position++
						goto l871
					l872:
						position, tokenIndex, depth = position871, tokenIndex871, depth871
						if buffer[position] != rune('O') {
							goto l866
						}
						position++
					}
				l871:
					{
						position873, tokenIndex873, depth873 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l874
						}
						position++
						goto l873
					l874:
						position, tokenIndex, depth = position873, tokenIndex873, depth873
						if buffer[position] != rune('T') {
							goto l866
						}
						position++
					}
				l873:
					depth--
					add(rulePegText, position868)
				}
				if !_rules[ruleAction65]() {
					goto l866
				}
				depth--
				add(ruleNot, position867)
			}
			return true
		l866:
			position, tokenIndex, depth = position866, tokenIndex866, depth866
			return false
		},
		/* 84 Equal <- <(<'='> Action66)> */
		func() bool {
			position875, tokenIndex875, depth875 := position, tokenIndex, depth
			{
				position876 := position
				depth++
				{
					position877 := position
					depth++
					if buffer[position] != rune('=') {
						goto l875
					}
					position++
					depth--
					add(rulePegText, position877)
				}
				if !_rules[ruleAction66]() {
					goto l875
				}
				depth--
				add(ruleEqual, position876)
			}
			return true
		l875:
			position, tokenIndex, depth = position875, tokenIndex875, depth875
			return false
		},
		/* 85 Less <- <(<'<'> Action67)> */
		func() bool {
			position878, tokenIndex878, depth878 := position, tokenIndex, depth
			{
				position879 := position
				depth++
				{
					position880 := position
					depth++
					if buffer[position] != rune('<') {
						goto l878
					}
					position++
					depth--
					add(rulePegText, position880)
				}
				if !_rules[ruleAction67]() {
					goto l878
				}
				depth--
				add(ruleLess, position879)
			}
			return true
		l878:
			position, tokenIndex, depth = position878, tokenIndex878, depth878
			return false
		},
		/* 86 LessOrEqual <- <(<('<' '=')> Action68)> */
		func() bool {
			position881, tokenIndex881, depth881 := position, tokenIndex, depth
			{
				position882 := position
				depth++
				{
					position883 := position
					depth++
					if buffer[position] != rune('<') {
						goto l881
					}
					position++
					if buffer[position] != rune('=') {
						goto l881
					}
					position++
					depth--
					add(rulePegText, position883)
				}
				if !_rules[ruleAction68]() {
					goto l881
				}
				depth--
				add(ruleLessOrEqual, position882)
			}
			return true
		l881:
			position, tokenIndex, depth = position881, tokenIndex881, depth881
			return false
		},
		/* 87 Greater <- <(<'>'> Action69)> */
		func() bool {
			position884, tokenIndex884, depth884 := position, tokenIndex, depth
			{
				position885 := position
				depth++
				{
					position886 := position
					depth++
					if buffer[position] != rune('>') {
						goto l884
					}
					position++
					depth--
					add(rulePegText, position886)
				}
				if !_rules[ruleAction69]() {
					goto l884
				}
				depth--
				add(ruleGreater, position885)
			}
			return true
		l884:
			position, tokenIndex, depth = position884, tokenIndex884, depth884
			return false
		},
		/* 88 GreaterOrEqual <- <(<('>' '=')> Action70)> */
		func() bool {
			position887, tokenIndex887, depth887 := position, tokenIndex, depth
			{
				position888 := position
				depth++
				{
					position889 := position
					depth++
					if buffer[position] != rune('>') {
						goto l887
					}
					position++
					if buffer[position] != rune('=') {
						goto l887
					}
					position++
					depth--
					add(rulePegText, position889)
				}
				if !_rules[ruleAction70]() {
					goto l887
				}
				depth--
				add(ruleGreaterOrEqual, position888)
			}
			return true
		l887:
			position, tokenIndex, depth = position887, tokenIndex887, depth887
			return false
		},
		/* 89 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action71)> */
		func() bool {
			position890, tokenIndex890, depth890 := position, tokenIndex, depth
			{
				position891 := position
				depth++
				{
					position892 := position
					depth++
					{
						position893, tokenIndex893, depth893 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l894
						}
						position++
						if buffer[position] != rune('=') {
							goto l894
						}
						position++
						goto l893
					l894:
						position, tokenIndex, depth = position893, tokenIndex893, depth893
						if buffer[position] != rune('<') {
							goto l890
						}
						position++
						if buffer[position] != rune('>') {
							goto l890
						}
						position++
					}
				l893:
					depth--
					add(rulePegText, position892)
				}
				if !_rules[ruleAction71]() {
					goto l890
				}
				depth--
				add(ruleNotEqual, position891)
			}
			return true
		l890:
			position, tokenIndex, depth = position890, tokenIndex890, depth890
			return false
		},
		/* 90 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action72)> */
		func() bool {
			position895, tokenIndex895, depth895 := position, tokenIndex, depth
			{
				position896 := position
				depth++
				{
					position897 := position
					depth++
					{
						position898, tokenIndex898, depth898 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l899
						}
						position++
						goto l898
					l899:
						position, tokenIndex, depth = position898, tokenIndex898, depth898
						if buffer[position] != rune('I') {
							goto l895
						}
						position++
					}
				l898:
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
							goto l895
						}
						position++
					}
				l900:
					depth--
					add(rulePegText, position897)
				}
				if !_rules[ruleAction72]() {
					goto l895
				}
				depth--
				add(ruleIs, position896)
			}
			return true
		l895:
			position, tokenIndex, depth = position895, tokenIndex895, depth895
			return false
		},
		/* 91 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action73)> */
		func() bool {
			position902, tokenIndex902, depth902 := position, tokenIndex, depth
			{
				position903 := position
				depth++
				{
					position904 := position
					depth++
					{
						position905, tokenIndex905, depth905 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l906
						}
						position++
						goto l905
					l906:
						position, tokenIndex, depth = position905, tokenIndex905, depth905
						if buffer[position] != rune('I') {
							goto l902
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
							goto l902
						}
						position++
					}
				l907:
					if !_rules[rulesp]() {
						goto l902
					}
					{
						position909, tokenIndex909, depth909 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l910
						}
						position++
						goto l909
					l910:
						position, tokenIndex, depth = position909, tokenIndex909, depth909
						if buffer[position] != rune('N') {
							goto l902
						}
						position++
					}
				l909:
					{
						position911, tokenIndex911, depth911 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l912
						}
						position++
						goto l911
					l912:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('O') {
							goto l902
						}
						position++
					}
				l911:
					{
						position913, tokenIndex913, depth913 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l914
						}
						position++
						goto l913
					l914:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('T') {
							goto l902
						}
						position++
					}
				l913:
					depth--
					add(rulePegText, position904)
				}
				if !_rules[ruleAction73]() {
					goto l902
				}
				depth--
				add(ruleIsNot, position903)
			}
			return true
		l902:
			position, tokenIndex, depth = position902, tokenIndex902, depth902
			return false
		},
		/* 92 Plus <- <(<'+'> Action74)> */
		func() bool {
			position915, tokenIndex915, depth915 := position, tokenIndex, depth
			{
				position916 := position
				depth++
				{
					position917 := position
					depth++
					if buffer[position] != rune('+') {
						goto l915
					}
					position++
					depth--
					add(rulePegText, position917)
				}
				if !_rules[ruleAction74]() {
					goto l915
				}
				depth--
				add(rulePlus, position916)
			}
			return true
		l915:
			position, tokenIndex, depth = position915, tokenIndex915, depth915
			return false
		},
		/* 93 Minus <- <(<'-'> Action75)> */
		func() bool {
			position918, tokenIndex918, depth918 := position, tokenIndex, depth
			{
				position919 := position
				depth++
				{
					position920 := position
					depth++
					if buffer[position] != rune('-') {
						goto l918
					}
					position++
					depth--
					add(rulePegText, position920)
				}
				if !_rules[ruleAction75]() {
					goto l918
				}
				depth--
				add(ruleMinus, position919)
			}
			return true
		l918:
			position, tokenIndex, depth = position918, tokenIndex918, depth918
			return false
		},
		/* 94 Multiply <- <(<'*'> Action76)> */
		func() bool {
			position921, tokenIndex921, depth921 := position, tokenIndex, depth
			{
				position922 := position
				depth++
				{
					position923 := position
					depth++
					if buffer[position] != rune('*') {
						goto l921
					}
					position++
					depth--
					add(rulePegText, position923)
				}
				if !_rules[ruleAction76]() {
					goto l921
				}
				depth--
				add(ruleMultiply, position922)
			}
			return true
		l921:
			position, tokenIndex, depth = position921, tokenIndex921, depth921
			return false
		},
		/* 95 Divide <- <(<'/'> Action77)> */
		func() bool {
			position924, tokenIndex924, depth924 := position, tokenIndex, depth
			{
				position925 := position
				depth++
				{
					position926 := position
					depth++
					if buffer[position] != rune('/') {
						goto l924
					}
					position++
					depth--
					add(rulePegText, position926)
				}
				if !_rules[ruleAction77]() {
					goto l924
				}
				depth--
				add(ruleDivide, position925)
			}
			return true
		l924:
			position, tokenIndex, depth = position924, tokenIndex924, depth924
			return false
		},
		/* 96 Modulo <- <(<'%'> Action78)> */
		func() bool {
			position927, tokenIndex927, depth927 := position, tokenIndex, depth
			{
				position928 := position
				depth++
				{
					position929 := position
					depth++
					if buffer[position] != rune('%') {
						goto l927
					}
					position++
					depth--
					add(rulePegText, position929)
				}
				if !_rules[ruleAction78]() {
					goto l927
				}
				depth--
				add(ruleModulo, position928)
			}
			return true
		l927:
			position, tokenIndex, depth = position927, tokenIndex927, depth927
			return false
		},
		/* 97 UnaryMinus <- <(<'-'> Action79)> */
		func() bool {
			position930, tokenIndex930, depth930 := position, tokenIndex, depth
			{
				position931 := position
				depth++
				{
					position932 := position
					depth++
					if buffer[position] != rune('-') {
						goto l930
					}
					position++
					depth--
					add(rulePegText, position932)
				}
				if !_rules[ruleAction79]() {
					goto l930
				}
				depth--
				add(ruleUnaryMinus, position931)
			}
			return true
		l930:
			position, tokenIndex, depth = position930, tokenIndex930, depth930
			return false
		},
		/* 98 Identifier <- <(<ident> Action80)> */
		func() bool {
			position933, tokenIndex933, depth933 := position, tokenIndex, depth
			{
				position934 := position
				depth++
				{
					position935 := position
					depth++
					if !_rules[ruleident]() {
						goto l933
					}
					depth--
					add(rulePegText, position935)
				}
				if !_rules[ruleAction80]() {
					goto l933
				}
				depth--
				add(ruleIdentifier, position934)
			}
			return true
		l933:
			position, tokenIndex, depth = position933, tokenIndex933, depth933
			return false
		},
		/* 99 TargetIdentifier <- <(<jsonPath> Action81)> */
		func() bool {
			position936, tokenIndex936, depth936 := position, tokenIndex, depth
			{
				position937 := position
				depth++
				{
					position938 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l936
					}
					depth--
					add(rulePegText, position938)
				}
				if !_rules[ruleAction81]() {
					goto l936
				}
				depth--
				add(ruleTargetIdentifier, position937)
			}
			return true
		l936:
			position, tokenIndex, depth = position936, tokenIndex936, depth936
			return false
		},
		/* 100 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position939, tokenIndex939, depth939 := position, tokenIndex, depth
			{
				position940 := position
				depth++
				{
					position941, tokenIndex941, depth941 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l942
					}
					position++
					goto l941
				l942:
					position, tokenIndex, depth = position941, tokenIndex941, depth941
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l939
					}
					position++
				}
			l941:
			l943:
				{
					position944, tokenIndex944, depth944 := position, tokenIndex, depth
					{
						position945, tokenIndex945, depth945 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l946
						}
						position++
						goto l945
					l946:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l947
						}
						position++
						goto l945
					l947:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l948
						}
						position++
						goto l945
					l948:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
						if buffer[position] != rune('_') {
							goto l944
						}
						position++
					}
				l945:
					goto l943
				l944:
					position, tokenIndex, depth = position944, tokenIndex944, depth944
				}
				depth--
				add(ruleident, position940)
			}
			return true
		l939:
			position, tokenIndex, depth = position939, tokenIndex939, depth939
			return false
		},
		/* 101 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position949, tokenIndex949, depth949 := position, tokenIndex, depth
			{
				position950 := position
				depth++
				{
					position951, tokenIndex951, depth951 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l952
					}
					position++
					goto l951
				l952:
					position, tokenIndex, depth = position951, tokenIndex951, depth951
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l949
					}
					position++
				}
			l951:
			l953:
				{
					position954, tokenIndex954, depth954 := position, tokenIndex, depth
					{
						position955, tokenIndex955, depth955 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l956
						}
						position++
						goto l955
					l956:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l957
						}
						position++
						goto l955
					l957:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l958
						}
						position++
						goto l955
					l958:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('_') {
							goto l959
						}
						position++
						goto l955
					l959:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('.') {
							goto l960
						}
						position++
						goto l955
					l960:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('[') {
							goto l961
						}
						position++
						goto l955
					l961:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune(']') {
							goto l962
						}
						position++
						goto l955
					l962:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('"') {
							goto l954
						}
						position++
					}
				l955:
					goto l953
				l954:
					position, tokenIndex, depth = position954, tokenIndex954, depth954
				}
				depth--
				add(rulejsonPath, position950)
			}
			return true
		l949:
			position, tokenIndex, depth = position949, tokenIndex949, depth949
			return false
		},
		/* 102 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position964 := position
				depth++
			l965:
				{
					position966, tokenIndex966, depth966 := position, tokenIndex, depth
					{
						position967, tokenIndex967, depth967 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l968
						}
						position++
						goto l967
					l968:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if buffer[position] != rune('\t') {
							goto l969
						}
						position++
						goto l967
					l969:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if buffer[position] != rune('\n') {
							goto l970
						}
						position++
						goto l967
					l970:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if buffer[position] != rune('\r') {
							goto l971
						}
						position++
						goto l967
					l971:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if !_rules[rulecomment]() {
							goto l966
						}
					}
				l967:
					goto l965
				l966:
					position, tokenIndex, depth = position966, tokenIndex966, depth966
				}
				depth--
				add(rulesp, position964)
			}
			return true
		},
		/* 103 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position972, tokenIndex972, depth972 := position, tokenIndex, depth
			{
				position973 := position
				depth++
				if buffer[position] != rune('-') {
					goto l972
				}
				position++
				if buffer[position] != rune('-') {
					goto l972
				}
				position++
			l974:
				{
					position975, tokenIndex975, depth975 := position, tokenIndex, depth
					{
						position976, tokenIndex976, depth976 := position, tokenIndex, depth
						{
							position977, tokenIndex977, depth977 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l978
							}
							position++
							goto l977
						l978:
							position, tokenIndex, depth = position977, tokenIndex977, depth977
							if buffer[position] != rune('\n') {
								goto l976
							}
							position++
						}
					l977:
						goto l975
					l976:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
					}
					if !matchDot() {
						goto l975
					}
					goto l974
				l975:
					position, tokenIndex, depth = position975, tokenIndex975, depth975
				}
				{
					position979, tokenIndex979, depth979 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l980
					}
					position++
					goto l979
				l980:
					position, tokenIndex, depth = position979, tokenIndex979, depth979
					if buffer[position] != rune('\n') {
						goto l972
					}
					position++
				}
			l979:
				depth--
				add(rulecomment, position973)
			}
			return true
		l972:
			position, tokenIndex, depth = position972, tokenIndex972, depth972
			return false
		},
		/* 105 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 106 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 107 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 108 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 109 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 110 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 111 Action6 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 112 Action7 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 113 Action8 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 114 Action9 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 115 Action10 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 116 Action11 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 117 Action12 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		nil,
		/* 119 Action13 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 120 Action14 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 121 Action15 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 122 Action16 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 123 Action17 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 124 Action18 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 125 Action19 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 126 Action20 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 127 Action21 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 128 Action22 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 129 Action23 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 130 Action24 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 131 Action25 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 132 Action26 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 133 Action27 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 134 Action28 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 135 Action29 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 136 Action30 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 137 Action31 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 138 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 139 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 140 Action34 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 141 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 142 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 143 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 144 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 145 Action39 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 146 Action40 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 147 Action41 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 148 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 149 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 150 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 151 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 152 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 153 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 154 Action48 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 155 Action49 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 156 Action50 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 157 Action51 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 158 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 159 Action53 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 160 Action54 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 161 Action55 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 162 Action56 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 163 Action57 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 164 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 165 Action59 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 166 Action60 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 167 Action61 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 168 Action62 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 169 Action63 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 170 Action64 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 171 Action65 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 172 Action66 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 173 Action67 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 174 Action68 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 175 Action69 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 176 Action70 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 177 Action71 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 178 Action72 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 179 Action73 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 180 Action74 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 181 Action75 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 182 Action76 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 183 Action77 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 184 Action78 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 185 Action79 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 186 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 187 Action81 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
	}
	p.rules = _rules
}

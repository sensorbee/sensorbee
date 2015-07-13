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
	rulePegText
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
	ruleAction80

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
	"PegText",
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
	"Action80",

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
	rules  [186]func() bool
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

			p.AssembleEmitter(begin, end)

		case ruleAction13:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction14:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction15:

			p.AssembleStreamEmitInterval()

		case ruleAction16:

			p.AssembleProjections(begin, end)

		case ruleAction17:

			p.AssembleAlias()

		case ruleAction18:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction19:

			p.AssembleInterval()

		case ruleAction20:

			p.AssembleInterval()

		case ruleAction21:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction22:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction23:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction24:

			p.EnsureAliasedStreamWindow()

		case ruleAction25:

			p.AssembleAliasedStreamWindow()

		case ruleAction26:

			p.AssembleStreamWindow()

		case ruleAction27:

			p.AssembleUDSFFuncApp()

		case ruleAction28:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction29:

			p.AssembleSourceSinkParam()

		case ruleAction30:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction31:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction32:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction33:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction34:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction38:

			p.AssembleUnaryPrefixOperation(begin, end)

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

			p.PushComponent(begin, end, Not)

		case ruleAction65:

			p.PushComponent(begin, end, Equal)

		case ruleAction66:

			p.PushComponent(begin, end, Less)

		case ruleAction67:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction68:

			p.PushComponent(begin, end, Greater)

		case ruleAction69:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction70:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction71:

			p.PushComponent(begin, end, Is)

		case ruleAction72:

			p.PushComponent(begin, end, IsNot)

		case ruleAction73:

			p.PushComponent(begin, end, Plus)

		case ruleAction74:

			p.PushComponent(begin, end, Minus)

		case ruleAction75:

			p.PushComponent(begin, end, Multiply)

		case ruleAction76:

			p.PushComponent(begin, end, Divide)

		case ruleAction77:

			p.PushComponent(begin, end, Modulo)

		case ruleAction78:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction80:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt / DropSourceStmt / DropStreamStmt)> */
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
			position21, tokenIndex21, depth21 := position, tokenIndex, depth
			{
				position22 := position
				depth++
				{
					position23, tokenIndex23, depth23 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l24
					}
					position++
					goto l23
				l24:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if buffer[position] != rune('S') {
						goto l21
					}
					position++
				}
			l23:
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l26
					}
					position++
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if buffer[position] != rune('E') {
						goto l21
					}
					position++
				}
			l25:
				{
					position27, tokenIndex27, depth27 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l28
					}
					position++
					goto l27
				l28:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					if buffer[position] != rune('L') {
						goto l21
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
						goto l21
					}
					position++
				}
			l29:
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l32
					}
					position++
					goto l31
				l32:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
					if buffer[position] != rune('C') {
						goto l21
					}
					position++
				}
			l31:
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l34
					}
					position++
					goto l33
				l34:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
					if buffer[position] != rune('T') {
						goto l21
					}
					position++
				}
			l33:
				if !_rules[rulesp]() {
					goto l21
				}
				if !_rules[ruleEmitter]() {
					goto l21
				}
				if !_rules[rulesp]() {
					goto l21
				}
				if !_rules[ruleProjections]() {
					goto l21
				}
				if !_rules[rulesp]() {
					goto l21
				}
				if !_rules[ruleWindowedFrom]() {
					goto l21
				}
				if !_rules[rulesp]() {
					goto l21
				}
				if !_rules[ruleFilter]() {
					goto l21
				}
				if !_rules[rulesp]() {
					goto l21
				}
				if !_rules[ruleGrouping]() {
					goto l21
				}
				if !_rules[rulesp]() {
					goto l21
				}
				if !_rules[ruleHaving]() {
					goto l21
				}
				if !_rules[rulesp]() {
					goto l21
				}
				if !_rules[ruleAction0]() {
					goto l21
				}
				depth--
				add(ruleSelectStmt, position22)
			}
			return true
		l21:
			position, tokenIndex, depth = position21, tokenIndex21, depth21
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
		/* 12 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position301, tokenIndex301, depth301 := position, tokenIndex, depth
			{
				position302 := position
				depth++
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('D') {
						goto l301
					}
					position++
				}
			l303:
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
						goto l301
					}
					position++
				}
			l305:
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l308
					}
					position++
					goto l307
				l308:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
					if buffer[position] != rune('O') {
						goto l301
					}
					position++
				}
			l307:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('P') {
						goto l301
					}
					position++
				}
			l309:
				if !_rules[rulesp]() {
					goto l301
				}
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l312
					}
					position++
					goto l311
				l312:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
					if buffer[position] != rune('S') {
						goto l301
					}
					position++
				}
			l311:
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('O') {
						goto l301
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('U') {
						goto l301
					}
					position++
				}
			l315:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('R') {
						goto l301
					}
					position++
				}
			l317:
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l320
					}
					position++
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if buffer[position] != rune('C') {
						goto l301
					}
					position++
				}
			l319:
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l322
					}
					position++
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if buffer[position] != rune('E') {
						goto l301
					}
					position++
				}
			l321:
				if !_rules[rulesp]() {
					goto l301
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l301
				}
				if !_rules[ruleAction10]() {
					goto l301
				}
				depth--
				add(ruleDropSourceStmt, position302)
			}
			return true
		l301:
			position, tokenIndex, depth = position301, tokenIndex301, depth301
			return false
		},
		/* 13 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action11)> */
		func() bool {
			position323, tokenIndex323, depth323 := position, tokenIndex, depth
			{
				position324 := position
				depth++
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l326
					}
					position++
					goto l325
				l326:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
					if buffer[position] != rune('D') {
						goto l323
					}
					position++
				}
			l325:
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l328
					}
					position++
					goto l327
				l328:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if buffer[position] != rune('R') {
						goto l323
					}
					position++
				}
			l327:
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l330
					}
					position++
					goto l329
				l330:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
					if buffer[position] != rune('O') {
						goto l323
					}
					position++
				}
			l329:
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('P') {
						goto l323
					}
					position++
				}
			l331:
				if !_rules[rulesp]() {
					goto l323
				}
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('S') {
						goto l323
					}
					position++
				}
			l333:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l336
					}
					position++
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if buffer[position] != rune('T') {
						goto l323
					}
					position++
				}
			l335:
				{
					position337, tokenIndex337, depth337 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l338
					}
					position++
					goto l337
				l338:
					position, tokenIndex, depth = position337, tokenIndex337, depth337
					if buffer[position] != rune('R') {
						goto l323
					}
					position++
				}
			l337:
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l340
					}
					position++
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if buffer[position] != rune('E') {
						goto l323
					}
					position++
				}
			l339:
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l342
					}
					position++
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					if buffer[position] != rune('A') {
						goto l323
					}
					position++
				}
			l341:
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('M') {
						goto l323
					}
					position++
				}
			l343:
				if !_rules[rulesp]() {
					goto l323
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l323
				}
				if !_rules[ruleAction11]() {
					goto l323
				}
				depth--
				add(ruleDropStreamStmt, position324)
			}
			return true
		l323:
			position, tokenIndex, depth = position323, tokenIndex323, depth323
			return false
		},
		/* 14 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action12)> */
		func() bool {
			position345, tokenIndex345, depth345 := position, tokenIndex, depth
			{
				position346 := position
				depth++
				{
					position347, tokenIndex347, depth347 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l348
					}
					goto l347
				l348:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if !_rules[ruleDSTREAM]() {
						goto l349
					}
					goto l347
				l349:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if !_rules[ruleRSTREAM]() {
						goto l345
					}
				}
			l347:
				{
					position350 := position
					depth++
					{
						position351, tokenIndex351, depth351 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l351
						}
						if buffer[position] != rune('[') {
							goto l351
						}
						position++
						if !_rules[rulesp]() {
							goto l351
						}
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
								goto l351
							}
							position++
						}
					l353:
						{
							position355, tokenIndex355, depth355 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l356
							}
							position++
							goto l355
						l356:
							position, tokenIndex, depth = position355, tokenIndex355, depth355
							if buffer[position] != rune('V') {
								goto l351
							}
							position++
						}
					l355:
						{
							position357, tokenIndex357, depth357 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l358
							}
							position++
							goto l357
						l358:
							position, tokenIndex, depth = position357, tokenIndex357, depth357
							if buffer[position] != rune('E') {
								goto l351
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
								goto l351
							}
							position++
						}
					l359:
						{
							position361, tokenIndex361, depth361 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l362
							}
							position++
							goto l361
						l362:
							position, tokenIndex, depth = position361, tokenIndex361, depth361
							if buffer[position] != rune('Y') {
								goto l351
							}
							position++
						}
					l361:
						if !_rules[rulesp]() {
							goto l351
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l351
						}
						if !_rules[rulesp]() {
							goto l351
						}
						if buffer[position] != rune(']') {
							goto l351
						}
						position++
						goto l352
					l351:
						position, tokenIndex, depth = position351, tokenIndex351, depth351
					}
				l352:
					depth--
					add(rulePegText, position350)
				}
				if !_rules[ruleAction12]() {
					goto l345
				}
				depth--
				add(ruleEmitter, position346)
			}
			return true
		l345:
			position, tokenIndex, depth = position345, tokenIndex345, depth345
			return false
		},
		/* 15 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position363, tokenIndex363, depth363 := position, tokenIndex, depth
			{
				position364 := position
				depth++
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l366
					}
				l367:
					{
						position368, tokenIndex368, depth368 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l368
						}
						if buffer[position] != rune(',') {
							goto l368
						}
						position++
						if !_rules[rulesp]() {
							goto l368
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l368
						}
						goto l367
					l368:
						position, tokenIndex, depth = position368, tokenIndex368, depth368
					}
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if !_rules[ruleTimeEmitterInterval]() {
						goto l369
					}
					goto l365
				l369:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if !_rules[ruleTupleEmitterInterval]() {
						goto l363
					}
				}
			l365:
				depth--
				add(ruleEmitterIntervals, position364)
			}
			return true
		l363:
			position, tokenIndex, depth = position363, tokenIndex363, depth363
			return false
		},
		/* 16 TimeEmitterInterval <- <(<TimeInterval> Action13)> */
		func() bool {
			position370, tokenIndex370, depth370 := position, tokenIndex, depth
			{
				position371 := position
				depth++
				{
					position372 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l370
					}
					depth--
					add(rulePegText, position372)
				}
				if !_rules[ruleAction13]() {
					goto l370
				}
				depth--
				add(ruleTimeEmitterInterval, position371)
			}
			return true
		l370:
			position, tokenIndex, depth = position370, tokenIndex370, depth370
			return false
		},
		/* 17 TupleEmitterInterval <- <(<TuplesInterval> Action14)> */
		func() bool {
			position373, tokenIndex373, depth373 := position, tokenIndex, depth
			{
				position374 := position
				depth++
				{
					position375 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l373
					}
					depth--
					add(rulePegText, position375)
				}
				if !_rules[ruleAction14]() {
					goto l373
				}
				depth--
				add(ruleTupleEmitterInterval, position374)
			}
			return true
		l373:
			position, tokenIndex, depth = position373, tokenIndex373, depth373
			return false
		},
		/* 18 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action15)> */
		func() bool {
			position376, tokenIndex376, depth376 := position, tokenIndex, depth
			{
				position377 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l376
				}
				if !_rules[rulesp]() {
					goto l376
				}
				{
					position378, tokenIndex378, depth378 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l379
					}
					position++
					goto l378
				l379:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
					if buffer[position] != rune('I') {
						goto l376
					}
					position++
				}
			l378:
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('N') {
						goto l376
					}
					position++
				}
			l380:
				if !_rules[rulesp]() {
					goto l376
				}
				if !_rules[ruleStream]() {
					goto l376
				}
				if !_rules[ruleAction15]() {
					goto l376
				}
				depth--
				add(ruleTupleEmitterFromInterval, position377)
			}
			return true
		l376:
			position, tokenIndex, depth = position376, tokenIndex376, depth376
			return false
		},
		/* 19 Projections <- <(<(Projection sp (',' sp Projection)*)> Action16)> */
		func() bool {
			position382, tokenIndex382, depth382 := position, tokenIndex, depth
			{
				position383 := position
				depth++
				{
					position384 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l382
					}
					if !_rules[rulesp]() {
						goto l382
					}
				l385:
					{
						position386, tokenIndex386, depth386 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l386
						}
						position++
						if !_rules[rulesp]() {
							goto l386
						}
						if !_rules[ruleProjection]() {
							goto l386
						}
						goto l385
					l386:
						position, tokenIndex, depth = position386, tokenIndex386, depth386
					}
					depth--
					add(rulePegText, position384)
				}
				if !_rules[ruleAction16]() {
					goto l382
				}
				depth--
				add(ruleProjections, position383)
			}
			return true
		l382:
			position, tokenIndex, depth = position382, tokenIndex382, depth382
			return false
		},
		/* 20 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position387, tokenIndex387, depth387 := position, tokenIndex, depth
			{
				position388 := position
				depth++
				{
					position389, tokenIndex389, depth389 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l390
					}
					goto l389
				l390:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
					if !_rules[ruleExpression]() {
						goto l391
					}
					goto l389
				l391:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
					if !_rules[ruleWildcard]() {
						goto l387
					}
				}
			l389:
				depth--
				add(ruleProjection, position388)
			}
			return true
		l387:
			position, tokenIndex, depth = position387, tokenIndex387, depth387
			return false
		},
		/* 21 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action17)> */
		func() bool {
			position392, tokenIndex392, depth392 := position, tokenIndex, depth
			{
				position393 := position
				depth++
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l395
					}
					goto l394
				l395:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
					if !_rules[ruleWildcard]() {
						goto l392
					}
				}
			l394:
				if !_rules[rulesp]() {
					goto l392
				}
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('A') {
						goto l392
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if buffer[position] != rune('S') {
						goto l392
					}
					position++
				}
			l398:
				if !_rules[rulesp]() {
					goto l392
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l392
				}
				if !_rules[ruleAction17]() {
					goto l392
				}
				depth--
				add(ruleAliasExpression, position393)
			}
			return true
		l392:
			position, tokenIndex, depth = position392, tokenIndex392, depth392
			return false
		},
		/* 22 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action18)> */
		func() bool {
			position400, tokenIndex400, depth400 := position, tokenIndex, depth
			{
				position401 := position
				depth++
				{
					position402 := position
					depth++
					{
						position403, tokenIndex403, depth403 := position, tokenIndex, depth
						{
							position405, tokenIndex405, depth405 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l406
							}
							position++
							goto l405
						l406:
							position, tokenIndex, depth = position405, tokenIndex405, depth405
							if buffer[position] != rune('F') {
								goto l403
							}
							position++
						}
					l405:
						{
							position407, tokenIndex407, depth407 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l408
							}
							position++
							goto l407
						l408:
							position, tokenIndex, depth = position407, tokenIndex407, depth407
							if buffer[position] != rune('R') {
								goto l403
							}
							position++
						}
					l407:
						{
							position409, tokenIndex409, depth409 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l410
							}
							position++
							goto l409
						l410:
							position, tokenIndex, depth = position409, tokenIndex409, depth409
							if buffer[position] != rune('O') {
								goto l403
							}
							position++
						}
					l409:
						{
							position411, tokenIndex411, depth411 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l412
							}
							position++
							goto l411
						l412:
							position, tokenIndex, depth = position411, tokenIndex411, depth411
							if buffer[position] != rune('M') {
								goto l403
							}
							position++
						}
					l411:
						if !_rules[rulesp]() {
							goto l403
						}
						if !_rules[ruleRelations]() {
							goto l403
						}
						if !_rules[rulesp]() {
							goto l403
						}
						goto l404
					l403:
						position, tokenIndex, depth = position403, tokenIndex403, depth403
					}
				l404:
					depth--
					add(rulePegText, position402)
				}
				if !_rules[ruleAction18]() {
					goto l400
				}
				depth--
				add(ruleWindowedFrom, position401)
			}
			return true
		l400:
			position, tokenIndex, depth = position400, tokenIndex400, depth400
			return false
		},
		/* 23 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l416
					}
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if !_rules[ruleTuplesInterval]() {
						goto l413
					}
				}
			l415:
				depth--
				add(ruleInterval, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 24 TimeInterval <- <(NumericLiteral sp SECONDS Action19)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l417
				}
				if !_rules[rulesp]() {
					goto l417
				}
				if !_rules[ruleSECONDS]() {
					goto l417
				}
				if !_rules[ruleAction19]() {
					goto l417
				}
				depth--
				add(ruleTimeInterval, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 25 TuplesInterval <- <(NumericLiteral sp TUPLES Action20)> */
		func() bool {
			position419, tokenIndex419, depth419 := position, tokenIndex, depth
			{
				position420 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l419
				}
				if !_rules[rulesp]() {
					goto l419
				}
				if !_rules[ruleTUPLES]() {
					goto l419
				}
				if !_rules[ruleAction20]() {
					goto l419
				}
				depth--
				add(ruleTuplesInterval, position420)
			}
			return true
		l419:
			position, tokenIndex, depth = position419, tokenIndex419, depth419
			return false
		},
		/* 26 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position421, tokenIndex421, depth421 := position, tokenIndex, depth
			{
				position422 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l421
				}
				if !_rules[rulesp]() {
					goto l421
				}
			l423:
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l424
					}
					position++
					if !_rules[rulesp]() {
						goto l424
					}
					if !_rules[ruleRelationLike]() {
						goto l424
					}
					goto l423
				l424:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
				}
				depth--
				add(ruleRelations, position422)
			}
			return true
		l421:
			position, tokenIndex, depth = position421, tokenIndex421, depth421
			return false
		},
		/* 27 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action21)> */
		func() bool {
			position425, tokenIndex425, depth425 := position, tokenIndex, depth
			{
				position426 := position
				depth++
				{
					position427 := position
					depth++
					{
						position428, tokenIndex428, depth428 := position, tokenIndex, depth
						{
							position430, tokenIndex430, depth430 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l431
							}
							position++
							goto l430
						l431:
							position, tokenIndex, depth = position430, tokenIndex430, depth430
							if buffer[position] != rune('W') {
								goto l428
							}
							position++
						}
					l430:
						{
							position432, tokenIndex432, depth432 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l433
							}
							position++
							goto l432
						l433:
							position, tokenIndex, depth = position432, tokenIndex432, depth432
							if buffer[position] != rune('H') {
								goto l428
							}
							position++
						}
					l432:
						{
							position434, tokenIndex434, depth434 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l435
							}
							position++
							goto l434
						l435:
							position, tokenIndex, depth = position434, tokenIndex434, depth434
							if buffer[position] != rune('E') {
								goto l428
							}
							position++
						}
					l434:
						{
							position436, tokenIndex436, depth436 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l437
							}
							position++
							goto l436
						l437:
							position, tokenIndex, depth = position436, tokenIndex436, depth436
							if buffer[position] != rune('R') {
								goto l428
							}
							position++
						}
					l436:
						{
							position438, tokenIndex438, depth438 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l439
							}
							position++
							goto l438
						l439:
							position, tokenIndex, depth = position438, tokenIndex438, depth438
							if buffer[position] != rune('E') {
								goto l428
							}
							position++
						}
					l438:
						if !_rules[rulesp]() {
							goto l428
						}
						if !_rules[ruleExpression]() {
							goto l428
						}
						goto l429
					l428:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
					}
				l429:
					depth--
					add(rulePegText, position427)
				}
				if !_rules[ruleAction21]() {
					goto l425
				}
				depth--
				add(ruleFilter, position426)
			}
			return true
		l425:
			position, tokenIndex, depth = position425, tokenIndex425, depth425
			return false
		},
		/* 28 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action22)> */
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
							if buffer[position] != rune('g') {
								goto l446
							}
							position++
							goto l445
						l446:
							position, tokenIndex, depth = position445, tokenIndex445, depth445
							if buffer[position] != rune('G') {
								goto l443
							}
							position++
						}
					l445:
						{
							position447, tokenIndex447, depth447 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l448
							}
							position++
							goto l447
						l448:
							position, tokenIndex, depth = position447, tokenIndex447, depth447
							if buffer[position] != rune('R') {
								goto l443
							}
							position++
						}
					l447:
						{
							position449, tokenIndex449, depth449 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l450
							}
							position++
							goto l449
						l450:
							position, tokenIndex, depth = position449, tokenIndex449, depth449
							if buffer[position] != rune('O') {
								goto l443
							}
							position++
						}
					l449:
						{
							position451, tokenIndex451, depth451 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l452
							}
							position++
							goto l451
						l452:
							position, tokenIndex, depth = position451, tokenIndex451, depth451
							if buffer[position] != rune('U') {
								goto l443
							}
							position++
						}
					l451:
						{
							position453, tokenIndex453, depth453 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l454
							}
							position++
							goto l453
						l454:
							position, tokenIndex, depth = position453, tokenIndex453, depth453
							if buffer[position] != rune('P') {
								goto l443
							}
							position++
						}
					l453:
						if !_rules[rulesp]() {
							goto l443
						}
						{
							position455, tokenIndex455, depth455 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l456
							}
							position++
							goto l455
						l456:
							position, tokenIndex, depth = position455, tokenIndex455, depth455
							if buffer[position] != rune('B') {
								goto l443
							}
							position++
						}
					l455:
						{
							position457, tokenIndex457, depth457 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l458
							}
							position++
							goto l457
						l458:
							position, tokenIndex, depth = position457, tokenIndex457, depth457
							if buffer[position] != rune('Y') {
								goto l443
							}
							position++
						}
					l457:
						if !_rules[rulesp]() {
							goto l443
						}
						if !_rules[ruleGroupList]() {
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
				add(ruleGrouping, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 29 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position459, tokenIndex459, depth459 := position, tokenIndex, depth
			{
				position460 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l459
				}
				if !_rules[rulesp]() {
					goto l459
				}
			l461:
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l462
					}
					position++
					if !_rules[rulesp]() {
						goto l462
					}
					if !_rules[ruleExpression]() {
						goto l462
					}
					goto l461
				l462:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
				}
				depth--
				add(ruleGroupList, position460)
			}
			return true
		l459:
			position, tokenIndex, depth = position459, tokenIndex459, depth459
			return false
		},
		/* 30 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action23)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				{
					position465 := position
					depth++
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						{
							position468, tokenIndex468, depth468 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l469
							}
							position++
							goto l468
						l469:
							position, tokenIndex, depth = position468, tokenIndex468, depth468
							if buffer[position] != rune('H') {
								goto l466
							}
							position++
						}
					l468:
						{
							position470, tokenIndex470, depth470 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l471
							}
							position++
							goto l470
						l471:
							position, tokenIndex, depth = position470, tokenIndex470, depth470
							if buffer[position] != rune('A') {
								goto l466
							}
							position++
						}
					l470:
						{
							position472, tokenIndex472, depth472 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l473
							}
							position++
							goto l472
						l473:
							position, tokenIndex, depth = position472, tokenIndex472, depth472
							if buffer[position] != rune('V') {
								goto l466
							}
							position++
						}
					l472:
						{
							position474, tokenIndex474, depth474 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l475
							}
							position++
							goto l474
						l475:
							position, tokenIndex, depth = position474, tokenIndex474, depth474
							if buffer[position] != rune('I') {
								goto l466
							}
							position++
						}
					l474:
						{
							position476, tokenIndex476, depth476 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l477
							}
							position++
							goto l476
						l477:
							position, tokenIndex, depth = position476, tokenIndex476, depth476
							if buffer[position] != rune('N') {
								goto l466
							}
							position++
						}
					l476:
						{
							position478, tokenIndex478, depth478 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l479
							}
							position++
							goto l478
						l479:
							position, tokenIndex, depth = position478, tokenIndex478, depth478
							if buffer[position] != rune('G') {
								goto l466
							}
							position++
						}
					l478:
						if !_rules[rulesp]() {
							goto l466
						}
						if !_rules[ruleExpression]() {
							goto l466
						}
						goto l467
					l466:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
					}
				l467:
					depth--
					add(rulePegText, position465)
				}
				if !_rules[ruleAction23]() {
					goto l463
				}
				depth--
				add(ruleHaving, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
			return false
		},
		/* 31 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action24))> */
		func() bool {
			position480, tokenIndex480, depth480 := position, tokenIndex, depth
			{
				position481 := position
				depth++
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l483
					}
					goto l482
				l483:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
					if !_rules[ruleStreamWindow]() {
						goto l480
					}
					if !_rules[ruleAction24]() {
						goto l480
					}
				}
			l482:
				depth--
				add(ruleRelationLike, position481)
			}
			return true
		l480:
			position, tokenIndex, depth = position480, tokenIndex480, depth480
			return false
		},
		/* 32 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action25)> */
		func() bool {
			position484, tokenIndex484, depth484 := position, tokenIndex, depth
			{
				position485 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l484
				}
				if !_rules[rulesp]() {
					goto l484
				}
				{
					position486, tokenIndex486, depth486 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l487
					}
					position++
					goto l486
				l487:
					position, tokenIndex, depth = position486, tokenIndex486, depth486
					if buffer[position] != rune('A') {
						goto l484
					}
					position++
				}
			l486:
				{
					position488, tokenIndex488, depth488 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l489
					}
					position++
					goto l488
				l489:
					position, tokenIndex, depth = position488, tokenIndex488, depth488
					if buffer[position] != rune('S') {
						goto l484
					}
					position++
				}
			l488:
				if !_rules[rulesp]() {
					goto l484
				}
				if !_rules[ruleIdentifier]() {
					goto l484
				}
				if !_rules[ruleAction25]() {
					goto l484
				}
				depth--
				add(ruleAliasedStreamWindow, position485)
			}
			return true
		l484:
			position, tokenIndex, depth = position484, tokenIndex484, depth484
			return false
		},
		/* 33 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action26)> */
		func() bool {
			position490, tokenIndex490, depth490 := position, tokenIndex, depth
			{
				position491 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l490
				}
				if !_rules[rulesp]() {
					goto l490
				}
				if buffer[position] != rune('[') {
					goto l490
				}
				position++
				if !_rules[rulesp]() {
					goto l490
				}
				{
					position492, tokenIndex492, depth492 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l493
					}
					position++
					goto l492
				l493:
					position, tokenIndex, depth = position492, tokenIndex492, depth492
					if buffer[position] != rune('R') {
						goto l490
					}
					position++
				}
			l492:
				{
					position494, tokenIndex494, depth494 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l495
					}
					position++
					goto l494
				l495:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if buffer[position] != rune('A') {
						goto l490
					}
					position++
				}
			l494:
				{
					position496, tokenIndex496, depth496 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l497
					}
					position++
					goto l496
				l497:
					position, tokenIndex, depth = position496, tokenIndex496, depth496
					if buffer[position] != rune('N') {
						goto l490
					}
					position++
				}
			l496:
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l499
					}
					position++
					goto l498
				l499:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if buffer[position] != rune('G') {
						goto l490
					}
					position++
				}
			l498:
				{
					position500, tokenIndex500, depth500 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l501
					}
					position++
					goto l500
				l501:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if buffer[position] != rune('E') {
						goto l490
					}
					position++
				}
			l500:
				if !_rules[rulesp]() {
					goto l490
				}
				if !_rules[ruleInterval]() {
					goto l490
				}
				if !_rules[rulesp]() {
					goto l490
				}
				if buffer[position] != rune(']') {
					goto l490
				}
				position++
				if !_rules[ruleAction26]() {
					goto l490
				}
				depth--
				add(ruleStreamWindow, position491)
			}
			return true
		l490:
			position, tokenIndex, depth = position490, tokenIndex490, depth490
			return false
		},
		/* 34 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l505
					}
					goto l504
				l505:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
					if !_rules[ruleStream]() {
						goto l502
					}
				}
			l504:
				depth--
				add(ruleStreamLike, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 35 UDSFFuncApp <- <(FuncApp Action27)> */
		func() bool {
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l506
				}
				if !_rules[ruleAction27]() {
					goto l506
				}
				depth--
				add(ruleUDSFFuncApp, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 36 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action28)> */
		func() bool {
			position508, tokenIndex508, depth508 := position, tokenIndex, depth
			{
				position509 := position
				depth++
				{
					position510 := position
					depth++
					{
						position511, tokenIndex511, depth511 := position, tokenIndex, depth
						{
							position513, tokenIndex513, depth513 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l514
							}
							position++
							goto l513
						l514:
							position, tokenIndex, depth = position513, tokenIndex513, depth513
							if buffer[position] != rune('W') {
								goto l511
							}
							position++
						}
					l513:
						{
							position515, tokenIndex515, depth515 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l516
							}
							position++
							goto l515
						l516:
							position, tokenIndex, depth = position515, tokenIndex515, depth515
							if buffer[position] != rune('I') {
								goto l511
							}
							position++
						}
					l515:
						{
							position517, tokenIndex517, depth517 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l518
							}
							position++
							goto l517
						l518:
							position, tokenIndex, depth = position517, tokenIndex517, depth517
							if buffer[position] != rune('T') {
								goto l511
							}
							position++
						}
					l517:
						{
							position519, tokenIndex519, depth519 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l520
							}
							position++
							goto l519
						l520:
							position, tokenIndex, depth = position519, tokenIndex519, depth519
							if buffer[position] != rune('H') {
								goto l511
							}
							position++
						}
					l519:
						if !_rules[rulesp]() {
							goto l511
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l511
						}
						if !_rules[rulesp]() {
							goto l511
						}
					l521:
						{
							position522, tokenIndex522, depth522 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l522
							}
							position++
							if !_rules[rulesp]() {
								goto l522
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l522
							}
							goto l521
						l522:
							position, tokenIndex, depth = position522, tokenIndex522, depth522
						}
						goto l512
					l511:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
					}
				l512:
					depth--
					add(rulePegText, position510)
				}
				if !_rules[ruleAction28]() {
					goto l508
				}
				depth--
				add(ruleSourceSinkSpecs, position509)
			}
			return true
		l508:
			position, tokenIndex, depth = position508, tokenIndex508, depth508
			return false
		},
		/* 37 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action29)> */
		func() bool {
			position523, tokenIndex523, depth523 := position, tokenIndex, depth
			{
				position524 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l523
				}
				if buffer[position] != rune('=') {
					goto l523
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l523
				}
				if !_rules[ruleAction29]() {
					goto l523
				}
				depth--
				add(ruleSourceSinkParam, position524)
			}
			return true
		l523:
			position, tokenIndex, depth = position523, tokenIndex523, depth523
			return false
		},
		/* 38 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position525, tokenIndex525, depth525 := position, tokenIndex, depth
			{
				position526 := position
				depth++
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l528
					}
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if !_rules[ruleLiteral]() {
						goto l525
					}
				}
			l527:
				depth--
				add(ruleSourceSinkParamVal, position526)
			}
			return true
		l525:
			position, tokenIndex, depth = position525, tokenIndex525, depth525
			return false
		},
		/* 39 PausedOpt <- <(<(Paused / Unpaused)?> Action30)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531 := position
					depth++
					{
						position532, tokenIndex532, depth532 := position, tokenIndex, depth
						{
							position534, tokenIndex534, depth534 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l535
							}
							goto l534
						l535:
							position, tokenIndex, depth = position534, tokenIndex534, depth534
							if !_rules[ruleUnpaused]() {
								goto l532
							}
						}
					l534:
						goto l533
					l532:
						position, tokenIndex, depth = position532, tokenIndex532, depth532
					}
				l533:
					depth--
					add(rulePegText, position531)
				}
				if !_rules[ruleAction30]() {
					goto l529
				}
				depth--
				add(rulePausedOpt, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 40 Expression <- <orExpr> */
		func() bool {
			position536, tokenIndex536, depth536 := position, tokenIndex, depth
			{
				position537 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l536
				}
				depth--
				add(ruleExpression, position537)
			}
			return true
		l536:
			position, tokenIndex, depth = position536, tokenIndex536, depth536
			return false
		},
		/* 41 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action31)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					position540 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l538
					}
					if !_rules[rulesp]() {
						goto l538
					}
					{
						position541, tokenIndex541, depth541 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l541
						}
						if !_rules[rulesp]() {
							goto l541
						}
						if !_rules[ruleandExpr]() {
							goto l541
						}
						goto l542
					l541:
						position, tokenIndex, depth = position541, tokenIndex541, depth541
					}
				l542:
					depth--
					add(rulePegText, position540)
				}
				if !_rules[ruleAction31]() {
					goto l538
				}
				depth--
				add(ruleorExpr, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 42 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action32)> */
		func() bool {
			position543, tokenIndex543, depth543 := position, tokenIndex, depth
			{
				position544 := position
				depth++
				{
					position545 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l543
					}
					if !_rules[rulesp]() {
						goto l543
					}
					{
						position546, tokenIndex546, depth546 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l546
						}
						if !_rules[rulesp]() {
							goto l546
						}
						if !_rules[rulenotExpr]() {
							goto l546
						}
						goto l547
					l546:
						position, tokenIndex, depth = position546, tokenIndex546, depth546
					}
				l547:
					depth--
					add(rulePegText, position545)
				}
				if !_rules[ruleAction32]() {
					goto l543
				}
				depth--
				add(ruleandExpr, position544)
			}
			return true
		l543:
			position, tokenIndex, depth = position543, tokenIndex543, depth543
			return false
		},
		/* 43 notExpr <- <(<((Not sp)? comparisonExpr)> Action33)> */
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
						if !_rules[ruleNot]() {
							goto l551
						}
						if !_rules[rulesp]() {
							goto l551
						}
						goto l552
					l551:
						position, tokenIndex, depth = position551, tokenIndex551, depth551
					}
				l552:
					if !_rules[rulecomparisonExpr]() {
						goto l548
					}
					depth--
					add(rulePegText, position550)
				}
				if !_rules[ruleAction33]() {
					goto l548
				}
				depth--
				add(rulenotExpr, position549)
			}
			return true
		l548:
			position, tokenIndex, depth = position548, tokenIndex548, depth548
			return false
		},
		/* 44 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action34)> */
		func() bool {
			position553, tokenIndex553, depth553 := position, tokenIndex, depth
			{
				position554 := position
				depth++
				{
					position555 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l553
					}
					if !_rules[rulesp]() {
						goto l553
					}
					{
						position556, tokenIndex556, depth556 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l556
						}
						if !_rules[rulesp]() {
							goto l556
						}
						if !_rules[ruleisExpr]() {
							goto l556
						}
						goto l557
					l556:
						position, tokenIndex, depth = position556, tokenIndex556, depth556
					}
				l557:
					depth--
					add(rulePegText, position555)
				}
				if !_rules[ruleAction34]() {
					goto l553
				}
				depth--
				add(rulecomparisonExpr, position554)
			}
			return true
		l553:
			position, tokenIndex, depth = position553, tokenIndex553, depth553
			return false
		},
		/* 45 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action35)> */
		func() bool {
			position558, tokenIndex558, depth558 := position, tokenIndex, depth
			{
				position559 := position
				depth++
				{
					position560 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l558
					}
					if !_rules[rulesp]() {
						goto l558
					}
					{
						position561, tokenIndex561, depth561 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l561
						}
						if !_rules[rulesp]() {
							goto l561
						}
						if !_rules[ruleNullLiteral]() {
							goto l561
						}
						goto l562
					l561:
						position, tokenIndex, depth = position561, tokenIndex561, depth561
					}
				l562:
					depth--
					add(rulePegText, position560)
				}
				if !_rules[ruleAction35]() {
					goto l558
				}
				depth--
				add(ruleisExpr, position559)
			}
			return true
		l558:
			position, tokenIndex, depth = position558, tokenIndex558, depth558
			return false
		},
		/* 46 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action36)> */
		func() bool {
			position563, tokenIndex563, depth563 := position, tokenIndex, depth
			{
				position564 := position
				depth++
				{
					position565 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l563
					}
					if !_rules[rulesp]() {
						goto l563
					}
					{
						position566, tokenIndex566, depth566 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l566
						}
						if !_rules[rulesp]() {
							goto l566
						}
						if !_rules[ruleproductExpr]() {
							goto l566
						}
						goto l567
					l566:
						position, tokenIndex, depth = position566, tokenIndex566, depth566
					}
				l567:
					depth--
					add(rulePegText, position565)
				}
				if !_rules[ruleAction36]() {
					goto l563
				}
				depth--
				add(ruletermExpr, position564)
			}
			return true
		l563:
			position, tokenIndex, depth = position563, tokenIndex563, depth563
			return false
		},
		/* 47 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action37)> */
		func() bool {
			position568, tokenIndex568, depth568 := position, tokenIndex, depth
			{
				position569 := position
				depth++
				{
					position570 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l568
					}
					if !_rules[rulesp]() {
						goto l568
					}
					{
						position571, tokenIndex571, depth571 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l571
						}
						if !_rules[rulesp]() {
							goto l571
						}
						if !_rules[ruleminusExpr]() {
							goto l571
						}
						goto l572
					l571:
						position, tokenIndex, depth = position571, tokenIndex571, depth571
					}
				l572:
					depth--
					add(rulePegText, position570)
				}
				if !_rules[ruleAction37]() {
					goto l568
				}
				depth--
				add(ruleproductExpr, position569)
			}
			return true
		l568:
			position, tokenIndex, depth = position568, tokenIndex568, depth568
			return false
		},
		/* 48 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action38)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l576
						}
						if !_rules[rulesp]() {
							goto l576
						}
						goto l577
					l576:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
					}
				l577:
					if !_rules[rulebaseExpr]() {
						goto l573
					}
					depth--
					add(rulePegText, position575)
				}
				if !_rules[ruleAction38]() {
					goto l573
				}
				depth--
				add(ruleminusExpr, position574)
			}
			return true
		l573:
			position, tokenIndex, depth = position573, tokenIndex573, depth573
			return false
		},
		/* 49 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position578, tokenIndex578, depth578 := position, tokenIndex, depth
			{
				position579 := position
				depth++
				{
					position580, tokenIndex580, depth580 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l581
					}
					position++
					if !_rules[rulesp]() {
						goto l581
					}
					if !_rules[ruleExpression]() {
						goto l581
					}
					if !_rules[rulesp]() {
						goto l581
					}
					if buffer[position] != rune(')') {
						goto l581
					}
					position++
					goto l580
				l581:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if !_rules[ruleBooleanLiteral]() {
						goto l582
					}
					goto l580
				l582:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if !_rules[ruleNullLiteral]() {
						goto l583
					}
					goto l580
				l583:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if !_rules[ruleFuncApp]() {
						goto l584
					}
					goto l580
				l584:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if !_rules[ruleRowMeta]() {
						goto l585
					}
					goto l580
				l585:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if !_rules[ruleRowValue]() {
						goto l586
					}
					goto l580
				l586:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if !_rules[ruleLiteral]() {
						goto l578
					}
				}
			l580:
				depth--
				add(rulebaseExpr, position579)
			}
			return true
		l578:
			position, tokenIndex, depth = position578, tokenIndex578, depth578
			return false
		},
		/* 50 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action39)> */
		func() bool {
			position587, tokenIndex587, depth587 := position, tokenIndex, depth
			{
				position588 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l587
				}
				if !_rules[rulesp]() {
					goto l587
				}
				if buffer[position] != rune('(') {
					goto l587
				}
				position++
				if !_rules[rulesp]() {
					goto l587
				}
				if !_rules[ruleFuncParams]() {
					goto l587
				}
				if !_rules[rulesp]() {
					goto l587
				}
				if buffer[position] != rune(')') {
					goto l587
				}
				position++
				if !_rules[ruleAction39]() {
					goto l587
				}
				depth--
				add(ruleFuncApp, position588)
			}
			return true
		l587:
			position, tokenIndex, depth = position587, tokenIndex587, depth587
			return false
		},
		/* 51 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action40)> */
		func() bool {
			position589, tokenIndex589, depth589 := position, tokenIndex, depth
			{
				position590 := position
				depth++
				{
					position591 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l589
					}
					if !_rules[rulesp]() {
						goto l589
					}
				l592:
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l593
						}
						position++
						if !_rules[rulesp]() {
							goto l593
						}
						if !_rules[ruleExpression]() {
							goto l593
						}
						goto l592
					l593:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
					}
					depth--
					add(rulePegText, position591)
				}
				if !_rules[ruleAction40]() {
					goto l589
				}
				depth--
				add(ruleFuncParams, position590)
			}
			return true
		l589:
			position, tokenIndex, depth = position589, tokenIndex589, depth589
			return false
		},
		/* 52 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position594, tokenIndex594, depth594 := position, tokenIndex, depth
			{
				position595 := position
				depth++
				{
					position596, tokenIndex596, depth596 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l597
					}
					goto l596
				l597:
					position, tokenIndex, depth = position596, tokenIndex596, depth596
					if !_rules[ruleNumericLiteral]() {
						goto l598
					}
					goto l596
				l598:
					position, tokenIndex, depth = position596, tokenIndex596, depth596
					if !_rules[ruleStringLiteral]() {
						goto l594
					}
				}
			l596:
				depth--
				add(ruleLiteral, position595)
			}
			return true
		l594:
			position, tokenIndex, depth = position594, tokenIndex594, depth594
			return false
		},
		/* 53 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				{
					position601, tokenIndex601, depth601 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l602
					}
					goto l601
				l602:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if !_rules[ruleNotEqual]() {
						goto l603
					}
					goto l601
				l603:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if !_rules[ruleLessOrEqual]() {
						goto l604
					}
					goto l601
				l604:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if !_rules[ruleLess]() {
						goto l605
					}
					goto l601
				l605:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if !_rules[ruleGreaterOrEqual]() {
						goto l606
					}
					goto l601
				l606:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if !_rules[ruleGreater]() {
						goto l607
					}
					goto l601
				l607:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if !_rules[ruleNotEqual]() {
						goto l599
					}
				}
			l601:
				depth--
				add(ruleComparisonOp, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 54 IsOp <- <(IsNot / Is)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610, tokenIndex610, depth610 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l611
					}
					goto l610
				l611:
					position, tokenIndex, depth = position610, tokenIndex610, depth610
					if !_rules[ruleIs]() {
						goto l608
					}
				}
			l610:
				depth--
				add(ruleIsOp, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 55 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position612, tokenIndex612, depth612 := position, tokenIndex, depth
			{
				position613 := position
				depth++
				{
					position614, tokenIndex614, depth614 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l615
					}
					goto l614
				l615:
					position, tokenIndex, depth = position614, tokenIndex614, depth614
					if !_rules[ruleMinus]() {
						goto l612
					}
				}
			l614:
				depth--
				add(rulePlusMinusOp, position613)
			}
			return true
		l612:
			position, tokenIndex, depth = position612, tokenIndex612, depth612
			return false
		},
		/* 56 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position616, tokenIndex616, depth616 := position, tokenIndex, depth
			{
				position617 := position
				depth++
				{
					position618, tokenIndex618, depth618 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l619
					}
					goto l618
				l619:
					position, tokenIndex, depth = position618, tokenIndex618, depth618
					if !_rules[ruleDivide]() {
						goto l620
					}
					goto l618
				l620:
					position, tokenIndex, depth = position618, tokenIndex618, depth618
					if !_rules[ruleModulo]() {
						goto l616
					}
				}
			l618:
				depth--
				add(ruleMultDivOp, position617)
			}
			return true
		l616:
			position, tokenIndex, depth = position616, tokenIndex616, depth616
			return false
		},
		/* 57 Stream <- <(<ident> Action41)> */
		func() bool {
			position621, tokenIndex621, depth621 := position, tokenIndex, depth
			{
				position622 := position
				depth++
				{
					position623 := position
					depth++
					if !_rules[ruleident]() {
						goto l621
					}
					depth--
					add(rulePegText, position623)
				}
				if !_rules[ruleAction41]() {
					goto l621
				}
				depth--
				add(ruleStream, position622)
			}
			return true
		l621:
			position, tokenIndex, depth = position621, tokenIndex621, depth621
			return false
		},
		/* 58 RowMeta <- <RowTimestamp> */
		func() bool {
			position624, tokenIndex624, depth624 := position, tokenIndex, depth
			{
				position625 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l624
				}
				depth--
				add(ruleRowMeta, position625)
			}
			return true
		l624:
			position, tokenIndex, depth = position624, tokenIndex624, depth624
			return false
		},
		/* 59 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action42)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
				{
					position628 := position
					depth++
					{
						position629, tokenIndex629, depth629 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l629
						}
						if buffer[position] != rune(':') {
							goto l629
						}
						position++
						goto l630
					l629:
						position, tokenIndex, depth = position629, tokenIndex629, depth629
					}
				l630:
					if buffer[position] != rune('t') {
						goto l626
					}
					position++
					if buffer[position] != rune('s') {
						goto l626
					}
					position++
					if buffer[position] != rune('(') {
						goto l626
					}
					position++
					if buffer[position] != rune(')') {
						goto l626
					}
					position++
					depth--
					add(rulePegText, position628)
				}
				if !_rules[ruleAction42]() {
					goto l626
				}
				depth--
				add(ruleRowTimestamp, position627)
			}
			return true
		l626:
			position, tokenIndex, depth = position626, tokenIndex626, depth626
			return false
		},
		/* 60 RowValue <- <(<((ident ':')? jsonPath)> Action43)> */
		func() bool {
			position631, tokenIndex631, depth631 := position, tokenIndex, depth
			{
				position632 := position
				depth++
				{
					position633 := position
					depth++
					{
						position634, tokenIndex634, depth634 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l634
						}
						if buffer[position] != rune(':') {
							goto l634
						}
						position++
						goto l635
					l634:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
					}
				l635:
					if !_rules[rulejsonPath]() {
						goto l631
					}
					depth--
					add(rulePegText, position633)
				}
				if !_rules[ruleAction43]() {
					goto l631
				}
				depth--
				add(ruleRowValue, position632)
			}
			return true
		l631:
			position, tokenIndex, depth = position631, tokenIndex631, depth631
			return false
		},
		/* 61 NumericLiteral <- <(<('-'? [0-9]+)> Action44)> */
		func() bool {
			position636, tokenIndex636, depth636 := position, tokenIndex, depth
			{
				position637 := position
				depth++
				{
					position638 := position
					depth++
					{
						position639, tokenIndex639, depth639 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l639
						}
						position++
						goto l640
					l639:
						position, tokenIndex, depth = position639, tokenIndex639, depth639
					}
				l640:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l636
					}
					position++
				l641:
					{
						position642, tokenIndex642, depth642 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l642
						}
						position++
						goto l641
					l642:
						position, tokenIndex, depth = position642, tokenIndex642, depth642
					}
					depth--
					add(rulePegText, position638)
				}
				if !_rules[ruleAction44]() {
					goto l636
				}
				depth--
				add(ruleNumericLiteral, position637)
			}
			return true
		l636:
			position, tokenIndex, depth = position636, tokenIndex636, depth636
			return false
		},
		/* 62 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action45)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645 := position
					depth++
					{
						position646, tokenIndex646, depth646 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l646
						}
						position++
						goto l647
					l646:
						position, tokenIndex, depth = position646, tokenIndex646, depth646
					}
				l647:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l643
					}
					position++
				l648:
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l649
						}
						position++
						goto l648
					l649:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
					}
					if buffer[position] != rune('.') {
						goto l643
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l643
					}
					position++
				l650:
					{
						position651, tokenIndex651, depth651 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l651
						}
						position++
						goto l650
					l651:
						position, tokenIndex, depth = position651, tokenIndex651, depth651
					}
					depth--
					add(rulePegText, position645)
				}
				if !_rules[ruleAction45]() {
					goto l643
				}
				depth--
				add(ruleFloatLiteral, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 63 Function <- <(<ident> Action46)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				{
					position654 := position
					depth++
					if !_rules[ruleident]() {
						goto l652
					}
					depth--
					add(rulePegText, position654)
				}
				if !_rules[ruleAction46]() {
					goto l652
				}
				depth--
				add(ruleFunction, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 64 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action47)> */
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
						if buffer[position] != rune('n') {
							goto l659
						}
						position++
						goto l658
					l659:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
						if buffer[position] != rune('N') {
							goto l655
						}
						position++
					}
				l658:
					{
						position660, tokenIndex660, depth660 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l661
						}
						position++
						goto l660
					l661:
						position, tokenIndex, depth = position660, tokenIndex660, depth660
						if buffer[position] != rune('U') {
							goto l655
						}
						position++
					}
				l660:
					{
						position662, tokenIndex662, depth662 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l663
						}
						position++
						goto l662
					l663:
						position, tokenIndex, depth = position662, tokenIndex662, depth662
						if buffer[position] != rune('L') {
							goto l655
						}
						position++
					}
				l662:
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l665
						}
						position++
						goto l664
					l665:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if buffer[position] != rune('L') {
							goto l655
						}
						position++
					}
				l664:
					depth--
					add(rulePegText, position657)
				}
				if !_rules[ruleAction47]() {
					goto l655
				}
				depth--
				add(ruleNullLiteral, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 65 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				{
					position668, tokenIndex668, depth668 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l669
					}
					goto l668
				l669:
					position, tokenIndex, depth = position668, tokenIndex668, depth668
					if !_rules[ruleFALSE]() {
						goto l666
					}
				}
			l668:
				depth--
				add(ruleBooleanLiteral, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 66 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action48)> */
		func() bool {
			position670, tokenIndex670, depth670 := position, tokenIndex, depth
			{
				position671 := position
				depth++
				{
					position672 := position
					depth++
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
							goto l670
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
							goto l670
						}
						position++
					}
				l675:
					{
						position677, tokenIndex677, depth677 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l678
						}
						position++
						goto l677
					l678:
						position, tokenIndex, depth = position677, tokenIndex677, depth677
						if buffer[position] != rune('U') {
							goto l670
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
							goto l670
						}
						position++
					}
				l679:
					depth--
					add(rulePegText, position672)
				}
				if !_rules[ruleAction48]() {
					goto l670
				}
				depth--
				add(ruleTRUE, position671)
			}
			return true
		l670:
			position, tokenIndex, depth = position670, tokenIndex670, depth670
			return false
		},
		/* 67 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action49)> */
		func() bool {
			position681, tokenIndex681, depth681 := position, tokenIndex, depth
			{
				position682 := position
				depth++
				{
					position683 := position
					depth++
					{
						position684, tokenIndex684, depth684 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l685
						}
						position++
						goto l684
					l685:
						position, tokenIndex, depth = position684, tokenIndex684, depth684
						if buffer[position] != rune('F') {
							goto l681
						}
						position++
					}
				l684:
					{
						position686, tokenIndex686, depth686 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
						if buffer[position] != rune('A') {
							goto l681
						}
						position++
					}
				l686:
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if buffer[position] != rune('L') {
							goto l681
						}
						position++
					}
				l688:
					{
						position690, tokenIndex690, depth690 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l691
						}
						position++
						goto l690
					l691:
						position, tokenIndex, depth = position690, tokenIndex690, depth690
						if buffer[position] != rune('S') {
							goto l681
						}
						position++
					}
				l690:
					{
						position692, tokenIndex692, depth692 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l693
						}
						position++
						goto l692
					l693:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						if buffer[position] != rune('E') {
							goto l681
						}
						position++
					}
				l692:
					depth--
					add(rulePegText, position683)
				}
				if !_rules[ruleAction49]() {
					goto l681
				}
				depth--
				add(ruleFALSE, position682)
			}
			return true
		l681:
			position, tokenIndex, depth = position681, tokenIndex681, depth681
			return false
		},
		/* 68 Wildcard <- <(<'*'> Action50)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				{
					position696 := position
					depth++
					if buffer[position] != rune('*') {
						goto l694
					}
					position++
					depth--
					add(rulePegText, position696)
				}
				if !_rules[ruleAction50]() {
					goto l694
				}
				depth--
				add(ruleWildcard, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 69 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action51)> */
		func() bool {
			position697, tokenIndex697, depth697 := position, tokenIndex, depth
			{
				position698 := position
				depth++
				{
					position699 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l697
					}
					position++
				l700:
					{
						position701, tokenIndex701, depth701 := position, tokenIndex, depth
						{
							position702, tokenIndex702, depth702 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l703
							}
							position++
							if buffer[position] != rune('\'') {
								goto l703
							}
							position++
							goto l702
						l703:
							position, tokenIndex, depth = position702, tokenIndex702, depth702
							{
								position704, tokenIndex704, depth704 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l704
								}
								position++
								goto l701
							l704:
								position, tokenIndex, depth = position704, tokenIndex704, depth704
							}
							if !matchDot() {
								goto l701
							}
						}
					l702:
						goto l700
					l701:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
					}
					if buffer[position] != rune('\'') {
						goto l697
					}
					position++
					depth--
					add(rulePegText, position699)
				}
				if !_rules[ruleAction51]() {
					goto l697
				}
				depth--
				add(ruleStringLiteral, position698)
			}
			return true
		l697:
			position, tokenIndex, depth = position697, tokenIndex697, depth697
			return false
		},
		/* 70 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action52)> */
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
						if buffer[position] != rune('i') {
							goto l709
						}
						position++
						goto l708
					l709:
						position, tokenIndex, depth = position708, tokenIndex708, depth708
						if buffer[position] != rune('I') {
							goto l705
						}
						position++
					}
				l708:
					{
						position710, tokenIndex710, depth710 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l711
						}
						position++
						goto l710
					l711:
						position, tokenIndex, depth = position710, tokenIndex710, depth710
						if buffer[position] != rune('S') {
							goto l705
						}
						position++
					}
				l710:
					{
						position712, tokenIndex712, depth712 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l713
						}
						position++
						goto l712
					l713:
						position, tokenIndex, depth = position712, tokenIndex712, depth712
						if buffer[position] != rune('T') {
							goto l705
						}
						position++
					}
				l712:
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l715
						}
						position++
						goto l714
					l715:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
						if buffer[position] != rune('R') {
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
						if buffer[position] != rune('a') {
							goto l719
						}
						position++
						goto l718
					l719:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
						if buffer[position] != rune('A') {
							goto l705
						}
						position++
					}
				l718:
					{
						position720, tokenIndex720, depth720 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l721
						}
						position++
						goto l720
					l721:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
						if buffer[position] != rune('M') {
							goto l705
						}
						position++
					}
				l720:
					depth--
					add(rulePegText, position707)
				}
				if !_rules[ruleAction52]() {
					goto l705
				}
				depth--
				add(ruleISTREAM, position706)
			}
			return true
		l705:
			position, tokenIndex, depth = position705, tokenIndex705, depth705
			return false
		},
		/* 71 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action53)> */
		func() bool {
			position722, tokenIndex722, depth722 := position, tokenIndex, depth
			{
				position723 := position
				depth++
				{
					position724 := position
					depth++
					{
						position725, tokenIndex725, depth725 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l726
						}
						position++
						goto l725
					l726:
						position, tokenIndex, depth = position725, tokenIndex725, depth725
						if buffer[position] != rune('D') {
							goto l722
						}
						position++
					}
				l725:
					{
						position727, tokenIndex727, depth727 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l728
						}
						position++
						goto l727
					l728:
						position, tokenIndex, depth = position727, tokenIndex727, depth727
						if buffer[position] != rune('S') {
							goto l722
						}
						position++
					}
				l727:
					{
						position729, tokenIndex729, depth729 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l730
						}
						position++
						goto l729
					l730:
						position, tokenIndex, depth = position729, tokenIndex729, depth729
						if buffer[position] != rune('T') {
							goto l722
						}
						position++
					}
				l729:
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l732
						}
						position++
						goto l731
					l732:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
						if buffer[position] != rune('R') {
							goto l722
						}
						position++
					}
				l731:
					{
						position733, tokenIndex733, depth733 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l734
						}
						position++
						goto l733
					l734:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
						if buffer[position] != rune('E') {
							goto l722
						}
						position++
					}
				l733:
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l736
						}
						position++
						goto l735
					l736:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if buffer[position] != rune('A') {
							goto l722
						}
						position++
					}
				l735:
					{
						position737, tokenIndex737, depth737 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l738
						}
						position++
						goto l737
					l738:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
						if buffer[position] != rune('M') {
							goto l722
						}
						position++
					}
				l737:
					depth--
					add(rulePegText, position724)
				}
				if !_rules[ruleAction53]() {
					goto l722
				}
				depth--
				add(ruleDSTREAM, position723)
			}
			return true
		l722:
			position, tokenIndex, depth = position722, tokenIndex722, depth722
			return false
		},
		/* 72 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action54)> */
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
						if buffer[position] != rune('r') {
							goto l743
						}
						position++
						goto l742
					l743:
						position, tokenIndex, depth = position742, tokenIndex742, depth742
						if buffer[position] != rune('R') {
							goto l739
						}
						position++
					}
				l742:
					{
						position744, tokenIndex744, depth744 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l745
						}
						position++
						goto l744
					l745:
						position, tokenIndex, depth = position744, tokenIndex744, depth744
						if buffer[position] != rune('S') {
							goto l739
						}
						position++
					}
				l744:
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
							goto l739
						}
						position++
					}
				l746:
					{
						position748, tokenIndex748, depth748 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l749
						}
						position++
						goto l748
					l749:
						position, tokenIndex, depth = position748, tokenIndex748, depth748
						if buffer[position] != rune('R') {
							goto l739
						}
						position++
					}
				l748:
					{
						position750, tokenIndex750, depth750 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l751
						}
						position++
						goto l750
					l751:
						position, tokenIndex, depth = position750, tokenIndex750, depth750
						if buffer[position] != rune('E') {
							goto l739
						}
						position++
					}
				l750:
					{
						position752, tokenIndex752, depth752 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l753
						}
						position++
						goto l752
					l753:
						position, tokenIndex, depth = position752, tokenIndex752, depth752
						if buffer[position] != rune('A') {
							goto l739
						}
						position++
					}
				l752:
					{
						position754, tokenIndex754, depth754 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l755
						}
						position++
						goto l754
					l755:
						position, tokenIndex, depth = position754, tokenIndex754, depth754
						if buffer[position] != rune('M') {
							goto l739
						}
						position++
					}
				l754:
					depth--
					add(rulePegText, position741)
				}
				if !_rules[ruleAction54]() {
					goto l739
				}
				depth--
				add(ruleRSTREAM, position740)
			}
			return true
		l739:
			position, tokenIndex, depth = position739, tokenIndex739, depth739
			return false
		},
		/* 73 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action55)> */
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
						if buffer[position] != rune('t') {
							goto l760
						}
						position++
						goto l759
					l760:
						position, tokenIndex, depth = position759, tokenIndex759, depth759
						if buffer[position] != rune('T') {
							goto l756
						}
						position++
					}
				l759:
					{
						position761, tokenIndex761, depth761 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l762
						}
						position++
						goto l761
					l762:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
						if buffer[position] != rune('U') {
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
						if buffer[position] != rune('l') {
							goto l766
						}
						position++
						goto l765
					l766:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
						if buffer[position] != rune('L') {
							goto l756
						}
						position++
					}
				l765:
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('E') {
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
					depth--
					add(rulePegText, position758)
				}
				if !_rules[ruleAction55]() {
					goto l756
				}
				depth--
				add(ruleTUPLES, position757)
			}
			return true
		l756:
			position, tokenIndex, depth = position756, tokenIndex756, depth756
			return false
		},
		/* 74 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action56)> */
		func() bool {
			position771, tokenIndex771, depth771 := position, tokenIndex, depth
			{
				position772 := position
				depth++
				{
					position773 := position
					depth++
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
							goto l771
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
							goto l771
						}
						position++
					}
				l776:
					{
						position778, tokenIndex778, depth778 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l779
						}
						position++
						goto l778
					l779:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						if buffer[position] != rune('C') {
							goto l771
						}
						position++
					}
				l778:
					{
						position780, tokenIndex780, depth780 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l781
						}
						position++
						goto l780
					l781:
						position, tokenIndex, depth = position780, tokenIndex780, depth780
						if buffer[position] != rune('O') {
							goto l771
						}
						position++
					}
				l780:
					{
						position782, tokenIndex782, depth782 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l783
						}
						position++
						goto l782
					l783:
						position, tokenIndex, depth = position782, tokenIndex782, depth782
						if buffer[position] != rune('N') {
							goto l771
						}
						position++
					}
				l782:
					{
						position784, tokenIndex784, depth784 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l785
						}
						position++
						goto l784
					l785:
						position, tokenIndex, depth = position784, tokenIndex784, depth784
						if buffer[position] != rune('D') {
							goto l771
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
							goto l771
						}
						position++
					}
				l786:
					depth--
					add(rulePegText, position773)
				}
				if !_rules[ruleAction56]() {
					goto l771
				}
				depth--
				add(ruleSECONDS, position772)
			}
			return true
		l771:
			position, tokenIndex, depth = position771, tokenIndex771, depth771
			return false
		},
		/* 75 StreamIdentifier <- <(<ident> Action57)> */
		func() bool {
			position788, tokenIndex788, depth788 := position, tokenIndex, depth
			{
				position789 := position
				depth++
				{
					position790 := position
					depth++
					if !_rules[ruleident]() {
						goto l788
					}
					depth--
					add(rulePegText, position790)
				}
				if !_rules[ruleAction57]() {
					goto l788
				}
				depth--
				add(ruleStreamIdentifier, position789)
			}
			return true
		l788:
			position, tokenIndex, depth = position788, tokenIndex788, depth788
			return false
		},
		/* 76 SourceSinkType <- <(<ident> Action58)> */
		func() bool {
			position791, tokenIndex791, depth791 := position, tokenIndex, depth
			{
				position792 := position
				depth++
				{
					position793 := position
					depth++
					if !_rules[ruleident]() {
						goto l791
					}
					depth--
					add(rulePegText, position793)
				}
				if !_rules[ruleAction58]() {
					goto l791
				}
				depth--
				add(ruleSourceSinkType, position792)
			}
			return true
		l791:
			position, tokenIndex, depth = position791, tokenIndex791, depth791
			return false
		},
		/* 77 SourceSinkParamKey <- <(<ident> Action59)> */
		func() bool {
			position794, tokenIndex794, depth794 := position, tokenIndex, depth
			{
				position795 := position
				depth++
				{
					position796 := position
					depth++
					if !_rules[ruleident]() {
						goto l794
					}
					depth--
					add(rulePegText, position796)
				}
				if !_rules[ruleAction59]() {
					goto l794
				}
				depth--
				add(ruleSourceSinkParamKey, position795)
			}
			return true
		l794:
			position, tokenIndex, depth = position794, tokenIndex794, depth794
			return false
		},
		/* 78 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action60)> */
		func() bool {
			position797, tokenIndex797, depth797 := position, tokenIndex, depth
			{
				position798 := position
				depth++
				{
					position799 := position
					depth++
					{
						position800, tokenIndex800, depth800 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l801
						}
						position++
						goto l800
					l801:
						position, tokenIndex, depth = position800, tokenIndex800, depth800
						if buffer[position] != rune('P') {
							goto l797
						}
						position++
					}
				l800:
					{
						position802, tokenIndex802, depth802 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l803
						}
						position++
						goto l802
					l803:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if buffer[position] != rune('A') {
							goto l797
						}
						position++
					}
				l802:
					{
						position804, tokenIndex804, depth804 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l805
						}
						position++
						goto l804
					l805:
						position, tokenIndex, depth = position804, tokenIndex804, depth804
						if buffer[position] != rune('U') {
							goto l797
						}
						position++
					}
				l804:
					{
						position806, tokenIndex806, depth806 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l807
						}
						position++
						goto l806
					l807:
						position, tokenIndex, depth = position806, tokenIndex806, depth806
						if buffer[position] != rune('S') {
							goto l797
						}
						position++
					}
				l806:
					{
						position808, tokenIndex808, depth808 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l809
						}
						position++
						goto l808
					l809:
						position, tokenIndex, depth = position808, tokenIndex808, depth808
						if buffer[position] != rune('E') {
							goto l797
						}
						position++
					}
				l808:
					{
						position810, tokenIndex810, depth810 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l811
						}
						position++
						goto l810
					l811:
						position, tokenIndex, depth = position810, tokenIndex810, depth810
						if buffer[position] != rune('D') {
							goto l797
						}
						position++
					}
				l810:
					depth--
					add(rulePegText, position799)
				}
				if !_rules[ruleAction60]() {
					goto l797
				}
				depth--
				add(rulePaused, position798)
			}
			return true
		l797:
			position, tokenIndex, depth = position797, tokenIndex797, depth797
			return false
		},
		/* 79 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action61)> */
		func() bool {
			position812, tokenIndex812, depth812 := position, tokenIndex, depth
			{
				position813 := position
				depth++
				{
					position814 := position
					depth++
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
							goto l812
						}
						position++
					}
				l815:
					{
						position817, tokenIndex817, depth817 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l818
						}
						position++
						goto l817
					l818:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						if buffer[position] != rune('N') {
							goto l812
						}
						position++
					}
				l817:
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
							goto l812
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
							goto l812
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
							goto l812
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
							goto l812
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
							goto l812
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
							goto l812
						}
						position++
					}
				l829:
					depth--
					add(rulePegText, position814)
				}
				if !_rules[ruleAction61]() {
					goto l812
				}
				depth--
				add(ruleUnpaused, position813)
			}
			return true
		l812:
			position, tokenIndex, depth = position812, tokenIndex812, depth812
			return false
		},
		/* 80 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action62)> */
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
						if buffer[position] != rune('o') {
							goto l835
						}
						position++
						goto l834
					l835:
						position, tokenIndex, depth = position834, tokenIndex834, depth834
						if buffer[position] != rune('O') {
							goto l831
						}
						position++
					}
				l834:
					{
						position836, tokenIndex836, depth836 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l837
						}
						position++
						goto l836
					l837:
						position, tokenIndex, depth = position836, tokenIndex836, depth836
						if buffer[position] != rune('R') {
							goto l831
						}
						position++
					}
				l836:
					depth--
					add(rulePegText, position833)
				}
				if !_rules[ruleAction62]() {
					goto l831
				}
				depth--
				add(ruleOr, position832)
			}
			return true
		l831:
			position, tokenIndex, depth = position831, tokenIndex831, depth831
			return false
		},
		/* 81 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action63)> */
		func() bool {
			position838, tokenIndex838, depth838 := position, tokenIndex, depth
			{
				position839 := position
				depth++
				{
					position840 := position
					depth++
					{
						position841, tokenIndex841, depth841 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l842
						}
						position++
						goto l841
					l842:
						position, tokenIndex, depth = position841, tokenIndex841, depth841
						if buffer[position] != rune('A') {
							goto l838
						}
						position++
					}
				l841:
					{
						position843, tokenIndex843, depth843 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l844
						}
						position++
						goto l843
					l844:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						if buffer[position] != rune('N') {
							goto l838
						}
						position++
					}
				l843:
					{
						position845, tokenIndex845, depth845 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l846
						}
						position++
						goto l845
					l846:
						position, tokenIndex, depth = position845, tokenIndex845, depth845
						if buffer[position] != rune('D') {
							goto l838
						}
						position++
					}
				l845:
					depth--
					add(rulePegText, position840)
				}
				if !_rules[ruleAction63]() {
					goto l838
				}
				depth--
				add(ruleAnd, position839)
			}
			return true
		l838:
			position, tokenIndex, depth = position838, tokenIndex838, depth838
			return false
		},
		/* 82 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action64)> */
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
						if buffer[position] != rune('n') {
							goto l851
						}
						position++
						goto l850
					l851:
						position, tokenIndex, depth = position850, tokenIndex850, depth850
						if buffer[position] != rune('N') {
							goto l847
						}
						position++
					}
				l850:
					{
						position852, tokenIndex852, depth852 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l853
						}
						position++
						goto l852
					l853:
						position, tokenIndex, depth = position852, tokenIndex852, depth852
						if buffer[position] != rune('O') {
							goto l847
						}
						position++
					}
				l852:
					{
						position854, tokenIndex854, depth854 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l855
						}
						position++
						goto l854
					l855:
						position, tokenIndex, depth = position854, tokenIndex854, depth854
						if buffer[position] != rune('T') {
							goto l847
						}
						position++
					}
				l854:
					depth--
					add(rulePegText, position849)
				}
				if !_rules[ruleAction64]() {
					goto l847
				}
				depth--
				add(ruleNot, position848)
			}
			return true
		l847:
			position, tokenIndex, depth = position847, tokenIndex847, depth847
			return false
		},
		/* 83 Equal <- <(<'='> Action65)> */
		func() bool {
			position856, tokenIndex856, depth856 := position, tokenIndex, depth
			{
				position857 := position
				depth++
				{
					position858 := position
					depth++
					if buffer[position] != rune('=') {
						goto l856
					}
					position++
					depth--
					add(rulePegText, position858)
				}
				if !_rules[ruleAction65]() {
					goto l856
				}
				depth--
				add(ruleEqual, position857)
			}
			return true
		l856:
			position, tokenIndex, depth = position856, tokenIndex856, depth856
			return false
		},
		/* 84 Less <- <(<'<'> Action66)> */
		func() bool {
			position859, tokenIndex859, depth859 := position, tokenIndex, depth
			{
				position860 := position
				depth++
				{
					position861 := position
					depth++
					if buffer[position] != rune('<') {
						goto l859
					}
					position++
					depth--
					add(rulePegText, position861)
				}
				if !_rules[ruleAction66]() {
					goto l859
				}
				depth--
				add(ruleLess, position860)
			}
			return true
		l859:
			position, tokenIndex, depth = position859, tokenIndex859, depth859
			return false
		},
		/* 85 LessOrEqual <- <(<('<' '=')> Action67)> */
		func() bool {
			position862, tokenIndex862, depth862 := position, tokenIndex, depth
			{
				position863 := position
				depth++
				{
					position864 := position
					depth++
					if buffer[position] != rune('<') {
						goto l862
					}
					position++
					if buffer[position] != rune('=') {
						goto l862
					}
					position++
					depth--
					add(rulePegText, position864)
				}
				if !_rules[ruleAction67]() {
					goto l862
				}
				depth--
				add(ruleLessOrEqual, position863)
			}
			return true
		l862:
			position, tokenIndex, depth = position862, tokenIndex862, depth862
			return false
		},
		/* 86 Greater <- <(<'>'> Action68)> */
		func() bool {
			position865, tokenIndex865, depth865 := position, tokenIndex, depth
			{
				position866 := position
				depth++
				{
					position867 := position
					depth++
					if buffer[position] != rune('>') {
						goto l865
					}
					position++
					depth--
					add(rulePegText, position867)
				}
				if !_rules[ruleAction68]() {
					goto l865
				}
				depth--
				add(ruleGreater, position866)
			}
			return true
		l865:
			position, tokenIndex, depth = position865, tokenIndex865, depth865
			return false
		},
		/* 87 GreaterOrEqual <- <(<('>' '=')> Action69)> */
		func() bool {
			position868, tokenIndex868, depth868 := position, tokenIndex, depth
			{
				position869 := position
				depth++
				{
					position870 := position
					depth++
					if buffer[position] != rune('>') {
						goto l868
					}
					position++
					if buffer[position] != rune('=') {
						goto l868
					}
					position++
					depth--
					add(rulePegText, position870)
				}
				if !_rules[ruleAction69]() {
					goto l868
				}
				depth--
				add(ruleGreaterOrEqual, position869)
			}
			return true
		l868:
			position, tokenIndex, depth = position868, tokenIndex868, depth868
			return false
		},
		/* 88 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action70)> */
		func() bool {
			position871, tokenIndex871, depth871 := position, tokenIndex, depth
			{
				position872 := position
				depth++
				{
					position873 := position
					depth++
					{
						position874, tokenIndex874, depth874 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l875
						}
						position++
						if buffer[position] != rune('=') {
							goto l875
						}
						position++
						goto l874
					l875:
						position, tokenIndex, depth = position874, tokenIndex874, depth874
						if buffer[position] != rune('<') {
							goto l871
						}
						position++
						if buffer[position] != rune('>') {
							goto l871
						}
						position++
					}
				l874:
					depth--
					add(rulePegText, position873)
				}
				if !_rules[ruleAction70]() {
					goto l871
				}
				depth--
				add(ruleNotEqual, position872)
			}
			return true
		l871:
			position, tokenIndex, depth = position871, tokenIndex871, depth871
			return false
		},
		/* 89 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action71)> */
		func() bool {
			position876, tokenIndex876, depth876 := position, tokenIndex, depth
			{
				position877 := position
				depth++
				{
					position878 := position
					depth++
					{
						position879, tokenIndex879, depth879 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l880
						}
						position++
						goto l879
					l880:
						position, tokenIndex, depth = position879, tokenIndex879, depth879
						if buffer[position] != rune('I') {
							goto l876
						}
						position++
					}
				l879:
					{
						position881, tokenIndex881, depth881 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l882
						}
						position++
						goto l881
					l882:
						position, tokenIndex, depth = position881, tokenIndex881, depth881
						if buffer[position] != rune('S') {
							goto l876
						}
						position++
					}
				l881:
					depth--
					add(rulePegText, position878)
				}
				if !_rules[ruleAction71]() {
					goto l876
				}
				depth--
				add(ruleIs, position877)
			}
			return true
		l876:
			position, tokenIndex, depth = position876, tokenIndex876, depth876
			return false
		},
		/* 90 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action72)> */
		func() bool {
			position883, tokenIndex883, depth883 := position, tokenIndex, depth
			{
				position884 := position
				depth++
				{
					position885 := position
					depth++
					{
						position886, tokenIndex886, depth886 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l887
						}
						position++
						goto l886
					l887:
						position, tokenIndex, depth = position886, tokenIndex886, depth886
						if buffer[position] != rune('I') {
							goto l883
						}
						position++
					}
				l886:
					{
						position888, tokenIndex888, depth888 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l889
						}
						position++
						goto l888
					l889:
						position, tokenIndex, depth = position888, tokenIndex888, depth888
						if buffer[position] != rune('S') {
							goto l883
						}
						position++
					}
				l888:
					if !_rules[rulesp]() {
						goto l883
					}
					{
						position890, tokenIndex890, depth890 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l891
						}
						position++
						goto l890
					l891:
						position, tokenIndex, depth = position890, tokenIndex890, depth890
						if buffer[position] != rune('N') {
							goto l883
						}
						position++
					}
				l890:
					{
						position892, tokenIndex892, depth892 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l893
						}
						position++
						goto l892
					l893:
						position, tokenIndex, depth = position892, tokenIndex892, depth892
						if buffer[position] != rune('O') {
							goto l883
						}
						position++
					}
				l892:
					{
						position894, tokenIndex894, depth894 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l895
						}
						position++
						goto l894
					l895:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
						if buffer[position] != rune('T') {
							goto l883
						}
						position++
					}
				l894:
					depth--
					add(rulePegText, position885)
				}
				if !_rules[ruleAction72]() {
					goto l883
				}
				depth--
				add(ruleIsNot, position884)
			}
			return true
		l883:
			position, tokenIndex, depth = position883, tokenIndex883, depth883
			return false
		},
		/* 91 Plus <- <(<'+'> Action73)> */
		func() bool {
			position896, tokenIndex896, depth896 := position, tokenIndex, depth
			{
				position897 := position
				depth++
				{
					position898 := position
					depth++
					if buffer[position] != rune('+') {
						goto l896
					}
					position++
					depth--
					add(rulePegText, position898)
				}
				if !_rules[ruleAction73]() {
					goto l896
				}
				depth--
				add(rulePlus, position897)
			}
			return true
		l896:
			position, tokenIndex, depth = position896, tokenIndex896, depth896
			return false
		},
		/* 92 Minus <- <(<'-'> Action74)> */
		func() bool {
			position899, tokenIndex899, depth899 := position, tokenIndex, depth
			{
				position900 := position
				depth++
				{
					position901 := position
					depth++
					if buffer[position] != rune('-') {
						goto l899
					}
					position++
					depth--
					add(rulePegText, position901)
				}
				if !_rules[ruleAction74]() {
					goto l899
				}
				depth--
				add(ruleMinus, position900)
			}
			return true
		l899:
			position, tokenIndex, depth = position899, tokenIndex899, depth899
			return false
		},
		/* 93 Multiply <- <(<'*'> Action75)> */
		func() bool {
			position902, tokenIndex902, depth902 := position, tokenIndex, depth
			{
				position903 := position
				depth++
				{
					position904 := position
					depth++
					if buffer[position] != rune('*') {
						goto l902
					}
					position++
					depth--
					add(rulePegText, position904)
				}
				if !_rules[ruleAction75]() {
					goto l902
				}
				depth--
				add(ruleMultiply, position903)
			}
			return true
		l902:
			position, tokenIndex, depth = position902, tokenIndex902, depth902
			return false
		},
		/* 94 Divide <- <(<'/'> Action76)> */
		func() bool {
			position905, tokenIndex905, depth905 := position, tokenIndex, depth
			{
				position906 := position
				depth++
				{
					position907 := position
					depth++
					if buffer[position] != rune('/') {
						goto l905
					}
					position++
					depth--
					add(rulePegText, position907)
				}
				if !_rules[ruleAction76]() {
					goto l905
				}
				depth--
				add(ruleDivide, position906)
			}
			return true
		l905:
			position, tokenIndex, depth = position905, tokenIndex905, depth905
			return false
		},
		/* 95 Modulo <- <(<'%'> Action77)> */
		func() bool {
			position908, tokenIndex908, depth908 := position, tokenIndex, depth
			{
				position909 := position
				depth++
				{
					position910 := position
					depth++
					if buffer[position] != rune('%') {
						goto l908
					}
					position++
					depth--
					add(rulePegText, position910)
				}
				if !_rules[ruleAction77]() {
					goto l908
				}
				depth--
				add(ruleModulo, position909)
			}
			return true
		l908:
			position, tokenIndex, depth = position908, tokenIndex908, depth908
			return false
		},
		/* 96 UnaryMinus <- <(<'-'> Action78)> */
		func() bool {
			position911, tokenIndex911, depth911 := position, tokenIndex, depth
			{
				position912 := position
				depth++
				{
					position913 := position
					depth++
					if buffer[position] != rune('-') {
						goto l911
					}
					position++
					depth--
					add(rulePegText, position913)
				}
				if !_rules[ruleAction78]() {
					goto l911
				}
				depth--
				add(ruleUnaryMinus, position912)
			}
			return true
		l911:
			position, tokenIndex, depth = position911, tokenIndex911, depth911
			return false
		},
		/* 97 Identifier <- <(<ident> Action79)> */
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
				if !_rules[ruleAction79]() {
					goto l914
				}
				depth--
				add(ruleIdentifier, position915)
			}
			return true
		l914:
			position, tokenIndex, depth = position914, tokenIndex914, depth914
			return false
		},
		/* 98 TargetIdentifier <- <(<jsonPath> Action80)> */
		func() bool {
			position917, tokenIndex917, depth917 := position, tokenIndex, depth
			{
				position918 := position
				depth++
				{
					position919 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l917
					}
					depth--
					add(rulePegText, position919)
				}
				if !_rules[ruleAction80]() {
					goto l917
				}
				depth--
				add(ruleTargetIdentifier, position918)
			}
			return true
		l917:
			position, tokenIndex, depth = position917, tokenIndex917, depth917
			return false
		},
		/* 99 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position920, tokenIndex920, depth920 := position, tokenIndex, depth
			{
				position921 := position
				depth++
				{
					position922, tokenIndex922, depth922 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l923
					}
					position++
					goto l922
				l923:
					position, tokenIndex, depth = position922, tokenIndex922, depth922
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l920
					}
					position++
				}
			l922:
			l924:
				{
					position925, tokenIndex925, depth925 := position, tokenIndex, depth
					{
						position926, tokenIndex926, depth926 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l928
						}
						position++
						goto l926
					l928:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l929
						}
						position++
						goto l926
					l929:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('_') {
							goto l925
						}
						position++
					}
				l926:
					goto l924
				l925:
					position, tokenIndex, depth = position925, tokenIndex925, depth925
				}
				depth--
				add(ruleident, position921)
			}
			return true
		l920:
			position, tokenIndex, depth = position920, tokenIndex920, depth920
			return false
		},
		/* 100 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position930, tokenIndex930, depth930 := position, tokenIndex, depth
			{
				position931 := position
				depth++
				{
					position932, tokenIndex932, depth932 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l933
					}
					position++
					goto l932
				l933:
					position, tokenIndex, depth = position932, tokenIndex932, depth932
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l930
					}
					position++
				}
			l932:
			l934:
				{
					position935, tokenIndex935, depth935 := position, tokenIndex, depth
					{
						position936, tokenIndex936, depth936 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l937
						}
						position++
						goto l936
					l937:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l938
						}
						position++
						goto l936
					l938:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l939
						}
						position++
						goto l936
					l939:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if buffer[position] != rune('_') {
							goto l940
						}
						position++
						goto l936
					l940:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if buffer[position] != rune('.') {
							goto l941
						}
						position++
						goto l936
					l941:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if buffer[position] != rune('[') {
							goto l942
						}
						position++
						goto l936
					l942:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if buffer[position] != rune(']') {
							goto l943
						}
						position++
						goto l936
					l943:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
						if buffer[position] != rune('"') {
							goto l935
						}
						position++
					}
				l936:
					goto l934
				l935:
					position, tokenIndex, depth = position935, tokenIndex935, depth935
				}
				depth--
				add(rulejsonPath, position931)
			}
			return true
		l930:
			position, tokenIndex, depth = position930, tokenIndex930, depth930
			return false
		},
		/* 101 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position945 := position
				depth++
			l946:
				{
					position947, tokenIndex947, depth947 := position, tokenIndex, depth
					{
						position948, tokenIndex948, depth948 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l949
						}
						position++
						goto l948
					l949:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
						if buffer[position] != rune('\t') {
							goto l950
						}
						position++
						goto l948
					l950:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
						if buffer[position] != rune('\n') {
							goto l951
						}
						position++
						goto l948
					l951:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
						if buffer[position] != rune('\r') {
							goto l952
						}
						position++
						goto l948
					l952:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
						if !_rules[rulecomment]() {
							goto l947
						}
					}
				l948:
					goto l946
				l947:
					position, tokenIndex, depth = position947, tokenIndex947, depth947
				}
				depth--
				add(rulesp, position945)
			}
			return true
		},
		/* 102 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position953, tokenIndex953, depth953 := position, tokenIndex, depth
			{
				position954 := position
				depth++
				if buffer[position] != rune('-') {
					goto l953
				}
				position++
				if buffer[position] != rune('-') {
					goto l953
				}
				position++
			l955:
				{
					position956, tokenIndex956, depth956 := position, tokenIndex, depth
					{
						position957, tokenIndex957, depth957 := position, tokenIndex, depth
						{
							position958, tokenIndex958, depth958 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l959
							}
							position++
							goto l958
						l959:
							position, tokenIndex, depth = position958, tokenIndex958, depth958
							if buffer[position] != rune('\n') {
								goto l957
							}
							position++
						}
					l958:
						goto l956
					l957:
						position, tokenIndex, depth = position957, tokenIndex957, depth957
					}
					if !matchDot() {
						goto l956
					}
					goto l955
				l956:
					position, tokenIndex, depth = position956, tokenIndex956, depth956
				}
				{
					position960, tokenIndex960, depth960 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l961
					}
					position++
					goto l960
				l961:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if buffer[position] != rune('\n') {
						goto l953
					}
					position++
				}
			l960:
				depth--
				add(rulecomment, position954)
			}
			return true
		l953:
			position, tokenIndex, depth = position953, tokenIndex953, depth953
			return false
		},
		/* 104 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 105 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 106 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 107 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 108 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 109 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 110 Action6 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 111 Action7 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 112 Action8 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 113 Action9 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 114 Action10 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 115 Action11 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		nil,
		/* 117 Action12 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 118 Action13 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 119 Action14 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 120 Action15 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 121 Action16 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 122 Action17 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 123 Action18 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 124 Action19 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 125 Action20 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 126 Action21 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 127 Action22 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 128 Action23 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 129 Action24 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 130 Action25 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 131 Action26 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 132 Action27 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 133 Action28 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 134 Action29 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 135 Action30 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 136 Action31 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 137 Action32 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 138 Action33 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 139 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 140 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 141 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 142 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 143 Action38 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 144 Action39 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 145 Action40 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 146 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 147 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 148 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 149 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 150 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 151 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 152 Action47 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 153 Action48 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 154 Action49 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 155 Action50 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 156 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 157 Action52 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 158 Action53 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 159 Action54 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 160 Action55 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 161 Action56 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 162 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 163 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 164 Action59 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 165 Action60 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 166 Action61 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 167 Action62 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 168 Action63 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 169 Action64 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 170 Action65 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 171 Action66 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 172 Action67 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 173 Action68 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 174 Action69 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 175 Action70 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 176 Action71 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 177 Action72 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 178 Action73 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 179 Action74 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 180 Action75 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 181 Action76 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 182 Action77 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 183 Action78 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 184 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 185 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
	}
	p.rules = _rules
}

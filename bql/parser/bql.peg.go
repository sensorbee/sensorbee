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
	rulePegText
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
	ruleAction82

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
	"PegText",
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
	"Action82",

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
	rules  [190]func() bool
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

			p.AssembleDropState()

		case ruleAction14:

			p.AssembleEmitter(begin, end)

		case ruleAction15:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction16:

			p.PushComponent(end, end, NewStream("*"))
			p.AssembleStreamEmitInterval()

		case ruleAction17:

			p.AssembleStreamEmitInterval()

		case ruleAction18:

			p.AssembleProjections(begin, end)

		case ruleAction19:

			p.AssembleAlias()

		case ruleAction20:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction21:

			p.AssembleInterval()

		case ruleAction22:

			p.AssembleInterval()

		case ruleAction23:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction24:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction25:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction26:

			p.EnsureAliasedStreamWindow()

		case ruleAction27:

			p.AssembleAliasedStreamWindow()

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

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction38:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction39:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction40:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction41:

			p.AssembleFuncApp()

		case ruleAction42:

			p.AssembleExpressions(begin, end)

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction44:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction49:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction50:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction51:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction52:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction53:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction54:

			p.PushComponent(begin, end, Istream)

		case ruleAction55:

			p.PushComponent(begin, end, Dstream)

		case ruleAction56:

			p.PushComponent(begin, end, Rstream)

		case ruleAction57:

			p.PushComponent(begin, end, Tuples)

		case ruleAction58:

			p.PushComponent(begin, end, Seconds)

		case ruleAction59:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction60:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction61:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction62:

			p.PushComponent(begin, end, Yes)

		case ruleAction63:

			p.PushComponent(begin, end, No)

		case ruleAction64:

			p.PushComponent(begin, end, Or)

		case ruleAction65:

			p.PushComponent(begin, end, And)

		case ruleAction66:

			p.PushComponent(begin, end, Not)

		case ruleAction67:

			p.PushComponent(begin, end, Equal)

		case ruleAction68:

			p.PushComponent(begin, end, Less)

		case ruleAction69:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction70:

			p.PushComponent(begin, end, Greater)

		case ruleAction71:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction72:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction73:

			p.PushComponent(begin, end, Is)

		case ruleAction74:

			p.PushComponent(begin, end, IsNot)

		case ruleAction75:

			p.PushComponent(begin, end, Plus)

		case ruleAction76:

			p.PushComponent(begin, end, Minus)

		case ruleAction77:

			p.PushComponent(begin, end, Multiply)

		case ruleAction78:

			p.PushComponent(begin, end, Divide)

		case ruleAction79:

			p.PushComponent(begin, end, Modulo)

		case ruleAction80:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction81:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction82:

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
		/* 1 Statement <- <(SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateSinkStmt / InsertIntoSelectStmt / InsertIntoFromStmt / CreateStateStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt / DropSourceStmt / DropStreamStmt / DropSinkStmt / DropStateStmt)> */
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
			position23, tokenIndex23, depth23 := position, tokenIndex, depth
			{
				position24 := position
				depth++
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l26
					}
					position++
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if buffer[position] != rune('S') {
						goto l23
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
						goto l23
					}
					position++
				}
			l27:
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l30
					}
					position++
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if buffer[position] != rune('L') {
						goto l23
					}
					position++
				}
			l29:
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l32
					}
					position++
					goto l31
				l32:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
					if buffer[position] != rune('E') {
						goto l23
					}
					position++
				}
			l31:
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l34
					}
					position++
					goto l33
				l34:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
					if buffer[position] != rune('C') {
						goto l23
					}
					position++
				}
			l33:
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l36
					}
					position++
					goto l35
				l36:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if buffer[position] != rune('T') {
						goto l23
					}
					position++
				}
			l35:
				if !_rules[rulesp]() {
					goto l23
				}
				if !_rules[ruleEmitter]() {
					goto l23
				}
				if !_rules[rulesp]() {
					goto l23
				}
				if !_rules[ruleProjections]() {
					goto l23
				}
				if !_rules[rulesp]() {
					goto l23
				}
				if !_rules[ruleWindowedFrom]() {
					goto l23
				}
				if !_rules[rulesp]() {
					goto l23
				}
				if !_rules[ruleFilter]() {
					goto l23
				}
				if !_rules[rulesp]() {
					goto l23
				}
				if !_rules[ruleGrouping]() {
					goto l23
				}
				if !_rules[rulesp]() {
					goto l23
				}
				if !_rules[ruleHaving]() {
					goto l23
				}
				if !_rules[rulesp]() {
					goto l23
				}
				if !_rules[ruleAction0]() {
					goto l23
				}
				depth--
				add(ruleSelectStmt, position24)
			}
			return true
		l23:
			position, tokenIndex, depth = position23, tokenIndex23, depth23
			return false
		},
		/* 3 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position37, tokenIndex37, depth37 := position, tokenIndex, depth
			{
				position38 := position
				depth++
				{
					position39, tokenIndex39, depth39 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l40
					}
					position++
					goto l39
				l40:
					position, tokenIndex, depth = position39, tokenIndex39, depth39
					if buffer[position] != rune('C') {
						goto l37
					}
					position++
				}
			l39:
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l42
					}
					position++
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if buffer[position] != rune('R') {
						goto l37
					}
					position++
				}
			l41:
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l44
					}
					position++
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if buffer[position] != rune('E') {
						goto l37
					}
					position++
				}
			l43:
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l46
					}
					position++
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if buffer[position] != rune('A') {
						goto l37
					}
					position++
				}
			l45:
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l48
					}
					position++
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					if buffer[position] != rune('T') {
						goto l37
					}
					position++
				}
			l47:
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l50
					}
					position++
					goto l49
				l50:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
					if buffer[position] != rune('E') {
						goto l37
					}
					position++
				}
			l49:
				if !_rules[rulesp]() {
					goto l37
				}
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l52
					}
					position++
					goto l51
				l52:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
					if buffer[position] != rune('S') {
						goto l37
					}
					position++
				}
			l51:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					if buffer[position] != rune('T') {
						goto l37
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
					if buffer[position] != rune('R') {
						goto l37
					}
					position++
				}
			l55:
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
					if buffer[position] != rune('E') {
						goto l37
					}
					position++
				}
			l57:
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
						goto l37
					}
					position++
				}
			l59:
				{
					position61, tokenIndex61, depth61 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l62
					}
					position++
					goto l61
				l62:
					position, tokenIndex, depth = position61, tokenIndex61, depth61
					if buffer[position] != rune('M') {
						goto l37
					}
					position++
				}
			l61:
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l37
				}
				if !_rules[rulesp]() {
					goto l37
				}
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l64
					}
					position++
					goto l63
				l64:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
					if buffer[position] != rune('A') {
						goto l37
					}
					position++
				}
			l63:
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
						goto l37
					}
					position++
				}
			l65:
				if !_rules[rulesp]() {
					goto l37
				}
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
						goto l37
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
						goto l37
					}
					position++
				}
			l69:
				{
					position71, tokenIndex71, depth71 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l72
					}
					position++
					goto l71
				l72:
					position, tokenIndex, depth = position71, tokenIndex71, depth71
					if buffer[position] != rune('L') {
						goto l37
					}
					position++
				}
			l71:
				{
					position73, tokenIndex73, depth73 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l74
					}
					position++
					goto l73
				l74:
					position, tokenIndex, depth = position73, tokenIndex73, depth73
					if buffer[position] != rune('E') {
						goto l37
					}
					position++
				}
			l73:
				{
					position75, tokenIndex75, depth75 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l76
					}
					position++
					goto l75
				l76:
					position, tokenIndex, depth = position75, tokenIndex75, depth75
					if buffer[position] != rune('C') {
						goto l37
					}
					position++
				}
			l75:
				{
					position77, tokenIndex77, depth77 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l78
					}
					position++
					goto l77
				l78:
					position, tokenIndex, depth = position77, tokenIndex77, depth77
					if buffer[position] != rune('T') {
						goto l37
					}
					position++
				}
			l77:
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleEmitter]() {
					goto l37
				}
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleProjections]() {
					goto l37
				}
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleWindowedFrom]() {
					goto l37
				}
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleFilter]() {
					goto l37
				}
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleGrouping]() {
					goto l37
				}
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleHaving]() {
					goto l37
				}
				if !_rules[rulesp]() {
					goto l37
				}
				if !_rules[ruleAction1]() {
					goto l37
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position38)
			}
			return true
		l37:
			position, tokenIndex, depth = position37, tokenIndex37, depth37
			return false
		},
		/* 4 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position79, tokenIndex79, depth79 := position, tokenIndex, depth
			{
				position80 := position
				depth++
				{
					position81, tokenIndex81, depth81 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l82
					}
					position++
					goto l81
				l82:
					position, tokenIndex, depth = position81, tokenIndex81, depth81
					if buffer[position] != rune('C') {
						goto l79
					}
					position++
				}
			l81:
				{
					position83, tokenIndex83, depth83 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l84
					}
					position++
					goto l83
				l84:
					position, tokenIndex, depth = position83, tokenIndex83, depth83
					if buffer[position] != rune('R') {
						goto l79
					}
					position++
				}
			l83:
				{
					position85, tokenIndex85, depth85 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l86
					}
					position++
					goto l85
				l86:
					position, tokenIndex, depth = position85, tokenIndex85, depth85
					if buffer[position] != rune('E') {
						goto l79
					}
					position++
				}
			l85:
				{
					position87, tokenIndex87, depth87 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l88
					}
					position++
					goto l87
				l88:
					position, tokenIndex, depth = position87, tokenIndex87, depth87
					if buffer[position] != rune('A') {
						goto l79
					}
					position++
				}
			l87:
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l90
					}
					position++
					goto l89
				l90:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
					if buffer[position] != rune('T') {
						goto l79
					}
					position++
				}
			l89:
				{
					position91, tokenIndex91, depth91 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l92
					}
					position++
					goto l91
				l92:
					position, tokenIndex, depth = position91, tokenIndex91, depth91
					if buffer[position] != rune('E') {
						goto l79
					}
					position++
				}
			l91:
				if !_rules[rulesp]() {
					goto l79
				}
				if !_rules[rulePausedOpt]() {
					goto l79
				}
				if !_rules[rulesp]() {
					goto l79
				}
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('S') {
						goto l79
					}
					position++
				}
			l93:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if buffer[position] != rune('O') {
						goto l79
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('U') {
						goto l79
					}
					position++
				}
			l97:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('R') {
						goto l79
					}
					position++
				}
			l99:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('C') {
						goto l79
					}
					position++
				}
			l101:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if buffer[position] != rune('E') {
						goto l79
					}
					position++
				}
			l103:
				if !_rules[rulesp]() {
					goto l79
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l79
				}
				if !_rules[rulesp]() {
					goto l79
				}
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('T') {
						goto l79
					}
					position++
				}
			l105:
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l108
					}
					position++
					goto l107
				l108:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
					if buffer[position] != rune('Y') {
						goto l79
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('P') {
						goto l79
					}
					position++
				}
			l109:
				{
					position111, tokenIndex111, depth111 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l112
					}
					position++
					goto l111
				l112:
					position, tokenIndex, depth = position111, tokenIndex111, depth111
					if buffer[position] != rune('E') {
						goto l79
					}
					position++
				}
			l111:
				if !_rules[rulesp]() {
					goto l79
				}
				if !_rules[ruleSourceSinkType]() {
					goto l79
				}
				if !_rules[rulesp]() {
					goto l79
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l79
				}
				if !_rules[ruleAction2]() {
					goto l79
				}
				depth--
				add(ruleCreateSourceStmt, position80)
			}
			return true
		l79:
			position, tokenIndex, depth = position79, tokenIndex79, depth79
			return false
		},
		/* 5 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position113, tokenIndex113, depth113 := position, tokenIndex, depth
			{
				position114 := position
				depth++
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('C') {
						goto l113
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('R') {
						goto l113
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('E') {
						goto l113
					}
					position++
				}
			l119:
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if buffer[position] != rune('A') {
						goto l113
					}
					position++
				}
			l121:
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('T') {
						goto l113
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('E') {
						goto l113
					}
					position++
				}
			l125:
				if !_rules[rulesp]() {
					goto l113
				}
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('S') {
						goto l113
					}
					position++
				}
			l127:
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('I') {
						goto l113
					}
					position++
				}
			l129:
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('N') {
						goto l113
					}
					position++
				}
			l131:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('K') {
						goto l113
					}
					position++
				}
			l133:
				if !_rules[rulesp]() {
					goto l113
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l113
				}
				if !_rules[rulesp]() {
					goto l113
				}
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('T') {
						goto l113
					}
					position++
				}
			l135:
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l138
					}
					position++
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if buffer[position] != rune('Y') {
						goto l113
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('P') {
						goto l113
					}
					position++
				}
			l139:
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l142
					}
					position++
					goto l141
				l142:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					if buffer[position] != rune('E') {
						goto l113
					}
					position++
				}
			l141:
				if !_rules[rulesp]() {
					goto l113
				}
				if !_rules[ruleSourceSinkType]() {
					goto l113
				}
				if !_rules[rulesp]() {
					goto l113
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l113
				}
				if !_rules[ruleAction3]() {
					goto l113
				}
				depth--
				add(ruleCreateSinkStmt, position114)
			}
			return true
		l113:
			position, tokenIndex, depth = position113, tokenIndex113, depth113
			return false
		},
		/* 6 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position143, tokenIndex143, depth143 := position, tokenIndex, depth
			{
				position144 := position
				depth++
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('C') {
						goto l143
					}
					position++
				}
			l145:
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('R') {
						goto l143
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('E') {
						goto l143
					}
					position++
				}
			l149:
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('A') {
						goto l143
					}
					position++
				}
			l151:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('T') {
						goto l143
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('E') {
						goto l143
					}
					position++
				}
			l155:
				if !_rules[rulesp]() {
					goto l143
				}
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('S') {
						goto l143
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
						goto l143
					}
					position++
				}
			l159:
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l162
					}
					position++
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if buffer[position] != rune('A') {
						goto l143
					}
					position++
				}
			l161:
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
						goto l143
					}
					position++
				}
			l163:
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('E') {
						goto l143
					}
					position++
				}
			l165:
				if !_rules[rulesp]() {
					goto l143
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l143
				}
				if !_rules[rulesp]() {
					goto l143
				}
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('T') {
						goto l143
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('Y') {
						goto l143
					}
					position++
				}
			l169:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('P') {
						goto l143
					}
					position++
				}
			l171:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('E') {
						goto l143
					}
					position++
				}
			l173:
				if !_rules[rulesp]() {
					goto l143
				}
				if !_rules[ruleSourceSinkType]() {
					goto l143
				}
				if !_rules[rulesp]() {
					goto l143
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l143
				}
				if !_rules[ruleAction4]() {
					goto l143
				}
				depth--
				add(ruleCreateStateStmt, position144)
			}
			return true
		l143:
			position, tokenIndex, depth = position143, tokenIndex143, depth143
			return false
		},
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action5)> */
		func() bool {
			position175, tokenIndex175, depth175 := position, tokenIndex, depth
			{
				position176 := position
				depth++
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('I') {
						goto l175
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('N') {
						goto l175
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('S') {
						goto l175
					}
					position++
				}
			l181:
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('E') {
						goto l175
					}
					position++
				}
			l183:
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('R') {
						goto l175
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
						goto l175
					}
					position++
				}
			l187:
				if !_rules[rulesp]() {
					goto l175
				}
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('I') {
						goto l175
					}
					position++
				}
			l189:
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('N') {
						goto l175
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
						goto l175
					}
					position++
				}
			l193:
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('O') {
						goto l175
					}
					position++
				}
			l195:
				if !_rules[rulesp]() {
					goto l175
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l175
				}
				if !_rules[rulesp]() {
					goto l175
				}
				if !_rules[ruleSelectStmt]() {
					goto l175
				}
				if !_rules[ruleAction5]() {
					goto l175
				}
				depth--
				add(ruleInsertIntoSelectStmt, position176)
			}
			return true
		l175:
			position, tokenIndex, depth = position175, tokenIndex175, depth175
			return false
		},
		/* 8 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action6)> */
		func() bool {
			position197, tokenIndex197, depth197 := position, tokenIndex, depth
			{
				position198 := position
				depth++
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('I') {
						goto l197
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('N') {
						goto l197
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('S') {
						goto l197
					}
					position++
				}
			l203:
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('E') {
						goto l197
					}
					position++
				}
			l205:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('R') {
						goto l197
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('T') {
						goto l197
					}
					position++
				}
			l209:
				if !_rules[rulesp]() {
					goto l197
				}
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('I') {
						goto l197
					}
					position++
				}
			l211:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('N') {
						goto l197
					}
					position++
				}
			l213:
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('T') {
						goto l197
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
						goto l197
					}
					position++
				}
			l217:
				if !_rules[rulesp]() {
					goto l197
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l197
				}
				if !_rules[rulesp]() {
					goto l197
				}
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('F') {
						goto l197
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
						goto l197
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('O') {
						goto l197
					}
					position++
				}
			l223:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('M') {
						goto l197
					}
					position++
				}
			l225:
				if !_rules[rulesp]() {
					goto l197
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l197
				}
				if !_rules[ruleAction6]() {
					goto l197
				}
				depth--
				add(ruleInsertIntoFromStmt, position198)
			}
			return true
		l197:
			position, tokenIndex, depth = position197, tokenIndex197, depth197
			return false
		},
		/* 9 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action7)> */
		func() bool {
			position227, tokenIndex227, depth227 := position, tokenIndex, depth
			{
				position228 := position
				depth++
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('P') {
						goto l227
					}
					position++
				}
			l229:
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if buffer[position] != rune('A') {
						goto l227
					}
					position++
				}
			l231:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('U') {
						goto l227
					}
					position++
				}
			l233:
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
						goto l227
					}
					position++
				}
			l235:
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if buffer[position] != rune('E') {
						goto l227
					}
					position++
				}
			l237:
				if !_rules[rulesp]() {
					goto l227
				}
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					if buffer[position] != rune('S') {
						goto l227
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('O') {
						goto l227
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('U') {
						goto l227
					}
					position++
				}
			l243:
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l246
					}
					position++
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if buffer[position] != rune('R') {
						goto l227
					}
					position++
				}
			l245:
				{
					position247, tokenIndex247, depth247 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l248
					}
					position++
					goto l247
				l248:
					position, tokenIndex, depth = position247, tokenIndex247, depth247
					if buffer[position] != rune('C') {
						goto l227
					}
					position++
				}
			l247:
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('E') {
						goto l227
					}
					position++
				}
			l249:
				if !_rules[rulesp]() {
					goto l227
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l227
				}
				if !_rules[ruleAction7]() {
					goto l227
				}
				depth--
				add(rulePauseSourceStmt, position228)
			}
			return true
		l227:
			position, tokenIndex, depth = position227, tokenIndex227, depth227
			return false
		},
		/* 10 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action8)> */
		func() bool {
			position251, tokenIndex251, depth251 := position, tokenIndex, depth
			{
				position252 := position
				depth++
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('R') {
						goto l251
					}
					position++
				}
			l253:
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l256
					}
					position++
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('E') {
						goto l251
					}
					position++
				}
			l255:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('S') {
						goto l251
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('U') {
						goto l251
					}
					position++
				}
			l259:
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('M') {
						goto l251
					}
					position++
				}
			l261:
				{
					position263, tokenIndex263, depth263 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if buffer[position] != rune('E') {
						goto l251
					}
					position++
				}
			l263:
				if !_rules[rulesp]() {
					goto l251
				}
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('S') {
						goto l251
					}
					position++
				}
			l265:
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l268
					}
					position++
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if buffer[position] != rune('O') {
						goto l251
					}
					position++
				}
			l267:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('U') {
						goto l251
					}
					position++
				}
			l269:
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('R') {
						goto l251
					}
					position++
				}
			l271:
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('C') {
						goto l251
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('E') {
						goto l251
					}
					position++
				}
			l275:
				if !_rules[rulesp]() {
					goto l251
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l251
				}
				if !_rules[ruleAction8]() {
					goto l251
				}
				depth--
				add(ruleResumeSourceStmt, position252)
			}
			return true
		l251:
			position, tokenIndex, depth = position251, tokenIndex251, depth251
			return false
		},
		/* 11 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action9)> */
		func() bool {
			position277, tokenIndex277, depth277 := position, tokenIndex, depth
			{
				position278 := position
				depth++
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('R') {
						goto l277
					}
					position++
				}
			l279:
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('E') {
						goto l277
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('W') {
						goto l277
					}
					position++
				}
			l283:
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('I') {
						goto l277
					}
					position++
				}
			l285:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('N') {
						goto l277
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('D') {
						goto l277
					}
					position++
				}
			l289:
				if !_rules[rulesp]() {
					goto l277
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
						goto l277
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
						goto l277
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
						goto l277
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
						goto l277
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
						goto l277
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
						goto l277
					}
					position++
				}
			l301:
				if !_rules[rulesp]() {
					goto l277
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l277
				}
				if !_rules[ruleAction9]() {
					goto l277
				}
				depth--
				add(ruleRewindSourceStmt, position278)
			}
			return true
		l277:
			position, tokenIndex, depth = position277, tokenIndex277, depth277
			return false
		},
		/* 12 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l306
					}
					position++
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if buffer[position] != rune('D') {
						goto l303
					}
					position++
				}
			l305:
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l308
					}
					position++
					goto l307
				l308:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
					if buffer[position] != rune('R') {
						goto l303
					}
					position++
				}
			l307:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('O') {
						goto l303
					}
					position++
				}
			l309:
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l312
					}
					position++
					goto l311
				l312:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
					if buffer[position] != rune('P') {
						goto l303
					}
					position++
				}
			l311:
				if !_rules[rulesp]() {
					goto l303
				}
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('S') {
						goto l303
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('O') {
						goto l303
					}
					position++
				}
			l315:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('U') {
						goto l303
					}
					position++
				}
			l317:
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l320
					}
					position++
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if buffer[position] != rune('R') {
						goto l303
					}
					position++
				}
			l319:
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l322
					}
					position++
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if buffer[position] != rune('C') {
						goto l303
					}
					position++
				}
			l321:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l324
					}
					position++
					goto l323
				l324:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					if buffer[position] != rune('E') {
						goto l303
					}
					position++
				}
			l323:
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
				add(ruleDropSourceStmt, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 13 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action11)> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l328
					}
					position++
					goto l327
				l328:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if buffer[position] != rune('D') {
						goto l325
					}
					position++
				}
			l327:
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l330
					}
					position++
					goto l329
				l330:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
					if buffer[position] != rune('R') {
						goto l325
					}
					position++
				}
			l329:
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('O') {
						goto l325
					}
					position++
				}
			l331:
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('P') {
						goto l325
					}
					position++
				}
			l333:
				if !_rules[rulesp]() {
					goto l325
				}
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l336
					}
					position++
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if buffer[position] != rune('S') {
						goto l325
					}
					position++
				}
			l335:
				{
					position337, tokenIndex337, depth337 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l338
					}
					position++
					goto l337
				l338:
					position, tokenIndex, depth = position337, tokenIndex337, depth337
					if buffer[position] != rune('T') {
						goto l325
					}
					position++
				}
			l337:
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l340
					}
					position++
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if buffer[position] != rune('R') {
						goto l325
					}
					position++
				}
			l339:
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l342
					}
					position++
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					if buffer[position] != rune('E') {
						goto l325
					}
					position++
				}
			l341:
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('A') {
						goto l325
					}
					position++
				}
			l343:
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l346
					}
					position++
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if buffer[position] != rune('M') {
						goto l325
					}
					position++
				}
			l345:
				if !_rules[rulesp]() {
					goto l325
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l325
				}
				if !_rules[ruleAction11]() {
					goto l325
				}
				depth--
				add(ruleDropStreamStmt, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
			return false
		},
		/* 14 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action12)> */
		func() bool {
			position347, tokenIndex347, depth347 := position, tokenIndex, depth
			{
				position348 := position
				depth++
				{
					position349, tokenIndex349, depth349 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l350
					}
					position++
					goto l349
				l350:
					position, tokenIndex, depth = position349, tokenIndex349, depth349
					if buffer[position] != rune('D') {
						goto l347
					}
					position++
				}
			l349:
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l352
					}
					position++
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					if buffer[position] != rune('R') {
						goto l347
					}
					position++
				}
			l351:
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l354
					}
					position++
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if buffer[position] != rune('O') {
						goto l347
					}
					position++
				}
			l353:
				{
					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l356
					}
					position++
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune('P') {
						goto l347
					}
					position++
				}
			l355:
				if !_rules[rulesp]() {
					goto l347
				}
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l358
					}
					position++
					goto l357
				l358:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
					if buffer[position] != rune('S') {
						goto l347
					}
					position++
				}
			l357:
				{
					position359, tokenIndex359, depth359 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l360
					}
					position++
					goto l359
				l360:
					position, tokenIndex, depth = position359, tokenIndex359, depth359
					if buffer[position] != rune('I') {
						goto l347
					}
					position++
				}
			l359:
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l362
					}
					position++
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if buffer[position] != rune('N') {
						goto l347
					}
					position++
				}
			l361:
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l364
					}
					position++
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					if buffer[position] != rune('K') {
						goto l347
					}
					position++
				}
			l363:
				if !_rules[rulesp]() {
					goto l347
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l347
				}
				if !_rules[ruleAction12]() {
					goto l347
				}
				depth--
				add(ruleDropSinkStmt, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 15 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position365, tokenIndex365, depth365 := position, tokenIndex, depth
			{
				position366 := position
				depth++
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					if buffer[position] != rune('D') {
						goto l365
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('R') {
						goto l365
					}
					position++
				}
			l369:
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l372
					}
					position++
					goto l371
				l372:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
					if buffer[position] != rune('O') {
						goto l365
					}
					position++
				}
			l371:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('P') {
						goto l365
					}
					position++
				}
			l373:
				if !_rules[rulesp]() {
					goto l365
				}
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
						goto l365
					}
					position++
				}
			l375:
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l378
					}
					position++
					goto l377
				l378:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
					if buffer[position] != rune('T') {
						goto l365
					}
					position++
				}
			l377:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('A') {
						goto l365
					}
					position++
				}
			l379:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('T') {
						goto l365
					}
					position++
				}
			l381:
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('E') {
						goto l365
					}
					position++
				}
			l383:
				if !_rules[rulesp]() {
					goto l365
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l365
				}
				if !_rules[ruleAction13]() {
					goto l365
				}
				depth--
				add(ruleDropStateStmt, position366)
			}
			return true
		l365:
			position, tokenIndex, depth = position365, tokenIndex365, depth365
			return false
		},
		/* 16 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) <(sp '[' sp (('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y')) sp EmitterIntervals sp ']')?> Action14)> */
		func() bool {
			position385, tokenIndex385, depth385 := position, tokenIndex, depth
			{
				position386 := position
				depth++
				{
					position387, tokenIndex387, depth387 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l388
					}
					goto l387
				l388:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
					if !_rules[ruleDSTREAM]() {
						goto l389
					}
					goto l387
				l389:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
					if !_rules[ruleRSTREAM]() {
						goto l385
					}
				}
			l387:
				{
					position390 := position
					depth++
					{
						position391, tokenIndex391, depth391 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l391
						}
						if buffer[position] != rune('[') {
							goto l391
						}
						position++
						if !_rules[rulesp]() {
							goto l391
						}
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
								goto l391
							}
							position++
						}
					l393:
						{
							position395, tokenIndex395, depth395 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l396
							}
							position++
							goto l395
						l396:
							position, tokenIndex, depth = position395, tokenIndex395, depth395
							if buffer[position] != rune('V') {
								goto l391
							}
							position++
						}
					l395:
						{
							position397, tokenIndex397, depth397 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l398
							}
							position++
							goto l397
						l398:
							position, tokenIndex, depth = position397, tokenIndex397, depth397
							if buffer[position] != rune('E') {
								goto l391
							}
							position++
						}
					l397:
						{
							position399, tokenIndex399, depth399 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l400
							}
							position++
							goto l399
						l400:
							position, tokenIndex, depth = position399, tokenIndex399, depth399
							if buffer[position] != rune('R') {
								goto l391
							}
							position++
						}
					l399:
						{
							position401, tokenIndex401, depth401 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l402
							}
							position++
							goto l401
						l402:
							position, tokenIndex, depth = position401, tokenIndex401, depth401
							if buffer[position] != rune('Y') {
								goto l391
							}
							position++
						}
					l401:
						if !_rules[rulesp]() {
							goto l391
						}
						if !_rules[ruleEmitterIntervals]() {
							goto l391
						}
						if !_rules[rulesp]() {
							goto l391
						}
						if buffer[position] != rune(']') {
							goto l391
						}
						position++
						goto l392
					l391:
						position, tokenIndex, depth = position391, tokenIndex391, depth391
					}
				l392:
					depth--
					add(rulePegText, position390)
				}
				if !_rules[ruleAction14]() {
					goto l385
				}
				depth--
				add(ruleEmitter, position386)
			}
			return true
		l385:
			position, tokenIndex, depth = position385, tokenIndex385, depth385
			return false
		},
		/* 17 EmitterIntervals <- <((TupleEmitterFromInterval (sp ',' sp TupleEmitterFromInterval)*) / TimeEmitterInterval / TupleEmitterInterval)> */
		func() bool {
			position403, tokenIndex403, depth403 := position, tokenIndex, depth
			{
				position404 := position
				depth++
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					if !_rules[ruleTupleEmitterFromInterval]() {
						goto l406
					}
				l407:
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l408
						}
						if buffer[position] != rune(',') {
							goto l408
						}
						position++
						if !_rules[rulesp]() {
							goto l408
						}
						if !_rules[ruleTupleEmitterFromInterval]() {
							goto l408
						}
						goto l407
					l408:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
					}
					goto l405
				l406:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if !_rules[ruleTimeEmitterInterval]() {
						goto l409
					}
					goto l405
				l409:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if !_rules[ruleTupleEmitterInterval]() {
						goto l403
					}
				}
			l405:
				depth--
				add(ruleEmitterIntervals, position404)
			}
			return true
		l403:
			position, tokenIndex, depth = position403, tokenIndex403, depth403
			return false
		},
		/* 18 TimeEmitterInterval <- <(<TimeInterval> Action15)> */
		func() bool {
			position410, tokenIndex410, depth410 := position, tokenIndex, depth
			{
				position411 := position
				depth++
				{
					position412 := position
					depth++
					if !_rules[ruleTimeInterval]() {
						goto l410
					}
					depth--
					add(rulePegText, position412)
				}
				if !_rules[ruleAction15]() {
					goto l410
				}
				depth--
				add(ruleTimeEmitterInterval, position411)
			}
			return true
		l410:
			position, tokenIndex, depth = position410, tokenIndex410, depth410
			return false
		},
		/* 19 TupleEmitterInterval <- <(<TuplesInterval> Action16)> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				{
					position415 := position
					depth++
					if !_rules[ruleTuplesInterval]() {
						goto l413
					}
					depth--
					add(rulePegText, position415)
				}
				if !_rules[ruleAction16]() {
					goto l413
				}
				depth--
				add(ruleTupleEmitterInterval, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 20 TupleEmitterFromInterval <- <(TuplesInterval sp (('i' / 'I') ('n' / 'N')) sp Stream Action17)> */
		func() bool {
			position416, tokenIndex416, depth416 := position, tokenIndex, depth
			{
				position417 := position
				depth++
				if !_rules[ruleTuplesInterval]() {
					goto l416
				}
				if !_rules[rulesp]() {
					goto l416
				}
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
					if buffer[position] != rune('I') {
						goto l416
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('N') {
						goto l416
					}
					position++
				}
			l420:
				if !_rules[rulesp]() {
					goto l416
				}
				if !_rules[ruleStream]() {
					goto l416
				}
				if !_rules[ruleAction17]() {
					goto l416
				}
				depth--
				add(ruleTupleEmitterFromInterval, position417)
			}
			return true
		l416:
			position, tokenIndex, depth = position416, tokenIndex416, depth416
			return false
		},
		/* 21 Projections <- <(<(Projection sp (',' sp Projection)*)> Action18)> */
		func() bool {
			position422, tokenIndex422, depth422 := position, tokenIndex, depth
			{
				position423 := position
				depth++
				{
					position424 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l422
					}
					if !_rules[rulesp]() {
						goto l422
					}
				l425:
					{
						position426, tokenIndex426, depth426 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l426
						}
						position++
						if !_rules[rulesp]() {
							goto l426
						}
						if !_rules[ruleProjection]() {
							goto l426
						}
						goto l425
					l426:
						position, tokenIndex, depth = position426, tokenIndex426, depth426
					}
					depth--
					add(rulePegText, position424)
				}
				if !_rules[ruleAction18]() {
					goto l422
				}
				depth--
				add(ruleProjections, position423)
			}
			return true
		l422:
			position, tokenIndex, depth = position422, tokenIndex422, depth422
			return false
		},
		/* 22 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position427, tokenIndex427, depth427 := position, tokenIndex, depth
			{
				position428 := position
				depth++
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l430
					}
					goto l429
				l430:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if !_rules[ruleExpression]() {
						goto l431
					}
					goto l429
				l431:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if !_rules[ruleWildcard]() {
						goto l427
					}
				}
			l429:
				depth--
				add(ruleProjection, position428)
			}
			return true
		l427:
			position, tokenIndex, depth = position427, tokenIndex427, depth427
			return false
		},
		/* 23 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action19)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l435
					}
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if !_rules[ruleWildcard]() {
						goto l432
					}
				}
			l434:
				if !_rules[rulesp]() {
					goto l432
				}
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l437
					}
					position++
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if buffer[position] != rune('A') {
						goto l432
					}
					position++
				}
			l436:
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if buffer[position] != rune('S') {
						goto l432
					}
					position++
				}
			l438:
				if !_rules[rulesp]() {
					goto l432
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l432
				}
				if !_rules[ruleAction19]() {
					goto l432
				}
				depth--
				add(ruleAliasExpression, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 24 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action20)> */
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
							if buffer[position] != rune('f') {
								goto l446
							}
							position++
							goto l445
						l446:
							position, tokenIndex, depth = position445, tokenIndex445, depth445
							if buffer[position] != rune('F') {
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
							if buffer[position] != rune('m') {
								goto l452
							}
							position++
							goto l451
						l452:
							position, tokenIndex, depth = position451, tokenIndex451, depth451
							if buffer[position] != rune('M') {
								goto l443
							}
							position++
						}
					l451:
						if !_rules[rulesp]() {
							goto l443
						}
						if !_rules[ruleRelations]() {
							goto l443
						}
						if !_rules[rulesp]() {
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
				if !_rules[ruleAction20]() {
					goto l440
				}
				depth--
				add(ruleWindowedFrom, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 25 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position453, tokenIndex453, depth453 := position, tokenIndex, depth
			{
				position454 := position
				depth++
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l456
					}
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if !_rules[ruleTuplesInterval]() {
						goto l453
					}
				}
			l455:
				depth--
				add(ruleInterval, position454)
			}
			return true
		l453:
			position, tokenIndex, depth = position453, tokenIndex453, depth453
			return false
		},
		/* 26 TimeInterval <- <(NumericLiteral sp SECONDS Action21)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l457
				}
				if !_rules[rulesp]() {
					goto l457
				}
				if !_rules[ruleSECONDS]() {
					goto l457
				}
				if !_rules[ruleAction21]() {
					goto l457
				}
				depth--
				add(ruleTimeInterval, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 27 TuplesInterval <- <(NumericLiteral sp TUPLES Action22)> */
		func() bool {
			position459, tokenIndex459, depth459 := position, tokenIndex, depth
			{
				position460 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l459
				}
				if !_rules[rulesp]() {
					goto l459
				}
				if !_rules[ruleTUPLES]() {
					goto l459
				}
				if !_rules[ruleAction22]() {
					goto l459
				}
				depth--
				add(ruleTuplesInterval, position460)
			}
			return true
		l459:
			position, tokenIndex, depth = position459, tokenIndex459, depth459
			return false
		},
		/* 28 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l461
				}
				if !_rules[rulesp]() {
					goto l461
				}
			l463:
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l464
					}
					position++
					if !_rules[rulesp]() {
						goto l464
					}
					if !_rules[ruleRelationLike]() {
						goto l464
					}
					goto l463
				l464:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
				}
				depth--
				add(ruleRelations, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 29 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action23)> */
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
							if buffer[position] != rune('w') {
								goto l471
							}
							position++
							goto l470
						l471:
							position, tokenIndex, depth = position470, tokenIndex470, depth470
							if buffer[position] != rune('W') {
								goto l468
							}
							position++
						}
					l470:
						{
							position472, tokenIndex472, depth472 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l473
							}
							position++
							goto l472
						l473:
							position, tokenIndex, depth = position472, tokenIndex472, depth472
							if buffer[position] != rune('H') {
								goto l468
							}
							position++
						}
					l472:
						{
							position474, tokenIndex474, depth474 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l475
							}
							position++
							goto l474
						l475:
							position, tokenIndex, depth = position474, tokenIndex474, depth474
							if buffer[position] != rune('E') {
								goto l468
							}
							position++
						}
					l474:
						{
							position476, tokenIndex476, depth476 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l477
							}
							position++
							goto l476
						l477:
							position, tokenIndex, depth = position476, tokenIndex476, depth476
							if buffer[position] != rune('R') {
								goto l468
							}
							position++
						}
					l476:
						{
							position478, tokenIndex478, depth478 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l479
							}
							position++
							goto l478
						l479:
							position, tokenIndex, depth = position478, tokenIndex478, depth478
							if buffer[position] != rune('E') {
								goto l468
							}
							position++
						}
					l478:
						if !_rules[rulesp]() {
							goto l468
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction23]() {
					goto l465
				}
				depth--
				add(ruleFilter, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 30 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action24)> */
		func() bool {
			position480, tokenIndex480, depth480 := position, tokenIndex, depth
			{
				position481 := position
				depth++
				{
					position482 := position
					depth++
					{
						position483, tokenIndex483, depth483 := position, tokenIndex, depth
						{
							position485, tokenIndex485, depth485 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l486
							}
							position++
							goto l485
						l486:
							position, tokenIndex, depth = position485, tokenIndex485, depth485
							if buffer[position] != rune('G') {
								goto l483
							}
							position++
						}
					l485:
						{
							position487, tokenIndex487, depth487 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l488
							}
							position++
							goto l487
						l488:
							position, tokenIndex, depth = position487, tokenIndex487, depth487
							if buffer[position] != rune('R') {
								goto l483
							}
							position++
						}
					l487:
						{
							position489, tokenIndex489, depth489 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l490
							}
							position++
							goto l489
						l490:
							position, tokenIndex, depth = position489, tokenIndex489, depth489
							if buffer[position] != rune('O') {
								goto l483
							}
							position++
						}
					l489:
						{
							position491, tokenIndex491, depth491 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l492
							}
							position++
							goto l491
						l492:
							position, tokenIndex, depth = position491, tokenIndex491, depth491
							if buffer[position] != rune('U') {
								goto l483
							}
							position++
						}
					l491:
						{
							position493, tokenIndex493, depth493 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l494
							}
							position++
							goto l493
						l494:
							position, tokenIndex, depth = position493, tokenIndex493, depth493
							if buffer[position] != rune('P') {
								goto l483
							}
							position++
						}
					l493:
						if !_rules[rulesp]() {
							goto l483
						}
						{
							position495, tokenIndex495, depth495 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l496
							}
							position++
							goto l495
						l496:
							position, tokenIndex, depth = position495, tokenIndex495, depth495
							if buffer[position] != rune('B') {
								goto l483
							}
							position++
						}
					l495:
						{
							position497, tokenIndex497, depth497 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l498
							}
							position++
							goto l497
						l498:
							position, tokenIndex, depth = position497, tokenIndex497, depth497
							if buffer[position] != rune('Y') {
								goto l483
							}
							position++
						}
					l497:
						if !_rules[rulesp]() {
							goto l483
						}
						if !_rules[ruleGroupList]() {
							goto l483
						}
						goto l484
					l483:
						position, tokenIndex, depth = position483, tokenIndex483, depth483
					}
				l484:
					depth--
					add(rulePegText, position482)
				}
				if !_rules[ruleAction24]() {
					goto l480
				}
				depth--
				add(ruleGrouping, position481)
			}
			return true
		l480:
			position, tokenIndex, depth = position480, tokenIndex480, depth480
			return false
		},
		/* 31 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l499
				}
				if !_rules[rulesp]() {
					goto l499
				}
			l501:
				{
					position502, tokenIndex502, depth502 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l502
					}
					position++
					if !_rules[rulesp]() {
						goto l502
					}
					if !_rules[ruleExpression]() {
						goto l502
					}
					goto l501
				l502:
					position, tokenIndex, depth = position502, tokenIndex502, depth502
				}
				depth--
				add(ruleGroupList, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 32 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action25)> */
		func() bool {
			position503, tokenIndex503, depth503 := position, tokenIndex, depth
			{
				position504 := position
				depth++
				{
					position505 := position
					depth++
					{
						position506, tokenIndex506, depth506 := position, tokenIndex, depth
						{
							position508, tokenIndex508, depth508 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l509
							}
							position++
							goto l508
						l509:
							position, tokenIndex, depth = position508, tokenIndex508, depth508
							if buffer[position] != rune('H') {
								goto l506
							}
							position++
						}
					l508:
						{
							position510, tokenIndex510, depth510 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l511
							}
							position++
							goto l510
						l511:
							position, tokenIndex, depth = position510, tokenIndex510, depth510
							if buffer[position] != rune('A') {
								goto l506
							}
							position++
						}
					l510:
						{
							position512, tokenIndex512, depth512 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l513
							}
							position++
							goto l512
						l513:
							position, tokenIndex, depth = position512, tokenIndex512, depth512
							if buffer[position] != rune('V') {
								goto l506
							}
							position++
						}
					l512:
						{
							position514, tokenIndex514, depth514 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l515
							}
							position++
							goto l514
						l515:
							position, tokenIndex, depth = position514, tokenIndex514, depth514
							if buffer[position] != rune('I') {
								goto l506
							}
							position++
						}
					l514:
						{
							position516, tokenIndex516, depth516 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l517
							}
							position++
							goto l516
						l517:
							position, tokenIndex, depth = position516, tokenIndex516, depth516
							if buffer[position] != rune('N') {
								goto l506
							}
							position++
						}
					l516:
						{
							position518, tokenIndex518, depth518 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l519
							}
							position++
							goto l518
						l519:
							position, tokenIndex, depth = position518, tokenIndex518, depth518
							if buffer[position] != rune('G') {
								goto l506
							}
							position++
						}
					l518:
						if !_rules[rulesp]() {
							goto l506
						}
						if !_rules[ruleExpression]() {
							goto l506
						}
						goto l507
					l506:
						position, tokenIndex, depth = position506, tokenIndex506, depth506
					}
				l507:
					depth--
					add(rulePegText, position505)
				}
				if !_rules[ruleAction25]() {
					goto l503
				}
				depth--
				add(ruleHaving, position504)
			}
			return true
		l503:
			position, tokenIndex, depth = position503, tokenIndex503, depth503
			return false
		},
		/* 33 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action26))> */
		func() bool {
			position520, tokenIndex520, depth520 := position, tokenIndex, depth
			{
				position521 := position
				depth++
				{
					position522, tokenIndex522, depth522 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l523
					}
					goto l522
				l523:
					position, tokenIndex, depth = position522, tokenIndex522, depth522
					if !_rules[ruleStreamWindow]() {
						goto l520
					}
					if !_rules[ruleAction26]() {
						goto l520
					}
				}
			l522:
				depth--
				add(ruleRelationLike, position521)
			}
			return true
		l520:
			position, tokenIndex, depth = position520, tokenIndex520, depth520
			return false
		},
		/* 34 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action27)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l524
				}
				if !_rules[rulesp]() {
					goto l524
				}
				{
					position526, tokenIndex526, depth526 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l527
					}
					position++
					goto l526
				l527:
					position, tokenIndex, depth = position526, tokenIndex526, depth526
					if buffer[position] != rune('A') {
						goto l524
					}
					position++
				}
			l526:
				{
					position528, tokenIndex528, depth528 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l529
					}
					position++
					goto l528
				l529:
					position, tokenIndex, depth = position528, tokenIndex528, depth528
					if buffer[position] != rune('S') {
						goto l524
					}
					position++
				}
			l528:
				if !_rules[rulesp]() {
					goto l524
				}
				if !_rules[ruleIdentifier]() {
					goto l524
				}
				if !_rules[ruleAction27]() {
					goto l524
				}
				depth--
				add(ruleAliasedStreamWindow, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 35 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action28)> */
		func() bool {
			position530, tokenIndex530, depth530 := position, tokenIndex, depth
			{
				position531 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l530
				}
				if !_rules[rulesp]() {
					goto l530
				}
				if buffer[position] != rune('[') {
					goto l530
				}
				position++
				if !_rules[rulesp]() {
					goto l530
				}
				{
					position532, tokenIndex532, depth532 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l533
					}
					position++
					goto l532
				l533:
					position, tokenIndex, depth = position532, tokenIndex532, depth532
					if buffer[position] != rune('R') {
						goto l530
					}
					position++
				}
			l532:
				{
					position534, tokenIndex534, depth534 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l535
					}
					position++
					goto l534
				l535:
					position, tokenIndex, depth = position534, tokenIndex534, depth534
					if buffer[position] != rune('A') {
						goto l530
					}
					position++
				}
			l534:
				{
					position536, tokenIndex536, depth536 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l537
					}
					position++
					goto l536
				l537:
					position, tokenIndex, depth = position536, tokenIndex536, depth536
					if buffer[position] != rune('N') {
						goto l530
					}
					position++
				}
			l536:
				{
					position538, tokenIndex538, depth538 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l539
					}
					position++
					goto l538
				l539:
					position, tokenIndex, depth = position538, tokenIndex538, depth538
					if buffer[position] != rune('G') {
						goto l530
					}
					position++
				}
			l538:
				{
					position540, tokenIndex540, depth540 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l541
					}
					position++
					goto l540
				l541:
					position, tokenIndex, depth = position540, tokenIndex540, depth540
					if buffer[position] != rune('E') {
						goto l530
					}
					position++
				}
			l540:
				if !_rules[rulesp]() {
					goto l530
				}
				if !_rules[ruleInterval]() {
					goto l530
				}
				if !_rules[rulesp]() {
					goto l530
				}
				if buffer[position] != rune(']') {
					goto l530
				}
				position++
				if !_rules[ruleAction28]() {
					goto l530
				}
				depth--
				add(ruleStreamWindow, position531)
			}
			return true
		l530:
			position, tokenIndex, depth = position530, tokenIndex530, depth530
			return false
		},
		/* 36 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position542, tokenIndex542, depth542 := position, tokenIndex, depth
			{
				position543 := position
				depth++
				{
					position544, tokenIndex544, depth544 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l545
					}
					goto l544
				l545:
					position, tokenIndex, depth = position544, tokenIndex544, depth544
					if !_rules[ruleStream]() {
						goto l542
					}
				}
			l544:
				depth--
				add(ruleStreamLike, position543)
			}
			return true
		l542:
			position, tokenIndex, depth = position542, tokenIndex542, depth542
			return false
		},
		/* 37 UDSFFuncApp <- <(FuncApp Action29)> */
		func() bool {
			position546, tokenIndex546, depth546 := position, tokenIndex, depth
			{
				position547 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l546
				}
				if !_rules[ruleAction29]() {
					goto l546
				}
				depth--
				add(ruleUDSFFuncApp, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 38 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action30)> */
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
							if buffer[position] != rune('w') {
								goto l554
							}
							position++
							goto l553
						l554:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if buffer[position] != rune('W') {
								goto l551
							}
							position++
						}
					l553:
						{
							position555, tokenIndex555, depth555 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l556
							}
							position++
							goto l555
						l556:
							position, tokenIndex, depth = position555, tokenIndex555, depth555
							if buffer[position] != rune('I') {
								goto l551
							}
							position++
						}
					l555:
						{
							position557, tokenIndex557, depth557 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l558
							}
							position++
							goto l557
						l558:
							position, tokenIndex, depth = position557, tokenIndex557, depth557
							if buffer[position] != rune('T') {
								goto l551
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
								goto l551
							}
							position++
						}
					l559:
						if !_rules[rulesp]() {
							goto l551
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l551
						}
						if !_rules[rulesp]() {
							goto l551
						}
					l561:
						{
							position562, tokenIndex562, depth562 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l562
							}
							position++
							if !_rules[rulesp]() {
								goto l562
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l562
							}
							goto l561
						l562:
							position, tokenIndex, depth = position562, tokenIndex562, depth562
						}
						goto l552
					l551:
						position, tokenIndex, depth = position551, tokenIndex551, depth551
					}
				l552:
					depth--
					add(rulePegText, position550)
				}
				if !_rules[ruleAction30]() {
					goto l548
				}
				depth--
				add(ruleSourceSinkSpecs, position549)
			}
			return true
		l548:
			position, tokenIndex, depth = position548, tokenIndex548, depth548
			return false
		},
		/* 39 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action31)> */
		func() bool {
			position563, tokenIndex563, depth563 := position, tokenIndex, depth
			{
				position564 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l563
				}
				if buffer[position] != rune('=') {
					goto l563
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l563
				}
				if !_rules[ruleAction31]() {
					goto l563
				}
				depth--
				add(ruleSourceSinkParam, position564)
			}
			return true
		l563:
			position, tokenIndex, depth = position563, tokenIndex563, depth563
			return false
		},
		/* 40 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567, tokenIndex567, depth567 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l568
					}
					goto l567
				l568:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleLiteral]() {
						goto l565
					}
				}
			l567:
				depth--
				add(ruleSourceSinkParamVal, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 41 PausedOpt <- <(<(Paused / Unpaused)?> Action32)> */
		func() bool {
			position569, tokenIndex569, depth569 := position, tokenIndex, depth
			{
				position570 := position
				depth++
				{
					position571 := position
					depth++
					{
						position572, tokenIndex572, depth572 := position, tokenIndex, depth
						{
							position574, tokenIndex574, depth574 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l575
							}
							goto l574
						l575:
							position, tokenIndex, depth = position574, tokenIndex574, depth574
							if !_rules[ruleUnpaused]() {
								goto l572
							}
						}
					l574:
						goto l573
					l572:
						position, tokenIndex, depth = position572, tokenIndex572, depth572
					}
				l573:
					depth--
					add(rulePegText, position571)
				}
				if !_rules[ruleAction32]() {
					goto l569
				}
				depth--
				add(rulePausedOpt, position570)
			}
			return true
		l569:
			position, tokenIndex, depth = position569, tokenIndex569, depth569
			return false
		},
		/* 42 Expression <- <orExpr> */
		func() bool {
			position576, tokenIndex576, depth576 := position, tokenIndex, depth
			{
				position577 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l576
				}
				depth--
				add(ruleExpression, position577)
			}
			return true
		l576:
			position, tokenIndex, depth = position576, tokenIndex576, depth576
			return false
		},
		/* 43 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action33)> */
		func() bool {
			position578, tokenIndex578, depth578 := position, tokenIndex, depth
			{
				position579 := position
				depth++
				{
					position580 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l578
					}
					if !_rules[rulesp]() {
						goto l578
					}
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l581
						}
						if !_rules[rulesp]() {
							goto l581
						}
						if !_rules[ruleandExpr]() {
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
				if !_rules[ruleAction33]() {
					goto l578
				}
				depth--
				add(ruleorExpr, position579)
			}
			return true
		l578:
			position, tokenIndex, depth = position578, tokenIndex578, depth578
			return false
		},
		/* 44 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action34)> */
		func() bool {
			position583, tokenIndex583, depth583 := position, tokenIndex, depth
			{
				position584 := position
				depth++
				{
					position585 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l583
					}
					if !_rules[rulesp]() {
						goto l583
					}
					{
						position586, tokenIndex586, depth586 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l586
						}
						if !_rules[rulesp]() {
							goto l586
						}
						if !_rules[rulenotExpr]() {
							goto l586
						}
						goto l587
					l586:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
					}
				l587:
					depth--
					add(rulePegText, position585)
				}
				if !_rules[ruleAction34]() {
					goto l583
				}
				depth--
				add(ruleandExpr, position584)
			}
			return true
		l583:
			position, tokenIndex, depth = position583, tokenIndex583, depth583
			return false
		},
		/* 45 notExpr <- <(<((Not sp)? comparisonExpr)> Action35)> */
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
						if !_rules[ruleNot]() {
							goto l591
						}
						if !_rules[rulesp]() {
							goto l591
						}
						goto l592
					l591:
						position, tokenIndex, depth = position591, tokenIndex591, depth591
					}
				l592:
					if !_rules[rulecomparisonExpr]() {
						goto l588
					}
					depth--
					add(rulePegText, position590)
				}
				if !_rules[ruleAction35]() {
					goto l588
				}
				depth--
				add(rulenotExpr, position589)
			}
			return true
		l588:
			position, tokenIndex, depth = position588, tokenIndex588, depth588
			return false
		},
		/* 46 comparisonExpr <- <(<(isExpr sp (ComparisonOp sp isExpr)?)> Action36)> */
		func() bool {
			position593, tokenIndex593, depth593 := position, tokenIndex, depth
			{
				position594 := position
				depth++
				{
					position595 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l593
					}
					if !_rules[rulesp]() {
						goto l593
					}
					{
						position596, tokenIndex596, depth596 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l596
						}
						if !_rules[rulesp]() {
							goto l596
						}
						if !_rules[ruleisExpr]() {
							goto l596
						}
						goto l597
					l596:
						position, tokenIndex, depth = position596, tokenIndex596, depth596
					}
				l597:
					depth--
					add(rulePegText, position595)
				}
				if !_rules[ruleAction36]() {
					goto l593
				}
				depth--
				add(rulecomparisonExpr, position594)
			}
			return true
		l593:
			position, tokenIndex, depth = position593, tokenIndex593, depth593
			return false
		},
		/* 47 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action37)> */
		func() bool {
			position598, tokenIndex598, depth598 := position, tokenIndex, depth
			{
				position599 := position
				depth++
				{
					position600 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l598
					}
					if !_rules[rulesp]() {
						goto l598
					}
					{
						position601, tokenIndex601, depth601 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l601
						}
						if !_rules[rulesp]() {
							goto l601
						}
						if !_rules[ruleNullLiteral]() {
							goto l601
						}
						goto l602
					l601:
						position, tokenIndex, depth = position601, tokenIndex601, depth601
					}
				l602:
					depth--
					add(rulePegText, position600)
				}
				if !_rules[ruleAction37]() {
					goto l598
				}
				depth--
				add(ruleisExpr, position599)
			}
			return true
		l598:
			position, tokenIndex, depth = position598, tokenIndex598, depth598
			return false
		},
		/* 48 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action38)> */
		func() bool {
			position603, tokenIndex603, depth603 := position, tokenIndex, depth
			{
				position604 := position
				depth++
				{
					position605 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l603
					}
					if !_rules[rulesp]() {
						goto l603
					}
					{
						position606, tokenIndex606, depth606 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l606
						}
						if !_rules[rulesp]() {
							goto l606
						}
						if !_rules[ruleproductExpr]() {
							goto l606
						}
						goto l607
					l606:
						position, tokenIndex, depth = position606, tokenIndex606, depth606
					}
				l607:
					depth--
					add(rulePegText, position605)
				}
				if !_rules[ruleAction38]() {
					goto l603
				}
				depth--
				add(ruletermExpr, position604)
			}
			return true
		l603:
			position, tokenIndex, depth = position603, tokenIndex603, depth603
			return false
		},
		/* 49 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr)?)> Action39)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l608
					}
					if !_rules[rulesp]() {
						goto l608
					}
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l611
						}
						if !_rules[rulesp]() {
							goto l611
						}
						if !_rules[ruleminusExpr]() {
							goto l611
						}
						goto l612
					l611:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
					}
				l612:
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction39]() {
					goto l608
				}
				depth--
				add(ruleproductExpr, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 50 minusExpr <- <(<((UnaryMinus sp)? baseExpr)> Action40)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l616
						}
						if !_rules[rulesp]() {
							goto l616
						}
						goto l617
					l616:
						position, tokenIndex, depth = position616, tokenIndex616, depth616
					}
				l617:
					if !_rules[rulebaseExpr]() {
						goto l613
					}
					depth--
					add(rulePegText, position615)
				}
				if !_rules[ruleAction40]() {
					goto l613
				}
				depth--
				add(ruleminusExpr, position614)
			}
			return true
		l613:
			position, tokenIndex, depth = position613, tokenIndex613, depth613
			return false
		},
		/* 51 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / FuncApp / RowMeta / RowValue / Literal)> */
		func() bool {
			position618, tokenIndex618, depth618 := position, tokenIndex, depth
			{
				position619 := position
				depth++
				{
					position620, tokenIndex620, depth620 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l621
					}
					position++
					if !_rules[rulesp]() {
						goto l621
					}
					if !_rules[ruleExpression]() {
						goto l621
					}
					if !_rules[rulesp]() {
						goto l621
					}
					if buffer[position] != rune(')') {
						goto l621
					}
					position++
					goto l620
				l621:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleBooleanLiteral]() {
						goto l622
					}
					goto l620
				l622:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleNullLiteral]() {
						goto l623
					}
					goto l620
				l623:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleFuncApp]() {
						goto l624
					}
					goto l620
				l624:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleRowMeta]() {
						goto l625
					}
					goto l620
				l625:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleRowValue]() {
						goto l626
					}
					goto l620
				l626:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if !_rules[ruleLiteral]() {
						goto l618
					}
				}
			l620:
				depth--
				add(rulebaseExpr, position619)
			}
			return true
		l618:
			position, tokenIndex, depth = position618, tokenIndex618, depth618
			return false
		},
		/* 52 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action41)> */
		func() bool {
			position627, tokenIndex627, depth627 := position, tokenIndex, depth
			{
				position628 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l627
				}
				if !_rules[rulesp]() {
					goto l627
				}
				if buffer[position] != rune('(') {
					goto l627
				}
				position++
				if !_rules[rulesp]() {
					goto l627
				}
				if !_rules[ruleFuncParams]() {
					goto l627
				}
				if !_rules[rulesp]() {
					goto l627
				}
				if buffer[position] != rune(')') {
					goto l627
				}
				position++
				if !_rules[ruleAction41]() {
					goto l627
				}
				depth--
				add(ruleFuncApp, position628)
			}
			return true
		l627:
			position, tokenIndex, depth = position627, tokenIndex627, depth627
			return false
		},
		/* 53 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action42)> */
		func() bool {
			position629, tokenIndex629, depth629 := position, tokenIndex, depth
			{
				position630 := position
				depth++
				{
					position631 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l629
					}
					if !_rules[rulesp]() {
						goto l629
					}
				l632:
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l633
						}
						position++
						if !_rules[rulesp]() {
							goto l633
						}
						if !_rules[ruleExpression]() {
							goto l633
						}
						goto l632
					l633:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
					}
					depth--
					add(rulePegText, position631)
				}
				if !_rules[ruleAction42]() {
					goto l629
				}
				depth--
				add(ruleFuncParams, position630)
			}
			return true
		l629:
			position, tokenIndex, depth = position629, tokenIndex629, depth629
			return false
		},
		/* 54 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position634, tokenIndex634, depth634 := position, tokenIndex, depth
			{
				position635 := position
				depth++
				{
					position636, tokenIndex636, depth636 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l637
					}
					goto l636
				l637:
					position, tokenIndex, depth = position636, tokenIndex636, depth636
					if !_rules[ruleNumericLiteral]() {
						goto l638
					}
					goto l636
				l638:
					position, tokenIndex, depth = position636, tokenIndex636, depth636
					if !_rules[ruleStringLiteral]() {
						goto l634
					}
				}
			l636:
				depth--
				add(ruleLiteral, position635)
			}
			return true
		l634:
			position, tokenIndex, depth = position634, tokenIndex634, depth634
			return false
		},
		/* 55 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position639, tokenIndex639, depth639 := position, tokenIndex, depth
			{
				position640 := position
				depth++
				{
					position641, tokenIndex641, depth641 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l642
					}
					goto l641
				l642:
					position, tokenIndex, depth = position641, tokenIndex641, depth641
					if !_rules[ruleNotEqual]() {
						goto l643
					}
					goto l641
				l643:
					position, tokenIndex, depth = position641, tokenIndex641, depth641
					if !_rules[ruleLessOrEqual]() {
						goto l644
					}
					goto l641
				l644:
					position, tokenIndex, depth = position641, tokenIndex641, depth641
					if !_rules[ruleLess]() {
						goto l645
					}
					goto l641
				l645:
					position, tokenIndex, depth = position641, tokenIndex641, depth641
					if !_rules[ruleGreaterOrEqual]() {
						goto l646
					}
					goto l641
				l646:
					position, tokenIndex, depth = position641, tokenIndex641, depth641
					if !_rules[ruleGreater]() {
						goto l647
					}
					goto l641
				l647:
					position, tokenIndex, depth = position641, tokenIndex641, depth641
					if !_rules[ruleNotEqual]() {
						goto l639
					}
				}
			l641:
				depth--
				add(ruleComparisonOp, position640)
			}
			return true
		l639:
			position, tokenIndex, depth = position639, tokenIndex639, depth639
			return false
		},
		/* 56 IsOp <- <(IsNot / Is)> */
		func() bool {
			position648, tokenIndex648, depth648 := position, tokenIndex, depth
			{
				position649 := position
				depth++
				{
					position650, tokenIndex650, depth650 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l651
					}
					goto l650
				l651:
					position, tokenIndex, depth = position650, tokenIndex650, depth650
					if !_rules[ruleIs]() {
						goto l648
					}
				}
			l650:
				depth--
				add(ruleIsOp, position649)
			}
			return true
		l648:
			position, tokenIndex, depth = position648, tokenIndex648, depth648
			return false
		},
		/* 57 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				{
					position654, tokenIndex654, depth654 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l655
					}
					goto l654
				l655:
					position, tokenIndex, depth = position654, tokenIndex654, depth654
					if !_rules[ruleMinus]() {
						goto l652
					}
				}
			l654:
				depth--
				add(rulePlusMinusOp, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 58 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				{
					position658, tokenIndex658, depth658 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l659
					}
					goto l658
				l659:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
					if !_rules[ruleDivide]() {
						goto l660
					}
					goto l658
				l660:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
					if !_rules[ruleModulo]() {
						goto l656
					}
				}
			l658:
				depth--
				add(ruleMultDivOp, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 59 Stream <- <(<ident> Action43)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				{
					position663 := position
					depth++
					if !_rules[ruleident]() {
						goto l661
					}
					depth--
					add(rulePegText, position663)
				}
				if !_rules[ruleAction43]() {
					goto l661
				}
				depth--
				add(ruleStream, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 60 RowMeta <- <RowTimestamp> */
		func() bool {
			position664, tokenIndex664, depth664 := position, tokenIndex, depth
			{
				position665 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l664
				}
				depth--
				add(ruleRowMeta, position665)
			}
			return true
		l664:
			position, tokenIndex, depth = position664, tokenIndex664, depth664
			return false
		},
		/* 61 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action44)> */
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
						if !_rules[ruleident]() {
							goto l669
						}
						if buffer[position] != rune(':') {
							goto l669
						}
						position++
						goto l670
					l669:
						position, tokenIndex, depth = position669, tokenIndex669, depth669
					}
				l670:
					if buffer[position] != rune('t') {
						goto l666
					}
					position++
					if buffer[position] != rune('s') {
						goto l666
					}
					position++
					if buffer[position] != rune('(') {
						goto l666
					}
					position++
					if buffer[position] != rune(')') {
						goto l666
					}
					position++
					depth--
					add(rulePegText, position668)
				}
				if !_rules[ruleAction44]() {
					goto l666
				}
				depth--
				add(ruleRowTimestamp, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 62 RowValue <- <(<((ident ':')? jsonPath)> Action45)> */
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
						if !_rules[ruleident]() {
							goto l674
						}
						if buffer[position] != rune(':') {
							goto l674
						}
						position++
						goto l675
					l674:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
					}
				l675:
					if !_rules[rulejsonPath]() {
						goto l671
					}
					depth--
					add(rulePegText, position673)
				}
				if !_rules[ruleAction45]() {
					goto l671
				}
				depth--
				add(ruleRowValue, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 63 NumericLiteral <- <(<('-'? [0-9]+)> Action46)> */
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
						if buffer[position] != rune('-') {
							goto l679
						}
						position++
						goto l680
					l679:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
					}
				l680:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l676
					}
					position++
				l681:
					{
						position682, tokenIndex682, depth682 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position682, tokenIndex682, depth682
					}
					depth--
					add(rulePegText, position678)
				}
				if !_rules[ruleAction46]() {
					goto l676
				}
				depth--
				add(ruleNumericLiteral, position677)
			}
			return true
		l676:
			position, tokenIndex, depth = position676, tokenIndex676, depth676
			return false
		},
		/* 64 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action47)> */
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
						if buffer[position] != rune('-') {
							goto l686
						}
						position++
						goto l687
					l686:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
					}
				l687:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l683
					}
					position++
				l688:
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
					}
					if buffer[position] != rune('.') {
						goto l683
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l683
					}
					position++
				l690:
					{
						position691, tokenIndex691, depth691 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l691
						}
						position++
						goto l690
					l691:
						position, tokenIndex, depth = position691, tokenIndex691, depth691
					}
					depth--
					add(rulePegText, position685)
				}
				if !_rules[ruleAction47]() {
					goto l683
				}
				depth--
				add(ruleFloatLiteral, position684)
			}
			return true
		l683:
			position, tokenIndex, depth = position683, tokenIndex683, depth683
			return false
		},
		/* 65 Function <- <(<ident> Action48)> */
		func() bool {
			position692, tokenIndex692, depth692 := position, tokenIndex, depth
			{
				position693 := position
				depth++
				{
					position694 := position
					depth++
					if !_rules[ruleident]() {
						goto l692
					}
					depth--
					add(rulePegText, position694)
				}
				if !_rules[ruleAction48]() {
					goto l692
				}
				depth--
				add(ruleFunction, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 66 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action49)> */
		func() bool {
			position695, tokenIndex695, depth695 := position, tokenIndex, depth
			{
				position696 := position
				depth++
				{
					position697 := position
					depth++
					{
						position698, tokenIndex698, depth698 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l699
						}
						position++
						goto l698
					l699:
						position, tokenIndex, depth = position698, tokenIndex698, depth698
						if buffer[position] != rune('N') {
							goto l695
						}
						position++
					}
				l698:
					{
						position700, tokenIndex700, depth700 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l701
						}
						position++
						goto l700
					l701:
						position, tokenIndex, depth = position700, tokenIndex700, depth700
						if buffer[position] != rune('U') {
							goto l695
						}
						position++
					}
				l700:
					{
						position702, tokenIndex702, depth702 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l703
						}
						position++
						goto l702
					l703:
						position, tokenIndex, depth = position702, tokenIndex702, depth702
						if buffer[position] != rune('L') {
							goto l695
						}
						position++
					}
				l702:
					{
						position704, tokenIndex704, depth704 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l705
						}
						position++
						goto l704
					l705:
						position, tokenIndex, depth = position704, tokenIndex704, depth704
						if buffer[position] != rune('L') {
							goto l695
						}
						position++
					}
				l704:
					depth--
					add(rulePegText, position697)
				}
				if !_rules[ruleAction49]() {
					goto l695
				}
				depth--
				add(ruleNullLiteral, position696)
			}
			return true
		l695:
			position, tokenIndex, depth = position695, tokenIndex695, depth695
			return false
		},
		/* 67 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position706, tokenIndex706, depth706 := position, tokenIndex, depth
			{
				position707 := position
				depth++
				{
					position708, tokenIndex708, depth708 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l709
					}
					goto l708
				l709:
					position, tokenIndex, depth = position708, tokenIndex708, depth708
					if !_rules[ruleFALSE]() {
						goto l706
					}
				}
			l708:
				depth--
				add(ruleBooleanLiteral, position707)
			}
			return true
		l706:
			position, tokenIndex, depth = position706, tokenIndex706, depth706
			return false
		},
		/* 68 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action50)> */
		func() bool {
			position710, tokenIndex710, depth710 := position, tokenIndex, depth
			{
				position711 := position
				depth++
				{
					position712 := position
					depth++
					{
						position713, tokenIndex713, depth713 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l714
						}
						position++
						goto l713
					l714:
						position, tokenIndex, depth = position713, tokenIndex713, depth713
						if buffer[position] != rune('T') {
							goto l710
						}
						position++
					}
				l713:
					{
						position715, tokenIndex715, depth715 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l716
						}
						position++
						goto l715
					l716:
						position, tokenIndex, depth = position715, tokenIndex715, depth715
						if buffer[position] != rune('R') {
							goto l710
						}
						position++
					}
				l715:
					{
						position717, tokenIndex717, depth717 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l718
						}
						position++
						goto l717
					l718:
						position, tokenIndex, depth = position717, tokenIndex717, depth717
						if buffer[position] != rune('U') {
							goto l710
						}
						position++
					}
				l717:
					{
						position719, tokenIndex719, depth719 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l720
						}
						position++
						goto l719
					l720:
						position, tokenIndex, depth = position719, tokenIndex719, depth719
						if buffer[position] != rune('E') {
							goto l710
						}
						position++
					}
				l719:
					depth--
					add(rulePegText, position712)
				}
				if !_rules[ruleAction50]() {
					goto l710
				}
				depth--
				add(ruleTRUE, position711)
			}
			return true
		l710:
			position, tokenIndex, depth = position710, tokenIndex710, depth710
			return false
		},
		/* 69 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action51)> */
		func() bool {
			position721, tokenIndex721, depth721 := position, tokenIndex, depth
			{
				position722 := position
				depth++
				{
					position723 := position
					depth++
					{
						position724, tokenIndex724, depth724 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l725
						}
						position++
						goto l724
					l725:
						position, tokenIndex, depth = position724, tokenIndex724, depth724
						if buffer[position] != rune('F') {
							goto l721
						}
						position++
					}
				l724:
					{
						position726, tokenIndex726, depth726 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l727
						}
						position++
						goto l726
					l727:
						position, tokenIndex, depth = position726, tokenIndex726, depth726
						if buffer[position] != rune('A') {
							goto l721
						}
						position++
					}
				l726:
					{
						position728, tokenIndex728, depth728 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l729
						}
						position++
						goto l728
					l729:
						position, tokenIndex, depth = position728, tokenIndex728, depth728
						if buffer[position] != rune('L') {
							goto l721
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
							goto l721
						}
						position++
					}
				l730:
					{
						position732, tokenIndex732, depth732 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l733
						}
						position++
						goto l732
					l733:
						position, tokenIndex, depth = position732, tokenIndex732, depth732
						if buffer[position] != rune('E') {
							goto l721
						}
						position++
					}
				l732:
					depth--
					add(rulePegText, position723)
				}
				if !_rules[ruleAction51]() {
					goto l721
				}
				depth--
				add(ruleFALSE, position722)
			}
			return true
		l721:
			position, tokenIndex, depth = position721, tokenIndex721, depth721
			return false
		},
		/* 70 Wildcard <- <(<'*'> Action52)> */
		func() bool {
			position734, tokenIndex734, depth734 := position, tokenIndex, depth
			{
				position735 := position
				depth++
				{
					position736 := position
					depth++
					if buffer[position] != rune('*') {
						goto l734
					}
					position++
					depth--
					add(rulePegText, position736)
				}
				if !_rules[ruleAction52]() {
					goto l734
				}
				depth--
				add(ruleWildcard, position735)
			}
			return true
		l734:
			position, tokenIndex, depth = position734, tokenIndex734, depth734
			return false
		},
		/* 71 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action53)> */
		func() bool {
			position737, tokenIndex737, depth737 := position, tokenIndex, depth
			{
				position738 := position
				depth++
				{
					position739 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l737
					}
					position++
				l740:
					{
						position741, tokenIndex741, depth741 := position, tokenIndex, depth
						{
							position742, tokenIndex742, depth742 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l743
							}
							position++
							if buffer[position] != rune('\'') {
								goto l743
							}
							position++
							goto l742
						l743:
							position, tokenIndex, depth = position742, tokenIndex742, depth742
							{
								position744, tokenIndex744, depth744 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l744
								}
								position++
								goto l741
							l744:
								position, tokenIndex, depth = position744, tokenIndex744, depth744
							}
							if !matchDot() {
								goto l741
							}
						}
					l742:
						goto l740
					l741:
						position, tokenIndex, depth = position741, tokenIndex741, depth741
					}
					if buffer[position] != rune('\'') {
						goto l737
					}
					position++
					depth--
					add(rulePegText, position739)
				}
				if !_rules[ruleAction53]() {
					goto l737
				}
				depth--
				add(ruleStringLiteral, position738)
			}
			return true
		l737:
			position, tokenIndex, depth = position737, tokenIndex737, depth737
			return false
		},
		/* 72 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action54)> */
		func() bool {
			position745, tokenIndex745, depth745 := position, tokenIndex, depth
			{
				position746 := position
				depth++
				{
					position747 := position
					depth++
					{
						position748, tokenIndex748, depth748 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l749
						}
						position++
						goto l748
					l749:
						position, tokenIndex, depth = position748, tokenIndex748, depth748
						if buffer[position] != rune('I') {
							goto l745
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
							goto l745
						}
						position++
					}
				l750:
					{
						position752, tokenIndex752, depth752 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l753
						}
						position++
						goto l752
					l753:
						position, tokenIndex, depth = position752, tokenIndex752, depth752
						if buffer[position] != rune('T') {
							goto l745
						}
						position++
					}
				l752:
					{
						position754, tokenIndex754, depth754 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l755
						}
						position++
						goto l754
					l755:
						position, tokenIndex, depth = position754, tokenIndex754, depth754
						if buffer[position] != rune('R') {
							goto l745
						}
						position++
					}
				l754:
					{
						position756, tokenIndex756, depth756 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l757
						}
						position++
						goto l756
					l757:
						position, tokenIndex, depth = position756, tokenIndex756, depth756
						if buffer[position] != rune('E') {
							goto l745
						}
						position++
					}
				l756:
					{
						position758, tokenIndex758, depth758 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l759
						}
						position++
						goto l758
					l759:
						position, tokenIndex, depth = position758, tokenIndex758, depth758
						if buffer[position] != rune('A') {
							goto l745
						}
						position++
					}
				l758:
					{
						position760, tokenIndex760, depth760 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l761
						}
						position++
						goto l760
					l761:
						position, tokenIndex, depth = position760, tokenIndex760, depth760
						if buffer[position] != rune('M') {
							goto l745
						}
						position++
					}
				l760:
					depth--
					add(rulePegText, position747)
				}
				if !_rules[ruleAction54]() {
					goto l745
				}
				depth--
				add(ruleISTREAM, position746)
			}
			return true
		l745:
			position, tokenIndex, depth = position745, tokenIndex745, depth745
			return false
		},
		/* 73 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action55)> */
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
						if buffer[position] != rune('d') {
							goto l766
						}
						position++
						goto l765
					l766:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
						if buffer[position] != rune('D') {
							goto l762
						}
						position++
					}
				l765:
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('S') {
							goto l762
						}
						position++
					}
				l767:
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('T') {
							goto l762
						}
						position++
					}
				l769:
					{
						position771, tokenIndex771, depth771 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l772
						}
						position++
						goto l771
					l772:
						position, tokenIndex, depth = position771, tokenIndex771, depth771
						if buffer[position] != rune('R') {
							goto l762
						}
						position++
					}
				l771:
					{
						position773, tokenIndex773, depth773 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l774
						}
						position++
						goto l773
					l774:
						position, tokenIndex, depth = position773, tokenIndex773, depth773
						if buffer[position] != rune('E') {
							goto l762
						}
						position++
					}
				l773:
					{
						position775, tokenIndex775, depth775 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l776
						}
						position++
						goto l775
					l776:
						position, tokenIndex, depth = position775, tokenIndex775, depth775
						if buffer[position] != rune('A') {
							goto l762
						}
						position++
					}
				l775:
					{
						position777, tokenIndex777, depth777 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l778
						}
						position++
						goto l777
					l778:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
						if buffer[position] != rune('M') {
							goto l762
						}
						position++
					}
				l777:
					depth--
					add(rulePegText, position764)
				}
				if !_rules[ruleAction55]() {
					goto l762
				}
				depth--
				add(ruleDSTREAM, position763)
			}
			return true
		l762:
			position, tokenIndex, depth = position762, tokenIndex762, depth762
			return false
		},
		/* 74 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action56)> */
		func() bool {
			position779, tokenIndex779, depth779 := position, tokenIndex, depth
			{
				position780 := position
				depth++
				{
					position781 := position
					depth++
					{
						position782, tokenIndex782, depth782 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l783
						}
						position++
						goto l782
					l783:
						position, tokenIndex, depth = position782, tokenIndex782, depth782
						if buffer[position] != rune('R') {
							goto l779
						}
						position++
					}
				l782:
					{
						position784, tokenIndex784, depth784 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l785
						}
						position++
						goto l784
					l785:
						position, tokenIndex, depth = position784, tokenIndex784, depth784
						if buffer[position] != rune('S') {
							goto l779
						}
						position++
					}
				l784:
					{
						position786, tokenIndex786, depth786 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l787
						}
						position++
						goto l786
					l787:
						position, tokenIndex, depth = position786, tokenIndex786, depth786
						if buffer[position] != rune('T') {
							goto l779
						}
						position++
					}
				l786:
					{
						position788, tokenIndex788, depth788 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l789
						}
						position++
						goto l788
					l789:
						position, tokenIndex, depth = position788, tokenIndex788, depth788
						if buffer[position] != rune('R') {
							goto l779
						}
						position++
					}
				l788:
					{
						position790, tokenIndex790, depth790 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l791
						}
						position++
						goto l790
					l791:
						position, tokenIndex, depth = position790, tokenIndex790, depth790
						if buffer[position] != rune('E') {
							goto l779
						}
						position++
					}
				l790:
					{
						position792, tokenIndex792, depth792 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l793
						}
						position++
						goto l792
					l793:
						position, tokenIndex, depth = position792, tokenIndex792, depth792
						if buffer[position] != rune('A') {
							goto l779
						}
						position++
					}
				l792:
					{
						position794, tokenIndex794, depth794 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l795
						}
						position++
						goto l794
					l795:
						position, tokenIndex, depth = position794, tokenIndex794, depth794
						if buffer[position] != rune('M') {
							goto l779
						}
						position++
					}
				l794:
					depth--
					add(rulePegText, position781)
				}
				if !_rules[ruleAction56]() {
					goto l779
				}
				depth--
				add(ruleRSTREAM, position780)
			}
			return true
		l779:
			position, tokenIndex, depth = position779, tokenIndex779, depth779
			return false
		},
		/* 75 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action57)> */
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
						if buffer[position] != rune('u') {
							goto l802
						}
						position++
						goto l801
					l802:
						position, tokenIndex, depth = position801, tokenIndex801, depth801
						if buffer[position] != rune('U') {
							goto l796
						}
						position++
					}
				l801:
					{
						position803, tokenIndex803, depth803 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l804
						}
						position++
						goto l803
					l804:
						position, tokenIndex, depth = position803, tokenIndex803, depth803
						if buffer[position] != rune('P') {
							goto l796
						}
						position++
					}
				l803:
					{
						position805, tokenIndex805, depth805 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l806
						}
						position++
						goto l805
					l806:
						position, tokenIndex, depth = position805, tokenIndex805, depth805
						if buffer[position] != rune('L') {
							goto l796
						}
						position++
					}
				l805:
					{
						position807, tokenIndex807, depth807 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l808
						}
						position++
						goto l807
					l808:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						if buffer[position] != rune('E') {
							goto l796
						}
						position++
					}
				l807:
					{
						position809, tokenIndex809, depth809 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l810
						}
						position++
						goto l809
					l810:
						position, tokenIndex, depth = position809, tokenIndex809, depth809
						if buffer[position] != rune('S') {
							goto l796
						}
						position++
					}
				l809:
					depth--
					add(rulePegText, position798)
				}
				if !_rules[ruleAction57]() {
					goto l796
				}
				depth--
				add(ruleTUPLES, position797)
			}
			return true
		l796:
			position, tokenIndex, depth = position796, tokenIndex796, depth796
			return false
		},
		/* 76 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action58)> */
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
						if buffer[position] != rune('s') {
							goto l815
						}
						position++
						goto l814
					l815:
						position, tokenIndex, depth = position814, tokenIndex814, depth814
						if buffer[position] != rune('S') {
							goto l811
						}
						position++
					}
				l814:
					{
						position816, tokenIndex816, depth816 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l817
						}
						position++
						goto l816
					l817:
						position, tokenIndex, depth = position816, tokenIndex816, depth816
						if buffer[position] != rune('E') {
							goto l811
						}
						position++
					}
				l816:
					{
						position818, tokenIndex818, depth818 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l819
						}
						position++
						goto l818
					l819:
						position, tokenIndex, depth = position818, tokenIndex818, depth818
						if buffer[position] != rune('C') {
							goto l811
						}
						position++
					}
				l818:
					{
						position820, tokenIndex820, depth820 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l821
						}
						position++
						goto l820
					l821:
						position, tokenIndex, depth = position820, tokenIndex820, depth820
						if buffer[position] != rune('O') {
							goto l811
						}
						position++
					}
				l820:
					{
						position822, tokenIndex822, depth822 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l823
						}
						position++
						goto l822
					l823:
						position, tokenIndex, depth = position822, tokenIndex822, depth822
						if buffer[position] != rune('N') {
							goto l811
						}
						position++
					}
				l822:
					{
						position824, tokenIndex824, depth824 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l825
						}
						position++
						goto l824
					l825:
						position, tokenIndex, depth = position824, tokenIndex824, depth824
						if buffer[position] != rune('D') {
							goto l811
						}
						position++
					}
				l824:
					{
						position826, tokenIndex826, depth826 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l827
						}
						position++
						goto l826
					l827:
						position, tokenIndex, depth = position826, tokenIndex826, depth826
						if buffer[position] != rune('S') {
							goto l811
						}
						position++
					}
				l826:
					depth--
					add(rulePegText, position813)
				}
				if !_rules[ruleAction58]() {
					goto l811
				}
				depth--
				add(ruleSECONDS, position812)
			}
			return true
		l811:
			position, tokenIndex, depth = position811, tokenIndex811, depth811
			return false
		},
		/* 77 StreamIdentifier <- <(<ident> Action59)> */
		func() bool {
			position828, tokenIndex828, depth828 := position, tokenIndex, depth
			{
				position829 := position
				depth++
				{
					position830 := position
					depth++
					if !_rules[ruleident]() {
						goto l828
					}
					depth--
					add(rulePegText, position830)
				}
				if !_rules[ruleAction59]() {
					goto l828
				}
				depth--
				add(ruleStreamIdentifier, position829)
			}
			return true
		l828:
			position, tokenIndex, depth = position828, tokenIndex828, depth828
			return false
		},
		/* 78 SourceSinkType <- <(<ident> Action60)> */
		func() bool {
			position831, tokenIndex831, depth831 := position, tokenIndex, depth
			{
				position832 := position
				depth++
				{
					position833 := position
					depth++
					if !_rules[ruleident]() {
						goto l831
					}
					depth--
					add(rulePegText, position833)
				}
				if !_rules[ruleAction60]() {
					goto l831
				}
				depth--
				add(ruleSourceSinkType, position832)
			}
			return true
		l831:
			position, tokenIndex, depth = position831, tokenIndex831, depth831
			return false
		},
		/* 79 SourceSinkParamKey <- <(<ident> Action61)> */
		func() bool {
			position834, tokenIndex834, depth834 := position, tokenIndex, depth
			{
				position835 := position
				depth++
				{
					position836 := position
					depth++
					if !_rules[ruleident]() {
						goto l834
					}
					depth--
					add(rulePegText, position836)
				}
				if !_rules[ruleAction61]() {
					goto l834
				}
				depth--
				add(ruleSourceSinkParamKey, position835)
			}
			return true
		l834:
			position, tokenIndex, depth = position834, tokenIndex834, depth834
			return false
		},
		/* 80 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action62)> */
		func() bool {
			position837, tokenIndex837, depth837 := position, tokenIndex, depth
			{
				position838 := position
				depth++
				{
					position839 := position
					depth++
					{
						position840, tokenIndex840, depth840 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l841
						}
						position++
						goto l840
					l841:
						position, tokenIndex, depth = position840, tokenIndex840, depth840
						if buffer[position] != rune('P') {
							goto l837
						}
						position++
					}
				l840:
					{
						position842, tokenIndex842, depth842 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l843
						}
						position++
						goto l842
					l843:
						position, tokenIndex, depth = position842, tokenIndex842, depth842
						if buffer[position] != rune('A') {
							goto l837
						}
						position++
					}
				l842:
					{
						position844, tokenIndex844, depth844 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l845
						}
						position++
						goto l844
					l845:
						position, tokenIndex, depth = position844, tokenIndex844, depth844
						if buffer[position] != rune('U') {
							goto l837
						}
						position++
					}
				l844:
					{
						position846, tokenIndex846, depth846 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l847
						}
						position++
						goto l846
					l847:
						position, tokenIndex, depth = position846, tokenIndex846, depth846
						if buffer[position] != rune('S') {
							goto l837
						}
						position++
					}
				l846:
					{
						position848, tokenIndex848, depth848 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l849
						}
						position++
						goto l848
					l849:
						position, tokenIndex, depth = position848, tokenIndex848, depth848
						if buffer[position] != rune('E') {
							goto l837
						}
						position++
					}
				l848:
					{
						position850, tokenIndex850, depth850 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l851
						}
						position++
						goto l850
					l851:
						position, tokenIndex, depth = position850, tokenIndex850, depth850
						if buffer[position] != rune('D') {
							goto l837
						}
						position++
					}
				l850:
					depth--
					add(rulePegText, position839)
				}
				if !_rules[ruleAction62]() {
					goto l837
				}
				depth--
				add(rulePaused, position838)
			}
			return true
		l837:
			position, tokenIndex, depth = position837, tokenIndex837, depth837
			return false
		},
		/* 81 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action63)> */
		func() bool {
			position852, tokenIndex852, depth852 := position, tokenIndex, depth
			{
				position853 := position
				depth++
				{
					position854 := position
					depth++
					{
						position855, tokenIndex855, depth855 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l856
						}
						position++
						goto l855
					l856:
						position, tokenIndex, depth = position855, tokenIndex855, depth855
						if buffer[position] != rune('U') {
							goto l852
						}
						position++
					}
				l855:
					{
						position857, tokenIndex857, depth857 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l858
						}
						position++
						goto l857
					l858:
						position, tokenIndex, depth = position857, tokenIndex857, depth857
						if buffer[position] != rune('N') {
							goto l852
						}
						position++
					}
				l857:
					{
						position859, tokenIndex859, depth859 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l860
						}
						position++
						goto l859
					l860:
						position, tokenIndex, depth = position859, tokenIndex859, depth859
						if buffer[position] != rune('P') {
							goto l852
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
							goto l852
						}
						position++
					}
				l861:
					{
						position863, tokenIndex863, depth863 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l864
						}
						position++
						goto l863
					l864:
						position, tokenIndex, depth = position863, tokenIndex863, depth863
						if buffer[position] != rune('U') {
							goto l852
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
							goto l852
						}
						position++
					}
				l865:
					{
						position867, tokenIndex867, depth867 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l868
						}
						position++
						goto l867
					l868:
						position, tokenIndex, depth = position867, tokenIndex867, depth867
						if buffer[position] != rune('E') {
							goto l852
						}
						position++
					}
				l867:
					{
						position869, tokenIndex869, depth869 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l870
						}
						position++
						goto l869
					l870:
						position, tokenIndex, depth = position869, tokenIndex869, depth869
						if buffer[position] != rune('D') {
							goto l852
						}
						position++
					}
				l869:
					depth--
					add(rulePegText, position854)
				}
				if !_rules[ruleAction63]() {
					goto l852
				}
				depth--
				add(ruleUnpaused, position853)
			}
			return true
		l852:
			position, tokenIndex, depth = position852, tokenIndex852, depth852
			return false
		},
		/* 82 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action64)> */
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
						if buffer[position] != rune('o') {
							goto l875
						}
						position++
						goto l874
					l875:
						position, tokenIndex, depth = position874, tokenIndex874, depth874
						if buffer[position] != rune('O') {
							goto l871
						}
						position++
					}
				l874:
					{
						position876, tokenIndex876, depth876 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l877
						}
						position++
						goto l876
					l877:
						position, tokenIndex, depth = position876, tokenIndex876, depth876
						if buffer[position] != rune('R') {
							goto l871
						}
						position++
					}
				l876:
					depth--
					add(rulePegText, position873)
				}
				if !_rules[ruleAction64]() {
					goto l871
				}
				depth--
				add(ruleOr, position872)
			}
			return true
		l871:
			position, tokenIndex, depth = position871, tokenIndex871, depth871
			return false
		},
		/* 83 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action65)> */
		func() bool {
			position878, tokenIndex878, depth878 := position, tokenIndex, depth
			{
				position879 := position
				depth++
				{
					position880 := position
					depth++
					{
						position881, tokenIndex881, depth881 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l882
						}
						position++
						goto l881
					l882:
						position, tokenIndex, depth = position881, tokenIndex881, depth881
						if buffer[position] != rune('A') {
							goto l878
						}
						position++
					}
				l881:
					{
						position883, tokenIndex883, depth883 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l884
						}
						position++
						goto l883
					l884:
						position, tokenIndex, depth = position883, tokenIndex883, depth883
						if buffer[position] != rune('N') {
							goto l878
						}
						position++
					}
				l883:
					{
						position885, tokenIndex885, depth885 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l886
						}
						position++
						goto l885
					l886:
						position, tokenIndex, depth = position885, tokenIndex885, depth885
						if buffer[position] != rune('D') {
							goto l878
						}
						position++
					}
				l885:
					depth--
					add(rulePegText, position880)
				}
				if !_rules[ruleAction65]() {
					goto l878
				}
				depth--
				add(ruleAnd, position879)
			}
			return true
		l878:
			position, tokenIndex, depth = position878, tokenIndex878, depth878
			return false
		},
		/* 84 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action66)> */
		func() bool {
			position887, tokenIndex887, depth887 := position, tokenIndex, depth
			{
				position888 := position
				depth++
				{
					position889 := position
					depth++
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
							goto l887
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
							goto l887
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
							goto l887
						}
						position++
					}
				l894:
					depth--
					add(rulePegText, position889)
				}
				if !_rules[ruleAction66]() {
					goto l887
				}
				depth--
				add(ruleNot, position888)
			}
			return true
		l887:
			position, tokenIndex, depth = position887, tokenIndex887, depth887
			return false
		},
		/* 85 Equal <- <(<'='> Action67)> */
		func() bool {
			position896, tokenIndex896, depth896 := position, tokenIndex, depth
			{
				position897 := position
				depth++
				{
					position898 := position
					depth++
					if buffer[position] != rune('=') {
						goto l896
					}
					position++
					depth--
					add(rulePegText, position898)
				}
				if !_rules[ruleAction67]() {
					goto l896
				}
				depth--
				add(ruleEqual, position897)
			}
			return true
		l896:
			position, tokenIndex, depth = position896, tokenIndex896, depth896
			return false
		},
		/* 86 Less <- <(<'<'> Action68)> */
		func() bool {
			position899, tokenIndex899, depth899 := position, tokenIndex, depth
			{
				position900 := position
				depth++
				{
					position901 := position
					depth++
					if buffer[position] != rune('<') {
						goto l899
					}
					position++
					depth--
					add(rulePegText, position901)
				}
				if !_rules[ruleAction68]() {
					goto l899
				}
				depth--
				add(ruleLess, position900)
			}
			return true
		l899:
			position, tokenIndex, depth = position899, tokenIndex899, depth899
			return false
		},
		/* 87 LessOrEqual <- <(<('<' '=')> Action69)> */
		func() bool {
			position902, tokenIndex902, depth902 := position, tokenIndex, depth
			{
				position903 := position
				depth++
				{
					position904 := position
					depth++
					if buffer[position] != rune('<') {
						goto l902
					}
					position++
					if buffer[position] != rune('=') {
						goto l902
					}
					position++
					depth--
					add(rulePegText, position904)
				}
				if !_rules[ruleAction69]() {
					goto l902
				}
				depth--
				add(ruleLessOrEqual, position903)
			}
			return true
		l902:
			position, tokenIndex, depth = position902, tokenIndex902, depth902
			return false
		},
		/* 88 Greater <- <(<'>'> Action70)> */
		func() bool {
			position905, tokenIndex905, depth905 := position, tokenIndex, depth
			{
				position906 := position
				depth++
				{
					position907 := position
					depth++
					if buffer[position] != rune('>') {
						goto l905
					}
					position++
					depth--
					add(rulePegText, position907)
				}
				if !_rules[ruleAction70]() {
					goto l905
				}
				depth--
				add(ruleGreater, position906)
			}
			return true
		l905:
			position, tokenIndex, depth = position905, tokenIndex905, depth905
			return false
		},
		/* 89 GreaterOrEqual <- <(<('>' '=')> Action71)> */
		func() bool {
			position908, tokenIndex908, depth908 := position, tokenIndex, depth
			{
				position909 := position
				depth++
				{
					position910 := position
					depth++
					if buffer[position] != rune('>') {
						goto l908
					}
					position++
					if buffer[position] != rune('=') {
						goto l908
					}
					position++
					depth--
					add(rulePegText, position910)
				}
				if !_rules[ruleAction71]() {
					goto l908
				}
				depth--
				add(ruleGreaterOrEqual, position909)
			}
			return true
		l908:
			position, tokenIndex, depth = position908, tokenIndex908, depth908
			return false
		},
		/* 90 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action72)> */
		func() bool {
			position911, tokenIndex911, depth911 := position, tokenIndex, depth
			{
				position912 := position
				depth++
				{
					position913 := position
					depth++
					{
						position914, tokenIndex914, depth914 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l915
						}
						position++
						if buffer[position] != rune('=') {
							goto l915
						}
						position++
						goto l914
					l915:
						position, tokenIndex, depth = position914, tokenIndex914, depth914
						if buffer[position] != rune('<') {
							goto l911
						}
						position++
						if buffer[position] != rune('>') {
							goto l911
						}
						position++
					}
				l914:
					depth--
					add(rulePegText, position913)
				}
				if !_rules[ruleAction72]() {
					goto l911
				}
				depth--
				add(ruleNotEqual, position912)
			}
			return true
		l911:
			position, tokenIndex, depth = position911, tokenIndex911, depth911
			return false
		},
		/* 91 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action73)> */
		func() bool {
			position916, tokenIndex916, depth916 := position, tokenIndex, depth
			{
				position917 := position
				depth++
				{
					position918 := position
					depth++
					{
						position919, tokenIndex919, depth919 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l920
						}
						position++
						goto l919
					l920:
						position, tokenIndex, depth = position919, tokenIndex919, depth919
						if buffer[position] != rune('I') {
							goto l916
						}
						position++
					}
				l919:
					{
						position921, tokenIndex921, depth921 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l922
						}
						position++
						goto l921
					l922:
						position, tokenIndex, depth = position921, tokenIndex921, depth921
						if buffer[position] != rune('S') {
							goto l916
						}
						position++
					}
				l921:
					depth--
					add(rulePegText, position918)
				}
				if !_rules[ruleAction73]() {
					goto l916
				}
				depth--
				add(ruleIs, position917)
			}
			return true
		l916:
			position, tokenIndex, depth = position916, tokenIndex916, depth916
			return false
		},
		/* 92 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action74)> */
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
						if buffer[position] != rune('i') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('I') {
							goto l923
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
							goto l923
						}
						position++
					}
				l928:
					if !_rules[rulesp]() {
						goto l923
					}
					{
						position930, tokenIndex930, depth930 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l931
						}
						position++
						goto l930
					l931:
						position, tokenIndex, depth = position930, tokenIndex930, depth930
						if buffer[position] != rune('N') {
							goto l923
						}
						position++
					}
				l930:
					{
						position932, tokenIndex932, depth932 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l933
						}
						position++
						goto l932
					l933:
						position, tokenIndex, depth = position932, tokenIndex932, depth932
						if buffer[position] != rune('O') {
							goto l923
						}
						position++
					}
				l932:
					{
						position934, tokenIndex934, depth934 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l935
						}
						position++
						goto l934
					l935:
						position, tokenIndex, depth = position934, tokenIndex934, depth934
						if buffer[position] != rune('T') {
							goto l923
						}
						position++
					}
				l934:
					depth--
					add(rulePegText, position925)
				}
				if !_rules[ruleAction74]() {
					goto l923
				}
				depth--
				add(ruleIsNot, position924)
			}
			return true
		l923:
			position, tokenIndex, depth = position923, tokenIndex923, depth923
			return false
		},
		/* 93 Plus <- <(<'+'> Action75)> */
		func() bool {
			position936, tokenIndex936, depth936 := position, tokenIndex, depth
			{
				position937 := position
				depth++
				{
					position938 := position
					depth++
					if buffer[position] != rune('+') {
						goto l936
					}
					position++
					depth--
					add(rulePegText, position938)
				}
				if !_rules[ruleAction75]() {
					goto l936
				}
				depth--
				add(rulePlus, position937)
			}
			return true
		l936:
			position, tokenIndex, depth = position936, tokenIndex936, depth936
			return false
		},
		/* 94 Minus <- <(<'-'> Action76)> */
		func() bool {
			position939, tokenIndex939, depth939 := position, tokenIndex, depth
			{
				position940 := position
				depth++
				{
					position941 := position
					depth++
					if buffer[position] != rune('-') {
						goto l939
					}
					position++
					depth--
					add(rulePegText, position941)
				}
				if !_rules[ruleAction76]() {
					goto l939
				}
				depth--
				add(ruleMinus, position940)
			}
			return true
		l939:
			position, tokenIndex, depth = position939, tokenIndex939, depth939
			return false
		},
		/* 95 Multiply <- <(<'*'> Action77)> */
		func() bool {
			position942, tokenIndex942, depth942 := position, tokenIndex, depth
			{
				position943 := position
				depth++
				{
					position944 := position
					depth++
					if buffer[position] != rune('*') {
						goto l942
					}
					position++
					depth--
					add(rulePegText, position944)
				}
				if !_rules[ruleAction77]() {
					goto l942
				}
				depth--
				add(ruleMultiply, position943)
			}
			return true
		l942:
			position, tokenIndex, depth = position942, tokenIndex942, depth942
			return false
		},
		/* 96 Divide <- <(<'/'> Action78)> */
		func() bool {
			position945, tokenIndex945, depth945 := position, tokenIndex, depth
			{
				position946 := position
				depth++
				{
					position947 := position
					depth++
					if buffer[position] != rune('/') {
						goto l945
					}
					position++
					depth--
					add(rulePegText, position947)
				}
				if !_rules[ruleAction78]() {
					goto l945
				}
				depth--
				add(ruleDivide, position946)
			}
			return true
		l945:
			position, tokenIndex, depth = position945, tokenIndex945, depth945
			return false
		},
		/* 97 Modulo <- <(<'%'> Action79)> */
		func() bool {
			position948, tokenIndex948, depth948 := position, tokenIndex, depth
			{
				position949 := position
				depth++
				{
					position950 := position
					depth++
					if buffer[position] != rune('%') {
						goto l948
					}
					position++
					depth--
					add(rulePegText, position950)
				}
				if !_rules[ruleAction79]() {
					goto l948
				}
				depth--
				add(ruleModulo, position949)
			}
			return true
		l948:
			position, tokenIndex, depth = position948, tokenIndex948, depth948
			return false
		},
		/* 98 UnaryMinus <- <(<'-'> Action80)> */
		func() bool {
			position951, tokenIndex951, depth951 := position, tokenIndex, depth
			{
				position952 := position
				depth++
				{
					position953 := position
					depth++
					if buffer[position] != rune('-') {
						goto l951
					}
					position++
					depth--
					add(rulePegText, position953)
				}
				if !_rules[ruleAction80]() {
					goto l951
				}
				depth--
				add(ruleUnaryMinus, position952)
			}
			return true
		l951:
			position, tokenIndex, depth = position951, tokenIndex951, depth951
			return false
		},
		/* 99 Identifier <- <(<ident> Action81)> */
		func() bool {
			position954, tokenIndex954, depth954 := position, tokenIndex, depth
			{
				position955 := position
				depth++
				{
					position956 := position
					depth++
					if !_rules[ruleident]() {
						goto l954
					}
					depth--
					add(rulePegText, position956)
				}
				if !_rules[ruleAction81]() {
					goto l954
				}
				depth--
				add(ruleIdentifier, position955)
			}
			return true
		l954:
			position, tokenIndex, depth = position954, tokenIndex954, depth954
			return false
		},
		/* 100 TargetIdentifier <- <(<jsonPath> Action82)> */
		func() bool {
			position957, tokenIndex957, depth957 := position, tokenIndex, depth
			{
				position958 := position
				depth++
				{
					position959 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l957
					}
					depth--
					add(rulePegText, position959)
				}
				if !_rules[ruleAction82]() {
					goto l957
				}
				depth--
				add(ruleTargetIdentifier, position958)
			}
			return true
		l957:
			position, tokenIndex, depth = position957, tokenIndex957, depth957
			return false
		},
		/* 101 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position960, tokenIndex960, depth960 := position, tokenIndex, depth
			{
				position961 := position
				depth++
				{
					position962, tokenIndex962, depth962 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l963
					}
					position++
					goto l962
				l963:
					position, tokenIndex, depth = position962, tokenIndex962, depth962
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l960
					}
					position++
				}
			l962:
			l964:
				{
					position965, tokenIndex965, depth965 := position, tokenIndex, depth
					{
						position966, tokenIndex966, depth966 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l967
						}
						position++
						goto l966
					l967:
						position, tokenIndex, depth = position966, tokenIndex966, depth966
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l968
						}
						position++
						goto l966
					l968:
						position, tokenIndex, depth = position966, tokenIndex966, depth966
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l969
						}
						position++
						goto l966
					l969:
						position, tokenIndex, depth = position966, tokenIndex966, depth966
						if buffer[position] != rune('_') {
							goto l965
						}
						position++
					}
				l966:
					goto l964
				l965:
					position, tokenIndex, depth = position965, tokenIndex965, depth965
				}
				depth--
				add(ruleident, position961)
			}
			return true
		l960:
			position, tokenIndex, depth = position960, tokenIndex960, depth960
			return false
		},
		/* 102 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position970, tokenIndex970, depth970 := position, tokenIndex, depth
			{
				position971 := position
				depth++
				{
					position972, tokenIndex972, depth972 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l973
					}
					position++
					goto l972
				l973:
					position, tokenIndex, depth = position972, tokenIndex972, depth972
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l970
					}
					position++
				}
			l972:
			l974:
				{
					position975, tokenIndex975, depth975 := position, tokenIndex, depth
					{
						position976, tokenIndex976, depth976 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l977
						}
						position++
						goto l976
					l977:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l978
						}
						position++
						goto l976
					l978:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l979
						}
						position++
						goto l976
					l979:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('_') {
							goto l980
						}
						position++
						goto l976
					l980:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('.') {
							goto l981
						}
						position++
						goto l976
					l981:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('[') {
							goto l982
						}
						position++
						goto l976
					l982:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune(']') {
							goto l983
						}
						position++
						goto l976
					l983:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('"') {
							goto l975
						}
						position++
					}
				l976:
					goto l974
				l975:
					position, tokenIndex, depth = position975, tokenIndex975, depth975
				}
				depth--
				add(rulejsonPath, position971)
			}
			return true
		l970:
			position, tokenIndex, depth = position970, tokenIndex970, depth970
			return false
		},
		/* 103 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position985 := position
				depth++
			l986:
				{
					position987, tokenIndex987, depth987 := position, tokenIndex, depth
					{
						position988, tokenIndex988, depth988 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l989
						}
						position++
						goto l988
					l989:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('\t') {
							goto l990
						}
						position++
						goto l988
					l990:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('\n') {
							goto l991
						}
						position++
						goto l988
					l991:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('\r') {
							goto l992
						}
						position++
						goto l988
					l992:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if !_rules[rulecomment]() {
							goto l987
						}
					}
				l988:
					goto l986
				l987:
					position, tokenIndex, depth = position987, tokenIndex987, depth987
				}
				depth--
				add(rulesp, position985)
			}
			return true
		},
		/* 104 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position993, tokenIndex993, depth993 := position, tokenIndex, depth
			{
				position994 := position
				depth++
				if buffer[position] != rune('-') {
					goto l993
				}
				position++
				if buffer[position] != rune('-') {
					goto l993
				}
				position++
			l995:
				{
					position996, tokenIndex996, depth996 := position, tokenIndex, depth
					{
						position997, tokenIndex997, depth997 := position, tokenIndex, depth
						{
							position998, tokenIndex998, depth998 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l999
							}
							position++
							goto l998
						l999:
							position, tokenIndex, depth = position998, tokenIndex998, depth998
							if buffer[position] != rune('\n') {
								goto l997
							}
							position++
						}
					l998:
						goto l996
					l997:
						position, tokenIndex, depth = position997, tokenIndex997, depth997
					}
					if !matchDot() {
						goto l996
					}
					goto l995
				l996:
					position, tokenIndex, depth = position996, tokenIndex996, depth996
				}
				{
					position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1001
					}
					position++
					goto l1000
				l1001:
					position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
					if buffer[position] != rune('\n') {
						goto l993
					}
					position++
				}
			l1000:
				depth--
				add(rulecomment, position994)
			}
			return true
		l993:
			position, tokenIndex, depth = position993, tokenIndex993, depth993
			return false
		},
		/* 106 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 107 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 108 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 109 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 110 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 111 Action5 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 112 Action6 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 113 Action7 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 114 Action8 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 115 Action9 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 116 Action10 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 117 Action11 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 118 Action12 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 119 Action13 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		nil,
		/* 121 Action14 <- <{
		    p.AssembleEmitter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 122 Action15 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 123 Action16 <- <{
		    p.PushComponent(end, end, NewStream("*"))
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 124 Action17 <- <{
		    p.AssembleStreamEmitInterval()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 125 Action18 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 126 Action19 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 127 Action20 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 128 Action21 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 129 Action22 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 130 Action23 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 131 Action24 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 132 Action25 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 133 Action26 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 134 Action27 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 135 Action28 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 136 Action29 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 137 Action30 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 138 Action31 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 139 Action32 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 140 Action33 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 141 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 142 Action35 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 143 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 144 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 145 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 146 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 147 Action40 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 148 Action41 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 149 Action42 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 150 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 151 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 152 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 153 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 154 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 155 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 156 Action49 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 157 Action50 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 158 Action51 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 159 Action52 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 160 Action53 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 161 Action54 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 162 Action55 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 163 Action56 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 164 Action57 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 165 Action58 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 166 Action59 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 167 Action60 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 168 Action61 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 169 Action62 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 170 Action63 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 171 Action64 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 172 Action65 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 173 Action66 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 174 Action67 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 175 Action68 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 176 Action69 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 177 Action70 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 178 Action71 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 179 Action72 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 180 Action73 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 181 Action74 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 182 Action75 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 183 Action76 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 184 Action77 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 185 Action78 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 186 Action79 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 187 Action80 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 188 Action81 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 189 Action82 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
	}
	p.rules = _rules
}

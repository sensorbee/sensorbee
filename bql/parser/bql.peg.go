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
	ruleotherOpExpr
	ruleisExpr
	ruletermExpr
	ruleproductExpr
	ruleminusExpr
	rulecastExpr
	rulebaseExpr
	ruleFuncTypeCast
	ruleFuncApp
	ruleFuncParams
	ruleLiteral
	ruleComparisonOp
	ruleOtherOp
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
	ruleType
	ruleBool
	ruleInt
	ruleFloat
	ruleString
	ruleBlob
	ruleTimestamp
	ruleArray
	ruleMap
	ruleOr
	ruleAnd
	ruleNot
	ruleEqual
	ruleLess
	ruleLessOrEqual
	ruleGreater
	ruleGreaterOrEqual
	ruleNotEqual
	ruleConcat
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
	ruleAction17
	rulePegText
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
	ruleAction87
	ruleAction88
	ruleAction89
	ruleAction90
	ruleAction91
	ruleAction92
	ruleAction93
	ruleAction94
	ruleAction95

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
	"otherOpExpr",
	"isExpr",
	"termExpr",
	"productExpr",
	"minusExpr",
	"castExpr",
	"baseExpr",
	"FuncTypeCast",
	"FuncApp",
	"FuncParams",
	"Literal",
	"ComparisonOp",
	"OtherOp",
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
	"Type",
	"Bool",
	"Int",
	"Float",
	"String",
	"Blob",
	"Timestamp",
	"Array",
	"Map",
	"Or",
	"And",
	"Not",
	"Equal",
	"Less",
	"LessOrEqual",
	"Greater",
	"GreaterOrEqual",
	"NotEqual",
	"Concat",
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
	"Action17",
	"PegText",
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
	"Action87",
	"Action88",
	"Action89",
	"Action90",
	"Action91",
	"Action92",
	"Action93",
	"Action94",
	"Action95",

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
	rules  [221]func() bool
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

			p.AssembleEmitter()

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

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction32:

			p.AssembleSourceSinkParam()

		case ruleAction33:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction34:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

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

			p.AssembleTypeCast(begin, end)

		case ruleAction44:

			p.AssembleTypeCast(begin, end)

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

			p.PushComponent(begin, end, Bool)

		case ruleAction69:

			p.PushComponent(begin, end, Int)

		case ruleAction70:

			p.PushComponent(begin, end, Float)

		case ruleAction71:

			p.PushComponent(begin, end, String)

		case ruleAction72:

			p.PushComponent(begin, end, Blob)

		case ruleAction73:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction74:

			p.PushComponent(begin, end, Array)

		case ruleAction75:

			p.PushComponent(begin, end, Map)

		case ruleAction76:

			p.PushComponent(begin, end, Or)

		case ruleAction77:

			p.PushComponent(begin, end, And)

		case ruleAction78:

			p.PushComponent(begin, end, Not)

		case ruleAction79:

			p.PushComponent(begin, end, Equal)

		case ruleAction80:

			p.PushComponent(begin, end, Less)

		case ruleAction81:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction82:

			p.PushComponent(begin, end, Greater)

		case ruleAction83:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction84:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction85:

			p.PushComponent(begin, end, Concat)

		case ruleAction86:

			p.PushComponent(begin, end, Is)

		case ruleAction87:

			p.PushComponent(begin, end, IsNot)

		case ruleAction88:

			p.PushComponent(begin, end, Plus)

		case ruleAction89:

			p.PushComponent(begin, end, Minus)

		case ruleAction90:

			p.PushComponent(begin, end, Multiply)

		case ruleAction91:

			p.PushComponent(begin, end, Divide)

		case ruleAction92:

			p.PushComponent(begin, end, Modulo)

		case ruleAction93:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction94:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction95:

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
		/* 7 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action1)> */
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
				if !_rules[ruleSelectStmt]() {
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
		/* 9 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
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
		/* 10 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
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
		/* 11 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action5)> */
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
		/* 12 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action6)> */
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
		/* 13 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action7)> */
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
		/* 14 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action8)> */
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
		/* 15 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action9)> */
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
		/* 16 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action10)> */
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
		/* 17 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action11)> */
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
		/* 18 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action12)> */
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
		/* 19 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
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
		/* 20 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action14)> */
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
		/* 21 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action15)> */
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
		/* 22 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action16)> */
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
		/* 23 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) Action17)> */
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
		/* 24 Projections <- <(<(Projection sp (',' sp Projection)*)> Action18)> */
		func() bool {
			position465, tokenIndex465, depth465 := position, tokenIndex, depth
			{
				position466 := position
				depth++
				{
					position467 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l465
					}
					if !_rules[rulesp]() {
						goto l465
					}
				l468:
					{
						position469, tokenIndex469, depth469 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l469
						}
						position++
						if !_rules[rulesp]() {
							goto l469
						}
						if !_rules[ruleProjection]() {
							goto l469
						}
						goto l468
					l469:
						position, tokenIndex, depth = position469, tokenIndex469, depth469
					}
					depth--
					add(rulePegText, position467)
				}
				if !_rules[ruleAction18]() {
					goto l465
				}
				depth--
				add(ruleProjections, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 25 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position470, tokenIndex470, depth470 := position, tokenIndex, depth
			{
				position471 := position
				depth++
				{
					position472, tokenIndex472, depth472 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l473
					}
					goto l472
				l473:
					position, tokenIndex, depth = position472, tokenIndex472, depth472
					if !_rules[ruleExpression]() {
						goto l474
					}
					goto l472
				l474:
					position, tokenIndex, depth = position472, tokenIndex472, depth472
					if !_rules[ruleWildcard]() {
						goto l470
					}
				}
			l472:
				depth--
				add(ruleProjection, position471)
			}
			return true
		l470:
			position, tokenIndex, depth = position470, tokenIndex470, depth470
			return false
		},
		/* 26 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action19)> */
		func() bool {
			position475, tokenIndex475, depth475 := position, tokenIndex, depth
			{
				position476 := position
				depth++
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l478
					}
					goto l477
				l478:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
					if !_rules[ruleWildcard]() {
						goto l475
					}
				}
			l477:
				if !_rules[rulesp]() {
					goto l475
				}
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l480
					}
					position++
					goto l479
				l480:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if buffer[position] != rune('A') {
						goto l475
					}
					position++
				}
			l479:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if buffer[position] != rune('S') {
						goto l475
					}
					position++
				}
			l481:
				if !_rules[rulesp]() {
					goto l475
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l475
				}
				if !_rules[ruleAction19]() {
					goto l475
				}
				depth--
				add(ruleAliasExpression, position476)
			}
			return true
		l475:
			position, tokenIndex, depth = position475, tokenIndex475, depth475
			return false
		},
		/* 27 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action20)> */
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
							if buffer[position] != rune('f') {
								goto l489
							}
							position++
							goto l488
						l489:
							position, tokenIndex, depth = position488, tokenIndex488, depth488
							if buffer[position] != rune('F') {
								goto l486
							}
							position++
						}
					l488:
						{
							position490, tokenIndex490, depth490 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l491
							}
							position++
							goto l490
						l491:
							position, tokenIndex, depth = position490, tokenIndex490, depth490
							if buffer[position] != rune('R') {
								goto l486
							}
							position++
						}
					l490:
						{
							position492, tokenIndex492, depth492 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l493
							}
							position++
							goto l492
						l493:
							position, tokenIndex, depth = position492, tokenIndex492, depth492
							if buffer[position] != rune('O') {
								goto l486
							}
							position++
						}
					l492:
						{
							position494, tokenIndex494, depth494 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l495
							}
							position++
							goto l494
						l495:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							if buffer[position] != rune('M') {
								goto l486
							}
							position++
						}
					l494:
						if !_rules[rulesp]() {
							goto l486
						}
						if !_rules[ruleRelations]() {
							goto l486
						}
						if !_rules[rulesp]() {
							goto l486
						}
						goto l487
					l486:
						position, tokenIndex, depth = position486, tokenIndex486, depth486
					}
				l487:
					depth--
					add(rulePegText, position485)
				}
				if !_rules[ruleAction20]() {
					goto l483
				}
				depth--
				add(ruleWindowedFrom, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 28 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position496, tokenIndex496, depth496 := position, tokenIndex, depth
			{
				position497 := position
				depth++
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l499
					}
					goto l498
				l499:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleTuplesInterval]() {
						goto l496
					}
				}
			l498:
				depth--
				add(ruleInterval, position497)
			}
			return true
		l496:
			position, tokenIndex, depth = position496, tokenIndex496, depth496
			return false
		},
		/* 29 TimeInterval <- <(NumericLiteral sp SECONDS Action21)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l500
				}
				if !_rules[rulesp]() {
					goto l500
				}
				if !_rules[ruleSECONDS]() {
					goto l500
				}
				if !_rules[ruleAction21]() {
					goto l500
				}
				depth--
				add(ruleTimeInterval, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 30 TuplesInterval <- <(NumericLiteral sp TUPLES Action22)> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l502
				}
				if !_rules[rulesp]() {
					goto l502
				}
				if !_rules[ruleTUPLES]() {
					goto l502
				}
				if !_rules[ruleAction22]() {
					goto l502
				}
				depth--
				add(ruleTuplesInterval, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 31 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position504, tokenIndex504, depth504 := position, tokenIndex, depth
			{
				position505 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l504
				}
				if !_rules[rulesp]() {
					goto l504
				}
			l506:
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l507
					}
					position++
					if !_rules[rulesp]() {
						goto l507
					}
					if !_rules[ruleRelationLike]() {
						goto l507
					}
					goto l506
				l507:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
				}
				depth--
				add(ruleRelations, position505)
			}
			return true
		l504:
			position, tokenIndex, depth = position504, tokenIndex504, depth504
			return false
		},
		/* 32 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action23)> */
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
							if buffer[position] != rune('h') {
								goto l516
							}
							position++
							goto l515
						l516:
							position, tokenIndex, depth = position515, tokenIndex515, depth515
							if buffer[position] != rune('H') {
								goto l511
							}
							position++
						}
					l515:
						{
							position517, tokenIndex517, depth517 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l518
							}
							position++
							goto l517
						l518:
							position, tokenIndex, depth = position517, tokenIndex517, depth517
							if buffer[position] != rune('E') {
								goto l511
							}
							position++
						}
					l517:
						{
							position519, tokenIndex519, depth519 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l520
							}
							position++
							goto l519
						l520:
							position, tokenIndex, depth = position519, tokenIndex519, depth519
							if buffer[position] != rune('R') {
								goto l511
							}
							position++
						}
					l519:
						{
							position521, tokenIndex521, depth521 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l522
							}
							position++
							goto l521
						l522:
							position, tokenIndex, depth = position521, tokenIndex521, depth521
							if buffer[position] != rune('E') {
								goto l511
							}
							position++
						}
					l521:
						if !_rules[rulesp]() {
							goto l511
						}
						if !_rules[ruleExpression]() {
							goto l511
						}
						goto l512
					l511:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
					}
				l512:
					depth--
					add(rulePegText, position510)
				}
				if !_rules[ruleAction23]() {
					goto l508
				}
				depth--
				add(ruleFilter, position509)
			}
			return true
		l508:
			position, tokenIndex, depth = position508, tokenIndex508, depth508
			return false
		},
		/* 33 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action24)> */
		func() bool {
			position523, tokenIndex523, depth523 := position, tokenIndex, depth
			{
				position524 := position
				depth++
				{
					position525 := position
					depth++
					{
						position526, tokenIndex526, depth526 := position, tokenIndex, depth
						{
							position528, tokenIndex528, depth528 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l529
							}
							position++
							goto l528
						l529:
							position, tokenIndex, depth = position528, tokenIndex528, depth528
							if buffer[position] != rune('G') {
								goto l526
							}
							position++
						}
					l528:
						{
							position530, tokenIndex530, depth530 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l531
							}
							position++
							goto l530
						l531:
							position, tokenIndex, depth = position530, tokenIndex530, depth530
							if buffer[position] != rune('R') {
								goto l526
							}
							position++
						}
					l530:
						{
							position532, tokenIndex532, depth532 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l533
							}
							position++
							goto l532
						l533:
							position, tokenIndex, depth = position532, tokenIndex532, depth532
							if buffer[position] != rune('O') {
								goto l526
							}
							position++
						}
					l532:
						{
							position534, tokenIndex534, depth534 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l535
							}
							position++
							goto l534
						l535:
							position, tokenIndex, depth = position534, tokenIndex534, depth534
							if buffer[position] != rune('U') {
								goto l526
							}
							position++
						}
					l534:
						{
							position536, tokenIndex536, depth536 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l537
							}
							position++
							goto l536
						l537:
							position, tokenIndex, depth = position536, tokenIndex536, depth536
							if buffer[position] != rune('P') {
								goto l526
							}
							position++
						}
					l536:
						if !_rules[rulesp]() {
							goto l526
						}
						{
							position538, tokenIndex538, depth538 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l539
							}
							position++
							goto l538
						l539:
							position, tokenIndex, depth = position538, tokenIndex538, depth538
							if buffer[position] != rune('B') {
								goto l526
							}
							position++
						}
					l538:
						{
							position540, tokenIndex540, depth540 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l541
							}
							position++
							goto l540
						l541:
							position, tokenIndex, depth = position540, tokenIndex540, depth540
							if buffer[position] != rune('Y') {
								goto l526
							}
							position++
						}
					l540:
						if !_rules[rulesp]() {
							goto l526
						}
						if !_rules[ruleGroupList]() {
							goto l526
						}
						goto l527
					l526:
						position, tokenIndex, depth = position526, tokenIndex526, depth526
					}
				l527:
					depth--
					add(rulePegText, position525)
				}
				if !_rules[ruleAction24]() {
					goto l523
				}
				depth--
				add(ruleGrouping, position524)
			}
			return true
		l523:
			position, tokenIndex, depth = position523, tokenIndex523, depth523
			return false
		},
		/* 34 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position542, tokenIndex542, depth542 := position, tokenIndex, depth
			{
				position543 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l542
				}
				if !_rules[rulesp]() {
					goto l542
				}
			l544:
				{
					position545, tokenIndex545, depth545 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l545
					}
					position++
					if !_rules[rulesp]() {
						goto l545
					}
					if !_rules[ruleExpression]() {
						goto l545
					}
					goto l544
				l545:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
				}
				depth--
				add(ruleGroupList, position543)
			}
			return true
		l542:
			position, tokenIndex, depth = position542, tokenIndex542, depth542
			return false
		},
		/* 35 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action25)> */
		func() bool {
			position546, tokenIndex546, depth546 := position, tokenIndex, depth
			{
				position547 := position
				depth++
				{
					position548 := position
					depth++
					{
						position549, tokenIndex549, depth549 := position, tokenIndex, depth
						{
							position551, tokenIndex551, depth551 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l552
							}
							position++
							goto l551
						l552:
							position, tokenIndex, depth = position551, tokenIndex551, depth551
							if buffer[position] != rune('H') {
								goto l549
							}
							position++
						}
					l551:
						{
							position553, tokenIndex553, depth553 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l554
							}
							position++
							goto l553
						l554:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if buffer[position] != rune('A') {
								goto l549
							}
							position++
						}
					l553:
						{
							position555, tokenIndex555, depth555 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l556
							}
							position++
							goto l555
						l556:
							position, tokenIndex, depth = position555, tokenIndex555, depth555
							if buffer[position] != rune('V') {
								goto l549
							}
							position++
						}
					l555:
						{
							position557, tokenIndex557, depth557 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l558
							}
							position++
							goto l557
						l558:
							position, tokenIndex, depth = position557, tokenIndex557, depth557
							if buffer[position] != rune('I') {
								goto l549
							}
							position++
						}
					l557:
						{
							position559, tokenIndex559, depth559 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l560
							}
							position++
							goto l559
						l560:
							position, tokenIndex, depth = position559, tokenIndex559, depth559
							if buffer[position] != rune('N') {
								goto l549
							}
							position++
						}
					l559:
						{
							position561, tokenIndex561, depth561 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l562
							}
							position++
							goto l561
						l562:
							position, tokenIndex, depth = position561, tokenIndex561, depth561
							if buffer[position] != rune('G') {
								goto l549
							}
							position++
						}
					l561:
						if !_rules[rulesp]() {
							goto l549
						}
						if !_rules[ruleExpression]() {
							goto l549
						}
						goto l550
					l549:
						position, tokenIndex, depth = position549, tokenIndex549, depth549
					}
				l550:
					depth--
					add(rulePegText, position548)
				}
				if !_rules[ruleAction25]() {
					goto l546
				}
				depth--
				add(ruleHaving, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 36 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action26))> */
		func() bool {
			position563, tokenIndex563, depth563 := position, tokenIndex, depth
			{
				position564 := position
				depth++
				{
					position565, tokenIndex565, depth565 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l566
					}
					goto l565
				l566:
					position, tokenIndex, depth = position565, tokenIndex565, depth565
					if !_rules[ruleStreamWindow]() {
						goto l563
					}
					if !_rules[ruleAction26]() {
						goto l563
					}
				}
			l565:
				depth--
				add(ruleRelationLike, position564)
			}
			return true
		l563:
			position, tokenIndex, depth = position563, tokenIndex563, depth563
			return false
		},
		/* 37 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action27)> */
		func() bool {
			position567, tokenIndex567, depth567 := position, tokenIndex, depth
			{
				position568 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l567
				}
				if !_rules[rulesp]() {
					goto l567
				}
				{
					position569, tokenIndex569, depth569 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l570
					}
					position++
					goto l569
				l570:
					position, tokenIndex, depth = position569, tokenIndex569, depth569
					if buffer[position] != rune('A') {
						goto l567
					}
					position++
				}
			l569:
				{
					position571, tokenIndex571, depth571 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l572
					}
					position++
					goto l571
				l572:
					position, tokenIndex, depth = position571, tokenIndex571, depth571
					if buffer[position] != rune('S') {
						goto l567
					}
					position++
				}
			l571:
				if !_rules[rulesp]() {
					goto l567
				}
				if !_rules[ruleIdentifier]() {
					goto l567
				}
				if !_rules[ruleAction27]() {
					goto l567
				}
				depth--
				add(ruleAliasedStreamWindow, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 38 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action28)> */
		func() bool {
			position573, tokenIndex573, depth573 := position, tokenIndex, depth
			{
				position574 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l573
				}
				if !_rules[rulesp]() {
					goto l573
				}
				if buffer[position] != rune('[') {
					goto l573
				}
				position++
				if !_rules[rulesp]() {
					goto l573
				}
				{
					position575, tokenIndex575, depth575 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l576
					}
					position++
					goto l575
				l576:
					position, tokenIndex, depth = position575, tokenIndex575, depth575
					if buffer[position] != rune('R') {
						goto l573
					}
					position++
				}
			l575:
				{
					position577, tokenIndex577, depth577 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l578
					}
					position++
					goto l577
				l578:
					position, tokenIndex, depth = position577, tokenIndex577, depth577
					if buffer[position] != rune('A') {
						goto l573
					}
					position++
				}
			l577:
				{
					position579, tokenIndex579, depth579 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l580
					}
					position++
					goto l579
				l580:
					position, tokenIndex, depth = position579, tokenIndex579, depth579
					if buffer[position] != rune('N') {
						goto l573
					}
					position++
				}
			l579:
				{
					position581, tokenIndex581, depth581 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l582
					}
					position++
					goto l581
				l582:
					position, tokenIndex, depth = position581, tokenIndex581, depth581
					if buffer[position] != rune('G') {
						goto l573
					}
					position++
				}
			l581:
				{
					position583, tokenIndex583, depth583 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l584
					}
					position++
					goto l583
				l584:
					position, tokenIndex, depth = position583, tokenIndex583, depth583
					if buffer[position] != rune('E') {
						goto l573
					}
					position++
				}
			l583:
				if !_rules[rulesp]() {
					goto l573
				}
				if !_rules[ruleInterval]() {
					goto l573
				}
				if !_rules[rulesp]() {
					goto l573
				}
				if buffer[position] != rune(']') {
					goto l573
				}
				position++
				if !_rules[ruleAction28]() {
					goto l573
				}
				depth--
				add(ruleStreamWindow, position574)
			}
			return true
		l573:
			position, tokenIndex, depth = position573, tokenIndex573, depth573
			return false
		},
		/* 39 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position585, tokenIndex585, depth585 := position, tokenIndex, depth
			{
				position586 := position
				depth++
				{
					position587, tokenIndex587, depth587 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l588
					}
					goto l587
				l588:
					position, tokenIndex, depth = position587, tokenIndex587, depth587
					if !_rules[ruleStream]() {
						goto l585
					}
				}
			l587:
				depth--
				add(ruleStreamLike, position586)
			}
			return true
		l585:
			position, tokenIndex, depth = position585, tokenIndex585, depth585
			return false
		},
		/* 40 UDSFFuncApp <- <(FuncApp Action29)> */
		func() bool {
			position589, tokenIndex589, depth589 := position, tokenIndex, depth
			{
				position590 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l589
				}
				if !_rules[ruleAction29]() {
					goto l589
				}
				depth--
				add(ruleUDSFFuncApp, position590)
			}
			return true
		l589:
			position, tokenIndex, depth = position589, tokenIndex589, depth589
			return false
		},
		/* 41 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action30)> */
		func() bool {
			position591, tokenIndex591, depth591 := position, tokenIndex, depth
			{
				position592 := position
				depth++
				{
					position593 := position
					depth++
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						{
							position596, tokenIndex596, depth596 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l597
							}
							position++
							goto l596
						l597:
							position, tokenIndex, depth = position596, tokenIndex596, depth596
							if buffer[position] != rune('W') {
								goto l594
							}
							position++
						}
					l596:
						{
							position598, tokenIndex598, depth598 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l599
							}
							position++
							goto l598
						l599:
							position, tokenIndex, depth = position598, tokenIndex598, depth598
							if buffer[position] != rune('I') {
								goto l594
							}
							position++
						}
					l598:
						{
							position600, tokenIndex600, depth600 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l601
							}
							position++
							goto l600
						l601:
							position, tokenIndex, depth = position600, tokenIndex600, depth600
							if buffer[position] != rune('T') {
								goto l594
							}
							position++
						}
					l600:
						{
							position602, tokenIndex602, depth602 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l603
							}
							position++
							goto l602
						l603:
							position, tokenIndex, depth = position602, tokenIndex602, depth602
							if buffer[position] != rune('H') {
								goto l594
							}
							position++
						}
					l602:
						if !_rules[rulesp]() {
							goto l594
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l594
						}
						if !_rules[rulesp]() {
							goto l594
						}
					l604:
						{
							position605, tokenIndex605, depth605 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l605
							}
							position++
							if !_rules[rulesp]() {
								goto l605
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l605
							}
							goto l604
						l605:
							position, tokenIndex, depth = position605, tokenIndex605, depth605
						}
						goto l595
					l594:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
					}
				l595:
					depth--
					add(rulePegText, position593)
				}
				if !_rules[ruleAction30]() {
					goto l591
				}
				depth--
				add(ruleSourceSinkSpecs, position592)
			}
			return true
		l591:
			position, tokenIndex, depth = position591, tokenIndex591, depth591
			return false
		},
		/* 42 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action31)> */
		func() bool {
			position606, tokenIndex606, depth606 := position, tokenIndex, depth
			{
				position607 := position
				depth++
				{
					position608 := position
					depth++
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('S') {
							goto l606
						}
						position++
					}
				l609:
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('E') {
							goto l606
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('T') {
							goto l606
						}
						position++
					}
				l613:
					if !_rules[rulesp]() {
						goto l606
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l606
					}
					if !_rules[rulesp]() {
						goto l606
					}
				l615:
					{
						position616, tokenIndex616, depth616 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l616
						}
						position++
						if !_rules[rulesp]() {
							goto l616
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l616
						}
						goto l615
					l616:
						position, tokenIndex, depth = position616, tokenIndex616, depth616
					}
					depth--
					add(rulePegText, position608)
				}
				if !_rules[ruleAction31]() {
					goto l606
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position607)
			}
			return true
		l606:
			position, tokenIndex, depth = position606, tokenIndex606, depth606
			return false
		},
		/* 43 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action32)> */
		func() bool {
			position617, tokenIndex617, depth617 := position, tokenIndex, depth
			{
				position618 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l617
				}
				if buffer[position] != rune('=') {
					goto l617
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l617
				}
				if !_rules[ruleAction32]() {
					goto l617
				}
				depth--
				add(ruleSourceSinkParam, position618)
			}
			return true
		l617:
			position, tokenIndex, depth = position617, tokenIndex617, depth617
			return false
		},
		/* 44 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position619, tokenIndex619, depth619 := position, tokenIndex, depth
			{
				position620 := position
				depth++
				{
					position621, tokenIndex621, depth621 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l622
					}
					goto l621
				l622:
					position, tokenIndex, depth = position621, tokenIndex621, depth621
					if !_rules[ruleLiteral]() {
						goto l619
					}
				}
			l621:
				depth--
				add(ruleSourceSinkParamVal, position620)
			}
			return true
		l619:
			position, tokenIndex, depth = position619, tokenIndex619, depth619
			return false
		},
		/* 45 PausedOpt <- <(<(Paused / Unpaused)?> Action33)> */
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
							if !_rules[rulePaused]() {
								goto l629
							}
							goto l628
						l629:
							position, tokenIndex, depth = position628, tokenIndex628, depth628
							if !_rules[ruleUnpaused]() {
								goto l626
							}
						}
					l628:
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
				add(rulePausedOpt, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 46 Expression <- <orExpr> */
		func() bool {
			position630, tokenIndex630, depth630 := position, tokenIndex, depth
			{
				position631 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l630
				}
				depth--
				add(ruleExpression, position631)
			}
			return true
		l630:
			position, tokenIndex, depth = position630, tokenIndex630, depth630
			return false
		},
		/* 47 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action34)> */
		func() bool {
			position632, tokenIndex632, depth632 := position, tokenIndex, depth
			{
				position633 := position
				depth++
				{
					position634 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l632
					}
					if !_rules[rulesp]() {
						goto l632
					}
					{
						position635, tokenIndex635, depth635 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l635
						}
						if !_rules[rulesp]() {
							goto l635
						}
						if !_rules[ruleandExpr]() {
							goto l635
						}
						goto l636
					l635:
						position, tokenIndex, depth = position635, tokenIndex635, depth635
					}
				l636:
					depth--
					add(rulePegText, position634)
				}
				if !_rules[ruleAction34]() {
					goto l632
				}
				depth--
				add(ruleorExpr, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 48 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action35)> */
		func() bool {
			position637, tokenIndex637, depth637 := position, tokenIndex, depth
			{
				position638 := position
				depth++
				{
					position639 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l637
					}
					if !_rules[rulesp]() {
						goto l637
					}
					{
						position640, tokenIndex640, depth640 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l640
						}
						if !_rules[rulesp]() {
							goto l640
						}
						if !_rules[rulenotExpr]() {
							goto l640
						}
						goto l641
					l640:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
					}
				l641:
					depth--
					add(rulePegText, position639)
				}
				if !_rules[ruleAction35]() {
					goto l637
				}
				depth--
				add(ruleandExpr, position638)
			}
			return true
		l637:
			position, tokenIndex, depth = position637, tokenIndex637, depth637
			return false
		},
		/* 49 notExpr <- <(<((Not sp)? comparisonExpr)> Action36)> */
		func() bool {
			position642, tokenIndex642, depth642 := position, tokenIndex, depth
			{
				position643 := position
				depth++
				{
					position644 := position
					depth++
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l645
						}
						if !_rules[rulesp]() {
							goto l645
						}
						goto l646
					l645:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
					}
				l646:
					if !_rules[rulecomparisonExpr]() {
						goto l642
					}
					depth--
					add(rulePegText, position644)
				}
				if !_rules[ruleAction36]() {
					goto l642
				}
				depth--
				add(rulenotExpr, position643)
			}
			return true
		l642:
			position, tokenIndex, depth = position642, tokenIndex642, depth642
			return false
		},
		/* 50 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action37)> */
		func() bool {
			position647, tokenIndex647, depth647 := position, tokenIndex, depth
			{
				position648 := position
				depth++
				{
					position649 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l647
					}
					if !_rules[rulesp]() {
						goto l647
					}
					{
						position650, tokenIndex650, depth650 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l650
						}
						if !_rules[rulesp]() {
							goto l650
						}
						if !_rules[ruleotherOpExpr]() {
							goto l650
						}
						goto l651
					l650:
						position, tokenIndex, depth = position650, tokenIndex650, depth650
					}
				l651:
					depth--
					add(rulePegText, position649)
				}
				if !_rules[ruleAction37]() {
					goto l647
				}
				depth--
				add(rulecomparisonExpr, position648)
			}
			return true
		l647:
			position, tokenIndex, depth = position647, tokenIndex647, depth647
			return false
		},
		/* 51 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action38)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				{
					position654 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l652
					}
					if !_rules[rulesp]() {
						goto l652
					}
				l655:
					{
						position656, tokenIndex656, depth656 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l656
						}
						if !_rules[rulesp]() {
							goto l656
						}
						if !_rules[ruleisExpr]() {
							goto l656
						}
						if !_rules[rulesp]() {
							goto l656
						}
						goto l655
					l656:
						position, tokenIndex, depth = position656, tokenIndex656, depth656
					}
					depth--
					add(rulePegText, position654)
				}
				if !_rules[ruleAction38]() {
					goto l652
				}
				depth--
				add(ruleotherOpExpr, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 52 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action39)> */
		func() bool {
			position657, tokenIndex657, depth657 := position, tokenIndex, depth
			{
				position658 := position
				depth++
				{
					position659 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l657
					}
					if !_rules[rulesp]() {
						goto l657
					}
					{
						position660, tokenIndex660, depth660 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l660
						}
						if !_rules[rulesp]() {
							goto l660
						}
						if !_rules[ruleNullLiteral]() {
							goto l660
						}
						goto l661
					l660:
						position, tokenIndex, depth = position660, tokenIndex660, depth660
					}
				l661:
					depth--
					add(rulePegText, position659)
				}
				if !_rules[ruleAction39]() {
					goto l657
				}
				depth--
				add(ruleisExpr, position658)
			}
			return true
		l657:
			position, tokenIndex, depth = position657, tokenIndex657, depth657
			return false
		},
		/* 53 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action40)> */
		func() bool {
			position662, tokenIndex662, depth662 := position, tokenIndex, depth
			{
				position663 := position
				depth++
				{
					position664 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l662
					}
					if !_rules[rulesp]() {
						goto l662
					}
				l665:
					{
						position666, tokenIndex666, depth666 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l666
						}
						if !_rules[rulesp]() {
							goto l666
						}
						if !_rules[ruleproductExpr]() {
							goto l666
						}
						if !_rules[rulesp]() {
							goto l666
						}
						goto l665
					l666:
						position, tokenIndex, depth = position666, tokenIndex666, depth666
					}
					depth--
					add(rulePegText, position664)
				}
				if !_rules[ruleAction40]() {
					goto l662
				}
				depth--
				add(ruletermExpr, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 54 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action41)> */
		func() bool {
			position667, tokenIndex667, depth667 := position, tokenIndex, depth
			{
				position668 := position
				depth++
				{
					position669 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l667
					}
					if !_rules[rulesp]() {
						goto l667
					}
				l670:
					{
						position671, tokenIndex671, depth671 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l671
						}
						if !_rules[rulesp]() {
							goto l671
						}
						if !_rules[ruleminusExpr]() {
							goto l671
						}
						if !_rules[rulesp]() {
							goto l671
						}
						goto l670
					l671:
						position, tokenIndex, depth = position671, tokenIndex671, depth671
					}
					depth--
					add(rulePegText, position669)
				}
				if !_rules[ruleAction41]() {
					goto l667
				}
				depth--
				add(ruleproductExpr, position668)
			}
			return true
		l667:
			position, tokenIndex, depth = position667, tokenIndex667, depth667
			return false
		},
		/* 55 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action42)> */
		func() bool {
			position672, tokenIndex672, depth672 := position, tokenIndex, depth
			{
				position673 := position
				depth++
				{
					position674 := position
					depth++
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l675
						}
						if !_rules[rulesp]() {
							goto l675
						}
						goto l676
					l675:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
					}
				l676:
					if !_rules[rulecastExpr]() {
						goto l672
					}
					depth--
					add(rulePegText, position674)
				}
				if !_rules[ruleAction42]() {
					goto l672
				}
				depth--
				add(ruleminusExpr, position673)
			}
			return true
		l672:
			position, tokenIndex, depth = position672, tokenIndex672, depth672
			return false
		},
		/* 56 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action43)> */
		func() bool {
			position677, tokenIndex677, depth677 := position, tokenIndex, depth
			{
				position678 := position
				depth++
				{
					position679 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l677
					}
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l680
						}
						if buffer[position] != rune(':') {
							goto l680
						}
						position++
						if buffer[position] != rune(':') {
							goto l680
						}
						position++
						if !_rules[rulesp]() {
							goto l680
						}
						if !_rules[ruleType]() {
							goto l680
						}
						goto l681
					l680:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
					}
				l681:
					depth--
					add(rulePegText, position679)
				}
				if !_rules[ruleAction43]() {
					goto l677
				}
				depth--
				add(rulecastExpr, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 57 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / Literal)> */
		func() bool {
			position682, tokenIndex682, depth682 := position, tokenIndex, depth
			{
				position683 := position
				depth++
				{
					position684, tokenIndex684, depth684 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l685
					}
					position++
					if !_rules[rulesp]() {
						goto l685
					}
					if !_rules[ruleExpression]() {
						goto l685
					}
					if !_rules[rulesp]() {
						goto l685
					}
					if buffer[position] != rune(')') {
						goto l685
					}
					position++
					goto l684
				l685:
					position, tokenIndex, depth = position684, tokenIndex684, depth684
					if !_rules[ruleBooleanLiteral]() {
						goto l686
					}
					goto l684
				l686:
					position, tokenIndex, depth = position684, tokenIndex684, depth684
					if !_rules[ruleNullLiteral]() {
						goto l687
					}
					goto l684
				l687:
					position, tokenIndex, depth = position684, tokenIndex684, depth684
					if !_rules[ruleRowMeta]() {
						goto l688
					}
					goto l684
				l688:
					position, tokenIndex, depth = position684, tokenIndex684, depth684
					if !_rules[ruleFuncTypeCast]() {
						goto l689
					}
					goto l684
				l689:
					position, tokenIndex, depth = position684, tokenIndex684, depth684
					if !_rules[ruleFuncApp]() {
						goto l690
					}
					goto l684
				l690:
					position, tokenIndex, depth = position684, tokenIndex684, depth684
					if !_rules[ruleRowValue]() {
						goto l691
					}
					goto l684
				l691:
					position, tokenIndex, depth = position684, tokenIndex684, depth684
					if !_rules[ruleLiteral]() {
						goto l682
					}
				}
			l684:
				depth--
				add(rulebaseExpr, position683)
			}
			return true
		l682:
			position, tokenIndex, depth = position682, tokenIndex682, depth682
			return false
		},
		/* 58 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action44)> */
		func() bool {
			position692, tokenIndex692, depth692 := position, tokenIndex, depth
			{
				position693 := position
				depth++
				{
					position694 := position
					depth++
					{
						position695, tokenIndex695, depth695 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
						if buffer[position] != rune('C') {
							goto l692
						}
						position++
					}
				l695:
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('A') {
							goto l692
						}
						position++
					}
				l697:
					{
						position699, tokenIndex699, depth699 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l700
						}
						position++
						goto l699
					l700:
						position, tokenIndex, depth = position699, tokenIndex699, depth699
						if buffer[position] != rune('S') {
							goto l692
						}
						position++
					}
				l699:
					{
						position701, tokenIndex701, depth701 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l702
						}
						position++
						goto l701
					l702:
						position, tokenIndex, depth = position701, tokenIndex701, depth701
						if buffer[position] != rune('T') {
							goto l692
						}
						position++
					}
				l701:
					if !_rules[rulesp]() {
						goto l692
					}
					if buffer[position] != rune('(') {
						goto l692
					}
					position++
					if !_rules[rulesp]() {
						goto l692
					}
					if !_rules[ruleExpression]() {
						goto l692
					}
					if !_rules[rulesp]() {
						goto l692
					}
					{
						position703, tokenIndex703, depth703 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l704
						}
						position++
						goto l703
					l704:
						position, tokenIndex, depth = position703, tokenIndex703, depth703
						if buffer[position] != rune('A') {
							goto l692
						}
						position++
					}
				l703:
					{
						position705, tokenIndex705, depth705 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l706
						}
						position++
						goto l705
					l706:
						position, tokenIndex, depth = position705, tokenIndex705, depth705
						if buffer[position] != rune('S') {
							goto l692
						}
						position++
					}
				l705:
					if !_rules[rulesp]() {
						goto l692
					}
					if !_rules[ruleType]() {
						goto l692
					}
					if !_rules[rulesp]() {
						goto l692
					}
					if buffer[position] != rune(')') {
						goto l692
					}
					position++
					depth--
					add(rulePegText, position694)
				}
				if !_rules[ruleAction44]() {
					goto l692
				}
				depth--
				add(ruleFuncTypeCast, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 59 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action45)> */
		func() bool {
			position707, tokenIndex707, depth707 := position, tokenIndex, depth
			{
				position708 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l707
				}
				if !_rules[rulesp]() {
					goto l707
				}
				if buffer[position] != rune('(') {
					goto l707
				}
				position++
				if !_rules[rulesp]() {
					goto l707
				}
				if !_rules[ruleFuncParams]() {
					goto l707
				}
				if !_rules[rulesp]() {
					goto l707
				}
				if buffer[position] != rune(')') {
					goto l707
				}
				position++
				if !_rules[ruleAction45]() {
					goto l707
				}
				depth--
				add(ruleFuncApp, position708)
			}
			return true
		l707:
			position, tokenIndex, depth = position707, tokenIndex707, depth707
			return false
		},
		/* 60 FuncParams <- <(<(Wildcard / (Expression sp (',' sp Expression)*)?)> Action46)> */
		func() bool {
			position709, tokenIndex709, depth709 := position, tokenIndex, depth
			{
				position710 := position
				depth++
				{
					position711 := position
					depth++
					{
						position712, tokenIndex712, depth712 := position, tokenIndex, depth
						if !_rules[ruleWildcard]() {
							goto l713
						}
						goto l712
					l713:
						position, tokenIndex, depth = position712, tokenIndex712, depth712
						{
							position714, tokenIndex714, depth714 := position, tokenIndex, depth
							if !_rules[ruleExpression]() {
								goto l714
							}
							if !_rules[rulesp]() {
								goto l714
							}
						l716:
							{
								position717, tokenIndex717, depth717 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l717
								}
								position++
								if !_rules[rulesp]() {
									goto l717
								}
								if !_rules[ruleExpression]() {
									goto l717
								}
								goto l716
							l717:
								position, tokenIndex, depth = position717, tokenIndex717, depth717
							}
							goto l715
						l714:
							position, tokenIndex, depth = position714, tokenIndex714, depth714
						}
					l715:
					}
				l712:
					depth--
					add(rulePegText, position711)
				}
				if !_rules[ruleAction46]() {
					goto l709
				}
				depth--
				add(ruleFuncParams, position710)
			}
			return true
		l709:
			position, tokenIndex, depth = position709, tokenIndex709, depth709
			return false
		},
		/* 61 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position718, tokenIndex718, depth718 := position, tokenIndex, depth
			{
				position719 := position
				depth++
				{
					position720, tokenIndex720, depth720 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l721
					}
					goto l720
				l721:
					position, tokenIndex, depth = position720, tokenIndex720, depth720
					if !_rules[ruleNumericLiteral]() {
						goto l722
					}
					goto l720
				l722:
					position, tokenIndex, depth = position720, tokenIndex720, depth720
					if !_rules[ruleStringLiteral]() {
						goto l718
					}
				}
			l720:
				depth--
				add(ruleLiteral, position719)
			}
			return true
		l718:
			position, tokenIndex, depth = position718, tokenIndex718, depth718
			return false
		},
		/* 62 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position723, tokenIndex723, depth723 := position, tokenIndex, depth
			{
				position724 := position
				depth++
				{
					position725, tokenIndex725, depth725 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l726
					}
					goto l725
				l726:
					position, tokenIndex, depth = position725, tokenIndex725, depth725
					if !_rules[ruleNotEqual]() {
						goto l727
					}
					goto l725
				l727:
					position, tokenIndex, depth = position725, tokenIndex725, depth725
					if !_rules[ruleLessOrEqual]() {
						goto l728
					}
					goto l725
				l728:
					position, tokenIndex, depth = position725, tokenIndex725, depth725
					if !_rules[ruleLess]() {
						goto l729
					}
					goto l725
				l729:
					position, tokenIndex, depth = position725, tokenIndex725, depth725
					if !_rules[ruleGreaterOrEqual]() {
						goto l730
					}
					goto l725
				l730:
					position, tokenIndex, depth = position725, tokenIndex725, depth725
					if !_rules[ruleGreater]() {
						goto l731
					}
					goto l725
				l731:
					position, tokenIndex, depth = position725, tokenIndex725, depth725
					if !_rules[ruleNotEqual]() {
						goto l723
					}
				}
			l725:
				depth--
				add(ruleComparisonOp, position724)
			}
			return true
		l723:
			position, tokenIndex, depth = position723, tokenIndex723, depth723
			return false
		},
		/* 63 OtherOp <- <Concat> */
		func() bool {
			position732, tokenIndex732, depth732 := position, tokenIndex, depth
			{
				position733 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l732
				}
				depth--
				add(ruleOtherOp, position733)
			}
			return true
		l732:
			position, tokenIndex, depth = position732, tokenIndex732, depth732
			return false
		},
		/* 64 IsOp <- <(IsNot / Is)> */
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
		/* 65 PlusMinusOp <- <(Plus / Minus)> */
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
		/* 66 MultDivOp <- <(Multiply / Divide / Modulo)> */
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
		/* 67 Stream <- <(<ident> Action47)> */
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
		/* 68 RowMeta <- <RowTimestamp> */
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
		/* 69 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action48)> */
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
		/* 70 RowValue <- <(<((ident ':' !':')? jsonPath)> Action49)> */
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
						{
							position762, tokenIndex762, depth762 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l762
							}
							position++
							goto l760
						l762:
							position, tokenIndex, depth = position762, tokenIndex762, depth762
						}
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
		/* 71 NumericLiteral <- <(<('-'? [0-9]+)> Action50)> */
		func() bool {
			position763, tokenIndex763, depth763 := position, tokenIndex, depth
			{
				position764 := position
				depth++
				{
					position765 := position
					depth++
					{
						position766, tokenIndex766, depth766 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l766
						}
						position++
						goto l767
					l766:
						position, tokenIndex, depth = position766, tokenIndex766, depth766
					}
				l767:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l763
					}
					position++
				l768:
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l769
						}
						position++
						goto l768
					l769:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
					}
					depth--
					add(rulePegText, position765)
				}
				if !_rules[ruleAction50]() {
					goto l763
				}
				depth--
				add(ruleNumericLiteral, position764)
			}
			return true
		l763:
			position, tokenIndex, depth = position763, tokenIndex763, depth763
			return false
		},
		/* 72 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action51)> */
		func() bool {
			position770, tokenIndex770, depth770 := position, tokenIndex, depth
			{
				position771 := position
				depth++
				{
					position772 := position
					depth++
					{
						position773, tokenIndex773, depth773 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l773
						}
						position++
						goto l774
					l773:
						position, tokenIndex, depth = position773, tokenIndex773, depth773
					}
				l774:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l770
					}
					position++
				l775:
					{
						position776, tokenIndex776, depth776 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l776
						}
						position++
						goto l775
					l776:
						position, tokenIndex, depth = position776, tokenIndex776, depth776
					}
					if buffer[position] != rune('.') {
						goto l770
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l770
					}
					position++
				l777:
					{
						position778, tokenIndex778, depth778 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l778
						}
						position++
						goto l777
					l778:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
					}
					depth--
					add(rulePegText, position772)
				}
				if !_rules[ruleAction51]() {
					goto l770
				}
				depth--
				add(ruleFloatLiteral, position771)
			}
			return true
		l770:
			position, tokenIndex, depth = position770, tokenIndex770, depth770
			return false
		},
		/* 73 Function <- <(<ident> Action52)> */
		func() bool {
			position779, tokenIndex779, depth779 := position, tokenIndex, depth
			{
				position780 := position
				depth++
				{
					position781 := position
					depth++
					if !_rules[ruleident]() {
						goto l779
					}
					depth--
					add(rulePegText, position781)
				}
				if !_rules[ruleAction52]() {
					goto l779
				}
				depth--
				add(ruleFunction, position780)
			}
			return true
		l779:
			position, tokenIndex, depth = position779, tokenIndex779, depth779
			return false
		},
		/* 74 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action53)> */
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
						if buffer[position] != rune('n') {
							goto l786
						}
						position++
						goto l785
					l786:
						position, tokenIndex, depth = position785, tokenIndex785, depth785
						if buffer[position] != rune('N') {
							goto l782
						}
						position++
					}
				l785:
					{
						position787, tokenIndex787, depth787 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l788
						}
						position++
						goto l787
					l788:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
						if buffer[position] != rune('U') {
							goto l782
						}
						position++
					}
				l787:
					{
						position789, tokenIndex789, depth789 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l790
						}
						position++
						goto l789
					l790:
						position, tokenIndex, depth = position789, tokenIndex789, depth789
						if buffer[position] != rune('L') {
							goto l782
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
							goto l782
						}
						position++
					}
				l791:
					depth--
					add(rulePegText, position784)
				}
				if !_rules[ruleAction53]() {
					goto l782
				}
				depth--
				add(ruleNullLiteral, position783)
			}
			return true
		l782:
			position, tokenIndex, depth = position782, tokenIndex782, depth782
			return false
		},
		/* 75 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position793, tokenIndex793, depth793 := position, tokenIndex, depth
			{
				position794 := position
				depth++
				{
					position795, tokenIndex795, depth795 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l796
					}
					goto l795
				l796:
					position, tokenIndex, depth = position795, tokenIndex795, depth795
					if !_rules[ruleFALSE]() {
						goto l793
					}
				}
			l795:
				depth--
				add(ruleBooleanLiteral, position794)
			}
			return true
		l793:
			position, tokenIndex, depth = position793, tokenIndex793, depth793
			return false
		},
		/* 76 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action54)> */
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
						if buffer[position] != rune('t') {
							goto l801
						}
						position++
						goto l800
					l801:
						position, tokenIndex, depth = position800, tokenIndex800, depth800
						if buffer[position] != rune('T') {
							goto l797
						}
						position++
					}
				l800:
					{
						position802, tokenIndex802, depth802 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l803
						}
						position++
						goto l802
					l803:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						if buffer[position] != rune('R') {
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
						if buffer[position] != rune('e') {
							goto l807
						}
						position++
						goto l806
					l807:
						position, tokenIndex, depth = position806, tokenIndex806, depth806
						if buffer[position] != rune('E') {
							goto l797
						}
						position++
					}
				l806:
					depth--
					add(rulePegText, position799)
				}
				if !_rules[ruleAction54]() {
					goto l797
				}
				depth--
				add(ruleTRUE, position798)
			}
			return true
		l797:
			position, tokenIndex, depth = position797, tokenIndex797, depth797
			return false
		},
		/* 77 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action55)> */
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
						if buffer[position] != rune('f') {
							goto l812
						}
						position++
						goto l811
					l812:
						position, tokenIndex, depth = position811, tokenIndex811, depth811
						if buffer[position] != rune('F') {
							goto l808
						}
						position++
					}
				l811:
					{
						position813, tokenIndex813, depth813 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l814
						}
						position++
						goto l813
					l814:
						position, tokenIndex, depth = position813, tokenIndex813, depth813
						if buffer[position] != rune('A') {
							goto l808
						}
						position++
					}
				l813:
					{
						position815, tokenIndex815, depth815 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l816
						}
						position++
						goto l815
					l816:
						position, tokenIndex, depth = position815, tokenIndex815, depth815
						if buffer[position] != rune('L') {
							goto l808
						}
						position++
					}
				l815:
					{
						position817, tokenIndex817, depth817 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l818
						}
						position++
						goto l817
					l818:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						if buffer[position] != rune('S') {
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
					depth--
					add(rulePegText, position810)
				}
				if !_rules[ruleAction55]() {
					goto l808
				}
				depth--
				add(ruleFALSE, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 78 Wildcard <- <(<'*'> Action56)> */
		func() bool {
			position821, tokenIndex821, depth821 := position, tokenIndex, depth
			{
				position822 := position
				depth++
				{
					position823 := position
					depth++
					if buffer[position] != rune('*') {
						goto l821
					}
					position++
					depth--
					add(rulePegText, position823)
				}
				if !_rules[ruleAction56]() {
					goto l821
				}
				depth--
				add(ruleWildcard, position822)
			}
			return true
		l821:
			position, tokenIndex, depth = position821, tokenIndex821, depth821
			return false
		},
		/* 79 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action57)> */
		func() bool {
			position824, tokenIndex824, depth824 := position, tokenIndex, depth
			{
				position825 := position
				depth++
				{
					position826 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l824
					}
					position++
				l827:
					{
						position828, tokenIndex828, depth828 := position, tokenIndex, depth
						{
							position829, tokenIndex829, depth829 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l830
							}
							position++
							if buffer[position] != rune('\'') {
								goto l830
							}
							position++
							goto l829
						l830:
							position, tokenIndex, depth = position829, tokenIndex829, depth829
							{
								position831, tokenIndex831, depth831 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l831
								}
								position++
								goto l828
							l831:
								position, tokenIndex, depth = position831, tokenIndex831, depth831
							}
							if !matchDot() {
								goto l828
							}
						}
					l829:
						goto l827
					l828:
						position, tokenIndex, depth = position828, tokenIndex828, depth828
					}
					if buffer[position] != rune('\'') {
						goto l824
					}
					position++
					depth--
					add(rulePegText, position826)
				}
				if !_rules[ruleAction57]() {
					goto l824
				}
				depth--
				add(ruleStringLiteral, position825)
			}
			return true
		l824:
			position, tokenIndex, depth = position824, tokenIndex824, depth824
			return false
		},
		/* 80 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action58)> */
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
					{
						position839, tokenIndex839, depth839 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l840
						}
						position++
						goto l839
					l840:
						position, tokenIndex, depth = position839, tokenIndex839, depth839
						if buffer[position] != rune('T') {
							goto l832
						}
						position++
					}
				l839:
					{
						position841, tokenIndex841, depth841 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l842
						}
						position++
						goto l841
					l842:
						position, tokenIndex, depth = position841, tokenIndex841, depth841
						if buffer[position] != rune('R') {
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
						if buffer[position] != rune('a') {
							goto l846
						}
						position++
						goto l845
					l846:
						position, tokenIndex, depth = position845, tokenIndex845, depth845
						if buffer[position] != rune('A') {
							goto l832
						}
						position++
					}
				l845:
					{
						position847, tokenIndex847, depth847 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l848
						}
						position++
						goto l847
					l848:
						position, tokenIndex, depth = position847, tokenIndex847, depth847
						if buffer[position] != rune('M') {
							goto l832
						}
						position++
					}
				l847:
					depth--
					add(rulePegText, position834)
				}
				if !_rules[ruleAction58]() {
					goto l832
				}
				depth--
				add(ruleISTREAM, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 81 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action59)> */
		func() bool {
			position849, tokenIndex849, depth849 := position, tokenIndex, depth
			{
				position850 := position
				depth++
				{
					position851 := position
					depth++
					{
						position852, tokenIndex852, depth852 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l853
						}
						position++
						goto l852
					l853:
						position, tokenIndex, depth = position852, tokenIndex852, depth852
						if buffer[position] != rune('D') {
							goto l849
						}
						position++
					}
				l852:
					{
						position854, tokenIndex854, depth854 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l855
						}
						position++
						goto l854
					l855:
						position, tokenIndex, depth = position854, tokenIndex854, depth854
						if buffer[position] != rune('S') {
							goto l849
						}
						position++
					}
				l854:
					{
						position856, tokenIndex856, depth856 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l857
						}
						position++
						goto l856
					l857:
						position, tokenIndex, depth = position856, tokenIndex856, depth856
						if buffer[position] != rune('T') {
							goto l849
						}
						position++
					}
				l856:
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l859
						}
						position++
						goto l858
					l859:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if buffer[position] != rune('R') {
							goto l849
						}
						position++
					}
				l858:
					{
						position860, tokenIndex860, depth860 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l861
						}
						position++
						goto l860
					l861:
						position, tokenIndex, depth = position860, tokenIndex860, depth860
						if buffer[position] != rune('E') {
							goto l849
						}
						position++
					}
				l860:
					{
						position862, tokenIndex862, depth862 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l863
						}
						position++
						goto l862
					l863:
						position, tokenIndex, depth = position862, tokenIndex862, depth862
						if buffer[position] != rune('A') {
							goto l849
						}
						position++
					}
				l862:
					{
						position864, tokenIndex864, depth864 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l865
						}
						position++
						goto l864
					l865:
						position, tokenIndex, depth = position864, tokenIndex864, depth864
						if buffer[position] != rune('M') {
							goto l849
						}
						position++
					}
				l864:
					depth--
					add(rulePegText, position851)
				}
				if !_rules[ruleAction59]() {
					goto l849
				}
				depth--
				add(ruleDSTREAM, position850)
			}
			return true
		l849:
			position, tokenIndex, depth = position849, tokenIndex849, depth849
			return false
		},
		/* 82 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action60)> */
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
						if buffer[position] != rune('r') {
							goto l870
						}
						position++
						goto l869
					l870:
						position, tokenIndex, depth = position869, tokenIndex869, depth869
						if buffer[position] != rune('R') {
							goto l866
						}
						position++
					}
				l869:
					{
						position871, tokenIndex871, depth871 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l872
						}
						position++
						goto l871
					l872:
						position, tokenIndex, depth = position871, tokenIndex871, depth871
						if buffer[position] != rune('S') {
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
					{
						position875, tokenIndex875, depth875 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l876
						}
						position++
						goto l875
					l876:
						position, tokenIndex, depth = position875, tokenIndex875, depth875
						if buffer[position] != rune('R') {
							goto l866
						}
						position++
					}
				l875:
					{
						position877, tokenIndex877, depth877 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l878
						}
						position++
						goto l877
					l878:
						position, tokenIndex, depth = position877, tokenIndex877, depth877
						if buffer[position] != rune('E') {
							goto l866
						}
						position++
					}
				l877:
					{
						position879, tokenIndex879, depth879 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l880
						}
						position++
						goto l879
					l880:
						position, tokenIndex, depth = position879, tokenIndex879, depth879
						if buffer[position] != rune('A') {
							goto l866
						}
						position++
					}
				l879:
					{
						position881, tokenIndex881, depth881 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l882
						}
						position++
						goto l881
					l882:
						position, tokenIndex, depth = position881, tokenIndex881, depth881
						if buffer[position] != rune('M') {
							goto l866
						}
						position++
					}
				l881:
					depth--
					add(rulePegText, position868)
				}
				if !_rules[ruleAction60]() {
					goto l866
				}
				depth--
				add(ruleRSTREAM, position867)
			}
			return true
		l866:
			position, tokenIndex, depth = position866, tokenIndex866, depth866
			return false
		},
		/* 83 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action61)> */
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
						if buffer[position] != rune('t') {
							goto l887
						}
						position++
						goto l886
					l887:
						position, tokenIndex, depth = position886, tokenIndex886, depth886
						if buffer[position] != rune('T') {
							goto l883
						}
						position++
					}
				l886:
					{
						position888, tokenIndex888, depth888 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l889
						}
						position++
						goto l888
					l889:
						position, tokenIndex, depth = position888, tokenIndex888, depth888
						if buffer[position] != rune('U') {
							goto l883
						}
						position++
					}
				l888:
					{
						position890, tokenIndex890, depth890 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l891
						}
						position++
						goto l890
					l891:
						position, tokenIndex, depth = position890, tokenIndex890, depth890
						if buffer[position] != rune('P') {
							goto l883
						}
						position++
					}
				l890:
					{
						position892, tokenIndex892, depth892 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l893
						}
						position++
						goto l892
					l893:
						position, tokenIndex, depth = position892, tokenIndex892, depth892
						if buffer[position] != rune('L') {
							goto l883
						}
						position++
					}
				l892:
					{
						position894, tokenIndex894, depth894 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l895
						}
						position++
						goto l894
					l895:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
						if buffer[position] != rune('E') {
							goto l883
						}
						position++
					}
				l894:
					{
						position896, tokenIndex896, depth896 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l897
						}
						position++
						goto l896
					l897:
						position, tokenIndex, depth = position896, tokenIndex896, depth896
						if buffer[position] != rune('S') {
							goto l883
						}
						position++
					}
				l896:
					depth--
					add(rulePegText, position885)
				}
				if !_rules[ruleAction61]() {
					goto l883
				}
				depth--
				add(ruleTUPLES, position884)
			}
			return true
		l883:
			position, tokenIndex, depth = position883, tokenIndex883, depth883
			return false
		},
		/* 84 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action62)> */
		func() bool {
			position898, tokenIndex898, depth898 := position, tokenIndex, depth
			{
				position899 := position
				depth++
				{
					position900 := position
					depth++
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
							goto l898
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
							goto l898
						}
						position++
					}
				l903:
					{
						position905, tokenIndex905, depth905 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l906
						}
						position++
						goto l905
					l906:
						position, tokenIndex, depth = position905, tokenIndex905, depth905
						if buffer[position] != rune('C') {
							goto l898
						}
						position++
					}
				l905:
					{
						position907, tokenIndex907, depth907 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l908
						}
						position++
						goto l907
					l908:
						position, tokenIndex, depth = position907, tokenIndex907, depth907
						if buffer[position] != rune('O') {
							goto l898
						}
						position++
					}
				l907:
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
							goto l898
						}
						position++
					}
				l909:
					{
						position911, tokenIndex911, depth911 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l912
						}
						position++
						goto l911
					l912:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('D') {
							goto l898
						}
						position++
					}
				l911:
					{
						position913, tokenIndex913, depth913 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l914
						}
						position++
						goto l913
					l914:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('S') {
							goto l898
						}
						position++
					}
				l913:
					depth--
					add(rulePegText, position900)
				}
				if !_rules[ruleAction62]() {
					goto l898
				}
				depth--
				add(ruleSECONDS, position899)
			}
			return true
		l898:
			position, tokenIndex, depth = position898, tokenIndex898, depth898
			return false
		},
		/* 85 StreamIdentifier <- <(<ident> Action63)> */
		func() bool {
			position915, tokenIndex915, depth915 := position, tokenIndex, depth
			{
				position916 := position
				depth++
				{
					position917 := position
					depth++
					if !_rules[ruleident]() {
						goto l915
					}
					depth--
					add(rulePegText, position917)
				}
				if !_rules[ruleAction63]() {
					goto l915
				}
				depth--
				add(ruleStreamIdentifier, position916)
			}
			return true
		l915:
			position, tokenIndex, depth = position915, tokenIndex915, depth915
			return false
		},
		/* 86 SourceSinkType <- <(<ident> Action64)> */
		func() bool {
			position918, tokenIndex918, depth918 := position, tokenIndex, depth
			{
				position919 := position
				depth++
				{
					position920 := position
					depth++
					if !_rules[ruleident]() {
						goto l918
					}
					depth--
					add(rulePegText, position920)
				}
				if !_rules[ruleAction64]() {
					goto l918
				}
				depth--
				add(ruleSourceSinkType, position919)
			}
			return true
		l918:
			position, tokenIndex, depth = position918, tokenIndex918, depth918
			return false
		},
		/* 87 SourceSinkParamKey <- <(<ident> Action65)> */
		func() bool {
			position921, tokenIndex921, depth921 := position, tokenIndex, depth
			{
				position922 := position
				depth++
				{
					position923 := position
					depth++
					if !_rules[ruleident]() {
						goto l921
					}
					depth--
					add(rulePegText, position923)
				}
				if !_rules[ruleAction65]() {
					goto l921
				}
				depth--
				add(ruleSourceSinkParamKey, position922)
			}
			return true
		l921:
			position, tokenIndex, depth = position921, tokenIndex921, depth921
			return false
		},
		/* 88 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action66)> */
		func() bool {
			position924, tokenIndex924, depth924 := position, tokenIndex, depth
			{
				position925 := position
				depth++
				{
					position926 := position
					depth++
					{
						position927, tokenIndex927, depth927 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l928
						}
						position++
						goto l927
					l928:
						position, tokenIndex, depth = position927, tokenIndex927, depth927
						if buffer[position] != rune('P') {
							goto l924
						}
						position++
					}
				l927:
					{
						position929, tokenIndex929, depth929 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l930
						}
						position++
						goto l929
					l930:
						position, tokenIndex, depth = position929, tokenIndex929, depth929
						if buffer[position] != rune('A') {
							goto l924
						}
						position++
					}
				l929:
					{
						position931, tokenIndex931, depth931 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l932
						}
						position++
						goto l931
					l932:
						position, tokenIndex, depth = position931, tokenIndex931, depth931
						if buffer[position] != rune('U') {
							goto l924
						}
						position++
					}
				l931:
					{
						position933, tokenIndex933, depth933 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l934
						}
						position++
						goto l933
					l934:
						position, tokenIndex, depth = position933, tokenIndex933, depth933
						if buffer[position] != rune('S') {
							goto l924
						}
						position++
					}
				l933:
					{
						position935, tokenIndex935, depth935 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l936
						}
						position++
						goto l935
					l936:
						position, tokenIndex, depth = position935, tokenIndex935, depth935
						if buffer[position] != rune('E') {
							goto l924
						}
						position++
					}
				l935:
					{
						position937, tokenIndex937, depth937 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l938
						}
						position++
						goto l937
					l938:
						position, tokenIndex, depth = position937, tokenIndex937, depth937
						if buffer[position] != rune('D') {
							goto l924
						}
						position++
					}
				l937:
					depth--
					add(rulePegText, position926)
				}
				if !_rules[ruleAction66]() {
					goto l924
				}
				depth--
				add(rulePaused, position925)
			}
			return true
		l924:
			position, tokenIndex, depth = position924, tokenIndex924, depth924
			return false
		},
		/* 89 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action67)> */
		func() bool {
			position939, tokenIndex939, depth939 := position, tokenIndex, depth
			{
				position940 := position
				depth++
				{
					position941 := position
					depth++
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
							goto l939
						}
						position++
					}
				l942:
					{
						position944, tokenIndex944, depth944 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l945
						}
						position++
						goto l944
					l945:
						position, tokenIndex, depth = position944, tokenIndex944, depth944
						if buffer[position] != rune('N') {
							goto l939
						}
						position++
					}
				l944:
					{
						position946, tokenIndex946, depth946 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l947
						}
						position++
						goto l946
					l947:
						position, tokenIndex, depth = position946, tokenIndex946, depth946
						if buffer[position] != rune('P') {
							goto l939
						}
						position++
					}
				l946:
					{
						position948, tokenIndex948, depth948 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l949
						}
						position++
						goto l948
					l949:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
						if buffer[position] != rune('A') {
							goto l939
						}
						position++
					}
				l948:
					{
						position950, tokenIndex950, depth950 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l951
						}
						position++
						goto l950
					l951:
						position, tokenIndex, depth = position950, tokenIndex950, depth950
						if buffer[position] != rune('U') {
							goto l939
						}
						position++
					}
				l950:
					{
						position952, tokenIndex952, depth952 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l953
						}
						position++
						goto l952
					l953:
						position, tokenIndex, depth = position952, tokenIndex952, depth952
						if buffer[position] != rune('S') {
							goto l939
						}
						position++
					}
				l952:
					{
						position954, tokenIndex954, depth954 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l955
						}
						position++
						goto l954
					l955:
						position, tokenIndex, depth = position954, tokenIndex954, depth954
						if buffer[position] != rune('E') {
							goto l939
						}
						position++
					}
				l954:
					{
						position956, tokenIndex956, depth956 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l957
						}
						position++
						goto l956
					l957:
						position, tokenIndex, depth = position956, tokenIndex956, depth956
						if buffer[position] != rune('D') {
							goto l939
						}
						position++
					}
				l956:
					depth--
					add(rulePegText, position941)
				}
				if !_rules[ruleAction67]() {
					goto l939
				}
				depth--
				add(ruleUnpaused, position940)
			}
			return true
		l939:
			position, tokenIndex, depth = position939, tokenIndex939, depth939
			return false
		},
		/* 90 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position958, tokenIndex958, depth958 := position, tokenIndex, depth
			{
				position959 := position
				depth++
				{
					position960, tokenIndex960, depth960 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l961
					}
					goto l960
				l961:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if !_rules[ruleInt]() {
						goto l962
					}
					goto l960
				l962:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if !_rules[ruleFloat]() {
						goto l963
					}
					goto l960
				l963:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if !_rules[ruleString]() {
						goto l964
					}
					goto l960
				l964:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if !_rules[ruleBlob]() {
						goto l965
					}
					goto l960
				l965:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if !_rules[ruleTimestamp]() {
						goto l966
					}
					goto l960
				l966:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if !_rules[ruleArray]() {
						goto l967
					}
					goto l960
				l967:
					position, tokenIndex, depth = position960, tokenIndex960, depth960
					if !_rules[ruleMap]() {
						goto l958
					}
				}
			l960:
				depth--
				add(ruleType, position959)
			}
			return true
		l958:
			position, tokenIndex, depth = position958, tokenIndex958, depth958
			return false
		},
		/* 91 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action68)> */
		func() bool {
			position968, tokenIndex968, depth968 := position, tokenIndex, depth
			{
				position969 := position
				depth++
				{
					position970 := position
					depth++
					{
						position971, tokenIndex971, depth971 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l972
						}
						position++
						goto l971
					l972:
						position, tokenIndex, depth = position971, tokenIndex971, depth971
						if buffer[position] != rune('B') {
							goto l968
						}
						position++
					}
				l971:
					{
						position973, tokenIndex973, depth973 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l974
						}
						position++
						goto l973
					l974:
						position, tokenIndex, depth = position973, tokenIndex973, depth973
						if buffer[position] != rune('O') {
							goto l968
						}
						position++
					}
				l973:
					{
						position975, tokenIndex975, depth975 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l976
						}
						position++
						goto l975
					l976:
						position, tokenIndex, depth = position975, tokenIndex975, depth975
						if buffer[position] != rune('O') {
							goto l968
						}
						position++
					}
				l975:
					{
						position977, tokenIndex977, depth977 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l978
						}
						position++
						goto l977
					l978:
						position, tokenIndex, depth = position977, tokenIndex977, depth977
						if buffer[position] != rune('L') {
							goto l968
						}
						position++
					}
				l977:
					depth--
					add(rulePegText, position970)
				}
				if !_rules[ruleAction68]() {
					goto l968
				}
				depth--
				add(ruleBool, position969)
			}
			return true
		l968:
			position, tokenIndex, depth = position968, tokenIndex968, depth968
			return false
		},
		/* 92 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action69)> */
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
						if buffer[position] != rune('n') {
							goto l985
						}
						position++
						goto l984
					l985:
						position, tokenIndex, depth = position984, tokenIndex984, depth984
						if buffer[position] != rune('N') {
							goto l979
						}
						position++
					}
				l984:
					{
						position986, tokenIndex986, depth986 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l987
						}
						position++
						goto l986
					l987:
						position, tokenIndex, depth = position986, tokenIndex986, depth986
						if buffer[position] != rune('T') {
							goto l979
						}
						position++
					}
				l986:
					depth--
					add(rulePegText, position981)
				}
				if !_rules[ruleAction69]() {
					goto l979
				}
				depth--
				add(ruleInt, position980)
			}
			return true
		l979:
			position, tokenIndex, depth = position979, tokenIndex979, depth979
			return false
		},
		/* 93 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action70)> */
		func() bool {
			position988, tokenIndex988, depth988 := position, tokenIndex, depth
			{
				position989 := position
				depth++
				{
					position990 := position
					depth++
					{
						position991, tokenIndex991, depth991 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l992
						}
						position++
						goto l991
					l992:
						position, tokenIndex, depth = position991, tokenIndex991, depth991
						if buffer[position] != rune('F') {
							goto l988
						}
						position++
					}
				l991:
					{
						position993, tokenIndex993, depth993 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l994
						}
						position++
						goto l993
					l994:
						position, tokenIndex, depth = position993, tokenIndex993, depth993
						if buffer[position] != rune('L') {
							goto l988
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
							goto l988
						}
						position++
					}
				l995:
					{
						position997, tokenIndex997, depth997 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l998
						}
						position++
						goto l997
					l998:
						position, tokenIndex, depth = position997, tokenIndex997, depth997
						if buffer[position] != rune('A') {
							goto l988
						}
						position++
					}
				l997:
					{
						position999, tokenIndex999, depth999 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1000
						}
						position++
						goto l999
					l1000:
						position, tokenIndex, depth = position999, tokenIndex999, depth999
						if buffer[position] != rune('T') {
							goto l988
						}
						position++
					}
				l999:
					depth--
					add(rulePegText, position990)
				}
				if !_rules[ruleAction70]() {
					goto l988
				}
				depth--
				add(ruleFloat, position989)
			}
			return true
		l988:
			position, tokenIndex, depth = position988, tokenIndex988, depth988
			return false
		},
		/* 94 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action71)> */
		func() bool {
			position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
			{
				position1002 := position
				depth++
				{
					position1003 := position
					depth++
					{
						position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1005
						}
						position++
						goto l1004
					l1005:
						position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
						if buffer[position] != rune('S') {
							goto l1001
						}
						position++
					}
				l1004:
					{
						position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1007
						}
						position++
						goto l1006
					l1007:
						position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
						if buffer[position] != rune('T') {
							goto l1001
						}
						position++
					}
				l1006:
					{
						position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1009
						}
						position++
						goto l1008
					l1009:
						position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
						if buffer[position] != rune('R') {
							goto l1001
						}
						position++
					}
				l1008:
					{
						position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1011
						}
						position++
						goto l1010
					l1011:
						position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
						if buffer[position] != rune('I') {
							goto l1001
						}
						position++
					}
				l1010:
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('N') {
							goto l1001
						}
						position++
					}
				l1012:
					{
						position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1015
						}
						position++
						goto l1014
					l1015:
						position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
						if buffer[position] != rune('G') {
							goto l1001
						}
						position++
					}
				l1014:
					depth--
					add(rulePegText, position1003)
				}
				if !_rules[ruleAction71]() {
					goto l1001
				}
				depth--
				add(ruleString, position1002)
			}
			return true
		l1001:
			position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
			return false
		},
		/* 95 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action72)> */
		func() bool {
			position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
			{
				position1017 := position
				depth++
				{
					position1018 := position
					depth++
					{
						position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1020
						}
						position++
						goto l1019
					l1020:
						position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
						if buffer[position] != rune('B') {
							goto l1016
						}
						position++
					}
				l1019:
					{
						position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1022
						}
						position++
						goto l1021
					l1022:
						position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
						if buffer[position] != rune('L') {
							goto l1016
						}
						position++
					}
				l1021:
					{
						position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1024
						}
						position++
						goto l1023
					l1024:
						position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
						if buffer[position] != rune('O') {
							goto l1016
						}
						position++
					}
				l1023:
					{
						position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1026
						}
						position++
						goto l1025
					l1026:
						position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
						if buffer[position] != rune('B') {
							goto l1016
						}
						position++
					}
				l1025:
					depth--
					add(rulePegText, position1018)
				}
				if !_rules[ruleAction72]() {
					goto l1016
				}
				depth--
				add(ruleBlob, position1017)
			}
			return true
		l1016:
			position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
			return false
		},
		/* 96 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action73)> */
		func() bool {
			position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
			{
				position1028 := position
				depth++
				{
					position1029 := position
					depth++
					{
						position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1031
						}
						position++
						goto l1030
					l1031:
						position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
						if buffer[position] != rune('T') {
							goto l1027
						}
						position++
					}
				l1030:
					{
						position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1033
						}
						position++
						goto l1032
					l1033:
						position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
						if buffer[position] != rune('I') {
							goto l1027
						}
						position++
					}
				l1032:
					{
						position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1035
						}
						position++
						goto l1034
					l1035:
						position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
						if buffer[position] != rune('M') {
							goto l1027
						}
						position++
					}
				l1034:
					{
						position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1037
						}
						position++
						goto l1036
					l1037:
						position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
						if buffer[position] != rune('E') {
							goto l1027
						}
						position++
					}
				l1036:
					{
						position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1039
						}
						position++
						goto l1038
					l1039:
						position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
						if buffer[position] != rune('S') {
							goto l1027
						}
						position++
					}
				l1038:
					{
						position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1041
						}
						position++
						goto l1040
					l1041:
						position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
						if buffer[position] != rune('T') {
							goto l1027
						}
						position++
					}
				l1040:
					{
						position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1043
						}
						position++
						goto l1042
					l1043:
						position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
						if buffer[position] != rune('A') {
							goto l1027
						}
						position++
					}
				l1042:
					{
						position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1045
						}
						position++
						goto l1044
					l1045:
						position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
						if buffer[position] != rune('M') {
							goto l1027
						}
						position++
					}
				l1044:
					{
						position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1047
						}
						position++
						goto l1046
					l1047:
						position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
						if buffer[position] != rune('P') {
							goto l1027
						}
						position++
					}
				l1046:
					depth--
					add(rulePegText, position1029)
				}
				if !_rules[ruleAction73]() {
					goto l1027
				}
				depth--
				add(ruleTimestamp, position1028)
			}
			return true
		l1027:
			position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
			return false
		},
		/* 97 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action74)> */
		func() bool {
			position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
			{
				position1049 := position
				depth++
				{
					position1050 := position
					depth++
					{
						position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1052
						}
						position++
						goto l1051
					l1052:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if buffer[position] != rune('A') {
							goto l1048
						}
						position++
					}
				l1051:
					{
						position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1054
						}
						position++
						goto l1053
					l1054:
						position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
						if buffer[position] != rune('R') {
							goto l1048
						}
						position++
					}
				l1053:
					{
						position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1056
						}
						position++
						goto l1055
					l1056:
						position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
						if buffer[position] != rune('R') {
							goto l1048
						}
						position++
					}
				l1055:
					{
						position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1058
						}
						position++
						goto l1057
					l1058:
						position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
						if buffer[position] != rune('A') {
							goto l1048
						}
						position++
					}
				l1057:
					{
						position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1060
						}
						position++
						goto l1059
					l1060:
						position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
						if buffer[position] != rune('Y') {
							goto l1048
						}
						position++
					}
				l1059:
					depth--
					add(rulePegText, position1050)
				}
				if !_rules[ruleAction74]() {
					goto l1048
				}
				depth--
				add(ruleArray, position1049)
			}
			return true
		l1048:
			position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
			return false
		},
		/* 98 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action75)> */
		func() bool {
			position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
			{
				position1062 := position
				depth++
				{
					position1063 := position
					depth++
					{
						position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1065
						}
						position++
						goto l1064
					l1065:
						position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
						if buffer[position] != rune('M') {
							goto l1061
						}
						position++
					}
				l1064:
					{
						position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1067
						}
						position++
						goto l1066
					l1067:
						position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
						if buffer[position] != rune('A') {
							goto l1061
						}
						position++
					}
				l1066:
					{
						position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1069
						}
						position++
						goto l1068
					l1069:
						position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
						if buffer[position] != rune('P') {
							goto l1061
						}
						position++
					}
				l1068:
					depth--
					add(rulePegText, position1063)
				}
				if !_rules[ruleAction75]() {
					goto l1061
				}
				depth--
				add(ruleMap, position1062)
			}
			return true
		l1061:
			position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
			return false
		},
		/* 99 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action76)> */
		func() bool {
			position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
			{
				position1071 := position
				depth++
				{
					position1072 := position
					depth++
					{
						position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1074
						}
						position++
						goto l1073
					l1074:
						position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
						if buffer[position] != rune('O') {
							goto l1070
						}
						position++
					}
				l1073:
					{
						position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1076
						}
						position++
						goto l1075
					l1076:
						position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
						if buffer[position] != rune('R') {
							goto l1070
						}
						position++
					}
				l1075:
					depth--
					add(rulePegText, position1072)
				}
				if !_rules[ruleAction76]() {
					goto l1070
				}
				depth--
				add(ruleOr, position1071)
			}
			return true
		l1070:
			position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
			return false
		},
		/* 100 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action77)> */
		func() bool {
			position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
			{
				position1078 := position
				depth++
				{
					position1079 := position
					depth++
					{
						position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1081
						}
						position++
						goto l1080
					l1081:
						position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
						if buffer[position] != rune('A') {
							goto l1077
						}
						position++
					}
				l1080:
					{
						position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1083
						}
						position++
						goto l1082
					l1083:
						position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
						if buffer[position] != rune('N') {
							goto l1077
						}
						position++
					}
				l1082:
					{
						position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1085
						}
						position++
						goto l1084
					l1085:
						position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
						if buffer[position] != rune('D') {
							goto l1077
						}
						position++
					}
				l1084:
					depth--
					add(rulePegText, position1079)
				}
				if !_rules[ruleAction77]() {
					goto l1077
				}
				depth--
				add(ruleAnd, position1078)
			}
			return true
		l1077:
			position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
			return false
		},
		/* 101 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action78)> */
		func() bool {
			position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
			{
				position1087 := position
				depth++
				{
					position1088 := position
					depth++
					{
						position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1090
						}
						position++
						goto l1089
					l1090:
						position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
						if buffer[position] != rune('N') {
							goto l1086
						}
						position++
					}
				l1089:
					{
						position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1092
						}
						position++
						goto l1091
					l1092:
						position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
						if buffer[position] != rune('O') {
							goto l1086
						}
						position++
					}
				l1091:
					{
						position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1094
						}
						position++
						goto l1093
					l1094:
						position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
						if buffer[position] != rune('T') {
							goto l1086
						}
						position++
					}
				l1093:
					depth--
					add(rulePegText, position1088)
				}
				if !_rules[ruleAction78]() {
					goto l1086
				}
				depth--
				add(ruleNot, position1087)
			}
			return true
		l1086:
			position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
			return false
		},
		/* 102 Equal <- <(<'='> Action79)> */
		func() bool {
			position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
			{
				position1096 := position
				depth++
				{
					position1097 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1095
					}
					position++
					depth--
					add(rulePegText, position1097)
				}
				if !_rules[ruleAction79]() {
					goto l1095
				}
				depth--
				add(ruleEqual, position1096)
			}
			return true
		l1095:
			position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
			return false
		},
		/* 103 Less <- <(<'<'> Action80)> */
		func() bool {
			position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
			{
				position1099 := position
				depth++
				{
					position1100 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1098
					}
					position++
					depth--
					add(rulePegText, position1100)
				}
				if !_rules[ruleAction80]() {
					goto l1098
				}
				depth--
				add(ruleLess, position1099)
			}
			return true
		l1098:
			position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
			return false
		},
		/* 104 LessOrEqual <- <(<('<' '=')> Action81)> */
		func() bool {
			position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
			{
				position1102 := position
				depth++
				{
					position1103 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1101
					}
					position++
					if buffer[position] != rune('=') {
						goto l1101
					}
					position++
					depth--
					add(rulePegText, position1103)
				}
				if !_rules[ruleAction81]() {
					goto l1101
				}
				depth--
				add(ruleLessOrEqual, position1102)
			}
			return true
		l1101:
			position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
			return false
		},
		/* 105 Greater <- <(<'>'> Action82)> */
		func() bool {
			position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
			{
				position1105 := position
				depth++
				{
					position1106 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1104
					}
					position++
					depth--
					add(rulePegText, position1106)
				}
				if !_rules[ruleAction82]() {
					goto l1104
				}
				depth--
				add(ruleGreater, position1105)
			}
			return true
		l1104:
			position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
			return false
		},
		/* 106 GreaterOrEqual <- <(<('>' '=')> Action83)> */
		func() bool {
			position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
			{
				position1108 := position
				depth++
				{
					position1109 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1107
					}
					position++
					if buffer[position] != rune('=') {
						goto l1107
					}
					position++
					depth--
					add(rulePegText, position1109)
				}
				if !_rules[ruleAction83]() {
					goto l1107
				}
				depth--
				add(ruleGreaterOrEqual, position1108)
			}
			return true
		l1107:
			position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
			return false
		},
		/* 107 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action84)> */
		func() bool {
			position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
			{
				position1111 := position
				depth++
				{
					position1112 := position
					depth++
					{
						position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1114
						}
						position++
						if buffer[position] != rune('=') {
							goto l1114
						}
						position++
						goto l1113
					l1114:
						position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
						if buffer[position] != rune('<') {
							goto l1110
						}
						position++
						if buffer[position] != rune('>') {
							goto l1110
						}
						position++
					}
				l1113:
					depth--
					add(rulePegText, position1112)
				}
				if !_rules[ruleAction84]() {
					goto l1110
				}
				depth--
				add(ruleNotEqual, position1111)
			}
			return true
		l1110:
			position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
			return false
		},
		/* 108 Concat <- <(<('|' '|')> Action85)> */
		func() bool {
			position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
			{
				position1116 := position
				depth++
				{
					position1117 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1115
					}
					position++
					if buffer[position] != rune('|') {
						goto l1115
					}
					position++
					depth--
					add(rulePegText, position1117)
				}
				if !_rules[ruleAction85]() {
					goto l1115
				}
				depth--
				add(ruleConcat, position1116)
			}
			return true
		l1115:
			position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
			return false
		},
		/* 109 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action86)> */
		func() bool {
			position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
			{
				position1119 := position
				depth++
				{
					position1120 := position
					depth++
					{
						position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1122
						}
						position++
						goto l1121
					l1122:
						position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
						if buffer[position] != rune('I') {
							goto l1118
						}
						position++
					}
				l1121:
					{
						position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1124
						}
						position++
						goto l1123
					l1124:
						position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
						if buffer[position] != rune('S') {
							goto l1118
						}
						position++
					}
				l1123:
					depth--
					add(rulePegText, position1120)
				}
				if !_rules[ruleAction86]() {
					goto l1118
				}
				depth--
				add(ruleIs, position1119)
			}
			return true
		l1118:
			position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
			return false
		},
		/* 110 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action87)> */
		func() bool {
			position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
			{
				position1126 := position
				depth++
				{
					position1127 := position
					depth++
					{
						position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1129
						}
						position++
						goto l1128
					l1129:
						position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
						if buffer[position] != rune('I') {
							goto l1125
						}
						position++
					}
				l1128:
					{
						position1130, tokenIndex1130, depth1130 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1131
						}
						position++
						goto l1130
					l1131:
						position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
						if buffer[position] != rune('S') {
							goto l1125
						}
						position++
					}
				l1130:
					if !_rules[rulesp]() {
						goto l1125
					}
					{
						position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1133
						}
						position++
						goto l1132
					l1133:
						position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
						if buffer[position] != rune('N') {
							goto l1125
						}
						position++
					}
				l1132:
					{
						position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1135
						}
						position++
						goto l1134
					l1135:
						position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
						if buffer[position] != rune('O') {
							goto l1125
						}
						position++
					}
				l1134:
					{
						position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1137
						}
						position++
						goto l1136
					l1137:
						position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
						if buffer[position] != rune('T') {
							goto l1125
						}
						position++
					}
				l1136:
					depth--
					add(rulePegText, position1127)
				}
				if !_rules[ruleAction87]() {
					goto l1125
				}
				depth--
				add(ruleIsNot, position1126)
			}
			return true
		l1125:
			position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
			return false
		},
		/* 111 Plus <- <(<'+'> Action88)> */
		func() bool {
			position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
			{
				position1139 := position
				depth++
				{
					position1140 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1138
					}
					position++
					depth--
					add(rulePegText, position1140)
				}
				if !_rules[ruleAction88]() {
					goto l1138
				}
				depth--
				add(rulePlus, position1139)
			}
			return true
		l1138:
			position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
			return false
		},
		/* 112 Minus <- <(<'-'> Action89)> */
		func() bool {
			position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
			{
				position1142 := position
				depth++
				{
					position1143 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1141
					}
					position++
					depth--
					add(rulePegText, position1143)
				}
				if !_rules[ruleAction89]() {
					goto l1141
				}
				depth--
				add(ruleMinus, position1142)
			}
			return true
		l1141:
			position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
			return false
		},
		/* 113 Multiply <- <(<'*'> Action90)> */
		func() bool {
			position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
			{
				position1145 := position
				depth++
				{
					position1146 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1144
					}
					position++
					depth--
					add(rulePegText, position1146)
				}
				if !_rules[ruleAction90]() {
					goto l1144
				}
				depth--
				add(ruleMultiply, position1145)
			}
			return true
		l1144:
			position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
			return false
		},
		/* 114 Divide <- <(<'/'> Action91)> */
		func() bool {
			position1147, tokenIndex1147, depth1147 := position, tokenIndex, depth
			{
				position1148 := position
				depth++
				{
					position1149 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1147
					}
					position++
					depth--
					add(rulePegText, position1149)
				}
				if !_rules[ruleAction91]() {
					goto l1147
				}
				depth--
				add(ruleDivide, position1148)
			}
			return true
		l1147:
			position, tokenIndex, depth = position1147, tokenIndex1147, depth1147
			return false
		},
		/* 115 Modulo <- <(<'%'> Action92)> */
		func() bool {
			position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
			{
				position1151 := position
				depth++
				{
					position1152 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1150
					}
					position++
					depth--
					add(rulePegText, position1152)
				}
				if !_rules[ruleAction92]() {
					goto l1150
				}
				depth--
				add(ruleModulo, position1151)
			}
			return true
		l1150:
			position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
			return false
		},
		/* 116 UnaryMinus <- <(<'-'> Action93)> */
		func() bool {
			position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
			{
				position1154 := position
				depth++
				{
					position1155 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1153
					}
					position++
					depth--
					add(rulePegText, position1155)
				}
				if !_rules[ruleAction93]() {
					goto l1153
				}
				depth--
				add(ruleUnaryMinus, position1154)
			}
			return true
		l1153:
			position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
			return false
		},
		/* 117 Identifier <- <(<ident> Action94)> */
		func() bool {
			position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
			{
				position1157 := position
				depth++
				{
					position1158 := position
					depth++
					if !_rules[ruleident]() {
						goto l1156
					}
					depth--
					add(rulePegText, position1158)
				}
				if !_rules[ruleAction94]() {
					goto l1156
				}
				depth--
				add(ruleIdentifier, position1157)
			}
			return true
		l1156:
			position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
			return false
		},
		/* 118 TargetIdentifier <- <(<jsonPath> Action95)> */
		func() bool {
			position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
			{
				position1160 := position
				depth++
				{
					position1161 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1159
					}
					depth--
					add(rulePegText, position1161)
				}
				if !_rules[ruleAction95]() {
					goto l1159
				}
				depth--
				add(ruleTargetIdentifier, position1160)
			}
			return true
		l1159:
			position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
			return false
		},
		/* 119 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
			{
				position1163 := position
				depth++
				{
					position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1165
					}
					position++
					goto l1164
				l1165:
					position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1162
					}
					position++
				}
			l1164:
			l1166:
				{
					position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
					{
						position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1169
						}
						position++
						goto l1168
					l1169:
						position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1170
						}
						position++
						goto l1168
					l1170:
						position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1171
						}
						position++
						goto l1168
					l1171:
						position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
						if buffer[position] != rune('_') {
							goto l1167
						}
						position++
					}
				l1168:
					goto l1166
				l1167:
					position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
				}
				depth--
				add(ruleident, position1163)
			}
			return true
		l1162:
			position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
			return false
		},
		/* 120 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
			{
				position1173 := position
				depth++
				{
					position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1175
					}
					position++
					goto l1174
				l1175:
					position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1172
					}
					position++
				}
			l1174:
			l1176:
				{
					position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
					{
						position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1179
						}
						position++
						goto l1178
					l1179:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1180
						}
						position++
						goto l1178
					l1180:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1181
						}
						position++
						goto l1178
					l1181:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('_') {
							goto l1182
						}
						position++
						goto l1178
					l1182:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('.') {
							goto l1183
						}
						position++
						goto l1178
					l1183:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('[') {
							goto l1184
						}
						position++
						goto l1178
					l1184:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune(']') {
							goto l1185
						}
						position++
						goto l1178
					l1185:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('"') {
							goto l1177
						}
						position++
					}
				l1178:
					goto l1176
				l1177:
					position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
				}
				depth--
				add(rulejsonPath, position1173)
			}
			return true
		l1172:
			position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
			return false
		},
		/* 121 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1187 := position
				depth++
			l1188:
				{
					position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
					{
						position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1191
						}
						position++
						goto l1190
					l1191:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('\t') {
							goto l1192
						}
						position++
						goto l1190
					l1192:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('\n') {
							goto l1193
						}
						position++
						goto l1190
					l1193:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('\r') {
							goto l1194
						}
						position++
						goto l1190
					l1194:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if !_rules[rulecomment]() {
							goto l1189
						}
					}
				l1190:
					goto l1188
				l1189:
					position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
				}
				depth--
				add(rulesp, position1187)
			}
			return true
		},
		/* 122 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
			{
				position1196 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1195
				}
				position++
				if buffer[position] != rune('-') {
					goto l1195
				}
				position++
			l1197:
				{
					position1198, tokenIndex1198, depth1198 := position, tokenIndex, depth
					{
						position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
						{
							position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1201
							}
							position++
							goto l1200
						l1201:
							position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
							if buffer[position] != rune('\n') {
								goto l1199
							}
							position++
						}
					l1200:
						goto l1198
					l1199:
						position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
					}
					if !matchDot() {
						goto l1198
					}
					goto l1197
				l1198:
					position, tokenIndex, depth = position1198, tokenIndex1198, depth1198
				}
				{
					position1202, tokenIndex1202, depth1202 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1203
					}
					position++
					goto l1202
				l1203:
					position, tokenIndex, depth = position1202, tokenIndex1202, depth1202
					if buffer[position] != rune('\n') {
						goto l1195
					}
					position++
				}
			l1202:
				depth--
				add(rulecomment, position1196)
			}
			return true
		l1195:
			position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
			return false
		},
		/* 124 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 125 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 126 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 127 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 128 Action4 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 129 Action5 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 130 Action6 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 131 Action7 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 132 Action8 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 133 Action9 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 134 Action10 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 135 Action11 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 136 Action12 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 137 Action13 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 138 Action14 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 139 Action15 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 140 Action16 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 141 Action17 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		nil,
		/* 143 Action18 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 144 Action19 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 145 Action20 <- <{
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
		/* 146 Action21 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 147 Action22 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 148 Action23 <- <{
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
		/* 149 Action24 <- <{
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
		/* 150 Action25 <- <{
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
		/* 151 Action26 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 152 Action27 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 153 Action28 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 154 Action29 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 155 Action30 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 156 Action31 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 157 Action32 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 158 Action33 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 159 Action34 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 160 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 161 Action36 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 162 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 163 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 164 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 165 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 166 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 167 Action42 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 168 Action43 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 169 Action44 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 170 Action45 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 171 Action46 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 172 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 173 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 174 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 175 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 176 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 177 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 178 Action53 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 179 Action54 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 180 Action55 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 181 Action56 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 182 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 183 Action58 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 184 Action59 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 185 Action60 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 186 Action61 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 187 Action62 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 188 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 189 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 190 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 191 Action66 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 192 Action67 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 193 Action68 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 194 Action69 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 195 Action70 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 196 Action71 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 197 Action72 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 198 Action73 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 199 Action74 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 200 Action75 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 201 Action76 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 202 Action77 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 203 Action78 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 204 Action79 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 205 Action80 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 206 Action81 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 207 Action82 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 208 Action83 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 209 Action84 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 210 Action85 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 211 Action86 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 212 Action87 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 213 Action88 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 214 Action89 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 215 Action90 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 216 Action91 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 217 Action92 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 218 Action93 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 219 Action94 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 220 Action95 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
	}
	p.rules = _rules
}

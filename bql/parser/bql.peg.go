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
	ruleSelectUnionStmt
	ruleCreateStreamAsSelectStmt
	ruleCreateStreamAsSelectUnionStmt
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
	rulePegText
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
	ruleAction96
	ruleAction97

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
	"SelectUnionStmt",
	"CreateStreamAsSelectStmt",
	"CreateStreamAsSelectUnionStmt",
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
	"PegText",
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
	"Action96",
	"Action97",

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
	rules  [225]func() bool
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

			p.AssembleSelectUnion(begin, end)

		case ruleAction2:

			p.AssembleCreateStreamAsSelect()

		case ruleAction3:

			p.AssembleCreateStreamAsSelectUnion()

		case ruleAction4:

			p.AssembleCreateSource()

		case ruleAction5:

			p.AssembleCreateSink()

		case ruleAction6:

			p.AssembleCreateState()

		case ruleAction7:

			p.AssembleUpdateState()

		case ruleAction8:

			p.AssembleUpdateSource()

		case ruleAction9:

			p.AssembleUpdateSink()

		case ruleAction10:

			p.AssembleInsertIntoSelect()

		case ruleAction11:

			p.AssembleInsertIntoFrom()

		case ruleAction12:

			p.AssemblePauseSource()

		case ruleAction13:

			p.AssembleResumeSource()

		case ruleAction14:

			p.AssembleRewindSource()

		case ruleAction15:

			p.AssembleDropSource()

		case ruleAction16:

			p.AssembleDropStream()

		case ruleAction17:

			p.AssembleDropSink()

		case ruleAction18:

			p.AssembleDropState()

		case ruleAction19:

			p.AssembleEmitter()

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

			p.AssembleBinaryOperation(begin, end)

		case ruleAction44:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction45:

			p.AssembleTypeCast(begin, end)

		case ruleAction46:

			p.AssembleTypeCast(begin, end)

		case ruleAction47:

			p.AssembleFuncApp()

		case ruleAction48:

			p.AssembleExpressions(begin, end)

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction53:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction54:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction55:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction56:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction57:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction58:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction59:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction60:

			p.PushComponent(begin, end, Istream)

		case ruleAction61:

			p.PushComponent(begin, end, Dstream)

		case ruleAction62:

			p.PushComponent(begin, end, Rstream)

		case ruleAction63:

			p.PushComponent(begin, end, Tuples)

		case ruleAction64:

			p.PushComponent(begin, end, Seconds)

		case ruleAction65:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction66:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction67:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction68:

			p.PushComponent(begin, end, Yes)

		case ruleAction69:

			p.PushComponent(begin, end, No)

		case ruleAction70:

			p.PushComponent(begin, end, Bool)

		case ruleAction71:

			p.PushComponent(begin, end, Int)

		case ruleAction72:

			p.PushComponent(begin, end, Float)

		case ruleAction73:

			p.PushComponent(begin, end, String)

		case ruleAction74:

			p.PushComponent(begin, end, Blob)

		case ruleAction75:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction76:

			p.PushComponent(begin, end, Array)

		case ruleAction77:

			p.PushComponent(begin, end, Map)

		case ruleAction78:

			p.PushComponent(begin, end, Or)

		case ruleAction79:

			p.PushComponent(begin, end, And)

		case ruleAction80:

			p.PushComponent(begin, end, Not)

		case ruleAction81:

			p.PushComponent(begin, end, Equal)

		case ruleAction82:

			p.PushComponent(begin, end, Less)

		case ruleAction83:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction84:

			p.PushComponent(begin, end, Greater)

		case ruleAction85:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction86:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction87:

			p.PushComponent(begin, end, Concat)

		case ruleAction88:

			p.PushComponent(begin, end, Is)

		case ruleAction89:

			p.PushComponent(begin, end, IsNot)

		case ruleAction90:

			p.PushComponent(begin, end, Plus)

		case ruleAction91:

			p.PushComponent(begin, end, Minus)

		case ruleAction92:

			p.PushComponent(begin, end, Multiply)

		case ruleAction93:

			p.PushComponent(begin, end, Divide)

		case ruleAction94:

			p.PushComponent(begin, end, Modulo)

		case ruleAction95:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction96:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction97:

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
		/* 1 Statement <- <(SelectUnionStmt / SelectStmt / SourceStmt / SinkStmt / StateStmt / StreamStmt)> */
		func() bool {
			position7, tokenIndex7, depth7 := position, tokenIndex, depth
			{
				position8 := position
				depth++
				{
					position9, tokenIndex9, depth9 := position, tokenIndex, depth
					if !_rules[ruleSelectUnionStmt]() {
						goto l10
					}
					goto l9
				l10:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleSelectStmt]() {
						goto l11
					}
					goto l9
				l11:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleSourceStmt]() {
						goto l12
					}
					goto l9
				l12:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleSinkStmt]() {
						goto l13
					}
					goto l9
				l13:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleStateStmt]() {
						goto l14
					}
					goto l9
				l14:
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
			position15, tokenIndex15, depth15 := position, tokenIndex, depth
			{
				position16 := position
				depth++
				{
					position17, tokenIndex17, depth17 := position, tokenIndex, depth
					if !_rules[ruleCreateSourceStmt]() {
						goto l18
					}
					goto l17
				l18:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleUpdateSourceStmt]() {
						goto l19
					}
					goto l17
				l19:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleDropSourceStmt]() {
						goto l20
					}
					goto l17
				l20:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[rulePauseSourceStmt]() {
						goto l21
					}
					goto l17
				l21:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleResumeSourceStmt]() {
						goto l22
					}
					goto l17
				l22:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleRewindSourceStmt]() {
						goto l15
					}
				}
			l17:
				depth--
				add(ruleSourceStmt, position16)
			}
			return true
		l15:
			position, tokenIndex, depth = position15, tokenIndex15, depth15
			return false
		},
		/* 3 SinkStmt <- <(CreateSinkStmt / UpdateSinkStmt / DropSinkStmt)> */
		func() bool {
			position23, tokenIndex23, depth23 := position, tokenIndex, depth
			{
				position24 := position
				depth++
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if !_rules[ruleCreateSinkStmt]() {
						goto l26
					}
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if !_rules[ruleUpdateSinkStmt]() {
						goto l27
					}
					goto l25
				l27:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if !_rules[ruleDropSinkStmt]() {
						goto l23
					}
				}
			l25:
				depth--
				add(ruleSinkStmt, position24)
			}
			return true
		l23:
			position, tokenIndex, depth = position23, tokenIndex23, depth23
			return false
		},
		/* 4 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt)> */
		func() bool {
			position28, tokenIndex28, depth28 := position, tokenIndex, depth
			{
				position29 := position
				depth++
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if !_rules[ruleCreateStateStmt]() {
						goto l31
					}
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleUpdateStateStmt]() {
						goto l32
					}
					goto l30
				l32:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleDropStateStmt]() {
						goto l28
					}
				}
			l30:
				depth--
				add(ruleStateStmt, position29)
			}
			return true
		l28:
			position, tokenIndex, depth = position28, tokenIndex28, depth28
			return false
		},
		/* 5 StreamStmt <- <(CreateStreamAsSelectUnionStmt / CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoSelectStmt / InsertIntoFromStmt)> */
		func() bool {
			position33, tokenIndex33, depth33 := position, tokenIndex, depth
			{
				position34 := position
				depth++
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l36
					}
					goto l35
				l36:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l37
					}
					goto l35
				l37:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleDropStreamStmt]() {
						goto l38
					}
					goto l35
				l38:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l39
					}
					goto l35
				l39:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l33
					}
				}
			l35:
				depth--
				add(ruleStreamStmt, position34)
			}
			return true
		l33:
			position, tokenIndex, depth = position33, tokenIndex33, depth33
			return false
		},
		/* 6 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
		func() bool {
			position40, tokenIndex40, depth40 := position, tokenIndex, depth
			{
				position41 := position
				depth++
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('S') {
						goto l40
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
						goto l40
					}
					position++
				}
			l44:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('L') {
						goto l40
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
						goto l40
					}
					position++
				}
			l48:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('C') {
						goto l40
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
						goto l40
					}
					position++
				}
			l52:
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
				if !_rules[ruleAction0]() {
					goto l40
				}
				depth--
				add(ruleSelectStmt, position41)
			}
			return true
		l40:
			position, tokenIndex, depth = position40, tokenIndex40, depth40
			return false
		},
		/* 7 SelectUnionStmt <- <(<(SelectStmt (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action1)> */
		func() bool {
			position54, tokenIndex54, depth54 := position, tokenIndex, depth
			{
				position55 := position
				depth++
				{
					position56 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l54
					}
					{
						position59, tokenIndex59, depth59 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l60
						}
						position++
						goto l59
					l60:
						position, tokenIndex, depth = position59, tokenIndex59, depth59
						if buffer[position] != rune('U') {
							goto l54
						}
						position++
					}
				l59:
					{
						position61, tokenIndex61, depth61 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l62
						}
						position++
						goto l61
					l62:
						position, tokenIndex, depth = position61, tokenIndex61, depth61
						if buffer[position] != rune('N') {
							goto l54
						}
						position++
					}
				l61:
					{
						position63, tokenIndex63, depth63 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l64
						}
						position++
						goto l63
					l64:
						position, tokenIndex, depth = position63, tokenIndex63, depth63
						if buffer[position] != rune('I') {
							goto l54
						}
						position++
					}
				l63:
					{
						position65, tokenIndex65, depth65 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l66
						}
						position++
						goto l65
					l66:
						position, tokenIndex, depth = position65, tokenIndex65, depth65
						if buffer[position] != rune('O') {
							goto l54
						}
						position++
					}
				l65:
					{
						position67, tokenIndex67, depth67 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l68
						}
						position++
						goto l67
					l68:
						position, tokenIndex, depth = position67, tokenIndex67, depth67
						if buffer[position] != rune('N') {
							goto l54
						}
						position++
					}
				l67:
					if !_rules[rulesp]() {
						goto l54
					}
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l70
						}
						position++
						goto l69
					l70:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
						if buffer[position] != rune('A') {
							goto l54
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
							goto l54
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
							goto l54
						}
						position++
					}
				l73:
					if !_rules[rulesp]() {
						goto l54
					}
					if !_rules[ruleSelectStmt]() {
						goto l54
					}
				l57:
					{
						position58, tokenIndex58, depth58 := position, tokenIndex, depth
						{
							position75, tokenIndex75, depth75 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l76
							}
							position++
							goto l75
						l76:
							position, tokenIndex, depth = position75, tokenIndex75, depth75
							if buffer[position] != rune('U') {
								goto l58
							}
							position++
						}
					l75:
						{
							position77, tokenIndex77, depth77 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l78
							}
							position++
							goto l77
						l78:
							position, tokenIndex, depth = position77, tokenIndex77, depth77
							if buffer[position] != rune('N') {
								goto l58
							}
							position++
						}
					l77:
						{
							position79, tokenIndex79, depth79 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l80
							}
							position++
							goto l79
						l80:
							position, tokenIndex, depth = position79, tokenIndex79, depth79
							if buffer[position] != rune('I') {
								goto l58
							}
							position++
						}
					l79:
						{
							position81, tokenIndex81, depth81 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l82
							}
							position++
							goto l81
						l82:
							position, tokenIndex, depth = position81, tokenIndex81, depth81
							if buffer[position] != rune('O') {
								goto l58
							}
							position++
						}
					l81:
						{
							position83, tokenIndex83, depth83 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l84
							}
							position++
							goto l83
						l84:
							position, tokenIndex, depth = position83, tokenIndex83, depth83
							if buffer[position] != rune('N') {
								goto l58
							}
							position++
						}
					l83:
						if !_rules[rulesp]() {
							goto l58
						}
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
								goto l58
							}
							position++
						}
					l85:
						{
							position87, tokenIndex87, depth87 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l88
							}
							position++
							goto l87
						l88:
							position, tokenIndex, depth = position87, tokenIndex87, depth87
							if buffer[position] != rune('L') {
								goto l58
							}
							position++
						}
					l87:
						{
							position89, tokenIndex89, depth89 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l90
							}
							position++
							goto l89
						l90:
							position, tokenIndex, depth = position89, tokenIndex89, depth89
							if buffer[position] != rune('L') {
								goto l58
							}
							position++
						}
					l89:
						if !_rules[rulesp]() {
							goto l58
						}
						if !_rules[ruleSelectStmt]() {
							goto l58
						}
						goto l57
					l58:
						position, tokenIndex, depth = position58, tokenIndex58, depth58
					}
					depth--
					add(rulePegText, position56)
				}
				if !_rules[ruleAction1]() {
					goto l54
				}
				depth--
				add(ruleSelectUnionStmt, position55)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 8 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action2)> */
		func() bool {
			position91, tokenIndex91, depth91 := position, tokenIndex, depth
			{
				position92 := position
				depth++
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('C') {
						goto l91
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
						goto l91
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('E') {
						goto l91
					}
					position++
				}
			l97:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('A') {
						goto l91
					}
					position++
				}
			l99:
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
						goto l91
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
						goto l91
					}
					position++
				}
			l103:
				if !_rules[rulesp]() {
					goto l91
				}
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('S') {
						goto l91
					}
					position++
				}
			l105:
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
						goto l91
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('R') {
						goto l91
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
						goto l91
					}
					position++
				}
			l111:
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l114
					}
					position++
					goto l113
				l114:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
					if buffer[position] != rune('A') {
						goto l91
					}
					position++
				}
			l113:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('M') {
						goto l91
					}
					position++
				}
			l115:
				if !_rules[rulesp]() {
					goto l91
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l91
				}
				if !_rules[rulesp]() {
					goto l91
				}
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
						goto l91
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('S') {
						goto l91
					}
					position++
				}
			l119:
				if !_rules[rulesp]() {
					goto l91
				}
				if !_rules[ruleSelectStmt]() {
					goto l91
				}
				if !_rules[ruleAction2]() {
					goto l91
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position92)
			}
			return true
		l91:
			position, tokenIndex, depth = position91, tokenIndex91, depth91
			return false
		},
		/* 9 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action3)> */
		func() bool {
			position121, tokenIndex121, depth121 := position, tokenIndex, depth
			{
				position122 := position
				depth++
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('C') {
						goto l121
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('R') {
						goto l121
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
						goto l121
					}
					position++
				}
			l127:
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('A') {
						goto l121
					}
					position++
				}
			l129:
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
						goto l121
					}
					position++
				}
			l131:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('E') {
						goto l121
					}
					position++
				}
			l133:
				if !_rules[rulesp]() {
					goto l121
				}
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('S') {
						goto l121
					}
					position++
				}
			l135:
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
						goto l121
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('R') {
						goto l121
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
						goto l121
					}
					position++
				}
			l141:
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					if buffer[position] != rune('A') {
						goto l121
					}
					position++
				}
			l143:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('M') {
						goto l121
					}
					position++
				}
			l145:
				if !_rules[rulesp]() {
					goto l121
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l121
				}
				if !_rules[rulesp]() {
					goto l121
				}
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
						goto l121
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('S') {
						goto l121
					}
					position++
				}
			l149:
				if !_rules[rulesp]() {
					goto l121
				}
				if !_rules[ruleSelectUnionStmt]() {
					goto l121
				}
				if !_rules[ruleAction3]() {
					goto l121
				}
				depth--
				add(ruleCreateStreamAsSelectUnionStmt, position122)
			}
			return true
		l121:
			position, tokenIndex, depth = position121, tokenIndex121, depth121
			return false
		},
		/* 10 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position151, tokenIndex151, depth151 := position, tokenIndex, depth
			{
				position152 := position
				depth++
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('C') {
						goto l151
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('R') {
						goto l151
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
						goto l151
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
						goto l151
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
						goto l151
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
						goto l151
					}
					position++
				}
			l163:
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[rulePausedOpt]() {
					goto l151
				}
				if !_rules[rulesp]() {
					goto l151
				}
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('S') {
						goto l151
					}
					position++
				}
			l165:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('O') {
						goto l151
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('U') {
						goto l151
					}
					position++
				}
			l169:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('R') {
						goto l151
					}
					position++
				}
			l171:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('C') {
						goto l151
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
						goto l151
					}
					position++
				}
			l175:
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l151
				}
				if !_rules[rulesp]() {
					goto l151
				}
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('T') {
						goto l151
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('Y') {
						goto l151
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
						goto l151
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
						goto l151
					}
					position++
				}
			l183:
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[ruleSourceSinkType]() {
					goto l151
				}
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l151
				}
				if !_rules[ruleAction4]() {
					goto l151
				}
				depth--
				add(ruleCreateSourceStmt, position152)
			}
			return true
		l151:
			position, tokenIndex, depth = position151, tokenIndex151, depth151
			return false
		},
		/* 11 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action5)> */
		func() bool {
			position185, tokenIndex185, depth185 := position, tokenIndex, depth
			{
				position186 := position
				depth++
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if buffer[position] != rune('C') {
						goto l185
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('R') {
						goto l185
					}
					position++
				}
			l189:
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('E') {
						goto l185
					}
					position++
				}
			l191:
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('A') {
						goto l185
					}
					position++
				}
			l193:
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('T') {
						goto l185
					}
					position++
				}
			l195:
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('E') {
						goto l185
					}
					position++
				}
			l197:
				if !_rules[rulesp]() {
					goto l185
				}
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
						goto l185
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('I') {
						goto l185
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('N') {
						goto l185
					}
					position++
				}
			l203:
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('K') {
						goto l185
					}
					position++
				}
			l205:
				if !_rules[rulesp]() {
					goto l185
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l185
				}
				if !_rules[rulesp]() {
					goto l185
				}
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
						goto l185
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('Y') {
						goto l185
					}
					position++
				}
			l209:
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('P') {
						goto l185
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
						goto l185
					}
					position++
				}
			l213:
				if !_rules[rulesp]() {
					goto l185
				}
				if !_rules[ruleSourceSinkType]() {
					goto l185
				}
				if !_rules[rulesp]() {
					goto l185
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l185
				}
				if !_rules[ruleAction5]() {
					goto l185
				}
				depth--
				add(ruleCreateSinkStmt, position186)
			}
			return true
		l185:
			position, tokenIndex, depth = position185, tokenIndex185, depth185
			return false
		},
		/* 12 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action6)> */
		func() bool {
			position215, tokenIndex215, depth215 := position, tokenIndex, depth
			{
				position216 := position
				depth++
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('C') {
						goto l215
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
						goto l215
					}
					position++
				}
			l219:
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if buffer[position] != rune('E') {
						goto l215
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('A') {
						goto l215
					}
					position++
				}
			l223:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('T') {
						goto l215
					}
					position++
				}
			l225:
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l228
					}
					position++
					goto l227
				l228:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
					if buffer[position] != rune('E') {
						goto l215
					}
					position++
				}
			l227:
				if !_rules[rulesp]() {
					goto l215
				}
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('S') {
						goto l215
					}
					position++
				}
			l229:
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if buffer[position] != rune('T') {
						goto l215
					}
					position++
				}
			l231:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('A') {
						goto l215
					}
					position++
				}
			l233:
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l236
					}
					position++
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if buffer[position] != rune('T') {
						goto l215
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
						goto l215
					}
					position++
				}
			l237:
				if !_rules[rulesp]() {
					goto l215
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l215
				}
				if !_rules[rulesp]() {
					goto l215
				}
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
						goto l215
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('Y') {
						goto l215
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('P') {
						goto l215
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
						goto l215
					}
					position++
				}
			l245:
				if !_rules[rulesp]() {
					goto l215
				}
				if !_rules[ruleSourceSinkType]() {
					goto l215
				}
				if !_rules[rulesp]() {
					goto l215
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l215
				}
				if !_rules[ruleAction6]() {
					goto l215
				}
				depth--
				add(ruleCreateStateStmt, position216)
			}
			return true
		l215:
			position, tokenIndex, depth = position215, tokenIndex215, depth215
			return false
		},
		/* 13 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action7)> */
		func() bool {
			position247, tokenIndex247, depth247 := position, tokenIndex, depth
			{
				position248 := position
				depth++
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('U') {
						goto l247
					}
					position++
				}
			l249:
				{
					position251, tokenIndex251, depth251 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l252
					}
					position++
					goto l251
				l252:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('P') {
						goto l247
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('D') {
						goto l247
					}
					position++
				}
			l253:
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l256
					}
					position++
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('A') {
						goto l247
					}
					position++
				}
			l255:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('T') {
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
					if buffer[position] != rune('t') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if buffer[position] != rune('T') {
						goto l247
					}
					position++
				}
			l263:
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('A') {
						goto l247
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
						goto l247
					}
					position++
				}
			l267:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('E') {
						goto l247
					}
					position++
				}
			l269:
				if !_rules[rulesp]() {
					goto l247
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l247
				}
				if !_rules[rulesp]() {
					goto l247
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l247
				}
				if !_rules[ruleAction7]() {
					goto l247
				}
				depth--
				add(ruleUpdateStateStmt, position248)
			}
			return true
		l247:
			position, tokenIndex, depth = position247, tokenIndex247, depth247
			return false
		},
		/* 14 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action8)> */
		func() bool {
			position271, tokenIndex271, depth271 := position, tokenIndex, depth
			{
				position272 := position
				depth++
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('U') {
						goto l271
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('P') {
						goto l271
					}
					position++
				}
			l275:
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('D') {
						goto l271
					}
					position++
				}
			l277:
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('A') {
						goto l271
					}
					position++
				}
			l279:
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('T') {
						goto l271
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('E') {
						goto l271
					}
					position++
				}
			l283:
				if !_rules[rulesp]() {
					goto l271
				}
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('S') {
						goto l271
					}
					position++
				}
			l285:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('O') {
						goto l271
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('U') {
						goto l271
					}
					position++
				}
			l289:
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('R') {
						goto l271
					}
					position++
				}
			l291:
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					if buffer[position] != rune('C') {
						goto l271
					}
					position++
				}
			l293:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('E') {
						goto l271
					}
					position++
				}
			l295:
				if !_rules[rulesp]() {
					goto l271
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l271
				}
				if !_rules[rulesp]() {
					goto l271
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l271
				}
				if !_rules[ruleAction8]() {
					goto l271
				}
				depth--
				add(ruleUpdateSourceStmt, position272)
			}
			return true
		l271:
			position, tokenIndex, depth = position271, tokenIndex271, depth271
			return false
		},
		/* 15 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position297, tokenIndex297, depth297 := position, tokenIndex, depth
			{
				position298 := position
				depth++
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l300
					}
					position++
					goto l299
				l300:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
					if buffer[position] != rune('U') {
						goto l297
					}
					position++
				}
			l299:
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l302
					}
					position++
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if buffer[position] != rune('P') {
						goto l297
					}
					position++
				}
			l301:
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
						goto l297
					}
					position++
				}
			l303:
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l306
					}
					position++
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if buffer[position] != rune('A') {
						goto l297
					}
					position++
				}
			l305:
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l308
					}
					position++
					goto l307
				l308:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
					if buffer[position] != rune('T') {
						goto l297
					}
					position++
				}
			l307:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('E') {
						goto l297
					}
					position++
				}
			l309:
				if !_rules[rulesp]() {
					goto l297
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
						goto l297
					}
					position++
				}
			l311:
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('I') {
						goto l297
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('N') {
						goto l297
					}
					position++
				}
			l315:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('K') {
						goto l297
					}
					position++
				}
			l317:
				if !_rules[rulesp]() {
					goto l297
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l297
				}
				if !_rules[rulesp]() {
					goto l297
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l297
				}
				if !_rules[ruleAction9]() {
					goto l297
				}
				depth--
				add(ruleUpdateSinkStmt, position298)
			}
			return true
		l297:
			position, tokenIndex, depth = position297, tokenIndex297, depth297
			return false
		},
		/* 16 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action10)> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l322
					}
					position++
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if buffer[position] != rune('I') {
						goto l319
					}
					position++
				}
			l321:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l324
					}
					position++
					goto l323
				l324:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					if buffer[position] != rune('N') {
						goto l319
					}
					position++
				}
			l323:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l326
					}
					position++
					goto l325
				l326:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
					if buffer[position] != rune('S') {
						goto l319
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
						goto l319
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
						goto l319
					}
					position++
				}
			l329:
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('T') {
						goto l319
					}
					position++
				}
			l331:
				if !_rules[rulesp]() {
					goto l319
				}
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('I') {
						goto l319
					}
					position++
				}
			l333:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l336
					}
					position++
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if buffer[position] != rune('N') {
						goto l319
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
						goto l319
					}
					position++
				}
			l337:
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l340
					}
					position++
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if buffer[position] != rune('O') {
						goto l319
					}
					position++
				}
			l339:
				if !_rules[rulesp]() {
					goto l319
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l319
				}
				if !_rules[rulesp]() {
					goto l319
				}
				if !_rules[ruleSelectStmt]() {
					goto l319
				}
				if !_rules[ruleAction10]() {
					goto l319
				}
				depth--
				add(ruleInsertIntoSelectStmt, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 17 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action11)> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
			{
				position342 := position
				depth++
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('I') {
						goto l341
					}
					position++
				}
			l343:
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l346
					}
					position++
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if buffer[position] != rune('N') {
						goto l341
					}
					position++
				}
			l345:
				{
					position347, tokenIndex347, depth347 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l348
					}
					position++
					goto l347
				l348:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if buffer[position] != rune('S') {
						goto l341
					}
					position++
				}
			l347:
				{
					position349, tokenIndex349, depth349 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l350
					}
					position++
					goto l349
				l350:
					position, tokenIndex, depth = position349, tokenIndex349, depth349
					if buffer[position] != rune('E') {
						goto l341
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
						goto l341
					}
					position++
				}
			l351:
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l354
					}
					position++
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if buffer[position] != rune('T') {
						goto l341
					}
					position++
				}
			l353:
				if !_rules[rulesp]() {
					goto l341
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
						goto l341
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
						goto l341
					}
					position++
				}
			l357:
				{
					position359, tokenIndex359, depth359 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l360
					}
					position++
					goto l359
				l360:
					position, tokenIndex, depth = position359, tokenIndex359, depth359
					if buffer[position] != rune('T') {
						goto l341
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
						goto l341
					}
					position++
				}
			l361:
				if !_rules[rulesp]() {
					goto l341
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l341
				}
				if !_rules[rulesp]() {
					goto l341
				}
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l364
					}
					position++
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					if buffer[position] != rune('F') {
						goto l341
					}
					position++
				}
			l363:
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('R') {
						goto l341
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
						goto l341
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('M') {
						goto l341
					}
					position++
				}
			l369:
				if !_rules[rulesp]() {
					goto l341
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l341
				}
				if !_rules[ruleAction11]() {
					goto l341
				}
				depth--
				add(ruleInsertIntoFromStmt, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
			return false
		},
		/* 18 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action12)> */
		func() bool {
			position371, tokenIndex371, depth371 := position, tokenIndex, depth
			{
				position372 := position
				depth++
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
						goto l371
					}
					position++
				}
			l373:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l376
					}
					position++
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					if buffer[position] != rune('A') {
						goto l371
					}
					position++
				}
			l375:
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l378
					}
					position++
					goto l377
				l378:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
					if buffer[position] != rune('U') {
						goto l371
					}
					position++
				}
			l377:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('S') {
						goto l371
					}
					position++
				}
			l379:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('E') {
						goto l371
					}
					position++
				}
			l381:
				if !_rules[rulesp]() {
					goto l371
				}
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('S') {
						goto l371
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('O') {
						goto l371
					}
					position++
				}
			l385:
				{
					position387, tokenIndex387, depth387 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l388
					}
					position++
					goto l387
				l388:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
					if buffer[position] != rune('U') {
						goto l371
					}
					position++
				}
			l387:
				{
					position389, tokenIndex389, depth389 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l390
					}
					position++
					goto l389
				l390:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
					if buffer[position] != rune('R') {
						goto l371
					}
					position++
				}
			l389:
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
					if buffer[position] != rune('C') {
						goto l371
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
						goto l371
					}
					position++
				}
			l393:
				if !_rules[rulesp]() {
					goto l371
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l371
				}
				if !_rules[ruleAction12]() {
					goto l371
				}
				depth--
				add(rulePauseSourceStmt, position372)
			}
			return true
		l371:
			position, tokenIndex, depth = position371, tokenIndex371, depth371
			return false
		},
		/* 19 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position395, tokenIndex395, depth395 := position, tokenIndex, depth
			{
				position396 := position
				depth++
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l398
					}
					position++
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if buffer[position] != rune('R') {
						goto l395
					}
					position++
				}
			l397:
				{
					position399, tokenIndex399, depth399 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l400
					}
					position++
					goto l399
				l400:
					position, tokenIndex, depth = position399, tokenIndex399, depth399
					if buffer[position] != rune('E') {
						goto l395
					}
					position++
				}
			l399:
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l402
					}
					position++
					goto l401
				l402:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if buffer[position] != rune('S') {
						goto l395
					}
					position++
				}
			l401:
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l404
					}
					position++
					goto l403
				l404:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
					if buffer[position] != rune('U') {
						goto l395
					}
					position++
				}
			l403:
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l406
					}
					position++
					goto l405
				l406:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if buffer[position] != rune('M') {
						goto l395
					}
					position++
				}
			l405:
				{
					position407, tokenIndex407, depth407 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l408
					}
					position++
					goto l407
				l408:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
					if buffer[position] != rune('E') {
						goto l395
					}
					position++
				}
			l407:
				if !_rules[rulesp]() {
					goto l395
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
						goto l395
					}
					position++
				}
			l409:
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l412
					}
					position++
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if buffer[position] != rune('O') {
						goto l395
					}
					position++
				}
			l411:
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l414
					}
					position++
					goto l413
				l414:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
					if buffer[position] != rune('U') {
						goto l395
					}
					position++
				}
			l413:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l416
					}
					position++
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if buffer[position] != rune('R') {
						goto l395
					}
					position++
				}
			l415:
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l418
					}
					position++
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if buffer[position] != rune('C') {
						goto l395
					}
					position++
				}
			l417:
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l420
					}
					position++
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if buffer[position] != rune('E') {
						goto l395
					}
					position++
				}
			l419:
				if !_rules[rulesp]() {
					goto l395
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l395
				}
				if !_rules[ruleAction13]() {
					goto l395
				}
				depth--
				add(ruleResumeSourceStmt, position396)
			}
			return true
		l395:
			position, tokenIndex, depth = position395, tokenIndex395, depth395
			return false
		},
		/* 20 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position421, tokenIndex421, depth421 := position, tokenIndex, depth
			{
				position422 := position
				depth++
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l424
					}
					position++
					goto l423
				l424:
					position, tokenIndex, depth = position423, tokenIndex423, depth423
					if buffer[position] != rune('R') {
						goto l421
					}
					position++
				}
			l423:
				{
					position425, tokenIndex425, depth425 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l426
					}
					position++
					goto l425
				l426:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if buffer[position] != rune('E') {
						goto l421
					}
					position++
				}
			l425:
				{
					position427, tokenIndex427, depth427 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l428
					}
					position++
					goto l427
				l428:
					position, tokenIndex, depth = position427, tokenIndex427, depth427
					if buffer[position] != rune('W') {
						goto l421
					}
					position++
				}
			l427:
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l430
					}
					position++
					goto l429
				l430:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if buffer[position] != rune('I') {
						goto l421
					}
					position++
				}
			l429:
				{
					position431, tokenIndex431, depth431 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l432
					}
					position++
					goto l431
				l432:
					position, tokenIndex, depth = position431, tokenIndex431, depth431
					if buffer[position] != rune('N') {
						goto l421
					}
					position++
				}
			l431:
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l434
					}
					position++
					goto l433
				l434:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if buffer[position] != rune('D') {
						goto l421
					}
					position++
				}
			l433:
				if !_rules[rulesp]() {
					goto l421
				}
				{
					position435, tokenIndex435, depth435 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l436
					}
					position++
					goto l435
				l436:
					position, tokenIndex, depth = position435, tokenIndex435, depth435
					if buffer[position] != rune('S') {
						goto l421
					}
					position++
				}
			l435:
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l438
					}
					position++
					goto l437
				l438:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
					if buffer[position] != rune('O') {
						goto l421
					}
					position++
				}
			l437:
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l440
					}
					position++
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if buffer[position] != rune('U') {
						goto l421
					}
					position++
				}
			l439:
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l442
					}
					position++
					goto l441
				l442:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
					if buffer[position] != rune('R') {
						goto l421
					}
					position++
				}
			l441:
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l444
					}
					position++
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if buffer[position] != rune('C') {
						goto l421
					}
					position++
				}
			l443:
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
						goto l421
					}
					position++
				}
			l445:
				if !_rules[rulesp]() {
					goto l421
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l421
				}
				if !_rules[ruleAction14]() {
					goto l421
				}
				depth--
				add(ruleRewindSourceStmt, position422)
			}
			return true
		l421:
			position, tokenIndex, depth = position421, tokenIndex421, depth421
			return false
		},
		/* 21 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position447, tokenIndex447, depth447 := position, tokenIndex, depth
			{
				position448 := position
				depth++
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l450
					}
					position++
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					if buffer[position] != rune('D') {
						goto l447
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
						goto l447
					}
					position++
				}
			l451:
				{
					position453, tokenIndex453, depth453 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l454
					}
					position++
					goto l453
				l454:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
					if buffer[position] != rune('O') {
						goto l447
					}
					position++
				}
			l453:
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l456
					}
					position++
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if buffer[position] != rune('P') {
						goto l447
					}
					position++
				}
			l455:
				if !_rules[rulesp]() {
					goto l447
				}
				{
					position457, tokenIndex457, depth457 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l458
					}
					position++
					goto l457
				l458:
					position, tokenIndex, depth = position457, tokenIndex457, depth457
					if buffer[position] != rune('S') {
						goto l447
					}
					position++
				}
			l457:
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l460
					}
					position++
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if buffer[position] != rune('O') {
						goto l447
					}
					position++
				}
			l459:
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l462
					}
					position++
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('U') {
						goto l447
					}
					position++
				}
			l461:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('R') {
						goto l447
					}
					position++
				}
			l463:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l466
					}
					position++
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					if buffer[position] != rune('C') {
						goto l447
					}
					position++
				}
			l465:
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l468
					}
					position++
					goto l467
				l468:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
					if buffer[position] != rune('E') {
						goto l447
					}
					position++
				}
			l467:
				if !_rules[rulesp]() {
					goto l447
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l447
				}
				if !_rules[ruleAction15]() {
					goto l447
				}
				depth--
				add(ruleDropSourceStmt, position448)
			}
			return true
		l447:
			position, tokenIndex, depth = position447, tokenIndex447, depth447
			return false
		},
		/* 22 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action16)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l472
					}
					position++
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if buffer[position] != rune('D') {
						goto l469
					}
					position++
				}
			l471:
				{
					position473, tokenIndex473, depth473 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l474
					}
					position++
					goto l473
				l474:
					position, tokenIndex, depth = position473, tokenIndex473, depth473
					if buffer[position] != rune('R') {
						goto l469
					}
					position++
				}
			l473:
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l476
					}
					position++
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if buffer[position] != rune('O') {
						goto l469
					}
					position++
				}
			l475:
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l478
					}
					position++
					goto l477
				l478:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
					if buffer[position] != rune('P') {
						goto l469
					}
					position++
				}
			l477:
				if !_rules[rulesp]() {
					goto l469
				}
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l480
					}
					position++
					goto l479
				l480:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if buffer[position] != rune('S') {
						goto l469
					}
					position++
				}
			l479:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if buffer[position] != rune('T') {
						goto l469
					}
					position++
				}
			l481:
				{
					position483, tokenIndex483, depth483 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l484
					}
					position++
					goto l483
				l484:
					position, tokenIndex, depth = position483, tokenIndex483, depth483
					if buffer[position] != rune('R') {
						goto l469
					}
					position++
				}
			l483:
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l486
					}
					position++
					goto l485
				l486:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					if buffer[position] != rune('E') {
						goto l469
					}
					position++
				}
			l485:
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l488
					}
					position++
					goto l487
				l488:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
					if buffer[position] != rune('A') {
						goto l469
					}
					position++
				}
			l487:
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l490
					}
					position++
					goto l489
				l490:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
					if buffer[position] != rune('M') {
						goto l469
					}
					position++
				}
			l489:
				if !_rules[rulesp]() {
					goto l469
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l469
				}
				if !_rules[ruleAction16]() {
					goto l469
				}
				depth--
				add(ruleDropStreamStmt, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 23 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action17)> */
		func() bool {
			position491, tokenIndex491, depth491 := position, tokenIndex, depth
			{
				position492 := position
				depth++
				{
					position493, tokenIndex493, depth493 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l494
					}
					position++
					goto l493
				l494:
					position, tokenIndex, depth = position493, tokenIndex493, depth493
					if buffer[position] != rune('D') {
						goto l491
					}
					position++
				}
			l493:
				{
					position495, tokenIndex495, depth495 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l496
					}
					position++
					goto l495
				l496:
					position, tokenIndex, depth = position495, tokenIndex495, depth495
					if buffer[position] != rune('R') {
						goto l491
					}
					position++
				}
			l495:
				{
					position497, tokenIndex497, depth497 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l498
					}
					position++
					goto l497
				l498:
					position, tokenIndex, depth = position497, tokenIndex497, depth497
					if buffer[position] != rune('O') {
						goto l491
					}
					position++
				}
			l497:
				{
					position499, tokenIndex499, depth499 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l500
					}
					position++
					goto l499
				l500:
					position, tokenIndex, depth = position499, tokenIndex499, depth499
					if buffer[position] != rune('P') {
						goto l491
					}
					position++
				}
			l499:
				if !_rules[rulesp]() {
					goto l491
				}
				{
					position501, tokenIndex501, depth501 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l502
					}
					position++
					goto l501
				l502:
					position, tokenIndex, depth = position501, tokenIndex501, depth501
					if buffer[position] != rune('S') {
						goto l491
					}
					position++
				}
			l501:
				{
					position503, tokenIndex503, depth503 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l504
					}
					position++
					goto l503
				l504:
					position, tokenIndex, depth = position503, tokenIndex503, depth503
					if buffer[position] != rune('I') {
						goto l491
					}
					position++
				}
			l503:
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l506
					}
					position++
					goto l505
				l506:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
					if buffer[position] != rune('N') {
						goto l491
					}
					position++
				}
			l505:
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('K') {
						goto l491
					}
					position++
				}
			l507:
				if !_rules[rulesp]() {
					goto l491
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l491
				}
				if !_rules[ruleAction17]() {
					goto l491
				}
				depth--
				add(ruleDropSinkStmt, position492)
			}
			return true
		l491:
			position, tokenIndex, depth = position491, tokenIndex491, depth491
			return false
		},
		/* 24 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action18)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if buffer[position] != rune('D') {
						goto l509
					}
					position++
				}
			l511:
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l514
					}
					position++
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if buffer[position] != rune('R') {
						goto l509
					}
					position++
				}
			l513:
				{
					position515, tokenIndex515, depth515 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l516
					}
					position++
					goto l515
				l516:
					position, tokenIndex, depth = position515, tokenIndex515, depth515
					if buffer[position] != rune('O') {
						goto l509
					}
					position++
				}
			l515:
				{
					position517, tokenIndex517, depth517 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l518
					}
					position++
					goto l517
				l518:
					position, tokenIndex, depth = position517, tokenIndex517, depth517
					if buffer[position] != rune('P') {
						goto l509
					}
					position++
				}
			l517:
				if !_rules[rulesp]() {
					goto l509
				}
				{
					position519, tokenIndex519, depth519 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l520
					}
					position++
					goto l519
				l520:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if buffer[position] != rune('S') {
						goto l509
					}
					position++
				}
			l519:
				{
					position521, tokenIndex521, depth521 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l522
					}
					position++
					goto l521
				l522:
					position, tokenIndex, depth = position521, tokenIndex521, depth521
					if buffer[position] != rune('T') {
						goto l509
					}
					position++
				}
			l521:
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
						goto l509
					}
					position++
				}
			l523:
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l526
					}
					position++
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if buffer[position] != rune('T') {
						goto l509
					}
					position++
				}
			l525:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l528
					}
					position++
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if buffer[position] != rune('E') {
						goto l509
					}
					position++
				}
			l527:
				if !_rules[rulesp]() {
					goto l509
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l509
				}
				if !_rules[ruleAction18]() {
					goto l509
				}
				depth--
				add(ruleDropStateStmt, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 25 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) Action19)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l532
					}
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleDSTREAM]() {
						goto l533
					}
					goto l531
				l533:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleRSTREAM]() {
						goto l529
					}
				}
			l531:
				if !_rules[ruleAction19]() {
					goto l529
				}
				depth--
				add(ruleEmitter, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 26 Projections <- <(<(Projection sp (',' sp Projection)*)> Action20)> */
		func() bool {
			position534, tokenIndex534, depth534 := position, tokenIndex, depth
			{
				position535 := position
				depth++
				{
					position536 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l534
					}
					if !_rules[rulesp]() {
						goto l534
					}
				l537:
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l538
						}
						position++
						if !_rules[rulesp]() {
							goto l538
						}
						if !_rules[ruleProjection]() {
							goto l538
						}
						goto l537
					l538:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
					}
					depth--
					add(rulePegText, position536)
				}
				if !_rules[ruleAction20]() {
					goto l534
				}
				depth--
				add(ruleProjections, position535)
			}
			return true
		l534:
			position, tokenIndex, depth = position534, tokenIndex534, depth534
			return false
		},
		/* 27 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position539, tokenIndex539, depth539 := position, tokenIndex, depth
			{
				position540 := position
				depth++
				{
					position541, tokenIndex541, depth541 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l542
					}
					goto l541
				l542:
					position, tokenIndex, depth = position541, tokenIndex541, depth541
					if !_rules[ruleExpression]() {
						goto l543
					}
					goto l541
				l543:
					position, tokenIndex, depth = position541, tokenIndex541, depth541
					if !_rules[ruleWildcard]() {
						goto l539
					}
				}
			l541:
				depth--
				add(ruleProjection, position540)
			}
			return true
		l539:
			position, tokenIndex, depth = position539, tokenIndex539, depth539
			return false
		},
		/* 28 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action21)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				{
					position546, tokenIndex546, depth546 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l547
					}
					goto l546
				l547:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
					if !_rules[ruleWildcard]() {
						goto l544
					}
				}
			l546:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position548, tokenIndex548, depth548 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l549
					}
					position++
					goto l548
				l549:
					position, tokenIndex, depth = position548, tokenIndex548, depth548
					if buffer[position] != rune('A') {
						goto l544
					}
					position++
				}
			l548:
				{
					position550, tokenIndex550, depth550 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l551
					}
					position++
					goto l550
				l551:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if buffer[position] != rune('S') {
						goto l544
					}
					position++
				}
			l550:
				if !_rules[rulesp]() {
					goto l544
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l544
				}
				if !_rules[ruleAction21]() {
					goto l544
				}
				depth--
				add(ruleAliasExpression, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 29 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action22)> */
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
							if buffer[position] != rune('f') {
								goto l558
							}
							position++
							goto l557
						l558:
							position, tokenIndex, depth = position557, tokenIndex557, depth557
							if buffer[position] != rune('F') {
								goto l555
							}
							position++
						}
					l557:
						{
							position559, tokenIndex559, depth559 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l560
							}
							position++
							goto l559
						l560:
							position, tokenIndex, depth = position559, tokenIndex559, depth559
							if buffer[position] != rune('R') {
								goto l555
							}
							position++
						}
					l559:
						{
							position561, tokenIndex561, depth561 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l562
							}
							position++
							goto l561
						l562:
							position, tokenIndex, depth = position561, tokenIndex561, depth561
							if buffer[position] != rune('O') {
								goto l555
							}
							position++
						}
					l561:
						{
							position563, tokenIndex563, depth563 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l564
							}
							position++
							goto l563
						l564:
							position, tokenIndex, depth = position563, tokenIndex563, depth563
							if buffer[position] != rune('M') {
								goto l555
							}
							position++
						}
					l563:
						if !_rules[rulesp]() {
							goto l555
						}
						if !_rules[ruleRelations]() {
							goto l555
						}
						if !_rules[rulesp]() {
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
				if !_rules[ruleAction22]() {
					goto l552
				}
				depth--
				add(ruleWindowedFrom, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 30 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567, tokenIndex567, depth567 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l568
					}
					goto l567
				l568:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleTuplesInterval]() {
						goto l565
					}
				}
			l567:
				depth--
				add(ruleInterval, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 31 TimeInterval <- <(NumericLiteral sp SECONDS Action23)> */
		func() bool {
			position569, tokenIndex569, depth569 := position, tokenIndex, depth
			{
				position570 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l569
				}
				if !_rules[rulesp]() {
					goto l569
				}
				if !_rules[ruleSECONDS]() {
					goto l569
				}
				if !_rules[ruleAction23]() {
					goto l569
				}
				depth--
				add(ruleTimeInterval, position570)
			}
			return true
		l569:
			position, tokenIndex, depth = position569, tokenIndex569, depth569
			return false
		},
		/* 32 TuplesInterval <- <(NumericLiteral sp TUPLES Action24)> */
		func() bool {
			position571, tokenIndex571, depth571 := position, tokenIndex, depth
			{
				position572 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l571
				}
				if !_rules[rulesp]() {
					goto l571
				}
				if !_rules[ruleTUPLES]() {
					goto l571
				}
				if !_rules[ruleAction24]() {
					goto l571
				}
				depth--
				add(ruleTuplesInterval, position572)
			}
			return true
		l571:
			position, tokenIndex, depth = position571, tokenIndex571, depth571
			return false
		},
		/* 33 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position573, tokenIndex573, depth573 := position, tokenIndex, depth
			{
				position574 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l573
				}
				if !_rules[rulesp]() {
					goto l573
				}
			l575:
				{
					position576, tokenIndex576, depth576 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l576
					}
					position++
					if !_rules[rulesp]() {
						goto l576
					}
					if !_rules[ruleRelationLike]() {
						goto l576
					}
					goto l575
				l576:
					position, tokenIndex, depth = position576, tokenIndex576, depth576
				}
				depth--
				add(ruleRelations, position574)
			}
			return true
		l573:
			position, tokenIndex, depth = position573, tokenIndex573, depth573
			return false
		},
		/* 34 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action25)> */
		func() bool {
			position577, tokenIndex577, depth577 := position, tokenIndex, depth
			{
				position578 := position
				depth++
				{
					position579 := position
					depth++
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						{
							position582, tokenIndex582, depth582 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l583
							}
							position++
							goto l582
						l583:
							position, tokenIndex, depth = position582, tokenIndex582, depth582
							if buffer[position] != rune('W') {
								goto l580
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
								goto l580
							}
							position++
						}
					l584:
						{
							position586, tokenIndex586, depth586 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l587
							}
							position++
							goto l586
						l587:
							position, tokenIndex, depth = position586, tokenIndex586, depth586
							if buffer[position] != rune('E') {
								goto l580
							}
							position++
						}
					l586:
						{
							position588, tokenIndex588, depth588 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l589
							}
							position++
							goto l588
						l589:
							position, tokenIndex, depth = position588, tokenIndex588, depth588
							if buffer[position] != rune('R') {
								goto l580
							}
							position++
						}
					l588:
						{
							position590, tokenIndex590, depth590 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l591
							}
							position++
							goto l590
						l591:
							position, tokenIndex, depth = position590, tokenIndex590, depth590
							if buffer[position] != rune('E') {
								goto l580
							}
							position++
						}
					l590:
						if !_rules[rulesp]() {
							goto l580
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction25]() {
					goto l577
				}
				depth--
				add(ruleFilter, position578)
			}
			return true
		l577:
			position, tokenIndex, depth = position577, tokenIndex577, depth577
			return false
		},
		/* 35 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action26)> */
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
						{
							position597, tokenIndex597, depth597 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l598
							}
							position++
							goto l597
						l598:
							position, tokenIndex, depth = position597, tokenIndex597, depth597
							if buffer[position] != rune('G') {
								goto l595
							}
							position++
						}
					l597:
						{
							position599, tokenIndex599, depth599 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l600
							}
							position++
							goto l599
						l600:
							position, tokenIndex, depth = position599, tokenIndex599, depth599
							if buffer[position] != rune('R') {
								goto l595
							}
							position++
						}
					l599:
						{
							position601, tokenIndex601, depth601 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l602
							}
							position++
							goto l601
						l602:
							position, tokenIndex, depth = position601, tokenIndex601, depth601
							if buffer[position] != rune('O') {
								goto l595
							}
							position++
						}
					l601:
						{
							position603, tokenIndex603, depth603 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l604
							}
							position++
							goto l603
						l604:
							position, tokenIndex, depth = position603, tokenIndex603, depth603
							if buffer[position] != rune('U') {
								goto l595
							}
							position++
						}
					l603:
						{
							position605, tokenIndex605, depth605 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l606
							}
							position++
							goto l605
						l606:
							position, tokenIndex, depth = position605, tokenIndex605, depth605
							if buffer[position] != rune('P') {
								goto l595
							}
							position++
						}
					l605:
						if !_rules[rulesp]() {
							goto l595
						}
						{
							position607, tokenIndex607, depth607 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l608
							}
							position++
							goto l607
						l608:
							position, tokenIndex, depth = position607, tokenIndex607, depth607
							if buffer[position] != rune('B') {
								goto l595
							}
							position++
						}
					l607:
						{
							position609, tokenIndex609, depth609 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l610
							}
							position++
							goto l609
						l610:
							position, tokenIndex, depth = position609, tokenIndex609, depth609
							if buffer[position] != rune('Y') {
								goto l595
							}
							position++
						}
					l609:
						if !_rules[rulesp]() {
							goto l595
						}
						if !_rules[ruleGroupList]() {
							goto l595
						}
						goto l596
					l595:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
					}
				l596:
					depth--
					add(rulePegText, position594)
				}
				if !_rules[ruleAction26]() {
					goto l592
				}
				depth--
				add(ruleGrouping, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 36 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position611, tokenIndex611, depth611 := position, tokenIndex, depth
			{
				position612 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l611
				}
				if !_rules[rulesp]() {
					goto l611
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
					if !_rules[ruleExpression]() {
						goto l614
					}
					goto l613
				l614:
					position, tokenIndex, depth = position614, tokenIndex614, depth614
				}
				depth--
				add(ruleGroupList, position612)
			}
			return true
		l611:
			position, tokenIndex, depth = position611, tokenIndex611, depth611
			return false
		},
		/* 37 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action27)> */
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
						{
							position620, tokenIndex620, depth620 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l621
							}
							position++
							goto l620
						l621:
							position, tokenIndex, depth = position620, tokenIndex620, depth620
							if buffer[position] != rune('H') {
								goto l618
							}
							position++
						}
					l620:
						{
							position622, tokenIndex622, depth622 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l623
							}
							position++
							goto l622
						l623:
							position, tokenIndex, depth = position622, tokenIndex622, depth622
							if buffer[position] != rune('A') {
								goto l618
							}
							position++
						}
					l622:
						{
							position624, tokenIndex624, depth624 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l625
							}
							position++
							goto l624
						l625:
							position, tokenIndex, depth = position624, tokenIndex624, depth624
							if buffer[position] != rune('V') {
								goto l618
							}
							position++
						}
					l624:
						{
							position626, tokenIndex626, depth626 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l627
							}
							position++
							goto l626
						l627:
							position, tokenIndex, depth = position626, tokenIndex626, depth626
							if buffer[position] != rune('I') {
								goto l618
							}
							position++
						}
					l626:
						{
							position628, tokenIndex628, depth628 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l629
							}
							position++
							goto l628
						l629:
							position, tokenIndex, depth = position628, tokenIndex628, depth628
							if buffer[position] != rune('N') {
								goto l618
							}
							position++
						}
					l628:
						{
							position630, tokenIndex630, depth630 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l631
							}
							position++
							goto l630
						l631:
							position, tokenIndex, depth = position630, tokenIndex630, depth630
							if buffer[position] != rune('G') {
								goto l618
							}
							position++
						}
					l630:
						if !_rules[rulesp]() {
							goto l618
						}
						if !_rules[ruleExpression]() {
							goto l618
						}
						goto l619
					l618:
						position, tokenIndex, depth = position618, tokenIndex618, depth618
					}
				l619:
					depth--
					add(rulePegText, position617)
				}
				if !_rules[ruleAction27]() {
					goto l615
				}
				depth--
				add(ruleHaving, position616)
			}
			return true
		l615:
			position, tokenIndex, depth = position615, tokenIndex615, depth615
			return false
		},
		/* 38 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action28))> */
		func() bool {
			position632, tokenIndex632, depth632 := position, tokenIndex, depth
			{
				position633 := position
				depth++
				{
					position634, tokenIndex634, depth634 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l635
					}
					goto l634
				l635:
					position, tokenIndex, depth = position634, tokenIndex634, depth634
					if !_rules[ruleStreamWindow]() {
						goto l632
					}
					if !_rules[ruleAction28]() {
						goto l632
					}
				}
			l634:
				depth--
				add(ruleRelationLike, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 39 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action29)> */
		func() bool {
			position636, tokenIndex636, depth636 := position, tokenIndex, depth
			{
				position637 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l636
				}
				if !_rules[rulesp]() {
					goto l636
				}
				{
					position638, tokenIndex638, depth638 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l639
					}
					position++
					goto l638
				l639:
					position, tokenIndex, depth = position638, tokenIndex638, depth638
					if buffer[position] != rune('A') {
						goto l636
					}
					position++
				}
			l638:
				{
					position640, tokenIndex640, depth640 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l641
					}
					position++
					goto l640
				l641:
					position, tokenIndex, depth = position640, tokenIndex640, depth640
					if buffer[position] != rune('S') {
						goto l636
					}
					position++
				}
			l640:
				if !_rules[rulesp]() {
					goto l636
				}
				if !_rules[ruleIdentifier]() {
					goto l636
				}
				if !_rules[ruleAction29]() {
					goto l636
				}
				depth--
				add(ruleAliasedStreamWindow, position637)
			}
			return true
		l636:
			position, tokenIndex, depth = position636, tokenIndex636, depth636
			return false
		},
		/* 40 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action30)> */
		func() bool {
			position642, tokenIndex642, depth642 := position, tokenIndex, depth
			{
				position643 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l642
				}
				if !_rules[rulesp]() {
					goto l642
				}
				if buffer[position] != rune('[') {
					goto l642
				}
				position++
				if !_rules[rulesp]() {
					goto l642
				}
				{
					position644, tokenIndex644, depth644 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l645
					}
					position++
					goto l644
				l645:
					position, tokenIndex, depth = position644, tokenIndex644, depth644
					if buffer[position] != rune('R') {
						goto l642
					}
					position++
				}
			l644:
				{
					position646, tokenIndex646, depth646 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l647
					}
					position++
					goto l646
				l647:
					position, tokenIndex, depth = position646, tokenIndex646, depth646
					if buffer[position] != rune('A') {
						goto l642
					}
					position++
				}
			l646:
				{
					position648, tokenIndex648, depth648 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l649
					}
					position++
					goto l648
				l649:
					position, tokenIndex, depth = position648, tokenIndex648, depth648
					if buffer[position] != rune('N') {
						goto l642
					}
					position++
				}
			l648:
				{
					position650, tokenIndex650, depth650 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l651
					}
					position++
					goto l650
				l651:
					position, tokenIndex, depth = position650, tokenIndex650, depth650
					if buffer[position] != rune('G') {
						goto l642
					}
					position++
				}
			l650:
				{
					position652, tokenIndex652, depth652 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l653
					}
					position++
					goto l652
				l653:
					position, tokenIndex, depth = position652, tokenIndex652, depth652
					if buffer[position] != rune('E') {
						goto l642
					}
					position++
				}
			l652:
				if !_rules[rulesp]() {
					goto l642
				}
				if !_rules[ruleInterval]() {
					goto l642
				}
				if !_rules[rulesp]() {
					goto l642
				}
				if buffer[position] != rune(']') {
					goto l642
				}
				position++
				if !_rules[ruleAction30]() {
					goto l642
				}
				depth--
				add(ruleStreamWindow, position643)
			}
			return true
		l642:
			position, tokenIndex, depth = position642, tokenIndex642, depth642
			return false
		},
		/* 41 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position654, tokenIndex654, depth654 := position, tokenIndex, depth
			{
				position655 := position
				depth++
				{
					position656, tokenIndex656, depth656 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l657
					}
					goto l656
				l657:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleStream]() {
						goto l654
					}
				}
			l656:
				depth--
				add(ruleStreamLike, position655)
			}
			return true
		l654:
			position, tokenIndex, depth = position654, tokenIndex654, depth654
			return false
		},
		/* 42 UDSFFuncApp <- <(FuncApp Action31)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l658
				}
				if !_rules[ruleAction31]() {
					goto l658
				}
				depth--
				add(ruleUDSFFuncApp, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 43 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action32)> */
		func() bool {
			position660, tokenIndex660, depth660 := position, tokenIndex, depth
			{
				position661 := position
				depth++
				{
					position662 := position
					depth++
					{
						position663, tokenIndex663, depth663 := position, tokenIndex, depth
						{
							position665, tokenIndex665, depth665 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l666
							}
							position++
							goto l665
						l666:
							position, tokenIndex, depth = position665, tokenIndex665, depth665
							if buffer[position] != rune('W') {
								goto l663
							}
							position++
						}
					l665:
						{
							position667, tokenIndex667, depth667 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l668
							}
							position++
							goto l667
						l668:
							position, tokenIndex, depth = position667, tokenIndex667, depth667
							if buffer[position] != rune('I') {
								goto l663
							}
							position++
						}
					l667:
						{
							position669, tokenIndex669, depth669 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l670
							}
							position++
							goto l669
						l670:
							position, tokenIndex, depth = position669, tokenIndex669, depth669
							if buffer[position] != rune('T') {
								goto l663
							}
							position++
						}
					l669:
						{
							position671, tokenIndex671, depth671 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l672
							}
							position++
							goto l671
						l672:
							position, tokenIndex, depth = position671, tokenIndex671, depth671
							if buffer[position] != rune('H') {
								goto l663
							}
							position++
						}
					l671:
						if !_rules[rulesp]() {
							goto l663
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l663
						}
						if !_rules[rulesp]() {
							goto l663
						}
					l673:
						{
							position674, tokenIndex674, depth674 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l674
							}
							position++
							if !_rules[rulesp]() {
								goto l674
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l674
							}
							goto l673
						l674:
							position, tokenIndex, depth = position674, tokenIndex674, depth674
						}
						goto l664
					l663:
						position, tokenIndex, depth = position663, tokenIndex663, depth663
					}
				l664:
					depth--
					add(rulePegText, position662)
				}
				if !_rules[ruleAction32]() {
					goto l660
				}
				depth--
				add(ruleSourceSinkSpecs, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 44 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action33)> */
		func() bool {
			position675, tokenIndex675, depth675 := position, tokenIndex, depth
			{
				position676 := position
				depth++
				{
					position677 := position
					depth++
					{
						position678, tokenIndex678, depth678 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l679
						}
						position++
						goto l678
					l679:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
						if buffer[position] != rune('S') {
							goto l675
						}
						position++
					}
				l678:
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l681
						}
						position++
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if buffer[position] != rune('E') {
							goto l675
						}
						position++
					}
				l680:
					{
						position682, tokenIndex682, depth682 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l683
						}
						position++
						goto l682
					l683:
						position, tokenIndex, depth = position682, tokenIndex682, depth682
						if buffer[position] != rune('T') {
							goto l675
						}
						position++
					}
				l682:
					if !_rules[rulesp]() {
						goto l675
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l675
					}
					if !_rules[rulesp]() {
						goto l675
					}
				l684:
					{
						position685, tokenIndex685, depth685 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l685
						}
						position++
						if !_rules[rulesp]() {
							goto l685
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l685
						}
						goto l684
					l685:
						position, tokenIndex, depth = position685, tokenIndex685, depth685
					}
					depth--
					add(rulePegText, position677)
				}
				if !_rules[ruleAction33]() {
					goto l675
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 45 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action34)> */
		func() bool {
			position686, tokenIndex686, depth686 := position, tokenIndex, depth
			{
				position687 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l686
				}
				if buffer[position] != rune('=') {
					goto l686
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l686
				}
				if !_rules[ruleAction34]() {
					goto l686
				}
				depth--
				add(ruleSourceSinkParam, position687)
			}
			return true
		l686:
			position, tokenIndex, depth = position686, tokenIndex686, depth686
			return false
		},
		/* 46 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position688, tokenIndex688, depth688 := position, tokenIndex, depth
			{
				position689 := position
				depth++
				{
					position690, tokenIndex690, depth690 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l691
					}
					goto l690
				l691:
					position, tokenIndex, depth = position690, tokenIndex690, depth690
					if !_rules[ruleLiteral]() {
						goto l688
					}
				}
			l690:
				depth--
				add(ruleSourceSinkParamVal, position689)
			}
			return true
		l688:
			position, tokenIndex, depth = position688, tokenIndex688, depth688
			return false
		},
		/* 47 PausedOpt <- <(<(Paused / Unpaused)?> Action35)> */
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
						{
							position697, tokenIndex697, depth697 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l698
							}
							goto l697
						l698:
							position, tokenIndex, depth = position697, tokenIndex697, depth697
							if !_rules[ruleUnpaused]() {
								goto l695
							}
						}
					l697:
						goto l696
					l695:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
					}
				l696:
					depth--
					add(rulePegText, position694)
				}
				if !_rules[ruleAction35]() {
					goto l692
				}
				depth--
				add(rulePausedOpt, position693)
			}
			return true
		l692:
			position, tokenIndex, depth = position692, tokenIndex692, depth692
			return false
		},
		/* 48 Expression <- <orExpr> */
		func() bool {
			position699, tokenIndex699, depth699 := position, tokenIndex, depth
			{
				position700 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l699
				}
				depth--
				add(ruleExpression, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 49 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action36)> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				{
					position703 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l701
					}
					if !_rules[rulesp]() {
						goto l701
					}
					{
						position704, tokenIndex704, depth704 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l704
						}
						if !_rules[rulesp]() {
							goto l704
						}
						if !_rules[ruleandExpr]() {
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
				if !_rules[ruleAction36]() {
					goto l701
				}
				depth--
				add(ruleorExpr, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 50 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action37)> */
		func() bool {
			position706, tokenIndex706, depth706 := position, tokenIndex, depth
			{
				position707 := position
				depth++
				{
					position708 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l706
					}
					if !_rules[rulesp]() {
						goto l706
					}
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l709
						}
						if !_rules[rulesp]() {
							goto l709
						}
						if !_rules[rulenotExpr]() {
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
				if !_rules[ruleAction37]() {
					goto l706
				}
				depth--
				add(ruleandExpr, position707)
			}
			return true
		l706:
			position, tokenIndex, depth = position706, tokenIndex706, depth706
			return false
		},
		/* 51 notExpr <- <(<((Not sp)? comparisonExpr)> Action38)> */
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
						if !_rules[ruleNot]() {
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
					if !_rules[rulecomparisonExpr]() {
						goto l711
					}
					depth--
					add(rulePegText, position713)
				}
				if !_rules[ruleAction38]() {
					goto l711
				}
				depth--
				add(rulenotExpr, position712)
			}
			return true
		l711:
			position, tokenIndex, depth = position711, tokenIndex711, depth711
			return false
		},
		/* 52 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action39)> */
		func() bool {
			position716, tokenIndex716, depth716 := position, tokenIndex, depth
			{
				position717 := position
				depth++
				{
					position718 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l716
					}
					if !_rules[rulesp]() {
						goto l716
					}
					{
						position719, tokenIndex719, depth719 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l719
						}
						if !_rules[rulesp]() {
							goto l719
						}
						if !_rules[ruleotherOpExpr]() {
							goto l719
						}
						goto l720
					l719:
						position, tokenIndex, depth = position719, tokenIndex719, depth719
					}
				l720:
					depth--
					add(rulePegText, position718)
				}
				if !_rules[ruleAction39]() {
					goto l716
				}
				depth--
				add(rulecomparisonExpr, position717)
			}
			return true
		l716:
			position, tokenIndex, depth = position716, tokenIndex716, depth716
			return false
		},
		/* 53 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action40)> */
		func() bool {
			position721, tokenIndex721, depth721 := position, tokenIndex, depth
			{
				position722 := position
				depth++
				{
					position723 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l721
					}
					if !_rules[rulesp]() {
						goto l721
					}
				l724:
					{
						position725, tokenIndex725, depth725 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l725
						}
						if !_rules[rulesp]() {
							goto l725
						}
						if !_rules[ruleisExpr]() {
							goto l725
						}
						if !_rules[rulesp]() {
							goto l725
						}
						goto l724
					l725:
						position, tokenIndex, depth = position725, tokenIndex725, depth725
					}
					depth--
					add(rulePegText, position723)
				}
				if !_rules[ruleAction40]() {
					goto l721
				}
				depth--
				add(ruleotherOpExpr, position722)
			}
			return true
		l721:
			position, tokenIndex, depth = position721, tokenIndex721, depth721
			return false
		},
		/* 54 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action41)> */
		func() bool {
			position726, tokenIndex726, depth726 := position, tokenIndex, depth
			{
				position727 := position
				depth++
				{
					position728 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l726
					}
					if !_rules[rulesp]() {
						goto l726
					}
					{
						position729, tokenIndex729, depth729 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l729
						}
						if !_rules[rulesp]() {
							goto l729
						}
						if !_rules[ruleNullLiteral]() {
							goto l729
						}
						goto l730
					l729:
						position, tokenIndex, depth = position729, tokenIndex729, depth729
					}
				l730:
					depth--
					add(rulePegText, position728)
				}
				if !_rules[ruleAction41]() {
					goto l726
				}
				depth--
				add(ruleisExpr, position727)
			}
			return true
		l726:
			position, tokenIndex, depth = position726, tokenIndex726, depth726
			return false
		},
		/* 55 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action42)> */
		func() bool {
			position731, tokenIndex731, depth731 := position, tokenIndex, depth
			{
				position732 := position
				depth++
				{
					position733 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l731
					}
					if !_rules[rulesp]() {
						goto l731
					}
				l734:
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l735
						}
						if !_rules[rulesp]() {
							goto l735
						}
						if !_rules[ruleproductExpr]() {
							goto l735
						}
						if !_rules[rulesp]() {
							goto l735
						}
						goto l734
					l735:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
					}
					depth--
					add(rulePegText, position733)
				}
				if !_rules[ruleAction42]() {
					goto l731
				}
				depth--
				add(ruletermExpr, position732)
			}
			return true
		l731:
			position, tokenIndex, depth = position731, tokenIndex731, depth731
			return false
		},
		/* 56 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action43)> */
		func() bool {
			position736, tokenIndex736, depth736 := position, tokenIndex, depth
			{
				position737 := position
				depth++
				{
					position738 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l736
					}
					if !_rules[rulesp]() {
						goto l736
					}
				l739:
					{
						position740, tokenIndex740, depth740 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l740
						}
						if !_rules[rulesp]() {
							goto l740
						}
						if !_rules[ruleminusExpr]() {
							goto l740
						}
						if !_rules[rulesp]() {
							goto l740
						}
						goto l739
					l740:
						position, tokenIndex, depth = position740, tokenIndex740, depth740
					}
					depth--
					add(rulePegText, position738)
				}
				if !_rules[ruleAction43]() {
					goto l736
				}
				depth--
				add(ruleproductExpr, position737)
			}
			return true
		l736:
			position, tokenIndex, depth = position736, tokenIndex736, depth736
			return false
		},
		/* 57 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action44)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l744
						}
						if !_rules[rulesp]() {
							goto l744
						}
						goto l745
					l744:
						position, tokenIndex, depth = position744, tokenIndex744, depth744
					}
				l745:
					if !_rules[rulecastExpr]() {
						goto l741
					}
					depth--
					add(rulePegText, position743)
				}
				if !_rules[ruleAction44]() {
					goto l741
				}
				depth--
				add(ruleminusExpr, position742)
			}
			return true
		l741:
			position, tokenIndex, depth = position741, tokenIndex741, depth741
			return false
		},
		/* 58 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action45)> */
		func() bool {
			position746, tokenIndex746, depth746 := position, tokenIndex, depth
			{
				position747 := position
				depth++
				{
					position748 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l746
					}
					{
						position749, tokenIndex749, depth749 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l749
						}
						if buffer[position] != rune(':') {
							goto l749
						}
						position++
						if buffer[position] != rune(':') {
							goto l749
						}
						position++
						if !_rules[rulesp]() {
							goto l749
						}
						if !_rules[ruleType]() {
							goto l749
						}
						goto l750
					l749:
						position, tokenIndex, depth = position749, tokenIndex749, depth749
					}
				l750:
					depth--
					add(rulePegText, position748)
				}
				if !_rules[ruleAction45]() {
					goto l746
				}
				depth--
				add(rulecastExpr, position747)
			}
			return true
		l746:
			position, tokenIndex, depth = position746, tokenIndex746, depth746
			return false
		},
		/* 59 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / Literal)> */
		func() bool {
			position751, tokenIndex751, depth751 := position, tokenIndex, depth
			{
				position752 := position
				depth++
				{
					position753, tokenIndex753, depth753 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l754
					}
					position++
					if !_rules[rulesp]() {
						goto l754
					}
					if !_rules[ruleExpression]() {
						goto l754
					}
					if !_rules[rulesp]() {
						goto l754
					}
					if buffer[position] != rune(')') {
						goto l754
					}
					position++
					goto l753
				l754:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if !_rules[ruleBooleanLiteral]() {
						goto l755
					}
					goto l753
				l755:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if !_rules[ruleNullLiteral]() {
						goto l756
					}
					goto l753
				l756:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if !_rules[ruleRowMeta]() {
						goto l757
					}
					goto l753
				l757:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if !_rules[ruleFuncTypeCast]() {
						goto l758
					}
					goto l753
				l758:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if !_rules[ruleFuncApp]() {
						goto l759
					}
					goto l753
				l759:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if !_rules[ruleRowValue]() {
						goto l760
					}
					goto l753
				l760:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if !_rules[ruleLiteral]() {
						goto l751
					}
				}
			l753:
				depth--
				add(rulebaseExpr, position752)
			}
			return true
		l751:
			position, tokenIndex, depth = position751, tokenIndex751, depth751
			return false
		},
		/* 60 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action46)> */
		func() bool {
			position761, tokenIndex761, depth761 := position, tokenIndex, depth
			{
				position762 := position
				depth++
				{
					position763 := position
					depth++
					{
						position764, tokenIndex764, depth764 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l765
						}
						position++
						goto l764
					l765:
						position, tokenIndex, depth = position764, tokenIndex764, depth764
						if buffer[position] != rune('C') {
							goto l761
						}
						position++
					}
				l764:
					{
						position766, tokenIndex766, depth766 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l767
						}
						position++
						goto l766
					l767:
						position, tokenIndex, depth = position766, tokenIndex766, depth766
						if buffer[position] != rune('A') {
							goto l761
						}
						position++
					}
				l766:
					{
						position768, tokenIndex768, depth768 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l769
						}
						position++
						goto l768
					l769:
						position, tokenIndex, depth = position768, tokenIndex768, depth768
						if buffer[position] != rune('S') {
							goto l761
						}
						position++
					}
				l768:
					{
						position770, tokenIndex770, depth770 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l771
						}
						position++
						goto l770
					l771:
						position, tokenIndex, depth = position770, tokenIndex770, depth770
						if buffer[position] != rune('T') {
							goto l761
						}
						position++
					}
				l770:
					if !_rules[rulesp]() {
						goto l761
					}
					if buffer[position] != rune('(') {
						goto l761
					}
					position++
					if !_rules[rulesp]() {
						goto l761
					}
					if !_rules[ruleExpression]() {
						goto l761
					}
					if !_rules[rulesp]() {
						goto l761
					}
					{
						position772, tokenIndex772, depth772 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l773
						}
						position++
						goto l772
					l773:
						position, tokenIndex, depth = position772, tokenIndex772, depth772
						if buffer[position] != rune('A') {
							goto l761
						}
						position++
					}
				l772:
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
							goto l761
						}
						position++
					}
				l774:
					if !_rules[rulesp]() {
						goto l761
					}
					if !_rules[ruleType]() {
						goto l761
					}
					if !_rules[rulesp]() {
						goto l761
					}
					if buffer[position] != rune(')') {
						goto l761
					}
					position++
					depth--
					add(rulePegText, position763)
				}
				if !_rules[ruleAction46]() {
					goto l761
				}
				depth--
				add(ruleFuncTypeCast, position762)
			}
			return true
		l761:
			position, tokenIndex, depth = position761, tokenIndex761, depth761
			return false
		},
		/* 61 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action47)> */
		func() bool {
			position776, tokenIndex776, depth776 := position, tokenIndex, depth
			{
				position777 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l776
				}
				if !_rules[rulesp]() {
					goto l776
				}
				if buffer[position] != rune('(') {
					goto l776
				}
				position++
				if !_rules[rulesp]() {
					goto l776
				}
				if !_rules[ruleFuncParams]() {
					goto l776
				}
				if !_rules[rulesp]() {
					goto l776
				}
				if buffer[position] != rune(')') {
					goto l776
				}
				position++
				if !_rules[ruleAction47]() {
					goto l776
				}
				depth--
				add(ruleFuncApp, position777)
			}
			return true
		l776:
			position, tokenIndex, depth = position776, tokenIndex776, depth776
			return false
		},
		/* 62 FuncParams <- <(<(Wildcard / (Expression sp (',' sp Expression)*)?)> Action48)> */
		func() bool {
			position778, tokenIndex778, depth778 := position, tokenIndex, depth
			{
				position779 := position
				depth++
				{
					position780 := position
					depth++
					{
						position781, tokenIndex781, depth781 := position, tokenIndex, depth
						if !_rules[ruleWildcard]() {
							goto l782
						}
						goto l781
					l782:
						position, tokenIndex, depth = position781, tokenIndex781, depth781
						{
							position783, tokenIndex783, depth783 := position, tokenIndex, depth
							if !_rules[ruleExpression]() {
								goto l783
							}
							if !_rules[rulesp]() {
								goto l783
							}
						l785:
							{
								position786, tokenIndex786, depth786 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l786
								}
								position++
								if !_rules[rulesp]() {
									goto l786
								}
								if !_rules[ruleExpression]() {
									goto l786
								}
								goto l785
							l786:
								position, tokenIndex, depth = position786, tokenIndex786, depth786
							}
							goto l784
						l783:
							position, tokenIndex, depth = position783, tokenIndex783, depth783
						}
					l784:
					}
				l781:
					depth--
					add(rulePegText, position780)
				}
				if !_rules[ruleAction48]() {
					goto l778
				}
				depth--
				add(ruleFuncParams, position779)
			}
			return true
		l778:
			position, tokenIndex, depth = position778, tokenIndex778, depth778
			return false
		},
		/* 63 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position787, tokenIndex787, depth787 := position, tokenIndex, depth
			{
				position788 := position
				depth++
				{
					position789, tokenIndex789, depth789 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l790
					}
					goto l789
				l790:
					position, tokenIndex, depth = position789, tokenIndex789, depth789
					if !_rules[ruleNumericLiteral]() {
						goto l791
					}
					goto l789
				l791:
					position, tokenIndex, depth = position789, tokenIndex789, depth789
					if !_rules[ruleStringLiteral]() {
						goto l787
					}
				}
			l789:
				depth--
				add(ruleLiteral, position788)
			}
			return true
		l787:
			position, tokenIndex, depth = position787, tokenIndex787, depth787
			return false
		},
		/* 64 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position792, tokenIndex792, depth792 := position, tokenIndex, depth
			{
				position793 := position
				depth++
				{
					position794, tokenIndex794, depth794 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l795
					}
					goto l794
				l795:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if !_rules[ruleNotEqual]() {
						goto l796
					}
					goto l794
				l796:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if !_rules[ruleLessOrEqual]() {
						goto l797
					}
					goto l794
				l797:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if !_rules[ruleLess]() {
						goto l798
					}
					goto l794
				l798:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if !_rules[ruleGreaterOrEqual]() {
						goto l799
					}
					goto l794
				l799:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if !_rules[ruleGreater]() {
						goto l800
					}
					goto l794
				l800:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if !_rules[ruleNotEqual]() {
						goto l792
					}
				}
			l794:
				depth--
				add(ruleComparisonOp, position793)
			}
			return true
		l792:
			position, tokenIndex, depth = position792, tokenIndex792, depth792
			return false
		},
		/* 65 OtherOp <- <Concat> */
		func() bool {
			position801, tokenIndex801, depth801 := position, tokenIndex, depth
			{
				position802 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l801
				}
				depth--
				add(ruleOtherOp, position802)
			}
			return true
		l801:
			position, tokenIndex, depth = position801, tokenIndex801, depth801
			return false
		},
		/* 66 IsOp <- <(IsNot / Is)> */
		func() bool {
			position803, tokenIndex803, depth803 := position, tokenIndex, depth
			{
				position804 := position
				depth++
				{
					position805, tokenIndex805, depth805 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l806
					}
					goto l805
				l806:
					position, tokenIndex, depth = position805, tokenIndex805, depth805
					if !_rules[ruleIs]() {
						goto l803
					}
				}
			l805:
				depth--
				add(ruleIsOp, position804)
			}
			return true
		l803:
			position, tokenIndex, depth = position803, tokenIndex803, depth803
			return false
		},
		/* 67 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position807, tokenIndex807, depth807 := position, tokenIndex, depth
			{
				position808 := position
				depth++
				{
					position809, tokenIndex809, depth809 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l810
					}
					goto l809
				l810:
					position, tokenIndex, depth = position809, tokenIndex809, depth809
					if !_rules[ruleMinus]() {
						goto l807
					}
				}
			l809:
				depth--
				add(rulePlusMinusOp, position808)
			}
			return true
		l807:
			position, tokenIndex, depth = position807, tokenIndex807, depth807
			return false
		},
		/* 68 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position811, tokenIndex811, depth811 := position, tokenIndex, depth
			{
				position812 := position
				depth++
				{
					position813, tokenIndex813, depth813 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l814
					}
					goto l813
				l814:
					position, tokenIndex, depth = position813, tokenIndex813, depth813
					if !_rules[ruleDivide]() {
						goto l815
					}
					goto l813
				l815:
					position, tokenIndex, depth = position813, tokenIndex813, depth813
					if !_rules[ruleModulo]() {
						goto l811
					}
				}
			l813:
				depth--
				add(ruleMultDivOp, position812)
			}
			return true
		l811:
			position, tokenIndex, depth = position811, tokenIndex811, depth811
			return false
		},
		/* 69 Stream <- <(<ident> Action49)> */
		func() bool {
			position816, tokenIndex816, depth816 := position, tokenIndex, depth
			{
				position817 := position
				depth++
				{
					position818 := position
					depth++
					if !_rules[ruleident]() {
						goto l816
					}
					depth--
					add(rulePegText, position818)
				}
				if !_rules[ruleAction49]() {
					goto l816
				}
				depth--
				add(ruleStream, position817)
			}
			return true
		l816:
			position, tokenIndex, depth = position816, tokenIndex816, depth816
			return false
		},
		/* 70 RowMeta <- <RowTimestamp> */
		func() bool {
			position819, tokenIndex819, depth819 := position, tokenIndex, depth
			{
				position820 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l819
				}
				depth--
				add(ruleRowMeta, position820)
			}
			return true
		l819:
			position, tokenIndex, depth = position819, tokenIndex819, depth819
			return false
		},
		/* 71 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action50)> */
		func() bool {
			position821, tokenIndex821, depth821 := position, tokenIndex, depth
			{
				position822 := position
				depth++
				{
					position823 := position
					depth++
					{
						position824, tokenIndex824, depth824 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l824
						}
						if buffer[position] != rune(':') {
							goto l824
						}
						position++
						goto l825
					l824:
						position, tokenIndex, depth = position824, tokenIndex824, depth824
					}
				l825:
					if buffer[position] != rune('t') {
						goto l821
					}
					position++
					if buffer[position] != rune('s') {
						goto l821
					}
					position++
					if buffer[position] != rune('(') {
						goto l821
					}
					position++
					if buffer[position] != rune(')') {
						goto l821
					}
					position++
					depth--
					add(rulePegText, position823)
				}
				if !_rules[ruleAction50]() {
					goto l821
				}
				depth--
				add(ruleRowTimestamp, position822)
			}
			return true
		l821:
			position, tokenIndex, depth = position821, tokenIndex821, depth821
			return false
		},
		/* 72 RowValue <- <(<((ident ':' !':')? jsonPath)> Action51)> */
		func() bool {
			position826, tokenIndex826, depth826 := position, tokenIndex, depth
			{
				position827 := position
				depth++
				{
					position828 := position
					depth++
					{
						position829, tokenIndex829, depth829 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l829
						}
						if buffer[position] != rune(':') {
							goto l829
						}
						position++
						{
							position831, tokenIndex831, depth831 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l831
							}
							position++
							goto l829
						l831:
							position, tokenIndex, depth = position831, tokenIndex831, depth831
						}
						goto l830
					l829:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
					}
				l830:
					if !_rules[rulejsonPath]() {
						goto l826
					}
					depth--
					add(rulePegText, position828)
				}
				if !_rules[ruleAction51]() {
					goto l826
				}
				depth--
				add(ruleRowValue, position827)
			}
			return true
		l826:
			position, tokenIndex, depth = position826, tokenIndex826, depth826
			return false
		},
		/* 73 NumericLiteral <- <(<('-'? [0-9]+)> Action52)> */
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
						if buffer[position] != rune('-') {
							goto l835
						}
						position++
						goto l836
					l835:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
					}
				l836:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l832
					}
					position++
				l837:
					{
						position838, tokenIndex838, depth838 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l838
						}
						position++
						goto l837
					l838:
						position, tokenIndex, depth = position838, tokenIndex838, depth838
					}
					depth--
					add(rulePegText, position834)
				}
				if !_rules[ruleAction52]() {
					goto l832
				}
				depth--
				add(ruleNumericLiteral, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 74 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action53)> */
		func() bool {
			position839, tokenIndex839, depth839 := position, tokenIndex, depth
			{
				position840 := position
				depth++
				{
					position841 := position
					depth++
					{
						position842, tokenIndex842, depth842 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l842
						}
						position++
						goto l843
					l842:
						position, tokenIndex, depth = position842, tokenIndex842, depth842
					}
				l843:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l839
					}
					position++
				l844:
					{
						position845, tokenIndex845, depth845 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l845
						}
						position++
						goto l844
					l845:
						position, tokenIndex, depth = position845, tokenIndex845, depth845
					}
					if buffer[position] != rune('.') {
						goto l839
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l839
					}
					position++
				l846:
					{
						position847, tokenIndex847, depth847 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l847
						}
						position++
						goto l846
					l847:
						position, tokenIndex, depth = position847, tokenIndex847, depth847
					}
					depth--
					add(rulePegText, position841)
				}
				if !_rules[ruleAction53]() {
					goto l839
				}
				depth--
				add(ruleFloatLiteral, position840)
			}
			return true
		l839:
			position, tokenIndex, depth = position839, tokenIndex839, depth839
			return false
		},
		/* 75 Function <- <(<ident> Action54)> */
		func() bool {
			position848, tokenIndex848, depth848 := position, tokenIndex, depth
			{
				position849 := position
				depth++
				{
					position850 := position
					depth++
					if !_rules[ruleident]() {
						goto l848
					}
					depth--
					add(rulePegText, position850)
				}
				if !_rules[ruleAction54]() {
					goto l848
				}
				depth--
				add(ruleFunction, position849)
			}
			return true
		l848:
			position, tokenIndex, depth = position848, tokenIndex848, depth848
			return false
		},
		/* 76 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action55)> */
		func() bool {
			position851, tokenIndex851, depth851 := position, tokenIndex, depth
			{
				position852 := position
				depth++
				{
					position853 := position
					depth++
					{
						position854, tokenIndex854, depth854 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l855
						}
						position++
						goto l854
					l855:
						position, tokenIndex, depth = position854, tokenIndex854, depth854
						if buffer[position] != rune('N') {
							goto l851
						}
						position++
					}
				l854:
					{
						position856, tokenIndex856, depth856 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l857
						}
						position++
						goto l856
					l857:
						position, tokenIndex, depth = position856, tokenIndex856, depth856
						if buffer[position] != rune('U') {
							goto l851
						}
						position++
					}
				l856:
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l859
						}
						position++
						goto l858
					l859:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if buffer[position] != rune('L') {
							goto l851
						}
						position++
					}
				l858:
					{
						position860, tokenIndex860, depth860 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l861
						}
						position++
						goto l860
					l861:
						position, tokenIndex, depth = position860, tokenIndex860, depth860
						if buffer[position] != rune('L') {
							goto l851
						}
						position++
					}
				l860:
					depth--
					add(rulePegText, position853)
				}
				if !_rules[ruleAction55]() {
					goto l851
				}
				depth--
				add(ruleNullLiteral, position852)
			}
			return true
		l851:
			position, tokenIndex, depth = position851, tokenIndex851, depth851
			return false
		},
		/* 77 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position862, tokenIndex862, depth862 := position, tokenIndex, depth
			{
				position863 := position
				depth++
				{
					position864, tokenIndex864, depth864 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l865
					}
					goto l864
				l865:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleFALSE]() {
						goto l862
					}
				}
			l864:
				depth--
				add(ruleBooleanLiteral, position863)
			}
			return true
		l862:
			position, tokenIndex, depth = position862, tokenIndex862, depth862
			return false
		},
		/* 78 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action56)> */
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
						if buffer[position] != rune('t') {
							goto l870
						}
						position++
						goto l869
					l870:
						position, tokenIndex, depth = position869, tokenIndex869, depth869
						if buffer[position] != rune('T') {
							goto l866
						}
						position++
					}
				l869:
					{
						position871, tokenIndex871, depth871 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l872
						}
						position++
						goto l871
					l872:
						position, tokenIndex, depth = position871, tokenIndex871, depth871
						if buffer[position] != rune('R') {
							goto l866
						}
						position++
					}
				l871:
					{
						position873, tokenIndex873, depth873 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l874
						}
						position++
						goto l873
					l874:
						position, tokenIndex, depth = position873, tokenIndex873, depth873
						if buffer[position] != rune('U') {
							goto l866
						}
						position++
					}
				l873:
					{
						position875, tokenIndex875, depth875 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l876
						}
						position++
						goto l875
					l876:
						position, tokenIndex, depth = position875, tokenIndex875, depth875
						if buffer[position] != rune('E') {
							goto l866
						}
						position++
					}
				l875:
					depth--
					add(rulePegText, position868)
				}
				if !_rules[ruleAction56]() {
					goto l866
				}
				depth--
				add(ruleTRUE, position867)
			}
			return true
		l866:
			position, tokenIndex, depth = position866, tokenIndex866, depth866
			return false
		},
		/* 79 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action57)> */
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
						if buffer[position] != rune('f') {
							goto l881
						}
						position++
						goto l880
					l881:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
						if buffer[position] != rune('F') {
							goto l877
						}
						position++
					}
				l880:
					{
						position882, tokenIndex882, depth882 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l883
						}
						position++
						goto l882
					l883:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
						if buffer[position] != rune('A') {
							goto l877
						}
						position++
					}
				l882:
					{
						position884, tokenIndex884, depth884 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l885
						}
						position++
						goto l884
					l885:
						position, tokenIndex, depth = position884, tokenIndex884, depth884
						if buffer[position] != rune('L') {
							goto l877
						}
						position++
					}
				l884:
					{
						position886, tokenIndex886, depth886 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l887
						}
						position++
						goto l886
					l887:
						position, tokenIndex, depth = position886, tokenIndex886, depth886
						if buffer[position] != rune('S') {
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
					depth--
					add(rulePegText, position879)
				}
				if !_rules[ruleAction57]() {
					goto l877
				}
				depth--
				add(ruleFALSE, position878)
			}
			return true
		l877:
			position, tokenIndex, depth = position877, tokenIndex877, depth877
			return false
		},
		/* 80 Wildcard <- <(<'*'> Action58)> */
		func() bool {
			position890, tokenIndex890, depth890 := position, tokenIndex, depth
			{
				position891 := position
				depth++
				{
					position892 := position
					depth++
					if buffer[position] != rune('*') {
						goto l890
					}
					position++
					depth--
					add(rulePegText, position892)
				}
				if !_rules[ruleAction58]() {
					goto l890
				}
				depth--
				add(ruleWildcard, position891)
			}
			return true
		l890:
			position, tokenIndex, depth = position890, tokenIndex890, depth890
			return false
		},
		/* 81 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action59)> */
		func() bool {
			position893, tokenIndex893, depth893 := position, tokenIndex, depth
			{
				position894 := position
				depth++
				{
					position895 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l893
					}
					position++
				l896:
					{
						position897, tokenIndex897, depth897 := position, tokenIndex, depth
						{
							position898, tokenIndex898, depth898 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l899
							}
							position++
							if buffer[position] != rune('\'') {
								goto l899
							}
							position++
							goto l898
						l899:
							position, tokenIndex, depth = position898, tokenIndex898, depth898
							{
								position900, tokenIndex900, depth900 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l900
								}
								position++
								goto l897
							l900:
								position, tokenIndex, depth = position900, tokenIndex900, depth900
							}
							if !matchDot() {
								goto l897
							}
						}
					l898:
						goto l896
					l897:
						position, tokenIndex, depth = position897, tokenIndex897, depth897
					}
					if buffer[position] != rune('\'') {
						goto l893
					}
					position++
					depth--
					add(rulePegText, position895)
				}
				if !_rules[ruleAction59]() {
					goto l893
				}
				depth--
				add(ruleStringLiteral, position894)
			}
			return true
		l893:
			position, tokenIndex, depth = position893, tokenIndex893, depth893
			return false
		},
		/* 82 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action60)> */
		func() bool {
			position901, tokenIndex901, depth901 := position, tokenIndex, depth
			{
				position902 := position
				depth++
				{
					position903 := position
					depth++
					{
						position904, tokenIndex904, depth904 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l905
						}
						position++
						goto l904
					l905:
						position, tokenIndex, depth = position904, tokenIndex904, depth904
						if buffer[position] != rune('I') {
							goto l901
						}
						position++
					}
				l904:
					{
						position906, tokenIndex906, depth906 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l907
						}
						position++
						goto l906
					l907:
						position, tokenIndex, depth = position906, tokenIndex906, depth906
						if buffer[position] != rune('S') {
							goto l901
						}
						position++
					}
				l906:
					{
						position908, tokenIndex908, depth908 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l909
						}
						position++
						goto l908
					l909:
						position, tokenIndex, depth = position908, tokenIndex908, depth908
						if buffer[position] != rune('T') {
							goto l901
						}
						position++
					}
				l908:
					{
						position910, tokenIndex910, depth910 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l911
						}
						position++
						goto l910
					l911:
						position, tokenIndex, depth = position910, tokenIndex910, depth910
						if buffer[position] != rune('R') {
							goto l901
						}
						position++
					}
				l910:
					{
						position912, tokenIndex912, depth912 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l913
						}
						position++
						goto l912
					l913:
						position, tokenIndex, depth = position912, tokenIndex912, depth912
						if buffer[position] != rune('E') {
							goto l901
						}
						position++
					}
				l912:
					{
						position914, tokenIndex914, depth914 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l915
						}
						position++
						goto l914
					l915:
						position, tokenIndex, depth = position914, tokenIndex914, depth914
						if buffer[position] != rune('A') {
							goto l901
						}
						position++
					}
				l914:
					{
						position916, tokenIndex916, depth916 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l917
						}
						position++
						goto l916
					l917:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						if buffer[position] != rune('M') {
							goto l901
						}
						position++
					}
				l916:
					depth--
					add(rulePegText, position903)
				}
				if !_rules[ruleAction60]() {
					goto l901
				}
				depth--
				add(ruleISTREAM, position902)
			}
			return true
		l901:
			position, tokenIndex, depth = position901, tokenIndex901, depth901
			return false
		},
		/* 83 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action61)> */
		func() bool {
			position918, tokenIndex918, depth918 := position, tokenIndex, depth
			{
				position919 := position
				depth++
				{
					position920 := position
					depth++
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
							goto l918
						}
						position++
					}
				l921:
					{
						position923, tokenIndex923, depth923 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l924
						}
						position++
						goto l923
					l924:
						position, tokenIndex, depth = position923, tokenIndex923, depth923
						if buffer[position] != rune('S') {
							goto l918
						}
						position++
					}
				l923:
					{
						position925, tokenIndex925, depth925 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l926
						}
						position++
						goto l925
					l926:
						position, tokenIndex, depth = position925, tokenIndex925, depth925
						if buffer[position] != rune('T') {
							goto l918
						}
						position++
					}
				l925:
					{
						position927, tokenIndex927, depth927 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l928
						}
						position++
						goto l927
					l928:
						position, tokenIndex, depth = position927, tokenIndex927, depth927
						if buffer[position] != rune('R') {
							goto l918
						}
						position++
					}
				l927:
					{
						position929, tokenIndex929, depth929 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l930
						}
						position++
						goto l929
					l930:
						position, tokenIndex, depth = position929, tokenIndex929, depth929
						if buffer[position] != rune('E') {
							goto l918
						}
						position++
					}
				l929:
					{
						position931, tokenIndex931, depth931 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l932
						}
						position++
						goto l931
					l932:
						position, tokenIndex, depth = position931, tokenIndex931, depth931
						if buffer[position] != rune('A') {
							goto l918
						}
						position++
					}
				l931:
					{
						position933, tokenIndex933, depth933 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l934
						}
						position++
						goto l933
					l934:
						position, tokenIndex, depth = position933, tokenIndex933, depth933
						if buffer[position] != rune('M') {
							goto l918
						}
						position++
					}
				l933:
					depth--
					add(rulePegText, position920)
				}
				if !_rules[ruleAction61]() {
					goto l918
				}
				depth--
				add(ruleDSTREAM, position919)
			}
			return true
		l918:
			position, tokenIndex, depth = position918, tokenIndex918, depth918
			return false
		},
		/* 84 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action62)> */
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
						if buffer[position] != rune('r') {
							goto l939
						}
						position++
						goto l938
					l939:
						position, tokenIndex, depth = position938, tokenIndex938, depth938
						if buffer[position] != rune('R') {
							goto l935
						}
						position++
					}
				l938:
					{
						position940, tokenIndex940, depth940 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l941
						}
						position++
						goto l940
					l941:
						position, tokenIndex, depth = position940, tokenIndex940, depth940
						if buffer[position] != rune('S') {
							goto l935
						}
						position++
					}
				l940:
					{
						position942, tokenIndex942, depth942 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l943
						}
						position++
						goto l942
					l943:
						position, tokenIndex, depth = position942, tokenIndex942, depth942
						if buffer[position] != rune('T') {
							goto l935
						}
						position++
					}
				l942:
					{
						position944, tokenIndex944, depth944 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l945
						}
						position++
						goto l944
					l945:
						position, tokenIndex, depth = position944, tokenIndex944, depth944
						if buffer[position] != rune('R') {
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
						if buffer[position] != rune('a') {
							goto l949
						}
						position++
						goto l948
					l949:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
						if buffer[position] != rune('A') {
							goto l935
						}
						position++
					}
				l948:
					{
						position950, tokenIndex950, depth950 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l951
						}
						position++
						goto l950
					l951:
						position, tokenIndex, depth = position950, tokenIndex950, depth950
						if buffer[position] != rune('M') {
							goto l935
						}
						position++
					}
				l950:
					depth--
					add(rulePegText, position937)
				}
				if !_rules[ruleAction62]() {
					goto l935
				}
				depth--
				add(ruleRSTREAM, position936)
			}
			return true
		l935:
			position, tokenIndex, depth = position935, tokenIndex935, depth935
			return false
		},
		/* 85 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action63)> */
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
						if buffer[position] != rune('t') {
							goto l956
						}
						position++
						goto l955
					l956:
						position, tokenIndex, depth = position955, tokenIndex955, depth955
						if buffer[position] != rune('T') {
							goto l952
						}
						position++
					}
				l955:
					{
						position957, tokenIndex957, depth957 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l958
						}
						position++
						goto l957
					l958:
						position, tokenIndex, depth = position957, tokenIndex957, depth957
						if buffer[position] != rune('U') {
							goto l952
						}
						position++
					}
				l957:
					{
						position959, tokenIndex959, depth959 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l960
						}
						position++
						goto l959
					l960:
						position, tokenIndex, depth = position959, tokenIndex959, depth959
						if buffer[position] != rune('P') {
							goto l952
						}
						position++
					}
				l959:
					{
						position961, tokenIndex961, depth961 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l962
						}
						position++
						goto l961
					l962:
						position, tokenIndex, depth = position961, tokenIndex961, depth961
						if buffer[position] != rune('L') {
							goto l952
						}
						position++
					}
				l961:
					{
						position963, tokenIndex963, depth963 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l964
						}
						position++
						goto l963
					l964:
						position, tokenIndex, depth = position963, tokenIndex963, depth963
						if buffer[position] != rune('E') {
							goto l952
						}
						position++
					}
				l963:
					{
						position965, tokenIndex965, depth965 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l966
						}
						position++
						goto l965
					l966:
						position, tokenIndex, depth = position965, tokenIndex965, depth965
						if buffer[position] != rune('S') {
							goto l952
						}
						position++
					}
				l965:
					depth--
					add(rulePegText, position954)
				}
				if !_rules[ruleAction63]() {
					goto l952
				}
				depth--
				add(ruleTUPLES, position953)
			}
			return true
		l952:
			position, tokenIndex, depth = position952, tokenIndex952, depth952
			return false
		},
		/* 86 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action64)> */
		func() bool {
			position967, tokenIndex967, depth967 := position, tokenIndex, depth
			{
				position968 := position
				depth++
				{
					position969 := position
					depth++
					{
						position970, tokenIndex970, depth970 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l971
						}
						position++
						goto l970
					l971:
						position, tokenIndex, depth = position970, tokenIndex970, depth970
						if buffer[position] != rune('S') {
							goto l967
						}
						position++
					}
				l970:
					{
						position972, tokenIndex972, depth972 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l973
						}
						position++
						goto l972
					l973:
						position, tokenIndex, depth = position972, tokenIndex972, depth972
						if buffer[position] != rune('E') {
							goto l967
						}
						position++
					}
				l972:
					{
						position974, tokenIndex974, depth974 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l975
						}
						position++
						goto l974
					l975:
						position, tokenIndex, depth = position974, tokenIndex974, depth974
						if buffer[position] != rune('C') {
							goto l967
						}
						position++
					}
				l974:
					{
						position976, tokenIndex976, depth976 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l977
						}
						position++
						goto l976
					l977:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('O') {
							goto l967
						}
						position++
					}
				l976:
					{
						position978, tokenIndex978, depth978 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l979
						}
						position++
						goto l978
					l979:
						position, tokenIndex, depth = position978, tokenIndex978, depth978
						if buffer[position] != rune('N') {
							goto l967
						}
						position++
					}
				l978:
					{
						position980, tokenIndex980, depth980 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l981
						}
						position++
						goto l980
					l981:
						position, tokenIndex, depth = position980, tokenIndex980, depth980
						if buffer[position] != rune('D') {
							goto l967
						}
						position++
					}
				l980:
					{
						position982, tokenIndex982, depth982 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l983
						}
						position++
						goto l982
					l983:
						position, tokenIndex, depth = position982, tokenIndex982, depth982
						if buffer[position] != rune('S') {
							goto l967
						}
						position++
					}
				l982:
					depth--
					add(rulePegText, position969)
				}
				if !_rules[ruleAction64]() {
					goto l967
				}
				depth--
				add(ruleSECONDS, position968)
			}
			return true
		l967:
			position, tokenIndex, depth = position967, tokenIndex967, depth967
			return false
		},
		/* 87 StreamIdentifier <- <(<ident> Action65)> */
		func() bool {
			position984, tokenIndex984, depth984 := position, tokenIndex, depth
			{
				position985 := position
				depth++
				{
					position986 := position
					depth++
					if !_rules[ruleident]() {
						goto l984
					}
					depth--
					add(rulePegText, position986)
				}
				if !_rules[ruleAction65]() {
					goto l984
				}
				depth--
				add(ruleStreamIdentifier, position985)
			}
			return true
		l984:
			position, tokenIndex, depth = position984, tokenIndex984, depth984
			return false
		},
		/* 88 SourceSinkType <- <(<ident> Action66)> */
		func() bool {
			position987, tokenIndex987, depth987 := position, tokenIndex, depth
			{
				position988 := position
				depth++
				{
					position989 := position
					depth++
					if !_rules[ruleident]() {
						goto l987
					}
					depth--
					add(rulePegText, position989)
				}
				if !_rules[ruleAction66]() {
					goto l987
				}
				depth--
				add(ruleSourceSinkType, position988)
			}
			return true
		l987:
			position, tokenIndex, depth = position987, tokenIndex987, depth987
			return false
		},
		/* 89 SourceSinkParamKey <- <(<ident> Action67)> */
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
				if !_rules[ruleAction67]() {
					goto l990
				}
				depth--
				add(ruleSourceSinkParamKey, position991)
			}
			return true
		l990:
			position, tokenIndex, depth = position990, tokenIndex990, depth990
			return false
		},
		/* 90 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action68)> */
		func() bool {
			position993, tokenIndex993, depth993 := position, tokenIndex, depth
			{
				position994 := position
				depth++
				{
					position995 := position
					depth++
					{
						position996, tokenIndex996, depth996 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l997
						}
						position++
						goto l996
					l997:
						position, tokenIndex, depth = position996, tokenIndex996, depth996
						if buffer[position] != rune('P') {
							goto l993
						}
						position++
					}
				l996:
					{
						position998, tokenIndex998, depth998 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l999
						}
						position++
						goto l998
					l999:
						position, tokenIndex, depth = position998, tokenIndex998, depth998
						if buffer[position] != rune('A') {
							goto l993
						}
						position++
					}
				l998:
					{
						position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1001
						}
						position++
						goto l1000
					l1001:
						position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
						if buffer[position] != rune('U') {
							goto l993
						}
						position++
					}
				l1000:
					{
						position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1003
						}
						position++
						goto l1002
					l1003:
						position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
						if buffer[position] != rune('S') {
							goto l993
						}
						position++
					}
				l1002:
					{
						position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1005
						}
						position++
						goto l1004
					l1005:
						position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
						if buffer[position] != rune('E') {
							goto l993
						}
						position++
					}
				l1004:
					{
						position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1007
						}
						position++
						goto l1006
					l1007:
						position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
						if buffer[position] != rune('D') {
							goto l993
						}
						position++
					}
				l1006:
					depth--
					add(rulePegText, position995)
				}
				if !_rules[ruleAction68]() {
					goto l993
				}
				depth--
				add(rulePaused, position994)
			}
			return true
		l993:
			position, tokenIndex, depth = position993, tokenIndex993, depth993
			return false
		},
		/* 91 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action69)> */
		func() bool {
			position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
			{
				position1009 := position
				depth++
				{
					position1010 := position
					depth++
					{
						position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1012
						}
						position++
						goto l1011
					l1012:
						position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
						if buffer[position] != rune('U') {
							goto l1008
						}
						position++
					}
				l1011:
					{
						position1013, tokenIndex1013, depth1013 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1014
						}
						position++
						goto l1013
					l1014:
						position, tokenIndex, depth = position1013, tokenIndex1013, depth1013
						if buffer[position] != rune('N') {
							goto l1008
						}
						position++
					}
				l1013:
					{
						position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1016
						}
						position++
						goto l1015
					l1016:
						position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
						if buffer[position] != rune('P') {
							goto l1008
						}
						position++
					}
				l1015:
					{
						position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1018
						}
						position++
						goto l1017
					l1018:
						position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
						if buffer[position] != rune('A') {
							goto l1008
						}
						position++
					}
				l1017:
					{
						position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1020
						}
						position++
						goto l1019
					l1020:
						position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
						if buffer[position] != rune('U') {
							goto l1008
						}
						position++
					}
				l1019:
					{
						position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1022
						}
						position++
						goto l1021
					l1022:
						position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
						if buffer[position] != rune('S') {
							goto l1008
						}
						position++
					}
				l1021:
					{
						position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1024
						}
						position++
						goto l1023
					l1024:
						position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
						if buffer[position] != rune('E') {
							goto l1008
						}
						position++
					}
				l1023:
					{
						position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1026
						}
						position++
						goto l1025
					l1026:
						position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
						if buffer[position] != rune('D') {
							goto l1008
						}
						position++
					}
				l1025:
					depth--
					add(rulePegText, position1010)
				}
				if !_rules[ruleAction69]() {
					goto l1008
				}
				depth--
				add(ruleUnpaused, position1009)
			}
			return true
		l1008:
			position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
			return false
		},
		/* 92 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
			{
				position1028 := position
				depth++
				{
					position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1030
					}
					goto l1029
				l1030:
					position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					if !_rules[ruleInt]() {
						goto l1031
					}
					goto l1029
				l1031:
					position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					if !_rules[ruleFloat]() {
						goto l1032
					}
					goto l1029
				l1032:
					position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					if !_rules[ruleString]() {
						goto l1033
					}
					goto l1029
				l1033:
					position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					if !_rules[ruleBlob]() {
						goto l1034
					}
					goto l1029
				l1034:
					position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					if !_rules[ruleTimestamp]() {
						goto l1035
					}
					goto l1029
				l1035:
					position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					if !_rules[ruleArray]() {
						goto l1036
					}
					goto l1029
				l1036:
					position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					if !_rules[ruleMap]() {
						goto l1027
					}
				}
			l1029:
				depth--
				add(ruleType, position1028)
			}
			return true
		l1027:
			position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
			return false
		},
		/* 93 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action70)> */
		func() bool {
			position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
			{
				position1038 := position
				depth++
				{
					position1039 := position
					depth++
					{
						position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1041
						}
						position++
						goto l1040
					l1041:
						position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
						if buffer[position] != rune('B') {
							goto l1037
						}
						position++
					}
				l1040:
					{
						position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1043
						}
						position++
						goto l1042
					l1043:
						position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
						if buffer[position] != rune('O') {
							goto l1037
						}
						position++
					}
				l1042:
					{
						position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1045
						}
						position++
						goto l1044
					l1045:
						position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
						if buffer[position] != rune('O') {
							goto l1037
						}
						position++
					}
				l1044:
					{
						position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1047
						}
						position++
						goto l1046
					l1047:
						position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
						if buffer[position] != rune('L') {
							goto l1037
						}
						position++
					}
				l1046:
					depth--
					add(rulePegText, position1039)
				}
				if !_rules[ruleAction70]() {
					goto l1037
				}
				depth--
				add(ruleBool, position1038)
			}
			return true
		l1037:
			position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
			return false
		},
		/* 94 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action71)> */
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
						if buffer[position] != rune('i') {
							goto l1052
						}
						position++
						goto l1051
					l1052:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if buffer[position] != rune('I') {
							goto l1048
						}
						position++
					}
				l1051:
					{
						position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1054
						}
						position++
						goto l1053
					l1054:
						position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
						if buffer[position] != rune('N') {
							goto l1048
						}
						position++
					}
				l1053:
					{
						position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1056
						}
						position++
						goto l1055
					l1056:
						position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
						if buffer[position] != rune('T') {
							goto l1048
						}
						position++
					}
				l1055:
					depth--
					add(rulePegText, position1050)
				}
				if !_rules[ruleAction71]() {
					goto l1048
				}
				depth--
				add(ruleInt, position1049)
			}
			return true
		l1048:
			position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
			return false
		},
		/* 95 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action72)> */
		func() bool {
			position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
			{
				position1058 := position
				depth++
				{
					position1059 := position
					depth++
					{
						position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1061
						}
						position++
						goto l1060
					l1061:
						position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
						if buffer[position] != rune('F') {
							goto l1057
						}
						position++
					}
				l1060:
					{
						position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1063
						}
						position++
						goto l1062
					l1063:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune('L') {
							goto l1057
						}
						position++
					}
				l1062:
					{
						position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1065
						}
						position++
						goto l1064
					l1065:
						position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
						if buffer[position] != rune('O') {
							goto l1057
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
							goto l1057
						}
						position++
					}
				l1066:
					{
						position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1069
						}
						position++
						goto l1068
					l1069:
						position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
						if buffer[position] != rune('T') {
							goto l1057
						}
						position++
					}
				l1068:
					depth--
					add(rulePegText, position1059)
				}
				if !_rules[ruleAction72]() {
					goto l1057
				}
				depth--
				add(ruleFloat, position1058)
			}
			return true
		l1057:
			position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
			return false
		},
		/* 96 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action73)> */
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
						if buffer[position] != rune('s') {
							goto l1074
						}
						position++
						goto l1073
					l1074:
						position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
						if buffer[position] != rune('S') {
							goto l1070
						}
						position++
					}
				l1073:
					{
						position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1076
						}
						position++
						goto l1075
					l1076:
						position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
						if buffer[position] != rune('T') {
							goto l1070
						}
						position++
					}
				l1075:
					{
						position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1078
						}
						position++
						goto l1077
					l1078:
						position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
						if buffer[position] != rune('R') {
							goto l1070
						}
						position++
					}
				l1077:
					{
						position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1080
						}
						position++
						goto l1079
					l1080:
						position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
						if buffer[position] != rune('I') {
							goto l1070
						}
						position++
					}
				l1079:
					{
						position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1082
						}
						position++
						goto l1081
					l1082:
						position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
						if buffer[position] != rune('N') {
							goto l1070
						}
						position++
					}
				l1081:
					{
						position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1084
						}
						position++
						goto l1083
					l1084:
						position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
						if buffer[position] != rune('G') {
							goto l1070
						}
						position++
					}
				l1083:
					depth--
					add(rulePegText, position1072)
				}
				if !_rules[ruleAction73]() {
					goto l1070
				}
				depth--
				add(ruleString, position1071)
			}
			return true
		l1070:
			position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
			return false
		},
		/* 97 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action74)> */
		func() bool {
			position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
			{
				position1086 := position
				depth++
				{
					position1087 := position
					depth++
					{
						position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1089
						}
						position++
						goto l1088
					l1089:
						position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
						if buffer[position] != rune('B') {
							goto l1085
						}
						position++
					}
				l1088:
					{
						position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1091
						}
						position++
						goto l1090
					l1091:
						position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
						if buffer[position] != rune('L') {
							goto l1085
						}
						position++
					}
				l1090:
					{
						position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1093
						}
						position++
						goto l1092
					l1093:
						position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
						if buffer[position] != rune('O') {
							goto l1085
						}
						position++
					}
				l1092:
					{
						position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1095
						}
						position++
						goto l1094
					l1095:
						position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
						if buffer[position] != rune('B') {
							goto l1085
						}
						position++
					}
				l1094:
					depth--
					add(rulePegText, position1087)
				}
				if !_rules[ruleAction74]() {
					goto l1085
				}
				depth--
				add(ruleBlob, position1086)
			}
			return true
		l1085:
			position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
			return false
		},
		/* 98 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action75)> */
		func() bool {
			position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
			{
				position1097 := position
				depth++
				{
					position1098 := position
					depth++
					{
						position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1100
						}
						position++
						goto l1099
					l1100:
						position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
						if buffer[position] != rune('T') {
							goto l1096
						}
						position++
					}
				l1099:
					{
						position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1102
						}
						position++
						goto l1101
					l1102:
						position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
						if buffer[position] != rune('I') {
							goto l1096
						}
						position++
					}
				l1101:
					{
						position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1104
						}
						position++
						goto l1103
					l1104:
						position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
						if buffer[position] != rune('M') {
							goto l1096
						}
						position++
					}
				l1103:
					{
						position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1106
						}
						position++
						goto l1105
					l1106:
						position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
						if buffer[position] != rune('E') {
							goto l1096
						}
						position++
					}
				l1105:
					{
						position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1108
						}
						position++
						goto l1107
					l1108:
						position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
						if buffer[position] != rune('S') {
							goto l1096
						}
						position++
					}
				l1107:
					{
						position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1110
						}
						position++
						goto l1109
					l1110:
						position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
						if buffer[position] != rune('T') {
							goto l1096
						}
						position++
					}
				l1109:
					{
						position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1112
						}
						position++
						goto l1111
					l1112:
						position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
						if buffer[position] != rune('A') {
							goto l1096
						}
						position++
					}
				l1111:
					{
						position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1114
						}
						position++
						goto l1113
					l1114:
						position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
						if buffer[position] != rune('M') {
							goto l1096
						}
						position++
					}
				l1113:
					{
						position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1116
						}
						position++
						goto l1115
					l1116:
						position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
						if buffer[position] != rune('P') {
							goto l1096
						}
						position++
					}
				l1115:
					depth--
					add(rulePegText, position1098)
				}
				if !_rules[ruleAction75]() {
					goto l1096
				}
				depth--
				add(ruleTimestamp, position1097)
			}
			return true
		l1096:
			position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
			return false
		},
		/* 99 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action76)> */
		func() bool {
			position1117, tokenIndex1117, depth1117 := position, tokenIndex, depth
			{
				position1118 := position
				depth++
				{
					position1119 := position
					depth++
					{
						position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1121
						}
						position++
						goto l1120
					l1121:
						position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
						if buffer[position] != rune('A') {
							goto l1117
						}
						position++
					}
				l1120:
					{
						position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1123
						}
						position++
						goto l1122
					l1123:
						position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
						if buffer[position] != rune('R') {
							goto l1117
						}
						position++
					}
				l1122:
					{
						position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1125
						}
						position++
						goto l1124
					l1125:
						position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
						if buffer[position] != rune('R') {
							goto l1117
						}
						position++
					}
				l1124:
					{
						position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1127
						}
						position++
						goto l1126
					l1127:
						position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
						if buffer[position] != rune('A') {
							goto l1117
						}
						position++
					}
				l1126:
					{
						position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1129
						}
						position++
						goto l1128
					l1129:
						position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
						if buffer[position] != rune('Y') {
							goto l1117
						}
						position++
					}
				l1128:
					depth--
					add(rulePegText, position1119)
				}
				if !_rules[ruleAction76]() {
					goto l1117
				}
				depth--
				add(ruleArray, position1118)
			}
			return true
		l1117:
			position, tokenIndex, depth = position1117, tokenIndex1117, depth1117
			return false
		},
		/* 100 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action77)> */
		func() bool {
			position1130, tokenIndex1130, depth1130 := position, tokenIndex, depth
			{
				position1131 := position
				depth++
				{
					position1132 := position
					depth++
					{
						position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1134
						}
						position++
						goto l1133
					l1134:
						position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
						if buffer[position] != rune('M') {
							goto l1130
						}
						position++
					}
				l1133:
					{
						position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
						if buffer[position] != rune('A') {
							goto l1130
						}
						position++
					}
				l1135:
					{
						position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
						if buffer[position] != rune('P') {
							goto l1130
						}
						position++
					}
				l1137:
					depth--
					add(rulePegText, position1132)
				}
				if !_rules[ruleAction77]() {
					goto l1130
				}
				depth--
				add(ruleMap, position1131)
			}
			return true
		l1130:
			position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
			return false
		},
		/* 101 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action78)> */
		func() bool {
			position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
			{
				position1140 := position
				depth++
				{
					position1141 := position
					depth++
					{
						position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1143
						}
						position++
						goto l1142
					l1143:
						position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
						if buffer[position] != rune('O') {
							goto l1139
						}
						position++
					}
				l1142:
					{
						position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1145
						}
						position++
						goto l1144
					l1145:
						position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
						if buffer[position] != rune('R') {
							goto l1139
						}
						position++
					}
				l1144:
					depth--
					add(rulePegText, position1141)
				}
				if !_rules[ruleAction78]() {
					goto l1139
				}
				depth--
				add(ruleOr, position1140)
			}
			return true
		l1139:
			position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
			return false
		},
		/* 102 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action79)> */
		func() bool {
			position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
			{
				position1147 := position
				depth++
				{
					position1148 := position
					depth++
					{
						position1149, tokenIndex1149, depth1149 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1150
						}
						position++
						goto l1149
					l1150:
						position, tokenIndex, depth = position1149, tokenIndex1149, depth1149
						if buffer[position] != rune('A') {
							goto l1146
						}
						position++
					}
				l1149:
					{
						position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1152
						}
						position++
						goto l1151
					l1152:
						position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
						if buffer[position] != rune('N') {
							goto l1146
						}
						position++
					}
				l1151:
					{
						position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1154
						}
						position++
						goto l1153
					l1154:
						position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
						if buffer[position] != rune('D') {
							goto l1146
						}
						position++
					}
				l1153:
					depth--
					add(rulePegText, position1148)
				}
				if !_rules[ruleAction79]() {
					goto l1146
				}
				depth--
				add(ruleAnd, position1147)
			}
			return true
		l1146:
			position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
			return false
		},
		/* 103 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action80)> */
		func() bool {
			position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
			{
				position1156 := position
				depth++
				{
					position1157 := position
					depth++
					{
						position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1159
						}
						position++
						goto l1158
					l1159:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						if buffer[position] != rune('N') {
							goto l1155
						}
						position++
					}
				l1158:
					{
						position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1161
						}
						position++
						goto l1160
					l1161:
						position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
						if buffer[position] != rune('O') {
							goto l1155
						}
						position++
					}
				l1160:
					{
						position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1163
						}
						position++
						goto l1162
					l1163:
						position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
						if buffer[position] != rune('T') {
							goto l1155
						}
						position++
					}
				l1162:
					depth--
					add(rulePegText, position1157)
				}
				if !_rules[ruleAction80]() {
					goto l1155
				}
				depth--
				add(ruleNot, position1156)
			}
			return true
		l1155:
			position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
			return false
		},
		/* 104 Equal <- <(<'='> Action81)> */
		func() bool {
			position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
			{
				position1165 := position
				depth++
				{
					position1166 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1164
					}
					position++
					depth--
					add(rulePegText, position1166)
				}
				if !_rules[ruleAction81]() {
					goto l1164
				}
				depth--
				add(ruleEqual, position1165)
			}
			return true
		l1164:
			position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
			return false
		},
		/* 105 Less <- <(<'<'> Action82)> */
		func() bool {
			position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
			{
				position1168 := position
				depth++
				{
					position1169 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1167
					}
					position++
					depth--
					add(rulePegText, position1169)
				}
				if !_rules[ruleAction82]() {
					goto l1167
				}
				depth--
				add(ruleLess, position1168)
			}
			return true
		l1167:
			position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
			return false
		},
		/* 106 LessOrEqual <- <(<('<' '=')> Action83)> */
		func() bool {
			position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
			{
				position1171 := position
				depth++
				{
					position1172 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1170
					}
					position++
					if buffer[position] != rune('=') {
						goto l1170
					}
					position++
					depth--
					add(rulePegText, position1172)
				}
				if !_rules[ruleAction83]() {
					goto l1170
				}
				depth--
				add(ruleLessOrEqual, position1171)
			}
			return true
		l1170:
			position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
			return false
		},
		/* 107 Greater <- <(<'>'> Action84)> */
		func() bool {
			position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
			{
				position1174 := position
				depth++
				{
					position1175 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1173
					}
					position++
					depth--
					add(rulePegText, position1175)
				}
				if !_rules[ruleAction84]() {
					goto l1173
				}
				depth--
				add(ruleGreater, position1174)
			}
			return true
		l1173:
			position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
			return false
		},
		/* 108 GreaterOrEqual <- <(<('>' '=')> Action85)> */
		func() bool {
			position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
			{
				position1177 := position
				depth++
				{
					position1178 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1176
					}
					position++
					if buffer[position] != rune('=') {
						goto l1176
					}
					position++
					depth--
					add(rulePegText, position1178)
				}
				if !_rules[ruleAction85]() {
					goto l1176
				}
				depth--
				add(ruleGreaterOrEqual, position1177)
			}
			return true
		l1176:
			position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
			return false
		},
		/* 109 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action86)> */
		func() bool {
			position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
			{
				position1180 := position
				depth++
				{
					position1181 := position
					depth++
					{
						position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1183
						}
						position++
						if buffer[position] != rune('=') {
							goto l1183
						}
						position++
						goto l1182
					l1183:
						position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
						if buffer[position] != rune('<') {
							goto l1179
						}
						position++
						if buffer[position] != rune('>') {
							goto l1179
						}
						position++
					}
				l1182:
					depth--
					add(rulePegText, position1181)
				}
				if !_rules[ruleAction86]() {
					goto l1179
				}
				depth--
				add(ruleNotEqual, position1180)
			}
			return true
		l1179:
			position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
			return false
		},
		/* 110 Concat <- <(<('|' '|')> Action87)> */
		func() bool {
			position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
			{
				position1185 := position
				depth++
				{
					position1186 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1184
					}
					position++
					if buffer[position] != rune('|') {
						goto l1184
					}
					position++
					depth--
					add(rulePegText, position1186)
				}
				if !_rules[ruleAction87]() {
					goto l1184
				}
				depth--
				add(ruleConcat, position1185)
			}
			return true
		l1184:
			position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
			return false
		},
		/* 111 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action88)> */
		func() bool {
			position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
			{
				position1188 := position
				depth++
				{
					position1189 := position
					depth++
					{
						position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1191
						}
						position++
						goto l1190
					l1191:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('I') {
							goto l1187
						}
						position++
					}
				l1190:
					{
						position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1193
						}
						position++
						goto l1192
					l1193:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						if buffer[position] != rune('S') {
							goto l1187
						}
						position++
					}
				l1192:
					depth--
					add(rulePegText, position1189)
				}
				if !_rules[ruleAction88]() {
					goto l1187
				}
				depth--
				add(ruleIs, position1188)
			}
			return true
		l1187:
			position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
			return false
		},
		/* 112 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action89)> */
		func() bool {
			position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
			{
				position1195 := position
				depth++
				{
					position1196 := position
					depth++
					{
						position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1198
						}
						position++
						goto l1197
					l1198:
						position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
						if buffer[position] != rune('I') {
							goto l1194
						}
						position++
					}
				l1197:
					{
						position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1200
						}
						position++
						goto l1199
					l1200:
						position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
						if buffer[position] != rune('S') {
							goto l1194
						}
						position++
					}
				l1199:
					if !_rules[rulesp]() {
						goto l1194
					}
					{
						position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1202
						}
						position++
						goto l1201
					l1202:
						position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
						if buffer[position] != rune('N') {
							goto l1194
						}
						position++
					}
				l1201:
					{
						position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1204
						}
						position++
						goto l1203
					l1204:
						position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
						if buffer[position] != rune('O') {
							goto l1194
						}
						position++
					}
				l1203:
					{
						position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1206
						}
						position++
						goto l1205
					l1206:
						position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
						if buffer[position] != rune('T') {
							goto l1194
						}
						position++
					}
				l1205:
					depth--
					add(rulePegText, position1196)
				}
				if !_rules[ruleAction89]() {
					goto l1194
				}
				depth--
				add(ruleIsNot, position1195)
			}
			return true
		l1194:
			position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
			return false
		},
		/* 113 Plus <- <(<'+'> Action90)> */
		func() bool {
			position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
			{
				position1208 := position
				depth++
				{
					position1209 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1207
					}
					position++
					depth--
					add(rulePegText, position1209)
				}
				if !_rules[ruleAction90]() {
					goto l1207
				}
				depth--
				add(rulePlus, position1208)
			}
			return true
		l1207:
			position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
			return false
		},
		/* 114 Minus <- <(<'-'> Action91)> */
		func() bool {
			position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
			{
				position1211 := position
				depth++
				{
					position1212 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1210
					}
					position++
					depth--
					add(rulePegText, position1212)
				}
				if !_rules[ruleAction91]() {
					goto l1210
				}
				depth--
				add(ruleMinus, position1211)
			}
			return true
		l1210:
			position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
			return false
		},
		/* 115 Multiply <- <(<'*'> Action92)> */
		func() bool {
			position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
			{
				position1214 := position
				depth++
				{
					position1215 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1213
					}
					position++
					depth--
					add(rulePegText, position1215)
				}
				if !_rules[ruleAction92]() {
					goto l1213
				}
				depth--
				add(ruleMultiply, position1214)
			}
			return true
		l1213:
			position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
			return false
		},
		/* 116 Divide <- <(<'/'> Action93)> */
		func() bool {
			position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
			{
				position1217 := position
				depth++
				{
					position1218 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1216
					}
					position++
					depth--
					add(rulePegText, position1218)
				}
				if !_rules[ruleAction93]() {
					goto l1216
				}
				depth--
				add(ruleDivide, position1217)
			}
			return true
		l1216:
			position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
			return false
		},
		/* 117 Modulo <- <(<'%'> Action94)> */
		func() bool {
			position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
			{
				position1220 := position
				depth++
				{
					position1221 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1219
					}
					position++
					depth--
					add(rulePegText, position1221)
				}
				if !_rules[ruleAction94]() {
					goto l1219
				}
				depth--
				add(ruleModulo, position1220)
			}
			return true
		l1219:
			position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
			return false
		},
		/* 118 UnaryMinus <- <(<'-'> Action95)> */
		func() bool {
			position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
			{
				position1223 := position
				depth++
				{
					position1224 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1222
					}
					position++
					depth--
					add(rulePegText, position1224)
				}
				if !_rules[ruleAction95]() {
					goto l1222
				}
				depth--
				add(ruleUnaryMinus, position1223)
			}
			return true
		l1222:
			position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
			return false
		},
		/* 119 Identifier <- <(<ident> Action96)> */
		func() bool {
			position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
			{
				position1226 := position
				depth++
				{
					position1227 := position
					depth++
					if !_rules[ruleident]() {
						goto l1225
					}
					depth--
					add(rulePegText, position1227)
				}
				if !_rules[ruleAction96]() {
					goto l1225
				}
				depth--
				add(ruleIdentifier, position1226)
			}
			return true
		l1225:
			position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
			return false
		},
		/* 120 TargetIdentifier <- <(<jsonPath> Action97)> */
		func() bool {
			position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
			{
				position1229 := position
				depth++
				{
					position1230 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1228
					}
					depth--
					add(rulePegText, position1230)
				}
				if !_rules[ruleAction97]() {
					goto l1228
				}
				depth--
				add(ruleTargetIdentifier, position1229)
			}
			return true
		l1228:
			position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
			return false
		},
		/* 121 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
			{
				position1232 := position
				depth++
				{
					position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1234
					}
					position++
					goto l1233
				l1234:
					position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1231
					}
					position++
				}
			l1233:
			l1235:
				{
					position1236, tokenIndex1236, depth1236 := position, tokenIndex, depth
					{
						position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1238
						}
						position++
						goto l1237
					l1238:
						position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1239
						}
						position++
						goto l1237
					l1239:
						position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1240
						}
						position++
						goto l1237
					l1240:
						position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
						if buffer[position] != rune('_') {
							goto l1236
						}
						position++
					}
				l1237:
					goto l1235
				l1236:
					position, tokenIndex, depth = position1236, tokenIndex1236, depth1236
				}
				depth--
				add(ruleident, position1232)
			}
			return true
		l1231:
			position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
			return false
		},
		/* 122 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
			{
				position1242 := position
				depth++
				{
					position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1244
					}
					position++
					goto l1243
				l1244:
					position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1241
					}
					position++
				}
			l1243:
			l1245:
				{
					position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
					{
						position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1248
						}
						position++
						goto l1247
					l1248:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1249
						}
						position++
						goto l1247
					l1249:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1250
						}
						position++
						goto l1247
					l1250:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if buffer[position] != rune('_') {
							goto l1251
						}
						position++
						goto l1247
					l1251:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if buffer[position] != rune('.') {
							goto l1252
						}
						position++
						goto l1247
					l1252:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if buffer[position] != rune('[') {
							goto l1253
						}
						position++
						goto l1247
					l1253:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if buffer[position] != rune(']') {
							goto l1254
						}
						position++
						goto l1247
					l1254:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if buffer[position] != rune('"') {
							goto l1246
						}
						position++
					}
				l1247:
					goto l1245
				l1246:
					position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
				}
				depth--
				add(rulejsonPath, position1242)
			}
			return true
		l1241:
			position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
			return false
		},
		/* 123 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1256 := position
				depth++
			l1257:
				{
					position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
					{
						position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1260
						}
						position++
						goto l1259
					l1260:
						position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
						if buffer[position] != rune('\t') {
							goto l1261
						}
						position++
						goto l1259
					l1261:
						position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
						if buffer[position] != rune('\n') {
							goto l1262
						}
						position++
						goto l1259
					l1262:
						position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
						if buffer[position] != rune('\r') {
							goto l1263
						}
						position++
						goto l1259
					l1263:
						position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
						if !_rules[rulecomment]() {
							goto l1258
						}
					}
				l1259:
					goto l1257
				l1258:
					position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
				}
				depth--
				add(rulesp, position1256)
			}
			return true
		},
		/* 124 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
			{
				position1265 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1264
				}
				position++
				if buffer[position] != rune('-') {
					goto l1264
				}
				position++
			l1266:
				{
					position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
					{
						position1268, tokenIndex1268, depth1268 := position, tokenIndex, depth
						{
							position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1270
							}
							position++
							goto l1269
						l1270:
							position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
							if buffer[position] != rune('\n') {
								goto l1268
							}
							position++
						}
					l1269:
						goto l1267
					l1268:
						position, tokenIndex, depth = position1268, tokenIndex1268, depth1268
					}
					if !matchDot() {
						goto l1267
					}
					goto l1266
				l1267:
					position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
				}
				{
					position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1272
					}
					position++
					goto l1271
				l1272:
					position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
					if buffer[position] != rune('\n') {
						goto l1264
					}
					position++
				}
			l1271:
				depth--
				add(rulecomment, position1265)
			}
			return true
		l1264:
			position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
			return false
		},
		/* 126 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 128 Action1 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 129 Action2 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 130 Action3 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 131 Action4 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 132 Action5 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 133 Action6 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 134 Action7 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 135 Action8 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 136 Action9 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 137 Action10 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 138 Action11 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 139 Action12 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 140 Action13 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 141 Action14 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 142 Action15 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 143 Action16 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 144 Action17 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 145 Action18 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 146 Action19 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 147 Action20 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 148 Action21 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 149 Action22 <- <{
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
		/* 150 Action23 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 151 Action24 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 152 Action25 <- <{
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
		/* 153 Action26 <- <{
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
		/* 154 Action27 <- <{
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
		/* 155 Action28 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 156 Action29 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 157 Action30 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 158 Action31 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 159 Action32 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 160 Action33 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 161 Action34 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 162 Action35 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 163 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 164 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 165 Action38 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 166 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 167 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 168 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 169 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 170 Action43 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 171 Action44 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 172 Action45 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 173 Action46 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 174 Action47 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 175 Action48 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 176 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 177 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 178 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 179 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 180 Action53 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 181 Action54 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 182 Action55 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 183 Action56 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 184 Action57 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 185 Action58 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 186 Action59 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 187 Action60 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 188 Action61 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 189 Action62 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 190 Action63 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 191 Action64 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 192 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 193 Action66 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 194 Action67 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 195 Action68 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 196 Action69 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 197 Action70 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 198 Action71 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 199 Action72 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 200 Action73 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 201 Action74 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 202 Action75 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 203 Action76 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 204 Action77 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 205 Action78 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 206 Action79 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 207 Action80 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 208 Action81 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 209 Action82 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 210 Action83 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 211 Action84 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 212 Action85 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 213 Action86 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 214 Action87 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 215 Action88 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 216 Action89 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 217 Action90 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 218 Action91 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 219 Action92 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 220 Action93 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 221 Action94 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 222 Action95 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 223 Action96 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 224 Action97 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
	}
	p.rules = _rules
}

package parser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint16

const (
	ruleUnknown pegRule = iota
	ruleSingleStatement
	ruleStatementWithRest
	ruleStatementWithoutRest
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
	ruleLoadStateStmt
	ruleLoadStateOrCreateStmt
	ruleSaveStateStmt
	ruleEmitter
	ruleEmitterOptions
	ruleEmitterOptionCombinations
	ruleEmitterLimit
	ruleEmitterSample
	ruleCountBasedSampling
	ruleRandomizedSampling
	ruleTimeBasedSampling
	ruleTimeBasedSamplingSeconds
	ruleTimeBasedSamplingMilliseconds
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
	ruleSetOptSpecs
	ruleSourceSinkParam
	ruleSourceSinkParamVal
	ruleParamLiteral
	ruleParamArrayExpr
	rulePausedOpt
	ruleExpressionOrWildcard
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
	ruleFuncAppWithOrderBy
	ruleFuncAppWithoutOrderBy
	ruleFuncParams
	ruleParamsOrder
	ruleSortedExpression
	ruleOrderDirectionOpt
	ruleArrayExpr
	ruleMapExpr
	ruleKeyValuePair
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
	ruleMILLISECONDS
	ruleStreamIdentifier
	ruleSourceSinkType
	ruleSourceSinkParamKey
	rulePaused
	ruleUnpaused
	ruleAscending
	ruleDescending
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
	rulejsonGetPath
	rulejsonSetPath
	rulejsonPathHead
	rulejsonGetPathNonHead
	rulejsonSetPathNonHead
	rulejsonMapSingleLevel
	rulejsonMapMultipleLevel
	rulejsonMapAccessString
	rulejsonMapAccessBracket
	rulesingleQuotedString
	rulejsonArrayAccess
	rulejsonNonNegativeArrayAccess
	rulejsonArraySlice
	rulejsonArrayPartialSlice
	rulejsonArrayFullSlice
	rulespElem
	rulesp
	rulespOpt
	rulecomment
	rulefinalComment
	rulePegText
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
	ruleAction98
	ruleAction99
	ruleAction100
	ruleAction101
	ruleAction102
	ruleAction103
	ruleAction104
	ruleAction105
	ruleAction106
	ruleAction107
	ruleAction108
	ruleAction109
	ruleAction110
	ruleAction111
	ruleAction112
	ruleAction113
	ruleAction114
	ruleAction115
	ruleAction116
	ruleAction117
	ruleAction118
	ruleAction119
	ruleAction120

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"SingleStatement",
	"StatementWithRest",
	"StatementWithoutRest",
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
	"LoadStateStmt",
	"LoadStateOrCreateStmt",
	"SaveStateStmt",
	"Emitter",
	"EmitterOptions",
	"EmitterOptionCombinations",
	"EmitterLimit",
	"EmitterSample",
	"CountBasedSampling",
	"RandomizedSampling",
	"TimeBasedSampling",
	"TimeBasedSamplingSeconds",
	"TimeBasedSamplingMilliseconds",
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
	"SetOptSpecs",
	"SourceSinkParam",
	"SourceSinkParamVal",
	"ParamLiteral",
	"ParamArrayExpr",
	"PausedOpt",
	"ExpressionOrWildcard",
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
	"FuncAppWithOrderBy",
	"FuncAppWithoutOrderBy",
	"FuncParams",
	"ParamsOrder",
	"SortedExpression",
	"OrderDirectionOpt",
	"ArrayExpr",
	"MapExpr",
	"KeyValuePair",
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
	"MILLISECONDS",
	"StreamIdentifier",
	"SourceSinkType",
	"SourceSinkParamKey",
	"Paused",
	"Unpaused",
	"Ascending",
	"Descending",
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
	"jsonGetPath",
	"jsonSetPath",
	"jsonPathHead",
	"jsonGetPathNonHead",
	"jsonSetPathNonHead",
	"jsonMapSingleLevel",
	"jsonMapMultipleLevel",
	"jsonMapAccessString",
	"jsonMapAccessBracket",
	"singleQuotedString",
	"jsonArrayAccess",
	"jsonNonNegativeArrayAccess",
	"jsonArraySlice",
	"jsonArrayPartialSlice",
	"jsonArrayFullSlice",
	"spElem",
	"sp",
	"spOpt",
	"comment",
	"finalComment",
	"PegText",
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
	"Action98",
	"Action99",
	"Action100",
	"Action101",
	"Action102",
	"Action103",
	"Action104",
	"Action105",
	"Action106",
	"Action107",
	"Action108",
	"Action109",
	"Action110",
	"Action111",
	"Action112",
	"Action113",
	"Action114",
	"Action115",
	"Action116",
	"Action117",
	"Action118",
	"Action119",
	"Action120",

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

type bqlPegBackend struct {
	parseStack

	Buffer string
	buffer []rune
	rules  [294]func() bool
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
	for i, c := range []rune(buffer) {
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
	p *bqlPegBackend
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

func (p *bqlPegBackend) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *bqlPegBackend) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *bqlPegBackend) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.IncludeTrailingWhitespace(begin, end)

		case ruleAction1:

			p.IncludeTrailingWhitespace(begin, end)

		case ruleAction2:

			p.AssembleSelect()

		case ruleAction3:

			p.AssembleSelectUnion(begin, end)

		case ruleAction4:

			p.AssembleCreateStreamAsSelect()

		case ruleAction5:

			p.AssembleCreateStreamAsSelectUnion()

		case ruleAction6:

			p.AssembleCreateSource()

		case ruleAction7:

			p.AssembleCreateSink()

		case ruleAction8:

			p.AssembleCreateState()

		case ruleAction9:

			p.AssembleUpdateState()

		case ruleAction10:

			p.AssembleUpdateSource()

		case ruleAction11:

			p.AssembleUpdateSink()

		case ruleAction12:

			p.AssembleInsertIntoSelect()

		case ruleAction13:

			p.AssembleInsertIntoFrom()

		case ruleAction14:

			p.AssemblePauseSource()

		case ruleAction15:

			p.AssembleResumeSource()

		case ruleAction16:

			p.AssembleRewindSource()

		case ruleAction17:

			p.AssembleDropSource()

		case ruleAction18:

			p.AssembleDropStream()

		case ruleAction19:

			p.AssembleDropSink()

		case ruleAction20:

			p.AssembleDropState()

		case ruleAction21:

			p.AssembleLoadState()

		case ruleAction22:

			p.AssembleLoadStateOrCreate()

		case ruleAction23:

			p.AssembleSaveState()

		case ruleAction24:

			p.AssembleEmitter()

		case ruleAction25:

			p.AssembleEmitterOptions(begin, end)

		case ruleAction26:

			p.AssembleEmitterLimit()

		case ruleAction27:

			p.AssembleEmitterSampling(CountBasedSampling, 1)

		case ruleAction28:

			p.AssembleEmitterSampling(RandomizedSampling, 1)

		case ruleAction29:

			p.AssembleEmitterSampling(TimeBasedSampling, 1)

		case ruleAction30:

			p.AssembleEmitterSampling(TimeBasedSampling, 0.001)

		case ruleAction31:

			p.AssembleProjections(begin, end)

		case ruleAction32:

			p.AssembleAlias()

		case ruleAction33:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction34:

			p.AssembleInterval()

		case ruleAction35:

			p.AssembleInterval()

		case ruleAction36:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction37:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction38:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction39:

			p.EnsureAliasedStreamWindow()

		case ruleAction40:

			p.AssembleAliasedStreamWindow()

		case ruleAction41:

			p.AssembleStreamWindow()

		case ruleAction42:

			p.AssembleUDSFFuncApp()

		case ruleAction43:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction44:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction45:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction46:

			p.AssembleSourceSinkParam()

		case ruleAction47:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction48:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction49:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction50:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction51:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction52:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction53:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction54:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction55:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction56:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction57:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction58:

			p.AssembleTypeCast(begin, end)

		case ruleAction59:

			p.AssembleTypeCast(begin, end)

		case ruleAction60:

			p.AssembleFuncApp()

		case ruleAction61:

			p.AssembleExpressions(begin, end)
			p.AssembleFuncApp()

		case ruleAction62:

			p.AssembleExpressions(begin, end)

		case ruleAction63:

			p.AssembleExpressions(begin, end)

		case ruleAction64:

			p.AssembleSortedExpression()

		case ruleAction65:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction66:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction67:

			p.AssembleMap(begin, end)

		case ruleAction68:

			p.AssembleKeyValuePair()

		case ruleAction69:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction70:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction71:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction72:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction73:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction74:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction75:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction76:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction77:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction78:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction80:

			p.PushComponent(begin, end, Istream)

		case ruleAction81:

			p.PushComponent(begin, end, Dstream)

		case ruleAction82:

			p.PushComponent(begin, end, Rstream)

		case ruleAction83:

			p.PushComponent(begin, end, Tuples)

		case ruleAction84:

			p.PushComponent(begin, end, Seconds)

		case ruleAction85:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction86:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction87:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction88:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction89:

			p.PushComponent(begin, end, Yes)

		case ruleAction90:

			p.PushComponent(begin, end, No)

		case ruleAction91:

			p.PushComponent(begin, end, Yes)

		case ruleAction92:

			p.PushComponent(begin, end, No)

		case ruleAction93:

			p.PushComponent(begin, end, Bool)

		case ruleAction94:

			p.PushComponent(begin, end, Int)

		case ruleAction95:

			p.PushComponent(begin, end, Float)

		case ruleAction96:

			p.PushComponent(begin, end, String)

		case ruleAction97:

			p.PushComponent(begin, end, Blob)

		case ruleAction98:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction99:

			p.PushComponent(begin, end, Array)

		case ruleAction100:

			p.PushComponent(begin, end, Map)

		case ruleAction101:

			p.PushComponent(begin, end, Or)

		case ruleAction102:

			p.PushComponent(begin, end, And)

		case ruleAction103:

			p.PushComponent(begin, end, Not)

		case ruleAction104:

			p.PushComponent(begin, end, Equal)

		case ruleAction105:

			p.PushComponent(begin, end, Less)

		case ruleAction106:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction107:

			p.PushComponent(begin, end, Greater)

		case ruleAction108:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction109:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction110:

			p.PushComponent(begin, end, Concat)

		case ruleAction111:

			p.PushComponent(begin, end, Is)

		case ruleAction112:

			p.PushComponent(begin, end, IsNot)

		case ruleAction113:

			p.PushComponent(begin, end, Plus)

		case ruleAction114:

			p.PushComponent(begin, end, Minus)

		case ruleAction115:

			p.PushComponent(begin, end, Multiply)

		case ruleAction116:

			p.PushComponent(begin, end, Divide)

		case ruleAction117:

			p.PushComponent(begin, end, Modulo)

		case ruleAction118:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction119:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction120:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		}
	}
	_, _, _, _ = buffer, text, begin, end
}

func (p *bqlPegBackend) Init() {
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
		/* 0 SingleStatement <- <(spOpt (StatementWithRest / StatementWithoutRest) !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[rulespOpt]() {
					goto l0
				}
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[ruleStatementWithRest]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleStatementWithoutRest]() {
						goto l0
					}
				}
			l2:
				{
					position4, tokenIndex4, depth4 := position, tokenIndex, depth
					if !matchDot() {
						goto l4
					}
					goto l0
				l4:
					position, tokenIndex, depth = position4, tokenIndex4, depth4
				}
				depth--
				add(ruleSingleStatement, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 StatementWithRest <- <(<(Statement spOpt ';' spOpt)> .* Action0)> */
		func() bool {
			position5, tokenIndex5, depth5 := position, tokenIndex, depth
			{
				position6 := position
				depth++
				{
					position7 := position
					depth++
					if !_rules[ruleStatement]() {
						goto l5
					}
					if !_rules[rulespOpt]() {
						goto l5
					}
					if buffer[position] != rune(';') {
						goto l5
					}
					position++
					if !_rules[rulespOpt]() {
						goto l5
					}
					depth--
					add(rulePegText, position7)
				}
			l8:
				{
					position9, tokenIndex9, depth9 := position, tokenIndex, depth
					if !matchDot() {
						goto l9
					}
					goto l8
				l9:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
				}
				if !_rules[ruleAction0]() {
					goto l5
				}
				depth--
				add(ruleStatementWithRest, position6)
			}
			return true
		l5:
			position, tokenIndex, depth = position5, tokenIndex5, depth5
			return false
		},
		/* 2 StatementWithoutRest <- <(<(Statement spOpt)> Action1)> */
		func() bool {
			position10, tokenIndex10, depth10 := position, tokenIndex, depth
			{
				position11 := position
				depth++
				{
					position12 := position
					depth++
					if !_rules[ruleStatement]() {
						goto l10
					}
					if !_rules[rulespOpt]() {
						goto l10
					}
					depth--
					add(rulePegText, position12)
				}
				if !_rules[ruleAction1]() {
					goto l10
				}
				depth--
				add(ruleStatementWithoutRest, position11)
			}
			return true
		l10:
			position, tokenIndex, depth = position10, tokenIndex10, depth10
			return false
		},
		/* 3 Statement <- <(SelectUnionStmt / SelectStmt / SourceStmt / SinkStmt / StateStmt / StreamStmt)> */
		func() bool {
			position13, tokenIndex13, depth13 := position, tokenIndex, depth
			{
				position14 := position
				depth++
				{
					position15, tokenIndex15, depth15 := position, tokenIndex, depth
					if !_rules[ruleSelectUnionStmt]() {
						goto l16
					}
					goto l15
				l16:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleSelectStmt]() {
						goto l17
					}
					goto l15
				l17:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleSourceStmt]() {
						goto l18
					}
					goto l15
				l18:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleSinkStmt]() {
						goto l19
					}
					goto l15
				l19:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleStateStmt]() {
						goto l20
					}
					goto l15
				l20:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleStreamStmt]() {
						goto l13
					}
				}
			l15:
				depth--
				add(ruleStatement, position14)
			}
			return true
		l13:
			position, tokenIndex, depth = position13, tokenIndex13, depth13
			return false
		},
		/* 4 SourceStmt <- <(CreateSourceStmt / UpdateSourceStmt / DropSourceStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt)> */
		func() bool {
			position21, tokenIndex21, depth21 := position, tokenIndex, depth
			{
				position22 := position
				depth++
				{
					position23, tokenIndex23, depth23 := position, tokenIndex, depth
					if !_rules[ruleCreateSourceStmt]() {
						goto l24
					}
					goto l23
				l24:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if !_rules[ruleUpdateSourceStmt]() {
						goto l25
					}
					goto l23
				l25:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if !_rules[ruleDropSourceStmt]() {
						goto l26
					}
					goto l23
				l26:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if !_rules[rulePauseSourceStmt]() {
						goto l27
					}
					goto l23
				l27:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if !_rules[ruleResumeSourceStmt]() {
						goto l28
					}
					goto l23
				l28:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if !_rules[ruleRewindSourceStmt]() {
						goto l21
					}
				}
			l23:
				depth--
				add(ruleSourceStmt, position22)
			}
			return true
		l21:
			position, tokenIndex, depth = position21, tokenIndex21, depth21
			return false
		},
		/* 5 SinkStmt <- <(CreateSinkStmt / UpdateSinkStmt / DropSinkStmt)> */
		func() bool {
			position29, tokenIndex29, depth29 := position, tokenIndex, depth
			{
				position30 := position
				depth++
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					if !_rules[ruleCreateSinkStmt]() {
						goto l32
					}
					goto l31
				l32:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
					if !_rules[ruleUpdateSinkStmt]() {
						goto l33
					}
					goto l31
				l33:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
					if !_rules[ruleDropSinkStmt]() {
						goto l29
					}
				}
			l31:
				depth--
				add(ruleSinkStmt, position30)
			}
			return true
		l29:
			position, tokenIndex, depth = position29, tokenIndex29, depth29
			return false
		},
		/* 6 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt / LoadStateOrCreateStmt / LoadStateStmt / SaveStateStmt)> */
		func() bool {
			position34, tokenIndex34, depth34 := position, tokenIndex, depth
			{
				position35 := position
				depth++
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					if !_rules[ruleCreateStateStmt]() {
						goto l37
					}
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if !_rules[ruleUpdateStateStmt]() {
						goto l38
					}
					goto l36
				l38:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if !_rules[ruleDropStateStmt]() {
						goto l39
					}
					goto l36
				l39:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if !_rules[ruleLoadStateOrCreateStmt]() {
						goto l40
					}
					goto l36
				l40:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if !_rules[ruleLoadStateStmt]() {
						goto l41
					}
					goto l36
				l41:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if !_rules[ruleSaveStateStmt]() {
						goto l34
					}
				}
			l36:
				depth--
				add(ruleStateStmt, position35)
			}
			return true
		l34:
			position, tokenIndex, depth = position34, tokenIndex34, depth34
			return false
		},
		/* 7 StreamStmt <- <(CreateStreamAsSelectUnionStmt / CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoSelectStmt / InsertIntoFromStmt)> */
		func() bool {
			position42, tokenIndex42, depth42 := position, tokenIndex, depth
			{
				position43 := position
				depth++
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l45
					}
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l46
					}
					goto l44
				l46:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if !_rules[ruleDropStreamStmt]() {
						goto l47
					}
					goto l44
				l47:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l48
					}
					goto l44
				l48:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l42
					}
				}
			l44:
				depth--
				add(ruleStreamStmt, position43)
			}
			return true
		l42:
			position, tokenIndex, depth = position42, tokenIndex42, depth42
			return false
		},
		/* 8 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') Emitter Projections WindowedFrom Filter Grouping Having Action2)> */
		func() bool {
			position49, tokenIndex49, depth49 := position, tokenIndex, depth
			{
				position50 := position
				depth++
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
						goto l49
					}
					position++
				}
			l51:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					if buffer[position] != rune('E') {
						goto l49
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
					if buffer[position] != rune('L') {
						goto l49
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
						goto l49
					}
					position++
				}
			l57:
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l60
					}
					position++
					goto l59
				l60:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
					if buffer[position] != rune('C') {
						goto l49
					}
					position++
				}
			l59:
				{
					position61, tokenIndex61, depth61 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l62
					}
					position++
					goto l61
				l62:
					position, tokenIndex, depth = position61, tokenIndex61, depth61
					if buffer[position] != rune('T') {
						goto l49
					}
					position++
				}
			l61:
				if !_rules[ruleEmitter]() {
					goto l49
				}
				if !_rules[ruleProjections]() {
					goto l49
				}
				if !_rules[ruleWindowedFrom]() {
					goto l49
				}
				if !_rules[ruleFilter]() {
					goto l49
				}
				if !_rules[ruleGrouping]() {
					goto l49
				}
				if !_rules[ruleHaving]() {
					goto l49
				}
				if !_rules[ruleAction2]() {
					goto l49
				}
				depth--
				add(ruleSelectStmt, position50)
			}
			return true
		l49:
			position, tokenIndex, depth = position49, tokenIndex49, depth49
			return false
		},
		/* 9 SelectUnionStmt <- <(<(SelectStmt (sp (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N')) sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action3)> */
		func() bool {
			position63, tokenIndex63, depth63 := position, tokenIndex, depth
			{
				position64 := position
				depth++
				{
					position65 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l63
					}
					if !_rules[rulesp]() {
						goto l63
					}
					{
						position68, tokenIndex68, depth68 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l69
						}
						position++
						goto l68
					l69:
						position, tokenIndex, depth = position68, tokenIndex68, depth68
						if buffer[position] != rune('U') {
							goto l63
						}
						position++
					}
				l68:
					{
						position70, tokenIndex70, depth70 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l71
						}
						position++
						goto l70
					l71:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
						if buffer[position] != rune('N') {
							goto l63
						}
						position++
					}
				l70:
					{
						position72, tokenIndex72, depth72 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l73
						}
						position++
						goto l72
					l73:
						position, tokenIndex, depth = position72, tokenIndex72, depth72
						if buffer[position] != rune('I') {
							goto l63
						}
						position++
					}
				l72:
					{
						position74, tokenIndex74, depth74 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l75
						}
						position++
						goto l74
					l75:
						position, tokenIndex, depth = position74, tokenIndex74, depth74
						if buffer[position] != rune('O') {
							goto l63
						}
						position++
					}
				l74:
					{
						position76, tokenIndex76, depth76 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l77
						}
						position++
						goto l76
					l77:
						position, tokenIndex, depth = position76, tokenIndex76, depth76
						if buffer[position] != rune('N') {
							goto l63
						}
						position++
					}
				l76:
					if !_rules[rulesp]() {
						goto l63
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
							goto l63
						}
						position++
					}
				l78:
					{
						position80, tokenIndex80, depth80 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l81
						}
						position++
						goto l80
					l81:
						position, tokenIndex, depth = position80, tokenIndex80, depth80
						if buffer[position] != rune('L') {
							goto l63
						}
						position++
					}
				l80:
					{
						position82, tokenIndex82, depth82 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l83
						}
						position++
						goto l82
					l83:
						position, tokenIndex, depth = position82, tokenIndex82, depth82
						if buffer[position] != rune('L') {
							goto l63
						}
						position++
					}
				l82:
					if !_rules[rulesp]() {
						goto l63
					}
					if !_rules[ruleSelectStmt]() {
						goto l63
					}
				l66:
					{
						position67, tokenIndex67, depth67 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l67
						}
						{
							position84, tokenIndex84, depth84 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l85
							}
							position++
							goto l84
						l85:
							position, tokenIndex, depth = position84, tokenIndex84, depth84
							if buffer[position] != rune('U') {
								goto l67
							}
							position++
						}
					l84:
						{
							position86, tokenIndex86, depth86 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l87
							}
							position++
							goto l86
						l87:
							position, tokenIndex, depth = position86, tokenIndex86, depth86
							if buffer[position] != rune('N') {
								goto l67
							}
							position++
						}
					l86:
						{
							position88, tokenIndex88, depth88 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l89
							}
							position++
							goto l88
						l89:
							position, tokenIndex, depth = position88, tokenIndex88, depth88
							if buffer[position] != rune('I') {
								goto l67
							}
							position++
						}
					l88:
						{
							position90, tokenIndex90, depth90 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l91
							}
							position++
							goto l90
						l91:
							position, tokenIndex, depth = position90, tokenIndex90, depth90
							if buffer[position] != rune('O') {
								goto l67
							}
							position++
						}
					l90:
						{
							position92, tokenIndex92, depth92 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l93
							}
							position++
							goto l92
						l93:
							position, tokenIndex, depth = position92, tokenIndex92, depth92
							if buffer[position] != rune('N') {
								goto l67
							}
							position++
						}
					l92:
						if !_rules[rulesp]() {
							goto l67
						}
						{
							position94, tokenIndex94, depth94 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l95
							}
							position++
							goto l94
						l95:
							position, tokenIndex, depth = position94, tokenIndex94, depth94
							if buffer[position] != rune('A') {
								goto l67
							}
							position++
						}
					l94:
						{
							position96, tokenIndex96, depth96 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l97
							}
							position++
							goto l96
						l97:
							position, tokenIndex, depth = position96, tokenIndex96, depth96
							if buffer[position] != rune('L') {
								goto l67
							}
							position++
						}
					l96:
						{
							position98, tokenIndex98, depth98 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l99
							}
							position++
							goto l98
						l99:
							position, tokenIndex, depth = position98, tokenIndex98, depth98
							if buffer[position] != rune('L') {
								goto l67
							}
							position++
						}
					l98:
						if !_rules[rulesp]() {
							goto l67
						}
						if !_rules[ruleSelectStmt]() {
							goto l67
						}
						goto l66
					l67:
						position, tokenIndex, depth = position67, tokenIndex67, depth67
					}
					depth--
					add(rulePegText, position65)
				}
				if !_rules[ruleAction3]() {
					goto l63
				}
				depth--
				add(ruleSelectUnionStmt, position64)
			}
			return true
		l63:
			position, tokenIndex, depth = position63, tokenIndex63, depth63
			return false
		},
		/* 10 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action4)> */
		func() bool {
			position100, tokenIndex100, depth100 := position, tokenIndex, depth
			{
				position101 := position
				depth++
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
						goto l100
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('R') {
						goto l100
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
						goto l100
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('T') {
						goto l100
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
						goto l100
					}
					position++
				}
			l112:
				if !_rules[rulesp]() {
					goto l100
				}
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('S') {
						goto l100
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('T') {
						goto l100
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
						goto l100
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
						goto l100
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
						goto l100
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('M') {
						goto l100
					}
					position++
				}
			l124:
				if !_rules[rulesp]() {
					goto l100
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l100
				}
				if !_rules[rulesp]() {
					goto l100
				}
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l126:
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
						goto l100
					}
					position++
				}
			l128:
				if !_rules[rulesp]() {
					goto l100
				}
				if !_rules[ruleSelectStmt]() {
					goto l100
				}
				if !_rules[ruleAction4]() {
					goto l100
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 11 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action5)> */
		func() bool {
			position130, tokenIndex130, depth130 := position, tokenIndex, depth
			{
				position131 := position
				depth++
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('C') {
						goto l130
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('R') {
						goto l130
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('E') {
						goto l130
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('T') {
						goto l130
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
						goto l130
					}
					position++
				}
			l142:
				if !_rules[rulesp]() {
					goto l130
				}
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('S') {
						goto l130
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('T') {
						goto l130
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
						goto l130
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
						goto l130
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
						goto l130
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('M') {
						goto l130
					}
					position++
				}
			l154:
				if !_rules[rulesp]() {
					goto l130
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l130
				}
				if !_rules[rulesp]() {
					goto l130
				}
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l156:
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
						goto l130
					}
					position++
				}
			l158:
				if !_rules[rulesp]() {
					goto l130
				}
				if !_rules[ruleSelectUnionStmt]() {
					goto l130
				}
				if !_rules[ruleAction5]() {
					goto l130
				}
				depth--
				add(ruleCreateStreamAsSelectUnionStmt, position131)
			}
			return true
		l130:
			position, tokenIndex, depth = position130, tokenIndex130, depth130
			return false
		},
		/* 12 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action6)> */
		func() bool {
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('C') {
						goto l160
					}
					position++
				}
			l162:
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('R') {
						goto l160
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
						goto l160
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('A') {
						goto l160
					}
					position++
				}
			l168:
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
						goto l160
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
						goto l160
					}
					position++
				}
			l172:
				if !_rules[rulePausedOpt]() {
					goto l160
				}
				if !_rules[rulesp]() {
					goto l160
				}
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('S') {
						goto l160
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('O') {
						goto l160
					}
					position++
				}
			l176:
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
						goto l160
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('R') {
						goto l160
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('C') {
						goto l160
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('E') {
						goto l160
					}
					position++
				}
			l184:
				if !_rules[rulesp]() {
					goto l160
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l160
				}
				if !_rules[rulesp]() {
					goto l160
				}
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
						goto l160
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('Y') {
						goto l160
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('P') {
						goto l160
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('E') {
						goto l160
					}
					position++
				}
			l192:
				if !_rules[rulesp]() {
					goto l160
				}
				if !_rules[ruleSourceSinkType]() {
					goto l160
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l160
				}
				if !_rules[ruleAction6]() {
					goto l160
				}
				depth--
				add(ruleCreateSourceStmt, position161)
			}
			return true
		l160:
			position, tokenIndex, depth = position160, tokenIndex160, depth160
			return false
		},
		/* 13 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action7)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('C') {
						goto l194
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('R') {
						goto l194
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
						goto l194
					}
					position++
				}
			l200:
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('A') {
						goto l194
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('T') {
						goto l194
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('E') {
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
					if buffer[position] != rune('s') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('S') {
						goto l194
					}
					position++
				}
			l208:
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
						goto l194
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
						goto l194
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('K') {
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
					if buffer[position] != rune('t') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('T') {
						goto l194
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('Y') {
						goto l194
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('P') {
						goto l194
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('E') {
						goto l194
					}
					position++
				}
			l222:
				if !_rules[rulesp]() {
					goto l194
				}
				if !_rules[ruleSourceSinkType]() {
					goto l194
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l194
				}
				if !_rules[ruleAction7]() {
					goto l194
				}
				depth--
				add(ruleCreateSinkStmt, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 14 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action8)> */
		func() bool {
			position224, tokenIndex224, depth224 := position, tokenIndex, depth
			{
				position225 := position
				depth++
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('C') {
						goto l224
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('R') {
						goto l224
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
						goto l224
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('A') {
						goto l224
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
						goto l224
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
						goto l224
					}
					position++
				}
			l236:
				if !_rules[rulesp]() {
					goto l224
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
						goto l224
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
						goto l224
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('A') {
						goto l224
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('T') {
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
				if !_rules[rulesp]() {
					goto l224
				}
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('Y') {
						goto l224
					}
					position++
				}
			l250:
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if buffer[position] != rune('P') {
						goto l224
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
						goto l224
					}
					position++
				}
			l254:
				if !_rules[rulesp]() {
					goto l224
				}
				if !_rules[ruleSourceSinkType]() {
					goto l224
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l224
				}
				if !_rules[ruleAction8]() {
					goto l224
				}
				depth--
				add(ruleCreateStateStmt, position225)
			}
			return true
		l224:
			position, tokenIndex, depth = position224, tokenIndex224, depth224
			return false
		},
		/* 15 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position256, tokenIndex256, depth256 := position, tokenIndex, depth
			{
				position257 := position
				depth++
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
						goto l256
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('P') {
						goto l256
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('D') {
						goto l256
					}
					position++
				}
			l262:
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('A') {
						goto l256
					}
					position++
				}
			l264:
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('E') {
						goto l256
					}
					position++
				}
			l268:
				if !_rules[rulesp]() {
					goto l256
				}
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if buffer[position] != rune('S') {
						goto l256
					}
					position++
				}
			l270:
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l273
					}
					position++
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l272:
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l275
					}
					position++
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('A') {
						goto l256
					}
					position++
				}
			l274:
				{
					position276, tokenIndex276, depth276 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l277
					}
					position++
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					if buffer[position] != rune('T') {
						goto l256
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
						goto l256
					}
					position++
				}
			l278:
				if !_rules[rulesp]() {
					goto l256
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l256
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l256
				}
				if !_rules[ruleAction9]() {
					goto l256
				}
				depth--
				add(ruleUpdateStateStmt, position257)
			}
			return true
		l256:
			position, tokenIndex, depth = position256, tokenIndex256, depth256
			return false
		},
		/* 16 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action10)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('U') {
						goto l280
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
					if buffer[position] != rune('P') {
						goto l280
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
						goto l280
					}
					position++
				}
			l286:
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l289
					}
					position++
					goto l288
				l289:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
					if buffer[position] != rune('A') {
						goto l280
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
						goto l280
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
						goto l280
					}
					position++
				}
			l292:
				if !_rules[rulesp]() {
					goto l280
				}
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l295
					}
					position++
					goto l294
				l295:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if buffer[position] != rune('S') {
						goto l280
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if buffer[position] != rune('O') {
						goto l280
					}
					position++
				}
			l296:
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					if buffer[position] != rune('U') {
						goto l280
					}
					position++
				}
			l298:
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l301
					}
					position++
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if buffer[position] != rune('R') {
						goto l280
					}
					position++
				}
			l300:
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l303
					}
					position++
					goto l302
				l303:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					if buffer[position] != rune('C') {
						goto l280
					}
					position++
				}
			l302:
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
					if buffer[position] != rune('E') {
						goto l280
					}
					position++
				}
			l304:
				if !_rules[rulesp]() {
					goto l280
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l280
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l280
				}
				if !_rules[ruleAction10]() {
					goto l280
				}
				depth--
				add(ruleUpdateSourceStmt, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 17 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier UpdateSourceSinkSpecs Action11)> */
		func() bool {
			position306, tokenIndex306, depth306 := position, tokenIndex, depth
			{
				position307 := position
				depth++
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
						goto l306
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
						goto l306
					}
					position++
				}
			l310:
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('D') {
						goto l306
					}
					position++
				}
			l312:
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l315
					}
					position++
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if buffer[position] != rune('A') {
						goto l306
					}
					position++
				}
			l314:
				{
					position316, tokenIndex316, depth316 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l317
					}
					position++
					goto l316
				l317:
					position, tokenIndex, depth = position316, tokenIndex316, depth316
					if buffer[position] != rune('T') {
						goto l306
					}
					position++
				}
			l316:
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
					if buffer[position] != rune('E') {
						goto l306
					}
					position++
				}
			l318:
				if !_rules[rulesp]() {
					goto l306
				}
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('S') {
						goto l306
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if buffer[position] != rune('I') {
						goto l306
					}
					position++
				}
			l322:
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l325
					}
					position++
					goto l324
				l325:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if buffer[position] != rune('N') {
						goto l306
					}
					position++
				}
			l324:
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l327
					}
					position++
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					if buffer[position] != rune('K') {
						goto l306
					}
					position++
				}
			l326:
				if !_rules[rulesp]() {
					goto l306
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l306
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l306
				}
				if !_rules[ruleAction11]() {
					goto l306
				}
				depth--
				add(ruleUpdateSinkStmt, position307)
			}
			return true
		l306:
			position, tokenIndex, depth = position306, tokenIndex306, depth306
			return false
		},
		/* 18 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action12)> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l331
					}
					position++
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					if buffer[position] != rune('I') {
						goto l328
					}
					position++
				}
			l330:
				{
					position332, tokenIndex332, depth332 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l333
					}
					position++
					goto l332
				l333:
					position, tokenIndex, depth = position332, tokenIndex332, depth332
					if buffer[position] != rune('N') {
						goto l328
					}
					position++
				}
			l332:
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
						goto l328
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
						goto l328
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
						goto l328
					}
					position++
				}
			l338:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l341
					}
					position++
					goto l340
				l341:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
					if buffer[position] != rune('T') {
						goto l328
					}
					position++
				}
			l340:
				if !_rules[rulesp]() {
					goto l328
				}
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if buffer[position] != rune('I') {
						goto l328
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
					if buffer[position] != rune('N') {
						goto l328
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					if buffer[position] != rune('T') {
						goto l328
					}
					position++
				}
			l346:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if buffer[position] != rune('O') {
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
				if !_rules[rulesp]() {
					goto l328
				}
				if !_rules[ruleSelectStmt]() {
					goto l328
				}
				if !_rules[ruleAction12]() {
					goto l328
				}
				depth--
				add(ruleInsertIntoSelectStmt, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 19 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action13)> */
		func() bool {
			position350, tokenIndex350, depth350 := position, tokenIndex, depth
			{
				position351 := position
				depth++
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('I') {
						goto l350
					}
					position++
				}
			l352:
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
					if buffer[position] != rune('N') {
						goto l350
					}
					position++
				}
			l354:
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
						goto l350
					}
					position++
				}
			l356:
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l359
					}
					position++
					goto l358
				l359:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					if buffer[position] != rune('E') {
						goto l350
					}
					position++
				}
			l358:
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if buffer[position] != rune('R') {
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
				if !_rules[rulesp]() {
					goto l350
				}
				{
					position364, tokenIndex364, depth364 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l365
					}
					position++
					goto l364
				l365:
					position, tokenIndex, depth = position364, tokenIndex364, depth364
					if buffer[position] != rune('I') {
						goto l350
					}
					position++
				}
			l364:
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l367
					}
					position++
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if buffer[position] != rune('N') {
						goto l350
					}
					position++
				}
			l366:
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if buffer[position] != rune('T') {
						goto l350
					}
					position++
				}
			l368:
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l371
					}
					position++
					goto l370
				l371:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
					if buffer[position] != rune('O') {
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
				if !_rules[rulesp]() {
					goto l350
				}
				{
					position372, tokenIndex372, depth372 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l373
					}
					position++
					goto l372
				l373:
					position, tokenIndex, depth = position372, tokenIndex372, depth372
					if buffer[position] != rune('F') {
						goto l350
					}
					position++
				}
			l372:
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('R') {
						goto l350
					}
					position++
				}
			l374:
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l377
					}
					position++
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					if buffer[position] != rune('O') {
						goto l350
					}
					position++
				}
			l376:
				{
					position378, tokenIndex378, depth378 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l379
					}
					position++
					goto l378
				l379:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
					if buffer[position] != rune('M') {
						goto l350
					}
					position++
				}
			l378:
				if !_rules[rulesp]() {
					goto l350
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l350
				}
				if !_rules[ruleAction13]() {
					goto l350
				}
				depth--
				add(ruleInsertIntoFromStmt, position351)
			}
			return true
		l350:
			position, tokenIndex, depth = position350, tokenIndex350, depth350
			return false
		},
		/* 20 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position380, tokenIndex380, depth380 := position, tokenIndex, depth
			{
				position381 := position
				depth++
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l383
					}
					position++
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if buffer[position] != rune('P') {
						goto l380
					}
					position++
				}
			l382:
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l385
					}
					position++
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if buffer[position] != rune('A') {
						goto l380
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('U') {
						goto l380
					}
					position++
				}
			l386:
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
						goto l380
					}
					position++
				}
			l388:
				{
					position390, tokenIndex390, depth390 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l391
					}
					position++
					goto l390
				l391:
					position, tokenIndex, depth = position390, tokenIndex390, depth390
					if buffer[position] != rune('E') {
						goto l380
					}
					position++
				}
			l390:
				if !_rules[rulesp]() {
					goto l380
				}
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if buffer[position] != rune('S') {
						goto l380
					}
					position++
				}
			l392:
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l395
					}
					position++
					goto l394
				l395:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
					if buffer[position] != rune('O') {
						goto l380
					}
					position++
				}
			l394:
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('U') {
						goto l380
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if buffer[position] != rune('R') {
						goto l380
					}
					position++
				}
			l398:
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l401
					}
					position++
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					if buffer[position] != rune('C') {
						goto l380
					}
					position++
				}
			l400:
				{
					position402, tokenIndex402, depth402 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l403
					}
					position++
					goto l402
				l403:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
					if buffer[position] != rune('E') {
						goto l380
					}
					position++
				}
			l402:
				if !_rules[rulesp]() {
					goto l380
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l380
				}
				if !_rules[ruleAction14]() {
					goto l380
				}
				depth--
				add(rulePauseSourceStmt, position381)
			}
			return true
		l380:
			position, tokenIndex, depth = position380, tokenIndex380, depth380
			return false
		},
		/* 21 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
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
						goto l404
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
						goto l404
					}
					position++
				}
			l408:
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
						goto l404
					}
					position++
				}
			l410:
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l413
					}
					position++
					goto l412
				l413:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
					if buffer[position] != rune('U') {
						goto l404
					}
					position++
				}
			l412:
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('M') {
						goto l404
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
						goto l404
					}
					position++
				}
			l416:
				if !_rules[rulesp]() {
					goto l404
				}
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
					if buffer[position] != rune('S') {
						goto l404
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('O') {
						goto l404
					}
					position++
				}
			l420:
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l423
					}
					position++
					goto l422
				l423:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
					if buffer[position] != rune('U') {
						goto l404
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
						goto l404
					}
					position++
				}
			l424:
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l427
					}
					position++
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if buffer[position] != rune('C') {
						goto l404
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
						goto l404
					}
					position++
				}
			l428:
				if !_rules[rulesp]() {
					goto l404
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l404
				}
				if !_rules[ruleAction15]() {
					goto l404
				}
				depth--
				add(ruleResumeSourceStmt, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 22 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position430, tokenIndex430, depth430 := position, tokenIndex, depth
			{
				position431 := position
				depth++
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('R') {
						goto l430
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
						goto l430
					}
					position++
				}
			l434:
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l437
					}
					position++
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if buffer[position] != rune('W') {
						goto l430
					}
					position++
				}
			l436:
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if buffer[position] != rune('I') {
						goto l430
					}
					position++
				}
			l438:
				{
					position440, tokenIndex440, depth440 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l441
					}
					position++
					goto l440
				l441:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if buffer[position] != rune('N') {
						goto l430
					}
					position++
				}
			l440:
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
						goto l430
					}
					position++
				}
			l442:
				if !_rules[rulesp]() {
					goto l430
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
						goto l430
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
						goto l430
					}
					position++
				}
			l446:
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
					if buffer[position] != rune('U') {
						goto l430
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('R') {
						goto l430
					}
					position++
				}
			l450:
				{
					position452, tokenIndex452, depth452 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l453
					}
					position++
					goto l452
				l453:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if buffer[position] != rune('C') {
						goto l430
					}
					position++
				}
			l452:
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l455
					}
					position++
					goto l454
				l455:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if buffer[position] != rune('E') {
						goto l430
					}
					position++
				}
			l454:
				if !_rules[rulesp]() {
					goto l430
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l430
				}
				if !_rules[ruleAction16]() {
					goto l430
				}
				depth--
				add(ruleRewindSourceStmt, position431)
			}
			return true
		l430:
			position, tokenIndex, depth = position430, tokenIndex430, depth430
			return false
		},
		/* 23 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action17)> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l459
					}
					position++
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if buffer[position] != rune('D') {
						goto l456
					}
					position++
				}
			l458:
				{
					position460, tokenIndex460, depth460 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l461
					}
					position++
					goto l460
				l461:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
					if buffer[position] != rune('R') {
						goto l456
					}
					position++
				}
			l460:
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l463
					}
					position++
					goto l462
				l463:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if buffer[position] != rune('O') {
						goto l456
					}
					position++
				}
			l462:
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l465
					}
					position++
					goto l464
				l465:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					if buffer[position] != rune('P') {
						goto l456
					}
					position++
				}
			l464:
				if !_rules[rulesp]() {
					goto l456
				}
				{
					position466, tokenIndex466, depth466 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l467
					}
					position++
					goto l466
				l467:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if buffer[position] != rune('S') {
						goto l456
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
						goto l456
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
						goto l456
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
						goto l456
					}
					position++
				}
			l472:
				{
					position474, tokenIndex474, depth474 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l475
					}
					position++
					goto l474
				l475:
					position, tokenIndex, depth = position474, tokenIndex474, depth474
					if buffer[position] != rune('C') {
						goto l456
					}
					position++
				}
			l474:
				{
					position476, tokenIndex476, depth476 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l477
					}
					position++
					goto l476
				l477:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if buffer[position] != rune('E') {
						goto l456
					}
					position++
				}
			l476:
				if !_rules[rulesp]() {
					goto l456
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l456
				}
				if !_rules[ruleAction17]() {
					goto l456
				}
				depth--
				add(ruleDropSourceStmt, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 24 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action18)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				{
					position480, tokenIndex480, depth480 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l481
					}
					position++
					goto l480
				l481:
					position, tokenIndex, depth = position480, tokenIndex480, depth480
					if buffer[position] != rune('D') {
						goto l478
					}
					position++
				}
			l480:
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l483
					}
					position++
					goto l482
				l483:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
					if buffer[position] != rune('R') {
						goto l478
					}
					position++
				}
			l482:
				{
					position484, tokenIndex484, depth484 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l485
					}
					position++
					goto l484
				l485:
					position, tokenIndex, depth = position484, tokenIndex484, depth484
					if buffer[position] != rune('O') {
						goto l478
					}
					position++
				}
			l484:
				{
					position486, tokenIndex486, depth486 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l487
					}
					position++
					goto l486
				l487:
					position, tokenIndex, depth = position486, tokenIndex486, depth486
					if buffer[position] != rune('P') {
						goto l478
					}
					position++
				}
			l486:
				if !_rules[rulesp]() {
					goto l478
				}
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
						goto l478
					}
					position++
				}
			l488:
				{
					position490, tokenIndex490, depth490 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l491
					}
					position++
					goto l490
				l491:
					position, tokenIndex, depth = position490, tokenIndex490, depth490
					if buffer[position] != rune('T') {
						goto l478
					}
					position++
				}
			l490:
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
						goto l478
					}
					position++
				}
			l492:
				{
					position494, tokenIndex494, depth494 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l495
					}
					position++
					goto l494
				l495:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if buffer[position] != rune('E') {
						goto l478
					}
					position++
				}
			l494:
				{
					position496, tokenIndex496, depth496 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l497
					}
					position++
					goto l496
				l497:
					position, tokenIndex, depth = position496, tokenIndex496, depth496
					if buffer[position] != rune('A') {
						goto l478
					}
					position++
				}
			l496:
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l499
					}
					position++
					goto l498
				l499:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if buffer[position] != rune('M') {
						goto l478
					}
					position++
				}
			l498:
				if !_rules[rulesp]() {
					goto l478
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l478
				}
				if !_rules[ruleAction18]() {
					goto l478
				}
				depth--
				add(ruleDropStreamStmt, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 25 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action19)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				{
					position502, tokenIndex502, depth502 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l503
					}
					position++
					goto l502
				l503:
					position, tokenIndex, depth = position502, tokenIndex502, depth502
					if buffer[position] != rune('D') {
						goto l500
					}
					position++
				}
			l502:
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l505
					}
					position++
					goto l504
				l505:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
					if buffer[position] != rune('R') {
						goto l500
					}
					position++
				}
			l504:
				{
					position506, tokenIndex506, depth506 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l507
					}
					position++
					goto l506
				l507:
					position, tokenIndex, depth = position506, tokenIndex506, depth506
					if buffer[position] != rune('O') {
						goto l500
					}
					position++
				}
			l506:
				{
					position508, tokenIndex508, depth508 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l509
					}
					position++
					goto l508
				l509:
					position, tokenIndex, depth = position508, tokenIndex508, depth508
					if buffer[position] != rune('P') {
						goto l500
					}
					position++
				}
			l508:
				if !_rules[rulesp]() {
					goto l500
				}
				{
					position510, tokenIndex510, depth510 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l511
					}
					position++
					goto l510
				l511:
					position, tokenIndex, depth = position510, tokenIndex510, depth510
					if buffer[position] != rune('S') {
						goto l500
					}
					position++
				}
			l510:
				{
					position512, tokenIndex512, depth512 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l513
					}
					position++
					goto l512
				l513:
					position, tokenIndex, depth = position512, tokenIndex512, depth512
					if buffer[position] != rune('I') {
						goto l500
					}
					position++
				}
			l512:
				{
					position514, tokenIndex514, depth514 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l515
					}
					position++
					goto l514
				l515:
					position, tokenIndex, depth = position514, tokenIndex514, depth514
					if buffer[position] != rune('N') {
						goto l500
					}
					position++
				}
			l514:
				{
					position516, tokenIndex516, depth516 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l517
					}
					position++
					goto l516
				l517:
					position, tokenIndex, depth = position516, tokenIndex516, depth516
					if buffer[position] != rune('K') {
						goto l500
					}
					position++
				}
			l516:
				if !_rules[rulesp]() {
					goto l500
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l500
				}
				if !_rules[ruleAction19]() {
					goto l500
				}
				depth--
				add(ruleDropSinkStmt, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 26 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action20)> */
		func() bool {
			position518, tokenIndex518, depth518 := position, tokenIndex, depth
			{
				position519 := position
				depth++
				{
					position520, tokenIndex520, depth520 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l521
					}
					position++
					goto l520
				l521:
					position, tokenIndex, depth = position520, tokenIndex520, depth520
					if buffer[position] != rune('D') {
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
					if buffer[position] != rune('p') {
						goto l527
					}
					position++
					goto l526
				l527:
					position, tokenIndex, depth = position526, tokenIndex526, depth526
					if buffer[position] != rune('P') {
						goto l518
					}
					position++
				}
			l526:
				if !_rules[rulesp]() {
					goto l518
				}
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
						goto l518
					}
					position++
				}
			l528:
				{
					position530, tokenIndex530, depth530 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l531
					}
					position++
					goto l530
				l531:
					position, tokenIndex, depth = position530, tokenIndex530, depth530
					if buffer[position] != rune('T') {
						goto l518
					}
					position++
				}
			l530:
				{
					position532, tokenIndex532, depth532 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l533
					}
					position++
					goto l532
				l533:
					position, tokenIndex, depth = position532, tokenIndex532, depth532
					if buffer[position] != rune('A') {
						goto l518
					}
					position++
				}
			l532:
				{
					position534, tokenIndex534, depth534 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l535
					}
					position++
					goto l534
				l535:
					position, tokenIndex, depth = position534, tokenIndex534, depth534
					if buffer[position] != rune('T') {
						goto l518
					}
					position++
				}
			l534:
				{
					position536, tokenIndex536, depth536 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l537
					}
					position++
					goto l536
				l537:
					position, tokenIndex, depth = position536, tokenIndex536, depth536
					if buffer[position] != rune('E') {
						goto l518
					}
					position++
				}
			l536:
				if !_rules[rulesp]() {
					goto l518
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l518
				}
				if !_rules[ruleAction20]() {
					goto l518
				}
				depth--
				add(ruleDropStateStmt, position519)
			}
			return true
		l518:
			position, tokenIndex, depth = position518, tokenIndex518, depth518
			return false
		},
		/* 27 LoadStateStmt <- <(('l' / 'L') ('o' / 'O') ('a' / 'A') ('d' / 'D') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SetOptSpecs Action21)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					position540, tokenIndex540, depth540 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l541
					}
					position++
					goto l540
				l541:
					position, tokenIndex, depth = position540, tokenIndex540, depth540
					if buffer[position] != rune('L') {
						goto l538
					}
					position++
				}
			l540:
				{
					position542, tokenIndex542, depth542 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l543
					}
					position++
					goto l542
				l543:
					position, tokenIndex, depth = position542, tokenIndex542, depth542
					if buffer[position] != rune('O') {
						goto l538
					}
					position++
				}
			l542:
				{
					position544, tokenIndex544, depth544 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l545
					}
					position++
					goto l544
				l545:
					position, tokenIndex, depth = position544, tokenIndex544, depth544
					if buffer[position] != rune('A') {
						goto l538
					}
					position++
				}
			l544:
				{
					position546, tokenIndex546, depth546 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l547
					}
					position++
					goto l546
				l547:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
					if buffer[position] != rune('D') {
						goto l538
					}
					position++
				}
			l546:
				if !_rules[rulesp]() {
					goto l538
				}
				{
					position548, tokenIndex548, depth548 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l549
					}
					position++
					goto l548
				l549:
					position, tokenIndex, depth = position548, tokenIndex548, depth548
					if buffer[position] != rune('S') {
						goto l538
					}
					position++
				}
			l548:
				{
					position550, tokenIndex550, depth550 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l551
					}
					position++
					goto l550
				l551:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if buffer[position] != rune('T') {
						goto l538
					}
					position++
				}
			l550:
				{
					position552, tokenIndex552, depth552 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l553
					}
					position++
					goto l552
				l553:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if buffer[position] != rune('A') {
						goto l538
					}
					position++
				}
			l552:
				{
					position554, tokenIndex554, depth554 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l555
					}
					position++
					goto l554
				l555:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
					if buffer[position] != rune('T') {
						goto l538
					}
					position++
				}
			l554:
				{
					position556, tokenIndex556, depth556 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l557
					}
					position++
					goto l556
				l557:
					position, tokenIndex, depth = position556, tokenIndex556, depth556
					if buffer[position] != rune('E') {
						goto l538
					}
					position++
				}
			l556:
				if !_rules[rulesp]() {
					goto l538
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l538
				}
				if !_rules[rulesp]() {
					goto l538
				}
				{
					position558, tokenIndex558, depth558 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l559
					}
					position++
					goto l558
				l559:
					position, tokenIndex, depth = position558, tokenIndex558, depth558
					if buffer[position] != rune('T') {
						goto l538
					}
					position++
				}
			l558:
				{
					position560, tokenIndex560, depth560 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l561
					}
					position++
					goto l560
				l561:
					position, tokenIndex, depth = position560, tokenIndex560, depth560
					if buffer[position] != rune('Y') {
						goto l538
					}
					position++
				}
			l560:
				{
					position562, tokenIndex562, depth562 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l563
					}
					position++
					goto l562
				l563:
					position, tokenIndex, depth = position562, tokenIndex562, depth562
					if buffer[position] != rune('P') {
						goto l538
					}
					position++
				}
			l562:
				{
					position564, tokenIndex564, depth564 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l565
					}
					position++
					goto l564
				l565:
					position, tokenIndex, depth = position564, tokenIndex564, depth564
					if buffer[position] != rune('E') {
						goto l538
					}
					position++
				}
			l564:
				if !_rules[rulesp]() {
					goto l538
				}
				if !_rules[ruleSourceSinkType]() {
					goto l538
				}
				if !_rules[ruleSetOptSpecs]() {
					goto l538
				}
				if !_rules[ruleAction21]() {
					goto l538
				}
				depth--
				add(ruleLoadStateStmt, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 28 LoadStateOrCreateStmt <- <(LoadStateStmt sp (('o' / 'O') ('r' / 'R')) sp (('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp (('i' / 'I') ('f' / 'F')) sp (('n' / 'N') ('o' / 'O') ('t' / 'T')) sp (('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S')) SourceSinkSpecs Action22)> */
		func() bool {
			position566, tokenIndex566, depth566 := position, tokenIndex, depth
			{
				position567 := position
				depth++
				if !_rules[ruleLoadStateStmt]() {
					goto l566
				}
				if !_rules[rulesp]() {
					goto l566
				}
				{
					position568, tokenIndex568, depth568 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l569
					}
					position++
					goto l568
				l569:
					position, tokenIndex, depth = position568, tokenIndex568, depth568
					if buffer[position] != rune('O') {
						goto l566
					}
					position++
				}
			l568:
				{
					position570, tokenIndex570, depth570 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l571
					}
					position++
					goto l570
				l571:
					position, tokenIndex, depth = position570, tokenIndex570, depth570
					if buffer[position] != rune('R') {
						goto l566
					}
					position++
				}
			l570:
				if !_rules[rulesp]() {
					goto l566
				}
				{
					position572, tokenIndex572, depth572 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l573
					}
					position++
					goto l572
				l573:
					position, tokenIndex, depth = position572, tokenIndex572, depth572
					if buffer[position] != rune('C') {
						goto l566
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
						goto l566
					}
					position++
				}
			l574:
				{
					position576, tokenIndex576, depth576 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l577
					}
					position++
					goto l576
				l577:
					position, tokenIndex, depth = position576, tokenIndex576, depth576
					if buffer[position] != rune('E') {
						goto l566
					}
					position++
				}
			l576:
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
						goto l566
					}
					position++
				}
			l578:
				{
					position580, tokenIndex580, depth580 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l581
					}
					position++
					goto l580
				l581:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if buffer[position] != rune('T') {
						goto l566
					}
					position++
				}
			l580:
				{
					position582, tokenIndex582, depth582 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l583
					}
					position++
					goto l582
				l583:
					position, tokenIndex, depth = position582, tokenIndex582, depth582
					if buffer[position] != rune('E') {
						goto l566
					}
					position++
				}
			l582:
				if !_rules[rulesp]() {
					goto l566
				}
				{
					position584, tokenIndex584, depth584 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l585
					}
					position++
					goto l584
				l585:
					position, tokenIndex, depth = position584, tokenIndex584, depth584
					if buffer[position] != rune('I') {
						goto l566
					}
					position++
				}
			l584:
				{
					position586, tokenIndex586, depth586 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l587
					}
					position++
					goto l586
				l587:
					position, tokenIndex, depth = position586, tokenIndex586, depth586
					if buffer[position] != rune('F') {
						goto l566
					}
					position++
				}
			l586:
				if !_rules[rulesp]() {
					goto l566
				}
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
						goto l566
					}
					position++
				}
			l588:
				{
					position590, tokenIndex590, depth590 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l591
					}
					position++
					goto l590
				l591:
					position, tokenIndex, depth = position590, tokenIndex590, depth590
					if buffer[position] != rune('O') {
						goto l566
					}
					position++
				}
			l590:
				{
					position592, tokenIndex592, depth592 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l593
					}
					position++
					goto l592
				l593:
					position, tokenIndex, depth = position592, tokenIndex592, depth592
					if buffer[position] != rune('T') {
						goto l566
					}
					position++
				}
			l592:
				if !_rules[rulesp]() {
					goto l566
				}
				{
					position594, tokenIndex594, depth594 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l595
					}
					position++
					goto l594
				l595:
					position, tokenIndex, depth = position594, tokenIndex594, depth594
					if buffer[position] != rune('E') {
						goto l566
					}
					position++
				}
			l594:
				{
					position596, tokenIndex596, depth596 := position, tokenIndex, depth
					if buffer[position] != rune('x') {
						goto l597
					}
					position++
					goto l596
				l597:
					position, tokenIndex, depth = position596, tokenIndex596, depth596
					if buffer[position] != rune('X') {
						goto l566
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
						goto l566
					}
					position++
				}
			l598:
				{
					position600, tokenIndex600, depth600 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l601
					}
					position++
					goto l600
				l601:
					position, tokenIndex, depth = position600, tokenIndex600, depth600
					if buffer[position] != rune('S') {
						goto l566
					}
					position++
				}
			l600:
				{
					position602, tokenIndex602, depth602 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l603
					}
					position++
					goto l602
				l603:
					position, tokenIndex, depth = position602, tokenIndex602, depth602
					if buffer[position] != rune('T') {
						goto l566
					}
					position++
				}
			l602:
				{
					position604, tokenIndex604, depth604 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l605
					}
					position++
					goto l604
				l605:
					position, tokenIndex, depth = position604, tokenIndex604, depth604
					if buffer[position] != rune('S') {
						goto l566
					}
					position++
				}
			l604:
				if !_rules[ruleSourceSinkSpecs]() {
					goto l566
				}
				if !_rules[ruleAction22]() {
					goto l566
				}
				depth--
				add(ruleLoadStateOrCreateStmt, position567)
			}
			return true
		l566:
			position, tokenIndex, depth = position566, tokenIndex566, depth566
			return false
		},
		/* 29 SaveStateStmt <- <(('s' / 'S') ('a' / 'A') ('v' / 'V') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action23)> */
		func() bool {
			position606, tokenIndex606, depth606 := position, tokenIndex, depth
			{
				position607 := position
				depth++
				{
					position608, tokenIndex608, depth608 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l609
					}
					position++
					goto l608
				l609:
					position, tokenIndex, depth = position608, tokenIndex608, depth608
					if buffer[position] != rune('S') {
						goto l606
					}
					position++
				}
			l608:
				{
					position610, tokenIndex610, depth610 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l611
					}
					position++
					goto l610
				l611:
					position, tokenIndex, depth = position610, tokenIndex610, depth610
					if buffer[position] != rune('A') {
						goto l606
					}
					position++
				}
			l610:
				{
					position612, tokenIndex612, depth612 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l613
					}
					position++
					goto l612
				l613:
					position, tokenIndex, depth = position612, tokenIndex612, depth612
					if buffer[position] != rune('V') {
						goto l606
					}
					position++
				}
			l612:
				{
					position614, tokenIndex614, depth614 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l615
					}
					position++
					goto l614
				l615:
					position, tokenIndex, depth = position614, tokenIndex614, depth614
					if buffer[position] != rune('E') {
						goto l606
					}
					position++
				}
			l614:
				if !_rules[rulesp]() {
					goto l606
				}
				{
					position616, tokenIndex616, depth616 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l617
					}
					position++
					goto l616
				l617:
					position, tokenIndex, depth = position616, tokenIndex616, depth616
					if buffer[position] != rune('S') {
						goto l606
					}
					position++
				}
			l616:
				{
					position618, tokenIndex618, depth618 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l619
					}
					position++
					goto l618
				l619:
					position, tokenIndex, depth = position618, tokenIndex618, depth618
					if buffer[position] != rune('T') {
						goto l606
					}
					position++
				}
			l618:
				{
					position620, tokenIndex620, depth620 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l621
					}
					position++
					goto l620
				l621:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if buffer[position] != rune('A') {
						goto l606
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
						goto l606
					}
					position++
				}
			l622:
				{
					position624, tokenIndex624, depth624 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l625
					}
					position++
					goto l624
				l625:
					position, tokenIndex, depth = position624, tokenIndex624, depth624
					if buffer[position] != rune('E') {
						goto l606
					}
					position++
				}
			l624:
				if !_rules[rulesp]() {
					goto l606
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l606
				}
				if !_rules[ruleAction23]() {
					goto l606
				}
				depth--
				add(ruleSaveStateStmt, position607)
			}
			return true
		l606:
			position, tokenIndex, depth = position606, tokenIndex606, depth606
			return false
		},
		/* 30 Emitter <- <(sp (ISTREAM / DSTREAM / RSTREAM) EmitterOptions Action24)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
				if !_rules[rulesp]() {
					goto l626
				}
				{
					position628, tokenIndex628, depth628 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l629
					}
					goto l628
				l629:
					position, tokenIndex, depth = position628, tokenIndex628, depth628
					if !_rules[ruleDSTREAM]() {
						goto l630
					}
					goto l628
				l630:
					position, tokenIndex, depth = position628, tokenIndex628, depth628
					if !_rules[ruleRSTREAM]() {
						goto l626
					}
				}
			l628:
				if !_rules[ruleEmitterOptions]() {
					goto l626
				}
				if !_rules[ruleAction24]() {
					goto l626
				}
				depth--
				add(ruleEmitter, position627)
			}
			return true
		l626:
			position, tokenIndex, depth = position626, tokenIndex626, depth626
			return false
		},
		/* 31 EmitterOptions <- <(<(spOpt '[' spOpt EmitterOptionCombinations spOpt ']')?> Action25)> */
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
						if !_rules[rulespOpt]() {
							goto l634
						}
						if buffer[position] != rune('[') {
							goto l634
						}
						position++
						if !_rules[rulespOpt]() {
							goto l634
						}
						if !_rules[ruleEmitterOptionCombinations]() {
							goto l634
						}
						if !_rules[rulespOpt]() {
							goto l634
						}
						if buffer[position] != rune(']') {
							goto l634
						}
						position++
						goto l635
					l634:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
					}
				l635:
					depth--
					add(rulePegText, position633)
				}
				if !_rules[ruleAction25]() {
					goto l631
				}
				depth--
				add(ruleEmitterOptions, position632)
			}
			return true
		l631:
			position, tokenIndex, depth = position631, tokenIndex631, depth631
			return false
		},
		/* 32 EmitterOptionCombinations <- <(EmitterLimit / (EmitterSample sp EmitterLimit) / EmitterSample)> */
		func() bool {
			position636, tokenIndex636, depth636 := position, tokenIndex, depth
			{
				position637 := position
				depth++
				{
					position638, tokenIndex638, depth638 := position, tokenIndex, depth
					if !_rules[ruleEmitterLimit]() {
						goto l639
					}
					goto l638
				l639:
					position, tokenIndex, depth = position638, tokenIndex638, depth638
					if !_rules[ruleEmitterSample]() {
						goto l640
					}
					if !_rules[rulesp]() {
						goto l640
					}
					if !_rules[ruleEmitterLimit]() {
						goto l640
					}
					goto l638
				l640:
					position, tokenIndex, depth = position638, tokenIndex638, depth638
					if !_rules[ruleEmitterSample]() {
						goto l636
					}
				}
			l638:
				depth--
				add(ruleEmitterOptionCombinations, position637)
			}
			return true
		l636:
			position, tokenIndex, depth = position636, tokenIndex636, depth636
			return false
		},
		/* 33 EmitterLimit <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') sp NumericLiteral Action26)> */
		func() bool {
			position641, tokenIndex641, depth641 := position, tokenIndex, depth
			{
				position642 := position
				depth++
				{
					position643, tokenIndex643, depth643 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l644
					}
					position++
					goto l643
				l644:
					position, tokenIndex, depth = position643, tokenIndex643, depth643
					if buffer[position] != rune('L') {
						goto l641
					}
					position++
				}
			l643:
				{
					position645, tokenIndex645, depth645 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l646
					}
					position++
					goto l645
				l646:
					position, tokenIndex, depth = position645, tokenIndex645, depth645
					if buffer[position] != rune('I') {
						goto l641
					}
					position++
				}
			l645:
				{
					position647, tokenIndex647, depth647 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l648
					}
					position++
					goto l647
				l648:
					position, tokenIndex, depth = position647, tokenIndex647, depth647
					if buffer[position] != rune('M') {
						goto l641
					}
					position++
				}
			l647:
				{
					position649, tokenIndex649, depth649 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l650
					}
					position++
					goto l649
				l650:
					position, tokenIndex, depth = position649, tokenIndex649, depth649
					if buffer[position] != rune('I') {
						goto l641
					}
					position++
				}
			l649:
				{
					position651, tokenIndex651, depth651 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l652
					}
					position++
					goto l651
				l652:
					position, tokenIndex, depth = position651, tokenIndex651, depth651
					if buffer[position] != rune('T') {
						goto l641
					}
					position++
				}
			l651:
				if !_rules[rulesp]() {
					goto l641
				}
				if !_rules[ruleNumericLiteral]() {
					goto l641
				}
				if !_rules[ruleAction26]() {
					goto l641
				}
				depth--
				add(ruleEmitterLimit, position642)
			}
			return true
		l641:
			position, tokenIndex, depth = position641, tokenIndex641, depth641
			return false
		},
		/* 34 EmitterSample <- <(CountBasedSampling / RandomizedSampling / TimeBasedSampling)> */
		func() bool {
			position653, tokenIndex653, depth653 := position, tokenIndex, depth
			{
				position654 := position
				depth++
				{
					position655, tokenIndex655, depth655 := position, tokenIndex, depth
					if !_rules[ruleCountBasedSampling]() {
						goto l656
					}
					goto l655
				l656:
					position, tokenIndex, depth = position655, tokenIndex655, depth655
					if !_rules[ruleRandomizedSampling]() {
						goto l657
					}
					goto l655
				l657:
					position, tokenIndex, depth = position655, tokenIndex655, depth655
					if !_rules[ruleTimeBasedSampling]() {
						goto l653
					}
				}
			l655:
				depth--
				add(ruleEmitterSample, position654)
			}
			return true
		l653:
			position, tokenIndex, depth = position653, tokenIndex653, depth653
			return false
		},
		/* 35 CountBasedSampling <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp NumericLiteral spOpt '-'? spOpt ((('s' / 'S') ('t' / 'T')) / (('n' / 'N') ('d' / 'D')) / (('r' / 'R') ('d' / 'D')) / (('t' / 'T') ('h' / 'H'))) sp (('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E')) Action27)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				{
					position660, tokenIndex660, depth660 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l661
					}
					position++
					goto l660
				l661:
					position, tokenIndex, depth = position660, tokenIndex660, depth660
					if buffer[position] != rune('E') {
						goto l658
					}
					position++
				}
			l660:
				{
					position662, tokenIndex662, depth662 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l663
					}
					position++
					goto l662
				l663:
					position, tokenIndex, depth = position662, tokenIndex662, depth662
					if buffer[position] != rune('V') {
						goto l658
					}
					position++
				}
			l662:
				{
					position664, tokenIndex664, depth664 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l665
					}
					position++
					goto l664
				l665:
					position, tokenIndex, depth = position664, tokenIndex664, depth664
					if buffer[position] != rune('E') {
						goto l658
					}
					position++
				}
			l664:
				{
					position666, tokenIndex666, depth666 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l667
					}
					position++
					goto l666
				l667:
					position, tokenIndex, depth = position666, tokenIndex666, depth666
					if buffer[position] != rune('R') {
						goto l658
					}
					position++
				}
			l666:
				{
					position668, tokenIndex668, depth668 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l669
					}
					position++
					goto l668
				l669:
					position, tokenIndex, depth = position668, tokenIndex668, depth668
					if buffer[position] != rune('Y') {
						goto l658
					}
					position++
				}
			l668:
				if !_rules[rulesp]() {
					goto l658
				}
				if !_rules[ruleNumericLiteral]() {
					goto l658
				}
				if !_rules[rulespOpt]() {
					goto l658
				}
				{
					position670, tokenIndex670, depth670 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l670
					}
					position++
					goto l671
				l670:
					position, tokenIndex, depth = position670, tokenIndex670, depth670
				}
			l671:
				if !_rules[rulespOpt]() {
					goto l658
				}
				{
					position672, tokenIndex672, depth672 := position, tokenIndex, depth
					{
						position674, tokenIndex674, depth674 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l675
						}
						position++
						goto l674
					l675:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
						if buffer[position] != rune('S') {
							goto l673
						}
						position++
					}
				l674:
					{
						position676, tokenIndex676, depth676 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l677
						}
						position++
						goto l676
					l677:
						position, tokenIndex, depth = position676, tokenIndex676, depth676
						if buffer[position] != rune('T') {
							goto l673
						}
						position++
					}
				l676:
					goto l672
				l673:
					position, tokenIndex, depth = position672, tokenIndex672, depth672
					{
						position679, tokenIndex679, depth679 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l680
						}
						position++
						goto l679
					l680:
						position, tokenIndex, depth = position679, tokenIndex679, depth679
						if buffer[position] != rune('N') {
							goto l678
						}
						position++
					}
				l679:
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if buffer[position] != rune('D') {
							goto l678
						}
						position++
					}
				l681:
					goto l672
				l678:
					position, tokenIndex, depth = position672, tokenIndex672, depth672
					{
						position684, tokenIndex684, depth684 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l685
						}
						position++
						goto l684
					l685:
						position, tokenIndex, depth = position684, tokenIndex684, depth684
						if buffer[position] != rune('R') {
							goto l683
						}
						position++
					}
				l684:
					{
						position686, tokenIndex686, depth686 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
						if buffer[position] != rune('D') {
							goto l683
						}
						position++
					}
				l686:
					goto l672
				l683:
					position, tokenIndex, depth = position672, tokenIndex672, depth672
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if buffer[position] != rune('T') {
							goto l658
						}
						position++
					}
				l688:
					{
						position690, tokenIndex690, depth690 := position, tokenIndex, depth
						if buffer[position] != rune('h') {
							goto l691
						}
						position++
						goto l690
					l691:
						position, tokenIndex, depth = position690, tokenIndex690, depth690
						if buffer[position] != rune('H') {
							goto l658
						}
						position++
					}
				l690:
				}
			l672:
				if !_rules[rulesp]() {
					goto l658
				}
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
						goto l658
					}
					position++
				}
			l692:
				{
					position694, tokenIndex694, depth694 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l695
					}
					position++
					goto l694
				l695:
					position, tokenIndex, depth = position694, tokenIndex694, depth694
					if buffer[position] != rune('U') {
						goto l658
					}
					position++
				}
			l694:
				{
					position696, tokenIndex696, depth696 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l697
					}
					position++
					goto l696
				l697:
					position, tokenIndex, depth = position696, tokenIndex696, depth696
					if buffer[position] != rune('P') {
						goto l658
					}
					position++
				}
			l696:
				{
					position698, tokenIndex698, depth698 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l699
					}
					position++
					goto l698
				l699:
					position, tokenIndex, depth = position698, tokenIndex698, depth698
					if buffer[position] != rune('L') {
						goto l658
					}
					position++
				}
			l698:
				{
					position700, tokenIndex700, depth700 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l701
					}
					position++
					goto l700
				l701:
					position, tokenIndex, depth = position700, tokenIndex700, depth700
					if buffer[position] != rune('E') {
						goto l658
					}
					position++
				}
			l700:
				if !_rules[ruleAction27]() {
					goto l658
				}
				depth--
				add(ruleCountBasedSampling, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 36 RandomizedSampling <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') sp (FloatLiteral / NumericLiteral) spOpt '%' Action28)> */
		func() bool {
			position702, tokenIndex702, depth702 := position, tokenIndex, depth
			{
				position703 := position
				depth++
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
						goto l702
					}
					position++
				}
			l704:
				{
					position706, tokenIndex706, depth706 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l707
					}
					position++
					goto l706
				l707:
					position, tokenIndex, depth = position706, tokenIndex706, depth706
					if buffer[position] != rune('A') {
						goto l702
					}
					position++
				}
			l706:
				{
					position708, tokenIndex708, depth708 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l709
					}
					position++
					goto l708
				l709:
					position, tokenIndex, depth = position708, tokenIndex708, depth708
					if buffer[position] != rune('M') {
						goto l702
					}
					position++
				}
			l708:
				{
					position710, tokenIndex710, depth710 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l711
					}
					position++
					goto l710
				l711:
					position, tokenIndex, depth = position710, tokenIndex710, depth710
					if buffer[position] != rune('P') {
						goto l702
					}
					position++
				}
			l710:
				{
					position712, tokenIndex712, depth712 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l713
					}
					position++
					goto l712
				l713:
					position, tokenIndex, depth = position712, tokenIndex712, depth712
					if buffer[position] != rune('L') {
						goto l702
					}
					position++
				}
			l712:
				{
					position714, tokenIndex714, depth714 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l715
					}
					position++
					goto l714
				l715:
					position, tokenIndex, depth = position714, tokenIndex714, depth714
					if buffer[position] != rune('E') {
						goto l702
					}
					position++
				}
			l714:
				if !_rules[rulesp]() {
					goto l702
				}
				{
					position716, tokenIndex716, depth716 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l717
					}
					goto l716
				l717:
					position, tokenIndex, depth = position716, tokenIndex716, depth716
					if !_rules[ruleNumericLiteral]() {
						goto l702
					}
				}
			l716:
				if !_rules[rulespOpt]() {
					goto l702
				}
				if buffer[position] != rune('%') {
					goto l702
				}
				position++
				if !_rules[ruleAction28]() {
					goto l702
				}
				depth--
				add(ruleRandomizedSampling, position703)
			}
			return true
		l702:
			position, tokenIndex, depth = position702, tokenIndex702, depth702
			return false
		},
		/* 37 TimeBasedSampling <- <(TimeBasedSamplingSeconds / TimeBasedSamplingMilliseconds)> */
		func() bool {
			position718, tokenIndex718, depth718 := position, tokenIndex, depth
			{
				position719 := position
				depth++
				{
					position720, tokenIndex720, depth720 := position, tokenIndex, depth
					if !_rules[ruleTimeBasedSamplingSeconds]() {
						goto l721
					}
					goto l720
				l721:
					position, tokenIndex, depth = position720, tokenIndex720, depth720
					if !_rules[ruleTimeBasedSamplingMilliseconds]() {
						goto l718
					}
				}
			l720:
				depth--
				add(ruleTimeBasedSampling, position719)
			}
			return true
		l718:
			position, tokenIndex, depth = position718, tokenIndex718, depth718
			return false
		},
		/* 38 TimeBasedSamplingSeconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action29)> */
		func() bool {
			position722, tokenIndex722, depth722 := position, tokenIndex, depth
			{
				position723 := position
				depth++
				{
					position724, tokenIndex724, depth724 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l725
					}
					position++
					goto l724
				l725:
					position, tokenIndex, depth = position724, tokenIndex724, depth724
					if buffer[position] != rune('E') {
						goto l722
					}
					position++
				}
			l724:
				{
					position726, tokenIndex726, depth726 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l727
					}
					position++
					goto l726
				l727:
					position, tokenIndex, depth = position726, tokenIndex726, depth726
					if buffer[position] != rune('V') {
						goto l722
					}
					position++
				}
			l726:
				{
					position728, tokenIndex728, depth728 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l729
					}
					position++
					goto l728
				l729:
					position, tokenIndex, depth = position728, tokenIndex728, depth728
					if buffer[position] != rune('E') {
						goto l722
					}
					position++
				}
			l728:
				{
					position730, tokenIndex730, depth730 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l731
					}
					position++
					goto l730
				l731:
					position, tokenIndex, depth = position730, tokenIndex730, depth730
					if buffer[position] != rune('R') {
						goto l722
					}
					position++
				}
			l730:
				{
					position732, tokenIndex732, depth732 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l733
					}
					position++
					goto l732
				l733:
					position, tokenIndex, depth = position732, tokenIndex732, depth732
					if buffer[position] != rune('Y') {
						goto l722
					}
					position++
				}
			l732:
				if !_rules[rulesp]() {
					goto l722
				}
				{
					position734, tokenIndex734, depth734 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l735
					}
					goto l734
				l735:
					position, tokenIndex, depth = position734, tokenIndex734, depth734
					if !_rules[ruleNumericLiteral]() {
						goto l722
					}
				}
			l734:
				if !_rules[rulesp]() {
					goto l722
				}
				{
					position736, tokenIndex736, depth736 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l737
					}
					position++
					goto l736
				l737:
					position, tokenIndex, depth = position736, tokenIndex736, depth736
					if buffer[position] != rune('S') {
						goto l722
					}
					position++
				}
			l736:
				{
					position738, tokenIndex738, depth738 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l739
					}
					position++
					goto l738
				l739:
					position, tokenIndex, depth = position738, tokenIndex738, depth738
					if buffer[position] != rune('E') {
						goto l722
					}
					position++
				}
			l738:
				{
					position740, tokenIndex740, depth740 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l741
					}
					position++
					goto l740
				l741:
					position, tokenIndex, depth = position740, tokenIndex740, depth740
					if buffer[position] != rune('C') {
						goto l722
					}
					position++
				}
			l740:
				{
					position742, tokenIndex742, depth742 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l743
					}
					position++
					goto l742
				l743:
					position, tokenIndex, depth = position742, tokenIndex742, depth742
					if buffer[position] != rune('O') {
						goto l722
					}
					position++
				}
			l742:
				{
					position744, tokenIndex744, depth744 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l745
					}
					position++
					goto l744
				l745:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
					if buffer[position] != rune('N') {
						goto l722
					}
					position++
				}
			l744:
				{
					position746, tokenIndex746, depth746 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l747
					}
					position++
					goto l746
				l747:
					position, tokenIndex, depth = position746, tokenIndex746, depth746
					if buffer[position] != rune('D') {
						goto l722
					}
					position++
				}
			l746:
				{
					position748, tokenIndex748, depth748 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l749
					}
					position++
					goto l748
				l749:
					position, tokenIndex, depth = position748, tokenIndex748, depth748
					if buffer[position] != rune('S') {
						goto l722
					}
					position++
				}
			l748:
				if !_rules[ruleAction29]() {
					goto l722
				}
				depth--
				add(ruleTimeBasedSamplingSeconds, position723)
			}
			return true
		l722:
			position, tokenIndex, depth = position722, tokenIndex722, depth722
			return false
		},
		/* 39 TimeBasedSamplingMilliseconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action30)> */
		func() bool {
			position750, tokenIndex750, depth750 := position, tokenIndex, depth
			{
				position751 := position
				depth++
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
						goto l750
					}
					position++
				}
			l752:
				{
					position754, tokenIndex754, depth754 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l755
					}
					position++
					goto l754
				l755:
					position, tokenIndex, depth = position754, tokenIndex754, depth754
					if buffer[position] != rune('V') {
						goto l750
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
						goto l750
					}
					position++
				}
			l756:
				{
					position758, tokenIndex758, depth758 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l759
					}
					position++
					goto l758
				l759:
					position, tokenIndex, depth = position758, tokenIndex758, depth758
					if buffer[position] != rune('R') {
						goto l750
					}
					position++
				}
			l758:
				{
					position760, tokenIndex760, depth760 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l761
					}
					position++
					goto l760
				l761:
					position, tokenIndex, depth = position760, tokenIndex760, depth760
					if buffer[position] != rune('Y') {
						goto l750
					}
					position++
				}
			l760:
				if !_rules[rulesp]() {
					goto l750
				}
				{
					position762, tokenIndex762, depth762 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l763
					}
					goto l762
				l763:
					position, tokenIndex, depth = position762, tokenIndex762, depth762
					if !_rules[ruleNumericLiteral]() {
						goto l750
					}
				}
			l762:
				if !_rules[rulesp]() {
					goto l750
				}
				{
					position764, tokenIndex764, depth764 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l765
					}
					position++
					goto l764
				l765:
					position, tokenIndex, depth = position764, tokenIndex764, depth764
					if buffer[position] != rune('M') {
						goto l750
					}
					position++
				}
			l764:
				{
					position766, tokenIndex766, depth766 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l767
					}
					position++
					goto l766
				l767:
					position, tokenIndex, depth = position766, tokenIndex766, depth766
					if buffer[position] != rune('I') {
						goto l750
					}
					position++
				}
			l766:
				{
					position768, tokenIndex768, depth768 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l769
					}
					position++
					goto l768
				l769:
					position, tokenIndex, depth = position768, tokenIndex768, depth768
					if buffer[position] != rune('L') {
						goto l750
					}
					position++
				}
			l768:
				{
					position770, tokenIndex770, depth770 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l771
					}
					position++
					goto l770
				l771:
					position, tokenIndex, depth = position770, tokenIndex770, depth770
					if buffer[position] != rune('L') {
						goto l750
					}
					position++
				}
			l770:
				{
					position772, tokenIndex772, depth772 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l773
					}
					position++
					goto l772
				l773:
					position, tokenIndex, depth = position772, tokenIndex772, depth772
					if buffer[position] != rune('I') {
						goto l750
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
						goto l750
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
						goto l750
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
						goto l750
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
						goto l750
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
						goto l750
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
						goto l750
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
						goto l750
					}
					position++
				}
			l786:
				if !_rules[ruleAction30]() {
					goto l750
				}
				depth--
				add(ruleTimeBasedSamplingMilliseconds, position751)
			}
			return true
		l750:
			position, tokenIndex, depth = position750, tokenIndex750, depth750
			return false
		},
		/* 40 Projections <- <(<(sp Projection (spOpt ',' spOpt Projection)*)> Action31)> */
		func() bool {
			position788, tokenIndex788, depth788 := position, tokenIndex, depth
			{
				position789 := position
				depth++
				{
					position790 := position
					depth++
					if !_rules[rulesp]() {
						goto l788
					}
					if !_rules[ruleProjection]() {
						goto l788
					}
				l791:
					{
						position792, tokenIndex792, depth792 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l792
						}
						if buffer[position] != rune(',') {
							goto l792
						}
						position++
						if !_rules[rulespOpt]() {
							goto l792
						}
						if !_rules[ruleProjection]() {
							goto l792
						}
						goto l791
					l792:
						position, tokenIndex, depth = position792, tokenIndex792, depth792
					}
					depth--
					add(rulePegText, position790)
				}
				if !_rules[ruleAction31]() {
					goto l788
				}
				depth--
				add(ruleProjections, position789)
			}
			return true
		l788:
			position, tokenIndex, depth = position788, tokenIndex788, depth788
			return false
		},
		/* 41 Projection <- <(AliasExpression / ExpressionOrWildcard)> */
		func() bool {
			position793, tokenIndex793, depth793 := position, tokenIndex, depth
			{
				position794 := position
				depth++
				{
					position795, tokenIndex795, depth795 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l796
					}
					goto l795
				l796:
					position, tokenIndex, depth = position795, tokenIndex795, depth795
					if !_rules[ruleExpressionOrWildcard]() {
						goto l793
					}
				}
			l795:
				depth--
				add(ruleProjection, position794)
			}
			return true
		l793:
			position, tokenIndex, depth = position793, tokenIndex793, depth793
			return false
		},
		/* 42 AliasExpression <- <(ExpressionOrWildcard sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action32)> */
		func() bool {
			position797, tokenIndex797, depth797 := position, tokenIndex, depth
			{
				position798 := position
				depth++
				if !_rules[ruleExpressionOrWildcard]() {
					goto l797
				}
				if !_rules[rulesp]() {
					goto l797
				}
				{
					position799, tokenIndex799, depth799 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l800
					}
					position++
					goto l799
				l800:
					position, tokenIndex, depth = position799, tokenIndex799, depth799
					if buffer[position] != rune('A') {
						goto l797
					}
					position++
				}
			l799:
				{
					position801, tokenIndex801, depth801 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l802
					}
					position++
					goto l801
				l802:
					position, tokenIndex, depth = position801, tokenIndex801, depth801
					if buffer[position] != rune('S') {
						goto l797
					}
					position++
				}
			l801:
				if !_rules[rulesp]() {
					goto l797
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l797
				}
				if !_rules[ruleAction32]() {
					goto l797
				}
				depth--
				add(ruleAliasExpression, position798)
			}
			return true
		l797:
			position, tokenIndex, depth = position797, tokenIndex797, depth797
			return false
		},
		/* 43 WindowedFrom <- <(<(sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp Relations)?> Action33)> */
		func() bool {
			position803, tokenIndex803, depth803 := position, tokenIndex, depth
			{
				position804 := position
				depth++
				{
					position805 := position
					depth++
					{
						position806, tokenIndex806, depth806 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l806
						}
						{
							position808, tokenIndex808, depth808 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l809
							}
							position++
							goto l808
						l809:
							position, tokenIndex, depth = position808, tokenIndex808, depth808
							if buffer[position] != rune('F') {
								goto l806
							}
							position++
						}
					l808:
						{
							position810, tokenIndex810, depth810 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l811
							}
							position++
							goto l810
						l811:
							position, tokenIndex, depth = position810, tokenIndex810, depth810
							if buffer[position] != rune('R') {
								goto l806
							}
							position++
						}
					l810:
						{
							position812, tokenIndex812, depth812 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l813
							}
							position++
							goto l812
						l813:
							position, tokenIndex, depth = position812, tokenIndex812, depth812
							if buffer[position] != rune('O') {
								goto l806
							}
							position++
						}
					l812:
						{
							position814, tokenIndex814, depth814 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l815
							}
							position++
							goto l814
						l815:
							position, tokenIndex, depth = position814, tokenIndex814, depth814
							if buffer[position] != rune('M') {
								goto l806
							}
							position++
						}
					l814:
						if !_rules[rulesp]() {
							goto l806
						}
						if !_rules[ruleRelations]() {
							goto l806
						}
						goto l807
					l806:
						position, tokenIndex, depth = position806, tokenIndex806, depth806
					}
				l807:
					depth--
					add(rulePegText, position805)
				}
				if !_rules[ruleAction33]() {
					goto l803
				}
				depth--
				add(ruleWindowedFrom, position804)
			}
			return true
		l803:
			position, tokenIndex, depth = position803, tokenIndex803, depth803
			return false
		},
		/* 44 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position816, tokenIndex816, depth816 := position, tokenIndex, depth
			{
				position817 := position
				depth++
				{
					position818, tokenIndex818, depth818 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l819
					}
					goto l818
				l819:
					position, tokenIndex, depth = position818, tokenIndex818, depth818
					if !_rules[ruleTuplesInterval]() {
						goto l816
					}
				}
			l818:
				depth--
				add(ruleInterval, position817)
			}
			return true
		l816:
			position, tokenIndex, depth = position816, tokenIndex816, depth816
			return false
		},
		/* 45 TimeInterval <- <((FloatLiteral / NumericLiteral) sp (SECONDS / MILLISECONDS) Action34)> */
		func() bool {
			position820, tokenIndex820, depth820 := position, tokenIndex, depth
			{
				position821 := position
				depth++
				{
					position822, tokenIndex822, depth822 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l823
					}
					goto l822
				l823:
					position, tokenIndex, depth = position822, tokenIndex822, depth822
					if !_rules[ruleNumericLiteral]() {
						goto l820
					}
				}
			l822:
				if !_rules[rulesp]() {
					goto l820
				}
				{
					position824, tokenIndex824, depth824 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l825
					}
					goto l824
				l825:
					position, tokenIndex, depth = position824, tokenIndex824, depth824
					if !_rules[ruleMILLISECONDS]() {
						goto l820
					}
				}
			l824:
				if !_rules[ruleAction34]() {
					goto l820
				}
				depth--
				add(ruleTimeInterval, position821)
			}
			return true
		l820:
			position, tokenIndex, depth = position820, tokenIndex820, depth820
			return false
		},
		/* 46 TuplesInterval <- <(NumericLiteral sp TUPLES Action35)> */
		func() bool {
			position826, tokenIndex826, depth826 := position, tokenIndex, depth
			{
				position827 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l826
				}
				if !_rules[rulesp]() {
					goto l826
				}
				if !_rules[ruleTUPLES]() {
					goto l826
				}
				if !_rules[ruleAction35]() {
					goto l826
				}
				depth--
				add(ruleTuplesInterval, position827)
			}
			return true
		l826:
			position, tokenIndex, depth = position826, tokenIndex826, depth826
			return false
		},
		/* 47 Relations <- <(RelationLike (spOpt ',' spOpt RelationLike)*)> */
		func() bool {
			position828, tokenIndex828, depth828 := position, tokenIndex, depth
			{
				position829 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l828
				}
			l830:
				{
					position831, tokenIndex831, depth831 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l831
					}
					if buffer[position] != rune(',') {
						goto l831
					}
					position++
					if !_rules[rulespOpt]() {
						goto l831
					}
					if !_rules[ruleRelationLike]() {
						goto l831
					}
					goto l830
				l831:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
				}
				depth--
				add(ruleRelations, position829)
			}
			return true
		l828:
			position, tokenIndex, depth = position828, tokenIndex828, depth828
			return false
		},
		/* 48 Filter <- <(<(sp (('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E')) sp Expression)?> Action36)> */
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
						if !_rules[rulesp]() {
							goto l835
						}
						{
							position837, tokenIndex837, depth837 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l838
							}
							position++
							goto l837
						l838:
							position, tokenIndex, depth = position837, tokenIndex837, depth837
							if buffer[position] != rune('W') {
								goto l835
							}
							position++
						}
					l837:
						{
							position839, tokenIndex839, depth839 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l840
							}
							position++
							goto l839
						l840:
							position, tokenIndex, depth = position839, tokenIndex839, depth839
							if buffer[position] != rune('H') {
								goto l835
							}
							position++
						}
					l839:
						{
							position841, tokenIndex841, depth841 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l842
							}
							position++
							goto l841
						l842:
							position, tokenIndex, depth = position841, tokenIndex841, depth841
							if buffer[position] != rune('E') {
								goto l835
							}
							position++
						}
					l841:
						{
							position843, tokenIndex843, depth843 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l844
							}
							position++
							goto l843
						l844:
							position, tokenIndex, depth = position843, tokenIndex843, depth843
							if buffer[position] != rune('R') {
								goto l835
							}
							position++
						}
					l843:
						{
							position845, tokenIndex845, depth845 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l846
							}
							position++
							goto l845
						l846:
							position, tokenIndex, depth = position845, tokenIndex845, depth845
							if buffer[position] != rune('E') {
								goto l835
							}
							position++
						}
					l845:
						if !_rules[rulesp]() {
							goto l835
						}
						if !_rules[ruleExpression]() {
							goto l835
						}
						goto l836
					l835:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
					}
				l836:
					depth--
					add(rulePegText, position834)
				}
				if !_rules[ruleAction36]() {
					goto l832
				}
				depth--
				add(ruleFilter, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 49 Grouping <- <(<(sp (('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P')) sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action37)> */
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
						if !_rules[rulesp]() {
							goto l850
						}
						{
							position852, tokenIndex852, depth852 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l853
							}
							position++
							goto l852
						l853:
							position, tokenIndex, depth = position852, tokenIndex852, depth852
							if buffer[position] != rune('G') {
								goto l850
							}
							position++
						}
					l852:
						{
							position854, tokenIndex854, depth854 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l855
							}
							position++
							goto l854
						l855:
							position, tokenIndex, depth = position854, tokenIndex854, depth854
							if buffer[position] != rune('R') {
								goto l850
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
								goto l850
							}
							position++
						}
					l856:
						{
							position858, tokenIndex858, depth858 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l859
							}
							position++
							goto l858
						l859:
							position, tokenIndex, depth = position858, tokenIndex858, depth858
							if buffer[position] != rune('U') {
								goto l850
							}
							position++
						}
					l858:
						{
							position860, tokenIndex860, depth860 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l861
							}
							position++
							goto l860
						l861:
							position, tokenIndex, depth = position860, tokenIndex860, depth860
							if buffer[position] != rune('P') {
								goto l850
							}
							position++
						}
					l860:
						if !_rules[rulesp]() {
							goto l850
						}
						{
							position862, tokenIndex862, depth862 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l863
							}
							position++
							goto l862
						l863:
							position, tokenIndex, depth = position862, tokenIndex862, depth862
							if buffer[position] != rune('B') {
								goto l850
							}
							position++
						}
					l862:
						{
							position864, tokenIndex864, depth864 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l865
							}
							position++
							goto l864
						l865:
							position, tokenIndex, depth = position864, tokenIndex864, depth864
							if buffer[position] != rune('Y') {
								goto l850
							}
							position++
						}
					l864:
						if !_rules[rulesp]() {
							goto l850
						}
						if !_rules[ruleGroupList]() {
							goto l850
						}
						goto l851
					l850:
						position, tokenIndex, depth = position850, tokenIndex850, depth850
					}
				l851:
					depth--
					add(rulePegText, position849)
				}
				if !_rules[ruleAction37]() {
					goto l847
				}
				depth--
				add(ruleGrouping, position848)
			}
			return true
		l847:
			position, tokenIndex, depth = position847, tokenIndex847, depth847
			return false
		},
		/* 50 GroupList <- <(Expression (spOpt ',' spOpt Expression)*)> */
		func() bool {
			position866, tokenIndex866, depth866 := position, tokenIndex, depth
			{
				position867 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l866
				}
			l868:
				{
					position869, tokenIndex869, depth869 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l869
					}
					if buffer[position] != rune(',') {
						goto l869
					}
					position++
					if !_rules[rulespOpt]() {
						goto l869
					}
					if !_rules[ruleExpression]() {
						goto l869
					}
					goto l868
				l869:
					position, tokenIndex, depth = position869, tokenIndex869, depth869
				}
				depth--
				add(ruleGroupList, position867)
			}
			return true
		l866:
			position, tokenIndex, depth = position866, tokenIndex866, depth866
			return false
		},
		/* 51 Having <- <(<(sp (('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G')) sp Expression)?> Action38)> */
		func() bool {
			position870, tokenIndex870, depth870 := position, tokenIndex, depth
			{
				position871 := position
				depth++
				{
					position872 := position
					depth++
					{
						position873, tokenIndex873, depth873 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l873
						}
						{
							position875, tokenIndex875, depth875 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l876
							}
							position++
							goto l875
						l876:
							position, tokenIndex, depth = position875, tokenIndex875, depth875
							if buffer[position] != rune('H') {
								goto l873
							}
							position++
						}
					l875:
						{
							position877, tokenIndex877, depth877 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l878
							}
							position++
							goto l877
						l878:
							position, tokenIndex, depth = position877, tokenIndex877, depth877
							if buffer[position] != rune('A') {
								goto l873
							}
							position++
						}
					l877:
						{
							position879, tokenIndex879, depth879 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l880
							}
							position++
							goto l879
						l880:
							position, tokenIndex, depth = position879, tokenIndex879, depth879
							if buffer[position] != rune('V') {
								goto l873
							}
							position++
						}
					l879:
						{
							position881, tokenIndex881, depth881 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l882
							}
							position++
							goto l881
						l882:
							position, tokenIndex, depth = position881, tokenIndex881, depth881
							if buffer[position] != rune('I') {
								goto l873
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
								goto l873
							}
							position++
						}
					l883:
						{
							position885, tokenIndex885, depth885 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l886
							}
							position++
							goto l885
						l886:
							position, tokenIndex, depth = position885, tokenIndex885, depth885
							if buffer[position] != rune('G') {
								goto l873
							}
							position++
						}
					l885:
						if !_rules[rulesp]() {
							goto l873
						}
						if !_rules[ruleExpression]() {
							goto l873
						}
						goto l874
					l873:
						position, tokenIndex, depth = position873, tokenIndex873, depth873
					}
				l874:
					depth--
					add(rulePegText, position872)
				}
				if !_rules[ruleAction38]() {
					goto l870
				}
				depth--
				add(ruleHaving, position871)
			}
			return true
		l870:
			position, tokenIndex, depth = position870, tokenIndex870, depth870
			return false
		},
		/* 52 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action39))> */
		func() bool {
			position887, tokenIndex887, depth887 := position, tokenIndex, depth
			{
				position888 := position
				depth++
				{
					position889, tokenIndex889, depth889 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l890
					}
					goto l889
				l890:
					position, tokenIndex, depth = position889, tokenIndex889, depth889
					if !_rules[ruleStreamWindow]() {
						goto l887
					}
					if !_rules[ruleAction39]() {
						goto l887
					}
				}
			l889:
				depth--
				add(ruleRelationLike, position888)
			}
			return true
		l887:
			position, tokenIndex, depth = position887, tokenIndex887, depth887
			return false
		},
		/* 53 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action40)> */
		func() bool {
			position891, tokenIndex891, depth891 := position, tokenIndex, depth
			{
				position892 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l891
				}
				if !_rules[rulesp]() {
					goto l891
				}
				{
					position893, tokenIndex893, depth893 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l894
					}
					position++
					goto l893
				l894:
					position, tokenIndex, depth = position893, tokenIndex893, depth893
					if buffer[position] != rune('A') {
						goto l891
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
						goto l891
					}
					position++
				}
			l895:
				if !_rules[rulesp]() {
					goto l891
				}
				if !_rules[ruleIdentifier]() {
					goto l891
				}
				if !_rules[ruleAction40]() {
					goto l891
				}
				depth--
				add(ruleAliasedStreamWindow, position892)
			}
			return true
		l891:
			position, tokenIndex, depth = position891, tokenIndex891, depth891
			return false
		},
		/* 54 StreamWindow <- <(StreamLike spOpt '[' spOpt (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval spOpt ']' Action41)> */
		func() bool {
			position897, tokenIndex897, depth897 := position, tokenIndex, depth
			{
				position898 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l897
				}
				if !_rules[rulespOpt]() {
					goto l897
				}
				if buffer[position] != rune('[') {
					goto l897
				}
				position++
				if !_rules[rulespOpt]() {
					goto l897
				}
				{
					position899, tokenIndex899, depth899 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l900
					}
					position++
					goto l899
				l900:
					position, tokenIndex, depth = position899, tokenIndex899, depth899
					if buffer[position] != rune('R') {
						goto l897
					}
					position++
				}
			l899:
				{
					position901, tokenIndex901, depth901 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l902
					}
					position++
					goto l901
				l902:
					position, tokenIndex, depth = position901, tokenIndex901, depth901
					if buffer[position] != rune('A') {
						goto l897
					}
					position++
				}
			l901:
				{
					position903, tokenIndex903, depth903 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l904
					}
					position++
					goto l903
				l904:
					position, tokenIndex, depth = position903, tokenIndex903, depth903
					if buffer[position] != rune('N') {
						goto l897
					}
					position++
				}
			l903:
				{
					position905, tokenIndex905, depth905 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l906
					}
					position++
					goto l905
				l906:
					position, tokenIndex, depth = position905, tokenIndex905, depth905
					if buffer[position] != rune('G') {
						goto l897
					}
					position++
				}
			l905:
				{
					position907, tokenIndex907, depth907 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l908
					}
					position++
					goto l907
				l908:
					position, tokenIndex, depth = position907, tokenIndex907, depth907
					if buffer[position] != rune('E') {
						goto l897
					}
					position++
				}
			l907:
				if !_rules[rulesp]() {
					goto l897
				}
				if !_rules[ruleInterval]() {
					goto l897
				}
				if !_rules[rulespOpt]() {
					goto l897
				}
				if buffer[position] != rune(']') {
					goto l897
				}
				position++
				if !_rules[ruleAction41]() {
					goto l897
				}
				depth--
				add(ruleStreamWindow, position898)
			}
			return true
		l897:
			position, tokenIndex, depth = position897, tokenIndex897, depth897
			return false
		},
		/* 55 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position909, tokenIndex909, depth909 := position, tokenIndex, depth
			{
				position910 := position
				depth++
				{
					position911, tokenIndex911, depth911 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l912
					}
					goto l911
				l912:
					position, tokenIndex, depth = position911, tokenIndex911, depth911
					if !_rules[ruleStream]() {
						goto l909
					}
				}
			l911:
				depth--
				add(ruleStreamLike, position910)
			}
			return true
		l909:
			position, tokenIndex, depth = position909, tokenIndex909, depth909
			return false
		},
		/* 56 UDSFFuncApp <- <(FuncAppWithoutOrderBy Action42)> */
		func() bool {
			position913, tokenIndex913, depth913 := position, tokenIndex, depth
			{
				position914 := position
				depth++
				if !_rules[ruleFuncAppWithoutOrderBy]() {
					goto l913
				}
				if !_rules[ruleAction42]() {
					goto l913
				}
				depth--
				add(ruleUDSFFuncApp, position914)
			}
			return true
		l913:
			position, tokenIndex, depth = position913, tokenIndex913, depth913
			return false
		},
		/* 57 SourceSinkSpecs <- <(<(sp (('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action43)> */
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
						if !_rules[rulesp]() {
							goto l918
						}
						{
							position920, tokenIndex920, depth920 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l921
							}
							position++
							goto l920
						l921:
							position, tokenIndex, depth = position920, tokenIndex920, depth920
							if buffer[position] != rune('W') {
								goto l918
							}
							position++
						}
					l920:
						{
							position922, tokenIndex922, depth922 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l923
							}
							position++
							goto l922
						l923:
							position, tokenIndex, depth = position922, tokenIndex922, depth922
							if buffer[position] != rune('I') {
								goto l918
							}
							position++
						}
					l922:
						{
							position924, tokenIndex924, depth924 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l925
							}
							position++
							goto l924
						l925:
							position, tokenIndex, depth = position924, tokenIndex924, depth924
							if buffer[position] != rune('T') {
								goto l918
							}
							position++
						}
					l924:
						{
							position926, tokenIndex926, depth926 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l927
							}
							position++
							goto l926
						l927:
							position, tokenIndex, depth = position926, tokenIndex926, depth926
							if buffer[position] != rune('H') {
								goto l918
							}
							position++
						}
					l926:
						if !_rules[rulesp]() {
							goto l918
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l918
						}
					l928:
						{
							position929, tokenIndex929, depth929 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l929
							}
							if buffer[position] != rune(',') {
								goto l929
							}
							position++
							if !_rules[rulespOpt]() {
								goto l929
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l929
							}
							goto l928
						l929:
							position, tokenIndex, depth = position929, tokenIndex929, depth929
						}
						goto l919
					l918:
						position, tokenIndex, depth = position918, tokenIndex918, depth918
					}
				l919:
					depth--
					add(rulePegText, position917)
				}
				if !_rules[ruleAction43]() {
					goto l915
				}
				depth--
				add(ruleSourceSinkSpecs, position916)
			}
			return true
		l915:
			position, tokenIndex, depth = position915, tokenIndex915, depth915
			return false
		},
		/* 58 UpdateSourceSinkSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)> Action44)> */
		func() bool {
			position930, tokenIndex930, depth930 := position, tokenIndex, depth
			{
				position931 := position
				depth++
				{
					position932 := position
					depth++
					if !_rules[rulesp]() {
						goto l930
					}
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
							goto l930
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
							goto l930
						}
						position++
					}
				l935:
					{
						position937, tokenIndex937, depth937 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l938
						}
						position++
						goto l937
					l938:
						position, tokenIndex, depth = position937, tokenIndex937, depth937
						if buffer[position] != rune('T') {
							goto l930
						}
						position++
					}
				l937:
					if !_rules[rulesp]() {
						goto l930
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l930
					}
				l939:
					{
						position940, tokenIndex940, depth940 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l940
						}
						if buffer[position] != rune(',') {
							goto l940
						}
						position++
						if !_rules[rulespOpt]() {
							goto l940
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l940
						}
						goto l939
					l940:
						position, tokenIndex, depth = position940, tokenIndex940, depth940
					}
					depth--
					add(rulePegText, position932)
				}
				if !_rules[ruleAction44]() {
					goto l930
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position931)
			}
			return true
		l930:
			position, tokenIndex, depth = position930, tokenIndex930, depth930
			return false
		},
		/* 59 SetOptSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action45)> */
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
						if !_rules[rulesp]() {
							goto l944
						}
						{
							position946, tokenIndex946, depth946 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l947
							}
							position++
							goto l946
						l947:
							position, tokenIndex, depth = position946, tokenIndex946, depth946
							if buffer[position] != rune('S') {
								goto l944
							}
							position++
						}
					l946:
						{
							position948, tokenIndex948, depth948 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l949
							}
							position++
							goto l948
						l949:
							position, tokenIndex, depth = position948, tokenIndex948, depth948
							if buffer[position] != rune('E') {
								goto l944
							}
							position++
						}
					l948:
						{
							position950, tokenIndex950, depth950 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l951
							}
							position++
							goto l950
						l951:
							position, tokenIndex, depth = position950, tokenIndex950, depth950
							if buffer[position] != rune('T') {
								goto l944
							}
							position++
						}
					l950:
						if !_rules[rulesp]() {
							goto l944
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l944
						}
					l952:
						{
							position953, tokenIndex953, depth953 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l953
							}
							if buffer[position] != rune(',') {
								goto l953
							}
							position++
							if !_rules[rulespOpt]() {
								goto l953
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l953
							}
							goto l952
						l953:
							position, tokenIndex, depth = position953, tokenIndex953, depth953
						}
						goto l945
					l944:
						position, tokenIndex, depth = position944, tokenIndex944, depth944
					}
				l945:
					depth--
					add(rulePegText, position943)
				}
				if !_rules[ruleAction45]() {
					goto l941
				}
				depth--
				add(ruleSetOptSpecs, position942)
			}
			return true
		l941:
			position, tokenIndex, depth = position941, tokenIndex941, depth941
			return false
		},
		/* 60 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action46)> */
		func() bool {
			position954, tokenIndex954, depth954 := position, tokenIndex, depth
			{
				position955 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l954
				}
				if buffer[position] != rune('=') {
					goto l954
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l954
				}
				if !_rules[ruleAction46]() {
					goto l954
				}
				depth--
				add(ruleSourceSinkParam, position955)
			}
			return true
		l954:
			position, tokenIndex, depth = position954, tokenIndex954, depth954
			return false
		},
		/* 61 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr)> */
		func() bool {
			position956, tokenIndex956, depth956 := position, tokenIndex, depth
			{
				position957 := position
				depth++
				{
					position958, tokenIndex958, depth958 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l959
					}
					goto l958
				l959:
					position, tokenIndex, depth = position958, tokenIndex958, depth958
					if !_rules[ruleParamArrayExpr]() {
						goto l956
					}
				}
			l958:
				depth--
				add(ruleSourceSinkParamVal, position957)
			}
			return true
		l956:
			position, tokenIndex, depth = position956, tokenIndex956, depth956
			return false
		},
		/* 62 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position960, tokenIndex960, depth960 := position, tokenIndex, depth
			{
				position961 := position
				depth++
				{
					position962, tokenIndex962, depth962 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l963
					}
					goto l962
				l963:
					position, tokenIndex, depth = position962, tokenIndex962, depth962
					if !_rules[ruleLiteral]() {
						goto l960
					}
				}
			l962:
				depth--
				add(ruleParamLiteral, position961)
			}
			return true
		l960:
			position, tokenIndex, depth = position960, tokenIndex960, depth960
			return false
		},
		/* 63 ParamArrayExpr <- <(<('[' spOpt (ParamLiteral (',' spOpt ParamLiteral)*)? spOpt ','? spOpt ']')> Action47)> */
		func() bool {
			position964, tokenIndex964, depth964 := position, tokenIndex, depth
			{
				position965 := position
				depth++
				{
					position966 := position
					depth++
					if buffer[position] != rune('[') {
						goto l964
					}
					position++
					if !_rules[rulespOpt]() {
						goto l964
					}
					{
						position967, tokenIndex967, depth967 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l967
						}
					l969:
						{
							position970, tokenIndex970, depth970 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l970
							}
							position++
							if !_rules[rulespOpt]() {
								goto l970
							}
							if !_rules[ruleParamLiteral]() {
								goto l970
							}
							goto l969
						l970:
							position, tokenIndex, depth = position970, tokenIndex970, depth970
						}
						goto l968
					l967:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
					}
				l968:
					if !_rules[rulespOpt]() {
						goto l964
					}
					{
						position971, tokenIndex971, depth971 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l971
						}
						position++
						goto l972
					l971:
						position, tokenIndex, depth = position971, tokenIndex971, depth971
					}
				l972:
					if !_rules[rulespOpt]() {
						goto l964
					}
					if buffer[position] != rune(']') {
						goto l964
					}
					position++
					depth--
					add(rulePegText, position966)
				}
				if !_rules[ruleAction47]() {
					goto l964
				}
				depth--
				add(ruleParamArrayExpr, position965)
			}
			return true
		l964:
			position, tokenIndex, depth = position964, tokenIndex964, depth964
			return false
		},
		/* 64 PausedOpt <- <(<(sp (Paused / Unpaused))?> Action48)> */
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
						if !_rules[rulesp]() {
							goto l976
						}
						{
							position978, tokenIndex978, depth978 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l979
							}
							goto l978
						l979:
							position, tokenIndex, depth = position978, tokenIndex978, depth978
							if !_rules[ruleUnpaused]() {
								goto l976
							}
						}
					l978:
						goto l977
					l976:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
					}
				l977:
					depth--
					add(rulePegText, position975)
				}
				if !_rules[ruleAction48]() {
					goto l973
				}
				depth--
				add(rulePausedOpt, position974)
			}
			return true
		l973:
			position, tokenIndex, depth = position973, tokenIndex973, depth973
			return false
		},
		/* 65 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position980, tokenIndex980, depth980 := position, tokenIndex, depth
			{
				position981 := position
				depth++
				{
					position982, tokenIndex982, depth982 := position, tokenIndex, depth
					if !_rules[ruleWildcard]() {
						goto l983
					}
					goto l982
				l983:
					position, tokenIndex, depth = position982, tokenIndex982, depth982
					if !_rules[ruleExpression]() {
						goto l980
					}
				}
			l982:
				depth--
				add(ruleExpressionOrWildcard, position981)
			}
			return true
		l980:
			position, tokenIndex, depth = position980, tokenIndex980, depth980
			return false
		},
		/* 66 Expression <- <orExpr> */
		func() bool {
			position984, tokenIndex984, depth984 := position, tokenIndex, depth
			{
				position985 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l984
				}
				depth--
				add(ruleExpression, position985)
			}
			return true
		l984:
			position, tokenIndex, depth = position984, tokenIndex984, depth984
			return false
		},
		/* 67 orExpr <- <(<(andExpr (sp Or sp andExpr)*)> Action49)> */
		func() bool {
			position986, tokenIndex986, depth986 := position, tokenIndex, depth
			{
				position987 := position
				depth++
				{
					position988 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l986
					}
				l989:
					{
						position990, tokenIndex990, depth990 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l990
						}
						if !_rules[ruleOr]() {
							goto l990
						}
						if !_rules[rulesp]() {
							goto l990
						}
						if !_rules[ruleandExpr]() {
							goto l990
						}
						goto l989
					l990:
						position, tokenIndex, depth = position990, tokenIndex990, depth990
					}
					depth--
					add(rulePegText, position988)
				}
				if !_rules[ruleAction49]() {
					goto l986
				}
				depth--
				add(ruleorExpr, position987)
			}
			return true
		l986:
			position, tokenIndex, depth = position986, tokenIndex986, depth986
			return false
		},
		/* 68 andExpr <- <(<(notExpr (sp And sp notExpr)*)> Action50)> */
		func() bool {
			position991, tokenIndex991, depth991 := position, tokenIndex, depth
			{
				position992 := position
				depth++
				{
					position993 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l991
					}
				l994:
					{
						position995, tokenIndex995, depth995 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l995
						}
						if !_rules[ruleAnd]() {
							goto l995
						}
						if !_rules[rulesp]() {
							goto l995
						}
						if !_rules[rulenotExpr]() {
							goto l995
						}
						goto l994
					l995:
						position, tokenIndex, depth = position995, tokenIndex995, depth995
					}
					depth--
					add(rulePegText, position993)
				}
				if !_rules[ruleAction50]() {
					goto l991
				}
				depth--
				add(ruleandExpr, position992)
			}
			return true
		l991:
			position, tokenIndex, depth = position991, tokenIndex991, depth991
			return false
		},
		/* 69 notExpr <- <(<((Not sp)? comparisonExpr)> Action51)> */
		func() bool {
			position996, tokenIndex996, depth996 := position, tokenIndex, depth
			{
				position997 := position
				depth++
				{
					position998 := position
					depth++
					{
						position999, tokenIndex999, depth999 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l999
						}
						if !_rules[rulesp]() {
							goto l999
						}
						goto l1000
					l999:
						position, tokenIndex, depth = position999, tokenIndex999, depth999
					}
				l1000:
					if !_rules[rulecomparisonExpr]() {
						goto l996
					}
					depth--
					add(rulePegText, position998)
				}
				if !_rules[ruleAction51]() {
					goto l996
				}
				depth--
				add(rulenotExpr, position997)
			}
			return true
		l996:
			position, tokenIndex, depth = position996, tokenIndex996, depth996
			return false
		},
		/* 70 comparisonExpr <- <(<(otherOpExpr (spOpt ComparisonOp spOpt otherOpExpr)?)> Action52)> */
		func() bool {
			position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
			{
				position1002 := position
				depth++
				{
					position1003 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l1001
					}
					{
						position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1004
						}
						if !_rules[ruleComparisonOp]() {
							goto l1004
						}
						if !_rules[rulespOpt]() {
							goto l1004
						}
						if !_rules[ruleotherOpExpr]() {
							goto l1004
						}
						goto l1005
					l1004:
						position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
					}
				l1005:
					depth--
					add(rulePegText, position1003)
				}
				if !_rules[ruleAction52]() {
					goto l1001
				}
				depth--
				add(rulecomparisonExpr, position1002)
			}
			return true
		l1001:
			position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
			return false
		},
		/* 71 otherOpExpr <- <(<(isExpr (spOpt OtherOp spOpt isExpr)*)> Action53)> */
		func() bool {
			position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
			{
				position1007 := position
				depth++
				{
					position1008 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l1006
					}
				l1009:
					{
						position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1010
						}
						if !_rules[ruleOtherOp]() {
							goto l1010
						}
						if !_rules[rulespOpt]() {
							goto l1010
						}
						if !_rules[ruleisExpr]() {
							goto l1010
						}
						goto l1009
					l1010:
						position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
					}
					depth--
					add(rulePegText, position1008)
				}
				if !_rules[ruleAction53]() {
					goto l1006
				}
				depth--
				add(ruleotherOpExpr, position1007)
			}
			return true
		l1006:
			position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
			return false
		},
		/* 72 isExpr <- <(<(termExpr (sp IsOp sp NullLiteral)?)> Action54)> */
		func() bool {
			position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
			{
				position1012 := position
				depth++
				{
					position1013 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l1011
					}
					{
						position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1014
						}
						if !_rules[ruleIsOp]() {
							goto l1014
						}
						if !_rules[rulesp]() {
							goto l1014
						}
						if !_rules[ruleNullLiteral]() {
							goto l1014
						}
						goto l1015
					l1014:
						position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
					}
				l1015:
					depth--
					add(rulePegText, position1013)
				}
				if !_rules[ruleAction54]() {
					goto l1011
				}
				depth--
				add(ruleisExpr, position1012)
			}
			return true
		l1011:
			position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
			return false
		},
		/* 73 termExpr <- <(<(productExpr (spOpt PlusMinusOp spOpt productExpr)*)> Action55)> */
		func() bool {
			position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
			{
				position1017 := position
				depth++
				{
					position1018 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l1016
					}
				l1019:
					{
						position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1020
						}
						if !_rules[rulePlusMinusOp]() {
							goto l1020
						}
						if !_rules[rulespOpt]() {
							goto l1020
						}
						if !_rules[ruleproductExpr]() {
							goto l1020
						}
						goto l1019
					l1020:
						position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
					}
					depth--
					add(rulePegText, position1018)
				}
				if !_rules[ruleAction55]() {
					goto l1016
				}
				depth--
				add(ruletermExpr, position1017)
			}
			return true
		l1016:
			position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
			return false
		},
		/* 74 productExpr <- <(<(minusExpr (spOpt MultDivOp spOpt minusExpr)*)> Action56)> */
		func() bool {
			position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
			{
				position1022 := position
				depth++
				{
					position1023 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l1021
					}
				l1024:
					{
						position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1025
						}
						if !_rules[ruleMultDivOp]() {
							goto l1025
						}
						if !_rules[rulespOpt]() {
							goto l1025
						}
						if !_rules[ruleminusExpr]() {
							goto l1025
						}
						goto l1024
					l1025:
						position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
					}
					depth--
					add(rulePegText, position1023)
				}
				if !_rules[ruleAction56]() {
					goto l1021
				}
				depth--
				add(ruleproductExpr, position1022)
			}
			return true
		l1021:
			position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
			return false
		},
		/* 75 minusExpr <- <(<((UnaryMinus spOpt)? castExpr)> Action57)> */
		func() bool {
			position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
			{
				position1027 := position
				depth++
				{
					position1028 := position
					depth++
					{
						position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l1029
						}
						if !_rules[rulespOpt]() {
							goto l1029
						}
						goto l1030
					l1029:
						position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
					}
				l1030:
					if !_rules[rulecastExpr]() {
						goto l1026
					}
					depth--
					add(rulePegText, position1028)
				}
				if !_rules[ruleAction57]() {
					goto l1026
				}
				depth--
				add(ruleminusExpr, position1027)
			}
			return true
		l1026:
			position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
			return false
		},
		/* 76 castExpr <- <(<(baseExpr (spOpt (':' ':') spOpt Type)?)> Action58)> */
		func() bool {
			position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
			{
				position1032 := position
				depth++
				{
					position1033 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l1031
					}
					{
						position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1034
						}
						if buffer[position] != rune(':') {
							goto l1034
						}
						position++
						if buffer[position] != rune(':') {
							goto l1034
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1034
						}
						if !_rules[ruleType]() {
							goto l1034
						}
						goto l1035
					l1034:
						position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
					}
				l1035:
					depth--
					add(rulePegText, position1033)
				}
				if !_rules[ruleAction58]() {
					goto l1031
				}
				depth--
				add(rulecastExpr, position1032)
			}
			return true
		l1031:
			position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
			return false
		},
		/* 77 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
			{
				position1037 := position
				depth++
				{
					position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l1039
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1039
					}
					if !_rules[ruleExpression]() {
						goto l1039
					}
					if !_rules[rulespOpt]() {
						goto l1039
					}
					if buffer[position] != rune(')') {
						goto l1039
					}
					position++
					goto l1038
				l1039:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleMapExpr]() {
						goto l1040
					}
					goto l1038
				l1040:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleBooleanLiteral]() {
						goto l1041
					}
					goto l1038
				l1041:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleNullLiteral]() {
						goto l1042
					}
					goto l1038
				l1042:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleRowMeta]() {
						goto l1043
					}
					goto l1038
				l1043:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleFuncTypeCast]() {
						goto l1044
					}
					goto l1038
				l1044:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleFuncApp]() {
						goto l1045
					}
					goto l1038
				l1045:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleRowValue]() {
						goto l1046
					}
					goto l1038
				l1046:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleArrayExpr]() {
						goto l1047
					}
					goto l1038
				l1047:
					position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
					if !_rules[ruleLiteral]() {
						goto l1036
					}
				}
			l1038:
				depth--
				add(rulebaseExpr, position1037)
			}
			return true
		l1036:
			position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
			return false
		},
		/* 78 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') spOpt '(' spOpt Expression sp (('a' / 'A') ('s' / 'S')) sp Type spOpt ')')> Action59)> */
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
						if buffer[position] != rune('c') {
							goto l1052
						}
						position++
						goto l1051
					l1052:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if buffer[position] != rune('C') {
							goto l1048
						}
						position++
					}
				l1051:
					{
						position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1054
						}
						position++
						goto l1053
					l1054:
						position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
						if buffer[position] != rune('A') {
							goto l1048
						}
						position++
					}
				l1053:
					{
						position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1056
						}
						position++
						goto l1055
					l1056:
						position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
						if buffer[position] != rune('S') {
							goto l1048
						}
						position++
					}
				l1055:
					{
						position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1058
						}
						position++
						goto l1057
					l1058:
						position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
						if buffer[position] != rune('T') {
							goto l1048
						}
						position++
					}
				l1057:
					if !_rules[rulespOpt]() {
						goto l1048
					}
					if buffer[position] != rune('(') {
						goto l1048
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1048
					}
					if !_rules[ruleExpression]() {
						goto l1048
					}
					if !_rules[rulesp]() {
						goto l1048
					}
					{
						position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1060
						}
						position++
						goto l1059
					l1060:
						position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
						if buffer[position] != rune('A') {
							goto l1048
						}
						position++
					}
				l1059:
					{
						position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1062
						}
						position++
						goto l1061
					l1062:
						position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
						if buffer[position] != rune('S') {
							goto l1048
						}
						position++
					}
				l1061:
					if !_rules[rulesp]() {
						goto l1048
					}
					if !_rules[ruleType]() {
						goto l1048
					}
					if !_rules[rulespOpt]() {
						goto l1048
					}
					if buffer[position] != rune(')') {
						goto l1048
					}
					position++
					depth--
					add(rulePegText, position1050)
				}
				if !_rules[ruleAction59]() {
					goto l1048
				}
				depth--
				add(ruleFuncTypeCast, position1049)
			}
			return true
		l1048:
			position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
			return false
		},
		/* 79 FuncApp <- <(FuncAppWithOrderBy / FuncAppWithoutOrderBy)> */
		func() bool {
			position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
			{
				position1064 := position
				depth++
				{
					position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
					if !_rules[ruleFuncAppWithOrderBy]() {
						goto l1066
					}
					goto l1065
				l1066:
					position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
					if !_rules[ruleFuncAppWithoutOrderBy]() {
						goto l1063
					}
				}
			l1065:
				depth--
				add(ruleFuncApp, position1064)
			}
			return true
		l1063:
			position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
			return false
		},
		/* 80 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action60)> */
		func() bool {
			position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
			{
				position1068 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1067
				}
				if !_rules[rulespOpt]() {
					goto l1067
				}
				if buffer[position] != rune('(') {
					goto l1067
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1067
				}
				if !_rules[ruleFuncParams]() {
					goto l1067
				}
				if !_rules[rulesp]() {
					goto l1067
				}
				if !_rules[ruleParamsOrder]() {
					goto l1067
				}
				if !_rules[rulespOpt]() {
					goto l1067
				}
				if buffer[position] != rune(')') {
					goto l1067
				}
				position++
				if !_rules[ruleAction60]() {
					goto l1067
				}
				depth--
				add(ruleFuncAppWithOrderBy, position1068)
			}
			return true
		l1067:
			position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
			return false
		},
		/* 81 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action61)> */
		func() bool {
			position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
			{
				position1070 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1069
				}
				if !_rules[rulespOpt]() {
					goto l1069
				}
				if buffer[position] != rune('(') {
					goto l1069
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1069
				}
				if !_rules[ruleFuncParams]() {
					goto l1069
				}
				{
					position1071 := position
					depth++
					if !_rules[rulespOpt]() {
						goto l1069
					}
					depth--
					add(rulePegText, position1071)
				}
				if buffer[position] != rune(')') {
					goto l1069
				}
				position++
				if !_rules[ruleAction61]() {
					goto l1069
				}
				depth--
				add(ruleFuncAppWithoutOrderBy, position1070)
			}
			return true
		l1069:
			position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
			return false
		},
		/* 82 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action62)> */
		func() bool {
			position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
			{
				position1073 := position
				depth++
				{
					position1074 := position
					depth++
					{
						position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1075
						}
					l1077:
						{
							position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1078
							}
							if buffer[position] != rune(',') {
								goto l1078
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1078
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1078
							}
							goto l1077
						l1078:
							position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
						}
						goto l1076
					l1075:
						position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
					}
				l1076:
					depth--
					add(rulePegText, position1074)
				}
				if !_rules[ruleAction62]() {
					goto l1072
				}
				depth--
				add(ruleFuncParams, position1073)
			}
			return true
		l1072:
			position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
			return false
		},
		/* 83 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action63)> */
		func() bool {
			position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
			{
				position1080 := position
				depth++
				{
					position1081 := position
					depth++
					{
						position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1083
						}
						position++
						goto l1082
					l1083:
						position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
						if buffer[position] != rune('O') {
							goto l1079
						}
						position++
					}
				l1082:
					{
						position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1085
						}
						position++
						goto l1084
					l1085:
						position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
						if buffer[position] != rune('R') {
							goto l1079
						}
						position++
					}
				l1084:
					{
						position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1087
						}
						position++
						goto l1086
					l1087:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if buffer[position] != rune('D') {
							goto l1079
						}
						position++
					}
				l1086:
					{
						position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1089
						}
						position++
						goto l1088
					l1089:
						position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
						if buffer[position] != rune('E') {
							goto l1079
						}
						position++
					}
				l1088:
					{
						position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1091
						}
						position++
						goto l1090
					l1091:
						position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
						if buffer[position] != rune('R') {
							goto l1079
						}
						position++
					}
				l1090:
					if !_rules[rulesp]() {
						goto l1079
					}
					{
						position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1093
						}
						position++
						goto l1092
					l1093:
						position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
						if buffer[position] != rune('B') {
							goto l1079
						}
						position++
					}
				l1092:
					{
						position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1095
						}
						position++
						goto l1094
					l1095:
						position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
						if buffer[position] != rune('Y') {
							goto l1079
						}
						position++
					}
				l1094:
					if !_rules[rulesp]() {
						goto l1079
					}
					if !_rules[ruleSortedExpression]() {
						goto l1079
					}
				l1096:
					{
						position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1097
						}
						if buffer[position] != rune(',') {
							goto l1097
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1097
						}
						if !_rules[ruleSortedExpression]() {
							goto l1097
						}
						goto l1096
					l1097:
						position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
					}
					depth--
					add(rulePegText, position1081)
				}
				if !_rules[ruleAction63]() {
					goto l1079
				}
				depth--
				add(ruleParamsOrder, position1080)
			}
			return true
		l1079:
			position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
			return false
		},
		/* 84 SortedExpression <- <(Expression OrderDirectionOpt Action64)> */
		func() bool {
			position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
			{
				position1099 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l1098
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1098
				}
				if !_rules[ruleAction64]() {
					goto l1098
				}
				depth--
				add(ruleSortedExpression, position1099)
			}
			return true
		l1098:
			position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
			return false
		},
		/* 85 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action65)> */
		func() bool {
			position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
			{
				position1101 := position
				depth++
				{
					position1102 := position
					depth++
					{
						position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1103
						}
						{
							position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
							if !_rules[ruleAscending]() {
								goto l1106
							}
							goto l1105
						l1106:
							position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
							if !_rules[ruleDescending]() {
								goto l1103
							}
						}
					l1105:
						goto l1104
					l1103:
						position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
					}
				l1104:
					depth--
					add(rulePegText, position1102)
				}
				if !_rules[ruleAction65]() {
					goto l1100
				}
				depth--
				add(ruleOrderDirectionOpt, position1101)
			}
			return true
		l1100:
			position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
			return false
		},
		/* 86 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action66)> */
		func() bool {
			position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
			{
				position1108 := position
				depth++
				{
					position1109 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1107
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1107
					}
					{
						position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1110
						}
					l1112:
						{
							position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1113
							}
							if buffer[position] != rune(',') {
								goto l1113
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1113
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1113
							}
							goto l1112
						l1113:
							position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
						}
						goto l1111
					l1110:
						position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
					}
				l1111:
					if !_rules[rulespOpt]() {
						goto l1107
					}
					{
						position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1114
						}
						position++
						goto l1115
					l1114:
						position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					}
				l1115:
					if !_rules[rulespOpt]() {
						goto l1107
					}
					if buffer[position] != rune(']') {
						goto l1107
					}
					position++
					depth--
					add(rulePegText, position1109)
				}
				if !_rules[ruleAction66]() {
					goto l1107
				}
				depth--
				add(ruleArrayExpr, position1108)
			}
			return true
		l1107:
			position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
			return false
		},
		/* 87 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action67)> */
		func() bool {
			position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
			{
				position1117 := position
				depth++
				{
					position1118 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1116
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1116
					}
					{
						position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l1119
						}
					l1121:
						{
							position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1122
							}
							if buffer[position] != rune(',') {
								goto l1122
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1122
							}
							if !_rules[ruleKeyValuePair]() {
								goto l1122
							}
							goto l1121
						l1122:
							position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
						}
						goto l1120
					l1119:
						position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
					}
				l1120:
					if !_rules[rulespOpt]() {
						goto l1116
					}
					if buffer[position] != rune('}') {
						goto l1116
					}
					position++
					depth--
					add(rulePegText, position1118)
				}
				if !_rules[ruleAction67]() {
					goto l1116
				}
				depth--
				add(ruleMapExpr, position1117)
			}
			return true
		l1116:
			position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
			return false
		},
		/* 88 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action68)> */
		func() bool {
			position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
			{
				position1124 := position
				depth++
				{
					position1125 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1123
					}
					if !_rules[rulespOpt]() {
						goto l1123
					}
					if buffer[position] != rune(':') {
						goto l1123
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1123
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1123
					}
					depth--
					add(rulePegText, position1125)
				}
				if !_rules[ruleAction68]() {
					goto l1123
				}
				depth--
				add(ruleKeyValuePair, position1124)
			}
			return true
		l1123:
			position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
			return false
		},
		/* 89 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
			{
				position1127 := position
				depth++
				{
					position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l1129
					}
					goto l1128
				l1129:
					position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
					if !_rules[ruleNumericLiteral]() {
						goto l1130
					}
					goto l1128
				l1130:
					position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
					if !_rules[ruleStringLiteral]() {
						goto l1126
					}
				}
			l1128:
				depth--
				add(ruleLiteral, position1127)
			}
			return true
		l1126:
			position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
			return false
		},
		/* 90 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
			{
				position1132 := position
				depth++
				{
					position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l1134
					}
					goto l1133
				l1134:
					position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
					if !_rules[ruleNotEqual]() {
						goto l1135
					}
					goto l1133
				l1135:
					position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
					if !_rules[ruleLessOrEqual]() {
						goto l1136
					}
					goto l1133
				l1136:
					position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
					if !_rules[ruleLess]() {
						goto l1137
					}
					goto l1133
				l1137:
					position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
					if !_rules[ruleGreaterOrEqual]() {
						goto l1138
					}
					goto l1133
				l1138:
					position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
					if !_rules[ruleGreater]() {
						goto l1139
					}
					goto l1133
				l1139:
					position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
					if !_rules[ruleNotEqual]() {
						goto l1131
					}
				}
			l1133:
				depth--
				add(ruleComparisonOp, position1132)
			}
			return true
		l1131:
			position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
			return false
		},
		/* 91 OtherOp <- <Concat> */
		func() bool {
			position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
			{
				position1141 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l1140
				}
				depth--
				add(ruleOtherOp, position1141)
			}
			return true
		l1140:
			position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
			return false
		},
		/* 92 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
			{
				position1143 := position
				depth++
				{
					position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l1145
					}
					goto l1144
				l1145:
					position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
					if !_rules[ruleIs]() {
						goto l1142
					}
				}
			l1144:
				depth--
				add(ruleIsOp, position1143)
			}
			return true
		l1142:
			position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
			return false
		},
		/* 93 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
			{
				position1147 := position
				depth++
				{
					position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l1149
					}
					goto l1148
				l1149:
					position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
					if !_rules[ruleMinus]() {
						goto l1146
					}
				}
			l1148:
				depth--
				add(rulePlusMinusOp, position1147)
			}
			return true
		l1146:
			position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
			return false
		},
		/* 94 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
			{
				position1151 := position
				depth++
				{
					position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l1153
					}
					goto l1152
				l1153:
					position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
					if !_rules[ruleDivide]() {
						goto l1154
					}
					goto l1152
				l1154:
					position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
					if !_rules[ruleModulo]() {
						goto l1150
					}
				}
			l1152:
				depth--
				add(ruleMultDivOp, position1151)
			}
			return true
		l1150:
			position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
			return false
		},
		/* 95 Stream <- <(<ident> Action69)> */
		func() bool {
			position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
			{
				position1156 := position
				depth++
				{
					position1157 := position
					depth++
					if !_rules[ruleident]() {
						goto l1155
					}
					depth--
					add(rulePegText, position1157)
				}
				if !_rules[ruleAction69]() {
					goto l1155
				}
				depth--
				add(ruleStream, position1156)
			}
			return true
		l1155:
			position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
			return false
		},
		/* 96 RowMeta <- <RowTimestamp> */
		func() bool {
			position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
			{
				position1159 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l1158
				}
				depth--
				add(ruleRowMeta, position1159)
			}
			return true
		l1158:
			position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
			return false
		},
		/* 97 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action70)> */
		func() bool {
			position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
			{
				position1161 := position
				depth++
				{
					position1162 := position
					depth++
					{
						position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1163
						}
						if buffer[position] != rune(':') {
							goto l1163
						}
						position++
						goto l1164
					l1163:
						position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
					}
				l1164:
					if buffer[position] != rune('t') {
						goto l1160
					}
					position++
					if buffer[position] != rune('s') {
						goto l1160
					}
					position++
					if buffer[position] != rune('(') {
						goto l1160
					}
					position++
					if buffer[position] != rune(')') {
						goto l1160
					}
					position++
					depth--
					add(rulePegText, position1162)
				}
				if !_rules[ruleAction70]() {
					goto l1160
				}
				depth--
				add(ruleRowTimestamp, position1161)
			}
			return true
		l1160:
			position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
			return false
		},
		/* 98 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action71)> */
		func() bool {
			position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
			{
				position1166 := position
				depth++
				{
					position1167 := position
					depth++
					{
						position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1168
						}
						if buffer[position] != rune(':') {
							goto l1168
						}
						position++
						{
							position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1170
							}
							position++
							goto l1168
						l1170:
							position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						}
						goto l1169
					l1168:
						position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
					}
				l1169:
					if !_rules[rulejsonGetPath]() {
						goto l1165
					}
					depth--
					add(rulePegText, position1167)
				}
				if !_rules[ruleAction71]() {
					goto l1165
				}
				depth--
				add(ruleRowValue, position1166)
			}
			return true
		l1165:
			position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
			return false
		},
		/* 99 NumericLiteral <- <(<('-'? [0-9]+)> Action72)> */
		func() bool {
			position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
			{
				position1172 := position
				depth++
				{
					position1173 := position
					depth++
					{
						position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1174
						}
						position++
						goto l1175
					l1174:
						position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
					}
				l1175:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1171
					}
					position++
				l1176:
					{
						position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1177
						}
						position++
						goto l1176
					l1177:
						position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
					}
					depth--
					add(rulePegText, position1173)
				}
				if !_rules[ruleAction72]() {
					goto l1171
				}
				depth--
				add(ruleNumericLiteral, position1172)
			}
			return true
		l1171:
			position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
			return false
		},
		/* 100 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action73)> */
		func() bool {
			position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
			{
				position1179 := position
				depth++
				{
					position1180 := position
					depth++
					{
						position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1181
						}
						position++
						goto l1182
					l1181:
						position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
					}
				l1182:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1178
					}
					position++
				l1183:
					{
						position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1184
						}
						position++
						goto l1183
					l1184:
						position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
					}
					if buffer[position] != rune('.') {
						goto l1178
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1178
					}
					position++
				l1185:
					{
						position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1186
						}
						position++
						goto l1185
					l1186:
						position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
					}
					depth--
					add(rulePegText, position1180)
				}
				if !_rules[ruleAction73]() {
					goto l1178
				}
				depth--
				add(ruleFloatLiteral, position1179)
			}
			return true
		l1178:
			position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
			return false
		},
		/* 101 Function <- <(<ident> Action74)> */
		func() bool {
			position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
			{
				position1188 := position
				depth++
				{
					position1189 := position
					depth++
					if !_rules[ruleident]() {
						goto l1187
					}
					depth--
					add(rulePegText, position1189)
				}
				if !_rules[ruleAction74]() {
					goto l1187
				}
				depth--
				add(ruleFunction, position1188)
			}
			return true
		l1187:
			position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
			return false
		},
		/* 102 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action75)> */
		func() bool {
			position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
			{
				position1191 := position
				depth++
				{
					position1192 := position
					depth++
					{
						position1193, tokenIndex1193, depth1193 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1194
						}
						position++
						goto l1193
					l1194:
						position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
						if buffer[position] != rune('N') {
							goto l1190
						}
						position++
					}
				l1193:
					{
						position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1196
						}
						position++
						goto l1195
					l1196:
						position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
						if buffer[position] != rune('U') {
							goto l1190
						}
						position++
					}
				l1195:
					{
						position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1198
						}
						position++
						goto l1197
					l1198:
						position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
						if buffer[position] != rune('L') {
							goto l1190
						}
						position++
					}
				l1197:
					{
						position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1200
						}
						position++
						goto l1199
					l1200:
						position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
						if buffer[position] != rune('L') {
							goto l1190
						}
						position++
					}
				l1199:
					depth--
					add(rulePegText, position1192)
				}
				if !_rules[ruleAction75]() {
					goto l1190
				}
				depth--
				add(ruleNullLiteral, position1191)
			}
			return true
		l1190:
			position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
			return false
		},
		/* 103 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
			{
				position1202 := position
				depth++
				{
					position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l1204
					}
					goto l1203
				l1204:
					position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
					if !_rules[ruleFALSE]() {
						goto l1201
					}
				}
			l1203:
				depth--
				add(ruleBooleanLiteral, position1202)
			}
			return true
		l1201:
			position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
			return false
		},
		/* 104 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action76)> */
		func() bool {
			position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
			{
				position1206 := position
				depth++
				{
					position1207 := position
					depth++
					{
						position1208, tokenIndex1208, depth1208 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1209
						}
						position++
						goto l1208
					l1209:
						position, tokenIndex, depth = position1208, tokenIndex1208, depth1208
						if buffer[position] != rune('T') {
							goto l1205
						}
						position++
					}
				l1208:
					{
						position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1211
						}
						position++
						goto l1210
					l1211:
						position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
						if buffer[position] != rune('R') {
							goto l1205
						}
						position++
					}
				l1210:
					{
						position1212, tokenIndex1212, depth1212 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1213
						}
						position++
						goto l1212
					l1213:
						position, tokenIndex, depth = position1212, tokenIndex1212, depth1212
						if buffer[position] != rune('U') {
							goto l1205
						}
						position++
					}
				l1212:
					{
						position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1215
						}
						position++
						goto l1214
					l1215:
						position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
						if buffer[position] != rune('E') {
							goto l1205
						}
						position++
					}
				l1214:
					depth--
					add(rulePegText, position1207)
				}
				if !_rules[ruleAction76]() {
					goto l1205
				}
				depth--
				add(ruleTRUE, position1206)
			}
			return true
		l1205:
			position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
			return false
		},
		/* 105 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action77)> */
		func() bool {
			position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
			{
				position1217 := position
				depth++
				{
					position1218 := position
					depth++
					{
						position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1220
						}
						position++
						goto l1219
					l1220:
						position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
						if buffer[position] != rune('F') {
							goto l1216
						}
						position++
					}
				l1219:
					{
						position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1222
						}
						position++
						goto l1221
					l1222:
						position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
						if buffer[position] != rune('A') {
							goto l1216
						}
						position++
					}
				l1221:
					{
						position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1224
						}
						position++
						goto l1223
					l1224:
						position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
						if buffer[position] != rune('L') {
							goto l1216
						}
						position++
					}
				l1223:
					{
						position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1226
						}
						position++
						goto l1225
					l1226:
						position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
						if buffer[position] != rune('S') {
							goto l1216
						}
						position++
					}
				l1225:
					{
						position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1228
						}
						position++
						goto l1227
					l1228:
						position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
						if buffer[position] != rune('E') {
							goto l1216
						}
						position++
					}
				l1227:
					depth--
					add(rulePegText, position1218)
				}
				if !_rules[ruleAction77]() {
					goto l1216
				}
				depth--
				add(ruleFALSE, position1217)
			}
			return true
		l1216:
			position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
			return false
		},
		/* 106 Wildcard <- <(<((ident ':' !':')? '*')> Action78)> */
		func() bool {
			position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
			{
				position1230 := position
				depth++
				{
					position1231 := position
					depth++
					{
						position1232, tokenIndex1232, depth1232 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1232
						}
						if buffer[position] != rune(':') {
							goto l1232
						}
						position++
						{
							position1234, tokenIndex1234, depth1234 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1234
							}
							position++
							goto l1232
						l1234:
							position, tokenIndex, depth = position1234, tokenIndex1234, depth1234
						}
						goto l1233
					l1232:
						position, tokenIndex, depth = position1232, tokenIndex1232, depth1232
					}
				l1233:
					if buffer[position] != rune('*') {
						goto l1229
					}
					position++
					depth--
					add(rulePegText, position1231)
				}
				if !_rules[ruleAction78]() {
					goto l1229
				}
				depth--
				add(ruleWildcard, position1230)
			}
			return true
		l1229:
			position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
			return false
		},
		/* 107 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action79)> */
		func() bool {
			position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
			{
				position1236 := position
				depth++
				{
					position1237 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l1235
					}
					position++
				l1238:
					{
						position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
						{
							position1240, tokenIndex1240, depth1240 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1241
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1241
							}
							position++
							goto l1240
						l1241:
							position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
							{
								position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1242
								}
								position++
								goto l1239
							l1242:
								position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
							}
							if !matchDot() {
								goto l1239
							}
						}
					l1240:
						goto l1238
					l1239:
						position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
					}
					if buffer[position] != rune('\'') {
						goto l1235
					}
					position++
					depth--
					add(rulePegText, position1237)
				}
				if !_rules[ruleAction79]() {
					goto l1235
				}
				depth--
				add(ruleStringLiteral, position1236)
			}
			return true
		l1235:
			position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
			return false
		},
		/* 108 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action80)> */
		func() bool {
			position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
			{
				position1244 := position
				depth++
				{
					position1245 := position
					depth++
					{
						position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1247
						}
						position++
						goto l1246
					l1247:
						position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
						if buffer[position] != rune('I') {
							goto l1243
						}
						position++
					}
				l1246:
					{
						position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1249
						}
						position++
						goto l1248
					l1249:
						position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
						if buffer[position] != rune('S') {
							goto l1243
						}
						position++
					}
				l1248:
					{
						position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1251
						}
						position++
						goto l1250
					l1251:
						position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
						if buffer[position] != rune('T') {
							goto l1243
						}
						position++
					}
				l1250:
					{
						position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1253
						}
						position++
						goto l1252
					l1253:
						position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
						if buffer[position] != rune('R') {
							goto l1243
						}
						position++
					}
				l1252:
					{
						position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1255
						}
						position++
						goto l1254
					l1255:
						position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
						if buffer[position] != rune('E') {
							goto l1243
						}
						position++
					}
				l1254:
					{
						position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1257
						}
						position++
						goto l1256
					l1257:
						position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
						if buffer[position] != rune('A') {
							goto l1243
						}
						position++
					}
				l1256:
					{
						position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1259
						}
						position++
						goto l1258
					l1259:
						position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
						if buffer[position] != rune('M') {
							goto l1243
						}
						position++
					}
				l1258:
					depth--
					add(rulePegText, position1245)
				}
				if !_rules[ruleAction80]() {
					goto l1243
				}
				depth--
				add(ruleISTREAM, position1244)
			}
			return true
		l1243:
			position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
			return false
		},
		/* 109 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action81)> */
		func() bool {
			position1260, tokenIndex1260, depth1260 := position, tokenIndex, depth
			{
				position1261 := position
				depth++
				{
					position1262 := position
					depth++
					{
						position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1264
						}
						position++
						goto l1263
					l1264:
						position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
						if buffer[position] != rune('D') {
							goto l1260
						}
						position++
					}
				l1263:
					{
						position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1266
						}
						position++
						goto l1265
					l1266:
						position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
						if buffer[position] != rune('S') {
							goto l1260
						}
						position++
					}
				l1265:
					{
						position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1268
						}
						position++
						goto l1267
					l1268:
						position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
						if buffer[position] != rune('T') {
							goto l1260
						}
						position++
					}
				l1267:
					{
						position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1270
						}
						position++
						goto l1269
					l1270:
						position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
						if buffer[position] != rune('R') {
							goto l1260
						}
						position++
					}
				l1269:
					{
						position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1272
						}
						position++
						goto l1271
					l1272:
						position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
						if buffer[position] != rune('E') {
							goto l1260
						}
						position++
					}
				l1271:
					{
						position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1274
						}
						position++
						goto l1273
					l1274:
						position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
						if buffer[position] != rune('A') {
							goto l1260
						}
						position++
					}
				l1273:
					{
						position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1276
						}
						position++
						goto l1275
					l1276:
						position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
						if buffer[position] != rune('M') {
							goto l1260
						}
						position++
					}
				l1275:
					depth--
					add(rulePegText, position1262)
				}
				if !_rules[ruleAction81]() {
					goto l1260
				}
				depth--
				add(ruleDSTREAM, position1261)
			}
			return true
		l1260:
			position, tokenIndex, depth = position1260, tokenIndex1260, depth1260
			return false
		},
		/* 110 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action82)> */
		func() bool {
			position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
			{
				position1278 := position
				depth++
				{
					position1279 := position
					depth++
					{
						position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1281
						}
						position++
						goto l1280
					l1281:
						position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
						if buffer[position] != rune('R') {
							goto l1277
						}
						position++
					}
				l1280:
					{
						position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1283
						}
						position++
						goto l1282
					l1283:
						position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
						if buffer[position] != rune('S') {
							goto l1277
						}
						position++
					}
				l1282:
					{
						position1284, tokenIndex1284, depth1284 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1285
						}
						position++
						goto l1284
					l1285:
						position, tokenIndex, depth = position1284, tokenIndex1284, depth1284
						if buffer[position] != rune('T') {
							goto l1277
						}
						position++
					}
				l1284:
					{
						position1286, tokenIndex1286, depth1286 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1287
						}
						position++
						goto l1286
					l1287:
						position, tokenIndex, depth = position1286, tokenIndex1286, depth1286
						if buffer[position] != rune('R') {
							goto l1277
						}
						position++
					}
				l1286:
					{
						position1288, tokenIndex1288, depth1288 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1289
						}
						position++
						goto l1288
					l1289:
						position, tokenIndex, depth = position1288, tokenIndex1288, depth1288
						if buffer[position] != rune('E') {
							goto l1277
						}
						position++
					}
				l1288:
					{
						position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1291
						}
						position++
						goto l1290
					l1291:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if buffer[position] != rune('A') {
							goto l1277
						}
						position++
					}
				l1290:
					{
						position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1293
						}
						position++
						goto l1292
					l1293:
						position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
						if buffer[position] != rune('M') {
							goto l1277
						}
						position++
					}
				l1292:
					depth--
					add(rulePegText, position1279)
				}
				if !_rules[ruleAction82]() {
					goto l1277
				}
				depth--
				add(ruleRSTREAM, position1278)
			}
			return true
		l1277:
			position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
			return false
		},
		/* 111 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action83)> */
		func() bool {
			position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
			{
				position1295 := position
				depth++
				{
					position1296 := position
					depth++
					{
						position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1298
						}
						position++
						goto l1297
					l1298:
						position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
						if buffer[position] != rune('T') {
							goto l1294
						}
						position++
					}
				l1297:
					{
						position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1300
						}
						position++
						goto l1299
					l1300:
						position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
						if buffer[position] != rune('U') {
							goto l1294
						}
						position++
					}
				l1299:
					{
						position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1302
						}
						position++
						goto l1301
					l1302:
						position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
						if buffer[position] != rune('P') {
							goto l1294
						}
						position++
					}
				l1301:
					{
						position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1304
						}
						position++
						goto l1303
					l1304:
						position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
						if buffer[position] != rune('L') {
							goto l1294
						}
						position++
					}
				l1303:
					{
						position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1306
						}
						position++
						goto l1305
					l1306:
						position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
						if buffer[position] != rune('E') {
							goto l1294
						}
						position++
					}
				l1305:
					{
						position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1308
						}
						position++
						goto l1307
					l1308:
						position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
						if buffer[position] != rune('S') {
							goto l1294
						}
						position++
					}
				l1307:
					depth--
					add(rulePegText, position1296)
				}
				if !_rules[ruleAction83]() {
					goto l1294
				}
				depth--
				add(ruleTUPLES, position1295)
			}
			return true
		l1294:
			position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
			return false
		},
		/* 112 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action84)> */
		func() bool {
			position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
			{
				position1310 := position
				depth++
				{
					position1311 := position
					depth++
					{
						position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1313
						}
						position++
						goto l1312
					l1313:
						position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
						if buffer[position] != rune('S') {
							goto l1309
						}
						position++
					}
				l1312:
					{
						position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1315
						}
						position++
						goto l1314
					l1315:
						position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
						if buffer[position] != rune('E') {
							goto l1309
						}
						position++
					}
				l1314:
					{
						position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1317
						}
						position++
						goto l1316
					l1317:
						position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
						if buffer[position] != rune('C') {
							goto l1309
						}
						position++
					}
				l1316:
					{
						position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1319
						}
						position++
						goto l1318
					l1319:
						position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
						if buffer[position] != rune('O') {
							goto l1309
						}
						position++
					}
				l1318:
					{
						position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1321
						}
						position++
						goto l1320
					l1321:
						position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
						if buffer[position] != rune('N') {
							goto l1309
						}
						position++
					}
				l1320:
					{
						position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1323
						}
						position++
						goto l1322
					l1323:
						position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
						if buffer[position] != rune('D') {
							goto l1309
						}
						position++
					}
				l1322:
					{
						position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1325
						}
						position++
						goto l1324
					l1325:
						position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
						if buffer[position] != rune('S') {
							goto l1309
						}
						position++
					}
				l1324:
					depth--
					add(rulePegText, position1311)
				}
				if !_rules[ruleAction84]() {
					goto l1309
				}
				depth--
				add(ruleSECONDS, position1310)
			}
			return true
		l1309:
			position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
			return false
		},
		/* 113 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action85)> */
		func() bool {
			position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
			{
				position1327 := position
				depth++
				{
					position1328 := position
					depth++
					{
						position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1330
						}
						position++
						goto l1329
					l1330:
						position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
						if buffer[position] != rune('M') {
							goto l1326
						}
						position++
					}
				l1329:
					{
						position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1332
						}
						position++
						goto l1331
					l1332:
						position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
						if buffer[position] != rune('I') {
							goto l1326
						}
						position++
					}
				l1331:
					{
						position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1334
						}
						position++
						goto l1333
					l1334:
						position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
						if buffer[position] != rune('L') {
							goto l1326
						}
						position++
					}
				l1333:
					{
						position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1336
						}
						position++
						goto l1335
					l1336:
						position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
						if buffer[position] != rune('L') {
							goto l1326
						}
						position++
					}
				l1335:
					{
						position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1338
						}
						position++
						goto l1337
					l1338:
						position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
						if buffer[position] != rune('I') {
							goto l1326
						}
						position++
					}
				l1337:
					{
						position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1340
						}
						position++
						goto l1339
					l1340:
						position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
						if buffer[position] != rune('S') {
							goto l1326
						}
						position++
					}
				l1339:
					{
						position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1342
						}
						position++
						goto l1341
					l1342:
						position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
						if buffer[position] != rune('E') {
							goto l1326
						}
						position++
					}
				l1341:
					{
						position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1344
						}
						position++
						goto l1343
					l1344:
						position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
						if buffer[position] != rune('C') {
							goto l1326
						}
						position++
					}
				l1343:
					{
						position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1346
						}
						position++
						goto l1345
					l1346:
						position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
						if buffer[position] != rune('O') {
							goto l1326
						}
						position++
					}
				l1345:
					{
						position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1348
						}
						position++
						goto l1347
					l1348:
						position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
						if buffer[position] != rune('N') {
							goto l1326
						}
						position++
					}
				l1347:
					{
						position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1350
						}
						position++
						goto l1349
					l1350:
						position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
						if buffer[position] != rune('D') {
							goto l1326
						}
						position++
					}
				l1349:
					{
						position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1352
						}
						position++
						goto l1351
					l1352:
						position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
						if buffer[position] != rune('S') {
							goto l1326
						}
						position++
					}
				l1351:
					depth--
					add(rulePegText, position1328)
				}
				if !_rules[ruleAction85]() {
					goto l1326
				}
				depth--
				add(ruleMILLISECONDS, position1327)
			}
			return true
		l1326:
			position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
			return false
		},
		/* 114 StreamIdentifier <- <(<ident> Action86)> */
		func() bool {
			position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
			{
				position1354 := position
				depth++
				{
					position1355 := position
					depth++
					if !_rules[ruleident]() {
						goto l1353
					}
					depth--
					add(rulePegText, position1355)
				}
				if !_rules[ruleAction86]() {
					goto l1353
				}
				depth--
				add(ruleStreamIdentifier, position1354)
			}
			return true
		l1353:
			position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
			return false
		},
		/* 115 SourceSinkType <- <(<ident> Action87)> */
		func() bool {
			position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
			{
				position1357 := position
				depth++
				{
					position1358 := position
					depth++
					if !_rules[ruleident]() {
						goto l1356
					}
					depth--
					add(rulePegText, position1358)
				}
				if !_rules[ruleAction87]() {
					goto l1356
				}
				depth--
				add(ruleSourceSinkType, position1357)
			}
			return true
		l1356:
			position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
			return false
		},
		/* 116 SourceSinkParamKey <- <(<ident> Action88)> */
		func() bool {
			position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
			{
				position1360 := position
				depth++
				{
					position1361 := position
					depth++
					if !_rules[ruleident]() {
						goto l1359
					}
					depth--
					add(rulePegText, position1361)
				}
				if !_rules[ruleAction88]() {
					goto l1359
				}
				depth--
				add(ruleSourceSinkParamKey, position1360)
			}
			return true
		l1359:
			position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
			return false
		},
		/* 117 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action89)> */
		func() bool {
			position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
			{
				position1363 := position
				depth++
				{
					position1364 := position
					depth++
					{
						position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1366
						}
						position++
						goto l1365
					l1366:
						position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
						if buffer[position] != rune('P') {
							goto l1362
						}
						position++
					}
				l1365:
					{
						position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1368
						}
						position++
						goto l1367
					l1368:
						position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
						if buffer[position] != rune('A') {
							goto l1362
						}
						position++
					}
				l1367:
					{
						position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1370
						}
						position++
						goto l1369
					l1370:
						position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
						if buffer[position] != rune('U') {
							goto l1362
						}
						position++
					}
				l1369:
					{
						position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1372
						}
						position++
						goto l1371
					l1372:
						position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
						if buffer[position] != rune('S') {
							goto l1362
						}
						position++
					}
				l1371:
					{
						position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1374
						}
						position++
						goto l1373
					l1374:
						position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
						if buffer[position] != rune('E') {
							goto l1362
						}
						position++
					}
				l1373:
					{
						position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1376
						}
						position++
						goto l1375
					l1376:
						position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
						if buffer[position] != rune('D') {
							goto l1362
						}
						position++
					}
				l1375:
					depth--
					add(rulePegText, position1364)
				}
				if !_rules[ruleAction89]() {
					goto l1362
				}
				depth--
				add(rulePaused, position1363)
			}
			return true
		l1362:
			position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
			return false
		},
		/* 118 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action90)> */
		func() bool {
			position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
			{
				position1378 := position
				depth++
				{
					position1379 := position
					depth++
					{
						position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1381
						}
						position++
						goto l1380
					l1381:
						position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
						if buffer[position] != rune('U') {
							goto l1377
						}
						position++
					}
				l1380:
					{
						position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1383
						}
						position++
						goto l1382
					l1383:
						position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
						if buffer[position] != rune('N') {
							goto l1377
						}
						position++
					}
				l1382:
					{
						position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1385
						}
						position++
						goto l1384
					l1385:
						position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
						if buffer[position] != rune('P') {
							goto l1377
						}
						position++
					}
				l1384:
					{
						position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1387
						}
						position++
						goto l1386
					l1387:
						position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
						if buffer[position] != rune('A') {
							goto l1377
						}
						position++
					}
				l1386:
					{
						position1388, tokenIndex1388, depth1388 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1389
						}
						position++
						goto l1388
					l1389:
						position, tokenIndex, depth = position1388, tokenIndex1388, depth1388
						if buffer[position] != rune('U') {
							goto l1377
						}
						position++
					}
				l1388:
					{
						position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1391
						}
						position++
						goto l1390
					l1391:
						position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
						if buffer[position] != rune('S') {
							goto l1377
						}
						position++
					}
				l1390:
					{
						position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1393
						}
						position++
						goto l1392
					l1393:
						position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
						if buffer[position] != rune('E') {
							goto l1377
						}
						position++
					}
				l1392:
					{
						position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1395
						}
						position++
						goto l1394
					l1395:
						position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
						if buffer[position] != rune('D') {
							goto l1377
						}
						position++
					}
				l1394:
					depth--
					add(rulePegText, position1379)
				}
				if !_rules[ruleAction90]() {
					goto l1377
				}
				depth--
				add(ruleUnpaused, position1378)
			}
			return true
		l1377:
			position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
			return false
		},
		/* 119 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action91)> */
		func() bool {
			position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
			{
				position1397 := position
				depth++
				{
					position1398 := position
					depth++
					{
						position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1400
						}
						position++
						goto l1399
					l1400:
						position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
						if buffer[position] != rune('A') {
							goto l1396
						}
						position++
					}
				l1399:
					{
						position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1402
						}
						position++
						goto l1401
					l1402:
						position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
						if buffer[position] != rune('S') {
							goto l1396
						}
						position++
					}
				l1401:
					{
						position1403, tokenIndex1403, depth1403 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1404
						}
						position++
						goto l1403
					l1404:
						position, tokenIndex, depth = position1403, tokenIndex1403, depth1403
						if buffer[position] != rune('C') {
							goto l1396
						}
						position++
					}
				l1403:
					depth--
					add(rulePegText, position1398)
				}
				if !_rules[ruleAction91]() {
					goto l1396
				}
				depth--
				add(ruleAscending, position1397)
			}
			return true
		l1396:
			position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
			return false
		},
		/* 120 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action92)> */
		func() bool {
			position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
			{
				position1406 := position
				depth++
				{
					position1407 := position
					depth++
					{
						position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1409
						}
						position++
						goto l1408
					l1409:
						position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
						if buffer[position] != rune('D') {
							goto l1405
						}
						position++
					}
				l1408:
					{
						position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1411
						}
						position++
						goto l1410
					l1411:
						position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
						if buffer[position] != rune('E') {
							goto l1405
						}
						position++
					}
				l1410:
					{
						position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1413
						}
						position++
						goto l1412
					l1413:
						position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
						if buffer[position] != rune('S') {
							goto l1405
						}
						position++
					}
				l1412:
					{
						position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1415
						}
						position++
						goto l1414
					l1415:
						position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
						if buffer[position] != rune('C') {
							goto l1405
						}
						position++
					}
				l1414:
					depth--
					add(rulePegText, position1407)
				}
				if !_rules[ruleAction92]() {
					goto l1405
				}
				depth--
				add(ruleDescending, position1406)
			}
			return true
		l1405:
			position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
			return false
		},
		/* 121 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
			{
				position1417 := position
				depth++
				{
					position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1419
					}
					goto l1418
				l1419:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
					if !_rules[ruleInt]() {
						goto l1420
					}
					goto l1418
				l1420:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
					if !_rules[ruleFloat]() {
						goto l1421
					}
					goto l1418
				l1421:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
					if !_rules[ruleString]() {
						goto l1422
					}
					goto l1418
				l1422:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
					if !_rules[ruleBlob]() {
						goto l1423
					}
					goto l1418
				l1423:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
					if !_rules[ruleTimestamp]() {
						goto l1424
					}
					goto l1418
				l1424:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
					if !_rules[ruleArray]() {
						goto l1425
					}
					goto l1418
				l1425:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
					if !_rules[ruleMap]() {
						goto l1416
					}
				}
			l1418:
				depth--
				add(ruleType, position1417)
			}
			return true
		l1416:
			position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
			return false
		},
		/* 122 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action93)> */
		func() bool {
			position1426, tokenIndex1426, depth1426 := position, tokenIndex, depth
			{
				position1427 := position
				depth++
				{
					position1428 := position
					depth++
					{
						position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1430
						}
						position++
						goto l1429
					l1430:
						position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
						if buffer[position] != rune('B') {
							goto l1426
						}
						position++
					}
				l1429:
					{
						position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1432
						}
						position++
						goto l1431
					l1432:
						position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
						if buffer[position] != rune('O') {
							goto l1426
						}
						position++
					}
				l1431:
					{
						position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1434
						}
						position++
						goto l1433
					l1434:
						position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
						if buffer[position] != rune('O') {
							goto l1426
						}
						position++
					}
				l1433:
					{
						position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1436
						}
						position++
						goto l1435
					l1436:
						position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
						if buffer[position] != rune('L') {
							goto l1426
						}
						position++
					}
				l1435:
					depth--
					add(rulePegText, position1428)
				}
				if !_rules[ruleAction93]() {
					goto l1426
				}
				depth--
				add(ruleBool, position1427)
			}
			return true
		l1426:
			position, tokenIndex, depth = position1426, tokenIndex1426, depth1426
			return false
		},
		/* 123 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action94)> */
		func() bool {
			position1437, tokenIndex1437, depth1437 := position, tokenIndex, depth
			{
				position1438 := position
				depth++
				{
					position1439 := position
					depth++
					{
						position1440, tokenIndex1440, depth1440 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1441
						}
						position++
						goto l1440
					l1441:
						position, tokenIndex, depth = position1440, tokenIndex1440, depth1440
						if buffer[position] != rune('I') {
							goto l1437
						}
						position++
					}
				l1440:
					{
						position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1443
						}
						position++
						goto l1442
					l1443:
						position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
						if buffer[position] != rune('N') {
							goto l1437
						}
						position++
					}
				l1442:
					{
						position1444, tokenIndex1444, depth1444 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1445
						}
						position++
						goto l1444
					l1445:
						position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
						if buffer[position] != rune('T') {
							goto l1437
						}
						position++
					}
				l1444:
					depth--
					add(rulePegText, position1439)
				}
				if !_rules[ruleAction94]() {
					goto l1437
				}
				depth--
				add(ruleInt, position1438)
			}
			return true
		l1437:
			position, tokenIndex, depth = position1437, tokenIndex1437, depth1437
			return false
		},
		/* 124 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action95)> */
		func() bool {
			position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
			{
				position1447 := position
				depth++
				{
					position1448 := position
					depth++
					{
						position1449, tokenIndex1449, depth1449 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1450
						}
						position++
						goto l1449
					l1450:
						position, tokenIndex, depth = position1449, tokenIndex1449, depth1449
						if buffer[position] != rune('F') {
							goto l1446
						}
						position++
					}
				l1449:
					{
						position1451, tokenIndex1451, depth1451 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1452
						}
						position++
						goto l1451
					l1452:
						position, tokenIndex, depth = position1451, tokenIndex1451, depth1451
						if buffer[position] != rune('L') {
							goto l1446
						}
						position++
					}
				l1451:
					{
						position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1454
						}
						position++
						goto l1453
					l1454:
						position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
						if buffer[position] != rune('O') {
							goto l1446
						}
						position++
					}
				l1453:
					{
						position1455, tokenIndex1455, depth1455 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1456
						}
						position++
						goto l1455
					l1456:
						position, tokenIndex, depth = position1455, tokenIndex1455, depth1455
						if buffer[position] != rune('A') {
							goto l1446
						}
						position++
					}
				l1455:
					{
						position1457, tokenIndex1457, depth1457 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1458
						}
						position++
						goto l1457
					l1458:
						position, tokenIndex, depth = position1457, tokenIndex1457, depth1457
						if buffer[position] != rune('T') {
							goto l1446
						}
						position++
					}
				l1457:
					depth--
					add(rulePegText, position1448)
				}
				if !_rules[ruleAction95]() {
					goto l1446
				}
				depth--
				add(ruleFloat, position1447)
			}
			return true
		l1446:
			position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
			return false
		},
		/* 125 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action96)> */
		func() bool {
			position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
			{
				position1460 := position
				depth++
				{
					position1461 := position
					depth++
					{
						position1462, tokenIndex1462, depth1462 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1463
						}
						position++
						goto l1462
					l1463:
						position, tokenIndex, depth = position1462, tokenIndex1462, depth1462
						if buffer[position] != rune('S') {
							goto l1459
						}
						position++
					}
				l1462:
					{
						position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1465
						}
						position++
						goto l1464
					l1465:
						position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
						if buffer[position] != rune('T') {
							goto l1459
						}
						position++
					}
				l1464:
					{
						position1466, tokenIndex1466, depth1466 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1467
						}
						position++
						goto l1466
					l1467:
						position, tokenIndex, depth = position1466, tokenIndex1466, depth1466
						if buffer[position] != rune('R') {
							goto l1459
						}
						position++
					}
				l1466:
					{
						position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1469
						}
						position++
						goto l1468
					l1469:
						position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
						if buffer[position] != rune('I') {
							goto l1459
						}
						position++
					}
				l1468:
					{
						position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1471
						}
						position++
						goto l1470
					l1471:
						position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
						if buffer[position] != rune('N') {
							goto l1459
						}
						position++
					}
				l1470:
					{
						position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1473
						}
						position++
						goto l1472
					l1473:
						position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
						if buffer[position] != rune('G') {
							goto l1459
						}
						position++
					}
				l1472:
					depth--
					add(rulePegText, position1461)
				}
				if !_rules[ruleAction96]() {
					goto l1459
				}
				depth--
				add(ruleString, position1460)
			}
			return true
		l1459:
			position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
			return false
		},
		/* 126 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action97)> */
		func() bool {
			position1474, tokenIndex1474, depth1474 := position, tokenIndex, depth
			{
				position1475 := position
				depth++
				{
					position1476 := position
					depth++
					{
						position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1478
						}
						position++
						goto l1477
					l1478:
						position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
						if buffer[position] != rune('B') {
							goto l1474
						}
						position++
					}
				l1477:
					{
						position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1480
						}
						position++
						goto l1479
					l1480:
						position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
						if buffer[position] != rune('L') {
							goto l1474
						}
						position++
					}
				l1479:
					{
						position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1482
						}
						position++
						goto l1481
					l1482:
						position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
						if buffer[position] != rune('O') {
							goto l1474
						}
						position++
					}
				l1481:
					{
						position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1484
						}
						position++
						goto l1483
					l1484:
						position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
						if buffer[position] != rune('B') {
							goto l1474
						}
						position++
					}
				l1483:
					depth--
					add(rulePegText, position1476)
				}
				if !_rules[ruleAction97]() {
					goto l1474
				}
				depth--
				add(ruleBlob, position1475)
			}
			return true
		l1474:
			position, tokenIndex, depth = position1474, tokenIndex1474, depth1474
			return false
		},
		/* 127 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action98)> */
		func() bool {
			position1485, tokenIndex1485, depth1485 := position, tokenIndex, depth
			{
				position1486 := position
				depth++
				{
					position1487 := position
					depth++
					{
						position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1489
						}
						position++
						goto l1488
					l1489:
						position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
						if buffer[position] != rune('T') {
							goto l1485
						}
						position++
					}
				l1488:
					{
						position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1491
						}
						position++
						goto l1490
					l1491:
						position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
						if buffer[position] != rune('I') {
							goto l1485
						}
						position++
					}
				l1490:
					{
						position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1493
						}
						position++
						goto l1492
					l1493:
						position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
						if buffer[position] != rune('M') {
							goto l1485
						}
						position++
					}
				l1492:
					{
						position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1495
						}
						position++
						goto l1494
					l1495:
						position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
						if buffer[position] != rune('E') {
							goto l1485
						}
						position++
					}
				l1494:
					{
						position1496, tokenIndex1496, depth1496 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1497
						}
						position++
						goto l1496
					l1497:
						position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
						if buffer[position] != rune('S') {
							goto l1485
						}
						position++
					}
				l1496:
					{
						position1498, tokenIndex1498, depth1498 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1499
						}
						position++
						goto l1498
					l1499:
						position, tokenIndex, depth = position1498, tokenIndex1498, depth1498
						if buffer[position] != rune('T') {
							goto l1485
						}
						position++
					}
				l1498:
					{
						position1500, tokenIndex1500, depth1500 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1501
						}
						position++
						goto l1500
					l1501:
						position, tokenIndex, depth = position1500, tokenIndex1500, depth1500
						if buffer[position] != rune('A') {
							goto l1485
						}
						position++
					}
				l1500:
					{
						position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1503
						}
						position++
						goto l1502
					l1503:
						position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
						if buffer[position] != rune('M') {
							goto l1485
						}
						position++
					}
				l1502:
					{
						position1504, tokenIndex1504, depth1504 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1505
						}
						position++
						goto l1504
					l1505:
						position, tokenIndex, depth = position1504, tokenIndex1504, depth1504
						if buffer[position] != rune('P') {
							goto l1485
						}
						position++
					}
				l1504:
					depth--
					add(rulePegText, position1487)
				}
				if !_rules[ruleAction98]() {
					goto l1485
				}
				depth--
				add(ruleTimestamp, position1486)
			}
			return true
		l1485:
			position, tokenIndex, depth = position1485, tokenIndex1485, depth1485
			return false
		},
		/* 128 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action99)> */
		func() bool {
			position1506, tokenIndex1506, depth1506 := position, tokenIndex, depth
			{
				position1507 := position
				depth++
				{
					position1508 := position
					depth++
					{
						position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1510
						}
						position++
						goto l1509
					l1510:
						position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
						if buffer[position] != rune('A') {
							goto l1506
						}
						position++
					}
				l1509:
					{
						position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1512
						}
						position++
						goto l1511
					l1512:
						position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
						if buffer[position] != rune('R') {
							goto l1506
						}
						position++
					}
				l1511:
					{
						position1513, tokenIndex1513, depth1513 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1514
						}
						position++
						goto l1513
					l1514:
						position, tokenIndex, depth = position1513, tokenIndex1513, depth1513
						if buffer[position] != rune('R') {
							goto l1506
						}
						position++
					}
				l1513:
					{
						position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1516
						}
						position++
						goto l1515
					l1516:
						position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
						if buffer[position] != rune('A') {
							goto l1506
						}
						position++
					}
				l1515:
					{
						position1517, tokenIndex1517, depth1517 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1518
						}
						position++
						goto l1517
					l1518:
						position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
						if buffer[position] != rune('Y') {
							goto l1506
						}
						position++
					}
				l1517:
					depth--
					add(rulePegText, position1508)
				}
				if !_rules[ruleAction99]() {
					goto l1506
				}
				depth--
				add(ruleArray, position1507)
			}
			return true
		l1506:
			position, tokenIndex, depth = position1506, tokenIndex1506, depth1506
			return false
		},
		/* 129 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action100)> */
		func() bool {
			position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
			{
				position1520 := position
				depth++
				{
					position1521 := position
					depth++
					{
						position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1523
						}
						position++
						goto l1522
					l1523:
						position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
						if buffer[position] != rune('M') {
							goto l1519
						}
						position++
					}
				l1522:
					{
						position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1525
						}
						position++
						goto l1524
					l1525:
						position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
						if buffer[position] != rune('A') {
							goto l1519
						}
						position++
					}
				l1524:
					{
						position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1527
						}
						position++
						goto l1526
					l1527:
						position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
						if buffer[position] != rune('P') {
							goto l1519
						}
						position++
					}
				l1526:
					depth--
					add(rulePegText, position1521)
				}
				if !_rules[ruleAction100]() {
					goto l1519
				}
				depth--
				add(ruleMap, position1520)
			}
			return true
		l1519:
			position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
			return false
		},
		/* 130 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action101)> */
		func() bool {
			position1528, tokenIndex1528, depth1528 := position, tokenIndex, depth
			{
				position1529 := position
				depth++
				{
					position1530 := position
					depth++
					{
						position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1532
						}
						position++
						goto l1531
					l1532:
						position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
						if buffer[position] != rune('O') {
							goto l1528
						}
						position++
					}
				l1531:
					{
						position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1534
						}
						position++
						goto l1533
					l1534:
						position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
						if buffer[position] != rune('R') {
							goto l1528
						}
						position++
					}
				l1533:
					depth--
					add(rulePegText, position1530)
				}
				if !_rules[ruleAction101]() {
					goto l1528
				}
				depth--
				add(ruleOr, position1529)
			}
			return true
		l1528:
			position, tokenIndex, depth = position1528, tokenIndex1528, depth1528
			return false
		},
		/* 131 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action102)> */
		func() bool {
			position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
			{
				position1536 := position
				depth++
				{
					position1537 := position
					depth++
					{
						position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1539
						}
						position++
						goto l1538
					l1539:
						position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
						if buffer[position] != rune('A') {
							goto l1535
						}
						position++
					}
				l1538:
					{
						position1540, tokenIndex1540, depth1540 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1541
						}
						position++
						goto l1540
					l1541:
						position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
						if buffer[position] != rune('N') {
							goto l1535
						}
						position++
					}
				l1540:
					{
						position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1543
						}
						position++
						goto l1542
					l1543:
						position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
						if buffer[position] != rune('D') {
							goto l1535
						}
						position++
					}
				l1542:
					depth--
					add(rulePegText, position1537)
				}
				if !_rules[ruleAction102]() {
					goto l1535
				}
				depth--
				add(ruleAnd, position1536)
			}
			return true
		l1535:
			position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
			return false
		},
		/* 132 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action103)> */
		func() bool {
			position1544, tokenIndex1544, depth1544 := position, tokenIndex, depth
			{
				position1545 := position
				depth++
				{
					position1546 := position
					depth++
					{
						position1547, tokenIndex1547, depth1547 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1548
						}
						position++
						goto l1547
					l1548:
						position, tokenIndex, depth = position1547, tokenIndex1547, depth1547
						if buffer[position] != rune('N') {
							goto l1544
						}
						position++
					}
				l1547:
					{
						position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1550
						}
						position++
						goto l1549
					l1550:
						position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
						if buffer[position] != rune('O') {
							goto l1544
						}
						position++
					}
				l1549:
					{
						position1551, tokenIndex1551, depth1551 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1552
						}
						position++
						goto l1551
					l1552:
						position, tokenIndex, depth = position1551, tokenIndex1551, depth1551
						if buffer[position] != rune('T') {
							goto l1544
						}
						position++
					}
				l1551:
					depth--
					add(rulePegText, position1546)
				}
				if !_rules[ruleAction103]() {
					goto l1544
				}
				depth--
				add(ruleNot, position1545)
			}
			return true
		l1544:
			position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
			return false
		},
		/* 133 Equal <- <(<'='> Action104)> */
		func() bool {
			position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
			{
				position1554 := position
				depth++
				{
					position1555 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1553
					}
					position++
					depth--
					add(rulePegText, position1555)
				}
				if !_rules[ruleAction104]() {
					goto l1553
				}
				depth--
				add(ruleEqual, position1554)
			}
			return true
		l1553:
			position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
			return false
		},
		/* 134 Less <- <(<'<'> Action105)> */
		func() bool {
			position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
			{
				position1557 := position
				depth++
				{
					position1558 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1556
					}
					position++
					depth--
					add(rulePegText, position1558)
				}
				if !_rules[ruleAction105]() {
					goto l1556
				}
				depth--
				add(ruleLess, position1557)
			}
			return true
		l1556:
			position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
			return false
		},
		/* 135 LessOrEqual <- <(<('<' '=')> Action106)> */
		func() bool {
			position1559, tokenIndex1559, depth1559 := position, tokenIndex, depth
			{
				position1560 := position
				depth++
				{
					position1561 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1559
					}
					position++
					if buffer[position] != rune('=') {
						goto l1559
					}
					position++
					depth--
					add(rulePegText, position1561)
				}
				if !_rules[ruleAction106]() {
					goto l1559
				}
				depth--
				add(ruleLessOrEqual, position1560)
			}
			return true
		l1559:
			position, tokenIndex, depth = position1559, tokenIndex1559, depth1559
			return false
		},
		/* 136 Greater <- <(<'>'> Action107)> */
		func() bool {
			position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
			{
				position1563 := position
				depth++
				{
					position1564 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1562
					}
					position++
					depth--
					add(rulePegText, position1564)
				}
				if !_rules[ruleAction107]() {
					goto l1562
				}
				depth--
				add(ruleGreater, position1563)
			}
			return true
		l1562:
			position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
			return false
		},
		/* 137 GreaterOrEqual <- <(<('>' '=')> Action108)> */
		func() bool {
			position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
			{
				position1566 := position
				depth++
				{
					position1567 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1565
					}
					position++
					if buffer[position] != rune('=') {
						goto l1565
					}
					position++
					depth--
					add(rulePegText, position1567)
				}
				if !_rules[ruleAction108]() {
					goto l1565
				}
				depth--
				add(ruleGreaterOrEqual, position1566)
			}
			return true
		l1565:
			position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
			return false
		},
		/* 138 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action109)> */
		func() bool {
			position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
			{
				position1569 := position
				depth++
				{
					position1570 := position
					depth++
					{
						position1571, tokenIndex1571, depth1571 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1572
						}
						position++
						if buffer[position] != rune('=') {
							goto l1572
						}
						position++
						goto l1571
					l1572:
						position, tokenIndex, depth = position1571, tokenIndex1571, depth1571
						if buffer[position] != rune('<') {
							goto l1568
						}
						position++
						if buffer[position] != rune('>') {
							goto l1568
						}
						position++
					}
				l1571:
					depth--
					add(rulePegText, position1570)
				}
				if !_rules[ruleAction109]() {
					goto l1568
				}
				depth--
				add(ruleNotEqual, position1569)
			}
			return true
		l1568:
			position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
			return false
		},
		/* 139 Concat <- <(<('|' '|')> Action110)> */
		func() bool {
			position1573, tokenIndex1573, depth1573 := position, tokenIndex, depth
			{
				position1574 := position
				depth++
				{
					position1575 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1573
					}
					position++
					if buffer[position] != rune('|') {
						goto l1573
					}
					position++
					depth--
					add(rulePegText, position1575)
				}
				if !_rules[ruleAction110]() {
					goto l1573
				}
				depth--
				add(ruleConcat, position1574)
			}
			return true
		l1573:
			position, tokenIndex, depth = position1573, tokenIndex1573, depth1573
			return false
		},
		/* 140 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action111)> */
		func() bool {
			position1576, tokenIndex1576, depth1576 := position, tokenIndex, depth
			{
				position1577 := position
				depth++
				{
					position1578 := position
					depth++
					{
						position1579, tokenIndex1579, depth1579 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1580
						}
						position++
						goto l1579
					l1580:
						position, tokenIndex, depth = position1579, tokenIndex1579, depth1579
						if buffer[position] != rune('I') {
							goto l1576
						}
						position++
					}
				l1579:
					{
						position1581, tokenIndex1581, depth1581 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1582
						}
						position++
						goto l1581
					l1582:
						position, tokenIndex, depth = position1581, tokenIndex1581, depth1581
						if buffer[position] != rune('S') {
							goto l1576
						}
						position++
					}
				l1581:
					depth--
					add(rulePegText, position1578)
				}
				if !_rules[ruleAction111]() {
					goto l1576
				}
				depth--
				add(ruleIs, position1577)
			}
			return true
		l1576:
			position, tokenIndex, depth = position1576, tokenIndex1576, depth1576
			return false
		},
		/* 141 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action112)> */
		func() bool {
			position1583, tokenIndex1583, depth1583 := position, tokenIndex, depth
			{
				position1584 := position
				depth++
				{
					position1585 := position
					depth++
					{
						position1586, tokenIndex1586, depth1586 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1587
						}
						position++
						goto l1586
					l1587:
						position, tokenIndex, depth = position1586, tokenIndex1586, depth1586
						if buffer[position] != rune('I') {
							goto l1583
						}
						position++
					}
				l1586:
					{
						position1588, tokenIndex1588, depth1588 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1589
						}
						position++
						goto l1588
					l1589:
						position, tokenIndex, depth = position1588, tokenIndex1588, depth1588
						if buffer[position] != rune('S') {
							goto l1583
						}
						position++
					}
				l1588:
					if !_rules[rulesp]() {
						goto l1583
					}
					{
						position1590, tokenIndex1590, depth1590 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1591
						}
						position++
						goto l1590
					l1591:
						position, tokenIndex, depth = position1590, tokenIndex1590, depth1590
						if buffer[position] != rune('N') {
							goto l1583
						}
						position++
					}
				l1590:
					{
						position1592, tokenIndex1592, depth1592 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1593
						}
						position++
						goto l1592
					l1593:
						position, tokenIndex, depth = position1592, tokenIndex1592, depth1592
						if buffer[position] != rune('O') {
							goto l1583
						}
						position++
					}
				l1592:
					{
						position1594, tokenIndex1594, depth1594 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1595
						}
						position++
						goto l1594
					l1595:
						position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
						if buffer[position] != rune('T') {
							goto l1583
						}
						position++
					}
				l1594:
					depth--
					add(rulePegText, position1585)
				}
				if !_rules[ruleAction112]() {
					goto l1583
				}
				depth--
				add(ruleIsNot, position1584)
			}
			return true
		l1583:
			position, tokenIndex, depth = position1583, tokenIndex1583, depth1583
			return false
		},
		/* 142 Plus <- <(<'+'> Action113)> */
		func() bool {
			position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
			{
				position1597 := position
				depth++
				{
					position1598 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1596
					}
					position++
					depth--
					add(rulePegText, position1598)
				}
				if !_rules[ruleAction113]() {
					goto l1596
				}
				depth--
				add(rulePlus, position1597)
			}
			return true
		l1596:
			position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
			return false
		},
		/* 143 Minus <- <(<'-'> Action114)> */
		func() bool {
			position1599, tokenIndex1599, depth1599 := position, tokenIndex, depth
			{
				position1600 := position
				depth++
				{
					position1601 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1599
					}
					position++
					depth--
					add(rulePegText, position1601)
				}
				if !_rules[ruleAction114]() {
					goto l1599
				}
				depth--
				add(ruleMinus, position1600)
			}
			return true
		l1599:
			position, tokenIndex, depth = position1599, tokenIndex1599, depth1599
			return false
		},
		/* 144 Multiply <- <(<'*'> Action115)> */
		func() bool {
			position1602, tokenIndex1602, depth1602 := position, tokenIndex, depth
			{
				position1603 := position
				depth++
				{
					position1604 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1602
					}
					position++
					depth--
					add(rulePegText, position1604)
				}
				if !_rules[ruleAction115]() {
					goto l1602
				}
				depth--
				add(ruleMultiply, position1603)
			}
			return true
		l1602:
			position, tokenIndex, depth = position1602, tokenIndex1602, depth1602
			return false
		},
		/* 145 Divide <- <(<'/'> Action116)> */
		func() bool {
			position1605, tokenIndex1605, depth1605 := position, tokenIndex, depth
			{
				position1606 := position
				depth++
				{
					position1607 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1605
					}
					position++
					depth--
					add(rulePegText, position1607)
				}
				if !_rules[ruleAction116]() {
					goto l1605
				}
				depth--
				add(ruleDivide, position1606)
			}
			return true
		l1605:
			position, tokenIndex, depth = position1605, tokenIndex1605, depth1605
			return false
		},
		/* 146 Modulo <- <(<'%'> Action117)> */
		func() bool {
			position1608, tokenIndex1608, depth1608 := position, tokenIndex, depth
			{
				position1609 := position
				depth++
				{
					position1610 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1608
					}
					position++
					depth--
					add(rulePegText, position1610)
				}
				if !_rules[ruleAction117]() {
					goto l1608
				}
				depth--
				add(ruleModulo, position1609)
			}
			return true
		l1608:
			position, tokenIndex, depth = position1608, tokenIndex1608, depth1608
			return false
		},
		/* 147 UnaryMinus <- <(<'-'> Action118)> */
		func() bool {
			position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
			{
				position1612 := position
				depth++
				{
					position1613 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1611
					}
					position++
					depth--
					add(rulePegText, position1613)
				}
				if !_rules[ruleAction118]() {
					goto l1611
				}
				depth--
				add(ruleUnaryMinus, position1612)
			}
			return true
		l1611:
			position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
			return false
		},
		/* 148 Identifier <- <(<ident> Action119)> */
		func() bool {
			position1614, tokenIndex1614, depth1614 := position, tokenIndex, depth
			{
				position1615 := position
				depth++
				{
					position1616 := position
					depth++
					if !_rules[ruleident]() {
						goto l1614
					}
					depth--
					add(rulePegText, position1616)
				}
				if !_rules[ruleAction119]() {
					goto l1614
				}
				depth--
				add(ruleIdentifier, position1615)
			}
			return true
		l1614:
			position, tokenIndex, depth = position1614, tokenIndex1614, depth1614
			return false
		},
		/* 149 TargetIdentifier <- <(<('*' / jsonSetPath)> Action120)> */
		func() bool {
			position1617, tokenIndex1617, depth1617 := position, tokenIndex, depth
			{
				position1618 := position
				depth++
				{
					position1619 := position
					depth++
					{
						position1620, tokenIndex1620, depth1620 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l1621
						}
						position++
						goto l1620
					l1621:
						position, tokenIndex, depth = position1620, tokenIndex1620, depth1620
						if !_rules[rulejsonSetPath]() {
							goto l1617
						}
					}
				l1620:
					depth--
					add(rulePegText, position1619)
				}
				if !_rules[ruleAction120]() {
					goto l1617
				}
				depth--
				add(ruleTargetIdentifier, position1618)
			}
			return true
		l1617:
			position, tokenIndex, depth = position1617, tokenIndex1617, depth1617
			return false
		},
		/* 150 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
			{
				position1623 := position
				depth++
				{
					position1624, tokenIndex1624, depth1624 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1625
					}
					position++
					goto l1624
				l1625:
					position, tokenIndex, depth = position1624, tokenIndex1624, depth1624
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1622
					}
					position++
				}
			l1624:
			l1626:
				{
					position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
					{
						position1628, tokenIndex1628, depth1628 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1629
						}
						position++
						goto l1628
					l1629:
						position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1630
						}
						position++
						goto l1628
					l1630:
						position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1631
						}
						position++
						goto l1628
					l1631:
						position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
						if buffer[position] != rune('_') {
							goto l1627
						}
						position++
					}
				l1628:
					goto l1626
				l1627:
					position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
				}
				depth--
				add(ruleident, position1623)
			}
			return true
		l1622:
			position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
			return false
		},
		/* 151 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1632, tokenIndex1632, depth1632 := position, tokenIndex, depth
			{
				position1633 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1632
				}
			l1634:
				{
					position1635, tokenIndex1635, depth1635 := position, tokenIndex, depth
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1635
					}
					goto l1634
				l1635:
					position, tokenIndex, depth = position1635, tokenIndex1635, depth1635
				}
				depth--
				add(rulejsonGetPath, position1633)
			}
			return true
		l1632:
			position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
			return false
		},
		/* 152 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1636, tokenIndex1636, depth1636 := position, tokenIndex, depth
			{
				position1637 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1636
				}
			l1638:
				{
					position1639, tokenIndex1639, depth1639 := position, tokenIndex, depth
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1639
					}
					goto l1638
				l1639:
					position, tokenIndex, depth = position1639, tokenIndex1639, depth1639
				}
				depth--
				add(rulejsonSetPath, position1637)
			}
			return true
		l1636:
			position, tokenIndex, depth = position1636, tokenIndex1636, depth1636
			return false
		},
		/* 153 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1640, tokenIndex1640, depth1640 := position, tokenIndex, depth
			{
				position1641 := position
				depth++
				{
					position1642, tokenIndex1642, depth1642 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1643
					}
					goto l1642
				l1643:
					position, tokenIndex, depth = position1642, tokenIndex1642, depth1642
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1640
					}
				}
			l1642:
				depth--
				add(rulejsonPathHead, position1641)
			}
			return true
		l1640:
			position, tokenIndex, depth = position1640, tokenIndex1640, depth1640
			return false
		},
		/* 154 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1644, tokenIndex1644, depth1644 := position, tokenIndex, depth
			{
				position1645 := position
				depth++
				{
					position1646, tokenIndex1646, depth1646 := position, tokenIndex, depth
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1647
					}
					goto l1646
				l1647:
					position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1648
					}
					goto l1646
				l1648:
					position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1649
					}
					goto l1646
				l1649:
					position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1650
					}
					goto l1646
				l1650:
					position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
					if !_rules[rulejsonArraySlice]() {
						goto l1651
					}
					goto l1646
				l1651:
					position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
					if !_rules[rulejsonArrayAccess]() {
						goto l1644
					}
				}
			l1646:
				depth--
				add(rulejsonGetPathNonHead, position1645)
			}
			return true
		l1644:
			position, tokenIndex, depth = position1644, tokenIndex1644, depth1644
			return false
		},
		/* 155 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
			{
				position1653 := position
				depth++
				{
					position1654, tokenIndex1654, depth1654 := position, tokenIndex, depth
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1655
					}
					goto l1654
				l1655:
					position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1652
					}
				}
			l1654:
				depth--
				add(rulejsonSetPathNonHead, position1653)
			}
			return true
		l1652:
			position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
			return false
		},
		/* 156 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1656, tokenIndex1656, depth1656 := position, tokenIndex, depth
			{
				position1657 := position
				depth++
				{
					position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1659
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1659
					}
					goto l1658
				l1659:
					position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1656
					}
				}
			l1658:
				depth--
				add(rulejsonMapSingleLevel, position1657)
			}
			return true
		l1656:
			position, tokenIndex, depth = position1656, tokenIndex1656, depth1656
			return false
		},
		/* 157 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
			{
				position1661 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1660
				}
				position++
				if buffer[position] != rune('.') {
					goto l1660
				}
				position++
				{
					position1662, tokenIndex1662, depth1662 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1663
					}
					goto l1662
				l1663:
					position, tokenIndex, depth = position1662, tokenIndex1662, depth1662
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1660
					}
				}
			l1662:
				depth--
				add(rulejsonMapMultipleLevel, position1661)
			}
			return true
		l1660:
			position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
			return false
		},
		/* 158 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1664, tokenIndex1664, depth1664 := position, tokenIndex, depth
			{
				position1665 := position
				depth++
				{
					position1666 := position
					depth++
					{
						position1667, tokenIndex1667, depth1667 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1668
						}
						position++
						goto l1667
					l1668:
						position, tokenIndex, depth = position1667, tokenIndex1667, depth1667
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1664
						}
						position++
					}
				l1667:
				l1669:
					{
						position1670, tokenIndex1670, depth1670 := position, tokenIndex, depth
						{
							position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1672
							}
							position++
							goto l1671
						l1672:
							position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1673
							}
							position++
							goto l1671
						l1673:
							position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1674
							}
							position++
							goto l1671
						l1674:
							position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
							if buffer[position] != rune('_') {
								goto l1670
							}
							position++
						}
					l1671:
						goto l1669
					l1670:
						position, tokenIndex, depth = position1670, tokenIndex1670, depth1670
					}
					depth--
					add(rulePegText, position1666)
				}
				depth--
				add(rulejsonMapAccessString, position1665)
			}
			return true
		l1664:
			position, tokenIndex, depth = position1664, tokenIndex1664, depth1664
			return false
		},
		/* 159 jsonMapAccessBracket <- <('[' singleQuotedString ']')> */
		func() bool {
			position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
			{
				position1676 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1675
				}
				position++
				if !_rules[rulesingleQuotedString]() {
					goto l1675
				}
				if buffer[position] != rune(']') {
					goto l1675
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1676)
			}
			return true
		l1675:
			position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
			return false
		},
		/* 160 singleQuotedString <- <('\'' <(('\'' '\'') / (!'\'' .))*> '\'')> */
		func() bool {
			position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
			{
				position1678 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l1677
				}
				position++
				{
					position1679 := position
					depth++
				l1680:
					{
						position1681, tokenIndex1681, depth1681 := position, tokenIndex, depth
						{
							position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1683
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1683
							}
							position++
							goto l1682
						l1683:
							position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
							{
								position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1684
								}
								position++
								goto l1681
							l1684:
								position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
							}
							if !matchDot() {
								goto l1681
							}
						}
					l1682:
						goto l1680
					l1681:
						position, tokenIndex, depth = position1681, tokenIndex1681, depth1681
					}
					depth--
					add(rulePegText, position1679)
				}
				if buffer[position] != rune('\'') {
					goto l1677
				}
				position++
				depth--
				add(rulesingleQuotedString, position1678)
			}
			return true
		l1677:
			position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
			return false
		},
		/* 161 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1685, tokenIndex1685, depth1685 := position, tokenIndex, depth
			{
				position1686 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1685
				}
				position++
				{
					position1687 := position
					depth++
					{
						position1688, tokenIndex1688, depth1688 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1688
						}
						position++
						goto l1689
					l1688:
						position, tokenIndex, depth = position1688, tokenIndex1688, depth1688
					}
				l1689:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1685
					}
					position++
				l1690:
					{
						position1691, tokenIndex1691, depth1691 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1691
						}
						position++
						goto l1690
					l1691:
						position, tokenIndex, depth = position1691, tokenIndex1691, depth1691
					}
					depth--
					add(rulePegText, position1687)
				}
				if buffer[position] != rune(']') {
					goto l1685
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1686)
			}
			return true
		l1685:
			position, tokenIndex, depth = position1685, tokenIndex1685, depth1685
			return false
		},
		/* 162 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
		func() bool {
			position1692, tokenIndex1692, depth1692 := position, tokenIndex, depth
			{
				position1693 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1692
				}
				position++
				{
					position1694 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1692
					}
					position++
				l1695:
					{
						position1696, tokenIndex1696, depth1696 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1696
						}
						position++
						goto l1695
					l1696:
						position, tokenIndex, depth = position1696, tokenIndex1696, depth1696
					}
					depth--
					add(rulePegText, position1694)
				}
				if buffer[position] != rune(']') {
					goto l1692
				}
				position++
				depth--
				add(rulejsonNonNegativeArrayAccess, position1693)
			}
			return true
		l1692:
			position, tokenIndex, depth = position1692, tokenIndex1692, depth1692
			return false
		},
		/* 163 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1697, tokenIndex1697, depth1697 := position, tokenIndex, depth
			{
				position1698 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1697
				}
				position++
				{
					position1699 := position
					depth++
					{
						position1700, tokenIndex1700, depth1700 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1700
						}
						position++
						goto l1701
					l1700:
						position, tokenIndex, depth = position1700, tokenIndex1700, depth1700
					}
				l1701:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1697
					}
					position++
				l1702:
					{
						position1703, tokenIndex1703, depth1703 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1703
						}
						position++
						goto l1702
					l1703:
						position, tokenIndex, depth = position1703, tokenIndex1703, depth1703
					}
					if buffer[position] != rune(':') {
						goto l1697
					}
					position++
					{
						position1704, tokenIndex1704, depth1704 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1704
						}
						position++
						goto l1705
					l1704:
						position, tokenIndex, depth = position1704, tokenIndex1704, depth1704
					}
				l1705:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1697
					}
					position++
				l1706:
					{
						position1707, tokenIndex1707, depth1707 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1707
						}
						position++
						goto l1706
					l1707:
						position, tokenIndex, depth = position1707, tokenIndex1707, depth1707
					}
					{
						position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1708
						}
						position++
						{
							position1710, tokenIndex1710, depth1710 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1710
							}
							position++
							goto l1711
						l1710:
							position, tokenIndex, depth = position1710, tokenIndex1710, depth1710
						}
					l1711:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1708
						}
						position++
					l1712:
						{
							position1713, tokenIndex1713, depth1713 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1713
							}
							position++
							goto l1712
						l1713:
							position, tokenIndex, depth = position1713, tokenIndex1713, depth1713
						}
						goto l1709
					l1708:
						position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
					}
				l1709:
					depth--
					add(rulePegText, position1699)
				}
				if buffer[position] != rune(']') {
					goto l1697
				}
				position++
				depth--
				add(rulejsonArraySlice, position1698)
			}
			return true
		l1697:
			position, tokenIndex, depth = position1697, tokenIndex1697, depth1697
			return false
		},
		/* 164 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1714, tokenIndex1714, depth1714 := position, tokenIndex, depth
			{
				position1715 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1714
				}
				position++
				{
					position1716 := position
					depth++
					{
						position1717, tokenIndex1717, depth1717 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1718
						}
						position++
						{
							position1719, tokenIndex1719, depth1719 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1719
							}
							position++
							goto l1720
						l1719:
							position, tokenIndex, depth = position1719, tokenIndex1719, depth1719
						}
					l1720:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1718
						}
						position++
					l1721:
						{
							position1722, tokenIndex1722, depth1722 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1722
							}
							position++
							goto l1721
						l1722:
							position, tokenIndex, depth = position1722, tokenIndex1722, depth1722
						}
						goto l1717
					l1718:
						position, tokenIndex, depth = position1717, tokenIndex1717, depth1717
						{
							position1723, tokenIndex1723, depth1723 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1723
							}
							position++
							goto l1724
						l1723:
							position, tokenIndex, depth = position1723, tokenIndex1723, depth1723
						}
					l1724:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1714
						}
						position++
					l1725:
						{
							position1726, tokenIndex1726, depth1726 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1726
							}
							position++
							goto l1725
						l1726:
							position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
						}
						if buffer[position] != rune(':') {
							goto l1714
						}
						position++
					}
				l1717:
					depth--
					add(rulePegText, position1716)
				}
				if buffer[position] != rune(']') {
					goto l1714
				}
				position++
				depth--
				add(rulejsonArrayPartialSlice, position1715)
			}
			return true
		l1714:
			position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
			return false
		},
		/* 165 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1727, tokenIndex1727, depth1727 := position, tokenIndex, depth
			{
				position1728 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1727
				}
				position++
				if buffer[position] != rune(':') {
					goto l1727
				}
				position++
				if buffer[position] != rune(']') {
					goto l1727
				}
				position++
				depth--
				add(rulejsonArrayFullSlice, position1728)
			}
			return true
		l1727:
			position, tokenIndex, depth = position1727, tokenIndex1727, depth1727
			return false
		},
		/* 166 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1729, tokenIndex1729, depth1729 := position, tokenIndex, depth
			{
				position1730 := position
				depth++
				{
					position1731, tokenIndex1731, depth1731 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l1732
					}
					position++
					goto l1731
				l1732:
					position, tokenIndex, depth = position1731, tokenIndex1731, depth1731
					if buffer[position] != rune('\t') {
						goto l1733
					}
					position++
					goto l1731
				l1733:
					position, tokenIndex, depth = position1731, tokenIndex1731, depth1731
					if buffer[position] != rune('\n') {
						goto l1734
					}
					position++
					goto l1731
				l1734:
					position, tokenIndex, depth = position1731, tokenIndex1731, depth1731
					if buffer[position] != rune('\r') {
						goto l1735
					}
					position++
					goto l1731
				l1735:
					position, tokenIndex, depth = position1731, tokenIndex1731, depth1731
					if !_rules[rulecomment]() {
						goto l1736
					}
					goto l1731
				l1736:
					position, tokenIndex, depth = position1731, tokenIndex1731, depth1731
					if !_rules[rulefinalComment]() {
						goto l1729
					}
				}
			l1731:
				depth--
				add(rulespElem, position1730)
			}
			return true
		l1729:
			position, tokenIndex, depth = position1729, tokenIndex1729, depth1729
			return false
		},
		/* 167 sp <- <spElem+> */
		func() bool {
			position1737, tokenIndex1737, depth1737 := position, tokenIndex, depth
			{
				position1738 := position
				depth++
				if !_rules[rulespElem]() {
					goto l1737
				}
			l1739:
				{
					position1740, tokenIndex1740, depth1740 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1740
					}
					goto l1739
				l1740:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
				}
				depth--
				add(rulesp, position1738)
			}
			return true
		l1737:
			position, tokenIndex, depth = position1737, tokenIndex1737, depth1737
			return false
		},
		/* 168 spOpt <- <spElem*> */
		func() bool {
			{
				position1742 := position
				depth++
			l1743:
				{
					position1744, tokenIndex1744, depth1744 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1744
					}
					goto l1743
				l1744:
					position, tokenIndex, depth = position1744, tokenIndex1744, depth1744
				}
				depth--
				add(rulespOpt, position1742)
			}
			return true
		},
		/* 169 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1745, tokenIndex1745, depth1745 := position, tokenIndex, depth
			{
				position1746 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1745
				}
				position++
				if buffer[position] != rune('-') {
					goto l1745
				}
				position++
			l1747:
				{
					position1748, tokenIndex1748, depth1748 := position, tokenIndex, depth
					{
						position1749, tokenIndex1749, depth1749 := position, tokenIndex, depth
						{
							position1750, tokenIndex1750, depth1750 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1751
							}
							position++
							goto l1750
						l1751:
							position, tokenIndex, depth = position1750, tokenIndex1750, depth1750
							if buffer[position] != rune('\n') {
								goto l1749
							}
							position++
						}
					l1750:
						goto l1748
					l1749:
						position, tokenIndex, depth = position1749, tokenIndex1749, depth1749
					}
					if !matchDot() {
						goto l1748
					}
					goto l1747
				l1748:
					position, tokenIndex, depth = position1748, tokenIndex1748, depth1748
				}
				{
					position1752, tokenIndex1752, depth1752 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1753
					}
					position++
					goto l1752
				l1753:
					position, tokenIndex, depth = position1752, tokenIndex1752, depth1752
					if buffer[position] != rune('\n') {
						goto l1745
					}
					position++
				}
			l1752:
				depth--
				add(rulecomment, position1746)
			}
			return true
		l1745:
			position, tokenIndex, depth = position1745, tokenIndex1745, depth1745
			return false
		},
		/* 170 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1754, tokenIndex1754, depth1754 := position, tokenIndex, depth
			{
				position1755 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1754
				}
				position++
				if buffer[position] != rune('-') {
					goto l1754
				}
				position++
			l1756:
				{
					position1757, tokenIndex1757, depth1757 := position, tokenIndex, depth
					{
						position1758, tokenIndex1758, depth1758 := position, tokenIndex, depth
						{
							position1759, tokenIndex1759, depth1759 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1760
							}
							position++
							goto l1759
						l1760:
							position, tokenIndex, depth = position1759, tokenIndex1759, depth1759
							if buffer[position] != rune('\n') {
								goto l1758
							}
							position++
						}
					l1759:
						goto l1757
					l1758:
						position, tokenIndex, depth = position1758, tokenIndex1758, depth1758
					}
					if !matchDot() {
						goto l1757
					}
					goto l1756
				l1757:
					position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
				}
				{
					position1761, tokenIndex1761, depth1761 := position, tokenIndex, depth
					if !matchDot() {
						goto l1761
					}
					goto l1754
				l1761:
					position, tokenIndex, depth = position1761, tokenIndex1761, depth1761
				}
				depth--
				add(rulefinalComment, position1755)
			}
			return true
		l1754:
			position, tokenIndex, depth = position1754, tokenIndex1754, depth1754
			return false
		},
		nil,
		/* 173 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 174 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 175 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 176 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 177 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 178 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 179 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 180 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 181 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 182 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 183 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 184 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 185 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 186 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 187 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 188 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 189 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 190 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 191 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 192 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 193 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 194 Action21 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 195 Action22 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 196 Action23 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 197 Action24 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 198 Action25 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 199 Action26 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 200 Action27 <- <{
		    p.AssembleEmitterSampling(CountBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 201 Action28 <- <{
		    p.AssembleEmitterSampling(RandomizedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 202 Action29 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 203 Action30 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 0.001)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 204 Action31 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 205 Action32 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 206 Action33 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 207 Action34 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 208 Action35 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 209 Action36 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 210 Action37 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 211 Action38 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 212 Action39 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 213 Action40 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 214 Action41 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 215 Action42 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 216 Action43 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 217 Action44 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 218 Action45 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 219 Action46 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 220 Action47 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 221 Action48 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 222 Action49 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 223 Action50 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 224 Action51 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 225 Action52 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 226 Action53 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 227 Action54 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 228 Action55 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 229 Action56 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 230 Action57 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 231 Action58 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 232 Action59 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 233 Action60 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 234 Action61 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 235 Action62 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 236 Action63 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 237 Action64 <- <{
		    p.AssembleSortedExpression()
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 238 Action65 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 239 Action66 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 240 Action67 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 241 Action68 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 242 Action69 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 243 Action70 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 244 Action71 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 245 Action72 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 246 Action73 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 247 Action74 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 248 Action75 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 249 Action76 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 250 Action77 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 251 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 252 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 253 Action80 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 254 Action81 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 255 Action82 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 256 Action83 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 257 Action84 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 258 Action85 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 259 Action86 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 260 Action87 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 261 Action88 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 262 Action89 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 263 Action90 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 264 Action91 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 265 Action92 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 266 Action93 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 267 Action94 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 268 Action95 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 269 Action96 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 270 Action97 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 271 Action98 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 272 Action99 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 273 Action100 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 274 Action101 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 275 Action102 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 276 Action103 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 277 Action104 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 278 Action105 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 279 Action106 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 280 Action107 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 281 Action108 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 282 Action109 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 283 Action110 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction110, position)
			}
			return true
		},
		/* 284 Action111 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction111, position)
			}
			return true
		},
		/* 285 Action112 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction112, position)
			}
			return true
		},
		/* 286 Action113 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction113, position)
			}
			return true
		},
		/* 287 Action114 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction114, position)
			}
			return true
		},
		/* 288 Action115 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction115, position)
			}
			return true
		},
		/* 289 Action116 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction116, position)
			}
			return true
		},
		/* 290 Action117 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction117, position)
			}
			return true
		},
		/* 291 Action118 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction118, position)
			}
			return true
		},
		/* 292 Action119 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction119, position)
			}
			return true
		},
		/* 293 Action120 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction120, position)
			}
			return true
		},
	}
	p.rules = _rules
}

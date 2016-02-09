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
	ruleEvalStmt
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
	ruleCapacitySpecOpt
	ruleSheddingSpecOpt
	ruleSheddingOption
	ruleSourceSinkSpecs
	ruleUpdateSourceSinkSpecs
	ruleSetOptSpecs
	ruleStateTagOpt
	ruleSourceSinkParam
	ruleSourceSinkParamVal
	ruleParamLiteral
	ruleParamArrayExpr
	ruleParamMapExpr
	ruleParamKeyValuePair
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
	ruleCase
	ruleConditionCase
	ruleExpressionCase
	ruleWhenThenPair
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
	ruleNonNegativeNumericLiteral
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
	ruleWait
	ruleDropOldest
	ruleDropNewest
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
	ruleAction121
	ruleAction122
	ruleAction123
	ruleAction124
	ruleAction125
	ruleAction126
	ruleAction127
	ruleAction128
	ruleAction129
	ruleAction130
	ruleAction131
	ruleAction132
	ruleAction133

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
	"EvalStmt",
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
	"CapacitySpecOpt",
	"SheddingSpecOpt",
	"SheddingOption",
	"SourceSinkSpecs",
	"UpdateSourceSinkSpecs",
	"SetOptSpecs",
	"StateTagOpt",
	"SourceSinkParam",
	"SourceSinkParamVal",
	"ParamLiteral",
	"ParamArrayExpr",
	"ParamMapExpr",
	"ParamKeyValuePair",
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
	"Case",
	"ConditionCase",
	"ExpressionCase",
	"WhenThenPair",
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
	"NonNegativeNumericLiteral",
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
	"Wait",
	"DropOldest",
	"DropNewest",
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
	"Action121",
	"Action122",
	"Action123",
	"Action124",
	"Action125",
	"Action126",
	"Action127",
	"Action128",
	"Action129",
	"Action130",
	"Action131",
	"Action132",
	"Action133",

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
	rules  [322]func() bool
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

			p.AssembleEval(begin, end)

		case ruleAction25:

			p.AssembleEmitter()

		case ruleAction26:

			p.AssembleEmitterOptions(begin, end)

		case ruleAction27:

			p.AssembleEmitterLimit()

		case ruleAction28:

			p.AssembleEmitterSampling(CountBasedSampling, 1)

		case ruleAction29:

			p.AssembleEmitterSampling(RandomizedSampling, 1)

		case ruleAction30:

			p.AssembleEmitterSampling(TimeBasedSampling, 1)

		case ruleAction31:

			p.AssembleEmitterSampling(TimeBasedSampling, 0.001)

		case ruleAction32:

			p.AssembleProjections(begin, end)

		case ruleAction33:

			p.AssembleAlias()

		case ruleAction34:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction35:

			p.AssembleInterval()

		case ruleAction36:

			p.AssembleInterval()

		case ruleAction37:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction38:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction39:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction40:

			p.EnsureAliasedStreamWindow()

		case ruleAction41:

			p.AssembleAliasedStreamWindow()

		case ruleAction42:

			p.AssembleStreamWindow()

		case ruleAction43:

			p.AssembleUDSFFuncApp()

		case ruleAction44:

			p.EnsureCapacitySpec(begin, end)

		case ruleAction45:

			p.EnsureSheddingSpec(begin, end)

		case ruleAction46:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction47:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction48:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction49:

			p.EnsureIdentifier(begin, end)

		case ruleAction50:

			p.AssembleSourceSinkParam()

		case ruleAction51:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction52:

			p.AssembleMap(begin, end)

		case ruleAction53:

			p.AssembleKeyValuePair()

		case ruleAction54:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction55:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction56:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction57:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction58:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction59:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction60:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction61:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction62:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction63:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction64:

			p.AssembleTypeCast(begin, end)

		case ruleAction65:

			p.AssembleTypeCast(begin, end)

		case ruleAction66:

			p.AssembleFuncApp()

		case ruleAction67:

			p.AssembleExpressions(begin, end)
			p.AssembleFuncApp()

		case ruleAction68:

			p.AssembleExpressions(begin, end)

		case ruleAction69:

			p.AssembleExpressions(begin, end)

		case ruleAction70:

			p.AssembleSortedExpression()

		case ruleAction71:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction72:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction73:

			p.AssembleMap(begin, end)

		case ruleAction74:

			p.AssembleKeyValuePair()

		case ruleAction75:

			p.AssembleConditionCase(begin, end)

		case ruleAction76:

			p.AssembleExpressionCase(begin, end)

		case ruleAction77:

			p.AssembleWhenThenPair()

		case ruleAction78:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction80:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction81:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction82:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction83:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction84:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction85:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction86:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction87:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction88:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction89:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction90:

			p.PushComponent(begin, end, Istream)

		case ruleAction91:

			p.PushComponent(begin, end, Dstream)

		case ruleAction92:

			p.PushComponent(begin, end, Rstream)

		case ruleAction93:

			p.PushComponent(begin, end, Tuples)

		case ruleAction94:

			p.PushComponent(begin, end, Seconds)

		case ruleAction95:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction96:

			p.PushComponent(begin, end, Wait)

		case ruleAction97:

			p.PushComponent(begin, end, DropOldest)

		case ruleAction98:

			p.PushComponent(begin, end, DropNewest)

		case ruleAction99:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction100:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction101:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction102:

			p.PushComponent(begin, end, Yes)

		case ruleAction103:

			p.PushComponent(begin, end, No)

		case ruleAction104:

			p.PushComponent(begin, end, Yes)

		case ruleAction105:

			p.PushComponent(begin, end, No)

		case ruleAction106:

			p.PushComponent(begin, end, Bool)

		case ruleAction107:

			p.PushComponent(begin, end, Int)

		case ruleAction108:

			p.PushComponent(begin, end, Float)

		case ruleAction109:

			p.PushComponent(begin, end, String)

		case ruleAction110:

			p.PushComponent(begin, end, Blob)

		case ruleAction111:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction112:

			p.PushComponent(begin, end, Array)

		case ruleAction113:

			p.PushComponent(begin, end, Map)

		case ruleAction114:

			p.PushComponent(begin, end, Or)

		case ruleAction115:

			p.PushComponent(begin, end, And)

		case ruleAction116:

			p.PushComponent(begin, end, Not)

		case ruleAction117:

			p.PushComponent(begin, end, Equal)

		case ruleAction118:

			p.PushComponent(begin, end, Less)

		case ruleAction119:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction120:

			p.PushComponent(begin, end, Greater)

		case ruleAction121:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction122:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction123:

			p.PushComponent(begin, end, Concat)

		case ruleAction124:

			p.PushComponent(begin, end, Is)

		case ruleAction125:

			p.PushComponent(begin, end, IsNot)

		case ruleAction126:

			p.PushComponent(begin, end, Plus)

		case ruleAction127:

			p.PushComponent(begin, end, Minus)

		case ruleAction128:

			p.PushComponent(begin, end, Multiply)

		case ruleAction129:

			p.PushComponent(begin, end, Divide)

		case ruleAction130:

			p.PushComponent(begin, end, Modulo)

		case ruleAction131:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction132:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction133:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
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
		/* 3 Statement <- <(SelectUnionStmt / SelectStmt / SourceStmt / SinkStmt / StateStmt / StreamStmt / EvalStmt)> */
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
						goto l21
					}
					goto l15
				l21:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleEvalStmt]() {
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
			position22, tokenIndex22, depth22 := position, tokenIndex, depth
			{
				position23 := position
				depth++
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if !_rules[ruleCreateSourceStmt]() {
						goto l25
					}
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleUpdateSourceStmt]() {
						goto l26
					}
					goto l24
				l26:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleDropSourceStmt]() {
						goto l27
					}
					goto l24
				l27:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[rulePauseSourceStmt]() {
						goto l28
					}
					goto l24
				l28:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleResumeSourceStmt]() {
						goto l29
					}
					goto l24
				l29:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleRewindSourceStmt]() {
						goto l22
					}
				}
			l24:
				depth--
				add(ruleSourceStmt, position23)
			}
			return true
		l22:
			position, tokenIndex, depth = position22, tokenIndex22, depth22
			return false
		},
		/* 5 SinkStmt <- <(CreateSinkStmt / UpdateSinkStmt / DropSinkStmt)> */
		func() bool {
			position30, tokenIndex30, depth30 := position, tokenIndex, depth
			{
				position31 := position
				depth++
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if !_rules[ruleCreateSinkStmt]() {
						goto l33
					}
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if !_rules[ruleUpdateSinkStmt]() {
						goto l34
					}
					goto l32
				l34:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if !_rules[ruleDropSinkStmt]() {
						goto l30
					}
				}
			l32:
				depth--
				add(ruleSinkStmt, position31)
			}
			return true
		l30:
			position, tokenIndex, depth = position30, tokenIndex30, depth30
			return false
		},
		/* 6 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt / LoadStateOrCreateStmt / LoadStateStmt / SaveStateStmt)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					if !_rules[ruleCreateStateStmt]() {
						goto l38
					}
					goto l37
				l38:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleUpdateStateStmt]() {
						goto l39
					}
					goto l37
				l39:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleDropStateStmt]() {
						goto l40
					}
					goto l37
				l40:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleLoadStateOrCreateStmt]() {
						goto l41
					}
					goto l37
				l41:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleLoadStateStmt]() {
						goto l42
					}
					goto l37
				l42:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleSaveStateStmt]() {
						goto l35
					}
				}
			l37:
				depth--
				add(ruleStateStmt, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 7 StreamStmt <- <(CreateStreamAsSelectUnionStmt / CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoSelectStmt / InsertIntoFromStmt)> */
		func() bool {
			position43, tokenIndex43, depth43 := position, tokenIndex, depth
			{
				position44 := position
				depth++
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l46
					}
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l47
					}
					goto l45
				l47:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if !_rules[ruleDropStreamStmt]() {
						goto l48
					}
					goto l45
				l48:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l49
					}
					goto l45
				l49:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l43
					}
				}
			l45:
				depth--
				add(ruleStreamStmt, position44)
			}
			return true
		l43:
			position, tokenIndex, depth = position43, tokenIndex43, depth43
			return false
		},
		/* 8 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') Emitter Projections WindowedFrom Filter Grouping Having Action2)> */
		func() bool {
			position50, tokenIndex50, depth50 := position, tokenIndex, depth
			{
				position51 := position
				depth++
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
						goto l50
					}
					position++
				}
			l52:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('E') {
						goto l50
					}
					position++
				}
			l54:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('L') {
						goto l50
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
						goto l50
					}
					position++
				}
			l58:
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l61
					}
					position++
					goto l60
				l61:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
					if buffer[position] != rune('C') {
						goto l50
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
						goto l50
					}
					position++
				}
			l62:
				if !_rules[ruleEmitter]() {
					goto l50
				}
				if !_rules[ruleProjections]() {
					goto l50
				}
				if !_rules[ruleWindowedFrom]() {
					goto l50
				}
				if !_rules[ruleFilter]() {
					goto l50
				}
				if !_rules[ruleGrouping]() {
					goto l50
				}
				if !_rules[ruleHaving]() {
					goto l50
				}
				if !_rules[ruleAction2]() {
					goto l50
				}
				depth--
				add(ruleSelectStmt, position51)
			}
			return true
		l50:
			position, tokenIndex, depth = position50, tokenIndex50, depth50
			return false
		},
		/* 9 SelectUnionStmt <- <(<(SelectStmt (sp (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N')) sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action3)> */
		func() bool {
			position64, tokenIndex64, depth64 := position, tokenIndex, depth
			{
				position65 := position
				depth++
				{
					position66 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l64
					}
					if !_rules[rulesp]() {
						goto l64
					}
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l70
						}
						position++
						goto l69
					l70:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
						if buffer[position] != rune('U') {
							goto l64
						}
						position++
					}
				l69:
					{
						position71, tokenIndex71, depth71 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l72
						}
						position++
						goto l71
					l72:
						position, tokenIndex, depth = position71, tokenIndex71, depth71
						if buffer[position] != rune('N') {
							goto l64
						}
						position++
					}
				l71:
					{
						position73, tokenIndex73, depth73 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l74
						}
						position++
						goto l73
					l74:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
						if buffer[position] != rune('I') {
							goto l64
						}
						position++
					}
				l73:
					{
						position75, tokenIndex75, depth75 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l76
						}
						position++
						goto l75
					l76:
						position, tokenIndex, depth = position75, tokenIndex75, depth75
						if buffer[position] != rune('O') {
							goto l64
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
							goto l64
						}
						position++
					}
				l77:
					if !_rules[rulesp]() {
						goto l64
					}
					{
						position79, tokenIndex79, depth79 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l80
						}
						position++
						goto l79
					l80:
						position, tokenIndex, depth = position79, tokenIndex79, depth79
						if buffer[position] != rune('A') {
							goto l64
						}
						position++
					}
				l79:
					{
						position81, tokenIndex81, depth81 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l82
						}
						position++
						goto l81
					l82:
						position, tokenIndex, depth = position81, tokenIndex81, depth81
						if buffer[position] != rune('L') {
							goto l64
						}
						position++
					}
				l81:
					{
						position83, tokenIndex83, depth83 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l84
						}
						position++
						goto l83
					l84:
						position, tokenIndex, depth = position83, tokenIndex83, depth83
						if buffer[position] != rune('L') {
							goto l64
						}
						position++
					}
				l83:
					if !_rules[rulesp]() {
						goto l64
					}
					if !_rules[ruleSelectStmt]() {
						goto l64
					}
				l67:
					{
						position68, tokenIndex68, depth68 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l68
						}
						{
							position85, tokenIndex85, depth85 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l86
							}
							position++
							goto l85
						l86:
							position, tokenIndex, depth = position85, tokenIndex85, depth85
							if buffer[position] != rune('U') {
								goto l68
							}
							position++
						}
					l85:
						{
							position87, tokenIndex87, depth87 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l88
							}
							position++
							goto l87
						l88:
							position, tokenIndex, depth = position87, tokenIndex87, depth87
							if buffer[position] != rune('N') {
								goto l68
							}
							position++
						}
					l87:
						{
							position89, tokenIndex89, depth89 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l90
							}
							position++
							goto l89
						l90:
							position, tokenIndex, depth = position89, tokenIndex89, depth89
							if buffer[position] != rune('I') {
								goto l68
							}
							position++
						}
					l89:
						{
							position91, tokenIndex91, depth91 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l92
							}
							position++
							goto l91
						l92:
							position, tokenIndex, depth = position91, tokenIndex91, depth91
							if buffer[position] != rune('O') {
								goto l68
							}
							position++
						}
					l91:
						{
							position93, tokenIndex93, depth93 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l94
							}
							position++
							goto l93
						l94:
							position, tokenIndex, depth = position93, tokenIndex93, depth93
							if buffer[position] != rune('N') {
								goto l68
							}
							position++
						}
					l93:
						if !_rules[rulesp]() {
							goto l68
						}
						{
							position95, tokenIndex95, depth95 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l96
							}
							position++
							goto l95
						l96:
							position, tokenIndex, depth = position95, tokenIndex95, depth95
							if buffer[position] != rune('A') {
								goto l68
							}
							position++
						}
					l95:
						{
							position97, tokenIndex97, depth97 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l98
							}
							position++
							goto l97
						l98:
							position, tokenIndex, depth = position97, tokenIndex97, depth97
							if buffer[position] != rune('L') {
								goto l68
							}
							position++
						}
					l97:
						{
							position99, tokenIndex99, depth99 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l100
							}
							position++
							goto l99
						l100:
							position, tokenIndex, depth = position99, tokenIndex99, depth99
							if buffer[position] != rune('L') {
								goto l68
							}
							position++
						}
					l99:
						if !_rules[rulesp]() {
							goto l68
						}
						if !_rules[ruleSelectStmt]() {
							goto l68
						}
						goto l67
					l68:
						position, tokenIndex, depth = position68, tokenIndex68, depth68
					}
					depth--
					add(rulePegText, position66)
				}
				if !_rules[ruleAction3]() {
					goto l64
				}
				depth--
				add(ruleSelectUnionStmt, position65)
			}
			return true
		l64:
			position, tokenIndex, depth = position64, tokenIndex64, depth64
			return false
		},
		/* 10 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action4)> */
		func() bool {
			position101, tokenIndex101, depth101 := position, tokenIndex, depth
			{
				position102 := position
				depth++
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
						goto l101
					}
					position++
				}
			l103:
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('R') {
						goto l101
					}
					position++
				}
			l105:
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l108
					}
					position++
					goto l107
				l108:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
					if buffer[position] != rune('E') {
						goto l101
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('A') {
						goto l101
					}
					position++
				}
			l109:
				{
					position111, tokenIndex111, depth111 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l112
					}
					position++
					goto l111
				l112:
					position, tokenIndex, depth = position111, tokenIndex111, depth111
					if buffer[position] != rune('T') {
						goto l101
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
						goto l101
					}
					position++
				}
			l113:
				if !_rules[rulesp]() {
					goto l101
				}
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('S') {
						goto l101
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('T') {
						goto l101
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
						goto l101
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
						goto l101
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
						goto l101
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('M') {
						goto l101
					}
					position++
				}
			l125:
				if !_rules[rulesp]() {
					goto l101
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l101
				}
				if !_rules[rulesp]() {
					goto l101
				}
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('A') {
						goto l101
					}
					position++
				}
			l127:
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
						goto l101
					}
					position++
				}
			l129:
				if !_rules[rulesp]() {
					goto l101
				}
				if !_rules[ruleSelectStmt]() {
					goto l101
				}
				if !_rules[ruleAction4]() {
					goto l101
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position102)
			}
			return true
		l101:
			position, tokenIndex, depth = position101, tokenIndex101, depth101
			return false
		},
		/* 11 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action5)> */
		func() bool {
			position131, tokenIndex131, depth131 := position, tokenIndex, depth
			{
				position132 := position
				depth++
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('C') {
						goto l131
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('R') {
						goto l131
					}
					position++
				}
			l135:
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l138
					}
					position++
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if buffer[position] != rune('E') {
						goto l131
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('A') {
						goto l131
					}
					position++
				}
			l139:
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l142
					}
					position++
					goto l141
				l142:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					if buffer[position] != rune('T') {
						goto l131
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
						goto l131
					}
					position++
				}
			l143:
				if !_rules[rulesp]() {
					goto l131
				}
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('S') {
						goto l131
					}
					position++
				}
			l145:
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('T') {
						goto l131
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
						goto l131
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
						goto l131
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
						goto l131
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('M') {
						goto l131
					}
					position++
				}
			l155:
				if !_rules[rulesp]() {
					goto l131
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l131
				}
				if !_rules[rulesp]() {
					goto l131
				}
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('A') {
						goto l131
					}
					position++
				}
			l157:
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
						goto l131
					}
					position++
				}
			l159:
				if !_rules[rulesp]() {
					goto l131
				}
				if !_rules[ruleSelectUnionStmt]() {
					goto l131
				}
				if !_rules[ruleAction5]() {
					goto l131
				}
				depth--
				add(ruleCreateStreamAsSelectUnionStmt, position132)
			}
			return true
		l131:
			position, tokenIndex, depth = position131, tokenIndex131, depth131
			return false
		},
		/* 12 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action6)> */
		func() bool {
			position161, tokenIndex161, depth161 := position, tokenIndex, depth
			{
				position162 := position
				depth++
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l164
					}
					position++
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if buffer[position] != rune('C') {
						goto l161
					}
					position++
				}
			l163:
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('R') {
						goto l161
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
						goto l161
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('A') {
						goto l161
					}
					position++
				}
			l169:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('T') {
						goto l161
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
						goto l161
					}
					position++
				}
			l173:
				if !_rules[rulePausedOpt]() {
					goto l161
				}
				if !_rules[rulesp]() {
					goto l161
				}
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('S') {
						goto l161
					}
					position++
				}
			l175:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('O') {
						goto l161
					}
					position++
				}
			l177:
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
						goto l161
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('R') {
						goto l161
					}
					position++
				}
			l181:
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('C') {
						goto l161
					}
					position++
				}
			l183:
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('E') {
						goto l161
					}
					position++
				}
			l185:
				if !_rules[rulesp]() {
					goto l161
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l161
				}
				if !_rules[rulesp]() {
					goto l161
				}
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
						goto l161
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('Y') {
						goto l161
					}
					position++
				}
			l189:
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l192
					}
					position++
					goto l191
				l192:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
					if buffer[position] != rune('P') {
						goto l161
					}
					position++
				}
			l191:
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('E') {
						goto l161
					}
					position++
				}
			l193:
				if !_rules[rulesp]() {
					goto l161
				}
				if !_rules[ruleSourceSinkType]() {
					goto l161
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l161
				}
				if !_rules[ruleAction6]() {
					goto l161
				}
				depth--
				add(ruleCreateSourceStmt, position162)
			}
			return true
		l161:
			position, tokenIndex, depth = position161, tokenIndex161, depth161
			return false
		},
		/* 13 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action7)> */
		func() bool {
			position195, tokenIndex195, depth195 := position, tokenIndex, depth
			{
				position196 := position
				depth++
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('C') {
						goto l195
					}
					position++
				}
			l197:
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('R') {
						goto l195
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('E') {
						goto l195
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('A') {
						goto l195
					}
					position++
				}
			l203:
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('T') {
						goto l195
					}
					position++
				}
			l205:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('E') {
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
					if buffer[position] != rune('s') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('S') {
						goto l195
					}
					position++
				}
			l209:
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
						goto l195
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
						goto l195
					}
					position++
				}
			l213:
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('K') {
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
					if buffer[position] != rune('t') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('T') {
						goto l195
					}
					position++
				}
			l217:
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('Y') {
						goto l195
					}
					position++
				}
			l219:
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l222
					}
					position++
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if buffer[position] != rune('P') {
						goto l195
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('E') {
						goto l195
					}
					position++
				}
			l223:
				if !_rules[rulesp]() {
					goto l195
				}
				if !_rules[ruleSourceSinkType]() {
					goto l195
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l195
				}
				if !_rules[ruleAction7]() {
					goto l195
				}
				depth--
				add(ruleCreateSinkStmt, position196)
			}
			return true
		l195:
			position, tokenIndex, depth = position195, tokenIndex195, depth195
			return false
		},
		/* 14 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action8)> */
		func() bool {
			position225, tokenIndex225, depth225 := position, tokenIndex, depth
			{
				position226 := position
				depth++
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l228
					}
					position++
					goto l227
				l228:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
					if buffer[position] != rune('C') {
						goto l225
					}
					position++
				}
			l227:
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('R') {
						goto l225
					}
					position++
				}
			l229:
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if buffer[position] != rune('E') {
						goto l225
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
						goto l225
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
						goto l225
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
						goto l225
					}
					position++
				}
			l237:
				if !_rules[rulesp]() {
					goto l225
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
						goto l225
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('T') {
						goto l225
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('A') {
						goto l225
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
				if !_rules[rulesp]() {
					goto l225
				}
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('T') {
						goto l225
					}
					position++
				}
			l249:
				{
					position251, tokenIndex251, depth251 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l252
					}
					position++
					goto l251
				l252:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('Y') {
						goto l225
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('P') {
						goto l225
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
						goto l225
					}
					position++
				}
			l255:
				if !_rules[rulesp]() {
					goto l225
				}
				if !_rules[ruleSourceSinkType]() {
					goto l225
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l225
				}
				if !_rules[ruleAction8]() {
					goto l225
				}
				depth--
				add(ruleCreateStateStmt, position226)
			}
			return true
		l225:
			position, tokenIndex, depth = position225, tokenIndex225, depth225
			return false
		},
		/* 15 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position257, tokenIndex257, depth257 := position, tokenIndex, depth
			{
				position258 := position
				depth++
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
						goto l257
					}
					position++
				}
			l259:
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('P') {
						goto l257
					}
					position++
				}
			l261:
				{
					position263, tokenIndex263, depth263 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if buffer[position] != rune('D') {
						goto l257
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
						goto l257
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
						goto l257
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
						goto l257
					}
					position++
				}
			l269:
				if !_rules[rulesp]() {
					goto l257
				}
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('S') {
						goto l257
					}
					position++
				}
			l271:
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('T') {
						goto l257
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('A') {
						goto l257
					}
					position++
				}
			l275:
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('T') {
						goto l257
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
						goto l257
					}
					position++
				}
			l279:
				if !_rules[rulesp]() {
					goto l257
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l257
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l257
				}
				if !_rules[ruleAction9]() {
					goto l257
				}
				depth--
				add(ruleUpdateStateStmt, position258)
			}
			return true
		l257:
			position, tokenIndex, depth = position257, tokenIndex257, depth257
			return false
		},
		/* 16 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action10)> */
		func() bool {
			position281, tokenIndex281, depth281 := position, tokenIndex, depth
			{
				position282 := position
				depth++
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('U') {
						goto l281
					}
					position++
				}
			l283:
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('P') {
						goto l281
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
						goto l281
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('A') {
						goto l281
					}
					position++
				}
			l289:
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('T') {
						goto l281
					}
					position++
				}
			l291:
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					if buffer[position] != rune('E') {
						goto l281
					}
					position++
				}
			l293:
				if !_rules[rulesp]() {
					goto l281
				}
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('S') {
						goto l281
					}
					position++
				}
			l295:
				{
					position297, tokenIndex297, depth297 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l298
					}
					position++
					goto l297
				l298:
					position, tokenIndex, depth = position297, tokenIndex297, depth297
					if buffer[position] != rune('O') {
						goto l281
					}
					position++
				}
			l297:
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
						goto l281
					}
					position++
				}
			l299:
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l302
					}
					position++
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if buffer[position] != rune('R') {
						goto l281
					}
					position++
				}
			l301:
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('C') {
						goto l281
					}
					position++
				}
			l303:
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l306
					}
					position++
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if buffer[position] != rune('E') {
						goto l281
					}
					position++
				}
			l305:
				if !_rules[rulesp]() {
					goto l281
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l281
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l281
				}
				if !_rules[ruleAction10]() {
					goto l281
				}
				depth--
				add(ruleUpdateSourceStmt, position282)
			}
			return true
		l281:
			position, tokenIndex, depth = position281, tokenIndex281, depth281
			return false
		},
		/* 17 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier UpdateSourceSinkSpecs Action11)> */
		func() bool {
			position307, tokenIndex307, depth307 := position, tokenIndex, depth
			{
				position308 := position
				depth++
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('U') {
						goto l307
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
						goto l307
					}
					position++
				}
			l311:
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('D') {
						goto l307
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('A') {
						goto l307
					}
					position++
				}
			l315:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('T') {
						goto l307
					}
					position++
				}
			l317:
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l320
					}
					position++
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if buffer[position] != rune('E') {
						goto l307
					}
					position++
				}
			l319:
				if !_rules[rulesp]() {
					goto l307
				}
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l322
					}
					position++
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if buffer[position] != rune('S') {
						goto l307
					}
					position++
				}
			l321:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l324
					}
					position++
					goto l323
				l324:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					if buffer[position] != rune('I') {
						goto l307
					}
					position++
				}
			l323:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l326
					}
					position++
					goto l325
				l326:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
					if buffer[position] != rune('N') {
						goto l307
					}
					position++
				}
			l325:
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l328
					}
					position++
					goto l327
				l328:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if buffer[position] != rune('K') {
						goto l307
					}
					position++
				}
			l327:
				if !_rules[rulesp]() {
					goto l307
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l307
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l307
				}
				if !_rules[ruleAction11]() {
					goto l307
				}
				depth--
				add(ruleUpdateSinkStmt, position308)
			}
			return true
		l307:
			position, tokenIndex, depth = position307, tokenIndex307, depth307
			return false
		},
		/* 18 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action12)> */
		func() bool {
			position329, tokenIndex329, depth329 := position, tokenIndex, depth
			{
				position330 := position
				depth++
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('I') {
						goto l329
					}
					position++
				}
			l331:
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('N') {
						goto l329
					}
					position++
				}
			l333:
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
						goto l329
					}
					position++
				}
			l335:
				{
					position337, tokenIndex337, depth337 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l338
					}
					position++
					goto l337
				l338:
					position, tokenIndex, depth = position337, tokenIndex337, depth337
					if buffer[position] != rune('E') {
						goto l329
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
						goto l329
					}
					position++
				}
			l339:
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l342
					}
					position++
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					if buffer[position] != rune('T') {
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
					if buffer[position] != rune('i') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('I') {
						goto l329
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
						goto l329
					}
					position++
				}
			l345:
				{
					position347, tokenIndex347, depth347 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l348
					}
					position++
					goto l347
				l348:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if buffer[position] != rune('T') {
						goto l329
					}
					position++
				}
			l347:
				{
					position349, tokenIndex349, depth349 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l350
					}
					position++
					goto l349
				l350:
					position, tokenIndex, depth = position349, tokenIndex349, depth349
					if buffer[position] != rune('O') {
						goto l329
					}
					position++
				}
			l349:
				if !_rules[rulesp]() {
					goto l329
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l329
				}
				if !_rules[rulesp]() {
					goto l329
				}
				if !_rules[ruleSelectStmt]() {
					goto l329
				}
				if !_rules[ruleAction12]() {
					goto l329
				}
				depth--
				add(ruleInsertIntoSelectStmt, position330)
			}
			return true
		l329:
			position, tokenIndex, depth = position329, tokenIndex329, depth329
			return false
		},
		/* 19 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action13)> */
		func() bool {
			position351, tokenIndex351, depth351 := position, tokenIndex, depth
			{
				position352 := position
				depth++
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l354
					}
					position++
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if buffer[position] != rune('I') {
						goto l351
					}
					position++
				}
			l353:
				{
					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l356
					}
					position++
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune('N') {
						goto l351
					}
					position++
				}
			l355:
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
						goto l351
					}
					position++
				}
			l357:
				{
					position359, tokenIndex359, depth359 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l360
					}
					position++
					goto l359
				l360:
					position, tokenIndex, depth = position359, tokenIndex359, depth359
					if buffer[position] != rune('E') {
						goto l351
					}
					position++
				}
			l359:
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l362
					}
					position++
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if buffer[position] != rune('R') {
						goto l351
					}
					position++
				}
			l361:
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l364
					}
					position++
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					if buffer[position] != rune('T') {
						goto l351
					}
					position++
				}
			l363:
				if !_rules[rulesp]() {
					goto l351
				}
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('I') {
						goto l351
					}
					position++
				}
			l365:
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					if buffer[position] != rune('N') {
						goto l351
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('T') {
						goto l351
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
						goto l351
					}
					position++
				}
			l371:
				if !_rules[rulesp]() {
					goto l351
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l351
				}
				if !_rules[rulesp]() {
					goto l351
				}
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('F') {
						goto l351
					}
					position++
				}
			l373:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l376
					}
					position++
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					if buffer[position] != rune('R') {
						goto l351
					}
					position++
				}
			l375:
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l378
					}
					position++
					goto l377
				l378:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
					if buffer[position] != rune('O') {
						goto l351
					}
					position++
				}
			l377:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('M') {
						goto l351
					}
					position++
				}
			l379:
				if !_rules[rulesp]() {
					goto l351
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l351
				}
				if !_rules[ruleAction13]() {
					goto l351
				}
				depth--
				add(ruleInsertIntoFromStmt, position352)
			}
			return true
		l351:
			position, tokenIndex, depth = position351, tokenIndex351, depth351
			return false
		},
		/* 20 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position381, tokenIndex381, depth381 := position, tokenIndex, depth
			{
				position382 := position
				depth++
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('P') {
						goto l381
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('A') {
						goto l381
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
						goto l381
					}
					position++
				}
			l387:
				{
					position389, tokenIndex389, depth389 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l390
					}
					position++
					goto l389
				l390:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
					if buffer[position] != rune('S') {
						goto l381
					}
					position++
				}
			l389:
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
					if buffer[position] != rune('E') {
						goto l381
					}
					position++
				}
			l391:
				if !_rules[rulesp]() {
					goto l381
				}
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l394
					}
					position++
					goto l393
				l394:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
					if buffer[position] != rune('S') {
						goto l381
					}
					position++
				}
			l393:
				{
					position395, tokenIndex395, depth395 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l396
					}
					position++
					goto l395
				l396:
					position, tokenIndex, depth = position395, tokenIndex395, depth395
					if buffer[position] != rune('O') {
						goto l381
					}
					position++
				}
			l395:
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l398
					}
					position++
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if buffer[position] != rune('U') {
						goto l381
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
						goto l381
					}
					position++
				}
			l399:
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l402
					}
					position++
					goto l401
				l402:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if buffer[position] != rune('C') {
						goto l381
					}
					position++
				}
			l401:
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l404
					}
					position++
					goto l403
				l404:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
					if buffer[position] != rune('E') {
						goto l381
					}
					position++
				}
			l403:
				if !_rules[rulesp]() {
					goto l381
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l381
				}
				if !_rules[ruleAction14]() {
					goto l381
				}
				depth--
				add(rulePauseSourceStmt, position382)
			}
			return true
		l381:
			position, tokenIndex, depth = position381, tokenIndex381, depth381
			return false
		},
		/* 21 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position405, tokenIndex405, depth405 := position, tokenIndex, depth
			{
				position406 := position
				depth++
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
						goto l405
					}
					position++
				}
			l407:
				{
					position409, tokenIndex409, depth409 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l410
					}
					position++
					goto l409
				l410:
					position, tokenIndex, depth = position409, tokenIndex409, depth409
					if buffer[position] != rune('E') {
						goto l405
					}
					position++
				}
			l409:
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l412
					}
					position++
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if buffer[position] != rune('S') {
						goto l405
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
						goto l405
					}
					position++
				}
			l413:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l416
					}
					position++
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if buffer[position] != rune('M') {
						goto l405
					}
					position++
				}
			l415:
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l418
					}
					position++
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if buffer[position] != rune('E') {
						goto l405
					}
					position++
				}
			l417:
				if !_rules[rulesp]() {
					goto l405
				}
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l420
					}
					position++
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if buffer[position] != rune('S') {
						goto l405
					}
					position++
				}
			l419:
				{
					position421, tokenIndex421, depth421 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l422
					}
					position++
					goto l421
				l422:
					position, tokenIndex, depth = position421, tokenIndex421, depth421
					if buffer[position] != rune('O') {
						goto l405
					}
					position++
				}
			l421:
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l424
					}
					position++
					goto l423
				l424:
					position, tokenIndex, depth = position423, tokenIndex423, depth423
					if buffer[position] != rune('U') {
						goto l405
					}
					position++
				}
			l423:
				{
					position425, tokenIndex425, depth425 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l426
					}
					position++
					goto l425
				l426:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if buffer[position] != rune('R') {
						goto l405
					}
					position++
				}
			l425:
				{
					position427, tokenIndex427, depth427 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l428
					}
					position++
					goto l427
				l428:
					position, tokenIndex, depth = position427, tokenIndex427, depth427
					if buffer[position] != rune('C') {
						goto l405
					}
					position++
				}
			l427:
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l430
					}
					position++
					goto l429
				l430:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if buffer[position] != rune('E') {
						goto l405
					}
					position++
				}
			l429:
				if !_rules[rulesp]() {
					goto l405
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l405
				}
				if !_rules[ruleAction15]() {
					goto l405
				}
				depth--
				add(ruleResumeSourceStmt, position406)
			}
			return true
		l405:
			position, tokenIndex, depth = position405, tokenIndex405, depth405
			return false
		},
		/* 22 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position431, tokenIndex431, depth431 := position, tokenIndex, depth
			{
				position432 := position
				depth++
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l434
					}
					position++
					goto l433
				l434:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if buffer[position] != rune('R') {
						goto l431
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
						goto l431
					}
					position++
				}
			l435:
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l438
					}
					position++
					goto l437
				l438:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
					if buffer[position] != rune('W') {
						goto l431
					}
					position++
				}
			l437:
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l440
					}
					position++
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if buffer[position] != rune('I') {
						goto l431
					}
					position++
				}
			l439:
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l442
					}
					position++
					goto l441
				l442:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
					if buffer[position] != rune('N') {
						goto l431
					}
					position++
				}
			l441:
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l444
					}
					position++
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if buffer[position] != rune('D') {
						goto l431
					}
					position++
				}
			l443:
				if !_rules[rulesp]() {
					goto l431
				}
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l446
					}
					position++
					goto l445
				l446:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
					if buffer[position] != rune('S') {
						goto l431
					}
					position++
				}
			l445:
				{
					position447, tokenIndex447, depth447 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l448
					}
					position++
					goto l447
				l448:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if buffer[position] != rune('O') {
						goto l431
					}
					position++
				}
			l447:
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l450
					}
					position++
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					if buffer[position] != rune('U') {
						goto l431
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
						goto l431
					}
					position++
				}
			l451:
				{
					position453, tokenIndex453, depth453 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l454
					}
					position++
					goto l453
				l454:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
					if buffer[position] != rune('C') {
						goto l431
					}
					position++
				}
			l453:
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l456
					}
					position++
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if buffer[position] != rune('E') {
						goto l431
					}
					position++
				}
			l455:
				if !_rules[rulesp]() {
					goto l431
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l431
				}
				if !_rules[ruleAction16]() {
					goto l431
				}
				depth--
				add(ruleRewindSourceStmt, position432)
			}
			return true
		l431:
			position, tokenIndex, depth = position431, tokenIndex431, depth431
			return false
		},
		/* 23 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action17)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l460
					}
					position++
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if buffer[position] != rune('D') {
						goto l457
					}
					position++
				}
			l459:
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l462
					}
					position++
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('R') {
						goto l457
					}
					position++
				}
			l461:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('O') {
						goto l457
					}
					position++
				}
			l463:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l466
					}
					position++
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					if buffer[position] != rune('P') {
						goto l457
					}
					position++
				}
			l465:
				if !_rules[rulesp]() {
					goto l457
				}
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l468
					}
					position++
					goto l467
				l468:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
					if buffer[position] != rune('S') {
						goto l457
					}
					position++
				}
			l467:
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l470
					}
					position++
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if buffer[position] != rune('O') {
						goto l457
					}
					position++
				}
			l469:
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l472
					}
					position++
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if buffer[position] != rune('U') {
						goto l457
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
						goto l457
					}
					position++
				}
			l473:
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l476
					}
					position++
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if buffer[position] != rune('C') {
						goto l457
					}
					position++
				}
			l475:
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l478
					}
					position++
					goto l477
				l478:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
					if buffer[position] != rune('E') {
						goto l457
					}
					position++
				}
			l477:
				if !_rules[rulesp]() {
					goto l457
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l457
				}
				if !_rules[ruleAction17]() {
					goto l457
				}
				depth--
				add(ruleDropSourceStmt, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 24 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action18)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if buffer[position] != rune('D') {
						goto l479
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
						goto l479
					}
					position++
				}
			l483:
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l486
					}
					position++
					goto l485
				l486:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					if buffer[position] != rune('O') {
						goto l479
					}
					position++
				}
			l485:
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l488
					}
					position++
					goto l487
				l488:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
					if buffer[position] != rune('P') {
						goto l479
					}
					position++
				}
			l487:
				if !_rules[rulesp]() {
					goto l479
				}
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l490
					}
					position++
					goto l489
				l490:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
					if buffer[position] != rune('S') {
						goto l479
					}
					position++
				}
			l489:
				{
					position491, tokenIndex491, depth491 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l492
					}
					position++
					goto l491
				l492:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if buffer[position] != rune('T') {
						goto l479
					}
					position++
				}
			l491:
				{
					position493, tokenIndex493, depth493 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l494
					}
					position++
					goto l493
				l494:
					position, tokenIndex, depth = position493, tokenIndex493, depth493
					if buffer[position] != rune('R') {
						goto l479
					}
					position++
				}
			l493:
				{
					position495, tokenIndex495, depth495 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l496
					}
					position++
					goto l495
				l496:
					position, tokenIndex, depth = position495, tokenIndex495, depth495
					if buffer[position] != rune('E') {
						goto l479
					}
					position++
				}
			l495:
				{
					position497, tokenIndex497, depth497 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l498
					}
					position++
					goto l497
				l498:
					position, tokenIndex, depth = position497, tokenIndex497, depth497
					if buffer[position] != rune('A') {
						goto l479
					}
					position++
				}
			l497:
				{
					position499, tokenIndex499, depth499 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l500
					}
					position++
					goto l499
				l500:
					position, tokenIndex, depth = position499, tokenIndex499, depth499
					if buffer[position] != rune('M') {
						goto l479
					}
					position++
				}
			l499:
				if !_rules[rulesp]() {
					goto l479
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l479
				}
				if !_rules[ruleAction18]() {
					goto l479
				}
				depth--
				add(ruleDropStreamStmt, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 25 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action19)> */
		func() bool {
			position501, tokenIndex501, depth501 := position, tokenIndex, depth
			{
				position502 := position
				depth++
				{
					position503, tokenIndex503, depth503 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l504
					}
					position++
					goto l503
				l504:
					position, tokenIndex, depth = position503, tokenIndex503, depth503
					if buffer[position] != rune('D') {
						goto l501
					}
					position++
				}
			l503:
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l506
					}
					position++
					goto l505
				l506:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
					if buffer[position] != rune('R') {
						goto l501
					}
					position++
				}
			l505:
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('O') {
						goto l501
					}
					position++
				}
			l507:
				{
					position509, tokenIndex509, depth509 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l510
					}
					position++
					goto l509
				l510:
					position, tokenIndex, depth = position509, tokenIndex509, depth509
					if buffer[position] != rune('P') {
						goto l501
					}
					position++
				}
			l509:
				if !_rules[rulesp]() {
					goto l501
				}
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if buffer[position] != rune('S') {
						goto l501
					}
					position++
				}
			l511:
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l514
					}
					position++
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if buffer[position] != rune('I') {
						goto l501
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
						goto l501
					}
					position++
				}
			l515:
				{
					position517, tokenIndex517, depth517 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l518
					}
					position++
					goto l517
				l518:
					position, tokenIndex, depth = position517, tokenIndex517, depth517
					if buffer[position] != rune('K') {
						goto l501
					}
					position++
				}
			l517:
				if !_rules[rulesp]() {
					goto l501
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l501
				}
				if !_rules[ruleAction19]() {
					goto l501
				}
				depth--
				add(ruleDropSinkStmt, position502)
			}
			return true
		l501:
			position, tokenIndex, depth = position501, tokenIndex501, depth501
			return false
		},
		/* 26 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action20)> */
		func() bool {
			position519, tokenIndex519, depth519 := position, tokenIndex, depth
			{
				position520 := position
				depth++
				{
					position521, tokenIndex521, depth521 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l522
					}
					position++
					goto l521
				l522:
					position, tokenIndex, depth = position521, tokenIndex521, depth521
					if buffer[position] != rune('D') {
						goto l519
					}
					position++
				}
			l521:
				{
					position523, tokenIndex523, depth523 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l524
					}
					position++
					goto l523
				l524:
					position, tokenIndex, depth = position523, tokenIndex523, depth523
					if buffer[position] != rune('R') {
						goto l519
					}
					position++
				}
			l523:
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l526
					}
					position++
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if buffer[position] != rune('O') {
						goto l519
					}
					position++
				}
			l525:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l528
					}
					position++
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if buffer[position] != rune('P') {
						goto l519
					}
					position++
				}
			l527:
				if !_rules[rulesp]() {
					goto l519
				}
				{
					position529, tokenIndex529, depth529 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l530
					}
					position++
					goto l529
				l530:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if buffer[position] != rune('S') {
						goto l519
					}
					position++
				}
			l529:
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l532
					}
					position++
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if buffer[position] != rune('T') {
						goto l519
					}
					position++
				}
			l531:
				{
					position533, tokenIndex533, depth533 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l534
					}
					position++
					goto l533
				l534:
					position, tokenIndex, depth = position533, tokenIndex533, depth533
					if buffer[position] != rune('A') {
						goto l519
					}
					position++
				}
			l533:
				{
					position535, tokenIndex535, depth535 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l536
					}
					position++
					goto l535
				l536:
					position, tokenIndex, depth = position535, tokenIndex535, depth535
					if buffer[position] != rune('T') {
						goto l519
					}
					position++
				}
			l535:
				{
					position537, tokenIndex537, depth537 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l538
					}
					position++
					goto l537
				l538:
					position, tokenIndex, depth = position537, tokenIndex537, depth537
					if buffer[position] != rune('E') {
						goto l519
					}
					position++
				}
			l537:
				if !_rules[rulesp]() {
					goto l519
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l519
				}
				if !_rules[ruleAction20]() {
					goto l519
				}
				depth--
				add(ruleDropStateStmt, position520)
			}
			return true
		l519:
			position, tokenIndex, depth = position519, tokenIndex519, depth519
			return false
		},
		/* 27 LoadStateStmt <- <(('l' / 'L') ('o' / 'O') ('a' / 'A') ('d' / 'D') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType StateTagOpt SetOptSpecs Action21)> */
		func() bool {
			position539, tokenIndex539, depth539 := position, tokenIndex, depth
			{
				position540 := position
				depth++
				{
					position541, tokenIndex541, depth541 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l542
					}
					position++
					goto l541
				l542:
					position, tokenIndex, depth = position541, tokenIndex541, depth541
					if buffer[position] != rune('L') {
						goto l539
					}
					position++
				}
			l541:
				{
					position543, tokenIndex543, depth543 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l544
					}
					position++
					goto l543
				l544:
					position, tokenIndex, depth = position543, tokenIndex543, depth543
					if buffer[position] != rune('O') {
						goto l539
					}
					position++
				}
			l543:
				{
					position545, tokenIndex545, depth545 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l546
					}
					position++
					goto l545
				l546:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if buffer[position] != rune('A') {
						goto l539
					}
					position++
				}
			l545:
				{
					position547, tokenIndex547, depth547 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l548
					}
					position++
					goto l547
				l548:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
					if buffer[position] != rune('D') {
						goto l539
					}
					position++
				}
			l547:
				if !_rules[rulesp]() {
					goto l539
				}
				{
					position549, tokenIndex549, depth549 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l550
					}
					position++
					goto l549
				l550:
					position, tokenIndex, depth = position549, tokenIndex549, depth549
					if buffer[position] != rune('S') {
						goto l539
					}
					position++
				}
			l549:
				{
					position551, tokenIndex551, depth551 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l552
					}
					position++
					goto l551
				l552:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if buffer[position] != rune('T') {
						goto l539
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
						goto l539
					}
					position++
				}
			l553:
				{
					position555, tokenIndex555, depth555 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l556
					}
					position++
					goto l555
				l556:
					position, tokenIndex, depth = position555, tokenIndex555, depth555
					if buffer[position] != rune('T') {
						goto l539
					}
					position++
				}
			l555:
				{
					position557, tokenIndex557, depth557 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l558
					}
					position++
					goto l557
				l558:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if buffer[position] != rune('E') {
						goto l539
					}
					position++
				}
			l557:
				if !_rules[rulesp]() {
					goto l539
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l539
				}
				if !_rules[rulesp]() {
					goto l539
				}
				{
					position559, tokenIndex559, depth559 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l560
					}
					position++
					goto l559
				l560:
					position, tokenIndex, depth = position559, tokenIndex559, depth559
					if buffer[position] != rune('T') {
						goto l539
					}
					position++
				}
			l559:
				{
					position561, tokenIndex561, depth561 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l562
					}
					position++
					goto l561
				l562:
					position, tokenIndex, depth = position561, tokenIndex561, depth561
					if buffer[position] != rune('Y') {
						goto l539
					}
					position++
				}
			l561:
				{
					position563, tokenIndex563, depth563 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l564
					}
					position++
					goto l563
				l564:
					position, tokenIndex, depth = position563, tokenIndex563, depth563
					if buffer[position] != rune('P') {
						goto l539
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
						goto l539
					}
					position++
				}
			l565:
				if !_rules[rulesp]() {
					goto l539
				}
				if !_rules[ruleSourceSinkType]() {
					goto l539
				}
				if !_rules[ruleStateTagOpt]() {
					goto l539
				}
				if !_rules[ruleSetOptSpecs]() {
					goto l539
				}
				if !_rules[ruleAction21]() {
					goto l539
				}
				depth--
				add(ruleLoadStateStmt, position540)
			}
			return true
		l539:
			position, tokenIndex, depth = position539, tokenIndex539, depth539
			return false
		},
		/* 28 LoadStateOrCreateStmt <- <(LoadStateStmt sp (('o' / 'O') ('r' / 'R')) sp (('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp (('i' / 'I') ('f' / 'F')) sp (('n' / 'N') ('o' / 'O') ('t' / 'T')) sp ((('s' / 'S') ('a' / 'A') ('v' / 'V') ('e' / 'E') ('d' / 'D')) / (('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S'))) SourceSinkSpecs Action22)> */
		func() bool {
			position567, tokenIndex567, depth567 := position, tokenIndex, depth
			{
				position568 := position
				depth++
				if !_rules[ruleLoadStateStmt]() {
					goto l567
				}
				if !_rules[rulesp]() {
					goto l567
				}
				{
					position569, tokenIndex569, depth569 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l570
					}
					position++
					goto l569
				l570:
					position, tokenIndex, depth = position569, tokenIndex569, depth569
					if buffer[position] != rune('O') {
						goto l567
					}
					position++
				}
			l569:
				{
					position571, tokenIndex571, depth571 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l572
					}
					position++
					goto l571
				l572:
					position, tokenIndex, depth = position571, tokenIndex571, depth571
					if buffer[position] != rune('R') {
						goto l567
					}
					position++
				}
			l571:
				if !_rules[rulesp]() {
					goto l567
				}
				{
					position573, tokenIndex573, depth573 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l574
					}
					position++
					goto l573
				l574:
					position, tokenIndex, depth = position573, tokenIndex573, depth573
					if buffer[position] != rune('C') {
						goto l567
					}
					position++
				}
			l573:
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
						goto l567
					}
					position++
				}
			l575:
				{
					position577, tokenIndex577, depth577 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l578
					}
					position++
					goto l577
				l578:
					position, tokenIndex, depth = position577, tokenIndex577, depth577
					if buffer[position] != rune('E') {
						goto l567
					}
					position++
				}
			l577:
				{
					position579, tokenIndex579, depth579 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l580
					}
					position++
					goto l579
				l580:
					position, tokenIndex, depth = position579, tokenIndex579, depth579
					if buffer[position] != rune('A') {
						goto l567
					}
					position++
				}
			l579:
				{
					position581, tokenIndex581, depth581 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l582
					}
					position++
					goto l581
				l582:
					position, tokenIndex, depth = position581, tokenIndex581, depth581
					if buffer[position] != rune('T') {
						goto l567
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
						goto l567
					}
					position++
				}
			l583:
				if !_rules[rulesp]() {
					goto l567
				}
				{
					position585, tokenIndex585, depth585 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l586
					}
					position++
					goto l585
				l586:
					position, tokenIndex, depth = position585, tokenIndex585, depth585
					if buffer[position] != rune('I') {
						goto l567
					}
					position++
				}
			l585:
				{
					position587, tokenIndex587, depth587 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l588
					}
					position++
					goto l587
				l588:
					position, tokenIndex, depth = position587, tokenIndex587, depth587
					if buffer[position] != rune('F') {
						goto l567
					}
					position++
				}
			l587:
				if !_rules[rulesp]() {
					goto l567
				}
				{
					position589, tokenIndex589, depth589 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l590
					}
					position++
					goto l589
				l590:
					position, tokenIndex, depth = position589, tokenIndex589, depth589
					if buffer[position] != rune('N') {
						goto l567
					}
					position++
				}
			l589:
				{
					position591, tokenIndex591, depth591 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l592
					}
					position++
					goto l591
				l592:
					position, tokenIndex, depth = position591, tokenIndex591, depth591
					if buffer[position] != rune('O') {
						goto l567
					}
					position++
				}
			l591:
				{
					position593, tokenIndex593, depth593 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l594
					}
					position++
					goto l593
				l594:
					position, tokenIndex, depth = position593, tokenIndex593, depth593
					if buffer[position] != rune('T') {
						goto l567
					}
					position++
				}
			l593:
				if !_rules[rulesp]() {
					goto l567
				}
				{
					position595, tokenIndex595, depth595 := position, tokenIndex, depth
					{
						position597, tokenIndex597, depth597 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l598
						}
						position++
						goto l597
					l598:
						position, tokenIndex, depth = position597, tokenIndex597, depth597
						if buffer[position] != rune('S') {
							goto l596
						}
						position++
					}
				l597:
					{
						position599, tokenIndex599, depth599 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l600
						}
						position++
						goto l599
					l600:
						position, tokenIndex, depth = position599, tokenIndex599, depth599
						if buffer[position] != rune('A') {
							goto l596
						}
						position++
					}
				l599:
					{
						position601, tokenIndex601, depth601 := position, tokenIndex, depth
						if buffer[position] != rune('v') {
							goto l602
						}
						position++
						goto l601
					l602:
						position, tokenIndex, depth = position601, tokenIndex601, depth601
						if buffer[position] != rune('V') {
							goto l596
						}
						position++
					}
				l601:
					{
						position603, tokenIndex603, depth603 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l604
						}
						position++
						goto l603
					l604:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
						if buffer[position] != rune('E') {
							goto l596
						}
						position++
					}
				l603:
					{
						position605, tokenIndex605, depth605 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l606
						}
						position++
						goto l605
					l606:
						position, tokenIndex, depth = position605, tokenIndex605, depth605
						if buffer[position] != rune('D') {
							goto l596
						}
						position++
					}
				l605:
					goto l595
				l596:
					position, tokenIndex, depth = position595, tokenIndex595, depth595
					{
						position607, tokenIndex607, depth607 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l608
						}
						position++
						goto l607
					l608:
						position, tokenIndex, depth = position607, tokenIndex607, depth607
						if buffer[position] != rune('E') {
							goto l567
						}
						position++
					}
				l607:
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('x') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('X') {
							goto l567
						}
						position++
					}
				l609:
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('I') {
							goto l567
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('S') {
							goto l567
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('T') {
							goto l567
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('S') {
							goto l567
						}
						position++
					}
				l617:
				}
			l595:
				if !_rules[ruleSourceSinkSpecs]() {
					goto l567
				}
				if !_rules[ruleAction22]() {
					goto l567
				}
				depth--
				add(ruleLoadStateOrCreateStmt, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 29 SaveStateStmt <- <(('s' / 'S') ('a' / 'A') ('v' / 'V') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier StateTagOpt Action23)> */
		func() bool {
			position619, tokenIndex619, depth619 := position, tokenIndex, depth
			{
				position620 := position
				depth++
				{
					position621, tokenIndex621, depth621 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l622
					}
					position++
					goto l621
				l622:
					position, tokenIndex, depth = position621, tokenIndex621, depth621
					if buffer[position] != rune('S') {
						goto l619
					}
					position++
				}
			l621:
				{
					position623, tokenIndex623, depth623 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l624
					}
					position++
					goto l623
				l624:
					position, tokenIndex, depth = position623, tokenIndex623, depth623
					if buffer[position] != rune('A') {
						goto l619
					}
					position++
				}
			l623:
				{
					position625, tokenIndex625, depth625 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l626
					}
					position++
					goto l625
				l626:
					position, tokenIndex, depth = position625, tokenIndex625, depth625
					if buffer[position] != rune('V') {
						goto l619
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
						goto l619
					}
					position++
				}
			l627:
				if !_rules[rulesp]() {
					goto l619
				}
				{
					position629, tokenIndex629, depth629 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l630
					}
					position++
					goto l629
				l630:
					position, tokenIndex, depth = position629, tokenIndex629, depth629
					if buffer[position] != rune('S') {
						goto l619
					}
					position++
				}
			l629:
				{
					position631, tokenIndex631, depth631 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l632
					}
					position++
					goto l631
				l632:
					position, tokenIndex, depth = position631, tokenIndex631, depth631
					if buffer[position] != rune('T') {
						goto l619
					}
					position++
				}
			l631:
				{
					position633, tokenIndex633, depth633 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l634
					}
					position++
					goto l633
				l634:
					position, tokenIndex, depth = position633, tokenIndex633, depth633
					if buffer[position] != rune('A') {
						goto l619
					}
					position++
				}
			l633:
				{
					position635, tokenIndex635, depth635 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l636
					}
					position++
					goto l635
				l636:
					position, tokenIndex, depth = position635, tokenIndex635, depth635
					if buffer[position] != rune('T') {
						goto l619
					}
					position++
				}
			l635:
				{
					position637, tokenIndex637, depth637 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l638
					}
					position++
					goto l637
				l638:
					position, tokenIndex, depth = position637, tokenIndex637, depth637
					if buffer[position] != rune('E') {
						goto l619
					}
					position++
				}
			l637:
				if !_rules[rulesp]() {
					goto l619
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l619
				}
				if !_rules[ruleStateTagOpt]() {
					goto l619
				}
				if !_rules[ruleAction23]() {
					goto l619
				}
				depth--
				add(ruleSaveStateStmt, position620)
			}
			return true
		l619:
			position, tokenIndex, depth = position619, tokenIndex619, depth619
			return false
		},
		/* 30 EvalStmt <- <(('e' / 'E') ('v' / 'V') ('a' / 'A') ('l' / 'L') sp Expression <(sp (('o' / 'O') ('n' / 'N')) sp MapExpr)?> Action24)> */
		func() bool {
			position639, tokenIndex639, depth639 := position, tokenIndex, depth
			{
				position640 := position
				depth++
				{
					position641, tokenIndex641, depth641 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l642
					}
					position++
					goto l641
				l642:
					position, tokenIndex, depth = position641, tokenIndex641, depth641
					if buffer[position] != rune('E') {
						goto l639
					}
					position++
				}
			l641:
				{
					position643, tokenIndex643, depth643 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l644
					}
					position++
					goto l643
				l644:
					position, tokenIndex, depth = position643, tokenIndex643, depth643
					if buffer[position] != rune('V') {
						goto l639
					}
					position++
				}
			l643:
				{
					position645, tokenIndex645, depth645 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l646
					}
					position++
					goto l645
				l646:
					position, tokenIndex, depth = position645, tokenIndex645, depth645
					if buffer[position] != rune('A') {
						goto l639
					}
					position++
				}
			l645:
				{
					position647, tokenIndex647, depth647 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l648
					}
					position++
					goto l647
				l648:
					position, tokenIndex, depth = position647, tokenIndex647, depth647
					if buffer[position] != rune('L') {
						goto l639
					}
					position++
				}
			l647:
				if !_rules[rulesp]() {
					goto l639
				}
				if !_rules[ruleExpression]() {
					goto l639
				}
				{
					position649 := position
					depth++
					{
						position650, tokenIndex650, depth650 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l650
						}
						{
							position652, tokenIndex652, depth652 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l653
							}
							position++
							goto l652
						l653:
							position, tokenIndex, depth = position652, tokenIndex652, depth652
							if buffer[position] != rune('O') {
								goto l650
							}
							position++
						}
					l652:
						{
							position654, tokenIndex654, depth654 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l655
							}
							position++
							goto l654
						l655:
							position, tokenIndex, depth = position654, tokenIndex654, depth654
							if buffer[position] != rune('N') {
								goto l650
							}
							position++
						}
					l654:
						if !_rules[rulesp]() {
							goto l650
						}
						if !_rules[ruleMapExpr]() {
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
				if !_rules[ruleAction24]() {
					goto l639
				}
				depth--
				add(ruleEvalStmt, position640)
			}
			return true
		l639:
			position, tokenIndex, depth = position639, tokenIndex639, depth639
			return false
		},
		/* 31 Emitter <- <(sp (ISTREAM / DSTREAM / RSTREAM) EmitterOptions Action25)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				if !_rules[rulesp]() {
					goto l656
				}
				{
					position658, tokenIndex658, depth658 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l659
					}
					goto l658
				l659:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
					if !_rules[ruleDSTREAM]() {
						goto l660
					}
					goto l658
				l660:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
					if !_rules[ruleRSTREAM]() {
						goto l656
					}
				}
			l658:
				if !_rules[ruleEmitterOptions]() {
					goto l656
				}
				if !_rules[ruleAction25]() {
					goto l656
				}
				depth--
				add(ruleEmitter, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 32 EmitterOptions <- <(<(spOpt '[' spOpt EmitterOptionCombinations spOpt ']')?> Action26)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				{
					position663 := position
					depth++
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l664
						}
						if buffer[position] != rune('[') {
							goto l664
						}
						position++
						if !_rules[rulespOpt]() {
							goto l664
						}
						if !_rules[ruleEmitterOptionCombinations]() {
							goto l664
						}
						if !_rules[rulespOpt]() {
							goto l664
						}
						if buffer[position] != rune(']') {
							goto l664
						}
						position++
						goto l665
					l664:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
					}
				l665:
					depth--
					add(rulePegText, position663)
				}
				if !_rules[ruleAction26]() {
					goto l661
				}
				depth--
				add(ruleEmitterOptions, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 33 EmitterOptionCombinations <- <(EmitterLimit / (EmitterSample sp EmitterLimit) / EmitterSample)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				{
					position668, tokenIndex668, depth668 := position, tokenIndex, depth
					if !_rules[ruleEmitterLimit]() {
						goto l669
					}
					goto l668
				l669:
					position, tokenIndex, depth = position668, tokenIndex668, depth668
					if !_rules[ruleEmitterSample]() {
						goto l670
					}
					if !_rules[rulesp]() {
						goto l670
					}
					if !_rules[ruleEmitterLimit]() {
						goto l670
					}
					goto l668
				l670:
					position, tokenIndex, depth = position668, tokenIndex668, depth668
					if !_rules[ruleEmitterSample]() {
						goto l666
					}
				}
			l668:
				depth--
				add(ruleEmitterOptionCombinations, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 34 EmitterLimit <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') sp NumericLiteral Action27)> */
		func() bool {
			position671, tokenIndex671, depth671 := position, tokenIndex, depth
			{
				position672 := position
				depth++
				{
					position673, tokenIndex673, depth673 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l674
					}
					position++
					goto l673
				l674:
					position, tokenIndex, depth = position673, tokenIndex673, depth673
					if buffer[position] != rune('L') {
						goto l671
					}
					position++
				}
			l673:
				{
					position675, tokenIndex675, depth675 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l676
					}
					position++
					goto l675
				l676:
					position, tokenIndex, depth = position675, tokenIndex675, depth675
					if buffer[position] != rune('I') {
						goto l671
					}
					position++
				}
			l675:
				{
					position677, tokenIndex677, depth677 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l678
					}
					position++
					goto l677
				l678:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if buffer[position] != rune('M') {
						goto l671
					}
					position++
				}
			l677:
				{
					position679, tokenIndex679, depth679 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l680
					}
					position++
					goto l679
				l680:
					position, tokenIndex, depth = position679, tokenIndex679, depth679
					if buffer[position] != rune('I') {
						goto l671
					}
					position++
				}
			l679:
				{
					position681, tokenIndex681, depth681 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l682
					}
					position++
					goto l681
				l682:
					position, tokenIndex, depth = position681, tokenIndex681, depth681
					if buffer[position] != rune('T') {
						goto l671
					}
					position++
				}
			l681:
				if !_rules[rulesp]() {
					goto l671
				}
				if !_rules[ruleNumericLiteral]() {
					goto l671
				}
				if !_rules[ruleAction27]() {
					goto l671
				}
				depth--
				add(ruleEmitterLimit, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 35 EmitterSample <- <(CountBasedSampling / RandomizedSampling / TimeBasedSampling)> */
		func() bool {
			position683, tokenIndex683, depth683 := position, tokenIndex, depth
			{
				position684 := position
				depth++
				{
					position685, tokenIndex685, depth685 := position, tokenIndex, depth
					if !_rules[ruleCountBasedSampling]() {
						goto l686
					}
					goto l685
				l686:
					position, tokenIndex, depth = position685, tokenIndex685, depth685
					if !_rules[ruleRandomizedSampling]() {
						goto l687
					}
					goto l685
				l687:
					position, tokenIndex, depth = position685, tokenIndex685, depth685
					if !_rules[ruleTimeBasedSampling]() {
						goto l683
					}
				}
			l685:
				depth--
				add(ruleEmitterSample, position684)
			}
			return true
		l683:
			position, tokenIndex, depth = position683, tokenIndex683, depth683
			return false
		},
		/* 36 CountBasedSampling <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp NumericLiteral spOpt '-'? spOpt ((('s' / 'S') ('t' / 'T')) / (('n' / 'N') ('d' / 'D')) / (('r' / 'R') ('d' / 'D')) / (('t' / 'T') ('h' / 'H'))) sp (('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E')) Action28)> */
		func() bool {
			position688, tokenIndex688, depth688 := position, tokenIndex, depth
			{
				position689 := position
				depth++
				{
					position690, tokenIndex690, depth690 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l691
					}
					position++
					goto l690
				l691:
					position, tokenIndex, depth = position690, tokenIndex690, depth690
					if buffer[position] != rune('E') {
						goto l688
					}
					position++
				}
			l690:
				{
					position692, tokenIndex692, depth692 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l693
					}
					position++
					goto l692
				l693:
					position, tokenIndex, depth = position692, tokenIndex692, depth692
					if buffer[position] != rune('V') {
						goto l688
					}
					position++
				}
			l692:
				{
					position694, tokenIndex694, depth694 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l695
					}
					position++
					goto l694
				l695:
					position, tokenIndex, depth = position694, tokenIndex694, depth694
					if buffer[position] != rune('E') {
						goto l688
					}
					position++
				}
			l694:
				{
					position696, tokenIndex696, depth696 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l697
					}
					position++
					goto l696
				l697:
					position, tokenIndex, depth = position696, tokenIndex696, depth696
					if buffer[position] != rune('R') {
						goto l688
					}
					position++
				}
			l696:
				{
					position698, tokenIndex698, depth698 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l699
					}
					position++
					goto l698
				l699:
					position, tokenIndex, depth = position698, tokenIndex698, depth698
					if buffer[position] != rune('Y') {
						goto l688
					}
					position++
				}
			l698:
				if !_rules[rulesp]() {
					goto l688
				}
				if !_rules[ruleNumericLiteral]() {
					goto l688
				}
				if !_rules[rulespOpt]() {
					goto l688
				}
				{
					position700, tokenIndex700, depth700 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l700
					}
					position++
					goto l701
				l700:
					position, tokenIndex, depth = position700, tokenIndex700, depth700
				}
			l701:
				if !_rules[rulespOpt]() {
					goto l688
				}
				{
					position702, tokenIndex702, depth702 := position, tokenIndex, depth
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
							goto l703
						}
						position++
					}
				l704:
					{
						position706, tokenIndex706, depth706 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l707
						}
						position++
						goto l706
					l707:
						position, tokenIndex, depth = position706, tokenIndex706, depth706
						if buffer[position] != rune('T') {
							goto l703
						}
						position++
					}
				l706:
					goto l702
				l703:
					position, tokenIndex, depth = position702, tokenIndex702, depth702
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l710
						}
						position++
						goto l709
					l710:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
						if buffer[position] != rune('N') {
							goto l708
						}
						position++
					}
				l709:
					{
						position711, tokenIndex711, depth711 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l712
						}
						position++
						goto l711
					l712:
						position, tokenIndex, depth = position711, tokenIndex711, depth711
						if buffer[position] != rune('D') {
							goto l708
						}
						position++
					}
				l711:
					goto l702
				l708:
					position, tokenIndex, depth = position702, tokenIndex702, depth702
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
							goto l713
						}
						position++
					}
				l714:
					{
						position716, tokenIndex716, depth716 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l717
						}
						position++
						goto l716
					l717:
						position, tokenIndex, depth = position716, tokenIndex716, depth716
						if buffer[position] != rune('D') {
							goto l713
						}
						position++
					}
				l716:
					goto l702
				l713:
					position, tokenIndex, depth = position702, tokenIndex702, depth702
					{
						position718, tokenIndex718, depth718 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l719
						}
						position++
						goto l718
					l719:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
						if buffer[position] != rune('T') {
							goto l688
						}
						position++
					}
				l718:
					{
						position720, tokenIndex720, depth720 := position, tokenIndex, depth
						if buffer[position] != rune('h') {
							goto l721
						}
						position++
						goto l720
					l721:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
						if buffer[position] != rune('H') {
							goto l688
						}
						position++
					}
				l720:
				}
			l702:
				if !_rules[rulesp]() {
					goto l688
				}
				{
					position722, tokenIndex722, depth722 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l723
					}
					position++
					goto l722
				l723:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if buffer[position] != rune('T') {
						goto l688
					}
					position++
				}
			l722:
				{
					position724, tokenIndex724, depth724 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l725
					}
					position++
					goto l724
				l725:
					position, tokenIndex, depth = position724, tokenIndex724, depth724
					if buffer[position] != rune('U') {
						goto l688
					}
					position++
				}
			l724:
				{
					position726, tokenIndex726, depth726 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l727
					}
					position++
					goto l726
				l727:
					position, tokenIndex, depth = position726, tokenIndex726, depth726
					if buffer[position] != rune('P') {
						goto l688
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
						goto l688
					}
					position++
				}
			l728:
				{
					position730, tokenIndex730, depth730 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l731
					}
					position++
					goto l730
				l731:
					position, tokenIndex, depth = position730, tokenIndex730, depth730
					if buffer[position] != rune('E') {
						goto l688
					}
					position++
				}
			l730:
				if !_rules[ruleAction28]() {
					goto l688
				}
				depth--
				add(ruleCountBasedSampling, position689)
			}
			return true
		l688:
			position, tokenIndex, depth = position688, tokenIndex688, depth688
			return false
		},
		/* 37 RandomizedSampling <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') sp (FloatLiteral / NumericLiteral) spOpt '%' Action29)> */
		func() bool {
			position732, tokenIndex732, depth732 := position, tokenIndex, depth
			{
				position733 := position
				depth++
				{
					position734, tokenIndex734, depth734 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l735
					}
					position++
					goto l734
				l735:
					position, tokenIndex, depth = position734, tokenIndex734, depth734
					if buffer[position] != rune('S') {
						goto l732
					}
					position++
				}
			l734:
				{
					position736, tokenIndex736, depth736 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l737
					}
					position++
					goto l736
				l737:
					position, tokenIndex, depth = position736, tokenIndex736, depth736
					if buffer[position] != rune('A') {
						goto l732
					}
					position++
				}
			l736:
				{
					position738, tokenIndex738, depth738 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l739
					}
					position++
					goto l738
				l739:
					position, tokenIndex, depth = position738, tokenIndex738, depth738
					if buffer[position] != rune('M') {
						goto l732
					}
					position++
				}
			l738:
				{
					position740, tokenIndex740, depth740 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l741
					}
					position++
					goto l740
				l741:
					position, tokenIndex, depth = position740, tokenIndex740, depth740
					if buffer[position] != rune('P') {
						goto l732
					}
					position++
				}
			l740:
				{
					position742, tokenIndex742, depth742 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l743
					}
					position++
					goto l742
				l743:
					position, tokenIndex, depth = position742, tokenIndex742, depth742
					if buffer[position] != rune('L') {
						goto l732
					}
					position++
				}
			l742:
				{
					position744, tokenIndex744, depth744 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l745
					}
					position++
					goto l744
				l745:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
					if buffer[position] != rune('E') {
						goto l732
					}
					position++
				}
			l744:
				if !_rules[rulesp]() {
					goto l732
				}
				{
					position746, tokenIndex746, depth746 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l747
					}
					goto l746
				l747:
					position, tokenIndex, depth = position746, tokenIndex746, depth746
					if !_rules[ruleNumericLiteral]() {
						goto l732
					}
				}
			l746:
				if !_rules[rulespOpt]() {
					goto l732
				}
				if buffer[position] != rune('%') {
					goto l732
				}
				position++
				if !_rules[ruleAction29]() {
					goto l732
				}
				depth--
				add(ruleRandomizedSampling, position733)
			}
			return true
		l732:
			position, tokenIndex, depth = position732, tokenIndex732, depth732
			return false
		},
		/* 38 TimeBasedSampling <- <(TimeBasedSamplingSeconds / TimeBasedSamplingMilliseconds)> */
		func() bool {
			position748, tokenIndex748, depth748 := position, tokenIndex, depth
			{
				position749 := position
				depth++
				{
					position750, tokenIndex750, depth750 := position, tokenIndex, depth
					if !_rules[ruleTimeBasedSamplingSeconds]() {
						goto l751
					}
					goto l750
				l751:
					position, tokenIndex, depth = position750, tokenIndex750, depth750
					if !_rules[ruleTimeBasedSamplingMilliseconds]() {
						goto l748
					}
				}
			l750:
				depth--
				add(ruleTimeBasedSampling, position749)
			}
			return true
		l748:
			position, tokenIndex, depth = position748, tokenIndex748, depth748
			return false
		},
		/* 39 TimeBasedSamplingSeconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action30)> */
		func() bool {
			position752, tokenIndex752, depth752 := position, tokenIndex, depth
			{
				position753 := position
				depth++
				{
					position754, tokenIndex754, depth754 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l755
					}
					position++
					goto l754
				l755:
					position, tokenIndex, depth = position754, tokenIndex754, depth754
					if buffer[position] != rune('E') {
						goto l752
					}
					position++
				}
			l754:
				{
					position756, tokenIndex756, depth756 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l757
					}
					position++
					goto l756
				l757:
					position, tokenIndex, depth = position756, tokenIndex756, depth756
					if buffer[position] != rune('V') {
						goto l752
					}
					position++
				}
			l756:
				{
					position758, tokenIndex758, depth758 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l759
					}
					position++
					goto l758
				l759:
					position, tokenIndex, depth = position758, tokenIndex758, depth758
					if buffer[position] != rune('E') {
						goto l752
					}
					position++
				}
			l758:
				{
					position760, tokenIndex760, depth760 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l761
					}
					position++
					goto l760
				l761:
					position, tokenIndex, depth = position760, tokenIndex760, depth760
					if buffer[position] != rune('R') {
						goto l752
					}
					position++
				}
			l760:
				{
					position762, tokenIndex762, depth762 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l763
					}
					position++
					goto l762
				l763:
					position, tokenIndex, depth = position762, tokenIndex762, depth762
					if buffer[position] != rune('Y') {
						goto l752
					}
					position++
				}
			l762:
				if !_rules[rulesp]() {
					goto l752
				}
				{
					position764, tokenIndex764, depth764 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l765
					}
					goto l764
				l765:
					position, tokenIndex, depth = position764, tokenIndex764, depth764
					if !_rules[ruleNumericLiteral]() {
						goto l752
					}
				}
			l764:
				if !_rules[rulesp]() {
					goto l752
				}
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
						goto l752
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
						goto l752
					}
					position++
				}
			l768:
				{
					position770, tokenIndex770, depth770 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l771
					}
					position++
					goto l770
				l771:
					position, tokenIndex, depth = position770, tokenIndex770, depth770
					if buffer[position] != rune('C') {
						goto l752
					}
					position++
				}
			l770:
				{
					position772, tokenIndex772, depth772 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l773
					}
					position++
					goto l772
				l773:
					position, tokenIndex, depth = position772, tokenIndex772, depth772
					if buffer[position] != rune('O') {
						goto l752
					}
					position++
				}
			l772:
				{
					position774, tokenIndex774, depth774 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l775
					}
					position++
					goto l774
				l775:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if buffer[position] != rune('N') {
						goto l752
					}
					position++
				}
			l774:
				{
					position776, tokenIndex776, depth776 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l777
					}
					position++
					goto l776
				l777:
					position, tokenIndex, depth = position776, tokenIndex776, depth776
					if buffer[position] != rune('D') {
						goto l752
					}
					position++
				}
			l776:
				{
					position778, tokenIndex778, depth778 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l779
					}
					position++
					goto l778
				l779:
					position, tokenIndex, depth = position778, tokenIndex778, depth778
					if buffer[position] != rune('S') {
						goto l752
					}
					position++
				}
			l778:
				if !_rules[ruleAction30]() {
					goto l752
				}
				depth--
				add(ruleTimeBasedSamplingSeconds, position753)
			}
			return true
		l752:
			position, tokenIndex, depth = position752, tokenIndex752, depth752
			return false
		},
		/* 40 TimeBasedSamplingMilliseconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action31)> */
		func() bool {
			position780, tokenIndex780, depth780 := position, tokenIndex, depth
			{
				position781 := position
				depth++
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
						goto l780
					}
					position++
				}
			l782:
				{
					position784, tokenIndex784, depth784 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l785
					}
					position++
					goto l784
				l785:
					position, tokenIndex, depth = position784, tokenIndex784, depth784
					if buffer[position] != rune('V') {
						goto l780
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
						goto l780
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
						goto l780
					}
					position++
				}
			l788:
				{
					position790, tokenIndex790, depth790 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l791
					}
					position++
					goto l790
				l791:
					position, tokenIndex, depth = position790, tokenIndex790, depth790
					if buffer[position] != rune('Y') {
						goto l780
					}
					position++
				}
			l790:
				if !_rules[rulesp]() {
					goto l780
				}
				{
					position792, tokenIndex792, depth792 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l793
					}
					goto l792
				l793:
					position, tokenIndex, depth = position792, tokenIndex792, depth792
					if !_rules[ruleNumericLiteral]() {
						goto l780
					}
				}
			l792:
				if !_rules[rulesp]() {
					goto l780
				}
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
						goto l780
					}
					position++
				}
			l794:
				{
					position796, tokenIndex796, depth796 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l797
					}
					position++
					goto l796
				l797:
					position, tokenIndex, depth = position796, tokenIndex796, depth796
					if buffer[position] != rune('I') {
						goto l780
					}
					position++
				}
			l796:
				{
					position798, tokenIndex798, depth798 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l799
					}
					position++
					goto l798
				l799:
					position, tokenIndex, depth = position798, tokenIndex798, depth798
					if buffer[position] != rune('L') {
						goto l780
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
						goto l780
					}
					position++
				}
			l800:
				{
					position802, tokenIndex802, depth802 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l803
					}
					position++
					goto l802
				l803:
					position, tokenIndex, depth = position802, tokenIndex802, depth802
					if buffer[position] != rune('I') {
						goto l780
					}
					position++
				}
			l802:
				{
					position804, tokenIndex804, depth804 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l805
					}
					position++
					goto l804
				l805:
					position, tokenIndex, depth = position804, tokenIndex804, depth804
					if buffer[position] != rune('S') {
						goto l780
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
						goto l780
					}
					position++
				}
			l806:
				{
					position808, tokenIndex808, depth808 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l809
					}
					position++
					goto l808
				l809:
					position, tokenIndex, depth = position808, tokenIndex808, depth808
					if buffer[position] != rune('C') {
						goto l780
					}
					position++
				}
			l808:
				{
					position810, tokenIndex810, depth810 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l811
					}
					position++
					goto l810
				l811:
					position, tokenIndex, depth = position810, tokenIndex810, depth810
					if buffer[position] != rune('O') {
						goto l780
					}
					position++
				}
			l810:
				{
					position812, tokenIndex812, depth812 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l813
					}
					position++
					goto l812
				l813:
					position, tokenIndex, depth = position812, tokenIndex812, depth812
					if buffer[position] != rune('N') {
						goto l780
					}
					position++
				}
			l812:
				{
					position814, tokenIndex814, depth814 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l815
					}
					position++
					goto l814
				l815:
					position, tokenIndex, depth = position814, tokenIndex814, depth814
					if buffer[position] != rune('D') {
						goto l780
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
						goto l780
					}
					position++
				}
			l816:
				if !_rules[ruleAction31]() {
					goto l780
				}
				depth--
				add(ruleTimeBasedSamplingMilliseconds, position781)
			}
			return true
		l780:
			position, tokenIndex, depth = position780, tokenIndex780, depth780
			return false
		},
		/* 41 Projections <- <(<(sp Projection (spOpt ',' spOpt Projection)*)> Action32)> */
		func() bool {
			position818, tokenIndex818, depth818 := position, tokenIndex, depth
			{
				position819 := position
				depth++
				{
					position820 := position
					depth++
					if !_rules[rulesp]() {
						goto l818
					}
					if !_rules[ruleProjection]() {
						goto l818
					}
				l821:
					{
						position822, tokenIndex822, depth822 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l822
						}
						if buffer[position] != rune(',') {
							goto l822
						}
						position++
						if !_rules[rulespOpt]() {
							goto l822
						}
						if !_rules[ruleProjection]() {
							goto l822
						}
						goto l821
					l822:
						position, tokenIndex, depth = position822, tokenIndex822, depth822
					}
					depth--
					add(rulePegText, position820)
				}
				if !_rules[ruleAction32]() {
					goto l818
				}
				depth--
				add(ruleProjections, position819)
			}
			return true
		l818:
			position, tokenIndex, depth = position818, tokenIndex818, depth818
			return false
		},
		/* 42 Projection <- <(AliasExpression / ExpressionOrWildcard)> */
		func() bool {
			position823, tokenIndex823, depth823 := position, tokenIndex, depth
			{
				position824 := position
				depth++
				{
					position825, tokenIndex825, depth825 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l826
					}
					goto l825
				l826:
					position, tokenIndex, depth = position825, tokenIndex825, depth825
					if !_rules[ruleExpressionOrWildcard]() {
						goto l823
					}
				}
			l825:
				depth--
				add(ruleProjection, position824)
			}
			return true
		l823:
			position, tokenIndex, depth = position823, tokenIndex823, depth823
			return false
		},
		/* 43 AliasExpression <- <(ExpressionOrWildcard sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action33)> */
		func() bool {
			position827, tokenIndex827, depth827 := position, tokenIndex, depth
			{
				position828 := position
				depth++
				if !_rules[ruleExpressionOrWildcard]() {
					goto l827
				}
				if !_rules[rulesp]() {
					goto l827
				}
				{
					position829, tokenIndex829, depth829 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l830
					}
					position++
					goto l829
				l830:
					position, tokenIndex, depth = position829, tokenIndex829, depth829
					if buffer[position] != rune('A') {
						goto l827
					}
					position++
				}
			l829:
				{
					position831, tokenIndex831, depth831 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l832
					}
					position++
					goto l831
				l832:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					if buffer[position] != rune('S') {
						goto l827
					}
					position++
				}
			l831:
				if !_rules[rulesp]() {
					goto l827
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l827
				}
				if !_rules[ruleAction33]() {
					goto l827
				}
				depth--
				add(ruleAliasExpression, position828)
			}
			return true
		l827:
			position, tokenIndex, depth = position827, tokenIndex827, depth827
			return false
		},
		/* 44 WindowedFrom <- <(<(sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp Relations)?> Action34)> */
		func() bool {
			position833, tokenIndex833, depth833 := position, tokenIndex, depth
			{
				position834 := position
				depth++
				{
					position835 := position
					depth++
					{
						position836, tokenIndex836, depth836 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l836
						}
						{
							position838, tokenIndex838, depth838 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l839
							}
							position++
							goto l838
						l839:
							position, tokenIndex, depth = position838, tokenIndex838, depth838
							if buffer[position] != rune('F') {
								goto l836
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
								goto l836
							}
							position++
						}
					l840:
						{
							position842, tokenIndex842, depth842 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l843
							}
							position++
							goto l842
						l843:
							position, tokenIndex, depth = position842, tokenIndex842, depth842
							if buffer[position] != rune('O') {
								goto l836
							}
							position++
						}
					l842:
						{
							position844, tokenIndex844, depth844 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l845
							}
							position++
							goto l844
						l845:
							position, tokenIndex, depth = position844, tokenIndex844, depth844
							if buffer[position] != rune('M') {
								goto l836
							}
							position++
						}
					l844:
						if !_rules[rulesp]() {
							goto l836
						}
						if !_rules[ruleRelations]() {
							goto l836
						}
						goto l837
					l836:
						position, tokenIndex, depth = position836, tokenIndex836, depth836
					}
				l837:
					depth--
					add(rulePegText, position835)
				}
				if !_rules[ruleAction34]() {
					goto l833
				}
				depth--
				add(ruleWindowedFrom, position834)
			}
			return true
		l833:
			position, tokenIndex, depth = position833, tokenIndex833, depth833
			return false
		},
		/* 45 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position846, tokenIndex846, depth846 := position, tokenIndex, depth
			{
				position847 := position
				depth++
				{
					position848, tokenIndex848, depth848 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l849
					}
					goto l848
				l849:
					position, tokenIndex, depth = position848, tokenIndex848, depth848
					if !_rules[ruleTuplesInterval]() {
						goto l846
					}
				}
			l848:
				depth--
				add(ruleInterval, position847)
			}
			return true
		l846:
			position, tokenIndex, depth = position846, tokenIndex846, depth846
			return false
		},
		/* 46 TimeInterval <- <((FloatLiteral / NumericLiteral) sp (SECONDS / MILLISECONDS) Action35)> */
		func() bool {
			position850, tokenIndex850, depth850 := position, tokenIndex, depth
			{
				position851 := position
				depth++
				{
					position852, tokenIndex852, depth852 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l853
					}
					goto l852
				l853:
					position, tokenIndex, depth = position852, tokenIndex852, depth852
					if !_rules[ruleNumericLiteral]() {
						goto l850
					}
				}
			l852:
				if !_rules[rulesp]() {
					goto l850
				}
				{
					position854, tokenIndex854, depth854 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l855
					}
					goto l854
				l855:
					position, tokenIndex, depth = position854, tokenIndex854, depth854
					if !_rules[ruleMILLISECONDS]() {
						goto l850
					}
				}
			l854:
				if !_rules[ruleAction35]() {
					goto l850
				}
				depth--
				add(ruleTimeInterval, position851)
			}
			return true
		l850:
			position, tokenIndex, depth = position850, tokenIndex850, depth850
			return false
		},
		/* 47 TuplesInterval <- <(NumericLiteral sp TUPLES Action36)> */
		func() bool {
			position856, tokenIndex856, depth856 := position, tokenIndex, depth
			{
				position857 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l856
				}
				if !_rules[rulesp]() {
					goto l856
				}
				if !_rules[ruleTUPLES]() {
					goto l856
				}
				if !_rules[ruleAction36]() {
					goto l856
				}
				depth--
				add(ruleTuplesInterval, position857)
			}
			return true
		l856:
			position, tokenIndex, depth = position856, tokenIndex856, depth856
			return false
		},
		/* 48 Relations <- <(RelationLike (spOpt ',' spOpt RelationLike)*)> */
		func() bool {
			position858, tokenIndex858, depth858 := position, tokenIndex, depth
			{
				position859 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l858
				}
			l860:
				{
					position861, tokenIndex861, depth861 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l861
					}
					if buffer[position] != rune(',') {
						goto l861
					}
					position++
					if !_rules[rulespOpt]() {
						goto l861
					}
					if !_rules[ruleRelationLike]() {
						goto l861
					}
					goto l860
				l861:
					position, tokenIndex, depth = position861, tokenIndex861, depth861
				}
				depth--
				add(ruleRelations, position859)
			}
			return true
		l858:
			position, tokenIndex, depth = position858, tokenIndex858, depth858
			return false
		},
		/* 49 Filter <- <(<(sp (('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E')) sp Expression)?> Action37)> */
		func() bool {
			position862, tokenIndex862, depth862 := position, tokenIndex, depth
			{
				position863 := position
				depth++
				{
					position864 := position
					depth++
					{
						position865, tokenIndex865, depth865 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l865
						}
						{
							position867, tokenIndex867, depth867 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l868
							}
							position++
							goto l867
						l868:
							position, tokenIndex, depth = position867, tokenIndex867, depth867
							if buffer[position] != rune('W') {
								goto l865
							}
							position++
						}
					l867:
						{
							position869, tokenIndex869, depth869 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l870
							}
							position++
							goto l869
						l870:
							position, tokenIndex, depth = position869, tokenIndex869, depth869
							if buffer[position] != rune('H') {
								goto l865
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
								goto l865
							}
							position++
						}
					l871:
						{
							position873, tokenIndex873, depth873 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l874
							}
							position++
							goto l873
						l874:
							position, tokenIndex, depth = position873, tokenIndex873, depth873
							if buffer[position] != rune('R') {
								goto l865
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
								goto l865
							}
							position++
						}
					l875:
						if !_rules[rulesp]() {
							goto l865
						}
						if !_rules[ruleExpression]() {
							goto l865
						}
						goto l866
					l865:
						position, tokenIndex, depth = position865, tokenIndex865, depth865
					}
				l866:
					depth--
					add(rulePegText, position864)
				}
				if !_rules[ruleAction37]() {
					goto l862
				}
				depth--
				add(ruleFilter, position863)
			}
			return true
		l862:
			position, tokenIndex, depth = position862, tokenIndex862, depth862
			return false
		},
		/* 50 Grouping <- <(<(sp (('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P')) sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action38)> */
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
						if !_rules[rulesp]() {
							goto l880
						}
						{
							position882, tokenIndex882, depth882 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l883
							}
							position++
							goto l882
						l883:
							position, tokenIndex, depth = position882, tokenIndex882, depth882
							if buffer[position] != rune('G') {
								goto l880
							}
							position++
						}
					l882:
						{
							position884, tokenIndex884, depth884 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l885
							}
							position++
							goto l884
						l885:
							position, tokenIndex, depth = position884, tokenIndex884, depth884
							if buffer[position] != rune('R') {
								goto l880
							}
							position++
						}
					l884:
						{
							position886, tokenIndex886, depth886 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l887
							}
							position++
							goto l886
						l887:
							position, tokenIndex, depth = position886, tokenIndex886, depth886
							if buffer[position] != rune('O') {
								goto l880
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
								goto l880
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
								goto l880
							}
							position++
						}
					l890:
						if !_rules[rulesp]() {
							goto l880
						}
						{
							position892, tokenIndex892, depth892 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l893
							}
							position++
							goto l892
						l893:
							position, tokenIndex, depth = position892, tokenIndex892, depth892
							if buffer[position] != rune('B') {
								goto l880
							}
							position++
						}
					l892:
						{
							position894, tokenIndex894, depth894 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l895
							}
							position++
							goto l894
						l895:
							position, tokenIndex, depth = position894, tokenIndex894, depth894
							if buffer[position] != rune('Y') {
								goto l880
							}
							position++
						}
					l894:
						if !_rules[rulesp]() {
							goto l880
						}
						if !_rules[ruleGroupList]() {
							goto l880
						}
						goto l881
					l880:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
					}
				l881:
					depth--
					add(rulePegText, position879)
				}
				if !_rules[ruleAction38]() {
					goto l877
				}
				depth--
				add(ruleGrouping, position878)
			}
			return true
		l877:
			position, tokenIndex, depth = position877, tokenIndex877, depth877
			return false
		},
		/* 51 GroupList <- <(Expression (spOpt ',' spOpt Expression)*)> */
		func() bool {
			position896, tokenIndex896, depth896 := position, tokenIndex, depth
			{
				position897 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l896
				}
			l898:
				{
					position899, tokenIndex899, depth899 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l899
					}
					if buffer[position] != rune(',') {
						goto l899
					}
					position++
					if !_rules[rulespOpt]() {
						goto l899
					}
					if !_rules[ruleExpression]() {
						goto l899
					}
					goto l898
				l899:
					position, tokenIndex, depth = position899, tokenIndex899, depth899
				}
				depth--
				add(ruleGroupList, position897)
			}
			return true
		l896:
			position, tokenIndex, depth = position896, tokenIndex896, depth896
			return false
		},
		/* 52 Having <- <(<(sp (('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G')) sp Expression)?> Action39)> */
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
						if !_rules[rulesp]() {
							goto l903
						}
						{
							position905, tokenIndex905, depth905 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l906
							}
							position++
							goto l905
						l906:
							position, tokenIndex, depth = position905, tokenIndex905, depth905
							if buffer[position] != rune('H') {
								goto l903
							}
							position++
						}
					l905:
						{
							position907, tokenIndex907, depth907 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l908
							}
							position++
							goto l907
						l908:
							position, tokenIndex, depth = position907, tokenIndex907, depth907
							if buffer[position] != rune('A') {
								goto l903
							}
							position++
						}
					l907:
						{
							position909, tokenIndex909, depth909 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l910
							}
							position++
							goto l909
						l910:
							position, tokenIndex, depth = position909, tokenIndex909, depth909
							if buffer[position] != rune('V') {
								goto l903
							}
							position++
						}
					l909:
						{
							position911, tokenIndex911, depth911 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l912
							}
							position++
							goto l911
						l912:
							position, tokenIndex, depth = position911, tokenIndex911, depth911
							if buffer[position] != rune('I') {
								goto l903
							}
							position++
						}
					l911:
						{
							position913, tokenIndex913, depth913 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l914
							}
							position++
							goto l913
						l914:
							position, tokenIndex, depth = position913, tokenIndex913, depth913
							if buffer[position] != rune('N') {
								goto l903
							}
							position++
						}
					l913:
						{
							position915, tokenIndex915, depth915 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l916
							}
							position++
							goto l915
						l916:
							position, tokenIndex, depth = position915, tokenIndex915, depth915
							if buffer[position] != rune('G') {
								goto l903
							}
							position++
						}
					l915:
						if !_rules[rulesp]() {
							goto l903
						}
						if !_rules[ruleExpression]() {
							goto l903
						}
						goto l904
					l903:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
					}
				l904:
					depth--
					add(rulePegText, position902)
				}
				if !_rules[ruleAction39]() {
					goto l900
				}
				depth--
				add(ruleHaving, position901)
			}
			return true
		l900:
			position, tokenIndex, depth = position900, tokenIndex900, depth900
			return false
		},
		/* 53 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action40))> */
		func() bool {
			position917, tokenIndex917, depth917 := position, tokenIndex, depth
			{
				position918 := position
				depth++
				{
					position919, tokenIndex919, depth919 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l920
					}
					goto l919
				l920:
					position, tokenIndex, depth = position919, tokenIndex919, depth919
					if !_rules[ruleStreamWindow]() {
						goto l917
					}
					if !_rules[ruleAction40]() {
						goto l917
					}
				}
			l919:
				depth--
				add(ruleRelationLike, position918)
			}
			return true
		l917:
			position, tokenIndex, depth = position917, tokenIndex917, depth917
			return false
		},
		/* 54 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action41)> */
		func() bool {
			position921, tokenIndex921, depth921 := position, tokenIndex, depth
			{
				position922 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l921
				}
				if !_rules[rulesp]() {
					goto l921
				}
				{
					position923, tokenIndex923, depth923 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l924
					}
					position++
					goto l923
				l924:
					position, tokenIndex, depth = position923, tokenIndex923, depth923
					if buffer[position] != rune('A') {
						goto l921
					}
					position++
				}
			l923:
				{
					position925, tokenIndex925, depth925 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l926
					}
					position++
					goto l925
				l926:
					position, tokenIndex, depth = position925, tokenIndex925, depth925
					if buffer[position] != rune('S') {
						goto l921
					}
					position++
				}
			l925:
				if !_rules[rulesp]() {
					goto l921
				}
				if !_rules[ruleIdentifier]() {
					goto l921
				}
				if !_rules[ruleAction41]() {
					goto l921
				}
				depth--
				add(ruleAliasedStreamWindow, position922)
			}
			return true
		l921:
			position, tokenIndex, depth = position921, tokenIndex921, depth921
			return false
		},
		/* 55 StreamWindow <- <(StreamLike spOpt '[' spOpt (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval CapacitySpecOpt SheddingSpecOpt spOpt ']' Action42)> */
		func() bool {
			position927, tokenIndex927, depth927 := position, tokenIndex, depth
			{
				position928 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l927
				}
				if !_rules[rulespOpt]() {
					goto l927
				}
				if buffer[position] != rune('[') {
					goto l927
				}
				position++
				if !_rules[rulespOpt]() {
					goto l927
				}
				{
					position929, tokenIndex929, depth929 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l930
					}
					position++
					goto l929
				l930:
					position, tokenIndex, depth = position929, tokenIndex929, depth929
					if buffer[position] != rune('R') {
						goto l927
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
						goto l927
					}
					position++
				}
			l931:
				{
					position933, tokenIndex933, depth933 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l934
					}
					position++
					goto l933
				l934:
					position, tokenIndex, depth = position933, tokenIndex933, depth933
					if buffer[position] != rune('N') {
						goto l927
					}
					position++
				}
			l933:
				{
					position935, tokenIndex935, depth935 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l936
					}
					position++
					goto l935
				l936:
					position, tokenIndex, depth = position935, tokenIndex935, depth935
					if buffer[position] != rune('G') {
						goto l927
					}
					position++
				}
			l935:
				{
					position937, tokenIndex937, depth937 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l938
					}
					position++
					goto l937
				l938:
					position, tokenIndex, depth = position937, tokenIndex937, depth937
					if buffer[position] != rune('E') {
						goto l927
					}
					position++
				}
			l937:
				if !_rules[rulesp]() {
					goto l927
				}
				if !_rules[ruleInterval]() {
					goto l927
				}
				if !_rules[ruleCapacitySpecOpt]() {
					goto l927
				}
				if !_rules[ruleSheddingSpecOpt]() {
					goto l927
				}
				if !_rules[rulespOpt]() {
					goto l927
				}
				if buffer[position] != rune(']') {
					goto l927
				}
				position++
				if !_rules[ruleAction42]() {
					goto l927
				}
				depth--
				add(ruleStreamWindow, position928)
			}
			return true
		l927:
			position, tokenIndex, depth = position927, tokenIndex927, depth927
			return false
		},
		/* 56 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position939, tokenIndex939, depth939 := position, tokenIndex, depth
			{
				position940 := position
				depth++
				{
					position941, tokenIndex941, depth941 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l942
					}
					goto l941
				l942:
					position, tokenIndex, depth = position941, tokenIndex941, depth941
					if !_rules[ruleStream]() {
						goto l939
					}
				}
			l941:
				depth--
				add(ruleStreamLike, position940)
			}
			return true
		l939:
			position, tokenIndex, depth = position939, tokenIndex939, depth939
			return false
		},
		/* 57 UDSFFuncApp <- <(FuncAppWithoutOrderBy Action43)> */
		func() bool {
			position943, tokenIndex943, depth943 := position, tokenIndex, depth
			{
				position944 := position
				depth++
				if !_rules[ruleFuncAppWithoutOrderBy]() {
					goto l943
				}
				if !_rules[ruleAction43]() {
					goto l943
				}
				depth--
				add(ruleUDSFFuncApp, position944)
			}
			return true
		l943:
			position, tokenIndex, depth = position943, tokenIndex943, depth943
			return false
		},
		/* 58 CapacitySpecOpt <- <(<(spOpt ',' spOpt (('b' / 'B') ('u' / 'U') ('f' / 'F') ('f' / 'F') ('e' / 'E') ('r' / 'R')) sp (('s' / 'S') ('i' / 'I') ('z' / 'Z') ('e' / 'E')) sp NonNegativeNumericLiteral)?> Action44)> */
		func() bool {
			position945, tokenIndex945, depth945 := position, tokenIndex, depth
			{
				position946 := position
				depth++
				{
					position947 := position
					depth++
					{
						position948, tokenIndex948, depth948 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l948
						}
						if buffer[position] != rune(',') {
							goto l948
						}
						position++
						if !_rules[rulespOpt]() {
							goto l948
						}
						{
							position950, tokenIndex950, depth950 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l951
							}
							position++
							goto l950
						l951:
							position, tokenIndex, depth = position950, tokenIndex950, depth950
							if buffer[position] != rune('B') {
								goto l948
							}
							position++
						}
					l950:
						{
							position952, tokenIndex952, depth952 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l953
							}
							position++
							goto l952
						l953:
							position, tokenIndex, depth = position952, tokenIndex952, depth952
							if buffer[position] != rune('U') {
								goto l948
							}
							position++
						}
					l952:
						{
							position954, tokenIndex954, depth954 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l955
							}
							position++
							goto l954
						l955:
							position, tokenIndex, depth = position954, tokenIndex954, depth954
							if buffer[position] != rune('F') {
								goto l948
							}
							position++
						}
					l954:
						{
							position956, tokenIndex956, depth956 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l957
							}
							position++
							goto l956
						l957:
							position, tokenIndex, depth = position956, tokenIndex956, depth956
							if buffer[position] != rune('F') {
								goto l948
							}
							position++
						}
					l956:
						{
							position958, tokenIndex958, depth958 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l959
							}
							position++
							goto l958
						l959:
							position, tokenIndex, depth = position958, tokenIndex958, depth958
							if buffer[position] != rune('E') {
								goto l948
							}
							position++
						}
					l958:
						{
							position960, tokenIndex960, depth960 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l961
							}
							position++
							goto l960
						l961:
							position, tokenIndex, depth = position960, tokenIndex960, depth960
							if buffer[position] != rune('R') {
								goto l948
							}
							position++
						}
					l960:
						if !_rules[rulesp]() {
							goto l948
						}
						{
							position962, tokenIndex962, depth962 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l963
							}
							position++
							goto l962
						l963:
							position, tokenIndex, depth = position962, tokenIndex962, depth962
							if buffer[position] != rune('S') {
								goto l948
							}
							position++
						}
					l962:
						{
							position964, tokenIndex964, depth964 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l965
							}
							position++
							goto l964
						l965:
							position, tokenIndex, depth = position964, tokenIndex964, depth964
							if buffer[position] != rune('I') {
								goto l948
							}
							position++
						}
					l964:
						{
							position966, tokenIndex966, depth966 := position, tokenIndex, depth
							if buffer[position] != rune('z') {
								goto l967
							}
							position++
							goto l966
						l967:
							position, tokenIndex, depth = position966, tokenIndex966, depth966
							if buffer[position] != rune('Z') {
								goto l948
							}
							position++
						}
					l966:
						{
							position968, tokenIndex968, depth968 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l969
							}
							position++
							goto l968
						l969:
							position, tokenIndex, depth = position968, tokenIndex968, depth968
							if buffer[position] != rune('E') {
								goto l948
							}
							position++
						}
					l968:
						if !_rules[rulesp]() {
							goto l948
						}
						if !_rules[ruleNonNegativeNumericLiteral]() {
							goto l948
						}
						goto l949
					l948:
						position, tokenIndex, depth = position948, tokenIndex948, depth948
					}
				l949:
					depth--
					add(rulePegText, position947)
				}
				if !_rules[ruleAction44]() {
					goto l945
				}
				depth--
				add(ruleCapacitySpecOpt, position946)
			}
			return true
		l945:
			position, tokenIndex, depth = position945, tokenIndex945, depth945
			return false
		},
		/* 59 SheddingSpecOpt <- <(<(spOpt ',' spOpt SheddingOption sp (('i' / 'I') ('f' / 'F')) sp (('f' / 'F') ('u' / 'U') ('l' / 'L') ('l' / 'L')))?> Action45)> */
		func() bool {
			position970, tokenIndex970, depth970 := position, tokenIndex, depth
			{
				position971 := position
				depth++
				{
					position972 := position
					depth++
					{
						position973, tokenIndex973, depth973 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l973
						}
						if buffer[position] != rune(',') {
							goto l973
						}
						position++
						if !_rules[rulespOpt]() {
							goto l973
						}
						if !_rules[ruleSheddingOption]() {
							goto l973
						}
						if !_rules[rulesp]() {
							goto l973
						}
						{
							position975, tokenIndex975, depth975 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l976
							}
							position++
							goto l975
						l976:
							position, tokenIndex, depth = position975, tokenIndex975, depth975
							if buffer[position] != rune('I') {
								goto l973
							}
							position++
						}
					l975:
						{
							position977, tokenIndex977, depth977 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l978
							}
							position++
							goto l977
						l978:
							position, tokenIndex, depth = position977, tokenIndex977, depth977
							if buffer[position] != rune('F') {
								goto l973
							}
							position++
						}
					l977:
						if !_rules[rulesp]() {
							goto l973
						}
						{
							position979, tokenIndex979, depth979 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l980
							}
							position++
							goto l979
						l980:
							position, tokenIndex, depth = position979, tokenIndex979, depth979
							if buffer[position] != rune('F') {
								goto l973
							}
							position++
						}
					l979:
						{
							position981, tokenIndex981, depth981 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l982
							}
							position++
							goto l981
						l982:
							position, tokenIndex, depth = position981, tokenIndex981, depth981
							if buffer[position] != rune('U') {
								goto l973
							}
							position++
						}
					l981:
						{
							position983, tokenIndex983, depth983 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l984
							}
							position++
							goto l983
						l984:
							position, tokenIndex, depth = position983, tokenIndex983, depth983
							if buffer[position] != rune('L') {
								goto l973
							}
							position++
						}
					l983:
						{
							position985, tokenIndex985, depth985 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l986
							}
							position++
							goto l985
						l986:
							position, tokenIndex, depth = position985, tokenIndex985, depth985
							if buffer[position] != rune('L') {
								goto l973
							}
							position++
						}
					l985:
						goto l974
					l973:
						position, tokenIndex, depth = position973, tokenIndex973, depth973
					}
				l974:
					depth--
					add(rulePegText, position972)
				}
				if !_rules[ruleAction45]() {
					goto l970
				}
				depth--
				add(ruleSheddingSpecOpt, position971)
			}
			return true
		l970:
			position, tokenIndex, depth = position970, tokenIndex970, depth970
			return false
		},
		/* 60 SheddingOption <- <(Wait / DropOldest / DropNewest)> */
		func() bool {
			position987, tokenIndex987, depth987 := position, tokenIndex, depth
			{
				position988 := position
				depth++
				{
					position989, tokenIndex989, depth989 := position, tokenIndex, depth
					if !_rules[ruleWait]() {
						goto l990
					}
					goto l989
				l990:
					position, tokenIndex, depth = position989, tokenIndex989, depth989
					if !_rules[ruleDropOldest]() {
						goto l991
					}
					goto l989
				l991:
					position, tokenIndex, depth = position989, tokenIndex989, depth989
					if !_rules[ruleDropNewest]() {
						goto l987
					}
				}
			l989:
				depth--
				add(ruleSheddingOption, position988)
			}
			return true
		l987:
			position, tokenIndex, depth = position987, tokenIndex987, depth987
			return false
		},
		/* 61 SourceSinkSpecs <- <(<(sp (('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action46)> */
		func() bool {
			position992, tokenIndex992, depth992 := position, tokenIndex, depth
			{
				position993 := position
				depth++
				{
					position994 := position
					depth++
					{
						position995, tokenIndex995, depth995 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l995
						}
						{
							position997, tokenIndex997, depth997 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l998
							}
							position++
							goto l997
						l998:
							position, tokenIndex, depth = position997, tokenIndex997, depth997
							if buffer[position] != rune('W') {
								goto l995
							}
							position++
						}
					l997:
						{
							position999, tokenIndex999, depth999 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l1000
							}
							position++
							goto l999
						l1000:
							position, tokenIndex, depth = position999, tokenIndex999, depth999
							if buffer[position] != rune('I') {
								goto l995
							}
							position++
						}
					l999:
						{
							position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1002
							}
							position++
							goto l1001
						l1002:
							position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
							if buffer[position] != rune('T') {
								goto l995
							}
							position++
						}
					l1001:
						{
							position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l1004
							}
							position++
							goto l1003
						l1004:
							position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
							if buffer[position] != rune('H') {
								goto l995
							}
							position++
						}
					l1003:
						if !_rules[rulesp]() {
							goto l995
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l995
						}
					l1005:
						{
							position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1006
							}
							if buffer[position] != rune(',') {
								goto l1006
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1006
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l1006
							}
							goto l1005
						l1006:
							position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
						}
						goto l996
					l995:
						position, tokenIndex, depth = position995, tokenIndex995, depth995
					}
				l996:
					depth--
					add(rulePegText, position994)
				}
				if !_rules[ruleAction46]() {
					goto l992
				}
				depth--
				add(ruleSourceSinkSpecs, position993)
			}
			return true
		l992:
			position, tokenIndex, depth = position992, tokenIndex992, depth992
			return false
		},
		/* 62 UpdateSourceSinkSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)> Action47)> */
		func() bool {
			position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
			{
				position1008 := position
				depth++
				{
					position1009 := position
					depth++
					if !_rules[rulesp]() {
						goto l1007
					}
					{
						position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1011
						}
						position++
						goto l1010
					l1011:
						position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
						if buffer[position] != rune('S') {
							goto l1007
						}
						position++
					}
				l1010:
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('E') {
							goto l1007
						}
						position++
					}
				l1012:
					{
						position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1015
						}
						position++
						goto l1014
					l1015:
						position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
						if buffer[position] != rune('T') {
							goto l1007
						}
						position++
					}
				l1014:
					if !_rules[rulesp]() {
						goto l1007
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l1007
					}
				l1016:
					{
						position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1017
						}
						if buffer[position] != rune(',') {
							goto l1017
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1017
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l1017
						}
						goto l1016
					l1017:
						position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
					}
					depth--
					add(rulePegText, position1009)
				}
				if !_rules[ruleAction47]() {
					goto l1007
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position1008)
			}
			return true
		l1007:
			position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
			return false
		},
		/* 63 SetOptSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action48)> */
		func() bool {
			position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
			{
				position1019 := position
				depth++
				{
					position1020 := position
					depth++
					{
						position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1021
						}
						{
							position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1024
							}
							position++
							goto l1023
						l1024:
							position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
							if buffer[position] != rune('S') {
								goto l1021
							}
							position++
						}
					l1023:
						{
							position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1026
							}
							position++
							goto l1025
						l1026:
							position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
							if buffer[position] != rune('E') {
								goto l1021
							}
							position++
						}
					l1025:
						{
							position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1028
							}
							position++
							goto l1027
						l1028:
							position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
							if buffer[position] != rune('T') {
								goto l1021
							}
							position++
						}
					l1027:
						if !_rules[rulesp]() {
							goto l1021
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l1021
						}
					l1029:
						{
							position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1030
							}
							if buffer[position] != rune(',') {
								goto l1030
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1030
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l1030
							}
							goto l1029
						l1030:
							position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
						}
						goto l1022
					l1021:
						position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
					}
				l1022:
					depth--
					add(rulePegText, position1020)
				}
				if !_rules[ruleAction48]() {
					goto l1018
				}
				depth--
				add(ruleSetOptSpecs, position1019)
			}
			return true
		l1018:
			position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
			return false
		},
		/* 64 StateTagOpt <- <(<(sp (('t' / 'T') ('a' / 'A') ('g' / 'G')) sp Identifier)?> Action49)> */
		func() bool {
			position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
			{
				position1032 := position
				depth++
				{
					position1033 := position
					depth++
					{
						position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1034
						}
						{
							position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1037
							}
							position++
							goto l1036
						l1037:
							position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
							if buffer[position] != rune('T') {
								goto l1034
							}
							position++
						}
					l1036:
						{
							position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1039
							}
							position++
							goto l1038
						l1039:
							position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
							if buffer[position] != rune('A') {
								goto l1034
							}
							position++
						}
					l1038:
						{
							position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l1041
							}
							position++
							goto l1040
						l1041:
							position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
							if buffer[position] != rune('G') {
								goto l1034
							}
							position++
						}
					l1040:
						if !_rules[rulesp]() {
							goto l1034
						}
						if !_rules[ruleIdentifier]() {
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
				if !_rules[ruleAction49]() {
					goto l1031
				}
				depth--
				add(ruleStateTagOpt, position1032)
			}
			return true
		l1031:
			position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
			return false
		},
		/* 65 SourceSinkParam <- <(SourceSinkParamKey spOpt '=' spOpt SourceSinkParamVal Action50)> */
		func() bool {
			position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
			{
				position1043 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l1042
				}
				if !_rules[rulespOpt]() {
					goto l1042
				}
				if buffer[position] != rune('=') {
					goto l1042
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1042
				}
				if !_rules[ruleSourceSinkParamVal]() {
					goto l1042
				}
				if !_rules[ruleAction50]() {
					goto l1042
				}
				depth--
				add(ruleSourceSinkParam, position1043)
			}
			return true
		l1042:
			position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
			return false
		},
		/* 66 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr / ParamMapExpr)> */
		func() bool {
			position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
			{
				position1045 := position
				depth++
				{
					position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l1047
					}
					goto l1046
				l1047:
					position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
					if !_rules[ruleParamArrayExpr]() {
						goto l1048
					}
					goto l1046
				l1048:
					position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
					if !_rules[ruleParamMapExpr]() {
						goto l1044
					}
				}
			l1046:
				depth--
				add(ruleSourceSinkParamVal, position1045)
			}
			return true
		l1044:
			position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
			return false
		},
		/* 67 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
			{
				position1050 := position
				depth++
				{
					position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l1052
					}
					goto l1051
				l1052:
					position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
					if !_rules[ruleLiteral]() {
						goto l1049
					}
				}
			l1051:
				depth--
				add(ruleParamLiteral, position1050)
			}
			return true
		l1049:
			position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
			return false
		},
		/* 68 ParamArrayExpr <- <(<('[' spOpt (ParamLiteral (',' spOpt ParamLiteral)*)? spOpt ','? spOpt ']')> Action51)> */
		func() bool {
			position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
			{
				position1054 := position
				depth++
				{
					position1055 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1053
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1053
					}
					{
						position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l1056
						}
					l1058:
						{
							position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l1059
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1059
							}
							if !_rules[ruleParamLiteral]() {
								goto l1059
							}
							goto l1058
						l1059:
							position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
						}
						goto l1057
					l1056:
						position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
					}
				l1057:
					if !_rules[rulespOpt]() {
						goto l1053
					}
					{
						position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1060
						}
						position++
						goto l1061
					l1060:
						position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
					}
				l1061:
					if !_rules[rulespOpt]() {
						goto l1053
					}
					if buffer[position] != rune(']') {
						goto l1053
					}
					position++
					depth--
					add(rulePegText, position1055)
				}
				if !_rules[ruleAction51]() {
					goto l1053
				}
				depth--
				add(ruleParamArrayExpr, position1054)
			}
			return true
		l1053:
			position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
			return false
		},
		/* 69 ParamMapExpr <- <(<('{' spOpt (ParamKeyValuePair (spOpt ',' spOpt ParamKeyValuePair)*)? spOpt '}')> Action52)> */
		func() bool {
			position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
			{
				position1063 := position
				depth++
				{
					position1064 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1062
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1062
					}
					{
						position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
						if !_rules[ruleParamKeyValuePair]() {
							goto l1065
						}
					l1067:
						{
							position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1068
							}
							if buffer[position] != rune(',') {
								goto l1068
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1068
							}
							if !_rules[ruleParamKeyValuePair]() {
								goto l1068
							}
							goto l1067
						l1068:
							position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
						}
						goto l1066
					l1065:
						position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
					}
				l1066:
					if !_rules[rulespOpt]() {
						goto l1062
					}
					if buffer[position] != rune('}') {
						goto l1062
					}
					position++
					depth--
					add(rulePegText, position1064)
				}
				if !_rules[ruleAction52]() {
					goto l1062
				}
				depth--
				add(ruleParamMapExpr, position1063)
			}
			return true
		l1062:
			position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
			return false
		},
		/* 70 ParamKeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ParamLiteral)> Action53)> */
		func() bool {
			position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
			{
				position1070 := position
				depth++
				{
					position1071 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1069
					}
					if !_rules[rulespOpt]() {
						goto l1069
					}
					if buffer[position] != rune(':') {
						goto l1069
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1069
					}
					if !_rules[ruleParamLiteral]() {
						goto l1069
					}
					depth--
					add(rulePegText, position1071)
				}
				if !_rules[ruleAction53]() {
					goto l1069
				}
				depth--
				add(ruleParamKeyValuePair, position1070)
			}
			return true
		l1069:
			position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
			return false
		},
		/* 71 PausedOpt <- <(<(sp (Paused / Unpaused))?> Action54)> */
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
						if !_rules[rulesp]() {
							goto l1075
						}
						{
							position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l1078
							}
							goto l1077
						l1078:
							position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
							if !_rules[ruleUnpaused]() {
								goto l1075
							}
						}
					l1077:
						goto l1076
					l1075:
						position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
					}
				l1076:
					depth--
					add(rulePegText, position1074)
				}
				if !_rules[ruleAction54]() {
					goto l1072
				}
				depth--
				add(rulePausedOpt, position1073)
			}
			return true
		l1072:
			position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
			return false
		},
		/* 72 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
			{
				position1080 := position
				depth++
				{
					position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
					if !_rules[ruleWildcard]() {
						goto l1082
					}
					goto l1081
				l1082:
					position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
					if !_rules[ruleExpression]() {
						goto l1079
					}
				}
			l1081:
				depth--
				add(ruleExpressionOrWildcard, position1080)
			}
			return true
		l1079:
			position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
			return false
		},
		/* 73 Expression <- <orExpr> */
		func() bool {
			position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
			{
				position1084 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l1083
				}
				depth--
				add(ruleExpression, position1084)
			}
			return true
		l1083:
			position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
			return false
		},
		/* 74 orExpr <- <(<(andExpr (sp Or sp andExpr)*)> Action55)> */
		func() bool {
			position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
			{
				position1086 := position
				depth++
				{
					position1087 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l1085
					}
				l1088:
					{
						position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1089
						}
						if !_rules[ruleOr]() {
							goto l1089
						}
						if !_rules[rulesp]() {
							goto l1089
						}
						if !_rules[ruleandExpr]() {
							goto l1089
						}
						goto l1088
					l1089:
						position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
					}
					depth--
					add(rulePegText, position1087)
				}
				if !_rules[ruleAction55]() {
					goto l1085
				}
				depth--
				add(ruleorExpr, position1086)
			}
			return true
		l1085:
			position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
			return false
		},
		/* 75 andExpr <- <(<(notExpr (sp And sp notExpr)*)> Action56)> */
		func() bool {
			position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
			{
				position1091 := position
				depth++
				{
					position1092 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l1090
					}
				l1093:
					{
						position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1094
						}
						if !_rules[ruleAnd]() {
							goto l1094
						}
						if !_rules[rulesp]() {
							goto l1094
						}
						if !_rules[rulenotExpr]() {
							goto l1094
						}
						goto l1093
					l1094:
						position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
					}
					depth--
					add(rulePegText, position1092)
				}
				if !_rules[ruleAction56]() {
					goto l1090
				}
				depth--
				add(ruleandExpr, position1091)
			}
			return true
		l1090:
			position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
			return false
		},
		/* 76 notExpr <- <(<((Not sp)? comparisonExpr)> Action57)> */
		func() bool {
			position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
			{
				position1096 := position
				depth++
				{
					position1097 := position
					depth++
					{
						position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l1098
						}
						if !_rules[rulesp]() {
							goto l1098
						}
						goto l1099
					l1098:
						position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
					}
				l1099:
					if !_rules[rulecomparisonExpr]() {
						goto l1095
					}
					depth--
					add(rulePegText, position1097)
				}
				if !_rules[ruleAction57]() {
					goto l1095
				}
				depth--
				add(rulenotExpr, position1096)
			}
			return true
		l1095:
			position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
			return false
		},
		/* 77 comparisonExpr <- <(<(otherOpExpr (spOpt ComparisonOp spOpt otherOpExpr)?)> Action58)> */
		func() bool {
			position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
			{
				position1101 := position
				depth++
				{
					position1102 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l1100
					}
					{
						position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1103
						}
						if !_rules[ruleComparisonOp]() {
							goto l1103
						}
						if !_rules[rulespOpt]() {
							goto l1103
						}
						if !_rules[ruleotherOpExpr]() {
							goto l1103
						}
						goto l1104
					l1103:
						position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
					}
				l1104:
					depth--
					add(rulePegText, position1102)
				}
				if !_rules[ruleAction58]() {
					goto l1100
				}
				depth--
				add(rulecomparisonExpr, position1101)
			}
			return true
		l1100:
			position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
			return false
		},
		/* 78 otherOpExpr <- <(<(isExpr (spOpt OtherOp spOpt isExpr)*)> Action59)> */
		func() bool {
			position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
			{
				position1106 := position
				depth++
				{
					position1107 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l1105
					}
				l1108:
					{
						position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1109
						}
						if !_rules[ruleOtherOp]() {
							goto l1109
						}
						if !_rules[rulespOpt]() {
							goto l1109
						}
						if !_rules[ruleisExpr]() {
							goto l1109
						}
						goto l1108
					l1109:
						position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
					}
					depth--
					add(rulePegText, position1107)
				}
				if !_rules[ruleAction59]() {
					goto l1105
				}
				depth--
				add(ruleotherOpExpr, position1106)
			}
			return true
		l1105:
			position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
			return false
		},
		/* 79 isExpr <- <(<(termExpr (sp IsOp sp NullLiteral)?)> Action60)> */
		func() bool {
			position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
			{
				position1111 := position
				depth++
				{
					position1112 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l1110
					}
					{
						position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1113
						}
						if !_rules[ruleIsOp]() {
							goto l1113
						}
						if !_rules[rulesp]() {
							goto l1113
						}
						if !_rules[ruleNullLiteral]() {
							goto l1113
						}
						goto l1114
					l1113:
						position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
					}
				l1114:
					depth--
					add(rulePegText, position1112)
				}
				if !_rules[ruleAction60]() {
					goto l1110
				}
				depth--
				add(ruleisExpr, position1111)
			}
			return true
		l1110:
			position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
			return false
		},
		/* 80 termExpr <- <(<(productExpr (spOpt PlusMinusOp spOpt productExpr)*)> Action61)> */
		func() bool {
			position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
			{
				position1116 := position
				depth++
				{
					position1117 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l1115
					}
				l1118:
					{
						position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1119
						}
						if !_rules[rulePlusMinusOp]() {
							goto l1119
						}
						if !_rules[rulespOpt]() {
							goto l1119
						}
						if !_rules[ruleproductExpr]() {
							goto l1119
						}
						goto l1118
					l1119:
						position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
					}
					depth--
					add(rulePegText, position1117)
				}
				if !_rules[ruleAction61]() {
					goto l1115
				}
				depth--
				add(ruletermExpr, position1116)
			}
			return true
		l1115:
			position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
			return false
		},
		/* 81 productExpr <- <(<(minusExpr (spOpt MultDivOp spOpt minusExpr)*)> Action62)> */
		func() bool {
			position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
			{
				position1121 := position
				depth++
				{
					position1122 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l1120
					}
				l1123:
					{
						position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1124
						}
						if !_rules[ruleMultDivOp]() {
							goto l1124
						}
						if !_rules[rulespOpt]() {
							goto l1124
						}
						if !_rules[ruleminusExpr]() {
							goto l1124
						}
						goto l1123
					l1124:
						position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
					}
					depth--
					add(rulePegText, position1122)
				}
				if !_rules[ruleAction62]() {
					goto l1120
				}
				depth--
				add(ruleproductExpr, position1121)
			}
			return true
		l1120:
			position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
			return false
		},
		/* 82 minusExpr <- <(<((UnaryMinus spOpt)? castExpr)> Action63)> */
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
						if !_rules[ruleUnaryMinus]() {
							goto l1128
						}
						if !_rules[rulespOpt]() {
							goto l1128
						}
						goto l1129
					l1128:
						position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
					}
				l1129:
					if !_rules[rulecastExpr]() {
						goto l1125
					}
					depth--
					add(rulePegText, position1127)
				}
				if !_rules[ruleAction63]() {
					goto l1125
				}
				depth--
				add(ruleminusExpr, position1126)
			}
			return true
		l1125:
			position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
			return false
		},
		/* 83 castExpr <- <(<(baseExpr (spOpt (':' ':') spOpt Type)?)> Action64)> */
		func() bool {
			position1130, tokenIndex1130, depth1130 := position, tokenIndex, depth
			{
				position1131 := position
				depth++
				{
					position1132 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l1130
					}
					{
						position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1133
						}
						if buffer[position] != rune(':') {
							goto l1133
						}
						position++
						if buffer[position] != rune(':') {
							goto l1133
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1133
						}
						if !_rules[ruleType]() {
							goto l1133
						}
						goto l1134
					l1133:
						position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
					}
				l1134:
					depth--
					add(rulePegText, position1132)
				}
				if !_rules[ruleAction64]() {
					goto l1130
				}
				depth--
				add(rulecastExpr, position1131)
			}
			return true
		l1130:
			position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
			return false
		},
		/* 84 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / Case / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
			{
				position1136 := position
				depth++
				{
					position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l1138
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1138
					}
					if !_rules[ruleExpression]() {
						goto l1138
					}
					if !_rules[rulespOpt]() {
						goto l1138
					}
					if buffer[position] != rune(')') {
						goto l1138
					}
					position++
					goto l1137
				l1138:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleMapExpr]() {
						goto l1139
					}
					goto l1137
				l1139:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleBooleanLiteral]() {
						goto l1140
					}
					goto l1137
				l1140:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleNullLiteral]() {
						goto l1141
					}
					goto l1137
				l1141:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleCase]() {
						goto l1142
					}
					goto l1137
				l1142:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleRowMeta]() {
						goto l1143
					}
					goto l1137
				l1143:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleFuncTypeCast]() {
						goto l1144
					}
					goto l1137
				l1144:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleFuncApp]() {
						goto l1145
					}
					goto l1137
				l1145:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleRowValue]() {
						goto l1146
					}
					goto l1137
				l1146:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleArrayExpr]() {
						goto l1147
					}
					goto l1137
				l1147:
					position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
					if !_rules[ruleLiteral]() {
						goto l1135
					}
				}
			l1137:
				depth--
				add(rulebaseExpr, position1136)
			}
			return true
		l1135:
			position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
			return false
		},
		/* 85 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') spOpt '(' spOpt Expression sp (('a' / 'A') ('s' / 'S')) sp Type spOpt ')')> Action65)> */
		func() bool {
			position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
			{
				position1149 := position
				depth++
				{
					position1150 := position
					depth++
					{
						position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1152
						}
						position++
						goto l1151
					l1152:
						position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
						if buffer[position] != rune('C') {
							goto l1148
						}
						position++
					}
				l1151:
					{
						position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1154
						}
						position++
						goto l1153
					l1154:
						position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
						if buffer[position] != rune('A') {
							goto l1148
						}
						position++
					}
				l1153:
					{
						position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1156
						}
						position++
						goto l1155
					l1156:
						position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
						if buffer[position] != rune('S') {
							goto l1148
						}
						position++
					}
				l1155:
					{
						position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1158
						}
						position++
						goto l1157
					l1158:
						position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
						if buffer[position] != rune('T') {
							goto l1148
						}
						position++
					}
				l1157:
					if !_rules[rulespOpt]() {
						goto l1148
					}
					if buffer[position] != rune('(') {
						goto l1148
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1148
					}
					if !_rules[ruleExpression]() {
						goto l1148
					}
					if !_rules[rulesp]() {
						goto l1148
					}
					{
						position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1160
						}
						position++
						goto l1159
					l1160:
						position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
						if buffer[position] != rune('A') {
							goto l1148
						}
						position++
					}
				l1159:
					{
						position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1162
						}
						position++
						goto l1161
					l1162:
						position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
						if buffer[position] != rune('S') {
							goto l1148
						}
						position++
					}
				l1161:
					if !_rules[rulesp]() {
						goto l1148
					}
					if !_rules[ruleType]() {
						goto l1148
					}
					if !_rules[rulespOpt]() {
						goto l1148
					}
					if buffer[position] != rune(')') {
						goto l1148
					}
					position++
					depth--
					add(rulePegText, position1150)
				}
				if !_rules[ruleAction65]() {
					goto l1148
				}
				depth--
				add(ruleFuncTypeCast, position1149)
			}
			return true
		l1148:
			position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
			return false
		},
		/* 86 FuncApp <- <(FuncAppWithOrderBy / FuncAppWithoutOrderBy)> */
		func() bool {
			position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
			{
				position1164 := position
				depth++
				{
					position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
					if !_rules[ruleFuncAppWithOrderBy]() {
						goto l1166
					}
					goto l1165
				l1166:
					position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
					if !_rules[ruleFuncAppWithoutOrderBy]() {
						goto l1163
					}
				}
			l1165:
				depth--
				add(ruleFuncApp, position1164)
			}
			return true
		l1163:
			position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
			return false
		},
		/* 87 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action66)> */
		func() bool {
			position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
			{
				position1168 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1167
				}
				if !_rules[rulespOpt]() {
					goto l1167
				}
				if buffer[position] != rune('(') {
					goto l1167
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1167
				}
				if !_rules[ruleFuncParams]() {
					goto l1167
				}
				if !_rules[rulesp]() {
					goto l1167
				}
				if !_rules[ruleParamsOrder]() {
					goto l1167
				}
				if !_rules[rulespOpt]() {
					goto l1167
				}
				if buffer[position] != rune(')') {
					goto l1167
				}
				position++
				if !_rules[ruleAction66]() {
					goto l1167
				}
				depth--
				add(ruleFuncAppWithOrderBy, position1168)
			}
			return true
		l1167:
			position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
			return false
		},
		/* 88 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action67)> */
		func() bool {
			position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
			{
				position1170 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1169
				}
				if !_rules[rulespOpt]() {
					goto l1169
				}
				if buffer[position] != rune('(') {
					goto l1169
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1169
				}
				if !_rules[ruleFuncParams]() {
					goto l1169
				}
				{
					position1171 := position
					depth++
					if !_rules[rulespOpt]() {
						goto l1169
					}
					depth--
					add(rulePegText, position1171)
				}
				if buffer[position] != rune(')') {
					goto l1169
				}
				position++
				if !_rules[ruleAction67]() {
					goto l1169
				}
				depth--
				add(ruleFuncAppWithoutOrderBy, position1170)
			}
			return true
		l1169:
			position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
			return false
		},
		/* 89 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action68)> */
		func() bool {
			position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
			{
				position1173 := position
				depth++
				{
					position1174 := position
					depth++
					{
						position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1175
						}
					l1177:
						{
							position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1178
							}
							if buffer[position] != rune(',') {
								goto l1178
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1178
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1178
							}
							goto l1177
						l1178:
							position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						}
						goto l1176
					l1175:
						position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
					}
				l1176:
					depth--
					add(rulePegText, position1174)
				}
				if !_rules[ruleAction68]() {
					goto l1172
				}
				depth--
				add(ruleFuncParams, position1173)
			}
			return true
		l1172:
			position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
			return false
		},
		/* 90 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action69)> */
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
						if buffer[position] != rune('o') {
							goto l1183
						}
						position++
						goto l1182
					l1183:
						position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
						if buffer[position] != rune('O') {
							goto l1179
						}
						position++
					}
				l1182:
					{
						position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1185
						}
						position++
						goto l1184
					l1185:
						position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
						if buffer[position] != rune('R') {
							goto l1179
						}
						position++
					}
				l1184:
					{
						position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1187
						}
						position++
						goto l1186
					l1187:
						position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
						if buffer[position] != rune('D') {
							goto l1179
						}
						position++
					}
				l1186:
					{
						position1188, tokenIndex1188, depth1188 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1189
						}
						position++
						goto l1188
					l1189:
						position, tokenIndex, depth = position1188, tokenIndex1188, depth1188
						if buffer[position] != rune('E') {
							goto l1179
						}
						position++
					}
				l1188:
					{
						position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1191
						}
						position++
						goto l1190
					l1191:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('R') {
							goto l1179
						}
						position++
					}
				l1190:
					if !_rules[rulesp]() {
						goto l1179
					}
					{
						position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1193
						}
						position++
						goto l1192
					l1193:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						if buffer[position] != rune('B') {
							goto l1179
						}
						position++
					}
				l1192:
					{
						position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1195
						}
						position++
						goto l1194
					l1195:
						position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
						if buffer[position] != rune('Y') {
							goto l1179
						}
						position++
					}
				l1194:
					if !_rules[rulesp]() {
						goto l1179
					}
					if !_rules[ruleSortedExpression]() {
						goto l1179
					}
				l1196:
					{
						position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1197
						}
						if buffer[position] != rune(',') {
							goto l1197
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1197
						}
						if !_rules[ruleSortedExpression]() {
							goto l1197
						}
						goto l1196
					l1197:
						position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
					}
					depth--
					add(rulePegText, position1181)
				}
				if !_rules[ruleAction69]() {
					goto l1179
				}
				depth--
				add(ruleParamsOrder, position1180)
			}
			return true
		l1179:
			position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
			return false
		},
		/* 91 SortedExpression <- <(Expression OrderDirectionOpt Action70)> */
		func() bool {
			position1198, tokenIndex1198, depth1198 := position, tokenIndex, depth
			{
				position1199 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l1198
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1198
				}
				if !_rules[ruleAction70]() {
					goto l1198
				}
				depth--
				add(ruleSortedExpression, position1199)
			}
			return true
		l1198:
			position, tokenIndex, depth = position1198, tokenIndex1198, depth1198
			return false
		},
		/* 92 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action71)> */
		func() bool {
			position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
			{
				position1201 := position
				depth++
				{
					position1202 := position
					depth++
					{
						position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1203
						}
						{
							position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
							if !_rules[ruleAscending]() {
								goto l1206
							}
							goto l1205
						l1206:
							position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
							if !_rules[ruleDescending]() {
								goto l1203
							}
						}
					l1205:
						goto l1204
					l1203:
						position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
					}
				l1204:
					depth--
					add(rulePegText, position1202)
				}
				if !_rules[ruleAction71]() {
					goto l1200
				}
				depth--
				add(ruleOrderDirectionOpt, position1201)
			}
			return true
		l1200:
			position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
			return false
		},
		/* 93 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action72)> */
		func() bool {
			position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
			{
				position1208 := position
				depth++
				{
					position1209 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1207
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1207
					}
					{
						position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1210
						}
					l1212:
						{
							position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1213
							}
							if buffer[position] != rune(',') {
								goto l1213
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1213
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1213
							}
							goto l1212
						l1213:
							position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
						}
						goto l1211
					l1210:
						position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
					}
				l1211:
					if !_rules[rulespOpt]() {
						goto l1207
					}
					{
						position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1214
						}
						position++
						goto l1215
					l1214:
						position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
					}
				l1215:
					if !_rules[rulespOpt]() {
						goto l1207
					}
					if buffer[position] != rune(']') {
						goto l1207
					}
					position++
					depth--
					add(rulePegText, position1209)
				}
				if !_rules[ruleAction72]() {
					goto l1207
				}
				depth--
				add(ruleArrayExpr, position1208)
			}
			return true
		l1207:
			position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
			return false
		},
		/* 94 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action73)> */
		func() bool {
			position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
			{
				position1217 := position
				depth++
				{
					position1218 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1216
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1216
					}
					{
						position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l1219
						}
					l1221:
						{
							position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1222
							}
							if buffer[position] != rune(',') {
								goto l1222
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1222
							}
							if !_rules[ruleKeyValuePair]() {
								goto l1222
							}
							goto l1221
						l1222:
							position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
						}
						goto l1220
					l1219:
						position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
					}
				l1220:
					if !_rules[rulespOpt]() {
						goto l1216
					}
					if buffer[position] != rune('}') {
						goto l1216
					}
					position++
					depth--
					add(rulePegText, position1218)
				}
				if !_rules[ruleAction73]() {
					goto l1216
				}
				depth--
				add(ruleMapExpr, position1217)
			}
			return true
		l1216:
			position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
			return false
		},
		/* 95 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action74)> */
		func() bool {
			position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
			{
				position1224 := position
				depth++
				{
					position1225 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1223
					}
					if !_rules[rulespOpt]() {
						goto l1223
					}
					if buffer[position] != rune(':') {
						goto l1223
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1223
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1223
					}
					depth--
					add(rulePegText, position1225)
				}
				if !_rules[ruleAction74]() {
					goto l1223
				}
				depth--
				add(ruleKeyValuePair, position1224)
			}
			return true
		l1223:
			position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
			return false
		},
		/* 96 Case <- <(ConditionCase / ExpressionCase)> */
		func() bool {
			position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
			{
				position1227 := position
				depth++
				{
					position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
					if !_rules[ruleConditionCase]() {
						goto l1229
					}
					goto l1228
				l1229:
					position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
					if !_rules[ruleExpressionCase]() {
						goto l1226
					}
				}
			l1228:
				depth--
				add(ruleCase, position1227)
			}
			return true
		l1226:
			position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
			return false
		},
		/* 97 ConditionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action75)> */
		func() bool {
			position1230, tokenIndex1230, depth1230 := position, tokenIndex, depth
			{
				position1231 := position
				depth++
				{
					position1232, tokenIndex1232, depth1232 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1233
					}
					position++
					goto l1232
				l1233:
					position, tokenIndex, depth = position1232, tokenIndex1232, depth1232
					if buffer[position] != rune('C') {
						goto l1230
					}
					position++
				}
			l1232:
				{
					position1234, tokenIndex1234, depth1234 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1235
					}
					position++
					goto l1234
				l1235:
					position, tokenIndex, depth = position1234, tokenIndex1234, depth1234
					if buffer[position] != rune('A') {
						goto l1230
					}
					position++
				}
			l1234:
				{
					position1236, tokenIndex1236, depth1236 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1237
					}
					position++
					goto l1236
				l1237:
					position, tokenIndex, depth = position1236, tokenIndex1236, depth1236
					if buffer[position] != rune('S') {
						goto l1230
					}
					position++
				}
			l1236:
				{
					position1238, tokenIndex1238, depth1238 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1239
					}
					position++
					goto l1238
				l1239:
					position, tokenIndex, depth = position1238, tokenIndex1238, depth1238
					if buffer[position] != rune('E') {
						goto l1230
					}
					position++
				}
			l1238:
				{
					position1240 := position
					depth++
					if !_rules[rulesp]() {
						goto l1230
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1230
					}
				l1241:
					{
						position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1242
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1242
						}
						goto l1241
					l1242:
						position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
					}
					{
						position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1243
						}
						{
							position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1246
							}
							position++
							goto l1245
						l1246:
							position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
							if buffer[position] != rune('E') {
								goto l1243
							}
							position++
						}
					l1245:
						{
							position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1248
							}
							position++
							goto l1247
						l1248:
							position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
							if buffer[position] != rune('L') {
								goto l1243
							}
							position++
						}
					l1247:
						{
							position1249, tokenIndex1249, depth1249 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1250
							}
							position++
							goto l1249
						l1250:
							position, tokenIndex, depth = position1249, tokenIndex1249, depth1249
							if buffer[position] != rune('S') {
								goto l1243
							}
							position++
						}
					l1249:
						{
							position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1252
							}
							position++
							goto l1251
						l1252:
							position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
							if buffer[position] != rune('E') {
								goto l1243
							}
							position++
						}
					l1251:
						if !_rules[rulesp]() {
							goto l1243
						}
						if !_rules[ruleExpression]() {
							goto l1243
						}
						goto l1244
					l1243:
						position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
					}
				l1244:
					if !_rules[rulesp]() {
						goto l1230
					}
					{
						position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1254
						}
						position++
						goto l1253
					l1254:
						position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
						if buffer[position] != rune('E') {
							goto l1230
						}
						position++
					}
				l1253:
					{
						position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1256
						}
						position++
						goto l1255
					l1256:
						position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
						if buffer[position] != rune('N') {
							goto l1230
						}
						position++
					}
				l1255:
					{
						position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1258
						}
						position++
						goto l1257
					l1258:
						position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
						if buffer[position] != rune('D') {
							goto l1230
						}
						position++
					}
				l1257:
					depth--
					add(rulePegText, position1240)
				}
				if !_rules[ruleAction75]() {
					goto l1230
				}
				depth--
				add(ruleConditionCase, position1231)
			}
			return true
		l1230:
			position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
			return false
		},
		/* 98 ExpressionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') sp Expression <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action76)> */
		func() bool {
			position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
			{
				position1260 := position
				depth++
				{
					position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1262
					}
					position++
					goto l1261
				l1262:
					position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
					if buffer[position] != rune('C') {
						goto l1259
					}
					position++
				}
			l1261:
				{
					position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1264
					}
					position++
					goto l1263
				l1264:
					position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
					if buffer[position] != rune('A') {
						goto l1259
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
						goto l1259
					}
					position++
				}
			l1265:
				{
					position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1268
					}
					position++
					goto l1267
				l1268:
					position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
					if buffer[position] != rune('E') {
						goto l1259
					}
					position++
				}
			l1267:
				if !_rules[rulesp]() {
					goto l1259
				}
				if !_rules[ruleExpression]() {
					goto l1259
				}
				{
					position1269 := position
					depth++
					if !_rules[rulesp]() {
						goto l1259
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1259
					}
				l1270:
					{
						position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1271
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1271
						}
						goto l1270
					l1271:
						position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
					}
					{
						position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1272
						}
						{
							position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1275
							}
							position++
							goto l1274
						l1275:
							position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
							if buffer[position] != rune('E') {
								goto l1272
							}
							position++
						}
					l1274:
						{
							position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1277
							}
							position++
							goto l1276
						l1277:
							position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
							if buffer[position] != rune('L') {
								goto l1272
							}
							position++
						}
					l1276:
						{
							position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1279
							}
							position++
							goto l1278
						l1279:
							position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
							if buffer[position] != rune('S') {
								goto l1272
							}
							position++
						}
					l1278:
						{
							position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1281
							}
							position++
							goto l1280
						l1281:
							position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
							if buffer[position] != rune('E') {
								goto l1272
							}
							position++
						}
					l1280:
						if !_rules[rulesp]() {
							goto l1272
						}
						if !_rules[ruleExpression]() {
							goto l1272
						}
						goto l1273
					l1272:
						position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
					}
				l1273:
					if !_rules[rulesp]() {
						goto l1259
					}
					{
						position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1283
						}
						position++
						goto l1282
					l1283:
						position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
						if buffer[position] != rune('E') {
							goto l1259
						}
						position++
					}
				l1282:
					{
						position1284, tokenIndex1284, depth1284 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1285
						}
						position++
						goto l1284
					l1285:
						position, tokenIndex, depth = position1284, tokenIndex1284, depth1284
						if buffer[position] != rune('N') {
							goto l1259
						}
						position++
					}
				l1284:
					{
						position1286, tokenIndex1286, depth1286 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1287
						}
						position++
						goto l1286
					l1287:
						position, tokenIndex, depth = position1286, tokenIndex1286, depth1286
						if buffer[position] != rune('D') {
							goto l1259
						}
						position++
					}
				l1286:
					depth--
					add(rulePegText, position1269)
				}
				if !_rules[ruleAction76]() {
					goto l1259
				}
				depth--
				add(ruleExpressionCase, position1260)
			}
			return true
		l1259:
			position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
			return false
		},
		/* 99 WhenThenPair <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('n' / 'N') sp Expression sp (('t' / 'T') ('h' / 'H') ('e' / 'E') ('n' / 'N')) sp ExpressionOrWildcard Action77)> */
		func() bool {
			position1288, tokenIndex1288, depth1288 := position, tokenIndex, depth
			{
				position1289 := position
				depth++
				{
					position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l1291
					}
					position++
					goto l1290
				l1291:
					position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
					if buffer[position] != rune('W') {
						goto l1288
					}
					position++
				}
			l1290:
				{
					position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
					if buffer[position] != rune('h') {
						goto l1293
					}
					position++
					goto l1292
				l1293:
					position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
					if buffer[position] != rune('H') {
						goto l1288
					}
					position++
				}
			l1292:
				{
					position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1295
					}
					position++
					goto l1294
				l1295:
					position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
					if buffer[position] != rune('E') {
						goto l1288
					}
					position++
				}
			l1294:
				{
					position1296, tokenIndex1296, depth1296 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1297
					}
					position++
					goto l1296
				l1297:
					position, tokenIndex, depth = position1296, tokenIndex1296, depth1296
					if buffer[position] != rune('N') {
						goto l1288
					}
					position++
				}
			l1296:
				if !_rules[rulesp]() {
					goto l1288
				}
				if !_rules[ruleExpression]() {
					goto l1288
				}
				if !_rules[rulesp]() {
					goto l1288
				}
				{
					position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1299
					}
					position++
					goto l1298
				l1299:
					position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
					if buffer[position] != rune('T') {
						goto l1288
					}
					position++
				}
			l1298:
				{
					position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
					if buffer[position] != rune('h') {
						goto l1301
					}
					position++
					goto l1300
				l1301:
					position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
					if buffer[position] != rune('H') {
						goto l1288
					}
					position++
				}
			l1300:
				{
					position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1303
					}
					position++
					goto l1302
				l1303:
					position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
					if buffer[position] != rune('E') {
						goto l1288
					}
					position++
				}
			l1302:
				{
					position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1305
					}
					position++
					goto l1304
				l1305:
					position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
					if buffer[position] != rune('N') {
						goto l1288
					}
					position++
				}
			l1304:
				if !_rules[rulesp]() {
					goto l1288
				}
				if !_rules[ruleExpressionOrWildcard]() {
					goto l1288
				}
				if !_rules[ruleAction77]() {
					goto l1288
				}
				depth--
				add(ruleWhenThenPair, position1289)
			}
			return true
		l1288:
			position, tokenIndex, depth = position1288, tokenIndex1288, depth1288
			return false
		},
		/* 100 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
			{
				position1307 := position
				depth++
				{
					position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l1309
					}
					goto l1308
				l1309:
					position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
					if !_rules[ruleNumericLiteral]() {
						goto l1310
					}
					goto l1308
				l1310:
					position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
					if !_rules[ruleStringLiteral]() {
						goto l1306
					}
				}
			l1308:
				depth--
				add(ruleLiteral, position1307)
			}
			return true
		l1306:
			position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
			return false
		},
		/* 101 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
			{
				position1312 := position
				depth++
				{
					position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l1314
					}
					goto l1313
				l1314:
					position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
					if !_rules[ruleNotEqual]() {
						goto l1315
					}
					goto l1313
				l1315:
					position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
					if !_rules[ruleLessOrEqual]() {
						goto l1316
					}
					goto l1313
				l1316:
					position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
					if !_rules[ruleLess]() {
						goto l1317
					}
					goto l1313
				l1317:
					position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
					if !_rules[ruleGreaterOrEqual]() {
						goto l1318
					}
					goto l1313
				l1318:
					position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
					if !_rules[ruleGreater]() {
						goto l1319
					}
					goto l1313
				l1319:
					position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
					if !_rules[ruleNotEqual]() {
						goto l1311
					}
				}
			l1313:
				depth--
				add(ruleComparisonOp, position1312)
			}
			return true
		l1311:
			position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
			return false
		},
		/* 102 OtherOp <- <Concat> */
		func() bool {
			position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
			{
				position1321 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l1320
				}
				depth--
				add(ruleOtherOp, position1321)
			}
			return true
		l1320:
			position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
			return false
		},
		/* 103 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
			{
				position1323 := position
				depth++
				{
					position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l1325
					}
					goto l1324
				l1325:
					position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
					if !_rules[ruleIs]() {
						goto l1322
					}
				}
			l1324:
				depth--
				add(ruleIsOp, position1323)
			}
			return true
		l1322:
			position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
			return false
		},
		/* 104 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
			{
				position1327 := position
				depth++
				{
					position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l1329
					}
					goto l1328
				l1329:
					position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
					if !_rules[ruleMinus]() {
						goto l1326
					}
				}
			l1328:
				depth--
				add(rulePlusMinusOp, position1327)
			}
			return true
		l1326:
			position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
			return false
		},
		/* 105 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1330, tokenIndex1330, depth1330 := position, tokenIndex, depth
			{
				position1331 := position
				depth++
				{
					position1332, tokenIndex1332, depth1332 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l1333
					}
					goto l1332
				l1333:
					position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
					if !_rules[ruleDivide]() {
						goto l1334
					}
					goto l1332
				l1334:
					position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
					if !_rules[ruleModulo]() {
						goto l1330
					}
				}
			l1332:
				depth--
				add(ruleMultDivOp, position1331)
			}
			return true
		l1330:
			position, tokenIndex, depth = position1330, tokenIndex1330, depth1330
			return false
		},
		/* 106 Stream <- <(<ident> Action78)> */
		func() bool {
			position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
			{
				position1336 := position
				depth++
				{
					position1337 := position
					depth++
					if !_rules[ruleident]() {
						goto l1335
					}
					depth--
					add(rulePegText, position1337)
				}
				if !_rules[ruleAction78]() {
					goto l1335
				}
				depth--
				add(ruleStream, position1336)
			}
			return true
		l1335:
			position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
			return false
		},
		/* 107 RowMeta <- <RowTimestamp> */
		func() bool {
			position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
			{
				position1339 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l1338
				}
				depth--
				add(ruleRowMeta, position1339)
			}
			return true
		l1338:
			position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
			return false
		},
		/* 108 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action79)> */
		func() bool {
			position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
			{
				position1341 := position
				depth++
				{
					position1342 := position
					depth++
					{
						position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1343
						}
						if buffer[position] != rune(':') {
							goto l1343
						}
						position++
						goto l1344
					l1343:
						position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					}
				l1344:
					if buffer[position] != rune('t') {
						goto l1340
					}
					position++
					if buffer[position] != rune('s') {
						goto l1340
					}
					position++
					if buffer[position] != rune('(') {
						goto l1340
					}
					position++
					if buffer[position] != rune(')') {
						goto l1340
					}
					position++
					depth--
					add(rulePegText, position1342)
				}
				if !_rules[ruleAction79]() {
					goto l1340
				}
				depth--
				add(ruleRowTimestamp, position1341)
			}
			return true
		l1340:
			position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
			return false
		},
		/* 109 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action80)> */
		func() bool {
			position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
			{
				position1346 := position
				depth++
				{
					position1347 := position
					depth++
					{
						position1348, tokenIndex1348, depth1348 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1348
						}
						if buffer[position] != rune(':') {
							goto l1348
						}
						position++
						{
							position1350, tokenIndex1350, depth1350 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1350
							}
							position++
							goto l1348
						l1350:
							position, tokenIndex, depth = position1350, tokenIndex1350, depth1350
						}
						goto l1349
					l1348:
						position, tokenIndex, depth = position1348, tokenIndex1348, depth1348
					}
				l1349:
					if !_rules[rulejsonGetPath]() {
						goto l1345
					}
					depth--
					add(rulePegText, position1347)
				}
				if !_rules[ruleAction80]() {
					goto l1345
				}
				depth--
				add(ruleRowValue, position1346)
			}
			return true
		l1345:
			position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
			return false
		},
		/* 110 NumericLiteral <- <(<('-'? [0-9]+)> Action81)> */
		func() bool {
			position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
			{
				position1352 := position
				depth++
				{
					position1353 := position
					depth++
					{
						position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1354
						}
						position++
						goto l1355
					l1354:
						position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
					}
				l1355:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1351
					}
					position++
				l1356:
					{
						position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1357
						}
						position++
						goto l1356
					l1357:
						position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
					}
					depth--
					add(rulePegText, position1353)
				}
				if !_rules[ruleAction81]() {
					goto l1351
				}
				depth--
				add(ruleNumericLiteral, position1352)
			}
			return true
		l1351:
			position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
			return false
		},
		/* 111 NonNegativeNumericLiteral <- <(<[0-9]+> Action82)> */
		func() bool {
			position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
			{
				position1359 := position
				depth++
				{
					position1360 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1358
					}
					position++
				l1361:
					{
						position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1362
						}
						position++
						goto l1361
					l1362:
						position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
					}
					depth--
					add(rulePegText, position1360)
				}
				if !_rules[ruleAction82]() {
					goto l1358
				}
				depth--
				add(ruleNonNegativeNumericLiteral, position1359)
			}
			return true
		l1358:
			position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
			return false
		},
		/* 112 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action83)> */
		func() bool {
			position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
			{
				position1364 := position
				depth++
				{
					position1365 := position
					depth++
					{
						position1366, tokenIndex1366, depth1366 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1366
						}
						position++
						goto l1367
					l1366:
						position, tokenIndex, depth = position1366, tokenIndex1366, depth1366
					}
				l1367:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1363
					}
					position++
				l1368:
					{
						position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1369
						}
						position++
						goto l1368
					l1369:
						position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
					}
					if buffer[position] != rune('.') {
						goto l1363
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1363
					}
					position++
				l1370:
					{
						position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1371
						}
						position++
						goto l1370
					l1371:
						position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
					}
					depth--
					add(rulePegText, position1365)
				}
				if !_rules[ruleAction83]() {
					goto l1363
				}
				depth--
				add(ruleFloatLiteral, position1364)
			}
			return true
		l1363:
			position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
			return false
		},
		/* 113 Function <- <(<ident> Action84)> */
		func() bool {
			position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
			{
				position1373 := position
				depth++
				{
					position1374 := position
					depth++
					if !_rules[ruleident]() {
						goto l1372
					}
					depth--
					add(rulePegText, position1374)
				}
				if !_rules[ruleAction84]() {
					goto l1372
				}
				depth--
				add(ruleFunction, position1373)
			}
			return true
		l1372:
			position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
			return false
		},
		/* 114 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action85)> */
		func() bool {
			position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
			{
				position1376 := position
				depth++
				{
					position1377 := position
					depth++
					{
						position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1379
						}
						position++
						goto l1378
					l1379:
						position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
						if buffer[position] != rune('N') {
							goto l1375
						}
						position++
					}
				l1378:
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
							goto l1375
						}
						position++
					}
				l1380:
					{
						position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1383
						}
						position++
						goto l1382
					l1383:
						position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
						if buffer[position] != rune('L') {
							goto l1375
						}
						position++
					}
				l1382:
					{
						position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1385
						}
						position++
						goto l1384
					l1385:
						position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
						if buffer[position] != rune('L') {
							goto l1375
						}
						position++
					}
				l1384:
					depth--
					add(rulePegText, position1377)
				}
				if !_rules[ruleAction85]() {
					goto l1375
				}
				depth--
				add(ruleNullLiteral, position1376)
			}
			return true
		l1375:
			position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
			return false
		},
		/* 115 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
			{
				position1387 := position
				depth++
				{
					position1388, tokenIndex1388, depth1388 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l1389
					}
					goto l1388
				l1389:
					position, tokenIndex, depth = position1388, tokenIndex1388, depth1388
					if !_rules[ruleFALSE]() {
						goto l1386
					}
				}
			l1388:
				depth--
				add(ruleBooleanLiteral, position1387)
			}
			return true
		l1386:
			position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
			return false
		},
		/* 116 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action86)> */
		func() bool {
			position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
			{
				position1391 := position
				depth++
				{
					position1392 := position
					depth++
					{
						position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1394
						}
						position++
						goto l1393
					l1394:
						position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
						if buffer[position] != rune('T') {
							goto l1390
						}
						position++
					}
				l1393:
					{
						position1395, tokenIndex1395, depth1395 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1396
						}
						position++
						goto l1395
					l1396:
						position, tokenIndex, depth = position1395, tokenIndex1395, depth1395
						if buffer[position] != rune('R') {
							goto l1390
						}
						position++
					}
				l1395:
					{
						position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1398
						}
						position++
						goto l1397
					l1398:
						position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
						if buffer[position] != rune('U') {
							goto l1390
						}
						position++
					}
				l1397:
					{
						position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1400
						}
						position++
						goto l1399
					l1400:
						position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
						if buffer[position] != rune('E') {
							goto l1390
						}
						position++
					}
				l1399:
					depth--
					add(rulePegText, position1392)
				}
				if !_rules[ruleAction86]() {
					goto l1390
				}
				depth--
				add(ruleTRUE, position1391)
			}
			return true
		l1390:
			position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
			return false
		},
		/* 117 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action87)> */
		func() bool {
			position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
			{
				position1402 := position
				depth++
				{
					position1403 := position
					depth++
					{
						position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1405
						}
						position++
						goto l1404
					l1405:
						position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
						if buffer[position] != rune('F') {
							goto l1401
						}
						position++
					}
				l1404:
					{
						position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1407
						}
						position++
						goto l1406
					l1407:
						position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
						if buffer[position] != rune('A') {
							goto l1401
						}
						position++
					}
				l1406:
					{
						position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1409
						}
						position++
						goto l1408
					l1409:
						position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
						if buffer[position] != rune('L') {
							goto l1401
						}
						position++
					}
				l1408:
					{
						position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1411
						}
						position++
						goto l1410
					l1411:
						position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
						if buffer[position] != rune('S') {
							goto l1401
						}
						position++
					}
				l1410:
					{
						position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1413
						}
						position++
						goto l1412
					l1413:
						position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
						if buffer[position] != rune('E') {
							goto l1401
						}
						position++
					}
				l1412:
					depth--
					add(rulePegText, position1403)
				}
				if !_rules[ruleAction87]() {
					goto l1401
				}
				depth--
				add(ruleFALSE, position1402)
			}
			return true
		l1401:
			position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
			return false
		},
		/* 118 Wildcard <- <(<((ident ':' !':')? '*')> Action88)> */
		func() bool {
			position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
			{
				position1415 := position
				depth++
				{
					position1416 := position
					depth++
					{
						position1417, tokenIndex1417, depth1417 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1417
						}
						if buffer[position] != rune(':') {
							goto l1417
						}
						position++
						{
							position1419, tokenIndex1419, depth1419 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1419
							}
							position++
							goto l1417
						l1419:
							position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
						}
						goto l1418
					l1417:
						position, tokenIndex, depth = position1417, tokenIndex1417, depth1417
					}
				l1418:
					if buffer[position] != rune('*') {
						goto l1414
					}
					position++
					depth--
					add(rulePegText, position1416)
				}
				if !_rules[ruleAction88]() {
					goto l1414
				}
				depth--
				add(ruleWildcard, position1415)
			}
			return true
		l1414:
			position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
			return false
		},
		/* 119 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action89)> */
		func() bool {
			position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
			{
				position1421 := position
				depth++
				{
					position1422 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l1420
					}
					position++
				l1423:
					{
						position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
						{
							position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1426
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1426
							}
							position++
							goto l1425
						l1426:
							position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
							{
								position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1427
								}
								position++
								goto l1424
							l1427:
								position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
							}
							if !matchDot() {
								goto l1424
							}
						}
					l1425:
						goto l1423
					l1424:
						position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
					}
					if buffer[position] != rune('\'') {
						goto l1420
					}
					position++
					depth--
					add(rulePegText, position1422)
				}
				if !_rules[ruleAction89]() {
					goto l1420
				}
				depth--
				add(ruleStringLiteral, position1421)
			}
			return true
		l1420:
			position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
			return false
		},
		/* 120 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action90)> */
		func() bool {
			position1428, tokenIndex1428, depth1428 := position, tokenIndex, depth
			{
				position1429 := position
				depth++
				{
					position1430 := position
					depth++
					{
						position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1432
						}
						position++
						goto l1431
					l1432:
						position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
						if buffer[position] != rune('I') {
							goto l1428
						}
						position++
					}
				l1431:
					{
						position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1434
						}
						position++
						goto l1433
					l1434:
						position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
						if buffer[position] != rune('S') {
							goto l1428
						}
						position++
					}
				l1433:
					{
						position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1436
						}
						position++
						goto l1435
					l1436:
						position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
						if buffer[position] != rune('T') {
							goto l1428
						}
						position++
					}
				l1435:
					{
						position1437, tokenIndex1437, depth1437 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1438
						}
						position++
						goto l1437
					l1438:
						position, tokenIndex, depth = position1437, tokenIndex1437, depth1437
						if buffer[position] != rune('R') {
							goto l1428
						}
						position++
					}
				l1437:
					{
						position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1440
						}
						position++
						goto l1439
					l1440:
						position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
						if buffer[position] != rune('E') {
							goto l1428
						}
						position++
					}
				l1439:
					{
						position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1442
						}
						position++
						goto l1441
					l1442:
						position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
						if buffer[position] != rune('A') {
							goto l1428
						}
						position++
					}
				l1441:
					{
						position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1444
						}
						position++
						goto l1443
					l1444:
						position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
						if buffer[position] != rune('M') {
							goto l1428
						}
						position++
					}
				l1443:
					depth--
					add(rulePegText, position1430)
				}
				if !_rules[ruleAction90]() {
					goto l1428
				}
				depth--
				add(ruleISTREAM, position1429)
			}
			return true
		l1428:
			position, tokenIndex, depth = position1428, tokenIndex1428, depth1428
			return false
		},
		/* 121 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action91)> */
		func() bool {
			position1445, tokenIndex1445, depth1445 := position, tokenIndex, depth
			{
				position1446 := position
				depth++
				{
					position1447 := position
					depth++
					{
						position1448, tokenIndex1448, depth1448 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1449
						}
						position++
						goto l1448
					l1449:
						position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
						if buffer[position] != rune('D') {
							goto l1445
						}
						position++
					}
				l1448:
					{
						position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1451
						}
						position++
						goto l1450
					l1451:
						position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
						if buffer[position] != rune('S') {
							goto l1445
						}
						position++
					}
				l1450:
					{
						position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1453
						}
						position++
						goto l1452
					l1453:
						position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
						if buffer[position] != rune('T') {
							goto l1445
						}
						position++
					}
				l1452:
					{
						position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1455
						}
						position++
						goto l1454
					l1455:
						position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
						if buffer[position] != rune('R') {
							goto l1445
						}
						position++
					}
				l1454:
					{
						position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1457
						}
						position++
						goto l1456
					l1457:
						position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
						if buffer[position] != rune('E') {
							goto l1445
						}
						position++
					}
				l1456:
					{
						position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1459
						}
						position++
						goto l1458
					l1459:
						position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
						if buffer[position] != rune('A') {
							goto l1445
						}
						position++
					}
				l1458:
					{
						position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1461
						}
						position++
						goto l1460
					l1461:
						position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
						if buffer[position] != rune('M') {
							goto l1445
						}
						position++
					}
				l1460:
					depth--
					add(rulePegText, position1447)
				}
				if !_rules[ruleAction91]() {
					goto l1445
				}
				depth--
				add(ruleDSTREAM, position1446)
			}
			return true
		l1445:
			position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
			return false
		},
		/* 122 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action92)> */
		func() bool {
			position1462, tokenIndex1462, depth1462 := position, tokenIndex, depth
			{
				position1463 := position
				depth++
				{
					position1464 := position
					depth++
					{
						position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1466
						}
						position++
						goto l1465
					l1466:
						position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
						if buffer[position] != rune('R') {
							goto l1462
						}
						position++
					}
				l1465:
					{
						position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1468
						}
						position++
						goto l1467
					l1468:
						position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
						if buffer[position] != rune('S') {
							goto l1462
						}
						position++
					}
				l1467:
					{
						position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1470
						}
						position++
						goto l1469
					l1470:
						position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
						if buffer[position] != rune('T') {
							goto l1462
						}
						position++
					}
				l1469:
					{
						position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1472
						}
						position++
						goto l1471
					l1472:
						position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
						if buffer[position] != rune('R') {
							goto l1462
						}
						position++
					}
				l1471:
					{
						position1473, tokenIndex1473, depth1473 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1474
						}
						position++
						goto l1473
					l1474:
						position, tokenIndex, depth = position1473, tokenIndex1473, depth1473
						if buffer[position] != rune('E') {
							goto l1462
						}
						position++
					}
				l1473:
					{
						position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1476
						}
						position++
						goto l1475
					l1476:
						position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
						if buffer[position] != rune('A') {
							goto l1462
						}
						position++
					}
				l1475:
					{
						position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1478
						}
						position++
						goto l1477
					l1478:
						position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
						if buffer[position] != rune('M') {
							goto l1462
						}
						position++
					}
				l1477:
					depth--
					add(rulePegText, position1464)
				}
				if !_rules[ruleAction92]() {
					goto l1462
				}
				depth--
				add(ruleRSTREAM, position1463)
			}
			return true
		l1462:
			position, tokenIndex, depth = position1462, tokenIndex1462, depth1462
			return false
		},
		/* 123 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action93)> */
		func() bool {
			position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
			{
				position1480 := position
				depth++
				{
					position1481 := position
					depth++
					{
						position1482, tokenIndex1482, depth1482 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1483
						}
						position++
						goto l1482
					l1483:
						position, tokenIndex, depth = position1482, tokenIndex1482, depth1482
						if buffer[position] != rune('T') {
							goto l1479
						}
						position++
					}
				l1482:
					{
						position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1485
						}
						position++
						goto l1484
					l1485:
						position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
						if buffer[position] != rune('U') {
							goto l1479
						}
						position++
					}
				l1484:
					{
						position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1487
						}
						position++
						goto l1486
					l1487:
						position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
						if buffer[position] != rune('P') {
							goto l1479
						}
						position++
					}
				l1486:
					{
						position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1489
						}
						position++
						goto l1488
					l1489:
						position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
						if buffer[position] != rune('L') {
							goto l1479
						}
						position++
					}
				l1488:
					{
						position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1491
						}
						position++
						goto l1490
					l1491:
						position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
						if buffer[position] != rune('E') {
							goto l1479
						}
						position++
					}
				l1490:
					{
						position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1493
						}
						position++
						goto l1492
					l1493:
						position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
						if buffer[position] != rune('S') {
							goto l1479
						}
						position++
					}
				l1492:
					depth--
					add(rulePegText, position1481)
				}
				if !_rules[ruleAction93]() {
					goto l1479
				}
				depth--
				add(ruleTUPLES, position1480)
			}
			return true
		l1479:
			position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
			return false
		},
		/* 124 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action94)> */
		func() bool {
			position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
			{
				position1495 := position
				depth++
				{
					position1496 := position
					depth++
					{
						position1497, tokenIndex1497, depth1497 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1498
						}
						position++
						goto l1497
					l1498:
						position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
						if buffer[position] != rune('S') {
							goto l1494
						}
						position++
					}
				l1497:
					{
						position1499, tokenIndex1499, depth1499 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1500
						}
						position++
						goto l1499
					l1500:
						position, tokenIndex, depth = position1499, tokenIndex1499, depth1499
						if buffer[position] != rune('E') {
							goto l1494
						}
						position++
					}
				l1499:
					{
						position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1502
						}
						position++
						goto l1501
					l1502:
						position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
						if buffer[position] != rune('C') {
							goto l1494
						}
						position++
					}
				l1501:
					{
						position1503, tokenIndex1503, depth1503 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1504
						}
						position++
						goto l1503
					l1504:
						position, tokenIndex, depth = position1503, tokenIndex1503, depth1503
						if buffer[position] != rune('O') {
							goto l1494
						}
						position++
					}
				l1503:
					{
						position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1506
						}
						position++
						goto l1505
					l1506:
						position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
						if buffer[position] != rune('N') {
							goto l1494
						}
						position++
					}
				l1505:
					{
						position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1508
						}
						position++
						goto l1507
					l1508:
						position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
						if buffer[position] != rune('D') {
							goto l1494
						}
						position++
					}
				l1507:
					{
						position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1510
						}
						position++
						goto l1509
					l1510:
						position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
						if buffer[position] != rune('S') {
							goto l1494
						}
						position++
					}
				l1509:
					depth--
					add(rulePegText, position1496)
				}
				if !_rules[ruleAction94]() {
					goto l1494
				}
				depth--
				add(ruleSECONDS, position1495)
			}
			return true
		l1494:
			position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
			return false
		},
		/* 125 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action95)> */
		func() bool {
			position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
			{
				position1512 := position
				depth++
				{
					position1513 := position
					depth++
					{
						position1514, tokenIndex1514, depth1514 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1515
						}
						position++
						goto l1514
					l1515:
						position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
						if buffer[position] != rune('M') {
							goto l1511
						}
						position++
					}
				l1514:
					{
						position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1517
						}
						position++
						goto l1516
					l1517:
						position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
						if buffer[position] != rune('I') {
							goto l1511
						}
						position++
					}
				l1516:
					{
						position1518, tokenIndex1518, depth1518 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1519
						}
						position++
						goto l1518
					l1519:
						position, tokenIndex, depth = position1518, tokenIndex1518, depth1518
						if buffer[position] != rune('L') {
							goto l1511
						}
						position++
					}
				l1518:
					{
						position1520, tokenIndex1520, depth1520 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1521
						}
						position++
						goto l1520
					l1521:
						position, tokenIndex, depth = position1520, tokenIndex1520, depth1520
						if buffer[position] != rune('L') {
							goto l1511
						}
						position++
					}
				l1520:
					{
						position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1523
						}
						position++
						goto l1522
					l1523:
						position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
						if buffer[position] != rune('I') {
							goto l1511
						}
						position++
					}
				l1522:
					{
						position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1525
						}
						position++
						goto l1524
					l1525:
						position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
						if buffer[position] != rune('S') {
							goto l1511
						}
						position++
					}
				l1524:
					{
						position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1527
						}
						position++
						goto l1526
					l1527:
						position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
						if buffer[position] != rune('E') {
							goto l1511
						}
						position++
					}
				l1526:
					{
						position1528, tokenIndex1528, depth1528 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1529
						}
						position++
						goto l1528
					l1529:
						position, tokenIndex, depth = position1528, tokenIndex1528, depth1528
						if buffer[position] != rune('C') {
							goto l1511
						}
						position++
					}
				l1528:
					{
						position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1531
						}
						position++
						goto l1530
					l1531:
						position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
						if buffer[position] != rune('O') {
							goto l1511
						}
						position++
					}
				l1530:
					{
						position1532, tokenIndex1532, depth1532 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1533
						}
						position++
						goto l1532
					l1533:
						position, tokenIndex, depth = position1532, tokenIndex1532, depth1532
						if buffer[position] != rune('N') {
							goto l1511
						}
						position++
					}
				l1532:
					{
						position1534, tokenIndex1534, depth1534 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1535
						}
						position++
						goto l1534
					l1535:
						position, tokenIndex, depth = position1534, tokenIndex1534, depth1534
						if buffer[position] != rune('D') {
							goto l1511
						}
						position++
					}
				l1534:
					{
						position1536, tokenIndex1536, depth1536 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1537
						}
						position++
						goto l1536
					l1537:
						position, tokenIndex, depth = position1536, tokenIndex1536, depth1536
						if buffer[position] != rune('S') {
							goto l1511
						}
						position++
					}
				l1536:
					depth--
					add(rulePegText, position1513)
				}
				if !_rules[ruleAction95]() {
					goto l1511
				}
				depth--
				add(ruleMILLISECONDS, position1512)
			}
			return true
		l1511:
			position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
			return false
		},
		/* 126 Wait <- <(<(('w' / 'W') ('a' / 'A') ('i' / 'I') ('t' / 'T'))> Action96)> */
		func() bool {
			position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
			{
				position1539 := position
				depth++
				{
					position1540 := position
					depth++
					{
						position1541, tokenIndex1541, depth1541 := position, tokenIndex, depth
						if buffer[position] != rune('w') {
							goto l1542
						}
						position++
						goto l1541
					l1542:
						position, tokenIndex, depth = position1541, tokenIndex1541, depth1541
						if buffer[position] != rune('W') {
							goto l1538
						}
						position++
					}
				l1541:
					{
						position1543, tokenIndex1543, depth1543 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1544
						}
						position++
						goto l1543
					l1544:
						position, tokenIndex, depth = position1543, tokenIndex1543, depth1543
						if buffer[position] != rune('A') {
							goto l1538
						}
						position++
					}
				l1543:
					{
						position1545, tokenIndex1545, depth1545 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1546
						}
						position++
						goto l1545
					l1546:
						position, tokenIndex, depth = position1545, tokenIndex1545, depth1545
						if buffer[position] != rune('I') {
							goto l1538
						}
						position++
					}
				l1545:
					{
						position1547, tokenIndex1547, depth1547 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1548
						}
						position++
						goto l1547
					l1548:
						position, tokenIndex, depth = position1547, tokenIndex1547, depth1547
						if buffer[position] != rune('T') {
							goto l1538
						}
						position++
					}
				l1547:
					depth--
					add(rulePegText, position1540)
				}
				if !_rules[ruleAction96]() {
					goto l1538
				}
				depth--
				add(ruleWait, position1539)
			}
			return true
		l1538:
			position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
			return false
		},
		/* 127 DropOldest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('o' / 'O') ('l' / 'L') ('d' / 'D') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action97)> */
		func() bool {
			position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
			{
				position1550 := position
				depth++
				{
					position1551 := position
					depth++
					{
						position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1553
						}
						position++
						goto l1552
					l1553:
						position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
						if buffer[position] != rune('D') {
							goto l1549
						}
						position++
					}
				l1552:
					{
						position1554, tokenIndex1554, depth1554 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1555
						}
						position++
						goto l1554
					l1555:
						position, tokenIndex, depth = position1554, tokenIndex1554, depth1554
						if buffer[position] != rune('R') {
							goto l1549
						}
						position++
					}
				l1554:
					{
						position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1557
						}
						position++
						goto l1556
					l1557:
						position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
						if buffer[position] != rune('O') {
							goto l1549
						}
						position++
					}
				l1556:
					{
						position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1559
						}
						position++
						goto l1558
					l1559:
						position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
						if buffer[position] != rune('P') {
							goto l1549
						}
						position++
					}
				l1558:
					if !_rules[rulesp]() {
						goto l1549
					}
					{
						position1560, tokenIndex1560, depth1560 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1561
						}
						position++
						goto l1560
					l1561:
						position, tokenIndex, depth = position1560, tokenIndex1560, depth1560
						if buffer[position] != rune('O') {
							goto l1549
						}
						position++
					}
				l1560:
					{
						position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1563
						}
						position++
						goto l1562
					l1563:
						position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
						if buffer[position] != rune('L') {
							goto l1549
						}
						position++
					}
				l1562:
					{
						position1564, tokenIndex1564, depth1564 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1565
						}
						position++
						goto l1564
					l1565:
						position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
						if buffer[position] != rune('D') {
							goto l1549
						}
						position++
					}
				l1564:
					{
						position1566, tokenIndex1566, depth1566 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1567
						}
						position++
						goto l1566
					l1567:
						position, tokenIndex, depth = position1566, tokenIndex1566, depth1566
						if buffer[position] != rune('E') {
							goto l1549
						}
						position++
					}
				l1566:
					{
						position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1569
						}
						position++
						goto l1568
					l1569:
						position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
						if buffer[position] != rune('S') {
							goto l1549
						}
						position++
					}
				l1568:
					{
						position1570, tokenIndex1570, depth1570 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1571
						}
						position++
						goto l1570
					l1571:
						position, tokenIndex, depth = position1570, tokenIndex1570, depth1570
						if buffer[position] != rune('T') {
							goto l1549
						}
						position++
					}
				l1570:
					depth--
					add(rulePegText, position1551)
				}
				if !_rules[ruleAction97]() {
					goto l1549
				}
				depth--
				add(ruleDropOldest, position1550)
			}
			return true
		l1549:
			position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
			return false
		},
		/* 128 DropNewest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('n' / 'N') ('e' / 'E') ('w' / 'W') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action98)> */
		func() bool {
			position1572, tokenIndex1572, depth1572 := position, tokenIndex, depth
			{
				position1573 := position
				depth++
				{
					position1574 := position
					depth++
					{
						position1575, tokenIndex1575, depth1575 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1576
						}
						position++
						goto l1575
					l1576:
						position, tokenIndex, depth = position1575, tokenIndex1575, depth1575
						if buffer[position] != rune('D') {
							goto l1572
						}
						position++
					}
				l1575:
					{
						position1577, tokenIndex1577, depth1577 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1578
						}
						position++
						goto l1577
					l1578:
						position, tokenIndex, depth = position1577, tokenIndex1577, depth1577
						if buffer[position] != rune('R') {
							goto l1572
						}
						position++
					}
				l1577:
					{
						position1579, tokenIndex1579, depth1579 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1580
						}
						position++
						goto l1579
					l1580:
						position, tokenIndex, depth = position1579, tokenIndex1579, depth1579
						if buffer[position] != rune('O') {
							goto l1572
						}
						position++
					}
				l1579:
					{
						position1581, tokenIndex1581, depth1581 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1582
						}
						position++
						goto l1581
					l1582:
						position, tokenIndex, depth = position1581, tokenIndex1581, depth1581
						if buffer[position] != rune('P') {
							goto l1572
						}
						position++
					}
				l1581:
					if !_rules[rulesp]() {
						goto l1572
					}
					{
						position1583, tokenIndex1583, depth1583 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1584
						}
						position++
						goto l1583
					l1584:
						position, tokenIndex, depth = position1583, tokenIndex1583, depth1583
						if buffer[position] != rune('N') {
							goto l1572
						}
						position++
					}
				l1583:
					{
						position1585, tokenIndex1585, depth1585 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1586
						}
						position++
						goto l1585
					l1586:
						position, tokenIndex, depth = position1585, tokenIndex1585, depth1585
						if buffer[position] != rune('E') {
							goto l1572
						}
						position++
					}
				l1585:
					{
						position1587, tokenIndex1587, depth1587 := position, tokenIndex, depth
						if buffer[position] != rune('w') {
							goto l1588
						}
						position++
						goto l1587
					l1588:
						position, tokenIndex, depth = position1587, tokenIndex1587, depth1587
						if buffer[position] != rune('W') {
							goto l1572
						}
						position++
					}
				l1587:
					{
						position1589, tokenIndex1589, depth1589 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1590
						}
						position++
						goto l1589
					l1590:
						position, tokenIndex, depth = position1589, tokenIndex1589, depth1589
						if buffer[position] != rune('E') {
							goto l1572
						}
						position++
					}
				l1589:
					{
						position1591, tokenIndex1591, depth1591 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1592
						}
						position++
						goto l1591
					l1592:
						position, tokenIndex, depth = position1591, tokenIndex1591, depth1591
						if buffer[position] != rune('S') {
							goto l1572
						}
						position++
					}
				l1591:
					{
						position1593, tokenIndex1593, depth1593 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1594
						}
						position++
						goto l1593
					l1594:
						position, tokenIndex, depth = position1593, tokenIndex1593, depth1593
						if buffer[position] != rune('T') {
							goto l1572
						}
						position++
					}
				l1593:
					depth--
					add(rulePegText, position1574)
				}
				if !_rules[ruleAction98]() {
					goto l1572
				}
				depth--
				add(ruleDropNewest, position1573)
			}
			return true
		l1572:
			position, tokenIndex, depth = position1572, tokenIndex1572, depth1572
			return false
		},
		/* 129 StreamIdentifier <- <(<ident> Action99)> */
		func() bool {
			position1595, tokenIndex1595, depth1595 := position, tokenIndex, depth
			{
				position1596 := position
				depth++
				{
					position1597 := position
					depth++
					if !_rules[ruleident]() {
						goto l1595
					}
					depth--
					add(rulePegText, position1597)
				}
				if !_rules[ruleAction99]() {
					goto l1595
				}
				depth--
				add(ruleStreamIdentifier, position1596)
			}
			return true
		l1595:
			position, tokenIndex, depth = position1595, tokenIndex1595, depth1595
			return false
		},
		/* 130 SourceSinkType <- <(<ident> Action100)> */
		func() bool {
			position1598, tokenIndex1598, depth1598 := position, tokenIndex, depth
			{
				position1599 := position
				depth++
				{
					position1600 := position
					depth++
					if !_rules[ruleident]() {
						goto l1598
					}
					depth--
					add(rulePegText, position1600)
				}
				if !_rules[ruleAction100]() {
					goto l1598
				}
				depth--
				add(ruleSourceSinkType, position1599)
			}
			return true
		l1598:
			position, tokenIndex, depth = position1598, tokenIndex1598, depth1598
			return false
		},
		/* 131 SourceSinkParamKey <- <(<ident> Action101)> */
		func() bool {
			position1601, tokenIndex1601, depth1601 := position, tokenIndex, depth
			{
				position1602 := position
				depth++
				{
					position1603 := position
					depth++
					if !_rules[ruleident]() {
						goto l1601
					}
					depth--
					add(rulePegText, position1603)
				}
				if !_rules[ruleAction101]() {
					goto l1601
				}
				depth--
				add(ruleSourceSinkParamKey, position1602)
			}
			return true
		l1601:
			position, tokenIndex, depth = position1601, tokenIndex1601, depth1601
			return false
		},
		/* 132 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action102)> */
		func() bool {
			position1604, tokenIndex1604, depth1604 := position, tokenIndex, depth
			{
				position1605 := position
				depth++
				{
					position1606 := position
					depth++
					{
						position1607, tokenIndex1607, depth1607 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1608
						}
						position++
						goto l1607
					l1608:
						position, tokenIndex, depth = position1607, tokenIndex1607, depth1607
						if buffer[position] != rune('P') {
							goto l1604
						}
						position++
					}
				l1607:
					{
						position1609, tokenIndex1609, depth1609 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1610
						}
						position++
						goto l1609
					l1610:
						position, tokenIndex, depth = position1609, tokenIndex1609, depth1609
						if buffer[position] != rune('A') {
							goto l1604
						}
						position++
					}
				l1609:
					{
						position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1612
						}
						position++
						goto l1611
					l1612:
						position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
						if buffer[position] != rune('U') {
							goto l1604
						}
						position++
					}
				l1611:
					{
						position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1614
						}
						position++
						goto l1613
					l1614:
						position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
						if buffer[position] != rune('S') {
							goto l1604
						}
						position++
					}
				l1613:
					{
						position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1616
						}
						position++
						goto l1615
					l1616:
						position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
						if buffer[position] != rune('E') {
							goto l1604
						}
						position++
					}
				l1615:
					{
						position1617, tokenIndex1617, depth1617 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1618
						}
						position++
						goto l1617
					l1618:
						position, tokenIndex, depth = position1617, tokenIndex1617, depth1617
						if buffer[position] != rune('D') {
							goto l1604
						}
						position++
					}
				l1617:
					depth--
					add(rulePegText, position1606)
				}
				if !_rules[ruleAction102]() {
					goto l1604
				}
				depth--
				add(rulePaused, position1605)
			}
			return true
		l1604:
			position, tokenIndex, depth = position1604, tokenIndex1604, depth1604
			return false
		},
		/* 133 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action103)> */
		func() bool {
			position1619, tokenIndex1619, depth1619 := position, tokenIndex, depth
			{
				position1620 := position
				depth++
				{
					position1621 := position
					depth++
					{
						position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1623
						}
						position++
						goto l1622
					l1623:
						position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
						if buffer[position] != rune('U') {
							goto l1619
						}
						position++
					}
				l1622:
					{
						position1624, tokenIndex1624, depth1624 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1625
						}
						position++
						goto l1624
					l1625:
						position, tokenIndex, depth = position1624, tokenIndex1624, depth1624
						if buffer[position] != rune('N') {
							goto l1619
						}
						position++
					}
				l1624:
					{
						position1626, tokenIndex1626, depth1626 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1627
						}
						position++
						goto l1626
					l1627:
						position, tokenIndex, depth = position1626, tokenIndex1626, depth1626
						if buffer[position] != rune('P') {
							goto l1619
						}
						position++
					}
				l1626:
					{
						position1628, tokenIndex1628, depth1628 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1629
						}
						position++
						goto l1628
					l1629:
						position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
						if buffer[position] != rune('A') {
							goto l1619
						}
						position++
					}
				l1628:
					{
						position1630, tokenIndex1630, depth1630 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1631
						}
						position++
						goto l1630
					l1631:
						position, tokenIndex, depth = position1630, tokenIndex1630, depth1630
						if buffer[position] != rune('U') {
							goto l1619
						}
						position++
					}
				l1630:
					{
						position1632, tokenIndex1632, depth1632 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1633
						}
						position++
						goto l1632
					l1633:
						position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
						if buffer[position] != rune('S') {
							goto l1619
						}
						position++
					}
				l1632:
					{
						position1634, tokenIndex1634, depth1634 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1635
						}
						position++
						goto l1634
					l1635:
						position, tokenIndex, depth = position1634, tokenIndex1634, depth1634
						if buffer[position] != rune('E') {
							goto l1619
						}
						position++
					}
				l1634:
					{
						position1636, tokenIndex1636, depth1636 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1637
						}
						position++
						goto l1636
					l1637:
						position, tokenIndex, depth = position1636, tokenIndex1636, depth1636
						if buffer[position] != rune('D') {
							goto l1619
						}
						position++
					}
				l1636:
					depth--
					add(rulePegText, position1621)
				}
				if !_rules[ruleAction103]() {
					goto l1619
				}
				depth--
				add(ruleUnpaused, position1620)
			}
			return true
		l1619:
			position, tokenIndex, depth = position1619, tokenIndex1619, depth1619
			return false
		},
		/* 134 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action104)> */
		func() bool {
			position1638, tokenIndex1638, depth1638 := position, tokenIndex, depth
			{
				position1639 := position
				depth++
				{
					position1640 := position
					depth++
					{
						position1641, tokenIndex1641, depth1641 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1642
						}
						position++
						goto l1641
					l1642:
						position, tokenIndex, depth = position1641, tokenIndex1641, depth1641
						if buffer[position] != rune('A') {
							goto l1638
						}
						position++
					}
				l1641:
					{
						position1643, tokenIndex1643, depth1643 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1644
						}
						position++
						goto l1643
					l1644:
						position, tokenIndex, depth = position1643, tokenIndex1643, depth1643
						if buffer[position] != rune('S') {
							goto l1638
						}
						position++
					}
				l1643:
					{
						position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1646
						}
						position++
						goto l1645
					l1646:
						position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
						if buffer[position] != rune('C') {
							goto l1638
						}
						position++
					}
				l1645:
					depth--
					add(rulePegText, position1640)
				}
				if !_rules[ruleAction104]() {
					goto l1638
				}
				depth--
				add(ruleAscending, position1639)
			}
			return true
		l1638:
			position, tokenIndex, depth = position1638, tokenIndex1638, depth1638
			return false
		},
		/* 135 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action105)> */
		func() bool {
			position1647, tokenIndex1647, depth1647 := position, tokenIndex, depth
			{
				position1648 := position
				depth++
				{
					position1649 := position
					depth++
					{
						position1650, tokenIndex1650, depth1650 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1651
						}
						position++
						goto l1650
					l1651:
						position, tokenIndex, depth = position1650, tokenIndex1650, depth1650
						if buffer[position] != rune('D') {
							goto l1647
						}
						position++
					}
				l1650:
					{
						position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1653
						}
						position++
						goto l1652
					l1653:
						position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
						if buffer[position] != rune('E') {
							goto l1647
						}
						position++
					}
				l1652:
					{
						position1654, tokenIndex1654, depth1654 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1655
						}
						position++
						goto l1654
					l1655:
						position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
						if buffer[position] != rune('S') {
							goto l1647
						}
						position++
					}
				l1654:
					{
						position1656, tokenIndex1656, depth1656 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1657
						}
						position++
						goto l1656
					l1657:
						position, tokenIndex, depth = position1656, tokenIndex1656, depth1656
						if buffer[position] != rune('C') {
							goto l1647
						}
						position++
					}
				l1656:
					depth--
					add(rulePegText, position1649)
				}
				if !_rules[ruleAction105]() {
					goto l1647
				}
				depth--
				add(ruleDescending, position1648)
			}
			return true
		l1647:
			position, tokenIndex, depth = position1647, tokenIndex1647, depth1647
			return false
		},
		/* 136 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
			{
				position1659 := position
				depth++
				{
					position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1661
					}
					goto l1660
				l1661:
					position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
					if !_rules[ruleInt]() {
						goto l1662
					}
					goto l1660
				l1662:
					position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
					if !_rules[ruleFloat]() {
						goto l1663
					}
					goto l1660
				l1663:
					position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
					if !_rules[ruleString]() {
						goto l1664
					}
					goto l1660
				l1664:
					position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
					if !_rules[ruleBlob]() {
						goto l1665
					}
					goto l1660
				l1665:
					position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
					if !_rules[ruleTimestamp]() {
						goto l1666
					}
					goto l1660
				l1666:
					position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
					if !_rules[ruleArray]() {
						goto l1667
					}
					goto l1660
				l1667:
					position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
					if !_rules[ruleMap]() {
						goto l1658
					}
				}
			l1660:
				depth--
				add(ruleType, position1659)
			}
			return true
		l1658:
			position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
			return false
		},
		/* 137 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action106)> */
		func() bool {
			position1668, tokenIndex1668, depth1668 := position, tokenIndex, depth
			{
				position1669 := position
				depth++
				{
					position1670 := position
					depth++
					{
						position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1672
						}
						position++
						goto l1671
					l1672:
						position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
						if buffer[position] != rune('B') {
							goto l1668
						}
						position++
					}
				l1671:
					{
						position1673, tokenIndex1673, depth1673 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1674
						}
						position++
						goto l1673
					l1674:
						position, tokenIndex, depth = position1673, tokenIndex1673, depth1673
						if buffer[position] != rune('O') {
							goto l1668
						}
						position++
					}
				l1673:
					{
						position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1676
						}
						position++
						goto l1675
					l1676:
						position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
						if buffer[position] != rune('O') {
							goto l1668
						}
						position++
					}
				l1675:
					{
						position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1678
						}
						position++
						goto l1677
					l1678:
						position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
						if buffer[position] != rune('L') {
							goto l1668
						}
						position++
					}
				l1677:
					depth--
					add(rulePegText, position1670)
				}
				if !_rules[ruleAction106]() {
					goto l1668
				}
				depth--
				add(ruleBool, position1669)
			}
			return true
		l1668:
			position, tokenIndex, depth = position1668, tokenIndex1668, depth1668
			return false
		},
		/* 138 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action107)> */
		func() bool {
			position1679, tokenIndex1679, depth1679 := position, tokenIndex, depth
			{
				position1680 := position
				depth++
				{
					position1681 := position
					depth++
					{
						position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1683
						}
						position++
						goto l1682
					l1683:
						position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
						if buffer[position] != rune('I') {
							goto l1679
						}
						position++
					}
				l1682:
					{
						position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1685
						}
						position++
						goto l1684
					l1685:
						position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
						if buffer[position] != rune('N') {
							goto l1679
						}
						position++
					}
				l1684:
					{
						position1686, tokenIndex1686, depth1686 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1687
						}
						position++
						goto l1686
					l1687:
						position, tokenIndex, depth = position1686, tokenIndex1686, depth1686
						if buffer[position] != rune('T') {
							goto l1679
						}
						position++
					}
				l1686:
					depth--
					add(rulePegText, position1681)
				}
				if !_rules[ruleAction107]() {
					goto l1679
				}
				depth--
				add(ruleInt, position1680)
			}
			return true
		l1679:
			position, tokenIndex, depth = position1679, tokenIndex1679, depth1679
			return false
		},
		/* 139 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action108)> */
		func() bool {
			position1688, tokenIndex1688, depth1688 := position, tokenIndex, depth
			{
				position1689 := position
				depth++
				{
					position1690 := position
					depth++
					{
						position1691, tokenIndex1691, depth1691 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1692
						}
						position++
						goto l1691
					l1692:
						position, tokenIndex, depth = position1691, tokenIndex1691, depth1691
						if buffer[position] != rune('F') {
							goto l1688
						}
						position++
					}
				l1691:
					{
						position1693, tokenIndex1693, depth1693 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1694
						}
						position++
						goto l1693
					l1694:
						position, tokenIndex, depth = position1693, tokenIndex1693, depth1693
						if buffer[position] != rune('L') {
							goto l1688
						}
						position++
					}
				l1693:
					{
						position1695, tokenIndex1695, depth1695 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1696
						}
						position++
						goto l1695
					l1696:
						position, tokenIndex, depth = position1695, tokenIndex1695, depth1695
						if buffer[position] != rune('O') {
							goto l1688
						}
						position++
					}
				l1695:
					{
						position1697, tokenIndex1697, depth1697 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1698
						}
						position++
						goto l1697
					l1698:
						position, tokenIndex, depth = position1697, tokenIndex1697, depth1697
						if buffer[position] != rune('A') {
							goto l1688
						}
						position++
					}
				l1697:
					{
						position1699, tokenIndex1699, depth1699 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1700
						}
						position++
						goto l1699
					l1700:
						position, tokenIndex, depth = position1699, tokenIndex1699, depth1699
						if buffer[position] != rune('T') {
							goto l1688
						}
						position++
					}
				l1699:
					depth--
					add(rulePegText, position1690)
				}
				if !_rules[ruleAction108]() {
					goto l1688
				}
				depth--
				add(ruleFloat, position1689)
			}
			return true
		l1688:
			position, tokenIndex, depth = position1688, tokenIndex1688, depth1688
			return false
		},
		/* 140 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action109)> */
		func() bool {
			position1701, tokenIndex1701, depth1701 := position, tokenIndex, depth
			{
				position1702 := position
				depth++
				{
					position1703 := position
					depth++
					{
						position1704, tokenIndex1704, depth1704 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1705
						}
						position++
						goto l1704
					l1705:
						position, tokenIndex, depth = position1704, tokenIndex1704, depth1704
						if buffer[position] != rune('S') {
							goto l1701
						}
						position++
					}
				l1704:
					{
						position1706, tokenIndex1706, depth1706 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1707
						}
						position++
						goto l1706
					l1707:
						position, tokenIndex, depth = position1706, tokenIndex1706, depth1706
						if buffer[position] != rune('T') {
							goto l1701
						}
						position++
					}
				l1706:
					{
						position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1709
						}
						position++
						goto l1708
					l1709:
						position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
						if buffer[position] != rune('R') {
							goto l1701
						}
						position++
					}
				l1708:
					{
						position1710, tokenIndex1710, depth1710 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1711
						}
						position++
						goto l1710
					l1711:
						position, tokenIndex, depth = position1710, tokenIndex1710, depth1710
						if buffer[position] != rune('I') {
							goto l1701
						}
						position++
					}
				l1710:
					{
						position1712, tokenIndex1712, depth1712 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1713
						}
						position++
						goto l1712
					l1713:
						position, tokenIndex, depth = position1712, tokenIndex1712, depth1712
						if buffer[position] != rune('N') {
							goto l1701
						}
						position++
					}
				l1712:
					{
						position1714, tokenIndex1714, depth1714 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1715
						}
						position++
						goto l1714
					l1715:
						position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
						if buffer[position] != rune('G') {
							goto l1701
						}
						position++
					}
				l1714:
					depth--
					add(rulePegText, position1703)
				}
				if !_rules[ruleAction109]() {
					goto l1701
				}
				depth--
				add(ruleString, position1702)
			}
			return true
		l1701:
			position, tokenIndex, depth = position1701, tokenIndex1701, depth1701
			return false
		},
		/* 141 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action110)> */
		func() bool {
			position1716, tokenIndex1716, depth1716 := position, tokenIndex, depth
			{
				position1717 := position
				depth++
				{
					position1718 := position
					depth++
					{
						position1719, tokenIndex1719, depth1719 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1720
						}
						position++
						goto l1719
					l1720:
						position, tokenIndex, depth = position1719, tokenIndex1719, depth1719
						if buffer[position] != rune('B') {
							goto l1716
						}
						position++
					}
				l1719:
					{
						position1721, tokenIndex1721, depth1721 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1722
						}
						position++
						goto l1721
					l1722:
						position, tokenIndex, depth = position1721, tokenIndex1721, depth1721
						if buffer[position] != rune('L') {
							goto l1716
						}
						position++
					}
				l1721:
					{
						position1723, tokenIndex1723, depth1723 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1724
						}
						position++
						goto l1723
					l1724:
						position, tokenIndex, depth = position1723, tokenIndex1723, depth1723
						if buffer[position] != rune('O') {
							goto l1716
						}
						position++
					}
				l1723:
					{
						position1725, tokenIndex1725, depth1725 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1726
						}
						position++
						goto l1725
					l1726:
						position, tokenIndex, depth = position1725, tokenIndex1725, depth1725
						if buffer[position] != rune('B') {
							goto l1716
						}
						position++
					}
				l1725:
					depth--
					add(rulePegText, position1718)
				}
				if !_rules[ruleAction110]() {
					goto l1716
				}
				depth--
				add(ruleBlob, position1717)
			}
			return true
		l1716:
			position, tokenIndex, depth = position1716, tokenIndex1716, depth1716
			return false
		},
		/* 142 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action111)> */
		func() bool {
			position1727, tokenIndex1727, depth1727 := position, tokenIndex, depth
			{
				position1728 := position
				depth++
				{
					position1729 := position
					depth++
					{
						position1730, tokenIndex1730, depth1730 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1731
						}
						position++
						goto l1730
					l1731:
						position, tokenIndex, depth = position1730, tokenIndex1730, depth1730
						if buffer[position] != rune('T') {
							goto l1727
						}
						position++
					}
				l1730:
					{
						position1732, tokenIndex1732, depth1732 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1733
						}
						position++
						goto l1732
					l1733:
						position, tokenIndex, depth = position1732, tokenIndex1732, depth1732
						if buffer[position] != rune('I') {
							goto l1727
						}
						position++
					}
				l1732:
					{
						position1734, tokenIndex1734, depth1734 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1735
						}
						position++
						goto l1734
					l1735:
						position, tokenIndex, depth = position1734, tokenIndex1734, depth1734
						if buffer[position] != rune('M') {
							goto l1727
						}
						position++
					}
				l1734:
					{
						position1736, tokenIndex1736, depth1736 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1737
						}
						position++
						goto l1736
					l1737:
						position, tokenIndex, depth = position1736, tokenIndex1736, depth1736
						if buffer[position] != rune('E') {
							goto l1727
						}
						position++
					}
				l1736:
					{
						position1738, tokenIndex1738, depth1738 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1739
						}
						position++
						goto l1738
					l1739:
						position, tokenIndex, depth = position1738, tokenIndex1738, depth1738
						if buffer[position] != rune('S') {
							goto l1727
						}
						position++
					}
				l1738:
					{
						position1740, tokenIndex1740, depth1740 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1741
						}
						position++
						goto l1740
					l1741:
						position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
						if buffer[position] != rune('T') {
							goto l1727
						}
						position++
					}
				l1740:
					{
						position1742, tokenIndex1742, depth1742 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1743
						}
						position++
						goto l1742
					l1743:
						position, tokenIndex, depth = position1742, tokenIndex1742, depth1742
						if buffer[position] != rune('A') {
							goto l1727
						}
						position++
					}
				l1742:
					{
						position1744, tokenIndex1744, depth1744 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1745
						}
						position++
						goto l1744
					l1745:
						position, tokenIndex, depth = position1744, tokenIndex1744, depth1744
						if buffer[position] != rune('M') {
							goto l1727
						}
						position++
					}
				l1744:
					{
						position1746, tokenIndex1746, depth1746 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1747
						}
						position++
						goto l1746
					l1747:
						position, tokenIndex, depth = position1746, tokenIndex1746, depth1746
						if buffer[position] != rune('P') {
							goto l1727
						}
						position++
					}
				l1746:
					depth--
					add(rulePegText, position1729)
				}
				if !_rules[ruleAction111]() {
					goto l1727
				}
				depth--
				add(ruleTimestamp, position1728)
			}
			return true
		l1727:
			position, tokenIndex, depth = position1727, tokenIndex1727, depth1727
			return false
		},
		/* 143 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action112)> */
		func() bool {
			position1748, tokenIndex1748, depth1748 := position, tokenIndex, depth
			{
				position1749 := position
				depth++
				{
					position1750 := position
					depth++
					{
						position1751, tokenIndex1751, depth1751 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1752
						}
						position++
						goto l1751
					l1752:
						position, tokenIndex, depth = position1751, tokenIndex1751, depth1751
						if buffer[position] != rune('A') {
							goto l1748
						}
						position++
					}
				l1751:
					{
						position1753, tokenIndex1753, depth1753 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1754
						}
						position++
						goto l1753
					l1754:
						position, tokenIndex, depth = position1753, tokenIndex1753, depth1753
						if buffer[position] != rune('R') {
							goto l1748
						}
						position++
					}
				l1753:
					{
						position1755, tokenIndex1755, depth1755 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1756
						}
						position++
						goto l1755
					l1756:
						position, tokenIndex, depth = position1755, tokenIndex1755, depth1755
						if buffer[position] != rune('R') {
							goto l1748
						}
						position++
					}
				l1755:
					{
						position1757, tokenIndex1757, depth1757 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1758
						}
						position++
						goto l1757
					l1758:
						position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
						if buffer[position] != rune('A') {
							goto l1748
						}
						position++
					}
				l1757:
					{
						position1759, tokenIndex1759, depth1759 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1760
						}
						position++
						goto l1759
					l1760:
						position, tokenIndex, depth = position1759, tokenIndex1759, depth1759
						if buffer[position] != rune('Y') {
							goto l1748
						}
						position++
					}
				l1759:
					depth--
					add(rulePegText, position1750)
				}
				if !_rules[ruleAction112]() {
					goto l1748
				}
				depth--
				add(ruleArray, position1749)
			}
			return true
		l1748:
			position, tokenIndex, depth = position1748, tokenIndex1748, depth1748
			return false
		},
		/* 144 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action113)> */
		func() bool {
			position1761, tokenIndex1761, depth1761 := position, tokenIndex, depth
			{
				position1762 := position
				depth++
				{
					position1763 := position
					depth++
					{
						position1764, tokenIndex1764, depth1764 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1765
						}
						position++
						goto l1764
					l1765:
						position, tokenIndex, depth = position1764, tokenIndex1764, depth1764
						if buffer[position] != rune('M') {
							goto l1761
						}
						position++
					}
				l1764:
					{
						position1766, tokenIndex1766, depth1766 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1767
						}
						position++
						goto l1766
					l1767:
						position, tokenIndex, depth = position1766, tokenIndex1766, depth1766
						if buffer[position] != rune('A') {
							goto l1761
						}
						position++
					}
				l1766:
					{
						position1768, tokenIndex1768, depth1768 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1769
						}
						position++
						goto l1768
					l1769:
						position, tokenIndex, depth = position1768, tokenIndex1768, depth1768
						if buffer[position] != rune('P') {
							goto l1761
						}
						position++
					}
				l1768:
					depth--
					add(rulePegText, position1763)
				}
				if !_rules[ruleAction113]() {
					goto l1761
				}
				depth--
				add(ruleMap, position1762)
			}
			return true
		l1761:
			position, tokenIndex, depth = position1761, tokenIndex1761, depth1761
			return false
		},
		/* 145 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action114)> */
		func() bool {
			position1770, tokenIndex1770, depth1770 := position, tokenIndex, depth
			{
				position1771 := position
				depth++
				{
					position1772 := position
					depth++
					{
						position1773, tokenIndex1773, depth1773 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1774
						}
						position++
						goto l1773
					l1774:
						position, tokenIndex, depth = position1773, tokenIndex1773, depth1773
						if buffer[position] != rune('O') {
							goto l1770
						}
						position++
					}
				l1773:
					{
						position1775, tokenIndex1775, depth1775 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1776
						}
						position++
						goto l1775
					l1776:
						position, tokenIndex, depth = position1775, tokenIndex1775, depth1775
						if buffer[position] != rune('R') {
							goto l1770
						}
						position++
					}
				l1775:
					depth--
					add(rulePegText, position1772)
				}
				if !_rules[ruleAction114]() {
					goto l1770
				}
				depth--
				add(ruleOr, position1771)
			}
			return true
		l1770:
			position, tokenIndex, depth = position1770, tokenIndex1770, depth1770
			return false
		},
		/* 146 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action115)> */
		func() bool {
			position1777, tokenIndex1777, depth1777 := position, tokenIndex, depth
			{
				position1778 := position
				depth++
				{
					position1779 := position
					depth++
					{
						position1780, tokenIndex1780, depth1780 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1781
						}
						position++
						goto l1780
					l1781:
						position, tokenIndex, depth = position1780, tokenIndex1780, depth1780
						if buffer[position] != rune('A') {
							goto l1777
						}
						position++
					}
				l1780:
					{
						position1782, tokenIndex1782, depth1782 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1783
						}
						position++
						goto l1782
					l1783:
						position, tokenIndex, depth = position1782, tokenIndex1782, depth1782
						if buffer[position] != rune('N') {
							goto l1777
						}
						position++
					}
				l1782:
					{
						position1784, tokenIndex1784, depth1784 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1785
						}
						position++
						goto l1784
					l1785:
						position, tokenIndex, depth = position1784, tokenIndex1784, depth1784
						if buffer[position] != rune('D') {
							goto l1777
						}
						position++
					}
				l1784:
					depth--
					add(rulePegText, position1779)
				}
				if !_rules[ruleAction115]() {
					goto l1777
				}
				depth--
				add(ruleAnd, position1778)
			}
			return true
		l1777:
			position, tokenIndex, depth = position1777, tokenIndex1777, depth1777
			return false
		},
		/* 147 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action116)> */
		func() bool {
			position1786, tokenIndex1786, depth1786 := position, tokenIndex, depth
			{
				position1787 := position
				depth++
				{
					position1788 := position
					depth++
					{
						position1789, tokenIndex1789, depth1789 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1790
						}
						position++
						goto l1789
					l1790:
						position, tokenIndex, depth = position1789, tokenIndex1789, depth1789
						if buffer[position] != rune('N') {
							goto l1786
						}
						position++
					}
				l1789:
					{
						position1791, tokenIndex1791, depth1791 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1792
						}
						position++
						goto l1791
					l1792:
						position, tokenIndex, depth = position1791, tokenIndex1791, depth1791
						if buffer[position] != rune('O') {
							goto l1786
						}
						position++
					}
				l1791:
					{
						position1793, tokenIndex1793, depth1793 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1794
						}
						position++
						goto l1793
					l1794:
						position, tokenIndex, depth = position1793, tokenIndex1793, depth1793
						if buffer[position] != rune('T') {
							goto l1786
						}
						position++
					}
				l1793:
					depth--
					add(rulePegText, position1788)
				}
				if !_rules[ruleAction116]() {
					goto l1786
				}
				depth--
				add(ruleNot, position1787)
			}
			return true
		l1786:
			position, tokenIndex, depth = position1786, tokenIndex1786, depth1786
			return false
		},
		/* 148 Equal <- <(<'='> Action117)> */
		func() bool {
			position1795, tokenIndex1795, depth1795 := position, tokenIndex, depth
			{
				position1796 := position
				depth++
				{
					position1797 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1795
					}
					position++
					depth--
					add(rulePegText, position1797)
				}
				if !_rules[ruleAction117]() {
					goto l1795
				}
				depth--
				add(ruleEqual, position1796)
			}
			return true
		l1795:
			position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
			return false
		},
		/* 149 Less <- <(<'<'> Action118)> */
		func() bool {
			position1798, tokenIndex1798, depth1798 := position, tokenIndex, depth
			{
				position1799 := position
				depth++
				{
					position1800 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1798
					}
					position++
					depth--
					add(rulePegText, position1800)
				}
				if !_rules[ruleAction118]() {
					goto l1798
				}
				depth--
				add(ruleLess, position1799)
			}
			return true
		l1798:
			position, tokenIndex, depth = position1798, tokenIndex1798, depth1798
			return false
		},
		/* 150 LessOrEqual <- <(<('<' '=')> Action119)> */
		func() bool {
			position1801, tokenIndex1801, depth1801 := position, tokenIndex, depth
			{
				position1802 := position
				depth++
				{
					position1803 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1801
					}
					position++
					if buffer[position] != rune('=') {
						goto l1801
					}
					position++
					depth--
					add(rulePegText, position1803)
				}
				if !_rules[ruleAction119]() {
					goto l1801
				}
				depth--
				add(ruleLessOrEqual, position1802)
			}
			return true
		l1801:
			position, tokenIndex, depth = position1801, tokenIndex1801, depth1801
			return false
		},
		/* 151 Greater <- <(<'>'> Action120)> */
		func() bool {
			position1804, tokenIndex1804, depth1804 := position, tokenIndex, depth
			{
				position1805 := position
				depth++
				{
					position1806 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1804
					}
					position++
					depth--
					add(rulePegText, position1806)
				}
				if !_rules[ruleAction120]() {
					goto l1804
				}
				depth--
				add(ruleGreater, position1805)
			}
			return true
		l1804:
			position, tokenIndex, depth = position1804, tokenIndex1804, depth1804
			return false
		},
		/* 152 GreaterOrEqual <- <(<('>' '=')> Action121)> */
		func() bool {
			position1807, tokenIndex1807, depth1807 := position, tokenIndex, depth
			{
				position1808 := position
				depth++
				{
					position1809 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1807
					}
					position++
					if buffer[position] != rune('=') {
						goto l1807
					}
					position++
					depth--
					add(rulePegText, position1809)
				}
				if !_rules[ruleAction121]() {
					goto l1807
				}
				depth--
				add(ruleGreaterOrEqual, position1808)
			}
			return true
		l1807:
			position, tokenIndex, depth = position1807, tokenIndex1807, depth1807
			return false
		},
		/* 153 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action122)> */
		func() bool {
			position1810, tokenIndex1810, depth1810 := position, tokenIndex, depth
			{
				position1811 := position
				depth++
				{
					position1812 := position
					depth++
					{
						position1813, tokenIndex1813, depth1813 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1814
						}
						position++
						if buffer[position] != rune('=') {
							goto l1814
						}
						position++
						goto l1813
					l1814:
						position, tokenIndex, depth = position1813, tokenIndex1813, depth1813
						if buffer[position] != rune('<') {
							goto l1810
						}
						position++
						if buffer[position] != rune('>') {
							goto l1810
						}
						position++
					}
				l1813:
					depth--
					add(rulePegText, position1812)
				}
				if !_rules[ruleAction122]() {
					goto l1810
				}
				depth--
				add(ruleNotEqual, position1811)
			}
			return true
		l1810:
			position, tokenIndex, depth = position1810, tokenIndex1810, depth1810
			return false
		},
		/* 154 Concat <- <(<('|' '|')> Action123)> */
		func() bool {
			position1815, tokenIndex1815, depth1815 := position, tokenIndex, depth
			{
				position1816 := position
				depth++
				{
					position1817 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1815
					}
					position++
					if buffer[position] != rune('|') {
						goto l1815
					}
					position++
					depth--
					add(rulePegText, position1817)
				}
				if !_rules[ruleAction123]() {
					goto l1815
				}
				depth--
				add(ruleConcat, position1816)
			}
			return true
		l1815:
			position, tokenIndex, depth = position1815, tokenIndex1815, depth1815
			return false
		},
		/* 155 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action124)> */
		func() bool {
			position1818, tokenIndex1818, depth1818 := position, tokenIndex, depth
			{
				position1819 := position
				depth++
				{
					position1820 := position
					depth++
					{
						position1821, tokenIndex1821, depth1821 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1822
						}
						position++
						goto l1821
					l1822:
						position, tokenIndex, depth = position1821, tokenIndex1821, depth1821
						if buffer[position] != rune('I') {
							goto l1818
						}
						position++
					}
				l1821:
					{
						position1823, tokenIndex1823, depth1823 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1824
						}
						position++
						goto l1823
					l1824:
						position, tokenIndex, depth = position1823, tokenIndex1823, depth1823
						if buffer[position] != rune('S') {
							goto l1818
						}
						position++
					}
				l1823:
					depth--
					add(rulePegText, position1820)
				}
				if !_rules[ruleAction124]() {
					goto l1818
				}
				depth--
				add(ruleIs, position1819)
			}
			return true
		l1818:
			position, tokenIndex, depth = position1818, tokenIndex1818, depth1818
			return false
		},
		/* 156 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action125)> */
		func() bool {
			position1825, tokenIndex1825, depth1825 := position, tokenIndex, depth
			{
				position1826 := position
				depth++
				{
					position1827 := position
					depth++
					{
						position1828, tokenIndex1828, depth1828 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1829
						}
						position++
						goto l1828
					l1829:
						position, tokenIndex, depth = position1828, tokenIndex1828, depth1828
						if buffer[position] != rune('I') {
							goto l1825
						}
						position++
					}
				l1828:
					{
						position1830, tokenIndex1830, depth1830 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1831
						}
						position++
						goto l1830
					l1831:
						position, tokenIndex, depth = position1830, tokenIndex1830, depth1830
						if buffer[position] != rune('S') {
							goto l1825
						}
						position++
					}
				l1830:
					if !_rules[rulesp]() {
						goto l1825
					}
					{
						position1832, tokenIndex1832, depth1832 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1833
						}
						position++
						goto l1832
					l1833:
						position, tokenIndex, depth = position1832, tokenIndex1832, depth1832
						if buffer[position] != rune('N') {
							goto l1825
						}
						position++
					}
				l1832:
					{
						position1834, tokenIndex1834, depth1834 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1835
						}
						position++
						goto l1834
					l1835:
						position, tokenIndex, depth = position1834, tokenIndex1834, depth1834
						if buffer[position] != rune('O') {
							goto l1825
						}
						position++
					}
				l1834:
					{
						position1836, tokenIndex1836, depth1836 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1837
						}
						position++
						goto l1836
					l1837:
						position, tokenIndex, depth = position1836, tokenIndex1836, depth1836
						if buffer[position] != rune('T') {
							goto l1825
						}
						position++
					}
				l1836:
					depth--
					add(rulePegText, position1827)
				}
				if !_rules[ruleAction125]() {
					goto l1825
				}
				depth--
				add(ruleIsNot, position1826)
			}
			return true
		l1825:
			position, tokenIndex, depth = position1825, tokenIndex1825, depth1825
			return false
		},
		/* 157 Plus <- <(<'+'> Action126)> */
		func() bool {
			position1838, tokenIndex1838, depth1838 := position, tokenIndex, depth
			{
				position1839 := position
				depth++
				{
					position1840 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1838
					}
					position++
					depth--
					add(rulePegText, position1840)
				}
				if !_rules[ruleAction126]() {
					goto l1838
				}
				depth--
				add(rulePlus, position1839)
			}
			return true
		l1838:
			position, tokenIndex, depth = position1838, tokenIndex1838, depth1838
			return false
		},
		/* 158 Minus <- <(<'-'> Action127)> */
		func() bool {
			position1841, tokenIndex1841, depth1841 := position, tokenIndex, depth
			{
				position1842 := position
				depth++
				{
					position1843 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1841
					}
					position++
					depth--
					add(rulePegText, position1843)
				}
				if !_rules[ruleAction127]() {
					goto l1841
				}
				depth--
				add(ruleMinus, position1842)
			}
			return true
		l1841:
			position, tokenIndex, depth = position1841, tokenIndex1841, depth1841
			return false
		},
		/* 159 Multiply <- <(<'*'> Action128)> */
		func() bool {
			position1844, tokenIndex1844, depth1844 := position, tokenIndex, depth
			{
				position1845 := position
				depth++
				{
					position1846 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1844
					}
					position++
					depth--
					add(rulePegText, position1846)
				}
				if !_rules[ruleAction128]() {
					goto l1844
				}
				depth--
				add(ruleMultiply, position1845)
			}
			return true
		l1844:
			position, tokenIndex, depth = position1844, tokenIndex1844, depth1844
			return false
		},
		/* 160 Divide <- <(<'/'> Action129)> */
		func() bool {
			position1847, tokenIndex1847, depth1847 := position, tokenIndex, depth
			{
				position1848 := position
				depth++
				{
					position1849 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1847
					}
					position++
					depth--
					add(rulePegText, position1849)
				}
				if !_rules[ruleAction129]() {
					goto l1847
				}
				depth--
				add(ruleDivide, position1848)
			}
			return true
		l1847:
			position, tokenIndex, depth = position1847, tokenIndex1847, depth1847
			return false
		},
		/* 161 Modulo <- <(<'%'> Action130)> */
		func() bool {
			position1850, tokenIndex1850, depth1850 := position, tokenIndex, depth
			{
				position1851 := position
				depth++
				{
					position1852 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1850
					}
					position++
					depth--
					add(rulePegText, position1852)
				}
				if !_rules[ruleAction130]() {
					goto l1850
				}
				depth--
				add(ruleModulo, position1851)
			}
			return true
		l1850:
			position, tokenIndex, depth = position1850, tokenIndex1850, depth1850
			return false
		},
		/* 162 UnaryMinus <- <(<'-'> Action131)> */
		func() bool {
			position1853, tokenIndex1853, depth1853 := position, tokenIndex, depth
			{
				position1854 := position
				depth++
				{
					position1855 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1853
					}
					position++
					depth--
					add(rulePegText, position1855)
				}
				if !_rules[ruleAction131]() {
					goto l1853
				}
				depth--
				add(ruleUnaryMinus, position1854)
			}
			return true
		l1853:
			position, tokenIndex, depth = position1853, tokenIndex1853, depth1853
			return false
		},
		/* 163 Identifier <- <(<ident> Action132)> */
		func() bool {
			position1856, tokenIndex1856, depth1856 := position, tokenIndex, depth
			{
				position1857 := position
				depth++
				{
					position1858 := position
					depth++
					if !_rules[ruleident]() {
						goto l1856
					}
					depth--
					add(rulePegText, position1858)
				}
				if !_rules[ruleAction132]() {
					goto l1856
				}
				depth--
				add(ruleIdentifier, position1857)
			}
			return true
		l1856:
			position, tokenIndex, depth = position1856, tokenIndex1856, depth1856
			return false
		},
		/* 164 TargetIdentifier <- <(<('*' / jsonSetPath)> Action133)> */
		func() bool {
			position1859, tokenIndex1859, depth1859 := position, tokenIndex, depth
			{
				position1860 := position
				depth++
				{
					position1861 := position
					depth++
					{
						position1862, tokenIndex1862, depth1862 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l1863
						}
						position++
						goto l1862
					l1863:
						position, tokenIndex, depth = position1862, tokenIndex1862, depth1862
						if !_rules[rulejsonSetPath]() {
							goto l1859
						}
					}
				l1862:
					depth--
					add(rulePegText, position1861)
				}
				if !_rules[ruleAction133]() {
					goto l1859
				}
				depth--
				add(ruleTargetIdentifier, position1860)
			}
			return true
		l1859:
			position, tokenIndex, depth = position1859, tokenIndex1859, depth1859
			return false
		},
		/* 165 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1864, tokenIndex1864, depth1864 := position, tokenIndex, depth
			{
				position1865 := position
				depth++
				{
					position1866, tokenIndex1866, depth1866 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1867
					}
					position++
					goto l1866
				l1867:
					position, tokenIndex, depth = position1866, tokenIndex1866, depth1866
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1864
					}
					position++
				}
			l1866:
			l1868:
				{
					position1869, tokenIndex1869, depth1869 := position, tokenIndex, depth
					{
						position1870, tokenIndex1870, depth1870 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1871
						}
						position++
						goto l1870
					l1871:
						position, tokenIndex, depth = position1870, tokenIndex1870, depth1870
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1872
						}
						position++
						goto l1870
					l1872:
						position, tokenIndex, depth = position1870, tokenIndex1870, depth1870
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1873
						}
						position++
						goto l1870
					l1873:
						position, tokenIndex, depth = position1870, tokenIndex1870, depth1870
						if buffer[position] != rune('_') {
							goto l1869
						}
						position++
					}
				l1870:
					goto l1868
				l1869:
					position, tokenIndex, depth = position1869, tokenIndex1869, depth1869
				}
				depth--
				add(ruleident, position1865)
			}
			return true
		l1864:
			position, tokenIndex, depth = position1864, tokenIndex1864, depth1864
			return false
		},
		/* 166 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1874, tokenIndex1874, depth1874 := position, tokenIndex, depth
			{
				position1875 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1874
				}
			l1876:
				{
					position1877, tokenIndex1877, depth1877 := position, tokenIndex, depth
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1877
					}
					goto l1876
				l1877:
					position, tokenIndex, depth = position1877, tokenIndex1877, depth1877
				}
				depth--
				add(rulejsonGetPath, position1875)
			}
			return true
		l1874:
			position, tokenIndex, depth = position1874, tokenIndex1874, depth1874
			return false
		},
		/* 167 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1878, tokenIndex1878, depth1878 := position, tokenIndex, depth
			{
				position1879 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1878
				}
			l1880:
				{
					position1881, tokenIndex1881, depth1881 := position, tokenIndex, depth
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1881
					}
					goto l1880
				l1881:
					position, tokenIndex, depth = position1881, tokenIndex1881, depth1881
				}
				depth--
				add(rulejsonSetPath, position1879)
			}
			return true
		l1878:
			position, tokenIndex, depth = position1878, tokenIndex1878, depth1878
			return false
		},
		/* 168 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1882, tokenIndex1882, depth1882 := position, tokenIndex, depth
			{
				position1883 := position
				depth++
				{
					position1884, tokenIndex1884, depth1884 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1885
					}
					goto l1884
				l1885:
					position, tokenIndex, depth = position1884, tokenIndex1884, depth1884
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1882
					}
				}
			l1884:
				depth--
				add(rulejsonPathHead, position1883)
			}
			return true
		l1882:
			position, tokenIndex, depth = position1882, tokenIndex1882, depth1882
			return false
		},
		/* 169 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1886, tokenIndex1886, depth1886 := position, tokenIndex, depth
			{
				position1887 := position
				depth++
				{
					position1888, tokenIndex1888, depth1888 := position, tokenIndex, depth
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1889
					}
					goto l1888
				l1889:
					position, tokenIndex, depth = position1888, tokenIndex1888, depth1888
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1890
					}
					goto l1888
				l1890:
					position, tokenIndex, depth = position1888, tokenIndex1888, depth1888
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1891
					}
					goto l1888
				l1891:
					position, tokenIndex, depth = position1888, tokenIndex1888, depth1888
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1892
					}
					goto l1888
				l1892:
					position, tokenIndex, depth = position1888, tokenIndex1888, depth1888
					if !_rules[rulejsonArraySlice]() {
						goto l1893
					}
					goto l1888
				l1893:
					position, tokenIndex, depth = position1888, tokenIndex1888, depth1888
					if !_rules[rulejsonArrayAccess]() {
						goto l1886
					}
				}
			l1888:
				depth--
				add(rulejsonGetPathNonHead, position1887)
			}
			return true
		l1886:
			position, tokenIndex, depth = position1886, tokenIndex1886, depth1886
			return false
		},
		/* 170 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1894, tokenIndex1894, depth1894 := position, tokenIndex, depth
			{
				position1895 := position
				depth++
				{
					position1896, tokenIndex1896, depth1896 := position, tokenIndex, depth
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1897
					}
					goto l1896
				l1897:
					position, tokenIndex, depth = position1896, tokenIndex1896, depth1896
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1894
					}
				}
			l1896:
				depth--
				add(rulejsonSetPathNonHead, position1895)
			}
			return true
		l1894:
			position, tokenIndex, depth = position1894, tokenIndex1894, depth1894
			return false
		},
		/* 171 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1898, tokenIndex1898, depth1898 := position, tokenIndex, depth
			{
				position1899 := position
				depth++
				{
					position1900, tokenIndex1900, depth1900 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1901
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1901
					}
					goto l1900
				l1901:
					position, tokenIndex, depth = position1900, tokenIndex1900, depth1900
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1898
					}
				}
			l1900:
				depth--
				add(rulejsonMapSingleLevel, position1899)
			}
			return true
		l1898:
			position, tokenIndex, depth = position1898, tokenIndex1898, depth1898
			return false
		},
		/* 172 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1902, tokenIndex1902, depth1902 := position, tokenIndex, depth
			{
				position1903 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1902
				}
				position++
				if buffer[position] != rune('.') {
					goto l1902
				}
				position++
				{
					position1904, tokenIndex1904, depth1904 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1905
					}
					goto l1904
				l1905:
					position, tokenIndex, depth = position1904, tokenIndex1904, depth1904
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1902
					}
				}
			l1904:
				depth--
				add(rulejsonMapMultipleLevel, position1903)
			}
			return true
		l1902:
			position, tokenIndex, depth = position1902, tokenIndex1902, depth1902
			return false
		},
		/* 173 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1906, tokenIndex1906, depth1906 := position, tokenIndex, depth
			{
				position1907 := position
				depth++
				{
					position1908 := position
					depth++
					{
						position1909, tokenIndex1909, depth1909 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1910
						}
						position++
						goto l1909
					l1910:
						position, tokenIndex, depth = position1909, tokenIndex1909, depth1909
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1906
						}
						position++
					}
				l1909:
				l1911:
					{
						position1912, tokenIndex1912, depth1912 := position, tokenIndex, depth
						{
							position1913, tokenIndex1913, depth1913 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1914
							}
							position++
							goto l1913
						l1914:
							position, tokenIndex, depth = position1913, tokenIndex1913, depth1913
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1915
							}
							position++
							goto l1913
						l1915:
							position, tokenIndex, depth = position1913, tokenIndex1913, depth1913
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1916
							}
							position++
							goto l1913
						l1916:
							position, tokenIndex, depth = position1913, tokenIndex1913, depth1913
							if buffer[position] != rune('_') {
								goto l1912
							}
							position++
						}
					l1913:
						goto l1911
					l1912:
						position, tokenIndex, depth = position1912, tokenIndex1912, depth1912
					}
					depth--
					add(rulePegText, position1908)
				}
				depth--
				add(rulejsonMapAccessString, position1907)
			}
			return true
		l1906:
			position, tokenIndex, depth = position1906, tokenIndex1906, depth1906
			return false
		},
		/* 174 jsonMapAccessBracket <- <('[' singleQuotedString ']')> */
		func() bool {
			position1917, tokenIndex1917, depth1917 := position, tokenIndex, depth
			{
				position1918 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1917
				}
				position++
				if !_rules[rulesingleQuotedString]() {
					goto l1917
				}
				if buffer[position] != rune(']') {
					goto l1917
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1918)
			}
			return true
		l1917:
			position, tokenIndex, depth = position1917, tokenIndex1917, depth1917
			return false
		},
		/* 175 singleQuotedString <- <('\'' <(('\'' '\'') / (!'\'' .))*> '\'')> */
		func() bool {
			position1919, tokenIndex1919, depth1919 := position, tokenIndex, depth
			{
				position1920 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l1919
				}
				position++
				{
					position1921 := position
					depth++
				l1922:
					{
						position1923, tokenIndex1923, depth1923 := position, tokenIndex, depth
						{
							position1924, tokenIndex1924, depth1924 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1925
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1925
							}
							position++
							goto l1924
						l1925:
							position, tokenIndex, depth = position1924, tokenIndex1924, depth1924
							{
								position1926, tokenIndex1926, depth1926 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1926
								}
								position++
								goto l1923
							l1926:
								position, tokenIndex, depth = position1926, tokenIndex1926, depth1926
							}
							if !matchDot() {
								goto l1923
							}
						}
					l1924:
						goto l1922
					l1923:
						position, tokenIndex, depth = position1923, tokenIndex1923, depth1923
					}
					depth--
					add(rulePegText, position1921)
				}
				if buffer[position] != rune('\'') {
					goto l1919
				}
				position++
				depth--
				add(rulesingleQuotedString, position1920)
			}
			return true
		l1919:
			position, tokenIndex, depth = position1919, tokenIndex1919, depth1919
			return false
		},
		/* 176 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1927, tokenIndex1927, depth1927 := position, tokenIndex, depth
			{
				position1928 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1927
				}
				position++
				{
					position1929 := position
					depth++
					{
						position1930, tokenIndex1930, depth1930 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1930
						}
						position++
						goto l1931
					l1930:
						position, tokenIndex, depth = position1930, tokenIndex1930, depth1930
					}
				l1931:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1927
					}
					position++
				l1932:
					{
						position1933, tokenIndex1933, depth1933 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1933
						}
						position++
						goto l1932
					l1933:
						position, tokenIndex, depth = position1933, tokenIndex1933, depth1933
					}
					depth--
					add(rulePegText, position1929)
				}
				if buffer[position] != rune(']') {
					goto l1927
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1928)
			}
			return true
		l1927:
			position, tokenIndex, depth = position1927, tokenIndex1927, depth1927
			return false
		},
		/* 177 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
		func() bool {
			position1934, tokenIndex1934, depth1934 := position, tokenIndex, depth
			{
				position1935 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1934
				}
				position++
				{
					position1936 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1934
					}
					position++
				l1937:
					{
						position1938, tokenIndex1938, depth1938 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1938
						}
						position++
						goto l1937
					l1938:
						position, tokenIndex, depth = position1938, tokenIndex1938, depth1938
					}
					depth--
					add(rulePegText, position1936)
				}
				if buffer[position] != rune(']') {
					goto l1934
				}
				position++
				depth--
				add(rulejsonNonNegativeArrayAccess, position1935)
			}
			return true
		l1934:
			position, tokenIndex, depth = position1934, tokenIndex1934, depth1934
			return false
		},
		/* 178 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1939, tokenIndex1939, depth1939 := position, tokenIndex, depth
			{
				position1940 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1939
				}
				position++
				{
					position1941 := position
					depth++
					{
						position1942, tokenIndex1942, depth1942 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1942
						}
						position++
						goto l1943
					l1942:
						position, tokenIndex, depth = position1942, tokenIndex1942, depth1942
					}
				l1943:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1939
					}
					position++
				l1944:
					{
						position1945, tokenIndex1945, depth1945 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1945
						}
						position++
						goto l1944
					l1945:
						position, tokenIndex, depth = position1945, tokenIndex1945, depth1945
					}
					if buffer[position] != rune(':') {
						goto l1939
					}
					position++
					{
						position1946, tokenIndex1946, depth1946 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1946
						}
						position++
						goto l1947
					l1946:
						position, tokenIndex, depth = position1946, tokenIndex1946, depth1946
					}
				l1947:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1939
					}
					position++
				l1948:
					{
						position1949, tokenIndex1949, depth1949 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1949
						}
						position++
						goto l1948
					l1949:
						position, tokenIndex, depth = position1949, tokenIndex1949, depth1949
					}
					{
						position1950, tokenIndex1950, depth1950 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1950
						}
						position++
						{
							position1952, tokenIndex1952, depth1952 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1952
							}
							position++
							goto l1953
						l1952:
							position, tokenIndex, depth = position1952, tokenIndex1952, depth1952
						}
					l1953:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1950
						}
						position++
					l1954:
						{
							position1955, tokenIndex1955, depth1955 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1955
							}
							position++
							goto l1954
						l1955:
							position, tokenIndex, depth = position1955, tokenIndex1955, depth1955
						}
						goto l1951
					l1950:
						position, tokenIndex, depth = position1950, tokenIndex1950, depth1950
					}
				l1951:
					depth--
					add(rulePegText, position1941)
				}
				if buffer[position] != rune(']') {
					goto l1939
				}
				position++
				depth--
				add(rulejsonArraySlice, position1940)
			}
			return true
		l1939:
			position, tokenIndex, depth = position1939, tokenIndex1939, depth1939
			return false
		},
		/* 179 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1956, tokenIndex1956, depth1956 := position, tokenIndex, depth
			{
				position1957 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1956
				}
				position++
				{
					position1958 := position
					depth++
					{
						position1959, tokenIndex1959, depth1959 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1960
						}
						position++
						{
							position1961, tokenIndex1961, depth1961 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1961
							}
							position++
							goto l1962
						l1961:
							position, tokenIndex, depth = position1961, tokenIndex1961, depth1961
						}
					l1962:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1960
						}
						position++
					l1963:
						{
							position1964, tokenIndex1964, depth1964 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1964
							}
							position++
							goto l1963
						l1964:
							position, tokenIndex, depth = position1964, tokenIndex1964, depth1964
						}
						goto l1959
					l1960:
						position, tokenIndex, depth = position1959, tokenIndex1959, depth1959
						{
							position1965, tokenIndex1965, depth1965 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1965
							}
							position++
							goto l1966
						l1965:
							position, tokenIndex, depth = position1965, tokenIndex1965, depth1965
						}
					l1966:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1956
						}
						position++
					l1967:
						{
							position1968, tokenIndex1968, depth1968 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1968
							}
							position++
							goto l1967
						l1968:
							position, tokenIndex, depth = position1968, tokenIndex1968, depth1968
						}
						if buffer[position] != rune(':') {
							goto l1956
						}
						position++
					}
				l1959:
					depth--
					add(rulePegText, position1958)
				}
				if buffer[position] != rune(']') {
					goto l1956
				}
				position++
				depth--
				add(rulejsonArrayPartialSlice, position1957)
			}
			return true
		l1956:
			position, tokenIndex, depth = position1956, tokenIndex1956, depth1956
			return false
		},
		/* 180 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1969, tokenIndex1969, depth1969 := position, tokenIndex, depth
			{
				position1970 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1969
				}
				position++
				if buffer[position] != rune(':') {
					goto l1969
				}
				position++
				if buffer[position] != rune(']') {
					goto l1969
				}
				position++
				depth--
				add(rulejsonArrayFullSlice, position1970)
			}
			return true
		l1969:
			position, tokenIndex, depth = position1969, tokenIndex1969, depth1969
			return false
		},
		/* 181 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1971, tokenIndex1971, depth1971 := position, tokenIndex, depth
			{
				position1972 := position
				depth++
				{
					position1973, tokenIndex1973, depth1973 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l1974
					}
					position++
					goto l1973
				l1974:
					position, tokenIndex, depth = position1973, tokenIndex1973, depth1973
					if buffer[position] != rune('\t') {
						goto l1975
					}
					position++
					goto l1973
				l1975:
					position, tokenIndex, depth = position1973, tokenIndex1973, depth1973
					if buffer[position] != rune('\n') {
						goto l1976
					}
					position++
					goto l1973
				l1976:
					position, tokenIndex, depth = position1973, tokenIndex1973, depth1973
					if buffer[position] != rune('\r') {
						goto l1977
					}
					position++
					goto l1973
				l1977:
					position, tokenIndex, depth = position1973, tokenIndex1973, depth1973
					if !_rules[rulecomment]() {
						goto l1978
					}
					goto l1973
				l1978:
					position, tokenIndex, depth = position1973, tokenIndex1973, depth1973
					if !_rules[rulefinalComment]() {
						goto l1971
					}
				}
			l1973:
				depth--
				add(rulespElem, position1972)
			}
			return true
		l1971:
			position, tokenIndex, depth = position1971, tokenIndex1971, depth1971
			return false
		},
		/* 182 sp <- <spElem+> */
		func() bool {
			position1979, tokenIndex1979, depth1979 := position, tokenIndex, depth
			{
				position1980 := position
				depth++
				if !_rules[rulespElem]() {
					goto l1979
				}
			l1981:
				{
					position1982, tokenIndex1982, depth1982 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1982
					}
					goto l1981
				l1982:
					position, tokenIndex, depth = position1982, tokenIndex1982, depth1982
				}
				depth--
				add(rulesp, position1980)
			}
			return true
		l1979:
			position, tokenIndex, depth = position1979, tokenIndex1979, depth1979
			return false
		},
		/* 183 spOpt <- <spElem*> */
		func() bool {
			{
				position1984 := position
				depth++
			l1985:
				{
					position1986, tokenIndex1986, depth1986 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1986
					}
					goto l1985
				l1986:
					position, tokenIndex, depth = position1986, tokenIndex1986, depth1986
				}
				depth--
				add(rulespOpt, position1984)
			}
			return true
		},
		/* 184 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1987, tokenIndex1987, depth1987 := position, tokenIndex, depth
			{
				position1988 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1987
				}
				position++
				if buffer[position] != rune('-') {
					goto l1987
				}
				position++
			l1989:
				{
					position1990, tokenIndex1990, depth1990 := position, tokenIndex, depth
					{
						position1991, tokenIndex1991, depth1991 := position, tokenIndex, depth
						{
							position1992, tokenIndex1992, depth1992 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1993
							}
							position++
							goto l1992
						l1993:
							position, tokenIndex, depth = position1992, tokenIndex1992, depth1992
							if buffer[position] != rune('\n') {
								goto l1991
							}
							position++
						}
					l1992:
						goto l1990
					l1991:
						position, tokenIndex, depth = position1991, tokenIndex1991, depth1991
					}
					if !matchDot() {
						goto l1990
					}
					goto l1989
				l1990:
					position, tokenIndex, depth = position1990, tokenIndex1990, depth1990
				}
				{
					position1994, tokenIndex1994, depth1994 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1995
					}
					position++
					goto l1994
				l1995:
					position, tokenIndex, depth = position1994, tokenIndex1994, depth1994
					if buffer[position] != rune('\n') {
						goto l1987
					}
					position++
				}
			l1994:
				depth--
				add(rulecomment, position1988)
			}
			return true
		l1987:
			position, tokenIndex, depth = position1987, tokenIndex1987, depth1987
			return false
		},
		/* 185 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1996, tokenIndex1996, depth1996 := position, tokenIndex, depth
			{
				position1997 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1996
				}
				position++
				if buffer[position] != rune('-') {
					goto l1996
				}
				position++
			l1998:
				{
					position1999, tokenIndex1999, depth1999 := position, tokenIndex, depth
					{
						position2000, tokenIndex2000, depth2000 := position, tokenIndex, depth
						{
							position2001, tokenIndex2001, depth2001 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l2002
							}
							position++
							goto l2001
						l2002:
							position, tokenIndex, depth = position2001, tokenIndex2001, depth2001
							if buffer[position] != rune('\n') {
								goto l2000
							}
							position++
						}
					l2001:
						goto l1999
					l2000:
						position, tokenIndex, depth = position2000, tokenIndex2000, depth2000
					}
					if !matchDot() {
						goto l1999
					}
					goto l1998
				l1999:
					position, tokenIndex, depth = position1999, tokenIndex1999, depth1999
				}
				{
					position2003, tokenIndex2003, depth2003 := position, tokenIndex, depth
					if !matchDot() {
						goto l2003
					}
					goto l1996
				l2003:
					position, tokenIndex, depth = position2003, tokenIndex2003, depth2003
				}
				depth--
				add(rulefinalComment, position1997)
			}
			return true
		l1996:
			position, tokenIndex, depth = position1996, tokenIndex1996, depth1996
			return false
		},
		nil,
		/* 188 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 189 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 190 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 191 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 192 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 193 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 194 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 195 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 196 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 197 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 198 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 199 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 200 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 201 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 202 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 203 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 204 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 205 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 206 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 207 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 208 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 209 Action21 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 210 Action22 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 211 Action23 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 212 Action24 <- <{
		    p.AssembleEval(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 213 Action25 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 214 Action26 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 215 Action27 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 216 Action28 <- <{
		    p.AssembleEmitterSampling(CountBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 217 Action29 <- <{
		    p.AssembleEmitterSampling(RandomizedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 218 Action30 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 219 Action31 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 0.001)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 220 Action32 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 221 Action33 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 222 Action34 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 223 Action35 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 224 Action36 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 225 Action37 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 226 Action38 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 227 Action39 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 228 Action40 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 229 Action41 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 230 Action42 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 231 Action43 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 232 Action44 <- <{
		    p.EnsureCapacitySpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 233 Action45 <- <{
		    p.EnsureSheddingSpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 234 Action46 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 235 Action47 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 236 Action48 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 237 Action49 <- <{
		    p.EnsureIdentifier(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 238 Action50 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 239 Action51 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 240 Action52 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 241 Action53 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 242 Action54 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 243 Action55 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 244 Action56 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 245 Action57 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 246 Action58 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 247 Action59 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 248 Action60 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 249 Action61 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 250 Action62 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 251 Action63 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 252 Action64 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 253 Action65 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 254 Action66 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 255 Action67 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 256 Action68 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 257 Action69 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 258 Action70 <- <{
		    p.AssembleSortedExpression()
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 259 Action71 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 260 Action72 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 261 Action73 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 262 Action74 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 263 Action75 <- <{
		    p.AssembleConditionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 264 Action76 <- <{
		    p.AssembleExpressionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 265 Action77 <- <{
		    p.AssembleWhenThenPair()
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 266 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 267 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 268 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 269 Action81 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 270 Action82 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 271 Action83 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 272 Action84 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 273 Action85 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 274 Action86 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 275 Action87 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 276 Action88 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 277 Action89 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 278 Action90 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 279 Action91 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 280 Action92 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 281 Action93 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 282 Action94 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 283 Action95 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 284 Action96 <- <{
		    p.PushComponent(begin, end, Wait)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 285 Action97 <- <{
		    p.PushComponent(begin, end, DropOldest)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 286 Action98 <- <{
		    p.PushComponent(begin, end, DropNewest)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 287 Action99 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 288 Action100 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 289 Action101 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 290 Action102 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 291 Action103 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 292 Action104 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 293 Action105 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 294 Action106 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 295 Action107 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 296 Action108 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 297 Action109 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 298 Action110 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction110, position)
			}
			return true
		},
		/* 299 Action111 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction111, position)
			}
			return true
		},
		/* 300 Action112 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction112, position)
			}
			return true
		},
		/* 301 Action113 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction113, position)
			}
			return true
		},
		/* 302 Action114 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction114, position)
			}
			return true
		},
		/* 303 Action115 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction115, position)
			}
			return true
		},
		/* 304 Action116 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction116, position)
			}
			return true
		},
		/* 305 Action117 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction117, position)
			}
			return true
		},
		/* 306 Action118 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction118, position)
			}
			return true
		},
		/* 307 Action119 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction119, position)
			}
			return true
		},
		/* 308 Action120 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction120, position)
			}
			return true
		},
		/* 309 Action121 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction121, position)
			}
			return true
		},
		/* 310 Action122 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction122, position)
			}
			return true
		},
		/* 311 Action123 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction123, position)
			}
			return true
		},
		/* 312 Action124 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction124, position)
			}
			return true
		},
		/* 313 Action125 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction125, position)
			}
			return true
		},
		/* 314 Action126 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction126, position)
			}
			return true
		},
		/* 315 Action127 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction127, position)
			}
			return true
		},
		/* 316 Action128 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction128, position)
			}
			return true
		},
		/* 317 Action129 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction129, position)
			}
			return true
		},
		/* 318 Action130 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction130, position)
			}
			return true
		},
		/* 319 Action131 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction131, position)
			}
			return true
		},
		/* 320 Action132 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction132, position)
			}
			return true
		},
		/* 321 Action133 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction133, position)
			}
			return true
		},
	}
	p.rules = _rules
}

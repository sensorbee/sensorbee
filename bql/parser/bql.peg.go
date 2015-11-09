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
	rules  [315]func() bool
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

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction76:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction77:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction78:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction80:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction81:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction82:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction83:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction84:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction85:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction86:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction87:

			p.PushComponent(begin, end, Istream)

		case ruleAction88:

			p.PushComponent(begin, end, Dstream)

		case ruleAction89:

			p.PushComponent(begin, end, Rstream)

		case ruleAction90:

			p.PushComponent(begin, end, Tuples)

		case ruleAction91:

			p.PushComponent(begin, end, Seconds)

		case ruleAction92:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction93:

			p.PushComponent(begin, end, Wait)

		case ruleAction94:

			p.PushComponent(begin, end, DropOldest)

		case ruleAction95:

			p.PushComponent(begin, end, DropNewest)

		case ruleAction96:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction97:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction98:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction99:

			p.PushComponent(begin, end, Yes)

		case ruleAction100:

			p.PushComponent(begin, end, No)

		case ruleAction101:

			p.PushComponent(begin, end, Yes)

		case ruleAction102:

			p.PushComponent(begin, end, No)

		case ruleAction103:

			p.PushComponent(begin, end, Bool)

		case ruleAction104:

			p.PushComponent(begin, end, Int)

		case ruleAction105:

			p.PushComponent(begin, end, Float)

		case ruleAction106:

			p.PushComponent(begin, end, String)

		case ruleAction107:

			p.PushComponent(begin, end, Blob)

		case ruleAction108:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction109:

			p.PushComponent(begin, end, Array)

		case ruleAction110:

			p.PushComponent(begin, end, Map)

		case ruleAction111:

			p.PushComponent(begin, end, Or)

		case ruleAction112:

			p.PushComponent(begin, end, And)

		case ruleAction113:

			p.PushComponent(begin, end, Not)

		case ruleAction114:

			p.PushComponent(begin, end, Equal)

		case ruleAction115:

			p.PushComponent(begin, end, Less)

		case ruleAction116:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction117:

			p.PushComponent(begin, end, Greater)

		case ruleAction118:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction119:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction120:

			p.PushComponent(begin, end, Concat)

		case ruleAction121:

			p.PushComponent(begin, end, Is)

		case ruleAction122:

			p.PushComponent(begin, end, IsNot)

		case ruleAction123:

			p.PushComponent(begin, end, Plus)

		case ruleAction124:

			p.PushComponent(begin, end, Minus)

		case ruleAction125:

			p.PushComponent(begin, end, Multiply)

		case ruleAction126:

			p.PushComponent(begin, end, Divide)

		case ruleAction127:

			p.PushComponent(begin, end, Modulo)

		case ruleAction128:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction129:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction130:

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
		/* 28 LoadStateOrCreateStmt <- <(LoadStateStmt sp (('o' / 'O') ('r' / 'R')) sp (('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp (('i' / 'I') ('f' / 'F')) sp (('n' / 'N') ('o' / 'O') ('t' / 'T')) sp (('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S')) SourceSinkSpecs Action22)> */
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
					if buffer[position] != rune('e') {
						goto l596
					}
					position++
					goto l595
				l596:
					position, tokenIndex, depth = position595, tokenIndex595, depth595
					if buffer[position] != rune('E') {
						goto l567
					}
					position++
				}
			l595:
				{
					position597, tokenIndex597, depth597 := position, tokenIndex, depth
					if buffer[position] != rune('x') {
						goto l598
					}
					position++
					goto l597
				l598:
					position, tokenIndex, depth = position597, tokenIndex597, depth597
					if buffer[position] != rune('X') {
						goto l567
					}
					position++
				}
			l597:
				{
					position599, tokenIndex599, depth599 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l600
					}
					position++
					goto l599
				l600:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if buffer[position] != rune('I') {
						goto l567
					}
					position++
				}
			l599:
				{
					position601, tokenIndex601, depth601 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l602
					}
					position++
					goto l601
				l602:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if buffer[position] != rune('S') {
						goto l567
					}
					position++
				}
			l601:
				{
					position603, tokenIndex603, depth603 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l604
					}
					position++
					goto l603
				l604:
					position, tokenIndex, depth = position603, tokenIndex603, depth603
					if buffer[position] != rune('T') {
						goto l567
					}
					position++
				}
			l603:
				{
					position605, tokenIndex605, depth605 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l606
					}
					position++
					goto l605
				l606:
					position, tokenIndex, depth = position605, tokenIndex605, depth605
					if buffer[position] != rune('S') {
						goto l567
					}
					position++
				}
			l605:
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
			position607, tokenIndex607, depth607 := position, tokenIndex, depth
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
						goto l607
					}
					position++
				}
			l609:
				{
					position611, tokenIndex611, depth611 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l612
					}
					position++
					goto l611
				l612:
					position, tokenIndex, depth = position611, tokenIndex611, depth611
					if buffer[position] != rune('A') {
						goto l607
					}
					position++
				}
			l611:
				{
					position613, tokenIndex613, depth613 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l614
					}
					position++
					goto l613
				l614:
					position, tokenIndex, depth = position613, tokenIndex613, depth613
					if buffer[position] != rune('V') {
						goto l607
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
						goto l607
					}
					position++
				}
			l615:
				if !_rules[rulesp]() {
					goto l607
				}
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
						goto l607
					}
					position++
				}
			l617:
				{
					position619, tokenIndex619, depth619 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l620
					}
					position++
					goto l619
				l620:
					position, tokenIndex, depth = position619, tokenIndex619, depth619
					if buffer[position] != rune('T') {
						goto l607
					}
					position++
				}
			l619:
				{
					position621, tokenIndex621, depth621 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l622
					}
					position++
					goto l621
				l622:
					position, tokenIndex, depth = position621, tokenIndex621, depth621
					if buffer[position] != rune('A') {
						goto l607
					}
					position++
				}
			l621:
				{
					position623, tokenIndex623, depth623 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l624
					}
					position++
					goto l623
				l624:
					position, tokenIndex, depth = position623, tokenIndex623, depth623
					if buffer[position] != rune('T') {
						goto l607
					}
					position++
				}
			l623:
				{
					position625, tokenIndex625, depth625 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l626
					}
					position++
					goto l625
				l626:
					position, tokenIndex, depth = position625, tokenIndex625, depth625
					if buffer[position] != rune('E') {
						goto l607
					}
					position++
				}
			l625:
				if !_rules[rulesp]() {
					goto l607
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l607
				}
				if !_rules[ruleStateTagOpt]() {
					goto l607
				}
				if !_rules[ruleAction23]() {
					goto l607
				}
				depth--
				add(ruleSaveStateStmt, position608)
			}
			return true
		l607:
			position, tokenIndex, depth = position607, tokenIndex607, depth607
			return false
		},
		/* 30 EvalStmt <- <(('e' / 'E') ('v' / 'V') ('a' / 'A') ('l' / 'L') sp Expression <(sp (('o' / 'O') ('n' / 'N')) sp MapExpr)?> Action24)> */
		func() bool {
			position627, tokenIndex627, depth627 := position, tokenIndex, depth
			{
				position628 := position
				depth++
				{
					position629, tokenIndex629, depth629 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l630
					}
					position++
					goto l629
				l630:
					position, tokenIndex, depth = position629, tokenIndex629, depth629
					if buffer[position] != rune('E') {
						goto l627
					}
					position++
				}
			l629:
				{
					position631, tokenIndex631, depth631 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l632
					}
					position++
					goto l631
				l632:
					position, tokenIndex, depth = position631, tokenIndex631, depth631
					if buffer[position] != rune('V') {
						goto l627
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
						goto l627
					}
					position++
				}
			l633:
				{
					position635, tokenIndex635, depth635 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l636
					}
					position++
					goto l635
				l636:
					position, tokenIndex, depth = position635, tokenIndex635, depth635
					if buffer[position] != rune('L') {
						goto l627
					}
					position++
				}
			l635:
				if !_rules[rulesp]() {
					goto l627
				}
				if !_rules[ruleExpression]() {
					goto l627
				}
				{
					position637 := position
					depth++
					{
						position638, tokenIndex638, depth638 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l638
						}
						{
							position640, tokenIndex640, depth640 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l641
							}
							position++
							goto l640
						l641:
							position, tokenIndex, depth = position640, tokenIndex640, depth640
							if buffer[position] != rune('O') {
								goto l638
							}
							position++
						}
					l640:
						{
							position642, tokenIndex642, depth642 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l643
							}
							position++
							goto l642
						l643:
							position, tokenIndex, depth = position642, tokenIndex642, depth642
							if buffer[position] != rune('N') {
								goto l638
							}
							position++
						}
					l642:
						if !_rules[rulesp]() {
							goto l638
						}
						if !_rules[ruleMapExpr]() {
							goto l638
						}
						goto l639
					l638:
						position, tokenIndex, depth = position638, tokenIndex638, depth638
					}
				l639:
					depth--
					add(rulePegText, position637)
				}
				if !_rules[ruleAction24]() {
					goto l627
				}
				depth--
				add(ruleEvalStmt, position628)
			}
			return true
		l627:
			position, tokenIndex, depth = position627, tokenIndex627, depth627
			return false
		},
		/* 31 Emitter <- <(sp (ISTREAM / DSTREAM / RSTREAM) EmitterOptions Action25)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				if !_rules[rulesp]() {
					goto l644
				}
				{
					position646, tokenIndex646, depth646 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l647
					}
					goto l646
				l647:
					position, tokenIndex, depth = position646, tokenIndex646, depth646
					if !_rules[ruleDSTREAM]() {
						goto l648
					}
					goto l646
				l648:
					position, tokenIndex, depth = position646, tokenIndex646, depth646
					if !_rules[ruleRSTREAM]() {
						goto l644
					}
				}
			l646:
				if !_rules[ruleEmitterOptions]() {
					goto l644
				}
				if !_rules[ruleAction25]() {
					goto l644
				}
				depth--
				add(ruleEmitter, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 32 EmitterOptions <- <(<(spOpt '[' spOpt EmitterOptionCombinations spOpt ']')?> Action26)> */
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
						if !_rules[rulespOpt]() {
							goto l652
						}
						if buffer[position] != rune('[') {
							goto l652
						}
						position++
						if !_rules[rulespOpt]() {
							goto l652
						}
						if !_rules[ruleEmitterOptionCombinations]() {
							goto l652
						}
						if !_rules[rulespOpt]() {
							goto l652
						}
						if buffer[position] != rune(']') {
							goto l652
						}
						position++
						goto l653
					l652:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
					}
				l653:
					depth--
					add(rulePegText, position651)
				}
				if !_rules[ruleAction26]() {
					goto l649
				}
				depth--
				add(ruleEmitterOptions, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 33 EmitterOptionCombinations <- <(EmitterLimit / (EmitterSample sp EmitterLimit) / EmitterSample)> */
		func() bool {
			position654, tokenIndex654, depth654 := position, tokenIndex, depth
			{
				position655 := position
				depth++
				{
					position656, tokenIndex656, depth656 := position, tokenIndex, depth
					if !_rules[ruleEmitterLimit]() {
						goto l657
					}
					goto l656
				l657:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleEmitterSample]() {
						goto l658
					}
					if !_rules[rulesp]() {
						goto l658
					}
					if !_rules[ruleEmitterLimit]() {
						goto l658
					}
					goto l656
				l658:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if !_rules[ruleEmitterSample]() {
						goto l654
					}
				}
			l656:
				depth--
				add(ruleEmitterOptionCombinations, position655)
			}
			return true
		l654:
			position, tokenIndex, depth = position654, tokenIndex654, depth654
			return false
		},
		/* 34 EmitterLimit <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') sp NumericLiteral Action27)> */
		func() bool {
			position659, tokenIndex659, depth659 := position, tokenIndex, depth
			{
				position660 := position
				depth++
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l662
					}
					position++
					goto l661
				l662:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
					if buffer[position] != rune('L') {
						goto l659
					}
					position++
				}
			l661:
				{
					position663, tokenIndex663, depth663 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l664
					}
					position++
					goto l663
				l664:
					position, tokenIndex, depth = position663, tokenIndex663, depth663
					if buffer[position] != rune('I') {
						goto l659
					}
					position++
				}
			l663:
				{
					position665, tokenIndex665, depth665 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l666
					}
					position++
					goto l665
				l666:
					position, tokenIndex, depth = position665, tokenIndex665, depth665
					if buffer[position] != rune('M') {
						goto l659
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
						goto l659
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
						goto l659
					}
					position++
				}
			l669:
				if !_rules[rulesp]() {
					goto l659
				}
				if !_rules[ruleNumericLiteral]() {
					goto l659
				}
				if !_rules[ruleAction27]() {
					goto l659
				}
				depth--
				add(ruleEmitterLimit, position660)
			}
			return true
		l659:
			position, tokenIndex, depth = position659, tokenIndex659, depth659
			return false
		},
		/* 35 EmitterSample <- <(CountBasedSampling / RandomizedSampling / TimeBasedSampling)> */
		func() bool {
			position671, tokenIndex671, depth671 := position, tokenIndex, depth
			{
				position672 := position
				depth++
				{
					position673, tokenIndex673, depth673 := position, tokenIndex, depth
					if !_rules[ruleCountBasedSampling]() {
						goto l674
					}
					goto l673
				l674:
					position, tokenIndex, depth = position673, tokenIndex673, depth673
					if !_rules[ruleRandomizedSampling]() {
						goto l675
					}
					goto l673
				l675:
					position, tokenIndex, depth = position673, tokenIndex673, depth673
					if !_rules[ruleTimeBasedSampling]() {
						goto l671
					}
				}
			l673:
				depth--
				add(ruleEmitterSample, position672)
			}
			return true
		l671:
			position, tokenIndex, depth = position671, tokenIndex671, depth671
			return false
		},
		/* 36 CountBasedSampling <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp NumericLiteral spOpt '-'? spOpt ((('s' / 'S') ('t' / 'T')) / (('n' / 'N') ('d' / 'D')) / (('r' / 'R') ('d' / 'D')) / (('t' / 'T') ('h' / 'H'))) sp (('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E')) Action28)> */
		func() bool {
			position676, tokenIndex676, depth676 := position, tokenIndex, depth
			{
				position677 := position
				depth++
				{
					position678, tokenIndex678, depth678 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l679
					}
					position++
					goto l678
				l679:
					position, tokenIndex, depth = position678, tokenIndex678, depth678
					if buffer[position] != rune('E') {
						goto l676
					}
					position++
				}
			l678:
				{
					position680, tokenIndex680, depth680 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l681
					}
					position++
					goto l680
				l681:
					position, tokenIndex, depth = position680, tokenIndex680, depth680
					if buffer[position] != rune('V') {
						goto l676
					}
					position++
				}
			l680:
				{
					position682, tokenIndex682, depth682 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l683
					}
					position++
					goto l682
				l683:
					position, tokenIndex, depth = position682, tokenIndex682, depth682
					if buffer[position] != rune('E') {
						goto l676
					}
					position++
				}
			l682:
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
						goto l676
					}
					position++
				}
			l684:
				{
					position686, tokenIndex686, depth686 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l687
					}
					position++
					goto l686
				l687:
					position, tokenIndex, depth = position686, tokenIndex686, depth686
					if buffer[position] != rune('Y') {
						goto l676
					}
					position++
				}
			l686:
				if !_rules[rulesp]() {
					goto l676
				}
				if !_rules[ruleNumericLiteral]() {
					goto l676
				}
				if !_rules[rulespOpt]() {
					goto l676
				}
				{
					position688, tokenIndex688, depth688 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l688
					}
					position++
					goto l689
				l688:
					position, tokenIndex, depth = position688, tokenIndex688, depth688
				}
			l689:
				if !_rules[rulespOpt]() {
					goto l676
				}
				{
					position690, tokenIndex690, depth690 := position, tokenIndex, depth
					{
						position692, tokenIndex692, depth692 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l693
						}
						position++
						goto l692
					l693:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						if buffer[position] != rune('S') {
							goto l691
						}
						position++
					}
				l692:
					{
						position694, tokenIndex694, depth694 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l695
						}
						position++
						goto l694
					l695:
						position, tokenIndex, depth = position694, tokenIndex694, depth694
						if buffer[position] != rune('T') {
							goto l691
						}
						position++
					}
				l694:
					goto l690
				l691:
					position, tokenIndex, depth = position690, tokenIndex690, depth690
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('N') {
							goto l696
						}
						position++
					}
				l697:
					{
						position699, tokenIndex699, depth699 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l700
						}
						position++
						goto l699
					l700:
						position, tokenIndex, depth = position699, tokenIndex699, depth699
						if buffer[position] != rune('D') {
							goto l696
						}
						position++
					}
				l699:
					goto l690
				l696:
					position, tokenIndex, depth = position690, tokenIndex690, depth690
					{
						position702, tokenIndex702, depth702 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l703
						}
						position++
						goto l702
					l703:
						position, tokenIndex, depth = position702, tokenIndex702, depth702
						if buffer[position] != rune('R') {
							goto l701
						}
						position++
					}
				l702:
					{
						position704, tokenIndex704, depth704 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l705
						}
						position++
						goto l704
					l705:
						position, tokenIndex, depth = position704, tokenIndex704, depth704
						if buffer[position] != rune('D') {
							goto l701
						}
						position++
					}
				l704:
					goto l690
				l701:
					position, tokenIndex, depth = position690, tokenIndex690, depth690
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
							goto l676
						}
						position++
					}
				l706:
					{
						position708, tokenIndex708, depth708 := position, tokenIndex, depth
						if buffer[position] != rune('h') {
							goto l709
						}
						position++
						goto l708
					l709:
						position, tokenIndex, depth = position708, tokenIndex708, depth708
						if buffer[position] != rune('H') {
							goto l676
						}
						position++
					}
				l708:
				}
			l690:
				if !_rules[rulesp]() {
					goto l676
				}
				{
					position710, tokenIndex710, depth710 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l711
					}
					position++
					goto l710
				l711:
					position, tokenIndex, depth = position710, tokenIndex710, depth710
					if buffer[position] != rune('T') {
						goto l676
					}
					position++
				}
			l710:
				{
					position712, tokenIndex712, depth712 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l713
					}
					position++
					goto l712
				l713:
					position, tokenIndex, depth = position712, tokenIndex712, depth712
					if buffer[position] != rune('U') {
						goto l676
					}
					position++
				}
			l712:
				{
					position714, tokenIndex714, depth714 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l715
					}
					position++
					goto l714
				l715:
					position, tokenIndex, depth = position714, tokenIndex714, depth714
					if buffer[position] != rune('P') {
						goto l676
					}
					position++
				}
			l714:
				{
					position716, tokenIndex716, depth716 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l717
					}
					position++
					goto l716
				l717:
					position, tokenIndex, depth = position716, tokenIndex716, depth716
					if buffer[position] != rune('L') {
						goto l676
					}
					position++
				}
			l716:
				{
					position718, tokenIndex718, depth718 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l719
					}
					position++
					goto l718
				l719:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if buffer[position] != rune('E') {
						goto l676
					}
					position++
				}
			l718:
				if !_rules[ruleAction28]() {
					goto l676
				}
				depth--
				add(ruleCountBasedSampling, position677)
			}
			return true
		l676:
			position, tokenIndex, depth = position676, tokenIndex676, depth676
			return false
		},
		/* 37 RandomizedSampling <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') sp (FloatLiteral / NumericLiteral) spOpt '%' Action29)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722, tokenIndex722, depth722 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l723
					}
					position++
					goto l722
				l723:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if buffer[position] != rune('S') {
						goto l720
					}
					position++
				}
			l722:
				{
					position724, tokenIndex724, depth724 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l725
					}
					position++
					goto l724
				l725:
					position, tokenIndex, depth = position724, tokenIndex724, depth724
					if buffer[position] != rune('A') {
						goto l720
					}
					position++
				}
			l724:
				{
					position726, tokenIndex726, depth726 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l727
					}
					position++
					goto l726
				l727:
					position, tokenIndex, depth = position726, tokenIndex726, depth726
					if buffer[position] != rune('M') {
						goto l720
					}
					position++
				}
			l726:
				{
					position728, tokenIndex728, depth728 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l729
					}
					position++
					goto l728
				l729:
					position, tokenIndex, depth = position728, tokenIndex728, depth728
					if buffer[position] != rune('P') {
						goto l720
					}
					position++
				}
			l728:
				{
					position730, tokenIndex730, depth730 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l731
					}
					position++
					goto l730
				l731:
					position, tokenIndex, depth = position730, tokenIndex730, depth730
					if buffer[position] != rune('L') {
						goto l720
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
						goto l720
					}
					position++
				}
			l732:
				if !_rules[rulesp]() {
					goto l720
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
						goto l720
					}
				}
			l734:
				if !_rules[rulespOpt]() {
					goto l720
				}
				if buffer[position] != rune('%') {
					goto l720
				}
				position++
				if !_rules[ruleAction29]() {
					goto l720
				}
				depth--
				add(ruleRandomizedSampling, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 38 TimeBasedSampling <- <(TimeBasedSamplingSeconds / TimeBasedSamplingMilliseconds)> */
		func() bool {
			position736, tokenIndex736, depth736 := position, tokenIndex, depth
			{
				position737 := position
				depth++
				{
					position738, tokenIndex738, depth738 := position, tokenIndex, depth
					if !_rules[ruleTimeBasedSamplingSeconds]() {
						goto l739
					}
					goto l738
				l739:
					position, tokenIndex, depth = position738, tokenIndex738, depth738
					if !_rules[ruleTimeBasedSamplingMilliseconds]() {
						goto l736
					}
				}
			l738:
				depth--
				add(ruleTimeBasedSampling, position737)
			}
			return true
		l736:
			position, tokenIndex, depth = position736, tokenIndex736, depth736
			return false
		},
		/* 39 TimeBasedSamplingSeconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action30)> */
		func() bool {
			position740, tokenIndex740, depth740 := position, tokenIndex, depth
			{
				position741 := position
				depth++
				{
					position742, tokenIndex742, depth742 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l743
					}
					position++
					goto l742
				l743:
					position, tokenIndex, depth = position742, tokenIndex742, depth742
					if buffer[position] != rune('E') {
						goto l740
					}
					position++
				}
			l742:
				{
					position744, tokenIndex744, depth744 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l745
					}
					position++
					goto l744
				l745:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
					if buffer[position] != rune('V') {
						goto l740
					}
					position++
				}
			l744:
				{
					position746, tokenIndex746, depth746 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l747
					}
					position++
					goto l746
				l747:
					position, tokenIndex, depth = position746, tokenIndex746, depth746
					if buffer[position] != rune('E') {
						goto l740
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
						goto l740
					}
					position++
				}
			l748:
				{
					position750, tokenIndex750, depth750 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l751
					}
					position++
					goto l750
				l751:
					position, tokenIndex, depth = position750, tokenIndex750, depth750
					if buffer[position] != rune('Y') {
						goto l740
					}
					position++
				}
			l750:
				if !_rules[rulesp]() {
					goto l740
				}
				{
					position752, tokenIndex752, depth752 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l753
					}
					goto l752
				l753:
					position, tokenIndex, depth = position752, tokenIndex752, depth752
					if !_rules[ruleNumericLiteral]() {
						goto l740
					}
				}
			l752:
				if !_rules[rulesp]() {
					goto l740
				}
				{
					position754, tokenIndex754, depth754 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l755
					}
					position++
					goto l754
				l755:
					position, tokenIndex, depth = position754, tokenIndex754, depth754
					if buffer[position] != rune('S') {
						goto l740
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
						goto l740
					}
					position++
				}
			l756:
				{
					position758, tokenIndex758, depth758 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l759
					}
					position++
					goto l758
				l759:
					position, tokenIndex, depth = position758, tokenIndex758, depth758
					if buffer[position] != rune('C') {
						goto l740
					}
					position++
				}
			l758:
				{
					position760, tokenIndex760, depth760 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l761
					}
					position++
					goto l760
				l761:
					position, tokenIndex, depth = position760, tokenIndex760, depth760
					if buffer[position] != rune('O') {
						goto l740
					}
					position++
				}
			l760:
				{
					position762, tokenIndex762, depth762 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l763
					}
					position++
					goto l762
				l763:
					position, tokenIndex, depth = position762, tokenIndex762, depth762
					if buffer[position] != rune('N') {
						goto l740
					}
					position++
				}
			l762:
				{
					position764, tokenIndex764, depth764 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l765
					}
					position++
					goto l764
				l765:
					position, tokenIndex, depth = position764, tokenIndex764, depth764
					if buffer[position] != rune('D') {
						goto l740
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
						goto l740
					}
					position++
				}
			l766:
				if !_rules[ruleAction30]() {
					goto l740
				}
				depth--
				add(ruleTimeBasedSamplingSeconds, position741)
			}
			return true
		l740:
			position, tokenIndex, depth = position740, tokenIndex740, depth740
			return false
		},
		/* 40 TimeBasedSamplingMilliseconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action31)> */
		func() bool {
			position768, tokenIndex768, depth768 := position, tokenIndex, depth
			{
				position769 := position
				depth++
				{
					position770, tokenIndex770, depth770 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l771
					}
					position++
					goto l770
				l771:
					position, tokenIndex, depth = position770, tokenIndex770, depth770
					if buffer[position] != rune('E') {
						goto l768
					}
					position++
				}
			l770:
				{
					position772, tokenIndex772, depth772 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l773
					}
					position++
					goto l772
				l773:
					position, tokenIndex, depth = position772, tokenIndex772, depth772
					if buffer[position] != rune('V') {
						goto l768
					}
					position++
				}
			l772:
				{
					position774, tokenIndex774, depth774 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l775
					}
					position++
					goto l774
				l775:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if buffer[position] != rune('E') {
						goto l768
					}
					position++
				}
			l774:
				{
					position776, tokenIndex776, depth776 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l777
					}
					position++
					goto l776
				l777:
					position, tokenIndex, depth = position776, tokenIndex776, depth776
					if buffer[position] != rune('R') {
						goto l768
					}
					position++
				}
			l776:
				{
					position778, tokenIndex778, depth778 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l779
					}
					position++
					goto l778
				l779:
					position, tokenIndex, depth = position778, tokenIndex778, depth778
					if buffer[position] != rune('Y') {
						goto l768
					}
					position++
				}
			l778:
				if !_rules[rulesp]() {
					goto l768
				}
				{
					position780, tokenIndex780, depth780 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l781
					}
					goto l780
				l781:
					position, tokenIndex, depth = position780, tokenIndex780, depth780
					if !_rules[ruleNumericLiteral]() {
						goto l768
					}
				}
			l780:
				if !_rules[rulesp]() {
					goto l768
				}
				{
					position782, tokenIndex782, depth782 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l783
					}
					position++
					goto l782
				l783:
					position, tokenIndex, depth = position782, tokenIndex782, depth782
					if buffer[position] != rune('M') {
						goto l768
					}
					position++
				}
			l782:
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
						goto l768
					}
					position++
				}
			l784:
				{
					position786, tokenIndex786, depth786 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l787
					}
					position++
					goto l786
				l787:
					position, tokenIndex, depth = position786, tokenIndex786, depth786
					if buffer[position] != rune('L') {
						goto l768
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
						goto l768
					}
					position++
				}
			l788:
				{
					position790, tokenIndex790, depth790 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l791
					}
					position++
					goto l790
				l791:
					position, tokenIndex, depth = position790, tokenIndex790, depth790
					if buffer[position] != rune('I') {
						goto l768
					}
					position++
				}
			l790:
				{
					position792, tokenIndex792, depth792 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l793
					}
					position++
					goto l792
				l793:
					position, tokenIndex, depth = position792, tokenIndex792, depth792
					if buffer[position] != rune('S') {
						goto l768
					}
					position++
				}
			l792:
				{
					position794, tokenIndex794, depth794 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l795
					}
					position++
					goto l794
				l795:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
					if buffer[position] != rune('E') {
						goto l768
					}
					position++
				}
			l794:
				{
					position796, tokenIndex796, depth796 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l797
					}
					position++
					goto l796
				l797:
					position, tokenIndex, depth = position796, tokenIndex796, depth796
					if buffer[position] != rune('C') {
						goto l768
					}
					position++
				}
			l796:
				{
					position798, tokenIndex798, depth798 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l799
					}
					position++
					goto l798
				l799:
					position, tokenIndex, depth = position798, tokenIndex798, depth798
					if buffer[position] != rune('O') {
						goto l768
					}
					position++
				}
			l798:
				{
					position800, tokenIndex800, depth800 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l801
					}
					position++
					goto l800
				l801:
					position, tokenIndex, depth = position800, tokenIndex800, depth800
					if buffer[position] != rune('N') {
						goto l768
					}
					position++
				}
			l800:
				{
					position802, tokenIndex802, depth802 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l803
					}
					position++
					goto l802
				l803:
					position, tokenIndex, depth = position802, tokenIndex802, depth802
					if buffer[position] != rune('D') {
						goto l768
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
						goto l768
					}
					position++
				}
			l804:
				if !_rules[ruleAction31]() {
					goto l768
				}
				depth--
				add(ruleTimeBasedSamplingMilliseconds, position769)
			}
			return true
		l768:
			position, tokenIndex, depth = position768, tokenIndex768, depth768
			return false
		},
		/* 41 Projections <- <(<(sp Projection (spOpt ',' spOpt Projection)*)> Action32)> */
		func() bool {
			position806, tokenIndex806, depth806 := position, tokenIndex, depth
			{
				position807 := position
				depth++
				{
					position808 := position
					depth++
					if !_rules[rulesp]() {
						goto l806
					}
					if !_rules[ruleProjection]() {
						goto l806
					}
				l809:
					{
						position810, tokenIndex810, depth810 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l810
						}
						if buffer[position] != rune(',') {
							goto l810
						}
						position++
						if !_rules[rulespOpt]() {
							goto l810
						}
						if !_rules[ruleProjection]() {
							goto l810
						}
						goto l809
					l810:
						position, tokenIndex, depth = position810, tokenIndex810, depth810
					}
					depth--
					add(rulePegText, position808)
				}
				if !_rules[ruleAction32]() {
					goto l806
				}
				depth--
				add(ruleProjections, position807)
			}
			return true
		l806:
			position, tokenIndex, depth = position806, tokenIndex806, depth806
			return false
		},
		/* 42 Projection <- <(AliasExpression / ExpressionOrWildcard)> */
		func() bool {
			position811, tokenIndex811, depth811 := position, tokenIndex, depth
			{
				position812 := position
				depth++
				{
					position813, tokenIndex813, depth813 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l814
					}
					goto l813
				l814:
					position, tokenIndex, depth = position813, tokenIndex813, depth813
					if !_rules[ruleExpressionOrWildcard]() {
						goto l811
					}
				}
			l813:
				depth--
				add(ruleProjection, position812)
			}
			return true
		l811:
			position, tokenIndex, depth = position811, tokenIndex811, depth811
			return false
		},
		/* 43 AliasExpression <- <(ExpressionOrWildcard sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action33)> */
		func() bool {
			position815, tokenIndex815, depth815 := position, tokenIndex, depth
			{
				position816 := position
				depth++
				if !_rules[ruleExpressionOrWildcard]() {
					goto l815
				}
				if !_rules[rulesp]() {
					goto l815
				}
				{
					position817, tokenIndex817, depth817 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l818
					}
					position++
					goto l817
				l818:
					position, tokenIndex, depth = position817, tokenIndex817, depth817
					if buffer[position] != rune('A') {
						goto l815
					}
					position++
				}
			l817:
				{
					position819, tokenIndex819, depth819 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l820
					}
					position++
					goto l819
				l820:
					position, tokenIndex, depth = position819, tokenIndex819, depth819
					if buffer[position] != rune('S') {
						goto l815
					}
					position++
				}
			l819:
				if !_rules[rulesp]() {
					goto l815
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l815
				}
				if !_rules[ruleAction33]() {
					goto l815
				}
				depth--
				add(ruleAliasExpression, position816)
			}
			return true
		l815:
			position, tokenIndex, depth = position815, tokenIndex815, depth815
			return false
		},
		/* 44 WindowedFrom <- <(<(sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp Relations)?> Action34)> */
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
						if !_rules[rulesp]() {
							goto l824
						}
						{
							position826, tokenIndex826, depth826 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l827
							}
							position++
							goto l826
						l827:
							position, tokenIndex, depth = position826, tokenIndex826, depth826
							if buffer[position] != rune('F') {
								goto l824
							}
							position++
						}
					l826:
						{
							position828, tokenIndex828, depth828 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l829
							}
							position++
							goto l828
						l829:
							position, tokenIndex, depth = position828, tokenIndex828, depth828
							if buffer[position] != rune('R') {
								goto l824
							}
							position++
						}
					l828:
						{
							position830, tokenIndex830, depth830 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l831
							}
							position++
							goto l830
						l831:
							position, tokenIndex, depth = position830, tokenIndex830, depth830
							if buffer[position] != rune('O') {
								goto l824
							}
							position++
						}
					l830:
						{
							position832, tokenIndex832, depth832 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l833
							}
							position++
							goto l832
						l833:
							position, tokenIndex, depth = position832, tokenIndex832, depth832
							if buffer[position] != rune('M') {
								goto l824
							}
							position++
						}
					l832:
						if !_rules[rulesp]() {
							goto l824
						}
						if !_rules[ruleRelations]() {
							goto l824
						}
						goto l825
					l824:
						position, tokenIndex, depth = position824, tokenIndex824, depth824
					}
				l825:
					depth--
					add(rulePegText, position823)
				}
				if !_rules[ruleAction34]() {
					goto l821
				}
				depth--
				add(ruleWindowedFrom, position822)
			}
			return true
		l821:
			position, tokenIndex, depth = position821, tokenIndex821, depth821
			return false
		},
		/* 45 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position834, tokenIndex834, depth834 := position, tokenIndex, depth
			{
				position835 := position
				depth++
				{
					position836, tokenIndex836, depth836 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l837
					}
					goto l836
				l837:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleTuplesInterval]() {
						goto l834
					}
				}
			l836:
				depth--
				add(ruleInterval, position835)
			}
			return true
		l834:
			position, tokenIndex, depth = position834, tokenIndex834, depth834
			return false
		},
		/* 46 TimeInterval <- <((FloatLiteral / NumericLiteral) sp (SECONDS / MILLISECONDS) Action35)> */
		func() bool {
			position838, tokenIndex838, depth838 := position, tokenIndex, depth
			{
				position839 := position
				depth++
				{
					position840, tokenIndex840, depth840 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l841
					}
					goto l840
				l841:
					position, tokenIndex, depth = position840, tokenIndex840, depth840
					if !_rules[ruleNumericLiteral]() {
						goto l838
					}
				}
			l840:
				if !_rules[rulesp]() {
					goto l838
				}
				{
					position842, tokenIndex842, depth842 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l843
					}
					goto l842
				l843:
					position, tokenIndex, depth = position842, tokenIndex842, depth842
					if !_rules[ruleMILLISECONDS]() {
						goto l838
					}
				}
			l842:
				if !_rules[ruleAction35]() {
					goto l838
				}
				depth--
				add(ruleTimeInterval, position839)
			}
			return true
		l838:
			position, tokenIndex, depth = position838, tokenIndex838, depth838
			return false
		},
		/* 47 TuplesInterval <- <(NumericLiteral sp TUPLES Action36)> */
		func() bool {
			position844, tokenIndex844, depth844 := position, tokenIndex, depth
			{
				position845 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l844
				}
				if !_rules[rulesp]() {
					goto l844
				}
				if !_rules[ruleTUPLES]() {
					goto l844
				}
				if !_rules[ruleAction36]() {
					goto l844
				}
				depth--
				add(ruleTuplesInterval, position845)
			}
			return true
		l844:
			position, tokenIndex, depth = position844, tokenIndex844, depth844
			return false
		},
		/* 48 Relations <- <(RelationLike (spOpt ',' spOpt RelationLike)*)> */
		func() bool {
			position846, tokenIndex846, depth846 := position, tokenIndex, depth
			{
				position847 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l846
				}
			l848:
				{
					position849, tokenIndex849, depth849 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l849
					}
					if buffer[position] != rune(',') {
						goto l849
					}
					position++
					if !_rules[rulespOpt]() {
						goto l849
					}
					if !_rules[ruleRelationLike]() {
						goto l849
					}
					goto l848
				l849:
					position, tokenIndex, depth = position849, tokenIndex849, depth849
				}
				depth--
				add(ruleRelations, position847)
			}
			return true
		l846:
			position, tokenIndex, depth = position846, tokenIndex846, depth846
			return false
		},
		/* 49 Filter <- <(<(sp (('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E')) sp Expression)?> Action37)> */
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
						if !_rules[rulesp]() {
							goto l853
						}
						{
							position855, tokenIndex855, depth855 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l856
							}
							position++
							goto l855
						l856:
							position, tokenIndex, depth = position855, tokenIndex855, depth855
							if buffer[position] != rune('W') {
								goto l853
							}
							position++
						}
					l855:
						{
							position857, tokenIndex857, depth857 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l858
							}
							position++
							goto l857
						l858:
							position, tokenIndex, depth = position857, tokenIndex857, depth857
							if buffer[position] != rune('H') {
								goto l853
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
								goto l853
							}
							position++
						}
					l859:
						{
							position861, tokenIndex861, depth861 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l862
							}
							position++
							goto l861
						l862:
							position, tokenIndex, depth = position861, tokenIndex861, depth861
							if buffer[position] != rune('R') {
								goto l853
							}
							position++
						}
					l861:
						{
							position863, tokenIndex863, depth863 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l864
							}
							position++
							goto l863
						l864:
							position, tokenIndex, depth = position863, tokenIndex863, depth863
							if buffer[position] != rune('E') {
								goto l853
							}
							position++
						}
					l863:
						if !_rules[rulesp]() {
							goto l853
						}
						if !_rules[ruleExpression]() {
							goto l853
						}
						goto l854
					l853:
						position, tokenIndex, depth = position853, tokenIndex853, depth853
					}
				l854:
					depth--
					add(rulePegText, position852)
				}
				if !_rules[ruleAction37]() {
					goto l850
				}
				depth--
				add(ruleFilter, position851)
			}
			return true
		l850:
			position, tokenIndex, depth = position850, tokenIndex850, depth850
			return false
		},
		/* 50 Grouping <- <(<(sp (('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P')) sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action38)> */
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
						if !_rules[rulesp]() {
							goto l868
						}
						{
							position870, tokenIndex870, depth870 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l871
							}
							position++
							goto l870
						l871:
							position, tokenIndex, depth = position870, tokenIndex870, depth870
							if buffer[position] != rune('G') {
								goto l868
							}
							position++
						}
					l870:
						{
							position872, tokenIndex872, depth872 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l873
							}
							position++
							goto l872
						l873:
							position, tokenIndex, depth = position872, tokenIndex872, depth872
							if buffer[position] != rune('R') {
								goto l868
							}
							position++
						}
					l872:
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
								goto l868
							}
							position++
						}
					l874:
						{
							position876, tokenIndex876, depth876 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l877
							}
							position++
							goto l876
						l877:
							position, tokenIndex, depth = position876, tokenIndex876, depth876
							if buffer[position] != rune('U') {
								goto l868
							}
							position++
						}
					l876:
						{
							position878, tokenIndex878, depth878 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l879
							}
							position++
							goto l878
						l879:
							position, tokenIndex, depth = position878, tokenIndex878, depth878
							if buffer[position] != rune('P') {
								goto l868
							}
							position++
						}
					l878:
						if !_rules[rulesp]() {
							goto l868
						}
						{
							position880, tokenIndex880, depth880 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l881
							}
							position++
							goto l880
						l881:
							position, tokenIndex, depth = position880, tokenIndex880, depth880
							if buffer[position] != rune('B') {
								goto l868
							}
							position++
						}
					l880:
						{
							position882, tokenIndex882, depth882 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l883
							}
							position++
							goto l882
						l883:
							position, tokenIndex, depth = position882, tokenIndex882, depth882
							if buffer[position] != rune('Y') {
								goto l868
							}
							position++
						}
					l882:
						if !_rules[rulesp]() {
							goto l868
						}
						if !_rules[ruleGroupList]() {
							goto l868
						}
						goto l869
					l868:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
					}
				l869:
					depth--
					add(rulePegText, position867)
				}
				if !_rules[ruleAction38]() {
					goto l865
				}
				depth--
				add(ruleGrouping, position866)
			}
			return true
		l865:
			position, tokenIndex, depth = position865, tokenIndex865, depth865
			return false
		},
		/* 51 GroupList <- <(Expression (spOpt ',' spOpt Expression)*)> */
		func() bool {
			position884, tokenIndex884, depth884 := position, tokenIndex, depth
			{
				position885 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l884
				}
			l886:
				{
					position887, tokenIndex887, depth887 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l887
					}
					if buffer[position] != rune(',') {
						goto l887
					}
					position++
					if !_rules[rulespOpt]() {
						goto l887
					}
					if !_rules[ruleExpression]() {
						goto l887
					}
					goto l886
				l887:
					position, tokenIndex, depth = position887, tokenIndex887, depth887
				}
				depth--
				add(ruleGroupList, position885)
			}
			return true
		l884:
			position, tokenIndex, depth = position884, tokenIndex884, depth884
			return false
		},
		/* 52 Having <- <(<(sp (('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G')) sp Expression)?> Action39)> */
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
						if !_rules[rulesp]() {
							goto l891
						}
						{
							position893, tokenIndex893, depth893 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l894
							}
							position++
							goto l893
						l894:
							position, tokenIndex, depth = position893, tokenIndex893, depth893
							if buffer[position] != rune('H') {
								goto l891
							}
							position++
						}
					l893:
						{
							position895, tokenIndex895, depth895 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l896
							}
							position++
							goto l895
						l896:
							position, tokenIndex, depth = position895, tokenIndex895, depth895
							if buffer[position] != rune('A') {
								goto l891
							}
							position++
						}
					l895:
						{
							position897, tokenIndex897, depth897 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l898
							}
							position++
							goto l897
						l898:
							position, tokenIndex, depth = position897, tokenIndex897, depth897
							if buffer[position] != rune('V') {
								goto l891
							}
							position++
						}
					l897:
						{
							position899, tokenIndex899, depth899 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l900
							}
							position++
							goto l899
						l900:
							position, tokenIndex, depth = position899, tokenIndex899, depth899
							if buffer[position] != rune('I') {
								goto l891
							}
							position++
						}
					l899:
						{
							position901, tokenIndex901, depth901 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l902
							}
							position++
							goto l901
						l902:
							position, tokenIndex, depth = position901, tokenIndex901, depth901
							if buffer[position] != rune('N') {
								goto l891
							}
							position++
						}
					l901:
						{
							position903, tokenIndex903, depth903 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l904
							}
							position++
							goto l903
						l904:
							position, tokenIndex, depth = position903, tokenIndex903, depth903
							if buffer[position] != rune('G') {
								goto l891
							}
							position++
						}
					l903:
						if !_rules[rulesp]() {
							goto l891
						}
						if !_rules[ruleExpression]() {
							goto l891
						}
						goto l892
					l891:
						position, tokenIndex, depth = position891, tokenIndex891, depth891
					}
				l892:
					depth--
					add(rulePegText, position890)
				}
				if !_rules[ruleAction39]() {
					goto l888
				}
				depth--
				add(ruleHaving, position889)
			}
			return true
		l888:
			position, tokenIndex, depth = position888, tokenIndex888, depth888
			return false
		},
		/* 53 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action40))> */
		func() bool {
			position905, tokenIndex905, depth905 := position, tokenIndex, depth
			{
				position906 := position
				depth++
				{
					position907, tokenIndex907, depth907 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l908
					}
					goto l907
				l908:
					position, tokenIndex, depth = position907, tokenIndex907, depth907
					if !_rules[ruleStreamWindow]() {
						goto l905
					}
					if !_rules[ruleAction40]() {
						goto l905
					}
				}
			l907:
				depth--
				add(ruleRelationLike, position906)
			}
			return true
		l905:
			position, tokenIndex, depth = position905, tokenIndex905, depth905
			return false
		},
		/* 54 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action41)> */
		func() bool {
			position909, tokenIndex909, depth909 := position, tokenIndex, depth
			{
				position910 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l909
				}
				if !_rules[rulesp]() {
					goto l909
				}
				{
					position911, tokenIndex911, depth911 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l912
					}
					position++
					goto l911
				l912:
					position, tokenIndex, depth = position911, tokenIndex911, depth911
					if buffer[position] != rune('A') {
						goto l909
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
						goto l909
					}
					position++
				}
			l913:
				if !_rules[rulesp]() {
					goto l909
				}
				if !_rules[ruleIdentifier]() {
					goto l909
				}
				if !_rules[ruleAction41]() {
					goto l909
				}
				depth--
				add(ruleAliasedStreamWindow, position910)
			}
			return true
		l909:
			position, tokenIndex, depth = position909, tokenIndex909, depth909
			return false
		},
		/* 55 StreamWindow <- <(StreamLike spOpt '[' spOpt (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval CapacitySpecOpt SheddingSpecOpt spOpt ']' Action42)> */
		func() bool {
			position915, tokenIndex915, depth915 := position, tokenIndex, depth
			{
				position916 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l915
				}
				if !_rules[rulespOpt]() {
					goto l915
				}
				if buffer[position] != rune('[') {
					goto l915
				}
				position++
				if !_rules[rulespOpt]() {
					goto l915
				}
				{
					position917, tokenIndex917, depth917 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l918
					}
					position++
					goto l917
				l918:
					position, tokenIndex, depth = position917, tokenIndex917, depth917
					if buffer[position] != rune('R') {
						goto l915
					}
					position++
				}
			l917:
				{
					position919, tokenIndex919, depth919 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l920
					}
					position++
					goto l919
				l920:
					position, tokenIndex, depth = position919, tokenIndex919, depth919
					if buffer[position] != rune('A') {
						goto l915
					}
					position++
				}
			l919:
				{
					position921, tokenIndex921, depth921 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l922
					}
					position++
					goto l921
				l922:
					position, tokenIndex, depth = position921, tokenIndex921, depth921
					if buffer[position] != rune('N') {
						goto l915
					}
					position++
				}
			l921:
				{
					position923, tokenIndex923, depth923 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l924
					}
					position++
					goto l923
				l924:
					position, tokenIndex, depth = position923, tokenIndex923, depth923
					if buffer[position] != rune('G') {
						goto l915
					}
					position++
				}
			l923:
				{
					position925, tokenIndex925, depth925 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l926
					}
					position++
					goto l925
				l926:
					position, tokenIndex, depth = position925, tokenIndex925, depth925
					if buffer[position] != rune('E') {
						goto l915
					}
					position++
				}
			l925:
				if !_rules[rulesp]() {
					goto l915
				}
				if !_rules[ruleInterval]() {
					goto l915
				}
				if !_rules[ruleCapacitySpecOpt]() {
					goto l915
				}
				if !_rules[ruleSheddingSpecOpt]() {
					goto l915
				}
				if !_rules[rulespOpt]() {
					goto l915
				}
				if buffer[position] != rune(']') {
					goto l915
				}
				position++
				if !_rules[ruleAction42]() {
					goto l915
				}
				depth--
				add(ruleStreamWindow, position916)
			}
			return true
		l915:
			position, tokenIndex, depth = position915, tokenIndex915, depth915
			return false
		},
		/* 56 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position927, tokenIndex927, depth927 := position, tokenIndex, depth
			{
				position928 := position
				depth++
				{
					position929, tokenIndex929, depth929 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l930
					}
					goto l929
				l930:
					position, tokenIndex, depth = position929, tokenIndex929, depth929
					if !_rules[ruleStream]() {
						goto l927
					}
				}
			l929:
				depth--
				add(ruleStreamLike, position928)
			}
			return true
		l927:
			position, tokenIndex, depth = position927, tokenIndex927, depth927
			return false
		},
		/* 57 UDSFFuncApp <- <(FuncAppWithoutOrderBy Action43)> */
		func() bool {
			position931, tokenIndex931, depth931 := position, tokenIndex, depth
			{
				position932 := position
				depth++
				if !_rules[ruleFuncAppWithoutOrderBy]() {
					goto l931
				}
				if !_rules[ruleAction43]() {
					goto l931
				}
				depth--
				add(ruleUDSFFuncApp, position932)
			}
			return true
		l931:
			position, tokenIndex, depth = position931, tokenIndex931, depth931
			return false
		},
		/* 58 CapacitySpecOpt <- <(<(spOpt ',' spOpt (('b' / 'B') ('u' / 'U') ('f' / 'F') ('f' / 'F') ('e' / 'E') ('r' / 'R')) sp (('s' / 'S') ('i' / 'I') ('z' / 'Z') ('e' / 'E')) sp NonNegativeNumericLiteral)?> Action44)> */
		func() bool {
			position933, tokenIndex933, depth933 := position, tokenIndex, depth
			{
				position934 := position
				depth++
				{
					position935 := position
					depth++
					{
						position936, tokenIndex936, depth936 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l936
						}
						if buffer[position] != rune(',') {
							goto l936
						}
						position++
						if !_rules[rulespOpt]() {
							goto l936
						}
						{
							position938, tokenIndex938, depth938 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l939
							}
							position++
							goto l938
						l939:
							position, tokenIndex, depth = position938, tokenIndex938, depth938
							if buffer[position] != rune('B') {
								goto l936
							}
							position++
						}
					l938:
						{
							position940, tokenIndex940, depth940 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l941
							}
							position++
							goto l940
						l941:
							position, tokenIndex, depth = position940, tokenIndex940, depth940
							if buffer[position] != rune('U') {
								goto l936
							}
							position++
						}
					l940:
						{
							position942, tokenIndex942, depth942 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l943
							}
							position++
							goto l942
						l943:
							position, tokenIndex, depth = position942, tokenIndex942, depth942
							if buffer[position] != rune('F') {
								goto l936
							}
							position++
						}
					l942:
						{
							position944, tokenIndex944, depth944 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l945
							}
							position++
							goto l944
						l945:
							position, tokenIndex, depth = position944, tokenIndex944, depth944
							if buffer[position] != rune('F') {
								goto l936
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
								goto l936
							}
							position++
						}
					l946:
						{
							position948, tokenIndex948, depth948 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l949
							}
							position++
							goto l948
						l949:
							position, tokenIndex, depth = position948, tokenIndex948, depth948
							if buffer[position] != rune('R') {
								goto l936
							}
							position++
						}
					l948:
						if !_rules[rulesp]() {
							goto l936
						}
						{
							position950, tokenIndex950, depth950 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l951
							}
							position++
							goto l950
						l951:
							position, tokenIndex, depth = position950, tokenIndex950, depth950
							if buffer[position] != rune('S') {
								goto l936
							}
							position++
						}
					l950:
						{
							position952, tokenIndex952, depth952 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l953
							}
							position++
							goto l952
						l953:
							position, tokenIndex, depth = position952, tokenIndex952, depth952
							if buffer[position] != rune('I') {
								goto l936
							}
							position++
						}
					l952:
						{
							position954, tokenIndex954, depth954 := position, tokenIndex, depth
							if buffer[position] != rune('z') {
								goto l955
							}
							position++
							goto l954
						l955:
							position, tokenIndex, depth = position954, tokenIndex954, depth954
							if buffer[position] != rune('Z') {
								goto l936
							}
							position++
						}
					l954:
						{
							position956, tokenIndex956, depth956 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l957
							}
							position++
							goto l956
						l957:
							position, tokenIndex, depth = position956, tokenIndex956, depth956
							if buffer[position] != rune('E') {
								goto l936
							}
							position++
						}
					l956:
						if !_rules[rulesp]() {
							goto l936
						}
						if !_rules[ruleNonNegativeNumericLiteral]() {
							goto l936
						}
						goto l937
					l936:
						position, tokenIndex, depth = position936, tokenIndex936, depth936
					}
				l937:
					depth--
					add(rulePegText, position935)
				}
				if !_rules[ruleAction44]() {
					goto l933
				}
				depth--
				add(ruleCapacitySpecOpt, position934)
			}
			return true
		l933:
			position, tokenIndex, depth = position933, tokenIndex933, depth933
			return false
		},
		/* 59 SheddingSpecOpt <- <(<(spOpt ',' spOpt SheddingOption sp (('i' / 'I') ('f' / 'F')) sp (('f' / 'F') ('u' / 'U') ('l' / 'L') ('l' / 'L')))?> Action45)> */
		func() bool {
			position958, tokenIndex958, depth958 := position, tokenIndex, depth
			{
				position959 := position
				depth++
				{
					position960 := position
					depth++
					{
						position961, tokenIndex961, depth961 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l961
						}
						if buffer[position] != rune(',') {
							goto l961
						}
						position++
						if !_rules[rulespOpt]() {
							goto l961
						}
						if !_rules[ruleSheddingOption]() {
							goto l961
						}
						if !_rules[rulesp]() {
							goto l961
						}
						{
							position963, tokenIndex963, depth963 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l964
							}
							position++
							goto l963
						l964:
							position, tokenIndex, depth = position963, tokenIndex963, depth963
							if buffer[position] != rune('I') {
								goto l961
							}
							position++
						}
					l963:
						{
							position965, tokenIndex965, depth965 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l966
							}
							position++
							goto l965
						l966:
							position, tokenIndex, depth = position965, tokenIndex965, depth965
							if buffer[position] != rune('F') {
								goto l961
							}
							position++
						}
					l965:
						if !_rules[rulesp]() {
							goto l961
						}
						{
							position967, tokenIndex967, depth967 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l968
							}
							position++
							goto l967
						l968:
							position, tokenIndex, depth = position967, tokenIndex967, depth967
							if buffer[position] != rune('F') {
								goto l961
							}
							position++
						}
					l967:
						{
							position969, tokenIndex969, depth969 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l970
							}
							position++
							goto l969
						l970:
							position, tokenIndex, depth = position969, tokenIndex969, depth969
							if buffer[position] != rune('U') {
								goto l961
							}
							position++
						}
					l969:
						{
							position971, tokenIndex971, depth971 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l972
							}
							position++
							goto l971
						l972:
							position, tokenIndex, depth = position971, tokenIndex971, depth971
							if buffer[position] != rune('L') {
								goto l961
							}
							position++
						}
					l971:
						{
							position973, tokenIndex973, depth973 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l974
							}
							position++
							goto l973
						l974:
							position, tokenIndex, depth = position973, tokenIndex973, depth973
							if buffer[position] != rune('L') {
								goto l961
							}
							position++
						}
					l973:
						goto l962
					l961:
						position, tokenIndex, depth = position961, tokenIndex961, depth961
					}
				l962:
					depth--
					add(rulePegText, position960)
				}
				if !_rules[ruleAction45]() {
					goto l958
				}
				depth--
				add(ruleSheddingSpecOpt, position959)
			}
			return true
		l958:
			position, tokenIndex, depth = position958, tokenIndex958, depth958
			return false
		},
		/* 60 SheddingOption <- <(Wait / DropOldest / DropNewest)> */
		func() bool {
			position975, tokenIndex975, depth975 := position, tokenIndex, depth
			{
				position976 := position
				depth++
				{
					position977, tokenIndex977, depth977 := position, tokenIndex, depth
					if !_rules[ruleWait]() {
						goto l978
					}
					goto l977
				l978:
					position, tokenIndex, depth = position977, tokenIndex977, depth977
					if !_rules[ruleDropOldest]() {
						goto l979
					}
					goto l977
				l979:
					position, tokenIndex, depth = position977, tokenIndex977, depth977
					if !_rules[ruleDropNewest]() {
						goto l975
					}
				}
			l977:
				depth--
				add(ruleSheddingOption, position976)
			}
			return true
		l975:
			position, tokenIndex, depth = position975, tokenIndex975, depth975
			return false
		},
		/* 61 SourceSinkSpecs <- <(<(sp (('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action46)> */
		func() bool {
			position980, tokenIndex980, depth980 := position, tokenIndex, depth
			{
				position981 := position
				depth++
				{
					position982 := position
					depth++
					{
						position983, tokenIndex983, depth983 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l983
						}
						{
							position985, tokenIndex985, depth985 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l986
							}
							position++
							goto l985
						l986:
							position, tokenIndex, depth = position985, tokenIndex985, depth985
							if buffer[position] != rune('W') {
								goto l983
							}
							position++
						}
					l985:
						{
							position987, tokenIndex987, depth987 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l988
							}
							position++
							goto l987
						l988:
							position, tokenIndex, depth = position987, tokenIndex987, depth987
							if buffer[position] != rune('I') {
								goto l983
							}
							position++
						}
					l987:
						{
							position989, tokenIndex989, depth989 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l990
							}
							position++
							goto l989
						l990:
							position, tokenIndex, depth = position989, tokenIndex989, depth989
							if buffer[position] != rune('T') {
								goto l983
							}
							position++
						}
					l989:
						{
							position991, tokenIndex991, depth991 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l992
							}
							position++
							goto l991
						l992:
							position, tokenIndex, depth = position991, tokenIndex991, depth991
							if buffer[position] != rune('H') {
								goto l983
							}
							position++
						}
					l991:
						if !_rules[rulesp]() {
							goto l983
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l983
						}
					l993:
						{
							position994, tokenIndex994, depth994 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l994
							}
							if buffer[position] != rune(',') {
								goto l994
							}
							position++
							if !_rules[rulespOpt]() {
								goto l994
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l994
							}
							goto l993
						l994:
							position, tokenIndex, depth = position994, tokenIndex994, depth994
						}
						goto l984
					l983:
						position, tokenIndex, depth = position983, tokenIndex983, depth983
					}
				l984:
					depth--
					add(rulePegText, position982)
				}
				if !_rules[ruleAction46]() {
					goto l980
				}
				depth--
				add(ruleSourceSinkSpecs, position981)
			}
			return true
		l980:
			position, tokenIndex, depth = position980, tokenIndex980, depth980
			return false
		},
		/* 62 UpdateSourceSinkSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)> Action47)> */
		func() bool {
			position995, tokenIndex995, depth995 := position, tokenIndex, depth
			{
				position996 := position
				depth++
				{
					position997 := position
					depth++
					if !_rules[rulesp]() {
						goto l995
					}
					{
						position998, tokenIndex998, depth998 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l999
						}
						position++
						goto l998
					l999:
						position, tokenIndex, depth = position998, tokenIndex998, depth998
						if buffer[position] != rune('S') {
							goto l995
						}
						position++
					}
				l998:
					{
						position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1001
						}
						position++
						goto l1000
					l1001:
						position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
						if buffer[position] != rune('E') {
							goto l995
						}
						position++
					}
				l1000:
					{
						position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1003
						}
						position++
						goto l1002
					l1003:
						position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
						if buffer[position] != rune('T') {
							goto l995
						}
						position++
					}
				l1002:
					if !_rules[rulesp]() {
						goto l995
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l995
					}
				l1004:
					{
						position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1005
						}
						if buffer[position] != rune(',') {
							goto l1005
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1005
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l1005
						}
						goto l1004
					l1005:
						position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
					}
					depth--
					add(rulePegText, position997)
				}
				if !_rules[ruleAction47]() {
					goto l995
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position996)
			}
			return true
		l995:
			position, tokenIndex, depth = position995, tokenIndex995, depth995
			return false
		},
		/* 63 SetOptSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action48)> */
		func() bool {
			position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
			{
				position1007 := position
				depth++
				{
					position1008 := position
					depth++
					{
						position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1009
						}
						{
							position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1012
							}
							position++
							goto l1011
						l1012:
							position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
							if buffer[position] != rune('S') {
								goto l1009
							}
							position++
						}
					l1011:
						{
							position1013, tokenIndex1013, depth1013 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1014
							}
							position++
							goto l1013
						l1014:
							position, tokenIndex, depth = position1013, tokenIndex1013, depth1013
							if buffer[position] != rune('E') {
								goto l1009
							}
							position++
						}
					l1013:
						{
							position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1016
							}
							position++
							goto l1015
						l1016:
							position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
							if buffer[position] != rune('T') {
								goto l1009
							}
							position++
						}
					l1015:
						if !_rules[rulesp]() {
							goto l1009
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l1009
						}
					l1017:
						{
							position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1018
							}
							if buffer[position] != rune(',') {
								goto l1018
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1018
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l1018
							}
							goto l1017
						l1018:
							position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
						}
						goto l1010
					l1009:
						position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
					}
				l1010:
					depth--
					add(rulePegText, position1008)
				}
				if !_rules[ruleAction48]() {
					goto l1006
				}
				depth--
				add(ruleSetOptSpecs, position1007)
			}
			return true
		l1006:
			position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
			return false
		},
		/* 64 StateTagOpt <- <(<(sp (('t' / 'T') ('a' / 'A') ('g' / 'G')) sp Identifier)?> Action49)> */
		func() bool {
			position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
			{
				position1020 := position
				depth++
				{
					position1021 := position
					depth++
					{
						position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1022
						}
						{
							position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1025
							}
							position++
							goto l1024
						l1025:
							position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
							if buffer[position] != rune('T') {
								goto l1022
							}
							position++
						}
					l1024:
						{
							position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1027
							}
							position++
							goto l1026
						l1027:
							position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
							if buffer[position] != rune('A') {
								goto l1022
							}
							position++
						}
					l1026:
						{
							position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l1029
							}
							position++
							goto l1028
						l1029:
							position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
							if buffer[position] != rune('G') {
								goto l1022
							}
							position++
						}
					l1028:
						if !_rules[rulesp]() {
							goto l1022
						}
						if !_rules[ruleIdentifier]() {
							goto l1022
						}
						goto l1023
					l1022:
						position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
					}
				l1023:
					depth--
					add(rulePegText, position1021)
				}
				if !_rules[ruleAction49]() {
					goto l1019
				}
				depth--
				add(ruleStateTagOpt, position1020)
			}
			return true
		l1019:
			position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
			return false
		},
		/* 65 SourceSinkParam <- <(SourceSinkParamKey spOpt '=' spOpt SourceSinkParamVal Action50)> */
		func() bool {
			position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
			{
				position1031 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l1030
				}
				if !_rules[rulespOpt]() {
					goto l1030
				}
				if buffer[position] != rune('=') {
					goto l1030
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1030
				}
				if !_rules[ruleSourceSinkParamVal]() {
					goto l1030
				}
				if !_rules[ruleAction50]() {
					goto l1030
				}
				depth--
				add(ruleSourceSinkParam, position1031)
			}
			return true
		l1030:
			position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
			return false
		},
		/* 66 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr / ParamMapExpr)> */
		func() bool {
			position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
			{
				position1033 := position
				depth++
				{
					position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l1035
					}
					goto l1034
				l1035:
					position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
					if !_rules[ruleParamArrayExpr]() {
						goto l1036
					}
					goto l1034
				l1036:
					position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
					if !_rules[ruleParamMapExpr]() {
						goto l1032
					}
				}
			l1034:
				depth--
				add(ruleSourceSinkParamVal, position1033)
			}
			return true
		l1032:
			position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
			return false
		},
		/* 67 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
			{
				position1038 := position
				depth++
				{
					position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l1040
					}
					goto l1039
				l1040:
					position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
					if !_rules[ruleLiteral]() {
						goto l1037
					}
				}
			l1039:
				depth--
				add(ruleParamLiteral, position1038)
			}
			return true
		l1037:
			position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
			return false
		},
		/* 68 ParamArrayExpr <- <(<('[' spOpt (ParamLiteral (',' spOpt ParamLiteral)*)? spOpt ','? spOpt ']')> Action51)> */
		func() bool {
			position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
			{
				position1042 := position
				depth++
				{
					position1043 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1041
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1041
					}
					{
						position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l1044
						}
					l1046:
						{
							position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l1047
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1047
							}
							if !_rules[ruleParamLiteral]() {
								goto l1047
							}
							goto l1046
						l1047:
							position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
						}
						goto l1045
					l1044:
						position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
					}
				l1045:
					if !_rules[rulespOpt]() {
						goto l1041
					}
					{
						position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1048
						}
						position++
						goto l1049
					l1048:
						position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
					}
				l1049:
					if !_rules[rulespOpt]() {
						goto l1041
					}
					if buffer[position] != rune(']') {
						goto l1041
					}
					position++
					depth--
					add(rulePegText, position1043)
				}
				if !_rules[ruleAction51]() {
					goto l1041
				}
				depth--
				add(ruleParamArrayExpr, position1042)
			}
			return true
		l1041:
			position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
			return false
		},
		/* 69 ParamMapExpr <- <(<('{' spOpt (ParamKeyValuePair (spOpt ',' spOpt ParamKeyValuePair)*)? spOpt '}')> Action52)> */
		func() bool {
			position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
			{
				position1051 := position
				depth++
				{
					position1052 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1050
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1050
					}
					{
						position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
						if !_rules[ruleParamKeyValuePair]() {
							goto l1053
						}
					l1055:
						{
							position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1056
							}
							if buffer[position] != rune(',') {
								goto l1056
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1056
							}
							if !_rules[ruleParamKeyValuePair]() {
								goto l1056
							}
							goto l1055
						l1056:
							position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
						}
						goto l1054
					l1053:
						position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
					}
				l1054:
					if !_rules[rulespOpt]() {
						goto l1050
					}
					if buffer[position] != rune('}') {
						goto l1050
					}
					position++
					depth--
					add(rulePegText, position1052)
				}
				if !_rules[ruleAction52]() {
					goto l1050
				}
				depth--
				add(ruleParamMapExpr, position1051)
			}
			return true
		l1050:
			position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
			return false
		},
		/* 70 ParamKeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ParamLiteral)> Action53)> */
		func() bool {
			position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
			{
				position1058 := position
				depth++
				{
					position1059 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1057
					}
					if !_rules[rulespOpt]() {
						goto l1057
					}
					if buffer[position] != rune(':') {
						goto l1057
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1057
					}
					if !_rules[ruleParamLiteral]() {
						goto l1057
					}
					depth--
					add(rulePegText, position1059)
				}
				if !_rules[ruleAction53]() {
					goto l1057
				}
				depth--
				add(ruleParamKeyValuePair, position1058)
			}
			return true
		l1057:
			position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
			return false
		},
		/* 71 PausedOpt <- <(<(sp (Paused / Unpaused))?> Action54)> */
		func() bool {
			position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
			{
				position1061 := position
				depth++
				{
					position1062 := position
					depth++
					{
						position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1063
						}
						{
							position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l1066
							}
							goto l1065
						l1066:
							position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
							if !_rules[ruleUnpaused]() {
								goto l1063
							}
						}
					l1065:
						goto l1064
					l1063:
						position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
					}
				l1064:
					depth--
					add(rulePegText, position1062)
				}
				if !_rules[ruleAction54]() {
					goto l1060
				}
				depth--
				add(rulePausedOpt, position1061)
			}
			return true
		l1060:
			position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
			return false
		},
		/* 72 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
			{
				position1068 := position
				depth++
				{
					position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
					if !_rules[ruleWildcard]() {
						goto l1070
					}
					goto l1069
				l1070:
					position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
					if !_rules[ruleExpression]() {
						goto l1067
					}
				}
			l1069:
				depth--
				add(ruleExpressionOrWildcard, position1068)
			}
			return true
		l1067:
			position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
			return false
		},
		/* 73 Expression <- <orExpr> */
		func() bool {
			position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
			{
				position1072 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l1071
				}
				depth--
				add(ruleExpression, position1072)
			}
			return true
		l1071:
			position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
			return false
		},
		/* 74 orExpr <- <(<(andExpr (sp Or sp andExpr)*)> Action55)> */
		func() bool {
			position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
			{
				position1074 := position
				depth++
				{
					position1075 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l1073
					}
				l1076:
					{
						position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1077
						}
						if !_rules[ruleOr]() {
							goto l1077
						}
						if !_rules[rulesp]() {
							goto l1077
						}
						if !_rules[ruleandExpr]() {
							goto l1077
						}
						goto l1076
					l1077:
						position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
					}
					depth--
					add(rulePegText, position1075)
				}
				if !_rules[ruleAction55]() {
					goto l1073
				}
				depth--
				add(ruleorExpr, position1074)
			}
			return true
		l1073:
			position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
			return false
		},
		/* 75 andExpr <- <(<(notExpr (sp And sp notExpr)*)> Action56)> */
		func() bool {
			position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
			{
				position1079 := position
				depth++
				{
					position1080 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l1078
					}
				l1081:
					{
						position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1082
						}
						if !_rules[ruleAnd]() {
							goto l1082
						}
						if !_rules[rulesp]() {
							goto l1082
						}
						if !_rules[rulenotExpr]() {
							goto l1082
						}
						goto l1081
					l1082:
						position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
					}
					depth--
					add(rulePegText, position1080)
				}
				if !_rules[ruleAction56]() {
					goto l1078
				}
				depth--
				add(ruleandExpr, position1079)
			}
			return true
		l1078:
			position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
			return false
		},
		/* 76 notExpr <- <(<((Not sp)? comparisonExpr)> Action57)> */
		func() bool {
			position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
			{
				position1084 := position
				depth++
				{
					position1085 := position
					depth++
					{
						position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l1086
						}
						if !_rules[rulesp]() {
							goto l1086
						}
						goto l1087
					l1086:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
					}
				l1087:
					if !_rules[rulecomparisonExpr]() {
						goto l1083
					}
					depth--
					add(rulePegText, position1085)
				}
				if !_rules[ruleAction57]() {
					goto l1083
				}
				depth--
				add(rulenotExpr, position1084)
			}
			return true
		l1083:
			position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
			return false
		},
		/* 77 comparisonExpr <- <(<(otherOpExpr (spOpt ComparisonOp spOpt otherOpExpr)?)> Action58)> */
		func() bool {
			position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
			{
				position1089 := position
				depth++
				{
					position1090 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l1088
					}
					{
						position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1091
						}
						if !_rules[ruleComparisonOp]() {
							goto l1091
						}
						if !_rules[rulespOpt]() {
							goto l1091
						}
						if !_rules[ruleotherOpExpr]() {
							goto l1091
						}
						goto l1092
					l1091:
						position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
					}
				l1092:
					depth--
					add(rulePegText, position1090)
				}
				if !_rules[ruleAction58]() {
					goto l1088
				}
				depth--
				add(rulecomparisonExpr, position1089)
			}
			return true
		l1088:
			position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
			return false
		},
		/* 78 otherOpExpr <- <(<(isExpr (spOpt OtherOp spOpt isExpr)*)> Action59)> */
		func() bool {
			position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
			{
				position1094 := position
				depth++
				{
					position1095 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l1093
					}
				l1096:
					{
						position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1097
						}
						if !_rules[ruleOtherOp]() {
							goto l1097
						}
						if !_rules[rulespOpt]() {
							goto l1097
						}
						if !_rules[ruleisExpr]() {
							goto l1097
						}
						goto l1096
					l1097:
						position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
					}
					depth--
					add(rulePegText, position1095)
				}
				if !_rules[ruleAction59]() {
					goto l1093
				}
				depth--
				add(ruleotherOpExpr, position1094)
			}
			return true
		l1093:
			position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
			return false
		},
		/* 79 isExpr <- <(<(termExpr (sp IsOp sp NullLiteral)?)> Action60)> */
		func() bool {
			position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
			{
				position1099 := position
				depth++
				{
					position1100 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l1098
					}
					{
						position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1101
						}
						if !_rules[ruleIsOp]() {
							goto l1101
						}
						if !_rules[rulesp]() {
							goto l1101
						}
						if !_rules[ruleNullLiteral]() {
							goto l1101
						}
						goto l1102
					l1101:
						position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
					}
				l1102:
					depth--
					add(rulePegText, position1100)
				}
				if !_rules[ruleAction60]() {
					goto l1098
				}
				depth--
				add(ruleisExpr, position1099)
			}
			return true
		l1098:
			position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
			return false
		},
		/* 80 termExpr <- <(<(productExpr (spOpt PlusMinusOp spOpt productExpr)*)> Action61)> */
		func() bool {
			position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
			{
				position1104 := position
				depth++
				{
					position1105 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l1103
					}
				l1106:
					{
						position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1107
						}
						if !_rules[rulePlusMinusOp]() {
							goto l1107
						}
						if !_rules[rulespOpt]() {
							goto l1107
						}
						if !_rules[ruleproductExpr]() {
							goto l1107
						}
						goto l1106
					l1107:
						position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
					}
					depth--
					add(rulePegText, position1105)
				}
				if !_rules[ruleAction61]() {
					goto l1103
				}
				depth--
				add(ruletermExpr, position1104)
			}
			return true
		l1103:
			position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
			return false
		},
		/* 81 productExpr <- <(<(minusExpr (spOpt MultDivOp spOpt minusExpr)*)> Action62)> */
		func() bool {
			position1108, tokenIndex1108, depth1108 := position, tokenIndex, depth
			{
				position1109 := position
				depth++
				{
					position1110 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l1108
					}
				l1111:
					{
						position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1112
						}
						if !_rules[ruleMultDivOp]() {
							goto l1112
						}
						if !_rules[rulespOpt]() {
							goto l1112
						}
						if !_rules[ruleminusExpr]() {
							goto l1112
						}
						goto l1111
					l1112:
						position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
					}
					depth--
					add(rulePegText, position1110)
				}
				if !_rules[ruleAction62]() {
					goto l1108
				}
				depth--
				add(ruleproductExpr, position1109)
			}
			return true
		l1108:
			position, tokenIndex, depth = position1108, tokenIndex1108, depth1108
			return false
		},
		/* 82 minusExpr <- <(<((UnaryMinus spOpt)? castExpr)> Action63)> */
		func() bool {
			position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
			{
				position1114 := position
				depth++
				{
					position1115 := position
					depth++
					{
						position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l1116
						}
						if !_rules[rulespOpt]() {
							goto l1116
						}
						goto l1117
					l1116:
						position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
					}
				l1117:
					if !_rules[rulecastExpr]() {
						goto l1113
					}
					depth--
					add(rulePegText, position1115)
				}
				if !_rules[ruleAction63]() {
					goto l1113
				}
				depth--
				add(ruleminusExpr, position1114)
			}
			return true
		l1113:
			position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
			return false
		},
		/* 83 castExpr <- <(<(baseExpr (spOpt (':' ':') spOpt Type)?)> Action64)> */
		func() bool {
			position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
			{
				position1119 := position
				depth++
				{
					position1120 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l1118
					}
					{
						position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1121
						}
						if buffer[position] != rune(':') {
							goto l1121
						}
						position++
						if buffer[position] != rune(':') {
							goto l1121
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1121
						}
						if !_rules[ruleType]() {
							goto l1121
						}
						goto l1122
					l1121:
						position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
					}
				l1122:
					depth--
					add(rulePegText, position1120)
				}
				if !_rules[ruleAction64]() {
					goto l1118
				}
				depth--
				add(rulecastExpr, position1119)
			}
			return true
		l1118:
			position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
			return false
		},
		/* 84 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
			{
				position1124 := position
				depth++
				{
					position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l1126
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1126
					}
					if !_rules[ruleExpression]() {
						goto l1126
					}
					if !_rules[rulespOpt]() {
						goto l1126
					}
					if buffer[position] != rune(')') {
						goto l1126
					}
					position++
					goto l1125
				l1126:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleMapExpr]() {
						goto l1127
					}
					goto l1125
				l1127:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleBooleanLiteral]() {
						goto l1128
					}
					goto l1125
				l1128:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleNullLiteral]() {
						goto l1129
					}
					goto l1125
				l1129:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleRowMeta]() {
						goto l1130
					}
					goto l1125
				l1130:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleFuncTypeCast]() {
						goto l1131
					}
					goto l1125
				l1131:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleFuncApp]() {
						goto l1132
					}
					goto l1125
				l1132:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleRowValue]() {
						goto l1133
					}
					goto l1125
				l1133:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleArrayExpr]() {
						goto l1134
					}
					goto l1125
				l1134:
					position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
					if !_rules[ruleLiteral]() {
						goto l1123
					}
				}
			l1125:
				depth--
				add(rulebaseExpr, position1124)
			}
			return true
		l1123:
			position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
			return false
		},
		/* 85 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') spOpt '(' spOpt Expression sp (('a' / 'A') ('s' / 'S')) sp Type spOpt ')')> Action65)> */
		func() bool {
			position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
			{
				position1136 := position
				depth++
				{
					position1137 := position
					depth++
					{
						position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1139
						}
						position++
						goto l1138
					l1139:
						position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
						if buffer[position] != rune('C') {
							goto l1135
						}
						position++
					}
				l1138:
					{
						position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1141
						}
						position++
						goto l1140
					l1141:
						position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
						if buffer[position] != rune('A') {
							goto l1135
						}
						position++
					}
				l1140:
					{
						position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1143
						}
						position++
						goto l1142
					l1143:
						position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
						if buffer[position] != rune('S') {
							goto l1135
						}
						position++
					}
				l1142:
					{
						position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1145
						}
						position++
						goto l1144
					l1145:
						position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
						if buffer[position] != rune('T') {
							goto l1135
						}
						position++
					}
				l1144:
					if !_rules[rulespOpt]() {
						goto l1135
					}
					if buffer[position] != rune('(') {
						goto l1135
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1135
					}
					if !_rules[ruleExpression]() {
						goto l1135
					}
					if !_rules[rulesp]() {
						goto l1135
					}
					{
						position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1147
						}
						position++
						goto l1146
					l1147:
						position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
						if buffer[position] != rune('A') {
							goto l1135
						}
						position++
					}
				l1146:
					{
						position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1149
						}
						position++
						goto l1148
					l1149:
						position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
						if buffer[position] != rune('S') {
							goto l1135
						}
						position++
					}
				l1148:
					if !_rules[rulesp]() {
						goto l1135
					}
					if !_rules[ruleType]() {
						goto l1135
					}
					if !_rules[rulespOpt]() {
						goto l1135
					}
					if buffer[position] != rune(')') {
						goto l1135
					}
					position++
					depth--
					add(rulePegText, position1137)
				}
				if !_rules[ruleAction65]() {
					goto l1135
				}
				depth--
				add(ruleFuncTypeCast, position1136)
			}
			return true
		l1135:
			position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
			return false
		},
		/* 86 FuncApp <- <(FuncAppWithOrderBy / FuncAppWithoutOrderBy)> */
		func() bool {
			position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
			{
				position1151 := position
				depth++
				{
					position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
					if !_rules[ruleFuncAppWithOrderBy]() {
						goto l1153
					}
					goto l1152
				l1153:
					position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
					if !_rules[ruleFuncAppWithoutOrderBy]() {
						goto l1150
					}
				}
			l1152:
				depth--
				add(ruleFuncApp, position1151)
			}
			return true
		l1150:
			position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
			return false
		},
		/* 87 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action66)> */
		func() bool {
			position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
			{
				position1155 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1154
				}
				if !_rules[rulespOpt]() {
					goto l1154
				}
				if buffer[position] != rune('(') {
					goto l1154
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1154
				}
				if !_rules[ruleFuncParams]() {
					goto l1154
				}
				if !_rules[rulesp]() {
					goto l1154
				}
				if !_rules[ruleParamsOrder]() {
					goto l1154
				}
				if !_rules[rulespOpt]() {
					goto l1154
				}
				if buffer[position] != rune(')') {
					goto l1154
				}
				position++
				if !_rules[ruleAction66]() {
					goto l1154
				}
				depth--
				add(ruleFuncAppWithOrderBy, position1155)
			}
			return true
		l1154:
			position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
			return false
		},
		/* 88 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action67)> */
		func() bool {
			position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
			{
				position1157 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1156
				}
				if !_rules[rulespOpt]() {
					goto l1156
				}
				if buffer[position] != rune('(') {
					goto l1156
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1156
				}
				if !_rules[ruleFuncParams]() {
					goto l1156
				}
				{
					position1158 := position
					depth++
					if !_rules[rulespOpt]() {
						goto l1156
					}
					depth--
					add(rulePegText, position1158)
				}
				if buffer[position] != rune(')') {
					goto l1156
				}
				position++
				if !_rules[ruleAction67]() {
					goto l1156
				}
				depth--
				add(ruleFuncAppWithoutOrderBy, position1157)
			}
			return true
		l1156:
			position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
			return false
		},
		/* 89 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action68)> */
		func() bool {
			position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
			{
				position1160 := position
				depth++
				{
					position1161 := position
					depth++
					{
						position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1162
						}
					l1164:
						{
							position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1165
							}
							if buffer[position] != rune(',') {
								goto l1165
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1165
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1165
							}
							goto l1164
						l1165:
							position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						}
						goto l1163
					l1162:
						position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
					}
				l1163:
					depth--
					add(rulePegText, position1161)
				}
				if !_rules[ruleAction68]() {
					goto l1159
				}
				depth--
				add(ruleFuncParams, position1160)
			}
			return true
		l1159:
			position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
			return false
		},
		/* 90 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action69)> */
		func() bool {
			position1166, tokenIndex1166, depth1166 := position, tokenIndex, depth
			{
				position1167 := position
				depth++
				{
					position1168 := position
					depth++
					{
						position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1170
						}
						position++
						goto l1169
					l1170:
						position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
						if buffer[position] != rune('O') {
							goto l1166
						}
						position++
					}
				l1169:
					{
						position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1172
						}
						position++
						goto l1171
					l1172:
						position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
						if buffer[position] != rune('R') {
							goto l1166
						}
						position++
					}
				l1171:
					{
						position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1174
						}
						position++
						goto l1173
					l1174:
						position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
						if buffer[position] != rune('D') {
							goto l1166
						}
						position++
					}
				l1173:
					{
						position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1176
						}
						position++
						goto l1175
					l1176:
						position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
						if buffer[position] != rune('E') {
							goto l1166
						}
						position++
					}
				l1175:
					{
						position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1178
						}
						position++
						goto l1177
					l1178:
						position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
						if buffer[position] != rune('R') {
							goto l1166
						}
						position++
					}
				l1177:
					if !_rules[rulesp]() {
						goto l1166
					}
					{
						position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1180
						}
						position++
						goto l1179
					l1180:
						position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
						if buffer[position] != rune('B') {
							goto l1166
						}
						position++
					}
				l1179:
					{
						position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1182
						}
						position++
						goto l1181
					l1182:
						position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
						if buffer[position] != rune('Y') {
							goto l1166
						}
						position++
					}
				l1181:
					if !_rules[rulesp]() {
						goto l1166
					}
					if !_rules[ruleSortedExpression]() {
						goto l1166
					}
				l1183:
					{
						position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1184
						}
						if buffer[position] != rune(',') {
							goto l1184
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1184
						}
						if !_rules[ruleSortedExpression]() {
							goto l1184
						}
						goto l1183
					l1184:
						position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
					}
					depth--
					add(rulePegText, position1168)
				}
				if !_rules[ruleAction69]() {
					goto l1166
				}
				depth--
				add(ruleParamsOrder, position1167)
			}
			return true
		l1166:
			position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
			return false
		},
		/* 91 SortedExpression <- <(Expression OrderDirectionOpt Action70)> */
		func() bool {
			position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
			{
				position1186 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l1185
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1185
				}
				if !_rules[ruleAction70]() {
					goto l1185
				}
				depth--
				add(ruleSortedExpression, position1186)
			}
			return true
		l1185:
			position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
			return false
		},
		/* 92 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action71)> */
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
						if !_rules[rulesp]() {
							goto l1190
						}
						{
							position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
							if !_rules[ruleAscending]() {
								goto l1193
							}
							goto l1192
						l1193:
							position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
							if !_rules[ruleDescending]() {
								goto l1190
							}
						}
					l1192:
						goto l1191
					l1190:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
					}
				l1191:
					depth--
					add(rulePegText, position1189)
				}
				if !_rules[ruleAction71]() {
					goto l1187
				}
				depth--
				add(ruleOrderDirectionOpt, position1188)
			}
			return true
		l1187:
			position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
			return false
		},
		/* 93 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action72)> */
		func() bool {
			position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
			{
				position1195 := position
				depth++
				{
					position1196 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1194
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1194
					}
					{
						position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1197
						}
					l1199:
						{
							position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1200
							}
							if buffer[position] != rune(',') {
								goto l1200
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1200
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1200
							}
							goto l1199
						l1200:
							position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
						}
						goto l1198
					l1197:
						position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
					}
				l1198:
					if !_rules[rulespOpt]() {
						goto l1194
					}
					{
						position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1201
						}
						position++
						goto l1202
					l1201:
						position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
					}
				l1202:
					if !_rules[rulespOpt]() {
						goto l1194
					}
					if buffer[position] != rune(']') {
						goto l1194
					}
					position++
					depth--
					add(rulePegText, position1196)
				}
				if !_rules[ruleAction72]() {
					goto l1194
				}
				depth--
				add(ruleArrayExpr, position1195)
			}
			return true
		l1194:
			position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
			return false
		},
		/* 94 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action73)> */
		func() bool {
			position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
			{
				position1204 := position
				depth++
				{
					position1205 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1203
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1203
					}
					{
						position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l1206
						}
					l1208:
						{
							position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1209
							}
							if buffer[position] != rune(',') {
								goto l1209
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1209
							}
							if !_rules[ruleKeyValuePair]() {
								goto l1209
							}
							goto l1208
						l1209:
							position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
						}
						goto l1207
					l1206:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
					}
				l1207:
					if !_rules[rulespOpt]() {
						goto l1203
					}
					if buffer[position] != rune('}') {
						goto l1203
					}
					position++
					depth--
					add(rulePegText, position1205)
				}
				if !_rules[ruleAction73]() {
					goto l1203
				}
				depth--
				add(ruleMapExpr, position1204)
			}
			return true
		l1203:
			position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
			return false
		},
		/* 95 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action74)> */
		func() bool {
			position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
			{
				position1211 := position
				depth++
				{
					position1212 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1210
					}
					if !_rules[rulespOpt]() {
						goto l1210
					}
					if buffer[position] != rune(':') {
						goto l1210
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1210
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1210
					}
					depth--
					add(rulePegText, position1212)
				}
				if !_rules[ruleAction74]() {
					goto l1210
				}
				depth--
				add(ruleKeyValuePair, position1211)
			}
			return true
		l1210:
			position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
			return false
		},
		/* 96 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
			{
				position1214 := position
				depth++
				{
					position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l1216
					}
					goto l1215
				l1216:
					position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
					if !_rules[ruleNumericLiteral]() {
						goto l1217
					}
					goto l1215
				l1217:
					position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
					if !_rules[ruleStringLiteral]() {
						goto l1213
					}
				}
			l1215:
				depth--
				add(ruleLiteral, position1214)
			}
			return true
		l1213:
			position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
			return false
		},
		/* 97 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1218, tokenIndex1218, depth1218 := position, tokenIndex, depth
			{
				position1219 := position
				depth++
				{
					position1220, tokenIndex1220, depth1220 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l1221
					}
					goto l1220
				l1221:
					position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
					if !_rules[ruleNotEqual]() {
						goto l1222
					}
					goto l1220
				l1222:
					position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
					if !_rules[ruleLessOrEqual]() {
						goto l1223
					}
					goto l1220
				l1223:
					position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
					if !_rules[ruleLess]() {
						goto l1224
					}
					goto l1220
				l1224:
					position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
					if !_rules[ruleGreaterOrEqual]() {
						goto l1225
					}
					goto l1220
				l1225:
					position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
					if !_rules[ruleGreater]() {
						goto l1226
					}
					goto l1220
				l1226:
					position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
					if !_rules[ruleNotEqual]() {
						goto l1218
					}
				}
			l1220:
				depth--
				add(ruleComparisonOp, position1219)
			}
			return true
		l1218:
			position, tokenIndex, depth = position1218, tokenIndex1218, depth1218
			return false
		},
		/* 98 OtherOp <- <Concat> */
		func() bool {
			position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
			{
				position1228 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l1227
				}
				depth--
				add(ruleOtherOp, position1228)
			}
			return true
		l1227:
			position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
			return false
		},
		/* 99 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
			{
				position1230 := position
				depth++
				{
					position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l1232
					}
					goto l1231
				l1232:
					position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
					if !_rules[ruleIs]() {
						goto l1229
					}
				}
			l1231:
				depth--
				add(ruleIsOp, position1230)
			}
			return true
		l1229:
			position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
			return false
		},
		/* 100 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
			{
				position1234 := position
				depth++
				{
					position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l1236
					}
					goto l1235
				l1236:
					position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
					if !_rules[ruleMinus]() {
						goto l1233
					}
				}
			l1235:
				depth--
				add(rulePlusMinusOp, position1234)
			}
			return true
		l1233:
			position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
			return false
		},
		/* 101 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
			{
				position1238 := position
				depth++
				{
					position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l1240
					}
					goto l1239
				l1240:
					position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
					if !_rules[ruleDivide]() {
						goto l1241
					}
					goto l1239
				l1241:
					position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
					if !_rules[ruleModulo]() {
						goto l1237
					}
				}
			l1239:
				depth--
				add(ruleMultDivOp, position1238)
			}
			return true
		l1237:
			position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
			return false
		},
		/* 102 Stream <- <(<ident> Action75)> */
		func() bool {
			position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
			{
				position1243 := position
				depth++
				{
					position1244 := position
					depth++
					if !_rules[ruleident]() {
						goto l1242
					}
					depth--
					add(rulePegText, position1244)
				}
				if !_rules[ruleAction75]() {
					goto l1242
				}
				depth--
				add(ruleStream, position1243)
			}
			return true
		l1242:
			position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
			return false
		},
		/* 103 RowMeta <- <RowTimestamp> */
		func() bool {
			position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
			{
				position1246 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l1245
				}
				depth--
				add(ruleRowMeta, position1246)
			}
			return true
		l1245:
			position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
			return false
		},
		/* 104 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action76)> */
		func() bool {
			position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
			{
				position1248 := position
				depth++
				{
					position1249 := position
					depth++
					{
						position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1250
						}
						if buffer[position] != rune(':') {
							goto l1250
						}
						position++
						goto l1251
					l1250:
						position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
					}
				l1251:
					if buffer[position] != rune('t') {
						goto l1247
					}
					position++
					if buffer[position] != rune('s') {
						goto l1247
					}
					position++
					if buffer[position] != rune('(') {
						goto l1247
					}
					position++
					if buffer[position] != rune(')') {
						goto l1247
					}
					position++
					depth--
					add(rulePegText, position1249)
				}
				if !_rules[ruleAction76]() {
					goto l1247
				}
				depth--
				add(ruleRowTimestamp, position1248)
			}
			return true
		l1247:
			position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
			return false
		},
		/* 105 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action77)> */
		func() bool {
			position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
			{
				position1253 := position
				depth++
				{
					position1254 := position
					depth++
					{
						position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1255
						}
						if buffer[position] != rune(':') {
							goto l1255
						}
						position++
						{
							position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1257
							}
							position++
							goto l1255
						l1257:
							position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
						}
						goto l1256
					l1255:
						position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
					}
				l1256:
					if !_rules[rulejsonGetPath]() {
						goto l1252
					}
					depth--
					add(rulePegText, position1254)
				}
				if !_rules[ruleAction77]() {
					goto l1252
				}
				depth--
				add(ruleRowValue, position1253)
			}
			return true
		l1252:
			position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
			return false
		},
		/* 106 NumericLiteral <- <(<('-'? [0-9]+)> Action78)> */
		func() bool {
			position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
			{
				position1259 := position
				depth++
				{
					position1260 := position
					depth++
					{
						position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1261
						}
						position++
						goto l1262
					l1261:
						position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
					}
				l1262:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1258
					}
					position++
				l1263:
					{
						position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1264
						}
						position++
						goto l1263
					l1264:
						position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
					}
					depth--
					add(rulePegText, position1260)
				}
				if !_rules[ruleAction78]() {
					goto l1258
				}
				depth--
				add(ruleNumericLiteral, position1259)
			}
			return true
		l1258:
			position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
			return false
		},
		/* 107 NonNegativeNumericLiteral <- <(<[0-9]+> Action79)> */
		func() bool {
			position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
			{
				position1266 := position
				depth++
				{
					position1267 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1265
					}
					position++
				l1268:
					{
						position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1269
						}
						position++
						goto l1268
					l1269:
						position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
					}
					depth--
					add(rulePegText, position1267)
				}
				if !_rules[ruleAction79]() {
					goto l1265
				}
				depth--
				add(ruleNonNegativeNumericLiteral, position1266)
			}
			return true
		l1265:
			position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
			return false
		},
		/* 108 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action80)> */
		func() bool {
			position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
			{
				position1271 := position
				depth++
				{
					position1272 := position
					depth++
					{
						position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1273
						}
						position++
						goto l1274
					l1273:
						position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
					}
				l1274:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1270
					}
					position++
				l1275:
					{
						position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1276
						}
						position++
						goto l1275
					l1276:
						position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
					}
					if buffer[position] != rune('.') {
						goto l1270
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1270
					}
					position++
				l1277:
					{
						position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1278
						}
						position++
						goto l1277
					l1278:
						position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
					}
					depth--
					add(rulePegText, position1272)
				}
				if !_rules[ruleAction80]() {
					goto l1270
				}
				depth--
				add(ruleFloatLiteral, position1271)
			}
			return true
		l1270:
			position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
			return false
		},
		/* 109 Function <- <(<ident> Action81)> */
		func() bool {
			position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
			{
				position1280 := position
				depth++
				{
					position1281 := position
					depth++
					if !_rules[ruleident]() {
						goto l1279
					}
					depth--
					add(rulePegText, position1281)
				}
				if !_rules[ruleAction81]() {
					goto l1279
				}
				depth--
				add(ruleFunction, position1280)
			}
			return true
		l1279:
			position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
			return false
		},
		/* 110 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action82)> */
		func() bool {
			position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
			{
				position1283 := position
				depth++
				{
					position1284 := position
					depth++
					{
						position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1286
						}
						position++
						goto l1285
					l1286:
						position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
						if buffer[position] != rune('N') {
							goto l1282
						}
						position++
					}
				l1285:
					{
						position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1288
						}
						position++
						goto l1287
					l1288:
						position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
						if buffer[position] != rune('U') {
							goto l1282
						}
						position++
					}
				l1287:
					{
						position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1290
						}
						position++
						goto l1289
					l1290:
						position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
						if buffer[position] != rune('L') {
							goto l1282
						}
						position++
					}
				l1289:
					{
						position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1292
						}
						position++
						goto l1291
					l1292:
						position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
						if buffer[position] != rune('L') {
							goto l1282
						}
						position++
					}
				l1291:
					depth--
					add(rulePegText, position1284)
				}
				if !_rules[ruleAction82]() {
					goto l1282
				}
				depth--
				add(ruleNullLiteral, position1283)
			}
			return true
		l1282:
			position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
			return false
		},
		/* 111 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
			{
				position1294 := position
				depth++
				{
					position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l1296
					}
					goto l1295
				l1296:
					position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
					if !_rules[ruleFALSE]() {
						goto l1293
					}
				}
			l1295:
				depth--
				add(ruleBooleanLiteral, position1294)
			}
			return true
		l1293:
			position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
			return false
		},
		/* 112 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action83)> */
		func() bool {
			position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
			{
				position1298 := position
				depth++
				{
					position1299 := position
					depth++
					{
						position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1301
						}
						position++
						goto l1300
					l1301:
						position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
						if buffer[position] != rune('T') {
							goto l1297
						}
						position++
					}
				l1300:
					{
						position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1303
						}
						position++
						goto l1302
					l1303:
						position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
						if buffer[position] != rune('R') {
							goto l1297
						}
						position++
					}
				l1302:
					{
						position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1305
						}
						position++
						goto l1304
					l1305:
						position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
						if buffer[position] != rune('U') {
							goto l1297
						}
						position++
					}
				l1304:
					{
						position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1307
						}
						position++
						goto l1306
					l1307:
						position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
						if buffer[position] != rune('E') {
							goto l1297
						}
						position++
					}
				l1306:
					depth--
					add(rulePegText, position1299)
				}
				if !_rules[ruleAction83]() {
					goto l1297
				}
				depth--
				add(ruleTRUE, position1298)
			}
			return true
		l1297:
			position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
			return false
		},
		/* 113 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action84)> */
		func() bool {
			position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
			{
				position1309 := position
				depth++
				{
					position1310 := position
					depth++
					{
						position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1312
						}
						position++
						goto l1311
					l1312:
						position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
						if buffer[position] != rune('F') {
							goto l1308
						}
						position++
					}
				l1311:
					{
						position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1314
						}
						position++
						goto l1313
					l1314:
						position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
						if buffer[position] != rune('A') {
							goto l1308
						}
						position++
					}
				l1313:
					{
						position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1316
						}
						position++
						goto l1315
					l1316:
						position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
						if buffer[position] != rune('L') {
							goto l1308
						}
						position++
					}
				l1315:
					{
						position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1318
						}
						position++
						goto l1317
					l1318:
						position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
						if buffer[position] != rune('S') {
							goto l1308
						}
						position++
					}
				l1317:
					{
						position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1320
						}
						position++
						goto l1319
					l1320:
						position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
						if buffer[position] != rune('E') {
							goto l1308
						}
						position++
					}
				l1319:
					depth--
					add(rulePegText, position1310)
				}
				if !_rules[ruleAction84]() {
					goto l1308
				}
				depth--
				add(ruleFALSE, position1309)
			}
			return true
		l1308:
			position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
			return false
		},
		/* 114 Wildcard <- <(<((ident ':' !':')? '*')> Action85)> */
		func() bool {
			position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
			{
				position1322 := position
				depth++
				{
					position1323 := position
					depth++
					{
						position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1324
						}
						if buffer[position] != rune(':') {
							goto l1324
						}
						position++
						{
							position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1326
							}
							position++
							goto l1324
						l1326:
							position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
						}
						goto l1325
					l1324:
						position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
					}
				l1325:
					if buffer[position] != rune('*') {
						goto l1321
					}
					position++
					depth--
					add(rulePegText, position1323)
				}
				if !_rules[ruleAction85]() {
					goto l1321
				}
				depth--
				add(ruleWildcard, position1322)
			}
			return true
		l1321:
			position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
			return false
		},
		/* 115 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action86)> */
		func() bool {
			position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
			{
				position1328 := position
				depth++
				{
					position1329 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l1327
					}
					position++
				l1330:
					{
						position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
						{
							position1332, tokenIndex1332, depth1332 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1333
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1333
							}
							position++
							goto l1332
						l1333:
							position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
							{
								position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1334
								}
								position++
								goto l1331
							l1334:
								position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
							}
							if !matchDot() {
								goto l1331
							}
						}
					l1332:
						goto l1330
					l1331:
						position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
					}
					if buffer[position] != rune('\'') {
						goto l1327
					}
					position++
					depth--
					add(rulePegText, position1329)
				}
				if !_rules[ruleAction86]() {
					goto l1327
				}
				depth--
				add(ruleStringLiteral, position1328)
			}
			return true
		l1327:
			position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
			return false
		},
		/* 116 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action87)> */
		func() bool {
			position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
			{
				position1336 := position
				depth++
				{
					position1337 := position
					depth++
					{
						position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1339
						}
						position++
						goto l1338
					l1339:
						position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
						if buffer[position] != rune('I') {
							goto l1335
						}
						position++
					}
				l1338:
					{
						position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1341
						}
						position++
						goto l1340
					l1341:
						position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
						if buffer[position] != rune('S') {
							goto l1335
						}
						position++
					}
				l1340:
					{
						position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1343
						}
						position++
						goto l1342
					l1343:
						position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
						if buffer[position] != rune('T') {
							goto l1335
						}
						position++
					}
				l1342:
					{
						position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1345
						}
						position++
						goto l1344
					l1345:
						position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
						if buffer[position] != rune('R') {
							goto l1335
						}
						position++
					}
				l1344:
					{
						position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1347
						}
						position++
						goto l1346
					l1347:
						position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
						if buffer[position] != rune('E') {
							goto l1335
						}
						position++
					}
				l1346:
					{
						position1348, tokenIndex1348, depth1348 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1349
						}
						position++
						goto l1348
					l1349:
						position, tokenIndex, depth = position1348, tokenIndex1348, depth1348
						if buffer[position] != rune('A') {
							goto l1335
						}
						position++
					}
				l1348:
					{
						position1350, tokenIndex1350, depth1350 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1351
						}
						position++
						goto l1350
					l1351:
						position, tokenIndex, depth = position1350, tokenIndex1350, depth1350
						if buffer[position] != rune('M') {
							goto l1335
						}
						position++
					}
				l1350:
					depth--
					add(rulePegText, position1337)
				}
				if !_rules[ruleAction87]() {
					goto l1335
				}
				depth--
				add(ruleISTREAM, position1336)
			}
			return true
		l1335:
			position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
			return false
		},
		/* 117 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action88)> */
		func() bool {
			position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
			{
				position1353 := position
				depth++
				{
					position1354 := position
					depth++
					{
						position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1356
						}
						position++
						goto l1355
					l1356:
						position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
						if buffer[position] != rune('D') {
							goto l1352
						}
						position++
					}
				l1355:
					{
						position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1358
						}
						position++
						goto l1357
					l1358:
						position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
						if buffer[position] != rune('S') {
							goto l1352
						}
						position++
					}
				l1357:
					{
						position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1360
						}
						position++
						goto l1359
					l1360:
						position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
						if buffer[position] != rune('T') {
							goto l1352
						}
						position++
					}
				l1359:
					{
						position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1362
						}
						position++
						goto l1361
					l1362:
						position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
						if buffer[position] != rune('R') {
							goto l1352
						}
						position++
					}
				l1361:
					{
						position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1364
						}
						position++
						goto l1363
					l1364:
						position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
						if buffer[position] != rune('E') {
							goto l1352
						}
						position++
					}
				l1363:
					{
						position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1366
						}
						position++
						goto l1365
					l1366:
						position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
						if buffer[position] != rune('A') {
							goto l1352
						}
						position++
					}
				l1365:
					{
						position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1368
						}
						position++
						goto l1367
					l1368:
						position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
						if buffer[position] != rune('M') {
							goto l1352
						}
						position++
					}
				l1367:
					depth--
					add(rulePegText, position1354)
				}
				if !_rules[ruleAction88]() {
					goto l1352
				}
				depth--
				add(ruleDSTREAM, position1353)
			}
			return true
		l1352:
			position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
			return false
		},
		/* 118 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action89)> */
		func() bool {
			position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
			{
				position1370 := position
				depth++
				{
					position1371 := position
					depth++
					{
						position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1373
						}
						position++
						goto l1372
					l1373:
						position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
						if buffer[position] != rune('R') {
							goto l1369
						}
						position++
					}
				l1372:
					{
						position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1375
						}
						position++
						goto l1374
					l1375:
						position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
						if buffer[position] != rune('S') {
							goto l1369
						}
						position++
					}
				l1374:
					{
						position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1377
						}
						position++
						goto l1376
					l1377:
						position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
						if buffer[position] != rune('T') {
							goto l1369
						}
						position++
					}
				l1376:
					{
						position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1379
						}
						position++
						goto l1378
					l1379:
						position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
						if buffer[position] != rune('R') {
							goto l1369
						}
						position++
					}
				l1378:
					{
						position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1381
						}
						position++
						goto l1380
					l1381:
						position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
						if buffer[position] != rune('E') {
							goto l1369
						}
						position++
					}
				l1380:
					{
						position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1383
						}
						position++
						goto l1382
					l1383:
						position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
						if buffer[position] != rune('A') {
							goto l1369
						}
						position++
					}
				l1382:
					{
						position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1385
						}
						position++
						goto l1384
					l1385:
						position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
						if buffer[position] != rune('M') {
							goto l1369
						}
						position++
					}
				l1384:
					depth--
					add(rulePegText, position1371)
				}
				if !_rules[ruleAction89]() {
					goto l1369
				}
				depth--
				add(ruleRSTREAM, position1370)
			}
			return true
		l1369:
			position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
			return false
		},
		/* 119 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action90)> */
		func() bool {
			position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
			{
				position1387 := position
				depth++
				{
					position1388 := position
					depth++
					{
						position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1390
						}
						position++
						goto l1389
					l1390:
						position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
						if buffer[position] != rune('T') {
							goto l1386
						}
						position++
					}
				l1389:
					{
						position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1392
						}
						position++
						goto l1391
					l1392:
						position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
						if buffer[position] != rune('U') {
							goto l1386
						}
						position++
					}
				l1391:
					{
						position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1394
						}
						position++
						goto l1393
					l1394:
						position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
						if buffer[position] != rune('P') {
							goto l1386
						}
						position++
					}
				l1393:
					{
						position1395, tokenIndex1395, depth1395 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1396
						}
						position++
						goto l1395
					l1396:
						position, tokenIndex, depth = position1395, tokenIndex1395, depth1395
						if buffer[position] != rune('L') {
							goto l1386
						}
						position++
					}
				l1395:
					{
						position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1398
						}
						position++
						goto l1397
					l1398:
						position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
						if buffer[position] != rune('E') {
							goto l1386
						}
						position++
					}
				l1397:
					{
						position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1400
						}
						position++
						goto l1399
					l1400:
						position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
						if buffer[position] != rune('S') {
							goto l1386
						}
						position++
					}
				l1399:
					depth--
					add(rulePegText, position1388)
				}
				if !_rules[ruleAction90]() {
					goto l1386
				}
				depth--
				add(ruleTUPLES, position1387)
			}
			return true
		l1386:
			position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
			return false
		},
		/* 120 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action91)> */
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
						if buffer[position] != rune('s') {
							goto l1405
						}
						position++
						goto l1404
					l1405:
						position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
						if buffer[position] != rune('S') {
							goto l1401
						}
						position++
					}
				l1404:
					{
						position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1407
						}
						position++
						goto l1406
					l1407:
						position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
						if buffer[position] != rune('E') {
							goto l1401
						}
						position++
					}
				l1406:
					{
						position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1409
						}
						position++
						goto l1408
					l1409:
						position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
						if buffer[position] != rune('C') {
							goto l1401
						}
						position++
					}
				l1408:
					{
						position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1411
						}
						position++
						goto l1410
					l1411:
						position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
						if buffer[position] != rune('O') {
							goto l1401
						}
						position++
					}
				l1410:
					{
						position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1413
						}
						position++
						goto l1412
					l1413:
						position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
						if buffer[position] != rune('N') {
							goto l1401
						}
						position++
					}
				l1412:
					{
						position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1415
						}
						position++
						goto l1414
					l1415:
						position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
						if buffer[position] != rune('D') {
							goto l1401
						}
						position++
					}
				l1414:
					{
						position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1417
						}
						position++
						goto l1416
					l1417:
						position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
						if buffer[position] != rune('S') {
							goto l1401
						}
						position++
					}
				l1416:
					depth--
					add(rulePegText, position1403)
				}
				if !_rules[ruleAction91]() {
					goto l1401
				}
				depth--
				add(ruleSECONDS, position1402)
			}
			return true
		l1401:
			position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
			return false
		},
		/* 121 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action92)> */
		func() bool {
			position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
			{
				position1419 := position
				depth++
				{
					position1420 := position
					depth++
					{
						position1421, tokenIndex1421, depth1421 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1422
						}
						position++
						goto l1421
					l1422:
						position, tokenIndex, depth = position1421, tokenIndex1421, depth1421
						if buffer[position] != rune('M') {
							goto l1418
						}
						position++
					}
				l1421:
					{
						position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1424
						}
						position++
						goto l1423
					l1424:
						position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
						if buffer[position] != rune('I') {
							goto l1418
						}
						position++
					}
				l1423:
					{
						position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1426
						}
						position++
						goto l1425
					l1426:
						position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
						if buffer[position] != rune('L') {
							goto l1418
						}
						position++
					}
				l1425:
					{
						position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1428
						}
						position++
						goto l1427
					l1428:
						position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
						if buffer[position] != rune('L') {
							goto l1418
						}
						position++
					}
				l1427:
					{
						position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1430
						}
						position++
						goto l1429
					l1430:
						position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
						if buffer[position] != rune('I') {
							goto l1418
						}
						position++
					}
				l1429:
					{
						position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1432
						}
						position++
						goto l1431
					l1432:
						position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
						if buffer[position] != rune('S') {
							goto l1418
						}
						position++
					}
				l1431:
					{
						position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1434
						}
						position++
						goto l1433
					l1434:
						position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
						if buffer[position] != rune('E') {
							goto l1418
						}
						position++
					}
				l1433:
					{
						position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1436
						}
						position++
						goto l1435
					l1436:
						position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
						if buffer[position] != rune('C') {
							goto l1418
						}
						position++
					}
				l1435:
					{
						position1437, tokenIndex1437, depth1437 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1438
						}
						position++
						goto l1437
					l1438:
						position, tokenIndex, depth = position1437, tokenIndex1437, depth1437
						if buffer[position] != rune('O') {
							goto l1418
						}
						position++
					}
				l1437:
					{
						position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1440
						}
						position++
						goto l1439
					l1440:
						position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
						if buffer[position] != rune('N') {
							goto l1418
						}
						position++
					}
				l1439:
					{
						position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1442
						}
						position++
						goto l1441
					l1442:
						position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
						if buffer[position] != rune('D') {
							goto l1418
						}
						position++
					}
				l1441:
					{
						position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1444
						}
						position++
						goto l1443
					l1444:
						position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
						if buffer[position] != rune('S') {
							goto l1418
						}
						position++
					}
				l1443:
					depth--
					add(rulePegText, position1420)
				}
				if !_rules[ruleAction92]() {
					goto l1418
				}
				depth--
				add(ruleMILLISECONDS, position1419)
			}
			return true
		l1418:
			position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
			return false
		},
		/* 122 Wait <- <(<(('w' / 'W') ('a' / 'A') ('i' / 'I') ('t' / 'T'))> Action93)> */
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
						if buffer[position] != rune('w') {
							goto l1449
						}
						position++
						goto l1448
					l1449:
						position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
						if buffer[position] != rune('W') {
							goto l1445
						}
						position++
					}
				l1448:
					{
						position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1451
						}
						position++
						goto l1450
					l1451:
						position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
						if buffer[position] != rune('A') {
							goto l1445
						}
						position++
					}
				l1450:
					{
						position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1453
						}
						position++
						goto l1452
					l1453:
						position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
						if buffer[position] != rune('I') {
							goto l1445
						}
						position++
					}
				l1452:
					{
						position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1455
						}
						position++
						goto l1454
					l1455:
						position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
						if buffer[position] != rune('T') {
							goto l1445
						}
						position++
					}
				l1454:
					depth--
					add(rulePegText, position1447)
				}
				if !_rules[ruleAction93]() {
					goto l1445
				}
				depth--
				add(ruleWait, position1446)
			}
			return true
		l1445:
			position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
			return false
		},
		/* 123 DropOldest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('o' / 'O') ('l' / 'L') ('d' / 'D') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action94)> */
		func() bool {
			position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
			{
				position1457 := position
				depth++
				{
					position1458 := position
					depth++
					{
						position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1460
						}
						position++
						goto l1459
					l1460:
						position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
						if buffer[position] != rune('D') {
							goto l1456
						}
						position++
					}
				l1459:
					{
						position1461, tokenIndex1461, depth1461 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1462
						}
						position++
						goto l1461
					l1462:
						position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
						if buffer[position] != rune('R') {
							goto l1456
						}
						position++
					}
				l1461:
					{
						position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1464
						}
						position++
						goto l1463
					l1464:
						position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
						if buffer[position] != rune('O') {
							goto l1456
						}
						position++
					}
				l1463:
					{
						position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1466
						}
						position++
						goto l1465
					l1466:
						position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
						if buffer[position] != rune('P') {
							goto l1456
						}
						position++
					}
				l1465:
					if !_rules[rulesp]() {
						goto l1456
					}
					{
						position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1468
						}
						position++
						goto l1467
					l1468:
						position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
						if buffer[position] != rune('O') {
							goto l1456
						}
						position++
					}
				l1467:
					{
						position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1470
						}
						position++
						goto l1469
					l1470:
						position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
						if buffer[position] != rune('L') {
							goto l1456
						}
						position++
					}
				l1469:
					{
						position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1472
						}
						position++
						goto l1471
					l1472:
						position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
						if buffer[position] != rune('D') {
							goto l1456
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
							goto l1456
						}
						position++
					}
				l1473:
					{
						position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1476
						}
						position++
						goto l1475
					l1476:
						position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
						if buffer[position] != rune('S') {
							goto l1456
						}
						position++
					}
				l1475:
					{
						position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1478
						}
						position++
						goto l1477
					l1478:
						position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
						if buffer[position] != rune('T') {
							goto l1456
						}
						position++
					}
				l1477:
					depth--
					add(rulePegText, position1458)
				}
				if !_rules[ruleAction94]() {
					goto l1456
				}
				depth--
				add(ruleDropOldest, position1457)
			}
			return true
		l1456:
			position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
			return false
		},
		/* 124 DropNewest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('n' / 'N') ('e' / 'E') ('w' / 'W') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action95)> */
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
						if buffer[position] != rune('d') {
							goto l1483
						}
						position++
						goto l1482
					l1483:
						position, tokenIndex, depth = position1482, tokenIndex1482, depth1482
						if buffer[position] != rune('D') {
							goto l1479
						}
						position++
					}
				l1482:
					{
						position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1485
						}
						position++
						goto l1484
					l1485:
						position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
						if buffer[position] != rune('R') {
							goto l1479
						}
						position++
					}
				l1484:
					{
						position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1487
						}
						position++
						goto l1486
					l1487:
						position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
						if buffer[position] != rune('O') {
							goto l1479
						}
						position++
					}
				l1486:
					{
						position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1489
						}
						position++
						goto l1488
					l1489:
						position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
						if buffer[position] != rune('P') {
							goto l1479
						}
						position++
					}
				l1488:
					if !_rules[rulesp]() {
						goto l1479
					}
					{
						position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1491
						}
						position++
						goto l1490
					l1491:
						position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
						if buffer[position] != rune('N') {
							goto l1479
						}
						position++
					}
				l1490:
					{
						position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1493
						}
						position++
						goto l1492
					l1493:
						position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
						if buffer[position] != rune('E') {
							goto l1479
						}
						position++
					}
				l1492:
					{
						position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
						if buffer[position] != rune('w') {
							goto l1495
						}
						position++
						goto l1494
					l1495:
						position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
						if buffer[position] != rune('W') {
							goto l1479
						}
						position++
					}
				l1494:
					{
						position1496, tokenIndex1496, depth1496 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1497
						}
						position++
						goto l1496
					l1497:
						position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
						if buffer[position] != rune('E') {
							goto l1479
						}
						position++
					}
				l1496:
					{
						position1498, tokenIndex1498, depth1498 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1499
						}
						position++
						goto l1498
					l1499:
						position, tokenIndex, depth = position1498, tokenIndex1498, depth1498
						if buffer[position] != rune('S') {
							goto l1479
						}
						position++
					}
				l1498:
					{
						position1500, tokenIndex1500, depth1500 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1501
						}
						position++
						goto l1500
					l1501:
						position, tokenIndex, depth = position1500, tokenIndex1500, depth1500
						if buffer[position] != rune('T') {
							goto l1479
						}
						position++
					}
				l1500:
					depth--
					add(rulePegText, position1481)
				}
				if !_rules[ruleAction95]() {
					goto l1479
				}
				depth--
				add(ruleDropNewest, position1480)
			}
			return true
		l1479:
			position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
			return false
		},
		/* 125 StreamIdentifier <- <(<ident> Action96)> */
		func() bool {
			position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
			{
				position1503 := position
				depth++
				{
					position1504 := position
					depth++
					if !_rules[ruleident]() {
						goto l1502
					}
					depth--
					add(rulePegText, position1504)
				}
				if !_rules[ruleAction96]() {
					goto l1502
				}
				depth--
				add(ruleStreamIdentifier, position1503)
			}
			return true
		l1502:
			position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
			return false
		},
		/* 126 SourceSinkType <- <(<ident> Action97)> */
		func() bool {
			position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
			{
				position1506 := position
				depth++
				{
					position1507 := position
					depth++
					if !_rules[ruleident]() {
						goto l1505
					}
					depth--
					add(rulePegText, position1507)
				}
				if !_rules[ruleAction97]() {
					goto l1505
				}
				depth--
				add(ruleSourceSinkType, position1506)
			}
			return true
		l1505:
			position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
			return false
		},
		/* 127 SourceSinkParamKey <- <(<ident> Action98)> */
		func() bool {
			position1508, tokenIndex1508, depth1508 := position, tokenIndex, depth
			{
				position1509 := position
				depth++
				{
					position1510 := position
					depth++
					if !_rules[ruleident]() {
						goto l1508
					}
					depth--
					add(rulePegText, position1510)
				}
				if !_rules[ruleAction98]() {
					goto l1508
				}
				depth--
				add(ruleSourceSinkParamKey, position1509)
			}
			return true
		l1508:
			position, tokenIndex, depth = position1508, tokenIndex1508, depth1508
			return false
		},
		/* 128 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action99)> */
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
						if buffer[position] != rune('p') {
							goto l1515
						}
						position++
						goto l1514
					l1515:
						position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
						if buffer[position] != rune('P') {
							goto l1511
						}
						position++
					}
				l1514:
					{
						position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1517
						}
						position++
						goto l1516
					l1517:
						position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
						if buffer[position] != rune('A') {
							goto l1511
						}
						position++
					}
				l1516:
					{
						position1518, tokenIndex1518, depth1518 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1519
						}
						position++
						goto l1518
					l1519:
						position, tokenIndex, depth = position1518, tokenIndex1518, depth1518
						if buffer[position] != rune('U') {
							goto l1511
						}
						position++
					}
				l1518:
					{
						position1520, tokenIndex1520, depth1520 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1521
						}
						position++
						goto l1520
					l1521:
						position, tokenIndex, depth = position1520, tokenIndex1520, depth1520
						if buffer[position] != rune('S') {
							goto l1511
						}
						position++
					}
				l1520:
					{
						position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1523
						}
						position++
						goto l1522
					l1523:
						position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
						if buffer[position] != rune('E') {
							goto l1511
						}
						position++
					}
				l1522:
					{
						position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1525
						}
						position++
						goto l1524
					l1525:
						position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
						if buffer[position] != rune('D') {
							goto l1511
						}
						position++
					}
				l1524:
					depth--
					add(rulePegText, position1513)
				}
				if !_rules[ruleAction99]() {
					goto l1511
				}
				depth--
				add(rulePaused, position1512)
			}
			return true
		l1511:
			position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
			return false
		},
		/* 129 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action100)> */
		func() bool {
			position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
			{
				position1527 := position
				depth++
				{
					position1528 := position
					depth++
					{
						position1529, tokenIndex1529, depth1529 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1530
						}
						position++
						goto l1529
					l1530:
						position, tokenIndex, depth = position1529, tokenIndex1529, depth1529
						if buffer[position] != rune('U') {
							goto l1526
						}
						position++
					}
				l1529:
					{
						position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1532
						}
						position++
						goto l1531
					l1532:
						position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
						if buffer[position] != rune('N') {
							goto l1526
						}
						position++
					}
				l1531:
					{
						position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1534
						}
						position++
						goto l1533
					l1534:
						position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
						if buffer[position] != rune('P') {
							goto l1526
						}
						position++
					}
				l1533:
					{
						position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1536
						}
						position++
						goto l1535
					l1536:
						position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
						if buffer[position] != rune('A') {
							goto l1526
						}
						position++
					}
				l1535:
					{
						position1537, tokenIndex1537, depth1537 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1538
						}
						position++
						goto l1537
					l1538:
						position, tokenIndex, depth = position1537, tokenIndex1537, depth1537
						if buffer[position] != rune('U') {
							goto l1526
						}
						position++
					}
				l1537:
					{
						position1539, tokenIndex1539, depth1539 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1540
						}
						position++
						goto l1539
					l1540:
						position, tokenIndex, depth = position1539, tokenIndex1539, depth1539
						if buffer[position] != rune('S') {
							goto l1526
						}
						position++
					}
				l1539:
					{
						position1541, tokenIndex1541, depth1541 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1542
						}
						position++
						goto l1541
					l1542:
						position, tokenIndex, depth = position1541, tokenIndex1541, depth1541
						if buffer[position] != rune('E') {
							goto l1526
						}
						position++
					}
				l1541:
					{
						position1543, tokenIndex1543, depth1543 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1544
						}
						position++
						goto l1543
					l1544:
						position, tokenIndex, depth = position1543, tokenIndex1543, depth1543
						if buffer[position] != rune('D') {
							goto l1526
						}
						position++
					}
				l1543:
					depth--
					add(rulePegText, position1528)
				}
				if !_rules[ruleAction100]() {
					goto l1526
				}
				depth--
				add(ruleUnpaused, position1527)
			}
			return true
		l1526:
			position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
			return false
		},
		/* 130 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action101)> */
		func() bool {
			position1545, tokenIndex1545, depth1545 := position, tokenIndex, depth
			{
				position1546 := position
				depth++
				{
					position1547 := position
					depth++
					{
						position1548, tokenIndex1548, depth1548 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1549
						}
						position++
						goto l1548
					l1549:
						position, tokenIndex, depth = position1548, tokenIndex1548, depth1548
						if buffer[position] != rune('A') {
							goto l1545
						}
						position++
					}
				l1548:
					{
						position1550, tokenIndex1550, depth1550 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1551
						}
						position++
						goto l1550
					l1551:
						position, tokenIndex, depth = position1550, tokenIndex1550, depth1550
						if buffer[position] != rune('S') {
							goto l1545
						}
						position++
					}
				l1550:
					{
						position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1553
						}
						position++
						goto l1552
					l1553:
						position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
						if buffer[position] != rune('C') {
							goto l1545
						}
						position++
					}
				l1552:
					depth--
					add(rulePegText, position1547)
				}
				if !_rules[ruleAction101]() {
					goto l1545
				}
				depth--
				add(ruleAscending, position1546)
			}
			return true
		l1545:
			position, tokenIndex, depth = position1545, tokenIndex1545, depth1545
			return false
		},
		/* 131 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action102)> */
		func() bool {
			position1554, tokenIndex1554, depth1554 := position, tokenIndex, depth
			{
				position1555 := position
				depth++
				{
					position1556 := position
					depth++
					{
						position1557, tokenIndex1557, depth1557 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1558
						}
						position++
						goto l1557
					l1558:
						position, tokenIndex, depth = position1557, tokenIndex1557, depth1557
						if buffer[position] != rune('D') {
							goto l1554
						}
						position++
					}
				l1557:
					{
						position1559, tokenIndex1559, depth1559 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1560
						}
						position++
						goto l1559
					l1560:
						position, tokenIndex, depth = position1559, tokenIndex1559, depth1559
						if buffer[position] != rune('E') {
							goto l1554
						}
						position++
					}
				l1559:
					{
						position1561, tokenIndex1561, depth1561 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1562
						}
						position++
						goto l1561
					l1562:
						position, tokenIndex, depth = position1561, tokenIndex1561, depth1561
						if buffer[position] != rune('S') {
							goto l1554
						}
						position++
					}
				l1561:
					{
						position1563, tokenIndex1563, depth1563 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1564
						}
						position++
						goto l1563
					l1564:
						position, tokenIndex, depth = position1563, tokenIndex1563, depth1563
						if buffer[position] != rune('C') {
							goto l1554
						}
						position++
					}
				l1563:
					depth--
					add(rulePegText, position1556)
				}
				if !_rules[ruleAction102]() {
					goto l1554
				}
				depth--
				add(ruleDescending, position1555)
			}
			return true
		l1554:
			position, tokenIndex, depth = position1554, tokenIndex1554, depth1554
			return false
		},
		/* 132 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
			{
				position1566 := position
				depth++
				{
					position1567, tokenIndex1567, depth1567 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1568
					}
					goto l1567
				l1568:
					position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
					if !_rules[ruleInt]() {
						goto l1569
					}
					goto l1567
				l1569:
					position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
					if !_rules[ruleFloat]() {
						goto l1570
					}
					goto l1567
				l1570:
					position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
					if !_rules[ruleString]() {
						goto l1571
					}
					goto l1567
				l1571:
					position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
					if !_rules[ruleBlob]() {
						goto l1572
					}
					goto l1567
				l1572:
					position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
					if !_rules[ruleTimestamp]() {
						goto l1573
					}
					goto l1567
				l1573:
					position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
					if !_rules[ruleArray]() {
						goto l1574
					}
					goto l1567
				l1574:
					position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
					if !_rules[ruleMap]() {
						goto l1565
					}
				}
			l1567:
				depth--
				add(ruleType, position1566)
			}
			return true
		l1565:
			position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
			return false
		},
		/* 133 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action103)> */
		func() bool {
			position1575, tokenIndex1575, depth1575 := position, tokenIndex, depth
			{
				position1576 := position
				depth++
				{
					position1577 := position
					depth++
					{
						position1578, tokenIndex1578, depth1578 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1579
						}
						position++
						goto l1578
					l1579:
						position, tokenIndex, depth = position1578, tokenIndex1578, depth1578
						if buffer[position] != rune('B') {
							goto l1575
						}
						position++
					}
				l1578:
					{
						position1580, tokenIndex1580, depth1580 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1581
						}
						position++
						goto l1580
					l1581:
						position, tokenIndex, depth = position1580, tokenIndex1580, depth1580
						if buffer[position] != rune('O') {
							goto l1575
						}
						position++
					}
				l1580:
					{
						position1582, tokenIndex1582, depth1582 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1583
						}
						position++
						goto l1582
					l1583:
						position, tokenIndex, depth = position1582, tokenIndex1582, depth1582
						if buffer[position] != rune('O') {
							goto l1575
						}
						position++
					}
				l1582:
					{
						position1584, tokenIndex1584, depth1584 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1585
						}
						position++
						goto l1584
					l1585:
						position, tokenIndex, depth = position1584, tokenIndex1584, depth1584
						if buffer[position] != rune('L') {
							goto l1575
						}
						position++
					}
				l1584:
					depth--
					add(rulePegText, position1577)
				}
				if !_rules[ruleAction103]() {
					goto l1575
				}
				depth--
				add(ruleBool, position1576)
			}
			return true
		l1575:
			position, tokenIndex, depth = position1575, tokenIndex1575, depth1575
			return false
		},
		/* 134 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action104)> */
		func() bool {
			position1586, tokenIndex1586, depth1586 := position, tokenIndex, depth
			{
				position1587 := position
				depth++
				{
					position1588 := position
					depth++
					{
						position1589, tokenIndex1589, depth1589 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1590
						}
						position++
						goto l1589
					l1590:
						position, tokenIndex, depth = position1589, tokenIndex1589, depth1589
						if buffer[position] != rune('I') {
							goto l1586
						}
						position++
					}
				l1589:
					{
						position1591, tokenIndex1591, depth1591 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1592
						}
						position++
						goto l1591
					l1592:
						position, tokenIndex, depth = position1591, tokenIndex1591, depth1591
						if buffer[position] != rune('N') {
							goto l1586
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
							goto l1586
						}
						position++
					}
				l1593:
					depth--
					add(rulePegText, position1588)
				}
				if !_rules[ruleAction104]() {
					goto l1586
				}
				depth--
				add(ruleInt, position1587)
			}
			return true
		l1586:
			position, tokenIndex, depth = position1586, tokenIndex1586, depth1586
			return false
		},
		/* 135 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action105)> */
		func() bool {
			position1595, tokenIndex1595, depth1595 := position, tokenIndex, depth
			{
				position1596 := position
				depth++
				{
					position1597 := position
					depth++
					{
						position1598, tokenIndex1598, depth1598 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1599
						}
						position++
						goto l1598
					l1599:
						position, tokenIndex, depth = position1598, tokenIndex1598, depth1598
						if buffer[position] != rune('F') {
							goto l1595
						}
						position++
					}
				l1598:
					{
						position1600, tokenIndex1600, depth1600 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1601
						}
						position++
						goto l1600
					l1601:
						position, tokenIndex, depth = position1600, tokenIndex1600, depth1600
						if buffer[position] != rune('L') {
							goto l1595
						}
						position++
					}
				l1600:
					{
						position1602, tokenIndex1602, depth1602 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1603
						}
						position++
						goto l1602
					l1603:
						position, tokenIndex, depth = position1602, tokenIndex1602, depth1602
						if buffer[position] != rune('O') {
							goto l1595
						}
						position++
					}
				l1602:
					{
						position1604, tokenIndex1604, depth1604 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1605
						}
						position++
						goto l1604
					l1605:
						position, tokenIndex, depth = position1604, tokenIndex1604, depth1604
						if buffer[position] != rune('A') {
							goto l1595
						}
						position++
					}
				l1604:
					{
						position1606, tokenIndex1606, depth1606 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1607
						}
						position++
						goto l1606
					l1607:
						position, tokenIndex, depth = position1606, tokenIndex1606, depth1606
						if buffer[position] != rune('T') {
							goto l1595
						}
						position++
					}
				l1606:
					depth--
					add(rulePegText, position1597)
				}
				if !_rules[ruleAction105]() {
					goto l1595
				}
				depth--
				add(ruleFloat, position1596)
			}
			return true
		l1595:
			position, tokenIndex, depth = position1595, tokenIndex1595, depth1595
			return false
		},
		/* 136 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action106)> */
		func() bool {
			position1608, tokenIndex1608, depth1608 := position, tokenIndex, depth
			{
				position1609 := position
				depth++
				{
					position1610 := position
					depth++
					{
						position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1612
						}
						position++
						goto l1611
					l1612:
						position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
						if buffer[position] != rune('S') {
							goto l1608
						}
						position++
					}
				l1611:
					{
						position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1614
						}
						position++
						goto l1613
					l1614:
						position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
						if buffer[position] != rune('T') {
							goto l1608
						}
						position++
					}
				l1613:
					{
						position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1616
						}
						position++
						goto l1615
					l1616:
						position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
						if buffer[position] != rune('R') {
							goto l1608
						}
						position++
					}
				l1615:
					{
						position1617, tokenIndex1617, depth1617 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1618
						}
						position++
						goto l1617
					l1618:
						position, tokenIndex, depth = position1617, tokenIndex1617, depth1617
						if buffer[position] != rune('I') {
							goto l1608
						}
						position++
					}
				l1617:
					{
						position1619, tokenIndex1619, depth1619 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1620
						}
						position++
						goto l1619
					l1620:
						position, tokenIndex, depth = position1619, tokenIndex1619, depth1619
						if buffer[position] != rune('N') {
							goto l1608
						}
						position++
					}
				l1619:
					{
						position1621, tokenIndex1621, depth1621 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1622
						}
						position++
						goto l1621
					l1622:
						position, tokenIndex, depth = position1621, tokenIndex1621, depth1621
						if buffer[position] != rune('G') {
							goto l1608
						}
						position++
					}
				l1621:
					depth--
					add(rulePegText, position1610)
				}
				if !_rules[ruleAction106]() {
					goto l1608
				}
				depth--
				add(ruleString, position1609)
			}
			return true
		l1608:
			position, tokenIndex, depth = position1608, tokenIndex1608, depth1608
			return false
		},
		/* 137 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action107)> */
		func() bool {
			position1623, tokenIndex1623, depth1623 := position, tokenIndex, depth
			{
				position1624 := position
				depth++
				{
					position1625 := position
					depth++
					{
						position1626, tokenIndex1626, depth1626 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1627
						}
						position++
						goto l1626
					l1627:
						position, tokenIndex, depth = position1626, tokenIndex1626, depth1626
						if buffer[position] != rune('B') {
							goto l1623
						}
						position++
					}
				l1626:
					{
						position1628, tokenIndex1628, depth1628 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1629
						}
						position++
						goto l1628
					l1629:
						position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
						if buffer[position] != rune('L') {
							goto l1623
						}
						position++
					}
				l1628:
					{
						position1630, tokenIndex1630, depth1630 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1631
						}
						position++
						goto l1630
					l1631:
						position, tokenIndex, depth = position1630, tokenIndex1630, depth1630
						if buffer[position] != rune('O') {
							goto l1623
						}
						position++
					}
				l1630:
					{
						position1632, tokenIndex1632, depth1632 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1633
						}
						position++
						goto l1632
					l1633:
						position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
						if buffer[position] != rune('B') {
							goto l1623
						}
						position++
					}
				l1632:
					depth--
					add(rulePegText, position1625)
				}
				if !_rules[ruleAction107]() {
					goto l1623
				}
				depth--
				add(ruleBlob, position1624)
			}
			return true
		l1623:
			position, tokenIndex, depth = position1623, tokenIndex1623, depth1623
			return false
		},
		/* 138 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action108)> */
		func() bool {
			position1634, tokenIndex1634, depth1634 := position, tokenIndex, depth
			{
				position1635 := position
				depth++
				{
					position1636 := position
					depth++
					{
						position1637, tokenIndex1637, depth1637 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1638
						}
						position++
						goto l1637
					l1638:
						position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
						if buffer[position] != rune('T') {
							goto l1634
						}
						position++
					}
				l1637:
					{
						position1639, tokenIndex1639, depth1639 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1640
						}
						position++
						goto l1639
					l1640:
						position, tokenIndex, depth = position1639, tokenIndex1639, depth1639
						if buffer[position] != rune('I') {
							goto l1634
						}
						position++
					}
				l1639:
					{
						position1641, tokenIndex1641, depth1641 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1642
						}
						position++
						goto l1641
					l1642:
						position, tokenIndex, depth = position1641, tokenIndex1641, depth1641
						if buffer[position] != rune('M') {
							goto l1634
						}
						position++
					}
				l1641:
					{
						position1643, tokenIndex1643, depth1643 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1644
						}
						position++
						goto l1643
					l1644:
						position, tokenIndex, depth = position1643, tokenIndex1643, depth1643
						if buffer[position] != rune('E') {
							goto l1634
						}
						position++
					}
				l1643:
					{
						position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1646
						}
						position++
						goto l1645
					l1646:
						position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
						if buffer[position] != rune('S') {
							goto l1634
						}
						position++
					}
				l1645:
					{
						position1647, tokenIndex1647, depth1647 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1648
						}
						position++
						goto l1647
					l1648:
						position, tokenIndex, depth = position1647, tokenIndex1647, depth1647
						if buffer[position] != rune('T') {
							goto l1634
						}
						position++
					}
				l1647:
					{
						position1649, tokenIndex1649, depth1649 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1650
						}
						position++
						goto l1649
					l1650:
						position, tokenIndex, depth = position1649, tokenIndex1649, depth1649
						if buffer[position] != rune('A') {
							goto l1634
						}
						position++
					}
				l1649:
					{
						position1651, tokenIndex1651, depth1651 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1652
						}
						position++
						goto l1651
					l1652:
						position, tokenIndex, depth = position1651, tokenIndex1651, depth1651
						if buffer[position] != rune('M') {
							goto l1634
						}
						position++
					}
				l1651:
					{
						position1653, tokenIndex1653, depth1653 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1654
						}
						position++
						goto l1653
					l1654:
						position, tokenIndex, depth = position1653, tokenIndex1653, depth1653
						if buffer[position] != rune('P') {
							goto l1634
						}
						position++
					}
				l1653:
					depth--
					add(rulePegText, position1636)
				}
				if !_rules[ruleAction108]() {
					goto l1634
				}
				depth--
				add(ruleTimestamp, position1635)
			}
			return true
		l1634:
			position, tokenIndex, depth = position1634, tokenIndex1634, depth1634
			return false
		},
		/* 139 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action109)> */
		func() bool {
			position1655, tokenIndex1655, depth1655 := position, tokenIndex, depth
			{
				position1656 := position
				depth++
				{
					position1657 := position
					depth++
					{
						position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1659
						}
						position++
						goto l1658
					l1659:
						position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
						if buffer[position] != rune('A') {
							goto l1655
						}
						position++
					}
				l1658:
					{
						position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1661
						}
						position++
						goto l1660
					l1661:
						position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
						if buffer[position] != rune('R') {
							goto l1655
						}
						position++
					}
				l1660:
					{
						position1662, tokenIndex1662, depth1662 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1663
						}
						position++
						goto l1662
					l1663:
						position, tokenIndex, depth = position1662, tokenIndex1662, depth1662
						if buffer[position] != rune('R') {
							goto l1655
						}
						position++
					}
				l1662:
					{
						position1664, tokenIndex1664, depth1664 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1665
						}
						position++
						goto l1664
					l1665:
						position, tokenIndex, depth = position1664, tokenIndex1664, depth1664
						if buffer[position] != rune('A') {
							goto l1655
						}
						position++
					}
				l1664:
					{
						position1666, tokenIndex1666, depth1666 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1667
						}
						position++
						goto l1666
					l1667:
						position, tokenIndex, depth = position1666, tokenIndex1666, depth1666
						if buffer[position] != rune('Y') {
							goto l1655
						}
						position++
					}
				l1666:
					depth--
					add(rulePegText, position1657)
				}
				if !_rules[ruleAction109]() {
					goto l1655
				}
				depth--
				add(ruleArray, position1656)
			}
			return true
		l1655:
			position, tokenIndex, depth = position1655, tokenIndex1655, depth1655
			return false
		},
		/* 140 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action110)> */
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
						if buffer[position] != rune('m') {
							goto l1672
						}
						position++
						goto l1671
					l1672:
						position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
						if buffer[position] != rune('M') {
							goto l1668
						}
						position++
					}
				l1671:
					{
						position1673, tokenIndex1673, depth1673 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1674
						}
						position++
						goto l1673
					l1674:
						position, tokenIndex, depth = position1673, tokenIndex1673, depth1673
						if buffer[position] != rune('A') {
							goto l1668
						}
						position++
					}
				l1673:
					{
						position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1676
						}
						position++
						goto l1675
					l1676:
						position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
						if buffer[position] != rune('P') {
							goto l1668
						}
						position++
					}
				l1675:
					depth--
					add(rulePegText, position1670)
				}
				if !_rules[ruleAction110]() {
					goto l1668
				}
				depth--
				add(ruleMap, position1669)
			}
			return true
		l1668:
			position, tokenIndex, depth = position1668, tokenIndex1668, depth1668
			return false
		},
		/* 141 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action111)> */
		func() bool {
			position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
			{
				position1678 := position
				depth++
				{
					position1679 := position
					depth++
					{
						position1680, tokenIndex1680, depth1680 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1681
						}
						position++
						goto l1680
					l1681:
						position, tokenIndex, depth = position1680, tokenIndex1680, depth1680
						if buffer[position] != rune('O') {
							goto l1677
						}
						position++
					}
				l1680:
					{
						position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1683
						}
						position++
						goto l1682
					l1683:
						position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
						if buffer[position] != rune('R') {
							goto l1677
						}
						position++
					}
				l1682:
					depth--
					add(rulePegText, position1679)
				}
				if !_rules[ruleAction111]() {
					goto l1677
				}
				depth--
				add(ruleOr, position1678)
			}
			return true
		l1677:
			position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
			return false
		},
		/* 142 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action112)> */
		func() bool {
			position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
			{
				position1685 := position
				depth++
				{
					position1686 := position
					depth++
					{
						position1687, tokenIndex1687, depth1687 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1688
						}
						position++
						goto l1687
					l1688:
						position, tokenIndex, depth = position1687, tokenIndex1687, depth1687
						if buffer[position] != rune('A') {
							goto l1684
						}
						position++
					}
				l1687:
					{
						position1689, tokenIndex1689, depth1689 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1690
						}
						position++
						goto l1689
					l1690:
						position, tokenIndex, depth = position1689, tokenIndex1689, depth1689
						if buffer[position] != rune('N') {
							goto l1684
						}
						position++
					}
				l1689:
					{
						position1691, tokenIndex1691, depth1691 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1692
						}
						position++
						goto l1691
					l1692:
						position, tokenIndex, depth = position1691, tokenIndex1691, depth1691
						if buffer[position] != rune('D') {
							goto l1684
						}
						position++
					}
				l1691:
					depth--
					add(rulePegText, position1686)
				}
				if !_rules[ruleAction112]() {
					goto l1684
				}
				depth--
				add(ruleAnd, position1685)
			}
			return true
		l1684:
			position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
			return false
		},
		/* 143 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action113)> */
		func() bool {
			position1693, tokenIndex1693, depth1693 := position, tokenIndex, depth
			{
				position1694 := position
				depth++
				{
					position1695 := position
					depth++
					{
						position1696, tokenIndex1696, depth1696 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1697
						}
						position++
						goto l1696
					l1697:
						position, tokenIndex, depth = position1696, tokenIndex1696, depth1696
						if buffer[position] != rune('N') {
							goto l1693
						}
						position++
					}
				l1696:
					{
						position1698, tokenIndex1698, depth1698 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1699
						}
						position++
						goto l1698
					l1699:
						position, tokenIndex, depth = position1698, tokenIndex1698, depth1698
						if buffer[position] != rune('O') {
							goto l1693
						}
						position++
					}
				l1698:
					{
						position1700, tokenIndex1700, depth1700 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1701
						}
						position++
						goto l1700
					l1701:
						position, tokenIndex, depth = position1700, tokenIndex1700, depth1700
						if buffer[position] != rune('T') {
							goto l1693
						}
						position++
					}
				l1700:
					depth--
					add(rulePegText, position1695)
				}
				if !_rules[ruleAction113]() {
					goto l1693
				}
				depth--
				add(ruleNot, position1694)
			}
			return true
		l1693:
			position, tokenIndex, depth = position1693, tokenIndex1693, depth1693
			return false
		},
		/* 144 Equal <- <(<'='> Action114)> */
		func() bool {
			position1702, tokenIndex1702, depth1702 := position, tokenIndex, depth
			{
				position1703 := position
				depth++
				{
					position1704 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1702
					}
					position++
					depth--
					add(rulePegText, position1704)
				}
				if !_rules[ruleAction114]() {
					goto l1702
				}
				depth--
				add(ruleEqual, position1703)
			}
			return true
		l1702:
			position, tokenIndex, depth = position1702, tokenIndex1702, depth1702
			return false
		},
		/* 145 Less <- <(<'<'> Action115)> */
		func() bool {
			position1705, tokenIndex1705, depth1705 := position, tokenIndex, depth
			{
				position1706 := position
				depth++
				{
					position1707 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1705
					}
					position++
					depth--
					add(rulePegText, position1707)
				}
				if !_rules[ruleAction115]() {
					goto l1705
				}
				depth--
				add(ruleLess, position1706)
			}
			return true
		l1705:
			position, tokenIndex, depth = position1705, tokenIndex1705, depth1705
			return false
		},
		/* 146 LessOrEqual <- <(<('<' '=')> Action116)> */
		func() bool {
			position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
			{
				position1709 := position
				depth++
				{
					position1710 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1708
					}
					position++
					if buffer[position] != rune('=') {
						goto l1708
					}
					position++
					depth--
					add(rulePegText, position1710)
				}
				if !_rules[ruleAction116]() {
					goto l1708
				}
				depth--
				add(ruleLessOrEqual, position1709)
			}
			return true
		l1708:
			position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
			return false
		},
		/* 147 Greater <- <(<'>'> Action117)> */
		func() bool {
			position1711, tokenIndex1711, depth1711 := position, tokenIndex, depth
			{
				position1712 := position
				depth++
				{
					position1713 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1711
					}
					position++
					depth--
					add(rulePegText, position1713)
				}
				if !_rules[ruleAction117]() {
					goto l1711
				}
				depth--
				add(ruleGreater, position1712)
			}
			return true
		l1711:
			position, tokenIndex, depth = position1711, tokenIndex1711, depth1711
			return false
		},
		/* 148 GreaterOrEqual <- <(<('>' '=')> Action118)> */
		func() bool {
			position1714, tokenIndex1714, depth1714 := position, tokenIndex, depth
			{
				position1715 := position
				depth++
				{
					position1716 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1714
					}
					position++
					if buffer[position] != rune('=') {
						goto l1714
					}
					position++
					depth--
					add(rulePegText, position1716)
				}
				if !_rules[ruleAction118]() {
					goto l1714
				}
				depth--
				add(ruleGreaterOrEqual, position1715)
			}
			return true
		l1714:
			position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
			return false
		},
		/* 149 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action119)> */
		func() bool {
			position1717, tokenIndex1717, depth1717 := position, tokenIndex, depth
			{
				position1718 := position
				depth++
				{
					position1719 := position
					depth++
					{
						position1720, tokenIndex1720, depth1720 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1721
						}
						position++
						if buffer[position] != rune('=') {
							goto l1721
						}
						position++
						goto l1720
					l1721:
						position, tokenIndex, depth = position1720, tokenIndex1720, depth1720
						if buffer[position] != rune('<') {
							goto l1717
						}
						position++
						if buffer[position] != rune('>') {
							goto l1717
						}
						position++
					}
				l1720:
					depth--
					add(rulePegText, position1719)
				}
				if !_rules[ruleAction119]() {
					goto l1717
				}
				depth--
				add(ruleNotEqual, position1718)
			}
			return true
		l1717:
			position, tokenIndex, depth = position1717, tokenIndex1717, depth1717
			return false
		},
		/* 150 Concat <- <(<('|' '|')> Action120)> */
		func() bool {
			position1722, tokenIndex1722, depth1722 := position, tokenIndex, depth
			{
				position1723 := position
				depth++
				{
					position1724 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1722
					}
					position++
					if buffer[position] != rune('|') {
						goto l1722
					}
					position++
					depth--
					add(rulePegText, position1724)
				}
				if !_rules[ruleAction120]() {
					goto l1722
				}
				depth--
				add(ruleConcat, position1723)
			}
			return true
		l1722:
			position, tokenIndex, depth = position1722, tokenIndex1722, depth1722
			return false
		},
		/* 151 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action121)> */
		func() bool {
			position1725, tokenIndex1725, depth1725 := position, tokenIndex, depth
			{
				position1726 := position
				depth++
				{
					position1727 := position
					depth++
					{
						position1728, tokenIndex1728, depth1728 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1729
						}
						position++
						goto l1728
					l1729:
						position, tokenIndex, depth = position1728, tokenIndex1728, depth1728
						if buffer[position] != rune('I') {
							goto l1725
						}
						position++
					}
				l1728:
					{
						position1730, tokenIndex1730, depth1730 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1731
						}
						position++
						goto l1730
					l1731:
						position, tokenIndex, depth = position1730, tokenIndex1730, depth1730
						if buffer[position] != rune('S') {
							goto l1725
						}
						position++
					}
				l1730:
					depth--
					add(rulePegText, position1727)
				}
				if !_rules[ruleAction121]() {
					goto l1725
				}
				depth--
				add(ruleIs, position1726)
			}
			return true
		l1725:
			position, tokenIndex, depth = position1725, tokenIndex1725, depth1725
			return false
		},
		/* 152 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action122)> */
		func() bool {
			position1732, tokenIndex1732, depth1732 := position, tokenIndex, depth
			{
				position1733 := position
				depth++
				{
					position1734 := position
					depth++
					{
						position1735, tokenIndex1735, depth1735 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1736
						}
						position++
						goto l1735
					l1736:
						position, tokenIndex, depth = position1735, tokenIndex1735, depth1735
						if buffer[position] != rune('I') {
							goto l1732
						}
						position++
					}
				l1735:
					{
						position1737, tokenIndex1737, depth1737 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1738
						}
						position++
						goto l1737
					l1738:
						position, tokenIndex, depth = position1737, tokenIndex1737, depth1737
						if buffer[position] != rune('S') {
							goto l1732
						}
						position++
					}
				l1737:
					if !_rules[rulesp]() {
						goto l1732
					}
					{
						position1739, tokenIndex1739, depth1739 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1740
						}
						position++
						goto l1739
					l1740:
						position, tokenIndex, depth = position1739, tokenIndex1739, depth1739
						if buffer[position] != rune('N') {
							goto l1732
						}
						position++
					}
				l1739:
					{
						position1741, tokenIndex1741, depth1741 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1742
						}
						position++
						goto l1741
					l1742:
						position, tokenIndex, depth = position1741, tokenIndex1741, depth1741
						if buffer[position] != rune('O') {
							goto l1732
						}
						position++
					}
				l1741:
					{
						position1743, tokenIndex1743, depth1743 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1744
						}
						position++
						goto l1743
					l1744:
						position, tokenIndex, depth = position1743, tokenIndex1743, depth1743
						if buffer[position] != rune('T') {
							goto l1732
						}
						position++
					}
				l1743:
					depth--
					add(rulePegText, position1734)
				}
				if !_rules[ruleAction122]() {
					goto l1732
				}
				depth--
				add(ruleIsNot, position1733)
			}
			return true
		l1732:
			position, tokenIndex, depth = position1732, tokenIndex1732, depth1732
			return false
		},
		/* 153 Plus <- <(<'+'> Action123)> */
		func() bool {
			position1745, tokenIndex1745, depth1745 := position, tokenIndex, depth
			{
				position1746 := position
				depth++
				{
					position1747 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1745
					}
					position++
					depth--
					add(rulePegText, position1747)
				}
				if !_rules[ruleAction123]() {
					goto l1745
				}
				depth--
				add(rulePlus, position1746)
			}
			return true
		l1745:
			position, tokenIndex, depth = position1745, tokenIndex1745, depth1745
			return false
		},
		/* 154 Minus <- <(<'-'> Action124)> */
		func() bool {
			position1748, tokenIndex1748, depth1748 := position, tokenIndex, depth
			{
				position1749 := position
				depth++
				{
					position1750 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1748
					}
					position++
					depth--
					add(rulePegText, position1750)
				}
				if !_rules[ruleAction124]() {
					goto l1748
				}
				depth--
				add(ruleMinus, position1749)
			}
			return true
		l1748:
			position, tokenIndex, depth = position1748, tokenIndex1748, depth1748
			return false
		},
		/* 155 Multiply <- <(<'*'> Action125)> */
		func() bool {
			position1751, tokenIndex1751, depth1751 := position, tokenIndex, depth
			{
				position1752 := position
				depth++
				{
					position1753 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1751
					}
					position++
					depth--
					add(rulePegText, position1753)
				}
				if !_rules[ruleAction125]() {
					goto l1751
				}
				depth--
				add(ruleMultiply, position1752)
			}
			return true
		l1751:
			position, tokenIndex, depth = position1751, tokenIndex1751, depth1751
			return false
		},
		/* 156 Divide <- <(<'/'> Action126)> */
		func() bool {
			position1754, tokenIndex1754, depth1754 := position, tokenIndex, depth
			{
				position1755 := position
				depth++
				{
					position1756 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1754
					}
					position++
					depth--
					add(rulePegText, position1756)
				}
				if !_rules[ruleAction126]() {
					goto l1754
				}
				depth--
				add(ruleDivide, position1755)
			}
			return true
		l1754:
			position, tokenIndex, depth = position1754, tokenIndex1754, depth1754
			return false
		},
		/* 157 Modulo <- <(<'%'> Action127)> */
		func() bool {
			position1757, tokenIndex1757, depth1757 := position, tokenIndex, depth
			{
				position1758 := position
				depth++
				{
					position1759 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1757
					}
					position++
					depth--
					add(rulePegText, position1759)
				}
				if !_rules[ruleAction127]() {
					goto l1757
				}
				depth--
				add(ruleModulo, position1758)
			}
			return true
		l1757:
			position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
			return false
		},
		/* 158 UnaryMinus <- <(<'-'> Action128)> */
		func() bool {
			position1760, tokenIndex1760, depth1760 := position, tokenIndex, depth
			{
				position1761 := position
				depth++
				{
					position1762 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1760
					}
					position++
					depth--
					add(rulePegText, position1762)
				}
				if !_rules[ruleAction128]() {
					goto l1760
				}
				depth--
				add(ruleUnaryMinus, position1761)
			}
			return true
		l1760:
			position, tokenIndex, depth = position1760, tokenIndex1760, depth1760
			return false
		},
		/* 159 Identifier <- <(<ident> Action129)> */
		func() bool {
			position1763, tokenIndex1763, depth1763 := position, tokenIndex, depth
			{
				position1764 := position
				depth++
				{
					position1765 := position
					depth++
					if !_rules[ruleident]() {
						goto l1763
					}
					depth--
					add(rulePegText, position1765)
				}
				if !_rules[ruleAction129]() {
					goto l1763
				}
				depth--
				add(ruleIdentifier, position1764)
			}
			return true
		l1763:
			position, tokenIndex, depth = position1763, tokenIndex1763, depth1763
			return false
		},
		/* 160 TargetIdentifier <- <(<('*' / jsonSetPath)> Action130)> */
		func() bool {
			position1766, tokenIndex1766, depth1766 := position, tokenIndex, depth
			{
				position1767 := position
				depth++
				{
					position1768 := position
					depth++
					{
						position1769, tokenIndex1769, depth1769 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l1770
						}
						position++
						goto l1769
					l1770:
						position, tokenIndex, depth = position1769, tokenIndex1769, depth1769
						if !_rules[rulejsonSetPath]() {
							goto l1766
						}
					}
				l1769:
					depth--
					add(rulePegText, position1768)
				}
				if !_rules[ruleAction130]() {
					goto l1766
				}
				depth--
				add(ruleTargetIdentifier, position1767)
			}
			return true
		l1766:
			position, tokenIndex, depth = position1766, tokenIndex1766, depth1766
			return false
		},
		/* 161 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1771, tokenIndex1771, depth1771 := position, tokenIndex, depth
			{
				position1772 := position
				depth++
				{
					position1773, tokenIndex1773, depth1773 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1774
					}
					position++
					goto l1773
				l1774:
					position, tokenIndex, depth = position1773, tokenIndex1773, depth1773
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1771
					}
					position++
				}
			l1773:
			l1775:
				{
					position1776, tokenIndex1776, depth1776 := position, tokenIndex, depth
					{
						position1777, tokenIndex1777, depth1777 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1778
						}
						position++
						goto l1777
					l1778:
						position, tokenIndex, depth = position1777, tokenIndex1777, depth1777
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1779
						}
						position++
						goto l1777
					l1779:
						position, tokenIndex, depth = position1777, tokenIndex1777, depth1777
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1780
						}
						position++
						goto l1777
					l1780:
						position, tokenIndex, depth = position1777, tokenIndex1777, depth1777
						if buffer[position] != rune('_') {
							goto l1776
						}
						position++
					}
				l1777:
					goto l1775
				l1776:
					position, tokenIndex, depth = position1776, tokenIndex1776, depth1776
				}
				depth--
				add(ruleident, position1772)
			}
			return true
		l1771:
			position, tokenIndex, depth = position1771, tokenIndex1771, depth1771
			return false
		},
		/* 162 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1781, tokenIndex1781, depth1781 := position, tokenIndex, depth
			{
				position1782 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1781
				}
			l1783:
				{
					position1784, tokenIndex1784, depth1784 := position, tokenIndex, depth
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1784
					}
					goto l1783
				l1784:
					position, tokenIndex, depth = position1784, tokenIndex1784, depth1784
				}
				depth--
				add(rulejsonGetPath, position1782)
			}
			return true
		l1781:
			position, tokenIndex, depth = position1781, tokenIndex1781, depth1781
			return false
		},
		/* 163 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1785, tokenIndex1785, depth1785 := position, tokenIndex, depth
			{
				position1786 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1785
				}
			l1787:
				{
					position1788, tokenIndex1788, depth1788 := position, tokenIndex, depth
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1788
					}
					goto l1787
				l1788:
					position, tokenIndex, depth = position1788, tokenIndex1788, depth1788
				}
				depth--
				add(rulejsonSetPath, position1786)
			}
			return true
		l1785:
			position, tokenIndex, depth = position1785, tokenIndex1785, depth1785
			return false
		},
		/* 164 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1789, tokenIndex1789, depth1789 := position, tokenIndex, depth
			{
				position1790 := position
				depth++
				{
					position1791, tokenIndex1791, depth1791 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1792
					}
					goto l1791
				l1792:
					position, tokenIndex, depth = position1791, tokenIndex1791, depth1791
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1789
					}
				}
			l1791:
				depth--
				add(rulejsonPathHead, position1790)
			}
			return true
		l1789:
			position, tokenIndex, depth = position1789, tokenIndex1789, depth1789
			return false
		},
		/* 165 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1793, tokenIndex1793, depth1793 := position, tokenIndex, depth
			{
				position1794 := position
				depth++
				{
					position1795, tokenIndex1795, depth1795 := position, tokenIndex, depth
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1796
					}
					goto l1795
				l1796:
					position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1797
					}
					goto l1795
				l1797:
					position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1798
					}
					goto l1795
				l1798:
					position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1799
					}
					goto l1795
				l1799:
					position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
					if !_rules[rulejsonArraySlice]() {
						goto l1800
					}
					goto l1795
				l1800:
					position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
					if !_rules[rulejsonArrayAccess]() {
						goto l1793
					}
				}
			l1795:
				depth--
				add(rulejsonGetPathNonHead, position1794)
			}
			return true
		l1793:
			position, tokenIndex, depth = position1793, tokenIndex1793, depth1793
			return false
		},
		/* 166 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1801, tokenIndex1801, depth1801 := position, tokenIndex, depth
			{
				position1802 := position
				depth++
				{
					position1803, tokenIndex1803, depth1803 := position, tokenIndex, depth
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1804
					}
					goto l1803
				l1804:
					position, tokenIndex, depth = position1803, tokenIndex1803, depth1803
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1801
					}
				}
			l1803:
				depth--
				add(rulejsonSetPathNonHead, position1802)
			}
			return true
		l1801:
			position, tokenIndex, depth = position1801, tokenIndex1801, depth1801
			return false
		},
		/* 167 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1805, tokenIndex1805, depth1805 := position, tokenIndex, depth
			{
				position1806 := position
				depth++
				{
					position1807, tokenIndex1807, depth1807 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1808
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1808
					}
					goto l1807
				l1808:
					position, tokenIndex, depth = position1807, tokenIndex1807, depth1807
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1805
					}
				}
			l1807:
				depth--
				add(rulejsonMapSingleLevel, position1806)
			}
			return true
		l1805:
			position, tokenIndex, depth = position1805, tokenIndex1805, depth1805
			return false
		},
		/* 168 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1809, tokenIndex1809, depth1809 := position, tokenIndex, depth
			{
				position1810 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1809
				}
				position++
				if buffer[position] != rune('.') {
					goto l1809
				}
				position++
				{
					position1811, tokenIndex1811, depth1811 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1812
					}
					goto l1811
				l1812:
					position, tokenIndex, depth = position1811, tokenIndex1811, depth1811
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1809
					}
				}
			l1811:
				depth--
				add(rulejsonMapMultipleLevel, position1810)
			}
			return true
		l1809:
			position, tokenIndex, depth = position1809, tokenIndex1809, depth1809
			return false
		},
		/* 169 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1813, tokenIndex1813, depth1813 := position, tokenIndex, depth
			{
				position1814 := position
				depth++
				{
					position1815 := position
					depth++
					{
						position1816, tokenIndex1816, depth1816 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1817
						}
						position++
						goto l1816
					l1817:
						position, tokenIndex, depth = position1816, tokenIndex1816, depth1816
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1813
						}
						position++
					}
				l1816:
				l1818:
					{
						position1819, tokenIndex1819, depth1819 := position, tokenIndex, depth
						{
							position1820, tokenIndex1820, depth1820 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1821
							}
							position++
							goto l1820
						l1821:
							position, tokenIndex, depth = position1820, tokenIndex1820, depth1820
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1822
							}
							position++
							goto l1820
						l1822:
							position, tokenIndex, depth = position1820, tokenIndex1820, depth1820
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1823
							}
							position++
							goto l1820
						l1823:
							position, tokenIndex, depth = position1820, tokenIndex1820, depth1820
							if buffer[position] != rune('_') {
								goto l1819
							}
							position++
						}
					l1820:
						goto l1818
					l1819:
						position, tokenIndex, depth = position1819, tokenIndex1819, depth1819
					}
					depth--
					add(rulePegText, position1815)
				}
				depth--
				add(rulejsonMapAccessString, position1814)
			}
			return true
		l1813:
			position, tokenIndex, depth = position1813, tokenIndex1813, depth1813
			return false
		},
		/* 170 jsonMapAccessBracket <- <('[' singleQuotedString ']')> */
		func() bool {
			position1824, tokenIndex1824, depth1824 := position, tokenIndex, depth
			{
				position1825 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1824
				}
				position++
				if !_rules[rulesingleQuotedString]() {
					goto l1824
				}
				if buffer[position] != rune(']') {
					goto l1824
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1825)
			}
			return true
		l1824:
			position, tokenIndex, depth = position1824, tokenIndex1824, depth1824
			return false
		},
		/* 171 singleQuotedString <- <('\'' <(('\'' '\'') / (!'\'' .))*> '\'')> */
		func() bool {
			position1826, tokenIndex1826, depth1826 := position, tokenIndex, depth
			{
				position1827 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l1826
				}
				position++
				{
					position1828 := position
					depth++
				l1829:
					{
						position1830, tokenIndex1830, depth1830 := position, tokenIndex, depth
						{
							position1831, tokenIndex1831, depth1831 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1832
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1832
							}
							position++
							goto l1831
						l1832:
							position, tokenIndex, depth = position1831, tokenIndex1831, depth1831
							{
								position1833, tokenIndex1833, depth1833 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1833
								}
								position++
								goto l1830
							l1833:
								position, tokenIndex, depth = position1833, tokenIndex1833, depth1833
							}
							if !matchDot() {
								goto l1830
							}
						}
					l1831:
						goto l1829
					l1830:
						position, tokenIndex, depth = position1830, tokenIndex1830, depth1830
					}
					depth--
					add(rulePegText, position1828)
				}
				if buffer[position] != rune('\'') {
					goto l1826
				}
				position++
				depth--
				add(rulesingleQuotedString, position1827)
			}
			return true
		l1826:
			position, tokenIndex, depth = position1826, tokenIndex1826, depth1826
			return false
		},
		/* 172 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1834, tokenIndex1834, depth1834 := position, tokenIndex, depth
			{
				position1835 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1834
				}
				position++
				{
					position1836 := position
					depth++
					{
						position1837, tokenIndex1837, depth1837 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1837
						}
						position++
						goto l1838
					l1837:
						position, tokenIndex, depth = position1837, tokenIndex1837, depth1837
					}
				l1838:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1834
					}
					position++
				l1839:
					{
						position1840, tokenIndex1840, depth1840 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1840
						}
						position++
						goto l1839
					l1840:
						position, tokenIndex, depth = position1840, tokenIndex1840, depth1840
					}
					depth--
					add(rulePegText, position1836)
				}
				if buffer[position] != rune(']') {
					goto l1834
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1835)
			}
			return true
		l1834:
			position, tokenIndex, depth = position1834, tokenIndex1834, depth1834
			return false
		},
		/* 173 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
		func() bool {
			position1841, tokenIndex1841, depth1841 := position, tokenIndex, depth
			{
				position1842 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1841
				}
				position++
				{
					position1843 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1841
					}
					position++
				l1844:
					{
						position1845, tokenIndex1845, depth1845 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1845
						}
						position++
						goto l1844
					l1845:
						position, tokenIndex, depth = position1845, tokenIndex1845, depth1845
					}
					depth--
					add(rulePegText, position1843)
				}
				if buffer[position] != rune(']') {
					goto l1841
				}
				position++
				depth--
				add(rulejsonNonNegativeArrayAccess, position1842)
			}
			return true
		l1841:
			position, tokenIndex, depth = position1841, tokenIndex1841, depth1841
			return false
		},
		/* 174 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1846, tokenIndex1846, depth1846 := position, tokenIndex, depth
			{
				position1847 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1846
				}
				position++
				{
					position1848 := position
					depth++
					{
						position1849, tokenIndex1849, depth1849 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1849
						}
						position++
						goto l1850
					l1849:
						position, tokenIndex, depth = position1849, tokenIndex1849, depth1849
					}
				l1850:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1846
					}
					position++
				l1851:
					{
						position1852, tokenIndex1852, depth1852 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1852
						}
						position++
						goto l1851
					l1852:
						position, tokenIndex, depth = position1852, tokenIndex1852, depth1852
					}
					if buffer[position] != rune(':') {
						goto l1846
					}
					position++
					{
						position1853, tokenIndex1853, depth1853 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1853
						}
						position++
						goto l1854
					l1853:
						position, tokenIndex, depth = position1853, tokenIndex1853, depth1853
					}
				l1854:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1846
					}
					position++
				l1855:
					{
						position1856, tokenIndex1856, depth1856 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1856
						}
						position++
						goto l1855
					l1856:
						position, tokenIndex, depth = position1856, tokenIndex1856, depth1856
					}
					{
						position1857, tokenIndex1857, depth1857 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1857
						}
						position++
						{
							position1859, tokenIndex1859, depth1859 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1859
							}
							position++
							goto l1860
						l1859:
							position, tokenIndex, depth = position1859, tokenIndex1859, depth1859
						}
					l1860:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1857
						}
						position++
					l1861:
						{
							position1862, tokenIndex1862, depth1862 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1862
							}
							position++
							goto l1861
						l1862:
							position, tokenIndex, depth = position1862, tokenIndex1862, depth1862
						}
						goto l1858
					l1857:
						position, tokenIndex, depth = position1857, tokenIndex1857, depth1857
					}
				l1858:
					depth--
					add(rulePegText, position1848)
				}
				if buffer[position] != rune(']') {
					goto l1846
				}
				position++
				depth--
				add(rulejsonArraySlice, position1847)
			}
			return true
		l1846:
			position, tokenIndex, depth = position1846, tokenIndex1846, depth1846
			return false
		},
		/* 175 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1863, tokenIndex1863, depth1863 := position, tokenIndex, depth
			{
				position1864 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1863
				}
				position++
				{
					position1865 := position
					depth++
					{
						position1866, tokenIndex1866, depth1866 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1867
						}
						position++
						{
							position1868, tokenIndex1868, depth1868 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1868
							}
							position++
							goto l1869
						l1868:
							position, tokenIndex, depth = position1868, tokenIndex1868, depth1868
						}
					l1869:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1867
						}
						position++
					l1870:
						{
							position1871, tokenIndex1871, depth1871 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1871
							}
							position++
							goto l1870
						l1871:
							position, tokenIndex, depth = position1871, tokenIndex1871, depth1871
						}
						goto l1866
					l1867:
						position, tokenIndex, depth = position1866, tokenIndex1866, depth1866
						{
							position1872, tokenIndex1872, depth1872 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1872
							}
							position++
							goto l1873
						l1872:
							position, tokenIndex, depth = position1872, tokenIndex1872, depth1872
						}
					l1873:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1863
						}
						position++
					l1874:
						{
							position1875, tokenIndex1875, depth1875 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1875
							}
							position++
							goto l1874
						l1875:
							position, tokenIndex, depth = position1875, tokenIndex1875, depth1875
						}
						if buffer[position] != rune(':') {
							goto l1863
						}
						position++
					}
				l1866:
					depth--
					add(rulePegText, position1865)
				}
				if buffer[position] != rune(']') {
					goto l1863
				}
				position++
				depth--
				add(rulejsonArrayPartialSlice, position1864)
			}
			return true
		l1863:
			position, tokenIndex, depth = position1863, tokenIndex1863, depth1863
			return false
		},
		/* 176 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1876, tokenIndex1876, depth1876 := position, tokenIndex, depth
			{
				position1877 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1876
				}
				position++
				if buffer[position] != rune(':') {
					goto l1876
				}
				position++
				if buffer[position] != rune(']') {
					goto l1876
				}
				position++
				depth--
				add(rulejsonArrayFullSlice, position1877)
			}
			return true
		l1876:
			position, tokenIndex, depth = position1876, tokenIndex1876, depth1876
			return false
		},
		/* 177 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1878, tokenIndex1878, depth1878 := position, tokenIndex, depth
			{
				position1879 := position
				depth++
				{
					position1880, tokenIndex1880, depth1880 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l1881
					}
					position++
					goto l1880
				l1881:
					position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
					if buffer[position] != rune('\t') {
						goto l1882
					}
					position++
					goto l1880
				l1882:
					position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
					if buffer[position] != rune('\n') {
						goto l1883
					}
					position++
					goto l1880
				l1883:
					position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
					if buffer[position] != rune('\r') {
						goto l1884
					}
					position++
					goto l1880
				l1884:
					position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
					if !_rules[rulecomment]() {
						goto l1885
					}
					goto l1880
				l1885:
					position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
					if !_rules[rulefinalComment]() {
						goto l1878
					}
				}
			l1880:
				depth--
				add(rulespElem, position1879)
			}
			return true
		l1878:
			position, tokenIndex, depth = position1878, tokenIndex1878, depth1878
			return false
		},
		/* 178 sp <- <spElem+> */
		func() bool {
			position1886, tokenIndex1886, depth1886 := position, tokenIndex, depth
			{
				position1887 := position
				depth++
				if !_rules[rulespElem]() {
					goto l1886
				}
			l1888:
				{
					position1889, tokenIndex1889, depth1889 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1889
					}
					goto l1888
				l1889:
					position, tokenIndex, depth = position1889, tokenIndex1889, depth1889
				}
				depth--
				add(rulesp, position1887)
			}
			return true
		l1886:
			position, tokenIndex, depth = position1886, tokenIndex1886, depth1886
			return false
		},
		/* 179 spOpt <- <spElem*> */
		func() bool {
			{
				position1891 := position
				depth++
			l1892:
				{
					position1893, tokenIndex1893, depth1893 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1893
					}
					goto l1892
				l1893:
					position, tokenIndex, depth = position1893, tokenIndex1893, depth1893
				}
				depth--
				add(rulespOpt, position1891)
			}
			return true
		},
		/* 180 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1894, tokenIndex1894, depth1894 := position, tokenIndex, depth
			{
				position1895 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1894
				}
				position++
				if buffer[position] != rune('-') {
					goto l1894
				}
				position++
			l1896:
				{
					position1897, tokenIndex1897, depth1897 := position, tokenIndex, depth
					{
						position1898, tokenIndex1898, depth1898 := position, tokenIndex, depth
						{
							position1899, tokenIndex1899, depth1899 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1900
							}
							position++
							goto l1899
						l1900:
							position, tokenIndex, depth = position1899, tokenIndex1899, depth1899
							if buffer[position] != rune('\n') {
								goto l1898
							}
							position++
						}
					l1899:
						goto l1897
					l1898:
						position, tokenIndex, depth = position1898, tokenIndex1898, depth1898
					}
					if !matchDot() {
						goto l1897
					}
					goto l1896
				l1897:
					position, tokenIndex, depth = position1897, tokenIndex1897, depth1897
				}
				{
					position1901, tokenIndex1901, depth1901 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1902
					}
					position++
					goto l1901
				l1902:
					position, tokenIndex, depth = position1901, tokenIndex1901, depth1901
					if buffer[position] != rune('\n') {
						goto l1894
					}
					position++
				}
			l1901:
				depth--
				add(rulecomment, position1895)
			}
			return true
		l1894:
			position, tokenIndex, depth = position1894, tokenIndex1894, depth1894
			return false
		},
		/* 181 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1903, tokenIndex1903, depth1903 := position, tokenIndex, depth
			{
				position1904 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1903
				}
				position++
				if buffer[position] != rune('-') {
					goto l1903
				}
				position++
			l1905:
				{
					position1906, tokenIndex1906, depth1906 := position, tokenIndex, depth
					{
						position1907, tokenIndex1907, depth1907 := position, tokenIndex, depth
						{
							position1908, tokenIndex1908, depth1908 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1909
							}
							position++
							goto l1908
						l1909:
							position, tokenIndex, depth = position1908, tokenIndex1908, depth1908
							if buffer[position] != rune('\n') {
								goto l1907
							}
							position++
						}
					l1908:
						goto l1906
					l1907:
						position, tokenIndex, depth = position1907, tokenIndex1907, depth1907
					}
					if !matchDot() {
						goto l1906
					}
					goto l1905
				l1906:
					position, tokenIndex, depth = position1906, tokenIndex1906, depth1906
				}
				{
					position1910, tokenIndex1910, depth1910 := position, tokenIndex, depth
					if !matchDot() {
						goto l1910
					}
					goto l1903
				l1910:
					position, tokenIndex, depth = position1910, tokenIndex1910, depth1910
				}
				depth--
				add(rulefinalComment, position1904)
			}
			return true
		l1903:
			position, tokenIndex, depth = position1903, tokenIndex1903, depth1903
			return false
		},
		nil,
		/* 184 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 185 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 186 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 187 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 188 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 189 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 190 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 191 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 192 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 193 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 194 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 195 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 196 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 197 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 198 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 199 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 200 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 201 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 202 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 203 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 204 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 205 Action21 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 206 Action22 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 207 Action23 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 208 Action24 <- <{
		    p.AssembleEval(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 209 Action25 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 210 Action26 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 211 Action27 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 212 Action28 <- <{
		    p.AssembleEmitterSampling(CountBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 213 Action29 <- <{
		    p.AssembleEmitterSampling(RandomizedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 214 Action30 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 215 Action31 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 0.001)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 216 Action32 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 217 Action33 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 218 Action34 <- <{
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
		/* 219 Action35 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 220 Action36 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 221 Action37 <- <{
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
		/* 222 Action38 <- <{
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
		/* 223 Action39 <- <{
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
		/* 224 Action40 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 225 Action41 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 226 Action42 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 227 Action43 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 228 Action44 <- <{
		    p.EnsureCapacitySpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 229 Action45 <- <{
		    p.EnsureSheddingSpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 230 Action46 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 231 Action47 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 232 Action48 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 233 Action49 <- <{
		    p.EnsureIdentifier(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 234 Action50 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 235 Action51 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 236 Action52 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 237 Action53 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 238 Action54 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 239 Action55 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 240 Action56 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 241 Action57 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 242 Action58 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 243 Action59 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 244 Action60 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 245 Action61 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 246 Action62 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 247 Action63 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 248 Action64 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 249 Action65 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 250 Action66 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 251 Action67 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 252 Action68 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 253 Action69 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 254 Action70 <- <{
		    p.AssembleSortedExpression()
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 255 Action71 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 256 Action72 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 257 Action73 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 258 Action74 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 259 Action75 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 260 Action76 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 261 Action77 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 262 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 263 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 264 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 265 Action81 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 266 Action82 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 267 Action83 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 268 Action84 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 269 Action85 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 270 Action86 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 271 Action87 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 272 Action88 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 273 Action89 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 274 Action90 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 275 Action91 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 276 Action92 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 277 Action93 <- <{
		    p.PushComponent(begin, end, Wait)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 278 Action94 <- <{
		    p.PushComponent(begin, end, DropOldest)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 279 Action95 <- <{
		    p.PushComponent(begin, end, DropNewest)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 280 Action96 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 281 Action97 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 282 Action98 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 283 Action99 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 284 Action100 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 285 Action101 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 286 Action102 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 287 Action103 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 288 Action104 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 289 Action105 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 290 Action106 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 291 Action107 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 292 Action108 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 293 Action109 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 294 Action110 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction110, position)
			}
			return true
		},
		/* 295 Action111 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction111, position)
			}
			return true
		},
		/* 296 Action112 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction112, position)
			}
			return true
		},
		/* 297 Action113 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction113, position)
			}
			return true
		},
		/* 298 Action114 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction114, position)
			}
			return true
		},
		/* 299 Action115 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction115, position)
			}
			return true
		},
		/* 300 Action116 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction116, position)
			}
			return true
		},
		/* 301 Action117 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction117, position)
			}
			return true
		},
		/* 302 Action118 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction118, position)
			}
			return true
		},
		/* 303 Action119 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction119, position)
			}
			return true
		},
		/* 304 Action120 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction120, position)
			}
			return true
		},
		/* 305 Action121 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction121, position)
			}
			return true
		},
		/* 306 Action122 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction122, position)
			}
			return true
		},
		/* 307 Action123 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction123, position)
			}
			return true
		},
		/* 308 Action124 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction124, position)
			}
			return true
		},
		/* 309 Action125 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction125, position)
			}
			return true
		},
		/* 310 Action126 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction126, position)
			}
			return true
		},
		/* 311 Action127 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction127, position)
			}
			return true
		},
		/* 312 Action128 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction128, position)
			}
			return true
		},
		/* 313 Action129 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction129, position)
			}
			return true
		},
		/* 314 Action130 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction130, position)
			}
			return true
		},
	}
	p.rules = _rules
}

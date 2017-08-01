package parser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

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
	ruleFuncAppSelector
	ruleFuncElemAccessor
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
	ruleMissing
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
	ruledoubleQuotedString
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
	ruleAction134
	ruleAction135
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
	"FuncAppSelector",
	"FuncElemAccessor",
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
	"Missing",
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
	"doubleQuotedString",
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
	"Action134",
	"Action135",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
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
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type bqlPegBackend struct {
	parseStack

	Buffer string
	buffer []rune
	rules  [326]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *bqlPegBackend) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *bqlPegBackend) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
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
	p   *bqlPegBackend
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *bqlPegBackend) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *bqlPegBackend) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
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

			p.AssembleInsertIntoFrom()

		case ruleAction13:

			p.AssemblePauseSource()

		case ruleAction14:

			p.AssembleResumeSource()

		case ruleAction15:

			p.AssembleRewindSource()

		case ruleAction16:

			p.AssembleDropSource()

		case ruleAction17:

			p.AssembleDropStream()

		case ruleAction18:

			p.AssembleDropSink()

		case ruleAction19:

			p.AssembleDropState()

		case ruleAction20:

			p.AssembleLoadState()

		case ruleAction21:

			p.AssembleLoadStateOrCreate()

		case ruleAction22:

			p.AssembleSaveState()

		case ruleAction23:

			p.AssembleEval(begin, end)

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

			p.EnsureCapacitySpec(begin, end)

		case ruleAction44:

			p.EnsureSheddingSpec(begin, end)

		case ruleAction45:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction46:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction47:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction48:

			p.EnsureIdentifier(begin, end)

		case ruleAction49:

			p.AssembleSourceSinkParam()

		case ruleAction50:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction51:

			p.AssembleMap(begin, end)

		case ruleAction52:

			p.AssembleKeyValuePair()

		case ruleAction53:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction54:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction55:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction56:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction57:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction58:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction59:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction60:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction61:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction62:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction63:

			p.AssembleTypeCast(begin, end)

		case ruleAction64:

			p.AssembleTypeCast(begin, end)

		case ruleAction65:

			p.AssembleFuncAppSelector()

		case ruleAction66:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRaw(substr))

		case ruleAction67:

			p.AssembleFuncApp()

		case ruleAction68:

			p.AssembleExpressions(begin, end)
			p.AssembleFuncApp()

		case ruleAction69:

			p.AssembleExpressions(begin, end)

		case ruleAction70:

			p.AssembleExpressions(begin, end)

		case ruleAction71:

			p.AssembleSortedExpression()

		case ruleAction72:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction73:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction74:

			p.AssembleMap(begin, end)

		case ruleAction75:

			p.AssembleKeyValuePair()

		case ruleAction76:

			p.AssembleConditionCase(begin, end)

		case ruleAction77:

			p.AssembleExpressionCase(begin, end)

		case ruleAction78:

			p.AssembleWhenThenPair()

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction80:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction81:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction82:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction83:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction84:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction85:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction86:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction87:

			p.PushComponent(begin, end, NewMissing())

		case ruleAction88:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction89:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction90:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction91:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction92:

			p.PushComponent(begin, end, Istream)

		case ruleAction93:

			p.PushComponent(begin, end, Dstream)

		case ruleAction94:

			p.PushComponent(begin, end, Rstream)

		case ruleAction95:

			p.PushComponent(begin, end, Tuples)

		case ruleAction96:

			p.PushComponent(begin, end, Seconds)

		case ruleAction97:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction98:

			p.PushComponent(begin, end, Wait)

		case ruleAction99:

			p.PushComponent(begin, end, DropOldest)

		case ruleAction100:

			p.PushComponent(begin, end, DropNewest)

		case ruleAction101:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction102:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction103:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction104:

			p.PushComponent(begin, end, Yes)

		case ruleAction105:

			p.PushComponent(begin, end, No)

		case ruleAction106:

			p.PushComponent(begin, end, Yes)

		case ruleAction107:

			p.PushComponent(begin, end, No)

		case ruleAction108:

			p.PushComponent(begin, end, Bool)

		case ruleAction109:

			p.PushComponent(begin, end, Int)

		case ruleAction110:

			p.PushComponent(begin, end, Float)

		case ruleAction111:

			p.PushComponent(begin, end, String)

		case ruleAction112:

			p.PushComponent(begin, end, Blob)

		case ruleAction113:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction114:

			p.PushComponent(begin, end, Array)

		case ruleAction115:

			p.PushComponent(begin, end, Map)

		case ruleAction116:

			p.PushComponent(begin, end, Or)

		case ruleAction117:

			p.PushComponent(begin, end, And)

		case ruleAction118:

			p.PushComponent(begin, end, Not)

		case ruleAction119:

			p.PushComponent(begin, end, Equal)

		case ruleAction120:

			p.PushComponent(begin, end, Less)

		case ruleAction121:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction122:

			p.PushComponent(begin, end, Greater)

		case ruleAction123:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction124:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction125:

			p.PushComponent(begin, end, Concat)

		case ruleAction126:

			p.PushComponent(begin, end, Is)

		case ruleAction127:

			p.PushComponent(begin, end, IsNot)

		case ruleAction128:

			p.PushComponent(begin, end, Plus)

		case ruleAction129:

			p.PushComponent(begin, end, Minus)

		case ruleAction130:

			p.PushComponent(begin, end, Multiply)

		case ruleAction131:

			p.PushComponent(begin, end, Divide)

		case ruleAction132:

			p.PushComponent(begin, end, Modulo)

		case ruleAction133:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction134:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction135:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *bqlPegBackend) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
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
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rulespOpt]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					if !_rules[ruleStatementWithRest]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex = position2, tokenIndex2
					if !_rules[ruleStatementWithoutRest]() {
						goto l0
					}
				}
			l2:
				{
					position4, tokenIndex4 := position, tokenIndex
					if !matchDot() {
						goto l4
					}
					goto l0
				l4:
					position, tokenIndex = position4, tokenIndex4
				}
				add(ruleSingleStatement, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 StatementWithRest <- <(<(Statement spOpt ';' spOpt)> .* Action0)> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
				{
					position7 := position
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
					add(rulePegText, position7)
				}
			l8:
				{
					position9, tokenIndex9 := position, tokenIndex
					if !matchDot() {
						goto l9
					}
					goto l8
				l9:
					position, tokenIndex = position9, tokenIndex9
				}
				if !_rules[ruleAction0]() {
					goto l5
				}
				add(ruleStatementWithRest, position6)
			}
			return true
		l5:
			position, tokenIndex = position5, tokenIndex5
			return false
		},
		/* 2 StatementWithoutRest <- <(<(Statement spOpt)> Action1)> */
		func() bool {
			position10, tokenIndex10 := position, tokenIndex
			{
				position11 := position
				{
					position12 := position
					if !_rules[ruleStatement]() {
						goto l10
					}
					if !_rules[rulespOpt]() {
						goto l10
					}
					add(rulePegText, position12)
				}
				if !_rules[ruleAction1]() {
					goto l10
				}
				add(ruleStatementWithoutRest, position11)
			}
			return true
		l10:
			position, tokenIndex = position10, tokenIndex10
			return false
		},
		/* 3 Statement <- <(SelectUnionStmt / SelectStmt / SourceStmt / SinkStmt / StateStmt / StreamStmt / EvalStmt)> */
		func() bool {
			position13, tokenIndex13 := position, tokenIndex
			{
				position14 := position
				{
					position15, tokenIndex15 := position, tokenIndex
					if !_rules[ruleSelectUnionStmt]() {
						goto l16
					}
					goto l15
				l16:
					position, tokenIndex = position15, tokenIndex15
					if !_rules[ruleSelectStmt]() {
						goto l17
					}
					goto l15
				l17:
					position, tokenIndex = position15, tokenIndex15
					if !_rules[ruleSourceStmt]() {
						goto l18
					}
					goto l15
				l18:
					position, tokenIndex = position15, tokenIndex15
					if !_rules[ruleSinkStmt]() {
						goto l19
					}
					goto l15
				l19:
					position, tokenIndex = position15, tokenIndex15
					if !_rules[ruleStateStmt]() {
						goto l20
					}
					goto l15
				l20:
					position, tokenIndex = position15, tokenIndex15
					if !_rules[ruleStreamStmt]() {
						goto l21
					}
					goto l15
				l21:
					position, tokenIndex = position15, tokenIndex15
					if !_rules[ruleEvalStmt]() {
						goto l13
					}
				}
			l15:
				add(ruleStatement, position14)
			}
			return true
		l13:
			position, tokenIndex = position13, tokenIndex13
			return false
		},
		/* 4 SourceStmt <- <(CreateSourceStmt / UpdateSourceStmt / DropSourceStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt)> */
		func() bool {
			position22, tokenIndex22 := position, tokenIndex
			{
				position23 := position
				{
					position24, tokenIndex24 := position, tokenIndex
					if !_rules[ruleCreateSourceStmt]() {
						goto l25
					}
					goto l24
				l25:
					position, tokenIndex = position24, tokenIndex24
					if !_rules[ruleUpdateSourceStmt]() {
						goto l26
					}
					goto l24
				l26:
					position, tokenIndex = position24, tokenIndex24
					if !_rules[ruleDropSourceStmt]() {
						goto l27
					}
					goto l24
				l27:
					position, tokenIndex = position24, tokenIndex24
					if !_rules[rulePauseSourceStmt]() {
						goto l28
					}
					goto l24
				l28:
					position, tokenIndex = position24, tokenIndex24
					if !_rules[ruleResumeSourceStmt]() {
						goto l29
					}
					goto l24
				l29:
					position, tokenIndex = position24, tokenIndex24
					if !_rules[ruleRewindSourceStmt]() {
						goto l22
					}
				}
			l24:
				add(ruleSourceStmt, position23)
			}
			return true
		l22:
			position, tokenIndex = position22, tokenIndex22
			return false
		},
		/* 5 SinkStmt <- <(CreateSinkStmt / UpdateSinkStmt / DropSinkStmt)> */
		func() bool {
			position30, tokenIndex30 := position, tokenIndex
			{
				position31 := position
				{
					position32, tokenIndex32 := position, tokenIndex
					if !_rules[ruleCreateSinkStmt]() {
						goto l33
					}
					goto l32
				l33:
					position, tokenIndex = position32, tokenIndex32
					if !_rules[ruleUpdateSinkStmt]() {
						goto l34
					}
					goto l32
				l34:
					position, tokenIndex = position32, tokenIndex32
					if !_rules[ruleDropSinkStmt]() {
						goto l30
					}
				}
			l32:
				add(ruleSinkStmt, position31)
			}
			return true
		l30:
			position, tokenIndex = position30, tokenIndex30
			return false
		},
		/* 6 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt / LoadStateOrCreateStmt / LoadStateStmt / SaveStateStmt)> */
		func() bool {
			position35, tokenIndex35 := position, tokenIndex
			{
				position36 := position
				{
					position37, tokenIndex37 := position, tokenIndex
					if !_rules[ruleCreateStateStmt]() {
						goto l38
					}
					goto l37
				l38:
					position, tokenIndex = position37, tokenIndex37
					if !_rules[ruleUpdateStateStmt]() {
						goto l39
					}
					goto l37
				l39:
					position, tokenIndex = position37, tokenIndex37
					if !_rules[ruleDropStateStmt]() {
						goto l40
					}
					goto l37
				l40:
					position, tokenIndex = position37, tokenIndex37
					if !_rules[ruleLoadStateOrCreateStmt]() {
						goto l41
					}
					goto l37
				l41:
					position, tokenIndex = position37, tokenIndex37
					if !_rules[ruleLoadStateStmt]() {
						goto l42
					}
					goto l37
				l42:
					position, tokenIndex = position37, tokenIndex37
					if !_rules[ruleSaveStateStmt]() {
						goto l35
					}
				}
			l37:
				add(ruleStateStmt, position36)
			}
			return true
		l35:
			position, tokenIndex = position35, tokenIndex35
			return false
		},
		/* 7 StreamStmt <- <(CreateStreamAsSelectUnionStmt / CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoFromStmt)> */
		func() bool {
			position43, tokenIndex43 := position, tokenIndex
			{
				position44 := position
				{
					position45, tokenIndex45 := position, tokenIndex
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l46
					}
					goto l45
				l46:
					position, tokenIndex = position45, tokenIndex45
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l47
					}
					goto l45
				l47:
					position, tokenIndex = position45, tokenIndex45
					if !_rules[ruleDropStreamStmt]() {
						goto l48
					}
					goto l45
				l48:
					position, tokenIndex = position45, tokenIndex45
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l43
					}
				}
			l45:
				add(ruleStreamStmt, position44)
			}
			return true
		l43:
			position, tokenIndex = position43, tokenIndex43
			return false
		},
		/* 8 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') Emitter Projections WindowedFrom Filter Grouping Having Action2)> */
		func() bool {
			position49, tokenIndex49 := position, tokenIndex
			{
				position50 := position
				{
					position51, tokenIndex51 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l52
					}
					position++
					goto l51
				l52:
					position, tokenIndex = position51, tokenIndex51
					if buffer[position] != rune('S') {
						goto l49
					}
					position++
				}
			l51:
				{
					position53, tokenIndex53 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex = position53, tokenIndex53
					if buffer[position] != rune('E') {
						goto l49
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex = position55, tokenIndex55
					if buffer[position] != rune('L') {
						goto l49
					}
					position++
				}
			l55:
				{
					position57, tokenIndex57 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex = position57, tokenIndex57
					if buffer[position] != rune('E') {
						goto l49
					}
					position++
				}
			l57:
				{
					position59, tokenIndex59 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l60
					}
					position++
					goto l59
				l60:
					position, tokenIndex = position59, tokenIndex59
					if buffer[position] != rune('C') {
						goto l49
					}
					position++
				}
			l59:
				{
					position61, tokenIndex61 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l62
					}
					position++
					goto l61
				l62:
					position, tokenIndex = position61, tokenIndex61
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
				add(ruleSelectStmt, position50)
			}
			return true
		l49:
			position, tokenIndex = position49, tokenIndex49
			return false
		},
		/* 9 SelectUnionStmt <- <(<(SelectStmt (sp (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N')) sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action3)> */
		func() bool {
			position63, tokenIndex63 := position, tokenIndex
			{
				position64 := position
				{
					position65 := position
					if !_rules[ruleSelectStmt]() {
						goto l63
					}
					if !_rules[rulesp]() {
						goto l63
					}
					{
						position68, tokenIndex68 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l69
						}
						position++
						goto l68
					l69:
						position, tokenIndex = position68, tokenIndex68
						if buffer[position] != rune('U') {
							goto l63
						}
						position++
					}
				l68:
					{
						position70, tokenIndex70 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l71
						}
						position++
						goto l70
					l71:
						position, tokenIndex = position70, tokenIndex70
						if buffer[position] != rune('N') {
							goto l63
						}
						position++
					}
				l70:
					{
						position72, tokenIndex72 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l73
						}
						position++
						goto l72
					l73:
						position, tokenIndex = position72, tokenIndex72
						if buffer[position] != rune('I') {
							goto l63
						}
						position++
					}
				l72:
					{
						position74, tokenIndex74 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l75
						}
						position++
						goto l74
					l75:
						position, tokenIndex = position74, tokenIndex74
						if buffer[position] != rune('O') {
							goto l63
						}
						position++
					}
				l74:
					{
						position76, tokenIndex76 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l77
						}
						position++
						goto l76
					l77:
						position, tokenIndex = position76, tokenIndex76
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
						position78, tokenIndex78 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l79
						}
						position++
						goto l78
					l79:
						position, tokenIndex = position78, tokenIndex78
						if buffer[position] != rune('A') {
							goto l63
						}
						position++
					}
				l78:
					{
						position80, tokenIndex80 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l81
						}
						position++
						goto l80
					l81:
						position, tokenIndex = position80, tokenIndex80
						if buffer[position] != rune('L') {
							goto l63
						}
						position++
					}
				l80:
					{
						position82, tokenIndex82 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l83
						}
						position++
						goto l82
					l83:
						position, tokenIndex = position82, tokenIndex82
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
						position67, tokenIndex67 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l67
						}
						{
							position84, tokenIndex84 := position, tokenIndex
							if buffer[position] != rune('u') {
								goto l85
							}
							position++
							goto l84
						l85:
							position, tokenIndex = position84, tokenIndex84
							if buffer[position] != rune('U') {
								goto l67
							}
							position++
						}
					l84:
						{
							position86, tokenIndex86 := position, tokenIndex
							if buffer[position] != rune('n') {
								goto l87
							}
							position++
							goto l86
						l87:
							position, tokenIndex = position86, tokenIndex86
							if buffer[position] != rune('N') {
								goto l67
							}
							position++
						}
					l86:
						{
							position88, tokenIndex88 := position, tokenIndex
							if buffer[position] != rune('i') {
								goto l89
							}
							position++
							goto l88
						l89:
							position, tokenIndex = position88, tokenIndex88
							if buffer[position] != rune('I') {
								goto l67
							}
							position++
						}
					l88:
						{
							position90, tokenIndex90 := position, tokenIndex
							if buffer[position] != rune('o') {
								goto l91
							}
							position++
							goto l90
						l91:
							position, tokenIndex = position90, tokenIndex90
							if buffer[position] != rune('O') {
								goto l67
							}
							position++
						}
					l90:
						{
							position92, tokenIndex92 := position, tokenIndex
							if buffer[position] != rune('n') {
								goto l93
							}
							position++
							goto l92
						l93:
							position, tokenIndex = position92, tokenIndex92
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
							position94, tokenIndex94 := position, tokenIndex
							if buffer[position] != rune('a') {
								goto l95
							}
							position++
							goto l94
						l95:
							position, tokenIndex = position94, tokenIndex94
							if buffer[position] != rune('A') {
								goto l67
							}
							position++
						}
					l94:
						{
							position96, tokenIndex96 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l97
							}
							position++
							goto l96
						l97:
							position, tokenIndex = position96, tokenIndex96
							if buffer[position] != rune('L') {
								goto l67
							}
							position++
						}
					l96:
						{
							position98, tokenIndex98 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l99
							}
							position++
							goto l98
						l99:
							position, tokenIndex = position98, tokenIndex98
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
						position, tokenIndex = position67, tokenIndex67
					}
					add(rulePegText, position65)
				}
				if !_rules[ruleAction3]() {
					goto l63
				}
				add(ruleSelectUnionStmt, position64)
			}
			return true
		l63:
			position, tokenIndex = position63, tokenIndex63
			return false
		},
		/* 10 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action4)> */
		func() bool {
			position100, tokenIndex100 := position, tokenIndex
			{
				position101 := position
				{
					position102, tokenIndex102 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex = position102, tokenIndex102
					if buffer[position] != rune('C') {
						goto l100
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex = position104, tokenIndex104
					if buffer[position] != rune('R') {
						goto l100
					}
					position++
				}
			l104:
				{
					position106, tokenIndex106 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex = position106, tokenIndex106
					if buffer[position] != rune('E') {
						goto l100
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex = position108, tokenIndex108
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex = position110, tokenIndex110
					if buffer[position] != rune('T') {
						goto l100
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex = position112, tokenIndex112
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
					position114, tokenIndex114 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex = position114, tokenIndex114
					if buffer[position] != rune('S') {
						goto l100
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex = position116, tokenIndex116
					if buffer[position] != rune('T') {
						goto l100
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex = position118, tokenIndex118
					if buffer[position] != rune('R') {
						goto l100
					}
					position++
				}
			l118:
				{
					position120, tokenIndex120 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex = position120, tokenIndex120
					if buffer[position] != rune('E') {
						goto l100
					}
					position++
				}
			l120:
				{
					position122, tokenIndex122 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex = position122, tokenIndex122
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex = position124, tokenIndex124
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
					position126, tokenIndex126 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex = position126, tokenIndex126
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex = position128, tokenIndex128
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
				add(ruleCreateStreamAsSelectStmt, position101)
			}
			return true
		l100:
			position, tokenIndex = position100, tokenIndex100
			return false
		},
		/* 11 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action5)> */
		func() bool {
			position130, tokenIndex130 := position, tokenIndex
			{
				position131 := position
				{
					position132, tokenIndex132 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex = position132, tokenIndex132
					if buffer[position] != rune('C') {
						goto l130
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex = position134, tokenIndex134
					if buffer[position] != rune('R') {
						goto l130
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex = position136, tokenIndex136
					if buffer[position] != rune('E') {
						goto l130
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex = position138, tokenIndex138
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex = position140, tokenIndex140
					if buffer[position] != rune('T') {
						goto l130
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex = position142, tokenIndex142
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
					position144, tokenIndex144 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex = position144, tokenIndex144
					if buffer[position] != rune('S') {
						goto l130
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex = position146, tokenIndex146
					if buffer[position] != rune('T') {
						goto l130
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex = position148, tokenIndex148
					if buffer[position] != rune('R') {
						goto l130
					}
					position++
				}
			l148:
				{
					position150, tokenIndex150 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex = position150, tokenIndex150
					if buffer[position] != rune('E') {
						goto l130
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex = position152, tokenIndex152
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex = position154, tokenIndex154
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
					position156, tokenIndex156 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex = position156, tokenIndex156
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l156:
				{
					position158, tokenIndex158 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex = position158, tokenIndex158
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
				add(ruleCreateStreamAsSelectUnionStmt, position131)
			}
			return true
		l130:
			position, tokenIndex = position130, tokenIndex130
			return false
		},
		/* 12 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action6)> */
		func() bool {
			position160, tokenIndex160 := position, tokenIndex
			{
				position161 := position
				{
					position162, tokenIndex162 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex = position162, tokenIndex162
					if buffer[position] != rune('C') {
						goto l160
					}
					position++
				}
			l162:
				{
					position164, tokenIndex164 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex = position164, tokenIndex164
					if buffer[position] != rune('R') {
						goto l160
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex = position166, tokenIndex166
					if buffer[position] != rune('E') {
						goto l160
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex = position168, tokenIndex168
					if buffer[position] != rune('A') {
						goto l160
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex = position170, tokenIndex170
					if buffer[position] != rune('T') {
						goto l160
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex = position172, tokenIndex172
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
					position174, tokenIndex174 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex = position174, tokenIndex174
					if buffer[position] != rune('S') {
						goto l160
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex = position176, tokenIndex176
					if buffer[position] != rune('O') {
						goto l160
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex = position178, tokenIndex178
					if buffer[position] != rune('U') {
						goto l160
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex = position180, tokenIndex180
					if buffer[position] != rune('R') {
						goto l160
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex = position182, tokenIndex182
					if buffer[position] != rune('C') {
						goto l160
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex = position184, tokenIndex184
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
					position186, tokenIndex186 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex = position186, tokenIndex186
					if buffer[position] != rune('T') {
						goto l160
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex = position188, tokenIndex188
					if buffer[position] != rune('Y') {
						goto l160
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex = position190, tokenIndex190
					if buffer[position] != rune('P') {
						goto l160
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex = position192, tokenIndex192
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
				add(ruleCreateSourceStmt, position161)
			}
			return true
		l160:
			position, tokenIndex = position160, tokenIndex160
			return false
		},
		/* 13 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action7)> */
		func() bool {
			position194, tokenIndex194 := position, tokenIndex
			{
				position195 := position
				{
					position196, tokenIndex196 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex = position196, tokenIndex196
					if buffer[position] != rune('C') {
						goto l194
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex = position198, tokenIndex198
					if buffer[position] != rune('R') {
						goto l194
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex = position200, tokenIndex200
					if buffer[position] != rune('E') {
						goto l194
					}
					position++
				}
			l200:
				{
					position202, tokenIndex202 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex = position202, tokenIndex202
					if buffer[position] != rune('A') {
						goto l194
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex = position204, tokenIndex204
					if buffer[position] != rune('T') {
						goto l194
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex = position206, tokenIndex206
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
					position208, tokenIndex208 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex = position208, tokenIndex208
					if buffer[position] != rune('S') {
						goto l194
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex = position210, tokenIndex210
					if buffer[position] != rune('I') {
						goto l194
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex = position212, tokenIndex212
					if buffer[position] != rune('N') {
						goto l194
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214 := position, tokenIndex
					if buffer[position] != rune('k') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex = position214, tokenIndex214
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
					position216, tokenIndex216 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex = position216, tokenIndex216
					if buffer[position] != rune('T') {
						goto l194
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex = position218, tokenIndex218
					if buffer[position] != rune('Y') {
						goto l194
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex = position220, tokenIndex220
					if buffer[position] != rune('P') {
						goto l194
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex = position222, tokenIndex222
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
				add(ruleCreateSinkStmt, position195)
			}
			return true
		l194:
			position, tokenIndex = position194, tokenIndex194
			return false
		},
		/* 14 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action8)> */
		func() bool {
			position224, tokenIndex224 := position, tokenIndex
			{
				position225 := position
				{
					position226, tokenIndex226 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex = position226, tokenIndex226
					if buffer[position] != rune('C') {
						goto l224
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex = position228, tokenIndex228
					if buffer[position] != rune('R') {
						goto l224
					}
					position++
				}
			l228:
				{
					position230, tokenIndex230 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex = position230, tokenIndex230
					if buffer[position] != rune('E') {
						goto l224
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex = position232, tokenIndex232
					if buffer[position] != rune('A') {
						goto l224
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex = position234, tokenIndex234
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex = position236, tokenIndex236
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
					position238, tokenIndex238 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex = position238, tokenIndex238
					if buffer[position] != rune('S') {
						goto l224
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex = position240, tokenIndex240
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex = position242, tokenIndex242
					if buffer[position] != rune('A') {
						goto l224
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex = position244, tokenIndex244
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex = position246, tokenIndex246
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
					position248, tokenIndex248 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex = position248, tokenIndex248
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex = position250, tokenIndex250
					if buffer[position] != rune('Y') {
						goto l224
					}
					position++
				}
			l250:
				{
					position252, tokenIndex252 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex = position252, tokenIndex252
					if buffer[position] != rune('P') {
						goto l224
					}
					position++
				}
			l252:
				{
					position254, tokenIndex254 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex = position254, tokenIndex254
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
				add(ruleCreateStateStmt, position225)
			}
			return true
		l224:
			position, tokenIndex = position224, tokenIndex224
			return false
		},
		/* 15 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position256, tokenIndex256 := position, tokenIndex
			{
				position257 := position
				{
					position258, tokenIndex258 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex = position258, tokenIndex258
					if buffer[position] != rune('U') {
						goto l256
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex = position260, tokenIndex260
					if buffer[position] != rune('P') {
						goto l256
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex = position262, tokenIndex262
					if buffer[position] != rune('D') {
						goto l256
					}
					position++
				}
			l262:
				{
					position264, tokenIndex264 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex = position264, tokenIndex264
					if buffer[position] != rune('A') {
						goto l256
					}
					position++
				}
			l264:
				{
					position266, tokenIndex266 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex = position266, tokenIndex266
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex = position268, tokenIndex268
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
					position270, tokenIndex270 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex = position270, tokenIndex270
					if buffer[position] != rune('S') {
						goto l256
					}
					position++
				}
			l270:
				{
					position272, tokenIndex272 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l273
					}
					position++
					goto l272
				l273:
					position, tokenIndex = position272, tokenIndex272
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l272:
				{
					position274, tokenIndex274 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l275
					}
					position++
					goto l274
				l275:
					position, tokenIndex = position274, tokenIndex274
					if buffer[position] != rune('A') {
						goto l256
					}
					position++
				}
			l274:
				{
					position276, tokenIndex276 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l277
					}
					position++
					goto l276
				l277:
					position, tokenIndex = position276, tokenIndex276
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l276:
				{
					position278, tokenIndex278 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l279
					}
					position++
					goto l278
				l279:
					position, tokenIndex = position278, tokenIndex278
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
				add(ruleUpdateStateStmt, position257)
			}
			return true
		l256:
			position, tokenIndex = position256, tokenIndex256
			return false
		},
		/* 16 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action10)> */
		func() bool {
			position280, tokenIndex280 := position, tokenIndex
			{
				position281 := position
				{
					position282, tokenIndex282 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex = position282, tokenIndex282
					if buffer[position] != rune('U') {
						goto l280
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex = position284, tokenIndex284
					if buffer[position] != rune('P') {
						goto l280
					}
					position++
				}
			l284:
				{
					position286, tokenIndex286 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l287
					}
					position++
					goto l286
				l287:
					position, tokenIndex = position286, tokenIndex286
					if buffer[position] != rune('D') {
						goto l280
					}
					position++
				}
			l286:
				{
					position288, tokenIndex288 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l289
					}
					position++
					goto l288
				l289:
					position, tokenIndex = position288, tokenIndex288
					if buffer[position] != rune('A') {
						goto l280
					}
					position++
				}
			l288:
				{
					position290, tokenIndex290 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l291
					}
					position++
					goto l290
				l291:
					position, tokenIndex = position290, tokenIndex290
					if buffer[position] != rune('T') {
						goto l280
					}
					position++
				}
			l290:
				{
					position292, tokenIndex292 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l293
					}
					position++
					goto l292
				l293:
					position, tokenIndex = position292, tokenIndex292
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
					position294, tokenIndex294 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l295
					}
					position++
					goto l294
				l295:
					position, tokenIndex = position294, tokenIndex294
					if buffer[position] != rune('S') {
						goto l280
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex = position296, tokenIndex296
					if buffer[position] != rune('O') {
						goto l280
					}
					position++
				}
			l296:
				{
					position298, tokenIndex298 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex = position298, tokenIndex298
					if buffer[position] != rune('U') {
						goto l280
					}
					position++
				}
			l298:
				{
					position300, tokenIndex300 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l301
					}
					position++
					goto l300
				l301:
					position, tokenIndex = position300, tokenIndex300
					if buffer[position] != rune('R') {
						goto l280
					}
					position++
				}
			l300:
				{
					position302, tokenIndex302 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l303
					}
					position++
					goto l302
				l303:
					position, tokenIndex = position302, tokenIndex302
					if buffer[position] != rune('C') {
						goto l280
					}
					position++
				}
			l302:
				{
					position304, tokenIndex304 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex = position304, tokenIndex304
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
				add(ruleUpdateSourceStmt, position281)
			}
			return true
		l280:
			position, tokenIndex = position280, tokenIndex280
			return false
		},
		/* 17 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier UpdateSourceSinkSpecs Action11)> */
		func() bool {
			position306, tokenIndex306 := position, tokenIndex
			{
				position307 := position
				{
					position308, tokenIndex308 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex = position308, tokenIndex308
					if buffer[position] != rune('U') {
						goto l306
					}
					position++
				}
			l308:
				{
					position310, tokenIndex310 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l311
					}
					position++
					goto l310
				l311:
					position, tokenIndex = position310, tokenIndex310
					if buffer[position] != rune('P') {
						goto l306
					}
					position++
				}
			l310:
				{
					position312, tokenIndex312 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex = position312, tokenIndex312
					if buffer[position] != rune('D') {
						goto l306
					}
					position++
				}
			l312:
				{
					position314, tokenIndex314 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l315
					}
					position++
					goto l314
				l315:
					position, tokenIndex = position314, tokenIndex314
					if buffer[position] != rune('A') {
						goto l306
					}
					position++
				}
			l314:
				{
					position316, tokenIndex316 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l317
					}
					position++
					goto l316
				l317:
					position, tokenIndex = position316, tokenIndex316
					if buffer[position] != rune('T') {
						goto l306
					}
					position++
				}
			l316:
				{
					position318, tokenIndex318 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex = position318, tokenIndex318
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
					position320, tokenIndex320 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex = position320, tokenIndex320
					if buffer[position] != rune('S') {
						goto l306
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex = position322, tokenIndex322
					if buffer[position] != rune('I') {
						goto l306
					}
					position++
				}
			l322:
				{
					position324, tokenIndex324 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l325
					}
					position++
					goto l324
				l325:
					position, tokenIndex = position324, tokenIndex324
					if buffer[position] != rune('N') {
						goto l306
					}
					position++
				}
			l324:
				{
					position326, tokenIndex326 := position, tokenIndex
					if buffer[position] != rune('k') {
						goto l327
					}
					position++
					goto l326
				l327:
					position, tokenIndex = position326, tokenIndex326
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
				add(ruleUpdateSinkStmt, position307)
			}
			return true
		l306:
			position, tokenIndex = position306, tokenIndex306
			return false
		},
		/* 18 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action12)> */
		func() bool {
			position328, tokenIndex328 := position, tokenIndex
			{
				position329 := position
				{
					position330, tokenIndex330 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l331
					}
					position++
					goto l330
				l331:
					position, tokenIndex = position330, tokenIndex330
					if buffer[position] != rune('I') {
						goto l328
					}
					position++
				}
			l330:
				{
					position332, tokenIndex332 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l333
					}
					position++
					goto l332
				l333:
					position, tokenIndex = position332, tokenIndex332
					if buffer[position] != rune('N') {
						goto l328
					}
					position++
				}
			l332:
				{
					position334, tokenIndex334 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex = position334, tokenIndex334
					if buffer[position] != rune('S') {
						goto l328
					}
					position++
				}
			l334:
				{
					position336, tokenIndex336 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex = position336, tokenIndex336
					if buffer[position] != rune('E') {
						goto l328
					}
					position++
				}
			l336:
				{
					position338, tokenIndex338 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l339
					}
					position++
					goto l338
				l339:
					position, tokenIndex = position338, tokenIndex338
					if buffer[position] != rune('R') {
						goto l328
					}
					position++
				}
			l338:
				{
					position340, tokenIndex340 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l341
					}
					position++
					goto l340
				l341:
					position, tokenIndex = position340, tokenIndex340
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
					position342, tokenIndex342 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex = position342, tokenIndex342
					if buffer[position] != rune('I') {
						goto l328
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex = position344, tokenIndex344
					if buffer[position] != rune('N') {
						goto l328
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex = position346, tokenIndex346
					if buffer[position] != rune('T') {
						goto l328
					}
					position++
				}
			l346:
				{
					position348, tokenIndex348 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex = position348, tokenIndex348
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
				{
					position350, tokenIndex350 := position, tokenIndex
					if buffer[position] != rune('f') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex = position350, tokenIndex350
					if buffer[position] != rune('F') {
						goto l328
					}
					position++
				}
			l350:
				{
					position352, tokenIndex352 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex = position352, tokenIndex352
					if buffer[position] != rune('R') {
						goto l328
					}
					position++
				}
			l352:
				{
					position354, tokenIndex354 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex = position354, tokenIndex354
					if buffer[position] != rune('O') {
						goto l328
					}
					position++
				}
			l354:
				{
					position356, tokenIndex356 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l357
					}
					position++
					goto l356
				l357:
					position, tokenIndex = position356, tokenIndex356
					if buffer[position] != rune('M') {
						goto l328
					}
					position++
				}
			l356:
				if !_rules[rulesp]() {
					goto l328
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l328
				}
				if !_rules[ruleAction12]() {
					goto l328
				}
				add(ruleInsertIntoFromStmt, position329)
			}
			return true
		l328:
			position, tokenIndex = position328, tokenIndex328
			return false
		},
		/* 19 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position358, tokenIndex358 := position, tokenIndex
			{
				position359 := position
				{
					position360, tokenIndex360 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex = position360, tokenIndex360
					if buffer[position] != rune('P') {
						goto l358
					}
					position++
				}
			l360:
				{
					position362, tokenIndex362 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l363
					}
					position++
					goto l362
				l363:
					position, tokenIndex = position362, tokenIndex362
					if buffer[position] != rune('A') {
						goto l358
					}
					position++
				}
			l362:
				{
					position364, tokenIndex364 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l365
					}
					position++
					goto l364
				l365:
					position, tokenIndex = position364, tokenIndex364
					if buffer[position] != rune('U') {
						goto l358
					}
					position++
				}
			l364:
				{
					position366, tokenIndex366 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l367
					}
					position++
					goto l366
				l367:
					position, tokenIndex = position366, tokenIndex366
					if buffer[position] != rune('S') {
						goto l358
					}
					position++
				}
			l366:
				{
					position368, tokenIndex368 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex = position368, tokenIndex368
					if buffer[position] != rune('E') {
						goto l358
					}
					position++
				}
			l368:
				if !_rules[rulesp]() {
					goto l358
				}
				{
					position370, tokenIndex370 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l371
					}
					position++
					goto l370
				l371:
					position, tokenIndex = position370, tokenIndex370
					if buffer[position] != rune('S') {
						goto l358
					}
					position++
				}
			l370:
				{
					position372, tokenIndex372 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l373
					}
					position++
					goto l372
				l373:
					position, tokenIndex = position372, tokenIndex372
					if buffer[position] != rune('O') {
						goto l358
					}
					position++
				}
			l372:
				{
					position374, tokenIndex374 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex = position374, tokenIndex374
					if buffer[position] != rune('U') {
						goto l358
					}
					position++
				}
			l374:
				{
					position376, tokenIndex376 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l377
					}
					position++
					goto l376
				l377:
					position, tokenIndex = position376, tokenIndex376
					if buffer[position] != rune('R') {
						goto l358
					}
					position++
				}
			l376:
				{
					position378, tokenIndex378 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l379
					}
					position++
					goto l378
				l379:
					position, tokenIndex = position378, tokenIndex378
					if buffer[position] != rune('C') {
						goto l358
					}
					position++
				}
			l378:
				{
					position380, tokenIndex380 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex = position380, tokenIndex380
					if buffer[position] != rune('E') {
						goto l358
					}
					position++
				}
			l380:
				if !_rules[rulesp]() {
					goto l358
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l358
				}
				if !_rules[ruleAction13]() {
					goto l358
				}
				add(rulePauseSourceStmt, position359)
			}
			return true
		l358:
			position, tokenIndex = position358, tokenIndex358
			return false
		},
		/* 20 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position382, tokenIndex382 := position, tokenIndex
			{
				position383 := position
				{
					position384, tokenIndex384 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l385
					}
					position++
					goto l384
				l385:
					position, tokenIndex = position384, tokenIndex384
					if buffer[position] != rune('R') {
						goto l382
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex = position386, tokenIndex386
					if buffer[position] != rune('E') {
						goto l382
					}
					position++
				}
			l386:
				{
					position388, tokenIndex388 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l389
					}
					position++
					goto l388
				l389:
					position, tokenIndex = position388, tokenIndex388
					if buffer[position] != rune('S') {
						goto l382
					}
					position++
				}
			l388:
				{
					position390, tokenIndex390 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l391
					}
					position++
					goto l390
				l391:
					position, tokenIndex = position390, tokenIndex390
					if buffer[position] != rune('U') {
						goto l382
					}
					position++
				}
			l390:
				{
					position392, tokenIndex392 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex = position392, tokenIndex392
					if buffer[position] != rune('M') {
						goto l382
					}
					position++
				}
			l392:
				{
					position394, tokenIndex394 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l395
					}
					position++
					goto l394
				l395:
					position, tokenIndex = position394, tokenIndex394
					if buffer[position] != rune('E') {
						goto l382
					}
					position++
				}
			l394:
				if !_rules[rulesp]() {
					goto l382
				}
				{
					position396, tokenIndex396 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex = position396, tokenIndex396
					if buffer[position] != rune('S') {
						goto l382
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex = position398, tokenIndex398
					if buffer[position] != rune('O') {
						goto l382
					}
					position++
				}
			l398:
				{
					position400, tokenIndex400 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l401
					}
					position++
					goto l400
				l401:
					position, tokenIndex = position400, tokenIndex400
					if buffer[position] != rune('U') {
						goto l382
					}
					position++
				}
			l400:
				{
					position402, tokenIndex402 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l403
					}
					position++
					goto l402
				l403:
					position, tokenIndex = position402, tokenIndex402
					if buffer[position] != rune('R') {
						goto l382
					}
					position++
				}
			l402:
				{
					position404, tokenIndex404 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l405
					}
					position++
					goto l404
				l405:
					position, tokenIndex = position404, tokenIndex404
					if buffer[position] != rune('C') {
						goto l382
					}
					position++
				}
			l404:
				{
					position406, tokenIndex406 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l407
					}
					position++
					goto l406
				l407:
					position, tokenIndex = position406, tokenIndex406
					if buffer[position] != rune('E') {
						goto l382
					}
					position++
				}
			l406:
				if !_rules[rulesp]() {
					goto l382
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l382
				}
				if !_rules[ruleAction14]() {
					goto l382
				}
				add(ruleResumeSourceStmt, position383)
			}
			return true
		l382:
			position, tokenIndex = position382, tokenIndex382
			return false
		},
		/* 21 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position408, tokenIndex408 := position, tokenIndex
			{
				position409 := position
				{
					position410, tokenIndex410 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l411
					}
					position++
					goto l410
				l411:
					position, tokenIndex = position410, tokenIndex410
					if buffer[position] != rune('R') {
						goto l408
					}
					position++
				}
			l410:
				{
					position412, tokenIndex412 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l413
					}
					position++
					goto l412
				l413:
					position, tokenIndex = position412, tokenIndex412
					if buffer[position] != rune('E') {
						goto l408
					}
					position++
				}
			l412:
				{
					position414, tokenIndex414 := position, tokenIndex
					if buffer[position] != rune('w') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex = position414, tokenIndex414
					if buffer[position] != rune('W') {
						goto l408
					}
					position++
				}
			l414:
				{
					position416, tokenIndex416 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l417
					}
					position++
					goto l416
				l417:
					position, tokenIndex = position416, tokenIndex416
					if buffer[position] != rune('I') {
						goto l408
					}
					position++
				}
			l416:
				{
					position418, tokenIndex418 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex = position418, tokenIndex418
					if buffer[position] != rune('N') {
						goto l408
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex = position420, tokenIndex420
					if buffer[position] != rune('D') {
						goto l408
					}
					position++
				}
			l420:
				if !_rules[rulesp]() {
					goto l408
				}
				{
					position422, tokenIndex422 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l423
					}
					position++
					goto l422
				l423:
					position, tokenIndex = position422, tokenIndex422
					if buffer[position] != rune('S') {
						goto l408
					}
					position++
				}
			l422:
				{
					position424, tokenIndex424 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l425
					}
					position++
					goto l424
				l425:
					position, tokenIndex = position424, tokenIndex424
					if buffer[position] != rune('O') {
						goto l408
					}
					position++
				}
			l424:
				{
					position426, tokenIndex426 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l427
					}
					position++
					goto l426
				l427:
					position, tokenIndex = position426, tokenIndex426
					if buffer[position] != rune('U') {
						goto l408
					}
					position++
				}
			l426:
				{
					position428, tokenIndex428 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l429
					}
					position++
					goto l428
				l429:
					position, tokenIndex = position428, tokenIndex428
					if buffer[position] != rune('R') {
						goto l408
					}
					position++
				}
			l428:
				{
					position430, tokenIndex430 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l431
					}
					position++
					goto l430
				l431:
					position, tokenIndex = position430, tokenIndex430
					if buffer[position] != rune('C') {
						goto l408
					}
					position++
				}
			l430:
				{
					position432, tokenIndex432 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex = position432, tokenIndex432
					if buffer[position] != rune('E') {
						goto l408
					}
					position++
				}
			l432:
				if !_rules[rulesp]() {
					goto l408
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l408
				}
				if !_rules[ruleAction15]() {
					goto l408
				}
				add(ruleRewindSourceStmt, position409)
			}
			return true
		l408:
			position, tokenIndex = position408, tokenIndex408
			return false
		},
		/* 22 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position434, tokenIndex434 := position, tokenIndex
			{
				position435 := position
				{
					position436, tokenIndex436 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l437
					}
					position++
					goto l436
				l437:
					position, tokenIndex = position436, tokenIndex436
					if buffer[position] != rune('D') {
						goto l434
					}
					position++
				}
			l436:
				{
					position438, tokenIndex438 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex = position438, tokenIndex438
					if buffer[position] != rune('R') {
						goto l434
					}
					position++
				}
			l438:
				{
					position440, tokenIndex440 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l441
					}
					position++
					goto l440
				l441:
					position, tokenIndex = position440, tokenIndex440
					if buffer[position] != rune('O') {
						goto l434
					}
					position++
				}
			l440:
				{
					position442, tokenIndex442 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l443
					}
					position++
					goto l442
				l443:
					position, tokenIndex = position442, tokenIndex442
					if buffer[position] != rune('P') {
						goto l434
					}
					position++
				}
			l442:
				if !_rules[rulesp]() {
					goto l434
				}
				{
					position444, tokenIndex444 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l445
					}
					position++
					goto l444
				l445:
					position, tokenIndex = position444, tokenIndex444
					if buffer[position] != rune('S') {
						goto l434
					}
					position++
				}
			l444:
				{
					position446, tokenIndex446 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex = position446, tokenIndex446
					if buffer[position] != rune('O') {
						goto l434
					}
					position++
				}
			l446:
				{
					position448, tokenIndex448 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex = position448, tokenIndex448
					if buffer[position] != rune('U') {
						goto l434
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex = position450, tokenIndex450
					if buffer[position] != rune('R') {
						goto l434
					}
					position++
				}
			l450:
				{
					position452, tokenIndex452 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l453
					}
					position++
					goto l452
				l453:
					position, tokenIndex = position452, tokenIndex452
					if buffer[position] != rune('C') {
						goto l434
					}
					position++
				}
			l452:
				{
					position454, tokenIndex454 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l455
					}
					position++
					goto l454
				l455:
					position, tokenIndex = position454, tokenIndex454
					if buffer[position] != rune('E') {
						goto l434
					}
					position++
				}
			l454:
				if !_rules[rulesp]() {
					goto l434
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l434
				}
				if !_rules[ruleAction16]() {
					goto l434
				}
				add(ruleDropSourceStmt, position435)
			}
			return true
		l434:
			position, tokenIndex = position434, tokenIndex434
			return false
		},
		/* 23 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action17)> */
		func() bool {
			position456, tokenIndex456 := position, tokenIndex
			{
				position457 := position
				{
					position458, tokenIndex458 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l459
					}
					position++
					goto l458
				l459:
					position, tokenIndex = position458, tokenIndex458
					if buffer[position] != rune('D') {
						goto l456
					}
					position++
				}
			l458:
				{
					position460, tokenIndex460 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l461
					}
					position++
					goto l460
				l461:
					position, tokenIndex = position460, tokenIndex460
					if buffer[position] != rune('R') {
						goto l456
					}
					position++
				}
			l460:
				{
					position462, tokenIndex462 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l463
					}
					position++
					goto l462
				l463:
					position, tokenIndex = position462, tokenIndex462
					if buffer[position] != rune('O') {
						goto l456
					}
					position++
				}
			l462:
				{
					position464, tokenIndex464 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l465
					}
					position++
					goto l464
				l465:
					position, tokenIndex = position464, tokenIndex464
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
					position466, tokenIndex466 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l467
					}
					position++
					goto l466
				l467:
					position, tokenIndex = position466, tokenIndex466
					if buffer[position] != rune('S') {
						goto l456
					}
					position++
				}
			l466:
				{
					position468, tokenIndex468 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l469
					}
					position++
					goto l468
				l469:
					position, tokenIndex = position468, tokenIndex468
					if buffer[position] != rune('T') {
						goto l456
					}
					position++
				}
			l468:
				{
					position470, tokenIndex470 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l471
					}
					position++
					goto l470
				l471:
					position, tokenIndex = position470, tokenIndex470
					if buffer[position] != rune('R') {
						goto l456
					}
					position++
				}
			l470:
				{
					position472, tokenIndex472 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l473
					}
					position++
					goto l472
				l473:
					position, tokenIndex = position472, tokenIndex472
					if buffer[position] != rune('E') {
						goto l456
					}
					position++
				}
			l472:
				{
					position474, tokenIndex474 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l475
					}
					position++
					goto l474
				l475:
					position, tokenIndex = position474, tokenIndex474
					if buffer[position] != rune('A') {
						goto l456
					}
					position++
				}
			l474:
				{
					position476, tokenIndex476 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l477
					}
					position++
					goto l476
				l477:
					position, tokenIndex = position476, tokenIndex476
					if buffer[position] != rune('M') {
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
				add(ruleDropStreamStmt, position457)
			}
			return true
		l456:
			position, tokenIndex = position456, tokenIndex456
			return false
		},
		/* 24 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action18)> */
		func() bool {
			position478, tokenIndex478 := position, tokenIndex
			{
				position479 := position
				{
					position480, tokenIndex480 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l481
					}
					position++
					goto l480
				l481:
					position, tokenIndex = position480, tokenIndex480
					if buffer[position] != rune('D') {
						goto l478
					}
					position++
				}
			l480:
				{
					position482, tokenIndex482 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l483
					}
					position++
					goto l482
				l483:
					position, tokenIndex = position482, tokenIndex482
					if buffer[position] != rune('R') {
						goto l478
					}
					position++
				}
			l482:
				{
					position484, tokenIndex484 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l485
					}
					position++
					goto l484
				l485:
					position, tokenIndex = position484, tokenIndex484
					if buffer[position] != rune('O') {
						goto l478
					}
					position++
				}
			l484:
				{
					position486, tokenIndex486 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l487
					}
					position++
					goto l486
				l487:
					position, tokenIndex = position486, tokenIndex486
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
					position488, tokenIndex488 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l489
					}
					position++
					goto l488
				l489:
					position, tokenIndex = position488, tokenIndex488
					if buffer[position] != rune('S') {
						goto l478
					}
					position++
				}
			l488:
				{
					position490, tokenIndex490 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l491
					}
					position++
					goto l490
				l491:
					position, tokenIndex = position490, tokenIndex490
					if buffer[position] != rune('I') {
						goto l478
					}
					position++
				}
			l490:
				{
					position492, tokenIndex492 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l493
					}
					position++
					goto l492
				l493:
					position, tokenIndex = position492, tokenIndex492
					if buffer[position] != rune('N') {
						goto l478
					}
					position++
				}
			l492:
				{
					position494, tokenIndex494 := position, tokenIndex
					if buffer[position] != rune('k') {
						goto l495
					}
					position++
					goto l494
				l495:
					position, tokenIndex = position494, tokenIndex494
					if buffer[position] != rune('K') {
						goto l478
					}
					position++
				}
			l494:
				if !_rules[rulesp]() {
					goto l478
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l478
				}
				if !_rules[ruleAction18]() {
					goto l478
				}
				add(ruleDropSinkStmt, position479)
			}
			return true
		l478:
			position, tokenIndex = position478, tokenIndex478
			return false
		},
		/* 25 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action19)> */
		func() bool {
			position496, tokenIndex496 := position, tokenIndex
			{
				position497 := position
				{
					position498, tokenIndex498 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l499
					}
					position++
					goto l498
				l499:
					position, tokenIndex = position498, tokenIndex498
					if buffer[position] != rune('D') {
						goto l496
					}
					position++
				}
			l498:
				{
					position500, tokenIndex500 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l501
					}
					position++
					goto l500
				l501:
					position, tokenIndex = position500, tokenIndex500
					if buffer[position] != rune('R') {
						goto l496
					}
					position++
				}
			l500:
				{
					position502, tokenIndex502 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l503
					}
					position++
					goto l502
				l503:
					position, tokenIndex = position502, tokenIndex502
					if buffer[position] != rune('O') {
						goto l496
					}
					position++
				}
			l502:
				{
					position504, tokenIndex504 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l505
					}
					position++
					goto l504
				l505:
					position, tokenIndex = position504, tokenIndex504
					if buffer[position] != rune('P') {
						goto l496
					}
					position++
				}
			l504:
				if !_rules[rulesp]() {
					goto l496
				}
				{
					position506, tokenIndex506 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l507
					}
					position++
					goto l506
				l507:
					position, tokenIndex = position506, tokenIndex506
					if buffer[position] != rune('S') {
						goto l496
					}
					position++
				}
			l506:
				{
					position508, tokenIndex508 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l509
					}
					position++
					goto l508
				l509:
					position, tokenIndex = position508, tokenIndex508
					if buffer[position] != rune('T') {
						goto l496
					}
					position++
				}
			l508:
				{
					position510, tokenIndex510 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l511
					}
					position++
					goto l510
				l511:
					position, tokenIndex = position510, tokenIndex510
					if buffer[position] != rune('A') {
						goto l496
					}
					position++
				}
			l510:
				{
					position512, tokenIndex512 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l513
					}
					position++
					goto l512
				l513:
					position, tokenIndex = position512, tokenIndex512
					if buffer[position] != rune('T') {
						goto l496
					}
					position++
				}
			l512:
				{
					position514, tokenIndex514 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l515
					}
					position++
					goto l514
				l515:
					position, tokenIndex = position514, tokenIndex514
					if buffer[position] != rune('E') {
						goto l496
					}
					position++
				}
			l514:
				if !_rules[rulesp]() {
					goto l496
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l496
				}
				if !_rules[ruleAction19]() {
					goto l496
				}
				add(ruleDropStateStmt, position497)
			}
			return true
		l496:
			position, tokenIndex = position496, tokenIndex496
			return false
		},
		/* 26 LoadStateStmt <- <(('l' / 'L') ('o' / 'O') ('a' / 'A') ('d' / 'D') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType StateTagOpt SetOptSpecs Action20)> */
		func() bool {
			position516, tokenIndex516 := position, tokenIndex
			{
				position517 := position
				{
					position518, tokenIndex518 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l519
					}
					position++
					goto l518
				l519:
					position, tokenIndex = position518, tokenIndex518
					if buffer[position] != rune('L') {
						goto l516
					}
					position++
				}
			l518:
				{
					position520, tokenIndex520 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l521
					}
					position++
					goto l520
				l521:
					position, tokenIndex = position520, tokenIndex520
					if buffer[position] != rune('O') {
						goto l516
					}
					position++
				}
			l520:
				{
					position522, tokenIndex522 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l523
					}
					position++
					goto l522
				l523:
					position, tokenIndex = position522, tokenIndex522
					if buffer[position] != rune('A') {
						goto l516
					}
					position++
				}
			l522:
				{
					position524, tokenIndex524 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l525
					}
					position++
					goto l524
				l525:
					position, tokenIndex = position524, tokenIndex524
					if buffer[position] != rune('D') {
						goto l516
					}
					position++
				}
			l524:
				if !_rules[rulesp]() {
					goto l516
				}
				{
					position526, tokenIndex526 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l527
					}
					position++
					goto l526
				l527:
					position, tokenIndex = position526, tokenIndex526
					if buffer[position] != rune('S') {
						goto l516
					}
					position++
				}
			l526:
				{
					position528, tokenIndex528 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l529
					}
					position++
					goto l528
				l529:
					position, tokenIndex = position528, tokenIndex528
					if buffer[position] != rune('T') {
						goto l516
					}
					position++
				}
			l528:
				{
					position530, tokenIndex530 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l531
					}
					position++
					goto l530
				l531:
					position, tokenIndex = position530, tokenIndex530
					if buffer[position] != rune('A') {
						goto l516
					}
					position++
				}
			l530:
				{
					position532, tokenIndex532 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l533
					}
					position++
					goto l532
				l533:
					position, tokenIndex = position532, tokenIndex532
					if buffer[position] != rune('T') {
						goto l516
					}
					position++
				}
			l532:
				{
					position534, tokenIndex534 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l535
					}
					position++
					goto l534
				l535:
					position, tokenIndex = position534, tokenIndex534
					if buffer[position] != rune('E') {
						goto l516
					}
					position++
				}
			l534:
				if !_rules[rulesp]() {
					goto l516
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l516
				}
				if !_rules[rulesp]() {
					goto l516
				}
				{
					position536, tokenIndex536 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l537
					}
					position++
					goto l536
				l537:
					position, tokenIndex = position536, tokenIndex536
					if buffer[position] != rune('T') {
						goto l516
					}
					position++
				}
			l536:
				{
					position538, tokenIndex538 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l539
					}
					position++
					goto l538
				l539:
					position, tokenIndex = position538, tokenIndex538
					if buffer[position] != rune('Y') {
						goto l516
					}
					position++
				}
			l538:
				{
					position540, tokenIndex540 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l541
					}
					position++
					goto l540
				l541:
					position, tokenIndex = position540, tokenIndex540
					if buffer[position] != rune('P') {
						goto l516
					}
					position++
				}
			l540:
				{
					position542, tokenIndex542 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l543
					}
					position++
					goto l542
				l543:
					position, tokenIndex = position542, tokenIndex542
					if buffer[position] != rune('E') {
						goto l516
					}
					position++
				}
			l542:
				if !_rules[rulesp]() {
					goto l516
				}
				if !_rules[ruleSourceSinkType]() {
					goto l516
				}
				if !_rules[ruleStateTagOpt]() {
					goto l516
				}
				if !_rules[ruleSetOptSpecs]() {
					goto l516
				}
				if !_rules[ruleAction20]() {
					goto l516
				}
				add(ruleLoadStateStmt, position517)
			}
			return true
		l516:
			position, tokenIndex = position516, tokenIndex516
			return false
		},
		/* 27 LoadStateOrCreateStmt <- <(LoadStateStmt sp (('o' / 'O') ('r' / 'R')) sp (('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp (('i' / 'I') ('f' / 'F')) sp (('n' / 'N') ('o' / 'O') ('t' / 'T')) sp ((('s' / 'S') ('a' / 'A') ('v' / 'V') ('e' / 'E') ('d' / 'D')) / (('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S'))) SourceSinkSpecs Action21)> */
		func() bool {
			position544, tokenIndex544 := position, tokenIndex
			{
				position545 := position
				if !_rules[ruleLoadStateStmt]() {
					goto l544
				}
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position546, tokenIndex546 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l547
					}
					position++
					goto l546
				l547:
					position, tokenIndex = position546, tokenIndex546
					if buffer[position] != rune('O') {
						goto l544
					}
					position++
				}
			l546:
				{
					position548, tokenIndex548 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l549
					}
					position++
					goto l548
				l549:
					position, tokenIndex = position548, tokenIndex548
					if buffer[position] != rune('R') {
						goto l544
					}
					position++
				}
			l548:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position550, tokenIndex550 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l551
					}
					position++
					goto l550
				l551:
					position, tokenIndex = position550, tokenIndex550
					if buffer[position] != rune('C') {
						goto l544
					}
					position++
				}
			l550:
				{
					position552, tokenIndex552 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l553
					}
					position++
					goto l552
				l553:
					position, tokenIndex = position552, tokenIndex552
					if buffer[position] != rune('R') {
						goto l544
					}
					position++
				}
			l552:
				{
					position554, tokenIndex554 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l555
					}
					position++
					goto l554
				l555:
					position, tokenIndex = position554, tokenIndex554
					if buffer[position] != rune('E') {
						goto l544
					}
					position++
				}
			l554:
				{
					position556, tokenIndex556 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l557
					}
					position++
					goto l556
				l557:
					position, tokenIndex = position556, tokenIndex556
					if buffer[position] != rune('A') {
						goto l544
					}
					position++
				}
			l556:
				{
					position558, tokenIndex558 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l559
					}
					position++
					goto l558
				l559:
					position, tokenIndex = position558, tokenIndex558
					if buffer[position] != rune('T') {
						goto l544
					}
					position++
				}
			l558:
				{
					position560, tokenIndex560 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l561
					}
					position++
					goto l560
				l561:
					position, tokenIndex = position560, tokenIndex560
					if buffer[position] != rune('E') {
						goto l544
					}
					position++
				}
			l560:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position562, tokenIndex562 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l563
					}
					position++
					goto l562
				l563:
					position, tokenIndex = position562, tokenIndex562
					if buffer[position] != rune('I') {
						goto l544
					}
					position++
				}
			l562:
				{
					position564, tokenIndex564 := position, tokenIndex
					if buffer[position] != rune('f') {
						goto l565
					}
					position++
					goto l564
				l565:
					position, tokenIndex = position564, tokenIndex564
					if buffer[position] != rune('F') {
						goto l544
					}
					position++
				}
			l564:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position566, tokenIndex566 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l567
					}
					position++
					goto l566
				l567:
					position, tokenIndex = position566, tokenIndex566
					if buffer[position] != rune('N') {
						goto l544
					}
					position++
				}
			l566:
				{
					position568, tokenIndex568 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l569
					}
					position++
					goto l568
				l569:
					position, tokenIndex = position568, tokenIndex568
					if buffer[position] != rune('O') {
						goto l544
					}
					position++
				}
			l568:
				{
					position570, tokenIndex570 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l571
					}
					position++
					goto l570
				l571:
					position, tokenIndex = position570, tokenIndex570
					if buffer[position] != rune('T') {
						goto l544
					}
					position++
				}
			l570:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position572, tokenIndex572 := position, tokenIndex
					{
						position574, tokenIndex574 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l575
						}
						position++
						goto l574
					l575:
						position, tokenIndex = position574, tokenIndex574
						if buffer[position] != rune('S') {
							goto l573
						}
						position++
					}
				l574:
					{
						position576, tokenIndex576 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l577
						}
						position++
						goto l576
					l577:
						position, tokenIndex = position576, tokenIndex576
						if buffer[position] != rune('A') {
							goto l573
						}
						position++
					}
				l576:
					{
						position578, tokenIndex578 := position, tokenIndex
						if buffer[position] != rune('v') {
							goto l579
						}
						position++
						goto l578
					l579:
						position, tokenIndex = position578, tokenIndex578
						if buffer[position] != rune('V') {
							goto l573
						}
						position++
					}
				l578:
					{
						position580, tokenIndex580 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l581
						}
						position++
						goto l580
					l581:
						position, tokenIndex = position580, tokenIndex580
						if buffer[position] != rune('E') {
							goto l573
						}
						position++
					}
				l580:
					{
						position582, tokenIndex582 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l583
						}
						position++
						goto l582
					l583:
						position, tokenIndex = position582, tokenIndex582
						if buffer[position] != rune('D') {
							goto l573
						}
						position++
					}
				l582:
					goto l572
				l573:
					position, tokenIndex = position572, tokenIndex572
					{
						position584, tokenIndex584 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l585
						}
						position++
						goto l584
					l585:
						position, tokenIndex = position584, tokenIndex584
						if buffer[position] != rune('E') {
							goto l544
						}
						position++
					}
				l584:
					{
						position586, tokenIndex586 := position, tokenIndex
						if buffer[position] != rune('x') {
							goto l587
						}
						position++
						goto l586
					l587:
						position, tokenIndex = position586, tokenIndex586
						if buffer[position] != rune('X') {
							goto l544
						}
						position++
					}
				l586:
					{
						position588, tokenIndex588 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l589
						}
						position++
						goto l588
					l589:
						position, tokenIndex = position588, tokenIndex588
						if buffer[position] != rune('I') {
							goto l544
						}
						position++
					}
				l588:
					{
						position590, tokenIndex590 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex = position590, tokenIndex590
						if buffer[position] != rune('S') {
							goto l544
						}
						position++
					}
				l590:
					{
						position592, tokenIndex592 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex = position592, tokenIndex592
						if buffer[position] != rune('T') {
							goto l544
						}
						position++
					}
				l592:
					{
						position594, tokenIndex594 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex = position594, tokenIndex594
						if buffer[position] != rune('S') {
							goto l544
						}
						position++
					}
				l594:
				}
			l572:
				if !_rules[ruleSourceSinkSpecs]() {
					goto l544
				}
				if !_rules[ruleAction21]() {
					goto l544
				}
				add(ruleLoadStateOrCreateStmt, position545)
			}
			return true
		l544:
			position, tokenIndex = position544, tokenIndex544
			return false
		},
		/* 28 SaveStateStmt <- <(('s' / 'S') ('a' / 'A') ('v' / 'V') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier StateTagOpt Action22)> */
		func() bool {
			position596, tokenIndex596 := position, tokenIndex
			{
				position597 := position
				{
					position598, tokenIndex598 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l599
					}
					position++
					goto l598
				l599:
					position, tokenIndex = position598, tokenIndex598
					if buffer[position] != rune('S') {
						goto l596
					}
					position++
				}
			l598:
				{
					position600, tokenIndex600 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l601
					}
					position++
					goto l600
				l601:
					position, tokenIndex = position600, tokenIndex600
					if buffer[position] != rune('A') {
						goto l596
					}
					position++
				}
			l600:
				{
					position602, tokenIndex602 := position, tokenIndex
					if buffer[position] != rune('v') {
						goto l603
					}
					position++
					goto l602
				l603:
					position, tokenIndex = position602, tokenIndex602
					if buffer[position] != rune('V') {
						goto l596
					}
					position++
				}
			l602:
				{
					position604, tokenIndex604 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l605
					}
					position++
					goto l604
				l605:
					position, tokenIndex = position604, tokenIndex604
					if buffer[position] != rune('E') {
						goto l596
					}
					position++
				}
			l604:
				if !_rules[rulesp]() {
					goto l596
				}
				{
					position606, tokenIndex606 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l607
					}
					position++
					goto l606
				l607:
					position, tokenIndex = position606, tokenIndex606
					if buffer[position] != rune('S') {
						goto l596
					}
					position++
				}
			l606:
				{
					position608, tokenIndex608 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l609
					}
					position++
					goto l608
				l609:
					position, tokenIndex = position608, tokenIndex608
					if buffer[position] != rune('T') {
						goto l596
					}
					position++
				}
			l608:
				{
					position610, tokenIndex610 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l611
					}
					position++
					goto l610
				l611:
					position, tokenIndex = position610, tokenIndex610
					if buffer[position] != rune('A') {
						goto l596
					}
					position++
				}
			l610:
				{
					position612, tokenIndex612 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l613
					}
					position++
					goto l612
				l613:
					position, tokenIndex = position612, tokenIndex612
					if buffer[position] != rune('T') {
						goto l596
					}
					position++
				}
			l612:
				{
					position614, tokenIndex614 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l615
					}
					position++
					goto l614
				l615:
					position, tokenIndex = position614, tokenIndex614
					if buffer[position] != rune('E') {
						goto l596
					}
					position++
				}
			l614:
				if !_rules[rulesp]() {
					goto l596
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l596
				}
				if !_rules[ruleStateTagOpt]() {
					goto l596
				}
				if !_rules[ruleAction22]() {
					goto l596
				}
				add(ruleSaveStateStmt, position597)
			}
			return true
		l596:
			position, tokenIndex = position596, tokenIndex596
			return false
		},
		/* 29 EvalStmt <- <(('e' / 'E') ('v' / 'V') ('a' / 'A') ('l' / 'L') sp Expression <(sp (('o' / 'O') ('n' / 'N')) sp MapExpr)?> Action23)> */
		func() bool {
			position616, tokenIndex616 := position, tokenIndex
			{
				position617 := position
				{
					position618, tokenIndex618 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l619
					}
					position++
					goto l618
				l619:
					position, tokenIndex = position618, tokenIndex618
					if buffer[position] != rune('E') {
						goto l616
					}
					position++
				}
			l618:
				{
					position620, tokenIndex620 := position, tokenIndex
					if buffer[position] != rune('v') {
						goto l621
					}
					position++
					goto l620
				l621:
					position, tokenIndex = position620, tokenIndex620
					if buffer[position] != rune('V') {
						goto l616
					}
					position++
				}
			l620:
				{
					position622, tokenIndex622 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l623
					}
					position++
					goto l622
				l623:
					position, tokenIndex = position622, tokenIndex622
					if buffer[position] != rune('A') {
						goto l616
					}
					position++
				}
			l622:
				{
					position624, tokenIndex624 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l625
					}
					position++
					goto l624
				l625:
					position, tokenIndex = position624, tokenIndex624
					if buffer[position] != rune('L') {
						goto l616
					}
					position++
				}
			l624:
				if !_rules[rulesp]() {
					goto l616
				}
				if !_rules[ruleExpression]() {
					goto l616
				}
				{
					position626 := position
					{
						position627, tokenIndex627 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l627
						}
						{
							position629, tokenIndex629 := position, tokenIndex
							if buffer[position] != rune('o') {
								goto l630
							}
							position++
							goto l629
						l630:
							position, tokenIndex = position629, tokenIndex629
							if buffer[position] != rune('O') {
								goto l627
							}
							position++
						}
					l629:
						{
							position631, tokenIndex631 := position, tokenIndex
							if buffer[position] != rune('n') {
								goto l632
							}
							position++
							goto l631
						l632:
							position, tokenIndex = position631, tokenIndex631
							if buffer[position] != rune('N') {
								goto l627
							}
							position++
						}
					l631:
						if !_rules[rulesp]() {
							goto l627
						}
						if !_rules[ruleMapExpr]() {
							goto l627
						}
						goto l628
					l627:
						position, tokenIndex = position627, tokenIndex627
					}
				l628:
					add(rulePegText, position626)
				}
				if !_rules[ruleAction23]() {
					goto l616
				}
				add(ruleEvalStmt, position617)
			}
			return true
		l616:
			position, tokenIndex = position616, tokenIndex616
			return false
		},
		/* 30 Emitter <- <(sp (ISTREAM / DSTREAM / RSTREAM) EmitterOptions Action24)> */
		func() bool {
			position633, tokenIndex633 := position, tokenIndex
			{
				position634 := position
				if !_rules[rulesp]() {
					goto l633
				}
				{
					position635, tokenIndex635 := position, tokenIndex
					if !_rules[ruleISTREAM]() {
						goto l636
					}
					goto l635
				l636:
					position, tokenIndex = position635, tokenIndex635
					if !_rules[ruleDSTREAM]() {
						goto l637
					}
					goto l635
				l637:
					position, tokenIndex = position635, tokenIndex635
					if !_rules[ruleRSTREAM]() {
						goto l633
					}
				}
			l635:
				if !_rules[ruleEmitterOptions]() {
					goto l633
				}
				if !_rules[ruleAction24]() {
					goto l633
				}
				add(ruleEmitter, position634)
			}
			return true
		l633:
			position, tokenIndex = position633, tokenIndex633
			return false
		},
		/* 31 EmitterOptions <- <(<(spOpt '[' spOpt EmitterOptionCombinations spOpt ']')?> Action25)> */
		func() bool {
			position638, tokenIndex638 := position, tokenIndex
			{
				position639 := position
				{
					position640 := position
					{
						position641, tokenIndex641 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l641
						}
						if buffer[position] != rune('[') {
							goto l641
						}
						position++
						if !_rules[rulespOpt]() {
							goto l641
						}
						if !_rules[ruleEmitterOptionCombinations]() {
							goto l641
						}
						if !_rules[rulespOpt]() {
							goto l641
						}
						if buffer[position] != rune(']') {
							goto l641
						}
						position++
						goto l642
					l641:
						position, tokenIndex = position641, tokenIndex641
					}
				l642:
					add(rulePegText, position640)
				}
				if !_rules[ruleAction25]() {
					goto l638
				}
				add(ruleEmitterOptions, position639)
			}
			return true
		l638:
			position, tokenIndex = position638, tokenIndex638
			return false
		},
		/* 32 EmitterOptionCombinations <- <(EmitterLimit / (EmitterSample sp EmitterLimit) / EmitterSample)> */
		func() bool {
			position643, tokenIndex643 := position, tokenIndex
			{
				position644 := position
				{
					position645, tokenIndex645 := position, tokenIndex
					if !_rules[ruleEmitterLimit]() {
						goto l646
					}
					goto l645
				l646:
					position, tokenIndex = position645, tokenIndex645
					if !_rules[ruleEmitterSample]() {
						goto l647
					}
					if !_rules[rulesp]() {
						goto l647
					}
					if !_rules[ruleEmitterLimit]() {
						goto l647
					}
					goto l645
				l647:
					position, tokenIndex = position645, tokenIndex645
					if !_rules[ruleEmitterSample]() {
						goto l643
					}
				}
			l645:
				add(ruleEmitterOptionCombinations, position644)
			}
			return true
		l643:
			position, tokenIndex = position643, tokenIndex643
			return false
		},
		/* 33 EmitterLimit <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') sp NumericLiteral Action26)> */
		func() bool {
			position648, tokenIndex648 := position, tokenIndex
			{
				position649 := position
				{
					position650, tokenIndex650 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l651
					}
					position++
					goto l650
				l651:
					position, tokenIndex = position650, tokenIndex650
					if buffer[position] != rune('L') {
						goto l648
					}
					position++
				}
			l650:
				{
					position652, tokenIndex652 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l653
					}
					position++
					goto l652
				l653:
					position, tokenIndex = position652, tokenIndex652
					if buffer[position] != rune('I') {
						goto l648
					}
					position++
				}
			l652:
				{
					position654, tokenIndex654 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l655
					}
					position++
					goto l654
				l655:
					position, tokenIndex = position654, tokenIndex654
					if buffer[position] != rune('M') {
						goto l648
					}
					position++
				}
			l654:
				{
					position656, tokenIndex656 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l657
					}
					position++
					goto l656
				l657:
					position, tokenIndex = position656, tokenIndex656
					if buffer[position] != rune('I') {
						goto l648
					}
					position++
				}
			l656:
				{
					position658, tokenIndex658 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l659
					}
					position++
					goto l658
				l659:
					position, tokenIndex = position658, tokenIndex658
					if buffer[position] != rune('T') {
						goto l648
					}
					position++
				}
			l658:
				if !_rules[rulesp]() {
					goto l648
				}
				if !_rules[ruleNumericLiteral]() {
					goto l648
				}
				if !_rules[ruleAction26]() {
					goto l648
				}
				add(ruleEmitterLimit, position649)
			}
			return true
		l648:
			position, tokenIndex = position648, tokenIndex648
			return false
		},
		/* 34 EmitterSample <- <(CountBasedSampling / RandomizedSampling / TimeBasedSampling)> */
		func() bool {
			position660, tokenIndex660 := position, tokenIndex
			{
				position661 := position
				{
					position662, tokenIndex662 := position, tokenIndex
					if !_rules[ruleCountBasedSampling]() {
						goto l663
					}
					goto l662
				l663:
					position, tokenIndex = position662, tokenIndex662
					if !_rules[ruleRandomizedSampling]() {
						goto l664
					}
					goto l662
				l664:
					position, tokenIndex = position662, tokenIndex662
					if !_rules[ruleTimeBasedSampling]() {
						goto l660
					}
				}
			l662:
				add(ruleEmitterSample, position661)
			}
			return true
		l660:
			position, tokenIndex = position660, tokenIndex660
			return false
		},
		/* 35 CountBasedSampling <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp NumericLiteral spOpt '-'? spOpt ((('s' / 'S') ('t' / 'T')) / (('n' / 'N') ('d' / 'D')) / (('r' / 'R') ('d' / 'D')) / (('t' / 'T') ('h' / 'H'))) sp (('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E')) Action27)> */
		func() bool {
			position665, tokenIndex665 := position, tokenIndex
			{
				position666 := position
				{
					position667, tokenIndex667 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l668
					}
					position++
					goto l667
				l668:
					position, tokenIndex = position667, tokenIndex667
					if buffer[position] != rune('E') {
						goto l665
					}
					position++
				}
			l667:
				{
					position669, tokenIndex669 := position, tokenIndex
					if buffer[position] != rune('v') {
						goto l670
					}
					position++
					goto l669
				l670:
					position, tokenIndex = position669, tokenIndex669
					if buffer[position] != rune('V') {
						goto l665
					}
					position++
				}
			l669:
				{
					position671, tokenIndex671 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l672
					}
					position++
					goto l671
				l672:
					position, tokenIndex = position671, tokenIndex671
					if buffer[position] != rune('E') {
						goto l665
					}
					position++
				}
			l671:
				{
					position673, tokenIndex673 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l674
					}
					position++
					goto l673
				l674:
					position, tokenIndex = position673, tokenIndex673
					if buffer[position] != rune('R') {
						goto l665
					}
					position++
				}
			l673:
				{
					position675, tokenIndex675 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l676
					}
					position++
					goto l675
				l676:
					position, tokenIndex = position675, tokenIndex675
					if buffer[position] != rune('Y') {
						goto l665
					}
					position++
				}
			l675:
				if !_rules[rulesp]() {
					goto l665
				}
				if !_rules[ruleNumericLiteral]() {
					goto l665
				}
				if !_rules[rulespOpt]() {
					goto l665
				}
				{
					position677, tokenIndex677 := position, tokenIndex
					if buffer[position] != rune('-') {
						goto l677
					}
					position++
					goto l678
				l677:
					position, tokenIndex = position677, tokenIndex677
				}
			l678:
				if !_rules[rulespOpt]() {
					goto l665
				}
				{
					position679, tokenIndex679 := position, tokenIndex
					{
						position681, tokenIndex681 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex = position681, tokenIndex681
						if buffer[position] != rune('S') {
							goto l680
						}
						position++
					}
				l681:
					{
						position683, tokenIndex683 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l684
						}
						position++
						goto l683
					l684:
						position, tokenIndex = position683, tokenIndex683
						if buffer[position] != rune('T') {
							goto l680
						}
						position++
					}
				l683:
					goto l679
				l680:
					position, tokenIndex = position679, tokenIndex679
					{
						position686, tokenIndex686 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex = position686, tokenIndex686
						if buffer[position] != rune('N') {
							goto l685
						}
						position++
					}
				l686:
					{
						position688, tokenIndex688 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex = position688, tokenIndex688
						if buffer[position] != rune('D') {
							goto l685
						}
						position++
					}
				l688:
					goto l679
				l685:
					position, tokenIndex = position679, tokenIndex679
					{
						position691, tokenIndex691 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l692
						}
						position++
						goto l691
					l692:
						position, tokenIndex = position691, tokenIndex691
						if buffer[position] != rune('R') {
							goto l690
						}
						position++
					}
				l691:
					{
						position693, tokenIndex693 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l694
						}
						position++
						goto l693
					l694:
						position, tokenIndex = position693, tokenIndex693
						if buffer[position] != rune('D') {
							goto l690
						}
						position++
					}
				l693:
					goto l679
				l690:
					position, tokenIndex = position679, tokenIndex679
					{
						position695, tokenIndex695 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex = position695, tokenIndex695
						if buffer[position] != rune('T') {
							goto l665
						}
						position++
					}
				l695:
					{
						position697, tokenIndex697 := position, tokenIndex
						if buffer[position] != rune('h') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex = position697, tokenIndex697
						if buffer[position] != rune('H') {
							goto l665
						}
						position++
					}
				l697:
				}
			l679:
				if !_rules[rulesp]() {
					goto l665
				}
				{
					position699, tokenIndex699 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l700
					}
					position++
					goto l699
				l700:
					position, tokenIndex = position699, tokenIndex699
					if buffer[position] != rune('T') {
						goto l665
					}
					position++
				}
			l699:
				{
					position701, tokenIndex701 := position, tokenIndex
					if buffer[position] != rune('u') {
						goto l702
					}
					position++
					goto l701
				l702:
					position, tokenIndex = position701, tokenIndex701
					if buffer[position] != rune('U') {
						goto l665
					}
					position++
				}
			l701:
				{
					position703, tokenIndex703 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l704
					}
					position++
					goto l703
				l704:
					position, tokenIndex = position703, tokenIndex703
					if buffer[position] != rune('P') {
						goto l665
					}
					position++
				}
			l703:
				{
					position705, tokenIndex705 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l706
					}
					position++
					goto l705
				l706:
					position, tokenIndex = position705, tokenIndex705
					if buffer[position] != rune('L') {
						goto l665
					}
					position++
				}
			l705:
				{
					position707, tokenIndex707 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l708
					}
					position++
					goto l707
				l708:
					position, tokenIndex = position707, tokenIndex707
					if buffer[position] != rune('E') {
						goto l665
					}
					position++
				}
			l707:
				if !_rules[ruleAction27]() {
					goto l665
				}
				add(ruleCountBasedSampling, position666)
			}
			return true
		l665:
			position, tokenIndex = position665, tokenIndex665
			return false
		},
		/* 36 RandomizedSampling <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') sp (FloatLiteral / NumericLiteral) spOpt '%' Action28)> */
		func() bool {
			position709, tokenIndex709 := position, tokenIndex
			{
				position710 := position
				{
					position711, tokenIndex711 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l712
					}
					position++
					goto l711
				l712:
					position, tokenIndex = position711, tokenIndex711
					if buffer[position] != rune('S') {
						goto l709
					}
					position++
				}
			l711:
				{
					position713, tokenIndex713 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l714
					}
					position++
					goto l713
				l714:
					position, tokenIndex = position713, tokenIndex713
					if buffer[position] != rune('A') {
						goto l709
					}
					position++
				}
			l713:
				{
					position715, tokenIndex715 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l716
					}
					position++
					goto l715
				l716:
					position, tokenIndex = position715, tokenIndex715
					if buffer[position] != rune('M') {
						goto l709
					}
					position++
				}
			l715:
				{
					position717, tokenIndex717 := position, tokenIndex
					if buffer[position] != rune('p') {
						goto l718
					}
					position++
					goto l717
				l718:
					position, tokenIndex = position717, tokenIndex717
					if buffer[position] != rune('P') {
						goto l709
					}
					position++
				}
			l717:
				{
					position719, tokenIndex719 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l720
					}
					position++
					goto l719
				l720:
					position, tokenIndex = position719, tokenIndex719
					if buffer[position] != rune('L') {
						goto l709
					}
					position++
				}
			l719:
				{
					position721, tokenIndex721 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l722
					}
					position++
					goto l721
				l722:
					position, tokenIndex = position721, tokenIndex721
					if buffer[position] != rune('E') {
						goto l709
					}
					position++
				}
			l721:
				if !_rules[rulesp]() {
					goto l709
				}
				{
					position723, tokenIndex723 := position, tokenIndex
					if !_rules[ruleFloatLiteral]() {
						goto l724
					}
					goto l723
				l724:
					position, tokenIndex = position723, tokenIndex723
					if !_rules[ruleNumericLiteral]() {
						goto l709
					}
				}
			l723:
				if !_rules[rulespOpt]() {
					goto l709
				}
				if buffer[position] != rune('%') {
					goto l709
				}
				position++
				if !_rules[ruleAction28]() {
					goto l709
				}
				add(ruleRandomizedSampling, position710)
			}
			return true
		l709:
			position, tokenIndex = position709, tokenIndex709
			return false
		},
		/* 37 TimeBasedSampling <- <(TimeBasedSamplingSeconds / TimeBasedSamplingMilliseconds)> */
		func() bool {
			position725, tokenIndex725 := position, tokenIndex
			{
				position726 := position
				{
					position727, tokenIndex727 := position, tokenIndex
					if !_rules[ruleTimeBasedSamplingSeconds]() {
						goto l728
					}
					goto l727
				l728:
					position, tokenIndex = position727, tokenIndex727
					if !_rules[ruleTimeBasedSamplingMilliseconds]() {
						goto l725
					}
				}
			l727:
				add(ruleTimeBasedSampling, position726)
			}
			return true
		l725:
			position, tokenIndex = position725, tokenIndex725
			return false
		},
		/* 38 TimeBasedSamplingSeconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action29)> */
		func() bool {
			position729, tokenIndex729 := position, tokenIndex
			{
				position730 := position
				{
					position731, tokenIndex731 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l732
					}
					position++
					goto l731
				l732:
					position, tokenIndex = position731, tokenIndex731
					if buffer[position] != rune('E') {
						goto l729
					}
					position++
				}
			l731:
				{
					position733, tokenIndex733 := position, tokenIndex
					if buffer[position] != rune('v') {
						goto l734
					}
					position++
					goto l733
				l734:
					position, tokenIndex = position733, tokenIndex733
					if buffer[position] != rune('V') {
						goto l729
					}
					position++
				}
			l733:
				{
					position735, tokenIndex735 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l736
					}
					position++
					goto l735
				l736:
					position, tokenIndex = position735, tokenIndex735
					if buffer[position] != rune('E') {
						goto l729
					}
					position++
				}
			l735:
				{
					position737, tokenIndex737 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l738
					}
					position++
					goto l737
				l738:
					position, tokenIndex = position737, tokenIndex737
					if buffer[position] != rune('R') {
						goto l729
					}
					position++
				}
			l737:
				{
					position739, tokenIndex739 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l740
					}
					position++
					goto l739
				l740:
					position, tokenIndex = position739, tokenIndex739
					if buffer[position] != rune('Y') {
						goto l729
					}
					position++
				}
			l739:
				if !_rules[rulesp]() {
					goto l729
				}
				{
					position741, tokenIndex741 := position, tokenIndex
					if !_rules[ruleFloatLiteral]() {
						goto l742
					}
					goto l741
				l742:
					position, tokenIndex = position741, tokenIndex741
					if !_rules[ruleNumericLiteral]() {
						goto l729
					}
				}
			l741:
				if !_rules[rulesp]() {
					goto l729
				}
				{
					position743, tokenIndex743 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l744
					}
					position++
					goto l743
				l744:
					position, tokenIndex = position743, tokenIndex743
					if buffer[position] != rune('S') {
						goto l729
					}
					position++
				}
			l743:
				{
					position745, tokenIndex745 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l746
					}
					position++
					goto l745
				l746:
					position, tokenIndex = position745, tokenIndex745
					if buffer[position] != rune('E') {
						goto l729
					}
					position++
				}
			l745:
				{
					position747, tokenIndex747 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l748
					}
					position++
					goto l747
				l748:
					position, tokenIndex = position747, tokenIndex747
					if buffer[position] != rune('C') {
						goto l729
					}
					position++
				}
			l747:
				{
					position749, tokenIndex749 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l750
					}
					position++
					goto l749
				l750:
					position, tokenIndex = position749, tokenIndex749
					if buffer[position] != rune('O') {
						goto l729
					}
					position++
				}
			l749:
				{
					position751, tokenIndex751 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l752
					}
					position++
					goto l751
				l752:
					position, tokenIndex = position751, tokenIndex751
					if buffer[position] != rune('N') {
						goto l729
					}
					position++
				}
			l751:
				{
					position753, tokenIndex753 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l754
					}
					position++
					goto l753
				l754:
					position, tokenIndex = position753, tokenIndex753
					if buffer[position] != rune('D') {
						goto l729
					}
					position++
				}
			l753:
				{
					position755, tokenIndex755 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l756
					}
					position++
					goto l755
				l756:
					position, tokenIndex = position755, tokenIndex755
					if buffer[position] != rune('S') {
						goto l729
					}
					position++
				}
			l755:
				if !_rules[ruleAction29]() {
					goto l729
				}
				add(ruleTimeBasedSamplingSeconds, position730)
			}
			return true
		l729:
			position, tokenIndex = position729, tokenIndex729
			return false
		},
		/* 39 TimeBasedSamplingMilliseconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action30)> */
		func() bool {
			position757, tokenIndex757 := position, tokenIndex
			{
				position758 := position
				{
					position759, tokenIndex759 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l760
					}
					position++
					goto l759
				l760:
					position, tokenIndex = position759, tokenIndex759
					if buffer[position] != rune('E') {
						goto l757
					}
					position++
				}
			l759:
				{
					position761, tokenIndex761 := position, tokenIndex
					if buffer[position] != rune('v') {
						goto l762
					}
					position++
					goto l761
				l762:
					position, tokenIndex = position761, tokenIndex761
					if buffer[position] != rune('V') {
						goto l757
					}
					position++
				}
			l761:
				{
					position763, tokenIndex763 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l764
					}
					position++
					goto l763
				l764:
					position, tokenIndex = position763, tokenIndex763
					if buffer[position] != rune('E') {
						goto l757
					}
					position++
				}
			l763:
				{
					position765, tokenIndex765 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l766
					}
					position++
					goto l765
				l766:
					position, tokenIndex = position765, tokenIndex765
					if buffer[position] != rune('R') {
						goto l757
					}
					position++
				}
			l765:
				{
					position767, tokenIndex767 := position, tokenIndex
					if buffer[position] != rune('y') {
						goto l768
					}
					position++
					goto l767
				l768:
					position, tokenIndex = position767, tokenIndex767
					if buffer[position] != rune('Y') {
						goto l757
					}
					position++
				}
			l767:
				if !_rules[rulesp]() {
					goto l757
				}
				{
					position769, tokenIndex769 := position, tokenIndex
					if !_rules[ruleFloatLiteral]() {
						goto l770
					}
					goto l769
				l770:
					position, tokenIndex = position769, tokenIndex769
					if !_rules[ruleNumericLiteral]() {
						goto l757
					}
				}
			l769:
				if !_rules[rulesp]() {
					goto l757
				}
				{
					position771, tokenIndex771 := position, tokenIndex
					if buffer[position] != rune('m') {
						goto l772
					}
					position++
					goto l771
				l772:
					position, tokenIndex = position771, tokenIndex771
					if buffer[position] != rune('M') {
						goto l757
					}
					position++
				}
			l771:
				{
					position773, tokenIndex773 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l774
					}
					position++
					goto l773
				l774:
					position, tokenIndex = position773, tokenIndex773
					if buffer[position] != rune('I') {
						goto l757
					}
					position++
				}
			l773:
				{
					position775, tokenIndex775 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l776
					}
					position++
					goto l775
				l776:
					position, tokenIndex = position775, tokenIndex775
					if buffer[position] != rune('L') {
						goto l757
					}
					position++
				}
			l775:
				{
					position777, tokenIndex777 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l778
					}
					position++
					goto l777
				l778:
					position, tokenIndex = position777, tokenIndex777
					if buffer[position] != rune('L') {
						goto l757
					}
					position++
				}
			l777:
				{
					position779, tokenIndex779 := position, tokenIndex
					if buffer[position] != rune('i') {
						goto l780
					}
					position++
					goto l779
				l780:
					position, tokenIndex = position779, tokenIndex779
					if buffer[position] != rune('I') {
						goto l757
					}
					position++
				}
			l779:
				{
					position781, tokenIndex781 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l782
					}
					position++
					goto l781
				l782:
					position, tokenIndex = position781, tokenIndex781
					if buffer[position] != rune('S') {
						goto l757
					}
					position++
				}
			l781:
				{
					position783, tokenIndex783 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l784
					}
					position++
					goto l783
				l784:
					position, tokenIndex = position783, tokenIndex783
					if buffer[position] != rune('E') {
						goto l757
					}
					position++
				}
			l783:
				{
					position785, tokenIndex785 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l786
					}
					position++
					goto l785
				l786:
					position, tokenIndex = position785, tokenIndex785
					if buffer[position] != rune('C') {
						goto l757
					}
					position++
				}
			l785:
				{
					position787, tokenIndex787 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l788
					}
					position++
					goto l787
				l788:
					position, tokenIndex = position787, tokenIndex787
					if buffer[position] != rune('O') {
						goto l757
					}
					position++
				}
			l787:
				{
					position789, tokenIndex789 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l790
					}
					position++
					goto l789
				l790:
					position, tokenIndex = position789, tokenIndex789
					if buffer[position] != rune('N') {
						goto l757
					}
					position++
				}
			l789:
				{
					position791, tokenIndex791 := position, tokenIndex
					if buffer[position] != rune('d') {
						goto l792
					}
					position++
					goto l791
				l792:
					position, tokenIndex = position791, tokenIndex791
					if buffer[position] != rune('D') {
						goto l757
					}
					position++
				}
			l791:
				{
					position793, tokenIndex793 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l794
					}
					position++
					goto l793
				l794:
					position, tokenIndex = position793, tokenIndex793
					if buffer[position] != rune('S') {
						goto l757
					}
					position++
				}
			l793:
				if !_rules[ruleAction30]() {
					goto l757
				}
				add(ruleTimeBasedSamplingMilliseconds, position758)
			}
			return true
		l757:
			position, tokenIndex = position757, tokenIndex757
			return false
		},
		/* 40 Projections <- <(<(sp Projection (spOpt ',' spOpt Projection)*)> Action31)> */
		func() bool {
			position795, tokenIndex795 := position, tokenIndex
			{
				position796 := position
				{
					position797 := position
					if !_rules[rulesp]() {
						goto l795
					}
					if !_rules[ruleProjection]() {
						goto l795
					}
				l798:
					{
						position799, tokenIndex799 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l799
						}
						if buffer[position] != rune(',') {
							goto l799
						}
						position++
						if !_rules[rulespOpt]() {
							goto l799
						}
						if !_rules[ruleProjection]() {
							goto l799
						}
						goto l798
					l799:
						position, tokenIndex = position799, tokenIndex799
					}
					add(rulePegText, position797)
				}
				if !_rules[ruleAction31]() {
					goto l795
				}
				add(ruleProjections, position796)
			}
			return true
		l795:
			position, tokenIndex = position795, tokenIndex795
			return false
		},
		/* 41 Projection <- <(AliasExpression / ExpressionOrWildcard)> */
		func() bool {
			position800, tokenIndex800 := position, tokenIndex
			{
				position801 := position
				{
					position802, tokenIndex802 := position, tokenIndex
					if !_rules[ruleAliasExpression]() {
						goto l803
					}
					goto l802
				l803:
					position, tokenIndex = position802, tokenIndex802
					if !_rules[ruleExpressionOrWildcard]() {
						goto l800
					}
				}
			l802:
				add(ruleProjection, position801)
			}
			return true
		l800:
			position, tokenIndex = position800, tokenIndex800
			return false
		},
		/* 42 AliasExpression <- <(ExpressionOrWildcard sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action32)> */
		func() bool {
			position804, tokenIndex804 := position, tokenIndex
			{
				position805 := position
				if !_rules[ruleExpressionOrWildcard]() {
					goto l804
				}
				if !_rules[rulesp]() {
					goto l804
				}
				{
					position806, tokenIndex806 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l807
					}
					position++
					goto l806
				l807:
					position, tokenIndex = position806, tokenIndex806
					if buffer[position] != rune('A') {
						goto l804
					}
					position++
				}
			l806:
				{
					position808, tokenIndex808 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l809
					}
					position++
					goto l808
				l809:
					position, tokenIndex = position808, tokenIndex808
					if buffer[position] != rune('S') {
						goto l804
					}
					position++
				}
			l808:
				if !_rules[rulesp]() {
					goto l804
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l804
				}
				if !_rules[ruleAction32]() {
					goto l804
				}
				add(ruleAliasExpression, position805)
			}
			return true
		l804:
			position, tokenIndex = position804, tokenIndex804
			return false
		},
		/* 43 WindowedFrom <- <(<(sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp Relations)?> Action33)> */
		func() bool {
			position810, tokenIndex810 := position, tokenIndex
			{
				position811 := position
				{
					position812 := position
					{
						position813, tokenIndex813 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l813
						}
						{
							position815, tokenIndex815 := position, tokenIndex
							if buffer[position] != rune('f') {
								goto l816
							}
							position++
							goto l815
						l816:
							position, tokenIndex = position815, tokenIndex815
							if buffer[position] != rune('F') {
								goto l813
							}
							position++
						}
					l815:
						{
							position817, tokenIndex817 := position, tokenIndex
							if buffer[position] != rune('r') {
								goto l818
							}
							position++
							goto l817
						l818:
							position, tokenIndex = position817, tokenIndex817
							if buffer[position] != rune('R') {
								goto l813
							}
							position++
						}
					l817:
						{
							position819, tokenIndex819 := position, tokenIndex
							if buffer[position] != rune('o') {
								goto l820
							}
							position++
							goto l819
						l820:
							position, tokenIndex = position819, tokenIndex819
							if buffer[position] != rune('O') {
								goto l813
							}
							position++
						}
					l819:
						{
							position821, tokenIndex821 := position, tokenIndex
							if buffer[position] != rune('m') {
								goto l822
							}
							position++
							goto l821
						l822:
							position, tokenIndex = position821, tokenIndex821
							if buffer[position] != rune('M') {
								goto l813
							}
							position++
						}
					l821:
						if !_rules[rulesp]() {
							goto l813
						}
						if !_rules[ruleRelations]() {
							goto l813
						}
						goto l814
					l813:
						position, tokenIndex = position813, tokenIndex813
					}
				l814:
					add(rulePegText, position812)
				}
				if !_rules[ruleAction33]() {
					goto l810
				}
				add(ruleWindowedFrom, position811)
			}
			return true
		l810:
			position, tokenIndex = position810, tokenIndex810
			return false
		},
		/* 44 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position823, tokenIndex823 := position, tokenIndex
			{
				position824 := position
				{
					position825, tokenIndex825 := position, tokenIndex
					if !_rules[ruleTimeInterval]() {
						goto l826
					}
					goto l825
				l826:
					position, tokenIndex = position825, tokenIndex825
					if !_rules[ruleTuplesInterval]() {
						goto l823
					}
				}
			l825:
				add(ruleInterval, position824)
			}
			return true
		l823:
			position, tokenIndex = position823, tokenIndex823
			return false
		},
		/* 45 TimeInterval <- <((FloatLiteral / NumericLiteral) sp (SECONDS / MILLISECONDS) Action34)> */
		func() bool {
			position827, tokenIndex827 := position, tokenIndex
			{
				position828 := position
				{
					position829, tokenIndex829 := position, tokenIndex
					if !_rules[ruleFloatLiteral]() {
						goto l830
					}
					goto l829
				l830:
					position, tokenIndex = position829, tokenIndex829
					if !_rules[ruleNumericLiteral]() {
						goto l827
					}
				}
			l829:
				if !_rules[rulesp]() {
					goto l827
				}
				{
					position831, tokenIndex831 := position, tokenIndex
					if !_rules[ruleSECONDS]() {
						goto l832
					}
					goto l831
				l832:
					position, tokenIndex = position831, tokenIndex831
					if !_rules[ruleMILLISECONDS]() {
						goto l827
					}
				}
			l831:
				if !_rules[ruleAction34]() {
					goto l827
				}
				add(ruleTimeInterval, position828)
			}
			return true
		l827:
			position, tokenIndex = position827, tokenIndex827
			return false
		},
		/* 46 TuplesInterval <- <(NumericLiteral sp TUPLES Action35)> */
		func() bool {
			position833, tokenIndex833 := position, tokenIndex
			{
				position834 := position
				if !_rules[ruleNumericLiteral]() {
					goto l833
				}
				if !_rules[rulesp]() {
					goto l833
				}
				if !_rules[ruleTUPLES]() {
					goto l833
				}
				if !_rules[ruleAction35]() {
					goto l833
				}
				add(ruleTuplesInterval, position834)
			}
			return true
		l833:
			position, tokenIndex = position833, tokenIndex833
			return false
		},
		/* 47 Relations <- <(RelationLike (spOpt ',' spOpt RelationLike)*)> */
		func() bool {
			position835, tokenIndex835 := position, tokenIndex
			{
				position836 := position
				if !_rules[ruleRelationLike]() {
					goto l835
				}
			l837:
				{
					position838, tokenIndex838 := position, tokenIndex
					if !_rules[rulespOpt]() {
						goto l838
					}
					if buffer[position] != rune(',') {
						goto l838
					}
					position++
					if !_rules[rulespOpt]() {
						goto l838
					}
					if !_rules[ruleRelationLike]() {
						goto l838
					}
					goto l837
				l838:
					position, tokenIndex = position838, tokenIndex838
				}
				add(ruleRelations, position836)
			}
			return true
		l835:
			position, tokenIndex = position835, tokenIndex835
			return false
		},
		/* 48 Filter <- <(<(sp (('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E')) sp Expression)?> Action36)> */
		func() bool {
			position839, tokenIndex839 := position, tokenIndex
			{
				position840 := position
				{
					position841 := position
					{
						position842, tokenIndex842 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l842
						}
						{
							position844, tokenIndex844 := position, tokenIndex
							if buffer[position] != rune('w') {
								goto l845
							}
							position++
							goto l844
						l845:
							position, tokenIndex = position844, tokenIndex844
							if buffer[position] != rune('W') {
								goto l842
							}
							position++
						}
					l844:
						{
							position846, tokenIndex846 := position, tokenIndex
							if buffer[position] != rune('h') {
								goto l847
							}
							position++
							goto l846
						l847:
							position, tokenIndex = position846, tokenIndex846
							if buffer[position] != rune('H') {
								goto l842
							}
							position++
						}
					l846:
						{
							position848, tokenIndex848 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l849
							}
							position++
							goto l848
						l849:
							position, tokenIndex = position848, tokenIndex848
							if buffer[position] != rune('E') {
								goto l842
							}
							position++
						}
					l848:
						{
							position850, tokenIndex850 := position, tokenIndex
							if buffer[position] != rune('r') {
								goto l851
							}
							position++
							goto l850
						l851:
							position, tokenIndex = position850, tokenIndex850
							if buffer[position] != rune('R') {
								goto l842
							}
							position++
						}
					l850:
						{
							position852, tokenIndex852 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l853
							}
							position++
							goto l852
						l853:
							position, tokenIndex = position852, tokenIndex852
							if buffer[position] != rune('E') {
								goto l842
							}
							position++
						}
					l852:
						if !_rules[rulesp]() {
							goto l842
						}
						if !_rules[ruleExpression]() {
							goto l842
						}
						goto l843
					l842:
						position, tokenIndex = position842, tokenIndex842
					}
				l843:
					add(rulePegText, position841)
				}
				if !_rules[ruleAction36]() {
					goto l839
				}
				add(ruleFilter, position840)
			}
			return true
		l839:
			position, tokenIndex = position839, tokenIndex839
			return false
		},
		/* 49 Grouping <- <(<(sp (('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P')) sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action37)> */
		func() bool {
			position854, tokenIndex854 := position, tokenIndex
			{
				position855 := position
				{
					position856 := position
					{
						position857, tokenIndex857 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l857
						}
						{
							position859, tokenIndex859 := position, tokenIndex
							if buffer[position] != rune('g') {
								goto l860
							}
							position++
							goto l859
						l860:
							position, tokenIndex = position859, tokenIndex859
							if buffer[position] != rune('G') {
								goto l857
							}
							position++
						}
					l859:
						{
							position861, tokenIndex861 := position, tokenIndex
							if buffer[position] != rune('r') {
								goto l862
							}
							position++
							goto l861
						l862:
							position, tokenIndex = position861, tokenIndex861
							if buffer[position] != rune('R') {
								goto l857
							}
							position++
						}
					l861:
						{
							position863, tokenIndex863 := position, tokenIndex
							if buffer[position] != rune('o') {
								goto l864
							}
							position++
							goto l863
						l864:
							position, tokenIndex = position863, tokenIndex863
							if buffer[position] != rune('O') {
								goto l857
							}
							position++
						}
					l863:
						{
							position865, tokenIndex865 := position, tokenIndex
							if buffer[position] != rune('u') {
								goto l866
							}
							position++
							goto l865
						l866:
							position, tokenIndex = position865, tokenIndex865
							if buffer[position] != rune('U') {
								goto l857
							}
							position++
						}
					l865:
						{
							position867, tokenIndex867 := position, tokenIndex
							if buffer[position] != rune('p') {
								goto l868
							}
							position++
							goto l867
						l868:
							position, tokenIndex = position867, tokenIndex867
							if buffer[position] != rune('P') {
								goto l857
							}
							position++
						}
					l867:
						if !_rules[rulesp]() {
							goto l857
						}
						{
							position869, tokenIndex869 := position, tokenIndex
							if buffer[position] != rune('b') {
								goto l870
							}
							position++
							goto l869
						l870:
							position, tokenIndex = position869, tokenIndex869
							if buffer[position] != rune('B') {
								goto l857
							}
							position++
						}
					l869:
						{
							position871, tokenIndex871 := position, tokenIndex
							if buffer[position] != rune('y') {
								goto l872
							}
							position++
							goto l871
						l872:
							position, tokenIndex = position871, tokenIndex871
							if buffer[position] != rune('Y') {
								goto l857
							}
							position++
						}
					l871:
						if !_rules[rulesp]() {
							goto l857
						}
						if !_rules[ruleGroupList]() {
							goto l857
						}
						goto l858
					l857:
						position, tokenIndex = position857, tokenIndex857
					}
				l858:
					add(rulePegText, position856)
				}
				if !_rules[ruleAction37]() {
					goto l854
				}
				add(ruleGrouping, position855)
			}
			return true
		l854:
			position, tokenIndex = position854, tokenIndex854
			return false
		},
		/* 50 GroupList <- <(Expression (spOpt ',' spOpt Expression)*)> */
		func() bool {
			position873, tokenIndex873 := position, tokenIndex
			{
				position874 := position
				if !_rules[ruleExpression]() {
					goto l873
				}
			l875:
				{
					position876, tokenIndex876 := position, tokenIndex
					if !_rules[rulespOpt]() {
						goto l876
					}
					if buffer[position] != rune(',') {
						goto l876
					}
					position++
					if !_rules[rulespOpt]() {
						goto l876
					}
					if !_rules[ruleExpression]() {
						goto l876
					}
					goto l875
				l876:
					position, tokenIndex = position876, tokenIndex876
				}
				add(ruleGroupList, position874)
			}
			return true
		l873:
			position, tokenIndex = position873, tokenIndex873
			return false
		},
		/* 51 Having <- <(<(sp (('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G')) sp Expression)?> Action38)> */
		func() bool {
			position877, tokenIndex877 := position, tokenIndex
			{
				position878 := position
				{
					position879 := position
					{
						position880, tokenIndex880 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l880
						}
						{
							position882, tokenIndex882 := position, tokenIndex
							if buffer[position] != rune('h') {
								goto l883
							}
							position++
							goto l882
						l883:
							position, tokenIndex = position882, tokenIndex882
							if buffer[position] != rune('H') {
								goto l880
							}
							position++
						}
					l882:
						{
							position884, tokenIndex884 := position, tokenIndex
							if buffer[position] != rune('a') {
								goto l885
							}
							position++
							goto l884
						l885:
							position, tokenIndex = position884, tokenIndex884
							if buffer[position] != rune('A') {
								goto l880
							}
							position++
						}
					l884:
						{
							position886, tokenIndex886 := position, tokenIndex
							if buffer[position] != rune('v') {
								goto l887
							}
							position++
							goto l886
						l887:
							position, tokenIndex = position886, tokenIndex886
							if buffer[position] != rune('V') {
								goto l880
							}
							position++
						}
					l886:
						{
							position888, tokenIndex888 := position, tokenIndex
							if buffer[position] != rune('i') {
								goto l889
							}
							position++
							goto l888
						l889:
							position, tokenIndex = position888, tokenIndex888
							if buffer[position] != rune('I') {
								goto l880
							}
							position++
						}
					l888:
						{
							position890, tokenIndex890 := position, tokenIndex
							if buffer[position] != rune('n') {
								goto l891
							}
							position++
							goto l890
						l891:
							position, tokenIndex = position890, tokenIndex890
							if buffer[position] != rune('N') {
								goto l880
							}
							position++
						}
					l890:
						{
							position892, tokenIndex892 := position, tokenIndex
							if buffer[position] != rune('g') {
								goto l893
							}
							position++
							goto l892
						l893:
							position, tokenIndex = position892, tokenIndex892
							if buffer[position] != rune('G') {
								goto l880
							}
							position++
						}
					l892:
						if !_rules[rulesp]() {
							goto l880
						}
						if !_rules[ruleExpression]() {
							goto l880
						}
						goto l881
					l880:
						position, tokenIndex = position880, tokenIndex880
					}
				l881:
					add(rulePegText, position879)
				}
				if !_rules[ruleAction38]() {
					goto l877
				}
				add(ruleHaving, position878)
			}
			return true
		l877:
			position, tokenIndex = position877, tokenIndex877
			return false
		},
		/* 52 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action39))> */
		func() bool {
			position894, tokenIndex894 := position, tokenIndex
			{
				position895 := position
				{
					position896, tokenIndex896 := position, tokenIndex
					if !_rules[ruleAliasedStreamWindow]() {
						goto l897
					}
					goto l896
				l897:
					position, tokenIndex = position896, tokenIndex896
					if !_rules[ruleStreamWindow]() {
						goto l894
					}
					if !_rules[ruleAction39]() {
						goto l894
					}
				}
			l896:
				add(ruleRelationLike, position895)
			}
			return true
		l894:
			position, tokenIndex = position894, tokenIndex894
			return false
		},
		/* 53 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action40)> */
		func() bool {
			position898, tokenIndex898 := position, tokenIndex
			{
				position899 := position
				if !_rules[ruleStreamWindow]() {
					goto l898
				}
				if !_rules[rulesp]() {
					goto l898
				}
				{
					position900, tokenIndex900 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l901
					}
					position++
					goto l900
				l901:
					position, tokenIndex = position900, tokenIndex900
					if buffer[position] != rune('A') {
						goto l898
					}
					position++
				}
			l900:
				{
					position902, tokenIndex902 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l903
					}
					position++
					goto l902
				l903:
					position, tokenIndex = position902, tokenIndex902
					if buffer[position] != rune('S') {
						goto l898
					}
					position++
				}
			l902:
				if !_rules[rulesp]() {
					goto l898
				}
				if !_rules[ruleIdentifier]() {
					goto l898
				}
				if !_rules[ruleAction40]() {
					goto l898
				}
				add(ruleAliasedStreamWindow, position899)
			}
			return true
		l898:
			position, tokenIndex = position898, tokenIndex898
			return false
		},
		/* 54 StreamWindow <- <(StreamLike spOpt '[' spOpt (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval CapacitySpecOpt SheddingSpecOpt spOpt ']' Action41)> */
		func() bool {
			position904, tokenIndex904 := position, tokenIndex
			{
				position905 := position
				if !_rules[ruleStreamLike]() {
					goto l904
				}
				if !_rules[rulespOpt]() {
					goto l904
				}
				if buffer[position] != rune('[') {
					goto l904
				}
				position++
				if !_rules[rulespOpt]() {
					goto l904
				}
				{
					position906, tokenIndex906 := position, tokenIndex
					if buffer[position] != rune('r') {
						goto l907
					}
					position++
					goto l906
				l907:
					position, tokenIndex = position906, tokenIndex906
					if buffer[position] != rune('R') {
						goto l904
					}
					position++
				}
			l906:
				{
					position908, tokenIndex908 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l909
					}
					position++
					goto l908
				l909:
					position, tokenIndex = position908, tokenIndex908
					if buffer[position] != rune('A') {
						goto l904
					}
					position++
				}
			l908:
				{
					position910, tokenIndex910 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l911
					}
					position++
					goto l910
				l911:
					position, tokenIndex = position910, tokenIndex910
					if buffer[position] != rune('N') {
						goto l904
					}
					position++
				}
			l910:
				{
					position912, tokenIndex912 := position, tokenIndex
					if buffer[position] != rune('g') {
						goto l913
					}
					position++
					goto l912
				l913:
					position, tokenIndex = position912, tokenIndex912
					if buffer[position] != rune('G') {
						goto l904
					}
					position++
				}
			l912:
				{
					position914, tokenIndex914 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l915
					}
					position++
					goto l914
				l915:
					position, tokenIndex = position914, tokenIndex914
					if buffer[position] != rune('E') {
						goto l904
					}
					position++
				}
			l914:
				if !_rules[rulesp]() {
					goto l904
				}
				if !_rules[ruleInterval]() {
					goto l904
				}
				if !_rules[ruleCapacitySpecOpt]() {
					goto l904
				}
				if !_rules[ruleSheddingSpecOpt]() {
					goto l904
				}
				if !_rules[rulespOpt]() {
					goto l904
				}
				if buffer[position] != rune(']') {
					goto l904
				}
				position++
				if !_rules[ruleAction41]() {
					goto l904
				}
				add(ruleStreamWindow, position905)
			}
			return true
		l904:
			position, tokenIndex = position904, tokenIndex904
			return false
		},
		/* 55 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position916, tokenIndex916 := position, tokenIndex
			{
				position917 := position
				{
					position918, tokenIndex918 := position, tokenIndex
					if !_rules[ruleUDSFFuncApp]() {
						goto l919
					}
					goto l918
				l919:
					position, tokenIndex = position918, tokenIndex918
					if !_rules[ruleStream]() {
						goto l916
					}
				}
			l918:
				add(ruleStreamLike, position917)
			}
			return true
		l916:
			position, tokenIndex = position916, tokenIndex916
			return false
		},
		/* 56 UDSFFuncApp <- <(FuncAppWithoutOrderBy Action42)> */
		func() bool {
			position920, tokenIndex920 := position, tokenIndex
			{
				position921 := position
				if !_rules[ruleFuncAppWithoutOrderBy]() {
					goto l920
				}
				if !_rules[ruleAction42]() {
					goto l920
				}
				add(ruleUDSFFuncApp, position921)
			}
			return true
		l920:
			position, tokenIndex = position920, tokenIndex920
			return false
		},
		/* 57 CapacitySpecOpt <- <(<(spOpt ',' spOpt (('b' / 'B') ('u' / 'U') ('f' / 'F') ('f' / 'F') ('e' / 'E') ('r' / 'R')) sp (('s' / 'S') ('i' / 'I') ('z' / 'Z') ('e' / 'E')) sp NonNegativeNumericLiteral)?> Action43)> */
		func() bool {
			position922, tokenIndex922 := position, tokenIndex
			{
				position923 := position
				{
					position924 := position
					{
						position925, tokenIndex925 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l925
						}
						if buffer[position] != rune(',') {
							goto l925
						}
						position++
						if !_rules[rulespOpt]() {
							goto l925
						}
						{
							position927, tokenIndex927 := position, tokenIndex
							if buffer[position] != rune('b') {
								goto l928
							}
							position++
							goto l927
						l928:
							position, tokenIndex = position927, tokenIndex927
							if buffer[position] != rune('B') {
								goto l925
							}
							position++
						}
					l927:
						{
							position929, tokenIndex929 := position, tokenIndex
							if buffer[position] != rune('u') {
								goto l930
							}
							position++
							goto l929
						l930:
							position, tokenIndex = position929, tokenIndex929
							if buffer[position] != rune('U') {
								goto l925
							}
							position++
						}
					l929:
						{
							position931, tokenIndex931 := position, tokenIndex
							if buffer[position] != rune('f') {
								goto l932
							}
							position++
							goto l931
						l932:
							position, tokenIndex = position931, tokenIndex931
							if buffer[position] != rune('F') {
								goto l925
							}
							position++
						}
					l931:
						{
							position933, tokenIndex933 := position, tokenIndex
							if buffer[position] != rune('f') {
								goto l934
							}
							position++
							goto l933
						l934:
							position, tokenIndex = position933, tokenIndex933
							if buffer[position] != rune('F') {
								goto l925
							}
							position++
						}
					l933:
						{
							position935, tokenIndex935 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l936
							}
							position++
							goto l935
						l936:
							position, tokenIndex = position935, tokenIndex935
							if buffer[position] != rune('E') {
								goto l925
							}
							position++
						}
					l935:
						{
							position937, tokenIndex937 := position, tokenIndex
							if buffer[position] != rune('r') {
								goto l938
							}
							position++
							goto l937
						l938:
							position, tokenIndex = position937, tokenIndex937
							if buffer[position] != rune('R') {
								goto l925
							}
							position++
						}
					l937:
						if !_rules[rulesp]() {
							goto l925
						}
						{
							position939, tokenIndex939 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l940
							}
							position++
							goto l939
						l940:
							position, tokenIndex = position939, tokenIndex939
							if buffer[position] != rune('S') {
								goto l925
							}
							position++
						}
					l939:
						{
							position941, tokenIndex941 := position, tokenIndex
							if buffer[position] != rune('i') {
								goto l942
							}
							position++
							goto l941
						l942:
							position, tokenIndex = position941, tokenIndex941
							if buffer[position] != rune('I') {
								goto l925
							}
							position++
						}
					l941:
						{
							position943, tokenIndex943 := position, tokenIndex
							if buffer[position] != rune('z') {
								goto l944
							}
							position++
							goto l943
						l944:
							position, tokenIndex = position943, tokenIndex943
							if buffer[position] != rune('Z') {
								goto l925
							}
							position++
						}
					l943:
						{
							position945, tokenIndex945 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l946
							}
							position++
							goto l945
						l946:
							position, tokenIndex = position945, tokenIndex945
							if buffer[position] != rune('E') {
								goto l925
							}
							position++
						}
					l945:
						if !_rules[rulesp]() {
							goto l925
						}
						if !_rules[ruleNonNegativeNumericLiteral]() {
							goto l925
						}
						goto l926
					l925:
						position, tokenIndex = position925, tokenIndex925
					}
				l926:
					add(rulePegText, position924)
				}
				if !_rules[ruleAction43]() {
					goto l922
				}
				add(ruleCapacitySpecOpt, position923)
			}
			return true
		l922:
			position, tokenIndex = position922, tokenIndex922
			return false
		},
		/* 58 SheddingSpecOpt <- <(<(spOpt ',' spOpt SheddingOption sp (('i' / 'I') ('f' / 'F')) sp (('f' / 'F') ('u' / 'U') ('l' / 'L') ('l' / 'L')))?> Action44)> */
		func() bool {
			position947, tokenIndex947 := position, tokenIndex
			{
				position948 := position
				{
					position949 := position
					{
						position950, tokenIndex950 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l950
						}
						if buffer[position] != rune(',') {
							goto l950
						}
						position++
						if !_rules[rulespOpt]() {
							goto l950
						}
						if !_rules[ruleSheddingOption]() {
							goto l950
						}
						if !_rules[rulesp]() {
							goto l950
						}
						{
							position952, tokenIndex952 := position, tokenIndex
							if buffer[position] != rune('i') {
								goto l953
							}
							position++
							goto l952
						l953:
							position, tokenIndex = position952, tokenIndex952
							if buffer[position] != rune('I') {
								goto l950
							}
							position++
						}
					l952:
						{
							position954, tokenIndex954 := position, tokenIndex
							if buffer[position] != rune('f') {
								goto l955
							}
							position++
							goto l954
						l955:
							position, tokenIndex = position954, tokenIndex954
							if buffer[position] != rune('F') {
								goto l950
							}
							position++
						}
					l954:
						if !_rules[rulesp]() {
							goto l950
						}
						{
							position956, tokenIndex956 := position, tokenIndex
							if buffer[position] != rune('f') {
								goto l957
							}
							position++
							goto l956
						l957:
							position, tokenIndex = position956, tokenIndex956
							if buffer[position] != rune('F') {
								goto l950
							}
							position++
						}
					l956:
						{
							position958, tokenIndex958 := position, tokenIndex
							if buffer[position] != rune('u') {
								goto l959
							}
							position++
							goto l958
						l959:
							position, tokenIndex = position958, tokenIndex958
							if buffer[position] != rune('U') {
								goto l950
							}
							position++
						}
					l958:
						{
							position960, tokenIndex960 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l961
							}
							position++
							goto l960
						l961:
							position, tokenIndex = position960, tokenIndex960
							if buffer[position] != rune('L') {
								goto l950
							}
							position++
						}
					l960:
						{
							position962, tokenIndex962 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l963
							}
							position++
							goto l962
						l963:
							position, tokenIndex = position962, tokenIndex962
							if buffer[position] != rune('L') {
								goto l950
							}
							position++
						}
					l962:
						goto l951
					l950:
						position, tokenIndex = position950, tokenIndex950
					}
				l951:
					add(rulePegText, position949)
				}
				if !_rules[ruleAction44]() {
					goto l947
				}
				add(ruleSheddingSpecOpt, position948)
			}
			return true
		l947:
			position, tokenIndex = position947, tokenIndex947
			return false
		},
		/* 59 SheddingOption <- <(Wait / DropOldest / DropNewest)> */
		func() bool {
			position964, tokenIndex964 := position, tokenIndex
			{
				position965 := position
				{
					position966, tokenIndex966 := position, tokenIndex
					if !_rules[ruleWait]() {
						goto l967
					}
					goto l966
				l967:
					position, tokenIndex = position966, tokenIndex966
					if !_rules[ruleDropOldest]() {
						goto l968
					}
					goto l966
				l968:
					position, tokenIndex = position966, tokenIndex966
					if !_rules[ruleDropNewest]() {
						goto l964
					}
				}
			l966:
				add(ruleSheddingOption, position965)
			}
			return true
		l964:
			position, tokenIndex = position964, tokenIndex964
			return false
		},
		/* 60 SourceSinkSpecs <- <(<(sp (('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action45)> */
		func() bool {
			position969, tokenIndex969 := position, tokenIndex
			{
				position970 := position
				{
					position971 := position
					{
						position972, tokenIndex972 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l972
						}
						{
							position974, tokenIndex974 := position, tokenIndex
							if buffer[position] != rune('w') {
								goto l975
							}
							position++
							goto l974
						l975:
							position, tokenIndex = position974, tokenIndex974
							if buffer[position] != rune('W') {
								goto l972
							}
							position++
						}
					l974:
						{
							position976, tokenIndex976 := position, tokenIndex
							if buffer[position] != rune('i') {
								goto l977
							}
							position++
							goto l976
						l977:
							position, tokenIndex = position976, tokenIndex976
							if buffer[position] != rune('I') {
								goto l972
							}
							position++
						}
					l976:
						{
							position978, tokenIndex978 := position, tokenIndex
							if buffer[position] != rune('t') {
								goto l979
							}
							position++
							goto l978
						l979:
							position, tokenIndex = position978, tokenIndex978
							if buffer[position] != rune('T') {
								goto l972
							}
							position++
						}
					l978:
						{
							position980, tokenIndex980 := position, tokenIndex
							if buffer[position] != rune('h') {
								goto l981
							}
							position++
							goto l980
						l981:
							position, tokenIndex = position980, tokenIndex980
							if buffer[position] != rune('H') {
								goto l972
							}
							position++
						}
					l980:
						if !_rules[rulesp]() {
							goto l972
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l972
						}
					l982:
						{
							position983, tokenIndex983 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l983
							}
							if buffer[position] != rune(',') {
								goto l983
							}
							position++
							if !_rules[rulespOpt]() {
								goto l983
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l983
							}
							goto l982
						l983:
							position, tokenIndex = position983, tokenIndex983
						}
						goto l973
					l972:
						position, tokenIndex = position972, tokenIndex972
					}
				l973:
					add(rulePegText, position971)
				}
				if !_rules[ruleAction45]() {
					goto l969
				}
				add(ruleSourceSinkSpecs, position970)
			}
			return true
		l969:
			position, tokenIndex = position969, tokenIndex969
			return false
		},
		/* 61 UpdateSourceSinkSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)> Action46)> */
		func() bool {
			position984, tokenIndex984 := position, tokenIndex
			{
				position985 := position
				{
					position986 := position
					if !_rules[rulesp]() {
						goto l984
					}
					{
						position987, tokenIndex987 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l988
						}
						position++
						goto l987
					l988:
						position, tokenIndex = position987, tokenIndex987
						if buffer[position] != rune('S') {
							goto l984
						}
						position++
					}
				l987:
					{
						position989, tokenIndex989 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l990
						}
						position++
						goto l989
					l990:
						position, tokenIndex = position989, tokenIndex989
						if buffer[position] != rune('E') {
							goto l984
						}
						position++
					}
				l989:
					{
						position991, tokenIndex991 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l992
						}
						position++
						goto l991
					l992:
						position, tokenIndex = position991, tokenIndex991
						if buffer[position] != rune('T') {
							goto l984
						}
						position++
					}
				l991:
					if !_rules[rulesp]() {
						goto l984
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l984
					}
				l993:
					{
						position994, tokenIndex994 := position, tokenIndex
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
						position, tokenIndex = position994, tokenIndex994
					}
					add(rulePegText, position986)
				}
				if !_rules[ruleAction46]() {
					goto l984
				}
				add(ruleUpdateSourceSinkSpecs, position985)
			}
			return true
		l984:
			position, tokenIndex = position984, tokenIndex984
			return false
		},
		/* 62 SetOptSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action47)> */
		func() bool {
			position995, tokenIndex995 := position, tokenIndex
			{
				position996 := position
				{
					position997 := position
					{
						position998, tokenIndex998 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l998
						}
						{
							position1000, tokenIndex1000 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l1001
							}
							position++
							goto l1000
						l1001:
							position, tokenIndex = position1000, tokenIndex1000
							if buffer[position] != rune('S') {
								goto l998
							}
							position++
						}
					l1000:
						{
							position1002, tokenIndex1002 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1003
							}
							position++
							goto l1002
						l1003:
							position, tokenIndex = position1002, tokenIndex1002
							if buffer[position] != rune('E') {
								goto l998
							}
							position++
						}
					l1002:
						{
							position1004, tokenIndex1004 := position, tokenIndex
							if buffer[position] != rune('t') {
								goto l1005
							}
							position++
							goto l1004
						l1005:
							position, tokenIndex = position1004, tokenIndex1004
							if buffer[position] != rune('T') {
								goto l998
							}
							position++
						}
					l1004:
						if !_rules[rulesp]() {
							goto l998
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l998
						}
					l1006:
						{
							position1007, tokenIndex1007 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1007
							}
							if buffer[position] != rune(',') {
								goto l1007
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1007
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l1007
							}
							goto l1006
						l1007:
							position, tokenIndex = position1007, tokenIndex1007
						}
						goto l999
					l998:
						position, tokenIndex = position998, tokenIndex998
					}
				l999:
					add(rulePegText, position997)
				}
				if !_rules[ruleAction47]() {
					goto l995
				}
				add(ruleSetOptSpecs, position996)
			}
			return true
		l995:
			position, tokenIndex = position995, tokenIndex995
			return false
		},
		/* 63 StateTagOpt <- <(<(sp (('t' / 'T') ('a' / 'A') ('g' / 'G')) sp Identifier)?> Action48)> */
		func() bool {
			position1008, tokenIndex1008 := position, tokenIndex
			{
				position1009 := position
				{
					position1010 := position
					{
						position1011, tokenIndex1011 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1011
						}
						{
							position1013, tokenIndex1013 := position, tokenIndex
							if buffer[position] != rune('t') {
								goto l1014
							}
							position++
							goto l1013
						l1014:
							position, tokenIndex = position1013, tokenIndex1013
							if buffer[position] != rune('T') {
								goto l1011
							}
							position++
						}
					l1013:
						{
							position1015, tokenIndex1015 := position, tokenIndex
							if buffer[position] != rune('a') {
								goto l1016
							}
							position++
							goto l1015
						l1016:
							position, tokenIndex = position1015, tokenIndex1015
							if buffer[position] != rune('A') {
								goto l1011
							}
							position++
						}
					l1015:
						{
							position1017, tokenIndex1017 := position, tokenIndex
							if buffer[position] != rune('g') {
								goto l1018
							}
							position++
							goto l1017
						l1018:
							position, tokenIndex = position1017, tokenIndex1017
							if buffer[position] != rune('G') {
								goto l1011
							}
							position++
						}
					l1017:
						if !_rules[rulesp]() {
							goto l1011
						}
						if !_rules[ruleIdentifier]() {
							goto l1011
						}
						goto l1012
					l1011:
						position, tokenIndex = position1011, tokenIndex1011
					}
				l1012:
					add(rulePegText, position1010)
				}
				if !_rules[ruleAction48]() {
					goto l1008
				}
				add(ruleStateTagOpt, position1009)
			}
			return true
		l1008:
			position, tokenIndex = position1008, tokenIndex1008
			return false
		},
		/* 64 SourceSinkParam <- <(SourceSinkParamKey spOpt '=' spOpt SourceSinkParamVal Action49)> */
		func() bool {
			position1019, tokenIndex1019 := position, tokenIndex
			{
				position1020 := position
				if !_rules[ruleSourceSinkParamKey]() {
					goto l1019
				}
				if !_rules[rulespOpt]() {
					goto l1019
				}
				if buffer[position] != rune('=') {
					goto l1019
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1019
				}
				if !_rules[ruleSourceSinkParamVal]() {
					goto l1019
				}
				if !_rules[ruleAction49]() {
					goto l1019
				}
				add(ruleSourceSinkParam, position1020)
			}
			return true
		l1019:
			position, tokenIndex = position1019, tokenIndex1019
			return false
		},
		/* 65 SourceSinkParamVal <- <ParamLiteral> */
		func() bool {
			position1021, tokenIndex1021 := position, tokenIndex
			{
				position1022 := position
				if !_rules[ruleParamLiteral]() {
					goto l1021
				}
				add(ruleSourceSinkParamVal, position1022)
			}
			return true
		l1021:
			position, tokenIndex = position1021, tokenIndex1021
			return false
		},
		/* 66 ParamLiteral <- <(BooleanLiteral / Literal / ParamArrayExpr / ParamMapExpr)> */
		func() bool {
			position1023, tokenIndex1023 := position, tokenIndex
			{
				position1024 := position
				{
					position1025, tokenIndex1025 := position, tokenIndex
					if !_rules[ruleBooleanLiteral]() {
						goto l1026
					}
					goto l1025
				l1026:
					position, tokenIndex = position1025, tokenIndex1025
					if !_rules[ruleLiteral]() {
						goto l1027
					}
					goto l1025
				l1027:
					position, tokenIndex = position1025, tokenIndex1025
					if !_rules[ruleParamArrayExpr]() {
						goto l1028
					}
					goto l1025
				l1028:
					position, tokenIndex = position1025, tokenIndex1025
					if !_rules[ruleParamMapExpr]() {
						goto l1023
					}
				}
			l1025:
				add(ruleParamLiteral, position1024)
			}
			return true
		l1023:
			position, tokenIndex = position1023, tokenIndex1023
			return false
		},
		/* 67 ParamArrayExpr <- <(<('[' spOpt (ParamLiteral (',' spOpt ParamLiteral)*)? spOpt ','? spOpt ']')> Action50)> */
		func() bool {
			position1029, tokenIndex1029 := position, tokenIndex
			{
				position1030 := position
				{
					position1031 := position
					if buffer[position] != rune('[') {
						goto l1029
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1029
					}
					{
						position1032, tokenIndex1032 := position, tokenIndex
						if !_rules[ruleParamLiteral]() {
							goto l1032
						}
					l1034:
						{
							position1035, tokenIndex1035 := position, tokenIndex
							if buffer[position] != rune(',') {
								goto l1035
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1035
							}
							if !_rules[ruleParamLiteral]() {
								goto l1035
							}
							goto l1034
						l1035:
							position, tokenIndex = position1035, tokenIndex1035
						}
						goto l1033
					l1032:
						position, tokenIndex = position1032, tokenIndex1032
					}
				l1033:
					if !_rules[rulespOpt]() {
						goto l1029
					}
					{
						position1036, tokenIndex1036 := position, tokenIndex
						if buffer[position] != rune(',') {
							goto l1036
						}
						position++
						goto l1037
					l1036:
						position, tokenIndex = position1036, tokenIndex1036
					}
				l1037:
					if !_rules[rulespOpt]() {
						goto l1029
					}
					if buffer[position] != rune(']') {
						goto l1029
					}
					position++
					add(rulePegText, position1031)
				}
				if !_rules[ruleAction50]() {
					goto l1029
				}
				add(ruleParamArrayExpr, position1030)
			}
			return true
		l1029:
			position, tokenIndex = position1029, tokenIndex1029
			return false
		},
		/* 68 ParamMapExpr <- <(<('{' spOpt (ParamKeyValuePair (spOpt ',' spOpt ParamKeyValuePair)*)? spOpt '}')> Action51)> */
		func() bool {
			position1038, tokenIndex1038 := position, tokenIndex
			{
				position1039 := position
				{
					position1040 := position
					if buffer[position] != rune('{') {
						goto l1038
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1038
					}
					{
						position1041, tokenIndex1041 := position, tokenIndex
						if !_rules[ruleParamKeyValuePair]() {
							goto l1041
						}
					l1043:
						{
							position1044, tokenIndex1044 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1044
							}
							if buffer[position] != rune(',') {
								goto l1044
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1044
							}
							if !_rules[ruleParamKeyValuePair]() {
								goto l1044
							}
							goto l1043
						l1044:
							position, tokenIndex = position1044, tokenIndex1044
						}
						goto l1042
					l1041:
						position, tokenIndex = position1041, tokenIndex1041
					}
				l1042:
					if !_rules[rulespOpt]() {
						goto l1038
					}
					if buffer[position] != rune('}') {
						goto l1038
					}
					position++
					add(rulePegText, position1040)
				}
				if !_rules[ruleAction51]() {
					goto l1038
				}
				add(ruleParamMapExpr, position1039)
			}
			return true
		l1038:
			position, tokenIndex = position1038, tokenIndex1038
			return false
		},
		/* 69 ParamKeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ParamLiteral)> Action52)> */
		func() bool {
			position1045, tokenIndex1045 := position, tokenIndex
			{
				position1046 := position
				{
					position1047 := position
					if !_rules[ruleStringLiteral]() {
						goto l1045
					}
					if !_rules[rulespOpt]() {
						goto l1045
					}
					if buffer[position] != rune(':') {
						goto l1045
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1045
					}
					if !_rules[ruleParamLiteral]() {
						goto l1045
					}
					add(rulePegText, position1047)
				}
				if !_rules[ruleAction52]() {
					goto l1045
				}
				add(ruleParamKeyValuePair, position1046)
			}
			return true
		l1045:
			position, tokenIndex = position1045, tokenIndex1045
			return false
		},
		/* 70 PausedOpt <- <(<(sp (Paused / Unpaused))?> Action53)> */
		func() bool {
			position1048, tokenIndex1048 := position, tokenIndex
			{
				position1049 := position
				{
					position1050 := position
					{
						position1051, tokenIndex1051 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1051
						}
						{
							position1053, tokenIndex1053 := position, tokenIndex
							if !_rules[rulePaused]() {
								goto l1054
							}
							goto l1053
						l1054:
							position, tokenIndex = position1053, tokenIndex1053
							if !_rules[ruleUnpaused]() {
								goto l1051
							}
						}
					l1053:
						goto l1052
					l1051:
						position, tokenIndex = position1051, tokenIndex1051
					}
				l1052:
					add(rulePegText, position1050)
				}
				if !_rules[ruleAction53]() {
					goto l1048
				}
				add(rulePausedOpt, position1049)
			}
			return true
		l1048:
			position, tokenIndex = position1048, tokenIndex1048
			return false
		},
		/* 71 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position1055, tokenIndex1055 := position, tokenIndex
			{
				position1056 := position
				{
					position1057, tokenIndex1057 := position, tokenIndex
					if !_rules[ruleWildcard]() {
						goto l1058
					}
					goto l1057
				l1058:
					position, tokenIndex = position1057, tokenIndex1057
					if !_rules[ruleExpression]() {
						goto l1055
					}
				}
			l1057:
				add(ruleExpressionOrWildcard, position1056)
			}
			return true
		l1055:
			position, tokenIndex = position1055, tokenIndex1055
			return false
		},
		/* 72 Expression <- <orExpr> */
		func() bool {
			position1059, tokenIndex1059 := position, tokenIndex
			{
				position1060 := position
				if !_rules[ruleorExpr]() {
					goto l1059
				}
				add(ruleExpression, position1060)
			}
			return true
		l1059:
			position, tokenIndex = position1059, tokenIndex1059
			return false
		},
		/* 73 orExpr <- <(<(andExpr (sp Or sp andExpr)*)> Action54)> */
		func() bool {
			position1061, tokenIndex1061 := position, tokenIndex
			{
				position1062 := position
				{
					position1063 := position
					if !_rules[ruleandExpr]() {
						goto l1061
					}
				l1064:
					{
						position1065, tokenIndex1065 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1065
						}
						if !_rules[ruleOr]() {
							goto l1065
						}
						if !_rules[rulesp]() {
							goto l1065
						}
						if !_rules[ruleandExpr]() {
							goto l1065
						}
						goto l1064
					l1065:
						position, tokenIndex = position1065, tokenIndex1065
					}
					add(rulePegText, position1063)
				}
				if !_rules[ruleAction54]() {
					goto l1061
				}
				add(ruleorExpr, position1062)
			}
			return true
		l1061:
			position, tokenIndex = position1061, tokenIndex1061
			return false
		},
		/* 74 andExpr <- <(<(notExpr (sp And sp notExpr)*)> Action55)> */
		func() bool {
			position1066, tokenIndex1066 := position, tokenIndex
			{
				position1067 := position
				{
					position1068 := position
					if !_rules[rulenotExpr]() {
						goto l1066
					}
				l1069:
					{
						position1070, tokenIndex1070 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1070
						}
						if !_rules[ruleAnd]() {
							goto l1070
						}
						if !_rules[rulesp]() {
							goto l1070
						}
						if !_rules[rulenotExpr]() {
							goto l1070
						}
						goto l1069
					l1070:
						position, tokenIndex = position1070, tokenIndex1070
					}
					add(rulePegText, position1068)
				}
				if !_rules[ruleAction55]() {
					goto l1066
				}
				add(ruleandExpr, position1067)
			}
			return true
		l1066:
			position, tokenIndex = position1066, tokenIndex1066
			return false
		},
		/* 75 notExpr <- <(<((Not sp)? comparisonExpr)> Action56)> */
		func() bool {
			position1071, tokenIndex1071 := position, tokenIndex
			{
				position1072 := position
				{
					position1073 := position
					{
						position1074, tokenIndex1074 := position, tokenIndex
						if !_rules[ruleNot]() {
							goto l1074
						}
						if !_rules[rulesp]() {
							goto l1074
						}
						goto l1075
					l1074:
						position, tokenIndex = position1074, tokenIndex1074
					}
				l1075:
					if !_rules[rulecomparisonExpr]() {
						goto l1071
					}
					add(rulePegText, position1073)
				}
				if !_rules[ruleAction56]() {
					goto l1071
				}
				add(rulenotExpr, position1072)
			}
			return true
		l1071:
			position, tokenIndex = position1071, tokenIndex1071
			return false
		},
		/* 76 comparisonExpr <- <(<(otherOpExpr (spOpt ComparisonOp spOpt otherOpExpr)?)> Action57)> */
		func() bool {
			position1076, tokenIndex1076 := position, tokenIndex
			{
				position1077 := position
				{
					position1078 := position
					if !_rules[ruleotherOpExpr]() {
						goto l1076
					}
					{
						position1079, tokenIndex1079 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1079
						}
						if !_rules[ruleComparisonOp]() {
							goto l1079
						}
						if !_rules[rulespOpt]() {
							goto l1079
						}
						if !_rules[ruleotherOpExpr]() {
							goto l1079
						}
						goto l1080
					l1079:
						position, tokenIndex = position1079, tokenIndex1079
					}
				l1080:
					add(rulePegText, position1078)
				}
				if !_rules[ruleAction57]() {
					goto l1076
				}
				add(rulecomparisonExpr, position1077)
			}
			return true
		l1076:
			position, tokenIndex = position1076, tokenIndex1076
			return false
		},
		/* 77 otherOpExpr <- <(<(isExpr (spOpt OtherOp spOpt isExpr)*)> Action58)> */
		func() bool {
			position1081, tokenIndex1081 := position, tokenIndex
			{
				position1082 := position
				{
					position1083 := position
					if !_rules[ruleisExpr]() {
						goto l1081
					}
				l1084:
					{
						position1085, tokenIndex1085 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1085
						}
						if !_rules[ruleOtherOp]() {
							goto l1085
						}
						if !_rules[rulespOpt]() {
							goto l1085
						}
						if !_rules[ruleisExpr]() {
							goto l1085
						}
						goto l1084
					l1085:
						position, tokenIndex = position1085, tokenIndex1085
					}
					add(rulePegText, position1083)
				}
				if !_rules[ruleAction58]() {
					goto l1081
				}
				add(ruleotherOpExpr, position1082)
			}
			return true
		l1081:
			position, tokenIndex = position1081, tokenIndex1081
			return false
		},
		/* 78 isExpr <- <(<((RowValue sp IsOp sp Missing) / (termExpr (sp IsOp sp NullLiteral)?))> Action59)> */
		func() bool {
			position1086, tokenIndex1086 := position, tokenIndex
			{
				position1087 := position
				{
					position1088 := position
					{
						position1089, tokenIndex1089 := position, tokenIndex
						if !_rules[ruleRowValue]() {
							goto l1090
						}
						if !_rules[rulesp]() {
							goto l1090
						}
						if !_rules[ruleIsOp]() {
							goto l1090
						}
						if !_rules[rulesp]() {
							goto l1090
						}
						if !_rules[ruleMissing]() {
							goto l1090
						}
						goto l1089
					l1090:
						position, tokenIndex = position1089, tokenIndex1089
						if !_rules[ruletermExpr]() {
							goto l1086
						}
						{
							position1091, tokenIndex1091 := position, tokenIndex
							if !_rules[rulesp]() {
								goto l1091
							}
							if !_rules[ruleIsOp]() {
								goto l1091
							}
							if !_rules[rulesp]() {
								goto l1091
							}
							if !_rules[ruleNullLiteral]() {
								goto l1091
							}
							goto l1092
						l1091:
							position, tokenIndex = position1091, tokenIndex1091
						}
					l1092:
					}
				l1089:
					add(rulePegText, position1088)
				}
				if !_rules[ruleAction59]() {
					goto l1086
				}
				add(ruleisExpr, position1087)
			}
			return true
		l1086:
			position, tokenIndex = position1086, tokenIndex1086
			return false
		},
		/* 79 termExpr <- <(<(productExpr (spOpt PlusMinusOp spOpt productExpr)*)> Action60)> */
		func() bool {
			position1093, tokenIndex1093 := position, tokenIndex
			{
				position1094 := position
				{
					position1095 := position
					if !_rules[ruleproductExpr]() {
						goto l1093
					}
				l1096:
					{
						position1097, tokenIndex1097 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1097
						}
						if !_rules[rulePlusMinusOp]() {
							goto l1097
						}
						if !_rules[rulespOpt]() {
							goto l1097
						}
						if !_rules[ruleproductExpr]() {
							goto l1097
						}
						goto l1096
					l1097:
						position, tokenIndex = position1097, tokenIndex1097
					}
					add(rulePegText, position1095)
				}
				if !_rules[ruleAction60]() {
					goto l1093
				}
				add(ruletermExpr, position1094)
			}
			return true
		l1093:
			position, tokenIndex = position1093, tokenIndex1093
			return false
		},
		/* 80 productExpr <- <(<(minusExpr (spOpt MultDivOp spOpt minusExpr)*)> Action61)> */
		func() bool {
			position1098, tokenIndex1098 := position, tokenIndex
			{
				position1099 := position
				{
					position1100 := position
					if !_rules[ruleminusExpr]() {
						goto l1098
					}
				l1101:
					{
						position1102, tokenIndex1102 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1102
						}
						if !_rules[ruleMultDivOp]() {
							goto l1102
						}
						if !_rules[rulespOpt]() {
							goto l1102
						}
						if !_rules[ruleminusExpr]() {
							goto l1102
						}
						goto l1101
					l1102:
						position, tokenIndex = position1102, tokenIndex1102
					}
					add(rulePegText, position1100)
				}
				if !_rules[ruleAction61]() {
					goto l1098
				}
				add(ruleproductExpr, position1099)
			}
			return true
		l1098:
			position, tokenIndex = position1098, tokenIndex1098
			return false
		},
		/* 81 minusExpr <- <(<((UnaryMinus spOpt)? castExpr)> Action62)> */
		func() bool {
			position1103, tokenIndex1103 := position, tokenIndex
			{
				position1104 := position
				{
					position1105 := position
					{
						position1106, tokenIndex1106 := position, tokenIndex
						if !_rules[ruleUnaryMinus]() {
							goto l1106
						}
						if !_rules[rulespOpt]() {
							goto l1106
						}
						goto l1107
					l1106:
						position, tokenIndex = position1106, tokenIndex1106
					}
				l1107:
					if !_rules[rulecastExpr]() {
						goto l1103
					}
					add(rulePegText, position1105)
				}
				if !_rules[ruleAction62]() {
					goto l1103
				}
				add(ruleminusExpr, position1104)
			}
			return true
		l1103:
			position, tokenIndex = position1103, tokenIndex1103
			return false
		},
		/* 82 castExpr <- <(<(baseExpr (spOpt (':' ':') spOpt Type)?)> Action63)> */
		func() bool {
			position1108, tokenIndex1108 := position, tokenIndex
			{
				position1109 := position
				{
					position1110 := position
					if !_rules[rulebaseExpr]() {
						goto l1108
					}
					{
						position1111, tokenIndex1111 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1111
						}
						if buffer[position] != rune(':') {
							goto l1111
						}
						position++
						if buffer[position] != rune(':') {
							goto l1111
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1111
						}
						if !_rules[ruleType]() {
							goto l1111
						}
						goto l1112
					l1111:
						position, tokenIndex = position1111, tokenIndex1111
					}
				l1112:
					add(rulePegText, position1110)
				}
				if !_rules[ruleAction63]() {
					goto l1108
				}
				add(rulecastExpr, position1109)
			}
			return true
		l1108:
			position, tokenIndex = position1108, tokenIndex1108
			return false
		},
		/* 83 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / Case / RowMeta / FuncTypeCast / FuncAppSelector / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position1113, tokenIndex1113 := position, tokenIndex
			{
				position1114 := position
				{
					position1115, tokenIndex1115 := position, tokenIndex
					if buffer[position] != rune('(') {
						goto l1116
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1116
					}
					if !_rules[ruleExpression]() {
						goto l1116
					}
					if !_rules[rulespOpt]() {
						goto l1116
					}
					if buffer[position] != rune(')') {
						goto l1116
					}
					position++
					goto l1115
				l1116:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleMapExpr]() {
						goto l1117
					}
					goto l1115
				l1117:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleBooleanLiteral]() {
						goto l1118
					}
					goto l1115
				l1118:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleNullLiteral]() {
						goto l1119
					}
					goto l1115
				l1119:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleCase]() {
						goto l1120
					}
					goto l1115
				l1120:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleRowMeta]() {
						goto l1121
					}
					goto l1115
				l1121:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleFuncTypeCast]() {
						goto l1122
					}
					goto l1115
				l1122:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleFuncAppSelector]() {
						goto l1123
					}
					goto l1115
				l1123:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleFuncApp]() {
						goto l1124
					}
					goto l1115
				l1124:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleRowValue]() {
						goto l1125
					}
					goto l1115
				l1125:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleArrayExpr]() {
						goto l1126
					}
					goto l1115
				l1126:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleLiteral]() {
						goto l1113
					}
				}
			l1115:
				add(rulebaseExpr, position1114)
			}
			return true
		l1113:
			position, tokenIndex = position1113, tokenIndex1113
			return false
		},
		/* 84 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') spOpt '(' spOpt Expression sp (('a' / 'A') ('s' / 'S')) sp Type spOpt ')')> Action64)> */
		func() bool {
			position1127, tokenIndex1127 := position, tokenIndex
			{
				position1128 := position
				{
					position1129 := position
					{
						position1130, tokenIndex1130 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1131
						}
						position++
						goto l1130
					l1131:
						position, tokenIndex = position1130, tokenIndex1130
						if buffer[position] != rune('C') {
							goto l1127
						}
						position++
					}
				l1130:
					{
						position1132, tokenIndex1132 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1133
						}
						position++
						goto l1132
					l1133:
						position, tokenIndex = position1132, tokenIndex1132
						if buffer[position] != rune('A') {
							goto l1127
						}
						position++
					}
				l1132:
					{
						position1134, tokenIndex1134 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1135
						}
						position++
						goto l1134
					l1135:
						position, tokenIndex = position1134, tokenIndex1134
						if buffer[position] != rune('S') {
							goto l1127
						}
						position++
					}
				l1134:
					{
						position1136, tokenIndex1136 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1137
						}
						position++
						goto l1136
					l1137:
						position, tokenIndex = position1136, tokenIndex1136
						if buffer[position] != rune('T') {
							goto l1127
						}
						position++
					}
				l1136:
					if !_rules[rulespOpt]() {
						goto l1127
					}
					if buffer[position] != rune('(') {
						goto l1127
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1127
					}
					if !_rules[ruleExpression]() {
						goto l1127
					}
					if !_rules[rulesp]() {
						goto l1127
					}
					{
						position1138, tokenIndex1138 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1139
						}
						position++
						goto l1138
					l1139:
						position, tokenIndex = position1138, tokenIndex1138
						if buffer[position] != rune('A') {
							goto l1127
						}
						position++
					}
				l1138:
					{
						position1140, tokenIndex1140 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1141
						}
						position++
						goto l1140
					l1141:
						position, tokenIndex = position1140, tokenIndex1140
						if buffer[position] != rune('S') {
							goto l1127
						}
						position++
					}
				l1140:
					if !_rules[rulesp]() {
						goto l1127
					}
					if !_rules[ruleType]() {
						goto l1127
					}
					if !_rules[rulespOpt]() {
						goto l1127
					}
					if buffer[position] != rune(')') {
						goto l1127
					}
					position++
					add(rulePegText, position1129)
				}
				if !_rules[ruleAction64]() {
					goto l1127
				}
				add(ruleFuncTypeCast, position1128)
			}
			return true
		l1127:
			position, tokenIndex = position1127, tokenIndex1127
			return false
		},
		/* 85 FuncApp <- <(FuncAppWithOrderBy / FuncAppWithoutOrderBy)> */
		func() bool {
			position1142, tokenIndex1142 := position, tokenIndex
			{
				position1143 := position
				{
					position1144, tokenIndex1144 := position, tokenIndex
					if !_rules[ruleFuncAppWithOrderBy]() {
						goto l1145
					}
					goto l1144
				l1145:
					position, tokenIndex = position1144, tokenIndex1144
					if !_rules[ruleFuncAppWithoutOrderBy]() {
						goto l1142
					}
				}
			l1144:
				add(ruleFuncApp, position1143)
			}
			return true
		l1142:
			position, tokenIndex = position1142, tokenIndex1142
			return false
		},
		/* 86 FuncAppSelector <- <(FuncApp FuncElemAccessor Action65)> */
		func() bool {
			position1146, tokenIndex1146 := position, tokenIndex
			{
				position1147 := position
				if !_rules[ruleFuncApp]() {
					goto l1146
				}
				if !_rules[ruleFuncElemAccessor]() {
					goto l1146
				}
				if !_rules[ruleAction65]() {
					goto l1146
				}
				add(ruleFuncAppSelector, position1147)
			}
			return true
		l1146:
			position, tokenIndex = position1146, tokenIndex1146
			return false
		},
		/* 87 FuncElemAccessor <- <(<jsonGetPathNonHead+> Action66)> */
		func() bool {
			position1148, tokenIndex1148 := position, tokenIndex
			{
				position1149 := position
				{
					position1150 := position
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1148
					}
				l1151:
					{
						position1152, tokenIndex1152 := position, tokenIndex
						if !_rules[rulejsonGetPathNonHead]() {
							goto l1152
						}
						goto l1151
					l1152:
						position, tokenIndex = position1152, tokenIndex1152
					}
					add(rulePegText, position1150)
				}
				if !_rules[ruleAction66]() {
					goto l1148
				}
				add(ruleFuncElemAccessor, position1149)
			}
			return true
		l1148:
			position, tokenIndex = position1148, tokenIndex1148
			return false
		},
		/* 88 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action67)> */
		func() bool {
			position1153, tokenIndex1153 := position, tokenIndex
			{
				position1154 := position
				if !_rules[ruleFunction]() {
					goto l1153
				}
				if !_rules[rulespOpt]() {
					goto l1153
				}
				if buffer[position] != rune('(') {
					goto l1153
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1153
				}
				if !_rules[ruleFuncParams]() {
					goto l1153
				}
				if !_rules[rulesp]() {
					goto l1153
				}
				if !_rules[ruleParamsOrder]() {
					goto l1153
				}
				if !_rules[rulespOpt]() {
					goto l1153
				}
				if buffer[position] != rune(')') {
					goto l1153
				}
				position++
				if !_rules[ruleAction67]() {
					goto l1153
				}
				add(ruleFuncAppWithOrderBy, position1154)
			}
			return true
		l1153:
			position, tokenIndex = position1153, tokenIndex1153
			return false
		},
		/* 89 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action68)> */
		func() bool {
			position1155, tokenIndex1155 := position, tokenIndex
			{
				position1156 := position
				if !_rules[ruleFunction]() {
					goto l1155
				}
				if !_rules[rulespOpt]() {
					goto l1155
				}
				if buffer[position] != rune('(') {
					goto l1155
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1155
				}
				if !_rules[ruleFuncParams]() {
					goto l1155
				}
				{
					position1157 := position
					if !_rules[rulespOpt]() {
						goto l1155
					}
					add(rulePegText, position1157)
				}
				if buffer[position] != rune(')') {
					goto l1155
				}
				position++
				if !_rules[ruleAction68]() {
					goto l1155
				}
				add(ruleFuncAppWithoutOrderBy, position1156)
			}
			return true
		l1155:
			position, tokenIndex = position1155, tokenIndex1155
			return false
		},
		/* 90 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action69)> */
		func() bool {
			position1158, tokenIndex1158 := position, tokenIndex
			{
				position1159 := position
				{
					position1160 := position
					{
						position1161, tokenIndex1161 := position, tokenIndex
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1161
						}
					l1163:
						{
							position1164, tokenIndex1164 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1164
							}
							if buffer[position] != rune(',') {
								goto l1164
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1164
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1164
							}
							goto l1163
						l1164:
							position, tokenIndex = position1164, tokenIndex1164
						}
						goto l1162
					l1161:
						position, tokenIndex = position1161, tokenIndex1161
					}
				l1162:
					add(rulePegText, position1160)
				}
				if !_rules[ruleAction69]() {
					goto l1158
				}
				add(ruleFuncParams, position1159)
			}
			return true
		l1158:
			position, tokenIndex = position1158, tokenIndex1158
			return false
		},
		/* 91 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action70)> */
		func() bool {
			position1165, tokenIndex1165 := position, tokenIndex
			{
				position1166 := position
				{
					position1167 := position
					{
						position1168, tokenIndex1168 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1169
						}
						position++
						goto l1168
					l1169:
						position, tokenIndex = position1168, tokenIndex1168
						if buffer[position] != rune('O') {
							goto l1165
						}
						position++
					}
				l1168:
					{
						position1170, tokenIndex1170 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1171
						}
						position++
						goto l1170
					l1171:
						position, tokenIndex = position1170, tokenIndex1170
						if buffer[position] != rune('R') {
							goto l1165
						}
						position++
					}
				l1170:
					{
						position1172, tokenIndex1172 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1173
						}
						position++
						goto l1172
					l1173:
						position, tokenIndex = position1172, tokenIndex1172
						if buffer[position] != rune('D') {
							goto l1165
						}
						position++
					}
				l1172:
					{
						position1174, tokenIndex1174 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1175
						}
						position++
						goto l1174
					l1175:
						position, tokenIndex = position1174, tokenIndex1174
						if buffer[position] != rune('E') {
							goto l1165
						}
						position++
					}
				l1174:
					{
						position1176, tokenIndex1176 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1177
						}
						position++
						goto l1176
					l1177:
						position, tokenIndex = position1176, tokenIndex1176
						if buffer[position] != rune('R') {
							goto l1165
						}
						position++
					}
				l1176:
					if !_rules[rulesp]() {
						goto l1165
					}
					{
						position1178, tokenIndex1178 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1179
						}
						position++
						goto l1178
					l1179:
						position, tokenIndex = position1178, tokenIndex1178
						if buffer[position] != rune('B') {
							goto l1165
						}
						position++
					}
				l1178:
					{
						position1180, tokenIndex1180 := position, tokenIndex
						if buffer[position] != rune('y') {
							goto l1181
						}
						position++
						goto l1180
					l1181:
						position, tokenIndex = position1180, tokenIndex1180
						if buffer[position] != rune('Y') {
							goto l1165
						}
						position++
					}
				l1180:
					if !_rules[rulesp]() {
						goto l1165
					}
					if !_rules[ruleSortedExpression]() {
						goto l1165
					}
				l1182:
					{
						position1183, tokenIndex1183 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1183
						}
						if buffer[position] != rune(',') {
							goto l1183
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1183
						}
						if !_rules[ruleSortedExpression]() {
							goto l1183
						}
						goto l1182
					l1183:
						position, tokenIndex = position1183, tokenIndex1183
					}
					add(rulePegText, position1167)
				}
				if !_rules[ruleAction70]() {
					goto l1165
				}
				add(ruleParamsOrder, position1166)
			}
			return true
		l1165:
			position, tokenIndex = position1165, tokenIndex1165
			return false
		},
		/* 92 SortedExpression <- <(Expression OrderDirectionOpt Action71)> */
		func() bool {
			position1184, tokenIndex1184 := position, tokenIndex
			{
				position1185 := position
				if !_rules[ruleExpression]() {
					goto l1184
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1184
				}
				if !_rules[ruleAction71]() {
					goto l1184
				}
				add(ruleSortedExpression, position1185)
			}
			return true
		l1184:
			position, tokenIndex = position1184, tokenIndex1184
			return false
		},
		/* 93 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action72)> */
		func() bool {
			position1186, tokenIndex1186 := position, tokenIndex
			{
				position1187 := position
				{
					position1188 := position
					{
						position1189, tokenIndex1189 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1189
						}
						{
							position1191, tokenIndex1191 := position, tokenIndex
							if !_rules[ruleAscending]() {
								goto l1192
							}
							goto l1191
						l1192:
							position, tokenIndex = position1191, tokenIndex1191
							if !_rules[ruleDescending]() {
								goto l1189
							}
						}
					l1191:
						goto l1190
					l1189:
						position, tokenIndex = position1189, tokenIndex1189
					}
				l1190:
					add(rulePegText, position1188)
				}
				if !_rules[ruleAction72]() {
					goto l1186
				}
				add(ruleOrderDirectionOpt, position1187)
			}
			return true
		l1186:
			position, tokenIndex = position1186, tokenIndex1186
			return false
		},
		/* 94 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action73)> */
		func() bool {
			position1193, tokenIndex1193 := position, tokenIndex
			{
				position1194 := position
				{
					position1195 := position
					if buffer[position] != rune('[') {
						goto l1193
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1193
					}
					{
						position1196, tokenIndex1196 := position, tokenIndex
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1196
						}
					l1198:
						{
							position1199, tokenIndex1199 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1199
							}
							if buffer[position] != rune(',') {
								goto l1199
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1199
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1199
							}
							goto l1198
						l1199:
							position, tokenIndex = position1199, tokenIndex1199
						}
						goto l1197
					l1196:
						position, tokenIndex = position1196, tokenIndex1196
					}
				l1197:
					if !_rules[rulespOpt]() {
						goto l1193
					}
					{
						position1200, tokenIndex1200 := position, tokenIndex
						if buffer[position] != rune(',') {
							goto l1200
						}
						position++
						goto l1201
					l1200:
						position, tokenIndex = position1200, tokenIndex1200
					}
				l1201:
					if !_rules[rulespOpt]() {
						goto l1193
					}
					if buffer[position] != rune(']') {
						goto l1193
					}
					position++
					add(rulePegText, position1195)
				}
				if !_rules[ruleAction73]() {
					goto l1193
				}
				add(ruleArrayExpr, position1194)
			}
			return true
		l1193:
			position, tokenIndex = position1193, tokenIndex1193
			return false
		},
		/* 95 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action74)> */
		func() bool {
			position1202, tokenIndex1202 := position, tokenIndex
			{
				position1203 := position
				{
					position1204 := position
					if buffer[position] != rune('{') {
						goto l1202
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1202
					}
					{
						position1205, tokenIndex1205 := position, tokenIndex
						if !_rules[ruleKeyValuePair]() {
							goto l1205
						}
					l1207:
						{
							position1208, tokenIndex1208 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1208
							}
							if buffer[position] != rune(',') {
								goto l1208
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1208
							}
							if !_rules[ruleKeyValuePair]() {
								goto l1208
							}
							goto l1207
						l1208:
							position, tokenIndex = position1208, tokenIndex1208
						}
						goto l1206
					l1205:
						position, tokenIndex = position1205, tokenIndex1205
					}
				l1206:
					if !_rules[rulespOpt]() {
						goto l1202
					}
					if buffer[position] != rune('}') {
						goto l1202
					}
					position++
					add(rulePegText, position1204)
				}
				if !_rules[ruleAction74]() {
					goto l1202
				}
				add(ruleMapExpr, position1203)
			}
			return true
		l1202:
			position, tokenIndex = position1202, tokenIndex1202
			return false
		},
		/* 96 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action75)> */
		func() bool {
			position1209, tokenIndex1209 := position, tokenIndex
			{
				position1210 := position
				{
					position1211 := position
					if !_rules[ruleStringLiteral]() {
						goto l1209
					}
					if !_rules[rulespOpt]() {
						goto l1209
					}
					if buffer[position] != rune(':') {
						goto l1209
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1209
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1209
					}
					add(rulePegText, position1211)
				}
				if !_rules[ruleAction75]() {
					goto l1209
				}
				add(ruleKeyValuePair, position1210)
			}
			return true
		l1209:
			position, tokenIndex = position1209, tokenIndex1209
			return false
		},
		/* 97 Case <- <(ConditionCase / ExpressionCase)> */
		func() bool {
			position1212, tokenIndex1212 := position, tokenIndex
			{
				position1213 := position
				{
					position1214, tokenIndex1214 := position, tokenIndex
					if !_rules[ruleConditionCase]() {
						goto l1215
					}
					goto l1214
				l1215:
					position, tokenIndex = position1214, tokenIndex1214
					if !_rules[ruleExpressionCase]() {
						goto l1212
					}
				}
			l1214:
				add(ruleCase, position1213)
			}
			return true
		l1212:
			position, tokenIndex = position1212, tokenIndex1212
			return false
		},
		/* 98 ConditionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action76)> */
		func() bool {
			position1216, tokenIndex1216 := position, tokenIndex
			{
				position1217 := position
				{
					position1218, tokenIndex1218 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l1219
					}
					position++
					goto l1218
				l1219:
					position, tokenIndex = position1218, tokenIndex1218
					if buffer[position] != rune('C') {
						goto l1216
					}
					position++
				}
			l1218:
				{
					position1220, tokenIndex1220 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l1221
					}
					position++
					goto l1220
				l1221:
					position, tokenIndex = position1220, tokenIndex1220
					if buffer[position] != rune('A') {
						goto l1216
					}
					position++
				}
			l1220:
				{
					position1222, tokenIndex1222 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l1223
					}
					position++
					goto l1222
				l1223:
					position, tokenIndex = position1222, tokenIndex1222
					if buffer[position] != rune('S') {
						goto l1216
					}
					position++
				}
			l1222:
				{
					position1224, tokenIndex1224 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1225
					}
					position++
					goto l1224
				l1225:
					position, tokenIndex = position1224, tokenIndex1224
					if buffer[position] != rune('E') {
						goto l1216
					}
					position++
				}
			l1224:
				{
					position1226 := position
					if !_rules[rulesp]() {
						goto l1216
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1216
					}
				l1227:
					{
						position1228, tokenIndex1228 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1228
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1228
						}
						goto l1227
					l1228:
						position, tokenIndex = position1228, tokenIndex1228
					}
					{
						position1229, tokenIndex1229 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1229
						}
						{
							position1231, tokenIndex1231 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1232
							}
							position++
							goto l1231
						l1232:
							position, tokenIndex = position1231, tokenIndex1231
							if buffer[position] != rune('E') {
								goto l1229
							}
							position++
						}
					l1231:
						{
							position1233, tokenIndex1233 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l1234
							}
							position++
							goto l1233
						l1234:
							position, tokenIndex = position1233, tokenIndex1233
							if buffer[position] != rune('L') {
								goto l1229
							}
							position++
						}
					l1233:
						{
							position1235, tokenIndex1235 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l1236
							}
							position++
							goto l1235
						l1236:
							position, tokenIndex = position1235, tokenIndex1235
							if buffer[position] != rune('S') {
								goto l1229
							}
							position++
						}
					l1235:
						{
							position1237, tokenIndex1237 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1238
							}
							position++
							goto l1237
						l1238:
							position, tokenIndex = position1237, tokenIndex1237
							if buffer[position] != rune('E') {
								goto l1229
							}
							position++
						}
					l1237:
						if !_rules[rulesp]() {
							goto l1229
						}
						if !_rules[ruleExpression]() {
							goto l1229
						}
						goto l1230
					l1229:
						position, tokenIndex = position1229, tokenIndex1229
					}
				l1230:
					if !_rules[rulesp]() {
						goto l1216
					}
					{
						position1239, tokenIndex1239 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1240
						}
						position++
						goto l1239
					l1240:
						position, tokenIndex = position1239, tokenIndex1239
						if buffer[position] != rune('E') {
							goto l1216
						}
						position++
					}
				l1239:
					{
						position1241, tokenIndex1241 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1242
						}
						position++
						goto l1241
					l1242:
						position, tokenIndex = position1241, tokenIndex1241
						if buffer[position] != rune('N') {
							goto l1216
						}
						position++
					}
				l1241:
					{
						position1243, tokenIndex1243 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1244
						}
						position++
						goto l1243
					l1244:
						position, tokenIndex = position1243, tokenIndex1243
						if buffer[position] != rune('D') {
							goto l1216
						}
						position++
					}
				l1243:
					add(rulePegText, position1226)
				}
				if !_rules[ruleAction76]() {
					goto l1216
				}
				add(ruleConditionCase, position1217)
			}
			return true
		l1216:
			position, tokenIndex = position1216, tokenIndex1216
			return false
		},
		/* 99 ExpressionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') sp Expression <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action77)> */
		func() bool {
			position1245, tokenIndex1245 := position, tokenIndex
			{
				position1246 := position
				{
					position1247, tokenIndex1247 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l1248
					}
					position++
					goto l1247
				l1248:
					position, tokenIndex = position1247, tokenIndex1247
					if buffer[position] != rune('C') {
						goto l1245
					}
					position++
				}
			l1247:
				{
					position1249, tokenIndex1249 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l1250
					}
					position++
					goto l1249
				l1250:
					position, tokenIndex = position1249, tokenIndex1249
					if buffer[position] != rune('A') {
						goto l1245
					}
					position++
				}
			l1249:
				{
					position1251, tokenIndex1251 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l1252
					}
					position++
					goto l1251
				l1252:
					position, tokenIndex = position1251, tokenIndex1251
					if buffer[position] != rune('S') {
						goto l1245
					}
					position++
				}
			l1251:
				{
					position1253, tokenIndex1253 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1254
					}
					position++
					goto l1253
				l1254:
					position, tokenIndex = position1253, tokenIndex1253
					if buffer[position] != rune('E') {
						goto l1245
					}
					position++
				}
			l1253:
				if !_rules[rulesp]() {
					goto l1245
				}
				if !_rules[ruleExpression]() {
					goto l1245
				}
				{
					position1255 := position
					if !_rules[rulesp]() {
						goto l1245
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1245
					}
				l1256:
					{
						position1257, tokenIndex1257 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1257
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1257
						}
						goto l1256
					l1257:
						position, tokenIndex = position1257, tokenIndex1257
					}
					{
						position1258, tokenIndex1258 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1258
						}
						{
							position1260, tokenIndex1260 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1261
							}
							position++
							goto l1260
						l1261:
							position, tokenIndex = position1260, tokenIndex1260
							if buffer[position] != rune('E') {
								goto l1258
							}
							position++
						}
					l1260:
						{
							position1262, tokenIndex1262 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l1263
							}
							position++
							goto l1262
						l1263:
							position, tokenIndex = position1262, tokenIndex1262
							if buffer[position] != rune('L') {
								goto l1258
							}
							position++
						}
					l1262:
						{
							position1264, tokenIndex1264 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l1265
							}
							position++
							goto l1264
						l1265:
							position, tokenIndex = position1264, tokenIndex1264
							if buffer[position] != rune('S') {
								goto l1258
							}
							position++
						}
					l1264:
						{
							position1266, tokenIndex1266 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1267
							}
							position++
							goto l1266
						l1267:
							position, tokenIndex = position1266, tokenIndex1266
							if buffer[position] != rune('E') {
								goto l1258
							}
							position++
						}
					l1266:
						if !_rules[rulesp]() {
							goto l1258
						}
						if !_rules[ruleExpression]() {
							goto l1258
						}
						goto l1259
					l1258:
						position, tokenIndex = position1258, tokenIndex1258
					}
				l1259:
					if !_rules[rulesp]() {
						goto l1245
					}
					{
						position1268, tokenIndex1268 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1269
						}
						position++
						goto l1268
					l1269:
						position, tokenIndex = position1268, tokenIndex1268
						if buffer[position] != rune('E') {
							goto l1245
						}
						position++
					}
				l1268:
					{
						position1270, tokenIndex1270 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1271
						}
						position++
						goto l1270
					l1271:
						position, tokenIndex = position1270, tokenIndex1270
						if buffer[position] != rune('N') {
							goto l1245
						}
						position++
					}
				l1270:
					{
						position1272, tokenIndex1272 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1273
						}
						position++
						goto l1272
					l1273:
						position, tokenIndex = position1272, tokenIndex1272
						if buffer[position] != rune('D') {
							goto l1245
						}
						position++
					}
				l1272:
					add(rulePegText, position1255)
				}
				if !_rules[ruleAction77]() {
					goto l1245
				}
				add(ruleExpressionCase, position1246)
			}
			return true
		l1245:
			position, tokenIndex = position1245, tokenIndex1245
			return false
		},
		/* 100 WhenThenPair <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('n' / 'N') sp Expression sp (('t' / 'T') ('h' / 'H') ('e' / 'E') ('n' / 'N')) sp ExpressionOrWildcard Action78)> */
		func() bool {
			position1274, tokenIndex1274 := position, tokenIndex
			{
				position1275 := position
				{
					position1276, tokenIndex1276 := position, tokenIndex
					if buffer[position] != rune('w') {
						goto l1277
					}
					position++
					goto l1276
				l1277:
					position, tokenIndex = position1276, tokenIndex1276
					if buffer[position] != rune('W') {
						goto l1274
					}
					position++
				}
			l1276:
				{
					position1278, tokenIndex1278 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l1279
					}
					position++
					goto l1278
				l1279:
					position, tokenIndex = position1278, tokenIndex1278
					if buffer[position] != rune('H') {
						goto l1274
					}
					position++
				}
			l1278:
				{
					position1280, tokenIndex1280 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1281
					}
					position++
					goto l1280
				l1281:
					position, tokenIndex = position1280, tokenIndex1280
					if buffer[position] != rune('E') {
						goto l1274
					}
					position++
				}
			l1280:
				{
					position1282, tokenIndex1282 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l1283
					}
					position++
					goto l1282
				l1283:
					position, tokenIndex = position1282, tokenIndex1282
					if buffer[position] != rune('N') {
						goto l1274
					}
					position++
				}
			l1282:
				if !_rules[rulesp]() {
					goto l1274
				}
				if !_rules[ruleExpression]() {
					goto l1274
				}
				if !_rules[rulesp]() {
					goto l1274
				}
				{
					position1284, tokenIndex1284 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l1285
					}
					position++
					goto l1284
				l1285:
					position, tokenIndex = position1284, tokenIndex1284
					if buffer[position] != rune('T') {
						goto l1274
					}
					position++
				}
			l1284:
				{
					position1286, tokenIndex1286 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l1287
					}
					position++
					goto l1286
				l1287:
					position, tokenIndex = position1286, tokenIndex1286
					if buffer[position] != rune('H') {
						goto l1274
					}
					position++
				}
			l1286:
				{
					position1288, tokenIndex1288 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1289
					}
					position++
					goto l1288
				l1289:
					position, tokenIndex = position1288, tokenIndex1288
					if buffer[position] != rune('E') {
						goto l1274
					}
					position++
				}
			l1288:
				{
					position1290, tokenIndex1290 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l1291
					}
					position++
					goto l1290
				l1291:
					position, tokenIndex = position1290, tokenIndex1290
					if buffer[position] != rune('N') {
						goto l1274
					}
					position++
				}
			l1290:
				if !_rules[rulesp]() {
					goto l1274
				}
				if !_rules[ruleExpressionOrWildcard]() {
					goto l1274
				}
				if !_rules[ruleAction78]() {
					goto l1274
				}
				add(ruleWhenThenPair, position1275)
			}
			return true
		l1274:
			position, tokenIndex = position1274, tokenIndex1274
			return false
		},
		/* 101 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1292, tokenIndex1292 := position, tokenIndex
			{
				position1293 := position
				{
					position1294, tokenIndex1294 := position, tokenIndex
					if !_rules[ruleFloatLiteral]() {
						goto l1295
					}
					goto l1294
				l1295:
					position, tokenIndex = position1294, tokenIndex1294
					if !_rules[ruleNumericLiteral]() {
						goto l1296
					}
					goto l1294
				l1296:
					position, tokenIndex = position1294, tokenIndex1294
					if !_rules[ruleStringLiteral]() {
						goto l1292
					}
				}
			l1294:
				add(ruleLiteral, position1293)
			}
			return true
		l1292:
			position, tokenIndex = position1292, tokenIndex1292
			return false
		},
		/* 102 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1297, tokenIndex1297 := position, tokenIndex
			{
				position1298 := position
				{
					position1299, tokenIndex1299 := position, tokenIndex
					if !_rules[ruleEqual]() {
						goto l1300
					}
					goto l1299
				l1300:
					position, tokenIndex = position1299, tokenIndex1299
					if !_rules[ruleNotEqual]() {
						goto l1301
					}
					goto l1299
				l1301:
					position, tokenIndex = position1299, tokenIndex1299
					if !_rules[ruleLessOrEqual]() {
						goto l1302
					}
					goto l1299
				l1302:
					position, tokenIndex = position1299, tokenIndex1299
					if !_rules[ruleLess]() {
						goto l1303
					}
					goto l1299
				l1303:
					position, tokenIndex = position1299, tokenIndex1299
					if !_rules[ruleGreaterOrEqual]() {
						goto l1304
					}
					goto l1299
				l1304:
					position, tokenIndex = position1299, tokenIndex1299
					if !_rules[ruleGreater]() {
						goto l1305
					}
					goto l1299
				l1305:
					position, tokenIndex = position1299, tokenIndex1299
					if !_rules[ruleNotEqual]() {
						goto l1297
					}
				}
			l1299:
				add(ruleComparisonOp, position1298)
			}
			return true
		l1297:
			position, tokenIndex = position1297, tokenIndex1297
			return false
		},
		/* 103 OtherOp <- <Concat> */
		func() bool {
			position1306, tokenIndex1306 := position, tokenIndex
			{
				position1307 := position
				if !_rules[ruleConcat]() {
					goto l1306
				}
				add(ruleOtherOp, position1307)
			}
			return true
		l1306:
			position, tokenIndex = position1306, tokenIndex1306
			return false
		},
		/* 104 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1308, tokenIndex1308 := position, tokenIndex
			{
				position1309 := position
				{
					position1310, tokenIndex1310 := position, tokenIndex
					if !_rules[ruleIsNot]() {
						goto l1311
					}
					goto l1310
				l1311:
					position, tokenIndex = position1310, tokenIndex1310
					if !_rules[ruleIs]() {
						goto l1308
					}
				}
			l1310:
				add(ruleIsOp, position1309)
			}
			return true
		l1308:
			position, tokenIndex = position1308, tokenIndex1308
			return false
		},
		/* 105 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1312, tokenIndex1312 := position, tokenIndex
			{
				position1313 := position
				{
					position1314, tokenIndex1314 := position, tokenIndex
					if !_rules[rulePlus]() {
						goto l1315
					}
					goto l1314
				l1315:
					position, tokenIndex = position1314, tokenIndex1314
					if !_rules[ruleMinus]() {
						goto l1312
					}
				}
			l1314:
				add(rulePlusMinusOp, position1313)
			}
			return true
		l1312:
			position, tokenIndex = position1312, tokenIndex1312
			return false
		},
		/* 106 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1316, tokenIndex1316 := position, tokenIndex
			{
				position1317 := position
				{
					position1318, tokenIndex1318 := position, tokenIndex
					if !_rules[ruleMultiply]() {
						goto l1319
					}
					goto l1318
				l1319:
					position, tokenIndex = position1318, tokenIndex1318
					if !_rules[ruleDivide]() {
						goto l1320
					}
					goto l1318
				l1320:
					position, tokenIndex = position1318, tokenIndex1318
					if !_rules[ruleModulo]() {
						goto l1316
					}
				}
			l1318:
				add(ruleMultDivOp, position1317)
			}
			return true
		l1316:
			position, tokenIndex = position1316, tokenIndex1316
			return false
		},
		/* 107 Stream <- <(<ident> Action79)> */
		func() bool {
			position1321, tokenIndex1321 := position, tokenIndex
			{
				position1322 := position
				{
					position1323 := position
					if !_rules[ruleident]() {
						goto l1321
					}
					add(rulePegText, position1323)
				}
				if !_rules[ruleAction79]() {
					goto l1321
				}
				add(ruleStream, position1322)
			}
			return true
		l1321:
			position, tokenIndex = position1321, tokenIndex1321
			return false
		},
		/* 108 RowMeta <- <RowTimestamp> */
		func() bool {
			position1324, tokenIndex1324 := position, tokenIndex
			{
				position1325 := position
				if !_rules[ruleRowTimestamp]() {
					goto l1324
				}
				add(ruleRowMeta, position1325)
			}
			return true
		l1324:
			position, tokenIndex = position1324, tokenIndex1324
			return false
		},
		/* 109 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action80)> */
		func() bool {
			position1326, tokenIndex1326 := position, tokenIndex
			{
				position1327 := position
				{
					position1328 := position
					{
						position1329, tokenIndex1329 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1329
						}
						if buffer[position] != rune(':') {
							goto l1329
						}
						position++
						goto l1330
					l1329:
						position, tokenIndex = position1329, tokenIndex1329
					}
				l1330:
					if buffer[position] != rune('t') {
						goto l1326
					}
					position++
					if buffer[position] != rune('s') {
						goto l1326
					}
					position++
					if buffer[position] != rune('(') {
						goto l1326
					}
					position++
					if buffer[position] != rune(')') {
						goto l1326
					}
					position++
					add(rulePegText, position1328)
				}
				if !_rules[ruleAction80]() {
					goto l1326
				}
				add(ruleRowTimestamp, position1327)
			}
			return true
		l1326:
			position, tokenIndex = position1326, tokenIndex1326
			return false
		},
		/* 110 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action81)> */
		func() bool {
			position1331, tokenIndex1331 := position, tokenIndex
			{
				position1332 := position
				{
					position1333 := position
					{
						position1334, tokenIndex1334 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1334
						}
						if buffer[position] != rune(':') {
							goto l1334
						}
						position++
						{
							position1336, tokenIndex1336 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l1336
							}
							position++
							goto l1334
						l1336:
							position, tokenIndex = position1336, tokenIndex1336
						}
						goto l1335
					l1334:
						position, tokenIndex = position1334, tokenIndex1334
					}
				l1335:
					if !_rules[rulejsonGetPath]() {
						goto l1331
					}
					add(rulePegText, position1333)
				}
				if !_rules[ruleAction81]() {
					goto l1331
				}
				add(ruleRowValue, position1332)
			}
			return true
		l1331:
			position, tokenIndex = position1331, tokenIndex1331
			return false
		},
		/* 111 NumericLiteral <- <(<('-'? [0-9]+)> Action82)> */
		func() bool {
			position1337, tokenIndex1337 := position, tokenIndex
			{
				position1338 := position
				{
					position1339 := position
					{
						position1340, tokenIndex1340 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1340
						}
						position++
						goto l1341
					l1340:
						position, tokenIndex = position1340, tokenIndex1340
					}
				l1341:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1337
					}
					position++
				l1342:
					{
						position1343, tokenIndex1343 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1343
						}
						position++
						goto l1342
					l1343:
						position, tokenIndex = position1343, tokenIndex1343
					}
					add(rulePegText, position1339)
				}
				if !_rules[ruleAction82]() {
					goto l1337
				}
				add(ruleNumericLiteral, position1338)
			}
			return true
		l1337:
			position, tokenIndex = position1337, tokenIndex1337
			return false
		},
		/* 112 NonNegativeNumericLiteral <- <(<[0-9]+> Action83)> */
		func() bool {
			position1344, tokenIndex1344 := position, tokenIndex
			{
				position1345 := position
				{
					position1346 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1344
					}
					position++
				l1347:
					{
						position1348, tokenIndex1348 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1348
						}
						position++
						goto l1347
					l1348:
						position, tokenIndex = position1348, tokenIndex1348
					}
					add(rulePegText, position1346)
				}
				if !_rules[ruleAction83]() {
					goto l1344
				}
				add(ruleNonNegativeNumericLiteral, position1345)
			}
			return true
		l1344:
			position, tokenIndex = position1344, tokenIndex1344
			return false
		},
		/* 113 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action84)> */
		func() bool {
			position1349, tokenIndex1349 := position, tokenIndex
			{
				position1350 := position
				{
					position1351 := position
					{
						position1352, tokenIndex1352 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1352
						}
						position++
						goto l1353
					l1352:
						position, tokenIndex = position1352, tokenIndex1352
					}
				l1353:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1349
					}
					position++
				l1354:
					{
						position1355, tokenIndex1355 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1355
						}
						position++
						goto l1354
					l1355:
						position, tokenIndex = position1355, tokenIndex1355
					}
					if buffer[position] != rune('.') {
						goto l1349
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1349
					}
					position++
				l1356:
					{
						position1357, tokenIndex1357 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1357
						}
						position++
						goto l1356
					l1357:
						position, tokenIndex = position1357, tokenIndex1357
					}
					add(rulePegText, position1351)
				}
				if !_rules[ruleAction84]() {
					goto l1349
				}
				add(ruleFloatLiteral, position1350)
			}
			return true
		l1349:
			position, tokenIndex = position1349, tokenIndex1349
			return false
		},
		/* 114 Function <- <(<ident> Action85)> */
		func() bool {
			position1358, tokenIndex1358 := position, tokenIndex
			{
				position1359 := position
				{
					position1360 := position
					if !_rules[ruleident]() {
						goto l1358
					}
					add(rulePegText, position1360)
				}
				if !_rules[ruleAction85]() {
					goto l1358
				}
				add(ruleFunction, position1359)
			}
			return true
		l1358:
			position, tokenIndex = position1358, tokenIndex1358
			return false
		},
		/* 115 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action86)> */
		func() bool {
			position1361, tokenIndex1361 := position, tokenIndex
			{
				position1362 := position
				{
					position1363 := position
					{
						position1364, tokenIndex1364 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1365
						}
						position++
						goto l1364
					l1365:
						position, tokenIndex = position1364, tokenIndex1364
						if buffer[position] != rune('N') {
							goto l1361
						}
						position++
					}
				l1364:
					{
						position1366, tokenIndex1366 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1367
						}
						position++
						goto l1366
					l1367:
						position, tokenIndex = position1366, tokenIndex1366
						if buffer[position] != rune('U') {
							goto l1361
						}
						position++
					}
				l1366:
					{
						position1368, tokenIndex1368 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1369
						}
						position++
						goto l1368
					l1369:
						position, tokenIndex = position1368, tokenIndex1368
						if buffer[position] != rune('L') {
							goto l1361
						}
						position++
					}
				l1368:
					{
						position1370, tokenIndex1370 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1371
						}
						position++
						goto l1370
					l1371:
						position, tokenIndex = position1370, tokenIndex1370
						if buffer[position] != rune('L') {
							goto l1361
						}
						position++
					}
				l1370:
					add(rulePegText, position1363)
				}
				if !_rules[ruleAction86]() {
					goto l1361
				}
				add(ruleNullLiteral, position1362)
			}
			return true
		l1361:
			position, tokenIndex = position1361, tokenIndex1361
			return false
		},
		/* 116 Missing <- <(<(('m' / 'M') ('i' / 'I') ('s' / 'S') ('s' / 'S') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action87)> */
		func() bool {
			position1372, tokenIndex1372 := position, tokenIndex
			{
				position1373 := position
				{
					position1374 := position
					{
						position1375, tokenIndex1375 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1376
						}
						position++
						goto l1375
					l1376:
						position, tokenIndex = position1375, tokenIndex1375
						if buffer[position] != rune('M') {
							goto l1372
						}
						position++
					}
				l1375:
					{
						position1377, tokenIndex1377 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1378
						}
						position++
						goto l1377
					l1378:
						position, tokenIndex = position1377, tokenIndex1377
						if buffer[position] != rune('I') {
							goto l1372
						}
						position++
					}
				l1377:
					{
						position1379, tokenIndex1379 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1380
						}
						position++
						goto l1379
					l1380:
						position, tokenIndex = position1379, tokenIndex1379
						if buffer[position] != rune('S') {
							goto l1372
						}
						position++
					}
				l1379:
					{
						position1381, tokenIndex1381 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1382
						}
						position++
						goto l1381
					l1382:
						position, tokenIndex = position1381, tokenIndex1381
						if buffer[position] != rune('S') {
							goto l1372
						}
						position++
					}
				l1381:
					{
						position1383, tokenIndex1383 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1384
						}
						position++
						goto l1383
					l1384:
						position, tokenIndex = position1383, tokenIndex1383
						if buffer[position] != rune('I') {
							goto l1372
						}
						position++
					}
				l1383:
					{
						position1385, tokenIndex1385 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1386
						}
						position++
						goto l1385
					l1386:
						position, tokenIndex = position1385, tokenIndex1385
						if buffer[position] != rune('N') {
							goto l1372
						}
						position++
					}
				l1385:
					{
						position1387, tokenIndex1387 := position, tokenIndex
						if buffer[position] != rune('g') {
							goto l1388
						}
						position++
						goto l1387
					l1388:
						position, tokenIndex = position1387, tokenIndex1387
						if buffer[position] != rune('G') {
							goto l1372
						}
						position++
					}
				l1387:
					add(rulePegText, position1374)
				}
				if !_rules[ruleAction87]() {
					goto l1372
				}
				add(ruleMissing, position1373)
			}
			return true
		l1372:
			position, tokenIndex = position1372, tokenIndex1372
			return false
		},
		/* 117 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1389, tokenIndex1389 := position, tokenIndex
			{
				position1390 := position
				{
					position1391, tokenIndex1391 := position, tokenIndex
					if !_rules[ruleTRUE]() {
						goto l1392
					}
					goto l1391
				l1392:
					position, tokenIndex = position1391, tokenIndex1391
					if !_rules[ruleFALSE]() {
						goto l1389
					}
				}
			l1391:
				add(ruleBooleanLiteral, position1390)
			}
			return true
		l1389:
			position, tokenIndex = position1389, tokenIndex1389
			return false
		},
		/* 118 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action88)> */
		func() bool {
			position1393, tokenIndex1393 := position, tokenIndex
			{
				position1394 := position
				{
					position1395 := position
					{
						position1396, tokenIndex1396 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1397
						}
						position++
						goto l1396
					l1397:
						position, tokenIndex = position1396, tokenIndex1396
						if buffer[position] != rune('T') {
							goto l1393
						}
						position++
					}
				l1396:
					{
						position1398, tokenIndex1398 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1399
						}
						position++
						goto l1398
					l1399:
						position, tokenIndex = position1398, tokenIndex1398
						if buffer[position] != rune('R') {
							goto l1393
						}
						position++
					}
				l1398:
					{
						position1400, tokenIndex1400 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1401
						}
						position++
						goto l1400
					l1401:
						position, tokenIndex = position1400, tokenIndex1400
						if buffer[position] != rune('U') {
							goto l1393
						}
						position++
					}
				l1400:
					{
						position1402, tokenIndex1402 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1403
						}
						position++
						goto l1402
					l1403:
						position, tokenIndex = position1402, tokenIndex1402
						if buffer[position] != rune('E') {
							goto l1393
						}
						position++
					}
				l1402:
					add(rulePegText, position1395)
				}
				if !_rules[ruleAction88]() {
					goto l1393
				}
				add(ruleTRUE, position1394)
			}
			return true
		l1393:
			position, tokenIndex = position1393, tokenIndex1393
			return false
		},
		/* 119 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action89)> */
		func() bool {
			position1404, tokenIndex1404 := position, tokenIndex
			{
				position1405 := position
				{
					position1406 := position
					{
						position1407, tokenIndex1407 := position, tokenIndex
						if buffer[position] != rune('f') {
							goto l1408
						}
						position++
						goto l1407
					l1408:
						position, tokenIndex = position1407, tokenIndex1407
						if buffer[position] != rune('F') {
							goto l1404
						}
						position++
					}
				l1407:
					{
						position1409, tokenIndex1409 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1410
						}
						position++
						goto l1409
					l1410:
						position, tokenIndex = position1409, tokenIndex1409
						if buffer[position] != rune('A') {
							goto l1404
						}
						position++
					}
				l1409:
					{
						position1411, tokenIndex1411 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1412
						}
						position++
						goto l1411
					l1412:
						position, tokenIndex = position1411, tokenIndex1411
						if buffer[position] != rune('L') {
							goto l1404
						}
						position++
					}
				l1411:
					{
						position1413, tokenIndex1413 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1414
						}
						position++
						goto l1413
					l1414:
						position, tokenIndex = position1413, tokenIndex1413
						if buffer[position] != rune('S') {
							goto l1404
						}
						position++
					}
				l1413:
					{
						position1415, tokenIndex1415 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1416
						}
						position++
						goto l1415
					l1416:
						position, tokenIndex = position1415, tokenIndex1415
						if buffer[position] != rune('E') {
							goto l1404
						}
						position++
					}
				l1415:
					add(rulePegText, position1406)
				}
				if !_rules[ruleAction89]() {
					goto l1404
				}
				add(ruleFALSE, position1405)
			}
			return true
		l1404:
			position, tokenIndex = position1404, tokenIndex1404
			return false
		},
		/* 120 Wildcard <- <(<((ident ':' !':')? '*')> Action90)> */
		func() bool {
			position1417, tokenIndex1417 := position, tokenIndex
			{
				position1418 := position
				{
					position1419 := position
					{
						position1420, tokenIndex1420 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1420
						}
						if buffer[position] != rune(':') {
							goto l1420
						}
						position++
						{
							position1422, tokenIndex1422 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l1422
							}
							position++
							goto l1420
						l1422:
							position, tokenIndex = position1422, tokenIndex1422
						}
						goto l1421
					l1420:
						position, tokenIndex = position1420, tokenIndex1420
					}
				l1421:
					if buffer[position] != rune('*') {
						goto l1417
					}
					position++
					add(rulePegText, position1419)
				}
				if !_rules[ruleAction90]() {
					goto l1417
				}
				add(ruleWildcard, position1418)
			}
			return true
		l1417:
			position, tokenIndex = position1417, tokenIndex1417
			return false
		},
		/* 121 StringLiteral <- <(<('"' (('"' '"') / (!'"' .))* '"')> Action91)> */
		func() bool {
			position1423, tokenIndex1423 := position, tokenIndex
			{
				position1424 := position
				{
					position1425 := position
					if buffer[position] != rune('"') {
						goto l1423
					}
					position++
				l1426:
					{
						position1427, tokenIndex1427 := position, tokenIndex
						{
							position1428, tokenIndex1428 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l1429
							}
							position++
							if buffer[position] != rune('"') {
								goto l1429
							}
							position++
							goto l1428
						l1429:
							position, tokenIndex = position1428, tokenIndex1428
							{
								position1430, tokenIndex1430 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l1430
								}
								position++
								goto l1427
							l1430:
								position, tokenIndex = position1430, tokenIndex1430
							}
							if !matchDot() {
								goto l1427
							}
						}
					l1428:
						goto l1426
					l1427:
						position, tokenIndex = position1427, tokenIndex1427
					}
					if buffer[position] != rune('"') {
						goto l1423
					}
					position++
					add(rulePegText, position1425)
				}
				if !_rules[ruleAction91]() {
					goto l1423
				}
				add(ruleStringLiteral, position1424)
			}
			return true
		l1423:
			position, tokenIndex = position1423, tokenIndex1423
			return false
		},
		/* 122 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action92)> */
		func() bool {
			position1431, tokenIndex1431 := position, tokenIndex
			{
				position1432 := position
				{
					position1433 := position
					{
						position1434, tokenIndex1434 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1435
						}
						position++
						goto l1434
					l1435:
						position, tokenIndex = position1434, tokenIndex1434
						if buffer[position] != rune('I') {
							goto l1431
						}
						position++
					}
				l1434:
					{
						position1436, tokenIndex1436 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1437
						}
						position++
						goto l1436
					l1437:
						position, tokenIndex = position1436, tokenIndex1436
						if buffer[position] != rune('S') {
							goto l1431
						}
						position++
					}
				l1436:
					{
						position1438, tokenIndex1438 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1439
						}
						position++
						goto l1438
					l1439:
						position, tokenIndex = position1438, tokenIndex1438
						if buffer[position] != rune('T') {
							goto l1431
						}
						position++
					}
				l1438:
					{
						position1440, tokenIndex1440 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1441
						}
						position++
						goto l1440
					l1441:
						position, tokenIndex = position1440, tokenIndex1440
						if buffer[position] != rune('R') {
							goto l1431
						}
						position++
					}
				l1440:
					{
						position1442, tokenIndex1442 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1443
						}
						position++
						goto l1442
					l1443:
						position, tokenIndex = position1442, tokenIndex1442
						if buffer[position] != rune('E') {
							goto l1431
						}
						position++
					}
				l1442:
					{
						position1444, tokenIndex1444 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1445
						}
						position++
						goto l1444
					l1445:
						position, tokenIndex = position1444, tokenIndex1444
						if buffer[position] != rune('A') {
							goto l1431
						}
						position++
					}
				l1444:
					{
						position1446, tokenIndex1446 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1447
						}
						position++
						goto l1446
					l1447:
						position, tokenIndex = position1446, tokenIndex1446
						if buffer[position] != rune('M') {
							goto l1431
						}
						position++
					}
				l1446:
					add(rulePegText, position1433)
				}
				if !_rules[ruleAction92]() {
					goto l1431
				}
				add(ruleISTREAM, position1432)
			}
			return true
		l1431:
			position, tokenIndex = position1431, tokenIndex1431
			return false
		},
		/* 123 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action93)> */
		func() bool {
			position1448, tokenIndex1448 := position, tokenIndex
			{
				position1449 := position
				{
					position1450 := position
					{
						position1451, tokenIndex1451 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1452
						}
						position++
						goto l1451
					l1452:
						position, tokenIndex = position1451, tokenIndex1451
						if buffer[position] != rune('D') {
							goto l1448
						}
						position++
					}
				l1451:
					{
						position1453, tokenIndex1453 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1454
						}
						position++
						goto l1453
					l1454:
						position, tokenIndex = position1453, tokenIndex1453
						if buffer[position] != rune('S') {
							goto l1448
						}
						position++
					}
				l1453:
					{
						position1455, tokenIndex1455 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1456
						}
						position++
						goto l1455
					l1456:
						position, tokenIndex = position1455, tokenIndex1455
						if buffer[position] != rune('T') {
							goto l1448
						}
						position++
					}
				l1455:
					{
						position1457, tokenIndex1457 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1458
						}
						position++
						goto l1457
					l1458:
						position, tokenIndex = position1457, tokenIndex1457
						if buffer[position] != rune('R') {
							goto l1448
						}
						position++
					}
				l1457:
					{
						position1459, tokenIndex1459 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1460
						}
						position++
						goto l1459
					l1460:
						position, tokenIndex = position1459, tokenIndex1459
						if buffer[position] != rune('E') {
							goto l1448
						}
						position++
					}
				l1459:
					{
						position1461, tokenIndex1461 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1462
						}
						position++
						goto l1461
					l1462:
						position, tokenIndex = position1461, tokenIndex1461
						if buffer[position] != rune('A') {
							goto l1448
						}
						position++
					}
				l1461:
					{
						position1463, tokenIndex1463 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1464
						}
						position++
						goto l1463
					l1464:
						position, tokenIndex = position1463, tokenIndex1463
						if buffer[position] != rune('M') {
							goto l1448
						}
						position++
					}
				l1463:
					add(rulePegText, position1450)
				}
				if !_rules[ruleAction93]() {
					goto l1448
				}
				add(ruleDSTREAM, position1449)
			}
			return true
		l1448:
			position, tokenIndex = position1448, tokenIndex1448
			return false
		},
		/* 124 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action94)> */
		func() bool {
			position1465, tokenIndex1465 := position, tokenIndex
			{
				position1466 := position
				{
					position1467 := position
					{
						position1468, tokenIndex1468 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1469
						}
						position++
						goto l1468
					l1469:
						position, tokenIndex = position1468, tokenIndex1468
						if buffer[position] != rune('R') {
							goto l1465
						}
						position++
					}
				l1468:
					{
						position1470, tokenIndex1470 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1471
						}
						position++
						goto l1470
					l1471:
						position, tokenIndex = position1470, tokenIndex1470
						if buffer[position] != rune('S') {
							goto l1465
						}
						position++
					}
				l1470:
					{
						position1472, tokenIndex1472 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1473
						}
						position++
						goto l1472
					l1473:
						position, tokenIndex = position1472, tokenIndex1472
						if buffer[position] != rune('T') {
							goto l1465
						}
						position++
					}
				l1472:
					{
						position1474, tokenIndex1474 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1475
						}
						position++
						goto l1474
					l1475:
						position, tokenIndex = position1474, tokenIndex1474
						if buffer[position] != rune('R') {
							goto l1465
						}
						position++
					}
				l1474:
					{
						position1476, tokenIndex1476 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1477
						}
						position++
						goto l1476
					l1477:
						position, tokenIndex = position1476, tokenIndex1476
						if buffer[position] != rune('E') {
							goto l1465
						}
						position++
					}
				l1476:
					{
						position1478, tokenIndex1478 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1479
						}
						position++
						goto l1478
					l1479:
						position, tokenIndex = position1478, tokenIndex1478
						if buffer[position] != rune('A') {
							goto l1465
						}
						position++
					}
				l1478:
					{
						position1480, tokenIndex1480 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1481
						}
						position++
						goto l1480
					l1481:
						position, tokenIndex = position1480, tokenIndex1480
						if buffer[position] != rune('M') {
							goto l1465
						}
						position++
					}
				l1480:
					add(rulePegText, position1467)
				}
				if !_rules[ruleAction94]() {
					goto l1465
				}
				add(ruleRSTREAM, position1466)
			}
			return true
		l1465:
			position, tokenIndex = position1465, tokenIndex1465
			return false
		},
		/* 125 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action95)> */
		func() bool {
			position1482, tokenIndex1482 := position, tokenIndex
			{
				position1483 := position
				{
					position1484 := position
					{
						position1485, tokenIndex1485 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1486
						}
						position++
						goto l1485
					l1486:
						position, tokenIndex = position1485, tokenIndex1485
						if buffer[position] != rune('T') {
							goto l1482
						}
						position++
					}
				l1485:
					{
						position1487, tokenIndex1487 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1488
						}
						position++
						goto l1487
					l1488:
						position, tokenIndex = position1487, tokenIndex1487
						if buffer[position] != rune('U') {
							goto l1482
						}
						position++
					}
				l1487:
					{
						position1489, tokenIndex1489 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1490
						}
						position++
						goto l1489
					l1490:
						position, tokenIndex = position1489, tokenIndex1489
						if buffer[position] != rune('P') {
							goto l1482
						}
						position++
					}
				l1489:
					{
						position1491, tokenIndex1491 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1492
						}
						position++
						goto l1491
					l1492:
						position, tokenIndex = position1491, tokenIndex1491
						if buffer[position] != rune('L') {
							goto l1482
						}
						position++
					}
				l1491:
					{
						position1493, tokenIndex1493 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1494
						}
						position++
						goto l1493
					l1494:
						position, tokenIndex = position1493, tokenIndex1493
						if buffer[position] != rune('E') {
							goto l1482
						}
						position++
					}
				l1493:
					{
						position1495, tokenIndex1495 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1496
						}
						position++
						goto l1495
					l1496:
						position, tokenIndex = position1495, tokenIndex1495
						if buffer[position] != rune('S') {
							goto l1482
						}
						position++
					}
				l1495:
					add(rulePegText, position1484)
				}
				if !_rules[ruleAction95]() {
					goto l1482
				}
				add(ruleTUPLES, position1483)
			}
			return true
		l1482:
			position, tokenIndex = position1482, tokenIndex1482
			return false
		},
		/* 126 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action96)> */
		func() bool {
			position1497, tokenIndex1497 := position, tokenIndex
			{
				position1498 := position
				{
					position1499 := position
					{
						position1500, tokenIndex1500 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1501
						}
						position++
						goto l1500
					l1501:
						position, tokenIndex = position1500, tokenIndex1500
						if buffer[position] != rune('S') {
							goto l1497
						}
						position++
					}
				l1500:
					{
						position1502, tokenIndex1502 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1503
						}
						position++
						goto l1502
					l1503:
						position, tokenIndex = position1502, tokenIndex1502
						if buffer[position] != rune('E') {
							goto l1497
						}
						position++
					}
				l1502:
					{
						position1504, tokenIndex1504 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1505
						}
						position++
						goto l1504
					l1505:
						position, tokenIndex = position1504, tokenIndex1504
						if buffer[position] != rune('C') {
							goto l1497
						}
						position++
					}
				l1504:
					{
						position1506, tokenIndex1506 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1507
						}
						position++
						goto l1506
					l1507:
						position, tokenIndex = position1506, tokenIndex1506
						if buffer[position] != rune('O') {
							goto l1497
						}
						position++
					}
				l1506:
					{
						position1508, tokenIndex1508 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1509
						}
						position++
						goto l1508
					l1509:
						position, tokenIndex = position1508, tokenIndex1508
						if buffer[position] != rune('N') {
							goto l1497
						}
						position++
					}
				l1508:
					{
						position1510, tokenIndex1510 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1511
						}
						position++
						goto l1510
					l1511:
						position, tokenIndex = position1510, tokenIndex1510
						if buffer[position] != rune('D') {
							goto l1497
						}
						position++
					}
				l1510:
					{
						position1512, tokenIndex1512 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1513
						}
						position++
						goto l1512
					l1513:
						position, tokenIndex = position1512, tokenIndex1512
						if buffer[position] != rune('S') {
							goto l1497
						}
						position++
					}
				l1512:
					add(rulePegText, position1499)
				}
				if !_rules[ruleAction96]() {
					goto l1497
				}
				add(ruleSECONDS, position1498)
			}
			return true
		l1497:
			position, tokenIndex = position1497, tokenIndex1497
			return false
		},
		/* 127 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action97)> */
		func() bool {
			position1514, tokenIndex1514 := position, tokenIndex
			{
				position1515 := position
				{
					position1516 := position
					{
						position1517, tokenIndex1517 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1518
						}
						position++
						goto l1517
					l1518:
						position, tokenIndex = position1517, tokenIndex1517
						if buffer[position] != rune('M') {
							goto l1514
						}
						position++
					}
				l1517:
					{
						position1519, tokenIndex1519 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1520
						}
						position++
						goto l1519
					l1520:
						position, tokenIndex = position1519, tokenIndex1519
						if buffer[position] != rune('I') {
							goto l1514
						}
						position++
					}
				l1519:
					{
						position1521, tokenIndex1521 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1522
						}
						position++
						goto l1521
					l1522:
						position, tokenIndex = position1521, tokenIndex1521
						if buffer[position] != rune('L') {
							goto l1514
						}
						position++
					}
				l1521:
					{
						position1523, tokenIndex1523 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1524
						}
						position++
						goto l1523
					l1524:
						position, tokenIndex = position1523, tokenIndex1523
						if buffer[position] != rune('L') {
							goto l1514
						}
						position++
					}
				l1523:
					{
						position1525, tokenIndex1525 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1526
						}
						position++
						goto l1525
					l1526:
						position, tokenIndex = position1525, tokenIndex1525
						if buffer[position] != rune('I') {
							goto l1514
						}
						position++
					}
				l1525:
					{
						position1527, tokenIndex1527 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1528
						}
						position++
						goto l1527
					l1528:
						position, tokenIndex = position1527, tokenIndex1527
						if buffer[position] != rune('S') {
							goto l1514
						}
						position++
					}
				l1527:
					{
						position1529, tokenIndex1529 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1530
						}
						position++
						goto l1529
					l1530:
						position, tokenIndex = position1529, tokenIndex1529
						if buffer[position] != rune('E') {
							goto l1514
						}
						position++
					}
				l1529:
					{
						position1531, tokenIndex1531 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1532
						}
						position++
						goto l1531
					l1532:
						position, tokenIndex = position1531, tokenIndex1531
						if buffer[position] != rune('C') {
							goto l1514
						}
						position++
					}
				l1531:
					{
						position1533, tokenIndex1533 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1534
						}
						position++
						goto l1533
					l1534:
						position, tokenIndex = position1533, tokenIndex1533
						if buffer[position] != rune('O') {
							goto l1514
						}
						position++
					}
				l1533:
					{
						position1535, tokenIndex1535 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1536
						}
						position++
						goto l1535
					l1536:
						position, tokenIndex = position1535, tokenIndex1535
						if buffer[position] != rune('N') {
							goto l1514
						}
						position++
					}
				l1535:
					{
						position1537, tokenIndex1537 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1538
						}
						position++
						goto l1537
					l1538:
						position, tokenIndex = position1537, tokenIndex1537
						if buffer[position] != rune('D') {
							goto l1514
						}
						position++
					}
				l1537:
					{
						position1539, tokenIndex1539 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1540
						}
						position++
						goto l1539
					l1540:
						position, tokenIndex = position1539, tokenIndex1539
						if buffer[position] != rune('S') {
							goto l1514
						}
						position++
					}
				l1539:
					add(rulePegText, position1516)
				}
				if !_rules[ruleAction97]() {
					goto l1514
				}
				add(ruleMILLISECONDS, position1515)
			}
			return true
		l1514:
			position, tokenIndex = position1514, tokenIndex1514
			return false
		},
		/* 128 Wait <- <(<(('w' / 'W') ('a' / 'A') ('i' / 'I') ('t' / 'T'))> Action98)> */
		func() bool {
			position1541, tokenIndex1541 := position, tokenIndex
			{
				position1542 := position
				{
					position1543 := position
					{
						position1544, tokenIndex1544 := position, tokenIndex
						if buffer[position] != rune('w') {
							goto l1545
						}
						position++
						goto l1544
					l1545:
						position, tokenIndex = position1544, tokenIndex1544
						if buffer[position] != rune('W') {
							goto l1541
						}
						position++
					}
				l1544:
					{
						position1546, tokenIndex1546 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1547
						}
						position++
						goto l1546
					l1547:
						position, tokenIndex = position1546, tokenIndex1546
						if buffer[position] != rune('A') {
							goto l1541
						}
						position++
					}
				l1546:
					{
						position1548, tokenIndex1548 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1549
						}
						position++
						goto l1548
					l1549:
						position, tokenIndex = position1548, tokenIndex1548
						if buffer[position] != rune('I') {
							goto l1541
						}
						position++
					}
				l1548:
					{
						position1550, tokenIndex1550 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1551
						}
						position++
						goto l1550
					l1551:
						position, tokenIndex = position1550, tokenIndex1550
						if buffer[position] != rune('T') {
							goto l1541
						}
						position++
					}
				l1550:
					add(rulePegText, position1543)
				}
				if !_rules[ruleAction98]() {
					goto l1541
				}
				add(ruleWait, position1542)
			}
			return true
		l1541:
			position, tokenIndex = position1541, tokenIndex1541
			return false
		},
		/* 129 DropOldest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('o' / 'O') ('l' / 'L') ('d' / 'D') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action99)> */
		func() bool {
			position1552, tokenIndex1552 := position, tokenIndex
			{
				position1553 := position
				{
					position1554 := position
					{
						position1555, tokenIndex1555 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1556
						}
						position++
						goto l1555
					l1556:
						position, tokenIndex = position1555, tokenIndex1555
						if buffer[position] != rune('D') {
							goto l1552
						}
						position++
					}
				l1555:
					{
						position1557, tokenIndex1557 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1558
						}
						position++
						goto l1557
					l1558:
						position, tokenIndex = position1557, tokenIndex1557
						if buffer[position] != rune('R') {
							goto l1552
						}
						position++
					}
				l1557:
					{
						position1559, tokenIndex1559 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1560
						}
						position++
						goto l1559
					l1560:
						position, tokenIndex = position1559, tokenIndex1559
						if buffer[position] != rune('O') {
							goto l1552
						}
						position++
					}
				l1559:
					{
						position1561, tokenIndex1561 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1562
						}
						position++
						goto l1561
					l1562:
						position, tokenIndex = position1561, tokenIndex1561
						if buffer[position] != rune('P') {
							goto l1552
						}
						position++
					}
				l1561:
					if !_rules[rulesp]() {
						goto l1552
					}
					{
						position1563, tokenIndex1563 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1564
						}
						position++
						goto l1563
					l1564:
						position, tokenIndex = position1563, tokenIndex1563
						if buffer[position] != rune('O') {
							goto l1552
						}
						position++
					}
				l1563:
					{
						position1565, tokenIndex1565 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1566
						}
						position++
						goto l1565
					l1566:
						position, tokenIndex = position1565, tokenIndex1565
						if buffer[position] != rune('L') {
							goto l1552
						}
						position++
					}
				l1565:
					{
						position1567, tokenIndex1567 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1568
						}
						position++
						goto l1567
					l1568:
						position, tokenIndex = position1567, tokenIndex1567
						if buffer[position] != rune('D') {
							goto l1552
						}
						position++
					}
				l1567:
					{
						position1569, tokenIndex1569 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1570
						}
						position++
						goto l1569
					l1570:
						position, tokenIndex = position1569, tokenIndex1569
						if buffer[position] != rune('E') {
							goto l1552
						}
						position++
					}
				l1569:
					{
						position1571, tokenIndex1571 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1572
						}
						position++
						goto l1571
					l1572:
						position, tokenIndex = position1571, tokenIndex1571
						if buffer[position] != rune('S') {
							goto l1552
						}
						position++
					}
				l1571:
					{
						position1573, tokenIndex1573 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1574
						}
						position++
						goto l1573
					l1574:
						position, tokenIndex = position1573, tokenIndex1573
						if buffer[position] != rune('T') {
							goto l1552
						}
						position++
					}
				l1573:
					add(rulePegText, position1554)
				}
				if !_rules[ruleAction99]() {
					goto l1552
				}
				add(ruleDropOldest, position1553)
			}
			return true
		l1552:
			position, tokenIndex = position1552, tokenIndex1552
			return false
		},
		/* 130 DropNewest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('n' / 'N') ('e' / 'E') ('w' / 'W') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action100)> */
		func() bool {
			position1575, tokenIndex1575 := position, tokenIndex
			{
				position1576 := position
				{
					position1577 := position
					{
						position1578, tokenIndex1578 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1579
						}
						position++
						goto l1578
					l1579:
						position, tokenIndex = position1578, tokenIndex1578
						if buffer[position] != rune('D') {
							goto l1575
						}
						position++
					}
				l1578:
					{
						position1580, tokenIndex1580 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1581
						}
						position++
						goto l1580
					l1581:
						position, tokenIndex = position1580, tokenIndex1580
						if buffer[position] != rune('R') {
							goto l1575
						}
						position++
					}
				l1580:
					{
						position1582, tokenIndex1582 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1583
						}
						position++
						goto l1582
					l1583:
						position, tokenIndex = position1582, tokenIndex1582
						if buffer[position] != rune('O') {
							goto l1575
						}
						position++
					}
				l1582:
					{
						position1584, tokenIndex1584 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1585
						}
						position++
						goto l1584
					l1585:
						position, tokenIndex = position1584, tokenIndex1584
						if buffer[position] != rune('P') {
							goto l1575
						}
						position++
					}
				l1584:
					if !_rules[rulesp]() {
						goto l1575
					}
					{
						position1586, tokenIndex1586 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1587
						}
						position++
						goto l1586
					l1587:
						position, tokenIndex = position1586, tokenIndex1586
						if buffer[position] != rune('N') {
							goto l1575
						}
						position++
					}
				l1586:
					{
						position1588, tokenIndex1588 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1589
						}
						position++
						goto l1588
					l1589:
						position, tokenIndex = position1588, tokenIndex1588
						if buffer[position] != rune('E') {
							goto l1575
						}
						position++
					}
				l1588:
					{
						position1590, tokenIndex1590 := position, tokenIndex
						if buffer[position] != rune('w') {
							goto l1591
						}
						position++
						goto l1590
					l1591:
						position, tokenIndex = position1590, tokenIndex1590
						if buffer[position] != rune('W') {
							goto l1575
						}
						position++
					}
				l1590:
					{
						position1592, tokenIndex1592 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1593
						}
						position++
						goto l1592
					l1593:
						position, tokenIndex = position1592, tokenIndex1592
						if buffer[position] != rune('E') {
							goto l1575
						}
						position++
					}
				l1592:
					{
						position1594, tokenIndex1594 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1595
						}
						position++
						goto l1594
					l1595:
						position, tokenIndex = position1594, tokenIndex1594
						if buffer[position] != rune('S') {
							goto l1575
						}
						position++
					}
				l1594:
					{
						position1596, tokenIndex1596 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1597
						}
						position++
						goto l1596
					l1597:
						position, tokenIndex = position1596, tokenIndex1596
						if buffer[position] != rune('T') {
							goto l1575
						}
						position++
					}
				l1596:
					add(rulePegText, position1577)
				}
				if !_rules[ruleAction100]() {
					goto l1575
				}
				add(ruleDropNewest, position1576)
			}
			return true
		l1575:
			position, tokenIndex = position1575, tokenIndex1575
			return false
		},
		/* 131 StreamIdentifier <- <(<ident> Action101)> */
		func() bool {
			position1598, tokenIndex1598 := position, tokenIndex
			{
				position1599 := position
				{
					position1600 := position
					if !_rules[ruleident]() {
						goto l1598
					}
					add(rulePegText, position1600)
				}
				if !_rules[ruleAction101]() {
					goto l1598
				}
				add(ruleStreamIdentifier, position1599)
			}
			return true
		l1598:
			position, tokenIndex = position1598, tokenIndex1598
			return false
		},
		/* 132 SourceSinkType <- <(<ident> Action102)> */
		func() bool {
			position1601, tokenIndex1601 := position, tokenIndex
			{
				position1602 := position
				{
					position1603 := position
					if !_rules[ruleident]() {
						goto l1601
					}
					add(rulePegText, position1603)
				}
				if !_rules[ruleAction102]() {
					goto l1601
				}
				add(ruleSourceSinkType, position1602)
			}
			return true
		l1601:
			position, tokenIndex = position1601, tokenIndex1601
			return false
		},
		/* 133 SourceSinkParamKey <- <(<ident> Action103)> */
		func() bool {
			position1604, tokenIndex1604 := position, tokenIndex
			{
				position1605 := position
				{
					position1606 := position
					if !_rules[ruleident]() {
						goto l1604
					}
					add(rulePegText, position1606)
				}
				if !_rules[ruleAction103]() {
					goto l1604
				}
				add(ruleSourceSinkParamKey, position1605)
			}
			return true
		l1604:
			position, tokenIndex = position1604, tokenIndex1604
			return false
		},
		/* 134 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action104)> */
		func() bool {
			position1607, tokenIndex1607 := position, tokenIndex
			{
				position1608 := position
				{
					position1609 := position
					{
						position1610, tokenIndex1610 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1611
						}
						position++
						goto l1610
					l1611:
						position, tokenIndex = position1610, tokenIndex1610
						if buffer[position] != rune('P') {
							goto l1607
						}
						position++
					}
				l1610:
					{
						position1612, tokenIndex1612 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1613
						}
						position++
						goto l1612
					l1613:
						position, tokenIndex = position1612, tokenIndex1612
						if buffer[position] != rune('A') {
							goto l1607
						}
						position++
					}
				l1612:
					{
						position1614, tokenIndex1614 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1615
						}
						position++
						goto l1614
					l1615:
						position, tokenIndex = position1614, tokenIndex1614
						if buffer[position] != rune('U') {
							goto l1607
						}
						position++
					}
				l1614:
					{
						position1616, tokenIndex1616 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1617
						}
						position++
						goto l1616
					l1617:
						position, tokenIndex = position1616, tokenIndex1616
						if buffer[position] != rune('S') {
							goto l1607
						}
						position++
					}
				l1616:
					{
						position1618, tokenIndex1618 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1619
						}
						position++
						goto l1618
					l1619:
						position, tokenIndex = position1618, tokenIndex1618
						if buffer[position] != rune('E') {
							goto l1607
						}
						position++
					}
				l1618:
					{
						position1620, tokenIndex1620 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1621
						}
						position++
						goto l1620
					l1621:
						position, tokenIndex = position1620, tokenIndex1620
						if buffer[position] != rune('D') {
							goto l1607
						}
						position++
					}
				l1620:
					add(rulePegText, position1609)
				}
				if !_rules[ruleAction104]() {
					goto l1607
				}
				add(rulePaused, position1608)
			}
			return true
		l1607:
			position, tokenIndex = position1607, tokenIndex1607
			return false
		},
		/* 135 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action105)> */
		func() bool {
			position1622, tokenIndex1622 := position, tokenIndex
			{
				position1623 := position
				{
					position1624 := position
					{
						position1625, tokenIndex1625 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1626
						}
						position++
						goto l1625
					l1626:
						position, tokenIndex = position1625, tokenIndex1625
						if buffer[position] != rune('U') {
							goto l1622
						}
						position++
					}
				l1625:
					{
						position1627, tokenIndex1627 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1628
						}
						position++
						goto l1627
					l1628:
						position, tokenIndex = position1627, tokenIndex1627
						if buffer[position] != rune('N') {
							goto l1622
						}
						position++
					}
				l1627:
					{
						position1629, tokenIndex1629 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1630
						}
						position++
						goto l1629
					l1630:
						position, tokenIndex = position1629, tokenIndex1629
						if buffer[position] != rune('P') {
							goto l1622
						}
						position++
					}
				l1629:
					{
						position1631, tokenIndex1631 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1632
						}
						position++
						goto l1631
					l1632:
						position, tokenIndex = position1631, tokenIndex1631
						if buffer[position] != rune('A') {
							goto l1622
						}
						position++
					}
				l1631:
					{
						position1633, tokenIndex1633 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1634
						}
						position++
						goto l1633
					l1634:
						position, tokenIndex = position1633, tokenIndex1633
						if buffer[position] != rune('U') {
							goto l1622
						}
						position++
					}
				l1633:
					{
						position1635, tokenIndex1635 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1636
						}
						position++
						goto l1635
					l1636:
						position, tokenIndex = position1635, tokenIndex1635
						if buffer[position] != rune('S') {
							goto l1622
						}
						position++
					}
				l1635:
					{
						position1637, tokenIndex1637 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1638
						}
						position++
						goto l1637
					l1638:
						position, tokenIndex = position1637, tokenIndex1637
						if buffer[position] != rune('E') {
							goto l1622
						}
						position++
					}
				l1637:
					{
						position1639, tokenIndex1639 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1640
						}
						position++
						goto l1639
					l1640:
						position, tokenIndex = position1639, tokenIndex1639
						if buffer[position] != rune('D') {
							goto l1622
						}
						position++
					}
				l1639:
					add(rulePegText, position1624)
				}
				if !_rules[ruleAction105]() {
					goto l1622
				}
				add(ruleUnpaused, position1623)
			}
			return true
		l1622:
			position, tokenIndex = position1622, tokenIndex1622
			return false
		},
		/* 136 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action106)> */
		func() bool {
			position1641, tokenIndex1641 := position, tokenIndex
			{
				position1642 := position
				{
					position1643 := position
					{
						position1644, tokenIndex1644 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1645
						}
						position++
						goto l1644
					l1645:
						position, tokenIndex = position1644, tokenIndex1644
						if buffer[position] != rune('A') {
							goto l1641
						}
						position++
					}
				l1644:
					{
						position1646, tokenIndex1646 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1647
						}
						position++
						goto l1646
					l1647:
						position, tokenIndex = position1646, tokenIndex1646
						if buffer[position] != rune('S') {
							goto l1641
						}
						position++
					}
				l1646:
					{
						position1648, tokenIndex1648 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1649
						}
						position++
						goto l1648
					l1649:
						position, tokenIndex = position1648, tokenIndex1648
						if buffer[position] != rune('C') {
							goto l1641
						}
						position++
					}
				l1648:
					add(rulePegText, position1643)
				}
				if !_rules[ruleAction106]() {
					goto l1641
				}
				add(ruleAscending, position1642)
			}
			return true
		l1641:
			position, tokenIndex = position1641, tokenIndex1641
			return false
		},
		/* 137 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action107)> */
		func() bool {
			position1650, tokenIndex1650 := position, tokenIndex
			{
				position1651 := position
				{
					position1652 := position
					{
						position1653, tokenIndex1653 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1654
						}
						position++
						goto l1653
					l1654:
						position, tokenIndex = position1653, tokenIndex1653
						if buffer[position] != rune('D') {
							goto l1650
						}
						position++
					}
				l1653:
					{
						position1655, tokenIndex1655 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1656
						}
						position++
						goto l1655
					l1656:
						position, tokenIndex = position1655, tokenIndex1655
						if buffer[position] != rune('E') {
							goto l1650
						}
						position++
					}
				l1655:
					{
						position1657, tokenIndex1657 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1658
						}
						position++
						goto l1657
					l1658:
						position, tokenIndex = position1657, tokenIndex1657
						if buffer[position] != rune('S') {
							goto l1650
						}
						position++
					}
				l1657:
					{
						position1659, tokenIndex1659 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1660
						}
						position++
						goto l1659
					l1660:
						position, tokenIndex = position1659, tokenIndex1659
						if buffer[position] != rune('C') {
							goto l1650
						}
						position++
					}
				l1659:
					add(rulePegText, position1652)
				}
				if !_rules[ruleAction107]() {
					goto l1650
				}
				add(ruleDescending, position1651)
			}
			return true
		l1650:
			position, tokenIndex = position1650, tokenIndex1650
			return false
		},
		/* 138 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1661, tokenIndex1661 := position, tokenIndex
			{
				position1662 := position
				{
					position1663, tokenIndex1663 := position, tokenIndex
					if !_rules[ruleBool]() {
						goto l1664
					}
					goto l1663
				l1664:
					position, tokenIndex = position1663, tokenIndex1663
					if !_rules[ruleInt]() {
						goto l1665
					}
					goto l1663
				l1665:
					position, tokenIndex = position1663, tokenIndex1663
					if !_rules[ruleFloat]() {
						goto l1666
					}
					goto l1663
				l1666:
					position, tokenIndex = position1663, tokenIndex1663
					if !_rules[ruleString]() {
						goto l1667
					}
					goto l1663
				l1667:
					position, tokenIndex = position1663, tokenIndex1663
					if !_rules[ruleBlob]() {
						goto l1668
					}
					goto l1663
				l1668:
					position, tokenIndex = position1663, tokenIndex1663
					if !_rules[ruleTimestamp]() {
						goto l1669
					}
					goto l1663
				l1669:
					position, tokenIndex = position1663, tokenIndex1663
					if !_rules[ruleArray]() {
						goto l1670
					}
					goto l1663
				l1670:
					position, tokenIndex = position1663, tokenIndex1663
					if !_rules[ruleMap]() {
						goto l1661
					}
				}
			l1663:
				add(ruleType, position1662)
			}
			return true
		l1661:
			position, tokenIndex = position1661, tokenIndex1661
			return false
		},
		/* 139 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action108)> */
		func() bool {
			position1671, tokenIndex1671 := position, tokenIndex
			{
				position1672 := position
				{
					position1673 := position
					{
						position1674, tokenIndex1674 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1675
						}
						position++
						goto l1674
					l1675:
						position, tokenIndex = position1674, tokenIndex1674
						if buffer[position] != rune('B') {
							goto l1671
						}
						position++
					}
				l1674:
					{
						position1676, tokenIndex1676 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1677
						}
						position++
						goto l1676
					l1677:
						position, tokenIndex = position1676, tokenIndex1676
						if buffer[position] != rune('O') {
							goto l1671
						}
						position++
					}
				l1676:
					{
						position1678, tokenIndex1678 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1679
						}
						position++
						goto l1678
					l1679:
						position, tokenIndex = position1678, tokenIndex1678
						if buffer[position] != rune('O') {
							goto l1671
						}
						position++
					}
				l1678:
					{
						position1680, tokenIndex1680 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1681
						}
						position++
						goto l1680
					l1681:
						position, tokenIndex = position1680, tokenIndex1680
						if buffer[position] != rune('L') {
							goto l1671
						}
						position++
					}
				l1680:
					add(rulePegText, position1673)
				}
				if !_rules[ruleAction108]() {
					goto l1671
				}
				add(ruleBool, position1672)
			}
			return true
		l1671:
			position, tokenIndex = position1671, tokenIndex1671
			return false
		},
		/* 140 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action109)> */
		func() bool {
			position1682, tokenIndex1682 := position, tokenIndex
			{
				position1683 := position
				{
					position1684 := position
					{
						position1685, tokenIndex1685 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1686
						}
						position++
						goto l1685
					l1686:
						position, tokenIndex = position1685, tokenIndex1685
						if buffer[position] != rune('I') {
							goto l1682
						}
						position++
					}
				l1685:
					{
						position1687, tokenIndex1687 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1688
						}
						position++
						goto l1687
					l1688:
						position, tokenIndex = position1687, tokenIndex1687
						if buffer[position] != rune('N') {
							goto l1682
						}
						position++
					}
				l1687:
					{
						position1689, tokenIndex1689 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1690
						}
						position++
						goto l1689
					l1690:
						position, tokenIndex = position1689, tokenIndex1689
						if buffer[position] != rune('T') {
							goto l1682
						}
						position++
					}
				l1689:
					add(rulePegText, position1684)
				}
				if !_rules[ruleAction109]() {
					goto l1682
				}
				add(ruleInt, position1683)
			}
			return true
		l1682:
			position, tokenIndex = position1682, tokenIndex1682
			return false
		},
		/* 141 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action110)> */
		func() bool {
			position1691, tokenIndex1691 := position, tokenIndex
			{
				position1692 := position
				{
					position1693 := position
					{
						position1694, tokenIndex1694 := position, tokenIndex
						if buffer[position] != rune('f') {
							goto l1695
						}
						position++
						goto l1694
					l1695:
						position, tokenIndex = position1694, tokenIndex1694
						if buffer[position] != rune('F') {
							goto l1691
						}
						position++
					}
				l1694:
					{
						position1696, tokenIndex1696 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1697
						}
						position++
						goto l1696
					l1697:
						position, tokenIndex = position1696, tokenIndex1696
						if buffer[position] != rune('L') {
							goto l1691
						}
						position++
					}
				l1696:
					{
						position1698, tokenIndex1698 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1699
						}
						position++
						goto l1698
					l1699:
						position, tokenIndex = position1698, tokenIndex1698
						if buffer[position] != rune('O') {
							goto l1691
						}
						position++
					}
				l1698:
					{
						position1700, tokenIndex1700 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1701
						}
						position++
						goto l1700
					l1701:
						position, tokenIndex = position1700, tokenIndex1700
						if buffer[position] != rune('A') {
							goto l1691
						}
						position++
					}
				l1700:
					{
						position1702, tokenIndex1702 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1703
						}
						position++
						goto l1702
					l1703:
						position, tokenIndex = position1702, tokenIndex1702
						if buffer[position] != rune('T') {
							goto l1691
						}
						position++
					}
				l1702:
					add(rulePegText, position1693)
				}
				if !_rules[ruleAction110]() {
					goto l1691
				}
				add(ruleFloat, position1692)
			}
			return true
		l1691:
			position, tokenIndex = position1691, tokenIndex1691
			return false
		},
		/* 142 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action111)> */
		func() bool {
			position1704, tokenIndex1704 := position, tokenIndex
			{
				position1705 := position
				{
					position1706 := position
					{
						position1707, tokenIndex1707 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1708
						}
						position++
						goto l1707
					l1708:
						position, tokenIndex = position1707, tokenIndex1707
						if buffer[position] != rune('S') {
							goto l1704
						}
						position++
					}
				l1707:
					{
						position1709, tokenIndex1709 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1710
						}
						position++
						goto l1709
					l1710:
						position, tokenIndex = position1709, tokenIndex1709
						if buffer[position] != rune('T') {
							goto l1704
						}
						position++
					}
				l1709:
					{
						position1711, tokenIndex1711 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1712
						}
						position++
						goto l1711
					l1712:
						position, tokenIndex = position1711, tokenIndex1711
						if buffer[position] != rune('R') {
							goto l1704
						}
						position++
					}
				l1711:
					{
						position1713, tokenIndex1713 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1714
						}
						position++
						goto l1713
					l1714:
						position, tokenIndex = position1713, tokenIndex1713
						if buffer[position] != rune('I') {
							goto l1704
						}
						position++
					}
				l1713:
					{
						position1715, tokenIndex1715 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1716
						}
						position++
						goto l1715
					l1716:
						position, tokenIndex = position1715, tokenIndex1715
						if buffer[position] != rune('N') {
							goto l1704
						}
						position++
					}
				l1715:
					{
						position1717, tokenIndex1717 := position, tokenIndex
						if buffer[position] != rune('g') {
							goto l1718
						}
						position++
						goto l1717
					l1718:
						position, tokenIndex = position1717, tokenIndex1717
						if buffer[position] != rune('G') {
							goto l1704
						}
						position++
					}
				l1717:
					add(rulePegText, position1706)
				}
				if !_rules[ruleAction111]() {
					goto l1704
				}
				add(ruleString, position1705)
			}
			return true
		l1704:
			position, tokenIndex = position1704, tokenIndex1704
			return false
		},
		/* 143 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action112)> */
		func() bool {
			position1719, tokenIndex1719 := position, tokenIndex
			{
				position1720 := position
				{
					position1721 := position
					{
						position1722, tokenIndex1722 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1723
						}
						position++
						goto l1722
					l1723:
						position, tokenIndex = position1722, tokenIndex1722
						if buffer[position] != rune('B') {
							goto l1719
						}
						position++
					}
				l1722:
					{
						position1724, tokenIndex1724 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1725
						}
						position++
						goto l1724
					l1725:
						position, tokenIndex = position1724, tokenIndex1724
						if buffer[position] != rune('L') {
							goto l1719
						}
						position++
					}
				l1724:
					{
						position1726, tokenIndex1726 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1727
						}
						position++
						goto l1726
					l1727:
						position, tokenIndex = position1726, tokenIndex1726
						if buffer[position] != rune('O') {
							goto l1719
						}
						position++
					}
				l1726:
					{
						position1728, tokenIndex1728 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1729
						}
						position++
						goto l1728
					l1729:
						position, tokenIndex = position1728, tokenIndex1728
						if buffer[position] != rune('B') {
							goto l1719
						}
						position++
					}
				l1728:
					add(rulePegText, position1721)
				}
				if !_rules[ruleAction112]() {
					goto l1719
				}
				add(ruleBlob, position1720)
			}
			return true
		l1719:
			position, tokenIndex = position1719, tokenIndex1719
			return false
		},
		/* 144 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action113)> */
		func() bool {
			position1730, tokenIndex1730 := position, tokenIndex
			{
				position1731 := position
				{
					position1732 := position
					{
						position1733, tokenIndex1733 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1734
						}
						position++
						goto l1733
					l1734:
						position, tokenIndex = position1733, tokenIndex1733
						if buffer[position] != rune('T') {
							goto l1730
						}
						position++
					}
				l1733:
					{
						position1735, tokenIndex1735 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1736
						}
						position++
						goto l1735
					l1736:
						position, tokenIndex = position1735, tokenIndex1735
						if buffer[position] != rune('I') {
							goto l1730
						}
						position++
					}
				l1735:
					{
						position1737, tokenIndex1737 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1738
						}
						position++
						goto l1737
					l1738:
						position, tokenIndex = position1737, tokenIndex1737
						if buffer[position] != rune('M') {
							goto l1730
						}
						position++
					}
				l1737:
					{
						position1739, tokenIndex1739 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1740
						}
						position++
						goto l1739
					l1740:
						position, tokenIndex = position1739, tokenIndex1739
						if buffer[position] != rune('E') {
							goto l1730
						}
						position++
					}
				l1739:
					{
						position1741, tokenIndex1741 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1742
						}
						position++
						goto l1741
					l1742:
						position, tokenIndex = position1741, tokenIndex1741
						if buffer[position] != rune('S') {
							goto l1730
						}
						position++
					}
				l1741:
					{
						position1743, tokenIndex1743 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1744
						}
						position++
						goto l1743
					l1744:
						position, tokenIndex = position1743, tokenIndex1743
						if buffer[position] != rune('T') {
							goto l1730
						}
						position++
					}
				l1743:
					{
						position1745, tokenIndex1745 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1746
						}
						position++
						goto l1745
					l1746:
						position, tokenIndex = position1745, tokenIndex1745
						if buffer[position] != rune('A') {
							goto l1730
						}
						position++
					}
				l1745:
					{
						position1747, tokenIndex1747 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1748
						}
						position++
						goto l1747
					l1748:
						position, tokenIndex = position1747, tokenIndex1747
						if buffer[position] != rune('M') {
							goto l1730
						}
						position++
					}
				l1747:
					{
						position1749, tokenIndex1749 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1750
						}
						position++
						goto l1749
					l1750:
						position, tokenIndex = position1749, tokenIndex1749
						if buffer[position] != rune('P') {
							goto l1730
						}
						position++
					}
				l1749:
					add(rulePegText, position1732)
				}
				if !_rules[ruleAction113]() {
					goto l1730
				}
				add(ruleTimestamp, position1731)
			}
			return true
		l1730:
			position, tokenIndex = position1730, tokenIndex1730
			return false
		},
		/* 145 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action114)> */
		func() bool {
			position1751, tokenIndex1751 := position, tokenIndex
			{
				position1752 := position
				{
					position1753 := position
					{
						position1754, tokenIndex1754 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1755
						}
						position++
						goto l1754
					l1755:
						position, tokenIndex = position1754, tokenIndex1754
						if buffer[position] != rune('A') {
							goto l1751
						}
						position++
					}
				l1754:
					{
						position1756, tokenIndex1756 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1757
						}
						position++
						goto l1756
					l1757:
						position, tokenIndex = position1756, tokenIndex1756
						if buffer[position] != rune('R') {
							goto l1751
						}
						position++
					}
				l1756:
					{
						position1758, tokenIndex1758 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1759
						}
						position++
						goto l1758
					l1759:
						position, tokenIndex = position1758, tokenIndex1758
						if buffer[position] != rune('R') {
							goto l1751
						}
						position++
					}
				l1758:
					{
						position1760, tokenIndex1760 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1761
						}
						position++
						goto l1760
					l1761:
						position, tokenIndex = position1760, tokenIndex1760
						if buffer[position] != rune('A') {
							goto l1751
						}
						position++
					}
				l1760:
					{
						position1762, tokenIndex1762 := position, tokenIndex
						if buffer[position] != rune('y') {
							goto l1763
						}
						position++
						goto l1762
					l1763:
						position, tokenIndex = position1762, tokenIndex1762
						if buffer[position] != rune('Y') {
							goto l1751
						}
						position++
					}
				l1762:
					add(rulePegText, position1753)
				}
				if !_rules[ruleAction114]() {
					goto l1751
				}
				add(ruleArray, position1752)
			}
			return true
		l1751:
			position, tokenIndex = position1751, tokenIndex1751
			return false
		},
		/* 146 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action115)> */
		func() bool {
			position1764, tokenIndex1764 := position, tokenIndex
			{
				position1765 := position
				{
					position1766 := position
					{
						position1767, tokenIndex1767 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1768
						}
						position++
						goto l1767
					l1768:
						position, tokenIndex = position1767, tokenIndex1767
						if buffer[position] != rune('M') {
							goto l1764
						}
						position++
					}
				l1767:
					{
						position1769, tokenIndex1769 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1770
						}
						position++
						goto l1769
					l1770:
						position, tokenIndex = position1769, tokenIndex1769
						if buffer[position] != rune('A') {
							goto l1764
						}
						position++
					}
				l1769:
					{
						position1771, tokenIndex1771 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1772
						}
						position++
						goto l1771
					l1772:
						position, tokenIndex = position1771, tokenIndex1771
						if buffer[position] != rune('P') {
							goto l1764
						}
						position++
					}
				l1771:
					add(rulePegText, position1766)
				}
				if !_rules[ruleAction115]() {
					goto l1764
				}
				add(ruleMap, position1765)
			}
			return true
		l1764:
			position, tokenIndex = position1764, tokenIndex1764
			return false
		},
		/* 147 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action116)> */
		func() bool {
			position1773, tokenIndex1773 := position, tokenIndex
			{
				position1774 := position
				{
					position1775 := position
					{
						position1776, tokenIndex1776 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1777
						}
						position++
						goto l1776
					l1777:
						position, tokenIndex = position1776, tokenIndex1776
						if buffer[position] != rune('O') {
							goto l1773
						}
						position++
					}
				l1776:
					{
						position1778, tokenIndex1778 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1779
						}
						position++
						goto l1778
					l1779:
						position, tokenIndex = position1778, tokenIndex1778
						if buffer[position] != rune('R') {
							goto l1773
						}
						position++
					}
				l1778:
					add(rulePegText, position1775)
				}
				if !_rules[ruleAction116]() {
					goto l1773
				}
				add(ruleOr, position1774)
			}
			return true
		l1773:
			position, tokenIndex = position1773, tokenIndex1773
			return false
		},
		/* 148 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action117)> */
		func() bool {
			position1780, tokenIndex1780 := position, tokenIndex
			{
				position1781 := position
				{
					position1782 := position
					{
						position1783, tokenIndex1783 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1784
						}
						position++
						goto l1783
					l1784:
						position, tokenIndex = position1783, tokenIndex1783
						if buffer[position] != rune('A') {
							goto l1780
						}
						position++
					}
				l1783:
					{
						position1785, tokenIndex1785 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1786
						}
						position++
						goto l1785
					l1786:
						position, tokenIndex = position1785, tokenIndex1785
						if buffer[position] != rune('N') {
							goto l1780
						}
						position++
					}
				l1785:
					{
						position1787, tokenIndex1787 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1788
						}
						position++
						goto l1787
					l1788:
						position, tokenIndex = position1787, tokenIndex1787
						if buffer[position] != rune('D') {
							goto l1780
						}
						position++
					}
				l1787:
					add(rulePegText, position1782)
				}
				if !_rules[ruleAction117]() {
					goto l1780
				}
				add(ruleAnd, position1781)
			}
			return true
		l1780:
			position, tokenIndex = position1780, tokenIndex1780
			return false
		},
		/* 149 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action118)> */
		func() bool {
			position1789, tokenIndex1789 := position, tokenIndex
			{
				position1790 := position
				{
					position1791 := position
					{
						position1792, tokenIndex1792 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1793
						}
						position++
						goto l1792
					l1793:
						position, tokenIndex = position1792, tokenIndex1792
						if buffer[position] != rune('N') {
							goto l1789
						}
						position++
					}
				l1792:
					{
						position1794, tokenIndex1794 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1795
						}
						position++
						goto l1794
					l1795:
						position, tokenIndex = position1794, tokenIndex1794
						if buffer[position] != rune('O') {
							goto l1789
						}
						position++
					}
				l1794:
					{
						position1796, tokenIndex1796 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1797
						}
						position++
						goto l1796
					l1797:
						position, tokenIndex = position1796, tokenIndex1796
						if buffer[position] != rune('T') {
							goto l1789
						}
						position++
					}
				l1796:
					add(rulePegText, position1791)
				}
				if !_rules[ruleAction118]() {
					goto l1789
				}
				add(ruleNot, position1790)
			}
			return true
		l1789:
			position, tokenIndex = position1789, tokenIndex1789
			return false
		},
		/* 150 Equal <- <(<'='> Action119)> */
		func() bool {
			position1798, tokenIndex1798 := position, tokenIndex
			{
				position1799 := position
				{
					position1800 := position
					if buffer[position] != rune('=') {
						goto l1798
					}
					position++
					add(rulePegText, position1800)
				}
				if !_rules[ruleAction119]() {
					goto l1798
				}
				add(ruleEqual, position1799)
			}
			return true
		l1798:
			position, tokenIndex = position1798, tokenIndex1798
			return false
		},
		/* 151 Less <- <(<'<'> Action120)> */
		func() bool {
			position1801, tokenIndex1801 := position, tokenIndex
			{
				position1802 := position
				{
					position1803 := position
					if buffer[position] != rune('<') {
						goto l1801
					}
					position++
					add(rulePegText, position1803)
				}
				if !_rules[ruleAction120]() {
					goto l1801
				}
				add(ruleLess, position1802)
			}
			return true
		l1801:
			position, tokenIndex = position1801, tokenIndex1801
			return false
		},
		/* 152 LessOrEqual <- <(<('<' '=')> Action121)> */
		func() bool {
			position1804, tokenIndex1804 := position, tokenIndex
			{
				position1805 := position
				{
					position1806 := position
					if buffer[position] != rune('<') {
						goto l1804
					}
					position++
					if buffer[position] != rune('=') {
						goto l1804
					}
					position++
					add(rulePegText, position1806)
				}
				if !_rules[ruleAction121]() {
					goto l1804
				}
				add(ruleLessOrEqual, position1805)
			}
			return true
		l1804:
			position, tokenIndex = position1804, tokenIndex1804
			return false
		},
		/* 153 Greater <- <(<'>'> Action122)> */
		func() bool {
			position1807, tokenIndex1807 := position, tokenIndex
			{
				position1808 := position
				{
					position1809 := position
					if buffer[position] != rune('>') {
						goto l1807
					}
					position++
					add(rulePegText, position1809)
				}
				if !_rules[ruleAction122]() {
					goto l1807
				}
				add(ruleGreater, position1808)
			}
			return true
		l1807:
			position, tokenIndex = position1807, tokenIndex1807
			return false
		},
		/* 154 GreaterOrEqual <- <(<('>' '=')> Action123)> */
		func() bool {
			position1810, tokenIndex1810 := position, tokenIndex
			{
				position1811 := position
				{
					position1812 := position
					if buffer[position] != rune('>') {
						goto l1810
					}
					position++
					if buffer[position] != rune('=') {
						goto l1810
					}
					position++
					add(rulePegText, position1812)
				}
				if !_rules[ruleAction123]() {
					goto l1810
				}
				add(ruleGreaterOrEqual, position1811)
			}
			return true
		l1810:
			position, tokenIndex = position1810, tokenIndex1810
			return false
		},
		/* 155 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action124)> */
		func() bool {
			position1813, tokenIndex1813 := position, tokenIndex
			{
				position1814 := position
				{
					position1815 := position
					{
						position1816, tokenIndex1816 := position, tokenIndex
						if buffer[position] != rune('!') {
							goto l1817
						}
						position++
						if buffer[position] != rune('=') {
							goto l1817
						}
						position++
						goto l1816
					l1817:
						position, tokenIndex = position1816, tokenIndex1816
						if buffer[position] != rune('<') {
							goto l1813
						}
						position++
						if buffer[position] != rune('>') {
							goto l1813
						}
						position++
					}
				l1816:
					add(rulePegText, position1815)
				}
				if !_rules[ruleAction124]() {
					goto l1813
				}
				add(ruleNotEqual, position1814)
			}
			return true
		l1813:
			position, tokenIndex = position1813, tokenIndex1813
			return false
		},
		/* 156 Concat <- <(<('|' '|')> Action125)> */
		func() bool {
			position1818, tokenIndex1818 := position, tokenIndex
			{
				position1819 := position
				{
					position1820 := position
					if buffer[position] != rune('|') {
						goto l1818
					}
					position++
					if buffer[position] != rune('|') {
						goto l1818
					}
					position++
					add(rulePegText, position1820)
				}
				if !_rules[ruleAction125]() {
					goto l1818
				}
				add(ruleConcat, position1819)
			}
			return true
		l1818:
			position, tokenIndex = position1818, tokenIndex1818
			return false
		},
		/* 157 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action126)> */
		func() bool {
			position1821, tokenIndex1821 := position, tokenIndex
			{
				position1822 := position
				{
					position1823 := position
					{
						position1824, tokenIndex1824 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1825
						}
						position++
						goto l1824
					l1825:
						position, tokenIndex = position1824, tokenIndex1824
						if buffer[position] != rune('I') {
							goto l1821
						}
						position++
					}
				l1824:
					{
						position1826, tokenIndex1826 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1827
						}
						position++
						goto l1826
					l1827:
						position, tokenIndex = position1826, tokenIndex1826
						if buffer[position] != rune('S') {
							goto l1821
						}
						position++
					}
				l1826:
					add(rulePegText, position1823)
				}
				if !_rules[ruleAction126]() {
					goto l1821
				}
				add(ruleIs, position1822)
			}
			return true
		l1821:
			position, tokenIndex = position1821, tokenIndex1821
			return false
		},
		/* 158 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action127)> */
		func() bool {
			position1828, tokenIndex1828 := position, tokenIndex
			{
				position1829 := position
				{
					position1830 := position
					{
						position1831, tokenIndex1831 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1832
						}
						position++
						goto l1831
					l1832:
						position, tokenIndex = position1831, tokenIndex1831
						if buffer[position] != rune('I') {
							goto l1828
						}
						position++
					}
				l1831:
					{
						position1833, tokenIndex1833 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1834
						}
						position++
						goto l1833
					l1834:
						position, tokenIndex = position1833, tokenIndex1833
						if buffer[position] != rune('S') {
							goto l1828
						}
						position++
					}
				l1833:
					if !_rules[rulesp]() {
						goto l1828
					}
					{
						position1835, tokenIndex1835 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1836
						}
						position++
						goto l1835
					l1836:
						position, tokenIndex = position1835, tokenIndex1835
						if buffer[position] != rune('N') {
							goto l1828
						}
						position++
					}
				l1835:
					{
						position1837, tokenIndex1837 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1838
						}
						position++
						goto l1837
					l1838:
						position, tokenIndex = position1837, tokenIndex1837
						if buffer[position] != rune('O') {
							goto l1828
						}
						position++
					}
				l1837:
					{
						position1839, tokenIndex1839 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1840
						}
						position++
						goto l1839
					l1840:
						position, tokenIndex = position1839, tokenIndex1839
						if buffer[position] != rune('T') {
							goto l1828
						}
						position++
					}
				l1839:
					add(rulePegText, position1830)
				}
				if !_rules[ruleAction127]() {
					goto l1828
				}
				add(ruleIsNot, position1829)
			}
			return true
		l1828:
			position, tokenIndex = position1828, tokenIndex1828
			return false
		},
		/* 159 Plus <- <(<'+'> Action128)> */
		func() bool {
			position1841, tokenIndex1841 := position, tokenIndex
			{
				position1842 := position
				{
					position1843 := position
					if buffer[position] != rune('+') {
						goto l1841
					}
					position++
					add(rulePegText, position1843)
				}
				if !_rules[ruleAction128]() {
					goto l1841
				}
				add(rulePlus, position1842)
			}
			return true
		l1841:
			position, tokenIndex = position1841, tokenIndex1841
			return false
		},
		/* 160 Minus <- <(<'-'> Action129)> */
		func() bool {
			position1844, tokenIndex1844 := position, tokenIndex
			{
				position1845 := position
				{
					position1846 := position
					if buffer[position] != rune('-') {
						goto l1844
					}
					position++
					add(rulePegText, position1846)
				}
				if !_rules[ruleAction129]() {
					goto l1844
				}
				add(ruleMinus, position1845)
			}
			return true
		l1844:
			position, tokenIndex = position1844, tokenIndex1844
			return false
		},
		/* 161 Multiply <- <(<'*'> Action130)> */
		func() bool {
			position1847, tokenIndex1847 := position, tokenIndex
			{
				position1848 := position
				{
					position1849 := position
					if buffer[position] != rune('*') {
						goto l1847
					}
					position++
					add(rulePegText, position1849)
				}
				if !_rules[ruleAction130]() {
					goto l1847
				}
				add(ruleMultiply, position1848)
			}
			return true
		l1847:
			position, tokenIndex = position1847, tokenIndex1847
			return false
		},
		/* 162 Divide <- <(<'/'> Action131)> */
		func() bool {
			position1850, tokenIndex1850 := position, tokenIndex
			{
				position1851 := position
				{
					position1852 := position
					if buffer[position] != rune('/') {
						goto l1850
					}
					position++
					add(rulePegText, position1852)
				}
				if !_rules[ruleAction131]() {
					goto l1850
				}
				add(ruleDivide, position1851)
			}
			return true
		l1850:
			position, tokenIndex = position1850, tokenIndex1850
			return false
		},
		/* 163 Modulo <- <(<'%'> Action132)> */
		func() bool {
			position1853, tokenIndex1853 := position, tokenIndex
			{
				position1854 := position
				{
					position1855 := position
					if buffer[position] != rune('%') {
						goto l1853
					}
					position++
					add(rulePegText, position1855)
				}
				if !_rules[ruleAction132]() {
					goto l1853
				}
				add(ruleModulo, position1854)
			}
			return true
		l1853:
			position, tokenIndex = position1853, tokenIndex1853
			return false
		},
		/* 164 UnaryMinus <- <(<'-'> Action133)> */
		func() bool {
			position1856, tokenIndex1856 := position, tokenIndex
			{
				position1857 := position
				{
					position1858 := position
					if buffer[position] != rune('-') {
						goto l1856
					}
					position++
					add(rulePegText, position1858)
				}
				if !_rules[ruleAction133]() {
					goto l1856
				}
				add(ruleUnaryMinus, position1857)
			}
			return true
		l1856:
			position, tokenIndex = position1856, tokenIndex1856
			return false
		},
		/* 165 Identifier <- <(<ident> Action134)> */
		func() bool {
			position1859, tokenIndex1859 := position, tokenIndex
			{
				position1860 := position
				{
					position1861 := position
					if !_rules[ruleident]() {
						goto l1859
					}
					add(rulePegText, position1861)
				}
				if !_rules[ruleAction134]() {
					goto l1859
				}
				add(ruleIdentifier, position1860)
			}
			return true
		l1859:
			position, tokenIndex = position1859, tokenIndex1859
			return false
		},
		/* 166 TargetIdentifier <- <(<('*' / jsonSetPath)> Action135)> */
		func() bool {
			position1862, tokenIndex1862 := position, tokenIndex
			{
				position1863 := position
				{
					position1864 := position
					{
						position1865, tokenIndex1865 := position, tokenIndex
						if buffer[position] != rune('*') {
							goto l1866
						}
						position++
						goto l1865
					l1866:
						position, tokenIndex = position1865, tokenIndex1865
						if !_rules[rulejsonSetPath]() {
							goto l1862
						}
					}
				l1865:
					add(rulePegText, position1864)
				}
				if !_rules[ruleAction135]() {
					goto l1862
				}
				add(ruleTargetIdentifier, position1863)
			}
			return true
		l1862:
			position, tokenIndex = position1862, tokenIndex1862
			return false
		},
		/* 167 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1867, tokenIndex1867 := position, tokenIndex
			{
				position1868 := position
				{
					position1869, tokenIndex1869 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1870
					}
					position++
					goto l1869
				l1870:
					position, tokenIndex = position1869, tokenIndex1869
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1867
					}
					position++
				}
			l1869:
			l1871:
				{
					position1872, tokenIndex1872 := position, tokenIndex
					{
						position1873, tokenIndex1873 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1874
						}
						position++
						goto l1873
					l1874:
						position, tokenIndex = position1873, tokenIndex1873
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1875
						}
						position++
						goto l1873
					l1875:
						position, tokenIndex = position1873, tokenIndex1873
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1876
						}
						position++
						goto l1873
					l1876:
						position, tokenIndex = position1873, tokenIndex1873
						if buffer[position] != rune('_') {
							goto l1872
						}
						position++
					}
				l1873:
					goto l1871
				l1872:
					position, tokenIndex = position1872, tokenIndex1872
				}
				add(ruleident, position1868)
			}
			return true
		l1867:
			position, tokenIndex = position1867, tokenIndex1867
			return false
		},
		/* 168 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1877, tokenIndex1877 := position, tokenIndex
			{
				position1878 := position
				if !_rules[rulejsonPathHead]() {
					goto l1877
				}
			l1879:
				{
					position1880, tokenIndex1880 := position, tokenIndex
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1880
					}
					goto l1879
				l1880:
					position, tokenIndex = position1880, tokenIndex1880
				}
				add(rulejsonGetPath, position1878)
			}
			return true
		l1877:
			position, tokenIndex = position1877, tokenIndex1877
			return false
		},
		/* 169 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1881, tokenIndex1881 := position, tokenIndex
			{
				position1882 := position
				if !_rules[rulejsonPathHead]() {
					goto l1881
				}
			l1883:
				{
					position1884, tokenIndex1884 := position, tokenIndex
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1884
					}
					goto l1883
				l1884:
					position, tokenIndex = position1884, tokenIndex1884
				}
				add(rulejsonSetPath, position1882)
			}
			return true
		l1881:
			position, tokenIndex = position1881, tokenIndex1881
			return false
		},
		/* 170 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1885, tokenIndex1885 := position, tokenIndex
			{
				position1886 := position
				{
					position1887, tokenIndex1887 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l1888
					}
					goto l1887
				l1888:
					position, tokenIndex = position1887, tokenIndex1887
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1885
					}
				}
			l1887:
				add(rulejsonPathHead, position1886)
			}
			return true
		l1885:
			position, tokenIndex = position1885, tokenIndex1885
			return false
		},
		/* 171 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1889, tokenIndex1889 := position, tokenIndex
			{
				position1890 := position
				{
					position1891, tokenIndex1891 := position, tokenIndex
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1892
					}
					goto l1891
				l1892:
					position, tokenIndex = position1891, tokenIndex1891
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1893
					}
					goto l1891
				l1893:
					position, tokenIndex = position1891, tokenIndex1891
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1894
					}
					goto l1891
				l1894:
					position, tokenIndex = position1891, tokenIndex1891
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1895
					}
					goto l1891
				l1895:
					position, tokenIndex = position1891, tokenIndex1891
					if !_rules[rulejsonArraySlice]() {
						goto l1896
					}
					goto l1891
				l1896:
					position, tokenIndex = position1891, tokenIndex1891
					if !_rules[rulejsonArrayAccess]() {
						goto l1889
					}
				}
			l1891:
				add(rulejsonGetPathNonHead, position1890)
			}
			return true
		l1889:
			position, tokenIndex = position1889, tokenIndex1889
			return false
		},
		/* 172 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1897, tokenIndex1897 := position, tokenIndex
			{
				position1898 := position
				{
					position1899, tokenIndex1899 := position, tokenIndex
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1900
					}
					goto l1899
				l1900:
					position, tokenIndex = position1899, tokenIndex1899
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1897
					}
				}
			l1899:
				add(rulejsonSetPathNonHead, position1898)
			}
			return true
		l1897:
			position, tokenIndex = position1897, tokenIndex1897
			return false
		},
		/* 173 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1901, tokenIndex1901 := position, tokenIndex
			{
				position1902 := position
				{
					position1903, tokenIndex1903 := position, tokenIndex
					if buffer[position] != rune('.') {
						goto l1904
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1904
					}
					goto l1903
				l1904:
					position, tokenIndex = position1903, tokenIndex1903
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1901
					}
				}
			l1903:
				add(rulejsonMapSingleLevel, position1902)
			}
			return true
		l1901:
			position, tokenIndex = position1901, tokenIndex1901
			return false
		},
		/* 174 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1905, tokenIndex1905 := position, tokenIndex
			{
				position1906 := position
				if buffer[position] != rune('.') {
					goto l1905
				}
				position++
				if buffer[position] != rune('.') {
					goto l1905
				}
				position++
				{
					position1907, tokenIndex1907 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l1908
					}
					goto l1907
				l1908:
					position, tokenIndex = position1907, tokenIndex1907
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1905
					}
				}
			l1907:
				add(rulejsonMapMultipleLevel, position1906)
			}
			return true
		l1905:
			position, tokenIndex = position1905, tokenIndex1905
			return false
		},
		/* 175 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1909, tokenIndex1909 := position, tokenIndex
			{
				position1910 := position
				{
					position1911 := position
					{
						position1912, tokenIndex1912 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1913
						}
						position++
						goto l1912
					l1913:
						position, tokenIndex = position1912, tokenIndex1912
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1909
						}
						position++
					}
				l1912:
				l1914:
					{
						position1915, tokenIndex1915 := position, tokenIndex
						{
							position1916, tokenIndex1916 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1917
							}
							position++
							goto l1916
						l1917:
							position, tokenIndex = position1916, tokenIndex1916
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1918
							}
							position++
							goto l1916
						l1918:
							position, tokenIndex = position1916, tokenIndex1916
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1919
							}
							position++
							goto l1916
						l1919:
							position, tokenIndex = position1916, tokenIndex1916
							if buffer[position] != rune('_') {
								goto l1915
							}
							position++
						}
					l1916:
						goto l1914
					l1915:
						position, tokenIndex = position1915, tokenIndex1915
					}
					add(rulePegText, position1911)
				}
				add(rulejsonMapAccessString, position1910)
			}
			return true
		l1909:
			position, tokenIndex = position1909, tokenIndex1909
			return false
		},
		/* 176 jsonMapAccessBracket <- <('[' doubleQuotedString ']')> */
		func() bool {
			position1920, tokenIndex1920 := position, tokenIndex
			{
				position1921 := position
				if buffer[position] != rune('[') {
					goto l1920
				}
				position++
				if !_rules[ruledoubleQuotedString]() {
					goto l1920
				}
				if buffer[position] != rune(']') {
					goto l1920
				}
				position++
				add(rulejsonMapAccessBracket, position1921)
			}
			return true
		l1920:
			position, tokenIndex = position1920, tokenIndex1920
			return false
		},
		/* 177 doubleQuotedString <- <('"' <(('"' '"') / (!'"' .))*> '"')> */
		func() bool {
			position1922, tokenIndex1922 := position, tokenIndex
			{
				position1923 := position
				if buffer[position] != rune('"') {
					goto l1922
				}
				position++
				{
					position1924 := position
				l1925:
					{
						position1926, tokenIndex1926 := position, tokenIndex
						{
							position1927, tokenIndex1927 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l1928
							}
							position++
							if buffer[position] != rune('"') {
								goto l1928
							}
							position++
							goto l1927
						l1928:
							position, tokenIndex = position1927, tokenIndex1927
							{
								position1929, tokenIndex1929 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l1929
								}
								position++
								goto l1926
							l1929:
								position, tokenIndex = position1929, tokenIndex1929
							}
							if !matchDot() {
								goto l1926
							}
						}
					l1927:
						goto l1925
					l1926:
						position, tokenIndex = position1926, tokenIndex1926
					}
					add(rulePegText, position1924)
				}
				if buffer[position] != rune('"') {
					goto l1922
				}
				position++
				add(ruledoubleQuotedString, position1923)
			}
			return true
		l1922:
			position, tokenIndex = position1922, tokenIndex1922
			return false
		},
		/* 178 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1930, tokenIndex1930 := position, tokenIndex
			{
				position1931 := position
				if buffer[position] != rune('[') {
					goto l1930
				}
				position++
				{
					position1932 := position
					{
						position1933, tokenIndex1933 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1933
						}
						position++
						goto l1934
					l1933:
						position, tokenIndex = position1933, tokenIndex1933
					}
				l1934:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1930
					}
					position++
				l1935:
					{
						position1936, tokenIndex1936 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1936
						}
						position++
						goto l1935
					l1936:
						position, tokenIndex = position1936, tokenIndex1936
					}
					add(rulePegText, position1932)
				}
				if buffer[position] != rune(']') {
					goto l1930
				}
				position++
				add(rulejsonArrayAccess, position1931)
			}
			return true
		l1930:
			position, tokenIndex = position1930, tokenIndex1930
			return false
		},
		/* 179 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
		func() bool {
			position1937, tokenIndex1937 := position, tokenIndex
			{
				position1938 := position
				if buffer[position] != rune('[') {
					goto l1937
				}
				position++
				{
					position1939 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1937
					}
					position++
				l1940:
					{
						position1941, tokenIndex1941 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1941
						}
						position++
						goto l1940
					l1941:
						position, tokenIndex = position1941, tokenIndex1941
					}
					add(rulePegText, position1939)
				}
				if buffer[position] != rune(']') {
					goto l1937
				}
				position++
				add(rulejsonNonNegativeArrayAccess, position1938)
			}
			return true
		l1937:
			position, tokenIndex = position1937, tokenIndex1937
			return false
		},
		/* 180 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1942, tokenIndex1942 := position, tokenIndex
			{
				position1943 := position
				if buffer[position] != rune('[') {
					goto l1942
				}
				position++
				{
					position1944 := position
					{
						position1945, tokenIndex1945 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1945
						}
						position++
						goto l1946
					l1945:
						position, tokenIndex = position1945, tokenIndex1945
					}
				l1946:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1942
					}
					position++
				l1947:
					{
						position1948, tokenIndex1948 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1948
						}
						position++
						goto l1947
					l1948:
						position, tokenIndex = position1948, tokenIndex1948
					}
					if buffer[position] != rune(':') {
						goto l1942
					}
					position++
					{
						position1949, tokenIndex1949 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1949
						}
						position++
						goto l1950
					l1949:
						position, tokenIndex = position1949, tokenIndex1949
					}
				l1950:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1942
					}
					position++
				l1951:
					{
						position1952, tokenIndex1952 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1952
						}
						position++
						goto l1951
					l1952:
						position, tokenIndex = position1952, tokenIndex1952
					}
					{
						position1953, tokenIndex1953 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l1953
						}
						position++
						{
							position1955, tokenIndex1955 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1955
							}
							position++
							goto l1956
						l1955:
							position, tokenIndex = position1955, tokenIndex1955
						}
					l1956:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1953
						}
						position++
					l1957:
						{
							position1958, tokenIndex1958 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1958
							}
							position++
							goto l1957
						l1958:
							position, tokenIndex = position1958, tokenIndex1958
						}
						goto l1954
					l1953:
						position, tokenIndex = position1953, tokenIndex1953
					}
				l1954:
					add(rulePegText, position1944)
				}
				if buffer[position] != rune(']') {
					goto l1942
				}
				position++
				add(rulejsonArraySlice, position1943)
			}
			return true
		l1942:
			position, tokenIndex = position1942, tokenIndex1942
			return false
		},
		/* 181 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1959, tokenIndex1959 := position, tokenIndex
			{
				position1960 := position
				if buffer[position] != rune('[') {
					goto l1959
				}
				position++
				{
					position1961 := position
					{
						position1962, tokenIndex1962 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l1963
						}
						position++
						{
							position1964, tokenIndex1964 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1964
							}
							position++
							goto l1965
						l1964:
							position, tokenIndex = position1964, tokenIndex1964
						}
					l1965:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1963
						}
						position++
					l1966:
						{
							position1967, tokenIndex1967 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1967
							}
							position++
							goto l1966
						l1967:
							position, tokenIndex = position1967, tokenIndex1967
						}
						goto l1962
					l1963:
						position, tokenIndex = position1962, tokenIndex1962
						{
							position1968, tokenIndex1968 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1968
							}
							position++
							goto l1969
						l1968:
							position, tokenIndex = position1968, tokenIndex1968
						}
					l1969:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1959
						}
						position++
					l1970:
						{
							position1971, tokenIndex1971 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1971
							}
							position++
							goto l1970
						l1971:
							position, tokenIndex = position1971, tokenIndex1971
						}
						if buffer[position] != rune(':') {
							goto l1959
						}
						position++
					}
				l1962:
					add(rulePegText, position1961)
				}
				if buffer[position] != rune(']') {
					goto l1959
				}
				position++
				add(rulejsonArrayPartialSlice, position1960)
			}
			return true
		l1959:
			position, tokenIndex = position1959, tokenIndex1959
			return false
		},
		/* 182 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1972, tokenIndex1972 := position, tokenIndex
			{
				position1973 := position
				if buffer[position] != rune('[') {
					goto l1972
				}
				position++
				if buffer[position] != rune(':') {
					goto l1972
				}
				position++
				if buffer[position] != rune(']') {
					goto l1972
				}
				position++
				add(rulejsonArrayFullSlice, position1973)
			}
			return true
		l1972:
			position, tokenIndex = position1972, tokenIndex1972
			return false
		},
		/* 183 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1974, tokenIndex1974 := position, tokenIndex
			{
				position1975 := position
				{
					position1976, tokenIndex1976 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l1977
					}
					position++
					goto l1976
				l1977:
					position, tokenIndex = position1976, tokenIndex1976
					if buffer[position] != rune('\t') {
						goto l1978
					}
					position++
					goto l1976
				l1978:
					position, tokenIndex = position1976, tokenIndex1976
					if buffer[position] != rune('\n') {
						goto l1979
					}
					position++
					goto l1976
				l1979:
					position, tokenIndex = position1976, tokenIndex1976
					if buffer[position] != rune('\r') {
						goto l1980
					}
					position++
					goto l1976
				l1980:
					position, tokenIndex = position1976, tokenIndex1976
					if !_rules[rulecomment]() {
						goto l1981
					}
					goto l1976
				l1981:
					position, tokenIndex = position1976, tokenIndex1976
					if !_rules[rulefinalComment]() {
						goto l1974
					}
				}
			l1976:
				add(rulespElem, position1975)
			}
			return true
		l1974:
			position, tokenIndex = position1974, tokenIndex1974
			return false
		},
		/* 184 sp <- <spElem+> */
		func() bool {
			position1982, tokenIndex1982 := position, tokenIndex
			{
				position1983 := position
				if !_rules[rulespElem]() {
					goto l1982
				}
			l1984:
				{
					position1985, tokenIndex1985 := position, tokenIndex
					if !_rules[rulespElem]() {
						goto l1985
					}
					goto l1984
				l1985:
					position, tokenIndex = position1985, tokenIndex1985
				}
				add(rulesp, position1983)
			}
			return true
		l1982:
			position, tokenIndex = position1982, tokenIndex1982
			return false
		},
		/* 185 spOpt <- <spElem*> */
		func() bool {
			{
				position1987 := position
			l1988:
				{
					position1989, tokenIndex1989 := position, tokenIndex
					if !_rules[rulespElem]() {
						goto l1989
					}
					goto l1988
				l1989:
					position, tokenIndex = position1989, tokenIndex1989
				}
				add(rulespOpt, position1987)
			}
			return true
		},
		/* 186 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1990, tokenIndex1990 := position, tokenIndex
			{
				position1991 := position
				if buffer[position] != rune('-') {
					goto l1990
				}
				position++
				if buffer[position] != rune('-') {
					goto l1990
				}
				position++
			l1992:
				{
					position1993, tokenIndex1993 := position, tokenIndex
					{
						position1994, tokenIndex1994 := position, tokenIndex
						{
							position1995, tokenIndex1995 := position, tokenIndex
							if buffer[position] != rune('\r') {
								goto l1996
							}
							position++
							goto l1995
						l1996:
							position, tokenIndex = position1995, tokenIndex1995
							if buffer[position] != rune('\n') {
								goto l1994
							}
							position++
						}
					l1995:
						goto l1993
					l1994:
						position, tokenIndex = position1994, tokenIndex1994
					}
					if !matchDot() {
						goto l1993
					}
					goto l1992
				l1993:
					position, tokenIndex = position1993, tokenIndex1993
				}
				{
					position1997, tokenIndex1997 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l1998
					}
					position++
					goto l1997
				l1998:
					position, tokenIndex = position1997, tokenIndex1997
					if buffer[position] != rune('\n') {
						goto l1990
					}
					position++
				}
			l1997:
				add(rulecomment, position1991)
			}
			return true
		l1990:
			position, tokenIndex = position1990, tokenIndex1990
			return false
		},
		/* 187 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1999, tokenIndex1999 := position, tokenIndex
			{
				position2000 := position
				if buffer[position] != rune('-') {
					goto l1999
				}
				position++
				if buffer[position] != rune('-') {
					goto l1999
				}
				position++
			l2001:
				{
					position2002, tokenIndex2002 := position, tokenIndex
					{
						position2003, tokenIndex2003 := position, tokenIndex
						{
							position2004, tokenIndex2004 := position, tokenIndex
							if buffer[position] != rune('\r') {
								goto l2005
							}
							position++
							goto l2004
						l2005:
							position, tokenIndex = position2004, tokenIndex2004
							if buffer[position] != rune('\n') {
								goto l2003
							}
							position++
						}
					l2004:
						goto l2002
					l2003:
						position, tokenIndex = position2003, tokenIndex2003
					}
					if !matchDot() {
						goto l2002
					}
					goto l2001
				l2002:
					position, tokenIndex = position2002, tokenIndex2002
				}
				{
					position2006, tokenIndex2006 := position, tokenIndex
					if !matchDot() {
						goto l2006
					}
					goto l1999
				l2006:
					position, tokenIndex = position2006, tokenIndex2006
				}
				add(rulefinalComment, position2000)
			}
			return true
		l1999:
			position, tokenIndex = position1999, tokenIndex1999
			return false
		},
		nil,
		/* 190 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 191 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 192 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 193 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 194 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 195 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 196 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 197 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 198 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 199 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 200 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 201 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 202 Action12 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 203 Action13 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 204 Action14 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 205 Action15 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 206 Action16 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 207 Action17 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 208 Action18 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 209 Action19 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 210 Action20 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 211 Action21 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 212 Action22 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 213 Action23 <- <{
		    p.AssembleEval(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 214 Action24 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 215 Action25 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 216 Action26 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 217 Action27 <- <{
		    p.AssembleEmitterSampling(CountBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 218 Action28 <- <{
		    p.AssembleEmitterSampling(RandomizedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 219 Action29 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 220 Action30 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 0.001)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 221 Action31 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 222 Action32 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 223 Action33 <- <{
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
		/* 224 Action34 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 225 Action35 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 226 Action36 <- <{
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
		/* 227 Action37 <- <{
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
		/* 228 Action38 <- <{
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
		/* 229 Action39 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 230 Action40 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 231 Action41 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 232 Action42 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 233 Action43 <- <{
		    p.EnsureCapacitySpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 234 Action44 <- <{
		    p.EnsureSheddingSpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 235 Action45 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 236 Action46 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 237 Action47 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 238 Action48 <- <{
		    p.EnsureIdentifier(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 239 Action49 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 240 Action50 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 241 Action51 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 242 Action52 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 243 Action53 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 244 Action54 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 245 Action55 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 246 Action56 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 247 Action57 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 248 Action58 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 249 Action59 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 250 Action60 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 251 Action61 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 252 Action62 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 253 Action63 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 254 Action64 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 255 Action65 <- <{
		    p.AssembleFuncAppSelector()
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 256 Action66 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRaw(substr))
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 257 Action67 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 258 Action68 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 259 Action69 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 260 Action70 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 261 Action71 <- <{
		    p.AssembleSortedExpression()
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 262 Action72 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 263 Action73 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 264 Action74 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 265 Action75 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 266 Action76 <- <{
		    p.AssembleConditionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 267 Action77 <- <{
		    p.AssembleExpressionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 268 Action78 <- <{
		    p.AssembleWhenThenPair()
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 269 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 270 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 271 Action81 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 272 Action82 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 273 Action83 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 274 Action84 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 275 Action85 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 276 Action86 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 277 Action87 <- <{
		    p.PushComponent(begin, end, NewMissing())
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 278 Action88 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 279 Action89 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 280 Action90 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 281 Action91 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 282 Action92 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 283 Action93 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 284 Action94 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 285 Action95 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 286 Action96 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 287 Action97 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 288 Action98 <- <{
		    p.PushComponent(begin, end, Wait)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 289 Action99 <- <{
		    p.PushComponent(begin, end, DropOldest)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 290 Action100 <- <{
		    p.PushComponent(begin, end, DropNewest)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 291 Action101 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 292 Action102 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 293 Action103 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 294 Action104 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 295 Action105 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 296 Action106 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 297 Action107 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 298 Action108 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 299 Action109 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 300 Action110 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction110, position)
			}
			return true
		},
		/* 301 Action111 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction111, position)
			}
			return true
		},
		/* 302 Action112 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction112, position)
			}
			return true
		},
		/* 303 Action113 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction113, position)
			}
			return true
		},
		/* 304 Action114 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction114, position)
			}
			return true
		},
		/* 305 Action115 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction115, position)
			}
			return true
		},
		/* 306 Action116 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction116, position)
			}
			return true
		},
		/* 307 Action117 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction117, position)
			}
			return true
		},
		/* 308 Action118 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction118, position)
			}
			return true
		},
		/* 309 Action119 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction119, position)
			}
			return true
		},
		/* 310 Action120 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction120, position)
			}
			return true
		},
		/* 311 Action121 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction121, position)
			}
			return true
		},
		/* 312 Action122 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction122, position)
			}
			return true
		},
		/* 313 Action123 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction123, position)
			}
			return true
		},
		/* 314 Action124 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction124, position)
			}
			return true
		},
		/* 315 Action125 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction125, position)
			}
			return true
		},
		/* 316 Action126 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction126, position)
			}
			return true
		},
		/* 317 Action127 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction127, position)
			}
			return true
		},
		/* 318 Action128 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction128, position)
			}
			return true
		},
		/* 319 Action129 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction129, position)
			}
			return true
		},
		/* 320 Action130 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction130, position)
			}
			return true
		},
		/* 321 Action131 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction131, position)
			}
			return true
		},
		/* 322 Action132 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction132, position)
			}
			return true
		},
		/* 323 Action133 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction133, position)
			}
			return true
		},
		/* 324 Action134 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction134, position)
			}
			return true
		},
		/* 325 Action135 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction135, position)
			}
			return true
		},
	}
	p.rules = _rules
}

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
	rules  [322]func() bool
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

			p.AssembleFuncApp()

		case ruleAction66:

			p.AssembleExpressions(begin, end)
			p.AssembleFuncApp()

		case ruleAction67:

			p.AssembleExpressions(begin, end)

		case ruleAction68:

			p.AssembleExpressions(begin, end)

		case ruleAction69:

			p.AssembleSortedExpression()

		case ruleAction70:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction71:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction72:

			p.AssembleMap(begin, end)

		case ruleAction73:

			p.AssembleKeyValuePair()

		case ruleAction74:

			p.AssembleConditionCase(begin, end)

		case ruleAction75:

			p.AssembleExpressionCase(begin, end)

		case ruleAction76:

			p.AssembleWhenThenPair()

		case ruleAction77:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction78:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction80:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction81:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction82:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction83:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction84:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction85:

			p.PushComponent(begin, end, NewMissing())

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
		/* 65 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr / ParamMapExpr)> */
		func() bool {
			position1021, tokenIndex1021 := position, tokenIndex
			{
				position1022 := position
				{
					position1023, tokenIndex1023 := position, tokenIndex
					if !_rules[ruleParamLiteral]() {
						goto l1024
					}
					goto l1023
				l1024:
					position, tokenIndex = position1023, tokenIndex1023
					if !_rules[ruleParamArrayExpr]() {
						goto l1025
					}
					goto l1023
				l1025:
					position, tokenIndex = position1023, tokenIndex1023
					if !_rules[ruleParamMapExpr]() {
						goto l1021
					}
				}
			l1023:
				add(ruleSourceSinkParamVal, position1022)
			}
			return true
		l1021:
			position, tokenIndex = position1021, tokenIndex1021
			return false
		},
		/* 66 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position1026, tokenIndex1026 := position, tokenIndex
			{
				position1027 := position
				{
					position1028, tokenIndex1028 := position, tokenIndex
					if !_rules[ruleBooleanLiteral]() {
						goto l1029
					}
					goto l1028
				l1029:
					position, tokenIndex = position1028, tokenIndex1028
					if !_rules[ruleLiteral]() {
						goto l1026
					}
				}
			l1028:
				add(ruleParamLiteral, position1027)
			}
			return true
		l1026:
			position, tokenIndex = position1026, tokenIndex1026
			return false
		},
		/* 67 ParamArrayExpr <- <(<('[' spOpt (ParamLiteral (',' spOpt ParamLiteral)*)? spOpt ','? spOpt ']')> Action50)> */
		func() bool {
			position1030, tokenIndex1030 := position, tokenIndex
			{
				position1031 := position
				{
					position1032 := position
					if buffer[position] != rune('[') {
						goto l1030
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1030
					}
					{
						position1033, tokenIndex1033 := position, tokenIndex
						if !_rules[ruleParamLiteral]() {
							goto l1033
						}
					l1035:
						{
							position1036, tokenIndex1036 := position, tokenIndex
							if buffer[position] != rune(',') {
								goto l1036
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1036
							}
							if !_rules[ruleParamLiteral]() {
								goto l1036
							}
							goto l1035
						l1036:
							position, tokenIndex = position1036, tokenIndex1036
						}
						goto l1034
					l1033:
						position, tokenIndex = position1033, tokenIndex1033
					}
				l1034:
					if !_rules[rulespOpt]() {
						goto l1030
					}
					{
						position1037, tokenIndex1037 := position, tokenIndex
						if buffer[position] != rune(',') {
							goto l1037
						}
						position++
						goto l1038
					l1037:
						position, tokenIndex = position1037, tokenIndex1037
					}
				l1038:
					if !_rules[rulespOpt]() {
						goto l1030
					}
					if buffer[position] != rune(']') {
						goto l1030
					}
					position++
					add(rulePegText, position1032)
				}
				if !_rules[ruleAction50]() {
					goto l1030
				}
				add(ruleParamArrayExpr, position1031)
			}
			return true
		l1030:
			position, tokenIndex = position1030, tokenIndex1030
			return false
		},
		/* 68 ParamMapExpr <- <(<('{' spOpt (ParamKeyValuePair (spOpt ',' spOpt ParamKeyValuePair)*)? spOpt '}')> Action51)> */
		func() bool {
			position1039, tokenIndex1039 := position, tokenIndex
			{
				position1040 := position
				{
					position1041 := position
					if buffer[position] != rune('{') {
						goto l1039
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1039
					}
					{
						position1042, tokenIndex1042 := position, tokenIndex
						if !_rules[ruleParamKeyValuePair]() {
							goto l1042
						}
					l1044:
						{
							position1045, tokenIndex1045 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1045
							}
							if buffer[position] != rune(',') {
								goto l1045
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1045
							}
							if !_rules[ruleParamKeyValuePair]() {
								goto l1045
							}
							goto l1044
						l1045:
							position, tokenIndex = position1045, tokenIndex1045
						}
						goto l1043
					l1042:
						position, tokenIndex = position1042, tokenIndex1042
					}
				l1043:
					if !_rules[rulespOpt]() {
						goto l1039
					}
					if buffer[position] != rune('}') {
						goto l1039
					}
					position++
					add(rulePegText, position1041)
				}
				if !_rules[ruleAction51]() {
					goto l1039
				}
				add(ruleParamMapExpr, position1040)
			}
			return true
		l1039:
			position, tokenIndex = position1039, tokenIndex1039
			return false
		},
		/* 69 ParamKeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ParamLiteral)> Action52)> */
		func() bool {
			position1046, tokenIndex1046 := position, tokenIndex
			{
				position1047 := position
				{
					position1048 := position
					if !_rules[ruleStringLiteral]() {
						goto l1046
					}
					if !_rules[rulespOpt]() {
						goto l1046
					}
					if buffer[position] != rune(':') {
						goto l1046
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1046
					}
					if !_rules[ruleParamLiteral]() {
						goto l1046
					}
					add(rulePegText, position1048)
				}
				if !_rules[ruleAction52]() {
					goto l1046
				}
				add(ruleParamKeyValuePair, position1047)
			}
			return true
		l1046:
			position, tokenIndex = position1046, tokenIndex1046
			return false
		},
		/* 70 PausedOpt <- <(<(sp (Paused / Unpaused))?> Action53)> */
		func() bool {
			position1049, tokenIndex1049 := position, tokenIndex
			{
				position1050 := position
				{
					position1051 := position
					{
						position1052, tokenIndex1052 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1052
						}
						{
							position1054, tokenIndex1054 := position, tokenIndex
							if !_rules[rulePaused]() {
								goto l1055
							}
							goto l1054
						l1055:
							position, tokenIndex = position1054, tokenIndex1054
							if !_rules[ruleUnpaused]() {
								goto l1052
							}
						}
					l1054:
						goto l1053
					l1052:
						position, tokenIndex = position1052, tokenIndex1052
					}
				l1053:
					add(rulePegText, position1051)
				}
				if !_rules[ruleAction53]() {
					goto l1049
				}
				add(rulePausedOpt, position1050)
			}
			return true
		l1049:
			position, tokenIndex = position1049, tokenIndex1049
			return false
		},
		/* 71 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position1056, tokenIndex1056 := position, tokenIndex
			{
				position1057 := position
				{
					position1058, tokenIndex1058 := position, tokenIndex
					if !_rules[ruleWildcard]() {
						goto l1059
					}
					goto l1058
				l1059:
					position, tokenIndex = position1058, tokenIndex1058
					if !_rules[ruleExpression]() {
						goto l1056
					}
				}
			l1058:
				add(ruleExpressionOrWildcard, position1057)
			}
			return true
		l1056:
			position, tokenIndex = position1056, tokenIndex1056
			return false
		},
		/* 72 Expression <- <orExpr> */
		func() bool {
			position1060, tokenIndex1060 := position, tokenIndex
			{
				position1061 := position
				if !_rules[ruleorExpr]() {
					goto l1060
				}
				add(ruleExpression, position1061)
			}
			return true
		l1060:
			position, tokenIndex = position1060, tokenIndex1060
			return false
		},
		/* 73 orExpr <- <(<(andExpr (sp Or sp andExpr)*)> Action54)> */
		func() bool {
			position1062, tokenIndex1062 := position, tokenIndex
			{
				position1063 := position
				{
					position1064 := position
					if !_rules[ruleandExpr]() {
						goto l1062
					}
				l1065:
					{
						position1066, tokenIndex1066 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1066
						}
						if !_rules[ruleOr]() {
							goto l1066
						}
						if !_rules[rulesp]() {
							goto l1066
						}
						if !_rules[ruleandExpr]() {
							goto l1066
						}
						goto l1065
					l1066:
						position, tokenIndex = position1066, tokenIndex1066
					}
					add(rulePegText, position1064)
				}
				if !_rules[ruleAction54]() {
					goto l1062
				}
				add(ruleorExpr, position1063)
			}
			return true
		l1062:
			position, tokenIndex = position1062, tokenIndex1062
			return false
		},
		/* 74 andExpr <- <(<(notExpr (sp And sp notExpr)*)> Action55)> */
		func() bool {
			position1067, tokenIndex1067 := position, tokenIndex
			{
				position1068 := position
				{
					position1069 := position
					if !_rules[rulenotExpr]() {
						goto l1067
					}
				l1070:
					{
						position1071, tokenIndex1071 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1071
						}
						if !_rules[ruleAnd]() {
							goto l1071
						}
						if !_rules[rulesp]() {
							goto l1071
						}
						if !_rules[rulenotExpr]() {
							goto l1071
						}
						goto l1070
					l1071:
						position, tokenIndex = position1071, tokenIndex1071
					}
					add(rulePegText, position1069)
				}
				if !_rules[ruleAction55]() {
					goto l1067
				}
				add(ruleandExpr, position1068)
			}
			return true
		l1067:
			position, tokenIndex = position1067, tokenIndex1067
			return false
		},
		/* 75 notExpr <- <(<((Not sp)? comparisonExpr)> Action56)> */
		func() bool {
			position1072, tokenIndex1072 := position, tokenIndex
			{
				position1073 := position
				{
					position1074 := position
					{
						position1075, tokenIndex1075 := position, tokenIndex
						if !_rules[ruleNot]() {
							goto l1075
						}
						if !_rules[rulesp]() {
							goto l1075
						}
						goto l1076
					l1075:
						position, tokenIndex = position1075, tokenIndex1075
					}
				l1076:
					if !_rules[rulecomparisonExpr]() {
						goto l1072
					}
					add(rulePegText, position1074)
				}
				if !_rules[ruleAction56]() {
					goto l1072
				}
				add(rulenotExpr, position1073)
			}
			return true
		l1072:
			position, tokenIndex = position1072, tokenIndex1072
			return false
		},
		/* 76 comparisonExpr <- <(<(otherOpExpr (spOpt ComparisonOp spOpt otherOpExpr)?)> Action57)> */
		func() bool {
			position1077, tokenIndex1077 := position, tokenIndex
			{
				position1078 := position
				{
					position1079 := position
					if !_rules[ruleotherOpExpr]() {
						goto l1077
					}
					{
						position1080, tokenIndex1080 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1080
						}
						if !_rules[ruleComparisonOp]() {
							goto l1080
						}
						if !_rules[rulespOpt]() {
							goto l1080
						}
						if !_rules[ruleotherOpExpr]() {
							goto l1080
						}
						goto l1081
					l1080:
						position, tokenIndex = position1080, tokenIndex1080
					}
				l1081:
					add(rulePegText, position1079)
				}
				if !_rules[ruleAction57]() {
					goto l1077
				}
				add(rulecomparisonExpr, position1078)
			}
			return true
		l1077:
			position, tokenIndex = position1077, tokenIndex1077
			return false
		},
		/* 77 otherOpExpr <- <(<(isExpr (spOpt OtherOp spOpt isExpr)*)> Action58)> */
		func() bool {
			position1082, tokenIndex1082 := position, tokenIndex
			{
				position1083 := position
				{
					position1084 := position
					if !_rules[ruleisExpr]() {
						goto l1082
					}
				l1085:
					{
						position1086, tokenIndex1086 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1086
						}
						if !_rules[ruleOtherOp]() {
							goto l1086
						}
						if !_rules[rulespOpt]() {
							goto l1086
						}
						if !_rules[ruleisExpr]() {
							goto l1086
						}
						goto l1085
					l1086:
						position, tokenIndex = position1086, tokenIndex1086
					}
					add(rulePegText, position1084)
				}
				if !_rules[ruleAction58]() {
					goto l1082
				}
				add(ruleotherOpExpr, position1083)
			}
			return true
		l1082:
			position, tokenIndex = position1082, tokenIndex1082
			return false
		},
		/* 78 isExpr <- <(<((RowValue sp IsOp sp Missing) / (termExpr (sp IsOp sp NullLiteral)?))> Action59)> */
		func() bool {
			position1087, tokenIndex1087 := position, tokenIndex
			{
				position1088 := position
				{
					position1089 := position
					{
						position1090, tokenIndex1090 := position, tokenIndex
						if !_rules[ruleRowValue]() {
							goto l1091
						}
						if !_rules[rulesp]() {
							goto l1091
						}
						if !_rules[ruleIsOp]() {
							goto l1091
						}
						if !_rules[rulesp]() {
							goto l1091
						}
						if !_rules[ruleMissing]() {
							goto l1091
						}
						goto l1090
					l1091:
						position, tokenIndex = position1090, tokenIndex1090
						if !_rules[ruletermExpr]() {
							goto l1087
						}
						{
							position1092, tokenIndex1092 := position, tokenIndex
							if !_rules[rulesp]() {
								goto l1092
							}
							if !_rules[ruleIsOp]() {
								goto l1092
							}
							if !_rules[rulesp]() {
								goto l1092
							}
							if !_rules[ruleNullLiteral]() {
								goto l1092
							}
							goto l1093
						l1092:
							position, tokenIndex = position1092, tokenIndex1092
						}
					l1093:
					}
				l1090:
					add(rulePegText, position1089)
				}
				if !_rules[ruleAction59]() {
					goto l1087
				}
				add(ruleisExpr, position1088)
			}
			return true
		l1087:
			position, tokenIndex = position1087, tokenIndex1087
			return false
		},
		/* 79 termExpr <- <(<(productExpr (spOpt PlusMinusOp spOpt productExpr)*)> Action60)> */
		func() bool {
			position1094, tokenIndex1094 := position, tokenIndex
			{
				position1095 := position
				{
					position1096 := position
					if !_rules[ruleproductExpr]() {
						goto l1094
					}
				l1097:
					{
						position1098, tokenIndex1098 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1098
						}
						if !_rules[rulePlusMinusOp]() {
							goto l1098
						}
						if !_rules[rulespOpt]() {
							goto l1098
						}
						if !_rules[ruleproductExpr]() {
							goto l1098
						}
						goto l1097
					l1098:
						position, tokenIndex = position1098, tokenIndex1098
					}
					add(rulePegText, position1096)
				}
				if !_rules[ruleAction60]() {
					goto l1094
				}
				add(ruletermExpr, position1095)
			}
			return true
		l1094:
			position, tokenIndex = position1094, tokenIndex1094
			return false
		},
		/* 80 productExpr <- <(<(minusExpr (spOpt MultDivOp spOpt minusExpr)*)> Action61)> */
		func() bool {
			position1099, tokenIndex1099 := position, tokenIndex
			{
				position1100 := position
				{
					position1101 := position
					if !_rules[ruleminusExpr]() {
						goto l1099
					}
				l1102:
					{
						position1103, tokenIndex1103 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1103
						}
						if !_rules[ruleMultDivOp]() {
							goto l1103
						}
						if !_rules[rulespOpt]() {
							goto l1103
						}
						if !_rules[ruleminusExpr]() {
							goto l1103
						}
						goto l1102
					l1103:
						position, tokenIndex = position1103, tokenIndex1103
					}
					add(rulePegText, position1101)
				}
				if !_rules[ruleAction61]() {
					goto l1099
				}
				add(ruleproductExpr, position1100)
			}
			return true
		l1099:
			position, tokenIndex = position1099, tokenIndex1099
			return false
		},
		/* 81 minusExpr <- <(<((UnaryMinus spOpt)? castExpr)> Action62)> */
		func() bool {
			position1104, tokenIndex1104 := position, tokenIndex
			{
				position1105 := position
				{
					position1106 := position
					{
						position1107, tokenIndex1107 := position, tokenIndex
						if !_rules[ruleUnaryMinus]() {
							goto l1107
						}
						if !_rules[rulespOpt]() {
							goto l1107
						}
						goto l1108
					l1107:
						position, tokenIndex = position1107, tokenIndex1107
					}
				l1108:
					if !_rules[rulecastExpr]() {
						goto l1104
					}
					add(rulePegText, position1106)
				}
				if !_rules[ruleAction62]() {
					goto l1104
				}
				add(ruleminusExpr, position1105)
			}
			return true
		l1104:
			position, tokenIndex = position1104, tokenIndex1104
			return false
		},
		/* 82 castExpr <- <(<(baseExpr (spOpt (':' ':') spOpt Type)?)> Action63)> */
		func() bool {
			position1109, tokenIndex1109 := position, tokenIndex
			{
				position1110 := position
				{
					position1111 := position
					if !_rules[rulebaseExpr]() {
						goto l1109
					}
					{
						position1112, tokenIndex1112 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1112
						}
						if buffer[position] != rune(':') {
							goto l1112
						}
						position++
						if buffer[position] != rune(':') {
							goto l1112
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1112
						}
						if !_rules[ruleType]() {
							goto l1112
						}
						goto l1113
					l1112:
						position, tokenIndex = position1112, tokenIndex1112
					}
				l1113:
					add(rulePegText, position1111)
				}
				if !_rules[ruleAction63]() {
					goto l1109
				}
				add(rulecastExpr, position1110)
			}
			return true
		l1109:
			position, tokenIndex = position1109, tokenIndex1109
			return false
		},
		/* 83 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / Case / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position1114, tokenIndex1114 := position, tokenIndex
			{
				position1115 := position
				{
					position1116, tokenIndex1116 := position, tokenIndex
					if buffer[position] != rune('(') {
						goto l1117
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1117
					}
					if !_rules[ruleExpression]() {
						goto l1117
					}
					if !_rules[rulespOpt]() {
						goto l1117
					}
					if buffer[position] != rune(')') {
						goto l1117
					}
					position++
					goto l1116
				l1117:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleMapExpr]() {
						goto l1118
					}
					goto l1116
				l1118:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleBooleanLiteral]() {
						goto l1119
					}
					goto l1116
				l1119:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleNullLiteral]() {
						goto l1120
					}
					goto l1116
				l1120:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleCase]() {
						goto l1121
					}
					goto l1116
				l1121:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleRowMeta]() {
						goto l1122
					}
					goto l1116
				l1122:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleFuncTypeCast]() {
						goto l1123
					}
					goto l1116
				l1123:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleFuncApp]() {
						goto l1124
					}
					goto l1116
				l1124:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleRowValue]() {
						goto l1125
					}
					goto l1116
				l1125:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleArrayExpr]() {
						goto l1126
					}
					goto l1116
				l1126:
					position, tokenIndex = position1116, tokenIndex1116
					if !_rules[ruleLiteral]() {
						goto l1114
					}
				}
			l1116:
				add(rulebaseExpr, position1115)
			}
			return true
		l1114:
			position, tokenIndex = position1114, tokenIndex1114
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
		/* 86 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action65)> */
		func() bool {
			position1146, tokenIndex1146 := position, tokenIndex
			{
				position1147 := position
				if !_rules[ruleFunction]() {
					goto l1146
				}
				if !_rules[rulespOpt]() {
					goto l1146
				}
				if buffer[position] != rune('(') {
					goto l1146
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1146
				}
				if !_rules[ruleFuncParams]() {
					goto l1146
				}
				if !_rules[rulesp]() {
					goto l1146
				}
				if !_rules[ruleParamsOrder]() {
					goto l1146
				}
				if !_rules[rulespOpt]() {
					goto l1146
				}
				if buffer[position] != rune(')') {
					goto l1146
				}
				position++
				if !_rules[ruleAction65]() {
					goto l1146
				}
				add(ruleFuncAppWithOrderBy, position1147)
			}
			return true
		l1146:
			position, tokenIndex = position1146, tokenIndex1146
			return false
		},
		/* 87 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action66)> */
		func() bool {
			position1148, tokenIndex1148 := position, tokenIndex
			{
				position1149 := position
				if !_rules[ruleFunction]() {
					goto l1148
				}
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
				if !_rules[ruleFuncParams]() {
					goto l1148
				}
				{
					position1150 := position
					if !_rules[rulespOpt]() {
						goto l1148
					}
					add(rulePegText, position1150)
				}
				if buffer[position] != rune(')') {
					goto l1148
				}
				position++
				if !_rules[ruleAction66]() {
					goto l1148
				}
				add(ruleFuncAppWithoutOrderBy, position1149)
			}
			return true
		l1148:
			position, tokenIndex = position1148, tokenIndex1148
			return false
		},
		/* 88 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action67)> */
		func() bool {
			position1151, tokenIndex1151 := position, tokenIndex
			{
				position1152 := position
				{
					position1153 := position
					{
						position1154, tokenIndex1154 := position, tokenIndex
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1154
						}
					l1156:
						{
							position1157, tokenIndex1157 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1157
							}
							if buffer[position] != rune(',') {
								goto l1157
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1157
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1157
							}
							goto l1156
						l1157:
							position, tokenIndex = position1157, tokenIndex1157
						}
						goto l1155
					l1154:
						position, tokenIndex = position1154, tokenIndex1154
					}
				l1155:
					add(rulePegText, position1153)
				}
				if !_rules[ruleAction67]() {
					goto l1151
				}
				add(ruleFuncParams, position1152)
			}
			return true
		l1151:
			position, tokenIndex = position1151, tokenIndex1151
			return false
		},
		/* 89 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action68)> */
		func() bool {
			position1158, tokenIndex1158 := position, tokenIndex
			{
				position1159 := position
				{
					position1160 := position
					{
						position1161, tokenIndex1161 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1162
						}
						position++
						goto l1161
					l1162:
						position, tokenIndex = position1161, tokenIndex1161
						if buffer[position] != rune('O') {
							goto l1158
						}
						position++
					}
				l1161:
					{
						position1163, tokenIndex1163 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1164
						}
						position++
						goto l1163
					l1164:
						position, tokenIndex = position1163, tokenIndex1163
						if buffer[position] != rune('R') {
							goto l1158
						}
						position++
					}
				l1163:
					{
						position1165, tokenIndex1165 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1166
						}
						position++
						goto l1165
					l1166:
						position, tokenIndex = position1165, tokenIndex1165
						if buffer[position] != rune('D') {
							goto l1158
						}
						position++
					}
				l1165:
					{
						position1167, tokenIndex1167 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1168
						}
						position++
						goto l1167
					l1168:
						position, tokenIndex = position1167, tokenIndex1167
						if buffer[position] != rune('E') {
							goto l1158
						}
						position++
					}
				l1167:
					{
						position1169, tokenIndex1169 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1170
						}
						position++
						goto l1169
					l1170:
						position, tokenIndex = position1169, tokenIndex1169
						if buffer[position] != rune('R') {
							goto l1158
						}
						position++
					}
				l1169:
					if !_rules[rulesp]() {
						goto l1158
					}
					{
						position1171, tokenIndex1171 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1172
						}
						position++
						goto l1171
					l1172:
						position, tokenIndex = position1171, tokenIndex1171
						if buffer[position] != rune('B') {
							goto l1158
						}
						position++
					}
				l1171:
					{
						position1173, tokenIndex1173 := position, tokenIndex
						if buffer[position] != rune('y') {
							goto l1174
						}
						position++
						goto l1173
					l1174:
						position, tokenIndex = position1173, tokenIndex1173
						if buffer[position] != rune('Y') {
							goto l1158
						}
						position++
					}
				l1173:
					if !_rules[rulesp]() {
						goto l1158
					}
					if !_rules[ruleSortedExpression]() {
						goto l1158
					}
				l1175:
					{
						position1176, tokenIndex1176 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1176
						}
						if buffer[position] != rune(',') {
							goto l1176
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1176
						}
						if !_rules[ruleSortedExpression]() {
							goto l1176
						}
						goto l1175
					l1176:
						position, tokenIndex = position1176, tokenIndex1176
					}
					add(rulePegText, position1160)
				}
				if !_rules[ruleAction68]() {
					goto l1158
				}
				add(ruleParamsOrder, position1159)
			}
			return true
		l1158:
			position, tokenIndex = position1158, tokenIndex1158
			return false
		},
		/* 90 SortedExpression <- <(Expression OrderDirectionOpt Action69)> */
		func() bool {
			position1177, tokenIndex1177 := position, tokenIndex
			{
				position1178 := position
				if !_rules[ruleExpression]() {
					goto l1177
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1177
				}
				if !_rules[ruleAction69]() {
					goto l1177
				}
				add(ruleSortedExpression, position1178)
			}
			return true
		l1177:
			position, tokenIndex = position1177, tokenIndex1177
			return false
		},
		/* 91 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action70)> */
		func() bool {
			position1179, tokenIndex1179 := position, tokenIndex
			{
				position1180 := position
				{
					position1181 := position
					{
						position1182, tokenIndex1182 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1182
						}
						{
							position1184, tokenIndex1184 := position, tokenIndex
							if !_rules[ruleAscending]() {
								goto l1185
							}
							goto l1184
						l1185:
							position, tokenIndex = position1184, tokenIndex1184
							if !_rules[ruleDescending]() {
								goto l1182
							}
						}
					l1184:
						goto l1183
					l1182:
						position, tokenIndex = position1182, tokenIndex1182
					}
				l1183:
					add(rulePegText, position1181)
				}
				if !_rules[ruleAction70]() {
					goto l1179
				}
				add(ruleOrderDirectionOpt, position1180)
			}
			return true
		l1179:
			position, tokenIndex = position1179, tokenIndex1179
			return false
		},
		/* 92 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action71)> */
		func() bool {
			position1186, tokenIndex1186 := position, tokenIndex
			{
				position1187 := position
				{
					position1188 := position
					if buffer[position] != rune('[') {
						goto l1186
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1186
					}
					{
						position1189, tokenIndex1189 := position, tokenIndex
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1189
						}
					l1191:
						{
							position1192, tokenIndex1192 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1192
							}
							if buffer[position] != rune(',') {
								goto l1192
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1192
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1192
							}
							goto l1191
						l1192:
							position, tokenIndex = position1192, tokenIndex1192
						}
						goto l1190
					l1189:
						position, tokenIndex = position1189, tokenIndex1189
					}
				l1190:
					if !_rules[rulespOpt]() {
						goto l1186
					}
					{
						position1193, tokenIndex1193 := position, tokenIndex
						if buffer[position] != rune(',') {
							goto l1193
						}
						position++
						goto l1194
					l1193:
						position, tokenIndex = position1193, tokenIndex1193
					}
				l1194:
					if !_rules[rulespOpt]() {
						goto l1186
					}
					if buffer[position] != rune(']') {
						goto l1186
					}
					position++
					add(rulePegText, position1188)
				}
				if !_rules[ruleAction71]() {
					goto l1186
				}
				add(ruleArrayExpr, position1187)
			}
			return true
		l1186:
			position, tokenIndex = position1186, tokenIndex1186
			return false
		},
		/* 93 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action72)> */
		func() bool {
			position1195, tokenIndex1195 := position, tokenIndex
			{
				position1196 := position
				{
					position1197 := position
					if buffer[position] != rune('{') {
						goto l1195
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1195
					}
					{
						position1198, tokenIndex1198 := position, tokenIndex
						if !_rules[ruleKeyValuePair]() {
							goto l1198
						}
					l1200:
						{
							position1201, tokenIndex1201 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1201
							}
							if buffer[position] != rune(',') {
								goto l1201
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1201
							}
							if !_rules[ruleKeyValuePair]() {
								goto l1201
							}
							goto l1200
						l1201:
							position, tokenIndex = position1201, tokenIndex1201
						}
						goto l1199
					l1198:
						position, tokenIndex = position1198, tokenIndex1198
					}
				l1199:
					if !_rules[rulespOpt]() {
						goto l1195
					}
					if buffer[position] != rune('}') {
						goto l1195
					}
					position++
					add(rulePegText, position1197)
				}
				if !_rules[ruleAction72]() {
					goto l1195
				}
				add(ruleMapExpr, position1196)
			}
			return true
		l1195:
			position, tokenIndex = position1195, tokenIndex1195
			return false
		},
		/* 94 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action73)> */
		func() bool {
			position1202, tokenIndex1202 := position, tokenIndex
			{
				position1203 := position
				{
					position1204 := position
					if !_rules[ruleStringLiteral]() {
						goto l1202
					}
					if !_rules[rulespOpt]() {
						goto l1202
					}
					if buffer[position] != rune(':') {
						goto l1202
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1202
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1202
					}
					add(rulePegText, position1204)
				}
				if !_rules[ruleAction73]() {
					goto l1202
				}
				add(ruleKeyValuePair, position1203)
			}
			return true
		l1202:
			position, tokenIndex = position1202, tokenIndex1202
			return false
		},
		/* 95 Case <- <(ConditionCase / ExpressionCase)> */
		func() bool {
			position1205, tokenIndex1205 := position, tokenIndex
			{
				position1206 := position
				{
					position1207, tokenIndex1207 := position, tokenIndex
					if !_rules[ruleConditionCase]() {
						goto l1208
					}
					goto l1207
				l1208:
					position, tokenIndex = position1207, tokenIndex1207
					if !_rules[ruleExpressionCase]() {
						goto l1205
					}
				}
			l1207:
				add(ruleCase, position1206)
			}
			return true
		l1205:
			position, tokenIndex = position1205, tokenIndex1205
			return false
		},
		/* 96 ConditionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action74)> */
		func() bool {
			position1209, tokenIndex1209 := position, tokenIndex
			{
				position1210 := position
				{
					position1211, tokenIndex1211 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l1212
					}
					position++
					goto l1211
				l1212:
					position, tokenIndex = position1211, tokenIndex1211
					if buffer[position] != rune('C') {
						goto l1209
					}
					position++
				}
			l1211:
				{
					position1213, tokenIndex1213 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l1214
					}
					position++
					goto l1213
				l1214:
					position, tokenIndex = position1213, tokenIndex1213
					if buffer[position] != rune('A') {
						goto l1209
					}
					position++
				}
			l1213:
				{
					position1215, tokenIndex1215 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l1216
					}
					position++
					goto l1215
				l1216:
					position, tokenIndex = position1215, tokenIndex1215
					if buffer[position] != rune('S') {
						goto l1209
					}
					position++
				}
			l1215:
				{
					position1217, tokenIndex1217 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1218
					}
					position++
					goto l1217
				l1218:
					position, tokenIndex = position1217, tokenIndex1217
					if buffer[position] != rune('E') {
						goto l1209
					}
					position++
				}
			l1217:
				{
					position1219 := position
					if !_rules[rulesp]() {
						goto l1209
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1209
					}
				l1220:
					{
						position1221, tokenIndex1221 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1221
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1221
						}
						goto l1220
					l1221:
						position, tokenIndex = position1221, tokenIndex1221
					}
					{
						position1222, tokenIndex1222 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1222
						}
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
								goto l1222
							}
							position++
						}
					l1224:
						{
							position1226, tokenIndex1226 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l1227
							}
							position++
							goto l1226
						l1227:
							position, tokenIndex = position1226, tokenIndex1226
							if buffer[position] != rune('L') {
								goto l1222
							}
							position++
						}
					l1226:
						{
							position1228, tokenIndex1228 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l1229
							}
							position++
							goto l1228
						l1229:
							position, tokenIndex = position1228, tokenIndex1228
							if buffer[position] != rune('S') {
								goto l1222
							}
							position++
						}
					l1228:
						{
							position1230, tokenIndex1230 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1231
							}
							position++
							goto l1230
						l1231:
							position, tokenIndex = position1230, tokenIndex1230
							if buffer[position] != rune('E') {
								goto l1222
							}
							position++
						}
					l1230:
						if !_rules[rulesp]() {
							goto l1222
						}
						if !_rules[ruleExpression]() {
							goto l1222
						}
						goto l1223
					l1222:
						position, tokenIndex = position1222, tokenIndex1222
					}
				l1223:
					if !_rules[rulesp]() {
						goto l1209
					}
					{
						position1232, tokenIndex1232 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1233
						}
						position++
						goto l1232
					l1233:
						position, tokenIndex = position1232, tokenIndex1232
						if buffer[position] != rune('E') {
							goto l1209
						}
						position++
					}
				l1232:
					{
						position1234, tokenIndex1234 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1235
						}
						position++
						goto l1234
					l1235:
						position, tokenIndex = position1234, tokenIndex1234
						if buffer[position] != rune('N') {
							goto l1209
						}
						position++
					}
				l1234:
					{
						position1236, tokenIndex1236 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1237
						}
						position++
						goto l1236
					l1237:
						position, tokenIndex = position1236, tokenIndex1236
						if buffer[position] != rune('D') {
							goto l1209
						}
						position++
					}
				l1236:
					add(rulePegText, position1219)
				}
				if !_rules[ruleAction74]() {
					goto l1209
				}
				add(ruleConditionCase, position1210)
			}
			return true
		l1209:
			position, tokenIndex = position1209, tokenIndex1209
			return false
		},
		/* 97 ExpressionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') sp Expression <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action75)> */
		func() bool {
			position1238, tokenIndex1238 := position, tokenIndex
			{
				position1239 := position
				{
					position1240, tokenIndex1240 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l1241
					}
					position++
					goto l1240
				l1241:
					position, tokenIndex = position1240, tokenIndex1240
					if buffer[position] != rune('C') {
						goto l1238
					}
					position++
				}
			l1240:
				{
					position1242, tokenIndex1242 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l1243
					}
					position++
					goto l1242
				l1243:
					position, tokenIndex = position1242, tokenIndex1242
					if buffer[position] != rune('A') {
						goto l1238
					}
					position++
				}
			l1242:
				{
					position1244, tokenIndex1244 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l1245
					}
					position++
					goto l1244
				l1245:
					position, tokenIndex = position1244, tokenIndex1244
					if buffer[position] != rune('S') {
						goto l1238
					}
					position++
				}
			l1244:
				{
					position1246, tokenIndex1246 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1247
					}
					position++
					goto l1246
				l1247:
					position, tokenIndex = position1246, tokenIndex1246
					if buffer[position] != rune('E') {
						goto l1238
					}
					position++
				}
			l1246:
				if !_rules[rulesp]() {
					goto l1238
				}
				if !_rules[ruleExpression]() {
					goto l1238
				}
				{
					position1248 := position
					if !_rules[rulesp]() {
						goto l1238
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1238
					}
				l1249:
					{
						position1250, tokenIndex1250 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1250
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1250
						}
						goto l1249
					l1250:
						position, tokenIndex = position1250, tokenIndex1250
					}
					{
						position1251, tokenIndex1251 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1251
						}
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
								goto l1251
							}
							position++
						}
					l1253:
						{
							position1255, tokenIndex1255 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l1256
							}
							position++
							goto l1255
						l1256:
							position, tokenIndex = position1255, tokenIndex1255
							if buffer[position] != rune('L') {
								goto l1251
							}
							position++
						}
					l1255:
						{
							position1257, tokenIndex1257 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l1258
							}
							position++
							goto l1257
						l1258:
							position, tokenIndex = position1257, tokenIndex1257
							if buffer[position] != rune('S') {
								goto l1251
							}
							position++
						}
					l1257:
						{
							position1259, tokenIndex1259 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1260
							}
							position++
							goto l1259
						l1260:
							position, tokenIndex = position1259, tokenIndex1259
							if buffer[position] != rune('E') {
								goto l1251
							}
							position++
						}
					l1259:
						if !_rules[rulesp]() {
							goto l1251
						}
						if !_rules[ruleExpression]() {
							goto l1251
						}
						goto l1252
					l1251:
						position, tokenIndex = position1251, tokenIndex1251
					}
				l1252:
					if !_rules[rulesp]() {
						goto l1238
					}
					{
						position1261, tokenIndex1261 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1262
						}
						position++
						goto l1261
					l1262:
						position, tokenIndex = position1261, tokenIndex1261
						if buffer[position] != rune('E') {
							goto l1238
						}
						position++
					}
				l1261:
					{
						position1263, tokenIndex1263 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1264
						}
						position++
						goto l1263
					l1264:
						position, tokenIndex = position1263, tokenIndex1263
						if buffer[position] != rune('N') {
							goto l1238
						}
						position++
					}
				l1263:
					{
						position1265, tokenIndex1265 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1266
						}
						position++
						goto l1265
					l1266:
						position, tokenIndex = position1265, tokenIndex1265
						if buffer[position] != rune('D') {
							goto l1238
						}
						position++
					}
				l1265:
					add(rulePegText, position1248)
				}
				if !_rules[ruleAction75]() {
					goto l1238
				}
				add(ruleExpressionCase, position1239)
			}
			return true
		l1238:
			position, tokenIndex = position1238, tokenIndex1238
			return false
		},
		/* 98 WhenThenPair <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('n' / 'N') sp Expression sp (('t' / 'T') ('h' / 'H') ('e' / 'E') ('n' / 'N')) sp ExpressionOrWildcard Action76)> */
		func() bool {
			position1267, tokenIndex1267 := position, tokenIndex
			{
				position1268 := position
				{
					position1269, tokenIndex1269 := position, tokenIndex
					if buffer[position] != rune('w') {
						goto l1270
					}
					position++
					goto l1269
				l1270:
					position, tokenIndex = position1269, tokenIndex1269
					if buffer[position] != rune('W') {
						goto l1267
					}
					position++
				}
			l1269:
				{
					position1271, tokenIndex1271 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l1272
					}
					position++
					goto l1271
				l1272:
					position, tokenIndex = position1271, tokenIndex1271
					if buffer[position] != rune('H') {
						goto l1267
					}
					position++
				}
			l1271:
				{
					position1273, tokenIndex1273 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1274
					}
					position++
					goto l1273
				l1274:
					position, tokenIndex = position1273, tokenIndex1273
					if buffer[position] != rune('E') {
						goto l1267
					}
					position++
				}
			l1273:
				{
					position1275, tokenIndex1275 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l1276
					}
					position++
					goto l1275
				l1276:
					position, tokenIndex = position1275, tokenIndex1275
					if buffer[position] != rune('N') {
						goto l1267
					}
					position++
				}
			l1275:
				if !_rules[rulesp]() {
					goto l1267
				}
				if !_rules[ruleExpression]() {
					goto l1267
				}
				if !_rules[rulesp]() {
					goto l1267
				}
				{
					position1277, tokenIndex1277 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l1278
					}
					position++
					goto l1277
				l1278:
					position, tokenIndex = position1277, tokenIndex1277
					if buffer[position] != rune('T') {
						goto l1267
					}
					position++
				}
			l1277:
				{
					position1279, tokenIndex1279 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l1280
					}
					position++
					goto l1279
				l1280:
					position, tokenIndex = position1279, tokenIndex1279
					if buffer[position] != rune('H') {
						goto l1267
					}
					position++
				}
			l1279:
				{
					position1281, tokenIndex1281 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1282
					}
					position++
					goto l1281
				l1282:
					position, tokenIndex = position1281, tokenIndex1281
					if buffer[position] != rune('E') {
						goto l1267
					}
					position++
				}
			l1281:
				{
					position1283, tokenIndex1283 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l1284
					}
					position++
					goto l1283
				l1284:
					position, tokenIndex = position1283, tokenIndex1283
					if buffer[position] != rune('N') {
						goto l1267
					}
					position++
				}
			l1283:
				if !_rules[rulesp]() {
					goto l1267
				}
				if !_rules[ruleExpressionOrWildcard]() {
					goto l1267
				}
				if !_rules[ruleAction76]() {
					goto l1267
				}
				add(ruleWhenThenPair, position1268)
			}
			return true
		l1267:
			position, tokenIndex = position1267, tokenIndex1267
			return false
		},
		/* 99 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1285, tokenIndex1285 := position, tokenIndex
			{
				position1286 := position
				{
					position1287, tokenIndex1287 := position, tokenIndex
					if !_rules[ruleFloatLiteral]() {
						goto l1288
					}
					goto l1287
				l1288:
					position, tokenIndex = position1287, tokenIndex1287
					if !_rules[ruleNumericLiteral]() {
						goto l1289
					}
					goto l1287
				l1289:
					position, tokenIndex = position1287, tokenIndex1287
					if !_rules[ruleStringLiteral]() {
						goto l1285
					}
				}
			l1287:
				add(ruleLiteral, position1286)
			}
			return true
		l1285:
			position, tokenIndex = position1285, tokenIndex1285
			return false
		},
		/* 100 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1290, tokenIndex1290 := position, tokenIndex
			{
				position1291 := position
				{
					position1292, tokenIndex1292 := position, tokenIndex
					if !_rules[ruleEqual]() {
						goto l1293
					}
					goto l1292
				l1293:
					position, tokenIndex = position1292, tokenIndex1292
					if !_rules[ruleNotEqual]() {
						goto l1294
					}
					goto l1292
				l1294:
					position, tokenIndex = position1292, tokenIndex1292
					if !_rules[ruleLessOrEqual]() {
						goto l1295
					}
					goto l1292
				l1295:
					position, tokenIndex = position1292, tokenIndex1292
					if !_rules[ruleLess]() {
						goto l1296
					}
					goto l1292
				l1296:
					position, tokenIndex = position1292, tokenIndex1292
					if !_rules[ruleGreaterOrEqual]() {
						goto l1297
					}
					goto l1292
				l1297:
					position, tokenIndex = position1292, tokenIndex1292
					if !_rules[ruleGreater]() {
						goto l1298
					}
					goto l1292
				l1298:
					position, tokenIndex = position1292, tokenIndex1292
					if !_rules[ruleNotEqual]() {
						goto l1290
					}
				}
			l1292:
				add(ruleComparisonOp, position1291)
			}
			return true
		l1290:
			position, tokenIndex = position1290, tokenIndex1290
			return false
		},
		/* 101 OtherOp <- <Concat> */
		func() bool {
			position1299, tokenIndex1299 := position, tokenIndex
			{
				position1300 := position
				if !_rules[ruleConcat]() {
					goto l1299
				}
				add(ruleOtherOp, position1300)
			}
			return true
		l1299:
			position, tokenIndex = position1299, tokenIndex1299
			return false
		},
		/* 102 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1301, tokenIndex1301 := position, tokenIndex
			{
				position1302 := position
				{
					position1303, tokenIndex1303 := position, tokenIndex
					if !_rules[ruleIsNot]() {
						goto l1304
					}
					goto l1303
				l1304:
					position, tokenIndex = position1303, tokenIndex1303
					if !_rules[ruleIs]() {
						goto l1301
					}
				}
			l1303:
				add(ruleIsOp, position1302)
			}
			return true
		l1301:
			position, tokenIndex = position1301, tokenIndex1301
			return false
		},
		/* 103 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1305, tokenIndex1305 := position, tokenIndex
			{
				position1306 := position
				{
					position1307, tokenIndex1307 := position, tokenIndex
					if !_rules[rulePlus]() {
						goto l1308
					}
					goto l1307
				l1308:
					position, tokenIndex = position1307, tokenIndex1307
					if !_rules[ruleMinus]() {
						goto l1305
					}
				}
			l1307:
				add(rulePlusMinusOp, position1306)
			}
			return true
		l1305:
			position, tokenIndex = position1305, tokenIndex1305
			return false
		},
		/* 104 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1309, tokenIndex1309 := position, tokenIndex
			{
				position1310 := position
				{
					position1311, tokenIndex1311 := position, tokenIndex
					if !_rules[ruleMultiply]() {
						goto l1312
					}
					goto l1311
				l1312:
					position, tokenIndex = position1311, tokenIndex1311
					if !_rules[ruleDivide]() {
						goto l1313
					}
					goto l1311
				l1313:
					position, tokenIndex = position1311, tokenIndex1311
					if !_rules[ruleModulo]() {
						goto l1309
					}
				}
			l1311:
				add(ruleMultDivOp, position1310)
			}
			return true
		l1309:
			position, tokenIndex = position1309, tokenIndex1309
			return false
		},
		/* 105 Stream <- <(<ident> Action77)> */
		func() bool {
			position1314, tokenIndex1314 := position, tokenIndex
			{
				position1315 := position
				{
					position1316 := position
					if !_rules[ruleident]() {
						goto l1314
					}
					add(rulePegText, position1316)
				}
				if !_rules[ruleAction77]() {
					goto l1314
				}
				add(ruleStream, position1315)
			}
			return true
		l1314:
			position, tokenIndex = position1314, tokenIndex1314
			return false
		},
		/* 106 RowMeta <- <RowTimestamp> */
		func() bool {
			position1317, tokenIndex1317 := position, tokenIndex
			{
				position1318 := position
				if !_rules[ruleRowTimestamp]() {
					goto l1317
				}
				add(ruleRowMeta, position1318)
			}
			return true
		l1317:
			position, tokenIndex = position1317, tokenIndex1317
			return false
		},
		/* 107 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action78)> */
		func() bool {
			position1319, tokenIndex1319 := position, tokenIndex
			{
				position1320 := position
				{
					position1321 := position
					{
						position1322, tokenIndex1322 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1322
						}
						if buffer[position] != rune(':') {
							goto l1322
						}
						position++
						goto l1323
					l1322:
						position, tokenIndex = position1322, tokenIndex1322
					}
				l1323:
					if buffer[position] != rune('t') {
						goto l1319
					}
					position++
					if buffer[position] != rune('s') {
						goto l1319
					}
					position++
					if buffer[position] != rune('(') {
						goto l1319
					}
					position++
					if buffer[position] != rune(')') {
						goto l1319
					}
					position++
					add(rulePegText, position1321)
				}
				if !_rules[ruleAction78]() {
					goto l1319
				}
				add(ruleRowTimestamp, position1320)
			}
			return true
		l1319:
			position, tokenIndex = position1319, tokenIndex1319
			return false
		},
		/* 108 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action79)> */
		func() bool {
			position1324, tokenIndex1324 := position, tokenIndex
			{
				position1325 := position
				{
					position1326 := position
					{
						position1327, tokenIndex1327 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1327
						}
						if buffer[position] != rune(':') {
							goto l1327
						}
						position++
						{
							position1329, tokenIndex1329 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l1329
							}
							position++
							goto l1327
						l1329:
							position, tokenIndex = position1329, tokenIndex1329
						}
						goto l1328
					l1327:
						position, tokenIndex = position1327, tokenIndex1327
					}
				l1328:
					if !_rules[rulejsonGetPath]() {
						goto l1324
					}
					add(rulePegText, position1326)
				}
				if !_rules[ruleAction79]() {
					goto l1324
				}
				add(ruleRowValue, position1325)
			}
			return true
		l1324:
			position, tokenIndex = position1324, tokenIndex1324
			return false
		},
		/* 109 NumericLiteral <- <(<('-'? [0-9]+)> Action80)> */
		func() bool {
			position1330, tokenIndex1330 := position, tokenIndex
			{
				position1331 := position
				{
					position1332 := position
					{
						position1333, tokenIndex1333 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1333
						}
						position++
						goto l1334
					l1333:
						position, tokenIndex = position1333, tokenIndex1333
					}
				l1334:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1330
					}
					position++
				l1335:
					{
						position1336, tokenIndex1336 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1336
						}
						position++
						goto l1335
					l1336:
						position, tokenIndex = position1336, tokenIndex1336
					}
					add(rulePegText, position1332)
				}
				if !_rules[ruleAction80]() {
					goto l1330
				}
				add(ruleNumericLiteral, position1331)
			}
			return true
		l1330:
			position, tokenIndex = position1330, tokenIndex1330
			return false
		},
		/* 110 NonNegativeNumericLiteral <- <(<[0-9]+> Action81)> */
		func() bool {
			position1337, tokenIndex1337 := position, tokenIndex
			{
				position1338 := position
				{
					position1339 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1337
					}
					position++
				l1340:
					{
						position1341, tokenIndex1341 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1341
						}
						position++
						goto l1340
					l1341:
						position, tokenIndex = position1341, tokenIndex1341
					}
					add(rulePegText, position1339)
				}
				if !_rules[ruleAction81]() {
					goto l1337
				}
				add(ruleNonNegativeNumericLiteral, position1338)
			}
			return true
		l1337:
			position, tokenIndex = position1337, tokenIndex1337
			return false
		},
		/* 111 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action82)> */
		func() bool {
			position1342, tokenIndex1342 := position, tokenIndex
			{
				position1343 := position
				{
					position1344 := position
					{
						position1345, tokenIndex1345 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1345
						}
						position++
						goto l1346
					l1345:
						position, tokenIndex = position1345, tokenIndex1345
					}
				l1346:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1342
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
					if buffer[position] != rune('.') {
						goto l1342
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1342
					}
					position++
				l1349:
					{
						position1350, tokenIndex1350 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1350
						}
						position++
						goto l1349
					l1350:
						position, tokenIndex = position1350, tokenIndex1350
					}
					add(rulePegText, position1344)
				}
				if !_rules[ruleAction82]() {
					goto l1342
				}
				add(ruleFloatLiteral, position1343)
			}
			return true
		l1342:
			position, tokenIndex = position1342, tokenIndex1342
			return false
		},
		/* 112 Function <- <(<ident> Action83)> */
		func() bool {
			position1351, tokenIndex1351 := position, tokenIndex
			{
				position1352 := position
				{
					position1353 := position
					if !_rules[ruleident]() {
						goto l1351
					}
					add(rulePegText, position1353)
				}
				if !_rules[ruleAction83]() {
					goto l1351
				}
				add(ruleFunction, position1352)
			}
			return true
		l1351:
			position, tokenIndex = position1351, tokenIndex1351
			return false
		},
		/* 113 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action84)> */
		func() bool {
			position1354, tokenIndex1354 := position, tokenIndex
			{
				position1355 := position
				{
					position1356 := position
					{
						position1357, tokenIndex1357 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1358
						}
						position++
						goto l1357
					l1358:
						position, tokenIndex = position1357, tokenIndex1357
						if buffer[position] != rune('N') {
							goto l1354
						}
						position++
					}
				l1357:
					{
						position1359, tokenIndex1359 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1360
						}
						position++
						goto l1359
					l1360:
						position, tokenIndex = position1359, tokenIndex1359
						if buffer[position] != rune('U') {
							goto l1354
						}
						position++
					}
				l1359:
					{
						position1361, tokenIndex1361 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1362
						}
						position++
						goto l1361
					l1362:
						position, tokenIndex = position1361, tokenIndex1361
						if buffer[position] != rune('L') {
							goto l1354
						}
						position++
					}
				l1361:
					{
						position1363, tokenIndex1363 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1364
						}
						position++
						goto l1363
					l1364:
						position, tokenIndex = position1363, tokenIndex1363
						if buffer[position] != rune('L') {
							goto l1354
						}
						position++
					}
				l1363:
					add(rulePegText, position1356)
				}
				if !_rules[ruleAction84]() {
					goto l1354
				}
				add(ruleNullLiteral, position1355)
			}
			return true
		l1354:
			position, tokenIndex = position1354, tokenIndex1354
			return false
		},
		/* 114 Missing <- <(<(('m' / 'M') ('i' / 'I') ('s' / 'S') ('s' / 'S') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action85)> */
		func() bool {
			position1365, tokenIndex1365 := position, tokenIndex
			{
				position1366 := position
				{
					position1367 := position
					{
						position1368, tokenIndex1368 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1369
						}
						position++
						goto l1368
					l1369:
						position, tokenIndex = position1368, tokenIndex1368
						if buffer[position] != rune('M') {
							goto l1365
						}
						position++
					}
				l1368:
					{
						position1370, tokenIndex1370 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1371
						}
						position++
						goto l1370
					l1371:
						position, tokenIndex = position1370, tokenIndex1370
						if buffer[position] != rune('I') {
							goto l1365
						}
						position++
					}
				l1370:
					{
						position1372, tokenIndex1372 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1373
						}
						position++
						goto l1372
					l1373:
						position, tokenIndex = position1372, tokenIndex1372
						if buffer[position] != rune('S') {
							goto l1365
						}
						position++
					}
				l1372:
					{
						position1374, tokenIndex1374 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1375
						}
						position++
						goto l1374
					l1375:
						position, tokenIndex = position1374, tokenIndex1374
						if buffer[position] != rune('S') {
							goto l1365
						}
						position++
					}
				l1374:
					{
						position1376, tokenIndex1376 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1377
						}
						position++
						goto l1376
					l1377:
						position, tokenIndex = position1376, tokenIndex1376
						if buffer[position] != rune('I') {
							goto l1365
						}
						position++
					}
				l1376:
					{
						position1378, tokenIndex1378 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1379
						}
						position++
						goto l1378
					l1379:
						position, tokenIndex = position1378, tokenIndex1378
						if buffer[position] != rune('N') {
							goto l1365
						}
						position++
					}
				l1378:
					{
						position1380, tokenIndex1380 := position, tokenIndex
						if buffer[position] != rune('g') {
							goto l1381
						}
						position++
						goto l1380
					l1381:
						position, tokenIndex = position1380, tokenIndex1380
						if buffer[position] != rune('G') {
							goto l1365
						}
						position++
					}
				l1380:
					add(rulePegText, position1367)
				}
				if !_rules[ruleAction85]() {
					goto l1365
				}
				add(ruleMissing, position1366)
			}
			return true
		l1365:
			position, tokenIndex = position1365, tokenIndex1365
			return false
		},
		/* 115 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1382, tokenIndex1382 := position, tokenIndex
			{
				position1383 := position
				{
					position1384, tokenIndex1384 := position, tokenIndex
					if !_rules[ruleTRUE]() {
						goto l1385
					}
					goto l1384
				l1385:
					position, tokenIndex = position1384, tokenIndex1384
					if !_rules[ruleFALSE]() {
						goto l1382
					}
				}
			l1384:
				add(ruleBooleanLiteral, position1383)
			}
			return true
		l1382:
			position, tokenIndex = position1382, tokenIndex1382
			return false
		},
		/* 116 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action86)> */
		func() bool {
			position1386, tokenIndex1386 := position, tokenIndex
			{
				position1387 := position
				{
					position1388 := position
					{
						position1389, tokenIndex1389 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1390
						}
						position++
						goto l1389
					l1390:
						position, tokenIndex = position1389, tokenIndex1389
						if buffer[position] != rune('T') {
							goto l1386
						}
						position++
					}
				l1389:
					{
						position1391, tokenIndex1391 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1392
						}
						position++
						goto l1391
					l1392:
						position, tokenIndex = position1391, tokenIndex1391
						if buffer[position] != rune('R') {
							goto l1386
						}
						position++
					}
				l1391:
					{
						position1393, tokenIndex1393 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1394
						}
						position++
						goto l1393
					l1394:
						position, tokenIndex = position1393, tokenIndex1393
						if buffer[position] != rune('U') {
							goto l1386
						}
						position++
					}
				l1393:
					{
						position1395, tokenIndex1395 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1396
						}
						position++
						goto l1395
					l1396:
						position, tokenIndex = position1395, tokenIndex1395
						if buffer[position] != rune('E') {
							goto l1386
						}
						position++
					}
				l1395:
					add(rulePegText, position1388)
				}
				if !_rules[ruleAction86]() {
					goto l1386
				}
				add(ruleTRUE, position1387)
			}
			return true
		l1386:
			position, tokenIndex = position1386, tokenIndex1386
			return false
		},
		/* 117 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action87)> */
		func() bool {
			position1397, tokenIndex1397 := position, tokenIndex
			{
				position1398 := position
				{
					position1399 := position
					{
						position1400, tokenIndex1400 := position, tokenIndex
						if buffer[position] != rune('f') {
							goto l1401
						}
						position++
						goto l1400
					l1401:
						position, tokenIndex = position1400, tokenIndex1400
						if buffer[position] != rune('F') {
							goto l1397
						}
						position++
					}
				l1400:
					{
						position1402, tokenIndex1402 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1403
						}
						position++
						goto l1402
					l1403:
						position, tokenIndex = position1402, tokenIndex1402
						if buffer[position] != rune('A') {
							goto l1397
						}
						position++
					}
				l1402:
					{
						position1404, tokenIndex1404 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1405
						}
						position++
						goto l1404
					l1405:
						position, tokenIndex = position1404, tokenIndex1404
						if buffer[position] != rune('L') {
							goto l1397
						}
						position++
					}
				l1404:
					{
						position1406, tokenIndex1406 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1407
						}
						position++
						goto l1406
					l1407:
						position, tokenIndex = position1406, tokenIndex1406
						if buffer[position] != rune('S') {
							goto l1397
						}
						position++
					}
				l1406:
					{
						position1408, tokenIndex1408 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1409
						}
						position++
						goto l1408
					l1409:
						position, tokenIndex = position1408, tokenIndex1408
						if buffer[position] != rune('E') {
							goto l1397
						}
						position++
					}
				l1408:
					add(rulePegText, position1399)
				}
				if !_rules[ruleAction87]() {
					goto l1397
				}
				add(ruleFALSE, position1398)
			}
			return true
		l1397:
			position, tokenIndex = position1397, tokenIndex1397
			return false
		},
		/* 118 Wildcard <- <(<((ident ':' !':')? '*')> Action88)> */
		func() bool {
			position1410, tokenIndex1410 := position, tokenIndex
			{
				position1411 := position
				{
					position1412 := position
					{
						position1413, tokenIndex1413 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1413
						}
						if buffer[position] != rune(':') {
							goto l1413
						}
						position++
						{
							position1415, tokenIndex1415 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l1415
							}
							position++
							goto l1413
						l1415:
							position, tokenIndex = position1415, tokenIndex1415
						}
						goto l1414
					l1413:
						position, tokenIndex = position1413, tokenIndex1413
					}
				l1414:
					if buffer[position] != rune('*') {
						goto l1410
					}
					position++
					add(rulePegText, position1412)
				}
				if !_rules[ruleAction88]() {
					goto l1410
				}
				add(ruleWildcard, position1411)
			}
			return true
		l1410:
			position, tokenIndex = position1410, tokenIndex1410
			return false
		},
		/* 119 StringLiteral <- <(<('"' (('"' '"') / (!'"' .))* '"')> Action89)> */
		func() bool {
			position1416, tokenIndex1416 := position, tokenIndex
			{
				position1417 := position
				{
					position1418 := position
					if buffer[position] != rune('"') {
						goto l1416
					}
					position++
				l1419:
					{
						position1420, tokenIndex1420 := position, tokenIndex
						{
							position1421, tokenIndex1421 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l1422
							}
							position++
							if buffer[position] != rune('"') {
								goto l1422
							}
							position++
							goto l1421
						l1422:
							position, tokenIndex = position1421, tokenIndex1421
							{
								position1423, tokenIndex1423 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l1423
								}
								position++
								goto l1420
							l1423:
								position, tokenIndex = position1423, tokenIndex1423
							}
							if !matchDot() {
								goto l1420
							}
						}
					l1421:
						goto l1419
					l1420:
						position, tokenIndex = position1420, tokenIndex1420
					}
					if buffer[position] != rune('"') {
						goto l1416
					}
					position++
					add(rulePegText, position1418)
				}
				if !_rules[ruleAction89]() {
					goto l1416
				}
				add(ruleStringLiteral, position1417)
			}
			return true
		l1416:
			position, tokenIndex = position1416, tokenIndex1416
			return false
		},
		/* 120 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action90)> */
		func() bool {
			position1424, tokenIndex1424 := position, tokenIndex
			{
				position1425 := position
				{
					position1426 := position
					{
						position1427, tokenIndex1427 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1428
						}
						position++
						goto l1427
					l1428:
						position, tokenIndex = position1427, tokenIndex1427
						if buffer[position] != rune('I') {
							goto l1424
						}
						position++
					}
				l1427:
					{
						position1429, tokenIndex1429 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1430
						}
						position++
						goto l1429
					l1430:
						position, tokenIndex = position1429, tokenIndex1429
						if buffer[position] != rune('S') {
							goto l1424
						}
						position++
					}
				l1429:
					{
						position1431, tokenIndex1431 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1432
						}
						position++
						goto l1431
					l1432:
						position, tokenIndex = position1431, tokenIndex1431
						if buffer[position] != rune('T') {
							goto l1424
						}
						position++
					}
				l1431:
					{
						position1433, tokenIndex1433 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1434
						}
						position++
						goto l1433
					l1434:
						position, tokenIndex = position1433, tokenIndex1433
						if buffer[position] != rune('R') {
							goto l1424
						}
						position++
					}
				l1433:
					{
						position1435, tokenIndex1435 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1436
						}
						position++
						goto l1435
					l1436:
						position, tokenIndex = position1435, tokenIndex1435
						if buffer[position] != rune('E') {
							goto l1424
						}
						position++
					}
				l1435:
					{
						position1437, tokenIndex1437 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1438
						}
						position++
						goto l1437
					l1438:
						position, tokenIndex = position1437, tokenIndex1437
						if buffer[position] != rune('A') {
							goto l1424
						}
						position++
					}
				l1437:
					{
						position1439, tokenIndex1439 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1440
						}
						position++
						goto l1439
					l1440:
						position, tokenIndex = position1439, tokenIndex1439
						if buffer[position] != rune('M') {
							goto l1424
						}
						position++
					}
				l1439:
					add(rulePegText, position1426)
				}
				if !_rules[ruleAction90]() {
					goto l1424
				}
				add(ruleISTREAM, position1425)
			}
			return true
		l1424:
			position, tokenIndex = position1424, tokenIndex1424
			return false
		},
		/* 121 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action91)> */
		func() bool {
			position1441, tokenIndex1441 := position, tokenIndex
			{
				position1442 := position
				{
					position1443 := position
					{
						position1444, tokenIndex1444 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1445
						}
						position++
						goto l1444
					l1445:
						position, tokenIndex = position1444, tokenIndex1444
						if buffer[position] != rune('D') {
							goto l1441
						}
						position++
					}
				l1444:
					{
						position1446, tokenIndex1446 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1447
						}
						position++
						goto l1446
					l1447:
						position, tokenIndex = position1446, tokenIndex1446
						if buffer[position] != rune('S') {
							goto l1441
						}
						position++
					}
				l1446:
					{
						position1448, tokenIndex1448 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1449
						}
						position++
						goto l1448
					l1449:
						position, tokenIndex = position1448, tokenIndex1448
						if buffer[position] != rune('T') {
							goto l1441
						}
						position++
					}
				l1448:
					{
						position1450, tokenIndex1450 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1451
						}
						position++
						goto l1450
					l1451:
						position, tokenIndex = position1450, tokenIndex1450
						if buffer[position] != rune('R') {
							goto l1441
						}
						position++
					}
				l1450:
					{
						position1452, tokenIndex1452 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1453
						}
						position++
						goto l1452
					l1453:
						position, tokenIndex = position1452, tokenIndex1452
						if buffer[position] != rune('E') {
							goto l1441
						}
						position++
					}
				l1452:
					{
						position1454, tokenIndex1454 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1455
						}
						position++
						goto l1454
					l1455:
						position, tokenIndex = position1454, tokenIndex1454
						if buffer[position] != rune('A') {
							goto l1441
						}
						position++
					}
				l1454:
					{
						position1456, tokenIndex1456 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1457
						}
						position++
						goto l1456
					l1457:
						position, tokenIndex = position1456, tokenIndex1456
						if buffer[position] != rune('M') {
							goto l1441
						}
						position++
					}
				l1456:
					add(rulePegText, position1443)
				}
				if !_rules[ruleAction91]() {
					goto l1441
				}
				add(ruleDSTREAM, position1442)
			}
			return true
		l1441:
			position, tokenIndex = position1441, tokenIndex1441
			return false
		},
		/* 122 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action92)> */
		func() bool {
			position1458, tokenIndex1458 := position, tokenIndex
			{
				position1459 := position
				{
					position1460 := position
					{
						position1461, tokenIndex1461 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1462
						}
						position++
						goto l1461
					l1462:
						position, tokenIndex = position1461, tokenIndex1461
						if buffer[position] != rune('R') {
							goto l1458
						}
						position++
					}
				l1461:
					{
						position1463, tokenIndex1463 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1464
						}
						position++
						goto l1463
					l1464:
						position, tokenIndex = position1463, tokenIndex1463
						if buffer[position] != rune('S') {
							goto l1458
						}
						position++
					}
				l1463:
					{
						position1465, tokenIndex1465 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1466
						}
						position++
						goto l1465
					l1466:
						position, tokenIndex = position1465, tokenIndex1465
						if buffer[position] != rune('T') {
							goto l1458
						}
						position++
					}
				l1465:
					{
						position1467, tokenIndex1467 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1468
						}
						position++
						goto l1467
					l1468:
						position, tokenIndex = position1467, tokenIndex1467
						if buffer[position] != rune('R') {
							goto l1458
						}
						position++
					}
				l1467:
					{
						position1469, tokenIndex1469 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1470
						}
						position++
						goto l1469
					l1470:
						position, tokenIndex = position1469, tokenIndex1469
						if buffer[position] != rune('E') {
							goto l1458
						}
						position++
					}
				l1469:
					{
						position1471, tokenIndex1471 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1472
						}
						position++
						goto l1471
					l1472:
						position, tokenIndex = position1471, tokenIndex1471
						if buffer[position] != rune('A') {
							goto l1458
						}
						position++
					}
				l1471:
					{
						position1473, tokenIndex1473 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1474
						}
						position++
						goto l1473
					l1474:
						position, tokenIndex = position1473, tokenIndex1473
						if buffer[position] != rune('M') {
							goto l1458
						}
						position++
					}
				l1473:
					add(rulePegText, position1460)
				}
				if !_rules[ruleAction92]() {
					goto l1458
				}
				add(ruleRSTREAM, position1459)
			}
			return true
		l1458:
			position, tokenIndex = position1458, tokenIndex1458
			return false
		},
		/* 123 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action93)> */
		func() bool {
			position1475, tokenIndex1475 := position, tokenIndex
			{
				position1476 := position
				{
					position1477 := position
					{
						position1478, tokenIndex1478 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1479
						}
						position++
						goto l1478
					l1479:
						position, tokenIndex = position1478, tokenIndex1478
						if buffer[position] != rune('T') {
							goto l1475
						}
						position++
					}
				l1478:
					{
						position1480, tokenIndex1480 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1481
						}
						position++
						goto l1480
					l1481:
						position, tokenIndex = position1480, tokenIndex1480
						if buffer[position] != rune('U') {
							goto l1475
						}
						position++
					}
				l1480:
					{
						position1482, tokenIndex1482 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1483
						}
						position++
						goto l1482
					l1483:
						position, tokenIndex = position1482, tokenIndex1482
						if buffer[position] != rune('P') {
							goto l1475
						}
						position++
					}
				l1482:
					{
						position1484, tokenIndex1484 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1485
						}
						position++
						goto l1484
					l1485:
						position, tokenIndex = position1484, tokenIndex1484
						if buffer[position] != rune('L') {
							goto l1475
						}
						position++
					}
				l1484:
					{
						position1486, tokenIndex1486 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1487
						}
						position++
						goto l1486
					l1487:
						position, tokenIndex = position1486, tokenIndex1486
						if buffer[position] != rune('E') {
							goto l1475
						}
						position++
					}
				l1486:
					{
						position1488, tokenIndex1488 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1489
						}
						position++
						goto l1488
					l1489:
						position, tokenIndex = position1488, tokenIndex1488
						if buffer[position] != rune('S') {
							goto l1475
						}
						position++
					}
				l1488:
					add(rulePegText, position1477)
				}
				if !_rules[ruleAction93]() {
					goto l1475
				}
				add(ruleTUPLES, position1476)
			}
			return true
		l1475:
			position, tokenIndex = position1475, tokenIndex1475
			return false
		},
		/* 124 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action94)> */
		func() bool {
			position1490, tokenIndex1490 := position, tokenIndex
			{
				position1491 := position
				{
					position1492 := position
					{
						position1493, tokenIndex1493 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1494
						}
						position++
						goto l1493
					l1494:
						position, tokenIndex = position1493, tokenIndex1493
						if buffer[position] != rune('S') {
							goto l1490
						}
						position++
					}
				l1493:
					{
						position1495, tokenIndex1495 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1496
						}
						position++
						goto l1495
					l1496:
						position, tokenIndex = position1495, tokenIndex1495
						if buffer[position] != rune('E') {
							goto l1490
						}
						position++
					}
				l1495:
					{
						position1497, tokenIndex1497 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1498
						}
						position++
						goto l1497
					l1498:
						position, tokenIndex = position1497, tokenIndex1497
						if buffer[position] != rune('C') {
							goto l1490
						}
						position++
					}
				l1497:
					{
						position1499, tokenIndex1499 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1500
						}
						position++
						goto l1499
					l1500:
						position, tokenIndex = position1499, tokenIndex1499
						if buffer[position] != rune('O') {
							goto l1490
						}
						position++
					}
				l1499:
					{
						position1501, tokenIndex1501 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1502
						}
						position++
						goto l1501
					l1502:
						position, tokenIndex = position1501, tokenIndex1501
						if buffer[position] != rune('N') {
							goto l1490
						}
						position++
					}
				l1501:
					{
						position1503, tokenIndex1503 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1504
						}
						position++
						goto l1503
					l1504:
						position, tokenIndex = position1503, tokenIndex1503
						if buffer[position] != rune('D') {
							goto l1490
						}
						position++
					}
				l1503:
					{
						position1505, tokenIndex1505 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1506
						}
						position++
						goto l1505
					l1506:
						position, tokenIndex = position1505, tokenIndex1505
						if buffer[position] != rune('S') {
							goto l1490
						}
						position++
					}
				l1505:
					add(rulePegText, position1492)
				}
				if !_rules[ruleAction94]() {
					goto l1490
				}
				add(ruleSECONDS, position1491)
			}
			return true
		l1490:
			position, tokenIndex = position1490, tokenIndex1490
			return false
		},
		/* 125 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action95)> */
		func() bool {
			position1507, tokenIndex1507 := position, tokenIndex
			{
				position1508 := position
				{
					position1509 := position
					{
						position1510, tokenIndex1510 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1511
						}
						position++
						goto l1510
					l1511:
						position, tokenIndex = position1510, tokenIndex1510
						if buffer[position] != rune('M') {
							goto l1507
						}
						position++
					}
				l1510:
					{
						position1512, tokenIndex1512 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1513
						}
						position++
						goto l1512
					l1513:
						position, tokenIndex = position1512, tokenIndex1512
						if buffer[position] != rune('I') {
							goto l1507
						}
						position++
					}
				l1512:
					{
						position1514, tokenIndex1514 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1515
						}
						position++
						goto l1514
					l1515:
						position, tokenIndex = position1514, tokenIndex1514
						if buffer[position] != rune('L') {
							goto l1507
						}
						position++
					}
				l1514:
					{
						position1516, tokenIndex1516 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1517
						}
						position++
						goto l1516
					l1517:
						position, tokenIndex = position1516, tokenIndex1516
						if buffer[position] != rune('L') {
							goto l1507
						}
						position++
					}
				l1516:
					{
						position1518, tokenIndex1518 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1519
						}
						position++
						goto l1518
					l1519:
						position, tokenIndex = position1518, tokenIndex1518
						if buffer[position] != rune('I') {
							goto l1507
						}
						position++
					}
				l1518:
					{
						position1520, tokenIndex1520 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1521
						}
						position++
						goto l1520
					l1521:
						position, tokenIndex = position1520, tokenIndex1520
						if buffer[position] != rune('S') {
							goto l1507
						}
						position++
					}
				l1520:
					{
						position1522, tokenIndex1522 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1523
						}
						position++
						goto l1522
					l1523:
						position, tokenIndex = position1522, tokenIndex1522
						if buffer[position] != rune('E') {
							goto l1507
						}
						position++
					}
				l1522:
					{
						position1524, tokenIndex1524 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1525
						}
						position++
						goto l1524
					l1525:
						position, tokenIndex = position1524, tokenIndex1524
						if buffer[position] != rune('C') {
							goto l1507
						}
						position++
					}
				l1524:
					{
						position1526, tokenIndex1526 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1527
						}
						position++
						goto l1526
					l1527:
						position, tokenIndex = position1526, tokenIndex1526
						if buffer[position] != rune('O') {
							goto l1507
						}
						position++
					}
				l1526:
					{
						position1528, tokenIndex1528 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1529
						}
						position++
						goto l1528
					l1529:
						position, tokenIndex = position1528, tokenIndex1528
						if buffer[position] != rune('N') {
							goto l1507
						}
						position++
					}
				l1528:
					{
						position1530, tokenIndex1530 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1531
						}
						position++
						goto l1530
					l1531:
						position, tokenIndex = position1530, tokenIndex1530
						if buffer[position] != rune('D') {
							goto l1507
						}
						position++
					}
				l1530:
					{
						position1532, tokenIndex1532 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1533
						}
						position++
						goto l1532
					l1533:
						position, tokenIndex = position1532, tokenIndex1532
						if buffer[position] != rune('S') {
							goto l1507
						}
						position++
					}
				l1532:
					add(rulePegText, position1509)
				}
				if !_rules[ruleAction95]() {
					goto l1507
				}
				add(ruleMILLISECONDS, position1508)
			}
			return true
		l1507:
			position, tokenIndex = position1507, tokenIndex1507
			return false
		},
		/* 126 Wait <- <(<(('w' / 'W') ('a' / 'A') ('i' / 'I') ('t' / 'T'))> Action96)> */
		func() bool {
			position1534, tokenIndex1534 := position, tokenIndex
			{
				position1535 := position
				{
					position1536 := position
					{
						position1537, tokenIndex1537 := position, tokenIndex
						if buffer[position] != rune('w') {
							goto l1538
						}
						position++
						goto l1537
					l1538:
						position, tokenIndex = position1537, tokenIndex1537
						if buffer[position] != rune('W') {
							goto l1534
						}
						position++
					}
				l1537:
					{
						position1539, tokenIndex1539 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1540
						}
						position++
						goto l1539
					l1540:
						position, tokenIndex = position1539, tokenIndex1539
						if buffer[position] != rune('A') {
							goto l1534
						}
						position++
					}
				l1539:
					{
						position1541, tokenIndex1541 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1542
						}
						position++
						goto l1541
					l1542:
						position, tokenIndex = position1541, tokenIndex1541
						if buffer[position] != rune('I') {
							goto l1534
						}
						position++
					}
				l1541:
					{
						position1543, tokenIndex1543 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1544
						}
						position++
						goto l1543
					l1544:
						position, tokenIndex = position1543, tokenIndex1543
						if buffer[position] != rune('T') {
							goto l1534
						}
						position++
					}
				l1543:
					add(rulePegText, position1536)
				}
				if !_rules[ruleAction96]() {
					goto l1534
				}
				add(ruleWait, position1535)
			}
			return true
		l1534:
			position, tokenIndex = position1534, tokenIndex1534
			return false
		},
		/* 127 DropOldest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('o' / 'O') ('l' / 'L') ('d' / 'D') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action97)> */
		func() bool {
			position1545, tokenIndex1545 := position, tokenIndex
			{
				position1546 := position
				{
					position1547 := position
					{
						position1548, tokenIndex1548 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1549
						}
						position++
						goto l1548
					l1549:
						position, tokenIndex = position1548, tokenIndex1548
						if buffer[position] != rune('D') {
							goto l1545
						}
						position++
					}
				l1548:
					{
						position1550, tokenIndex1550 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1551
						}
						position++
						goto l1550
					l1551:
						position, tokenIndex = position1550, tokenIndex1550
						if buffer[position] != rune('R') {
							goto l1545
						}
						position++
					}
				l1550:
					{
						position1552, tokenIndex1552 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1553
						}
						position++
						goto l1552
					l1553:
						position, tokenIndex = position1552, tokenIndex1552
						if buffer[position] != rune('O') {
							goto l1545
						}
						position++
					}
				l1552:
					{
						position1554, tokenIndex1554 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1555
						}
						position++
						goto l1554
					l1555:
						position, tokenIndex = position1554, tokenIndex1554
						if buffer[position] != rune('P') {
							goto l1545
						}
						position++
					}
				l1554:
					if !_rules[rulesp]() {
						goto l1545
					}
					{
						position1556, tokenIndex1556 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1557
						}
						position++
						goto l1556
					l1557:
						position, tokenIndex = position1556, tokenIndex1556
						if buffer[position] != rune('O') {
							goto l1545
						}
						position++
					}
				l1556:
					{
						position1558, tokenIndex1558 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1559
						}
						position++
						goto l1558
					l1559:
						position, tokenIndex = position1558, tokenIndex1558
						if buffer[position] != rune('L') {
							goto l1545
						}
						position++
					}
				l1558:
					{
						position1560, tokenIndex1560 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1561
						}
						position++
						goto l1560
					l1561:
						position, tokenIndex = position1560, tokenIndex1560
						if buffer[position] != rune('D') {
							goto l1545
						}
						position++
					}
				l1560:
					{
						position1562, tokenIndex1562 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1563
						}
						position++
						goto l1562
					l1563:
						position, tokenIndex = position1562, tokenIndex1562
						if buffer[position] != rune('E') {
							goto l1545
						}
						position++
					}
				l1562:
					{
						position1564, tokenIndex1564 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1565
						}
						position++
						goto l1564
					l1565:
						position, tokenIndex = position1564, tokenIndex1564
						if buffer[position] != rune('S') {
							goto l1545
						}
						position++
					}
				l1564:
					{
						position1566, tokenIndex1566 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1567
						}
						position++
						goto l1566
					l1567:
						position, tokenIndex = position1566, tokenIndex1566
						if buffer[position] != rune('T') {
							goto l1545
						}
						position++
					}
				l1566:
					add(rulePegText, position1547)
				}
				if !_rules[ruleAction97]() {
					goto l1545
				}
				add(ruleDropOldest, position1546)
			}
			return true
		l1545:
			position, tokenIndex = position1545, tokenIndex1545
			return false
		},
		/* 128 DropNewest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('n' / 'N') ('e' / 'E') ('w' / 'W') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action98)> */
		func() bool {
			position1568, tokenIndex1568 := position, tokenIndex
			{
				position1569 := position
				{
					position1570 := position
					{
						position1571, tokenIndex1571 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1572
						}
						position++
						goto l1571
					l1572:
						position, tokenIndex = position1571, tokenIndex1571
						if buffer[position] != rune('D') {
							goto l1568
						}
						position++
					}
				l1571:
					{
						position1573, tokenIndex1573 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1574
						}
						position++
						goto l1573
					l1574:
						position, tokenIndex = position1573, tokenIndex1573
						if buffer[position] != rune('R') {
							goto l1568
						}
						position++
					}
				l1573:
					{
						position1575, tokenIndex1575 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1576
						}
						position++
						goto l1575
					l1576:
						position, tokenIndex = position1575, tokenIndex1575
						if buffer[position] != rune('O') {
							goto l1568
						}
						position++
					}
				l1575:
					{
						position1577, tokenIndex1577 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1578
						}
						position++
						goto l1577
					l1578:
						position, tokenIndex = position1577, tokenIndex1577
						if buffer[position] != rune('P') {
							goto l1568
						}
						position++
					}
				l1577:
					if !_rules[rulesp]() {
						goto l1568
					}
					{
						position1579, tokenIndex1579 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1580
						}
						position++
						goto l1579
					l1580:
						position, tokenIndex = position1579, tokenIndex1579
						if buffer[position] != rune('N') {
							goto l1568
						}
						position++
					}
				l1579:
					{
						position1581, tokenIndex1581 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1582
						}
						position++
						goto l1581
					l1582:
						position, tokenIndex = position1581, tokenIndex1581
						if buffer[position] != rune('E') {
							goto l1568
						}
						position++
					}
				l1581:
					{
						position1583, tokenIndex1583 := position, tokenIndex
						if buffer[position] != rune('w') {
							goto l1584
						}
						position++
						goto l1583
					l1584:
						position, tokenIndex = position1583, tokenIndex1583
						if buffer[position] != rune('W') {
							goto l1568
						}
						position++
					}
				l1583:
					{
						position1585, tokenIndex1585 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1586
						}
						position++
						goto l1585
					l1586:
						position, tokenIndex = position1585, tokenIndex1585
						if buffer[position] != rune('E') {
							goto l1568
						}
						position++
					}
				l1585:
					{
						position1587, tokenIndex1587 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1588
						}
						position++
						goto l1587
					l1588:
						position, tokenIndex = position1587, tokenIndex1587
						if buffer[position] != rune('S') {
							goto l1568
						}
						position++
					}
				l1587:
					{
						position1589, tokenIndex1589 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1590
						}
						position++
						goto l1589
					l1590:
						position, tokenIndex = position1589, tokenIndex1589
						if buffer[position] != rune('T') {
							goto l1568
						}
						position++
					}
				l1589:
					add(rulePegText, position1570)
				}
				if !_rules[ruleAction98]() {
					goto l1568
				}
				add(ruleDropNewest, position1569)
			}
			return true
		l1568:
			position, tokenIndex = position1568, tokenIndex1568
			return false
		},
		/* 129 StreamIdentifier <- <(<ident> Action99)> */
		func() bool {
			position1591, tokenIndex1591 := position, tokenIndex
			{
				position1592 := position
				{
					position1593 := position
					if !_rules[ruleident]() {
						goto l1591
					}
					add(rulePegText, position1593)
				}
				if !_rules[ruleAction99]() {
					goto l1591
				}
				add(ruleStreamIdentifier, position1592)
			}
			return true
		l1591:
			position, tokenIndex = position1591, tokenIndex1591
			return false
		},
		/* 130 SourceSinkType <- <(<ident> Action100)> */
		func() bool {
			position1594, tokenIndex1594 := position, tokenIndex
			{
				position1595 := position
				{
					position1596 := position
					if !_rules[ruleident]() {
						goto l1594
					}
					add(rulePegText, position1596)
				}
				if !_rules[ruleAction100]() {
					goto l1594
				}
				add(ruleSourceSinkType, position1595)
			}
			return true
		l1594:
			position, tokenIndex = position1594, tokenIndex1594
			return false
		},
		/* 131 SourceSinkParamKey <- <(<ident> Action101)> */
		func() bool {
			position1597, tokenIndex1597 := position, tokenIndex
			{
				position1598 := position
				{
					position1599 := position
					if !_rules[ruleident]() {
						goto l1597
					}
					add(rulePegText, position1599)
				}
				if !_rules[ruleAction101]() {
					goto l1597
				}
				add(ruleSourceSinkParamKey, position1598)
			}
			return true
		l1597:
			position, tokenIndex = position1597, tokenIndex1597
			return false
		},
		/* 132 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action102)> */
		func() bool {
			position1600, tokenIndex1600 := position, tokenIndex
			{
				position1601 := position
				{
					position1602 := position
					{
						position1603, tokenIndex1603 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1604
						}
						position++
						goto l1603
					l1604:
						position, tokenIndex = position1603, tokenIndex1603
						if buffer[position] != rune('P') {
							goto l1600
						}
						position++
					}
				l1603:
					{
						position1605, tokenIndex1605 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1606
						}
						position++
						goto l1605
					l1606:
						position, tokenIndex = position1605, tokenIndex1605
						if buffer[position] != rune('A') {
							goto l1600
						}
						position++
					}
				l1605:
					{
						position1607, tokenIndex1607 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1608
						}
						position++
						goto l1607
					l1608:
						position, tokenIndex = position1607, tokenIndex1607
						if buffer[position] != rune('U') {
							goto l1600
						}
						position++
					}
				l1607:
					{
						position1609, tokenIndex1609 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1610
						}
						position++
						goto l1609
					l1610:
						position, tokenIndex = position1609, tokenIndex1609
						if buffer[position] != rune('S') {
							goto l1600
						}
						position++
					}
				l1609:
					{
						position1611, tokenIndex1611 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1612
						}
						position++
						goto l1611
					l1612:
						position, tokenIndex = position1611, tokenIndex1611
						if buffer[position] != rune('E') {
							goto l1600
						}
						position++
					}
				l1611:
					{
						position1613, tokenIndex1613 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1614
						}
						position++
						goto l1613
					l1614:
						position, tokenIndex = position1613, tokenIndex1613
						if buffer[position] != rune('D') {
							goto l1600
						}
						position++
					}
				l1613:
					add(rulePegText, position1602)
				}
				if !_rules[ruleAction102]() {
					goto l1600
				}
				add(rulePaused, position1601)
			}
			return true
		l1600:
			position, tokenIndex = position1600, tokenIndex1600
			return false
		},
		/* 133 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action103)> */
		func() bool {
			position1615, tokenIndex1615 := position, tokenIndex
			{
				position1616 := position
				{
					position1617 := position
					{
						position1618, tokenIndex1618 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1619
						}
						position++
						goto l1618
					l1619:
						position, tokenIndex = position1618, tokenIndex1618
						if buffer[position] != rune('U') {
							goto l1615
						}
						position++
					}
				l1618:
					{
						position1620, tokenIndex1620 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1621
						}
						position++
						goto l1620
					l1621:
						position, tokenIndex = position1620, tokenIndex1620
						if buffer[position] != rune('N') {
							goto l1615
						}
						position++
					}
				l1620:
					{
						position1622, tokenIndex1622 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1623
						}
						position++
						goto l1622
					l1623:
						position, tokenIndex = position1622, tokenIndex1622
						if buffer[position] != rune('P') {
							goto l1615
						}
						position++
					}
				l1622:
					{
						position1624, tokenIndex1624 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1625
						}
						position++
						goto l1624
					l1625:
						position, tokenIndex = position1624, tokenIndex1624
						if buffer[position] != rune('A') {
							goto l1615
						}
						position++
					}
				l1624:
					{
						position1626, tokenIndex1626 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1627
						}
						position++
						goto l1626
					l1627:
						position, tokenIndex = position1626, tokenIndex1626
						if buffer[position] != rune('U') {
							goto l1615
						}
						position++
					}
				l1626:
					{
						position1628, tokenIndex1628 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1629
						}
						position++
						goto l1628
					l1629:
						position, tokenIndex = position1628, tokenIndex1628
						if buffer[position] != rune('S') {
							goto l1615
						}
						position++
					}
				l1628:
					{
						position1630, tokenIndex1630 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1631
						}
						position++
						goto l1630
					l1631:
						position, tokenIndex = position1630, tokenIndex1630
						if buffer[position] != rune('E') {
							goto l1615
						}
						position++
					}
				l1630:
					{
						position1632, tokenIndex1632 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1633
						}
						position++
						goto l1632
					l1633:
						position, tokenIndex = position1632, tokenIndex1632
						if buffer[position] != rune('D') {
							goto l1615
						}
						position++
					}
				l1632:
					add(rulePegText, position1617)
				}
				if !_rules[ruleAction103]() {
					goto l1615
				}
				add(ruleUnpaused, position1616)
			}
			return true
		l1615:
			position, tokenIndex = position1615, tokenIndex1615
			return false
		},
		/* 134 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action104)> */
		func() bool {
			position1634, tokenIndex1634 := position, tokenIndex
			{
				position1635 := position
				{
					position1636 := position
					{
						position1637, tokenIndex1637 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1638
						}
						position++
						goto l1637
					l1638:
						position, tokenIndex = position1637, tokenIndex1637
						if buffer[position] != rune('A') {
							goto l1634
						}
						position++
					}
				l1637:
					{
						position1639, tokenIndex1639 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1640
						}
						position++
						goto l1639
					l1640:
						position, tokenIndex = position1639, tokenIndex1639
						if buffer[position] != rune('S') {
							goto l1634
						}
						position++
					}
				l1639:
					{
						position1641, tokenIndex1641 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1642
						}
						position++
						goto l1641
					l1642:
						position, tokenIndex = position1641, tokenIndex1641
						if buffer[position] != rune('C') {
							goto l1634
						}
						position++
					}
				l1641:
					add(rulePegText, position1636)
				}
				if !_rules[ruleAction104]() {
					goto l1634
				}
				add(ruleAscending, position1635)
			}
			return true
		l1634:
			position, tokenIndex = position1634, tokenIndex1634
			return false
		},
		/* 135 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action105)> */
		func() bool {
			position1643, tokenIndex1643 := position, tokenIndex
			{
				position1644 := position
				{
					position1645 := position
					{
						position1646, tokenIndex1646 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1647
						}
						position++
						goto l1646
					l1647:
						position, tokenIndex = position1646, tokenIndex1646
						if buffer[position] != rune('D') {
							goto l1643
						}
						position++
					}
				l1646:
					{
						position1648, tokenIndex1648 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1649
						}
						position++
						goto l1648
					l1649:
						position, tokenIndex = position1648, tokenIndex1648
						if buffer[position] != rune('E') {
							goto l1643
						}
						position++
					}
				l1648:
					{
						position1650, tokenIndex1650 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1651
						}
						position++
						goto l1650
					l1651:
						position, tokenIndex = position1650, tokenIndex1650
						if buffer[position] != rune('S') {
							goto l1643
						}
						position++
					}
				l1650:
					{
						position1652, tokenIndex1652 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1653
						}
						position++
						goto l1652
					l1653:
						position, tokenIndex = position1652, tokenIndex1652
						if buffer[position] != rune('C') {
							goto l1643
						}
						position++
					}
				l1652:
					add(rulePegText, position1645)
				}
				if !_rules[ruleAction105]() {
					goto l1643
				}
				add(ruleDescending, position1644)
			}
			return true
		l1643:
			position, tokenIndex = position1643, tokenIndex1643
			return false
		},
		/* 136 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1654, tokenIndex1654 := position, tokenIndex
			{
				position1655 := position
				{
					position1656, tokenIndex1656 := position, tokenIndex
					if !_rules[ruleBool]() {
						goto l1657
					}
					goto l1656
				l1657:
					position, tokenIndex = position1656, tokenIndex1656
					if !_rules[ruleInt]() {
						goto l1658
					}
					goto l1656
				l1658:
					position, tokenIndex = position1656, tokenIndex1656
					if !_rules[ruleFloat]() {
						goto l1659
					}
					goto l1656
				l1659:
					position, tokenIndex = position1656, tokenIndex1656
					if !_rules[ruleString]() {
						goto l1660
					}
					goto l1656
				l1660:
					position, tokenIndex = position1656, tokenIndex1656
					if !_rules[ruleBlob]() {
						goto l1661
					}
					goto l1656
				l1661:
					position, tokenIndex = position1656, tokenIndex1656
					if !_rules[ruleTimestamp]() {
						goto l1662
					}
					goto l1656
				l1662:
					position, tokenIndex = position1656, tokenIndex1656
					if !_rules[ruleArray]() {
						goto l1663
					}
					goto l1656
				l1663:
					position, tokenIndex = position1656, tokenIndex1656
					if !_rules[ruleMap]() {
						goto l1654
					}
				}
			l1656:
				add(ruleType, position1655)
			}
			return true
		l1654:
			position, tokenIndex = position1654, tokenIndex1654
			return false
		},
		/* 137 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action106)> */
		func() bool {
			position1664, tokenIndex1664 := position, tokenIndex
			{
				position1665 := position
				{
					position1666 := position
					{
						position1667, tokenIndex1667 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1668
						}
						position++
						goto l1667
					l1668:
						position, tokenIndex = position1667, tokenIndex1667
						if buffer[position] != rune('B') {
							goto l1664
						}
						position++
					}
				l1667:
					{
						position1669, tokenIndex1669 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1670
						}
						position++
						goto l1669
					l1670:
						position, tokenIndex = position1669, tokenIndex1669
						if buffer[position] != rune('O') {
							goto l1664
						}
						position++
					}
				l1669:
					{
						position1671, tokenIndex1671 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1672
						}
						position++
						goto l1671
					l1672:
						position, tokenIndex = position1671, tokenIndex1671
						if buffer[position] != rune('O') {
							goto l1664
						}
						position++
					}
				l1671:
					{
						position1673, tokenIndex1673 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1674
						}
						position++
						goto l1673
					l1674:
						position, tokenIndex = position1673, tokenIndex1673
						if buffer[position] != rune('L') {
							goto l1664
						}
						position++
					}
				l1673:
					add(rulePegText, position1666)
				}
				if !_rules[ruleAction106]() {
					goto l1664
				}
				add(ruleBool, position1665)
			}
			return true
		l1664:
			position, tokenIndex = position1664, tokenIndex1664
			return false
		},
		/* 138 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action107)> */
		func() bool {
			position1675, tokenIndex1675 := position, tokenIndex
			{
				position1676 := position
				{
					position1677 := position
					{
						position1678, tokenIndex1678 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1679
						}
						position++
						goto l1678
					l1679:
						position, tokenIndex = position1678, tokenIndex1678
						if buffer[position] != rune('I') {
							goto l1675
						}
						position++
					}
				l1678:
					{
						position1680, tokenIndex1680 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1681
						}
						position++
						goto l1680
					l1681:
						position, tokenIndex = position1680, tokenIndex1680
						if buffer[position] != rune('N') {
							goto l1675
						}
						position++
					}
				l1680:
					{
						position1682, tokenIndex1682 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1683
						}
						position++
						goto l1682
					l1683:
						position, tokenIndex = position1682, tokenIndex1682
						if buffer[position] != rune('T') {
							goto l1675
						}
						position++
					}
				l1682:
					add(rulePegText, position1677)
				}
				if !_rules[ruleAction107]() {
					goto l1675
				}
				add(ruleInt, position1676)
			}
			return true
		l1675:
			position, tokenIndex = position1675, tokenIndex1675
			return false
		},
		/* 139 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action108)> */
		func() bool {
			position1684, tokenIndex1684 := position, tokenIndex
			{
				position1685 := position
				{
					position1686 := position
					{
						position1687, tokenIndex1687 := position, tokenIndex
						if buffer[position] != rune('f') {
							goto l1688
						}
						position++
						goto l1687
					l1688:
						position, tokenIndex = position1687, tokenIndex1687
						if buffer[position] != rune('F') {
							goto l1684
						}
						position++
					}
				l1687:
					{
						position1689, tokenIndex1689 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1690
						}
						position++
						goto l1689
					l1690:
						position, tokenIndex = position1689, tokenIndex1689
						if buffer[position] != rune('L') {
							goto l1684
						}
						position++
					}
				l1689:
					{
						position1691, tokenIndex1691 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1692
						}
						position++
						goto l1691
					l1692:
						position, tokenIndex = position1691, tokenIndex1691
						if buffer[position] != rune('O') {
							goto l1684
						}
						position++
					}
				l1691:
					{
						position1693, tokenIndex1693 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1694
						}
						position++
						goto l1693
					l1694:
						position, tokenIndex = position1693, tokenIndex1693
						if buffer[position] != rune('A') {
							goto l1684
						}
						position++
					}
				l1693:
					{
						position1695, tokenIndex1695 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1696
						}
						position++
						goto l1695
					l1696:
						position, tokenIndex = position1695, tokenIndex1695
						if buffer[position] != rune('T') {
							goto l1684
						}
						position++
					}
				l1695:
					add(rulePegText, position1686)
				}
				if !_rules[ruleAction108]() {
					goto l1684
				}
				add(ruleFloat, position1685)
			}
			return true
		l1684:
			position, tokenIndex = position1684, tokenIndex1684
			return false
		},
		/* 140 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action109)> */
		func() bool {
			position1697, tokenIndex1697 := position, tokenIndex
			{
				position1698 := position
				{
					position1699 := position
					{
						position1700, tokenIndex1700 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1701
						}
						position++
						goto l1700
					l1701:
						position, tokenIndex = position1700, tokenIndex1700
						if buffer[position] != rune('S') {
							goto l1697
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
							goto l1697
						}
						position++
					}
				l1702:
					{
						position1704, tokenIndex1704 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1705
						}
						position++
						goto l1704
					l1705:
						position, tokenIndex = position1704, tokenIndex1704
						if buffer[position] != rune('R') {
							goto l1697
						}
						position++
					}
				l1704:
					{
						position1706, tokenIndex1706 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1707
						}
						position++
						goto l1706
					l1707:
						position, tokenIndex = position1706, tokenIndex1706
						if buffer[position] != rune('I') {
							goto l1697
						}
						position++
					}
				l1706:
					{
						position1708, tokenIndex1708 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1709
						}
						position++
						goto l1708
					l1709:
						position, tokenIndex = position1708, tokenIndex1708
						if buffer[position] != rune('N') {
							goto l1697
						}
						position++
					}
				l1708:
					{
						position1710, tokenIndex1710 := position, tokenIndex
						if buffer[position] != rune('g') {
							goto l1711
						}
						position++
						goto l1710
					l1711:
						position, tokenIndex = position1710, tokenIndex1710
						if buffer[position] != rune('G') {
							goto l1697
						}
						position++
					}
				l1710:
					add(rulePegText, position1699)
				}
				if !_rules[ruleAction109]() {
					goto l1697
				}
				add(ruleString, position1698)
			}
			return true
		l1697:
			position, tokenIndex = position1697, tokenIndex1697
			return false
		},
		/* 141 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action110)> */
		func() bool {
			position1712, tokenIndex1712 := position, tokenIndex
			{
				position1713 := position
				{
					position1714 := position
					{
						position1715, tokenIndex1715 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1716
						}
						position++
						goto l1715
					l1716:
						position, tokenIndex = position1715, tokenIndex1715
						if buffer[position] != rune('B') {
							goto l1712
						}
						position++
					}
				l1715:
					{
						position1717, tokenIndex1717 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1718
						}
						position++
						goto l1717
					l1718:
						position, tokenIndex = position1717, tokenIndex1717
						if buffer[position] != rune('L') {
							goto l1712
						}
						position++
					}
				l1717:
					{
						position1719, tokenIndex1719 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1720
						}
						position++
						goto l1719
					l1720:
						position, tokenIndex = position1719, tokenIndex1719
						if buffer[position] != rune('O') {
							goto l1712
						}
						position++
					}
				l1719:
					{
						position1721, tokenIndex1721 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1722
						}
						position++
						goto l1721
					l1722:
						position, tokenIndex = position1721, tokenIndex1721
						if buffer[position] != rune('B') {
							goto l1712
						}
						position++
					}
				l1721:
					add(rulePegText, position1714)
				}
				if !_rules[ruleAction110]() {
					goto l1712
				}
				add(ruleBlob, position1713)
			}
			return true
		l1712:
			position, tokenIndex = position1712, tokenIndex1712
			return false
		},
		/* 142 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action111)> */
		func() bool {
			position1723, tokenIndex1723 := position, tokenIndex
			{
				position1724 := position
				{
					position1725 := position
					{
						position1726, tokenIndex1726 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1727
						}
						position++
						goto l1726
					l1727:
						position, tokenIndex = position1726, tokenIndex1726
						if buffer[position] != rune('T') {
							goto l1723
						}
						position++
					}
				l1726:
					{
						position1728, tokenIndex1728 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1729
						}
						position++
						goto l1728
					l1729:
						position, tokenIndex = position1728, tokenIndex1728
						if buffer[position] != rune('I') {
							goto l1723
						}
						position++
					}
				l1728:
					{
						position1730, tokenIndex1730 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1731
						}
						position++
						goto l1730
					l1731:
						position, tokenIndex = position1730, tokenIndex1730
						if buffer[position] != rune('M') {
							goto l1723
						}
						position++
					}
				l1730:
					{
						position1732, tokenIndex1732 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1733
						}
						position++
						goto l1732
					l1733:
						position, tokenIndex = position1732, tokenIndex1732
						if buffer[position] != rune('E') {
							goto l1723
						}
						position++
					}
				l1732:
					{
						position1734, tokenIndex1734 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1735
						}
						position++
						goto l1734
					l1735:
						position, tokenIndex = position1734, tokenIndex1734
						if buffer[position] != rune('S') {
							goto l1723
						}
						position++
					}
				l1734:
					{
						position1736, tokenIndex1736 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1737
						}
						position++
						goto l1736
					l1737:
						position, tokenIndex = position1736, tokenIndex1736
						if buffer[position] != rune('T') {
							goto l1723
						}
						position++
					}
				l1736:
					{
						position1738, tokenIndex1738 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1739
						}
						position++
						goto l1738
					l1739:
						position, tokenIndex = position1738, tokenIndex1738
						if buffer[position] != rune('A') {
							goto l1723
						}
						position++
					}
				l1738:
					{
						position1740, tokenIndex1740 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1741
						}
						position++
						goto l1740
					l1741:
						position, tokenIndex = position1740, tokenIndex1740
						if buffer[position] != rune('M') {
							goto l1723
						}
						position++
					}
				l1740:
					{
						position1742, tokenIndex1742 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1743
						}
						position++
						goto l1742
					l1743:
						position, tokenIndex = position1742, tokenIndex1742
						if buffer[position] != rune('P') {
							goto l1723
						}
						position++
					}
				l1742:
					add(rulePegText, position1725)
				}
				if !_rules[ruleAction111]() {
					goto l1723
				}
				add(ruleTimestamp, position1724)
			}
			return true
		l1723:
			position, tokenIndex = position1723, tokenIndex1723
			return false
		},
		/* 143 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action112)> */
		func() bool {
			position1744, tokenIndex1744 := position, tokenIndex
			{
				position1745 := position
				{
					position1746 := position
					{
						position1747, tokenIndex1747 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1748
						}
						position++
						goto l1747
					l1748:
						position, tokenIndex = position1747, tokenIndex1747
						if buffer[position] != rune('A') {
							goto l1744
						}
						position++
					}
				l1747:
					{
						position1749, tokenIndex1749 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1750
						}
						position++
						goto l1749
					l1750:
						position, tokenIndex = position1749, tokenIndex1749
						if buffer[position] != rune('R') {
							goto l1744
						}
						position++
					}
				l1749:
					{
						position1751, tokenIndex1751 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1752
						}
						position++
						goto l1751
					l1752:
						position, tokenIndex = position1751, tokenIndex1751
						if buffer[position] != rune('R') {
							goto l1744
						}
						position++
					}
				l1751:
					{
						position1753, tokenIndex1753 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1754
						}
						position++
						goto l1753
					l1754:
						position, tokenIndex = position1753, tokenIndex1753
						if buffer[position] != rune('A') {
							goto l1744
						}
						position++
					}
				l1753:
					{
						position1755, tokenIndex1755 := position, tokenIndex
						if buffer[position] != rune('y') {
							goto l1756
						}
						position++
						goto l1755
					l1756:
						position, tokenIndex = position1755, tokenIndex1755
						if buffer[position] != rune('Y') {
							goto l1744
						}
						position++
					}
				l1755:
					add(rulePegText, position1746)
				}
				if !_rules[ruleAction112]() {
					goto l1744
				}
				add(ruleArray, position1745)
			}
			return true
		l1744:
			position, tokenIndex = position1744, tokenIndex1744
			return false
		},
		/* 144 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action113)> */
		func() bool {
			position1757, tokenIndex1757 := position, tokenIndex
			{
				position1758 := position
				{
					position1759 := position
					{
						position1760, tokenIndex1760 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1761
						}
						position++
						goto l1760
					l1761:
						position, tokenIndex = position1760, tokenIndex1760
						if buffer[position] != rune('M') {
							goto l1757
						}
						position++
					}
				l1760:
					{
						position1762, tokenIndex1762 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1763
						}
						position++
						goto l1762
					l1763:
						position, tokenIndex = position1762, tokenIndex1762
						if buffer[position] != rune('A') {
							goto l1757
						}
						position++
					}
				l1762:
					{
						position1764, tokenIndex1764 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1765
						}
						position++
						goto l1764
					l1765:
						position, tokenIndex = position1764, tokenIndex1764
						if buffer[position] != rune('P') {
							goto l1757
						}
						position++
					}
				l1764:
					add(rulePegText, position1759)
				}
				if !_rules[ruleAction113]() {
					goto l1757
				}
				add(ruleMap, position1758)
			}
			return true
		l1757:
			position, tokenIndex = position1757, tokenIndex1757
			return false
		},
		/* 145 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action114)> */
		func() bool {
			position1766, tokenIndex1766 := position, tokenIndex
			{
				position1767 := position
				{
					position1768 := position
					{
						position1769, tokenIndex1769 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1770
						}
						position++
						goto l1769
					l1770:
						position, tokenIndex = position1769, tokenIndex1769
						if buffer[position] != rune('O') {
							goto l1766
						}
						position++
					}
				l1769:
					{
						position1771, tokenIndex1771 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1772
						}
						position++
						goto l1771
					l1772:
						position, tokenIndex = position1771, tokenIndex1771
						if buffer[position] != rune('R') {
							goto l1766
						}
						position++
					}
				l1771:
					add(rulePegText, position1768)
				}
				if !_rules[ruleAction114]() {
					goto l1766
				}
				add(ruleOr, position1767)
			}
			return true
		l1766:
			position, tokenIndex = position1766, tokenIndex1766
			return false
		},
		/* 146 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action115)> */
		func() bool {
			position1773, tokenIndex1773 := position, tokenIndex
			{
				position1774 := position
				{
					position1775 := position
					{
						position1776, tokenIndex1776 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1777
						}
						position++
						goto l1776
					l1777:
						position, tokenIndex = position1776, tokenIndex1776
						if buffer[position] != rune('A') {
							goto l1773
						}
						position++
					}
				l1776:
					{
						position1778, tokenIndex1778 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1779
						}
						position++
						goto l1778
					l1779:
						position, tokenIndex = position1778, tokenIndex1778
						if buffer[position] != rune('N') {
							goto l1773
						}
						position++
					}
				l1778:
					{
						position1780, tokenIndex1780 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1781
						}
						position++
						goto l1780
					l1781:
						position, tokenIndex = position1780, tokenIndex1780
						if buffer[position] != rune('D') {
							goto l1773
						}
						position++
					}
				l1780:
					add(rulePegText, position1775)
				}
				if !_rules[ruleAction115]() {
					goto l1773
				}
				add(ruleAnd, position1774)
			}
			return true
		l1773:
			position, tokenIndex = position1773, tokenIndex1773
			return false
		},
		/* 147 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action116)> */
		func() bool {
			position1782, tokenIndex1782 := position, tokenIndex
			{
				position1783 := position
				{
					position1784 := position
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
							goto l1782
						}
						position++
					}
				l1785:
					{
						position1787, tokenIndex1787 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1788
						}
						position++
						goto l1787
					l1788:
						position, tokenIndex = position1787, tokenIndex1787
						if buffer[position] != rune('O') {
							goto l1782
						}
						position++
					}
				l1787:
					{
						position1789, tokenIndex1789 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1790
						}
						position++
						goto l1789
					l1790:
						position, tokenIndex = position1789, tokenIndex1789
						if buffer[position] != rune('T') {
							goto l1782
						}
						position++
					}
				l1789:
					add(rulePegText, position1784)
				}
				if !_rules[ruleAction116]() {
					goto l1782
				}
				add(ruleNot, position1783)
			}
			return true
		l1782:
			position, tokenIndex = position1782, tokenIndex1782
			return false
		},
		/* 148 Equal <- <(<'='> Action117)> */
		func() bool {
			position1791, tokenIndex1791 := position, tokenIndex
			{
				position1792 := position
				{
					position1793 := position
					if buffer[position] != rune('=') {
						goto l1791
					}
					position++
					add(rulePegText, position1793)
				}
				if !_rules[ruleAction117]() {
					goto l1791
				}
				add(ruleEqual, position1792)
			}
			return true
		l1791:
			position, tokenIndex = position1791, tokenIndex1791
			return false
		},
		/* 149 Less <- <(<'<'> Action118)> */
		func() bool {
			position1794, tokenIndex1794 := position, tokenIndex
			{
				position1795 := position
				{
					position1796 := position
					if buffer[position] != rune('<') {
						goto l1794
					}
					position++
					add(rulePegText, position1796)
				}
				if !_rules[ruleAction118]() {
					goto l1794
				}
				add(ruleLess, position1795)
			}
			return true
		l1794:
			position, tokenIndex = position1794, tokenIndex1794
			return false
		},
		/* 150 LessOrEqual <- <(<('<' '=')> Action119)> */
		func() bool {
			position1797, tokenIndex1797 := position, tokenIndex
			{
				position1798 := position
				{
					position1799 := position
					if buffer[position] != rune('<') {
						goto l1797
					}
					position++
					if buffer[position] != rune('=') {
						goto l1797
					}
					position++
					add(rulePegText, position1799)
				}
				if !_rules[ruleAction119]() {
					goto l1797
				}
				add(ruleLessOrEqual, position1798)
			}
			return true
		l1797:
			position, tokenIndex = position1797, tokenIndex1797
			return false
		},
		/* 151 Greater <- <(<'>'> Action120)> */
		func() bool {
			position1800, tokenIndex1800 := position, tokenIndex
			{
				position1801 := position
				{
					position1802 := position
					if buffer[position] != rune('>') {
						goto l1800
					}
					position++
					add(rulePegText, position1802)
				}
				if !_rules[ruleAction120]() {
					goto l1800
				}
				add(ruleGreater, position1801)
			}
			return true
		l1800:
			position, tokenIndex = position1800, tokenIndex1800
			return false
		},
		/* 152 GreaterOrEqual <- <(<('>' '=')> Action121)> */
		func() bool {
			position1803, tokenIndex1803 := position, tokenIndex
			{
				position1804 := position
				{
					position1805 := position
					if buffer[position] != rune('>') {
						goto l1803
					}
					position++
					if buffer[position] != rune('=') {
						goto l1803
					}
					position++
					add(rulePegText, position1805)
				}
				if !_rules[ruleAction121]() {
					goto l1803
				}
				add(ruleGreaterOrEqual, position1804)
			}
			return true
		l1803:
			position, tokenIndex = position1803, tokenIndex1803
			return false
		},
		/* 153 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action122)> */
		func() bool {
			position1806, tokenIndex1806 := position, tokenIndex
			{
				position1807 := position
				{
					position1808 := position
					{
						position1809, tokenIndex1809 := position, tokenIndex
						if buffer[position] != rune('!') {
							goto l1810
						}
						position++
						if buffer[position] != rune('=') {
							goto l1810
						}
						position++
						goto l1809
					l1810:
						position, tokenIndex = position1809, tokenIndex1809
						if buffer[position] != rune('<') {
							goto l1806
						}
						position++
						if buffer[position] != rune('>') {
							goto l1806
						}
						position++
					}
				l1809:
					add(rulePegText, position1808)
				}
				if !_rules[ruleAction122]() {
					goto l1806
				}
				add(ruleNotEqual, position1807)
			}
			return true
		l1806:
			position, tokenIndex = position1806, tokenIndex1806
			return false
		},
		/* 154 Concat <- <(<('|' '|')> Action123)> */
		func() bool {
			position1811, tokenIndex1811 := position, tokenIndex
			{
				position1812 := position
				{
					position1813 := position
					if buffer[position] != rune('|') {
						goto l1811
					}
					position++
					if buffer[position] != rune('|') {
						goto l1811
					}
					position++
					add(rulePegText, position1813)
				}
				if !_rules[ruleAction123]() {
					goto l1811
				}
				add(ruleConcat, position1812)
			}
			return true
		l1811:
			position, tokenIndex = position1811, tokenIndex1811
			return false
		},
		/* 155 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action124)> */
		func() bool {
			position1814, tokenIndex1814 := position, tokenIndex
			{
				position1815 := position
				{
					position1816 := position
					{
						position1817, tokenIndex1817 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1818
						}
						position++
						goto l1817
					l1818:
						position, tokenIndex = position1817, tokenIndex1817
						if buffer[position] != rune('I') {
							goto l1814
						}
						position++
					}
				l1817:
					{
						position1819, tokenIndex1819 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1820
						}
						position++
						goto l1819
					l1820:
						position, tokenIndex = position1819, tokenIndex1819
						if buffer[position] != rune('S') {
							goto l1814
						}
						position++
					}
				l1819:
					add(rulePegText, position1816)
				}
				if !_rules[ruleAction124]() {
					goto l1814
				}
				add(ruleIs, position1815)
			}
			return true
		l1814:
			position, tokenIndex = position1814, tokenIndex1814
			return false
		},
		/* 156 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action125)> */
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
					if !_rules[rulesp]() {
						goto l1821
					}
					{
						position1828, tokenIndex1828 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1829
						}
						position++
						goto l1828
					l1829:
						position, tokenIndex = position1828, tokenIndex1828
						if buffer[position] != rune('N') {
							goto l1821
						}
						position++
					}
				l1828:
					{
						position1830, tokenIndex1830 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1831
						}
						position++
						goto l1830
					l1831:
						position, tokenIndex = position1830, tokenIndex1830
						if buffer[position] != rune('O') {
							goto l1821
						}
						position++
					}
				l1830:
					{
						position1832, tokenIndex1832 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1833
						}
						position++
						goto l1832
					l1833:
						position, tokenIndex = position1832, tokenIndex1832
						if buffer[position] != rune('T') {
							goto l1821
						}
						position++
					}
				l1832:
					add(rulePegText, position1823)
				}
				if !_rules[ruleAction125]() {
					goto l1821
				}
				add(ruleIsNot, position1822)
			}
			return true
		l1821:
			position, tokenIndex = position1821, tokenIndex1821
			return false
		},
		/* 157 Plus <- <(<'+'> Action126)> */
		func() bool {
			position1834, tokenIndex1834 := position, tokenIndex
			{
				position1835 := position
				{
					position1836 := position
					if buffer[position] != rune('+') {
						goto l1834
					}
					position++
					add(rulePegText, position1836)
				}
				if !_rules[ruleAction126]() {
					goto l1834
				}
				add(rulePlus, position1835)
			}
			return true
		l1834:
			position, tokenIndex = position1834, tokenIndex1834
			return false
		},
		/* 158 Minus <- <(<'-'> Action127)> */
		func() bool {
			position1837, tokenIndex1837 := position, tokenIndex
			{
				position1838 := position
				{
					position1839 := position
					if buffer[position] != rune('-') {
						goto l1837
					}
					position++
					add(rulePegText, position1839)
				}
				if !_rules[ruleAction127]() {
					goto l1837
				}
				add(ruleMinus, position1838)
			}
			return true
		l1837:
			position, tokenIndex = position1837, tokenIndex1837
			return false
		},
		/* 159 Multiply <- <(<'*'> Action128)> */
		func() bool {
			position1840, tokenIndex1840 := position, tokenIndex
			{
				position1841 := position
				{
					position1842 := position
					if buffer[position] != rune('*') {
						goto l1840
					}
					position++
					add(rulePegText, position1842)
				}
				if !_rules[ruleAction128]() {
					goto l1840
				}
				add(ruleMultiply, position1841)
			}
			return true
		l1840:
			position, tokenIndex = position1840, tokenIndex1840
			return false
		},
		/* 160 Divide <- <(<'/'> Action129)> */
		func() bool {
			position1843, tokenIndex1843 := position, tokenIndex
			{
				position1844 := position
				{
					position1845 := position
					if buffer[position] != rune('/') {
						goto l1843
					}
					position++
					add(rulePegText, position1845)
				}
				if !_rules[ruleAction129]() {
					goto l1843
				}
				add(ruleDivide, position1844)
			}
			return true
		l1843:
			position, tokenIndex = position1843, tokenIndex1843
			return false
		},
		/* 161 Modulo <- <(<'%'> Action130)> */
		func() bool {
			position1846, tokenIndex1846 := position, tokenIndex
			{
				position1847 := position
				{
					position1848 := position
					if buffer[position] != rune('%') {
						goto l1846
					}
					position++
					add(rulePegText, position1848)
				}
				if !_rules[ruleAction130]() {
					goto l1846
				}
				add(ruleModulo, position1847)
			}
			return true
		l1846:
			position, tokenIndex = position1846, tokenIndex1846
			return false
		},
		/* 162 UnaryMinus <- <(<'-'> Action131)> */
		func() bool {
			position1849, tokenIndex1849 := position, tokenIndex
			{
				position1850 := position
				{
					position1851 := position
					if buffer[position] != rune('-') {
						goto l1849
					}
					position++
					add(rulePegText, position1851)
				}
				if !_rules[ruleAction131]() {
					goto l1849
				}
				add(ruleUnaryMinus, position1850)
			}
			return true
		l1849:
			position, tokenIndex = position1849, tokenIndex1849
			return false
		},
		/* 163 Identifier <- <(<ident> Action132)> */
		func() bool {
			position1852, tokenIndex1852 := position, tokenIndex
			{
				position1853 := position
				{
					position1854 := position
					if !_rules[ruleident]() {
						goto l1852
					}
					add(rulePegText, position1854)
				}
				if !_rules[ruleAction132]() {
					goto l1852
				}
				add(ruleIdentifier, position1853)
			}
			return true
		l1852:
			position, tokenIndex = position1852, tokenIndex1852
			return false
		},
		/* 164 TargetIdentifier <- <(<('*' / jsonSetPath)> Action133)> */
		func() bool {
			position1855, tokenIndex1855 := position, tokenIndex
			{
				position1856 := position
				{
					position1857 := position
					{
						position1858, tokenIndex1858 := position, tokenIndex
						if buffer[position] != rune('*') {
							goto l1859
						}
						position++
						goto l1858
					l1859:
						position, tokenIndex = position1858, tokenIndex1858
						if !_rules[rulejsonSetPath]() {
							goto l1855
						}
					}
				l1858:
					add(rulePegText, position1857)
				}
				if !_rules[ruleAction133]() {
					goto l1855
				}
				add(ruleTargetIdentifier, position1856)
			}
			return true
		l1855:
			position, tokenIndex = position1855, tokenIndex1855
			return false
		},
		/* 165 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1860, tokenIndex1860 := position, tokenIndex
			{
				position1861 := position
				{
					position1862, tokenIndex1862 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1863
					}
					position++
					goto l1862
				l1863:
					position, tokenIndex = position1862, tokenIndex1862
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1860
					}
					position++
				}
			l1862:
			l1864:
				{
					position1865, tokenIndex1865 := position, tokenIndex
					{
						position1866, tokenIndex1866 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1867
						}
						position++
						goto l1866
					l1867:
						position, tokenIndex = position1866, tokenIndex1866
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1868
						}
						position++
						goto l1866
					l1868:
						position, tokenIndex = position1866, tokenIndex1866
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1869
						}
						position++
						goto l1866
					l1869:
						position, tokenIndex = position1866, tokenIndex1866
						if buffer[position] != rune('_') {
							goto l1865
						}
						position++
					}
				l1866:
					goto l1864
				l1865:
					position, tokenIndex = position1865, tokenIndex1865
				}
				add(ruleident, position1861)
			}
			return true
		l1860:
			position, tokenIndex = position1860, tokenIndex1860
			return false
		},
		/* 166 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1870, tokenIndex1870 := position, tokenIndex
			{
				position1871 := position
				if !_rules[rulejsonPathHead]() {
					goto l1870
				}
			l1872:
				{
					position1873, tokenIndex1873 := position, tokenIndex
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1873
					}
					goto l1872
				l1873:
					position, tokenIndex = position1873, tokenIndex1873
				}
				add(rulejsonGetPath, position1871)
			}
			return true
		l1870:
			position, tokenIndex = position1870, tokenIndex1870
			return false
		},
		/* 167 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1874, tokenIndex1874 := position, tokenIndex
			{
				position1875 := position
				if !_rules[rulejsonPathHead]() {
					goto l1874
				}
			l1876:
				{
					position1877, tokenIndex1877 := position, tokenIndex
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1877
					}
					goto l1876
				l1877:
					position, tokenIndex = position1877, tokenIndex1877
				}
				add(rulejsonSetPath, position1875)
			}
			return true
		l1874:
			position, tokenIndex = position1874, tokenIndex1874
			return false
		},
		/* 168 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1878, tokenIndex1878 := position, tokenIndex
			{
				position1879 := position
				{
					position1880, tokenIndex1880 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l1881
					}
					goto l1880
				l1881:
					position, tokenIndex = position1880, tokenIndex1880
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1878
					}
				}
			l1880:
				add(rulejsonPathHead, position1879)
			}
			return true
		l1878:
			position, tokenIndex = position1878, tokenIndex1878
			return false
		},
		/* 169 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1882, tokenIndex1882 := position, tokenIndex
			{
				position1883 := position
				{
					position1884, tokenIndex1884 := position, tokenIndex
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1885
					}
					goto l1884
				l1885:
					position, tokenIndex = position1884, tokenIndex1884
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1886
					}
					goto l1884
				l1886:
					position, tokenIndex = position1884, tokenIndex1884
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1887
					}
					goto l1884
				l1887:
					position, tokenIndex = position1884, tokenIndex1884
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1888
					}
					goto l1884
				l1888:
					position, tokenIndex = position1884, tokenIndex1884
					if !_rules[rulejsonArraySlice]() {
						goto l1889
					}
					goto l1884
				l1889:
					position, tokenIndex = position1884, tokenIndex1884
					if !_rules[rulejsonArrayAccess]() {
						goto l1882
					}
				}
			l1884:
				add(rulejsonGetPathNonHead, position1883)
			}
			return true
		l1882:
			position, tokenIndex = position1882, tokenIndex1882
			return false
		},
		/* 170 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1890, tokenIndex1890 := position, tokenIndex
			{
				position1891 := position
				{
					position1892, tokenIndex1892 := position, tokenIndex
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1893
					}
					goto l1892
				l1893:
					position, tokenIndex = position1892, tokenIndex1892
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1890
					}
				}
			l1892:
				add(rulejsonSetPathNonHead, position1891)
			}
			return true
		l1890:
			position, tokenIndex = position1890, tokenIndex1890
			return false
		},
		/* 171 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1894, tokenIndex1894 := position, tokenIndex
			{
				position1895 := position
				{
					position1896, tokenIndex1896 := position, tokenIndex
					if buffer[position] != rune('.') {
						goto l1897
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1897
					}
					goto l1896
				l1897:
					position, tokenIndex = position1896, tokenIndex1896
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1894
					}
				}
			l1896:
				add(rulejsonMapSingleLevel, position1895)
			}
			return true
		l1894:
			position, tokenIndex = position1894, tokenIndex1894
			return false
		},
		/* 172 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1898, tokenIndex1898 := position, tokenIndex
			{
				position1899 := position
				if buffer[position] != rune('.') {
					goto l1898
				}
				position++
				if buffer[position] != rune('.') {
					goto l1898
				}
				position++
				{
					position1900, tokenIndex1900 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l1901
					}
					goto l1900
				l1901:
					position, tokenIndex = position1900, tokenIndex1900
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1898
					}
				}
			l1900:
				add(rulejsonMapMultipleLevel, position1899)
			}
			return true
		l1898:
			position, tokenIndex = position1898, tokenIndex1898
			return false
		},
		/* 173 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1902, tokenIndex1902 := position, tokenIndex
			{
				position1903 := position
				{
					position1904 := position
					{
						position1905, tokenIndex1905 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1906
						}
						position++
						goto l1905
					l1906:
						position, tokenIndex = position1905, tokenIndex1905
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1902
						}
						position++
					}
				l1905:
				l1907:
					{
						position1908, tokenIndex1908 := position, tokenIndex
						{
							position1909, tokenIndex1909 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1910
							}
							position++
							goto l1909
						l1910:
							position, tokenIndex = position1909, tokenIndex1909
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1911
							}
							position++
							goto l1909
						l1911:
							position, tokenIndex = position1909, tokenIndex1909
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1912
							}
							position++
							goto l1909
						l1912:
							position, tokenIndex = position1909, tokenIndex1909
							if buffer[position] != rune('_') {
								goto l1908
							}
							position++
						}
					l1909:
						goto l1907
					l1908:
						position, tokenIndex = position1908, tokenIndex1908
					}
					add(rulePegText, position1904)
				}
				add(rulejsonMapAccessString, position1903)
			}
			return true
		l1902:
			position, tokenIndex = position1902, tokenIndex1902
			return false
		},
		/* 174 jsonMapAccessBracket <- <('[' doubleQuotedString ']')> */
		func() bool {
			position1913, tokenIndex1913 := position, tokenIndex
			{
				position1914 := position
				if buffer[position] != rune('[') {
					goto l1913
				}
				position++
				if !_rules[ruledoubleQuotedString]() {
					goto l1913
				}
				if buffer[position] != rune(']') {
					goto l1913
				}
				position++
				add(rulejsonMapAccessBracket, position1914)
			}
			return true
		l1913:
			position, tokenIndex = position1913, tokenIndex1913
			return false
		},
		/* 175 doubleQuotedString <- <('"' <(('"' '"') / (!'"' .))*> '"')> */
		func() bool {
			position1915, tokenIndex1915 := position, tokenIndex
			{
				position1916 := position
				if buffer[position] != rune('"') {
					goto l1915
				}
				position++
				{
					position1917 := position
				l1918:
					{
						position1919, tokenIndex1919 := position, tokenIndex
						{
							position1920, tokenIndex1920 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l1921
							}
							position++
							if buffer[position] != rune('"') {
								goto l1921
							}
							position++
							goto l1920
						l1921:
							position, tokenIndex = position1920, tokenIndex1920
							{
								position1922, tokenIndex1922 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l1922
								}
								position++
								goto l1919
							l1922:
								position, tokenIndex = position1922, tokenIndex1922
							}
							if !matchDot() {
								goto l1919
							}
						}
					l1920:
						goto l1918
					l1919:
						position, tokenIndex = position1919, tokenIndex1919
					}
					add(rulePegText, position1917)
				}
				if buffer[position] != rune('"') {
					goto l1915
				}
				position++
				add(ruledoubleQuotedString, position1916)
			}
			return true
		l1915:
			position, tokenIndex = position1915, tokenIndex1915
			return false
		},
		/* 176 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1923, tokenIndex1923 := position, tokenIndex
			{
				position1924 := position
				if buffer[position] != rune('[') {
					goto l1923
				}
				position++
				{
					position1925 := position
					{
						position1926, tokenIndex1926 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1926
						}
						position++
						goto l1927
					l1926:
						position, tokenIndex = position1926, tokenIndex1926
					}
				l1927:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1923
					}
					position++
				l1928:
					{
						position1929, tokenIndex1929 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1929
						}
						position++
						goto l1928
					l1929:
						position, tokenIndex = position1929, tokenIndex1929
					}
					add(rulePegText, position1925)
				}
				if buffer[position] != rune(']') {
					goto l1923
				}
				position++
				add(rulejsonArrayAccess, position1924)
			}
			return true
		l1923:
			position, tokenIndex = position1923, tokenIndex1923
			return false
		},
		/* 177 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
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
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1930
					}
					position++
				l1933:
					{
						position1934, tokenIndex1934 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1934
						}
						position++
						goto l1933
					l1934:
						position, tokenIndex = position1934, tokenIndex1934
					}
					add(rulePegText, position1932)
				}
				if buffer[position] != rune(']') {
					goto l1930
				}
				position++
				add(rulejsonNonNegativeArrayAccess, position1931)
			}
			return true
		l1930:
			position, tokenIndex = position1930, tokenIndex1930
			return false
		},
		/* 178 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1935, tokenIndex1935 := position, tokenIndex
			{
				position1936 := position
				if buffer[position] != rune('[') {
					goto l1935
				}
				position++
				{
					position1937 := position
					{
						position1938, tokenIndex1938 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1938
						}
						position++
						goto l1939
					l1938:
						position, tokenIndex = position1938, tokenIndex1938
					}
				l1939:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1935
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
					if buffer[position] != rune(':') {
						goto l1935
					}
					position++
					{
						position1942, tokenIndex1942 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1942
						}
						position++
						goto l1943
					l1942:
						position, tokenIndex = position1942, tokenIndex1942
					}
				l1943:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1935
					}
					position++
				l1944:
					{
						position1945, tokenIndex1945 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1945
						}
						position++
						goto l1944
					l1945:
						position, tokenIndex = position1945, tokenIndex1945
					}
					{
						position1946, tokenIndex1946 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l1946
						}
						position++
						{
							position1948, tokenIndex1948 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1948
							}
							position++
							goto l1949
						l1948:
							position, tokenIndex = position1948, tokenIndex1948
						}
					l1949:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1946
						}
						position++
					l1950:
						{
							position1951, tokenIndex1951 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1951
							}
							position++
							goto l1950
						l1951:
							position, tokenIndex = position1951, tokenIndex1951
						}
						goto l1947
					l1946:
						position, tokenIndex = position1946, tokenIndex1946
					}
				l1947:
					add(rulePegText, position1937)
				}
				if buffer[position] != rune(']') {
					goto l1935
				}
				position++
				add(rulejsonArraySlice, position1936)
			}
			return true
		l1935:
			position, tokenIndex = position1935, tokenIndex1935
			return false
		},
		/* 179 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1952, tokenIndex1952 := position, tokenIndex
			{
				position1953 := position
				if buffer[position] != rune('[') {
					goto l1952
				}
				position++
				{
					position1954 := position
					{
						position1955, tokenIndex1955 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l1956
						}
						position++
						{
							position1957, tokenIndex1957 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1957
							}
							position++
							goto l1958
						l1957:
							position, tokenIndex = position1957, tokenIndex1957
						}
					l1958:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1956
						}
						position++
					l1959:
						{
							position1960, tokenIndex1960 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1960
							}
							position++
							goto l1959
						l1960:
							position, tokenIndex = position1960, tokenIndex1960
						}
						goto l1955
					l1956:
						position, tokenIndex = position1955, tokenIndex1955
						{
							position1961, tokenIndex1961 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1961
							}
							position++
							goto l1962
						l1961:
							position, tokenIndex = position1961, tokenIndex1961
						}
					l1962:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1952
						}
						position++
					l1963:
						{
							position1964, tokenIndex1964 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1964
							}
							position++
							goto l1963
						l1964:
							position, tokenIndex = position1964, tokenIndex1964
						}
						if buffer[position] != rune(':') {
							goto l1952
						}
						position++
					}
				l1955:
					add(rulePegText, position1954)
				}
				if buffer[position] != rune(']') {
					goto l1952
				}
				position++
				add(rulejsonArrayPartialSlice, position1953)
			}
			return true
		l1952:
			position, tokenIndex = position1952, tokenIndex1952
			return false
		},
		/* 180 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1965, tokenIndex1965 := position, tokenIndex
			{
				position1966 := position
				if buffer[position] != rune('[') {
					goto l1965
				}
				position++
				if buffer[position] != rune(':') {
					goto l1965
				}
				position++
				if buffer[position] != rune(']') {
					goto l1965
				}
				position++
				add(rulejsonArrayFullSlice, position1966)
			}
			return true
		l1965:
			position, tokenIndex = position1965, tokenIndex1965
			return false
		},
		/* 181 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1967, tokenIndex1967 := position, tokenIndex
			{
				position1968 := position
				{
					position1969, tokenIndex1969 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l1970
					}
					position++
					goto l1969
				l1970:
					position, tokenIndex = position1969, tokenIndex1969
					if buffer[position] != rune('\t') {
						goto l1971
					}
					position++
					goto l1969
				l1971:
					position, tokenIndex = position1969, tokenIndex1969
					if buffer[position] != rune('\n') {
						goto l1972
					}
					position++
					goto l1969
				l1972:
					position, tokenIndex = position1969, tokenIndex1969
					if buffer[position] != rune('\r') {
						goto l1973
					}
					position++
					goto l1969
				l1973:
					position, tokenIndex = position1969, tokenIndex1969
					if !_rules[rulecomment]() {
						goto l1974
					}
					goto l1969
				l1974:
					position, tokenIndex = position1969, tokenIndex1969
					if !_rules[rulefinalComment]() {
						goto l1967
					}
				}
			l1969:
				add(rulespElem, position1968)
			}
			return true
		l1967:
			position, tokenIndex = position1967, tokenIndex1967
			return false
		},
		/* 182 sp <- <spElem+> */
		func() bool {
			position1975, tokenIndex1975 := position, tokenIndex
			{
				position1976 := position
				if !_rules[rulespElem]() {
					goto l1975
				}
			l1977:
				{
					position1978, tokenIndex1978 := position, tokenIndex
					if !_rules[rulespElem]() {
						goto l1978
					}
					goto l1977
				l1978:
					position, tokenIndex = position1978, tokenIndex1978
				}
				add(rulesp, position1976)
			}
			return true
		l1975:
			position, tokenIndex = position1975, tokenIndex1975
			return false
		},
		/* 183 spOpt <- <spElem*> */
		func() bool {
			{
				position1980 := position
			l1981:
				{
					position1982, tokenIndex1982 := position, tokenIndex
					if !_rules[rulespElem]() {
						goto l1982
					}
					goto l1981
				l1982:
					position, tokenIndex = position1982, tokenIndex1982
				}
				add(rulespOpt, position1980)
			}
			return true
		},
		/* 184 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1983, tokenIndex1983 := position, tokenIndex
			{
				position1984 := position
				if buffer[position] != rune('-') {
					goto l1983
				}
				position++
				if buffer[position] != rune('-') {
					goto l1983
				}
				position++
			l1985:
				{
					position1986, tokenIndex1986 := position, tokenIndex
					{
						position1987, tokenIndex1987 := position, tokenIndex
						{
							position1988, tokenIndex1988 := position, tokenIndex
							if buffer[position] != rune('\r') {
								goto l1989
							}
							position++
							goto l1988
						l1989:
							position, tokenIndex = position1988, tokenIndex1988
							if buffer[position] != rune('\n') {
								goto l1987
							}
							position++
						}
					l1988:
						goto l1986
					l1987:
						position, tokenIndex = position1987, tokenIndex1987
					}
					if !matchDot() {
						goto l1986
					}
					goto l1985
				l1986:
					position, tokenIndex = position1986, tokenIndex1986
				}
				{
					position1990, tokenIndex1990 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l1991
					}
					position++
					goto l1990
				l1991:
					position, tokenIndex = position1990, tokenIndex1990
					if buffer[position] != rune('\n') {
						goto l1983
					}
					position++
				}
			l1990:
				add(rulecomment, position1984)
			}
			return true
		l1983:
			position, tokenIndex = position1983, tokenIndex1983
			return false
		},
		/* 185 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1992, tokenIndex1992 := position, tokenIndex
			{
				position1993 := position
				if buffer[position] != rune('-') {
					goto l1992
				}
				position++
				if buffer[position] != rune('-') {
					goto l1992
				}
				position++
			l1994:
				{
					position1995, tokenIndex1995 := position, tokenIndex
					{
						position1996, tokenIndex1996 := position, tokenIndex
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
								goto l1996
							}
							position++
						}
					l1997:
						goto l1995
					l1996:
						position, tokenIndex = position1996, tokenIndex1996
					}
					if !matchDot() {
						goto l1995
					}
					goto l1994
				l1995:
					position, tokenIndex = position1995, tokenIndex1995
				}
				{
					position1999, tokenIndex1999 := position, tokenIndex
					if !matchDot() {
						goto l1999
					}
					goto l1992
				l1999:
					position, tokenIndex = position1999, tokenIndex1999
				}
				add(rulefinalComment, position1993)
			}
			return true
		l1992:
			position, tokenIndex = position1992, tokenIndex1992
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
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 201 Action13 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 202 Action14 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 203 Action15 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 204 Action16 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 205 Action17 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 206 Action18 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 207 Action19 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 208 Action20 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 209 Action21 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 210 Action22 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 211 Action23 <- <{
		    p.AssembleEval(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 212 Action24 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 213 Action25 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 214 Action26 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 215 Action27 <- <{
		    p.AssembleEmitterSampling(CountBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 216 Action28 <- <{
		    p.AssembleEmitterSampling(RandomizedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 217 Action29 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 218 Action30 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 0.001)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 219 Action31 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 220 Action32 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 221 Action33 <- <{
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
		/* 222 Action34 <- <{
		    p.AssembleInterval()
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
		/* 225 Action37 <- <{
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
		/* 226 Action38 <- <{
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
		/* 227 Action39 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 228 Action40 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 229 Action41 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 230 Action42 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 231 Action43 <- <{
		    p.EnsureCapacitySpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 232 Action44 <- <{
		    p.EnsureSheddingSpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 233 Action45 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
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
		    p.EnsureIdentifier(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 237 Action49 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 238 Action50 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 239 Action51 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 240 Action52 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 241 Action53 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 242 Action54 <- <{
		    p.AssembleBinaryOperation(begin, end)
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
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 245 Action57 <- <{
		    p.AssembleBinaryOperation(begin, end)
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
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 251 Action63 <- <{
		    p.AssembleTypeCast(begin, end)
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
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 254 Action66 <- <{
		    p.AssembleExpressions(begin, end)
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
		    p.AssembleSortedExpression()
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 258 Action70 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 259 Action71 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 260 Action72 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 261 Action73 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 262 Action74 <- <{
		    p.AssembleConditionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 263 Action75 <- <{
		    p.AssembleExpressionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 264 Action76 <- <{
		    p.AssembleWhenThenPair()
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 265 Action77 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 266 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 267 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 268 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
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
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 271 Action83 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 272 Action84 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 273 Action85 <- <{
		    p.PushComponent(begin, end, NewMissing())
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

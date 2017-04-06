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
		/* 83 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / Case / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
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
					if !_rules[ruleFuncApp]() {
						goto l1123
					}
					goto l1115
				l1123:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleRowValue]() {
						goto l1124
					}
					goto l1115
				l1124:
					position, tokenIndex = position1115, tokenIndex1115
					if !_rules[ruleArrayExpr]() {
						goto l1125
					}
					goto l1115
				l1125:
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
			position1126, tokenIndex1126 := position, tokenIndex
			{
				position1127 := position
				{
					position1128 := position
					{
						position1129, tokenIndex1129 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1130
						}
						position++
						goto l1129
					l1130:
						position, tokenIndex = position1129, tokenIndex1129
						if buffer[position] != rune('C') {
							goto l1126
						}
						position++
					}
				l1129:
					{
						position1131, tokenIndex1131 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1132
						}
						position++
						goto l1131
					l1132:
						position, tokenIndex = position1131, tokenIndex1131
						if buffer[position] != rune('A') {
							goto l1126
						}
						position++
					}
				l1131:
					{
						position1133, tokenIndex1133 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1134
						}
						position++
						goto l1133
					l1134:
						position, tokenIndex = position1133, tokenIndex1133
						if buffer[position] != rune('S') {
							goto l1126
						}
						position++
					}
				l1133:
					{
						position1135, tokenIndex1135 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex = position1135, tokenIndex1135
						if buffer[position] != rune('T') {
							goto l1126
						}
						position++
					}
				l1135:
					if !_rules[rulespOpt]() {
						goto l1126
					}
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
					if !_rules[rulesp]() {
						goto l1126
					}
					{
						position1137, tokenIndex1137 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex = position1137, tokenIndex1137
						if buffer[position] != rune('A') {
							goto l1126
						}
						position++
					}
				l1137:
					{
						position1139, tokenIndex1139 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1140
						}
						position++
						goto l1139
					l1140:
						position, tokenIndex = position1139, tokenIndex1139
						if buffer[position] != rune('S') {
							goto l1126
						}
						position++
					}
				l1139:
					if !_rules[rulesp]() {
						goto l1126
					}
					if !_rules[ruleType]() {
						goto l1126
					}
					if !_rules[rulespOpt]() {
						goto l1126
					}
					if buffer[position] != rune(')') {
						goto l1126
					}
					position++
					add(rulePegText, position1128)
				}
				if !_rules[ruleAction64]() {
					goto l1126
				}
				add(ruleFuncTypeCast, position1127)
			}
			return true
		l1126:
			position, tokenIndex = position1126, tokenIndex1126
			return false
		},
		/* 85 FuncApp <- <(FuncAppWithOrderBy / FuncAppWithoutOrderBy)> */
		func() bool {
			position1141, tokenIndex1141 := position, tokenIndex
			{
				position1142 := position
				{
					position1143, tokenIndex1143 := position, tokenIndex
					if !_rules[ruleFuncAppWithOrderBy]() {
						goto l1144
					}
					goto l1143
				l1144:
					position, tokenIndex = position1143, tokenIndex1143
					if !_rules[ruleFuncAppWithoutOrderBy]() {
						goto l1141
					}
				}
			l1143:
				add(ruleFuncApp, position1142)
			}
			return true
		l1141:
			position, tokenIndex = position1141, tokenIndex1141
			return false
		},
		/* 86 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action65)> */
		func() bool {
			position1145, tokenIndex1145 := position, tokenIndex
			{
				position1146 := position
				if !_rules[ruleFunction]() {
					goto l1145
				}
				if !_rules[rulespOpt]() {
					goto l1145
				}
				if buffer[position] != rune('(') {
					goto l1145
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1145
				}
				if !_rules[ruleFuncParams]() {
					goto l1145
				}
				if !_rules[rulesp]() {
					goto l1145
				}
				if !_rules[ruleParamsOrder]() {
					goto l1145
				}
				if !_rules[rulespOpt]() {
					goto l1145
				}
				if buffer[position] != rune(')') {
					goto l1145
				}
				position++
				if !_rules[ruleAction65]() {
					goto l1145
				}
				add(ruleFuncAppWithOrderBy, position1146)
			}
			return true
		l1145:
			position, tokenIndex = position1145, tokenIndex1145
			return false
		},
		/* 87 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action66)> */
		func() bool {
			position1147, tokenIndex1147 := position, tokenIndex
			{
				position1148 := position
				if !_rules[ruleFunction]() {
					goto l1147
				}
				if !_rules[rulespOpt]() {
					goto l1147
				}
				if buffer[position] != rune('(') {
					goto l1147
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1147
				}
				if !_rules[ruleFuncParams]() {
					goto l1147
				}
				{
					position1149 := position
					if !_rules[rulespOpt]() {
						goto l1147
					}
					add(rulePegText, position1149)
				}
				if buffer[position] != rune(')') {
					goto l1147
				}
				position++
				if !_rules[ruleAction66]() {
					goto l1147
				}
				add(ruleFuncAppWithoutOrderBy, position1148)
			}
			return true
		l1147:
			position, tokenIndex = position1147, tokenIndex1147
			return false
		},
		/* 88 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action67)> */
		func() bool {
			position1150, tokenIndex1150 := position, tokenIndex
			{
				position1151 := position
				{
					position1152 := position
					{
						position1153, tokenIndex1153 := position, tokenIndex
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1153
						}
					l1155:
						{
							position1156, tokenIndex1156 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1156
							}
							if buffer[position] != rune(',') {
								goto l1156
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1156
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1156
							}
							goto l1155
						l1156:
							position, tokenIndex = position1156, tokenIndex1156
						}
						goto l1154
					l1153:
						position, tokenIndex = position1153, tokenIndex1153
					}
				l1154:
					add(rulePegText, position1152)
				}
				if !_rules[ruleAction67]() {
					goto l1150
				}
				add(ruleFuncParams, position1151)
			}
			return true
		l1150:
			position, tokenIndex = position1150, tokenIndex1150
			return false
		},
		/* 89 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action68)> */
		func() bool {
			position1157, tokenIndex1157 := position, tokenIndex
			{
				position1158 := position
				{
					position1159 := position
					{
						position1160, tokenIndex1160 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1161
						}
						position++
						goto l1160
					l1161:
						position, tokenIndex = position1160, tokenIndex1160
						if buffer[position] != rune('O') {
							goto l1157
						}
						position++
					}
				l1160:
					{
						position1162, tokenIndex1162 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1163
						}
						position++
						goto l1162
					l1163:
						position, tokenIndex = position1162, tokenIndex1162
						if buffer[position] != rune('R') {
							goto l1157
						}
						position++
					}
				l1162:
					{
						position1164, tokenIndex1164 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1165
						}
						position++
						goto l1164
					l1165:
						position, tokenIndex = position1164, tokenIndex1164
						if buffer[position] != rune('D') {
							goto l1157
						}
						position++
					}
				l1164:
					{
						position1166, tokenIndex1166 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1167
						}
						position++
						goto l1166
					l1167:
						position, tokenIndex = position1166, tokenIndex1166
						if buffer[position] != rune('E') {
							goto l1157
						}
						position++
					}
				l1166:
					{
						position1168, tokenIndex1168 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1169
						}
						position++
						goto l1168
					l1169:
						position, tokenIndex = position1168, tokenIndex1168
						if buffer[position] != rune('R') {
							goto l1157
						}
						position++
					}
				l1168:
					if !_rules[rulesp]() {
						goto l1157
					}
					{
						position1170, tokenIndex1170 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1171
						}
						position++
						goto l1170
					l1171:
						position, tokenIndex = position1170, tokenIndex1170
						if buffer[position] != rune('B') {
							goto l1157
						}
						position++
					}
				l1170:
					{
						position1172, tokenIndex1172 := position, tokenIndex
						if buffer[position] != rune('y') {
							goto l1173
						}
						position++
						goto l1172
					l1173:
						position, tokenIndex = position1172, tokenIndex1172
						if buffer[position] != rune('Y') {
							goto l1157
						}
						position++
					}
				l1172:
					if !_rules[rulesp]() {
						goto l1157
					}
					if !_rules[ruleSortedExpression]() {
						goto l1157
					}
				l1174:
					{
						position1175, tokenIndex1175 := position, tokenIndex
						if !_rules[rulespOpt]() {
							goto l1175
						}
						if buffer[position] != rune(',') {
							goto l1175
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1175
						}
						if !_rules[ruleSortedExpression]() {
							goto l1175
						}
						goto l1174
					l1175:
						position, tokenIndex = position1175, tokenIndex1175
					}
					add(rulePegText, position1159)
				}
				if !_rules[ruleAction68]() {
					goto l1157
				}
				add(ruleParamsOrder, position1158)
			}
			return true
		l1157:
			position, tokenIndex = position1157, tokenIndex1157
			return false
		},
		/* 90 SortedExpression <- <(Expression OrderDirectionOpt Action69)> */
		func() bool {
			position1176, tokenIndex1176 := position, tokenIndex
			{
				position1177 := position
				if !_rules[ruleExpression]() {
					goto l1176
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1176
				}
				if !_rules[ruleAction69]() {
					goto l1176
				}
				add(ruleSortedExpression, position1177)
			}
			return true
		l1176:
			position, tokenIndex = position1176, tokenIndex1176
			return false
		},
		/* 91 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action70)> */
		func() bool {
			position1178, tokenIndex1178 := position, tokenIndex
			{
				position1179 := position
				{
					position1180 := position
					{
						position1181, tokenIndex1181 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1181
						}
						{
							position1183, tokenIndex1183 := position, tokenIndex
							if !_rules[ruleAscending]() {
								goto l1184
							}
							goto l1183
						l1184:
							position, tokenIndex = position1183, tokenIndex1183
							if !_rules[ruleDescending]() {
								goto l1181
							}
						}
					l1183:
						goto l1182
					l1181:
						position, tokenIndex = position1181, tokenIndex1181
					}
				l1182:
					add(rulePegText, position1180)
				}
				if !_rules[ruleAction70]() {
					goto l1178
				}
				add(ruleOrderDirectionOpt, position1179)
			}
			return true
		l1178:
			position, tokenIndex = position1178, tokenIndex1178
			return false
		},
		/* 92 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action71)> */
		func() bool {
			position1185, tokenIndex1185 := position, tokenIndex
			{
				position1186 := position
				{
					position1187 := position
					if buffer[position] != rune('[') {
						goto l1185
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1185
					}
					{
						position1188, tokenIndex1188 := position, tokenIndex
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1188
						}
					l1190:
						{
							position1191, tokenIndex1191 := position, tokenIndex
							if !_rules[rulespOpt]() {
								goto l1191
							}
							if buffer[position] != rune(',') {
								goto l1191
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1191
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1191
							}
							goto l1190
						l1191:
							position, tokenIndex = position1191, tokenIndex1191
						}
						goto l1189
					l1188:
						position, tokenIndex = position1188, tokenIndex1188
					}
				l1189:
					if !_rules[rulespOpt]() {
						goto l1185
					}
					{
						position1192, tokenIndex1192 := position, tokenIndex
						if buffer[position] != rune(',') {
							goto l1192
						}
						position++
						goto l1193
					l1192:
						position, tokenIndex = position1192, tokenIndex1192
					}
				l1193:
					if !_rules[rulespOpt]() {
						goto l1185
					}
					if buffer[position] != rune(']') {
						goto l1185
					}
					position++
					add(rulePegText, position1187)
				}
				if !_rules[ruleAction71]() {
					goto l1185
				}
				add(ruleArrayExpr, position1186)
			}
			return true
		l1185:
			position, tokenIndex = position1185, tokenIndex1185
			return false
		},
		/* 93 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action72)> */
		func() bool {
			position1194, tokenIndex1194 := position, tokenIndex
			{
				position1195 := position
				{
					position1196 := position
					if buffer[position] != rune('{') {
						goto l1194
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1194
					}
					{
						position1197, tokenIndex1197 := position, tokenIndex
						if !_rules[ruleKeyValuePair]() {
							goto l1197
						}
					l1199:
						{
							position1200, tokenIndex1200 := position, tokenIndex
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
							if !_rules[ruleKeyValuePair]() {
								goto l1200
							}
							goto l1199
						l1200:
							position, tokenIndex = position1200, tokenIndex1200
						}
						goto l1198
					l1197:
						position, tokenIndex = position1197, tokenIndex1197
					}
				l1198:
					if !_rules[rulespOpt]() {
						goto l1194
					}
					if buffer[position] != rune('}') {
						goto l1194
					}
					position++
					add(rulePegText, position1196)
				}
				if !_rules[ruleAction72]() {
					goto l1194
				}
				add(ruleMapExpr, position1195)
			}
			return true
		l1194:
			position, tokenIndex = position1194, tokenIndex1194
			return false
		},
		/* 94 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action73)> */
		func() bool {
			position1201, tokenIndex1201 := position, tokenIndex
			{
				position1202 := position
				{
					position1203 := position
					if !_rules[ruleStringLiteral]() {
						goto l1201
					}
					if !_rules[rulespOpt]() {
						goto l1201
					}
					if buffer[position] != rune(':') {
						goto l1201
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1201
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1201
					}
					add(rulePegText, position1203)
				}
				if !_rules[ruleAction73]() {
					goto l1201
				}
				add(ruleKeyValuePair, position1202)
			}
			return true
		l1201:
			position, tokenIndex = position1201, tokenIndex1201
			return false
		},
		/* 95 Case <- <(ConditionCase / ExpressionCase)> */
		func() bool {
			position1204, tokenIndex1204 := position, tokenIndex
			{
				position1205 := position
				{
					position1206, tokenIndex1206 := position, tokenIndex
					if !_rules[ruleConditionCase]() {
						goto l1207
					}
					goto l1206
				l1207:
					position, tokenIndex = position1206, tokenIndex1206
					if !_rules[ruleExpressionCase]() {
						goto l1204
					}
				}
			l1206:
				add(ruleCase, position1205)
			}
			return true
		l1204:
			position, tokenIndex = position1204, tokenIndex1204
			return false
		},
		/* 96 ConditionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action74)> */
		func() bool {
			position1208, tokenIndex1208 := position, tokenIndex
			{
				position1209 := position
				{
					position1210, tokenIndex1210 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l1211
					}
					position++
					goto l1210
				l1211:
					position, tokenIndex = position1210, tokenIndex1210
					if buffer[position] != rune('C') {
						goto l1208
					}
					position++
				}
			l1210:
				{
					position1212, tokenIndex1212 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l1213
					}
					position++
					goto l1212
				l1213:
					position, tokenIndex = position1212, tokenIndex1212
					if buffer[position] != rune('A') {
						goto l1208
					}
					position++
				}
			l1212:
				{
					position1214, tokenIndex1214 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l1215
					}
					position++
					goto l1214
				l1215:
					position, tokenIndex = position1214, tokenIndex1214
					if buffer[position] != rune('S') {
						goto l1208
					}
					position++
				}
			l1214:
				{
					position1216, tokenIndex1216 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1217
					}
					position++
					goto l1216
				l1217:
					position, tokenIndex = position1216, tokenIndex1216
					if buffer[position] != rune('E') {
						goto l1208
					}
					position++
				}
			l1216:
				{
					position1218 := position
					if !_rules[rulesp]() {
						goto l1208
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1208
					}
				l1219:
					{
						position1220, tokenIndex1220 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1220
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1220
						}
						goto l1219
					l1220:
						position, tokenIndex = position1220, tokenIndex1220
					}
					{
						position1221, tokenIndex1221 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1221
						}
						{
							position1223, tokenIndex1223 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1224
							}
							position++
							goto l1223
						l1224:
							position, tokenIndex = position1223, tokenIndex1223
							if buffer[position] != rune('E') {
								goto l1221
							}
							position++
						}
					l1223:
						{
							position1225, tokenIndex1225 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l1226
							}
							position++
							goto l1225
						l1226:
							position, tokenIndex = position1225, tokenIndex1225
							if buffer[position] != rune('L') {
								goto l1221
							}
							position++
						}
					l1225:
						{
							position1227, tokenIndex1227 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l1228
							}
							position++
							goto l1227
						l1228:
							position, tokenIndex = position1227, tokenIndex1227
							if buffer[position] != rune('S') {
								goto l1221
							}
							position++
						}
					l1227:
						{
							position1229, tokenIndex1229 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1230
							}
							position++
							goto l1229
						l1230:
							position, tokenIndex = position1229, tokenIndex1229
							if buffer[position] != rune('E') {
								goto l1221
							}
							position++
						}
					l1229:
						if !_rules[rulesp]() {
							goto l1221
						}
						if !_rules[ruleExpression]() {
							goto l1221
						}
						goto l1222
					l1221:
						position, tokenIndex = position1221, tokenIndex1221
					}
				l1222:
					if !_rules[rulesp]() {
						goto l1208
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
							goto l1208
						}
						position++
					}
				l1231:
					{
						position1233, tokenIndex1233 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1234
						}
						position++
						goto l1233
					l1234:
						position, tokenIndex = position1233, tokenIndex1233
						if buffer[position] != rune('N') {
							goto l1208
						}
						position++
					}
				l1233:
					{
						position1235, tokenIndex1235 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1236
						}
						position++
						goto l1235
					l1236:
						position, tokenIndex = position1235, tokenIndex1235
						if buffer[position] != rune('D') {
							goto l1208
						}
						position++
					}
				l1235:
					add(rulePegText, position1218)
				}
				if !_rules[ruleAction74]() {
					goto l1208
				}
				add(ruleConditionCase, position1209)
			}
			return true
		l1208:
			position, tokenIndex = position1208, tokenIndex1208
			return false
		},
		/* 97 ExpressionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') sp Expression <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action75)> */
		func() bool {
			position1237, tokenIndex1237 := position, tokenIndex
			{
				position1238 := position
				{
					position1239, tokenIndex1239 := position, tokenIndex
					if buffer[position] != rune('c') {
						goto l1240
					}
					position++
					goto l1239
				l1240:
					position, tokenIndex = position1239, tokenIndex1239
					if buffer[position] != rune('C') {
						goto l1237
					}
					position++
				}
			l1239:
				{
					position1241, tokenIndex1241 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l1242
					}
					position++
					goto l1241
				l1242:
					position, tokenIndex = position1241, tokenIndex1241
					if buffer[position] != rune('A') {
						goto l1237
					}
					position++
				}
			l1241:
				{
					position1243, tokenIndex1243 := position, tokenIndex
					if buffer[position] != rune('s') {
						goto l1244
					}
					position++
					goto l1243
				l1244:
					position, tokenIndex = position1243, tokenIndex1243
					if buffer[position] != rune('S') {
						goto l1237
					}
					position++
				}
			l1243:
				{
					position1245, tokenIndex1245 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1246
					}
					position++
					goto l1245
				l1246:
					position, tokenIndex = position1245, tokenIndex1245
					if buffer[position] != rune('E') {
						goto l1237
					}
					position++
				}
			l1245:
				if !_rules[rulesp]() {
					goto l1237
				}
				if !_rules[ruleExpression]() {
					goto l1237
				}
				{
					position1247 := position
					if !_rules[rulesp]() {
						goto l1237
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1237
					}
				l1248:
					{
						position1249, tokenIndex1249 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1249
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1249
						}
						goto l1248
					l1249:
						position, tokenIndex = position1249, tokenIndex1249
					}
					{
						position1250, tokenIndex1250 := position, tokenIndex
						if !_rules[rulesp]() {
							goto l1250
						}
						{
							position1252, tokenIndex1252 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1253
							}
							position++
							goto l1252
						l1253:
							position, tokenIndex = position1252, tokenIndex1252
							if buffer[position] != rune('E') {
								goto l1250
							}
							position++
						}
					l1252:
						{
							position1254, tokenIndex1254 := position, tokenIndex
							if buffer[position] != rune('l') {
								goto l1255
							}
							position++
							goto l1254
						l1255:
							position, tokenIndex = position1254, tokenIndex1254
							if buffer[position] != rune('L') {
								goto l1250
							}
							position++
						}
					l1254:
						{
							position1256, tokenIndex1256 := position, tokenIndex
							if buffer[position] != rune('s') {
								goto l1257
							}
							position++
							goto l1256
						l1257:
							position, tokenIndex = position1256, tokenIndex1256
							if buffer[position] != rune('S') {
								goto l1250
							}
							position++
						}
					l1256:
						{
							position1258, tokenIndex1258 := position, tokenIndex
							if buffer[position] != rune('e') {
								goto l1259
							}
							position++
							goto l1258
						l1259:
							position, tokenIndex = position1258, tokenIndex1258
							if buffer[position] != rune('E') {
								goto l1250
							}
							position++
						}
					l1258:
						if !_rules[rulesp]() {
							goto l1250
						}
						if !_rules[ruleExpression]() {
							goto l1250
						}
						goto l1251
					l1250:
						position, tokenIndex = position1250, tokenIndex1250
					}
				l1251:
					if !_rules[rulesp]() {
						goto l1237
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
							goto l1237
						}
						position++
					}
				l1260:
					{
						position1262, tokenIndex1262 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1263
						}
						position++
						goto l1262
					l1263:
						position, tokenIndex = position1262, tokenIndex1262
						if buffer[position] != rune('N') {
							goto l1237
						}
						position++
					}
				l1262:
					{
						position1264, tokenIndex1264 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1265
						}
						position++
						goto l1264
					l1265:
						position, tokenIndex = position1264, tokenIndex1264
						if buffer[position] != rune('D') {
							goto l1237
						}
						position++
					}
				l1264:
					add(rulePegText, position1247)
				}
				if !_rules[ruleAction75]() {
					goto l1237
				}
				add(ruleExpressionCase, position1238)
			}
			return true
		l1237:
			position, tokenIndex = position1237, tokenIndex1237
			return false
		},
		/* 98 WhenThenPair <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('n' / 'N') sp Expression sp (('t' / 'T') ('h' / 'H') ('e' / 'E') ('n' / 'N')) sp ExpressionOrWildcard Action76)> */
		func() bool {
			position1266, tokenIndex1266 := position, tokenIndex
			{
				position1267 := position
				{
					position1268, tokenIndex1268 := position, tokenIndex
					if buffer[position] != rune('w') {
						goto l1269
					}
					position++
					goto l1268
				l1269:
					position, tokenIndex = position1268, tokenIndex1268
					if buffer[position] != rune('W') {
						goto l1266
					}
					position++
				}
			l1268:
				{
					position1270, tokenIndex1270 := position, tokenIndex
					if buffer[position] != rune('h') {
						goto l1271
					}
					position++
					goto l1270
				l1271:
					position, tokenIndex = position1270, tokenIndex1270
					if buffer[position] != rune('H') {
						goto l1266
					}
					position++
				}
			l1270:
				{
					position1272, tokenIndex1272 := position, tokenIndex
					if buffer[position] != rune('e') {
						goto l1273
					}
					position++
					goto l1272
				l1273:
					position, tokenIndex = position1272, tokenIndex1272
					if buffer[position] != rune('E') {
						goto l1266
					}
					position++
				}
			l1272:
				{
					position1274, tokenIndex1274 := position, tokenIndex
					if buffer[position] != rune('n') {
						goto l1275
					}
					position++
					goto l1274
				l1275:
					position, tokenIndex = position1274, tokenIndex1274
					if buffer[position] != rune('N') {
						goto l1266
					}
					position++
				}
			l1274:
				if !_rules[rulesp]() {
					goto l1266
				}
				if !_rules[ruleExpression]() {
					goto l1266
				}
				if !_rules[rulesp]() {
					goto l1266
				}
				{
					position1276, tokenIndex1276 := position, tokenIndex
					if buffer[position] != rune('t') {
						goto l1277
					}
					position++
					goto l1276
				l1277:
					position, tokenIndex = position1276, tokenIndex1276
					if buffer[position] != rune('T') {
						goto l1266
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
						goto l1266
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
						goto l1266
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
						goto l1266
					}
					position++
				}
			l1282:
				if !_rules[rulesp]() {
					goto l1266
				}
				if !_rules[ruleExpressionOrWildcard]() {
					goto l1266
				}
				if !_rules[ruleAction76]() {
					goto l1266
				}
				add(ruleWhenThenPair, position1267)
			}
			return true
		l1266:
			position, tokenIndex = position1266, tokenIndex1266
			return false
		},
		/* 99 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1284, tokenIndex1284 := position, tokenIndex
			{
				position1285 := position
				{
					position1286, tokenIndex1286 := position, tokenIndex
					if !_rules[ruleFloatLiteral]() {
						goto l1287
					}
					goto l1286
				l1287:
					position, tokenIndex = position1286, tokenIndex1286
					if !_rules[ruleNumericLiteral]() {
						goto l1288
					}
					goto l1286
				l1288:
					position, tokenIndex = position1286, tokenIndex1286
					if !_rules[ruleStringLiteral]() {
						goto l1284
					}
				}
			l1286:
				add(ruleLiteral, position1285)
			}
			return true
		l1284:
			position, tokenIndex = position1284, tokenIndex1284
			return false
		},
		/* 100 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1289, tokenIndex1289 := position, tokenIndex
			{
				position1290 := position
				{
					position1291, tokenIndex1291 := position, tokenIndex
					if !_rules[ruleEqual]() {
						goto l1292
					}
					goto l1291
				l1292:
					position, tokenIndex = position1291, tokenIndex1291
					if !_rules[ruleNotEqual]() {
						goto l1293
					}
					goto l1291
				l1293:
					position, tokenIndex = position1291, tokenIndex1291
					if !_rules[ruleLessOrEqual]() {
						goto l1294
					}
					goto l1291
				l1294:
					position, tokenIndex = position1291, tokenIndex1291
					if !_rules[ruleLess]() {
						goto l1295
					}
					goto l1291
				l1295:
					position, tokenIndex = position1291, tokenIndex1291
					if !_rules[ruleGreaterOrEqual]() {
						goto l1296
					}
					goto l1291
				l1296:
					position, tokenIndex = position1291, tokenIndex1291
					if !_rules[ruleGreater]() {
						goto l1297
					}
					goto l1291
				l1297:
					position, tokenIndex = position1291, tokenIndex1291
					if !_rules[ruleNotEqual]() {
						goto l1289
					}
				}
			l1291:
				add(ruleComparisonOp, position1290)
			}
			return true
		l1289:
			position, tokenIndex = position1289, tokenIndex1289
			return false
		},
		/* 101 OtherOp <- <Concat> */
		func() bool {
			position1298, tokenIndex1298 := position, tokenIndex
			{
				position1299 := position
				if !_rules[ruleConcat]() {
					goto l1298
				}
				add(ruleOtherOp, position1299)
			}
			return true
		l1298:
			position, tokenIndex = position1298, tokenIndex1298
			return false
		},
		/* 102 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1300, tokenIndex1300 := position, tokenIndex
			{
				position1301 := position
				{
					position1302, tokenIndex1302 := position, tokenIndex
					if !_rules[ruleIsNot]() {
						goto l1303
					}
					goto l1302
				l1303:
					position, tokenIndex = position1302, tokenIndex1302
					if !_rules[ruleIs]() {
						goto l1300
					}
				}
			l1302:
				add(ruleIsOp, position1301)
			}
			return true
		l1300:
			position, tokenIndex = position1300, tokenIndex1300
			return false
		},
		/* 103 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1304, tokenIndex1304 := position, tokenIndex
			{
				position1305 := position
				{
					position1306, tokenIndex1306 := position, tokenIndex
					if !_rules[rulePlus]() {
						goto l1307
					}
					goto l1306
				l1307:
					position, tokenIndex = position1306, tokenIndex1306
					if !_rules[ruleMinus]() {
						goto l1304
					}
				}
			l1306:
				add(rulePlusMinusOp, position1305)
			}
			return true
		l1304:
			position, tokenIndex = position1304, tokenIndex1304
			return false
		},
		/* 104 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1308, tokenIndex1308 := position, tokenIndex
			{
				position1309 := position
				{
					position1310, tokenIndex1310 := position, tokenIndex
					if !_rules[ruleMultiply]() {
						goto l1311
					}
					goto l1310
				l1311:
					position, tokenIndex = position1310, tokenIndex1310
					if !_rules[ruleDivide]() {
						goto l1312
					}
					goto l1310
				l1312:
					position, tokenIndex = position1310, tokenIndex1310
					if !_rules[ruleModulo]() {
						goto l1308
					}
				}
			l1310:
				add(ruleMultDivOp, position1309)
			}
			return true
		l1308:
			position, tokenIndex = position1308, tokenIndex1308
			return false
		},
		/* 105 Stream <- <(<ident> Action77)> */
		func() bool {
			position1313, tokenIndex1313 := position, tokenIndex
			{
				position1314 := position
				{
					position1315 := position
					if !_rules[ruleident]() {
						goto l1313
					}
					add(rulePegText, position1315)
				}
				if !_rules[ruleAction77]() {
					goto l1313
				}
				add(ruleStream, position1314)
			}
			return true
		l1313:
			position, tokenIndex = position1313, tokenIndex1313
			return false
		},
		/* 106 RowMeta <- <RowTimestamp> */
		func() bool {
			position1316, tokenIndex1316 := position, tokenIndex
			{
				position1317 := position
				if !_rules[ruleRowTimestamp]() {
					goto l1316
				}
				add(ruleRowMeta, position1317)
			}
			return true
		l1316:
			position, tokenIndex = position1316, tokenIndex1316
			return false
		},
		/* 107 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action78)> */
		func() bool {
			position1318, tokenIndex1318 := position, tokenIndex
			{
				position1319 := position
				{
					position1320 := position
					{
						position1321, tokenIndex1321 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1321
						}
						if buffer[position] != rune(':') {
							goto l1321
						}
						position++
						goto l1322
					l1321:
						position, tokenIndex = position1321, tokenIndex1321
					}
				l1322:
					if buffer[position] != rune('t') {
						goto l1318
					}
					position++
					if buffer[position] != rune('s') {
						goto l1318
					}
					position++
					if buffer[position] != rune('(') {
						goto l1318
					}
					position++
					if buffer[position] != rune(')') {
						goto l1318
					}
					position++
					add(rulePegText, position1320)
				}
				if !_rules[ruleAction78]() {
					goto l1318
				}
				add(ruleRowTimestamp, position1319)
			}
			return true
		l1318:
			position, tokenIndex = position1318, tokenIndex1318
			return false
		},
		/* 108 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action79)> */
		func() bool {
			position1323, tokenIndex1323 := position, tokenIndex
			{
				position1324 := position
				{
					position1325 := position
					{
						position1326, tokenIndex1326 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1326
						}
						if buffer[position] != rune(':') {
							goto l1326
						}
						position++
						{
							position1328, tokenIndex1328 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l1328
							}
							position++
							goto l1326
						l1328:
							position, tokenIndex = position1328, tokenIndex1328
						}
						goto l1327
					l1326:
						position, tokenIndex = position1326, tokenIndex1326
					}
				l1327:
					if !_rules[rulejsonGetPath]() {
						goto l1323
					}
					add(rulePegText, position1325)
				}
				if !_rules[ruleAction79]() {
					goto l1323
				}
				add(ruleRowValue, position1324)
			}
			return true
		l1323:
			position, tokenIndex = position1323, tokenIndex1323
			return false
		},
		/* 109 NumericLiteral <- <(<('-'? [0-9]+)> Action80)> */
		func() bool {
			position1329, tokenIndex1329 := position, tokenIndex
			{
				position1330 := position
				{
					position1331 := position
					{
						position1332, tokenIndex1332 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1332
						}
						position++
						goto l1333
					l1332:
						position, tokenIndex = position1332, tokenIndex1332
					}
				l1333:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1329
					}
					position++
				l1334:
					{
						position1335, tokenIndex1335 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1335
						}
						position++
						goto l1334
					l1335:
						position, tokenIndex = position1335, tokenIndex1335
					}
					add(rulePegText, position1331)
				}
				if !_rules[ruleAction80]() {
					goto l1329
				}
				add(ruleNumericLiteral, position1330)
			}
			return true
		l1329:
			position, tokenIndex = position1329, tokenIndex1329
			return false
		},
		/* 110 NonNegativeNumericLiteral <- <(<[0-9]+> Action81)> */
		func() bool {
			position1336, tokenIndex1336 := position, tokenIndex
			{
				position1337 := position
				{
					position1338 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1336
					}
					position++
				l1339:
					{
						position1340, tokenIndex1340 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1340
						}
						position++
						goto l1339
					l1340:
						position, tokenIndex = position1340, tokenIndex1340
					}
					add(rulePegText, position1338)
				}
				if !_rules[ruleAction81]() {
					goto l1336
				}
				add(ruleNonNegativeNumericLiteral, position1337)
			}
			return true
		l1336:
			position, tokenIndex = position1336, tokenIndex1336
			return false
		},
		/* 111 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action82)> */
		func() bool {
			position1341, tokenIndex1341 := position, tokenIndex
			{
				position1342 := position
				{
					position1343 := position
					{
						position1344, tokenIndex1344 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1344
						}
						position++
						goto l1345
					l1344:
						position, tokenIndex = position1344, tokenIndex1344
					}
				l1345:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1341
					}
					position++
				l1346:
					{
						position1347, tokenIndex1347 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1347
						}
						position++
						goto l1346
					l1347:
						position, tokenIndex = position1347, tokenIndex1347
					}
					if buffer[position] != rune('.') {
						goto l1341
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1341
					}
					position++
				l1348:
					{
						position1349, tokenIndex1349 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1349
						}
						position++
						goto l1348
					l1349:
						position, tokenIndex = position1349, tokenIndex1349
					}
					add(rulePegText, position1343)
				}
				if !_rules[ruleAction82]() {
					goto l1341
				}
				add(ruleFloatLiteral, position1342)
			}
			return true
		l1341:
			position, tokenIndex = position1341, tokenIndex1341
			return false
		},
		/* 112 Function <- <(<ident> Action83)> */
		func() bool {
			position1350, tokenIndex1350 := position, tokenIndex
			{
				position1351 := position
				{
					position1352 := position
					if !_rules[ruleident]() {
						goto l1350
					}
					add(rulePegText, position1352)
				}
				if !_rules[ruleAction83]() {
					goto l1350
				}
				add(ruleFunction, position1351)
			}
			return true
		l1350:
			position, tokenIndex = position1350, tokenIndex1350
			return false
		},
		/* 113 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action84)> */
		func() bool {
			position1353, tokenIndex1353 := position, tokenIndex
			{
				position1354 := position
				{
					position1355 := position
					{
						position1356, tokenIndex1356 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1357
						}
						position++
						goto l1356
					l1357:
						position, tokenIndex = position1356, tokenIndex1356
						if buffer[position] != rune('N') {
							goto l1353
						}
						position++
					}
				l1356:
					{
						position1358, tokenIndex1358 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1359
						}
						position++
						goto l1358
					l1359:
						position, tokenIndex = position1358, tokenIndex1358
						if buffer[position] != rune('U') {
							goto l1353
						}
						position++
					}
				l1358:
					{
						position1360, tokenIndex1360 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1361
						}
						position++
						goto l1360
					l1361:
						position, tokenIndex = position1360, tokenIndex1360
						if buffer[position] != rune('L') {
							goto l1353
						}
						position++
					}
				l1360:
					{
						position1362, tokenIndex1362 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1363
						}
						position++
						goto l1362
					l1363:
						position, tokenIndex = position1362, tokenIndex1362
						if buffer[position] != rune('L') {
							goto l1353
						}
						position++
					}
				l1362:
					add(rulePegText, position1355)
				}
				if !_rules[ruleAction84]() {
					goto l1353
				}
				add(ruleNullLiteral, position1354)
			}
			return true
		l1353:
			position, tokenIndex = position1353, tokenIndex1353
			return false
		},
		/* 114 Missing <- <(<(('m' / 'M') ('i' / 'I') ('s' / 'S') ('s' / 'S') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action85)> */
		func() bool {
			position1364, tokenIndex1364 := position, tokenIndex
			{
				position1365 := position
				{
					position1366 := position
					{
						position1367, tokenIndex1367 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1368
						}
						position++
						goto l1367
					l1368:
						position, tokenIndex = position1367, tokenIndex1367
						if buffer[position] != rune('M') {
							goto l1364
						}
						position++
					}
				l1367:
					{
						position1369, tokenIndex1369 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1370
						}
						position++
						goto l1369
					l1370:
						position, tokenIndex = position1369, tokenIndex1369
						if buffer[position] != rune('I') {
							goto l1364
						}
						position++
					}
				l1369:
					{
						position1371, tokenIndex1371 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1372
						}
						position++
						goto l1371
					l1372:
						position, tokenIndex = position1371, tokenIndex1371
						if buffer[position] != rune('S') {
							goto l1364
						}
						position++
					}
				l1371:
					{
						position1373, tokenIndex1373 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1374
						}
						position++
						goto l1373
					l1374:
						position, tokenIndex = position1373, tokenIndex1373
						if buffer[position] != rune('S') {
							goto l1364
						}
						position++
					}
				l1373:
					{
						position1375, tokenIndex1375 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1376
						}
						position++
						goto l1375
					l1376:
						position, tokenIndex = position1375, tokenIndex1375
						if buffer[position] != rune('I') {
							goto l1364
						}
						position++
					}
				l1375:
					{
						position1377, tokenIndex1377 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1378
						}
						position++
						goto l1377
					l1378:
						position, tokenIndex = position1377, tokenIndex1377
						if buffer[position] != rune('N') {
							goto l1364
						}
						position++
					}
				l1377:
					{
						position1379, tokenIndex1379 := position, tokenIndex
						if buffer[position] != rune('g') {
							goto l1380
						}
						position++
						goto l1379
					l1380:
						position, tokenIndex = position1379, tokenIndex1379
						if buffer[position] != rune('G') {
							goto l1364
						}
						position++
					}
				l1379:
					add(rulePegText, position1366)
				}
				if !_rules[ruleAction85]() {
					goto l1364
				}
				add(ruleMissing, position1365)
			}
			return true
		l1364:
			position, tokenIndex = position1364, tokenIndex1364
			return false
		},
		/* 115 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1381, tokenIndex1381 := position, tokenIndex
			{
				position1382 := position
				{
					position1383, tokenIndex1383 := position, tokenIndex
					if !_rules[ruleTRUE]() {
						goto l1384
					}
					goto l1383
				l1384:
					position, tokenIndex = position1383, tokenIndex1383
					if !_rules[ruleFALSE]() {
						goto l1381
					}
				}
			l1383:
				add(ruleBooleanLiteral, position1382)
			}
			return true
		l1381:
			position, tokenIndex = position1381, tokenIndex1381
			return false
		},
		/* 116 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action86)> */
		func() bool {
			position1385, tokenIndex1385 := position, tokenIndex
			{
				position1386 := position
				{
					position1387 := position
					{
						position1388, tokenIndex1388 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1389
						}
						position++
						goto l1388
					l1389:
						position, tokenIndex = position1388, tokenIndex1388
						if buffer[position] != rune('T') {
							goto l1385
						}
						position++
					}
				l1388:
					{
						position1390, tokenIndex1390 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1391
						}
						position++
						goto l1390
					l1391:
						position, tokenIndex = position1390, tokenIndex1390
						if buffer[position] != rune('R') {
							goto l1385
						}
						position++
					}
				l1390:
					{
						position1392, tokenIndex1392 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1393
						}
						position++
						goto l1392
					l1393:
						position, tokenIndex = position1392, tokenIndex1392
						if buffer[position] != rune('U') {
							goto l1385
						}
						position++
					}
				l1392:
					{
						position1394, tokenIndex1394 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1395
						}
						position++
						goto l1394
					l1395:
						position, tokenIndex = position1394, tokenIndex1394
						if buffer[position] != rune('E') {
							goto l1385
						}
						position++
					}
				l1394:
					add(rulePegText, position1387)
				}
				if !_rules[ruleAction86]() {
					goto l1385
				}
				add(ruleTRUE, position1386)
			}
			return true
		l1385:
			position, tokenIndex = position1385, tokenIndex1385
			return false
		},
		/* 117 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action87)> */
		func() bool {
			position1396, tokenIndex1396 := position, tokenIndex
			{
				position1397 := position
				{
					position1398 := position
					{
						position1399, tokenIndex1399 := position, tokenIndex
						if buffer[position] != rune('f') {
							goto l1400
						}
						position++
						goto l1399
					l1400:
						position, tokenIndex = position1399, tokenIndex1399
						if buffer[position] != rune('F') {
							goto l1396
						}
						position++
					}
				l1399:
					{
						position1401, tokenIndex1401 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1402
						}
						position++
						goto l1401
					l1402:
						position, tokenIndex = position1401, tokenIndex1401
						if buffer[position] != rune('A') {
							goto l1396
						}
						position++
					}
				l1401:
					{
						position1403, tokenIndex1403 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1404
						}
						position++
						goto l1403
					l1404:
						position, tokenIndex = position1403, tokenIndex1403
						if buffer[position] != rune('L') {
							goto l1396
						}
						position++
					}
				l1403:
					{
						position1405, tokenIndex1405 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1406
						}
						position++
						goto l1405
					l1406:
						position, tokenIndex = position1405, tokenIndex1405
						if buffer[position] != rune('S') {
							goto l1396
						}
						position++
					}
				l1405:
					{
						position1407, tokenIndex1407 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1408
						}
						position++
						goto l1407
					l1408:
						position, tokenIndex = position1407, tokenIndex1407
						if buffer[position] != rune('E') {
							goto l1396
						}
						position++
					}
				l1407:
					add(rulePegText, position1398)
				}
				if !_rules[ruleAction87]() {
					goto l1396
				}
				add(ruleFALSE, position1397)
			}
			return true
		l1396:
			position, tokenIndex = position1396, tokenIndex1396
			return false
		},
		/* 118 Wildcard <- <(<((ident ':' !':')? '*')> Action88)> */
		func() bool {
			position1409, tokenIndex1409 := position, tokenIndex
			{
				position1410 := position
				{
					position1411 := position
					{
						position1412, tokenIndex1412 := position, tokenIndex
						if !_rules[ruleident]() {
							goto l1412
						}
						if buffer[position] != rune(':') {
							goto l1412
						}
						position++
						{
							position1414, tokenIndex1414 := position, tokenIndex
							if buffer[position] != rune(':') {
								goto l1414
							}
							position++
							goto l1412
						l1414:
							position, tokenIndex = position1414, tokenIndex1414
						}
						goto l1413
					l1412:
						position, tokenIndex = position1412, tokenIndex1412
					}
				l1413:
					if buffer[position] != rune('*') {
						goto l1409
					}
					position++
					add(rulePegText, position1411)
				}
				if !_rules[ruleAction88]() {
					goto l1409
				}
				add(ruleWildcard, position1410)
			}
			return true
		l1409:
			position, tokenIndex = position1409, tokenIndex1409
			return false
		},
		/* 119 StringLiteral <- <(<('"' (('"' '"') / (!'"' .))* '"')> Action89)> */
		func() bool {
			position1415, tokenIndex1415 := position, tokenIndex
			{
				position1416 := position
				{
					position1417 := position
					if buffer[position] != rune('"') {
						goto l1415
					}
					position++
				l1418:
					{
						position1419, tokenIndex1419 := position, tokenIndex
						{
							position1420, tokenIndex1420 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l1421
							}
							position++
							if buffer[position] != rune('"') {
								goto l1421
							}
							position++
							goto l1420
						l1421:
							position, tokenIndex = position1420, tokenIndex1420
							{
								position1422, tokenIndex1422 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l1422
								}
								position++
								goto l1419
							l1422:
								position, tokenIndex = position1422, tokenIndex1422
							}
							if !matchDot() {
								goto l1419
							}
						}
					l1420:
						goto l1418
					l1419:
						position, tokenIndex = position1419, tokenIndex1419
					}
					if buffer[position] != rune('"') {
						goto l1415
					}
					position++
					add(rulePegText, position1417)
				}
				if !_rules[ruleAction89]() {
					goto l1415
				}
				add(ruleStringLiteral, position1416)
			}
			return true
		l1415:
			position, tokenIndex = position1415, tokenIndex1415
			return false
		},
		/* 120 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action90)> */
		func() bool {
			position1423, tokenIndex1423 := position, tokenIndex
			{
				position1424 := position
				{
					position1425 := position
					{
						position1426, tokenIndex1426 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1427
						}
						position++
						goto l1426
					l1427:
						position, tokenIndex = position1426, tokenIndex1426
						if buffer[position] != rune('I') {
							goto l1423
						}
						position++
					}
				l1426:
					{
						position1428, tokenIndex1428 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1429
						}
						position++
						goto l1428
					l1429:
						position, tokenIndex = position1428, tokenIndex1428
						if buffer[position] != rune('S') {
							goto l1423
						}
						position++
					}
				l1428:
					{
						position1430, tokenIndex1430 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1431
						}
						position++
						goto l1430
					l1431:
						position, tokenIndex = position1430, tokenIndex1430
						if buffer[position] != rune('T') {
							goto l1423
						}
						position++
					}
				l1430:
					{
						position1432, tokenIndex1432 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1433
						}
						position++
						goto l1432
					l1433:
						position, tokenIndex = position1432, tokenIndex1432
						if buffer[position] != rune('R') {
							goto l1423
						}
						position++
					}
				l1432:
					{
						position1434, tokenIndex1434 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1435
						}
						position++
						goto l1434
					l1435:
						position, tokenIndex = position1434, tokenIndex1434
						if buffer[position] != rune('E') {
							goto l1423
						}
						position++
					}
				l1434:
					{
						position1436, tokenIndex1436 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1437
						}
						position++
						goto l1436
					l1437:
						position, tokenIndex = position1436, tokenIndex1436
						if buffer[position] != rune('A') {
							goto l1423
						}
						position++
					}
				l1436:
					{
						position1438, tokenIndex1438 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1439
						}
						position++
						goto l1438
					l1439:
						position, tokenIndex = position1438, tokenIndex1438
						if buffer[position] != rune('M') {
							goto l1423
						}
						position++
					}
				l1438:
					add(rulePegText, position1425)
				}
				if !_rules[ruleAction90]() {
					goto l1423
				}
				add(ruleISTREAM, position1424)
			}
			return true
		l1423:
			position, tokenIndex = position1423, tokenIndex1423
			return false
		},
		/* 121 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action91)> */
		func() bool {
			position1440, tokenIndex1440 := position, tokenIndex
			{
				position1441 := position
				{
					position1442 := position
					{
						position1443, tokenIndex1443 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1444
						}
						position++
						goto l1443
					l1444:
						position, tokenIndex = position1443, tokenIndex1443
						if buffer[position] != rune('D') {
							goto l1440
						}
						position++
					}
				l1443:
					{
						position1445, tokenIndex1445 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1446
						}
						position++
						goto l1445
					l1446:
						position, tokenIndex = position1445, tokenIndex1445
						if buffer[position] != rune('S') {
							goto l1440
						}
						position++
					}
				l1445:
					{
						position1447, tokenIndex1447 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1448
						}
						position++
						goto l1447
					l1448:
						position, tokenIndex = position1447, tokenIndex1447
						if buffer[position] != rune('T') {
							goto l1440
						}
						position++
					}
				l1447:
					{
						position1449, tokenIndex1449 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1450
						}
						position++
						goto l1449
					l1450:
						position, tokenIndex = position1449, tokenIndex1449
						if buffer[position] != rune('R') {
							goto l1440
						}
						position++
					}
				l1449:
					{
						position1451, tokenIndex1451 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1452
						}
						position++
						goto l1451
					l1452:
						position, tokenIndex = position1451, tokenIndex1451
						if buffer[position] != rune('E') {
							goto l1440
						}
						position++
					}
				l1451:
					{
						position1453, tokenIndex1453 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1454
						}
						position++
						goto l1453
					l1454:
						position, tokenIndex = position1453, tokenIndex1453
						if buffer[position] != rune('A') {
							goto l1440
						}
						position++
					}
				l1453:
					{
						position1455, tokenIndex1455 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1456
						}
						position++
						goto l1455
					l1456:
						position, tokenIndex = position1455, tokenIndex1455
						if buffer[position] != rune('M') {
							goto l1440
						}
						position++
					}
				l1455:
					add(rulePegText, position1442)
				}
				if !_rules[ruleAction91]() {
					goto l1440
				}
				add(ruleDSTREAM, position1441)
			}
			return true
		l1440:
			position, tokenIndex = position1440, tokenIndex1440
			return false
		},
		/* 122 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action92)> */
		func() bool {
			position1457, tokenIndex1457 := position, tokenIndex
			{
				position1458 := position
				{
					position1459 := position
					{
						position1460, tokenIndex1460 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1461
						}
						position++
						goto l1460
					l1461:
						position, tokenIndex = position1460, tokenIndex1460
						if buffer[position] != rune('R') {
							goto l1457
						}
						position++
					}
				l1460:
					{
						position1462, tokenIndex1462 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1463
						}
						position++
						goto l1462
					l1463:
						position, tokenIndex = position1462, tokenIndex1462
						if buffer[position] != rune('S') {
							goto l1457
						}
						position++
					}
				l1462:
					{
						position1464, tokenIndex1464 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1465
						}
						position++
						goto l1464
					l1465:
						position, tokenIndex = position1464, tokenIndex1464
						if buffer[position] != rune('T') {
							goto l1457
						}
						position++
					}
				l1464:
					{
						position1466, tokenIndex1466 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1467
						}
						position++
						goto l1466
					l1467:
						position, tokenIndex = position1466, tokenIndex1466
						if buffer[position] != rune('R') {
							goto l1457
						}
						position++
					}
				l1466:
					{
						position1468, tokenIndex1468 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1469
						}
						position++
						goto l1468
					l1469:
						position, tokenIndex = position1468, tokenIndex1468
						if buffer[position] != rune('E') {
							goto l1457
						}
						position++
					}
				l1468:
					{
						position1470, tokenIndex1470 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1471
						}
						position++
						goto l1470
					l1471:
						position, tokenIndex = position1470, tokenIndex1470
						if buffer[position] != rune('A') {
							goto l1457
						}
						position++
					}
				l1470:
					{
						position1472, tokenIndex1472 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1473
						}
						position++
						goto l1472
					l1473:
						position, tokenIndex = position1472, tokenIndex1472
						if buffer[position] != rune('M') {
							goto l1457
						}
						position++
					}
				l1472:
					add(rulePegText, position1459)
				}
				if !_rules[ruleAction92]() {
					goto l1457
				}
				add(ruleRSTREAM, position1458)
			}
			return true
		l1457:
			position, tokenIndex = position1457, tokenIndex1457
			return false
		},
		/* 123 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action93)> */
		func() bool {
			position1474, tokenIndex1474 := position, tokenIndex
			{
				position1475 := position
				{
					position1476 := position
					{
						position1477, tokenIndex1477 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1478
						}
						position++
						goto l1477
					l1478:
						position, tokenIndex = position1477, tokenIndex1477
						if buffer[position] != rune('T') {
							goto l1474
						}
						position++
					}
				l1477:
					{
						position1479, tokenIndex1479 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1480
						}
						position++
						goto l1479
					l1480:
						position, tokenIndex = position1479, tokenIndex1479
						if buffer[position] != rune('U') {
							goto l1474
						}
						position++
					}
				l1479:
					{
						position1481, tokenIndex1481 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1482
						}
						position++
						goto l1481
					l1482:
						position, tokenIndex = position1481, tokenIndex1481
						if buffer[position] != rune('P') {
							goto l1474
						}
						position++
					}
				l1481:
					{
						position1483, tokenIndex1483 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1484
						}
						position++
						goto l1483
					l1484:
						position, tokenIndex = position1483, tokenIndex1483
						if buffer[position] != rune('L') {
							goto l1474
						}
						position++
					}
				l1483:
					{
						position1485, tokenIndex1485 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1486
						}
						position++
						goto l1485
					l1486:
						position, tokenIndex = position1485, tokenIndex1485
						if buffer[position] != rune('E') {
							goto l1474
						}
						position++
					}
				l1485:
					{
						position1487, tokenIndex1487 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1488
						}
						position++
						goto l1487
					l1488:
						position, tokenIndex = position1487, tokenIndex1487
						if buffer[position] != rune('S') {
							goto l1474
						}
						position++
					}
				l1487:
					add(rulePegText, position1476)
				}
				if !_rules[ruleAction93]() {
					goto l1474
				}
				add(ruleTUPLES, position1475)
			}
			return true
		l1474:
			position, tokenIndex = position1474, tokenIndex1474
			return false
		},
		/* 124 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action94)> */
		func() bool {
			position1489, tokenIndex1489 := position, tokenIndex
			{
				position1490 := position
				{
					position1491 := position
					{
						position1492, tokenIndex1492 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1493
						}
						position++
						goto l1492
					l1493:
						position, tokenIndex = position1492, tokenIndex1492
						if buffer[position] != rune('S') {
							goto l1489
						}
						position++
					}
				l1492:
					{
						position1494, tokenIndex1494 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1495
						}
						position++
						goto l1494
					l1495:
						position, tokenIndex = position1494, tokenIndex1494
						if buffer[position] != rune('E') {
							goto l1489
						}
						position++
					}
				l1494:
					{
						position1496, tokenIndex1496 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1497
						}
						position++
						goto l1496
					l1497:
						position, tokenIndex = position1496, tokenIndex1496
						if buffer[position] != rune('C') {
							goto l1489
						}
						position++
					}
				l1496:
					{
						position1498, tokenIndex1498 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1499
						}
						position++
						goto l1498
					l1499:
						position, tokenIndex = position1498, tokenIndex1498
						if buffer[position] != rune('O') {
							goto l1489
						}
						position++
					}
				l1498:
					{
						position1500, tokenIndex1500 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1501
						}
						position++
						goto l1500
					l1501:
						position, tokenIndex = position1500, tokenIndex1500
						if buffer[position] != rune('N') {
							goto l1489
						}
						position++
					}
				l1500:
					{
						position1502, tokenIndex1502 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1503
						}
						position++
						goto l1502
					l1503:
						position, tokenIndex = position1502, tokenIndex1502
						if buffer[position] != rune('D') {
							goto l1489
						}
						position++
					}
				l1502:
					{
						position1504, tokenIndex1504 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1505
						}
						position++
						goto l1504
					l1505:
						position, tokenIndex = position1504, tokenIndex1504
						if buffer[position] != rune('S') {
							goto l1489
						}
						position++
					}
				l1504:
					add(rulePegText, position1491)
				}
				if !_rules[ruleAction94]() {
					goto l1489
				}
				add(ruleSECONDS, position1490)
			}
			return true
		l1489:
			position, tokenIndex = position1489, tokenIndex1489
			return false
		},
		/* 125 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action95)> */
		func() bool {
			position1506, tokenIndex1506 := position, tokenIndex
			{
				position1507 := position
				{
					position1508 := position
					{
						position1509, tokenIndex1509 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1510
						}
						position++
						goto l1509
					l1510:
						position, tokenIndex = position1509, tokenIndex1509
						if buffer[position] != rune('M') {
							goto l1506
						}
						position++
					}
				l1509:
					{
						position1511, tokenIndex1511 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1512
						}
						position++
						goto l1511
					l1512:
						position, tokenIndex = position1511, tokenIndex1511
						if buffer[position] != rune('I') {
							goto l1506
						}
						position++
					}
				l1511:
					{
						position1513, tokenIndex1513 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1514
						}
						position++
						goto l1513
					l1514:
						position, tokenIndex = position1513, tokenIndex1513
						if buffer[position] != rune('L') {
							goto l1506
						}
						position++
					}
				l1513:
					{
						position1515, tokenIndex1515 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1516
						}
						position++
						goto l1515
					l1516:
						position, tokenIndex = position1515, tokenIndex1515
						if buffer[position] != rune('L') {
							goto l1506
						}
						position++
					}
				l1515:
					{
						position1517, tokenIndex1517 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1518
						}
						position++
						goto l1517
					l1518:
						position, tokenIndex = position1517, tokenIndex1517
						if buffer[position] != rune('I') {
							goto l1506
						}
						position++
					}
				l1517:
					{
						position1519, tokenIndex1519 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1520
						}
						position++
						goto l1519
					l1520:
						position, tokenIndex = position1519, tokenIndex1519
						if buffer[position] != rune('S') {
							goto l1506
						}
						position++
					}
				l1519:
					{
						position1521, tokenIndex1521 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1522
						}
						position++
						goto l1521
					l1522:
						position, tokenIndex = position1521, tokenIndex1521
						if buffer[position] != rune('E') {
							goto l1506
						}
						position++
					}
				l1521:
					{
						position1523, tokenIndex1523 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1524
						}
						position++
						goto l1523
					l1524:
						position, tokenIndex = position1523, tokenIndex1523
						if buffer[position] != rune('C') {
							goto l1506
						}
						position++
					}
				l1523:
					{
						position1525, tokenIndex1525 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1526
						}
						position++
						goto l1525
					l1526:
						position, tokenIndex = position1525, tokenIndex1525
						if buffer[position] != rune('O') {
							goto l1506
						}
						position++
					}
				l1525:
					{
						position1527, tokenIndex1527 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1528
						}
						position++
						goto l1527
					l1528:
						position, tokenIndex = position1527, tokenIndex1527
						if buffer[position] != rune('N') {
							goto l1506
						}
						position++
					}
				l1527:
					{
						position1529, tokenIndex1529 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1530
						}
						position++
						goto l1529
					l1530:
						position, tokenIndex = position1529, tokenIndex1529
						if buffer[position] != rune('D') {
							goto l1506
						}
						position++
					}
				l1529:
					{
						position1531, tokenIndex1531 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1532
						}
						position++
						goto l1531
					l1532:
						position, tokenIndex = position1531, tokenIndex1531
						if buffer[position] != rune('S') {
							goto l1506
						}
						position++
					}
				l1531:
					add(rulePegText, position1508)
				}
				if !_rules[ruleAction95]() {
					goto l1506
				}
				add(ruleMILLISECONDS, position1507)
			}
			return true
		l1506:
			position, tokenIndex = position1506, tokenIndex1506
			return false
		},
		/* 126 Wait <- <(<(('w' / 'W') ('a' / 'A') ('i' / 'I') ('t' / 'T'))> Action96)> */
		func() bool {
			position1533, tokenIndex1533 := position, tokenIndex
			{
				position1534 := position
				{
					position1535 := position
					{
						position1536, tokenIndex1536 := position, tokenIndex
						if buffer[position] != rune('w') {
							goto l1537
						}
						position++
						goto l1536
					l1537:
						position, tokenIndex = position1536, tokenIndex1536
						if buffer[position] != rune('W') {
							goto l1533
						}
						position++
					}
				l1536:
					{
						position1538, tokenIndex1538 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1539
						}
						position++
						goto l1538
					l1539:
						position, tokenIndex = position1538, tokenIndex1538
						if buffer[position] != rune('A') {
							goto l1533
						}
						position++
					}
				l1538:
					{
						position1540, tokenIndex1540 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1541
						}
						position++
						goto l1540
					l1541:
						position, tokenIndex = position1540, tokenIndex1540
						if buffer[position] != rune('I') {
							goto l1533
						}
						position++
					}
				l1540:
					{
						position1542, tokenIndex1542 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1543
						}
						position++
						goto l1542
					l1543:
						position, tokenIndex = position1542, tokenIndex1542
						if buffer[position] != rune('T') {
							goto l1533
						}
						position++
					}
				l1542:
					add(rulePegText, position1535)
				}
				if !_rules[ruleAction96]() {
					goto l1533
				}
				add(ruleWait, position1534)
			}
			return true
		l1533:
			position, tokenIndex = position1533, tokenIndex1533
			return false
		},
		/* 127 DropOldest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('o' / 'O') ('l' / 'L') ('d' / 'D') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action97)> */
		func() bool {
			position1544, tokenIndex1544 := position, tokenIndex
			{
				position1545 := position
				{
					position1546 := position
					{
						position1547, tokenIndex1547 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1548
						}
						position++
						goto l1547
					l1548:
						position, tokenIndex = position1547, tokenIndex1547
						if buffer[position] != rune('D') {
							goto l1544
						}
						position++
					}
				l1547:
					{
						position1549, tokenIndex1549 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1550
						}
						position++
						goto l1549
					l1550:
						position, tokenIndex = position1549, tokenIndex1549
						if buffer[position] != rune('R') {
							goto l1544
						}
						position++
					}
				l1549:
					{
						position1551, tokenIndex1551 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1552
						}
						position++
						goto l1551
					l1552:
						position, tokenIndex = position1551, tokenIndex1551
						if buffer[position] != rune('O') {
							goto l1544
						}
						position++
					}
				l1551:
					{
						position1553, tokenIndex1553 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1554
						}
						position++
						goto l1553
					l1554:
						position, tokenIndex = position1553, tokenIndex1553
						if buffer[position] != rune('P') {
							goto l1544
						}
						position++
					}
				l1553:
					if !_rules[rulesp]() {
						goto l1544
					}
					{
						position1555, tokenIndex1555 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1556
						}
						position++
						goto l1555
					l1556:
						position, tokenIndex = position1555, tokenIndex1555
						if buffer[position] != rune('O') {
							goto l1544
						}
						position++
					}
				l1555:
					{
						position1557, tokenIndex1557 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1558
						}
						position++
						goto l1557
					l1558:
						position, tokenIndex = position1557, tokenIndex1557
						if buffer[position] != rune('L') {
							goto l1544
						}
						position++
					}
				l1557:
					{
						position1559, tokenIndex1559 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1560
						}
						position++
						goto l1559
					l1560:
						position, tokenIndex = position1559, tokenIndex1559
						if buffer[position] != rune('D') {
							goto l1544
						}
						position++
					}
				l1559:
					{
						position1561, tokenIndex1561 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1562
						}
						position++
						goto l1561
					l1562:
						position, tokenIndex = position1561, tokenIndex1561
						if buffer[position] != rune('E') {
							goto l1544
						}
						position++
					}
				l1561:
					{
						position1563, tokenIndex1563 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1564
						}
						position++
						goto l1563
					l1564:
						position, tokenIndex = position1563, tokenIndex1563
						if buffer[position] != rune('S') {
							goto l1544
						}
						position++
					}
				l1563:
					{
						position1565, tokenIndex1565 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1566
						}
						position++
						goto l1565
					l1566:
						position, tokenIndex = position1565, tokenIndex1565
						if buffer[position] != rune('T') {
							goto l1544
						}
						position++
					}
				l1565:
					add(rulePegText, position1546)
				}
				if !_rules[ruleAction97]() {
					goto l1544
				}
				add(ruleDropOldest, position1545)
			}
			return true
		l1544:
			position, tokenIndex = position1544, tokenIndex1544
			return false
		},
		/* 128 DropNewest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('n' / 'N') ('e' / 'E') ('w' / 'W') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action98)> */
		func() bool {
			position1567, tokenIndex1567 := position, tokenIndex
			{
				position1568 := position
				{
					position1569 := position
					{
						position1570, tokenIndex1570 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1571
						}
						position++
						goto l1570
					l1571:
						position, tokenIndex = position1570, tokenIndex1570
						if buffer[position] != rune('D') {
							goto l1567
						}
						position++
					}
				l1570:
					{
						position1572, tokenIndex1572 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1573
						}
						position++
						goto l1572
					l1573:
						position, tokenIndex = position1572, tokenIndex1572
						if buffer[position] != rune('R') {
							goto l1567
						}
						position++
					}
				l1572:
					{
						position1574, tokenIndex1574 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1575
						}
						position++
						goto l1574
					l1575:
						position, tokenIndex = position1574, tokenIndex1574
						if buffer[position] != rune('O') {
							goto l1567
						}
						position++
					}
				l1574:
					{
						position1576, tokenIndex1576 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1577
						}
						position++
						goto l1576
					l1577:
						position, tokenIndex = position1576, tokenIndex1576
						if buffer[position] != rune('P') {
							goto l1567
						}
						position++
					}
				l1576:
					if !_rules[rulesp]() {
						goto l1567
					}
					{
						position1578, tokenIndex1578 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1579
						}
						position++
						goto l1578
					l1579:
						position, tokenIndex = position1578, tokenIndex1578
						if buffer[position] != rune('N') {
							goto l1567
						}
						position++
					}
				l1578:
					{
						position1580, tokenIndex1580 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1581
						}
						position++
						goto l1580
					l1581:
						position, tokenIndex = position1580, tokenIndex1580
						if buffer[position] != rune('E') {
							goto l1567
						}
						position++
					}
				l1580:
					{
						position1582, tokenIndex1582 := position, tokenIndex
						if buffer[position] != rune('w') {
							goto l1583
						}
						position++
						goto l1582
					l1583:
						position, tokenIndex = position1582, tokenIndex1582
						if buffer[position] != rune('W') {
							goto l1567
						}
						position++
					}
				l1582:
					{
						position1584, tokenIndex1584 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1585
						}
						position++
						goto l1584
					l1585:
						position, tokenIndex = position1584, tokenIndex1584
						if buffer[position] != rune('E') {
							goto l1567
						}
						position++
					}
				l1584:
					{
						position1586, tokenIndex1586 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1587
						}
						position++
						goto l1586
					l1587:
						position, tokenIndex = position1586, tokenIndex1586
						if buffer[position] != rune('S') {
							goto l1567
						}
						position++
					}
				l1586:
					{
						position1588, tokenIndex1588 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1589
						}
						position++
						goto l1588
					l1589:
						position, tokenIndex = position1588, tokenIndex1588
						if buffer[position] != rune('T') {
							goto l1567
						}
						position++
					}
				l1588:
					add(rulePegText, position1569)
				}
				if !_rules[ruleAction98]() {
					goto l1567
				}
				add(ruleDropNewest, position1568)
			}
			return true
		l1567:
			position, tokenIndex = position1567, tokenIndex1567
			return false
		},
		/* 129 StreamIdentifier <- <(<ident> Action99)> */
		func() bool {
			position1590, tokenIndex1590 := position, tokenIndex
			{
				position1591 := position
				{
					position1592 := position
					if !_rules[ruleident]() {
						goto l1590
					}
					add(rulePegText, position1592)
				}
				if !_rules[ruleAction99]() {
					goto l1590
				}
				add(ruleStreamIdentifier, position1591)
			}
			return true
		l1590:
			position, tokenIndex = position1590, tokenIndex1590
			return false
		},
		/* 130 SourceSinkType <- <(<ident> Action100)> */
		func() bool {
			position1593, tokenIndex1593 := position, tokenIndex
			{
				position1594 := position
				{
					position1595 := position
					if !_rules[ruleident]() {
						goto l1593
					}
					add(rulePegText, position1595)
				}
				if !_rules[ruleAction100]() {
					goto l1593
				}
				add(ruleSourceSinkType, position1594)
			}
			return true
		l1593:
			position, tokenIndex = position1593, tokenIndex1593
			return false
		},
		/* 131 SourceSinkParamKey <- <(<ident> Action101)> */
		func() bool {
			position1596, tokenIndex1596 := position, tokenIndex
			{
				position1597 := position
				{
					position1598 := position
					if !_rules[ruleident]() {
						goto l1596
					}
					add(rulePegText, position1598)
				}
				if !_rules[ruleAction101]() {
					goto l1596
				}
				add(ruleSourceSinkParamKey, position1597)
			}
			return true
		l1596:
			position, tokenIndex = position1596, tokenIndex1596
			return false
		},
		/* 132 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action102)> */
		func() bool {
			position1599, tokenIndex1599 := position, tokenIndex
			{
				position1600 := position
				{
					position1601 := position
					{
						position1602, tokenIndex1602 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1603
						}
						position++
						goto l1602
					l1603:
						position, tokenIndex = position1602, tokenIndex1602
						if buffer[position] != rune('P') {
							goto l1599
						}
						position++
					}
				l1602:
					{
						position1604, tokenIndex1604 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1605
						}
						position++
						goto l1604
					l1605:
						position, tokenIndex = position1604, tokenIndex1604
						if buffer[position] != rune('A') {
							goto l1599
						}
						position++
					}
				l1604:
					{
						position1606, tokenIndex1606 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1607
						}
						position++
						goto l1606
					l1607:
						position, tokenIndex = position1606, tokenIndex1606
						if buffer[position] != rune('U') {
							goto l1599
						}
						position++
					}
				l1606:
					{
						position1608, tokenIndex1608 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1609
						}
						position++
						goto l1608
					l1609:
						position, tokenIndex = position1608, tokenIndex1608
						if buffer[position] != rune('S') {
							goto l1599
						}
						position++
					}
				l1608:
					{
						position1610, tokenIndex1610 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1611
						}
						position++
						goto l1610
					l1611:
						position, tokenIndex = position1610, tokenIndex1610
						if buffer[position] != rune('E') {
							goto l1599
						}
						position++
					}
				l1610:
					{
						position1612, tokenIndex1612 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1613
						}
						position++
						goto l1612
					l1613:
						position, tokenIndex = position1612, tokenIndex1612
						if buffer[position] != rune('D') {
							goto l1599
						}
						position++
					}
				l1612:
					add(rulePegText, position1601)
				}
				if !_rules[ruleAction102]() {
					goto l1599
				}
				add(rulePaused, position1600)
			}
			return true
		l1599:
			position, tokenIndex = position1599, tokenIndex1599
			return false
		},
		/* 133 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action103)> */
		func() bool {
			position1614, tokenIndex1614 := position, tokenIndex
			{
				position1615 := position
				{
					position1616 := position
					{
						position1617, tokenIndex1617 := position, tokenIndex
						if buffer[position] != rune('u') {
							goto l1618
						}
						position++
						goto l1617
					l1618:
						position, tokenIndex = position1617, tokenIndex1617
						if buffer[position] != rune('U') {
							goto l1614
						}
						position++
					}
				l1617:
					{
						position1619, tokenIndex1619 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1620
						}
						position++
						goto l1619
					l1620:
						position, tokenIndex = position1619, tokenIndex1619
						if buffer[position] != rune('N') {
							goto l1614
						}
						position++
					}
				l1619:
					{
						position1621, tokenIndex1621 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1622
						}
						position++
						goto l1621
					l1622:
						position, tokenIndex = position1621, tokenIndex1621
						if buffer[position] != rune('P') {
							goto l1614
						}
						position++
					}
				l1621:
					{
						position1623, tokenIndex1623 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1624
						}
						position++
						goto l1623
					l1624:
						position, tokenIndex = position1623, tokenIndex1623
						if buffer[position] != rune('A') {
							goto l1614
						}
						position++
					}
				l1623:
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
							goto l1614
						}
						position++
					}
				l1625:
					{
						position1627, tokenIndex1627 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1628
						}
						position++
						goto l1627
					l1628:
						position, tokenIndex = position1627, tokenIndex1627
						if buffer[position] != rune('S') {
							goto l1614
						}
						position++
					}
				l1627:
					{
						position1629, tokenIndex1629 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1630
						}
						position++
						goto l1629
					l1630:
						position, tokenIndex = position1629, tokenIndex1629
						if buffer[position] != rune('E') {
							goto l1614
						}
						position++
					}
				l1629:
					{
						position1631, tokenIndex1631 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1632
						}
						position++
						goto l1631
					l1632:
						position, tokenIndex = position1631, tokenIndex1631
						if buffer[position] != rune('D') {
							goto l1614
						}
						position++
					}
				l1631:
					add(rulePegText, position1616)
				}
				if !_rules[ruleAction103]() {
					goto l1614
				}
				add(ruleUnpaused, position1615)
			}
			return true
		l1614:
			position, tokenIndex = position1614, tokenIndex1614
			return false
		},
		/* 134 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action104)> */
		func() bool {
			position1633, tokenIndex1633 := position, tokenIndex
			{
				position1634 := position
				{
					position1635 := position
					{
						position1636, tokenIndex1636 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1637
						}
						position++
						goto l1636
					l1637:
						position, tokenIndex = position1636, tokenIndex1636
						if buffer[position] != rune('A') {
							goto l1633
						}
						position++
					}
				l1636:
					{
						position1638, tokenIndex1638 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1639
						}
						position++
						goto l1638
					l1639:
						position, tokenIndex = position1638, tokenIndex1638
						if buffer[position] != rune('S') {
							goto l1633
						}
						position++
					}
				l1638:
					{
						position1640, tokenIndex1640 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1641
						}
						position++
						goto l1640
					l1641:
						position, tokenIndex = position1640, tokenIndex1640
						if buffer[position] != rune('C') {
							goto l1633
						}
						position++
					}
				l1640:
					add(rulePegText, position1635)
				}
				if !_rules[ruleAction104]() {
					goto l1633
				}
				add(ruleAscending, position1634)
			}
			return true
		l1633:
			position, tokenIndex = position1633, tokenIndex1633
			return false
		},
		/* 135 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action105)> */
		func() bool {
			position1642, tokenIndex1642 := position, tokenIndex
			{
				position1643 := position
				{
					position1644 := position
					{
						position1645, tokenIndex1645 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1646
						}
						position++
						goto l1645
					l1646:
						position, tokenIndex = position1645, tokenIndex1645
						if buffer[position] != rune('D') {
							goto l1642
						}
						position++
					}
				l1645:
					{
						position1647, tokenIndex1647 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1648
						}
						position++
						goto l1647
					l1648:
						position, tokenIndex = position1647, tokenIndex1647
						if buffer[position] != rune('E') {
							goto l1642
						}
						position++
					}
				l1647:
					{
						position1649, tokenIndex1649 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1650
						}
						position++
						goto l1649
					l1650:
						position, tokenIndex = position1649, tokenIndex1649
						if buffer[position] != rune('S') {
							goto l1642
						}
						position++
					}
				l1649:
					{
						position1651, tokenIndex1651 := position, tokenIndex
						if buffer[position] != rune('c') {
							goto l1652
						}
						position++
						goto l1651
					l1652:
						position, tokenIndex = position1651, tokenIndex1651
						if buffer[position] != rune('C') {
							goto l1642
						}
						position++
					}
				l1651:
					add(rulePegText, position1644)
				}
				if !_rules[ruleAction105]() {
					goto l1642
				}
				add(ruleDescending, position1643)
			}
			return true
		l1642:
			position, tokenIndex = position1642, tokenIndex1642
			return false
		},
		/* 136 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1653, tokenIndex1653 := position, tokenIndex
			{
				position1654 := position
				{
					position1655, tokenIndex1655 := position, tokenIndex
					if !_rules[ruleBool]() {
						goto l1656
					}
					goto l1655
				l1656:
					position, tokenIndex = position1655, tokenIndex1655
					if !_rules[ruleInt]() {
						goto l1657
					}
					goto l1655
				l1657:
					position, tokenIndex = position1655, tokenIndex1655
					if !_rules[ruleFloat]() {
						goto l1658
					}
					goto l1655
				l1658:
					position, tokenIndex = position1655, tokenIndex1655
					if !_rules[ruleString]() {
						goto l1659
					}
					goto l1655
				l1659:
					position, tokenIndex = position1655, tokenIndex1655
					if !_rules[ruleBlob]() {
						goto l1660
					}
					goto l1655
				l1660:
					position, tokenIndex = position1655, tokenIndex1655
					if !_rules[ruleTimestamp]() {
						goto l1661
					}
					goto l1655
				l1661:
					position, tokenIndex = position1655, tokenIndex1655
					if !_rules[ruleArray]() {
						goto l1662
					}
					goto l1655
				l1662:
					position, tokenIndex = position1655, tokenIndex1655
					if !_rules[ruleMap]() {
						goto l1653
					}
				}
			l1655:
				add(ruleType, position1654)
			}
			return true
		l1653:
			position, tokenIndex = position1653, tokenIndex1653
			return false
		},
		/* 137 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action106)> */
		func() bool {
			position1663, tokenIndex1663 := position, tokenIndex
			{
				position1664 := position
				{
					position1665 := position
					{
						position1666, tokenIndex1666 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1667
						}
						position++
						goto l1666
					l1667:
						position, tokenIndex = position1666, tokenIndex1666
						if buffer[position] != rune('B') {
							goto l1663
						}
						position++
					}
				l1666:
					{
						position1668, tokenIndex1668 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1669
						}
						position++
						goto l1668
					l1669:
						position, tokenIndex = position1668, tokenIndex1668
						if buffer[position] != rune('O') {
							goto l1663
						}
						position++
					}
				l1668:
					{
						position1670, tokenIndex1670 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1671
						}
						position++
						goto l1670
					l1671:
						position, tokenIndex = position1670, tokenIndex1670
						if buffer[position] != rune('O') {
							goto l1663
						}
						position++
					}
				l1670:
					{
						position1672, tokenIndex1672 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1673
						}
						position++
						goto l1672
					l1673:
						position, tokenIndex = position1672, tokenIndex1672
						if buffer[position] != rune('L') {
							goto l1663
						}
						position++
					}
				l1672:
					add(rulePegText, position1665)
				}
				if !_rules[ruleAction106]() {
					goto l1663
				}
				add(ruleBool, position1664)
			}
			return true
		l1663:
			position, tokenIndex = position1663, tokenIndex1663
			return false
		},
		/* 138 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action107)> */
		func() bool {
			position1674, tokenIndex1674 := position, tokenIndex
			{
				position1675 := position
				{
					position1676 := position
					{
						position1677, tokenIndex1677 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1678
						}
						position++
						goto l1677
					l1678:
						position, tokenIndex = position1677, tokenIndex1677
						if buffer[position] != rune('I') {
							goto l1674
						}
						position++
					}
				l1677:
					{
						position1679, tokenIndex1679 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1680
						}
						position++
						goto l1679
					l1680:
						position, tokenIndex = position1679, tokenIndex1679
						if buffer[position] != rune('N') {
							goto l1674
						}
						position++
					}
				l1679:
					{
						position1681, tokenIndex1681 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1682
						}
						position++
						goto l1681
					l1682:
						position, tokenIndex = position1681, tokenIndex1681
						if buffer[position] != rune('T') {
							goto l1674
						}
						position++
					}
				l1681:
					add(rulePegText, position1676)
				}
				if !_rules[ruleAction107]() {
					goto l1674
				}
				add(ruleInt, position1675)
			}
			return true
		l1674:
			position, tokenIndex = position1674, tokenIndex1674
			return false
		},
		/* 139 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action108)> */
		func() bool {
			position1683, tokenIndex1683 := position, tokenIndex
			{
				position1684 := position
				{
					position1685 := position
					{
						position1686, tokenIndex1686 := position, tokenIndex
						if buffer[position] != rune('f') {
							goto l1687
						}
						position++
						goto l1686
					l1687:
						position, tokenIndex = position1686, tokenIndex1686
						if buffer[position] != rune('F') {
							goto l1683
						}
						position++
					}
				l1686:
					{
						position1688, tokenIndex1688 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1689
						}
						position++
						goto l1688
					l1689:
						position, tokenIndex = position1688, tokenIndex1688
						if buffer[position] != rune('L') {
							goto l1683
						}
						position++
					}
				l1688:
					{
						position1690, tokenIndex1690 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1691
						}
						position++
						goto l1690
					l1691:
						position, tokenIndex = position1690, tokenIndex1690
						if buffer[position] != rune('O') {
							goto l1683
						}
						position++
					}
				l1690:
					{
						position1692, tokenIndex1692 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1693
						}
						position++
						goto l1692
					l1693:
						position, tokenIndex = position1692, tokenIndex1692
						if buffer[position] != rune('A') {
							goto l1683
						}
						position++
					}
				l1692:
					{
						position1694, tokenIndex1694 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1695
						}
						position++
						goto l1694
					l1695:
						position, tokenIndex = position1694, tokenIndex1694
						if buffer[position] != rune('T') {
							goto l1683
						}
						position++
					}
				l1694:
					add(rulePegText, position1685)
				}
				if !_rules[ruleAction108]() {
					goto l1683
				}
				add(ruleFloat, position1684)
			}
			return true
		l1683:
			position, tokenIndex = position1683, tokenIndex1683
			return false
		},
		/* 140 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action109)> */
		func() bool {
			position1696, tokenIndex1696 := position, tokenIndex
			{
				position1697 := position
				{
					position1698 := position
					{
						position1699, tokenIndex1699 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1700
						}
						position++
						goto l1699
					l1700:
						position, tokenIndex = position1699, tokenIndex1699
						if buffer[position] != rune('S') {
							goto l1696
						}
						position++
					}
				l1699:
					{
						position1701, tokenIndex1701 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1702
						}
						position++
						goto l1701
					l1702:
						position, tokenIndex = position1701, tokenIndex1701
						if buffer[position] != rune('T') {
							goto l1696
						}
						position++
					}
				l1701:
					{
						position1703, tokenIndex1703 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1704
						}
						position++
						goto l1703
					l1704:
						position, tokenIndex = position1703, tokenIndex1703
						if buffer[position] != rune('R') {
							goto l1696
						}
						position++
					}
				l1703:
					{
						position1705, tokenIndex1705 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1706
						}
						position++
						goto l1705
					l1706:
						position, tokenIndex = position1705, tokenIndex1705
						if buffer[position] != rune('I') {
							goto l1696
						}
						position++
					}
				l1705:
					{
						position1707, tokenIndex1707 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1708
						}
						position++
						goto l1707
					l1708:
						position, tokenIndex = position1707, tokenIndex1707
						if buffer[position] != rune('N') {
							goto l1696
						}
						position++
					}
				l1707:
					{
						position1709, tokenIndex1709 := position, tokenIndex
						if buffer[position] != rune('g') {
							goto l1710
						}
						position++
						goto l1709
					l1710:
						position, tokenIndex = position1709, tokenIndex1709
						if buffer[position] != rune('G') {
							goto l1696
						}
						position++
					}
				l1709:
					add(rulePegText, position1698)
				}
				if !_rules[ruleAction109]() {
					goto l1696
				}
				add(ruleString, position1697)
			}
			return true
		l1696:
			position, tokenIndex = position1696, tokenIndex1696
			return false
		},
		/* 141 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action110)> */
		func() bool {
			position1711, tokenIndex1711 := position, tokenIndex
			{
				position1712 := position
				{
					position1713 := position
					{
						position1714, tokenIndex1714 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1715
						}
						position++
						goto l1714
					l1715:
						position, tokenIndex = position1714, tokenIndex1714
						if buffer[position] != rune('B') {
							goto l1711
						}
						position++
					}
				l1714:
					{
						position1716, tokenIndex1716 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l1717
						}
						position++
						goto l1716
					l1717:
						position, tokenIndex = position1716, tokenIndex1716
						if buffer[position] != rune('L') {
							goto l1711
						}
						position++
					}
				l1716:
					{
						position1718, tokenIndex1718 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1719
						}
						position++
						goto l1718
					l1719:
						position, tokenIndex = position1718, tokenIndex1718
						if buffer[position] != rune('O') {
							goto l1711
						}
						position++
					}
				l1718:
					{
						position1720, tokenIndex1720 := position, tokenIndex
						if buffer[position] != rune('b') {
							goto l1721
						}
						position++
						goto l1720
					l1721:
						position, tokenIndex = position1720, tokenIndex1720
						if buffer[position] != rune('B') {
							goto l1711
						}
						position++
					}
				l1720:
					add(rulePegText, position1713)
				}
				if !_rules[ruleAction110]() {
					goto l1711
				}
				add(ruleBlob, position1712)
			}
			return true
		l1711:
			position, tokenIndex = position1711, tokenIndex1711
			return false
		},
		/* 142 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action111)> */
		func() bool {
			position1722, tokenIndex1722 := position, tokenIndex
			{
				position1723 := position
				{
					position1724 := position
					{
						position1725, tokenIndex1725 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1726
						}
						position++
						goto l1725
					l1726:
						position, tokenIndex = position1725, tokenIndex1725
						if buffer[position] != rune('T') {
							goto l1722
						}
						position++
					}
				l1725:
					{
						position1727, tokenIndex1727 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1728
						}
						position++
						goto l1727
					l1728:
						position, tokenIndex = position1727, tokenIndex1727
						if buffer[position] != rune('I') {
							goto l1722
						}
						position++
					}
				l1727:
					{
						position1729, tokenIndex1729 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1730
						}
						position++
						goto l1729
					l1730:
						position, tokenIndex = position1729, tokenIndex1729
						if buffer[position] != rune('M') {
							goto l1722
						}
						position++
					}
				l1729:
					{
						position1731, tokenIndex1731 := position, tokenIndex
						if buffer[position] != rune('e') {
							goto l1732
						}
						position++
						goto l1731
					l1732:
						position, tokenIndex = position1731, tokenIndex1731
						if buffer[position] != rune('E') {
							goto l1722
						}
						position++
					}
				l1731:
					{
						position1733, tokenIndex1733 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1734
						}
						position++
						goto l1733
					l1734:
						position, tokenIndex = position1733, tokenIndex1733
						if buffer[position] != rune('S') {
							goto l1722
						}
						position++
					}
				l1733:
					{
						position1735, tokenIndex1735 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1736
						}
						position++
						goto l1735
					l1736:
						position, tokenIndex = position1735, tokenIndex1735
						if buffer[position] != rune('T') {
							goto l1722
						}
						position++
					}
				l1735:
					{
						position1737, tokenIndex1737 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1738
						}
						position++
						goto l1737
					l1738:
						position, tokenIndex = position1737, tokenIndex1737
						if buffer[position] != rune('A') {
							goto l1722
						}
						position++
					}
				l1737:
					{
						position1739, tokenIndex1739 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1740
						}
						position++
						goto l1739
					l1740:
						position, tokenIndex = position1739, tokenIndex1739
						if buffer[position] != rune('M') {
							goto l1722
						}
						position++
					}
				l1739:
					{
						position1741, tokenIndex1741 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1742
						}
						position++
						goto l1741
					l1742:
						position, tokenIndex = position1741, tokenIndex1741
						if buffer[position] != rune('P') {
							goto l1722
						}
						position++
					}
				l1741:
					add(rulePegText, position1724)
				}
				if !_rules[ruleAction111]() {
					goto l1722
				}
				add(ruleTimestamp, position1723)
			}
			return true
		l1722:
			position, tokenIndex = position1722, tokenIndex1722
			return false
		},
		/* 143 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action112)> */
		func() bool {
			position1743, tokenIndex1743 := position, tokenIndex
			{
				position1744 := position
				{
					position1745 := position
					{
						position1746, tokenIndex1746 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1747
						}
						position++
						goto l1746
					l1747:
						position, tokenIndex = position1746, tokenIndex1746
						if buffer[position] != rune('A') {
							goto l1743
						}
						position++
					}
				l1746:
					{
						position1748, tokenIndex1748 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1749
						}
						position++
						goto l1748
					l1749:
						position, tokenIndex = position1748, tokenIndex1748
						if buffer[position] != rune('R') {
							goto l1743
						}
						position++
					}
				l1748:
					{
						position1750, tokenIndex1750 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1751
						}
						position++
						goto l1750
					l1751:
						position, tokenIndex = position1750, tokenIndex1750
						if buffer[position] != rune('R') {
							goto l1743
						}
						position++
					}
				l1750:
					{
						position1752, tokenIndex1752 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1753
						}
						position++
						goto l1752
					l1753:
						position, tokenIndex = position1752, tokenIndex1752
						if buffer[position] != rune('A') {
							goto l1743
						}
						position++
					}
				l1752:
					{
						position1754, tokenIndex1754 := position, tokenIndex
						if buffer[position] != rune('y') {
							goto l1755
						}
						position++
						goto l1754
					l1755:
						position, tokenIndex = position1754, tokenIndex1754
						if buffer[position] != rune('Y') {
							goto l1743
						}
						position++
					}
				l1754:
					add(rulePegText, position1745)
				}
				if !_rules[ruleAction112]() {
					goto l1743
				}
				add(ruleArray, position1744)
			}
			return true
		l1743:
			position, tokenIndex = position1743, tokenIndex1743
			return false
		},
		/* 144 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action113)> */
		func() bool {
			position1756, tokenIndex1756 := position, tokenIndex
			{
				position1757 := position
				{
					position1758 := position
					{
						position1759, tokenIndex1759 := position, tokenIndex
						if buffer[position] != rune('m') {
							goto l1760
						}
						position++
						goto l1759
					l1760:
						position, tokenIndex = position1759, tokenIndex1759
						if buffer[position] != rune('M') {
							goto l1756
						}
						position++
					}
				l1759:
					{
						position1761, tokenIndex1761 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1762
						}
						position++
						goto l1761
					l1762:
						position, tokenIndex = position1761, tokenIndex1761
						if buffer[position] != rune('A') {
							goto l1756
						}
						position++
					}
				l1761:
					{
						position1763, tokenIndex1763 := position, tokenIndex
						if buffer[position] != rune('p') {
							goto l1764
						}
						position++
						goto l1763
					l1764:
						position, tokenIndex = position1763, tokenIndex1763
						if buffer[position] != rune('P') {
							goto l1756
						}
						position++
					}
				l1763:
					add(rulePegText, position1758)
				}
				if !_rules[ruleAction113]() {
					goto l1756
				}
				add(ruleMap, position1757)
			}
			return true
		l1756:
			position, tokenIndex = position1756, tokenIndex1756
			return false
		},
		/* 145 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action114)> */
		func() bool {
			position1765, tokenIndex1765 := position, tokenIndex
			{
				position1766 := position
				{
					position1767 := position
					{
						position1768, tokenIndex1768 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1769
						}
						position++
						goto l1768
					l1769:
						position, tokenIndex = position1768, tokenIndex1768
						if buffer[position] != rune('O') {
							goto l1765
						}
						position++
					}
				l1768:
					{
						position1770, tokenIndex1770 := position, tokenIndex
						if buffer[position] != rune('r') {
							goto l1771
						}
						position++
						goto l1770
					l1771:
						position, tokenIndex = position1770, tokenIndex1770
						if buffer[position] != rune('R') {
							goto l1765
						}
						position++
					}
				l1770:
					add(rulePegText, position1767)
				}
				if !_rules[ruleAction114]() {
					goto l1765
				}
				add(ruleOr, position1766)
			}
			return true
		l1765:
			position, tokenIndex = position1765, tokenIndex1765
			return false
		},
		/* 146 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action115)> */
		func() bool {
			position1772, tokenIndex1772 := position, tokenIndex
			{
				position1773 := position
				{
					position1774 := position
					{
						position1775, tokenIndex1775 := position, tokenIndex
						if buffer[position] != rune('a') {
							goto l1776
						}
						position++
						goto l1775
					l1776:
						position, tokenIndex = position1775, tokenIndex1775
						if buffer[position] != rune('A') {
							goto l1772
						}
						position++
					}
				l1775:
					{
						position1777, tokenIndex1777 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1778
						}
						position++
						goto l1777
					l1778:
						position, tokenIndex = position1777, tokenIndex1777
						if buffer[position] != rune('N') {
							goto l1772
						}
						position++
					}
				l1777:
					{
						position1779, tokenIndex1779 := position, tokenIndex
						if buffer[position] != rune('d') {
							goto l1780
						}
						position++
						goto l1779
					l1780:
						position, tokenIndex = position1779, tokenIndex1779
						if buffer[position] != rune('D') {
							goto l1772
						}
						position++
					}
				l1779:
					add(rulePegText, position1774)
				}
				if !_rules[ruleAction115]() {
					goto l1772
				}
				add(ruleAnd, position1773)
			}
			return true
		l1772:
			position, tokenIndex = position1772, tokenIndex1772
			return false
		},
		/* 147 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action116)> */
		func() bool {
			position1781, tokenIndex1781 := position, tokenIndex
			{
				position1782 := position
				{
					position1783 := position
					{
						position1784, tokenIndex1784 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1785
						}
						position++
						goto l1784
					l1785:
						position, tokenIndex = position1784, tokenIndex1784
						if buffer[position] != rune('N') {
							goto l1781
						}
						position++
					}
				l1784:
					{
						position1786, tokenIndex1786 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1787
						}
						position++
						goto l1786
					l1787:
						position, tokenIndex = position1786, tokenIndex1786
						if buffer[position] != rune('O') {
							goto l1781
						}
						position++
					}
				l1786:
					{
						position1788, tokenIndex1788 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1789
						}
						position++
						goto l1788
					l1789:
						position, tokenIndex = position1788, tokenIndex1788
						if buffer[position] != rune('T') {
							goto l1781
						}
						position++
					}
				l1788:
					add(rulePegText, position1783)
				}
				if !_rules[ruleAction116]() {
					goto l1781
				}
				add(ruleNot, position1782)
			}
			return true
		l1781:
			position, tokenIndex = position1781, tokenIndex1781
			return false
		},
		/* 148 Equal <- <(<'='> Action117)> */
		func() bool {
			position1790, tokenIndex1790 := position, tokenIndex
			{
				position1791 := position
				{
					position1792 := position
					if buffer[position] != rune('=') {
						goto l1790
					}
					position++
					add(rulePegText, position1792)
				}
				if !_rules[ruleAction117]() {
					goto l1790
				}
				add(ruleEqual, position1791)
			}
			return true
		l1790:
			position, tokenIndex = position1790, tokenIndex1790
			return false
		},
		/* 149 Less <- <(<'<'> Action118)> */
		func() bool {
			position1793, tokenIndex1793 := position, tokenIndex
			{
				position1794 := position
				{
					position1795 := position
					if buffer[position] != rune('<') {
						goto l1793
					}
					position++
					add(rulePegText, position1795)
				}
				if !_rules[ruleAction118]() {
					goto l1793
				}
				add(ruleLess, position1794)
			}
			return true
		l1793:
			position, tokenIndex = position1793, tokenIndex1793
			return false
		},
		/* 150 LessOrEqual <- <(<('<' '=')> Action119)> */
		func() bool {
			position1796, tokenIndex1796 := position, tokenIndex
			{
				position1797 := position
				{
					position1798 := position
					if buffer[position] != rune('<') {
						goto l1796
					}
					position++
					if buffer[position] != rune('=') {
						goto l1796
					}
					position++
					add(rulePegText, position1798)
				}
				if !_rules[ruleAction119]() {
					goto l1796
				}
				add(ruleLessOrEqual, position1797)
			}
			return true
		l1796:
			position, tokenIndex = position1796, tokenIndex1796
			return false
		},
		/* 151 Greater <- <(<'>'> Action120)> */
		func() bool {
			position1799, tokenIndex1799 := position, tokenIndex
			{
				position1800 := position
				{
					position1801 := position
					if buffer[position] != rune('>') {
						goto l1799
					}
					position++
					add(rulePegText, position1801)
				}
				if !_rules[ruleAction120]() {
					goto l1799
				}
				add(ruleGreater, position1800)
			}
			return true
		l1799:
			position, tokenIndex = position1799, tokenIndex1799
			return false
		},
		/* 152 GreaterOrEqual <- <(<('>' '=')> Action121)> */
		func() bool {
			position1802, tokenIndex1802 := position, tokenIndex
			{
				position1803 := position
				{
					position1804 := position
					if buffer[position] != rune('>') {
						goto l1802
					}
					position++
					if buffer[position] != rune('=') {
						goto l1802
					}
					position++
					add(rulePegText, position1804)
				}
				if !_rules[ruleAction121]() {
					goto l1802
				}
				add(ruleGreaterOrEqual, position1803)
			}
			return true
		l1802:
			position, tokenIndex = position1802, tokenIndex1802
			return false
		},
		/* 153 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action122)> */
		func() bool {
			position1805, tokenIndex1805 := position, tokenIndex
			{
				position1806 := position
				{
					position1807 := position
					{
						position1808, tokenIndex1808 := position, tokenIndex
						if buffer[position] != rune('!') {
							goto l1809
						}
						position++
						if buffer[position] != rune('=') {
							goto l1809
						}
						position++
						goto l1808
					l1809:
						position, tokenIndex = position1808, tokenIndex1808
						if buffer[position] != rune('<') {
							goto l1805
						}
						position++
						if buffer[position] != rune('>') {
							goto l1805
						}
						position++
					}
				l1808:
					add(rulePegText, position1807)
				}
				if !_rules[ruleAction122]() {
					goto l1805
				}
				add(ruleNotEqual, position1806)
			}
			return true
		l1805:
			position, tokenIndex = position1805, tokenIndex1805
			return false
		},
		/* 154 Concat <- <(<('|' '|')> Action123)> */
		func() bool {
			position1810, tokenIndex1810 := position, tokenIndex
			{
				position1811 := position
				{
					position1812 := position
					if buffer[position] != rune('|') {
						goto l1810
					}
					position++
					if buffer[position] != rune('|') {
						goto l1810
					}
					position++
					add(rulePegText, position1812)
				}
				if !_rules[ruleAction123]() {
					goto l1810
				}
				add(ruleConcat, position1811)
			}
			return true
		l1810:
			position, tokenIndex = position1810, tokenIndex1810
			return false
		},
		/* 155 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action124)> */
		func() bool {
			position1813, tokenIndex1813 := position, tokenIndex
			{
				position1814 := position
				{
					position1815 := position
					{
						position1816, tokenIndex1816 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1817
						}
						position++
						goto l1816
					l1817:
						position, tokenIndex = position1816, tokenIndex1816
						if buffer[position] != rune('I') {
							goto l1813
						}
						position++
					}
				l1816:
					{
						position1818, tokenIndex1818 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1819
						}
						position++
						goto l1818
					l1819:
						position, tokenIndex = position1818, tokenIndex1818
						if buffer[position] != rune('S') {
							goto l1813
						}
						position++
					}
				l1818:
					add(rulePegText, position1815)
				}
				if !_rules[ruleAction124]() {
					goto l1813
				}
				add(ruleIs, position1814)
			}
			return true
		l1813:
			position, tokenIndex = position1813, tokenIndex1813
			return false
		},
		/* 156 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action125)> */
		func() bool {
			position1820, tokenIndex1820 := position, tokenIndex
			{
				position1821 := position
				{
					position1822 := position
					{
						position1823, tokenIndex1823 := position, tokenIndex
						if buffer[position] != rune('i') {
							goto l1824
						}
						position++
						goto l1823
					l1824:
						position, tokenIndex = position1823, tokenIndex1823
						if buffer[position] != rune('I') {
							goto l1820
						}
						position++
					}
				l1823:
					{
						position1825, tokenIndex1825 := position, tokenIndex
						if buffer[position] != rune('s') {
							goto l1826
						}
						position++
						goto l1825
					l1826:
						position, tokenIndex = position1825, tokenIndex1825
						if buffer[position] != rune('S') {
							goto l1820
						}
						position++
					}
				l1825:
					if !_rules[rulesp]() {
						goto l1820
					}
					{
						position1827, tokenIndex1827 := position, tokenIndex
						if buffer[position] != rune('n') {
							goto l1828
						}
						position++
						goto l1827
					l1828:
						position, tokenIndex = position1827, tokenIndex1827
						if buffer[position] != rune('N') {
							goto l1820
						}
						position++
					}
				l1827:
					{
						position1829, tokenIndex1829 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l1830
						}
						position++
						goto l1829
					l1830:
						position, tokenIndex = position1829, tokenIndex1829
						if buffer[position] != rune('O') {
							goto l1820
						}
						position++
					}
				l1829:
					{
						position1831, tokenIndex1831 := position, tokenIndex
						if buffer[position] != rune('t') {
							goto l1832
						}
						position++
						goto l1831
					l1832:
						position, tokenIndex = position1831, tokenIndex1831
						if buffer[position] != rune('T') {
							goto l1820
						}
						position++
					}
				l1831:
					add(rulePegText, position1822)
				}
				if !_rules[ruleAction125]() {
					goto l1820
				}
				add(ruleIsNot, position1821)
			}
			return true
		l1820:
			position, tokenIndex = position1820, tokenIndex1820
			return false
		},
		/* 157 Plus <- <(<'+'> Action126)> */
		func() bool {
			position1833, tokenIndex1833 := position, tokenIndex
			{
				position1834 := position
				{
					position1835 := position
					if buffer[position] != rune('+') {
						goto l1833
					}
					position++
					add(rulePegText, position1835)
				}
				if !_rules[ruleAction126]() {
					goto l1833
				}
				add(rulePlus, position1834)
			}
			return true
		l1833:
			position, tokenIndex = position1833, tokenIndex1833
			return false
		},
		/* 158 Minus <- <(<'-'> Action127)> */
		func() bool {
			position1836, tokenIndex1836 := position, tokenIndex
			{
				position1837 := position
				{
					position1838 := position
					if buffer[position] != rune('-') {
						goto l1836
					}
					position++
					add(rulePegText, position1838)
				}
				if !_rules[ruleAction127]() {
					goto l1836
				}
				add(ruleMinus, position1837)
			}
			return true
		l1836:
			position, tokenIndex = position1836, tokenIndex1836
			return false
		},
		/* 159 Multiply <- <(<'*'> Action128)> */
		func() bool {
			position1839, tokenIndex1839 := position, tokenIndex
			{
				position1840 := position
				{
					position1841 := position
					if buffer[position] != rune('*') {
						goto l1839
					}
					position++
					add(rulePegText, position1841)
				}
				if !_rules[ruleAction128]() {
					goto l1839
				}
				add(ruleMultiply, position1840)
			}
			return true
		l1839:
			position, tokenIndex = position1839, tokenIndex1839
			return false
		},
		/* 160 Divide <- <(<'/'> Action129)> */
		func() bool {
			position1842, tokenIndex1842 := position, tokenIndex
			{
				position1843 := position
				{
					position1844 := position
					if buffer[position] != rune('/') {
						goto l1842
					}
					position++
					add(rulePegText, position1844)
				}
				if !_rules[ruleAction129]() {
					goto l1842
				}
				add(ruleDivide, position1843)
			}
			return true
		l1842:
			position, tokenIndex = position1842, tokenIndex1842
			return false
		},
		/* 161 Modulo <- <(<'%'> Action130)> */
		func() bool {
			position1845, tokenIndex1845 := position, tokenIndex
			{
				position1846 := position
				{
					position1847 := position
					if buffer[position] != rune('%') {
						goto l1845
					}
					position++
					add(rulePegText, position1847)
				}
				if !_rules[ruleAction130]() {
					goto l1845
				}
				add(ruleModulo, position1846)
			}
			return true
		l1845:
			position, tokenIndex = position1845, tokenIndex1845
			return false
		},
		/* 162 UnaryMinus <- <(<'-'> Action131)> */
		func() bool {
			position1848, tokenIndex1848 := position, tokenIndex
			{
				position1849 := position
				{
					position1850 := position
					if buffer[position] != rune('-') {
						goto l1848
					}
					position++
					add(rulePegText, position1850)
				}
				if !_rules[ruleAction131]() {
					goto l1848
				}
				add(ruleUnaryMinus, position1849)
			}
			return true
		l1848:
			position, tokenIndex = position1848, tokenIndex1848
			return false
		},
		/* 163 Identifier <- <(<ident> Action132)> */
		func() bool {
			position1851, tokenIndex1851 := position, tokenIndex
			{
				position1852 := position
				{
					position1853 := position
					if !_rules[ruleident]() {
						goto l1851
					}
					add(rulePegText, position1853)
				}
				if !_rules[ruleAction132]() {
					goto l1851
				}
				add(ruleIdentifier, position1852)
			}
			return true
		l1851:
			position, tokenIndex = position1851, tokenIndex1851
			return false
		},
		/* 164 TargetIdentifier <- <(<('*' / jsonSetPath)> Action133)> */
		func() bool {
			position1854, tokenIndex1854 := position, tokenIndex
			{
				position1855 := position
				{
					position1856 := position
					{
						position1857, tokenIndex1857 := position, tokenIndex
						if buffer[position] != rune('*') {
							goto l1858
						}
						position++
						goto l1857
					l1858:
						position, tokenIndex = position1857, tokenIndex1857
						if !_rules[rulejsonSetPath]() {
							goto l1854
						}
					}
				l1857:
					add(rulePegText, position1856)
				}
				if !_rules[ruleAction133]() {
					goto l1854
				}
				add(ruleTargetIdentifier, position1855)
			}
			return true
		l1854:
			position, tokenIndex = position1854, tokenIndex1854
			return false
		},
		/* 165 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1859, tokenIndex1859 := position, tokenIndex
			{
				position1860 := position
				{
					position1861, tokenIndex1861 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1862
					}
					position++
					goto l1861
				l1862:
					position, tokenIndex = position1861, tokenIndex1861
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1859
					}
					position++
				}
			l1861:
			l1863:
				{
					position1864, tokenIndex1864 := position, tokenIndex
					{
						position1865, tokenIndex1865 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1866
						}
						position++
						goto l1865
					l1866:
						position, tokenIndex = position1865, tokenIndex1865
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1867
						}
						position++
						goto l1865
					l1867:
						position, tokenIndex = position1865, tokenIndex1865
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1868
						}
						position++
						goto l1865
					l1868:
						position, tokenIndex = position1865, tokenIndex1865
						if buffer[position] != rune('_') {
							goto l1864
						}
						position++
					}
				l1865:
					goto l1863
				l1864:
					position, tokenIndex = position1864, tokenIndex1864
				}
				add(ruleident, position1860)
			}
			return true
		l1859:
			position, tokenIndex = position1859, tokenIndex1859
			return false
		},
		/* 166 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1869, tokenIndex1869 := position, tokenIndex
			{
				position1870 := position
				if !_rules[rulejsonPathHead]() {
					goto l1869
				}
			l1871:
				{
					position1872, tokenIndex1872 := position, tokenIndex
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1872
					}
					goto l1871
				l1872:
					position, tokenIndex = position1872, tokenIndex1872
				}
				add(rulejsonGetPath, position1870)
			}
			return true
		l1869:
			position, tokenIndex = position1869, tokenIndex1869
			return false
		},
		/* 167 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1873, tokenIndex1873 := position, tokenIndex
			{
				position1874 := position
				if !_rules[rulejsonPathHead]() {
					goto l1873
				}
			l1875:
				{
					position1876, tokenIndex1876 := position, tokenIndex
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1876
					}
					goto l1875
				l1876:
					position, tokenIndex = position1876, tokenIndex1876
				}
				add(rulejsonSetPath, position1874)
			}
			return true
		l1873:
			position, tokenIndex = position1873, tokenIndex1873
			return false
		},
		/* 168 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1877, tokenIndex1877 := position, tokenIndex
			{
				position1878 := position
				{
					position1879, tokenIndex1879 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l1880
					}
					goto l1879
				l1880:
					position, tokenIndex = position1879, tokenIndex1879
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1877
					}
				}
			l1879:
				add(rulejsonPathHead, position1878)
			}
			return true
		l1877:
			position, tokenIndex = position1877, tokenIndex1877
			return false
		},
		/* 169 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1881, tokenIndex1881 := position, tokenIndex
			{
				position1882 := position
				{
					position1883, tokenIndex1883 := position, tokenIndex
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1884
					}
					goto l1883
				l1884:
					position, tokenIndex = position1883, tokenIndex1883
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1885
					}
					goto l1883
				l1885:
					position, tokenIndex = position1883, tokenIndex1883
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1886
					}
					goto l1883
				l1886:
					position, tokenIndex = position1883, tokenIndex1883
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1887
					}
					goto l1883
				l1887:
					position, tokenIndex = position1883, tokenIndex1883
					if !_rules[rulejsonArraySlice]() {
						goto l1888
					}
					goto l1883
				l1888:
					position, tokenIndex = position1883, tokenIndex1883
					if !_rules[rulejsonArrayAccess]() {
						goto l1881
					}
				}
			l1883:
				add(rulejsonGetPathNonHead, position1882)
			}
			return true
		l1881:
			position, tokenIndex = position1881, tokenIndex1881
			return false
		},
		/* 170 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1889, tokenIndex1889 := position, tokenIndex
			{
				position1890 := position
				{
					position1891, tokenIndex1891 := position, tokenIndex
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1892
					}
					goto l1891
				l1892:
					position, tokenIndex = position1891, tokenIndex1891
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1889
					}
				}
			l1891:
				add(rulejsonSetPathNonHead, position1890)
			}
			return true
		l1889:
			position, tokenIndex = position1889, tokenIndex1889
			return false
		},
		/* 171 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1893, tokenIndex1893 := position, tokenIndex
			{
				position1894 := position
				{
					position1895, tokenIndex1895 := position, tokenIndex
					if buffer[position] != rune('.') {
						goto l1896
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1896
					}
					goto l1895
				l1896:
					position, tokenIndex = position1895, tokenIndex1895
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1893
					}
				}
			l1895:
				add(rulejsonMapSingleLevel, position1894)
			}
			return true
		l1893:
			position, tokenIndex = position1893, tokenIndex1893
			return false
		},
		/* 172 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1897, tokenIndex1897 := position, tokenIndex
			{
				position1898 := position
				if buffer[position] != rune('.') {
					goto l1897
				}
				position++
				if buffer[position] != rune('.') {
					goto l1897
				}
				position++
				{
					position1899, tokenIndex1899 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l1900
					}
					goto l1899
				l1900:
					position, tokenIndex = position1899, tokenIndex1899
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1897
					}
				}
			l1899:
				add(rulejsonMapMultipleLevel, position1898)
			}
			return true
		l1897:
			position, tokenIndex = position1897, tokenIndex1897
			return false
		},
		/* 173 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1901, tokenIndex1901 := position, tokenIndex
			{
				position1902 := position
				{
					position1903 := position
					{
						position1904, tokenIndex1904 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1905
						}
						position++
						goto l1904
					l1905:
						position, tokenIndex = position1904, tokenIndex1904
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1901
						}
						position++
					}
				l1904:
				l1906:
					{
						position1907, tokenIndex1907 := position, tokenIndex
						{
							position1908, tokenIndex1908 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1909
							}
							position++
							goto l1908
						l1909:
							position, tokenIndex = position1908, tokenIndex1908
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1910
							}
							position++
							goto l1908
						l1910:
							position, tokenIndex = position1908, tokenIndex1908
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1911
							}
							position++
							goto l1908
						l1911:
							position, tokenIndex = position1908, tokenIndex1908
							if buffer[position] != rune('_') {
								goto l1907
							}
							position++
						}
					l1908:
						goto l1906
					l1907:
						position, tokenIndex = position1907, tokenIndex1907
					}
					add(rulePegText, position1903)
				}
				add(rulejsonMapAccessString, position1902)
			}
			return true
		l1901:
			position, tokenIndex = position1901, tokenIndex1901
			return false
		},
		/* 174 jsonMapAccessBracket <- <('[' doubleQuotedString ']')> */
		func() bool {
			position1912, tokenIndex1912 := position, tokenIndex
			{
				position1913 := position
				if buffer[position] != rune('[') {
					goto l1912
				}
				position++
				if !_rules[ruledoubleQuotedString]() {
					goto l1912
				}
				if buffer[position] != rune(']') {
					goto l1912
				}
				position++
				add(rulejsonMapAccessBracket, position1913)
			}
			return true
		l1912:
			position, tokenIndex = position1912, tokenIndex1912
			return false
		},
		/* 175 doubleQuotedString <- <('"' <(('"' '"') / (!'"' .))*> '"')> */
		func() bool {
			position1914, tokenIndex1914 := position, tokenIndex
			{
				position1915 := position
				if buffer[position] != rune('"') {
					goto l1914
				}
				position++
				{
					position1916 := position
				l1917:
					{
						position1918, tokenIndex1918 := position, tokenIndex
						{
							position1919, tokenIndex1919 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l1920
							}
							position++
							if buffer[position] != rune('"') {
								goto l1920
							}
							position++
							goto l1919
						l1920:
							position, tokenIndex = position1919, tokenIndex1919
							{
								position1921, tokenIndex1921 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l1921
								}
								position++
								goto l1918
							l1921:
								position, tokenIndex = position1921, tokenIndex1921
							}
							if !matchDot() {
								goto l1918
							}
						}
					l1919:
						goto l1917
					l1918:
						position, tokenIndex = position1918, tokenIndex1918
					}
					add(rulePegText, position1916)
				}
				if buffer[position] != rune('"') {
					goto l1914
				}
				position++
				add(ruledoubleQuotedString, position1915)
			}
			return true
		l1914:
			position, tokenIndex = position1914, tokenIndex1914
			return false
		},
		/* 176 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1922, tokenIndex1922 := position, tokenIndex
			{
				position1923 := position
				if buffer[position] != rune('[') {
					goto l1922
				}
				position++
				{
					position1924 := position
					{
						position1925, tokenIndex1925 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1925
						}
						position++
						goto l1926
					l1925:
						position, tokenIndex = position1925, tokenIndex1925
					}
				l1926:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1922
					}
					position++
				l1927:
					{
						position1928, tokenIndex1928 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1928
						}
						position++
						goto l1927
					l1928:
						position, tokenIndex = position1928, tokenIndex1928
					}
					add(rulePegText, position1924)
				}
				if buffer[position] != rune(']') {
					goto l1922
				}
				position++
				add(rulejsonArrayAccess, position1923)
			}
			return true
		l1922:
			position, tokenIndex = position1922, tokenIndex1922
			return false
		},
		/* 177 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
		func() bool {
			position1929, tokenIndex1929 := position, tokenIndex
			{
				position1930 := position
				if buffer[position] != rune('[') {
					goto l1929
				}
				position++
				{
					position1931 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1929
					}
					position++
				l1932:
					{
						position1933, tokenIndex1933 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1933
						}
						position++
						goto l1932
					l1933:
						position, tokenIndex = position1933, tokenIndex1933
					}
					add(rulePegText, position1931)
				}
				if buffer[position] != rune(']') {
					goto l1929
				}
				position++
				add(rulejsonNonNegativeArrayAccess, position1930)
			}
			return true
		l1929:
			position, tokenIndex = position1929, tokenIndex1929
			return false
		},
		/* 178 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1934, tokenIndex1934 := position, tokenIndex
			{
				position1935 := position
				if buffer[position] != rune('[') {
					goto l1934
				}
				position++
				{
					position1936 := position
					{
						position1937, tokenIndex1937 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1937
						}
						position++
						goto l1938
					l1937:
						position, tokenIndex = position1937, tokenIndex1937
					}
				l1938:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1934
					}
					position++
				l1939:
					{
						position1940, tokenIndex1940 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1940
						}
						position++
						goto l1939
					l1940:
						position, tokenIndex = position1940, tokenIndex1940
					}
					if buffer[position] != rune(':') {
						goto l1934
					}
					position++
					{
						position1941, tokenIndex1941 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l1941
						}
						position++
						goto l1942
					l1941:
						position, tokenIndex = position1941, tokenIndex1941
					}
				l1942:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1934
					}
					position++
				l1943:
					{
						position1944, tokenIndex1944 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1944
						}
						position++
						goto l1943
					l1944:
						position, tokenIndex = position1944, tokenIndex1944
					}
					{
						position1945, tokenIndex1945 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l1945
						}
						position++
						{
							position1947, tokenIndex1947 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1947
							}
							position++
							goto l1948
						l1947:
							position, tokenIndex = position1947, tokenIndex1947
						}
					l1948:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1945
						}
						position++
					l1949:
						{
							position1950, tokenIndex1950 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1950
							}
							position++
							goto l1949
						l1950:
							position, tokenIndex = position1950, tokenIndex1950
						}
						goto l1946
					l1945:
						position, tokenIndex = position1945, tokenIndex1945
					}
				l1946:
					add(rulePegText, position1936)
				}
				if buffer[position] != rune(']') {
					goto l1934
				}
				position++
				add(rulejsonArraySlice, position1935)
			}
			return true
		l1934:
			position, tokenIndex = position1934, tokenIndex1934
			return false
		},
		/* 179 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1951, tokenIndex1951 := position, tokenIndex
			{
				position1952 := position
				if buffer[position] != rune('[') {
					goto l1951
				}
				position++
				{
					position1953 := position
					{
						position1954, tokenIndex1954 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l1955
						}
						position++
						{
							position1956, tokenIndex1956 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1956
							}
							position++
							goto l1957
						l1956:
							position, tokenIndex = position1956, tokenIndex1956
						}
					l1957:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1955
						}
						position++
					l1958:
						{
							position1959, tokenIndex1959 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1959
							}
							position++
							goto l1958
						l1959:
							position, tokenIndex = position1959, tokenIndex1959
						}
						goto l1954
					l1955:
						position, tokenIndex = position1954, tokenIndex1954
						{
							position1960, tokenIndex1960 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l1960
							}
							position++
							goto l1961
						l1960:
							position, tokenIndex = position1960, tokenIndex1960
						}
					l1961:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1951
						}
						position++
					l1962:
						{
							position1963, tokenIndex1963 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1963
							}
							position++
							goto l1962
						l1963:
							position, tokenIndex = position1963, tokenIndex1963
						}
						if buffer[position] != rune(':') {
							goto l1951
						}
						position++
					}
				l1954:
					add(rulePegText, position1953)
				}
				if buffer[position] != rune(']') {
					goto l1951
				}
				position++
				add(rulejsonArrayPartialSlice, position1952)
			}
			return true
		l1951:
			position, tokenIndex = position1951, tokenIndex1951
			return false
		},
		/* 180 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1964, tokenIndex1964 := position, tokenIndex
			{
				position1965 := position
				if buffer[position] != rune('[') {
					goto l1964
				}
				position++
				if buffer[position] != rune(':') {
					goto l1964
				}
				position++
				if buffer[position] != rune(']') {
					goto l1964
				}
				position++
				add(rulejsonArrayFullSlice, position1965)
			}
			return true
		l1964:
			position, tokenIndex = position1964, tokenIndex1964
			return false
		},
		/* 181 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1966, tokenIndex1966 := position, tokenIndex
			{
				position1967 := position
				{
					position1968, tokenIndex1968 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l1969
					}
					position++
					goto l1968
				l1969:
					position, tokenIndex = position1968, tokenIndex1968
					if buffer[position] != rune('\t') {
						goto l1970
					}
					position++
					goto l1968
				l1970:
					position, tokenIndex = position1968, tokenIndex1968
					if buffer[position] != rune('\n') {
						goto l1971
					}
					position++
					goto l1968
				l1971:
					position, tokenIndex = position1968, tokenIndex1968
					if buffer[position] != rune('\r') {
						goto l1972
					}
					position++
					goto l1968
				l1972:
					position, tokenIndex = position1968, tokenIndex1968
					if !_rules[rulecomment]() {
						goto l1973
					}
					goto l1968
				l1973:
					position, tokenIndex = position1968, tokenIndex1968
					if !_rules[rulefinalComment]() {
						goto l1966
					}
				}
			l1968:
				add(rulespElem, position1967)
			}
			return true
		l1966:
			position, tokenIndex = position1966, tokenIndex1966
			return false
		},
		/* 182 sp <- <spElem+> */
		func() bool {
			position1974, tokenIndex1974 := position, tokenIndex
			{
				position1975 := position
				if !_rules[rulespElem]() {
					goto l1974
				}
			l1976:
				{
					position1977, tokenIndex1977 := position, tokenIndex
					if !_rules[rulespElem]() {
						goto l1977
					}
					goto l1976
				l1977:
					position, tokenIndex = position1977, tokenIndex1977
				}
				add(rulesp, position1975)
			}
			return true
		l1974:
			position, tokenIndex = position1974, tokenIndex1974
			return false
		},
		/* 183 spOpt <- <spElem*> */
		func() bool {
			{
				position1979 := position
			l1980:
				{
					position1981, tokenIndex1981 := position, tokenIndex
					if !_rules[rulespElem]() {
						goto l1981
					}
					goto l1980
				l1981:
					position, tokenIndex = position1981, tokenIndex1981
				}
				add(rulespOpt, position1979)
			}
			return true
		},
		/* 184 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1982, tokenIndex1982 := position, tokenIndex
			{
				position1983 := position
				if buffer[position] != rune('-') {
					goto l1982
				}
				position++
				if buffer[position] != rune('-') {
					goto l1982
				}
				position++
			l1984:
				{
					position1985, tokenIndex1985 := position, tokenIndex
					{
						position1986, tokenIndex1986 := position, tokenIndex
						{
							position1987, tokenIndex1987 := position, tokenIndex
							if buffer[position] != rune('\r') {
								goto l1988
							}
							position++
							goto l1987
						l1988:
							position, tokenIndex = position1987, tokenIndex1987
							if buffer[position] != rune('\n') {
								goto l1986
							}
							position++
						}
					l1987:
						goto l1985
					l1986:
						position, tokenIndex = position1986, tokenIndex1986
					}
					if !matchDot() {
						goto l1985
					}
					goto l1984
				l1985:
					position, tokenIndex = position1985, tokenIndex1985
				}
				{
					position1989, tokenIndex1989 := position, tokenIndex
					if buffer[position] != rune('\r') {
						goto l1990
					}
					position++
					goto l1989
				l1990:
					position, tokenIndex = position1989, tokenIndex1989
					if buffer[position] != rune('\n') {
						goto l1982
					}
					position++
				}
			l1989:
				add(rulecomment, position1983)
			}
			return true
		l1982:
			position, tokenIndex = position1982, tokenIndex1982
			return false
		},
		/* 185 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1991, tokenIndex1991 := position, tokenIndex
			{
				position1992 := position
				if buffer[position] != rune('-') {
					goto l1991
				}
				position++
				if buffer[position] != rune('-') {
					goto l1991
				}
				position++
			l1993:
				{
					position1994, tokenIndex1994 := position, tokenIndex
					{
						position1995, tokenIndex1995 := position, tokenIndex
						{
							position1996, tokenIndex1996 := position, tokenIndex
							if buffer[position] != rune('\r') {
								goto l1997
							}
							position++
							goto l1996
						l1997:
							position, tokenIndex = position1996, tokenIndex1996
							if buffer[position] != rune('\n') {
								goto l1995
							}
							position++
						}
					l1996:
						goto l1994
					l1995:
						position, tokenIndex = position1995, tokenIndex1995
					}
					if !matchDot() {
						goto l1994
					}
					goto l1993
				l1994:
					position, tokenIndex = position1994, tokenIndex1994
				}
				{
					position1998, tokenIndex1998 := position, tokenIndex
					if !matchDot() {
						goto l1998
					}
					goto l1991
				l1998:
					position, tokenIndex = position1998, tokenIndex1998
				}
				add(rulefinalComment, position1992)
			}
			return true
		l1991:
			position, tokenIndex = position1991, tokenIndex1991
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

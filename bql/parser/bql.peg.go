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
	rules  [289]func() bool
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

			p.AssembleEmitterSampling(CountBasedSampling)

		case ruleAction28:

			p.AssembleEmitterSampling(RandomizedSampling)

		case ruleAction29:

			p.AssembleProjections(begin, end)

		case ruleAction30:

			p.AssembleAlias()

		case ruleAction31:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction32:

			p.AssembleInterval()

		case ruleAction33:

			p.AssembleInterval()

		case ruleAction34:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction35:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction36:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction37:

			p.EnsureAliasedStreamWindow()

		case ruleAction38:

			p.AssembleAliasedStreamWindow()

		case ruleAction39:

			p.AssembleStreamWindow()

		case ruleAction40:

			p.AssembleUDSFFuncApp()

		case ruleAction41:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction42:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction43:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction44:

			p.AssembleSourceSinkParam()

		case ruleAction45:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction46:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction47:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction48:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction49:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction50:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction51:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction52:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction53:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction54:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction55:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction56:

			p.AssembleTypeCast(begin, end)

		case ruleAction57:

			p.AssembleTypeCast(begin, end)

		case ruleAction58:

			p.AssembleFuncApp()

		case ruleAction59:

			p.AssembleExpressions(begin, end)
			p.AssembleFuncApp()

		case ruleAction60:

			p.AssembleExpressions(begin, end)

		case ruleAction61:

			p.AssembleExpressions(begin, end)

		case ruleAction62:

			p.AssembleSortedExpression()

		case ruleAction63:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction64:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction65:

			p.AssembleMap(begin, end)

		case ruleAction66:

			p.AssembleKeyValuePair()

		case ruleAction67:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction68:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction69:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction70:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction71:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction72:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction73:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction74:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction75:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction76:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction77:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction78:

			p.PushComponent(begin, end, Istream)

		case ruleAction79:

			p.PushComponent(begin, end, Dstream)

		case ruleAction80:

			p.PushComponent(begin, end, Rstream)

		case ruleAction81:

			p.PushComponent(begin, end, Tuples)

		case ruleAction82:

			p.PushComponent(begin, end, Seconds)

		case ruleAction83:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction84:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction85:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction86:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction87:

			p.PushComponent(begin, end, Yes)

		case ruleAction88:

			p.PushComponent(begin, end, No)

		case ruleAction89:

			p.PushComponent(begin, end, Yes)

		case ruleAction90:

			p.PushComponent(begin, end, No)

		case ruleAction91:

			p.PushComponent(begin, end, Bool)

		case ruleAction92:

			p.PushComponent(begin, end, Int)

		case ruleAction93:

			p.PushComponent(begin, end, Float)

		case ruleAction94:

			p.PushComponent(begin, end, String)

		case ruleAction95:

			p.PushComponent(begin, end, Blob)

		case ruleAction96:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction97:

			p.PushComponent(begin, end, Array)

		case ruleAction98:

			p.PushComponent(begin, end, Map)

		case ruleAction99:

			p.PushComponent(begin, end, Or)

		case ruleAction100:

			p.PushComponent(begin, end, And)

		case ruleAction101:

			p.PushComponent(begin, end, Not)

		case ruleAction102:

			p.PushComponent(begin, end, Equal)

		case ruleAction103:

			p.PushComponent(begin, end, Less)

		case ruleAction104:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction105:

			p.PushComponent(begin, end, Greater)

		case ruleAction106:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction107:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction108:

			p.PushComponent(begin, end, Concat)

		case ruleAction109:

			p.PushComponent(begin, end, Is)

		case ruleAction110:

			p.PushComponent(begin, end, IsNot)

		case ruleAction111:

			p.PushComponent(begin, end, Plus)

		case ruleAction112:

			p.PushComponent(begin, end, Minus)

		case ruleAction113:

			p.PushComponent(begin, end, Multiply)

		case ruleAction114:

			p.PushComponent(begin, end, Divide)

		case ruleAction115:

			p.PushComponent(begin, end, Modulo)

		case ruleAction116:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction117:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction118:

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
		/* 34 EmitterSample <- <(CountBasedSampling / RandomizedSampling)> */
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
			position657, tokenIndex657, depth657 := position, tokenIndex, depth
			{
				position658 := position
				depth++
				{
					position659, tokenIndex659, depth659 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l660
					}
					position++
					goto l659
				l660:
					position, tokenIndex, depth = position659, tokenIndex659, depth659
					if buffer[position] != rune('E') {
						goto l657
					}
					position++
				}
			l659:
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l662
					}
					position++
					goto l661
				l662:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
					if buffer[position] != rune('V') {
						goto l657
					}
					position++
				}
			l661:
				{
					position663, tokenIndex663, depth663 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l664
					}
					position++
					goto l663
				l664:
					position, tokenIndex, depth = position663, tokenIndex663, depth663
					if buffer[position] != rune('E') {
						goto l657
					}
					position++
				}
			l663:
				{
					position665, tokenIndex665, depth665 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l666
					}
					position++
					goto l665
				l666:
					position, tokenIndex, depth = position665, tokenIndex665, depth665
					if buffer[position] != rune('R') {
						goto l657
					}
					position++
				}
			l665:
				{
					position667, tokenIndex667, depth667 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l668
					}
					position++
					goto l667
				l668:
					position, tokenIndex, depth = position667, tokenIndex667, depth667
					if buffer[position] != rune('Y') {
						goto l657
					}
					position++
				}
			l667:
				if !_rules[rulesp]() {
					goto l657
				}
				if !_rules[ruleNumericLiteral]() {
					goto l657
				}
				if !_rules[rulespOpt]() {
					goto l657
				}
				{
					position669, tokenIndex669, depth669 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l669
					}
					position++
					goto l670
				l669:
					position, tokenIndex, depth = position669, tokenIndex669, depth669
				}
			l670:
				if !_rules[rulespOpt]() {
					goto l657
				}
				{
					position671, tokenIndex671, depth671 := position, tokenIndex, depth
					{
						position673, tokenIndex673, depth673 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l674
						}
						position++
						goto l673
					l674:
						position, tokenIndex, depth = position673, tokenIndex673, depth673
						if buffer[position] != rune('S') {
							goto l672
						}
						position++
					}
				l673:
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l676
						}
						position++
						goto l675
					l676:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('T') {
							goto l672
						}
						position++
					}
				l675:
					goto l671
				l672:
					position, tokenIndex, depth = position671, tokenIndex671, depth671
					{
						position678, tokenIndex678, depth678 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l679
						}
						position++
						goto l678
					l679:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
						if buffer[position] != rune('N') {
							goto l677
						}
						position++
					}
				l678:
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l681
						}
						position++
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if buffer[position] != rune('D') {
							goto l677
						}
						position++
					}
				l680:
					goto l671
				l677:
					position, tokenIndex, depth = position671, tokenIndex671, depth671
					{
						position683, tokenIndex683, depth683 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l684
						}
						position++
						goto l683
					l684:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
						if buffer[position] != rune('R') {
							goto l682
						}
						position++
					}
				l683:
					{
						position685, tokenIndex685, depth685 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l686
						}
						position++
						goto l685
					l686:
						position, tokenIndex, depth = position685, tokenIndex685, depth685
						if buffer[position] != rune('D') {
							goto l682
						}
						position++
					}
				l685:
					goto l671
				l682:
					position, tokenIndex, depth = position671, tokenIndex671, depth671
					{
						position687, tokenIndex687, depth687 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l688
						}
						position++
						goto l687
					l688:
						position, tokenIndex, depth = position687, tokenIndex687, depth687
						if buffer[position] != rune('T') {
							goto l657
						}
						position++
					}
				l687:
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if buffer[position] != rune('h') {
							goto l690
						}
						position++
						goto l689
					l690:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
						if buffer[position] != rune('H') {
							goto l657
						}
						position++
					}
				l689:
				}
			l671:
				if !_rules[rulesp]() {
					goto l657
				}
				{
					position691, tokenIndex691, depth691 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l692
					}
					position++
					goto l691
				l692:
					position, tokenIndex, depth = position691, tokenIndex691, depth691
					if buffer[position] != rune('T') {
						goto l657
					}
					position++
				}
			l691:
				{
					position693, tokenIndex693, depth693 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l694
					}
					position++
					goto l693
				l694:
					position, tokenIndex, depth = position693, tokenIndex693, depth693
					if buffer[position] != rune('U') {
						goto l657
					}
					position++
				}
			l693:
				{
					position695, tokenIndex695, depth695 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l696
					}
					position++
					goto l695
				l696:
					position, tokenIndex, depth = position695, tokenIndex695, depth695
					if buffer[position] != rune('P') {
						goto l657
					}
					position++
				}
			l695:
				{
					position697, tokenIndex697, depth697 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l698
					}
					position++
					goto l697
				l698:
					position, tokenIndex, depth = position697, tokenIndex697, depth697
					if buffer[position] != rune('L') {
						goto l657
					}
					position++
				}
			l697:
				{
					position699, tokenIndex699, depth699 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l700
					}
					position++
					goto l699
				l700:
					position, tokenIndex, depth = position699, tokenIndex699, depth699
					if buffer[position] != rune('E') {
						goto l657
					}
					position++
				}
			l699:
				if !_rules[ruleAction27]() {
					goto l657
				}
				depth--
				add(ruleCountBasedSampling, position658)
			}
			return true
		l657:
			position, tokenIndex, depth = position657, tokenIndex657, depth657
			return false
		},
		/* 36 RandomizedSampling <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') sp NumericLiteral spOpt '%' Action28)> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				{
					position703, tokenIndex703, depth703 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l704
					}
					position++
					goto l703
				l704:
					position, tokenIndex, depth = position703, tokenIndex703, depth703
					if buffer[position] != rune('S') {
						goto l701
					}
					position++
				}
			l703:
				{
					position705, tokenIndex705, depth705 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l706
					}
					position++
					goto l705
				l706:
					position, tokenIndex, depth = position705, tokenIndex705, depth705
					if buffer[position] != rune('A') {
						goto l701
					}
					position++
				}
			l705:
				{
					position707, tokenIndex707, depth707 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l708
					}
					position++
					goto l707
				l708:
					position, tokenIndex, depth = position707, tokenIndex707, depth707
					if buffer[position] != rune('M') {
						goto l701
					}
					position++
				}
			l707:
				{
					position709, tokenIndex709, depth709 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l710
					}
					position++
					goto l709
				l710:
					position, tokenIndex, depth = position709, tokenIndex709, depth709
					if buffer[position] != rune('P') {
						goto l701
					}
					position++
				}
			l709:
				{
					position711, tokenIndex711, depth711 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l712
					}
					position++
					goto l711
				l712:
					position, tokenIndex, depth = position711, tokenIndex711, depth711
					if buffer[position] != rune('L') {
						goto l701
					}
					position++
				}
			l711:
				{
					position713, tokenIndex713, depth713 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l714
					}
					position++
					goto l713
				l714:
					position, tokenIndex, depth = position713, tokenIndex713, depth713
					if buffer[position] != rune('E') {
						goto l701
					}
					position++
				}
			l713:
				if !_rules[rulesp]() {
					goto l701
				}
				if !_rules[ruleNumericLiteral]() {
					goto l701
				}
				if !_rules[rulespOpt]() {
					goto l701
				}
				if buffer[position] != rune('%') {
					goto l701
				}
				position++
				if !_rules[ruleAction28]() {
					goto l701
				}
				depth--
				add(ruleRandomizedSampling, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 37 Projections <- <(<(sp Projection (spOpt ',' spOpt Projection)*)> Action29)> */
		func() bool {
			position715, tokenIndex715, depth715 := position, tokenIndex, depth
			{
				position716 := position
				depth++
				{
					position717 := position
					depth++
					if !_rules[rulesp]() {
						goto l715
					}
					if !_rules[ruleProjection]() {
						goto l715
					}
				l718:
					{
						position719, tokenIndex719, depth719 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l719
						}
						if buffer[position] != rune(',') {
							goto l719
						}
						position++
						if !_rules[rulespOpt]() {
							goto l719
						}
						if !_rules[ruleProjection]() {
							goto l719
						}
						goto l718
					l719:
						position, tokenIndex, depth = position719, tokenIndex719, depth719
					}
					depth--
					add(rulePegText, position717)
				}
				if !_rules[ruleAction29]() {
					goto l715
				}
				depth--
				add(ruleProjections, position716)
			}
			return true
		l715:
			position, tokenIndex, depth = position715, tokenIndex715, depth715
			return false
		},
		/* 38 Projection <- <(AliasExpression / ExpressionOrWildcard)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722, tokenIndex722, depth722 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l723
					}
					goto l722
				l723:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleExpressionOrWildcard]() {
						goto l720
					}
				}
			l722:
				depth--
				add(ruleProjection, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 39 AliasExpression <- <(ExpressionOrWildcard sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action30)> */
		func() bool {
			position724, tokenIndex724, depth724 := position, tokenIndex, depth
			{
				position725 := position
				depth++
				if !_rules[ruleExpressionOrWildcard]() {
					goto l724
				}
				if !_rules[rulesp]() {
					goto l724
				}
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
						goto l724
					}
					position++
				}
			l726:
				{
					position728, tokenIndex728, depth728 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l729
					}
					position++
					goto l728
				l729:
					position, tokenIndex, depth = position728, tokenIndex728, depth728
					if buffer[position] != rune('S') {
						goto l724
					}
					position++
				}
			l728:
				if !_rules[rulesp]() {
					goto l724
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l724
				}
				if !_rules[ruleAction30]() {
					goto l724
				}
				depth--
				add(ruleAliasExpression, position725)
			}
			return true
		l724:
			position, tokenIndex, depth = position724, tokenIndex724, depth724
			return false
		},
		/* 40 WindowedFrom <- <(<(sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp Relations)?> Action31)> */
		func() bool {
			position730, tokenIndex730, depth730 := position, tokenIndex, depth
			{
				position731 := position
				depth++
				{
					position732 := position
					depth++
					{
						position733, tokenIndex733, depth733 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l733
						}
						{
							position735, tokenIndex735, depth735 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l736
							}
							position++
							goto l735
						l736:
							position, tokenIndex, depth = position735, tokenIndex735, depth735
							if buffer[position] != rune('F') {
								goto l733
							}
							position++
						}
					l735:
						{
							position737, tokenIndex737, depth737 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l738
							}
							position++
							goto l737
						l738:
							position, tokenIndex, depth = position737, tokenIndex737, depth737
							if buffer[position] != rune('R') {
								goto l733
							}
							position++
						}
					l737:
						{
							position739, tokenIndex739, depth739 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l740
							}
							position++
							goto l739
						l740:
							position, tokenIndex, depth = position739, tokenIndex739, depth739
							if buffer[position] != rune('O') {
								goto l733
							}
							position++
						}
					l739:
						{
							position741, tokenIndex741, depth741 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l742
							}
							position++
							goto l741
						l742:
							position, tokenIndex, depth = position741, tokenIndex741, depth741
							if buffer[position] != rune('M') {
								goto l733
							}
							position++
						}
					l741:
						if !_rules[rulesp]() {
							goto l733
						}
						if !_rules[ruleRelations]() {
							goto l733
						}
						goto l734
					l733:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
					}
				l734:
					depth--
					add(rulePegText, position732)
				}
				if !_rules[ruleAction31]() {
					goto l730
				}
				depth--
				add(ruleWindowedFrom, position731)
			}
			return true
		l730:
			position, tokenIndex, depth = position730, tokenIndex730, depth730
			return false
		},
		/* 41 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position743, tokenIndex743, depth743 := position, tokenIndex, depth
			{
				position744 := position
				depth++
				{
					position745, tokenIndex745, depth745 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l746
					}
					goto l745
				l746:
					position, tokenIndex, depth = position745, tokenIndex745, depth745
					if !_rules[ruleTuplesInterval]() {
						goto l743
					}
				}
			l745:
				depth--
				add(ruleInterval, position744)
			}
			return true
		l743:
			position, tokenIndex, depth = position743, tokenIndex743, depth743
			return false
		},
		/* 42 TimeInterval <- <(NumericLiteral sp (SECONDS / MILLISECONDS) Action32)> */
		func() bool {
			position747, tokenIndex747, depth747 := position, tokenIndex, depth
			{
				position748 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l747
				}
				if !_rules[rulesp]() {
					goto l747
				}
				{
					position749, tokenIndex749, depth749 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l750
					}
					goto l749
				l750:
					position, tokenIndex, depth = position749, tokenIndex749, depth749
					if !_rules[ruleMILLISECONDS]() {
						goto l747
					}
				}
			l749:
				if !_rules[ruleAction32]() {
					goto l747
				}
				depth--
				add(ruleTimeInterval, position748)
			}
			return true
		l747:
			position, tokenIndex, depth = position747, tokenIndex747, depth747
			return false
		},
		/* 43 TuplesInterval <- <(NumericLiteral sp TUPLES Action33)> */
		func() bool {
			position751, tokenIndex751, depth751 := position, tokenIndex, depth
			{
				position752 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l751
				}
				if !_rules[rulesp]() {
					goto l751
				}
				if !_rules[ruleTUPLES]() {
					goto l751
				}
				if !_rules[ruleAction33]() {
					goto l751
				}
				depth--
				add(ruleTuplesInterval, position752)
			}
			return true
		l751:
			position, tokenIndex, depth = position751, tokenIndex751, depth751
			return false
		},
		/* 44 Relations <- <(RelationLike (spOpt ',' spOpt RelationLike)*)> */
		func() bool {
			position753, tokenIndex753, depth753 := position, tokenIndex, depth
			{
				position754 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l753
				}
			l755:
				{
					position756, tokenIndex756, depth756 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l756
					}
					if buffer[position] != rune(',') {
						goto l756
					}
					position++
					if !_rules[rulespOpt]() {
						goto l756
					}
					if !_rules[ruleRelationLike]() {
						goto l756
					}
					goto l755
				l756:
					position, tokenIndex, depth = position756, tokenIndex756, depth756
				}
				depth--
				add(ruleRelations, position754)
			}
			return true
		l753:
			position, tokenIndex, depth = position753, tokenIndex753, depth753
			return false
		},
		/* 45 Filter <- <(<(sp (('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E')) sp Expression)?> Action34)> */
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
						if !_rules[rulesp]() {
							goto l760
						}
						{
							position762, tokenIndex762, depth762 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l763
							}
							position++
							goto l762
						l763:
							position, tokenIndex, depth = position762, tokenIndex762, depth762
							if buffer[position] != rune('W') {
								goto l760
							}
							position++
						}
					l762:
						{
							position764, tokenIndex764, depth764 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l765
							}
							position++
							goto l764
						l765:
							position, tokenIndex, depth = position764, tokenIndex764, depth764
							if buffer[position] != rune('H') {
								goto l760
							}
							position++
						}
					l764:
						{
							position766, tokenIndex766, depth766 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l767
							}
							position++
							goto l766
						l767:
							position, tokenIndex, depth = position766, tokenIndex766, depth766
							if buffer[position] != rune('E') {
								goto l760
							}
							position++
						}
					l766:
						{
							position768, tokenIndex768, depth768 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l769
							}
							position++
							goto l768
						l769:
							position, tokenIndex, depth = position768, tokenIndex768, depth768
							if buffer[position] != rune('R') {
								goto l760
							}
							position++
						}
					l768:
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
								goto l760
							}
							position++
						}
					l770:
						if !_rules[rulesp]() {
							goto l760
						}
						if !_rules[ruleExpression]() {
							goto l760
						}
						goto l761
					l760:
						position, tokenIndex, depth = position760, tokenIndex760, depth760
					}
				l761:
					depth--
					add(rulePegText, position759)
				}
				if !_rules[ruleAction34]() {
					goto l757
				}
				depth--
				add(ruleFilter, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 46 Grouping <- <(<(sp (('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P')) sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action35)> */
		func() bool {
			position772, tokenIndex772, depth772 := position, tokenIndex, depth
			{
				position773 := position
				depth++
				{
					position774 := position
					depth++
					{
						position775, tokenIndex775, depth775 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l775
						}
						{
							position777, tokenIndex777, depth777 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l778
							}
							position++
							goto l777
						l778:
							position, tokenIndex, depth = position777, tokenIndex777, depth777
							if buffer[position] != rune('G') {
								goto l775
							}
							position++
						}
					l777:
						{
							position779, tokenIndex779, depth779 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l780
							}
							position++
							goto l779
						l780:
							position, tokenIndex, depth = position779, tokenIndex779, depth779
							if buffer[position] != rune('R') {
								goto l775
							}
							position++
						}
					l779:
						{
							position781, tokenIndex781, depth781 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l782
							}
							position++
							goto l781
						l782:
							position, tokenIndex, depth = position781, tokenIndex781, depth781
							if buffer[position] != rune('O') {
								goto l775
							}
							position++
						}
					l781:
						{
							position783, tokenIndex783, depth783 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l784
							}
							position++
							goto l783
						l784:
							position, tokenIndex, depth = position783, tokenIndex783, depth783
							if buffer[position] != rune('U') {
								goto l775
							}
							position++
						}
					l783:
						{
							position785, tokenIndex785, depth785 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l786
							}
							position++
							goto l785
						l786:
							position, tokenIndex, depth = position785, tokenIndex785, depth785
							if buffer[position] != rune('P') {
								goto l775
							}
							position++
						}
					l785:
						if !_rules[rulesp]() {
							goto l775
						}
						{
							position787, tokenIndex787, depth787 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l788
							}
							position++
							goto l787
						l788:
							position, tokenIndex, depth = position787, tokenIndex787, depth787
							if buffer[position] != rune('B') {
								goto l775
							}
							position++
						}
					l787:
						{
							position789, tokenIndex789, depth789 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l790
							}
							position++
							goto l789
						l790:
							position, tokenIndex, depth = position789, tokenIndex789, depth789
							if buffer[position] != rune('Y') {
								goto l775
							}
							position++
						}
					l789:
						if !_rules[rulesp]() {
							goto l775
						}
						if !_rules[ruleGroupList]() {
							goto l775
						}
						goto l776
					l775:
						position, tokenIndex, depth = position775, tokenIndex775, depth775
					}
				l776:
					depth--
					add(rulePegText, position774)
				}
				if !_rules[ruleAction35]() {
					goto l772
				}
				depth--
				add(ruleGrouping, position773)
			}
			return true
		l772:
			position, tokenIndex, depth = position772, tokenIndex772, depth772
			return false
		},
		/* 47 GroupList <- <(Expression (spOpt ',' spOpt Expression)*)> */
		func() bool {
			position791, tokenIndex791, depth791 := position, tokenIndex, depth
			{
				position792 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l791
				}
			l793:
				{
					position794, tokenIndex794, depth794 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l794
					}
					if buffer[position] != rune(',') {
						goto l794
					}
					position++
					if !_rules[rulespOpt]() {
						goto l794
					}
					if !_rules[ruleExpression]() {
						goto l794
					}
					goto l793
				l794:
					position, tokenIndex, depth = position794, tokenIndex794, depth794
				}
				depth--
				add(ruleGroupList, position792)
			}
			return true
		l791:
			position, tokenIndex, depth = position791, tokenIndex791, depth791
			return false
		},
		/* 48 Having <- <(<(sp (('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G')) sp Expression)?> Action36)> */
		func() bool {
			position795, tokenIndex795, depth795 := position, tokenIndex, depth
			{
				position796 := position
				depth++
				{
					position797 := position
					depth++
					{
						position798, tokenIndex798, depth798 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l798
						}
						{
							position800, tokenIndex800, depth800 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l801
							}
							position++
							goto l800
						l801:
							position, tokenIndex, depth = position800, tokenIndex800, depth800
							if buffer[position] != rune('H') {
								goto l798
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
								goto l798
							}
							position++
						}
					l802:
						{
							position804, tokenIndex804, depth804 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l805
							}
							position++
							goto l804
						l805:
							position, tokenIndex, depth = position804, tokenIndex804, depth804
							if buffer[position] != rune('V') {
								goto l798
							}
							position++
						}
					l804:
						{
							position806, tokenIndex806, depth806 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l807
							}
							position++
							goto l806
						l807:
							position, tokenIndex, depth = position806, tokenIndex806, depth806
							if buffer[position] != rune('I') {
								goto l798
							}
							position++
						}
					l806:
						{
							position808, tokenIndex808, depth808 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l809
							}
							position++
							goto l808
						l809:
							position, tokenIndex, depth = position808, tokenIndex808, depth808
							if buffer[position] != rune('N') {
								goto l798
							}
							position++
						}
					l808:
						{
							position810, tokenIndex810, depth810 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l811
							}
							position++
							goto l810
						l811:
							position, tokenIndex, depth = position810, tokenIndex810, depth810
							if buffer[position] != rune('G') {
								goto l798
							}
							position++
						}
					l810:
						if !_rules[rulesp]() {
							goto l798
						}
						if !_rules[ruleExpression]() {
							goto l798
						}
						goto l799
					l798:
						position, tokenIndex, depth = position798, tokenIndex798, depth798
					}
				l799:
					depth--
					add(rulePegText, position797)
				}
				if !_rules[ruleAction36]() {
					goto l795
				}
				depth--
				add(ruleHaving, position796)
			}
			return true
		l795:
			position, tokenIndex, depth = position795, tokenIndex795, depth795
			return false
		},
		/* 49 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action37))> */
		func() bool {
			position812, tokenIndex812, depth812 := position, tokenIndex, depth
			{
				position813 := position
				depth++
				{
					position814, tokenIndex814, depth814 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l815
					}
					goto l814
				l815:
					position, tokenIndex, depth = position814, tokenIndex814, depth814
					if !_rules[ruleStreamWindow]() {
						goto l812
					}
					if !_rules[ruleAction37]() {
						goto l812
					}
				}
			l814:
				depth--
				add(ruleRelationLike, position813)
			}
			return true
		l812:
			position, tokenIndex, depth = position812, tokenIndex812, depth812
			return false
		},
		/* 50 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action38)> */
		func() bool {
			position816, tokenIndex816, depth816 := position, tokenIndex, depth
			{
				position817 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l816
				}
				if !_rules[rulesp]() {
					goto l816
				}
				{
					position818, tokenIndex818, depth818 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l819
					}
					position++
					goto l818
				l819:
					position, tokenIndex, depth = position818, tokenIndex818, depth818
					if buffer[position] != rune('A') {
						goto l816
					}
					position++
				}
			l818:
				{
					position820, tokenIndex820, depth820 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l821
					}
					position++
					goto l820
				l821:
					position, tokenIndex, depth = position820, tokenIndex820, depth820
					if buffer[position] != rune('S') {
						goto l816
					}
					position++
				}
			l820:
				if !_rules[rulesp]() {
					goto l816
				}
				if !_rules[ruleIdentifier]() {
					goto l816
				}
				if !_rules[ruleAction38]() {
					goto l816
				}
				depth--
				add(ruleAliasedStreamWindow, position817)
			}
			return true
		l816:
			position, tokenIndex, depth = position816, tokenIndex816, depth816
			return false
		},
		/* 51 StreamWindow <- <(StreamLike spOpt '[' spOpt (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval spOpt ']' Action39)> */
		func() bool {
			position822, tokenIndex822, depth822 := position, tokenIndex, depth
			{
				position823 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l822
				}
				if !_rules[rulespOpt]() {
					goto l822
				}
				if buffer[position] != rune('[') {
					goto l822
				}
				position++
				if !_rules[rulespOpt]() {
					goto l822
				}
				{
					position824, tokenIndex824, depth824 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l825
					}
					position++
					goto l824
				l825:
					position, tokenIndex, depth = position824, tokenIndex824, depth824
					if buffer[position] != rune('R') {
						goto l822
					}
					position++
				}
			l824:
				{
					position826, tokenIndex826, depth826 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l827
					}
					position++
					goto l826
				l827:
					position, tokenIndex, depth = position826, tokenIndex826, depth826
					if buffer[position] != rune('A') {
						goto l822
					}
					position++
				}
			l826:
				{
					position828, tokenIndex828, depth828 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l829
					}
					position++
					goto l828
				l829:
					position, tokenIndex, depth = position828, tokenIndex828, depth828
					if buffer[position] != rune('N') {
						goto l822
					}
					position++
				}
			l828:
				{
					position830, tokenIndex830, depth830 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l831
					}
					position++
					goto l830
				l831:
					position, tokenIndex, depth = position830, tokenIndex830, depth830
					if buffer[position] != rune('G') {
						goto l822
					}
					position++
				}
			l830:
				{
					position832, tokenIndex832, depth832 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l833
					}
					position++
					goto l832
				l833:
					position, tokenIndex, depth = position832, tokenIndex832, depth832
					if buffer[position] != rune('E') {
						goto l822
					}
					position++
				}
			l832:
				if !_rules[rulesp]() {
					goto l822
				}
				if !_rules[ruleInterval]() {
					goto l822
				}
				if !_rules[rulespOpt]() {
					goto l822
				}
				if buffer[position] != rune(']') {
					goto l822
				}
				position++
				if !_rules[ruleAction39]() {
					goto l822
				}
				depth--
				add(ruleStreamWindow, position823)
			}
			return true
		l822:
			position, tokenIndex, depth = position822, tokenIndex822, depth822
			return false
		},
		/* 52 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position834, tokenIndex834, depth834 := position, tokenIndex, depth
			{
				position835 := position
				depth++
				{
					position836, tokenIndex836, depth836 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l837
					}
					goto l836
				l837:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleStream]() {
						goto l834
					}
				}
			l836:
				depth--
				add(ruleStreamLike, position835)
			}
			return true
		l834:
			position, tokenIndex, depth = position834, tokenIndex834, depth834
			return false
		},
		/* 53 UDSFFuncApp <- <(FuncAppWithoutOrderBy Action40)> */
		func() bool {
			position838, tokenIndex838, depth838 := position, tokenIndex, depth
			{
				position839 := position
				depth++
				if !_rules[ruleFuncAppWithoutOrderBy]() {
					goto l838
				}
				if !_rules[ruleAction40]() {
					goto l838
				}
				depth--
				add(ruleUDSFFuncApp, position839)
			}
			return true
		l838:
			position, tokenIndex, depth = position838, tokenIndex838, depth838
			return false
		},
		/* 54 SourceSinkSpecs <- <(<(sp (('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action41)> */
		func() bool {
			position840, tokenIndex840, depth840 := position, tokenIndex, depth
			{
				position841 := position
				depth++
				{
					position842 := position
					depth++
					{
						position843, tokenIndex843, depth843 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l843
						}
						{
							position845, tokenIndex845, depth845 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l846
							}
							position++
							goto l845
						l846:
							position, tokenIndex, depth = position845, tokenIndex845, depth845
							if buffer[position] != rune('W') {
								goto l843
							}
							position++
						}
					l845:
						{
							position847, tokenIndex847, depth847 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l848
							}
							position++
							goto l847
						l848:
							position, tokenIndex, depth = position847, tokenIndex847, depth847
							if buffer[position] != rune('I') {
								goto l843
							}
							position++
						}
					l847:
						{
							position849, tokenIndex849, depth849 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l850
							}
							position++
							goto l849
						l850:
							position, tokenIndex, depth = position849, tokenIndex849, depth849
							if buffer[position] != rune('T') {
								goto l843
							}
							position++
						}
					l849:
						{
							position851, tokenIndex851, depth851 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l852
							}
							position++
							goto l851
						l852:
							position, tokenIndex, depth = position851, tokenIndex851, depth851
							if buffer[position] != rune('H') {
								goto l843
							}
							position++
						}
					l851:
						if !_rules[rulesp]() {
							goto l843
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l843
						}
					l853:
						{
							position854, tokenIndex854, depth854 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l854
							}
							if buffer[position] != rune(',') {
								goto l854
							}
							position++
							if !_rules[rulespOpt]() {
								goto l854
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l854
							}
							goto l853
						l854:
							position, tokenIndex, depth = position854, tokenIndex854, depth854
						}
						goto l844
					l843:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
					}
				l844:
					depth--
					add(rulePegText, position842)
				}
				if !_rules[ruleAction41]() {
					goto l840
				}
				depth--
				add(ruleSourceSinkSpecs, position841)
			}
			return true
		l840:
			position, tokenIndex, depth = position840, tokenIndex840, depth840
			return false
		},
		/* 55 UpdateSourceSinkSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)> Action42)> */
		func() bool {
			position855, tokenIndex855, depth855 := position, tokenIndex, depth
			{
				position856 := position
				depth++
				{
					position857 := position
					depth++
					if !_rules[rulesp]() {
						goto l855
					}
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l859
						}
						position++
						goto l858
					l859:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
						if buffer[position] != rune('S') {
							goto l855
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
							goto l855
						}
						position++
					}
				l860:
					{
						position862, tokenIndex862, depth862 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l863
						}
						position++
						goto l862
					l863:
						position, tokenIndex, depth = position862, tokenIndex862, depth862
						if buffer[position] != rune('T') {
							goto l855
						}
						position++
					}
				l862:
					if !_rules[rulesp]() {
						goto l855
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l855
					}
				l864:
					{
						position865, tokenIndex865, depth865 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l865
						}
						if buffer[position] != rune(',') {
							goto l865
						}
						position++
						if !_rules[rulespOpt]() {
							goto l865
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l865
						}
						goto l864
					l865:
						position, tokenIndex, depth = position865, tokenIndex865, depth865
					}
					depth--
					add(rulePegText, position857)
				}
				if !_rules[ruleAction42]() {
					goto l855
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position856)
			}
			return true
		l855:
			position, tokenIndex, depth = position855, tokenIndex855, depth855
			return false
		},
		/* 56 SetOptSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action43)> */
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
						if !_rules[rulesp]() {
							goto l869
						}
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
								goto l869
							}
							position++
						}
					l871:
						{
							position873, tokenIndex873, depth873 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l874
							}
							position++
							goto l873
						l874:
							position, tokenIndex, depth = position873, tokenIndex873, depth873
							if buffer[position] != rune('E') {
								goto l869
							}
							position++
						}
					l873:
						{
							position875, tokenIndex875, depth875 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l876
							}
							position++
							goto l875
						l876:
							position, tokenIndex, depth = position875, tokenIndex875, depth875
							if buffer[position] != rune('T') {
								goto l869
							}
							position++
						}
					l875:
						if !_rules[rulesp]() {
							goto l869
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l869
						}
					l877:
						{
							position878, tokenIndex878, depth878 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l878
							}
							if buffer[position] != rune(',') {
								goto l878
							}
							position++
							if !_rules[rulespOpt]() {
								goto l878
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l878
							}
							goto l877
						l878:
							position, tokenIndex, depth = position878, tokenIndex878, depth878
						}
						goto l870
					l869:
						position, tokenIndex, depth = position869, tokenIndex869, depth869
					}
				l870:
					depth--
					add(rulePegText, position868)
				}
				if !_rules[ruleAction43]() {
					goto l866
				}
				depth--
				add(ruleSetOptSpecs, position867)
			}
			return true
		l866:
			position, tokenIndex, depth = position866, tokenIndex866, depth866
			return false
		},
		/* 57 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action44)> */
		func() bool {
			position879, tokenIndex879, depth879 := position, tokenIndex, depth
			{
				position880 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l879
				}
				if buffer[position] != rune('=') {
					goto l879
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l879
				}
				if !_rules[ruleAction44]() {
					goto l879
				}
				depth--
				add(ruleSourceSinkParam, position880)
			}
			return true
		l879:
			position, tokenIndex, depth = position879, tokenIndex879, depth879
			return false
		},
		/* 58 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr)> */
		func() bool {
			position881, tokenIndex881, depth881 := position, tokenIndex, depth
			{
				position882 := position
				depth++
				{
					position883, tokenIndex883, depth883 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l884
					}
					goto l883
				l884:
					position, tokenIndex, depth = position883, tokenIndex883, depth883
					if !_rules[ruleParamArrayExpr]() {
						goto l881
					}
				}
			l883:
				depth--
				add(ruleSourceSinkParamVal, position882)
			}
			return true
		l881:
			position, tokenIndex, depth = position881, tokenIndex881, depth881
			return false
		},
		/* 59 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position885, tokenIndex885, depth885 := position, tokenIndex, depth
			{
				position886 := position
				depth++
				{
					position887, tokenIndex887, depth887 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l888
					}
					goto l887
				l888:
					position, tokenIndex, depth = position887, tokenIndex887, depth887
					if !_rules[ruleLiteral]() {
						goto l885
					}
				}
			l887:
				depth--
				add(ruleParamLiteral, position886)
			}
			return true
		l885:
			position, tokenIndex, depth = position885, tokenIndex885, depth885
			return false
		},
		/* 60 ParamArrayExpr <- <(<('[' spOpt (ParamLiteral (',' spOpt ParamLiteral)*)? spOpt ','? spOpt ']')> Action45)> */
		func() bool {
			position889, tokenIndex889, depth889 := position, tokenIndex, depth
			{
				position890 := position
				depth++
				{
					position891 := position
					depth++
					if buffer[position] != rune('[') {
						goto l889
					}
					position++
					if !_rules[rulespOpt]() {
						goto l889
					}
					{
						position892, tokenIndex892, depth892 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l892
						}
					l894:
						{
							position895, tokenIndex895, depth895 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l895
							}
							position++
							if !_rules[rulespOpt]() {
								goto l895
							}
							if !_rules[ruleParamLiteral]() {
								goto l895
							}
							goto l894
						l895:
							position, tokenIndex, depth = position895, tokenIndex895, depth895
						}
						goto l893
					l892:
						position, tokenIndex, depth = position892, tokenIndex892, depth892
					}
				l893:
					if !_rules[rulespOpt]() {
						goto l889
					}
					{
						position896, tokenIndex896, depth896 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l896
						}
						position++
						goto l897
					l896:
						position, tokenIndex, depth = position896, tokenIndex896, depth896
					}
				l897:
					if !_rules[rulespOpt]() {
						goto l889
					}
					if buffer[position] != rune(']') {
						goto l889
					}
					position++
					depth--
					add(rulePegText, position891)
				}
				if !_rules[ruleAction45]() {
					goto l889
				}
				depth--
				add(ruleParamArrayExpr, position890)
			}
			return true
		l889:
			position, tokenIndex, depth = position889, tokenIndex889, depth889
			return false
		},
		/* 61 PausedOpt <- <(<(sp (Paused / Unpaused))?> Action46)> */
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
						if !_rules[rulesp]() {
							goto l901
						}
						{
							position903, tokenIndex903, depth903 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l904
							}
							goto l903
						l904:
							position, tokenIndex, depth = position903, tokenIndex903, depth903
							if !_rules[ruleUnpaused]() {
								goto l901
							}
						}
					l903:
						goto l902
					l901:
						position, tokenIndex, depth = position901, tokenIndex901, depth901
					}
				l902:
					depth--
					add(rulePegText, position900)
				}
				if !_rules[ruleAction46]() {
					goto l898
				}
				depth--
				add(rulePausedOpt, position899)
			}
			return true
		l898:
			position, tokenIndex, depth = position898, tokenIndex898, depth898
			return false
		},
		/* 62 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position905, tokenIndex905, depth905 := position, tokenIndex, depth
			{
				position906 := position
				depth++
				{
					position907, tokenIndex907, depth907 := position, tokenIndex, depth
					if !_rules[ruleWildcard]() {
						goto l908
					}
					goto l907
				l908:
					position, tokenIndex, depth = position907, tokenIndex907, depth907
					if !_rules[ruleExpression]() {
						goto l905
					}
				}
			l907:
				depth--
				add(ruleExpressionOrWildcard, position906)
			}
			return true
		l905:
			position, tokenIndex, depth = position905, tokenIndex905, depth905
			return false
		},
		/* 63 Expression <- <orExpr> */
		func() bool {
			position909, tokenIndex909, depth909 := position, tokenIndex, depth
			{
				position910 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l909
				}
				depth--
				add(ruleExpression, position910)
			}
			return true
		l909:
			position, tokenIndex, depth = position909, tokenIndex909, depth909
			return false
		},
		/* 64 orExpr <- <(<(andExpr (sp Or sp andExpr)*)> Action47)> */
		func() bool {
			position911, tokenIndex911, depth911 := position, tokenIndex, depth
			{
				position912 := position
				depth++
				{
					position913 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l911
					}
				l914:
					{
						position915, tokenIndex915, depth915 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l915
						}
						if !_rules[ruleOr]() {
							goto l915
						}
						if !_rules[rulesp]() {
							goto l915
						}
						if !_rules[ruleandExpr]() {
							goto l915
						}
						goto l914
					l915:
						position, tokenIndex, depth = position915, tokenIndex915, depth915
					}
					depth--
					add(rulePegText, position913)
				}
				if !_rules[ruleAction47]() {
					goto l911
				}
				depth--
				add(ruleorExpr, position912)
			}
			return true
		l911:
			position, tokenIndex, depth = position911, tokenIndex911, depth911
			return false
		},
		/* 65 andExpr <- <(<(notExpr (sp And sp notExpr)*)> Action48)> */
		func() bool {
			position916, tokenIndex916, depth916 := position, tokenIndex, depth
			{
				position917 := position
				depth++
				{
					position918 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l916
					}
				l919:
					{
						position920, tokenIndex920, depth920 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l920
						}
						if !_rules[ruleAnd]() {
							goto l920
						}
						if !_rules[rulesp]() {
							goto l920
						}
						if !_rules[rulenotExpr]() {
							goto l920
						}
						goto l919
					l920:
						position, tokenIndex, depth = position920, tokenIndex920, depth920
					}
					depth--
					add(rulePegText, position918)
				}
				if !_rules[ruleAction48]() {
					goto l916
				}
				depth--
				add(ruleandExpr, position917)
			}
			return true
		l916:
			position, tokenIndex, depth = position916, tokenIndex916, depth916
			return false
		},
		/* 66 notExpr <- <(<((Not sp)? comparisonExpr)> Action49)> */
		func() bool {
			position921, tokenIndex921, depth921 := position, tokenIndex, depth
			{
				position922 := position
				depth++
				{
					position923 := position
					depth++
					{
						position924, tokenIndex924, depth924 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l924
						}
						if !_rules[rulesp]() {
							goto l924
						}
						goto l925
					l924:
						position, tokenIndex, depth = position924, tokenIndex924, depth924
					}
				l925:
					if !_rules[rulecomparisonExpr]() {
						goto l921
					}
					depth--
					add(rulePegText, position923)
				}
				if !_rules[ruleAction49]() {
					goto l921
				}
				depth--
				add(rulenotExpr, position922)
			}
			return true
		l921:
			position, tokenIndex, depth = position921, tokenIndex921, depth921
			return false
		},
		/* 67 comparisonExpr <- <(<(otherOpExpr (spOpt ComparisonOp spOpt otherOpExpr)?)> Action50)> */
		func() bool {
			position926, tokenIndex926, depth926 := position, tokenIndex, depth
			{
				position927 := position
				depth++
				{
					position928 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l926
					}
					{
						position929, tokenIndex929, depth929 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l929
						}
						if !_rules[ruleComparisonOp]() {
							goto l929
						}
						if !_rules[rulespOpt]() {
							goto l929
						}
						if !_rules[ruleotherOpExpr]() {
							goto l929
						}
						goto l930
					l929:
						position, tokenIndex, depth = position929, tokenIndex929, depth929
					}
				l930:
					depth--
					add(rulePegText, position928)
				}
				if !_rules[ruleAction50]() {
					goto l926
				}
				depth--
				add(rulecomparisonExpr, position927)
			}
			return true
		l926:
			position, tokenIndex, depth = position926, tokenIndex926, depth926
			return false
		},
		/* 68 otherOpExpr <- <(<(isExpr (spOpt OtherOp spOpt isExpr)*)> Action51)> */
		func() bool {
			position931, tokenIndex931, depth931 := position, tokenIndex, depth
			{
				position932 := position
				depth++
				{
					position933 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l931
					}
				l934:
					{
						position935, tokenIndex935, depth935 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l935
						}
						if !_rules[ruleOtherOp]() {
							goto l935
						}
						if !_rules[rulespOpt]() {
							goto l935
						}
						if !_rules[ruleisExpr]() {
							goto l935
						}
						goto l934
					l935:
						position, tokenIndex, depth = position935, tokenIndex935, depth935
					}
					depth--
					add(rulePegText, position933)
				}
				if !_rules[ruleAction51]() {
					goto l931
				}
				depth--
				add(ruleotherOpExpr, position932)
			}
			return true
		l931:
			position, tokenIndex, depth = position931, tokenIndex931, depth931
			return false
		},
		/* 69 isExpr <- <(<(termExpr (sp IsOp sp NullLiteral)?)> Action52)> */
		func() bool {
			position936, tokenIndex936, depth936 := position, tokenIndex, depth
			{
				position937 := position
				depth++
				{
					position938 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l936
					}
					{
						position939, tokenIndex939, depth939 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l939
						}
						if !_rules[ruleIsOp]() {
							goto l939
						}
						if !_rules[rulesp]() {
							goto l939
						}
						if !_rules[ruleNullLiteral]() {
							goto l939
						}
						goto l940
					l939:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
					}
				l940:
					depth--
					add(rulePegText, position938)
				}
				if !_rules[ruleAction52]() {
					goto l936
				}
				depth--
				add(ruleisExpr, position937)
			}
			return true
		l936:
			position, tokenIndex, depth = position936, tokenIndex936, depth936
			return false
		},
		/* 70 termExpr <- <(<(productExpr (spOpt PlusMinusOp spOpt productExpr)*)> Action53)> */
		func() bool {
			position941, tokenIndex941, depth941 := position, tokenIndex, depth
			{
				position942 := position
				depth++
				{
					position943 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l941
					}
				l944:
					{
						position945, tokenIndex945, depth945 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l945
						}
						if !_rules[rulePlusMinusOp]() {
							goto l945
						}
						if !_rules[rulespOpt]() {
							goto l945
						}
						if !_rules[ruleproductExpr]() {
							goto l945
						}
						goto l944
					l945:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
					}
					depth--
					add(rulePegText, position943)
				}
				if !_rules[ruleAction53]() {
					goto l941
				}
				depth--
				add(ruletermExpr, position942)
			}
			return true
		l941:
			position, tokenIndex, depth = position941, tokenIndex941, depth941
			return false
		},
		/* 71 productExpr <- <(<(minusExpr (spOpt MultDivOp spOpt minusExpr)*)> Action54)> */
		func() bool {
			position946, tokenIndex946, depth946 := position, tokenIndex, depth
			{
				position947 := position
				depth++
				{
					position948 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l946
					}
				l949:
					{
						position950, tokenIndex950, depth950 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l950
						}
						if !_rules[ruleMultDivOp]() {
							goto l950
						}
						if !_rules[rulespOpt]() {
							goto l950
						}
						if !_rules[ruleminusExpr]() {
							goto l950
						}
						goto l949
					l950:
						position, tokenIndex, depth = position950, tokenIndex950, depth950
					}
					depth--
					add(rulePegText, position948)
				}
				if !_rules[ruleAction54]() {
					goto l946
				}
				depth--
				add(ruleproductExpr, position947)
			}
			return true
		l946:
			position, tokenIndex, depth = position946, tokenIndex946, depth946
			return false
		},
		/* 72 minusExpr <- <(<((UnaryMinus spOpt)? castExpr)> Action55)> */
		func() bool {
			position951, tokenIndex951, depth951 := position, tokenIndex, depth
			{
				position952 := position
				depth++
				{
					position953 := position
					depth++
					{
						position954, tokenIndex954, depth954 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l954
						}
						if !_rules[rulespOpt]() {
							goto l954
						}
						goto l955
					l954:
						position, tokenIndex, depth = position954, tokenIndex954, depth954
					}
				l955:
					if !_rules[rulecastExpr]() {
						goto l951
					}
					depth--
					add(rulePegText, position953)
				}
				if !_rules[ruleAction55]() {
					goto l951
				}
				depth--
				add(ruleminusExpr, position952)
			}
			return true
		l951:
			position, tokenIndex, depth = position951, tokenIndex951, depth951
			return false
		},
		/* 73 castExpr <- <(<(baseExpr (spOpt (':' ':') spOpt Type)?)> Action56)> */
		func() bool {
			position956, tokenIndex956, depth956 := position, tokenIndex, depth
			{
				position957 := position
				depth++
				{
					position958 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l956
					}
					{
						position959, tokenIndex959, depth959 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l959
						}
						if buffer[position] != rune(':') {
							goto l959
						}
						position++
						if buffer[position] != rune(':') {
							goto l959
						}
						position++
						if !_rules[rulespOpt]() {
							goto l959
						}
						if !_rules[ruleType]() {
							goto l959
						}
						goto l960
					l959:
						position, tokenIndex, depth = position959, tokenIndex959, depth959
					}
				l960:
					depth--
					add(rulePegText, position958)
				}
				if !_rules[ruleAction56]() {
					goto l956
				}
				depth--
				add(rulecastExpr, position957)
			}
			return true
		l956:
			position, tokenIndex, depth = position956, tokenIndex956, depth956
			return false
		},
		/* 74 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position961, tokenIndex961, depth961 := position, tokenIndex, depth
			{
				position962 := position
				depth++
				{
					position963, tokenIndex963, depth963 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l964
					}
					position++
					if !_rules[rulespOpt]() {
						goto l964
					}
					if !_rules[ruleExpression]() {
						goto l964
					}
					if !_rules[rulespOpt]() {
						goto l964
					}
					if buffer[position] != rune(')') {
						goto l964
					}
					position++
					goto l963
				l964:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleMapExpr]() {
						goto l965
					}
					goto l963
				l965:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleBooleanLiteral]() {
						goto l966
					}
					goto l963
				l966:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleNullLiteral]() {
						goto l967
					}
					goto l963
				l967:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleRowMeta]() {
						goto l968
					}
					goto l963
				l968:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleFuncTypeCast]() {
						goto l969
					}
					goto l963
				l969:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleFuncApp]() {
						goto l970
					}
					goto l963
				l970:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleRowValue]() {
						goto l971
					}
					goto l963
				l971:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleArrayExpr]() {
						goto l972
					}
					goto l963
				l972:
					position, tokenIndex, depth = position963, tokenIndex963, depth963
					if !_rules[ruleLiteral]() {
						goto l961
					}
				}
			l963:
				depth--
				add(rulebaseExpr, position962)
			}
			return true
		l961:
			position, tokenIndex, depth = position961, tokenIndex961, depth961
			return false
		},
		/* 75 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') spOpt '(' spOpt Expression sp (('a' / 'A') ('s' / 'S')) sp Type spOpt ')')> Action57)> */
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
						if buffer[position] != rune('c') {
							goto l977
						}
						position++
						goto l976
					l977:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('C') {
							goto l973
						}
						position++
					}
				l976:
					{
						position978, tokenIndex978, depth978 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l979
						}
						position++
						goto l978
					l979:
						position, tokenIndex, depth = position978, tokenIndex978, depth978
						if buffer[position] != rune('A') {
							goto l973
						}
						position++
					}
				l978:
					{
						position980, tokenIndex980, depth980 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l981
						}
						position++
						goto l980
					l981:
						position, tokenIndex, depth = position980, tokenIndex980, depth980
						if buffer[position] != rune('S') {
							goto l973
						}
						position++
					}
				l980:
					{
						position982, tokenIndex982, depth982 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l983
						}
						position++
						goto l982
					l983:
						position, tokenIndex, depth = position982, tokenIndex982, depth982
						if buffer[position] != rune('T') {
							goto l973
						}
						position++
					}
				l982:
					if !_rules[rulespOpt]() {
						goto l973
					}
					if buffer[position] != rune('(') {
						goto l973
					}
					position++
					if !_rules[rulespOpt]() {
						goto l973
					}
					if !_rules[ruleExpression]() {
						goto l973
					}
					if !_rules[rulesp]() {
						goto l973
					}
					{
						position984, tokenIndex984, depth984 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l985
						}
						position++
						goto l984
					l985:
						position, tokenIndex, depth = position984, tokenIndex984, depth984
						if buffer[position] != rune('A') {
							goto l973
						}
						position++
					}
				l984:
					{
						position986, tokenIndex986, depth986 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l987
						}
						position++
						goto l986
					l987:
						position, tokenIndex, depth = position986, tokenIndex986, depth986
						if buffer[position] != rune('S') {
							goto l973
						}
						position++
					}
				l986:
					if !_rules[rulesp]() {
						goto l973
					}
					if !_rules[ruleType]() {
						goto l973
					}
					if !_rules[rulespOpt]() {
						goto l973
					}
					if buffer[position] != rune(')') {
						goto l973
					}
					position++
					depth--
					add(rulePegText, position975)
				}
				if !_rules[ruleAction57]() {
					goto l973
				}
				depth--
				add(ruleFuncTypeCast, position974)
			}
			return true
		l973:
			position, tokenIndex, depth = position973, tokenIndex973, depth973
			return false
		},
		/* 76 FuncApp <- <(FuncAppWithOrderBy / FuncAppWithoutOrderBy)> */
		func() bool {
			position988, tokenIndex988, depth988 := position, tokenIndex, depth
			{
				position989 := position
				depth++
				{
					position990, tokenIndex990, depth990 := position, tokenIndex, depth
					if !_rules[ruleFuncAppWithOrderBy]() {
						goto l991
					}
					goto l990
				l991:
					position, tokenIndex, depth = position990, tokenIndex990, depth990
					if !_rules[ruleFuncAppWithoutOrderBy]() {
						goto l988
					}
				}
			l990:
				depth--
				add(ruleFuncApp, position989)
			}
			return true
		l988:
			position, tokenIndex, depth = position988, tokenIndex988, depth988
			return false
		},
		/* 77 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action58)> */
		func() bool {
			position992, tokenIndex992, depth992 := position, tokenIndex, depth
			{
				position993 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l992
				}
				if !_rules[rulespOpt]() {
					goto l992
				}
				if buffer[position] != rune('(') {
					goto l992
				}
				position++
				if !_rules[rulespOpt]() {
					goto l992
				}
				if !_rules[ruleFuncParams]() {
					goto l992
				}
				if !_rules[rulesp]() {
					goto l992
				}
				if !_rules[ruleParamsOrder]() {
					goto l992
				}
				if !_rules[rulespOpt]() {
					goto l992
				}
				if buffer[position] != rune(')') {
					goto l992
				}
				position++
				if !_rules[ruleAction58]() {
					goto l992
				}
				depth--
				add(ruleFuncAppWithOrderBy, position993)
			}
			return true
		l992:
			position, tokenIndex, depth = position992, tokenIndex992, depth992
			return false
		},
		/* 78 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action59)> */
		func() bool {
			position994, tokenIndex994, depth994 := position, tokenIndex, depth
			{
				position995 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l994
				}
				if !_rules[rulespOpt]() {
					goto l994
				}
				if buffer[position] != rune('(') {
					goto l994
				}
				position++
				if !_rules[rulespOpt]() {
					goto l994
				}
				if !_rules[ruleFuncParams]() {
					goto l994
				}
				{
					position996 := position
					depth++
					if !_rules[rulespOpt]() {
						goto l994
					}
					depth--
					add(rulePegText, position996)
				}
				if buffer[position] != rune(')') {
					goto l994
				}
				position++
				if !_rules[ruleAction59]() {
					goto l994
				}
				depth--
				add(ruleFuncAppWithoutOrderBy, position995)
			}
			return true
		l994:
			position, tokenIndex, depth = position994, tokenIndex994, depth994
			return false
		},
		/* 79 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action60)> */
		func() bool {
			position997, tokenIndex997, depth997 := position, tokenIndex, depth
			{
				position998 := position
				depth++
				{
					position999 := position
					depth++
					{
						position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1000
						}
					l1002:
						{
							position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1003
							}
							if buffer[position] != rune(',') {
								goto l1003
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1003
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1003
							}
							goto l1002
						l1003:
							position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
						}
						goto l1001
					l1000:
						position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
					}
				l1001:
					depth--
					add(rulePegText, position999)
				}
				if !_rules[ruleAction60]() {
					goto l997
				}
				depth--
				add(ruleFuncParams, position998)
			}
			return true
		l997:
			position, tokenIndex, depth = position997, tokenIndex997, depth997
			return false
		},
		/* 80 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action61)> */
		func() bool {
			position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
			{
				position1005 := position
				depth++
				{
					position1006 := position
					depth++
					{
						position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1008
						}
						position++
						goto l1007
					l1008:
						position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
						if buffer[position] != rune('O') {
							goto l1004
						}
						position++
					}
				l1007:
					{
						position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1010
						}
						position++
						goto l1009
					l1010:
						position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
						if buffer[position] != rune('R') {
							goto l1004
						}
						position++
					}
				l1009:
					{
						position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1012
						}
						position++
						goto l1011
					l1012:
						position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
						if buffer[position] != rune('D') {
							goto l1004
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
							goto l1004
						}
						position++
					}
				l1013:
					{
						position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1016
						}
						position++
						goto l1015
					l1016:
						position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
						if buffer[position] != rune('R') {
							goto l1004
						}
						position++
					}
				l1015:
					if !_rules[rulesp]() {
						goto l1004
					}
					{
						position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1018
						}
						position++
						goto l1017
					l1018:
						position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
						if buffer[position] != rune('B') {
							goto l1004
						}
						position++
					}
				l1017:
					{
						position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1020
						}
						position++
						goto l1019
					l1020:
						position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
						if buffer[position] != rune('Y') {
							goto l1004
						}
						position++
					}
				l1019:
					if !_rules[rulesp]() {
						goto l1004
					}
					if !_rules[ruleSortedExpression]() {
						goto l1004
					}
				l1021:
					{
						position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1022
						}
						if buffer[position] != rune(',') {
							goto l1022
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1022
						}
						if !_rules[ruleSortedExpression]() {
							goto l1022
						}
						goto l1021
					l1022:
						position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
					}
					depth--
					add(rulePegText, position1006)
				}
				if !_rules[ruleAction61]() {
					goto l1004
				}
				depth--
				add(ruleParamsOrder, position1005)
			}
			return true
		l1004:
			position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
			return false
		},
		/* 81 SortedExpression <- <(Expression OrderDirectionOpt Action62)> */
		func() bool {
			position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
			{
				position1024 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l1023
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1023
				}
				if !_rules[ruleAction62]() {
					goto l1023
				}
				depth--
				add(ruleSortedExpression, position1024)
			}
			return true
		l1023:
			position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
			return false
		},
		/* 82 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action63)> */
		func() bool {
			position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
			{
				position1026 := position
				depth++
				{
					position1027 := position
					depth++
					{
						position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1028
						}
						{
							position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
							if !_rules[ruleAscending]() {
								goto l1031
							}
							goto l1030
						l1031:
							position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
							if !_rules[ruleDescending]() {
								goto l1028
							}
						}
					l1030:
						goto l1029
					l1028:
						position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
					}
				l1029:
					depth--
					add(rulePegText, position1027)
				}
				if !_rules[ruleAction63]() {
					goto l1025
				}
				depth--
				add(ruleOrderDirectionOpt, position1026)
			}
			return true
		l1025:
			position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
			return false
		},
		/* 83 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action64)> */
		func() bool {
			position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
			{
				position1033 := position
				depth++
				{
					position1034 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1032
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1032
					}
					{
						position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1035
						}
					l1037:
						{
							position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1038
							}
							if buffer[position] != rune(',') {
								goto l1038
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1038
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1038
							}
							goto l1037
						l1038:
							position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
						}
						goto l1036
					l1035:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
					}
				l1036:
					if !_rules[rulespOpt]() {
						goto l1032
					}
					{
						position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1039
						}
						position++
						goto l1040
					l1039:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
					}
				l1040:
					if !_rules[rulespOpt]() {
						goto l1032
					}
					if buffer[position] != rune(']') {
						goto l1032
					}
					position++
					depth--
					add(rulePegText, position1034)
				}
				if !_rules[ruleAction64]() {
					goto l1032
				}
				depth--
				add(ruleArrayExpr, position1033)
			}
			return true
		l1032:
			position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
			return false
		},
		/* 84 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action65)> */
		func() bool {
			position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
			{
				position1042 := position
				depth++
				{
					position1043 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1041
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1041
					}
					{
						position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l1044
						}
					l1046:
						{
							position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1047
							}
							if buffer[position] != rune(',') {
								goto l1047
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1047
							}
							if !_rules[ruleKeyValuePair]() {
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
					if buffer[position] != rune('}') {
						goto l1041
					}
					position++
					depth--
					add(rulePegText, position1043)
				}
				if !_rules[ruleAction65]() {
					goto l1041
				}
				depth--
				add(ruleMapExpr, position1042)
			}
			return true
		l1041:
			position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
			return false
		},
		/* 85 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action66)> */
		func() bool {
			position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
			{
				position1049 := position
				depth++
				{
					position1050 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1048
					}
					if !_rules[rulespOpt]() {
						goto l1048
					}
					if buffer[position] != rune(':') {
						goto l1048
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1048
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1048
					}
					depth--
					add(rulePegText, position1050)
				}
				if !_rules[ruleAction66]() {
					goto l1048
				}
				depth--
				add(ruleKeyValuePair, position1049)
			}
			return true
		l1048:
			position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
			return false
		},
		/* 86 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
			{
				position1052 := position
				depth++
				{
					position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l1054
					}
					goto l1053
				l1054:
					position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
					if !_rules[ruleNumericLiteral]() {
						goto l1055
					}
					goto l1053
				l1055:
					position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
					if !_rules[ruleStringLiteral]() {
						goto l1051
					}
				}
			l1053:
				depth--
				add(ruleLiteral, position1052)
			}
			return true
		l1051:
			position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
			return false
		},
		/* 87 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
			{
				position1057 := position
				depth++
				{
					position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l1059
					}
					goto l1058
				l1059:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if !_rules[ruleNotEqual]() {
						goto l1060
					}
					goto l1058
				l1060:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if !_rules[ruleLessOrEqual]() {
						goto l1061
					}
					goto l1058
				l1061:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if !_rules[ruleLess]() {
						goto l1062
					}
					goto l1058
				l1062:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if !_rules[ruleGreaterOrEqual]() {
						goto l1063
					}
					goto l1058
				l1063:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if !_rules[ruleGreater]() {
						goto l1064
					}
					goto l1058
				l1064:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if !_rules[ruleNotEqual]() {
						goto l1056
					}
				}
			l1058:
				depth--
				add(ruleComparisonOp, position1057)
			}
			return true
		l1056:
			position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
			return false
		},
		/* 88 OtherOp <- <Concat> */
		func() bool {
			position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
			{
				position1066 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l1065
				}
				depth--
				add(ruleOtherOp, position1066)
			}
			return true
		l1065:
			position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
			return false
		},
		/* 89 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
			{
				position1068 := position
				depth++
				{
					position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l1070
					}
					goto l1069
				l1070:
					position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
					if !_rules[ruleIs]() {
						goto l1067
					}
				}
			l1069:
				depth--
				add(ruleIsOp, position1068)
			}
			return true
		l1067:
			position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
			return false
		},
		/* 90 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
			{
				position1072 := position
				depth++
				{
					position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l1074
					}
					goto l1073
				l1074:
					position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
					if !_rules[ruleMinus]() {
						goto l1071
					}
				}
			l1073:
				depth--
				add(rulePlusMinusOp, position1072)
			}
			return true
		l1071:
			position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
			return false
		},
		/* 91 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
			{
				position1076 := position
				depth++
				{
					position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l1078
					}
					goto l1077
				l1078:
					position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
					if !_rules[ruleDivide]() {
						goto l1079
					}
					goto l1077
				l1079:
					position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
					if !_rules[ruleModulo]() {
						goto l1075
					}
				}
			l1077:
				depth--
				add(ruleMultDivOp, position1076)
			}
			return true
		l1075:
			position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
			return false
		},
		/* 92 Stream <- <(<ident> Action67)> */
		func() bool {
			position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
			{
				position1081 := position
				depth++
				{
					position1082 := position
					depth++
					if !_rules[ruleident]() {
						goto l1080
					}
					depth--
					add(rulePegText, position1082)
				}
				if !_rules[ruleAction67]() {
					goto l1080
				}
				depth--
				add(ruleStream, position1081)
			}
			return true
		l1080:
			position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
			return false
		},
		/* 93 RowMeta <- <RowTimestamp> */
		func() bool {
			position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
			{
				position1084 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l1083
				}
				depth--
				add(ruleRowMeta, position1084)
			}
			return true
		l1083:
			position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
			return false
		},
		/* 94 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action68)> */
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
						if !_rules[ruleident]() {
							goto l1088
						}
						if buffer[position] != rune(':') {
							goto l1088
						}
						position++
						goto l1089
					l1088:
						position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
					}
				l1089:
					if buffer[position] != rune('t') {
						goto l1085
					}
					position++
					if buffer[position] != rune('s') {
						goto l1085
					}
					position++
					if buffer[position] != rune('(') {
						goto l1085
					}
					position++
					if buffer[position] != rune(')') {
						goto l1085
					}
					position++
					depth--
					add(rulePegText, position1087)
				}
				if !_rules[ruleAction68]() {
					goto l1085
				}
				depth--
				add(ruleRowTimestamp, position1086)
			}
			return true
		l1085:
			position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
			return false
		},
		/* 95 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action69)> */
		func() bool {
			position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
			{
				position1091 := position
				depth++
				{
					position1092 := position
					depth++
					{
						position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1093
						}
						if buffer[position] != rune(':') {
							goto l1093
						}
						position++
						{
							position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1095
							}
							position++
							goto l1093
						l1095:
							position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
						}
						goto l1094
					l1093:
						position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
					}
				l1094:
					if !_rules[rulejsonGetPath]() {
						goto l1090
					}
					depth--
					add(rulePegText, position1092)
				}
				if !_rules[ruleAction69]() {
					goto l1090
				}
				depth--
				add(ruleRowValue, position1091)
			}
			return true
		l1090:
			position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
			return false
		},
		/* 96 NumericLiteral <- <(<('-'? [0-9]+)> Action70)> */
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
						if buffer[position] != rune('-') {
							goto l1099
						}
						position++
						goto l1100
					l1099:
						position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
					}
				l1100:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1096
					}
					position++
				l1101:
					{
						position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1102
						}
						position++
						goto l1101
					l1102:
						position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
					}
					depth--
					add(rulePegText, position1098)
				}
				if !_rules[ruleAction70]() {
					goto l1096
				}
				depth--
				add(ruleNumericLiteral, position1097)
			}
			return true
		l1096:
			position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
			return false
		},
		/* 97 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action71)> */
		func() bool {
			position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
			{
				position1104 := position
				depth++
				{
					position1105 := position
					depth++
					{
						position1106, tokenIndex1106, depth1106 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1106
						}
						position++
						goto l1107
					l1106:
						position, tokenIndex, depth = position1106, tokenIndex1106, depth1106
					}
				l1107:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1103
					}
					position++
				l1108:
					{
						position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1109
						}
						position++
						goto l1108
					l1109:
						position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
					}
					if buffer[position] != rune('.') {
						goto l1103
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1103
					}
					position++
				l1110:
					{
						position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1111
						}
						position++
						goto l1110
					l1111:
						position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					}
					depth--
					add(rulePegText, position1105)
				}
				if !_rules[ruleAction71]() {
					goto l1103
				}
				depth--
				add(ruleFloatLiteral, position1104)
			}
			return true
		l1103:
			position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
			return false
		},
		/* 98 Function <- <(<ident> Action72)> */
		func() bool {
			position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
			{
				position1113 := position
				depth++
				{
					position1114 := position
					depth++
					if !_rules[ruleident]() {
						goto l1112
					}
					depth--
					add(rulePegText, position1114)
				}
				if !_rules[ruleAction72]() {
					goto l1112
				}
				depth--
				add(ruleFunction, position1113)
			}
			return true
		l1112:
			position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
			return false
		},
		/* 99 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action73)> */
		func() bool {
			position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
			{
				position1116 := position
				depth++
				{
					position1117 := position
					depth++
					{
						position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1119
						}
						position++
						goto l1118
					l1119:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						if buffer[position] != rune('N') {
							goto l1115
						}
						position++
					}
				l1118:
					{
						position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1121
						}
						position++
						goto l1120
					l1121:
						position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
						if buffer[position] != rune('U') {
							goto l1115
						}
						position++
					}
				l1120:
					{
						position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1123
						}
						position++
						goto l1122
					l1123:
						position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
						if buffer[position] != rune('L') {
							goto l1115
						}
						position++
					}
				l1122:
					{
						position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1125
						}
						position++
						goto l1124
					l1125:
						position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
						if buffer[position] != rune('L') {
							goto l1115
						}
						position++
					}
				l1124:
					depth--
					add(rulePegText, position1117)
				}
				if !_rules[ruleAction73]() {
					goto l1115
				}
				depth--
				add(ruleNullLiteral, position1116)
			}
			return true
		l1115:
			position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
			return false
		},
		/* 100 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
			{
				position1127 := position
				depth++
				{
					position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l1129
					}
					goto l1128
				l1129:
					position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
					if !_rules[ruleFALSE]() {
						goto l1126
					}
				}
			l1128:
				depth--
				add(ruleBooleanLiteral, position1127)
			}
			return true
		l1126:
			position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
			return false
		},
		/* 101 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action74)> */
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
						if buffer[position] != rune('t') {
							goto l1134
						}
						position++
						goto l1133
					l1134:
						position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
						if buffer[position] != rune('T') {
							goto l1130
						}
						position++
					}
				l1133:
					{
						position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
						if buffer[position] != rune('R') {
							goto l1130
						}
						position++
					}
				l1135:
					{
						position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
						if buffer[position] != rune('U') {
							goto l1130
						}
						position++
					}
				l1137:
					{
						position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1140
						}
						position++
						goto l1139
					l1140:
						position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
						if buffer[position] != rune('E') {
							goto l1130
						}
						position++
					}
				l1139:
					depth--
					add(rulePegText, position1132)
				}
				if !_rules[ruleAction74]() {
					goto l1130
				}
				depth--
				add(ruleTRUE, position1131)
			}
			return true
		l1130:
			position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
			return false
		},
		/* 102 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action75)> */
		func() bool {
			position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
			{
				position1142 := position
				depth++
				{
					position1143 := position
					depth++
					{
						position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1145
						}
						position++
						goto l1144
					l1145:
						position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
						if buffer[position] != rune('F') {
							goto l1141
						}
						position++
					}
				l1144:
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
							goto l1141
						}
						position++
					}
				l1146:
					{
						position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1149
						}
						position++
						goto l1148
					l1149:
						position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
						if buffer[position] != rune('L') {
							goto l1141
						}
						position++
					}
				l1148:
					{
						position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1151
						}
						position++
						goto l1150
					l1151:
						position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
						if buffer[position] != rune('S') {
							goto l1141
						}
						position++
					}
				l1150:
					{
						position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1153
						}
						position++
						goto l1152
					l1153:
						position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
						if buffer[position] != rune('E') {
							goto l1141
						}
						position++
					}
				l1152:
					depth--
					add(rulePegText, position1143)
				}
				if !_rules[ruleAction75]() {
					goto l1141
				}
				depth--
				add(ruleFALSE, position1142)
			}
			return true
		l1141:
			position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
			return false
		},
		/* 103 Wildcard <- <(<((ident ':' !':')? '*')> Action76)> */
		func() bool {
			position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
			{
				position1155 := position
				depth++
				{
					position1156 := position
					depth++
					{
						position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1157
						}
						if buffer[position] != rune(':') {
							goto l1157
						}
						position++
						{
							position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1159
							}
							position++
							goto l1157
						l1159:
							position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
						}
						goto l1158
					l1157:
						position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
					}
				l1158:
					if buffer[position] != rune('*') {
						goto l1154
					}
					position++
					depth--
					add(rulePegText, position1156)
				}
				if !_rules[ruleAction76]() {
					goto l1154
				}
				depth--
				add(ruleWildcard, position1155)
			}
			return true
		l1154:
			position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
			return false
		},
		/* 104 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action77)> */
		func() bool {
			position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
			{
				position1161 := position
				depth++
				{
					position1162 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l1160
					}
					position++
				l1163:
					{
						position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
						{
							position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1166
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1166
							}
							position++
							goto l1165
						l1166:
							position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
							{
								position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1167
								}
								position++
								goto l1164
							l1167:
								position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
							}
							if !matchDot() {
								goto l1164
							}
						}
					l1165:
						goto l1163
					l1164:
						position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
					}
					if buffer[position] != rune('\'') {
						goto l1160
					}
					position++
					depth--
					add(rulePegText, position1162)
				}
				if !_rules[ruleAction77]() {
					goto l1160
				}
				depth--
				add(ruleStringLiteral, position1161)
			}
			return true
		l1160:
			position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
			return false
		},
		/* 105 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action78)> */
		func() bool {
			position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
			{
				position1169 := position
				depth++
				{
					position1170 := position
					depth++
					{
						position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1172
						}
						position++
						goto l1171
					l1172:
						position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
						if buffer[position] != rune('I') {
							goto l1168
						}
						position++
					}
				l1171:
					{
						position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1174
						}
						position++
						goto l1173
					l1174:
						position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
						if buffer[position] != rune('S') {
							goto l1168
						}
						position++
					}
				l1173:
					{
						position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1176
						}
						position++
						goto l1175
					l1176:
						position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
						if buffer[position] != rune('T') {
							goto l1168
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
							goto l1168
						}
						position++
					}
				l1177:
					{
						position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1180
						}
						position++
						goto l1179
					l1180:
						position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
						if buffer[position] != rune('E') {
							goto l1168
						}
						position++
					}
				l1179:
					{
						position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1182
						}
						position++
						goto l1181
					l1182:
						position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
						if buffer[position] != rune('A') {
							goto l1168
						}
						position++
					}
				l1181:
					{
						position1183, tokenIndex1183, depth1183 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1184
						}
						position++
						goto l1183
					l1184:
						position, tokenIndex, depth = position1183, tokenIndex1183, depth1183
						if buffer[position] != rune('M') {
							goto l1168
						}
						position++
					}
				l1183:
					depth--
					add(rulePegText, position1170)
				}
				if !_rules[ruleAction78]() {
					goto l1168
				}
				depth--
				add(ruleISTREAM, position1169)
			}
			return true
		l1168:
			position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
			return false
		},
		/* 106 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action79)> */
		func() bool {
			position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
			{
				position1186 := position
				depth++
				{
					position1187 := position
					depth++
					{
						position1188, tokenIndex1188, depth1188 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1189
						}
						position++
						goto l1188
					l1189:
						position, tokenIndex, depth = position1188, tokenIndex1188, depth1188
						if buffer[position] != rune('D') {
							goto l1185
						}
						position++
					}
				l1188:
					{
						position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1191
						}
						position++
						goto l1190
					l1191:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('S') {
							goto l1185
						}
						position++
					}
				l1190:
					{
						position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1193
						}
						position++
						goto l1192
					l1193:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						if buffer[position] != rune('T') {
							goto l1185
						}
						position++
					}
				l1192:
					{
						position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1195
						}
						position++
						goto l1194
					l1195:
						position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
						if buffer[position] != rune('R') {
							goto l1185
						}
						position++
					}
				l1194:
					{
						position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1197
						}
						position++
						goto l1196
					l1197:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
						if buffer[position] != rune('E') {
							goto l1185
						}
						position++
					}
				l1196:
					{
						position1198, tokenIndex1198, depth1198 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1199
						}
						position++
						goto l1198
					l1199:
						position, tokenIndex, depth = position1198, tokenIndex1198, depth1198
						if buffer[position] != rune('A') {
							goto l1185
						}
						position++
					}
				l1198:
					{
						position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1201
						}
						position++
						goto l1200
					l1201:
						position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
						if buffer[position] != rune('M') {
							goto l1185
						}
						position++
					}
				l1200:
					depth--
					add(rulePegText, position1187)
				}
				if !_rules[ruleAction79]() {
					goto l1185
				}
				depth--
				add(ruleDSTREAM, position1186)
			}
			return true
		l1185:
			position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
			return false
		},
		/* 107 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action80)> */
		func() bool {
			position1202, tokenIndex1202, depth1202 := position, tokenIndex, depth
			{
				position1203 := position
				depth++
				{
					position1204 := position
					depth++
					{
						position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1206
						}
						position++
						goto l1205
					l1206:
						position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
						if buffer[position] != rune('R') {
							goto l1202
						}
						position++
					}
				l1205:
					{
						position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1208
						}
						position++
						goto l1207
					l1208:
						position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
						if buffer[position] != rune('S') {
							goto l1202
						}
						position++
					}
				l1207:
					{
						position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1210
						}
						position++
						goto l1209
					l1210:
						position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
						if buffer[position] != rune('T') {
							goto l1202
						}
						position++
					}
				l1209:
					{
						position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1212
						}
						position++
						goto l1211
					l1212:
						position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
						if buffer[position] != rune('R') {
							goto l1202
						}
						position++
					}
				l1211:
					{
						position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1214
						}
						position++
						goto l1213
					l1214:
						position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
						if buffer[position] != rune('E') {
							goto l1202
						}
						position++
					}
				l1213:
					{
						position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1216
						}
						position++
						goto l1215
					l1216:
						position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
						if buffer[position] != rune('A') {
							goto l1202
						}
						position++
					}
				l1215:
					{
						position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1218
						}
						position++
						goto l1217
					l1218:
						position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
						if buffer[position] != rune('M') {
							goto l1202
						}
						position++
					}
				l1217:
					depth--
					add(rulePegText, position1204)
				}
				if !_rules[ruleAction80]() {
					goto l1202
				}
				depth--
				add(ruleRSTREAM, position1203)
			}
			return true
		l1202:
			position, tokenIndex, depth = position1202, tokenIndex1202, depth1202
			return false
		},
		/* 108 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action81)> */
		func() bool {
			position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
			{
				position1220 := position
				depth++
				{
					position1221 := position
					depth++
					{
						position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1223
						}
						position++
						goto l1222
					l1223:
						position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
						if buffer[position] != rune('T') {
							goto l1219
						}
						position++
					}
				l1222:
					{
						position1224, tokenIndex1224, depth1224 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1225
						}
						position++
						goto l1224
					l1225:
						position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
						if buffer[position] != rune('U') {
							goto l1219
						}
						position++
					}
				l1224:
					{
						position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1227
						}
						position++
						goto l1226
					l1227:
						position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
						if buffer[position] != rune('P') {
							goto l1219
						}
						position++
					}
				l1226:
					{
						position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1229
						}
						position++
						goto l1228
					l1229:
						position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
						if buffer[position] != rune('L') {
							goto l1219
						}
						position++
					}
				l1228:
					{
						position1230, tokenIndex1230, depth1230 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1231
						}
						position++
						goto l1230
					l1231:
						position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
						if buffer[position] != rune('E') {
							goto l1219
						}
						position++
					}
				l1230:
					{
						position1232, tokenIndex1232, depth1232 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1233
						}
						position++
						goto l1232
					l1233:
						position, tokenIndex, depth = position1232, tokenIndex1232, depth1232
						if buffer[position] != rune('S') {
							goto l1219
						}
						position++
					}
				l1232:
					depth--
					add(rulePegText, position1221)
				}
				if !_rules[ruleAction81]() {
					goto l1219
				}
				depth--
				add(ruleTUPLES, position1220)
			}
			return true
		l1219:
			position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
			return false
		},
		/* 109 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action82)> */
		func() bool {
			position1234, tokenIndex1234, depth1234 := position, tokenIndex, depth
			{
				position1235 := position
				depth++
				{
					position1236 := position
					depth++
					{
						position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1238
						}
						position++
						goto l1237
					l1238:
						position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
						if buffer[position] != rune('S') {
							goto l1234
						}
						position++
					}
				l1237:
					{
						position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1240
						}
						position++
						goto l1239
					l1240:
						position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
						if buffer[position] != rune('E') {
							goto l1234
						}
						position++
					}
				l1239:
					{
						position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1242
						}
						position++
						goto l1241
					l1242:
						position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
						if buffer[position] != rune('C') {
							goto l1234
						}
						position++
					}
				l1241:
					{
						position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1244
						}
						position++
						goto l1243
					l1244:
						position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
						if buffer[position] != rune('O') {
							goto l1234
						}
						position++
					}
				l1243:
					{
						position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1246
						}
						position++
						goto l1245
					l1246:
						position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
						if buffer[position] != rune('N') {
							goto l1234
						}
						position++
					}
				l1245:
					{
						position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1248
						}
						position++
						goto l1247
					l1248:
						position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
						if buffer[position] != rune('D') {
							goto l1234
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
							goto l1234
						}
						position++
					}
				l1249:
					depth--
					add(rulePegText, position1236)
				}
				if !_rules[ruleAction82]() {
					goto l1234
				}
				depth--
				add(ruleSECONDS, position1235)
			}
			return true
		l1234:
			position, tokenIndex, depth = position1234, tokenIndex1234, depth1234
			return false
		},
		/* 110 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action83)> */
		func() bool {
			position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
			{
				position1252 := position
				depth++
				{
					position1253 := position
					depth++
					{
						position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1255
						}
						position++
						goto l1254
					l1255:
						position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
						if buffer[position] != rune('M') {
							goto l1251
						}
						position++
					}
				l1254:
					{
						position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1257
						}
						position++
						goto l1256
					l1257:
						position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
						if buffer[position] != rune('I') {
							goto l1251
						}
						position++
					}
				l1256:
					{
						position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1259
						}
						position++
						goto l1258
					l1259:
						position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
						if buffer[position] != rune('L') {
							goto l1251
						}
						position++
					}
				l1258:
					{
						position1260, tokenIndex1260, depth1260 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1261
						}
						position++
						goto l1260
					l1261:
						position, tokenIndex, depth = position1260, tokenIndex1260, depth1260
						if buffer[position] != rune('L') {
							goto l1251
						}
						position++
					}
				l1260:
					{
						position1262, tokenIndex1262, depth1262 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1263
						}
						position++
						goto l1262
					l1263:
						position, tokenIndex, depth = position1262, tokenIndex1262, depth1262
						if buffer[position] != rune('I') {
							goto l1251
						}
						position++
					}
				l1262:
					{
						position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1265
						}
						position++
						goto l1264
					l1265:
						position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
						if buffer[position] != rune('S') {
							goto l1251
						}
						position++
					}
				l1264:
					{
						position1266, tokenIndex1266, depth1266 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1267
						}
						position++
						goto l1266
					l1267:
						position, tokenIndex, depth = position1266, tokenIndex1266, depth1266
						if buffer[position] != rune('E') {
							goto l1251
						}
						position++
					}
				l1266:
					{
						position1268, tokenIndex1268, depth1268 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1269
						}
						position++
						goto l1268
					l1269:
						position, tokenIndex, depth = position1268, tokenIndex1268, depth1268
						if buffer[position] != rune('C') {
							goto l1251
						}
						position++
					}
				l1268:
					{
						position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1271
						}
						position++
						goto l1270
					l1271:
						position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
						if buffer[position] != rune('O') {
							goto l1251
						}
						position++
					}
				l1270:
					{
						position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1273
						}
						position++
						goto l1272
					l1273:
						position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
						if buffer[position] != rune('N') {
							goto l1251
						}
						position++
					}
				l1272:
					{
						position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1275
						}
						position++
						goto l1274
					l1275:
						position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
						if buffer[position] != rune('D') {
							goto l1251
						}
						position++
					}
				l1274:
					{
						position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1277
						}
						position++
						goto l1276
					l1277:
						position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
						if buffer[position] != rune('S') {
							goto l1251
						}
						position++
					}
				l1276:
					depth--
					add(rulePegText, position1253)
				}
				if !_rules[ruleAction83]() {
					goto l1251
				}
				depth--
				add(ruleMILLISECONDS, position1252)
			}
			return true
		l1251:
			position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
			return false
		},
		/* 111 StreamIdentifier <- <(<ident> Action84)> */
		func() bool {
			position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
			{
				position1279 := position
				depth++
				{
					position1280 := position
					depth++
					if !_rules[ruleident]() {
						goto l1278
					}
					depth--
					add(rulePegText, position1280)
				}
				if !_rules[ruleAction84]() {
					goto l1278
				}
				depth--
				add(ruleStreamIdentifier, position1279)
			}
			return true
		l1278:
			position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
			return false
		},
		/* 112 SourceSinkType <- <(<ident> Action85)> */
		func() bool {
			position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
			{
				position1282 := position
				depth++
				{
					position1283 := position
					depth++
					if !_rules[ruleident]() {
						goto l1281
					}
					depth--
					add(rulePegText, position1283)
				}
				if !_rules[ruleAction85]() {
					goto l1281
				}
				depth--
				add(ruleSourceSinkType, position1282)
			}
			return true
		l1281:
			position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
			return false
		},
		/* 113 SourceSinkParamKey <- <(<ident> Action86)> */
		func() bool {
			position1284, tokenIndex1284, depth1284 := position, tokenIndex, depth
			{
				position1285 := position
				depth++
				{
					position1286 := position
					depth++
					if !_rules[ruleident]() {
						goto l1284
					}
					depth--
					add(rulePegText, position1286)
				}
				if !_rules[ruleAction86]() {
					goto l1284
				}
				depth--
				add(ruleSourceSinkParamKey, position1285)
			}
			return true
		l1284:
			position, tokenIndex, depth = position1284, tokenIndex1284, depth1284
			return false
		},
		/* 114 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action87)> */
		func() bool {
			position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
			{
				position1288 := position
				depth++
				{
					position1289 := position
					depth++
					{
						position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1291
						}
						position++
						goto l1290
					l1291:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if buffer[position] != rune('P') {
							goto l1287
						}
						position++
					}
				l1290:
					{
						position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1293
						}
						position++
						goto l1292
					l1293:
						position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
						if buffer[position] != rune('A') {
							goto l1287
						}
						position++
					}
				l1292:
					{
						position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1295
						}
						position++
						goto l1294
					l1295:
						position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
						if buffer[position] != rune('U') {
							goto l1287
						}
						position++
					}
				l1294:
					{
						position1296, tokenIndex1296, depth1296 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1297
						}
						position++
						goto l1296
					l1297:
						position, tokenIndex, depth = position1296, tokenIndex1296, depth1296
						if buffer[position] != rune('S') {
							goto l1287
						}
						position++
					}
				l1296:
					{
						position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1299
						}
						position++
						goto l1298
					l1299:
						position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
						if buffer[position] != rune('E') {
							goto l1287
						}
						position++
					}
				l1298:
					{
						position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1301
						}
						position++
						goto l1300
					l1301:
						position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
						if buffer[position] != rune('D') {
							goto l1287
						}
						position++
					}
				l1300:
					depth--
					add(rulePegText, position1289)
				}
				if !_rules[ruleAction87]() {
					goto l1287
				}
				depth--
				add(rulePaused, position1288)
			}
			return true
		l1287:
			position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
			return false
		},
		/* 115 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action88)> */
		func() bool {
			position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
			{
				position1303 := position
				depth++
				{
					position1304 := position
					depth++
					{
						position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1306
						}
						position++
						goto l1305
					l1306:
						position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
						if buffer[position] != rune('U') {
							goto l1302
						}
						position++
					}
				l1305:
					{
						position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1308
						}
						position++
						goto l1307
					l1308:
						position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
						if buffer[position] != rune('N') {
							goto l1302
						}
						position++
					}
				l1307:
					{
						position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1310
						}
						position++
						goto l1309
					l1310:
						position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
						if buffer[position] != rune('P') {
							goto l1302
						}
						position++
					}
				l1309:
					{
						position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1312
						}
						position++
						goto l1311
					l1312:
						position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
						if buffer[position] != rune('A') {
							goto l1302
						}
						position++
					}
				l1311:
					{
						position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1314
						}
						position++
						goto l1313
					l1314:
						position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
						if buffer[position] != rune('U') {
							goto l1302
						}
						position++
					}
				l1313:
					{
						position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1316
						}
						position++
						goto l1315
					l1316:
						position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
						if buffer[position] != rune('S') {
							goto l1302
						}
						position++
					}
				l1315:
					{
						position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1318
						}
						position++
						goto l1317
					l1318:
						position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
						if buffer[position] != rune('E') {
							goto l1302
						}
						position++
					}
				l1317:
					{
						position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1320
						}
						position++
						goto l1319
					l1320:
						position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
						if buffer[position] != rune('D') {
							goto l1302
						}
						position++
					}
				l1319:
					depth--
					add(rulePegText, position1304)
				}
				if !_rules[ruleAction88]() {
					goto l1302
				}
				depth--
				add(ruleUnpaused, position1303)
			}
			return true
		l1302:
			position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
			return false
		},
		/* 116 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action89)> */
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
						if buffer[position] != rune('a') {
							goto l1325
						}
						position++
						goto l1324
					l1325:
						position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
						if buffer[position] != rune('A') {
							goto l1321
						}
						position++
					}
				l1324:
					{
						position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1327
						}
						position++
						goto l1326
					l1327:
						position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
						if buffer[position] != rune('S') {
							goto l1321
						}
						position++
					}
				l1326:
					{
						position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1329
						}
						position++
						goto l1328
					l1329:
						position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
						if buffer[position] != rune('C') {
							goto l1321
						}
						position++
					}
				l1328:
					depth--
					add(rulePegText, position1323)
				}
				if !_rules[ruleAction89]() {
					goto l1321
				}
				depth--
				add(ruleAscending, position1322)
			}
			return true
		l1321:
			position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
			return false
		},
		/* 117 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action90)> */
		func() bool {
			position1330, tokenIndex1330, depth1330 := position, tokenIndex, depth
			{
				position1331 := position
				depth++
				{
					position1332 := position
					depth++
					{
						position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1334
						}
						position++
						goto l1333
					l1334:
						position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
						if buffer[position] != rune('D') {
							goto l1330
						}
						position++
					}
				l1333:
					{
						position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1336
						}
						position++
						goto l1335
					l1336:
						position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
						if buffer[position] != rune('E') {
							goto l1330
						}
						position++
					}
				l1335:
					{
						position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1338
						}
						position++
						goto l1337
					l1338:
						position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
						if buffer[position] != rune('S') {
							goto l1330
						}
						position++
					}
				l1337:
					{
						position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1340
						}
						position++
						goto l1339
					l1340:
						position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
						if buffer[position] != rune('C') {
							goto l1330
						}
						position++
					}
				l1339:
					depth--
					add(rulePegText, position1332)
				}
				if !_rules[ruleAction90]() {
					goto l1330
				}
				depth--
				add(ruleDescending, position1331)
			}
			return true
		l1330:
			position, tokenIndex, depth = position1330, tokenIndex1330, depth1330
			return false
		},
		/* 118 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
			{
				position1342 := position
				depth++
				{
					position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1344
					}
					goto l1343
				l1344:
					position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					if !_rules[ruleInt]() {
						goto l1345
					}
					goto l1343
				l1345:
					position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					if !_rules[ruleFloat]() {
						goto l1346
					}
					goto l1343
				l1346:
					position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					if !_rules[ruleString]() {
						goto l1347
					}
					goto l1343
				l1347:
					position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					if !_rules[ruleBlob]() {
						goto l1348
					}
					goto l1343
				l1348:
					position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					if !_rules[ruleTimestamp]() {
						goto l1349
					}
					goto l1343
				l1349:
					position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					if !_rules[ruleArray]() {
						goto l1350
					}
					goto l1343
				l1350:
					position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					if !_rules[ruleMap]() {
						goto l1341
					}
				}
			l1343:
				depth--
				add(ruleType, position1342)
			}
			return true
		l1341:
			position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
			return false
		},
		/* 119 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action91)> */
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
						if buffer[position] != rune('b') {
							goto l1355
						}
						position++
						goto l1354
					l1355:
						position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
						if buffer[position] != rune('B') {
							goto l1351
						}
						position++
					}
				l1354:
					{
						position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1357
						}
						position++
						goto l1356
					l1357:
						position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
						if buffer[position] != rune('O') {
							goto l1351
						}
						position++
					}
				l1356:
					{
						position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1359
						}
						position++
						goto l1358
					l1359:
						position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
						if buffer[position] != rune('O') {
							goto l1351
						}
						position++
					}
				l1358:
					{
						position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1361
						}
						position++
						goto l1360
					l1361:
						position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
						if buffer[position] != rune('L') {
							goto l1351
						}
						position++
					}
				l1360:
					depth--
					add(rulePegText, position1353)
				}
				if !_rules[ruleAction91]() {
					goto l1351
				}
				depth--
				add(ruleBool, position1352)
			}
			return true
		l1351:
			position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
			return false
		},
		/* 120 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action92)> */
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
						if buffer[position] != rune('i') {
							goto l1366
						}
						position++
						goto l1365
					l1366:
						position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
						if buffer[position] != rune('I') {
							goto l1362
						}
						position++
					}
				l1365:
					{
						position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1368
						}
						position++
						goto l1367
					l1368:
						position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
						if buffer[position] != rune('N') {
							goto l1362
						}
						position++
					}
				l1367:
					{
						position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1370
						}
						position++
						goto l1369
					l1370:
						position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
						if buffer[position] != rune('T') {
							goto l1362
						}
						position++
					}
				l1369:
					depth--
					add(rulePegText, position1364)
				}
				if !_rules[ruleAction92]() {
					goto l1362
				}
				depth--
				add(ruleInt, position1363)
			}
			return true
		l1362:
			position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
			return false
		},
		/* 121 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action93)> */
		func() bool {
			position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
			{
				position1372 := position
				depth++
				{
					position1373 := position
					depth++
					{
						position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1375
						}
						position++
						goto l1374
					l1375:
						position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
						if buffer[position] != rune('F') {
							goto l1371
						}
						position++
					}
				l1374:
					{
						position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1377
						}
						position++
						goto l1376
					l1377:
						position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
						if buffer[position] != rune('L') {
							goto l1371
						}
						position++
					}
				l1376:
					{
						position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1379
						}
						position++
						goto l1378
					l1379:
						position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
						if buffer[position] != rune('O') {
							goto l1371
						}
						position++
					}
				l1378:
					{
						position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1381
						}
						position++
						goto l1380
					l1381:
						position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
						if buffer[position] != rune('A') {
							goto l1371
						}
						position++
					}
				l1380:
					{
						position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1383
						}
						position++
						goto l1382
					l1383:
						position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
						if buffer[position] != rune('T') {
							goto l1371
						}
						position++
					}
				l1382:
					depth--
					add(rulePegText, position1373)
				}
				if !_rules[ruleAction93]() {
					goto l1371
				}
				depth--
				add(ruleFloat, position1372)
			}
			return true
		l1371:
			position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
			return false
		},
		/* 122 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action94)> */
		func() bool {
			position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
			{
				position1385 := position
				depth++
				{
					position1386 := position
					depth++
					{
						position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1388
						}
						position++
						goto l1387
					l1388:
						position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
						if buffer[position] != rune('S') {
							goto l1384
						}
						position++
					}
				l1387:
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
							goto l1384
						}
						position++
					}
				l1389:
					{
						position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1392
						}
						position++
						goto l1391
					l1392:
						position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
						if buffer[position] != rune('R') {
							goto l1384
						}
						position++
					}
				l1391:
					{
						position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1394
						}
						position++
						goto l1393
					l1394:
						position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
						if buffer[position] != rune('I') {
							goto l1384
						}
						position++
					}
				l1393:
					{
						position1395, tokenIndex1395, depth1395 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1396
						}
						position++
						goto l1395
					l1396:
						position, tokenIndex, depth = position1395, tokenIndex1395, depth1395
						if buffer[position] != rune('N') {
							goto l1384
						}
						position++
					}
				l1395:
					{
						position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1398
						}
						position++
						goto l1397
					l1398:
						position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
						if buffer[position] != rune('G') {
							goto l1384
						}
						position++
					}
				l1397:
					depth--
					add(rulePegText, position1386)
				}
				if !_rules[ruleAction94]() {
					goto l1384
				}
				depth--
				add(ruleString, position1385)
			}
			return true
		l1384:
			position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
			return false
		},
		/* 123 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action95)> */
		func() bool {
			position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
			{
				position1400 := position
				depth++
				{
					position1401 := position
					depth++
					{
						position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1403
						}
						position++
						goto l1402
					l1403:
						position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
						if buffer[position] != rune('B') {
							goto l1399
						}
						position++
					}
				l1402:
					{
						position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1405
						}
						position++
						goto l1404
					l1405:
						position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
						if buffer[position] != rune('L') {
							goto l1399
						}
						position++
					}
				l1404:
					{
						position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1407
						}
						position++
						goto l1406
					l1407:
						position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
						if buffer[position] != rune('O') {
							goto l1399
						}
						position++
					}
				l1406:
					{
						position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1409
						}
						position++
						goto l1408
					l1409:
						position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
						if buffer[position] != rune('B') {
							goto l1399
						}
						position++
					}
				l1408:
					depth--
					add(rulePegText, position1401)
				}
				if !_rules[ruleAction95]() {
					goto l1399
				}
				depth--
				add(ruleBlob, position1400)
			}
			return true
		l1399:
			position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
			return false
		},
		/* 124 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action96)> */
		func() bool {
			position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
			{
				position1411 := position
				depth++
				{
					position1412 := position
					depth++
					{
						position1413, tokenIndex1413, depth1413 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1414
						}
						position++
						goto l1413
					l1414:
						position, tokenIndex, depth = position1413, tokenIndex1413, depth1413
						if buffer[position] != rune('T') {
							goto l1410
						}
						position++
					}
				l1413:
					{
						position1415, tokenIndex1415, depth1415 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1416
						}
						position++
						goto l1415
					l1416:
						position, tokenIndex, depth = position1415, tokenIndex1415, depth1415
						if buffer[position] != rune('I') {
							goto l1410
						}
						position++
					}
				l1415:
					{
						position1417, tokenIndex1417, depth1417 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1418
						}
						position++
						goto l1417
					l1418:
						position, tokenIndex, depth = position1417, tokenIndex1417, depth1417
						if buffer[position] != rune('M') {
							goto l1410
						}
						position++
					}
				l1417:
					{
						position1419, tokenIndex1419, depth1419 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1420
						}
						position++
						goto l1419
					l1420:
						position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
						if buffer[position] != rune('E') {
							goto l1410
						}
						position++
					}
				l1419:
					{
						position1421, tokenIndex1421, depth1421 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1422
						}
						position++
						goto l1421
					l1422:
						position, tokenIndex, depth = position1421, tokenIndex1421, depth1421
						if buffer[position] != rune('S') {
							goto l1410
						}
						position++
					}
				l1421:
					{
						position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1424
						}
						position++
						goto l1423
					l1424:
						position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
						if buffer[position] != rune('T') {
							goto l1410
						}
						position++
					}
				l1423:
					{
						position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1426
						}
						position++
						goto l1425
					l1426:
						position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
						if buffer[position] != rune('A') {
							goto l1410
						}
						position++
					}
				l1425:
					{
						position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1428
						}
						position++
						goto l1427
					l1428:
						position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
						if buffer[position] != rune('M') {
							goto l1410
						}
						position++
					}
				l1427:
					{
						position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1430
						}
						position++
						goto l1429
					l1430:
						position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
						if buffer[position] != rune('P') {
							goto l1410
						}
						position++
					}
				l1429:
					depth--
					add(rulePegText, position1412)
				}
				if !_rules[ruleAction96]() {
					goto l1410
				}
				depth--
				add(ruleTimestamp, position1411)
			}
			return true
		l1410:
			position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
			return false
		},
		/* 125 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action97)> */
		func() bool {
			position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
			{
				position1432 := position
				depth++
				{
					position1433 := position
					depth++
					{
						position1434, tokenIndex1434, depth1434 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1435
						}
						position++
						goto l1434
					l1435:
						position, tokenIndex, depth = position1434, tokenIndex1434, depth1434
						if buffer[position] != rune('A') {
							goto l1431
						}
						position++
					}
				l1434:
					{
						position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1437
						}
						position++
						goto l1436
					l1437:
						position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
						if buffer[position] != rune('R') {
							goto l1431
						}
						position++
					}
				l1436:
					{
						position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1439
						}
						position++
						goto l1438
					l1439:
						position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
						if buffer[position] != rune('R') {
							goto l1431
						}
						position++
					}
				l1438:
					{
						position1440, tokenIndex1440, depth1440 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1441
						}
						position++
						goto l1440
					l1441:
						position, tokenIndex, depth = position1440, tokenIndex1440, depth1440
						if buffer[position] != rune('A') {
							goto l1431
						}
						position++
					}
				l1440:
					{
						position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1443
						}
						position++
						goto l1442
					l1443:
						position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
						if buffer[position] != rune('Y') {
							goto l1431
						}
						position++
					}
				l1442:
					depth--
					add(rulePegText, position1433)
				}
				if !_rules[ruleAction97]() {
					goto l1431
				}
				depth--
				add(ruleArray, position1432)
			}
			return true
		l1431:
			position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
			return false
		},
		/* 126 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action98)> */
		func() bool {
			position1444, tokenIndex1444, depth1444 := position, tokenIndex, depth
			{
				position1445 := position
				depth++
				{
					position1446 := position
					depth++
					{
						position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1448
						}
						position++
						goto l1447
					l1448:
						position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
						if buffer[position] != rune('M') {
							goto l1444
						}
						position++
					}
				l1447:
					{
						position1449, tokenIndex1449, depth1449 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1450
						}
						position++
						goto l1449
					l1450:
						position, tokenIndex, depth = position1449, tokenIndex1449, depth1449
						if buffer[position] != rune('A') {
							goto l1444
						}
						position++
					}
				l1449:
					{
						position1451, tokenIndex1451, depth1451 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1452
						}
						position++
						goto l1451
					l1452:
						position, tokenIndex, depth = position1451, tokenIndex1451, depth1451
						if buffer[position] != rune('P') {
							goto l1444
						}
						position++
					}
				l1451:
					depth--
					add(rulePegText, position1446)
				}
				if !_rules[ruleAction98]() {
					goto l1444
				}
				depth--
				add(ruleMap, position1445)
			}
			return true
		l1444:
			position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
			return false
		},
		/* 127 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action99)> */
		func() bool {
			position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
			{
				position1454 := position
				depth++
				{
					position1455 := position
					depth++
					{
						position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1457
						}
						position++
						goto l1456
					l1457:
						position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
						if buffer[position] != rune('O') {
							goto l1453
						}
						position++
					}
				l1456:
					{
						position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1459
						}
						position++
						goto l1458
					l1459:
						position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
						if buffer[position] != rune('R') {
							goto l1453
						}
						position++
					}
				l1458:
					depth--
					add(rulePegText, position1455)
				}
				if !_rules[ruleAction99]() {
					goto l1453
				}
				depth--
				add(ruleOr, position1454)
			}
			return true
		l1453:
			position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
			return false
		},
		/* 128 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action100)> */
		func() bool {
			position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
			{
				position1461 := position
				depth++
				{
					position1462 := position
					depth++
					{
						position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1464
						}
						position++
						goto l1463
					l1464:
						position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
						if buffer[position] != rune('A') {
							goto l1460
						}
						position++
					}
				l1463:
					{
						position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1466
						}
						position++
						goto l1465
					l1466:
						position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
						if buffer[position] != rune('N') {
							goto l1460
						}
						position++
					}
				l1465:
					{
						position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1468
						}
						position++
						goto l1467
					l1468:
						position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
						if buffer[position] != rune('D') {
							goto l1460
						}
						position++
					}
				l1467:
					depth--
					add(rulePegText, position1462)
				}
				if !_rules[ruleAction100]() {
					goto l1460
				}
				depth--
				add(ruleAnd, position1461)
			}
			return true
		l1460:
			position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
			return false
		},
		/* 129 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action101)> */
		func() bool {
			position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
			{
				position1470 := position
				depth++
				{
					position1471 := position
					depth++
					{
						position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1473
						}
						position++
						goto l1472
					l1473:
						position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
						if buffer[position] != rune('N') {
							goto l1469
						}
						position++
					}
				l1472:
					{
						position1474, tokenIndex1474, depth1474 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1475
						}
						position++
						goto l1474
					l1475:
						position, tokenIndex, depth = position1474, tokenIndex1474, depth1474
						if buffer[position] != rune('O') {
							goto l1469
						}
						position++
					}
				l1474:
					{
						position1476, tokenIndex1476, depth1476 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1477
						}
						position++
						goto l1476
					l1477:
						position, tokenIndex, depth = position1476, tokenIndex1476, depth1476
						if buffer[position] != rune('T') {
							goto l1469
						}
						position++
					}
				l1476:
					depth--
					add(rulePegText, position1471)
				}
				if !_rules[ruleAction101]() {
					goto l1469
				}
				depth--
				add(ruleNot, position1470)
			}
			return true
		l1469:
			position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
			return false
		},
		/* 130 Equal <- <(<'='> Action102)> */
		func() bool {
			position1478, tokenIndex1478, depth1478 := position, tokenIndex, depth
			{
				position1479 := position
				depth++
				{
					position1480 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1478
					}
					position++
					depth--
					add(rulePegText, position1480)
				}
				if !_rules[ruleAction102]() {
					goto l1478
				}
				depth--
				add(ruleEqual, position1479)
			}
			return true
		l1478:
			position, tokenIndex, depth = position1478, tokenIndex1478, depth1478
			return false
		},
		/* 131 Less <- <(<'<'> Action103)> */
		func() bool {
			position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
			{
				position1482 := position
				depth++
				{
					position1483 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1481
					}
					position++
					depth--
					add(rulePegText, position1483)
				}
				if !_rules[ruleAction103]() {
					goto l1481
				}
				depth--
				add(ruleLess, position1482)
			}
			return true
		l1481:
			position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
			return false
		},
		/* 132 LessOrEqual <- <(<('<' '=')> Action104)> */
		func() bool {
			position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
			{
				position1485 := position
				depth++
				{
					position1486 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1484
					}
					position++
					if buffer[position] != rune('=') {
						goto l1484
					}
					position++
					depth--
					add(rulePegText, position1486)
				}
				if !_rules[ruleAction104]() {
					goto l1484
				}
				depth--
				add(ruleLessOrEqual, position1485)
			}
			return true
		l1484:
			position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
			return false
		},
		/* 133 Greater <- <(<'>'> Action105)> */
		func() bool {
			position1487, tokenIndex1487, depth1487 := position, tokenIndex, depth
			{
				position1488 := position
				depth++
				{
					position1489 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1487
					}
					position++
					depth--
					add(rulePegText, position1489)
				}
				if !_rules[ruleAction105]() {
					goto l1487
				}
				depth--
				add(ruleGreater, position1488)
			}
			return true
		l1487:
			position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
			return false
		},
		/* 134 GreaterOrEqual <- <(<('>' '=')> Action106)> */
		func() bool {
			position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
			{
				position1491 := position
				depth++
				{
					position1492 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1490
					}
					position++
					if buffer[position] != rune('=') {
						goto l1490
					}
					position++
					depth--
					add(rulePegText, position1492)
				}
				if !_rules[ruleAction106]() {
					goto l1490
				}
				depth--
				add(ruleGreaterOrEqual, position1491)
			}
			return true
		l1490:
			position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
			return false
		},
		/* 135 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action107)> */
		func() bool {
			position1493, tokenIndex1493, depth1493 := position, tokenIndex, depth
			{
				position1494 := position
				depth++
				{
					position1495 := position
					depth++
					{
						position1496, tokenIndex1496, depth1496 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1497
						}
						position++
						if buffer[position] != rune('=') {
							goto l1497
						}
						position++
						goto l1496
					l1497:
						position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
						if buffer[position] != rune('<') {
							goto l1493
						}
						position++
						if buffer[position] != rune('>') {
							goto l1493
						}
						position++
					}
				l1496:
					depth--
					add(rulePegText, position1495)
				}
				if !_rules[ruleAction107]() {
					goto l1493
				}
				depth--
				add(ruleNotEqual, position1494)
			}
			return true
		l1493:
			position, tokenIndex, depth = position1493, tokenIndex1493, depth1493
			return false
		},
		/* 136 Concat <- <(<('|' '|')> Action108)> */
		func() bool {
			position1498, tokenIndex1498, depth1498 := position, tokenIndex, depth
			{
				position1499 := position
				depth++
				{
					position1500 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1498
					}
					position++
					if buffer[position] != rune('|') {
						goto l1498
					}
					position++
					depth--
					add(rulePegText, position1500)
				}
				if !_rules[ruleAction108]() {
					goto l1498
				}
				depth--
				add(ruleConcat, position1499)
			}
			return true
		l1498:
			position, tokenIndex, depth = position1498, tokenIndex1498, depth1498
			return false
		},
		/* 137 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action109)> */
		func() bool {
			position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
			{
				position1502 := position
				depth++
				{
					position1503 := position
					depth++
					{
						position1504, tokenIndex1504, depth1504 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1505
						}
						position++
						goto l1504
					l1505:
						position, tokenIndex, depth = position1504, tokenIndex1504, depth1504
						if buffer[position] != rune('I') {
							goto l1501
						}
						position++
					}
				l1504:
					{
						position1506, tokenIndex1506, depth1506 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1507
						}
						position++
						goto l1506
					l1507:
						position, tokenIndex, depth = position1506, tokenIndex1506, depth1506
						if buffer[position] != rune('S') {
							goto l1501
						}
						position++
					}
				l1506:
					depth--
					add(rulePegText, position1503)
				}
				if !_rules[ruleAction109]() {
					goto l1501
				}
				depth--
				add(ruleIs, position1502)
			}
			return true
		l1501:
			position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
			return false
		},
		/* 138 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action110)> */
		func() bool {
			position1508, tokenIndex1508, depth1508 := position, tokenIndex, depth
			{
				position1509 := position
				depth++
				{
					position1510 := position
					depth++
					{
						position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1512
						}
						position++
						goto l1511
					l1512:
						position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
						if buffer[position] != rune('I') {
							goto l1508
						}
						position++
					}
				l1511:
					{
						position1513, tokenIndex1513, depth1513 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1514
						}
						position++
						goto l1513
					l1514:
						position, tokenIndex, depth = position1513, tokenIndex1513, depth1513
						if buffer[position] != rune('S') {
							goto l1508
						}
						position++
					}
				l1513:
					if !_rules[rulesp]() {
						goto l1508
					}
					{
						position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1516
						}
						position++
						goto l1515
					l1516:
						position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
						if buffer[position] != rune('N') {
							goto l1508
						}
						position++
					}
				l1515:
					{
						position1517, tokenIndex1517, depth1517 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1518
						}
						position++
						goto l1517
					l1518:
						position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
						if buffer[position] != rune('O') {
							goto l1508
						}
						position++
					}
				l1517:
					{
						position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1520
						}
						position++
						goto l1519
					l1520:
						position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
						if buffer[position] != rune('T') {
							goto l1508
						}
						position++
					}
				l1519:
					depth--
					add(rulePegText, position1510)
				}
				if !_rules[ruleAction110]() {
					goto l1508
				}
				depth--
				add(ruleIsNot, position1509)
			}
			return true
		l1508:
			position, tokenIndex, depth = position1508, tokenIndex1508, depth1508
			return false
		},
		/* 139 Plus <- <(<'+'> Action111)> */
		func() bool {
			position1521, tokenIndex1521, depth1521 := position, tokenIndex, depth
			{
				position1522 := position
				depth++
				{
					position1523 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1521
					}
					position++
					depth--
					add(rulePegText, position1523)
				}
				if !_rules[ruleAction111]() {
					goto l1521
				}
				depth--
				add(rulePlus, position1522)
			}
			return true
		l1521:
			position, tokenIndex, depth = position1521, tokenIndex1521, depth1521
			return false
		},
		/* 140 Minus <- <(<'-'> Action112)> */
		func() bool {
			position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
			{
				position1525 := position
				depth++
				{
					position1526 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1524
					}
					position++
					depth--
					add(rulePegText, position1526)
				}
				if !_rules[ruleAction112]() {
					goto l1524
				}
				depth--
				add(ruleMinus, position1525)
			}
			return true
		l1524:
			position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
			return false
		},
		/* 141 Multiply <- <(<'*'> Action113)> */
		func() bool {
			position1527, tokenIndex1527, depth1527 := position, tokenIndex, depth
			{
				position1528 := position
				depth++
				{
					position1529 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1527
					}
					position++
					depth--
					add(rulePegText, position1529)
				}
				if !_rules[ruleAction113]() {
					goto l1527
				}
				depth--
				add(ruleMultiply, position1528)
			}
			return true
		l1527:
			position, tokenIndex, depth = position1527, tokenIndex1527, depth1527
			return false
		},
		/* 142 Divide <- <(<'/'> Action114)> */
		func() bool {
			position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
			{
				position1531 := position
				depth++
				{
					position1532 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1530
					}
					position++
					depth--
					add(rulePegText, position1532)
				}
				if !_rules[ruleAction114]() {
					goto l1530
				}
				depth--
				add(ruleDivide, position1531)
			}
			return true
		l1530:
			position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
			return false
		},
		/* 143 Modulo <- <(<'%'> Action115)> */
		func() bool {
			position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
			{
				position1534 := position
				depth++
				{
					position1535 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1533
					}
					position++
					depth--
					add(rulePegText, position1535)
				}
				if !_rules[ruleAction115]() {
					goto l1533
				}
				depth--
				add(ruleModulo, position1534)
			}
			return true
		l1533:
			position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
			return false
		},
		/* 144 UnaryMinus <- <(<'-'> Action116)> */
		func() bool {
			position1536, tokenIndex1536, depth1536 := position, tokenIndex, depth
			{
				position1537 := position
				depth++
				{
					position1538 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1536
					}
					position++
					depth--
					add(rulePegText, position1538)
				}
				if !_rules[ruleAction116]() {
					goto l1536
				}
				depth--
				add(ruleUnaryMinus, position1537)
			}
			return true
		l1536:
			position, tokenIndex, depth = position1536, tokenIndex1536, depth1536
			return false
		},
		/* 145 Identifier <- <(<ident> Action117)> */
		func() bool {
			position1539, tokenIndex1539, depth1539 := position, tokenIndex, depth
			{
				position1540 := position
				depth++
				{
					position1541 := position
					depth++
					if !_rules[ruleident]() {
						goto l1539
					}
					depth--
					add(rulePegText, position1541)
				}
				if !_rules[ruleAction117]() {
					goto l1539
				}
				depth--
				add(ruleIdentifier, position1540)
			}
			return true
		l1539:
			position, tokenIndex, depth = position1539, tokenIndex1539, depth1539
			return false
		},
		/* 146 TargetIdentifier <- <(<jsonSetPath> Action118)> */
		func() bool {
			position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
			{
				position1543 := position
				depth++
				{
					position1544 := position
					depth++
					if !_rules[rulejsonSetPath]() {
						goto l1542
					}
					depth--
					add(rulePegText, position1544)
				}
				if !_rules[ruleAction118]() {
					goto l1542
				}
				depth--
				add(ruleTargetIdentifier, position1543)
			}
			return true
		l1542:
			position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
			return false
		},
		/* 147 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1545, tokenIndex1545, depth1545 := position, tokenIndex, depth
			{
				position1546 := position
				depth++
				{
					position1547, tokenIndex1547, depth1547 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1548
					}
					position++
					goto l1547
				l1548:
					position, tokenIndex, depth = position1547, tokenIndex1547, depth1547
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1545
					}
					position++
				}
			l1547:
			l1549:
				{
					position1550, tokenIndex1550, depth1550 := position, tokenIndex, depth
					{
						position1551, tokenIndex1551, depth1551 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1552
						}
						position++
						goto l1551
					l1552:
						position, tokenIndex, depth = position1551, tokenIndex1551, depth1551
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1553
						}
						position++
						goto l1551
					l1553:
						position, tokenIndex, depth = position1551, tokenIndex1551, depth1551
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1554
						}
						position++
						goto l1551
					l1554:
						position, tokenIndex, depth = position1551, tokenIndex1551, depth1551
						if buffer[position] != rune('_') {
							goto l1550
						}
						position++
					}
				l1551:
					goto l1549
				l1550:
					position, tokenIndex, depth = position1550, tokenIndex1550, depth1550
				}
				depth--
				add(ruleident, position1546)
			}
			return true
		l1545:
			position, tokenIndex, depth = position1545, tokenIndex1545, depth1545
			return false
		},
		/* 148 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1555, tokenIndex1555, depth1555 := position, tokenIndex, depth
			{
				position1556 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1555
				}
			l1557:
				{
					position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1558
					}
					goto l1557
				l1558:
					position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
				}
				depth--
				add(rulejsonGetPath, position1556)
			}
			return true
		l1555:
			position, tokenIndex, depth = position1555, tokenIndex1555, depth1555
			return false
		},
		/* 149 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1559, tokenIndex1559, depth1559 := position, tokenIndex, depth
			{
				position1560 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1559
				}
			l1561:
				{
					position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1562
					}
					goto l1561
				l1562:
					position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
				}
				depth--
				add(rulejsonSetPath, position1560)
			}
			return true
		l1559:
			position, tokenIndex, depth = position1559, tokenIndex1559, depth1559
			return false
		},
		/* 150 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1563, tokenIndex1563, depth1563 := position, tokenIndex, depth
			{
				position1564 := position
				depth++
				{
					position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1566
					}
					goto l1565
				l1566:
					position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1563
					}
				}
			l1565:
				depth--
				add(rulejsonPathHead, position1564)
			}
			return true
		l1563:
			position, tokenIndex, depth = position1563, tokenIndex1563, depth1563
			return false
		},
		/* 151 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1567, tokenIndex1567, depth1567 := position, tokenIndex, depth
			{
				position1568 := position
				depth++
				{
					position1569, tokenIndex1569, depth1569 := position, tokenIndex, depth
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1570
					}
					goto l1569
				l1570:
					position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1571
					}
					goto l1569
				l1571:
					position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1572
					}
					goto l1569
				l1572:
					position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1573
					}
					goto l1569
				l1573:
					position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
					if !_rules[rulejsonArraySlice]() {
						goto l1574
					}
					goto l1569
				l1574:
					position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
					if !_rules[rulejsonArrayAccess]() {
						goto l1567
					}
				}
			l1569:
				depth--
				add(rulejsonGetPathNonHead, position1568)
			}
			return true
		l1567:
			position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
			return false
		},
		/* 152 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1575, tokenIndex1575, depth1575 := position, tokenIndex, depth
			{
				position1576 := position
				depth++
				{
					position1577, tokenIndex1577, depth1577 := position, tokenIndex, depth
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1578
					}
					goto l1577
				l1578:
					position, tokenIndex, depth = position1577, tokenIndex1577, depth1577
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1575
					}
				}
			l1577:
				depth--
				add(rulejsonSetPathNonHead, position1576)
			}
			return true
		l1575:
			position, tokenIndex, depth = position1575, tokenIndex1575, depth1575
			return false
		},
		/* 153 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1579, tokenIndex1579, depth1579 := position, tokenIndex, depth
			{
				position1580 := position
				depth++
				{
					position1581, tokenIndex1581, depth1581 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1582
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1582
					}
					goto l1581
				l1582:
					position, tokenIndex, depth = position1581, tokenIndex1581, depth1581
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1579
					}
				}
			l1581:
				depth--
				add(rulejsonMapSingleLevel, position1580)
			}
			return true
		l1579:
			position, tokenIndex, depth = position1579, tokenIndex1579, depth1579
			return false
		},
		/* 154 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1583, tokenIndex1583, depth1583 := position, tokenIndex, depth
			{
				position1584 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1583
				}
				position++
				if buffer[position] != rune('.') {
					goto l1583
				}
				position++
				{
					position1585, tokenIndex1585, depth1585 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1586
					}
					goto l1585
				l1586:
					position, tokenIndex, depth = position1585, tokenIndex1585, depth1585
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1583
					}
				}
			l1585:
				depth--
				add(rulejsonMapMultipleLevel, position1584)
			}
			return true
		l1583:
			position, tokenIndex, depth = position1583, tokenIndex1583, depth1583
			return false
		},
		/* 155 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1587, tokenIndex1587, depth1587 := position, tokenIndex, depth
			{
				position1588 := position
				depth++
				{
					position1589 := position
					depth++
					{
						position1590, tokenIndex1590, depth1590 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1591
						}
						position++
						goto l1590
					l1591:
						position, tokenIndex, depth = position1590, tokenIndex1590, depth1590
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1587
						}
						position++
					}
				l1590:
				l1592:
					{
						position1593, tokenIndex1593, depth1593 := position, tokenIndex, depth
						{
							position1594, tokenIndex1594, depth1594 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1595
							}
							position++
							goto l1594
						l1595:
							position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1596
							}
							position++
							goto l1594
						l1596:
							position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1597
							}
							position++
							goto l1594
						l1597:
							position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
							if buffer[position] != rune('_') {
								goto l1593
							}
							position++
						}
					l1594:
						goto l1592
					l1593:
						position, tokenIndex, depth = position1593, tokenIndex1593, depth1593
					}
					depth--
					add(rulePegText, position1589)
				}
				depth--
				add(rulejsonMapAccessString, position1588)
			}
			return true
		l1587:
			position, tokenIndex, depth = position1587, tokenIndex1587, depth1587
			return false
		},
		/* 156 jsonMapAccessBracket <- <('[' singleQuotedString ']')> */
		func() bool {
			position1598, tokenIndex1598, depth1598 := position, tokenIndex, depth
			{
				position1599 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1598
				}
				position++
				if !_rules[rulesingleQuotedString]() {
					goto l1598
				}
				if buffer[position] != rune(']') {
					goto l1598
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1599)
			}
			return true
		l1598:
			position, tokenIndex, depth = position1598, tokenIndex1598, depth1598
			return false
		},
		/* 157 singleQuotedString <- <('\'' <(('\'' '\'') / (!'\'' .))*> '\'')> */
		func() bool {
			position1600, tokenIndex1600, depth1600 := position, tokenIndex, depth
			{
				position1601 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l1600
				}
				position++
				{
					position1602 := position
					depth++
				l1603:
					{
						position1604, tokenIndex1604, depth1604 := position, tokenIndex, depth
						{
							position1605, tokenIndex1605, depth1605 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1606
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1606
							}
							position++
							goto l1605
						l1606:
							position, tokenIndex, depth = position1605, tokenIndex1605, depth1605
							{
								position1607, tokenIndex1607, depth1607 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1607
								}
								position++
								goto l1604
							l1607:
								position, tokenIndex, depth = position1607, tokenIndex1607, depth1607
							}
							if !matchDot() {
								goto l1604
							}
						}
					l1605:
						goto l1603
					l1604:
						position, tokenIndex, depth = position1604, tokenIndex1604, depth1604
					}
					depth--
					add(rulePegText, position1602)
				}
				if buffer[position] != rune('\'') {
					goto l1600
				}
				position++
				depth--
				add(rulesingleQuotedString, position1601)
			}
			return true
		l1600:
			position, tokenIndex, depth = position1600, tokenIndex1600, depth1600
			return false
		},
		/* 158 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1608, tokenIndex1608, depth1608 := position, tokenIndex, depth
			{
				position1609 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1608
				}
				position++
				{
					position1610 := position
					depth++
					{
						position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1611
						}
						position++
						goto l1612
					l1611:
						position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
					}
				l1612:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1608
					}
					position++
				l1613:
					{
						position1614, tokenIndex1614, depth1614 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1614
						}
						position++
						goto l1613
					l1614:
						position, tokenIndex, depth = position1614, tokenIndex1614, depth1614
					}
					depth--
					add(rulePegText, position1610)
				}
				if buffer[position] != rune(']') {
					goto l1608
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1609)
			}
			return true
		l1608:
			position, tokenIndex, depth = position1608, tokenIndex1608, depth1608
			return false
		},
		/* 159 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
		func() bool {
			position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
			{
				position1616 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1615
				}
				position++
				{
					position1617 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1615
					}
					position++
				l1618:
					{
						position1619, tokenIndex1619, depth1619 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1619
						}
						position++
						goto l1618
					l1619:
						position, tokenIndex, depth = position1619, tokenIndex1619, depth1619
					}
					depth--
					add(rulePegText, position1617)
				}
				if buffer[position] != rune(']') {
					goto l1615
				}
				position++
				depth--
				add(rulejsonNonNegativeArrayAccess, position1616)
			}
			return true
		l1615:
			position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
			return false
		},
		/* 160 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1620, tokenIndex1620, depth1620 := position, tokenIndex, depth
			{
				position1621 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1620
				}
				position++
				{
					position1622 := position
					depth++
					{
						position1623, tokenIndex1623, depth1623 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1623
						}
						position++
						goto l1624
					l1623:
						position, tokenIndex, depth = position1623, tokenIndex1623, depth1623
					}
				l1624:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1620
					}
					position++
				l1625:
					{
						position1626, tokenIndex1626, depth1626 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1626
						}
						position++
						goto l1625
					l1626:
						position, tokenIndex, depth = position1626, tokenIndex1626, depth1626
					}
					if buffer[position] != rune(':') {
						goto l1620
					}
					position++
					{
						position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1627
						}
						position++
						goto l1628
					l1627:
						position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
					}
				l1628:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1620
					}
					position++
				l1629:
					{
						position1630, tokenIndex1630, depth1630 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1630
						}
						position++
						goto l1629
					l1630:
						position, tokenIndex, depth = position1630, tokenIndex1630, depth1630
					}
					{
						position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1631
						}
						position++
						{
							position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1633
							}
							position++
							goto l1634
						l1633:
							position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
						}
					l1634:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1631
						}
						position++
					l1635:
						{
							position1636, tokenIndex1636, depth1636 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1636
							}
							position++
							goto l1635
						l1636:
							position, tokenIndex, depth = position1636, tokenIndex1636, depth1636
						}
						goto l1632
					l1631:
						position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
					}
				l1632:
					depth--
					add(rulePegText, position1622)
				}
				if buffer[position] != rune(']') {
					goto l1620
				}
				position++
				depth--
				add(rulejsonArraySlice, position1621)
			}
			return true
		l1620:
			position, tokenIndex, depth = position1620, tokenIndex1620, depth1620
			return false
		},
		/* 161 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1637, tokenIndex1637, depth1637 := position, tokenIndex, depth
			{
				position1638 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1637
				}
				position++
				{
					position1639 := position
					depth++
					{
						position1640, tokenIndex1640, depth1640 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1641
						}
						position++
						{
							position1642, tokenIndex1642, depth1642 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1642
							}
							position++
							goto l1643
						l1642:
							position, tokenIndex, depth = position1642, tokenIndex1642, depth1642
						}
					l1643:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1641
						}
						position++
					l1644:
						{
							position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1645
							}
							position++
							goto l1644
						l1645:
							position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
						}
						goto l1640
					l1641:
						position, tokenIndex, depth = position1640, tokenIndex1640, depth1640
						{
							position1646, tokenIndex1646, depth1646 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1646
							}
							position++
							goto l1647
						l1646:
							position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
						}
					l1647:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1637
						}
						position++
					l1648:
						{
							position1649, tokenIndex1649, depth1649 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1649
							}
							position++
							goto l1648
						l1649:
							position, tokenIndex, depth = position1649, tokenIndex1649, depth1649
						}
						if buffer[position] != rune(':') {
							goto l1637
						}
						position++
					}
				l1640:
					depth--
					add(rulePegText, position1639)
				}
				if buffer[position] != rune(']') {
					goto l1637
				}
				position++
				depth--
				add(rulejsonArrayPartialSlice, position1638)
			}
			return true
		l1637:
			position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
			return false
		},
		/* 162 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1650, tokenIndex1650, depth1650 := position, tokenIndex, depth
			{
				position1651 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1650
				}
				position++
				if buffer[position] != rune(':') {
					goto l1650
				}
				position++
				if buffer[position] != rune(']') {
					goto l1650
				}
				position++
				depth--
				add(rulejsonArrayFullSlice, position1651)
			}
			return true
		l1650:
			position, tokenIndex, depth = position1650, tokenIndex1650, depth1650
			return false
		},
		/* 163 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
			{
				position1653 := position
				depth++
				{
					position1654, tokenIndex1654, depth1654 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l1655
					}
					position++
					goto l1654
				l1655:
					position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
					if buffer[position] != rune('\t') {
						goto l1656
					}
					position++
					goto l1654
				l1656:
					position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
					if buffer[position] != rune('\n') {
						goto l1657
					}
					position++
					goto l1654
				l1657:
					position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
					if buffer[position] != rune('\r') {
						goto l1658
					}
					position++
					goto l1654
				l1658:
					position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
					if !_rules[rulecomment]() {
						goto l1659
					}
					goto l1654
				l1659:
					position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
					if !_rules[rulefinalComment]() {
						goto l1652
					}
				}
			l1654:
				depth--
				add(rulespElem, position1653)
			}
			return true
		l1652:
			position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
			return false
		},
		/* 164 sp <- <spElem+> */
		func() bool {
			position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
			{
				position1661 := position
				depth++
				if !_rules[rulespElem]() {
					goto l1660
				}
			l1662:
				{
					position1663, tokenIndex1663, depth1663 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1663
					}
					goto l1662
				l1663:
					position, tokenIndex, depth = position1663, tokenIndex1663, depth1663
				}
				depth--
				add(rulesp, position1661)
			}
			return true
		l1660:
			position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
			return false
		},
		/* 165 spOpt <- <spElem*> */
		func() bool {
			{
				position1665 := position
				depth++
			l1666:
				{
					position1667, tokenIndex1667, depth1667 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1667
					}
					goto l1666
				l1667:
					position, tokenIndex, depth = position1667, tokenIndex1667, depth1667
				}
				depth--
				add(rulespOpt, position1665)
			}
			return true
		},
		/* 166 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1668, tokenIndex1668, depth1668 := position, tokenIndex, depth
			{
				position1669 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1668
				}
				position++
				if buffer[position] != rune('-') {
					goto l1668
				}
				position++
			l1670:
				{
					position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
					{
						position1672, tokenIndex1672, depth1672 := position, tokenIndex, depth
						{
							position1673, tokenIndex1673, depth1673 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1674
							}
							position++
							goto l1673
						l1674:
							position, tokenIndex, depth = position1673, tokenIndex1673, depth1673
							if buffer[position] != rune('\n') {
								goto l1672
							}
							position++
						}
					l1673:
						goto l1671
					l1672:
						position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
					}
					if !matchDot() {
						goto l1671
					}
					goto l1670
				l1671:
					position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
				}
				{
					position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1676
					}
					position++
					goto l1675
				l1676:
					position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
					if buffer[position] != rune('\n') {
						goto l1668
					}
					position++
				}
			l1675:
				depth--
				add(rulecomment, position1669)
			}
			return true
		l1668:
			position, tokenIndex, depth = position1668, tokenIndex1668, depth1668
			return false
		},
		/* 167 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
			{
				position1678 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1677
				}
				position++
				if buffer[position] != rune('-') {
					goto l1677
				}
				position++
			l1679:
				{
					position1680, tokenIndex1680, depth1680 := position, tokenIndex, depth
					{
						position1681, tokenIndex1681, depth1681 := position, tokenIndex, depth
						{
							position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1683
							}
							position++
							goto l1682
						l1683:
							position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
							if buffer[position] != rune('\n') {
								goto l1681
							}
							position++
						}
					l1682:
						goto l1680
					l1681:
						position, tokenIndex, depth = position1681, tokenIndex1681, depth1681
					}
					if !matchDot() {
						goto l1680
					}
					goto l1679
				l1680:
					position, tokenIndex, depth = position1680, tokenIndex1680, depth1680
				}
				{
					position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
					if !matchDot() {
						goto l1684
					}
					goto l1677
				l1684:
					position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
				}
				depth--
				add(rulefinalComment, position1678)
			}
			return true
		l1677:
			position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
			return false
		},
		nil,
		/* 170 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 171 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 172 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 173 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 174 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 175 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 176 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 177 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 178 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 179 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 180 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 181 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 182 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 183 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 184 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 185 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 186 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 187 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 188 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 189 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 190 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 191 Action21 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 192 Action22 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 193 Action23 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 194 Action24 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 195 Action25 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 196 Action26 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 197 Action27 <- <{
		    p.AssembleEmitterSampling(CountBasedSampling)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 198 Action28 <- <{
		    p.AssembleEmitterSampling(RandomizedSampling)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 199 Action29 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 200 Action30 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 201 Action31 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 202 Action32 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 203 Action33 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 204 Action34 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 205 Action35 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 206 Action36 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 207 Action37 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 208 Action38 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 209 Action39 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 210 Action40 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 211 Action41 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 212 Action42 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 213 Action43 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 214 Action44 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 215 Action45 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 216 Action46 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 217 Action47 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 218 Action48 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 219 Action49 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 220 Action50 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 221 Action51 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 222 Action52 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 223 Action53 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 224 Action54 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 225 Action55 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 226 Action56 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 227 Action57 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 228 Action58 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 229 Action59 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 230 Action60 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 231 Action61 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 232 Action62 <- <{
		    p.AssembleSortedExpression()
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 233 Action63 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 234 Action64 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 235 Action65 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 236 Action66 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 237 Action67 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 238 Action68 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 239 Action69 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 240 Action70 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 241 Action71 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 242 Action72 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 243 Action73 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 244 Action74 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 245 Action75 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 246 Action76 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 247 Action77 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 248 Action78 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 249 Action79 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 250 Action80 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 251 Action81 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 252 Action82 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 253 Action83 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 254 Action84 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 255 Action85 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 256 Action86 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 257 Action87 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 258 Action88 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 259 Action89 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 260 Action90 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 261 Action91 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 262 Action92 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 263 Action93 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 264 Action94 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 265 Action95 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 266 Action96 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 267 Action97 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 268 Action98 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 269 Action99 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 270 Action100 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 271 Action101 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 272 Action102 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 273 Action103 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 274 Action104 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 275 Action105 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 276 Action106 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 277 Action107 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 278 Action108 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 279 Action109 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 280 Action110 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction110, position)
			}
			return true
		},
		/* 281 Action111 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction111, position)
			}
			return true
		},
		/* 282 Action112 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction112, position)
			}
			return true
		},
		/* 283 Action113 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction113, position)
			}
			return true
		},
		/* 284 Action114 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction114, position)
			}
			return true
		},
		/* 285 Action115 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction115, position)
			}
			return true
		},
		/* 286 Action116 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction116, position)
			}
			return true
		},
		/* 287 Action117 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction117, position)
			}
			return true
		},
		/* 288 Action118 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction118, position)
			}
			return true
		},
	}
	p.rules = _rules
}

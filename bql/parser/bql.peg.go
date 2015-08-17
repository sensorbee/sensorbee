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
	ruleEmitterLimit
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
	ruleFuncParams
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
	rulejsonPathHead
	rulejsonPathNonHead
	rulejsonMapAccessString
	rulejsonMapAccessBracket
	rulejsonArrayAccess
	rulesp
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
	"EmitterLimit",
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
	"FuncParams",
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
	"jsonPathHead",
	"jsonPathNonHead",
	"jsonMapAccessString",
	"jsonMapAccessBracket",
	"jsonArrayAccess",
	"sp",
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
	rules  [259]func() bool
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

			p.AssembleProjections(begin, end)

		case ruleAction28:

			p.AssembleAlias()

		case ruleAction29:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction30:

			p.AssembleInterval()

		case ruleAction31:

			p.AssembleInterval()

		case ruleAction32:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction33:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction34:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction35:

			p.EnsureAliasedStreamWindow()

		case ruleAction36:

			p.AssembleAliasedStreamWindow()

		case ruleAction37:

			p.AssembleStreamWindow()

		case ruleAction38:

			p.AssembleUDSFFuncApp()

		case ruleAction39:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction40:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction41:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction42:

			p.AssembleSourceSinkParam()

		case ruleAction43:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction44:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction45:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction46:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction47:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction48:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction49:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction50:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction51:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction52:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction53:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction54:

			p.AssembleTypeCast(begin, end)

		case ruleAction55:

			p.AssembleTypeCast(begin, end)

		case ruleAction56:

			p.AssembleFuncApp()

		case ruleAction57:

			p.AssembleExpressions(begin, end)

		case ruleAction58:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction59:

			p.AssembleMap(begin, end)

		case ruleAction60:

			p.AssembleKeyValuePair()

		case ruleAction61:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction62:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction63:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction64:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction65:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction66:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction67:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction68:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction69:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction70:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction71:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction72:

			p.PushComponent(begin, end, Istream)

		case ruleAction73:

			p.PushComponent(begin, end, Dstream)

		case ruleAction74:

			p.PushComponent(begin, end, Rstream)

		case ruleAction75:

			p.PushComponent(begin, end, Tuples)

		case ruleAction76:

			p.PushComponent(begin, end, Seconds)

		case ruleAction77:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction78:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction80:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction81:

			p.PushComponent(begin, end, Yes)

		case ruleAction82:

			p.PushComponent(begin, end, No)

		case ruleAction83:

			p.PushComponent(begin, end, Bool)

		case ruleAction84:

			p.PushComponent(begin, end, Int)

		case ruleAction85:

			p.PushComponent(begin, end, Float)

		case ruleAction86:

			p.PushComponent(begin, end, String)

		case ruleAction87:

			p.PushComponent(begin, end, Blob)

		case ruleAction88:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction89:

			p.PushComponent(begin, end, Array)

		case ruleAction90:

			p.PushComponent(begin, end, Map)

		case ruleAction91:

			p.PushComponent(begin, end, Or)

		case ruleAction92:

			p.PushComponent(begin, end, And)

		case ruleAction93:

			p.PushComponent(begin, end, Not)

		case ruleAction94:

			p.PushComponent(begin, end, Equal)

		case ruleAction95:

			p.PushComponent(begin, end, Less)

		case ruleAction96:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction97:

			p.PushComponent(begin, end, Greater)

		case ruleAction98:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction99:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction100:

			p.PushComponent(begin, end, Concat)

		case ruleAction101:

			p.PushComponent(begin, end, Is)

		case ruleAction102:

			p.PushComponent(begin, end, IsNot)

		case ruleAction103:

			p.PushComponent(begin, end, Plus)

		case ruleAction104:

			p.PushComponent(begin, end, Minus)

		case ruleAction105:

			p.PushComponent(begin, end, Multiply)

		case ruleAction106:

			p.PushComponent(begin, end, Divide)

		case ruleAction107:

			p.PushComponent(begin, end, Modulo)

		case ruleAction108:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction109:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction110:

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
		/* 0 SingleStatement <- <(sp (StatementWithRest / StatementWithoutRest) !.)> */
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
		/* 1 StatementWithRest <- <(<(Statement sp ';' sp)> .* Action0)> */
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
					if !_rules[rulesp]() {
						goto l5
					}
					if buffer[position] != rune(';') {
						goto l5
					}
					position++
					if !_rules[rulesp]() {
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
		/* 2 StatementWithoutRest <- <(<(Statement sp)> Action1)> */
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
					if !_rules[rulesp]() {
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
		/* 8 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action2)> */
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
				if !_rules[rulesp]() {
					goto l49
				}
				if !_rules[ruleEmitter]() {
					goto l49
				}
				if !_rules[rulesp]() {
					goto l49
				}
				if !_rules[ruleProjections]() {
					goto l49
				}
				if !_rules[rulesp]() {
					goto l49
				}
				if !_rules[ruleWindowedFrom]() {
					goto l49
				}
				if !_rules[rulesp]() {
					goto l49
				}
				if !_rules[ruleFilter]() {
					goto l49
				}
				if !_rules[rulesp]() {
					goto l49
				}
				if !_rules[ruleGrouping]() {
					goto l49
				}
				if !_rules[rulesp]() {
					goto l49
				}
				if !_rules[ruleHaving]() {
					goto l49
				}
				if !_rules[rulesp]() {
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
		/* 9 SelectUnionStmt <- <(<(SelectStmt (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action3)> */
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
		/* 12 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action6)> */
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
				if !_rules[rulesp]() {
					goto l160
				}
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
				if !_rules[rulesp]() {
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
		/* 13 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action7)> */
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
				if !_rules[rulesp]() {
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
		/* 14 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action8)> */
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
				if !_rules[rulesp]() {
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
		/* 15 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action9)> */
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
				if !_rules[rulesp]() {
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
		/* 16 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action10)> */
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
				if !_rules[rulesp]() {
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
		/* 17 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action11)> */
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
				if !_rules[rulesp]() {
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
		/* 27 LoadStateStmt <- <(('l' / 'L') ('o' / 'O') ('a' / 'A') ('d' / 'D') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SetOptSpecs Action21)> */
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
				if !_rules[rulesp]() {
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
		/* 28 LoadStateOrCreateStmt <- <(LoadStateStmt sp (('o' / 'O') ('r' / 'R')) sp (('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp (('i' / 'I') ('f' / 'F')) sp (('n' / 'N') ('o' / 'O') ('t' / 'T')) sp (('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S')) sp SourceSinkSpecs Action22)> */
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
				if !_rules[rulesp]() {
					goto l566
				}
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
		/* 30 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) sp EmitterOptions Action24)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
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
				if !_rules[rulesp]() {
					goto l626
				}
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
		/* 31 EmitterOptions <- <(<('[' sp EmitterLimit sp ']')?> Action25)> */
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
						if buffer[position] != rune('[') {
							goto l634
						}
						position++
						if !_rules[rulesp]() {
							goto l634
						}
						if !_rules[ruleEmitterLimit]() {
							goto l634
						}
						if !_rules[rulesp]() {
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
		/* 32 EmitterLimit <- <('L' 'I' 'M' 'I' 'T' sp NumericLiteral Action26)> */
		func() bool {
			position636, tokenIndex636, depth636 := position, tokenIndex, depth
			{
				position637 := position
				depth++
				if buffer[position] != rune('L') {
					goto l636
				}
				position++
				if buffer[position] != rune('I') {
					goto l636
				}
				position++
				if buffer[position] != rune('M') {
					goto l636
				}
				position++
				if buffer[position] != rune('I') {
					goto l636
				}
				position++
				if buffer[position] != rune('T') {
					goto l636
				}
				position++
				if !_rules[rulesp]() {
					goto l636
				}
				if !_rules[ruleNumericLiteral]() {
					goto l636
				}
				if !_rules[ruleAction26]() {
					goto l636
				}
				depth--
				add(ruleEmitterLimit, position637)
			}
			return true
		l636:
			position, tokenIndex, depth = position636, tokenIndex636, depth636
			return false
		},
		/* 33 Projections <- <(<(Projection sp (',' sp Projection)*)> Action27)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l638
					}
					if !_rules[rulesp]() {
						goto l638
					}
				l641:
					{
						position642, tokenIndex642, depth642 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l642
						}
						position++
						if !_rules[rulesp]() {
							goto l642
						}
						if !_rules[ruleProjection]() {
							goto l642
						}
						goto l641
					l642:
						position, tokenIndex, depth = position642, tokenIndex642, depth642
					}
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction27]() {
					goto l638
				}
				depth--
				add(ruleProjections, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 34 Projection <- <(AliasExpression / ExpressionOrWildcard)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645, tokenIndex645, depth645 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l646
					}
					goto l645
				l646:
					position, tokenIndex, depth = position645, tokenIndex645, depth645
					if !_rules[ruleExpressionOrWildcard]() {
						goto l643
					}
				}
			l645:
				depth--
				add(ruleProjection, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 35 AliasExpression <- <(ExpressionOrWildcard sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action28)> */
		func() bool {
			position647, tokenIndex647, depth647 := position, tokenIndex, depth
			{
				position648 := position
				depth++
				if !_rules[ruleExpressionOrWildcard]() {
					goto l647
				}
				if !_rules[rulesp]() {
					goto l647
				}
				{
					position649, tokenIndex649, depth649 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l650
					}
					position++
					goto l649
				l650:
					position, tokenIndex, depth = position649, tokenIndex649, depth649
					if buffer[position] != rune('A') {
						goto l647
					}
					position++
				}
			l649:
				{
					position651, tokenIndex651, depth651 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l652
					}
					position++
					goto l651
				l652:
					position, tokenIndex, depth = position651, tokenIndex651, depth651
					if buffer[position] != rune('S') {
						goto l647
					}
					position++
				}
			l651:
				if !_rules[rulesp]() {
					goto l647
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l647
				}
				if !_rules[ruleAction28]() {
					goto l647
				}
				depth--
				add(ruleAliasExpression, position648)
			}
			return true
		l647:
			position, tokenIndex, depth = position647, tokenIndex647, depth647
			return false
		},
		/* 36 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action29)> */
		func() bool {
			position653, tokenIndex653, depth653 := position, tokenIndex, depth
			{
				position654 := position
				depth++
				{
					position655 := position
					depth++
					{
						position656, tokenIndex656, depth656 := position, tokenIndex, depth
						{
							position658, tokenIndex658, depth658 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l659
							}
							position++
							goto l658
						l659:
							position, tokenIndex, depth = position658, tokenIndex658, depth658
							if buffer[position] != rune('F') {
								goto l656
							}
							position++
						}
					l658:
						{
							position660, tokenIndex660, depth660 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l661
							}
							position++
							goto l660
						l661:
							position, tokenIndex, depth = position660, tokenIndex660, depth660
							if buffer[position] != rune('R') {
								goto l656
							}
							position++
						}
					l660:
						{
							position662, tokenIndex662, depth662 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l663
							}
							position++
							goto l662
						l663:
							position, tokenIndex, depth = position662, tokenIndex662, depth662
							if buffer[position] != rune('O') {
								goto l656
							}
							position++
						}
					l662:
						{
							position664, tokenIndex664, depth664 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l665
							}
							position++
							goto l664
						l665:
							position, tokenIndex, depth = position664, tokenIndex664, depth664
							if buffer[position] != rune('M') {
								goto l656
							}
							position++
						}
					l664:
						if !_rules[rulesp]() {
							goto l656
						}
						if !_rules[ruleRelations]() {
							goto l656
						}
						if !_rules[rulesp]() {
							goto l656
						}
						goto l657
					l656:
						position, tokenIndex, depth = position656, tokenIndex656, depth656
					}
				l657:
					depth--
					add(rulePegText, position655)
				}
				if !_rules[ruleAction29]() {
					goto l653
				}
				depth--
				add(ruleWindowedFrom, position654)
			}
			return true
		l653:
			position, tokenIndex, depth = position653, tokenIndex653, depth653
			return false
		},
		/* 37 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				{
					position668, tokenIndex668, depth668 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l669
					}
					goto l668
				l669:
					position, tokenIndex, depth = position668, tokenIndex668, depth668
					if !_rules[ruleTuplesInterval]() {
						goto l666
					}
				}
			l668:
				depth--
				add(ruleInterval, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 38 TimeInterval <- <(NumericLiteral sp (SECONDS / MILLISECONDS) Action30)> */
		func() bool {
			position670, tokenIndex670, depth670 := position, tokenIndex, depth
			{
				position671 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l670
				}
				if !_rules[rulesp]() {
					goto l670
				}
				{
					position672, tokenIndex672, depth672 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l673
					}
					goto l672
				l673:
					position, tokenIndex, depth = position672, tokenIndex672, depth672
					if !_rules[ruleMILLISECONDS]() {
						goto l670
					}
				}
			l672:
				if !_rules[ruleAction30]() {
					goto l670
				}
				depth--
				add(ruleTimeInterval, position671)
			}
			return true
		l670:
			position, tokenIndex, depth = position670, tokenIndex670, depth670
			return false
		},
		/* 39 TuplesInterval <- <(NumericLiteral sp TUPLES Action31)> */
		func() bool {
			position674, tokenIndex674, depth674 := position, tokenIndex, depth
			{
				position675 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l674
				}
				if !_rules[rulesp]() {
					goto l674
				}
				if !_rules[ruleTUPLES]() {
					goto l674
				}
				if !_rules[ruleAction31]() {
					goto l674
				}
				depth--
				add(ruleTuplesInterval, position675)
			}
			return true
		l674:
			position, tokenIndex, depth = position674, tokenIndex674, depth674
			return false
		},
		/* 40 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position676, tokenIndex676, depth676 := position, tokenIndex, depth
			{
				position677 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l676
				}
				if !_rules[rulesp]() {
					goto l676
				}
			l678:
				{
					position679, tokenIndex679, depth679 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l679
					}
					position++
					if !_rules[rulesp]() {
						goto l679
					}
					if !_rules[ruleRelationLike]() {
						goto l679
					}
					goto l678
				l679:
					position, tokenIndex, depth = position679, tokenIndex679, depth679
				}
				depth--
				add(ruleRelations, position677)
			}
			return true
		l676:
			position, tokenIndex, depth = position676, tokenIndex676, depth676
			return false
		},
		/* 41 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action32)> */
		func() bool {
			position680, tokenIndex680, depth680 := position, tokenIndex, depth
			{
				position681 := position
				depth++
				{
					position682 := position
					depth++
					{
						position683, tokenIndex683, depth683 := position, tokenIndex, depth
						{
							position685, tokenIndex685, depth685 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l686
							}
							position++
							goto l685
						l686:
							position, tokenIndex, depth = position685, tokenIndex685, depth685
							if buffer[position] != rune('W') {
								goto l683
							}
							position++
						}
					l685:
						{
							position687, tokenIndex687, depth687 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l688
							}
							position++
							goto l687
						l688:
							position, tokenIndex, depth = position687, tokenIndex687, depth687
							if buffer[position] != rune('H') {
								goto l683
							}
							position++
						}
					l687:
						{
							position689, tokenIndex689, depth689 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l690
							}
							position++
							goto l689
						l690:
							position, tokenIndex, depth = position689, tokenIndex689, depth689
							if buffer[position] != rune('E') {
								goto l683
							}
							position++
						}
					l689:
						{
							position691, tokenIndex691, depth691 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l692
							}
							position++
							goto l691
						l692:
							position, tokenIndex, depth = position691, tokenIndex691, depth691
							if buffer[position] != rune('R') {
								goto l683
							}
							position++
						}
					l691:
						{
							position693, tokenIndex693, depth693 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l694
							}
							position++
							goto l693
						l694:
							position, tokenIndex, depth = position693, tokenIndex693, depth693
							if buffer[position] != rune('E') {
								goto l683
							}
							position++
						}
					l693:
						if !_rules[rulesp]() {
							goto l683
						}
						if !_rules[ruleExpression]() {
							goto l683
						}
						goto l684
					l683:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
					}
				l684:
					depth--
					add(rulePegText, position682)
				}
				if !_rules[ruleAction32]() {
					goto l680
				}
				depth--
				add(ruleFilter, position681)
			}
			return true
		l680:
			position, tokenIndex, depth = position680, tokenIndex680, depth680
			return false
		},
		/* 42 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action33)> */
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
						{
							position700, tokenIndex700, depth700 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l701
							}
							position++
							goto l700
						l701:
							position, tokenIndex, depth = position700, tokenIndex700, depth700
							if buffer[position] != rune('G') {
								goto l698
							}
							position++
						}
					l700:
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
								goto l698
							}
							position++
						}
					l702:
						{
							position704, tokenIndex704, depth704 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l705
							}
							position++
							goto l704
						l705:
							position, tokenIndex, depth = position704, tokenIndex704, depth704
							if buffer[position] != rune('O') {
								goto l698
							}
							position++
						}
					l704:
						{
							position706, tokenIndex706, depth706 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l707
							}
							position++
							goto l706
						l707:
							position, tokenIndex, depth = position706, tokenIndex706, depth706
							if buffer[position] != rune('U') {
								goto l698
							}
							position++
						}
					l706:
						{
							position708, tokenIndex708, depth708 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l709
							}
							position++
							goto l708
						l709:
							position, tokenIndex, depth = position708, tokenIndex708, depth708
							if buffer[position] != rune('P') {
								goto l698
							}
							position++
						}
					l708:
						if !_rules[rulesp]() {
							goto l698
						}
						{
							position710, tokenIndex710, depth710 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l711
							}
							position++
							goto l710
						l711:
							position, tokenIndex, depth = position710, tokenIndex710, depth710
							if buffer[position] != rune('B') {
								goto l698
							}
							position++
						}
					l710:
						{
							position712, tokenIndex712, depth712 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l713
							}
							position++
							goto l712
						l713:
							position, tokenIndex, depth = position712, tokenIndex712, depth712
							if buffer[position] != rune('Y') {
								goto l698
							}
							position++
						}
					l712:
						if !_rules[rulesp]() {
							goto l698
						}
						if !_rules[ruleGroupList]() {
							goto l698
						}
						goto l699
					l698:
						position, tokenIndex, depth = position698, tokenIndex698, depth698
					}
				l699:
					depth--
					add(rulePegText, position697)
				}
				if !_rules[ruleAction33]() {
					goto l695
				}
				depth--
				add(ruleGrouping, position696)
			}
			return true
		l695:
			position, tokenIndex, depth = position695, tokenIndex695, depth695
			return false
		},
		/* 43 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position714, tokenIndex714, depth714 := position, tokenIndex, depth
			{
				position715 := position
				depth++
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
				depth--
				add(ruleGroupList, position715)
			}
			return true
		l714:
			position, tokenIndex, depth = position714, tokenIndex714, depth714
			return false
		},
		/* 44 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action34)> */
		func() bool {
			position718, tokenIndex718, depth718 := position, tokenIndex, depth
			{
				position719 := position
				depth++
				{
					position720 := position
					depth++
					{
						position721, tokenIndex721, depth721 := position, tokenIndex, depth
						{
							position723, tokenIndex723, depth723 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l724
							}
							position++
							goto l723
						l724:
							position, tokenIndex, depth = position723, tokenIndex723, depth723
							if buffer[position] != rune('H') {
								goto l721
							}
							position++
						}
					l723:
						{
							position725, tokenIndex725, depth725 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l726
							}
							position++
							goto l725
						l726:
							position, tokenIndex, depth = position725, tokenIndex725, depth725
							if buffer[position] != rune('A') {
								goto l721
							}
							position++
						}
					l725:
						{
							position727, tokenIndex727, depth727 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l728
							}
							position++
							goto l727
						l728:
							position, tokenIndex, depth = position727, tokenIndex727, depth727
							if buffer[position] != rune('V') {
								goto l721
							}
							position++
						}
					l727:
						{
							position729, tokenIndex729, depth729 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l730
							}
							position++
							goto l729
						l730:
							position, tokenIndex, depth = position729, tokenIndex729, depth729
							if buffer[position] != rune('I') {
								goto l721
							}
							position++
						}
					l729:
						{
							position731, tokenIndex731, depth731 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l732
							}
							position++
							goto l731
						l732:
							position, tokenIndex, depth = position731, tokenIndex731, depth731
							if buffer[position] != rune('N') {
								goto l721
							}
							position++
						}
					l731:
						{
							position733, tokenIndex733, depth733 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l734
							}
							position++
							goto l733
						l734:
							position, tokenIndex, depth = position733, tokenIndex733, depth733
							if buffer[position] != rune('G') {
								goto l721
							}
							position++
						}
					l733:
						if !_rules[rulesp]() {
							goto l721
						}
						if !_rules[ruleExpression]() {
							goto l721
						}
						goto l722
					l721:
						position, tokenIndex, depth = position721, tokenIndex721, depth721
					}
				l722:
					depth--
					add(rulePegText, position720)
				}
				if !_rules[ruleAction34]() {
					goto l718
				}
				depth--
				add(ruleHaving, position719)
			}
			return true
		l718:
			position, tokenIndex, depth = position718, tokenIndex718, depth718
			return false
		},
		/* 45 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action35))> */
		func() bool {
			position735, tokenIndex735, depth735 := position, tokenIndex, depth
			{
				position736 := position
				depth++
				{
					position737, tokenIndex737, depth737 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l738
					}
					goto l737
				l738:
					position, tokenIndex, depth = position737, tokenIndex737, depth737
					if !_rules[ruleStreamWindow]() {
						goto l735
					}
					if !_rules[ruleAction35]() {
						goto l735
					}
				}
			l737:
				depth--
				add(ruleRelationLike, position736)
			}
			return true
		l735:
			position, tokenIndex, depth = position735, tokenIndex735, depth735
			return false
		},
		/* 46 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action36)> */
		func() bool {
			position739, tokenIndex739, depth739 := position, tokenIndex, depth
			{
				position740 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l739
				}
				if !_rules[rulesp]() {
					goto l739
				}
				{
					position741, tokenIndex741, depth741 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l742
					}
					position++
					goto l741
				l742:
					position, tokenIndex, depth = position741, tokenIndex741, depth741
					if buffer[position] != rune('A') {
						goto l739
					}
					position++
				}
			l741:
				{
					position743, tokenIndex743, depth743 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l744
					}
					position++
					goto l743
				l744:
					position, tokenIndex, depth = position743, tokenIndex743, depth743
					if buffer[position] != rune('S') {
						goto l739
					}
					position++
				}
			l743:
				if !_rules[rulesp]() {
					goto l739
				}
				if !_rules[ruleIdentifier]() {
					goto l739
				}
				if !_rules[ruleAction36]() {
					goto l739
				}
				depth--
				add(ruleAliasedStreamWindow, position740)
			}
			return true
		l739:
			position, tokenIndex, depth = position739, tokenIndex739, depth739
			return false
		},
		/* 47 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action37)> */
		func() bool {
			position745, tokenIndex745, depth745 := position, tokenIndex, depth
			{
				position746 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l745
				}
				if !_rules[rulesp]() {
					goto l745
				}
				if buffer[position] != rune('[') {
					goto l745
				}
				position++
				if !_rules[rulesp]() {
					goto l745
				}
				{
					position747, tokenIndex747, depth747 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l748
					}
					position++
					goto l747
				l748:
					position, tokenIndex, depth = position747, tokenIndex747, depth747
					if buffer[position] != rune('R') {
						goto l745
					}
					position++
				}
			l747:
				{
					position749, tokenIndex749, depth749 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l750
					}
					position++
					goto l749
				l750:
					position, tokenIndex, depth = position749, tokenIndex749, depth749
					if buffer[position] != rune('A') {
						goto l745
					}
					position++
				}
			l749:
				{
					position751, tokenIndex751, depth751 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l752
					}
					position++
					goto l751
				l752:
					position, tokenIndex, depth = position751, tokenIndex751, depth751
					if buffer[position] != rune('N') {
						goto l745
					}
					position++
				}
			l751:
				{
					position753, tokenIndex753, depth753 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l754
					}
					position++
					goto l753
				l754:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if buffer[position] != rune('G') {
						goto l745
					}
					position++
				}
			l753:
				{
					position755, tokenIndex755, depth755 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l756
					}
					position++
					goto l755
				l756:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if buffer[position] != rune('E') {
						goto l745
					}
					position++
				}
			l755:
				if !_rules[rulesp]() {
					goto l745
				}
				if !_rules[ruleInterval]() {
					goto l745
				}
				if !_rules[rulesp]() {
					goto l745
				}
				if buffer[position] != rune(']') {
					goto l745
				}
				position++
				if !_rules[ruleAction37]() {
					goto l745
				}
				depth--
				add(ruleStreamWindow, position746)
			}
			return true
		l745:
			position, tokenIndex, depth = position745, tokenIndex745, depth745
			return false
		},
		/* 48 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position757, tokenIndex757, depth757 := position, tokenIndex, depth
			{
				position758 := position
				depth++
				{
					position759, tokenIndex759, depth759 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l760
					}
					goto l759
				l760:
					position, tokenIndex, depth = position759, tokenIndex759, depth759
					if !_rules[ruleStream]() {
						goto l757
					}
				}
			l759:
				depth--
				add(ruleStreamLike, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 49 UDSFFuncApp <- <(FuncApp Action38)> */
		func() bool {
			position761, tokenIndex761, depth761 := position, tokenIndex, depth
			{
				position762 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l761
				}
				if !_rules[ruleAction38]() {
					goto l761
				}
				depth--
				add(ruleUDSFFuncApp, position762)
			}
			return true
		l761:
			position, tokenIndex, depth = position761, tokenIndex761, depth761
			return false
		},
		/* 50 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action39)> */
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
						{
							position768, tokenIndex768, depth768 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l769
							}
							position++
							goto l768
						l769:
							position, tokenIndex, depth = position768, tokenIndex768, depth768
							if buffer[position] != rune('W') {
								goto l766
							}
							position++
						}
					l768:
						{
							position770, tokenIndex770, depth770 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l771
							}
							position++
							goto l770
						l771:
							position, tokenIndex, depth = position770, tokenIndex770, depth770
							if buffer[position] != rune('I') {
								goto l766
							}
							position++
						}
					l770:
						{
							position772, tokenIndex772, depth772 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l773
							}
							position++
							goto l772
						l773:
							position, tokenIndex, depth = position772, tokenIndex772, depth772
							if buffer[position] != rune('T') {
								goto l766
							}
							position++
						}
					l772:
						{
							position774, tokenIndex774, depth774 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l775
							}
							position++
							goto l774
						l775:
							position, tokenIndex, depth = position774, tokenIndex774, depth774
							if buffer[position] != rune('H') {
								goto l766
							}
							position++
						}
					l774:
						if !_rules[rulesp]() {
							goto l766
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l766
						}
						if !_rules[rulesp]() {
							goto l766
						}
					l776:
						{
							position777, tokenIndex777, depth777 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l777
							}
							position++
							if !_rules[rulesp]() {
								goto l777
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l777
							}
							goto l776
						l777:
							position, tokenIndex, depth = position777, tokenIndex777, depth777
						}
						goto l767
					l766:
						position, tokenIndex, depth = position766, tokenIndex766, depth766
					}
				l767:
					depth--
					add(rulePegText, position765)
				}
				if !_rules[ruleAction39]() {
					goto l763
				}
				depth--
				add(ruleSourceSinkSpecs, position764)
			}
			return true
		l763:
			position, tokenIndex, depth = position763, tokenIndex763, depth763
			return false
		},
		/* 51 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action40)> */
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
						if buffer[position] != rune('s') {
							goto l782
						}
						position++
						goto l781
					l782:
						position, tokenIndex, depth = position781, tokenIndex781, depth781
						if buffer[position] != rune('S') {
							goto l778
						}
						position++
					}
				l781:
					{
						position783, tokenIndex783, depth783 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l784
						}
						position++
						goto l783
					l784:
						position, tokenIndex, depth = position783, tokenIndex783, depth783
						if buffer[position] != rune('E') {
							goto l778
						}
						position++
					}
				l783:
					{
						position785, tokenIndex785, depth785 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l786
						}
						position++
						goto l785
					l786:
						position, tokenIndex, depth = position785, tokenIndex785, depth785
						if buffer[position] != rune('T') {
							goto l778
						}
						position++
					}
				l785:
					if !_rules[rulesp]() {
						goto l778
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l778
					}
					if !_rules[rulesp]() {
						goto l778
					}
				l787:
					{
						position788, tokenIndex788, depth788 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l788
						}
						position++
						if !_rules[rulesp]() {
							goto l788
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l788
						}
						goto l787
					l788:
						position, tokenIndex, depth = position788, tokenIndex788, depth788
					}
					depth--
					add(rulePegText, position780)
				}
				if !_rules[ruleAction40]() {
					goto l778
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position779)
			}
			return true
		l778:
			position, tokenIndex, depth = position778, tokenIndex778, depth778
			return false
		},
		/* 52 SetOptSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action41)> */
		func() bool {
			position789, tokenIndex789, depth789 := position, tokenIndex, depth
			{
				position790 := position
				depth++
				{
					position791 := position
					depth++
					{
						position792, tokenIndex792, depth792 := position, tokenIndex, depth
						{
							position794, tokenIndex794, depth794 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l795
							}
							position++
							goto l794
						l795:
							position, tokenIndex, depth = position794, tokenIndex794, depth794
							if buffer[position] != rune('S') {
								goto l792
							}
							position++
						}
					l794:
						{
							position796, tokenIndex796, depth796 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l797
							}
							position++
							goto l796
						l797:
							position, tokenIndex, depth = position796, tokenIndex796, depth796
							if buffer[position] != rune('E') {
								goto l792
							}
							position++
						}
					l796:
						{
							position798, tokenIndex798, depth798 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l799
							}
							position++
							goto l798
						l799:
							position, tokenIndex, depth = position798, tokenIndex798, depth798
							if buffer[position] != rune('T') {
								goto l792
							}
							position++
						}
					l798:
						if !_rules[rulesp]() {
							goto l792
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l792
						}
						if !_rules[rulesp]() {
							goto l792
						}
					l800:
						{
							position801, tokenIndex801, depth801 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l801
							}
							position++
							if !_rules[rulesp]() {
								goto l801
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l801
							}
							goto l800
						l801:
							position, tokenIndex, depth = position801, tokenIndex801, depth801
						}
						goto l793
					l792:
						position, tokenIndex, depth = position792, tokenIndex792, depth792
					}
				l793:
					depth--
					add(rulePegText, position791)
				}
				if !_rules[ruleAction41]() {
					goto l789
				}
				depth--
				add(ruleSetOptSpecs, position790)
			}
			return true
		l789:
			position, tokenIndex, depth = position789, tokenIndex789, depth789
			return false
		},
		/* 53 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action42)> */
		func() bool {
			position802, tokenIndex802, depth802 := position, tokenIndex, depth
			{
				position803 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l802
				}
				if buffer[position] != rune('=') {
					goto l802
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l802
				}
				if !_rules[ruleAction42]() {
					goto l802
				}
				depth--
				add(ruleSourceSinkParam, position803)
			}
			return true
		l802:
			position, tokenIndex, depth = position802, tokenIndex802, depth802
			return false
		},
		/* 54 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr)> */
		func() bool {
			position804, tokenIndex804, depth804 := position, tokenIndex, depth
			{
				position805 := position
				depth++
				{
					position806, tokenIndex806, depth806 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l807
					}
					goto l806
				l807:
					position, tokenIndex, depth = position806, tokenIndex806, depth806
					if !_rules[ruleParamArrayExpr]() {
						goto l804
					}
				}
			l806:
				depth--
				add(ruleSourceSinkParamVal, position805)
			}
			return true
		l804:
			position, tokenIndex, depth = position804, tokenIndex804, depth804
			return false
		},
		/* 55 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position808, tokenIndex808, depth808 := position, tokenIndex, depth
			{
				position809 := position
				depth++
				{
					position810, tokenIndex810, depth810 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l811
					}
					goto l810
				l811:
					position, tokenIndex, depth = position810, tokenIndex810, depth810
					if !_rules[ruleLiteral]() {
						goto l808
					}
				}
			l810:
				depth--
				add(ruleParamLiteral, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 56 ParamArrayExpr <- <(<('[' sp (ParamLiteral (',' sp ParamLiteral)*)? sp ','? sp ']')> Action43)> */
		func() bool {
			position812, tokenIndex812, depth812 := position, tokenIndex, depth
			{
				position813 := position
				depth++
				{
					position814 := position
					depth++
					if buffer[position] != rune('[') {
						goto l812
					}
					position++
					if !_rules[rulesp]() {
						goto l812
					}
					{
						position815, tokenIndex815, depth815 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l815
						}
					l817:
						{
							position818, tokenIndex818, depth818 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l818
							}
							position++
							if !_rules[rulesp]() {
								goto l818
							}
							if !_rules[ruleParamLiteral]() {
								goto l818
							}
							goto l817
						l818:
							position, tokenIndex, depth = position818, tokenIndex818, depth818
						}
						goto l816
					l815:
						position, tokenIndex, depth = position815, tokenIndex815, depth815
					}
				l816:
					if !_rules[rulesp]() {
						goto l812
					}
					{
						position819, tokenIndex819, depth819 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l819
						}
						position++
						goto l820
					l819:
						position, tokenIndex, depth = position819, tokenIndex819, depth819
					}
				l820:
					if !_rules[rulesp]() {
						goto l812
					}
					if buffer[position] != rune(']') {
						goto l812
					}
					position++
					depth--
					add(rulePegText, position814)
				}
				if !_rules[ruleAction43]() {
					goto l812
				}
				depth--
				add(ruleParamArrayExpr, position813)
			}
			return true
		l812:
			position, tokenIndex, depth = position812, tokenIndex812, depth812
			return false
		},
		/* 57 PausedOpt <- <(<(Paused / Unpaused)?> Action44)> */
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
						{
							position826, tokenIndex826, depth826 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l827
							}
							goto l826
						l827:
							position, tokenIndex, depth = position826, tokenIndex826, depth826
							if !_rules[ruleUnpaused]() {
								goto l824
							}
						}
					l826:
						goto l825
					l824:
						position, tokenIndex, depth = position824, tokenIndex824, depth824
					}
				l825:
					depth--
					add(rulePegText, position823)
				}
				if !_rules[ruleAction44]() {
					goto l821
				}
				depth--
				add(rulePausedOpt, position822)
			}
			return true
		l821:
			position, tokenIndex, depth = position821, tokenIndex821, depth821
			return false
		},
		/* 58 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position828, tokenIndex828, depth828 := position, tokenIndex, depth
			{
				position829 := position
				depth++
				{
					position830, tokenIndex830, depth830 := position, tokenIndex, depth
					if !_rules[ruleWildcard]() {
						goto l831
					}
					goto l830
				l831:
					position, tokenIndex, depth = position830, tokenIndex830, depth830
					if !_rules[ruleExpression]() {
						goto l828
					}
				}
			l830:
				depth--
				add(ruleExpressionOrWildcard, position829)
			}
			return true
		l828:
			position, tokenIndex, depth = position828, tokenIndex828, depth828
			return false
		},
		/* 59 Expression <- <orExpr> */
		func() bool {
			position832, tokenIndex832, depth832 := position, tokenIndex, depth
			{
				position833 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l832
				}
				depth--
				add(ruleExpression, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 60 orExpr <- <(<(andExpr sp (Or sp andExpr)*)> Action45)> */
		func() bool {
			position834, tokenIndex834, depth834 := position, tokenIndex, depth
			{
				position835 := position
				depth++
				{
					position836 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l834
					}
					if !_rules[rulesp]() {
						goto l834
					}
				l837:
					{
						position838, tokenIndex838, depth838 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l838
						}
						if !_rules[rulesp]() {
							goto l838
						}
						if !_rules[ruleandExpr]() {
							goto l838
						}
						goto l837
					l838:
						position, tokenIndex, depth = position838, tokenIndex838, depth838
					}
					depth--
					add(rulePegText, position836)
				}
				if !_rules[ruleAction45]() {
					goto l834
				}
				depth--
				add(ruleorExpr, position835)
			}
			return true
		l834:
			position, tokenIndex, depth = position834, tokenIndex834, depth834
			return false
		},
		/* 61 andExpr <- <(<(notExpr sp (And sp notExpr)*)> Action46)> */
		func() bool {
			position839, tokenIndex839, depth839 := position, tokenIndex, depth
			{
				position840 := position
				depth++
				{
					position841 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l839
					}
					if !_rules[rulesp]() {
						goto l839
					}
				l842:
					{
						position843, tokenIndex843, depth843 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l843
						}
						if !_rules[rulesp]() {
							goto l843
						}
						if !_rules[rulenotExpr]() {
							goto l843
						}
						goto l842
					l843:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
					}
					depth--
					add(rulePegText, position841)
				}
				if !_rules[ruleAction46]() {
					goto l839
				}
				depth--
				add(ruleandExpr, position840)
			}
			return true
		l839:
			position, tokenIndex, depth = position839, tokenIndex839, depth839
			return false
		},
		/* 62 notExpr <- <(<((Not sp)? comparisonExpr)> Action47)> */
		func() bool {
			position844, tokenIndex844, depth844 := position, tokenIndex, depth
			{
				position845 := position
				depth++
				{
					position846 := position
					depth++
					{
						position847, tokenIndex847, depth847 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l847
						}
						if !_rules[rulesp]() {
							goto l847
						}
						goto l848
					l847:
						position, tokenIndex, depth = position847, tokenIndex847, depth847
					}
				l848:
					if !_rules[rulecomparisonExpr]() {
						goto l844
					}
					depth--
					add(rulePegText, position846)
				}
				if !_rules[ruleAction47]() {
					goto l844
				}
				depth--
				add(rulenotExpr, position845)
			}
			return true
		l844:
			position, tokenIndex, depth = position844, tokenIndex844, depth844
			return false
		},
		/* 63 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action48)> */
		func() bool {
			position849, tokenIndex849, depth849 := position, tokenIndex, depth
			{
				position850 := position
				depth++
				{
					position851 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l849
					}
					if !_rules[rulesp]() {
						goto l849
					}
					{
						position852, tokenIndex852, depth852 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l852
						}
						if !_rules[rulesp]() {
							goto l852
						}
						if !_rules[ruleotherOpExpr]() {
							goto l852
						}
						goto l853
					l852:
						position, tokenIndex, depth = position852, tokenIndex852, depth852
					}
				l853:
					depth--
					add(rulePegText, position851)
				}
				if !_rules[ruleAction48]() {
					goto l849
				}
				depth--
				add(rulecomparisonExpr, position850)
			}
			return true
		l849:
			position, tokenIndex, depth = position849, tokenIndex849, depth849
			return false
		},
		/* 64 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action49)> */
		func() bool {
			position854, tokenIndex854, depth854 := position, tokenIndex, depth
			{
				position855 := position
				depth++
				{
					position856 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l854
					}
					if !_rules[rulesp]() {
						goto l854
					}
				l857:
					{
						position858, tokenIndex858, depth858 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l858
						}
						if !_rules[rulesp]() {
							goto l858
						}
						if !_rules[ruleisExpr]() {
							goto l858
						}
						if !_rules[rulesp]() {
							goto l858
						}
						goto l857
					l858:
						position, tokenIndex, depth = position858, tokenIndex858, depth858
					}
					depth--
					add(rulePegText, position856)
				}
				if !_rules[ruleAction49]() {
					goto l854
				}
				depth--
				add(ruleotherOpExpr, position855)
			}
			return true
		l854:
			position, tokenIndex, depth = position854, tokenIndex854, depth854
			return false
		},
		/* 65 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action50)> */
		func() bool {
			position859, tokenIndex859, depth859 := position, tokenIndex, depth
			{
				position860 := position
				depth++
				{
					position861 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l859
					}
					if !_rules[rulesp]() {
						goto l859
					}
					{
						position862, tokenIndex862, depth862 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l862
						}
						if !_rules[rulesp]() {
							goto l862
						}
						if !_rules[ruleNullLiteral]() {
							goto l862
						}
						goto l863
					l862:
						position, tokenIndex, depth = position862, tokenIndex862, depth862
					}
				l863:
					depth--
					add(rulePegText, position861)
				}
				if !_rules[ruleAction50]() {
					goto l859
				}
				depth--
				add(ruleisExpr, position860)
			}
			return true
		l859:
			position, tokenIndex, depth = position859, tokenIndex859, depth859
			return false
		},
		/* 66 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action51)> */
		func() bool {
			position864, tokenIndex864, depth864 := position, tokenIndex, depth
			{
				position865 := position
				depth++
				{
					position866 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l864
					}
					if !_rules[rulesp]() {
						goto l864
					}
				l867:
					{
						position868, tokenIndex868, depth868 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l868
						}
						if !_rules[rulesp]() {
							goto l868
						}
						if !_rules[ruleproductExpr]() {
							goto l868
						}
						if !_rules[rulesp]() {
							goto l868
						}
						goto l867
					l868:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
					}
					depth--
					add(rulePegText, position866)
				}
				if !_rules[ruleAction51]() {
					goto l864
				}
				depth--
				add(ruletermExpr, position865)
			}
			return true
		l864:
			position, tokenIndex, depth = position864, tokenIndex864, depth864
			return false
		},
		/* 67 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action52)> */
		func() bool {
			position869, tokenIndex869, depth869 := position, tokenIndex, depth
			{
				position870 := position
				depth++
				{
					position871 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l869
					}
					if !_rules[rulesp]() {
						goto l869
					}
				l872:
					{
						position873, tokenIndex873, depth873 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l873
						}
						if !_rules[rulesp]() {
							goto l873
						}
						if !_rules[ruleminusExpr]() {
							goto l873
						}
						if !_rules[rulesp]() {
							goto l873
						}
						goto l872
					l873:
						position, tokenIndex, depth = position873, tokenIndex873, depth873
					}
					depth--
					add(rulePegText, position871)
				}
				if !_rules[ruleAction52]() {
					goto l869
				}
				depth--
				add(ruleproductExpr, position870)
			}
			return true
		l869:
			position, tokenIndex, depth = position869, tokenIndex869, depth869
			return false
		},
		/* 68 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action53)> */
		func() bool {
			position874, tokenIndex874, depth874 := position, tokenIndex, depth
			{
				position875 := position
				depth++
				{
					position876 := position
					depth++
					{
						position877, tokenIndex877, depth877 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l877
						}
						if !_rules[rulesp]() {
							goto l877
						}
						goto l878
					l877:
						position, tokenIndex, depth = position877, tokenIndex877, depth877
					}
				l878:
					if !_rules[rulecastExpr]() {
						goto l874
					}
					depth--
					add(rulePegText, position876)
				}
				if !_rules[ruleAction53]() {
					goto l874
				}
				depth--
				add(ruleminusExpr, position875)
			}
			return true
		l874:
			position, tokenIndex, depth = position874, tokenIndex874, depth874
			return false
		},
		/* 69 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action54)> */
		func() bool {
			position879, tokenIndex879, depth879 := position, tokenIndex, depth
			{
				position880 := position
				depth++
				{
					position881 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l879
					}
					{
						position882, tokenIndex882, depth882 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l882
						}
						if buffer[position] != rune(':') {
							goto l882
						}
						position++
						if buffer[position] != rune(':') {
							goto l882
						}
						position++
						if !_rules[rulesp]() {
							goto l882
						}
						if !_rules[ruleType]() {
							goto l882
						}
						goto l883
					l882:
						position, tokenIndex, depth = position882, tokenIndex882, depth882
					}
				l883:
					depth--
					add(rulePegText, position881)
				}
				if !_rules[ruleAction54]() {
					goto l879
				}
				depth--
				add(rulecastExpr, position880)
			}
			return true
		l879:
			position, tokenIndex, depth = position879, tokenIndex879, depth879
			return false
		},
		/* 70 baseExpr <- <(('(' sp Expression sp ')') / MapExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position884, tokenIndex884, depth884 := position, tokenIndex, depth
			{
				position885 := position
				depth++
				{
					position886, tokenIndex886, depth886 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l887
					}
					position++
					if !_rules[rulesp]() {
						goto l887
					}
					if !_rules[ruleExpression]() {
						goto l887
					}
					if !_rules[rulesp]() {
						goto l887
					}
					if buffer[position] != rune(')') {
						goto l887
					}
					position++
					goto l886
				l887:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleMapExpr]() {
						goto l888
					}
					goto l886
				l888:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleBooleanLiteral]() {
						goto l889
					}
					goto l886
				l889:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleNullLiteral]() {
						goto l890
					}
					goto l886
				l890:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleRowMeta]() {
						goto l891
					}
					goto l886
				l891:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleFuncTypeCast]() {
						goto l892
					}
					goto l886
				l892:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleFuncApp]() {
						goto l893
					}
					goto l886
				l893:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleRowValue]() {
						goto l894
					}
					goto l886
				l894:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleArrayExpr]() {
						goto l895
					}
					goto l886
				l895:
					position, tokenIndex, depth = position886, tokenIndex886, depth886
					if !_rules[ruleLiteral]() {
						goto l884
					}
				}
			l886:
				depth--
				add(rulebaseExpr, position885)
			}
			return true
		l884:
			position, tokenIndex, depth = position884, tokenIndex884, depth884
			return false
		},
		/* 71 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action55)> */
		func() bool {
			position896, tokenIndex896, depth896 := position, tokenIndex, depth
			{
				position897 := position
				depth++
				{
					position898 := position
					depth++
					{
						position899, tokenIndex899, depth899 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l900
						}
						position++
						goto l899
					l900:
						position, tokenIndex, depth = position899, tokenIndex899, depth899
						if buffer[position] != rune('C') {
							goto l896
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
							goto l896
						}
						position++
					}
				l901:
					{
						position903, tokenIndex903, depth903 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l904
						}
						position++
						goto l903
					l904:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if buffer[position] != rune('S') {
							goto l896
						}
						position++
					}
				l903:
					{
						position905, tokenIndex905, depth905 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l906
						}
						position++
						goto l905
					l906:
						position, tokenIndex, depth = position905, tokenIndex905, depth905
						if buffer[position] != rune('T') {
							goto l896
						}
						position++
					}
				l905:
					if !_rules[rulesp]() {
						goto l896
					}
					if buffer[position] != rune('(') {
						goto l896
					}
					position++
					if !_rules[rulesp]() {
						goto l896
					}
					if !_rules[ruleExpression]() {
						goto l896
					}
					if !_rules[rulesp]() {
						goto l896
					}
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
							goto l896
						}
						position++
					}
				l907:
					{
						position909, tokenIndex909, depth909 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l910
						}
						position++
						goto l909
					l910:
						position, tokenIndex, depth = position909, tokenIndex909, depth909
						if buffer[position] != rune('S') {
							goto l896
						}
						position++
					}
				l909:
					if !_rules[rulesp]() {
						goto l896
					}
					if !_rules[ruleType]() {
						goto l896
					}
					if !_rules[rulesp]() {
						goto l896
					}
					if buffer[position] != rune(')') {
						goto l896
					}
					position++
					depth--
					add(rulePegText, position898)
				}
				if !_rules[ruleAction55]() {
					goto l896
				}
				depth--
				add(ruleFuncTypeCast, position897)
			}
			return true
		l896:
			position, tokenIndex, depth = position896, tokenIndex896, depth896
			return false
		},
		/* 72 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action56)> */
		func() bool {
			position911, tokenIndex911, depth911 := position, tokenIndex, depth
			{
				position912 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l911
				}
				if !_rules[rulesp]() {
					goto l911
				}
				if buffer[position] != rune('(') {
					goto l911
				}
				position++
				if !_rules[rulesp]() {
					goto l911
				}
				if !_rules[ruleFuncParams]() {
					goto l911
				}
				if !_rules[rulesp]() {
					goto l911
				}
				if buffer[position] != rune(')') {
					goto l911
				}
				position++
				if !_rules[ruleAction56]() {
					goto l911
				}
				depth--
				add(ruleFuncApp, position912)
			}
			return true
		l911:
			position, tokenIndex, depth = position911, tokenIndex911, depth911
			return false
		},
		/* 73 FuncParams <- <(<(ExpressionOrWildcard sp (',' sp ExpressionOrWildcard)*)?> Action57)> */
		func() bool {
			position913, tokenIndex913, depth913 := position, tokenIndex, depth
			{
				position914 := position
				depth++
				{
					position915 := position
					depth++
					{
						position916, tokenIndex916, depth916 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l916
						}
						if !_rules[rulesp]() {
							goto l916
						}
					l918:
						{
							position919, tokenIndex919, depth919 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l919
							}
							position++
							if !_rules[rulesp]() {
								goto l919
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l919
							}
							goto l918
						l919:
							position, tokenIndex, depth = position919, tokenIndex919, depth919
						}
						goto l917
					l916:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
					}
				l917:
					depth--
					add(rulePegText, position915)
				}
				if !_rules[ruleAction57]() {
					goto l913
				}
				depth--
				add(ruleFuncParams, position914)
			}
			return true
		l913:
			position, tokenIndex, depth = position913, tokenIndex913, depth913
			return false
		},
		/* 74 ArrayExpr <- <(<('[' sp (ExpressionOrWildcard (',' sp ExpressionOrWildcard)*)? sp ','? sp ']')> Action58)> */
		func() bool {
			position920, tokenIndex920, depth920 := position, tokenIndex, depth
			{
				position921 := position
				depth++
				{
					position922 := position
					depth++
					if buffer[position] != rune('[') {
						goto l920
					}
					position++
					if !_rules[rulesp]() {
						goto l920
					}
					{
						position923, tokenIndex923, depth923 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l923
						}
					l925:
						{
							position926, tokenIndex926, depth926 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l926
							}
							position++
							if !_rules[rulesp]() {
								goto l926
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l926
							}
							goto l925
						l926:
							position, tokenIndex, depth = position926, tokenIndex926, depth926
						}
						goto l924
					l923:
						position, tokenIndex, depth = position923, tokenIndex923, depth923
					}
				l924:
					if !_rules[rulesp]() {
						goto l920
					}
					{
						position927, tokenIndex927, depth927 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l927
						}
						position++
						goto l928
					l927:
						position, tokenIndex, depth = position927, tokenIndex927, depth927
					}
				l928:
					if !_rules[rulesp]() {
						goto l920
					}
					if buffer[position] != rune(']') {
						goto l920
					}
					position++
					depth--
					add(rulePegText, position922)
				}
				if !_rules[ruleAction58]() {
					goto l920
				}
				depth--
				add(ruleArrayExpr, position921)
			}
			return true
		l920:
			position, tokenIndex, depth = position920, tokenIndex920, depth920
			return false
		},
		/* 75 MapExpr <- <(<('{' sp (KeyValuePair (',' sp KeyValuePair)*)? sp '}')> Action59)> */
		func() bool {
			position929, tokenIndex929, depth929 := position, tokenIndex, depth
			{
				position930 := position
				depth++
				{
					position931 := position
					depth++
					if buffer[position] != rune('{') {
						goto l929
					}
					position++
					if !_rules[rulesp]() {
						goto l929
					}
					{
						position932, tokenIndex932, depth932 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l932
						}
					l934:
						{
							position935, tokenIndex935, depth935 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l935
							}
							position++
							if !_rules[rulesp]() {
								goto l935
							}
							if !_rules[ruleKeyValuePair]() {
								goto l935
							}
							goto l934
						l935:
							position, tokenIndex, depth = position935, tokenIndex935, depth935
						}
						goto l933
					l932:
						position, tokenIndex, depth = position932, tokenIndex932, depth932
					}
				l933:
					if !_rules[rulesp]() {
						goto l929
					}
					if buffer[position] != rune('}') {
						goto l929
					}
					position++
					depth--
					add(rulePegText, position931)
				}
				if !_rules[ruleAction59]() {
					goto l929
				}
				depth--
				add(ruleMapExpr, position930)
			}
			return true
		l929:
			position, tokenIndex, depth = position929, tokenIndex929, depth929
			return false
		},
		/* 76 KeyValuePair <- <(<(StringLiteral sp ':' sp ExpressionOrWildcard)> Action60)> */
		func() bool {
			position936, tokenIndex936, depth936 := position, tokenIndex, depth
			{
				position937 := position
				depth++
				{
					position938 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l936
					}
					if !_rules[rulesp]() {
						goto l936
					}
					if buffer[position] != rune(':') {
						goto l936
					}
					position++
					if !_rules[rulesp]() {
						goto l936
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l936
					}
					depth--
					add(rulePegText, position938)
				}
				if !_rules[ruleAction60]() {
					goto l936
				}
				depth--
				add(ruleKeyValuePair, position937)
			}
			return true
		l936:
			position, tokenIndex, depth = position936, tokenIndex936, depth936
			return false
		},
		/* 77 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position939, tokenIndex939, depth939 := position, tokenIndex, depth
			{
				position940 := position
				depth++
				{
					position941, tokenIndex941, depth941 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l942
					}
					goto l941
				l942:
					position, tokenIndex, depth = position941, tokenIndex941, depth941
					if !_rules[ruleNumericLiteral]() {
						goto l943
					}
					goto l941
				l943:
					position, tokenIndex, depth = position941, tokenIndex941, depth941
					if !_rules[ruleStringLiteral]() {
						goto l939
					}
				}
			l941:
				depth--
				add(ruleLiteral, position940)
			}
			return true
		l939:
			position, tokenIndex, depth = position939, tokenIndex939, depth939
			return false
		},
		/* 78 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position944, tokenIndex944, depth944 := position, tokenIndex, depth
			{
				position945 := position
				depth++
				{
					position946, tokenIndex946, depth946 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l947
					}
					goto l946
				l947:
					position, tokenIndex, depth = position946, tokenIndex946, depth946
					if !_rules[ruleNotEqual]() {
						goto l948
					}
					goto l946
				l948:
					position, tokenIndex, depth = position946, tokenIndex946, depth946
					if !_rules[ruleLessOrEqual]() {
						goto l949
					}
					goto l946
				l949:
					position, tokenIndex, depth = position946, tokenIndex946, depth946
					if !_rules[ruleLess]() {
						goto l950
					}
					goto l946
				l950:
					position, tokenIndex, depth = position946, tokenIndex946, depth946
					if !_rules[ruleGreaterOrEqual]() {
						goto l951
					}
					goto l946
				l951:
					position, tokenIndex, depth = position946, tokenIndex946, depth946
					if !_rules[ruleGreater]() {
						goto l952
					}
					goto l946
				l952:
					position, tokenIndex, depth = position946, tokenIndex946, depth946
					if !_rules[ruleNotEqual]() {
						goto l944
					}
				}
			l946:
				depth--
				add(ruleComparisonOp, position945)
			}
			return true
		l944:
			position, tokenIndex, depth = position944, tokenIndex944, depth944
			return false
		},
		/* 79 OtherOp <- <Concat> */
		func() bool {
			position953, tokenIndex953, depth953 := position, tokenIndex, depth
			{
				position954 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l953
				}
				depth--
				add(ruleOtherOp, position954)
			}
			return true
		l953:
			position, tokenIndex, depth = position953, tokenIndex953, depth953
			return false
		},
		/* 80 IsOp <- <(IsNot / Is)> */
		func() bool {
			position955, tokenIndex955, depth955 := position, tokenIndex, depth
			{
				position956 := position
				depth++
				{
					position957, tokenIndex957, depth957 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l958
					}
					goto l957
				l958:
					position, tokenIndex, depth = position957, tokenIndex957, depth957
					if !_rules[ruleIs]() {
						goto l955
					}
				}
			l957:
				depth--
				add(ruleIsOp, position956)
			}
			return true
		l955:
			position, tokenIndex, depth = position955, tokenIndex955, depth955
			return false
		},
		/* 81 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position959, tokenIndex959, depth959 := position, tokenIndex, depth
			{
				position960 := position
				depth++
				{
					position961, tokenIndex961, depth961 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l962
					}
					goto l961
				l962:
					position, tokenIndex, depth = position961, tokenIndex961, depth961
					if !_rules[ruleMinus]() {
						goto l959
					}
				}
			l961:
				depth--
				add(rulePlusMinusOp, position960)
			}
			return true
		l959:
			position, tokenIndex, depth = position959, tokenIndex959, depth959
			return false
		},
		/* 82 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position963, tokenIndex963, depth963 := position, tokenIndex, depth
			{
				position964 := position
				depth++
				{
					position965, tokenIndex965, depth965 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l966
					}
					goto l965
				l966:
					position, tokenIndex, depth = position965, tokenIndex965, depth965
					if !_rules[ruleDivide]() {
						goto l967
					}
					goto l965
				l967:
					position, tokenIndex, depth = position965, tokenIndex965, depth965
					if !_rules[ruleModulo]() {
						goto l963
					}
				}
			l965:
				depth--
				add(ruleMultDivOp, position964)
			}
			return true
		l963:
			position, tokenIndex, depth = position963, tokenIndex963, depth963
			return false
		},
		/* 83 Stream <- <(<ident> Action61)> */
		func() bool {
			position968, tokenIndex968, depth968 := position, tokenIndex, depth
			{
				position969 := position
				depth++
				{
					position970 := position
					depth++
					if !_rules[ruleident]() {
						goto l968
					}
					depth--
					add(rulePegText, position970)
				}
				if !_rules[ruleAction61]() {
					goto l968
				}
				depth--
				add(ruleStream, position969)
			}
			return true
		l968:
			position, tokenIndex, depth = position968, tokenIndex968, depth968
			return false
		},
		/* 84 RowMeta <- <RowTimestamp> */
		func() bool {
			position971, tokenIndex971, depth971 := position, tokenIndex, depth
			{
				position972 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l971
				}
				depth--
				add(ruleRowMeta, position972)
			}
			return true
		l971:
			position, tokenIndex, depth = position971, tokenIndex971, depth971
			return false
		},
		/* 85 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action62)> */
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
						if !_rules[ruleident]() {
							goto l976
						}
						if buffer[position] != rune(':') {
							goto l976
						}
						position++
						goto l977
					l976:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
					}
				l977:
					if buffer[position] != rune('t') {
						goto l973
					}
					position++
					if buffer[position] != rune('s') {
						goto l973
					}
					position++
					if buffer[position] != rune('(') {
						goto l973
					}
					position++
					if buffer[position] != rune(')') {
						goto l973
					}
					position++
					depth--
					add(rulePegText, position975)
				}
				if !_rules[ruleAction62]() {
					goto l973
				}
				depth--
				add(ruleRowTimestamp, position974)
			}
			return true
		l973:
			position, tokenIndex, depth = position973, tokenIndex973, depth973
			return false
		},
		/* 86 RowValue <- <(<((ident ':' !':')? jsonPath)> Action63)> */
		func() bool {
			position978, tokenIndex978, depth978 := position, tokenIndex, depth
			{
				position979 := position
				depth++
				{
					position980 := position
					depth++
					{
						position981, tokenIndex981, depth981 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l981
						}
						if buffer[position] != rune(':') {
							goto l981
						}
						position++
						{
							position983, tokenIndex983, depth983 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l983
							}
							position++
							goto l981
						l983:
							position, tokenIndex, depth = position983, tokenIndex983, depth983
						}
						goto l982
					l981:
						position, tokenIndex, depth = position981, tokenIndex981, depth981
					}
				l982:
					if !_rules[rulejsonPath]() {
						goto l978
					}
					depth--
					add(rulePegText, position980)
				}
				if !_rules[ruleAction63]() {
					goto l978
				}
				depth--
				add(ruleRowValue, position979)
			}
			return true
		l978:
			position, tokenIndex, depth = position978, tokenIndex978, depth978
			return false
		},
		/* 87 NumericLiteral <- <(<('-'? [0-9]+)> Action64)> */
		func() bool {
			position984, tokenIndex984, depth984 := position, tokenIndex, depth
			{
				position985 := position
				depth++
				{
					position986 := position
					depth++
					{
						position987, tokenIndex987, depth987 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l987
						}
						position++
						goto l988
					l987:
						position, tokenIndex, depth = position987, tokenIndex987, depth987
					}
				l988:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l984
					}
					position++
				l989:
					{
						position990, tokenIndex990, depth990 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l990
						}
						position++
						goto l989
					l990:
						position, tokenIndex, depth = position990, tokenIndex990, depth990
					}
					depth--
					add(rulePegText, position986)
				}
				if !_rules[ruleAction64]() {
					goto l984
				}
				depth--
				add(ruleNumericLiteral, position985)
			}
			return true
		l984:
			position, tokenIndex, depth = position984, tokenIndex984, depth984
			return false
		},
		/* 88 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action65)> */
		func() bool {
			position991, tokenIndex991, depth991 := position, tokenIndex, depth
			{
				position992 := position
				depth++
				{
					position993 := position
					depth++
					{
						position994, tokenIndex994, depth994 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l994
						}
						position++
						goto l995
					l994:
						position, tokenIndex, depth = position994, tokenIndex994, depth994
					}
				l995:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l991
					}
					position++
				l996:
					{
						position997, tokenIndex997, depth997 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l997
						}
						position++
						goto l996
					l997:
						position, tokenIndex, depth = position997, tokenIndex997, depth997
					}
					if buffer[position] != rune('.') {
						goto l991
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l991
					}
					position++
				l998:
					{
						position999, tokenIndex999, depth999 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l999
						}
						position++
						goto l998
					l999:
						position, tokenIndex, depth = position999, tokenIndex999, depth999
					}
					depth--
					add(rulePegText, position993)
				}
				if !_rules[ruleAction65]() {
					goto l991
				}
				depth--
				add(ruleFloatLiteral, position992)
			}
			return true
		l991:
			position, tokenIndex, depth = position991, tokenIndex991, depth991
			return false
		},
		/* 89 Function <- <(<ident> Action66)> */
		func() bool {
			position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
			{
				position1001 := position
				depth++
				{
					position1002 := position
					depth++
					if !_rules[ruleident]() {
						goto l1000
					}
					depth--
					add(rulePegText, position1002)
				}
				if !_rules[ruleAction66]() {
					goto l1000
				}
				depth--
				add(ruleFunction, position1001)
			}
			return true
		l1000:
			position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
			return false
		},
		/* 90 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action67)> */
		func() bool {
			position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
			{
				position1004 := position
				depth++
				{
					position1005 := position
					depth++
					{
						position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1007
						}
						position++
						goto l1006
					l1007:
						position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
						if buffer[position] != rune('N') {
							goto l1003
						}
						position++
					}
				l1006:
					{
						position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1009
						}
						position++
						goto l1008
					l1009:
						position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
						if buffer[position] != rune('U') {
							goto l1003
						}
						position++
					}
				l1008:
					{
						position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1011
						}
						position++
						goto l1010
					l1011:
						position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
						if buffer[position] != rune('L') {
							goto l1003
						}
						position++
					}
				l1010:
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('L') {
							goto l1003
						}
						position++
					}
				l1012:
					depth--
					add(rulePegText, position1005)
				}
				if !_rules[ruleAction67]() {
					goto l1003
				}
				depth--
				add(ruleNullLiteral, position1004)
			}
			return true
		l1003:
			position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
			return false
		},
		/* 91 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
			{
				position1015 := position
				depth++
				{
					position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l1017
					}
					goto l1016
				l1017:
					position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
					if !_rules[ruleFALSE]() {
						goto l1014
					}
				}
			l1016:
				depth--
				add(ruleBooleanLiteral, position1015)
			}
			return true
		l1014:
			position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
			return false
		},
		/* 92 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action68)> */
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
						if buffer[position] != rune('t') {
							goto l1022
						}
						position++
						goto l1021
					l1022:
						position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
						if buffer[position] != rune('T') {
							goto l1018
						}
						position++
					}
				l1021:
					{
						position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1024
						}
						position++
						goto l1023
					l1024:
						position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
						if buffer[position] != rune('R') {
							goto l1018
						}
						position++
					}
				l1023:
					{
						position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1026
						}
						position++
						goto l1025
					l1026:
						position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
						if buffer[position] != rune('U') {
							goto l1018
						}
						position++
					}
				l1025:
					{
						position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1028
						}
						position++
						goto l1027
					l1028:
						position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
						if buffer[position] != rune('E') {
							goto l1018
						}
						position++
					}
				l1027:
					depth--
					add(rulePegText, position1020)
				}
				if !_rules[ruleAction68]() {
					goto l1018
				}
				depth--
				add(ruleTRUE, position1019)
			}
			return true
		l1018:
			position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
			return false
		},
		/* 93 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action69)> */
		func() bool {
			position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
			{
				position1030 := position
				depth++
				{
					position1031 := position
					depth++
					{
						position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1033
						}
						position++
						goto l1032
					l1033:
						position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
						if buffer[position] != rune('F') {
							goto l1029
						}
						position++
					}
				l1032:
					{
						position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1035
						}
						position++
						goto l1034
					l1035:
						position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
						if buffer[position] != rune('A') {
							goto l1029
						}
						position++
					}
				l1034:
					{
						position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1037
						}
						position++
						goto l1036
					l1037:
						position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
						if buffer[position] != rune('L') {
							goto l1029
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
							goto l1029
						}
						position++
					}
				l1038:
					{
						position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1041
						}
						position++
						goto l1040
					l1041:
						position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
						if buffer[position] != rune('E') {
							goto l1029
						}
						position++
					}
				l1040:
					depth--
					add(rulePegText, position1031)
				}
				if !_rules[ruleAction69]() {
					goto l1029
				}
				depth--
				add(ruleFALSE, position1030)
			}
			return true
		l1029:
			position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
			return false
		},
		/* 94 Wildcard <- <(<((ident ':' !':')? '*')> Action70)> */
		func() bool {
			position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
			{
				position1043 := position
				depth++
				{
					position1044 := position
					depth++
					{
						position1045, tokenIndex1045, depth1045 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1045
						}
						if buffer[position] != rune(':') {
							goto l1045
						}
						position++
						{
							position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1047
							}
							position++
							goto l1045
						l1047:
							position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
						}
						goto l1046
					l1045:
						position, tokenIndex, depth = position1045, tokenIndex1045, depth1045
					}
				l1046:
					if buffer[position] != rune('*') {
						goto l1042
					}
					position++
					depth--
					add(rulePegText, position1044)
				}
				if !_rules[ruleAction70]() {
					goto l1042
				}
				depth--
				add(ruleWildcard, position1043)
			}
			return true
		l1042:
			position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
			return false
		},
		/* 95 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action71)> */
		func() bool {
			position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
			{
				position1049 := position
				depth++
				{
					position1050 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l1048
					}
					position++
				l1051:
					{
						position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
						{
							position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1054
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1054
							}
							position++
							goto l1053
						l1054:
							position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
							{
								position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1055
								}
								position++
								goto l1052
							l1055:
								position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
							}
							if !matchDot() {
								goto l1052
							}
						}
					l1053:
						goto l1051
					l1052:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
					}
					if buffer[position] != rune('\'') {
						goto l1048
					}
					position++
					depth--
					add(rulePegText, position1050)
				}
				if !_rules[ruleAction71]() {
					goto l1048
				}
				depth--
				add(ruleStringLiteral, position1049)
			}
			return true
		l1048:
			position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
			return false
		},
		/* 96 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action72)> */
		func() bool {
			position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
			{
				position1057 := position
				depth++
				{
					position1058 := position
					depth++
					{
						position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1060
						}
						position++
						goto l1059
					l1060:
						position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
						if buffer[position] != rune('I') {
							goto l1056
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
							goto l1056
						}
						position++
					}
				l1061:
					{
						position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1064
						}
						position++
						goto l1063
					l1064:
						position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
						if buffer[position] != rune('T') {
							goto l1056
						}
						position++
					}
				l1063:
					{
						position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1066
						}
						position++
						goto l1065
					l1066:
						position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
						if buffer[position] != rune('R') {
							goto l1056
						}
						position++
					}
				l1065:
					{
						position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1068
						}
						position++
						goto l1067
					l1068:
						position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
						if buffer[position] != rune('E') {
							goto l1056
						}
						position++
					}
				l1067:
					{
						position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1070
						}
						position++
						goto l1069
					l1070:
						position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
						if buffer[position] != rune('A') {
							goto l1056
						}
						position++
					}
				l1069:
					{
						position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1072
						}
						position++
						goto l1071
					l1072:
						position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
						if buffer[position] != rune('M') {
							goto l1056
						}
						position++
					}
				l1071:
					depth--
					add(rulePegText, position1058)
				}
				if !_rules[ruleAction72]() {
					goto l1056
				}
				depth--
				add(ruleISTREAM, position1057)
			}
			return true
		l1056:
			position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
			return false
		},
		/* 97 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action73)> */
		func() bool {
			position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
			{
				position1074 := position
				depth++
				{
					position1075 := position
					depth++
					{
						position1076, tokenIndex1076, depth1076 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1077
						}
						position++
						goto l1076
					l1077:
						position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
						if buffer[position] != rune('D') {
							goto l1073
						}
						position++
					}
				l1076:
					{
						position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1079
						}
						position++
						goto l1078
					l1079:
						position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
						if buffer[position] != rune('S') {
							goto l1073
						}
						position++
					}
				l1078:
					{
						position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1081
						}
						position++
						goto l1080
					l1081:
						position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
						if buffer[position] != rune('T') {
							goto l1073
						}
						position++
					}
				l1080:
					{
						position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1083
						}
						position++
						goto l1082
					l1083:
						position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
						if buffer[position] != rune('R') {
							goto l1073
						}
						position++
					}
				l1082:
					{
						position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1085
						}
						position++
						goto l1084
					l1085:
						position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
						if buffer[position] != rune('E') {
							goto l1073
						}
						position++
					}
				l1084:
					{
						position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1087
						}
						position++
						goto l1086
					l1087:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if buffer[position] != rune('A') {
							goto l1073
						}
						position++
					}
				l1086:
					{
						position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1089
						}
						position++
						goto l1088
					l1089:
						position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
						if buffer[position] != rune('M') {
							goto l1073
						}
						position++
					}
				l1088:
					depth--
					add(rulePegText, position1075)
				}
				if !_rules[ruleAction73]() {
					goto l1073
				}
				depth--
				add(ruleDSTREAM, position1074)
			}
			return true
		l1073:
			position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
			return false
		},
		/* 98 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action74)> */
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
						if buffer[position] != rune('r') {
							goto l1094
						}
						position++
						goto l1093
					l1094:
						position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
						if buffer[position] != rune('R') {
							goto l1090
						}
						position++
					}
				l1093:
					{
						position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1096
						}
						position++
						goto l1095
					l1096:
						position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
						if buffer[position] != rune('S') {
							goto l1090
						}
						position++
					}
				l1095:
					{
						position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1098
						}
						position++
						goto l1097
					l1098:
						position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
						if buffer[position] != rune('T') {
							goto l1090
						}
						position++
					}
				l1097:
					{
						position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1100
						}
						position++
						goto l1099
					l1100:
						position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
						if buffer[position] != rune('R') {
							goto l1090
						}
						position++
					}
				l1099:
					{
						position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1102
						}
						position++
						goto l1101
					l1102:
						position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
						if buffer[position] != rune('E') {
							goto l1090
						}
						position++
					}
				l1101:
					{
						position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1104
						}
						position++
						goto l1103
					l1104:
						position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
						if buffer[position] != rune('A') {
							goto l1090
						}
						position++
					}
				l1103:
					{
						position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1106
						}
						position++
						goto l1105
					l1106:
						position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
						if buffer[position] != rune('M') {
							goto l1090
						}
						position++
					}
				l1105:
					depth--
					add(rulePegText, position1092)
				}
				if !_rules[ruleAction74]() {
					goto l1090
				}
				depth--
				add(ruleRSTREAM, position1091)
			}
			return true
		l1090:
			position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
			return false
		},
		/* 99 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action75)> */
		func() bool {
			position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
			{
				position1108 := position
				depth++
				{
					position1109 := position
					depth++
					{
						position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1111
						}
						position++
						goto l1110
					l1111:
						position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
						if buffer[position] != rune('T') {
							goto l1107
						}
						position++
					}
				l1110:
					{
						position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1113
						}
						position++
						goto l1112
					l1113:
						position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
						if buffer[position] != rune('U') {
							goto l1107
						}
						position++
					}
				l1112:
					{
						position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1115
						}
						position++
						goto l1114
					l1115:
						position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
						if buffer[position] != rune('P') {
							goto l1107
						}
						position++
					}
				l1114:
					{
						position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1117
						}
						position++
						goto l1116
					l1117:
						position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
						if buffer[position] != rune('L') {
							goto l1107
						}
						position++
					}
				l1116:
					{
						position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1119
						}
						position++
						goto l1118
					l1119:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						if buffer[position] != rune('E') {
							goto l1107
						}
						position++
					}
				l1118:
					{
						position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1121
						}
						position++
						goto l1120
					l1121:
						position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
						if buffer[position] != rune('S') {
							goto l1107
						}
						position++
					}
				l1120:
					depth--
					add(rulePegText, position1109)
				}
				if !_rules[ruleAction75]() {
					goto l1107
				}
				depth--
				add(ruleTUPLES, position1108)
			}
			return true
		l1107:
			position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
			return false
		},
		/* 100 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action76)> */
		func() bool {
			position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
			{
				position1123 := position
				depth++
				{
					position1124 := position
					depth++
					{
						position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1126
						}
						position++
						goto l1125
					l1126:
						position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
						if buffer[position] != rune('S') {
							goto l1122
						}
						position++
					}
				l1125:
					{
						position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1128
						}
						position++
						goto l1127
					l1128:
						position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
						if buffer[position] != rune('E') {
							goto l1122
						}
						position++
					}
				l1127:
					{
						position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1130
						}
						position++
						goto l1129
					l1130:
						position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
						if buffer[position] != rune('C') {
							goto l1122
						}
						position++
					}
				l1129:
					{
						position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1132
						}
						position++
						goto l1131
					l1132:
						position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
						if buffer[position] != rune('O') {
							goto l1122
						}
						position++
					}
				l1131:
					{
						position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1134
						}
						position++
						goto l1133
					l1134:
						position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
						if buffer[position] != rune('N') {
							goto l1122
						}
						position++
					}
				l1133:
					{
						position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
						if buffer[position] != rune('D') {
							goto l1122
						}
						position++
					}
				l1135:
					{
						position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
						if buffer[position] != rune('S') {
							goto l1122
						}
						position++
					}
				l1137:
					depth--
					add(rulePegText, position1124)
				}
				if !_rules[ruleAction76]() {
					goto l1122
				}
				depth--
				add(ruleSECONDS, position1123)
			}
			return true
		l1122:
			position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
			return false
		},
		/* 101 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action77)> */
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
						if buffer[position] != rune('m') {
							goto l1143
						}
						position++
						goto l1142
					l1143:
						position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
						if buffer[position] != rune('M') {
							goto l1139
						}
						position++
					}
				l1142:
					{
						position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1145
						}
						position++
						goto l1144
					l1145:
						position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
						if buffer[position] != rune('I') {
							goto l1139
						}
						position++
					}
				l1144:
					{
						position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1147
						}
						position++
						goto l1146
					l1147:
						position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
						if buffer[position] != rune('L') {
							goto l1139
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
							goto l1139
						}
						position++
					}
				l1148:
					{
						position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1151
						}
						position++
						goto l1150
					l1151:
						position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
						if buffer[position] != rune('I') {
							goto l1139
						}
						position++
					}
				l1150:
					{
						position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1153
						}
						position++
						goto l1152
					l1153:
						position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
						if buffer[position] != rune('S') {
							goto l1139
						}
						position++
					}
				l1152:
					{
						position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1155
						}
						position++
						goto l1154
					l1155:
						position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
						if buffer[position] != rune('E') {
							goto l1139
						}
						position++
					}
				l1154:
					{
						position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1157
						}
						position++
						goto l1156
					l1157:
						position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
						if buffer[position] != rune('C') {
							goto l1139
						}
						position++
					}
				l1156:
					{
						position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1159
						}
						position++
						goto l1158
					l1159:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						if buffer[position] != rune('O') {
							goto l1139
						}
						position++
					}
				l1158:
					{
						position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1161
						}
						position++
						goto l1160
					l1161:
						position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
						if buffer[position] != rune('N') {
							goto l1139
						}
						position++
					}
				l1160:
					{
						position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1163
						}
						position++
						goto l1162
					l1163:
						position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
						if buffer[position] != rune('D') {
							goto l1139
						}
						position++
					}
				l1162:
					{
						position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1165
						}
						position++
						goto l1164
					l1165:
						position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
						if buffer[position] != rune('S') {
							goto l1139
						}
						position++
					}
				l1164:
					depth--
					add(rulePegText, position1141)
				}
				if !_rules[ruleAction77]() {
					goto l1139
				}
				depth--
				add(ruleMILLISECONDS, position1140)
			}
			return true
		l1139:
			position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
			return false
		},
		/* 102 StreamIdentifier <- <(<ident> Action78)> */
		func() bool {
			position1166, tokenIndex1166, depth1166 := position, tokenIndex, depth
			{
				position1167 := position
				depth++
				{
					position1168 := position
					depth++
					if !_rules[ruleident]() {
						goto l1166
					}
					depth--
					add(rulePegText, position1168)
				}
				if !_rules[ruleAction78]() {
					goto l1166
				}
				depth--
				add(ruleStreamIdentifier, position1167)
			}
			return true
		l1166:
			position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
			return false
		},
		/* 103 SourceSinkType <- <(<ident> Action79)> */
		func() bool {
			position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
			{
				position1170 := position
				depth++
				{
					position1171 := position
					depth++
					if !_rules[ruleident]() {
						goto l1169
					}
					depth--
					add(rulePegText, position1171)
				}
				if !_rules[ruleAction79]() {
					goto l1169
				}
				depth--
				add(ruleSourceSinkType, position1170)
			}
			return true
		l1169:
			position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
			return false
		},
		/* 104 SourceSinkParamKey <- <(<ident> Action80)> */
		func() bool {
			position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
			{
				position1173 := position
				depth++
				{
					position1174 := position
					depth++
					if !_rules[ruleident]() {
						goto l1172
					}
					depth--
					add(rulePegText, position1174)
				}
				if !_rules[ruleAction80]() {
					goto l1172
				}
				depth--
				add(ruleSourceSinkParamKey, position1173)
			}
			return true
		l1172:
			position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
			return false
		},
		/* 105 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action81)> */
		func() bool {
			position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
			{
				position1176 := position
				depth++
				{
					position1177 := position
					depth++
					{
						position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1179
						}
						position++
						goto l1178
					l1179:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('P') {
							goto l1175
						}
						position++
					}
				l1178:
					{
						position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1181
						}
						position++
						goto l1180
					l1181:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						if buffer[position] != rune('A') {
							goto l1175
						}
						position++
					}
				l1180:
					{
						position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1183
						}
						position++
						goto l1182
					l1183:
						position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
						if buffer[position] != rune('U') {
							goto l1175
						}
						position++
					}
				l1182:
					{
						position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1185
						}
						position++
						goto l1184
					l1185:
						position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
						if buffer[position] != rune('S') {
							goto l1175
						}
						position++
					}
				l1184:
					{
						position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1187
						}
						position++
						goto l1186
					l1187:
						position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
						if buffer[position] != rune('E') {
							goto l1175
						}
						position++
					}
				l1186:
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
							goto l1175
						}
						position++
					}
				l1188:
					depth--
					add(rulePegText, position1177)
				}
				if !_rules[ruleAction81]() {
					goto l1175
				}
				depth--
				add(rulePaused, position1176)
			}
			return true
		l1175:
			position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
			return false
		},
		/* 106 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action82)> */
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
						if buffer[position] != rune('u') {
							goto l1194
						}
						position++
						goto l1193
					l1194:
						position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
						if buffer[position] != rune('U') {
							goto l1190
						}
						position++
					}
				l1193:
					{
						position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1196
						}
						position++
						goto l1195
					l1196:
						position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
						if buffer[position] != rune('N') {
							goto l1190
						}
						position++
					}
				l1195:
					{
						position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1198
						}
						position++
						goto l1197
					l1198:
						position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
						if buffer[position] != rune('P') {
							goto l1190
						}
						position++
					}
				l1197:
					{
						position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1200
						}
						position++
						goto l1199
					l1200:
						position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
						if buffer[position] != rune('A') {
							goto l1190
						}
						position++
					}
				l1199:
					{
						position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1202
						}
						position++
						goto l1201
					l1202:
						position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
						if buffer[position] != rune('U') {
							goto l1190
						}
						position++
					}
				l1201:
					{
						position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1204
						}
						position++
						goto l1203
					l1204:
						position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
						if buffer[position] != rune('S') {
							goto l1190
						}
						position++
					}
				l1203:
					{
						position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1206
						}
						position++
						goto l1205
					l1206:
						position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
						if buffer[position] != rune('E') {
							goto l1190
						}
						position++
					}
				l1205:
					{
						position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1208
						}
						position++
						goto l1207
					l1208:
						position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
						if buffer[position] != rune('D') {
							goto l1190
						}
						position++
					}
				l1207:
					depth--
					add(rulePegText, position1192)
				}
				if !_rules[ruleAction82]() {
					goto l1190
				}
				depth--
				add(ruleUnpaused, position1191)
			}
			return true
		l1190:
			position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
			return false
		},
		/* 107 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
			{
				position1210 := position
				depth++
				{
					position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1212
					}
					goto l1211
				l1212:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if !_rules[ruleInt]() {
						goto l1213
					}
					goto l1211
				l1213:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if !_rules[ruleFloat]() {
						goto l1214
					}
					goto l1211
				l1214:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if !_rules[ruleString]() {
						goto l1215
					}
					goto l1211
				l1215:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if !_rules[ruleBlob]() {
						goto l1216
					}
					goto l1211
				l1216:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if !_rules[ruleTimestamp]() {
						goto l1217
					}
					goto l1211
				l1217:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if !_rules[ruleArray]() {
						goto l1218
					}
					goto l1211
				l1218:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if !_rules[ruleMap]() {
						goto l1209
					}
				}
			l1211:
				depth--
				add(ruleType, position1210)
			}
			return true
		l1209:
			position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
			return false
		},
		/* 108 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action83)> */
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
						if buffer[position] != rune('b') {
							goto l1223
						}
						position++
						goto l1222
					l1223:
						position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
						if buffer[position] != rune('B') {
							goto l1219
						}
						position++
					}
				l1222:
					{
						position1224, tokenIndex1224, depth1224 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1225
						}
						position++
						goto l1224
					l1225:
						position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
						if buffer[position] != rune('O') {
							goto l1219
						}
						position++
					}
				l1224:
					{
						position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1227
						}
						position++
						goto l1226
					l1227:
						position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
						if buffer[position] != rune('O') {
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
					depth--
					add(rulePegText, position1221)
				}
				if !_rules[ruleAction83]() {
					goto l1219
				}
				depth--
				add(ruleBool, position1220)
			}
			return true
		l1219:
			position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
			return false
		},
		/* 109 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action84)> */
		func() bool {
			position1230, tokenIndex1230, depth1230 := position, tokenIndex, depth
			{
				position1231 := position
				depth++
				{
					position1232 := position
					depth++
					{
						position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1234
						}
						position++
						goto l1233
					l1234:
						position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
						if buffer[position] != rune('I') {
							goto l1230
						}
						position++
					}
				l1233:
					{
						position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1236
						}
						position++
						goto l1235
					l1236:
						position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
						if buffer[position] != rune('N') {
							goto l1230
						}
						position++
					}
				l1235:
					{
						position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1238
						}
						position++
						goto l1237
					l1238:
						position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
						if buffer[position] != rune('T') {
							goto l1230
						}
						position++
					}
				l1237:
					depth--
					add(rulePegText, position1232)
				}
				if !_rules[ruleAction84]() {
					goto l1230
				}
				depth--
				add(ruleInt, position1231)
			}
			return true
		l1230:
			position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
			return false
		},
		/* 110 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action85)> */
		func() bool {
			position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
			{
				position1240 := position
				depth++
				{
					position1241 := position
					depth++
					{
						position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1243
						}
						position++
						goto l1242
					l1243:
						position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
						if buffer[position] != rune('F') {
							goto l1239
						}
						position++
					}
				l1242:
					{
						position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1245
						}
						position++
						goto l1244
					l1245:
						position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
						if buffer[position] != rune('L') {
							goto l1239
						}
						position++
					}
				l1244:
					{
						position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1247
						}
						position++
						goto l1246
					l1247:
						position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
						if buffer[position] != rune('O') {
							goto l1239
						}
						position++
					}
				l1246:
					{
						position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1249
						}
						position++
						goto l1248
					l1249:
						position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
						if buffer[position] != rune('A') {
							goto l1239
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
							goto l1239
						}
						position++
					}
				l1250:
					depth--
					add(rulePegText, position1241)
				}
				if !_rules[ruleAction85]() {
					goto l1239
				}
				depth--
				add(ruleFloat, position1240)
			}
			return true
		l1239:
			position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
			return false
		},
		/* 111 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action86)> */
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
						if buffer[position] != rune('s') {
							goto l1256
						}
						position++
						goto l1255
					l1256:
						position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
						if buffer[position] != rune('S') {
							goto l1252
						}
						position++
					}
				l1255:
					{
						position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1258
						}
						position++
						goto l1257
					l1258:
						position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
						if buffer[position] != rune('T') {
							goto l1252
						}
						position++
					}
				l1257:
					{
						position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1260
						}
						position++
						goto l1259
					l1260:
						position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
						if buffer[position] != rune('R') {
							goto l1252
						}
						position++
					}
				l1259:
					{
						position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1262
						}
						position++
						goto l1261
					l1262:
						position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
						if buffer[position] != rune('I') {
							goto l1252
						}
						position++
					}
				l1261:
					{
						position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1264
						}
						position++
						goto l1263
					l1264:
						position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
						if buffer[position] != rune('N') {
							goto l1252
						}
						position++
					}
				l1263:
					{
						position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1266
						}
						position++
						goto l1265
					l1266:
						position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
						if buffer[position] != rune('G') {
							goto l1252
						}
						position++
					}
				l1265:
					depth--
					add(rulePegText, position1254)
				}
				if !_rules[ruleAction86]() {
					goto l1252
				}
				depth--
				add(ruleString, position1253)
			}
			return true
		l1252:
			position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
			return false
		},
		/* 112 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action87)> */
		func() bool {
			position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
			{
				position1268 := position
				depth++
				{
					position1269 := position
					depth++
					{
						position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1271
						}
						position++
						goto l1270
					l1271:
						position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
						if buffer[position] != rune('B') {
							goto l1267
						}
						position++
					}
				l1270:
					{
						position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1273
						}
						position++
						goto l1272
					l1273:
						position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
						if buffer[position] != rune('L') {
							goto l1267
						}
						position++
					}
				l1272:
					{
						position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1275
						}
						position++
						goto l1274
					l1275:
						position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
						if buffer[position] != rune('O') {
							goto l1267
						}
						position++
					}
				l1274:
					{
						position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1277
						}
						position++
						goto l1276
					l1277:
						position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
						if buffer[position] != rune('B') {
							goto l1267
						}
						position++
					}
				l1276:
					depth--
					add(rulePegText, position1269)
				}
				if !_rules[ruleAction87]() {
					goto l1267
				}
				depth--
				add(ruleBlob, position1268)
			}
			return true
		l1267:
			position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
			return false
		},
		/* 113 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action88)> */
		func() bool {
			position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
			{
				position1279 := position
				depth++
				{
					position1280 := position
					depth++
					{
						position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1282
						}
						position++
						goto l1281
					l1282:
						position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
						if buffer[position] != rune('T') {
							goto l1278
						}
						position++
					}
				l1281:
					{
						position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1284
						}
						position++
						goto l1283
					l1284:
						position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
						if buffer[position] != rune('I') {
							goto l1278
						}
						position++
					}
				l1283:
					{
						position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1286
						}
						position++
						goto l1285
					l1286:
						position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
						if buffer[position] != rune('M') {
							goto l1278
						}
						position++
					}
				l1285:
					{
						position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1288
						}
						position++
						goto l1287
					l1288:
						position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
						if buffer[position] != rune('E') {
							goto l1278
						}
						position++
					}
				l1287:
					{
						position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1290
						}
						position++
						goto l1289
					l1290:
						position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
						if buffer[position] != rune('S') {
							goto l1278
						}
						position++
					}
				l1289:
					{
						position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1292
						}
						position++
						goto l1291
					l1292:
						position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
						if buffer[position] != rune('T') {
							goto l1278
						}
						position++
					}
				l1291:
					{
						position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1294
						}
						position++
						goto l1293
					l1294:
						position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
						if buffer[position] != rune('A') {
							goto l1278
						}
						position++
					}
				l1293:
					{
						position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1296
						}
						position++
						goto l1295
					l1296:
						position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
						if buffer[position] != rune('M') {
							goto l1278
						}
						position++
					}
				l1295:
					{
						position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1298
						}
						position++
						goto l1297
					l1298:
						position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
						if buffer[position] != rune('P') {
							goto l1278
						}
						position++
					}
				l1297:
					depth--
					add(rulePegText, position1280)
				}
				if !_rules[ruleAction88]() {
					goto l1278
				}
				depth--
				add(ruleTimestamp, position1279)
			}
			return true
		l1278:
			position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
			return false
		},
		/* 114 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action89)> */
		func() bool {
			position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
			{
				position1300 := position
				depth++
				{
					position1301 := position
					depth++
					{
						position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1303
						}
						position++
						goto l1302
					l1303:
						position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
						if buffer[position] != rune('A') {
							goto l1299
						}
						position++
					}
				l1302:
					{
						position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1305
						}
						position++
						goto l1304
					l1305:
						position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
						if buffer[position] != rune('R') {
							goto l1299
						}
						position++
					}
				l1304:
					{
						position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1307
						}
						position++
						goto l1306
					l1307:
						position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
						if buffer[position] != rune('R') {
							goto l1299
						}
						position++
					}
				l1306:
					{
						position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1309
						}
						position++
						goto l1308
					l1309:
						position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
						if buffer[position] != rune('A') {
							goto l1299
						}
						position++
					}
				l1308:
					{
						position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1311
						}
						position++
						goto l1310
					l1311:
						position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
						if buffer[position] != rune('Y') {
							goto l1299
						}
						position++
					}
				l1310:
					depth--
					add(rulePegText, position1301)
				}
				if !_rules[ruleAction89]() {
					goto l1299
				}
				depth--
				add(ruleArray, position1300)
			}
			return true
		l1299:
			position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
			return false
		},
		/* 115 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action90)> */
		func() bool {
			position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
			{
				position1313 := position
				depth++
				{
					position1314 := position
					depth++
					{
						position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1316
						}
						position++
						goto l1315
					l1316:
						position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
						if buffer[position] != rune('M') {
							goto l1312
						}
						position++
					}
				l1315:
					{
						position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1318
						}
						position++
						goto l1317
					l1318:
						position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
						if buffer[position] != rune('A') {
							goto l1312
						}
						position++
					}
				l1317:
					{
						position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1320
						}
						position++
						goto l1319
					l1320:
						position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
						if buffer[position] != rune('P') {
							goto l1312
						}
						position++
					}
				l1319:
					depth--
					add(rulePegText, position1314)
				}
				if !_rules[ruleAction90]() {
					goto l1312
				}
				depth--
				add(ruleMap, position1313)
			}
			return true
		l1312:
			position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
			return false
		},
		/* 116 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action91)> */
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
						if buffer[position] != rune('o') {
							goto l1325
						}
						position++
						goto l1324
					l1325:
						position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
						if buffer[position] != rune('O') {
							goto l1321
						}
						position++
					}
				l1324:
					{
						position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1327
						}
						position++
						goto l1326
					l1327:
						position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
						if buffer[position] != rune('R') {
							goto l1321
						}
						position++
					}
				l1326:
					depth--
					add(rulePegText, position1323)
				}
				if !_rules[ruleAction91]() {
					goto l1321
				}
				depth--
				add(ruleOr, position1322)
			}
			return true
		l1321:
			position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
			return false
		},
		/* 117 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action92)> */
		func() bool {
			position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
			{
				position1329 := position
				depth++
				{
					position1330 := position
					depth++
					{
						position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1332
						}
						position++
						goto l1331
					l1332:
						position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
						if buffer[position] != rune('A') {
							goto l1328
						}
						position++
					}
				l1331:
					{
						position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1334
						}
						position++
						goto l1333
					l1334:
						position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
						if buffer[position] != rune('N') {
							goto l1328
						}
						position++
					}
				l1333:
					{
						position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1336
						}
						position++
						goto l1335
					l1336:
						position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
						if buffer[position] != rune('D') {
							goto l1328
						}
						position++
					}
				l1335:
					depth--
					add(rulePegText, position1330)
				}
				if !_rules[ruleAction92]() {
					goto l1328
				}
				depth--
				add(ruleAnd, position1329)
			}
			return true
		l1328:
			position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
			return false
		},
		/* 118 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action93)> */
		func() bool {
			position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
			{
				position1338 := position
				depth++
				{
					position1339 := position
					depth++
					{
						position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1341
						}
						position++
						goto l1340
					l1341:
						position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
						if buffer[position] != rune('N') {
							goto l1337
						}
						position++
					}
				l1340:
					{
						position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1343
						}
						position++
						goto l1342
					l1343:
						position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
						if buffer[position] != rune('O') {
							goto l1337
						}
						position++
					}
				l1342:
					{
						position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1345
						}
						position++
						goto l1344
					l1345:
						position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
						if buffer[position] != rune('T') {
							goto l1337
						}
						position++
					}
				l1344:
					depth--
					add(rulePegText, position1339)
				}
				if !_rules[ruleAction93]() {
					goto l1337
				}
				depth--
				add(ruleNot, position1338)
			}
			return true
		l1337:
			position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
			return false
		},
		/* 119 Equal <- <(<'='> Action94)> */
		func() bool {
			position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
			{
				position1347 := position
				depth++
				{
					position1348 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1346
					}
					position++
					depth--
					add(rulePegText, position1348)
				}
				if !_rules[ruleAction94]() {
					goto l1346
				}
				depth--
				add(ruleEqual, position1347)
			}
			return true
		l1346:
			position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
			return false
		},
		/* 120 Less <- <(<'<'> Action95)> */
		func() bool {
			position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
			{
				position1350 := position
				depth++
				{
					position1351 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1349
					}
					position++
					depth--
					add(rulePegText, position1351)
				}
				if !_rules[ruleAction95]() {
					goto l1349
				}
				depth--
				add(ruleLess, position1350)
			}
			return true
		l1349:
			position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
			return false
		},
		/* 121 LessOrEqual <- <(<('<' '=')> Action96)> */
		func() bool {
			position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
			{
				position1353 := position
				depth++
				{
					position1354 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1352
					}
					position++
					if buffer[position] != rune('=') {
						goto l1352
					}
					position++
					depth--
					add(rulePegText, position1354)
				}
				if !_rules[ruleAction96]() {
					goto l1352
				}
				depth--
				add(ruleLessOrEqual, position1353)
			}
			return true
		l1352:
			position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
			return false
		},
		/* 122 Greater <- <(<'>'> Action97)> */
		func() bool {
			position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
			{
				position1356 := position
				depth++
				{
					position1357 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1355
					}
					position++
					depth--
					add(rulePegText, position1357)
				}
				if !_rules[ruleAction97]() {
					goto l1355
				}
				depth--
				add(ruleGreater, position1356)
			}
			return true
		l1355:
			position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
			return false
		},
		/* 123 GreaterOrEqual <- <(<('>' '=')> Action98)> */
		func() bool {
			position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
			{
				position1359 := position
				depth++
				{
					position1360 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1358
					}
					position++
					if buffer[position] != rune('=') {
						goto l1358
					}
					position++
					depth--
					add(rulePegText, position1360)
				}
				if !_rules[ruleAction98]() {
					goto l1358
				}
				depth--
				add(ruleGreaterOrEqual, position1359)
			}
			return true
		l1358:
			position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
			return false
		},
		/* 124 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action99)> */
		func() bool {
			position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
			{
				position1362 := position
				depth++
				{
					position1363 := position
					depth++
					{
						position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1365
						}
						position++
						if buffer[position] != rune('=') {
							goto l1365
						}
						position++
						goto l1364
					l1365:
						position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
						if buffer[position] != rune('<') {
							goto l1361
						}
						position++
						if buffer[position] != rune('>') {
							goto l1361
						}
						position++
					}
				l1364:
					depth--
					add(rulePegText, position1363)
				}
				if !_rules[ruleAction99]() {
					goto l1361
				}
				depth--
				add(ruleNotEqual, position1362)
			}
			return true
		l1361:
			position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
			return false
		},
		/* 125 Concat <- <(<('|' '|')> Action100)> */
		func() bool {
			position1366, tokenIndex1366, depth1366 := position, tokenIndex, depth
			{
				position1367 := position
				depth++
				{
					position1368 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1366
					}
					position++
					if buffer[position] != rune('|') {
						goto l1366
					}
					position++
					depth--
					add(rulePegText, position1368)
				}
				if !_rules[ruleAction100]() {
					goto l1366
				}
				depth--
				add(ruleConcat, position1367)
			}
			return true
		l1366:
			position, tokenIndex, depth = position1366, tokenIndex1366, depth1366
			return false
		},
		/* 126 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action101)> */
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
						if buffer[position] != rune('i') {
							goto l1373
						}
						position++
						goto l1372
					l1373:
						position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
						if buffer[position] != rune('I') {
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
					depth--
					add(rulePegText, position1371)
				}
				if !_rules[ruleAction101]() {
					goto l1369
				}
				depth--
				add(ruleIs, position1370)
			}
			return true
		l1369:
			position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
			return false
		},
		/* 127 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action102)> */
		func() bool {
			position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
			{
				position1377 := position
				depth++
				{
					position1378 := position
					depth++
					{
						position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1380
						}
						position++
						goto l1379
					l1380:
						position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
						if buffer[position] != rune('I') {
							goto l1376
						}
						position++
					}
				l1379:
					{
						position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1382
						}
						position++
						goto l1381
					l1382:
						position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
						if buffer[position] != rune('S') {
							goto l1376
						}
						position++
					}
				l1381:
					if !_rules[rulesp]() {
						goto l1376
					}
					{
						position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1384
						}
						position++
						goto l1383
					l1384:
						position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
						if buffer[position] != rune('N') {
							goto l1376
						}
						position++
					}
				l1383:
					{
						position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1386
						}
						position++
						goto l1385
					l1386:
						position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
						if buffer[position] != rune('O') {
							goto l1376
						}
						position++
					}
				l1385:
					{
						position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1388
						}
						position++
						goto l1387
					l1388:
						position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
						if buffer[position] != rune('T') {
							goto l1376
						}
						position++
					}
				l1387:
					depth--
					add(rulePegText, position1378)
				}
				if !_rules[ruleAction102]() {
					goto l1376
				}
				depth--
				add(ruleIsNot, position1377)
			}
			return true
		l1376:
			position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
			return false
		},
		/* 128 Plus <- <(<'+'> Action103)> */
		func() bool {
			position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
			{
				position1390 := position
				depth++
				{
					position1391 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1389
					}
					position++
					depth--
					add(rulePegText, position1391)
				}
				if !_rules[ruleAction103]() {
					goto l1389
				}
				depth--
				add(rulePlus, position1390)
			}
			return true
		l1389:
			position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
			return false
		},
		/* 129 Minus <- <(<'-'> Action104)> */
		func() bool {
			position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
			{
				position1393 := position
				depth++
				{
					position1394 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1392
					}
					position++
					depth--
					add(rulePegText, position1394)
				}
				if !_rules[ruleAction104]() {
					goto l1392
				}
				depth--
				add(ruleMinus, position1393)
			}
			return true
		l1392:
			position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
			return false
		},
		/* 130 Multiply <- <(<'*'> Action105)> */
		func() bool {
			position1395, tokenIndex1395, depth1395 := position, tokenIndex, depth
			{
				position1396 := position
				depth++
				{
					position1397 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1395
					}
					position++
					depth--
					add(rulePegText, position1397)
				}
				if !_rules[ruleAction105]() {
					goto l1395
				}
				depth--
				add(ruleMultiply, position1396)
			}
			return true
		l1395:
			position, tokenIndex, depth = position1395, tokenIndex1395, depth1395
			return false
		},
		/* 131 Divide <- <(<'/'> Action106)> */
		func() bool {
			position1398, tokenIndex1398, depth1398 := position, tokenIndex, depth
			{
				position1399 := position
				depth++
				{
					position1400 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1398
					}
					position++
					depth--
					add(rulePegText, position1400)
				}
				if !_rules[ruleAction106]() {
					goto l1398
				}
				depth--
				add(ruleDivide, position1399)
			}
			return true
		l1398:
			position, tokenIndex, depth = position1398, tokenIndex1398, depth1398
			return false
		},
		/* 132 Modulo <- <(<'%'> Action107)> */
		func() bool {
			position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
			{
				position1402 := position
				depth++
				{
					position1403 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1401
					}
					position++
					depth--
					add(rulePegText, position1403)
				}
				if !_rules[ruleAction107]() {
					goto l1401
				}
				depth--
				add(ruleModulo, position1402)
			}
			return true
		l1401:
			position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
			return false
		},
		/* 133 UnaryMinus <- <(<'-'> Action108)> */
		func() bool {
			position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
			{
				position1405 := position
				depth++
				{
					position1406 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1404
					}
					position++
					depth--
					add(rulePegText, position1406)
				}
				if !_rules[ruleAction108]() {
					goto l1404
				}
				depth--
				add(ruleUnaryMinus, position1405)
			}
			return true
		l1404:
			position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
			return false
		},
		/* 134 Identifier <- <(<ident> Action109)> */
		func() bool {
			position1407, tokenIndex1407, depth1407 := position, tokenIndex, depth
			{
				position1408 := position
				depth++
				{
					position1409 := position
					depth++
					if !_rules[ruleident]() {
						goto l1407
					}
					depth--
					add(rulePegText, position1409)
				}
				if !_rules[ruleAction109]() {
					goto l1407
				}
				depth--
				add(ruleIdentifier, position1408)
			}
			return true
		l1407:
			position, tokenIndex, depth = position1407, tokenIndex1407, depth1407
			return false
		},
		/* 135 TargetIdentifier <- <(<jsonPath> Action110)> */
		func() bool {
			position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
			{
				position1411 := position
				depth++
				{
					position1412 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1410
					}
					depth--
					add(rulePegText, position1412)
				}
				if !_rules[ruleAction110]() {
					goto l1410
				}
				depth--
				add(ruleTargetIdentifier, position1411)
			}
			return true
		l1410:
			position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
			return false
		},
		/* 136 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1413, tokenIndex1413, depth1413 := position, tokenIndex, depth
			{
				position1414 := position
				depth++
				{
					position1415, tokenIndex1415, depth1415 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1416
					}
					position++
					goto l1415
				l1416:
					position, tokenIndex, depth = position1415, tokenIndex1415, depth1415
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1413
					}
					position++
				}
			l1415:
			l1417:
				{
					position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
					{
						position1419, tokenIndex1419, depth1419 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1420
						}
						position++
						goto l1419
					l1420:
						position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1421
						}
						position++
						goto l1419
					l1421:
						position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1422
						}
						position++
						goto l1419
					l1422:
						position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
						if buffer[position] != rune('_') {
							goto l1418
						}
						position++
					}
				l1419:
					goto l1417
				l1418:
					position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
				}
				depth--
				add(ruleident, position1414)
			}
			return true
		l1413:
			position, tokenIndex, depth = position1413, tokenIndex1413, depth1413
			return false
		},
		/* 137 jsonPath <- <(jsonPathHead jsonPathNonHead*)> */
		func() bool {
			position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
			{
				position1424 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1423
				}
			l1425:
				{
					position1426, tokenIndex1426, depth1426 := position, tokenIndex, depth
					if !_rules[rulejsonPathNonHead]() {
						goto l1426
					}
					goto l1425
				l1426:
					position, tokenIndex, depth = position1426, tokenIndex1426, depth1426
				}
				depth--
				add(rulejsonPath, position1424)
			}
			return true
		l1423:
			position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
			return false
		},
		/* 138 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
			{
				position1428 := position
				depth++
				{
					position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1430
					}
					goto l1429
				l1430:
					position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1427
					}
				}
			l1429:
				depth--
				add(rulejsonPathHead, position1428)
			}
			return true
		l1427:
			position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
			return false
		},
		/* 139 jsonPathNonHead <- <(('.' jsonMapAccessString) / jsonMapAccessBracket / jsonArrayAccess)> */
		func() bool {
			position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
			{
				position1432 := position
				depth++
				{
					position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1434
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1434
					}
					goto l1433
				l1434:
					position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1435
					}
					goto l1433
				l1435:
					position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
					if !_rules[rulejsonArrayAccess]() {
						goto l1431
					}
				}
			l1433:
				depth--
				add(rulejsonPathNonHead, position1432)
			}
			return true
		l1431:
			position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
			return false
		},
		/* 140 jsonMapAccessString <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
			{
				position1437 := position
				depth++
				{
					position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1439
					}
					position++
					goto l1438
				l1439:
					position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1436
					}
					position++
				}
			l1438:
			l1440:
				{
					position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
					{
						position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1443
						}
						position++
						goto l1442
					l1443:
						position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1444
						}
						position++
						goto l1442
					l1444:
						position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1445
						}
						position++
						goto l1442
					l1445:
						position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
						if buffer[position] != rune('_') {
							goto l1441
						}
						position++
					}
				l1442:
					goto l1440
				l1441:
					position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
				}
				depth--
				add(rulejsonMapAccessString, position1437)
			}
			return true
		l1436:
			position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
			return false
		},
		/* 141 jsonMapAccessBracket <- <('[' '\'' (('\'' '\'') / (!'\'' .))* '\'' ']')> */
		func() bool {
			position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
			{
				position1447 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1446
				}
				position++
				if buffer[position] != rune('\'') {
					goto l1446
				}
				position++
			l1448:
				{
					position1449, tokenIndex1449, depth1449 := position, tokenIndex, depth
					{
						position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
						if buffer[position] != rune('\'') {
							goto l1451
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1451
						}
						position++
						goto l1450
					l1451:
						position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
						{
							position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1452
							}
							position++
							goto l1449
						l1452:
							position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
						}
						if !matchDot() {
							goto l1449
						}
					}
				l1450:
					goto l1448
				l1449:
					position, tokenIndex, depth = position1449, tokenIndex1449, depth1449
				}
				if buffer[position] != rune('\'') {
					goto l1446
				}
				position++
				if buffer[position] != rune(']') {
					goto l1446
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1447)
			}
			return true
		l1446:
			position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
			return false
		},
		/* 142 jsonArrayAccess <- <('[' [0-9]+ ']')> */
		func() bool {
			position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
			{
				position1454 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1453
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1453
				}
				position++
			l1455:
				{
					position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1456
					}
					position++
					goto l1455
				l1456:
					position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
				}
				if buffer[position] != rune(']') {
					goto l1453
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1454)
			}
			return true
		l1453:
			position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
			return false
		},
		/* 143 sp <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)*> */
		func() bool {
			{
				position1458 := position
				depth++
			l1459:
				{
					position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
					{
						position1461, tokenIndex1461, depth1461 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1462
						}
						position++
						goto l1461
					l1462:
						position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
						if buffer[position] != rune('\t') {
							goto l1463
						}
						position++
						goto l1461
					l1463:
						position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
						if buffer[position] != rune('\n') {
							goto l1464
						}
						position++
						goto l1461
					l1464:
						position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
						if buffer[position] != rune('\r') {
							goto l1465
						}
						position++
						goto l1461
					l1465:
						position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
						if !_rules[rulecomment]() {
							goto l1466
						}
						goto l1461
					l1466:
						position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
						if !_rules[rulefinalComment]() {
							goto l1460
						}
					}
				l1461:
					goto l1459
				l1460:
					position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
				}
				depth--
				add(rulesp, position1458)
			}
			return true
		},
		/* 144 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
			{
				position1468 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1467
				}
				position++
				if buffer[position] != rune('-') {
					goto l1467
				}
				position++
			l1469:
				{
					position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
					{
						position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
						{
							position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1473
							}
							position++
							goto l1472
						l1473:
							position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
							if buffer[position] != rune('\n') {
								goto l1471
							}
							position++
						}
					l1472:
						goto l1470
					l1471:
						position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
					}
					if !matchDot() {
						goto l1470
					}
					goto l1469
				l1470:
					position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
				}
				{
					position1474, tokenIndex1474, depth1474 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1475
					}
					position++
					goto l1474
				l1475:
					position, tokenIndex, depth = position1474, tokenIndex1474, depth1474
					if buffer[position] != rune('\n') {
						goto l1467
					}
					position++
				}
			l1474:
				depth--
				add(rulecomment, position1468)
			}
			return true
		l1467:
			position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
			return false
		},
		/* 145 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1476, tokenIndex1476, depth1476 := position, tokenIndex, depth
			{
				position1477 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1476
				}
				position++
				if buffer[position] != rune('-') {
					goto l1476
				}
				position++
			l1478:
				{
					position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
					{
						position1480, tokenIndex1480, depth1480 := position, tokenIndex, depth
						{
							position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1482
							}
							position++
							goto l1481
						l1482:
							position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
							if buffer[position] != rune('\n') {
								goto l1480
							}
							position++
						}
					l1481:
						goto l1479
					l1480:
						position, tokenIndex, depth = position1480, tokenIndex1480, depth1480
					}
					if !matchDot() {
						goto l1479
					}
					goto l1478
				l1479:
					position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
				}
				{
					position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
					if !matchDot() {
						goto l1483
					}
					goto l1476
				l1483:
					position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
				}
				depth--
				add(rulefinalComment, position1477)
			}
			return true
		l1476:
			position, tokenIndex, depth = position1476, tokenIndex1476, depth1476
			return false
		},
		nil,
		/* 148 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 149 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 150 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 151 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 152 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 153 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 154 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 155 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 156 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 157 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 158 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 159 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 160 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 161 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 162 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 163 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 164 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 165 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 166 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 167 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 168 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 169 Action21 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 170 Action22 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 171 Action23 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 172 Action24 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 173 Action25 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 174 Action26 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 175 Action27 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 176 Action28 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 177 Action29 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 178 Action30 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 179 Action31 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 180 Action32 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 181 Action33 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 182 Action34 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 183 Action35 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 184 Action36 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 185 Action37 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 186 Action38 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 187 Action39 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 188 Action40 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 189 Action41 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 190 Action42 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 191 Action43 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 192 Action44 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 193 Action45 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 194 Action46 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 195 Action47 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 196 Action48 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 197 Action49 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 198 Action50 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 199 Action51 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 200 Action52 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 201 Action53 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 202 Action54 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 203 Action55 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 204 Action56 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 205 Action57 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 206 Action58 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 207 Action59 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 208 Action60 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 209 Action61 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 210 Action62 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 211 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 212 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 213 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 214 Action66 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 215 Action67 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 216 Action68 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 217 Action69 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 218 Action70 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 219 Action71 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 220 Action72 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 221 Action73 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 222 Action74 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 223 Action75 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 224 Action76 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 225 Action77 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 226 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 227 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 228 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 229 Action81 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 230 Action82 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 231 Action83 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 232 Action84 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 233 Action85 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 234 Action86 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 235 Action87 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 236 Action88 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 237 Action89 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 238 Action90 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 239 Action91 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 240 Action92 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 241 Action93 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 242 Action94 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 243 Action95 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 244 Action96 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 245 Action97 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 246 Action98 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 247 Action99 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 248 Action100 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 249 Action101 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 250 Action102 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 251 Action103 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 252 Action104 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 253 Action105 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 254 Action106 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 255 Action107 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 256 Action108 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 257 Action109 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 258 Action110 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction110, position)
			}
			return true
		},
	}
	p.rules = _rules
}

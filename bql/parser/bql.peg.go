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
	ruleStar
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
	"Star",
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
	rules  [258]func() bool
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

			p.AssembleEmitter()

		case ruleAction24:

			p.AssembleEmitterOptions(begin, end)

		case ruleAction25:

			p.AssembleEmitterLimit()

		case ruleAction26:

			p.AssembleProjections(begin, end)

		case ruleAction27:

			p.AssembleAlias()

		case ruleAction28:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction29:

			p.AssembleInterval()

		case ruleAction30:

			p.AssembleInterval()

		case ruleAction31:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction32:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction33:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction34:

			p.EnsureAliasedStreamWindow()

		case ruleAction35:

			p.AssembleAliasedStreamWindow()

		case ruleAction36:

			p.AssembleStreamWindow()

		case ruleAction37:

			p.AssembleUDSFFuncApp()

		case ruleAction38:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction39:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction40:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction41:

			p.AssembleSourceSinkParam()

		case ruleAction42:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction43:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction44:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction45:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction46:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction47:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction48:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction49:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction50:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction51:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction52:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction53:

			p.AssembleTypeCast(begin, end)

		case ruleAction54:

			p.AssembleTypeCast(begin, end)

		case ruleAction55:

			p.AssembleFuncApp()

		case ruleAction56:

			p.AssembleExpressions(begin, end)

		case ruleAction57:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction58:

			p.AssembleMap(begin, end)

		case ruleAction59:

			p.AssembleKeyValuePair()

		case ruleAction60:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction61:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction62:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction63:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction64:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction65:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction66:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction67:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction68:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction69:

			p.PushComponent(begin, end, NewWildcard(""))

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
		/* 6 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt / LoadStateOrCreateStmt / LoadStateStmt)> */
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
			position41, tokenIndex41, depth41 := position, tokenIndex, depth
			{
				position42 := position
				depth++
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l44
					}
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l45
					}
					goto l43
				l45:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if !_rules[ruleDropStreamStmt]() {
						goto l46
					}
					goto l43
				l46:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l47
					}
					goto l43
				l47:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l41
					}
				}
			l43:
				depth--
				add(ruleStreamStmt, position42)
			}
			return true
		l41:
			position, tokenIndex, depth = position41, tokenIndex41, depth41
			return false
		},
		/* 8 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action2)> */
		func() bool {
			position48, tokenIndex48, depth48 := position, tokenIndex, depth
			{
				position49 := position
				depth++
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('S') {
						goto l48
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('E') {
						goto l48
					}
					position++
				}
			l52:
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('L') {
						goto l48
					}
					position++
				}
			l54:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('E') {
						goto l48
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('C') {
						goto l48
					}
					position++
				}
			l58:
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l61
					}
					position++
					goto l60
				l61:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
					if buffer[position] != rune('T') {
						goto l48
					}
					position++
				}
			l60:
				if !_rules[rulesp]() {
					goto l48
				}
				if !_rules[ruleEmitter]() {
					goto l48
				}
				if !_rules[rulesp]() {
					goto l48
				}
				if !_rules[ruleProjections]() {
					goto l48
				}
				if !_rules[rulesp]() {
					goto l48
				}
				if !_rules[ruleWindowedFrom]() {
					goto l48
				}
				if !_rules[rulesp]() {
					goto l48
				}
				if !_rules[ruleFilter]() {
					goto l48
				}
				if !_rules[rulesp]() {
					goto l48
				}
				if !_rules[ruleGrouping]() {
					goto l48
				}
				if !_rules[rulesp]() {
					goto l48
				}
				if !_rules[ruleHaving]() {
					goto l48
				}
				if !_rules[rulesp]() {
					goto l48
				}
				if !_rules[ruleAction2]() {
					goto l48
				}
				depth--
				add(ruleSelectStmt, position49)
			}
			return true
		l48:
			position, tokenIndex, depth = position48, tokenIndex48, depth48
			return false
		},
		/* 9 SelectUnionStmt <- <(<(SelectStmt (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action3)> */
		func() bool {
			position62, tokenIndex62, depth62 := position, tokenIndex, depth
			{
				position63 := position
				depth++
				{
					position64 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l62
					}
					{
						position67, tokenIndex67, depth67 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l68
						}
						position++
						goto l67
					l68:
						position, tokenIndex, depth = position67, tokenIndex67, depth67
						if buffer[position] != rune('U') {
							goto l62
						}
						position++
					}
				l67:
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l70
						}
						position++
						goto l69
					l70:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
						if buffer[position] != rune('N') {
							goto l62
						}
						position++
					}
				l69:
					{
						position71, tokenIndex71, depth71 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l72
						}
						position++
						goto l71
					l72:
						position, tokenIndex, depth = position71, tokenIndex71, depth71
						if buffer[position] != rune('I') {
							goto l62
						}
						position++
					}
				l71:
					{
						position73, tokenIndex73, depth73 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l74
						}
						position++
						goto l73
					l74:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
						if buffer[position] != rune('O') {
							goto l62
						}
						position++
					}
				l73:
					{
						position75, tokenIndex75, depth75 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l76
						}
						position++
						goto l75
					l76:
						position, tokenIndex, depth = position75, tokenIndex75, depth75
						if buffer[position] != rune('N') {
							goto l62
						}
						position++
					}
				l75:
					if !_rules[rulesp]() {
						goto l62
					}
					{
						position77, tokenIndex77, depth77 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l78
						}
						position++
						goto l77
					l78:
						position, tokenIndex, depth = position77, tokenIndex77, depth77
						if buffer[position] != rune('A') {
							goto l62
						}
						position++
					}
				l77:
					{
						position79, tokenIndex79, depth79 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l80
						}
						position++
						goto l79
					l80:
						position, tokenIndex, depth = position79, tokenIndex79, depth79
						if buffer[position] != rune('L') {
							goto l62
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
							goto l62
						}
						position++
					}
				l81:
					if !_rules[rulesp]() {
						goto l62
					}
					if !_rules[ruleSelectStmt]() {
						goto l62
					}
				l65:
					{
						position66, tokenIndex66, depth66 := position, tokenIndex, depth
						{
							position83, tokenIndex83, depth83 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l84
							}
							position++
							goto l83
						l84:
							position, tokenIndex, depth = position83, tokenIndex83, depth83
							if buffer[position] != rune('U') {
								goto l66
							}
							position++
						}
					l83:
						{
							position85, tokenIndex85, depth85 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l86
							}
							position++
							goto l85
						l86:
							position, tokenIndex, depth = position85, tokenIndex85, depth85
							if buffer[position] != rune('N') {
								goto l66
							}
							position++
						}
					l85:
						{
							position87, tokenIndex87, depth87 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l88
							}
							position++
							goto l87
						l88:
							position, tokenIndex, depth = position87, tokenIndex87, depth87
							if buffer[position] != rune('I') {
								goto l66
							}
							position++
						}
					l87:
						{
							position89, tokenIndex89, depth89 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l90
							}
							position++
							goto l89
						l90:
							position, tokenIndex, depth = position89, tokenIndex89, depth89
							if buffer[position] != rune('O') {
								goto l66
							}
							position++
						}
					l89:
						{
							position91, tokenIndex91, depth91 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l92
							}
							position++
							goto l91
						l92:
							position, tokenIndex, depth = position91, tokenIndex91, depth91
							if buffer[position] != rune('N') {
								goto l66
							}
							position++
						}
					l91:
						if !_rules[rulesp]() {
							goto l66
						}
						{
							position93, tokenIndex93, depth93 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l94
							}
							position++
							goto l93
						l94:
							position, tokenIndex, depth = position93, tokenIndex93, depth93
							if buffer[position] != rune('A') {
								goto l66
							}
							position++
						}
					l93:
						{
							position95, tokenIndex95, depth95 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l96
							}
							position++
							goto l95
						l96:
							position, tokenIndex, depth = position95, tokenIndex95, depth95
							if buffer[position] != rune('L') {
								goto l66
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
								goto l66
							}
							position++
						}
					l97:
						if !_rules[rulesp]() {
							goto l66
						}
						if !_rules[ruleSelectStmt]() {
							goto l66
						}
						goto l65
					l66:
						position, tokenIndex, depth = position66, tokenIndex66, depth66
					}
					depth--
					add(rulePegText, position64)
				}
				if !_rules[ruleAction3]() {
					goto l62
				}
				depth--
				add(ruleSelectUnionStmt, position63)
			}
			return true
		l62:
			position, tokenIndex, depth = position62, tokenIndex62, depth62
			return false
		},
		/* 10 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action4)> */
		func() bool {
			position99, tokenIndex99, depth99 := position, tokenIndex, depth
			{
				position100 := position
				depth++
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('C') {
						goto l99
					}
					position++
				}
			l101:
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l104
					}
					position++
					goto l103
				l104:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
					if buffer[position] != rune('R') {
						goto l99
					}
					position++
				}
			l103:
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('E') {
						goto l99
					}
					position++
				}
			l105:
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l108
					}
					position++
					goto l107
				l108:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
					if buffer[position] != rune('A') {
						goto l99
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('T') {
						goto l99
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
						goto l99
					}
					position++
				}
			l111:
				if !_rules[rulesp]() {
					goto l99
				}
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l114
					}
					position++
					goto l113
				l114:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
					if buffer[position] != rune('S') {
						goto l99
					}
					position++
				}
			l113:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('T') {
						goto l99
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('R') {
						goto l99
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('E') {
						goto l99
					}
					position++
				}
			l119:
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if buffer[position] != rune('A') {
						goto l99
					}
					position++
				}
			l121:
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('M') {
						goto l99
					}
					position++
				}
			l123:
				if !_rules[rulesp]() {
					goto l99
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l99
				}
				if !_rules[rulesp]() {
					goto l99
				}
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('A') {
						goto l99
					}
					position++
				}
			l125:
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('S') {
						goto l99
					}
					position++
				}
			l127:
				if !_rules[rulesp]() {
					goto l99
				}
				if !_rules[ruleSelectStmt]() {
					goto l99
				}
				if !_rules[ruleAction4]() {
					goto l99
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position100)
			}
			return true
		l99:
			position, tokenIndex, depth = position99, tokenIndex99, depth99
			return false
		},
		/* 11 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action5)> */
		func() bool {
			position129, tokenIndex129, depth129 := position, tokenIndex, depth
			{
				position130 := position
				depth++
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('C') {
						goto l129
					}
					position++
				}
			l131:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l134
					}
					position++
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					if buffer[position] != rune('R') {
						goto l129
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('E') {
						goto l129
					}
					position++
				}
			l135:
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l138
					}
					position++
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if buffer[position] != rune('A') {
						goto l129
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('T') {
						goto l129
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
						goto l129
					}
					position++
				}
			l141:
				if !_rules[rulesp]() {
					goto l129
				}
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					if buffer[position] != rune('S') {
						goto l129
					}
					position++
				}
			l143:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('T') {
						goto l129
					}
					position++
				}
			l145:
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('R') {
						goto l129
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('E') {
						goto l129
					}
					position++
				}
			l149:
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('A') {
						goto l129
					}
					position++
				}
			l151:
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('M') {
						goto l129
					}
					position++
				}
			l153:
				if !_rules[rulesp]() {
					goto l129
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l129
				}
				if !_rules[rulesp]() {
					goto l129
				}
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('A') {
						goto l129
					}
					position++
				}
			l155:
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('S') {
						goto l129
					}
					position++
				}
			l157:
				if !_rules[rulesp]() {
					goto l129
				}
				if !_rules[ruleSelectUnionStmt]() {
					goto l129
				}
				if !_rules[ruleAction5]() {
					goto l129
				}
				depth--
				add(ruleCreateStreamAsSelectUnionStmt, position130)
			}
			return true
		l129:
			position, tokenIndex, depth = position129, tokenIndex129, depth129
			return false
		},
		/* 12 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action6)> */
		func() bool {
			position159, tokenIndex159, depth159 := position, tokenIndex, depth
			{
				position160 := position
				depth++
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l162
					}
					position++
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if buffer[position] != rune('C') {
						goto l159
					}
					position++
				}
			l161:
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l164
					}
					position++
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if buffer[position] != rune('R') {
						goto l159
					}
					position++
				}
			l163:
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('E') {
						goto l159
					}
					position++
				}
			l165:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('A') {
						goto l159
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('T') {
						goto l159
					}
					position++
				}
			l169:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('E') {
						goto l159
					}
					position++
				}
			l171:
				if !_rules[rulesp]() {
					goto l159
				}
				if !_rules[rulePausedOpt]() {
					goto l159
				}
				if !_rules[rulesp]() {
					goto l159
				}
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('S') {
						goto l159
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('O') {
						goto l159
					}
					position++
				}
			l175:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('U') {
						goto l159
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('R') {
						goto l159
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('C') {
						goto l159
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
						goto l159
					}
					position++
				}
			l183:
				if !_rules[rulesp]() {
					goto l159
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l159
				}
				if !_rules[rulesp]() {
					goto l159
				}
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('T') {
						goto l159
					}
					position++
				}
			l185:
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if buffer[position] != rune('Y') {
						goto l159
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('P') {
						goto l159
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
						goto l159
					}
					position++
				}
			l191:
				if !_rules[rulesp]() {
					goto l159
				}
				if !_rules[ruleSourceSinkType]() {
					goto l159
				}
				if !_rules[rulesp]() {
					goto l159
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l159
				}
				if !_rules[ruleAction6]() {
					goto l159
				}
				depth--
				add(ruleCreateSourceStmt, position160)
			}
			return true
		l159:
			position, tokenIndex, depth = position159, tokenIndex159, depth159
			return false
		},
		/* 13 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action7)> */
		func() bool {
			position193, tokenIndex193, depth193 := position, tokenIndex, depth
			{
				position194 := position
				depth++
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('C') {
						goto l193
					}
					position++
				}
			l195:
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l198
					}
					position++
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if buffer[position] != rune('R') {
						goto l193
					}
					position++
				}
			l197:
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('E') {
						goto l193
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('A') {
						goto l193
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('T') {
						goto l193
					}
					position++
				}
			l203:
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('E') {
						goto l193
					}
					position++
				}
			l205:
				if !_rules[rulesp]() {
					goto l193
				}
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('S') {
						goto l193
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('I') {
						goto l193
					}
					position++
				}
			l209:
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('N') {
						goto l193
					}
					position++
				}
			l211:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('K') {
						goto l193
					}
					position++
				}
			l213:
				if !_rules[rulesp]() {
					goto l193
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l193
				}
				if !_rules[rulesp]() {
					goto l193
				}
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('T') {
						goto l193
					}
					position++
				}
			l215:
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('Y') {
						goto l193
					}
					position++
				}
			l217:
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('P') {
						goto l193
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
						goto l193
					}
					position++
				}
			l221:
				if !_rules[rulesp]() {
					goto l193
				}
				if !_rules[ruleSourceSinkType]() {
					goto l193
				}
				if !_rules[rulesp]() {
					goto l193
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l193
				}
				if !_rules[ruleAction7]() {
					goto l193
				}
				depth--
				add(ruleCreateSinkStmt, position194)
			}
			return true
		l193:
			position, tokenIndex, depth = position193, tokenIndex193, depth193
			return false
		},
		/* 14 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action8)> */
		func() bool {
			position223, tokenIndex223, depth223 := position, tokenIndex, depth
			{
				position224 := position
				depth++
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('C') {
						goto l223
					}
					position++
				}
			l225:
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l228
					}
					position++
					goto l227
				l228:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
					if buffer[position] != rune('R') {
						goto l223
					}
					position++
				}
			l227:
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('E') {
						goto l223
					}
					position++
				}
			l229:
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l232
					}
					position++
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if buffer[position] != rune('A') {
						goto l223
					}
					position++
				}
			l231:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('T') {
						goto l223
					}
					position++
				}
			l233:
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l236
					}
					position++
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if buffer[position] != rune('E') {
						goto l223
					}
					position++
				}
			l235:
				if !_rules[rulesp]() {
					goto l223
				}
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if buffer[position] != rune('S') {
						goto l223
					}
					position++
				}
			l237:
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
						goto l223
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('A') {
						goto l223
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('T') {
						goto l223
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
						goto l223
					}
					position++
				}
			l245:
				if !_rules[rulesp]() {
					goto l223
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l223
				}
				if !_rules[rulesp]() {
					goto l223
				}
				{
					position247, tokenIndex247, depth247 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l248
					}
					position++
					goto l247
				l248:
					position, tokenIndex, depth = position247, tokenIndex247, depth247
					if buffer[position] != rune('T') {
						goto l223
					}
					position++
				}
			l247:
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('Y') {
						goto l223
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
						goto l223
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('E') {
						goto l223
					}
					position++
				}
			l253:
				if !_rules[rulesp]() {
					goto l223
				}
				if !_rules[ruleSourceSinkType]() {
					goto l223
				}
				if !_rules[rulesp]() {
					goto l223
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l223
				}
				if !_rules[ruleAction8]() {
					goto l223
				}
				depth--
				add(ruleCreateStateStmt, position224)
			}
			return true
		l223:
			position, tokenIndex, depth = position223, tokenIndex223, depth223
			return false
		},
		/* 15 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position255, tokenIndex255, depth255 := position, tokenIndex, depth
			{
				position256 := position
				depth++
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('U') {
						goto l255
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('P') {
						goto l255
					}
					position++
				}
			l259:
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('D') {
						goto l255
					}
					position++
				}
			l261:
				{
					position263, tokenIndex263, depth263 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l264
					}
					position++
					goto l263
				l264:
					position, tokenIndex, depth = position263, tokenIndex263, depth263
					if buffer[position] != rune('A') {
						goto l255
					}
					position++
				}
			l263:
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('T') {
						goto l255
					}
					position++
				}
			l265:
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l268
					}
					position++
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if buffer[position] != rune('E') {
						goto l255
					}
					position++
				}
			l267:
				if !_rules[rulesp]() {
					goto l255
				}
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('S') {
						goto l255
					}
					position++
				}
			l269:
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('T') {
						goto l255
					}
					position++
				}
			l271:
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('A') {
						goto l255
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('T') {
						goto l255
					}
					position++
				}
			l275:
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('E') {
						goto l255
					}
					position++
				}
			l277:
				if !_rules[rulesp]() {
					goto l255
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l255
				}
				if !_rules[rulesp]() {
					goto l255
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l255
				}
				if !_rules[ruleAction9]() {
					goto l255
				}
				depth--
				add(ruleUpdateStateStmt, position256)
			}
			return true
		l255:
			position, tokenIndex, depth = position255, tokenIndex255, depth255
			return false
		},
		/* 16 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action10)> */
		func() bool {
			position279, tokenIndex279, depth279 := position, tokenIndex, depth
			{
				position280 := position
				depth++
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('U') {
						goto l279
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('P') {
						goto l279
					}
					position++
				}
			l283:
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('D') {
						goto l279
					}
					position++
				}
			l285:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('A') {
						goto l279
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('T') {
						goto l279
					}
					position++
				}
			l289:
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('E') {
						goto l279
					}
					position++
				}
			l291:
				if !_rules[rulesp]() {
					goto l279
				}
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					if buffer[position] != rune('S') {
						goto l279
					}
					position++
				}
			l293:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('O') {
						goto l279
					}
					position++
				}
			l295:
				{
					position297, tokenIndex297, depth297 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l298
					}
					position++
					goto l297
				l298:
					position, tokenIndex, depth = position297, tokenIndex297, depth297
					if buffer[position] != rune('U') {
						goto l279
					}
					position++
				}
			l297:
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l300
					}
					position++
					goto l299
				l300:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
					if buffer[position] != rune('R') {
						goto l279
					}
					position++
				}
			l299:
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l302
					}
					position++
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if buffer[position] != rune('C') {
						goto l279
					}
					position++
				}
			l301:
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('E') {
						goto l279
					}
					position++
				}
			l303:
				if !_rules[rulesp]() {
					goto l279
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l279
				}
				if !_rules[rulesp]() {
					goto l279
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l279
				}
				if !_rules[ruleAction10]() {
					goto l279
				}
				depth--
				add(ruleUpdateSourceStmt, position280)
			}
			return true
		l279:
			position, tokenIndex, depth = position279, tokenIndex279, depth279
			return false
		},
		/* 17 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action11)> */
		func() bool {
			position305, tokenIndex305, depth305 := position, tokenIndex, depth
			{
				position306 := position
				depth++
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l308
					}
					position++
					goto l307
				l308:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
					if buffer[position] != rune('U') {
						goto l305
					}
					position++
				}
			l307:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('P') {
						goto l305
					}
					position++
				}
			l309:
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l312
					}
					position++
					goto l311
				l312:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
					if buffer[position] != rune('D') {
						goto l305
					}
					position++
				}
			l311:
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('A') {
						goto l305
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('T') {
						goto l305
					}
					position++
				}
			l315:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('E') {
						goto l305
					}
					position++
				}
			l317:
				if !_rules[rulesp]() {
					goto l305
				}
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l320
					}
					position++
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if buffer[position] != rune('S') {
						goto l305
					}
					position++
				}
			l319:
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
						goto l305
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
						goto l305
					}
					position++
				}
			l323:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l326
					}
					position++
					goto l325
				l326:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
					if buffer[position] != rune('K') {
						goto l305
					}
					position++
				}
			l325:
				if !_rules[rulesp]() {
					goto l305
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l305
				}
				if !_rules[rulesp]() {
					goto l305
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l305
				}
				if !_rules[ruleAction11]() {
					goto l305
				}
				depth--
				add(ruleUpdateSinkStmt, position306)
			}
			return true
		l305:
			position, tokenIndex, depth = position305, tokenIndex305, depth305
			return false
		},
		/* 18 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action12)> */
		func() bool {
			position327, tokenIndex327, depth327 := position, tokenIndex, depth
			{
				position328 := position
				depth++
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l330
					}
					position++
					goto l329
				l330:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
					if buffer[position] != rune('I') {
						goto l327
					}
					position++
				}
			l329:
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('N') {
						goto l327
					}
					position++
				}
			l331:
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('S') {
						goto l327
					}
					position++
				}
			l333:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l336
					}
					position++
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if buffer[position] != rune('E') {
						goto l327
					}
					position++
				}
			l335:
				{
					position337, tokenIndex337, depth337 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l338
					}
					position++
					goto l337
				l338:
					position, tokenIndex, depth = position337, tokenIndex337, depth337
					if buffer[position] != rune('R') {
						goto l327
					}
					position++
				}
			l337:
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l340
					}
					position++
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if buffer[position] != rune('T') {
						goto l327
					}
					position++
				}
			l339:
				if !_rules[rulesp]() {
					goto l327
				}
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l342
					}
					position++
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					if buffer[position] != rune('I') {
						goto l327
					}
					position++
				}
			l341:
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('N') {
						goto l327
					}
					position++
				}
			l343:
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l346
					}
					position++
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if buffer[position] != rune('T') {
						goto l327
					}
					position++
				}
			l345:
				{
					position347, tokenIndex347, depth347 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l348
					}
					position++
					goto l347
				l348:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if buffer[position] != rune('O') {
						goto l327
					}
					position++
				}
			l347:
				if !_rules[rulesp]() {
					goto l327
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l327
				}
				if !_rules[rulesp]() {
					goto l327
				}
				if !_rules[ruleSelectStmt]() {
					goto l327
				}
				if !_rules[ruleAction12]() {
					goto l327
				}
				depth--
				add(ruleInsertIntoSelectStmt, position328)
			}
			return true
		l327:
			position, tokenIndex, depth = position327, tokenIndex327, depth327
			return false
		},
		/* 19 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action13)> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l352
					}
					position++
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					if buffer[position] != rune('I') {
						goto l349
					}
					position++
				}
			l351:
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l354
					}
					position++
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if buffer[position] != rune('N') {
						goto l349
					}
					position++
				}
			l353:
				{
					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l356
					}
					position++
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune('S') {
						goto l349
					}
					position++
				}
			l355:
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l358
					}
					position++
					goto l357
				l358:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
					if buffer[position] != rune('E') {
						goto l349
					}
					position++
				}
			l357:
				{
					position359, tokenIndex359, depth359 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l360
					}
					position++
					goto l359
				l360:
					position, tokenIndex, depth = position359, tokenIndex359, depth359
					if buffer[position] != rune('R') {
						goto l349
					}
					position++
				}
			l359:
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l362
					}
					position++
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if buffer[position] != rune('T') {
						goto l349
					}
					position++
				}
			l361:
				if !_rules[rulesp]() {
					goto l349
				}
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l364
					}
					position++
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					if buffer[position] != rune('I') {
						goto l349
					}
					position++
				}
			l363:
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('N') {
						goto l349
					}
					position++
				}
			l365:
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					if buffer[position] != rune('T') {
						goto l349
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('O') {
						goto l349
					}
					position++
				}
			l369:
				if !_rules[rulesp]() {
					goto l349
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l349
				}
				if !_rules[rulesp]() {
					goto l349
				}
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l372
					}
					position++
					goto l371
				l372:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
					if buffer[position] != rune('F') {
						goto l349
					}
					position++
				}
			l371:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('R') {
						goto l349
					}
					position++
				}
			l373:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l376
					}
					position++
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					if buffer[position] != rune('O') {
						goto l349
					}
					position++
				}
			l375:
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l378
					}
					position++
					goto l377
				l378:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
					if buffer[position] != rune('M') {
						goto l349
					}
					position++
				}
			l377:
				if !_rules[rulesp]() {
					goto l349
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l349
				}
				if !_rules[ruleAction13]() {
					goto l349
				}
				depth--
				add(ruleInsertIntoFromStmt, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 20 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position379, tokenIndex379, depth379 := position, tokenIndex, depth
			{
				position380 := position
				depth++
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('P') {
						goto l379
					}
					position++
				}
			l381:
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('A') {
						goto l379
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('U') {
						goto l379
					}
					position++
				}
			l385:
				{
					position387, tokenIndex387, depth387 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l388
					}
					position++
					goto l387
				l388:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
					if buffer[position] != rune('S') {
						goto l379
					}
					position++
				}
			l387:
				{
					position389, tokenIndex389, depth389 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l390
					}
					position++
					goto l389
				l390:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
					if buffer[position] != rune('E') {
						goto l379
					}
					position++
				}
			l389:
				if !_rules[rulesp]() {
					goto l379
				}
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
					if buffer[position] != rune('S') {
						goto l379
					}
					position++
				}
			l391:
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l394
					}
					position++
					goto l393
				l394:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
					if buffer[position] != rune('O') {
						goto l379
					}
					position++
				}
			l393:
				{
					position395, tokenIndex395, depth395 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l396
					}
					position++
					goto l395
				l396:
					position, tokenIndex, depth = position395, tokenIndex395, depth395
					if buffer[position] != rune('U') {
						goto l379
					}
					position++
				}
			l395:
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
						goto l379
					}
					position++
				}
			l397:
				{
					position399, tokenIndex399, depth399 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l400
					}
					position++
					goto l399
				l400:
					position, tokenIndex, depth = position399, tokenIndex399, depth399
					if buffer[position] != rune('C') {
						goto l379
					}
					position++
				}
			l399:
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l402
					}
					position++
					goto l401
				l402:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if buffer[position] != rune('E') {
						goto l379
					}
					position++
				}
			l401:
				if !_rules[rulesp]() {
					goto l379
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l379
				}
				if !_rules[ruleAction14]() {
					goto l379
				}
				depth--
				add(rulePauseSourceStmt, position380)
			}
			return true
		l379:
			position, tokenIndex, depth = position379, tokenIndex379, depth379
			return false
		},
		/* 21 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position403, tokenIndex403, depth403 := position, tokenIndex, depth
			{
				position404 := position
				depth++
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l406
					}
					position++
					goto l405
				l406:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if buffer[position] != rune('R') {
						goto l403
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
						goto l403
					}
					position++
				}
			l407:
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
						goto l403
					}
					position++
				}
			l409:
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l412
					}
					position++
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if buffer[position] != rune('U') {
						goto l403
					}
					position++
				}
			l411:
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l414
					}
					position++
					goto l413
				l414:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
					if buffer[position] != rune('M') {
						goto l403
					}
					position++
				}
			l413:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l416
					}
					position++
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if buffer[position] != rune('E') {
						goto l403
					}
					position++
				}
			l415:
				if !_rules[rulesp]() {
					goto l403
				}
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l418
					}
					position++
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if buffer[position] != rune('S') {
						goto l403
					}
					position++
				}
			l417:
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l420
					}
					position++
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if buffer[position] != rune('O') {
						goto l403
					}
					position++
				}
			l419:
				{
					position421, tokenIndex421, depth421 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l422
					}
					position++
					goto l421
				l422:
					position, tokenIndex, depth = position421, tokenIndex421, depth421
					if buffer[position] != rune('U') {
						goto l403
					}
					position++
				}
			l421:
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
						goto l403
					}
					position++
				}
			l423:
				{
					position425, tokenIndex425, depth425 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l426
					}
					position++
					goto l425
				l426:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if buffer[position] != rune('C') {
						goto l403
					}
					position++
				}
			l425:
				{
					position427, tokenIndex427, depth427 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l428
					}
					position++
					goto l427
				l428:
					position, tokenIndex, depth = position427, tokenIndex427, depth427
					if buffer[position] != rune('E') {
						goto l403
					}
					position++
				}
			l427:
				if !_rules[rulesp]() {
					goto l403
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l403
				}
				if !_rules[ruleAction15]() {
					goto l403
				}
				depth--
				add(ruleResumeSourceStmt, position404)
			}
			return true
		l403:
			position, tokenIndex, depth = position403, tokenIndex403, depth403
			return false
		},
		/* 22 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position429, tokenIndex429, depth429 := position, tokenIndex, depth
			{
				position430 := position
				depth++
				{
					position431, tokenIndex431, depth431 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l432
					}
					position++
					goto l431
				l432:
					position, tokenIndex, depth = position431, tokenIndex431, depth431
					if buffer[position] != rune('R') {
						goto l429
					}
					position++
				}
			l431:
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l434
					}
					position++
					goto l433
				l434:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if buffer[position] != rune('E') {
						goto l429
					}
					position++
				}
			l433:
				{
					position435, tokenIndex435, depth435 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l436
					}
					position++
					goto l435
				l436:
					position, tokenIndex, depth = position435, tokenIndex435, depth435
					if buffer[position] != rune('W') {
						goto l429
					}
					position++
				}
			l435:
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l438
					}
					position++
					goto l437
				l438:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
					if buffer[position] != rune('I') {
						goto l429
					}
					position++
				}
			l437:
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l440
					}
					position++
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if buffer[position] != rune('N') {
						goto l429
					}
					position++
				}
			l439:
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l442
					}
					position++
					goto l441
				l442:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
					if buffer[position] != rune('D') {
						goto l429
					}
					position++
				}
			l441:
				if !_rules[rulesp]() {
					goto l429
				}
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l444
					}
					position++
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if buffer[position] != rune('S') {
						goto l429
					}
					position++
				}
			l443:
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l446
					}
					position++
					goto l445
				l446:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
					if buffer[position] != rune('O') {
						goto l429
					}
					position++
				}
			l445:
				{
					position447, tokenIndex447, depth447 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l448
					}
					position++
					goto l447
				l448:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if buffer[position] != rune('U') {
						goto l429
					}
					position++
				}
			l447:
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l450
					}
					position++
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					if buffer[position] != rune('R') {
						goto l429
					}
					position++
				}
			l449:
				{
					position451, tokenIndex451, depth451 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l452
					}
					position++
					goto l451
				l452:
					position, tokenIndex, depth = position451, tokenIndex451, depth451
					if buffer[position] != rune('C') {
						goto l429
					}
					position++
				}
			l451:
				{
					position453, tokenIndex453, depth453 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l454
					}
					position++
					goto l453
				l454:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
					if buffer[position] != rune('E') {
						goto l429
					}
					position++
				}
			l453:
				if !_rules[rulesp]() {
					goto l429
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l429
				}
				if !_rules[ruleAction16]() {
					goto l429
				}
				depth--
				add(ruleRewindSourceStmt, position430)
			}
			return true
		l429:
			position, tokenIndex, depth = position429, tokenIndex429, depth429
			return false
		},
		/* 23 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action17)> */
		func() bool {
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				{
					position457, tokenIndex457, depth457 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l458
					}
					position++
					goto l457
				l458:
					position, tokenIndex, depth = position457, tokenIndex457, depth457
					if buffer[position] != rune('D') {
						goto l455
					}
					position++
				}
			l457:
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l460
					}
					position++
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if buffer[position] != rune('R') {
						goto l455
					}
					position++
				}
			l459:
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l462
					}
					position++
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('O') {
						goto l455
					}
					position++
				}
			l461:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('P') {
						goto l455
					}
					position++
				}
			l463:
				if !_rules[rulesp]() {
					goto l455
				}
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l466
					}
					position++
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					if buffer[position] != rune('S') {
						goto l455
					}
					position++
				}
			l465:
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l468
					}
					position++
					goto l467
				l468:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
					if buffer[position] != rune('O') {
						goto l455
					}
					position++
				}
			l467:
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l470
					}
					position++
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if buffer[position] != rune('U') {
						goto l455
					}
					position++
				}
			l469:
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l472
					}
					position++
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if buffer[position] != rune('R') {
						goto l455
					}
					position++
				}
			l471:
				{
					position473, tokenIndex473, depth473 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l474
					}
					position++
					goto l473
				l474:
					position, tokenIndex, depth = position473, tokenIndex473, depth473
					if buffer[position] != rune('C') {
						goto l455
					}
					position++
				}
			l473:
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l476
					}
					position++
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if buffer[position] != rune('E') {
						goto l455
					}
					position++
				}
			l475:
				if !_rules[rulesp]() {
					goto l455
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l455
				}
				if !_rules[ruleAction17]() {
					goto l455
				}
				depth--
				add(ruleDropSourceStmt, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
			return false
		},
		/* 24 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action18)> */
		func() bool {
			position477, tokenIndex477, depth477 := position, tokenIndex, depth
			{
				position478 := position
				depth++
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l480
					}
					position++
					goto l479
				l480:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if buffer[position] != rune('D') {
						goto l477
					}
					position++
				}
			l479:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if buffer[position] != rune('R') {
						goto l477
					}
					position++
				}
			l481:
				{
					position483, tokenIndex483, depth483 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l484
					}
					position++
					goto l483
				l484:
					position, tokenIndex, depth = position483, tokenIndex483, depth483
					if buffer[position] != rune('O') {
						goto l477
					}
					position++
				}
			l483:
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l486
					}
					position++
					goto l485
				l486:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					if buffer[position] != rune('P') {
						goto l477
					}
					position++
				}
			l485:
				if !_rules[rulesp]() {
					goto l477
				}
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l488
					}
					position++
					goto l487
				l488:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
					if buffer[position] != rune('S') {
						goto l477
					}
					position++
				}
			l487:
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l490
					}
					position++
					goto l489
				l490:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
					if buffer[position] != rune('T') {
						goto l477
					}
					position++
				}
			l489:
				{
					position491, tokenIndex491, depth491 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l492
					}
					position++
					goto l491
				l492:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if buffer[position] != rune('R') {
						goto l477
					}
					position++
				}
			l491:
				{
					position493, tokenIndex493, depth493 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l494
					}
					position++
					goto l493
				l494:
					position, tokenIndex, depth = position493, tokenIndex493, depth493
					if buffer[position] != rune('E') {
						goto l477
					}
					position++
				}
			l493:
				{
					position495, tokenIndex495, depth495 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l496
					}
					position++
					goto l495
				l496:
					position, tokenIndex, depth = position495, tokenIndex495, depth495
					if buffer[position] != rune('A') {
						goto l477
					}
					position++
				}
			l495:
				{
					position497, tokenIndex497, depth497 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l498
					}
					position++
					goto l497
				l498:
					position, tokenIndex, depth = position497, tokenIndex497, depth497
					if buffer[position] != rune('M') {
						goto l477
					}
					position++
				}
			l497:
				if !_rules[rulesp]() {
					goto l477
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l477
				}
				if !_rules[ruleAction18]() {
					goto l477
				}
				depth--
				add(ruleDropStreamStmt, position478)
			}
			return true
		l477:
			position, tokenIndex, depth = position477, tokenIndex477, depth477
			return false
		},
		/* 25 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action19)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				{
					position501, tokenIndex501, depth501 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l502
					}
					position++
					goto l501
				l502:
					position, tokenIndex, depth = position501, tokenIndex501, depth501
					if buffer[position] != rune('D') {
						goto l499
					}
					position++
				}
			l501:
				{
					position503, tokenIndex503, depth503 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l504
					}
					position++
					goto l503
				l504:
					position, tokenIndex, depth = position503, tokenIndex503, depth503
					if buffer[position] != rune('R') {
						goto l499
					}
					position++
				}
			l503:
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l506
					}
					position++
					goto l505
				l506:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
					if buffer[position] != rune('O') {
						goto l499
					}
					position++
				}
			l505:
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('P') {
						goto l499
					}
					position++
				}
			l507:
				if !_rules[rulesp]() {
					goto l499
				}
				{
					position509, tokenIndex509, depth509 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l510
					}
					position++
					goto l509
				l510:
					position, tokenIndex, depth = position509, tokenIndex509, depth509
					if buffer[position] != rune('S') {
						goto l499
					}
					position++
				}
			l509:
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if buffer[position] != rune('I') {
						goto l499
					}
					position++
				}
			l511:
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l514
					}
					position++
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if buffer[position] != rune('N') {
						goto l499
					}
					position++
				}
			l513:
				{
					position515, tokenIndex515, depth515 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l516
					}
					position++
					goto l515
				l516:
					position, tokenIndex, depth = position515, tokenIndex515, depth515
					if buffer[position] != rune('K') {
						goto l499
					}
					position++
				}
			l515:
				if !_rules[rulesp]() {
					goto l499
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l499
				}
				if !_rules[ruleAction19]() {
					goto l499
				}
				depth--
				add(ruleDropSinkStmt, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 26 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action20)> */
		func() bool {
			position517, tokenIndex517, depth517 := position, tokenIndex, depth
			{
				position518 := position
				depth++
				{
					position519, tokenIndex519, depth519 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l520
					}
					position++
					goto l519
				l520:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if buffer[position] != rune('D') {
						goto l517
					}
					position++
				}
			l519:
				{
					position521, tokenIndex521, depth521 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l522
					}
					position++
					goto l521
				l522:
					position, tokenIndex, depth = position521, tokenIndex521, depth521
					if buffer[position] != rune('R') {
						goto l517
					}
					position++
				}
			l521:
				{
					position523, tokenIndex523, depth523 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l524
					}
					position++
					goto l523
				l524:
					position, tokenIndex, depth = position523, tokenIndex523, depth523
					if buffer[position] != rune('O') {
						goto l517
					}
					position++
				}
			l523:
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l526
					}
					position++
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if buffer[position] != rune('P') {
						goto l517
					}
					position++
				}
			l525:
				if !_rules[rulesp]() {
					goto l517
				}
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l528
					}
					position++
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if buffer[position] != rune('S') {
						goto l517
					}
					position++
				}
			l527:
				{
					position529, tokenIndex529, depth529 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l530
					}
					position++
					goto l529
				l530:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if buffer[position] != rune('T') {
						goto l517
					}
					position++
				}
			l529:
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l532
					}
					position++
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if buffer[position] != rune('A') {
						goto l517
					}
					position++
				}
			l531:
				{
					position533, tokenIndex533, depth533 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l534
					}
					position++
					goto l533
				l534:
					position, tokenIndex, depth = position533, tokenIndex533, depth533
					if buffer[position] != rune('T') {
						goto l517
					}
					position++
				}
			l533:
				{
					position535, tokenIndex535, depth535 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l536
					}
					position++
					goto l535
				l536:
					position, tokenIndex, depth = position535, tokenIndex535, depth535
					if buffer[position] != rune('E') {
						goto l517
					}
					position++
				}
			l535:
				if !_rules[rulesp]() {
					goto l517
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l517
				}
				if !_rules[ruleAction20]() {
					goto l517
				}
				depth--
				add(ruleDropStateStmt, position518)
			}
			return true
		l517:
			position, tokenIndex, depth = position517, tokenIndex517, depth517
			return false
		},
		/* 27 LoadStateStmt <- <(('l' / 'L') ('o' / 'O') ('a' / 'A') ('d' / 'D') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SetOptSpecs Action21)> */
		func() bool {
			position537, tokenIndex537, depth537 := position, tokenIndex, depth
			{
				position538 := position
				depth++
				{
					position539, tokenIndex539, depth539 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l540
					}
					position++
					goto l539
				l540:
					position, tokenIndex, depth = position539, tokenIndex539, depth539
					if buffer[position] != rune('L') {
						goto l537
					}
					position++
				}
			l539:
				{
					position541, tokenIndex541, depth541 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l542
					}
					position++
					goto l541
				l542:
					position, tokenIndex, depth = position541, tokenIndex541, depth541
					if buffer[position] != rune('O') {
						goto l537
					}
					position++
				}
			l541:
				{
					position543, tokenIndex543, depth543 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l544
					}
					position++
					goto l543
				l544:
					position, tokenIndex, depth = position543, tokenIndex543, depth543
					if buffer[position] != rune('A') {
						goto l537
					}
					position++
				}
			l543:
				{
					position545, tokenIndex545, depth545 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l546
					}
					position++
					goto l545
				l546:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if buffer[position] != rune('D') {
						goto l537
					}
					position++
				}
			l545:
				if !_rules[rulesp]() {
					goto l537
				}
				{
					position547, tokenIndex547, depth547 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l548
					}
					position++
					goto l547
				l548:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
					if buffer[position] != rune('S') {
						goto l537
					}
					position++
				}
			l547:
				{
					position549, tokenIndex549, depth549 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l550
					}
					position++
					goto l549
				l550:
					position, tokenIndex, depth = position549, tokenIndex549, depth549
					if buffer[position] != rune('T') {
						goto l537
					}
					position++
				}
			l549:
				{
					position551, tokenIndex551, depth551 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l552
					}
					position++
					goto l551
				l552:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if buffer[position] != rune('A') {
						goto l537
					}
					position++
				}
			l551:
				{
					position553, tokenIndex553, depth553 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l554
					}
					position++
					goto l553
				l554:
					position, tokenIndex, depth = position553, tokenIndex553, depth553
					if buffer[position] != rune('T') {
						goto l537
					}
					position++
				}
			l553:
				{
					position555, tokenIndex555, depth555 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l556
					}
					position++
					goto l555
				l556:
					position, tokenIndex, depth = position555, tokenIndex555, depth555
					if buffer[position] != rune('E') {
						goto l537
					}
					position++
				}
			l555:
				if !_rules[rulesp]() {
					goto l537
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l537
				}
				if !_rules[rulesp]() {
					goto l537
				}
				{
					position557, tokenIndex557, depth557 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l558
					}
					position++
					goto l557
				l558:
					position, tokenIndex, depth = position557, tokenIndex557, depth557
					if buffer[position] != rune('T') {
						goto l537
					}
					position++
				}
			l557:
				{
					position559, tokenIndex559, depth559 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l560
					}
					position++
					goto l559
				l560:
					position, tokenIndex, depth = position559, tokenIndex559, depth559
					if buffer[position] != rune('Y') {
						goto l537
					}
					position++
				}
			l559:
				{
					position561, tokenIndex561, depth561 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l562
					}
					position++
					goto l561
				l562:
					position, tokenIndex, depth = position561, tokenIndex561, depth561
					if buffer[position] != rune('P') {
						goto l537
					}
					position++
				}
			l561:
				{
					position563, tokenIndex563, depth563 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l564
					}
					position++
					goto l563
				l564:
					position, tokenIndex, depth = position563, tokenIndex563, depth563
					if buffer[position] != rune('E') {
						goto l537
					}
					position++
				}
			l563:
				if !_rules[rulesp]() {
					goto l537
				}
				if !_rules[ruleSourceSinkType]() {
					goto l537
				}
				if !_rules[rulesp]() {
					goto l537
				}
				if !_rules[ruleSetOptSpecs]() {
					goto l537
				}
				if !_rules[ruleAction21]() {
					goto l537
				}
				depth--
				add(ruleLoadStateStmt, position538)
			}
			return true
		l537:
			position, tokenIndex, depth = position537, tokenIndex537, depth537
			return false
		},
		/* 28 LoadStateOrCreateStmt <- <(LoadStateStmt sp (('o' / 'O') ('r' / 'R')) sp (('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp (('i' / 'I') ('f' / 'F')) sp (('n' / 'N') ('o' / 'O') ('t' / 'T')) sp (('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S')) sp SourceSinkSpecs Action22)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				if !_rules[ruleLoadStateStmt]() {
					goto l565
				}
				if !_rules[rulesp]() {
					goto l565
				}
				{
					position567, tokenIndex567, depth567 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l568
					}
					position++
					goto l567
				l568:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if buffer[position] != rune('O') {
						goto l565
					}
					position++
				}
			l567:
				{
					position569, tokenIndex569, depth569 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l570
					}
					position++
					goto l569
				l570:
					position, tokenIndex, depth = position569, tokenIndex569, depth569
					if buffer[position] != rune('R') {
						goto l565
					}
					position++
				}
			l569:
				if !_rules[rulesp]() {
					goto l565
				}
				{
					position571, tokenIndex571, depth571 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l572
					}
					position++
					goto l571
				l572:
					position, tokenIndex, depth = position571, tokenIndex571, depth571
					if buffer[position] != rune('C') {
						goto l565
					}
					position++
				}
			l571:
				{
					position573, tokenIndex573, depth573 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l574
					}
					position++
					goto l573
				l574:
					position, tokenIndex, depth = position573, tokenIndex573, depth573
					if buffer[position] != rune('R') {
						goto l565
					}
					position++
				}
			l573:
				{
					position575, tokenIndex575, depth575 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l576
					}
					position++
					goto l575
				l576:
					position, tokenIndex, depth = position575, tokenIndex575, depth575
					if buffer[position] != rune('E') {
						goto l565
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
						goto l565
					}
					position++
				}
			l577:
				{
					position579, tokenIndex579, depth579 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l580
					}
					position++
					goto l579
				l580:
					position, tokenIndex, depth = position579, tokenIndex579, depth579
					if buffer[position] != rune('T') {
						goto l565
					}
					position++
				}
			l579:
				{
					position581, tokenIndex581, depth581 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l582
					}
					position++
					goto l581
				l582:
					position, tokenIndex, depth = position581, tokenIndex581, depth581
					if buffer[position] != rune('E') {
						goto l565
					}
					position++
				}
			l581:
				if !_rules[rulesp]() {
					goto l565
				}
				{
					position583, tokenIndex583, depth583 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l584
					}
					position++
					goto l583
				l584:
					position, tokenIndex, depth = position583, tokenIndex583, depth583
					if buffer[position] != rune('I') {
						goto l565
					}
					position++
				}
			l583:
				{
					position585, tokenIndex585, depth585 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l586
					}
					position++
					goto l585
				l586:
					position, tokenIndex, depth = position585, tokenIndex585, depth585
					if buffer[position] != rune('F') {
						goto l565
					}
					position++
				}
			l585:
				if !_rules[rulesp]() {
					goto l565
				}
				{
					position587, tokenIndex587, depth587 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l588
					}
					position++
					goto l587
				l588:
					position, tokenIndex, depth = position587, tokenIndex587, depth587
					if buffer[position] != rune('N') {
						goto l565
					}
					position++
				}
			l587:
				{
					position589, tokenIndex589, depth589 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l590
					}
					position++
					goto l589
				l590:
					position, tokenIndex, depth = position589, tokenIndex589, depth589
					if buffer[position] != rune('O') {
						goto l565
					}
					position++
				}
			l589:
				{
					position591, tokenIndex591, depth591 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l592
					}
					position++
					goto l591
				l592:
					position, tokenIndex, depth = position591, tokenIndex591, depth591
					if buffer[position] != rune('T') {
						goto l565
					}
					position++
				}
			l591:
				if !_rules[rulesp]() {
					goto l565
				}
				{
					position593, tokenIndex593, depth593 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l594
					}
					position++
					goto l593
				l594:
					position, tokenIndex, depth = position593, tokenIndex593, depth593
					if buffer[position] != rune('E') {
						goto l565
					}
					position++
				}
			l593:
				{
					position595, tokenIndex595, depth595 := position, tokenIndex, depth
					if buffer[position] != rune('x') {
						goto l596
					}
					position++
					goto l595
				l596:
					position, tokenIndex, depth = position595, tokenIndex595, depth595
					if buffer[position] != rune('X') {
						goto l565
					}
					position++
				}
			l595:
				{
					position597, tokenIndex597, depth597 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l598
					}
					position++
					goto l597
				l598:
					position, tokenIndex, depth = position597, tokenIndex597, depth597
					if buffer[position] != rune('I') {
						goto l565
					}
					position++
				}
			l597:
				{
					position599, tokenIndex599, depth599 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l600
					}
					position++
					goto l599
				l600:
					position, tokenIndex, depth = position599, tokenIndex599, depth599
					if buffer[position] != rune('S') {
						goto l565
					}
					position++
				}
			l599:
				{
					position601, tokenIndex601, depth601 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l602
					}
					position++
					goto l601
				l602:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
					if buffer[position] != rune('T') {
						goto l565
					}
					position++
				}
			l601:
				{
					position603, tokenIndex603, depth603 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l604
					}
					position++
					goto l603
				l604:
					position, tokenIndex, depth = position603, tokenIndex603, depth603
					if buffer[position] != rune('S') {
						goto l565
					}
					position++
				}
			l603:
				if !_rules[rulesp]() {
					goto l565
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l565
				}
				if !_rules[ruleAction22]() {
					goto l565
				}
				depth--
				add(ruleLoadStateOrCreateStmt, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 29 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) sp EmitterOptions Action23)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				{
					position607, tokenIndex607, depth607 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l608
					}
					goto l607
				l608:
					position, tokenIndex, depth = position607, tokenIndex607, depth607
					if !_rules[ruleDSTREAM]() {
						goto l609
					}
					goto l607
				l609:
					position, tokenIndex, depth = position607, tokenIndex607, depth607
					if !_rules[ruleRSTREAM]() {
						goto l605
					}
				}
			l607:
				if !_rules[rulesp]() {
					goto l605
				}
				if !_rules[ruleEmitterOptions]() {
					goto l605
				}
				if !_rules[ruleAction23]() {
					goto l605
				}
				depth--
				add(ruleEmitter, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 30 EmitterOptions <- <(<('[' sp EmitterLimit sp ']')?> Action24)> */
		func() bool {
			position610, tokenIndex610, depth610 := position, tokenIndex, depth
			{
				position611 := position
				depth++
				{
					position612 := position
					depth++
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('[') {
							goto l613
						}
						position++
						if !_rules[rulesp]() {
							goto l613
						}
						if !_rules[ruleEmitterLimit]() {
							goto l613
						}
						if !_rules[rulesp]() {
							goto l613
						}
						if buffer[position] != rune(']') {
							goto l613
						}
						position++
						goto l614
					l613:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
					}
				l614:
					depth--
					add(rulePegText, position612)
				}
				if !_rules[ruleAction24]() {
					goto l610
				}
				depth--
				add(ruleEmitterOptions, position611)
			}
			return true
		l610:
			position, tokenIndex, depth = position610, tokenIndex610, depth610
			return false
		},
		/* 31 EmitterLimit <- <('L' 'I' 'M' 'I' 'T' sp NumericLiteral Action25)> */
		func() bool {
			position615, tokenIndex615, depth615 := position, tokenIndex, depth
			{
				position616 := position
				depth++
				if buffer[position] != rune('L') {
					goto l615
				}
				position++
				if buffer[position] != rune('I') {
					goto l615
				}
				position++
				if buffer[position] != rune('M') {
					goto l615
				}
				position++
				if buffer[position] != rune('I') {
					goto l615
				}
				position++
				if buffer[position] != rune('T') {
					goto l615
				}
				position++
				if !_rules[rulesp]() {
					goto l615
				}
				if !_rules[ruleNumericLiteral]() {
					goto l615
				}
				if !_rules[ruleAction25]() {
					goto l615
				}
				depth--
				add(ruleEmitterLimit, position616)
			}
			return true
		l615:
			position, tokenIndex, depth = position615, tokenIndex615, depth615
			return false
		},
		/* 32 Projections <- <(<(Projection sp (',' sp Projection)*)> Action26)> */
		func() bool {
			position617, tokenIndex617, depth617 := position, tokenIndex, depth
			{
				position618 := position
				depth++
				{
					position619 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l617
					}
					if !_rules[rulesp]() {
						goto l617
					}
				l620:
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l621
						}
						position++
						if !_rules[rulesp]() {
							goto l621
						}
						if !_rules[ruleProjection]() {
							goto l621
						}
						goto l620
					l621:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
					}
					depth--
					add(rulePegText, position619)
				}
				if !_rules[ruleAction26]() {
					goto l617
				}
				depth--
				add(ruleProjections, position618)
			}
			return true
		l617:
			position, tokenIndex, depth = position617, tokenIndex617, depth617
			return false
		},
		/* 33 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position622, tokenIndex622, depth622 := position, tokenIndex, depth
			{
				position623 := position
				depth++
				{
					position624, tokenIndex624, depth624 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l625
					}
					goto l624
				l625:
					position, tokenIndex, depth = position624, tokenIndex624, depth624
					if !_rules[ruleExpression]() {
						goto l626
					}
					goto l624
				l626:
					position, tokenIndex, depth = position624, tokenIndex624, depth624
					if !_rules[ruleWildcard]() {
						goto l622
					}
				}
			l624:
				depth--
				add(ruleProjection, position623)
			}
			return true
		l622:
			position, tokenIndex, depth = position622, tokenIndex622, depth622
			return false
		},
		/* 34 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action27)> */
		func() bool {
			position627, tokenIndex627, depth627 := position, tokenIndex, depth
			{
				position628 := position
				depth++
				{
					position629, tokenIndex629, depth629 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l630
					}
					goto l629
				l630:
					position, tokenIndex, depth = position629, tokenIndex629, depth629
					if !_rules[ruleWildcard]() {
						goto l627
					}
				}
			l629:
				if !_rules[rulesp]() {
					goto l627
				}
				{
					position631, tokenIndex631, depth631 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l632
					}
					position++
					goto l631
				l632:
					position, tokenIndex, depth = position631, tokenIndex631, depth631
					if buffer[position] != rune('A') {
						goto l627
					}
					position++
				}
			l631:
				{
					position633, tokenIndex633, depth633 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l634
					}
					position++
					goto l633
				l634:
					position, tokenIndex, depth = position633, tokenIndex633, depth633
					if buffer[position] != rune('S') {
						goto l627
					}
					position++
				}
			l633:
				if !_rules[rulesp]() {
					goto l627
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l627
				}
				if !_rules[ruleAction27]() {
					goto l627
				}
				depth--
				add(ruleAliasExpression, position628)
			}
			return true
		l627:
			position, tokenIndex, depth = position627, tokenIndex627, depth627
			return false
		},
		/* 35 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action28)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				{
					position637 := position
					depth++
					{
						position638, tokenIndex638, depth638 := position, tokenIndex, depth
						{
							position640, tokenIndex640, depth640 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l641
							}
							position++
							goto l640
						l641:
							position, tokenIndex, depth = position640, tokenIndex640, depth640
							if buffer[position] != rune('F') {
								goto l638
							}
							position++
						}
					l640:
						{
							position642, tokenIndex642, depth642 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l643
							}
							position++
							goto l642
						l643:
							position, tokenIndex, depth = position642, tokenIndex642, depth642
							if buffer[position] != rune('R') {
								goto l638
							}
							position++
						}
					l642:
						{
							position644, tokenIndex644, depth644 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l645
							}
							position++
							goto l644
						l645:
							position, tokenIndex, depth = position644, tokenIndex644, depth644
							if buffer[position] != rune('O') {
								goto l638
							}
							position++
						}
					l644:
						{
							position646, tokenIndex646, depth646 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l647
							}
							position++
							goto l646
						l647:
							position, tokenIndex, depth = position646, tokenIndex646, depth646
							if buffer[position] != rune('M') {
								goto l638
							}
							position++
						}
					l646:
						if !_rules[rulesp]() {
							goto l638
						}
						if !_rules[ruleRelations]() {
							goto l638
						}
						if !_rules[rulesp]() {
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
				if !_rules[ruleAction28]() {
					goto l635
				}
				depth--
				add(ruleWindowedFrom, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 36 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position648, tokenIndex648, depth648 := position, tokenIndex, depth
			{
				position649 := position
				depth++
				{
					position650, tokenIndex650, depth650 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l651
					}
					goto l650
				l651:
					position, tokenIndex, depth = position650, tokenIndex650, depth650
					if !_rules[ruleTuplesInterval]() {
						goto l648
					}
				}
			l650:
				depth--
				add(ruleInterval, position649)
			}
			return true
		l648:
			position, tokenIndex, depth = position648, tokenIndex648, depth648
			return false
		},
		/* 37 TimeInterval <- <(NumericLiteral sp (SECONDS / MILLISECONDS) Action29)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l652
				}
				if !_rules[rulesp]() {
					goto l652
				}
				{
					position654, tokenIndex654, depth654 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l655
					}
					goto l654
				l655:
					position, tokenIndex, depth = position654, tokenIndex654, depth654
					if !_rules[ruleMILLISECONDS]() {
						goto l652
					}
				}
			l654:
				if !_rules[ruleAction29]() {
					goto l652
				}
				depth--
				add(ruleTimeInterval, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 38 TuplesInterval <- <(NumericLiteral sp TUPLES Action30)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l656
				}
				if !_rules[rulesp]() {
					goto l656
				}
				if !_rules[ruleTUPLES]() {
					goto l656
				}
				if !_rules[ruleAction30]() {
					goto l656
				}
				depth--
				add(ruleTuplesInterval, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 39 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l658
				}
				if !_rules[rulesp]() {
					goto l658
				}
			l660:
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l661
					}
					position++
					if !_rules[rulesp]() {
						goto l661
					}
					if !_rules[ruleRelationLike]() {
						goto l661
					}
					goto l660
				l661:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
				}
				depth--
				add(ruleRelations, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 40 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action31)> */
		func() bool {
			position662, tokenIndex662, depth662 := position, tokenIndex, depth
			{
				position663 := position
				depth++
				{
					position664 := position
					depth++
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						{
							position667, tokenIndex667, depth667 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l668
							}
							position++
							goto l667
						l668:
							position, tokenIndex, depth = position667, tokenIndex667, depth667
							if buffer[position] != rune('W') {
								goto l665
							}
							position++
						}
					l667:
						{
							position669, tokenIndex669, depth669 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l670
							}
							position++
							goto l669
						l670:
							position, tokenIndex, depth = position669, tokenIndex669, depth669
							if buffer[position] != rune('H') {
								goto l665
							}
							position++
						}
					l669:
						{
							position671, tokenIndex671, depth671 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l672
							}
							position++
							goto l671
						l672:
							position, tokenIndex, depth = position671, tokenIndex671, depth671
							if buffer[position] != rune('E') {
								goto l665
							}
							position++
						}
					l671:
						{
							position673, tokenIndex673, depth673 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l674
							}
							position++
							goto l673
						l674:
							position, tokenIndex, depth = position673, tokenIndex673, depth673
							if buffer[position] != rune('R') {
								goto l665
							}
							position++
						}
					l673:
						{
							position675, tokenIndex675, depth675 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l676
							}
							position++
							goto l675
						l676:
							position, tokenIndex, depth = position675, tokenIndex675, depth675
							if buffer[position] != rune('E') {
								goto l665
							}
							position++
						}
					l675:
						if !_rules[rulesp]() {
							goto l665
						}
						if !_rules[ruleExpression]() {
							goto l665
						}
						goto l666
					l665:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
					}
				l666:
					depth--
					add(rulePegText, position664)
				}
				if !_rules[ruleAction31]() {
					goto l662
				}
				depth--
				add(ruleFilter, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 41 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action32)> */
		func() bool {
			position677, tokenIndex677, depth677 := position, tokenIndex, depth
			{
				position678 := position
				depth++
				{
					position679 := position
					depth++
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						{
							position682, tokenIndex682, depth682 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l683
							}
							position++
							goto l682
						l683:
							position, tokenIndex, depth = position682, tokenIndex682, depth682
							if buffer[position] != rune('G') {
								goto l680
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
								goto l680
							}
							position++
						}
					l684:
						{
							position686, tokenIndex686, depth686 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l687
							}
							position++
							goto l686
						l687:
							position, tokenIndex, depth = position686, tokenIndex686, depth686
							if buffer[position] != rune('O') {
								goto l680
							}
							position++
						}
					l686:
						{
							position688, tokenIndex688, depth688 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l689
							}
							position++
							goto l688
						l689:
							position, tokenIndex, depth = position688, tokenIndex688, depth688
							if buffer[position] != rune('U') {
								goto l680
							}
							position++
						}
					l688:
						{
							position690, tokenIndex690, depth690 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l691
							}
							position++
							goto l690
						l691:
							position, tokenIndex, depth = position690, tokenIndex690, depth690
							if buffer[position] != rune('P') {
								goto l680
							}
							position++
						}
					l690:
						if !_rules[rulesp]() {
							goto l680
						}
						{
							position692, tokenIndex692, depth692 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l693
							}
							position++
							goto l692
						l693:
							position, tokenIndex, depth = position692, tokenIndex692, depth692
							if buffer[position] != rune('B') {
								goto l680
							}
							position++
						}
					l692:
						{
							position694, tokenIndex694, depth694 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l695
							}
							position++
							goto l694
						l695:
							position, tokenIndex, depth = position694, tokenIndex694, depth694
							if buffer[position] != rune('Y') {
								goto l680
							}
							position++
						}
					l694:
						if !_rules[rulesp]() {
							goto l680
						}
						if !_rules[ruleGroupList]() {
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
				if !_rules[ruleAction32]() {
					goto l677
				}
				depth--
				add(ruleGrouping, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 42 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position696, tokenIndex696, depth696 := position, tokenIndex, depth
			{
				position697 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l696
				}
				if !_rules[rulesp]() {
					goto l696
				}
			l698:
				{
					position699, tokenIndex699, depth699 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l699
					}
					position++
					if !_rules[rulesp]() {
						goto l699
					}
					if !_rules[ruleExpression]() {
						goto l699
					}
					goto l698
				l699:
					position, tokenIndex, depth = position699, tokenIndex699, depth699
				}
				depth--
				add(ruleGroupList, position697)
			}
			return true
		l696:
			position, tokenIndex, depth = position696, tokenIndex696, depth696
			return false
		},
		/* 43 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action33)> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				{
					position702 := position
					depth++
					{
						position703, tokenIndex703, depth703 := position, tokenIndex, depth
						{
							position705, tokenIndex705, depth705 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l706
							}
							position++
							goto l705
						l706:
							position, tokenIndex, depth = position705, tokenIndex705, depth705
							if buffer[position] != rune('H') {
								goto l703
							}
							position++
						}
					l705:
						{
							position707, tokenIndex707, depth707 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l708
							}
							position++
							goto l707
						l708:
							position, tokenIndex, depth = position707, tokenIndex707, depth707
							if buffer[position] != rune('A') {
								goto l703
							}
							position++
						}
					l707:
						{
							position709, tokenIndex709, depth709 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l710
							}
							position++
							goto l709
						l710:
							position, tokenIndex, depth = position709, tokenIndex709, depth709
							if buffer[position] != rune('V') {
								goto l703
							}
							position++
						}
					l709:
						{
							position711, tokenIndex711, depth711 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l712
							}
							position++
							goto l711
						l712:
							position, tokenIndex, depth = position711, tokenIndex711, depth711
							if buffer[position] != rune('I') {
								goto l703
							}
							position++
						}
					l711:
						{
							position713, tokenIndex713, depth713 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l714
							}
							position++
							goto l713
						l714:
							position, tokenIndex, depth = position713, tokenIndex713, depth713
							if buffer[position] != rune('N') {
								goto l703
							}
							position++
						}
					l713:
						{
							position715, tokenIndex715, depth715 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l716
							}
							position++
							goto l715
						l716:
							position, tokenIndex, depth = position715, tokenIndex715, depth715
							if buffer[position] != rune('G') {
								goto l703
							}
							position++
						}
					l715:
						if !_rules[rulesp]() {
							goto l703
						}
						if !_rules[ruleExpression]() {
							goto l703
						}
						goto l704
					l703:
						position, tokenIndex, depth = position703, tokenIndex703, depth703
					}
				l704:
					depth--
					add(rulePegText, position702)
				}
				if !_rules[ruleAction33]() {
					goto l700
				}
				depth--
				add(ruleHaving, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 44 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action34))> */
		func() bool {
			position717, tokenIndex717, depth717 := position, tokenIndex, depth
			{
				position718 := position
				depth++
				{
					position719, tokenIndex719, depth719 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l720
					}
					goto l719
				l720:
					position, tokenIndex, depth = position719, tokenIndex719, depth719
					if !_rules[ruleStreamWindow]() {
						goto l717
					}
					if !_rules[ruleAction34]() {
						goto l717
					}
				}
			l719:
				depth--
				add(ruleRelationLike, position718)
			}
			return true
		l717:
			position, tokenIndex, depth = position717, tokenIndex717, depth717
			return false
		},
		/* 45 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action35)> */
		func() bool {
			position721, tokenIndex721, depth721 := position, tokenIndex, depth
			{
				position722 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l721
				}
				if !_rules[rulesp]() {
					goto l721
				}
				{
					position723, tokenIndex723, depth723 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l724
					}
					position++
					goto l723
				l724:
					position, tokenIndex, depth = position723, tokenIndex723, depth723
					if buffer[position] != rune('A') {
						goto l721
					}
					position++
				}
			l723:
				{
					position725, tokenIndex725, depth725 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l726
					}
					position++
					goto l725
				l726:
					position, tokenIndex, depth = position725, tokenIndex725, depth725
					if buffer[position] != rune('S') {
						goto l721
					}
					position++
				}
			l725:
				if !_rules[rulesp]() {
					goto l721
				}
				if !_rules[ruleIdentifier]() {
					goto l721
				}
				if !_rules[ruleAction35]() {
					goto l721
				}
				depth--
				add(ruleAliasedStreamWindow, position722)
			}
			return true
		l721:
			position, tokenIndex, depth = position721, tokenIndex721, depth721
			return false
		},
		/* 46 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action36)> */
		func() bool {
			position727, tokenIndex727, depth727 := position, tokenIndex, depth
			{
				position728 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l727
				}
				if !_rules[rulesp]() {
					goto l727
				}
				if buffer[position] != rune('[') {
					goto l727
				}
				position++
				if !_rules[rulesp]() {
					goto l727
				}
				{
					position729, tokenIndex729, depth729 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l730
					}
					position++
					goto l729
				l730:
					position, tokenIndex, depth = position729, tokenIndex729, depth729
					if buffer[position] != rune('R') {
						goto l727
					}
					position++
				}
			l729:
				{
					position731, tokenIndex731, depth731 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l732
					}
					position++
					goto l731
				l732:
					position, tokenIndex, depth = position731, tokenIndex731, depth731
					if buffer[position] != rune('A') {
						goto l727
					}
					position++
				}
			l731:
				{
					position733, tokenIndex733, depth733 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l734
					}
					position++
					goto l733
				l734:
					position, tokenIndex, depth = position733, tokenIndex733, depth733
					if buffer[position] != rune('N') {
						goto l727
					}
					position++
				}
			l733:
				{
					position735, tokenIndex735, depth735 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l736
					}
					position++
					goto l735
				l736:
					position, tokenIndex, depth = position735, tokenIndex735, depth735
					if buffer[position] != rune('G') {
						goto l727
					}
					position++
				}
			l735:
				{
					position737, tokenIndex737, depth737 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l738
					}
					position++
					goto l737
				l738:
					position, tokenIndex, depth = position737, tokenIndex737, depth737
					if buffer[position] != rune('E') {
						goto l727
					}
					position++
				}
			l737:
				if !_rules[rulesp]() {
					goto l727
				}
				if !_rules[ruleInterval]() {
					goto l727
				}
				if !_rules[rulesp]() {
					goto l727
				}
				if buffer[position] != rune(']') {
					goto l727
				}
				position++
				if !_rules[ruleAction36]() {
					goto l727
				}
				depth--
				add(ruleStreamWindow, position728)
			}
			return true
		l727:
			position, tokenIndex, depth = position727, tokenIndex727, depth727
			return false
		},
		/* 47 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position739, tokenIndex739, depth739 := position, tokenIndex, depth
			{
				position740 := position
				depth++
				{
					position741, tokenIndex741, depth741 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l742
					}
					goto l741
				l742:
					position, tokenIndex, depth = position741, tokenIndex741, depth741
					if !_rules[ruleStream]() {
						goto l739
					}
				}
			l741:
				depth--
				add(ruleStreamLike, position740)
			}
			return true
		l739:
			position, tokenIndex, depth = position739, tokenIndex739, depth739
			return false
		},
		/* 48 UDSFFuncApp <- <(FuncApp Action37)> */
		func() bool {
			position743, tokenIndex743, depth743 := position, tokenIndex, depth
			{
				position744 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l743
				}
				if !_rules[ruleAction37]() {
					goto l743
				}
				depth--
				add(ruleUDSFFuncApp, position744)
			}
			return true
		l743:
			position, tokenIndex, depth = position743, tokenIndex743, depth743
			return false
		},
		/* 49 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action38)> */
		func() bool {
			position745, tokenIndex745, depth745 := position, tokenIndex, depth
			{
				position746 := position
				depth++
				{
					position747 := position
					depth++
					{
						position748, tokenIndex748, depth748 := position, tokenIndex, depth
						{
							position750, tokenIndex750, depth750 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l751
							}
							position++
							goto l750
						l751:
							position, tokenIndex, depth = position750, tokenIndex750, depth750
							if buffer[position] != rune('W') {
								goto l748
							}
							position++
						}
					l750:
						{
							position752, tokenIndex752, depth752 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l753
							}
							position++
							goto l752
						l753:
							position, tokenIndex, depth = position752, tokenIndex752, depth752
							if buffer[position] != rune('I') {
								goto l748
							}
							position++
						}
					l752:
						{
							position754, tokenIndex754, depth754 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l755
							}
							position++
							goto l754
						l755:
							position, tokenIndex, depth = position754, tokenIndex754, depth754
							if buffer[position] != rune('T') {
								goto l748
							}
							position++
						}
					l754:
						{
							position756, tokenIndex756, depth756 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l757
							}
							position++
							goto l756
						l757:
							position, tokenIndex, depth = position756, tokenIndex756, depth756
							if buffer[position] != rune('H') {
								goto l748
							}
							position++
						}
					l756:
						if !_rules[rulesp]() {
							goto l748
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l748
						}
						if !_rules[rulesp]() {
							goto l748
						}
					l758:
						{
							position759, tokenIndex759, depth759 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l759
							}
							position++
							if !_rules[rulesp]() {
								goto l759
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l759
							}
							goto l758
						l759:
							position, tokenIndex, depth = position759, tokenIndex759, depth759
						}
						goto l749
					l748:
						position, tokenIndex, depth = position748, tokenIndex748, depth748
					}
				l749:
					depth--
					add(rulePegText, position747)
				}
				if !_rules[ruleAction38]() {
					goto l745
				}
				depth--
				add(ruleSourceSinkSpecs, position746)
			}
			return true
		l745:
			position, tokenIndex, depth = position745, tokenIndex745, depth745
			return false
		},
		/* 50 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action39)> */
		func() bool {
			position760, tokenIndex760, depth760 := position, tokenIndex, depth
			{
				position761 := position
				depth++
				{
					position762 := position
					depth++
					{
						position763, tokenIndex763, depth763 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l764
						}
						position++
						goto l763
					l764:
						position, tokenIndex, depth = position763, tokenIndex763, depth763
						if buffer[position] != rune('S') {
							goto l760
						}
						position++
					}
				l763:
					{
						position765, tokenIndex765, depth765 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l766
						}
						position++
						goto l765
					l766:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
						if buffer[position] != rune('E') {
							goto l760
						}
						position++
					}
				l765:
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('T') {
							goto l760
						}
						position++
					}
				l767:
					if !_rules[rulesp]() {
						goto l760
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l760
					}
					if !_rules[rulesp]() {
						goto l760
					}
				l769:
					{
						position770, tokenIndex770, depth770 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l770
						}
						position++
						if !_rules[rulesp]() {
							goto l770
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l770
						}
						goto l769
					l770:
						position, tokenIndex, depth = position770, tokenIndex770, depth770
					}
					depth--
					add(rulePegText, position762)
				}
				if !_rules[ruleAction39]() {
					goto l760
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position761)
			}
			return true
		l760:
			position, tokenIndex, depth = position760, tokenIndex760, depth760
			return false
		},
		/* 51 SetOptSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action40)> */
		func() bool {
			position771, tokenIndex771, depth771 := position, tokenIndex, depth
			{
				position772 := position
				depth++
				{
					position773 := position
					depth++
					{
						position774, tokenIndex774, depth774 := position, tokenIndex, depth
						{
							position776, tokenIndex776, depth776 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l777
							}
							position++
							goto l776
						l777:
							position, tokenIndex, depth = position776, tokenIndex776, depth776
							if buffer[position] != rune('S') {
								goto l774
							}
							position++
						}
					l776:
						{
							position778, tokenIndex778, depth778 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l779
							}
							position++
							goto l778
						l779:
							position, tokenIndex, depth = position778, tokenIndex778, depth778
							if buffer[position] != rune('E') {
								goto l774
							}
							position++
						}
					l778:
						{
							position780, tokenIndex780, depth780 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l781
							}
							position++
							goto l780
						l781:
							position, tokenIndex, depth = position780, tokenIndex780, depth780
							if buffer[position] != rune('T') {
								goto l774
							}
							position++
						}
					l780:
						if !_rules[rulesp]() {
							goto l774
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l774
						}
						if !_rules[rulesp]() {
							goto l774
						}
					l782:
						{
							position783, tokenIndex783, depth783 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l783
							}
							position++
							if !_rules[rulesp]() {
								goto l783
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l783
							}
							goto l782
						l783:
							position, tokenIndex, depth = position783, tokenIndex783, depth783
						}
						goto l775
					l774:
						position, tokenIndex, depth = position774, tokenIndex774, depth774
					}
				l775:
					depth--
					add(rulePegText, position773)
				}
				if !_rules[ruleAction40]() {
					goto l771
				}
				depth--
				add(ruleSetOptSpecs, position772)
			}
			return true
		l771:
			position, tokenIndex, depth = position771, tokenIndex771, depth771
			return false
		},
		/* 52 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action41)> */
		func() bool {
			position784, tokenIndex784, depth784 := position, tokenIndex, depth
			{
				position785 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l784
				}
				if buffer[position] != rune('=') {
					goto l784
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l784
				}
				if !_rules[ruleAction41]() {
					goto l784
				}
				depth--
				add(ruleSourceSinkParam, position785)
			}
			return true
		l784:
			position, tokenIndex, depth = position784, tokenIndex784, depth784
			return false
		},
		/* 53 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr)> */
		func() bool {
			position786, tokenIndex786, depth786 := position, tokenIndex, depth
			{
				position787 := position
				depth++
				{
					position788, tokenIndex788, depth788 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l789
					}
					goto l788
				l789:
					position, tokenIndex, depth = position788, tokenIndex788, depth788
					if !_rules[ruleParamArrayExpr]() {
						goto l786
					}
				}
			l788:
				depth--
				add(ruleSourceSinkParamVal, position787)
			}
			return true
		l786:
			position, tokenIndex, depth = position786, tokenIndex786, depth786
			return false
		},
		/* 54 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position790, tokenIndex790, depth790 := position, tokenIndex, depth
			{
				position791 := position
				depth++
				{
					position792, tokenIndex792, depth792 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l793
					}
					goto l792
				l793:
					position, tokenIndex, depth = position792, tokenIndex792, depth792
					if !_rules[ruleLiteral]() {
						goto l790
					}
				}
			l792:
				depth--
				add(ruleParamLiteral, position791)
			}
			return true
		l790:
			position, tokenIndex, depth = position790, tokenIndex790, depth790
			return false
		},
		/* 55 ParamArrayExpr <- <(<('[' sp (ParamLiteral (',' sp ParamLiteral)*)? sp ','? sp ']')> Action42)> */
		func() bool {
			position794, tokenIndex794, depth794 := position, tokenIndex, depth
			{
				position795 := position
				depth++
				{
					position796 := position
					depth++
					if buffer[position] != rune('[') {
						goto l794
					}
					position++
					if !_rules[rulesp]() {
						goto l794
					}
					{
						position797, tokenIndex797, depth797 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l797
						}
					l799:
						{
							position800, tokenIndex800, depth800 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l800
							}
							position++
							if !_rules[rulesp]() {
								goto l800
							}
							if !_rules[ruleParamLiteral]() {
								goto l800
							}
							goto l799
						l800:
							position, tokenIndex, depth = position800, tokenIndex800, depth800
						}
						goto l798
					l797:
						position, tokenIndex, depth = position797, tokenIndex797, depth797
					}
				l798:
					if !_rules[rulesp]() {
						goto l794
					}
					{
						position801, tokenIndex801, depth801 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l801
						}
						position++
						goto l802
					l801:
						position, tokenIndex, depth = position801, tokenIndex801, depth801
					}
				l802:
					if !_rules[rulesp]() {
						goto l794
					}
					if buffer[position] != rune(']') {
						goto l794
					}
					position++
					depth--
					add(rulePegText, position796)
				}
				if !_rules[ruleAction42]() {
					goto l794
				}
				depth--
				add(ruleParamArrayExpr, position795)
			}
			return true
		l794:
			position, tokenIndex, depth = position794, tokenIndex794, depth794
			return false
		},
		/* 56 PausedOpt <- <(<(Paused / Unpaused)?> Action43)> */
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
						{
							position808, tokenIndex808, depth808 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l809
							}
							goto l808
						l809:
							position, tokenIndex, depth = position808, tokenIndex808, depth808
							if !_rules[ruleUnpaused]() {
								goto l806
							}
						}
					l808:
						goto l807
					l806:
						position, tokenIndex, depth = position806, tokenIndex806, depth806
					}
				l807:
					depth--
					add(rulePegText, position805)
				}
				if !_rules[ruleAction43]() {
					goto l803
				}
				depth--
				add(rulePausedOpt, position804)
			}
			return true
		l803:
			position, tokenIndex, depth = position803, tokenIndex803, depth803
			return false
		},
		/* 57 Expression <- <orExpr> */
		func() bool {
			position810, tokenIndex810, depth810 := position, tokenIndex, depth
			{
				position811 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l810
				}
				depth--
				add(ruleExpression, position811)
			}
			return true
		l810:
			position, tokenIndex, depth = position810, tokenIndex810, depth810
			return false
		},
		/* 58 orExpr <- <(<(andExpr sp (Or sp andExpr)*)> Action44)> */
		func() bool {
			position812, tokenIndex812, depth812 := position, tokenIndex, depth
			{
				position813 := position
				depth++
				{
					position814 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l812
					}
					if !_rules[rulesp]() {
						goto l812
					}
				l815:
					{
						position816, tokenIndex816, depth816 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l816
						}
						if !_rules[rulesp]() {
							goto l816
						}
						if !_rules[ruleandExpr]() {
							goto l816
						}
						goto l815
					l816:
						position, tokenIndex, depth = position816, tokenIndex816, depth816
					}
					depth--
					add(rulePegText, position814)
				}
				if !_rules[ruleAction44]() {
					goto l812
				}
				depth--
				add(ruleorExpr, position813)
			}
			return true
		l812:
			position, tokenIndex, depth = position812, tokenIndex812, depth812
			return false
		},
		/* 59 andExpr <- <(<(notExpr sp (And sp notExpr)*)> Action45)> */
		func() bool {
			position817, tokenIndex817, depth817 := position, tokenIndex, depth
			{
				position818 := position
				depth++
				{
					position819 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l817
					}
					if !_rules[rulesp]() {
						goto l817
					}
				l820:
					{
						position821, tokenIndex821, depth821 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l821
						}
						if !_rules[rulesp]() {
							goto l821
						}
						if !_rules[rulenotExpr]() {
							goto l821
						}
						goto l820
					l821:
						position, tokenIndex, depth = position821, tokenIndex821, depth821
					}
					depth--
					add(rulePegText, position819)
				}
				if !_rules[ruleAction45]() {
					goto l817
				}
				depth--
				add(ruleandExpr, position818)
			}
			return true
		l817:
			position, tokenIndex, depth = position817, tokenIndex817, depth817
			return false
		},
		/* 60 notExpr <- <(<((Not sp)? comparisonExpr)> Action46)> */
		func() bool {
			position822, tokenIndex822, depth822 := position, tokenIndex, depth
			{
				position823 := position
				depth++
				{
					position824 := position
					depth++
					{
						position825, tokenIndex825, depth825 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l825
						}
						if !_rules[rulesp]() {
							goto l825
						}
						goto l826
					l825:
						position, tokenIndex, depth = position825, tokenIndex825, depth825
					}
				l826:
					if !_rules[rulecomparisonExpr]() {
						goto l822
					}
					depth--
					add(rulePegText, position824)
				}
				if !_rules[ruleAction46]() {
					goto l822
				}
				depth--
				add(rulenotExpr, position823)
			}
			return true
		l822:
			position, tokenIndex, depth = position822, tokenIndex822, depth822
			return false
		},
		/* 61 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action47)> */
		func() bool {
			position827, tokenIndex827, depth827 := position, tokenIndex, depth
			{
				position828 := position
				depth++
				{
					position829 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l827
					}
					if !_rules[rulesp]() {
						goto l827
					}
					{
						position830, tokenIndex830, depth830 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l830
						}
						if !_rules[rulesp]() {
							goto l830
						}
						if !_rules[ruleotherOpExpr]() {
							goto l830
						}
						goto l831
					l830:
						position, tokenIndex, depth = position830, tokenIndex830, depth830
					}
				l831:
					depth--
					add(rulePegText, position829)
				}
				if !_rules[ruleAction47]() {
					goto l827
				}
				depth--
				add(rulecomparisonExpr, position828)
			}
			return true
		l827:
			position, tokenIndex, depth = position827, tokenIndex827, depth827
			return false
		},
		/* 62 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action48)> */
		func() bool {
			position832, tokenIndex832, depth832 := position, tokenIndex, depth
			{
				position833 := position
				depth++
				{
					position834 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l832
					}
					if !_rules[rulesp]() {
						goto l832
					}
				l835:
					{
						position836, tokenIndex836, depth836 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l836
						}
						if !_rules[rulesp]() {
							goto l836
						}
						if !_rules[ruleisExpr]() {
							goto l836
						}
						if !_rules[rulesp]() {
							goto l836
						}
						goto l835
					l836:
						position, tokenIndex, depth = position836, tokenIndex836, depth836
					}
					depth--
					add(rulePegText, position834)
				}
				if !_rules[ruleAction48]() {
					goto l832
				}
				depth--
				add(ruleotherOpExpr, position833)
			}
			return true
		l832:
			position, tokenIndex, depth = position832, tokenIndex832, depth832
			return false
		},
		/* 63 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action49)> */
		func() bool {
			position837, tokenIndex837, depth837 := position, tokenIndex, depth
			{
				position838 := position
				depth++
				{
					position839 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l837
					}
					if !_rules[rulesp]() {
						goto l837
					}
					{
						position840, tokenIndex840, depth840 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l840
						}
						if !_rules[rulesp]() {
							goto l840
						}
						if !_rules[ruleNullLiteral]() {
							goto l840
						}
						goto l841
					l840:
						position, tokenIndex, depth = position840, tokenIndex840, depth840
					}
				l841:
					depth--
					add(rulePegText, position839)
				}
				if !_rules[ruleAction49]() {
					goto l837
				}
				depth--
				add(ruleisExpr, position838)
			}
			return true
		l837:
			position, tokenIndex, depth = position837, tokenIndex837, depth837
			return false
		},
		/* 64 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action50)> */
		func() bool {
			position842, tokenIndex842, depth842 := position, tokenIndex, depth
			{
				position843 := position
				depth++
				{
					position844 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l842
					}
					if !_rules[rulesp]() {
						goto l842
					}
				l845:
					{
						position846, tokenIndex846, depth846 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l846
						}
						if !_rules[rulesp]() {
							goto l846
						}
						if !_rules[ruleproductExpr]() {
							goto l846
						}
						if !_rules[rulesp]() {
							goto l846
						}
						goto l845
					l846:
						position, tokenIndex, depth = position846, tokenIndex846, depth846
					}
					depth--
					add(rulePegText, position844)
				}
				if !_rules[ruleAction50]() {
					goto l842
				}
				depth--
				add(ruletermExpr, position843)
			}
			return true
		l842:
			position, tokenIndex, depth = position842, tokenIndex842, depth842
			return false
		},
		/* 65 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action51)> */
		func() bool {
			position847, tokenIndex847, depth847 := position, tokenIndex, depth
			{
				position848 := position
				depth++
				{
					position849 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l847
					}
					if !_rules[rulesp]() {
						goto l847
					}
				l850:
					{
						position851, tokenIndex851, depth851 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l851
						}
						if !_rules[rulesp]() {
							goto l851
						}
						if !_rules[ruleminusExpr]() {
							goto l851
						}
						if !_rules[rulesp]() {
							goto l851
						}
						goto l850
					l851:
						position, tokenIndex, depth = position851, tokenIndex851, depth851
					}
					depth--
					add(rulePegText, position849)
				}
				if !_rules[ruleAction51]() {
					goto l847
				}
				depth--
				add(ruleproductExpr, position848)
			}
			return true
		l847:
			position, tokenIndex, depth = position847, tokenIndex847, depth847
			return false
		},
		/* 66 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action52)> */
		func() bool {
			position852, tokenIndex852, depth852 := position, tokenIndex, depth
			{
				position853 := position
				depth++
				{
					position854 := position
					depth++
					{
						position855, tokenIndex855, depth855 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l855
						}
						if !_rules[rulesp]() {
							goto l855
						}
						goto l856
					l855:
						position, tokenIndex, depth = position855, tokenIndex855, depth855
					}
				l856:
					if !_rules[rulecastExpr]() {
						goto l852
					}
					depth--
					add(rulePegText, position854)
				}
				if !_rules[ruleAction52]() {
					goto l852
				}
				depth--
				add(ruleminusExpr, position853)
			}
			return true
		l852:
			position, tokenIndex, depth = position852, tokenIndex852, depth852
			return false
		},
		/* 67 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action53)> */
		func() bool {
			position857, tokenIndex857, depth857 := position, tokenIndex, depth
			{
				position858 := position
				depth++
				{
					position859 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l857
					}
					{
						position860, tokenIndex860, depth860 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l860
						}
						if buffer[position] != rune(':') {
							goto l860
						}
						position++
						if buffer[position] != rune(':') {
							goto l860
						}
						position++
						if !_rules[rulesp]() {
							goto l860
						}
						if !_rules[ruleType]() {
							goto l860
						}
						goto l861
					l860:
						position, tokenIndex, depth = position860, tokenIndex860, depth860
					}
				l861:
					depth--
					add(rulePegText, position859)
				}
				if !_rules[ruleAction53]() {
					goto l857
				}
				depth--
				add(rulecastExpr, position858)
			}
			return true
		l857:
			position, tokenIndex, depth = position857, tokenIndex857, depth857
			return false
		},
		/* 68 baseExpr <- <(('(' sp Expression sp ')') / MapExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position862, tokenIndex862, depth862 := position, tokenIndex, depth
			{
				position863 := position
				depth++
				{
					position864, tokenIndex864, depth864 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l865
					}
					position++
					if !_rules[rulesp]() {
						goto l865
					}
					if !_rules[ruleExpression]() {
						goto l865
					}
					if !_rules[rulesp]() {
						goto l865
					}
					if buffer[position] != rune(')') {
						goto l865
					}
					position++
					goto l864
				l865:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleMapExpr]() {
						goto l866
					}
					goto l864
				l866:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleBooleanLiteral]() {
						goto l867
					}
					goto l864
				l867:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleNullLiteral]() {
						goto l868
					}
					goto l864
				l868:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleRowMeta]() {
						goto l869
					}
					goto l864
				l869:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleFuncTypeCast]() {
						goto l870
					}
					goto l864
				l870:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleFuncApp]() {
						goto l871
					}
					goto l864
				l871:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleRowValue]() {
						goto l872
					}
					goto l864
				l872:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleArrayExpr]() {
						goto l873
					}
					goto l864
				l873:
					position, tokenIndex, depth = position864, tokenIndex864, depth864
					if !_rules[ruleLiteral]() {
						goto l862
					}
				}
			l864:
				depth--
				add(rulebaseExpr, position863)
			}
			return true
		l862:
			position, tokenIndex, depth = position862, tokenIndex862, depth862
			return false
		},
		/* 69 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action54)> */
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
						if buffer[position] != rune('c') {
							goto l878
						}
						position++
						goto l877
					l878:
						position, tokenIndex, depth = position877, tokenIndex877, depth877
						if buffer[position] != rune('C') {
							goto l874
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
							goto l874
						}
						position++
					}
				l879:
					{
						position881, tokenIndex881, depth881 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l882
						}
						position++
						goto l881
					l882:
						position, tokenIndex, depth = position881, tokenIndex881, depth881
						if buffer[position] != rune('S') {
							goto l874
						}
						position++
					}
				l881:
					{
						position883, tokenIndex883, depth883 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l884
						}
						position++
						goto l883
					l884:
						position, tokenIndex, depth = position883, tokenIndex883, depth883
						if buffer[position] != rune('T') {
							goto l874
						}
						position++
					}
				l883:
					if !_rules[rulesp]() {
						goto l874
					}
					if buffer[position] != rune('(') {
						goto l874
					}
					position++
					if !_rules[rulesp]() {
						goto l874
					}
					if !_rules[ruleExpression]() {
						goto l874
					}
					if !_rules[rulesp]() {
						goto l874
					}
					{
						position885, tokenIndex885, depth885 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l886
						}
						position++
						goto l885
					l886:
						position, tokenIndex, depth = position885, tokenIndex885, depth885
						if buffer[position] != rune('A') {
							goto l874
						}
						position++
					}
				l885:
					{
						position887, tokenIndex887, depth887 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l888
						}
						position++
						goto l887
					l888:
						position, tokenIndex, depth = position887, tokenIndex887, depth887
						if buffer[position] != rune('S') {
							goto l874
						}
						position++
					}
				l887:
					if !_rules[rulesp]() {
						goto l874
					}
					if !_rules[ruleType]() {
						goto l874
					}
					if !_rules[rulesp]() {
						goto l874
					}
					if buffer[position] != rune(')') {
						goto l874
					}
					position++
					depth--
					add(rulePegText, position876)
				}
				if !_rules[ruleAction54]() {
					goto l874
				}
				depth--
				add(ruleFuncTypeCast, position875)
			}
			return true
		l874:
			position, tokenIndex, depth = position874, tokenIndex874, depth874
			return false
		},
		/* 70 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action55)> */
		func() bool {
			position889, tokenIndex889, depth889 := position, tokenIndex, depth
			{
				position890 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l889
				}
				if !_rules[rulesp]() {
					goto l889
				}
				if buffer[position] != rune('(') {
					goto l889
				}
				position++
				if !_rules[rulesp]() {
					goto l889
				}
				if !_rules[ruleFuncParams]() {
					goto l889
				}
				if !_rules[rulesp]() {
					goto l889
				}
				if buffer[position] != rune(')') {
					goto l889
				}
				position++
				if !_rules[ruleAction55]() {
					goto l889
				}
				depth--
				add(ruleFuncApp, position890)
			}
			return true
		l889:
			position, tokenIndex, depth = position889, tokenIndex889, depth889
			return false
		},
		/* 71 FuncParams <- <(<(Star / (Expression sp (',' sp Expression)*)?)> Action56)> */
		func() bool {
			position891, tokenIndex891, depth891 := position, tokenIndex, depth
			{
				position892 := position
				depth++
				{
					position893 := position
					depth++
					{
						position894, tokenIndex894, depth894 := position, tokenIndex, depth
						if !_rules[ruleStar]() {
							goto l895
						}
						goto l894
					l895:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
						{
							position896, tokenIndex896, depth896 := position, tokenIndex, depth
							if !_rules[ruleExpression]() {
								goto l896
							}
							if !_rules[rulesp]() {
								goto l896
							}
						l898:
							{
								position899, tokenIndex899, depth899 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l899
								}
								position++
								if !_rules[rulesp]() {
									goto l899
								}
								if !_rules[ruleExpression]() {
									goto l899
								}
								goto l898
							l899:
								position, tokenIndex, depth = position899, tokenIndex899, depth899
							}
							goto l897
						l896:
							position, tokenIndex, depth = position896, tokenIndex896, depth896
						}
					l897:
					}
				l894:
					depth--
					add(rulePegText, position893)
				}
				if !_rules[ruleAction56]() {
					goto l891
				}
				depth--
				add(ruleFuncParams, position892)
			}
			return true
		l891:
			position, tokenIndex, depth = position891, tokenIndex891, depth891
			return false
		},
		/* 72 ArrayExpr <- <(<('[' sp (Expression (',' sp Expression)*)? sp ','? sp ']')> Action57)> */
		func() bool {
			position900, tokenIndex900, depth900 := position, tokenIndex, depth
			{
				position901 := position
				depth++
				{
					position902 := position
					depth++
					if buffer[position] != rune('[') {
						goto l900
					}
					position++
					if !_rules[rulesp]() {
						goto l900
					}
					{
						position903, tokenIndex903, depth903 := position, tokenIndex, depth
						if !_rules[ruleExpression]() {
							goto l903
						}
					l905:
						{
							position906, tokenIndex906, depth906 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l906
							}
							position++
							if !_rules[rulesp]() {
								goto l906
							}
							if !_rules[ruleExpression]() {
								goto l906
							}
							goto l905
						l906:
							position, tokenIndex, depth = position906, tokenIndex906, depth906
						}
						goto l904
					l903:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
					}
				l904:
					if !_rules[rulesp]() {
						goto l900
					}
					{
						position907, tokenIndex907, depth907 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l907
						}
						position++
						goto l908
					l907:
						position, tokenIndex, depth = position907, tokenIndex907, depth907
					}
				l908:
					if !_rules[rulesp]() {
						goto l900
					}
					if buffer[position] != rune(']') {
						goto l900
					}
					position++
					depth--
					add(rulePegText, position902)
				}
				if !_rules[ruleAction57]() {
					goto l900
				}
				depth--
				add(ruleArrayExpr, position901)
			}
			return true
		l900:
			position, tokenIndex, depth = position900, tokenIndex900, depth900
			return false
		},
		/* 73 MapExpr <- <(<('{' sp (KeyValuePair (',' sp KeyValuePair)*)? sp '}')> Action58)> */
		func() bool {
			position909, tokenIndex909, depth909 := position, tokenIndex, depth
			{
				position910 := position
				depth++
				{
					position911 := position
					depth++
					if buffer[position] != rune('{') {
						goto l909
					}
					position++
					if !_rules[rulesp]() {
						goto l909
					}
					{
						position912, tokenIndex912, depth912 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l912
						}
					l914:
						{
							position915, tokenIndex915, depth915 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l915
							}
							position++
							if !_rules[rulesp]() {
								goto l915
							}
							if !_rules[ruleKeyValuePair]() {
								goto l915
							}
							goto l914
						l915:
							position, tokenIndex, depth = position915, tokenIndex915, depth915
						}
						goto l913
					l912:
						position, tokenIndex, depth = position912, tokenIndex912, depth912
					}
				l913:
					if !_rules[rulesp]() {
						goto l909
					}
					if buffer[position] != rune('}') {
						goto l909
					}
					position++
					depth--
					add(rulePegText, position911)
				}
				if !_rules[ruleAction58]() {
					goto l909
				}
				depth--
				add(ruleMapExpr, position910)
			}
			return true
		l909:
			position, tokenIndex, depth = position909, tokenIndex909, depth909
			return false
		},
		/* 74 KeyValuePair <- <(<(StringLiteral sp ':' sp Expression)> Action59)> */
		func() bool {
			position916, tokenIndex916, depth916 := position, tokenIndex, depth
			{
				position917 := position
				depth++
				{
					position918 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l916
					}
					if !_rules[rulesp]() {
						goto l916
					}
					if buffer[position] != rune(':') {
						goto l916
					}
					position++
					if !_rules[rulesp]() {
						goto l916
					}
					if !_rules[ruleExpression]() {
						goto l916
					}
					depth--
					add(rulePegText, position918)
				}
				if !_rules[ruleAction59]() {
					goto l916
				}
				depth--
				add(ruleKeyValuePair, position917)
			}
			return true
		l916:
			position, tokenIndex, depth = position916, tokenIndex916, depth916
			return false
		},
		/* 75 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position919, tokenIndex919, depth919 := position, tokenIndex, depth
			{
				position920 := position
				depth++
				{
					position921, tokenIndex921, depth921 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l922
					}
					goto l921
				l922:
					position, tokenIndex, depth = position921, tokenIndex921, depth921
					if !_rules[ruleNumericLiteral]() {
						goto l923
					}
					goto l921
				l923:
					position, tokenIndex, depth = position921, tokenIndex921, depth921
					if !_rules[ruleStringLiteral]() {
						goto l919
					}
				}
			l921:
				depth--
				add(ruleLiteral, position920)
			}
			return true
		l919:
			position, tokenIndex, depth = position919, tokenIndex919, depth919
			return false
		},
		/* 76 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position924, tokenIndex924, depth924 := position, tokenIndex, depth
			{
				position925 := position
				depth++
				{
					position926, tokenIndex926, depth926 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l927
					}
					goto l926
				l927:
					position, tokenIndex, depth = position926, tokenIndex926, depth926
					if !_rules[ruleNotEqual]() {
						goto l928
					}
					goto l926
				l928:
					position, tokenIndex, depth = position926, tokenIndex926, depth926
					if !_rules[ruleLessOrEqual]() {
						goto l929
					}
					goto l926
				l929:
					position, tokenIndex, depth = position926, tokenIndex926, depth926
					if !_rules[ruleLess]() {
						goto l930
					}
					goto l926
				l930:
					position, tokenIndex, depth = position926, tokenIndex926, depth926
					if !_rules[ruleGreaterOrEqual]() {
						goto l931
					}
					goto l926
				l931:
					position, tokenIndex, depth = position926, tokenIndex926, depth926
					if !_rules[ruleGreater]() {
						goto l932
					}
					goto l926
				l932:
					position, tokenIndex, depth = position926, tokenIndex926, depth926
					if !_rules[ruleNotEqual]() {
						goto l924
					}
				}
			l926:
				depth--
				add(ruleComparisonOp, position925)
			}
			return true
		l924:
			position, tokenIndex, depth = position924, tokenIndex924, depth924
			return false
		},
		/* 77 OtherOp <- <Concat> */
		func() bool {
			position933, tokenIndex933, depth933 := position, tokenIndex, depth
			{
				position934 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l933
				}
				depth--
				add(ruleOtherOp, position934)
			}
			return true
		l933:
			position, tokenIndex, depth = position933, tokenIndex933, depth933
			return false
		},
		/* 78 IsOp <- <(IsNot / Is)> */
		func() bool {
			position935, tokenIndex935, depth935 := position, tokenIndex, depth
			{
				position936 := position
				depth++
				{
					position937, tokenIndex937, depth937 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l938
					}
					goto l937
				l938:
					position, tokenIndex, depth = position937, tokenIndex937, depth937
					if !_rules[ruleIs]() {
						goto l935
					}
				}
			l937:
				depth--
				add(ruleIsOp, position936)
			}
			return true
		l935:
			position, tokenIndex, depth = position935, tokenIndex935, depth935
			return false
		},
		/* 79 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position939, tokenIndex939, depth939 := position, tokenIndex, depth
			{
				position940 := position
				depth++
				{
					position941, tokenIndex941, depth941 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l942
					}
					goto l941
				l942:
					position, tokenIndex, depth = position941, tokenIndex941, depth941
					if !_rules[ruleMinus]() {
						goto l939
					}
				}
			l941:
				depth--
				add(rulePlusMinusOp, position940)
			}
			return true
		l939:
			position, tokenIndex, depth = position939, tokenIndex939, depth939
			return false
		},
		/* 80 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position943, tokenIndex943, depth943 := position, tokenIndex, depth
			{
				position944 := position
				depth++
				{
					position945, tokenIndex945, depth945 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l946
					}
					goto l945
				l946:
					position, tokenIndex, depth = position945, tokenIndex945, depth945
					if !_rules[ruleDivide]() {
						goto l947
					}
					goto l945
				l947:
					position, tokenIndex, depth = position945, tokenIndex945, depth945
					if !_rules[ruleModulo]() {
						goto l943
					}
				}
			l945:
				depth--
				add(ruleMultDivOp, position944)
			}
			return true
		l943:
			position, tokenIndex, depth = position943, tokenIndex943, depth943
			return false
		},
		/* 81 Stream <- <(<ident> Action60)> */
		func() bool {
			position948, tokenIndex948, depth948 := position, tokenIndex, depth
			{
				position949 := position
				depth++
				{
					position950 := position
					depth++
					if !_rules[ruleident]() {
						goto l948
					}
					depth--
					add(rulePegText, position950)
				}
				if !_rules[ruleAction60]() {
					goto l948
				}
				depth--
				add(ruleStream, position949)
			}
			return true
		l948:
			position, tokenIndex, depth = position948, tokenIndex948, depth948
			return false
		},
		/* 82 RowMeta <- <RowTimestamp> */
		func() bool {
			position951, tokenIndex951, depth951 := position, tokenIndex, depth
			{
				position952 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l951
				}
				depth--
				add(ruleRowMeta, position952)
			}
			return true
		l951:
			position, tokenIndex, depth = position951, tokenIndex951, depth951
			return false
		},
		/* 83 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action61)> */
		func() bool {
			position953, tokenIndex953, depth953 := position, tokenIndex, depth
			{
				position954 := position
				depth++
				{
					position955 := position
					depth++
					{
						position956, tokenIndex956, depth956 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l956
						}
						if buffer[position] != rune(':') {
							goto l956
						}
						position++
						goto l957
					l956:
						position, tokenIndex, depth = position956, tokenIndex956, depth956
					}
				l957:
					if buffer[position] != rune('t') {
						goto l953
					}
					position++
					if buffer[position] != rune('s') {
						goto l953
					}
					position++
					if buffer[position] != rune('(') {
						goto l953
					}
					position++
					if buffer[position] != rune(')') {
						goto l953
					}
					position++
					depth--
					add(rulePegText, position955)
				}
				if !_rules[ruleAction61]() {
					goto l953
				}
				depth--
				add(ruleRowTimestamp, position954)
			}
			return true
		l953:
			position, tokenIndex, depth = position953, tokenIndex953, depth953
			return false
		},
		/* 84 RowValue <- <(<((ident ':' !':')? jsonPath)> Action62)> */
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
						if !_rules[ruleident]() {
							goto l961
						}
						if buffer[position] != rune(':') {
							goto l961
						}
						position++
						{
							position963, tokenIndex963, depth963 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l963
							}
							position++
							goto l961
						l963:
							position, tokenIndex, depth = position963, tokenIndex963, depth963
						}
						goto l962
					l961:
						position, tokenIndex, depth = position961, tokenIndex961, depth961
					}
				l962:
					if !_rules[rulejsonPath]() {
						goto l958
					}
					depth--
					add(rulePegText, position960)
				}
				if !_rules[ruleAction62]() {
					goto l958
				}
				depth--
				add(ruleRowValue, position959)
			}
			return true
		l958:
			position, tokenIndex, depth = position958, tokenIndex958, depth958
			return false
		},
		/* 85 NumericLiteral <- <(<('-'? [0-9]+)> Action63)> */
		func() bool {
			position964, tokenIndex964, depth964 := position, tokenIndex, depth
			{
				position965 := position
				depth++
				{
					position966 := position
					depth++
					{
						position967, tokenIndex967, depth967 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l967
						}
						position++
						goto l968
					l967:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
					}
				l968:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l964
					}
					position++
				l969:
					{
						position970, tokenIndex970, depth970 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l970
						}
						position++
						goto l969
					l970:
						position, tokenIndex, depth = position970, tokenIndex970, depth970
					}
					depth--
					add(rulePegText, position966)
				}
				if !_rules[ruleAction63]() {
					goto l964
				}
				depth--
				add(ruleNumericLiteral, position965)
			}
			return true
		l964:
			position, tokenIndex, depth = position964, tokenIndex964, depth964
			return false
		},
		/* 86 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action64)> */
		func() bool {
			position971, tokenIndex971, depth971 := position, tokenIndex, depth
			{
				position972 := position
				depth++
				{
					position973 := position
					depth++
					{
						position974, tokenIndex974, depth974 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l974
						}
						position++
						goto l975
					l974:
						position, tokenIndex, depth = position974, tokenIndex974, depth974
					}
				l975:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l971
					}
					position++
				l976:
					{
						position977, tokenIndex977, depth977 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l977
						}
						position++
						goto l976
					l977:
						position, tokenIndex, depth = position977, tokenIndex977, depth977
					}
					if buffer[position] != rune('.') {
						goto l971
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l971
					}
					position++
				l978:
					{
						position979, tokenIndex979, depth979 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l979
						}
						position++
						goto l978
					l979:
						position, tokenIndex, depth = position979, tokenIndex979, depth979
					}
					depth--
					add(rulePegText, position973)
				}
				if !_rules[ruleAction64]() {
					goto l971
				}
				depth--
				add(ruleFloatLiteral, position972)
			}
			return true
		l971:
			position, tokenIndex, depth = position971, tokenIndex971, depth971
			return false
		},
		/* 87 Function <- <(<ident> Action65)> */
		func() bool {
			position980, tokenIndex980, depth980 := position, tokenIndex, depth
			{
				position981 := position
				depth++
				{
					position982 := position
					depth++
					if !_rules[ruleident]() {
						goto l980
					}
					depth--
					add(rulePegText, position982)
				}
				if !_rules[ruleAction65]() {
					goto l980
				}
				depth--
				add(ruleFunction, position981)
			}
			return true
		l980:
			position, tokenIndex, depth = position980, tokenIndex980, depth980
			return false
		},
		/* 88 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action66)> */
		func() bool {
			position983, tokenIndex983, depth983 := position, tokenIndex, depth
			{
				position984 := position
				depth++
				{
					position985 := position
					depth++
					{
						position986, tokenIndex986, depth986 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l987
						}
						position++
						goto l986
					l987:
						position, tokenIndex, depth = position986, tokenIndex986, depth986
						if buffer[position] != rune('N') {
							goto l983
						}
						position++
					}
				l986:
					{
						position988, tokenIndex988, depth988 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l989
						}
						position++
						goto l988
					l989:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('U') {
							goto l983
						}
						position++
					}
				l988:
					{
						position990, tokenIndex990, depth990 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l991
						}
						position++
						goto l990
					l991:
						position, tokenIndex, depth = position990, tokenIndex990, depth990
						if buffer[position] != rune('L') {
							goto l983
						}
						position++
					}
				l990:
					{
						position992, tokenIndex992, depth992 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l993
						}
						position++
						goto l992
					l993:
						position, tokenIndex, depth = position992, tokenIndex992, depth992
						if buffer[position] != rune('L') {
							goto l983
						}
						position++
					}
				l992:
					depth--
					add(rulePegText, position985)
				}
				if !_rules[ruleAction66]() {
					goto l983
				}
				depth--
				add(ruleNullLiteral, position984)
			}
			return true
		l983:
			position, tokenIndex, depth = position983, tokenIndex983, depth983
			return false
		},
		/* 89 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position994, tokenIndex994, depth994 := position, tokenIndex, depth
			{
				position995 := position
				depth++
				{
					position996, tokenIndex996, depth996 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l997
					}
					goto l996
				l997:
					position, tokenIndex, depth = position996, tokenIndex996, depth996
					if !_rules[ruleFALSE]() {
						goto l994
					}
				}
			l996:
				depth--
				add(ruleBooleanLiteral, position995)
			}
			return true
		l994:
			position, tokenIndex, depth = position994, tokenIndex994, depth994
			return false
		},
		/* 90 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action67)> */
		func() bool {
			position998, tokenIndex998, depth998 := position, tokenIndex, depth
			{
				position999 := position
				depth++
				{
					position1000 := position
					depth++
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
							goto l998
						}
						position++
					}
				l1001:
					{
						position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1004
						}
						position++
						goto l1003
					l1004:
						position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
						if buffer[position] != rune('R') {
							goto l998
						}
						position++
					}
				l1003:
					{
						position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1006
						}
						position++
						goto l1005
					l1006:
						position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
						if buffer[position] != rune('U') {
							goto l998
						}
						position++
					}
				l1005:
					{
						position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1008
						}
						position++
						goto l1007
					l1008:
						position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
						if buffer[position] != rune('E') {
							goto l998
						}
						position++
					}
				l1007:
					depth--
					add(rulePegText, position1000)
				}
				if !_rules[ruleAction67]() {
					goto l998
				}
				depth--
				add(ruleTRUE, position999)
			}
			return true
		l998:
			position, tokenIndex, depth = position998, tokenIndex998, depth998
			return false
		},
		/* 91 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action68)> */
		func() bool {
			position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
			{
				position1010 := position
				depth++
				{
					position1011 := position
					depth++
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('F') {
							goto l1009
						}
						position++
					}
				l1012:
					{
						position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1015
						}
						position++
						goto l1014
					l1015:
						position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
						if buffer[position] != rune('A') {
							goto l1009
						}
						position++
					}
				l1014:
					{
						position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1017
						}
						position++
						goto l1016
					l1017:
						position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
						if buffer[position] != rune('L') {
							goto l1009
						}
						position++
					}
				l1016:
					{
						position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1019
						}
						position++
						goto l1018
					l1019:
						position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
						if buffer[position] != rune('S') {
							goto l1009
						}
						position++
					}
				l1018:
					{
						position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1021
						}
						position++
						goto l1020
					l1021:
						position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
						if buffer[position] != rune('E') {
							goto l1009
						}
						position++
					}
				l1020:
					depth--
					add(rulePegText, position1011)
				}
				if !_rules[ruleAction68]() {
					goto l1009
				}
				depth--
				add(ruleFALSE, position1010)
			}
			return true
		l1009:
			position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
			return false
		},
		/* 92 Star <- <(<'*'> Action69)> */
		func() bool {
			position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
			{
				position1023 := position
				depth++
				{
					position1024 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1022
					}
					position++
					depth--
					add(rulePegText, position1024)
				}
				if !_rules[ruleAction69]() {
					goto l1022
				}
				depth--
				add(ruleStar, position1023)
			}
			return true
		l1022:
			position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
			return false
		},
		/* 93 Wildcard <- <(<((ident ':' !':')? '*')> Action70)> */
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
						if !_rules[ruleident]() {
							goto l1028
						}
						if buffer[position] != rune(':') {
							goto l1028
						}
						position++
						{
							position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1030
							}
							position++
							goto l1028
						l1030:
							position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
						}
						goto l1029
					l1028:
						position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
					}
				l1029:
					if buffer[position] != rune('*') {
						goto l1025
					}
					position++
					depth--
					add(rulePegText, position1027)
				}
				if !_rules[ruleAction70]() {
					goto l1025
				}
				depth--
				add(ruleWildcard, position1026)
			}
			return true
		l1025:
			position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
			return false
		},
		/* 94 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action71)> */
		func() bool {
			position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
			{
				position1032 := position
				depth++
				{
					position1033 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l1031
					}
					position++
				l1034:
					{
						position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
						{
							position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1037
							}
							position++
							if buffer[position] != rune('\'') {
								goto l1037
							}
							position++
							goto l1036
						l1037:
							position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
							{
								position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l1038
								}
								position++
								goto l1035
							l1038:
								position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
							}
							if !matchDot() {
								goto l1035
							}
						}
					l1036:
						goto l1034
					l1035:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
					}
					if buffer[position] != rune('\'') {
						goto l1031
					}
					position++
					depth--
					add(rulePegText, position1033)
				}
				if !_rules[ruleAction71]() {
					goto l1031
				}
				depth--
				add(ruleStringLiteral, position1032)
			}
			return true
		l1031:
			position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
			return false
		},
		/* 95 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action72)> */
		func() bool {
			position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
			{
				position1040 := position
				depth++
				{
					position1041 := position
					depth++
					{
						position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1043
						}
						position++
						goto l1042
					l1043:
						position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
						if buffer[position] != rune('I') {
							goto l1039
						}
						position++
					}
				l1042:
					{
						position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1045
						}
						position++
						goto l1044
					l1045:
						position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
						if buffer[position] != rune('S') {
							goto l1039
						}
						position++
					}
				l1044:
					{
						position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1047
						}
						position++
						goto l1046
					l1047:
						position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
						if buffer[position] != rune('T') {
							goto l1039
						}
						position++
					}
				l1046:
					{
						position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1049
						}
						position++
						goto l1048
					l1049:
						position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
						if buffer[position] != rune('R') {
							goto l1039
						}
						position++
					}
				l1048:
					{
						position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1051
						}
						position++
						goto l1050
					l1051:
						position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
						if buffer[position] != rune('E') {
							goto l1039
						}
						position++
					}
				l1050:
					{
						position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1053
						}
						position++
						goto l1052
					l1053:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
						if buffer[position] != rune('A') {
							goto l1039
						}
						position++
					}
				l1052:
					{
						position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1055
						}
						position++
						goto l1054
					l1055:
						position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
						if buffer[position] != rune('M') {
							goto l1039
						}
						position++
					}
				l1054:
					depth--
					add(rulePegText, position1041)
				}
				if !_rules[ruleAction72]() {
					goto l1039
				}
				depth--
				add(ruleISTREAM, position1040)
			}
			return true
		l1039:
			position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
			return false
		},
		/* 96 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action73)> */
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
						if buffer[position] != rune('d') {
							goto l1060
						}
						position++
						goto l1059
					l1060:
						position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
						if buffer[position] != rune('D') {
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
				if !_rules[ruleAction73]() {
					goto l1056
				}
				depth--
				add(ruleDSTREAM, position1057)
			}
			return true
		l1056:
			position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
			return false
		},
		/* 97 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action74)> */
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
						if buffer[position] != rune('r') {
							goto l1077
						}
						position++
						goto l1076
					l1077:
						position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
						if buffer[position] != rune('R') {
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
				if !_rules[ruleAction74]() {
					goto l1073
				}
				depth--
				add(ruleRSTREAM, position1074)
			}
			return true
		l1073:
			position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
			return false
		},
		/* 98 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action75)> */
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
						if buffer[position] != rune('t') {
							goto l1094
						}
						position++
						goto l1093
					l1094:
						position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
						if buffer[position] != rune('T') {
							goto l1090
						}
						position++
					}
				l1093:
					{
						position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1096
						}
						position++
						goto l1095
					l1096:
						position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
						if buffer[position] != rune('U') {
							goto l1090
						}
						position++
					}
				l1095:
					{
						position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1098
						}
						position++
						goto l1097
					l1098:
						position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
						if buffer[position] != rune('P') {
							goto l1090
						}
						position++
					}
				l1097:
					{
						position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1100
						}
						position++
						goto l1099
					l1100:
						position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
						if buffer[position] != rune('L') {
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
						if buffer[position] != rune('s') {
							goto l1104
						}
						position++
						goto l1103
					l1104:
						position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
						if buffer[position] != rune('S') {
							goto l1090
						}
						position++
					}
				l1103:
					depth--
					add(rulePegText, position1092)
				}
				if !_rules[ruleAction75]() {
					goto l1090
				}
				depth--
				add(ruleTUPLES, position1091)
			}
			return true
		l1090:
			position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
			return false
		},
		/* 99 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action76)> */
		func() bool {
			position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
			{
				position1106 := position
				depth++
				{
					position1107 := position
					depth++
					{
						position1108, tokenIndex1108, depth1108 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1109
						}
						position++
						goto l1108
					l1109:
						position, tokenIndex, depth = position1108, tokenIndex1108, depth1108
						if buffer[position] != rune('S') {
							goto l1105
						}
						position++
					}
				l1108:
					{
						position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1111
						}
						position++
						goto l1110
					l1111:
						position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
						if buffer[position] != rune('E') {
							goto l1105
						}
						position++
					}
				l1110:
					{
						position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1113
						}
						position++
						goto l1112
					l1113:
						position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
						if buffer[position] != rune('C') {
							goto l1105
						}
						position++
					}
				l1112:
					{
						position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1115
						}
						position++
						goto l1114
					l1115:
						position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
						if buffer[position] != rune('O') {
							goto l1105
						}
						position++
					}
				l1114:
					{
						position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1117
						}
						position++
						goto l1116
					l1117:
						position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
						if buffer[position] != rune('N') {
							goto l1105
						}
						position++
					}
				l1116:
					{
						position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1119
						}
						position++
						goto l1118
					l1119:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						if buffer[position] != rune('D') {
							goto l1105
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
							goto l1105
						}
						position++
					}
				l1120:
					depth--
					add(rulePegText, position1107)
				}
				if !_rules[ruleAction76]() {
					goto l1105
				}
				depth--
				add(ruleSECONDS, position1106)
			}
			return true
		l1105:
			position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
			return false
		},
		/* 100 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action77)> */
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
						if buffer[position] != rune('m') {
							goto l1126
						}
						position++
						goto l1125
					l1126:
						position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
						if buffer[position] != rune('M') {
							goto l1122
						}
						position++
					}
				l1125:
					{
						position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1128
						}
						position++
						goto l1127
					l1128:
						position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
						if buffer[position] != rune('I') {
							goto l1122
						}
						position++
					}
				l1127:
					{
						position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1130
						}
						position++
						goto l1129
					l1130:
						position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
						if buffer[position] != rune('L') {
							goto l1122
						}
						position++
					}
				l1129:
					{
						position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1132
						}
						position++
						goto l1131
					l1132:
						position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
						if buffer[position] != rune('L') {
							goto l1122
						}
						position++
					}
				l1131:
					{
						position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1134
						}
						position++
						goto l1133
					l1134:
						position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
						if buffer[position] != rune('I') {
							goto l1122
						}
						position++
					}
				l1133:
					{
						position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
						if buffer[position] != rune('S') {
							goto l1122
						}
						position++
					}
				l1135:
					{
						position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
						if buffer[position] != rune('E') {
							goto l1122
						}
						position++
					}
				l1137:
					{
						position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1140
						}
						position++
						goto l1139
					l1140:
						position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
						if buffer[position] != rune('C') {
							goto l1122
						}
						position++
					}
				l1139:
					{
						position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1142
						}
						position++
						goto l1141
					l1142:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
						if buffer[position] != rune('O') {
							goto l1122
						}
						position++
					}
				l1141:
					{
						position1143, tokenIndex1143, depth1143 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1144
						}
						position++
						goto l1143
					l1144:
						position, tokenIndex, depth = position1143, tokenIndex1143, depth1143
						if buffer[position] != rune('N') {
							goto l1122
						}
						position++
					}
				l1143:
					{
						position1145, tokenIndex1145, depth1145 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1146
						}
						position++
						goto l1145
					l1146:
						position, tokenIndex, depth = position1145, tokenIndex1145, depth1145
						if buffer[position] != rune('D') {
							goto l1122
						}
						position++
					}
				l1145:
					{
						position1147, tokenIndex1147, depth1147 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1148
						}
						position++
						goto l1147
					l1148:
						position, tokenIndex, depth = position1147, tokenIndex1147, depth1147
						if buffer[position] != rune('S') {
							goto l1122
						}
						position++
					}
				l1147:
					depth--
					add(rulePegText, position1124)
				}
				if !_rules[ruleAction77]() {
					goto l1122
				}
				depth--
				add(ruleMILLISECONDS, position1123)
			}
			return true
		l1122:
			position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
			return false
		},
		/* 101 StreamIdentifier <- <(<ident> Action78)> */
		func() bool {
			position1149, tokenIndex1149, depth1149 := position, tokenIndex, depth
			{
				position1150 := position
				depth++
				{
					position1151 := position
					depth++
					if !_rules[ruleident]() {
						goto l1149
					}
					depth--
					add(rulePegText, position1151)
				}
				if !_rules[ruleAction78]() {
					goto l1149
				}
				depth--
				add(ruleStreamIdentifier, position1150)
			}
			return true
		l1149:
			position, tokenIndex, depth = position1149, tokenIndex1149, depth1149
			return false
		},
		/* 102 SourceSinkType <- <(<ident> Action79)> */
		func() bool {
			position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
			{
				position1153 := position
				depth++
				{
					position1154 := position
					depth++
					if !_rules[ruleident]() {
						goto l1152
					}
					depth--
					add(rulePegText, position1154)
				}
				if !_rules[ruleAction79]() {
					goto l1152
				}
				depth--
				add(ruleSourceSinkType, position1153)
			}
			return true
		l1152:
			position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
			return false
		},
		/* 103 SourceSinkParamKey <- <(<ident> Action80)> */
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
				if !_rules[ruleAction80]() {
					goto l1155
				}
				depth--
				add(ruleSourceSinkParamKey, position1156)
			}
			return true
		l1155:
			position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
			return false
		},
		/* 104 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action81)> */
		func() bool {
			position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
			{
				position1159 := position
				depth++
				{
					position1160 := position
					depth++
					{
						position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1162
						}
						position++
						goto l1161
					l1162:
						position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
						if buffer[position] != rune('P') {
							goto l1158
						}
						position++
					}
				l1161:
					{
						position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1164
						}
						position++
						goto l1163
					l1164:
						position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
						if buffer[position] != rune('A') {
							goto l1158
						}
						position++
					}
				l1163:
					{
						position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1166
						}
						position++
						goto l1165
					l1166:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						if buffer[position] != rune('U') {
							goto l1158
						}
						position++
					}
				l1165:
					{
						position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1168
						}
						position++
						goto l1167
					l1168:
						position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
						if buffer[position] != rune('S') {
							goto l1158
						}
						position++
					}
				l1167:
					{
						position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1170
						}
						position++
						goto l1169
					l1170:
						position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
						if buffer[position] != rune('E') {
							goto l1158
						}
						position++
					}
				l1169:
					{
						position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1172
						}
						position++
						goto l1171
					l1172:
						position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
						if buffer[position] != rune('D') {
							goto l1158
						}
						position++
					}
				l1171:
					depth--
					add(rulePegText, position1160)
				}
				if !_rules[ruleAction81]() {
					goto l1158
				}
				depth--
				add(rulePaused, position1159)
			}
			return true
		l1158:
			position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
			return false
		},
		/* 105 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action82)> */
		func() bool {
			position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
			{
				position1174 := position
				depth++
				{
					position1175 := position
					depth++
					{
						position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1177
						}
						position++
						goto l1176
					l1177:
						position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
						if buffer[position] != rune('U') {
							goto l1173
						}
						position++
					}
				l1176:
					{
						position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1179
						}
						position++
						goto l1178
					l1179:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('N') {
							goto l1173
						}
						position++
					}
				l1178:
					{
						position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1181
						}
						position++
						goto l1180
					l1181:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						if buffer[position] != rune('P') {
							goto l1173
						}
						position++
					}
				l1180:
					{
						position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1183
						}
						position++
						goto l1182
					l1183:
						position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
						if buffer[position] != rune('A') {
							goto l1173
						}
						position++
					}
				l1182:
					{
						position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1185
						}
						position++
						goto l1184
					l1185:
						position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
						if buffer[position] != rune('U') {
							goto l1173
						}
						position++
					}
				l1184:
					{
						position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1187
						}
						position++
						goto l1186
					l1187:
						position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
						if buffer[position] != rune('S') {
							goto l1173
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
							goto l1173
						}
						position++
					}
				l1188:
					{
						position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1191
						}
						position++
						goto l1190
					l1191:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('D') {
							goto l1173
						}
						position++
					}
				l1190:
					depth--
					add(rulePegText, position1175)
				}
				if !_rules[ruleAction82]() {
					goto l1173
				}
				depth--
				add(ruleUnpaused, position1174)
			}
			return true
		l1173:
			position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
			return false
		},
		/* 106 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
			{
				position1193 := position
				depth++
				{
					position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1195
					}
					goto l1194
				l1195:
					position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
					if !_rules[ruleInt]() {
						goto l1196
					}
					goto l1194
				l1196:
					position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
					if !_rules[ruleFloat]() {
						goto l1197
					}
					goto l1194
				l1197:
					position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
					if !_rules[ruleString]() {
						goto l1198
					}
					goto l1194
				l1198:
					position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
					if !_rules[ruleBlob]() {
						goto l1199
					}
					goto l1194
				l1199:
					position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
					if !_rules[ruleTimestamp]() {
						goto l1200
					}
					goto l1194
				l1200:
					position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
					if !_rules[ruleArray]() {
						goto l1201
					}
					goto l1194
				l1201:
					position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
					if !_rules[ruleMap]() {
						goto l1192
					}
				}
			l1194:
				depth--
				add(ruleType, position1193)
			}
			return true
		l1192:
			position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
			return false
		},
		/* 107 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action83)> */
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
						if buffer[position] != rune('b') {
							goto l1206
						}
						position++
						goto l1205
					l1206:
						position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
						if buffer[position] != rune('B') {
							goto l1202
						}
						position++
					}
				l1205:
					{
						position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1208
						}
						position++
						goto l1207
					l1208:
						position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
						if buffer[position] != rune('O') {
							goto l1202
						}
						position++
					}
				l1207:
					{
						position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1210
						}
						position++
						goto l1209
					l1210:
						position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
						if buffer[position] != rune('O') {
							goto l1202
						}
						position++
					}
				l1209:
					{
						position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1212
						}
						position++
						goto l1211
					l1212:
						position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
						if buffer[position] != rune('L') {
							goto l1202
						}
						position++
					}
				l1211:
					depth--
					add(rulePegText, position1204)
				}
				if !_rules[ruleAction83]() {
					goto l1202
				}
				depth--
				add(ruleBool, position1203)
			}
			return true
		l1202:
			position, tokenIndex, depth = position1202, tokenIndex1202, depth1202
			return false
		},
		/* 108 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action84)> */
		func() bool {
			position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
			{
				position1214 := position
				depth++
				{
					position1215 := position
					depth++
					{
						position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1217
						}
						position++
						goto l1216
					l1217:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if buffer[position] != rune('I') {
							goto l1213
						}
						position++
					}
				l1216:
					{
						position1218, tokenIndex1218, depth1218 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1219
						}
						position++
						goto l1218
					l1219:
						position, tokenIndex, depth = position1218, tokenIndex1218, depth1218
						if buffer[position] != rune('N') {
							goto l1213
						}
						position++
					}
				l1218:
					{
						position1220, tokenIndex1220, depth1220 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1221
						}
						position++
						goto l1220
					l1221:
						position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
						if buffer[position] != rune('T') {
							goto l1213
						}
						position++
					}
				l1220:
					depth--
					add(rulePegText, position1215)
				}
				if !_rules[ruleAction84]() {
					goto l1213
				}
				depth--
				add(ruleInt, position1214)
			}
			return true
		l1213:
			position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
			return false
		},
		/* 109 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action85)> */
		func() bool {
			position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
			{
				position1223 := position
				depth++
				{
					position1224 := position
					depth++
					{
						position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1226
						}
						position++
						goto l1225
					l1226:
						position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
						if buffer[position] != rune('F') {
							goto l1222
						}
						position++
					}
				l1225:
					{
						position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1228
						}
						position++
						goto l1227
					l1228:
						position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
						if buffer[position] != rune('L') {
							goto l1222
						}
						position++
					}
				l1227:
					{
						position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1230
						}
						position++
						goto l1229
					l1230:
						position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
						if buffer[position] != rune('O') {
							goto l1222
						}
						position++
					}
				l1229:
					{
						position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1232
						}
						position++
						goto l1231
					l1232:
						position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
						if buffer[position] != rune('A') {
							goto l1222
						}
						position++
					}
				l1231:
					{
						position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1234
						}
						position++
						goto l1233
					l1234:
						position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
						if buffer[position] != rune('T') {
							goto l1222
						}
						position++
					}
				l1233:
					depth--
					add(rulePegText, position1224)
				}
				if !_rules[ruleAction85]() {
					goto l1222
				}
				depth--
				add(ruleFloat, position1223)
			}
			return true
		l1222:
			position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
			return false
		},
		/* 110 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action86)> */
		func() bool {
			position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
			{
				position1236 := position
				depth++
				{
					position1237 := position
					depth++
					{
						position1238, tokenIndex1238, depth1238 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1239
						}
						position++
						goto l1238
					l1239:
						position, tokenIndex, depth = position1238, tokenIndex1238, depth1238
						if buffer[position] != rune('S') {
							goto l1235
						}
						position++
					}
				l1238:
					{
						position1240, tokenIndex1240, depth1240 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1241
						}
						position++
						goto l1240
					l1241:
						position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
						if buffer[position] != rune('T') {
							goto l1235
						}
						position++
					}
				l1240:
					{
						position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1243
						}
						position++
						goto l1242
					l1243:
						position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
						if buffer[position] != rune('R') {
							goto l1235
						}
						position++
					}
				l1242:
					{
						position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1245
						}
						position++
						goto l1244
					l1245:
						position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
						if buffer[position] != rune('I') {
							goto l1235
						}
						position++
					}
				l1244:
					{
						position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1247
						}
						position++
						goto l1246
					l1247:
						position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
						if buffer[position] != rune('N') {
							goto l1235
						}
						position++
					}
				l1246:
					{
						position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1249
						}
						position++
						goto l1248
					l1249:
						position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
						if buffer[position] != rune('G') {
							goto l1235
						}
						position++
					}
				l1248:
					depth--
					add(rulePegText, position1237)
				}
				if !_rules[ruleAction86]() {
					goto l1235
				}
				depth--
				add(ruleString, position1236)
			}
			return true
		l1235:
			position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
			return false
		},
		/* 111 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action87)> */
		func() bool {
			position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
			{
				position1251 := position
				depth++
				{
					position1252 := position
					depth++
					{
						position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1254
						}
						position++
						goto l1253
					l1254:
						position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
						if buffer[position] != rune('B') {
							goto l1250
						}
						position++
					}
				l1253:
					{
						position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1256
						}
						position++
						goto l1255
					l1256:
						position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
						if buffer[position] != rune('L') {
							goto l1250
						}
						position++
					}
				l1255:
					{
						position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1258
						}
						position++
						goto l1257
					l1258:
						position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
						if buffer[position] != rune('O') {
							goto l1250
						}
						position++
					}
				l1257:
					{
						position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1260
						}
						position++
						goto l1259
					l1260:
						position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
						if buffer[position] != rune('B') {
							goto l1250
						}
						position++
					}
				l1259:
					depth--
					add(rulePegText, position1252)
				}
				if !_rules[ruleAction87]() {
					goto l1250
				}
				depth--
				add(ruleBlob, position1251)
			}
			return true
		l1250:
			position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
			return false
		},
		/* 112 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action88)> */
		func() bool {
			position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
			{
				position1262 := position
				depth++
				{
					position1263 := position
					depth++
					{
						position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1265
						}
						position++
						goto l1264
					l1265:
						position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
						if buffer[position] != rune('T') {
							goto l1261
						}
						position++
					}
				l1264:
					{
						position1266, tokenIndex1266, depth1266 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1267
						}
						position++
						goto l1266
					l1267:
						position, tokenIndex, depth = position1266, tokenIndex1266, depth1266
						if buffer[position] != rune('I') {
							goto l1261
						}
						position++
					}
				l1266:
					{
						position1268, tokenIndex1268, depth1268 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1269
						}
						position++
						goto l1268
					l1269:
						position, tokenIndex, depth = position1268, tokenIndex1268, depth1268
						if buffer[position] != rune('M') {
							goto l1261
						}
						position++
					}
				l1268:
					{
						position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1271
						}
						position++
						goto l1270
					l1271:
						position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
						if buffer[position] != rune('E') {
							goto l1261
						}
						position++
					}
				l1270:
					{
						position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1273
						}
						position++
						goto l1272
					l1273:
						position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
						if buffer[position] != rune('S') {
							goto l1261
						}
						position++
					}
				l1272:
					{
						position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1275
						}
						position++
						goto l1274
					l1275:
						position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
						if buffer[position] != rune('T') {
							goto l1261
						}
						position++
					}
				l1274:
					{
						position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1277
						}
						position++
						goto l1276
					l1277:
						position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
						if buffer[position] != rune('A') {
							goto l1261
						}
						position++
					}
				l1276:
					{
						position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1279
						}
						position++
						goto l1278
					l1279:
						position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
						if buffer[position] != rune('M') {
							goto l1261
						}
						position++
					}
				l1278:
					{
						position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1281
						}
						position++
						goto l1280
					l1281:
						position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
						if buffer[position] != rune('P') {
							goto l1261
						}
						position++
					}
				l1280:
					depth--
					add(rulePegText, position1263)
				}
				if !_rules[ruleAction88]() {
					goto l1261
				}
				depth--
				add(ruleTimestamp, position1262)
			}
			return true
		l1261:
			position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
			return false
		},
		/* 113 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action89)> */
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
						if buffer[position] != rune('a') {
							goto l1286
						}
						position++
						goto l1285
					l1286:
						position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
						if buffer[position] != rune('A') {
							goto l1282
						}
						position++
					}
				l1285:
					{
						position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1288
						}
						position++
						goto l1287
					l1288:
						position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
						if buffer[position] != rune('R') {
							goto l1282
						}
						position++
					}
				l1287:
					{
						position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1290
						}
						position++
						goto l1289
					l1290:
						position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
						if buffer[position] != rune('R') {
							goto l1282
						}
						position++
					}
				l1289:
					{
						position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1292
						}
						position++
						goto l1291
					l1292:
						position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
						if buffer[position] != rune('A') {
							goto l1282
						}
						position++
					}
				l1291:
					{
						position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1294
						}
						position++
						goto l1293
					l1294:
						position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
						if buffer[position] != rune('Y') {
							goto l1282
						}
						position++
					}
				l1293:
					depth--
					add(rulePegText, position1284)
				}
				if !_rules[ruleAction89]() {
					goto l1282
				}
				depth--
				add(ruleArray, position1283)
			}
			return true
		l1282:
			position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
			return false
		},
		/* 114 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action90)> */
		func() bool {
			position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
			{
				position1296 := position
				depth++
				{
					position1297 := position
					depth++
					{
						position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1299
						}
						position++
						goto l1298
					l1299:
						position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
						if buffer[position] != rune('M') {
							goto l1295
						}
						position++
					}
				l1298:
					{
						position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1301
						}
						position++
						goto l1300
					l1301:
						position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
						if buffer[position] != rune('A') {
							goto l1295
						}
						position++
					}
				l1300:
					{
						position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1303
						}
						position++
						goto l1302
					l1303:
						position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
						if buffer[position] != rune('P') {
							goto l1295
						}
						position++
					}
				l1302:
					depth--
					add(rulePegText, position1297)
				}
				if !_rules[ruleAction90]() {
					goto l1295
				}
				depth--
				add(ruleMap, position1296)
			}
			return true
		l1295:
			position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
			return false
		},
		/* 115 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action91)> */
		func() bool {
			position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
			{
				position1305 := position
				depth++
				{
					position1306 := position
					depth++
					{
						position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1308
						}
						position++
						goto l1307
					l1308:
						position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
						if buffer[position] != rune('O') {
							goto l1304
						}
						position++
					}
				l1307:
					{
						position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1310
						}
						position++
						goto l1309
					l1310:
						position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
						if buffer[position] != rune('R') {
							goto l1304
						}
						position++
					}
				l1309:
					depth--
					add(rulePegText, position1306)
				}
				if !_rules[ruleAction91]() {
					goto l1304
				}
				depth--
				add(ruleOr, position1305)
			}
			return true
		l1304:
			position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
			return false
		},
		/* 116 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action92)> */
		func() bool {
			position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
			{
				position1312 := position
				depth++
				{
					position1313 := position
					depth++
					{
						position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1315
						}
						position++
						goto l1314
					l1315:
						position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
						if buffer[position] != rune('A') {
							goto l1311
						}
						position++
					}
				l1314:
					{
						position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1317
						}
						position++
						goto l1316
					l1317:
						position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
						if buffer[position] != rune('N') {
							goto l1311
						}
						position++
					}
				l1316:
					{
						position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1319
						}
						position++
						goto l1318
					l1319:
						position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
						if buffer[position] != rune('D') {
							goto l1311
						}
						position++
					}
				l1318:
					depth--
					add(rulePegText, position1313)
				}
				if !_rules[ruleAction92]() {
					goto l1311
				}
				depth--
				add(ruleAnd, position1312)
			}
			return true
		l1311:
			position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
			return false
		},
		/* 117 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action93)> */
		func() bool {
			position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
			{
				position1321 := position
				depth++
				{
					position1322 := position
					depth++
					{
						position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1324
						}
						position++
						goto l1323
					l1324:
						position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
						if buffer[position] != rune('N') {
							goto l1320
						}
						position++
					}
				l1323:
					{
						position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1326
						}
						position++
						goto l1325
					l1326:
						position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
						if buffer[position] != rune('O') {
							goto l1320
						}
						position++
					}
				l1325:
					{
						position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1328
						}
						position++
						goto l1327
					l1328:
						position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
						if buffer[position] != rune('T') {
							goto l1320
						}
						position++
					}
				l1327:
					depth--
					add(rulePegText, position1322)
				}
				if !_rules[ruleAction93]() {
					goto l1320
				}
				depth--
				add(ruleNot, position1321)
			}
			return true
		l1320:
			position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
			return false
		},
		/* 118 Equal <- <(<'='> Action94)> */
		func() bool {
			position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
			{
				position1330 := position
				depth++
				{
					position1331 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1329
					}
					position++
					depth--
					add(rulePegText, position1331)
				}
				if !_rules[ruleAction94]() {
					goto l1329
				}
				depth--
				add(ruleEqual, position1330)
			}
			return true
		l1329:
			position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
			return false
		},
		/* 119 Less <- <(<'<'> Action95)> */
		func() bool {
			position1332, tokenIndex1332, depth1332 := position, tokenIndex, depth
			{
				position1333 := position
				depth++
				{
					position1334 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1332
					}
					position++
					depth--
					add(rulePegText, position1334)
				}
				if !_rules[ruleAction95]() {
					goto l1332
				}
				depth--
				add(ruleLess, position1333)
			}
			return true
		l1332:
			position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
			return false
		},
		/* 120 LessOrEqual <- <(<('<' '=')> Action96)> */
		func() bool {
			position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
			{
				position1336 := position
				depth++
				{
					position1337 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1335
					}
					position++
					if buffer[position] != rune('=') {
						goto l1335
					}
					position++
					depth--
					add(rulePegText, position1337)
				}
				if !_rules[ruleAction96]() {
					goto l1335
				}
				depth--
				add(ruleLessOrEqual, position1336)
			}
			return true
		l1335:
			position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
			return false
		},
		/* 121 Greater <- <(<'>'> Action97)> */
		func() bool {
			position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
			{
				position1339 := position
				depth++
				{
					position1340 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1338
					}
					position++
					depth--
					add(rulePegText, position1340)
				}
				if !_rules[ruleAction97]() {
					goto l1338
				}
				depth--
				add(ruleGreater, position1339)
			}
			return true
		l1338:
			position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
			return false
		},
		/* 122 GreaterOrEqual <- <(<('>' '=')> Action98)> */
		func() bool {
			position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
			{
				position1342 := position
				depth++
				{
					position1343 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1341
					}
					position++
					if buffer[position] != rune('=') {
						goto l1341
					}
					position++
					depth--
					add(rulePegText, position1343)
				}
				if !_rules[ruleAction98]() {
					goto l1341
				}
				depth--
				add(ruleGreaterOrEqual, position1342)
			}
			return true
		l1341:
			position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
			return false
		},
		/* 123 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action99)> */
		func() bool {
			position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
			{
				position1345 := position
				depth++
				{
					position1346 := position
					depth++
					{
						position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1348
						}
						position++
						if buffer[position] != rune('=') {
							goto l1348
						}
						position++
						goto l1347
					l1348:
						position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
						if buffer[position] != rune('<') {
							goto l1344
						}
						position++
						if buffer[position] != rune('>') {
							goto l1344
						}
						position++
					}
				l1347:
					depth--
					add(rulePegText, position1346)
				}
				if !_rules[ruleAction99]() {
					goto l1344
				}
				depth--
				add(ruleNotEqual, position1345)
			}
			return true
		l1344:
			position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
			return false
		},
		/* 124 Concat <- <(<('|' '|')> Action100)> */
		func() bool {
			position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
			{
				position1350 := position
				depth++
				{
					position1351 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1349
					}
					position++
					if buffer[position] != rune('|') {
						goto l1349
					}
					position++
					depth--
					add(rulePegText, position1351)
				}
				if !_rules[ruleAction100]() {
					goto l1349
				}
				depth--
				add(ruleConcat, position1350)
			}
			return true
		l1349:
			position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
			return false
		},
		/* 125 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action101)> */
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
						if buffer[position] != rune('i') {
							goto l1356
						}
						position++
						goto l1355
					l1356:
						position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
						if buffer[position] != rune('I') {
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
					depth--
					add(rulePegText, position1354)
				}
				if !_rules[ruleAction101]() {
					goto l1352
				}
				depth--
				add(ruleIs, position1353)
			}
			return true
		l1352:
			position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
			return false
		},
		/* 126 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action102)> */
		func() bool {
			position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
			{
				position1360 := position
				depth++
				{
					position1361 := position
					depth++
					{
						position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1363
						}
						position++
						goto l1362
					l1363:
						position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
						if buffer[position] != rune('I') {
							goto l1359
						}
						position++
					}
				l1362:
					{
						position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1365
						}
						position++
						goto l1364
					l1365:
						position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
						if buffer[position] != rune('S') {
							goto l1359
						}
						position++
					}
				l1364:
					if !_rules[rulesp]() {
						goto l1359
					}
					{
						position1366, tokenIndex1366, depth1366 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1367
						}
						position++
						goto l1366
					l1367:
						position, tokenIndex, depth = position1366, tokenIndex1366, depth1366
						if buffer[position] != rune('N') {
							goto l1359
						}
						position++
					}
				l1366:
					{
						position1368, tokenIndex1368, depth1368 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1369
						}
						position++
						goto l1368
					l1369:
						position, tokenIndex, depth = position1368, tokenIndex1368, depth1368
						if buffer[position] != rune('O') {
							goto l1359
						}
						position++
					}
				l1368:
					{
						position1370, tokenIndex1370, depth1370 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1371
						}
						position++
						goto l1370
					l1371:
						position, tokenIndex, depth = position1370, tokenIndex1370, depth1370
						if buffer[position] != rune('T') {
							goto l1359
						}
						position++
					}
				l1370:
					depth--
					add(rulePegText, position1361)
				}
				if !_rules[ruleAction102]() {
					goto l1359
				}
				depth--
				add(ruleIsNot, position1360)
			}
			return true
		l1359:
			position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
			return false
		},
		/* 127 Plus <- <(<'+'> Action103)> */
		func() bool {
			position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
			{
				position1373 := position
				depth++
				{
					position1374 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1372
					}
					position++
					depth--
					add(rulePegText, position1374)
				}
				if !_rules[ruleAction103]() {
					goto l1372
				}
				depth--
				add(rulePlus, position1373)
			}
			return true
		l1372:
			position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
			return false
		},
		/* 128 Minus <- <(<'-'> Action104)> */
		func() bool {
			position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
			{
				position1376 := position
				depth++
				{
					position1377 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1375
					}
					position++
					depth--
					add(rulePegText, position1377)
				}
				if !_rules[ruleAction104]() {
					goto l1375
				}
				depth--
				add(ruleMinus, position1376)
			}
			return true
		l1375:
			position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
			return false
		},
		/* 129 Multiply <- <(<'*'> Action105)> */
		func() bool {
			position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
			{
				position1379 := position
				depth++
				{
					position1380 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1378
					}
					position++
					depth--
					add(rulePegText, position1380)
				}
				if !_rules[ruleAction105]() {
					goto l1378
				}
				depth--
				add(ruleMultiply, position1379)
			}
			return true
		l1378:
			position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
			return false
		},
		/* 130 Divide <- <(<'/'> Action106)> */
		func() bool {
			position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
			{
				position1382 := position
				depth++
				{
					position1383 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1381
					}
					position++
					depth--
					add(rulePegText, position1383)
				}
				if !_rules[ruleAction106]() {
					goto l1381
				}
				depth--
				add(ruleDivide, position1382)
			}
			return true
		l1381:
			position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
			return false
		},
		/* 131 Modulo <- <(<'%'> Action107)> */
		func() bool {
			position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
			{
				position1385 := position
				depth++
				{
					position1386 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1384
					}
					position++
					depth--
					add(rulePegText, position1386)
				}
				if !_rules[ruleAction107]() {
					goto l1384
				}
				depth--
				add(ruleModulo, position1385)
			}
			return true
		l1384:
			position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
			return false
		},
		/* 132 UnaryMinus <- <(<'-'> Action108)> */
		func() bool {
			position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
			{
				position1388 := position
				depth++
				{
					position1389 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1387
					}
					position++
					depth--
					add(rulePegText, position1389)
				}
				if !_rules[ruleAction108]() {
					goto l1387
				}
				depth--
				add(ruleUnaryMinus, position1388)
			}
			return true
		l1387:
			position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
			return false
		},
		/* 133 Identifier <- <(<ident> Action109)> */
		func() bool {
			position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
			{
				position1391 := position
				depth++
				{
					position1392 := position
					depth++
					if !_rules[ruleident]() {
						goto l1390
					}
					depth--
					add(rulePegText, position1392)
				}
				if !_rules[ruleAction109]() {
					goto l1390
				}
				depth--
				add(ruleIdentifier, position1391)
			}
			return true
		l1390:
			position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
			return false
		},
		/* 134 TargetIdentifier <- <(<jsonPath> Action110)> */
		func() bool {
			position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
			{
				position1394 := position
				depth++
				{
					position1395 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1393
					}
					depth--
					add(rulePegText, position1395)
				}
				if !_rules[ruleAction110]() {
					goto l1393
				}
				depth--
				add(ruleTargetIdentifier, position1394)
			}
			return true
		l1393:
			position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
			return false
		},
		/* 135 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
			{
				position1397 := position
				depth++
				{
					position1398, tokenIndex1398, depth1398 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1399
					}
					position++
					goto l1398
				l1399:
					position, tokenIndex, depth = position1398, tokenIndex1398, depth1398
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1396
					}
					position++
				}
			l1398:
			l1400:
				{
					position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
					{
						position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1403
						}
						position++
						goto l1402
					l1403:
						position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1404
						}
						position++
						goto l1402
					l1404:
						position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1405
						}
						position++
						goto l1402
					l1405:
						position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
						if buffer[position] != rune('_') {
							goto l1401
						}
						position++
					}
				l1402:
					goto l1400
				l1401:
					position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
				}
				depth--
				add(ruleident, position1397)
			}
			return true
		l1396:
			position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
			return false
		},
		/* 136 jsonPath <- <(jsonPathHead jsonPathNonHead*)> */
		func() bool {
			position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
			{
				position1407 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1406
				}
			l1408:
				{
					position1409, tokenIndex1409, depth1409 := position, tokenIndex, depth
					if !_rules[rulejsonPathNonHead]() {
						goto l1409
					}
					goto l1408
				l1409:
					position, tokenIndex, depth = position1409, tokenIndex1409, depth1409
				}
				depth--
				add(rulejsonPath, position1407)
			}
			return true
		l1406:
			position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
			return false
		},
		/* 137 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
			{
				position1411 := position
				depth++
				{
					position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1413
					}
					goto l1412
				l1413:
					position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1410
					}
				}
			l1412:
				depth--
				add(rulejsonPathHead, position1411)
			}
			return true
		l1410:
			position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
			return false
		},
		/* 138 jsonPathNonHead <- <(('.' jsonMapAccessString) / jsonMapAccessBracket / jsonArrayAccess)> */
		func() bool {
			position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
			{
				position1415 := position
				depth++
				{
					position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1417
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1417
					}
					goto l1416
				l1417:
					position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1418
					}
					goto l1416
				l1418:
					position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
					if !_rules[rulejsonArrayAccess]() {
						goto l1414
					}
				}
			l1416:
				depth--
				add(rulejsonPathNonHead, position1415)
			}
			return true
		l1414:
			position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
			return false
		},
		/* 139 jsonMapAccessString <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1419, tokenIndex1419, depth1419 := position, tokenIndex, depth
			{
				position1420 := position
				depth++
				{
					position1421, tokenIndex1421, depth1421 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1422
					}
					position++
					goto l1421
				l1422:
					position, tokenIndex, depth = position1421, tokenIndex1421, depth1421
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1419
					}
					position++
				}
			l1421:
			l1423:
				{
					position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
					{
						position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1426
						}
						position++
						goto l1425
					l1426:
						position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1427
						}
						position++
						goto l1425
					l1427:
						position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1428
						}
						position++
						goto l1425
					l1428:
						position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
						if buffer[position] != rune('_') {
							goto l1424
						}
						position++
					}
				l1425:
					goto l1423
				l1424:
					position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
				}
				depth--
				add(rulejsonMapAccessString, position1420)
			}
			return true
		l1419:
			position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
			return false
		},
		/* 140 jsonMapAccessBracket <- <('[' '\'' (('\'' '\'') / (!'\'' .))* '\'' ']')> */
		func() bool {
			position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
			{
				position1430 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1429
				}
				position++
				if buffer[position] != rune('\'') {
					goto l1429
				}
				position++
			l1431:
				{
					position1432, tokenIndex1432, depth1432 := position, tokenIndex, depth
					{
						position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
						if buffer[position] != rune('\'') {
							goto l1434
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1434
						}
						position++
						goto l1433
					l1434:
						position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
						{
							position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1435
							}
							position++
							goto l1432
						l1435:
							position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
						}
						if !matchDot() {
							goto l1432
						}
					}
				l1433:
					goto l1431
				l1432:
					position, tokenIndex, depth = position1432, tokenIndex1432, depth1432
				}
				if buffer[position] != rune('\'') {
					goto l1429
				}
				position++
				if buffer[position] != rune(']') {
					goto l1429
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1430)
			}
			return true
		l1429:
			position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
			return false
		},
		/* 141 jsonArrayAccess <- <('[' [0-9]+ ']')> */
		func() bool {
			position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
			{
				position1437 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1436
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1436
				}
				position++
			l1438:
				{
					position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1439
					}
					position++
					goto l1438
				l1439:
					position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
				}
				if buffer[position] != rune(']') {
					goto l1436
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1437)
			}
			return true
		l1436:
			position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
			return false
		},
		/* 142 sp <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)*> */
		func() bool {
			{
				position1441 := position
				depth++
			l1442:
				{
					position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
					{
						position1444, tokenIndex1444, depth1444 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1445
						}
						position++
						goto l1444
					l1445:
						position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
						if buffer[position] != rune('\t') {
							goto l1446
						}
						position++
						goto l1444
					l1446:
						position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
						if buffer[position] != rune('\n') {
							goto l1447
						}
						position++
						goto l1444
					l1447:
						position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
						if buffer[position] != rune('\r') {
							goto l1448
						}
						position++
						goto l1444
					l1448:
						position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
						if !_rules[rulecomment]() {
							goto l1449
						}
						goto l1444
					l1449:
						position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
						if !_rules[rulefinalComment]() {
							goto l1443
						}
					}
				l1444:
					goto l1442
				l1443:
					position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
				}
				depth--
				add(rulesp, position1441)
			}
			return true
		},
		/* 143 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
			{
				position1451 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1450
				}
				position++
				if buffer[position] != rune('-') {
					goto l1450
				}
				position++
			l1452:
				{
					position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
					{
						position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
						{
							position1455, tokenIndex1455, depth1455 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1456
							}
							position++
							goto l1455
						l1456:
							position, tokenIndex, depth = position1455, tokenIndex1455, depth1455
							if buffer[position] != rune('\n') {
								goto l1454
							}
							position++
						}
					l1455:
						goto l1453
					l1454:
						position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
					}
					if !matchDot() {
						goto l1453
					}
					goto l1452
				l1453:
					position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
				}
				{
					position1457, tokenIndex1457, depth1457 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1458
					}
					position++
					goto l1457
				l1458:
					position, tokenIndex, depth = position1457, tokenIndex1457, depth1457
					if buffer[position] != rune('\n') {
						goto l1450
					}
					position++
				}
			l1457:
				depth--
				add(rulecomment, position1451)
			}
			return true
		l1450:
			position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
			return false
		},
		/* 144 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
			{
				position1460 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1459
				}
				position++
				if buffer[position] != rune('-') {
					goto l1459
				}
				position++
			l1461:
				{
					position1462, tokenIndex1462, depth1462 := position, tokenIndex, depth
					{
						position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
						{
							position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1465
							}
							position++
							goto l1464
						l1465:
							position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
							if buffer[position] != rune('\n') {
								goto l1463
							}
							position++
						}
					l1464:
						goto l1462
					l1463:
						position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
					}
					if !matchDot() {
						goto l1462
					}
					goto l1461
				l1462:
					position, tokenIndex, depth = position1462, tokenIndex1462, depth1462
				}
				{
					position1466, tokenIndex1466, depth1466 := position, tokenIndex, depth
					if !matchDot() {
						goto l1466
					}
					goto l1459
				l1466:
					position, tokenIndex, depth = position1466, tokenIndex1466, depth1466
				}
				depth--
				add(rulefinalComment, position1460)
			}
			return true
		l1459:
			position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
			return false
		},
		nil,
		/* 147 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 148 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 149 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 150 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 151 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 152 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 153 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 154 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 155 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 156 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 157 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 158 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 159 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 160 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 161 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 162 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 163 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 164 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 165 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 166 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 167 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 168 Action21 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 169 Action22 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 170 Action23 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 171 Action24 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 172 Action25 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 173 Action26 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 174 Action27 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 175 Action28 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 176 Action29 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 177 Action30 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 178 Action31 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 179 Action32 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 180 Action33 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 181 Action34 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 182 Action35 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 183 Action36 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 184 Action37 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 185 Action38 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 186 Action39 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 187 Action40 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 188 Action41 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 189 Action42 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 190 Action43 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 191 Action44 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 192 Action45 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 193 Action46 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 194 Action47 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 195 Action48 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 196 Action49 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 197 Action50 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 198 Action51 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 199 Action52 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 200 Action53 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 201 Action54 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 202 Action55 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 203 Action56 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 204 Action57 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 205 Action58 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 206 Action59 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 207 Action60 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 208 Action61 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 209 Action62 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 210 Action63 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 211 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 212 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 213 Action66 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 214 Action67 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 215 Action68 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 216 Action69 <- <{
		    p.PushComponent(begin, end, NewWildcard(""))
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 217 Action70 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 218 Action71 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 219 Action72 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 220 Action73 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 221 Action74 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 222 Action75 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 223 Action76 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 224 Action77 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 225 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 226 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 227 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 228 Action81 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 229 Action82 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 230 Action83 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 231 Action84 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 232 Action85 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 233 Action86 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 234 Action87 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 235 Action88 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 236 Action89 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 237 Action90 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 238 Action91 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 239 Action92 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 240 Action93 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 241 Action94 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 242 Action95 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 243 Action96 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 244 Action97 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 245 Action98 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 246 Action99 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 247 Action100 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 248 Action101 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 249 Action102 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 250 Action103 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 251 Action104 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 252 Action105 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 253 Action106 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 254 Action107 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 255 Action108 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 256 Action109 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 257 Action110 <- <{
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

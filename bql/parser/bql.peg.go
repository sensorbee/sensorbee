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
	rules  [252]func() bool
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

			p.AssembleEmitter()

		case ruleAction22:

			p.AssembleEmitterOptions(begin, end)

		case ruleAction23:

			p.AssembleEmitterLimit()

		case ruleAction24:

			p.AssembleProjections(begin, end)

		case ruleAction25:

			p.AssembleAlias()

		case ruleAction26:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction27:

			p.AssembleInterval()

		case ruleAction28:

			p.AssembleInterval()

		case ruleAction29:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction30:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction31:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction32:

			p.EnsureAliasedStreamWindow()

		case ruleAction33:

			p.AssembleAliasedStreamWindow()

		case ruleAction34:

			p.AssembleStreamWindow()

		case ruleAction35:

			p.AssembleUDSFFuncApp()

		case ruleAction36:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction37:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction38:

			p.AssembleSourceSinkParam()

		case ruleAction39:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction40:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction41:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction42:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction43:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction44:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction45:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction46:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction47:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction48:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction49:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction50:

			p.AssembleTypeCast(begin, end)

		case ruleAction51:

			p.AssembleTypeCast(begin, end)

		case ruleAction52:

			p.AssembleFuncApp()

		case ruleAction53:

			p.AssembleExpressions(begin, end)

		case ruleAction54:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction55:

			p.AssembleMap(begin, end)

		case ruleAction56:

			p.AssembleKeyValuePair()

		case ruleAction57:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction58:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction59:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction60:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction61:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction62:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction63:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction64:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction65:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction66:

			p.PushComponent(begin, end, NewWildcard(""))

		case ruleAction67:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction68:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction69:

			p.PushComponent(begin, end, Istream)

		case ruleAction70:

			p.PushComponent(begin, end, Dstream)

		case ruleAction71:

			p.PushComponent(begin, end, Rstream)

		case ruleAction72:

			p.PushComponent(begin, end, Tuples)

		case ruleAction73:

			p.PushComponent(begin, end, Seconds)

		case ruleAction74:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction75:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction76:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction77:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction78:

			p.PushComponent(begin, end, Yes)

		case ruleAction79:

			p.PushComponent(begin, end, No)

		case ruleAction80:

			p.PushComponent(begin, end, Bool)

		case ruleAction81:

			p.PushComponent(begin, end, Int)

		case ruleAction82:

			p.PushComponent(begin, end, Float)

		case ruleAction83:

			p.PushComponent(begin, end, String)

		case ruleAction84:

			p.PushComponent(begin, end, Blob)

		case ruleAction85:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction86:

			p.PushComponent(begin, end, Array)

		case ruleAction87:

			p.PushComponent(begin, end, Map)

		case ruleAction88:

			p.PushComponent(begin, end, Or)

		case ruleAction89:

			p.PushComponent(begin, end, And)

		case ruleAction90:

			p.PushComponent(begin, end, Not)

		case ruleAction91:

			p.PushComponent(begin, end, Equal)

		case ruleAction92:

			p.PushComponent(begin, end, Less)

		case ruleAction93:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction94:

			p.PushComponent(begin, end, Greater)

		case ruleAction95:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction96:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction97:

			p.PushComponent(begin, end, Concat)

		case ruleAction98:

			p.PushComponent(begin, end, Is)

		case ruleAction99:

			p.PushComponent(begin, end, IsNot)

		case ruleAction100:

			p.PushComponent(begin, end, Plus)

		case ruleAction101:

			p.PushComponent(begin, end, Minus)

		case ruleAction102:

			p.PushComponent(begin, end, Multiply)

		case ruleAction103:

			p.PushComponent(begin, end, Divide)

		case ruleAction104:

			p.PushComponent(begin, end, Modulo)

		case ruleAction105:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction106:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction107:

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
		/* 6 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt)> */
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
			position39, tokenIndex39, depth39 := position, tokenIndex, depth
			{
				position40 := position
				depth++
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l42
					}
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l43
					}
					goto l41
				l43:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if !_rules[ruleDropStreamStmt]() {
						goto l44
					}
					goto l41
				l44:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l45
					}
					goto l41
				l45:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l39
					}
				}
			l41:
				depth--
				add(ruleStreamStmt, position40)
			}
			return true
		l39:
			position, tokenIndex, depth = position39, tokenIndex39, depth39
			return false
		},
		/* 8 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action2)> */
		func() bool {
			position46, tokenIndex46, depth46 := position, tokenIndex, depth
			{
				position47 := position
				depth++
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('S') {
						goto l46
					}
					position++
				}
			l48:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('E') {
						goto l46
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('L') {
						goto l46
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
						goto l46
					}
					position++
				}
			l54:
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l57
					}
					position++
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					if buffer[position] != rune('C') {
						goto l46
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('T') {
						goto l46
					}
					position++
				}
			l58:
				if !_rules[rulesp]() {
					goto l46
				}
				if !_rules[ruleEmitter]() {
					goto l46
				}
				if !_rules[rulesp]() {
					goto l46
				}
				if !_rules[ruleProjections]() {
					goto l46
				}
				if !_rules[rulesp]() {
					goto l46
				}
				if !_rules[ruleWindowedFrom]() {
					goto l46
				}
				if !_rules[rulesp]() {
					goto l46
				}
				if !_rules[ruleFilter]() {
					goto l46
				}
				if !_rules[rulesp]() {
					goto l46
				}
				if !_rules[ruleGrouping]() {
					goto l46
				}
				if !_rules[rulesp]() {
					goto l46
				}
				if !_rules[ruleHaving]() {
					goto l46
				}
				if !_rules[rulesp]() {
					goto l46
				}
				if !_rules[ruleAction2]() {
					goto l46
				}
				depth--
				add(ruleSelectStmt, position47)
			}
			return true
		l46:
			position, tokenIndex, depth = position46, tokenIndex46, depth46
			return false
		},
		/* 9 SelectUnionStmt <- <(<(SelectStmt (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action3)> */
		func() bool {
			position60, tokenIndex60, depth60 := position, tokenIndex, depth
			{
				position61 := position
				depth++
				{
					position62 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l60
					}
					{
						position65, tokenIndex65, depth65 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l66
						}
						position++
						goto l65
					l66:
						position, tokenIndex, depth = position65, tokenIndex65, depth65
						if buffer[position] != rune('U') {
							goto l60
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
							goto l60
						}
						position++
					}
				l67:
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l70
						}
						position++
						goto l69
					l70:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
						if buffer[position] != rune('I') {
							goto l60
						}
						position++
					}
				l69:
					{
						position71, tokenIndex71, depth71 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l72
						}
						position++
						goto l71
					l72:
						position, tokenIndex, depth = position71, tokenIndex71, depth71
						if buffer[position] != rune('O') {
							goto l60
						}
						position++
					}
				l71:
					{
						position73, tokenIndex73, depth73 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l74
						}
						position++
						goto l73
					l74:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
						if buffer[position] != rune('N') {
							goto l60
						}
						position++
					}
				l73:
					if !_rules[rulesp]() {
						goto l60
					}
					{
						position75, tokenIndex75, depth75 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l76
						}
						position++
						goto l75
					l76:
						position, tokenIndex, depth = position75, tokenIndex75, depth75
						if buffer[position] != rune('A') {
							goto l60
						}
						position++
					}
				l75:
					{
						position77, tokenIndex77, depth77 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l78
						}
						position++
						goto l77
					l78:
						position, tokenIndex, depth = position77, tokenIndex77, depth77
						if buffer[position] != rune('L') {
							goto l60
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
							goto l60
						}
						position++
					}
				l79:
					if !_rules[rulesp]() {
						goto l60
					}
					if !_rules[ruleSelectStmt]() {
						goto l60
					}
				l63:
					{
						position64, tokenIndex64, depth64 := position, tokenIndex, depth
						{
							position81, tokenIndex81, depth81 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l82
							}
							position++
							goto l81
						l82:
							position, tokenIndex, depth = position81, tokenIndex81, depth81
							if buffer[position] != rune('U') {
								goto l64
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
								goto l64
							}
							position++
						}
					l83:
						{
							position85, tokenIndex85, depth85 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l86
							}
							position++
							goto l85
						l86:
							position, tokenIndex, depth = position85, tokenIndex85, depth85
							if buffer[position] != rune('I') {
								goto l64
							}
							position++
						}
					l85:
						{
							position87, tokenIndex87, depth87 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l88
							}
							position++
							goto l87
						l88:
							position, tokenIndex, depth = position87, tokenIndex87, depth87
							if buffer[position] != rune('O') {
								goto l64
							}
							position++
						}
					l87:
						{
							position89, tokenIndex89, depth89 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l90
							}
							position++
							goto l89
						l90:
							position, tokenIndex, depth = position89, tokenIndex89, depth89
							if buffer[position] != rune('N') {
								goto l64
							}
							position++
						}
					l89:
						if !_rules[rulesp]() {
							goto l64
						}
						{
							position91, tokenIndex91, depth91 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l92
							}
							position++
							goto l91
						l92:
							position, tokenIndex, depth = position91, tokenIndex91, depth91
							if buffer[position] != rune('A') {
								goto l64
							}
							position++
						}
					l91:
						{
							position93, tokenIndex93, depth93 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l94
							}
							position++
							goto l93
						l94:
							position, tokenIndex, depth = position93, tokenIndex93, depth93
							if buffer[position] != rune('L') {
								goto l64
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
								goto l64
							}
							position++
						}
					l95:
						if !_rules[rulesp]() {
							goto l64
						}
						if !_rules[ruleSelectStmt]() {
							goto l64
						}
						goto l63
					l64:
						position, tokenIndex, depth = position64, tokenIndex64, depth64
					}
					depth--
					add(rulePegText, position62)
				}
				if !_rules[ruleAction3]() {
					goto l60
				}
				depth--
				add(ruleSelectUnionStmt, position61)
			}
			return true
		l60:
			position, tokenIndex, depth = position60, tokenIndex60, depth60
			return false
		},
		/* 10 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action4)> */
		func() bool {
			position97, tokenIndex97, depth97 := position, tokenIndex, depth
			{
				position98 := position
				depth++
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('C') {
						goto l97
					}
					position++
				}
			l99:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('R') {
						goto l97
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
						goto l97
					}
					position++
				}
			l103:
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('A') {
						goto l97
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
						goto l97
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('E') {
						goto l97
					}
					position++
				}
			l109:
				if !_rules[rulesp]() {
					goto l97
				}
				{
					position111, tokenIndex111, depth111 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l112
					}
					position++
					goto l111
				l112:
					position, tokenIndex, depth = position111, tokenIndex111, depth111
					if buffer[position] != rune('S') {
						goto l97
					}
					position++
				}
			l111:
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l114
					}
					position++
					goto l113
				l114:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
					if buffer[position] != rune('T') {
						goto l97
					}
					position++
				}
			l113:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('R') {
						goto l97
					}
					position++
				}
			l115:
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('E') {
						goto l97
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('A') {
						goto l97
					}
					position++
				}
			l119:
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l122
					}
					position++
					goto l121
				l122:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
					if buffer[position] != rune('M') {
						goto l97
					}
					position++
				}
			l121:
				if !_rules[rulesp]() {
					goto l97
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l97
				}
				if !_rules[rulesp]() {
					goto l97
				}
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
						goto l97
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('S') {
						goto l97
					}
					position++
				}
			l125:
				if !_rules[rulesp]() {
					goto l97
				}
				if !_rules[ruleSelectStmt]() {
					goto l97
				}
				if !_rules[ruleAction4]() {
					goto l97
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position98)
			}
			return true
		l97:
			position, tokenIndex, depth = position97, tokenIndex97, depth97
			return false
		},
		/* 11 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action5)> */
		func() bool {
			position127, tokenIndex127, depth127 := position, tokenIndex, depth
			{
				position128 := position
				depth++
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('C') {
						goto l127
					}
					position++
				}
			l129:
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('R') {
						goto l127
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
						goto l127
					}
					position++
				}
			l133:
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('A') {
						goto l127
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
						goto l127
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('E') {
						goto l127
					}
					position++
				}
			l139:
				if !_rules[rulesp]() {
					goto l127
				}
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l142
					}
					position++
					goto l141
				l142:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					if buffer[position] != rune('S') {
						goto l127
					}
					position++
				}
			l141:
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					if buffer[position] != rune('T') {
						goto l127
					}
					position++
				}
			l143:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('R') {
						goto l127
					}
					position++
				}
			l145:
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('E') {
						goto l127
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('A') {
						goto l127
					}
					position++
				}
			l149:
				{
					position151, tokenIndex151, depth151 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l152
					}
					position++
					goto l151
				l152:
					position, tokenIndex, depth = position151, tokenIndex151, depth151
					if buffer[position] != rune('M') {
						goto l127
					}
					position++
				}
			l151:
				if !_rules[rulesp]() {
					goto l127
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l127
				}
				if !_rules[rulesp]() {
					goto l127
				}
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
						goto l127
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('S') {
						goto l127
					}
					position++
				}
			l155:
				if !_rules[rulesp]() {
					goto l127
				}
				if !_rules[ruleSelectUnionStmt]() {
					goto l127
				}
				if !_rules[ruleAction5]() {
					goto l127
				}
				depth--
				add(ruleCreateStreamAsSelectUnionStmt, position128)
			}
			return true
		l127:
			position, tokenIndex, depth = position127, tokenIndex127, depth127
			return false
		},
		/* 12 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action6)> */
		func() bool {
			position157, tokenIndex157, depth157 := position, tokenIndex, depth
			{
				position158 := position
				depth++
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l160
					}
					position++
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if buffer[position] != rune('C') {
						goto l157
					}
					position++
				}
			l159:
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l162
					}
					position++
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if buffer[position] != rune('R') {
						goto l157
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
						goto l157
					}
					position++
				}
			l163:
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('A') {
						goto l157
					}
					position++
				}
			l165:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('T') {
						goto l157
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('E') {
						goto l157
					}
					position++
				}
			l169:
				if !_rules[rulesp]() {
					goto l157
				}
				if !_rules[rulePausedOpt]() {
					goto l157
				}
				if !_rules[rulesp]() {
					goto l157
				}
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('S') {
						goto l157
					}
					position++
				}
			l171:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('O') {
						goto l157
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('U') {
						goto l157
					}
					position++
				}
			l175:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('R') {
						goto l157
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('C') {
						goto l157
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('E') {
						goto l157
					}
					position++
				}
			l181:
				if !_rules[rulesp]() {
					goto l157
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l157
				}
				if !_rules[rulesp]() {
					goto l157
				}
				{
					position183, tokenIndex183, depth183 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l184
					}
					position++
					goto l183
				l184:
					position, tokenIndex, depth = position183, tokenIndex183, depth183
					if buffer[position] != rune('T') {
						goto l157
					}
					position++
				}
			l183:
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l186
					}
					position++
					goto l185
				l186:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
					if buffer[position] != rune('Y') {
						goto l157
					}
					position++
				}
			l185:
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if buffer[position] != rune('P') {
						goto l157
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('E') {
						goto l157
					}
					position++
				}
			l189:
				if !_rules[rulesp]() {
					goto l157
				}
				if !_rules[ruleSourceSinkType]() {
					goto l157
				}
				if !_rules[rulesp]() {
					goto l157
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l157
				}
				if !_rules[ruleAction6]() {
					goto l157
				}
				depth--
				add(ruleCreateSourceStmt, position158)
			}
			return true
		l157:
			position, tokenIndex, depth = position157, tokenIndex157, depth157
			return false
		},
		/* 13 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action7)> */
		func() bool {
			position191, tokenIndex191, depth191 := position, tokenIndex, depth
			{
				position192 := position
				depth++
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('C') {
						goto l191
					}
					position++
				}
			l193:
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('R') {
						goto l191
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
						goto l191
					}
					position++
				}
			l197:
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('A') {
						goto l191
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('T') {
						goto l191
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('E') {
						goto l191
					}
					position++
				}
			l203:
				if !_rules[rulesp]() {
					goto l191
				}
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('S') {
						goto l191
					}
					position++
				}
			l205:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('I') {
						goto l191
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('N') {
						goto l191
					}
					position++
				}
			l209:
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('K') {
						goto l191
					}
					position++
				}
			l211:
				if !_rules[rulesp]() {
					goto l191
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l191
				}
				if !_rules[rulesp]() {
					goto l191
				}
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('T') {
						goto l191
					}
					position++
				}
			l213:
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l216
					}
					position++
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if buffer[position] != rune('Y') {
						goto l191
					}
					position++
				}
			l215:
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('P') {
						goto l191
					}
					position++
				}
			l217:
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('E') {
						goto l191
					}
					position++
				}
			l219:
				if !_rules[rulesp]() {
					goto l191
				}
				if !_rules[ruleSourceSinkType]() {
					goto l191
				}
				if !_rules[rulesp]() {
					goto l191
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l191
				}
				if !_rules[ruleAction7]() {
					goto l191
				}
				depth--
				add(ruleCreateSinkStmt, position192)
			}
			return true
		l191:
			position, tokenIndex, depth = position191, tokenIndex191, depth191
			return false
		},
		/* 14 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action8)> */
		func() bool {
			position221, tokenIndex221, depth221 := position, tokenIndex, depth
			{
				position222 := position
				depth++
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('C') {
						goto l221
					}
					position++
				}
			l223:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('R') {
						goto l221
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
						goto l221
					}
					position++
				}
			l227:
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('A') {
						goto l221
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
						goto l221
					}
					position++
				}
			l231:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l234
					}
					position++
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					if buffer[position] != rune('E') {
						goto l221
					}
					position++
				}
			l233:
				if !_rules[rulesp]() {
					goto l221
				}
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l236
					}
					position++
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					if buffer[position] != rune('S') {
						goto l221
					}
					position++
				}
			l235:
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l238
					}
					position++
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					if buffer[position] != rune('T') {
						goto l221
					}
					position++
				}
			l237:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l240
					}
					position++
					goto l239
				l240:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
					if buffer[position] != rune('A') {
						goto l221
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
						goto l221
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('E') {
						goto l221
					}
					position++
				}
			l243:
				if !_rules[rulesp]() {
					goto l221
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l221
				}
				if !_rules[rulesp]() {
					goto l221
				}
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
						goto l221
					}
					position++
				}
			l245:
				{
					position247, tokenIndex247, depth247 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l248
					}
					position++
					goto l247
				l248:
					position, tokenIndex, depth = position247, tokenIndex247, depth247
					if buffer[position] != rune('Y') {
						goto l221
					}
					position++
				}
			l247:
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('P') {
						goto l221
					}
					position++
				}
			l249:
				{
					position251, tokenIndex251, depth251 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l252
					}
					position++
					goto l251
				l252:
					position, tokenIndex, depth = position251, tokenIndex251, depth251
					if buffer[position] != rune('E') {
						goto l221
					}
					position++
				}
			l251:
				if !_rules[rulesp]() {
					goto l221
				}
				if !_rules[ruleSourceSinkType]() {
					goto l221
				}
				if !_rules[rulesp]() {
					goto l221
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l221
				}
				if !_rules[ruleAction8]() {
					goto l221
				}
				depth--
				add(ruleCreateStateStmt, position222)
			}
			return true
		l221:
			position, tokenIndex, depth = position221, tokenIndex221, depth221
			return false
		},
		/* 15 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position253, tokenIndex253, depth253 := position, tokenIndex, depth
			{
				position254 := position
				depth++
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l256
					}
					position++
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('U') {
						goto l253
					}
					position++
				}
			l255:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('P') {
						goto l253
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('D') {
						goto l253
					}
					position++
				}
			l259:
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('A') {
						goto l253
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
						goto l253
					}
					position++
				}
			l263:
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('E') {
						goto l253
					}
					position++
				}
			l265:
				if !_rules[rulesp]() {
					goto l253
				}
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l268
					}
					position++
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					if buffer[position] != rune('S') {
						goto l253
					}
					position++
				}
			l267:
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l270
					}
					position++
					goto l269
				l270:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
					if buffer[position] != rune('T') {
						goto l253
					}
					position++
				}
			l269:
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l272
					}
					position++
					goto l271
				l272:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
					if buffer[position] != rune('A') {
						goto l253
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
						goto l253
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('E') {
						goto l253
					}
					position++
				}
			l275:
				if !_rules[rulesp]() {
					goto l253
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l253
				}
				if !_rules[rulesp]() {
					goto l253
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l253
				}
				if !_rules[ruleAction9]() {
					goto l253
				}
				depth--
				add(ruleUpdateStateStmt, position254)
			}
			return true
		l253:
			position, tokenIndex, depth = position253, tokenIndex253, depth253
			return false
		},
		/* 16 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action10)> */
		func() bool {
			position277, tokenIndex277, depth277 := position, tokenIndex, depth
			{
				position278 := position
				depth++
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('U') {
						goto l277
					}
					position++
				}
			l279:
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('P') {
						goto l277
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('D') {
						goto l277
					}
					position++
				}
			l283:
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('A') {
						goto l277
					}
					position++
				}
			l285:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('T') {
						goto l277
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('E') {
						goto l277
					}
					position++
				}
			l289:
				if !_rules[rulesp]() {
					goto l277
				}
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('S') {
						goto l277
					}
					position++
				}
			l291:
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					if buffer[position] != rune('O') {
						goto l277
					}
					position++
				}
			l293:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('U') {
						goto l277
					}
					position++
				}
			l295:
				{
					position297, tokenIndex297, depth297 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l298
					}
					position++
					goto l297
				l298:
					position, tokenIndex, depth = position297, tokenIndex297, depth297
					if buffer[position] != rune('R') {
						goto l277
					}
					position++
				}
			l297:
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l300
					}
					position++
					goto l299
				l300:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
					if buffer[position] != rune('C') {
						goto l277
					}
					position++
				}
			l299:
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l302
					}
					position++
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if buffer[position] != rune('E') {
						goto l277
					}
					position++
				}
			l301:
				if !_rules[rulesp]() {
					goto l277
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l277
				}
				if !_rules[rulesp]() {
					goto l277
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l277
				}
				if !_rules[ruleAction10]() {
					goto l277
				}
				depth--
				add(ruleUpdateSourceStmt, position278)
			}
			return true
		l277:
			position, tokenIndex, depth = position277, tokenIndex277, depth277
			return false
		},
		/* 17 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action11)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l306
					}
					position++
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if buffer[position] != rune('U') {
						goto l303
					}
					position++
				}
			l305:
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l308
					}
					position++
					goto l307
				l308:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
					if buffer[position] != rune('P') {
						goto l303
					}
					position++
				}
			l307:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('D') {
						goto l303
					}
					position++
				}
			l309:
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l312
					}
					position++
					goto l311
				l312:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
					if buffer[position] != rune('A') {
						goto l303
					}
					position++
				}
			l311:
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('T') {
						goto l303
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('E') {
						goto l303
					}
					position++
				}
			l315:
				if !_rules[rulesp]() {
					goto l303
				}
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('S') {
						goto l303
					}
					position++
				}
			l317:
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l320
					}
					position++
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					if buffer[position] != rune('I') {
						goto l303
					}
					position++
				}
			l319:
				{
					position321, tokenIndex321, depth321 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l322
					}
					position++
					goto l321
				l322:
					position, tokenIndex, depth = position321, tokenIndex321, depth321
					if buffer[position] != rune('N') {
						goto l303
					}
					position++
				}
			l321:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l324
					}
					position++
					goto l323
				l324:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					if buffer[position] != rune('K') {
						goto l303
					}
					position++
				}
			l323:
				if !_rules[rulesp]() {
					goto l303
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l303
				}
				if !_rules[rulesp]() {
					goto l303
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l303
				}
				if !_rules[ruleAction11]() {
					goto l303
				}
				depth--
				add(ruleUpdateSinkStmt, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 18 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action12)> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l328
					}
					position++
					goto l327
				l328:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if buffer[position] != rune('I') {
						goto l325
					}
					position++
				}
			l327:
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l330
					}
					position++
					goto l329
				l330:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
					if buffer[position] != rune('N') {
						goto l325
					}
					position++
				}
			l329:
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('S') {
						goto l325
					}
					position++
				}
			l331:
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('E') {
						goto l325
					}
					position++
				}
			l333:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l336
					}
					position++
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if buffer[position] != rune('R') {
						goto l325
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
						goto l325
					}
					position++
				}
			l337:
				if !_rules[rulesp]() {
					goto l325
				}
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l340
					}
					position++
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if buffer[position] != rune('I') {
						goto l325
					}
					position++
				}
			l339:
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l342
					}
					position++
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					if buffer[position] != rune('N') {
						goto l325
					}
					position++
				}
			l341:
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l344
					}
					position++
					goto l343
				l344:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
					if buffer[position] != rune('T') {
						goto l325
					}
					position++
				}
			l343:
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l346
					}
					position++
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if buffer[position] != rune('O') {
						goto l325
					}
					position++
				}
			l345:
				if !_rules[rulesp]() {
					goto l325
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l325
				}
				if !_rules[rulesp]() {
					goto l325
				}
				if !_rules[ruleSelectStmt]() {
					goto l325
				}
				if !_rules[ruleAction12]() {
					goto l325
				}
				depth--
				add(ruleInsertIntoSelectStmt, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
			return false
		},
		/* 19 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action13)> */
		func() bool {
			position347, tokenIndex347, depth347 := position, tokenIndex, depth
			{
				position348 := position
				depth++
				{
					position349, tokenIndex349, depth349 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l350
					}
					position++
					goto l349
				l350:
					position, tokenIndex, depth = position349, tokenIndex349, depth349
					if buffer[position] != rune('I') {
						goto l347
					}
					position++
				}
			l349:
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l352
					}
					position++
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					if buffer[position] != rune('N') {
						goto l347
					}
					position++
				}
			l351:
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l354
					}
					position++
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if buffer[position] != rune('S') {
						goto l347
					}
					position++
				}
			l353:
				{
					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l356
					}
					position++
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune('E') {
						goto l347
					}
					position++
				}
			l355:
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l358
					}
					position++
					goto l357
				l358:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
					if buffer[position] != rune('R') {
						goto l347
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
						goto l347
					}
					position++
				}
			l359:
				if !_rules[rulesp]() {
					goto l347
				}
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l362
					}
					position++
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if buffer[position] != rune('I') {
						goto l347
					}
					position++
				}
			l361:
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l364
					}
					position++
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					if buffer[position] != rune('N') {
						goto l347
					}
					position++
				}
			l363:
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('T') {
						goto l347
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
						goto l347
					}
					position++
				}
			l367:
				if !_rules[rulesp]() {
					goto l347
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l347
				}
				if !_rules[rulesp]() {
					goto l347
				}
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('F') {
						goto l347
					}
					position++
				}
			l369:
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l372
					}
					position++
					goto l371
				l372:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
					if buffer[position] != rune('R') {
						goto l347
					}
					position++
				}
			l371:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('O') {
						goto l347
					}
					position++
				}
			l373:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l376
					}
					position++
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					if buffer[position] != rune('M') {
						goto l347
					}
					position++
				}
			l375:
				if !_rules[rulesp]() {
					goto l347
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l347
				}
				if !_rules[ruleAction13]() {
					goto l347
				}
				depth--
				add(ruleInsertIntoFromStmt, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 20 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position377, tokenIndex377, depth377 := position, tokenIndex, depth
			{
				position378 := position
				depth++
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('P') {
						goto l377
					}
					position++
				}
			l379:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('A') {
						goto l377
					}
					position++
				}
			l381:
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('U') {
						goto l377
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('S') {
						goto l377
					}
					position++
				}
			l385:
				{
					position387, tokenIndex387, depth387 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l388
					}
					position++
					goto l387
				l388:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
					if buffer[position] != rune('E') {
						goto l377
					}
					position++
				}
			l387:
				if !_rules[rulesp]() {
					goto l377
				}
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
						goto l377
					}
					position++
				}
			l389:
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
					if buffer[position] != rune('O') {
						goto l377
					}
					position++
				}
			l391:
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l394
					}
					position++
					goto l393
				l394:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
					if buffer[position] != rune('U') {
						goto l377
					}
					position++
				}
			l393:
				{
					position395, tokenIndex395, depth395 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l396
					}
					position++
					goto l395
				l396:
					position, tokenIndex, depth = position395, tokenIndex395, depth395
					if buffer[position] != rune('R') {
						goto l377
					}
					position++
				}
			l395:
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l398
					}
					position++
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if buffer[position] != rune('C') {
						goto l377
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
						goto l377
					}
					position++
				}
			l399:
				if !_rules[rulesp]() {
					goto l377
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l377
				}
				if !_rules[ruleAction14]() {
					goto l377
				}
				depth--
				add(rulePauseSourceStmt, position378)
			}
			return true
		l377:
			position, tokenIndex, depth = position377, tokenIndex377, depth377
			return false
		},
		/* 21 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position401, tokenIndex401, depth401 := position, tokenIndex, depth
			{
				position402 := position
				depth++
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l404
					}
					position++
					goto l403
				l404:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
					if buffer[position] != rune('R') {
						goto l401
					}
					position++
				}
			l403:
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l406
					}
					position++
					goto l405
				l406:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if buffer[position] != rune('E') {
						goto l401
					}
					position++
				}
			l405:
				{
					position407, tokenIndex407, depth407 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l408
					}
					position++
					goto l407
				l408:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
					if buffer[position] != rune('S') {
						goto l401
					}
					position++
				}
			l407:
				{
					position409, tokenIndex409, depth409 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l410
					}
					position++
					goto l409
				l410:
					position, tokenIndex, depth = position409, tokenIndex409, depth409
					if buffer[position] != rune('U') {
						goto l401
					}
					position++
				}
			l409:
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l412
					}
					position++
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if buffer[position] != rune('M') {
						goto l401
					}
					position++
				}
			l411:
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l414
					}
					position++
					goto l413
				l414:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
					if buffer[position] != rune('E') {
						goto l401
					}
					position++
				}
			l413:
				if !_rules[rulesp]() {
					goto l401
				}
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l416
					}
					position++
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if buffer[position] != rune('S') {
						goto l401
					}
					position++
				}
			l415:
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l418
					}
					position++
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if buffer[position] != rune('O') {
						goto l401
					}
					position++
				}
			l417:
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l420
					}
					position++
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if buffer[position] != rune('U') {
						goto l401
					}
					position++
				}
			l419:
				{
					position421, tokenIndex421, depth421 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l422
					}
					position++
					goto l421
				l422:
					position, tokenIndex, depth = position421, tokenIndex421, depth421
					if buffer[position] != rune('R') {
						goto l401
					}
					position++
				}
			l421:
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l424
					}
					position++
					goto l423
				l424:
					position, tokenIndex, depth = position423, tokenIndex423, depth423
					if buffer[position] != rune('C') {
						goto l401
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
						goto l401
					}
					position++
				}
			l425:
				if !_rules[rulesp]() {
					goto l401
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l401
				}
				if !_rules[ruleAction15]() {
					goto l401
				}
				depth--
				add(ruleResumeSourceStmt, position402)
			}
			return true
		l401:
			position, tokenIndex, depth = position401, tokenIndex401, depth401
			return false
		},
		/* 22 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position427, tokenIndex427, depth427 := position, tokenIndex, depth
			{
				position428 := position
				depth++
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l430
					}
					position++
					goto l429
				l430:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if buffer[position] != rune('R') {
						goto l427
					}
					position++
				}
			l429:
				{
					position431, tokenIndex431, depth431 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l432
					}
					position++
					goto l431
				l432:
					position, tokenIndex, depth = position431, tokenIndex431, depth431
					if buffer[position] != rune('E') {
						goto l427
					}
					position++
				}
			l431:
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l434
					}
					position++
					goto l433
				l434:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if buffer[position] != rune('W') {
						goto l427
					}
					position++
				}
			l433:
				{
					position435, tokenIndex435, depth435 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l436
					}
					position++
					goto l435
				l436:
					position, tokenIndex, depth = position435, tokenIndex435, depth435
					if buffer[position] != rune('I') {
						goto l427
					}
					position++
				}
			l435:
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l438
					}
					position++
					goto l437
				l438:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
					if buffer[position] != rune('N') {
						goto l427
					}
					position++
				}
			l437:
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l440
					}
					position++
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if buffer[position] != rune('D') {
						goto l427
					}
					position++
				}
			l439:
				if !_rules[rulesp]() {
					goto l427
				}
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l442
					}
					position++
					goto l441
				l442:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
					if buffer[position] != rune('S') {
						goto l427
					}
					position++
				}
			l441:
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l444
					}
					position++
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if buffer[position] != rune('O') {
						goto l427
					}
					position++
				}
			l443:
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l446
					}
					position++
					goto l445
				l446:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
					if buffer[position] != rune('U') {
						goto l427
					}
					position++
				}
			l445:
				{
					position447, tokenIndex447, depth447 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l448
					}
					position++
					goto l447
				l448:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if buffer[position] != rune('R') {
						goto l427
					}
					position++
				}
			l447:
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l450
					}
					position++
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					if buffer[position] != rune('C') {
						goto l427
					}
					position++
				}
			l449:
				{
					position451, tokenIndex451, depth451 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l452
					}
					position++
					goto l451
				l452:
					position, tokenIndex, depth = position451, tokenIndex451, depth451
					if buffer[position] != rune('E') {
						goto l427
					}
					position++
				}
			l451:
				if !_rules[rulesp]() {
					goto l427
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l427
				}
				if !_rules[ruleAction16]() {
					goto l427
				}
				depth--
				add(ruleRewindSourceStmt, position428)
			}
			return true
		l427:
			position, tokenIndex, depth = position427, tokenIndex427, depth427
			return false
		},
		/* 23 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action17)> */
		func() bool {
			position453, tokenIndex453, depth453 := position, tokenIndex, depth
			{
				position454 := position
				depth++
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l456
					}
					position++
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if buffer[position] != rune('D') {
						goto l453
					}
					position++
				}
			l455:
				{
					position457, tokenIndex457, depth457 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l458
					}
					position++
					goto l457
				l458:
					position, tokenIndex, depth = position457, tokenIndex457, depth457
					if buffer[position] != rune('R') {
						goto l453
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
						goto l453
					}
					position++
				}
			l459:
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l462
					}
					position++
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('P') {
						goto l453
					}
					position++
				}
			l461:
				if !_rules[rulesp]() {
					goto l453
				}
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('S') {
						goto l453
					}
					position++
				}
			l463:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l466
					}
					position++
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					if buffer[position] != rune('O') {
						goto l453
					}
					position++
				}
			l465:
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l468
					}
					position++
					goto l467
				l468:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
					if buffer[position] != rune('U') {
						goto l453
					}
					position++
				}
			l467:
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l470
					}
					position++
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if buffer[position] != rune('R') {
						goto l453
					}
					position++
				}
			l469:
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l472
					}
					position++
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if buffer[position] != rune('C') {
						goto l453
					}
					position++
				}
			l471:
				{
					position473, tokenIndex473, depth473 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l474
					}
					position++
					goto l473
				l474:
					position, tokenIndex, depth = position473, tokenIndex473, depth473
					if buffer[position] != rune('E') {
						goto l453
					}
					position++
				}
			l473:
				if !_rules[rulesp]() {
					goto l453
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l453
				}
				if !_rules[ruleAction17]() {
					goto l453
				}
				depth--
				add(ruleDropSourceStmt, position454)
			}
			return true
		l453:
			position, tokenIndex, depth = position453, tokenIndex453, depth453
			return false
		},
		/* 24 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action18)> */
		func() bool {
			position475, tokenIndex475, depth475 := position, tokenIndex, depth
			{
				position476 := position
				depth++
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l478
					}
					position++
					goto l477
				l478:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
					if buffer[position] != rune('D') {
						goto l475
					}
					position++
				}
			l477:
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l480
					}
					position++
					goto l479
				l480:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if buffer[position] != rune('R') {
						goto l475
					}
					position++
				}
			l479:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if buffer[position] != rune('O') {
						goto l475
					}
					position++
				}
			l481:
				{
					position483, tokenIndex483, depth483 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l484
					}
					position++
					goto l483
				l484:
					position, tokenIndex, depth = position483, tokenIndex483, depth483
					if buffer[position] != rune('P') {
						goto l475
					}
					position++
				}
			l483:
				if !_rules[rulesp]() {
					goto l475
				}
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l486
					}
					position++
					goto l485
				l486:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					if buffer[position] != rune('S') {
						goto l475
					}
					position++
				}
			l485:
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l488
					}
					position++
					goto l487
				l488:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
					if buffer[position] != rune('T') {
						goto l475
					}
					position++
				}
			l487:
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l490
					}
					position++
					goto l489
				l490:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
					if buffer[position] != rune('R') {
						goto l475
					}
					position++
				}
			l489:
				{
					position491, tokenIndex491, depth491 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l492
					}
					position++
					goto l491
				l492:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
					if buffer[position] != rune('E') {
						goto l475
					}
					position++
				}
			l491:
				{
					position493, tokenIndex493, depth493 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l494
					}
					position++
					goto l493
				l494:
					position, tokenIndex, depth = position493, tokenIndex493, depth493
					if buffer[position] != rune('A') {
						goto l475
					}
					position++
				}
			l493:
				{
					position495, tokenIndex495, depth495 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l496
					}
					position++
					goto l495
				l496:
					position, tokenIndex, depth = position495, tokenIndex495, depth495
					if buffer[position] != rune('M') {
						goto l475
					}
					position++
				}
			l495:
				if !_rules[rulesp]() {
					goto l475
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l475
				}
				if !_rules[ruleAction18]() {
					goto l475
				}
				depth--
				add(ruleDropStreamStmt, position476)
			}
			return true
		l475:
			position, tokenIndex, depth = position475, tokenIndex475, depth475
			return false
		},
		/* 25 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action19)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				{
					position499, tokenIndex499, depth499 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l500
					}
					position++
					goto l499
				l500:
					position, tokenIndex, depth = position499, tokenIndex499, depth499
					if buffer[position] != rune('D') {
						goto l497
					}
					position++
				}
			l499:
				{
					position501, tokenIndex501, depth501 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l502
					}
					position++
					goto l501
				l502:
					position, tokenIndex, depth = position501, tokenIndex501, depth501
					if buffer[position] != rune('R') {
						goto l497
					}
					position++
				}
			l501:
				{
					position503, tokenIndex503, depth503 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l504
					}
					position++
					goto l503
				l504:
					position, tokenIndex, depth = position503, tokenIndex503, depth503
					if buffer[position] != rune('O') {
						goto l497
					}
					position++
				}
			l503:
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l506
					}
					position++
					goto l505
				l506:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
					if buffer[position] != rune('P') {
						goto l497
					}
					position++
				}
			l505:
				if !_rules[rulesp]() {
					goto l497
				}
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('S') {
						goto l497
					}
					position++
				}
			l507:
				{
					position509, tokenIndex509, depth509 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l510
					}
					position++
					goto l509
				l510:
					position, tokenIndex, depth = position509, tokenIndex509, depth509
					if buffer[position] != rune('I') {
						goto l497
					}
					position++
				}
			l509:
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if buffer[position] != rune('N') {
						goto l497
					}
					position++
				}
			l511:
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l514
					}
					position++
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if buffer[position] != rune('K') {
						goto l497
					}
					position++
				}
			l513:
				if !_rules[rulesp]() {
					goto l497
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l497
				}
				if !_rules[ruleAction19]() {
					goto l497
				}
				depth--
				add(ruleDropSinkStmt, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 26 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action20)> */
		func() bool {
			position515, tokenIndex515, depth515 := position, tokenIndex, depth
			{
				position516 := position
				depth++
				{
					position517, tokenIndex517, depth517 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l518
					}
					position++
					goto l517
				l518:
					position, tokenIndex, depth = position517, tokenIndex517, depth517
					if buffer[position] != rune('D') {
						goto l515
					}
					position++
				}
			l517:
				{
					position519, tokenIndex519, depth519 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l520
					}
					position++
					goto l519
				l520:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if buffer[position] != rune('R') {
						goto l515
					}
					position++
				}
			l519:
				{
					position521, tokenIndex521, depth521 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l522
					}
					position++
					goto l521
				l522:
					position, tokenIndex, depth = position521, tokenIndex521, depth521
					if buffer[position] != rune('O') {
						goto l515
					}
					position++
				}
			l521:
				{
					position523, tokenIndex523, depth523 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l524
					}
					position++
					goto l523
				l524:
					position, tokenIndex, depth = position523, tokenIndex523, depth523
					if buffer[position] != rune('P') {
						goto l515
					}
					position++
				}
			l523:
				if !_rules[rulesp]() {
					goto l515
				}
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l526
					}
					position++
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if buffer[position] != rune('S') {
						goto l515
					}
					position++
				}
			l525:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l528
					}
					position++
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if buffer[position] != rune('T') {
						goto l515
					}
					position++
				}
			l527:
				{
					position529, tokenIndex529, depth529 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l530
					}
					position++
					goto l529
				l530:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					if buffer[position] != rune('A') {
						goto l515
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
						goto l515
					}
					position++
				}
			l531:
				{
					position533, tokenIndex533, depth533 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l534
					}
					position++
					goto l533
				l534:
					position, tokenIndex, depth = position533, tokenIndex533, depth533
					if buffer[position] != rune('E') {
						goto l515
					}
					position++
				}
			l533:
				if !_rules[rulesp]() {
					goto l515
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l515
				}
				if !_rules[ruleAction20]() {
					goto l515
				}
				depth--
				add(ruleDropStateStmt, position516)
			}
			return true
		l515:
			position, tokenIndex, depth = position515, tokenIndex515, depth515
			return false
		},
		/* 27 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) sp EmitterOptions Action21)> */
		func() bool {
			position535, tokenIndex535, depth535 := position, tokenIndex, depth
			{
				position536 := position
				depth++
				{
					position537, tokenIndex537, depth537 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l538
					}
					goto l537
				l538:
					position, tokenIndex, depth = position537, tokenIndex537, depth537
					if !_rules[ruleDSTREAM]() {
						goto l539
					}
					goto l537
				l539:
					position, tokenIndex, depth = position537, tokenIndex537, depth537
					if !_rules[ruleRSTREAM]() {
						goto l535
					}
				}
			l537:
				if !_rules[rulesp]() {
					goto l535
				}
				if !_rules[ruleEmitterOptions]() {
					goto l535
				}
				if !_rules[ruleAction21]() {
					goto l535
				}
				depth--
				add(ruleEmitter, position536)
			}
			return true
		l535:
			position, tokenIndex, depth = position535, tokenIndex535, depth535
			return false
		},
		/* 28 EmitterOptions <- <(<('[' sp EmitterLimit sp ']')?> Action22)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					{
						position543, tokenIndex543, depth543 := position, tokenIndex, depth
						if buffer[position] != rune('[') {
							goto l543
						}
						position++
						if !_rules[rulesp]() {
							goto l543
						}
						if !_rules[ruleEmitterLimit]() {
							goto l543
						}
						if !_rules[rulesp]() {
							goto l543
						}
						if buffer[position] != rune(']') {
							goto l543
						}
						position++
						goto l544
					l543:
						position, tokenIndex, depth = position543, tokenIndex543, depth543
					}
				l544:
					depth--
					add(rulePegText, position542)
				}
				if !_rules[ruleAction22]() {
					goto l540
				}
				depth--
				add(ruleEmitterOptions, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 29 EmitterLimit <- <('L' 'I' 'M' 'I' 'T' sp NumericLiteral Action23)> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				if buffer[position] != rune('L') {
					goto l545
				}
				position++
				if buffer[position] != rune('I') {
					goto l545
				}
				position++
				if buffer[position] != rune('M') {
					goto l545
				}
				position++
				if buffer[position] != rune('I') {
					goto l545
				}
				position++
				if buffer[position] != rune('T') {
					goto l545
				}
				position++
				if !_rules[rulesp]() {
					goto l545
				}
				if !_rules[ruleNumericLiteral]() {
					goto l545
				}
				if !_rules[ruleAction23]() {
					goto l545
				}
				depth--
				add(ruleEmitterLimit, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 30 Projections <- <(<(Projection sp (',' sp Projection)*)> Action24)> */
		func() bool {
			position547, tokenIndex547, depth547 := position, tokenIndex, depth
			{
				position548 := position
				depth++
				{
					position549 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l547
					}
					if !_rules[rulesp]() {
						goto l547
					}
				l550:
					{
						position551, tokenIndex551, depth551 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l551
						}
						position++
						if !_rules[rulesp]() {
							goto l551
						}
						if !_rules[ruleProjection]() {
							goto l551
						}
						goto l550
					l551:
						position, tokenIndex, depth = position551, tokenIndex551, depth551
					}
					depth--
					add(rulePegText, position549)
				}
				if !_rules[ruleAction24]() {
					goto l547
				}
				depth--
				add(ruleProjections, position548)
			}
			return true
		l547:
			position, tokenIndex, depth = position547, tokenIndex547, depth547
			return false
		},
		/* 31 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position552, tokenIndex552, depth552 := position, tokenIndex, depth
			{
				position553 := position
				depth++
				{
					position554, tokenIndex554, depth554 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l555
					}
					goto l554
				l555:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
					if !_rules[ruleExpression]() {
						goto l556
					}
					goto l554
				l556:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
					if !_rules[ruleWildcard]() {
						goto l552
					}
				}
			l554:
				depth--
				add(ruleProjection, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 32 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action25)> */
		func() bool {
			position557, tokenIndex557, depth557 := position, tokenIndex, depth
			{
				position558 := position
				depth++
				{
					position559, tokenIndex559, depth559 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l560
					}
					goto l559
				l560:
					position, tokenIndex, depth = position559, tokenIndex559, depth559
					if !_rules[ruleWildcard]() {
						goto l557
					}
				}
			l559:
				if !_rules[rulesp]() {
					goto l557
				}
				{
					position561, tokenIndex561, depth561 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l562
					}
					position++
					goto l561
				l562:
					position, tokenIndex, depth = position561, tokenIndex561, depth561
					if buffer[position] != rune('A') {
						goto l557
					}
					position++
				}
			l561:
				{
					position563, tokenIndex563, depth563 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l564
					}
					position++
					goto l563
				l564:
					position, tokenIndex, depth = position563, tokenIndex563, depth563
					if buffer[position] != rune('S') {
						goto l557
					}
					position++
				}
			l563:
				if !_rules[rulesp]() {
					goto l557
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l557
				}
				if !_rules[ruleAction25]() {
					goto l557
				}
				depth--
				add(ruleAliasExpression, position558)
			}
			return true
		l557:
			position, tokenIndex, depth = position557, tokenIndex557, depth557
			return false
		},
		/* 33 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action26)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567 := position
					depth++
					{
						position568, tokenIndex568, depth568 := position, tokenIndex, depth
						{
							position570, tokenIndex570, depth570 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l571
							}
							position++
							goto l570
						l571:
							position, tokenIndex, depth = position570, tokenIndex570, depth570
							if buffer[position] != rune('F') {
								goto l568
							}
							position++
						}
					l570:
						{
							position572, tokenIndex572, depth572 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l573
							}
							position++
							goto l572
						l573:
							position, tokenIndex, depth = position572, tokenIndex572, depth572
							if buffer[position] != rune('R') {
								goto l568
							}
							position++
						}
					l572:
						{
							position574, tokenIndex574, depth574 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l575
							}
							position++
							goto l574
						l575:
							position, tokenIndex, depth = position574, tokenIndex574, depth574
							if buffer[position] != rune('O') {
								goto l568
							}
							position++
						}
					l574:
						{
							position576, tokenIndex576, depth576 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l577
							}
							position++
							goto l576
						l577:
							position, tokenIndex, depth = position576, tokenIndex576, depth576
							if buffer[position] != rune('M') {
								goto l568
							}
							position++
						}
					l576:
						if !_rules[rulesp]() {
							goto l568
						}
						if !_rules[ruleRelations]() {
							goto l568
						}
						if !_rules[rulesp]() {
							goto l568
						}
						goto l569
					l568:
						position, tokenIndex, depth = position568, tokenIndex568, depth568
					}
				l569:
					depth--
					add(rulePegText, position567)
				}
				if !_rules[ruleAction26]() {
					goto l565
				}
				depth--
				add(ruleWindowedFrom, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 34 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position578, tokenIndex578, depth578 := position, tokenIndex, depth
			{
				position579 := position
				depth++
				{
					position580, tokenIndex580, depth580 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l581
					}
					goto l580
				l581:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
					if !_rules[ruleTuplesInterval]() {
						goto l578
					}
				}
			l580:
				depth--
				add(ruleInterval, position579)
			}
			return true
		l578:
			position, tokenIndex, depth = position578, tokenIndex578, depth578
			return false
		},
		/* 35 TimeInterval <- <(NumericLiteral sp (SECONDS / MILLISECONDS) Action27)> */
		func() bool {
			position582, tokenIndex582, depth582 := position, tokenIndex, depth
			{
				position583 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l582
				}
				if !_rules[rulesp]() {
					goto l582
				}
				{
					position584, tokenIndex584, depth584 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l585
					}
					goto l584
				l585:
					position, tokenIndex, depth = position584, tokenIndex584, depth584
					if !_rules[ruleMILLISECONDS]() {
						goto l582
					}
				}
			l584:
				if !_rules[ruleAction27]() {
					goto l582
				}
				depth--
				add(ruleTimeInterval, position583)
			}
			return true
		l582:
			position, tokenIndex, depth = position582, tokenIndex582, depth582
			return false
		},
		/* 36 TuplesInterval <- <(NumericLiteral sp TUPLES Action28)> */
		func() bool {
			position586, tokenIndex586, depth586 := position, tokenIndex, depth
			{
				position587 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l586
				}
				if !_rules[rulesp]() {
					goto l586
				}
				if !_rules[ruleTUPLES]() {
					goto l586
				}
				if !_rules[ruleAction28]() {
					goto l586
				}
				depth--
				add(ruleTuplesInterval, position587)
			}
			return true
		l586:
			position, tokenIndex, depth = position586, tokenIndex586, depth586
			return false
		},
		/* 37 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position588, tokenIndex588, depth588 := position, tokenIndex, depth
			{
				position589 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l588
				}
				if !_rules[rulesp]() {
					goto l588
				}
			l590:
				{
					position591, tokenIndex591, depth591 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l591
					}
					position++
					if !_rules[rulesp]() {
						goto l591
					}
					if !_rules[ruleRelationLike]() {
						goto l591
					}
					goto l590
				l591:
					position, tokenIndex, depth = position591, tokenIndex591, depth591
				}
				depth--
				add(ruleRelations, position589)
			}
			return true
		l588:
			position, tokenIndex, depth = position588, tokenIndex588, depth588
			return false
		},
		/* 38 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action29)> */
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
							if buffer[position] != rune('w') {
								goto l598
							}
							position++
							goto l597
						l598:
							position, tokenIndex, depth = position597, tokenIndex597, depth597
							if buffer[position] != rune('W') {
								goto l595
							}
							position++
						}
					l597:
						{
							position599, tokenIndex599, depth599 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l600
							}
							position++
							goto l599
						l600:
							position, tokenIndex, depth = position599, tokenIndex599, depth599
							if buffer[position] != rune('H') {
								goto l595
							}
							position++
						}
					l599:
						{
							position601, tokenIndex601, depth601 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l602
							}
							position++
							goto l601
						l602:
							position, tokenIndex, depth = position601, tokenIndex601, depth601
							if buffer[position] != rune('E') {
								goto l595
							}
							position++
						}
					l601:
						{
							position603, tokenIndex603, depth603 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l604
							}
							position++
							goto l603
						l604:
							position, tokenIndex, depth = position603, tokenIndex603, depth603
							if buffer[position] != rune('R') {
								goto l595
							}
							position++
						}
					l603:
						{
							position605, tokenIndex605, depth605 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l606
							}
							position++
							goto l605
						l606:
							position, tokenIndex, depth = position605, tokenIndex605, depth605
							if buffer[position] != rune('E') {
								goto l595
							}
							position++
						}
					l605:
						if !_rules[rulesp]() {
							goto l595
						}
						if !_rules[ruleExpression]() {
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
				if !_rules[ruleAction29]() {
					goto l592
				}
				depth--
				add(ruleFilter, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 39 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action30)> */
		func() bool {
			position607, tokenIndex607, depth607 := position, tokenIndex, depth
			{
				position608 := position
				depth++
				{
					position609 := position
					depth++
					{
						position610, tokenIndex610, depth610 := position, tokenIndex, depth
						{
							position612, tokenIndex612, depth612 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l613
							}
							position++
							goto l612
						l613:
							position, tokenIndex, depth = position612, tokenIndex612, depth612
							if buffer[position] != rune('G') {
								goto l610
							}
							position++
						}
					l612:
						{
							position614, tokenIndex614, depth614 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l615
							}
							position++
							goto l614
						l615:
							position, tokenIndex, depth = position614, tokenIndex614, depth614
							if buffer[position] != rune('R') {
								goto l610
							}
							position++
						}
					l614:
						{
							position616, tokenIndex616, depth616 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l617
							}
							position++
							goto l616
						l617:
							position, tokenIndex, depth = position616, tokenIndex616, depth616
							if buffer[position] != rune('O') {
								goto l610
							}
							position++
						}
					l616:
						{
							position618, tokenIndex618, depth618 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l619
							}
							position++
							goto l618
						l619:
							position, tokenIndex, depth = position618, tokenIndex618, depth618
							if buffer[position] != rune('U') {
								goto l610
							}
							position++
						}
					l618:
						{
							position620, tokenIndex620, depth620 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l621
							}
							position++
							goto l620
						l621:
							position, tokenIndex, depth = position620, tokenIndex620, depth620
							if buffer[position] != rune('P') {
								goto l610
							}
							position++
						}
					l620:
						if !_rules[rulesp]() {
							goto l610
						}
						{
							position622, tokenIndex622, depth622 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l623
							}
							position++
							goto l622
						l623:
							position, tokenIndex, depth = position622, tokenIndex622, depth622
							if buffer[position] != rune('B') {
								goto l610
							}
							position++
						}
					l622:
						{
							position624, tokenIndex624, depth624 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l625
							}
							position++
							goto l624
						l625:
							position, tokenIndex, depth = position624, tokenIndex624, depth624
							if buffer[position] != rune('Y') {
								goto l610
							}
							position++
						}
					l624:
						if !_rules[rulesp]() {
							goto l610
						}
						if !_rules[ruleGroupList]() {
							goto l610
						}
						goto l611
					l610:
						position, tokenIndex, depth = position610, tokenIndex610, depth610
					}
				l611:
					depth--
					add(rulePegText, position609)
				}
				if !_rules[ruleAction30]() {
					goto l607
				}
				depth--
				add(ruleGrouping, position608)
			}
			return true
		l607:
			position, tokenIndex, depth = position607, tokenIndex607, depth607
			return false
		},
		/* 40 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l626
				}
				if !_rules[rulesp]() {
					goto l626
				}
			l628:
				{
					position629, tokenIndex629, depth629 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l629
					}
					position++
					if !_rules[rulesp]() {
						goto l629
					}
					if !_rules[ruleExpression]() {
						goto l629
					}
					goto l628
				l629:
					position, tokenIndex, depth = position629, tokenIndex629, depth629
				}
				depth--
				add(ruleGroupList, position627)
			}
			return true
		l626:
			position, tokenIndex, depth = position626, tokenIndex626, depth626
			return false
		},
		/* 41 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action31)> */
		func() bool {
			position630, tokenIndex630, depth630 := position, tokenIndex, depth
			{
				position631 := position
				depth++
				{
					position632 := position
					depth++
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						{
							position635, tokenIndex635, depth635 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l636
							}
							position++
							goto l635
						l636:
							position, tokenIndex, depth = position635, tokenIndex635, depth635
							if buffer[position] != rune('H') {
								goto l633
							}
							position++
						}
					l635:
						{
							position637, tokenIndex637, depth637 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l638
							}
							position++
							goto l637
						l638:
							position, tokenIndex, depth = position637, tokenIndex637, depth637
							if buffer[position] != rune('A') {
								goto l633
							}
							position++
						}
					l637:
						{
							position639, tokenIndex639, depth639 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l640
							}
							position++
							goto l639
						l640:
							position, tokenIndex, depth = position639, tokenIndex639, depth639
							if buffer[position] != rune('V') {
								goto l633
							}
							position++
						}
					l639:
						{
							position641, tokenIndex641, depth641 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l642
							}
							position++
							goto l641
						l642:
							position, tokenIndex, depth = position641, tokenIndex641, depth641
							if buffer[position] != rune('I') {
								goto l633
							}
							position++
						}
					l641:
						{
							position643, tokenIndex643, depth643 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l644
							}
							position++
							goto l643
						l644:
							position, tokenIndex, depth = position643, tokenIndex643, depth643
							if buffer[position] != rune('N') {
								goto l633
							}
							position++
						}
					l643:
						{
							position645, tokenIndex645, depth645 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l646
							}
							position++
							goto l645
						l646:
							position, tokenIndex, depth = position645, tokenIndex645, depth645
							if buffer[position] != rune('G') {
								goto l633
							}
							position++
						}
					l645:
						if !_rules[rulesp]() {
							goto l633
						}
						if !_rules[ruleExpression]() {
							goto l633
						}
						goto l634
					l633:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
					}
				l634:
					depth--
					add(rulePegText, position632)
				}
				if !_rules[ruleAction31]() {
					goto l630
				}
				depth--
				add(ruleHaving, position631)
			}
			return true
		l630:
			position, tokenIndex, depth = position630, tokenIndex630, depth630
			return false
		},
		/* 42 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action32))> */
		func() bool {
			position647, tokenIndex647, depth647 := position, tokenIndex, depth
			{
				position648 := position
				depth++
				{
					position649, tokenIndex649, depth649 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l650
					}
					goto l649
				l650:
					position, tokenIndex, depth = position649, tokenIndex649, depth649
					if !_rules[ruleStreamWindow]() {
						goto l647
					}
					if !_rules[ruleAction32]() {
						goto l647
					}
				}
			l649:
				depth--
				add(ruleRelationLike, position648)
			}
			return true
		l647:
			position, tokenIndex, depth = position647, tokenIndex647, depth647
			return false
		},
		/* 43 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action33)> */
		func() bool {
			position651, tokenIndex651, depth651 := position, tokenIndex, depth
			{
				position652 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l651
				}
				if !_rules[rulesp]() {
					goto l651
				}
				{
					position653, tokenIndex653, depth653 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l654
					}
					position++
					goto l653
				l654:
					position, tokenIndex, depth = position653, tokenIndex653, depth653
					if buffer[position] != rune('A') {
						goto l651
					}
					position++
				}
			l653:
				{
					position655, tokenIndex655, depth655 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l656
					}
					position++
					goto l655
				l656:
					position, tokenIndex, depth = position655, tokenIndex655, depth655
					if buffer[position] != rune('S') {
						goto l651
					}
					position++
				}
			l655:
				if !_rules[rulesp]() {
					goto l651
				}
				if !_rules[ruleIdentifier]() {
					goto l651
				}
				if !_rules[ruleAction33]() {
					goto l651
				}
				depth--
				add(ruleAliasedStreamWindow, position652)
			}
			return true
		l651:
			position, tokenIndex, depth = position651, tokenIndex651, depth651
			return false
		},
		/* 44 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action34)> */
		func() bool {
			position657, tokenIndex657, depth657 := position, tokenIndex, depth
			{
				position658 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l657
				}
				if !_rules[rulesp]() {
					goto l657
				}
				if buffer[position] != rune('[') {
					goto l657
				}
				position++
				if !_rules[rulesp]() {
					goto l657
				}
				{
					position659, tokenIndex659, depth659 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l660
					}
					position++
					goto l659
				l660:
					position, tokenIndex, depth = position659, tokenIndex659, depth659
					if buffer[position] != rune('R') {
						goto l657
					}
					position++
				}
			l659:
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l662
					}
					position++
					goto l661
				l662:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
					if buffer[position] != rune('A') {
						goto l657
					}
					position++
				}
			l661:
				{
					position663, tokenIndex663, depth663 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l664
					}
					position++
					goto l663
				l664:
					position, tokenIndex, depth = position663, tokenIndex663, depth663
					if buffer[position] != rune('N') {
						goto l657
					}
					position++
				}
			l663:
				{
					position665, tokenIndex665, depth665 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l666
					}
					position++
					goto l665
				l666:
					position, tokenIndex, depth = position665, tokenIndex665, depth665
					if buffer[position] != rune('G') {
						goto l657
					}
					position++
				}
			l665:
				{
					position667, tokenIndex667, depth667 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l668
					}
					position++
					goto l667
				l668:
					position, tokenIndex, depth = position667, tokenIndex667, depth667
					if buffer[position] != rune('E') {
						goto l657
					}
					position++
				}
			l667:
				if !_rules[rulesp]() {
					goto l657
				}
				if !_rules[ruleInterval]() {
					goto l657
				}
				if !_rules[rulesp]() {
					goto l657
				}
				if buffer[position] != rune(']') {
					goto l657
				}
				position++
				if !_rules[ruleAction34]() {
					goto l657
				}
				depth--
				add(ruleStreamWindow, position658)
			}
			return true
		l657:
			position, tokenIndex, depth = position657, tokenIndex657, depth657
			return false
		},
		/* 45 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position669, tokenIndex669, depth669 := position, tokenIndex, depth
			{
				position670 := position
				depth++
				{
					position671, tokenIndex671, depth671 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l672
					}
					goto l671
				l672:
					position, tokenIndex, depth = position671, tokenIndex671, depth671
					if !_rules[ruleStream]() {
						goto l669
					}
				}
			l671:
				depth--
				add(ruleStreamLike, position670)
			}
			return true
		l669:
			position, tokenIndex, depth = position669, tokenIndex669, depth669
			return false
		},
		/* 46 UDSFFuncApp <- <(FuncApp Action35)> */
		func() bool {
			position673, tokenIndex673, depth673 := position, tokenIndex, depth
			{
				position674 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l673
				}
				if !_rules[ruleAction35]() {
					goto l673
				}
				depth--
				add(ruleUDSFFuncApp, position674)
			}
			return true
		l673:
			position, tokenIndex, depth = position673, tokenIndex673, depth673
			return false
		},
		/* 47 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action36)> */
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
						{
							position680, tokenIndex680, depth680 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l681
							}
							position++
							goto l680
						l681:
							position, tokenIndex, depth = position680, tokenIndex680, depth680
							if buffer[position] != rune('W') {
								goto l678
							}
							position++
						}
					l680:
						{
							position682, tokenIndex682, depth682 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l683
							}
							position++
							goto l682
						l683:
							position, tokenIndex, depth = position682, tokenIndex682, depth682
							if buffer[position] != rune('I') {
								goto l678
							}
							position++
						}
					l682:
						{
							position684, tokenIndex684, depth684 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l685
							}
							position++
							goto l684
						l685:
							position, tokenIndex, depth = position684, tokenIndex684, depth684
							if buffer[position] != rune('T') {
								goto l678
							}
							position++
						}
					l684:
						{
							position686, tokenIndex686, depth686 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l687
							}
							position++
							goto l686
						l687:
							position, tokenIndex, depth = position686, tokenIndex686, depth686
							if buffer[position] != rune('H') {
								goto l678
							}
							position++
						}
					l686:
						if !_rules[rulesp]() {
							goto l678
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l678
						}
						if !_rules[rulesp]() {
							goto l678
						}
					l688:
						{
							position689, tokenIndex689, depth689 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l689
							}
							position++
							if !_rules[rulesp]() {
								goto l689
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l689
							}
							goto l688
						l689:
							position, tokenIndex, depth = position689, tokenIndex689, depth689
						}
						goto l679
					l678:
						position, tokenIndex, depth = position678, tokenIndex678, depth678
					}
				l679:
					depth--
					add(rulePegText, position677)
				}
				if !_rules[ruleAction36]() {
					goto l675
				}
				depth--
				add(ruleSourceSinkSpecs, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 48 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action37)> */
		func() bool {
			position690, tokenIndex690, depth690 := position, tokenIndex, depth
			{
				position691 := position
				depth++
				{
					position692 := position
					depth++
					{
						position693, tokenIndex693, depth693 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l694
						}
						position++
						goto l693
					l694:
						position, tokenIndex, depth = position693, tokenIndex693, depth693
						if buffer[position] != rune('S') {
							goto l690
						}
						position++
					}
				l693:
					{
						position695, tokenIndex695, depth695 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
						if buffer[position] != rune('E') {
							goto l690
						}
						position++
					}
				l695:
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('T') {
							goto l690
						}
						position++
					}
				l697:
					if !_rules[rulesp]() {
						goto l690
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l690
					}
					if !_rules[rulesp]() {
						goto l690
					}
				l699:
					{
						position700, tokenIndex700, depth700 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l700
						}
						position++
						if !_rules[rulesp]() {
							goto l700
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l700
						}
						goto l699
					l700:
						position, tokenIndex, depth = position700, tokenIndex700, depth700
					}
					depth--
					add(rulePegText, position692)
				}
				if !_rules[ruleAction37]() {
					goto l690
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position691)
			}
			return true
		l690:
			position, tokenIndex, depth = position690, tokenIndex690, depth690
			return false
		},
		/* 49 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action38)> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l701
				}
				if buffer[position] != rune('=') {
					goto l701
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l701
				}
				if !_rules[ruleAction38]() {
					goto l701
				}
				depth--
				add(ruleSourceSinkParam, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 50 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr)> */
		func() bool {
			position703, tokenIndex703, depth703 := position, tokenIndex, depth
			{
				position704 := position
				depth++
				{
					position705, tokenIndex705, depth705 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l706
					}
					goto l705
				l706:
					position, tokenIndex, depth = position705, tokenIndex705, depth705
					if !_rules[ruleParamArrayExpr]() {
						goto l703
					}
				}
			l705:
				depth--
				add(ruleSourceSinkParamVal, position704)
			}
			return true
		l703:
			position, tokenIndex, depth = position703, tokenIndex703, depth703
			return false
		},
		/* 51 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position707, tokenIndex707, depth707 := position, tokenIndex, depth
			{
				position708 := position
				depth++
				{
					position709, tokenIndex709, depth709 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l710
					}
					goto l709
				l710:
					position, tokenIndex, depth = position709, tokenIndex709, depth709
					if !_rules[ruleLiteral]() {
						goto l707
					}
				}
			l709:
				depth--
				add(ruleParamLiteral, position708)
			}
			return true
		l707:
			position, tokenIndex, depth = position707, tokenIndex707, depth707
			return false
		},
		/* 52 ParamArrayExpr <- <(<('[' sp (ParamLiteral (',' sp ParamLiteral)*)? sp ','? sp ']')> Action39)> */
		func() bool {
			position711, tokenIndex711, depth711 := position, tokenIndex, depth
			{
				position712 := position
				depth++
				{
					position713 := position
					depth++
					if buffer[position] != rune('[') {
						goto l711
					}
					position++
					if !_rules[rulesp]() {
						goto l711
					}
					{
						position714, tokenIndex714, depth714 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
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
							if !_rules[ruleParamLiteral]() {
								goto l717
							}
							goto l716
						l717:
							position, tokenIndex, depth = position717, tokenIndex717, depth717
						}
						goto l715
					l714:
						position, tokenIndex, depth = position714, tokenIndex714, depth714
					}
				l715:
					if !_rules[rulesp]() {
						goto l711
					}
					{
						position718, tokenIndex718, depth718 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l718
						}
						position++
						goto l719
					l718:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
					}
				l719:
					if !_rules[rulesp]() {
						goto l711
					}
					if buffer[position] != rune(']') {
						goto l711
					}
					position++
					depth--
					add(rulePegText, position713)
				}
				if !_rules[ruleAction39]() {
					goto l711
				}
				depth--
				add(ruleParamArrayExpr, position712)
			}
			return true
		l711:
			position, tokenIndex, depth = position711, tokenIndex711, depth711
			return false
		},
		/* 53 PausedOpt <- <(<(Paused / Unpaused)?> Action40)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722 := position
					depth++
					{
						position723, tokenIndex723, depth723 := position, tokenIndex, depth
						{
							position725, tokenIndex725, depth725 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l726
							}
							goto l725
						l726:
							position, tokenIndex, depth = position725, tokenIndex725, depth725
							if !_rules[ruleUnpaused]() {
								goto l723
							}
						}
					l725:
						goto l724
					l723:
						position, tokenIndex, depth = position723, tokenIndex723, depth723
					}
				l724:
					depth--
					add(rulePegText, position722)
				}
				if !_rules[ruleAction40]() {
					goto l720
				}
				depth--
				add(rulePausedOpt, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 54 Expression <- <orExpr> */
		func() bool {
			position727, tokenIndex727, depth727 := position, tokenIndex, depth
			{
				position728 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l727
				}
				depth--
				add(ruleExpression, position728)
			}
			return true
		l727:
			position, tokenIndex, depth = position727, tokenIndex727, depth727
			return false
		},
		/* 55 orExpr <- <(<(andExpr sp (Or sp andExpr)*)> Action41)> */
		func() bool {
			position729, tokenIndex729, depth729 := position, tokenIndex, depth
			{
				position730 := position
				depth++
				{
					position731 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l729
					}
					if !_rules[rulesp]() {
						goto l729
					}
				l732:
					{
						position733, tokenIndex733, depth733 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l733
						}
						if !_rules[rulesp]() {
							goto l733
						}
						if !_rules[ruleandExpr]() {
							goto l733
						}
						goto l732
					l733:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
					}
					depth--
					add(rulePegText, position731)
				}
				if !_rules[ruleAction41]() {
					goto l729
				}
				depth--
				add(ruleorExpr, position730)
			}
			return true
		l729:
			position, tokenIndex, depth = position729, tokenIndex729, depth729
			return false
		},
		/* 56 andExpr <- <(<(notExpr sp (And sp notExpr)*)> Action42)> */
		func() bool {
			position734, tokenIndex734, depth734 := position, tokenIndex, depth
			{
				position735 := position
				depth++
				{
					position736 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l734
					}
					if !_rules[rulesp]() {
						goto l734
					}
				l737:
					{
						position738, tokenIndex738, depth738 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l738
						}
						if !_rules[rulesp]() {
							goto l738
						}
						if !_rules[rulenotExpr]() {
							goto l738
						}
						goto l737
					l738:
						position, tokenIndex, depth = position738, tokenIndex738, depth738
					}
					depth--
					add(rulePegText, position736)
				}
				if !_rules[ruleAction42]() {
					goto l734
				}
				depth--
				add(ruleandExpr, position735)
			}
			return true
		l734:
			position, tokenIndex, depth = position734, tokenIndex734, depth734
			return false
		},
		/* 57 notExpr <- <(<((Not sp)? comparisonExpr)> Action43)> */
		func() bool {
			position739, tokenIndex739, depth739 := position, tokenIndex, depth
			{
				position740 := position
				depth++
				{
					position741 := position
					depth++
					{
						position742, tokenIndex742, depth742 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l742
						}
						if !_rules[rulesp]() {
							goto l742
						}
						goto l743
					l742:
						position, tokenIndex, depth = position742, tokenIndex742, depth742
					}
				l743:
					if !_rules[rulecomparisonExpr]() {
						goto l739
					}
					depth--
					add(rulePegText, position741)
				}
				if !_rules[ruleAction43]() {
					goto l739
				}
				depth--
				add(rulenotExpr, position740)
			}
			return true
		l739:
			position, tokenIndex, depth = position739, tokenIndex739, depth739
			return false
		},
		/* 58 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action44)> */
		func() bool {
			position744, tokenIndex744, depth744 := position, tokenIndex, depth
			{
				position745 := position
				depth++
				{
					position746 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l744
					}
					if !_rules[rulesp]() {
						goto l744
					}
					{
						position747, tokenIndex747, depth747 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l747
						}
						if !_rules[rulesp]() {
							goto l747
						}
						if !_rules[ruleotherOpExpr]() {
							goto l747
						}
						goto l748
					l747:
						position, tokenIndex, depth = position747, tokenIndex747, depth747
					}
				l748:
					depth--
					add(rulePegText, position746)
				}
				if !_rules[ruleAction44]() {
					goto l744
				}
				depth--
				add(rulecomparisonExpr, position745)
			}
			return true
		l744:
			position, tokenIndex, depth = position744, tokenIndex744, depth744
			return false
		},
		/* 59 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action45)> */
		func() bool {
			position749, tokenIndex749, depth749 := position, tokenIndex, depth
			{
				position750 := position
				depth++
				{
					position751 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l749
					}
					if !_rules[rulesp]() {
						goto l749
					}
				l752:
					{
						position753, tokenIndex753, depth753 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l753
						}
						if !_rules[rulesp]() {
							goto l753
						}
						if !_rules[ruleisExpr]() {
							goto l753
						}
						if !_rules[rulesp]() {
							goto l753
						}
						goto l752
					l753:
						position, tokenIndex, depth = position753, tokenIndex753, depth753
					}
					depth--
					add(rulePegText, position751)
				}
				if !_rules[ruleAction45]() {
					goto l749
				}
				depth--
				add(ruleotherOpExpr, position750)
			}
			return true
		l749:
			position, tokenIndex, depth = position749, tokenIndex749, depth749
			return false
		},
		/* 60 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action46)> */
		func() bool {
			position754, tokenIndex754, depth754 := position, tokenIndex, depth
			{
				position755 := position
				depth++
				{
					position756 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l754
					}
					if !_rules[rulesp]() {
						goto l754
					}
					{
						position757, tokenIndex757, depth757 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l757
						}
						if !_rules[rulesp]() {
							goto l757
						}
						if !_rules[ruleNullLiteral]() {
							goto l757
						}
						goto l758
					l757:
						position, tokenIndex, depth = position757, tokenIndex757, depth757
					}
				l758:
					depth--
					add(rulePegText, position756)
				}
				if !_rules[ruleAction46]() {
					goto l754
				}
				depth--
				add(ruleisExpr, position755)
			}
			return true
		l754:
			position, tokenIndex, depth = position754, tokenIndex754, depth754
			return false
		},
		/* 61 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action47)> */
		func() bool {
			position759, tokenIndex759, depth759 := position, tokenIndex, depth
			{
				position760 := position
				depth++
				{
					position761 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l759
					}
					if !_rules[rulesp]() {
						goto l759
					}
				l762:
					{
						position763, tokenIndex763, depth763 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l763
						}
						if !_rules[rulesp]() {
							goto l763
						}
						if !_rules[ruleproductExpr]() {
							goto l763
						}
						if !_rules[rulesp]() {
							goto l763
						}
						goto l762
					l763:
						position, tokenIndex, depth = position763, tokenIndex763, depth763
					}
					depth--
					add(rulePegText, position761)
				}
				if !_rules[ruleAction47]() {
					goto l759
				}
				depth--
				add(ruletermExpr, position760)
			}
			return true
		l759:
			position, tokenIndex, depth = position759, tokenIndex759, depth759
			return false
		},
		/* 62 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action48)> */
		func() bool {
			position764, tokenIndex764, depth764 := position, tokenIndex, depth
			{
				position765 := position
				depth++
				{
					position766 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l764
					}
					if !_rules[rulesp]() {
						goto l764
					}
				l767:
					{
						position768, tokenIndex768, depth768 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l768
						}
						if !_rules[rulesp]() {
							goto l768
						}
						if !_rules[ruleminusExpr]() {
							goto l768
						}
						if !_rules[rulesp]() {
							goto l768
						}
						goto l767
					l768:
						position, tokenIndex, depth = position768, tokenIndex768, depth768
					}
					depth--
					add(rulePegText, position766)
				}
				if !_rules[ruleAction48]() {
					goto l764
				}
				depth--
				add(ruleproductExpr, position765)
			}
			return true
		l764:
			position, tokenIndex, depth = position764, tokenIndex764, depth764
			return false
		},
		/* 63 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action49)> */
		func() bool {
			position769, tokenIndex769, depth769 := position, tokenIndex, depth
			{
				position770 := position
				depth++
				{
					position771 := position
					depth++
					{
						position772, tokenIndex772, depth772 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l772
						}
						if !_rules[rulesp]() {
							goto l772
						}
						goto l773
					l772:
						position, tokenIndex, depth = position772, tokenIndex772, depth772
					}
				l773:
					if !_rules[rulecastExpr]() {
						goto l769
					}
					depth--
					add(rulePegText, position771)
				}
				if !_rules[ruleAction49]() {
					goto l769
				}
				depth--
				add(ruleminusExpr, position770)
			}
			return true
		l769:
			position, tokenIndex, depth = position769, tokenIndex769, depth769
			return false
		},
		/* 64 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action50)> */
		func() bool {
			position774, tokenIndex774, depth774 := position, tokenIndex, depth
			{
				position775 := position
				depth++
				{
					position776 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l774
					}
					{
						position777, tokenIndex777, depth777 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l777
						}
						if buffer[position] != rune(':') {
							goto l777
						}
						position++
						if buffer[position] != rune(':') {
							goto l777
						}
						position++
						if !_rules[rulesp]() {
							goto l777
						}
						if !_rules[ruleType]() {
							goto l777
						}
						goto l778
					l777:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
					}
				l778:
					depth--
					add(rulePegText, position776)
				}
				if !_rules[ruleAction50]() {
					goto l774
				}
				depth--
				add(rulecastExpr, position775)
			}
			return true
		l774:
			position, tokenIndex, depth = position774, tokenIndex774, depth774
			return false
		},
		/* 65 baseExpr <- <(('(' sp Expression sp ')') / MapExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position779, tokenIndex779, depth779 := position, tokenIndex, depth
			{
				position780 := position
				depth++
				{
					position781, tokenIndex781, depth781 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l782
					}
					position++
					if !_rules[rulesp]() {
						goto l782
					}
					if !_rules[ruleExpression]() {
						goto l782
					}
					if !_rules[rulesp]() {
						goto l782
					}
					if buffer[position] != rune(')') {
						goto l782
					}
					position++
					goto l781
				l782:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleMapExpr]() {
						goto l783
					}
					goto l781
				l783:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleBooleanLiteral]() {
						goto l784
					}
					goto l781
				l784:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleNullLiteral]() {
						goto l785
					}
					goto l781
				l785:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleRowMeta]() {
						goto l786
					}
					goto l781
				l786:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleFuncTypeCast]() {
						goto l787
					}
					goto l781
				l787:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleFuncApp]() {
						goto l788
					}
					goto l781
				l788:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleRowValue]() {
						goto l789
					}
					goto l781
				l789:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleArrayExpr]() {
						goto l790
					}
					goto l781
				l790:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if !_rules[ruleLiteral]() {
						goto l779
					}
				}
			l781:
				depth--
				add(rulebaseExpr, position780)
			}
			return true
		l779:
			position, tokenIndex, depth = position779, tokenIndex779, depth779
			return false
		},
		/* 66 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action51)> */
		func() bool {
			position791, tokenIndex791, depth791 := position, tokenIndex, depth
			{
				position792 := position
				depth++
				{
					position793 := position
					depth++
					{
						position794, tokenIndex794, depth794 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l795
						}
						position++
						goto l794
					l795:
						position, tokenIndex, depth = position794, tokenIndex794, depth794
						if buffer[position] != rune('C') {
							goto l791
						}
						position++
					}
				l794:
					{
						position796, tokenIndex796, depth796 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l797
						}
						position++
						goto l796
					l797:
						position, tokenIndex, depth = position796, tokenIndex796, depth796
						if buffer[position] != rune('A') {
							goto l791
						}
						position++
					}
				l796:
					{
						position798, tokenIndex798, depth798 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l799
						}
						position++
						goto l798
					l799:
						position, tokenIndex, depth = position798, tokenIndex798, depth798
						if buffer[position] != rune('S') {
							goto l791
						}
						position++
					}
				l798:
					{
						position800, tokenIndex800, depth800 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l801
						}
						position++
						goto l800
					l801:
						position, tokenIndex, depth = position800, tokenIndex800, depth800
						if buffer[position] != rune('T') {
							goto l791
						}
						position++
					}
				l800:
					if !_rules[rulesp]() {
						goto l791
					}
					if buffer[position] != rune('(') {
						goto l791
					}
					position++
					if !_rules[rulesp]() {
						goto l791
					}
					if !_rules[ruleExpression]() {
						goto l791
					}
					if !_rules[rulesp]() {
						goto l791
					}
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
							goto l791
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
							goto l791
						}
						position++
					}
				l804:
					if !_rules[rulesp]() {
						goto l791
					}
					if !_rules[ruleType]() {
						goto l791
					}
					if !_rules[rulesp]() {
						goto l791
					}
					if buffer[position] != rune(')') {
						goto l791
					}
					position++
					depth--
					add(rulePegText, position793)
				}
				if !_rules[ruleAction51]() {
					goto l791
				}
				depth--
				add(ruleFuncTypeCast, position792)
			}
			return true
		l791:
			position, tokenIndex, depth = position791, tokenIndex791, depth791
			return false
		},
		/* 67 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action52)> */
		func() bool {
			position806, tokenIndex806, depth806 := position, tokenIndex, depth
			{
				position807 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l806
				}
				if !_rules[rulesp]() {
					goto l806
				}
				if buffer[position] != rune('(') {
					goto l806
				}
				position++
				if !_rules[rulesp]() {
					goto l806
				}
				if !_rules[ruleFuncParams]() {
					goto l806
				}
				if !_rules[rulesp]() {
					goto l806
				}
				if buffer[position] != rune(')') {
					goto l806
				}
				position++
				if !_rules[ruleAction52]() {
					goto l806
				}
				depth--
				add(ruleFuncApp, position807)
			}
			return true
		l806:
			position, tokenIndex, depth = position806, tokenIndex806, depth806
			return false
		},
		/* 68 FuncParams <- <(<(Star / (Expression sp (',' sp Expression)*)?)> Action53)> */
		func() bool {
			position808, tokenIndex808, depth808 := position, tokenIndex, depth
			{
				position809 := position
				depth++
				{
					position810 := position
					depth++
					{
						position811, tokenIndex811, depth811 := position, tokenIndex, depth
						if !_rules[ruleStar]() {
							goto l812
						}
						goto l811
					l812:
						position, tokenIndex, depth = position811, tokenIndex811, depth811
						{
							position813, tokenIndex813, depth813 := position, tokenIndex, depth
							if !_rules[ruleExpression]() {
								goto l813
							}
							if !_rules[rulesp]() {
								goto l813
							}
						l815:
							{
								position816, tokenIndex816, depth816 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l816
								}
								position++
								if !_rules[rulesp]() {
									goto l816
								}
								if !_rules[ruleExpression]() {
									goto l816
								}
								goto l815
							l816:
								position, tokenIndex, depth = position816, tokenIndex816, depth816
							}
							goto l814
						l813:
							position, tokenIndex, depth = position813, tokenIndex813, depth813
						}
					l814:
					}
				l811:
					depth--
					add(rulePegText, position810)
				}
				if !_rules[ruleAction53]() {
					goto l808
				}
				depth--
				add(ruleFuncParams, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 69 ArrayExpr <- <(<('[' sp (Expression (',' sp Expression)*)? sp ','? sp ']')> Action54)> */
		func() bool {
			position817, tokenIndex817, depth817 := position, tokenIndex, depth
			{
				position818 := position
				depth++
				{
					position819 := position
					depth++
					if buffer[position] != rune('[') {
						goto l817
					}
					position++
					if !_rules[rulesp]() {
						goto l817
					}
					{
						position820, tokenIndex820, depth820 := position, tokenIndex, depth
						if !_rules[ruleExpression]() {
							goto l820
						}
					l822:
						{
							position823, tokenIndex823, depth823 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l823
							}
							position++
							if !_rules[rulesp]() {
								goto l823
							}
							if !_rules[ruleExpression]() {
								goto l823
							}
							goto l822
						l823:
							position, tokenIndex, depth = position823, tokenIndex823, depth823
						}
						goto l821
					l820:
						position, tokenIndex, depth = position820, tokenIndex820, depth820
					}
				l821:
					if !_rules[rulesp]() {
						goto l817
					}
					{
						position824, tokenIndex824, depth824 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l824
						}
						position++
						goto l825
					l824:
						position, tokenIndex, depth = position824, tokenIndex824, depth824
					}
				l825:
					if !_rules[rulesp]() {
						goto l817
					}
					if buffer[position] != rune(']') {
						goto l817
					}
					position++
					depth--
					add(rulePegText, position819)
				}
				if !_rules[ruleAction54]() {
					goto l817
				}
				depth--
				add(ruleArrayExpr, position818)
			}
			return true
		l817:
			position, tokenIndex, depth = position817, tokenIndex817, depth817
			return false
		},
		/* 70 MapExpr <- <(<('{' sp (KeyValuePair (',' sp KeyValuePair)*)? sp '}')> Action55)> */
		func() bool {
			position826, tokenIndex826, depth826 := position, tokenIndex, depth
			{
				position827 := position
				depth++
				{
					position828 := position
					depth++
					if buffer[position] != rune('{') {
						goto l826
					}
					position++
					if !_rules[rulesp]() {
						goto l826
					}
					{
						position829, tokenIndex829, depth829 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l829
						}
					l831:
						{
							position832, tokenIndex832, depth832 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l832
							}
							position++
							if !_rules[rulesp]() {
								goto l832
							}
							if !_rules[ruleKeyValuePair]() {
								goto l832
							}
							goto l831
						l832:
							position, tokenIndex, depth = position832, tokenIndex832, depth832
						}
						goto l830
					l829:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
					}
				l830:
					if !_rules[rulesp]() {
						goto l826
					}
					if buffer[position] != rune('}') {
						goto l826
					}
					position++
					depth--
					add(rulePegText, position828)
				}
				if !_rules[ruleAction55]() {
					goto l826
				}
				depth--
				add(ruleMapExpr, position827)
			}
			return true
		l826:
			position, tokenIndex, depth = position826, tokenIndex826, depth826
			return false
		},
		/* 71 KeyValuePair <- <(<(StringLiteral sp ':' sp Expression)> Action56)> */
		func() bool {
			position833, tokenIndex833, depth833 := position, tokenIndex, depth
			{
				position834 := position
				depth++
				{
					position835 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l833
					}
					if !_rules[rulesp]() {
						goto l833
					}
					if buffer[position] != rune(':') {
						goto l833
					}
					position++
					if !_rules[rulesp]() {
						goto l833
					}
					if !_rules[ruleExpression]() {
						goto l833
					}
					depth--
					add(rulePegText, position835)
				}
				if !_rules[ruleAction56]() {
					goto l833
				}
				depth--
				add(ruleKeyValuePair, position834)
			}
			return true
		l833:
			position, tokenIndex, depth = position833, tokenIndex833, depth833
			return false
		},
		/* 72 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position836, tokenIndex836, depth836 := position, tokenIndex, depth
			{
				position837 := position
				depth++
				{
					position838, tokenIndex838, depth838 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l839
					}
					goto l838
				l839:
					position, tokenIndex, depth = position838, tokenIndex838, depth838
					if !_rules[ruleNumericLiteral]() {
						goto l840
					}
					goto l838
				l840:
					position, tokenIndex, depth = position838, tokenIndex838, depth838
					if !_rules[ruleStringLiteral]() {
						goto l836
					}
				}
			l838:
				depth--
				add(ruleLiteral, position837)
			}
			return true
		l836:
			position, tokenIndex, depth = position836, tokenIndex836, depth836
			return false
		},
		/* 73 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position841, tokenIndex841, depth841 := position, tokenIndex, depth
			{
				position842 := position
				depth++
				{
					position843, tokenIndex843, depth843 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l844
					}
					goto l843
				l844:
					position, tokenIndex, depth = position843, tokenIndex843, depth843
					if !_rules[ruleNotEqual]() {
						goto l845
					}
					goto l843
				l845:
					position, tokenIndex, depth = position843, tokenIndex843, depth843
					if !_rules[ruleLessOrEqual]() {
						goto l846
					}
					goto l843
				l846:
					position, tokenIndex, depth = position843, tokenIndex843, depth843
					if !_rules[ruleLess]() {
						goto l847
					}
					goto l843
				l847:
					position, tokenIndex, depth = position843, tokenIndex843, depth843
					if !_rules[ruleGreaterOrEqual]() {
						goto l848
					}
					goto l843
				l848:
					position, tokenIndex, depth = position843, tokenIndex843, depth843
					if !_rules[ruleGreater]() {
						goto l849
					}
					goto l843
				l849:
					position, tokenIndex, depth = position843, tokenIndex843, depth843
					if !_rules[ruleNotEqual]() {
						goto l841
					}
				}
			l843:
				depth--
				add(ruleComparisonOp, position842)
			}
			return true
		l841:
			position, tokenIndex, depth = position841, tokenIndex841, depth841
			return false
		},
		/* 74 OtherOp <- <Concat> */
		func() bool {
			position850, tokenIndex850, depth850 := position, tokenIndex, depth
			{
				position851 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l850
				}
				depth--
				add(ruleOtherOp, position851)
			}
			return true
		l850:
			position, tokenIndex, depth = position850, tokenIndex850, depth850
			return false
		},
		/* 75 IsOp <- <(IsNot / Is)> */
		func() bool {
			position852, tokenIndex852, depth852 := position, tokenIndex, depth
			{
				position853 := position
				depth++
				{
					position854, tokenIndex854, depth854 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l855
					}
					goto l854
				l855:
					position, tokenIndex, depth = position854, tokenIndex854, depth854
					if !_rules[ruleIs]() {
						goto l852
					}
				}
			l854:
				depth--
				add(ruleIsOp, position853)
			}
			return true
		l852:
			position, tokenIndex, depth = position852, tokenIndex852, depth852
			return false
		},
		/* 76 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position856, tokenIndex856, depth856 := position, tokenIndex, depth
			{
				position857 := position
				depth++
				{
					position858, tokenIndex858, depth858 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l859
					}
					goto l858
				l859:
					position, tokenIndex, depth = position858, tokenIndex858, depth858
					if !_rules[ruleMinus]() {
						goto l856
					}
				}
			l858:
				depth--
				add(rulePlusMinusOp, position857)
			}
			return true
		l856:
			position, tokenIndex, depth = position856, tokenIndex856, depth856
			return false
		},
		/* 77 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position860, tokenIndex860, depth860 := position, tokenIndex, depth
			{
				position861 := position
				depth++
				{
					position862, tokenIndex862, depth862 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l863
					}
					goto l862
				l863:
					position, tokenIndex, depth = position862, tokenIndex862, depth862
					if !_rules[ruleDivide]() {
						goto l864
					}
					goto l862
				l864:
					position, tokenIndex, depth = position862, tokenIndex862, depth862
					if !_rules[ruleModulo]() {
						goto l860
					}
				}
			l862:
				depth--
				add(ruleMultDivOp, position861)
			}
			return true
		l860:
			position, tokenIndex, depth = position860, tokenIndex860, depth860
			return false
		},
		/* 78 Stream <- <(<ident> Action57)> */
		func() bool {
			position865, tokenIndex865, depth865 := position, tokenIndex, depth
			{
				position866 := position
				depth++
				{
					position867 := position
					depth++
					if !_rules[ruleident]() {
						goto l865
					}
					depth--
					add(rulePegText, position867)
				}
				if !_rules[ruleAction57]() {
					goto l865
				}
				depth--
				add(ruleStream, position866)
			}
			return true
		l865:
			position, tokenIndex, depth = position865, tokenIndex865, depth865
			return false
		},
		/* 79 RowMeta <- <RowTimestamp> */
		func() bool {
			position868, tokenIndex868, depth868 := position, tokenIndex, depth
			{
				position869 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l868
				}
				depth--
				add(ruleRowMeta, position869)
			}
			return true
		l868:
			position, tokenIndex, depth = position868, tokenIndex868, depth868
			return false
		},
		/* 80 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action58)> */
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
						if !_rules[ruleident]() {
							goto l873
						}
						if buffer[position] != rune(':') {
							goto l873
						}
						position++
						goto l874
					l873:
						position, tokenIndex, depth = position873, tokenIndex873, depth873
					}
				l874:
					if buffer[position] != rune('t') {
						goto l870
					}
					position++
					if buffer[position] != rune('s') {
						goto l870
					}
					position++
					if buffer[position] != rune('(') {
						goto l870
					}
					position++
					if buffer[position] != rune(')') {
						goto l870
					}
					position++
					depth--
					add(rulePegText, position872)
				}
				if !_rules[ruleAction58]() {
					goto l870
				}
				depth--
				add(ruleRowTimestamp, position871)
			}
			return true
		l870:
			position, tokenIndex, depth = position870, tokenIndex870, depth870
			return false
		},
		/* 81 RowValue <- <(<((ident ':' !':')? jsonPath)> Action59)> */
		func() bool {
			position875, tokenIndex875, depth875 := position, tokenIndex, depth
			{
				position876 := position
				depth++
				{
					position877 := position
					depth++
					{
						position878, tokenIndex878, depth878 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l878
						}
						if buffer[position] != rune(':') {
							goto l878
						}
						position++
						{
							position880, tokenIndex880, depth880 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l880
							}
							position++
							goto l878
						l880:
							position, tokenIndex, depth = position880, tokenIndex880, depth880
						}
						goto l879
					l878:
						position, tokenIndex, depth = position878, tokenIndex878, depth878
					}
				l879:
					if !_rules[rulejsonPath]() {
						goto l875
					}
					depth--
					add(rulePegText, position877)
				}
				if !_rules[ruleAction59]() {
					goto l875
				}
				depth--
				add(ruleRowValue, position876)
			}
			return true
		l875:
			position, tokenIndex, depth = position875, tokenIndex875, depth875
			return false
		},
		/* 82 NumericLiteral <- <(<('-'? [0-9]+)> Action60)> */
		func() bool {
			position881, tokenIndex881, depth881 := position, tokenIndex, depth
			{
				position882 := position
				depth++
				{
					position883 := position
					depth++
					{
						position884, tokenIndex884, depth884 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l884
						}
						position++
						goto l885
					l884:
						position, tokenIndex, depth = position884, tokenIndex884, depth884
					}
				l885:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l881
					}
					position++
				l886:
					{
						position887, tokenIndex887, depth887 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l887
						}
						position++
						goto l886
					l887:
						position, tokenIndex, depth = position887, tokenIndex887, depth887
					}
					depth--
					add(rulePegText, position883)
				}
				if !_rules[ruleAction60]() {
					goto l881
				}
				depth--
				add(ruleNumericLiteral, position882)
			}
			return true
		l881:
			position, tokenIndex, depth = position881, tokenIndex881, depth881
			return false
		},
		/* 83 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action61)> */
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
						if buffer[position] != rune('-') {
							goto l891
						}
						position++
						goto l892
					l891:
						position, tokenIndex, depth = position891, tokenIndex891, depth891
					}
				l892:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l888
					}
					position++
				l893:
					{
						position894, tokenIndex894, depth894 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l894
						}
						position++
						goto l893
					l894:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
					}
					if buffer[position] != rune('.') {
						goto l888
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l888
					}
					position++
				l895:
					{
						position896, tokenIndex896, depth896 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l896
						}
						position++
						goto l895
					l896:
						position, tokenIndex, depth = position896, tokenIndex896, depth896
					}
					depth--
					add(rulePegText, position890)
				}
				if !_rules[ruleAction61]() {
					goto l888
				}
				depth--
				add(ruleFloatLiteral, position889)
			}
			return true
		l888:
			position, tokenIndex, depth = position888, tokenIndex888, depth888
			return false
		},
		/* 84 Function <- <(<ident> Action62)> */
		func() bool {
			position897, tokenIndex897, depth897 := position, tokenIndex, depth
			{
				position898 := position
				depth++
				{
					position899 := position
					depth++
					if !_rules[ruleident]() {
						goto l897
					}
					depth--
					add(rulePegText, position899)
				}
				if !_rules[ruleAction62]() {
					goto l897
				}
				depth--
				add(ruleFunction, position898)
			}
			return true
		l897:
			position, tokenIndex, depth = position897, tokenIndex897, depth897
			return false
		},
		/* 85 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action63)> */
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
						if buffer[position] != rune('n') {
							goto l904
						}
						position++
						goto l903
					l904:
						position, tokenIndex, depth = position903, tokenIndex903, depth903
						if buffer[position] != rune('N') {
							goto l900
						}
						position++
					}
				l903:
					{
						position905, tokenIndex905, depth905 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l906
						}
						position++
						goto l905
					l906:
						position, tokenIndex, depth = position905, tokenIndex905, depth905
						if buffer[position] != rune('U') {
							goto l900
						}
						position++
					}
				l905:
					{
						position907, tokenIndex907, depth907 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l908
						}
						position++
						goto l907
					l908:
						position, tokenIndex, depth = position907, tokenIndex907, depth907
						if buffer[position] != rune('L') {
							goto l900
						}
						position++
					}
				l907:
					{
						position909, tokenIndex909, depth909 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l910
						}
						position++
						goto l909
					l910:
						position, tokenIndex, depth = position909, tokenIndex909, depth909
						if buffer[position] != rune('L') {
							goto l900
						}
						position++
					}
				l909:
					depth--
					add(rulePegText, position902)
				}
				if !_rules[ruleAction63]() {
					goto l900
				}
				depth--
				add(ruleNullLiteral, position901)
			}
			return true
		l900:
			position, tokenIndex, depth = position900, tokenIndex900, depth900
			return false
		},
		/* 86 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position911, tokenIndex911, depth911 := position, tokenIndex, depth
			{
				position912 := position
				depth++
				{
					position913, tokenIndex913, depth913 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l914
					}
					goto l913
				l914:
					position, tokenIndex, depth = position913, tokenIndex913, depth913
					if !_rules[ruleFALSE]() {
						goto l911
					}
				}
			l913:
				depth--
				add(ruleBooleanLiteral, position912)
			}
			return true
		l911:
			position, tokenIndex, depth = position911, tokenIndex911, depth911
			return false
		},
		/* 87 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action64)> */
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
						if buffer[position] != rune('t') {
							goto l919
						}
						position++
						goto l918
					l919:
						position, tokenIndex, depth = position918, tokenIndex918, depth918
						if buffer[position] != rune('T') {
							goto l915
						}
						position++
					}
				l918:
					{
						position920, tokenIndex920, depth920 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l921
						}
						position++
						goto l920
					l921:
						position, tokenIndex, depth = position920, tokenIndex920, depth920
						if buffer[position] != rune('R') {
							goto l915
						}
						position++
					}
				l920:
					{
						position922, tokenIndex922, depth922 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l923
						}
						position++
						goto l922
					l923:
						position, tokenIndex, depth = position922, tokenIndex922, depth922
						if buffer[position] != rune('U') {
							goto l915
						}
						position++
					}
				l922:
					{
						position924, tokenIndex924, depth924 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l925
						}
						position++
						goto l924
					l925:
						position, tokenIndex, depth = position924, tokenIndex924, depth924
						if buffer[position] != rune('E') {
							goto l915
						}
						position++
					}
				l924:
					depth--
					add(rulePegText, position917)
				}
				if !_rules[ruleAction64]() {
					goto l915
				}
				depth--
				add(ruleTRUE, position916)
			}
			return true
		l915:
			position, tokenIndex, depth = position915, tokenIndex915, depth915
			return false
		},
		/* 88 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action65)> */
		func() bool {
			position926, tokenIndex926, depth926 := position, tokenIndex, depth
			{
				position927 := position
				depth++
				{
					position928 := position
					depth++
					{
						position929, tokenIndex929, depth929 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l930
						}
						position++
						goto l929
					l930:
						position, tokenIndex, depth = position929, tokenIndex929, depth929
						if buffer[position] != rune('F') {
							goto l926
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
							goto l926
						}
						position++
					}
				l931:
					{
						position933, tokenIndex933, depth933 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l934
						}
						position++
						goto l933
					l934:
						position, tokenIndex, depth = position933, tokenIndex933, depth933
						if buffer[position] != rune('L') {
							goto l926
						}
						position++
					}
				l933:
					{
						position935, tokenIndex935, depth935 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l936
						}
						position++
						goto l935
					l936:
						position, tokenIndex, depth = position935, tokenIndex935, depth935
						if buffer[position] != rune('S') {
							goto l926
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
							goto l926
						}
						position++
					}
				l937:
					depth--
					add(rulePegText, position928)
				}
				if !_rules[ruleAction65]() {
					goto l926
				}
				depth--
				add(ruleFALSE, position927)
			}
			return true
		l926:
			position, tokenIndex, depth = position926, tokenIndex926, depth926
			return false
		},
		/* 89 Star <- <(<'*'> Action66)> */
		func() bool {
			position939, tokenIndex939, depth939 := position, tokenIndex, depth
			{
				position940 := position
				depth++
				{
					position941 := position
					depth++
					if buffer[position] != rune('*') {
						goto l939
					}
					position++
					depth--
					add(rulePegText, position941)
				}
				if !_rules[ruleAction66]() {
					goto l939
				}
				depth--
				add(ruleStar, position940)
			}
			return true
		l939:
			position, tokenIndex, depth = position939, tokenIndex939, depth939
			return false
		},
		/* 90 Wildcard <- <(<((ident ':' !':')? '*')> Action67)> */
		func() bool {
			position942, tokenIndex942, depth942 := position, tokenIndex, depth
			{
				position943 := position
				depth++
				{
					position944 := position
					depth++
					{
						position945, tokenIndex945, depth945 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l945
						}
						if buffer[position] != rune(':') {
							goto l945
						}
						position++
						{
							position947, tokenIndex947, depth947 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l947
							}
							position++
							goto l945
						l947:
							position, tokenIndex, depth = position947, tokenIndex947, depth947
						}
						goto l946
					l945:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
					}
				l946:
					if buffer[position] != rune('*') {
						goto l942
					}
					position++
					depth--
					add(rulePegText, position944)
				}
				if !_rules[ruleAction67]() {
					goto l942
				}
				depth--
				add(ruleWildcard, position943)
			}
			return true
		l942:
			position, tokenIndex, depth = position942, tokenIndex942, depth942
			return false
		},
		/* 91 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action68)> */
		func() bool {
			position948, tokenIndex948, depth948 := position, tokenIndex, depth
			{
				position949 := position
				depth++
				{
					position950 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l948
					}
					position++
				l951:
					{
						position952, tokenIndex952, depth952 := position, tokenIndex, depth
						{
							position953, tokenIndex953, depth953 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l954
							}
							position++
							if buffer[position] != rune('\'') {
								goto l954
							}
							position++
							goto l953
						l954:
							position, tokenIndex, depth = position953, tokenIndex953, depth953
							{
								position955, tokenIndex955, depth955 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l955
								}
								position++
								goto l952
							l955:
								position, tokenIndex, depth = position955, tokenIndex955, depth955
							}
							if !matchDot() {
								goto l952
							}
						}
					l953:
						goto l951
					l952:
						position, tokenIndex, depth = position952, tokenIndex952, depth952
					}
					if buffer[position] != rune('\'') {
						goto l948
					}
					position++
					depth--
					add(rulePegText, position950)
				}
				if !_rules[ruleAction68]() {
					goto l948
				}
				depth--
				add(ruleStringLiteral, position949)
			}
			return true
		l948:
			position, tokenIndex, depth = position948, tokenIndex948, depth948
			return false
		},
		/* 92 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action69)> */
		func() bool {
			position956, tokenIndex956, depth956 := position, tokenIndex, depth
			{
				position957 := position
				depth++
				{
					position958 := position
					depth++
					{
						position959, tokenIndex959, depth959 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l960
						}
						position++
						goto l959
					l960:
						position, tokenIndex, depth = position959, tokenIndex959, depth959
						if buffer[position] != rune('I') {
							goto l956
						}
						position++
					}
				l959:
					{
						position961, tokenIndex961, depth961 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l962
						}
						position++
						goto l961
					l962:
						position, tokenIndex, depth = position961, tokenIndex961, depth961
						if buffer[position] != rune('S') {
							goto l956
						}
						position++
					}
				l961:
					{
						position963, tokenIndex963, depth963 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l964
						}
						position++
						goto l963
					l964:
						position, tokenIndex, depth = position963, tokenIndex963, depth963
						if buffer[position] != rune('T') {
							goto l956
						}
						position++
					}
				l963:
					{
						position965, tokenIndex965, depth965 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l966
						}
						position++
						goto l965
					l966:
						position, tokenIndex, depth = position965, tokenIndex965, depth965
						if buffer[position] != rune('R') {
							goto l956
						}
						position++
					}
				l965:
					{
						position967, tokenIndex967, depth967 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l968
						}
						position++
						goto l967
					l968:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if buffer[position] != rune('E') {
							goto l956
						}
						position++
					}
				l967:
					{
						position969, tokenIndex969, depth969 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l970
						}
						position++
						goto l969
					l970:
						position, tokenIndex, depth = position969, tokenIndex969, depth969
						if buffer[position] != rune('A') {
							goto l956
						}
						position++
					}
				l969:
					{
						position971, tokenIndex971, depth971 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l972
						}
						position++
						goto l971
					l972:
						position, tokenIndex, depth = position971, tokenIndex971, depth971
						if buffer[position] != rune('M') {
							goto l956
						}
						position++
					}
				l971:
					depth--
					add(rulePegText, position958)
				}
				if !_rules[ruleAction69]() {
					goto l956
				}
				depth--
				add(ruleISTREAM, position957)
			}
			return true
		l956:
			position, tokenIndex, depth = position956, tokenIndex956, depth956
			return false
		},
		/* 93 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action70)> */
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
						if buffer[position] != rune('d') {
							goto l977
						}
						position++
						goto l976
					l977:
						position, tokenIndex, depth = position976, tokenIndex976, depth976
						if buffer[position] != rune('D') {
							goto l973
						}
						position++
					}
				l976:
					{
						position978, tokenIndex978, depth978 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l979
						}
						position++
						goto l978
					l979:
						position, tokenIndex, depth = position978, tokenIndex978, depth978
						if buffer[position] != rune('S') {
							goto l973
						}
						position++
					}
				l978:
					{
						position980, tokenIndex980, depth980 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l981
						}
						position++
						goto l980
					l981:
						position, tokenIndex, depth = position980, tokenIndex980, depth980
						if buffer[position] != rune('T') {
							goto l973
						}
						position++
					}
				l980:
					{
						position982, tokenIndex982, depth982 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l983
						}
						position++
						goto l982
					l983:
						position, tokenIndex, depth = position982, tokenIndex982, depth982
						if buffer[position] != rune('R') {
							goto l973
						}
						position++
					}
				l982:
					{
						position984, tokenIndex984, depth984 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l985
						}
						position++
						goto l984
					l985:
						position, tokenIndex, depth = position984, tokenIndex984, depth984
						if buffer[position] != rune('E') {
							goto l973
						}
						position++
					}
				l984:
					{
						position986, tokenIndex986, depth986 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l987
						}
						position++
						goto l986
					l987:
						position, tokenIndex, depth = position986, tokenIndex986, depth986
						if buffer[position] != rune('A') {
							goto l973
						}
						position++
					}
				l986:
					{
						position988, tokenIndex988, depth988 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l989
						}
						position++
						goto l988
					l989:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('M') {
							goto l973
						}
						position++
					}
				l988:
					depth--
					add(rulePegText, position975)
				}
				if !_rules[ruleAction70]() {
					goto l973
				}
				depth--
				add(ruleDSTREAM, position974)
			}
			return true
		l973:
			position, tokenIndex, depth = position973, tokenIndex973, depth973
			return false
		},
		/* 94 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action71)> */
		func() bool {
			position990, tokenIndex990, depth990 := position, tokenIndex, depth
			{
				position991 := position
				depth++
				{
					position992 := position
					depth++
					{
						position993, tokenIndex993, depth993 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l994
						}
						position++
						goto l993
					l994:
						position, tokenIndex, depth = position993, tokenIndex993, depth993
						if buffer[position] != rune('R') {
							goto l990
						}
						position++
					}
				l993:
					{
						position995, tokenIndex995, depth995 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l996
						}
						position++
						goto l995
					l996:
						position, tokenIndex, depth = position995, tokenIndex995, depth995
						if buffer[position] != rune('S') {
							goto l990
						}
						position++
					}
				l995:
					{
						position997, tokenIndex997, depth997 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l998
						}
						position++
						goto l997
					l998:
						position, tokenIndex, depth = position997, tokenIndex997, depth997
						if buffer[position] != rune('T') {
							goto l990
						}
						position++
					}
				l997:
					{
						position999, tokenIndex999, depth999 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1000
						}
						position++
						goto l999
					l1000:
						position, tokenIndex, depth = position999, tokenIndex999, depth999
						if buffer[position] != rune('R') {
							goto l990
						}
						position++
					}
				l999:
					{
						position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1002
						}
						position++
						goto l1001
					l1002:
						position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
						if buffer[position] != rune('E') {
							goto l990
						}
						position++
					}
				l1001:
					{
						position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1004
						}
						position++
						goto l1003
					l1004:
						position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
						if buffer[position] != rune('A') {
							goto l990
						}
						position++
					}
				l1003:
					{
						position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1006
						}
						position++
						goto l1005
					l1006:
						position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
						if buffer[position] != rune('M') {
							goto l990
						}
						position++
					}
				l1005:
					depth--
					add(rulePegText, position992)
				}
				if !_rules[ruleAction71]() {
					goto l990
				}
				depth--
				add(ruleRSTREAM, position991)
			}
			return true
		l990:
			position, tokenIndex, depth = position990, tokenIndex990, depth990
			return false
		},
		/* 95 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action72)> */
		func() bool {
			position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
			{
				position1008 := position
				depth++
				{
					position1009 := position
					depth++
					{
						position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1011
						}
						position++
						goto l1010
					l1011:
						position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
						if buffer[position] != rune('T') {
							goto l1007
						}
						position++
					}
				l1010:
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1013
						}
						position++
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						if buffer[position] != rune('U') {
							goto l1007
						}
						position++
					}
				l1012:
					{
						position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1015
						}
						position++
						goto l1014
					l1015:
						position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
						if buffer[position] != rune('P') {
							goto l1007
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
							goto l1007
						}
						position++
					}
				l1016:
					{
						position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1019
						}
						position++
						goto l1018
					l1019:
						position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
						if buffer[position] != rune('E') {
							goto l1007
						}
						position++
					}
				l1018:
					{
						position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1021
						}
						position++
						goto l1020
					l1021:
						position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
						if buffer[position] != rune('S') {
							goto l1007
						}
						position++
					}
				l1020:
					depth--
					add(rulePegText, position1009)
				}
				if !_rules[ruleAction72]() {
					goto l1007
				}
				depth--
				add(ruleTUPLES, position1008)
			}
			return true
		l1007:
			position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
			return false
		},
		/* 96 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action73)> */
		func() bool {
			position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
			{
				position1023 := position
				depth++
				{
					position1024 := position
					depth++
					{
						position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1026
						}
						position++
						goto l1025
					l1026:
						position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
						if buffer[position] != rune('S') {
							goto l1022
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
							goto l1022
						}
						position++
					}
				l1027:
					{
						position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1030
						}
						position++
						goto l1029
					l1030:
						position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
						if buffer[position] != rune('C') {
							goto l1022
						}
						position++
					}
				l1029:
					{
						position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1032
						}
						position++
						goto l1031
					l1032:
						position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
						if buffer[position] != rune('O') {
							goto l1022
						}
						position++
					}
				l1031:
					{
						position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1034
						}
						position++
						goto l1033
					l1034:
						position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
						if buffer[position] != rune('N') {
							goto l1022
						}
						position++
					}
				l1033:
					{
						position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1036
						}
						position++
						goto l1035
					l1036:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						if buffer[position] != rune('D') {
							goto l1022
						}
						position++
					}
				l1035:
					{
						position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1038
						}
						position++
						goto l1037
					l1038:
						position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
						if buffer[position] != rune('S') {
							goto l1022
						}
						position++
					}
				l1037:
					depth--
					add(rulePegText, position1024)
				}
				if !_rules[ruleAction73]() {
					goto l1022
				}
				depth--
				add(ruleSECONDS, position1023)
			}
			return true
		l1022:
			position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
			return false
		},
		/* 97 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action74)> */
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
						if buffer[position] != rune('m') {
							goto l1043
						}
						position++
						goto l1042
					l1043:
						position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
						if buffer[position] != rune('M') {
							goto l1039
						}
						position++
					}
				l1042:
					{
						position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1045
						}
						position++
						goto l1044
					l1045:
						position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
						if buffer[position] != rune('I') {
							goto l1039
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
							goto l1039
						}
						position++
					}
				l1046:
					{
						position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1049
						}
						position++
						goto l1048
					l1049:
						position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
						if buffer[position] != rune('L') {
							goto l1039
						}
						position++
					}
				l1048:
					{
						position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1051
						}
						position++
						goto l1050
					l1051:
						position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
						if buffer[position] != rune('I') {
							goto l1039
						}
						position++
					}
				l1050:
					{
						position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1053
						}
						position++
						goto l1052
					l1053:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
						if buffer[position] != rune('S') {
							goto l1039
						}
						position++
					}
				l1052:
					{
						position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1055
						}
						position++
						goto l1054
					l1055:
						position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
						if buffer[position] != rune('E') {
							goto l1039
						}
						position++
					}
				l1054:
					{
						position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1057
						}
						position++
						goto l1056
					l1057:
						position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
						if buffer[position] != rune('C') {
							goto l1039
						}
						position++
					}
				l1056:
					{
						position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1059
						}
						position++
						goto l1058
					l1059:
						position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
						if buffer[position] != rune('O') {
							goto l1039
						}
						position++
					}
				l1058:
					{
						position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1061
						}
						position++
						goto l1060
					l1061:
						position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
						if buffer[position] != rune('N') {
							goto l1039
						}
						position++
					}
				l1060:
					{
						position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1063
						}
						position++
						goto l1062
					l1063:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune('D') {
							goto l1039
						}
						position++
					}
				l1062:
					{
						position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1065
						}
						position++
						goto l1064
					l1065:
						position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
						if buffer[position] != rune('S') {
							goto l1039
						}
						position++
					}
				l1064:
					depth--
					add(rulePegText, position1041)
				}
				if !_rules[ruleAction74]() {
					goto l1039
				}
				depth--
				add(ruleMILLISECONDS, position1040)
			}
			return true
		l1039:
			position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
			return false
		},
		/* 98 StreamIdentifier <- <(<ident> Action75)> */
		func() bool {
			position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
			{
				position1067 := position
				depth++
				{
					position1068 := position
					depth++
					if !_rules[ruleident]() {
						goto l1066
					}
					depth--
					add(rulePegText, position1068)
				}
				if !_rules[ruleAction75]() {
					goto l1066
				}
				depth--
				add(ruleStreamIdentifier, position1067)
			}
			return true
		l1066:
			position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
			return false
		},
		/* 99 SourceSinkType <- <(<ident> Action76)> */
		func() bool {
			position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
			{
				position1070 := position
				depth++
				{
					position1071 := position
					depth++
					if !_rules[ruleident]() {
						goto l1069
					}
					depth--
					add(rulePegText, position1071)
				}
				if !_rules[ruleAction76]() {
					goto l1069
				}
				depth--
				add(ruleSourceSinkType, position1070)
			}
			return true
		l1069:
			position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
			return false
		},
		/* 100 SourceSinkParamKey <- <(<ident> Action77)> */
		func() bool {
			position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
			{
				position1073 := position
				depth++
				{
					position1074 := position
					depth++
					if !_rules[ruleident]() {
						goto l1072
					}
					depth--
					add(rulePegText, position1074)
				}
				if !_rules[ruleAction77]() {
					goto l1072
				}
				depth--
				add(ruleSourceSinkParamKey, position1073)
			}
			return true
		l1072:
			position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
			return false
		},
		/* 101 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action78)> */
		func() bool {
			position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
			{
				position1076 := position
				depth++
				{
					position1077 := position
					depth++
					{
						position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1079
						}
						position++
						goto l1078
					l1079:
						position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
						if buffer[position] != rune('P') {
							goto l1075
						}
						position++
					}
				l1078:
					{
						position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1081
						}
						position++
						goto l1080
					l1081:
						position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
						if buffer[position] != rune('A') {
							goto l1075
						}
						position++
					}
				l1080:
					{
						position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1083
						}
						position++
						goto l1082
					l1083:
						position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
						if buffer[position] != rune('U') {
							goto l1075
						}
						position++
					}
				l1082:
					{
						position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1085
						}
						position++
						goto l1084
					l1085:
						position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
						if buffer[position] != rune('S') {
							goto l1075
						}
						position++
					}
				l1084:
					{
						position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1087
						}
						position++
						goto l1086
					l1087:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if buffer[position] != rune('E') {
							goto l1075
						}
						position++
					}
				l1086:
					{
						position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1089
						}
						position++
						goto l1088
					l1089:
						position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
						if buffer[position] != rune('D') {
							goto l1075
						}
						position++
					}
				l1088:
					depth--
					add(rulePegText, position1077)
				}
				if !_rules[ruleAction78]() {
					goto l1075
				}
				depth--
				add(rulePaused, position1076)
			}
			return true
		l1075:
			position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
			return false
		},
		/* 102 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action79)> */
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
						if buffer[position] != rune('u') {
							goto l1094
						}
						position++
						goto l1093
					l1094:
						position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
						if buffer[position] != rune('U') {
							goto l1090
						}
						position++
					}
				l1093:
					{
						position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1096
						}
						position++
						goto l1095
					l1096:
						position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
						if buffer[position] != rune('N') {
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
						if buffer[position] != rune('a') {
							goto l1100
						}
						position++
						goto l1099
					l1100:
						position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
						if buffer[position] != rune('A') {
							goto l1090
						}
						position++
					}
				l1099:
					{
						position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1102
						}
						position++
						goto l1101
					l1102:
						position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
						if buffer[position] != rune('U') {
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
							goto l1090
						}
						position++
					}
				l1105:
					{
						position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1108
						}
						position++
						goto l1107
					l1108:
						position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
						if buffer[position] != rune('D') {
							goto l1090
						}
						position++
					}
				l1107:
					depth--
					add(rulePegText, position1092)
				}
				if !_rules[ruleAction79]() {
					goto l1090
				}
				depth--
				add(ruleUnpaused, position1091)
			}
			return true
		l1090:
			position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
			return false
		},
		/* 103 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
			{
				position1110 := position
				depth++
				{
					position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1112
					}
					goto l1111
				l1112:
					position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					if !_rules[ruleInt]() {
						goto l1113
					}
					goto l1111
				l1113:
					position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					if !_rules[ruleFloat]() {
						goto l1114
					}
					goto l1111
				l1114:
					position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					if !_rules[ruleString]() {
						goto l1115
					}
					goto l1111
				l1115:
					position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					if !_rules[ruleBlob]() {
						goto l1116
					}
					goto l1111
				l1116:
					position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					if !_rules[ruleTimestamp]() {
						goto l1117
					}
					goto l1111
				l1117:
					position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					if !_rules[ruleArray]() {
						goto l1118
					}
					goto l1111
				l1118:
					position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
					if !_rules[ruleMap]() {
						goto l1109
					}
				}
			l1111:
				depth--
				add(ruleType, position1110)
			}
			return true
		l1109:
			position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
			return false
		},
		/* 104 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action80)> */
		func() bool {
			position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
			{
				position1120 := position
				depth++
				{
					position1121 := position
					depth++
					{
						position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1123
						}
						position++
						goto l1122
					l1123:
						position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
						if buffer[position] != rune('B') {
							goto l1119
						}
						position++
					}
				l1122:
					{
						position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1125
						}
						position++
						goto l1124
					l1125:
						position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
						if buffer[position] != rune('O') {
							goto l1119
						}
						position++
					}
				l1124:
					{
						position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1127
						}
						position++
						goto l1126
					l1127:
						position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
						if buffer[position] != rune('O') {
							goto l1119
						}
						position++
					}
				l1126:
					{
						position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1129
						}
						position++
						goto l1128
					l1129:
						position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
						if buffer[position] != rune('L') {
							goto l1119
						}
						position++
					}
				l1128:
					depth--
					add(rulePegText, position1121)
				}
				if !_rules[ruleAction80]() {
					goto l1119
				}
				depth--
				add(ruleBool, position1120)
			}
			return true
		l1119:
			position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
			return false
		},
		/* 105 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action81)> */
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
						if buffer[position] != rune('i') {
							goto l1134
						}
						position++
						goto l1133
					l1134:
						position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
						if buffer[position] != rune('I') {
							goto l1130
						}
						position++
					}
				l1133:
					{
						position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
						if buffer[position] != rune('N') {
							goto l1130
						}
						position++
					}
				l1135:
					{
						position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
						if buffer[position] != rune('T') {
							goto l1130
						}
						position++
					}
				l1137:
					depth--
					add(rulePegText, position1132)
				}
				if !_rules[ruleAction81]() {
					goto l1130
				}
				depth--
				add(ruleInt, position1131)
			}
			return true
		l1130:
			position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
			return false
		},
		/* 106 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action82)> */
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
						if buffer[position] != rune('f') {
							goto l1143
						}
						position++
						goto l1142
					l1143:
						position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
						if buffer[position] != rune('F') {
							goto l1139
						}
						position++
					}
				l1142:
					{
						position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1145
						}
						position++
						goto l1144
					l1145:
						position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
						if buffer[position] != rune('L') {
							goto l1139
						}
						position++
					}
				l1144:
					{
						position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1147
						}
						position++
						goto l1146
					l1147:
						position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
						if buffer[position] != rune('O') {
							goto l1139
						}
						position++
					}
				l1146:
					{
						position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1149
						}
						position++
						goto l1148
					l1149:
						position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
						if buffer[position] != rune('A') {
							goto l1139
						}
						position++
					}
				l1148:
					{
						position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1151
						}
						position++
						goto l1150
					l1151:
						position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
						if buffer[position] != rune('T') {
							goto l1139
						}
						position++
					}
				l1150:
					depth--
					add(rulePegText, position1141)
				}
				if !_rules[ruleAction82]() {
					goto l1139
				}
				depth--
				add(ruleFloat, position1140)
			}
			return true
		l1139:
			position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
			return false
		},
		/* 107 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action83)> */
		func() bool {
			position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
			{
				position1153 := position
				depth++
				{
					position1154 := position
					depth++
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
							goto l1152
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
							goto l1152
						}
						position++
					}
				l1157:
					{
						position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1160
						}
						position++
						goto l1159
					l1160:
						position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
						if buffer[position] != rune('R') {
							goto l1152
						}
						position++
					}
				l1159:
					{
						position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1162
						}
						position++
						goto l1161
					l1162:
						position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
						if buffer[position] != rune('I') {
							goto l1152
						}
						position++
					}
				l1161:
					{
						position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1164
						}
						position++
						goto l1163
					l1164:
						position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
						if buffer[position] != rune('N') {
							goto l1152
						}
						position++
					}
				l1163:
					{
						position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1166
						}
						position++
						goto l1165
					l1166:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						if buffer[position] != rune('G') {
							goto l1152
						}
						position++
					}
				l1165:
					depth--
					add(rulePegText, position1154)
				}
				if !_rules[ruleAction83]() {
					goto l1152
				}
				depth--
				add(ruleString, position1153)
			}
			return true
		l1152:
			position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
			return false
		},
		/* 108 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action84)> */
		func() bool {
			position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
			{
				position1168 := position
				depth++
				{
					position1169 := position
					depth++
					{
						position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1171
						}
						position++
						goto l1170
					l1171:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						if buffer[position] != rune('B') {
							goto l1167
						}
						position++
					}
				l1170:
					{
						position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1173
						}
						position++
						goto l1172
					l1173:
						position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
						if buffer[position] != rune('L') {
							goto l1167
						}
						position++
					}
				l1172:
					{
						position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1175
						}
						position++
						goto l1174
					l1175:
						position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
						if buffer[position] != rune('O') {
							goto l1167
						}
						position++
					}
				l1174:
					{
						position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1177
						}
						position++
						goto l1176
					l1177:
						position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
						if buffer[position] != rune('B') {
							goto l1167
						}
						position++
					}
				l1176:
					depth--
					add(rulePegText, position1169)
				}
				if !_rules[ruleAction84]() {
					goto l1167
				}
				depth--
				add(ruleBlob, position1168)
			}
			return true
		l1167:
			position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
			return false
		},
		/* 109 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action85)> */
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
						if buffer[position] != rune('t') {
							goto l1182
						}
						position++
						goto l1181
					l1182:
						position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
						if buffer[position] != rune('T') {
							goto l1178
						}
						position++
					}
				l1181:
					{
						position1183, tokenIndex1183, depth1183 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1184
						}
						position++
						goto l1183
					l1184:
						position, tokenIndex, depth = position1183, tokenIndex1183, depth1183
						if buffer[position] != rune('I') {
							goto l1178
						}
						position++
					}
				l1183:
					{
						position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1186
						}
						position++
						goto l1185
					l1186:
						position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
						if buffer[position] != rune('M') {
							goto l1178
						}
						position++
					}
				l1185:
					{
						position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1188
						}
						position++
						goto l1187
					l1188:
						position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
						if buffer[position] != rune('E') {
							goto l1178
						}
						position++
					}
				l1187:
					{
						position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1190
						}
						position++
						goto l1189
					l1190:
						position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
						if buffer[position] != rune('S') {
							goto l1178
						}
						position++
					}
				l1189:
					{
						position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1192
						}
						position++
						goto l1191
					l1192:
						position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
						if buffer[position] != rune('T') {
							goto l1178
						}
						position++
					}
				l1191:
					{
						position1193, tokenIndex1193, depth1193 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1194
						}
						position++
						goto l1193
					l1194:
						position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
						if buffer[position] != rune('A') {
							goto l1178
						}
						position++
					}
				l1193:
					{
						position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1196
						}
						position++
						goto l1195
					l1196:
						position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
						if buffer[position] != rune('M') {
							goto l1178
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
							goto l1178
						}
						position++
					}
				l1197:
					depth--
					add(rulePegText, position1180)
				}
				if !_rules[ruleAction85]() {
					goto l1178
				}
				depth--
				add(ruleTimestamp, position1179)
			}
			return true
		l1178:
			position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
			return false
		},
		/* 110 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action86)> */
		func() bool {
			position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
			{
				position1200 := position
				depth++
				{
					position1201 := position
					depth++
					{
						position1202, tokenIndex1202, depth1202 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1203
						}
						position++
						goto l1202
					l1203:
						position, tokenIndex, depth = position1202, tokenIndex1202, depth1202
						if buffer[position] != rune('A') {
							goto l1199
						}
						position++
					}
				l1202:
					{
						position1204, tokenIndex1204, depth1204 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1205
						}
						position++
						goto l1204
					l1205:
						position, tokenIndex, depth = position1204, tokenIndex1204, depth1204
						if buffer[position] != rune('R') {
							goto l1199
						}
						position++
					}
				l1204:
					{
						position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1207
						}
						position++
						goto l1206
					l1207:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						if buffer[position] != rune('R') {
							goto l1199
						}
						position++
					}
				l1206:
					{
						position1208, tokenIndex1208, depth1208 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1209
						}
						position++
						goto l1208
					l1209:
						position, tokenIndex, depth = position1208, tokenIndex1208, depth1208
						if buffer[position] != rune('A') {
							goto l1199
						}
						position++
					}
				l1208:
					{
						position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1211
						}
						position++
						goto l1210
					l1211:
						position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
						if buffer[position] != rune('Y') {
							goto l1199
						}
						position++
					}
				l1210:
					depth--
					add(rulePegText, position1201)
				}
				if !_rules[ruleAction86]() {
					goto l1199
				}
				depth--
				add(ruleArray, position1200)
			}
			return true
		l1199:
			position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
			return false
		},
		/* 111 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action87)> */
		func() bool {
			position1212, tokenIndex1212, depth1212 := position, tokenIndex, depth
			{
				position1213 := position
				depth++
				{
					position1214 := position
					depth++
					{
						position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1216
						}
						position++
						goto l1215
					l1216:
						position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
						if buffer[position] != rune('M') {
							goto l1212
						}
						position++
					}
				l1215:
					{
						position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1218
						}
						position++
						goto l1217
					l1218:
						position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
						if buffer[position] != rune('A') {
							goto l1212
						}
						position++
					}
				l1217:
					{
						position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1220
						}
						position++
						goto l1219
					l1220:
						position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
						if buffer[position] != rune('P') {
							goto l1212
						}
						position++
					}
				l1219:
					depth--
					add(rulePegText, position1214)
				}
				if !_rules[ruleAction87]() {
					goto l1212
				}
				depth--
				add(ruleMap, position1213)
			}
			return true
		l1212:
			position, tokenIndex, depth = position1212, tokenIndex1212, depth1212
			return false
		},
		/* 112 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action88)> */
		func() bool {
			position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
			{
				position1222 := position
				depth++
				{
					position1223 := position
					depth++
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
							goto l1221
						}
						position++
					}
				l1224:
					{
						position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1227
						}
						position++
						goto l1226
					l1227:
						position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
						if buffer[position] != rune('R') {
							goto l1221
						}
						position++
					}
				l1226:
					depth--
					add(rulePegText, position1223)
				}
				if !_rules[ruleAction88]() {
					goto l1221
				}
				depth--
				add(ruleOr, position1222)
			}
			return true
		l1221:
			position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
			return false
		},
		/* 113 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action89)> */
		func() bool {
			position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
			{
				position1229 := position
				depth++
				{
					position1230 := position
					depth++
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
							goto l1228
						}
						position++
					}
				l1231:
					{
						position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1234
						}
						position++
						goto l1233
					l1234:
						position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
						if buffer[position] != rune('N') {
							goto l1228
						}
						position++
					}
				l1233:
					{
						position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1236
						}
						position++
						goto l1235
					l1236:
						position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
						if buffer[position] != rune('D') {
							goto l1228
						}
						position++
					}
				l1235:
					depth--
					add(rulePegText, position1230)
				}
				if !_rules[ruleAction89]() {
					goto l1228
				}
				depth--
				add(ruleAnd, position1229)
			}
			return true
		l1228:
			position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
			return false
		},
		/* 114 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action90)> */
		func() bool {
			position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
			{
				position1238 := position
				depth++
				{
					position1239 := position
					depth++
					{
						position1240, tokenIndex1240, depth1240 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1241
						}
						position++
						goto l1240
					l1241:
						position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
						if buffer[position] != rune('N') {
							goto l1237
						}
						position++
					}
				l1240:
					{
						position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1243
						}
						position++
						goto l1242
					l1243:
						position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
						if buffer[position] != rune('O') {
							goto l1237
						}
						position++
					}
				l1242:
					{
						position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1245
						}
						position++
						goto l1244
					l1245:
						position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
						if buffer[position] != rune('T') {
							goto l1237
						}
						position++
					}
				l1244:
					depth--
					add(rulePegText, position1239)
				}
				if !_rules[ruleAction90]() {
					goto l1237
				}
				depth--
				add(ruleNot, position1238)
			}
			return true
		l1237:
			position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
			return false
		},
		/* 115 Equal <- <(<'='> Action91)> */
		func() bool {
			position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
			{
				position1247 := position
				depth++
				{
					position1248 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1246
					}
					position++
					depth--
					add(rulePegText, position1248)
				}
				if !_rules[ruleAction91]() {
					goto l1246
				}
				depth--
				add(ruleEqual, position1247)
			}
			return true
		l1246:
			position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
			return false
		},
		/* 116 Less <- <(<'<'> Action92)> */
		func() bool {
			position1249, tokenIndex1249, depth1249 := position, tokenIndex, depth
			{
				position1250 := position
				depth++
				{
					position1251 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1249
					}
					position++
					depth--
					add(rulePegText, position1251)
				}
				if !_rules[ruleAction92]() {
					goto l1249
				}
				depth--
				add(ruleLess, position1250)
			}
			return true
		l1249:
			position, tokenIndex, depth = position1249, tokenIndex1249, depth1249
			return false
		},
		/* 117 LessOrEqual <- <(<('<' '=')> Action93)> */
		func() bool {
			position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
			{
				position1253 := position
				depth++
				{
					position1254 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1252
					}
					position++
					if buffer[position] != rune('=') {
						goto l1252
					}
					position++
					depth--
					add(rulePegText, position1254)
				}
				if !_rules[ruleAction93]() {
					goto l1252
				}
				depth--
				add(ruleLessOrEqual, position1253)
			}
			return true
		l1252:
			position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
			return false
		},
		/* 118 Greater <- <(<'>'> Action94)> */
		func() bool {
			position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
			{
				position1256 := position
				depth++
				{
					position1257 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1255
					}
					position++
					depth--
					add(rulePegText, position1257)
				}
				if !_rules[ruleAction94]() {
					goto l1255
				}
				depth--
				add(ruleGreater, position1256)
			}
			return true
		l1255:
			position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
			return false
		},
		/* 119 GreaterOrEqual <- <(<('>' '=')> Action95)> */
		func() bool {
			position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
			{
				position1259 := position
				depth++
				{
					position1260 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1258
					}
					position++
					if buffer[position] != rune('=') {
						goto l1258
					}
					position++
					depth--
					add(rulePegText, position1260)
				}
				if !_rules[ruleAction95]() {
					goto l1258
				}
				depth--
				add(ruleGreaterOrEqual, position1259)
			}
			return true
		l1258:
			position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
			return false
		},
		/* 120 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action96)> */
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
						if buffer[position] != rune('!') {
							goto l1265
						}
						position++
						if buffer[position] != rune('=') {
							goto l1265
						}
						position++
						goto l1264
					l1265:
						position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
						if buffer[position] != rune('<') {
							goto l1261
						}
						position++
						if buffer[position] != rune('>') {
							goto l1261
						}
						position++
					}
				l1264:
					depth--
					add(rulePegText, position1263)
				}
				if !_rules[ruleAction96]() {
					goto l1261
				}
				depth--
				add(ruleNotEqual, position1262)
			}
			return true
		l1261:
			position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
			return false
		},
		/* 121 Concat <- <(<('|' '|')> Action97)> */
		func() bool {
			position1266, tokenIndex1266, depth1266 := position, tokenIndex, depth
			{
				position1267 := position
				depth++
				{
					position1268 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1266
					}
					position++
					if buffer[position] != rune('|') {
						goto l1266
					}
					position++
					depth--
					add(rulePegText, position1268)
				}
				if !_rules[ruleAction97]() {
					goto l1266
				}
				depth--
				add(ruleConcat, position1267)
			}
			return true
		l1266:
			position, tokenIndex, depth = position1266, tokenIndex1266, depth1266
			return false
		},
		/* 122 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action98)> */
		func() bool {
			position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
			{
				position1270 := position
				depth++
				{
					position1271 := position
					depth++
					{
						position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1273
						}
						position++
						goto l1272
					l1273:
						position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
						if buffer[position] != rune('I') {
							goto l1269
						}
						position++
					}
				l1272:
					{
						position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1275
						}
						position++
						goto l1274
					l1275:
						position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
						if buffer[position] != rune('S') {
							goto l1269
						}
						position++
					}
				l1274:
					depth--
					add(rulePegText, position1271)
				}
				if !_rules[ruleAction98]() {
					goto l1269
				}
				depth--
				add(ruleIs, position1270)
			}
			return true
		l1269:
			position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
			return false
		},
		/* 123 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action99)> */
		func() bool {
			position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
			{
				position1277 := position
				depth++
				{
					position1278 := position
					depth++
					{
						position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1280
						}
						position++
						goto l1279
					l1280:
						position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
						if buffer[position] != rune('I') {
							goto l1276
						}
						position++
					}
				l1279:
					{
						position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1282
						}
						position++
						goto l1281
					l1282:
						position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
						if buffer[position] != rune('S') {
							goto l1276
						}
						position++
					}
				l1281:
					if !_rules[rulesp]() {
						goto l1276
					}
					{
						position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1284
						}
						position++
						goto l1283
					l1284:
						position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
						if buffer[position] != rune('N') {
							goto l1276
						}
						position++
					}
				l1283:
					{
						position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1286
						}
						position++
						goto l1285
					l1286:
						position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
						if buffer[position] != rune('O') {
							goto l1276
						}
						position++
					}
				l1285:
					{
						position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1288
						}
						position++
						goto l1287
					l1288:
						position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
						if buffer[position] != rune('T') {
							goto l1276
						}
						position++
					}
				l1287:
					depth--
					add(rulePegText, position1278)
				}
				if !_rules[ruleAction99]() {
					goto l1276
				}
				depth--
				add(ruleIsNot, position1277)
			}
			return true
		l1276:
			position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
			return false
		},
		/* 124 Plus <- <(<'+'> Action100)> */
		func() bool {
			position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
			{
				position1290 := position
				depth++
				{
					position1291 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1289
					}
					position++
					depth--
					add(rulePegText, position1291)
				}
				if !_rules[ruleAction100]() {
					goto l1289
				}
				depth--
				add(rulePlus, position1290)
			}
			return true
		l1289:
			position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
			return false
		},
		/* 125 Minus <- <(<'-'> Action101)> */
		func() bool {
			position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
			{
				position1293 := position
				depth++
				{
					position1294 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1292
					}
					position++
					depth--
					add(rulePegText, position1294)
				}
				if !_rules[ruleAction101]() {
					goto l1292
				}
				depth--
				add(ruleMinus, position1293)
			}
			return true
		l1292:
			position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
			return false
		},
		/* 126 Multiply <- <(<'*'> Action102)> */
		func() bool {
			position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
			{
				position1296 := position
				depth++
				{
					position1297 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1295
					}
					position++
					depth--
					add(rulePegText, position1297)
				}
				if !_rules[ruleAction102]() {
					goto l1295
				}
				depth--
				add(ruleMultiply, position1296)
			}
			return true
		l1295:
			position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
			return false
		},
		/* 127 Divide <- <(<'/'> Action103)> */
		func() bool {
			position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
			{
				position1299 := position
				depth++
				{
					position1300 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1298
					}
					position++
					depth--
					add(rulePegText, position1300)
				}
				if !_rules[ruleAction103]() {
					goto l1298
				}
				depth--
				add(ruleDivide, position1299)
			}
			return true
		l1298:
			position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
			return false
		},
		/* 128 Modulo <- <(<'%'> Action104)> */
		func() bool {
			position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
			{
				position1302 := position
				depth++
				{
					position1303 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1301
					}
					position++
					depth--
					add(rulePegText, position1303)
				}
				if !_rules[ruleAction104]() {
					goto l1301
				}
				depth--
				add(ruleModulo, position1302)
			}
			return true
		l1301:
			position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
			return false
		},
		/* 129 UnaryMinus <- <(<'-'> Action105)> */
		func() bool {
			position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
			{
				position1305 := position
				depth++
				{
					position1306 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1304
					}
					position++
					depth--
					add(rulePegText, position1306)
				}
				if !_rules[ruleAction105]() {
					goto l1304
				}
				depth--
				add(ruleUnaryMinus, position1305)
			}
			return true
		l1304:
			position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
			return false
		},
		/* 130 Identifier <- <(<ident> Action106)> */
		func() bool {
			position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
			{
				position1308 := position
				depth++
				{
					position1309 := position
					depth++
					if !_rules[ruleident]() {
						goto l1307
					}
					depth--
					add(rulePegText, position1309)
				}
				if !_rules[ruleAction106]() {
					goto l1307
				}
				depth--
				add(ruleIdentifier, position1308)
			}
			return true
		l1307:
			position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
			return false
		},
		/* 131 TargetIdentifier <- <(<jsonPath> Action107)> */
		func() bool {
			position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
			{
				position1311 := position
				depth++
				{
					position1312 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1310
					}
					depth--
					add(rulePegText, position1312)
				}
				if !_rules[ruleAction107]() {
					goto l1310
				}
				depth--
				add(ruleTargetIdentifier, position1311)
			}
			return true
		l1310:
			position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
			return false
		},
		/* 132 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
			{
				position1314 := position
				depth++
				{
					position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1316
					}
					position++
					goto l1315
				l1316:
					position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1313
					}
					position++
				}
			l1315:
			l1317:
				{
					position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
					{
						position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1320
						}
						position++
						goto l1319
					l1320:
						position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1321
						}
						position++
						goto l1319
					l1321:
						position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1322
						}
						position++
						goto l1319
					l1322:
						position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
						if buffer[position] != rune('_') {
							goto l1318
						}
						position++
					}
				l1319:
					goto l1317
				l1318:
					position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
				}
				depth--
				add(ruleident, position1314)
			}
			return true
		l1313:
			position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
			return false
		},
		/* 133 jsonPath <- <(jsonPathHead jsonPathNonHead*)> */
		func() bool {
			position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
			{
				position1324 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1323
				}
			l1325:
				{
					position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
					if !_rules[rulejsonPathNonHead]() {
						goto l1326
					}
					goto l1325
				l1326:
					position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
				}
				depth--
				add(rulejsonPath, position1324)
			}
			return true
		l1323:
			position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
			return false
		},
		/* 134 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
			{
				position1328 := position
				depth++
				{
					position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1330
					}
					goto l1329
				l1330:
					position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1327
					}
				}
			l1329:
				depth--
				add(rulejsonPathHead, position1328)
			}
			return true
		l1327:
			position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
			return false
		},
		/* 135 jsonPathNonHead <- <(('.' jsonMapAccessString) / jsonMapAccessBracket / jsonArrayAccess)> */
		func() bool {
			position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
			{
				position1332 := position
				depth++
				{
					position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1334
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1334
					}
					goto l1333
				l1334:
					position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1335
					}
					goto l1333
				l1335:
					position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
					if !_rules[rulejsonArrayAccess]() {
						goto l1331
					}
				}
			l1333:
				depth--
				add(rulejsonPathNonHead, position1332)
			}
			return true
		l1331:
			position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
			return false
		},
		/* 136 jsonMapAccessString <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1336, tokenIndex1336, depth1336 := position, tokenIndex, depth
			{
				position1337 := position
				depth++
				{
					position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1339
					}
					position++
					goto l1338
				l1339:
					position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1336
					}
					position++
				}
			l1338:
			l1340:
				{
					position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
					{
						position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1343
						}
						position++
						goto l1342
					l1343:
						position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1344
						}
						position++
						goto l1342
					l1344:
						position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1345
						}
						position++
						goto l1342
					l1345:
						position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
						if buffer[position] != rune('_') {
							goto l1341
						}
						position++
					}
				l1342:
					goto l1340
				l1341:
					position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
				}
				depth--
				add(rulejsonMapAccessString, position1337)
			}
			return true
		l1336:
			position, tokenIndex, depth = position1336, tokenIndex1336, depth1336
			return false
		},
		/* 137 jsonMapAccessBracket <- <('[' '\'' (('\'' '\'') / (!'\'' .))* '\'' ']')> */
		func() bool {
			position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
			{
				position1347 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1346
				}
				position++
				if buffer[position] != rune('\'') {
					goto l1346
				}
				position++
			l1348:
				{
					position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
					{
						position1350, tokenIndex1350, depth1350 := position, tokenIndex, depth
						if buffer[position] != rune('\'') {
							goto l1351
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1351
						}
						position++
						goto l1350
					l1351:
						position, tokenIndex, depth = position1350, tokenIndex1350, depth1350
						{
							position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1352
							}
							position++
							goto l1349
						l1352:
							position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
						}
						if !matchDot() {
							goto l1349
						}
					}
				l1350:
					goto l1348
				l1349:
					position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
				}
				if buffer[position] != rune('\'') {
					goto l1346
				}
				position++
				if buffer[position] != rune(']') {
					goto l1346
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1347)
			}
			return true
		l1346:
			position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
			return false
		},
		/* 138 jsonArrayAccess <- <('[' [0-9]+ ']')> */
		func() bool {
			position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
			{
				position1354 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1353
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1353
				}
				position++
			l1355:
				{
					position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1356
					}
					position++
					goto l1355
				l1356:
					position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
				}
				if buffer[position] != rune(']') {
					goto l1353
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1354)
			}
			return true
		l1353:
			position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
			return false
		},
		/* 139 sp <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)*> */
		func() bool {
			{
				position1358 := position
				depth++
			l1359:
				{
					position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
					{
						position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1362
						}
						position++
						goto l1361
					l1362:
						position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
						if buffer[position] != rune('\t') {
							goto l1363
						}
						position++
						goto l1361
					l1363:
						position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
						if buffer[position] != rune('\n') {
							goto l1364
						}
						position++
						goto l1361
					l1364:
						position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
						if buffer[position] != rune('\r') {
							goto l1365
						}
						position++
						goto l1361
					l1365:
						position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
						if !_rules[rulecomment]() {
							goto l1366
						}
						goto l1361
					l1366:
						position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
						if !_rules[rulefinalComment]() {
							goto l1360
						}
					}
				l1361:
					goto l1359
				l1360:
					position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
				}
				depth--
				add(rulesp, position1358)
			}
			return true
		},
		/* 140 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
			{
				position1368 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1367
				}
				position++
				if buffer[position] != rune('-') {
					goto l1367
				}
				position++
			l1369:
				{
					position1370, tokenIndex1370, depth1370 := position, tokenIndex, depth
					{
						position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
						{
							position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1373
							}
							position++
							goto l1372
						l1373:
							position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
							if buffer[position] != rune('\n') {
								goto l1371
							}
							position++
						}
					l1372:
						goto l1370
					l1371:
						position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
					}
					if !matchDot() {
						goto l1370
					}
					goto l1369
				l1370:
					position, tokenIndex, depth = position1370, tokenIndex1370, depth1370
				}
				{
					position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1375
					}
					position++
					goto l1374
				l1375:
					position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
					if buffer[position] != rune('\n') {
						goto l1367
					}
					position++
				}
			l1374:
				depth--
				add(rulecomment, position1368)
			}
			return true
		l1367:
			position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
			return false
		},
		/* 141 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
			{
				position1377 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1376
				}
				position++
				if buffer[position] != rune('-') {
					goto l1376
				}
				position++
			l1378:
				{
					position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
					{
						position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
						{
							position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1382
							}
							position++
							goto l1381
						l1382:
							position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
							if buffer[position] != rune('\n') {
								goto l1380
							}
							position++
						}
					l1381:
						goto l1379
					l1380:
						position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
					}
					if !matchDot() {
						goto l1379
					}
					goto l1378
				l1379:
					position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
				}
				{
					position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
					if !matchDot() {
						goto l1383
					}
					goto l1376
				l1383:
					position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
				}
				depth--
				add(rulefinalComment, position1377)
			}
			return true
		l1376:
			position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
			return false
		},
		nil,
		/* 144 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 145 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 146 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 147 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 148 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 149 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 150 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 151 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 152 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 153 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 154 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 155 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 156 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 157 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 158 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 159 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 160 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 161 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 162 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 163 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 164 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 165 Action21 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 166 Action22 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 167 Action23 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 168 Action24 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 169 Action25 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 170 Action26 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 171 Action27 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 172 Action28 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 173 Action29 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 174 Action30 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 175 Action31 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 176 Action32 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 177 Action33 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 178 Action34 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 179 Action35 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 180 Action36 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 181 Action37 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 182 Action38 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 183 Action39 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 184 Action40 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 185 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 186 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 187 Action43 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 188 Action44 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 189 Action45 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 190 Action46 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 191 Action47 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 192 Action48 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 193 Action49 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 194 Action50 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 195 Action51 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 196 Action52 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 197 Action53 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 198 Action54 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 199 Action55 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 200 Action56 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 201 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 202 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 203 Action59 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 204 Action60 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 205 Action61 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 206 Action62 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 207 Action63 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 208 Action64 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 209 Action65 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 210 Action66 <- <{
		    p.PushComponent(begin, end, NewWildcard(""))
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 211 Action67 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 212 Action68 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 213 Action69 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 214 Action70 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 215 Action71 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 216 Action72 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 217 Action73 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 218 Action74 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 219 Action75 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 220 Action76 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 221 Action77 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 222 Action78 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 223 Action79 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 224 Action80 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 225 Action81 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 226 Action82 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 227 Action83 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 228 Action84 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 229 Action85 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 230 Action86 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 231 Action87 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 232 Action88 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 233 Action89 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 234 Action90 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 235 Action91 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 236 Action92 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 237 Action93 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 238 Action94 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 239 Action95 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 240 Action96 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 241 Action97 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 242 Action98 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 243 Action99 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 244 Action100 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 245 Action101 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 246 Action102 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 247 Action103 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 248 Action104 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 249 Action105 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 250 Action106 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 251 Action107 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
	}
	p.rules = _rules
}

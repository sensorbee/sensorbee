package parser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

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
	rules  [248]func() bool
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
	for i, c := range buffer[0:] {
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

			p.AssembleProjections(begin, end)

		case ruleAction23:

			p.AssembleAlias()

		case ruleAction24:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction25:

			p.AssembleInterval()

		case ruleAction26:

			p.AssembleInterval()

		case ruleAction27:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction28:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction29:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction30:

			p.EnsureAliasedStreamWindow()

		case ruleAction31:

			p.AssembleAliasedStreamWindow()

		case ruleAction32:

			p.AssembleStreamWindow()

		case ruleAction33:

			p.AssembleUDSFFuncApp()

		case ruleAction34:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction35:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction36:

			p.AssembleSourceSinkParam()

		case ruleAction37:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction38:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction39:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction40:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction41:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction42:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction43:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction44:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction45:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction46:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction47:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction48:

			p.AssembleTypeCast(begin, end)

		case ruleAction49:

			p.AssembleTypeCast(begin, end)

		case ruleAction50:

			p.AssembleFuncApp()

		case ruleAction51:

			p.AssembleExpressions(begin, end)

		case ruleAction52:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction53:

			p.AssembleMap(begin, end)

		case ruleAction54:

			p.AssembleKeyValuePair()

		case ruleAction55:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction56:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction57:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction58:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction59:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction60:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction61:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction62:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction63:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction64:

			p.PushComponent(begin, end, NewWildcard(""))

		case ruleAction65:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction66:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction67:

			p.PushComponent(begin, end, Istream)

		case ruleAction68:

			p.PushComponent(begin, end, Dstream)

		case ruleAction69:

			p.PushComponent(begin, end, Rstream)

		case ruleAction70:

			p.PushComponent(begin, end, Tuples)

		case ruleAction71:

			p.PushComponent(begin, end, Seconds)

		case ruleAction72:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction73:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction74:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction75:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction76:

			p.PushComponent(begin, end, Yes)

		case ruleAction77:

			p.PushComponent(begin, end, No)

		case ruleAction78:

			p.PushComponent(begin, end, Bool)

		case ruleAction79:

			p.PushComponent(begin, end, Int)

		case ruleAction80:

			p.PushComponent(begin, end, Float)

		case ruleAction81:

			p.PushComponent(begin, end, String)

		case ruleAction82:

			p.PushComponent(begin, end, Blob)

		case ruleAction83:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction84:

			p.PushComponent(begin, end, Array)

		case ruleAction85:

			p.PushComponent(begin, end, Map)

		case ruleAction86:

			p.PushComponent(begin, end, Or)

		case ruleAction87:

			p.PushComponent(begin, end, And)

		case ruleAction88:

			p.PushComponent(begin, end, Not)

		case ruleAction89:

			p.PushComponent(begin, end, Equal)

		case ruleAction90:

			p.PushComponent(begin, end, Less)

		case ruleAction91:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction92:

			p.PushComponent(begin, end, Greater)

		case ruleAction93:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction94:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction95:

			p.PushComponent(begin, end, Concat)

		case ruleAction96:

			p.PushComponent(begin, end, Is)

		case ruleAction97:

			p.PushComponent(begin, end, IsNot)

		case ruleAction98:

			p.PushComponent(begin, end, Plus)

		case ruleAction99:

			p.PushComponent(begin, end, Minus)

		case ruleAction100:

			p.PushComponent(begin, end, Multiply)

		case ruleAction101:

			p.PushComponent(begin, end, Divide)

		case ruleAction102:

			p.PushComponent(begin, end, Modulo)

		case ruleAction103:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction104:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction105:

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
		/* 27 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) Action21)> */
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
		/* 28 Projections <- <(<(Projection sp (',' sp Projection)*)> Action22)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l540
					}
					if !_rules[rulesp]() {
						goto l540
					}
				l543:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l544
						}
						position++
						if !_rules[rulesp]() {
							goto l544
						}
						if !_rules[ruleProjection]() {
							goto l544
						}
						goto l543
					l544:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
					}
					depth--
					add(rulePegText, position542)
				}
				if !_rules[ruleAction22]() {
					goto l540
				}
				depth--
				add(ruleProjections, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 29 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				{
					position547, tokenIndex547, depth547 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l548
					}
					goto l547
				l548:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
					if !_rules[ruleExpression]() {
						goto l549
					}
					goto l547
				l549:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
					if !_rules[ruleWildcard]() {
						goto l545
					}
				}
			l547:
				depth--
				add(ruleProjection, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 30 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action23)> */
		func() bool {
			position550, tokenIndex550, depth550 := position, tokenIndex, depth
			{
				position551 := position
				depth++
				{
					position552, tokenIndex552, depth552 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l553
					}
					goto l552
				l553:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if !_rules[ruleWildcard]() {
						goto l550
					}
				}
			l552:
				if !_rules[rulesp]() {
					goto l550
				}
				{
					position554, tokenIndex554, depth554 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l555
					}
					position++
					goto l554
				l555:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
					if buffer[position] != rune('A') {
						goto l550
					}
					position++
				}
			l554:
				{
					position556, tokenIndex556, depth556 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l557
					}
					position++
					goto l556
				l557:
					position, tokenIndex, depth = position556, tokenIndex556, depth556
					if buffer[position] != rune('S') {
						goto l550
					}
					position++
				}
			l556:
				if !_rules[rulesp]() {
					goto l550
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l550
				}
				if !_rules[ruleAction23]() {
					goto l550
				}
				depth--
				add(ruleAliasExpression, position551)
			}
			return true
		l550:
			position, tokenIndex, depth = position550, tokenIndex550, depth550
			return false
		},
		/* 31 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action24)> */
		func() bool {
			position558, tokenIndex558, depth558 := position, tokenIndex, depth
			{
				position559 := position
				depth++
				{
					position560 := position
					depth++
					{
						position561, tokenIndex561, depth561 := position, tokenIndex, depth
						{
							position563, tokenIndex563, depth563 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l564
							}
							position++
							goto l563
						l564:
							position, tokenIndex, depth = position563, tokenIndex563, depth563
							if buffer[position] != rune('F') {
								goto l561
							}
							position++
						}
					l563:
						{
							position565, tokenIndex565, depth565 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l566
							}
							position++
							goto l565
						l566:
							position, tokenIndex, depth = position565, tokenIndex565, depth565
							if buffer[position] != rune('R') {
								goto l561
							}
							position++
						}
					l565:
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
								goto l561
							}
							position++
						}
					l567:
						{
							position569, tokenIndex569, depth569 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l570
							}
							position++
							goto l569
						l570:
							position, tokenIndex, depth = position569, tokenIndex569, depth569
							if buffer[position] != rune('M') {
								goto l561
							}
							position++
						}
					l569:
						if !_rules[rulesp]() {
							goto l561
						}
						if !_rules[ruleRelations]() {
							goto l561
						}
						if !_rules[rulesp]() {
							goto l561
						}
						goto l562
					l561:
						position, tokenIndex, depth = position561, tokenIndex561, depth561
					}
				l562:
					depth--
					add(rulePegText, position560)
				}
				if !_rules[ruleAction24]() {
					goto l558
				}
				depth--
				add(ruleWindowedFrom, position559)
			}
			return true
		l558:
			position, tokenIndex, depth = position558, tokenIndex558, depth558
			return false
		},
		/* 32 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position571, tokenIndex571, depth571 := position, tokenIndex, depth
			{
				position572 := position
				depth++
				{
					position573, tokenIndex573, depth573 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l574
					}
					goto l573
				l574:
					position, tokenIndex, depth = position573, tokenIndex573, depth573
					if !_rules[ruleTuplesInterval]() {
						goto l571
					}
				}
			l573:
				depth--
				add(ruleInterval, position572)
			}
			return true
		l571:
			position, tokenIndex, depth = position571, tokenIndex571, depth571
			return false
		},
		/* 33 TimeInterval <- <(NumericLiteral sp (SECONDS / MILLISECONDS) Action25)> */
		func() bool {
			position575, tokenIndex575, depth575 := position, tokenIndex, depth
			{
				position576 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l575
				}
				if !_rules[rulesp]() {
					goto l575
				}
				{
					position577, tokenIndex577, depth577 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l578
					}
					goto l577
				l578:
					position, tokenIndex, depth = position577, tokenIndex577, depth577
					if !_rules[ruleMILLISECONDS]() {
						goto l575
					}
				}
			l577:
				if !_rules[ruleAction25]() {
					goto l575
				}
				depth--
				add(ruleTimeInterval, position576)
			}
			return true
		l575:
			position, tokenIndex, depth = position575, tokenIndex575, depth575
			return false
		},
		/* 34 TuplesInterval <- <(NumericLiteral sp TUPLES Action26)> */
		func() bool {
			position579, tokenIndex579, depth579 := position, tokenIndex, depth
			{
				position580 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l579
				}
				if !_rules[rulesp]() {
					goto l579
				}
				if !_rules[ruleTUPLES]() {
					goto l579
				}
				if !_rules[ruleAction26]() {
					goto l579
				}
				depth--
				add(ruleTuplesInterval, position580)
			}
			return true
		l579:
			position, tokenIndex, depth = position579, tokenIndex579, depth579
			return false
		},
		/* 35 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position581, tokenIndex581, depth581 := position, tokenIndex, depth
			{
				position582 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l581
				}
				if !_rules[rulesp]() {
					goto l581
				}
			l583:
				{
					position584, tokenIndex584, depth584 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l584
					}
					position++
					if !_rules[rulesp]() {
						goto l584
					}
					if !_rules[ruleRelationLike]() {
						goto l584
					}
					goto l583
				l584:
					position, tokenIndex, depth = position584, tokenIndex584, depth584
				}
				depth--
				add(ruleRelations, position582)
			}
			return true
		l581:
			position, tokenIndex, depth = position581, tokenIndex581, depth581
			return false
		},
		/* 36 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action27)> */
		func() bool {
			position585, tokenIndex585, depth585 := position, tokenIndex, depth
			{
				position586 := position
				depth++
				{
					position587 := position
					depth++
					{
						position588, tokenIndex588, depth588 := position, tokenIndex, depth
						{
							position590, tokenIndex590, depth590 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l591
							}
							position++
							goto l590
						l591:
							position, tokenIndex, depth = position590, tokenIndex590, depth590
							if buffer[position] != rune('W') {
								goto l588
							}
							position++
						}
					l590:
						{
							position592, tokenIndex592, depth592 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l593
							}
							position++
							goto l592
						l593:
							position, tokenIndex, depth = position592, tokenIndex592, depth592
							if buffer[position] != rune('H') {
								goto l588
							}
							position++
						}
					l592:
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
								goto l588
							}
							position++
						}
					l594:
						{
							position596, tokenIndex596, depth596 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l597
							}
							position++
							goto l596
						l597:
							position, tokenIndex, depth = position596, tokenIndex596, depth596
							if buffer[position] != rune('R') {
								goto l588
							}
							position++
						}
					l596:
						{
							position598, tokenIndex598, depth598 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l599
							}
							position++
							goto l598
						l599:
							position, tokenIndex, depth = position598, tokenIndex598, depth598
							if buffer[position] != rune('E') {
								goto l588
							}
							position++
						}
					l598:
						if !_rules[rulesp]() {
							goto l588
						}
						if !_rules[ruleExpression]() {
							goto l588
						}
						goto l589
					l588:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
					}
				l589:
					depth--
					add(rulePegText, position587)
				}
				if !_rules[ruleAction27]() {
					goto l585
				}
				depth--
				add(ruleFilter, position586)
			}
			return true
		l585:
			position, tokenIndex, depth = position585, tokenIndex585, depth585
			return false
		},
		/* 37 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action28)> */
		func() bool {
			position600, tokenIndex600, depth600 := position, tokenIndex, depth
			{
				position601 := position
				depth++
				{
					position602 := position
					depth++
					{
						position603, tokenIndex603, depth603 := position, tokenIndex, depth
						{
							position605, tokenIndex605, depth605 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l606
							}
							position++
							goto l605
						l606:
							position, tokenIndex, depth = position605, tokenIndex605, depth605
							if buffer[position] != rune('G') {
								goto l603
							}
							position++
						}
					l605:
						{
							position607, tokenIndex607, depth607 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l608
							}
							position++
							goto l607
						l608:
							position, tokenIndex, depth = position607, tokenIndex607, depth607
							if buffer[position] != rune('R') {
								goto l603
							}
							position++
						}
					l607:
						{
							position609, tokenIndex609, depth609 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l610
							}
							position++
							goto l609
						l610:
							position, tokenIndex, depth = position609, tokenIndex609, depth609
							if buffer[position] != rune('O') {
								goto l603
							}
							position++
						}
					l609:
						{
							position611, tokenIndex611, depth611 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l612
							}
							position++
							goto l611
						l612:
							position, tokenIndex, depth = position611, tokenIndex611, depth611
							if buffer[position] != rune('U') {
								goto l603
							}
							position++
						}
					l611:
						{
							position613, tokenIndex613, depth613 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l614
							}
							position++
							goto l613
						l614:
							position, tokenIndex, depth = position613, tokenIndex613, depth613
							if buffer[position] != rune('P') {
								goto l603
							}
							position++
						}
					l613:
						if !_rules[rulesp]() {
							goto l603
						}
						{
							position615, tokenIndex615, depth615 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l616
							}
							position++
							goto l615
						l616:
							position, tokenIndex, depth = position615, tokenIndex615, depth615
							if buffer[position] != rune('B') {
								goto l603
							}
							position++
						}
					l615:
						{
							position617, tokenIndex617, depth617 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l618
							}
							position++
							goto l617
						l618:
							position, tokenIndex, depth = position617, tokenIndex617, depth617
							if buffer[position] != rune('Y') {
								goto l603
							}
							position++
						}
					l617:
						if !_rules[rulesp]() {
							goto l603
						}
						if !_rules[ruleGroupList]() {
							goto l603
						}
						goto l604
					l603:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
					}
				l604:
					depth--
					add(rulePegText, position602)
				}
				if !_rules[ruleAction28]() {
					goto l600
				}
				depth--
				add(ruleGrouping, position601)
			}
			return true
		l600:
			position, tokenIndex, depth = position600, tokenIndex600, depth600
			return false
		},
		/* 38 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position619, tokenIndex619, depth619 := position, tokenIndex, depth
			{
				position620 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l619
				}
				if !_rules[rulesp]() {
					goto l619
				}
			l621:
				{
					position622, tokenIndex622, depth622 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l622
					}
					position++
					if !_rules[rulesp]() {
						goto l622
					}
					if !_rules[ruleExpression]() {
						goto l622
					}
					goto l621
				l622:
					position, tokenIndex, depth = position622, tokenIndex622, depth622
				}
				depth--
				add(ruleGroupList, position620)
			}
			return true
		l619:
			position, tokenIndex, depth = position619, tokenIndex619, depth619
			return false
		},
		/* 39 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action29)> */
		func() bool {
			position623, tokenIndex623, depth623 := position, tokenIndex, depth
			{
				position624 := position
				depth++
				{
					position625 := position
					depth++
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						{
							position628, tokenIndex628, depth628 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l629
							}
							position++
							goto l628
						l629:
							position, tokenIndex, depth = position628, tokenIndex628, depth628
							if buffer[position] != rune('H') {
								goto l626
							}
							position++
						}
					l628:
						{
							position630, tokenIndex630, depth630 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l631
							}
							position++
							goto l630
						l631:
							position, tokenIndex, depth = position630, tokenIndex630, depth630
							if buffer[position] != rune('A') {
								goto l626
							}
							position++
						}
					l630:
						{
							position632, tokenIndex632, depth632 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l633
							}
							position++
							goto l632
						l633:
							position, tokenIndex, depth = position632, tokenIndex632, depth632
							if buffer[position] != rune('V') {
								goto l626
							}
							position++
						}
					l632:
						{
							position634, tokenIndex634, depth634 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l635
							}
							position++
							goto l634
						l635:
							position, tokenIndex, depth = position634, tokenIndex634, depth634
							if buffer[position] != rune('I') {
								goto l626
							}
							position++
						}
					l634:
						{
							position636, tokenIndex636, depth636 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l637
							}
							position++
							goto l636
						l637:
							position, tokenIndex, depth = position636, tokenIndex636, depth636
							if buffer[position] != rune('N') {
								goto l626
							}
							position++
						}
					l636:
						{
							position638, tokenIndex638, depth638 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l639
							}
							position++
							goto l638
						l639:
							position, tokenIndex, depth = position638, tokenIndex638, depth638
							if buffer[position] != rune('G') {
								goto l626
							}
							position++
						}
					l638:
						if !_rules[rulesp]() {
							goto l626
						}
						if !_rules[ruleExpression]() {
							goto l626
						}
						goto l627
					l626:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
					}
				l627:
					depth--
					add(rulePegText, position625)
				}
				if !_rules[ruleAction29]() {
					goto l623
				}
				depth--
				add(ruleHaving, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 40 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action30))> */
		func() bool {
			position640, tokenIndex640, depth640 := position, tokenIndex, depth
			{
				position641 := position
				depth++
				{
					position642, tokenIndex642, depth642 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l643
					}
					goto l642
				l643:
					position, tokenIndex, depth = position642, tokenIndex642, depth642
					if !_rules[ruleStreamWindow]() {
						goto l640
					}
					if !_rules[ruleAction30]() {
						goto l640
					}
				}
			l642:
				depth--
				add(ruleRelationLike, position641)
			}
			return true
		l640:
			position, tokenIndex, depth = position640, tokenIndex640, depth640
			return false
		},
		/* 41 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action31)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l644
				}
				if !_rules[rulesp]() {
					goto l644
				}
				{
					position646, tokenIndex646, depth646 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l647
					}
					position++
					goto l646
				l647:
					position, tokenIndex, depth = position646, tokenIndex646, depth646
					if buffer[position] != rune('A') {
						goto l644
					}
					position++
				}
			l646:
				{
					position648, tokenIndex648, depth648 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l649
					}
					position++
					goto l648
				l649:
					position, tokenIndex, depth = position648, tokenIndex648, depth648
					if buffer[position] != rune('S') {
						goto l644
					}
					position++
				}
			l648:
				if !_rules[rulesp]() {
					goto l644
				}
				if !_rules[ruleIdentifier]() {
					goto l644
				}
				if !_rules[ruleAction31]() {
					goto l644
				}
				depth--
				add(ruleAliasedStreamWindow, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 42 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action32)> */
		func() bool {
			position650, tokenIndex650, depth650 := position, tokenIndex, depth
			{
				position651 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l650
				}
				if !_rules[rulesp]() {
					goto l650
				}
				if buffer[position] != rune('[') {
					goto l650
				}
				position++
				if !_rules[rulesp]() {
					goto l650
				}
				{
					position652, tokenIndex652, depth652 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l653
					}
					position++
					goto l652
				l653:
					position, tokenIndex, depth = position652, tokenIndex652, depth652
					if buffer[position] != rune('R') {
						goto l650
					}
					position++
				}
			l652:
				{
					position654, tokenIndex654, depth654 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l655
					}
					position++
					goto l654
				l655:
					position, tokenIndex, depth = position654, tokenIndex654, depth654
					if buffer[position] != rune('A') {
						goto l650
					}
					position++
				}
			l654:
				{
					position656, tokenIndex656, depth656 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l657
					}
					position++
					goto l656
				l657:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if buffer[position] != rune('N') {
						goto l650
					}
					position++
				}
			l656:
				{
					position658, tokenIndex658, depth658 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l659
					}
					position++
					goto l658
				l659:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
					if buffer[position] != rune('G') {
						goto l650
					}
					position++
				}
			l658:
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
						goto l650
					}
					position++
				}
			l660:
				if !_rules[rulesp]() {
					goto l650
				}
				if !_rules[ruleInterval]() {
					goto l650
				}
				if !_rules[rulesp]() {
					goto l650
				}
				if buffer[position] != rune(']') {
					goto l650
				}
				position++
				if !_rules[ruleAction32]() {
					goto l650
				}
				depth--
				add(ruleStreamWindow, position651)
			}
			return true
		l650:
			position, tokenIndex, depth = position650, tokenIndex650, depth650
			return false
		},
		/* 43 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position662, tokenIndex662, depth662 := position, tokenIndex, depth
			{
				position663 := position
				depth++
				{
					position664, tokenIndex664, depth664 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l665
					}
					goto l664
				l665:
					position, tokenIndex, depth = position664, tokenIndex664, depth664
					if !_rules[ruleStream]() {
						goto l662
					}
				}
			l664:
				depth--
				add(ruleStreamLike, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 44 UDSFFuncApp <- <(FuncApp Action33)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l666
				}
				if !_rules[ruleAction33]() {
					goto l666
				}
				depth--
				add(ruleUDSFFuncApp, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 45 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action34)> */
		func() bool {
			position668, tokenIndex668, depth668 := position, tokenIndex, depth
			{
				position669 := position
				depth++
				{
					position670 := position
					depth++
					{
						position671, tokenIndex671, depth671 := position, tokenIndex, depth
						{
							position673, tokenIndex673, depth673 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l674
							}
							position++
							goto l673
						l674:
							position, tokenIndex, depth = position673, tokenIndex673, depth673
							if buffer[position] != rune('W') {
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
							if buffer[position] != rune('t') {
								goto l678
							}
							position++
							goto l677
						l678:
							position, tokenIndex, depth = position677, tokenIndex677, depth677
							if buffer[position] != rune('T') {
								goto l671
							}
							position++
						}
					l677:
						{
							position679, tokenIndex679, depth679 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l680
							}
							position++
							goto l679
						l680:
							position, tokenIndex, depth = position679, tokenIndex679, depth679
							if buffer[position] != rune('H') {
								goto l671
							}
							position++
						}
					l679:
						if !_rules[rulesp]() {
							goto l671
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l671
						}
						if !_rules[rulesp]() {
							goto l671
						}
					l681:
						{
							position682, tokenIndex682, depth682 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l682
							}
							position++
							if !_rules[rulesp]() {
								goto l682
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l682
							}
							goto l681
						l682:
							position, tokenIndex, depth = position682, tokenIndex682, depth682
						}
						goto l672
					l671:
						position, tokenIndex, depth = position671, tokenIndex671, depth671
					}
				l672:
					depth--
					add(rulePegText, position670)
				}
				if !_rules[ruleAction34]() {
					goto l668
				}
				depth--
				add(ruleSourceSinkSpecs, position669)
			}
			return true
		l668:
			position, tokenIndex, depth = position668, tokenIndex668, depth668
			return false
		},
		/* 46 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action35)> */
		func() bool {
			position683, tokenIndex683, depth683 := position, tokenIndex, depth
			{
				position684 := position
				depth++
				{
					position685 := position
					depth++
					{
						position686, tokenIndex686, depth686 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
						if buffer[position] != rune('S') {
							goto l683
						}
						position++
					}
				l686:
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if buffer[position] != rune('E') {
							goto l683
						}
						position++
					}
				l688:
					{
						position690, tokenIndex690, depth690 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l691
						}
						position++
						goto l690
					l691:
						position, tokenIndex, depth = position690, tokenIndex690, depth690
						if buffer[position] != rune('T') {
							goto l683
						}
						position++
					}
				l690:
					if !_rules[rulesp]() {
						goto l683
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l683
					}
					if !_rules[rulesp]() {
						goto l683
					}
				l692:
					{
						position693, tokenIndex693, depth693 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l693
						}
						position++
						if !_rules[rulesp]() {
							goto l693
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l693
						}
						goto l692
					l693:
						position, tokenIndex, depth = position693, tokenIndex693, depth693
					}
					depth--
					add(rulePegText, position685)
				}
				if !_rules[ruleAction35]() {
					goto l683
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position684)
			}
			return true
		l683:
			position, tokenIndex, depth = position683, tokenIndex683, depth683
			return false
		},
		/* 47 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action36)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l694
				}
				if buffer[position] != rune('=') {
					goto l694
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l694
				}
				if !_rules[ruleAction36]() {
					goto l694
				}
				depth--
				add(ruleSourceSinkParam, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 48 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr)> */
		func() bool {
			position696, tokenIndex696, depth696 := position, tokenIndex, depth
			{
				position697 := position
				depth++
				{
					position698, tokenIndex698, depth698 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l699
					}
					goto l698
				l699:
					position, tokenIndex, depth = position698, tokenIndex698, depth698
					if !_rules[ruleParamArrayExpr]() {
						goto l696
					}
				}
			l698:
				depth--
				add(ruleSourceSinkParamVal, position697)
			}
			return true
		l696:
			position, tokenIndex, depth = position696, tokenIndex696, depth696
			return false
		},
		/* 49 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				{
					position702, tokenIndex702, depth702 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l703
					}
					goto l702
				l703:
					position, tokenIndex, depth = position702, tokenIndex702, depth702
					if !_rules[ruleLiteral]() {
						goto l700
					}
				}
			l702:
				depth--
				add(ruleParamLiteral, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 50 ParamArrayExpr <- <(<('[' sp (ParamLiteral (',' sp ParamLiteral)*)? sp ','? sp ']')> Action37)> */
		func() bool {
			position704, tokenIndex704, depth704 := position, tokenIndex, depth
			{
				position705 := position
				depth++
				{
					position706 := position
					depth++
					if buffer[position] != rune('[') {
						goto l704
					}
					position++
					if !_rules[rulesp]() {
						goto l704
					}
					{
						position707, tokenIndex707, depth707 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l707
						}
					l709:
						{
							position710, tokenIndex710, depth710 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l710
							}
							position++
							if !_rules[rulesp]() {
								goto l710
							}
							if !_rules[ruleParamLiteral]() {
								goto l710
							}
							goto l709
						l710:
							position, tokenIndex, depth = position710, tokenIndex710, depth710
						}
						goto l708
					l707:
						position, tokenIndex, depth = position707, tokenIndex707, depth707
					}
				l708:
					if !_rules[rulesp]() {
						goto l704
					}
					{
						position711, tokenIndex711, depth711 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l711
						}
						position++
						goto l712
					l711:
						position, tokenIndex, depth = position711, tokenIndex711, depth711
					}
				l712:
					if !_rules[rulesp]() {
						goto l704
					}
					if buffer[position] != rune(']') {
						goto l704
					}
					position++
					depth--
					add(rulePegText, position706)
				}
				if !_rules[ruleAction37]() {
					goto l704
				}
				depth--
				add(ruleParamArrayExpr, position705)
			}
			return true
		l704:
			position, tokenIndex, depth = position704, tokenIndex704, depth704
			return false
		},
		/* 51 PausedOpt <- <(<(Paused / Unpaused)?> Action38)> */
		func() bool {
			position713, tokenIndex713, depth713 := position, tokenIndex, depth
			{
				position714 := position
				depth++
				{
					position715 := position
					depth++
					{
						position716, tokenIndex716, depth716 := position, tokenIndex, depth
						{
							position718, tokenIndex718, depth718 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l719
							}
							goto l718
						l719:
							position, tokenIndex, depth = position718, tokenIndex718, depth718
							if !_rules[ruleUnpaused]() {
								goto l716
							}
						}
					l718:
						goto l717
					l716:
						position, tokenIndex, depth = position716, tokenIndex716, depth716
					}
				l717:
					depth--
					add(rulePegText, position715)
				}
				if !_rules[ruleAction38]() {
					goto l713
				}
				depth--
				add(rulePausedOpt, position714)
			}
			return true
		l713:
			position, tokenIndex, depth = position713, tokenIndex713, depth713
			return false
		},
		/* 52 Expression <- <orExpr> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l720
				}
				depth--
				add(ruleExpression, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 53 orExpr <- <(<(andExpr sp (Or sp andExpr)*)> Action39)> */
		func() bool {
			position722, tokenIndex722, depth722 := position, tokenIndex, depth
			{
				position723 := position
				depth++
				{
					position724 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l722
					}
					if !_rules[rulesp]() {
						goto l722
					}
				l725:
					{
						position726, tokenIndex726, depth726 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l726
						}
						if !_rules[rulesp]() {
							goto l726
						}
						if !_rules[ruleandExpr]() {
							goto l726
						}
						goto l725
					l726:
						position, tokenIndex, depth = position726, tokenIndex726, depth726
					}
					depth--
					add(rulePegText, position724)
				}
				if !_rules[ruleAction39]() {
					goto l722
				}
				depth--
				add(ruleorExpr, position723)
			}
			return true
		l722:
			position, tokenIndex, depth = position722, tokenIndex722, depth722
			return false
		},
		/* 54 andExpr <- <(<(notExpr sp (And sp notExpr)*)> Action40)> */
		func() bool {
			position727, tokenIndex727, depth727 := position, tokenIndex, depth
			{
				position728 := position
				depth++
				{
					position729 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l727
					}
					if !_rules[rulesp]() {
						goto l727
					}
				l730:
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l731
						}
						if !_rules[rulesp]() {
							goto l731
						}
						if !_rules[rulenotExpr]() {
							goto l731
						}
						goto l730
					l731:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
					}
					depth--
					add(rulePegText, position729)
				}
				if !_rules[ruleAction40]() {
					goto l727
				}
				depth--
				add(ruleandExpr, position728)
			}
			return true
		l727:
			position, tokenIndex, depth = position727, tokenIndex727, depth727
			return false
		},
		/* 55 notExpr <- <(<((Not sp)? comparisonExpr)> Action41)> */
		func() bool {
			position732, tokenIndex732, depth732 := position, tokenIndex, depth
			{
				position733 := position
				depth++
				{
					position734 := position
					depth++
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l735
						}
						if !_rules[rulesp]() {
							goto l735
						}
						goto l736
					l735:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
					}
				l736:
					if !_rules[rulecomparisonExpr]() {
						goto l732
					}
					depth--
					add(rulePegText, position734)
				}
				if !_rules[ruleAction41]() {
					goto l732
				}
				depth--
				add(rulenotExpr, position733)
			}
			return true
		l732:
			position, tokenIndex, depth = position732, tokenIndex732, depth732
			return false
		},
		/* 56 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action42)> */
		func() bool {
			position737, tokenIndex737, depth737 := position, tokenIndex, depth
			{
				position738 := position
				depth++
				{
					position739 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l737
					}
					if !_rules[rulesp]() {
						goto l737
					}
					{
						position740, tokenIndex740, depth740 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l740
						}
						if !_rules[rulesp]() {
							goto l740
						}
						if !_rules[ruleotherOpExpr]() {
							goto l740
						}
						goto l741
					l740:
						position, tokenIndex, depth = position740, tokenIndex740, depth740
					}
				l741:
					depth--
					add(rulePegText, position739)
				}
				if !_rules[ruleAction42]() {
					goto l737
				}
				depth--
				add(rulecomparisonExpr, position738)
			}
			return true
		l737:
			position, tokenIndex, depth = position737, tokenIndex737, depth737
			return false
		},
		/* 57 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action43)> */
		func() bool {
			position742, tokenIndex742, depth742 := position, tokenIndex, depth
			{
				position743 := position
				depth++
				{
					position744 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l742
					}
					if !_rules[rulesp]() {
						goto l742
					}
				l745:
					{
						position746, tokenIndex746, depth746 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l746
						}
						if !_rules[rulesp]() {
							goto l746
						}
						if !_rules[ruleisExpr]() {
							goto l746
						}
						if !_rules[rulesp]() {
							goto l746
						}
						goto l745
					l746:
						position, tokenIndex, depth = position746, tokenIndex746, depth746
					}
					depth--
					add(rulePegText, position744)
				}
				if !_rules[ruleAction43]() {
					goto l742
				}
				depth--
				add(ruleotherOpExpr, position743)
			}
			return true
		l742:
			position, tokenIndex, depth = position742, tokenIndex742, depth742
			return false
		},
		/* 58 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action44)> */
		func() bool {
			position747, tokenIndex747, depth747 := position, tokenIndex, depth
			{
				position748 := position
				depth++
				{
					position749 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l747
					}
					if !_rules[rulesp]() {
						goto l747
					}
					{
						position750, tokenIndex750, depth750 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l750
						}
						if !_rules[rulesp]() {
							goto l750
						}
						if !_rules[ruleNullLiteral]() {
							goto l750
						}
						goto l751
					l750:
						position, tokenIndex, depth = position750, tokenIndex750, depth750
					}
				l751:
					depth--
					add(rulePegText, position749)
				}
				if !_rules[ruleAction44]() {
					goto l747
				}
				depth--
				add(ruleisExpr, position748)
			}
			return true
		l747:
			position, tokenIndex, depth = position747, tokenIndex747, depth747
			return false
		},
		/* 59 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action45)> */
		func() bool {
			position752, tokenIndex752, depth752 := position, tokenIndex, depth
			{
				position753 := position
				depth++
				{
					position754 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l752
					}
					if !_rules[rulesp]() {
						goto l752
					}
				l755:
					{
						position756, tokenIndex756, depth756 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l756
						}
						if !_rules[rulesp]() {
							goto l756
						}
						if !_rules[ruleproductExpr]() {
							goto l756
						}
						if !_rules[rulesp]() {
							goto l756
						}
						goto l755
					l756:
						position, tokenIndex, depth = position756, tokenIndex756, depth756
					}
					depth--
					add(rulePegText, position754)
				}
				if !_rules[ruleAction45]() {
					goto l752
				}
				depth--
				add(ruletermExpr, position753)
			}
			return true
		l752:
			position, tokenIndex, depth = position752, tokenIndex752, depth752
			return false
		},
		/* 60 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action46)> */
		func() bool {
			position757, tokenIndex757, depth757 := position, tokenIndex, depth
			{
				position758 := position
				depth++
				{
					position759 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l757
					}
					if !_rules[rulesp]() {
						goto l757
					}
				l760:
					{
						position761, tokenIndex761, depth761 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l761
						}
						if !_rules[rulesp]() {
							goto l761
						}
						if !_rules[ruleminusExpr]() {
							goto l761
						}
						if !_rules[rulesp]() {
							goto l761
						}
						goto l760
					l761:
						position, tokenIndex, depth = position761, tokenIndex761, depth761
					}
					depth--
					add(rulePegText, position759)
				}
				if !_rules[ruleAction46]() {
					goto l757
				}
				depth--
				add(ruleproductExpr, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 61 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action47)> */
		func() bool {
			position762, tokenIndex762, depth762 := position, tokenIndex, depth
			{
				position763 := position
				depth++
				{
					position764 := position
					depth++
					{
						position765, tokenIndex765, depth765 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l765
						}
						if !_rules[rulesp]() {
							goto l765
						}
						goto l766
					l765:
						position, tokenIndex, depth = position765, tokenIndex765, depth765
					}
				l766:
					if !_rules[rulecastExpr]() {
						goto l762
					}
					depth--
					add(rulePegText, position764)
				}
				if !_rules[ruleAction47]() {
					goto l762
				}
				depth--
				add(ruleminusExpr, position763)
			}
			return true
		l762:
			position, tokenIndex, depth = position762, tokenIndex762, depth762
			return false
		},
		/* 62 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action48)> */
		func() bool {
			position767, tokenIndex767, depth767 := position, tokenIndex, depth
			{
				position768 := position
				depth++
				{
					position769 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l767
					}
					{
						position770, tokenIndex770, depth770 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l770
						}
						if buffer[position] != rune(':') {
							goto l770
						}
						position++
						if buffer[position] != rune(':') {
							goto l770
						}
						position++
						if !_rules[rulesp]() {
							goto l770
						}
						if !_rules[ruleType]() {
							goto l770
						}
						goto l771
					l770:
						position, tokenIndex, depth = position770, tokenIndex770, depth770
					}
				l771:
					depth--
					add(rulePegText, position769)
				}
				if !_rules[ruleAction48]() {
					goto l767
				}
				depth--
				add(rulecastExpr, position768)
			}
			return true
		l767:
			position, tokenIndex, depth = position767, tokenIndex767, depth767
			return false
		},
		/* 63 baseExpr <- <(('(' sp Expression sp ')') / MapExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position772, tokenIndex772, depth772 := position, tokenIndex, depth
			{
				position773 := position
				depth++
				{
					position774, tokenIndex774, depth774 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l775
					}
					position++
					if !_rules[rulesp]() {
						goto l775
					}
					if !_rules[ruleExpression]() {
						goto l775
					}
					if !_rules[rulesp]() {
						goto l775
					}
					if buffer[position] != rune(')') {
						goto l775
					}
					position++
					goto l774
				l775:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleMapExpr]() {
						goto l776
					}
					goto l774
				l776:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleBooleanLiteral]() {
						goto l777
					}
					goto l774
				l777:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleNullLiteral]() {
						goto l778
					}
					goto l774
				l778:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleRowMeta]() {
						goto l779
					}
					goto l774
				l779:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleFuncTypeCast]() {
						goto l780
					}
					goto l774
				l780:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleFuncApp]() {
						goto l781
					}
					goto l774
				l781:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleRowValue]() {
						goto l782
					}
					goto l774
				l782:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleArrayExpr]() {
						goto l783
					}
					goto l774
				l783:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleLiteral]() {
						goto l772
					}
				}
			l774:
				depth--
				add(rulebaseExpr, position773)
			}
			return true
		l772:
			position, tokenIndex, depth = position772, tokenIndex772, depth772
			return false
		},
		/* 64 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action49)> */
		func() bool {
			position784, tokenIndex784, depth784 := position, tokenIndex, depth
			{
				position785 := position
				depth++
				{
					position786 := position
					depth++
					{
						position787, tokenIndex787, depth787 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l788
						}
						position++
						goto l787
					l788:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
						if buffer[position] != rune('C') {
							goto l784
						}
						position++
					}
				l787:
					{
						position789, tokenIndex789, depth789 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l790
						}
						position++
						goto l789
					l790:
						position, tokenIndex, depth = position789, tokenIndex789, depth789
						if buffer[position] != rune('A') {
							goto l784
						}
						position++
					}
				l789:
					{
						position791, tokenIndex791, depth791 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l792
						}
						position++
						goto l791
					l792:
						position, tokenIndex, depth = position791, tokenIndex791, depth791
						if buffer[position] != rune('S') {
							goto l784
						}
						position++
					}
				l791:
					{
						position793, tokenIndex793, depth793 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l794
						}
						position++
						goto l793
					l794:
						position, tokenIndex, depth = position793, tokenIndex793, depth793
						if buffer[position] != rune('T') {
							goto l784
						}
						position++
					}
				l793:
					if !_rules[rulesp]() {
						goto l784
					}
					if buffer[position] != rune('(') {
						goto l784
					}
					position++
					if !_rules[rulesp]() {
						goto l784
					}
					if !_rules[ruleExpression]() {
						goto l784
					}
					if !_rules[rulesp]() {
						goto l784
					}
					{
						position795, tokenIndex795, depth795 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l796
						}
						position++
						goto l795
					l796:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						if buffer[position] != rune('A') {
							goto l784
						}
						position++
					}
				l795:
					{
						position797, tokenIndex797, depth797 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l798
						}
						position++
						goto l797
					l798:
						position, tokenIndex, depth = position797, tokenIndex797, depth797
						if buffer[position] != rune('S') {
							goto l784
						}
						position++
					}
				l797:
					if !_rules[rulesp]() {
						goto l784
					}
					if !_rules[ruleType]() {
						goto l784
					}
					if !_rules[rulesp]() {
						goto l784
					}
					if buffer[position] != rune(')') {
						goto l784
					}
					position++
					depth--
					add(rulePegText, position786)
				}
				if !_rules[ruleAction49]() {
					goto l784
				}
				depth--
				add(ruleFuncTypeCast, position785)
			}
			return true
		l784:
			position, tokenIndex, depth = position784, tokenIndex784, depth784
			return false
		},
		/* 65 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action50)> */
		func() bool {
			position799, tokenIndex799, depth799 := position, tokenIndex, depth
			{
				position800 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l799
				}
				if !_rules[rulesp]() {
					goto l799
				}
				if buffer[position] != rune('(') {
					goto l799
				}
				position++
				if !_rules[rulesp]() {
					goto l799
				}
				if !_rules[ruleFuncParams]() {
					goto l799
				}
				if !_rules[rulesp]() {
					goto l799
				}
				if buffer[position] != rune(')') {
					goto l799
				}
				position++
				if !_rules[ruleAction50]() {
					goto l799
				}
				depth--
				add(ruleFuncApp, position800)
			}
			return true
		l799:
			position, tokenIndex, depth = position799, tokenIndex799, depth799
			return false
		},
		/* 66 FuncParams <- <(<(Star / (Expression sp (',' sp Expression)*)?)> Action51)> */
		func() bool {
			position801, tokenIndex801, depth801 := position, tokenIndex, depth
			{
				position802 := position
				depth++
				{
					position803 := position
					depth++
					{
						position804, tokenIndex804, depth804 := position, tokenIndex, depth
						if !_rules[ruleStar]() {
							goto l805
						}
						goto l804
					l805:
						position, tokenIndex, depth = position804, tokenIndex804, depth804
						{
							position806, tokenIndex806, depth806 := position, tokenIndex, depth
							if !_rules[ruleExpression]() {
								goto l806
							}
							if !_rules[rulesp]() {
								goto l806
							}
						l808:
							{
								position809, tokenIndex809, depth809 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l809
								}
								position++
								if !_rules[rulesp]() {
									goto l809
								}
								if !_rules[ruleExpression]() {
									goto l809
								}
								goto l808
							l809:
								position, tokenIndex, depth = position809, tokenIndex809, depth809
							}
							goto l807
						l806:
							position, tokenIndex, depth = position806, tokenIndex806, depth806
						}
					l807:
					}
				l804:
					depth--
					add(rulePegText, position803)
				}
				if !_rules[ruleAction51]() {
					goto l801
				}
				depth--
				add(ruleFuncParams, position802)
			}
			return true
		l801:
			position, tokenIndex, depth = position801, tokenIndex801, depth801
			return false
		},
		/* 67 ArrayExpr <- <(<('[' sp (Expression (',' sp Expression)*)? sp ','? sp ']')> Action52)> */
		func() bool {
			position810, tokenIndex810, depth810 := position, tokenIndex, depth
			{
				position811 := position
				depth++
				{
					position812 := position
					depth++
					if buffer[position] != rune('[') {
						goto l810
					}
					position++
					if !_rules[rulesp]() {
						goto l810
					}
					{
						position813, tokenIndex813, depth813 := position, tokenIndex, depth
						if !_rules[ruleExpression]() {
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
					if !_rules[rulesp]() {
						goto l810
					}
					{
						position817, tokenIndex817, depth817 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l817
						}
						position++
						goto l818
					l817:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
					}
				l818:
					if !_rules[rulesp]() {
						goto l810
					}
					if buffer[position] != rune(']') {
						goto l810
					}
					position++
					depth--
					add(rulePegText, position812)
				}
				if !_rules[ruleAction52]() {
					goto l810
				}
				depth--
				add(ruleArrayExpr, position811)
			}
			return true
		l810:
			position, tokenIndex, depth = position810, tokenIndex810, depth810
			return false
		},
		/* 68 MapExpr <- <(<('{' sp (KeyValuePair (',' sp KeyValuePair)*)? sp '}')> Action53)> */
		func() bool {
			position819, tokenIndex819, depth819 := position, tokenIndex, depth
			{
				position820 := position
				depth++
				{
					position821 := position
					depth++
					if buffer[position] != rune('{') {
						goto l819
					}
					position++
					if !_rules[rulesp]() {
						goto l819
					}
					{
						position822, tokenIndex822, depth822 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l822
						}
					l824:
						{
							position825, tokenIndex825, depth825 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l825
							}
							position++
							if !_rules[rulesp]() {
								goto l825
							}
							if !_rules[ruleKeyValuePair]() {
								goto l825
							}
							goto l824
						l825:
							position, tokenIndex, depth = position825, tokenIndex825, depth825
						}
						goto l823
					l822:
						position, tokenIndex, depth = position822, tokenIndex822, depth822
					}
				l823:
					if !_rules[rulesp]() {
						goto l819
					}
					if buffer[position] != rune('}') {
						goto l819
					}
					position++
					depth--
					add(rulePegText, position821)
				}
				if !_rules[ruleAction53]() {
					goto l819
				}
				depth--
				add(ruleMapExpr, position820)
			}
			return true
		l819:
			position, tokenIndex, depth = position819, tokenIndex819, depth819
			return false
		},
		/* 69 KeyValuePair <- <(<(StringLiteral sp ':' sp Expression)> Action54)> */
		func() bool {
			position826, tokenIndex826, depth826 := position, tokenIndex, depth
			{
				position827 := position
				depth++
				{
					position828 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l826
					}
					if !_rules[rulesp]() {
						goto l826
					}
					if buffer[position] != rune(':') {
						goto l826
					}
					position++
					if !_rules[rulesp]() {
						goto l826
					}
					if !_rules[ruleExpression]() {
						goto l826
					}
					depth--
					add(rulePegText, position828)
				}
				if !_rules[ruleAction54]() {
					goto l826
				}
				depth--
				add(ruleKeyValuePair, position827)
			}
			return true
		l826:
			position, tokenIndex, depth = position826, tokenIndex826, depth826
			return false
		},
		/* 70 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position829, tokenIndex829, depth829 := position, tokenIndex, depth
			{
				position830 := position
				depth++
				{
					position831, tokenIndex831, depth831 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l832
					}
					goto l831
				l832:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					if !_rules[ruleNumericLiteral]() {
						goto l833
					}
					goto l831
				l833:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					if !_rules[ruleStringLiteral]() {
						goto l829
					}
				}
			l831:
				depth--
				add(ruleLiteral, position830)
			}
			return true
		l829:
			position, tokenIndex, depth = position829, tokenIndex829, depth829
			return false
		},
		/* 71 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position834, tokenIndex834, depth834 := position, tokenIndex, depth
			{
				position835 := position
				depth++
				{
					position836, tokenIndex836, depth836 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l837
					}
					goto l836
				l837:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleNotEqual]() {
						goto l838
					}
					goto l836
				l838:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleLessOrEqual]() {
						goto l839
					}
					goto l836
				l839:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleLess]() {
						goto l840
					}
					goto l836
				l840:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleGreaterOrEqual]() {
						goto l841
					}
					goto l836
				l841:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleGreater]() {
						goto l842
					}
					goto l836
				l842:
					position, tokenIndex, depth = position836, tokenIndex836, depth836
					if !_rules[ruleNotEqual]() {
						goto l834
					}
				}
			l836:
				depth--
				add(ruleComparisonOp, position835)
			}
			return true
		l834:
			position, tokenIndex, depth = position834, tokenIndex834, depth834
			return false
		},
		/* 72 OtherOp <- <Concat> */
		func() bool {
			position843, tokenIndex843, depth843 := position, tokenIndex, depth
			{
				position844 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l843
				}
				depth--
				add(ruleOtherOp, position844)
			}
			return true
		l843:
			position, tokenIndex, depth = position843, tokenIndex843, depth843
			return false
		},
		/* 73 IsOp <- <(IsNot / Is)> */
		func() bool {
			position845, tokenIndex845, depth845 := position, tokenIndex, depth
			{
				position846 := position
				depth++
				{
					position847, tokenIndex847, depth847 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l848
					}
					goto l847
				l848:
					position, tokenIndex, depth = position847, tokenIndex847, depth847
					if !_rules[ruleIs]() {
						goto l845
					}
				}
			l847:
				depth--
				add(ruleIsOp, position846)
			}
			return true
		l845:
			position, tokenIndex, depth = position845, tokenIndex845, depth845
			return false
		},
		/* 74 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position849, tokenIndex849, depth849 := position, tokenIndex, depth
			{
				position850 := position
				depth++
				{
					position851, tokenIndex851, depth851 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l852
					}
					goto l851
				l852:
					position, tokenIndex, depth = position851, tokenIndex851, depth851
					if !_rules[ruleMinus]() {
						goto l849
					}
				}
			l851:
				depth--
				add(rulePlusMinusOp, position850)
			}
			return true
		l849:
			position, tokenIndex, depth = position849, tokenIndex849, depth849
			return false
		},
		/* 75 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position853, tokenIndex853, depth853 := position, tokenIndex, depth
			{
				position854 := position
				depth++
				{
					position855, tokenIndex855, depth855 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l856
					}
					goto l855
				l856:
					position, tokenIndex, depth = position855, tokenIndex855, depth855
					if !_rules[ruleDivide]() {
						goto l857
					}
					goto l855
				l857:
					position, tokenIndex, depth = position855, tokenIndex855, depth855
					if !_rules[ruleModulo]() {
						goto l853
					}
				}
			l855:
				depth--
				add(ruleMultDivOp, position854)
			}
			return true
		l853:
			position, tokenIndex, depth = position853, tokenIndex853, depth853
			return false
		},
		/* 76 Stream <- <(<ident> Action55)> */
		func() bool {
			position858, tokenIndex858, depth858 := position, tokenIndex, depth
			{
				position859 := position
				depth++
				{
					position860 := position
					depth++
					if !_rules[ruleident]() {
						goto l858
					}
					depth--
					add(rulePegText, position860)
				}
				if !_rules[ruleAction55]() {
					goto l858
				}
				depth--
				add(ruleStream, position859)
			}
			return true
		l858:
			position, tokenIndex, depth = position858, tokenIndex858, depth858
			return false
		},
		/* 77 RowMeta <- <RowTimestamp> */
		func() bool {
			position861, tokenIndex861, depth861 := position, tokenIndex, depth
			{
				position862 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l861
				}
				depth--
				add(ruleRowMeta, position862)
			}
			return true
		l861:
			position, tokenIndex, depth = position861, tokenIndex861, depth861
			return false
		},
		/* 78 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action56)> */
		func() bool {
			position863, tokenIndex863, depth863 := position, tokenIndex, depth
			{
				position864 := position
				depth++
				{
					position865 := position
					depth++
					{
						position866, tokenIndex866, depth866 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l866
						}
						if buffer[position] != rune(':') {
							goto l866
						}
						position++
						goto l867
					l866:
						position, tokenIndex, depth = position866, tokenIndex866, depth866
					}
				l867:
					if buffer[position] != rune('t') {
						goto l863
					}
					position++
					if buffer[position] != rune('s') {
						goto l863
					}
					position++
					if buffer[position] != rune('(') {
						goto l863
					}
					position++
					if buffer[position] != rune(')') {
						goto l863
					}
					position++
					depth--
					add(rulePegText, position865)
				}
				if !_rules[ruleAction56]() {
					goto l863
				}
				depth--
				add(ruleRowTimestamp, position864)
			}
			return true
		l863:
			position, tokenIndex, depth = position863, tokenIndex863, depth863
			return false
		},
		/* 79 RowValue <- <(<((ident ':' !':')? jsonPath)> Action57)> */
		func() bool {
			position868, tokenIndex868, depth868 := position, tokenIndex, depth
			{
				position869 := position
				depth++
				{
					position870 := position
					depth++
					{
						position871, tokenIndex871, depth871 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l871
						}
						if buffer[position] != rune(':') {
							goto l871
						}
						position++
						{
							position873, tokenIndex873, depth873 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l873
							}
							position++
							goto l871
						l873:
							position, tokenIndex, depth = position873, tokenIndex873, depth873
						}
						goto l872
					l871:
						position, tokenIndex, depth = position871, tokenIndex871, depth871
					}
				l872:
					if !_rules[rulejsonPath]() {
						goto l868
					}
					depth--
					add(rulePegText, position870)
				}
				if !_rules[ruleAction57]() {
					goto l868
				}
				depth--
				add(ruleRowValue, position869)
			}
			return true
		l868:
			position, tokenIndex, depth = position868, tokenIndex868, depth868
			return false
		},
		/* 80 NumericLiteral <- <(<('-'? [0-9]+)> Action58)> */
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
						if buffer[position] != rune('-') {
							goto l877
						}
						position++
						goto l878
					l877:
						position, tokenIndex, depth = position877, tokenIndex877, depth877
					}
				l878:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l874
					}
					position++
				l879:
					{
						position880, tokenIndex880, depth880 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l880
						}
						position++
						goto l879
					l880:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
					}
					depth--
					add(rulePegText, position876)
				}
				if !_rules[ruleAction58]() {
					goto l874
				}
				depth--
				add(ruleNumericLiteral, position875)
			}
			return true
		l874:
			position, tokenIndex, depth = position874, tokenIndex874, depth874
			return false
		},
		/* 81 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action59)> */
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
					if buffer[position] != rune('.') {
						goto l881
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l881
					}
					position++
				l888:
					{
						position889, tokenIndex889, depth889 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l889
						}
						position++
						goto l888
					l889:
						position, tokenIndex, depth = position889, tokenIndex889, depth889
					}
					depth--
					add(rulePegText, position883)
				}
				if !_rules[ruleAction59]() {
					goto l881
				}
				depth--
				add(ruleFloatLiteral, position882)
			}
			return true
		l881:
			position, tokenIndex, depth = position881, tokenIndex881, depth881
			return false
		},
		/* 82 Function <- <(<ident> Action60)> */
		func() bool {
			position890, tokenIndex890, depth890 := position, tokenIndex, depth
			{
				position891 := position
				depth++
				{
					position892 := position
					depth++
					if !_rules[ruleident]() {
						goto l890
					}
					depth--
					add(rulePegText, position892)
				}
				if !_rules[ruleAction60]() {
					goto l890
				}
				depth--
				add(ruleFunction, position891)
			}
			return true
		l890:
			position, tokenIndex, depth = position890, tokenIndex890, depth890
			return false
		},
		/* 83 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action61)> */
		func() bool {
			position893, tokenIndex893, depth893 := position, tokenIndex, depth
			{
				position894 := position
				depth++
				{
					position895 := position
					depth++
					{
						position896, tokenIndex896, depth896 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l897
						}
						position++
						goto l896
					l897:
						position, tokenIndex, depth = position896, tokenIndex896, depth896
						if buffer[position] != rune('N') {
							goto l893
						}
						position++
					}
				l896:
					{
						position898, tokenIndex898, depth898 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l899
						}
						position++
						goto l898
					l899:
						position, tokenIndex, depth = position898, tokenIndex898, depth898
						if buffer[position] != rune('U') {
							goto l893
						}
						position++
					}
				l898:
					{
						position900, tokenIndex900, depth900 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l901
						}
						position++
						goto l900
					l901:
						position, tokenIndex, depth = position900, tokenIndex900, depth900
						if buffer[position] != rune('L') {
							goto l893
						}
						position++
					}
				l900:
					{
						position902, tokenIndex902, depth902 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l903
						}
						position++
						goto l902
					l903:
						position, tokenIndex, depth = position902, tokenIndex902, depth902
						if buffer[position] != rune('L') {
							goto l893
						}
						position++
					}
				l902:
					depth--
					add(rulePegText, position895)
				}
				if !_rules[ruleAction61]() {
					goto l893
				}
				depth--
				add(ruleNullLiteral, position894)
			}
			return true
		l893:
			position, tokenIndex, depth = position893, tokenIndex893, depth893
			return false
		},
		/* 84 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position904, tokenIndex904, depth904 := position, tokenIndex, depth
			{
				position905 := position
				depth++
				{
					position906, tokenIndex906, depth906 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l907
					}
					goto l906
				l907:
					position, tokenIndex, depth = position906, tokenIndex906, depth906
					if !_rules[ruleFALSE]() {
						goto l904
					}
				}
			l906:
				depth--
				add(ruleBooleanLiteral, position905)
			}
			return true
		l904:
			position, tokenIndex, depth = position904, tokenIndex904, depth904
			return false
		},
		/* 85 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action62)> */
		func() bool {
			position908, tokenIndex908, depth908 := position, tokenIndex, depth
			{
				position909 := position
				depth++
				{
					position910 := position
					depth++
					{
						position911, tokenIndex911, depth911 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l912
						}
						position++
						goto l911
					l912:
						position, tokenIndex, depth = position911, tokenIndex911, depth911
						if buffer[position] != rune('T') {
							goto l908
						}
						position++
					}
				l911:
					{
						position913, tokenIndex913, depth913 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l914
						}
						position++
						goto l913
					l914:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
						if buffer[position] != rune('R') {
							goto l908
						}
						position++
					}
				l913:
					{
						position915, tokenIndex915, depth915 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l916
						}
						position++
						goto l915
					l916:
						position, tokenIndex, depth = position915, tokenIndex915, depth915
						if buffer[position] != rune('U') {
							goto l908
						}
						position++
					}
				l915:
					{
						position917, tokenIndex917, depth917 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l918
						}
						position++
						goto l917
					l918:
						position, tokenIndex, depth = position917, tokenIndex917, depth917
						if buffer[position] != rune('E') {
							goto l908
						}
						position++
					}
				l917:
					depth--
					add(rulePegText, position910)
				}
				if !_rules[ruleAction62]() {
					goto l908
				}
				depth--
				add(ruleTRUE, position909)
			}
			return true
		l908:
			position, tokenIndex, depth = position908, tokenIndex908, depth908
			return false
		},
		/* 86 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action63)> */
		func() bool {
			position919, tokenIndex919, depth919 := position, tokenIndex, depth
			{
				position920 := position
				depth++
				{
					position921 := position
					depth++
					{
						position922, tokenIndex922, depth922 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l923
						}
						position++
						goto l922
					l923:
						position, tokenIndex, depth = position922, tokenIndex922, depth922
						if buffer[position] != rune('F') {
							goto l919
						}
						position++
					}
				l922:
					{
						position924, tokenIndex924, depth924 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l925
						}
						position++
						goto l924
					l925:
						position, tokenIndex, depth = position924, tokenIndex924, depth924
						if buffer[position] != rune('A') {
							goto l919
						}
						position++
					}
				l924:
					{
						position926, tokenIndex926, depth926 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('L') {
							goto l919
						}
						position++
					}
				l926:
					{
						position928, tokenIndex928, depth928 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l929
						}
						position++
						goto l928
					l929:
						position, tokenIndex, depth = position928, tokenIndex928, depth928
						if buffer[position] != rune('S') {
							goto l919
						}
						position++
					}
				l928:
					{
						position930, tokenIndex930, depth930 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l931
						}
						position++
						goto l930
					l931:
						position, tokenIndex, depth = position930, tokenIndex930, depth930
						if buffer[position] != rune('E') {
							goto l919
						}
						position++
					}
				l930:
					depth--
					add(rulePegText, position921)
				}
				if !_rules[ruleAction63]() {
					goto l919
				}
				depth--
				add(ruleFALSE, position920)
			}
			return true
		l919:
			position, tokenIndex, depth = position919, tokenIndex919, depth919
			return false
		},
		/* 87 Star <- <(<'*'> Action64)> */
		func() bool {
			position932, tokenIndex932, depth932 := position, tokenIndex, depth
			{
				position933 := position
				depth++
				{
					position934 := position
					depth++
					if buffer[position] != rune('*') {
						goto l932
					}
					position++
					depth--
					add(rulePegText, position934)
				}
				if !_rules[ruleAction64]() {
					goto l932
				}
				depth--
				add(ruleStar, position933)
			}
			return true
		l932:
			position, tokenIndex, depth = position932, tokenIndex932, depth932
			return false
		},
		/* 88 Wildcard <- <(<((ident ':' !':')? '*')> Action65)> */
		func() bool {
			position935, tokenIndex935, depth935 := position, tokenIndex, depth
			{
				position936 := position
				depth++
				{
					position937 := position
					depth++
					{
						position938, tokenIndex938, depth938 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l938
						}
						if buffer[position] != rune(':') {
							goto l938
						}
						position++
						{
							position940, tokenIndex940, depth940 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l940
							}
							position++
							goto l938
						l940:
							position, tokenIndex, depth = position940, tokenIndex940, depth940
						}
						goto l939
					l938:
						position, tokenIndex, depth = position938, tokenIndex938, depth938
					}
				l939:
					if buffer[position] != rune('*') {
						goto l935
					}
					position++
					depth--
					add(rulePegText, position937)
				}
				if !_rules[ruleAction65]() {
					goto l935
				}
				depth--
				add(ruleWildcard, position936)
			}
			return true
		l935:
			position, tokenIndex, depth = position935, tokenIndex935, depth935
			return false
		},
		/* 89 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action66)> */
		func() bool {
			position941, tokenIndex941, depth941 := position, tokenIndex, depth
			{
				position942 := position
				depth++
				{
					position943 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l941
					}
					position++
				l944:
					{
						position945, tokenIndex945, depth945 := position, tokenIndex, depth
						{
							position946, tokenIndex946, depth946 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l947
							}
							position++
							if buffer[position] != rune('\'') {
								goto l947
							}
							position++
							goto l946
						l947:
							position, tokenIndex, depth = position946, tokenIndex946, depth946
							{
								position948, tokenIndex948, depth948 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l948
								}
								position++
								goto l945
							l948:
								position, tokenIndex, depth = position948, tokenIndex948, depth948
							}
							if !matchDot() {
								goto l945
							}
						}
					l946:
						goto l944
					l945:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
					}
					if buffer[position] != rune('\'') {
						goto l941
					}
					position++
					depth--
					add(rulePegText, position943)
				}
				if !_rules[ruleAction66]() {
					goto l941
				}
				depth--
				add(ruleStringLiteral, position942)
			}
			return true
		l941:
			position, tokenIndex, depth = position941, tokenIndex941, depth941
			return false
		},
		/* 90 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action67)> */
		func() bool {
			position949, tokenIndex949, depth949 := position, tokenIndex, depth
			{
				position950 := position
				depth++
				{
					position951 := position
					depth++
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
							goto l949
						}
						position++
					}
				l952:
					{
						position954, tokenIndex954, depth954 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l955
						}
						position++
						goto l954
					l955:
						position, tokenIndex, depth = position954, tokenIndex954, depth954
						if buffer[position] != rune('S') {
							goto l949
						}
						position++
					}
				l954:
					{
						position956, tokenIndex956, depth956 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l957
						}
						position++
						goto l956
					l957:
						position, tokenIndex, depth = position956, tokenIndex956, depth956
						if buffer[position] != rune('T') {
							goto l949
						}
						position++
					}
				l956:
					{
						position958, tokenIndex958, depth958 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l959
						}
						position++
						goto l958
					l959:
						position, tokenIndex, depth = position958, tokenIndex958, depth958
						if buffer[position] != rune('R') {
							goto l949
						}
						position++
					}
				l958:
					{
						position960, tokenIndex960, depth960 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l961
						}
						position++
						goto l960
					l961:
						position, tokenIndex, depth = position960, tokenIndex960, depth960
						if buffer[position] != rune('E') {
							goto l949
						}
						position++
					}
				l960:
					{
						position962, tokenIndex962, depth962 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l963
						}
						position++
						goto l962
					l963:
						position, tokenIndex, depth = position962, tokenIndex962, depth962
						if buffer[position] != rune('A') {
							goto l949
						}
						position++
					}
				l962:
					{
						position964, tokenIndex964, depth964 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l965
						}
						position++
						goto l964
					l965:
						position, tokenIndex, depth = position964, tokenIndex964, depth964
						if buffer[position] != rune('M') {
							goto l949
						}
						position++
					}
				l964:
					depth--
					add(rulePegText, position951)
				}
				if !_rules[ruleAction67]() {
					goto l949
				}
				depth--
				add(ruleISTREAM, position950)
			}
			return true
		l949:
			position, tokenIndex, depth = position949, tokenIndex949, depth949
			return false
		},
		/* 91 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action68)> */
		func() bool {
			position966, tokenIndex966, depth966 := position, tokenIndex, depth
			{
				position967 := position
				depth++
				{
					position968 := position
					depth++
					{
						position969, tokenIndex969, depth969 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l970
						}
						position++
						goto l969
					l970:
						position, tokenIndex, depth = position969, tokenIndex969, depth969
						if buffer[position] != rune('D') {
							goto l966
						}
						position++
					}
				l969:
					{
						position971, tokenIndex971, depth971 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l972
						}
						position++
						goto l971
					l972:
						position, tokenIndex, depth = position971, tokenIndex971, depth971
						if buffer[position] != rune('S') {
							goto l966
						}
						position++
					}
				l971:
					{
						position973, tokenIndex973, depth973 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l974
						}
						position++
						goto l973
					l974:
						position, tokenIndex, depth = position973, tokenIndex973, depth973
						if buffer[position] != rune('T') {
							goto l966
						}
						position++
					}
				l973:
					{
						position975, tokenIndex975, depth975 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l976
						}
						position++
						goto l975
					l976:
						position, tokenIndex, depth = position975, tokenIndex975, depth975
						if buffer[position] != rune('R') {
							goto l966
						}
						position++
					}
				l975:
					{
						position977, tokenIndex977, depth977 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l978
						}
						position++
						goto l977
					l978:
						position, tokenIndex, depth = position977, tokenIndex977, depth977
						if buffer[position] != rune('E') {
							goto l966
						}
						position++
					}
				l977:
					{
						position979, tokenIndex979, depth979 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l980
						}
						position++
						goto l979
					l980:
						position, tokenIndex, depth = position979, tokenIndex979, depth979
						if buffer[position] != rune('A') {
							goto l966
						}
						position++
					}
				l979:
					{
						position981, tokenIndex981, depth981 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l982
						}
						position++
						goto l981
					l982:
						position, tokenIndex, depth = position981, tokenIndex981, depth981
						if buffer[position] != rune('M') {
							goto l966
						}
						position++
					}
				l981:
					depth--
					add(rulePegText, position968)
				}
				if !_rules[ruleAction68]() {
					goto l966
				}
				depth--
				add(ruleDSTREAM, position967)
			}
			return true
		l966:
			position, tokenIndex, depth = position966, tokenIndex966, depth966
			return false
		},
		/* 92 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action69)> */
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
						if buffer[position] != rune('r') {
							goto l987
						}
						position++
						goto l986
					l987:
						position, tokenIndex, depth = position986, tokenIndex986, depth986
						if buffer[position] != rune('R') {
							goto l983
						}
						position++
					}
				l986:
					{
						position988, tokenIndex988, depth988 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l989
						}
						position++
						goto l988
					l989:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('S') {
							goto l983
						}
						position++
					}
				l988:
					{
						position990, tokenIndex990, depth990 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l991
						}
						position++
						goto l990
					l991:
						position, tokenIndex, depth = position990, tokenIndex990, depth990
						if buffer[position] != rune('T') {
							goto l983
						}
						position++
					}
				l990:
					{
						position992, tokenIndex992, depth992 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l993
						}
						position++
						goto l992
					l993:
						position, tokenIndex, depth = position992, tokenIndex992, depth992
						if buffer[position] != rune('R') {
							goto l983
						}
						position++
					}
				l992:
					{
						position994, tokenIndex994, depth994 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l995
						}
						position++
						goto l994
					l995:
						position, tokenIndex, depth = position994, tokenIndex994, depth994
						if buffer[position] != rune('E') {
							goto l983
						}
						position++
					}
				l994:
					{
						position996, tokenIndex996, depth996 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l997
						}
						position++
						goto l996
					l997:
						position, tokenIndex, depth = position996, tokenIndex996, depth996
						if buffer[position] != rune('A') {
							goto l983
						}
						position++
					}
				l996:
					{
						position998, tokenIndex998, depth998 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l999
						}
						position++
						goto l998
					l999:
						position, tokenIndex, depth = position998, tokenIndex998, depth998
						if buffer[position] != rune('M') {
							goto l983
						}
						position++
					}
				l998:
					depth--
					add(rulePegText, position985)
				}
				if !_rules[ruleAction69]() {
					goto l983
				}
				depth--
				add(ruleRSTREAM, position984)
			}
			return true
		l983:
			position, tokenIndex, depth = position983, tokenIndex983, depth983
			return false
		},
		/* 93 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action70)> */
		func() bool {
			position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
			{
				position1001 := position
				depth++
				{
					position1002 := position
					depth++
					{
						position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1004
						}
						position++
						goto l1003
					l1004:
						position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
						if buffer[position] != rune('T') {
							goto l1000
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
							goto l1000
						}
						position++
					}
				l1005:
					{
						position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1008
						}
						position++
						goto l1007
					l1008:
						position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
						if buffer[position] != rune('P') {
							goto l1000
						}
						position++
					}
				l1007:
					{
						position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1010
						}
						position++
						goto l1009
					l1010:
						position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
						if buffer[position] != rune('L') {
							goto l1000
						}
						position++
					}
				l1009:
					{
						position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1012
						}
						position++
						goto l1011
					l1012:
						position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
						if buffer[position] != rune('E') {
							goto l1000
						}
						position++
					}
				l1011:
					{
						position1013, tokenIndex1013, depth1013 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1014
						}
						position++
						goto l1013
					l1014:
						position, tokenIndex, depth = position1013, tokenIndex1013, depth1013
						if buffer[position] != rune('S') {
							goto l1000
						}
						position++
					}
				l1013:
					depth--
					add(rulePegText, position1002)
				}
				if !_rules[ruleAction70]() {
					goto l1000
				}
				depth--
				add(ruleTUPLES, position1001)
			}
			return true
		l1000:
			position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
			return false
		},
		/* 94 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action71)> */
		func() bool {
			position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
			{
				position1016 := position
				depth++
				{
					position1017 := position
					depth++
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
							goto l1015
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
							goto l1015
						}
						position++
					}
				l1020:
					{
						position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1023
						}
						position++
						goto l1022
					l1023:
						position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
						if buffer[position] != rune('C') {
							goto l1015
						}
						position++
					}
				l1022:
					{
						position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1025
						}
						position++
						goto l1024
					l1025:
						position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
						if buffer[position] != rune('O') {
							goto l1015
						}
						position++
					}
				l1024:
					{
						position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1027
						}
						position++
						goto l1026
					l1027:
						position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
						if buffer[position] != rune('N') {
							goto l1015
						}
						position++
					}
				l1026:
					{
						position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1029
						}
						position++
						goto l1028
					l1029:
						position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
						if buffer[position] != rune('D') {
							goto l1015
						}
						position++
					}
				l1028:
					{
						position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1031
						}
						position++
						goto l1030
					l1031:
						position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
						if buffer[position] != rune('S') {
							goto l1015
						}
						position++
					}
				l1030:
					depth--
					add(rulePegText, position1017)
				}
				if !_rules[ruleAction71]() {
					goto l1015
				}
				depth--
				add(ruleSECONDS, position1016)
			}
			return true
		l1015:
			position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
			return false
		},
		/* 95 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action72)> */
		func() bool {
			position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
			{
				position1033 := position
				depth++
				{
					position1034 := position
					depth++
					{
						position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1036
						}
						position++
						goto l1035
					l1036:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						if buffer[position] != rune('M') {
							goto l1032
						}
						position++
					}
				l1035:
					{
						position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1038
						}
						position++
						goto l1037
					l1038:
						position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
						if buffer[position] != rune('I') {
							goto l1032
						}
						position++
					}
				l1037:
					{
						position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1040
						}
						position++
						goto l1039
					l1040:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if buffer[position] != rune('L') {
							goto l1032
						}
						position++
					}
				l1039:
					{
						position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1042
						}
						position++
						goto l1041
					l1042:
						position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
						if buffer[position] != rune('L') {
							goto l1032
						}
						position++
					}
				l1041:
					{
						position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1044
						}
						position++
						goto l1043
					l1044:
						position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
						if buffer[position] != rune('I') {
							goto l1032
						}
						position++
					}
				l1043:
					{
						position1045, tokenIndex1045, depth1045 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1046
						}
						position++
						goto l1045
					l1046:
						position, tokenIndex, depth = position1045, tokenIndex1045, depth1045
						if buffer[position] != rune('S') {
							goto l1032
						}
						position++
					}
				l1045:
					{
						position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1048
						}
						position++
						goto l1047
					l1048:
						position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
						if buffer[position] != rune('E') {
							goto l1032
						}
						position++
					}
				l1047:
					{
						position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1050
						}
						position++
						goto l1049
					l1050:
						position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
						if buffer[position] != rune('C') {
							goto l1032
						}
						position++
					}
				l1049:
					{
						position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1052
						}
						position++
						goto l1051
					l1052:
						position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
						if buffer[position] != rune('O') {
							goto l1032
						}
						position++
					}
				l1051:
					{
						position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1054
						}
						position++
						goto l1053
					l1054:
						position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
						if buffer[position] != rune('N') {
							goto l1032
						}
						position++
					}
				l1053:
					{
						position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1056
						}
						position++
						goto l1055
					l1056:
						position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
						if buffer[position] != rune('D') {
							goto l1032
						}
						position++
					}
				l1055:
					{
						position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1058
						}
						position++
						goto l1057
					l1058:
						position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
						if buffer[position] != rune('S') {
							goto l1032
						}
						position++
					}
				l1057:
					depth--
					add(rulePegText, position1034)
				}
				if !_rules[ruleAction72]() {
					goto l1032
				}
				depth--
				add(ruleMILLISECONDS, position1033)
			}
			return true
		l1032:
			position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
			return false
		},
		/* 96 StreamIdentifier <- <(<ident> Action73)> */
		func() bool {
			position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
			{
				position1060 := position
				depth++
				{
					position1061 := position
					depth++
					if !_rules[ruleident]() {
						goto l1059
					}
					depth--
					add(rulePegText, position1061)
				}
				if !_rules[ruleAction73]() {
					goto l1059
				}
				depth--
				add(ruleStreamIdentifier, position1060)
			}
			return true
		l1059:
			position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
			return false
		},
		/* 97 SourceSinkType <- <(<ident> Action74)> */
		func() bool {
			position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
			{
				position1063 := position
				depth++
				{
					position1064 := position
					depth++
					if !_rules[ruleident]() {
						goto l1062
					}
					depth--
					add(rulePegText, position1064)
				}
				if !_rules[ruleAction74]() {
					goto l1062
				}
				depth--
				add(ruleSourceSinkType, position1063)
			}
			return true
		l1062:
			position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
			return false
		},
		/* 98 SourceSinkParamKey <- <(<ident> Action75)> */
		func() bool {
			position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
			{
				position1066 := position
				depth++
				{
					position1067 := position
					depth++
					if !_rules[ruleident]() {
						goto l1065
					}
					depth--
					add(rulePegText, position1067)
				}
				if !_rules[ruleAction75]() {
					goto l1065
				}
				depth--
				add(ruleSourceSinkParamKey, position1066)
			}
			return true
		l1065:
			position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
			return false
		},
		/* 99 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action76)> */
		func() bool {
			position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
			{
				position1069 := position
				depth++
				{
					position1070 := position
					depth++
					{
						position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1072
						}
						position++
						goto l1071
					l1072:
						position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
						if buffer[position] != rune('P') {
							goto l1068
						}
						position++
					}
				l1071:
					{
						position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1074
						}
						position++
						goto l1073
					l1074:
						position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
						if buffer[position] != rune('A') {
							goto l1068
						}
						position++
					}
				l1073:
					{
						position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1076
						}
						position++
						goto l1075
					l1076:
						position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
						if buffer[position] != rune('U') {
							goto l1068
						}
						position++
					}
				l1075:
					{
						position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1078
						}
						position++
						goto l1077
					l1078:
						position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
						if buffer[position] != rune('S') {
							goto l1068
						}
						position++
					}
				l1077:
					{
						position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1080
						}
						position++
						goto l1079
					l1080:
						position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
						if buffer[position] != rune('E') {
							goto l1068
						}
						position++
					}
				l1079:
					{
						position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1082
						}
						position++
						goto l1081
					l1082:
						position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
						if buffer[position] != rune('D') {
							goto l1068
						}
						position++
					}
				l1081:
					depth--
					add(rulePegText, position1070)
				}
				if !_rules[ruleAction76]() {
					goto l1068
				}
				depth--
				add(rulePaused, position1069)
			}
			return true
		l1068:
			position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
			return false
		},
		/* 100 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action77)> */
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
						if buffer[position] != rune('u') {
							goto l1087
						}
						position++
						goto l1086
					l1087:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
						if buffer[position] != rune('U') {
							goto l1083
						}
						position++
					}
				l1086:
					{
						position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1089
						}
						position++
						goto l1088
					l1089:
						position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
						if buffer[position] != rune('N') {
							goto l1083
						}
						position++
					}
				l1088:
					{
						position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1091
						}
						position++
						goto l1090
					l1091:
						position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
						if buffer[position] != rune('P') {
							goto l1083
						}
						position++
					}
				l1090:
					{
						position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1093
						}
						position++
						goto l1092
					l1093:
						position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
						if buffer[position] != rune('A') {
							goto l1083
						}
						position++
					}
				l1092:
					{
						position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1095
						}
						position++
						goto l1094
					l1095:
						position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
						if buffer[position] != rune('U') {
							goto l1083
						}
						position++
					}
				l1094:
					{
						position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1097
						}
						position++
						goto l1096
					l1097:
						position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
						if buffer[position] != rune('S') {
							goto l1083
						}
						position++
					}
				l1096:
					{
						position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1099
						}
						position++
						goto l1098
					l1099:
						position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
						if buffer[position] != rune('E') {
							goto l1083
						}
						position++
					}
				l1098:
					{
						position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1101
						}
						position++
						goto l1100
					l1101:
						position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
						if buffer[position] != rune('D') {
							goto l1083
						}
						position++
					}
				l1100:
					depth--
					add(rulePegText, position1085)
				}
				if !_rules[ruleAction77]() {
					goto l1083
				}
				depth--
				add(ruleUnpaused, position1084)
			}
			return true
		l1083:
			position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
			return false
		},
		/* 101 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
			{
				position1103 := position
				depth++
				{
					position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1105
					}
					goto l1104
				l1105:
					position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
					if !_rules[ruleInt]() {
						goto l1106
					}
					goto l1104
				l1106:
					position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
					if !_rules[ruleFloat]() {
						goto l1107
					}
					goto l1104
				l1107:
					position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
					if !_rules[ruleString]() {
						goto l1108
					}
					goto l1104
				l1108:
					position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
					if !_rules[ruleBlob]() {
						goto l1109
					}
					goto l1104
				l1109:
					position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
					if !_rules[ruleTimestamp]() {
						goto l1110
					}
					goto l1104
				l1110:
					position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
					if !_rules[ruleArray]() {
						goto l1111
					}
					goto l1104
				l1111:
					position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
					if !_rules[ruleMap]() {
						goto l1102
					}
				}
			l1104:
				depth--
				add(ruleType, position1103)
			}
			return true
		l1102:
			position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
			return false
		},
		/* 102 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action78)> */
		func() bool {
			position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
			{
				position1113 := position
				depth++
				{
					position1114 := position
					depth++
					{
						position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1116
						}
						position++
						goto l1115
					l1116:
						position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
						if buffer[position] != rune('B') {
							goto l1112
						}
						position++
					}
				l1115:
					{
						position1117, tokenIndex1117, depth1117 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1118
						}
						position++
						goto l1117
					l1118:
						position, tokenIndex, depth = position1117, tokenIndex1117, depth1117
						if buffer[position] != rune('O') {
							goto l1112
						}
						position++
					}
				l1117:
					{
						position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1120
						}
						position++
						goto l1119
					l1120:
						position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
						if buffer[position] != rune('O') {
							goto l1112
						}
						position++
					}
				l1119:
					{
						position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1122
						}
						position++
						goto l1121
					l1122:
						position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
						if buffer[position] != rune('L') {
							goto l1112
						}
						position++
					}
				l1121:
					depth--
					add(rulePegText, position1114)
				}
				if !_rules[ruleAction78]() {
					goto l1112
				}
				depth--
				add(ruleBool, position1113)
			}
			return true
		l1112:
			position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
			return false
		},
		/* 103 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action79)> */
		func() bool {
			position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
			{
				position1124 := position
				depth++
				{
					position1125 := position
					depth++
					{
						position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1127
						}
						position++
						goto l1126
					l1127:
						position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
						if buffer[position] != rune('I') {
							goto l1123
						}
						position++
					}
				l1126:
					{
						position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1129
						}
						position++
						goto l1128
					l1129:
						position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
						if buffer[position] != rune('N') {
							goto l1123
						}
						position++
					}
				l1128:
					{
						position1130, tokenIndex1130, depth1130 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1131
						}
						position++
						goto l1130
					l1131:
						position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
						if buffer[position] != rune('T') {
							goto l1123
						}
						position++
					}
				l1130:
					depth--
					add(rulePegText, position1125)
				}
				if !_rules[ruleAction79]() {
					goto l1123
				}
				depth--
				add(ruleInt, position1124)
			}
			return true
		l1123:
			position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
			return false
		},
		/* 104 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action80)> */
		func() bool {
			position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
			{
				position1133 := position
				depth++
				{
					position1134 := position
					depth++
					{
						position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
						if buffer[position] != rune('F') {
							goto l1132
						}
						position++
					}
				l1135:
					{
						position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
						if buffer[position] != rune('L') {
							goto l1132
						}
						position++
					}
				l1137:
					{
						position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1140
						}
						position++
						goto l1139
					l1140:
						position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
						if buffer[position] != rune('O') {
							goto l1132
						}
						position++
					}
				l1139:
					{
						position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1142
						}
						position++
						goto l1141
					l1142:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
						if buffer[position] != rune('A') {
							goto l1132
						}
						position++
					}
				l1141:
					{
						position1143, tokenIndex1143, depth1143 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1144
						}
						position++
						goto l1143
					l1144:
						position, tokenIndex, depth = position1143, tokenIndex1143, depth1143
						if buffer[position] != rune('T') {
							goto l1132
						}
						position++
					}
				l1143:
					depth--
					add(rulePegText, position1134)
				}
				if !_rules[ruleAction80]() {
					goto l1132
				}
				depth--
				add(ruleFloat, position1133)
			}
			return true
		l1132:
			position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
			return false
		},
		/* 105 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action81)> */
		func() bool {
			position1145, tokenIndex1145, depth1145 := position, tokenIndex, depth
			{
				position1146 := position
				depth++
				{
					position1147 := position
					depth++
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
							goto l1145
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
							goto l1145
						}
						position++
					}
				l1150:
					{
						position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1153
						}
						position++
						goto l1152
					l1153:
						position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
						if buffer[position] != rune('R') {
							goto l1145
						}
						position++
					}
				l1152:
					{
						position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1155
						}
						position++
						goto l1154
					l1155:
						position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
						if buffer[position] != rune('I') {
							goto l1145
						}
						position++
					}
				l1154:
					{
						position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1157
						}
						position++
						goto l1156
					l1157:
						position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
						if buffer[position] != rune('N') {
							goto l1145
						}
						position++
					}
				l1156:
					{
						position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1159
						}
						position++
						goto l1158
					l1159:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						if buffer[position] != rune('G') {
							goto l1145
						}
						position++
					}
				l1158:
					depth--
					add(rulePegText, position1147)
				}
				if !_rules[ruleAction81]() {
					goto l1145
				}
				depth--
				add(ruleString, position1146)
			}
			return true
		l1145:
			position, tokenIndex, depth = position1145, tokenIndex1145, depth1145
			return false
		},
		/* 106 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action82)> */
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
						if buffer[position] != rune('b') {
							goto l1164
						}
						position++
						goto l1163
					l1164:
						position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
						if buffer[position] != rune('B') {
							goto l1160
						}
						position++
					}
				l1163:
					{
						position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1166
						}
						position++
						goto l1165
					l1166:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						if buffer[position] != rune('L') {
							goto l1160
						}
						position++
					}
				l1165:
					{
						position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1168
						}
						position++
						goto l1167
					l1168:
						position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
						if buffer[position] != rune('O') {
							goto l1160
						}
						position++
					}
				l1167:
					{
						position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1170
						}
						position++
						goto l1169
					l1170:
						position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
						if buffer[position] != rune('B') {
							goto l1160
						}
						position++
					}
				l1169:
					depth--
					add(rulePegText, position1162)
				}
				if !_rules[ruleAction82]() {
					goto l1160
				}
				depth--
				add(ruleBlob, position1161)
			}
			return true
		l1160:
			position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
			return false
		},
		/* 107 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action83)> */
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
						if buffer[position] != rune('t') {
							goto l1175
						}
						position++
						goto l1174
					l1175:
						position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
						if buffer[position] != rune('T') {
							goto l1171
						}
						position++
					}
				l1174:
					{
						position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1177
						}
						position++
						goto l1176
					l1177:
						position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
						if buffer[position] != rune('I') {
							goto l1171
						}
						position++
					}
				l1176:
					{
						position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1179
						}
						position++
						goto l1178
					l1179:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('M') {
							goto l1171
						}
						position++
					}
				l1178:
					{
						position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1181
						}
						position++
						goto l1180
					l1181:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						if buffer[position] != rune('E') {
							goto l1171
						}
						position++
					}
				l1180:
					{
						position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1183
						}
						position++
						goto l1182
					l1183:
						position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
						if buffer[position] != rune('S') {
							goto l1171
						}
						position++
					}
				l1182:
					{
						position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1185
						}
						position++
						goto l1184
					l1185:
						position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
						if buffer[position] != rune('T') {
							goto l1171
						}
						position++
					}
				l1184:
					{
						position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1187
						}
						position++
						goto l1186
					l1187:
						position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
						if buffer[position] != rune('A') {
							goto l1171
						}
						position++
					}
				l1186:
					{
						position1188, tokenIndex1188, depth1188 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1189
						}
						position++
						goto l1188
					l1189:
						position, tokenIndex, depth = position1188, tokenIndex1188, depth1188
						if buffer[position] != rune('M') {
							goto l1171
						}
						position++
					}
				l1188:
					{
						position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1191
						}
						position++
						goto l1190
					l1191:
						position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						if buffer[position] != rune('P') {
							goto l1171
						}
						position++
					}
				l1190:
					depth--
					add(rulePegText, position1173)
				}
				if !_rules[ruleAction83]() {
					goto l1171
				}
				depth--
				add(ruleTimestamp, position1172)
			}
			return true
		l1171:
			position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
			return false
		},
		/* 108 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action84)> */
		func() bool {
			position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
			{
				position1193 := position
				depth++
				{
					position1194 := position
					depth++
					{
						position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1196
						}
						position++
						goto l1195
					l1196:
						position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
						if buffer[position] != rune('A') {
							goto l1192
						}
						position++
					}
				l1195:
					{
						position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1198
						}
						position++
						goto l1197
					l1198:
						position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
						if buffer[position] != rune('R') {
							goto l1192
						}
						position++
					}
				l1197:
					{
						position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1200
						}
						position++
						goto l1199
					l1200:
						position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
						if buffer[position] != rune('R') {
							goto l1192
						}
						position++
					}
				l1199:
					{
						position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1202
						}
						position++
						goto l1201
					l1202:
						position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
						if buffer[position] != rune('A') {
							goto l1192
						}
						position++
					}
				l1201:
					{
						position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1204
						}
						position++
						goto l1203
					l1204:
						position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
						if buffer[position] != rune('Y') {
							goto l1192
						}
						position++
					}
				l1203:
					depth--
					add(rulePegText, position1194)
				}
				if !_rules[ruleAction84]() {
					goto l1192
				}
				depth--
				add(ruleArray, position1193)
			}
			return true
		l1192:
			position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
			return false
		},
		/* 109 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action85)> */
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
						if buffer[position] != rune('m') {
							goto l1209
						}
						position++
						goto l1208
					l1209:
						position, tokenIndex, depth = position1208, tokenIndex1208, depth1208
						if buffer[position] != rune('M') {
							goto l1205
						}
						position++
					}
				l1208:
					{
						position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1211
						}
						position++
						goto l1210
					l1211:
						position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
						if buffer[position] != rune('A') {
							goto l1205
						}
						position++
					}
				l1210:
					{
						position1212, tokenIndex1212, depth1212 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1213
						}
						position++
						goto l1212
					l1213:
						position, tokenIndex, depth = position1212, tokenIndex1212, depth1212
						if buffer[position] != rune('P') {
							goto l1205
						}
						position++
					}
				l1212:
					depth--
					add(rulePegText, position1207)
				}
				if !_rules[ruleAction85]() {
					goto l1205
				}
				depth--
				add(ruleMap, position1206)
			}
			return true
		l1205:
			position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
			return false
		},
		/* 110 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action86)> */
		func() bool {
			position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
			{
				position1215 := position
				depth++
				{
					position1216 := position
					depth++
					{
						position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1218
						}
						position++
						goto l1217
					l1218:
						position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
						if buffer[position] != rune('O') {
							goto l1214
						}
						position++
					}
				l1217:
					{
						position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1220
						}
						position++
						goto l1219
					l1220:
						position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
						if buffer[position] != rune('R') {
							goto l1214
						}
						position++
					}
				l1219:
					depth--
					add(rulePegText, position1216)
				}
				if !_rules[ruleAction86]() {
					goto l1214
				}
				depth--
				add(ruleOr, position1215)
			}
			return true
		l1214:
			position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
			return false
		},
		/* 111 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action87)> */
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
						if buffer[position] != rune('a') {
							goto l1225
						}
						position++
						goto l1224
					l1225:
						position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
						if buffer[position] != rune('A') {
							goto l1221
						}
						position++
					}
				l1224:
					{
						position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1227
						}
						position++
						goto l1226
					l1227:
						position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
						if buffer[position] != rune('N') {
							goto l1221
						}
						position++
					}
				l1226:
					{
						position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1229
						}
						position++
						goto l1228
					l1229:
						position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
						if buffer[position] != rune('D') {
							goto l1221
						}
						position++
					}
				l1228:
					depth--
					add(rulePegText, position1223)
				}
				if !_rules[ruleAction87]() {
					goto l1221
				}
				depth--
				add(ruleAnd, position1222)
			}
			return true
		l1221:
			position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
			return false
		},
		/* 112 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action88)> */
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
						if buffer[position] != rune('n') {
							goto l1234
						}
						position++
						goto l1233
					l1234:
						position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
						if buffer[position] != rune('N') {
							goto l1230
						}
						position++
					}
				l1233:
					{
						position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1236
						}
						position++
						goto l1235
					l1236:
						position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
						if buffer[position] != rune('O') {
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
				if !_rules[ruleAction88]() {
					goto l1230
				}
				depth--
				add(ruleNot, position1231)
			}
			return true
		l1230:
			position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
			return false
		},
		/* 113 Equal <- <(<'='> Action89)> */
		func() bool {
			position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
			{
				position1240 := position
				depth++
				{
					position1241 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1239
					}
					position++
					depth--
					add(rulePegText, position1241)
				}
				if !_rules[ruleAction89]() {
					goto l1239
				}
				depth--
				add(ruleEqual, position1240)
			}
			return true
		l1239:
			position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
			return false
		},
		/* 114 Less <- <(<'<'> Action90)> */
		func() bool {
			position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
			{
				position1243 := position
				depth++
				{
					position1244 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1242
					}
					position++
					depth--
					add(rulePegText, position1244)
				}
				if !_rules[ruleAction90]() {
					goto l1242
				}
				depth--
				add(ruleLess, position1243)
			}
			return true
		l1242:
			position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
			return false
		},
		/* 115 LessOrEqual <- <(<('<' '=')> Action91)> */
		func() bool {
			position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
			{
				position1246 := position
				depth++
				{
					position1247 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1245
					}
					position++
					if buffer[position] != rune('=') {
						goto l1245
					}
					position++
					depth--
					add(rulePegText, position1247)
				}
				if !_rules[ruleAction91]() {
					goto l1245
				}
				depth--
				add(ruleLessOrEqual, position1246)
			}
			return true
		l1245:
			position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
			return false
		},
		/* 116 Greater <- <(<'>'> Action92)> */
		func() bool {
			position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
			{
				position1249 := position
				depth++
				{
					position1250 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1248
					}
					position++
					depth--
					add(rulePegText, position1250)
				}
				if !_rules[ruleAction92]() {
					goto l1248
				}
				depth--
				add(ruleGreater, position1249)
			}
			return true
		l1248:
			position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
			return false
		},
		/* 117 GreaterOrEqual <- <(<('>' '=')> Action93)> */
		func() bool {
			position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
			{
				position1252 := position
				depth++
				{
					position1253 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1251
					}
					position++
					if buffer[position] != rune('=') {
						goto l1251
					}
					position++
					depth--
					add(rulePegText, position1253)
				}
				if !_rules[ruleAction93]() {
					goto l1251
				}
				depth--
				add(ruleGreaterOrEqual, position1252)
			}
			return true
		l1251:
			position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
			return false
		},
		/* 118 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action94)> */
		func() bool {
			position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
			{
				position1255 := position
				depth++
				{
					position1256 := position
					depth++
					{
						position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1258
						}
						position++
						if buffer[position] != rune('=') {
							goto l1258
						}
						position++
						goto l1257
					l1258:
						position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
						if buffer[position] != rune('<') {
							goto l1254
						}
						position++
						if buffer[position] != rune('>') {
							goto l1254
						}
						position++
					}
				l1257:
					depth--
					add(rulePegText, position1256)
				}
				if !_rules[ruleAction94]() {
					goto l1254
				}
				depth--
				add(ruleNotEqual, position1255)
			}
			return true
		l1254:
			position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
			return false
		},
		/* 119 Concat <- <(<('|' '|')> Action95)> */
		func() bool {
			position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
			{
				position1260 := position
				depth++
				{
					position1261 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1259
					}
					position++
					if buffer[position] != rune('|') {
						goto l1259
					}
					position++
					depth--
					add(rulePegText, position1261)
				}
				if !_rules[ruleAction95]() {
					goto l1259
				}
				depth--
				add(ruleConcat, position1260)
			}
			return true
		l1259:
			position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
			return false
		},
		/* 120 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action96)> */
		func() bool {
			position1262, tokenIndex1262, depth1262 := position, tokenIndex, depth
			{
				position1263 := position
				depth++
				{
					position1264 := position
					depth++
					{
						position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1266
						}
						position++
						goto l1265
					l1266:
						position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
						if buffer[position] != rune('I') {
							goto l1262
						}
						position++
					}
				l1265:
					{
						position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1268
						}
						position++
						goto l1267
					l1268:
						position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
						if buffer[position] != rune('S') {
							goto l1262
						}
						position++
					}
				l1267:
					depth--
					add(rulePegText, position1264)
				}
				if !_rules[ruleAction96]() {
					goto l1262
				}
				depth--
				add(ruleIs, position1263)
			}
			return true
		l1262:
			position, tokenIndex, depth = position1262, tokenIndex1262, depth1262
			return false
		},
		/* 121 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action97)> */
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
					if !_rules[rulesp]() {
						goto l1269
					}
					{
						position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1277
						}
						position++
						goto l1276
					l1277:
						position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
						if buffer[position] != rune('N') {
							goto l1269
						}
						position++
					}
				l1276:
					{
						position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1279
						}
						position++
						goto l1278
					l1279:
						position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
						if buffer[position] != rune('O') {
							goto l1269
						}
						position++
					}
				l1278:
					{
						position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1281
						}
						position++
						goto l1280
					l1281:
						position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
						if buffer[position] != rune('T') {
							goto l1269
						}
						position++
					}
				l1280:
					depth--
					add(rulePegText, position1271)
				}
				if !_rules[ruleAction97]() {
					goto l1269
				}
				depth--
				add(ruleIsNot, position1270)
			}
			return true
		l1269:
			position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
			return false
		},
		/* 122 Plus <- <(<'+'> Action98)> */
		func() bool {
			position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
			{
				position1283 := position
				depth++
				{
					position1284 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1282
					}
					position++
					depth--
					add(rulePegText, position1284)
				}
				if !_rules[ruleAction98]() {
					goto l1282
				}
				depth--
				add(rulePlus, position1283)
			}
			return true
		l1282:
			position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
			return false
		},
		/* 123 Minus <- <(<'-'> Action99)> */
		func() bool {
			position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
			{
				position1286 := position
				depth++
				{
					position1287 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1285
					}
					position++
					depth--
					add(rulePegText, position1287)
				}
				if !_rules[ruleAction99]() {
					goto l1285
				}
				depth--
				add(ruleMinus, position1286)
			}
			return true
		l1285:
			position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
			return false
		},
		/* 124 Multiply <- <(<'*'> Action100)> */
		func() bool {
			position1288, tokenIndex1288, depth1288 := position, tokenIndex, depth
			{
				position1289 := position
				depth++
				{
					position1290 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1288
					}
					position++
					depth--
					add(rulePegText, position1290)
				}
				if !_rules[ruleAction100]() {
					goto l1288
				}
				depth--
				add(ruleMultiply, position1289)
			}
			return true
		l1288:
			position, tokenIndex, depth = position1288, tokenIndex1288, depth1288
			return false
		},
		/* 125 Divide <- <(<'/'> Action101)> */
		func() bool {
			position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
			{
				position1292 := position
				depth++
				{
					position1293 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1291
					}
					position++
					depth--
					add(rulePegText, position1293)
				}
				if !_rules[ruleAction101]() {
					goto l1291
				}
				depth--
				add(ruleDivide, position1292)
			}
			return true
		l1291:
			position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
			return false
		},
		/* 126 Modulo <- <(<'%'> Action102)> */
		func() bool {
			position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
			{
				position1295 := position
				depth++
				{
					position1296 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1294
					}
					position++
					depth--
					add(rulePegText, position1296)
				}
				if !_rules[ruleAction102]() {
					goto l1294
				}
				depth--
				add(ruleModulo, position1295)
			}
			return true
		l1294:
			position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
			return false
		},
		/* 127 UnaryMinus <- <(<'-'> Action103)> */
		func() bool {
			position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
			{
				position1298 := position
				depth++
				{
					position1299 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1297
					}
					position++
					depth--
					add(rulePegText, position1299)
				}
				if !_rules[ruleAction103]() {
					goto l1297
				}
				depth--
				add(ruleUnaryMinus, position1298)
			}
			return true
		l1297:
			position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
			return false
		},
		/* 128 Identifier <- <(<ident> Action104)> */
		func() bool {
			position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
			{
				position1301 := position
				depth++
				{
					position1302 := position
					depth++
					if !_rules[ruleident]() {
						goto l1300
					}
					depth--
					add(rulePegText, position1302)
				}
				if !_rules[ruleAction104]() {
					goto l1300
				}
				depth--
				add(ruleIdentifier, position1301)
			}
			return true
		l1300:
			position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
			return false
		},
		/* 129 TargetIdentifier <- <(<jsonPath> Action105)> */
		func() bool {
			position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
			{
				position1304 := position
				depth++
				{
					position1305 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1303
					}
					depth--
					add(rulePegText, position1305)
				}
				if !_rules[ruleAction105]() {
					goto l1303
				}
				depth--
				add(ruleTargetIdentifier, position1304)
			}
			return true
		l1303:
			position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
			return false
		},
		/* 130 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
			{
				position1307 := position
				depth++
				{
					position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1309
					}
					position++
					goto l1308
				l1309:
					position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1306
					}
					position++
				}
			l1308:
			l1310:
				{
					position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
					{
						position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1313
						}
						position++
						goto l1312
					l1313:
						position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1314
						}
						position++
						goto l1312
					l1314:
						position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1315
						}
						position++
						goto l1312
					l1315:
						position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
						if buffer[position] != rune('_') {
							goto l1311
						}
						position++
					}
				l1312:
					goto l1310
				l1311:
					position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
				}
				depth--
				add(ruleident, position1307)
			}
			return true
		l1306:
			position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
			return false
		},
		/* 131 jsonPath <- <(jsonPathHead jsonPathNonHead*)> */
		func() bool {
			position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
			{
				position1317 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1316
				}
			l1318:
				{
					position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
					if !_rules[rulejsonPathNonHead]() {
						goto l1319
					}
					goto l1318
				l1319:
					position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
				}
				depth--
				add(rulejsonPath, position1317)
			}
			return true
		l1316:
			position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
			return false
		},
		/* 132 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
			{
				position1321 := position
				depth++
				{
					position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1323
					}
					goto l1322
				l1323:
					position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1320
					}
				}
			l1322:
				depth--
				add(rulejsonPathHead, position1321)
			}
			return true
		l1320:
			position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
			return false
		},
		/* 133 jsonPathNonHead <- <(('.' jsonMapAccessString) / jsonMapAccessBracket / jsonArrayAccess)> */
		func() bool {
			position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
			{
				position1325 := position
				depth++
				{
					position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1327
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1327
					}
					goto l1326
				l1327:
					position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1328
					}
					goto l1326
				l1328:
					position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
					if !_rules[rulejsonArrayAccess]() {
						goto l1324
					}
				}
			l1326:
				depth--
				add(rulejsonPathNonHead, position1325)
			}
			return true
		l1324:
			position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
			return false
		},
		/* 134 jsonMapAccessString <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
			{
				position1330 := position
				depth++
				{
					position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1332
					}
					position++
					goto l1331
				l1332:
					position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1329
					}
					position++
				}
			l1331:
			l1333:
				{
					position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
					{
						position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1336
						}
						position++
						goto l1335
					l1336:
						position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1337
						}
						position++
						goto l1335
					l1337:
						position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1338
						}
						position++
						goto l1335
					l1338:
						position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
						if buffer[position] != rune('_') {
							goto l1334
						}
						position++
					}
				l1335:
					goto l1333
				l1334:
					position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
				}
				depth--
				add(rulejsonMapAccessString, position1330)
			}
			return true
		l1329:
			position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
			return false
		},
		/* 135 jsonMapAccessBracket <- <('[' '\'' (('\'' '\'') / (!'\'' .))* '\'' ']')> */
		func() bool {
			position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
			{
				position1340 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1339
				}
				position++
				if buffer[position] != rune('\'') {
					goto l1339
				}
				position++
			l1341:
				{
					position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
					{
						position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
						if buffer[position] != rune('\'') {
							goto l1344
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1344
						}
						position++
						goto l1343
					l1344:
						position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
						{
							position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l1345
							}
							position++
							goto l1342
						l1345:
							position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
						}
						if !matchDot() {
							goto l1342
						}
					}
				l1343:
					goto l1341
				l1342:
					position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
				}
				if buffer[position] != rune('\'') {
					goto l1339
				}
				position++
				if buffer[position] != rune(']') {
					goto l1339
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1340)
			}
			return true
		l1339:
			position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
			return false
		},
		/* 136 jsonArrayAccess <- <('[' [0-9]+ ']')> */
		func() bool {
			position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
			{
				position1347 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1346
				}
				position++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1346
				}
				position++
			l1348:
				{
					position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1349
					}
					position++
					goto l1348
				l1349:
					position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
				}
				if buffer[position] != rune(']') {
					goto l1346
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1347)
			}
			return true
		l1346:
			position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
			return false
		},
		/* 137 sp <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)*> */
		func() bool {
			{
				position1351 := position
				depth++
			l1352:
				{
					position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
					{
						position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1355
						}
						position++
						goto l1354
					l1355:
						position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
						if buffer[position] != rune('\t') {
							goto l1356
						}
						position++
						goto l1354
					l1356:
						position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
						if buffer[position] != rune('\n') {
							goto l1357
						}
						position++
						goto l1354
					l1357:
						position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
						if buffer[position] != rune('\r') {
							goto l1358
						}
						position++
						goto l1354
					l1358:
						position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
						if !_rules[rulecomment]() {
							goto l1359
						}
						goto l1354
					l1359:
						position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
						if !_rules[rulefinalComment]() {
							goto l1353
						}
					}
				l1354:
					goto l1352
				l1353:
					position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
				}
				depth--
				add(rulesp, position1351)
			}
			return true
		},
		/* 138 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
			{
				position1361 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1360
				}
				position++
				if buffer[position] != rune('-') {
					goto l1360
				}
				position++
			l1362:
				{
					position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
					{
						position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
						{
							position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1366
							}
							position++
							goto l1365
						l1366:
							position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
							if buffer[position] != rune('\n') {
								goto l1364
							}
							position++
						}
					l1365:
						goto l1363
					l1364:
						position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
					}
					if !matchDot() {
						goto l1363
					}
					goto l1362
				l1363:
					position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
				}
				{
					position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1368
					}
					position++
					goto l1367
				l1368:
					position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
					if buffer[position] != rune('\n') {
						goto l1360
					}
					position++
				}
			l1367:
				depth--
				add(rulecomment, position1361)
			}
			return true
		l1360:
			position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
			return false
		},
		/* 139 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
			{
				position1370 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1369
				}
				position++
				if buffer[position] != rune('-') {
					goto l1369
				}
				position++
			l1371:
				{
					position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
					{
						position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
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
								goto l1373
							}
							position++
						}
					l1374:
						goto l1372
					l1373:
						position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
					}
					if !matchDot() {
						goto l1372
					}
					goto l1371
				l1372:
					position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
				}
				{
					position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
					if !matchDot() {
						goto l1376
					}
					goto l1369
				l1376:
					position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
				}
				depth--
				add(rulefinalComment, position1370)
			}
			return true
		l1369:
			position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
			return false
		},
		nil,
		/* 142 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 143 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 144 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 145 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 146 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 147 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 148 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 149 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 150 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 151 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 152 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 153 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 154 Action12 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 155 Action13 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 156 Action14 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 157 Action15 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 158 Action16 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 159 Action17 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 160 Action18 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 161 Action19 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 162 Action20 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 163 Action21 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 164 Action22 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 165 Action23 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 166 Action24 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 167 Action25 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 168 Action26 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 169 Action27 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 170 Action28 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 171 Action29 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 172 Action30 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 173 Action31 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 174 Action32 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 175 Action33 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 176 Action34 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 177 Action35 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 178 Action36 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 179 Action37 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 180 Action38 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 181 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 182 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 183 Action41 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 184 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 185 Action43 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 186 Action44 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 187 Action45 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 188 Action46 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 189 Action47 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 190 Action48 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 191 Action49 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 192 Action50 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 193 Action51 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 194 Action52 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 195 Action53 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 196 Action54 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 197 Action55 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 198 Action56 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 199 Action57 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 200 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 201 Action59 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 202 Action60 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 203 Action61 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 204 Action62 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 205 Action63 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 206 Action64 <- <{
		    p.PushComponent(begin, end, NewWildcard(""))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 207 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 208 Action66 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 209 Action67 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 210 Action68 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 211 Action69 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 212 Action70 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 213 Action71 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 214 Action72 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 215 Action73 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 216 Action74 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 217 Action75 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 218 Action76 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 219 Action77 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 220 Action78 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 221 Action79 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 222 Action80 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 223 Action81 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 224 Action82 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 225 Action83 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 226 Action84 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 227 Action85 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 228 Action86 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 229 Action87 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 230 Action88 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 231 Action89 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 232 Action90 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 233 Action91 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 234 Action92 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 235 Action93 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 236 Action94 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 237 Action95 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 238 Action96 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 239 Action97 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 240 Action98 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 241 Action99 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 242 Action100 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 243 Action101 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 244 Action102 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 245 Action103 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 246 Action104 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 247 Action105 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
	}
	p.rules = _rules
}

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
	ruleStatements
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
	rulesp
	rulecomment
	ruleAction0
	rulePegText
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

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Statements",
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
	"sp",
	"comment",
	"Action0",
	"PegText",
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

type bqlPeg struct {
	parseStack

	Buffer string
	buffer []rune
	rules  [231]func() bool
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
	p *bqlPeg
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

func (p *bqlPeg) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *bqlPeg) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *bqlPeg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.AssembleSelect()

		case ruleAction1:

			p.AssembleSelectUnion(begin, end)

		case ruleAction2:

			p.AssembleCreateStreamAsSelect()

		case ruleAction3:

			p.AssembleCreateStreamAsSelectUnion()

		case ruleAction4:

			p.AssembleCreateSource()

		case ruleAction5:

			p.AssembleCreateSink()

		case ruleAction6:

			p.AssembleCreateState()

		case ruleAction7:

			p.AssembleUpdateState()

		case ruleAction8:

			p.AssembleUpdateSource()

		case ruleAction9:

			p.AssembleUpdateSink()

		case ruleAction10:

			p.AssembleInsertIntoSelect()

		case ruleAction11:

			p.AssembleInsertIntoFrom()

		case ruleAction12:

			p.AssemblePauseSource()

		case ruleAction13:

			p.AssembleResumeSource()

		case ruleAction14:

			p.AssembleRewindSource()

		case ruleAction15:

			p.AssembleDropSource()

		case ruleAction16:

			p.AssembleDropStream()

		case ruleAction17:

			p.AssembleDropSink()

		case ruleAction18:

			p.AssembleDropState()

		case ruleAction19:

			p.AssembleEmitter()

		case ruleAction20:

			p.AssembleProjections(begin, end)

		case ruleAction21:

			p.AssembleAlias()

		case ruleAction22:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction23:

			p.AssembleInterval()

		case ruleAction24:

			p.AssembleInterval()

		case ruleAction25:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction26:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction27:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction28:

			p.EnsureAliasedStreamWindow()

		case ruleAction29:

			p.AssembleAliasedStreamWindow()

		case ruleAction30:

			p.AssembleStreamWindow()

		case ruleAction31:

			p.AssembleUDSFFuncApp()

		case ruleAction32:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction33:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction34:

			p.AssembleSourceSinkParam()

		case ruleAction35:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction38:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction39:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction40:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction41:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction42:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction43:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction44:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction45:

			p.AssembleTypeCast(begin, end)

		case ruleAction46:

			p.AssembleTypeCast(begin, end)

		case ruleAction47:

			p.AssembleFuncApp()

		case ruleAction48:

			p.AssembleExpressions(begin, end)

		case ruleAction49:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction53:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction54:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction55:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction56:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction57:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction58:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction59:

			p.PushComponent(begin, end, NewWildcard(""))

		case ruleAction60:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction61:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction62:

			p.PushComponent(begin, end, Istream)

		case ruleAction63:

			p.PushComponent(begin, end, Dstream)

		case ruleAction64:

			p.PushComponent(begin, end, Rstream)

		case ruleAction65:

			p.PushComponent(begin, end, Tuples)

		case ruleAction66:

			p.PushComponent(begin, end, Seconds)

		case ruleAction67:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction68:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction69:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction70:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction71:

			p.PushComponent(begin, end, Yes)

		case ruleAction72:

			p.PushComponent(begin, end, No)

		case ruleAction73:

			p.PushComponent(begin, end, Bool)

		case ruleAction74:

			p.PushComponent(begin, end, Int)

		case ruleAction75:

			p.PushComponent(begin, end, Float)

		case ruleAction76:

			p.PushComponent(begin, end, String)

		case ruleAction77:

			p.PushComponent(begin, end, Blob)

		case ruleAction78:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction79:

			p.PushComponent(begin, end, Array)

		case ruleAction80:

			p.PushComponent(begin, end, Map)

		case ruleAction81:

			p.PushComponent(begin, end, Or)

		case ruleAction82:

			p.PushComponent(begin, end, And)

		case ruleAction83:

			p.PushComponent(begin, end, Not)

		case ruleAction84:

			p.PushComponent(begin, end, Equal)

		case ruleAction85:

			p.PushComponent(begin, end, Less)

		case ruleAction86:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction87:

			p.PushComponent(begin, end, Greater)

		case ruleAction88:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction89:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction90:

			p.PushComponent(begin, end, Concat)

		case ruleAction91:

			p.PushComponent(begin, end, Is)

		case ruleAction92:

			p.PushComponent(begin, end, IsNot)

		case ruleAction93:

			p.PushComponent(begin, end, Plus)

		case ruleAction94:

			p.PushComponent(begin, end, Minus)

		case ruleAction95:

			p.PushComponent(begin, end, Multiply)

		case ruleAction96:

			p.PushComponent(begin, end, Divide)

		case ruleAction97:

			p.PushComponent(begin, end, Modulo)

		case ruleAction98:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction99:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction100:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		}
	}
	_, _, _, _ = buffer, text, begin, end
}

func (p *bqlPeg) Init() {
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
		/* 0 Statements <- <(sp ((Statement sp ';' .*) / Statement) !.)> */
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
					if !_rules[ruleStatement]() {
						goto l3
					}
					if !_rules[rulesp]() {
						goto l3
					}
					if buffer[position] != rune(';') {
						goto l3
					}
					position++
				l4:
					{
						position5, tokenIndex5, depth5 := position, tokenIndex, depth
						if !matchDot() {
							goto l5
						}
						goto l4
					l5:
						position, tokenIndex, depth = position5, tokenIndex5, depth5
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleStatement]() {
						goto l0
					}
				}
			l2:
				{
					position6, tokenIndex6, depth6 := position, tokenIndex, depth
					if !matchDot() {
						goto l6
					}
					goto l0
				l6:
					position, tokenIndex, depth = position6, tokenIndex6, depth6
				}
				depth--
				add(ruleStatements, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Statement <- <(SelectUnionStmt / SelectStmt / SourceStmt / SinkStmt / StateStmt / StreamStmt)> */
		func() bool {
			position7, tokenIndex7, depth7 := position, tokenIndex, depth
			{
				position8 := position
				depth++
				{
					position9, tokenIndex9, depth9 := position, tokenIndex, depth
					if !_rules[ruleSelectUnionStmt]() {
						goto l10
					}
					goto l9
				l10:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleSelectStmt]() {
						goto l11
					}
					goto l9
				l11:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleSourceStmt]() {
						goto l12
					}
					goto l9
				l12:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleSinkStmt]() {
						goto l13
					}
					goto l9
				l13:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleStateStmt]() {
						goto l14
					}
					goto l9
				l14:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
					if !_rules[ruleStreamStmt]() {
						goto l7
					}
				}
			l9:
				depth--
				add(ruleStatement, position8)
			}
			return true
		l7:
			position, tokenIndex, depth = position7, tokenIndex7, depth7
			return false
		},
		/* 2 SourceStmt <- <(CreateSourceStmt / UpdateSourceStmt / DropSourceStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt)> */
		func() bool {
			position15, tokenIndex15, depth15 := position, tokenIndex, depth
			{
				position16 := position
				depth++
				{
					position17, tokenIndex17, depth17 := position, tokenIndex, depth
					if !_rules[ruleCreateSourceStmt]() {
						goto l18
					}
					goto l17
				l18:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleUpdateSourceStmt]() {
						goto l19
					}
					goto l17
				l19:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleDropSourceStmt]() {
						goto l20
					}
					goto l17
				l20:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[rulePauseSourceStmt]() {
						goto l21
					}
					goto l17
				l21:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleResumeSourceStmt]() {
						goto l22
					}
					goto l17
				l22:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
					if !_rules[ruleRewindSourceStmt]() {
						goto l15
					}
				}
			l17:
				depth--
				add(ruleSourceStmt, position16)
			}
			return true
		l15:
			position, tokenIndex, depth = position15, tokenIndex15, depth15
			return false
		},
		/* 3 SinkStmt <- <(CreateSinkStmt / UpdateSinkStmt / DropSinkStmt)> */
		func() bool {
			position23, tokenIndex23, depth23 := position, tokenIndex, depth
			{
				position24 := position
				depth++
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					if !_rules[ruleCreateSinkStmt]() {
						goto l26
					}
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if !_rules[ruleUpdateSinkStmt]() {
						goto l27
					}
					goto l25
				l27:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if !_rules[ruleDropSinkStmt]() {
						goto l23
					}
				}
			l25:
				depth--
				add(ruleSinkStmt, position24)
			}
			return true
		l23:
			position, tokenIndex, depth = position23, tokenIndex23, depth23
			return false
		},
		/* 4 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt)> */
		func() bool {
			position28, tokenIndex28, depth28 := position, tokenIndex, depth
			{
				position29 := position
				depth++
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if !_rules[ruleCreateStateStmt]() {
						goto l31
					}
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleUpdateStateStmt]() {
						goto l32
					}
					goto l30
				l32:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if !_rules[ruleDropStateStmt]() {
						goto l28
					}
				}
			l30:
				depth--
				add(ruleStateStmt, position29)
			}
			return true
		l28:
			position, tokenIndex, depth = position28, tokenIndex28, depth28
			return false
		},
		/* 5 StreamStmt <- <(CreateStreamAsSelectUnionStmt / CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoSelectStmt / InsertIntoFromStmt)> */
		func() bool {
			position33, tokenIndex33, depth33 := position, tokenIndex, depth
			{
				position34 := position
				depth++
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l36
					}
					goto l35
				l36:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l37
					}
					goto l35
				l37:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleDropStreamStmt]() {
						goto l38
					}
					goto l35
				l38:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l39
					}
					goto l35
				l39:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l33
					}
				}
			l35:
				depth--
				add(ruleStreamStmt, position34)
			}
			return true
		l33:
			position, tokenIndex, depth = position33, tokenIndex33, depth33
			return false
		},
		/* 6 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
		func() bool {
			position40, tokenIndex40, depth40 := position, tokenIndex, depth
			{
				position41 := position
				depth++
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('S') {
						goto l40
					}
					position++
				}
			l42:
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l45
					}
					position++
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if buffer[position] != rune('E') {
						goto l40
					}
					position++
				}
			l44:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('L') {
						goto l40
					}
					position++
				}
			l46:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('E') {
						goto l40
					}
					position++
				}
			l48:
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('C') {
						goto l40
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('T') {
						goto l40
					}
					position++
				}
			l52:
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleEmitter]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleProjections]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleWindowedFrom]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleFilter]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleGrouping]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleHaving]() {
					goto l40
				}
				if !_rules[rulesp]() {
					goto l40
				}
				if !_rules[ruleAction0]() {
					goto l40
				}
				depth--
				add(ruleSelectStmt, position41)
			}
			return true
		l40:
			position, tokenIndex, depth = position40, tokenIndex40, depth40
			return false
		},
		/* 7 SelectUnionStmt <- <(<(SelectStmt (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action1)> */
		func() bool {
			position54, tokenIndex54, depth54 := position, tokenIndex, depth
			{
				position55 := position
				depth++
				{
					position56 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l54
					}
					{
						position59, tokenIndex59, depth59 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l60
						}
						position++
						goto l59
					l60:
						position, tokenIndex, depth = position59, tokenIndex59, depth59
						if buffer[position] != rune('U') {
							goto l54
						}
						position++
					}
				l59:
					{
						position61, tokenIndex61, depth61 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l62
						}
						position++
						goto l61
					l62:
						position, tokenIndex, depth = position61, tokenIndex61, depth61
						if buffer[position] != rune('N') {
							goto l54
						}
						position++
					}
				l61:
					{
						position63, tokenIndex63, depth63 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l64
						}
						position++
						goto l63
					l64:
						position, tokenIndex, depth = position63, tokenIndex63, depth63
						if buffer[position] != rune('I') {
							goto l54
						}
						position++
					}
				l63:
					{
						position65, tokenIndex65, depth65 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l66
						}
						position++
						goto l65
					l66:
						position, tokenIndex, depth = position65, tokenIndex65, depth65
						if buffer[position] != rune('O') {
							goto l54
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
							goto l54
						}
						position++
					}
				l67:
					if !_rules[rulesp]() {
						goto l54
					}
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l70
						}
						position++
						goto l69
					l70:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
						if buffer[position] != rune('A') {
							goto l54
						}
						position++
					}
				l69:
					{
						position71, tokenIndex71, depth71 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l72
						}
						position++
						goto l71
					l72:
						position, tokenIndex, depth = position71, tokenIndex71, depth71
						if buffer[position] != rune('L') {
							goto l54
						}
						position++
					}
				l71:
					{
						position73, tokenIndex73, depth73 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l74
						}
						position++
						goto l73
					l74:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
						if buffer[position] != rune('L') {
							goto l54
						}
						position++
					}
				l73:
					if !_rules[rulesp]() {
						goto l54
					}
					if !_rules[ruleSelectStmt]() {
						goto l54
					}
				l57:
					{
						position58, tokenIndex58, depth58 := position, tokenIndex, depth
						{
							position75, tokenIndex75, depth75 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l76
							}
							position++
							goto l75
						l76:
							position, tokenIndex, depth = position75, tokenIndex75, depth75
							if buffer[position] != rune('U') {
								goto l58
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
								goto l58
							}
							position++
						}
					l77:
						{
							position79, tokenIndex79, depth79 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l80
							}
							position++
							goto l79
						l80:
							position, tokenIndex, depth = position79, tokenIndex79, depth79
							if buffer[position] != rune('I') {
								goto l58
							}
							position++
						}
					l79:
						{
							position81, tokenIndex81, depth81 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l82
							}
							position++
							goto l81
						l82:
							position, tokenIndex, depth = position81, tokenIndex81, depth81
							if buffer[position] != rune('O') {
								goto l58
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
								goto l58
							}
							position++
						}
					l83:
						if !_rules[rulesp]() {
							goto l58
						}
						{
							position85, tokenIndex85, depth85 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l86
							}
							position++
							goto l85
						l86:
							position, tokenIndex, depth = position85, tokenIndex85, depth85
							if buffer[position] != rune('A') {
								goto l58
							}
							position++
						}
					l85:
						{
							position87, tokenIndex87, depth87 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l88
							}
							position++
							goto l87
						l88:
							position, tokenIndex, depth = position87, tokenIndex87, depth87
							if buffer[position] != rune('L') {
								goto l58
							}
							position++
						}
					l87:
						{
							position89, tokenIndex89, depth89 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l90
							}
							position++
							goto l89
						l90:
							position, tokenIndex, depth = position89, tokenIndex89, depth89
							if buffer[position] != rune('L') {
								goto l58
							}
							position++
						}
					l89:
						if !_rules[rulesp]() {
							goto l58
						}
						if !_rules[ruleSelectStmt]() {
							goto l58
						}
						goto l57
					l58:
						position, tokenIndex, depth = position58, tokenIndex58, depth58
					}
					depth--
					add(rulePegText, position56)
				}
				if !_rules[ruleAction1]() {
					goto l54
				}
				depth--
				add(ruleSelectUnionStmt, position55)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 8 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action2)> */
		func() bool {
			position91, tokenIndex91, depth91 := position, tokenIndex, depth
			{
				position92 := position
				depth++
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l94
					}
					position++
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					if buffer[position] != rune('C') {
						goto l91
					}
					position++
				}
			l93:
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
					if buffer[position] != rune('R') {
						goto l91
					}
					position++
				}
			l95:
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l98
					}
					position++
					goto l97
				l98:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
					if buffer[position] != rune('E') {
						goto l91
					}
					position++
				}
			l97:
				{
					position99, tokenIndex99, depth99 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l100
					}
					position++
					goto l99
				l100:
					position, tokenIndex, depth = position99, tokenIndex99, depth99
					if buffer[position] != rune('A') {
						goto l91
					}
					position++
				}
			l99:
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l102
					}
					position++
					goto l101
				l102:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
					if buffer[position] != rune('T') {
						goto l91
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
						goto l91
					}
					position++
				}
			l103:
				if !_rules[rulesp]() {
					goto l91
				}
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l106
					}
					position++
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if buffer[position] != rune('S') {
						goto l91
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
						goto l91
					}
					position++
				}
			l107:
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l110
					}
					position++
					goto l109
				l110:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
					if buffer[position] != rune('R') {
						goto l91
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
						goto l91
					}
					position++
				}
			l111:
				{
					position113, tokenIndex113, depth113 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l114
					}
					position++
					goto l113
				l114:
					position, tokenIndex, depth = position113, tokenIndex113, depth113
					if buffer[position] != rune('A') {
						goto l91
					}
					position++
				}
			l113:
				{
					position115, tokenIndex115, depth115 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l116
					}
					position++
					goto l115
				l116:
					position, tokenIndex, depth = position115, tokenIndex115, depth115
					if buffer[position] != rune('M') {
						goto l91
					}
					position++
				}
			l115:
				if !_rules[rulesp]() {
					goto l91
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l91
				}
				if !_rules[rulesp]() {
					goto l91
				}
				{
					position117, tokenIndex117, depth117 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l118
					}
					position++
					goto l117
				l118:
					position, tokenIndex, depth = position117, tokenIndex117, depth117
					if buffer[position] != rune('A') {
						goto l91
					}
					position++
				}
			l117:
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l120
					}
					position++
					goto l119
				l120:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
					if buffer[position] != rune('S') {
						goto l91
					}
					position++
				}
			l119:
				if !_rules[rulesp]() {
					goto l91
				}
				if !_rules[ruleSelectStmt]() {
					goto l91
				}
				if !_rules[ruleAction2]() {
					goto l91
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position92)
			}
			return true
		l91:
			position, tokenIndex, depth = position91, tokenIndex91, depth91
			return false
		},
		/* 9 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action3)> */
		func() bool {
			position121, tokenIndex121, depth121 := position, tokenIndex, depth
			{
				position122 := position
				depth++
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l124
					}
					position++
					goto l123
				l124:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
					if buffer[position] != rune('C') {
						goto l121
					}
					position++
				}
			l123:
				{
					position125, tokenIndex125, depth125 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l126
					}
					position++
					goto l125
				l126:
					position, tokenIndex, depth = position125, tokenIndex125, depth125
					if buffer[position] != rune('R') {
						goto l121
					}
					position++
				}
			l125:
				{
					position127, tokenIndex127, depth127 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l128
					}
					position++
					goto l127
				l128:
					position, tokenIndex, depth = position127, tokenIndex127, depth127
					if buffer[position] != rune('E') {
						goto l121
					}
					position++
				}
			l127:
				{
					position129, tokenIndex129, depth129 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l130
					}
					position++
					goto l129
				l130:
					position, tokenIndex, depth = position129, tokenIndex129, depth129
					if buffer[position] != rune('A') {
						goto l121
					}
					position++
				}
			l129:
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l132
					}
					position++
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if buffer[position] != rune('T') {
						goto l121
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
						goto l121
					}
					position++
				}
			l133:
				if !_rules[rulesp]() {
					goto l121
				}
				{
					position135, tokenIndex135, depth135 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l136
					}
					position++
					goto l135
				l136:
					position, tokenIndex, depth = position135, tokenIndex135, depth135
					if buffer[position] != rune('S') {
						goto l121
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
						goto l121
					}
					position++
				}
			l137:
				{
					position139, tokenIndex139, depth139 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l140
					}
					position++
					goto l139
				l140:
					position, tokenIndex, depth = position139, tokenIndex139, depth139
					if buffer[position] != rune('R') {
						goto l121
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
						goto l121
					}
					position++
				}
			l141:
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l144
					}
					position++
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					if buffer[position] != rune('A') {
						goto l121
					}
					position++
				}
			l143:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l146
					}
					position++
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					if buffer[position] != rune('M') {
						goto l121
					}
					position++
				}
			l145:
				if !_rules[rulesp]() {
					goto l121
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l121
				}
				if !_rules[rulesp]() {
					goto l121
				}
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l148
					}
					position++
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if buffer[position] != rune('A') {
						goto l121
					}
					position++
				}
			l147:
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l150
					}
					position++
					goto l149
				l150:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
					if buffer[position] != rune('S') {
						goto l121
					}
					position++
				}
			l149:
				if !_rules[rulesp]() {
					goto l121
				}
				if !_rules[ruleSelectUnionStmt]() {
					goto l121
				}
				if !_rules[ruleAction3]() {
					goto l121
				}
				depth--
				add(ruleCreateStreamAsSelectUnionStmt, position122)
			}
			return true
		l121:
			position, tokenIndex, depth = position121, tokenIndex121, depth121
			return false
		},
		/* 10 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position151, tokenIndex151, depth151 := position, tokenIndex, depth
			{
				position152 := position
				depth++
				{
					position153, tokenIndex153, depth153 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l154
					}
					position++
					goto l153
				l154:
					position, tokenIndex, depth = position153, tokenIndex153, depth153
					if buffer[position] != rune('C') {
						goto l151
					}
					position++
				}
			l153:
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l156
					}
					position++
					goto l155
				l156:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
					if buffer[position] != rune('R') {
						goto l151
					}
					position++
				}
			l155:
				{
					position157, tokenIndex157, depth157 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l158
					}
					position++
					goto l157
				l158:
					position, tokenIndex, depth = position157, tokenIndex157, depth157
					if buffer[position] != rune('E') {
						goto l151
					}
					position++
				}
			l157:
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l160
					}
					position++
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if buffer[position] != rune('A') {
						goto l151
					}
					position++
				}
			l159:
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l162
					}
					position++
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if buffer[position] != rune('T') {
						goto l151
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
						goto l151
					}
					position++
				}
			l163:
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[rulePausedOpt]() {
					goto l151
				}
				if !_rules[rulesp]() {
					goto l151
				}
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					if buffer[position] != rune('S') {
						goto l151
					}
					position++
				}
			l165:
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l168
					}
					position++
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('O') {
						goto l151
					}
					position++
				}
			l167:
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l170
					}
					position++
					goto l169
				l170:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
					if buffer[position] != rune('U') {
						goto l151
					}
					position++
				}
			l169:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l172
					}
					position++
					goto l171
				l172:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
					if buffer[position] != rune('R') {
						goto l151
					}
					position++
				}
			l171:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l174
					}
					position++
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if buffer[position] != rune('C') {
						goto l151
					}
					position++
				}
			l173:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l176
					}
					position++
					goto l175
				l176:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
					if buffer[position] != rune('E') {
						goto l151
					}
					position++
				}
			l175:
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l151
				}
				if !_rules[rulesp]() {
					goto l151
				}
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l178
					}
					position++
					goto l177
				l178:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
					if buffer[position] != rune('T') {
						goto l151
					}
					position++
				}
			l177:
				{
					position179, tokenIndex179, depth179 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l180
					}
					position++
					goto l179
				l180:
					position, tokenIndex, depth = position179, tokenIndex179, depth179
					if buffer[position] != rune('Y') {
						goto l151
					}
					position++
				}
			l179:
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l182
					}
					position++
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if buffer[position] != rune('P') {
						goto l151
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
						goto l151
					}
					position++
				}
			l183:
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[ruleSourceSinkType]() {
					goto l151
				}
				if !_rules[rulesp]() {
					goto l151
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l151
				}
				if !_rules[ruleAction4]() {
					goto l151
				}
				depth--
				add(ruleCreateSourceStmt, position152)
			}
			return true
		l151:
			position, tokenIndex, depth = position151, tokenIndex151, depth151
			return false
		},
		/* 11 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action5)> */
		func() bool {
			position185, tokenIndex185, depth185 := position, tokenIndex, depth
			{
				position186 := position
				depth++
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l188
					}
					position++
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					if buffer[position] != rune('C') {
						goto l185
					}
					position++
				}
			l187:
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l190
					}
					position++
					goto l189
				l190:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
					if buffer[position] != rune('R') {
						goto l185
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
						goto l185
					}
					position++
				}
			l191:
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l194
					}
					position++
					goto l193
				l194:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
					if buffer[position] != rune('A') {
						goto l185
					}
					position++
				}
			l193:
				{
					position195, tokenIndex195, depth195 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l196
					}
					position++
					goto l195
				l196:
					position, tokenIndex, depth = position195, tokenIndex195, depth195
					if buffer[position] != rune('T') {
						goto l185
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
						goto l185
					}
					position++
				}
			l197:
				if !_rules[rulesp]() {
					goto l185
				}
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l200
					}
					position++
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if buffer[position] != rune('S') {
						goto l185
					}
					position++
				}
			l199:
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l202
					}
					position++
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if buffer[position] != rune('I') {
						goto l185
					}
					position++
				}
			l201:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l204
					}
					position++
					goto l203
				l204:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
					if buffer[position] != rune('N') {
						goto l185
					}
					position++
				}
			l203:
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l206
					}
					position++
					goto l205
				l206:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
					if buffer[position] != rune('K') {
						goto l185
					}
					position++
				}
			l205:
				if !_rules[rulesp]() {
					goto l185
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l185
				}
				if !_rules[rulesp]() {
					goto l185
				}
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l208
					}
					position++
					goto l207
				l208:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
					if buffer[position] != rune('T') {
						goto l185
					}
					position++
				}
			l207:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l210
					}
					position++
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if buffer[position] != rune('Y') {
						goto l185
					}
					position++
				}
			l209:
				{
					position211, tokenIndex211, depth211 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l212
					}
					position++
					goto l211
				l212:
					position, tokenIndex, depth = position211, tokenIndex211, depth211
					if buffer[position] != rune('P') {
						goto l185
					}
					position++
				}
			l211:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l214
					}
					position++
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if buffer[position] != rune('E') {
						goto l185
					}
					position++
				}
			l213:
				if !_rules[rulesp]() {
					goto l185
				}
				if !_rules[ruleSourceSinkType]() {
					goto l185
				}
				if !_rules[rulesp]() {
					goto l185
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l185
				}
				if !_rules[ruleAction5]() {
					goto l185
				}
				depth--
				add(ruleCreateSinkStmt, position186)
			}
			return true
		l185:
			position, tokenIndex, depth = position185, tokenIndex185, depth185
			return false
		},
		/* 12 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action6)> */
		func() bool {
			position215, tokenIndex215, depth215 := position, tokenIndex, depth
			{
				position216 := position
				depth++
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l218
					}
					position++
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if buffer[position] != rune('C') {
						goto l215
					}
					position++
				}
			l217:
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l220
					}
					position++
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if buffer[position] != rune('R') {
						goto l215
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
						goto l215
					}
					position++
				}
			l221:
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l224
					}
					position++
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					if buffer[position] != rune('A') {
						goto l215
					}
					position++
				}
			l223:
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l226
					}
					position++
					goto l225
				l226:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
					if buffer[position] != rune('T') {
						goto l215
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
						goto l215
					}
					position++
				}
			l227:
				if !_rules[rulesp]() {
					goto l215
				}
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l230
					}
					position++
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if buffer[position] != rune('S') {
						goto l215
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
						goto l215
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
						goto l215
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
						goto l215
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
						goto l215
					}
					position++
				}
			l237:
				if !_rules[rulesp]() {
					goto l215
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l215
				}
				if !_rules[rulesp]() {
					goto l215
				}
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
						goto l215
					}
					position++
				}
			l239:
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l242
					}
					position++
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if buffer[position] != rune('Y') {
						goto l215
					}
					position++
				}
			l241:
				{
					position243, tokenIndex243, depth243 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l244
					}
					position++
					goto l243
				l244:
					position, tokenIndex, depth = position243, tokenIndex243, depth243
					if buffer[position] != rune('P') {
						goto l215
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
						goto l215
					}
					position++
				}
			l245:
				if !_rules[rulesp]() {
					goto l215
				}
				if !_rules[ruleSourceSinkType]() {
					goto l215
				}
				if !_rules[rulesp]() {
					goto l215
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l215
				}
				if !_rules[ruleAction6]() {
					goto l215
				}
				depth--
				add(ruleCreateStateStmt, position216)
			}
			return true
		l215:
			position, tokenIndex, depth = position215, tokenIndex215, depth215
			return false
		},
		/* 13 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action7)> */
		func() bool {
			position247, tokenIndex247, depth247 := position, tokenIndex, depth
			{
				position248 := position
				depth++
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l250
					}
					position++
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if buffer[position] != rune('U') {
						goto l247
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
						goto l247
					}
					position++
				}
			l251:
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l254
					}
					position++
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if buffer[position] != rune('D') {
						goto l247
					}
					position++
				}
			l253:
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l256
					}
					position++
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('A') {
						goto l247
					}
					position++
				}
			l255:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l258
					}
					position++
					goto l257
				l258:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
					if buffer[position] != rune('T') {
						goto l247
					}
					position++
				}
			l257:
				{
					position259, tokenIndex259, depth259 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l260
					}
					position++
					goto l259
				l260:
					position, tokenIndex, depth = position259, tokenIndex259, depth259
					if buffer[position] != rune('E') {
						goto l247
					}
					position++
				}
			l259:
				if !_rules[rulesp]() {
					goto l247
				}
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l262
					}
					position++
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if buffer[position] != rune('S') {
						goto l247
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
						goto l247
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
						goto l247
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
						goto l247
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
						goto l247
					}
					position++
				}
			l269:
				if !_rules[rulesp]() {
					goto l247
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l247
				}
				if !_rules[rulesp]() {
					goto l247
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l247
				}
				if !_rules[ruleAction7]() {
					goto l247
				}
				depth--
				add(ruleUpdateStateStmt, position248)
			}
			return true
		l247:
			position, tokenIndex, depth = position247, tokenIndex247, depth247
			return false
		},
		/* 14 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action8)> */
		func() bool {
			position271, tokenIndex271, depth271 := position, tokenIndex, depth
			{
				position272 := position
				depth++
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('U') {
						goto l271
					}
					position++
				}
			l273:
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l276
					}
					position++
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					if buffer[position] != rune('P') {
						goto l271
					}
					position++
				}
			l275:
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l278
					}
					position++
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					if buffer[position] != rune('D') {
						goto l271
					}
					position++
				}
			l277:
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l280
					}
					position++
					goto l279
				l280:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
					if buffer[position] != rune('A') {
						goto l271
					}
					position++
				}
			l279:
				{
					position281, tokenIndex281, depth281 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l282
					}
					position++
					goto l281
				l282:
					position, tokenIndex, depth = position281, tokenIndex281, depth281
					if buffer[position] != rune('T') {
						goto l271
					}
					position++
				}
			l281:
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l284
					}
					position++
					goto l283
				l284:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
					if buffer[position] != rune('E') {
						goto l271
					}
					position++
				}
			l283:
				if !_rules[rulesp]() {
					goto l271
				}
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l286
					}
					position++
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('S') {
						goto l271
					}
					position++
				}
			l285:
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l288
					}
					position++
					goto l287
				l288:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
					if buffer[position] != rune('O') {
						goto l271
					}
					position++
				}
			l287:
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l290
					}
					position++
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					if buffer[position] != rune('U') {
						goto l271
					}
					position++
				}
			l289:
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l292
					}
					position++
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					if buffer[position] != rune('R') {
						goto l271
					}
					position++
				}
			l291:
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l294
					}
					position++
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					if buffer[position] != rune('C') {
						goto l271
					}
					position++
				}
			l293:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l296
					}
					position++
					goto l295
				l296:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
					if buffer[position] != rune('E') {
						goto l271
					}
					position++
				}
			l295:
				if !_rules[rulesp]() {
					goto l271
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l271
				}
				if !_rules[rulesp]() {
					goto l271
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l271
				}
				if !_rules[ruleAction8]() {
					goto l271
				}
				depth--
				add(ruleUpdateSourceStmt, position272)
			}
			return true
		l271:
			position, tokenIndex, depth = position271, tokenIndex271, depth271
			return false
		},
		/* 15 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position297, tokenIndex297, depth297 := position, tokenIndex, depth
			{
				position298 := position
				depth++
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
						goto l297
					}
					position++
				}
			l299:
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l302
					}
					position++
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					if buffer[position] != rune('P') {
						goto l297
					}
					position++
				}
			l301:
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('D') {
						goto l297
					}
					position++
				}
			l303:
				{
					position305, tokenIndex305, depth305 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l306
					}
					position++
					goto l305
				l306:
					position, tokenIndex, depth = position305, tokenIndex305, depth305
					if buffer[position] != rune('A') {
						goto l297
					}
					position++
				}
			l305:
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l308
					}
					position++
					goto l307
				l308:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
					if buffer[position] != rune('T') {
						goto l297
					}
					position++
				}
			l307:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l310
					}
					position++
					goto l309
				l310:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
					if buffer[position] != rune('E') {
						goto l297
					}
					position++
				}
			l309:
				if !_rules[rulesp]() {
					goto l297
				}
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l312
					}
					position++
					goto l311
				l312:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
					if buffer[position] != rune('S') {
						goto l297
					}
					position++
				}
			l311:
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l314
					}
					position++
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					if buffer[position] != rune('I') {
						goto l297
					}
					position++
				}
			l313:
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l316
					}
					position++
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					if buffer[position] != rune('N') {
						goto l297
					}
					position++
				}
			l315:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l318
					}
					position++
					goto l317
				l318:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
					if buffer[position] != rune('K') {
						goto l297
					}
					position++
				}
			l317:
				if !_rules[rulesp]() {
					goto l297
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l297
				}
				if !_rules[rulesp]() {
					goto l297
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l297
				}
				if !_rules[ruleAction9]() {
					goto l297
				}
				depth--
				add(ruleUpdateSinkStmt, position298)
			}
			return true
		l297:
			position, tokenIndex, depth = position297, tokenIndex297, depth297
			return false
		},
		/* 16 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action10)> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
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
						goto l319
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
						goto l319
					}
					position++
				}
			l323:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l326
					}
					position++
					goto l325
				l326:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
					if buffer[position] != rune('S') {
						goto l319
					}
					position++
				}
			l325:
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l328
					}
					position++
					goto l327
				l328:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if buffer[position] != rune('E') {
						goto l319
					}
					position++
				}
			l327:
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l330
					}
					position++
					goto l329
				l330:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
					if buffer[position] != rune('R') {
						goto l319
					}
					position++
				}
			l329:
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l332
					}
					position++
					goto l331
				l332:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
					if buffer[position] != rune('T') {
						goto l319
					}
					position++
				}
			l331:
				if !_rules[rulesp]() {
					goto l319
				}
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l334
					}
					position++
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if buffer[position] != rune('I') {
						goto l319
					}
					position++
				}
			l333:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l336
					}
					position++
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if buffer[position] != rune('N') {
						goto l319
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
						goto l319
					}
					position++
				}
			l337:
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l340
					}
					position++
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					if buffer[position] != rune('O') {
						goto l319
					}
					position++
				}
			l339:
				if !_rules[rulesp]() {
					goto l319
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l319
				}
				if !_rules[rulesp]() {
					goto l319
				}
				if !_rules[ruleSelectStmt]() {
					goto l319
				}
				if !_rules[ruleAction10]() {
					goto l319
				}
				depth--
				add(ruleInsertIntoSelectStmt, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 17 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action11)> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
			{
				position342 := position
				depth++
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
						goto l341
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
						goto l341
					}
					position++
				}
			l345:
				{
					position347, tokenIndex347, depth347 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l348
					}
					position++
					goto l347
				l348:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if buffer[position] != rune('S') {
						goto l341
					}
					position++
				}
			l347:
				{
					position349, tokenIndex349, depth349 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l350
					}
					position++
					goto l349
				l350:
					position, tokenIndex, depth = position349, tokenIndex349, depth349
					if buffer[position] != rune('E') {
						goto l341
					}
					position++
				}
			l349:
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l352
					}
					position++
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					if buffer[position] != rune('R') {
						goto l341
					}
					position++
				}
			l351:
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l354
					}
					position++
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if buffer[position] != rune('T') {
						goto l341
					}
					position++
				}
			l353:
				if !_rules[rulesp]() {
					goto l341
				}
				{
					position355, tokenIndex355, depth355 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l356
					}
					position++
					goto l355
				l356:
					position, tokenIndex, depth = position355, tokenIndex355, depth355
					if buffer[position] != rune('I') {
						goto l341
					}
					position++
				}
			l355:
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l358
					}
					position++
					goto l357
				l358:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
					if buffer[position] != rune('N') {
						goto l341
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
						goto l341
					}
					position++
				}
			l359:
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l362
					}
					position++
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if buffer[position] != rune('O') {
						goto l341
					}
					position++
				}
			l361:
				if !_rules[rulesp]() {
					goto l341
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l341
				}
				if !_rules[rulesp]() {
					goto l341
				}
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l364
					}
					position++
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					if buffer[position] != rune('F') {
						goto l341
					}
					position++
				}
			l363:
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('R') {
						goto l341
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
						goto l341
					}
					position++
				}
			l367:
				{
					position369, tokenIndex369, depth369 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l370
					}
					position++
					goto l369
				l370:
					position, tokenIndex, depth = position369, tokenIndex369, depth369
					if buffer[position] != rune('M') {
						goto l341
					}
					position++
				}
			l369:
				if !_rules[rulesp]() {
					goto l341
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l341
				}
				if !_rules[ruleAction11]() {
					goto l341
				}
				depth--
				add(ruleInsertIntoFromStmt, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
			return false
		},
		/* 18 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action12)> */
		func() bool {
			position371, tokenIndex371, depth371 := position, tokenIndex, depth
			{
				position372 := position
				depth++
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('P') {
						goto l371
					}
					position++
				}
			l373:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l376
					}
					position++
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					if buffer[position] != rune('A') {
						goto l371
					}
					position++
				}
			l375:
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l378
					}
					position++
					goto l377
				l378:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
					if buffer[position] != rune('U') {
						goto l371
					}
					position++
				}
			l377:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('S') {
						goto l371
					}
					position++
				}
			l379:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('E') {
						goto l371
					}
					position++
				}
			l381:
				if !_rules[rulesp]() {
					goto l371
				}
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('S') {
						goto l371
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('O') {
						goto l371
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
						goto l371
					}
					position++
				}
			l387:
				{
					position389, tokenIndex389, depth389 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l390
					}
					position++
					goto l389
				l390:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
					if buffer[position] != rune('R') {
						goto l371
					}
					position++
				}
			l389:
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
					if buffer[position] != rune('C') {
						goto l371
					}
					position++
				}
			l391:
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l394
					}
					position++
					goto l393
				l394:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
					if buffer[position] != rune('E') {
						goto l371
					}
					position++
				}
			l393:
				if !_rules[rulesp]() {
					goto l371
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l371
				}
				if !_rules[ruleAction12]() {
					goto l371
				}
				depth--
				add(rulePauseSourceStmt, position372)
			}
			return true
		l371:
			position, tokenIndex, depth = position371, tokenIndex371, depth371
			return false
		},
		/* 19 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position395, tokenIndex395, depth395 := position, tokenIndex, depth
			{
				position396 := position
				depth++
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
						goto l395
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
						goto l395
					}
					position++
				}
			l399:
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l402
					}
					position++
					goto l401
				l402:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if buffer[position] != rune('S') {
						goto l395
					}
					position++
				}
			l401:
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l404
					}
					position++
					goto l403
				l404:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
					if buffer[position] != rune('U') {
						goto l395
					}
					position++
				}
			l403:
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l406
					}
					position++
					goto l405
				l406:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
					if buffer[position] != rune('M') {
						goto l395
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
						goto l395
					}
					position++
				}
			l407:
				if !_rules[rulesp]() {
					goto l395
				}
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
						goto l395
					}
					position++
				}
			l409:
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l412
					}
					position++
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if buffer[position] != rune('O') {
						goto l395
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
						goto l395
					}
					position++
				}
			l413:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l416
					}
					position++
					goto l415
				l416:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
					if buffer[position] != rune('R') {
						goto l395
					}
					position++
				}
			l415:
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l418
					}
					position++
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if buffer[position] != rune('C') {
						goto l395
					}
					position++
				}
			l417:
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l420
					}
					position++
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if buffer[position] != rune('E') {
						goto l395
					}
					position++
				}
			l419:
				if !_rules[rulesp]() {
					goto l395
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l395
				}
				if !_rules[ruleAction13]() {
					goto l395
				}
				depth--
				add(ruleResumeSourceStmt, position396)
			}
			return true
		l395:
			position, tokenIndex, depth = position395, tokenIndex395, depth395
			return false
		},
		/* 20 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position421, tokenIndex421, depth421 := position, tokenIndex, depth
			{
				position422 := position
				depth++
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
						goto l421
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
						goto l421
					}
					position++
				}
			l425:
				{
					position427, tokenIndex427, depth427 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l428
					}
					position++
					goto l427
				l428:
					position, tokenIndex, depth = position427, tokenIndex427, depth427
					if buffer[position] != rune('W') {
						goto l421
					}
					position++
				}
			l427:
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l430
					}
					position++
					goto l429
				l430:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if buffer[position] != rune('I') {
						goto l421
					}
					position++
				}
			l429:
				{
					position431, tokenIndex431, depth431 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l432
					}
					position++
					goto l431
				l432:
					position, tokenIndex, depth = position431, tokenIndex431, depth431
					if buffer[position] != rune('N') {
						goto l421
					}
					position++
				}
			l431:
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l434
					}
					position++
					goto l433
				l434:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if buffer[position] != rune('D') {
						goto l421
					}
					position++
				}
			l433:
				if !_rules[rulesp]() {
					goto l421
				}
				{
					position435, tokenIndex435, depth435 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l436
					}
					position++
					goto l435
				l436:
					position, tokenIndex, depth = position435, tokenIndex435, depth435
					if buffer[position] != rune('S') {
						goto l421
					}
					position++
				}
			l435:
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l438
					}
					position++
					goto l437
				l438:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
					if buffer[position] != rune('O') {
						goto l421
					}
					position++
				}
			l437:
				{
					position439, tokenIndex439, depth439 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l440
					}
					position++
					goto l439
				l440:
					position, tokenIndex, depth = position439, tokenIndex439, depth439
					if buffer[position] != rune('U') {
						goto l421
					}
					position++
				}
			l439:
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l442
					}
					position++
					goto l441
				l442:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
					if buffer[position] != rune('R') {
						goto l421
					}
					position++
				}
			l441:
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l444
					}
					position++
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if buffer[position] != rune('C') {
						goto l421
					}
					position++
				}
			l443:
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l446
					}
					position++
					goto l445
				l446:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
					if buffer[position] != rune('E') {
						goto l421
					}
					position++
				}
			l445:
				if !_rules[rulesp]() {
					goto l421
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l421
				}
				if !_rules[ruleAction14]() {
					goto l421
				}
				depth--
				add(ruleRewindSourceStmt, position422)
			}
			return true
		l421:
			position, tokenIndex, depth = position421, tokenIndex421, depth421
			return false
		},
		/* 21 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position447, tokenIndex447, depth447 := position, tokenIndex, depth
			{
				position448 := position
				depth++
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l450
					}
					position++
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					if buffer[position] != rune('D') {
						goto l447
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
						goto l447
					}
					position++
				}
			l451:
				{
					position453, tokenIndex453, depth453 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l454
					}
					position++
					goto l453
				l454:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
					if buffer[position] != rune('O') {
						goto l447
					}
					position++
				}
			l453:
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l456
					}
					position++
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if buffer[position] != rune('P') {
						goto l447
					}
					position++
				}
			l455:
				if !_rules[rulesp]() {
					goto l447
				}
				{
					position457, tokenIndex457, depth457 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l458
					}
					position++
					goto l457
				l458:
					position, tokenIndex, depth = position457, tokenIndex457, depth457
					if buffer[position] != rune('S') {
						goto l447
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
						goto l447
					}
					position++
				}
			l459:
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l462
					}
					position++
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('U') {
						goto l447
					}
					position++
				}
			l461:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('R') {
						goto l447
					}
					position++
				}
			l463:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l466
					}
					position++
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					if buffer[position] != rune('C') {
						goto l447
					}
					position++
				}
			l465:
				{
					position467, tokenIndex467, depth467 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l468
					}
					position++
					goto l467
				l468:
					position, tokenIndex, depth = position467, tokenIndex467, depth467
					if buffer[position] != rune('E') {
						goto l447
					}
					position++
				}
			l467:
				if !_rules[rulesp]() {
					goto l447
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l447
				}
				if !_rules[ruleAction15]() {
					goto l447
				}
				depth--
				add(ruleDropSourceStmt, position448)
			}
			return true
		l447:
			position, tokenIndex, depth = position447, tokenIndex447, depth447
			return false
		},
		/* 22 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action16)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				{
					position471, tokenIndex471, depth471 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l472
					}
					position++
					goto l471
				l472:
					position, tokenIndex, depth = position471, tokenIndex471, depth471
					if buffer[position] != rune('D') {
						goto l469
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
						goto l469
					}
					position++
				}
			l473:
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l476
					}
					position++
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if buffer[position] != rune('O') {
						goto l469
					}
					position++
				}
			l475:
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l478
					}
					position++
					goto l477
				l478:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
					if buffer[position] != rune('P') {
						goto l469
					}
					position++
				}
			l477:
				if !_rules[rulesp]() {
					goto l469
				}
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l480
					}
					position++
					goto l479
				l480:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if buffer[position] != rune('S') {
						goto l469
					}
					position++
				}
			l479:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
					if buffer[position] != rune('T') {
						goto l469
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
						goto l469
					}
					position++
				}
			l483:
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l486
					}
					position++
					goto l485
				l486:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					if buffer[position] != rune('E') {
						goto l469
					}
					position++
				}
			l485:
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l488
					}
					position++
					goto l487
				l488:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
					if buffer[position] != rune('A') {
						goto l469
					}
					position++
				}
			l487:
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l490
					}
					position++
					goto l489
				l490:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
					if buffer[position] != rune('M') {
						goto l469
					}
					position++
				}
			l489:
				if !_rules[rulesp]() {
					goto l469
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l469
				}
				if !_rules[ruleAction16]() {
					goto l469
				}
				depth--
				add(ruleDropStreamStmt, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 23 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action17)> */
		func() bool {
			position491, tokenIndex491, depth491 := position, tokenIndex, depth
			{
				position492 := position
				depth++
				{
					position493, tokenIndex493, depth493 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l494
					}
					position++
					goto l493
				l494:
					position, tokenIndex, depth = position493, tokenIndex493, depth493
					if buffer[position] != rune('D') {
						goto l491
					}
					position++
				}
			l493:
				{
					position495, tokenIndex495, depth495 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l496
					}
					position++
					goto l495
				l496:
					position, tokenIndex, depth = position495, tokenIndex495, depth495
					if buffer[position] != rune('R') {
						goto l491
					}
					position++
				}
			l495:
				{
					position497, tokenIndex497, depth497 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l498
					}
					position++
					goto l497
				l498:
					position, tokenIndex, depth = position497, tokenIndex497, depth497
					if buffer[position] != rune('O') {
						goto l491
					}
					position++
				}
			l497:
				{
					position499, tokenIndex499, depth499 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l500
					}
					position++
					goto l499
				l500:
					position, tokenIndex, depth = position499, tokenIndex499, depth499
					if buffer[position] != rune('P') {
						goto l491
					}
					position++
				}
			l499:
				if !_rules[rulesp]() {
					goto l491
				}
				{
					position501, tokenIndex501, depth501 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l502
					}
					position++
					goto l501
				l502:
					position, tokenIndex, depth = position501, tokenIndex501, depth501
					if buffer[position] != rune('S') {
						goto l491
					}
					position++
				}
			l501:
				{
					position503, tokenIndex503, depth503 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l504
					}
					position++
					goto l503
				l504:
					position, tokenIndex, depth = position503, tokenIndex503, depth503
					if buffer[position] != rune('I') {
						goto l491
					}
					position++
				}
			l503:
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l506
					}
					position++
					goto l505
				l506:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
					if buffer[position] != rune('N') {
						goto l491
					}
					position++
				}
			l505:
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('K') {
						goto l491
					}
					position++
				}
			l507:
				if !_rules[rulesp]() {
					goto l491
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l491
				}
				if !_rules[ruleAction17]() {
					goto l491
				}
				depth--
				add(ruleDropSinkStmt, position492)
			}
			return true
		l491:
			position, tokenIndex, depth = position491, tokenIndex491, depth491
			return false
		},
		/* 24 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action18)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if buffer[position] != rune('D') {
						goto l509
					}
					position++
				}
			l511:
				{
					position513, tokenIndex513, depth513 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l514
					}
					position++
					goto l513
				l514:
					position, tokenIndex, depth = position513, tokenIndex513, depth513
					if buffer[position] != rune('R') {
						goto l509
					}
					position++
				}
			l513:
				{
					position515, tokenIndex515, depth515 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l516
					}
					position++
					goto l515
				l516:
					position, tokenIndex, depth = position515, tokenIndex515, depth515
					if buffer[position] != rune('O') {
						goto l509
					}
					position++
				}
			l515:
				{
					position517, tokenIndex517, depth517 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l518
					}
					position++
					goto l517
				l518:
					position, tokenIndex, depth = position517, tokenIndex517, depth517
					if buffer[position] != rune('P') {
						goto l509
					}
					position++
				}
			l517:
				if !_rules[rulesp]() {
					goto l509
				}
				{
					position519, tokenIndex519, depth519 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l520
					}
					position++
					goto l519
				l520:
					position, tokenIndex, depth = position519, tokenIndex519, depth519
					if buffer[position] != rune('S') {
						goto l509
					}
					position++
				}
			l519:
				{
					position521, tokenIndex521, depth521 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l522
					}
					position++
					goto l521
				l522:
					position, tokenIndex, depth = position521, tokenIndex521, depth521
					if buffer[position] != rune('T') {
						goto l509
					}
					position++
				}
			l521:
				{
					position523, tokenIndex523, depth523 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l524
					}
					position++
					goto l523
				l524:
					position, tokenIndex, depth = position523, tokenIndex523, depth523
					if buffer[position] != rune('A') {
						goto l509
					}
					position++
				}
			l523:
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l526
					}
					position++
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if buffer[position] != rune('T') {
						goto l509
					}
					position++
				}
			l525:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l528
					}
					position++
					goto l527
				l528:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
					if buffer[position] != rune('E') {
						goto l509
					}
					position++
				}
			l527:
				if !_rules[rulesp]() {
					goto l509
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l509
				}
				if !_rules[ruleAction18]() {
					goto l509
				}
				depth--
				add(ruleDropStateStmt, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 25 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) Action19)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l532
					}
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleDSTREAM]() {
						goto l533
					}
					goto l531
				l533:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					if !_rules[ruleRSTREAM]() {
						goto l529
					}
				}
			l531:
				if !_rules[ruleAction19]() {
					goto l529
				}
				depth--
				add(ruleEmitter, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 26 Projections <- <(<(Projection sp (',' sp Projection)*)> Action20)> */
		func() bool {
			position534, tokenIndex534, depth534 := position, tokenIndex, depth
			{
				position535 := position
				depth++
				{
					position536 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l534
					}
					if !_rules[rulesp]() {
						goto l534
					}
				l537:
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l538
						}
						position++
						if !_rules[rulesp]() {
							goto l538
						}
						if !_rules[ruleProjection]() {
							goto l538
						}
						goto l537
					l538:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
					}
					depth--
					add(rulePegText, position536)
				}
				if !_rules[ruleAction20]() {
					goto l534
				}
				depth--
				add(ruleProjections, position535)
			}
			return true
		l534:
			position, tokenIndex, depth = position534, tokenIndex534, depth534
			return false
		},
		/* 27 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position539, tokenIndex539, depth539 := position, tokenIndex, depth
			{
				position540 := position
				depth++
				{
					position541, tokenIndex541, depth541 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l542
					}
					goto l541
				l542:
					position, tokenIndex, depth = position541, tokenIndex541, depth541
					if !_rules[ruleExpression]() {
						goto l543
					}
					goto l541
				l543:
					position, tokenIndex, depth = position541, tokenIndex541, depth541
					if !_rules[ruleWildcard]() {
						goto l539
					}
				}
			l541:
				depth--
				add(ruleProjection, position540)
			}
			return true
		l539:
			position, tokenIndex, depth = position539, tokenIndex539, depth539
			return false
		},
		/* 28 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action21)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				{
					position546, tokenIndex546, depth546 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l547
					}
					goto l546
				l547:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
					if !_rules[ruleWildcard]() {
						goto l544
					}
				}
			l546:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position548, tokenIndex548, depth548 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l549
					}
					position++
					goto l548
				l549:
					position, tokenIndex, depth = position548, tokenIndex548, depth548
					if buffer[position] != rune('A') {
						goto l544
					}
					position++
				}
			l548:
				{
					position550, tokenIndex550, depth550 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l551
					}
					position++
					goto l550
				l551:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if buffer[position] != rune('S') {
						goto l544
					}
					position++
				}
			l550:
				if !_rules[rulesp]() {
					goto l544
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l544
				}
				if !_rules[ruleAction21]() {
					goto l544
				}
				depth--
				add(ruleAliasExpression, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 29 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action22)> */
		func() bool {
			position552, tokenIndex552, depth552 := position, tokenIndex, depth
			{
				position553 := position
				depth++
				{
					position554 := position
					depth++
					{
						position555, tokenIndex555, depth555 := position, tokenIndex, depth
						{
							position557, tokenIndex557, depth557 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l558
							}
							position++
							goto l557
						l558:
							position, tokenIndex, depth = position557, tokenIndex557, depth557
							if buffer[position] != rune('F') {
								goto l555
							}
							position++
						}
					l557:
						{
							position559, tokenIndex559, depth559 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l560
							}
							position++
							goto l559
						l560:
							position, tokenIndex, depth = position559, tokenIndex559, depth559
							if buffer[position] != rune('R') {
								goto l555
							}
							position++
						}
					l559:
						{
							position561, tokenIndex561, depth561 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l562
							}
							position++
							goto l561
						l562:
							position, tokenIndex, depth = position561, tokenIndex561, depth561
							if buffer[position] != rune('O') {
								goto l555
							}
							position++
						}
					l561:
						{
							position563, tokenIndex563, depth563 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l564
							}
							position++
							goto l563
						l564:
							position, tokenIndex, depth = position563, tokenIndex563, depth563
							if buffer[position] != rune('M') {
								goto l555
							}
							position++
						}
					l563:
						if !_rules[rulesp]() {
							goto l555
						}
						if !_rules[ruleRelations]() {
							goto l555
						}
						if !_rules[rulesp]() {
							goto l555
						}
						goto l556
					l555:
						position, tokenIndex, depth = position555, tokenIndex555, depth555
					}
				l556:
					depth--
					add(rulePegText, position554)
				}
				if !_rules[ruleAction22]() {
					goto l552
				}
				depth--
				add(ruleWindowedFrom, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 30 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position565, tokenIndex565, depth565 := position, tokenIndex, depth
			{
				position566 := position
				depth++
				{
					position567, tokenIndex567, depth567 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l568
					}
					goto l567
				l568:
					position, tokenIndex, depth = position567, tokenIndex567, depth567
					if !_rules[ruleTuplesInterval]() {
						goto l565
					}
				}
			l567:
				depth--
				add(ruleInterval, position566)
			}
			return true
		l565:
			position, tokenIndex, depth = position565, tokenIndex565, depth565
			return false
		},
		/* 31 TimeInterval <- <(NumericLiteral sp (SECONDS / MILLISECONDS) Action23)> */
		func() bool {
			position569, tokenIndex569, depth569 := position, tokenIndex, depth
			{
				position570 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l569
				}
				if !_rules[rulesp]() {
					goto l569
				}
				{
					position571, tokenIndex571, depth571 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l572
					}
					goto l571
				l572:
					position, tokenIndex, depth = position571, tokenIndex571, depth571
					if !_rules[ruleMILLISECONDS]() {
						goto l569
					}
				}
			l571:
				if !_rules[ruleAction23]() {
					goto l569
				}
				depth--
				add(ruleTimeInterval, position570)
			}
			return true
		l569:
			position, tokenIndex, depth = position569, tokenIndex569, depth569
			return false
		},
		/* 32 TuplesInterval <- <(NumericLiteral sp TUPLES Action24)> */
		func() bool {
			position573, tokenIndex573, depth573 := position, tokenIndex, depth
			{
				position574 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l573
				}
				if !_rules[rulesp]() {
					goto l573
				}
				if !_rules[ruleTUPLES]() {
					goto l573
				}
				if !_rules[ruleAction24]() {
					goto l573
				}
				depth--
				add(ruleTuplesInterval, position574)
			}
			return true
		l573:
			position, tokenIndex, depth = position573, tokenIndex573, depth573
			return false
		},
		/* 33 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position575, tokenIndex575, depth575 := position, tokenIndex, depth
			{
				position576 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l575
				}
				if !_rules[rulesp]() {
					goto l575
				}
			l577:
				{
					position578, tokenIndex578, depth578 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l578
					}
					position++
					if !_rules[rulesp]() {
						goto l578
					}
					if !_rules[ruleRelationLike]() {
						goto l578
					}
					goto l577
				l578:
					position, tokenIndex, depth = position578, tokenIndex578, depth578
				}
				depth--
				add(ruleRelations, position576)
			}
			return true
		l575:
			position, tokenIndex, depth = position575, tokenIndex575, depth575
			return false
		},
		/* 34 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action25)> */
		func() bool {
			position579, tokenIndex579, depth579 := position, tokenIndex, depth
			{
				position580 := position
				depth++
				{
					position581 := position
					depth++
					{
						position582, tokenIndex582, depth582 := position, tokenIndex, depth
						{
							position584, tokenIndex584, depth584 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l585
							}
							position++
							goto l584
						l585:
							position, tokenIndex, depth = position584, tokenIndex584, depth584
							if buffer[position] != rune('W') {
								goto l582
							}
							position++
						}
					l584:
						{
							position586, tokenIndex586, depth586 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l587
							}
							position++
							goto l586
						l587:
							position, tokenIndex, depth = position586, tokenIndex586, depth586
							if buffer[position] != rune('H') {
								goto l582
							}
							position++
						}
					l586:
						{
							position588, tokenIndex588, depth588 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l589
							}
							position++
							goto l588
						l589:
							position, tokenIndex, depth = position588, tokenIndex588, depth588
							if buffer[position] != rune('E') {
								goto l582
							}
							position++
						}
					l588:
						{
							position590, tokenIndex590, depth590 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l591
							}
							position++
							goto l590
						l591:
							position, tokenIndex, depth = position590, tokenIndex590, depth590
							if buffer[position] != rune('R') {
								goto l582
							}
							position++
						}
					l590:
						{
							position592, tokenIndex592, depth592 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l593
							}
							position++
							goto l592
						l593:
							position, tokenIndex, depth = position592, tokenIndex592, depth592
							if buffer[position] != rune('E') {
								goto l582
							}
							position++
						}
					l592:
						if !_rules[rulesp]() {
							goto l582
						}
						if !_rules[ruleExpression]() {
							goto l582
						}
						goto l583
					l582:
						position, tokenIndex, depth = position582, tokenIndex582, depth582
					}
				l583:
					depth--
					add(rulePegText, position581)
				}
				if !_rules[ruleAction25]() {
					goto l579
				}
				depth--
				add(ruleFilter, position580)
			}
			return true
		l579:
			position, tokenIndex, depth = position579, tokenIndex579, depth579
			return false
		},
		/* 35 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action26)> */
		func() bool {
			position594, tokenIndex594, depth594 := position, tokenIndex, depth
			{
				position595 := position
				depth++
				{
					position596 := position
					depth++
					{
						position597, tokenIndex597, depth597 := position, tokenIndex, depth
						{
							position599, tokenIndex599, depth599 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l600
							}
							position++
							goto l599
						l600:
							position, tokenIndex, depth = position599, tokenIndex599, depth599
							if buffer[position] != rune('G') {
								goto l597
							}
							position++
						}
					l599:
						{
							position601, tokenIndex601, depth601 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l602
							}
							position++
							goto l601
						l602:
							position, tokenIndex, depth = position601, tokenIndex601, depth601
							if buffer[position] != rune('R') {
								goto l597
							}
							position++
						}
					l601:
						{
							position603, tokenIndex603, depth603 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l604
							}
							position++
							goto l603
						l604:
							position, tokenIndex, depth = position603, tokenIndex603, depth603
							if buffer[position] != rune('O') {
								goto l597
							}
							position++
						}
					l603:
						{
							position605, tokenIndex605, depth605 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l606
							}
							position++
							goto l605
						l606:
							position, tokenIndex, depth = position605, tokenIndex605, depth605
							if buffer[position] != rune('U') {
								goto l597
							}
							position++
						}
					l605:
						{
							position607, tokenIndex607, depth607 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l608
							}
							position++
							goto l607
						l608:
							position, tokenIndex, depth = position607, tokenIndex607, depth607
							if buffer[position] != rune('P') {
								goto l597
							}
							position++
						}
					l607:
						if !_rules[rulesp]() {
							goto l597
						}
						{
							position609, tokenIndex609, depth609 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l610
							}
							position++
							goto l609
						l610:
							position, tokenIndex, depth = position609, tokenIndex609, depth609
							if buffer[position] != rune('B') {
								goto l597
							}
							position++
						}
					l609:
						{
							position611, tokenIndex611, depth611 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l612
							}
							position++
							goto l611
						l612:
							position, tokenIndex, depth = position611, tokenIndex611, depth611
							if buffer[position] != rune('Y') {
								goto l597
							}
							position++
						}
					l611:
						if !_rules[rulesp]() {
							goto l597
						}
						if !_rules[ruleGroupList]() {
							goto l597
						}
						goto l598
					l597:
						position, tokenIndex, depth = position597, tokenIndex597, depth597
					}
				l598:
					depth--
					add(rulePegText, position596)
				}
				if !_rules[ruleAction26]() {
					goto l594
				}
				depth--
				add(ruleGrouping, position595)
			}
			return true
		l594:
			position, tokenIndex, depth = position594, tokenIndex594, depth594
			return false
		},
		/* 36 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position613, tokenIndex613, depth613 := position, tokenIndex, depth
			{
				position614 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l613
				}
				if !_rules[rulesp]() {
					goto l613
				}
			l615:
				{
					position616, tokenIndex616, depth616 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l616
					}
					position++
					if !_rules[rulesp]() {
						goto l616
					}
					if !_rules[ruleExpression]() {
						goto l616
					}
					goto l615
				l616:
					position, tokenIndex, depth = position616, tokenIndex616, depth616
				}
				depth--
				add(ruleGroupList, position614)
			}
			return true
		l613:
			position, tokenIndex, depth = position613, tokenIndex613, depth613
			return false
		},
		/* 37 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action27)> */
		func() bool {
			position617, tokenIndex617, depth617 := position, tokenIndex, depth
			{
				position618 := position
				depth++
				{
					position619 := position
					depth++
					{
						position620, tokenIndex620, depth620 := position, tokenIndex, depth
						{
							position622, tokenIndex622, depth622 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l623
							}
							position++
							goto l622
						l623:
							position, tokenIndex, depth = position622, tokenIndex622, depth622
							if buffer[position] != rune('H') {
								goto l620
							}
							position++
						}
					l622:
						{
							position624, tokenIndex624, depth624 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l625
							}
							position++
							goto l624
						l625:
							position, tokenIndex, depth = position624, tokenIndex624, depth624
							if buffer[position] != rune('A') {
								goto l620
							}
							position++
						}
					l624:
						{
							position626, tokenIndex626, depth626 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l627
							}
							position++
							goto l626
						l627:
							position, tokenIndex, depth = position626, tokenIndex626, depth626
							if buffer[position] != rune('V') {
								goto l620
							}
							position++
						}
					l626:
						{
							position628, tokenIndex628, depth628 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l629
							}
							position++
							goto l628
						l629:
							position, tokenIndex, depth = position628, tokenIndex628, depth628
							if buffer[position] != rune('I') {
								goto l620
							}
							position++
						}
					l628:
						{
							position630, tokenIndex630, depth630 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l631
							}
							position++
							goto l630
						l631:
							position, tokenIndex, depth = position630, tokenIndex630, depth630
							if buffer[position] != rune('N') {
								goto l620
							}
							position++
						}
					l630:
						{
							position632, tokenIndex632, depth632 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l633
							}
							position++
							goto l632
						l633:
							position, tokenIndex, depth = position632, tokenIndex632, depth632
							if buffer[position] != rune('G') {
								goto l620
							}
							position++
						}
					l632:
						if !_rules[rulesp]() {
							goto l620
						}
						if !_rules[ruleExpression]() {
							goto l620
						}
						goto l621
					l620:
						position, tokenIndex, depth = position620, tokenIndex620, depth620
					}
				l621:
					depth--
					add(rulePegText, position619)
				}
				if !_rules[ruleAction27]() {
					goto l617
				}
				depth--
				add(ruleHaving, position618)
			}
			return true
		l617:
			position, tokenIndex, depth = position617, tokenIndex617, depth617
			return false
		},
		/* 38 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action28))> */
		func() bool {
			position634, tokenIndex634, depth634 := position, tokenIndex, depth
			{
				position635 := position
				depth++
				{
					position636, tokenIndex636, depth636 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l637
					}
					goto l636
				l637:
					position, tokenIndex, depth = position636, tokenIndex636, depth636
					if !_rules[ruleStreamWindow]() {
						goto l634
					}
					if !_rules[ruleAction28]() {
						goto l634
					}
				}
			l636:
				depth--
				add(ruleRelationLike, position635)
			}
			return true
		l634:
			position, tokenIndex, depth = position634, tokenIndex634, depth634
			return false
		},
		/* 39 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action29)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l638
				}
				if !_rules[rulesp]() {
					goto l638
				}
				{
					position640, tokenIndex640, depth640 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l641
					}
					position++
					goto l640
				l641:
					position, tokenIndex, depth = position640, tokenIndex640, depth640
					if buffer[position] != rune('A') {
						goto l638
					}
					position++
				}
			l640:
				{
					position642, tokenIndex642, depth642 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l643
					}
					position++
					goto l642
				l643:
					position, tokenIndex, depth = position642, tokenIndex642, depth642
					if buffer[position] != rune('S') {
						goto l638
					}
					position++
				}
			l642:
				if !_rules[rulesp]() {
					goto l638
				}
				if !_rules[ruleIdentifier]() {
					goto l638
				}
				if !_rules[ruleAction29]() {
					goto l638
				}
				depth--
				add(ruleAliasedStreamWindow, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 40 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action30)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l644
				}
				if !_rules[rulesp]() {
					goto l644
				}
				if buffer[position] != rune('[') {
					goto l644
				}
				position++
				if !_rules[rulesp]() {
					goto l644
				}
				{
					position646, tokenIndex646, depth646 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l647
					}
					position++
					goto l646
				l647:
					position, tokenIndex, depth = position646, tokenIndex646, depth646
					if buffer[position] != rune('R') {
						goto l644
					}
					position++
				}
			l646:
				{
					position648, tokenIndex648, depth648 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l649
					}
					position++
					goto l648
				l649:
					position, tokenIndex, depth = position648, tokenIndex648, depth648
					if buffer[position] != rune('A') {
						goto l644
					}
					position++
				}
			l648:
				{
					position650, tokenIndex650, depth650 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l651
					}
					position++
					goto l650
				l651:
					position, tokenIndex, depth = position650, tokenIndex650, depth650
					if buffer[position] != rune('N') {
						goto l644
					}
					position++
				}
			l650:
				{
					position652, tokenIndex652, depth652 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l653
					}
					position++
					goto l652
				l653:
					position, tokenIndex, depth = position652, tokenIndex652, depth652
					if buffer[position] != rune('G') {
						goto l644
					}
					position++
				}
			l652:
				{
					position654, tokenIndex654, depth654 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l655
					}
					position++
					goto l654
				l655:
					position, tokenIndex, depth = position654, tokenIndex654, depth654
					if buffer[position] != rune('E') {
						goto l644
					}
					position++
				}
			l654:
				if !_rules[rulesp]() {
					goto l644
				}
				if !_rules[ruleInterval]() {
					goto l644
				}
				if !_rules[rulesp]() {
					goto l644
				}
				if buffer[position] != rune(']') {
					goto l644
				}
				position++
				if !_rules[ruleAction30]() {
					goto l644
				}
				depth--
				add(ruleStreamWindow, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 41 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				{
					position658, tokenIndex658, depth658 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l659
					}
					goto l658
				l659:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
					if !_rules[ruleStream]() {
						goto l656
					}
				}
			l658:
				depth--
				add(ruleStreamLike, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 42 UDSFFuncApp <- <(FuncApp Action31)> */
		func() bool {
			position660, tokenIndex660, depth660 := position, tokenIndex, depth
			{
				position661 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l660
				}
				if !_rules[ruleAction31]() {
					goto l660
				}
				depth--
				add(ruleUDSFFuncApp, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 43 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action32)> */
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
							if buffer[position] != rune('i') {
								goto l670
							}
							position++
							goto l669
						l670:
							position, tokenIndex, depth = position669, tokenIndex669, depth669
							if buffer[position] != rune('I') {
								goto l665
							}
							position++
						}
					l669:
						{
							position671, tokenIndex671, depth671 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l672
							}
							position++
							goto l671
						l672:
							position, tokenIndex, depth = position671, tokenIndex671, depth671
							if buffer[position] != rune('T') {
								goto l665
							}
							position++
						}
					l671:
						{
							position673, tokenIndex673, depth673 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l674
							}
							position++
							goto l673
						l674:
							position, tokenIndex, depth = position673, tokenIndex673, depth673
							if buffer[position] != rune('H') {
								goto l665
							}
							position++
						}
					l673:
						if !_rules[rulesp]() {
							goto l665
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l665
						}
						if !_rules[rulesp]() {
							goto l665
						}
					l675:
						{
							position676, tokenIndex676, depth676 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l676
							}
							position++
							if !_rules[rulesp]() {
								goto l676
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l676
							}
							goto l675
						l676:
							position, tokenIndex, depth = position676, tokenIndex676, depth676
						}
						goto l666
					l665:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
					}
				l666:
					depth--
					add(rulePegText, position664)
				}
				if !_rules[ruleAction32]() {
					goto l662
				}
				depth--
				add(ruleSourceSinkSpecs, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 44 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action33)> */
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
						if buffer[position] != rune('s') {
							goto l681
						}
						position++
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if buffer[position] != rune('S') {
							goto l677
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
							goto l677
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
							goto l677
						}
						position++
					}
				l684:
					if !_rules[rulesp]() {
						goto l677
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l677
					}
					if !_rules[rulesp]() {
						goto l677
					}
				l686:
					{
						position687, tokenIndex687, depth687 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l687
						}
						position++
						if !_rules[rulesp]() {
							goto l687
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l687
						}
						goto l686
					l687:
						position, tokenIndex, depth = position687, tokenIndex687, depth687
					}
					depth--
					add(rulePegText, position679)
				}
				if !_rules[ruleAction33]() {
					goto l677
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 45 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action34)> */
		func() bool {
			position688, tokenIndex688, depth688 := position, tokenIndex, depth
			{
				position689 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l688
				}
				if buffer[position] != rune('=') {
					goto l688
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l688
				}
				if !_rules[ruleAction34]() {
					goto l688
				}
				depth--
				add(ruleSourceSinkParam, position689)
			}
			return true
		l688:
			position, tokenIndex, depth = position688, tokenIndex688, depth688
			return false
		},
		/* 46 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position690, tokenIndex690, depth690 := position, tokenIndex, depth
			{
				position691 := position
				depth++
				{
					position692, tokenIndex692, depth692 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l693
					}
					goto l692
				l693:
					position, tokenIndex, depth = position692, tokenIndex692, depth692
					if !_rules[ruleLiteral]() {
						goto l690
					}
				}
			l692:
				depth--
				add(ruleSourceSinkParamVal, position691)
			}
			return true
		l690:
			position, tokenIndex, depth = position690, tokenIndex690, depth690
			return false
		},
		/* 47 PausedOpt <- <(<(Paused / Unpaused)?> Action35)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				{
					position696 := position
					depth++
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						{
							position699, tokenIndex699, depth699 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l700
							}
							goto l699
						l700:
							position, tokenIndex, depth = position699, tokenIndex699, depth699
							if !_rules[ruleUnpaused]() {
								goto l697
							}
						}
					l699:
						goto l698
					l697:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
					}
				l698:
					depth--
					add(rulePegText, position696)
				}
				if !_rules[ruleAction35]() {
					goto l694
				}
				depth--
				add(rulePausedOpt, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 48 Expression <- <orExpr> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l701
				}
				depth--
				add(ruleExpression, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 49 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action36)> */
		func() bool {
			position703, tokenIndex703, depth703 := position, tokenIndex, depth
			{
				position704 := position
				depth++
				{
					position705 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l703
					}
					if !_rules[rulesp]() {
						goto l703
					}
					{
						position706, tokenIndex706, depth706 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l706
						}
						if !_rules[rulesp]() {
							goto l706
						}
						if !_rules[ruleandExpr]() {
							goto l706
						}
						goto l707
					l706:
						position, tokenIndex, depth = position706, tokenIndex706, depth706
					}
				l707:
					depth--
					add(rulePegText, position705)
				}
				if !_rules[ruleAction36]() {
					goto l703
				}
				depth--
				add(ruleorExpr, position704)
			}
			return true
		l703:
			position, tokenIndex, depth = position703, tokenIndex703, depth703
			return false
		},
		/* 50 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action37)> */
		func() bool {
			position708, tokenIndex708, depth708 := position, tokenIndex, depth
			{
				position709 := position
				depth++
				{
					position710 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l708
					}
					if !_rules[rulesp]() {
						goto l708
					}
					{
						position711, tokenIndex711, depth711 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l711
						}
						if !_rules[rulesp]() {
							goto l711
						}
						if !_rules[rulenotExpr]() {
							goto l711
						}
						goto l712
					l711:
						position, tokenIndex, depth = position711, tokenIndex711, depth711
					}
				l712:
					depth--
					add(rulePegText, position710)
				}
				if !_rules[ruleAction37]() {
					goto l708
				}
				depth--
				add(ruleandExpr, position709)
			}
			return true
		l708:
			position, tokenIndex, depth = position708, tokenIndex708, depth708
			return false
		},
		/* 51 notExpr <- <(<((Not sp)? comparisonExpr)> Action38)> */
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
						if !_rules[ruleNot]() {
							goto l716
						}
						if !_rules[rulesp]() {
							goto l716
						}
						goto l717
					l716:
						position, tokenIndex, depth = position716, tokenIndex716, depth716
					}
				l717:
					if !_rules[rulecomparisonExpr]() {
						goto l713
					}
					depth--
					add(rulePegText, position715)
				}
				if !_rules[ruleAction38]() {
					goto l713
				}
				depth--
				add(rulenotExpr, position714)
			}
			return true
		l713:
			position, tokenIndex, depth = position713, tokenIndex713, depth713
			return false
		},
		/* 52 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action39)> */
		func() bool {
			position718, tokenIndex718, depth718 := position, tokenIndex, depth
			{
				position719 := position
				depth++
				{
					position720 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l718
					}
					if !_rules[rulesp]() {
						goto l718
					}
					{
						position721, tokenIndex721, depth721 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l721
						}
						if !_rules[rulesp]() {
							goto l721
						}
						if !_rules[ruleotherOpExpr]() {
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
				if !_rules[ruleAction39]() {
					goto l718
				}
				depth--
				add(rulecomparisonExpr, position719)
			}
			return true
		l718:
			position, tokenIndex, depth = position718, tokenIndex718, depth718
			return false
		},
		/* 53 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action40)> */
		func() bool {
			position723, tokenIndex723, depth723 := position, tokenIndex, depth
			{
				position724 := position
				depth++
				{
					position725 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l723
					}
					if !_rules[rulesp]() {
						goto l723
					}
				l726:
					{
						position727, tokenIndex727, depth727 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l727
						}
						if !_rules[rulesp]() {
							goto l727
						}
						if !_rules[ruleisExpr]() {
							goto l727
						}
						if !_rules[rulesp]() {
							goto l727
						}
						goto l726
					l727:
						position, tokenIndex, depth = position727, tokenIndex727, depth727
					}
					depth--
					add(rulePegText, position725)
				}
				if !_rules[ruleAction40]() {
					goto l723
				}
				depth--
				add(ruleotherOpExpr, position724)
			}
			return true
		l723:
			position, tokenIndex, depth = position723, tokenIndex723, depth723
			return false
		},
		/* 54 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action41)> */
		func() bool {
			position728, tokenIndex728, depth728 := position, tokenIndex, depth
			{
				position729 := position
				depth++
				{
					position730 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l728
					}
					if !_rules[rulesp]() {
						goto l728
					}
					{
						position731, tokenIndex731, depth731 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l731
						}
						if !_rules[rulesp]() {
							goto l731
						}
						if !_rules[ruleNullLiteral]() {
							goto l731
						}
						goto l732
					l731:
						position, tokenIndex, depth = position731, tokenIndex731, depth731
					}
				l732:
					depth--
					add(rulePegText, position730)
				}
				if !_rules[ruleAction41]() {
					goto l728
				}
				depth--
				add(ruleisExpr, position729)
			}
			return true
		l728:
			position, tokenIndex, depth = position728, tokenIndex728, depth728
			return false
		},
		/* 55 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action42)> */
		func() bool {
			position733, tokenIndex733, depth733 := position, tokenIndex, depth
			{
				position734 := position
				depth++
				{
					position735 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l733
					}
					if !_rules[rulesp]() {
						goto l733
					}
				l736:
					{
						position737, tokenIndex737, depth737 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l737
						}
						if !_rules[rulesp]() {
							goto l737
						}
						if !_rules[ruleproductExpr]() {
							goto l737
						}
						if !_rules[rulesp]() {
							goto l737
						}
						goto l736
					l737:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
					}
					depth--
					add(rulePegText, position735)
				}
				if !_rules[ruleAction42]() {
					goto l733
				}
				depth--
				add(ruletermExpr, position734)
			}
			return true
		l733:
			position, tokenIndex, depth = position733, tokenIndex733, depth733
			return false
		},
		/* 56 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action43)> */
		func() bool {
			position738, tokenIndex738, depth738 := position, tokenIndex, depth
			{
				position739 := position
				depth++
				{
					position740 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l738
					}
					if !_rules[rulesp]() {
						goto l738
					}
				l741:
					{
						position742, tokenIndex742, depth742 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l742
						}
						if !_rules[rulesp]() {
							goto l742
						}
						if !_rules[ruleminusExpr]() {
							goto l742
						}
						if !_rules[rulesp]() {
							goto l742
						}
						goto l741
					l742:
						position, tokenIndex, depth = position742, tokenIndex742, depth742
					}
					depth--
					add(rulePegText, position740)
				}
				if !_rules[ruleAction43]() {
					goto l738
				}
				depth--
				add(ruleproductExpr, position739)
			}
			return true
		l738:
			position, tokenIndex, depth = position738, tokenIndex738, depth738
			return false
		},
		/* 57 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action44)> */
		func() bool {
			position743, tokenIndex743, depth743 := position, tokenIndex, depth
			{
				position744 := position
				depth++
				{
					position745 := position
					depth++
					{
						position746, tokenIndex746, depth746 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l746
						}
						if !_rules[rulesp]() {
							goto l746
						}
						goto l747
					l746:
						position, tokenIndex, depth = position746, tokenIndex746, depth746
					}
				l747:
					if !_rules[rulecastExpr]() {
						goto l743
					}
					depth--
					add(rulePegText, position745)
				}
				if !_rules[ruleAction44]() {
					goto l743
				}
				depth--
				add(ruleminusExpr, position744)
			}
			return true
		l743:
			position, tokenIndex, depth = position743, tokenIndex743, depth743
			return false
		},
		/* 58 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action45)> */
		func() bool {
			position748, tokenIndex748, depth748 := position, tokenIndex, depth
			{
				position749 := position
				depth++
				{
					position750 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l748
					}
					{
						position751, tokenIndex751, depth751 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l751
						}
						if buffer[position] != rune(':') {
							goto l751
						}
						position++
						if buffer[position] != rune(':') {
							goto l751
						}
						position++
						if !_rules[rulesp]() {
							goto l751
						}
						if !_rules[ruleType]() {
							goto l751
						}
						goto l752
					l751:
						position, tokenIndex, depth = position751, tokenIndex751, depth751
					}
				l752:
					depth--
					add(rulePegText, position750)
				}
				if !_rules[ruleAction45]() {
					goto l748
				}
				depth--
				add(rulecastExpr, position749)
			}
			return true
		l748:
			position, tokenIndex, depth = position748, tokenIndex748, depth748
			return false
		},
		/* 59 baseExpr <- <(('(' sp Expression sp ')') / ArrayExpr / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / Literal)> */
		func() bool {
			position753, tokenIndex753, depth753 := position, tokenIndex, depth
			{
				position754 := position
				depth++
				{
					position755, tokenIndex755, depth755 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l756
					}
					position++
					if !_rules[rulesp]() {
						goto l756
					}
					if !_rules[ruleExpression]() {
						goto l756
					}
					if !_rules[rulesp]() {
						goto l756
					}
					if buffer[position] != rune(')') {
						goto l756
					}
					position++
					goto l755
				l756:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleArrayExpr]() {
						goto l757
					}
					goto l755
				l757:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleBooleanLiteral]() {
						goto l758
					}
					goto l755
				l758:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleNullLiteral]() {
						goto l759
					}
					goto l755
				l759:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleRowMeta]() {
						goto l760
					}
					goto l755
				l760:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleFuncTypeCast]() {
						goto l761
					}
					goto l755
				l761:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleFuncApp]() {
						goto l762
					}
					goto l755
				l762:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleRowValue]() {
						goto l763
					}
					goto l755
				l763:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if !_rules[ruleLiteral]() {
						goto l753
					}
				}
			l755:
				depth--
				add(rulebaseExpr, position754)
			}
			return true
		l753:
			position, tokenIndex, depth = position753, tokenIndex753, depth753
			return false
		},
		/* 60 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action46)> */
		func() bool {
			position764, tokenIndex764, depth764 := position, tokenIndex, depth
			{
				position765 := position
				depth++
				{
					position766 := position
					depth++
					{
						position767, tokenIndex767, depth767 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l768
						}
						position++
						goto l767
					l768:
						position, tokenIndex, depth = position767, tokenIndex767, depth767
						if buffer[position] != rune('C') {
							goto l764
						}
						position++
					}
				l767:
					{
						position769, tokenIndex769, depth769 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l770
						}
						position++
						goto l769
					l770:
						position, tokenIndex, depth = position769, tokenIndex769, depth769
						if buffer[position] != rune('A') {
							goto l764
						}
						position++
					}
				l769:
					{
						position771, tokenIndex771, depth771 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l772
						}
						position++
						goto l771
					l772:
						position, tokenIndex, depth = position771, tokenIndex771, depth771
						if buffer[position] != rune('S') {
							goto l764
						}
						position++
					}
				l771:
					{
						position773, tokenIndex773, depth773 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l774
						}
						position++
						goto l773
					l774:
						position, tokenIndex, depth = position773, tokenIndex773, depth773
						if buffer[position] != rune('T') {
							goto l764
						}
						position++
					}
				l773:
					if !_rules[rulesp]() {
						goto l764
					}
					if buffer[position] != rune('(') {
						goto l764
					}
					position++
					if !_rules[rulesp]() {
						goto l764
					}
					if !_rules[ruleExpression]() {
						goto l764
					}
					if !_rules[rulesp]() {
						goto l764
					}
					{
						position775, tokenIndex775, depth775 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l776
						}
						position++
						goto l775
					l776:
						position, tokenIndex, depth = position775, tokenIndex775, depth775
						if buffer[position] != rune('A') {
							goto l764
						}
						position++
					}
				l775:
					{
						position777, tokenIndex777, depth777 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l778
						}
						position++
						goto l777
					l778:
						position, tokenIndex, depth = position777, tokenIndex777, depth777
						if buffer[position] != rune('S') {
							goto l764
						}
						position++
					}
				l777:
					if !_rules[rulesp]() {
						goto l764
					}
					if !_rules[ruleType]() {
						goto l764
					}
					if !_rules[rulesp]() {
						goto l764
					}
					if buffer[position] != rune(')') {
						goto l764
					}
					position++
					depth--
					add(rulePegText, position766)
				}
				if !_rules[ruleAction46]() {
					goto l764
				}
				depth--
				add(ruleFuncTypeCast, position765)
			}
			return true
		l764:
			position, tokenIndex, depth = position764, tokenIndex764, depth764
			return false
		},
		/* 61 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action47)> */
		func() bool {
			position779, tokenIndex779, depth779 := position, tokenIndex, depth
			{
				position780 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l779
				}
				if !_rules[rulesp]() {
					goto l779
				}
				if buffer[position] != rune('(') {
					goto l779
				}
				position++
				if !_rules[rulesp]() {
					goto l779
				}
				if !_rules[ruleFuncParams]() {
					goto l779
				}
				if !_rules[rulesp]() {
					goto l779
				}
				if buffer[position] != rune(')') {
					goto l779
				}
				position++
				if !_rules[ruleAction47]() {
					goto l779
				}
				depth--
				add(ruleFuncApp, position780)
			}
			return true
		l779:
			position, tokenIndex, depth = position779, tokenIndex779, depth779
			return false
		},
		/* 62 FuncParams <- <(<(Star / (Expression sp (',' sp Expression)*)?)> Action48)> */
		func() bool {
			position781, tokenIndex781, depth781 := position, tokenIndex, depth
			{
				position782 := position
				depth++
				{
					position783 := position
					depth++
					{
						position784, tokenIndex784, depth784 := position, tokenIndex, depth
						if !_rules[ruleStar]() {
							goto l785
						}
						goto l784
					l785:
						position, tokenIndex, depth = position784, tokenIndex784, depth784
						{
							position786, tokenIndex786, depth786 := position, tokenIndex, depth
							if !_rules[ruleExpression]() {
								goto l786
							}
							if !_rules[rulesp]() {
								goto l786
							}
						l788:
							{
								position789, tokenIndex789, depth789 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l789
								}
								position++
								if !_rules[rulesp]() {
									goto l789
								}
								if !_rules[ruleExpression]() {
									goto l789
								}
								goto l788
							l789:
								position, tokenIndex, depth = position789, tokenIndex789, depth789
							}
							goto l787
						l786:
							position, tokenIndex, depth = position786, tokenIndex786, depth786
						}
					l787:
					}
				l784:
					depth--
					add(rulePegText, position783)
				}
				if !_rules[ruleAction48]() {
					goto l781
				}
				depth--
				add(ruleFuncParams, position782)
			}
			return true
		l781:
			position, tokenIndex, depth = position781, tokenIndex781, depth781
			return false
		},
		/* 63 ArrayExpr <- <(<('[' sp (Expression (',' sp Expression)*)? sp ']')> Action49)> */
		func() bool {
			position790, tokenIndex790, depth790 := position, tokenIndex, depth
			{
				position791 := position
				depth++
				{
					position792 := position
					depth++
					if buffer[position] != rune('[') {
						goto l790
					}
					position++
					if !_rules[rulesp]() {
						goto l790
					}
					{
						position793, tokenIndex793, depth793 := position, tokenIndex, depth
						if !_rules[ruleExpression]() {
							goto l793
						}
					l795:
						{
							position796, tokenIndex796, depth796 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l796
							}
							position++
							if !_rules[rulesp]() {
								goto l796
							}
							if !_rules[ruleExpression]() {
								goto l796
							}
							goto l795
						l796:
							position, tokenIndex, depth = position796, tokenIndex796, depth796
						}
						goto l794
					l793:
						position, tokenIndex, depth = position793, tokenIndex793, depth793
					}
				l794:
					if !_rules[rulesp]() {
						goto l790
					}
					if buffer[position] != rune(']') {
						goto l790
					}
					position++
					depth--
					add(rulePegText, position792)
				}
				if !_rules[ruleAction49]() {
					goto l790
				}
				depth--
				add(ruleArrayExpr, position791)
			}
			return true
		l790:
			position, tokenIndex, depth = position790, tokenIndex790, depth790
			return false
		},
		/* 64 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position797, tokenIndex797, depth797 := position, tokenIndex, depth
			{
				position798 := position
				depth++
				{
					position799, tokenIndex799, depth799 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l800
					}
					goto l799
				l800:
					position, tokenIndex, depth = position799, tokenIndex799, depth799
					if !_rules[ruleNumericLiteral]() {
						goto l801
					}
					goto l799
				l801:
					position, tokenIndex, depth = position799, tokenIndex799, depth799
					if !_rules[ruleStringLiteral]() {
						goto l797
					}
				}
			l799:
				depth--
				add(ruleLiteral, position798)
			}
			return true
		l797:
			position, tokenIndex, depth = position797, tokenIndex797, depth797
			return false
		},
		/* 65 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position802, tokenIndex802, depth802 := position, tokenIndex, depth
			{
				position803 := position
				depth++
				{
					position804, tokenIndex804, depth804 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l805
					}
					goto l804
				l805:
					position, tokenIndex, depth = position804, tokenIndex804, depth804
					if !_rules[ruleNotEqual]() {
						goto l806
					}
					goto l804
				l806:
					position, tokenIndex, depth = position804, tokenIndex804, depth804
					if !_rules[ruleLessOrEqual]() {
						goto l807
					}
					goto l804
				l807:
					position, tokenIndex, depth = position804, tokenIndex804, depth804
					if !_rules[ruleLess]() {
						goto l808
					}
					goto l804
				l808:
					position, tokenIndex, depth = position804, tokenIndex804, depth804
					if !_rules[ruleGreaterOrEqual]() {
						goto l809
					}
					goto l804
				l809:
					position, tokenIndex, depth = position804, tokenIndex804, depth804
					if !_rules[ruleGreater]() {
						goto l810
					}
					goto l804
				l810:
					position, tokenIndex, depth = position804, tokenIndex804, depth804
					if !_rules[ruleNotEqual]() {
						goto l802
					}
				}
			l804:
				depth--
				add(ruleComparisonOp, position803)
			}
			return true
		l802:
			position, tokenIndex, depth = position802, tokenIndex802, depth802
			return false
		},
		/* 66 OtherOp <- <Concat> */
		func() bool {
			position811, tokenIndex811, depth811 := position, tokenIndex, depth
			{
				position812 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l811
				}
				depth--
				add(ruleOtherOp, position812)
			}
			return true
		l811:
			position, tokenIndex, depth = position811, tokenIndex811, depth811
			return false
		},
		/* 67 IsOp <- <(IsNot / Is)> */
		func() bool {
			position813, tokenIndex813, depth813 := position, tokenIndex, depth
			{
				position814 := position
				depth++
				{
					position815, tokenIndex815, depth815 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l816
					}
					goto l815
				l816:
					position, tokenIndex, depth = position815, tokenIndex815, depth815
					if !_rules[ruleIs]() {
						goto l813
					}
				}
			l815:
				depth--
				add(ruleIsOp, position814)
			}
			return true
		l813:
			position, tokenIndex, depth = position813, tokenIndex813, depth813
			return false
		},
		/* 68 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position817, tokenIndex817, depth817 := position, tokenIndex, depth
			{
				position818 := position
				depth++
				{
					position819, tokenIndex819, depth819 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l820
					}
					goto l819
				l820:
					position, tokenIndex, depth = position819, tokenIndex819, depth819
					if !_rules[ruleMinus]() {
						goto l817
					}
				}
			l819:
				depth--
				add(rulePlusMinusOp, position818)
			}
			return true
		l817:
			position, tokenIndex, depth = position817, tokenIndex817, depth817
			return false
		},
		/* 69 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position821, tokenIndex821, depth821 := position, tokenIndex, depth
			{
				position822 := position
				depth++
				{
					position823, tokenIndex823, depth823 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l824
					}
					goto l823
				l824:
					position, tokenIndex, depth = position823, tokenIndex823, depth823
					if !_rules[ruleDivide]() {
						goto l825
					}
					goto l823
				l825:
					position, tokenIndex, depth = position823, tokenIndex823, depth823
					if !_rules[ruleModulo]() {
						goto l821
					}
				}
			l823:
				depth--
				add(ruleMultDivOp, position822)
			}
			return true
		l821:
			position, tokenIndex, depth = position821, tokenIndex821, depth821
			return false
		},
		/* 70 Stream <- <(<ident> Action50)> */
		func() bool {
			position826, tokenIndex826, depth826 := position, tokenIndex, depth
			{
				position827 := position
				depth++
				{
					position828 := position
					depth++
					if !_rules[ruleident]() {
						goto l826
					}
					depth--
					add(rulePegText, position828)
				}
				if !_rules[ruleAction50]() {
					goto l826
				}
				depth--
				add(ruleStream, position827)
			}
			return true
		l826:
			position, tokenIndex, depth = position826, tokenIndex826, depth826
			return false
		},
		/* 71 RowMeta <- <RowTimestamp> */
		func() bool {
			position829, tokenIndex829, depth829 := position, tokenIndex, depth
			{
				position830 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l829
				}
				depth--
				add(ruleRowMeta, position830)
			}
			return true
		l829:
			position, tokenIndex, depth = position829, tokenIndex829, depth829
			return false
		},
		/* 72 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action51)> */
		func() bool {
			position831, tokenIndex831, depth831 := position, tokenIndex, depth
			{
				position832 := position
				depth++
				{
					position833 := position
					depth++
					{
						position834, tokenIndex834, depth834 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l834
						}
						if buffer[position] != rune(':') {
							goto l834
						}
						position++
						goto l835
					l834:
						position, tokenIndex, depth = position834, tokenIndex834, depth834
					}
				l835:
					if buffer[position] != rune('t') {
						goto l831
					}
					position++
					if buffer[position] != rune('s') {
						goto l831
					}
					position++
					if buffer[position] != rune('(') {
						goto l831
					}
					position++
					if buffer[position] != rune(')') {
						goto l831
					}
					position++
					depth--
					add(rulePegText, position833)
				}
				if !_rules[ruleAction51]() {
					goto l831
				}
				depth--
				add(ruleRowTimestamp, position832)
			}
			return true
		l831:
			position, tokenIndex, depth = position831, tokenIndex831, depth831
			return false
		},
		/* 73 RowValue <- <(<((ident ':' !':')? jsonPath)> Action52)> */
		func() bool {
			position836, tokenIndex836, depth836 := position, tokenIndex, depth
			{
				position837 := position
				depth++
				{
					position838 := position
					depth++
					{
						position839, tokenIndex839, depth839 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l839
						}
						if buffer[position] != rune(':') {
							goto l839
						}
						position++
						{
							position841, tokenIndex841, depth841 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l841
							}
							position++
							goto l839
						l841:
							position, tokenIndex, depth = position841, tokenIndex841, depth841
						}
						goto l840
					l839:
						position, tokenIndex, depth = position839, tokenIndex839, depth839
					}
				l840:
					if !_rules[rulejsonPath]() {
						goto l836
					}
					depth--
					add(rulePegText, position838)
				}
				if !_rules[ruleAction52]() {
					goto l836
				}
				depth--
				add(ruleRowValue, position837)
			}
			return true
		l836:
			position, tokenIndex, depth = position836, tokenIndex836, depth836
			return false
		},
		/* 74 NumericLiteral <- <(<('-'? [0-9]+)> Action53)> */
		func() bool {
			position842, tokenIndex842, depth842 := position, tokenIndex, depth
			{
				position843 := position
				depth++
				{
					position844 := position
					depth++
					{
						position845, tokenIndex845, depth845 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l845
						}
						position++
						goto l846
					l845:
						position, tokenIndex, depth = position845, tokenIndex845, depth845
					}
				l846:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l842
					}
					position++
				l847:
					{
						position848, tokenIndex848, depth848 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l848
						}
						position++
						goto l847
					l848:
						position, tokenIndex, depth = position848, tokenIndex848, depth848
					}
					depth--
					add(rulePegText, position844)
				}
				if !_rules[ruleAction53]() {
					goto l842
				}
				depth--
				add(ruleNumericLiteral, position843)
			}
			return true
		l842:
			position, tokenIndex, depth = position842, tokenIndex842, depth842
			return false
		},
		/* 75 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action54)> */
		func() bool {
			position849, tokenIndex849, depth849 := position, tokenIndex, depth
			{
				position850 := position
				depth++
				{
					position851 := position
					depth++
					{
						position852, tokenIndex852, depth852 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l852
						}
						position++
						goto l853
					l852:
						position, tokenIndex, depth = position852, tokenIndex852, depth852
					}
				l853:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l849
					}
					position++
				l854:
					{
						position855, tokenIndex855, depth855 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l855
						}
						position++
						goto l854
					l855:
						position, tokenIndex, depth = position855, tokenIndex855, depth855
					}
					if buffer[position] != rune('.') {
						goto l849
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l849
					}
					position++
				l856:
					{
						position857, tokenIndex857, depth857 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l857
						}
						position++
						goto l856
					l857:
						position, tokenIndex, depth = position857, tokenIndex857, depth857
					}
					depth--
					add(rulePegText, position851)
				}
				if !_rules[ruleAction54]() {
					goto l849
				}
				depth--
				add(ruleFloatLiteral, position850)
			}
			return true
		l849:
			position, tokenIndex, depth = position849, tokenIndex849, depth849
			return false
		},
		/* 76 Function <- <(<ident> Action55)> */
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
				add(ruleFunction, position859)
			}
			return true
		l858:
			position, tokenIndex, depth = position858, tokenIndex858, depth858
			return false
		},
		/* 77 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action56)> */
		func() bool {
			position861, tokenIndex861, depth861 := position, tokenIndex, depth
			{
				position862 := position
				depth++
				{
					position863 := position
					depth++
					{
						position864, tokenIndex864, depth864 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l865
						}
						position++
						goto l864
					l865:
						position, tokenIndex, depth = position864, tokenIndex864, depth864
						if buffer[position] != rune('N') {
							goto l861
						}
						position++
					}
				l864:
					{
						position866, tokenIndex866, depth866 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l867
						}
						position++
						goto l866
					l867:
						position, tokenIndex, depth = position866, tokenIndex866, depth866
						if buffer[position] != rune('U') {
							goto l861
						}
						position++
					}
				l866:
					{
						position868, tokenIndex868, depth868 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l869
						}
						position++
						goto l868
					l869:
						position, tokenIndex, depth = position868, tokenIndex868, depth868
						if buffer[position] != rune('L') {
							goto l861
						}
						position++
					}
				l868:
					{
						position870, tokenIndex870, depth870 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l871
						}
						position++
						goto l870
					l871:
						position, tokenIndex, depth = position870, tokenIndex870, depth870
						if buffer[position] != rune('L') {
							goto l861
						}
						position++
					}
				l870:
					depth--
					add(rulePegText, position863)
				}
				if !_rules[ruleAction56]() {
					goto l861
				}
				depth--
				add(ruleNullLiteral, position862)
			}
			return true
		l861:
			position, tokenIndex, depth = position861, tokenIndex861, depth861
			return false
		},
		/* 78 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position872, tokenIndex872, depth872 := position, tokenIndex, depth
			{
				position873 := position
				depth++
				{
					position874, tokenIndex874, depth874 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l875
					}
					goto l874
				l875:
					position, tokenIndex, depth = position874, tokenIndex874, depth874
					if !_rules[ruleFALSE]() {
						goto l872
					}
				}
			l874:
				depth--
				add(ruleBooleanLiteral, position873)
			}
			return true
		l872:
			position, tokenIndex, depth = position872, tokenIndex872, depth872
			return false
		},
		/* 79 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action57)> */
		func() bool {
			position876, tokenIndex876, depth876 := position, tokenIndex, depth
			{
				position877 := position
				depth++
				{
					position878 := position
					depth++
					{
						position879, tokenIndex879, depth879 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l880
						}
						position++
						goto l879
					l880:
						position, tokenIndex, depth = position879, tokenIndex879, depth879
						if buffer[position] != rune('T') {
							goto l876
						}
						position++
					}
				l879:
					{
						position881, tokenIndex881, depth881 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l882
						}
						position++
						goto l881
					l882:
						position, tokenIndex, depth = position881, tokenIndex881, depth881
						if buffer[position] != rune('R') {
							goto l876
						}
						position++
					}
				l881:
					{
						position883, tokenIndex883, depth883 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l884
						}
						position++
						goto l883
					l884:
						position, tokenIndex, depth = position883, tokenIndex883, depth883
						if buffer[position] != rune('U') {
							goto l876
						}
						position++
					}
				l883:
					{
						position885, tokenIndex885, depth885 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l886
						}
						position++
						goto l885
					l886:
						position, tokenIndex, depth = position885, tokenIndex885, depth885
						if buffer[position] != rune('E') {
							goto l876
						}
						position++
					}
				l885:
					depth--
					add(rulePegText, position878)
				}
				if !_rules[ruleAction57]() {
					goto l876
				}
				depth--
				add(ruleTRUE, position877)
			}
			return true
		l876:
			position, tokenIndex, depth = position876, tokenIndex876, depth876
			return false
		},
		/* 80 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action58)> */
		func() bool {
			position887, tokenIndex887, depth887 := position, tokenIndex, depth
			{
				position888 := position
				depth++
				{
					position889 := position
					depth++
					{
						position890, tokenIndex890, depth890 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l891
						}
						position++
						goto l890
					l891:
						position, tokenIndex, depth = position890, tokenIndex890, depth890
						if buffer[position] != rune('F') {
							goto l887
						}
						position++
					}
				l890:
					{
						position892, tokenIndex892, depth892 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l893
						}
						position++
						goto l892
					l893:
						position, tokenIndex, depth = position892, tokenIndex892, depth892
						if buffer[position] != rune('A') {
							goto l887
						}
						position++
					}
				l892:
					{
						position894, tokenIndex894, depth894 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l895
						}
						position++
						goto l894
					l895:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
						if buffer[position] != rune('L') {
							goto l887
						}
						position++
					}
				l894:
					{
						position896, tokenIndex896, depth896 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l897
						}
						position++
						goto l896
					l897:
						position, tokenIndex, depth = position896, tokenIndex896, depth896
						if buffer[position] != rune('S') {
							goto l887
						}
						position++
					}
				l896:
					{
						position898, tokenIndex898, depth898 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l899
						}
						position++
						goto l898
					l899:
						position, tokenIndex, depth = position898, tokenIndex898, depth898
						if buffer[position] != rune('E') {
							goto l887
						}
						position++
					}
				l898:
					depth--
					add(rulePegText, position889)
				}
				if !_rules[ruleAction58]() {
					goto l887
				}
				depth--
				add(ruleFALSE, position888)
			}
			return true
		l887:
			position, tokenIndex, depth = position887, tokenIndex887, depth887
			return false
		},
		/* 81 Star <- <(<'*'> Action59)> */
		func() bool {
			position900, tokenIndex900, depth900 := position, tokenIndex, depth
			{
				position901 := position
				depth++
				{
					position902 := position
					depth++
					if buffer[position] != rune('*') {
						goto l900
					}
					position++
					depth--
					add(rulePegText, position902)
				}
				if !_rules[ruleAction59]() {
					goto l900
				}
				depth--
				add(ruleStar, position901)
			}
			return true
		l900:
			position, tokenIndex, depth = position900, tokenIndex900, depth900
			return false
		},
		/* 82 Wildcard <- <(<((ident ':' !':')? '*')> Action60)> */
		func() bool {
			position903, tokenIndex903, depth903 := position, tokenIndex, depth
			{
				position904 := position
				depth++
				{
					position905 := position
					depth++
					{
						position906, tokenIndex906, depth906 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l906
						}
						if buffer[position] != rune(':') {
							goto l906
						}
						position++
						{
							position908, tokenIndex908, depth908 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l908
							}
							position++
							goto l906
						l908:
							position, tokenIndex, depth = position908, tokenIndex908, depth908
						}
						goto l907
					l906:
						position, tokenIndex, depth = position906, tokenIndex906, depth906
					}
				l907:
					if buffer[position] != rune('*') {
						goto l903
					}
					position++
					depth--
					add(rulePegText, position905)
				}
				if !_rules[ruleAction60]() {
					goto l903
				}
				depth--
				add(ruleWildcard, position904)
			}
			return true
		l903:
			position, tokenIndex, depth = position903, tokenIndex903, depth903
			return false
		},
		/* 83 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action61)> */
		func() bool {
			position909, tokenIndex909, depth909 := position, tokenIndex, depth
			{
				position910 := position
				depth++
				{
					position911 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l909
					}
					position++
				l912:
					{
						position913, tokenIndex913, depth913 := position, tokenIndex, depth
						{
							position914, tokenIndex914, depth914 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l915
							}
							position++
							if buffer[position] != rune('\'') {
								goto l915
							}
							position++
							goto l914
						l915:
							position, tokenIndex, depth = position914, tokenIndex914, depth914
							{
								position916, tokenIndex916, depth916 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l916
								}
								position++
								goto l913
							l916:
								position, tokenIndex, depth = position916, tokenIndex916, depth916
							}
							if !matchDot() {
								goto l913
							}
						}
					l914:
						goto l912
					l913:
						position, tokenIndex, depth = position913, tokenIndex913, depth913
					}
					if buffer[position] != rune('\'') {
						goto l909
					}
					position++
					depth--
					add(rulePegText, position911)
				}
				if !_rules[ruleAction61]() {
					goto l909
				}
				depth--
				add(ruleStringLiteral, position910)
			}
			return true
		l909:
			position, tokenIndex, depth = position909, tokenIndex909, depth909
			return false
		},
		/* 84 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action62)> */
		func() bool {
			position917, tokenIndex917, depth917 := position, tokenIndex, depth
			{
				position918 := position
				depth++
				{
					position919 := position
					depth++
					{
						position920, tokenIndex920, depth920 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l921
						}
						position++
						goto l920
					l921:
						position, tokenIndex, depth = position920, tokenIndex920, depth920
						if buffer[position] != rune('I') {
							goto l917
						}
						position++
					}
				l920:
					{
						position922, tokenIndex922, depth922 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l923
						}
						position++
						goto l922
					l923:
						position, tokenIndex, depth = position922, tokenIndex922, depth922
						if buffer[position] != rune('S') {
							goto l917
						}
						position++
					}
				l922:
					{
						position924, tokenIndex924, depth924 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l925
						}
						position++
						goto l924
					l925:
						position, tokenIndex, depth = position924, tokenIndex924, depth924
						if buffer[position] != rune('T') {
							goto l917
						}
						position++
					}
				l924:
					{
						position926, tokenIndex926, depth926 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('R') {
							goto l917
						}
						position++
					}
				l926:
					{
						position928, tokenIndex928, depth928 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l929
						}
						position++
						goto l928
					l929:
						position, tokenIndex, depth = position928, tokenIndex928, depth928
						if buffer[position] != rune('E') {
							goto l917
						}
						position++
					}
				l928:
					{
						position930, tokenIndex930, depth930 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l931
						}
						position++
						goto l930
					l931:
						position, tokenIndex, depth = position930, tokenIndex930, depth930
						if buffer[position] != rune('A') {
							goto l917
						}
						position++
					}
				l930:
					{
						position932, tokenIndex932, depth932 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l933
						}
						position++
						goto l932
					l933:
						position, tokenIndex, depth = position932, tokenIndex932, depth932
						if buffer[position] != rune('M') {
							goto l917
						}
						position++
					}
				l932:
					depth--
					add(rulePegText, position919)
				}
				if !_rules[ruleAction62]() {
					goto l917
				}
				depth--
				add(ruleISTREAM, position918)
			}
			return true
		l917:
			position, tokenIndex, depth = position917, tokenIndex917, depth917
			return false
		},
		/* 85 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action63)> */
		func() bool {
			position934, tokenIndex934, depth934 := position, tokenIndex, depth
			{
				position935 := position
				depth++
				{
					position936 := position
					depth++
					{
						position937, tokenIndex937, depth937 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l938
						}
						position++
						goto l937
					l938:
						position, tokenIndex, depth = position937, tokenIndex937, depth937
						if buffer[position] != rune('D') {
							goto l934
						}
						position++
					}
				l937:
					{
						position939, tokenIndex939, depth939 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l940
						}
						position++
						goto l939
					l940:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						if buffer[position] != rune('S') {
							goto l934
						}
						position++
					}
				l939:
					{
						position941, tokenIndex941, depth941 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l942
						}
						position++
						goto l941
					l942:
						position, tokenIndex, depth = position941, tokenIndex941, depth941
						if buffer[position] != rune('T') {
							goto l934
						}
						position++
					}
				l941:
					{
						position943, tokenIndex943, depth943 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l944
						}
						position++
						goto l943
					l944:
						position, tokenIndex, depth = position943, tokenIndex943, depth943
						if buffer[position] != rune('R') {
							goto l934
						}
						position++
					}
				l943:
					{
						position945, tokenIndex945, depth945 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l946
						}
						position++
						goto l945
					l946:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
						if buffer[position] != rune('E') {
							goto l934
						}
						position++
					}
				l945:
					{
						position947, tokenIndex947, depth947 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l948
						}
						position++
						goto l947
					l948:
						position, tokenIndex, depth = position947, tokenIndex947, depth947
						if buffer[position] != rune('A') {
							goto l934
						}
						position++
					}
				l947:
					{
						position949, tokenIndex949, depth949 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l950
						}
						position++
						goto l949
					l950:
						position, tokenIndex, depth = position949, tokenIndex949, depth949
						if buffer[position] != rune('M') {
							goto l934
						}
						position++
					}
				l949:
					depth--
					add(rulePegText, position936)
				}
				if !_rules[ruleAction63]() {
					goto l934
				}
				depth--
				add(ruleDSTREAM, position935)
			}
			return true
		l934:
			position, tokenIndex, depth = position934, tokenIndex934, depth934
			return false
		},
		/* 86 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action64)> */
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
						if buffer[position] != rune('r') {
							goto l955
						}
						position++
						goto l954
					l955:
						position, tokenIndex, depth = position954, tokenIndex954, depth954
						if buffer[position] != rune('R') {
							goto l951
						}
						position++
					}
				l954:
					{
						position956, tokenIndex956, depth956 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l957
						}
						position++
						goto l956
					l957:
						position, tokenIndex, depth = position956, tokenIndex956, depth956
						if buffer[position] != rune('S') {
							goto l951
						}
						position++
					}
				l956:
					{
						position958, tokenIndex958, depth958 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l959
						}
						position++
						goto l958
					l959:
						position, tokenIndex, depth = position958, tokenIndex958, depth958
						if buffer[position] != rune('T') {
							goto l951
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
							goto l951
						}
						position++
					}
				l960:
					{
						position962, tokenIndex962, depth962 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l963
						}
						position++
						goto l962
					l963:
						position, tokenIndex, depth = position962, tokenIndex962, depth962
						if buffer[position] != rune('E') {
							goto l951
						}
						position++
					}
				l962:
					{
						position964, tokenIndex964, depth964 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l965
						}
						position++
						goto l964
					l965:
						position, tokenIndex, depth = position964, tokenIndex964, depth964
						if buffer[position] != rune('A') {
							goto l951
						}
						position++
					}
				l964:
					{
						position966, tokenIndex966, depth966 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l967
						}
						position++
						goto l966
					l967:
						position, tokenIndex, depth = position966, tokenIndex966, depth966
						if buffer[position] != rune('M') {
							goto l951
						}
						position++
					}
				l966:
					depth--
					add(rulePegText, position953)
				}
				if !_rules[ruleAction64]() {
					goto l951
				}
				depth--
				add(ruleRSTREAM, position952)
			}
			return true
		l951:
			position, tokenIndex, depth = position951, tokenIndex951, depth951
			return false
		},
		/* 87 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action65)> */
		func() bool {
			position968, tokenIndex968, depth968 := position, tokenIndex, depth
			{
				position969 := position
				depth++
				{
					position970 := position
					depth++
					{
						position971, tokenIndex971, depth971 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l972
						}
						position++
						goto l971
					l972:
						position, tokenIndex, depth = position971, tokenIndex971, depth971
						if buffer[position] != rune('T') {
							goto l968
						}
						position++
					}
				l971:
					{
						position973, tokenIndex973, depth973 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l974
						}
						position++
						goto l973
					l974:
						position, tokenIndex, depth = position973, tokenIndex973, depth973
						if buffer[position] != rune('U') {
							goto l968
						}
						position++
					}
				l973:
					{
						position975, tokenIndex975, depth975 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l976
						}
						position++
						goto l975
					l976:
						position, tokenIndex, depth = position975, tokenIndex975, depth975
						if buffer[position] != rune('P') {
							goto l968
						}
						position++
					}
				l975:
					{
						position977, tokenIndex977, depth977 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l978
						}
						position++
						goto l977
					l978:
						position, tokenIndex, depth = position977, tokenIndex977, depth977
						if buffer[position] != rune('L') {
							goto l968
						}
						position++
					}
				l977:
					{
						position979, tokenIndex979, depth979 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l980
						}
						position++
						goto l979
					l980:
						position, tokenIndex, depth = position979, tokenIndex979, depth979
						if buffer[position] != rune('E') {
							goto l968
						}
						position++
					}
				l979:
					{
						position981, tokenIndex981, depth981 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l982
						}
						position++
						goto l981
					l982:
						position, tokenIndex, depth = position981, tokenIndex981, depth981
						if buffer[position] != rune('S') {
							goto l968
						}
						position++
					}
				l981:
					depth--
					add(rulePegText, position970)
				}
				if !_rules[ruleAction65]() {
					goto l968
				}
				depth--
				add(ruleTUPLES, position969)
			}
			return true
		l968:
			position, tokenIndex, depth = position968, tokenIndex968, depth968
			return false
		},
		/* 88 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action66)> */
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
						if buffer[position] != rune('s') {
							goto l987
						}
						position++
						goto l986
					l987:
						position, tokenIndex, depth = position986, tokenIndex986, depth986
						if buffer[position] != rune('S') {
							goto l983
						}
						position++
					}
				l986:
					{
						position988, tokenIndex988, depth988 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l989
						}
						position++
						goto l988
					l989:
						position, tokenIndex, depth = position988, tokenIndex988, depth988
						if buffer[position] != rune('E') {
							goto l983
						}
						position++
					}
				l988:
					{
						position990, tokenIndex990, depth990 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l991
						}
						position++
						goto l990
					l991:
						position, tokenIndex, depth = position990, tokenIndex990, depth990
						if buffer[position] != rune('C') {
							goto l983
						}
						position++
					}
				l990:
					{
						position992, tokenIndex992, depth992 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l993
						}
						position++
						goto l992
					l993:
						position, tokenIndex, depth = position992, tokenIndex992, depth992
						if buffer[position] != rune('O') {
							goto l983
						}
						position++
					}
				l992:
					{
						position994, tokenIndex994, depth994 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l995
						}
						position++
						goto l994
					l995:
						position, tokenIndex, depth = position994, tokenIndex994, depth994
						if buffer[position] != rune('N') {
							goto l983
						}
						position++
					}
				l994:
					{
						position996, tokenIndex996, depth996 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l997
						}
						position++
						goto l996
					l997:
						position, tokenIndex, depth = position996, tokenIndex996, depth996
						if buffer[position] != rune('D') {
							goto l983
						}
						position++
					}
				l996:
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
							goto l983
						}
						position++
					}
				l998:
					depth--
					add(rulePegText, position985)
				}
				if !_rules[ruleAction66]() {
					goto l983
				}
				depth--
				add(ruleSECONDS, position984)
			}
			return true
		l983:
			position, tokenIndex, depth = position983, tokenIndex983, depth983
			return false
		},
		/* 89 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action67)> */
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
						if buffer[position] != rune('m') {
							goto l1004
						}
						position++
						goto l1003
					l1004:
						position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
						if buffer[position] != rune('M') {
							goto l1000
						}
						position++
					}
				l1003:
					{
						position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1006
						}
						position++
						goto l1005
					l1006:
						position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
						if buffer[position] != rune('I') {
							goto l1000
						}
						position++
					}
				l1005:
					{
						position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1008
						}
						position++
						goto l1007
					l1008:
						position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
						if buffer[position] != rune('L') {
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
						if buffer[position] != rune('i') {
							goto l1012
						}
						position++
						goto l1011
					l1012:
						position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
						if buffer[position] != rune('I') {
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
					{
						position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1016
						}
						position++
						goto l1015
					l1016:
						position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
						if buffer[position] != rune('E') {
							goto l1000
						}
						position++
					}
				l1015:
					{
						position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1018
						}
						position++
						goto l1017
					l1018:
						position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
						if buffer[position] != rune('C') {
							goto l1000
						}
						position++
					}
				l1017:
					{
						position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1020
						}
						position++
						goto l1019
					l1020:
						position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
						if buffer[position] != rune('O') {
							goto l1000
						}
						position++
					}
				l1019:
					{
						position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1022
						}
						position++
						goto l1021
					l1022:
						position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
						if buffer[position] != rune('N') {
							goto l1000
						}
						position++
					}
				l1021:
					{
						position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1024
						}
						position++
						goto l1023
					l1024:
						position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
						if buffer[position] != rune('D') {
							goto l1000
						}
						position++
					}
				l1023:
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
							goto l1000
						}
						position++
					}
				l1025:
					depth--
					add(rulePegText, position1002)
				}
				if !_rules[ruleAction67]() {
					goto l1000
				}
				depth--
				add(ruleMILLISECONDS, position1001)
			}
			return true
		l1000:
			position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
			return false
		},
		/* 90 StreamIdentifier <- <(<ident> Action68)> */
		func() bool {
			position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
			{
				position1028 := position
				depth++
				{
					position1029 := position
					depth++
					if !_rules[ruleident]() {
						goto l1027
					}
					depth--
					add(rulePegText, position1029)
				}
				if !_rules[ruleAction68]() {
					goto l1027
				}
				depth--
				add(ruleStreamIdentifier, position1028)
			}
			return true
		l1027:
			position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
			return false
		},
		/* 91 SourceSinkType <- <(<ident> Action69)> */
		func() bool {
			position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
			{
				position1031 := position
				depth++
				{
					position1032 := position
					depth++
					if !_rules[ruleident]() {
						goto l1030
					}
					depth--
					add(rulePegText, position1032)
				}
				if !_rules[ruleAction69]() {
					goto l1030
				}
				depth--
				add(ruleSourceSinkType, position1031)
			}
			return true
		l1030:
			position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
			return false
		},
		/* 92 SourceSinkParamKey <- <(<ident> Action70)> */
		func() bool {
			position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
			{
				position1034 := position
				depth++
				{
					position1035 := position
					depth++
					if !_rules[ruleident]() {
						goto l1033
					}
					depth--
					add(rulePegText, position1035)
				}
				if !_rules[ruleAction70]() {
					goto l1033
				}
				depth--
				add(ruleSourceSinkParamKey, position1034)
			}
			return true
		l1033:
			position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
			return false
		},
		/* 93 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action71)> */
		func() bool {
			position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
			{
				position1037 := position
				depth++
				{
					position1038 := position
					depth++
					{
						position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1040
						}
						position++
						goto l1039
					l1040:
						position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
						if buffer[position] != rune('P') {
							goto l1036
						}
						position++
					}
				l1039:
					{
						position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1042
						}
						position++
						goto l1041
					l1042:
						position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
						if buffer[position] != rune('A') {
							goto l1036
						}
						position++
					}
				l1041:
					{
						position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1044
						}
						position++
						goto l1043
					l1044:
						position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
						if buffer[position] != rune('U') {
							goto l1036
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
							goto l1036
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
							goto l1036
						}
						position++
					}
				l1047:
					{
						position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1050
						}
						position++
						goto l1049
					l1050:
						position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
						if buffer[position] != rune('D') {
							goto l1036
						}
						position++
					}
				l1049:
					depth--
					add(rulePegText, position1038)
				}
				if !_rules[ruleAction71]() {
					goto l1036
				}
				depth--
				add(rulePaused, position1037)
			}
			return true
		l1036:
			position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
			return false
		},
		/* 94 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action72)> */
		func() bool {
			position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
			{
				position1052 := position
				depth++
				{
					position1053 := position
					depth++
					{
						position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1055
						}
						position++
						goto l1054
					l1055:
						position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
						if buffer[position] != rune('U') {
							goto l1051
						}
						position++
					}
				l1054:
					{
						position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1057
						}
						position++
						goto l1056
					l1057:
						position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
						if buffer[position] != rune('N') {
							goto l1051
						}
						position++
					}
				l1056:
					{
						position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1059
						}
						position++
						goto l1058
					l1059:
						position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
						if buffer[position] != rune('P') {
							goto l1051
						}
						position++
					}
				l1058:
					{
						position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1061
						}
						position++
						goto l1060
					l1061:
						position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
						if buffer[position] != rune('A') {
							goto l1051
						}
						position++
					}
				l1060:
					{
						position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1063
						}
						position++
						goto l1062
					l1063:
						position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
						if buffer[position] != rune('U') {
							goto l1051
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
							goto l1051
						}
						position++
					}
				l1064:
					{
						position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1067
						}
						position++
						goto l1066
					l1067:
						position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
						if buffer[position] != rune('E') {
							goto l1051
						}
						position++
					}
				l1066:
					{
						position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1069
						}
						position++
						goto l1068
					l1069:
						position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
						if buffer[position] != rune('D') {
							goto l1051
						}
						position++
					}
				l1068:
					depth--
					add(rulePegText, position1053)
				}
				if !_rules[ruleAction72]() {
					goto l1051
				}
				depth--
				add(ruleUnpaused, position1052)
			}
			return true
		l1051:
			position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
			return false
		},
		/* 95 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
			{
				position1071 := position
				depth++
				{
					position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1073
					}
					goto l1072
				l1073:
					position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
					if !_rules[ruleInt]() {
						goto l1074
					}
					goto l1072
				l1074:
					position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
					if !_rules[ruleFloat]() {
						goto l1075
					}
					goto l1072
				l1075:
					position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
					if !_rules[ruleString]() {
						goto l1076
					}
					goto l1072
				l1076:
					position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
					if !_rules[ruleBlob]() {
						goto l1077
					}
					goto l1072
				l1077:
					position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
					if !_rules[ruleTimestamp]() {
						goto l1078
					}
					goto l1072
				l1078:
					position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
					if !_rules[ruleArray]() {
						goto l1079
					}
					goto l1072
				l1079:
					position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
					if !_rules[ruleMap]() {
						goto l1070
					}
				}
			l1072:
				depth--
				add(ruleType, position1071)
			}
			return true
		l1070:
			position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
			return false
		},
		/* 96 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action73)> */
		func() bool {
			position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
			{
				position1081 := position
				depth++
				{
					position1082 := position
					depth++
					{
						position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1084
						}
						position++
						goto l1083
					l1084:
						position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
						if buffer[position] != rune('B') {
							goto l1080
						}
						position++
					}
				l1083:
					{
						position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1086
						}
						position++
						goto l1085
					l1086:
						position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
						if buffer[position] != rune('O') {
							goto l1080
						}
						position++
					}
				l1085:
					{
						position1087, tokenIndex1087, depth1087 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1088
						}
						position++
						goto l1087
					l1088:
						position, tokenIndex, depth = position1087, tokenIndex1087, depth1087
						if buffer[position] != rune('O') {
							goto l1080
						}
						position++
					}
				l1087:
					{
						position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1090
						}
						position++
						goto l1089
					l1090:
						position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
						if buffer[position] != rune('L') {
							goto l1080
						}
						position++
					}
				l1089:
					depth--
					add(rulePegText, position1082)
				}
				if !_rules[ruleAction73]() {
					goto l1080
				}
				depth--
				add(ruleBool, position1081)
			}
			return true
		l1080:
			position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
			return false
		},
		/* 97 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action74)> */
		func() bool {
			position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
			{
				position1092 := position
				depth++
				{
					position1093 := position
					depth++
					{
						position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1095
						}
						position++
						goto l1094
					l1095:
						position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
						if buffer[position] != rune('I') {
							goto l1091
						}
						position++
					}
				l1094:
					{
						position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1097
						}
						position++
						goto l1096
					l1097:
						position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
						if buffer[position] != rune('N') {
							goto l1091
						}
						position++
					}
				l1096:
					{
						position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1099
						}
						position++
						goto l1098
					l1099:
						position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
						if buffer[position] != rune('T') {
							goto l1091
						}
						position++
					}
				l1098:
					depth--
					add(rulePegText, position1093)
				}
				if !_rules[ruleAction74]() {
					goto l1091
				}
				depth--
				add(ruleInt, position1092)
			}
			return true
		l1091:
			position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
			return false
		},
		/* 98 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action75)> */
		func() bool {
			position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
			{
				position1101 := position
				depth++
				{
					position1102 := position
					depth++
					{
						position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1104
						}
						position++
						goto l1103
					l1104:
						position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
						if buffer[position] != rune('F') {
							goto l1100
						}
						position++
					}
				l1103:
					{
						position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1106
						}
						position++
						goto l1105
					l1106:
						position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
						if buffer[position] != rune('L') {
							goto l1100
						}
						position++
					}
				l1105:
					{
						position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1108
						}
						position++
						goto l1107
					l1108:
						position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
						if buffer[position] != rune('O') {
							goto l1100
						}
						position++
					}
				l1107:
					{
						position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1110
						}
						position++
						goto l1109
					l1110:
						position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
						if buffer[position] != rune('A') {
							goto l1100
						}
						position++
					}
				l1109:
					{
						position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1112
						}
						position++
						goto l1111
					l1112:
						position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
						if buffer[position] != rune('T') {
							goto l1100
						}
						position++
					}
				l1111:
					depth--
					add(rulePegText, position1102)
				}
				if !_rules[ruleAction75]() {
					goto l1100
				}
				depth--
				add(ruleFloat, position1101)
			}
			return true
		l1100:
			position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
			return false
		},
		/* 99 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action76)> */
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
						if buffer[position] != rune('s') {
							goto l1117
						}
						position++
						goto l1116
					l1117:
						position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
						if buffer[position] != rune('S') {
							goto l1113
						}
						position++
					}
				l1116:
					{
						position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1119
						}
						position++
						goto l1118
					l1119:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						if buffer[position] != rune('T') {
							goto l1113
						}
						position++
					}
				l1118:
					{
						position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1121
						}
						position++
						goto l1120
					l1121:
						position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
						if buffer[position] != rune('R') {
							goto l1113
						}
						position++
					}
				l1120:
					{
						position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1123
						}
						position++
						goto l1122
					l1123:
						position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
						if buffer[position] != rune('I') {
							goto l1113
						}
						position++
					}
				l1122:
					{
						position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1125
						}
						position++
						goto l1124
					l1125:
						position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
						if buffer[position] != rune('N') {
							goto l1113
						}
						position++
					}
				l1124:
					{
						position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1127
						}
						position++
						goto l1126
					l1127:
						position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
						if buffer[position] != rune('G') {
							goto l1113
						}
						position++
					}
				l1126:
					depth--
					add(rulePegText, position1115)
				}
				if !_rules[ruleAction76]() {
					goto l1113
				}
				depth--
				add(ruleString, position1114)
			}
			return true
		l1113:
			position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
			return false
		},
		/* 100 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action77)> */
		func() bool {
			position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
			{
				position1129 := position
				depth++
				{
					position1130 := position
					depth++
					{
						position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1132
						}
						position++
						goto l1131
					l1132:
						position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
						if buffer[position] != rune('B') {
							goto l1128
						}
						position++
					}
				l1131:
					{
						position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1134
						}
						position++
						goto l1133
					l1134:
						position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
						if buffer[position] != rune('L') {
							goto l1128
						}
						position++
					}
				l1133:
					{
						position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1136
						}
						position++
						goto l1135
					l1136:
						position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
						if buffer[position] != rune('O') {
							goto l1128
						}
						position++
					}
				l1135:
					{
						position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1138
						}
						position++
						goto l1137
					l1138:
						position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
						if buffer[position] != rune('B') {
							goto l1128
						}
						position++
					}
				l1137:
					depth--
					add(rulePegText, position1130)
				}
				if !_rules[ruleAction77]() {
					goto l1128
				}
				depth--
				add(ruleBlob, position1129)
			}
			return true
		l1128:
			position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
			return false
		},
		/* 101 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action78)> */
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
						if buffer[position] != rune('t') {
							goto l1143
						}
						position++
						goto l1142
					l1143:
						position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
						if buffer[position] != rune('T') {
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
						if buffer[position] != rune('m') {
							goto l1147
						}
						position++
						goto l1146
					l1147:
						position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
						if buffer[position] != rune('M') {
							goto l1139
						}
						position++
					}
				l1146:
					{
						position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1149
						}
						position++
						goto l1148
					l1149:
						position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
						if buffer[position] != rune('E') {
							goto l1139
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
							goto l1139
						}
						position++
					}
				l1150:
					{
						position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1153
						}
						position++
						goto l1152
					l1153:
						position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
						if buffer[position] != rune('T') {
							goto l1139
						}
						position++
					}
				l1152:
					{
						position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1155
						}
						position++
						goto l1154
					l1155:
						position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
						if buffer[position] != rune('A') {
							goto l1139
						}
						position++
					}
				l1154:
					{
						position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1157
						}
						position++
						goto l1156
					l1157:
						position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
						if buffer[position] != rune('M') {
							goto l1139
						}
						position++
					}
				l1156:
					{
						position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1159
						}
						position++
						goto l1158
					l1159:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						if buffer[position] != rune('P') {
							goto l1139
						}
						position++
					}
				l1158:
					depth--
					add(rulePegText, position1141)
				}
				if !_rules[ruleAction78]() {
					goto l1139
				}
				depth--
				add(ruleTimestamp, position1140)
			}
			return true
		l1139:
			position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
			return false
		},
		/* 102 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action79)> */
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
						if buffer[position] != rune('a') {
							goto l1164
						}
						position++
						goto l1163
					l1164:
						position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
						if buffer[position] != rune('A') {
							goto l1160
						}
						position++
					}
				l1163:
					{
						position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1166
						}
						position++
						goto l1165
					l1166:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						if buffer[position] != rune('R') {
							goto l1160
						}
						position++
					}
				l1165:
					{
						position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1168
						}
						position++
						goto l1167
					l1168:
						position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
						if buffer[position] != rune('R') {
							goto l1160
						}
						position++
					}
				l1167:
					{
						position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1170
						}
						position++
						goto l1169
					l1170:
						position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
						if buffer[position] != rune('A') {
							goto l1160
						}
						position++
					}
				l1169:
					{
						position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1172
						}
						position++
						goto l1171
					l1172:
						position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
						if buffer[position] != rune('Y') {
							goto l1160
						}
						position++
					}
				l1171:
					depth--
					add(rulePegText, position1162)
				}
				if !_rules[ruleAction79]() {
					goto l1160
				}
				depth--
				add(ruleArray, position1161)
			}
			return true
		l1160:
			position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
			return false
		},
		/* 103 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action80)> */
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
						if buffer[position] != rune('m') {
							goto l1177
						}
						position++
						goto l1176
					l1177:
						position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
						if buffer[position] != rune('M') {
							goto l1173
						}
						position++
					}
				l1176:
					{
						position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1179
						}
						position++
						goto l1178
					l1179:
						position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
						if buffer[position] != rune('A') {
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
					depth--
					add(rulePegText, position1175)
				}
				if !_rules[ruleAction80]() {
					goto l1173
				}
				depth--
				add(ruleMap, position1174)
			}
			return true
		l1173:
			position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
			return false
		},
		/* 104 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action81)> */
		func() bool {
			position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
			{
				position1183 := position
				depth++
				{
					position1184 := position
					depth++
					{
						position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1186
						}
						position++
						goto l1185
					l1186:
						position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
						if buffer[position] != rune('O') {
							goto l1182
						}
						position++
					}
				l1185:
					{
						position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1188
						}
						position++
						goto l1187
					l1188:
						position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
						if buffer[position] != rune('R') {
							goto l1182
						}
						position++
					}
				l1187:
					depth--
					add(rulePegText, position1184)
				}
				if !_rules[ruleAction81]() {
					goto l1182
				}
				depth--
				add(ruleOr, position1183)
			}
			return true
		l1182:
			position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
			return false
		},
		/* 105 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action82)> */
		func() bool {
			position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
			{
				position1190 := position
				depth++
				{
					position1191 := position
					depth++
					{
						position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1193
						}
						position++
						goto l1192
					l1193:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						if buffer[position] != rune('A') {
							goto l1189
						}
						position++
					}
				l1192:
					{
						position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1195
						}
						position++
						goto l1194
					l1195:
						position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
						if buffer[position] != rune('N') {
							goto l1189
						}
						position++
					}
				l1194:
					{
						position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1197
						}
						position++
						goto l1196
					l1197:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
						if buffer[position] != rune('D') {
							goto l1189
						}
						position++
					}
				l1196:
					depth--
					add(rulePegText, position1191)
				}
				if !_rules[ruleAction82]() {
					goto l1189
				}
				depth--
				add(ruleAnd, position1190)
			}
			return true
		l1189:
			position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
			return false
		},
		/* 106 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action83)> */
		func() bool {
			position1198, tokenIndex1198, depth1198 := position, tokenIndex, depth
			{
				position1199 := position
				depth++
				{
					position1200 := position
					depth++
					{
						position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1202
						}
						position++
						goto l1201
					l1202:
						position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
						if buffer[position] != rune('N') {
							goto l1198
						}
						position++
					}
				l1201:
					{
						position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1204
						}
						position++
						goto l1203
					l1204:
						position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
						if buffer[position] != rune('O') {
							goto l1198
						}
						position++
					}
				l1203:
					{
						position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1206
						}
						position++
						goto l1205
					l1206:
						position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
						if buffer[position] != rune('T') {
							goto l1198
						}
						position++
					}
				l1205:
					depth--
					add(rulePegText, position1200)
				}
				if !_rules[ruleAction83]() {
					goto l1198
				}
				depth--
				add(ruleNot, position1199)
			}
			return true
		l1198:
			position, tokenIndex, depth = position1198, tokenIndex1198, depth1198
			return false
		},
		/* 107 Equal <- <(<'='> Action84)> */
		func() bool {
			position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
			{
				position1208 := position
				depth++
				{
					position1209 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1207
					}
					position++
					depth--
					add(rulePegText, position1209)
				}
				if !_rules[ruleAction84]() {
					goto l1207
				}
				depth--
				add(ruleEqual, position1208)
			}
			return true
		l1207:
			position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
			return false
		},
		/* 108 Less <- <(<'<'> Action85)> */
		func() bool {
			position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
			{
				position1211 := position
				depth++
				{
					position1212 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1210
					}
					position++
					depth--
					add(rulePegText, position1212)
				}
				if !_rules[ruleAction85]() {
					goto l1210
				}
				depth--
				add(ruleLess, position1211)
			}
			return true
		l1210:
			position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
			return false
		},
		/* 109 LessOrEqual <- <(<('<' '=')> Action86)> */
		func() bool {
			position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
			{
				position1214 := position
				depth++
				{
					position1215 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1213
					}
					position++
					if buffer[position] != rune('=') {
						goto l1213
					}
					position++
					depth--
					add(rulePegText, position1215)
				}
				if !_rules[ruleAction86]() {
					goto l1213
				}
				depth--
				add(ruleLessOrEqual, position1214)
			}
			return true
		l1213:
			position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
			return false
		},
		/* 110 Greater <- <(<'>'> Action87)> */
		func() bool {
			position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
			{
				position1217 := position
				depth++
				{
					position1218 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1216
					}
					position++
					depth--
					add(rulePegText, position1218)
				}
				if !_rules[ruleAction87]() {
					goto l1216
				}
				depth--
				add(ruleGreater, position1217)
			}
			return true
		l1216:
			position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
			return false
		},
		/* 111 GreaterOrEqual <- <(<('>' '=')> Action88)> */
		func() bool {
			position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
			{
				position1220 := position
				depth++
				{
					position1221 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1219
					}
					position++
					if buffer[position] != rune('=') {
						goto l1219
					}
					position++
					depth--
					add(rulePegText, position1221)
				}
				if !_rules[ruleAction88]() {
					goto l1219
				}
				depth--
				add(ruleGreaterOrEqual, position1220)
			}
			return true
		l1219:
			position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
			return false
		},
		/* 112 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action89)> */
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
						if buffer[position] != rune('!') {
							goto l1226
						}
						position++
						if buffer[position] != rune('=') {
							goto l1226
						}
						position++
						goto l1225
					l1226:
						position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
						if buffer[position] != rune('<') {
							goto l1222
						}
						position++
						if buffer[position] != rune('>') {
							goto l1222
						}
						position++
					}
				l1225:
					depth--
					add(rulePegText, position1224)
				}
				if !_rules[ruleAction89]() {
					goto l1222
				}
				depth--
				add(ruleNotEqual, position1223)
			}
			return true
		l1222:
			position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
			return false
		},
		/* 113 Concat <- <(<('|' '|')> Action90)> */
		func() bool {
			position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
			{
				position1228 := position
				depth++
				{
					position1229 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1227
					}
					position++
					if buffer[position] != rune('|') {
						goto l1227
					}
					position++
					depth--
					add(rulePegText, position1229)
				}
				if !_rules[ruleAction90]() {
					goto l1227
				}
				depth--
				add(ruleConcat, position1228)
			}
			return true
		l1227:
			position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
			return false
		},
		/* 114 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action91)> */
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
						if buffer[position] != rune('s') {
							goto l1236
						}
						position++
						goto l1235
					l1236:
						position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
						if buffer[position] != rune('S') {
							goto l1230
						}
						position++
					}
				l1235:
					depth--
					add(rulePegText, position1232)
				}
				if !_rules[ruleAction91]() {
					goto l1230
				}
				depth--
				add(ruleIs, position1231)
			}
			return true
		l1230:
			position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
			return false
		},
		/* 115 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action92)> */
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
						if buffer[position] != rune('i') {
							goto l1241
						}
						position++
						goto l1240
					l1241:
						position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
						if buffer[position] != rune('I') {
							goto l1237
						}
						position++
					}
				l1240:
					{
						position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1243
						}
						position++
						goto l1242
					l1243:
						position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
						if buffer[position] != rune('S') {
							goto l1237
						}
						position++
					}
				l1242:
					if !_rules[rulesp]() {
						goto l1237
					}
					{
						position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1245
						}
						position++
						goto l1244
					l1245:
						position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
						if buffer[position] != rune('N') {
							goto l1237
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
							goto l1237
						}
						position++
					}
				l1246:
					{
						position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1249
						}
						position++
						goto l1248
					l1249:
						position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
						if buffer[position] != rune('T') {
							goto l1237
						}
						position++
					}
				l1248:
					depth--
					add(rulePegText, position1239)
				}
				if !_rules[ruleAction92]() {
					goto l1237
				}
				depth--
				add(ruleIsNot, position1238)
			}
			return true
		l1237:
			position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
			return false
		},
		/* 116 Plus <- <(<'+'> Action93)> */
		func() bool {
			position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
			{
				position1251 := position
				depth++
				{
					position1252 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1250
					}
					position++
					depth--
					add(rulePegText, position1252)
				}
				if !_rules[ruleAction93]() {
					goto l1250
				}
				depth--
				add(rulePlus, position1251)
			}
			return true
		l1250:
			position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
			return false
		},
		/* 117 Minus <- <(<'-'> Action94)> */
		func() bool {
			position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
			{
				position1254 := position
				depth++
				{
					position1255 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1253
					}
					position++
					depth--
					add(rulePegText, position1255)
				}
				if !_rules[ruleAction94]() {
					goto l1253
				}
				depth--
				add(ruleMinus, position1254)
			}
			return true
		l1253:
			position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
			return false
		},
		/* 118 Multiply <- <(<'*'> Action95)> */
		func() bool {
			position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
			{
				position1257 := position
				depth++
				{
					position1258 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1256
					}
					position++
					depth--
					add(rulePegText, position1258)
				}
				if !_rules[ruleAction95]() {
					goto l1256
				}
				depth--
				add(ruleMultiply, position1257)
			}
			return true
		l1256:
			position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
			return false
		},
		/* 119 Divide <- <(<'/'> Action96)> */
		func() bool {
			position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
			{
				position1260 := position
				depth++
				{
					position1261 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1259
					}
					position++
					depth--
					add(rulePegText, position1261)
				}
				if !_rules[ruleAction96]() {
					goto l1259
				}
				depth--
				add(ruleDivide, position1260)
			}
			return true
		l1259:
			position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
			return false
		},
		/* 120 Modulo <- <(<'%'> Action97)> */
		func() bool {
			position1262, tokenIndex1262, depth1262 := position, tokenIndex, depth
			{
				position1263 := position
				depth++
				{
					position1264 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1262
					}
					position++
					depth--
					add(rulePegText, position1264)
				}
				if !_rules[ruleAction97]() {
					goto l1262
				}
				depth--
				add(ruleModulo, position1263)
			}
			return true
		l1262:
			position, tokenIndex, depth = position1262, tokenIndex1262, depth1262
			return false
		},
		/* 121 UnaryMinus <- <(<'-'> Action98)> */
		func() bool {
			position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
			{
				position1266 := position
				depth++
				{
					position1267 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1265
					}
					position++
					depth--
					add(rulePegText, position1267)
				}
				if !_rules[ruleAction98]() {
					goto l1265
				}
				depth--
				add(ruleUnaryMinus, position1266)
			}
			return true
		l1265:
			position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
			return false
		},
		/* 122 Identifier <- <(<ident> Action99)> */
		func() bool {
			position1268, tokenIndex1268, depth1268 := position, tokenIndex, depth
			{
				position1269 := position
				depth++
				{
					position1270 := position
					depth++
					if !_rules[ruleident]() {
						goto l1268
					}
					depth--
					add(rulePegText, position1270)
				}
				if !_rules[ruleAction99]() {
					goto l1268
				}
				depth--
				add(ruleIdentifier, position1269)
			}
			return true
		l1268:
			position, tokenIndex, depth = position1268, tokenIndex1268, depth1268
			return false
		},
		/* 123 TargetIdentifier <- <(<jsonPath> Action100)> */
		func() bool {
			position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
			{
				position1272 := position
				depth++
				{
					position1273 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1271
					}
					depth--
					add(rulePegText, position1273)
				}
				if !_rules[ruleAction100]() {
					goto l1271
				}
				depth--
				add(ruleTargetIdentifier, position1272)
			}
			return true
		l1271:
			position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
			return false
		},
		/* 124 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
			{
				position1275 := position
				depth++
				{
					position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1277
					}
					position++
					goto l1276
				l1277:
					position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1274
					}
					position++
				}
			l1276:
			l1278:
				{
					position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
					{
						position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1281
						}
						position++
						goto l1280
					l1281:
						position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1282
						}
						position++
						goto l1280
					l1282:
						position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1283
						}
						position++
						goto l1280
					l1283:
						position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
						if buffer[position] != rune('_') {
							goto l1279
						}
						position++
					}
				l1280:
					goto l1278
				l1279:
					position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
				}
				depth--
				add(ruleident, position1275)
			}
			return true
		l1274:
			position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
			return false
		},
		/* 125 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1284, tokenIndex1284, depth1284 := position, tokenIndex, depth
			{
				position1285 := position
				depth++
				{
					position1286, tokenIndex1286, depth1286 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1287
					}
					position++
					goto l1286
				l1287:
					position, tokenIndex, depth = position1286, tokenIndex1286, depth1286
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1284
					}
					position++
				}
			l1286:
			l1288:
				{
					position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
					{
						position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1291
						}
						position++
						goto l1290
					l1291:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1292
						}
						position++
						goto l1290
					l1292:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1293
						}
						position++
						goto l1290
					l1293:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if buffer[position] != rune('_') {
							goto l1294
						}
						position++
						goto l1290
					l1294:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if buffer[position] != rune('.') {
							goto l1295
						}
						position++
						goto l1290
					l1295:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if buffer[position] != rune('[') {
							goto l1296
						}
						position++
						goto l1290
					l1296:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if buffer[position] != rune(']') {
							goto l1297
						}
						position++
						goto l1290
					l1297:
						position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
						if buffer[position] != rune('"') {
							goto l1289
						}
						position++
					}
				l1290:
					goto l1288
				l1289:
					position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
				}
				depth--
				add(rulejsonPath, position1285)
			}
			return true
		l1284:
			position, tokenIndex, depth = position1284, tokenIndex1284, depth1284
			return false
		},
		/* 126 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1299 := position
				depth++
			l1300:
				{
					position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
					{
						position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1303
						}
						position++
						goto l1302
					l1303:
						position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
						if buffer[position] != rune('\t') {
							goto l1304
						}
						position++
						goto l1302
					l1304:
						position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
						if buffer[position] != rune('\n') {
							goto l1305
						}
						position++
						goto l1302
					l1305:
						position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
						if buffer[position] != rune('\r') {
							goto l1306
						}
						position++
						goto l1302
					l1306:
						position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
						if !_rules[rulecomment]() {
							goto l1301
						}
					}
				l1302:
					goto l1300
				l1301:
					position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
				}
				depth--
				add(rulesp, position1299)
			}
			return true
		},
		/* 127 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
			{
				position1308 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1307
				}
				position++
				if buffer[position] != rune('-') {
					goto l1307
				}
				position++
			l1309:
				{
					position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
					{
						position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
						{
							position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1313
							}
							position++
							goto l1312
						l1313:
							position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
							if buffer[position] != rune('\n') {
								goto l1311
							}
							position++
						}
					l1312:
						goto l1310
					l1311:
						position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
					}
					if !matchDot() {
						goto l1310
					}
					goto l1309
				l1310:
					position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
				}
				{
					position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1315
					}
					position++
					goto l1314
				l1315:
					position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
					if buffer[position] != rune('\n') {
						goto l1307
					}
					position++
				}
			l1314:
				depth--
				add(rulecomment, position1308)
			}
			return true
		l1307:
			position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
			return false
		},
		/* 129 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 131 Action1 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 132 Action2 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 133 Action3 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 134 Action4 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 135 Action5 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 136 Action6 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 137 Action7 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 138 Action8 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 139 Action9 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 140 Action10 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 141 Action11 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 142 Action12 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 143 Action13 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 144 Action14 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 145 Action15 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 146 Action16 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 147 Action17 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 148 Action18 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 149 Action19 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 150 Action20 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 151 Action21 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 152 Action22 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 153 Action23 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 154 Action24 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 155 Action25 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 156 Action26 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 157 Action27 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 158 Action28 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 159 Action29 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 160 Action30 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 161 Action31 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 162 Action32 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 163 Action33 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 164 Action34 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 165 Action35 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 166 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 167 Action37 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 168 Action38 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 169 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 170 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 171 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 172 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 173 Action43 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 174 Action44 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 175 Action45 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 176 Action46 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 177 Action47 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 178 Action48 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 179 Action49 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 180 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 181 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 182 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 183 Action53 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 184 Action54 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 185 Action55 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 186 Action56 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 187 Action57 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 188 Action58 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 189 Action59 <- <{
		    p.PushComponent(begin, end, NewWildcard(""))
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 190 Action60 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 191 Action61 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 192 Action62 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 193 Action63 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 194 Action64 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 195 Action65 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 196 Action66 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 197 Action67 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 198 Action68 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 199 Action69 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 200 Action70 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 201 Action71 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 202 Action72 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 203 Action73 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 204 Action74 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 205 Action75 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 206 Action76 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 207 Action77 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 208 Action78 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 209 Action79 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 210 Action80 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 211 Action81 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 212 Action82 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 213 Action83 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 214 Action84 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 215 Action85 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 216 Action86 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 217 Action87 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 218 Action88 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 219 Action89 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 220 Action90 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 221 Action91 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 222 Action92 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 223 Action93 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 224 Action94 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 225 Action95 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 226 Action96 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 227 Action97 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 228 Action98 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 229 Action99 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 230 Action100 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
	}
	p.rules = _rules
}

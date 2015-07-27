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
	rules  [223]func() bool
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

			p.AssembleCreateSource()

		case ruleAction4:

			p.AssembleCreateSink()

		case ruleAction5:

			p.AssembleCreateState()

		case ruleAction6:

			p.AssembleUpdateState()

		case ruleAction7:

			p.AssembleUpdateSource()

		case ruleAction8:

			p.AssembleUpdateSink()

		case ruleAction9:

			p.AssembleInsertIntoSelect()

		case ruleAction10:

			p.AssembleInsertIntoFrom()

		case ruleAction11:

			p.AssemblePauseSource()

		case ruleAction12:

			p.AssembleResumeSource()

		case ruleAction13:

			p.AssembleRewindSource()

		case ruleAction14:

			p.AssembleDropSource()

		case ruleAction15:

			p.AssembleDropStream()

		case ruleAction16:

			p.AssembleDropSink()

		case ruleAction17:

			p.AssembleDropState()

		case ruleAction18:

			p.AssembleEmitter()

		case ruleAction19:

			p.AssembleProjections(begin, end)

		case ruleAction20:

			p.AssembleAlias()

		case ruleAction21:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction22:

			p.AssembleInterval()

		case ruleAction23:

			p.AssembleInterval()

		case ruleAction24:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction25:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction26:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction27:

			p.EnsureAliasedStreamWindow()

		case ruleAction28:

			p.AssembleAliasedStreamWindow()

		case ruleAction29:

			p.AssembleStreamWindow()

		case ruleAction30:

			p.AssembleUDSFFuncApp()

		case ruleAction31:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction32:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction33:

			p.AssembleSourceSinkParam()

		case ruleAction34:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction35:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction36:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction37:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction38:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction39:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction40:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction41:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction42:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction43:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction44:

			p.AssembleTypeCast(begin, end)

		case ruleAction45:

			p.AssembleTypeCast(begin, end)

		case ruleAction46:

			p.AssembleFuncApp()

		case ruleAction47:

			p.AssembleExpressions(begin, end)

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction49:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction50:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction51:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction52:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction53:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction54:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction55:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction56:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction57:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction58:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction59:

			p.PushComponent(begin, end, Istream)

		case ruleAction60:

			p.PushComponent(begin, end, Dstream)

		case ruleAction61:

			p.PushComponent(begin, end, Rstream)

		case ruleAction62:

			p.PushComponent(begin, end, Tuples)

		case ruleAction63:

			p.PushComponent(begin, end, Seconds)

		case ruleAction64:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction65:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction66:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction67:

			p.PushComponent(begin, end, Yes)

		case ruleAction68:

			p.PushComponent(begin, end, No)

		case ruleAction69:

			p.PushComponent(begin, end, Bool)

		case ruleAction70:

			p.PushComponent(begin, end, Int)

		case ruleAction71:

			p.PushComponent(begin, end, Float)

		case ruleAction72:

			p.PushComponent(begin, end, String)

		case ruleAction73:

			p.PushComponent(begin, end, Blob)

		case ruleAction74:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction75:

			p.PushComponent(begin, end, Array)

		case ruleAction76:

			p.PushComponent(begin, end, Map)

		case ruleAction77:

			p.PushComponent(begin, end, Or)

		case ruleAction78:

			p.PushComponent(begin, end, And)

		case ruleAction79:

			p.PushComponent(begin, end, Not)

		case ruleAction80:

			p.PushComponent(begin, end, Equal)

		case ruleAction81:

			p.PushComponent(begin, end, Less)

		case ruleAction82:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction83:

			p.PushComponent(begin, end, Greater)

		case ruleAction84:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction85:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction86:

			p.PushComponent(begin, end, Concat)

		case ruleAction87:

			p.PushComponent(begin, end, Is)

		case ruleAction88:

			p.PushComponent(begin, end, IsNot)

		case ruleAction89:

			p.PushComponent(begin, end, Plus)

		case ruleAction90:

			p.PushComponent(begin, end, Minus)

		case ruleAction91:

			p.PushComponent(begin, end, Multiply)

		case ruleAction92:

			p.PushComponent(begin, end, Divide)

		case ruleAction93:

			p.PushComponent(begin, end, Modulo)

		case ruleAction94:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction95:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction96:

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
		/* 5 StreamStmt <- <(CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoSelectStmt / InsertIntoFromStmt)> */
		func() bool {
			position33, tokenIndex33, depth33 := position, tokenIndex, depth
			{
				position34 := position
				depth++
				{
					position35, tokenIndex35, depth35 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l36
					}
					goto l35
				l36:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleDropStreamStmt]() {
						goto l37
					}
					goto l35
				l37:
					position, tokenIndex, depth = position35, tokenIndex35, depth35
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l38
					}
					goto l35
				l38:
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
			position39, tokenIndex39, depth39 := position, tokenIndex, depth
			{
				position40 := position
				depth++
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l42
					}
					position++
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if buffer[position] != rune('S') {
						goto l39
					}
					position++
				}
			l41:
				{
					position43, tokenIndex43, depth43 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l44
					}
					position++
					goto l43
				l44:
					position, tokenIndex, depth = position43, tokenIndex43, depth43
					if buffer[position] != rune('E') {
						goto l39
					}
					position++
				}
			l43:
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l46
					}
					position++
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if buffer[position] != rune('L') {
						goto l39
					}
					position++
				}
			l45:
				{
					position47, tokenIndex47, depth47 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l48
					}
					position++
					goto l47
				l48:
					position, tokenIndex, depth = position47, tokenIndex47, depth47
					if buffer[position] != rune('E') {
						goto l39
					}
					position++
				}
			l47:
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l50
					}
					position++
					goto l49
				l50:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
					if buffer[position] != rune('C') {
						goto l39
					}
					position++
				}
			l49:
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l52
					}
					position++
					goto l51
				l52:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
					if buffer[position] != rune('T') {
						goto l39
					}
					position++
				}
			l51:
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleEmitter]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleProjections]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleWindowedFrom]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleFilter]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleGrouping]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleHaving]() {
					goto l39
				}
				if !_rules[rulesp]() {
					goto l39
				}
				if !_rules[ruleAction0]() {
					goto l39
				}
				depth--
				add(ruleSelectStmt, position40)
			}
			return true
		l39:
			position, tokenIndex, depth = position39, tokenIndex39, depth39
			return false
		},
		/* 7 SelectUnionStmt <- <(<(SelectStmt (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action1)> */
		func() bool {
			position53, tokenIndex53, depth53 := position, tokenIndex, depth
			{
				position54 := position
				depth++
				{
					position55 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l53
					}
					{
						position58, tokenIndex58, depth58 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l59
						}
						position++
						goto l58
					l59:
						position, tokenIndex, depth = position58, tokenIndex58, depth58
						if buffer[position] != rune('U') {
							goto l53
						}
						position++
					}
				l58:
					{
						position60, tokenIndex60, depth60 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l61
						}
						position++
						goto l60
					l61:
						position, tokenIndex, depth = position60, tokenIndex60, depth60
						if buffer[position] != rune('N') {
							goto l53
						}
						position++
					}
				l60:
					{
						position62, tokenIndex62, depth62 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l63
						}
						position++
						goto l62
					l63:
						position, tokenIndex, depth = position62, tokenIndex62, depth62
						if buffer[position] != rune('I') {
							goto l53
						}
						position++
					}
				l62:
					{
						position64, tokenIndex64, depth64 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l65
						}
						position++
						goto l64
					l65:
						position, tokenIndex, depth = position64, tokenIndex64, depth64
						if buffer[position] != rune('O') {
							goto l53
						}
						position++
					}
				l64:
					{
						position66, tokenIndex66, depth66 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l67
						}
						position++
						goto l66
					l67:
						position, tokenIndex, depth = position66, tokenIndex66, depth66
						if buffer[position] != rune('N') {
							goto l53
						}
						position++
					}
				l66:
					if !_rules[rulesp]() {
						goto l53
					}
					{
						position68, tokenIndex68, depth68 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l69
						}
						position++
						goto l68
					l69:
						position, tokenIndex, depth = position68, tokenIndex68, depth68
						if buffer[position] != rune('A') {
							goto l53
						}
						position++
					}
				l68:
					{
						position70, tokenIndex70, depth70 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l71
						}
						position++
						goto l70
					l71:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
						if buffer[position] != rune('L') {
							goto l53
						}
						position++
					}
				l70:
					{
						position72, tokenIndex72, depth72 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l73
						}
						position++
						goto l72
					l73:
						position, tokenIndex, depth = position72, tokenIndex72, depth72
						if buffer[position] != rune('L') {
							goto l53
						}
						position++
					}
				l72:
					if !_rules[rulesp]() {
						goto l53
					}
					if !_rules[ruleSelectStmt]() {
						goto l53
					}
				l56:
					{
						position57, tokenIndex57, depth57 := position, tokenIndex, depth
						{
							position74, tokenIndex74, depth74 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l75
							}
							position++
							goto l74
						l75:
							position, tokenIndex, depth = position74, tokenIndex74, depth74
							if buffer[position] != rune('U') {
								goto l57
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
								goto l57
							}
							position++
						}
					l76:
						{
							position78, tokenIndex78, depth78 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l79
							}
							position++
							goto l78
						l79:
							position, tokenIndex, depth = position78, tokenIndex78, depth78
							if buffer[position] != rune('I') {
								goto l57
							}
							position++
						}
					l78:
						{
							position80, tokenIndex80, depth80 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l81
							}
							position++
							goto l80
						l81:
							position, tokenIndex, depth = position80, tokenIndex80, depth80
							if buffer[position] != rune('O') {
								goto l57
							}
							position++
						}
					l80:
						{
							position82, tokenIndex82, depth82 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l83
							}
							position++
							goto l82
						l83:
							position, tokenIndex, depth = position82, tokenIndex82, depth82
							if buffer[position] != rune('N') {
								goto l57
							}
							position++
						}
					l82:
						if !_rules[rulesp]() {
							goto l57
						}
						{
							position84, tokenIndex84, depth84 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l85
							}
							position++
							goto l84
						l85:
							position, tokenIndex, depth = position84, tokenIndex84, depth84
							if buffer[position] != rune('A') {
								goto l57
							}
							position++
						}
					l84:
						{
							position86, tokenIndex86, depth86 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l87
							}
							position++
							goto l86
						l87:
							position, tokenIndex, depth = position86, tokenIndex86, depth86
							if buffer[position] != rune('L') {
								goto l57
							}
							position++
						}
					l86:
						{
							position88, tokenIndex88, depth88 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l89
							}
							position++
							goto l88
						l89:
							position, tokenIndex, depth = position88, tokenIndex88, depth88
							if buffer[position] != rune('L') {
								goto l57
							}
							position++
						}
					l88:
						if !_rules[rulesp]() {
							goto l57
						}
						if !_rules[ruleSelectStmt]() {
							goto l57
						}
						goto l56
					l57:
						position, tokenIndex, depth = position57, tokenIndex57, depth57
					}
					depth--
					add(rulePegText, position55)
				}
				if !_rules[ruleAction1]() {
					goto l53
				}
				depth--
				add(ruleSelectUnionStmt, position54)
			}
			return true
		l53:
			position, tokenIndex, depth = position53, tokenIndex53, depth53
			return false
		},
		/* 8 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action2)> */
		func() bool {
			position90, tokenIndex90, depth90 := position, tokenIndex, depth
			{
				position91 := position
				depth++
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('C') {
						goto l90
					}
					position++
				}
			l92:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('R') {
						goto l90
					}
					position++
				}
			l94:
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('E') {
						goto l90
					}
					position++
				}
			l96:
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('A') {
						goto l90
					}
					position++
				}
			l98:
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l101
					}
					position++
					goto l100
				l101:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
					if buffer[position] != rune('T') {
						goto l90
					}
					position++
				}
			l100:
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('E') {
						goto l90
					}
					position++
				}
			l102:
				if !_rules[rulesp]() {
					goto l90
				}
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('S') {
						goto l90
					}
					position++
				}
			l104:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('T') {
						goto l90
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('R') {
						goto l90
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('E') {
						goto l90
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('A') {
						goto l90
					}
					position++
				}
			l112:
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('M') {
						goto l90
					}
					position++
				}
			l114:
				if !_rules[rulesp]() {
					goto l90
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l90
				}
				if !_rules[rulesp]() {
					goto l90
				}
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('A') {
						goto l90
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('S') {
						goto l90
					}
					position++
				}
			l118:
				if !_rules[rulesp]() {
					goto l90
				}
				if !_rules[ruleSelectStmt]() {
					goto l90
				}
				if !_rules[ruleAction2]() {
					goto l90
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position91)
			}
			return true
		l90:
			position, tokenIndex, depth = position90, tokenIndex90, depth90
			return false
		},
		/* 9 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
		func() bool {
			position120, tokenIndex120, depth120 := position, tokenIndex, depth
			{
				position121 := position
				depth++
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('C') {
						goto l120
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('R') {
						goto l120
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('E') {
						goto l120
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('A') {
						goto l120
					}
					position++
				}
			l128:
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('T') {
						goto l120
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('E') {
						goto l120
					}
					position++
				}
			l132:
				if !_rules[rulesp]() {
					goto l120
				}
				if !_rules[rulePausedOpt]() {
					goto l120
				}
				if !_rules[rulesp]() {
					goto l120
				}
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('S') {
						goto l120
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('O') {
						goto l120
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('U') {
						goto l120
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('R') {
						goto l120
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('C') {
						goto l120
					}
					position++
				}
			l142:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('E') {
						goto l120
					}
					position++
				}
			l144:
				if !_rules[rulesp]() {
					goto l120
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l120
				}
				if !_rules[rulesp]() {
					goto l120
				}
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
						goto l120
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('Y') {
						goto l120
					}
					position++
				}
			l148:
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('P') {
						goto l120
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('E') {
						goto l120
					}
					position++
				}
			l152:
				if !_rules[rulesp]() {
					goto l120
				}
				if !_rules[ruleSourceSinkType]() {
					goto l120
				}
				if !_rules[rulesp]() {
					goto l120
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l120
				}
				if !_rules[ruleAction3]() {
					goto l120
				}
				depth--
				add(ruleCreateSourceStmt, position121)
			}
			return true
		l120:
			position, tokenIndex, depth = position120, tokenIndex120, depth120
			return false
		},
		/* 10 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action4)> */
		func() bool {
			position154, tokenIndex154, depth154 := position, tokenIndex, depth
			{
				position155 := position
				depth++
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('C') {
						goto l154
					}
					position++
				}
			l156:
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if buffer[position] != rune('R') {
						goto l154
					}
					position++
				}
			l158:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('E') {
						goto l154
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('A') {
						goto l154
					}
					position++
				}
			l162:
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('T') {
						goto l154
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
						goto l154
					}
					position++
				}
			l166:
				if !_rules[rulesp]() {
					goto l154
				}
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('S') {
						goto l154
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('I') {
						goto l154
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('N') {
						goto l154
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('K') {
						goto l154
					}
					position++
				}
			l174:
				if !_rules[rulesp]() {
					goto l154
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l154
				}
				if !_rules[rulesp]() {
					goto l154
				}
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('T') {
						goto l154
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('Y') {
						goto l154
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('P') {
						goto l154
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('E') {
						goto l154
					}
					position++
				}
			l182:
				if !_rules[rulesp]() {
					goto l154
				}
				if !_rules[ruleSourceSinkType]() {
					goto l154
				}
				if !_rules[rulesp]() {
					goto l154
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l154
				}
				if !_rules[ruleAction4]() {
					goto l154
				}
				depth--
				add(ruleCreateSinkStmt, position155)
			}
			return true
		l154:
			position, tokenIndex, depth = position154, tokenIndex154, depth154
			return false
		},
		/* 11 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action5)> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('C') {
						goto l184
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('R') {
						goto l184
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('E') {
						goto l184
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('A') {
						goto l184
					}
					position++
				}
			l192:
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('T') {
						goto l184
					}
					position++
				}
			l194:
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('E') {
						goto l184
					}
					position++
				}
			l196:
				if !_rules[rulesp]() {
					goto l184
				}
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('S') {
						goto l184
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('T') {
						goto l184
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
						goto l184
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
						goto l184
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
						goto l184
					}
					position++
				}
			l206:
				if !_rules[rulesp]() {
					goto l184
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l184
				}
				if !_rules[rulesp]() {
					goto l184
				}
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('T') {
						goto l184
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('Y') {
						goto l184
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('P') {
						goto l184
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('E') {
						goto l184
					}
					position++
				}
			l214:
				if !_rules[rulesp]() {
					goto l184
				}
				if !_rules[ruleSourceSinkType]() {
					goto l184
				}
				if !_rules[rulesp]() {
					goto l184
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l184
				}
				if !_rules[ruleAction5]() {
					goto l184
				}
				depth--
				add(ruleCreateStateStmt, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
		/* 12 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action6)> */
		func() bool {
			position216, tokenIndex216, depth216 := position, tokenIndex, depth
			{
				position217 := position
				depth++
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('U') {
						goto l216
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
						goto l216
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('D') {
						goto l216
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('A') {
						goto l216
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('T') {
						goto l216
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('E') {
						goto l216
					}
					position++
				}
			l228:
				if !_rules[rulesp]() {
					goto l216
				}
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
					if buffer[position] != rune('S') {
						goto l216
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('T') {
						goto l216
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('A') {
						goto l216
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('T') {
						goto l216
					}
					position++
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('E') {
						goto l216
					}
					position++
				}
			l238:
				if !_rules[rulesp]() {
					goto l216
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l216
				}
				if !_rules[rulesp]() {
					goto l216
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l216
				}
				if !_rules[ruleAction6]() {
					goto l216
				}
				depth--
				add(ruleUpdateStateStmt, position217)
			}
			return true
		l216:
			position, tokenIndex, depth = position216, tokenIndex216, depth216
			return false
		},
		/* 13 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action7)> */
		func() bool {
			position240, tokenIndex240, depth240 := position, tokenIndex, depth
			{
				position241 := position
				depth++
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('U') {
						goto l240
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('P') {
						goto l240
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('D') {
						goto l240
					}
					position++
				}
			l246:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('A') {
						goto l240
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('T') {
						goto l240
					}
					position++
				}
			l250:
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if buffer[position] != rune('E') {
						goto l240
					}
					position++
				}
			l252:
				if !_rules[rulesp]() {
					goto l240
				}
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if buffer[position] != rune('S') {
						goto l240
					}
					position++
				}
			l254:
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l257
					}
					position++
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if buffer[position] != rune('O') {
						goto l240
					}
					position++
				}
			l256:
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
						goto l240
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('R') {
						goto l240
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('C') {
						goto l240
					}
					position++
				}
			l262:
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('E') {
						goto l240
					}
					position++
				}
			l264:
				if !_rules[rulesp]() {
					goto l240
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l240
				}
				if !_rules[rulesp]() {
					goto l240
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l240
				}
				if !_rules[ruleAction7]() {
					goto l240
				}
				depth--
				add(ruleUpdateSourceStmt, position241)
			}
			return true
		l240:
			position, tokenIndex, depth = position240, tokenIndex240, depth240
			return false
		},
		/* 14 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp UpdateSourceSinkSpecs Action8)> */
		func() bool {
			position266, tokenIndex266, depth266 := position, tokenIndex, depth
			{
				position267 := position
				depth++
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('U') {
						goto l266
					}
					position++
				}
			l268:
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if buffer[position] != rune('P') {
						goto l266
					}
					position++
				}
			l270:
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l273
					}
					position++
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if buffer[position] != rune('D') {
						goto l266
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
						goto l266
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
						goto l266
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
						goto l266
					}
					position++
				}
			l278:
				if !_rules[rulesp]() {
					goto l266
				}
				{
					position280, tokenIndex280, depth280 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l281
					}
					position++
					goto l280
				l281:
					position, tokenIndex, depth = position280, tokenIndex280, depth280
					if buffer[position] != rune('S') {
						goto l266
					}
					position++
				}
			l280:
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('I') {
						goto l266
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
					if buffer[position] != rune('N') {
						goto l266
					}
					position++
				}
			l284:
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l287
					}
					position++
					goto l286
				l287:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
					if buffer[position] != rune('K') {
						goto l266
					}
					position++
				}
			l286:
				if !_rules[rulesp]() {
					goto l266
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l266
				}
				if !_rules[rulesp]() {
					goto l266
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l266
				}
				if !_rules[ruleAction8]() {
					goto l266
				}
				depth--
				add(ruleUpdateSinkStmt, position267)
			}
			return true
		l266:
			position, tokenIndex, depth = position266, tokenIndex266, depth266
			return false
		},
		/* 15 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action9)> */
		func() bool {
			position288, tokenIndex288, depth288 := position, tokenIndex, depth
			{
				position289 := position
				depth++
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l291
					}
					position++
					goto l290
				l291:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
					if buffer[position] != rune('I') {
						goto l288
					}
					position++
				}
			l290:
				{
					position292, tokenIndex292, depth292 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l293
					}
					position++
					goto l292
				l293:
					position, tokenIndex, depth = position292, tokenIndex292, depth292
					if buffer[position] != rune('N') {
						goto l288
					}
					position++
				}
			l292:
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
						goto l288
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if buffer[position] != rune('E') {
						goto l288
					}
					position++
				}
			l296:
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					if buffer[position] != rune('R') {
						goto l288
					}
					position++
				}
			l298:
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l301
					}
					position++
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if buffer[position] != rune('T') {
						goto l288
					}
					position++
				}
			l300:
				if !_rules[rulesp]() {
					goto l288
				}
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l303
					}
					position++
					goto l302
				l303:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					if buffer[position] != rune('I') {
						goto l288
					}
					position++
				}
			l302:
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
					if buffer[position] != rune('N') {
						goto l288
					}
					position++
				}
			l304:
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l307
					}
					position++
					goto l306
				l307:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
					if buffer[position] != rune('T') {
						goto l288
					}
					position++
				}
			l306:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					if buffer[position] != rune('O') {
						goto l288
					}
					position++
				}
			l308:
				if !_rules[rulesp]() {
					goto l288
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l288
				}
				if !_rules[rulesp]() {
					goto l288
				}
				if !_rules[ruleSelectStmt]() {
					goto l288
				}
				if !_rules[ruleAction9]() {
					goto l288
				}
				depth--
				add(ruleInsertIntoSelectStmt, position289)
			}
			return true
		l288:
			position, tokenIndex, depth = position288, tokenIndex288, depth288
			return false
		},
		/* 16 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action10)> */
		func() bool {
			position310, tokenIndex310, depth310 := position, tokenIndex, depth
			{
				position311 := position
				depth++
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('I') {
						goto l310
					}
					position++
				}
			l312:
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l315
					}
					position++
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if buffer[position] != rune('N') {
						goto l310
					}
					position++
				}
			l314:
				{
					position316, tokenIndex316, depth316 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l317
					}
					position++
					goto l316
				l317:
					position, tokenIndex, depth = position316, tokenIndex316, depth316
					if buffer[position] != rune('S') {
						goto l310
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
						goto l310
					}
					position++
				}
			l318:
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('R') {
						goto l310
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if buffer[position] != rune('T') {
						goto l310
					}
					position++
				}
			l322:
				if !_rules[rulesp]() {
					goto l310
				}
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l325
					}
					position++
					goto l324
				l325:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if buffer[position] != rune('I') {
						goto l310
					}
					position++
				}
			l324:
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l327
					}
					position++
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					if buffer[position] != rune('N') {
						goto l310
					}
					position++
				}
			l326:
				{
					position328, tokenIndex328, depth328 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l329
					}
					position++
					goto l328
				l329:
					position, tokenIndex, depth = position328, tokenIndex328, depth328
					if buffer[position] != rune('T') {
						goto l310
					}
					position++
				}
			l328:
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l331
					}
					position++
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					if buffer[position] != rune('O') {
						goto l310
					}
					position++
				}
			l330:
				if !_rules[rulesp]() {
					goto l310
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l310
				}
				if !_rules[rulesp]() {
					goto l310
				}
				{
					position332, tokenIndex332, depth332 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l333
					}
					position++
					goto l332
				l333:
					position, tokenIndex, depth = position332, tokenIndex332, depth332
					if buffer[position] != rune('F') {
						goto l310
					}
					position++
				}
			l332:
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if buffer[position] != rune('R') {
						goto l310
					}
					position++
				}
			l334:
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if buffer[position] != rune('O') {
						goto l310
					}
					position++
				}
			l336:
				{
					position338, tokenIndex338, depth338 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l339
					}
					position++
					goto l338
				l339:
					position, tokenIndex, depth = position338, tokenIndex338, depth338
					if buffer[position] != rune('M') {
						goto l310
					}
					position++
				}
			l338:
				if !_rules[rulesp]() {
					goto l310
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l310
				}
				if !_rules[ruleAction10]() {
					goto l310
				}
				depth--
				add(ruleInsertIntoFromStmt, position311)
			}
			return true
		l310:
			position, tokenIndex, depth = position310, tokenIndex310, depth310
			return false
		},
		/* 17 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action11)> */
		func() bool {
			position340, tokenIndex340, depth340 := position, tokenIndex, depth
			{
				position341 := position
				depth++
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if buffer[position] != rune('P') {
						goto l340
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
					if buffer[position] != rune('A') {
						goto l340
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					if buffer[position] != rune('U') {
						goto l340
					}
					position++
				}
			l346:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if buffer[position] != rune('S') {
						goto l340
					}
					position++
				}
			l348:
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if buffer[position] != rune('E') {
						goto l340
					}
					position++
				}
			l350:
				if !_rules[rulesp]() {
					goto l340
				}
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('S') {
						goto l340
					}
					position++
				}
			l352:
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
					if buffer[position] != rune('O') {
						goto l340
					}
					position++
				}
			l354:
				{
					position356, tokenIndex356, depth356 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l357
					}
					position++
					goto l356
				l357:
					position, tokenIndex, depth = position356, tokenIndex356, depth356
					if buffer[position] != rune('U') {
						goto l340
					}
					position++
				}
			l356:
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l359
					}
					position++
					goto l358
				l359:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					if buffer[position] != rune('R') {
						goto l340
					}
					position++
				}
			l358:
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if buffer[position] != rune('C') {
						goto l340
					}
					position++
				}
			l360:
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l363
					}
					position++
					goto l362
				l363:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
					if buffer[position] != rune('E') {
						goto l340
					}
					position++
				}
			l362:
				if !_rules[rulesp]() {
					goto l340
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l340
				}
				if !_rules[ruleAction11]() {
					goto l340
				}
				depth--
				add(rulePauseSourceStmt, position341)
			}
			return true
		l340:
			position, tokenIndex, depth = position340, tokenIndex340, depth340
			return false
		},
		/* 18 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action12)> */
		func() bool {
			position364, tokenIndex364, depth364 := position, tokenIndex, depth
			{
				position365 := position
				depth++
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l367
					}
					position++
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if buffer[position] != rune('R') {
						goto l364
					}
					position++
				}
			l366:
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if buffer[position] != rune('E') {
						goto l364
					}
					position++
				}
			l368:
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l371
					}
					position++
					goto l370
				l371:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
					if buffer[position] != rune('S') {
						goto l364
					}
					position++
				}
			l370:
				{
					position372, tokenIndex372, depth372 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l373
					}
					position++
					goto l372
				l373:
					position, tokenIndex, depth = position372, tokenIndex372, depth372
					if buffer[position] != rune('U') {
						goto l364
					}
					position++
				}
			l372:
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('M') {
						goto l364
					}
					position++
				}
			l374:
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l377
					}
					position++
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					if buffer[position] != rune('E') {
						goto l364
					}
					position++
				}
			l376:
				if !_rules[rulesp]() {
					goto l364
				}
				{
					position378, tokenIndex378, depth378 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l379
					}
					position++
					goto l378
				l379:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
					if buffer[position] != rune('S') {
						goto l364
					}
					position++
				}
			l378:
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('O') {
						goto l364
					}
					position++
				}
			l380:
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l383
					}
					position++
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if buffer[position] != rune('U') {
						goto l364
					}
					position++
				}
			l382:
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l385
					}
					position++
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if buffer[position] != rune('R') {
						goto l364
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('C') {
						goto l364
					}
					position++
				}
			l386:
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l389
					}
					position++
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					if buffer[position] != rune('E') {
						goto l364
					}
					position++
				}
			l388:
				if !_rules[rulesp]() {
					goto l364
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l364
				}
				if !_rules[ruleAction12]() {
					goto l364
				}
				depth--
				add(ruleResumeSourceStmt, position365)
			}
			return true
		l364:
			position, tokenIndex, depth = position364, tokenIndex364, depth364
			return false
		},
		/* 19 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if buffer[position] != rune('R') {
						goto l390
					}
					position++
				}
			l392:
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l395
					}
					position++
					goto l394
				l395:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
					if buffer[position] != rune('E') {
						goto l390
					}
					position++
				}
			l394:
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('W') {
						goto l390
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if buffer[position] != rune('I') {
						goto l390
					}
					position++
				}
			l398:
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l401
					}
					position++
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					if buffer[position] != rune('N') {
						goto l390
					}
					position++
				}
			l400:
				{
					position402, tokenIndex402, depth402 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l403
					}
					position++
					goto l402
				l403:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
					if buffer[position] != rune('D') {
						goto l390
					}
					position++
				}
			l402:
				if !_rules[rulesp]() {
					goto l390
				}
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l405
					}
					position++
					goto l404
				l405:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					if buffer[position] != rune('S') {
						goto l390
					}
					position++
				}
			l404:
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l407
					}
					position++
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if buffer[position] != rune('O') {
						goto l390
					}
					position++
				}
			l406:
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l409
					}
					position++
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					if buffer[position] != rune('U') {
						goto l390
					}
					position++
				}
			l408:
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l411
					}
					position++
					goto l410
				l411:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
					if buffer[position] != rune('R') {
						goto l390
					}
					position++
				}
			l410:
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l413
					}
					position++
					goto l412
				l413:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
					if buffer[position] != rune('C') {
						goto l390
					}
					position++
				}
			l412:
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('E') {
						goto l390
					}
					position++
				}
			l414:
				if !_rules[rulesp]() {
					goto l390
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l390
				}
				if !_rules[ruleAction13]() {
					goto l390
				}
				depth--
				add(ruleRewindSourceStmt, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 20 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position416, tokenIndex416, depth416 := position, tokenIndex, depth
			{
				position417 := position
				depth++
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
					if buffer[position] != rune('D') {
						goto l416
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('R') {
						goto l416
					}
					position++
				}
			l420:
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l423
					}
					position++
					goto l422
				l423:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
					if buffer[position] != rune('O') {
						goto l416
					}
					position++
				}
			l422:
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l425
					}
					position++
					goto l424
				l425:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
					if buffer[position] != rune('P') {
						goto l416
					}
					position++
				}
			l424:
				if !_rules[rulesp]() {
					goto l416
				}
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l427
					}
					position++
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if buffer[position] != rune('S') {
						goto l416
					}
					position++
				}
			l426:
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l429
					}
					position++
					goto l428
				l429:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
					if buffer[position] != rune('O') {
						goto l416
					}
					position++
				}
			l428:
				{
					position430, tokenIndex430, depth430 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l431
					}
					position++
					goto l430
				l431:
					position, tokenIndex, depth = position430, tokenIndex430, depth430
					if buffer[position] != rune('U') {
						goto l416
					}
					position++
				}
			l430:
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
						goto l416
					}
					position++
				}
			l432:
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l435
					}
					position++
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if buffer[position] != rune('C') {
						goto l416
					}
					position++
				}
			l434:
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l437
					}
					position++
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if buffer[position] != rune('E') {
						goto l416
					}
					position++
				}
			l436:
				if !_rules[rulesp]() {
					goto l416
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l416
				}
				if !_rules[ruleAction14]() {
					goto l416
				}
				depth--
				add(ruleDropSourceStmt, position417)
			}
			return true
		l416:
			position, tokenIndex, depth = position416, tokenIndex416, depth416
			return false
		},
		/* 21 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action15)> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				{
					position440, tokenIndex440, depth440 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l441
					}
					position++
					goto l440
				l441:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if buffer[position] != rune('D') {
						goto l438
					}
					position++
				}
			l440:
				{
					position442, tokenIndex442, depth442 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l443
					}
					position++
					goto l442
				l443:
					position, tokenIndex, depth = position442, tokenIndex442, depth442
					if buffer[position] != rune('R') {
						goto l438
					}
					position++
				}
			l442:
				{
					position444, tokenIndex444, depth444 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l445
					}
					position++
					goto l444
				l445:
					position, tokenIndex, depth = position444, tokenIndex444, depth444
					if buffer[position] != rune('O') {
						goto l438
					}
					position++
				}
			l444:
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
					if buffer[position] != rune('P') {
						goto l438
					}
					position++
				}
			l446:
				if !_rules[rulesp]() {
					goto l438
				}
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
					if buffer[position] != rune('S') {
						goto l438
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('T') {
						goto l438
					}
					position++
				}
			l450:
				{
					position452, tokenIndex452, depth452 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l453
					}
					position++
					goto l452
				l453:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if buffer[position] != rune('R') {
						goto l438
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
						goto l438
					}
					position++
				}
			l454:
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l457
					}
					position++
					goto l456
				l457:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if buffer[position] != rune('A') {
						goto l438
					}
					position++
				}
			l456:
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l459
					}
					position++
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if buffer[position] != rune('M') {
						goto l438
					}
					position++
				}
			l458:
				if !_rules[rulesp]() {
					goto l438
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l438
				}
				if !_rules[ruleAction15]() {
					goto l438
				}
				depth--
				add(ruleDropStreamStmt, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 22 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action16)> */
		func() bool {
			position460, tokenIndex460, depth460 := position, tokenIndex, depth
			{
				position461 := position
				depth++
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l463
					}
					position++
					goto l462
				l463:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if buffer[position] != rune('D') {
						goto l460
					}
					position++
				}
			l462:
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l465
					}
					position++
					goto l464
				l465:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					if buffer[position] != rune('R') {
						goto l460
					}
					position++
				}
			l464:
				{
					position466, tokenIndex466, depth466 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l467
					}
					position++
					goto l466
				l467:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if buffer[position] != rune('O') {
						goto l460
					}
					position++
				}
			l466:
				{
					position468, tokenIndex468, depth468 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l469
					}
					position++
					goto l468
				l469:
					position, tokenIndex, depth = position468, tokenIndex468, depth468
					if buffer[position] != rune('P') {
						goto l460
					}
					position++
				}
			l468:
				if !_rules[rulesp]() {
					goto l460
				}
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l471
					}
					position++
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if buffer[position] != rune('S') {
						goto l460
					}
					position++
				}
			l470:
				{
					position472, tokenIndex472, depth472 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l473
					}
					position++
					goto l472
				l473:
					position, tokenIndex, depth = position472, tokenIndex472, depth472
					if buffer[position] != rune('I') {
						goto l460
					}
					position++
				}
			l472:
				{
					position474, tokenIndex474, depth474 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l475
					}
					position++
					goto l474
				l475:
					position, tokenIndex, depth = position474, tokenIndex474, depth474
					if buffer[position] != rune('N') {
						goto l460
					}
					position++
				}
			l474:
				{
					position476, tokenIndex476, depth476 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l477
					}
					position++
					goto l476
				l477:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if buffer[position] != rune('K') {
						goto l460
					}
					position++
				}
			l476:
				if !_rules[rulesp]() {
					goto l460
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l460
				}
				if !_rules[ruleAction16]() {
					goto l460
				}
				depth--
				add(ruleDropSinkStmt, position461)
			}
			return true
		l460:
			position, tokenIndex, depth = position460, tokenIndex460, depth460
			return false
		},
		/* 23 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action17)> */
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
					if buffer[position] != rune('a') {
						goto l493
					}
					position++
					goto l492
				l493:
					position, tokenIndex, depth = position492, tokenIndex492, depth492
					if buffer[position] != rune('A') {
						goto l478
					}
					position++
				}
			l492:
				{
					position494, tokenIndex494, depth494 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l495
					}
					position++
					goto l494
				l495:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if buffer[position] != rune('T') {
						goto l478
					}
					position++
				}
			l494:
				{
					position496, tokenIndex496, depth496 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l497
					}
					position++
					goto l496
				l497:
					position, tokenIndex, depth = position496, tokenIndex496, depth496
					if buffer[position] != rune('E') {
						goto l478
					}
					position++
				}
			l496:
				if !_rules[rulesp]() {
					goto l478
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l478
				}
				if !_rules[ruleAction17]() {
					goto l478
				}
				depth--
				add(ruleDropStateStmt, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 24 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) Action18)> */
		func() bool {
			position498, tokenIndex498, depth498 := position, tokenIndex, depth
			{
				position499 := position
				depth++
				{
					position500, tokenIndex500, depth500 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l501
					}
					goto l500
				l501:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if !_rules[ruleDSTREAM]() {
						goto l502
					}
					goto l500
				l502:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if !_rules[ruleRSTREAM]() {
						goto l498
					}
				}
			l500:
				if !_rules[ruleAction18]() {
					goto l498
				}
				depth--
				add(ruleEmitter, position499)
			}
			return true
		l498:
			position, tokenIndex, depth = position498, tokenIndex498, depth498
			return false
		},
		/* 25 Projections <- <(<(Projection sp (',' sp Projection)*)> Action19)> */
		func() bool {
			position503, tokenIndex503, depth503 := position, tokenIndex, depth
			{
				position504 := position
				depth++
				{
					position505 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l503
					}
					if !_rules[rulesp]() {
						goto l503
					}
				l506:
					{
						position507, tokenIndex507, depth507 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l507
						}
						position++
						if !_rules[rulesp]() {
							goto l507
						}
						if !_rules[ruleProjection]() {
							goto l507
						}
						goto l506
					l507:
						position, tokenIndex, depth = position507, tokenIndex507, depth507
					}
					depth--
					add(rulePegText, position505)
				}
				if !_rules[ruleAction19]() {
					goto l503
				}
				depth--
				add(ruleProjections, position504)
			}
			return true
		l503:
			position, tokenIndex, depth = position503, tokenIndex503, depth503
			return false
		},
		/* 26 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position508, tokenIndex508, depth508 := position, tokenIndex, depth
			{
				position509 := position
				depth++
				{
					position510, tokenIndex510, depth510 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l511
					}
					goto l510
				l511:
					position, tokenIndex, depth = position510, tokenIndex510, depth510
					if !_rules[ruleExpression]() {
						goto l512
					}
					goto l510
				l512:
					position, tokenIndex, depth = position510, tokenIndex510, depth510
					if !_rules[ruleWildcard]() {
						goto l508
					}
				}
			l510:
				depth--
				add(ruleProjection, position509)
			}
			return true
		l508:
			position, tokenIndex, depth = position508, tokenIndex508, depth508
			return false
		},
		/* 27 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action20)> */
		func() bool {
			position513, tokenIndex513, depth513 := position, tokenIndex, depth
			{
				position514 := position
				depth++
				{
					position515, tokenIndex515, depth515 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l516
					}
					goto l515
				l516:
					position, tokenIndex, depth = position515, tokenIndex515, depth515
					if !_rules[ruleWildcard]() {
						goto l513
					}
				}
			l515:
				if !_rules[rulesp]() {
					goto l513
				}
				{
					position517, tokenIndex517, depth517 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l518
					}
					position++
					goto l517
				l518:
					position, tokenIndex, depth = position517, tokenIndex517, depth517
					if buffer[position] != rune('A') {
						goto l513
					}
					position++
				}
			l517:
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
						goto l513
					}
					position++
				}
			l519:
				if !_rules[rulesp]() {
					goto l513
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l513
				}
				if !_rules[ruleAction20]() {
					goto l513
				}
				depth--
				add(ruleAliasExpression, position514)
			}
			return true
		l513:
			position, tokenIndex, depth = position513, tokenIndex513, depth513
			return false
		},
		/* 28 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action21)> */
		func() bool {
			position521, tokenIndex521, depth521 := position, tokenIndex, depth
			{
				position522 := position
				depth++
				{
					position523 := position
					depth++
					{
						position524, tokenIndex524, depth524 := position, tokenIndex, depth
						{
							position526, tokenIndex526, depth526 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l527
							}
							position++
							goto l526
						l527:
							position, tokenIndex, depth = position526, tokenIndex526, depth526
							if buffer[position] != rune('F') {
								goto l524
							}
							position++
						}
					l526:
						{
							position528, tokenIndex528, depth528 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l529
							}
							position++
							goto l528
						l529:
							position, tokenIndex, depth = position528, tokenIndex528, depth528
							if buffer[position] != rune('R') {
								goto l524
							}
							position++
						}
					l528:
						{
							position530, tokenIndex530, depth530 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l531
							}
							position++
							goto l530
						l531:
							position, tokenIndex, depth = position530, tokenIndex530, depth530
							if buffer[position] != rune('O') {
								goto l524
							}
							position++
						}
					l530:
						{
							position532, tokenIndex532, depth532 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l533
							}
							position++
							goto l532
						l533:
							position, tokenIndex, depth = position532, tokenIndex532, depth532
							if buffer[position] != rune('M') {
								goto l524
							}
							position++
						}
					l532:
						if !_rules[rulesp]() {
							goto l524
						}
						if !_rules[ruleRelations]() {
							goto l524
						}
						if !_rules[rulesp]() {
							goto l524
						}
						goto l525
					l524:
						position, tokenIndex, depth = position524, tokenIndex524, depth524
					}
				l525:
					depth--
					add(rulePegText, position523)
				}
				if !_rules[ruleAction21]() {
					goto l521
				}
				depth--
				add(ruleWindowedFrom, position522)
			}
			return true
		l521:
			position, tokenIndex, depth = position521, tokenIndex521, depth521
			return false
		},
		/* 29 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position534, tokenIndex534, depth534 := position, tokenIndex, depth
			{
				position535 := position
				depth++
				{
					position536, tokenIndex536, depth536 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l537
					}
					goto l536
				l537:
					position, tokenIndex, depth = position536, tokenIndex536, depth536
					if !_rules[ruleTuplesInterval]() {
						goto l534
					}
				}
			l536:
				depth--
				add(ruleInterval, position535)
			}
			return true
		l534:
			position, tokenIndex, depth = position534, tokenIndex534, depth534
			return false
		},
		/* 30 TimeInterval <- <(NumericLiteral sp SECONDS Action22)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l538
				}
				if !_rules[rulesp]() {
					goto l538
				}
				if !_rules[ruleSECONDS]() {
					goto l538
				}
				if !_rules[ruleAction22]() {
					goto l538
				}
				depth--
				add(ruleTimeInterval, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 31 TuplesInterval <- <(NumericLiteral sp TUPLES Action23)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l540
				}
				if !_rules[rulesp]() {
					goto l540
				}
				if !_rules[ruleTUPLES]() {
					goto l540
				}
				if !_rules[ruleAction23]() {
					goto l540
				}
				depth--
				add(ruleTuplesInterval, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 32 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position542, tokenIndex542, depth542 := position, tokenIndex, depth
			{
				position543 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l542
				}
				if !_rules[rulesp]() {
					goto l542
				}
			l544:
				{
					position545, tokenIndex545, depth545 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l545
					}
					position++
					if !_rules[rulesp]() {
						goto l545
					}
					if !_rules[ruleRelationLike]() {
						goto l545
					}
					goto l544
				l545:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
				}
				depth--
				add(ruleRelations, position543)
			}
			return true
		l542:
			position, tokenIndex, depth = position542, tokenIndex542, depth542
			return false
		},
		/* 33 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action24)> */
		func() bool {
			position546, tokenIndex546, depth546 := position, tokenIndex, depth
			{
				position547 := position
				depth++
				{
					position548 := position
					depth++
					{
						position549, tokenIndex549, depth549 := position, tokenIndex, depth
						{
							position551, tokenIndex551, depth551 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l552
							}
							position++
							goto l551
						l552:
							position, tokenIndex, depth = position551, tokenIndex551, depth551
							if buffer[position] != rune('W') {
								goto l549
							}
							position++
						}
					l551:
						{
							position553, tokenIndex553, depth553 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l554
							}
							position++
							goto l553
						l554:
							position, tokenIndex, depth = position553, tokenIndex553, depth553
							if buffer[position] != rune('H') {
								goto l549
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
								goto l549
							}
							position++
						}
					l555:
						{
							position557, tokenIndex557, depth557 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l558
							}
							position++
							goto l557
						l558:
							position, tokenIndex, depth = position557, tokenIndex557, depth557
							if buffer[position] != rune('R') {
								goto l549
							}
							position++
						}
					l557:
						{
							position559, tokenIndex559, depth559 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l560
							}
							position++
							goto l559
						l560:
							position, tokenIndex, depth = position559, tokenIndex559, depth559
							if buffer[position] != rune('E') {
								goto l549
							}
							position++
						}
					l559:
						if !_rules[rulesp]() {
							goto l549
						}
						if !_rules[ruleExpression]() {
							goto l549
						}
						goto l550
					l549:
						position, tokenIndex, depth = position549, tokenIndex549, depth549
					}
				l550:
					depth--
					add(rulePegText, position548)
				}
				if !_rules[ruleAction24]() {
					goto l546
				}
				depth--
				add(ruleFilter, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 34 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action25)> */
		func() bool {
			position561, tokenIndex561, depth561 := position, tokenIndex, depth
			{
				position562 := position
				depth++
				{
					position563 := position
					depth++
					{
						position564, tokenIndex564, depth564 := position, tokenIndex, depth
						{
							position566, tokenIndex566, depth566 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l567
							}
							position++
							goto l566
						l567:
							position, tokenIndex, depth = position566, tokenIndex566, depth566
							if buffer[position] != rune('G') {
								goto l564
							}
							position++
						}
					l566:
						{
							position568, tokenIndex568, depth568 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l569
							}
							position++
							goto l568
						l569:
							position, tokenIndex, depth = position568, tokenIndex568, depth568
							if buffer[position] != rune('R') {
								goto l564
							}
							position++
						}
					l568:
						{
							position570, tokenIndex570, depth570 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l571
							}
							position++
							goto l570
						l571:
							position, tokenIndex, depth = position570, tokenIndex570, depth570
							if buffer[position] != rune('O') {
								goto l564
							}
							position++
						}
					l570:
						{
							position572, tokenIndex572, depth572 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l573
							}
							position++
							goto l572
						l573:
							position, tokenIndex, depth = position572, tokenIndex572, depth572
							if buffer[position] != rune('U') {
								goto l564
							}
							position++
						}
					l572:
						{
							position574, tokenIndex574, depth574 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l575
							}
							position++
							goto l574
						l575:
							position, tokenIndex, depth = position574, tokenIndex574, depth574
							if buffer[position] != rune('P') {
								goto l564
							}
							position++
						}
					l574:
						if !_rules[rulesp]() {
							goto l564
						}
						{
							position576, tokenIndex576, depth576 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l577
							}
							position++
							goto l576
						l577:
							position, tokenIndex, depth = position576, tokenIndex576, depth576
							if buffer[position] != rune('B') {
								goto l564
							}
							position++
						}
					l576:
						{
							position578, tokenIndex578, depth578 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l579
							}
							position++
							goto l578
						l579:
							position, tokenIndex, depth = position578, tokenIndex578, depth578
							if buffer[position] != rune('Y') {
								goto l564
							}
							position++
						}
					l578:
						if !_rules[rulesp]() {
							goto l564
						}
						if !_rules[ruleGroupList]() {
							goto l564
						}
						goto l565
					l564:
						position, tokenIndex, depth = position564, tokenIndex564, depth564
					}
				l565:
					depth--
					add(rulePegText, position563)
				}
				if !_rules[ruleAction25]() {
					goto l561
				}
				depth--
				add(ruleGrouping, position562)
			}
			return true
		l561:
			position, tokenIndex, depth = position561, tokenIndex561, depth561
			return false
		},
		/* 35 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position580, tokenIndex580, depth580 := position, tokenIndex, depth
			{
				position581 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l580
				}
				if !_rules[rulesp]() {
					goto l580
				}
			l582:
				{
					position583, tokenIndex583, depth583 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l583
					}
					position++
					if !_rules[rulesp]() {
						goto l583
					}
					if !_rules[ruleExpression]() {
						goto l583
					}
					goto l582
				l583:
					position, tokenIndex, depth = position583, tokenIndex583, depth583
				}
				depth--
				add(ruleGroupList, position581)
			}
			return true
		l580:
			position, tokenIndex, depth = position580, tokenIndex580, depth580
			return false
		},
		/* 36 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action26)> */
		func() bool {
			position584, tokenIndex584, depth584 := position, tokenIndex, depth
			{
				position585 := position
				depth++
				{
					position586 := position
					depth++
					{
						position587, tokenIndex587, depth587 := position, tokenIndex, depth
						{
							position589, tokenIndex589, depth589 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l590
							}
							position++
							goto l589
						l590:
							position, tokenIndex, depth = position589, tokenIndex589, depth589
							if buffer[position] != rune('H') {
								goto l587
							}
							position++
						}
					l589:
						{
							position591, tokenIndex591, depth591 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l592
							}
							position++
							goto l591
						l592:
							position, tokenIndex, depth = position591, tokenIndex591, depth591
							if buffer[position] != rune('A') {
								goto l587
							}
							position++
						}
					l591:
						{
							position593, tokenIndex593, depth593 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l594
							}
							position++
							goto l593
						l594:
							position, tokenIndex, depth = position593, tokenIndex593, depth593
							if buffer[position] != rune('V') {
								goto l587
							}
							position++
						}
					l593:
						{
							position595, tokenIndex595, depth595 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l596
							}
							position++
							goto l595
						l596:
							position, tokenIndex, depth = position595, tokenIndex595, depth595
							if buffer[position] != rune('I') {
								goto l587
							}
							position++
						}
					l595:
						{
							position597, tokenIndex597, depth597 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l598
							}
							position++
							goto l597
						l598:
							position, tokenIndex, depth = position597, tokenIndex597, depth597
							if buffer[position] != rune('N') {
								goto l587
							}
							position++
						}
					l597:
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
								goto l587
							}
							position++
						}
					l599:
						if !_rules[rulesp]() {
							goto l587
						}
						if !_rules[ruleExpression]() {
							goto l587
						}
						goto l588
					l587:
						position, tokenIndex, depth = position587, tokenIndex587, depth587
					}
				l588:
					depth--
					add(rulePegText, position586)
				}
				if !_rules[ruleAction26]() {
					goto l584
				}
				depth--
				add(ruleHaving, position585)
			}
			return true
		l584:
			position, tokenIndex, depth = position584, tokenIndex584, depth584
			return false
		},
		/* 37 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action27))> */
		func() bool {
			position601, tokenIndex601, depth601 := position, tokenIndex, depth
			{
				position602 := position
				depth++
				{
					position603, tokenIndex603, depth603 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l604
					}
					goto l603
				l604:
					position, tokenIndex, depth = position603, tokenIndex603, depth603
					if !_rules[ruleStreamWindow]() {
						goto l601
					}
					if !_rules[ruleAction27]() {
						goto l601
					}
				}
			l603:
				depth--
				add(ruleRelationLike, position602)
			}
			return true
		l601:
			position, tokenIndex, depth = position601, tokenIndex601, depth601
			return false
		},
		/* 38 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action28)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l605
				}
				if !_rules[rulesp]() {
					goto l605
				}
				{
					position607, tokenIndex607, depth607 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l608
					}
					position++
					goto l607
				l608:
					position, tokenIndex, depth = position607, tokenIndex607, depth607
					if buffer[position] != rune('A') {
						goto l605
					}
					position++
				}
			l607:
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
						goto l605
					}
					position++
				}
			l609:
				if !_rules[rulesp]() {
					goto l605
				}
				if !_rules[ruleIdentifier]() {
					goto l605
				}
				if !_rules[ruleAction28]() {
					goto l605
				}
				depth--
				add(ruleAliasedStreamWindow, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 39 StreamWindow <- <(StreamLike sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action29)> */
		func() bool {
			position611, tokenIndex611, depth611 := position, tokenIndex, depth
			{
				position612 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l611
				}
				if !_rules[rulesp]() {
					goto l611
				}
				if buffer[position] != rune('[') {
					goto l611
				}
				position++
				if !_rules[rulesp]() {
					goto l611
				}
				{
					position613, tokenIndex613, depth613 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l614
					}
					position++
					goto l613
				l614:
					position, tokenIndex, depth = position613, tokenIndex613, depth613
					if buffer[position] != rune('R') {
						goto l611
					}
					position++
				}
			l613:
				{
					position615, tokenIndex615, depth615 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l616
					}
					position++
					goto l615
				l616:
					position, tokenIndex, depth = position615, tokenIndex615, depth615
					if buffer[position] != rune('A') {
						goto l611
					}
					position++
				}
			l615:
				{
					position617, tokenIndex617, depth617 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l618
					}
					position++
					goto l617
				l618:
					position, tokenIndex, depth = position617, tokenIndex617, depth617
					if buffer[position] != rune('N') {
						goto l611
					}
					position++
				}
			l617:
				{
					position619, tokenIndex619, depth619 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l620
					}
					position++
					goto l619
				l620:
					position, tokenIndex, depth = position619, tokenIndex619, depth619
					if buffer[position] != rune('G') {
						goto l611
					}
					position++
				}
			l619:
				{
					position621, tokenIndex621, depth621 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l622
					}
					position++
					goto l621
				l622:
					position, tokenIndex, depth = position621, tokenIndex621, depth621
					if buffer[position] != rune('E') {
						goto l611
					}
					position++
				}
			l621:
				if !_rules[rulesp]() {
					goto l611
				}
				if !_rules[ruleInterval]() {
					goto l611
				}
				if !_rules[rulesp]() {
					goto l611
				}
				if buffer[position] != rune(']') {
					goto l611
				}
				position++
				if !_rules[ruleAction29]() {
					goto l611
				}
				depth--
				add(ruleStreamWindow, position612)
			}
			return true
		l611:
			position, tokenIndex, depth = position611, tokenIndex611, depth611
			return false
		},
		/* 40 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position623, tokenIndex623, depth623 := position, tokenIndex, depth
			{
				position624 := position
				depth++
				{
					position625, tokenIndex625, depth625 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l626
					}
					goto l625
				l626:
					position, tokenIndex, depth = position625, tokenIndex625, depth625
					if !_rules[ruleStream]() {
						goto l623
					}
				}
			l625:
				depth--
				add(ruleStreamLike, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 41 UDSFFuncApp <- <(FuncApp Action30)> */
		func() bool {
			position627, tokenIndex627, depth627 := position, tokenIndex, depth
			{
				position628 := position
				depth++
				if !_rules[ruleFuncApp]() {
					goto l627
				}
				if !_rules[ruleAction30]() {
					goto l627
				}
				depth--
				add(ruleUDSFFuncApp, position628)
			}
			return true
		l627:
			position, tokenIndex, depth = position627, tokenIndex627, depth627
			return false
		},
		/* 42 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action31)> */
		func() bool {
			position629, tokenIndex629, depth629 := position, tokenIndex, depth
			{
				position630 := position
				depth++
				{
					position631 := position
					depth++
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						{
							position634, tokenIndex634, depth634 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l635
							}
							position++
							goto l634
						l635:
							position, tokenIndex, depth = position634, tokenIndex634, depth634
							if buffer[position] != rune('W') {
								goto l632
							}
							position++
						}
					l634:
						{
							position636, tokenIndex636, depth636 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l637
							}
							position++
							goto l636
						l637:
							position, tokenIndex, depth = position636, tokenIndex636, depth636
							if buffer[position] != rune('I') {
								goto l632
							}
							position++
						}
					l636:
						{
							position638, tokenIndex638, depth638 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l639
							}
							position++
							goto l638
						l639:
							position, tokenIndex, depth = position638, tokenIndex638, depth638
							if buffer[position] != rune('T') {
								goto l632
							}
							position++
						}
					l638:
						{
							position640, tokenIndex640, depth640 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l641
							}
							position++
							goto l640
						l641:
							position, tokenIndex, depth = position640, tokenIndex640, depth640
							if buffer[position] != rune('H') {
								goto l632
							}
							position++
						}
					l640:
						if !_rules[rulesp]() {
							goto l632
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l632
						}
						if !_rules[rulesp]() {
							goto l632
						}
					l642:
						{
							position643, tokenIndex643, depth643 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l643
							}
							position++
							if !_rules[rulesp]() {
								goto l643
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l643
							}
							goto l642
						l643:
							position, tokenIndex, depth = position643, tokenIndex643, depth643
						}
						goto l633
					l632:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
					}
				l633:
					depth--
					add(rulePegText, position631)
				}
				if !_rules[ruleAction31]() {
					goto l629
				}
				depth--
				add(ruleSourceSinkSpecs, position630)
			}
			return true
		l629:
			position, tokenIndex, depth = position629, tokenIndex629, depth629
			return false
		},
		/* 43 UpdateSourceSinkSpecs <- <(<(('s' / 'S') ('e' / 'E') ('t' / 'T') sp SourceSinkParam sp (',' sp SourceSinkParam)*)> Action32)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				{
					position646 := position
					depth++
					{
						position647, tokenIndex647, depth647 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l648
						}
						position++
						goto l647
					l648:
						position, tokenIndex, depth = position647, tokenIndex647, depth647
						if buffer[position] != rune('S') {
							goto l644
						}
						position++
					}
				l647:
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l650
						}
						position++
						goto l649
					l650:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						if buffer[position] != rune('E') {
							goto l644
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
							goto l644
						}
						position++
					}
				l651:
					if !_rules[rulesp]() {
						goto l644
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l644
					}
					if !_rules[rulesp]() {
						goto l644
					}
				l653:
					{
						position654, tokenIndex654, depth654 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l654
						}
						position++
						if !_rules[rulesp]() {
							goto l654
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l654
						}
						goto l653
					l654:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
					}
					depth--
					add(rulePegText, position646)
				}
				if !_rules[ruleAction32]() {
					goto l644
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 44 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action33)> */
		func() bool {
			position655, tokenIndex655, depth655 := position, tokenIndex, depth
			{
				position656 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l655
				}
				if buffer[position] != rune('=') {
					goto l655
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l655
				}
				if !_rules[ruleAction33]() {
					goto l655
				}
				depth--
				add(ruleSourceSinkParam, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 45 SourceSinkParamVal <- <(BooleanLiteral / Literal)> */
		func() bool {
			position657, tokenIndex657, depth657 := position, tokenIndex, depth
			{
				position658 := position
				depth++
				{
					position659, tokenIndex659, depth659 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l660
					}
					goto l659
				l660:
					position, tokenIndex, depth = position659, tokenIndex659, depth659
					if !_rules[ruleLiteral]() {
						goto l657
					}
				}
			l659:
				depth--
				add(ruleSourceSinkParamVal, position658)
			}
			return true
		l657:
			position, tokenIndex, depth = position657, tokenIndex657, depth657
			return false
		},
		/* 46 PausedOpt <- <(<(Paused / Unpaused)?> Action34)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				{
					position663 := position
					depth++
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						{
							position666, tokenIndex666, depth666 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l667
							}
							goto l666
						l667:
							position, tokenIndex, depth = position666, tokenIndex666, depth666
							if !_rules[ruleUnpaused]() {
								goto l664
							}
						}
					l666:
						goto l665
					l664:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
					}
				l665:
					depth--
					add(rulePegText, position663)
				}
				if !_rules[ruleAction34]() {
					goto l661
				}
				depth--
				add(rulePausedOpt, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 47 Expression <- <orExpr> */
		func() bool {
			position668, tokenIndex668, depth668 := position, tokenIndex, depth
			{
				position669 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l668
				}
				depth--
				add(ruleExpression, position669)
			}
			return true
		l668:
			position, tokenIndex, depth = position668, tokenIndex668, depth668
			return false
		},
		/* 48 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action35)> */
		func() bool {
			position670, tokenIndex670, depth670 := position, tokenIndex, depth
			{
				position671 := position
				depth++
				{
					position672 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l670
					}
					if !_rules[rulesp]() {
						goto l670
					}
					{
						position673, tokenIndex673, depth673 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l673
						}
						if !_rules[rulesp]() {
							goto l673
						}
						if !_rules[ruleandExpr]() {
							goto l673
						}
						goto l674
					l673:
						position, tokenIndex, depth = position673, tokenIndex673, depth673
					}
				l674:
					depth--
					add(rulePegText, position672)
				}
				if !_rules[ruleAction35]() {
					goto l670
				}
				depth--
				add(ruleorExpr, position671)
			}
			return true
		l670:
			position, tokenIndex, depth = position670, tokenIndex670, depth670
			return false
		},
		/* 49 andExpr <- <(<(notExpr sp (And sp notExpr)?)> Action36)> */
		func() bool {
			position675, tokenIndex675, depth675 := position, tokenIndex, depth
			{
				position676 := position
				depth++
				{
					position677 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l675
					}
					if !_rules[rulesp]() {
						goto l675
					}
					{
						position678, tokenIndex678, depth678 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l678
						}
						if !_rules[rulesp]() {
							goto l678
						}
						if !_rules[rulenotExpr]() {
							goto l678
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
				add(ruleandExpr, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 50 notExpr <- <(<((Not sp)? comparisonExpr)> Action37)> */
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
						if !_rules[ruleNot]() {
							goto l683
						}
						if !_rules[rulesp]() {
							goto l683
						}
						goto l684
					l683:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
					}
				l684:
					if !_rules[rulecomparisonExpr]() {
						goto l680
					}
					depth--
					add(rulePegText, position682)
				}
				if !_rules[ruleAction37]() {
					goto l680
				}
				depth--
				add(rulenotExpr, position681)
			}
			return true
		l680:
			position, tokenIndex, depth = position680, tokenIndex680, depth680
			return false
		},
		/* 51 comparisonExpr <- <(<(otherOpExpr sp (ComparisonOp sp otherOpExpr)?)> Action38)> */
		func() bool {
			position685, tokenIndex685, depth685 := position, tokenIndex, depth
			{
				position686 := position
				depth++
				{
					position687 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l685
					}
					if !_rules[rulesp]() {
						goto l685
					}
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l688
						}
						if !_rules[rulesp]() {
							goto l688
						}
						if !_rules[ruleotherOpExpr]() {
							goto l688
						}
						goto l689
					l688:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
					}
				l689:
					depth--
					add(rulePegText, position687)
				}
				if !_rules[ruleAction38]() {
					goto l685
				}
				depth--
				add(rulecomparisonExpr, position686)
			}
			return true
		l685:
			position, tokenIndex, depth = position685, tokenIndex685, depth685
			return false
		},
		/* 52 otherOpExpr <- <(<(isExpr sp (OtherOp sp isExpr sp)*)> Action39)> */
		func() bool {
			position690, tokenIndex690, depth690 := position, tokenIndex, depth
			{
				position691 := position
				depth++
				{
					position692 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l690
					}
					if !_rules[rulesp]() {
						goto l690
					}
				l693:
					{
						position694, tokenIndex694, depth694 := position, tokenIndex, depth
						if !_rules[ruleOtherOp]() {
							goto l694
						}
						if !_rules[rulesp]() {
							goto l694
						}
						if !_rules[ruleisExpr]() {
							goto l694
						}
						if !_rules[rulesp]() {
							goto l694
						}
						goto l693
					l694:
						position, tokenIndex, depth = position694, tokenIndex694, depth694
					}
					depth--
					add(rulePegText, position692)
				}
				if !_rules[ruleAction39]() {
					goto l690
				}
				depth--
				add(ruleotherOpExpr, position691)
			}
			return true
		l690:
			position, tokenIndex, depth = position690, tokenIndex690, depth690
			return false
		},
		/* 53 isExpr <- <(<(termExpr sp (IsOp sp NullLiteral)?)> Action40)> */
		func() bool {
			position695, tokenIndex695, depth695 := position, tokenIndex, depth
			{
				position696 := position
				depth++
				{
					position697 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l695
					}
					if !_rules[rulesp]() {
						goto l695
					}
					{
						position698, tokenIndex698, depth698 := position, tokenIndex, depth
						if !_rules[ruleIsOp]() {
							goto l698
						}
						if !_rules[rulesp]() {
							goto l698
						}
						if !_rules[ruleNullLiteral]() {
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
				if !_rules[ruleAction40]() {
					goto l695
				}
				depth--
				add(ruleisExpr, position696)
			}
			return true
		l695:
			position, tokenIndex, depth = position695, tokenIndex695, depth695
			return false
		},
		/* 54 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr sp)*)> Action41)> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				{
					position702 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l700
					}
					if !_rules[rulesp]() {
						goto l700
					}
				l703:
					{
						position704, tokenIndex704, depth704 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l704
						}
						if !_rules[rulesp]() {
							goto l704
						}
						if !_rules[ruleproductExpr]() {
							goto l704
						}
						if !_rules[rulesp]() {
							goto l704
						}
						goto l703
					l704:
						position, tokenIndex, depth = position704, tokenIndex704, depth704
					}
					depth--
					add(rulePegText, position702)
				}
				if !_rules[ruleAction41]() {
					goto l700
				}
				depth--
				add(ruletermExpr, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 55 productExpr <- <(<(minusExpr sp (MultDivOp sp minusExpr sp)*)> Action42)> */
		func() bool {
			position705, tokenIndex705, depth705 := position, tokenIndex, depth
			{
				position706 := position
				depth++
				{
					position707 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l705
					}
					if !_rules[rulesp]() {
						goto l705
					}
				l708:
					{
						position709, tokenIndex709, depth709 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l709
						}
						if !_rules[rulesp]() {
							goto l709
						}
						if !_rules[ruleminusExpr]() {
							goto l709
						}
						if !_rules[rulesp]() {
							goto l709
						}
						goto l708
					l709:
						position, tokenIndex, depth = position709, tokenIndex709, depth709
					}
					depth--
					add(rulePegText, position707)
				}
				if !_rules[ruleAction42]() {
					goto l705
				}
				depth--
				add(ruleproductExpr, position706)
			}
			return true
		l705:
			position, tokenIndex, depth = position705, tokenIndex705, depth705
			return false
		},
		/* 56 minusExpr <- <(<((UnaryMinus sp)? castExpr)> Action43)> */
		func() bool {
			position710, tokenIndex710, depth710 := position, tokenIndex, depth
			{
				position711 := position
				depth++
				{
					position712 := position
					depth++
					{
						position713, tokenIndex713, depth713 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l713
						}
						if !_rules[rulesp]() {
							goto l713
						}
						goto l714
					l713:
						position, tokenIndex, depth = position713, tokenIndex713, depth713
					}
				l714:
					if !_rules[rulecastExpr]() {
						goto l710
					}
					depth--
					add(rulePegText, position712)
				}
				if !_rules[ruleAction43]() {
					goto l710
				}
				depth--
				add(ruleminusExpr, position711)
			}
			return true
		l710:
			position, tokenIndex, depth = position710, tokenIndex710, depth710
			return false
		},
		/* 57 castExpr <- <(<(baseExpr (sp (':' ':') sp Type)?)> Action44)> */
		func() bool {
			position715, tokenIndex715, depth715 := position, tokenIndex, depth
			{
				position716 := position
				depth++
				{
					position717 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l715
					}
					{
						position718, tokenIndex718, depth718 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l718
						}
						if buffer[position] != rune(':') {
							goto l718
						}
						position++
						if buffer[position] != rune(':') {
							goto l718
						}
						position++
						if !_rules[rulesp]() {
							goto l718
						}
						if !_rules[ruleType]() {
							goto l718
						}
						goto l719
					l718:
						position, tokenIndex, depth = position718, tokenIndex718, depth718
					}
				l719:
					depth--
					add(rulePegText, position717)
				}
				if !_rules[ruleAction44]() {
					goto l715
				}
				depth--
				add(rulecastExpr, position716)
			}
			return true
		l715:
			position, tokenIndex, depth = position715, tokenIndex715, depth715
			return false
		},
		/* 58 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / NullLiteral / RowMeta / FuncTypeCast / FuncApp / RowValue / Literal)> */
		func() bool {
			position720, tokenIndex720, depth720 := position, tokenIndex, depth
			{
				position721 := position
				depth++
				{
					position722, tokenIndex722, depth722 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l723
					}
					position++
					if !_rules[rulesp]() {
						goto l723
					}
					if !_rules[ruleExpression]() {
						goto l723
					}
					if !_rules[rulesp]() {
						goto l723
					}
					if buffer[position] != rune(')') {
						goto l723
					}
					position++
					goto l722
				l723:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleBooleanLiteral]() {
						goto l724
					}
					goto l722
				l724:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleNullLiteral]() {
						goto l725
					}
					goto l722
				l725:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleRowMeta]() {
						goto l726
					}
					goto l722
				l726:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleFuncTypeCast]() {
						goto l727
					}
					goto l722
				l727:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleFuncApp]() {
						goto l728
					}
					goto l722
				l728:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleRowValue]() {
						goto l729
					}
					goto l722
				l729:
					position, tokenIndex, depth = position722, tokenIndex722, depth722
					if !_rules[ruleLiteral]() {
						goto l720
					}
				}
			l722:
				depth--
				add(rulebaseExpr, position721)
			}
			return true
		l720:
			position, tokenIndex, depth = position720, tokenIndex720, depth720
			return false
		},
		/* 59 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') sp '(' sp Expression sp (('a' / 'A') ('s' / 'S')) sp Type sp ')')> Action45)> */
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
						if buffer[position] != rune('c') {
							goto l734
						}
						position++
						goto l733
					l734:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
						if buffer[position] != rune('C') {
							goto l730
						}
						position++
					}
				l733:
					{
						position735, tokenIndex735, depth735 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l736
						}
						position++
						goto l735
					l736:
						position, tokenIndex, depth = position735, tokenIndex735, depth735
						if buffer[position] != rune('A') {
							goto l730
						}
						position++
					}
				l735:
					{
						position737, tokenIndex737, depth737 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l738
						}
						position++
						goto l737
					l738:
						position, tokenIndex, depth = position737, tokenIndex737, depth737
						if buffer[position] != rune('S') {
							goto l730
						}
						position++
					}
				l737:
					{
						position739, tokenIndex739, depth739 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l740
						}
						position++
						goto l739
					l740:
						position, tokenIndex, depth = position739, tokenIndex739, depth739
						if buffer[position] != rune('T') {
							goto l730
						}
						position++
					}
				l739:
					if !_rules[rulesp]() {
						goto l730
					}
					if buffer[position] != rune('(') {
						goto l730
					}
					position++
					if !_rules[rulesp]() {
						goto l730
					}
					if !_rules[ruleExpression]() {
						goto l730
					}
					if !_rules[rulesp]() {
						goto l730
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
							goto l730
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
							goto l730
						}
						position++
					}
				l743:
					if !_rules[rulesp]() {
						goto l730
					}
					if !_rules[ruleType]() {
						goto l730
					}
					if !_rules[rulesp]() {
						goto l730
					}
					if buffer[position] != rune(')') {
						goto l730
					}
					position++
					depth--
					add(rulePegText, position732)
				}
				if !_rules[ruleAction45]() {
					goto l730
				}
				depth--
				add(ruleFuncTypeCast, position731)
			}
			return true
		l730:
			position, tokenIndex, depth = position730, tokenIndex730, depth730
			return false
		},
		/* 60 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action46)> */
		func() bool {
			position745, tokenIndex745, depth745 := position, tokenIndex, depth
			{
				position746 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l745
				}
				if !_rules[rulesp]() {
					goto l745
				}
				if buffer[position] != rune('(') {
					goto l745
				}
				position++
				if !_rules[rulesp]() {
					goto l745
				}
				if !_rules[ruleFuncParams]() {
					goto l745
				}
				if !_rules[rulesp]() {
					goto l745
				}
				if buffer[position] != rune(')') {
					goto l745
				}
				position++
				if !_rules[ruleAction46]() {
					goto l745
				}
				depth--
				add(ruleFuncApp, position746)
			}
			return true
		l745:
			position, tokenIndex, depth = position745, tokenIndex745, depth745
			return false
		},
		/* 61 FuncParams <- <(<(Wildcard / (Expression sp (',' sp Expression)*)?)> Action47)> */
		func() bool {
			position747, tokenIndex747, depth747 := position, tokenIndex, depth
			{
				position748 := position
				depth++
				{
					position749 := position
					depth++
					{
						position750, tokenIndex750, depth750 := position, tokenIndex, depth
						if !_rules[ruleWildcard]() {
							goto l751
						}
						goto l750
					l751:
						position, tokenIndex, depth = position750, tokenIndex750, depth750
						{
							position752, tokenIndex752, depth752 := position, tokenIndex, depth
							if !_rules[ruleExpression]() {
								goto l752
							}
							if !_rules[rulesp]() {
								goto l752
							}
						l754:
							{
								position755, tokenIndex755, depth755 := position, tokenIndex, depth
								if buffer[position] != rune(',') {
									goto l755
								}
								position++
								if !_rules[rulesp]() {
									goto l755
								}
								if !_rules[ruleExpression]() {
									goto l755
								}
								goto l754
							l755:
								position, tokenIndex, depth = position755, tokenIndex755, depth755
							}
							goto l753
						l752:
							position, tokenIndex, depth = position752, tokenIndex752, depth752
						}
					l753:
					}
				l750:
					depth--
					add(rulePegText, position749)
				}
				if !_rules[ruleAction47]() {
					goto l747
				}
				depth--
				add(ruleFuncParams, position748)
			}
			return true
		l747:
			position, tokenIndex, depth = position747, tokenIndex747, depth747
			return false
		},
		/* 62 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position756, tokenIndex756, depth756 := position, tokenIndex, depth
			{
				position757 := position
				depth++
				{
					position758, tokenIndex758, depth758 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l759
					}
					goto l758
				l759:
					position, tokenIndex, depth = position758, tokenIndex758, depth758
					if !_rules[ruleNumericLiteral]() {
						goto l760
					}
					goto l758
				l760:
					position, tokenIndex, depth = position758, tokenIndex758, depth758
					if !_rules[ruleStringLiteral]() {
						goto l756
					}
				}
			l758:
				depth--
				add(ruleLiteral, position757)
			}
			return true
		l756:
			position, tokenIndex, depth = position756, tokenIndex756, depth756
			return false
		},
		/* 63 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position761, tokenIndex761, depth761 := position, tokenIndex, depth
			{
				position762 := position
				depth++
				{
					position763, tokenIndex763, depth763 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l764
					}
					goto l763
				l764:
					position, tokenIndex, depth = position763, tokenIndex763, depth763
					if !_rules[ruleNotEqual]() {
						goto l765
					}
					goto l763
				l765:
					position, tokenIndex, depth = position763, tokenIndex763, depth763
					if !_rules[ruleLessOrEqual]() {
						goto l766
					}
					goto l763
				l766:
					position, tokenIndex, depth = position763, tokenIndex763, depth763
					if !_rules[ruleLess]() {
						goto l767
					}
					goto l763
				l767:
					position, tokenIndex, depth = position763, tokenIndex763, depth763
					if !_rules[ruleGreaterOrEqual]() {
						goto l768
					}
					goto l763
				l768:
					position, tokenIndex, depth = position763, tokenIndex763, depth763
					if !_rules[ruleGreater]() {
						goto l769
					}
					goto l763
				l769:
					position, tokenIndex, depth = position763, tokenIndex763, depth763
					if !_rules[ruleNotEqual]() {
						goto l761
					}
				}
			l763:
				depth--
				add(ruleComparisonOp, position762)
			}
			return true
		l761:
			position, tokenIndex, depth = position761, tokenIndex761, depth761
			return false
		},
		/* 64 OtherOp <- <Concat> */
		func() bool {
			position770, tokenIndex770, depth770 := position, tokenIndex, depth
			{
				position771 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l770
				}
				depth--
				add(ruleOtherOp, position771)
			}
			return true
		l770:
			position, tokenIndex, depth = position770, tokenIndex770, depth770
			return false
		},
		/* 65 IsOp <- <(IsNot / Is)> */
		func() bool {
			position772, tokenIndex772, depth772 := position, tokenIndex, depth
			{
				position773 := position
				depth++
				{
					position774, tokenIndex774, depth774 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l775
					}
					goto l774
				l775:
					position, tokenIndex, depth = position774, tokenIndex774, depth774
					if !_rules[ruleIs]() {
						goto l772
					}
				}
			l774:
				depth--
				add(ruleIsOp, position773)
			}
			return true
		l772:
			position, tokenIndex, depth = position772, tokenIndex772, depth772
			return false
		},
		/* 66 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position776, tokenIndex776, depth776 := position, tokenIndex, depth
			{
				position777 := position
				depth++
				{
					position778, tokenIndex778, depth778 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l779
					}
					goto l778
				l779:
					position, tokenIndex, depth = position778, tokenIndex778, depth778
					if !_rules[ruleMinus]() {
						goto l776
					}
				}
			l778:
				depth--
				add(rulePlusMinusOp, position777)
			}
			return true
		l776:
			position, tokenIndex, depth = position776, tokenIndex776, depth776
			return false
		},
		/* 67 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position780, tokenIndex780, depth780 := position, tokenIndex, depth
			{
				position781 := position
				depth++
				{
					position782, tokenIndex782, depth782 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l783
					}
					goto l782
				l783:
					position, tokenIndex, depth = position782, tokenIndex782, depth782
					if !_rules[ruleDivide]() {
						goto l784
					}
					goto l782
				l784:
					position, tokenIndex, depth = position782, tokenIndex782, depth782
					if !_rules[ruleModulo]() {
						goto l780
					}
				}
			l782:
				depth--
				add(ruleMultDivOp, position781)
			}
			return true
		l780:
			position, tokenIndex, depth = position780, tokenIndex780, depth780
			return false
		},
		/* 68 Stream <- <(<ident> Action48)> */
		func() bool {
			position785, tokenIndex785, depth785 := position, tokenIndex, depth
			{
				position786 := position
				depth++
				{
					position787 := position
					depth++
					if !_rules[ruleident]() {
						goto l785
					}
					depth--
					add(rulePegText, position787)
				}
				if !_rules[ruleAction48]() {
					goto l785
				}
				depth--
				add(ruleStream, position786)
			}
			return true
		l785:
			position, tokenIndex, depth = position785, tokenIndex785, depth785
			return false
		},
		/* 69 RowMeta <- <RowTimestamp> */
		func() bool {
			position788, tokenIndex788, depth788 := position, tokenIndex, depth
			{
				position789 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l788
				}
				depth--
				add(ruleRowMeta, position789)
			}
			return true
		l788:
			position, tokenIndex, depth = position788, tokenIndex788, depth788
			return false
		},
		/* 70 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action49)> */
		func() bool {
			position790, tokenIndex790, depth790 := position, tokenIndex, depth
			{
				position791 := position
				depth++
				{
					position792 := position
					depth++
					{
						position793, tokenIndex793, depth793 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l793
						}
						if buffer[position] != rune(':') {
							goto l793
						}
						position++
						goto l794
					l793:
						position, tokenIndex, depth = position793, tokenIndex793, depth793
					}
				l794:
					if buffer[position] != rune('t') {
						goto l790
					}
					position++
					if buffer[position] != rune('s') {
						goto l790
					}
					position++
					if buffer[position] != rune('(') {
						goto l790
					}
					position++
					if buffer[position] != rune(')') {
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
				add(ruleRowTimestamp, position791)
			}
			return true
		l790:
			position, tokenIndex, depth = position790, tokenIndex790, depth790
			return false
		},
		/* 71 RowValue <- <(<((ident ':' !':')? jsonPath)> Action50)> */
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
						if !_rules[ruleident]() {
							goto l798
						}
						if buffer[position] != rune(':') {
							goto l798
						}
						position++
						{
							position800, tokenIndex800, depth800 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l800
							}
							position++
							goto l798
						l800:
							position, tokenIndex, depth = position800, tokenIndex800, depth800
						}
						goto l799
					l798:
						position, tokenIndex, depth = position798, tokenIndex798, depth798
					}
				l799:
					if !_rules[rulejsonPath]() {
						goto l795
					}
					depth--
					add(rulePegText, position797)
				}
				if !_rules[ruleAction50]() {
					goto l795
				}
				depth--
				add(ruleRowValue, position796)
			}
			return true
		l795:
			position, tokenIndex, depth = position795, tokenIndex795, depth795
			return false
		},
		/* 72 NumericLiteral <- <(<('-'? [0-9]+)> Action51)> */
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
						if buffer[position] != rune('-') {
							goto l804
						}
						position++
						goto l805
					l804:
						position, tokenIndex, depth = position804, tokenIndex804, depth804
					}
				l805:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l801
					}
					position++
				l806:
					{
						position807, tokenIndex807, depth807 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l807
						}
						position++
						goto l806
					l807:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
					}
					depth--
					add(rulePegText, position803)
				}
				if !_rules[ruleAction51]() {
					goto l801
				}
				depth--
				add(ruleNumericLiteral, position802)
			}
			return true
		l801:
			position, tokenIndex, depth = position801, tokenIndex801, depth801
			return false
		},
		/* 73 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action52)> */
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
						if buffer[position] != rune('-') {
							goto l811
						}
						position++
						goto l812
					l811:
						position, tokenIndex, depth = position811, tokenIndex811, depth811
					}
				l812:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l808
					}
					position++
				l813:
					{
						position814, tokenIndex814, depth814 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l814
						}
						position++
						goto l813
					l814:
						position, tokenIndex, depth = position814, tokenIndex814, depth814
					}
					if buffer[position] != rune('.') {
						goto l808
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l808
					}
					position++
				l815:
					{
						position816, tokenIndex816, depth816 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l816
						}
						position++
						goto l815
					l816:
						position, tokenIndex, depth = position816, tokenIndex816, depth816
					}
					depth--
					add(rulePegText, position810)
				}
				if !_rules[ruleAction52]() {
					goto l808
				}
				depth--
				add(ruleFloatLiteral, position809)
			}
			return true
		l808:
			position, tokenIndex, depth = position808, tokenIndex808, depth808
			return false
		},
		/* 74 Function <- <(<ident> Action53)> */
		func() bool {
			position817, tokenIndex817, depth817 := position, tokenIndex, depth
			{
				position818 := position
				depth++
				{
					position819 := position
					depth++
					if !_rules[ruleident]() {
						goto l817
					}
					depth--
					add(rulePegText, position819)
				}
				if !_rules[ruleAction53]() {
					goto l817
				}
				depth--
				add(ruleFunction, position818)
			}
			return true
		l817:
			position, tokenIndex, depth = position817, tokenIndex817, depth817
			return false
		},
		/* 75 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action54)> */
		func() bool {
			position820, tokenIndex820, depth820 := position, tokenIndex, depth
			{
				position821 := position
				depth++
				{
					position822 := position
					depth++
					{
						position823, tokenIndex823, depth823 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l824
						}
						position++
						goto l823
					l824:
						position, tokenIndex, depth = position823, tokenIndex823, depth823
						if buffer[position] != rune('N') {
							goto l820
						}
						position++
					}
				l823:
					{
						position825, tokenIndex825, depth825 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l826
						}
						position++
						goto l825
					l826:
						position, tokenIndex, depth = position825, tokenIndex825, depth825
						if buffer[position] != rune('U') {
							goto l820
						}
						position++
					}
				l825:
					{
						position827, tokenIndex827, depth827 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l828
						}
						position++
						goto l827
					l828:
						position, tokenIndex, depth = position827, tokenIndex827, depth827
						if buffer[position] != rune('L') {
							goto l820
						}
						position++
					}
				l827:
					{
						position829, tokenIndex829, depth829 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l830
						}
						position++
						goto l829
					l830:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						if buffer[position] != rune('L') {
							goto l820
						}
						position++
					}
				l829:
					depth--
					add(rulePegText, position822)
				}
				if !_rules[ruleAction54]() {
					goto l820
				}
				depth--
				add(ruleNullLiteral, position821)
			}
			return true
		l820:
			position, tokenIndex, depth = position820, tokenIndex820, depth820
			return false
		},
		/* 76 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position831, tokenIndex831, depth831 := position, tokenIndex, depth
			{
				position832 := position
				depth++
				{
					position833, tokenIndex833, depth833 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l834
					}
					goto l833
				l834:
					position, tokenIndex, depth = position833, tokenIndex833, depth833
					if !_rules[ruleFALSE]() {
						goto l831
					}
				}
			l833:
				depth--
				add(ruleBooleanLiteral, position832)
			}
			return true
		l831:
			position, tokenIndex, depth = position831, tokenIndex831, depth831
			return false
		},
		/* 77 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action55)> */
		func() bool {
			position835, tokenIndex835, depth835 := position, tokenIndex, depth
			{
				position836 := position
				depth++
				{
					position837 := position
					depth++
					{
						position838, tokenIndex838, depth838 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l839
						}
						position++
						goto l838
					l839:
						position, tokenIndex, depth = position838, tokenIndex838, depth838
						if buffer[position] != rune('T') {
							goto l835
						}
						position++
					}
				l838:
					{
						position840, tokenIndex840, depth840 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l841
						}
						position++
						goto l840
					l841:
						position, tokenIndex, depth = position840, tokenIndex840, depth840
						if buffer[position] != rune('R') {
							goto l835
						}
						position++
					}
				l840:
					{
						position842, tokenIndex842, depth842 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l843
						}
						position++
						goto l842
					l843:
						position, tokenIndex, depth = position842, tokenIndex842, depth842
						if buffer[position] != rune('U') {
							goto l835
						}
						position++
					}
				l842:
					{
						position844, tokenIndex844, depth844 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l845
						}
						position++
						goto l844
					l845:
						position, tokenIndex, depth = position844, tokenIndex844, depth844
						if buffer[position] != rune('E') {
							goto l835
						}
						position++
					}
				l844:
					depth--
					add(rulePegText, position837)
				}
				if !_rules[ruleAction55]() {
					goto l835
				}
				depth--
				add(ruleTRUE, position836)
			}
			return true
		l835:
			position, tokenIndex, depth = position835, tokenIndex835, depth835
			return false
		},
		/* 78 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action56)> */
		func() bool {
			position846, tokenIndex846, depth846 := position, tokenIndex, depth
			{
				position847 := position
				depth++
				{
					position848 := position
					depth++
					{
						position849, tokenIndex849, depth849 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l850
						}
						position++
						goto l849
					l850:
						position, tokenIndex, depth = position849, tokenIndex849, depth849
						if buffer[position] != rune('F') {
							goto l846
						}
						position++
					}
				l849:
					{
						position851, tokenIndex851, depth851 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l852
						}
						position++
						goto l851
					l852:
						position, tokenIndex, depth = position851, tokenIndex851, depth851
						if buffer[position] != rune('A') {
							goto l846
						}
						position++
					}
				l851:
					{
						position853, tokenIndex853, depth853 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l854
						}
						position++
						goto l853
					l854:
						position, tokenIndex, depth = position853, tokenIndex853, depth853
						if buffer[position] != rune('L') {
							goto l846
						}
						position++
					}
				l853:
					{
						position855, tokenIndex855, depth855 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l856
						}
						position++
						goto l855
					l856:
						position, tokenIndex, depth = position855, tokenIndex855, depth855
						if buffer[position] != rune('S') {
							goto l846
						}
						position++
					}
				l855:
					{
						position857, tokenIndex857, depth857 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l858
						}
						position++
						goto l857
					l858:
						position, tokenIndex, depth = position857, tokenIndex857, depth857
						if buffer[position] != rune('E') {
							goto l846
						}
						position++
					}
				l857:
					depth--
					add(rulePegText, position848)
				}
				if !_rules[ruleAction56]() {
					goto l846
				}
				depth--
				add(ruleFALSE, position847)
			}
			return true
		l846:
			position, tokenIndex, depth = position846, tokenIndex846, depth846
			return false
		},
		/* 79 Wildcard <- <(<'*'> Action57)> */
		func() bool {
			position859, tokenIndex859, depth859 := position, tokenIndex, depth
			{
				position860 := position
				depth++
				{
					position861 := position
					depth++
					if buffer[position] != rune('*') {
						goto l859
					}
					position++
					depth--
					add(rulePegText, position861)
				}
				if !_rules[ruleAction57]() {
					goto l859
				}
				depth--
				add(ruleWildcard, position860)
			}
			return true
		l859:
			position, tokenIndex, depth = position859, tokenIndex859, depth859
			return false
		},
		/* 80 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action58)> */
		func() bool {
			position862, tokenIndex862, depth862 := position, tokenIndex, depth
			{
				position863 := position
				depth++
				{
					position864 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l862
					}
					position++
				l865:
					{
						position866, tokenIndex866, depth866 := position, tokenIndex, depth
						{
							position867, tokenIndex867, depth867 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l868
							}
							position++
							if buffer[position] != rune('\'') {
								goto l868
							}
							position++
							goto l867
						l868:
							position, tokenIndex, depth = position867, tokenIndex867, depth867
							{
								position869, tokenIndex869, depth869 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l869
								}
								position++
								goto l866
							l869:
								position, tokenIndex, depth = position869, tokenIndex869, depth869
							}
							if !matchDot() {
								goto l866
							}
						}
					l867:
						goto l865
					l866:
						position, tokenIndex, depth = position866, tokenIndex866, depth866
					}
					if buffer[position] != rune('\'') {
						goto l862
					}
					position++
					depth--
					add(rulePegText, position864)
				}
				if !_rules[ruleAction58]() {
					goto l862
				}
				depth--
				add(ruleStringLiteral, position863)
			}
			return true
		l862:
			position, tokenIndex, depth = position862, tokenIndex862, depth862
			return false
		},
		/* 81 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action59)> */
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
						if buffer[position] != rune('i') {
							goto l874
						}
						position++
						goto l873
					l874:
						position, tokenIndex, depth = position873, tokenIndex873, depth873
						if buffer[position] != rune('I') {
							goto l870
						}
						position++
					}
				l873:
					{
						position875, tokenIndex875, depth875 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l876
						}
						position++
						goto l875
					l876:
						position, tokenIndex, depth = position875, tokenIndex875, depth875
						if buffer[position] != rune('S') {
							goto l870
						}
						position++
					}
				l875:
					{
						position877, tokenIndex877, depth877 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l878
						}
						position++
						goto l877
					l878:
						position, tokenIndex, depth = position877, tokenIndex877, depth877
						if buffer[position] != rune('T') {
							goto l870
						}
						position++
					}
				l877:
					{
						position879, tokenIndex879, depth879 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l880
						}
						position++
						goto l879
					l880:
						position, tokenIndex, depth = position879, tokenIndex879, depth879
						if buffer[position] != rune('R') {
							goto l870
						}
						position++
					}
				l879:
					{
						position881, tokenIndex881, depth881 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l882
						}
						position++
						goto l881
					l882:
						position, tokenIndex, depth = position881, tokenIndex881, depth881
						if buffer[position] != rune('E') {
							goto l870
						}
						position++
					}
				l881:
					{
						position883, tokenIndex883, depth883 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l884
						}
						position++
						goto l883
					l884:
						position, tokenIndex, depth = position883, tokenIndex883, depth883
						if buffer[position] != rune('A') {
							goto l870
						}
						position++
					}
				l883:
					{
						position885, tokenIndex885, depth885 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l886
						}
						position++
						goto l885
					l886:
						position, tokenIndex, depth = position885, tokenIndex885, depth885
						if buffer[position] != rune('M') {
							goto l870
						}
						position++
					}
				l885:
					depth--
					add(rulePegText, position872)
				}
				if !_rules[ruleAction59]() {
					goto l870
				}
				depth--
				add(ruleISTREAM, position871)
			}
			return true
		l870:
			position, tokenIndex, depth = position870, tokenIndex870, depth870
			return false
		},
		/* 82 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action60)> */
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
						if buffer[position] != rune('d') {
							goto l891
						}
						position++
						goto l890
					l891:
						position, tokenIndex, depth = position890, tokenIndex890, depth890
						if buffer[position] != rune('D') {
							goto l887
						}
						position++
					}
				l890:
					{
						position892, tokenIndex892, depth892 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l893
						}
						position++
						goto l892
					l893:
						position, tokenIndex, depth = position892, tokenIndex892, depth892
						if buffer[position] != rune('S') {
							goto l887
						}
						position++
					}
				l892:
					{
						position894, tokenIndex894, depth894 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l895
						}
						position++
						goto l894
					l895:
						position, tokenIndex, depth = position894, tokenIndex894, depth894
						if buffer[position] != rune('T') {
							goto l887
						}
						position++
					}
				l894:
					{
						position896, tokenIndex896, depth896 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l897
						}
						position++
						goto l896
					l897:
						position, tokenIndex, depth = position896, tokenIndex896, depth896
						if buffer[position] != rune('R') {
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
					{
						position900, tokenIndex900, depth900 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l901
						}
						position++
						goto l900
					l901:
						position, tokenIndex, depth = position900, tokenIndex900, depth900
						if buffer[position] != rune('A') {
							goto l887
						}
						position++
					}
				l900:
					{
						position902, tokenIndex902, depth902 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l903
						}
						position++
						goto l902
					l903:
						position, tokenIndex, depth = position902, tokenIndex902, depth902
						if buffer[position] != rune('M') {
							goto l887
						}
						position++
					}
				l902:
					depth--
					add(rulePegText, position889)
				}
				if !_rules[ruleAction60]() {
					goto l887
				}
				depth--
				add(ruleDSTREAM, position888)
			}
			return true
		l887:
			position, tokenIndex, depth = position887, tokenIndex887, depth887
			return false
		},
		/* 83 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action61)> */
		func() bool {
			position904, tokenIndex904, depth904 := position, tokenIndex, depth
			{
				position905 := position
				depth++
				{
					position906 := position
					depth++
					{
						position907, tokenIndex907, depth907 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l908
						}
						position++
						goto l907
					l908:
						position, tokenIndex, depth = position907, tokenIndex907, depth907
						if buffer[position] != rune('R') {
							goto l904
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
							goto l904
						}
						position++
					}
				l909:
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
							goto l904
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
							goto l904
						}
						position++
					}
				l913:
					{
						position915, tokenIndex915, depth915 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l916
						}
						position++
						goto l915
					l916:
						position, tokenIndex, depth = position915, tokenIndex915, depth915
						if buffer[position] != rune('E') {
							goto l904
						}
						position++
					}
				l915:
					{
						position917, tokenIndex917, depth917 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l918
						}
						position++
						goto l917
					l918:
						position, tokenIndex, depth = position917, tokenIndex917, depth917
						if buffer[position] != rune('A') {
							goto l904
						}
						position++
					}
				l917:
					{
						position919, tokenIndex919, depth919 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l920
						}
						position++
						goto l919
					l920:
						position, tokenIndex, depth = position919, tokenIndex919, depth919
						if buffer[position] != rune('M') {
							goto l904
						}
						position++
					}
				l919:
					depth--
					add(rulePegText, position906)
				}
				if !_rules[ruleAction61]() {
					goto l904
				}
				depth--
				add(ruleRSTREAM, position905)
			}
			return true
		l904:
			position, tokenIndex, depth = position904, tokenIndex904, depth904
			return false
		},
		/* 84 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action62)> */
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
						if buffer[position] != rune('t') {
							goto l925
						}
						position++
						goto l924
					l925:
						position, tokenIndex, depth = position924, tokenIndex924, depth924
						if buffer[position] != rune('T') {
							goto l921
						}
						position++
					}
				l924:
					{
						position926, tokenIndex926, depth926 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l927
						}
						position++
						goto l926
					l927:
						position, tokenIndex, depth = position926, tokenIndex926, depth926
						if buffer[position] != rune('U') {
							goto l921
						}
						position++
					}
				l926:
					{
						position928, tokenIndex928, depth928 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l929
						}
						position++
						goto l928
					l929:
						position, tokenIndex, depth = position928, tokenIndex928, depth928
						if buffer[position] != rune('P') {
							goto l921
						}
						position++
					}
				l928:
					{
						position930, tokenIndex930, depth930 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l931
						}
						position++
						goto l930
					l931:
						position, tokenIndex, depth = position930, tokenIndex930, depth930
						if buffer[position] != rune('L') {
							goto l921
						}
						position++
					}
				l930:
					{
						position932, tokenIndex932, depth932 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l933
						}
						position++
						goto l932
					l933:
						position, tokenIndex, depth = position932, tokenIndex932, depth932
						if buffer[position] != rune('E') {
							goto l921
						}
						position++
					}
				l932:
					{
						position934, tokenIndex934, depth934 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l935
						}
						position++
						goto l934
					l935:
						position, tokenIndex, depth = position934, tokenIndex934, depth934
						if buffer[position] != rune('S') {
							goto l921
						}
						position++
					}
				l934:
					depth--
					add(rulePegText, position923)
				}
				if !_rules[ruleAction62]() {
					goto l921
				}
				depth--
				add(ruleTUPLES, position922)
			}
			return true
		l921:
			position, tokenIndex, depth = position921, tokenIndex921, depth921
			return false
		},
		/* 85 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action63)> */
		func() bool {
			position936, tokenIndex936, depth936 := position, tokenIndex, depth
			{
				position937 := position
				depth++
				{
					position938 := position
					depth++
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
							goto l936
						}
						position++
					}
				l939:
					{
						position941, tokenIndex941, depth941 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l942
						}
						position++
						goto l941
					l942:
						position, tokenIndex, depth = position941, tokenIndex941, depth941
						if buffer[position] != rune('E') {
							goto l936
						}
						position++
					}
				l941:
					{
						position943, tokenIndex943, depth943 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l944
						}
						position++
						goto l943
					l944:
						position, tokenIndex, depth = position943, tokenIndex943, depth943
						if buffer[position] != rune('C') {
							goto l936
						}
						position++
					}
				l943:
					{
						position945, tokenIndex945, depth945 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l946
						}
						position++
						goto l945
					l946:
						position, tokenIndex, depth = position945, tokenIndex945, depth945
						if buffer[position] != rune('O') {
							goto l936
						}
						position++
					}
				l945:
					{
						position947, tokenIndex947, depth947 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l948
						}
						position++
						goto l947
					l948:
						position, tokenIndex, depth = position947, tokenIndex947, depth947
						if buffer[position] != rune('N') {
							goto l936
						}
						position++
					}
				l947:
					{
						position949, tokenIndex949, depth949 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l950
						}
						position++
						goto l949
					l950:
						position, tokenIndex, depth = position949, tokenIndex949, depth949
						if buffer[position] != rune('D') {
							goto l936
						}
						position++
					}
				l949:
					{
						position951, tokenIndex951, depth951 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l952
						}
						position++
						goto l951
					l952:
						position, tokenIndex, depth = position951, tokenIndex951, depth951
						if buffer[position] != rune('S') {
							goto l936
						}
						position++
					}
				l951:
					depth--
					add(rulePegText, position938)
				}
				if !_rules[ruleAction63]() {
					goto l936
				}
				depth--
				add(ruleSECONDS, position937)
			}
			return true
		l936:
			position, tokenIndex, depth = position936, tokenIndex936, depth936
			return false
		},
		/* 86 StreamIdentifier <- <(<ident> Action64)> */
		func() bool {
			position953, tokenIndex953, depth953 := position, tokenIndex, depth
			{
				position954 := position
				depth++
				{
					position955 := position
					depth++
					if !_rules[ruleident]() {
						goto l953
					}
					depth--
					add(rulePegText, position955)
				}
				if !_rules[ruleAction64]() {
					goto l953
				}
				depth--
				add(ruleStreamIdentifier, position954)
			}
			return true
		l953:
			position, tokenIndex, depth = position953, tokenIndex953, depth953
			return false
		},
		/* 87 SourceSinkType <- <(<ident> Action65)> */
		func() bool {
			position956, tokenIndex956, depth956 := position, tokenIndex, depth
			{
				position957 := position
				depth++
				{
					position958 := position
					depth++
					if !_rules[ruleident]() {
						goto l956
					}
					depth--
					add(rulePegText, position958)
				}
				if !_rules[ruleAction65]() {
					goto l956
				}
				depth--
				add(ruleSourceSinkType, position957)
			}
			return true
		l956:
			position, tokenIndex, depth = position956, tokenIndex956, depth956
			return false
		},
		/* 88 SourceSinkParamKey <- <(<ident> Action66)> */
		func() bool {
			position959, tokenIndex959, depth959 := position, tokenIndex, depth
			{
				position960 := position
				depth++
				{
					position961 := position
					depth++
					if !_rules[ruleident]() {
						goto l959
					}
					depth--
					add(rulePegText, position961)
				}
				if !_rules[ruleAction66]() {
					goto l959
				}
				depth--
				add(ruleSourceSinkParamKey, position960)
			}
			return true
		l959:
			position, tokenIndex, depth = position959, tokenIndex959, depth959
			return false
		},
		/* 89 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action67)> */
		func() bool {
			position962, tokenIndex962, depth962 := position, tokenIndex, depth
			{
				position963 := position
				depth++
				{
					position964 := position
					depth++
					{
						position965, tokenIndex965, depth965 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l966
						}
						position++
						goto l965
					l966:
						position, tokenIndex, depth = position965, tokenIndex965, depth965
						if buffer[position] != rune('P') {
							goto l962
						}
						position++
					}
				l965:
					{
						position967, tokenIndex967, depth967 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l968
						}
						position++
						goto l967
					l968:
						position, tokenIndex, depth = position967, tokenIndex967, depth967
						if buffer[position] != rune('A') {
							goto l962
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
							goto l962
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
							goto l962
						}
						position++
					}
				l971:
					{
						position973, tokenIndex973, depth973 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l974
						}
						position++
						goto l973
					l974:
						position, tokenIndex, depth = position973, tokenIndex973, depth973
						if buffer[position] != rune('E') {
							goto l962
						}
						position++
					}
				l973:
					{
						position975, tokenIndex975, depth975 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l976
						}
						position++
						goto l975
					l976:
						position, tokenIndex, depth = position975, tokenIndex975, depth975
						if buffer[position] != rune('D') {
							goto l962
						}
						position++
					}
				l975:
					depth--
					add(rulePegText, position964)
				}
				if !_rules[ruleAction67]() {
					goto l962
				}
				depth--
				add(rulePaused, position963)
			}
			return true
		l962:
			position, tokenIndex, depth = position962, tokenIndex962, depth962
			return false
		},
		/* 90 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action68)> */
		func() bool {
			position977, tokenIndex977, depth977 := position, tokenIndex, depth
			{
				position978 := position
				depth++
				{
					position979 := position
					depth++
					{
						position980, tokenIndex980, depth980 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l981
						}
						position++
						goto l980
					l981:
						position, tokenIndex, depth = position980, tokenIndex980, depth980
						if buffer[position] != rune('U') {
							goto l977
						}
						position++
					}
				l980:
					{
						position982, tokenIndex982, depth982 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l983
						}
						position++
						goto l982
					l983:
						position, tokenIndex, depth = position982, tokenIndex982, depth982
						if buffer[position] != rune('N') {
							goto l977
						}
						position++
					}
				l982:
					{
						position984, tokenIndex984, depth984 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l985
						}
						position++
						goto l984
					l985:
						position, tokenIndex, depth = position984, tokenIndex984, depth984
						if buffer[position] != rune('P') {
							goto l977
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
							goto l977
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
							goto l977
						}
						position++
					}
				l988:
					{
						position990, tokenIndex990, depth990 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l991
						}
						position++
						goto l990
					l991:
						position, tokenIndex, depth = position990, tokenIndex990, depth990
						if buffer[position] != rune('S') {
							goto l977
						}
						position++
					}
				l990:
					{
						position992, tokenIndex992, depth992 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l993
						}
						position++
						goto l992
					l993:
						position, tokenIndex, depth = position992, tokenIndex992, depth992
						if buffer[position] != rune('E') {
							goto l977
						}
						position++
					}
				l992:
					{
						position994, tokenIndex994, depth994 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l995
						}
						position++
						goto l994
					l995:
						position, tokenIndex, depth = position994, tokenIndex994, depth994
						if buffer[position] != rune('D') {
							goto l977
						}
						position++
					}
				l994:
					depth--
					add(rulePegText, position979)
				}
				if !_rules[ruleAction68]() {
					goto l977
				}
				depth--
				add(ruleUnpaused, position978)
			}
			return true
		l977:
			position, tokenIndex, depth = position977, tokenIndex977, depth977
			return false
		},
		/* 91 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position996, tokenIndex996, depth996 := position, tokenIndex, depth
			{
				position997 := position
				depth++
				{
					position998, tokenIndex998, depth998 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l999
					}
					goto l998
				l999:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if !_rules[ruleInt]() {
						goto l1000
					}
					goto l998
				l1000:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if !_rules[ruleFloat]() {
						goto l1001
					}
					goto l998
				l1001:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if !_rules[ruleString]() {
						goto l1002
					}
					goto l998
				l1002:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if !_rules[ruleBlob]() {
						goto l1003
					}
					goto l998
				l1003:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if !_rules[ruleTimestamp]() {
						goto l1004
					}
					goto l998
				l1004:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if !_rules[ruleArray]() {
						goto l1005
					}
					goto l998
				l1005:
					position, tokenIndex, depth = position998, tokenIndex998, depth998
					if !_rules[ruleMap]() {
						goto l996
					}
				}
			l998:
				depth--
				add(ruleType, position997)
			}
			return true
		l996:
			position, tokenIndex, depth = position996, tokenIndex996, depth996
			return false
		},
		/* 92 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action69)> */
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
						if buffer[position] != rune('b') {
							goto l1010
						}
						position++
						goto l1009
					l1010:
						position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
						if buffer[position] != rune('B') {
							goto l1006
						}
						position++
					}
				l1009:
					{
						position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1012
						}
						position++
						goto l1011
					l1012:
						position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
						if buffer[position] != rune('O') {
							goto l1006
						}
						position++
					}
				l1011:
					{
						position1013, tokenIndex1013, depth1013 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1014
						}
						position++
						goto l1013
					l1014:
						position, tokenIndex, depth = position1013, tokenIndex1013, depth1013
						if buffer[position] != rune('O') {
							goto l1006
						}
						position++
					}
				l1013:
					{
						position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1016
						}
						position++
						goto l1015
					l1016:
						position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
						if buffer[position] != rune('L') {
							goto l1006
						}
						position++
					}
				l1015:
					depth--
					add(rulePegText, position1008)
				}
				if !_rules[ruleAction69]() {
					goto l1006
				}
				depth--
				add(ruleBool, position1007)
			}
			return true
		l1006:
			position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
			return false
		},
		/* 93 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action70)> */
		func() bool {
			position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
			{
				position1018 := position
				depth++
				{
					position1019 := position
					depth++
					{
						position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1021
						}
						position++
						goto l1020
					l1021:
						position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
						if buffer[position] != rune('I') {
							goto l1017
						}
						position++
					}
				l1020:
					{
						position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1023
						}
						position++
						goto l1022
					l1023:
						position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
						if buffer[position] != rune('N') {
							goto l1017
						}
						position++
					}
				l1022:
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
							goto l1017
						}
						position++
					}
				l1024:
					depth--
					add(rulePegText, position1019)
				}
				if !_rules[ruleAction70]() {
					goto l1017
				}
				depth--
				add(ruleInt, position1018)
			}
			return true
		l1017:
			position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
			return false
		},
		/* 94 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action71)> */
		func() bool {
			position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
			{
				position1027 := position
				depth++
				{
					position1028 := position
					depth++
					{
						position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1030
						}
						position++
						goto l1029
					l1030:
						position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
						if buffer[position] != rune('F') {
							goto l1026
						}
						position++
					}
				l1029:
					{
						position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1032
						}
						position++
						goto l1031
					l1032:
						position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
						if buffer[position] != rune('L') {
							goto l1026
						}
						position++
					}
				l1031:
					{
						position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1034
						}
						position++
						goto l1033
					l1034:
						position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
						if buffer[position] != rune('O') {
							goto l1026
						}
						position++
					}
				l1033:
					{
						position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1036
						}
						position++
						goto l1035
					l1036:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						if buffer[position] != rune('A') {
							goto l1026
						}
						position++
					}
				l1035:
					{
						position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1038
						}
						position++
						goto l1037
					l1038:
						position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
						if buffer[position] != rune('T') {
							goto l1026
						}
						position++
					}
				l1037:
					depth--
					add(rulePegText, position1028)
				}
				if !_rules[ruleAction71]() {
					goto l1026
				}
				depth--
				add(ruleFloat, position1027)
			}
			return true
		l1026:
			position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
			return false
		},
		/* 95 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action72)> */
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
						if buffer[position] != rune('s') {
							goto l1043
						}
						position++
						goto l1042
					l1043:
						position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
						if buffer[position] != rune('S') {
							goto l1039
						}
						position++
					}
				l1042:
					{
						position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1045
						}
						position++
						goto l1044
					l1045:
						position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
						if buffer[position] != rune('T') {
							goto l1039
						}
						position++
					}
				l1044:
					{
						position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1047
						}
						position++
						goto l1046
					l1047:
						position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
						if buffer[position] != rune('R') {
							goto l1039
						}
						position++
					}
				l1046:
					{
						position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1049
						}
						position++
						goto l1048
					l1049:
						position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
						if buffer[position] != rune('I') {
							goto l1039
						}
						position++
					}
				l1048:
					{
						position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1051
						}
						position++
						goto l1050
					l1051:
						position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
						if buffer[position] != rune('N') {
							goto l1039
						}
						position++
					}
				l1050:
					{
						position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1053
						}
						position++
						goto l1052
					l1053:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
						if buffer[position] != rune('G') {
							goto l1039
						}
						position++
					}
				l1052:
					depth--
					add(rulePegText, position1041)
				}
				if !_rules[ruleAction72]() {
					goto l1039
				}
				depth--
				add(ruleString, position1040)
			}
			return true
		l1039:
			position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
			return false
		},
		/* 96 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action73)> */
		func() bool {
			position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
			{
				position1055 := position
				depth++
				{
					position1056 := position
					depth++
					{
						position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1058
						}
						position++
						goto l1057
					l1058:
						position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
						if buffer[position] != rune('B') {
							goto l1054
						}
						position++
					}
				l1057:
					{
						position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1060
						}
						position++
						goto l1059
					l1060:
						position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
						if buffer[position] != rune('L') {
							goto l1054
						}
						position++
					}
				l1059:
					{
						position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1062
						}
						position++
						goto l1061
					l1062:
						position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
						if buffer[position] != rune('O') {
							goto l1054
						}
						position++
					}
				l1061:
					{
						position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1064
						}
						position++
						goto l1063
					l1064:
						position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
						if buffer[position] != rune('B') {
							goto l1054
						}
						position++
					}
				l1063:
					depth--
					add(rulePegText, position1056)
				}
				if !_rules[ruleAction73]() {
					goto l1054
				}
				depth--
				add(ruleBlob, position1055)
			}
			return true
		l1054:
			position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
			return false
		},
		/* 97 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action74)> */
		func() bool {
			position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
			{
				position1066 := position
				depth++
				{
					position1067 := position
					depth++
					{
						position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1069
						}
						position++
						goto l1068
					l1069:
						position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
						if buffer[position] != rune('T') {
							goto l1065
						}
						position++
					}
				l1068:
					{
						position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1071
						}
						position++
						goto l1070
					l1071:
						position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
						if buffer[position] != rune('I') {
							goto l1065
						}
						position++
					}
				l1070:
					{
						position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1073
						}
						position++
						goto l1072
					l1073:
						position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
						if buffer[position] != rune('M') {
							goto l1065
						}
						position++
					}
				l1072:
					{
						position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1075
						}
						position++
						goto l1074
					l1075:
						position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
						if buffer[position] != rune('E') {
							goto l1065
						}
						position++
					}
				l1074:
					{
						position1076, tokenIndex1076, depth1076 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1077
						}
						position++
						goto l1076
					l1077:
						position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
						if buffer[position] != rune('S') {
							goto l1065
						}
						position++
					}
				l1076:
					{
						position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1079
						}
						position++
						goto l1078
					l1079:
						position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
						if buffer[position] != rune('T') {
							goto l1065
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
							goto l1065
						}
						position++
					}
				l1080:
					{
						position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1083
						}
						position++
						goto l1082
					l1083:
						position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
						if buffer[position] != rune('M') {
							goto l1065
						}
						position++
					}
				l1082:
					{
						position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1085
						}
						position++
						goto l1084
					l1085:
						position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
						if buffer[position] != rune('P') {
							goto l1065
						}
						position++
					}
				l1084:
					depth--
					add(rulePegText, position1067)
				}
				if !_rules[ruleAction74]() {
					goto l1065
				}
				depth--
				add(ruleTimestamp, position1066)
			}
			return true
		l1065:
			position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
			return false
		},
		/* 98 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action75)> */
		func() bool {
			position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
			{
				position1087 := position
				depth++
				{
					position1088 := position
					depth++
					{
						position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1090
						}
						position++
						goto l1089
					l1090:
						position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
						if buffer[position] != rune('A') {
							goto l1086
						}
						position++
					}
				l1089:
					{
						position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1092
						}
						position++
						goto l1091
					l1092:
						position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
						if buffer[position] != rune('R') {
							goto l1086
						}
						position++
					}
				l1091:
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
							goto l1086
						}
						position++
					}
				l1093:
					{
						position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1096
						}
						position++
						goto l1095
					l1096:
						position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
						if buffer[position] != rune('A') {
							goto l1086
						}
						position++
					}
				l1095:
					{
						position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1098
						}
						position++
						goto l1097
					l1098:
						position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
						if buffer[position] != rune('Y') {
							goto l1086
						}
						position++
					}
				l1097:
					depth--
					add(rulePegText, position1088)
				}
				if !_rules[ruleAction75]() {
					goto l1086
				}
				depth--
				add(ruleArray, position1087)
			}
			return true
		l1086:
			position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
			return false
		},
		/* 99 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action76)> */
		func() bool {
			position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
			{
				position1100 := position
				depth++
				{
					position1101 := position
					depth++
					{
						position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1103
						}
						position++
						goto l1102
					l1103:
						position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
						if buffer[position] != rune('M') {
							goto l1099
						}
						position++
					}
				l1102:
					{
						position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1105
						}
						position++
						goto l1104
					l1105:
						position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
						if buffer[position] != rune('A') {
							goto l1099
						}
						position++
					}
				l1104:
					{
						position1106, tokenIndex1106, depth1106 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1107
						}
						position++
						goto l1106
					l1107:
						position, tokenIndex, depth = position1106, tokenIndex1106, depth1106
						if buffer[position] != rune('P') {
							goto l1099
						}
						position++
					}
				l1106:
					depth--
					add(rulePegText, position1101)
				}
				if !_rules[ruleAction76]() {
					goto l1099
				}
				depth--
				add(ruleMap, position1100)
			}
			return true
		l1099:
			position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
			return false
		},
		/* 100 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action77)> */
		func() bool {
			position1108, tokenIndex1108, depth1108 := position, tokenIndex, depth
			{
				position1109 := position
				depth++
				{
					position1110 := position
					depth++
					{
						position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1112
						}
						position++
						goto l1111
					l1112:
						position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
						if buffer[position] != rune('O') {
							goto l1108
						}
						position++
					}
				l1111:
					{
						position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1114
						}
						position++
						goto l1113
					l1114:
						position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
						if buffer[position] != rune('R') {
							goto l1108
						}
						position++
					}
				l1113:
					depth--
					add(rulePegText, position1110)
				}
				if !_rules[ruleAction77]() {
					goto l1108
				}
				depth--
				add(ruleOr, position1109)
			}
			return true
		l1108:
			position, tokenIndex, depth = position1108, tokenIndex1108, depth1108
			return false
		},
		/* 101 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action78)> */
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
						if buffer[position] != rune('a') {
							goto l1119
						}
						position++
						goto l1118
					l1119:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						if buffer[position] != rune('A') {
							goto l1115
						}
						position++
					}
				l1118:
					{
						position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1121
						}
						position++
						goto l1120
					l1121:
						position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
						if buffer[position] != rune('N') {
							goto l1115
						}
						position++
					}
				l1120:
					{
						position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1123
						}
						position++
						goto l1122
					l1123:
						position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
						if buffer[position] != rune('D') {
							goto l1115
						}
						position++
					}
				l1122:
					depth--
					add(rulePegText, position1117)
				}
				if !_rules[ruleAction78]() {
					goto l1115
				}
				depth--
				add(ruleAnd, position1116)
			}
			return true
		l1115:
			position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
			return false
		},
		/* 102 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action79)> */
		func() bool {
			position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
			{
				position1125 := position
				depth++
				{
					position1126 := position
					depth++
					{
						position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1128
						}
						position++
						goto l1127
					l1128:
						position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
						if buffer[position] != rune('N') {
							goto l1124
						}
						position++
					}
				l1127:
					{
						position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1130
						}
						position++
						goto l1129
					l1130:
						position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
						if buffer[position] != rune('O') {
							goto l1124
						}
						position++
					}
				l1129:
					{
						position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1132
						}
						position++
						goto l1131
					l1132:
						position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
						if buffer[position] != rune('T') {
							goto l1124
						}
						position++
					}
				l1131:
					depth--
					add(rulePegText, position1126)
				}
				if !_rules[ruleAction79]() {
					goto l1124
				}
				depth--
				add(ruleNot, position1125)
			}
			return true
		l1124:
			position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
			return false
		},
		/* 103 Equal <- <(<'='> Action80)> */
		func() bool {
			position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
			{
				position1134 := position
				depth++
				{
					position1135 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1133
					}
					position++
					depth--
					add(rulePegText, position1135)
				}
				if !_rules[ruleAction80]() {
					goto l1133
				}
				depth--
				add(ruleEqual, position1134)
			}
			return true
		l1133:
			position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
			return false
		},
		/* 104 Less <- <(<'<'> Action81)> */
		func() bool {
			position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
			{
				position1137 := position
				depth++
				{
					position1138 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1136
					}
					position++
					depth--
					add(rulePegText, position1138)
				}
				if !_rules[ruleAction81]() {
					goto l1136
				}
				depth--
				add(ruleLess, position1137)
			}
			return true
		l1136:
			position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
			return false
		},
		/* 105 LessOrEqual <- <(<('<' '=')> Action82)> */
		func() bool {
			position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
			{
				position1140 := position
				depth++
				{
					position1141 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1139
					}
					position++
					if buffer[position] != rune('=') {
						goto l1139
					}
					position++
					depth--
					add(rulePegText, position1141)
				}
				if !_rules[ruleAction82]() {
					goto l1139
				}
				depth--
				add(ruleLessOrEqual, position1140)
			}
			return true
		l1139:
			position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
			return false
		},
		/* 106 Greater <- <(<'>'> Action83)> */
		func() bool {
			position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
			{
				position1143 := position
				depth++
				{
					position1144 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1142
					}
					position++
					depth--
					add(rulePegText, position1144)
				}
				if !_rules[ruleAction83]() {
					goto l1142
				}
				depth--
				add(ruleGreater, position1143)
			}
			return true
		l1142:
			position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
			return false
		},
		/* 107 GreaterOrEqual <- <(<('>' '=')> Action84)> */
		func() bool {
			position1145, tokenIndex1145, depth1145 := position, tokenIndex, depth
			{
				position1146 := position
				depth++
				{
					position1147 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1145
					}
					position++
					if buffer[position] != rune('=') {
						goto l1145
					}
					position++
					depth--
					add(rulePegText, position1147)
				}
				if !_rules[ruleAction84]() {
					goto l1145
				}
				depth--
				add(ruleGreaterOrEqual, position1146)
			}
			return true
		l1145:
			position, tokenIndex, depth = position1145, tokenIndex1145, depth1145
			return false
		},
		/* 108 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action85)> */
		func() bool {
			position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
			{
				position1149 := position
				depth++
				{
					position1150 := position
					depth++
					{
						position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1152
						}
						position++
						if buffer[position] != rune('=') {
							goto l1152
						}
						position++
						goto l1151
					l1152:
						position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
						if buffer[position] != rune('<') {
							goto l1148
						}
						position++
						if buffer[position] != rune('>') {
							goto l1148
						}
						position++
					}
				l1151:
					depth--
					add(rulePegText, position1150)
				}
				if !_rules[ruleAction85]() {
					goto l1148
				}
				depth--
				add(ruleNotEqual, position1149)
			}
			return true
		l1148:
			position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
			return false
		},
		/* 109 Concat <- <(<('|' '|')> Action86)> */
		func() bool {
			position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
			{
				position1154 := position
				depth++
				{
					position1155 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1153
					}
					position++
					if buffer[position] != rune('|') {
						goto l1153
					}
					position++
					depth--
					add(rulePegText, position1155)
				}
				if !_rules[ruleAction86]() {
					goto l1153
				}
				depth--
				add(ruleConcat, position1154)
			}
			return true
		l1153:
			position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
			return false
		},
		/* 110 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action87)> */
		func() bool {
			position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
			{
				position1157 := position
				depth++
				{
					position1158 := position
					depth++
					{
						position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1160
						}
						position++
						goto l1159
					l1160:
						position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
						if buffer[position] != rune('I') {
							goto l1156
						}
						position++
					}
				l1159:
					{
						position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1162
						}
						position++
						goto l1161
					l1162:
						position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
						if buffer[position] != rune('S') {
							goto l1156
						}
						position++
					}
				l1161:
					depth--
					add(rulePegText, position1158)
				}
				if !_rules[ruleAction87]() {
					goto l1156
				}
				depth--
				add(ruleIs, position1157)
			}
			return true
		l1156:
			position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
			return false
		},
		/* 111 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action88)> */
		func() bool {
			position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
			{
				position1164 := position
				depth++
				{
					position1165 := position
					depth++
					{
						position1166, tokenIndex1166, depth1166 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1167
						}
						position++
						goto l1166
					l1167:
						position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
						if buffer[position] != rune('I') {
							goto l1163
						}
						position++
					}
				l1166:
					{
						position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1169
						}
						position++
						goto l1168
					l1169:
						position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
						if buffer[position] != rune('S') {
							goto l1163
						}
						position++
					}
				l1168:
					if !_rules[rulesp]() {
						goto l1163
					}
					{
						position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1171
						}
						position++
						goto l1170
					l1171:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						if buffer[position] != rune('N') {
							goto l1163
						}
						position++
					}
				l1170:
					{
						position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1173
						}
						position++
						goto l1172
					l1173:
						position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
						if buffer[position] != rune('O') {
							goto l1163
						}
						position++
					}
				l1172:
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
							goto l1163
						}
						position++
					}
				l1174:
					depth--
					add(rulePegText, position1165)
				}
				if !_rules[ruleAction88]() {
					goto l1163
				}
				depth--
				add(ruleIsNot, position1164)
			}
			return true
		l1163:
			position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
			return false
		},
		/* 112 Plus <- <(<'+'> Action89)> */
		func() bool {
			position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
			{
				position1177 := position
				depth++
				{
					position1178 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1176
					}
					position++
					depth--
					add(rulePegText, position1178)
				}
				if !_rules[ruleAction89]() {
					goto l1176
				}
				depth--
				add(rulePlus, position1177)
			}
			return true
		l1176:
			position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
			return false
		},
		/* 113 Minus <- <(<'-'> Action90)> */
		func() bool {
			position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
			{
				position1180 := position
				depth++
				{
					position1181 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1179
					}
					position++
					depth--
					add(rulePegText, position1181)
				}
				if !_rules[ruleAction90]() {
					goto l1179
				}
				depth--
				add(ruleMinus, position1180)
			}
			return true
		l1179:
			position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
			return false
		},
		/* 114 Multiply <- <(<'*'> Action91)> */
		func() bool {
			position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
			{
				position1183 := position
				depth++
				{
					position1184 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1182
					}
					position++
					depth--
					add(rulePegText, position1184)
				}
				if !_rules[ruleAction91]() {
					goto l1182
				}
				depth--
				add(ruleMultiply, position1183)
			}
			return true
		l1182:
			position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
			return false
		},
		/* 115 Divide <- <(<'/'> Action92)> */
		func() bool {
			position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
			{
				position1186 := position
				depth++
				{
					position1187 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1185
					}
					position++
					depth--
					add(rulePegText, position1187)
				}
				if !_rules[ruleAction92]() {
					goto l1185
				}
				depth--
				add(ruleDivide, position1186)
			}
			return true
		l1185:
			position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
			return false
		},
		/* 116 Modulo <- <(<'%'> Action93)> */
		func() bool {
			position1188, tokenIndex1188, depth1188 := position, tokenIndex, depth
			{
				position1189 := position
				depth++
				{
					position1190 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1188
					}
					position++
					depth--
					add(rulePegText, position1190)
				}
				if !_rules[ruleAction93]() {
					goto l1188
				}
				depth--
				add(ruleModulo, position1189)
			}
			return true
		l1188:
			position, tokenIndex, depth = position1188, tokenIndex1188, depth1188
			return false
		},
		/* 117 UnaryMinus <- <(<'-'> Action94)> */
		func() bool {
			position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
			{
				position1192 := position
				depth++
				{
					position1193 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1191
					}
					position++
					depth--
					add(rulePegText, position1193)
				}
				if !_rules[ruleAction94]() {
					goto l1191
				}
				depth--
				add(ruleUnaryMinus, position1192)
			}
			return true
		l1191:
			position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
			return false
		},
		/* 118 Identifier <- <(<ident> Action95)> */
		func() bool {
			position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
			{
				position1195 := position
				depth++
				{
					position1196 := position
					depth++
					if !_rules[ruleident]() {
						goto l1194
					}
					depth--
					add(rulePegText, position1196)
				}
				if !_rules[ruleAction95]() {
					goto l1194
				}
				depth--
				add(ruleIdentifier, position1195)
			}
			return true
		l1194:
			position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
			return false
		},
		/* 119 TargetIdentifier <- <(<jsonPath> Action96)> */
		func() bool {
			position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
			{
				position1198 := position
				depth++
				{
					position1199 := position
					depth++
					if !_rules[rulejsonPath]() {
						goto l1197
					}
					depth--
					add(rulePegText, position1199)
				}
				if !_rules[ruleAction96]() {
					goto l1197
				}
				depth--
				add(ruleTargetIdentifier, position1198)
			}
			return true
		l1197:
			position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
			return false
		},
		/* 120 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
			{
				position1201 := position
				depth++
				{
					position1202, tokenIndex1202, depth1202 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1203
					}
					position++
					goto l1202
				l1203:
					position, tokenIndex, depth = position1202, tokenIndex1202, depth1202
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1200
					}
					position++
				}
			l1202:
			l1204:
				{
					position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
					{
						position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1207
						}
						position++
						goto l1206
					l1207:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1208
						}
						position++
						goto l1206
					l1208:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1209
						}
						position++
						goto l1206
					l1209:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						if buffer[position] != rune('_') {
							goto l1205
						}
						position++
					}
				l1206:
					goto l1204
				l1205:
					position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
				}
				depth--
				add(ruleident, position1201)
			}
			return true
		l1200:
			position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
			return false
		},
		/* 121 jsonPath <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_' / '.' / '[' / ']' / '"')*)> */
		func() bool {
			position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
			{
				position1211 := position
				depth++
				{
					position1212, tokenIndex1212, depth1212 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1213
					}
					position++
					goto l1212
				l1213:
					position, tokenIndex, depth = position1212, tokenIndex1212, depth1212
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1210
					}
					position++
				}
			l1212:
			l1214:
				{
					position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
					{
						position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1217
						}
						position++
						goto l1216
					l1217:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1218
						}
						position++
						goto l1216
					l1218:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1219
						}
						position++
						goto l1216
					l1219:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if buffer[position] != rune('_') {
							goto l1220
						}
						position++
						goto l1216
					l1220:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if buffer[position] != rune('.') {
							goto l1221
						}
						position++
						goto l1216
					l1221:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if buffer[position] != rune('[') {
							goto l1222
						}
						position++
						goto l1216
					l1222:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if buffer[position] != rune(']') {
							goto l1223
						}
						position++
						goto l1216
					l1223:
						position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
						if buffer[position] != rune('"') {
							goto l1215
						}
						position++
					}
				l1216:
					goto l1214
				l1215:
					position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
				}
				depth--
				add(rulejsonPath, position1211)
			}
			return true
		l1210:
			position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
			return false
		},
		/* 122 sp <- <(' ' / '\t' / '\n' / '\r' / comment)*> */
		func() bool {
			{
				position1225 := position
				depth++
			l1226:
				{
					position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
					{
						position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l1229
						}
						position++
						goto l1228
					l1229:
						position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
						if buffer[position] != rune('\t') {
							goto l1230
						}
						position++
						goto l1228
					l1230:
						position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
						if buffer[position] != rune('\n') {
							goto l1231
						}
						position++
						goto l1228
					l1231:
						position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
						if buffer[position] != rune('\r') {
							goto l1232
						}
						position++
						goto l1228
					l1232:
						position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
						if !_rules[rulecomment]() {
							goto l1227
						}
					}
				l1228:
					goto l1226
				l1227:
					position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
				}
				depth--
				add(rulesp, position1225)
			}
			return true
		},
		/* 123 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
			{
				position1234 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1233
				}
				position++
				if buffer[position] != rune('-') {
					goto l1233
				}
				position++
			l1235:
				{
					position1236, tokenIndex1236, depth1236 := position, tokenIndex, depth
					{
						position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
						{
							position1238, tokenIndex1238, depth1238 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1239
							}
							position++
							goto l1238
						l1239:
							position, tokenIndex, depth = position1238, tokenIndex1238, depth1238
							if buffer[position] != rune('\n') {
								goto l1237
							}
							position++
						}
					l1238:
						goto l1236
					l1237:
						position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
					}
					if !matchDot() {
						goto l1236
					}
					goto l1235
				l1236:
					position, tokenIndex, depth = position1236, tokenIndex1236, depth1236
				}
				{
					position1240, tokenIndex1240, depth1240 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1241
					}
					position++
					goto l1240
				l1241:
					position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
					if buffer[position] != rune('\n') {
						goto l1233
					}
					position++
				}
			l1240:
				depth--
				add(rulecomment, position1234)
			}
			return true
		l1233:
			position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
			return false
		},
		/* 125 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 127 Action1 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 128 Action2 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 129 Action3 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 130 Action4 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 131 Action5 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 132 Action6 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 133 Action7 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 134 Action8 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 135 Action9 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 136 Action10 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 137 Action11 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 138 Action12 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 139 Action13 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 140 Action14 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 141 Action15 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 142 Action16 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 143 Action17 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 144 Action18 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 145 Action19 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 146 Action20 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 147 Action21 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 148 Action22 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 149 Action23 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 150 Action24 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 151 Action25 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 152 Action26 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 153 Action27 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 154 Action28 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 155 Action29 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 156 Action30 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 157 Action31 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 158 Action32 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 159 Action33 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 160 Action34 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 161 Action35 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 162 Action36 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 163 Action37 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 164 Action38 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 165 Action39 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 166 Action40 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 167 Action41 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 168 Action42 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 169 Action43 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 170 Action44 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 171 Action45 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 172 Action46 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 173 Action47 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 174 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 175 Action49 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 176 Action50 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 177 Action51 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 178 Action52 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 179 Action53 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 180 Action54 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 181 Action55 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 182 Action56 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 183 Action57 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 184 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 185 Action59 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 186 Action60 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 187 Action61 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 188 Action62 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 189 Action63 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 190 Action64 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 191 Action65 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 192 Action66 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 193 Action67 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 194 Action68 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 195 Action69 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 196 Action70 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 197 Action71 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 198 Action72 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 199 Action73 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 200 Action74 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 201 Action75 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 202 Action76 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 203 Action77 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 204 Action78 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 205 Action79 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 206 Action80 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 207 Action81 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 208 Action82 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 209 Action83 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 210 Action84 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 211 Action85 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 212 Action86 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 213 Action87 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 214 Action88 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 215 Action89 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 216 Action90 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 217 Action91 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 218 Action92 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 219 Action93 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 220 Action94 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 221 Action95 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 222 Action96 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
	}
	p.rules = _rules
}

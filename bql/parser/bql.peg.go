package parser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 4

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleStatement
	ruleSelectStmt
	ruleCreateStreamAsSelectStmt
	ruleCreateSourceStmt
	ruleCreateSinkStmt
	ruleCreateStreamFromSourceStmt
	ruleCreateStreamFromSourceExtStmt
	ruleInsertIntoSelectStmt
	ruleEmitter
	ruleProjections
	ruleProjection
	ruleAliasExpression
	ruleWindowedFrom
	ruleDefWindowedFrom
	ruleInterval
	ruleRelations
	ruleDefRelations
	ruleFilter
	ruleGrouping
	ruleGroupList
	ruleHaving
	ruleRelationLike
	ruleDefRelationLike
	ruleAliasedStreamWindow
	ruleDefAliasedStreamWindow
	ruleStreamWindow
	ruleDefStreamWindow
	ruleSourceSinkSpecs
	ruleSourceSinkParam
	ruleExpression
	ruleorExpr
	ruleandExpr
	rulecomparisonExpr
	ruletermExpr
	ruleproductExpr
	rulebaseExpr
	ruleFuncApp
	ruleFuncParams
	ruleLiteral
	ruleComparisonOp
	rulePlusMinusOp
	ruleMultDivOp
	ruleStream
	ruleRowValue
	ruleNumericLiteral
	ruleFloatLiteral
	ruleFunction
	ruleBooleanLiteral
	ruleTRUE
	ruleFALSE
	ruleWildcard
	ruleStringLiteral
	ruleISTREAM
	ruleDSTREAM
	ruleRSTREAM
	ruleIntervalUnit
	ruleTUPLES
	ruleSECONDS
	ruleStreamIdentifier
	ruleSourceSinkType
	ruleSourceSinkParamKey
	ruleSourceSinkParamVal
	ruleOr
	ruleAnd
	ruleEqual
	ruleLess
	ruleLessOrEqual
	ruleGreater
	ruleGreaterOrEqual
	ruleNotEqual
	rulePlus
	ruleMinus
	ruleMultiply
	ruleDivide
	ruleModulo
	ruleIdentifier
	ruleident
	rulesp
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	rulePegText
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

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Statement",
	"SelectStmt",
	"CreateStreamAsSelectStmt",
	"CreateSourceStmt",
	"CreateSinkStmt",
	"CreateStreamFromSourceStmt",
	"CreateStreamFromSourceExtStmt",
	"InsertIntoSelectStmt",
	"Emitter",
	"Projections",
	"Projection",
	"AliasExpression",
	"WindowedFrom",
	"DefWindowedFrom",
	"Interval",
	"Relations",
	"DefRelations",
	"Filter",
	"Grouping",
	"GroupList",
	"Having",
	"RelationLike",
	"DefRelationLike",
	"AliasedStreamWindow",
	"DefAliasedStreamWindow",
	"StreamWindow",
	"DefStreamWindow",
	"SourceSinkSpecs",
	"SourceSinkParam",
	"Expression",
	"orExpr",
	"andExpr",
	"comparisonExpr",
	"termExpr",
	"productExpr",
	"baseExpr",
	"FuncApp",
	"FuncParams",
	"Literal",
	"ComparisonOp",
	"PlusMinusOp",
	"MultDivOp",
	"Stream",
	"RowValue",
	"NumericLiteral",
	"FloatLiteral",
	"Function",
	"BooleanLiteral",
	"TRUE",
	"FALSE",
	"Wildcard",
	"StringLiteral",
	"ISTREAM",
	"DSTREAM",
	"RSTREAM",
	"IntervalUnit",
	"TUPLES",
	"SECONDS",
	"StreamIdentifier",
	"SourceSinkType",
	"SourceSinkParamKey",
	"SourceSinkParamVal",
	"Or",
	"And",
	"Equal",
	"Less",
	"LessOrEqual",
	"Greater",
	"GreaterOrEqual",
	"NotEqual",
	"Plus",
	"Minus",
	"Multiply",
	"Divide",
	"Modulo",
	"Identifier",
	"ident",
	"sp",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"PegText",
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

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next, depth int)
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
type token16 struct {
	pegRule
	begin, end, next int16
}

func (t *token16) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token16) isParentOf(u token16) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token16) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token16) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens16 struct {
	tree    []token16
	ordered [][]token16
}

func (t *tokens16) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens16) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens16) Order() [][]token16 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int16, 1, math.MaxInt16)
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

	ordered, pool := make([][]token16, len(depths)), make([]token16, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int16(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state16 struct {
	token16
	depths []int16
	leaf   bool
}

func (t *tokens16) AST() *node32 {
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

func (t *tokens16) PreOrder() (<-chan state16, [][]token16) {
	s, ordered := make(chan state16, 6), t.Order()
	go func() {
		var states [8]state16
		for i, _ := range states {
			states[i].depths = make([]int16, len(ordered))
		}
		depths, state, depth := make([]int16, len(ordered)), 0, 1
		write := func(t token16, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int16(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token16 = ordered[0][0]
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
							write(token16{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token16{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
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
					write(token16{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
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

func (t *tokens16) PrintSyntax() {
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

func (t *tokens16) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens16) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token16{pegRule: rule, begin: int16(begin), end: int16(end), next: int16(depth)}
}

func (t *tokens16) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens16) Error() []token32 {
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

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next int32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
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
		token.next = int32(i)
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
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int32(depth), leaf
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

func (t *tokens32) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token32{pegRule: rule, begin: int32(begin), end: int32(end), next: int32(depth)}
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

func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}

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
	rules  [143]func() bool
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
	buffer, begin, end := p.Buffer, 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)

		case ruleAction0:

			p.AssembleSelect()

		case ruleAction1:

			p.AssembleCreateStreamAsSelect()

		case ruleAction2:

			p.AssembleCreateSource()

		case ruleAction3:

			p.AssembleCreateSink()

		case ruleAction4:

			p.AssembleCreateStreamFromSource()

		case ruleAction5:

			p.AssembleCreateStreamFromSourceExt()

		case ruleAction6:

			p.AssembleInsertIntoSelect()

		case ruleAction7:

			p.AssembleEmitter()

		case ruleAction8:

			p.AssembleProjections(begin, end)

		case ruleAction9:

			p.AssembleAlias()

		case ruleAction10:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction11:

			p.AssembleWindowedFrom(begin, end)

		case ruleAction12:

			p.AssembleInterval()

		case ruleAction13:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction14:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction15:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction16:

			p.EnsureAliasedStreamWindow()

		case ruleAction17:

			p.EnsureAliasedStreamWindow()

		case ruleAction18:

			p.AssembleAliasedStreamWindow()

		case ruleAction19:

			p.AssembleAliasedStreamWindow()

		case ruleAction20:

			p.AssembleStreamWindow()

		case ruleAction21:

			p.AssembleStreamWindow()

		case ruleAction22:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction23:

			p.AssembleSourceSinkParam()

		case ruleAction24:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction25:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction26:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction27:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction28:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction29:

			p.AssembleFuncApp()

		case ruleAction30:

			p.AssembleExpressions(begin, end)

		case ruleAction31:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction32:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction33:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction34:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction35:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction36:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction37:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction38:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction40:

			p.PushComponent(begin, end, Istream)

		case ruleAction41:

			p.PushComponent(begin, end, Dstream)

		case ruleAction42:

			p.PushComponent(begin, end, Rstream)

		case ruleAction43:

			p.PushComponent(begin, end, Tuples)

		case ruleAction44:

			p.PushComponent(begin, end, Seconds)

		case ruleAction45:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction46:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction47:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction48:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamVal(substr))

		case ruleAction49:

			p.PushComponent(begin, end, Or)

		case ruleAction50:

			p.PushComponent(begin, end, And)

		case ruleAction51:

			p.PushComponent(begin, end, Equal)

		case ruleAction52:

			p.PushComponent(begin, end, Less)

		case ruleAction53:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction54:

			p.PushComponent(begin, end, Greater)

		case ruleAction55:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction56:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction57:

			p.PushComponent(begin, end, Plus)

		case ruleAction58:

			p.PushComponent(begin, end, Minus)

		case ruleAction59:

			p.PushComponent(begin, end, Multiply)

		case ruleAction60:

			p.PushComponent(begin, end, Divide)

		case ruleAction61:

			p.PushComponent(begin, end, Modulo)

		case ruleAction62:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		}
	}
	_, _, _ = buffer, begin, end
}

func (p *bqlPeg) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens16{tree: make([]token16, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := 0, 0, 0, p.buffer, p.rules

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

	add := func(rule pegRule, begin int) {
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
		/* 0 Statement <- <((SelectStmt / CreateStreamAsSelectStmt / CreateSourceStmt / CreateStreamFromSourceStmt / CreateStreamFromSourceExtStmt / CreateSinkStmt / InsertIntoSelectStmt) !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[ruleSelectStmt]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l4
					}
					goto l2
				l4:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleCreateSourceStmt]() {
						goto l5
					}
					goto l2
				l5:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleCreateStreamFromSourceStmt]() {
						goto l6
					}
					goto l2
				l6:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleCreateStreamFromSourceExtStmt]() {
						goto l7
					}
					goto l2
				l7:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleCreateSinkStmt]() {
						goto l8
					}
					goto l2
				l8:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleInsertIntoSelectStmt]() {
						goto l0
					}
				}
			l2:
				{
					position9, tokenIndex9, depth9 := position, tokenIndex, depth
					if !matchDot() {
						goto l9
					}
					goto l0
				l9:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
				}
				depth--
				add(ruleStatement, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Projections sp DefWindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
		func() bool {
			position10, tokenIndex10, depth10 := position, tokenIndex, depth
			{
				position11 := position
				depth++
				{
					position12, tokenIndex12, depth12 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l13
					}
					position++
					goto l12
				l13:
					position, tokenIndex, depth = position12, tokenIndex12, depth12
					if buffer[position] != rune('S') {
						goto l10
					}
					position++
				}
			l12:
				{
					position14, tokenIndex14, depth14 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l15
					}
					position++
					goto l14
				l15:
					position, tokenIndex, depth = position14, tokenIndex14, depth14
					if buffer[position] != rune('E') {
						goto l10
					}
					position++
				}
			l14:
				{
					position16, tokenIndex16, depth16 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l17
					}
					position++
					goto l16
				l17:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
					if buffer[position] != rune('L') {
						goto l10
					}
					position++
				}
			l16:
				{
					position18, tokenIndex18, depth18 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l19
					}
					position++
					goto l18
				l19:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					if buffer[position] != rune('E') {
						goto l10
					}
					position++
				}
			l18:
				{
					position20, tokenIndex20, depth20 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l21
					}
					position++
					goto l20
				l21:
					position, tokenIndex, depth = position20, tokenIndex20, depth20
					if buffer[position] != rune('C') {
						goto l10
					}
					position++
				}
			l20:
				{
					position22, tokenIndex22, depth22 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l23
					}
					position++
					goto l22
				l23:
					position, tokenIndex, depth = position22, tokenIndex22, depth22
					if buffer[position] != rune('T') {
						goto l10
					}
					position++
				}
			l22:
				if !_rules[rulesp]() {
					goto l10
				}
				if !_rules[ruleProjections]() {
					goto l10
				}
				if !_rules[rulesp]() {
					goto l10
				}
				if !_rules[ruleDefWindowedFrom]() {
					goto l10
				}
				if !_rules[rulesp]() {
					goto l10
				}
				if !_rules[ruleFilter]() {
					goto l10
				}
				if !_rules[rulesp]() {
					goto l10
				}
				if !_rules[ruleGrouping]() {
					goto l10
				}
				if !_rules[rulesp]() {
					goto l10
				}
				if !_rules[ruleHaving]() {
					goto l10
				}
				if !_rules[rulesp]() {
					goto l10
				}
				if !_rules[ruleAction0]() {
					goto l10
				}
				depth--
				add(ruleSelectStmt, position11)
			}
			return true
		l10:
			position, tokenIndex, depth = position10, tokenIndex10, depth10
			return false
		},
		/* 2 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp Emitter sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
		func() bool {
			position24, tokenIndex24, depth24 := position, tokenIndex, depth
			{
				position25 := position
				depth++
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l27
					}
					position++
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					if buffer[position] != rune('C') {
						goto l24
					}
					position++
				}
			l26:
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l29
					}
					position++
					goto l28
				l29:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
					if buffer[position] != rune('R') {
						goto l24
					}
					position++
				}
			l28:
				{
					position30, tokenIndex30, depth30 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l31
					}
					position++
					goto l30
				l31:
					position, tokenIndex, depth = position30, tokenIndex30, depth30
					if buffer[position] != rune('E') {
						goto l24
					}
					position++
				}
			l30:
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l33
					}
					position++
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if buffer[position] != rune('A') {
						goto l24
					}
					position++
				}
			l32:
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l35
					}
					position++
					goto l34
				l35:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
					if buffer[position] != rune('T') {
						goto l24
					}
					position++
				}
			l34:
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l37
					}
					position++
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if buffer[position] != rune('E') {
						goto l24
					}
					position++
				}
			l36:
				if !_rules[rulesp]() {
					goto l24
				}
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l39
					}
					position++
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if buffer[position] != rune('S') {
						goto l24
					}
					position++
				}
			l38:
				{
					position40, tokenIndex40, depth40 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l41
					}
					position++
					goto l40
				l41:
					position, tokenIndex, depth = position40, tokenIndex40, depth40
					if buffer[position] != rune('T') {
						goto l24
					}
					position++
				}
			l40:
				{
					position42, tokenIndex42, depth42 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l43
					}
					position++
					goto l42
				l43:
					position, tokenIndex, depth = position42, tokenIndex42, depth42
					if buffer[position] != rune('R') {
						goto l24
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
						goto l24
					}
					position++
				}
			l44:
				{
					position46, tokenIndex46, depth46 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex, depth = position46, tokenIndex46, depth46
					if buffer[position] != rune('A') {
						goto l24
					}
					position++
				}
			l46:
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l49
					}
					position++
					goto l48
				l49:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
					if buffer[position] != rune('M') {
						goto l24
					}
					position++
				}
			l48:
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				{
					position50, tokenIndex50, depth50 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l51
					}
					position++
					goto l50
				l51:
					position, tokenIndex, depth = position50, tokenIndex50, depth50
					if buffer[position] != rune('A') {
						goto l24
					}
					position++
				}
			l50:
				{
					position52, tokenIndex52, depth52 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex, depth = position52, tokenIndex52, depth52
					if buffer[position] != rune('S') {
						goto l24
					}
					position++
				}
			l52:
				if !_rules[rulesp]() {
					goto l24
				}
				{
					position54, tokenIndex54, depth54 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l55
					}
					position++
					goto l54
				l55:
					position, tokenIndex, depth = position54, tokenIndex54, depth54
					if buffer[position] != rune('S') {
						goto l24
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
						goto l24
					}
					position++
				}
			l56:
				{
					position58, tokenIndex58, depth58 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex, depth = position58, tokenIndex58, depth58
					if buffer[position] != rune('L') {
						goto l24
					}
					position++
				}
			l58:
				{
					position60, tokenIndex60, depth60 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l61
					}
					position++
					goto l60
				l61:
					position, tokenIndex, depth = position60, tokenIndex60, depth60
					if buffer[position] != rune('E') {
						goto l24
					}
					position++
				}
			l60:
				{
					position62, tokenIndex62, depth62 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex, depth = position62, tokenIndex62, depth62
					if buffer[position] != rune('C') {
						goto l24
					}
					position++
				}
			l62:
				{
					position64, tokenIndex64, depth64 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l65
					}
					position++
					goto l64
				l65:
					position, tokenIndex, depth = position64, tokenIndex64, depth64
					if buffer[position] != rune('T') {
						goto l24
					}
					position++
				}
			l64:
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleEmitter]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleProjections]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleWindowedFrom]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleFilter]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleGrouping]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleHaving]() {
					goto l24
				}
				if !_rules[rulesp]() {
					goto l24
				}
				if !_rules[ruleAction1]() {
					goto l24
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position25)
			}
			return true
		l24:
			position, tokenIndex, depth = position24, tokenIndex24, depth24
			return false
		},
		/* 3 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
		func() bool {
			position66, tokenIndex66, depth66 := position, tokenIndex, depth
			{
				position67 := position
				depth++
				{
					position68, tokenIndex68, depth68 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l69
					}
					position++
					goto l68
				l69:
					position, tokenIndex, depth = position68, tokenIndex68, depth68
					if buffer[position] != rune('C') {
						goto l66
					}
					position++
				}
			l68:
				{
					position70, tokenIndex70, depth70 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l71
					}
					position++
					goto l70
				l71:
					position, tokenIndex, depth = position70, tokenIndex70, depth70
					if buffer[position] != rune('R') {
						goto l66
					}
					position++
				}
			l70:
				{
					position72, tokenIndex72, depth72 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l73
					}
					position++
					goto l72
				l73:
					position, tokenIndex, depth = position72, tokenIndex72, depth72
					if buffer[position] != rune('E') {
						goto l66
					}
					position++
				}
			l72:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					if buffer[position] != rune('A') {
						goto l66
					}
					position++
				}
			l74:
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l77
					}
					position++
					goto l76
				l77:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
					if buffer[position] != rune('T') {
						goto l66
					}
					position++
				}
			l76:
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l79
					}
					position++
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if buffer[position] != rune('E') {
						goto l66
					}
					position++
				}
			l78:
				if !_rules[rulesp]() {
					goto l66
				}
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
					if buffer[position] != rune('S') {
						goto l66
					}
					position++
				}
			l80:
				{
					position82, tokenIndex82, depth82 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l83
					}
					position++
					goto l82
				l83:
					position, tokenIndex, depth = position82, tokenIndex82, depth82
					if buffer[position] != rune('O') {
						goto l66
					}
					position++
				}
			l82:
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
						goto l66
					}
					position++
				}
			l84:
				{
					position86, tokenIndex86, depth86 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l87
					}
					position++
					goto l86
				l87:
					position, tokenIndex, depth = position86, tokenIndex86, depth86
					if buffer[position] != rune('R') {
						goto l66
					}
					position++
				}
			l86:
				{
					position88, tokenIndex88, depth88 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l89
					}
					position++
					goto l88
				l89:
					position, tokenIndex, depth = position88, tokenIndex88, depth88
					if buffer[position] != rune('C') {
						goto l66
					}
					position++
				}
			l88:
				{
					position90, tokenIndex90, depth90 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l91
					}
					position++
					goto l90
				l91:
					position, tokenIndex, depth = position90, tokenIndex90, depth90
					if buffer[position] != rune('E') {
						goto l66
					}
					position++
				}
			l90:
				if !_rules[rulesp]() {
					goto l66
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l66
				}
				if !_rules[rulesp]() {
					goto l66
				}
				{
					position92, tokenIndex92, depth92 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l93
					}
					position++
					goto l92
				l93:
					position, tokenIndex, depth = position92, tokenIndex92, depth92
					if buffer[position] != rune('T') {
						goto l66
					}
					position++
				}
			l92:
				{
					position94, tokenIndex94, depth94 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l95
					}
					position++
					goto l94
				l95:
					position, tokenIndex, depth = position94, tokenIndex94, depth94
					if buffer[position] != rune('Y') {
						goto l66
					}
					position++
				}
			l94:
				{
					position96, tokenIndex96, depth96 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l97
					}
					position++
					goto l96
				l97:
					position, tokenIndex, depth = position96, tokenIndex96, depth96
					if buffer[position] != rune('P') {
						goto l66
					}
					position++
				}
			l96:
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l99
					}
					position++
					goto l98
				l99:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
					if buffer[position] != rune('E') {
						goto l66
					}
					position++
				}
			l98:
				if !_rules[rulesp]() {
					goto l66
				}
				if !_rules[ruleSourceSinkType]() {
					goto l66
				}
				if !_rules[rulesp]() {
					goto l66
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l66
				}
				if !_rules[ruleAction2]() {
					goto l66
				}
				depth--
				add(ruleCreateSourceStmt, position67)
			}
			return true
		l66:
			position, tokenIndex, depth = position66, tokenIndex66, depth66
			return false
		},
		/* 4 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
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
					if buffer[position] != rune('i') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('I') {
						goto l100
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('N') {
						goto l100
					}
					position++
				}
			l118:
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
					if buffer[position] != rune('K') {
						goto l100
					}
					position++
				}
			l120:
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
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('T') {
						goto l100
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('Y') {
						goto l100
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('P') {
						goto l100
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('E') {
						goto l100
					}
					position++
				}
			l128:
				if !_rules[rulesp]() {
					goto l100
				}
				if !_rules[ruleSourceSinkType]() {
					goto l100
				}
				if !_rules[rulesp]() {
					goto l100
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l100
				}
				if !_rules[ruleAction3]() {
					goto l100
				}
				depth--
				add(ruleCreateSinkStmt, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 5 CreateStreamFromSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action4)> */
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
					if buffer[position] != rune('f') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('F') {
						goto l130
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
						goto l130
					}
					position++
				}
			l158:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('O') {
						goto l130
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('M') {
						goto l130
					}
					position++
				}
			l162:
				if !_rules[rulesp]() {
					goto l130
				}
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('S') {
						goto l130
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('O') {
						goto l130
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('U') {
						goto l130
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('R') {
						goto l130
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('C') {
						goto l130
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('E') {
						goto l130
					}
					position++
				}
			l174:
				if !_rules[rulesp]() {
					goto l130
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l130
				}
				if !_rules[ruleAction4]() {
					goto l130
				}
				depth--
				add(ruleCreateStreamFromSourceStmt, position131)
			}
			return true
		l130:
			position, tokenIndex, depth = position130, tokenIndex130, depth130
			return false
		},
		/* 6 CreateStreamFromSourceExtStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp SourceSinkType sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp SourceSinkSpecs Action5)> */
		func() bool {
			position176, tokenIndex176, depth176 := position, tokenIndex, depth
			{
				position177 := position
				depth++
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('C') {
						goto l176
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
						goto l176
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
						goto l176
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('A') {
						goto l176
					}
					position++
				}
			l184:
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
						goto l176
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('E') {
						goto l176
					}
					position++
				}
			l188:
				if !_rules[rulesp]() {
					goto l176
				}
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('S') {
						goto l176
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('T') {
						goto l176
					}
					position++
				}
			l192:
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('R') {
						goto l176
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
						goto l176
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('A') {
						goto l176
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('M') {
						goto l176
					}
					position++
				}
			l200:
				if !_rules[rulesp]() {
					goto l176
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l176
				}
				if !_rules[rulesp]() {
					goto l176
				}
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('F') {
						goto l176
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('R') {
						goto l176
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('O') {
						goto l176
					}
					position++
				}
			l206:
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('M') {
						goto l176
					}
					position++
				}
			l208:
				if !_rules[rulesp]() {
					goto l176
				}
				if !_rules[ruleSourceSinkType]() {
					goto l176
				}
				if !_rules[rulesp]() {
					goto l176
				}
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('S') {
						goto l176
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('O') {
						goto l176
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('U') {
						goto l176
					}
					position++
				}
			l214:
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('R') {
						goto l176
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('C') {
						goto l176
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('E') {
						goto l176
					}
					position++
				}
			l220:
				if !_rules[rulesp]() {
					goto l176
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l176
				}
				if !_rules[ruleAction5]() {
					goto l176
				}
				depth--
				add(ruleCreateStreamFromSourceExtStmt, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp SelectStmt Action6)> */
		func() bool {
			position222, tokenIndex222, depth222 := position, tokenIndex, depth
			{
				position223 := position
				depth++
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('I') {
						goto l222
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('N') {
						goto l222
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('S') {
						goto l222
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
						goto l222
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('R') {
						goto l222
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
						goto l222
					}
					position++
				}
			l234:
				if !_rules[rulesp]() {
					goto l222
				}
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('I') {
						goto l222
					}
					position++
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('N') {
						goto l222
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
						goto l222
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('O') {
						goto l222
					}
					position++
				}
			l242:
				if !_rules[rulesp]() {
					goto l222
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l222
				}
				if !_rules[rulesp]() {
					goto l222
				}
				if !_rules[ruleSelectStmt]() {
					goto l222
				}
				if !_rules[ruleAction6]() {
					goto l222
				}
				depth--
				add(ruleInsertIntoSelectStmt, position223)
			}
			return true
		l222:
			position, tokenIndex, depth = position222, tokenIndex222, depth222
			return false
		},
		/* 8 Emitter <- <((ISTREAM / DSTREAM / RSTREAM) Action7)> */
		func() bool {
			position244, tokenIndex244, depth244 := position, tokenIndex, depth
			{
				position245 := position
				depth++
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l247
					}
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if !_rules[ruleDSTREAM]() {
						goto l248
					}
					goto l246
				l248:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if !_rules[ruleRSTREAM]() {
						goto l244
					}
				}
			l246:
				if !_rules[ruleAction7]() {
					goto l244
				}
				depth--
				add(ruleEmitter, position245)
			}
			return true
		l244:
			position, tokenIndex, depth = position244, tokenIndex244, depth244
			return false
		},
		/* 9 Projections <- <(<(Projection sp (',' sp Projection)*)> Action8)> */
		func() bool {
			position249, tokenIndex249, depth249 := position, tokenIndex, depth
			{
				position250 := position
				depth++
				{
					position251 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l249
					}
					if !_rules[rulesp]() {
						goto l249
					}
				l252:
					{
						position253, tokenIndex253, depth253 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l253
						}
						position++
						if !_rules[rulesp]() {
							goto l253
						}
						if !_rules[ruleProjection]() {
							goto l253
						}
						goto l252
					l253:
						position, tokenIndex, depth = position253, tokenIndex253, depth253
					}
					depth--
					add(rulePegText, position251)
				}
				if !_rules[ruleAction8]() {
					goto l249
				}
				depth--
				add(ruleProjections, position250)
			}
			return true
		l249:
			position, tokenIndex, depth = position249, tokenIndex249, depth249
			return false
		},
		/* 10 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position254, tokenIndex254, depth254 := position, tokenIndex, depth
			{
				position255 := position
				depth++
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l257
					}
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if !_rules[ruleExpression]() {
						goto l258
					}
					goto l256
				l258:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if !_rules[ruleWildcard]() {
						goto l254
					}
				}
			l256:
				depth--
				add(ruleProjection, position255)
			}
			return true
		l254:
			position, tokenIndex, depth = position254, tokenIndex254, depth254
			return false
		},
		/* 11 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action9)> */
		func() bool {
			position259, tokenIndex259, depth259 := position, tokenIndex, depth
			{
				position260 := position
				depth++
				{
					position261, tokenIndex261, depth261 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l262
					}
					goto l261
				l262:
					position, tokenIndex, depth = position261, tokenIndex261, depth261
					if !_rules[ruleWildcard]() {
						goto l259
					}
				}
			l261:
				if !_rules[rulesp]() {
					goto l259
				}
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
						goto l259
					}
					position++
				}
			l263:
				{
					position265, tokenIndex265, depth265 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l266
					}
					position++
					goto l265
				l266:
					position, tokenIndex, depth = position265, tokenIndex265, depth265
					if buffer[position] != rune('S') {
						goto l259
					}
					position++
				}
			l265:
				if !_rules[rulesp]() {
					goto l259
				}
				if !_rules[ruleIdentifier]() {
					goto l259
				}
				if !_rules[ruleAction9]() {
					goto l259
				}
				depth--
				add(ruleAliasExpression, position260)
			}
			return true
		l259:
			position, tokenIndex, depth = position259, tokenIndex259, depth259
			return false
		},
		/* 12 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp)?> Action10)> */
		func() bool {
			position267, tokenIndex267, depth267 := position, tokenIndex, depth
			{
				position268 := position
				depth++
				{
					position269 := position
					depth++
					{
						position270, tokenIndex270, depth270 := position, tokenIndex, depth
						{
							position272, tokenIndex272, depth272 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l273
							}
							position++
							goto l272
						l273:
							position, tokenIndex, depth = position272, tokenIndex272, depth272
							if buffer[position] != rune('F') {
								goto l270
							}
							position++
						}
					l272:
						{
							position274, tokenIndex274, depth274 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l275
							}
							position++
							goto l274
						l275:
							position, tokenIndex, depth = position274, tokenIndex274, depth274
							if buffer[position] != rune('R') {
								goto l270
							}
							position++
						}
					l274:
						{
							position276, tokenIndex276, depth276 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l277
							}
							position++
							goto l276
						l277:
							position, tokenIndex, depth = position276, tokenIndex276, depth276
							if buffer[position] != rune('O') {
								goto l270
							}
							position++
						}
					l276:
						{
							position278, tokenIndex278, depth278 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l279
							}
							position++
							goto l278
						l279:
							position, tokenIndex, depth = position278, tokenIndex278, depth278
							if buffer[position] != rune('M') {
								goto l270
							}
							position++
						}
					l278:
						if !_rules[rulesp]() {
							goto l270
						}
						if !_rules[ruleRelations]() {
							goto l270
						}
						if !_rules[rulesp]() {
							goto l270
						}
						goto l271
					l270:
						position, tokenIndex, depth = position270, tokenIndex270, depth270
					}
				l271:
					depth--
					add(rulePegText, position269)
				}
				if !_rules[ruleAction10]() {
					goto l267
				}
				depth--
				add(ruleWindowedFrom, position268)
			}
			return true
		l267:
			position, tokenIndex, depth = position267, tokenIndex267, depth267
			return false
		},
		/* 13 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefRelations sp)?> Action11)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				{
					position282 := position
					depth++
					{
						position283, tokenIndex283, depth283 := position, tokenIndex, depth
						{
							position285, tokenIndex285, depth285 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l286
							}
							position++
							goto l285
						l286:
							position, tokenIndex, depth = position285, tokenIndex285, depth285
							if buffer[position] != rune('F') {
								goto l283
							}
							position++
						}
					l285:
						{
							position287, tokenIndex287, depth287 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l288
							}
							position++
							goto l287
						l288:
							position, tokenIndex, depth = position287, tokenIndex287, depth287
							if buffer[position] != rune('R') {
								goto l283
							}
							position++
						}
					l287:
						{
							position289, tokenIndex289, depth289 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l290
							}
							position++
							goto l289
						l290:
							position, tokenIndex, depth = position289, tokenIndex289, depth289
							if buffer[position] != rune('O') {
								goto l283
							}
							position++
						}
					l289:
						{
							position291, tokenIndex291, depth291 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l292
							}
							position++
							goto l291
						l292:
							position, tokenIndex, depth = position291, tokenIndex291, depth291
							if buffer[position] != rune('M') {
								goto l283
							}
							position++
						}
					l291:
						if !_rules[rulesp]() {
							goto l283
						}
						if !_rules[ruleDefRelations]() {
							goto l283
						}
						if !_rules[rulesp]() {
							goto l283
						}
						goto l284
					l283:
						position, tokenIndex, depth = position283, tokenIndex283, depth283
					}
				l284:
					depth--
					add(rulePegText, position282)
				}
				if !_rules[ruleAction11]() {
					goto l280
				}
				depth--
				add(ruleDefWindowedFrom, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 14 Interval <- <(NumericLiteral sp IntervalUnit Action12)> */
		func() bool {
			position293, tokenIndex293, depth293 := position, tokenIndex, depth
			{
				position294 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l293
				}
				if !_rules[rulesp]() {
					goto l293
				}
				if !_rules[ruleIntervalUnit]() {
					goto l293
				}
				if !_rules[ruleAction12]() {
					goto l293
				}
				depth--
				add(ruleInterval, position294)
			}
			return true
		l293:
			position, tokenIndex, depth = position293, tokenIndex293, depth293
			return false
		},
		/* 15 Relations <- <(RelationLike sp (',' sp RelationLike)*)> */
		func() bool {
			position295, tokenIndex295, depth295 := position, tokenIndex, depth
			{
				position296 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l295
				}
				if !_rules[rulesp]() {
					goto l295
				}
			l297:
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l298
					}
					position++
					if !_rules[rulesp]() {
						goto l298
					}
					if !_rules[ruleRelationLike]() {
						goto l298
					}
					goto l297
				l298:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
				}
				depth--
				add(ruleRelations, position296)
			}
			return true
		l295:
			position, tokenIndex, depth = position295, tokenIndex295, depth295
			return false
		},
		/* 16 DefRelations <- <(DefRelationLike sp (',' sp DefRelationLike)*)> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{
				position300 := position
				depth++
				if !_rules[ruleDefRelationLike]() {
					goto l299
				}
				if !_rules[rulesp]() {
					goto l299
				}
			l301:
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l302
					}
					position++
					if !_rules[rulesp]() {
						goto l302
					}
					if !_rules[ruleDefRelationLike]() {
						goto l302
					}
					goto l301
				l302:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
				}
				depth--
				add(ruleDefRelations, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 17 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action13)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				{
					position305 := position
					depth++
					{
						position306, tokenIndex306, depth306 := position, tokenIndex, depth
						{
							position308, tokenIndex308, depth308 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l309
							}
							position++
							goto l308
						l309:
							position, tokenIndex, depth = position308, tokenIndex308, depth308
							if buffer[position] != rune('W') {
								goto l306
							}
							position++
						}
					l308:
						{
							position310, tokenIndex310, depth310 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l311
							}
							position++
							goto l310
						l311:
							position, tokenIndex, depth = position310, tokenIndex310, depth310
							if buffer[position] != rune('H') {
								goto l306
							}
							position++
						}
					l310:
						{
							position312, tokenIndex312, depth312 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l313
							}
							position++
							goto l312
						l313:
							position, tokenIndex, depth = position312, tokenIndex312, depth312
							if buffer[position] != rune('E') {
								goto l306
							}
							position++
						}
					l312:
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('R') {
								goto l306
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('E') {
								goto l306
							}
							position++
						}
					l316:
						if !_rules[rulesp]() {
							goto l306
						}
						if !_rules[ruleExpression]() {
							goto l306
						}
						goto l307
					l306:
						position, tokenIndex, depth = position306, tokenIndex306, depth306
					}
				l307:
					depth--
					add(rulePegText, position305)
				}
				if !_rules[ruleAction13]() {
					goto l303
				}
				depth--
				add(ruleFilter, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 18 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action14)> */
		func() bool {
			position318, tokenIndex318, depth318 := position, tokenIndex, depth
			{
				position319 := position
				depth++
				{
					position320 := position
					depth++
					{
						position321, tokenIndex321, depth321 := position, tokenIndex, depth
						{
							position323, tokenIndex323, depth323 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l324
							}
							position++
							goto l323
						l324:
							position, tokenIndex, depth = position323, tokenIndex323, depth323
							if buffer[position] != rune('G') {
								goto l321
							}
							position++
						}
					l323:
						{
							position325, tokenIndex325, depth325 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l326
							}
							position++
							goto l325
						l326:
							position, tokenIndex, depth = position325, tokenIndex325, depth325
							if buffer[position] != rune('R') {
								goto l321
							}
							position++
						}
					l325:
						{
							position327, tokenIndex327, depth327 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l328
							}
							position++
							goto l327
						l328:
							position, tokenIndex, depth = position327, tokenIndex327, depth327
							if buffer[position] != rune('O') {
								goto l321
							}
							position++
						}
					l327:
						{
							position329, tokenIndex329, depth329 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l330
							}
							position++
							goto l329
						l330:
							position, tokenIndex, depth = position329, tokenIndex329, depth329
							if buffer[position] != rune('U') {
								goto l321
							}
							position++
						}
					l329:
						{
							position331, tokenIndex331, depth331 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l332
							}
							position++
							goto l331
						l332:
							position, tokenIndex, depth = position331, tokenIndex331, depth331
							if buffer[position] != rune('P') {
								goto l321
							}
							position++
						}
					l331:
						if !_rules[rulesp]() {
							goto l321
						}
						{
							position333, tokenIndex333, depth333 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l334
							}
							position++
							goto l333
						l334:
							position, tokenIndex, depth = position333, tokenIndex333, depth333
							if buffer[position] != rune('B') {
								goto l321
							}
							position++
						}
					l333:
						{
							position335, tokenIndex335, depth335 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l336
							}
							position++
							goto l335
						l336:
							position, tokenIndex, depth = position335, tokenIndex335, depth335
							if buffer[position] != rune('Y') {
								goto l321
							}
							position++
						}
					l335:
						if !_rules[rulesp]() {
							goto l321
						}
						if !_rules[ruleGroupList]() {
							goto l321
						}
						goto l322
					l321:
						position, tokenIndex, depth = position321, tokenIndex321, depth321
					}
				l322:
					depth--
					add(rulePegText, position320)
				}
				if !_rules[ruleAction14]() {
					goto l318
				}
				depth--
				add(ruleGrouping, position319)
			}
			return true
		l318:
			position, tokenIndex, depth = position318, tokenIndex318, depth318
			return false
		},
		/* 19 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l337
				}
				if !_rules[rulesp]() {
					goto l337
				}
			l339:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l340
					}
					position++
					if !_rules[rulesp]() {
						goto l340
					}
					if !_rules[ruleExpression]() {
						goto l340
					}
					goto l339
				l340:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
				}
				depth--
				add(ruleGroupList, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 20 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action15)> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
			{
				position342 := position
				depth++
				{
					position343 := position
					depth++
					{
						position344, tokenIndex344, depth344 := position, tokenIndex, depth
						{
							position346, tokenIndex346, depth346 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l347
							}
							position++
							goto l346
						l347:
							position, tokenIndex, depth = position346, tokenIndex346, depth346
							if buffer[position] != rune('H') {
								goto l344
							}
							position++
						}
					l346:
						{
							position348, tokenIndex348, depth348 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l349
							}
							position++
							goto l348
						l349:
							position, tokenIndex, depth = position348, tokenIndex348, depth348
							if buffer[position] != rune('A') {
								goto l344
							}
							position++
						}
					l348:
						{
							position350, tokenIndex350, depth350 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l351
							}
							position++
							goto l350
						l351:
							position, tokenIndex, depth = position350, tokenIndex350, depth350
							if buffer[position] != rune('V') {
								goto l344
							}
							position++
						}
					l350:
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
								goto l344
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
								goto l344
							}
							position++
						}
					l354:
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l357
							}
							position++
							goto l356
						l357:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
							if buffer[position] != rune('G') {
								goto l344
							}
							position++
						}
					l356:
						if !_rules[rulesp]() {
							goto l344
						}
						if !_rules[ruleExpression]() {
							goto l344
						}
						goto l345
					l344:
						position, tokenIndex, depth = position344, tokenIndex344, depth344
					}
				l345:
					depth--
					add(rulePegText, position343)
				}
				if !_rules[ruleAction15]() {
					goto l341
				}
				depth--
				add(ruleHaving, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
			return false
		},
		/* 21 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action16))> */
		func() bool {
			position358, tokenIndex358, depth358 := position, tokenIndex, depth
			{
				position359 := position
				depth++
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l361
					}
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if !_rules[ruleStreamWindow]() {
						goto l358
					}
					if !_rules[ruleAction16]() {
						goto l358
					}
				}
			l360:
				depth--
				add(ruleRelationLike, position359)
			}
			return true
		l358:
			position, tokenIndex, depth = position358, tokenIndex358, depth358
			return false
		},
		/* 22 DefRelationLike <- <(DefAliasedStreamWindow / (DefStreamWindow Action17))> */
		func() bool {
			position362, tokenIndex362, depth362 := position, tokenIndex, depth
			{
				position363 := position
				depth++
				{
					position364, tokenIndex364, depth364 := position, tokenIndex, depth
					if !_rules[ruleDefAliasedStreamWindow]() {
						goto l365
					}
					goto l364
				l365:
					position, tokenIndex, depth = position364, tokenIndex364, depth364
					if !_rules[ruleDefStreamWindow]() {
						goto l362
					}
					if !_rules[ruleAction17]() {
						goto l362
					}
				}
			l364:
				depth--
				add(ruleDefRelationLike, position363)
			}
			return true
		l362:
			position, tokenIndex, depth = position362, tokenIndex362, depth362
			return false
		},
		/* 23 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action18)> */
		func() bool {
			position366, tokenIndex366, depth366 := position, tokenIndex, depth
			{
				position367 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l366
				}
				if !_rules[rulesp]() {
					goto l366
				}
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if buffer[position] != rune('A') {
						goto l366
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
						goto l366
					}
					position++
				}
			l370:
				if !_rules[rulesp]() {
					goto l366
				}
				if !_rules[ruleIdentifier]() {
					goto l366
				}
				if !_rules[ruleAction18]() {
					goto l366
				}
				depth--
				add(ruleAliasedStreamWindow, position367)
			}
			return true
		l366:
			position, tokenIndex, depth = position366, tokenIndex366, depth366
			return false
		},
		/* 24 DefAliasedStreamWindow <- <(DefStreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action19)> */
		func() bool {
			position372, tokenIndex372, depth372 := position, tokenIndex, depth
			{
				position373 := position
				depth++
				if !_rules[ruleDefStreamWindow]() {
					goto l372
				}
				if !_rules[rulesp]() {
					goto l372
				}
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('A') {
						goto l372
					}
					position++
				}
			l374:
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l377
					}
					position++
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					if buffer[position] != rune('S') {
						goto l372
					}
					position++
				}
			l376:
				if !_rules[rulesp]() {
					goto l372
				}
				if !_rules[ruleIdentifier]() {
					goto l372
				}
				if !_rules[ruleAction19]() {
					goto l372
				}
				depth--
				add(ruleDefAliasedStreamWindow, position373)
			}
			return true
		l372:
			position, tokenIndex, depth = position372, tokenIndex372, depth372
			return false
		},
		/* 25 StreamWindow <- <(Stream sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']' Action20)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				if !_rules[ruleStream]() {
					goto l378
				}
				if !_rules[rulesp]() {
					goto l378
				}
				if buffer[position] != rune('[') {
					goto l378
				}
				position++
				if !_rules[rulesp]() {
					goto l378
				}
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('R') {
						goto l378
					}
					position++
				}
			l380:
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l383
					}
					position++
					goto l382
				l383:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
					if buffer[position] != rune('A') {
						goto l378
					}
					position++
				}
			l382:
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l385
					}
					position++
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if buffer[position] != rune('N') {
						goto l378
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('G') {
						goto l378
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
						goto l378
					}
					position++
				}
			l388:
				if !_rules[rulesp]() {
					goto l378
				}
				if !_rules[ruleInterval]() {
					goto l378
				}
				if !_rules[rulesp]() {
					goto l378
				}
				if buffer[position] != rune(']') {
					goto l378
				}
				position++
				if !_rules[ruleAction20]() {
					goto l378
				}
				depth--
				add(ruleStreamWindow, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 26 DefStreamWindow <- <(Stream (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval sp ']')? Action21)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				if !_rules[ruleStream]() {
					goto l390
				}
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l392
					}
					if buffer[position] != rune('[') {
						goto l392
					}
					position++
					if !_rules[rulesp]() {
						goto l392
					}
					{
						position394, tokenIndex394, depth394 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l395
						}
						position++
						goto l394
					l395:
						position, tokenIndex, depth = position394, tokenIndex394, depth394
						if buffer[position] != rune('R') {
							goto l392
						}
						position++
					}
				l394:
					{
						position396, tokenIndex396, depth396 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l397
						}
						position++
						goto l396
					l397:
						position, tokenIndex, depth = position396, tokenIndex396, depth396
						if buffer[position] != rune('A') {
							goto l392
						}
						position++
					}
				l396:
					{
						position398, tokenIndex398, depth398 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l399
						}
						position++
						goto l398
					l399:
						position, tokenIndex, depth = position398, tokenIndex398, depth398
						if buffer[position] != rune('N') {
							goto l392
						}
						position++
					}
				l398:
					{
						position400, tokenIndex400, depth400 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l401
						}
						position++
						goto l400
					l401:
						position, tokenIndex, depth = position400, tokenIndex400, depth400
						if buffer[position] != rune('G') {
							goto l392
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
							goto l392
						}
						position++
					}
				l402:
					if !_rules[rulesp]() {
						goto l392
					}
					if !_rules[ruleInterval]() {
						goto l392
					}
					if !_rules[rulesp]() {
						goto l392
					}
					if buffer[position] != rune(']') {
						goto l392
					}
					position++
					goto l393
				l392:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
				}
			l393:
				if !_rules[ruleAction21]() {
					goto l390
				}
				depth--
				add(ruleDefStreamWindow, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 27 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action22)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				{
					position406 := position
					depth++
					{
						position407, tokenIndex407, depth407 := position, tokenIndex, depth
						{
							position409, tokenIndex409, depth409 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l410
							}
							position++
							goto l409
						l410:
							position, tokenIndex, depth = position409, tokenIndex409, depth409
							if buffer[position] != rune('W') {
								goto l407
							}
							position++
						}
					l409:
						{
							position411, tokenIndex411, depth411 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l412
							}
							position++
							goto l411
						l412:
							position, tokenIndex, depth = position411, tokenIndex411, depth411
							if buffer[position] != rune('I') {
								goto l407
							}
							position++
						}
					l411:
						{
							position413, tokenIndex413, depth413 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l414
							}
							position++
							goto l413
						l414:
							position, tokenIndex, depth = position413, tokenIndex413, depth413
							if buffer[position] != rune('T') {
								goto l407
							}
							position++
						}
					l413:
						{
							position415, tokenIndex415, depth415 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l416
							}
							position++
							goto l415
						l416:
							position, tokenIndex, depth = position415, tokenIndex415, depth415
							if buffer[position] != rune('H') {
								goto l407
							}
							position++
						}
					l415:
						if !_rules[rulesp]() {
							goto l407
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l407
						}
						if !_rules[rulesp]() {
							goto l407
						}
					l417:
						{
							position418, tokenIndex418, depth418 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l418
							}
							position++
							if !_rules[rulesp]() {
								goto l418
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l418
							}
							goto l417
						l418:
							position, tokenIndex, depth = position418, tokenIndex418, depth418
						}
						goto l408
					l407:
						position, tokenIndex, depth = position407, tokenIndex407, depth407
					}
				l408:
					depth--
					add(rulePegText, position406)
				}
				if !_rules[ruleAction22]() {
					goto l404
				}
				depth--
				add(ruleSourceSinkSpecs, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 28 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action23)> */
		func() bool {
			position419, tokenIndex419, depth419 := position, tokenIndex, depth
			{
				position420 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l419
				}
				if buffer[position] != rune('=') {
					goto l419
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l419
				}
				if !_rules[ruleAction23]() {
					goto l419
				}
				depth--
				add(ruleSourceSinkParam, position420)
			}
			return true
		l419:
			position, tokenIndex, depth = position419, tokenIndex419, depth419
			return false
		},
		/* 29 Expression <- <orExpr> */
		func() bool {
			position421, tokenIndex421, depth421 := position, tokenIndex, depth
			{
				position422 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l421
				}
				depth--
				add(ruleExpression, position422)
			}
			return true
		l421:
			position, tokenIndex, depth = position421, tokenIndex421, depth421
			return false
		},
		/* 30 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action24)> */
		func() bool {
			position423, tokenIndex423, depth423 := position, tokenIndex, depth
			{
				position424 := position
				depth++
				{
					position425 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l423
					}
					if !_rules[rulesp]() {
						goto l423
					}
					{
						position426, tokenIndex426, depth426 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l426
						}
						if !_rules[rulesp]() {
							goto l426
						}
						if !_rules[ruleandExpr]() {
							goto l426
						}
						goto l427
					l426:
						position, tokenIndex, depth = position426, tokenIndex426, depth426
					}
				l427:
					depth--
					add(rulePegText, position425)
				}
				if !_rules[ruleAction24]() {
					goto l423
				}
				depth--
				add(ruleorExpr, position424)
			}
			return true
		l423:
			position, tokenIndex, depth = position423, tokenIndex423, depth423
			return false
		},
		/* 31 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action25)> */
		func() bool {
			position428, tokenIndex428, depth428 := position, tokenIndex, depth
			{
				position429 := position
				depth++
				{
					position430 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l428
					}
					if !_rules[rulesp]() {
						goto l428
					}
					{
						position431, tokenIndex431, depth431 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l431
						}
						if !_rules[rulesp]() {
							goto l431
						}
						if !_rules[rulecomparisonExpr]() {
							goto l431
						}
						goto l432
					l431:
						position, tokenIndex, depth = position431, tokenIndex431, depth431
					}
				l432:
					depth--
					add(rulePegText, position430)
				}
				if !_rules[ruleAction25]() {
					goto l428
				}
				depth--
				add(ruleandExpr, position429)
			}
			return true
		l428:
			position, tokenIndex, depth = position428, tokenIndex428, depth428
			return false
		},
		/* 32 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action26)> */
		func() bool {
			position433, tokenIndex433, depth433 := position, tokenIndex, depth
			{
				position434 := position
				depth++
				{
					position435 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l433
					}
					if !_rules[rulesp]() {
						goto l433
					}
					{
						position436, tokenIndex436, depth436 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l436
						}
						if !_rules[rulesp]() {
							goto l436
						}
						if !_rules[ruletermExpr]() {
							goto l436
						}
						goto l437
					l436:
						position, tokenIndex, depth = position436, tokenIndex436, depth436
					}
				l437:
					depth--
					add(rulePegText, position435)
				}
				if !_rules[ruleAction26]() {
					goto l433
				}
				depth--
				add(rulecomparisonExpr, position434)
			}
			return true
		l433:
			position, tokenIndex, depth = position433, tokenIndex433, depth433
			return false
		},
		/* 33 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action27)> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				{
					position440 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l438
					}
					if !_rules[rulesp]() {
						goto l438
					}
					{
						position441, tokenIndex441, depth441 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l441
						}
						if !_rules[rulesp]() {
							goto l441
						}
						if !_rules[ruleproductExpr]() {
							goto l441
						}
						goto l442
					l441:
						position, tokenIndex, depth = position441, tokenIndex441, depth441
					}
				l442:
					depth--
					add(rulePegText, position440)
				}
				if !_rules[ruleAction27]() {
					goto l438
				}
				depth--
				add(ruletermExpr, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 34 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action28)> */
		func() bool {
			position443, tokenIndex443, depth443 := position, tokenIndex, depth
			{
				position444 := position
				depth++
				{
					position445 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l443
					}
					if !_rules[rulesp]() {
						goto l443
					}
					{
						position446, tokenIndex446, depth446 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l446
						}
						if !_rules[rulesp]() {
							goto l446
						}
						if !_rules[rulebaseExpr]() {
							goto l446
						}
						goto l447
					l446:
						position, tokenIndex, depth = position446, tokenIndex446, depth446
					}
				l447:
					depth--
					add(rulePegText, position445)
				}
				if !_rules[ruleAction28]() {
					goto l443
				}
				depth--
				add(ruleproductExpr, position444)
			}
			return true
		l443:
			position, tokenIndex, depth = position443, tokenIndex443, depth443
			return false
		},
		/* 35 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position448, tokenIndex448, depth448 := position, tokenIndex, depth
			{
				position449 := position
				depth++
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l451
					}
					position++
					if !_rules[rulesp]() {
						goto l451
					}
					if !_rules[ruleExpression]() {
						goto l451
					}
					if !_rules[rulesp]() {
						goto l451
					}
					if buffer[position] != rune(')') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if !_rules[ruleBooleanLiteral]() {
						goto l452
					}
					goto l450
				l452:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if !_rules[ruleFuncApp]() {
						goto l453
					}
					goto l450
				l453:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if !_rules[ruleRowValue]() {
						goto l454
					}
					goto l450
				l454:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if !_rules[ruleLiteral]() {
						goto l448
					}
				}
			l450:
				depth--
				add(rulebaseExpr, position449)
			}
			return true
		l448:
			position, tokenIndex, depth = position448, tokenIndex448, depth448
			return false
		},
		/* 36 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action29)> */
		func() bool {
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l455
				}
				if !_rules[rulesp]() {
					goto l455
				}
				if buffer[position] != rune('(') {
					goto l455
				}
				position++
				if !_rules[rulesp]() {
					goto l455
				}
				if !_rules[ruleFuncParams]() {
					goto l455
				}
				if !_rules[rulesp]() {
					goto l455
				}
				if buffer[position] != rune(')') {
					goto l455
				}
				position++
				if !_rules[ruleAction29]() {
					goto l455
				}
				depth--
				add(ruleFuncApp, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
			return false
		},
		/* 37 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action30)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				{
					position459 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l457
					}
					if !_rules[rulesp]() {
						goto l457
					}
				l460:
					{
						position461, tokenIndex461, depth461 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l461
						}
						position++
						if !_rules[rulesp]() {
							goto l461
						}
						if !_rules[ruleExpression]() {
							goto l461
						}
						goto l460
					l461:
						position, tokenIndex, depth = position461, tokenIndex461, depth461
					}
					depth--
					add(rulePegText, position459)
				}
				if !_rules[ruleAction30]() {
					goto l457
				}
				depth--
				add(ruleFuncParams, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 38 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l465
					}
					goto l464
				l465:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					if !_rules[ruleNumericLiteral]() {
						goto l466
					}
					goto l464
				l466:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					if !_rules[ruleStringLiteral]() {
						goto l462
					}
				}
			l464:
				depth--
				add(ruleLiteral, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 39 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l470
					}
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if !_rules[ruleNotEqual]() {
						goto l471
					}
					goto l469
				l471:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if !_rules[ruleLessOrEqual]() {
						goto l472
					}
					goto l469
				l472:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if !_rules[ruleLess]() {
						goto l473
					}
					goto l469
				l473:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if !_rules[ruleGreaterOrEqual]() {
						goto l474
					}
					goto l469
				l474:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if !_rules[ruleGreater]() {
						goto l475
					}
					goto l469
				l475:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					if !_rules[ruleNotEqual]() {
						goto l467
					}
				}
			l469:
				depth--
				add(ruleComparisonOp, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 40 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position476, tokenIndex476, depth476 := position, tokenIndex, depth
			{
				position477 := position
				depth++
				{
					position478, tokenIndex478, depth478 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l479
					}
					goto l478
				l479:
					position, tokenIndex, depth = position478, tokenIndex478, depth478
					if !_rules[ruleMinus]() {
						goto l476
					}
				}
			l478:
				depth--
				add(rulePlusMinusOp, position477)
			}
			return true
		l476:
			position, tokenIndex, depth = position476, tokenIndex476, depth476
			return false
		},
		/* 41 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position480, tokenIndex480, depth480 := position, tokenIndex, depth
			{
				position481 := position
				depth++
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l483
					}
					goto l482
				l483:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
					if !_rules[ruleDivide]() {
						goto l484
					}
					goto l482
				l484:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
					if !_rules[ruleModulo]() {
						goto l480
					}
				}
			l482:
				depth--
				add(ruleMultDivOp, position481)
			}
			return true
		l480:
			position, tokenIndex, depth = position480, tokenIndex480, depth480
			return false
		},
		/* 42 Stream <- <(<ident> Action31)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				{
					position487 := position
					depth++
					if !_rules[ruleident]() {
						goto l485
					}
					depth--
					add(rulePegText, position487)
				}
				if !_rules[ruleAction31]() {
					goto l485
				}
				depth--
				add(ruleStream, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 43 RowValue <- <(<((ident '.')? ident)> Action32)> */
		func() bool {
			position488, tokenIndex488, depth488 := position, tokenIndex, depth
			{
				position489 := position
				depth++
				{
					position490 := position
					depth++
					{
						position491, tokenIndex491, depth491 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l491
						}
						if buffer[position] != rune('.') {
							goto l491
						}
						position++
						goto l492
					l491:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
					}
				l492:
					if !_rules[ruleident]() {
						goto l488
					}
					depth--
					add(rulePegText, position490)
				}
				if !_rules[ruleAction32]() {
					goto l488
				}
				depth--
				add(ruleRowValue, position489)
			}
			return true
		l488:
			position, tokenIndex, depth = position488, tokenIndex488, depth488
			return false
		},
		/* 44 NumericLiteral <- <(<('-'? [0-9]+)> Action33)> */
		func() bool {
			position493, tokenIndex493, depth493 := position, tokenIndex, depth
			{
				position494 := position
				depth++
				{
					position495 := position
					depth++
					{
						position496, tokenIndex496, depth496 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l496
						}
						position++
						goto l497
					l496:
						position, tokenIndex, depth = position496, tokenIndex496, depth496
					}
				l497:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l493
					}
					position++
				l498:
					{
						position499, tokenIndex499, depth499 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l499
						}
						position++
						goto l498
					l499:
						position, tokenIndex, depth = position499, tokenIndex499, depth499
					}
					depth--
					add(rulePegText, position495)
				}
				if !_rules[ruleAction33]() {
					goto l493
				}
				depth--
				add(ruleNumericLiteral, position494)
			}
			return true
		l493:
			position, tokenIndex, depth = position493, tokenIndex493, depth493
			return false
		},
		/* 45 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action34)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				{
					position502 := position
					depth++
					{
						position503, tokenIndex503, depth503 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l503
						}
						position++
						goto l504
					l503:
						position, tokenIndex, depth = position503, tokenIndex503, depth503
					}
				l504:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l500
					}
					position++
				l505:
					{
						position506, tokenIndex506, depth506 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l506
						}
						position++
						goto l505
					l506:
						position, tokenIndex, depth = position506, tokenIndex506, depth506
					}
					if buffer[position] != rune('.') {
						goto l500
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l500
					}
					position++
				l507:
					{
						position508, tokenIndex508, depth508 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l508
						}
						position++
						goto l507
					l508:
						position, tokenIndex, depth = position508, tokenIndex508, depth508
					}
					depth--
					add(rulePegText, position502)
				}
				if !_rules[ruleAction34]() {
					goto l500
				}
				depth--
				add(ruleFloatLiteral, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 46 Function <- <(<ident> Action35)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				{
					position511 := position
					depth++
					if !_rules[ruleident]() {
						goto l509
					}
					depth--
					add(rulePegText, position511)
				}
				if !_rules[ruleAction35]() {
					goto l509
				}
				depth--
				add(ruleFunction, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 47 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position512, tokenIndex512, depth512 := position, tokenIndex, depth
			{
				position513 := position
				depth++
				{
					position514, tokenIndex514, depth514 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l515
					}
					goto l514
				l515:
					position, tokenIndex, depth = position514, tokenIndex514, depth514
					if !_rules[ruleFALSE]() {
						goto l512
					}
				}
			l514:
				depth--
				add(ruleBooleanLiteral, position513)
			}
			return true
		l512:
			position, tokenIndex, depth = position512, tokenIndex512, depth512
			return false
		},
		/* 48 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action36)> */
		func() bool {
			position516, tokenIndex516, depth516 := position, tokenIndex, depth
			{
				position517 := position
				depth++
				{
					position518 := position
					depth++
					{
						position519, tokenIndex519, depth519 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l520
						}
						position++
						goto l519
					l520:
						position, tokenIndex, depth = position519, tokenIndex519, depth519
						if buffer[position] != rune('T') {
							goto l516
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
							goto l516
						}
						position++
					}
				l521:
					{
						position523, tokenIndex523, depth523 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l524
						}
						position++
						goto l523
					l524:
						position, tokenIndex, depth = position523, tokenIndex523, depth523
						if buffer[position] != rune('U') {
							goto l516
						}
						position++
					}
				l523:
					{
						position525, tokenIndex525, depth525 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l526
						}
						position++
						goto l525
					l526:
						position, tokenIndex, depth = position525, tokenIndex525, depth525
						if buffer[position] != rune('E') {
							goto l516
						}
						position++
					}
				l525:
					depth--
					add(rulePegText, position518)
				}
				if !_rules[ruleAction36]() {
					goto l516
				}
				depth--
				add(ruleTRUE, position517)
			}
			return true
		l516:
			position, tokenIndex, depth = position516, tokenIndex516, depth516
			return false
		},
		/* 49 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action37)> */
		func() bool {
			position527, tokenIndex527, depth527 := position, tokenIndex, depth
			{
				position528 := position
				depth++
				{
					position529 := position
					depth++
					{
						position530, tokenIndex530, depth530 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l531
						}
						position++
						goto l530
					l531:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
						if buffer[position] != rune('F') {
							goto l527
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
							goto l527
						}
						position++
					}
				l532:
					{
						position534, tokenIndex534, depth534 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l535
						}
						position++
						goto l534
					l535:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if buffer[position] != rune('L') {
							goto l527
						}
						position++
					}
				l534:
					{
						position536, tokenIndex536, depth536 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l537
						}
						position++
						goto l536
					l537:
						position, tokenIndex, depth = position536, tokenIndex536, depth536
						if buffer[position] != rune('S') {
							goto l527
						}
						position++
					}
				l536:
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l539
						}
						position++
						goto l538
					l539:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
						if buffer[position] != rune('E') {
							goto l527
						}
						position++
					}
				l538:
					depth--
					add(rulePegText, position529)
				}
				if !_rules[ruleAction37]() {
					goto l527
				}
				depth--
				add(ruleFALSE, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 50 Wildcard <- <(<'*'> Action38)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					if buffer[position] != rune('*') {
						goto l540
					}
					position++
					depth--
					add(rulePegText, position542)
				}
				if !_rules[ruleAction38]() {
					goto l540
				}
				depth--
				add(ruleWildcard, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 51 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action39)> */
		func() bool {
			position543, tokenIndex543, depth543 := position, tokenIndex, depth
			{
				position544 := position
				depth++
				{
					position545 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l543
					}
					position++
				l546:
					{
						position547, tokenIndex547, depth547 := position, tokenIndex, depth
						{
							position548, tokenIndex548, depth548 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l549
							}
							position++
							if buffer[position] != rune('\'') {
								goto l549
							}
							position++
							goto l548
						l549:
							position, tokenIndex, depth = position548, tokenIndex548, depth548
							{
								position550, tokenIndex550, depth550 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l550
								}
								position++
								goto l547
							l550:
								position, tokenIndex, depth = position550, tokenIndex550, depth550
							}
							if !matchDot() {
								goto l547
							}
						}
					l548:
						goto l546
					l547:
						position, tokenIndex, depth = position547, tokenIndex547, depth547
					}
					if buffer[position] != rune('\'') {
						goto l543
					}
					position++
					depth--
					add(rulePegText, position545)
				}
				if !_rules[ruleAction39]() {
					goto l543
				}
				depth--
				add(ruleStringLiteral, position544)
			}
			return true
		l543:
			position, tokenIndex, depth = position543, tokenIndex543, depth543
			return false
		},
		/* 52 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action40)> */
		func() bool {
			position551, tokenIndex551, depth551 := position, tokenIndex, depth
			{
				position552 := position
				depth++
				{
					position553 := position
					depth++
					{
						position554, tokenIndex554, depth554 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l555
						}
						position++
						goto l554
					l555:
						position, tokenIndex, depth = position554, tokenIndex554, depth554
						if buffer[position] != rune('I') {
							goto l551
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
							goto l551
						}
						position++
					}
				l556:
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
							goto l551
						}
						position++
					}
				l558:
					{
						position560, tokenIndex560, depth560 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l561
						}
						position++
						goto l560
					l561:
						position, tokenIndex, depth = position560, tokenIndex560, depth560
						if buffer[position] != rune('R') {
							goto l551
						}
						position++
					}
				l560:
					{
						position562, tokenIndex562, depth562 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l563
						}
						position++
						goto l562
					l563:
						position, tokenIndex, depth = position562, tokenIndex562, depth562
						if buffer[position] != rune('E') {
							goto l551
						}
						position++
					}
				l562:
					{
						position564, tokenIndex564, depth564 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l565
						}
						position++
						goto l564
					l565:
						position, tokenIndex, depth = position564, tokenIndex564, depth564
						if buffer[position] != rune('A') {
							goto l551
						}
						position++
					}
				l564:
					{
						position566, tokenIndex566, depth566 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l567
						}
						position++
						goto l566
					l567:
						position, tokenIndex, depth = position566, tokenIndex566, depth566
						if buffer[position] != rune('M') {
							goto l551
						}
						position++
					}
				l566:
					depth--
					add(rulePegText, position553)
				}
				if !_rules[ruleAction40]() {
					goto l551
				}
				depth--
				add(ruleISTREAM, position552)
			}
			return true
		l551:
			position, tokenIndex, depth = position551, tokenIndex551, depth551
			return false
		},
		/* 53 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action41)> */
		func() bool {
			position568, tokenIndex568, depth568 := position, tokenIndex, depth
			{
				position569 := position
				depth++
				{
					position570 := position
					depth++
					{
						position571, tokenIndex571, depth571 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l572
						}
						position++
						goto l571
					l572:
						position, tokenIndex, depth = position571, tokenIndex571, depth571
						if buffer[position] != rune('D') {
							goto l568
						}
						position++
					}
				l571:
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l574
						}
						position++
						goto l573
					l574:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
						if buffer[position] != rune('S') {
							goto l568
						}
						position++
					}
				l573:
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l576
						}
						position++
						goto l575
					l576:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
						if buffer[position] != rune('T') {
							goto l568
						}
						position++
					}
				l575:
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('R') {
							goto l568
						}
						position++
					}
				l577:
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l580
						}
						position++
						goto l579
					l580:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
						if buffer[position] != rune('E') {
							goto l568
						}
						position++
					}
				l579:
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l582
						}
						position++
						goto l581
					l582:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
						if buffer[position] != rune('A') {
							goto l568
						}
						position++
					}
				l581:
					{
						position583, tokenIndex583, depth583 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l584
						}
						position++
						goto l583
					l584:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
						if buffer[position] != rune('M') {
							goto l568
						}
						position++
					}
				l583:
					depth--
					add(rulePegText, position570)
				}
				if !_rules[ruleAction41]() {
					goto l568
				}
				depth--
				add(ruleDSTREAM, position569)
			}
			return true
		l568:
			position, tokenIndex, depth = position568, tokenIndex568, depth568
			return false
		},
		/* 54 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action42)> */
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
						if buffer[position] != rune('r') {
							goto l589
						}
						position++
						goto l588
					l589:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
						if buffer[position] != rune('R') {
							goto l585
						}
						position++
					}
				l588:
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
						if buffer[position] != rune('S') {
							goto l585
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
							goto l585
						}
						position++
					}
				l592:
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
						if buffer[position] != rune('R') {
							goto l585
						}
						position++
					}
				l594:
					{
						position596, tokenIndex596, depth596 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l597
						}
						position++
						goto l596
					l597:
						position, tokenIndex, depth = position596, tokenIndex596, depth596
						if buffer[position] != rune('E') {
							goto l585
						}
						position++
					}
				l596:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l599
						}
						position++
						goto l598
					l599:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
						if buffer[position] != rune('A') {
							goto l585
						}
						position++
					}
				l598:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l601
						}
						position++
						goto l600
					l601:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
						if buffer[position] != rune('M') {
							goto l585
						}
						position++
					}
				l600:
					depth--
					add(rulePegText, position587)
				}
				if !_rules[ruleAction42]() {
					goto l585
				}
				depth--
				add(ruleRSTREAM, position586)
			}
			return true
		l585:
			position, tokenIndex, depth = position585, tokenIndex585, depth585
			return false
		},
		/* 55 IntervalUnit <- <(TUPLES / SECONDS)> */
		func() bool {
			position602, tokenIndex602, depth602 := position, tokenIndex, depth
			{
				position603 := position
				depth++
				{
					position604, tokenIndex604, depth604 := position, tokenIndex, depth
					if !_rules[ruleTUPLES]() {
						goto l605
					}
					goto l604
				l605:
					position, tokenIndex, depth = position604, tokenIndex604, depth604
					if !_rules[ruleSECONDS]() {
						goto l602
					}
				}
			l604:
				depth--
				add(ruleIntervalUnit, position603)
			}
			return true
		l602:
			position, tokenIndex, depth = position602, tokenIndex602, depth602
			return false
		},
		/* 56 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action43)> */
		func() bool {
			position606, tokenIndex606, depth606 := position, tokenIndex, depth
			{
				position607 := position
				depth++
				{
					position608 := position
					depth++
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('T') {
							goto l606
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
							goto l606
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
							goto l606
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('L') {
							goto l606
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('E') {
							goto l606
						}
						position++
					}
				l617:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('S') {
							goto l606
						}
						position++
					}
				l619:
					depth--
					add(rulePegText, position608)
				}
				if !_rules[ruleAction43]() {
					goto l606
				}
				depth--
				add(ruleTUPLES, position607)
			}
			return true
		l606:
			position, tokenIndex, depth = position606, tokenIndex606, depth606
			return false
		},
		/* 57 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action44)> */
		func() bool {
			position621, tokenIndex621, depth621 := position, tokenIndex, depth
			{
				position622 := position
				depth++
				{
					position623 := position
					depth++
					{
						position624, tokenIndex624, depth624 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('S') {
							goto l621
						}
						position++
					}
				l624:
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('E') {
							goto l621
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('C') {
							goto l621
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('O') {
							goto l621
						}
						position++
					}
				l630:
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l633
						}
						position++
						goto l632
					l633:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
						if buffer[position] != rune('N') {
							goto l621
						}
						position++
					}
				l632:
					{
						position634, tokenIndex634, depth634 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l635
						}
						position++
						goto l634
					l635:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
						if buffer[position] != rune('D') {
							goto l621
						}
						position++
					}
				l634:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l637
						}
						position++
						goto l636
					l637:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
						if buffer[position] != rune('S') {
							goto l621
						}
						position++
					}
				l636:
					depth--
					add(rulePegText, position623)
				}
				if !_rules[ruleAction44]() {
					goto l621
				}
				depth--
				add(ruleSECONDS, position622)
			}
			return true
		l621:
			position, tokenIndex, depth = position621, tokenIndex621, depth621
			return false
		},
		/* 58 StreamIdentifier <- <(<ident> Action45)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					if !_rules[ruleident]() {
						goto l638
					}
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction45]() {
					goto l638
				}
				depth--
				add(ruleStreamIdentifier, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 59 SourceSinkType <- <(<ident> Action46)> */
		func() bool {
			position641, tokenIndex641, depth641 := position, tokenIndex, depth
			{
				position642 := position
				depth++
				{
					position643 := position
					depth++
					if !_rules[ruleident]() {
						goto l641
					}
					depth--
					add(rulePegText, position643)
				}
				if !_rules[ruleAction46]() {
					goto l641
				}
				depth--
				add(ruleSourceSinkType, position642)
			}
			return true
		l641:
			position, tokenIndex, depth = position641, tokenIndex641, depth641
			return false
		},
		/* 60 SourceSinkParamKey <- <(<ident> Action47)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				{
					position646 := position
					depth++
					if !_rules[ruleident]() {
						goto l644
					}
					depth--
					add(rulePegText, position646)
				}
				if !_rules[ruleAction47]() {
					goto l644
				}
				depth--
				add(ruleSourceSinkParamKey, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 61 SourceSinkParamVal <- <(<([a-z] / [A-Z] / [0-9] / '_')+> Action48)> */
		func() bool {
			position647, tokenIndex647, depth647 := position, tokenIndex, depth
			{
				position648 := position
				depth++
				{
					position649 := position
					depth++
					{
						position652, tokenIndex652, depth652 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l653
						}
						position++
						goto l652
					l653:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l654
						}
						position++
						goto l652
					l654:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l655
						}
						position++
						goto l652
					l655:
						position, tokenIndex, depth = position652, tokenIndex652, depth652
						if buffer[position] != rune('_') {
							goto l647
						}
						position++
					}
				l652:
				l650:
					{
						position651, tokenIndex651, depth651 := position, tokenIndex, depth
						{
							position656, tokenIndex656, depth656 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l657
							}
							position++
							goto l656
						l657:
							position, tokenIndex, depth = position656, tokenIndex656, depth656
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l658
							}
							position++
							goto l656
						l658:
							position, tokenIndex, depth = position656, tokenIndex656, depth656
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l659
							}
							position++
							goto l656
						l659:
							position, tokenIndex, depth = position656, tokenIndex656, depth656
							if buffer[position] != rune('_') {
								goto l651
							}
							position++
						}
					l656:
						goto l650
					l651:
						position, tokenIndex, depth = position651, tokenIndex651, depth651
					}
					depth--
					add(rulePegText, position649)
				}
				if !_rules[ruleAction48]() {
					goto l647
				}
				depth--
				add(ruleSourceSinkParamVal, position648)
			}
			return true
		l647:
			position, tokenIndex, depth = position647, tokenIndex647, depth647
			return false
		},
		/* 62 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action49)> */
		func() bool {
			position660, tokenIndex660, depth660 := position, tokenIndex, depth
			{
				position661 := position
				depth++
				{
					position662 := position
					depth++
					{
						position663, tokenIndex663, depth663 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l664
						}
						position++
						goto l663
					l664:
						position, tokenIndex, depth = position663, tokenIndex663, depth663
						if buffer[position] != rune('O') {
							goto l660
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
							goto l660
						}
						position++
					}
				l665:
					depth--
					add(rulePegText, position662)
				}
				if !_rules[ruleAction49]() {
					goto l660
				}
				depth--
				add(ruleOr, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 63 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action50)> */
		func() bool {
			position667, tokenIndex667, depth667 := position, tokenIndex, depth
			{
				position668 := position
				depth++
				{
					position669 := position
					depth++
					{
						position670, tokenIndex670, depth670 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l671
						}
						position++
						goto l670
					l671:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
						if buffer[position] != rune('A') {
							goto l667
						}
						position++
					}
				l670:
					{
						position672, tokenIndex672, depth672 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l673
						}
						position++
						goto l672
					l673:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						if buffer[position] != rune('N') {
							goto l667
						}
						position++
					}
				l672:
					{
						position674, tokenIndex674, depth674 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l675
						}
						position++
						goto l674
					l675:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
						if buffer[position] != rune('D') {
							goto l667
						}
						position++
					}
				l674:
					depth--
					add(rulePegText, position669)
				}
				if !_rules[ruleAction50]() {
					goto l667
				}
				depth--
				add(ruleAnd, position668)
			}
			return true
		l667:
			position, tokenIndex, depth = position667, tokenIndex667, depth667
			return false
		},
		/* 64 Equal <- <(<'='> Action51)> */
		func() bool {
			position676, tokenIndex676, depth676 := position, tokenIndex, depth
			{
				position677 := position
				depth++
				{
					position678 := position
					depth++
					if buffer[position] != rune('=') {
						goto l676
					}
					position++
					depth--
					add(rulePegText, position678)
				}
				if !_rules[ruleAction51]() {
					goto l676
				}
				depth--
				add(ruleEqual, position677)
			}
			return true
		l676:
			position, tokenIndex, depth = position676, tokenIndex676, depth676
			return false
		},
		/* 65 Less <- <(<'<'> Action52)> */
		func() bool {
			position679, tokenIndex679, depth679 := position, tokenIndex, depth
			{
				position680 := position
				depth++
				{
					position681 := position
					depth++
					if buffer[position] != rune('<') {
						goto l679
					}
					position++
					depth--
					add(rulePegText, position681)
				}
				if !_rules[ruleAction52]() {
					goto l679
				}
				depth--
				add(ruleLess, position680)
			}
			return true
		l679:
			position, tokenIndex, depth = position679, tokenIndex679, depth679
			return false
		},
		/* 66 LessOrEqual <- <(<('<' '=')> Action53)> */
		func() bool {
			position682, tokenIndex682, depth682 := position, tokenIndex, depth
			{
				position683 := position
				depth++
				{
					position684 := position
					depth++
					if buffer[position] != rune('<') {
						goto l682
					}
					position++
					if buffer[position] != rune('=') {
						goto l682
					}
					position++
					depth--
					add(rulePegText, position684)
				}
				if !_rules[ruleAction53]() {
					goto l682
				}
				depth--
				add(ruleLessOrEqual, position683)
			}
			return true
		l682:
			position, tokenIndex, depth = position682, tokenIndex682, depth682
			return false
		},
		/* 67 Greater <- <(<'>'> Action54)> */
		func() bool {
			position685, tokenIndex685, depth685 := position, tokenIndex, depth
			{
				position686 := position
				depth++
				{
					position687 := position
					depth++
					if buffer[position] != rune('>') {
						goto l685
					}
					position++
					depth--
					add(rulePegText, position687)
				}
				if !_rules[ruleAction54]() {
					goto l685
				}
				depth--
				add(ruleGreater, position686)
			}
			return true
		l685:
			position, tokenIndex, depth = position685, tokenIndex685, depth685
			return false
		},
		/* 68 GreaterOrEqual <- <(<('>' '=')> Action55)> */
		func() bool {
			position688, tokenIndex688, depth688 := position, tokenIndex, depth
			{
				position689 := position
				depth++
				{
					position690 := position
					depth++
					if buffer[position] != rune('>') {
						goto l688
					}
					position++
					if buffer[position] != rune('=') {
						goto l688
					}
					position++
					depth--
					add(rulePegText, position690)
				}
				if !_rules[ruleAction55]() {
					goto l688
				}
				depth--
				add(ruleGreaterOrEqual, position689)
			}
			return true
		l688:
			position, tokenIndex, depth = position688, tokenIndex688, depth688
			return false
		},
		/* 69 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action56)> */
		func() bool {
			position691, tokenIndex691, depth691 := position, tokenIndex, depth
			{
				position692 := position
				depth++
				{
					position693 := position
					depth++
					{
						position694, tokenIndex694, depth694 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l695
						}
						position++
						if buffer[position] != rune('=') {
							goto l695
						}
						position++
						goto l694
					l695:
						position, tokenIndex, depth = position694, tokenIndex694, depth694
						if buffer[position] != rune('<') {
							goto l691
						}
						position++
						if buffer[position] != rune('>') {
							goto l691
						}
						position++
					}
				l694:
					depth--
					add(rulePegText, position693)
				}
				if !_rules[ruleAction56]() {
					goto l691
				}
				depth--
				add(ruleNotEqual, position692)
			}
			return true
		l691:
			position, tokenIndex, depth = position691, tokenIndex691, depth691
			return false
		},
		/* 70 Plus <- <(<'+'> Action57)> */
		func() bool {
			position696, tokenIndex696, depth696 := position, tokenIndex, depth
			{
				position697 := position
				depth++
				{
					position698 := position
					depth++
					if buffer[position] != rune('+') {
						goto l696
					}
					position++
					depth--
					add(rulePegText, position698)
				}
				if !_rules[ruleAction57]() {
					goto l696
				}
				depth--
				add(rulePlus, position697)
			}
			return true
		l696:
			position, tokenIndex, depth = position696, tokenIndex696, depth696
			return false
		},
		/* 71 Minus <- <(<'-'> Action58)> */
		func() bool {
			position699, tokenIndex699, depth699 := position, tokenIndex, depth
			{
				position700 := position
				depth++
				{
					position701 := position
					depth++
					if buffer[position] != rune('-') {
						goto l699
					}
					position++
					depth--
					add(rulePegText, position701)
				}
				if !_rules[ruleAction58]() {
					goto l699
				}
				depth--
				add(ruleMinus, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 72 Multiply <- <(<'*'> Action59)> */
		func() bool {
			position702, tokenIndex702, depth702 := position, tokenIndex, depth
			{
				position703 := position
				depth++
				{
					position704 := position
					depth++
					if buffer[position] != rune('*') {
						goto l702
					}
					position++
					depth--
					add(rulePegText, position704)
				}
				if !_rules[ruleAction59]() {
					goto l702
				}
				depth--
				add(ruleMultiply, position703)
			}
			return true
		l702:
			position, tokenIndex, depth = position702, tokenIndex702, depth702
			return false
		},
		/* 73 Divide <- <(<'/'> Action60)> */
		func() bool {
			position705, tokenIndex705, depth705 := position, tokenIndex, depth
			{
				position706 := position
				depth++
				{
					position707 := position
					depth++
					if buffer[position] != rune('/') {
						goto l705
					}
					position++
					depth--
					add(rulePegText, position707)
				}
				if !_rules[ruleAction60]() {
					goto l705
				}
				depth--
				add(ruleDivide, position706)
			}
			return true
		l705:
			position, tokenIndex, depth = position705, tokenIndex705, depth705
			return false
		},
		/* 74 Modulo <- <(<'%'> Action61)> */
		func() bool {
			position708, tokenIndex708, depth708 := position, tokenIndex, depth
			{
				position709 := position
				depth++
				{
					position710 := position
					depth++
					if buffer[position] != rune('%') {
						goto l708
					}
					position++
					depth--
					add(rulePegText, position710)
				}
				if !_rules[ruleAction61]() {
					goto l708
				}
				depth--
				add(ruleModulo, position709)
			}
			return true
		l708:
			position, tokenIndex, depth = position708, tokenIndex708, depth708
			return false
		},
		/* 75 Identifier <- <(<ident> Action62)> */
		func() bool {
			position711, tokenIndex711, depth711 := position, tokenIndex, depth
			{
				position712 := position
				depth++
				{
					position713 := position
					depth++
					if !_rules[ruleident]() {
						goto l711
					}
					depth--
					add(rulePegText, position713)
				}
				if !_rules[ruleAction62]() {
					goto l711
				}
				depth--
				add(ruleIdentifier, position712)
			}
			return true
		l711:
			position, tokenIndex, depth = position711, tokenIndex711, depth711
			return false
		},
		/* 76 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position714, tokenIndex714, depth714 := position, tokenIndex, depth
			{
				position715 := position
				depth++
				{
					position716, tokenIndex716, depth716 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l717
					}
					position++
					goto l716
				l717:
					position, tokenIndex, depth = position716, tokenIndex716, depth716
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l714
					}
					position++
				}
			l716:
			l718:
				{
					position719, tokenIndex719, depth719 := position, tokenIndex, depth
					{
						position720, tokenIndex720, depth720 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l721
						}
						position++
						goto l720
					l721:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l722
						}
						position++
						goto l720
					l722:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l723
						}
						position++
						goto l720
					l723:
						position, tokenIndex, depth = position720, tokenIndex720, depth720
						if buffer[position] != rune('_') {
							goto l719
						}
						position++
					}
				l720:
					goto l718
				l719:
					position, tokenIndex, depth = position719, tokenIndex719, depth719
				}
				depth--
				add(ruleident, position715)
			}
			return true
		l714:
			position, tokenIndex, depth = position714, tokenIndex714, depth714
			return false
		},
		/* 77 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position725 := position
				depth++
			l726:
				{
					position727, tokenIndex727, depth727 := position, tokenIndex, depth
					{
						position728, tokenIndex728, depth728 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l729
						}
						position++
						goto l728
					l729:
						position, tokenIndex, depth = position728, tokenIndex728, depth728
						if buffer[position] != rune('\t') {
							goto l730
						}
						position++
						goto l728
					l730:
						position, tokenIndex, depth = position728, tokenIndex728, depth728
						if buffer[position] != rune('\n') {
							goto l727
						}
						position++
					}
				l728:
					goto l726
				l727:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
				}
				depth--
				add(rulesp, position725)
			}
			return true
		},
		/* 79 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 80 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 81 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 82 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 83 Action4 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 84 Action5 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 85 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 86 Action7 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 88 Action8 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 89 Action9 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 90 Action10 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 91 Action11 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 92 Action12 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 93 Action13 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 94 Action14 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 95 Action15 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 96 Action16 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 97 Action17 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 98 Action18 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 99 Action19 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 100 Action20 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 101 Action21 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 102 Action22 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 103 Action23 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 104 Action24 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 105 Action25 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 106 Action26 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 107 Action27 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 108 Action28 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 109 Action29 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 110 Action30 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 111 Action31 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 112 Action32 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 113 Action33 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 114 Action34 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 115 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 116 Action36 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 117 Action37 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 118 Action38 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 119 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 120 Action40 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 121 Action41 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 122 Action42 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 123 Action43 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 124 Action44 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 125 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 126 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 127 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 128 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamVal(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 129 Action49 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 130 Action50 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 131 Action51 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 132 Action52 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 133 Action53 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 134 Action54 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 135 Action55 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 136 Action56 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 137 Action57 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 138 Action58 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 139 Action59 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 140 Action60 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 141 Action61 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 142 Action62 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
	}
	p.rules = _rules
}

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
	ruleEmitProjections
	ruleProjections
	ruleProjection
	ruleAliasExpression
	ruleWindowedFrom
	ruleRange
	ruleWindowedRelations
	ruleFilter
	ruleGrouping
	ruleGroupList
	ruleHaving
	ruleWindowedRelationLike
	ruleAliasWindowedRelation
	ruleWindowedRelation
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
	ruleRelation
	ruleRowValue
	ruleNumericLiteral
	ruleFloatLiteral
	ruleFunction
	ruleBooleanLiteral
	ruleTRUE
	ruleFALSE
	ruleWildcard
	ruleStringLiteral
	ruleEmitter
	ruleISTREAM
	ruleDSTREAM
	ruleRSTREAM
	ruleRangeUnit
	ruleTUPLES
	ruleSECONDS
	ruleSourceSinkName
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
	"EmitProjections",
	"Projections",
	"Projection",
	"AliasExpression",
	"WindowedFrom",
	"Range",
	"WindowedRelations",
	"Filter",
	"Grouping",
	"GroupList",
	"Having",
	"WindowedRelationLike",
	"AliasWindowedRelation",
	"WindowedRelation",
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
	"Relation",
	"RowValue",
	"NumericLiteral",
	"FloatLiteral",
	"Function",
	"BooleanLiteral",
	"TRUE",
	"FALSE",
	"Wildcard",
	"StringLiteral",
	"Emitter",
	"ISTREAM",
	"DSTREAM",
	"RSTREAM",
	"RangeUnit",
	"TUPLES",
	"SECONDS",
	"SourceSinkName",
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
	rules  [135]func() bool
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

			p.AssembleEmitProjections()

		case ruleAction8:

			p.AssembleProjections(begin, end)

		case ruleAction9:

			p.AssembleAlias()

		case ruleAction10:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction11:

			p.AssembleRange()

		case ruleAction12:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction13:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction14:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction15:

			p.EnsureAliasWindowedRelation()

		case ruleAction16:

			p.AssembleAliasWindowedRelation()

		case ruleAction17:

			p.AssembleWindowedRelation()

		case ruleAction18:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction19:

			p.AssembleSourceSinkParam()

		case ruleAction20:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction21:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction22:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction23:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction24:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction25:

			p.AssembleFuncApp()

		case ruleAction26:

			p.AssembleExpressions(begin, end)

		case ruleAction27:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRelation(substr))

		case ruleAction28:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction29:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction30:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction31:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction32:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction33:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction34:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction35:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction36:

			p.PushComponent(begin, end, Istream)

		case ruleAction37:

			p.PushComponent(begin, end, Dstream)

		case ruleAction38:

			p.PushComponent(begin, end, Rstream)

		case ruleAction39:

			p.PushComponent(begin, end, Tuples)

		case ruleAction40:

			p.PushComponent(begin, end, Seconds)

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkName(substr))

		case ruleAction42:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction43:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction44:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamVal(substr))

		case ruleAction45:

			p.PushComponent(begin, end, Or)

		case ruleAction46:

			p.PushComponent(begin, end, And)

		case ruleAction47:

			p.PushComponent(begin, end, Equal)

		case ruleAction48:

			p.PushComponent(begin, end, Less)

		case ruleAction49:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction50:

			p.PushComponent(begin, end, Greater)

		case ruleAction51:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction52:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction53:

			p.PushComponent(begin, end, Plus)

		case ruleAction54:

			p.PushComponent(begin, end, Minus)

		case ruleAction55:

			p.PushComponent(begin, end, Multiply)

		case ruleAction56:

			p.PushComponent(begin, end, Divide)

		case ruleAction57:

			p.PushComponent(begin, end, Modulo)

		case ruleAction58:

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
		/* 1 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Projections sp WindowedFrom sp Filter sp Grouping sp Having sp Action0)> */
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
				if !_rules[ruleWindowedFrom]() {
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
		/* 2 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp Relation sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp EmitProjections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
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
				if !_rules[ruleRelation]() {
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
				if !_rules[ruleEmitProjections]() {
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
		/* 3 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp SourceSinkName sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action2)> */
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
				if !_rules[ruleSourceSinkName]() {
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
		/* 4 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp SourceSinkName sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType sp SourceSinkSpecs Action3)> */
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
				if !_rules[ruleSourceSinkName]() {
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
		/* 5 CreateStreamFromSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp Relation sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp SourceSinkName Action4)> */
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
				if !_rules[ruleRelation]() {
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
				if !_rules[ruleSourceSinkName]() {
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
		/* 6 CreateStreamFromSourceExtStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp Relation sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp SourceSinkType sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp SourceSinkSpecs Action5)> */
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
				if !_rules[ruleRelation]() {
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
		/* 7 InsertIntoSelectStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp SourceSinkName sp SelectStmt Action6)> */
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
				if !_rules[ruleSourceSinkName]() {
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
		/* 8 EmitProjections <- <(Emitter sp Projections Action7)> */
		func() bool {
			position244, tokenIndex244, depth244 := position, tokenIndex, depth
			{
				position245 := position
				depth++
				if !_rules[ruleEmitter]() {
					goto l244
				}
				if !_rules[rulesp]() {
					goto l244
				}
				if !_rules[ruleProjections]() {
					goto l244
				}
				if !_rules[ruleAction7]() {
					goto l244
				}
				depth--
				add(ruleEmitProjections, position245)
			}
			return true
		l244:
			position, tokenIndex, depth = position244, tokenIndex244, depth244
			return false
		},
		/* 9 Projections <- <(<(Projection sp (',' sp Projection)*)> Action8)> */
		func() bool {
			position246, tokenIndex246, depth246 := position, tokenIndex, depth
			{
				position247 := position
				depth++
				{
					position248 := position
					depth++
					if !_rules[ruleProjection]() {
						goto l246
					}
					if !_rules[rulesp]() {
						goto l246
					}
				l249:
					{
						position250, tokenIndex250, depth250 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l250
						}
						position++
						if !_rules[rulesp]() {
							goto l250
						}
						if !_rules[ruleProjection]() {
							goto l250
						}
						goto l249
					l250:
						position, tokenIndex, depth = position250, tokenIndex250, depth250
					}
					depth--
					add(rulePegText, position248)
				}
				if !_rules[ruleAction8]() {
					goto l246
				}
				depth--
				add(ruleProjections, position247)
			}
			return true
		l246:
			position, tokenIndex, depth = position246, tokenIndex246, depth246
			return false
		},
		/* 10 Projection <- <(AliasExpression / Expression / Wildcard)> */
		func() bool {
			position251, tokenIndex251, depth251 := position, tokenIndex, depth
			{
				position252 := position
				depth++
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l254
					}
					goto l253
				l254:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if !_rules[ruleExpression]() {
						goto l255
					}
					goto l253
				l255:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
					if !_rules[ruleWildcard]() {
						goto l251
					}
				}
			l253:
				depth--
				add(ruleProjection, position252)
			}
			return true
		l251:
			position, tokenIndex, depth = position251, tokenIndex251, depth251
			return false
		},
		/* 11 AliasExpression <- <((Expression / Wildcard) sp (('a' / 'A') ('s' / 'S')) sp Identifier Action9)> */
		func() bool {
			position256, tokenIndex256, depth256 := position, tokenIndex, depth
			{
				position257 := position
				depth++
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if !_rules[ruleExpression]() {
						goto l259
					}
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if !_rules[ruleWildcard]() {
						goto l256
					}
				}
			l258:
				if !_rules[rulesp]() {
					goto l256
				}
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('A') {
						goto l256
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('S') {
						goto l256
					}
					position++
				}
			l262:
				if !_rules[rulesp]() {
					goto l256
				}
				if !_rules[ruleIdentifier]() {
					goto l256
				}
				if !_rules[ruleAction9]() {
					goto l256
				}
				depth--
				add(ruleAliasExpression, position257)
			}
			return true
		l256:
			position, tokenIndex, depth = position256, tokenIndex256, depth256
			return false
		},
		/* 12 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp WindowedRelations sp)?> Action10)> */
		func() bool {
			position264, tokenIndex264, depth264 := position, tokenIndex, depth
			{
				position265 := position
				depth++
				{
					position266 := position
					depth++
					{
						position267, tokenIndex267, depth267 := position, tokenIndex, depth
						{
							position269, tokenIndex269, depth269 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l270
							}
							position++
							goto l269
						l270:
							position, tokenIndex, depth = position269, tokenIndex269, depth269
							if buffer[position] != rune('F') {
								goto l267
							}
							position++
						}
					l269:
						{
							position271, tokenIndex271, depth271 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l272
							}
							position++
							goto l271
						l272:
							position, tokenIndex, depth = position271, tokenIndex271, depth271
							if buffer[position] != rune('R') {
								goto l267
							}
							position++
						}
					l271:
						{
							position273, tokenIndex273, depth273 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l274
							}
							position++
							goto l273
						l274:
							position, tokenIndex, depth = position273, tokenIndex273, depth273
							if buffer[position] != rune('O') {
								goto l267
							}
							position++
						}
					l273:
						{
							position275, tokenIndex275, depth275 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l276
							}
							position++
							goto l275
						l276:
							position, tokenIndex, depth = position275, tokenIndex275, depth275
							if buffer[position] != rune('M') {
								goto l267
							}
							position++
						}
					l275:
						if !_rules[rulesp]() {
							goto l267
						}
						if !_rules[ruleWindowedRelations]() {
							goto l267
						}
						if !_rules[rulesp]() {
							goto l267
						}
						goto l268
					l267:
						position, tokenIndex, depth = position267, tokenIndex267, depth267
					}
				l268:
					depth--
					add(rulePegText, position266)
				}
				if !_rules[ruleAction10]() {
					goto l264
				}
				depth--
				add(ruleWindowedFrom, position265)
			}
			return true
		l264:
			position, tokenIndex, depth = position264, tokenIndex264, depth264
			return false
		},
		/* 13 Range <- <(NumericLiteral sp RangeUnit Action11)> */
		func() bool {
			position277, tokenIndex277, depth277 := position, tokenIndex, depth
			{
				position278 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l277
				}
				if !_rules[rulesp]() {
					goto l277
				}
				if !_rules[ruleRangeUnit]() {
					goto l277
				}
				if !_rules[ruleAction11]() {
					goto l277
				}
				depth--
				add(ruleRange, position278)
			}
			return true
		l277:
			position, tokenIndex, depth = position277, tokenIndex277, depth277
			return false
		},
		/* 14 WindowedRelations <- <(WindowedRelationLike sp (',' sp WindowedRelationLike)*)> */
		func() bool {
			position279, tokenIndex279, depth279 := position, tokenIndex, depth
			{
				position280 := position
				depth++
				if !_rules[ruleWindowedRelationLike]() {
					goto l279
				}
				if !_rules[rulesp]() {
					goto l279
				}
			l281:
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l282
					}
					position++
					if !_rules[rulesp]() {
						goto l282
					}
					if !_rules[ruleWindowedRelationLike]() {
						goto l282
					}
					goto l281
				l282:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
				}
				depth--
				add(ruleWindowedRelations, position280)
			}
			return true
		l279:
			position, tokenIndex, depth = position279, tokenIndex279, depth279
			return false
		},
		/* 15 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action12)> */
		func() bool {
			position283, tokenIndex283, depth283 := position, tokenIndex, depth
			{
				position284 := position
				depth++
				{
					position285 := position
					depth++
					{
						position286, tokenIndex286, depth286 := position, tokenIndex, depth
						{
							position288, tokenIndex288, depth288 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l289
							}
							position++
							goto l288
						l289:
							position, tokenIndex, depth = position288, tokenIndex288, depth288
							if buffer[position] != rune('W') {
								goto l286
							}
							position++
						}
					l288:
						{
							position290, tokenIndex290, depth290 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l291
							}
							position++
							goto l290
						l291:
							position, tokenIndex, depth = position290, tokenIndex290, depth290
							if buffer[position] != rune('H') {
								goto l286
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
								goto l286
							}
							position++
						}
					l292:
						{
							position294, tokenIndex294, depth294 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l295
							}
							position++
							goto l294
						l295:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							if buffer[position] != rune('R') {
								goto l286
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
								goto l286
							}
							position++
						}
					l296:
						if !_rules[rulesp]() {
							goto l286
						}
						if !_rules[ruleExpression]() {
							goto l286
						}
						goto l287
					l286:
						position, tokenIndex, depth = position286, tokenIndex286, depth286
					}
				l287:
					depth--
					add(rulePegText, position285)
				}
				if !_rules[ruleAction12]() {
					goto l283
				}
				depth--
				add(ruleFilter, position284)
			}
			return true
		l283:
			position, tokenIndex, depth = position283, tokenIndex283, depth283
			return false
		},
		/* 16 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action13)> */
		func() bool {
			position298, tokenIndex298, depth298 := position, tokenIndex, depth
			{
				position299 := position
				depth++
				{
					position300 := position
					depth++
					{
						position301, tokenIndex301, depth301 := position, tokenIndex, depth
						{
							position303, tokenIndex303, depth303 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l304
							}
							position++
							goto l303
						l304:
							position, tokenIndex, depth = position303, tokenIndex303, depth303
							if buffer[position] != rune('G') {
								goto l301
							}
							position++
						}
					l303:
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l306
							}
							position++
							goto l305
						l306:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
							if buffer[position] != rune('R') {
								goto l301
							}
							position++
						}
					l305:
						{
							position307, tokenIndex307, depth307 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l308
							}
							position++
							goto l307
						l308:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if buffer[position] != rune('O') {
								goto l301
							}
							position++
						}
					l307:
						{
							position309, tokenIndex309, depth309 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l310
							}
							position++
							goto l309
						l310:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if buffer[position] != rune('U') {
								goto l301
							}
							position++
						}
					l309:
						{
							position311, tokenIndex311, depth311 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l312
							}
							position++
							goto l311
						l312:
							position, tokenIndex, depth = position311, tokenIndex311, depth311
							if buffer[position] != rune('P') {
								goto l301
							}
							position++
						}
					l311:
						if !_rules[rulesp]() {
							goto l301
						}
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l314
							}
							position++
							goto l313
						l314:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
							if buffer[position] != rune('B') {
								goto l301
							}
							position++
						}
					l313:
						{
							position315, tokenIndex315, depth315 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l316
							}
							position++
							goto l315
						l316:
							position, tokenIndex, depth = position315, tokenIndex315, depth315
							if buffer[position] != rune('Y') {
								goto l301
							}
							position++
						}
					l315:
						if !_rules[rulesp]() {
							goto l301
						}
						if !_rules[ruleGroupList]() {
							goto l301
						}
						goto l302
					l301:
						position, tokenIndex, depth = position301, tokenIndex301, depth301
					}
				l302:
					depth--
					add(rulePegText, position300)
				}
				if !_rules[ruleAction13]() {
					goto l298
				}
				depth--
				add(ruleGrouping, position299)
			}
			return true
		l298:
			position, tokenIndex, depth = position298, tokenIndex298, depth298
			return false
		},
		/* 17 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position317, tokenIndex317, depth317 := position, tokenIndex, depth
			{
				position318 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l317
				}
				if !_rules[rulesp]() {
					goto l317
				}
			l319:
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l320
					}
					position++
					if !_rules[rulesp]() {
						goto l320
					}
					if !_rules[ruleExpression]() {
						goto l320
					}
					goto l319
				l320:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
				}
				depth--
				add(ruleGroupList, position318)
			}
			return true
		l317:
			position, tokenIndex, depth = position317, tokenIndex317, depth317
			return false
		},
		/* 18 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action14)> */
		func() bool {
			position321, tokenIndex321, depth321 := position, tokenIndex, depth
			{
				position322 := position
				depth++
				{
					position323 := position
					depth++
					{
						position324, tokenIndex324, depth324 := position, tokenIndex, depth
						{
							position326, tokenIndex326, depth326 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l327
							}
							position++
							goto l326
						l327:
							position, tokenIndex, depth = position326, tokenIndex326, depth326
							if buffer[position] != rune('H') {
								goto l324
							}
							position++
						}
					l326:
						{
							position328, tokenIndex328, depth328 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l329
							}
							position++
							goto l328
						l329:
							position, tokenIndex, depth = position328, tokenIndex328, depth328
							if buffer[position] != rune('A') {
								goto l324
							}
							position++
						}
					l328:
						{
							position330, tokenIndex330, depth330 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l331
							}
							position++
							goto l330
						l331:
							position, tokenIndex, depth = position330, tokenIndex330, depth330
							if buffer[position] != rune('V') {
								goto l324
							}
							position++
						}
					l330:
						{
							position332, tokenIndex332, depth332 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l333
							}
							position++
							goto l332
						l333:
							position, tokenIndex, depth = position332, tokenIndex332, depth332
							if buffer[position] != rune('I') {
								goto l324
							}
							position++
						}
					l332:
						{
							position334, tokenIndex334, depth334 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l335
							}
							position++
							goto l334
						l335:
							position, tokenIndex, depth = position334, tokenIndex334, depth334
							if buffer[position] != rune('N') {
								goto l324
							}
							position++
						}
					l334:
						{
							position336, tokenIndex336, depth336 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l337
							}
							position++
							goto l336
						l337:
							position, tokenIndex, depth = position336, tokenIndex336, depth336
							if buffer[position] != rune('G') {
								goto l324
							}
							position++
						}
					l336:
						if !_rules[rulesp]() {
							goto l324
						}
						if !_rules[ruleExpression]() {
							goto l324
						}
						goto l325
					l324:
						position, tokenIndex, depth = position324, tokenIndex324, depth324
					}
				l325:
					depth--
					add(rulePegText, position323)
				}
				if !_rules[ruleAction14]() {
					goto l321
				}
				depth--
				add(ruleHaving, position322)
			}
			return true
		l321:
			position, tokenIndex, depth = position321, tokenIndex321, depth321
			return false
		},
		/* 19 WindowedRelationLike <- <(AliasWindowedRelation / (WindowedRelation Action15))> */
		func() bool {
			position338, tokenIndex338, depth338 := position, tokenIndex, depth
			{
				position339 := position
				depth++
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if !_rules[ruleAliasWindowedRelation]() {
						goto l341
					}
					goto l340
				l341:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
					if !_rules[ruleWindowedRelation]() {
						goto l338
					}
					if !_rules[ruleAction15]() {
						goto l338
					}
				}
			l340:
				depth--
				add(ruleWindowedRelationLike, position339)
			}
			return true
		l338:
			position, tokenIndex, depth = position338, tokenIndex338, depth338
			return false
		},
		/* 20 AliasWindowedRelation <- <(WindowedRelation sp (('a' / 'A') ('s' / 'S')) sp Identifier Action16)> */
		func() bool {
			position342, tokenIndex342, depth342 := position, tokenIndex, depth
			{
				position343 := position
				depth++
				if !_rules[ruleWindowedRelation]() {
					goto l342
				}
				if !_rules[rulesp]() {
					goto l342
				}
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
						goto l342
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					if buffer[position] != rune('S') {
						goto l342
					}
					position++
				}
			l346:
				if !_rules[rulesp]() {
					goto l342
				}
				if !_rules[ruleIdentifier]() {
					goto l342
				}
				if !_rules[ruleAction16]() {
					goto l342
				}
				depth--
				add(ruleAliasWindowedRelation, position343)
			}
			return true
		l342:
			position, tokenIndex, depth = position342, tokenIndex342, depth342
			return false
		},
		/* 21 WindowedRelation <- <(Relation sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Range sp ']' Action17)> */
		func() bool {
			position348, tokenIndex348, depth348 := position, tokenIndex, depth
			{
				position349 := position
				depth++
				if !_rules[ruleRelation]() {
					goto l348
				}
				if !_rules[rulesp]() {
					goto l348
				}
				if buffer[position] != rune('[') {
					goto l348
				}
				position++
				if !_rules[rulesp]() {
					goto l348
				}
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if buffer[position] != rune('R') {
						goto l348
					}
					position++
				}
			l350:
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('A') {
						goto l348
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
						goto l348
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
						goto l348
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
						goto l348
					}
					position++
				}
			l358:
				if !_rules[rulesp]() {
					goto l348
				}
				if !_rules[ruleRange]() {
					goto l348
				}
				if !_rules[rulesp]() {
					goto l348
				}
				if buffer[position] != rune(']') {
					goto l348
				}
				position++
				if !_rules[ruleAction17]() {
					goto l348
				}
				depth--
				add(ruleWindowedRelation, position349)
			}
			return true
		l348:
			position, tokenIndex, depth = position348, tokenIndex348, depth348
			return false
		},
		/* 22 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action18)> */
		func() bool {
			position360, tokenIndex360, depth360 := position, tokenIndex, depth
			{
				position361 := position
				depth++
				{
					position362 := position
					depth++
					{
						position363, tokenIndex363, depth363 := position, tokenIndex, depth
						{
							position365, tokenIndex365, depth365 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l366
							}
							position++
							goto l365
						l366:
							position, tokenIndex, depth = position365, tokenIndex365, depth365
							if buffer[position] != rune('W') {
								goto l363
							}
							position++
						}
					l365:
						{
							position367, tokenIndex367, depth367 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l368
							}
							position++
							goto l367
						l368:
							position, tokenIndex, depth = position367, tokenIndex367, depth367
							if buffer[position] != rune('I') {
								goto l363
							}
							position++
						}
					l367:
						{
							position369, tokenIndex369, depth369 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l370
							}
							position++
							goto l369
						l370:
							position, tokenIndex, depth = position369, tokenIndex369, depth369
							if buffer[position] != rune('T') {
								goto l363
							}
							position++
						}
					l369:
						{
							position371, tokenIndex371, depth371 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l372
							}
							position++
							goto l371
						l372:
							position, tokenIndex, depth = position371, tokenIndex371, depth371
							if buffer[position] != rune('H') {
								goto l363
							}
							position++
						}
					l371:
						if !_rules[rulesp]() {
							goto l363
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l363
						}
						if !_rules[rulesp]() {
							goto l363
						}
					l373:
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l374
							}
							position++
							if !_rules[rulesp]() {
								goto l374
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l374
							}
							goto l373
						l374:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
						}
						goto l364
					l363:
						position, tokenIndex, depth = position363, tokenIndex363, depth363
					}
				l364:
					depth--
					add(rulePegText, position362)
				}
				if !_rules[ruleAction18]() {
					goto l360
				}
				depth--
				add(ruleSourceSinkSpecs, position361)
			}
			return true
		l360:
			position, tokenIndex, depth = position360, tokenIndex360, depth360
			return false
		},
		/* 23 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action19)> */
		func() bool {
			position375, tokenIndex375, depth375 := position, tokenIndex, depth
			{
				position376 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l375
				}
				if buffer[position] != rune('=') {
					goto l375
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l375
				}
				if !_rules[ruleAction19]() {
					goto l375
				}
				depth--
				add(ruleSourceSinkParam, position376)
			}
			return true
		l375:
			position, tokenIndex, depth = position375, tokenIndex375, depth375
			return false
		},
		/* 24 Expression <- <orExpr> */
		func() bool {
			position377, tokenIndex377, depth377 := position, tokenIndex, depth
			{
				position378 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l377
				}
				depth--
				add(ruleExpression, position378)
			}
			return true
		l377:
			position, tokenIndex, depth = position377, tokenIndex377, depth377
			return false
		},
		/* 25 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action20)> */
		func() bool {
			position379, tokenIndex379, depth379 := position, tokenIndex, depth
			{
				position380 := position
				depth++
				{
					position381 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l379
					}
					if !_rules[rulesp]() {
						goto l379
					}
					{
						position382, tokenIndex382, depth382 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l382
						}
						if !_rules[rulesp]() {
							goto l382
						}
						if !_rules[ruleandExpr]() {
							goto l382
						}
						goto l383
					l382:
						position, tokenIndex, depth = position382, tokenIndex382, depth382
					}
				l383:
					depth--
					add(rulePegText, position381)
				}
				if !_rules[ruleAction20]() {
					goto l379
				}
				depth--
				add(ruleorExpr, position380)
			}
			return true
		l379:
			position, tokenIndex, depth = position379, tokenIndex379, depth379
			return false
		},
		/* 26 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action21)> */
		func() bool {
			position384, tokenIndex384, depth384 := position, tokenIndex, depth
			{
				position385 := position
				depth++
				{
					position386 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l384
					}
					if !_rules[rulesp]() {
						goto l384
					}
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l387
						}
						if !_rules[rulesp]() {
							goto l387
						}
						if !_rules[rulecomparisonExpr]() {
							goto l387
						}
						goto l388
					l387:
						position, tokenIndex, depth = position387, tokenIndex387, depth387
					}
				l388:
					depth--
					add(rulePegText, position386)
				}
				if !_rules[ruleAction21]() {
					goto l384
				}
				depth--
				add(ruleandExpr, position385)
			}
			return true
		l384:
			position, tokenIndex, depth = position384, tokenIndex384, depth384
			return false
		},
		/* 27 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action22)> */
		func() bool {
			position389, tokenIndex389, depth389 := position, tokenIndex, depth
			{
				position390 := position
				depth++
				{
					position391 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l389
					}
					if !_rules[rulesp]() {
						goto l389
					}
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l392
						}
						if !_rules[rulesp]() {
							goto l392
						}
						if !_rules[ruletermExpr]() {
							goto l392
						}
						goto l393
					l392:
						position, tokenIndex, depth = position392, tokenIndex392, depth392
					}
				l393:
					depth--
					add(rulePegText, position391)
				}
				if !_rules[ruleAction22]() {
					goto l389
				}
				depth--
				add(rulecomparisonExpr, position390)
			}
			return true
		l389:
			position, tokenIndex, depth = position389, tokenIndex389, depth389
			return false
		},
		/* 28 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action23)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				{
					position396 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l394
					}
					if !_rules[rulesp]() {
						goto l394
					}
					{
						position397, tokenIndex397, depth397 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l397
						}
						if !_rules[rulesp]() {
							goto l397
						}
						if !_rules[ruleproductExpr]() {
							goto l397
						}
						goto l398
					l397:
						position, tokenIndex, depth = position397, tokenIndex397, depth397
					}
				l398:
					depth--
					add(rulePegText, position396)
				}
				if !_rules[ruleAction23]() {
					goto l394
				}
				depth--
				add(ruletermExpr, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 29 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action24)> */
		func() bool {
			position399, tokenIndex399, depth399 := position, tokenIndex, depth
			{
				position400 := position
				depth++
				{
					position401 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l399
					}
					if !_rules[rulesp]() {
						goto l399
					}
					{
						position402, tokenIndex402, depth402 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l402
						}
						if !_rules[rulesp]() {
							goto l402
						}
						if !_rules[rulebaseExpr]() {
							goto l402
						}
						goto l403
					l402:
						position, tokenIndex, depth = position402, tokenIndex402, depth402
					}
				l403:
					depth--
					add(rulePegText, position401)
				}
				if !_rules[ruleAction24]() {
					goto l399
				}
				depth--
				add(ruleproductExpr, position400)
			}
			return true
		l399:
			position, tokenIndex, depth = position399, tokenIndex399, depth399
			return false
		},
		/* 30 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l407
					}
					position++
					if !_rules[rulesp]() {
						goto l407
					}
					if !_rules[ruleExpression]() {
						goto l407
					}
					if !_rules[rulesp]() {
						goto l407
					}
					if buffer[position] != rune(')') {
						goto l407
					}
					position++
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if !_rules[ruleBooleanLiteral]() {
						goto l408
					}
					goto l406
				l408:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if !_rules[ruleFuncApp]() {
						goto l409
					}
					goto l406
				l409:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if !_rules[ruleRowValue]() {
						goto l410
					}
					goto l406
				l410:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if !_rules[ruleLiteral]() {
						goto l404
					}
				}
			l406:
				depth--
				add(rulebaseExpr, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 31 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action25)> */
		func() bool {
			position411, tokenIndex411, depth411 := position, tokenIndex, depth
			{
				position412 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l411
				}
				if !_rules[rulesp]() {
					goto l411
				}
				if buffer[position] != rune('(') {
					goto l411
				}
				position++
				if !_rules[rulesp]() {
					goto l411
				}
				if !_rules[ruleFuncParams]() {
					goto l411
				}
				if !_rules[rulesp]() {
					goto l411
				}
				if buffer[position] != rune(')') {
					goto l411
				}
				position++
				if !_rules[ruleAction25]() {
					goto l411
				}
				depth--
				add(ruleFuncApp, position412)
			}
			return true
		l411:
			position, tokenIndex, depth = position411, tokenIndex411, depth411
			return false
		},
		/* 32 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action26)> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				{
					position415 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l413
					}
					if !_rules[rulesp]() {
						goto l413
					}
				l416:
					{
						position417, tokenIndex417, depth417 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l417
						}
						position++
						if !_rules[rulesp]() {
							goto l417
						}
						if !_rules[ruleExpression]() {
							goto l417
						}
						goto l416
					l417:
						position, tokenIndex, depth = position417, tokenIndex417, depth417
					}
					depth--
					add(rulePegText, position415)
				}
				if !_rules[ruleAction26]() {
					goto l413
				}
				depth--
				add(ruleFuncParams, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 33 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position418, tokenIndex418, depth418 := position, tokenIndex, depth
			{
				position419 := position
				depth++
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l421
					}
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if !_rules[ruleNumericLiteral]() {
						goto l422
					}
					goto l420
				l422:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if !_rules[ruleStringLiteral]() {
						goto l418
					}
				}
			l420:
				depth--
				add(ruleLiteral, position419)
			}
			return true
		l418:
			position, tokenIndex, depth = position418, tokenIndex418, depth418
			return false
		},
		/* 34 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position423, tokenIndex423, depth423 := position, tokenIndex, depth
			{
				position424 := position
				depth++
				{
					position425, tokenIndex425, depth425 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l426
					}
					goto l425
				l426:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if !_rules[ruleNotEqual]() {
						goto l427
					}
					goto l425
				l427:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if !_rules[ruleLessOrEqual]() {
						goto l428
					}
					goto l425
				l428:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if !_rules[ruleLess]() {
						goto l429
					}
					goto l425
				l429:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if !_rules[ruleGreaterOrEqual]() {
						goto l430
					}
					goto l425
				l430:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if !_rules[ruleGreater]() {
						goto l431
					}
					goto l425
				l431:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if !_rules[ruleNotEqual]() {
						goto l423
					}
				}
			l425:
				depth--
				add(ruleComparisonOp, position424)
			}
			return true
		l423:
			position, tokenIndex, depth = position423, tokenIndex423, depth423
			return false
		},
		/* 35 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l435
					}
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if !_rules[ruleMinus]() {
						goto l432
					}
				}
			l434:
				depth--
				add(rulePlusMinusOp, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 36 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l439
					}
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if !_rules[ruleDivide]() {
						goto l440
					}
					goto l438
				l440:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if !_rules[ruleModulo]() {
						goto l436
					}
				}
			l438:
				depth--
				add(ruleMultDivOp, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
			return false
		},
		/* 37 Relation <- <(<ident> Action27)> */
		func() bool {
			position441, tokenIndex441, depth441 := position, tokenIndex, depth
			{
				position442 := position
				depth++
				{
					position443 := position
					depth++
					if !_rules[ruleident]() {
						goto l441
					}
					depth--
					add(rulePegText, position443)
				}
				if !_rules[ruleAction27]() {
					goto l441
				}
				depth--
				add(ruleRelation, position442)
			}
			return true
		l441:
			position, tokenIndex, depth = position441, tokenIndex441, depth441
			return false
		},
		/* 38 RowValue <- <(<((ident '.')? ident)> Action28)> */
		func() bool {
			position444, tokenIndex444, depth444 := position, tokenIndex, depth
			{
				position445 := position
				depth++
				{
					position446 := position
					depth++
					{
						position447, tokenIndex447, depth447 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l447
						}
						if buffer[position] != rune('.') {
							goto l447
						}
						position++
						goto l448
					l447:
						position, tokenIndex, depth = position447, tokenIndex447, depth447
					}
				l448:
					if !_rules[ruleident]() {
						goto l444
					}
					depth--
					add(rulePegText, position446)
				}
				if !_rules[ruleAction28]() {
					goto l444
				}
				depth--
				add(ruleRowValue, position445)
			}
			return true
		l444:
			position, tokenIndex, depth = position444, tokenIndex444, depth444
			return false
		},
		/* 39 NumericLiteral <- <(<('-'? [0-9]+)> Action29)> */
		func() bool {
			position449, tokenIndex449, depth449 := position, tokenIndex, depth
			{
				position450 := position
				depth++
				{
					position451 := position
					depth++
					{
						position452, tokenIndex452, depth452 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l452
						}
						position++
						goto l453
					l452:
						position, tokenIndex, depth = position452, tokenIndex452, depth452
					}
				l453:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l449
					}
					position++
				l454:
					{
						position455, tokenIndex455, depth455 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l455
						}
						position++
						goto l454
					l455:
						position, tokenIndex, depth = position455, tokenIndex455, depth455
					}
					depth--
					add(rulePegText, position451)
				}
				if !_rules[ruleAction29]() {
					goto l449
				}
				depth--
				add(ruleNumericLiteral, position450)
			}
			return true
		l449:
			position, tokenIndex, depth = position449, tokenIndex449, depth449
			return false
		},
		/* 40 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action30)> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				{
					position458 := position
					depth++
					{
						position459, tokenIndex459, depth459 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l459
						}
						position++
						goto l460
					l459:
						position, tokenIndex, depth = position459, tokenIndex459, depth459
					}
				l460:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l456
					}
					position++
				l461:
					{
						position462, tokenIndex462, depth462 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l462
						}
						position++
						goto l461
					l462:
						position, tokenIndex, depth = position462, tokenIndex462, depth462
					}
					if buffer[position] != rune('.') {
						goto l456
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l456
					}
					position++
				l463:
					{
						position464, tokenIndex464, depth464 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l464
						}
						position++
						goto l463
					l464:
						position, tokenIndex, depth = position464, tokenIndex464, depth464
					}
					depth--
					add(rulePegText, position458)
				}
				if !_rules[ruleAction30]() {
					goto l456
				}
				depth--
				add(ruleFloatLiteral, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 41 Function <- <(<ident> Action31)> */
		func() bool {
			position465, tokenIndex465, depth465 := position, tokenIndex, depth
			{
				position466 := position
				depth++
				{
					position467 := position
					depth++
					if !_rules[ruleident]() {
						goto l465
					}
					depth--
					add(rulePegText, position467)
				}
				if !_rules[ruleAction31]() {
					goto l465
				}
				depth--
				add(ruleFunction, position466)
			}
			return true
		l465:
			position, tokenIndex, depth = position465, tokenIndex465, depth465
			return false
		},
		/* 42 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position468, tokenIndex468, depth468 := position, tokenIndex, depth
			{
				position469 := position
				depth++
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l471
					}
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if !_rules[ruleFALSE]() {
						goto l468
					}
				}
			l470:
				depth--
				add(ruleBooleanLiteral, position469)
			}
			return true
		l468:
			position, tokenIndex, depth = position468, tokenIndex468, depth468
			return false
		},
		/* 43 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action32)> */
		func() bool {
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				{
					position474 := position
					depth++
					{
						position475, tokenIndex475, depth475 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l476
						}
						position++
						goto l475
					l476:
						position, tokenIndex, depth = position475, tokenIndex475, depth475
						if buffer[position] != rune('T') {
							goto l472
						}
						position++
					}
				l475:
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l478
						}
						position++
						goto l477
					l478:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
						if buffer[position] != rune('R') {
							goto l472
						}
						position++
					}
				l477:
					{
						position479, tokenIndex479, depth479 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l480
						}
						position++
						goto l479
					l480:
						position, tokenIndex, depth = position479, tokenIndex479, depth479
						if buffer[position] != rune('U') {
							goto l472
						}
						position++
					}
				l479:
					{
						position481, tokenIndex481, depth481 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l482
						}
						position++
						goto l481
					l482:
						position, tokenIndex, depth = position481, tokenIndex481, depth481
						if buffer[position] != rune('E') {
							goto l472
						}
						position++
					}
				l481:
					depth--
					add(rulePegText, position474)
				}
				if !_rules[ruleAction32]() {
					goto l472
				}
				depth--
				add(ruleTRUE, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 44 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action33)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				{
					position485 := position
					depth++
					{
						position486, tokenIndex486, depth486 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l487
						}
						position++
						goto l486
					l487:
						position, tokenIndex, depth = position486, tokenIndex486, depth486
						if buffer[position] != rune('F') {
							goto l483
						}
						position++
					}
				l486:
					{
						position488, tokenIndex488, depth488 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l489
						}
						position++
						goto l488
					l489:
						position, tokenIndex, depth = position488, tokenIndex488, depth488
						if buffer[position] != rune('A') {
							goto l483
						}
						position++
					}
				l488:
					{
						position490, tokenIndex490, depth490 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l491
						}
						position++
						goto l490
					l491:
						position, tokenIndex, depth = position490, tokenIndex490, depth490
						if buffer[position] != rune('L') {
							goto l483
						}
						position++
					}
				l490:
					{
						position492, tokenIndex492, depth492 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l493
						}
						position++
						goto l492
					l493:
						position, tokenIndex, depth = position492, tokenIndex492, depth492
						if buffer[position] != rune('S') {
							goto l483
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
							goto l483
						}
						position++
					}
				l494:
					depth--
					add(rulePegText, position485)
				}
				if !_rules[ruleAction33]() {
					goto l483
				}
				depth--
				add(ruleFALSE, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 45 Wildcard <- <(<'*'> Action34)> */
		func() bool {
			position496, tokenIndex496, depth496 := position, tokenIndex, depth
			{
				position497 := position
				depth++
				{
					position498 := position
					depth++
					if buffer[position] != rune('*') {
						goto l496
					}
					position++
					depth--
					add(rulePegText, position498)
				}
				if !_rules[ruleAction34]() {
					goto l496
				}
				depth--
				add(ruleWildcard, position497)
			}
			return true
		l496:
			position, tokenIndex, depth = position496, tokenIndex496, depth496
			return false
		},
		/* 46 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action35)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				{
					position501 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l499
					}
					position++
				l502:
					{
						position503, tokenIndex503, depth503 := position, tokenIndex, depth
						{
							position504, tokenIndex504, depth504 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l505
							}
							position++
							if buffer[position] != rune('\'') {
								goto l505
							}
							position++
							goto l504
						l505:
							position, tokenIndex, depth = position504, tokenIndex504, depth504
							{
								position506, tokenIndex506, depth506 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l506
								}
								position++
								goto l503
							l506:
								position, tokenIndex, depth = position506, tokenIndex506, depth506
							}
							if !matchDot() {
								goto l503
							}
						}
					l504:
						goto l502
					l503:
						position, tokenIndex, depth = position503, tokenIndex503, depth503
					}
					if buffer[position] != rune('\'') {
						goto l499
					}
					position++
					depth--
					add(rulePegText, position501)
				}
				if !_rules[ruleAction35]() {
					goto l499
				}
				depth--
				add(ruleStringLiteral, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 47 Emitter <- <(ISTREAM / DSTREAM / RSTREAM)> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				{
					position509, tokenIndex509, depth509 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l510
					}
					goto l509
				l510:
					position, tokenIndex, depth = position509, tokenIndex509, depth509
					if !_rules[ruleDSTREAM]() {
						goto l511
					}
					goto l509
				l511:
					position, tokenIndex, depth = position509, tokenIndex509, depth509
					if !_rules[ruleRSTREAM]() {
						goto l507
					}
				}
			l509:
				depth--
				add(ruleEmitter, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 48 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action36)> */
		func() bool {
			position512, tokenIndex512, depth512 := position, tokenIndex, depth
			{
				position513 := position
				depth++
				{
					position514 := position
					depth++
					{
						position515, tokenIndex515, depth515 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l516
						}
						position++
						goto l515
					l516:
						position, tokenIndex, depth = position515, tokenIndex515, depth515
						if buffer[position] != rune('I') {
							goto l512
						}
						position++
					}
				l515:
					{
						position517, tokenIndex517, depth517 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l518
						}
						position++
						goto l517
					l518:
						position, tokenIndex, depth = position517, tokenIndex517, depth517
						if buffer[position] != rune('S') {
							goto l512
						}
						position++
					}
				l517:
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
							goto l512
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
							goto l512
						}
						position++
					}
				l521:
					{
						position523, tokenIndex523, depth523 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l524
						}
						position++
						goto l523
					l524:
						position, tokenIndex, depth = position523, tokenIndex523, depth523
						if buffer[position] != rune('E') {
							goto l512
						}
						position++
					}
				l523:
					{
						position525, tokenIndex525, depth525 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l526
						}
						position++
						goto l525
					l526:
						position, tokenIndex, depth = position525, tokenIndex525, depth525
						if buffer[position] != rune('A') {
							goto l512
						}
						position++
					}
				l525:
					{
						position527, tokenIndex527, depth527 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l528
						}
						position++
						goto l527
					l528:
						position, tokenIndex, depth = position527, tokenIndex527, depth527
						if buffer[position] != rune('M') {
							goto l512
						}
						position++
					}
				l527:
					depth--
					add(rulePegText, position514)
				}
				if !_rules[ruleAction36]() {
					goto l512
				}
				depth--
				add(ruleISTREAM, position513)
			}
			return true
		l512:
			position, tokenIndex, depth = position512, tokenIndex512, depth512
			return false
		},
		/* 49 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action37)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531 := position
					depth++
					{
						position532, tokenIndex532, depth532 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l533
						}
						position++
						goto l532
					l533:
						position, tokenIndex, depth = position532, tokenIndex532, depth532
						if buffer[position] != rune('D') {
							goto l529
						}
						position++
					}
				l532:
					{
						position534, tokenIndex534, depth534 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l535
						}
						position++
						goto l534
					l535:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if buffer[position] != rune('S') {
							goto l529
						}
						position++
					}
				l534:
					{
						position536, tokenIndex536, depth536 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l537
						}
						position++
						goto l536
					l537:
						position, tokenIndex, depth = position536, tokenIndex536, depth536
						if buffer[position] != rune('T') {
							goto l529
						}
						position++
					}
				l536:
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l539
						}
						position++
						goto l538
					l539:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
						if buffer[position] != rune('R') {
							goto l529
						}
						position++
					}
				l538:
					{
						position540, tokenIndex540, depth540 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l541
						}
						position++
						goto l540
					l541:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
						if buffer[position] != rune('E') {
							goto l529
						}
						position++
					}
				l540:
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l543
						}
						position++
						goto l542
					l543:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
						if buffer[position] != rune('A') {
							goto l529
						}
						position++
					}
				l542:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l545
						}
						position++
						goto l544
					l545:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
						if buffer[position] != rune('M') {
							goto l529
						}
						position++
					}
				l544:
					depth--
					add(rulePegText, position531)
				}
				if !_rules[ruleAction37]() {
					goto l529
				}
				depth--
				add(ruleDSTREAM, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 50 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action38)> */
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
						if buffer[position] != rune('r') {
							goto l550
						}
						position++
						goto l549
					l550:
						position, tokenIndex, depth = position549, tokenIndex549, depth549
						if buffer[position] != rune('R') {
							goto l546
						}
						position++
					}
				l549:
					{
						position551, tokenIndex551, depth551 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l552
						}
						position++
						goto l551
					l552:
						position, tokenIndex, depth = position551, tokenIndex551, depth551
						if buffer[position] != rune('S') {
							goto l546
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
							goto l546
						}
						position++
					}
				l553:
					{
						position555, tokenIndex555, depth555 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l556
						}
						position++
						goto l555
					l556:
						position, tokenIndex, depth = position555, tokenIndex555, depth555
						if buffer[position] != rune('R') {
							goto l546
						}
						position++
					}
				l555:
					{
						position557, tokenIndex557, depth557 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l558
						}
						position++
						goto l557
					l558:
						position, tokenIndex, depth = position557, tokenIndex557, depth557
						if buffer[position] != rune('E') {
							goto l546
						}
						position++
					}
				l557:
					{
						position559, tokenIndex559, depth559 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l560
						}
						position++
						goto l559
					l560:
						position, tokenIndex, depth = position559, tokenIndex559, depth559
						if buffer[position] != rune('A') {
							goto l546
						}
						position++
					}
				l559:
					{
						position561, tokenIndex561, depth561 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l562
						}
						position++
						goto l561
					l562:
						position, tokenIndex, depth = position561, tokenIndex561, depth561
						if buffer[position] != rune('M') {
							goto l546
						}
						position++
					}
				l561:
					depth--
					add(rulePegText, position548)
				}
				if !_rules[ruleAction38]() {
					goto l546
				}
				depth--
				add(ruleRSTREAM, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 51 RangeUnit <- <(TUPLES / SECONDS)> */
		func() bool {
			position563, tokenIndex563, depth563 := position, tokenIndex, depth
			{
				position564 := position
				depth++
				{
					position565, tokenIndex565, depth565 := position, tokenIndex, depth
					if !_rules[ruleTUPLES]() {
						goto l566
					}
					goto l565
				l566:
					position, tokenIndex, depth = position565, tokenIndex565, depth565
					if !_rules[ruleSECONDS]() {
						goto l563
					}
				}
			l565:
				depth--
				add(ruleRangeUnit, position564)
			}
			return true
		l563:
			position, tokenIndex, depth = position563, tokenIndex563, depth563
			return false
		},
		/* 52 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action39)> */
		func() bool {
			position567, tokenIndex567, depth567 := position, tokenIndex, depth
			{
				position568 := position
				depth++
				{
					position569 := position
					depth++
					{
						position570, tokenIndex570, depth570 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l571
						}
						position++
						goto l570
					l571:
						position, tokenIndex, depth = position570, tokenIndex570, depth570
						if buffer[position] != rune('T') {
							goto l567
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
							goto l567
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
							goto l567
						}
						position++
					}
				l574:
					{
						position576, tokenIndex576, depth576 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l577
						}
						position++
						goto l576
					l577:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						if buffer[position] != rune('L') {
							goto l567
						}
						position++
					}
				l576:
					{
						position578, tokenIndex578, depth578 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l579
						}
						position++
						goto l578
					l579:
						position, tokenIndex, depth = position578, tokenIndex578, depth578
						if buffer[position] != rune('E') {
							goto l567
						}
						position++
					}
				l578:
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l581
						}
						position++
						goto l580
					l581:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
						if buffer[position] != rune('S') {
							goto l567
						}
						position++
					}
				l580:
					depth--
					add(rulePegText, position569)
				}
				if !_rules[ruleAction39]() {
					goto l567
				}
				depth--
				add(ruleTUPLES, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 53 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action40)> */
		func() bool {
			position582, tokenIndex582, depth582 := position, tokenIndex, depth
			{
				position583 := position
				depth++
				{
					position584 := position
					depth++
					{
						position585, tokenIndex585, depth585 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l586
						}
						position++
						goto l585
					l586:
						position, tokenIndex, depth = position585, tokenIndex585, depth585
						if buffer[position] != rune('S') {
							goto l582
						}
						position++
					}
				l585:
					{
						position587, tokenIndex587, depth587 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l588
						}
						position++
						goto l587
					l588:
						position, tokenIndex, depth = position587, tokenIndex587, depth587
						if buffer[position] != rune('E') {
							goto l582
						}
						position++
					}
				l587:
					{
						position589, tokenIndex589, depth589 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l590
						}
						position++
						goto l589
					l590:
						position, tokenIndex, depth = position589, tokenIndex589, depth589
						if buffer[position] != rune('C') {
							goto l582
						}
						position++
					}
				l589:
					{
						position591, tokenIndex591, depth591 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l592
						}
						position++
						goto l591
					l592:
						position, tokenIndex, depth = position591, tokenIndex591, depth591
						if buffer[position] != rune('O') {
							goto l582
						}
						position++
					}
				l591:
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l594
						}
						position++
						goto l593
					l594:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
						if buffer[position] != rune('N') {
							goto l582
						}
						position++
					}
				l593:
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l596
						}
						position++
						goto l595
					l596:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
						if buffer[position] != rune('D') {
							goto l582
						}
						position++
					}
				l595:
					{
						position597, tokenIndex597, depth597 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l598
						}
						position++
						goto l597
					l598:
						position, tokenIndex, depth = position597, tokenIndex597, depth597
						if buffer[position] != rune('S') {
							goto l582
						}
						position++
					}
				l597:
					depth--
					add(rulePegText, position584)
				}
				if !_rules[ruleAction40]() {
					goto l582
				}
				depth--
				add(ruleSECONDS, position583)
			}
			return true
		l582:
			position, tokenIndex, depth = position582, tokenIndex582, depth582
			return false
		},
		/* 54 SourceSinkName <- <(<ident> Action41)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				{
					position601 := position
					depth++
					if !_rules[ruleident]() {
						goto l599
					}
					depth--
					add(rulePegText, position601)
				}
				if !_rules[ruleAction41]() {
					goto l599
				}
				depth--
				add(ruleSourceSinkName, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 55 SourceSinkType <- <(<ident> Action42)> */
		func() bool {
			position602, tokenIndex602, depth602 := position, tokenIndex, depth
			{
				position603 := position
				depth++
				{
					position604 := position
					depth++
					if !_rules[ruleident]() {
						goto l602
					}
					depth--
					add(rulePegText, position604)
				}
				if !_rules[ruleAction42]() {
					goto l602
				}
				depth--
				add(ruleSourceSinkType, position603)
			}
			return true
		l602:
			position, tokenIndex, depth = position602, tokenIndex602, depth602
			return false
		},
		/* 56 SourceSinkParamKey <- <(<ident> Action43)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				{
					position607 := position
					depth++
					if !_rules[ruleident]() {
						goto l605
					}
					depth--
					add(rulePegText, position607)
				}
				if !_rules[ruleAction43]() {
					goto l605
				}
				depth--
				add(ruleSourceSinkParamKey, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 57 SourceSinkParamVal <- <(<([a-z] / [A-Z] / [0-9] / '_')+> Action44)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l615
						}
						position++
						goto l613
					l615:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l616
						}
						position++
						goto l613
					l616:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('_') {
							goto l608
						}
						position++
					}
				l613:
				l611:
					{
						position612, tokenIndex612, depth612 := position, tokenIndex, depth
						{
							position617, tokenIndex617, depth617 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l618
							}
							position++
							goto l617
						l618:
							position, tokenIndex, depth = position617, tokenIndex617, depth617
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l619
							}
							position++
							goto l617
						l619:
							position, tokenIndex, depth = position617, tokenIndex617, depth617
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l620
							}
							position++
							goto l617
						l620:
							position, tokenIndex, depth = position617, tokenIndex617, depth617
							if buffer[position] != rune('_') {
								goto l612
							}
							position++
						}
					l617:
						goto l611
					l612:
						position, tokenIndex, depth = position612, tokenIndex612, depth612
					}
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction44]() {
					goto l608
				}
				depth--
				add(ruleSourceSinkParamVal, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 58 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action45)> */
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
						if buffer[position] != rune('o') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('O') {
							goto l621
						}
						position++
					}
				l624:
					{
						position626, tokenIndex626, depth626 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('R') {
							goto l621
						}
						position++
					}
				l626:
					depth--
					add(rulePegText, position623)
				}
				if !_rules[ruleAction45]() {
					goto l621
				}
				depth--
				add(ruleOr, position622)
			}
			return true
		l621:
			position, tokenIndex, depth = position621, tokenIndex621, depth621
			return false
		},
		/* 59 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action46)> */
		func() bool {
			position628, tokenIndex628, depth628 := position, tokenIndex, depth
			{
				position629 := position
				depth++
				{
					position630 := position
					depth++
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
							goto l628
						}
						position++
					}
				l631:
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l634
						}
						position++
						goto l633
					l634:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
						if buffer[position] != rune('N') {
							goto l628
						}
						position++
					}
				l633:
					{
						position635, tokenIndex635, depth635 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l636
						}
						position++
						goto l635
					l636:
						position, tokenIndex, depth = position635, tokenIndex635, depth635
						if buffer[position] != rune('D') {
							goto l628
						}
						position++
					}
				l635:
					depth--
					add(rulePegText, position630)
				}
				if !_rules[ruleAction46]() {
					goto l628
				}
				depth--
				add(ruleAnd, position629)
			}
			return true
		l628:
			position, tokenIndex, depth = position628, tokenIndex628, depth628
			return false
		},
		/* 60 Equal <- <(<'='> Action47)> */
		func() bool {
			position637, tokenIndex637, depth637 := position, tokenIndex, depth
			{
				position638 := position
				depth++
				{
					position639 := position
					depth++
					if buffer[position] != rune('=') {
						goto l637
					}
					position++
					depth--
					add(rulePegText, position639)
				}
				if !_rules[ruleAction47]() {
					goto l637
				}
				depth--
				add(ruleEqual, position638)
			}
			return true
		l637:
			position, tokenIndex, depth = position637, tokenIndex637, depth637
			return false
		},
		/* 61 Less <- <(<'<'> Action48)> */
		func() bool {
			position640, tokenIndex640, depth640 := position, tokenIndex, depth
			{
				position641 := position
				depth++
				{
					position642 := position
					depth++
					if buffer[position] != rune('<') {
						goto l640
					}
					position++
					depth--
					add(rulePegText, position642)
				}
				if !_rules[ruleAction48]() {
					goto l640
				}
				depth--
				add(ruleLess, position641)
			}
			return true
		l640:
			position, tokenIndex, depth = position640, tokenIndex640, depth640
			return false
		},
		/* 62 LessOrEqual <- <(<('<' '=')> Action49)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645 := position
					depth++
					if buffer[position] != rune('<') {
						goto l643
					}
					position++
					if buffer[position] != rune('=') {
						goto l643
					}
					position++
					depth--
					add(rulePegText, position645)
				}
				if !_rules[ruleAction49]() {
					goto l643
				}
				depth--
				add(ruleLessOrEqual, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 63 Greater <- <(<'>'> Action50)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				{
					position648 := position
					depth++
					if buffer[position] != rune('>') {
						goto l646
					}
					position++
					depth--
					add(rulePegText, position648)
				}
				if !_rules[ruleAction50]() {
					goto l646
				}
				depth--
				add(ruleGreater, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 64 GreaterOrEqual <- <(<('>' '=')> Action51)> */
		func() bool {
			position649, tokenIndex649, depth649 := position, tokenIndex, depth
			{
				position650 := position
				depth++
				{
					position651 := position
					depth++
					if buffer[position] != rune('>') {
						goto l649
					}
					position++
					if buffer[position] != rune('=') {
						goto l649
					}
					position++
					depth--
					add(rulePegText, position651)
				}
				if !_rules[ruleAction51]() {
					goto l649
				}
				depth--
				add(ruleGreaterOrEqual, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 65 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action52)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				{
					position654 := position
					depth++
					{
						position655, tokenIndex655, depth655 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l656
						}
						position++
						if buffer[position] != rune('=') {
							goto l656
						}
						position++
						goto l655
					l656:
						position, tokenIndex, depth = position655, tokenIndex655, depth655
						if buffer[position] != rune('<') {
							goto l652
						}
						position++
						if buffer[position] != rune('>') {
							goto l652
						}
						position++
					}
				l655:
					depth--
					add(rulePegText, position654)
				}
				if !_rules[ruleAction52]() {
					goto l652
				}
				depth--
				add(ruleNotEqual, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 66 Plus <- <(<'+'> Action53)> */
		func() bool {
			position657, tokenIndex657, depth657 := position, tokenIndex, depth
			{
				position658 := position
				depth++
				{
					position659 := position
					depth++
					if buffer[position] != rune('+') {
						goto l657
					}
					position++
					depth--
					add(rulePegText, position659)
				}
				if !_rules[ruleAction53]() {
					goto l657
				}
				depth--
				add(rulePlus, position658)
			}
			return true
		l657:
			position, tokenIndex, depth = position657, tokenIndex657, depth657
			return false
		},
		/* 67 Minus <- <(<'-'> Action54)> */
		func() bool {
			position660, tokenIndex660, depth660 := position, tokenIndex, depth
			{
				position661 := position
				depth++
				{
					position662 := position
					depth++
					if buffer[position] != rune('-') {
						goto l660
					}
					position++
					depth--
					add(rulePegText, position662)
				}
				if !_rules[ruleAction54]() {
					goto l660
				}
				depth--
				add(ruleMinus, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 68 Multiply <- <(<'*'> Action55)> */
		func() bool {
			position663, tokenIndex663, depth663 := position, tokenIndex, depth
			{
				position664 := position
				depth++
				{
					position665 := position
					depth++
					if buffer[position] != rune('*') {
						goto l663
					}
					position++
					depth--
					add(rulePegText, position665)
				}
				if !_rules[ruleAction55]() {
					goto l663
				}
				depth--
				add(ruleMultiply, position664)
			}
			return true
		l663:
			position, tokenIndex, depth = position663, tokenIndex663, depth663
			return false
		},
		/* 69 Divide <- <(<'/'> Action56)> */
		func() bool {
			position666, tokenIndex666, depth666 := position, tokenIndex, depth
			{
				position667 := position
				depth++
				{
					position668 := position
					depth++
					if buffer[position] != rune('/') {
						goto l666
					}
					position++
					depth--
					add(rulePegText, position668)
				}
				if !_rules[ruleAction56]() {
					goto l666
				}
				depth--
				add(ruleDivide, position667)
			}
			return true
		l666:
			position, tokenIndex, depth = position666, tokenIndex666, depth666
			return false
		},
		/* 70 Modulo <- <(<'%'> Action57)> */
		func() bool {
			position669, tokenIndex669, depth669 := position, tokenIndex, depth
			{
				position670 := position
				depth++
				{
					position671 := position
					depth++
					if buffer[position] != rune('%') {
						goto l669
					}
					position++
					depth--
					add(rulePegText, position671)
				}
				if !_rules[ruleAction57]() {
					goto l669
				}
				depth--
				add(ruleModulo, position670)
			}
			return true
		l669:
			position, tokenIndex, depth = position669, tokenIndex669, depth669
			return false
		},
		/* 71 Identifier <- <(<ident> Action58)> */
		func() bool {
			position672, tokenIndex672, depth672 := position, tokenIndex, depth
			{
				position673 := position
				depth++
				{
					position674 := position
					depth++
					if !_rules[ruleident]() {
						goto l672
					}
					depth--
					add(rulePegText, position674)
				}
				if !_rules[ruleAction58]() {
					goto l672
				}
				depth--
				add(ruleIdentifier, position673)
			}
			return true
		l672:
			position, tokenIndex, depth = position672, tokenIndex672, depth672
			return false
		},
		/* 72 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position675, tokenIndex675, depth675 := position, tokenIndex, depth
			{
				position676 := position
				depth++
				{
					position677, tokenIndex677, depth677 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l678
					}
					position++
					goto l677
				l678:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l675
					}
					position++
				}
			l677:
			l679:
				{
					position680, tokenIndex680, depth680 := position, tokenIndex, depth
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l683
						}
						position++
						goto l681
					l683:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l684
						}
						position++
						goto l681
					l684:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if buffer[position] != rune('_') {
							goto l680
						}
						position++
					}
				l681:
					goto l679
				l680:
					position, tokenIndex, depth = position680, tokenIndex680, depth680
				}
				depth--
				add(ruleident, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 73 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position686 := position
				depth++
			l687:
				{
					position688, tokenIndex688, depth688 := position, tokenIndex, depth
					{
						position689, tokenIndex689, depth689 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l690
						}
						position++
						goto l689
					l690:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
						if buffer[position] != rune('\t') {
							goto l691
						}
						position++
						goto l689
					l691:
						position, tokenIndex, depth = position689, tokenIndex689, depth689
						if buffer[position] != rune('\n') {
							goto l688
						}
						position++
					}
				l689:
					goto l687
				l688:
					position, tokenIndex, depth = position688, tokenIndex688, depth688
				}
				depth--
				add(rulesp, position686)
			}
			return true
		},
		/* 75 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 76 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 77 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 78 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 79 Action4 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 80 Action5 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 81 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 82 Action7 <- <{
		    p.AssembleEmitProjections()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 84 Action8 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 85 Action9 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 86 Action10 <- <{
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
		/* 87 Action11 <- <{
		    p.AssembleRange()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 88 Action12 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 89 Action13 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 90 Action14 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 91 Action15 <- <{
		    p.EnsureAliasWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 92 Action16 <- <{
		    p.AssembleAliasWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 93 Action17 <- <{
		    p.AssembleWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 94 Action18 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 95 Action19 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 96 Action20 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 97 Action21 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 98 Action22 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 99 Action23 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 100 Action24 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 101 Action25 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 102 Action26 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 103 Action27 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRelation(substr))
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 104 Action28 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 105 Action29 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 106 Action30 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 107 Action31 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 108 Action32 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 109 Action33 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 110 Action34 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 111 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 112 Action36 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 113 Action37 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 114 Action38 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 115 Action39 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 116 Action40 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 117 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkName(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 118 Action42 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 119 Action43 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 120 Action44 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamVal(substr))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 121 Action45 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 122 Action46 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 123 Action47 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 124 Action48 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 125 Action49 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 126 Action50 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 127 Action51 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 128 Action52 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 129 Action53 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 130 Action54 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 131 Action55 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 132 Action56 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 133 Action57 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 134 Action58 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
	}
	p.rules = _rules
}

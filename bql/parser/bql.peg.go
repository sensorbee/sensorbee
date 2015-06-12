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
	ruleDefWindowedFrom
	ruleRange
	ruleWindowedRelations
	ruleDefWindowedRelations
	ruleFilter
	ruleGrouping
	ruleGroupList
	ruleHaving
	ruleWindowedRelationLike
	ruleDefWindowedRelationLike
	ruleAliasWindowedRelation
	ruleDefAliasWindowedRelation
	ruleWindowedRelation
	ruleDefWindowedRelation
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
	"EmitProjections",
	"Projections",
	"Projection",
	"AliasExpression",
	"WindowedFrom",
	"DefWindowedFrom",
	"Range",
	"WindowedRelations",
	"DefWindowedRelations",
	"Filter",
	"Grouping",
	"GroupList",
	"Having",
	"WindowedRelationLike",
	"DefWindowedRelationLike",
	"AliasWindowedRelation",
	"DefAliasWindowedRelation",
	"WindowedRelation",
	"DefWindowedRelation",
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
	rules  [144]func() bool
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

			p.AssembleWindowedFrom(begin, end)

		case ruleAction12:

			p.AssembleRange()

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

			p.EnsureAliasWindowedRelation()

		case ruleAction17:

			p.EnsureAliasWindowedRelation()

		case ruleAction18:

			p.AssembleAliasWindowedRelation()

		case ruleAction19:

			p.AssembleAliasWindowedRelation()

		case ruleAction20:

			p.AssembleWindowedRelation()

		case ruleAction21:

			p.AssembleWindowedRelation()

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
			p.PushComponent(begin, end, NewRelation(substr))

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
		/* 2 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp (('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T')) sp EmitProjections sp WindowedFrom sp Filter sp Grouping sp Having sp Action1)> */
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
		/* 13 DefWindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp DefWindowedRelations sp)?> Action11)> */
		func() bool {
			position277, tokenIndex277, depth277 := position, tokenIndex, depth
			{
				position278 := position
				depth++
				{
					position279 := position
					depth++
					{
						position280, tokenIndex280, depth280 := position, tokenIndex, depth
						{
							position282, tokenIndex282, depth282 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l283
							}
							position++
							goto l282
						l283:
							position, tokenIndex, depth = position282, tokenIndex282, depth282
							if buffer[position] != rune('F') {
								goto l280
							}
							position++
						}
					l282:
						{
							position284, tokenIndex284, depth284 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l285
							}
							position++
							goto l284
						l285:
							position, tokenIndex, depth = position284, tokenIndex284, depth284
							if buffer[position] != rune('R') {
								goto l280
							}
							position++
						}
					l284:
						{
							position286, tokenIndex286, depth286 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l287
							}
							position++
							goto l286
						l287:
							position, tokenIndex, depth = position286, tokenIndex286, depth286
							if buffer[position] != rune('O') {
								goto l280
							}
							position++
						}
					l286:
						{
							position288, tokenIndex288, depth288 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l289
							}
							position++
							goto l288
						l289:
							position, tokenIndex, depth = position288, tokenIndex288, depth288
							if buffer[position] != rune('M') {
								goto l280
							}
							position++
						}
					l288:
						if !_rules[rulesp]() {
							goto l280
						}
						if !_rules[ruleDefWindowedRelations]() {
							goto l280
						}
						if !_rules[rulesp]() {
							goto l280
						}
						goto l281
					l280:
						position, tokenIndex, depth = position280, tokenIndex280, depth280
					}
				l281:
					depth--
					add(rulePegText, position279)
				}
				if !_rules[ruleAction11]() {
					goto l277
				}
				depth--
				add(ruleDefWindowedFrom, position278)
			}
			return true
		l277:
			position, tokenIndex, depth = position277, tokenIndex277, depth277
			return false
		},
		/* 14 Range <- <(NumericLiteral sp RangeUnit Action12)> */
		func() bool {
			position290, tokenIndex290, depth290 := position, tokenIndex, depth
			{
				position291 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l290
				}
				if !_rules[rulesp]() {
					goto l290
				}
				if !_rules[ruleRangeUnit]() {
					goto l290
				}
				if !_rules[ruleAction12]() {
					goto l290
				}
				depth--
				add(ruleRange, position291)
			}
			return true
		l290:
			position, tokenIndex, depth = position290, tokenIndex290, depth290
			return false
		},
		/* 15 WindowedRelations <- <(WindowedRelationLike sp (',' sp WindowedRelationLike)*)> */
		func() bool {
			position292, tokenIndex292, depth292 := position, tokenIndex, depth
			{
				position293 := position
				depth++
				if !_rules[ruleWindowedRelationLike]() {
					goto l292
				}
				if !_rules[rulesp]() {
					goto l292
				}
			l294:
				{
					position295, tokenIndex295, depth295 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l295
					}
					position++
					if !_rules[rulesp]() {
						goto l295
					}
					if !_rules[ruleWindowedRelationLike]() {
						goto l295
					}
					goto l294
				l295:
					position, tokenIndex, depth = position295, tokenIndex295, depth295
				}
				depth--
				add(ruleWindowedRelations, position293)
			}
			return true
		l292:
			position, tokenIndex, depth = position292, tokenIndex292, depth292
			return false
		},
		/* 16 DefWindowedRelations <- <(DefWindowedRelationLike sp (',' sp DefWindowedRelationLike)*)> */
		func() bool {
			position296, tokenIndex296, depth296 := position, tokenIndex, depth
			{
				position297 := position
				depth++
				if !_rules[ruleDefWindowedRelationLike]() {
					goto l296
				}
				if !_rules[rulesp]() {
					goto l296
				}
			l298:
				{
					position299, tokenIndex299, depth299 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l299
					}
					position++
					if !_rules[rulesp]() {
						goto l299
					}
					if !_rules[ruleDefWindowedRelationLike]() {
						goto l299
					}
					goto l298
				l299:
					position, tokenIndex, depth = position299, tokenIndex299, depth299
				}
				depth--
				add(ruleDefWindowedRelations, position297)
			}
			return true
		l296:
			position, tokenIndex, depth = position296, tokenIndex296, depth296
			return false
		},
		/* 17 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action13)> */
		func() bool {
			position300, tokenIndex300, depth300 := position, tokenIndex, depth
			{
				position301 := position
				depth++
				{
					position302 := position
					depth++
					{
						position303, tokenIndex303, depth303 := position, tokenIndex, depth
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l306
							}
							position++
							goto l305
						l306:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
							if buffer[position] != rune('W') {
								goto l303
							}
							position++
						}
					l305:
						{
							position307, tokenIndex307, depth307 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l308
							}
							position++
							goto l307
						l308:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if buffer[position] != rune('H') {
								goto l303
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
								goto l303
							}
							position++
						}
					l309:
						{
							position311, tokenIndex311, depth311 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l312
							}
							position++
							goto l311
						l312:
							position, tokenIndex, depth = position311, tokenIndex311, depth311
							if buffer[position] != rune('R') {
								goto l303
							}
							position++
						}
					l311:
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l314
							}
							position++
							goto l313
						l314:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
							if buffer[position] != rune('E') {
								goto l303
							}
							position++
						}
					l313:
						if !_rules[rulesp]() {
							goto l303
						}
						if !_rules[ruleExpression]() {
							goto l303
						}
						goto l304
					l303:
						position, tokenIndex, depth = position303, tokenIndex303, depth303
					}
				l304:
					depth--
					add(rulePegText, position302)
				}
				if !_rules[ruleAction13]() {
					goto l300
				}
				depth--
				add(ruleFilter, position301)
			}
			return true
		l300:
			position, tokenIndex, depth = position300, tokenIndex300, depth300
			return false
		},
		/* 18 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action14)> */
		func() bool {
			position315, tokenIndex315, depth315 := position, tokenIndex, depth
			{
				position316 := position
				depth++
				{
					position317 := position
					depth++
					{
						position318, tokenIndex318, depth318 := position, tokenIndex, depth
						{
							position320, tokenIndex320, depth320 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l321
							}
							position++
							goto l320
						l321:
							position, tokenIndex, depth = position320, tokenIndex320, depth320
							if buffer[position] != rune('G') {
								goto l318
							}
							position++
						}
					l320:
						{
							position322, tokenIndex322, depth322 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l323
							}
							position++
							goto l322
						l323:
							position, tokenIndex, depth = position322, tokenIndex322, depth322
							if buffer[position] != rune('R') {
								goto l318
							}
							position++
						}
					l322:
						{
							position324, tokenIndex324, depth324 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l325
							}
							position++
							goto l324
						l325:
							position, tokenIndex, depth = position324, tokenIndex324, depth324
							if buffer[position] != rune('O') {
								goto l318
							}
							position++
						}
					l324:
						{
							position326, tokenIndex326, depth326 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l327
							}
							position++
							goto l326
						l327:
							position, tokenIndex, depth = position326, tokenIndex326, depth326
							if buffer[position] != rune('U') {
								goto l318
							}
							position++
						}
					l326:
						{
							position328, tokenIndex328, depth328 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l329
							}
							position++
							goto l328
						l329:
							position, tokenIndex, depth = position328, tokenIndex328, depth328
							if buffer[position] != rune('P') {
								goto l318
							}
							position++
						}
					l328:
						if !_rules[rulesp]() {
							goto l318
						}
						{
							position330, tokenIndex330, depth330 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l331
							}
							position++
							goto l330
						l331:
							position, tokenIndex, depth = position330, tokenIndex330, depth330
							if buffer[position] != rune('B') {
								goto l318
							}
							position++
						}
					l330:
						{
							position332, tokenIndex332, depth332 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l333
							}
							position++
							goto l332
						l333:
							position, tokenIndex, depth = position332, tokenIndex332, depth332
							if buffer[position] != rune('Y') {
								goto l318
							}
							position++
						}
					l332:
						if !_rules[rulesp]() {
							goto l318
						}
						if !_rules[ruleGroupList]() {
							goto l318
						}
						goto l319
					l318:
						position, tokenIndex, depth = position318, tokenIndex318, depth318
					}
				l319:
					depth--
					add(rulePegText, position317)
				}
				if !_rules[ruleAction14]() {
					goto l315
				}
				depth--
				add(ruleGrouping, position316)
			}
			return true
		l315:
			position, tokenIndex, depth = position315, tokenIndex315, depth315
			return false
		},
		/* 19 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position334, tokenIndex334, depth334 := position, tokenIndex, depth
			{
				position335 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l334
				}
				if !_rules[rulesp]() {
					goto l334
				}
			l336:
				{
					position337, tokenIndex337, depth337 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l337
					}
					position++
					if !_rules[rulesp]() {
						goto l337
					}
					if !_rules[ruleExpression]() {
						goto l337
					}
					goto l336
				l337:
					position, tokenIndex, depth = position337, tokenIndex337, depth337
				}
				depth--
				add(ruleGroupList, position335)
			}
			return true
		l334:
			position, tokenIndex, depth = position334, tokenIndex334, depth334
			return false
		},
		/* 20 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action15)> */
		func() bool {
			position338, tokenIndex338, depth338 := position, tokenIndex, depth
			{
				position339 := position
				depth++
				{
					position340 := position
					depth++
					{
						position341, tokenIndex341, depth341 := position, tokenIndex, depth
						{
							position343, tokenIndex343, depth343 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l344
							}
							position++
							goto l343
						l344:
							position, tokenIndex, depth = position343, tokenIndex343, depth343
							if buffer[position] != rune('H') {
								goto l341
							}
							position++
						}
					l343:
						{
							position345, tokenIndex345, depth345 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l346
							}
							position++
							goto l345
						l346:
							position, tokenIndex, depth = position345, tokenIndex345, depth345
							if buffer[position] != rune('A') {
								goto l341
							}
							position++
						}
					l345:
						{
							position347, tokenIndex347, depth347 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l348
							}
							position++
							goto l347
						l348:
							position, tokenIndex, depth = position347, tokenIndex347, depth347
							if buffer[position] != rune('V') {
								goto l341
							}
							position++
						}
					l347:
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
								goto l341
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
								goto l341
							}
							position++
						}
					l351:
						{
							position353, tokenIndex353, depth353 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l354
							}
							position++
							goto l353
						l354:
							position, tokenIndex, depth = position353, tokenIndex353, depth353
							if buffer[position] != rune('G') {
								goto l341
							}
							position++
						}
					l353:
						if !_rules[rulesp]() {
							goto l341
						}
						if !_rules[ruleExpression]() {
							goto l341
						}
						goto l342
					l341:
						position, tokenIndex, depth = position341, tokenIndex341, depth341
					}
				l342:
					depth--
					add(rulePegText, position340)
				}
				if !_rules[ruleAction15]() {
					goto l338
				}
				depth--
				add(ruleHaving, position339)
			}
			return true
		l338:
			position, tokenIndex, depth = position338, tokenIndex338, depth338
			return false
		},
		/* 21 WindowedRelationLike <- <(AliasWindowedRelation / (WindowedRelation Action16))> */
		func() bool {
			position355, tokenIndex355, depth355 := position, tokenIndex, depth
			{
				position356 := position
				depth++
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if !_rules[ruleAliasWindowedRelation]() {
						goto l358
					}
					goto l357
				l358:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
					if !_rules[ruleWindowedRelation]() {
						goto l355
					}
					if !_rules[ruleAction16]() {
						goto l355
					}
				}
			l357:
				depth--
				add(ruleWindowedRelationLike, position356)
			}
			return true
		l355:
			position, tokenIndex, depth = position355, tokenIndex355, depth355
			return false
		},
		/* 22 DefWindowedRelationLike <- <(DefAliasWindowedRelation / (DefWindowedRelation Action17))> */
		func() bool {
			position359, tokenIndex359, depth359 := position, tokenIndex, depth
			{
				position360 := position
				depth++
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if !_rules[ruleDefAliasWindowedRelation]() {
						goto l362
					}
					goto l361
				l362:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
					if !_rules[ruleDefWindowedRelation]() {
						goto l359
					}
					if !_rules[ruleAction17]() {
						goto l359
					}
				}
			l361:
				depth--
				add(ruleDefWindowedRelationLike, position360)
			}
			return true
		l359:
			position, tokenIndex, depth = position359, tokenIndex359, depth359
			return false
		},
		/* 23 AliasWindowedRelation <- <(WindowedRelation sp (('a' / 'A') ('s' / 'S')) sp Identifier Action18)> */
		func() bool {
			position363, tokenIndex363, depth363 := position, tokenIndex, depth
			{
				position364 := position
				depth++
				if !_rules[ruleWindowedRelation]() {
					goto l363
				}
				if !_rules[rulesp]() {
					goto l363
				}
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if buffer[position] != rune('A') {
						goto l363
					}
					position++
				}
			l365:
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l368
					}
					position++
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					if buffer[position] != rune('S') {
						goto l363
					}
					position++
				}
			l367:
				if !_rules[rulesp]() {
					goto l363
				}
				if !_rules[ruleIdentifier]() {
					goto l363
				}
				if !_rules[ruleAction18]() {
					goto l363
				}
				depth--
				add(ruleAliasWindowedRelation, position364)
			}
			return true
		l363:
			position, tokenIndex, depth = position363, tokenIndex363, depth363
			return false
		},
		/* 24 DefAliasWindowedRelation <- <(DefWindowedRelation sp (('a' / 'A') ('s' / 'S')) sp Identifier Action19)> */
		func() bool {
			position369, tokenIndex369, depth369 := position, tokenIndex, depth
			{
				position370 := position
				depth++
				if !_rules[ruleDefWindowedRelation]() {
					goto l369
				}
				if !_rules[rulesp]() {
					goto l369
				}
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l372
					}
					position++
					goto l371
				l372:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
					if buffer[position] != rune('A') {
						goto l369
					}
					position++
				}
			l371:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('S') {
						goto l369
					}
					position++
				}
			l373:
				if !_rules[rulesp]() {
					goto l369
				}
				if !_rules[ruleIdentifier]() {
					goto l369
				}
				if !_rules[ruleAction19]() {
					goto l369
				}
				depth--
				add(ruleDefAliasWindowedRelation, position370)
			}
			return true
		l369:
			position, tokenIndex, depth = position369, tokenIndex369, depth369
			return false
		},
		/* 25 WindowedRelation <- <(Relation sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Range sp ']' Action20)> */
		func() bool {
			position375, tokenIndex375, depth375 := position, tokenIndex, depth
			{
				position376 := position
				depth++
				if !_rules[ruleRelation]() {
					goto l375
				}
				if !_rules[rulesp]() {
					goto l375
				}
				if buffer[position] != rune('[') {
					goto l375
				}
				position++
				if !_rules[rulesp]() {
					goto l375
				}
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l378
					}
					position++
					goto l377
				l378:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
					if buffer[position] != rune('R') {
						goto l375
					}
					position++
				}
			l377:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('A') {
						goto l375
					}
					position++
				}
			l379:
				{
					position381, tokenIndex381, depth381 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l382
					}
					position++
					goto l381
				l382:
					position, tokenIndex, depth = position381, tokenIndex381, depth381
					if buffer[position] != rune('N') {
						goto l375
					}
					position++
				}
			l381:
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l384
					}
					position++
					goto l383
				l384:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
					if buffer[position] != rune('G') {
						goto l375
					}
					position++
				}
			l383:
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l386
					}
					position++
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if buffer[position] != rune('E') {
						goto l375
					}
					position++
				}
			l385:
				if !_rules[rulesp]() {
					goto l375
				}
				if !_rules[ruleRange]() {
					goto l375
				}
				if !_rules[rulesp]() {
					goto l375
				}
				if buffer[position] != rune(']') {
					goto l375
				}
				position++
				if !_rules[ruleAction20]() {
					goto l375
				}
				depth--
				add(ruleWindowedRelation, position376)
			}
			return true
		l375:
			position, tokenIndex, depth = position375, tokenIndex375, depth375
			return false
		},
		/* 26 DefWindowedRelation <- <(Relation (sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Range sp ']')? Action21)> */
		func() bool {
			position387, tokenIndex387, depth387 := position, tokenIndex, depth
			{
				position388 := position
				depth++
				if !_rules[ruleRelation]() {
					goto l387
				}
				{
					position389, tokenIndex389, depth389 := position, tokenIndex, depth
					if !_rules[rulesp]() {
						goto l389
					}
					if buffer[position] != rune('[') {
						goto l389
					}
					position++
					if !_rules[rulesp]() {
						goto l389
					}
					{
						position391, tokenIndex391, depth391 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l392
						}
						position++
						goto l391
					l392:
						position, tokenIndex, depth = position391, tokenIndex391, depth391
						if buffer[position] != rune('R') {
							goto l389
						}
						position++
					}
				l391:
					{
						position393, tokenIndex393, depth393 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l394
						}
						position++
						goto l393
					l394:
						position, tokenIndex, depth = position393, tokenIndex393, depth393
						if buffer[position] != rune('A') {
							goto l389
						}
						position++
					}
				l393:
					{
						position395, tokenIndex395, depth395 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l396
						}
						position++
						goto l395
					l396:
						position, tokenIndex, depth = position395, tokenIndex395, depth395
						if buffer[position] != rune('N') {
							goto l389
						}
						position++
					}
				l395:
					{
						position397, tokenIndex397, depth397 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l398
						}
						position++
						goto l397
					l398:
						position, tokenIndex, depth = position397, tokenIndex397, depth397
						if buffer[position] != rune('G') {
							goto l389
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
							goto l389
						}
						position++
					}
				l399:
					if !_rules[rulesp]() {
						goto l389
					}
					if !_rules[ruleRange]() {
						goto l389
					}
					if !_rules[rulesp]() {
						goto l389
					}
					if buffer[position] != rune(']') {
						goto l389
					}
					position++
					goto l390
				l389:
					position, tokenIndex, depth = position389, tokenIndex389, depth389
				}
			l390:
				if !_rules[ruleAction21]() {
					goto l387
				}
				depth--
				add(ruleDefWindowedRelation, position388)
			}
			return true
		l387:
			position, tokenIndex, depth = position387, tokenIndex387, depth387
			return false
		},
		/* 27 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action22)> */
		func() bool {
			position401, tokenIndex401, depth401 := position, tokenIndex, depth
			{
				position402 := position
				depth++
				{
					position403 := position
					depth++
					{
						position404, tokenIndex404, depth404 := position, tokenIndex, depth
						{
							position406, tokenIndex406, depth406 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l407
							}
							position++
							goto l406
						l407:
							position, tokenIndex, depth = position406, tokenIndex406, depth406
							if buffer[position] != rune('W') {
								goto l404
							}
							position++
						}
					l406:
						{
							position408, tokenIndex408, depth408 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l409
							}
							position++
							goto l408
						l409:
							position, tokenIndex, depth = position408, tokenIndex408, depth408
							if buffer[position] != rune('I') {
								goto l404
							}
							position++
						}
					l408:
						{
							position410, tokenIndex410, depth410 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l411
							}
							position++
							goto l410
						l411:
							position, tokenIndex, depth = position410, tokenIndex410, depth410
							if buffer[position] != rune('T') {
								goto l404
							}
							position++
						}
					l410:
						{
							position412, tokenIndex412, depth412 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l413
							}
							position++
							goto l412
						l413:
							position, tokenIndex, depth = position412, tokenIndex412, depth412
							if buffer[position] != rune('H') {
								goto l404
							}
							position++
						}
					l412:
						if !_rules[rulesp]() {
							goto l404
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l404
						}
						if !_rules[rulesp]() {
							goto l404
						}
					l414:
						{
							position415, tokenIndex415, depth415 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l415
							}
							position++
							if !_rules[rulesp]() {
								goto l415
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l415
							}
							goto l414
						l415:
							position, tokenIndex, depth = position415, tokenIndex415, depth415
						}
						goto l405
					l404:
						position, tokenIndex, depth = position404, tokenIndex404, depth404
					}
				l405:
					depth--
					add(rulePegText, position403)
				}
				if !_rules[ruleAction22]() {
					goto l401
				}
				depth--
				add(ruleSourceSinkSpecs, position402)
			}
			return true
		l401:
			position, tokenIndex, depth = position401, tokenIndex401, depth401
			return false
		},
		/* 28 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action23)> */
		func() bool {
			position416, tokenIndex416, depth416 := position, tokenIndex, depth
			{
				position417 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l416
				}
				if buffer[position] != rune('=') {
					goto l416
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l416
				}
				if !_rules[ruleAction23]() {
					goto l416
				}
				depth--
				add(ruleSourceSinkParam, position417)
			}
			return true
		l416:
			position, tokenIndex, depth = position416, tokenIndex416, depth416
			return false
		},
		/* 29 Expression <- <orExpr> */
		func() bool {
			position418, tokenIndex418, depth418 := position, tokenIndex, depth
			{
				position419 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l418
				}
				depth--
				add(ruleExpression, position419)
			}
			return true
		l418:
			position, tokenIndex, depth = position418, tokenIndex418, depth418
			return false
		},
		/* 30 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action24)> */
		func() bool {
			position420, tokenIndex420, depth420 := position, tokenIndex, depth
			{
				position421 := position
				depth++
				{
					position422 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l420
					}
					if !_rules[rulesp]() {
						goto l420
					}
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l423
						}
						if !_rules[rulesp]() {
							goto l423
						}
						if !_rules[ruleandExpr]() {
							goto l423
						}
						goto l424
					l423:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
					}
				l424:
					depth--
					add(rulePegText, position422)
				}
				if !_rules[ruleAction24]() {
					goto l420
				}
				depth--
				add(ruleorExpr, position421)
			}
			return true
		l420:
			position, tokenIndex, depth = position420, tokenIndex420, depth420
			return false
		},
		/* 31 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action25)> */
		func() bool {
			position425, tokenIndex425, depth425 := position, tokenIndex, depth
			{
				position426 := position
				depth++
				{
					position427 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l425
					}
					if !_rules[rulesp]() {
						goto l425
					}
					{
						position428, tokenIndex428, depth428 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l428
						}
						if !_rules[rulesp]() {
							goto l428
						}
						if !_rules[rulecomparisonExpr]() {
							goto l428
						}
						goto l429
					l428:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
					}
				l429:
					depth--
					add(rulePegText, position427)
				}
				if !_rules[ruleAction25]() {
					goto l425
				}
				depth--
				add(ruleandExpr, position426)
			}
			return true
		l425:
			position, tokenIndex, depth = position425, tokenIndex425, depth425
			return false
		},
		/* 32 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action26)> */
		func() bool {
			position430, tokenIndex430, depth430 := position, tokenIndex, depth
			{
				position431 := position
				depth++
				{
					position432 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l430
					}
					if !_rules[rulesp]() {
						goto l430
					}
					{
						position433, tokenIndex433, depth433 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l433
						}
						if !_rules[rulesp]() {
							goto l433
						}
						if !_rules[ruletermExpr]() {
							goto l433
						}
						goto l434
					l433:
						position, tokenIndex, depth = position433, tokenIndex433, depth433
					}
				l434:
					depth--
					add(rulePegText, position432)
				}
				if !_rules[ruleAction26]() {
					goto l430
				}
				depth--
				add(rulecomparisonExpr, position431)
			}
			return true
		l430:
			position, tokenIndex, depth = position430, tokenIndex430, depth430
			return false
		},
		/* 33 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action27)> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				{
					position437 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l435
					}
					if !_rules[rulesp]() {
						goto l435
					}
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l438
						}
						if !_rules[rulesp]() {
							goto l438
						}
						if !_rules[ruleproductExpr]() {
							goto l438
						}
						goto l439
					l438:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
					}
				l439:
					depth--
					add(rulePegText, position437)
				}
				if !_rules[ruleAction27]() {
					goto l435
				}
				depth--
				add(ruletermExpr, position436)
			}
			return true
		l435:
			position, tokenIndex, depth = position435, tokenIndex435, depth435
			return false
		},
		/* 34 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action28)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				{
					position442 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l440
					}
					if !_rules[rulesp]() {
						goto l440
					}
					{
						position443, tokenIndex443, depth443 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l443
						}
						if !_rules[rulesp]() {
							goto l443
						}
						if !_rules[rulebaseExpr]() {
							goto l443
						}
						goto l444
					l443:
						position, tokenIndex, depth = position443, tokenIndex443, depth443
					}
				l444:
					depth--
					add(rulePegText, position442)
				}
				if !_rules[ruleAction28]() {
					goto l440
				}
				depth--
				add(ruleproductExpr, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 35 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / RowValue / Literal)> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				{
					position447, tokenIndex447, depth447 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l448
					}
					position++
					if !_rules[rulesp]() {
						goto l448
					}
					if !_rules[ruleExpression]() {
						goto l448
					}
					if !_rules[rulesp]() {
						goto l448
					}
					if buffer[position] != rune(')') {
						goto l448
					}
					position++
					goto l447
				l448:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if !_rules[ruleBooleanLiteral]() {
						goto l449
					}
					goto l447
				l449:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if !_rules[ruleFuncApp]() {
						goto l450
					}
					goto l447
				l450:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if !_rules[ruleRowValue]() {
						goto l451
					}
					goto l447
				l451:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
					if !_rules[ruleLiteral]() {
						goto l445
					}
				}
			l447:
				depth--
				add(rulebaseExpr, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 36 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action29)> */
		func() bool {
			position452, tokenIndex452, depth452 := position, tokenIndex, depth
			{
				position453 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l452
				}
				if !_rules[rulesp]() {
					goto l452
				}
				if buffer[position] != rune('(') {
					goto l452
				}
				position++
				if !_rules[rulesp]() {
					goto l452
				}
				if !_rules[ruleFuncParams]() {
					goto l452
				}
				if !_rules[rulesp]() {
					goto l452
				}
				if buffer[position] != rune(')') {
					goto l452
				}
				position++
				if !_rules[ruleAction29]() {
					goto l452
				}
				depth--
				add(ruleFuncApp, position453)
			}
			return true
		l452:
			position, tokenIndex, depth = position452, tokenIndex452, depth452
			return false
		},
		/* 37 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action30)> */
		func() bool {
			position454, tokenIndex454, depth454 := position, tokenIndex, depth
			{
				position455 := position
				depth++
				{
					position456 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l454
					}
					if !_rules[rulesp]() {
						goto l454
					}
				l457:
					{
						position458, tokenIndex458, depth458 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l458
						}
						position++
						if !_rules[rulesp]() {
							goto l458
						}
						if !_rules[ruleExpression]() {
							goto l458
						}
						goto l457
					l458:
						position, tokenIndex, depth = position458, tokenIndex458, depth458
					}
					depth--
					add(rulePegText, position456)
				}
				if !_rules[ruleAction30]() {
					goto l454
				}
				depth--
				add(ruleFuncParams, position455)
			}
			return true
		l454:
			position, tokenIndex, depth = position454, tokenIndex454, depth454
			return false
		},
		/* 38 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position459, tokenIndex459, depth459 := position, tokenIndex, depth
			{
				position460 := position
				depth++
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l462
					}
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if !_rules[ruleNumericLiteral]() {
						goto l463
					}
					goto l461
				l463:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if !_rules[ruleStringLiteral]() {
						goto l459
					}
				}
			l461:
				depth--
				add(ruleLiteral, position460)
			}
			return true
		l459:
			position, tokenIndex, depth = position459, tokenIndex459, depth459
			return false
		},
		/* 39 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position464, tokenIndex464, depth464 := position, tokenIndex, depth
			{
				position465 := position
				depth++
				{
					position466, tokenIndex466, depth466 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l467
					}
					goto l466
				l467:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if !_rules[ruleNotEqual]() {
						goto l468
					}
					goto l466
				l468:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if !_rules[ruleLessOrEqual]() {
						goto l469
					}
					goto l466
				l469:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if !_rules[ruleLess]() {
						goto l470
					}
					goto l466
				l470:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if !_rules[ruleGreaterOrEqual]() {
						goto l471
					}
					goto l466
				l471:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if !_rules[ruleGreater]() {
						goto l472
					}
					goto l466
				l472:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if !_rules[ruleNotEqual]() {
						goto l464
					}
				}
			l466:
				depth--
				add(ruleComparisonOp, position465)
			}
			return true
		l464:
			position, tokenIndex, depth = position464, tokenIndex464, depth464
			return false
		},
		/* 40 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position473, tokenIndex473, depth473 := position, tokenIndex, depth
			{
				position474 := position
				depth++
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l476
					}
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !_rules[ruleMinus]() {
						goto l473
					}
				}
			l475:
				depth--
				add(rulePlusMinusOp, position474)
			}
			return true
		l473:
			position, tokenIndex, depth = position473, tokenIndex473, depth473
			return false
		},
		/* 41 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position477, tokenIndex477, depth477 := position, tokenIndex, depth
			{
				position478 := position
				depth++
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l480
					}
					goto l479
				l480:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if !_rules[ruleDivide]() {
						goto l481
					}
					goto l479
				l481:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
					if !_rules[ruleModulo]() {
						goto l477
					}
				}
			l479:
				depth--
				add(ruleMultDivOp, position478)
			}
			return true
		l477:
			position, tokenIndex, depth = position477, tokenIndex477, depth477
			return false
		},
		/* 42 Relation <- <(<ident> Action31)> */
		func() bool {
			position482, tokenIndex482, depth482 := position, tokenIndex, depth
			{
				position483 := position
				depth++
				{
					position484 := position
					depth++
					if !_rules[ruleident]() {
						goto l482
					}
					depth--
					add(rulePegText, position484)
				}
				if !_rules[ruleAction31]() {
					goto l482
				}
				depth--
				add(ruleRelation, position483)
			}
			return true
		l482:
			position, tokenIndex, depth = position482, tokenIndex482, depth482
			return false
		},
		/* 43 RowValue <- <(<((ident '.')? ident)> Action32)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				{
					position487 := position
					depth++
					{
						position488, tokenIndex488, depth488 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l488
						}
						if buffer[position] != rune('.') {
							goto l488
						}
						position++
						goto l489
					l488:
						position, tokenIndex, depth = position488, tokenIndex488, depth488
					}
				l489:
					if !_rules[ruleident]() {
						goto l485
					}
					depth--
					add(rulePegText, position487)
				}
				if !_rules[ruleAction32]() {
					goto l485
				}
				depth--
				add(ruleRowValue, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 44 NumericLiteral <- <(<('-'? [0-9]+)> Action33)> */
		func() bool {
			position490, tokenIndex490, depth490 := position, tokenIndex, depth
			{
				position491 := position
				depth++
				{
					position492 := position
					depth++
					{
						position493, tokenIndex493, depth493 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l493
						}
						position++
						goto l494
					l493:
						position, tokenIndex, depth = position493, tokenIndex493, depth493
					}
				l494:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l490
					}
					position++
				l495:
					{
						position496, tokenIndex496, depth496 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l496
						}
						position++
						goto l495
					l496:
						position, tokenIndex, depth = position496, tokenIndex496, depth496
					}
					depth--
					add(rulePegText, position492)
				}
				if !_rules[ruleAction33]() {
					goto l490
				}
				depth--
				add(ruleNumericLiteral, position491)
			}
			return true
		l490:
			position, tokenIndex, depth = position490, tokenIndex490, depth490
			return false
		},
		/* 45 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action34)> */
		func() bool {
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				{
					position499 := position
					depth++
					{
						position500, tokenIndex500, depth500 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l500
						}
						position++
						goto l501
					l500:
						position, tokenIndex, depth = position500, tokenIndex500, depth500
					}
				l501:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l497
					}
					position++
				l502:
					{
						position503, tokenIndex503, depth503 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l503
						}
						position++
						goto l502
					l503:
						position, tokenIndex, depth = position503, tokenIndex503, depth503
					}
					if buffer[position] != rune('.') {
						goto l497
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l497
					}
					position++
				l504:
					{
						position505, tokenIndex505, depth505 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l505
						}
						position++
						goto l504
					l505:
						position, tokenIndex, depth = position505, tokenIndex505, depth505
					}
					depth--
					add(rulePegText, position499)
				}
				if !_rules[ruleAction34]() {
					goto l497
				}
				depth--
				add(ruleFloatLiteral, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 46 Function <- <(<ident> Action35)> */
		func() bool {
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				{
					position508 := position
					depth++
					if !_rules[ruleident]() {
						goto l506
					}
					depth--
					add(rulePegText, position508)
				}
				if !_rules[ruleAction35]() {
					goto l506
				}
				depth--
				add(ruleFunction, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 47 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				{
					position511, tokenIndex511, depth511 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l512
					}
					goto l511
				l512:
					position, tokenIndex, depth = position511, tokenIndex511, depth511
					if !_rules[ruleFALSE]() {
						goto l509
					}
				}
			l511:
				depth--
				add(ruleBooleanLiteral, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 48 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action36)> */
		func() bool {
			position513, tokenIndex513, depth513 := position, tokenIndex, depth
			{
				position514 := position
				depth++
				{
					position515 := position
					depth++
					{
						position516, tokenIndex516, depth516 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l517
						}
						position++
						goto l516
					l517:
						position, tokenIndex, depth = position516, tokenIndex516, depth516
						if buffer[position] != rune('T') {
							goto l513
						}
						position++
					}
				l516:
					{
						position518, tokenIndex518, depth518 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l519
						}
						position++
						goto l518
					l519:
						position, tokenIndex, depth = position518, tokenIndex518, depth518
						if buffer[position] != rune('R') {
							goto l513
						}
						position++
					}
				l518:
					{
						position520, tokenIndex520, depth520 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l521
						}
						position++
						goto l520
					l521:
						position, tokenIndex, depth = position520, tokenIndex520, depth520
						if buffer[position] != rune('U') {
							goto l513
						}
						position++
					}
				l520:
					{
						position522, tokenIndex522, depth522 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l523
						}
						position++
						goto l522
					l523:
						position, tokenIndex, depth = position522, tokenIndex522, depth522
						if buffer[position] != rune('E') {
							goto l513
						}
						position++
					}
				l522:
					depth--
					add(rulePegText, position515)
				}
				if !_rules[ruleAction36]() {
					goto l513
				}
				depth--
				add(ruleTRUE, position514)
			}
			return true
		l513:
			position, tokenIndex, depth = position513, tokenIndex513, depth513
			return false
		},
		/* 49 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action37)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				{
					position526 := position
					depth++
					{
						position527, tokenIndex527, depth527 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l528
						}
						position++
						goto l527
					l528:
						position, tokenIndex, depth = position527, tokenIndex527, depth527
						if buffer[position] != rune('F') {
							goto l524
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
							goto l524
						}
						position++
					}
				l529:
					{
						position531, tokenIndex531, depth531 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l532
						}
						position++
						goto l531
					l532:
						position, tokenIndex, depth = position531, tokenIndex531, depth531
						if buffer[position] != rune('L') {
							goto l524
						}
						position++
					}
				l531:
					{
						position533, tokenIndex533, depth533 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l534
						}
						position++
						goto l533
					l534:
						position, tokenIndex, depth = position533, tokenIndex533, depth533
						if buffer[position] != rune('S') {
							goto l524
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
							goto l524
						}
						position++
					}
				l535:
					depth--
					add(rulePegText, position526)
				}
				if !_rules[ruleAction37]() {
					goto l524
				}
				depth--
				add(ruleFALSE, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 50 Wildcard <- <(<'*'> Action38)> */
		func() bool {
			position537, tokenIndex537, depth537 := position, tokenIndex, depth
			{
				position538 := position
				depth++
				{
					position539 := position
					depth++
					if buffer[position] != rune('*') {
						goto l537
					}
					position++
					depth--
					add(rulePegText, position539)
				}
				if !_rules[ruleAction38]() {
					goto l537
				}
				depth--
				add(ruleWildcard, position538)
			}
			return true
		l537:
			position, tokenIndex, depth = position537, tokenIndex537, depth537
			return false
		},
		/* 51 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action39)> */
		func() bool {
			position540, tokenIndex540, depth540 := position, tokenIndex, depth
			{
				position541 := position
				depth++
				{
					position542 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l540
					}
					position++
				l543:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						{
							position545, tokenIndex545, depth545 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l546
							}
							position++
							if buffer[position] != rune('\'') {
								goto l546
							}
							position++
							goto l545
						l546:
							position, tokenIndex, depth = position545, tokenIndex545, depth545
							{
								position547, tokenIndex547, depth547 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l547
								}
								position++
								goto l544
							l547:
								position, tokenIndex, depth = position547, tokenIndex547, depth547
							}
							if !matchDot() {
								goto l544
							}
						}
					l545:
						goto l543
					l544:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
					}
					if buffer[position] != rune('\'') {
						goto l540
					}
					position++
					depth--
					add(rulePegText, position542)
				}
				if !_rules[ruleAction39]() {
					goto l540
				}
				depth--
				add(ruleStringLiteral, position541)
			}
			return true
		l540:
			position, tokenIndex, depth = position540, tokenIndex540, depth540
			return false
		},
		/* 52 Emitter <- <(ISTREAM / DSTREAM / RSTREAM)> */
		func() bool {
			position548, tokenIndex548, depth548 := position, tokenIndex, depth
			{
				position549 := position
				depth++
				{
					position550, tokenIndex550, depth550 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l551
					}
					goto l550
				l551:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleDSTREAM]() {
						goto l552
					}
					goto l550
				l552:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if !_rules[ruleRSTREAM]() {
						goto l548
					}
				}
			l550:
				depth--
				add(ruleEmitter, position549)
			}
			return true
		l548:
			position, tokenIndex, depth = position548, tokenIndex548, depth548
			return false
		},
		/* 53 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action40)> */
		func() bool {
			position553, tokenIndex553, depth553 := position, tokenIndex, depth
			{
				position554 := position
				depth++
				{
					position555 := position
					depth++
					{
						position556, tokenIndex556, depth556 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l557
						}
						position++
						goto l556
					l557:
						position, tokenIndex, depth = position556, tokenIndex556, depth556
						if buffer[position] != rune('I') {
							goto l553
						}
						position++
					}
				l556:
					{
						position558, tokenIndex558, depth558 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l559
						}
						position++
						goto l558
					l559:
						position, tokenIndex, depth = position558, tokenIndex558, depth558
						if buffer[position] != rune('S') {
							goto l553
						}
						position++
					}
				l558:
					{
						position560, tokenIndex560, depth560 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l561
						}
						position++
						goto l560
					l561:
						position, tokenIndex, depth = position560, tokenIndex560, depth560
						if buffer[position] != rune('T') {
							goto l553
						}
						position++
					}
				l560:
					{
						position562, tokenIndex562, depth562 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l563
						}
						position++
						goto l562
					l563:
						position, tokenIndex, depth = position562, tokenIndex562, depth562
						if buffer[position] != rune('R') {
							goto l553
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
							goto l553
						}
						position++
					}
				l564:
					{
						position566, tokenIndex566, depth566 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l567
						}
						position++
						goto l566
					l567:
						position, tokenIndex, depth = position566, tokenIndex566, depth566
						if buffer[position] != rune('A') {
							goto l553
						}
						position++
					}
				l566:
					{
						position568, tokenIndex568, depth568 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l569
						}
						position++
						goto l568
					l569:
						position, tokenIndex, depth = position568, tokenIndex568, depth568
						if buffer[position] != rune('M') {
							goto l553
						}
						position++
					}
				l568:
					depth--
					add(rulePegText, position555)
				}
				if !_rules[ruleAction40]() {
					goto l553
				}
				depth--
				add(ruleISTREAM, position554)
			}
			return true
		l553:
			position, tokenIndex, depth = position553, tokenIndex553, depth553
			return false
		},
		/* 54 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action41)> */
		func() bool {
			position570, tokenIndex570, depth570 := position, tokenIndex, depth
			{
				position571 := position
				depth++
				{
					position572 := position
					depth++
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l574
						}
						position++
						goto l573
					l574:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
						if buffer[position] != rune('D') {
							goto l570
						}
						position++
					}
				l573:
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l576
						}
						position++
						goto l575
					l576:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
						if buffer[position] != rune('S') {
							goto l570
						}
						position++
					}
				l575:
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('T') {
							goto l570
						}
						position++
					}
				l577:
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l580
						}
						position++
						goto l579
					l580:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
						if buffer[position] != rune('R') {
							goto l570
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
							goto l570
						}
						position++
					}
				l581:
					{
						position583, tokenIndex583, depth583 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l584
						}
						position++
						goto l583
					l584:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
						if buffer[position] != rune('A') {
							goto l570
						}
						position++
					}
				l583:
					{
						position585, tokenIndex585, depth585 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l586
						}
						position++
						goto l585
					l586:
						position, tokenIndex, depth = position585, tokenIndex585, depth585
						if buffer[position] != rune('M') {
							goto l570
						}
						position++
					}
				l585:
					depth--
					add(rulePegText, position572)
				}
				if !_rules[ruleAction41]() {
					goto l570
				}
				depth--
				add(ruleDSTREAM, position571)
			}
			return true
		l570:
			position, tokenIndex, depth = position570, tokenIndex570, depth570
			return false
		},
		/* 55 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action42)> */
		func() bool {
			position587, tokenIndex587, depth587 := position, tokenIndex, depth
			{
				position588 := position
				depth++
				{
					position589 := position
					depth++
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
							goto l587
						}
						position++
					}
				l590:
					{
						position592, tokenIndex592, depth592 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex, depth = position592, tokenIndex592, depth592
						if buffer[position] != rune('S') {
							goto l587
						}
						position++
					}
				l592:
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
						if buffer[position] != rune('T') {
							goto l587
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
							goto l587
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
							goto l587
						}
						position++
					}
				l598:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l601
						}
						position++
						goto l600
					l601:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
						if buffer[position] != rune('A') {
							goto l587
						}
						position++
					}
				l600:
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l603
						}
						position++
						goto l602
					l603:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if buffer[position] != rune('M') {
							goto l587
						}
						position++
					}
				l602:
					depth--
					add(rulePegText, position589)
				}
				if !_rules[ruleAction42]() {
					goto l587
				}
				depth--
				add(ruleRSTREAM, position588)
			}
			return true
		l587:
			position, tokenIndex, depth = position587, tokenIndex587, depth587
			return false
		},
		/* 56 RangeUnit <- <(TUPLES / SECONDS)> */
		func() bool {
			position604, tokenIndex604, depth604 := position, tokenIndex, depth
			{
				position605 := position
				depth++
				{
					position606, tokenIndex606, depth606 := position, tokenIndex, depth
					if !_rules[ruleTUPLES]() {
						goto l607
					}
					goto l606
				l607:
					position, tokenIndex, depth = position606, tokenIndex606, depth606
					if !_rules[ruleSECONDS]() {
						goto l604
					}
				}
			l606:
				depth--
				add(ruleRangeUnit, position605)
			}
			return true
		l604:
			position, tokenIndex, depth = position604, tokenIndex604, depth604
			return false
		},
		/* 57 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action43)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('T') {
							goto l608
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('U') {
							goto l608
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('P') {
							goto l608
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('L') {
							goto l608
						}
						position++
					}
				l617:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('E') {
							goto l608
						}
						position++
					}
				l619:
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l622
						}
						position++
						goto l621
					l622:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
						if buffer[position] != rune('S') {
							goto l608
						}
						position++
					}
				l621:
					depth--
					add(rulePegText, position610)
				}
				if !_rules[ruleAction43]() {
					goto l608
				}
				depth--
				add(ruleTUPLES, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 58 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action44)> */
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
						if buffer[position] != rune('s') {
							goto l627
						}
						position++
						goto l626
					l627:
						position, tokenIndex, depth = position626, tokenIndex626, depth626
						if buffer[position] != rune('S') {
							goto l623
						}
						position++
					}
				l626:
					{
						position628, tokenIndex628, depth628 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position628, tokenIndex628, depth628
						if buffer[position] != rune('E') {
							goto l623
						}
						position++
					}
				l628:
					{
						position630, tokenIndex630, depth630 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l631
						}
						position++
						goto l630
					l631:
						position, tokenIndex, depth = position630, tokenIndex630, depth630
						if buffer[position] != rune('C') {
							goto l623
						}
						position++
					}
				l630:
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l633
						}
						position++
						goto l632
					l633:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
						if buffer[position] != rune('O') {
							goto l623
						}
						position++
					}
				l632:
					{
						position634, tokenIndex634, depth634 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l635
						}
						position++
						goto l634
					l635:
						position, tokenIndex, depth = position634, tokenIndex634, depth634
						if buffer[position] != rune('N') {
							goto l623
						}
						position++
					}
				l634:
					{
						position636, tokenIndex636, depth636 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l637
						}
						position++
						goto l636
					l637:
						position, tokenIndex, depth = position636, tokenIndex636, depth636
						if buffer[position] != rune('D') {
							goto l623
						}
						position++
					}
				l636:
					{
						position638, tokenIndex638, depth638 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l639
						}
						position++
						goto l638
					l639:
						position, tokenIndex, depth = position638, tokenIndex638, depth638
						if buffer[position] != rune('S') {
							goto l623
						}
						position++
					}
				l638:
					depth--
					add(rulePegText, position625)
				}
				if !_rules[ruleAction44]() {
					goto l623
				}
				depth--
				add(ruleSECONDS, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 59 StreamIdentifier <- <(<ident> Action45)> */
		func() bool {
			position640, tokenIndex640, depth640 := position, tokenIndex, depth
			{
				position641 := position
				depth++
				{
					position642 := position
					depth++
					if !_rules[ruleident]() {
						goto l640
					}
					depth--
					add(rulePegText, position642)
				}
				if !_rules[ruleAction45]() {
					goto l640
				}
				depth--
				add(ruleStreamIdentifier, position641)
			}
			return true
		l640:
			position, tokenIndex, depth = position640, tokenIndex640, depth640
			return false
		},
		/* 60 SourceSinkType <- <(<ident> Action46)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645 := position
					depth++
					if !_rules[ruleident]() {
						goto l643
					}
					depth--
					add(rulePegText, position645)
				}
				if !_rules[ruleAction46]() {
					goto l643
				}
				depth--
				add(ruleSourceSinkType, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 61 SourceSinkParamKey <- <(<ident> Action47)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				{
					position648 := position
					depth++
					if !_rules[ruleident]() {
						goto l646
					}
					depth--
					add(rulePegText, position648)
				}
				if !_rules[ruleAction47]() {
					goto l646
				}
				depth--
				add(ruleSourceSinkParamKey, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 62 SourceSinkParamVal <- <(<([a-z] / [A-Z] / [0-9] / '_')+> Action48)> */
		func() bool {
			position649, tokenIndex649, depth649 := position, tokenIndex, depth
			{
				position650 := position
				depth++
				{
					position651 := position
					depth++
					{
						position654, tokenIndex654, depth654 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l655
						}
						position++
						goto l654
					l655:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l656
						}
						position++
						goto l654
					l656:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l657
						}
						position++
						goto l654
					l657:
						position, tokenIndex, depth = position654, tokenIndex654, depth654
						if buffer[position] != rune('_') {
							goto l649
						}
						position++
					}
				l654:
				l652:
					{
						position653, tokenIndex653, depth653 := position, tokenIndex, depth
						{
							position658, tokenIndex658, depth658 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l659
							}
							position++
							goto l658
						l659:
							position, tokenIndex, depth = position658, tokenIndex658, depth658
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l660
							}
							position++
							goto l658
						l660:
							position, tokenIndex, depth = position658, tokenIndex658, depth658
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l661
							}
							position++
							goto l658
						l661:
							position, tokenIndex, depth = position658, tokenIndex658, depth658
							if buffer[position] != rune('_') {
								goto l653
							}
							position++
						}
					l658:
						goto l652
					l653:
						position, tokenIndex, depth = position653, tokenIndex653, depth653
					}
					depth--
					add(rulePegText, position651)
				}
				if !_rules[ruleAction48]() {
					goto l649
				}
				depth--
				add(ruleSourceSinkParamVal, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 63 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action49)> */
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
						if buffer[position] != rune('o') {
							goto l666
						}
						position++
						goto l665
					l666:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if buffer[position] != rune('O') {
							goto l662
						}
						position++
					}
				l665:
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if buffer[position] != rune('R') {
							goto l662
						}
						position++
					}
				l667:
					depth--
					add(rulePegText, position664)
				}
				if !_rules[ruleAction49]() {
					goto l662
				}
				depth--
				add(ruleOr, position663)
			}
			return true
		l662:
			position, tokenIndex, depth = position662, tokenIndex662, depth662
			return false
		},
		/* 64 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action50)> */
		func() bool {
			position669, tokenIndex669, depth669 := position, tokenIndex, depth
			{
				position670 := position
				depth++
				{
					position671 := position
					depth++
					{
						position672, tokenIndex672, depth672 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l673
						}
						position++
						goto l672
					l673:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						if buffer[position] != rune('A') {
							goto l669
						}
						position++
					}
				l672:
					{
						position674, tokenIndex674, depth674 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l675
						}
						position++
						goto l674
					l675:
						position, tokenIndex, depth = position674, tokenIndex674, depth674
						if buffer[position] != rune('N') {
							goto l669
						}
						position++
					}
				l674:
					{
						position676, tokenIndex676, depth676 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l677
						}
						position++
						goto l676
					l677:
						position, tokenIndex, depth = position676, tokenIndex676, depth676
						if buffer[position] != rune('D') {
							goto l669
						}
						position++
					}
				l676:
					depth--
					add(rulePegText, position671)
				}
				if !_rules[ruleAction50]() {
					goto l669
				}
				depth--
				add(ruleAnd, position670)
			}
			return true
		l669:
			position, tokenIndex, depth = position669, tokenIndex669, depth669
			return false
		},
		/* 65 Equal <- <(<'='> Action51)> */
		func() bool {
			position678, tokenIndex678, depth678 := position, tokenIndex, depth
			{
				position679 := position
				depth++
				{
					position680 := position
					depth++
					if buffer[position] != rune('=') {
						goto l678
					}
					position++
					depth--
					add(rulePegText, position680)
				}
				if !_rules[ruleAction51]() {
					goto l678
				}
				depth--
				add(ruleEqual, position679)
			}
			return true
		l678:
			position, tokenIndex, depth = position678, tokenIndex678, depth678
			return false
		},
		/* 66 Less <- <(<'<'> Action52)> */
		func() bool {
			position681, tokenIndex681, depth681 := position, tokenIndex, depth
			{
				position682 := position
				depth++
				{
					position683 := position
					depth++
					if buffer[position] != rune('<') {
						goto l681
					}
					position++
					depth--
					add(rulePegText, position683)
				}
				if !_rules[ruleAction52]() {
					goto l681
				}
				depth--
				add(ruleLess, position682)
			}
			return true
		l681:
			position, tokenIndex, depth = position681, tokenIndex681, depth681
			return false
		},
		/* 67 LessOrEqual <- <(<('<' '=')> Action53)> */
		func() bool {
			position684, tokenIndex684, depth684 := position, tokenIndex, depth
			{
				position685 := position
				depth++
				{
					position686 := position
					depth++
					if buffer[position] != rune('<') {
						goto l684
					}
					position++
					if buffer[position] != rune('=') {
						goto l684
					}
					position++
					depth--
					add(rulePegText, position686)
				}
				if !_rules[ruleAction53]() {
					goto l684
				}
				depth--
				add(ruleLessOrEqual, position685)
			}
			return true
		l684:
			position, tokenIndex, depth = position684, tokenIndex684, depth684
			return false
		},
		/* 68 Greater <- <(<'>'> Action54)> */
		func() bool {
			position687, tokenIndex687, depth687 := position, tokenIndex, depth
			{
				position688 := position
				depth++
				{
					position689 := position
					depth++
					if buffer[position] != rune('>') {
						goto l687
					}
					position++
					depth--
					add(rulePegText, position689)
				}
				if !_rules[ruleAction54]() {
					goto l687
				}
				depth--
				add(ruleGreater, position688)
			}
			return true
		l687:
			position, tokenIndex, depth = position687, tokenIndex687, depth687
			return false
		},
		/* 69 GreaterOrEqual <- <(<('>' '=')> Action55)> */
		func() bool {
			position690, tokenIndex690, depth690 := position, tokenIndex, depth
			{
				position691 := position
				depth++
				{
					position692 := position
					depth++
					if buffer[position] != rune('>') {
						goto l690
					}
					position++
					if buffer[position] != rune('=') {
						goto l690
					}
					position++
					depth--
					add(rulePegText, position692)
				}
				if !_rules[ruleAction55]() {
					goto l690
				}
				depth--
				add(ruleGreaterOrEqual, position691)
			}
			return true
		l690:
			position, tokenIndex, depth = position690, tokenIndex690, depth690
			return false
		},
		/* 70 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action56)> */
		func() bool {
			position693, tokenIndex693, depth693 := position, tokenIndex, depth
			{
				position694 := position
				depth++
				{
					position695 := position
					depth++
					{
						position696, tokenIndex696, depth696 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l697
						}
						position++
						if buffer[position] != rune('=') {
							goto l697
						}
						position++
						goto l696
					l697:
						position, tokenIndex, depth = position696, tokenIndex696, depth696
						if buffer[position] != rune('<') {
							goto l693
						}
						position++
						if buffer[position] != rune('>') {
							goto l693
						}
						position++
					}
				l696:
					depth--
					add(rulePegText, position695)
				}
				if !_rules[ruleAction56]() {
					goto l693
				}
				depth--
				add(ruleNotEqual, position694)
			}
			return true
		l693:
			position, tokenIndex, depth = position693, tokenIndex693, depth693
			return false
		},
		/* 71 Plus <- <(<'+'> Action57)> */
		func() bool {
			position698, tokenIndex698, depth698 := position, tokenIndex, depth
			{
				position699 := position
				depth++
				{
					position700 := position
					depth++
					if buffer[position] != rune('+') {
						goto l698
					}
					position++
					depth--
					add(rulePegText, position700)
				}
				if !_rules[ruleAction57]() {
					goto l698
				}
				depth--
				add(rulePlus, position699)
			}
			return true
		l698:
			position, tokenIndex, depth = position698, tokenIndex698, depth698
			return false
		},
		/* 72 Minus <- <(<'-'> Action58)> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				{
					position703 := position
					depth++
					if buffer[position] != rune('-') {
						goto l701
					}
					position++
					depth--
					add(rulePegText, position703)
				}
				if !_rules[ruleAction58]() {
					goto l701
				}
				depth--
				add(ruleMinus, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 73 Multiply <- <(<'*'> Action59)> */
		func() bool {
			position704, tokenIndex704, depth704 := position, tokenIndex, depth
			{
				position705 := position
				depth++
				{
					position706 := position
					depth++
					if buffer[position] != rune('*') {
						goto l704
					}
					position++
					depth--
					add(rulePegText, position706)
				}
				if !_rules[ruleAction59]() {
					goto l704
				}
				depth--
				add(ruleMultiply, position705)
			}
			return true
		l704:
			position, tokenIndex, depth = position704, tokenIndex704, depth704
			return false
		},
		/* 74 Divide <- <(<'/'> Action60)> */
		func() bool {
			position707, tokenIndex707, depth707 := position, tokenIndex, depth
			{
				position708 := position
				depth++
				{
					position709 := position
					depth++
					if buffer[position] != rune('/') {
						goto l707
					}
					position++
					depth--
					add(rulePegText, position709)
				}
				if !_rules[ruleAction60]() {
					goto l707
				}
				depth--
				add(ruleDivide, position708)
			}
			return true
		l707:
			position, tokenIndex, depth = position707, tokenIndex707, depth707
			return false
		},
		/* 75 Modulo <- <(<'%'> Action61)> */
		func() bool {
			position710, tokenIndex710, depth710 := position, tokenIndex, depth
			{
				position711 := position
				depth++
				{
					position712 := position
					depth++
					if buffer[position] != rune('%') {
						goto l710
					}
					position++
					depth--
					add(rulePegText, position712)
				}
				if !_rules[ruleAction61]() {
					goto l710
				}
				depth--
				add(ruleModulo, position711)
			}
			return true
		l710:
			position, tokenIndex, depth = position710, tokenIndex710, depth710
			return false
		},
		/* 76 Identifier <- <(<ident> Action62)> */
		func() bool {
			position713, tokenIndex713, depth713 := position, tokenIndex, depth
			{
				position714 := position
				depth++
				{
					position715 := position
					depth++
					if !_rules[ruleident]() {
						goto l713
					}
					depth--
					add(rulePegText, position715)
				}
				if !_rules[ruleAction62]() {
					goto l713
				}
				depth--
				add(ruleIdentifier, position714)
			}
			return true
		l713:
			position, tokenIndex, depth = position713, tokenIndex713, depth713
			return false
		},
		/* 77 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position716, tokenIndex716, depth716 := position, tokenIndex, depth
			{
				position717 := position
				depth++
				{
					position718, tokenIndex718, depth718 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l719
					}
					position++
					goto l718
				l719:
					position, tokenIndex, depth = position718, tokenIndex718, depth718
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l716
					}
					position++
				}
			l718:
			l720:
				{
					position721, tokenIndex721, depth721 := position, tokenIndex, depth
					{
						position722, tokenIndex722, depth722 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l723
						}
						position++
						goto l722
					l723:
						position, tokenIndex, depth = position722, tokenIndex722, depth722
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l724
						}
						position++
						goto l722
					l724:
						position, tokenIndex, depth = position722, tokenIndex722, depth722
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l725
						}
						position++
						goto l722
					l725:
						position, tokenIndex, depth = position722, tokenIndex722, depth722
						if buffer[position] != rune('_') {
							goto l721
						}
						position++
					}
				l722:
					goto l720
				l721:
					position, tokenIndex, depth = position721, tokenIndex721, depth721
				}
				depth--
				add(ruleident, position717)
			}
			return true
		l716:
			position, tokenIndex, depth = position716, tokenIndex716, depth716
			return false
		},
		/* 78 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position727 := position
				depth++
			l728:
				{
					position729, tokenIndex729, depth729 := position, tokenIndex, depth
					{
						position730, tokenIndex730, depth730 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l731
						}
						position++
						goto l730
					l731:
						position, tokenIndex, depth = position730, tokenIndex730, depth730
						if buffer[position] != rune('\t') {
							goto l732
						}
						position++
						goto l730
					l732:
						position, tokenIndex, depth = position730, tokenIndex730, depth730
						if buffer[position] != rune('\n') {
							goto l729
						}
						position++
					}
				l730:
					goto l728
				l729:
					position, tokenIndex, depth = position729, tokenIndex729, depth729
				}
				depth--
				add(rulesp, position727)
			}
			return true
		},
		/* 80 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 81 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 82 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 83 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 84 Action4 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 85 Action5 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 86 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 87 Action7 <- <{
		    p.AssembleEmitProjections()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 89 Action8 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 90 Action9 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 91 Action10 <- <{
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
		/* 92 Action11 <- <{
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 93 Action12 <- <{
		    p.AssembleRange()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 94 Action13 <- <{
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
		/* 95 Action14 <- <{
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
		/* 96 Action15 <- <{
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
		/* 97 Action16 <- <{
		    p.EnsureAliasWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 98 Action17 <- <{
		    p.EnsureAliasWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 99 Action18 <- <{
		    p.AssembleAliasWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 100 Action19 <- <{
		    p.AssembleAliasWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 101 Action20 <- <{
		    p.AssembleWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 102 Action21 <- <{
		    p.AssembleWindowedRelation()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 103 Action22 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 104 Action23 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 105 Action24 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 106 Action25 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 107 Action26 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 108 Action27 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 109 Action28 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 110 Action29 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 111 Action30 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 112 Action31 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRelation(substr))
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 113 Action32 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 114 Action33 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 115 Action34 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 116 Action35 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 117 Action36 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 118 Action37 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 119 Action38 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 120 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 121 Action40 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 122 Action41 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 123 Action42 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 124 Action43 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 125 Action44 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 126 Action45 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 127 Action46 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 128 Action47 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 129 Action48 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamVal(substr))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 130 Action49 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 131 Action50 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 132 Action51 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 133 Action52 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 134 Action53 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 135 Action54 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 136 Action55 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 137 Action56 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 138 Action57 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 139 Action58 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 140 Action59 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 141 Action60 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 142 Action61 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 143 Action62 <- <{
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

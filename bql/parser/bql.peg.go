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
	ruleWindowedFrom
	ruleRange
	ruleFrom
	ruleRelations
	ruleFilter
	ruleGrouping
	ruleGroupList
	ruleHaving
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
	ruleColumnName
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
	"WindowedFrom",
	"Range",
	"From",
	"Relations",
	"Filter",
	"Grouping",
	"GroupList",
	"Having",
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
	"ColumnName",
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
	rules  [127]func() bool
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

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction10:

			p.AssembleRange()

		case ruleAction11:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleFrom(begin, end)

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

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction16:

			p.AssembleSourceSinkParam()

		case ruleAction17:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction18:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction19:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction20:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction21:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction22:

			p.AssembleFuncApp()

		case ruleAction23:

			p.AssembleExpressions(begin, end)

		case ruleAction24:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRelation(substr))

		case ruleAction25:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewColumnName(substr))

		case ruleAction26:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction27:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction28:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction29:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction30:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction31:

			p.PushComponent(begin, end, NewWildcard())

		case ruleAction32:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction33:

			p.PushComponent(begin, end, Istream)

		case ruleAction34:

			p.PushComponent(begin, end, Dstream)

		case ruleAction35:

			p.PushComponent(begin, end, Rstream)

		case ruleAction36:

			p.PushComponent(begin, end, Tuples)

		case ruleAction37:

			p.PushComponent(begin, end, Seconds)

		case ruleAction38:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkName(substr))

		case ruleAction39:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction40:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction41:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamVal(substr))

		case ruleAction42:

			p.PushComponent(begin, end, Or)

		case ruleAction43:

			p.PushComponent(begin, end, And)

		case ruleAction44:

			p.PushComponent(begin, end, Equal)

		case ruleAction45:

			p.PushComponent(begin, end, Less)

		case ruleAction46:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction47:

			p.PushComponent(begin, end, Greater)

		case ruleAction48:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction49:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction50:

			p.PushComponent(begin, end, Plus)

		case ruleAction51:

			p.PushComponent(begin, end, Minus)

		case ruleAction52:

			p.PushComponent(begin, end, Multiply)

		case ruleAction53:

			p.PushComponent(begin, end, Divide)

		case ruleAction54:

			p.PushComponent(begin, end, Modulo)

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
		/* 1 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') sp Projections sp From sp Filter sp Grouping sp Having sp Action0)> */
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
				if !_rules[ruleFrom]() {
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
		/* 8 EmitProjections <- <(Emitter sp '(' Projections ')' Action7)> */
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
				if buffer[position] != rune('(') {
					goto l244
				}
				position++
				if !_rules[ruleProjections]() {
					goto l244
				}
				if buffer[position] != rune(')') {
					goto l244
				}
				position++
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
		/* 10 Projection <- <Expression> */
		func() bool {
			position251, tokenIndex251, depth251 := position, tokenIndex, depth
			{
				position252 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l251
				}
				depth--
				add(ruleProjection, position252)
			}
			return true
		l251:
			position, tokenIndex, depth = position251, tokenIndex251, depth251
			return false
		},
		/* 11 WindowedFrom <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations sp '[' sp (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Range sp ']')?> Action9)> */
		func() bool {
			position253, tokenIndex253, depth253 := position, tokenIndex, depth
			{
				position254 := position
				depth++
				{
					position255 := position
					depth++
					{
						position256, tokenIndex256, depth256 := position, tokenIndex, depth
						{
							position258, tokenIndex258, depth258 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l259
							}
							position++
							goto l258
						l259:
							position, tokenIndex, depth = position258, tokenIndex258, depth258
							if buffer[position] != rune('F') {
								goto l256
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
								goto l256
							}
							position++
						}
					l260:
						{
							position262, tokenIndex262, depth262 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l263
							}
							position++
							goto l262
						l263:
							position, tokenIndex, depth = position262, tokenIndex262, depth262
							if buffer[position] != rune('O') {
								goto l256
							}
							position++
						}
					l262:
						{
							position264, tokenIndex264, depth264 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l265
							}
							position++
							goto l264
						l265:
							position, tokenIndex, depth = position264, tokenIndex264, depth264
							if buffer[position] != rune('M') {
								goto l256
							}
							position++
						}
					l264:
						if !_rules[rulesp]() {
							goto l256
						}
						if !_rules[ruleRelations]() {
							goto l256
						}
						if !_rules[rulesp]() {
							goto l256
						}
						if buffer[position] != rune('[') {
							goto l256
						}
						position++
						if !_rules[rulesp]() {
							goto l256
						}
						{
							position266, tokenIndex266, depth266 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l267
							}
							position++
							goto l266
						l267:
							position, tokenIndex, depth = position266, tokenIndex266, depth266
							if buffer[position] != rune('R') {
								goto l256
							}
							position++
						}
					l266:
						{
							position268, tokenIndex268, depth268 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l269
							}
							position++
							goto l268
						l269:
							position, tokenIndex, depth = position268, tokenIndex268, depth268
							if buffer[position] != rune('A') {
								goto l256
							}
							position++
						}
					l268:
						{
							position270, tokenIndex270, depth270 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l271
							}
							position++
							goto l270
						l271:
							position, tokenIndex, depth = position270, tokenIndex270, depth270
							if buffer[position] != rune('N') {
								goto l256
							}
							position++
						}
					l270:
						{
							position272, tokenIndex272, depth272 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l273
							}
							position++
							goto l272
						l273:
							position, tokenIndex, depth = position272, tokenIndex272, depth272
							if buffer[position] != rune('G') {
								goto l256
							}
							position++
						}
					l272:
						{
							position274, tokenIndex274, depth274 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l275
							}
							position++
							goto l274
						l275:
							position, tokenIndex, depth = position274, tokenIndex274, depth274
							if buffer[position] != rune('E') {
								goto l256
							}
							position++
						}
					l274:
						if !_rules[rulesp]() {
							goto l256
						}
						if !_rules[ruleRange]() {
							goto l256
						}
						if !_rules[rulesp]() {
							goto l256
						}
						if buffer[position] != rune(']') {
							goto l256
						}
						position++
						goto l257
					l256:
						position, tokenIndex, depth = position256, tokenIndex256, depth256
					}
				l257:
					depth--
					add(rulePegText, position255)
				}
				if !_rules[ruleAction9]() {
					goto l253
				}
				depth--
				add(ruleWindowedFrom, position254)
			}
			return true
		l253:
			position, tokenIndex, depth = position253, tokenIndex253, depth253
			return false
		},
		/* 12 Range <- <(NumericLiteral sp RangeUnit Action10)> */
		func() bool {
			position276, tokenIndex276, depth276 := position, tokenIndex, depth
			{
				position277 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l276
				}
				if !_rules[rulesp]() {
					goto l276
				}
				if !_rules[ruleRangeUnit]() {
					goto l276
				}
				if !_rules[ruleAction10]() {
					goto l276
				}
				depth--
				add(ruleRange, position277)
			}
			return true
		l276:
			position, tokenIndex, depth = position276, tokenIndex276, depth276
			return false
		},
		/* 13 From <- <(<(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') sp Relations)?> Action11)> */
		func() bool {
			position278, tokenIndex278, depth278 := position, tokenIndex, depth
			{
				position279 := position
				depth++
				{
					position280 := position
					depth++
					{
						position281, tokenIndex281, depth281 := position, tokenIndex, depth
						{
							position283, tokenIndex283, depth283 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l284
							}
							position++
							goto l283
						l284:
							position, tokenIndex, depth = position283, tokenIndex283, depth283
							if buffer[position] != rune('F') {
								goto l281
							}
							position++
						}
					l283:
						{
							position285, tokenIndex285, depth285 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l286
							}
							position++
							goto l285
						l286:
							position, tokenIndex, depth = position285, tokenIndex285, depth285
							if buffer[position] != rune('R') {
								goto l281
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
								goto l281
							}
							position++
						}
					l287:
						{
							position289, tokenIndex289, depth289 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l290
							}
							position++
							goto l289
						l290:
							position, tokenIndex, depth = position289, tokenIndex289, depth289
							if buffer[position] != rune('M') {
								goto l281
							}
							position++
						}
					l289:
						if !_rules[rulesp]() {
							goto l281
						}
						if !_rules[ruleRelations]() {
							goto l281
						}
						goto l282
					l281:
						position, tokenIndex, depth = position281, tokenIndex281, depth281
					}
				l282:
					depth--
					add(rulePegText, position280)
				}
				if !_rules[ruleAction11]() {
					goto l278
				}
				depth--
				add(ruleFrom, position279)
			}
			return true
		l278:
			position, tokenIndex, depth = position278, tokenIndex278, depth278
			return false
		},
		/* 14 Relations <- <(Relation sp (',' sp Relation)*)> */
		func() bool {
			position291, tokenIndex291, depth291 := position, tokenIndex, depth
			{
				position292 := position
				depth++
				if !_rules[ruleRelation]() {
					goto l291
				}
				if !_rules[rulesp]() {
					goto l291
				}
			l293:
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l294
					}
					position++
					if !_rules[rulesp]() {
						goto l294
					}
					if !_rules[ruleRelation]() {
						goto l294
					}
					goto l293
				l294:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
				}
				depth--
				add(ruleRelations, position292)
			}
			return true
		l291:
			position, tokenIndex, depth = position291, tokenIndex291, depth291
			return false
		},
		/* 15 Filter <- <(<(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') sp Expression)?> Action12)> */
		func() bool {
			position295, tokenIndex295, depth295 := position, tokenIndex, depth
			{
				position296 := position
				depth++
				{
					position297 := position
					depth++
					{
						position298, tokenIndex298, depth298 := position, tokenIndex, depth
						{
							position300, tokenIndex300, depth300 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l301
							}
							position++
							goto l300
						l301:
							position, tokenIndex, depth = position300, tokenIndex300, depth300
							if buffer[position] != rune('W') {
								goto l298
							}
							position++
						}
					l300:
						{
							position302, tokenIndex302, depth302 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l303
							}
							position++
							goto l302
						l303:
							position, tokenIndex, depth = position302, tokenIndex302, depth302
							if buffer[position] != rune('H') {
								goto l298
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
								goto l298
							}
							position++
						}
					l304:
						{
							position306, tokenIndex306, depth306 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							goto l306
						l307:
							position, tokenIndex, depth = position306, tokenIndex306, depth306
							if buffer[position] != rune('R') {
								goto l298
							}
							position++
						}
					l306:
						{
							position308, tokenIndex308, depth308 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l309
							}
							position++
							goto l308
						l309:
							position, tokenIndex, depth = position308, tokenIndex308, depth308
							if buffer[position] != rune('E') {
								goto l298
							}
							position++
						}
					l308:
						if !_rules[rulesp]() {
							goto l298
						}
						if !_rules[ruleExpression]() {
							goto l298
						}
						goto l299
					l298:
						position, tokenIndex, depth = position298, tokenIndex298, depth298
					}
				l299:
					depth--
					add(rulePegText, position297)
				}
				if !_rules[ruleAction12]() {
					goto l295
				}
				depth--
				add(ruleFilter, position296)
			}
			return true
		l295:
			position, tokenIndex, depth = position295, tokenIndex295, depth295
			return false
		},
		/* 16 Grouping <- <(<(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action13)> */
		func() bool {
			position310, tokenIndex310, depth310 := position, tokenIndex, depth
			{
				position311 := position
				depth++
				{
					position312 := position
					depth++
					{
						position313, tokenIndex313, depth313 := position, tokenIndex, depth
						{
							position315, tokenIndex315, depth315 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l316
							}
							position++
							goto l315
						l316:
							position, tokenIndex, depth = position315, tokenIndex315, depth315
							if buffer[position] != rune('G') {
								goto l313
							}
							position++
						}
					l315:
						{
							position317, tokenIndex317, depth317 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l318
							}
							position++
							goto l317
						l318:
							position, tokenIndex, depth = position317, tokenIndex317, depth317
							if buffer[position] != rune('R') {
								goto l313
							}
							position++
						}
					l317:
						{
							position319, tokenIndex319, depth319 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l320
							}
							position++
							goto l319
						l320:
							position, tokenIndex, depth = position319, tokenIndex319, depth319
							if buffer[position] != rune('O') {
								goto l313
							}
							position++
						}
					l319:
						{
							position321, tokenIndex321, depth321 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l322
							}
							position++
							goto l321
						l322:
							position, tokenIndex, depth = position321, tokenIndex321, depth321
							if buffer[position] != rune('U') {
								goto l313
							}
							position++
						}
					l321:
						{
							position323, tokenIndex323, depth323 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l324
							}
							position++
							goto l323
						l324:
							position, tokenIndex, depth = position323, tokenIndex323, depth323
							if buffer[position] != rune('P') {
								goto l313
							}
							position++
						}
					l323:
						if !_rules[rulesp]() {
							goto l313
						}
						{
							position325, tokenIndex325, depth325 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l326
							}
							position++
							goto l325
						l326:
							position, tokenIndex, depth = position325, tokenIndex325, depth325
							if buffer[position] != rune('B') {
								goto l313
							}
							position++
						}
					l325:
						{
							position327, tokenIndex327, depth327 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l328
							}
							position++
							goto l327
						l328:
							position, tokenIndex, depth = position327, tokenIndex327, depth327
							if buffer[position] != rune('Y') {
								goto l313
							}
							position++
						}
					l327:
						if !_rules[rulesp]() {
							goto l313
						}
						if !_rules[ruleGroupList]() {
							goto l313
						}
						goto l314
					l313:
						position, tokenIndex, depth = position313, tokenIndex313, depth313
					}
				l314:
					depth--
					add(rulePegText, position312)
				}
				if !_rules[ruleAction13]() {
					goto l310
				}
				depth--
				add(ruleGrouping, position311)
			}
			return true
		l310:
			position, tokenIndex, depth = position310, tokenIndex310, depth310
			return false
		},
		/* 17 GroupList <- <(Expression sp (',' sp Expression)*)> */
		func() bool {
			position329, tokenIndex329, depth329 := position, tokenIndex, depth
			{
				position330 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l329
				}
				if !_rules[rulesp]() {
					goto l329
				}
			l331:
				{
					position332, tokenIndex332, depth332 := position, tokenIndex, depth
					if buffer[position] != rune(',') {
						goto l332
					}
					position++
					if !_rules[rulesp]() {
						goto l332
					}
					if !_rules[ruleExpression]() {
						goto l332
					}
					goto l331
				l332:
					position, tokenIndex, depth = position332, tokenIndex332, depth332
				}
				depth--
				add(ruleGroupList, position330)
			}
			return true
		l329:
			position, tokenIndex, depth = position329, tokenIndex329, depth329
			return false
		},
		/* 18 Having <- <(<(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') sp Expression)?> Action14)> */
		func() bool {
			position333, tokenIndex333, depth333 := position, tokenIndex, depth
			{
				position334 := position
				depth++
				{
					position335 := position
					depth++
					{
						position336, tokenIndex336, depth336 := position, tokenIndex, depth
						{
							position338, tokenIndex338, depth338 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l339
							}
							position++
							goto l338
						l339:
							position, tokenIndex, depth = position338, tokenIndex338, depth338
							if buffer[position] != rune('H') {
								goto l336
							}
							position++
						}
					l338:
						{
							position340, tokenIndex340, depth340 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l341
							}
							position++
							goto l340
						l341:
							position, tokenIndex, depth = position340, tokenIndex340, depth340
							if buffer[position] != rune('A') {
								goto l336
							}
							position++
						}
					l340:
						{
							position342, tokenIndex342, depth342 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l343
							}
							position++
							goto l342
						l343:
							position, tokenIndex, depth = position342, tokenIndex342, depth342
							if buffer[position] != rune('V') {
								goto l336
							}
							position++
						}
					l342:
						{
							position344, tokenIndex344, depth344 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l345
							}
							position++
							goto l344
						l345:
							position, tokenIndex, depth = position344, tokenIndex344, depth344
							if buffer[position] != rune('I') {
								goto l336
							}
							position++
						}
					l344:
						{
							position346, tokenIndex346, depth346 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l347
							}
							position++
							goto l346
						l347:
							position, tokenIndex, depth = position346, tokenIndex346, depth346
							if buffer[position] != rune('N') {
								goto l336
							}
							position++
						}
					l346:
						{
							position348, tokenIndex348, depth348 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l349
							}
							position++
							goto l348
						l349:
							position, tokenIndex, depth = position348, tokenIndex348, depth348
							if buffer[position] != rune('G') {
								goto l336
							}
							position++
						}
					l348:
						if !_rules[rulesp]() {
							goto l336
						}
						if !_rules[ruleExpression]() {
							goto l336
						}
						goto l337
					l336:
						position, tokenIndex, depth = position336, tokenIndex336, depth336
					}
				l337:
					depth--
					add(rulePegText, position335)
				}
				if !_rules[ruleAction14]() {
					goto l333
				}
				depth--
				add(ruleHaving, position334)
			}
			return true
		l333:
			position, tokenIndex, depth = position333, tokenIndex333, depth333
			return false
		},
		/* 19 SourceSinkSpecs <- <(<(('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H') sp SourceSinkParam sp (',' sp SourceSinkParam)*)?> Action15)> */
		func() bool {
			position350, tokenIndex350, depth350 := position, tokenIndex, depth
			{
				position351 := position
				depth++
				{
					position352 := position
					depth++
					{
						position353, tokenIndex353, depth353 := position, tokenIndex, depth
						{
							position355, tokenIndex355, depth355 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l356
							}
							position++
							goto l355
						l356:
							position, tokenIndex, depth = position355, tokenIndex355, depth355
							if buffer[position] != rune('W') {
								goto l353
							}
							position++
						}
					l355:
						{
							position357, tokenIndex357, depth357 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l358
							}
							position++
							goto l357
						l358:
							position, tokenIndex, depth = position357, tokenIndex357, depth357
							if buffer[position] != rune('I') {
								goto l353
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
								goto l353
							}
							position++
						}
					l359:
						{
							position361, tokenIndex361, depth361 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l362
							}
							position++
							goto l361
						l362:
							position, tokenIndex, depth = position361, tokenIndex361, depth361
							if buffer[position] != rune('H') {
								goto l353
							}
							position++
						}
					l361:
						if !_rules[rulesp]() {
							goto l353
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l353
						}
						if !_rules[rulesp]() {
							goto l353
						}
					l363:
						{
							position364, tokenIndex364, depth364 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l364
							}
							position++
							if !_rules[rulesp]() {
								goto l364
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l364
							}
							goto l363
						l364:
							position, tokenIndex, depth = position364, tokenIndex364, depth364
						}
						goto l354
					l353:
						position, tokenIndex, depth = position353, tokenIndex353, depth353
					}
				l354:
					depth--
					add(rulePegText, position352)
				}
				if !_rules[ruleAction15]() {
					goto l350
				}
				depth--
				add(ruleSourceSinkSpecs, position351)
			}
			return true
		l350:
			position, tokenIndex, depth = position350, tokenIndex350, depth350
			return false
		},
		/* 20 SourceSinkParam <- <(SourceSinkParamKey '=' SourceSinkParamVal Action16)> */
		func() bool {
			position365, tokenIndex365, depth365 := position, tokenIndex, depth
			{
				position366 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l365
				}
				if buffer[position] != rune('=') {
					goto l365
				}
				position++
				if !_rules[ruleSourceSinkParamVal]() {
					goto l365
				}
				if !_rules[ruleAction16]() {
					goto l365
				}
				depth--
				add(ruleSourceSinkParam, position366)
			}
			return true
		l365:
			position, tokenIndex, depth = position365, tokenIndex365, depth365
			return false
		},
		/* 21 Expression <- <orExpr> */
		func() bool {
			position367, tokenIndex367, depth367 := position, tokenIndex, depth
			{
				position368 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l367
				}
				depth--
				add(ruleExpression, position368)
			}
			return true
		l367:
			position, tokenIndex, depth = position367, tokenIndex367, depth367
			return false
		},
		/* 22 orExpr <- <(<(andExpr sp (Or sp andExpr)?)> Action17)> */
		func() bool {
			position369, tokenIndex369, depth369 := position, tokenIndex, depth
			{
				position370 := position
				depth++
				{
					position371 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l369
					}
					if !_rules[rulesp]() {
						goto l369
					}
					{
						position372, tokenIndex372, depth372 := position, tokenIndex, depth
						if !_rules[ruleOr]() {
							goto l372
						}
						if !_rules[rulesp]() {
							goto l372
						}
						if !_rules[ruleandExpr]() {
							goto l372
						}
						goto l373
					l372:
						position, tokenIndex, depth = position372, tokenIndex372, depth372
					}
				l373:
					depth--
					add(rulePegText, position371)
				}
				if !_rules[ruleAction17]() {
					goto l369
				}
				depth--
				add(ruleorExpr, position370)
			}
			return true
		l369:
			position, tokenIndex, depth = position369, tokenIndex369, depth369
			return false
		},
		/* 23 andExpr <- <(<(comparisonExpr sp (And sp comparisonExpr)?)> Action18)> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				{
					position376 := position
					depth++
					if !_rules[rulecomparisonExpr]() {
						goto l374
					}
					if !_rules[rulesp]() {
						goto l374
					}
					{
						position377, tokenIndex377, depth377 := position, tokenIndex, depth
						if !_rules[ruleAnd]() {
							goto l377
						}
						if !_rules[rulesp]() {
							goto l377
						}
						if !_rules[rulecomparisonExpr]() {
							goto l377
						}
						goto l378
					l377:
						position, tokenIndex, depth = position377, tokenIndex377, depth377
					}
				l378:
					depth--
					add(rulePegText, position376)
				}
				if !_rules[ruleAction18]() {
					goto l374
				}
				depth--
				add(ruleandExpr, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
			return false
		},
		/* 24 comparisonExpr <- <(<(termExpr sp (ComparisonOp sp termExpr)?)> Action19)> */
		func() bool {
			position379, tokenIndex379, depth379 := position, tokenIndex, depth
			{
				position380 := position
				depth++
				{
					position381 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l379
					}
					if !_rules[rulesp]() {
						goto l379
					}
					{
						position382, tokenIndex382, depth382 := position, tokenIndex, depth
						if !_rules[ruleComparisonOp]() {
							goto l382
						}
						if !_rules[rulesp]() {
							goto l382
						}
						if !_rules[ruletermExpr]() {
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
				if !_rules[ruleAction19]() {
					goto l379
				}
				depth--
				add(rulecomparisonExpr, position380)
			}
			return true
		l379:
			position, tokenIndex, depth = position379, tokenIndex379, depth379
			return false
		},
		/* 25 termExpr <- <(<(productExpr sp (PlusMinusOp sp productExpr)?)> Action20)> */
		func() bool {
			position384, tokenIndex384, depth384 := position, tokenIndex, depth
			{
				position385 := position
				depth++
				{
					position386 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l384
					}
					if !_rules[rulesp]() {
						goto l384
					}
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if !_rules[rulePlusMinusOp]() {
							goto l387
						}
						if !_rules[rulesp]() {
							goto l387
						}
						if !_rules[ruleproductExpr]() {
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
				if !_rules[ruleAction20]() {
					goto l384
				}
				depth--
				add(ruletermExpr, position385)
			}
			return true
		l384:
			position, tokenIndex, depth = position384, tokenIndex384, depth384
			return false
		},
		/* 26 productExpr <- <(<(baseExpr sp (MultDivOp sp baseExpr)?)> Action21)> */
		func() bool {
			position389, tokenIndex389, depth389 := position, tokenIndex, depth
			{
				position390 := position
				depth++
				{
					position391 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l389
					}
					if !_rules[rulesp]() {
						goto l389
					}
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						if !_rules[ruleMultDivOp]() {
							goto l392
						}
						if !_rules[rulesp]() {
							goto l392
						}
						if !_rules[rulebaseExpr]() {
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
				if !_rules[ruleAction21]() {
					goto l389
				}
				depth--
				add(ruleproductExpr, position390)
			}
			return true
		l389:
			position, tokenIndex, depth = position389, tokenIndex389, depth389
			return false
		},
		/* 27 baseExpr <- <(('(' sp Expression sp ')') / BooleanLiteral / FuncApp / ColumnName / Wildcard / Literal)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l397
					}
					position++
					if !_rules[rulesp]() {
						goto l397
					}
					if !_rules[ruleExpression]() {
						goto l397
					}
					if !_rules[rulesp]() {
						goto l397
					}
					if buffer[position] != rune(')') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if !_rules[ruleBooleanLiteral]() {
						goto l398
					}
					goto l396
				l398:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if !_rules[ruleFuncApp]() {
						goto l399
					}
					goto l396
				l399:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if !_rules[ruleColumnName]() {
						goto l400
					}
					goto l396
				l400:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if !_rules[ruleWildcard]() {
						goto l401
					}
					goto l396
				l401:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if !_rules[ruleLiteral]() {
						goto l394
					}
				}
			l396:
				depth--
				add(rulebaseExpr, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 28 FuncApp <- <(Function sp '(' sp FuncParams sp ')' Action22)> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l402
				}
				if !_rules[rulesp]() {
					goto l402
				}
				if buffer[position] != rune('(') {
					goto l402
				}
				position++
				if !_rules[rulesp]() {
					goto l402
				}
				if !_rules[ruleFuncParams]() {
					goto l402
				}
				if !_rules[rulesp]() {
					goto l402
				}
				if buffer[position] != rune(')') {
					goto l402
				}
				position++
				if !_rules[ruleAction22]() {
					goto l402
				}
				depth--
				add(ruleFuncApp, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 29 FuncParams <- <(<(Expression sp (',' sp Expression)*)> Action23)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				{
					position406 := position
					depth++
					if !_rules[ruleExpression]() {
						goto l404
					}
					if !_rules[rulesp]() {
						goto l404
					}
				l407:
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l408
						}
						position++
						if !_rules[rulesp]() {
							goto l408
						}
						if !_rules[ruleExpression]() {
							goto l408
						}
						goto l407
					l408:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
					}
					depth--
					add(rulePegText, position406)
				}
				if !_rules[ruleAction23]() {
					goto l404
				}
				depth--
				add(ruleFuncParams, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 30 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position409, tokenIndex409, depth409 := position, tokenIndex, depth
			{
				position410 := position
				depth++
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l412
					}
					goto l411
				l412:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if !_rules[ruleNumericLiteral]() {
						goto l413
					}
					goto l411
				l413:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
					if !_rules[ruleStringLiteral]() {
						goto l409
					}
				}
			l411:
				depth--
				add(ruleLiteral, position410)
			}
			return true
		l409:
			position, tokenIndex, depth = position409, tokenIndex409, depth409
			return false
		},
		/* 31 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position414, tokenIndex414, depth414 := position, tokenIndex, depth
			{
				position415 := position
				depth++
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l417
					}
					goto l416
				l417:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if !_rules[ruleNotEqual]() {
						goto l418
					}
					goto l416
				l418:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if !_rules[ruleLessOrEqual]() {
						goto l419
					}
					goto l416
				l419:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if !_rules[ruleLess]() {
						goto l420
					}
					goto l416
				l420:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if !_rules[ruleGreaterOrEqual]() {
						goto l421
					}
					goto l416
				l421:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if !_rules[ruleGreater]() {
						goto l422
					}
					goto l416
				l422:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if !_rules[ruleNotEqual]() {
						goto l414
					}
				}
			l416:
				depth--
				add(ruleComparisonOp, position415)
			}
			return true
		l414:
			position, tokenIndex, depth = position414, tokenIndex414, depth414
			return false
		},
		/* 32 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position423, tokenIndex423, depth423 := position, tokenIndex, depth
			{
				position424 := position
				depth++
				{
					position425, tokenIndex425, depth425 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l426
					}
					goto l425
				l426:
					position, tokenIndex, depth = position425, tokenIndex425, depth425
					if !_rules[ruleMinus]() {
						goto l423
					}
				}
			l425:
				depth--
				add(rulePlusMinusOp, position424)
			}
			return true
		l423:
			position, tokenIndex, depth = position423, tokenIndex423, depth423
			return false
		},
		/* 33 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position427, tokenIndex427, depth427 := position, tokenIndex, depth
			{
				position428 := position
				depth++
				{
					position429, tokenIndex429, depth429 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l430
					}
					goto l429
				l430:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if !_rules[ruleDivide]() {
						goto l431
					}
					goto l429
				l431:
					position, tokenIndex, depth = position429, tokenIndex429, depth429
					if !_rules[ruleModulo]() {
						goto l427
					}
				}
			l429:
				depth--
				add(ruleMultDivOp, position428)
			}
			return true
		l427:
			position, tokenIndex, depth = position427, tokenIndex427, depth427
			return false
		},
		/* 34 Relation <- <(<ident> Action24)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				{
					position434 := position
					depth++
					if !_rules[ruleident]() {
						goto l432
					}
					depth--
					add(rulePegText, position434)
				}
				if !_rules[ruleAction24]() {
					goto l432
				}
				depth--
				add(ruleRelation, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 35 ColumnName <- <(<ident> Action25)> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				{
					position437 := position
					depth++
					if !_rules[ruleident]() {
						goto l435
					}
					depth--
					add(rulePegText, position437)
				}
				if !_rules[ruleAction25]() {
					goto l435
				}
				depth--
				add(ruleColumnName, position436)
			}
			return true
		l435:
			position, tokenIndex, depth = position435, tokenIndex435, depth435
			return false
		},
		/* 36 NumericLiteral <- <(<('-'? [0-9]+)> Action26)> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				{
					position440 := position
					depth++
					{
						position441, tokenIndex441, depth441 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l441
						}
						position++
						goto l442
					l441:
						position, tokenIndex, depth = position441, tokenIndex441, depth441
					}
				l442:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l438
					}
					position++
				l443:
					{
						position444, tokenIndex444, depth444 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l444
						}
						position++
						goto l443
					l444:
						position, tokenIndex, depth = position444, tokenIndex444, depth444
					}
					depth--
					add(rulePegText, position440)
				}
				if !_rules[ruleAction26]() {
					goto l438
				}
				depth--
				add(ruleNumericLiteral, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 37 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action27)> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				{
					position447 := position
					depth++
					{
						position448, tokenIndex448, depth448 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l448
						}
						position++
						goto l449
					l448:
						position, tokenIndex, depth = position448, tokenIndex448, depth448
					}
				l449:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l445
					}
					position++
				l450:
					{
						position451, tokenIndex451, depth451 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l451
						}
						position++
						goto l450
					l451:
						position, tokenIndex, depth = position451, tokenIndex451, depth451
					}
					if buffer[position] != rune('.') {
						goto l445
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l445
					}
					position++
				l452:
					{
						position453, tokenIndex453, depth453 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l453
						}
						position++
						goto l452
					l453:
						position, tokenIndex, depth = position453, tokenIndex453, depth453
					}
					depth--
					add(rulePegText, position447)
				}
				if !_rules[ruleAction27]() {
					goto l445
				}
				depth--
				add(ruleFloatLiteral, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 38 Function <- <(<ident> Action28)> */
		func() bool {
			position454, tokenIndex454, depth454 := position, tokenIndex, depth
			{
				position455 := position
				depth++
				{
					position456 := position
					depth++
					if !_rules[ruleident]() {
						goto l454
					}
					depth--
					add(rulePegText, position456)
				}
				if !_rules[ruleAction28]() {
					goto l454
				}
				depth--
				add(ruleFunction, position455)
			}
			return true
		l454:
			position, tokenIndex, depth = position454, tokenIndex454, depth454
			return false
		},
		/* 39 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				{
					position459, tokenIndex459, depth459 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l460
					}
					goto l459
				l460:
					position, tokenIndex, depth = position459, tokenIndex459, depth459
					if !_rules[ruleFALSE]() {
						goto l457
					}
				}
			l459:
				depth--
				add(ruleBooleanLiteral, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 40 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action29)> */
		func() bool {
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				{
					position463 := position
					depth++
					{
						position464, tokenIndex464, depth464 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l465
						}
						position++
						goto l464
					l465:
						position, tokenIndex, depth = position464, tokenIndex464, depth464
						if buffer[position] != rune('T') {
							goto l461
						}
						position++
					}
				l464:
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l467
						}
						position++
						goto l466
					l467:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
						if buffer[position] != rune('R') {
							goto l461
						}
						position++
					}
				l466:
					{
						position468, tokenIndex468, depth468 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l469
						}
						position++
						goto l468
					l469:
						position, tokenIndex, depth = position468, tokenIndex468, depth468
						if buffer[position] != rune('U') {
							goto l461
						}
						position++
					}
				l468:
					{
						position470, tokenIndex470, depth470 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l471
						}
						position++
						goto l470
					l471:
						position, tokenIndex, depth = position470, tokenIndex470, depth470
						if buffer[position] != rune('E') {
							goto l461
						}
						position++
					}
				l470:
					depth--
					add(rulePegText, position463)
				}
				if !_rules[ruleAction29]() {
					goto l461
				}
				depth--
				add(ruleTRUE, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 41 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action30)> */
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
						if buffer[position] != rune('f') {
							goto l476
						}
						position++
						goto l475
					l476:
						position, tokenIndex, depth = position475, tokenIndex475, depth475
						if buffer[position] != rune('F') {
							goto l472
						}
						position++
					}
				l475:
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l478
						}
						position++
						goto l477
					l478:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
						if buffer[position] != rune('A') {
							goto l472
						}
						position++
					}
				l477:
					{
						position479, tokenIndex479, depth479 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l480
						}
						position++
						goto l479
					l480:
						position, tokenIndex, depth = position479, tokenIndex479, depth479
						if buffer[position] != rune('L') {
							goto l472
						}
						position++
					}
				l479:
					{
						position481, tokenIndex481, depth481 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l482
						}
						position++
						goto l481
					l482:
						position, tokenIndex, depth = position481, tokenIndex481, depth481
						if buffer[position] != rune('S') {
							goto l472
						}
						position++
					}
				l481:
					{
						position483, tokenIndex483, depth483 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l484
						}
						position++
						goto l483
					l484:
						position, tokenIndex, depth = position483, tokenIndex483, depth483
						if buffer[position] != rune('E') {
							goto l472
						}
						position++
					}
				l483:
					depth--
					add(rulePegText, position474)
				}
				if !_rules[ruleAction30]() {
					goto l472
				}
				depth--
				add(ruleFALSE, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 42 Wildcard <- <(<'*'> Action31)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				{
					position487 := position
					depth++
					if buffer[position] != rune('*') {
						goto l485
					}
					position++
					depth--
					add(rulePegText, position487)
				}
				if !_rules[ruleAction31]() {
					goto l485
				}
				depth--
				add(ruleWildcard, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 43 StringLiteral <- <(<('\'' (('\'' '\'') / (!'\'' .))* '\'')> Action32)> */
		func() bool {
			position488, tokenIndex488, depth488 := position, tokenIndex, depth
			{
				position489 := position
				depth++
				{
					position490 := position
					depth++
					if buffer[position] != rune('\'') {
						goto l488
					}
					position++
				l491:
					{
						position492, tokenIndex492, depth492 := position, tokenIndex, depth
						{
							position493, tokenIndex493, depth493 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l494
							}
							position++
							if buffer[position] != rune('\'') {
								goto l494
							}
							position++
							goto l493
						l494:
							position, tokenIndex, depth = position493, tokenIndex493, depth493
							{
								position495, tokenIndex495, depth495 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l495
								}
								position++
								goto l492
							l495:
								position, tokenIndex, depth = position495, tokenIndex495, depth495
							}
							if !matchDot() {
								goto l492
							}
						}
					l493:
						goto l491
					l492:
						position, tokenIndex, depth = position492, tokenIndex492, depth492
					}
					if buffer[position] != rune('\'') {
						goto l488
					}
					position++
					depth--
					add(rulePegText, position490)
				}
				if !_rules[ruleAction32]() {
					goto l488
				}
				depth--
				add(ruleStringLiteral, position489)
			}
			return true
		l488:
			position, tokenIndex, depth = position488, tokenIndex488, depth488
			return false
		},
		/* 44 Emitter <- <(ISTREAM / DSTREAM / RSTREAM)> */
		func() bool {
			position496, tokenIndex496, depth496 := position, tokenIndex, depth
			{
				position497 := position
				depth++
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l499
					}
					goto l498
				l499:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleDSTREAM]() {
						goto l500
					}
					goto l498
				l500:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if !_rules[ruleRSTREAM]() {
						goto l496
					}
				}
			l498:
				depth--
				add(ruleEmitter, position497)
			}
			return true
		l496:
			position, tokenIndex, depth = position496, tokenIndex496, depth496
			return false
		},
		/* 45 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action33)> */
		func() bool {
			position501, tokenIndex501, depth501 := position, tokenIndex, depth
			{
				position502 := position
				depth++
				{
					position503 := position
					depth++
					{
						position504, tokenIndex504, depth504 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l505
						}
						position++
						goto l504
					l505:
						position, tokenIndex, depth = position504, tokenIndex504, depth504
						if buffer[position] != rune('I') {
							goto l501
						}
						position++
					}
				l504:
					{
						position506, tokenIndex506, depth506 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l507
						}
						position++
						goto l506
					l507:
						position, tokenIndex, depth = position506, tokenIndex506, depth506
						if buffer[position] != rune('S') {
							goto l501
						}
						position++
					}
				l506:
					{
						position508, tokenIndex508, depth508 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l509
						}
						position++
						goto l508
					l509:
						position, tokenIndex, depth = position508, tokenIndex508, depth508
						if buffer[position] != rune('T') {
							goto l501
						}
						position++
					}
				l508:
					{
						position510, tokenIndex510, depth510 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l511
						}
						position++
						goto l510
					l511:
						position, tokenIndex, depth = position510, tokenIndex510, depth510
						if buffer[position] != rune('R') {
							goto l501
						}
						position++
					}
				l510:
					{
						position512, tokenIndex512, depth512 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l513
						}
						position++
						goto l512
					l513:
						position, tokenIndex, depth = position512, tokenIndex512, depth512
						if buffer[position] != rune('E') {
							goto l501
						}
						position++
					}
				l512:
					{
						position514, tokenIndex514, depth514 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l515
						}
						position++
						goto l514
					l515:
						position, tokenIndex, depth = position514, tokenIndex514, depth514
						if buffer[position] != rune('A') {
							goto l501
						}
						position++
					}
				l514:
					{
						position516, tokenIndex516, depth516 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l517
						}
						position++
						goto l516
					l517:
						position, tokenIndex, depth = position516, tokenIndex516, depth516
						if buffer[position] != rune('M') {
							goto l501
						}
						position++
					}
				l516:
					depth--
					add(rulePegText, position503)
				}
				if !_rules[ruleAction33]() {
					goto l501
				}
				depth--
				add(ruleISTREAM, position502)
			}
			return true
		l501:
			position, tokenIndex, depth = position501, tokenIndex501, depth501
			return false
		},
		/* 46 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action34)> */
		func() bool {
			position518, tokenIndex518, depth518 := position, tokenIndex, depth
			{
				position519 := position
				depth++
				{
					position520 := position
					depth++
					{
						position521, tokenIndex521, depth521 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l522
						}
						position++
						goto l521
					l522:
						position, tokenIndex, depth = position521, tokenIndex521, depth521
						if buffer[position] != rune('D') {
							goto l518
						}
						position++
					}
				l521:
					{
						position523, tokenIndex523, depth523 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l524
						}
						position++
						goto l523
					l524:
						position, tokenIndex, depth = position523, tokenIndex523, depth523
						if buffer[position] != rune('S') {
							goto l518
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
							goto l518
						}
						position++
					}
				l525:
					{
						position527, tokenIndex527, depth527 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l528
						}
						position++
						goto l527
					l528:
						position, tokenIndex, depth = position527, tokenIndex527, depth527
						if buffer[position] != rune('R') {
							goto l518
						}
						position++
					}
				l527:
					{
						position529, tokenIndex529, depth529 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l530
						}
						position++
						goto l529
					l530:
						position, tokenIndex, depth = position529, tokenIndex529, depth529
						if buffer[position] != rune('E') {
							goto l518
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
							goto l518
						}
						position++
					}
				l531:
					{
						position533, tokenIndex533, depth533 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l534
						}
						position++
						goto l533
					l534:
						position, tokenIndex, depth = position533, tokenIndex533, depth533
						if buffer[position] != rune('M') {
							goto l518
						}
						position++
					}
				l533:
					depth--
					add(rulePegText, position520)
				}
				if !_rules[ruleAction34]() {
					goto l518
				}
				depth--
				add(ruleDSTREAM, position519)
			}
			return true
		l518:
			position, tokenIndex, depth = position518, tokenIndex518, depth518
			return false
		},
		/* 47 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action35)> */
		func() bool {
			position535, tokenIndex535, depth535 := position, tokenIndex, depth
			{
				position536 := position
				depth++
				{
					position537 := position
					depth++
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
							goto l535
						}
						position++
					}
				l538:
					{
						position540, tokenIndex540, depth540 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l541
						}
						position++
						goto l540
					l541:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
						if buffer[position] != rune('S') {
							goto l535
						}
						position++
					}
				l540:
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l543
						}
						position++
						goto l542
					l543:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
						if buffer[position] != rune('T') {
							goto l535
						}
						position++
					}
				l542:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l545
						}
						position++
						goto l544
					l545:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
						if buffer[position] != rune('R') {
							goto l535
						}
						position++
					}
				l544:
					{
						position546, tokenIndex546, depth546 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l547
						}
						position++
						goto l546
					l547:
						position, tokenIndex, depth = position546, tokenIndex546, depth546
						if buffer[position] != rune('E') {
							goto l535
						}
						position++
					}
				l546:
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
							goto l535
						}
						position++
					}
				l548:
					{
						position550, tokenIndex550, depth550 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l551
						}
						position++
						goto l550
					l551:
						position, tokenIndex, depth = position550, tokenIndex550, depth550
						if buffer[position] != rune('M') {
							goto l535
						}
						position++
					}
				l550:
					depth--
					add(rulePegText, position537)
				}
				if !_rules[ruleAction35]() {
					goto l535
				}
				depth--
				add(ruleRSTREAM, position536)
			}
			return true
		l535:
			position, tokenIndex, depth = position535, tokenIndex535, depth535
			return false
		},
		/* 48 RangeUnit <- <(TUPLES / SECONDS)> */
		func() bool {
			position552, tokenIndex552, depth552 := position, tokenIndex, depth
			{
				position553 := position
				depth++
				{
					position554, tokenIndex554, depth554 := position, tokenIndex, depth
					if !_rules[ruleTUPLES]() {
						goto l555
					}
					goto l554
				l555:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
					if !_rules[ruleSECONDS]() {
						goto l552
					}
				}
			l554:
				depth--
				add(ruleRangeUnit, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 49 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action36)> */
		func() bool {
			position556, tokenIndex556, depth556 := position, tokenIndex, depth
			{
				position557 := position
				depth++
				{
					position558 := position
					depth++
					{
						position559, tokenIndex559, depth559 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l560
						}
						position++
						goto l559
					l560:
						position, tokenIndex, depth = position559, tokenIndex559, depth559
						if buffer[position] != rune('T') {
							goto l556
						}
						position++
					}
				l559:
					{
						position561, tokenIndex561, depth561 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l562
						}
						position++
						goto l561
					l562:
						position, tokenIndex, depth = position561, tokenIndex561, depth561
						if buffer[position] != rune('U') {
							goto l556
						}
						position++
					}
				l561:
					{
						position563, tokenIndex563, depth563 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l564
						}
						position++
						goto l563
					l564:
						position, tokenIndex, depth = position563, tokenIndex563, depth563
						if buffer[position] != rune('P') {
							goto l556
						}
						position++
					}
				l563:
					{
						position565, tokenIndex565, depth565 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l566
						}
						position++
						goto l565
					l566:
						position, tokenIndex, depth = position565, tokenIndex565, depth565
						if buffer[position] != rune('L') {
							goto l556
						}
						position++
					}
				l565:
					{
						position567, tokenIndex567, depth567 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l568
						}
						position++
						goto l567
					l568:
						position, tokenIndex, depth = position567, tokenIndex567, depth567
						if buffer[position] != rune('E') {
							goto l556
						}
						position++
					}
				l567:
					{
						position569, tokenIndex569, depth569 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l570
						}
						position++
						goto l569
					l570:
						position, tokenIndex, depth = position569, tokenIndex569, depth569
						if buffer[position] != rune('S') {
							goto l556
						}
						position++
					}
				l569:
					depth--
					add(rulePegText, position558)
				}
				if !_rules[ruleAction36]() {
					goto l556
				}
				depth--
				add(ruleTUPLES, position557)
			}
			return true
		l556:
			position, tokenIndex, depth = position556, tokenIndex556, depth556
			return false
		},
		/* 50 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action37)> */
		func() bool {
			position571, tokenIndex571, depth571 := position, tokenIndex, depth
			{
				position572 := position
				depth++
				{
					position573 := position
					depth++
					{
						position574, tokenIndex574, depth574 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l575
						}
						position++
						goto l574
					l575:
						position, tokenIndex, depth = position574, tokenIndex574, depth574
						if buffer[position] != rune('S') {
							goto l571
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
							goto l571
						}
						position++
					}
				l576:
					{
						position578, tokenIndex578, depth578 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l579
						}
						position++
						goto l578
					l579:
						position, tokenIndex, depth = position578, tokenIndex578, depth578
						if buffer[position] != rune('C') {
							goto l571
						}
						position++
					}
				l578:
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l581
						}
						position++
						goto l580
					l581:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
						if buffer[position] != rune('O') {
							goto l571
						}
						position++
					}
				l580:
					{
						position582, tokenIndex582, depth582 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l583
						}
						position++
						goto l582
					l583:
						position, tokenIndex, depth = position582, tokenIndex582, depth582
						if buffer[position] != rune('N') {
							goto l571
						}
						position++
					}
				l582:
					{
						position584, tokenIndex584, depth584 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l585
						}
						position++
						goto l584
					l585:
						position, tokenIndex, depth = position584, tokenIndex584, depth584
						if buffer[position] != rune('D') {
							goto l571
						}
						position++
					}
				l584:
					{
						position586, tokenIndex586, depth586 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l587
						}
						position++
						goto l586
					l587:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
						if buffer[position] != rune('S') {
							goto l571
						}
						position++
					}
				l586:
					depth--
					add(rulePegText, position573)
				}
				if !_rules[ruleAction37]() {
					goto l571
				}
				depth--
				add(ruleSECONDS, position572)
			}
			return true
		l571:
			position, tokenIndex, depth = position571, tokenIndex571, depth571
			return false
		},
		/* 51 SourceSinkName <- <(<ident> Action38)> */
		func() bool {
			position588, tokenIndex588, depth588 := position, tokenIndex, depth
			{
				position589 := position
				depth++
				{
					position590 := position
					depth++
					if !_rules[ruleident]() {
						goto l588
					}
					depth--
					add(rulePegText, position590)
				}
				if !_rules[ruleAction38]() {
					goto l588
				}
				depth--
				add(ruleSourceSinkName, position589)
			}
			return true
		l588:
			position, tokenIndex, depth = position588, tokenIndex588, depth588
			return false
		},
		/* 52 SourceSinkType <- <(<ident> Action39)> */
		func() bool {
			position591, tokenIndex591, depth591 := position, tokenIndex, depth
			{
				position592 := position
				depth++
				{
					position593 := position
					depth++
					if !_rules[ruleident]() {
						goto l591
					}
					depth--
					add(rulePegText, position593)
				}
				if !_rules[ruleAction39]() {
					goto l591
				}
				depth--
				add(ruleSourceSinkType, position592)
			}
			return true
		l591:
			position, tokenIndex, depth = position591, tokenIndex591, depth591
			return false
		},
		/* 53 SourceSinkParamKey <- <(<ident> Action40)> */
		func() bool {
			position594, tokenIndex594, depth594 := position, tokenIndex, depth
			{
				position595 := position
				depth++
				{
					position596 := position
					depth++
					if !_rules[ruleident]() {
						goto l594
					}
					depth--
					add(rulePegText, position596)
				}
				if !_rules[ruleAction40]() {
					goto l594
				}
				depth--
				add(ruleSourceSinkParamKey, position595)
			}
			return true
		l594:
			position, tokenIndex, depth = position594, tokenIndex594, depth594
			return false
		},
		/* 54 SourceSinkParamVal <- <(<([a-z] / [A-Z] / [0-9] / '_')+> Action41)> */
		func() bool {
			position597, tokenIndex597, depth597 := position, tokenIndex, depth
			{
				position598 := position
				depth++
				{
					position599 := position
					depth++
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l603
						}
						position++
						goto l602
					l603:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l604
						}
						position++
						goto l602
					l604:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l605
						}
						position++
						goto l602
					l605:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if buffer[position] != rune('_') {
							goto l597
						}
						position++
					}
				l602:
				l600:
					{
						position601, tokenIndex601, depth601 := position, tokenIndex, depth
						{
							position606, tokenIndex606, depth606 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l607
							}
							position++
							goto l606
						l607:
							position, tokenIndex, depth = position606, tokenIndex606, depth606
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l608
							}
							position++
							goto l606
						l608:
							position, tokenIndex, depth = position606, tokenIndex606, depth606
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l609
							}
							position++
							goto l606
						l609:
							position, tokenIndex, depth = position606, tokenIndex606, depth606
							if buffer[position] != rune('_') {
								goto l601
							}
							position++
						}
					l606:
						goto l600
					l601:
						position, tokenIndex, depth = position601, tokenIndex601, depth601
					}
					depth--
					add(rulePegText, position599)
				}
				if !_rules[ruleAction41]() {
					goto l597
				}
				depth--
				add(ruleSourceSinkParamVal, position598)
			}
			return true
		l597:
			position, tokenIndex, depth = position597, tokenIndex597, depth597
			return false
		},
		/* 55 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action42)> */
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
						if buffer[position] != rune('o') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('O') {
							goto l610
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('R') {
							goto l610
						}
						position++
					}
				l615:
					depth--
					add(rulePegText, position612)
				}
				if !_rules[ruleAction42]() {
					goto l610
				}
				depth--
				add(ruleOr, position611)
			}
			return true
		l610:
			position, tokenIndex, depth = position610, tokenIndex610, depth610
			return false
		},
		/* 56 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action43)> */
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
						if buffer[position] != rune('a') {
							goto l621
						}
						position++
						goto l620
					l621:
						position, tokenIndex, depth = position620, tokenIndex620, depth620
						if buffer[position] != rune('A') {
							goto l617
						}
						position++
					}
				l620:
					{
						position622, tokenIndex622, depth622 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l623
						}
						position++
						goto l622
					l623:
						position, tokenIndex, depth = position622, tokenIndex622, depth622
						if buffer[position] != rune('N') {
							goto l617
						}
						position++
					}
				l622:
					{
						position624, tokenIndex624, depth624 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l625
						}
						position++
						goto l624
					l625:
						position, tokenIndex, depth = position624, tokenIndex624, depth624
						if buffer[position] != rune('D') {
							goto l617
						}
						position++
					}
				l624:
					depth--
					add(rulePegText, position619)
				}
				if !_rules[ruleAction43]() {
					goto l617
				}
				depth--
				add(ruleAnd, position618)
			}
			return true
		l617:
			position, tokenIndex, depth = position617, tokenIndex617, depth617
			return false
		},
		/* 57 Equal <- <(<'='> Action44)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
				{
					position628 := position
					depth++
					if buffer[position] != rune('=') {
						goto l626
					}
					position++
					depth--
					add(rulePegText, position628)
				}
				if !_rules[ruleAction44]() {
					goto l626
				}
				depth--
				add(ruleEqual, position627)
			}
			return true
		l626:
			position, tokenIndex, depth = position626, tokenIndex626, depth626
			return false
		},
		/* 58 Less <- <(<'<'> Action45)> */
		func() bool {
			position629, tokenIndex629, depth629 := position, tokenIndex, depth
			{
				position630 := position
				depth++
				{
					position631 := position
					depth++
					if buffer[position] != rune('<') {
						goto l629
					}
					position++
					depth--
					add(rulePegText, position631)
				}
				if !_rules[ruleAction45]() {
					goto l629
				}
				depth--
				add(ruleLess, position630)
			}
			return true
		l629:
			position, tokenIndex, depth = position629, tokenIndex629, depth629
			return false
		},
		/* 59 LessOrEqual <- <(<('<' '=')> Action46)> */
		func() bool {
			position632, tokenIndex632, depth632 := position, tokenIndex, depth
			{
				position633 := position
				depth++
				{
					position634 := position
					depth++
					if buffer[position] != rune('<') {
						goto l632
					}
					position++
					if buffer[position] != rune('=') {
						goto l632
					}
					position++
					depth--
					add(rulePegText, position634)
				}
				if !_rules[ruleAction46]() {
					goto l632
				}
				depth--
				add(ruleLessOrEqual, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 60 Greater <- <(<'>'> Action47)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				{
					position637 := position
					depth++
					if buffer[position] != rune('>') {
						goto l635
					}
					position++
					depth--
					add(rulePegText, position637)
				}
				if !_rules[ruleAction47]() {
					goto l635
				}
				depth--
				add(ruleGreater, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 61 GreaterOrEqual <- <(<('>' '=')> Action48)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					if buffer[position] != rune('>') {
						goto l638
					}
					position++
					if buffer[position] != rune('=') {
						goto l638
					}
					position++
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction48]() {
					goto l638
				}
				depth--
				add(ruleGreaterOrEqual, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 62 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action49)> */
		func() bool {
			position641, tokenIndex641, depth641 := position, tokenIndex, depth
			{
				position642 := position
				depth++
				{
					position643 := position
					depth++
					{
						position644, tokenIndex644, depth644 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l645
						}
						position++
						if buffer[position] != rune('=') {
							goto l645
						}
						position++
						goto l644
					l645:
						position, tokenIndex, depth = position644, tokenIndex644, depth644
						if buffer[position] != rune('<') {
							goto l641
						}
						position++
						if buffer[position] != rune('>') {
							goto l641
						}
						position++
					}
				l644:
					depth--
					add(rulePegText, position643)
				}
				if !_rules[ruleAction49]() {
					goto l641
				}
				depth--
				add(ruleNotEqual, position642)
			}
			return true
		l641:
			position, tokenIndex, depth = position641, tokenIndex641, depth641
			return false
		},
		/* 63 Plus <- <(<'+'> Action50)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				{
					position648 := position
					depth++
					if buffer[position] != rune('+') {
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
				add(rulePlus, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 64 Minus <- <(<'-'> Action51)> */
		func() bool {
			position649, tokenIndex649, depth649 := position, tokenIndex, depth
			{
				position650 := position
				depth++
				{
					position651 := position
					depth++
					if buffer[position] != rune('-') {
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
				add(ruleMinus, position650)
			}
			return true
		l649:
			position, tokenIndex, depth = position649, tokenIndex649, depth649
			return false
		},
		/* 65 Multiply <- <(<'*'> Action52)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				{
					position654 := position
					depth++
					if buffer[position] != rune('*') {
						goto l652
					}
					position++
					depth--
					add(rulePegText, position654)
				}
				if !_rules[ruleAction52]() {
					goto l652
				}
				depth--
				add(ruleMultiply, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 66 Divide <- <(<'/'> Action53)> */
		func() bool {
			position655, tokenIndex655, depth655 := position, tokenIndex, depth
			{
				position656 := position
				depth++
				{
					position657 := position
					depth++
					if buffer[position] != rune('/') {
						goto l655
					}
					position++
					depth--
					add(rulePegText, position657)
				}
				if !_rules[ruleAction53]() {
					goto l655
				}
				depth--
				add(ruleDivide, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 67 Modulo <- <(<'%'> Action54)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				{
					position660 := position
					depth++
					if buffer[position] != rune('%') {
						goto l658
					}
					position++
					depth--
					add(rulePegText, position660)
				}
				if !_rules[ruleAction54]() {
					goto l658
				}
				depth--
				add(ruleModulo, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 68 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position661, tokenIndex661, depth661 := position, tokenIndex, depth
			{
				position662 := position
				depth++
				{
					position663, tokenIndex663, depth663 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l664
					}
					position++
					goto l663
				l664:
					position, tokenIndex, depth = position663, tokenIndex663, depth663
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l661
					}
					position++
				}
			l663:
			l665:
				{
					position666, tokenIndex666, depth666 := position, tokenIndex, depth
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l668
						}
						position++
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l669
						}
						position++
						goto l667
					l669:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l670
						}
						position++
						goto l667
					l670:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						if buffer[position] != rune('_') {
							goto l666
						}
						position++
					}
				l667:
					goto l665
				l666:
					position, tokenIndex, depth = position666, tokenIndex666, depth666
				}
				depth--
				add(ruleident, position662)
			}
			return true
		l661:
			position, tokenIndex, depth = position661, tokenIndex661, depth661
			return false
		},
		/* 69 sp <- <(' ' / '\t' / '\n')*> */
		func() bool {
			{
				position672 := position
				depth++
			l673:
				{
					position674, tokenIndex674, depth674 := position, tokenIndex, depth
					{
						position675, tokenIndex675, depth675 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l676
						}
						position++
						goto l675
					l676:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('\t') {
							goto l677
						}
						position++
						goto l675
					l677:
						position, tokenIndex, depth = position675, tokenIndex675, depth675
						if buffer[position] != rune('\n') {
							goto l674
						}
						position++
					}
				l675:
					goto l673
				l674:
					position, tokenIndex, depth = position674, tokenIndex674, depth674
				}
				depth--
				add(rulesp, position672)
			}
			return true
		},
		/* 71 Action0 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 72 Action1 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 73 Action2 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 74 Action3 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 75 Action4 <- <{
		    p.AssembleCreateStreamFromSource()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 76 Action5 <- <{
		    p.AssembleCreateStreamFromSourceExt()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 77 Action6 <- <{
		    p.AssembleInsertIntoSelect()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 78 Action7 <- <{
		    p.AssembleEmitProjections()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 80 Action8 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 81 Action9 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 82 Action10 <- <{
		    p.AssembleRange()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 83 Action11 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 84 Action12 <- <{
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
		/* 85 Action13 <- <{
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
		/* 86 Action14 <- <{
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
		/* 87 Action15 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 88 Action16 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 89 Action17 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 90 Action18 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 91 Action19 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 92 Action20 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 93 Action21 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 94 Action22 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 95 Action23 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 96 Action24 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRelation(substr))
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 97 Action25 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewColumnName(substr))
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 98 Action26 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 99 Action27 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 100 Action28 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 101 Action29 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 102 Action30 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 103 Action31 <- <{
		    p.PushComponent(begin, end, NewWildcard())
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 104 Action32 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 105 Action33 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 106 Action34 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 107 Action35 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 108 Action36 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 109 Action37 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 110 Action38 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkName(substr))
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 111 Action39 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 112 Action40 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 113 Action41 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamVal(substr))
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 114 Action42 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 115 Action43 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 116 Action44 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 117 Action45 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 118 Action46 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 119 Action47 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 120 Action48 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 121 Action49 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 122 Action50 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 123 Action51 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 124 Action52 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 125 Action53 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 126 Action54 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
	}
	p.rules = _rules
}

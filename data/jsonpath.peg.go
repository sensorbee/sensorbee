package data

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulejsonPath
	rulejsonPathHead
	rulejsonPathNonHead
	rulejsonMapSingleLevel
	rulejsonMapMultipleLevel
	rulejsonMapAccessString
	rulejsonMapAccessBracket
	rulesingleQuotedString
	ruledoubleQuotedString
	rulejsonArrayAccess
	rulejsonArraySlice
	rulejsonArrayPartialSlice
	rulejsonArrayFullSlice
	ruleAction0
	ruleAction1
	ruleAction2
	rulePegText
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"jsonPath",
	"jsonPathHead",
	"jsonPathNonHead",
	"jsonMapSingleLevel",
	"jsonMapMultipleLevel",
	"jsonMapAccessString",
	"jsonMapAccessBracket",
	"singleQuotedString",
	"doubleQuotedString",
	"jsonArrayAccess",
	"jsonArraySlice",
	"jsonArrayPartialSlice",
	"jsonArrayFullSlice",
	"Action0",
	"Action1",
	"Action2",
	"PegText",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",

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

type jsonPeg struct {
	components []extractor
	lastKey    string

	Buffer string
	buffer []rune
	rules  [25]func() bool
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
	p *jsonPeg
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

func (p *jsonPeg) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *jsonPeg) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *jsonPeg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.addMapAccess(p.lastKey)

		case ruleAction1:

			p.addMapAccess(p.lastKey)

		case ruleAction2:

			p.addRecursiveAccess(p.lastKey)

		case ruleAction3:

			substr := string([]rune(buffer)[begin:end])
			p.lastKey = substr

		case ruleAction4:

			substr := string([]rune(buffer)[begin:end])
			p.lastKey = strings.Replace(substr, "''", "'", -1)

		case ruleAction5:

			substr := string([]rune(buffer)[begin:end])
			p.lastKey = strings.Replace(substr, "\"\"", "\"", -1)

		case ruleAction6:

			substr := string([]rune(buffer)[begin:end])
			p.addArrayAccess(substr)

		case ruleAction7:

			substr := string([]rune(buffer)[begin:end])
			p.addArraySlice(substr)

		case ruleAction8:

			substr := string([]rune(buffer)[begin:end])
			p.addArraySlice(substr)

		case ruleAction9:

			p.addArraySlice("0:")

		}
	}
	_, _, _, _ = buffer, text, begin, end
}

func (p *jsonPeg) Init() {
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
		/* 0 jsonPath <- <(jsonPathHead jsonPathNonHead* !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l0
				}
			l2:
				{
					position3, tokenIndex3, depth3 := position, tokenIndex, depth
					if !_rules[rulejsonPathNonHead]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
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
				add(rulejsonPath, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 jsonPathHead <- <((jsonMapAccessString / jsonMapAccessBracket) Action0)> */
		func() bool {
			position5, tokenIndex5, depth5 := position, tokenIndex, depth
			{
				position6 := position
				depth++
				{
					position7, tokenIndex7, depth7 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l8
					}
					goto l7
				l8:
					position, tokenIndex, depth = position7, tokenIndex7, depth7
					if !_rules[rulejsonMapAccessBracket]() {
						goto l5
					}
				}
			l7:
				if !_rules[ruleAction0]() {
					goto l5
				}
				depth--
				add(rulejsonPathHead, position6)
			}
			return true
		l5:
			position, tokenIndex, depth = position5, tokenIndex5, depth5
			return false
		},
		/* 2 jsonPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position9, tokenIndex9, depth9 := position, tokenIndex, depth
			{
				position10 := position
				depth++
				{
					position11, tokenIndex11, depth11 := position, tokenIndex, depth
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l12
					}
					goto l11
				l12:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulejsonMapSingleLevel]() {
						goto l13
					}
					goto l11
				l13:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulejsonArrayFullSlice]() {
						goto l14
					}
					goto l11
				l14:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l15
					}
					goto l11
				l15:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulejsonArraySlice]() {
						goto l16
					}
					goto l11
				l16:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
					if !_rules[rulejsonArrayAccess]() {
						goto l9
					}
				}
			l11:
				depth--
				add(rulejsonPathNonHead, position10)
			}
			return true
		l9:
			position, tokenIndex, depth = position9, tokenIndex9, depth9
			return false
		},
		/* 3 jsonMapSingleLevel <- <((('.' jsonMapAccessString) / jsonMapAccessBracket) Action1)> */
		func() bool {
			position17, tokenIndex17, depth17 := position, tokenIndex, depth
			{
				position18 := position
				depth++
				{
					position19, tokenIndex19, depth19 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l20
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l20
					}
					goto l19
				l20:
					position, tokenIndex, depth = position19, tokenIndex19, depth19
					if !_rules[rulejsonMapAccessBracket]() {
						goto l17
					}
				}
			l19:
				if !_rules[ruleAction1]() {
					goto l17
				}
				depth--
				add(rulejsonMapSingleLevel, position18)
			}
			return true
		l17:
			position, tokenIndex, depth = position17, tokenIndex17, depth17
			return false
		},
		/* 4 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket) Action2)> */
		func() bool {
			position21, tokenIndex21, depth21 := position, tokenIndex, depth
			{
				position22 := position
				depth++
				if buffer[position] != rune('.') {
					goto l21
				}
				position++
				if buffer[position] != rune('.') {
					goto l21
				}
				position++
				{
					position23, tokenIndex23, depth23 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l24
					}
					goto l23
				l24:
					position, tokenIndex, depth = position23, tokenIndex23, depth23
					if !_rules[rulejsonMapAccessBracket]() {
						goto l21
					}
				}
			l23:
				if !_rules[ruleAction2]() {
					goto l21
				}
				depth--
				add(rulejsonMapMultipleLevel, position22)
			}
			return true
		l21:
			position, tokenIndex, depth = position21, tokenIndex21, depth21
			return false
		},
		/* 5 jsonMapAccessString <- <(<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> Action3)> */
		func() bool {
			position25, tokenIndex25, depth25 := position, tokenIndex, depth
			{
				position26 := position
				depth++
				{
					position27 := position
					depth++
					{
						position28, tokenIndex28, depth28 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l29
						}
						position++
						goto l28
					l29:
						position, tokenIndex, depth = position28, tokenIndex28, depth28
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l25
						}
						position++
					}
				l28:
				l30:
					{
						position31, tokenIndex31, depth31 := position, tokenIndex, depth
						{
							position32, tokenIndex32, depth32 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l33
							}
							position++
							goto l32
						l33:
							position, tokenIndex, depth = position32, tokenIndex32, depth32
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l34
							}
							position++
							goto l32
						l34:
							position, tokenIndex, depth = position32, tokenIndex32, depth32
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l35
							}
							position++
							goto l32
						l35:
							position, tokenIndex, depth = position32, tokenIndex32, depth32
							if buffer[position] != rune('_') {
								goto l31
							}
							position++
						}
					l32:
						goto l30
					l31:
						position, tokenIndex, depth = position31, tokenIndex31, depth31
					}
					depth--
					add(rulePegText, position27)
				}
				if !_rules[ruleAction3]() {
					goto l25
				}
				depth--
				add(rulejsonMapAccessString, position26)
			}
			return true
		l25:
			position, tokenIndex, depth = position25, tokenIndex25, depth25
			return false
		},
		/* 6 jsonMapAccessBracket <- <('[' (singleQuotedString / doubleQuotedString) ']')> */
		func() bool {
			position36, tokenIndex36, depth36 := position, tokenIndex, depth
			{
				position37 := position
				depth++
				if buffer[position] != rune('[') {
					goto l36
				}
				position++
				{
					position38, tokenIndex38, depth38 := position, tokenIndex, depth
					if !_rules[rulesingleQuotedString]() {
						goto l39
					}
					goto l38
				l39:
					position, tokenIndex, depth = position38, tokenIndex38, depth38
					if !_rules[ruledoubleQuotedString]() {
						goto l36
					}
				}
			l38:
				if buffer[position] != rune(']') {
					goto l36
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position37)
			}
			return true
		l36:
			position, tokenIndex, depth = position36, tokenIndex36, depth36
			return false
		},
		/* 7 singleQuotedString <- <('\'' <(('\'' '\'') / (!'\'' .))*> '\'' Action4)> */
		func() bool {
			position40, tokenIndex40, depth40 := position, tokenIndex, depth
			{
				position41 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l40
				}
				position++
				{
					position42 := position
					depth++
				l43:
					{
						position44, tokenIndex44, depth44 := position, tokenIndex, depth
						{
							position45, tokenIndex45, depth45 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l46
							}
							position++
							if buffer[position] != rune('\'') {
								goto l46
							}
							position++
							goto l45
						l46:
							position, tokenIndex, depth = position45, tokenIndex45, depth45
							{
								position47, tokenIndex47, depth47 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l47
								}
								position++
								goto l44
							l47:
								position, tokenIndex, depth = position47, tokenIndex47, depth47
							}
							if !matchDot() {
								goto l44
							}
						}
					l45:
						goto l43
					l44:
						position, tokenIndex, depth = position44, tokenIndex44, depth44
					}
					depth--
					add(rulePegText, position42)
				}
				if buffer[position] != rune('\'') {
					goto l40
				}
				position++
				if !_rules[ruleAction4]() {
					goto l40
				}
				depth--
				add(rulesingleQuotedString, position41)
			}
			return true
		l40:
			position, tokenIndex, depth = position40, tokenIndex40, depth40
			return false
		},
		/* 8 doubleQuotedString <- <('"' <(('"' '"') / (!'"' .))*> '"' Action5)> */
		func() bool {
			position48, tokenIndex48, depth48 := position, tokenIndex, depth
			{
				position49 := position
				depth++
				if buffer[position] != rune('"') {
					goto l48
				}
				position++
				{
					position50 := position
					depth++
				l51:
					{
						position52, tokenIndex52, depth52 := position, tokenIndex, depth
						{
							position53, tokenIndex53, depth53 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l54
							}
							position++
							if buffer[position] != rune('"') {
								goto l54
							}
							position++
							goto l53
						l54:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							{
								position55, tokenIndex55, depth55 := position, tokenIndex, depth
								if buffer[position] != rune('"') {
									goto l55
								}
								position++
								goto l52
							l55:
								position, tokenIndex, depth = position55, tokenIndex55, depth55
							}
							if !matchDot() {
								goto l52
							}
						}
					l53:
						goto l51
					l52:
						position, tokenIndex, depth = position52, tokenIndex52, depth52
					}
					depth--
					add(rulePegText, position50)
				}
				if buffer[position] != rune('"') {
					goto l48
				}
				position++
				if !_rules[ruleAction5]() {
					goto l48
				}
				depth--
				add(ruledoubleQuotedString, position49)
			}
			return true
		l48:
			position, tokenIndex, depth = position48, tokenIndex48, depth48
			return false
		},
		/* 9 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']' Action6)> */
		func() bool {
			position56, tokenIndex56, depth56 := position, tokenIndex, depth
			{
				position57 := position
				depth++
				if buffer[position] != rune('[') {
					goto l56
				}
				position++
				{
					position58 := position
					depth++
					{
						position59, tokenIndex59, depth59 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l59
						}
						position++
						goto l60
					l59:
						position, tokenIndex, depth = position59, tokenIndex59, depth59
					}
				l60:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l56
					}
					position++
				l61:
					{
						position62, tokenIndex62, depth62 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l62
						}
						position++
						goto l61
					l62:
						position, tokenIndex, depth = position62, tokenIndex62, depth62
					}
					depth--
					add(rulePegText, position58)
				}
				if buffer[position] != rune(']') {
					goto l56
				}
				position++
				if !_rules[ruleAction6]() {
					goto l56
				}
				depth--
				add(rulejsonArrayAccess, position57)
			}
			return true
		l56:
			position, tokenIndex, depth = position56, tokenIndex56, depth56
			return false
		},
		/* 10 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']' Action7)> */
		func() bool {
			position63, tokenIndex63, depth63 := position, tokenIndex, depth
			{
				position64 := position
				depth++
				if buffer[position] != rune('[') {
					goto l63
				}
				position++
				{
					position65 := position
					depth++
					{
						position66, tokenIndex66, depth66 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l66
						}
						position++
						goto l67
					l66:
						position, tokenIndex, depth = position66, tokenIndex66, depth66
					}
				l67:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l63
					}
					position++
				l68:
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l69
						}
						position++
						goto l68
					l69:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
					}
					if buffer[position] != rune(':') {
						goto l63
					}
					position++
					{
						position70, tokenIndex70, depth70 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l70
						}
						position++
						goto l71
					l70:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
					}
				l71:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l63
					}
					position++
				l72:
					{
						position73, tokenIndex73, depth73 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l73
						}
						position++
						goto l72
					l73:
						position, tokenIndex, depth = position73, tokenIndex73, depth73
					}
					{
						position74, tokenIndex74, depth74 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l74
						}
						position++
						{
							position76, tokenIndex76, depth76 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l76
							}
							position++
							goto l77
						l76:
							position, tokenIndex, depth = position76, tokenIndex76, depth76
						}
					l77:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l74
						}
						position++
					l78:
						{
							position79, tokenIndex79, depth79 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l79
							}
							position++
							goto l78
						l79:
							position, tokenIndex, depth = position79, tokenIndex79, depth79
						}
						goto l75
					l74:
						position, tokenIndex, depth = position74, tokenIndex74, depth74
					}
				l75:
					depth--
					add(rulePegText, position65)
				}
				if buffer[position] != rune(']') {
					goto l63
				}
				position++
				if !_rules[ruleAction7]() {
					goto l63
				}
				depth--
				add(rulejsonArraySlice, position64)
			}
			return true
		l63:
			position, tokenIndex, depth = position63, tokenIndex63, depth63
			return false
		},
		/* 11 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']' Action8)> */
		func() bool {
			position80, tokenIndex80, depth80 := position, tokenIndex, depth
			{
				position81 := position
				depth++
				if buffer[position] != rune('[') {
					goto l80
				}
				position++
				{
					position82 := position
					depth++
					{
						position83, tokenIndex83, depth83 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l84
						}
						position++
						{
							position85, tokenIndex85, depth85 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l85
							}
							position++
							goto l86
						l85:
							position, tokenIndex, depth = position85, tokenIndex85, depth85
						}
					l86:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l84
						}
						position++
					l87:
						{
							position88, tokenIndex88, depth88 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l88
							}
							position++
							goto l87
						l88:
							position, tokenIndex, depth = position88, tokenIndex88, depth88
						}
						goto l83
					l84:
						position, tokenIndex, depth = position83, tokenIndex83, depth83
						{
							position89, tokenIndex89, depth89 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l89
							}
							position++
							goto l90
						l89:
							position, tokenIndex, depth = position89, tokenIndex89, depth89
						}
					l90:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l80
						}
						position++
					l91:
						{
							position92, tokenIndex92, depth92 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l92
							}
							position++
							goto l91
						l92:
							position, tokenIndex, depth = position92, tokenIndex92, depth92
						}
						if buffer[position] != rune(':') {
							goto l80
						}
						position++
					}
				l83:
					depth--
					add(rulePegText, position82)
				}
				if buffer[position] != rune(']') {
					goto l80
				}
				position++
				if !_rules[ruleAction8]() {
					goto l80
				}
				depth--
				add(rulejsonArrayPartialSlice, position81)
			}
			return true
		l80:
			position, tokenIndex, depth = position80, tokenIndex80, depth80
			return false
		},
		/* 12 jsonArrayFullSlice <- <('[' ':' ']' Action9)> */
		func() bool {
			position93, tokenIndex93, depth93 := position, tokenIndex, depth
			{
				position94 := position
				depth++
				if buffer[position] != rune('[') {
					goto l93
				}
				position++
				if buffer[position] != rune(':') {
					goto l93
				}
				position++
				if buffer[position] != rune(']') {
					goto l93
				}
				position++
				if !_rules[ruleAction9]() {
					goto l93
				}
				depth--
				add(rulejsonArrayFullSlice, position94)
			}
			return true
		l93:
			position, tokenIndex, depth = position93, tokenIndex93, depth93
			return false
		},
		/* 14 Action0 <- <{
		    p.addMapAccess(p.lastKey)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 15 Action1 <- <{
		    p.addMapAccess(p.lastKey)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 16 Action2 <- <{
		    p.addRecursiveAccess(p.lastKey)
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		nil,
		/* 18 Action3 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.lastKey = substr
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 19 Action4 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.lastKey = strings.Replace(substr, "''", "'", -1)
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 20 Action5 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.lastKey = strings.Replace(substr, "\"\"", "\"", -1)
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 21 Action6 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.addArrayAccess(substr)
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 22 Action7 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.addArraySlice(substr)
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 23 Action8 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.addArraySlice(substr)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 24 Action9 <- <{
		    p.addArraySlice("0:")
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
	}
	p.rules = _rules
}

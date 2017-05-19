package data

import (
	"strings"
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

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
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
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
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type jsonPeg struct {
	components []extractor
	lastKey    string

	Buffer string
	buffer []rune
	rules  [25]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *jsonPeg) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *jsonPeg) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
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
	p   *jsonPeg
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *jsonPeg) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *jsonPeg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
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
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *jsonPeg) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
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
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rulejsonPathHead]() {
					goto l0
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[rulejsonPathNonHead]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position4, tokenIndex4 := position, tokenIndex
					if !matchDot() {
						goto l4
					}
					goto l0
				l4:
					position, tokenIndex = position4, tokenIndex4
				}
				add(rulejsonPath, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 jsonPathHead <- <((jsonMapAccessString / jsonMapAccessBracket) Action0)> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
				{
					position7, tokenIndex7 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l8
					}
					goto l7
				l8:
					position, tokenIndex = position7, tokenIndex7
					if !_rules[rulejsonMapAccessBracket]() {
						goto l5
					}
				}
			l7:
				if !_rules[ruleAction0]() {
					goto l5
				}
				add(rulejsonPathHead, position6)
			}
			return true
		l5:
			position, tokenIndex = position5, tokenIndex5
			return false
		},
		/* 2 jsonPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position9, tokenIndex9 := position, tokenIndex
			{
				position10 := position
				{
					position11, tokenIndex11 := position, tokenIndex
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l12
					}
					goto l11
				l12:
					position, tokenIndex = position11, tokenIndex11
					if !_rules[rulejsonMapSingleLevel]() {
						goto l13
					}
					goto l11
				l13:
					position, tokenIndex = position11, tokenIndex11
					if !_rules[rulejsonArrayFullSlice]() {
						goto l14
					}
					goto l11
				l14:
					position, tokenIndex = position11, tokenIndex11
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l15
					}
					goto l11
				l15:
					position, tokenIndex = position11, tokenIndex11
					if !_rules[rulejsonArraySlice]() {
						goto l16
					}
					goto l11
				l16:
					position, tokenIndex = position11, tokenIndex11
					if !_rules[rulejsonArrayAccess]() {
						goto l9
					}
				}
			l11:
				add(rulejsonPathNonHead, position10)
			}
			return true
		l9:
			position, tokenIndex = position9, tokenIndex9
			return false
		},
		/* 3 jsonMapSingleLevel <- <((('.' jsonMapAccessString) / jsonMapAccessBracket) Action1)> */
		func() bool {
			position17, tokenIndex17 := position, tokenIndex
			{
				position18 := position
				{
					position19, tokenIndex19 := position, tokenIndex
					if buffer[position] != rune('.') {
						goto l20
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l20
					}
					goto l19
				l20:
					position, tokenIndex = position19, tokenIndex19
					if !_rules[rulejsonMapAccessBracket]() {
						goto l17
					}
				}
			l19:
				if !_rules[ruleAction1]() {
					goto l17
				}
				add(rulejsonMapSingleLevel, position18)
			}
			return true
		l17:
			position, tokenIndex = position17, tokenIndex17
			return false
		},
		/* 4 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket) Action2)> */
		func() bool {
			position21, tokenIndex21 := position, tokenIndex
			{
				position22 := position
				if buffer[position] != rune('.') {
					goto l21
				}
				position++
				if buffer[position] != rune('.') {
					goto l21
				}
				position++
				{
					position23, tokenIndex23 := position, tokenIndex
					if !_rules[rulejsonMapAccessString]() {
						goto l24
					}
					goto l23
				l24:
					position, tokenIndex = position23, tokenIndex23
					if !_rules[rulejsonMapAccessBracket]() {
						goto l21
					}
				}
			l23:
				if !_rules[ruleAction2]() {
					goto l21
				}
				add(rulejsonMapMultipleLevel, position22)
			}
			return true
		l21:
			position, tokenIndex = position21, tokenIndex21
			return false
		},
		/* 5 jsonMapAccessString <- <(<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> Action3)> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				{
					position27 := position
					{
						position28, tokenIndex28 := position, tokenIndex
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l29
						}
						position++
						goto l28
					l29:
						position, tokenIndex = position28, tokenIndex28
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l25
						}
						position++
					}
				l28:
				l30:
					{
						position31, tokenIndex31 := position, tokenIndex
						{
							position32, tokenIndex32 := position, tokenIndex
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l33
							}
							position++
							goto l32
						l33:
							position, tokenIndex = position32, tokenIndex32
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l34
							}
							position++
							goto l32
						l34:
							position, tokenIndex = position32, tokenIndex32
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l35
							}
							position++
							goto l32
						l35:
							position, tokenIndex = position32, tokenIndex32
							if buffer[position] != rune('_') {
								goto l31
							}
							position++
						}
					l32:
						goto l30
					l31:
						position, tokenIndex = position31, tokenIndex31
					}
					add(rulePegText, position27)
				}
				if !_rules[ruleAction3]() {
					goto l25
				}
				add(rulejsonMapAccessString, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 6 jsonMapAccessBracket <- <('[' (singleQuotedString / doubleQuotedString) ']')> */
		func() bool {
			position36, tokenIndex36 := position, tokenIndex
			{
				position37 := position
				if buffer[position] != rune('[') {
					goto l36
				}
				position++
				{
					position38, tokenIndex38 := position, tokenIndex
					if !_rules[rulesingleQuotedString]() {
						goto l39
					}
					goto l38
				l39:
					position, tokenIndex = position38, tokenIndex38
					if !_rules[ruledoubleQuotedString]() {
						goto l36
					}
				}
			l38:
				if buffer[position] != rune(']') {
					goto l36
				}
				position++
				add(rulejsonMapAccessBracket, position37)
			}
			return true
		l36:
			position, tokenIndex = position36, tokenIndex36
			return false
		},
		/* 7 singleQuotedString <- <('\'' <(('\'' '\'') / (!'\'' .))*> '\'' Action4)> */
		func() bool {
			position40, tokenIndex40 := position, tokenIndex
			{
				position41 := position
				if buffer[position] != rune('\'') {
					goto l40
				}
				position++
				{
					position42 := position
				l43:
					{
						position44, tokenIndex44 := position, tokenIndex
						{
							position45, tokenIndex45 := position, tokenIndex
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
							position, tokenIndex = position45, tokenIndex45
							{
								position47, tokenIndex47 := position, tokenIndex
								if buffer[position] != rune('\'') {
									goto l47
								}
								position++
								goto l44
							l47:
								position, tokenIndex = position47, tokenIndex47
							}
							if !matchDot() {
								goto l44
							}
						}
					l45:
						goto l43
					l44:
						position, tokenIndex = position44, tokenIndex44
					}
					add(rulePegText, position42)
				}
				if buffer[position] != rune('\'') {
					goto l40
				}
				position++
				if !_rules[ruleAction4]() {
					goto l40
				}
				add(rulesingleQuotedString, position41)
			}
			return true
		l40:
			position, tokenIndex = position40, tokenIndex40
			return false
		},
		/* 8 doubleQuotedString <- <('"' <(('"' '"') / (!'"' .))*> '"' Action5)> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				if buffer[position] != rune('"') {
					goto l48
				}
				position++
				{
					position50 := position
				l51:
					{
						position52, tokenIndex52 := position, tokenIndex
						{
							position53, tokenIndex53 := position, tokenIndex
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
							position, tokenIndex = position53, tokenIndex53
							{
								position55, tokenIndex55 := position, tokenIndex
								if buffer[position] != rune('"') {
									goto l55
								}
								position++
								goto l52
							l55:
								position, tokenIndex = position55, tokenIndex55
							}
							if !matchDot() {
								goto l52
							}
						}
					l53:
						goto l51
					l52:
						position, tokenIndex = position52, tokenIndex52
					}
					add(rulePegText, position50)
				}
				if buffer[position] != rune('"') {
					goto l48
				}
				position++
				if !_rules[ruleAction5]() {
					goto l48
				}
				add(ruledoubleQuotedString, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 9 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']' Action6)> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				if buffer[position] != rune('[') {
					goto l56
				}
				position++
				{
					position58 := position
					{
						position59, tokenIndex59 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l59
						}
						position++
						goto l60
					l59:
						position, tokenIndex = position59, tokenIndex59
					}
				l60:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l56
					}
					position++
				l61:
					{
						position62, tokenIndex62 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l62
						}
						position++
						goto l61
					l62:
						position, tokenIndex = position62, tokenIndex62
					}
					add(rulePegText, position58)
				}
				if buffer[position] != rune(']') {
					goto l56
				}
				position++
				if !_rules[ruleAction6]() {
					goto l56
				}
				add(rulejsonArrayAccess, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 10 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']' Action7)> */
		func() bool {
			position63, tokenIndex63 := position, tokenIndex
			{
				position64 := position
				if buffer[position] != rune('[') {
					goto l63
				}
				position++
				{
					position65 := position
					{
						position66, tokenIndex66 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l66
						}
						position++
						goto l67
					l66:
						position, tokenIndex = position66, tokenIndex66
					}
				l67:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l63
					}
					position++
				l68:
					{
						position69, tokenIndex69 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l69
						}
						position++
						goto l68
					l69:
						position, tokenIndex = position69, tokenIndex69
					}
					if buffer[position] != rune(':') {
						goto l63
					}
					position++
					{
						position70, tokenIndex70 := position, tokenIndex
						if buffer[position] != rune('-') {
							goto l70
						}
						position++
						goto l71
					l70:
						position, tokenIndex = position70, tokenIndex70
					}
				l71:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l63
					}
					position++
				l72:
					{
						position73, tokenIndex73 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l73
						}
						position++
						goto l72
					l73:
						position, tokenIndex = position73, tokenIndex73
					}
					{
						position74, tokenIndex74 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l74
						}
						position++
						{
							position76, tokenIndex76 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l76
							}
							position++
							goto l77
						l76:
							position, tokenIndex = position76, tokenIndex76
						}
					l77:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l74
						}
						position++
					l78:
						{
							position79, tokenIndex79 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l79
							}
							position++
							goto l78
						l79:
							position, tokenIndex = position79, tokenIndex79
						}
						goto l75
					l74:
						position, tokenIndex = position74, tokenIndex74
					}
				l75:
					add(rulePegText, position65)
				}
				if buffer[position] != rune(']') {
					goto l63
				}
				position++
				if !_rules[ruleAction7]() {
					goto l63
				}
				add(rulejsonArraySlice, position64)
			}
			return true
		l63:
			position, tokenIndex = position63, tokenIndex63
			return false
		},
		/* 11 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']' Action8)> */
		func() bool {
			position80, tokenIndex80 := position, tokenIndex
			{
				position81 := position
				if buffer[position] != rune('[') {
					goto l80
				}
				position++
				{
					position82 := position
					{
						position83, tokenIndex83 := position, tokenIndex
						if buffer[position] != rune(':') {
							goto l84
						}
						position++
						{
							position85, tokenIndex85 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l85
							}
							position++
							goto l86
						l85:
							position, tokenIndex = position85, tokenIndex85
						}
					l86:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l84
						}
						position++
					l87:
						{
							position88, tokenIndex88 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l88
							}
							position++
							goto l87
						l88:
							position, tokenIndex = position88, tokenIndex88
						}
						goto l83
					l84:
						position, tokenIndex = position83, tokenIndex83
						{
							position89, tokenIndex89 := position, tokenIndex
							if buffer[position] != rune('-') {
								goto l89
							}
							position++
							goto l90
						l89:
							position, tokenIndex = position89, tokenIndex89
						}
					l90:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l80
						}
						position++
					l91:
						{
							position92, tokenIndex92 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l92
							}
							position++
							goto l91
						l92:
							position, tokenIndex = position92, tokenIndex92
						}
						if buffer[position] != rune(':') {
							goto l80
						}
						position++
					}
				l83:
					add(rulePegText, position82)
				}
				if buffer[position] != rune(']') {
					goto l80
				}
				position++
				if !_rules[ruleAction8]() {
					goto l80
				}
				add(rulejsonArrayPartialSlice, position81)
			}
			return true
		l80:
			position, tokenIndex = position80, tokenIndex80
			return false
		},
		/* 12 jsonArrayFullSlice <- <('[' ':' ']' Action9)> */
		func() bool {
			position93, tokenIndex93 := position, tokenIndex
			{
				position94 := position
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
				add(rulejsonArrayFullSlice, position94)
			}
			return true
		l93:
			position, tokenIndex = position93, tokenIndex93
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

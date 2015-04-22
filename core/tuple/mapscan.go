package tuple

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

var reArrayPath = regexp.MustCompile(`^([^\[]+)?(\[[0-9]+\])?$`)

func split(s string) []string {
	i := 0
	a := []string{}
	t := ""
	rs := []rune(s)
	l := len(rs)
	for i < l {
		r := rs[i]
		switch r {
		case '\\':
			i++
			if i < l {
				t += string(rs[i])
			}
		case '.':
			if t != "" {
				a = append(a, t)
				t = ""
			}
		case '[':
			if i < l-1 {
				nr := rs[i+1]
				if nr == '"' || nr == '\'' {
					inbr := splitBracket(rs, i+2, nr)
					if inbr != "" {
						if t != "" {
							a = append(a, t)
						}
						a = append(a, inbr)
						t = ""
						i += 1 + len(inbr) + 2 // " + inner bracket + "]
						break
					}
				}
			}
			t += string(r)
		default:
			t += string(r)
		}
		i++
	}
	if t != "" {
		a = append(a, t)
	}
	return a
}

func splitBracket(rs []rune, i int, b rune) string {
	l := len(rs)
	t := ""
	for i < l {
		r := rs[i]
		if r == b {
			if i < l-1 {
				if rs[i+1] == ']' {
					return t
				}
			}
		} else {
			t += string(r)
		}
		i++
	}
	return ""
}

func scanMap(m Map, p string, t *Value) (err error) {
	if p == "" {
		return errors.New("empty key is not supported")
	}
	var v Value
	mm := m
	for _, token := range split(p) {
		sl := reArrayPath.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return errors.New("invalid path phrase")
		}
		ss := sl[0]
		if ss[1] != "" {
			mv := mm[ss[1]]
			if mv == nil {
				return errors.New("not found the key in map: " + ss[1])
			}
			v = mv
			if mv.Type() == TypeMap {
				mm, _ = mv.Map()
			}
		}
		// get array index number
		if ss[2] != "" {
			i64, err := strconv.ParseInt(ss[2][1:len(ss[2])-1], 10, 64)
			if err != nil {
				return errors.New("invalid array index number: " + token)
			}
			if i64 > math.MaxInt32 {
				return errors.New("overflow index number: " + token)
			}
			a, err := v.Array()
			if err != nil {
				return err
			}
			i := int(i64)
			if i >= len(a) {
				return errors.New("out of range access: " + token)
			}
			v = a[i]
		}
	}
	*t = v
	return nil
}

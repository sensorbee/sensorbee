package tuple

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile(`^([^\[]+)?(\[[0-9]+\])?$`)

func toError(v interface{}) error {
	if v != nil {
		if e, ok := v.(error); ok {
			return e
		}
		if e, ok := v.(string); ok {
			return errors.New(e)
		}
		return errors.New("Unknown error")
	}
	return nil
}

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

func scanMap(m Map, p string, t *Value) (err error) {
	if p == "" {
		return errors.New("empty key is not supported")
	}
	var v Value
	for _, token := range split(p) {
		sl := re.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return errors.New("invalid path phrase")
		}
		ss := sl[0]
		if ss[1] != "" {
			mv := m[ss[1]]
			if mv == nil {
				return errors.New("not found the key in map")
			}
			v = mv
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
			if v.Type() != TypeArray {
				return errors.New("invalid array path phrase: " + token)
			}
			i := int(i64)
			a, _ := v.Array()
			found := false
			for n, av := range a {
				if n == i {
					found = true
					v = av
					break
				}
			}
			if !found {
				return errors.New("out of range access: " + token)
			}
		}
	}
	*t = v
	return nil
}

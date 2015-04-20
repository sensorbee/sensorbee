package tuple

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

var reArray = regexp.MustCompile(`^([^\[]+)?(\[[0-9]+\])?$`)

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
		sl := reArray.FindAllStringSubmatch(token, -1)
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

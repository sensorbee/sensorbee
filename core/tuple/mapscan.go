package tuple

import (
	"errors"
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
		case '/':
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

func ScanMap(m Map, p string, t *Value) (err error) {
	if p == "" {
		return errors.New("empty path is not supported")
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
				return errors.New("invalid path phrase")
			}
			v = mv
		}
		// get array index number
		if ss[2] != "" {
			i, err := strconv.Atoi(ss[2][1 : len(ss[2])-1])
			if err != nil {
				return errors.New("invalid array index phrase: " + ss[2])
			}
			if v.Type() != TypeArray {
				return errors.New("invalid array path phrase: " + ss[1])
			}
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
				return errors.New("out of range access: " + ss[1] + ss[2])
			}
		}
	}
	*t = v
	return nil
}

package tuple

import (
	"errors"
	"regexp"
)

var re = regexp.MustCompile("^([^\\[]+)?(\\[[0-9]+\\])?$")

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

func ScanTree(v interface{}, p string, t interface{}) (err error) {
	defer func() {
		if err == nil {
			err = toError(recover())
		}
	}()
	if p == "" {
		return errors.New("invalid path")
	}
	//var ok bool
	for _, token := range split(p) {
		sl := re.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return errors.New("invalid path")
		}
	}
	return nil
}

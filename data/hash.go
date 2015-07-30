package data

import (
	"fmt"
	"hash/fnv"
	"io"
	"sort"
)

type HashValue uint64

func Hash(v Value) HashValue {
	h := fnv.New64a()
	updateHash(v, h)
	return HashValue(h.Sum64())
}

func HashEqual(v1 Value, v2 Value) bool {
	return v1.Type() == v2.Type() && Hash(v1) == Hash(v2)
}

func updateHash(v Value, h io.Writer) {
	switch v.Type() {
	case TypeBlob:
		// to decode: read digits until colon, then that many bytes
		b, _ := v.asBlob()
		h.Write([]byte(fmt.Sprintf("B%d:", len(b))))
		h.Write(b)
	case TypeBool:
		// to decode: read one digit
		b, _ := v.asBool()
		if b {
			h.Write([]byte("b1"))
		} else {
			h.Write([]byte("b0"))
		}
	case TypeFloat:
		// to decode: read until semi-colon. if any character other
		// than [-0-9] is contained, it is a float.
		f, _ := v.asFloat()
		// %v will automatically drop the ".0" for integer-like floats
		h.Write([]byte(fmt.Sprintf("n%v;", f)))
	case TypeInt:
		// to decode: read until semi-colon
		i, _ := v.asInt()
		h.Write([]byte(fmt.Sprintf("n%d;", i)))
	case TypeNull:
		h.Write([]byte("N"))
	case TypeString:
		// to decode: read digits until colon, then that many bytes
		s, _ := v.asString()
		h.Write([]byte(fmt.Sprintf("s%d:%s", len(s), s)))
	case TypeTimestamp:
		// to decode: read digits until colon, then that many bytes
		s := v.String()
		l := len(s)
		h.Write([]byte(fmt.Sprintf("t%d:%s", l-2, s[1:l-1])))
	case TypeArray:
		// to decode: read digits until colon, then recurse that
		// many times
		a, _ := v.asArray()
		h.Write([]byte(fmt.Sprintf("a%d:", len(a))))
		for _, item := range a {
			updateHash(item, h)
		}
	case TypeMap:
		// to decode: read digits until colon, then recurse twice
		// that many times to read key (string) and value
		m, _ := v.asMap()
		h.Write([]byte(fmt.Sprintf("m%d:", len(m))))
		keys := make(sort.StringSlice, 0, len(m))
		for key := range m {
			keys = append(keys, key)
		}
		keys.Sort()
		for _, key := range keys {
			updateHash(String(key), h)
			updateHash(m[key], h)
		}
	}
}

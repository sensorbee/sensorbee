package data

import (
	"bytes"
	"hash/fnv"
	"io"
	"sort"
	"strconv"
)

type HashValue uint64

func Hash(v Value) HashValue {
	h := fnv.New64a()
	updateHash(v, h)
	return HashValue(h.Sum64())
}

func Equal(v1 Value, v2 Value) bool {
	lType := v1.Type()
	rType := v2.Type()
	// cases in which we need a byte array comparison
	if lType == rType || // same type
		(lType == TypeFloat && rType == TypeInt) || // float vs. int
		(lType == TypeInt && rType == TypeFloat) { // int vs. float
		// compare based on the string representation
		// (this is exact, not probabilistic)
		var left, right bytes.Buffer
		updateHash(v1, &left)
		updateHash(v2, &right)
		return left.String() == right.String()
	}
	// if we arrive here, types are so different that the values
	// cannot possibly be equal
	return false
}

func updateHash(v Value, h io.Writer) {
	buffer := make([]byte, 1, 5)
	switch v.Type() {
	case TypeBlob:
		buffer[0] = 'B'
		// to decode: read digits until colon, then that many bytes
		b, _ := v.asBlob()
		buffer = strconv.AppendInt(buffer, int64(len(b)), 10)
		buffer = append(buffer, ':')
		h.Write(buffer)
		h.Write(b)
	case TypeBool:
		buffer[0] = 'b'
		// to decode: read one digit
		b, _ := v.asBool()
		if b {
			buffer = append(buffer, '1')
		} else {
			buffer = append(buffer, '0')
		}
		h.Write(buffer)
	case TypeFloat:
		buffer[0] = 'n'
		// to decode: read until semi-colon. if any character other
		// than [-0-9] is contained, it is a float.
		f, _ := v.asFloat()
		buffer = strconv.AppendFloat(buffer, f, 'g', -1, 64)
		buffer = append(buffer, ';')
		h.Write(buffer)
	case TypeInt:
		buffer[0] = 'n'
		// to decode: read until semi-colon
		i, _ := v.asInt()
		buffer = strconv.AppendInt(buffer, i, 10)
		buffer = append(buffer, ';')
		h.Write(buffer)
	case TypeNull:
		buffer[0] = 'N'
		h.Write(buffer)
	case TypeString:
		buffer[0] = 's'
		// to decode: read digits until colon, then that many bytes
		s, _ := v.asString()
		buffer = strconv.AppendInt(buffer, int64(len(s)), 10)
		buffer = append(buffer, ':')
		h.Write(buffer)
		h.Write([]byte(s))
	case TypeTimestamp:
		buffer[0] = 't'
		// to decode: read digits until colon, then that many bytes
		s := v.String()
		l := len(s)
		buffer = strconv.AppendInt(buffer, int64(l-2), 10)
		buffer = append(buffer, ':')
		h.Write(buffer)
		h.Write([]byte(s[1 : l-1]))
	case TypeArray:
		buffer[0] = 'a'
		// to decode: read digits until colon, then recurse that
		// many times
		a, _ := v.asArray()
		buffer = strconv.AppendInt(buffer, int64(len(a)), 10)
		buffer = append(buffer, ':')
		h.Write(buffer)
		for _, item := range a {
			updateHash(item, h)
		}
	case TypeMap:
		buffer[0] = 'm'
		// to decode: read digits until colon, then recurse twice
		// that many times to read key (string) and value
		m, _ := v.asMap()
		buffer = strconv.AppendInt(buffer, int64(len(m)), 10)
		buffer = append(buffer, ':')
		h.Write(buffer)
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

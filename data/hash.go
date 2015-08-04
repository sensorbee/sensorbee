package data

import (
	"bytes"
	"hash/fnv"
	"io"
	"math"
	"sort"
	"sync/atomic"
)

type HashValue uint64

func Hash(v Value) HashValue {
	h := fnv.New64a()
	buffer := make([]byte, 0, 16)
	updateHash(v, h, buffer)
	return HashValue(h.Sum64())
}

func Equal(v1 Value, v2 Value) bool {
	lType := v1.Type()
	rType := v2.Type()
	// cases in which we need a byte array comparison
	if lType == rType || // same type
		(lType == TypeFloat && rType == TypeInt) || // float vs. int
		(lType == TypeInt && rType == TypeFloat) { // int vs. float
	} else {
		// if we arrive here, types are so different that the values
		// cannot possibly be equal
		return false
	}

	switch lType {
	case TypeNull:
		// As long as Null is compared within a map or a array, Null = Null
		// should be true.
		return true

	case TypeBool:
		lhs, _ := v1.asBool()
		rhs, _ := v2.asBool()
		return lhs == rhs

	case TypeInt:
		lhs, _ := v1.asInt()
		if rType == TypeFloat {
			rhs, _ := v2.asFloat()
			return float64(lhs) == rhs
		}
		rhs, _ := v2.asInt()
		return lhs == rhs

	case TypeFloat:
		lhs, _ := v1.asFloat()
		if rType == TypeInt {
			rhs, _ := v2.asInt()
			return lhs == float64(rhs)
		}
		rhs, _ := v2.asFloat()
		return lhs == rhs // NaN == NaN is false

	case TypeString:
		lhs, _ := v1.asString()
		rhs, _ := v2.asString()
		return lhs == rhs

	case TypeBlob:
		lhs, _ := v1.asBlob()
		rhs, _ := v2.asBlob()
		return bytes.Equal(lhs, rhs)

	case TypeTimestamp:
		lhs, _ := v1.asTimestamp()
		rhs, _ := v2.asTimestamp()
		return lhs.Equal(rhs)

	case TypeArray:
		lhs, _ := v1.asArray()
		rhs, _ := v2.asArray()
		if len(lhs) != len(rhs) {
			return false
		}
		for i, l := range lhs {
			if !Equal(l, rhs[i]) {
				return false
			}
		}
		return true

	case TypeMap:
		lhs, _ := v1.asMap()
		rhs, _ := v2.asMap()
		if len(lhs) != len(rhs) {
			return false
		}
		for k, l := range lhs {
			r, ok := rhs[k]
			if !ok {
				return false
			}
			if !Equal(l, r) {
				return false
			}
		}
		return true

	default:
		// no such case, though
		return false
	}
}

func appendInt32(b []byte, t TypeID, i int32) []byte {
	return append(b, byte(t),
		byte(i&0xff),
		byte((i>>8)&0xff),
		byte((i>>16)&0xff),
		byte((i>>24)&0xff),
	)
}

func appendInt64(b []byte, t TypeID, i int64) []byte {
	return append(b, byte(t),
		byte(i&0xff),
		byte((i>>8)&0xff),
		byte((i>>16)&0xff),
		byte((i>>24)&0xff),
		byte((i>>32)&0xff),
		byte((i>>40)&0xff),
		byte((i>>48)&0xff),
		byte((i>>56)&0xff),
	)
}

var (
	nullHashCounter int64
)

func updateHash(v Value, h io.Writer, buffer []byte) []byte {
	switch v.Type() {
	case TypeNull:
		buffer = appendInt64(buffer, TypeNull, 0)
		h.Write(buffer)

	case TypeBool:
		b, _ := v.asBool()
		if b {
			buffer = append(buffer, byte(TypeBool), 1)
		} else {
			buffer = append(buffer, byte(TypeBool), 0)
		}
		h.Write(buffer)

	case TypeInt:
		i, _ := v.asInt()
		buffer = appendInt64(buffer, TypeInt, i)
		h.Write(buffer)

	case TypeFloat:
		f, _ := v.asFloat()
		if float64(int64(f)) == f {
			return updateHash(Int(f), h, buffer)
		}

		if math.IsNaN(f) {
			// NaN is processed as Null with a unique counter which results in
			// generating different hash values for each NaN.
			cnt := atomic.AddInt64(&nullHashCounter, 1)
			buffer = appendInt64(buffer, TypeNull, cnt)
		} else {
			buffer = appendInt64(buffer, TypeFloat, int64(math.Float64bits(f)))
		}
		h.Write(buffer)

	case TypeString:
		s, _ := v.asString()
		buffer = appendInt32(buffer, TypeString, int32(len(s)))
		h.Write(buffer)
		io.WriteString(h, s)

	case TypeBlob:
		b, _ := v.asBlob()
		buffer = appendInt32(buffer, TypeBlob, int32(len(b)))
		h.Write(buffer)
		h.Write(b)

	case TypeTimestamp:
		t, _ := v.asTimestamp()
		buffer = appendInt64(buffer, TypeTimestamp, t.Unix())
		// TODO: This TypeInt isn't necessary.
		buffer = appendInt32(buffer, TypeInt, int32(t.Nanosecond()/1000)) // Use microseconds
		h.Write(buffer)

	case TypeArray:
		a, _ := v.asArray()
		buffer = appendInt32(buffer, TypeArray, int32(len(a)))
		h.Write(buffer)
		for _, item := range a {
			buffer = updateHash(item, h, buffer[:0])
		}

	case TypeMap:
		m, _ := v.asMap()
		buffer = appendInt32(buffer, TypeMap, int32(len(m)))
		h.Write(buffer)

		// TODO: reduce this allocation
		keys := make(sort.StringSlice, 0, len(m))
		for key := range m {
			keys = append(keys, key)
		}
		keys.Sort()
		for _, key := range keys {
			buffer = updateHash(String(key), h, buffer[:0])
			buffer = updateHash(m[key], h, buffer[:0])
		}
	}
	return buffer
}

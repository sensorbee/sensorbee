package data

import (
	"hash/fnv"
	"io"
	"math"
	"sort"
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
		return Hash(v1) == Hash(v2)
	}
	// if we arrive here, types are so different that the values
	// cannot possibly be equal
	return false
}

func appendInt32(b []byte, i int32) []byte {
	return append(b,
		byte(i&0xff),
		byte((i>>8)&0xff),
		byte((i>>16)&0xff),
		byte((i>>24)&0xff),
	)
}

func appendInt64(b []byte, i int64) []byte {
	return append(b,
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

func updateHash(v Value, h io.Writer, buffer []byte) []byte {
	switch v.Type() {
	case TypeBlob:
		b, _ := v.asBlob()
		buffer = append(buffer, byte(TypeBlob))
		buffer = appendInt32(buffer, int32(len(b)))
		h.Write(buffer)
		h.Write(b)
	case TypeBool:
		b, _ := v.asBool()
		if b {
			buffer = append(buffer, byte(TypeBool), 1)
		} else {
			buffer = append(buffer, byte(TypeBool), 0)
		}
		h.Write(buffer)
	case TypeFloat:
		f, _ := v.asFloat()
		if float64(int64(f)) == f {
			return updateHash(Int(f), h, buffer)
		}
		buffer = append(buffer, byte(TypeFloat))
		buffer = appendInt64(buffer, int64(math.Float64bits(f)))
		h.Write(buffer)
	case TypeInt:
		i, _ := v.asInt()
		buffer = append(buffer, byte(TypeInt))
		buffer = appendInt64(buffer, i)
		h.Write(buffer)
	case TypeNull:
		buffer = append(buffer, byte(TypeNull))
		h.Write(buffer)
	case TypeString:
		s, _ := v.asString()
		buffer = append(buffer, byte(TypeString))
		buffer = appendInt32(buffer, int32(len(s)))
		h.Write(buffer)
		io.WriteString(h, s)
	case TypeTimestamp:
		t, _ := v.asTimestamp()
		buffer = append(buffer, byte(TypeTimestamp))
		buffer = appendInt64(buffer, t.Unix())
		buffer = appendInt32(buffer, int32(t.Nanosecond()/1000)) // Use microseconds
		h.Write(buffer)
	case TypeArray:
		a, _ := v.asArray()
		buffer = append(buffer, byte(TypeArray))
		buffer = appendInt32(buffer, int32(len(a)))
		h.Write(buffer)
		for _, item := range a {
			buffer = updateHash(item, h, buffer[:0])
		}
	case TypeMap:
		m, _ := v.asMap()
		buffer = append(buffer, byte(TypeMap))
		buffer = appendInt32(buffer, int32(len(m)))
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

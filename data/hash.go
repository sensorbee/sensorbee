package data

import (
	"bytes"
	"hash/fnv"
	"io"
	"math"
	"sync/atomic"
)

// HashValue is a type which hold hash values of Values.
type HashValue uint64

// Hash computes a hash value of a Value. A hash value of a Float having an
// integer value is computed as a Int's hash value so that Hash(Float(2.0))
// equals Hash(Int(2)). The hash value of Null is always the same. Hash values
// of NaNs always varies. For example, Hash(Float(math.NaN())) isn't equal to
// Hash(Float(math.NaN())). Therefore, if an array or a map has a NaN, the hash
// value changes everytime calling Hash function.
func Hash(v Value) HashValue {
	h := fnv.New64a()
	buffer := make([]byte, 0, 16)
	updateHash(v, h, buffer)
	return HashValue(h.Sum64())
}

// Equal tests equality of two Values. When comparing a Float and an Int, Int is
// implicitly converted to float so that they can be true when the Float has an
// integer value. For example, Equal(Float(2.0), Int(2)) is true.
//
// If either one is NaN, Equal always returns false. Equal(Null, Null) is true
// although this is inconsistent with the three-valued logic.
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
	i *= 16777619 // multiply fnv.prime32 due to the same reason as appendInt64
	return append(b, byte(t),
		byte(i&0xff),
		byte((i>>8)&0xff),
		byte((i>>16)&0xff),
		byte((i>>24)&0xff),
	)
}

func appendInt64(b []byte, t TypeID, i int64) []byte {
	// Because FNV-64a doesn't seem to work well with small numbers,
	// fnv.prime64 is manually multiplied beforehand.
	i *= 1099511628211
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
			buffer = append(buffer, byte(TypeBool), 0xaa)
		} else {
			buffer = append(buffer, byte(TypeBool), 0x55)
		}
		// 0xaa and 0x55 is to make better distribution of hash values.
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

		// In terms of reducing hash value collisions, computing the hash value
		// over sorted key-value pairs is the best. However, some Maps such as
		// feature vectors of machine learning sometimes have over 10,000
		// elements. Because Hash is called on each tuple, this method could be
		// a serious bottle-neck.
		//
		// The following implementation computes consistent hash values very
		// quickly by sacrificing the possibility of collisions. Instead of
		// sorting keys, it sums up all hash values computed from each key
		// value pair. Therefore, when numbers of elements in two Maps are
		// the same, possibility of collisions will increase. However, the hash
		// value is 64bit and SensorBee usually process less than 1B tuples
		// at once, the possibility is considered sufficiently low.

		var upper uint32
		var lower uint64

		subHash := fnv.New64a() // TODO: reduce this allocation
		for k, v := range m {
			subHash.Reset()

			// Because values usually vary more than keys, hash values of values
			// should be computed first to make better distribution of hash
			// values.
			buffer = updateHash(v, subHash, buffer[:0])
			buffer = updateHash(String(k), subHash, buffer[:0])
			sh := subHash.Sum64()
			if sh+lower < sh|lower { // carried
				upper++
			}
			lower += sh
		}

		buffer = appendInt32(buffer[:0], TypeMap, int32(len(m)))
		buffer = appendInt32(buffer, TypeMap, int32(upper))
		buffer = appendInt64(buffer, TypeMap, int64(lower))
		h.Write(buffer)
	}
	return buffer
}

package tuple

import (
	"github.com/mattn/go-scan"
	"time"
)

type Value interface {
	Type() TypeID
	ToInt(path String) Int
	ToFloat(path String) Float
	ToString(path String) String
	ToTimestamp(path String) Timestamp
	Array() Array
	ToArray(path String) Array
	Map() Map
	ToMap(path String) Map
}

const (
	TypeUnknown TypeID = iota
	TypeInt
	TypeFloat
	TypeString
	TypeBlob
	TypeTimesamp
	TypeArray
	TypeMap
)

func (t TypeID) String() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	case TypeBlob:
		return "blob"
	case TypeTimesamp:
		return "timestamp"
	case TypeArray:
		return "array"
	case TypeMap:
		return "map"
	default:
		return "unknown"
	}
}

type TypeID int
type Int int64
type Float float64
type String string
type Blob []byte
type Timestamp time.Time
type Array []Value
type Map map[string]Value

func (i Int) Type() TypeID {
	return TypeInt
}

func (f Float) Type() TypeID {
	return TypeFloat
}

func (s String) Type() TypeID {
	return TypeString
}

func (b Blob) Type() TypeID {
	return TypeBlob
}

func (t Timestamp) Type() TypeID {
	return TypeTimesamp
}

func (a Array) Type() TypeID {
	return TypeArray
}

func (m Map) Type() TypeID {
	return TypeMap
}

func (m Map) ToInt(path string) (Int, error) {
	var i Int
	err := scan.ScanTree(m.toMapInterface(), path, &i)
	return i, err
}

func (m Map) ToFloat(path string) (Float, error) {
	var f Float
	err := scan.ScanTree(m.toMapInterface(), path, &f)
	return f, err
}

func (m Map) ToString(path string) (String, error) {
	var s String
	err := scan.ScanTree(m.toMapInterface(), path, &s)
	return s, err
}

func (m Map) ToBlob(path string) (Blob, error) {
	var b Blob
	err := scan.ScanTree(m.toMapInterface(), path, &b)
	return b, err
}

func (a Array) Array() Array {
	return a
}

func (m Map) ToArray(path string) (Array, error) {
	var a Array
	err := scan.ScanTree(m.toMapInterface(), path, &a)
	return a, err
}

func (m Map) Map() Map {
	return m
}

func (m Map) ToMap(path string) (Map, error) {
	var mm Map
	err := scan.ScanTree(m.toMapInterface(), path, &mm)
	return mm, err
}

func (a Array) toArrayInterface() []interface{} {
	t := []interface{}{}
	for _, v := range a {
		var e interface{}
		switch v.Type() {
		case TypeArray:
			e = v.Array().toArrayInterface()
		case TypeMap:
			e = v.Map().toMapInterface()
		default:
			e = v
		}
		t = append(t, e)
	}
	return t
}

func (m Map) toMapInterface() map[string]interface{} {
	t := map[string]interface{}{}
	for k, v := range m {
		switch v.Type() {
		case TypeArray:
			t[k] = v.Array().toArrayInterface()
		case TypeMap:
			t[k] = v.Map().toMapInterface()
		default:
			t[k] = v
		}
	}
	return t
}

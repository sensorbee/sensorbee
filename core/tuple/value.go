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
	ToArray(path String) Array
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

func (v Value) ToInt(path string) (Value, error) {
	var i Int
	err := scan.ScanTree(v.toMapInterface(), path, &i)
	return i, err
}

func (v *Value) ToFloat(path string) (Value, error) {
	var f Float
	err := scan.ScanTree(v.toMapInterface(), path, &f)
	return f, err
}

func (v *Value) ToString(path string) (Value, error) {
	var s String
	err := scan.ScanTree(v.toMapInterface(), path, &s)
	return s, err
}

func (v *Value) ToBlob(path string) (Value, error) {
	var b Blob
	err := scan.ScanTree(v.toMapInterface(), path, &b)
	return b, err
}

func (v *Value) ToArray(path string) (Value, error) {
	var a Array
	err := scan.ScanTree(v.toMapInterface(), path, &a)
	return a, err
}

func (v *Value) ToMap(path string) (Value, error) {
	var m Map
	err := scan.ScanTree(v.toMapInterface(), path, &m)
	return m, err
}

func (a Map) toArrayInterface() []interface{} {
	t := []interface{}{}
	for _, v := range a {
		var e interface{}
		switch v.Type() {
		case TypeArray:
			e = v.ToArray().toArrayInterface()
		case TypeMap:
			e = v.ToMap().toMapInterface()
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

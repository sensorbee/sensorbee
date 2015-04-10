package tuple

import (
	"errors"
	"fmt"
	"github.com/mattn/go-scan"
	"time"
)

type Value interface {
	Type() TypeID
	Bool() (Bool, error)
	Int() (Int, error)
	Float() (Float, error)
	String() (String, error)
	Blob() (Blob, error)
	Timestamp() (Timestamp, error)
	Array() (Array, error)
	Map() (Map, error)
}

// TODO: Provide NewMap(map[string]interface{}) Map

func castError(from TypeID, to TypeID) error {
	return errors.New(fmt.Sprintf("unsupported cast %v from %v", to.String(), from.String()))
}

type TypeID int

const (
	TypeUnknown TypeID = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeBlob
	TypeTimestamp
	TypeArray
	TypeMap
)

func (t TypeID) String() string {
	switch t {
	case TypeBool:
		return "bool"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	case TypeBlob:
		return "blob"
	case TypeTimestamp:
		return "timestamp"
	case TypeArray:
		return "array"
	case TypeMap:
		return "map"
	default:
		return "unknown"
	}
}

type Bool bool

func (b Bool) Type() TypeID {
	return TypeBool
}

func (b Bool) Bool() (Bool, error) {
	return b, nil
}

func (b Bool) Int() (Int, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Bool) Float() (Float, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Bool) String() (String, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Bool) Blob() (Blob, error) {
	return nil, castError(b.Type(), TypeBlob)
}

func (b Bool) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(b.Type(), TypeTimestamp)
}

func (b Bool) Array() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Bool) Map() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

type Int int64

func (i Int) Type() TypeID {
	return TypeInt
}

func (i Int) Bool() (Bool, error) {
	return false, castError(i.Type(), TypeBool)
}

func (i Int) Int() (Int, error) {
	return i, nil
}

func (i Int) Float() (Float, error) {
	return Float(i), nil
}

func (i Int) String() (String, error) {
	return String(fmt.Sprint(i)), nil
}

func (i Int) Blob() (Blob, error) {
	return nil, castError(i.Type(), TypeBlob)
}

func (i Int) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(i.Type(), TypeTimestamp)
}

func (i Int) Array() (Array, error) {
	return nil, castError(i.Type(), TypeArray)
}

func (i Int) Map() (Map, error) {
	return nil, castError(i.Type(), TypeMap)
}

type Float float64

func (f Float) Type() TypeID {
	return TypeFloat
}

func (f Float) Bool() (Bool, error) {
	return false, castError(f.Type(), TypeBool)
}

func (f Float) Int() (Int, error) {
	return Int(f), nil
}

func (f Float) Float() (Float, error) {
	return f, nil
}

func (f Float) String() (String, error) {
	return String(fmt.Sprint(f)), nil
}

func (f Float) Blob() (Blob, error) {
	return nil, castError(f.Type(), TypeFloat)
}

func (f Float) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(f.Type(), TypeTimestamp)
}

func (f Float) Array() (Array, error) {
	return nil, castError(f.Type(), TypeArray)
}

func (f Float) Map() (Map, error) {
	return nil, castError(f.Type(), TypeMap)
}

type String string

func (s String) Type() TypeID {
	return TypeString
}

func (s String) Bool() (Bool, error) {
	return false, castError(s.Type(), TypeBool)
}

func (s String) Int() (Int, error) {
	return 0, castError(s.Type(), TypeInt)
}

func (s String) Float() (Float, error) {
	return 0, castError(s.Type(), TypeFloat)
}

func (s String) String() (String, error) {
	return s, nil
}

func (s String) Blob() (Blob, error) {
	return nil, castError(s.Type(), TypeBlob)
}

func (s String) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(s.Type(), TypeTimestamp)
}

func (s String) Array() (Array, error) {
	return nil, castError(s.Type(), TypeArray)
}

func (s String) Map() (Map, error) {
	return nil, castError(s.Type(), TypeMap)
}

type Blob []byte

func (b Blob) Type() TypeID {
	return TypeBlob
}

func (b Blob) Bool() (Bool, error) {
	return false, castError(b.Type(), TypeBool)
}

func (b Blob) Int() (Int, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Blob) Float() (Float, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Blob) String() (String, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Blob) Blob() (Blob, error) {
	return b, nil
}

func (b Blob) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(b.Type(), TypeTimestamp)
}

func (b Blob) Array() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Blob) Map() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

type Timestamp time.Time

func (t Timestamp) Type() TypeID {
	return TypeTimestamp
}

func (t Timestamp) Bool() (Bool, error) {
	return false, castError(t.Type(), TypeBool)
}

func (t Timestamp) Int() (Int, error) {
	return 0, castError(t.Type(), TypeInt)
}

func (t Timestamp) Float() (Float, error) {
	return 0, castError(t.Type(), TypeFloat)
}

func (t Timestamp) String() (String, error) {
	return "", castError(t.Type(), TypeString)
}

func (t Timestamp) Blob() (Blob, error) {
	return nil, castError(t.Type(), TypeBlob)
}

func (t Timestamp) Timestamp() (Timestamp, error) {
	return t, nil
}

func (t Timestamp) Array() (Array, error) {
	return nil, castError(t.Type(), TypeArray)
}

func (t Timestamp) Map() (Map, error) {
	return nil, castError(t.Type(), TypeMap)
}

type Array []Value

func (a Array) Type() TypeID {
	return TypeArray
}

func (a Array) Bool() (Bool, error) {
	return false, castError(a.Type(), TypeBool)
}

func (a Array) Int() (Int, error) {
	return 0, castError(a.Type(), TypeInt)
}

func (a Array) Float() (Float, error) {
	return 0, castError(a.Type(), TypeFloat)
}

func (a Array) String() (String, error) {
	return "", castError(a.Type(), TypeString)
}

func (a Array) Blob() (Blob, error) {
	return nil, castError(a.Type(), TypeBlob)
}

func (a Array) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(a.Type(), TypeTimestamp)
}

func (a Array) Array() (Array, error) {
	return a, nil
}

func (a Array) Map() (Map, error) {
	return nil, castError(a.Type(), TypeMap)
}

func (a Array) toArrayInterface() []interface{} {
	t := make([]interface{}, len(a))
	for idx, value := range a {
		var element interface{}
		switch value.Type() {
		case TypeArray:
			a, _ := value.Array()
			element = a.toArrayInterface()
		case TypeMap:
			m, _ := value.Map()
			element = m.toMapInterface()
		default:
			element = value
		}
		t[idx] = element
	}
	return t
}

type Map map[string]Value

func (m Map) Type() TypeID {
	return TypeMap
}

func (m Map) Bool() (Bool, error) {
	return false, castError(m.Type(), TypeBool)
}

func (m Map) Int() (Int, error) {
	return 0, castError(m.Type(), TypeInt)
}

func (m Map) Float() (Float, error) {
	return 0, castError(m.Type(), TypeFloat)
}

func (m Map) String() (String, error) {
	return "", castError(m.Type(), TypeString)
}

func (m Map) Blob() (Blob, error) {
	return nil, castError(m.Type(), TypeBlob)
}

func (m Map) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(m.Type(), TypeTimestamp)
}

func (m Map) Array() (Array, error) {
	return nil, castError(m.Type(), TypeArray)
}

func (m Map) Map() (Map, error) {
	return m, nil
}

func (m Map) Get(path string) (Value, error) {
	// TODO: support json path manually
	var v Value
	err := scan.ScanTree(m.toMapInterface(), path, &v)
	return v, err
}

// toMapInterface converts Map to map[string]interface{}.
// This is only for go-scan.
func (m Map) toMapInterface() map[string]interface{} {
	t := make(map[string]interface{}, len(m))
	for k, v := range m {
		switch v.Type() {
		case TypeArray:
			a, _ := v.Array()
			t[k] = a.toArrayInterface()
		case TypeMap:
			m, _ := v.Map()
			t[k] = m.toMapInterface()
		default:
			t[k] = v
		}
	}
	return t
}

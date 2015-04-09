package tuple

import (
	"errors"
	"fmt"
	"github.com/mattn/go-scan"
	"time"
)

type Value interface {
	Type() TypeID
	Int() (Int, error)
	Float() (Float, error)
	String() (String, error)
	Blob() (Blob, error)
	Timestamp() (Timestamp, error)
	Array() (Array, error)
	Map() (Map, error)
}

// TODO: Provide NewMap(map[string]interface{}) Map

func Error(from TypeID, to TypeID) error {
	return errors.New(fmt.Sprintf("unsupported cast %v from %v", to.String(), from.String()))
}

type TypeID int

const (
	TypeUnknown TypeID = iota
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

type Int int64

func (i Int) Type() TypeID {
	return TypeInt
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
	return nil, Error(i.Type(), TypeBlob)
}

func (i Int) Timestamp() (Timestamp, error) {
	return Timestamp{}, Error(i.Type(), TypeTimestamp)
}

func (i Int) Array() (Array, error) {
	return nil, Error(i.Type(), TypeArray)
}

func (i Int) Map() (Map, error) {
	return nil, Error(i.Type(), TypeMap)
}

type Float float64

func (f Float) Type() TypeID {
	return TypeFloat
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
	return nil, Error(f.Type(), TypeFloat)
}

func (f Float) Timestamp() (Timestamp, error) {
	return Timestamp{}, Error(f.Type(), TypeTimestamp)
}

func (f Float) Array() (Array, error) {
	return nil, Error(f.Type(), TypeArray)
}

func (f Float) Map() (Map, error) {
	return nil, Error(f.Type(), TypeMap)
}

type String string

func (s String) Type() TypeID {
	return TypeString
}

func (s String) Int() (Int, error) {
	return 0, Error(s.Type(), TypeInt)
}

func (s String) Float() (Float, error) {
	return 0, Error(s.Type(), TypeFloat)
}

func (s String) String() (String, error) {
	return s, nil
}

func (s String) Blob() (Blob, error) {
	return nil, Error(s.Type(), TypeBlob)
}

func (s String) Timestamp() (Timestamp, error) {
	return Timestamp{}, Error(s.Type(), TypeTimestamp)
}

func (s String) Array() (Array, error) {
	return nil, Error(s.Type(), TypeArray)
}

func (s String) Map() (Map, error) {
	return nil, Error(s.Type(), TypeMap)
}

type Blob []byte

func (b Blob) Type() TypeID {
	return TypeBlob
}

func (b Blob) Int() (Int, error) {
	return 0, Error(b.Type(), TypeInt)
}

func (b Blob) Float() (Float, error) {
	return 0, Error(b.Type(), TypeFloat)
}

func (b Blob) String() (String, error) {
	return "", Error(b.Type(), TypeString)
}

func (b Blob) Blob() (Blob, error) {
	return b, nil
}

func (b Blob) Timestamp() (Timestamp, error) {
	return Timestamp{}, Error(b.Type(), TypeTimestamp)
}

func (b Blob) Array() (Array, error) {
	return nil, Error(b.Type(), TypeArray)
}

func (b Blob) Map() (Map, error) {
	return nil, Error(b.Type(), TypeMap)
}

type Timestamp time.Time

func (t Timestamp) Type() TypeID {
	return TypeTimestamp
}

func (t Timestamp) Int() (Int, error) {
	return 0, Error(t.Type(), TypeInt)
}

func (t Timestamp) Float() (Float, error) {
	return 0, Error(t.Type(), TypeFloat)
}

func (t Timestamp) String() (String, error) {
	return "", Error(t.Type(), TypeString)
}

func (t Timestamp) Blob() (Blob, error) {
	return nil, Error(t.Type(), TypeBlob)
}

func (t Timestamp) Timestamp() (Timestamp, error) {
	return t, nil
}

func (t Timestamp) Array() (Array, error) {
	return nil, Error(t.Type(), TypeArray)
}

func (t Timestamp) Map() (Map, error) {
	return nil, Error(t.Type(), TypeMap)
}

type Array []Value

func (a Array) Type() TypeID {
	return TypeArray
}

func (a Array) Int() (Int, error) {
	return 0, Error(a.Type(), TypeInt)
}

func (a Array) Float() (Float, error) {
	return 0, Error(a.Type(), TypeFloat)
}

func (a Array) String() (String, error) {
	return "", Error(a.Type(), TypeString)
}

func (a Array) Blob() (Blob, error) {
	return nil, Error(a.Type(), TypeBlob)
}

func (a Array) Timestamp() (Timestamp, error) {
	return Timestamp{}, Error(a.Type(), TypeTimestamp)
}

func (a Array) Array() (Array, error) {
	return a, nil
}

func (a Array) Map() (Map, error) {
	return nil, Error(a.Type(), TypeMap)
}

func (a Array) toArrayInterface() []interface{} {
	t := []interface{}{}
	for _, v := range a {
		var e interface{}
		switch v.Type() {
		case TypeArray:
			a, _ := v.Array()
			e = a.toArrayInterface()
		case TypeMap:
			m, _ := v.Map()
			e = m.toMapInterface()
		default:
			e = v
		}
		t = append(t, e)
	}
	return t
}

type Map map[string]Value

func (m Map) Type() TypeID {
	return TypeMap
}

func (m Map) Int() (Int, error) {
	return 0, Error(m.Type(), TypeInt)
}

func (m Map) Float() (Float, error) {
	return 0, Error(m.Type(), TypeFloat)
}

func (m Map) String() (String, error) {
	return "", Error(m.Type(), TypeString)
}

func (m Map) Blob() (Blob, error) {
	return nil, Error(m.Type(), TypeBlob)
}

func (m Map) Timestamp() (Timestamp, error) {
	return Timestamp{}, Error(m.Type(), TypeTimestamp)
}

func (m Map) Array() (Array, error) {
	return nil, Error(m.Type(), TypeArray)
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
	t := map[string]interface{}{}
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

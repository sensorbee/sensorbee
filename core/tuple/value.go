package tuple

import (
	"fmt"
	"github.com/mattn/go-scan"
	"time"
)

type Value interface {
	Type() TypeID
	Int() Int
	Float() Float
	String() String
	Blob() Blob
	Timestamp() Timestamp
	Array() Array
	Map() Map
}

// TODO: Provide NewMap(map[string]interface{}) Map
// TODO: Need Implemet Test Code

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

func (i Int) Int() Int {
	return i
}

func (i Int) Float() Float {
	return Float(i)
}

func (i Int) String() String {
	return String(fmt.Sprint(i))
}

func (i Int) Blob() Blob {
	// TODO: This method should return an error instead of panic
	panic("unsupported conversion")
}

func (i Int) Timestamp() Timestamp {
	panic("unsupported conversion")
}

func (i Int) Array() Array {
	panic("unsupported conversion")
}

func (i Int) Map() Map {
	panic("unsupported conversion")
}

type Float float64

func (f Float) Type() TypeID {
	return TypeFloat
}

func (f Float) Int() Int {
	return Int(f)
}

func (f Float) Float() Float {
	return f
}

func (f Float) String() String {
	return String(fmt.Sprint(f))
}

func (f Float) Blob() Blob {
	panic("unsupported conversion")
}

func (f Float) Timestamp() Timestamp {
	panic("unsupported conversion")
}

func (f Float) Array() Array {
	panic("unsupported conversion")
}

func (f Float) Map() Map {
	panic("unsupported conversion")
}

type String string

func (s String) Type() TypeID {
	return TypeString
}

func (s String) Int() Int {
	panic("unsupported conversion")
}

func (s String) Float() Float {
	panic("unsupported conversion")
}

func (s String) String() String {
	return s
}

func (s String) Blob() Blob {
	panic("unsupported conversion")
}

func (s String) Timestamp() Timestamp {
	panic("unsupported conversion")
}

func (s String) Array() Array {
	panic("unsupported conversion")
}

func (s String) Map() Map {
	panic("unsupported conversion")
}

type Blob []byte

func (b Blob) Type() TypeID {
	return TypeBlob
}

func (b Blob) Int() Int {
	panic("unsupported conversion")
}

func (b Blob) Float() Float {
	panic("unsupported conversion")
}

func (b Blob) String() String {
	panic("unsupported conversion")
}

func (b Blob) Blob() Blob {
	return b
}

func (b Blob) Timestamp() Timestamp {
	panic("unsupported conversion")
}

func (b Blob) Array() Array {
	panic("unsupported conversion")
}

func (b Blob) Map() Map {
	panic("unsupported conversion")
}

type Timestamp time.Time

func (t Timestamp) Type() TypeID {
	return TypeTimestamp
}

func (t Timestamp) Int() Int {
	panic("unsupported conversion")
}

func (t Timestamp) Float() Float {
	panic("unsupported conversion")
}

func (t Timestamp) String() String {
	panic("unsupported conversion")
}

func (t Timestamp) Blob() Blob {
	panic("unsupported conversion")
}

func (t Timestamp) Timestamp() Timestamp {
	return t
}

func (t Timestamp) Array() Array {
	panic("unsupported conversion")
}

func (t Timestamp) Map() Map {
	panic("unsupported conversion")
}

type Array []Value

func (a Array) Type() TypeID {
	return TypeArray
}

func (a Array) Int() Int {
	panic("unsupported conversion")
}

func (a Array) Float() Float {
	panic("unsupported conversion")
}

func (a Array) String() String {
	panic("unsupported conversion")
}

func (a Array) Blob() Blob {
	panic("unsupported conversion")
}

func (a Array) Timestamp() Timestamp {
	panic("unsupported conversion")
}

func (a Array) Array() Array {
	return a
}

func (a Array) Map() Map {
	panic("unsupported conversion")
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

type Map map[string]Value

func (m Map) Type() TypeID {
	return TypeMap
}

func (m Map) Int() Int {
	panic("unsupported conversion")
}

func (m Map) Float() Float {
	panic("unsupported conversion")
}

func (m Map) String() String {
	panic("unsupported conversion")
}

func (m Map) Blob() Blob {
	panic("unsupported conversion")
}

func (m Map) Timestamp() Timestamp {
	panic("unsupported conversion")
}

func (m Map) Array() Array {
	panic("unsupported conversion")
}

func (m Map) Map() Map {
	return m
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
			t[k] = v.Array().toArrayInterface()
		case TypeMap:
			t[k] = v.Map().toMapInterface()
		default:
			t[k] = v
		}
	}
	return t
}

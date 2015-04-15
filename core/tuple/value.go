package tuple

import (
	"errors"
	"fmt"
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
	clone() Value
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

package tuple

import (
	"errors"
	"fmt"
	"time"
)

// Value is the generic interface for all data that can be stored
// inside a Tuple. Since we assume the data not to conform to any
// schema, data can have any shape and it can also change within a
// stream from one Tuple to the next. Therefore we need to be
// careful with respect to type conversions. A Value obtained, e.g.,
// by Map.Get should always be converted using the appropriate method
// and error checking must be done.
//
// Example:
//  i, err := val.AsInt()
//  if err != nil { ... }
type Value interface {
	Type() TypeID
	AsBool() (bool, error)
	AsInt() (int64, error)
	AsFloat() (float64, error)
	AsString() (string, error)
	AsBlob() ([]byte, error)
	AsTimestamp() (time.Time, error)
	AsArray() (Array, error)
	AsMap() (Map, error)
	clone() Value
}

// TODO: Provide NewMap(map[string]interface{}) Map

func castError(from TypeID, to TypeID) error {
	return errors.New(fmt.Sprintf("unsupported cast %v from %v", to.String(), from.String()))
}

type TypeID int

const (
	TypeUnknown TypeID = iota
	TypeNull
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
	case TypeNull:
		return "null"
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

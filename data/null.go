package data

import (
	"time"
)

type Null struct{}

func (n Null) Type() TypeID {
	return TypeNull
}

func (n Null) asBool() (bool, error) {
	return false, castError(n.Type(), TypeBool)
}

func (n Null) asInt() (int64, error) {
	return 0, castError(n.Type(), TypeInt)
}

func (n Null) asFloat() (float64, error) {
	return 0, castError(n.Type(), TypeFloat)
}

func (n Null) asString() (string, error) {
	return "", castError(n.Type(), TypeString)
}

func (n Null) asBlob() ([]byte, error) {
	return nil, castError(n.Type(), TypeBlob)
}

func (n Null) asTimestamp() (time.Time, error) {
	return time.Time{}, castError(n.Type(), TypeTimestamp)
}

func (n Null) asArray() (Array, error) {
	return nil, castError(n.Type(), TypeArray)
}

func (n Null) asMap() (Map, error) {
	return nil, castError(n.Type(), TypeMap)
}

func (n Null) clone() Value {
	return Null{}
}

func (n Null) MarshalJSON() ([]byte, error) {
	// the JSON serialization is defined via the String()
	// return value as defined below
	return []byte(n.String()), nil
}

func (n Null) String() string {
	return "null"
}

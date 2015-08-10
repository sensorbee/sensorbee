package data

import (
	"time"
)

// Null corresponds to null in JSON. It can be assigned to Value interface.
// Null is provided for Null Object pattern and it should always be used
// instead of nil.
type Null struct{}

// Type returns TypeID of Null. It's always TypeNull.
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

// MarshalJSON marshals Null to JSON.
func (n Null) MarshalJSON() ([]byte, error) {
	// the JSON serialization is defined via the String()
	// return value as defined below
	return []byte(n.String()), nil
}

// String returns JSON representation of a Null.
func (n Null) String() string {
	return "null"
}

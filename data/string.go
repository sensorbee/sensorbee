package data

import (
	"encoding/json"
	"fmt"
	"time"
)

// String is a string. It can be assigned to Value interface.
type String string

// Type returns TypeID of String. It's always TypeString.
func (s String) Type() TypeID {
	return TypeString
}

func (s String) asBool() (bool, error) {
	return false, castError(s.Type(), TypeBool)
}

func (s String) asInt() (int64, error) {
	return 0, castError(s.Type(), TypeInt)
}

func (s String) asFloat() (float64, error) {
	return 0, castError(s.Type(), TypeFloat)
}

func (s String) asString() (string, error) {
	return string(s), nil
}

func (s String) asBlob() ([]byte, error) {
	return nil, castError(s.Type(), TypeBlob)
}

func (s String) asTimestamp() (time.Time, error) {
	return time.Time{}, castError(s.Type(), TypeTimestamp)
}

func (s String) asArray() (Array, error) {
	return nil, castError(s.Type(), TypeArray)
}

func (s String) asMap() (Map, error) {
	return nil, castError(s.Type(), TypeMap)
}

func (s String) clone() Value {
	return String(s)
}

// String returns JSON representation of a String. A string "a" will be marshaled
// as `"a"` (double quotes are included), not `a`. To obtain a plain string
// without double quotes, use ToString function.
func (s String) String() string {
	// the String return value is defined via the
	// default JSON serialization
	bytes, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("(unserializable string: %v)", err)
	}
	return string(bytes)
}

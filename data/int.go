package data

import (
	"encoding/json"
	"fmt"
	"time"
)

// Int is an integer. It can be assigned to Value interface. A value more than
// 2^53 - 1 cannot exactly be marshaled to JSON because some languages like
// JavaScript or Lua only has float as a numeric type.
type Int int64

// Type returns TypeID of Int. It's always TypeInt.
func (i Int) Type() TypeID {
	return TypeInt
}

func (i Int) asBool() (bool, error) {
	return false, castError(i.Type(), TypeBool)
}

func (i Int) asInt() (int64, error) {
	return int64(i), nil
}

func (i Int) asFloat() (float64, error) {
	return 0, castError(i.Type(), TypeFloat)
}

func (i Int) asString() (string, error) {
	return "", castError(i.Type(), TypeString)
}

func (i Int) asBlob() ([]byte, error) {
	return nil, castError(i.Type(), TypeBlob)
}

func (i Int) asTimestamp() (time.Time, error) {
	return time.Time{}, castError(i.Type(), TypeTimestamp)
}

func (i Int) asArray() (Array, error) {
	return nil, castError(i.Type(), TypeArray)
}

func (i Int) asMap() (Map, error) {
	return nil, castError(i.Type(), TypeMap)
}

func (i Int) clone() Value {
	return Int(i)
}

// String returns JSON representation of an Int.
func (i Int) String() string {
	// the String return value is defined via the
	// default JSON serialization
	bytes, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("(unserializable int: %v)", err)
	}
	return string(bytes)
}

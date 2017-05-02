package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Array is an array of Values. It can be assigned to Value interface.
type Array []Value

// Type returns TypeID of Array. It's always TypeArray.
func (a Array) Type() TypeID {
	return TypeArray
}

func (a Array) asBool() (bool, error) {
	return false, castError(a.Type(), TypeBool)
}

func (a Array) asInt() (int64, error) {
	return 0, castError(a.Type(), TypeInt)
}

func (a Array) asFloat() (float64, error) {
	return 0, castError(a.Type(), TypeFloat)
}

func (a Array) asString() (string, error) {
	return "", castError(a.Type(), TypeString)
}

func (a Array) asBlob() ([]byte, error) {
	return nil, castError(a.Type(), TypeBlob)
}

func (a Array) asTimestamp() (time.Time, error) {
	return time.Time{}, castError(a.Type(), TypeTimestamp)
}

func (a Array) asArray() (Array, error) {
	return a, nil
}

func (a Array) asMap() (Map, error) {
	return nil, castError(a.Type(), TypeMap)
}

func (a Array) clone() Value {
	return a.Copy()
}

// String returns JSON representation of an Array.
func (a Array) String() string {
	// the String return value is defined via the
	// default JSON serialization
	bytes, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("(unserializable array: %v)", err)
	}
	return string(bytes)
}

// UnmarshalJSON reconstructs an Array from JSON.
func (a *Array) UnmarshalJSON(data []byte) error {
	var j []interface{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&j); err != nil {
		return err
	}

	newArray, err := NewArray(j)
	if err != nil {
		return err
	}
	*a = newArray
	return nil
}

// Copy performs deep copy of an Array. The Array returned from this method can
// safely be modified without affecting the original.
func (a Array) Copy() Array {
	out := make(Array, len(a))
	for idx, val := range a {
		out[idx] = val.clone()
	}
	return out
}

package data

import (
	"fmt"
	"time"
)

type Bool bool

func (b Bool) Type() TypeID {
	return TypeBool
}

func (b Bool) asBool() (bool, error) {
	return bool(b), nil
}

func (b Bool) asInt() (int64, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Bool) asFloat() (float64, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Bool) asString() (string, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Bool) asBlob() ([]byte, error) {
	return nil, castError(b.Type(), TypeBlob)
}

func (b Bool) asTimestamp() (time.Time, error) {
	return time.Time{}, castError(b.Type(), TypeTimestamp)
}

func (b Bool) asArray() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Bool) asMap() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

func (b Bool) clone() Value {
	return Bool(b)
}

func (b Bool) String() string {
	return fmt.Sprintf("%#v", b)
}

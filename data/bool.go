package data

import (
	"fmt"
	"time"
)

// Bool is a boolean value. It can be assigned to Value interface.
type Bool bool

const (
	// True is a constant having true value of Bool type.
	True Bool = true

	// False is a constant having false value of Bool type.
	False Bool = false
)

// Type returns TypeID of Bool. It's always TypeBool.
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

// String returns JSON representation of a Bool, which is true or false.
func (b Bool) String() string {
	return fmt.Sprintf("%#v", b)
}

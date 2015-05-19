package tuple

import (
	"fmt"
	"time"
)

type Bool bool

func (b Bool) Type() TypeID {
	return TypeBool
}

func (b Bool) AsBool() (bool, error) {
	return bool(b), nil
}

func (b Bool) AsInt() (int64, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Bool) AsFloat() (float64, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Bool) AsString() (string, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Bool) AsBlob() ([]byte, error) {
	return nil, castError(b.Type(), TypeBlob)
}

func (b Bool) AsTimestamp() (time.Time, error) {
	return time.Time{}, castError(b.Type(), TypeTimestamp)
}

func (b Bool) AsArray() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Bool) AsMap() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

func (b Bool) clone() Value {
	return Bool(b)
}

func (b Bool) String() string {
	return fmt.Sprintf("%#v", b)
}

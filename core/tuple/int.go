package tuple

import (
	"fmt"
	"time"
)

type Int int64

func (i Int) Type() TypeID {
	return TypeInt
}

func (i Int) AsBool() (bool, error) {
	return false, castError(i.Type(), TypeBool)
}

func (i Int) AsInt() (int64, error) {
	return int64(i), nil
}

func (i Int) AsFloat() (float64, error) {
	return float64(i), nil
}

func (i Int) AsString() (string, error) {
	return fmt.Sprint(i), nil
}

func (i Int) AsBlob() ([]byte, error) {
	return nil, castError(i.Type(), TypeBlob)
}

func (i Int) AsTimestamp() (time.Time, error) {
	return time.Time{}, castError(i.Type(), TypeTimestamp)
}

func (i Int) AsArray() (Array, error) {
	return nil, castError(i.Type(), TypeArray)
}

func (i Int) AsMap() (Map, error) {
	return nil, castError(i.Type(), TypeMap)
}

func (i Int) clone() Value {
	return Int(i)
}

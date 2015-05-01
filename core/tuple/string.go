package tuple

import (
	"time"
)

type String string

func (s String) Type() TypeID {
	return TypeString
}

func (s String) AsBool() (bool, error) {
	return false, castError(s.Type(), TypeBool)
}

func (s String) AsInt() (int64, error) {
	return 0, castError(s.Type(), TypeInt)
}

func (s String) AsFloat() (float64, error) {
	return 0, castError(s.Type(), TypeFloat)
}

func (s String) AsString() (string, error) {
	return string(s), nil
}

func (s String) AsBlob() ([]byte, error) {
	return nil, castError(s.Type(), TypeBlob)
}

func (s String) AsTimestamp() (time.Time, error) {
	return time.Time{}, castError(s.Type(), TypeTimestamp)
}

func (s String) AsArray() (Array, error) {
	return nil, castError(s.Type(), TypeArray)
}

func (s String) AsMap() (Map, error) {
	return nil, castError(s.Type(), TypeMap)
}

func (s String) clone() Value {
	return String(s)
}

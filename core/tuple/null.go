package tuple

import (
	"time"
)

type Null struct{}

func (n Null) Type() TypeID {
	return TypeNull
}

func (n Null) AsBool() (bool, error) {
	return false, castError(n.Type(), TypeBool)
}

func (n Null) AsInt() (int64, error) {
	return 0, castError(n.Type(), TypeInt)
}

func (n Null) AsFloat() (float64, error) {
	return 0, castError(n.Type(), TypeFloat)
}

func (n Null) AsString() (string, error) {
	return "", castError(n.Type(), TypeString)
}

func (n Null) AsBlob() ([]byte, error) {
	return nil, castError(n.Type(), TypeBlob)
}

func (n Null) AsTimestamp() (time.Time, error) {
	return time.Time{}, castError(n.Type(), TypeTimestamp)
}

func (n Null) AsArray() (Array, error) {
	return nil, castError(n.Type(), TypeArray)
}

func (n Null) AsMap() (Map, error) {
	return nil, castError(n.Type(), TypeMap)
}

func (n Null) clone() Value {
	return Null{}
}

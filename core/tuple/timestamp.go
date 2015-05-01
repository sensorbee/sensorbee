package tuple

import (
	"time"
)

type Timestamp time.Time

func (t Timestamp) Type() TypeID {
	return TypeTimestamp
}

func (t Timestamp) AsBool() (Bool, error) {
	return false, castError(t.Type(), TypeBool)
}

func (t Timestamp) AsInt() (Int, error) {
	return 0, castError(t.Type(), TypeInt)
}

func (t Timestamp) AsFloat() (Float, error) {
	return 0, castError(t.Type(), TypeFloat)
}

func (t Timestamp) AsString() (String, error) {
	return "", castError(t.Type(), TypeString)
}

func (t Timestamp) AsBlob() (Blob, error) {
	return nil, castError(t.Type(), TypeBlob)
}

func (t Timestamp) AsTimestamp() (Timestamp, error) {
	return t, nil
}

func (t Timestamp) AsArray() (Array, error) {
	return nil, castError(t.Type(), TypeArray)
}

func (t Timestamp) AsMap() (Map, error) {
	return nil, castError(t.Type(), TypeMap)
}

func (t Timestamp) clone() Value {
	return Timestamp(t)
}

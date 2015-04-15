package tuple

import (
	"time"
)

type Timestamp time.Time

func (t Timestamp) Type() TypeID {
	return TypeTimestamp
}

func (t Timestamp) Bool() (Bool, error) {
	return false, castError(t.Type(), TypeBool)
}

func (t Timestamp) Int() (Int, error) {
	return 0, castError(t.Type(), TypeInt)
}

func (t Timestamp) Float() (Float, error) {
	return 0, castError(t.Type(), TypeFloat)
}

func (t Timestamp) String() (String, error) {
	return "", castError(t.Type(), TypeString)
}

func (t Timestamp) Blob() (Blob, error) {
	return nil, castError(t.Type(), TypeBlob)
}

func (t Timestamp) Timestamp() (Timestamp, error) {
	return t, nil
}

func (t Timestamp) Array() (Array, error) {
	return nil, castError(t.Type(), TypeArray)
}

func (t Timestamp) Map() (Map, error) {
	return nil, castError(t.Type(), TypeMap)
}

func (t Timestamp) clone() Value {
	return Timestamp(t)
}

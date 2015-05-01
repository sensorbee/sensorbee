package tuple

import (
	"time"
)

type Blob []byte

func (b Blob) Type() TypeID {
	return TypeBlob
}

func (b Blob) AsBool() (bool, error) {
	return false, castError(b.Type(), TypeBool)
}

func (b Blob) AsInt() (int64, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Blob) AsFloat() (float64, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Blob) AsString() (string, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Blob) AsBlob() ([]byte, error) {
	return b, nil
}

func (b Blob) AsTimestamp() (time.Time, error) {
	return time.Time{}, castError(b.Type(), TypeTimestamp)
}

func (b Blob) AsArray() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Blob) AsMap() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

func (b Blob) clone() Value {
	out := make([]byte, len(b))
	for idx, val := range b {
		out[idx] = val
	}
	return Blob(out)
}

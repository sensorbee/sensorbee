package tuple

import (
	"fmt"
	"time"
)

type Float float64

func (f Float) Type() TypeID {
	return TypeFloat
}

func (f Float) AsBool() (bool, error) {
	return false, castError(f.Type(), TypeBool)
}

func (f Float) AsInt() (int64, error) {
	return int64(f), nil
}

func (f Float) AsFloat() (float64, error) {
	return float64(f), nil
}

func (f Float) AsString() (string, error) {
	return fmt.Sprint(f), nil
}

func (f Float) AsBlob() ([]byte, error) {
	return nil, castError(f.Type(), TypeFloat)
}

func (f Float) AsTimestamp() (time.Time, error) {
	return time.Time{}, castError(f.Type(), TypeTimestamp)
}

func (f Float) AsArray() (Array, error) {
	return nil, castError(f.Type(), TypeArray)
}

func (f Float) AsMap() (Map, error) {
	return nil, castError(f.Type(), TypeMap)
}

func (f Float) clone() Value {
	return Float(f)
}

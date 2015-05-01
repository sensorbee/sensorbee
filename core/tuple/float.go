package tuple

import (
	"fmt"
)

type Float float64

func (f Float) Type() TypeID {
	return TypeFloat
}

func (f Float) AsBool() (Bool, error) {
	return false, castError(f.Type(), TypeBool)
}

func (f Float) AsInt() (Int, error) {
	return Int(f), nil
}

func (f Float) AsFloat() (Float, error) {
	return f, nil
}

func (f Float) AsString() (String, error) {
	return String(fmt.Sprint(f)), nil
}

func (f Float) AsBlob() (Blob, error) {
	return nil, castError(f.Type(), TypeFloat)
}

func (f Float) AsTimestamp() (Timestamp, error) {
	return Timestamp{}, castError(f.Type(), TypeTimestamp)
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

package tuple

import (
	"fmt"
)

type Int int64

func (i Int) Type() TypeID {
	return TypeInt
}

func (i Int) AsBool() (Bool, error) {
	return false, castError(i.Type(), TypeBool)
}

func (i Int) AsInt() (Int, error) {
	return i, nil
}

func (i Int) AsFloat() (Float, error) {
	return Float(i), nil
}

func (i Int) AsString() (String, error) {
	return String(fmt.Sprint(i)), nil
}

func (i Int) AsBlob() (Blob, error) {
	return nil, castError(i.Type(), TypeBlob)
}

func (i Int) AsTimestamp() (Timestamp, error) {
	return Timestamp{}, castError(i.Type(), TypeTimestamp)
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

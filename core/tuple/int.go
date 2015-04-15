package tuple

import (
	"fmt"
)

type Int int64

func (i Int) Type() TypeID {
	return TypeInt
}

func (i Int) Bool() (Bool, error) {
	return false, castError(i.Type(), TypeBool)
}

func (i Int) Int() (Int, error) {
	return i, nil
}

func (i Int) Float() (Float, error) {
	return Float(i), nil
}

func (i Int) String() (String, error) {
	return String(fmt.Sprint(i)), nil
}

func (i Int) Blob() (Blob, error) {
	return nil, castError(i.Type(), TypeBlob)
}

func (i Int) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(i.Type(), TypeTimestamp)
}

func (i Int) Array() (Array, error) {
	return nil, castError(i.Type(), TypeArray)
}

func (i Int) Map() (Map, error) {
	return nil, castError(i.Type(), TypeMap)
}

func (i Int) clone() Value {
	return Int(i)
}

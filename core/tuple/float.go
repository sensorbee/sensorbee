package tuple

import (
	"fmt"
)

type Float float64

func (f Float) Type() TypeID {
	return TypeFloat
}

func (f Float) Bool() (Bool, error) {
	return false, castError(f.Type(), TypeBool)
}

func (f Float) Int() (Int, error) {
	return Int(f), nil
}

func (f Float) Float() (Float, error) {
	return f, nil
}

func (f Float) String() (String, error) {
	return String(fmt.Sprint(f)), nil
}

func (f Float) Blob() (Blob, error) {
	return nil, castError(f.Type(), TypeFloat)
}

func (f Float) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(f.Type(), TypeTimestamp)
}

func (f Float) Array() (Array, error) {
	return nil, castError(f.Type(), TypeArray)
}

func (f Float) Map() (Map, error) {
	return nil, castError(f.Type(), TypeMap)
}

func (f Float) clone() Value {
	return Float(f)
}

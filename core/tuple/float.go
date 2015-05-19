package tuple

import (
	"fmt"
	"math"
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
	return 0, castError(f.Type(), TypeInt)
}

func (f Float) AsFloat() (float64, error) {
	return float64(f), nil
}

func (f Float) AsString() (string, error) {
	return "", castError(f.Type(), TypeString)
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

func (f Float) MarshalJSON() ([]byte, error) {
	// the JSON serialization is defined via the String()
	// return value as defined below
	return []byte(f.String()), nil
}

func (f Float) String() string {
	fl := float64(f)
	// "NaN and Infinity regardless of sign are represented
	// as the String null." (ECMA-262)
	// (The default JSON serializer will return an error instead,
	// cf. <https://github.com/golang/go/issues/3480>)
	if math.IsNaN(fl) {
		return "null"
	} else if math.IsInf(fl, 0) {
		return "null"
	}
	return fmt.Sprintf("%#v", f)
}

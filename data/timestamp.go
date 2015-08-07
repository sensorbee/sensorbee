package data

import (
	"time"
)

type Timestamp time.Time

func (t Timestamp) Type() TypeID {
	return TypeTimestamp
}

func (t Timestamp) asBool() (bool, error) {
	return false, castError(t.Type(), TypeBool)
}

func (t Timestamp) asInt() (int64, error) {
	return 0, castError(t.Type(), TypeInt)
}

func (t Timestamp) asFloat() (float64, error) {
	return 0, castError(t.Type(), TypeFloat)
}

func (t Timestamp) asString() (string, error) {
	return "", castError(t.Type(), TypeString)
}

func (t Timestamp) asBlob() ([]byte, error) {
	return nil, castError(t.Type(), TypeBlob)
}

func (t Timestamp) asTimestamp() (time.Time, error) {
	return time.Time(t), nil
}

func (t Timestamp) asArray() (Array, error) {
	return nil, castError(t.Type(), TypeArray)
}

func (t Timestamp) asMap() (Map, error) {
	return nil, castError(t.Type(), TypeMap)
}

func (t Timestamp) clone() Value {
	return Timestamp(t)
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	// the JSON serialization is defined via the String()
	// return value as defined below
	return []byte(t.String()), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	ts, err := time.Parse(`"`+time.RFC3339Nano+`"`, str)
	if err == nil {
		*t = Timestamp(ts)
		return nil
	}
	ts, err = time.Parse(`"`+time.RFC3339+`"`, str)
	if err == nil {
		*t = Timestamp(ts)
		return nil
	}
	return err
}

func (t Timestamp) String() string {
	s, _ := ToString(t)
	return `"` + s + `"`
}

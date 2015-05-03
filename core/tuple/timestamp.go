package tuple

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Timestamp time.Time

func (t Timestamp) Type() TypeID {
	return TypeTimestamp
}

func (t Timestamp) AsBool() (bool, error) {
	return false, castError(t.Type(), TypeBool)
}

func (t Timestamp) AsInt() (int64, error) {
	return 0, castError(t.Type(), TypeInt)
}

func (t Timestamp) AsFloat() (float64, error) {
	return 0, castError(t.Type(), TypeFloat)
}

func (t Timestamp) AsString() (string, error) {
	return "", castError(t.Type(), TypeString)
}

func (t Timestamp) AsBlob() ([]byte, error) {
	return nil, castError(t.Type(), TypeBlob)
}

func (t Timestamp) AsTimestamp() (time.Time, error) {
	return time.Time(t), nil
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

// timeExt is msgpack extension struct for Timestamp
type timeExt struct {
}

func (t *timeExt) WriteExt(v interface{}) []byte {
	buf := new(bytes.Buffer)
	var nanoTime int64
	switch vt := v.(type) {
	case time.Time:
		nanoTime = vt.UTC().UnixNano()
	case *time.Time:
		nanoTime = vt.UTC().UnixNano()
	default:
		return nil
	}
	err := binary.Write(buf, binary.LittleEndian, nanoTime)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func (t *timeExt) ReadExt(dst interface{}, src []byte) {
	buf := bytes.NewReader(src)
	var nanoTime int64
	binary.Read(buf, binary.LittleEndian, &nanoTime)
	tt := dst.(*time.Time)
	*tt = time.Unix(0, nanoTime).UTC()
}

func (t *timeExt) ConvertExt(v interface{}) interface{} {
	switch vt := v.(type) {
	case time.Time:
		return vt.UTC().UnixNano()
	case *time.Time:
		return vt.UTC().UnixNano()
	default:
		return nil
	}
}
func (t *timeExt) UpdateExt(dst interface{}, src interface{}) {
	tt := dst.(*time.Time)
	switch vt := src.(type) {
	case int64:
		*tt = time.Unix(0, vt).UTC()
	case uint64:
		*tt = time.Unix(0, int64(vt)).UTC()
	default:
		// do nothing
	}
}

package data

import (
	"encoding/json"
	"fmt"
	"time"
)

// Blob is a binary large object which may have any type of byte data.
// It can be assigned to Value interface.
type Blob []byte

// Type returns TypeID of Blob. It's always TypeBlob.
func (b Blob) Type() TypeID {
	return TypeBlob
}

func (b Blob) asBool() (bool, error) {
	return false, castError(b.Type(), TypeBool)
}

func (b Blob) asInt() (int64, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Blob) asFloat() (float64, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Blob) asString() (string, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Blob) asBlob() ([]byte, error) {
	return b, nil
}

func (b Blob) asTimestamp() (time.Time, error) {
	return time.Time{}, castError(b.Type(), TypeTimestamp)
}

func (b Blob) asArray() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Blob) asMap() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

func (b Blob) clone() Value {
	out := make([]byte, len(b))
	for idx, val := range b {
		out[idx] = val
	}
	return Blob(out)
}

// MarshalJSON marshals a Blob to JSON. nil will be encoded as null. It doesn't
// encode a Blob in base64 format.
func (b Blob) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}

	// To avoid base64 encoding done by Go's encoding/json, []byte must be
	// converted to string.
	// TODO: reduce this extra copy
	return json.Marshal(string(b))
}

// UnmarshalJSON unmarshals a JSON string to a Blob. It doesn't parse base64
// format.
func (b *Blob) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*b = []byte(s)
	return nil
}

// Stringreturns JSON representation of a Blob. Blob is marshaled as a string.
func (b Blob) String() string {
	bytes, err := b.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("(unserializable blob: %v)", err)
	}
	return string(bytes)
}

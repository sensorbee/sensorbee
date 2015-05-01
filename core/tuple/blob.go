package tuple

type Blob []byte

func (b Blob) Type() TypeID {
	return TypeBlob
}

func (b Blob) AsBool() (Bool, error) {
	return false, castError(b.Type(), TypeBool)
}

func (b Blob) AsInt() (Int, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Blob) AsFloat() (Float, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Blob) AsString() (String, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Blob) AsBlob() (Blob, error) {
	return b, nil
}

func (b Blob) AsTimestamp() (Timestamp, error) {
	return Timestamp{}, castError(b.Type(), TypeTimestamp)
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

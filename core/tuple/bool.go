package tuple

type Bool bool

func (b Bool) Type() TypeID {
	return TypeBool
}

func (b Bool) AsBool() (Bool, error) {
	return b, nil
}

func (b Bool) AsInt() (Int, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Bool) AsFloat() (Float, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Bool) AsString() (String, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Bool) AsBlob() (Blob, error) {
	return nil, castError(b.Type(), TypeBlob)
}

func (b Bool) AsTimestamp() (Timestamp, error) {
	return Timestamp{}, castError(b.Type(), TypeTimestamp)
}

func (b Bool) AsArray() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Bool) AsMap() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

func (b Bool) clone() Value {
	return Bool(b)
}

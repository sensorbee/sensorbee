package tuple

type String string

func (s String) Type() TypeID {
	return TypeString
}

func (s String) AsBool() (Bool, error) {
	return false, castError(s.Type(), TypeBool)
}

func (s String) AsInt() (Int, error) {
	return 0, castError(s.Type(), TypeInt)
}

func (s String) AsFloat() (Float, error) {
	return 0, castError(s.Type(), TypeFloat)
}

func (s String) AsString() (String, error) {
	return s, nil
}

func (s String) AsBlob() (Blob, error) {
	return nil, castError(s.Type(), TypeBlob)
}

func (s String) AsTimestamp() (Timestamp, error) {
	return Timestamp{}, castError(s.Type(), TypeTimestamp)
}

func (s String) AsArray() (Array, error) {
	return nil, castError(s.Type(), TypeArray)
}

func (s String) AsMap() (Map, error) {
	return nil, castError(s.Type(), TypeMap)
}

func (s String) clone() Value {
	return String(s)
}

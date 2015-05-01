package tuple

type Array []Value

func (a Array) Type() TypeID {
	return TypeArray
}

func (a Array) AsBool() (Bool, error) {
	return false, castError(a.Type(), TypeBool)
}

func (a Array) AsInt() (Int, error) {
	return 0, castError(a.Type(), TypeInt)
}

func (a Array) AsFloat() (Float, error) {
	return 0, castError(a.Type(), TypeFloat)
}

func (a Array) AsString() (String, error) {
	return "", castError(a.Type(), TypeString)
}

func (a Array) AsBlob() (Blob, error) {
	return nil, castError(a.Type(), TypeBlob)
}

func (a Array) AsTimestamp() (Timestamp, error) {
	return Timestamp{}, castError(a.Type(), TypeTimestamp)
}

func (a Array) AsArray() (Array, error) {
	return a, nil
}

func (a Array) AsMap() (Map, error) {
	return nil, castError(a.Type(), TypeMap)
}

func (a Array) clone() Value {
	out := make([]Value, len(a))
	for idx, val := range a {
		out[idx] = val.clone()
	}
	return Array(out)
}

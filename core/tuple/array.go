package tuple

type Array []Value

func (a Array) Type() TypeID {
	return TypeArray
}

func (a Array) Bool() (Bool, error) {
	return false, castError(a.Type(), TypeBool)
}

func (a Array) Int() (Int, error) {
	return 0, castError(a.Type(), TypeInt)
}

func (a Array) Float() (Float, error) {
	return 0, castError(a.Type(), TypeFloat)
}

func (a Array) String() (String, error) {
	return "", castError(a.Type(), TypeString)
}

func (a Array) Blob() (Blob, error) {
	return nil, castError(a.Type(), TypeBlob)
}

func (a Array) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(a.Type(), TypeTimestamp)
}

func (a Array) Array() (Array, error) {
	return a, nil
}

func (a Array) Map() (Map, error) {
	return nil, castError(a.Type(), TypeMap)
}

func (a Array) clone() Value {
	out := make([]Value, len(a))
	for idx, val := range a {
		out[idx] = val.clone()
	}
	return Array(out)
}

package tuple

type Bool bool

func (b Bool) Type() TypeID {
	return TypeBool
}

func (b Bool) Bool() (Bool, error) {
	return b, nil
}

func (b Bool) Int() (Int, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Bool) Float() (Float, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Bool) String() (String, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Bool) Blob() (Blob, error) {
	return nil, castError(b.Type(), TypeBlob)
}

func (b Bool) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(b.Type(), TypeTimestamp)
}

func (b Bool) Array() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Bool) Map() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

func (b Bool) clone() Value {
	return Bool(b)
}

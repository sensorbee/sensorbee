package tuple

type Null struct{}

func (n Null) Type() TypeID {
	return TypeNull
}

func (n Null) AsBool() (Bool, error) {
	return false, castError(n.Type(), TypeBool)
}

func (n Null) AsInt() (Int, error) {
	return 0, castError(n.Type(), TypeInt)
}

func (n Null) AsFloat() (Float, error) {
	return 0, castError(n.Type(), TypeFloat)
}

func (n Null) AsString() (String, error) {
	return "", castError(n.Type(), TypeString)
}

func (n Null) AsBlob() (Blob, error) {
	return nil, castError(n.Type(), TypeBlob)
}

func (n Null) AsTimestamp() (Timestamp, error) {
	return Timestamp{}, castError(n.Type(), TypeTimestamp)
}

func (n Null) AsArray() (Array, error) {
	return nil, castError(n.Type(), TypeArray)
}

func (n Null) AsMap() (Map, error) {
	return nil, castError(n.Type(), TypeMap)
}

func (n Null) clone() Value {
	return Null{}
}

package tuple

type Blob []byte

func (b Blob) Type() TypeID {
	return TypeBlob
}

func (b Blob) Bool() (Bool, error) {
	return false, castError(b.Type(), TypeBool)
}

func (b Blob) Int() (Int, error) {
	return 0, castError(b.Type(), TypeInt)
}

func (b Blob) Float() (Float, error) {
	return 0, castError(b.Type(), TypeFloat)
}

func (b Blob) String() (String, error) {
	return "", castError(b.Type(), TypeString)
}

func (b Blob) Blob() (Blob, error) {
	return b, nil
}

func (b Blob) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(b.Type(), TypeTimestamp)
}

func (b Blob) Array() (Array, error) {
	return nil, castError(b.Type(), TypeArray)
}

func (b Blob) Map() (Map, error) {
	return nil, castError(b.Type(), TypeMap)
}

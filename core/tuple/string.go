package tuple

type String string

func (s String) Type() TypeID {
	return TypeString
}

func (s String) Bool() (Bool, error) {
	return false, castError(s.Type(), TypeBool)
}

func (s String) Int() (Int, error) {
	return 0, castError(s.Type(), TypeInt)
}

func (s String) Float() (Float, error) {
	return 0, castError(s.Type(), TypeFloat)
}

func (s String) String() (String, error) {
	return s, nil
}

func (s String) Blob() (Blob, error) {
	return nil, castError(s.Type(), TypeBlob)
}

func (s String) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(s.Type(), TypeTimestamp)
}

func (s String) Array() (Array, error) {
	return nil, castError(s.Type(), TypeArray)
}

func (s String) Map() (Map, error) {
	return nil, castError(s.Type(), TypeMap)
}

func (s String) clone() Value {
	return String(s)
}

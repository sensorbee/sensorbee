package tuple

type Null struct{}

func (n Null) Type() TypeID {
	return TypeNull
}

func (n Null) Bool() (Bool, error) {
	return false, castError(n.Type(), TypeBool)
}

func (n Null) Int() (Int, error) {
	return 0, castError(n.Type(), TypeInt)
}

func (n Null) Float() (Float, error) {
	return 0, castError(n.Type(), TypeFloat)
}

func (n Null) String() (String, error) {
	return "", castError(n.Type(), TypeString)
}

func (n Null) Blob() (Blob, error) {
	return nil, castError(n.Type(), TypeBlob)
}

func (n Null) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(n.Type(), TypeTimestamp)
}

func (n Null) Array() (Array, error) {
	return nil, castError(n.Type(), TypeArray)
}

func (n Null) Map() (Map, error) {
	return nil, castError(n.Type(), TypeMap)
}

func (n Null) clone() Value {
	return Null{}
}

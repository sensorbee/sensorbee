package tuple

type Map map[string]Value

func (m Map) Type() TypeID {
	return TypeMap
}

func (m Map) Bool() (Bool, error) {
	return false, castError(m.Type(), TypeBool)
}

func (m Map) Int() (Int, error) {
	return 0, castError(m.Type(), TypeInt)
}

func (m Map) Float() (Float, error) {
	return 0, castError(m.Type(), TypeFloat)
}

func (m Map) String() (String, error) {
	return "", castError(m.Type(), TypeString)
}

func (m Map) Blob() (Blob, error) {
	return nil, castError(m.Type(), TypeBlob)
}

func (m Map) Timestamp() (Timestamp, error) {
	return Timestamp{}, castError(m.Type(), TypeTimestamp)
}

func (m Map) Array() (Array, error) {
	return nil, castError(m.Type(), TypeArray)
}

func (m Map) Map() (Map, error) {
	return m, nil
}

func (m Map) clone() Value {
	return m.Copy()
}

func (m Map) Copy() Map {
	out := make(map[string]Value, len(m))
	for key, val := range m {
		out[key] = val.clone()
	}
	return Map(out)
}

func (m Map) Get(path string) (Value, error) {
	// TODO: support json path manually
	var v Value
	err := ScanTree(m.toMapInterface(), path, &v)
	return v, err
}

// toMapInterface converts Map to map[string]interface{}.
// This is only for go-scan.
func (m Map) toMapInterface() map[string]interface{} {
	t := make(map[string]interface{}, len(m))
	for k, v := range m {
		switch v.Type() {
		case TypeBlob:
			b, _ := v.Blob()
			t[k] = &b
		case TypeArray:
			a, _ := v.Array()
			t[k] = a.toArrayInterface()
		case TypeMap:
			m, _ := v.Map()
			t[k] = m.toMapInterface()
		default:
			t[k] = v
		}
	}
	return t
}

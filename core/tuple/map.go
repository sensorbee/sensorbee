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

// Get value(s) from a structured Map followed by the path expression.
// Return error when the path expression is invalid or the path is not found
// in Map keys.
// Value interface can be cast for each types using Value interface's methods.
// Example:
//  v, err := map.Get("path")
//  s, err := v.String() // cast to String
//
// Path Expression Example:
// Given the following Map structure
//  Map{
//  	"store": Map{
//  		"name": String("store name"),
//  		"book": Array([]Value{
//  			Map{
//  				"title": String("book name"),
//  			},
//  		}),
//  	},
//  }
// To get values, access the following path expressions
//  `store`              -> get store's Map
//  `store.name`         -> get "store name"
//  `store.book[0].title -> get "book name"
// or
//  `["store"]`                     -> get store's Map
//  `["store"]["name"]`             -> get "store name"
//  `["store"]["book"][0]["title"]` -> get "book name"
//
func (m Map) Get(path string) (Value, error) {
	var v Value
	err := scanMap(m, path, &v)
	return v, err
}

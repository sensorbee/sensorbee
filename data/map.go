package data

import (
	"encoding/json"
	"fmt"
	"time"
)

// Map is a map of Values. It can be assigned to Value interface. Only
// string keys are allowed.
type Map map[string]Value

// Type returns TypeID of Map. It's always TypeMap.
func (m Map) Type() TypeID {
	return TypeMap
}

func (m Map) asBool() (bool, error) {
	return false, castError(m.Type(), TypeBool)
}

func (m Map) asInt() (int64, error) {
	return 0, castError(m.Type(), TypeInt)
}

func (m Map) asFloat() (float64, error) {
	return 0, castError(m.Type(), TypeFloat)
}

func (m Map) asString() (string, error) {
	return "", castError(m.Type(), TypeString)
}

func (m Map) asBlob() ([]byte, error) {
	return nil, castError(m.Type(), TypeBlob)
}

func (m Map) asTimestamp() (time.Time, error) {
	return time.Time{}, castError(m.Type(), TypeTimestamp)
}

func (m Map) asArray() (Array, error) {
	return nil, castError(m.Type(), TypeArray)
}

func (m Map) asMap() (Map, error) {
	return m, nil
}

func (m Map) clone() Value {
	return m.Copy()
}

// String returns JSON representation of a Map.
func (m Map) String() string {
	// the String return value is defined via the
	// default JSON serialization
	bytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("(unserializable map: %v)", err)
	}
	return string(bytes)
}

// UnmarshalJSON reconstructs a Map from JSON.
func (m *Map) UnmarshalJSON(data []byte) error {
	var j map[string]interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	newMap, err := NewMap(j)
	if err != nil {
		return err
	}
	*m = newMap
	return nil
}

// Copy performs deep copy of a Map. The Map returned from this method can
// safely be modified without affecting the original.
func (m Map) Copy() Map {
	out := make(map[string]Value, len(m))
	for key, val := range m {
		out[key] = val.clone()
	}
	return Map(out)
}

// Get returns value(s) from a structured Map as addressed by the
// given path expression. Returns an error when the path expression
// is invalid or the path is not found in the Map.
//
// Type conversion can be done for each type using the Value
// interface's methods.
//
// Example:
//  p, err := CompilePath("path")
//  if err != nil { ... }
//  v, err := map.Get(p)
//  if err != nil { ... }
//  s, err := v.asString() // cast to String
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
func (m Map) Get(path Path) (Value, error) {
	return path.evaluate(m)
}

// Set sets a value in a structured Map as addressed by the
// given path expression. It is possible to set items even
// if the parent levels do not exist (`mkdir -p` behavior),
// i.e., you can set an item at "stores[17].owner.name" also
// when there is no "stores" key in the top-level Map. If a
// list index is higher than the current length, intermediate
// items will be created as Null items.
// Set returns an error when the path expression is invalid or
// when one of the intermediate components already exists
// but is not a map/list.
func (m Map) Set(path Path, val Value) error {
	return path.set(m, val)
}

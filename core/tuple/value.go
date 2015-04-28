package tuple

import (
	"errors"
	"fmt"
	"github.com/ugorji/go/codec"
	"time"
)

// Value is the generic interface for all data that can be stored
// inside a Tuple. Since we assume the data not to conform to any
// schema, data can have any shape and it can also change within a
// stream from one Tuple to the next. Therefore we need to be
// careful with respect to type conversions. A Value obtained, e.g.,
// by Map.Get should always be converted using the appropriate method
// and error checking must be done.
//
// Example:
//  i, err := val.AsInt()
//  if err != nil { ... }
type Value interface {
	Type() TypeID
	AsBool() (bool, error)
	AsInt() (int64, error)
	AsFloat() (float64, error)
	AsString() (string, error)
	AsBlob() ([]byte, error)
	AsTimestamp() (time.Time, error)
	AsArray() (Array, error)
	AsMap() (Map, error)
	clone() Value
}

// TODO: Provide NewMap(map[string]interface{}) Map

func castError(from TypeID, to TypeID) error {
	return errors.New(fmt.Sprintf("unsupported cast %v from %v", to.String(), from.String()))
}

type TypeID int

const (
	TypeUnknown TypeID = iota
	TypeNull
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeBlob
	TypeTimestamp
	TypeArray
	TypeMap
)

func (t TypeID) String() string {
	switch t {
	case TypeNull:
		return "null"
	case TypeBool:
		return "bool"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	case TypeBlob:
		return "blob"
	case TypeTimestamp:
		return "timestamp"
	case TypeArray:
		return "array"
	case TypeMap:
		return "map"
	default:
		return "unknown"
	}
}

var mh = &codec.MsgpackHandle{RawToString: true}

func UnmarshalMsgpack(b []byte) (Map, error) {
	var m map[interface{}]interface{}
	dec := codec.NewDecoderBytes(b, mh)
	dec.Decode(&m)

	return newMap(m)
}

func newMap(m map[interface{}]interface{}) (Map, error) {
	fmt.Println(m)
	result := Map{}
	for k, v := range m {
		key, ok := k.(string)
		if !ok {
			return nil, errors.New("Non string type key is not supported")
		}
		switch v.(type) {
		case []interface{}:
			innerArray, err := newArray(v.([]interface{}))
			if err != nil {
				return nil, err // TODO is it OK to return nil?
			}
			result[key] = Array(innerArray)
		case map[interface{}]interface{}:
			innerMap, err := newMap(v.(map[interface{}]interface{}))
			if err != nil {
				return nil, err // TODO is it OK to return nil?
			}
			result[key] = Map(innerMap)
		case bool:
			result[key] = Bool(v.(bool))
		case int:
			result[key] = Int(v.(int))
		case int8:
			result[key] = Int(v.(int8))
		case int16:
			result[key] = Int(v.(int16))
		case int32:
			result[key] = Int(v.(int32))
		case int64:
			result[key] = Int(v.(int64))
		case float32:
			result[key] = Float(v.(float32))
		case float64:
			result[key] = Float(v.(float64))
		case string:
			result[key] = String(v.(string))
		case []byte:
			result[key] = Blob(v.([]byte))
		case nil:
			result[key] = nil
		}
	}
	return result, nil
}

func newArray(a []interface{}) ([]Value, error) {
	result := make([]Value, len(a))
	for i, v := range a {
		switch v.(type) {
		case []interface{}:
			innerArray, err := newArray(v.([]interface{}))
			if err != nil {
				return nil, err
			}
			result[i] = Array(innerArray)
		case map[interface{}]interface{}:
			innerMap, err := newMap(v.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			result[i] = Map(innerMap)
		case bool:
			result[i] = Bool(v.(bool))
		case int:
			result[i] = Int(v.(int))
		case int8:
			result[i] = Int(v.(int8))
		case int16:
			result[i] = Int(v.(int16))
		case int32:
			result[i] = Int(v.(int32))
		case int64:
			result[i] = Int(v.(int64))
		case float32:
			result[i] = Float(v.(float32))
		case float64:
			result[i] = Float(v.(float64))
		case string:
			result[i] = String(v.(string))
		case []byte:
			result[i] = Blob(v.([]byte))
		case nil:
			result[i] = nil
		}
	}
	return result, nil
}

func MarshalMsgpack(m Map) ([]byte, error) {
	// iMap, err := newIMap(m)
	// if err != nil {
	// 	return nil, err
	// }
	var out []byte
	enc := codec.NewEncoderBytes(&out, mh)
	enc.Encode(m) // TODO is it unnecessary to use iMap?

	return out, nil
}

func newIMap(m Map) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for k, v := range m {
		switch v.Type() {
		case TypeBool:
			result[k], _ = v.AsBool()
		case TypeInt:
			result[k], _ = v.AsInt()
		case TypeFloat:
			result[k], _ = v.AsFloat()
		case TypeString:
			result[k], _ = v.AsString()
		case TypeBlob:
			result[k], _ = v.AsBlob()
		case TypeTimestamp:
			result[k], _ = v.AsTimestamp()
		case TypeArray:
			innerArray, _ := v.AsArray()
			result[k], _ = newIArray(innerArray)
		case TypeMap:
			innerMap, _ := v.AsMap()
			result[k], _ = newIMap(innerMap)
		}
	}
	return result, nil
}

func newIArray(a Array) ([]interface{}, error) {
	result := make([]interface{}, len(a))
	for i, v := range a {
		switch v.Type() {
		case TypeBool:
			result[i], _ = v.AsBool()
		case TypeInt:
			result[i], _ = v.AsInt()
		case TypeFloat:
			result[i], _ = v.AsFloat()
		case TypeString:
			result[i], _ = v.AsString()
		case TypeBlob:
			result[i], _ = v.AsBlob()
		case TypeArray:
			innerArray, _ := v.AsArray()
			result[i], _ = newIArray(innerArray)
		case TypeMap:
			innerMap, _ := v.AsMap()
			result[i], _ = newIMap(innerMap)
		}
	}
	return result, nil
}

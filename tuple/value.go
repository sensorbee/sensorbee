package tuple

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"math"
	"reflect"
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
//  i, err := val.asInt()
//  if err != nil { ... }
type Value interface {
	// Type returns the actual type of a Value. (Note that this is
	// faster than a type switch.) If a Value has `Type() == Type{X}`,
	// then it can be assumed that the `As{X}()` conversion will
	// not fail.
	Type() TypeID
	asBool() (bool, error)
	asInt() (int64, error)
	asFloat() (float64, error)
	asString() (string, error)
	asBlob() ([]byte, error)
	asTimestamp() (time.Time, error)
	asArray() (Array, error)
	asMap() (Map, error)
	clone() Value
	String() string
}

func castError(from TypeID, to TypeID) error {
	return fmt.Errorf("unsupported cast %v from %v", to.String(), from.String())
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

var msgpackHandle = &codec.MsgpackHandle{}

func init() {
	msgpackHandle.MapType = reflect.TypeOf(map[string]interface{}(nil))
	msgpackHandle.RawToString = true
	msgpackHandle.WriteExt = false
}

// UnmarshalMsgpack returns a Map object from a byte array encoded
// by msgpack serialization. The byte is expected to decode key-value
// style map. Returns an error when value type is not supported in SensorBee.
func UnmarshalMsgpack(b []byte) (Map, error) {
	var m map[string]interface{}
	dec := codec.NewDecoderBytes(b, msgpackHandle)
	dec.Decode(&m)

	return NewMap(m)
}

// NewMap returns a Map object from map[string]interface{}.
// Returns an error when value type is not supported in SensorBee.
//
// Example:
// The following sample interface{} will be converted to mapSample Map.
//   var sample = map[string]interface{}{
//      "bool":   true,
//      "int":    int64(1),
//      "float":  float64(0.1),
//      "string": "homhom",
//      "time":   time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC),
//      "array": []interface{}{true, 10, "inarray",
//          map[string]interface{}{
//              "mapinarray": "arraymap",
//          }},
//      "map": map[string]interface{}{
//          "map_a": "a",
//          "map_b": 2,
//      },
//      "byte": []byte("test byte"),
//      "null": nil,
//  }
//  var mapSample = Map{
//      "bool":   Bool(true),
//      "int":    Int(1),
//      "float":  Float(0.1),
//      "string": String("homhom"),
//      "time":   Timestamp(time.Date(2015, time.May, 1, 14, 27, 0, 0, time.UTC)),
//      "array": Array([]Value{Bool(true), Int(10), String("inarray"),
//          Map{
//              "mapinarray": String("arraymap"),
//          }}),
//      "map": Map{
//          "map_a": String("a"),
//          "map_b": Int(2),
//      },
//      "byte": Blob([]byte("test byte")),
//      "null": Null{},
//  }
//
func NewMap(m map[string]interface{}) (Map, error) {
	result := Map{}
	for k, v := range m {
		value, err := NewValue(v)
		if err != nil {
			return nil, err
		}
		result[k] = value
	}
	return result, nil
}

// NewArray returns a Array object from []interface{}.
// Returns an error when value type is not supported in SensorBee.
func NewArray(a []interface{}) (Array, error) {
	result := make([]Value, len(a))
	for i, v := range a {
		value, err := NewValue(v)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
}

// NewValue returns a Value object from interface{}.
// Returns an error when value type is not supported in SensorBee.
func NewValue(v interface{}) (result Value, err error) {
	switch vt := v.(type) {
	case []interface{}:
		a, err := NewArray(vt)
		if err != nil {
			return nil, err
		}
		result = a
	case map[string]interface{}:
		m, err := NewMap(vt)
		if err != nil {
			return nil, err
		}
		result = m
	case bool:
		result = Bool(vt)
	case int:
		result = Int(vt)
	case int8:
		result = Int(vt)
	case int16:
		result = Int(vt)
	case int32:
		result = Int(vt)
	case int64:
		result = Int(vt)
	case uint:
		if vt > math.MaxInt64 {
			err = fmt.Errorf("an int value must be less than 2^63: %v", vt)
			break
		}
		result = Int(vt)
	case uint8:
		result = Int(vt)
	case uint16:
		result = Int(vt)
	case uint32:
		result = Int(vt)
	case uint64:
		if vt > math.MaxInt64 {
			err = fmt.Errorf("an int value must be less than 2^63: %v", vt)
			break
		}
		result = Int(vt)
	case float32:
		result = Float(vt)
	case float64:
		result = Float(vt)
	case time.Time:
		result = Timestamp(vt)
	case string:
		result = String(vt)
	case []byte:
		result = Blob(vt)
	case nil:
		result = Null{}

	// support some tuple types for convenience
	case Array:
		result = vt
	case Map:
		result = vt
	case Blob:
		result = vt
	case Timestamp:
		result = vt
	default:
		err = fmt.Errorf("unsupported type %T", v)
	}
	return result, err
}

// MarshalMsgpack returns a byte array encoded by msgpack serialization
// from a Map object. Returns an error when msgpack serialization failed.
func MarshalMsgpack(m Map) ([]byte, error) {
	iMap := newIMap(m)
	var out []byte
	enc := codec.NewEncoderBytes(&out, msgpackHandle)
	err := enc.Encode(iMap)

	return out, err
}

func newIMap(m Map) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range m {
		value := newIValue(v)
		result[k] = value
	}
	return result
}

func newIArray(a Array) []interface{} {
	result := make([]interface{}, len(a))
	for i, v := range a {
		value := newIValue(v)
		result[i] = value
	}
	return result
}

func newIValue(v Value) interface{} {
	var result interface{}
	switch v.Type() {
	case TypeBool:
		result, _ = v.asBool()
	case TypeInt:
		result, _ = v.asInt()
	case TypeFloat:
		result, _ = v.asFloat()
	case TypeString:
		result, _ = v.asString()
	case TypeBlob:
		result, _ = v.asBlob()
	case TypeTimestamp:
		result, _ = ToInt(v)
	case TypeArray:
		innerArray, _ := v.asArray()
		result = newIArray(innerArray)
	case TypeMap:
		innerMap, _ := v.asMap()
		result = newIMap(innerMap)
	case TypeNull:
		result = nil
	default:
		//do nothing
	}
	return result
}

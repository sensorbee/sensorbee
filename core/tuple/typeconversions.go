package tuple

import (
	"fmt"
)

// ToBool converts a given Value to a bool, if possible. The conversion
// rules are similar to those in Python:
//
//  * Null: false
//  * Bool: actual boolean value
//  * Int: true if non-zero
//  * Float: true if non-zero
//  * String: true if non-empty
//  * Blob: true if non-empty
//  * Timestamp: true if IsZero() is false
//  * Array: true if non-empty
//  * Map: true if non-empty
func ToBool(v Value) (bool, error) {
	defaultValue := false
	switch v.Type() {
	case TypeNull:
		return false, nil
	case TypeBool:
		val, e := v.AsBool()
		if e != nil {
			return defaultValue, e
		}
		return val, nil
	case TypeInt:
		val, e := v.AsInt()
		if e != nil {
			return defaultValue, e
		}
		return val != 0, nil
	case TypeFloat:
		val, e := v.AsFloat()
		if e != nil {
			return defaultValue, e
		}
		return val != 0.0, nil
	case TypeString:
		val, e := v.AsString()
		if e != nil {
			return defaultValue, e
		}
		return len(val) > 0, nil
	case TypeBlob:
		val, e := v.AsBlob()
		if e != nil {
			return defaultValue, e
		}
		return len(val) > 0, nil
	case TypeTimestamp:
		val, e := v.AsTimestamp()
		if e != nil {
			return defaultValue, e
		}
		return !val.IsZero(), nil
	case TypeArray:
		val, e := v.AsArray()
		if e != nil {
			return defaultValue, e
		}
		return len(val) > 0, nil
	case TypeMap:
		val, e := v.AsMap()
		if e != nil {
			return defaultValue, e
		}
		return len(val) > 0, nil
	default:
		return defaultValue,
			fmt.Errorf("cannot convert %T to bool", v)
	}
}

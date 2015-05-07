package tuple

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// Constants defining the largest/smallest float64 numbers that can
// be converted.
const (
	MaxConvFloat64 = float64(math.MaxInt64)
	MinConvFloat64 = float64(math.MinInt64)
)

// ToBool converts a given Value to a bool, if possible. The conversion
// rules are similar to those in Python:
//
//  * Null: false
//  * Bool: actual boolean value
//  * Int: true if non-zero
//  * Float: true if non-zero and not NaN
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
		return val != 0.0 && !math.IsNaN(val), nil
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

// ToInt converts a given Value to an int64, if possible. The conversion
// rules are as follows:
//
//  * Null: (error)
//  * Bool: 0 if false, 1 if true
//  * Int: actual value
//  * Float: conversion as done by int64(value)
//    (values outside of valid int64 bounds will lead to an error)
//  * String: parsed integer with base 0 as per strconv.ParseInt
//    (values outside of valid int64 bounds will lead to an error)
//  * Blob: (error)
//  * Timestamp: the number of microseconds elapsed since
//    January 1, 1970 UTC.
//  * Array: (error)
//  * Map: (error)
func ToInt(v Value) (int64, error) {
	defaultValue := int64(0)
	switch v.Type() {
	case TypeBool:
		val, e := v.AsBool()
		if e != nil {
			return defaultValue, e
		}
		if val {
			return 1, nil
		}
		return 0, nil
	case TypeInt:
		val, e := v.AsInt()
		if e != nil {
			return defaultValue, e
		}
		return val, nil
	case TypeFloat:
		val, e := v.AsFloat()
		if e != nil {
			return defaultValue, e
		}
		if val >= MinConvFloat64 && val <= MaxConvFloat64 {
			return int64(val), nil
		}
		return defaultValue,
			fmt.Errorf("%v is out of bounds for int64 conversion", val)
	case TypeString:
		val, e := v.AsString()
		if e != nil {
			return defaultValue, e
		}
		return strconv.ParseInt(val, 0, 64)
	case TypeTimestamp:
		val, e := v.AsTimestamp()
		if e != nil {
			return defaultValue, e
		}
		// we use the method below instead of UnixNano()/1000 because
		// the latter may be out of range for some timestamps
		seconds := val.Unix()
		secondsAsMicroseconds := seconds * 1000 * 1000
		microsecondPart := int64(val.Nanosecond() / 1000)
		return secondsAsMicroseconds + microsecondPart, nil
	default:
		return defaultValue,
			fmt.Errorf("cannot convert %T to int64", v)
	}
}

// ToFloat converts a given Value to a float64, if possible. The conversion
// rules are as follows:
//
//  * Null: (error)
//  * Bool: 0.0 if false, 1.0 if true
//  * Int: conversion as done by float64(value)
//  * Float: actual value
//  * String: parsed float as per strconv.ParseFloat
//    (values outside of valid float64 bounds will lead to an error)
//  * Blob: (error)
//  * Timestamp: the number of seconds (not microseconds!) elapsed since
//    January 1, 1970 UTC, with a decimal part
//  * Array: (error)
//  * Map: (error)
func ToFloat(v Value) (float64, error) {
	defaultValue := float64(0)
	switch v.Type() {
	case TypeBool:
		val, e := v.AsBool()
		if e != nil {
			return defaultValue, e
		}
		if val {
			return 1.0, nil
		}
		return 0.0, nil
	case TypeInt:
		val, e := v.AsInt()
		if e != nil {
			return defaultValue, e
		}
		return float64(val), nil
	case TypeFloat:
		val, e := v.AsFloat()
		if e != nil {
			return defaultValue, e
		}
		return val, nil
	case TypeString:
		val, e := v.AsString()
		if e != nil {
			return defaultValue, e
		}
		return strconv.ParseFloat(val, 64)
	case TypeTimestamp:
		val, e := v.AsTimestamp()
		if e != nil {
			return defaultValue, e
		}
		// We want to compute `val.UnixNano()/1e9`, but sometimes `UnixNano()`
		// is not defined, so we switch to `val.Unix() + val.Nanosecond()/1e9`.
		// Note that due to numerical issues, this sometimes yields different
		// results within the range of machine precision.
		return float64(val.Unix()) + float64(val.Nanosecond())/1e9, nil
	default:
		return defaultValue,
			fmt.Errorf("cannot convert %T to float64", v)
	}
}

// ToString converts a given Value to a string. The conversion
// rules are as follows:
//
//  * Null: "null"
//  * Bool, Int, Float, String: Go's "%v" representation
//  * Timestamp: ISO 8601 representation, see time.RFC3339
//  * other: Go's "%#v" representation
func ToString(v Value) (string, error) {
	defaultValue := ""
	switch v.Type() {
	case TypeNull:
		return "null", nil
	case TypeBool, TypeInt, TypeFloat, TypeString:
		return fmt.Sprintf("%v", v), nil
	case TypeTimestamp:
		val, e := v.AsTimestamp()
		if e != nil {
			return defaultValue, e
		}
		return val.Format(time.RFC3339), nil
	default:
		return fmt.Sprintf("%#v", v), nil
	}
}

// ToTime converts a given Value to a time.Time struct, if possible.
// The conversion rules are as follows:
//
//  * Null: zero time (this is *not* the time with Unix time 0!)
//  * Int: Time with the given Unix time in microseconds
//  * Float: Time with the given Unix time in seconds, where the decimal
//    part will be considered as a part of a second
//    (values outside of valid int64 bounds will lead to an error)
//  * String: Time with the given RFC3339/ISO8601 representation
//  * Timestamp: actual time
//  * other: (error)
func ToTime(v Value) (time.Time, error) {
	defaultValue := time.Time{}
	switch v.Type() {
	case TypeNull:
		return defaultValue, nil
	case TypeInt:
		val, e := v.AsInt()
		if e != nil {
			return defaultValue, e
		}
		micro := int64(1000 * 1000)
		return time.Unix(val/micro, (val%micro)*1000), nil
	case TypeFloat:
		val, e := v.AsFloat()
		if e != nil {
			return defaultValue, e
		}
		if val >= MinConvFloat64 && val <= MaxConvFloat64 {
			// say val is 3.7 or -4.6
			integralPart := int64(val)                 // 3 or -4
			decimalPart := val - float64(integralPart) // 0.7 or -0.6
			ns := int64(1e9 * decimalPart)             // nanosecond part
			return time.Unix(integralPart, ns), nil
		}
		return defaultValue,
			fmt.Errorf("%v is out of bounds for int64 conversion", val)
	case TypeString:
		val, e := v.AsString()
		if e != nil {
			return defaultValue, e
		}
		return time.Parse(time.RFC3339, val)
	case TypeTimestamp:
		val, e := v.AsTimestamp()
		if e != nil {
			return defaultValue, e
		}
		return val, nil
	default:
		return defaultValue,
			fmt.Errorf("cannot convert %T to Time", v)
	}
}

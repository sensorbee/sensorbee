package builtin

import (
	"bytes"
	"fmt"
	"math"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sort"
)

// singleParamAggFunc is a template for aggregate functions that
// have exactly one parameter
type singleParamAggFunc struct {
	aggFun func([]data.Value) (data.Value, error)
}

func (f *singleParamAggFunc) Accept(arity int) bool {
	return arity == 1
}

func (f *singleParamAggFunc) IsAggregationParameter(k int) bool {
	return k == 0
}

func (f *singleParamAggFunc) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("function takes exactly one argument")
	}
	arr, err := data.AsArray(args[0])
	if err != nil {
		return nil, fmt.Errorf("function needs array input, not %T", args[0])
	}
	return f.aggFun(arr)
}

// twoParamAggFunc is a template for aggregate functions that
// have exactly two (aggregation) parameters
type twoParamAggFunc struct {
	aggFun func([]data.Value, []data.Value) (data.Value, error)
}

func (f *twoParamAggFunc) Accept(arity int) bool {
	return arity == 1 || arity == 2
}

func (f *twoParamAggFunc) IsAggregationParameter(k int) bool {
	return k == 0 || k == 1
}

func (f *twoParamAggFunc) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("function takes exactly two arguments")
	}
	arr1, err := data.AsArray(args[0])
	if err != nil {
		return nil, fmt.Errorf("function needs array input, not %T", args[0])
	}
	arr2, err := data.AsArray(args[1])
	if err != nil {
		return nil, fmt.Errorf("function needs array input, not %T", args[1])
	}
	return f.aggFun(arr1, arr2)
}

// countFunc is an aggregate function that counts the number
// of non-null values passed in.
//
// It can be used in BQL as `count`.
//
//  Input: anything (aggregated)
//  Return Type: Int
var countFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		// count() is O(n) in the spirit of PostgreSQL
		c := int64(0)
		for _, item := range arr {
			if item.Type() != data.TypeNull {
				c++
			}
		}
		return data.Int(c), nil
	},
}

// arrayAggFunc is an aggregate function that concatenates
// input values (including nulls), into an array.
//
// It can be used in BQL as `array_agg`.
//
//  Input: any (aggregated)
//  Return Type: Array (Null on empty input)
var arrayAggFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		return data.Array(arr), nil
	},
}

// avgFunc is an aggregate function that computes the average
// of all input values. Null values are ignored, non-numeric
// values lead to an error.
//
// It can be used in BQL as `avg`.
//
//  Input: Int or Float (aggregated)
//  Return Type: Float (Null on empty input)
var avgFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		sum := float64(0.0)
		count := int64(0)
		for _, item := range arr {
			if item.Type() == data.TypeInt {
				i, _ := data.AsInt(item)
				sum += float64(i)
				count++
			} else if item.Type() == data.TypeFloat {
				f, _ := data.AsFloat(item)
				sum += f
				count++
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a number",
					item, item)
			}
		}
		if count == 0 {
			// only null inputs
			return data.Null{}, nil
		}
		return data.Float(sum / float64(count)), nil
	},
}

// medianFunc is an aggregate function that computes the median
// of all input values. Null values are ignored, non-numeric
// values lead to an error.
//
// It can be used in BQL as `median`.
//
//  Input: Int or Float (aggregated)
//  Return Type: Float (Null on empty input)
var medianFunc udf.UDF = &singleParamAggFunc{
	// this uses a simple algorithm that first sorts all numbers,
	// then takes the one in the middle, i.e., it has time
	// complexity O(n log n)
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		// collect all non-null numeric values
		floatVals := make([]float64, 0, len(arr))
		for _, item := range arr {
			if item.Type() == data.TypeInt {
				i, _ := data.AsInt(item)
				floatVals = append(floatVals, float64(i))
			} else if item.Type() == data.TypeFloat {
				f, _ := data.AsFloat(item)
				floatVals = append(floatVals, f)
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a number",
					item, item)
			}
		}
		if len(floatVals) == 0 {
			// only null inputs
			return data.Null{}, nil
		}
		// sort the input
		sort.Float64s(floatVals)
		// take the middle element (or the average of two middle elements)
		middle := len(floatVals) / 2
		result := floatVals[middle]
		if len(floatVals)%2 == 0 {
			result = (result + floatVals[middle-1]) / 2
		}
		return data.Float(result), nil
	},
}

// skipping bit_and and bit_or here since they are quite low-level

// boolAndFunc is an aggregate function that returns true if
// all input values are true, false otherwise. Null values are
// ignored, non-boolean values lead to an error.
//
// It can be used in BQL as `bool_and`.
//
//  Input: Bool (aggregated)
//  Return Type: Bool (Null on empty input)
var boolAndFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		result := true
		onlyNulls := true
		for _, item := range arr {
			if item.Type() == data.TypeBool {
				b, _ := data.AsBool(item)
				if !b {
					result = b
					// note that if we break here, we will not notice
					// if there are un-boolable values further below
					// and therefore become dependent on the order
					// of rows, which is not good. therefore we do
					// not break here.
				}
				onlyNulls = false
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a bool",
					item, item)
			}
		}
		if onlyNulls {
			return data.Null{}, nil
		}
		return data.Bool(result), nil
	},
}

// boolOrFunc is an aggregate function that returns true if at least
// one of the input values is true, false otherwise. Null values are
// ignored, non-boolean values lead to an error.
//
// It can be used in BQL as `bool_or`.
//
//  Input: Bool (aggregated)
//  Return Type: Bool (Null on empty input)
var boolOrFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		result := false
		onlyNulls := true
		for _, item := range arr {
			if item.Type() == data.TypeBool {
				b, _ := data.AsBool(item)
				if b {
					result = b
					// note that if we break here, we will not notice
					// if there are un-boolable values further below
					// and therefore become dependent on the order
					// of rows, which is not good. therefore we do
					// not break here.
				}
				onlyNulls = false
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a bool",
					item, item)
			}
		}
		if onlyNulls {
			return data.Null{}, nil
		}
		return data.Bool(result), nil
	},
}

// jsonObjectAggFunc is an aggregate functions that returns
// a JSON object with key/value pairs from the two input
// aggregate parameters. If both key and value are Null, the
// pair is ignored. If only the value is Null, it is still added
// with the corresponding key. If only the key is Null, this is
// an error. If the key is not a string, this is an error.
//
// It can be used in BQL as `json_object_agg`.
//
//  Input: String (aggregated), any (aggregated)
//  Return Type: Map (Null on empty input)
var jsonObjectAggFunc udf.UDF = &twoParamAggFunc{
	aggFun: func(keys []data.Value, values []data.Value) (data.Value, error) {
		if len(keys) == 0 && len(values) == 0 {
			return data.Null{}, nil
		} else if len(keys) != len(values) {
			return nil, fmt.Errorf("inputs must have same length (%d != %d)",
				len(keys), len(values))
		}
		result := make(data.Map, len(keys))
		for idx, key := range keys {
			value := values[idx]
			if key.Type() == data.TypeString {
				s, _ := data.AsString(key)
				if _, exists := result[s]; exists {
					return nil, fmt.Errorf("key '%s' appears multiple times", s)
				}
				result[s] = value
			} else if key.Type() == data.TypeNull && value.Type() == data.TypeNull {
				continue
			} else if key.Type() == data.TypeNull {
				return nil, fmt.Errorf("key is null but value (%s) is not", value)
			} else {
				// DO NOT outdent this block, even if golint recommends it,
				// or we will never reach the next iteration
				return nil, fmt.Errorf("cannot interpret %s (%T) as a string",
					key, key)
			}
		}
		return result, nil
	},
}

// maxFunc is an aggregate function that computes the maximum
// value of all input values. Null values are ignored, non-numeric
// values lead to an error.
//
// It can be used in BQL as `max`.
//
//  Input: Int or Float (aggregated)
//  Return Type: same as maximal input value (Null on empty input)
var maxFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		// deal with the case of leading nulls and only nulls
		firstNonNull := -1
		for i, item := range arr {
			if item.Type() != data.TypeNull {
				firstNonNull = i
				break
			}
		}
		if firstNonNull == -1 {
			return data.Null{}, nil
		}
		// if we have timestamp-shaped data
		if arr[firstNonNull].Type() == data.TypeTimestamp {
			maxTime, _ := data.AsTimestamp(arr[firstNonNull])
			for _, item := range arr[firstNonNull:] {
				if item.Type() == data.TypeTimestamp {
					t, _ := data.AsTimestamp(item)
					if maxTime.Sub(t).Seconds() < 0 {
						maxTime = t
					}
				} else if item.Type() == data.TypeNull {
					continue
				} else {
					return nil, fmt.Errorf("cannot interpret %s (%T) as a timestamp",
						item, item)
				}
			}
			return data.Timestamp(maxTime), nil
		}
		// else: numeric
		maxFloat := -float64(math.MaxFloat64)
		maxInt := int64(math.MinInt64)
		for _, item := range arr[firstNonNull:] {
			if item.Type() == data.TypeInt {
				i, _ := data.AsInt(item)
				if i > maxInt {
					maxInt = i
				}
			} else if item.Type() == data.TypeFloat {
				f, _ := data.AsFloat(item)
				if f > maxFloat {
					maxFloat = f
				}
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a number",
					item, item)
			}
		}
		if float64(maxInt) >= maxFloat {
			return data.Int(maxInt), nil
		}
		return data.Float(maxFloat), nil
	},
}

// minFunc is an aggregate function that computes the minimum
// value of all input values. Null values are ignored, non-numeric
// values lead to an error.
//
// It can be used in BQL as `min`.
//
//  Input: Int or Float (aggregated)
//  Return Type: same as minimal input value (Null on empty input)
var minFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		// deal with the case of leading nulls and only nulls
		firstNonNull := -1
		for i, item := range arr {
			if item.Type() != data.TypeNull {
				firstNonNull = i
				break
			}
		}
		if firstNonNull == -1 {
			return data.Null{}, nil
		}
		// if we have timestamp-shaped data
		if arr[firstNonNull].Type() == data.TypeTimestamp {
			minTime, _ := data.AsTimestamp(arr[firstNonNull])
			for _, item := range arr[firstNonNull:] {
				if item.Type() == data.TypeTimestamp {
					t, _ := data.AsTimestamp(item)
					if minTime.Sub(t).Seconds() > 0 {
						minTime = t
					}
				} else if item.Type() == data.TypeNull {
					continue
				} else {
					return nil, fmt.Errorf("cannot interpret %s (%T) as a timestamp",
						item, item)
				}
			}
			return data.Timestamp(minTime), nil
		}
		// else: numeric
		minFloat := float64(math.MaxFloat64)
		minInt := int64(math.MaxInt64)
		for _, item := range arr[firstNonNull:] {
			if item.Type() == data.TypeInt {
				i, _ := data.AsInt(item)
				if i < minInt {
					minInt = i
				}
			} else if item.Type() == data.TypeFloat {
				f, _ := data.AsFloat(item)
				if f < minFloat {
					minFloat = f
				}
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a number",
					item, item)
			}
		}
		if float64(minInt) <= minFloat {
			return data.Int(minInt), nil
		}
		return data.Float(minFloat), nil
	},
}

type stringAggFuncTmpl struct {
}

func (f *stringAggFuncTmpl) Accept(arity int) bool {
	return arity == 1 || arity == 2
}

func (f *stringAggFuncTmpl) IsAggregationParameter(k int) bool {
	return k == 0
}

func (f *stringAggFuncTmpl) Call(ctx *core.Context, args ...data.Value) (data.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("function takes exactly two arguments")
	}
	arr, err := data.AsArray(args[0])
	if err != nil {
		return nil, fmt.Errorf("function needs array input, not %T", args[0])
	}
	delim, err := data.AsString(args[1])
	if err != nil {
		return nil, fmt.Errorf("function needs string input, not %T", args[1])
	}
	if len(arr) == 0 {
		return data.Null{}, nil
	}
	var buffer bytes.Buffer
	for _, item := range arr {
		if item.Type() == data.TypeString {
			if buffer.Len() > 0 {
				buffer.WriteString(delim)
			}
			s, _ := data.AsString(item)
			buffer.WriteString(s)
		} else if item.Type() == data.TypeNull {
			continue
		} else {
			return nil, fmt.Errorf("cannot interpret %s (%T) as a string",
				item, item)
		}
	}
	return data.String(buffer.String()), nil
}

// stringAggFunc(expr, delim) is an aggregate function that
// concatenates its input values into a string, separated by
// a delimiter. Null values are ignored, non-string values
// lead to an error.
var stringAggFunc udf.UDF = &stringAggFuncTmpl{}

// sumFunc is an aggregate function that computes the sum
// of all input values. Null values are ignored, non-numeric
// values lead to an error.
//
// It can be used in BQL as `sum`.
//
//  Input: Int or Float (aggregated)
//  Return Type: Float if the input contains a Float, Int otherwise
//   (Null on empty input)
var sumFunc udf.UDF = &singleParamAggFunc{
	aggFun: func(arr []data.Value) (data.Value, error) {
		if len(arr) == 0 {
			return data.Null{}, nil
		}
		sum := float64(0.0)
		intSum := int64(0)
		hadFloat := false
		onlyNulls := true
		for _, item := range arr {
			if item.Type() == data.TypeInt {
				i, _ := data.AsInt(item)
				// if intSum overflows here, so be it. maybe later
				// additions will fix the situation again. if we
				// try to detect this here and return an error, we
				// become dependent on the input order of numbers.
				intSum += i
				f := float64(i)
				sum += f
				onlyNulls = false
			} else if item.Type() == data.TypeFloat {
				f, _ := data.AsFloat(item)
				sum += f
				hadFloat = true
				onlyNulls = false
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a number",
					item, item)
			}
		}
		if onlyNulls {
			return data.Null{}, nil
		}
		if !hadFloat {
			// if we had only integers, return the integer sum
			// (this is better than converting the float sum
			// back to int64 because we inherit Go's way of dealing
			// with overflows)
			return data.Int(intSum), nil
		}
		return data.Float(sum), nil
	},
}

// skipping xmlagg here since we have no XML data type

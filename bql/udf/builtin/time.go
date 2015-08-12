package builtin

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type diffUsFuncTmpl struct {
	twoParamFunc
}

func (f *diffUsFuncTmpl) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 2 {
		return nil, fmt.Errorf("function takes exactly two arguments")
	}
	arg1 := args[0]
	arg2 := args[1]
	if arg1.Type() == data.TypeNull || arg2.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg1.Type() == data.TypeTimestamp && arg2.Type() == data.TypeTimestamp {
		t1, _ := data.AsTimestamp(arg1)
		t2, _ := data.AsTimestamp(arg2)
		dur := t2.Sub(t1)
		return data.Int(dur.Nanoseconds() / 1000), nil
	}
	return nil, fmt.Errorf("cannot interpret %s and/or %s as timestamp", arg1, arg2)
}

// diffUsFunc(t1, t2) computes the (signed) temporal distance from
// timestamp t1 to t2 in microseconds.
// See also: time.Time.Sub.
//
// It can be used in BQL as `distance_us`.
//
//  Input: 2 * Timestamp
//  Return Type: Int
var diffUsFunc udf.UDF = &diffUsFuncTmpl{}

// clockTimestampFunc returns the local time (in UTC) as a Timestamp.
// See also: time.Now
//
// It can be used in BQL as `clock_timestamp`.
//
//  Input: None
//  Return Type: Timestamp
var clockTimestampFunc = udf.MustConvertGeneric(func() time.Time {
	return time.Now().In(time.UTC)
})

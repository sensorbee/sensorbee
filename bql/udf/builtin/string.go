package builtin

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"regexp"
	"strings"
)

// singleParamStringFunc is a template for functions that
// have a single string parameter as input.
type singleParamStringFunc struct {
	singleParamFunc
	strFun func(string) data.Value
}

func (f *singleParamStringFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) != 1 {
		return nil, fmt.Errorf("function takes exactly one argument")
	}
	arg := args[0]
	if arg.Type() == data.TypeNull {
		return data.Null{}, nil
	} else if arg.Type() == data.TypeString {
		s, _ := data.AsString(arg)
		return f.strFun(s), nil
	}
	return nil, fmt.Errorf("cannot interpret %s as a string", arg)
}

// twoParamStringFunc is a template for functions that
// have two string parameters as input.
type twoParamStringFunc struct {
	twoParamFunc
	strFun func(string, string) data.Value
}

func (f *twoParamStringFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
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
	} else if arg1.Type() == data.TypeString && arg2.Type() == data.TypeString {
		s1, _ := data.AsString(arg1)
		s2, _ := data.AsString(arg2)
		return f.strFun(s1, s2), nil
	}
	return nil, fmt.Errorf("cannot interpret %s and/or %s as string", arg1, arg2)
}

// bitLengthFunc computes the number of bits in a string.
//
// It can be used in BQL as `bit_length`.
//
//  Input: String
//  Return Type: Int
var bitLengthFunc udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		return data.Int(len(s) * 8)
	},
}

// charLengthFunc computes the number of characters in a string.
//
// It can be used in BQL as `char_length`.
//
//  Input: String
//  Return Type: Int
var charLengthFunc udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		return data.Int(len([]rune(s)))
	},
}

// lowerFunc computes the lowercase version of a string as per
// the Unicode standard.
// See also: strings.ToLower
//
// It can be used in BQL as `lower`.
//
//  Input: String
//  Return Type: String
var lowerFunc udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		return data.String(strings.ToLower(s))
	},
}

// byteLengthFunc computes the number of bytes in a string.
//
// It can be used in BQL as `octet_length`.
//
//  Input: String
//  Return Type: Int
var octetLengthFunc udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		return data.Int(len(s))
	},
}

// upperFunc computes the uppercase version of a string as per
// the Unicode standard.
// See also: strings.ToUpper
//
// It can be used in BQL as `upper`.
//
//  Input: String
//  Return Type: String
var upperFunc udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		return data.String(strings.ToUpper(s))
	},
}

type overlayFuncTmpl struct {
}

func (f *overlayFuncTmpl) Accept(arity int) bool {
	return arity == 3 || arity == 4
}

func (f *overlayFuncTmpl) IsAggregationParameter(k int) bool {
	return false
}

func (f *overlayFuncTmpl) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) < 3 || len(args) > 4 {
		return nil, fmt.Errorf("function takes three or four arguments")
	}
	s := ""
	switch args[0].Type() {
	default:
		return nil, fmt.Errorf("cannot interpret %s as a string", args[0])
	case data.TypeNull:
		return data.Null{}, nil
	case data.TypeString:
		s, _ = data.AsString(args[0])
	}
	repl := ""
	switch args[1].Type() {
	default:
		return nil, fmt.Errorf("cannot interpret %s as a string", args[1])
	case data.TypeNull:
		return data.Null{}, nil
	case data.TypeString:
		repl, _ = data.AsString(args[1])
	}
	from := int64(0)
	switch args[2].Type() {
	default:
		return nil, fmt.Errorf("cannot interpret %s as an integer", args[2])
	case data.TypeNull:
		return data.Null{}, nil
	case data.TypeInt:
		from, _ = data.AsInt(args[2])
	}
	if from < 1 {
		return nil, fmt.Errorf("`from` parameter must be at least 1")
	}
	length := int64(len([]rune(repl)))
	if len(args) == 4 {
		switch args[3].Type() {
		default:
			return nil, fmt.Errorf("cannot interpret %s as an integer", args[3])
		case data.TypeNull:
			return data.Null{}, nil
		case data.TypeInt:
			length, _ = data.AsInt(args[3])
		}
		if length < 0 {
			return nil, fmt.Errorf("`for` parameter must be at least 0")
		}
	}
	sRunes := []rune(s)
	if from > int64(len(sRunes)+1) {
		from = int64(len(sRunes) + 1)
	}
	result := string(sRunes[:from-1]) + repl
	if from+length <= int64(len(sRunes)+1) {
		result += string(sRunes[from-1+length:])
	}
	return data.String(result), nil
}

// overlayFunc(s, repl, from, [for]) replaces `for` characters in `s`
// in `s` with the string `repl`, starting at `from` (1-based counting).
// If `for` is not given, the length of `repl` is used.
// See also: SQL's `overlay(string placing string from int for int)`
//
// It can be used in BQL as `overlay`.
//
//  Input: String, String, Int, [Int]
//  Return Type: String
var overlayFunc udf.UDF = &overlayFuncTmpl{}

// strposFunc(str, substr) returns the index of the first occurence
// of `substr` in `str` (1-based) or 0 if it is not found.
// See also: SQL's `position(substring in string)`
//
// It can be used in BQL as `strpos`.
//
//  Input: 2 * String
//  Return Type: Int
var strposFunc udf.UDF = &twoParamStringFunc{
	strFun: func(str, substr string) data.Value {
		return data.Int(strings.Index(str, substr) + 1)
	},
}

type substringFuncTmpl struct {
	twoParamFunc
}

func (f *substringFuncTmpl) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("function takes two or three arguments")
	}
	// convert first argument to string
	str := ""
	if args[0].Type() == data.TypeNull {
		return data.Null{}, nil
	} else if args[0].Type() == data.TypeString {
		str, _ = data.AsString(args[0])
	} else {
		return nil, fmt.Errorf("cannot interpret %s as a string", args[0])
	}
	// substring(str, regex) or substring(str, from)
	if len(args) == 2 {
		if args[1].Type() == data.TypeNull {
			return data.Null{}, nil
		} else if args[1].Type() == data.TypeString {
			expr, _ := data.AsString(args[1])
			re := regexp.MustCompile(expr)
			return data.String(re.FindString(str)), nil
		} else if args[1].Type() == data.TypeInt {
			from, _ := data.AsInt(args[1])
			if from < 1 {
				return nil, fmt.Errorf("`from` parameter must be at least 1")
			}
			sRunes := []rune(str)
			if from > int64(len(sRunes)+1) {
				from = int64(len(sRunes) + 1)
			}
			return data.String(sRunes[from-1:]), nil
		}
		return nil, fmt.Errorf("cannot interpret %s as string or integer", args[1])
	}
	// if we arrive here, we have 3 arguments: substring(str, from, for)
	from := int64(0)
	if args[1].Type() == data.TypeNull {
		return data.Null{}, nil
	} else if args[1].Type() == data.TypeInt {
		from, _ = data.AsInt(args[1])
	} else {
		return nil, fmt.Errorf("cannot interpret %s as integer", args[1])
	}
	if from < 1 {
		return nil, fmt.Errorf("`from` parameter must be at least 1")
	}
	length := int64(0)
	if args[2].Type() == data.TypeNull {
		return data.Null{}, nil
	} else if args[2].Type() == data.TypeInt {
		length, _ = data.AsInt(args[2])
	} else {
		return nil, fmt.Errorf("cannot interpret %s as integer", args[2])
	}
	if length < 0 {
		return nil, fmt.Errorf("`for` parameter must be at least 0")
	}
	sRunes := []rune(str)
	if from > int64(len(sRunes)+1) {
		from = int64(len(sRunes) + 1)
	}
	maxIdx := length + from
	if maxIdx > int64(len(sRunes)+1) {
		maxIdx = int64(len(sRunes) + 1)
	}
	return data.String(sRunes[from-1 : maxIdx-1]), nil
}

// substringFunc(str, reg) returns the part of `str` matching
// the regular expression `reg`.
// See also: regex.Regex.FindString and SQL's
// `substring(string from pattern)`
//
// substringFunc(str, from, [for]) returns the `for` characters
// of `str` starting from the `from` index (1-based). If `for` is
// not given, everything until the end of `str` is returned.
// See also: SQL's: `substring(string from int for int)`
//
// It can be used in BQL as `substring`.
//
//  Input: String, String  or  String, Int, [Int]
//  Return Type: String
var substringFunc udf.UDF = &substringFuncTmpl{}

// ltrimSpaceFunc removes whitespace (" ", \t, \n, \r) from
// the beginning of a string.
//
// It can be used in BQL as `ltrim`.
//
//  Input: String
//  Return Type: String
var ltrimSpaceFunc udf.UDF = &singleParamStringFunc{
	strFun: func(str string) data.Value {
		return data.String(strings.TrimLeft(str, " \t\n\r"))
	},
}

// ltrimFunc(str, chars) removes the longest string containing
// only characters from `chars` from the beginning of `str`.
// See also: strings.TrimLeft
//
// It can be used in BQL as `ltrim`.
//
//  Input: 2 * String
//  Return Type: String
var ltrimFunc udf.UDF = &twoParamStringFunc{
	strFun: func(str, characters string) data.Value {
		return data.String(strings.TrimLeft(str, characters))
	},
}

// rtrimSpaceFunc removes whitespace (" ", \t, \n, \r) from
// the end of a string.
//
// It can be used in BQL as `rtrim`.
//
//  Input: String
//  Return Type: String
var rtrimSpaceFunc udf.UDF = &singleParamStringFunc{
	strFun: func(str string) data.Value {
		return data.String(strings.TrimRight(str, " \t\n\r"))
	},
}

// rtrimFunc(str, chars) removes the longest string containing
// only characters from `chars` from the end of `str`.
// See also: strings.TrimRight
//
// It can be used in BQL as `rtrim`.
//
//  Input: 2 * String
//  Return Type: String
var rtrimFunc udf.UDF = &twoParamStringFunc{
	strFun: func(str, characters string) data.Value {
		return data.String(strings.TrimRight(str, characters))
	},
}

// btrimSpaceFunc removes whitespace (" ", \t, \n, \r) from
// the beginning and end of a string.
//
// It can be used in BQL as `btrim`.
//
//  Input: String
//  Return Type: String
var btrimSpaceFunc udf.UDF = &singleParamStringFunc{
	strFun: func(str string) data.Value {
		return data.String(strings.Trim(str, " \t\n\r"))
	},
}

// btrimFunc(str, chars) removes the longest string containing
// only characters from `chars` from the beginning and end of `str`.
// See also: strings.Trim
//
// It can be used in BQL as `btrim`.
//
//  Input: 2 * String
//  Return Type: String
var btrimFunc udf.UDF = &twoParamStringFunc{
	strFun: func(str, characters string) data.Value {
		return data.String(strings.Trim(str, characters))
	},
}

type variadicFunc struct {
	minParams int
	varFun    func(args ...data.Value) (data.Value, error)
}

func (f *variadicFunc) Accept(arity int) bool {
	return arity >= f.minParams
}

func (f *variadicFunc) IsAggregationParameter(k int) bool {
	return false
}

func (f *variadicFunc) Call(ctx *core.Context, args ...data.Value) (val data.Value, err error) {
	if len(args) < f.minParams {
		return nil, fmt.Errorf("function takes at least %d parameters", f.minParams)
	}
	return f.varFun(args...)

}

// concatFunc concatenates all strings given as input arguments.
// Null values are ignored (treated like ""), non-string arguments
// lead to an error.
//
// It can be used in BQL as `concat`.
//
//  Input: n * String
//  Return Type: String
var concatFunc udf.UDF = &variadicFunc{
	minParams: 1,
	varFun: func(args ...data.Value) (data.Value, error) {
		var buffer bytes.Buffer

		for _, item := range args {
			if item.Type() == data.TypeString {
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
	},
}

// concatWsFunc concatenates all but the first strings given as input
// arguments. The first argument is used as the separator. Null values
// are ignored, non-string arguments lead to an error.
//
// It can be used in BQL as `concat_ws`.
//
//  Input: (1 + n) * String
//  Return Type: String
var concatWsFunc udf.UDF = &variadicFunc{
	minParams: 2,
	varFun: func(args ...data.Value) (data.Value, error) {
		delim := ""
		if args[0].Type() == data.TypeNull {
			return data.Null{}, nil
		} else if args[0].Type() == data.TypeString {
			delim, _ = data.AsString(args[0])
		} else {
			return nil, fmt.Errorf("cannot interpret %s (%T) as a string",
				args[0], args[0])
		}

		var buffer bytes.Buffer

		for _, item := range args[1:] {
			if item.Type() == data.TypeString {
				s, _ := data.AsString(item)
				if buffer.Len() > 0 {
					buffer.WriteString(delim)
				}
				buffer.WriteString(s)
			} else if item.Type() == data.TypeNull {
				continue
			} else {
				return nil, fmt.Errorf("cannot interpret %s (%T) as a string",
					item, item)
			}
		}

		return data.String(buffer.String()), nil
	},
}

// formatFunc formats all but the first argument according
// to a format string (the first argument).
// See also: fmt.Sprintf
//
// It can be used in BQL as `format`.
//
//  Input: String, Any
//  Return Type: String
var formatFunc udf.UDF = &variadicFunc{
	minParams: 2,
	varFun: func(args ...data.Value) (data.Value, error) {
		fmtStr := ""
		if args[0].Type() == data.TypeNull {
			return data.Null{}, nil
		} else if args[0].Type() == data.TypeString {
			fmtStr, _ = data.AsString(args[0])
		} else {
			return nil, fmt.Errorf("cannot interpret %s (%T) as a string",
				args[0], args[0])
		}

		fmtArgs := make([]interface{}, len(args)-1)
		for pos, item := range args[1:] {
			if item.Type() == data.TypeString {
				// to get rid of the additional quotes added by String(),
				// we have to use a native string, not data.String
				s, _ := data.AsString(item)
				fmtArgs[pos] = s
			} else if item.Type() == data.TypeFloat {
				// the user should decide on float formatting himself,
				// not on our String() function (in particular because
				// that one turns Inf/NaN to "null")
				f, _ := data.AsFloat(item)
				fmtArgs[pos] = f
			} else {
				fmtArgs[pos] = item
			}
		}

		return data.String(fmt.Sprintf(fmtStr, fmtArgs...)), nil
	},
}

// md5Func computes the MD5 checksum of a string in hexadecimal
// format.
// See also: crypto/md5.Sum
//
// It can be used in BQL as `md5`.
//
//  Input: String
//  Return Type: String
var md5Func udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		sum := md5.Sum([]byte(s))
		return data.String(fmt.Sprintf("%x", sum))
	},
}

// sha1Func computes the SHA-1 checksum of a string in hexadecimal
// format.
// See also: crypto/sha1.Sum
//
// It can be used in BQL as `sha1`.
//
//  Input: String
//  Return Type: String
var sha1Func udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		sum := sha1.Sum([]byte(s))
		return data.String(fmt.Sprintf("%x", sum))
	},
}

// sha256Func computes the SHA-256 checksum of a string in hexadecimal
// format.
// See also: crypto/sha256.Sum
//
// It can be used in BQL as `sha256`.
//
//  Input: String
//  Return Type: String
var sha256Func udf.UDF = &singleParamStringFunc{
	strFun: func(s string) data.Value {
		sum := sha256.Sum256([]byte(s))
		return data.String(fmt.Sprintf("%x", sum))
	},
}

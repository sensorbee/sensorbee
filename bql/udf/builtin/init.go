package builtin

import (
	"pfi/sensorbee/sensorbee/bql/udf"
)

func init() {
	// numeric functions
	udf.RegisterGlobalUDF("abs", absFunc)
	udf.RegisterGlobalUDF("cbrt", cbrtFunc)
	udf.RegisterGlobalUDF("ceil", ceilFunc)
	udf.RegisterGlobalUDF("degrees", degreesFunc)
	udf.RegisterGlobalUDF("div", divFunc)
	udf.RegisterGlobalUDF("exp", expFunc)
	udf.RegisterGlobalUDF("floor", floorFunc)
	udf.RegisterGlobalUDF("ln", lnFunc)
	udf.RegisterGlobalUDF("log", &arityDispatcher{
		unary: logFunc, binary: logBaseFunc})
	udf.RegisterGlobalUDF("mod", modFunc)
	udf.RegisterGlobalUDF("pi", piFunc)
	udf.RegisterGlobalUDF("power", powFunc)
	udf.RegisterGlobalUDF("radians", radiansFunc)
	udf.RegisterGlobalUDF("round", roundFunc)
	udf.RegisterGlobalUDF("sign", signFunc)
	udf.RegisterGlobalUDF("sqrt", sqrtFunc)
	udf.RegisterGlobalUDF("trunc", truncFunc)
	udf.RegisterGlobalUDF("width_bucket", widthBucketFunc)
	// random functions
	udf.RegisterGlobalUDF("random", randomFunc)
	udf.RegisterGlobalUDF("setseed", setseedFunc)
	// trigonometric functions
	udf.RegisterGlobalUDF("acos", acosFunc)
	udf.RegisterGlobalUDF("asin", asinFunc)
	udf.RegisterGlobalUDF("atan", atanFunc)
	udf.RegisterGlobalUDF("cos", cosFunc)
	udf.RegisterGlobalUDF("cot", cotFunc)
	udf.RegisterGlobalUDF("sin", sinFunc)
	udf.RegisterGlobalUDF("tan", tanFunc)
	// string functions
	udf.RegisterGlobalUDF("bit_length", bitLengthFunc)
	udf.RegisterGlobalUDF("btrim", &arityDispatcher{
		unary: btrimSpaceFunc, binary: btrimFunc})
	udf.RegisterGlobalUDF("char_length", charLengthFunc)
	udf.RegisterGlobalUDF("concat", concatFunc)
	udf.RegisterGlobalUDF("concat_ws", concatWsFunc)
	udf.RegisterGlobalUDF("format", formatFunc)
	udf.RegisterGlobalUDF("lower", lowerFunc)
	udf.RegisterGlobalUDF("ltrim", &arityDispatcher{
		unary: ltrimSpaceFunc, binary: ltrimFunc})
	udf.RegisterGlobalUDF("md5", md5Func)
	udf.RegisterGlobalUDF("octet_length", octetLengthFunc)
	udf.RegisterGlobalUDF("overlay", &arityDispatcher{
		ternary: overlayFunc, quaternary: overlayFunc})
	udf.RegisterGlobalUDF("rtrim", &arityDispatcher{
		unary: rtrimSpaceFunc, binary: rtrimFunc})
	udf.RegisterGlobalUDF("sha1", sha1Func)
	udf.RegisterGlobalUDF("sha256", sha256Func)
	udf.RegisterGlobalUDF("strpos", strposFunc)
	udf.RegisterGlobalUDF("substring", &arityDispatcher{
		binary: substringFunc, ternary: substringFunc})
	udf.RegisterGlobalUDF("upper", upperFunc)
	// time functions
	udf.RegisterGlobalUDF("distance_us", diffUsFunc)
	// aggregate functions
	udf.RegisterGlobalUDF("array_agg", arrayAggFunc)
	udf.RegisterGlobalUDF("avg", avgFunc)
	udf.RegisterGlobalUDF("count", countFunc)
	udf.RegisterGlobalUDF("bool_and", boolAndFunc)
	udf.RegisterGlobalUDF("bool_or", boolOrFunc)
	udf.RegisterGlobalUDF("json_object_agg", jsonObjectAggFunc)
	udf.RegisterGlobalUDF("max", maxFunc)
	udf.RegisterGlobalUDF("min", minFunc)
	udf.RegisterGlobalUDF("string_agg", stringAggFunc)
	udf.RegisterGlobalUDF("sum", sumFunc)
	// other functions
	udf.RegisterGlobalUDF("coalesce", coalesceFunc)
}

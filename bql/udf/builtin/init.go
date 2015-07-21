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
	udf.RegisterGlobalUDF("log", &unaryBinaryDispatcher{logFunc, logBaseFunc})
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
	// aggregate functions
	udf.RegisterGlobalUDF("count", countFunc)
}

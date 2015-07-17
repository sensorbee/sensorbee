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
}

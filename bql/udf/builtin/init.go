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
}

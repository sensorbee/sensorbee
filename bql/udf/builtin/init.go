package builtin

import (
	"pfi/sensorbee/sensorbee/bql/udf"
)

func init() {
	// numeric functions
	udf.RegisterGlobalUDF("abs", AbsFunc())
	udf.RegisterGlobalUDF("cbrt", CbrtFunc())
	udf.RegisterGlobalUDF("ceil", CeilFunc())
	udf.RegisterGlobalUDF("degrees", DegreesFunc())
	udf.RegisterGlobalUDF("exp", ExpFunc())
	udf.RegisterGlobalUDF("floor", FloorFunc())
}

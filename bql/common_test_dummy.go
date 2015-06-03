package bql_test

import (
	"pfi/sensorbee/sensorbee/core"
)

func newTestContext(config core.Configuration) *core.Context {
	return &core.Context{
		Logger: core.NewConsolePrintLogger(),
		Config: config,
	}
}

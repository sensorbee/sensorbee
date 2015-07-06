package api

import (
	"github.com/gocraft/web"
)

type APIContext struct {
	*Context
}

func SetUpAPIRouter(prefix string, router *web.Router, route func(prefix string, r *web.Router)) {
	root := router.Subrouter(APIContext{}, "/api/v1")

	SetUpTopologiesRouter(prefix, root)

	if route != nil {
		route(prefix, root)
	}
}

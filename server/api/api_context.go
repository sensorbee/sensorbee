package api

import (
	"github.com/gocraft/web"
)

type APIContext struct {
	*Context
}

func SetUpAPIRouter(prefix string, router *web.Router) {
	SetUpAPIRouterWithCustomRoute(prefix, router, nil)
}

func SetUpAPIRouterWithCustomRoute(prefix string, router *web.Router, route func(prefix string, r *web.Router)) {
	root := router.Subrouter(APIContext{}, "/api/v1") // TODO need to set version like hawk??

	SetUpBQLRouter(prefix, root)

	if route != nil {
		route(prefix, root)
	}
}

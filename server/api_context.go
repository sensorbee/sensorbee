package server

import (
	"github.com/gocraft/web"
)

// APIContext is a base context of all API controllers.
type APIContext struct {
	*Context
}

// SetUpAPIRouter sets up a router for APIs with user defined custom route.
// Subrouters needs to have APIContext as their first field.
func SetUpAPIRouter(prefix string, router *web.Router, route func(prefix string, r *web.Router)) {
	root := router.Subrouter(APIContext{}, "/api/v1")

	setUpTopologiesRouter(prefix, root)
	setUpServerStatusRouter(prefix, root)

	if route != nil {
		route(prefix, root)
	}
}

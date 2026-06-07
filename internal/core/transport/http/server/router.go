package core_http_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("1")
	ApiVersion2 = ApiVersion("2")
	ApiVersion3 = ApiVersion("3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewAPIVersionRouter(version ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux: http.NewServeMux(),
		apiVersion: version,
	}
}

func (v *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		v.Handle(pattern, route.Handler)
	}
}